"""`config-export` (no-agent job): nightly, password-free snapshot of Legacy's config, pushed to Git.

What is snapshotted (an explicit allowlist, not a profile tarball: `hermes profile export` packs
the whole home, including sessions, hub caches and bundled skills, about 12 MB):

  config.yaml          secret-looking keys blanked, ${ENV} references kept
  SOUL.md, memories/MEMORY.md, memories/USER.md
  skills/index.json    every skill: name, category, origin (bundled / user / external)
  skills/user/...      the contents of skills Legacy created or edited (not bundled, not external)
  cron/jobs.json       job definitions only (no run history, no outputs)

Before anything is written to Git the whole snapshot is scanned: known token shapes, private-key
headers, password/token assignments, and every secret value this pod actually holds (the
environment and auth.json). One hit aborts the run with the file and rule name, never the value.

Git write access is a deploy key mounted from the optional Secret `onemind-hermes-export`
(keys: ssh-privatekey, repo, optional branch) at /var/run/legacy-export. Without it the script
fails loudly (a silently skipped backup is worse than an alert); the job stays paused until the
owner creates the Secret. SSH goes to ssh.github.com:443, which the pod's egress already allows.

Usage: legacy_config_export.py [--dry-run] [--out DIR]
  --dry-run  build and scan the snapshot, print a summary, touch no Git
"""

import hashlib
import io
import json
import os
import re
import shutil
import subprocess
import sys
import tempfile
from pathlib import Path

HOME = Path(os.environ.get("HERMES_HOME", "/opt/data"))
EXPORT_DIR = Path(os.environ.get("LEGACY_EXPORT_SECRET_DIR", "/var/run/legacy-export"))
DEST = "legacy"  # directory inside the target repository
EXTERNAL_SKILLS = HOME / "onemind-skills"

# ssh.github.com:443 host key, from https://api.github.com/meta (SHA256:+DiY3wvvV6TuJJhbpZisF/zLDA0zPMSvHdkr4UvCOqU)
KNOWN_HOSTS = "[ssh.github.com]:443 ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIOMqqnkVzrm0SdG6UOoqKLsabgH5C9okWi0dh2l9GKJl\n"

SECRET_KEY = re.compile(r"(token|secret|password|passwd|api[_-]?key|credential|private[_-]?key|auth)", re.I)
ENV_REF = re.compile(r"^\$\{?[A-Z][A-Z0-9_]*\}?$|^(Bearer )?\$\{[A-Z][A-Z0-9_]*\}$")

SCAN_RULES = {
    "openai-style key": re.compile(r"\bsk-[A-Za-z0-9_-]{20,}"),
    "github token": re.compile(r"\b(ghp|gho|ghu|ghs|ghr)_[A-Za-z0-9]{30,}|github_pat_[A-Za-z0-9_]{30,}"),
    "slack token": re.compile(r"\bxox[abprs]-[A-Za-z0-9-]{10,}"),
    "aws access key": re.compile(r"\b(AKIA|ASIA)[0-9A-Z]{16}\b"),
    "private key block": re.compile(r"-----BEGIN [A-Z ]*PRIVATE KEY-----"),
    "jwt": re.compile(r"\beyJ[A-Za-z0-9_-]{10,}\.[A-Za-z0-9_-]{10,}\.[A-Za-z0-9_-]{10,}"),
    "telegram bot token": re.compile(r"\b\d{8,10}:[A-Za-z0-9_-]{35}\b"),
    "credential assignment": re.compile(
        r"(?i)\b(password|passwd|secret|api[_-]?key|access[_-]?token|auth[_-]?token)\b[\"']?\s*[:=]\s*[\"']?"
        r"(?!<redacted>|\$\{|\$[A-Z]|true\b|false\b|null\b|none\b|\*+|your[_-]|changeme|example|placeholder)[A-Za-z0-9_\-+/=.]{12,}"),
}


def die(msg, code=1):
    print("config-export: " + msg)
    sys.exit(code)


def redact(node, key=""):
    if isinstance(node, dict):
        return {k: redact(v, str(k)) for k, v in node.items()}
    if isinstance(node, list):
        return [redact(v, key) for v in node]
    if isinstance(node, str) and node and SECRET_KEY.search(key) and not ENV_REF.match(node.strip()):
        return "<redacted>"
    return node


def build(snap):
    from ruamel.yaml import YAML  # PyYAML is not in Hermes' venv

    n = 0

    def put(rel, data):
        nonlocal n
        path = snap / rel
        path.parent.mkdir(parents=True, exist_ok=True)
        path.write_text(data, encoding="utf-8")
        n += 1

    cfg = HOME / "config.yaml"
    if cfg.is_file():
        yml = YAML(typ="safe")
        yml.default_flow_style = False
        buf = io.StringIO()
        yml.dump(redact(yml.load(cfg.read_text()) or {}), buf)
        put("config.yaml", buf.getvalue())
    for rel in ("SOUL.md", "memories/MEMORY.md", "memories/USER.md"):
        if (HOME / rel).is_file():
            put(rel, (HOME / rel).read_text(encoding="utf-8", errors="replace"))

    bundled = set()
    manifest = HOME / "skills" / ".bundled_manifest"
    if manifest.is_file():
        bundled = {ln.split(":", 1)[0].strip() for ln in manifest.read_text().splitlines() if ln.strip()}
    index = []
    for root, base, origin in ((HOME / "skills", HOME / "skills", None), (EXTERNAL_SKILLS, EXTERNAL_SKILLS, "external")):
        if not root.is_dir():
            continue
        for skill_md in sorted(root.rglob("SKILL.md")):
            sdir = skill_md.parent
            rel = sdir.relative_to(base)
            if ".hub" in rel.parts or ".archive" in rel.parts:
                continue
            kind = origin or ("bundled" if sdir.name in bundled else "user")
            index.append({"name": sdir.name, "category": rel.parts[0] if len(rel.parts) > 1 else "", "origin": kind})
            if kind == "user":
                for f in sorted(sdir.rglob("*")):
                    if f.is_file() and f.stat().st_size <= 200_000 and "__pycache__" not in f.parts:
                        try:
                            put("skills/user/%s/%s" % (rel, f.relative_to(sdir)), f.read_text(encoding="utf-8"))
                        except UnicodeDecodeError:
                            pass  # binary file: not snapshotted
    put("skills/index.json", json.dumps(sorted(index, key=lambda s: (s["origin"], s["name"])), indent=2) + "\n")

    keep = ("name", "schedule_display", "prompt", "skills", "model", "provider", "script", "no_agent", "monitor_script",
            "deliver", "enabled", "state", "paused_reason", "reasoning_effort", "repeat")
    jobs_file = HOME / "cron" / "jobs.json"
    jobs = []
    if jobs_file.is_file():
        raw = json.loads(jobs_file.read_text())
        for job in raw.get("jobs", raw) if isinstance(raw, dict) else raw:
            jobs.append({k: job.get(k) for k in keep if k in job})
    put("cron/jobs.json", json.dumps(sorted(jobs, key=lambda j: j.get("name", "")), indent=2) + "\n")
    return n


def held_secrets():
    """Secret values this pod holds: scanned for verbatim, so a leak is caught whatever its shape."""
    vals = set()
    for k, v in os.environ.items():
        if re.search(r"KEY|TOKEN|SECRET|PASSWORD|PASSWD|CREDENTIAL", k, re.I) and len(v) >= 8:
            vals.add(v)
    auth = HOME / "auth.json"
    if auth.is_file():
        def leaves(x):
            if isinstance(x, dict):
                for v in x.values():
                    yield from leaves(v)
            elif isinstance(x, list):
                for v in x:
                    yield from leaves(v)
            elif isinstance(x, str) and len(x) >= 16 and " " not in x:
                yield x
        try:
            vals.update(leaves(json.loads(auth.read_text())))
        except ValueError:
            pass
    return vals


def scan(snap):
    secrets = held_secrets()
    hits = []
    for f in sorted(p for p in snap.rglob("*") if p.is_file()):
        text = f.read_text(encoding="utf-8", errors="replace")
        rel = str(f.relative_to(snap))
        for name, rx in SCAN_RULES.items():
            if rx.search(text):
                hits.append("%s: %s" % (rel, name))
        if any(s in text for s in secrets):
            hits.append("%s: verbatim secret held by this pod" % rel)
    return hits


def run(cmd, **kw):
    return subprocess.run(cmd, check=True, capture_output=True, text=True, timeout=120, **kw)


def push(snap, files):
    key = EXPORT_DIR / "ssh-privatekey"
    repo_f = EXPORT_DIR / "repo"
    if not key.is_file() or not repo_f.is_file():
        die("no Git write credential: Secret onemind-hermes-export (ssh-privatekey, repo) is not mounted at %s" % EXPORT_DIR)
    repo = repo_f.read_text().strip()
    branch_f = EXPORT_DIR / "branch"
    branch = branch_f.read_text().strip() if branch_f.is_file() else "main"
    work = Path(tempfile.mkdtemp(prefix="legacy-export-"))
    try:
        kh = work / "known_hosts"
        kh.write_text(KNOWN_HOSTS)
        ssh_key = work / "key"
        shutil.copyfile(key, ssh_key)  # ssh refuses a group-readable key; the Secret mount may be 0444
        ssh_key.chmod(0o600)
        env = dict(os.environ, GIT_SSH_COMMAND=(
            "ssh -i %s -o IdentitiesOnly=yes -o StrictHostKeyChecking=yes -o UserKnownHostsFile=%s -o BatchMode=yes" % (ssh_key, kh)),
            GIT_TERMINAL_PROMPT="0")
        clone = work / "repo"
        try:
            run(["git", "clone", "--quiet", "--depth", "1", "--branch", branch, repo, str(clone)], env=env)
        except subprocess.CalledProcessError:
            # empty repository or missing branch: start a fresh history on that branch
            run(["git", "init", "--quiet", "-b", branch, str(clone)], env=env)
            run(["git", "-C", str(clone), "remote", "add", "origin", repo], env=env)
        target = clone / DEST
        shutil.rmtree(target, ignore_errors=True)
        shutil.copytree(snap, target)
        run(["git", "-C", str(clone), "add", "-A", DEST], env=env)
        if not run(["git", "-C", str(clone), "status", "--porcelain", DEST], env=env).stdout.strip():
            return None  # nothing changed: silent
        digest = hashlib.sha256("".join(sorted(files)).encode()).hexdigest()[:8]
        run(["git", "-C", str(clone), "-c", "user.name=Legacy", "-c", "user.email=legacy@onemindos.dev",
             "commit", "--quiet", "-m", "chore(legacy): nightly config snapshot %s (%d files, secret-scanned)" % (digest, len(files))], env=env)
        run(["git", "-C", str(clone), "push", "--quiet", "origin", "HEAD:refs/heads/" + branch], env=env)
        return run(["git", "-C", str(clone), "rev-parse", "--short", "HEAD"], env=env).stdout.strip()
    except subprocess.CalledProcessError as exc:
        die("git step failed (%s): %s" % (" ".join(exc.cmd[3:5]), (exc.stderr or "")[-200:].replace(str(work), "<tmp>")))
    finally:
        shutil.rmtree(work, ignore_errors=True)


def main():
    args = sys.argv[1:]
    dry = "--dry-run" in args
    out = Path(args[args.index("--out") + 1]) if "--out" in args else None
    snap = out or Path(tempfile.mkdtemp(prefix="legacy-snapshot-"))
    snap.mkdir(parents=True, exist_ok=True)
    try:
        count = build(snap)
        hits = scan(snap)
        if hits:
            die("ABORTED, secret scan found %d problem(s); nothing was written to Git:\n  %s" % (len(hits), "\n  ".join(hits)))
        files = [str(p.relative_to(snap)) for p in sorted(snap.rglob("*")) if p.is_file()]
        if dry:
            print("config-export dry run: %d files, secret scan clean" % count)
            for f in files:
                print("  " + f)
            return 0
        sha = push(snap, files)
        if sha:
            print("config-export: snapshot committed (%s), %d files, secret scan clean" % (sha, len(files)))
        return 0
    finally:
        if out is None:
            shutil.rmtree(snap, ignore_errors=True)


if __name__ == "__main__":
    sys.exit(main())
