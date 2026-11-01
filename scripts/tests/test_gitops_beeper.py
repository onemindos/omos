"""Deployment and credential boundaries for the owner-approved Beeper batch."""
import importlib.util
import json
import shutil
from pathlib import Path
import subprocess
import unittest
from unittest.mock import patch
import yaml

ROOT = Path(__file__).resolve().parents[2]
APP = ROOT / "platform/beeper-bridges"
TYPES = {"whatsapp", "telegram", "gmessages", "instagram", "twitter"}


class BeeperDeploymentContracts(unittest.TestCase):
    @classmethod
    def setUpClass(cls):
        # These contracts inspect authored resources and run without cluster CLIs.
        cls.resources = []
        for path in sorted(APP.glob("*.yaml")):
            if path.name != "kustomization.yaml":
                cls.resources.extend(item for item in yaml.safe_load_all(path.read_text()) if item)
        spec = importlib.util.spec_from_file_location("prepare_beeper", ROOT / "scripts/prepare_beeper_secrets.py")
        cls.helper = importlib.util.module_from_spec(spec)
        spec.loader.exec_module(cls.helper)

    @unittest.skipUnless(shutil.which("kubectl"), "Kustomize render requires kubectl; source contracts always run")
    def test_kustomize_render_contains_all_authored_resources(self):
        rendered = subprocess.check_output(["kubectl", "kustomize", str(APP)], text=True)
        resources = [item for item in yaml.safe_load_all(rendered) if item]
        expected = {(r["kind"], r["metadata"]["name"]) for r in self.resources}
        actual = {(r["kind"], r["metadata"]["name"]) for r in resources}
        self.assertTrue(expected <= actual)
        self.assertFalse(any(r["kind"] == "Secret" for r in resources))

    def test_all_five_deploy_together_without_public_or_cluster_authority(self):
        deployments = [r for r in self.resources if r["kind"] == "Deployment"]
        self.assertEqual({d["metadata"]["name"] for d in deployments}, {"beeper-" + t for t in TYPES})
        self.assertFalse(any(r["kind"] in {"Ingress", "Service", "RoleBinding", "ClusterRoleBinding"} for r in self.resources))
        for deployment in deployments:
            self.assertEqual(deployment["spec"]["replicas"], 1)
            self.assertEqual(deployment["spec"]["strategy"]["type"], "Recreate")
            pod = deployment["spec"]["template"]["spec"]
            self.assertFalse(pod["automountServiceAccountToken"])
            self.assertTrue(pod["securityContext"]["runAsNonRoot"])
            for container in pod["containers"] + pod.get("initContainers", []):
                self.assertIn("@sha256:", container["image"])
                self.assertFalse(container["securityContext"]["allowPrivilegeEscalation"])
                self.assertIn("ALL", container["securityContext"]["capabilities"]["drop"])
                self.assertIn("requests", container["resources"])
                self.assertIn("limits", container["resources"])

    def test_retained_state_has_separate_crash_consistent_backups(self):
        claims = [r for r in self.resources if r["kind"] == "PersistentVolumeClaim"]
        backups = [r for r in self.resources if r["kind"] == "ReplicationSource"]
        self.assertEqual(len(claims), 5)
        self.assertEqual({r["spec"]["sourcePVC"] for r in backups}, {c["metadata"]["name"] for c in claims})
        for claim in claims:
            self.assertEqual(claim["spec"]["storageClassName"], "do-block-storage-retain")
            self.assertIn("Prune=false", claim["metadata"]["annotations"]["argocd.argoproj.io/sync-options"])
        for backup in backups:
            restic = backup["spec"]["restic"]
            self.assertEqual(restic["copyMethod"], "Snapshot")
            self.assertEqual(restic["storageClassName"], "do-block-storage")
            self.assertEqual(restic["cacheStorageClassName"], "do-block-storage")

    def test_batch_activation_is_explicit_and_mcp_content_access_stays_disabled(self):
        argo = yaml.safe_load((ROOT / "platform/argocd/platform-beeper-bridges.yaml").read_text())
        self.assertNotIn("automated", argo["spec"]["syncPolicy"])
        apps = yaml.safe_load((ROOT / "platform/hermes/mcp/apps.yaml").read_text())["apps"]
        beeper = next(a for a in apps if a["id"] == "beeper")
        self.assertFalse(beeper["enabled"])
        self.assertTrue(beeper["blocked_reason"])
        self.assertFalse(any(r["kind"] == "Secret" for r in self.resources))

    def test_connector_lock_is_complete_and_checksums_are_fixed(self):
        lock = json.loads((APP / "artifacts-lock.json").read_text())
        self.assertEqual(set(lock), TYPES)
        for name, entry in lock.items():
            self.assertEqual(entry["binary"], "mautrix-" + name)
            self.assertRegex(entry["url"], r"^https://mau\.dev/mautrix/[^/]+/-/jobs/[0-9]+/artifacts/raw/")
            self.assertRegex(entry["sha256"], r"^[0-9a-f]{64}$")
            self.assertRegex(entry["commit"], r"^[0-9a-f]{40}$")

    def test_private_payloads_require_correct_auth_and_distinct_backup_repositories(self):
        token = "bat_" + "x" * 32
        config = {"environments": {"prod": {"access_token": token}}}
        backup = {"S3_ACCESS_KEY": "fixture-access", "S3_SECRET_KEY": "fixture-secret",
                  "S3_BUCKET": "fixture-bucket", "S3_ENDPOINT": "https://storage.invalid"}
        payloads = self.helper.prepare_payloads(config, "123", "a" * 32, backup)
        self.assertEqual(len(payloads), 7)
        self.assertEqual(payloads["beeper-bridge-auth"]["MATRIX_ACCESS_TOKEN"], token)
        repositories = [p["RESTIC_REPOSITORY"] for p in payloads.values() if "RESTIC_REPOSITORY" in p]
        passwords = [p["RESTIC_PASSWORD"] for p in payloads.values() if "RESTIC_PASSWORD" in p]
        self.assertEqual(len(set(repositories)), 5)
        self.assertEqual(len(set(passwords)), 5)
        with self.assertRaises(ValueError):
            self.helper.prepare_payloads({}, "123", "a" * 32, backup)
        with self.assertRaises(ValueError):
            self.helper.prepare_payloads(config, "not-numeric", "a" * 32, backup)

    def test_sealing_failure_never_echoes_private_subprocess_diagnostics(self):
        result = subprocess.CompletedProcess([], 1, stdout="", stderr="sensitive diagnostic fixture")
        with patch.object(self.helper.subprocess, "run", return_value=result):
            with self.assertRaisesRegex(RuntimeError, "^Sealing failed; no plaintext output was written$"):
                self.helper.seal("test", {"private": "fixture-value"}, Path("/public/cert"))


if __name__ == "__main__":
    unittest.main()
