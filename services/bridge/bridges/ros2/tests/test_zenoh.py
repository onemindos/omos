"""Zenoh transport against the scripted robot over real Zenoh sessions on localhost.

The scripted robot opens a Zenoh router session, as zenoh-bridge-ros2dds does by default
(router mode since v0.11), and the bridge connects as a client."""

import asyncio
import json
import socket
import threading

import pytest

from cmdtest import build, new_key
from ros2bridge import msgs
from ros2bridge.bridge import Bridge, BridgeConfig
from ros2bridge.fence import Fence
from ros2bridge.geo import LocalFrame
from ros2bridge.zenoh_transport import ZenohTransport, key
from turtlebot_sim import TurtlebotSim, open_session, serve

zenoh = pytest.importorskip("zenoh")
PUB, PRIV = new_key()
FRAME = LocalFrame(30.2672, -97.7431, 150.0, 90.0)


def free_port():
    s = socket.socket()
    s.bind(("127.0.0.1", 0))
    p = s.getsockname()[1]
    s.close()
    return p


@pytest.fixture
def robot():
    port = free_port()
    sim = TurtlebotSim("sim_tb1", x=0.0, y=0.0, arena_m=20.0)
    sess = open_session("router", [f"tcp/127.0.0.1:{port}"], [])
    stop = threading.Event()
    th = threading.Thread(target=serve, args=(sim, sess, stop), daemon=True)
    th.start()
    yield sim, f"tcp/127.0.0.1:{port}"
    stop.set()
    th.join(5)
    sess.close()


class Bus:
    def __init__(self):
        self.msgs = []

    async def publish(self, subject, data, msg_id=""):
        self.msgs.append((subject, json.loads(data)))


def test_key_mapping():
    assert key("", "odom") == "odom" and key("/bot1", "/odom") == "bot1/odom"
    assert key("bot1/", "navigate_to_pose/_action/send_goal") == "bot1/navigate_to_pose/_action/send_goal"


def test_transport_odom_goal_cancel(robot):
    sim, ep = robot

    async def go():
        t = ZenohTransport([ep], "client", namespace="sim_tb1", query_timeout=3.0)
        odoms, bats = [], []
        await t.start(odoms.append, lambda p: None, bats.append)
        for _ in range(50):
            if odoms and bats:
                break
            await asyncio.sleep(0.1)
        r = await t.send_goal(msgs.PoseStamped(msgs.Header(msgs.Time.now(), "map"), msgs.Pose.from_xy_yaw(1.0, 0.0, 0)))
        rej = await t.send_goal(msgs.PoseStamped(msgs.Header(msgs.Time.now(), "map"), msgs.Pose.from_xy_yaw(99.0, 0, 0)))
        await asyncio.sleep(1.0)
        moved = sim.x
        hold = await t.hold(None)
        est = await t.estop()
        stats = t.stats()
        await t.close()
        return odoms, bats, r, rej, moved, hold, est, stats

    odoms, bats, r, rej, moved, hold, est, stats = asyncio.run(go())
    assert odoms and isinstance(odoms[0], msgs.Odometry) and odoms[0].child_frame_id == "base_footprint"
    assert bats and 0 < bats[0].percentage <= 1
    assert r.accepted and len(r.goal_id) == 32 and not rej.accepted
    assert moved > 0.05  # the scripted robot drove toward the goal
    assert "return_code=0" in hold and "cmd_vel zeroed" in est
    assert sim.cancels == 2 and sim.twists >= 3 and stats["decode_errors"] == 0


def test_bridge_end_to_end_over_zenoh(robot):
    """Signed goto -> zenoh action goal -> scripted robot drives -> telemetry shows arrival."""
    sim, ep = robot
    fence = Fence.parse(json.dumps({"keep_in": [{"id": "yard", "circle": {"lat": 30.2672, "lon": -97.7431, "radius_m": 30}}]}))

    async def go():
        bus = Bus()
        t = ZenohTransport([ep], "client", namespace="sim_tb1")
        b = Bridge(BridgeConfig("sim-tb1", FRAME, fence=fence, keys={"policy-test": PUB}, position_interval=0.2), bus, t)
        await t.start(b.on_odom, b.on_pose, b.on_battery)
        for _ in range(50):
            if b.has_pos:
                break
            await asyncio.sleep(0.1)
        lat, lon, _ = FRAME.to_wgs84(0.6, 0.3)
        await b.handle_command(*build(PRIV, verb="goto", data={"lat": lat, "lon": lon, "hae_m": 150.0}, task_id="tsk-z"))
        for _ in range(150):
            if b.goal is None and b.mode == "idle":
                break
            await asyncio.sleep(0.1)
        await t.close()
        return b, bus

    b, bus = asyncio.run(go())
    outs = [m["data"]["outcome"] for s, m in bus.msgs if s == "evt.robot.sim-tb1.command"]
    assert outs == ["acknowledged", "executed"]
    tasks = [m["data"]["state"] for s, m in bus.msgs if s == "evt.robot.sim-tb1.task"]
    assert tasks == ["received", "executing", "done_ok"]
    pos = [m for s, m in bus.msgs if s == "tel.fleet.robot.sim-tb1.position"]
    assert len(pos) > 3 and abs(b.x - 0.6) < 0.35 and abs(b.y - 0.3) < 0.35
