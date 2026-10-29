"""FEATURES D5 against the real policy service (omos PR #554), not a stub.

Runs only when the policy binary and its source tree are given (the service is not on main yet):

    POLICY_BIN=/path/to/policy  POLICY_SRC=/path/to/omos-checkout-of-#554  pytest test_policy_integration.py

The server uses three accounts like the A19 layout (omos PR #562): ONEMIND (policy, entity
stub, admin recorder), CLIENT (human `zeus`), AGENTS (`hermes` with the A19 permission list, and
`rogue`, an agent identity that is allowed to publish svc.policy.> so the policy service's own
agent_cannot_confirm check is reached). Passwords are random per run; nothing is committed.

Proves:
1. A Hermes proposal (through the plugin tool) waits for confirmation; no cmd.fleet.> message.
2. The plugin has no confirm path; Hermes' own credential cannot publish svc.policy.*.confirm or
   cmd.fleet.> (server permission violation); an agent identity that reaches confirm is refused
   agent_cannot_confirm (stop-line test 2).
3. Only after the human `zeus` confirms does cmd.fleet.drone.sim-x1.goto appear, and every
   cmd.fleet.> message has a human evt.so.action.approved (stop-line test 4 invariant).
4. "Follow that person": the proposal is denied forbidden_action and recorded (stop-line test 6).
"""
from __future__ import annotations

import asyncio
import json
import os
import secrets
import subprocess
import threading
import time
from pathlib import Path

import pytest

from conftest import FakeCtx, NatsServer, free_port, load_plugin, nats_server_bin

P = load_plugin()
POLICY_BIN = os.environ.get("POLICY_BIN")
POLICY_SRC = os.environ.get("POLICY_SRC")

pytestmark = pytest.mark.skipif(not (POLICY_BIN and POLICY_SRC and nats_server_bin()),
                                reason="POLICY_BIN, POLICY_SRC and nats-server are required")

USERS = ["policy", "entity", "admin", "zeus", "hermes", "rogue"]


def accounts_conf(pw):
    return f"""
accounts {{
  ONEMIND {{
    jetstream: enabled
    users: [
      {{ user: policy, password: "{pw['policy']}", permissions: {{
          publish: {{ allow: ["cmd.>", "evt.so.action.>", "svc.entity.get", "$JS.API.>", "_INBOX.>"] }}
          subscribe: {{ allow: ["svc.policy.>", "evt.drone.>", "evt.robot.>", "evt.vehicle.>", "evt.mission.>", "evt.entity.>", "_INBOX.>", "$KV.>"] }}
          allow_responses: true }} }}
      {{ user: entity, password: "{pw['entity']}", permissions: {{
          publish: {{ allow: ["_INBOX.>"] }}, subscribe: {{ allow: ["svc.entity.>"] }}, allow_responses: true }} }}
      {{ user: admin, password: "{pw['admin']}" }}
    ]
    exports: [ {{ service: "svc.policy.>" }}, {{ service: "svc.entity.>" }} ]
  }}
  CLIENT {{
    users: [
      {{ user: zeus, password: "{pw['zeus']}", permissions: {{
          publish: {{ allow: ["svc.policy.>", "svc.entity.>", "cmd.map.>"] }}, subscribe: {{ allow: ["_INBOX.>", "cmd.map.>"] }} }} }}
    ]
    imports: [ {{ service: {{ account: ONEMIND, subject: "svc.policy.>" }}, share: true }} ]
  }}
  AGENTS {{
    users: [
      {{ user: hermes, password: "{pw['hermes']}", permissions: {{
          publish: {{ allow: ["agents.*.hermes.>", "cmd.map.>", "svc.entity.>", "svc.layer.list", "svc.layer.meta",
                             "svc.policy.command.propose", "svc.policy.task.propose", "svc.policy.action.get"] }}
          subscribe: {{ allow: ["agents.prompt.hermes.>", "_INBOX.>"] }} }} }}
      {{ user: rogue, password: "{pw['rogue']}", permissions: {{
          publish: {{ allow: ["svc.policy.>"] }}, subscribe: {{ allow: ["_INBOX.>"] }} }} }}
    ]
    imports: [
      {{ service: {{ account: ONEMIND, subject: "svc.policy.command.propose" }}, share: true }}
      {{ service: {{ account: ONEMIND, subject: "svc.policy.task.propose" }}, share: true }}
      {{ service: {{ account: ONEMIND, subject: "svc.policy.action.get" }}, share: true }}
      {{ service: {{ account: ONEMIND, subject: "svc.policy.command.confirm" }}, share: true }}
      {{ service: {{ account: ONEMIND, subject: "svc.policy.task.confirm" }}, share: true }}
      {{ service: {{ account: ONEMIND, subject: "svc.entity.>" }}, share: true }}
    ]
  }}
}}
"""
# Note: AGENTS imports confirm only so that `rogue` can reach it; the A19 layout does not import
# it at all, and `hermes` is not allowed to publish it either way.


DIRECTORY = {"own_account": "ONEMIND", "agent_accounts": ["AGENTS"], "principals": [
    {"account": "CLIENT", "user": "zeus", "kind": "user", "id": "user:zeus"},
    {"account": "AGENTS", "user": "hermes", "kind": "agent", "id": "agent:hermes/zeus/legacy"},
    {"account": "AGENTS", "user": "rogue", "kind": "user", "id": "user:rogue"},  # mislabelled on purpose
]}


def tuples():
    model = json.loads((Path(POLICY_SRC) / "platform/security/openfga/onemind-actions.json").read_text())
    del model
    names = ["asset.goto", "asset.hold", "asset.rtl", "asset.land", "asset.orbit", "asset.observe", "asset.estop",
             "asset.disarm", "task.assign", "task.abort", "entity.star"]
    ts = [{"user": "user:zeus", "relation": "operator", "object": "asset:sim-x1"},
          {"user": "user:rogue", "relation": "operator", "object": "asset:sim-x1"}]
    for n in names:
        obj = "action_type:" + n
        for u in ("user:zeus", "user:rogue"):
            ts += [{"user": u, "relation": r, "object": obj} for r in ("proposer", "approver")]
        ts.append({"user": "agent:hermes/zeus/legacy", "relation": "proposer", "object": obj})
    return ts


ENTITIES = {
    "sim-x1": {"entity_id": "sim-x1", "entity_type": "drone", "kind": "multirotor",
               "location": {"lat": 30.2700, "lon": -97.7400, "hae_m": 180.0},
               "health": {"connection": "online"}, "indicators": {"simulated": True},
               "capabilities": ["onemind.task.goto.v1", "onemind.task.orbit.v1", "onemind.task.rtl.v1"]},
    "trk-p1": {"entity_id": "trk-p1", "entity_type": "track", "kind": "person",
               "location": {"lat": 30.2805, "lon": -97.7430}, "indicators": {"simulated": True}},
}


class Bg:
    """Async helpers on their own loop: entity stub, recorders, raw clients."""

    def __init__(self):
        self.loop = asyncio.new_event_loop()
        threading.Thread(target=self.loop.run_forever, daemon=True).start()
        self.conns = []

    def run(self, coro, t=15):
        return asyncio.run_coroutine_threadsafe(coro, self.loop).result(t)

    async def connect(self, url, errors=None):
        import nats

        async def err_cb(e):
            if errors is not None:
                errors.append(str(e))

        nc = await nats.connect(url, error_cb=err_cb, max_reconnect_attempts=0)
        self.conns.append(nc)
        return nc

    def close(self):
        for nc in self.conns:
            try:
                self.run(nc.close(), 5)
            except Exception:
                pass
        self.loop.call_soon_threadsafe(self.loop.stop)


@pytest.fixture(scope="module")
def stack(tmp_path_factory):
    base = tmp_path_factory.mktemp("policy-int")
    pw = {u: secrets.token_hex(16) for u in USERS}
    srv = NatsServer(base / "nats", accounts_conf(pw))
    bg = Bg()
    cfg = base / "cfg"
    cfg.mkdir()
    (cfg / "principals.json").write_text(json.dumps(DIRECTORY))
    (cfg / "tuples.json").write_text(json.dumps(tuples()))
    rec = {"cmd": [], "audit": []}

    async def setup():
        import nats  # noqa: F401

        admin = await bg.connect(srv.url("admin", pw["admin"]))
        js = admin.jetstream()
        from nats.js.api import StreamConfig, KeyValueConfig, DiscardPolicy

        await js.add_stream(StreamConfig(name="CMD", subjects=["cmd.>"], discard=DiscardPolicy.NEW, duplicate_window=120))
        await js.add_stream(StreamConfig(name="EVT", subjects=["evt.so.>", "evt.drone.>", "evt.robot.>"], duplicate_window=120))
        await js.add_stream(StreamConfig(name="MISSION", subjects=["evt.mission.>"]))
        kv = await js.create_key_value(KeyValueConfig(bucket="geofences", history=10))
        await kv.put("zone-home", json.dumps({"entity_id": "zone-home", "entity_type": "zone", "live": True, "geo": {
            "shape": "circle", "kind": "keep_in", "radius_m": 20000,
            "coords": {"type": "Point", "coordinates": [-97.7431, 30.2672]}}}).encode())

        async def on_cmd(m):
            rec["cmd"].append((m.subject, json.loads(m.data)))

        async def on_audit(m):
            rec["audit"].append((m.subject, json.loads(m.data)))

        await admin.subscribe("cmd.>", cb=on_cmd)
        await admin.subscribe("evt.so.action.>", cb=on_audit)
        ent = await bg.connect(srv.url("entity", pw["entity"]))

        async def get(m):
            q = json.loads(m.data)
            e = ENTITIES.get(q.get("entity_id"))
            await m.respond(json.dumps({"entity": e} if e else {"error": {"code": "not_found"}}).encode())

        await ent.subscribe("svc.entity.get", cb=get)
        await admin.flush()
        await ent.flush()

    bg.run(setup())
    env = dict(os.environ, NATS_URL=srv.url(), NATS_USER="policy", NATS_PASSWORD=pw["policy"], POLICY_AUTHZ="local",
               POLICY_MODEL_FILE=str(Path(POLICY_SRC) / "platform/security/openfga/onemind-actions.json"),
               POLICY_TUPLES_FILE=str(cfg / "tuples.json"), POLICY_SCHEMA_DIR=str(Path(POLICY_SRC) / "services/schemas"),
               POLICY_PRINCIPALS_FILE=str(cfg / "principals.json"), HEALTH_ADDR=f"127.0.0.1:{free_port()}")
    log = open(base / "policy.log", "w")
    proc = subprocess.Popen([POLICY_BIN], env=env, stdout=log, stderr=subprocess.STDOUT)
    deadline = time.time() + 15
    while time.time() < deadline and "policy service ready" not in (base / "policy.log").read_text():
        if proc.poll() is not None:
            raise RuntimeError("policy exited: " + (base / "policy.log").read_text())
        time.sleep(0.1)
    hermes_env = {"ONEMIND_BUS_URL": srv.url(), "ONEMIND_BUS_USER": "hermes", "ONEMIND_BUS_PASSWORD": pw["hermes"],
                  "ONEMIND_BUS_ON_BEHALF_OF": "user:zeus", "ONEMIND_BUS_TIMEOUT_S": "5", "ONEMIND_BUS_MAP_ACK_WINDOW_S": "0.5"}
    ctx = FakeCtx()
    tools = P.register(ctx, env=hermes_env)
    s = type("Stack", (), {})()
    s.srv, s.pw, s.bg, s.rec, s.ctx, s.base, s.policy_log = srv, pw, bg, rec, ctx, base, base / "policy.log"
    yield s
    tools.client.close()
    proc.terminate()
    proc.wait(10)
    log.close()
    bg.close()
    srv.stop()
    out = os.environ.get("ONEMIND_BUS_EVIDENCE_DIR")
    if out:
        Path(out).mkdir(parents=True, exist_ok=True)
        with open(Path(out) / "d5-policy-integration-bus.jsonl", "w") as f:
            for kind in ("audit", "cmd"):
                for subj, body in rec[kind]:
                    f.write(json.dumps({"stream": kind, "subject": subj, "body": body}) + "\n")
        (Path(out) / "d5-policy-service.log").write_text((base / "policy.log").read_text())
        (Path(out) / "d5-nats-server.log").write_text(srv.log_text())


def _request(stack, user, subject, body, timeout=3):
    errors = []

    async def go():
        nc = await stack.bg.connect(stack.srv.url(user, stack.pw[user]), errors)
        try:
            m = await nc.request(subject, json.dumps(body).encode(), timeout=timeout)
            return json.loads(m.data or b"{}")
        except Exception as e:
            return {"_exception": type(e).__name__}
        finally:
            await asyncio.sleep(0.2)
            await nc.close()

    r = stack.bg.run(go())
    return r, errors


def _audit(stack, event, action_id):
    return [b for s, b in stack.rec["audit"] if s == f"evt.so.action.{event}" and b["data"]["action_id"] == action_id]


GOTO = {"action_type": "asset.goto", "target": {"entity_type": "drone", "entity_id": "sim-x1"},
        "params": {"lat": 30.2750, "lon": -97.7420, "hae_m": 200}, "reason": "look at the barn roof"}


def test_1_hermes_proposal_waits_for_a_human(stack):
    r = stack.ctx.dispatch("onemind_command_propose", GOTO, task_id="sess-d5", user_task="send sim-x1 over the barn")
    assert r["status"] == "awaiting_confirmation", json.dumps(r)
    assert r["approvals_needed"] == 1 and r["level"] == "A3"
    stack.action = r
    time.sleep(1.5)
    assert [c for c in stack.rec["cmd"] if c[0].startswith("cmd.fleet.")] == []
    prop = _audit(stack, "proposed", r["action_id"])
    assert prop and prop[0]["data"]["requested_by"] == {"kind": "agent", "id": "agent:hermes/zeus/legacy",
                                                         "on_behalf_of": "user:zeus"}
    assert prop[0]["data"]["actor_verified"] is True and prop[0]["data"]["agent_session"] == "sess-d5"
    st = stack.ctx.dispatch("onemind_action_status", {"action_id": r["action_id"]})
    assert st["action"]["status"] == "awaiting_confirmation"


def test_2_hermes_cannot_confirm(stack):
    a = stack.action
    # (a) the plugin has no confirm tool and its client refuses the subject
    assert not [t for t in stack.ctx.tools if "confirm" in t]
    with pytest.raises(P.bus.SubjectNotAllowed):
        P.bus.check_request_subject("svc.policy.command.confirm")
    body = {"v": 1, "action_id": a["action_id"], "params_hash": a["params_hash"], "decision": "approve"}
    # (b) Hermes' own NATS credential: the server refuses the publish
    r, errs = _request(stack, "hermes", "svc.policy.command.confirm", body, timeout=1.5)
    assert "_exception" in r and any("permissions violation" in e.lower() for e in errs), (r, errs)
    r, errs = _request(stack, "hermes", "cmd.fleet.drone.sim-x1.goto", {"lat": 1}, timeout=1)
    assert any("permissions violation" in e.lower() for e in errs)
    # (c) an agent identity that can reach confirm is refused by the policy service
    r, _ = _request(stack, "rogue", "svc.policy.command.confirm", body)
    assert r["denial"]["code"] == "agent_cannot_confirm", r
    time.sleep(0.5)
    assert [c for c in stack.rec["cmd"] if c[0].startswith("cmd.fleet.")] == []
    refused = [b for s, b in stack.rec["audit"] if b["data"].get("code") == "agent_cannot_confirm"]
    assert refused, "the refused confirmation is audited"


def test_3_only_a_human_confirm_dispatches(stack):
    a = stack.action
    body = {"v": 1, "action_id": a["action_id"], "params_hash": a["params_hash"], "decision": "approve"}
    r, _ = _request(stack, "zeus", "svc.policy.command.confirm", body)
    assert r["status"] in ("approved", "dispatched") and r.get("command_id"), r
    deadline = time.time() + 5
    while time.time() < deadline and not [c for c in stack.rec["cmd"] if c[0] == "cmd.fleet.drone.sim-x1.goto"]:
        time.sleep(0.1)
    cmds = [c for c in stack.rec["cmd"] if c[0] == "cmd.fleet.drone.sim-x1.goto"]
    assert len(cmds) == 1
    approved = _audit(stack, "approved", a["action_id"])
    assert approved and approved[0]["data"]["approvals"][0]["by"] == "user:zeus"


def test_4_follow_that_person_is_denied_and_recorded(stack):
    r = stack.ctx.dispatch("onemind_command_propose", {
        "action_type": "asset.orbit", "target": {"entity_type": "drone", "entity_id": "sim-x1"},
        "params": {"center": {"entity_id": "trk-p1"}, "radius_m": 30, "hae_m": 60, "duration_s": 600},
        "reason": "operator asked to follow that person"}, task_id="sess-d5-person")
    assert r["denial"]["code"] == "forbidden_action", r
    time.sleep(0.5)
    denied = _audit(stack, "denied", r["action_id"])
    assert denied and denied[0]["data"]["code"] == "forbidden_action"


def test_5_invariant_no_fleet_command_without_a_human_approval(stack):
    time.sleep(0.5)
    approved_ids = {b["data"]["action_id"] for s, b in stack.rec["audit"] if s == "evt.so.action.approved"
                    and all(ap["by"].startswith("user:") for ap in b["data"]["approvals"])}
    fleet = [b for s, b in stack.rec["cmd"] if s.startswith("cmd.fleet.")]
    assert fleet, "the run did dispatch one command"
    for c in fleet:
        aid = c.get("action_id")
        assert aid in approved_ids, c
    log = stack.srv.log_text()
    assert "\"AGENTS/user:hermes\" - Publish Violation - Subject \"cmd.fleet.drone.sim-x1.goto\"" in log
    assert "\"AGENTS/user:hermes\" - Publish Violation - Subject \"svc.policy.command.confirm\"" in log
