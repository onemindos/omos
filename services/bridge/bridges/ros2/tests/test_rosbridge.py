"""rosbridge fallback against a fake rosbridge v2 server (websockets)."""

import asyncio
import json

import pytest

from ros2bridge import msgs
from ros2bridge.rosbridge_transport import RosbridgeTransport, odometry_from_json

websockets = pytest.importorskip("websockets")

ODOM = {"header": {"stamp": {"sec": 5, "nanosec": 6}, "frame_id": "odom"}, "child_frame_id": "base_footprint",
        "pose": {"pose": {"position": {"x": 1.0, "y": 2.0, "z": 0.0}, "orientation": {"x": 0, "y": 0, "z": 0, "w": 1}},
                 "covariance": [0.0] * 36},
        "twist": {"twist": {"linear": {"x": 0.1, "y": 0, "z": 0}, "angular": {"x": 0, "y": 0, "z": 0.2}},
                  "covariance": [0.0] * 36}}


def test_odometry_json_decode():
    o = odometry_from_json(ODOM)
    assert (o.pose.x, o.pose.y, o.twist.lx, o.twist.az, o.header.stamp.sec) == (1.0, 2.0, 0.1, 0.2, 5)
    with pytest.raises(ValueError):
        odometry_from_json({"pose": {"pose": {"position": {"x": "1"}}}})


def test_transport_against_fake_server():
    async def go():
        got = []

        async def handler(ws):
            async for raw in ws:
                m = json.loads(raw)
                got.append(m)
                if m["op"] == "subscribe" and m["topic"] == "/bot1/odom":
                    await ws.send(json.dumps({"op": "publish", "topic": "/bot1/odom", "msg": ODOM}))
                    await ws.send("not json")
                    await ws.send(json.dumps({"op": "publish", "topic": "/bot1/battery_state",
                                              "msg": {"percentage": 0.5, "voltage": 12.0, "present": True}}))

        async with websockets.serve(handler, "127.0.0.1", 0) as srv:
            port = srv.sockets[0].getsockname()[1]
            t = RosbridgeTransport(f"ws://127.0.0.1:{port}", namespace="bot1")
            odoms, bats = [], []
            await t.start(odoms.append, lambda p: None, bats.append)
            assert await t.wait_connected(5)
            for _ in range(50):
                if odoms and bats:
                    break
                await asyncio.sleep(0.05)
            r = await t.send_goal(msgs.PoseStamped(msgs.Header(msgs.Time(1, 2), "map"), msgs.Pose.from_xy_yaw(3, 4, 0)))
            await t.hold(msgs.PoseStamped(msgs.Header(), msgs.Pose()))
            await asyncio.sleep(0.2)  # let the server receive everything sent so far
            n_before = len(got)
            await t.estop(msgs.PoseStamped(msgs.Header(), msgs.Pose.from_xy_yaw(7, 8, 0)))
            await asyncio.sleep(0.2)
            got_estop = got[n_before:]
            stats = t.stats()
            await t.close()
            return got, odoms, bats, r, stats, got_estop

    got, odoms, bats, r, stats, est = asyncio.run(go())
    ops = [(m["op"], m.get("topic")) for m in got]
    assert ("subscribe", "/bot1/odom") in ops and ("subscribe", "/bot1/battery_state") in ops
    assert ("advertise", "/bot1/goal_pose") in ops and ("advertise", "/bot1/cmd_vel") in ops
    goal = next(m for m in got if m["op"] == "publish" and m["topic"] == "/bot1/goal_pose")
    assert goal["msg"]["pose"]["position"] == {"x": 3, "y": 4, "z": 0.0} and goal["msg"]["header"]["frame_id"] == "map"
    zero = [m for m in got if m["op"] == "publish" and m["topic"] == "/bot1/cmd_vel"]
    assert len(zero) == 4 and zero[0]["msg"]["linear"]["x"] == 0.0
    assert odoms[0].pose.x == 1.0 and bats[0].percentage == 0.5
    assert r.accepted and "no goal ack" in r.detail and stats["decode_errors"] == 1
    # estop: zero cmd_vel first, then re-target the goal to the current pose and cancel all goals
    assert [(m["op"], m.get("topic") or m.get("service")) for m in est] == [
        ("publish", "/bot1/cmd_vel"), ("publish", "/bot1/goal_pose"),
        ("call_service", "/bot1/navigate_to_pose/_action/cancel_goal"),
        ("publish", "/bot1/cmd_vel"), ("publish", "/bot1/cmd_vel")]
    assert est[1]["msg"]["pose"]["position"] == {"x": 7, "y": 8, "z": 0.0}
    assert est[2]["type"] == "action_msgs/srv/CancelGoal"
    assert est[2]["args"]["goal_info"]["goal_id"]["uuid"] == [0] * 16


def test_send_without_connection_raises():
    from ros2bridge.transport import TransportError

    async def go():
        t = RosbridgeTransport("ws://127.0.0.1:1")
        with pytest.raises(TransportError):
            await t.estop()
    asyncio.run(go())
