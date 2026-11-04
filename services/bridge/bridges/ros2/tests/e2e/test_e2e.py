"""End-to-end test (FEATURES H3 acceptance, asset-command.md section 9 style).

Starts a disposable nats-server with per-identity accounts (ONEMIND, CLIENT with a
``share: true`` import of svc.policy.>, AGENTS), the policy service from omos PR #554
(``POLICY_AUTHZ=local``), a stub entity service, the scripted Zenoh robot and this bridge as
separate processes, then drives the whole action flow as a human on the CLIENT account:

  propose asset.goto -> confirm -> policy signs cmd.fleet.robot.sim-tb1.goto -> bridge
  verifies -> Zenoh NavigateToPose goal -> robot drives -> tel.fleet.robot.sim-tb1.position
  arrives at the target -> evt.robot.sim-tb1.command executed -> evt.so.action.executed.

Plus the refusals: agent cannot confirm, a client cannot publish cmd.fleet.>, replayed and
unsigned commands, a target the policy allows but the bridge's own fence refuses, a
policy-side geofence denial, a non-sim target, and the e-stop latch.

Runs only when NATS_SERVER (nats-server binary), POLICY_BIN (policy binary built from
PR #554) and POLICY_SRC (a checkout containing services/schemas and
platform/security/openfga) are set; tests/e2e/run.sh sets them. Passwords and the policy
signing key are generated per run in a temp dir; nothing is committed.
"""

from __future__ import annotations

import asyncio
import base64
import json
import os
import secrets
import shutil
import socket
import subprocess
import sys
import tempfile
import time
import uuid

import pytest

NATS_SERVER = os.environ.get("NATS_SERVER", "")
POLICY_BIN = os.environ.get("POLICY_BIN", "")
POLICY_SRC = os.environ.get("POLICY_SRC", "")
pytestmark = pytest.mark.skipif(not (NATS_SERVER and POLICY_BIN and POLICY_SRC),
                                reason="set NATS_SERVER, POLICY_BIN and POLICY_SRC (tests/e2e/run.sh)")

HERE = os.path.dirname(os.path.abspath(__file__))
PKG = os.path.dirname(os.path.dirname(HERE))
sys.path.insert(0, PKG)

ORIGIN = (30.2672, -97.7431, 150.0)
ASSET = "sim-tb1"
# Owner run against a real ROS 2 robot or Gazebo (OWNER-STEPS.md): point the bridge at an
# external zenoh-bridge-ros2dds instead of starting the scripted robot.
EXT_ZENOH = os.environ.get("E2E_ZENOH_CONNECT", "")
NAMESPACE = os.environ.get("E2E_ROS_NAMESPACE", "sim_tb1")
POSE_TOPIC = os.environ.get("E2E_POSE_TOPIC", "")
GOAL = (float(os.environ.get("E2E_GOAL_X", "1.2")), float(os.environ.get("E2E_GOAL_Y", "0.6")))
ARRIVE_S = float(os.environ.get("E2E_ARRIVE_TIMEOUT", "30"))


def free_port() -> int:
    s = socket.socket()
    s.bind(("127.0.0.1", 0))
    p = s.getsockname()[1]
    s.close()
    return p


def nats_conf(port: int, store: str, pw: dict) -> str:
    return f"""
listen: "127.0.0.1:{port}"
jetstream {{ store_dir: "{store}", max_mem: 64MB, max_file: 256MB }}
accounts {{
  ONEMIND {{
    jetstream: enabled
    users: [
      {{ user: policy, password: "{pw['policy']}", permissions: {{
          publish: {{ allow: ["cmd.>", "evt.so.action.>", "svc.entity.get", "$JS.API.>", "_INBOX.>"] }}
          subscribe: {{ allow: ["svc.policy.>", "evt.drone.>", "evt.robot.>", "evt.vehicle.>", "evt.mission.>", "evt.entity.>", "_INBOX.>"] }}
          allow_responses: true }} }}
      {{ user: entity, password: "{pw['entity']}", permissions: {{
          publish: {{ allow: ["_INBOX.>"] }}, subscribe: {{ allow: ["svc.entity.>", "tel.fleet.>"] }}, allow_responses: true }} }}
      {{ user: bridge, password: "{pw['bridge']}", permissions: {{
          publish: {{ allow: ["tel.fleet.robot.{ASSET}.>", "evt.robot.{ASSET}.>"] }}
          subscribe: {{ allow: ["cmd.fleet.robot.{ASSET}.>", "_INBOX.>"] }} }} }}
      {{ user: admin, password: "{pw['admin']}" }}
    ]
    exports: [ {{ service: "svc.policy.>" }} ]
  }}
  CLIENT {{
    users: [
      {{ user: zeus, password: "{pw['zeus']}", permissions: {{
          publish: {{ allow: ["svc.policy.>", "cmd.map.>"] }}, subscribe: {{ allow: ["_INBOX.>"] }} }} }}
    ]
    imports: [ {{ service: {{ account: ONEMIND, subject: "svc.policy.>" }}, share: true }} ]
  }}
  AGENTS {{
    users: [
      {{ user: hermes, password: "{pw['hermes']}", permissions: {{
          publish: {{ allow: ["svc.policy.command.propose", "svc.policy.task.propose", "svc.policy.action.get", "svc.policy.command.confirm"] }}
          subscribe: {{ allow: ["_INBOX.>"] }} }} }}
    ]
    imports: [ {{ service: {{ account: ONEMIND, subject: "svc.policy.>" }}, share: true }} ]
  }}
}}
"""


def zone(entity_id: str, kind: str, radius_m: float) -> bytes:
    return json.dumps({"entity_id": entity_id, "live": True,
                       "geo": {"shape": "circle", "kind": kind, "radius_m": radius_m,
                               "coords": {"type": "Point", "coordinates": [ORIGIN[1], ORIGIN[0]]}},
                       "properties": {"applies_to": [ASSET]}}).encode()


class Proc:
    def __init__(self, name, args, env, logdir):
        self.name = name
        self.log = open(os.path.join(logdir, f"{name}.log"), "wb")
        self.p = subprocess.Popen(args, env=env, stdout=self.log, stderr=subprocess.STDOUT)

    def stop(self):
        if self.p.poll() is None:
            self.p.terminate()
            try:
                self.p.wait(5)
            except subprocess.TimeoutExpired:
                self.p.kill()
        self.log.close()


@pytest.fixture
def stack():
    from cryptography.hazmat.primitives import serialization
    from cryptography.hazmat.primitives.asymmetric.ed25519 import Ed25519PrivateKey

    evidence = os.environ.get("E2E_EVIDENCE_DIR", "")
    work = tempfile.mkdtemp(prefix="h3-e2e-")
    os.chmod(work, 0o700)
    logdir = evidence or work
    os.makedirs(logdir, exist_ok=True)
    pw = {k: secrets.token_urlsafe(18) for k in ("policy", "entity", "bridge", "admin", "zeus", "hermes")}
    nport, zport, hport = free_port(), free_port(), free_port()
    with open(os.path.join(work, "nats.conf"), "w") as f:
        f.write(nats_conf(nport, os.path.join(work, "js"), pw))
    key = Ed25519PrivateKey.generate()
    keyfile = os.path.join(work, "policy-signing.pem")
    with open(keyfile, "wb") as f:
        f.write(key.private_bytes(serialization.Encoding.PEM, serialization.PrivateFormat.PKCS8,
                                  serialization.NoEncryption()))
    pub_b64 = base64.b64encode(key.public_key().public_bytes_raw()).decode()
    with open(os.path.join(work, "principals.json"), "w") as f:
        json.dump({"own_account": "ONEMIND", "agent_accounts": ["AGENTS"], "principals": [
            {"account": "CLIENT", "user": "zeus", "kind": "user", "id": "user:zeus"},
            {"account": "AGENTS", "user": "hermes", "kind": "agent", "id": "agent:hermes/zeus/legacy"}]}, f)
    tuples = [{"user": "user:zeus", "relation": "operator", "object": f"asset:{ASSET}"},
              {"user": "user:zeus", "relation": "operator", "object": "asset:tb1"}]
    for verb in ("goto", "hold", "rtl", "estop", "estop_release"):
        obj = f"action_type:asset.{verb}"
        tuples += [{"user": "user:zeus", "relation": r, "object": obj} for r in ("proposer", "approver", "second_approver")]
        tuples.append({"user": "agent:hermes/zeus/legacy", "relation": "proposer", "object": obj})
    with open(os.path.join(work, "tuples.json"), "w") as f:
        json.dump(tuples, f)
    fence_path = os.path.join(work, "bridge-fence.json")
    with open(fence_path, "w") as f:  # the bridge's own fence is tighter (30 m) than the policy's zone (200 m)
        json.dump({"keep_in": [{"id": "yard", "circle": {"lat": ORIGIN[0], "lon": ORIGIN[1], "radius_m": 30}}]}, f)

    procs = []
    procs.append(Proc("nats-server", [NATS_SERVER, "-c", os.path.join(work, "nats.conf")], dict(os.environ), logdir))
    url = f"nats://127.0.0.1:{nport}"
    for _ in range(100):
        try:
            socket.create_connection(("127.0.0.1", nport), 0.2).close()
            break
        except OSError:
            time.sleep(0.05)

    async def setup_js():
        import nats

        nc = await nats.connect(url, user="admin", password=pw["admin"])
        js = nc.jetstream()
        from nats.js.api import StreamConfig

        await js.add_stream(StreamConfig(name="CMD", subjects=["cmd.>"], duplicate_window=120))
        await js.add_stream(StreamConfig(name="EVT", subjects=["evt.so.>", "evt.drone.>", "evt.robot.>"], duplicate_window=120))
        await js.add_stream(StreamConfig(name="TEL", subjects=["tel.fleet.>"], max_age=3600))
        kv = await js.create_key_value(bucket="geofences", history=10)
        await kv.put("yard", zone("zone-yard", "keep_in", 200.0))
        await nc.close()

    asyncio.run(setup_js())
    py = sys.executable
    base_env = {k: v for k, v in os.environ.items() if not k.startswith("NATS_")}
    procs.append(Proc("policy", [POLICY_BIN], dict(base_env, NATS_URL=url, NATS_USER="policy", NATS_PASSWORD=pw["policy"],
                                                   POLICY_AUTHZ="local",
                                                   POLICY_MODEL_FILE=os.path.join(POLICY_SRC, "platform/security/openfga/onemind-actions.json"),
                                                   POLICY_TUPLES_FILE=os.path.join(work, "tuples.json"),
                                                   POLICY_SCHEMA_DIR=os.path.join(POLICY_SRC, "services/schemas"),
                                                   POLICY_PRINCIPALS_FILE=os.path.join(work, "principals.json"),
                                                   POLICY_SIGNING_KEY_FILE=keyfile, POLICY_SIGNING_KID="policy-1",
                                                   HEALTH_ADDR=f"127.0.0.1:{hport}"), logdir))
    if not EXT_ZENOH:
        procs.append(Proc("robot-sim", [py, "-u", os.path.join(PKG, "sim", "turtlebot_sim.py"), "--mode", "router",
                                        "--listen", f"tcp/127.0.0.1:{zport}", "--namespace", NAMESPACE,
                                        "--log", os.path.join(logdir, "robot-sim.jsonl")], base_env, logdir))
    procs.append(Proc("bridge-ros2", [py, "-u", "-m", "ros2bridge"], dict(
        base_env, PYTHONPATH=PKG, NATS_URL=url, NATS_USER="bridge", NATS_PASSWORD=pw["bridge"], ASSET_ID=ASSET,
        TRANSPORT="zenoh", ZENOH_CONNECT=EXT_ZENOH or f"tcp/127.0.0.1:{zport}", ROS_NAMESPACE=NAMESPACE,
        POSE_TOPIC=POSE_TOPIC,
        MAP_ORIGIN_LAT=str(ORIGIN[0]), MAP_ORIGIN_LON=str(ORIGIN[1]), MAP_ORIGIN_HAE=str(ORIGIN[2]),
        POLICY_PUBKEYS=f"policy-1={pub_b64}", GEOFENCE_FILE=fence_path, POSITION_INTERVAL="200ms"), logdir))
    yield {"url": url, "pw": pw, "key": key, "work": work, "logdir": logdir, "procs": procs}
    for p in reversed(procs):
        p.stop()
    shutil.rmtree(work, ignore_errors=True)


def test_h3_end_to_end(stack):
    asyncio.run(_flow(stack))


async def _flow(st):
    import nats
    from nats.errors import TimeoutError as NatsTimeout

    from ros2bridge import command as C
    from ros2bridge.geo import LocalFrame

    url, pw = st["url"], st["pw"]
    frame = LocalFrame(*ORIGIN, 90.0)
    transcript = open(os.path.join(st["logdir"], "transcript.jsonl"), "w")
    seen: list = []
    perm_errors: list = []

    admin = await nats.connect(url, user="admin", password=pw["admin"])

    async def rec(m):
        try:
            data = json.loads(m.data)
        except ValueError:
            data = m.data.decode(errors="replace")
        seen.append((time.time(), m.subject, data))
        if not m.subject.startswith("tel.") or len([1 for _, s, _ in seen if s == m.subject]) % 10 == 1:
            transcript.write(json.dumps({"t": round(time.time(), 3), "subject": m.subject, "data": data}) + "\n")
            transcript.flush()

    await admin.subscribe(">", cb=rec)

    # Stub entity service: answers svc.entity.get from the bridge's own telemetry.
    entity = await nats.connect(url, user="entity", password=pw["entity"])
    latest: dict = {}

    async def on_tel(m):
        latest[m.subject.split(".")[3]] = json.loads(m.data)

    async def on_get(m):
        req = json.loads(m.data)
        t = latest.get(req.get("entity_id"))
        if t is None and req.get("entity_id") != "tb1":
            await m.respond(json.dumps({"error": {"code": "not_found"}}).encode())
            return
        t = t or latest.get(ASSET)
        await m.respond(json.dumps({"entity": {
            "entity_id": req["entity_id"], "entity_type": "robot", "kind": "ugv",
            "location": {"lat": t["lat"], "lon": t["lon"], "hae_m": t["hae"]},
            "health": {"connection": "online"},
            "capabilities": ["onemind.task.goto.v1", "onemind.task.hold.v1", "onemind.task.rtl.v1"],
            "indicators": {"simulated": req["entity_id"].startswith("sim-")}}}).encode())

    await entity.subscribe(f"tel.fleet.robot.{ASSET}.position", cb=on_tel)
    await entity.subscribe("svc.entity.get", cb=on_get)

    async def zeus_err(e):
        perm_errors.append(str(e))

    zeus = await nats.connect(url, user="zeus", password=pw["zeus"], error_cb=zeus_err)
    hermes = await nats.connect(url, user="hermes", password=pw["hermes"])

    async def wait_for(pred, timeout, what):
        end = time.time() + timeout
        while time.time() < end:
            r = pred()
            if r:
                return r
            await asyncio.sleep(0.1)
        raise AssertionError(f"timed out waiting for {what}")

    def on(subject, since=0.0):
        return [d for t, s, d in seen if s == subject and t >= since]

    def outcome_for(cmd_id):
        return [d["data"] for d in on(f"evt.robot.{ASSET}.command") if d["data"]["command_id"] == cmd_id]

    async def propose(conn, params, verb="goto", entity_id=ASSET, kind="user", who="user:zeus"):
        req = {"v": 1, "request_id": str(uuid.uuid4()), "action_type": f"asset.{verb}", "action_version": 1,
               "target": {"entity_type": "robot", "entity_id": entity_id}, "params": params,
               "requested_by": {"kind": kind, "id": who, **({"on_behalf_of": "user:zeus"} if kind == "agent" else {})},
               "reason": "H3 e2e"}
        r = await conn.request("svc.policy.command.propose", json.dumps(req).encode(), timeout=5)
        return json.loads(r.data)

    async def confirm(conn, p):
        r = await conn.request("svc.policy.command.confirm", json.dumps(
            {"v": 1, "action_id": p["action_id"], "params_hash": p["params_hash"], "decision": "approve"}).encode(), timeout=5)
        return json.loads(r.data)

    def goto_params(x, y):
        lat, lon, _ = frame.to_wgs84(x, y)
        return {"lat": round(lat, 8), "lon": round(lon, 8), "hae_m": ORIGIN[2]}

    report = {}
    # 0. robot telemetry reaches the bus
    await wait_for(lambda: on(f"tel.fleet.robot.{ASSET}.position") and latest.get(ASSET), 20, "robot telemetry")
    await wait_for(lambda: on(f"tel.fleet.robot.{ASSET}.status"), 5, "robot status")

    # 1. the acceptance path: propose, human confirm, robot moves and arrives
    p = await propose(zeus, goto_params(*GOAL))
    assert p["status"] == "awaiting_confirmation", p
    c = await confirm(zeus, p)
    assert c["command_id"], c
    out = await wait_for(lambda: [o for o in outcome_for(c["command_id"]) if o["outcome"] == "executed"], 10, "executed")
    ack = [o for o in outcome_for(c["command_id"]) if o["outcome"] == "acknowledged"]
    assert ack and ack[0]["sig_verified"] is True, ack
    tx, ty = GOAL
    tlat, tlon, _ = frame.to_wgs84(tx, ty)
    from ros2bridge.fence import distance

    arrived = await wait_for(lambda: [d for d in on(f"tel.fleet.robot.{ASSET}.position")
                                      if distance(d["lat"], d["lon"], tlat, tlon) < 0.35], ARRIVE_S, "robot arrival")
    await wait_for(lambda: [d for d in on("evt.so.action.executed") if d["data"]["action_id"] == p["action_id"]], 10,
                   "evt.so.action.executed")
    report["goto"] = {"action_id": p["action_id"], "command_id": c["command_id"],
                      "arrived_within_m": round(distance(arrived[0]["lat"], arrived[0]["lon"], tlat, tlon), 3)}
    captured = next(d for t, s, d in seen if s == f"cmd.fleet.robot.{ASSET}.goto" and d["id"] == c["command_id"])

    # 2. an agent may propose but never confirm
    pa = await propose(hermes, goto_params(0.5, 0.0), kind="agent", who="agent:hermes/zeus/legacy")
    assert pa["status"] == "awaiting_confirmation", pa
    ca = await confirm(hermes, pa)
    assert ca["denial"]["code"] == "agent_cannot_confirm", ca
    report["agent_confirm"] = ca["denial"]["code"]

    # 3. a client cannot publish cmd.fleet.> (server permission), so the bridge never sees it
    n_before = len(on(f"evt.robot.{ASSET}.command"))
    await zeus.publish(f"cmd.fleet.robot.{ASSET}.goto", json.dumps(captured).encode())
    await zeus.flush()
    await asyncio.sleep(0.5)
    assert any("permissions violation" in e.lower() for e in perm_errors), perm_errors
    assert len(on(f"evt.robot.{ASSET}.command")) == n_before
    report["client_publish_cmd"] = "permissions violation; bridge received nothing"

    # 4. replaying the policy's own signed command is refused (nonce)
    t0 = time.time()
    await admin.publish(captured["subject"], json.dumps(captured).encode())
    r = await wait_for(lambda: [d["data"] for d in on(f"evt.robot.{ASSET}.command", t0)
                                if d["data"]["command_id"] == c["command_id"] and d["data"]["outcome"] == "refused"], 5, "replay refusal")
    assert r[0]["code"] == C.REPLAYED, r
    report["replay"] = r[0]["code"]

    # 5. an unsigned forged command is refused
    forged = {k: v for k, v in captured.items() if k != "sig"}
    forged["id"] = str(uuid.uuid4())
    t0 = time.time()
    await admin.publish(captured["subject"], json.dumps(forged).encode())
    r = await wait_for(lambda: [d["data"] for d in on(f"evt.robot.{ASSET}.command", t0) if d["data"]["command_id"] == forged["id"]], 5, "unsigned refusal")
    assert r[0]["code"] == C.BAD_SIGNATURE, r
    report["unsigned"] = r[0]["code"]

    # 6. inside the policy's 200 m zone but outside the bridge's 30 m fence: the bridge refuses
    p6 = await propose(zeus, goto_params(60.0, 0.0))
    c6 = await confirm(zeus, p6)
    r = await wait_for(lambda: [o for o in outcome_for(c6.get("command_id", "-")) if o["outcome"] == "refused"], 10, "bridge fence refusal")
    assert (r[0]["code"], r[0]["detail"]) == (C.PRECONDITION, "keep_in"), r
    await wait_for(lambda: [d for d in on("evt.so.action.refused") if d["data"]["action_id"] == p6["action_id"]], 10,
                   "evt.so.action.refused")
    report["bridge_fence"] = f"{r[0]['code']}/{r[0]['detail']} (policy approved {p6['action_id']})"

    # 7. outside the policy's zone: denied before any command exists
    p7 = await propose(zeus, goto_params(500.0, 0.0))
    assert p7["status"] == "denied" and p7["denial"]["code"] == "precondition_failed", p7
    report["policy_fence"] = p7["denial"]["code"]

    # 8. a non-sim target is denied by policy
    p8 = await propose(zeus, goto_params(1.0, 0.0), entity_id="tb1")
    assert p8["status"] == "denied" and p8["denial"]["code"] == "not_simulated", p8
    report["non_sim"] = p8["denial"]["code"]

    # 9. e-stop through the policy, then a goto is refused by the bridge's latch
    p9 = await propose(zeus, {}, verb="estop")
    c9 = await confirm(zeus, p9)
    await wait_for(lambda: [o for o in outcome_for(c9["command_id"]) if o["outcome"] == "executed"], 10, "estop executed")
    p10 = await propose(zeus, goto_params(0.0, 1.0))
    c10 = await confirm(zeus, p10)
    r = await wait_for(lambda: [o for o in outcome_for(c10.get("command_id", "-")) if o["outcome"] == "refused"], 10, "estop latch")
    assert r[0]["detail"] == "estopped", r
    report["estop_latch"] = r[0]["detail"]

    # Whole-run invariant: every cmd.fleet message follows a human approval of its action.
    approved = {d["data"]["action_id"] for d in on("evt.so.action.approved")}
    cmds = [d for t, s, d in seen if s.startswith("cmd.fleet.")]
    assert cmds and all(d["action_id"] in approved for d in cmds if d.get("source") == "service:policy")
    report["invariant"] = f"{len([d for d in cmds if d.get('source') == 'service:policy'])} policy commands, all human-approved"
    with open(os.path.join(st["logdir"], "e2e-report.json"), "w") as f:
        json.dump(report, f, indent=2)
    transcript.close()
    for c_ in (admin, entity, zeus, hermes):
        try:
            await c_.close()
        except Exception:
            pass
    del NatsTimeout
