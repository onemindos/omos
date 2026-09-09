"""
OneMind OS — Legacy Robot Commander
=====================================
Legacy-side interface for commanding DimOS robots over the NATS fabric.

Usage from within Legacy (or any NATS client on the fabric):

    from omos.dimos.nats_bridge.legacy_commander import RobotCommander

    robot = RobotCommander("go2-alpha")  # matches ROBOT_ID in DimOS
    await robot.move(linear_x=0.5, angular_z=0.2)
    state = await robot.get_state()

Or via NATS CLI from the pod:
    nats pub fabric.robot.go2-alpha.cmd '{"linear":{"x":0.5},"angular":{"z":0.3}}'
    nats sub fabric.robot.go2-alpha.odom   # watch telemetry live
    nats sub fabric.robot.go2-alpha.cot    # watch TAK CoT events live
"""

from __future__ import annotations

import asyncio
import json
import time
from typing import Any

import nats


class RobotCommander:
    """
    Async interface for Legacy to command any DimOS robot on the fabric.

    Each method is fire-and-forget (pub) or request-reply (req).
    """

    def __init__(
        self,
        robot_id: str,
        nats_url: str = "nats://localhost:4222",
    ):
        self.robot_id = robot_id
        self.nats_url = nats_url
        self._nc: nats.aio.client.Client | None = None

        self._subj_cmd    = f"fabric.robot.{robot_id}.cmd"
        self._subj_odom   = f"fabric.robot.{robot_id}.odom"
        self._subj_cot    = f"fabric.robot.{robot_id}.cot"
        self._subj_status = f"fabric.robot.{robot_id}.status"

    async def connect(self) -> None:
        if not self._nc:
            self._nc = await nats.connect(self.nats_url)

    async def disconnect(self) -> None:
        if self._nc:
            await self._nc.drain()
            self._nc = None

    # ------------------------------------------------------------------
    # Commands
    # ------------------------------------------------------------------

    async def move(
        self,
        linear_x: float = 0.0,
        linear_y: float = 0.0,
        linear_z: float = 0.0,
        angular_x: float = 0.0,
        angular_y: float = 0.0,
        angular_z: float = 0.0,
    ) -> None:
        """Publish a velocity command to the robot."""
        await self._ensure_connected()
        payload = json.dumps({
            "linear":  {"x": linear_x,  "y": linear_y,  "z": linear_z},
            "angular": {"x": angular_x, "y": angular_y, "z": angular_z},
            "ts": time.time(),
            "source": "legacy",
        }).encode()
        await self._nc.publish(self._subj_cmd, payload)

    async def stop(self) -> None:
        """Immediately stop the robot (zero velocity)."""
        await self.move()

    async def forward(self, speed: float = 0.5, duration_s: float = 0.0) -> None:
        """Move forward at `speed` m/s. If duration_s > 0, stop after that time."""
        await self.move(linear_x=speed)
        if duration_s > 0:
            await asyncio.sleep(duration_s)
            await self.stop()

    async def turn(self, angular_z: float = 0.5, duration_s: float = 0.0) -> None:
        """Turn in place. Positive = counter-clockwise."""
        await self.move(angular_z=angular_z)
        if duration_s > 0:
            await asyncio.sleep(duration_s)
            await self.stop()

    # ------------------------------------------------------------------
    # Telemetry
    # ------------------------------------------------------------------

    async def get_state(self, timeout: float = 2.0) -> dict[str, Any] | None:
        """
        Subscribe to odometry for one message and return robot state.
        Returns None if robot is offline or not publishing.
        """
        await self._ensure_connected()
        result: dict | None = None
        event = asyncio.Event()

        async def handler(msg: nats.aio.msg.Msg):
            nonlocal result
            try:
                result = json.loads(msg.data.decode())
            except Exception:
                pass
            event.set()

        sub = await self._nc.subscribe(self._subj_odom, cb=handler)
        try:
            await asyncio.wait_for(event.wait(), timeout=timeout)
        except asyncio.TimeoutError:
            pass
        finally:
            await sub.unsubscribe()
        return result

    async def watch_telemetry(
        self,
        on_odom: callable | None = None,
        on_cot:  callable | None = None,
        duration_s: float = 10.0,
    ) -> None:
        """
        Watch robot telemetry streams for `duration_s` seconds.
        Calls on_odom(dict) and/or on_cot(str) for each message.
        """
        await self._ensure_connected()
        subs = []

        if on_odom:
            async def _odom_handler(msg):
                on_odom(json.loads(msg.data.decode()))
            subs.append(await self._nc.subscribe(self._subj_odom, cb=_odom_handler))

        if on_cot:
            async def _cot_handler(msg):
                on_cot(msg.data.decode())
            subs.append(await self._nc.subscribe(self._subj_cot, cb=_cot_handler))

        await asyncio.sleep(duration_s)

        for sub in subs:
            await sub.unsubscribe()

    # ------------------------------------------------------------------
    # Discovery
    # ------------------------------------------------------------------

    @staticmethod
    async def discover_robots(
        nats_url: str = "nats://localhost:4222",
        timeout: float = 3.0,
    ) -> list[str]:
        """
        Listen on fabric.robot.*.status for `timeout` seconds and return
        the list of online robot_ids.
        """
        nc = await nats.connect(nats_url)
        found: list[str] = []

        async def handler(msg: nats.aio.msg.Msg):
            try:
                data = json.loads(msg.data.decode())
                rid = data.get("robot_id")
                if rid and rid not in found:
                    found.append(rid)
            except Exception:
                pass

        sub = await nc.subscribe("fabric.robot.*.status", cb=handler)
        await asyncio.sleep(timeout)
        await sub.unsubscribe()
        await nc.drain()
        return found

    # ------------------------------------------------------------------
    # Internal
    # ------------------------------------------------------------------

    async def _ensure_connected(self) -> None:
        if not self._nc:
            await self.connect()