"""rosbridge fallback transport (rosbridge_suite v2 protocol, BSD-3-Clause server).

Uses the documented ops ``subscribe``, ``advertise`` and ``publish``
(RobotWebTools/rosbridge_suite ROSBRIDGE_PROTOCOL.md). Messages are JSON with ROS 2 field
names (``header.stamp.sec`` / ``nanosec``). Goals are published on the ``goal_pose`` topic
(geometry_msgs/msg/PoseStamped), so there is no accepted/rejected answer: the bridge
reports ``executed`` once the publish was written to the socket and says so in the
outcome. Prefer the Zenoh transport where it can run.
"""

from __future__ import annotations

import asyncio
import itertools
import json
import logging
import os

from . import msgs
from .transport import GoalResult, Topics, Transport, TransportError

log = logging.getLogger("ros2bridge.rosbridge")


def _n(d, *path, default=0.0):
    for p in path:
        if not isinstance(d, dict) or p not in d:
            return default
        d = d[p]
    return d


def _f(v) -> float:
    if v is None:  # rosbridge sends NaN as null
        return float("nan")
    if isinstance(v, bool) or not isinstance(v, (int, float)):
        raise ValueError(f"not a number: {v!r}")
    return float(v)


def _header(d) -> msgs.Header:
    return msgs.Header(msgs.Time(int(_n(d, "stamp", "sec", default=0)), int(_n(d, "stamp", "nanosec", default=0))),
                       str(_n(d, "frame_id", default="")))


def _pose(d) -> msgs.Pose:
    p, o = _n(d, "position", default={}), _n(d, "orientation", default={})
    return msgs.Pose(_f(_n(p, "x")), _f(_n(p, "y")), _f(_n(p, "z")),
                     _f(_n(o, "x")), _f(_n(o, "y")), _f(_n(o, "z")), _f(_n(o, "w", default=1.0)))


def odometry_from_json(m: dict) -> msgs.Odometry:
    tw = _n(m, "twist", "twist", default={})
    lin, ang = _n(tw, "linear", default={}), _n(tw, "angular", default={})
    return msgs.Odometry(_header(_n(m, "header", default={})), str(_n(m, "child_frame_id", default="")),
                         _pose(_n(m, "pose", "pose", default={})),
                         msgs.Twist(_f(_n(lin, "x")), _f(_n(lin, "y")), _f(_n(lin, "z")),
                                    _f(_n(ang, "x")), _f(_n(ang, "y")), _f(_n(ang, "z"))))


def pose_cov_from_json(m: dict) -> msgs.PoseStamped:
    return msgs.PoseStamped(_header(_n(m, "header", default={})), _pose(_n(m, "pose", "pose", default={})))


def battery_from_json(m: dict) -> msgs.BatteryState:
    return msgs.BatteryState(_header(_n(m, "header", default={})), _f(m.get("voltage", float("nan"))),
                             _f(m.get("percentage", float("nan"))), bool(m.get("present", True)))


def pose_stamped_to_json(p: msgs.PoseStamped) -> dict:
    return {"header": {"stamp": {"sec": p.header.stamp.sec, "nanosec": p.header.stamp.nanosec},
                       "frame_id": p.header.frame_id},
            "pose": {"position": {"x": p.pose.x, "y": p.pose.y, "z": p.pose.z},
                     "orientation": {"x": p.pose.qx, "y": p.pose.qy, "z": p.pose.qz, "w": p.pose.qw}}}


ZERO_TWIST = {"linear": {"x": 0.0, "y": 0.0, "z": 0.0}, "angular": {"x": 0.0, "y": 0.0, "z": 0.0}}


def _topic(namespace: str, name: str) -> str:
    ns, name = namespace.strip("/"), name.strip("/")
    return f"/{ns}/{name}" if ns else f"/{name}"


class RosbridgeTransport(Transport):
    name = "rosbridge"

    def __init__(self, url: str, namespace: str = "", topics: Topics | None = None, reconnect_s: float = 2.0):
        self.url, self.namespace, self.topics = url, namespace, topics or Topics()
        self.reconnect_s = reconnect_s
        self._ws = None
        self._task = None
        self._ids = itertools.count(1)
        self._connected = asyncio.Event()
        self._closing = False
        self._rx = 0
        self._decode_errors = 0

    def t(self, name: str) -> str:
        return _topic(self.namespace, name)

    async def start(self, on_odom, on_pose, on_battery) -> None:
        self._cbs = {self.t(self.topics.odom): (odometry_from_json, on_odom, "nav_msgs/msg/Odometry")}
        if self.topics.pose:
            self._cbs[self.t(self.topics.pose)] = (pose_cov_from_json, on_pose,
                                                   "geometry_msgs/msg/PoseWithCovarianceStamped")
        if self.topics.battery:
            self._cbs[self.t(self.topics.battery)] = (battery_from_json, on_battery, "sensor_msgs/msg/BatteryState")
        self._task = asyncio.create_task(self._run())

    async def _send(self, msg: dict) -> None:
        ws = self._ws
        if ws is None:
            raise TransportError("rosbridge not connected")
        try:
            await ws.send(json.dumps(msg))
        except Exception as e:  # websockets raises ConnectionClosed subclasses
            raise TransportError(f"rosbridge send: {e}") from e

    async def _run(self) -> None:
        import websockets

        while not self._closing:
            try:
                async with websockets.connect(self.url, open_timeout=5, max_size=2**22) as ws:
                    self._ws = ws
                    for topic, (_, _, typ) in self._cbs.items():
                        await ws.send(json.dumps({"op": "subscribe", "id": f"sub-{next(self._ids)}", "topic": topic,
                                                  "type": typ, "queue_length": 1}))
                    for topic, typ in ((self.t(self.topics.goal_topic), "geometry_msgs/msg/PoseStamped"),
                                       (self.t(self.topics.cmd_vel), "geometry_msgs/msg/Twist")):
                        await ws.send(json.dumps({"op": "advertise", "id": f"adv-{next(self._ids)}",
                                                  "topic": topic, "type": typ}))
                    self._connected.set()
                    log.info("rosbridge connected %s", self.url)
                    async for raw in ws:
                        self._on_message(raw)
            except asyncio.CancelledError:
                raise
            except Exception as e:
                log.warning("rosbridge %s: %s; retrying in %.1fs", self.url, e, self.reconnect_s)
            finally:
                self._ws = None
                self._connected.clear()
            if not self._closing:
                await asyncio.sleep(self.reconnect_s)

    def _on_message(self, raw) -> None:
        self._rx += 1
        try:
            m = json.loads(raw)
            if m.get("op") != "publish":
                return
            entry = self._cbs.get(m.get("topic"))
            if entry is None:
                return
            decode, cb, _ = entry
            cb(decode(m.get("msg") or {}))
        except (ValueError, TypeError, AttributeError) as e:
            self._decode_errors += 1
            if self._decode_errors in (1, 10, 100) or self._decode_errors % 1000 == 0:
                log.warning("undecodable rosbridge message: %s (%d so far)", e, self._decode_errors)

    async def send_goal(self, goal) -> GoalResult:
        await self._send({"op": "publish", "topic": self.t(self.topics.goal_topic), "msg": pose_stamped_to_json(goal)})
        return GoalResult(True, os.urandom(16).hex(), f"rosbridge publish {self.t(self.topics.goal_topic)} (no goal ack)")

    async def hold(self, here) -> str:
        # No action cancel on this path: re-target Nav2 at the current pose, then stop.
        await self._send({"op": "publish", "topic": self.t(self.topics.goal_topic), "msg": pose_stamped_to_json(here)})
        await self._send({"op": "publish", "topic": self.t(self.topics.cmd_vel), "msg": ZERO_TWIST})
        return f"goal re-targeted to the current pose on {self.t(self.topics.goal_topic)}; cmd_vel zeroed"

    async def estop(self) -> str:
        for _ in range(3):
            await self._send({"op": "publish", "topic": self.t(self.topics.cmd_vel), "msg": ZERO_TWIST})
        return f"cmd_vel zeroed on {self.t(self.topics.cmd_vel)} (rosbridge cannot cancel the Nav2 goal)"

    async def wait_connected(self, timeout: float) -> bool:
        try:
            await asyncio.wait_for(self._connected.wait(), timeout)
            return True
        except asyncio.TimeoutError:
            return False

    async def close(self) -> None:
        self._closing = True
        if self._task:
            self._task.cancel()
            try:
                await self._task
            except (asyncio.CancelledError, Exception):
                pass

    def stats(self) -> dict:
        return {"rx": self._rx, "decode_errors": self._decode_errors}
