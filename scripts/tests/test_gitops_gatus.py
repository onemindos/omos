"""Gatus (flagship F7): manual sync, pinned image, small resources, hardened pod, config as code."""
from pathlib import Path
import re
import shutil
import subprocess
import unittest

import yaml

ROOT = Path(__file__).resolve().parents[2]
APP = ROOT / "platform/gatus"


def render():
    if not shutil.which("kubectl"):
        raise unittest.SkipTest("Kustomize render needs kubectl")
    return list(yaml.safe_load_all(subprocess.check_output(["kubectl", "kustomize", str(APP)], text=True)))


class GatusApp(unittest.TestCase):
    def test_application_is_manual_sync_and_in_the_gatus_namespace(self):
        app = yaml.safe_load((ROOT / "platform/argocd/platform-gatus.yaml").read_text())
        self.assertNotIn("automated", app["spec"].get("syncPolicy", {}))
        self.assertEqual(app["spec"]["source"]["path"], "platform/gatus")
        self.assertEqual(app["spec"]["destination"]["namespace"], "gatus")
        self.assertEqual(app["metadata"]["annotations"]["argocd.argoproj.io/sync-wave"], "5")

    def test_namespace_uses_the_peer_labels(self):
        ns = yaml.safe_load((APP / "namespace.yaml").read_text())
        self.assertEqual(ns["metadata"]["labels"]["onemind.io/layer"], "guard")
        for mode in ("enforce", "audit", "warn"):
            self.assertEqual(ns["metadata"]["labels"][f"pod-security.kubernetes.io/{mode}"], "restricted")

    def test_workload_is_pinned_small_hardened_and_stateless(self):
        docs = render()
        kinds = [d["kind"] for d in docs]
        self.assertNotIn("PersistentVolumeClaim", kinds)
        self.assertNotIn("StatefulSet", kinds)
        deploy = next(d for d in docs if d["kind"] == "Deployment")
        pod = deploy["spec"]["template"]["spec"]
        container = pod["containers"][0]
        self.assertRegex(container["image"], r"^twinproduction/gatus:v\d+\.\d+\.\d+@sha256:[0-9a-f]{64}$")
        self.assertNotIn(":latest", container["image"])
        res = container["resources"]
        self.assertEqual(res["requests"], {"cpu": "10m", "memory": "48Mi"})
        self.assertEqual(res["limits"], {"cpu": "200m", "memory": "128Mi"})
        self.assertTrue(pod["securityContext"]["runAsNonRoot"])
        self.assertEqual(pod["securityContext"]["seccompProfile"]["type"], "RuntimeDefault")
        self.assertFalse(pod["automountServiceAccountToken"])
        sc = container["securityContext"]
        self.assertFalse(sc["allowPrivilegeEscalation"])
        self.assertTrue(sc["readOnlyRootFilesystem"])
        self.assertEqual(sc["capabilities"]["drop"], ["ALL"])
        pools = pod["affinity"]["nodeAffinity"]["preferredDuringSchedulingIgnoredDuringExecution"][0]["preference"]["matchExpressions"][0]
        self.assertEqual(pools["values"], ["onemind-ecosystem-cpu"])

    def test_config_is_in_memory_with_metrics_and_merged_endpoint_files(self):
        config = yaml.safe_load((APP / "config.yaml").read_text())
        self.assertTrue(config["metrics"])
        self.assertEqual(config["storage"]["type"], "memory")
        docs = render()
        cm = next(d for d in docs if d["kind"] == "ConfigMap")
        self.assertEqual(sorted(cm["data"]), ["config.yaml", "endpoints.yaml"])
        # The pod mounts the whole directory, so Gatus merges every file.
        env = {e["name"]: e["value"] for e in next(d for d in docs if d["kind"] == "Deployment")["spec"]["template"]["spec"]["containers"][0]["env"]}
        self.assertEqual(env["GATUS_CONFIG_PATH"], "/config")

    def test_generated_endpoints_are_unique_and_no_deliberate_failure_remains(self):
        generated = yaml.safe_load((APP / "endpoints.yaml").read_text())["endpoints"]
        self.assertGreaterEqual(len(generated), 10)
        self.assertEqual(len({e["name"] for e in generated}), len(generated))
        self.assertFalse((APP / "test-endpoints.yaml").exists())

    def test_ui_is_tailnet_only_and_metrics_are_discovered(self):
        docs = render()
        ingress = next(d for d in docs if d["kind"] == "Ingress")
        self.assertEqual(ingress["spec"]["ingressClassName"], "tailscale")
        self.assertEqual(ingress["metadata"]["annotations"]["tailscale.com/hostname"], "status")
        monitor = next(d for d in docs if d["kind"] == "ServiceMonitor")
        self.assertEqual(monitor["metadata"]["labels"]["release"], "monitoring")
        self.assertEqual(monitor["spec"]["endpoints"][0]["path"], "/metrics")
        service = next(d for d in docs if d["kind"] == "Service")
        self.assertEqual(monitor["spec"]["selector"]["matchLabels"].items() <= service["metadata"]["labels"].items(), True)

    def test_alert_rules_cover_down_and_certificate_expiry_without_telegram(self):
        text = (APP / "prometheusrule.yaml").read_text()
        rule = yaml.safe_load(text)
        alerts = {r["alert"]: r for g in rule["spec"]["groups"] for r in g["rules"]}
        self.assertEqual(alerts["GatusEndpointDown"]["for"], "5m")
        self.assertIn("gatus_results_endpoint_success", alerts["GatusEndpointDown"]["expr"])
        self.assertIn("1209600", alerts["GatusCertificateExpiring"]["expr"])
        self.assertNotIn("receiver:", text)  # routing stays with the default Alertmanager route

    def test_ingress_policy_admits_only_the_three_readers(self):
        policy = yaml.safe_load((APP / "network.yaml").read_text())
        self.assertEqual(policy["spec"]["policyTypes"], ["Ingress"])
        namespaces = sorted(
            peer["namespaceSelector"]["matchLabels"]["kubernetes.io/metadata.name"] for peer in policy["spec"]["ingress"][0]["from"]
        )
        self.assertEqual(namespaces, ["backstage", "monitoring", "tailscale"])

    def test_no_secret_or_credential_is_committed_here(self):
        for path in APP.iterdir():
            if path.is_file():
                text = path.read_text()
                self.assertNotRegex(text, r"(?i)\b(kind: Secret|password:|api[_-]?key:|token:)")
                self.assertIsNone(re.search(r"https?://[^/\s]+:[^/\s]+@", text))


if __name__ == "__main__":
    unittest.main()
