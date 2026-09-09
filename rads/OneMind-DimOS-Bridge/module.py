"""
OneMind OS — DimOS NATS Bridge Module
======================================
Connects any DimOS robot to the OneMind NATS fabric.

What it does:
  • Subscribes to robot odometry/pose streams from DimOS
  • Publishes robot state to NATS:  fabric.robot.<robot_id>.odom
  • Converts robot pose → TAK CoT XML and publishes to:  fabric.robot.<robot_id>.cot
  • Subscribes to NATS command topic:  fabric.robot.<robot_id>.cmd
    and forwards vel commands back into DimOS as Twist messages
  • Exposes MCP skill: send_velocity_command (callable by Legacy)

Usage — drop into a DimOS blueprint:

    from omos.dimos.nats_bridge.module import NATSBridgeModule
    from dimos.core.coordination.blueprints import autoconnect
    from dimos.robot.unitree.go2.blueprints.smart.unitree_go2 import unitree_go2

    blueprint = autoconnect(
        unitree_go2(),
        NATSBridgeModule.blueprint(robot_id="go2-alpha"),
    )
    blueprint.build().loop()

Requirements:
    pip install nats-py   (nats-py >= 2.6)
    NATS_URL env var or pass nats_url= kwarg (default: nats://localhost:4222)
"""

from __future__ import annotations

import asyncio
import json
import math
import os
import threading
import time
import uuid
from datetime import datetime, timezone

import nats

from dimos.core.core import rpc
from dimos.core.module import Module
from dimos.core.stream import In, Out
from dimos.msgs.geometry_msgs.Twist import Twist
from dimos.msgs.nav_msgs.Odometry import Odometry


# ---------------------------------------------------------------------------
# CoT helper
# ---------------------------------------------------------------------------

def _quat_to_yaw_deg(qx: float, qy: float, qz: float, qw: float) -> float:
    """Quaternion → heading in degrees (0 = North, clockwise positive)."""
    yaw_rad = math.atan2(2.0 * (qw * qz + qx * qy),
                         1.0 - 2.0 * (qy * qy + qz * qz))
    return math.degrees(yaw_rad) % 360.0


def _build_cot(robot_id: str, odom: Odometry, uid_prefix: str = "OMOS-ROBOT") -> str:
    """
    Build a Cursor-on-Target (CoT) XML event from robot odometry.

    NOTE: DimOS odometry is in the robot's local frame (x/y metres from
    start pose).  Without a GPS fix we emit (lat=0, lon=0) as a placeholder.
    Swap in real GPS when available.
    """
    now = datetime.now(timezone.utc)
    stale = now.replace(second=now.second + 30)   # stale in 30 s
    ts   = now.strftime("%Y-%m-%dT%H:%M:%S.%f")[:-3] + "Z"
    st   = stale.strftime("%Y-%m-%dT%H:%M:%S.%f")[:-3] + "Z"

    uid  = f"{uid_prefix}-{robot_id}"
    hae  = odom.z          # height above ellipsoid ≈ z for ground robots
    course = _quat_to_yaw_deg(
        odom.orientation.x, odom.orientation.y,
        odom.orientation.z, odom.orientation.w,
    )
    speed  = math.hypot(odom.vx, odom.vy)

    # TAK type: a-n-G-E-V (ground vehicle, neutral)
    cot_type = "a-n-G-E-V"

    return f"""<?xml version="1.0" encoding="utf-8"?>
<event version="2.0"
       uid="{uid}"
       type="{cot_type}"
       how="m-g"
       time="{ts}"
       start="{ts}"
       stale="{st}">
  <point lat="0.0" lon="0.0" hae="{hae:.2f}" ce="9999999.0" le="9999999.0"/>
  <detail>
    <contact callsign="{robot_id}"/>
    <status battery="unknown"/>
    <track course="{course:.1f}" speed="{speed:.3f}"/>
    <remarks>DimOS robot on OneMind fabric | x={odom.x:.2f} y={odom.y:.2f} z={odom.z:.2f}</remarks>
    <__group name="OneMind-Robots" role="Robot"/>
  </detail>
</event>"""


# ---------------------------------------------------------------------------
# Module
# ---------------------------------------------------------------------------

class NATSBridgeModule(Module):
    """
    DimOS Module — bridges robot streams to the OneMind NATS fabric.

    Streams consumed:
        odom_in : Odometry  — robot pose + velocity (from nav stack)

    Streams produced:
        cmd_vel : Twist     — velocity command sourced from NATS

    Configuration (pass via constructor kwargs or env vars):
        robot_id   : str   — unique robot name on the fabric  (default: "robot-{uuid4[:8]}")
        nats_url   : str   — NATS server URL                  (default: NATS_URL env or nats://localhost:4222)
        publish_hz : float — max publish rate to NATS         (default: 10 Hz)
        cot_hz     : float — CoT publish rate                 (default: 1 Hz)
    """

    # DimOS stream declarations
    odom_in: In[Odometry]
    cmd_vel: Out[Twist]

    def __init__(
        self,
        robot_id: str | None = None,
        nats_url: str | None = None,
        publish_hz: float = 10.0,
        cot_hz: float = 1.0,
        **kwargs,
    ):
        super().__init__(**kwargs)
        self.robot_id   = robot_id or f"robot-{str(uuid.uuid4())[:8]}"
        self.nats_url   = nats_url or os.environ.get("NATS_URL", "nats://localhost:4222")
        self._pub_interval = 1.0 / max(publish_hz, 0.1)
        self._cot_interval = 1.0 / max(cot_hz, 0.01)

        self._nc: nats.aio.client.Client | None = None
        self._loop: asyncio.AbstractEventLoop | None = None
        self._thread: threading.Thread | None = None
        self._last_odom: Odometry | None = None
        self._last_pub_ts: float = 0.0
        self._last_cot_ts: float = 0.0

        # Subjects
        self._subj_odom = f"fabric.robot.{self.robot_id}.odom"
        self._subj_cot  = f"fabric.robot.{self.robot_id}.cot"
        self._subj_cmd  = f"fabric.robot.{self.robot_id}.cmd"
        self._subj_status = f"fabric.robot.{self.robot_id}.status"

    # ------------------------------------------------------------------
    # Lifecycle
    # ------------------------------------------------------------------

    @rpc
    def start(self) -> None:
        """Start the NATS bridge (called by DimOS blueprint engine)."""
        self._loop = asyncio.new_event_loop()
        self._thread = threading.Thread(
            target=self._run_event_loop, daemon=True, name=f"nats-bridge-{self.robot_id}"
        )
        self._thread.start()

        # Subscribe to incoming odometry from the DimOS stream
        self.odom_in.subscribe(self._on_odom)

        print(f"[NATSBridge] Started — robot_id={self.robot_id}  nats={self.nats_url}")
        print(f"[NATSBridge]   publish → {self._subj_odom}")
        print(f"[NATSBridge]   cot     → {self._subj_cot}")
        print(f"[NATSBridge]   listen  ← {self._subj_cmd}")

    @rpc
    def stop(self) -> None:
        """Gracefully disconnect from NATS."""
        if self._nc and self._loop:
            asyncio.run_coroutine_threadsafe(self._nc.drain(), self._loop)
        print(f"[NATSBridge] Stopped — {self.robot_id}")

    # ------------------------------------------------------------------
    # Internal
    # ------------------------------------------------------------------

    def _run_event_loop(self) -> None:
        asyncio.set_event_loop(self._loop)
        self._loop.run_until_complete(self._async_main())

    async def _async_main(self) -> None:
        self._nc = await nats.connect(self.nats_url)
        print(f"[NATSBridge] Connected to NATS @ {self.nats_url}")

        # Publish online status
        await self._nc.publish(
            self._subj_status,
            json.dumps({
                "robot_id": self.robot_id,
                "status": "online",
                "ts": time.time(),
                "nats_bridge": "omos-dimos-v1",
            }).encode(),
        )

        # Subscribe to incoming velocity commands from Legacy / NATS
        await self._nc.subscribe(self._subj_cmd, cb=self._on_nats_cmd)

        # Keep loop alive
        while True:
            await asyncio.sleep(1.0)

    def _on_odom(self, odom: Odometry) -> None:
        """Called by DimOS reactive stream on every new odometry message."""
        self._last_odom = odom
        now = time.monotonic()

        # Rate-limit NATS publish
        if now - self._last_pub_ts >= self._pub_interval:
            self._last_pub_ts = now
            self._publish_odom(odom)

        # Rate-limit CoT publish
        if now - self._last_cot_ts >= self._cot_interval:
            self._last_cot_ts = now
            self._publish_cot(odom)

    def _publish_odom(self, odom: Odometry) -> None:
        if not self._nc or not self._loop:
            return
        payload = json.dumps({
            "robot_id": self.robot_id,
            "ts": odom.ts,
            "frame_id": odom.frame_id,
            "x": odom.x, "y": odom.y, "z": odom.z,
            "roll": odom.roll, "pitch": odom.pitch, "yaw": odom.yaw,
            "vx": odom.vx, "vy": odom.vy, "vz": odom.vz,
            "wx": odom.wx, "wy": odom.wy, "wz": odom.wz,
            "qx": odom.orientation.x, "qy": odom.orientation.y,
            "qz": odom.orientation.z, "qw": odom.orientation.w,
        }).encode()
        asyncio.run_coroutine_threadsafe(
            self._nc.publish(self._subj_odom, payload), self._loop
        )

    def _publish_cot(self, odom: Odometry) -> None:
        if not self._nc or not self._loop:
            return
        cot_xml = _build_cot(self.robot_id, odom).encode()
        asyncio.run_coroutine_threadsafe(
            self._nc.publish(self._subj_cot, cot_xml), self._loop
        )

    async def _on_nats_cmd(self, msg: nats.aio.msg.Msg) -> None:
        """
        Receive a velocity command from NATS and push it into the DimOS
        cmd_vel stream so the robot actually moves.

        Expected payload (JSON):
            { "linear": {"x": 0.5, "y": 0.0, "z": 0.0},
              "angular": {"x": 0.0, "y": 0.0, "z": 0.3} }
        """
        try:
            data = json.loads(msg.data.decode())
            lin  = data.get("linear",  {})
            ang  = data.get("angular", {})
            twist = Twist(
                linear=[lin.get("x", 0.0), lin.get("y", 0.0), lin.get("z", 0.0)],
                angular=[ang.get("x", 0.0), ang.get("y", 0.0), ang.get("z", 0.3)],
            )
            self.cmd_vel.publish(twist)
            print(f"[NATSBridge] CMD ← vx={twist.linear.x:.2f} wz={twist.angular.z:.2f}")
        except Exception as e:
            print(f"[NATSBridge] CMD parse error: {e}")

    # ------------------------------------------------------------------
    # MCP Skills (callable by Legacy via dimos mcp call ...)
    # ------------------------------------------------------------------

    @rpc
    def send_velocity_command(
        self,
        linear_x: float = 0.0,
        linear_y: float = 0.0,
        angular_z: float = 0.0,
    ) -> str:
        """
        MCP skill: directly command the robot's velocity.
        Called by Legacy: dimos mcp call send_velocity_command --arg linear_x=0.5
        """
        twist = Twist(
            linear=[linear_x, linear_y, 0.0],
            angular=[0.0, 0.0, angular_z],
        )
        self.cmd_vel.publish(twist)
        return f"OK: vx={linear_x:.2f} vy={linear_y:.2f} wz={angular_z:.2f}"

    @rpc
    def get_robot_state(self) -> dict:
        """MCP skill: return current robot pose/velocity snapshot."""
        if self._last_odom is None:
            return {"status": "no_data", "robot_id": self.robot_id}
        o = self._last_odom
        return {
            "robot_id": self.robot_id,
            "ts": o.ts,
            "position": {"x": o.x, "y": o.y, "z": o.z},
            "orientation": {"roll": o.roll, "pitch": o.pitch, "yaw": o.yaw},
            "velocity": {"vx": o.vx, "vy": o.vy, "wz": o.wz},
            "nats_subjects": {
                "odom": self._subj_odom,
                "cot":  self._subj_cot,
                "cmd":  self._subj_cmd,
            },
        }