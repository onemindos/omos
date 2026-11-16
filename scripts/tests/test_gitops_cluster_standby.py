"""cluster_standby.py: parents paused first, originals recorded once, wake restores in reverse order."""
import importlib.util
import json
from pathlib import Path
import unittest
from unittest import mock

ROOT = Path(__file__).resolve().parents[2]
spec = importlib.util.spec_from_file_location("cluster_standby", ROOT / "scripts/ops/cluster_standby.py")
cs = importlib.util.module_from_spec(spec)
spec.loader.exec_module(cs)


def app(name, ns, resources=(), automated=True, annotations=None):
    sp = {"automated": {"selfHeal": True}} if automated else {}
    return {"metadata": {"name": name, "annotations": annotations or {}},
            "spec": {"destination": {"namespace": ns}, "syncPolicy": sp},
            "status": {"resources": [dict(r) for r in resources]}}


def workload(name, replicas, annotations=None):
    return {"metadata": {"name": name, "annotations": annotations or {}}, "spec": {"replicas": replicas}}


class Standby(unittest.TestCase):
    def setUp(self):
        self.calls = []
        patcher = mock.patch.object(cs, "kubectl", side_effect=lambda *a, **k: self.calls.append(a) or "")
        patcher.start()
        self.addCleanup(patcher.stop)

    def test_every_parent_is_paused_before_any_child(self):
        # platform-security owns ecosystem-openfga; alphabetically it sorts after it. The 2026-11-16 dry run
        # showed a parent paused after its child re-enables the child, so parents must go first.
        apps = [
            app("ecosystem-openfga", "onemind-core"),
            app("platform-security", "argocd", [{"kind": "Application", "name": "ecosystem-openfga"}]),
            app("platform-n8n", "n8n"),
            app("platform-onemind-hermes", "hermes"),
        ]
        with mock.patch.object(cs, "kget", return_value=apps):
            names = [a["metadata"]["name"] for a in cs.apps_to_pause(["n8n", "onemind-core"])]
        self.assertEqual(names[0], "platform-security")
        self.assertIn("ecosystem-openfga", names)
        self.assertNotIn("platform-onemind-hermes", names, "Hermes must keep auto-sync")

    def test_reverted_workload_is_scaled_again_without_losing_its_original_count(self):
        reverted = workload("openfga", 1, {cs.A_REPLICAS: "2"})
        with mock.patch.object(cs, "apps_to_pause", return_value=[]), \
             mock.patch.object(cs, "cronjobs", return_value=[]), \
             mock.patch.object(cs, "replication_sources", return_value=[]), \
             mock.patch.object(cs, "workloads", return_value=[("onemind-core", "deployment", reverted)]), \
             mock.patch("builtins.print"):
            cs.do_standby(["onemind-core"], apply=True)
        self.assertIn(("scale", "deployment", "openfga", "-n", "onemind-core", "--replicas=0"), self.calls)
        self.assertFalse([c for c in self.calls if c[0] == "annotate"], "original count must not be overwritten")

    def test_re_enabled_app_is_paused_again_keeping_the_first_record(self):
        original = json.dumps({"automated": {"selfHeal": True}, "retry": {"limit": 3}})
        a = app("ecosystem-openfga", "onemind-core", annotations={cs.A_SYNC: original})
        with mock.patch.object(cs, "apps_to_pause", return_value=[a]), \
             mock.patch.object(cs, "cronjobs", return_value=[]), \
             mock.patch.object(cs, "replication_sources", return_value=[]), \
             mock.patch.object(cs, "workloads", return_value=[]), \
             mock.patch("builtins.print"):
            cs.do_standby(["onemind-core"], apply=True)
        self.assertEqual([c[0] for c in self.calls], ["patch"])

    def test_wake_restores_parents_last(self):
        sp = json.dumps({"automated": {"selfHeal": True}})
        parent = app("app-of-apps", "argocd", automated=False, annotations={cs.A_SYNC: sp})
        child = app("platform-n8n", "n8n", automated=False, annotations={cs.A_SYNC: sp})

        def fake_annotated(kind, key):
            return [parent, child] if kind.startswith("applications") else []
        with mock.patch.object(cs, "annotated", side_effect=fake_annotated), mock.patch("builtins.print"):
            cs.do_wake(apply=True)
        patched = [c[2] for c in self.calls if c[:2] == ("patch", "application")]
        self.assertEqual(patched, ["platform-n8n", "app-of-apps"])

    def test_keep_list_holds_hermes_and_every_fail_closed_webhook_namespace(self):
        for ns in ("hermes", "tailscale", "auth", "onemind-data", "argocd",
                   "cnpg-system", "kyverno", "cert-manager", "ingress-nginx"):
            self.assertIn(ns, cs.KEEP_NAMESPACES)


if __name__ == "__main__":
    unittest.main()
