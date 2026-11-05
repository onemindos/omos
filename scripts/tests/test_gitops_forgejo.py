"""Backup failure recovery and deployment/credential boundaries, no live API calls."""
import importlib.util
from pathlib import Path
import shutil
import subprocess
import unittest

import yaml

ROOT = Path(__file__).resolve().parents[2]
APP = ROOT / "platform/forgejo"


def load(name, path):
    spec = importlib.util.spec_from_file_location(name, path)
    module = importlib.util.module_from_spec(spec)
    spec.loader.exec_module(module)
    return module


BACKUP = load("forgejo_backup", APP / "scripts/coordinated_backup.py")
SECRETS = load("forgejo_secrets", ROOT / "scripts/prepare_forgejo_secrets.py")


class FakeAPI:
    def __init__(self, replicas=1, backup_fails=False, paused=False):
        self.replicas = replicas
        self.calls = []
        self.marker = None
        self.backup_fails = backup_fails
        self.paused = paused

    def __call__(self, path, patch=None):
        self.calls.append((path, patch))
        if path == BACKUP.SCALE:
            if patch:
                self.replicas = patch["spec"]["replicas"]
            return {"metadata": {"resourceVersion": "42"}, "spec": {"replicas": self.replicas}}
        if path == BACKUP.PODS:
            return {"items": []}
        if path == BACKUP.SOURCE:
            if patch:
                self.marker = patch["spec"]["trigger"]["manual"]
                if self.backup_fails:
                    raise RuntimeError("fixture-only API failure")
            return {"spec": {"paused": self.paused}, "status": {"lastManualSync": self.marker}}
        raise AssertionError("Unscoped API endpoint: " + path)


def immediate_wait(check, timeout):
    if not check():
        raise TimeoutError("fixture-only timeout")


class ForgejoBackup(unittest.TestCase):
    def test_success_quiesces_before_trigger_and_restores_prior_replicas(self):
        api = FakeAPI()
        BACKUP.coordinate(api, immediate_wait)
        mutations = [(p, body) for p, body in api.calls if body]
        self.assertEqual([p for p, _ in mutations], [BACKUP.SCALE, BACKUP.SOURCE, BACKUP.SCALE])
        self.assertEqual(api.replicas, 1)
        self.assertEqual(mutations[0][1]["spec"]["replicas"], 0)
        self.assertIn("resourceVersion", mutations[0][1]["metadata"])

    def test_backup_failure_restores_writer_and_never_claims_success(self):
        api = FakeAPI(backup_fails=True)
        with self.assertRaises(RuntimeError):
            BACKUP.coordinate(api, immediate_wait)
        self.assertEqual(api.replicas, 1)

    def test_wait_timeout_restores_writer(self):
        api = FakeAPI()
        def timeout(check, seconds):
            raise TimeoutError("fixture-only")
        with self.assertRaises(TimeoutError):
            BACKUP.coordinate(api, timeout)
        self.assertEqual(api.replicas, 1)
        self.assertIsNone(api.marker)

    def test_lost_pause_response_still_restores_writer(self):
        api = FakeAPI()
        def lost_response(path, patch=None):
            result = api(path, patch)
            if path == BACKUP.SCALE and patch and patch["spec"]["replicas"] == 0:
                raise ConnectionError("fixture-only response lost after server commit")
            return result
        with self.assertRaises(ConnectionError):
            BACKUP.coordinate(lost_response, immediate_wait)
        self.assertEqual(api.replicas, 1)

    def test_intentionally_stopped_writer_is_not_started(self):
        api = FakeAPI(replicas=0)
        BACKUP.coordinate(api, immediate_wait)
        self.assertEqual(api.replicas, 0)
        self.assertFalse(any(p == BACKUP.SCALE and patch for p, patch in api.calls))

    def test_paused_source_and_unexpected_replica_count_fail_before_scaling(self):
        for api in [FakeAPI(paused=True), FakeAPI(replicas=2)]:
            with self.assertRaises(ValueError):
                BACKUP.coordinate(api, immediate_wait)
            self.assertFalse(any(patch for _, patch in api.calls))

    def test_external_replica_change_is_preserved(self):
        api = FakeAPI()
        def external_change(check, seconds):
            self.assertTrue(check())
            if api.marker:
                api.replicas = 2
        BACKUP.coordinate(api, external_change)
        self.assertEqual(api.replicas, 2)

    def test_single_pvc_and_scoped_rbac_no_agent_admin_or_public_registration(self):
        docs = [d for p in APP.glob("*.yaml") if p.name != "kustomization.yaml"
                for d in yaml.safe_load_all(p.read_text()) if d]
        deploy = next(d for d in docs if d["kind"] == "Deployment")
        pod = deploy["spec"]["template"]["spec"]
        self.assertFalse(pod["automountServiceAccountToken"])
        self.assertEqual(deploy["spec"]["strategy"]["type"], "Recreate")
        env = {e["name"]: e["value"] for e in pod["containers"][0]["env"]}
        self.assertEqual(env["FORGEJO__actions__ENABLED"], "false")
        self.assertEqual(env["FORGEJO__service__DISABLE_REGISTRATION"], "true")
        self.assertEqual(env["FORGEJO__server__DISABLE_SSH"], "true")
        claims = [d for d in docs if d["kind"] == "PersistentVolumeClaim"]
        self.assertEqual(len(claims), 1)
        self.assertEqual(claims[0]["spec"]["storageClassName"], "do-block-storage-retain")
        source = next(d for d in docs if d["kind"] == "ReplicationSource")
        self.assertFalse(source["spec"]["trigger"])
        self.assertEqual(source["spec"]["restic"]["copyMethod"], "Snapshot")
        for key in ["storageClassName", "cacheStorageClassName"]:
            self.assertEqual(source["spec"]["restic"][key], "do-block-storage")
        role = next(d for d in docs if d["kind"] == "Role")
        for rule in role["rules"]:
            self.assertNotIn("*", rule["verbs"])
            self.assertNotIn("secrets", rule["resources"])
            if "patch" in rule["verbs"]:
                self.assertEqual(len(rule["resourceNames"]), 1)
        app = yaml.safe_load((ROOT / "platform/argocd/platform-forgejo.yaml").read_text())
        self.assertNotIn("automated", app["spec"]["syncPolicy"])

    def test_secret_payload_uses_separate_repository_and_nonreused_keys(self):
        good = {"RESTIC_REPOSITORY": "s3:https://storage.invalid/bucket/backups/restic/forgejo/data",
                "RESTIC_PASSWORD": "fixture-only-" * 4,
                "AWS_ACCESS_KEY_ID": "fixture-access", "AWS_SECRET_ACCESS_KEY": "fixture-secret"}
        data = SECRETS.payloads(good)
        self.assertEqual(len(data), 2)
        values = list(data["forgejo-keys"].values())
        self.assertEqual(len(values), len(set(values)))
        with self.assertRaises(ValueError):
            SECRETS.payloads({**good, "RESTIC_REPOSITORY": "s3:https://storage.invalid/shared"})

    @unittest.skipUnless(shutil.which("kubectl"), "Kustomize render needs kubectl; source/behavior checks always run")
    def test_kustomize_renders_without_plaintext_secrets(self):
        docs = list(yaml.safe_load_all(subprocess.check_output(["kubectl", "kustomize", str(APP)], text=True)))
        self.assertTrue(any(d["kind"] == "Deployment" for d in docs))
        self.assertFalse(any(d["kind"] == "Secret" for d in docs))


if __name__ == "__main__":
    unittest.main()
