"""Guards for the pending Alertmanager receiver (flagship F3).

The receiver depends on a SealedSecret that Zeus creates. Until it exists the AlertmanagerConfig
must not be reachable from anything Argo CD syncs, because a bad Alertmanager config can stop
all alerting.
"""
import os
import re
import shutil
import subprocess
import tempfile
import unittest
from pathlib import Path

import yaml

ROOT = Path(__file__).resolve().parents[2]
MONITORING = ROOT / "platform/monitoring"
PENDING = MONITORING / "pending"
PENDING_FILE = PENDING / "alertmanagerconfig-onemind-backstage.yaml.disabled"
# What Argo CD's directory source treats as manifests.
MANIFEST_SUFFIXES = {".yaml", ".yml", ".json", ".jsonnet"}
SECRET_NAME = "onemind-alert-webhook"
SECRET_MANIFESTS = [
    MONITORING / "onemind-alert-webhook.sealedsecret.yaml",
    ROOT / "platform/backstage/alert-webhook.sealedsecret.yaml",
]
# (path relative to the repo root, namespace, key) for the two SealedSecrets the helper writes.
SECRET_SPECS = [
    ("platform/monitoring/onemind-alert-webhook.sealedsecret.yaml", "monitoring", "token"),
    ("platform/backstage/alert-webhook.sealedsecret.yaml", "backstage", "ONEMIND_ALERT_WEBHOOK_TOKEN"),
]
HELPER = ROOT / "scripts/ops/seal-alert-webhook.sh"
HELPER_NAME = "seal-alert-webhook"


def sealed_secret_problems(path, namespace, key):
    """Why `path` is not the expected SealedSecret (an empty list means it is)."""
    if not path.is_file():
        return [f"{path} is missing"]
    docs = [d for d in yaml.safe_load_all(path.read_text()) if d]
    if len(docs) != 1 or not isinstance(docs[0], dict):
        return [f"{path} must hold exactly one document"]
    doc = docs[0]
    problems = []
    if doc.get("kind") != "SealedSecret" or doc.get("apiVersion") != "bitnami.com/v1alpha1":
        problems.append(f"{path} is not a bitnami.com/v1alpha1 SealedSecret")
    for plain in ("data", "stringData"):
        if plain in doc:
            problems.append(f"{path} has a top-level {plain} block (plaintext Secret?)")
    meta = doc.get("metadata") or {}
    if meta.get("name") != SECRET_NAME or meta.get("namespace") != namespace:
        problems.append(f"{path} must be {namespace}/{SECRET_NAME}")
    spec = doc.get("spec") or {}
    encrypted = spec.get("encryptedData") or {}
    if set(encrypted) != {key}:
        problems.append(f"{path} encryptedData must hold exactly the key {key}")
    for value in encrypted.values():
        if not (isinstance(value, str) and value.startswith("Ag") and len(value) >= 100):
            problems.append(f"{path} value does not look sealed")
    if (spec.get("template") or {}).get("data"):
        problems.append(f"{path} template carries plaintext data")
    return problems


def receiver_is_reachable(root):
    """True when a synced manifest under platform/monitoring defines the onemind-backstage receiver."""
    monitoring = root / "platform/monitoring"
    pending = monitoring / "pending"
    for path in monitoring.rglob("*"):
        if not path.is_file() or path.suffix not in MANIFEST_SUFFIXES or pending in path.parents:
            continue
        for doc in yaml.safe_load_all(path.read_text()):
            if (
                isinstance(doc, dict)
                and doc.get("kind") == "AlertmanagerConfig"
                and doc["metadata"]["name"] == "onemind-backstage"
            ):
                return True
    return False


def enable_problems(root):
    """What is wrong if the receiver is live: it needs both valid SealedSecrets and the wiring."""
    if not receiver_is_reachable(root):
        return []
    problems = []
    for rel, namespace, key in SECRET_SPECS:
        problems += sealed_secret_problems(root / rel, namespace, key)
    kustomization = root / "platform/backstage/kustomization.yaml"
    listed = kustomization.read_text() if kustomization.is_file() else ""
    if "alert-webhook.sealedsecret.yaml" not in listed:
        problems.append("platform/backstage/kustomization.yaml does not list alert-webhook.sealedsecret.yaml")
    return problems


def load():
    return yaml.safe_load(PENDING_FILE.read_text())


class PendingAlertmanagerConfig(unittest.TestCase):
    def test_pending_folder_contains_nothing_argo_would_sync(self):
        offenders = [p for p in PENDING.rglob("*") if p.is_file() and p.suffix in MANIFEST_SUFFIXES]
        self.assertEqual([], offenders)
        self.assertTrue(PENDING_FILE.is_file())
        self.assertTrue((PENDING / "README.md").is_file())

    def test_platform_monitoring_application_still_reads_this_tree_recursively(self):
        app = yaml.safe_load((ROOT / "platform/argocd/platform-monitoring.yaml").read_text())
        source = app["spec"]["source"]
        self.assertEqual("platform/monitoring", source["path"])
        self.assertTrue(source["directory"]["recurse"])

    def test_no_synced_manifest_defines_or_lists_the_receiver(self):
        if all(p.is_file() for p in SECRET_MANIFESTS):
            self.skipTest("secrets exist; the receiver may be enabled")
        synced = [p for p in MONITORING.rglob("*") if p.is_file() and p.suffix in MANIFEST_SUFFIXES]
        for path in synced:
            text = path.read_text()
            for doc in yaml.safe_load_all(text):
                if isinstance(doc, dict) and doc.get("kind") == "AlertmanagerConfig":
                    self.assertNotEqual("onemind-backstage", doc["metadata"]["name"], str(path))
        # Nothing may reference the pending file by name (kustomization resources, Argo includes).
        for path in ROOT.joinpath("platform").rglob("*"):
            if path.is_file() and path.suffix in MANIFEST_SUFFIXES | {".md"} and path.parent != PENDING:
                self.assertNotIn("alertmanagerconfig-onemind-backstage", path.read_text(), str(path))
        # The webhook secret is not wired into the live backstage manifests yet either.
        self.assertNotIn(SECRET_NAME, (ROOT / "platform/backstage/deployment.yaml").read_text())
        self.assertNotIn("alert-webhook", (ROOT / "platform/backstage/kustomization.yaml").read_text())

    def test_receiver_shape_is_safe(self):
        doc = load()
        self.assertEqual("AlertmanagerConfig", doc["kind"])
        self.assertEqual("monitoring", doc["metadata"]["namespace"])
        route = doc["spec"]["route"]
        # continue: true keeps every other route working exactly as before.
        self.assertTrue(route["continue"])
        severity = [m for m in route["matchers"] if m["name"] == "severity"]
        self.assertEqual([{"name": "severity", "matchType": "=~", "value": "critical|warning"}], severity)
        receiver = next(r for r in doc["spec"]["receivers"] if r["name"] == route["receiver"])
        hook = receiver["webhookConfigs"][0]
        self.assertTrue(hook["sendResolved"])
        self.assertEqual(
            "http://backstage.backstage.svc.cluster.local:7007/api/events/http/alertmanager", hook["url"]
        )
        creds = hook["httpConfig"]["authorization"]["credentials"]
        self.assertEqual({"name": SECRET_NAME, "key": "token"}, creds)
        # The token is only ever a Secret reference, never inline.
        self.assertNotIn("urlSecret", hook)
        text = PENDING_FILE.read_text()
        for forbidden in ("Bearer ey", "token:", "password"):
            self.assertNotIn(forbidden, text.replace("key: token", ""))


class EnableGuard(unittest.TestCase):
    """The receiver may only become reachable together with both valid SealedSecrets."""

    def tree(self, receiver_dir="platform/monitoring", secrets=True, plaintext=False, wired=True):
        tmp = Path(tempfile.mkdtemp())
        self.addCleanup(shutil.rmtree, tmp, ignore_errors=True)
        (tmp / "platform/monitoring/pending").mkdir(parents=True)
        (tmp / "platform/backstage").mkdir(parents=True)
        receiver = (tmp / receiver_dir) / "alertmanagerconfig-onemind-backstage.yaml"
        receiver.parent.mkdir(parents=True, exist_ok=True)
        receiver.write_text(
            "apiVersion: monitoring.coreos.com/v1alpha1\nkind: AlertmanagerConfig\n"
            "metadata:\n  name: onemind-backstage\n  namespace: monitoring\n"
        )
        if secrets:
            for rel, namespace, key in SECRET_SPECS:
                body = (
                    "apiVersion: v1\nkind: Secret\nmetadata:\n"
                    f"  name: {SECRET_NAME}\n  namespace: {namespace}\nstringData:\n  {key}: x\n"
                    if plaintext
                    else "apiVersion: bitnami.com/v1alpha1\nkind: SealedSecret\nmetadata:\n"
                    f"  name: {SECRET_NAME}\n  namespace: {namespace}\nspec:\n  encryptedData:\n"
                    f"    {key}: Ag{'A' * 200}\n"
                )
                (tmp / rel).write_text(body)
        (tmp / "platform/backstage/kustomization.yaml").write_text(
            "resources:\n" + ("  - alert-webhook.sealedsecret.yaml\n" if wired else "  - deployment.yaml\n")
        )
        return tmp

    def test_parked_receiver_needs_nothing(self):
        self.assertEqual([], enable_problems(self.tree(receiver_dir="platform/monitoring/pending", secrets=False)))

    def test_reachable_receiver_without_secrets_fails(self):
        problems = enable_problems(self.tree(secrets=False))
        self.assertEqual(2, len([p for p in problems if "is missing" in p]), problems)

    def test_reachable_receiver_with_a_plaintext_secret_fails(self):
        self.assertTrue(enable_problems(self.tree(plaintext=True)))

    def test_reachable_receiver_without_backstage_wiring_fails(self):
        problems = enable_problems(self.tree(wired=False))
        self.assertEqual(1, len(problems), problems)

    def test_reachable_receiver_with_both_sealed_secrets_passes(self):
        self.assertEqual([], enable_problems(self.tree()))

    def test_repository_state_is_consistent(self):
        self.assertEqual([], enable_problems(ROOT))

    def test_secret_manifests_that_exist_are_valid_sealed_secrets(self):
        for rel, namespace, key in SECRET_SPECS:
            if (ROOT / rel).exists():
                self.assertEqual([], sealed_secret_problems(ROOT / rel, namespace, key))


class SealHelper(unittest.TestCase):
    """scripts/ops/seal-alert-webhook.sh is an operator tool, never part of what Argo syncs."""

    def test_no_argo_source_can_reach_the_helper(self):
        def paths(node):
            if isinstance(node, dict):
                for key, value in node.items():
                    if key == "path" and isinstance(value, str):
                        yield value
                    yield from paths(value)
            elif isinstance(node, list):
                for item in node:
                    yield from paths(item)

        checked = 0
        for folder in ("platform/argocd", "platform/ecosystem"):
            for path in (ROOT / folder).rglob("*"):
                if not path.is_file() or path.suffix not in {".yaml", ".yml"}:
                    continue
                try:
                    docs = list(yaml.safe_load_all(path.read_text()))
                except yaml.YAMLError:
                    continue  # Helm templates and values are not Argo Applications
                for doc in docs:
                    if not isinstance(doc, dict) or doc.get("kind") not in {"Application", "ApplicationSet"}:
                        continue
                    checked += 1
                    for value in paths(doc):
                        normal = value.strip().removeprefix("./").strip("/")
                        self.assertNotIn(normal, {"", ".", "scripts", "scripts/ops"}, f"{path}: {value}")
                        self.assertFalse(normal.startswith("scripts/"), f"{path}: {value}")
        self.assertGreater(checked, 10)

    def test_no_synced_manifest_mentions_the_helper(self):
        for path in (ROOT / "platform").rglob("*"):
            if path.is_file() and path.suffix in MANIFEST_SUFFIXES:
                self.assertNotIn(HELPER_NAME, path.read_text(), str(path))

    def test_helper_lives_outside_every_synced_tree(self):
        self.assertTrue(HELPER.is_file())
        self.assertEqual(("scripts", "ops"), HELPER.relative_to(ROOT).parts[:2])
        self.assertTrue(os.access(HELPER, os.X_OK), "helper must be executable")

    def test_helper_is_safe_by_construction(self):
        text = HELPER.read_text()
        self.assertTrue(text.startswith("#!/usr/bin/env bash\n"))
        self.assertIn("\nset -euo pipefail\n", text)
        self.assertNotRegex(text, r"(?m)^\s*set\s+-[a-z]*x")
        self.assertNotRegex(text, r"bash\s+-[a-z]*x")
        # The token is never echoed, written, or put in a command line.
        self.assertNotRegex(text, r"(echo|cat)[^\n|]*\$\{?TOKEN\b")
        self.assertNotRegex(text, r"printf[^\n|]*\$\{?TOKEN\b[^\n|]*>")
        for forbidden in ("--from-literal", "kubectl create secret", "kubectl apply", "kubectl delete", "mktemp -u"):
            self.assertNotIn(forbidden, text)
        # Repo sealing convention (CLAUDE.md) and refusal to overwrite.
        for required in ("--scope namespace-wide", "--format yaml", "--fetch-cert", "noclobber", "refusing to overwrite"):
            self.assertIn(required, text)
        # Dry run exits before any token exists: the generation comes after the --apply gate.
        gate = text.index('if [ "$APPLY" -ne 1 ]')
        generate = text.index('TOKEN="$(openssl rand')
        self.assertLess(gate, generate)
        self.assertIn("exit 0", text[gate:generate])

    def run_dry(self, outputs_exist=False):
        bash = shutil.which("bash")
        openssl = shutil.which("openssl")
        if not (bash and openssl):
            self.skipTest("bash and openssl are required")
        tmp = Path(tempfile.mkdtemp())
        self.addCleanup(shutil.rmtree, tmp, ignore_errors=True)
        (tmp / "scripts/ops").mkdir(parents=True)
        (tmp / "platform/monitoring").mkdir(parents=True)
        (tmp / "platform/backstage").mkdir(parents=True)
        script = tmp / "scripts/ops/seal-alert-webhook.sh"
        shutil.copy(HELPER, script)
        stubs = tmp / "stubs"
        stubs.mkdir()
        # Stand-ins for the cluster tools: they answer only the read-only calls a dry run makes.
        (stubs / "kubectl").write_text(
            '#!/bin/sh\n[ "$1 $2 $3 $4" = "config get-contexts -o name" ] && { echo do-nyc3-omos-cloud; exit 0; }\nexit 9\n'
        )
        (stubs / "kubeseal").write_text(
            '#!/bin/sh\ncase "$*" in *--fetch-cert*) echo "-----BEGIN CERTIFICATE-----"; exit 0;; esac\necho "unexpected kubeseal call" >&2; exit 9\n'
        )
        for stub in stubs.iterdir():
            stub.chmod(0o755)
        if outputs_exist:
            (tmp / SECRET_SPECS[0][0]).write_text("existing\n")
        env = {"PATH": f"{stubs}:{Path(openssl).parent}:/usr/bin:/bin", "HOME": str(tmp)}
        result = subprocess.run(
            [bash, str(script)], env=env, capture_output=True, text=True, cwd=tmp, check=False
        )
        written = sorted(p.name for p in (tmp / "platform").rglob("*") if p.is_file())
        return result, written

    def test_default_run_is_a_dry_run_that_writes_nothing(self):
        result, written = self.run_dry()
        self.assertEqual(0, result.returncode, result.stdout + result.stderr)
        self.assertIn("Dry run complete", result.stdout)
        self.assertEqual([], written)
        self.assertNotRegex(result.stdout + result.stderr, r"[0-9a-f]{64}")

    def test_existing_output_is_refused_and_left_untouched(self):
        result, written = self.run_dry(outputs_exist=True)
        self.assertNotEqual(0, result.returncode)
        self.assertIn("nothing was written", result.stderr)
        self.assertEqual(["onemind-alert-webhook.sealedsecret.yaml"], written)


if __name__ == "__main__":
    unittest.main()
