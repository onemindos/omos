#!/usr/bin/env python3
"""Portable local/Linux release gates; Git owns intent and Argo owns reconciliation.

Never apply/patch/delete workloads, change Argo policy, or merge/bypass a PR.
Only checked-in, reviewed profiles are supported. Receipts/logs stay outside Git.
"""

import argparse
from contextlib import contextmanager
from datetime import datetime, timezone
import fcntl
import json
import os
from pathlib import Path
import re
import shutil
import subprocess
import sys
import tempfile
import time

import yaml

SHA = re.compile(r"[a-f0-9]{40}")
DIGEST = re.compile(r"sha256:[a-f0-9]{64}")
RELEASE_GATES = {"immutable-yarn", "typescript", "identity-permission-token-tests", "techdocs", "backend-build"}


class GateError(Exception):
    pass


def require(condition, reason):
    if not condition:
        raise GateError(reason)


def utc():
    return datetime.now(timezone.utc).isoformat()


def write_json(path, value):
    with open(path, "x", encoding="utf8", opener=lambda p, f: os.open(p, f, 0o600)) as output:
        json.dump(value, output, indent=2)
        output.write("\n")


@contextmanager
def repo_lock(root):
    common = subprocess.check_output(["git", "rev-parse", "--git-common-dir"], cwd=root, text=True).strip()
    with open((root / common).resolve() / "onemind-release.lock", "a", encoding="utf8") as lock:
        try:
            fcntl.flock(lock, fcntl.LOCK_EX | fcntl.LOCK_NB)
        except BlockingIOError:
            raise GateError("Another pipeline command holds the repository lock") from None
        try:
            yield
        finally:
            fcntl.flock(lock, fcntl.LOCK_UN)


class Run:
    def __init__(self, root, output=None):
        self.root = root
        if output:
            self.output = Path(output).resolve()
            require(not self.output.is_relative_to(root), "Run artifacts must be outside the repository")
            self.output.mkdir(mode=0o700, parents=True, exist_ok=False)
        else:
            self.output = Path(tempfile.mkdtemp(prefix="onemind-gitops-"))
        self.steps = []

    def command(self, label, args, *, env=None, required=True, timeout=600, cwd=None):
        log = self.output / (str(len(self.steps) + 1).zfill(2) + "-" + label + ".log")
        with open(log, "xb", opener=lambda p, f: os.open(p, f, 0o600)) as output:
            started = utc()
            try:
                proc = subprocess.run(args, cwd=cwd or self.root, env=env, stdout=output,
                                      stderr=subprocess.STDOUT, timeout=timeout)
                code = proc.returncode
            except (OSError, subprocess.TimeoutExpired):
                code = -1
        self.steps.append({"name": label, "startedAt": started, "finishedAt": utc(),
                           "passed": code == 0, "required": required, "exitCode": code, "log": log.name})
        print(f"{'PASS' if code == 0 else 'FAIL'} {label}", flush=True)
        if required:
            require(code == 0, f"Gate {label} failed; inspect private run log {log}")
        return log

    def git(self, *args):
        try:
            return subprocess.check_output(["git", *args], cwd=self.root, text=True,
                                           stderr=subprocess.DEVNULL, timeout=120).strip()
        except (subprocess.SubprocessError, OSError):
            raise GateError("Git precondition failed") from None

    def clean(self):
        require(not self.git("status", "--porcelain", "--untracked-files=normal"),
                "Commit or preserve local changes before running the pipeline")

    def identity(self, p):
        require(self.git("remote", "get-url", "origin") in (
            "git@github.com:onemindos/omos.git", p["repository"], p["repository"] + ".git"),
            "Origin is not the approved GitOps repository")

    def receipt(self, phase, **extra):
        value = {"schemaVersion": "onemind.gitops-run.v1", "phase": phase, "observedAt": utc(),
                 "commit": self.git("rev-parse", "HEAD"), "githubChecks": "not_asserted_by_local_executor",
                 "steps": self.steps, **extra}
        write_json(self.output / "run.json", value)
        print("Run evidence: " + str(self.output), flush=True)
        return value


def profile(root, name):
    profiles = json.loads((root / "scripts/gitops/profiles.json").read_text())
    require(name in profiles, "Unknown profile; onboard it through a reviewed PR")
    p = profiles[name]
    for key in ("manifest", "gitopsPath", "applicationFile", "sourcePath", "buildWorkflow"):
        require((root / p[key]).resolve().is_relative_to(root), "Profile path escapes repository")
    require(p["repository"] == "https://github.com/onemindos/omos", "Unapproved repository")
    return p


def quality(run, base, p):
    run.clean()
    base = run.git("rev-parse", base)
    require(SHA.fullmatch(base), "Invalid base commit")
    env = {**os.environ, "YAML_BASE_SHA": base}
    run.command("diff-check", ["git", "diff", "--check", base, "HEAD"])
    for script in ("validate_changed_yaml", "validate_manifests", "validate_helm_renders"):
        run.command(script, [sys.executable, f"scripts/{script}.py"], env=env)
    shells = run.git("ls-files", "*.sh").splitlines()
    shells = [s for s in shells if not s.startswith("vault/")]
    if shells:
        run.command("shellcheck", ["shellcheck", "--severity=error", "-x", "-s", "bash", *shells])
    run.command("terraform-format", ["terraform", "-chdir=infra", "fmt", "-check", "-recursive"])
    run.command("terraform-init-no-backend", ["terraform", "-chdir=infra", "init", "-backend=false", "-input=false", "-lockfile=readonly"])
    run.command("terraform-validate", ["terraform", "-chdir=infra", "validate"])
    run.command("pipeline-tests", [sys.executable, "-m", "unittest", "discover", "-s", "scripts/tests", "-p", "test_gitops*.py", "-v"])
    run.command("image-release-contracts", [sys.executable, "services/onemind-portal/scripts/test-release-image.py"])
    if run.git("diff", "--name-only", base, "HEAD", "--", p["sourcePath"]):
        version = subprocess.check_output(["node", "--version"], text=True).strip()
        require(re.match(r"v(22|24)\.", version), "Use supported Node 22/24 in PATH for portal source checks")
        cwd = run.root / p["sourcePath"]
        yarn = ["node", ".yarn/releases/yarn-4.13.0.cjs"]
        run.command("portal-immutable-install", [*yarn, "install", "--immutable"], cwd=cwd)
        run.command("portal-types", [*yarn, "tsc"], cwd=cwd)
        run.command("portal-lint", [*yarn, "lint"], cwd=cwd)
        run.command("portal-security-tests", [*yarn, "test", "--watch=false", "--runInBand",
                    "plugins/plugin-onemind-permissions-backend/src", "plugins/plugin-onemind-auth/src/scopedFetch.test.ts",
                    "plugins/plugin-onemind-auth-backend/src/subject.test.ts", "plugins/plugin-onemind-common-backend/src/evidenceHttp.test.ts", "plugins/plugin-onemind-services-backend/src/prometheusAlerts.test.ts", "plugins/plugin-onemind-services-backend/src/prometheusRoutes.test.ts"], cwd=cwd)
        run.command("portal-layout-tests", [*yarn, "workspace", "app", "test", "--watch=false", "--runInBand",
                    "ComponentOverview.test.tsx", "componentOverviewModule.test.ts", "catalogCompositionModules.test.ts", "prometheusApi.test.ts", "prometheusPlugin.test.ts", "PrometheusAlerts.test.tsx"], cwd=cwd)
        run.command("catalog-audit-tests", ["node", "--test", "scripts/test-audit-catalog.cjs"], cwd=cwd)
    run.command("topology-catalog-tests", ["node", "--test", "scripts/test-topology.cjs"], cwd=run.root / p["sourcePath"])
    run.command("commit-secrets", ["gitleaks", "git", ".", "--config", ".gitleaks.toml", "--redact",
                                  "--no-banner", f"--log-opts={base}..HEAD"])
    # Inspect the exact tracked release tree, not host caches/ignored local credentials.
    with tempfile.TemporaryDirectory(prefix="onemind-tracked-gate-") as tmp:
        archive = Path(tmp) / "tracked.tar"
        with open(archive, "wb") as output:
            subprocess.run(["git", "archive", "HEAD"], cwd=run.root, stdout=output, check=True)
        tree = Path(tmp) / "tree"
        tree.mkdir()
        subprocess.run(["tar", "-xf", str(archive), "-C", str(tree)], check=True)
        # Scan from the archived root so relative paths match the existing
        # baseline; an absolute scan target makes all nine old findings look new.
        run.command("tracked-tree-secrets", ["gitleaks", "dir", ".", "--config", str(run.root / ".gitleaks.toml"),
                    "--baseline-path", str(run.root / ".gitleaks-baseline.json"), "--redact", "--no-banner"], cwd=tree)
    # Historical whole-tree YAML debt is diagnostic, not silently called passing.
    # Changed YAML, shell, Terraform, schemas and secret gates above remain required.
    run.command("whole-tree-diagnostic", ["make", "validate"], required=False)
    run.clean()
    return base


def release_contract(release, p):
    require(isinstance(release, dict), "Invalid release record")
    commit = release.get("commit", "")
    digest = release.get("digest", "")
    require(isinstance(commit, str) and SHA.fullmatch(commit), "Invalid release commit")
    require(isinstance(digest, str) and DIGEST.fullmatch(digest), "Invalid published digest")
    require(release.get("published") is True and release.get("platform") == p["platform"], "Not a published target-platform image")
    require(release.get("image") == p["imageRepository"] + ":sha-" + commit, "Unexpected registry/tag")
    require(RELEASE_GATES.issubset(set(release.get("gates", []))), "Missing Linux release gates")
    return p["imageRepository"] + "@" + digest


def artifact(run, release_file, p):
    release = json.loads(Path(release_file).read_text())
    image = release_contract(release, p)
    config_log = run.command("registry-source-platform", ["docker", "buildx", "imagetools", "inspect", image,
                                                           "--format", "{{json .Image}}"])
    config = json.loads(config_log.read_text())
    os_name, architecture = p["platform"].split("/")
    require(config.get("os") == os_name and config.get("architecture") == architecture, "Registry platform mismatch")
    labels = config.get("config", {}).get("Labels", {})
    require(labels.get("org.opencontainers.image.revision") == release["commit"], "Registry source revision mismatch")
    require(labels.get("org.opencontainers.image.source") == p["repository"], "Registry source repository mismatch")
    raw = run.command("registry-runtime-digests", ["docker", "buildx", "imagetools", "inspect", "--raw", image])
    manifest = json.loads(raw.read_text())
    allowed = {release["digest"]}
    if "manifests" in manifest:
        runtime = [m for m in manifest["manifests"] if m.get("platform", {}).get("os") == os_name
                   and m.get("platform", {}).get("architecture") == architecture]
        require(len(runtime) == 1 and DIGEST.fullmatch(runtime[0].get("digest", "")), "No unique runtime platform manifest")
        allowed.add(runtime[0]["digest"])
    release["allowedRuntimeDigests"] = sorted(allowed)
    require(run.git("merge-base", release["commit"], "HEAD") == release["commit"], "Release source is not in this checkout's history")
    require(not run.git("diff", "--name-only", release["commit"], "HEAD", "--", p["sourcePath"], p["buildWorkflow"]),
            "Source changed since this image was built; rebuild instead of reusing it")
    return release, image


def changed_image(text, p, image):
    resources = list(yaml.safe_load_all(text))
    matches = [d for d in resources if d and d.get("kind") == "Deployment"
               and d.get("metadata", {}).get("name") == p["deployment"]
               and d.get("metadata", {}).get("namespace") == p["namespace"]]
    require(len(matches) == 1, "No unique profile Deployment in manifest")
    containers = [c for c in matches[0]["spec"]["template"]["spec"]["containers"] if c["name"] == p["container"]]
    require(len(containers) == 1, "No unique profile container")
    old = containers[0]["image"]
    require(old.startswith(p["imageRepository"] + "@sha256:"), "Unexpected existing image repository")
    require(old != image, "This digest is already desired; no promotion needed")
    require(text.count(old) == 1, "Ambiguous image replacement")
    return text.replace(old, image), old


def promote(run, p, args):
    run.clean()
    require(re.fullmatch(r"codex/[a-z0-9][a-z0-9/-]{0,100}", args.branch or ""), "Use an explicit codex/ promotion branch")
    run.command("fetch-main", ["git", "fetch", "origin", "main"])
    release, image = artifact(run, args.release, p)
    require(run.git("merge-base", release["commit"], "origin/main") == release["commit"], "Merge accepted source before promotion")
    require(not run.git("diff", "--name-only", release["commit"], "origin/main", "--", p["sourcePath"], p["buildWorkflow"]),
            "Main has different source; rebuild it")
    run.command("promotion-branch", ["git", "switch", "-c", args.branch, "origin/main"])
    path = run.root / p["manifest"]
    text, previous = changed_image(path.read_text(), p, image)
    path.write_text(text)
    run.command("stage-pin", ["git", "add", "--", p["manifest"]])
    run.command("commit-pin", ["git", "commit", "-m", f"deploy({args.profile}): pin verified {release['commit'][:12]} image"])
    quality(run, "origin/main", p)
    run.command("push-pin", ["git", "push", "-u", "origin", args.branch])
    body = run.output / "pr-body.md"
    body.write_text(f"## Verified GitOps promotion\nSource: `{release['commit']}`\n\nImage: `{image}`\n\n"
                    f"Previous: `{previous}`\n\nApplication: `{p['application']}`. Argo owns rollout; no imperative workload change.\n\n"
                    "Required local gates and registry source/platform checks passed; full-tree diagnostic debt stays explicit. "
                    "GitHub-hosted CI is unavailable, not green. Human review/normal merge is still required. "
                    "After merge use pipeline verify, then browser/semantic/security acceptance. "
                    "Rollback is a reviewed Git revert, never a deployment patch or automatic data rollback.\n")
    log = run.command("open-promotion-pr", ["gh", "pr", "create", "--base", "main", "--head", args.branch,
                      "--title", f"deploy({args.profile}): verified immutable image promotion", "--body-file", str(body)])
    url = log.read_text().strip()
    require(re.fullmatch(r"https://github.com/onemindos/omos/pull/\d+", url), "PR URL not confirmed; inspect run evidence")
    return {"pullRequest": url, "previousImage": previous, "image": image, "sourceCommit": release["commit"], "approval": "required_before_merge"}


def application_contract(app, p):
    spec = app.get("spec", {})
    source = spec.get("source", {})
    require(source.get("repoURL", "").removesuffix(".git") == p["repository"] and source.get("targetRevision") == "main"
            and source.get("path") == p["gitopsPath"], "Application owner/source differs from approved profile")
    require(spec.get("destination", {}).get("namespace") == p["namespace"]
            and spec.get("destination", {}).get("server") == "https://kubernetes.default.svc", "Application destination mismatch")


def rollout_state(app, deployment, pods, p, image, revision, allowed_digests=None):
    application_contract(app, p)
    status = app.get("status", {})
    if (status.get("sync", {}).get("status") != "Synced" or status.get("health", {}).get("status") != "Healthy"
            or status.get("sync", {}).get("revision") != revision
            or status.get("operationState", {}).get("phase") not in (None, "Succeeded")):
        return False
    spec, observed = deployment.get("spec", {}), deployment.get("status", {})
    desired = spec.get("replicas", 1)
    if desired != p["replicas"] or desired <= 0:
        return False
    if observed.get("observedGeneration", -1) < deployment.get("metadata", {}).get("generation", 0):
        return False
    if any(observed.get(k, 0) != desired for k in ("replicas", "readyReplicas", "updatedReplicas", "availableReplicas")):
        return False
    containers = spec.get("template", {}).get("spec", {}).get("containers", [])
    if [c.get("image") for c in containers if c.get("name") == p["container"]] != [image]:
        return False
    active = [pod for pod in pods.get("items", []) if not pod.get("metadata", {}).get("deletionTimestamp")]
    if len(active) != desired:
        return False
    for pod in active:
        if not any(c.get("type") == "Ready" and c.get("status") == "True" for c in pod.get("status", {}).get("conditions", [])):
            return False
        named = [c for c in pod.get("spec", {}).get("containers", []) if c.get("name") == p["container"]]
        states = [c for c in pod.get("status", {}).get("containerStatuses", []) if c.get("name") == p["container"]]
        if len(named) != 1 or named[0].get("image") != image or len(states) != 1 or not states[0].get("ready") or not states[0].get("imageID"):
            return False
        actual_digest = states[0]["imageID"].rsplit("@", 1)[-1].removeprefix("containerd://")
        if actual_digest not in (allowed_digests or {image.rsplit("@", 1)[-1]}):
            return False
    return True


def kube_get(run, p, *args):
    # Fixed read-only verb/context; kubectl never receives credentials in argv.
    try:
        out = subprocess.check_output(["kubectl", "--context", p["context"], "get", *args, "-o", "json"],
                                      cwd=run.root, stderr=subprocess.DEVNULL, timeout=30)
        return json.loads(out)
    except (subprocess.SubprocessError, OSError, ValueError):
        raise GateError("Read-only Kubernetes observation failed") from None


def verify(run, p, args):
    require(SHA.fullmatch(args.revision or ""), "Specify the actual merged GitOps commit")
    run.command("fetch-verification-main", ["git", "fetch", "origin", "main"])
    release, image = artifact(run, args.release, p)
    deadline = time.monotonic() + args.timeout
    while True:
        app = kube_get(run, p, "application", p["application"], "-n", p["argoNamespace"])
        deploy = kube_get(run, p, "deployment", p["deployment"], "-n", p["namespace"])
        pods = kube_get(run, p, "pods", "-n", p["namespace"], "-l", p["podSelector"])
        actual = app.get("status", {}).get("sync", {}).get("revision")
        if actual and actual != args.revision and SHA.fullmatch(actual):
            # Later unrelated Git commits are fine only when the entire owned
            # GitOps path and Application are unchanged. Never approve a superseded rollout.
            if run.git("merge-base", args.revision, actual) == args.revision:
                require(not run.git("diff", "--name-only", args.revision, actual, "--", p["gitopsPath"], p["applicationFile"]),
                        "A newer commit superseded this rollout; review the newer desired state")
                accepted = actual
            else:
                accepted = args.revision
        else:
            accepted = args.revision
        if rollout_state(app, deploy, pods, p, image, accepted, set(release["allowedRuntimeDigests"])):
            write_json(run.output / "application.json", app)
            write_json(run.output / "deployment.json", deploy)
            write_json(run.output / "pods.json", pods)
            return {"image": image, "expectedRevision": args.revision, "observedRevision": accepted,
                    "gitopsRollout": "verified", "semanticBrowserSecurityAcceptance": "still_required"}
        require(time.monotonic() < deadline, "Rollout timeout; inspect Argo/events, do not patch, force sync or auto-rollback")
        time.sleep(min(15, max(0, deadline - time.monotonic())))


def retirement_plan(run, p):
    app = kube_get(run, p, "application", p["application"], "-n", p["argoNamespace"])
    application_contract(app, p)
    resources = app.get("status", {}).get("resources", [])
    require(resources, "Resource inventory unavailable; never infer an empty application")
    claims = kube_get(run, p, "pvc", "-n", p["namespace"])
    return {"application": p["application"], "protectedProfile": p["protected"], "executionPermitted": False,
            "resources": [{k: r.get(k) for k in ("group", "kind", "namespace", "name")} for r in resources],
            "namespacePVCsToRetain": [r["metadata"]["name"] for r in claims.get("items", [])],
            "requiredBeforeRetirement": ["dependency/catalog review", "fresh app-consistent backup and isolated restore",
                "explicit resource-level approval", "reviewed Git PR with retention policy", "Argo prune/finalizer review", "rollback window"],
            "warning": "Plan only. No delete, scale, prune approval, finalizer removal, or namespace/PVC deletion is implemented."}


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("command", choices=("check", "publish", "artifact", "promote", "verify", "retirement-plan"))
    parser.add_argument("--profile", default="portal")
    parser.add_argument("--base", default="origin/main")
    parser.add_argument("--release", type=Path)
    parser.add_argument("--branch")
    parser.add_argument("--revision")
    parser.add_argument("--timeout", type=int, default=600)
    parser.add_argument("--output", type=Path)
    args = parser.parse_args()
    if args.command in ("artifact", "promote", "verify") and not args.release:
        parser.error("--release is required for this command")
    run = None
    try:
        require(1 <= args.timeout <= 1800, "Timeout must be 1–1800 seconds")
        root = Path(subprocess.check_output(["git", "rev-parse", "--show-toplevel"], text=True).strip()).resolve()
        run = Run(root, args.output)
        p = profile(root, args.profile)
        with repo_lock(root):
            run.clean()
            run.identity(p)
            if args.command in ("check", "publish"):
                quality(run, args.base, p)
                result = {"requiredLocalGates": "passed"}
                if args.command == "publish":
                    require(shutil.disk_usage(root).free >= 8 * 1024**3, "Less than 8GiB host free; choose capacity before building")
                    build = run.output / "image"
                    run.command("linux-release", ["bash", f"{p['sourcePath']}/scripts/release-image.sh", "--push", "--output-dir", str(build)], timeout=3600)
                    _, image = artifact(run, build / "release.json", p)
                    result.update({"release": str(build / "release.json"), "image": image, "deployed": False})
            elif args.command == "artifact":
                _, image = artifact(run, args.release, p)
                result = {"image": image, "registrySourcePlatform": "verified", "deployed": False}
            elif args.command == "promote":
                result = promote(run, p, args)
            elif args.command == "verify":
                result = verify(run, p, args)
            else:
                result = retirement_plan(run, p)
            run.receipt(args.command + "_passed", **result)
        return 0
    except (GateError, OSError, ValueError, TypeError, KeyError, subprocess.SubprocessError) as exc:
        # Raw provider/command exception strings can include sensitive context.
        message = str(exc) if isinstance(exc, GateError) else "Pipeline input/tool failure; inspect private logs"
        print("STOP: " + message, file=sys.stderr)
        if run:
            run.receipt("failed", reason=message)
        return 1


if __name__ == "__main__":
    sys.exit(main())
