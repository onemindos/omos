"""Tool gateway (platform/tool-gateway) contracts: the fragment renderer and the committed output."""
from pathlib import Path
import importlib.util
import unittest

import yaml

ROOT = Path(__file__).resolve().parents[2]
BASE = ROOT / "platform/tool-gateway"
spec = importlib.util.spec_from_file_location("render_gateway_config", ROOT / "scripts/gateway/render_gateway_config.py")
render = importlib.util.module_from_spec(spec)
spec.loader.exec_module(render)

GW = {
    "port": 8080,
    "clients": {"hermes": {"env": "GATEWAY_KEY_HERMES", "secret": "gateway-clients", "key": "hermes"}},
    "groups": {"ops": {"path": "/ops/mcp", "clients": ["hermes"]}, "dev": {"path": "/dev/mcp", "clients": ["hermes"]}},
}


def frag(name="grafana", **kw):
    base = {
        "name": name, "groups": ["ops"], "_file": f"{name}.yaml",
        "env": {f"TARGET_{name.upper()}_TOKEN": {"secret": f"gateway-target-{name}", "key": "token"}},
        "target": {"mcp": {"host": f"http://{name}:8000/mcp"},
                   "policies": {"backendAuth": {"key": "${TARGET_%s_TOKEN}" % name.upper()}}},
    }
    base.update(kw)
    return base


class Renderer(unittest.TestCase):
    def test_committed_output_is_current(self):
        self.assertEqual(render.OUT.read_text(), render.dump(render.render(*render.load())))

    def test_every_route_requires_a_client_key(self):
        values = render.render(*render.load())
        routes = values["config"]["binds"][0]["listeners"][0]["routes"]
        self.assertTrue(routes)
        for r in routes:
            # The agentgateway default lets keyless requests through: strict is mandatory (Phase 0, A7).
            self.assertEqual(r["policies"]["apiKey"]["mode"], "strict", r["name"])
            self.assertTrue(r["policies"]["apiKey"]["keys"], r["name"])

    def test_no_plaintext_credentials_in_routes(self):
        routes = render.render(*render.load())["config"]["binds"][0]["listeners"][0]["routes"]

        def keys(node):
            if isinstance(node, dict):
                for k, v in node.items():
                    if k == "key":
                        yield v
                    yield from keys(v)
            elif isinstance(node, list):
                for v in node:
                    yield from keys(v)

        found = list(keys(routes))
        self.assertTrue(found)
        for value in found:
            self.assertRegex(value, r"^\$\{[A-Z0-9_]+\}$")

    def test_groups_without_targets_are_skipped_and_sorted(self):
        out = render.render(GW, [frag("zeta"), frag("alpha")])
        routes = out["config"]["binds"][0]["listeners"][0]["routes"]
        self.assertEqual([r["name"] for r in routes], ["ops"])
        targets = routes[0]["backends"][0]["mcp"]["targets"]
        self.assertEqual([t["name"] for t in targets], ["alpha", "zeta"])

    def test_rejects_bad_fragments(self):
        bad = [
            frag(name="my-tool", **{"_file": "my-tool.yaml"}),  # hyphen
            frag(groups=["nope"]),
            frag(**{"_file": "other.yaml"}),
            frag(env={"OTHER_TOKEN": {"secret": "s", "key": "k"}}),  # wrong prefix
            frag(target={"mcp": {"host": "http://x/mcp"}, "policies": {"backendAuth": {"key": "${UNDECLARED}"}}}),
            frag(target={"mcp": {"host": "http://x/mcp"}, "policies": {"x": "price$5"}}),  # bare dollar
            frag(target={"name": "x", "mcp": {"host": "http://x/mcp"}}),
        ]
        for f in bad:
            with self.subTest(f=f["_file"]), self.assertRaises(render.RenderError):
                render.render(GW, [f])

    def test_escaped_dollar_is_allowed(self):
        out = render.render(GW, [frag(target={"mcp": {"host": "http://x/mcp"}, "policies": {"x": "a$$b"}})])
        self.assertTrue(out["config"]["binds"])


class Manifests(unittest.TestCase):
    def docs(self, rel):
        return [d for d in yaml.safe_load_all((BASE / rel).read_text()) if d]

    def test_only_hermes_may_call_the_gateway(self):
        gw = next(d for d in self.docs("manifests/network.yaml") if d["metadata"]["name"] == "agentgateway")
        peers = [(f["namespaceSelector"]["matchLabels"]["kubernetes.io/metadata.name"],
                  f.get("podSelector", {}).get("matchLabels", {}).get("app.kubernetes.io/name"))
                 for rule in gw["spec"]["ingress"] for f in rule["from"]]
        self.assertEqual(peers, [("hermes", "onemind-hermes"), ("monitoring", "prometheus")])

    def test_default_deny_both_directions(self):
        d = next(d for d in self.docs("manifests/network.yaml") if d["metadata"]["name"] == "default-deny")
        self.assertEqual(d["spec"]["podSelector"], {})
        self.assertEqual(set(d["spec"]["policyTypes"]), {"Ingress", "Egress"})

    def test_chart_pinned_and_values_have_no_lists_that_helm_would_replace(self):
        app = yaml.safe_load((ROOT / "platform/argocd/platform-tool-gateway.yaml").read_text())
        chart = app["spec"]["sources"][0]
        self.assertEqual((chart["chart"], chart["targetRevision"]), ("agentgateway-standalone", "v1.7.0"))
        values = yaml.safe_load((BASE / "values.yaml").read_text())
        self.assertNotIn("extraEnv", values)
        self.assertNotIn("binds", values.get("config", {}))
        self.assertIn("@sha256:", values["image"]["tag"])
        self.assertEqual(values["nodeSelector"], {"onemind.io/compute-role": "ecosystem"})


if __name__ == "__main__":
    unittest.main()
