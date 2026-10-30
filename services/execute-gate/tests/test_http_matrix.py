import io
import json
import threading
import urllib.request

from execute_gate.audit import Audit
from execute_gate.gate import Gate
from execute_gate.matrix import events_from_sync
from execute_gate.policy import Policy
from execute_gate.server import Server, make_handler
from execute_gate.store import MemoryStore

ROOM = "!a:onemindos.dev"


class N:
    def send(self, text):
        pass


def test_http_endpoint_and_bearer():
    gate = Gate(MemoryStore(), Policy(), "@zeus:onemindos.dev", ROOM, "@gate:onemindos.dev", N(), Audit(None, io.StringIO()))
    srv = Server(("127.0.0.1", 0), make_handler(gate, "s3cret", lambda: True))
    threading.Thread(target=srv.serve_forever, daemon=True).start()
    base = f"http://127.0.0.1:{srv.server_port}"
    body = json.dumps({"subject": {"id": "hermes"}, "action": {"name": "tools/call"},
                       "resource": {"id": "k/x"}, "context": {"arguments": {"a": 1}}}).encode()

    def post(auth):
        req = urllib.request.Request(base + "/access/v1/evaluation", data=body, method="POST")
        if auth:
            req.add_header("Authorization", auth)
        try:
            with urllib.request.urlopen(req) as r:
                return r.status, json.loads(r.read())
        except urllib.error.HTTPError as e:
            return e.code, json.loads(e.read())

    try:
        assert post(None)[0] == 401
        assert post("Bearer wrong")[0] == 401
        code, data = post("Bearer s3cret")
        assert code == 200 and data["decision"] is False
        with urllib.request.urlopen(base + "/readyz") as r:
            assert r.status == 200
    finally:
        srv.shutdown()


def test_events_from_sync_extracts_text_and_edits():
    resp = {"rooms": {"join": {ROOM: {"timeline": {"events": [
        {"type": "m.room.message", "event_id": "$1", "sender": "@zeus:onemindos.dev",
         "content": {"msgtype": "m.text", "body": "execute AB12"}},
        {"type": "m.room.message", "event_id": "$2", "sender": "@zeus:onemindos.dev",
         "content": {"msgtype": "m.text", "body": "* execute AB12", "m.new_content": {"body": "execute AB12"},
                     "m.relates_to": {"rel_type": "m.replace", "event_id": "$1"}}},
        {"type": "m.room.member", "event_id": "$3", "sender": "@x:y", "content": {}},
    ]}}}}}
    got = list(events_from_sync(resp, ROOM))
    assert [g[0] for g in got] == ["$1", "$2"]
    assert got[0][4] is False and got[1][4] is True
    assert list(events_from_sync(resp, "!other:x")) == []
