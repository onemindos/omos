"""Safety gates for the staged native Karakeep deployment, without image execution."""
from __future__ import annotations

import ipaddress
import json
import shutil
import subprocess
import unittest
from pathlib import Path

import yaml

ROOT = Path(__file__).resolve().parents[2]
APP = ROOT / "platform/karakeep"


class KarakeepPreparation(unittest.TestCase):
    @classmethod
    def setUpClass(cls):
        cls.kustomization = yaml.safe_load((APP / "kustomization.yaml").read_text())
        cls.docs = []
        for resource in cls.kustomization["resources"]:
            cls.docs.extend(yaml.safe_load_all((APP / resource).read_text()))
        cls.deployments = {d["metadata"]["name"]: d for d in cls.docs if d["kind"] == "Deployment"}
        cls.policies = {d["metadata"]["name"]: d["spec"] for d in cls.docs if d["kind"] == "NetworkPolicy"}

    def test_application_is_manual_and_no_process_or_migration_can_start(self):
        app = yaml.safe_load((ROOT / "platform/argocd/platform-karakeep.yaml").read_text())
        self.assertNotIn("automated", app["spec"]["syncPolicy"])
        self.assertEqual(app["spec"]["source"]["path"], "platform/karakeep")
        self.assertEqual(set(self.deployments), {"web", "chrome", "meilisearch"})
        self.assertTrue(all(d["spec"]["replicas"] == 0 for d in self.deployments.values()))
        self.assertFalse(any(d["kind"] in {"Job", "CronJob", "Cluster", "Database"} for d in self.docs))

    def test_images_match_locked_native_source_and_have_immutable_digests(self):
        lock = json.loads((APP / "upstream.lock.json").read_text())
        self.assertRegex(lock["upstream"]["commit"], r"^[0-9a-f]{40}$")
        self.assertEqual(lock["upstream"]["nativePath"], "kubernetes")
        for name, deployment in self.deployments.items():
            image = deployment["spec"]["template"]["spec"]["containers"][0]["image"]
            self.assertEqual(image, lock["images"][name])
            self.assertRegex(image, r":[^:@]+@sha256:[0-9a-f]{64}$")
            self.assertNotRegex(image, r":(latest|release)@")

    def test_only_private_tailnet_web_listener_is_exposed(self):
        services = [d for d in self.docs if d["kind"] == "Service"]
        self.assertTrue(all(d["spec"].get("type", "ClusterIP") == "ClusterIP" for d in services))
        ingress = next(d for d in self.docs if d["kind"] == "Ingress")
        self.assertEqual(ingress["spec"]["ingressClassName"], "tailscale")
        self.assertEqual(ingress["metadata"]["annotations"], {"tailscale.com/hostname": "karakeep"})
        self.assertEqual(ingress["spec"]["defaultBackend"]["service"]["name"], "web")
        self.assertEqual(self.kustomization["namespace"], "karakeep")
        self.assertEqual(set(self.kustomization), {"apiVersion", "kind", "namespace", "resources", "configMapGenerator"})
        generator, = self.kustomization["configMapGenerator"]
        self.assertEqual(generator["name"], "karakeep-configuration")
        configuration = dict(item.split("=", 1) for item in generator["literals"])
        self.assertEqual(configuration["NEXTAUTH_URL"], "https://karakeep.tail717ba0.ts.net")
        self.assertEqual(configuration["DISABLE_SIGNUPS"], "true")
        self.assertEqual(configuration["INFERENCE_ENABLE_AUTO_TAGGING"], "false")
        self.assertNotIn("OPENAI_API_KEY", configuration)

    @unittest.skipUnless(shutil.which("kubectl"), "Kustomize render requires kubectl; source safety checks always run")
    def test_native_kustomize_render_preserves_identity_and_images(self):
        rendered = list(yaml.safe_load_all(subprocess.check_output(
            ["kubectl", "kustomize", str(APP)], text=True)))
        expected = {(d["kind"], d["metadata"]["name"]) for d in self.docs}
        actual = {(d["kind"], d["metadata"]["name"]) for d in rendered}
        self.assertTrue(expected <= actual)
        self.assertEqual(len(actual - expected), 1)
        for doc in rendered:
            if doc["kind"] != "Namespace":
                self.assertEqual(doc["metadata"]["namespace"], "karakeep")
            if doc["kind"] == "Deployment":
                source = self.deployments[doc["metadata"]["name"]]
                self.assertEqual([c["image"] for c in doc["spec"]["template"]["spec"]["containers"]],
                                 [c["image"] for c in source["spec"]["template"]["spec"]["containers"]])

    def test_crawler_public_fetch_cannot_reach_private_cluster_or_metadata_addresses(self):
        deny = self.policies["default-deny"]
        self.assertEqual(deny["podSelector"], {})
        self.assertEqual(set(deny["policyTypes"]), {"Ingress", "Egress"})
        self.assertNotIn("ingress", deny)
        self.assertNotIn("egress", deny)
        fetch = self.policies["public-web-fetch"]
        self.assertEqual(fetch["podSelector"]["matchExpressions"][0]["values"], ["web", "chrome"])
        entry = fetch["egress"][0]
        self.assertEqual({p["port"] for p in entry["ports"]}, {80, 443})
        block = entry["to"][0]["ipBlock"]
        excluded = [ipaddress.ip_network(c) for c in block["except"]]
        for address in ["10.109.3.74", "172.16.0.1", "192.168.1.1", "100.100.100.100", "169.254.169.254", "127.0.0.1"]:
            self.assertTrue(any(ipaddress.ip_address(address) in net for net in excluded), address)
        self.assertFalse(any(ipaddress.ip_address("1.1.1.1") in net for net in excluded))
        origin = self.policies["web"]["ingress"][0]["from"][0]
        self.assertEqual(origin["namespaceSelector"]["matchLabels"]["kubernetes.io/metadata.name"], "tailscale")
        self.assertEqual(origin["podSelector"]["matchLabels"]["tailscale.com/parent-resource-ns"], "karakeep")
        for name in ["chrome", "meilisearch"]:
            peers = self.policies[name]["ingress"][0]["from"]
            self.assertEqual(peers, [{"podSelector": {"matchLabels": {"app.kubernetes.io/name": "karakeep", "app.kubernetes.io/component": "web"}}}])

    def test_persistent_state_is_retained_and_has_only_one_writer(self):
        pvcs = [d for d in self.docs if d["kind"] == "PersistentVolumeClaim"]
        self.assertEqual({d["metadata"]["name"] for d in pvcs}, {"data-pvc", "meilisearch-pvc"})
        for pvc in pvcs:
            protections = pvc["metadata"]["annotations"]["argocd.argoproj.io/sync-options"]
            for option in ["Prune=false", "Delete=false", "Replace=false", "Force=false"]:
                self.assertIn(option, protections)
            self.assertEqual(pvc["spec"]["accessModes"], ["ReadWriteOnce"])
        for name in ["web", "meilisearch"]:
            self.assertEqual(self.deployments[name]["spec"]["strategy"]["type"], "Recreate")
        self.assertIn("No Karakeep backup is claimed", (APP / "BACKUP.md").read_text())

    def test_runtime_has_budgets_probes_and_no_api_tokens_or_plaintext_secrets(self):
        self.assertFalse(any(d["kind"] == "Secret" for d in self.docs))
        sealed = next(d for d in self.docs if d["kind"] == "SealedSecret")
        self.assertEqual(set(sealed["spec"]["encryptedData"]), {"NEXTAUTH_SECRET", "MEILI_MASTER_KEY"})
        self.assertEqual(sealed["metadata"]["namespace"], "karakeep")
        self.assertNotIn("sealedsecrets.bitnami.com/namespace-wide", sealed["metadata"].get("annotations", {}))
        for deployment in self.deployments.values():
            pod = deployment["spec"]["template"]["spec"]
            self.assertFalse(pod["automountServiceAccountToken"])
            self.assertTrue(pod["securityContext"]["runAsNonRoot"])
            for container in pod["containers"]:
                self.assertIn("memory", container["resources"]["requests"])
                self.assertIn("cpu", container["resources"]["requests"])
                self.assertIn("memory", container["resources"]["limits"])
                self.assertIn("cpu", container["resources"]["limits"])
                self.assertFalse(container["securityContext"]["allowPrivilegeEscalation"])
                for name in ["startupProbe", "readinessProbe", "livenessProbe"]:
                    self.assertIn(name, container)


if __name__ == "__main__":
    unittest.main()
