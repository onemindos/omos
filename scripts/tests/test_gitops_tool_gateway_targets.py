"""PR B: gateway target fragments that stay off until their key is sealed, and the owner key session script.

A target with `enabled_by` is rendered only while its SealedSecret file is committed, and its deployable manifests (backend, egress,
namespace rules) wait in platform/tool-gateway/pending/<name>/ until scripts/gateway/keys-session.sh moves them into manifests/.
"""
from pathlib import Path
import hashlib
import http.server
import importlib.util
import json
import os
import re
import shutil
import subprocess
import sys
import tempfile
import threading
import unittest

import yaml

ROOT = Path(__file__).resolve().parents[2]
BASE = ROOT / "platform/tool-gateway"
SCRIPT = ROOT / "scripts/gateway/keys-session.sh"
BOOT = ROOT / "scripts/gateway/gate_bootstrap.py"
spec = importlib.util.spec_from_file_location("render_gateway_config", ROOT / "scripts/gateway/render_gateway_config.py")
render = importlib.util.module_from_spec(spec)
spec.loader.exec_module(render)

NEW_TARGETS = {"grafana": "ops", "argocd": "ops", "github": "dev", "doapps": "cloud", "dodoks": "cloud", "donetworking": "cloud",
               "cloudflare": "cloud", "n8n": "work", "homeassistant": "home", "cloudtak": "tak"}


def fragments():
    return {p.stem: yaml.safe_load(p.read_text()) for p in sorted((BASE / "targets").glob("*.yaml"))}


def pending_name(secret):
    return secret.removeprefix("gateway-target-") if hasattr(str, "removeprefix") else secret[len("gateway-target-"):]


class Pending(unittest.TestCase):
    def test_pending_dir_exists_exactly_while_the_secret_is_not_sealed(self):
        for name, frag in fragments().items():
            for secret in frag.get("enabled_by") or []:
                sealed = (BASE / f"manifests/{secret}.sealedsecret.yaml").is_file()
                pend = (BASE / "pending" / pending_name(secret) / "manifests").is_dir()
                self.assertNotEqual(sealed, pend, f"{name}: {secret} sealed={sealed} but pending manifests present={pend}")

    def test_new_targets_are_gated_and_declare_every_secret_they_use(self):
        frags = fragments()
        for name, group in NEW_TARGETS.items():
            f = frags[name]
            self.assertEqual(f["groups"], [group], name)
            secrets = {r["secret"] for r in (f.get("env") or {}).values()}
            self.assertTrue(f.get("enabled_by"), f"{name}: needs enabled_by")
            self.assertTrue(secrets <= set(f["enabled_by"]), name)
            for s in f["enabled_by"]:
                self.assertRegex(s, r"^gateway-target-[a-z0-9]+$")

    def test_nothing_new_is_rendered_while_keys_are_missing(self):
        gw, targets, classes = render.load()
        values = render.render(gw, targets, classes)
        text = yaml.safe_dump(values)
        for name in NEW_TARGETS:
            self.assertNotIn(f"name: {name}\n", text)
        self.assertEqual({e["name"] for e in values["extraEnv"] if "TARGET_" in e["name"]},
                         {v for t in targets if render.is_active(t) for v in (t.get("env") or {})})
        # an absent key must not leave a Secret reference behind: the gateway pod would not start
        for e in values["extraEnv"]:
            ref = e["valueFrom"]["secretKeyRef"]
            self.assertTrue(ref["name"] == "gateway-clients" or (BASE / f"manifests/{ref['name']}.sealedsecret.yaml").is_file(), ref)

    def test_everything_enabled_renders_and_reaches_the_expected_groups(self):
        gw, targets, classes = render.load()
        orig = render.sealed_exists
        render.sealed_exists = lambda s: True
        try:
            values = render.render(gw, targets, classes)
            policy = yaml.safe_load(yaml.safe_load(render.gate_policy_text(targets, classes))["data"]["policy.yaml"])
        finally:
            render.sealed_exists = orig
        routes = {r["name"]: r for r in values["config"]["binds"][0]["listeners"][0]["routes"]}
        for name, group in NEW_TARGETS.items():
            names = [t["name"] for t in routes[group]["backends"][0]["mcp"]["targets"]]
            self.assertIn(name, names)
        self.assertIn("homeassistant/Hass*", policy["physical"])
        self.assertIn("cloudflare/execute", policy["destructive"])
        for g in ("ops", "dev", "cloud", "home", "work", "tak"):
            self.assertTrue(routes[g]["policies"]["mcpAuthorization"]["rules"], g)
        env = {e["name"]: e["valueFrom"]["secretKeyRef"] for e in values["extraEnv"]}
        self.assertEqual(env["TARGET_ARGOCD_MCP_AUTH"], {"name": "gateway-target-argocd", "key": "mcp_auth"})
        self.assertEqual(env["TARGET_DOAPPS_TOKEN"]["name"], "gateway-target-digitalocean")

    def test_pending_manifests_are_restricted_pinned_and_namespaced(self):
        for pend in sorted((BASE / "pending").glob("*/manifests")):
            for path in sorted(pend.rglob("*.yaml")):
                for d in yaml.safe_load_all(path.read_text()):
                    if not d:
                        continue
                    self.assertIn(d["metadata"].get("namespace"), {"tool-gateway", "argocd", "n8n", "cloudtak"}, path)
                    if d["kind"] == "Deployment":
                        pod = d["spec"]["template"]["spec"]
                        c = pod["containers"][0]
                        self.assertRegex(c["image"], r"@sha256:[0-9a-f]{64}$", path)
                        self.assertTrue(c["resources"]["requests"] and c["resources"]["limits"], path)
                        self.assertTrue(c["securityContext"]["readOnlyRootFilesystem"], path)
                        self.assertEqual(c["securityContext"]["capabilities"], {"drop": ["ALL"]}, path)
                        self.assertTrue(pod["securityContext"]["runAsNonRoot"], path)
                        self.assertFalse(pod["automountServiceAccountToken"], path)
                        self.assertEqual(pod["nodeSelector"], {"onemind.io/compute-role": "ecosystem"})
                        self.assertEqual(d["spec"]["replicas"], 1)
                    if d["kind"] == "NetworkPolicy":
                        self.assertNotEqual(d["spec"]["podSelector"], {}, path)

    def test_backends_match_their_fragments_and_secrets(self):
        frags = fragments()
        for name in ("grafana", "argocd"):
            dep = next(d for d in yaml.safe_load_all((BASE / f"pending/{name}/manifests/backends/{name}/deployment.yaml").read_text())
                       if d and d["kind"] == "Deployment")
            svc = next(d for d in yaml.safe_load_all((BASE / f"pending/{name}/manifests/backends/{name}/deployment.yaml").read_text())
                       if d and d["kind"] == "Service")
            host = frags[name]["target"]["mcp"]["host"]
            self.assertEqual(host, "http://mcp-%s.tool-gateway.svc.cluster.local:%d/mcp" % (name, svc["spec"]["ports"][0]["port"]))
            refs = {e["valueFrom"]["secretKeyRef"]["key"] for e in dep["spec"]["template"]["spec"]["containers"][0]["env"]
                    if "valueFrom" in e}
            keys = sealed_keys_in_script()[f"gateway-target-{name}"]
            self.assertTrue(refs <= keys, (name, refs, keys))
        # the grafana backend is read-only until someone removes the flag on purpose
        g = yaml.safe_load((BASE / "pending/grafana/manifests/backends/grafana/deployment.yaml").read_text().split("\n---")[0])
        self.assertIn("-disable-write", g["spec"]["template"]["spec"]["containers"][0]["args"])
        a = yaml.safe_load((BASE / "pending/argocd/manifests/backends/argocd/deployment.yaml").read_text().split("\n---")[0])
        env = {e["name"]: e.get("value") for e in a["spec"]["template"]["spec"]["containers"][0]["env"]}
        self.assertEqual(env["MCP_READ_ONLY"], "true")

    def test_every_fragment_secret_key_is_sealed_by_the_script(self):
        script_keys = sealed_keys_in_script()
        for name, f in fragments().items():
            for ref in (f.get("env") or {}).values():
                if ref["secret"].startswith("gateway-target-") and f.get("enabled_by"):
                    self.assertIn(ref["key"], script_keys[ref["secret"]], f"{name}: {ref}")

    def test_argocd_account_legacy_is_read_only(self):
        vals = yaml.safe_load((ROOT / "platform/values/argocd.yaml").read_text())
        cm = vals["configs"]["cm"]
        self.assertEqual(cm["accounts.legacy"], "apiKey")
        policy = [l.strip() for l in vals["configs"]["rbac"]["policy.csv"].splitlines() if "legacy" in l and not l.strip().startswith("#")]
        for line in policy:
            self.assertNotRegex(line, r"sync|action|delete|override|create|update|\*, \*|exec")
        self.assertIn("g, legacy, role:legacy", policy)

    def test_karakeep_is_blocked_not_wired(self):
        self.assertNotIn("karakeep", fragments())
        self.assertTrue((BASE / "pending/karakeep/README.md").is_file())
        self.assertFalse((BASE / "pending/karakeep/manifests").exists())

    def test_cloudtak_spec_is_get_only_and_in_sync(self):
        cm = yaml.safe_load((BASE / "manifests/openapi-cloudtak.configmap.yaml").read_text())
        spec_ = json.loads(cm["data"]["openapi.yaml"])
        for path, item in spec_["paths"].items():
            self.assertEqual(set(item), {"get"}, path)
        vols = yaml.safe_load((BASE / "values.yaml").read_text())
        self.assertIn({"name": "openapi-cloudtak", "configMap": {"name": "gateway-openapi-cloudtak"}}, vols["extraVolumes"])
        out = subprocess.run([sys.executable, str(ROOT / "scripts/gateway/trim_openapi_get.py"), "--check"], capture_output=True, text=True)
        self.assertEqual(out.returncode, 0, out.stdout + out.stderr)


def sealed_keys_in_script():
    """{secret name: set(keys)} read from the `seal <name> <ns> <file> ... -v/-f <key>` calls of keys-session.sh."""
    keys = {}
    text = SCRIPT.read_text().replace("\\\n", " ")
    for m in re.finditer(r"^\s*seal (gateway-target-[a-z0-9]+) tool-gateway \S+(.*)$", text, re.M):
        keys.setdefault(m.group(1), set()).update(re.findall(r"-[vf] ([A-Za-z0-9._-]+) ", m.group(2) + " "))
    return keys


@unittest.skipUnless(shutil.which("kubectl") and shutil.which("bash"), "needs kubectl (client-side only) and bash")
class KeysSessionDryRun(unittest.TestCase):
    SECRET_MARKERS = ["MARKER-GRAFANA", "MARKER-ARGOCD", "MARKER-GITHUB", "MARKER-DO", "MARKER-CF", "MARKER-N8N-MCP", "MARKER-N8N-API",
                      "MARKER-HA", "MARKER-KARAKEEP", "MARKER-PORTAL", "MARKER-P12PASS", "MARKER-SSHKEYFILE", "MARKER-P12FILE"]

    def tree(self):
        h = hashlib.sha256()
        for p in sorted(list((BASE).rglob("*")) + [ROOT / "platform/hermes/onemind-hermes/kustomization.yaml"]):
            if p.is_file():
                h.update(str(p).encode() + p.read_bytes())
        return h.hexdigest()

    def answers(self, tmp):
        pad = lambda tag, n=45: tag + "x" * (n - len(tag))  # noqa: E731
        (tmp / "key").write_text("MARKER-SSHKEYFILE\n")
        (tmp / "p12").write_text("MARKER-P12FILE\n")
        lines = [
            "", "",                                   # ghcr: continue, reuse
            "",                                       # nats
            "",                                       # gate
            "", pad("MARKER-GRAFANA"),
            "", pad("MARKER-ARGOCD"),
            "", pad("MARKER-GITHUB"),
            "", pad("MARKER-DO"),
            "", pad("MARKER-CF"),
            "", pad("MARKER-N8N-MCP"), pad("MARKER-N8N-API"),
            "", pad("MARKER-HA", 120),
            "", "",                                   # cloudtak reuse
            "", pad("MARKER-KARAKEEP"),
            "", pad("MARKER-PORTAL"),
            "", str(tmp / "p12"), "MARKER-P12PASS",
            "", str(tmp / "key"), "",                 # legacy-ssh, no known_hosts
            "", "", "", str(tmp / "key"),             # legacy-export: ssh mode, default repo
        ]
        return "\n".join(lines) + "\n"

    def run_script(self, tmp, args, stdin):
        return subprocess.run(["bash", str(SCRIPT), "--dry-run", "--out", str(tmp / "out")] + args, input=stdin, capture_output=True,
                              text=True, env={**os.environ, "KUBE_CONTEXT": "does-not-exist"}, cwd=str(ROOT))

    def test_all_steps_dry_run(self):
        before = self.tree()
        with tempfile.TemporaryDirectory() as t:
            tmp = Path(t)
            r = self.run_script(tmp, ["all"], self.answers(tmp))
            self.assertEqual(r.returncode, 0, r.stderr)
            repo = tmp / "out/repo"
            out_text = r.stdout + r.stderr
            written = {p.name for p in (repo / "platform/tool-gateway/manifests").glob("*.sealedsecret.yaml")}
            for want in ("ghcr", "gate-nats", "gate-matrix", "gateway-target-grafana", "gateway-target-argocd", "gateway-target-github",
                         "gateway-target-digitalocean", "gateway-target-cloudflare", "gateway-target-n8n", "gateway-target-homeassistant",
                         "gateway-target-cloudtak", "gateway-target-karakeep", "gateway-target-takportal", "gateway-target-takserver"):
                self.assertIn(want + ".sealedsecret.yaml", written)
            hm = repo / "platform/hermes/onemind-hermes"
            self.assertTrue((hm / "legacy-ssh.sealedsecret.yaml").is_file() and (hm / "export.sealedsecret.yaml").is_file())
            kust = (hm / "kustomization.yaml").read_text()
            self.assertIn("  - legacy-ssh.sealedsecret.yaml", kust)
            self.assertIn("  - export.sealedsecret.yaml", kust)
            # activation moved the pending manifests, the gate turned on, the render is current, karakeep is not wired
            tg = repo / "platform/tool-gateway"
            self.assertEqual(sorted(p.name for p in (tg / "pending").iterdir()), ["karakeep"])
            self.assertTrue((tg / "manifests/backends/grafana/deployment.yaml").is_file())
            dep = yaml.safe_load((tg / "manifests/execute-gate/deployment.yaml").read_text().split("\n---")[0])
            self.assertEqual(dep["spec"]["replicas"], 1)
            chk = subprocess.run([sys.executable, str(repo / "scripts/gateway/render_gateway_config.py"), "--check"], capture_output=True, text=True)
            self.assertEqual(chk.returncode, 0, chk.stderr)
            vals = yaml.safe_load((tg / "values.generated.yaml").read_text())
            names = {t_["name"] for r_ in vals["config"]["binds"][0]["listeners"][0]["routes"] for b in r_["backends"] for t_ in b["mcp"]["targets"]}
            self.assertTrue(set(NEW_TARGETS) <= names)
            # no value anywhere: not in the output, not in any written file
            blobs = [out_text] + [p.read_text(errors="ignore") for p in repo.rglob("*.sealedsecret.yaml")]
            for marker in self.SECRET_MARKERS:
                for b in blobs:
                    self.assertNotIn(marker, b)
            self.assertIn("bytes", out_text)  # lengths are reported
        self.assertEqual(self.tree(), before, "a dry run must not touch the repository")

    def test_one_target_can_be_rerun(self):
        with tempfile.TemporaryDirectory() as t:
            tmp = Path(t)
            token = "MARKER-GRAFANA" + "x" * 30
            for _ in range(2):
                r = self.run_script(tmp, ["grafana"], "\n" + token + "\n")
                self.assertEqual(r.returncode, 0, r.stderr)
            self.assertTrue((tmp / "out/repo/platform/tool-gateway/manifests/gateway-target-grafana.sealedsecret.yaml").is_file())

    def test_skip_and_short_values(self):
        with tempfile.TemporaryDirectory() as t:
            tmp = Path(t)
            r = self.run_script(tmp, ["github"], "s\n")
            self.assertEqual(r.returncode, 0)
            self.assertFalse((tmp / "out/repo/platform/tool-gateway/manifests/gateway-target-github.sealedsecret.yaml").exists())
            r = self.run_script(tmp, ["github"], "\nshort\n")
            self.assertFalse((tmp / "out/repo/platform/tool-gateway/manifests/gateway-target-github.sealedsecret.yaml").exists())
            self.assertNotIn("short\n", r.stdout)

    def test_unknown_step_is_refused(self):
        with tempfile.TemporaryDirectory() as t:
            r = self.run_script(Path(t), ["nextcloud"], "")
            self.assertNotEqual(r.returncode, 0)
            self.assertIn("unknown step", r.stderr)


class ScriptHygiene(unittest.TestCase):
    def test_values_never_reach_argv_history_or_stdout(self):
        text = SCRIPT.read_text()
        self.assertIn("read -r -s", text)
        for bad in ("--from-literal", "echo ", " -p ", "set -x", "history"):
            code = [l for l in text.splitlines() if not l.lstrip().startswith("#")]
            self.assertFalse(any(bad in l for l in code if "printf '%s' \"$pat\"" not in l and "say " not in l and "-p \"" not in l),
                             f"unexpected {bad!r}")
        self.assertIn("kubeseal --cert", text)
        self.assertIn("--scope namespace-wide", text)
        self.assertNotIn("nextcloud", text.lower())
        self.assertNotIn("paperclip", text.lower())
        r = subprocess.run(["bash", "-n", str(SCRIPT)], capture_output=True, text=True)
        self.assertEqual(r.returncode, 0, r.stderr)


class FakeSynapse(http.server.BaseHTTPRequestHandler):
    owner_exists = True
    encrypted = False
    seen = []

    def log_message(self, *a):
        pass

    def _send(self, code, obj):
        raw = json.dumps(obj).encode()
        self.send_response(code)
        self.send_header("Content-Type", "application/json")
        self.send_header("Content-Length", str(len(raw)))
        self.end_headers()
        self.wfile.write(raw)

    def do_GET(self):
        FakeSynapse.seen.append(("GET", self.path))
        if "/profile/" in self.path:
            return self._send(200 if FakeSynapse.owner_exists else 404, {})
        if "m.room.encryption" in self.path:
            return self._send(200 if FakeSynapse.encrypted else 404, {})
        self._send(404, {})

    def do_POST(self):
        body = json.loads(self.rfile.read(int(self.headers.get("Content-Length", 0))) or b"{}")
        FakeSynapse.seen.append(("POST", self.path, body))
        if self.path.endswith("/login"):
            return self._send(200, {"access_token": "syt_FAKE_TOKEN_VALUE"})
        if self.path.endswith("/createRoom"):
            if self.headers.get("Authorization") != "Bearer syt_FAKE_TOKEN_VALUE":
                return self._send(401, {})
            return self._send(200, {"room_id": "!fakeroom:onemindos.dev"})
        self._send(404, {})


class GateBootstrap(unittest.TestCase):
    def setUp(self):
        FakeSynapse.owner_exists, FakeSynapse.encrypted, FakeSynapse.seen = True, False, []
        self.srv = http.server.ThreadingHTTPServer(("127.0.0.1", 0), FakeSynapse)
        threading.Thread(target=self.srv.serve_forever, daemon=True).start()
        self.addCleanup(self.srv.server_close)
        self.addCleanup(self.srv.shutdown)
        self.base = "http://127.0.0.1:%d" % self.srv.server_address[1]

    def run_boot(self, register, pw="PW-SECRET-VALUE-123"):
        cfg = {"user": "gate", "password": pw, "owner": "@zeus:onemindos.dev", "base": self.base, "register_cmd": register}
        return subprocess.run([sys.executable, "-I", str(BOOT)], input=json.dumps(cfg), capture_output=True, text=True)

    OK = [sys.executable, "-c", "import sys; d=sys.stdin.read().split(); assert len(d)==2 and d[0]==d[1]"]

    def test_creates_private_unencrypted_room_and_returns_token_without_logging_it(self):
        r = self.run_boot(self.OK)
        self.assertEqual(r.returncode, 0, r.stderr)
        out = json.loads(r.stdout)
        self.assertEqual(out, {"token": "syt_FAKE_TOKEN_VALUE", "room_id": "!fakeroom:onemindos.dev", "mxid": "@gate:onemindos.dev"})
        self.assertNotIn("syt_FAKE", r.stderr)
        self.assertNotIn("PW-SECRET", r.stderr + r.stdout)
        room = next(s[2] for s in FakeSynapse.seen if s[1].endswith("/createRoom"))
        self.assertEqual(room["preset"], "private_chat")
        self.assertEqual(room["invite"], ["@zeus:onemindos.dev"])
        self.assertNotIn("initial_state", room)  # no m.room.encryption event
        login = next(s[2] for s in FakeSynapse.seen if s[1].endswith("/login"))
        self.assertEqual(login["password"], "PW-SECRET-VALUE-123")

    def test_default_register_command_takes_the_password_on_stdin_not_argv(self):
        src = BOOT.read_text()
        self.assertIn('"register_new_matrix_user", "-c", "/config/homeserver.yaml", "-u", user, "--no-admin", base', src)
        self.assertNotIn('"-p"', src)
        self.assertIn("input=pw", src)

    def test_missing_owner_stops_before_registering(self):
        FakeSynapse.owner_exists = False
        r = self.run_boot([sys.executable, "-c", "open('/nonexistent-dir/x','w')"])
        self.assertEqual(r.returncode, 2)
        self.assertFalse(any(s[1].endswith("/login") for s in FakeSynapse.seen))

    def test_existing_bot_user_is_reported_and_the_password_is_not(self):
        r = self.run_boot([sys.executable, "-c", "import sys; print('ERROR! User ID already taken.', file=sys.stderr); sys.exit(1)"])
        self.assertEqual(r.returncode, 3)
        self.assertNotIn("PW-SECRET", r.stderr)

    def test_encrypted_room_is_refused(self):
        FakeSynapse.encrypted = True
        r = self.run_boot(self.OK)
        self.assertEqual(r.returncode, 7)
        self.assertEqual(r.stdout, "")


if __name__ == "__main__":
    unittest.main()
