"""Tool gateway Kubernetes backend (platform/tool-gateway/manifests/backends/kubernetes): the RBAC ceiling and the pins."""
from pathlib import Path
import re
import unittest

import yaml

ROOT = Path(__file__).resolve().parents[2]
BACKEND = ROOT / "platform/tool-gateway/manifests/backends/kubernetes"
CLASSIFICATION = ROOT / "platform/tool-gateway/policy/classification/kubernetes.yaml"


def docs(name):
    return [d for d in yaml.safe_load_all((BACKEND / name).read_text()) if d]


def rules():
    out = []
    for d in docs("rbac.yaml"):
        if d["kind"] in ("ClusterRole", "Role"):
            out += [(d["metadata"]["name"], r) for r in d["rules"]]
    return out


class Rbac(unittest.TestCase):
    def test_no_secret_access_and_no_wildcard_core_resources(self):
        for name, r in rules():
            groups = r.get("apiGroups", [])
            res = r.get("resources", [])
            self.assertNotIn("secrets", res, name)
            if "" in groups or "*" in groups:
                self.assertNotIn("*", res, f"{name}: wildcard resources in the core group would include Secrets")
            self.assertNotIn("*", r["verbs"], name)
            self.assertNotIn("*", groups, name)

    def test_only_pods_can_be_deleted(self):
        for name, r in rules():
            if {"delete", "deletecollection"} & set(r["verbs"]):
                self.assertEqual((r["apiGroups"], r["resources"]), ([""], ["pods"]), name)

    def test_no_exec_proxy_rbac_or_create(self):
        for name, r in rules():
            self.assertNotIn("create", r["verbs"], name)
            for res in r.get("resources", []):
                self.assertFalse(res.startswith(("pods/exec", "pods/attach", "pods/portforward", "nodes/proxy", "pods/eviction")), name)
            self.assertNotIn("rbac.authorization.k8s.io", r.get("apiGroups", []), name)

    def test_serviceaccount_bindings_stay_on_k8s_operator(self):
        for d in docs("rbac.yaml"):
            if d["kind"] in ("ClusterRoleBinding", "RoleBinding"):
                self.assertEqual([(s["kind"], s["name"], s["namespace"]) for s in d["subjects"]],
                                 [("ServiceAccount", "k8s-operator", "tool-gateway")])


class Workload(unittest.TestCase):
    def setUp(self):
        self.dep = next(d for d in docs("deployment.yaml") if d["kind"] == "Deployment")
        self.pod = self.dep["spec"]["template"]["spec"]
        self.c = self.pod["containers"][0]

    def test_image_pinned_by_digest(self):
        self.assertRegex(self.c["image"], r"@sha256:[0-9a-f]{64}$")

    def test_resources_and_hardening(self):
        self.assertEqual(self.pod["nodeSelector"], {"onemind.io/compute-role": "ecosystem"})
        for kind in ("requests", "limits"):
            self.assertTrue({"cpu", "memory"} <= set(self.c["resources"][kind]))
        sc = self.c["securityContext"]
        self.assertTrue(sc["readOnlyRootFilesystem"])
        self.assertFalse(sc["allowPrivilegeEscalation"])
        self.assertTrue(self.pod["securityContext"]["runAsNonRoot"])

    def test_config_denies_secrets_and_disables_exec(self):
        cfg = next(d for d in docs("config.yaml"))["data"]["config.toml"]
        self.assertRegex(cfg, r'kind = "Secret"')
        for tool in ("pods_exec", "pods_run"):
            self.assertIn(f'"{tool}"', cfg)

    def test_selftest_configmap_matches_rbac(self):
        cm = docs("selftest.yaml")[0]
        self.assertEqual((cm["kind"], cm["metadata"]["name"], cm["metadata"]["namespace"]), ("ConfigMap", "gateway-selftest", "tool-gateway"))
        role = next(d for d in docs("rbac.yaml") if d["kind"] == "Role")
        self.assertEqual(role["rules"][0]["resourceNames"], ["gateway-selftest"])


class Classification(unittest.TestCase):
    def test_destructive_tools_are_not_exposed(self):
        c = yaml.safe_load(CLASSIFICATION.read_text())
        exposed = set(c["read"]) | set(c["write"])
        for tool in ("pods_delete", "resources_delete", "resources_scale"):
            self.assertIn(tool, c["destructive"])
            self.assertNotIn(tool, exposed)
        self.assertFalse(c.get("guards"), "argument predicates are inert in agentgateway v1.7.0")
        self.assertFalse(any(re.search(r"delete|exec|run", t) for t in exposed))


if __name__ == "__main__":
    unittest.main()
