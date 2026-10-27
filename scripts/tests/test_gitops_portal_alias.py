"""Contracts for the portal tailnet alias: it redirects to the single Backstage origin."""
from __future__ import annotations

import re
import unittest
from pathlib import Path
from urllib.parse import urlsplit

import yaml

ROOT = Path(__file__).resolve().parents[2]
APP = ROOT / "platform/backstage"


class PortalAlias(unittest.TestCase):
    @classmethod
    def setUpClass(cls):
        cls.kustomization = yaml.safe_load((APP / "kustomization.yaml").read_text())
        cls.docs = {d["kind"]: d for d in yaml.safe_load_all((APP / "portal-redirect.yaml").read_text())}
        backstage = next(yaml.safe_load_all((APP / "deployment.yaml").read_text()))
        env = backstage["spec"]["template"]["spec"]["containers"][0]["env"]
        cls.portal_url = next(e["value"] for e in env if e["name"] == "PORTAL_URL")
        cls.conf = (APP / "portal-redirect.conf").read_text()

    def test_alias_and_its_generated_config_are_registered(self):
        self.assertIn("portal-redirect.yaml", self.kustomization["resources"])
        generator = next(g for g in self.kustomization["configMapGenerator"] if g["name"] == "portal-redirect")
        self.assertEqual((generator["namespace"], generator["files"]), ("backstage", ["nginx.conf=portal-redirect.conf"]))
        volumes = self.docs["Deployment"]["spec"]["template"]["spec"]["volumes"]
        self.assertIn({"name": "config", "configMap": {"name": "portal-redirect"}}, volumes)

    def test_alias_temporarily_redirects_to_portal_url_and_never_serves_backstage(self):
        # Backstage derives baseUrl, CORS, CSP 'self' and the OIDC callback from PORTAL_URL alone.
        self.assertEqual(re.findall(r"return (30\d) (\S+);", self.conf), [("302", self.portal_url + "$request_uri")])
        self.assertNotIn("proxy_pass", self.conf)
        ingress = self.docs["Ingress"]
        alias = ingress["metadata"]["annotations"]["tailscale.com/hostname"]
        self.assertEqual((ingress["spec"]["ingressClassName"], ingress["spec"]["tls"]), ("tailscale", [{"hosts": [alias]}]))
        self.assertNotEqual(urlsplit(self.portal_url).hostname.split(".")[0], alias)  # would redirect to itself
        backend = ingress["spec"]["defaultBackend"]["service"]["name"]
        self.assertEqual(backend, self.docs["Service"]["metadata"]["name"])
        self.assertNotEqual(backend, "backstage")

    def test_only_the_alias_proxy_reaches_nginx_and_nginx_has_no_egress(self):
        ingress, policy = self.docs["Ingress"]["metadata"], self.docs["NetworkPolicy"]["spec"]
        deployment = self.docs["Deployment"]["spec"]
        self.assertEqual(policy["podSelector"]["matchLabels"], deployment["selector"]["matchLabels"])
        self.assertEqual(set(policy["policyTypes"]), {"Ingress", "Egress"})
        self.assertNotIn("egress", policy)
        rule, = policy["ingress"]
        self.assertEqual(rule["from"], [{
            "namespaceSelector": {"matchLabels": {"kubernetes.io/metadata.name": "tailscale"}},
            "podSelector": {"matchLabels": {
                "tailscale.com/parent-resource": ingress["name"],
                "tailscale.com/parent-resource-ns": ingress["namespace"],
                "tailscale.com/parent-resource-type": "ingress",
            }},
        }])
        container, = deployment["template"]["spec"]["containers"]
        self.assertEqual(rule["ports"], [{"protocol": "TCP", "port": container["ports"][0]["containerPort"]}])

    def test_redirect_pod_stays_out_of_backstage_selectors(self):
        # The backstage Service, its PodDisruptionBudget and postgres-ingress select app=backstage.
        labels = self.docs["Deployment"]["spec"]["template"]["metadata"]["labels"]
        self.assertNotEqual(labels.get("app"), "backstage")
        self.assertFalse(self.docs["Deployment"]["spec"]["template"]["spec"]["automountServiceAccountToken"])


if __name__ == "__main__":
    unittest.main()
