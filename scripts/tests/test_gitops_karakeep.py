"""GitOps launch, private identity and recovery contracts for native Karakeep."""
from __future__ import annotations

import ipaddress
import json
import shutil
import subprocess
import sys
import types
import unittest
from pathlib import Path
from unittest.mock import patch

import yaml

ROOT = Path(__file__).resolve().parents[2]
APP = ROOT / "platform/karakeep"


class KarakeepDeployment(unittest.TestCase):
    @classmethod
    def setUpClass(cls):
        cls.kustomization = yaml.safe_load((APP / "kustomization.yaml").read_text())
        cls.docs = []
        for resource in cls.kustomization["resources"]:
            cls.docs.extend(yaml.safe_load_all((APP / resource).read_text()))
        cls.deployments = {d["metadata"]["name"]: d for d in cls.docs if d["kind"] == "Deployment"}
        cls.policies = {d["metadata"]["name"]: d["spec"] for d in cls.docs if d["kind"] == "NetworkPolicy"}

    def test_authorized_first_start_waits_for_identity_and_keeps_one_sqlite_writer(self):
        app = yaml.safe_load((ROOT / "platform/argocd/platform-karakeep.yaml").read_text())
        self.assertEqual(app["spec"]["syncPolicy"]["automated"], {"prune": True, "selfHeal": True})
        self.assertEqual(app["spec"]["source"]["path"], "platform/karakeep")
        self.assertEqual(set(self.deployments), {"web", "chrome", "meilisearch"})
        self.assertTrue(all(d["spec"]["replicas"] == 1 for d in self.deployments.values()))
        hook, = [d for d in self.docs if d["kind"] == "Job"]
        self.assertEqual(hook["metadata"]["namespace"], "auth")
        self.assertEqual(hook["metadata"]["annotations"]["argocd.argoproj.io/hook"], "Sync")
        self.assertLess(int(hook["metadata"]["annotations"]["argocd.argoproj.io/sync-wave"]),
                        int(self.deployments["web"]["metadata"]["annotations"]["argocd.argoproj.io/sync-wave"]))
        self.assertLessEqual(hook["spec"]["backoffLimit"], 3)
        self.assertEqual(hook["spec"]["activeDeadlineSeconds"], 600)
        self.assertFalse(any(d["kind"] in {"Cluster", "Database", "Role", "RoleBinding", "ClusterRole", "ClusterRoleBinding"}
                             for d in self.docs))

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
        # A parent namespace transformer must not move the auth-only hook/secret.
        self.assertNotIn("namespace", self.kustomization)
        generator, = self.kustomization["configMapGenerator"]
        self.assertEqual(generator["name"], "karakeep-configuration")
        self.assertEqual(generator["namespace"], "karakeep")
        configuration = dict(item.split("=", 1) for item in generator["literals"])
        self.assertEqual(configuration["NEXTAUTH_URL"], "https://karakeep.tail717ba0.ts.net")
        self.assertEqual(configuration["DISABLE_PASSWORD_AUTH"], "true")
        # First OIDC enrollment is possible only for the one bound Authentik owner.
        self.assertEqual(configuration["DISABLE_SIGNUPS"], "false")
        self.assertEqual(configuration["OAUTH_CLIENT_ID"], "onemind-karakeep")
        self.assertEqual(configuration["OAUTH_ALLOW_DANGEROUS_EMAIL_ACCOUNT_LINKING"], "false")
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
            if doc["kind"] in {"Job", "SealedSecret"} and doc["metadata"]["name"] in {
                "karakeep-authentik-oidc-v1", "karakeep-oidc-bootstrap"
            }:
                self.assertEqual(doc["metadata"]["namespace"], "auth")
            elif doc["kind"] != "Namespace":
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
        sealed = {d["metadata"]["name"]: d for d in self.docs if d["kind"] == "SealedSecret"}
        self.assertEqual(set(sealed["karakeep-secrets"]["spec"]["encryptedData"]), {"NEXTAUTH_SECRET", "MEILI_MASTER_KEY"})
        self.assertEqual(set(sealed["karakeep-oidc"]["spec"]["encryptedData"]), {"OAUTH_CLIENT_SECRET"})
        self.assertEqual(set(sealed["karakeep-oidc-bootstrap"]["spec"]["encryptedData"]), {"OAUTH_CLIENT_SECRET"})
        for secret in sealed.values():
            self.assertNotIn("sealedsecrets.bitnami.com/namespace-wide", secret["metadata"].get("annotations", {}))
            self.assertNotIn("sealedsecrets.bitnami.com/cluster-wide", secret["metadata"].get("annotations", {}))
        web = self.deployments["web"]["spec"]["template"]["spec"]["containers"][0]
        self.assertEqual({e["secretRef"]["name"] for e in web["envFrom"] if "secretRef" in e},
                         {"karakeep-secrets", "karakeep-oidc"})
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

    def test_backup_copies_entire_claims_with_distinct_encrypted_repositories_and_no_workload_rights(self):
        sources = [d for d in self.docs if d["kind"] == "ReplicationSource"]
        self.assertEqual({d["spec"]["sourcePVC"] for d in sources}, {"data-pvc", "meilisearch-pvc"})
        self.assertEqual(len({d["spec"]["restic"]["repository"] for d in sources}), 2)
        self.assertEqual(len({d["spec"]["trigger"]["schedule"] for d in sources}), 2)
        for source in sources:
            spec = source["spec"]
            self.assertFalse(spec["paused"])
            self.assertEqual(spec["restic"]["copyMethod"], "Snapshot")
            self.assertEqual(spec["restic"]["moverSecurityContext"]["runAsUser"], 1000)
            self.assertEqual(spec["restic"]["retain"]["daily"], 7)
            self.assertEqual(spec["restic"]["retain"]["weekly"], 4)
            self.assertEqual(spec["restic"]["retain"]["monthly"], 3)
            secret = next(d for d in self.docs if d["kind"] == "SealedSecret"
                          and d["metadata"]["name"] == spec["restic"]["repository"])
            self.assertEqual(set(secret["spec"]["encryptedData"]),
                             {"RESTIC_REPOSITORY", "RESTIC_PASSWORD", "AWS_ACCESS_KEY_ID", "AWS_SECRET_ACCESS_KEY", "AWS_DEFAULT_REGION"})
        self.assertIn("crash-consistent", (APP / "BACKUP.md").read_text())

    def run_identity_hook(self, *, uid=None, bindings=None):
        """Run the exact Argo hook against a bounded in-memory identity API."""
        hook, = [d for d in self.docs if d["kind"] == "Job"]
        container, = hook["spec"]["template"]["spec"]["containers"]
        self.assertFalse(hook["spec"]["template"]["spec"]["automountServiceAccountToken"])
        self.assertEqual({e["valueFrom"]["secretKeyRef"]["name"] for e in container["env"]},
                         {"authentik-secret", "karakeep-oidc-bootstrap"})
        owner_uid = "f73e10d872325ed63330a138782f12d3e8252e00a05466e2f0180bcc885a42c4"
        calls = []
        state = {"bindings": list(bindings or [])}

        class Client:
            def __init__(self, **kwargs):
                pass

            def __enter__(self):
                return self

            def __exit__(self, *_):
                return False

            def request(self, method, url, json=None):
                path = url.split("/api/v3/", 1)[1]
                calls.append((method, path, json))
                if method == "GET" and path.startswith("core/users/"):
                    data = {"results": [{"pk": 4, "username": "akadmin", "uid": uid or owner_uid, "is_active": True}]}
                elif method == "GET" and path.startswith("providers/oauth2/"):
                    data = {"results": [{"pk": 1, "client_id": "onemind-backstage", "signing_key": "existing-key",
                                         "authorization_flow": "existing-flow", "invalidation_flow": "existing-logout",
                                         "property_mappings": ["existing-mapping"], "sub_mode": "hashed_user_id",
                                         "issuer_mode": "per_provider", "access_token_validity": "minutes=10"}]}
                elif method == "GET" and path.startswith("core/applications/"):
                    return types.SimpleNamespace(status_code=404)
                elif method == "GET" and path.startswith("policies/bindings/"):
                    data = {"results": state["bindings"]}
                elif method == "POST" and path == "providers/oauth2/":
                    data = {"pk": 42, **json}
                elif method == "POST" and path == "core/applications/":
                    data = {"pk": "karakeep-app", **json}
                elif method == "POST" and path == "policies/bindings/":
                    state["bindings"].append(json)
                    data = json
                else:
                    raise AssertionError((method, path))
                return types.SimpleNamespace(status_code=200, is_success=True, content=b"present", json=lambda: data)

        with patch.dict(sys.modules, {"httpx": types.SimpleNamespace(Client=Client)}), \
             patch.dict("os.environ", {"AUTH_TOKEN": "unit-test-only", "OAUTH_CLIENT_SECRET": "unit-test-only"}):
            exec(compile(container["args"][0], "karakeep-auth-hook", "exec"), {})
        return calls

    def test_identity_bootstrap_grants_only_existing_owner_access(self):
        calls = self.run_identity_hook()
        mutations = [(method, path, body) for method, path, body in calls if method != "GET"]
        self.assertEqual({path for _, path, _ in mutations},
                         {"providers/oauth2/", "core/applications/", "policies/bindings/"})
        provider = next(body for _, path, body in mutations if path == "providers/oauth2/")
        self.assertEqual(provider["client_id"], "onemind-karakeep")
        self.assertEqual(provider["authorization_flow"], "existing-flow")
        self.assertEqual(provider["redirect_uris"], [{"matching_mode": "strict",
                         "url": "https://karakeep.tail717ba0.ts.net/api/auth/callback/custom", "redirect_uri_type": "authorization"}])
        app = next(body for _, path, body in mutations if path == "core/applications/")
        self.assertEqual(app["policy_engine_mode"], "all")
        binding = next(body for _, path, body in mutations if path == "policies/bindings/")
        self.assertEqual(binding, {"target": "karakeep-app", "user": 4, "enabled": True, "negate": False, "order": 0})

    def test_identity_bootstrap_refuses_replaced_owner_or_extra_access(self):
        with self.assertRaisesRegex(RuntimeError, "owner identity"):
            self.run_identity_hook(uid="replaced-identity")
        with self.assertRaisesRegex(RuntimeError, "Unexpected Karakeep application bindings"):
            self.run_identity_hook(bindings=[{"user": 99, "enabled": True}])


if __name__ == "__main__":
    unittest.main()
