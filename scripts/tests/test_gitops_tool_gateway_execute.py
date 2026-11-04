"""Phase 4 execute gate wiring: renderer output, the execute-gate manifests, NATS objects and the held Hermes entry.

The live behaviour of the extAuthz rule set (list works, call needs the gate, gate down means 403) was verified against the
pinned agentgateway image on 2026-11-03; docs/operations/tool-gateway.md records the steps. These tests pin the Git side.
"""
from pathlib import Path
import copy
import hashlib
import importlib.util
import io
import re
import sys
import unittest

import yaml

ROOT = Path(__file__).resolve().parents[2]
BASE = ROOT / "platform/tool-gateway"
GATE = BASE / "manifests/execute-gate"
spec = importlib.util.spec_from_file_location("render_gateway_config", ROOT / "scripts/gateway/render_gateway_config.py")
render = importlib.util.module_from_spec(spec)
spec.loader.exec_module(render)

# The sealed files the gate cannot start without. replicas must be 1 exactly when all three exist.
GATE_SEALED = ("gate-matrix", "gate-nats", "ghcr")


def routes(values):
    return {r["name"]: r for r in values["config"]["binds"][0]["listeners"][0]["routes"]}


def docs(path):
    return [d for d in yaml.safe_load_all(Path(path).read_text()) if d]


def deployment():
    return next(d for d in docs(GATE / "deployment.yaml") if d["kind"] == "Deployment")


class Renderer(unittest.TestCase):
    def setUp(self):
        self.gw, self.targets, self.classes = render.load()
        self.values = render.render(self.gw, self.targets, self.classes)
        self.routes = routes(self.values)

    def test_generated_files_are_current(self):
        for path, text in render.outputs().items():
            self.assertEqual(path.read_text(), text, f"{path.relative_to(ROOT)} is stale: run render_gateway_config.py")

    def test_only_the_execute_route_calls_the_gate(self):
        for name, r in self.routes.items():
            if name == "execute":
                self.assertIn("extAuthz", r["policies"])
            else:
                self.assertNotIn("extAuthz", r["policies"], name)
                self.assertFalse(any("extauthz" in str(rule) for rule in r["policies"]["mcpAuthorization"]["rules"]), name)

    def test_gate_rules_list_gated_tools_and_gate_calls_only(self):
        ex = self.routes["execute"]["policies"]
        rules = ex["mcpAuthorization"]["rules"]
        # Visibility is allow rules only; the single require rule is keyed on the method, because a require on
        # extauthz.decision alone is evaluated for tools/list too, where no decision exists, and empties the list.
        self.assertEqual(rules[-1], {"require": 'mcp.methodName != "tools/call" || extauthz.decision == true'})
        self.assertTrue(all(set(r) == {"allow"} for r in rules[:-1]))
        active = {t["name"] for t in self.targets if render.is_active(t)}
        gated = {n for n, c in self.classes.items() if n in active and render.classification.has_gated(c)}
        self.assertEqual({re.search(r'target == "(\w+)"', r["allow"]).group(1) for r in rules[:-1]}, gated)
        for r in rules[:-1]:
            body = r["allow"]
            self.assertNotIn("require", r)
            # no exposed (read/write) tool name may appear in an execute allow rule
            tgt = re.search(r'target == "(\w+)"', body).group(1)
            for cls in render.classification.EXPOSED:
                for n in self.classes[tgt][cls]:
                    self.assertNotIn('"%s"' % n, body)

    def test_gate_is_fail_closed_and_asked_for_tool_calls_only(self):
        cond = self.routes["execute"]["policies"]["extAuthz"]["conditional"]
        self.assertEqual(len(cond), 1)
        c = cond[0]
        self.assertEqual(c["condition"], "has(mcp.tool)")
        self.assertEqual(c["failureMode"], {"denyWithStatus": 403})
        self.assertNotIn("cache", c)  # a cached allow would outlive the single-use approval
        self.assertEqual(c["host"], "execute-gate.tool-gateway.svc.cluster.local:8080")
        http = c["protocol"]["http"]
        self.assertEqual(http["path"], '"/access/v1/evaluation"')
        self.assertEqual(http["metadata"], {"decision": "json(response.body).decision"})
        # a tool without `arguments` must not make the CEL error (that would deny before the gate is asked)
        self.assertIn("has(mcp.tool.arguments)", http["body"])
        self.assertIn("apiKey.client", http["body"])
        self.assertIn('mcp.tool.target + "/" + mcp.tool.name', http["body"])
        svc = next(d for d in docs(GATE / "deployment.yaml") if d["kind"] == "Service")
        self.assertEqual(c["host"], "%s.tool-gateway.svc.cluster.local:%d" % (svc["metadata"]["name"], svc["spec"]["ports"][0]["port"]))

    def test_gate_absent_changes_nothing_outside_execute(self):
        """The execute-gate pod and its Secrets are not needed by the gateway: the other routes and the pod env are
        exactly what they are with no gate at all (the route is deny-all, as before Phase 4)."""
        gw = copy.deepcopy(self.gw)
        del gw["groups"]["execute"]["authz"]
        baseline = render.render(gw, self.targets, self.classes)
        self.assertEqual(baseline["extraEnv"], self.values["extraEnv"])
        base_routes = routes(baseline)
        self.assertEqual(set(base_routes), set(self.routes))
        for name in self.routes:
            if name != "execute":
                self.assertEqual(self.routes[name], base_routes[name], name)
        self.assertEqual(base_routes["execute"]["policies"]["mcpAuthorization"], {"rules": [{"require": "false"}]})
        self.assertNotIn("extAuthz", base_routes["execute"]["policies"])
        for env in self.values["extraEnv"]:
            self.assertNotIn(env["valueFrom"]["secretKeyRef"]["name"], GATE_SEALED)

    def test_authz_is_validated(self):
        for bad in ({"host": "no-port"}, {"host": "x:8080", "extra": 1}, {}, "x"):
            gw = copy.deepcopy(self.gw)
            gw["groups"]["execute"]["authz"] = bad
            with self.assertRaises(render.RenderError):
                render.render(gw, self.targets, self.classes)

    def test_only_the_execute_group_may_have_authz_effect(self):
        gw = copy.deepcopy(self.gw)
        gw["groups"]["ops"]["authz"] = {"host": "x.y:8080"}
        out = routes(render.render(gw, self.targets, self.classes))
        self.assertNotIn("extAuthz", out["ops"]["policies"])

    def test_physical_tools_reach_the_gate_policy(self):
        classes = copy.deepcopy(self.classes)
        classes["onemind"]["physical"] = ["lock_*"]
        text = render.gate_policy_text(self.targets, classes)
        policy = yaml.safe_load(yaml.safe_load(text)["data"]["policy.yaml"])
        self.assertEqual(policy["physical"], ["onemind/lock_*"])
        self.assertIn("kubernetes/resources_delete", policy["destructive"])

    def test_policy_hash_annotation_matches_the_configmap(self):
        sha = hashlib.sha256((GATE / "policy.configmap.yaml").read_text().encode()).hexdigest()
        ann = deployment()["spec"]["template"]["metadata"]["annotations"]["onemind.io/policy-sha256"]
        self.assertEqual(ann, sha)


class GateManifests(unittest.TestCase):
    def test_replicas_follow_the_sealed_secrets(self):
        present = [(BASE / "manifests" / f"{n}.sealedsecret.yaml").is_file() for n in GATE_SEALED]
        want = 1 if all(present) else 0
        self.assertEqual(deployment()["spec"]["replicas"], want,
                         "replicas must be 1 exactly when gate-matrix, gate-nats and ghcr are all sealed "
                         "(scripts/gateway/keys-session.sh sets it); the pod cannot start without them")

    def test_pod_is_pinned_restricted_and_bounded(self):
        d = deployment()
        pod = d["spec"]["template"]["spec"]
        c = pod["containers"][0]
        self.assertRegex(c["image"], r"^ghcr\.io/onemindos/execute-gate:sha-[0-9a-f]{40}@sha256:[0-9a-f]{64}$")
        self.assertEqual(c["securityContext"]["readOnlyRootFilesystem"], True)
        self.assertEqual(c["securityContext"]["capabilities"], {"drop": ["ALL"]})
        self.assertEqual(c["securityContext"]["allowPrivilegeEscalation"], False)
        self.assertTrue(pod["securityContext"]["runAsNonRoot"])
        self.assertEqual(pod["securityContext"]["seccompProfile"]["type"], "RuntimeDefault")
        self.assertFalse(pod["automountServiceAccountToken"])
        self.assertEqual(pod["nodeSelector"], {"onemind.io/compute-role": "ecosystem"})
        self.assertEqual(pod["imagePullSecrets"], [{"name": "ghcr"}])
        self.assertTrue(c["resources"]["requests"] and c["resources"]["limits"])
        self.assertEqual(d["spec"]["strategy"]["type"], "Recreate")  # two pods would run two Matrix sync loops

    def test_secret_references_and_env_match_the_service(self):
        c = deployment()["spec"]["template"]["spec"]["containers"][0]
        env = {e["name"]: e for e in c["env"]}
        refs = {n: (e["valueFrom"]["secretKeyRef"]["name"], e["valueFrom"]["secretKeyRef"]["key"])
                for n, e in env.items() if "valueFrom" in e}
        self.assertEqual(refs, {
            "MATRIX_ACCESS_TOKEN": ("gate-matrix", "token"),
            "APPROVALS_ROOM_ID": ("gate-matrix", "room_id"),
            "NATS_USER": ("gate-nats", "user"),
            "NATS_PASSWORD": ("gate-nats", "password"),
        })
        self.assertEqual(env["OWNER_MXID"]["value"], "@zeus:onemindos.dev")
        self.assertEqual(env["STORE"]["value"], "nats")
        self.assertEqual(env["KV_BUCKET"]["value"], "exec-approvals")
        self.assertEqual(env["ALLOWED_SUBJECTS"]["value"], "hermes")
        self.assertEqual(env["POLICY_FILE"]["value"], "/etc/execute-gate/policy.yaml")
        self.assertTrue(env["NATS_URL"]["value"].startswith("tls://nats.nats.svc.cluster.local:"))
        self.assertEqual(env["SSL_CERT_FILE"]["value"], "/etc/nats-ca/ca.crt")
        # every variable the service requires or the policy needs is set (services/execute-gate/README.md)
        for required in ("APPROVALS_ROOM_ID", "MATRIX_ACCESS_TOKEN", "POLICY_FILE"):
            self.assertIn(required, env)
        self.assertNotIn("GATE_AUTH_TOKEN", env)  # NetworkPolicy is the caller check; no shared token to rotate

    def test_service_and_probes(self):
        c = deployment()["spec"]["template"]["spec"]["containers"][0]
        self.assertEqual(c["readinessProbe"]["httpGet"]["path"], "/readyz")
        self.assertEqual(c["livenessProbe"]["httpGet"]["path"], "/healthz")
        svc = next(d for d in docs(GATE / "deployment.yaml") if d["kind"] == "Service")
        self.assertEqual(svc["spec"]["selector"], {"app.kubernetes.io/name": "execute-gate"})
        self.assertEqual(svc["spec"]["ports"][0]["port"], 8080)

    def test_network_policies_both_sides(self):
        gate = docs(GATE / "network.yaml")[0]["spec"]
        peers = gate["ingress"]
        self.assertEqual(len(peers), 1)
        self.assertEqual(peers[0]["from"], [{"podSelector": {"matchLabels": {
            "app.kubernetes.io/name": "agentgateway-standalone", "app.kubernetes.io/instance": "agentgateway"}}}])
        self.assertEqual(peers[0]["ports"], [{"protocol": "TCP", "port": 8080}])
        ports = {(r["to"][0]["namespaceSelector"]["matchLabels"]["kubernetes.io/metadata.name"], p["port"])
                 for r in gate["egress"] for p in r["ports"]}
        self.assertEqual(ports, {("kube-system", 53), ("matrix", 8008), ("nats", 4222)})
        egress = docs(BASE / "manifests/egress-execute-gate.yaml")[0]["spec"]
        self.assertEqual(egress["policyTypes"], ["Egress"])
        self.assertEqual(egress["podSelector"]["matchLabels"]["app.kubernetes.io/name"], "agentgateway-standalone")
        self.assertEqual(egress["egress"][0]["to"], [{"podSelector": {"matchLabels": {"app.kubernetes.io/name": "execute-gate"}}}])
        self.assertEqual(egress["egress"][0]["ports"], [{"protocol": "TCP", "port": 8080}])

    def test_nats_admits_only_the_gate_pod_from_tool_gateway(self):
        clients = next(d for d in docs(ROOT / "platform/nats/network-policy.yaml") if d["metadata"]["name"] == "nats-clients")
        peers = [f for r in clients["spec"]["ingress"] for f in r["from"]]
        mine = [f for f in peers if f.get("namespaceSelector", {}).get("matchLabels", {}).get("kubernetes.io/metadata.name") == "tool-gateway"]
        self.assertEqual(len(mine), 1)
        self.assertEqual(mine[0]["podSelector"]["matchLabels"], {"app.kubernetes.io/name": "execute-gate"})

    def test_kv_bucket_and_audit_subject(self):
        kv = yaml.safe_load((ROOT / "platform/nats/streams/kv-exec-approvals.yaml").read_text())
        self.assertEqual(kv["kind"], "KeyValue")
        self.assertEqual(kv["spec"]["bucket"], "exec-approvals")
        self.assertEqual(kv["spec"]["replicas"], 3)
        self.assertEqual(kv["spec"]["storage"], "file")
        self.assertTrue(kv["spec"]["preventDelete"])
        self.assertEqual(kv["spec"]["history"], 1)  # the gate compares revisions, never reads history
        evt = yaml.safe_load((ROOT / "platform/nats/streams/evt.yaml").read_text())
        self.assertIn("evt.toolgw.>", evt["spec"]["subjects"])
        # the subjects the service publishes (evt.{domain}.{entity}.{event}) fall under it
        audit = (ROOT / "services/execute-gate/execute_gate/audit.py").read_text()
        self.assertIn('f"evt.toolgw.approval.{event}"', audit)


class HeldHermesEntry(unittest.TestCase):
    def test_execute_server_is_in_the_defaults_but_held(self):
        conf = yaml.safe_load((ROOT / "platform/hermes/onemind-hermes/config.defaults.yaml").read_text())
        ex = conf["mcp_servers"]["execute"]
        self.assertEqual(ex["url"], "http://agentgateway.tool-gateway.svc.cluster.local:8080/execute/mcp")
        self.assertEqual(ex["trust"], "untrusted")
        self.assertEqual(ex["headers"], {"Authorization": "Bearer ${ONEMIND_GATEWAY_KEY}"})
        applier = (ROOT / "platform/hermes/onemind-hermes/mcp/apply_mcp.py").read_text()
        self.assertRegex(applier, r'HOLD = \["execute"\]')
        self.assertIn('name in HOLD and name not in INCLUDE', applier)


class GateAcceptsWhatTheGatewaySends(unittest.TestCase):
    """The AuthZEN body captured from the pinned gateway image (2026-11-03) is understood by the service."""

    BODY = {"subject": {"type": "agent", "id": "hermes"}, "resource": {"type": "tool", "id": "kubernetes/pods_delete"},
            "action": {"name": "tools/call"}, "context": {"arguments": {"name": "x", "namespace": "y"}}}

    def setUp(self):
        sys.path.insert(0, str(ROOT / "services/execute-gate"))
        self.addCleanup(sys.path.remove, str(ROOT / "services/execute-gate"))

    def test_decide_denies_and_opens_a_request(self):
        from execute_gate.audit import Audit
        from execute_gate.gate import Gate
        from execute_gate.policy import Policy
        from execute_gate.server import decide
        from execute_gate.store import MemoryStore

        class Sent:
            def __init__(self):
                self.text = []

            def send(self, t):
                self.text.append(t)

        sent = Sent()
        gate = Gate(MemoryStore(), Policy(), "@zeus:onemindos.dev", "!r:onemindos.dev", "@gate:onemindos.dev", sent,
                    Audit(None, io.StringIO()))
        out = decide(gate, copy.deepcopy(self.BODY))
        self.assertFalse(out["decision"])
        self.assertEqual(out["context"]["reason"], "approval_required")
        self.assertRegex(out["context"]["request_id"], r"^[0-9A-F]{4}$")
        self.assertIn("kubernetes/pods_delete", sent.text[0])
        # a call with no arguments hashes like an empty object (the gateway sends {} for it)
        empty = copy.deepcopy(self.BODY)
        empty["context"]["arguments"] = {}
        self.assertFalse(decide(gate, empty)["decision"])


if __name__ == "__main__":
    unittest.main()
