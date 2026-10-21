"""Offline source contracts; optional pinned-chart render checks, never deployment."""

import argparse
from pathlib import Path
import subprocess
import unittest

import yaml

ROOT = Path(__file__).resolve().parents[2]
PILOT = ROOT / "platform/librechat/pilot"
APP = yaml.safe_load((ROOT / "platform/argocd/platform-librechat-pilot.yaml").read_text())
VALUES = yaml.safe_load((PILOT / "values.yaml").read_text())
CONFIG = yaml.safe_load(VALUES["librechat"]["configYamlContent"])
SUPPORT = [doc for doc in yaml.safe_load_all((PILOT / "network-policy.yaml").read_text()) if doc]
RENDERED = []
REQUIRED_KEYS = {
    "MONGO_URI", "CREDS_KEY", "CREDS_IV", "JWT_SECRET", "JWT_REFRESH_SECRET",
    "OPENID_CLIENT_ID", "OPENID_CLIENT_SECRET", "OPENID_SESSION_SECRET", "HERMES_API_KEY",
}
IMAGE = ("registry.librechat.ai/librechat-ai/librechat:"
         "v0.8.8@sha256:5b706bcee1ca154a1e2987303aaeca0eb27f97919cccee4bb604226c2f8abf61")


class SourceContracts(unittest.TestCase):
    def test_manual_gitops_no_cascade(self):
        self.assertNotIn("automated", APP["spec"]["syncPolicy"])
        self.assertNotIn("finalizers", APP["metadata"])
        self.assertIn("Prune=false", APP["spec"]["syncPolicy"]["syncOptions"])
        self.assertEqual(APP["spec"]["sources"][1]["path"], "platform/librechat/pilot")
        self.assertEqual(APP["spec"]["destination"]["namespace"], "librechat-pilot")

    def test_pinned_chart_image_and_no_bundled_dependencies(self):
        self.assertEqual(APP["spec"]["sources"][0]["targetRevision"], "2.0.16")
        self.assertEqual(VALUES["image"]["tag"], IMAGE.split("librechat:", 1)[1])
        for dep in ("mongodb", "redis", "meilisearch", "librechat-rag-api", "langfuseFanout"):
            self.assertIs(VALUES[dep]["enabled"], False)

    def test_each_secret_key_required_not_plaintext(self):
        # No bulk envFrom Secret that could override the security config.
        self.assertEqual(VALUES["global"]["librechat"]["existingSecretName"], "")
        refs = VALUES["global"]["librechat"]["env"]
        self.assertEqual({e["name"] for e in refs}, REQUIRED_KEYS)
        for entry in refs:
            self.assertNotIn("value", entry)
            self.assertEqual(entry["valueFrom"]["secretKeyRef"], {
                "name": "librechat-runtime", "key": entry["name"], "optional": False})
        self.assertFalse(REQUIRED_KEYS & VALUES["librechat"]["configEnv"].keys())

    def test_private_oidc_required_group_and_no_email_login(self):
        env = VALUES["librechat"]["configEnv"]
        self.assertEqual(env["ENDPOINTS"], "custom")
        self.assertEqual(env["OPENID_REQUIRED_ROLE"], "onemind-librechat-pilot")
        self.assertEqual(env["OPENID_REQUIRED_ROLE_TOKEN_KIND"], "id")
        self.assertEqual(env["OPENID_REQUIRED_ROLE_PARAMETER_PATH"], "groups")
        for key in ("ALLOW_EMAIL_LOGIN", "ALLOW_REGISTRATION", "ALLOW_SHARED_LINKS",
                    "ALLOW_SHARED_LINKS_PUBLIC", "OPENID_REUSE_TOKENS"):
            self.assertEqual(env[key], "false")
        self.assertTrue(env["OPENID_ISSUER"].startswith("https://onemind-auth.tail717ba0.ts.net/"))
        self.assertEqual(env["DOMAIN_SERVER"], "https://librechat.tail717ba0.ts.net")
        self.assertNotIn("OPENID_ADMIN_ROLE", env)

    def test_real_agent_not_direct_provider_or_hidden_calls(self):
        endpoint, = CONFIG["endpoints"]["custom"]
        self.assertEqual(endpoint["name"], "Legacy")
        self.assertEqual(endpoint["apiKey"], "${HERMES_API_KEY}")
        self.assertEqual(endpoint["baseURL"], "http://hermes-dashboard.hermes.svc.cluster.local:8642/v1")
        self.assertEqual(endpoint["models"], {"default": ["hermes-agent"], "fetch": False})
        self.assertIs(endpoint["titleConvo"], False)
        self.assertIs(endpoint["summarize"], False)
        self.assertEqual(CONFIG["endpoints"]["allowedAddresses"],
                         ["hermes-dashboard.hermes.svc.cluster.local:8642"])
        self.assertNotIn("mcpServers", CONFIG)

    def test_header_is_user_and_conversation_scoped(self):
        endpoint, = CONFIG["endpoints"]["custom"]
        self.assertEqual(endpoint["headers"], {"X-Hermes-Session-Key":
            "librechat:{{LIBRECHAT_USER_ID}}:{{LIBRECHAT_BODY_CONVERSATIONID}}"})
        # This tests configuration only, not runtime authorization/isolation.

    def test_unproven_files_disabled_and_no_fictitious_billing(self):
        for endpoint in ("default", "Legacy"):
            self.assertIs(CONFIG["fileConfig"]["endpoints"][endpoint]["disabled"], True)
        self.assertIs(CONFIG["balance"]["enabled"], False)
        self.assertNotIn("tokenConfig", CONFIG["endpoints"]["custom"][0])

    def test_single_replica_persistence_recreate_and_probes(self):
        self.assertEqual(VALUES["replicaCount"], 1)
        self.assertEqual(VALUES["updateStrategy"]["type"], "Recreate")
        self.assertIs(VALUES["librechat"]["imageVolume"]["enabled"], True)
        self.assertEqual(VALUES["readinessProbe"]["httpGet"]["path"], "/readyz")
        self.assertEqual(VALUES["livenessProbe"]["httpGet"]["path"], "/livez")
        self.assertIs(VALUES["serviceAccount"]["automount"], False)
        self.assertIs(VALUES["securityContext"]["allowPrivilegeEscalation"], False)

    def test_default_deny_and_only_intended_network_paths(self):
        deny, allow = SUPPORT
        self.assertEqual(deny["spec"], {"podSelector": {}, "policyTypes": ["Ingress", "Egress"]})
        spec = allow["spec"]
        self.assertEqual(spec["podSelector"]["matchLabels"]["app.kubernetes.io/name"], "librechat-pilot")
        self.assertEqual(len(spec["egress"]), 4)
        self.assertEqual({port["port"] for e in spec["egress"] for port in e["ports"]}, {53, 443, 8642, 27017})
        self.assertNotIn("ipBlock", yaml.safe_dump(spec))
        for e in spec["egress"]:
            self.assertTrue(e["to"])
            self.assertTrue(all(peer for peer in e["to"]))
        mesh = spec["egress"][2]["to"][0]
        self.assertEqual(mesh["namespaceSelector"]["matchLabels"]["kubernetes.io/metadata.name"], "tailscale")
        self.assertEqual(mesh["podSelector"]["matchLabels"]["tailscale.com/parent-resource"], "egress")
        self.assertEqual(spec["ingress"][0]["ports"], [{"protocol": "TCP", "port": 3080}])

    def test_no_new_rbac_or_plaintext_secret_in_support(self):
        for f in PILOT.glob("*.yaml"):
            if f.name == "values.yaml":
                continue
            for doc in yaml.safe_load_all(f.read_text()):
                self.assertNotIn(doc["kind"], ("Secret", "Role", "ClusterRole", "RoleBinding", "ClusterRoleBinding"))
        ingress = yaml.safe_load((PILOT / "ingress.yaml").read_text())
        self.assertEqual(ingress["spec"]["ingressClassName"], "tailscale")


@unittest.skipUnless(RENDERED, "provide --rendered or --render-with-helm for real chart checks")
class RenderContracts(unittest.TestCase):
    def test_chart_workload_matches_policy_and_requires_each_key(self):
        workloads = [d for d in RENDERED if d["kind"] in ("Deployment", "StatefulSet")]
        self.assertEqual(len(workloads), 1)
        deploy = workloads[0]
        pod = deploy["spec"]["template"]
        self.assertEqual(pod["metadata"]["labels"]["app.kubernetes.io/name"], "librechat-pilot")
        self.assertEqual(deploy["spec"]["strategy"]["type"], "Recreate")
        self.assertEqual(deploy["spec"]["replicas"], 1)
        container, = pod["spec"]["containers"]
        self.assertEqual(container["image"], IMAGE)
        self.assertEqual({e["name"] for e in container["env"]}, REQUIRED_KEYS)
        self.assertFalse(any("secretRef" in e for e in container["envFrom"]))
        self.assertEqual(container["readinessProbe"]["httpGet"]["path"], "/readyz")
        self.assertIn("requests", container["resources"])
        sa, = [d for d in RENDERED if d["kind"] == "ServiceAccount"]
        self.assertIs(sa["automountServiceAccountToken"], False)
        self.assertEqual(pod["spec"]["serviceAccountName"], sa["metadata"]["name"])
        self.assertNotIn("ClusterRoleBinding", {d["kind"] for d in RENDERED})
        self.assertNotIn("Secret", {d["kind"] for d in RENDERED})

    def test_config_and_image_pvc_actually_render(self):
        config, = [d for d in RENDERED if d["kind"] == "ConfigMap" and "librechat.yaml" in d["data"]]
        self.assertEqual(yaml.safe_load(config["data"]["librechat.yaml"]), CONFIG)
        pvc, = [d for d in RENDERED if d["kind"] == "PersistentVolumeClaim"]
        self.assertEqual(pvc["spec"]["accessModes"], ["ReadWriteOnce"])


if __name__ == "__main__":
    parser = argparse.ArgumentParser(description=__doc__)
    group = parser.add_mutually_exclusive_group()
    group.add_argument("--rendered", type=Path)
    group.add_argument("--render-with-helm", action="store_true")
    args = parser.parse_args()
    if args.render_with_helm:
        rendered_text = subprocess.check_output([
            "helm", "template", "librechat-pilot", "oci://ghcr.io/librechat-ai/librechat-chart/librechat",
            "--version", "2.0.16", "--skip-tests", "--namespace", "librechat-pilot", "-f", str(PILOT / "values.yaml"),
        ], text=True, timeout=120)
    elif args.rendered:
        rendered_text = args.rendered.read_text()
    else:
        rendered_text = ""
    RENDERED.extend(d for d in yaml.safe_load_all(rendered_text) if isinstance(d, dict) and d.get("kind"))
    # skipUnless was evaluated before command-line rendering populated the list.
    if RENDERED:
        RenderContracts.__unittest_skip__ = False
    unittest.main(argv=[__file__], verbosity=2)
