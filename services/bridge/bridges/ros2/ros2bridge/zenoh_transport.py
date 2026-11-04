"""Zenoh transport: talks to zenoh-bridge-ros2dds (eclipse-zenoh Python, Apache-2.0 OR EPL-2.0).

Key expressions: zenoh-bridge-ros2dds routes a ROS interface ``/name`` to the Zenoh key
``name``, prefixed with its configured ``namespace`` ("a cmd_vel topic in the robot will
be seen as namespace/cmd_vel outside the robot", DEFAULT_CONFIG.json5 of
eclipse-zenoh/zenoh-plugin-ros2dds 1.10.2). Actions are routed as queryables
``{action}/_action/send_goal`` and ``{action}/_action/cancel_goal`` with CDR request and
reply payloads. The action key names and the absence of a DDS request header in the
payload are not stated in that README and are **unverified** until the owner run in
OWNER-STEPS.md (the scripted robot in sim/ implements them as described here).
"""

from __future__ import annotations

import asyncio
import logging
import os
import time

from . import msgs
from .cdr import CDRError
from .transport import GoalResult, Topics, Transport, TransportError, TransportTimeout

log = logging.getLogger("ros2bridge.zenoh")


def key(namespace: str, name: str) -> str:
    ns, name = namespace.strip("/"), name.strip("/")
    return f"{ns}/{name}" if ns else name


class ZenohTransport(Transport):
    name = "zenoh"

    def __init__(self, connect: list[str], mode: str = "client", listen: list[str] | None = None,
                 namespace: str = "", topics: Topics | None = None, query_timeout: float = 3.0,
                 multicast_scouting: bool = False):
        self.connect, self.mode, self.listen = connect, mode, listen or []
        self.namespace, self.topics = namespace, topics or Topics()
        self.query_timeout = query_timeout
        self.multicast = multicast_scouting
        self._s = None
        self._subs = []
        self._cmd_vel = None
        self._decode_errors = 0
        self._last_err_log = 0.0
        self._rx = 0

    def k(self, name: str) -> str:
        return key(self.namespace, name)

    def _config(self):
        import json

        import zenoh

        cfg = zenoh.Config()
        cfg.insert_json5("mode", json.dumps(self.mode))
        if self.connect:
            cfg.insert_json5("connect/endpoints", json.dumps(self.connect))
        if self.listen:
            cfg.insert_json5("listen/endpoints", json.dumps(self.listen))
        cfg.insert_json5("scouting/multicast/enabled", "true" if self.multicast else "false")
        return cfg

    def _deliver(self, loop, decode, cb, what):
        def handler(sample):
            self._rx += 1
            try:
                m = decode(bytes(sample.payload.to_bytes()))
            except (CDRError, ValueError) as e:
                self._decode_errors += 1
                now = time.monotonic()
                if now - self._last_err_log > 10:
                    self._last_err_log = now
                    log.warning("undecodable %s on %s: %s (%d so far)", what, sample.key_expr, e, self._decode_errors)
                return
            loop.call_soon_threadsafe(cb, m)
        return handler

    async def start(self, on_odom, on_pose, on_battery) -> None:
        import zenoh

        loop = asyncio.get_running_loop()
        self._s = await asyncio.to_thread(zenoh.open, self._config())
        t = self.topics
        self._subs.append(self._s.declare_subscriber(self.k(t.odom), self._deliver(loop, msgs.decode_odometry, on_odom, "odometry")))
        if t.pose:
            self._subs.append(self._s.declare_subscriber(
                self.k(t.pose), self._deliver(loop, msgs.decode_pose_with_covariance_stamped, on_pose, "pose")))
        if t.battery:
            self._subs.append(self._s.declare_subscriber(
                self.k(t.battery), self._deliver(loop, msgs.decode_battery_state, on_battery, "battery")))
        self._cmd_vel = self._s.declare_publisher(self.k(t.cmd_vel))
        log.info("zenoh %s session, connect=%s listen=%s; keys %s, %s, %s, %s/_action/*",
                 self.mode, self.connect, self.listen, self.k(t.odom), self.k(t.battery), self.k(t.cmd_vel),
                 self.k(t.nav_action))

    def _query(self, k: str, payload: bytes) -> list[bytes]:
        if self._s is None:
            raise TransportError("zenoh session not open")
        out = []
        for reply in self._s.get(k, payload=payload, timeout=self.query_timeout):
            if reply.ok is not None:
                out.append(bytes(reply.ok.payload.to_bytes()))
            else:
                raise TransportError(f"{k}: error reply {reply.err.payload.to_string()!r}")
        return out

    async def send_goal(self, goal) -> GoalResult:
        gid = os.urandom(16)
        payload = msgs.encode_send_goal_request(msgs.SendGoalRequest(gid, goal, ""))
        k = self.k(self.topics.nav_action + "/_action/send_goal")
        replies = await asyncio.to_thread(self._query, k, payload)
        if not replies:
            raise TransportTimeout(f"no reply on {k} within {self.query_timeout}s (is Nav2 running?)")
        try:
            resp = msgs.decode_send_goal_response(replies[0])
        except CDRError as e:
            raise TransportError(f"{k}: undecodable reply: {e}") from e
        return GoalResult(resp.accepted, gid.hex(), k)

    async def _cancel_all(self) -> str:
        k = self.k(self.topics.nav_action + "/_action/cancel_goal")
        replies = await asyncio.to_thread(self._query, k, msgs.encode_cancel_goal_request())
        if not replies:
            raise TransportTimeout(f"no reply on {k} within {self.query_timeout}s")
        resp = msgs.decode_cancel_goal_response(replies[0])
        return f"{k} return_code={resp.return_code} canceling={len(resp.goals)}"

    def _zero(self, n: int = 3) -> None:
        if self._cmd_vel is None:
            raise TransportError("zenoh session not open")
        z = msgs.encode_twist(msgs.Twist())
        for _ in range(n):
            self._cmd_vel.put(z)

    async def hold(self, here) -> str:
        detail = await self._cancel_all()
        self._zero(1)
        return detail

    async def estop(self, here=None) -> str:
        # Zero velocity first, without waiting on anything; then cancel goals best effort.
        self._zero(3)
        try:
            return "cmd_vel zeroed; " + await self._cancel_all()
        except (TransportTimeout, TransportError) as e:
            return f"cmd_vel zeroed; cancel_goal failed: {e}"

    async def close(self) -> None:
        if self._s is not None:
            s, self._s = self._s, None
            await asyncio.to_thread(s.close)

    def stats(self) -> dict:
        return {"rx": self._rx, "decode_errors": self._decode_errors}
