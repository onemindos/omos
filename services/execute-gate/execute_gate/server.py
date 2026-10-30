"""HTTP front (AuthZEN evaluation endpoint) and process wiring. Stdlib http.server only."""

from __future__ import annotations

import hmac
import json
import os
import signal
import socketserver
import sys
import threading
from http.server import BaseHTTPRequestHandler, ThreadingHTTPServer

from .audit import Audit
from .gate import Gate
from .matrix import MatrixClient, sync_loop
from .policy import load_policy
from .store import MemoryStore, NatsKvStore

MAX_BODY = 1 << 20


class Server(ThreadingHTTPServer):
    daemon_threads = True

    def server_bind(self):  # skip http.server's reverse-DNS getfqdn(), which can stall for seconds
        socketserver.TCPServer.server_bind(self)
        self.server_name, self.server_port = str(self.server_address[0]), self.server_address[1]


def decide(gate: Gate, body: dict) -> dict:
    """AuthZEN evaluation: subject/action/resource/context -> {"decision": bool, "context": {...}}."""
    subject = (body.get("subject") or {}).get("id")
    action = (body.get("action") or {}).get("name")
    res = body.get("resource") or {}
    ctx = body.get("context") or {}
    if action != "tools/call":
        gate.audit.emit("denied", reason="unsupported_action", action=str(action))
        return {"decision": False, "context": {"reason": "unsupported_action"}}
    resource = res.get("id")
    args = ctx.get("arguments")
    if args is None:
        args = {}
    ok, reason, rid = gate.evaluate(subject, resource, args)
    out = {"reason": reason}
    if rid:
        out["request_id"] = rid
    return {"decision": bool(ok), "context": out}


def make_handler(gate: Gate, auth_token: str | None, ready):
    class H(BaseHTTPRequestHandler):
        protocol_version = "HTTP/1.1"

        def log_message(self, *a):  # audit lines are the log
            pass

        def _send(self, code, obj):
            raw = json.dumps(obj).encode()
            self.send_response(code)
            self.send_header("Content-Type", "application/json")
            self.send_header("Content-Length", str(len(raw)))
            self.end_headers()
            self.wfile.write(raw)

        def do_GET(self):
            if self.path == "/healthz":
                return self._send(200, {"ok": True})
            if self.path == "/readyz":
                good = ready()
                return self._send(200 if good else 503, {"ready": good})
            self._send(404, {"error": "not found"})

        def do_POST(self):
            if self.path != "/access/v1/evaluation":
                return self._send(404, {"decision": False})
            try:
                n = int(self.headers.get("Content-Length", "0"))
            except ValueError:
                n = 0
            if n <= 0 or n > MAX_BODY:
                self.close_connection = True
                return self._send(400, {"decision": False})
            raw = self.rfile.read(n)
            if auth_token:
                got = self.headers.get("Authorization", "")
                if not hmac.compare_digest(got.encode(), ("Bearer " + auth_token).encode()):
                    return self._send(401, {"decision": False})
            try:
                body = json.loads(raw)
                if not isinstance(body, dict):
                    raise ValueError
            except (ValueError, TypeError):
                return self._send(400, {"decision": False})
            try:
                self._send(200, decide(gate, body))
            except Exception:
                self._send(200, {"decision": False, "context": {"reason": "error"}})

    return H


def _env(name, default=None, required=False):
    v = os.environ.get(name, default)
    if required and not v:
        sys.exit(f"missing required environment variable {name}")
    return v


def main() -> None:
    owner = _env("OWNER_MXID", "@zeus:onemindos.dev")
    room = _env("APPROVALS_ROOM_ID", required=True)
    token = _env("MATRIX_ACCESS_TOKEN", required=True)
    hs = _env("MATRIX_HOMESERVER", "http://synapse.matrix.svc.cluster.local:8008")
    subjects = frozenset(s for s in _env("ALLOWED_SUBJECTS", "hermes").split(",") if s)
    policy = load_policy(_env("POLICY_FILE"), subjects)
    client = MatrixClient(hs, token, room)
    bot = _env("BOT_MXID") or client.whoami()
    if bot == owner:
        sys.exit("bot MXID must differ from OWNER_MXID")

    if _env("STORE", "nats") == "memory":
        store, publish = MemoryStore(), None
    else:
        store = NatsKvStore(_env("NATS_URL", "nats://nats.nats.svc.cluster.local:4222"),
                            _env("NATS_USER"), _env("NATS_PASSWORD"), _env("KV_BUCKET", "exec-approvals"))
        publish = store.publish
    gate = Gate(store, policy, owner, room, bot, client, Audit(publish),
                ttl=int(_env("APPROVAL_TTL_SECONDS", "600")))

    stop = threading.Event()
    th = threading.Thread(target=sync_loop, args=(client, gate, stop), name="matrix-sync", daemon=True)
    th.start()
    httpd = Server(("0.0.0.0", int(_env("PORT", "8080"))),
                                make_handler(gate, _env("GATE_AUTH_TOKEN"), lambda: store.healthy() and th.is_alive()))

    def shutdown(*_):
        stop.set()
        threading.Thread(target=httpd.shutdown, daemon=True).start()

    signal.signal(signal.SIGTERM, shutdown)
    signal.signal(signal.SIGINT, shutdown)
    gate.audit.emit("started", owner=owner, bot=bot, ttl=gate.ttl)
    httpd.serve_forever()


if __name__ == "__main__":
    main()
