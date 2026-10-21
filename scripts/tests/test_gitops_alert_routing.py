"""Guards for the pending Alertmanager receiver (flagship F3).

The receiver depends on a SealedSecret that Zeus creates. Until it exists the AlertmanagerConfig
must not be reachable from anything Argo CD syncs, because a bad Alertmanager config can stop
all alerting.
"""
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


if __name__ == "__main__":
    unittest.main()
