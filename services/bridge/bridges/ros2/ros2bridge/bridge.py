"""The bridge core: ROS 2 robot state to fleet telemetry, verified fleet commands to ROS 2.

Subjects (locked grammar; asset-command.md section 1, robotics-cot-profile.md section 1):
  publishes  tel.fleet.robot.{id}.position   doctrine section 9.2 geo-envelope
             tel.fleet.robot.{id}.status     same envelope, link/battery/mode in meta
             evt.robot.{id}.command          command outcome (asset-command.md section 6.4)
             evt.robot.{id}.task             task status for task-bearing commands (mission-task.md section 7)
  subscribes cmd.fleet.robot.{id}.>          core subscription only: no queued replay (section 4)

The command flow mirrors the MAVLink bridge (omos PR #550): verify (command.py), check
robot state and the bridge's own geofence, report ``acknowledged``, execute, then report
``executed`` (the robot accepted the goal) or ``failed``.
"""

from __future__ import annotations

import asyncio
import json
import logging
import math
import uuid
from dataclasses import dataclass, field
from datetime import datetime, timedelta, timezone

from . import command as C
from . import msgs
from .fence import Fence
from .fence import check as fence_check
from .geo import LocalFrame
from .transport import Transport, TransportError, TransportTimeout

log = logging.getLogger("ros2bridge")

MAX_TRACKED_TASKS = 4096
MOVING_MPS = 0.05  # below this a ground robot is "stationary" (turtlebot3 burger max is 0.22 m/s)


@dataclass
class BridgeConfig:
    asset_id: str
    frame: LocalFrame
    callsign: str = ""
    source: str = "bridge-ros2"
    kind: str = "ugv"
    fence: Fence | None = None
    keys: dict = field(default_factory=dict)
    sig_required: bool = True
    max_ttl: timedelta = timedelta(seconds=60)
    clock_skew: timedelta = timedelta(seconds=5)
    position_interval: float = 0.5
    status_interval: float = 1.0
    stale_s: int = 10
    link_timeout: float = 3.0
    goal_tolerance_m: float = 0.3
    home: tuple | None = None  # (lat, lon); default: first position heard
    map_frame: str = "map"
    use_odom_position: bool = True  # False when a map-frame pose topic (amcl_pose) is configured

    def __post_init__(self):
        if not self.asset_id.startswith("sim-"):
            raise ValueError(f"ASSET_ID {self.asset_id!r} must start with sim- (simulated assets only, decision 0004 #13)")
        if not self.callsign:
            self.callsign = self.asset_id.upper()


class Bus:
    """Publishes to NATS. msg_id, when set, becomes the Nats-Msg-Id header so JetStream
    de-duplicates (EVT has a 2 m duplicate window, platform/nats/streams/evt.yaml)."""

    async def publish(self, subject: str, data: bytes, msg_id: str = "") -> None:  # pragma: no cover
        raise NotImplementedError


def rfc3339ms(t: datetime) -> str:
    """RFC 3339 UTC with milliseconds, as the omos envelopes require."""
    u = t.astimezone(timezone.utc)
    return u.strftime("%Y-%m-%dT%H:%M:%S.") + f"{u.microsecond // 1000:03d}Z"


def _rnd(v: float, places: int) -> float:
    if v is None or not math.isfinite(v):
        return -1.0
    return round(v, places) + 0.0


class Bridge:
    def __init__(self, cfg: BridgeConfig, bus: Bus, transport: Transport, clock=C.utcnow):
        self.cfg, self.bus, self.transport, self.clock = cfg, bus, transport, clock
        self.verifier = C.Verifier(C.Config("robot", cfg.asset_id, cfg.keys, cfg.sig_required,
                                            cfg.max_ttl, cfg.clock_skew))
        if not cfg.sig_required:
            log.warning("[%s] SIG_REQUIRED=false; unsigned commands are accepted (local development only, "
                        "asset-command.md section 6.3)", cfg.source)
        if cfg.fence is None or not cfg.fence.keep_in:
            log.warning("[%s] no keep-in zone; goto and rtl are refused (fail closed)", cfg.source)
        # robot state
        self.x = self.y = self.z = self.yaw = 0.0
        self.vx = self.vy = 0.0
        self.has_pos = False
        self.pos_at: datetime | None = None
        self.heard_at: datetime | None = None  # last message of any kind from the robot
        self.last_pos_pub: datetime | None = None
        self.battery_pct = -1
        self.voltage = -1.0
        self.mode = "idle"
        self.estopped = False
        self.goal: dict | None = None  # {"x","y","lat","lon","task_env","goal_id"}
        self.home = cfg.home
        self._move_lock = asyncio.Lock()
        self._task_ver: dict[str, int] = {}
        self._tasks: set = set()
        self.counters = {"tel": 0, "commands": 0, "refused": 0, "executed": 0, "failed": 0}

    # ---- subjects ----
    def position_subject(self) -> str:
        return f"tel.fleet.robot.{self.cfg.asset_id}.position"

    def status_subject(self) -> str:
        return f"tel.fleet.robot.{self.cfg.asset_id}.status"

    def command_subject(self) -> str:
        return f"cmd.fleet.robot.{self.cfg.asset_id}.>"

    def outcome_subject(self) -> str:
        return f"evt.robot.{self.cfg.asset_id}.command"

    def task_subject(self) -> str:
        return f"evt.robot.{self.cfg.asset_id}.task"

    # ---- helpers ----
    def _spawn(self, coro) -> None:
        t = asyncio.ensure_future(coro)
        self._tasks.add(t)
        t.add_done_callback(self._tasks.discard)

    async def _publish(self, subject: str, payload, msg_id: str = "") -> None:
        try:
            await self.bus.publish(subject, json.dumps(payload, separators=(",", ":")).encode(), msg_id)
        except Exception as e:  # never let a bus error kill the robot loop
            log.warning("[%s] publish %s: %s", self.cfg.source, subject, e)

    def online(self, now: datetime) -> bool:
        return self.heard_at is not None and (now - self.heard_at).total_seconds() <= self.cfg.link_timeout

    def latlon(self) -> tuple[float, float, float]:
        return self.cfg.frame.to_wgs84(self.x, self.y, self.z)

    # ---- robot input (called on the asyncio loop) ----
    def on_odom(self, o: msgs.Odometry) -> None:
        now = self.clock()
        self.heard_at = now
        self.vx, self.vy = o.twist.lx, o.twist.ly
        if self.cfg.use_odom_position:
            self._set_pose(o.pose, now)

    def on_pose(self, p: msgs.PoseStamped) -> None:
        now = self.clock()
        self.heard_at = now
        if not self.cfg.use_odom_position:
            self._set_pose(p.pose, now)

    def on_battery(self, b: msgs.BatteryState) -> None:
        self.heard_at = self.clock()
        if b.percentage is not None and math.isfinite(b.percentage):
            self.battery_pct = max(0, min(100, int(round(b.percentage * 100))))
        if b.voltage is not None and math.isfinite(b.voltage):
            self.voltage = b.voltage

    def _set_pose(self, pose: msgs.Pose, now: datetime) -> None:
        vals = (pose.x, pose.y, pose.z, pose.qx, pose.qy, pose.qz, pose.qw)
        if not all(math.isfinite(v) for v in vals):
            return
        self.x, self.y, self.z, self.yaw = pose.x, pose.y, pose.z, pose.yaw()
        self.has_pos, self.pos_at = True, now
        if self.home is None:
            lat, lon, _ = self.latlon()
            self.home = (lat, lon)
            log.info("[%s] home set to the first position heard: %.7f, %.7f", self.cfg.source, lat, lon)
        self._check_goal()
        if self.last_pos_pub is None or (now - self.last_pos_pub).total_seconds() >= self.cfg.position_interval:
            self.last_pos_pub = now
            self.counters["tel"] += 1
            self._spawn(self._publish(self.position_subject(), self.position_envelope(now)))

    def _check_goal(self) -> None:
        g = self.goal
        if g is None:
            return
        if math.hypot(g["x"] - self.x, g["y"] - self.y) <= self.cfg.goal_tolerance_m:
            log.info("[%s] goal reached (%s)", self.cfg.source, g["verb"])
            self.goal = None
            self.mode = "idle"
            env = g.get("env")
            if env and env.get("task_id"):
                self._spawn(self._report_task(env, "done_ok", None))

    # ---- telemetry ----
    def _speed(self) -> float:
        return math.hypot(self.vx, self.vy)

    def position_envelope(self, now: datetime) -> dict:
        lat, lon, hae = self.latlon()
        speed = self._speed()
        heading = self.cfg.frame.yaw_to_bearing(self.yaw)
        course = heading if speed >= MOVING_MPS else -1.0
        if speed >= MOVING_MPS and self.vx < 0:
            course = (heading + 180.0) % 360.0  # reversing
        return {
            "uid": self.cfg.asset_id, "type": "fleet.robot", "lat": round(lat, 8), "lon": round(lon, 8),
            "hae": _rnd(hae, 2), "course": _rnd(course, 1), "speed": _rnd(speed, 3),
            "callsign": self.cfg.callsign, "stale": self.cfg.stale_s, "t": int(now.timestamp() * 1000),
            "meta": {
                "kind": self.cfg.kind, "heading_deg": _rnd(heading, 1), "battery_pct": self.battery_pct,
                "mode": self.mode, "estopped": self.estopped,
                "map_x_m": _rnd(self.x, 3), "map_y_m": _rnd(self.y, 3), "frame": self.cfg.map_frame,
                "hae_source": "map_origin", "simulated": True, "source": self.cfg.source,
                "transport": self.transport.name,
            },
        }

    def status_envelope(self, now: datetime) -> dict | None:
        if not self.has_pos:
            return None  # section 9.2 requires lat/lon
        lat, lon, hae = self.latlon()
        meta = {
            "kind": self.cfg.kind, "battery_pct": self.battery_pct, "voltage_v": _rnd(self.voltage, 2),
            "mode": self.mode, "estopped": self.estopped,
            "connection": "online" if self.online(now) else "offline",
            "link_rssi": -1,
            "last_heard_t": int(self.heard_at.timestamp() * 1000) if self.heard_at else -1,
            "position_age_ms": int((now - self.pos_at).total_seconds() * 1000) if self.pos_at else -1,
            "transport": self.transport.name, "simulated": True, "source": self.cfg.source,
            "commands_accepted": ",".join(C.VERBS),
        }
        meta.update({f"transport_{k}": v for k, v in self.transport.stats().items()})
        if self.goal:
            meta["goal_lat"], meta["goal_lon"] = round(self.goal["lat"], 8), round(self.goal["lon"], 8)
        if self.home:
            meta["home_lat"], meta["home_lon"] = round(self.home[0], 8), round(self.home[1], 8)
        return {"uid": self.cfg.asset_id, "type": "fleet.robot", "lat": round(lat, 8), "lon": round(lon, 8),
                "hae": _rnd(hae, 2), "course": -1, "speed": _rnd(self._speed(), 3), "callsign": self.cfg.callsign,
                "stale": self.cfg.stale_s, "t": int(now.timestamp() * 1000), "meta": meta}

    async def publish_status(self) -> None:
        env = self.status_envelope(self.clock())
        if env is not None:
            await self._publish(self.status_subject(), env)

    async def run_status(self, stop: asyncio.Event) -> None:
        while not stop.is_set():
            await self.publish_status()
            try:
                await asyncio.wait_for(stop.wait(), self.cfg.status_interval)
            except asyncio.TimeoutError:
                pass

    # ---- outcomes ----
    async def _event(self, subject: str, typ: str, data: dict, env: dict | None, msg_id: str = "") -> None:
        eid = str(uuid.uuid4())
        out = {"v": 1, "id": eid, "ts": rfc3339ms(self.clock()), "type": typ, "source": self.cfg.source,
               "subject": subject, "data": data}
        if isinstance(env, dict):
            if C._UUID.match(str(env.get("id", ""))):
                out["causation"] = env["id"]
            if C._UUID.match(str(env.get("corr", ""))):
                out["corr"] = env["corr"]
            if isinstance(env.get("trace"), str) and env["trace"]:
                out["trace"] = env["trace"]
        await self._publish(subject, out, msg_id or eid)

    async def report(self, env: dict | None, verb: str, outcome: str, code: str = "", message: str = "",
                     detail: str = "", ros: dict | None = None, sig_ok: bool | None = None) -> None:
        data = {"action_id": "", "command_id": "", "outcome": outcome, "verb": verb, "simulated": True}
        if isinstance(env, dict):
            data["action_id"] = env.get("action_id", "") if isinstance(env.get("action_id"), str) else ""
            data["command_id"] = env.get("id", "") if isinstance(env.get("id"), str) else ""
            if isinstance(env.get("task_id"), str) and env["task_id"]:
                data["task_id"] = env["task_id"]
        for k, v in (("code", code), ("detail", detail), ("message", message)):
            if v:
                data[k] = v
        if ros:
            data["ros"] = ros
        if sig_ok is not None:
            data["sig_verified"] = sig_ok
        if outcome in self.counters:
            self.counters[outcome] += 1
        log.info("[%s] command %s %s outcome=%s code=%s %s", self.cfg.source, data["command_id"], verb, outcome,
                 code, message)
        await self._event(self.outcome_subject(), "evt.robot.command", data, env)
        if isinstance(env, dict) and data.get("task_id"):
            state, err = {"acknowledged": ("received", None), "executed": ("executing", None),
                          "refused": ("done_not_ok", {"code": "rejected", "message": f"{code}: {message}"}),
                          "failed": ("done_not_ok", {"code": "failed", "message": f"{code}: {message}"}),
                          }.get(outcome, (None, None))
            if state == "executing" and verb in ("hold", "estop", "estop_release"):
                state = "done_ok"
            if state:
                await self._report_task(env, state, err)

    def _next_task_version(self, task_id: str) -> int:
        if task_id not in self._task_ver and len(self._task_ver) >= MAX_TRACKED_TASKS:
            self._task_ver.pop(next(iter(self._task_ver)))
        v = self._task_ver.get(task_id, 0) + 1
        self._task_ver[task_id] = v
        return v

    async def _report_task(self, env: dict, state: str, err) -> None:
        tid = env["task_id"]
        ver = self._next_task_version(tid)
        await self._event(self.task_subject(), "evt.robot.task",
                          {"task_id": tid, "state": state, "status_version": ver, "error": err,
                           "at_t": int(self.clock().timestamp() * 1000), "simulated": True},
                          env, msg_id=f"{tid}:{ver}")

    # ---- commands ----
    def _precond(self, detail: str, message: str) -> C.Refusal:
        return C.Refusal(C.PRECONDITION, message, detail)

    def _precheck(self, res: C.Result, now: datetime) -> tuple[float, float] | None:
        """Robot-state checks. Returns the target (lat, lon) for movement verbs."""
        verb, p = res.verb.name, res.params
        if verb == "estop":
            return None  # a stop is never blocked (asset-command.md section 3)
        if not self.online(now):
            raise self._precond("offline", f"nothing heard from the robot within {self.cfg.link_timeout}s")
        if verb in ("goto", "rtl"):
            if self.estopped:
                raise self._precond("estopped", "robot is e-stopped; estop_release (two humans) first")
            if not self.has_pos:
                raise self._precond("keep_in", "robot position unknown")
            if verb == "goto":
                target = (p.lat, p.lon)
            else:
                if self.home is None:
                    raise self._precond("keep_in", "home position unknown")
                target = self.home
            v = fence_check(self.cfg.fence, *target)
            if v is not None:
                raise self._precond(v.detail, f"{verb} target: {v.reason}")
            return target
        return None

    def _goal_pose(self, lat: float, lon: float, bearing: float | None) -> tuple[msgs.PoseStamped, float, float]:
        gx, gy = self.cfg.frame.from_wgs84(lat, lon)
        if bearing is not None:
            yaw = self.cfg.frame.bearing_to_yaw(bearing)
        else:  # face the direction of travel
            yaw = math.atan2(gy - self.y, gx - self.x) if math.hypot(gx - self.x, gy - self.y) > 1e-3 else self.yaw
        return msgs.PoseStamped(msgs.Header(msgs.Time.now(), self.cfg.map_frame),
                                msgs.Pose.from_xy_yaw(gx, gy, yaw)), gx, gy

    async def handle_command(self, subject: str, data: bytes) -> None:
        self.counters["commands"] += 1
        prefix = f"cmd.fleet.robot.{self.cfg.asset_id}."
        verb = subject[len(prefix):] if subject.startswith(prefix) else subject
        try:
            env = C.decode(data)
        except C.Refusal as r:
            await self.report(None, verb, "refused", r.code, r.message)
            return
        now = self.clock()
        try:
            res = self.verifier.verify(subject, env, now)
        except C.Refusal as r:
            await self.report(env if isinstance(env, dict) else None, verb, "refused", r.code, r.message, r.detail)
            return
        if not res.signed:
            log.warning("[%s] command %s accepted without a signature (SIG_REQUIRED=false)", self.cfg.source, env["id"])
        sig_ok = res.signed
        try:
            target = self._precheck(res, now)
        except C.Refusal as r:
            await self.report(env, verb, "refused", r.code, r.message, r.detail, sig_ok=sig_ok)
            return
        if env.get("dry_run") is True:
            await self.report(env, verb, "acknowledged", message="dry_run: verified, not sent to the robot", sig_ok=sig_ok)
            return
        await self.report(env, verb, "acknowledged", sig_ok=sig_ok)

        if verb == "estop":
            await self._execute(env, res, target)
            return
        async with self._move_lock:
            if not self.clock() < res.expires_at:
                await self.report(env, verb, "failed", C.EXPIRED, "expired while waiting for the previous command")
                return
            await self._execute(env, res, target)

    async def _execute(self, env: dict, res: C.Result, target) -> None:
        verb, p, tp = res.verb.name, res.params, self.transport.name
        ros = {"transport": tp}
        try:
            if verb in ("goto", "rtl"):
                pose, gx, gy = self._goal_pose(target[0], target[1], p.yaw_deg if verb == "goto" else None)
                ros.update(interface=self.transport.topics.nav_action if tp == "zenoh" else self.transport.topics.goal_topic,
                           frame=self.cfg.map_frame, x=round(gx, 3), y=round(gy, 3))
                r = await self.transport.send_goal(pose)
                ros["goal_id"] = r.goal_id
                if not r.accepted:
                    await self.report(env, verb, "failed", C.VEHICLE_REJECTED, f"goal rejected ({r.detail})", ros=ros)
                    return
                self.goal = {"verb": verb, "x": gx, "y": gy, "lat": target[0], "lon": target[1], "env": env,
                             "goal_id": r.goal_id}
                self.mode = "navigating" if verb == "goto" else "rtl"
                msg = r.detail
                if verb == "goto" and p.speed_mps is not None:
                    msg += "; speed_mps not applied (NavigateToPose has no per-goal speed)"
                await self.report(env, verb, "executed", message=msg, ros=ros)
            elif verb == "hold":
                here = msgs.PoseStamped(msgs.Header(msgs.Time.now(), self.cfg.map_frame),
                                        msgs.Pose.from_xy_yaw(self.x, self.y, self.yaw))
                detail = await self.transport.hold(here)
                self.goal, self.mode = None, "holding"
                await self.report(env, verb, "executed", message=detail, ros=ros)
            elif verb == "estop":
                self.estopped, self.goal, self.mode = True, None, "estop"  # latch before anything can fail
                detail = await self.transport.estop()
                await self.report(env, verb, "executed", message=detail, ros=ros)
            elif verb == "estop_release":
                self.estopped = False
                self.mode = "idle"
                await self.report(env, verb, "executed", message="e-stop latch released; robot stays stopped", ros=ros)
            else:  # pragma: no cover - VERBS and this switch are kept in step by tests
                await self.report(env, verb, "failed", C.UNKNOWN_ACTION, f"no ROS mapping for {verb}")
        except TransportTimeout as e:
            await self.report(env, verb, "failed", C.NO_ACK, str(e), ros=ros)
        except TransportError as e:
            await self.report(env, verb, "failed", C.NO_ACK, str(e), detail="not_sent", ros=ros)
