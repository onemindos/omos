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
cl = render.classification
spec2 = importlib.util.spec_from_file_location("classify_tools", ROOT / "scripts/gateway/classify_tools.py")
classify = importlib.util.module_from_spec(spec2)
spec2.loader.exec_module(classify)

GW = {
    "port": 8080,
    "anchor": "grafana",
    "clients": {"hermes": {"env": "GATEWAY_KEY_HERMES", "secret": "gateway-clients", "key": "hermes"}},
    "groups": {
        "ops": {"path": "/ops/mcp", "clients": ["hermes"]},
        "dev": {"path": "/dev/mcp", "clients": ["hermes"]},
        "identity": {"kind": "deferred", "path": "/identity/mcp", "clients": ["hermes"]},
        "execute": {"kind": "execute", "path": "/execute/mcp", "clients": ["hermes"]},
    },
}


def klass(target="grafana", **kw):
    base = {"target": target, "read": [], "write": [], "destructive": [], "physical": [], "guards": [], "proposed": []}
    base.update(kw)
    return base


def routes_of(out):
    return {r["name"]: r for r in out["config"]["binds"][0]["listeners"][0]["routes"]}


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

    def test_all_groups_are_rendered_and_targets_sorted(self):
        out = render.render({**GW, "anchor": "alpha"}, [frag("zeta"), frag("alpha")], {"zeta": klass("zeta"), "alpha": klass("alpha")})
        routes = out["config"]["binds"][0]["listeners"][0]["routes"]
        self.assertEqual([r["name"] for r in routes], ["ops", "dev", "identity", "execute"])
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
                render.render(GW, [f], {f.get("name", "x"): klass(f.get("name", "x"))})

    def test_escaped_dollar_is_allowed(self):
        out = render.render(GW, [frag(target={"mcp": {"host": "http://x/mcp"}, "policies": {"x": "a$$b"}})],
                            {"grafana": klass()})
        self.assertTrue(out["config"]["binds"])

    def test_render_is_idempotent(self):
        inputs = render.load()
        first = render.dump(render.render(*inputs))
        self.assertEqual(first, render.dump(render.render(*inputs)))
        self.assertEqual(first, render.dump(render.render(*render.load())))

    def test_every_group_has_a_strict_route_and_authorization(self):
        gw, targets, classes = render.load()
        routes = routes_of(render.render(gw, targets, classes))
        self.assertEqual(set(routes), set(gw["groups"]))
        for want in ("ops", "dev", "cloud", "data", "tak", "home", "work", "identity", "execute"):
            self.assertIn(want, routes)
        for name, r in routes.items():
            self.assertEqual(r["policies"]["apiKey"]["mode"], "strict", name)
            # A route with no mcpAuthorization lists every tool (Phase 0, A1).
            self.assertTrue(r["policies"]["mcpAuthorization"]["rules"], name)

    def test_execute_and_deferred_are_denied(self):
        for gw, targets, classes in ((render.load()),
                                     (GW, [frag()], {"grafana": klass(read=["a"], destructive=["delete_x"])})):
            routes = routes_of(render.render(gw, targets, classes))
            for name in ("execute", "identity"):
                self.assertEqual(routes[name]["policies"]["mcpAuthorization"], {"rules": [{"require": "false"}]}, name)

    def test_unlisted_destructive_and_physical_tools_are_not_listed(self):
        c = klass(read=["get_*", "list"], write=["create_note"], destructive=["delete_note"], physical=["unlock_door"])
        out = render.render(GW, [frag()], {"grafana": c})
        rules = routes_of(out)["ops"]["policies"]["mcpAuthorization"]["rules"]
        self.assertEqual(len(rules), 1)
        expr = rules[0]["allow"]
        for listed in ("list", "create_note"):
            self.assertIn('"%s"' % listed, expr)
        self.assertIn("get_.*", expr)
        for hidden in ("delete_note", "unlock_door", "drop_table"):
            self.assertNotIn(hidden, expr)
        self.assertIn('mcp.tool.target == "grafana"', expr)
        self.assertNotIn("$ ", expr.replace("$$", ""))  # a literal dollar must be escaped

    def test_target_with_only_gated_tools_exposes_nothing_on_its_group(self):
        out = render.render(GW, [frag()], {"grafana": klass(destructive=["delete_x"])})
        self.assertEqual(routes_of(out)["ops"]["policies"]["mcpAuthorization"], {"rules": [{"require": "false"}]})

    def test_gated_targets_ride_on_execute_and_empty_groups_use_the_anchor(self):
        out = render.render(GW, [frag()], {"grafana": klass(read=["a"], physical=["unlock"])})
        routes = routes_of(out)
        names = lambda r: [t["name"] for t in r["backends"][0]["mcp"]["targets"]]  # noqa: E731
        self.assertEqual(names(routes["execute"]), ["grafana"])
        self.assertEqual(names(routes["dev"]), ["grafana"])  # anchor, hidden by require: "false"
        self.assertEqual(routes["dev"]["policies"]["mcpAuthorization"], {"rules": [{"require": "false"}]})

    def test_missing_or_orphan_classification_is_rejected(self):
        with self.assertRaises(render.RenderError):
            render.render(GW, [frag()], {})
        with self.assertRaises(render.RenderError):
            render.render(GW, [frag()], {"grafana": klass(), "other": klass("other")})

    def test_fragment_cannot_join_execute_or_deferred_groups(self):
        for g in ("execute", "identity"):
            with self.assertRaises(render.RenderError):
                render.render(GW, [frag(groups=[g])], {"grafana": klass()})

    def test_guard_hides_the_dangerous_arguments(self):
        c = klass(write=["scale"], guards=[{"tool": "scale", "class": "destructive", "when": "mcp.tool.arguments.replicas == 0"}])
        expr = routes_of(render.render(GW, [frag()], {"grafana": c}))["ops"]["policies"]["mcpAuthorization"]["rules"][0]["allow"]
        self.assertIn('mcp.tool.name == "scale" && !(mcp.tool.arguments.replicas == 0)', expr)


class Classification(unittest.TestCase):
    def parse(self, text, name="t.yaml"):
        import tempfile
        with tempfile.TemporaryDirectory() as d:
            p = Path(d) / name
            p.write_text(text)
            return cl.parse(p)

    def test_committed_files_are_valid_and_onemind_is_read_only(self):
        classes = cl.load_all(render.CLASSIFICATION)
        self.assertIn("onemind", classes)
        # Every onemind-mcp tool is read only except the tile38 adapter's gated write and destructive tools.
        self.assertEqual(len(classes["onemind"]["read"]), 26)
        self.assertEqual(classes["onemind"]["write"], ["tile38_set_point"])
        self.assertEqual(classes["onemind"]["destructive"], ["tile38_delete_object", "tile38_drop_collection"])
        self.assertEqual(classes["onemind"]["physical"], [])

    def test_rejects_bad_files(self):
        bad = [
            "target: t\nread: [a]\nwrite: [a]\n",  # two classes
            "target: t\nread: [get_*]\ndestructive: [get_x]\n",  # gated name swallowed by a glob
            "target: t\nread: [a]\nproposed: [a]\n",  # unreviewed proposal
            "target: t\nread: [a]\nextra: 1\n",
            "target: other\nread: [a]\n",  # file name mismatch
            "target: t\nread: [\"a b\"]\n",
            "target: t\nwrite: [s]\nguards: [{tool: nope, class: destructive, when: x}]\n",
        ]
        for text in bad:
            with self.subTest(text=text), self.assertRaises(cl.PolicyError):
                self.parse(text.replace("\\n", "\n"))

    def test_unlisted_is_destructive(self):
        c = klass(read=["get_*"], write=["set_x"])
        self.assertEqual(cl.classify_name("get_pod", c), "read")
        self.assertEqual(cl.classify_name("set_x", c), "write")
        self.assertEqual(cl.classify_name("anything_else", c), "destructive")

    def test_classify_tools_heuristics_fail_closed(self):
        p = classify.propose
        self.assertEqual(p("pods_list"), "read")
        self.assertEqual(p("get_dashboard"), "read")
        self.assertEqual(p("create_annotation"), "write")
        self.assertEqual(p("deleteDashboard"), "destructive")
        self.assertEqual(p("resources_scale"), "destructive")
        self.assertEqual(p("lock_door"), "physical")
        self.assertEqual(p("get_and_delete"), "destructive")
        self.assertEqual(p("frobnicate"), "destructive")
        self.assertEqual(p("frobnicate", {"readOnlyHint": True}), "read")

    def test_classify_tools_drift(self):
        c = klass(read=["get_*"], write=["gone"])
        new, stale = classify.compare(["get_a", "mystery"], c)
        self.assertEqual((new, stale), (["mystery"], ["gone"]))


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
