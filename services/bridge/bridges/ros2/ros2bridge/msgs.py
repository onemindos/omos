"""ROS 2 message layouts used by the bridge, as plain CDR (see cdr.py).

Field order follows the ROS 2 interface definitions (Jazzy):
  builtin_interfaces/msg/Time          int32 sec, uint32 nanosec
  std_msgs/msg/Header                  Time stamp, string frame_id
  geometry_msgs/msg/Pose               Point position (3 x float64), Quaternion orientation (4 x float64)
  geometry_msgs/msg/PoseStamped        Header, Pose
  geometry_msgs/msg/PoseWithCovarianceStamped  Header, Pose, float64[36]
  geometry_msgs/msg/Twist              Vector3 linear, Vector3 angular (float64)
  nav_msgs/msg/Odometry                Header, string child_frame_id, Pose, float64[36], Twist, float64[36]
  sensor_msgs/msg/BatteryState         Header, float32 voltage, temperature, current, charge, capacity,
                                       design_capacity, percentage, uint8 x3, bool present,
                                       float32[] cell_voltage, float32[] cell_temperature, string location,
                                       string serial_number
  nav2_msgs/action/NavigateToPose      goal: PoseStamped pose, string behavior_tree
  action_msgs SendGoal request         unique_identifier_msgs/UUID goal_id (uint8[16]), goal
  action_msgs SendGoal response        bool accepted, Time stamp
  action_msgs/srv/CancelGoal           request GoalInfo{UUID, Time}; response int8 return_code, GoalInfo[]
"""

from __future__ import annotations

import math
import time as _systime
from dataclasses import dataclass, field

from .cdr import Reader, Writer


@dataclass
class Time:
    sec: int = 0
    nanosec: int = 0

    @staticmethod
    def now() -> "Time":
        t = _systime.time_ns()
        return Time(t // 1_000_000_000, t % 1_000_000_000)

    def to_seconds(self) -> float:
        return self.sec + self.nanosec / 1e9


@dataclass
class Header:
    stamp: Time = field(default_factory=Time)
    frame_id: str = ""


@dataclass
class Pose:
    x: float = 0.0
    y: float = 0.0
    z: float = 0.0
    qx: float = 0.0
    qy: float = 0.0
    qz: float = 0.0
    qw: float = 1.0

    def yaw(self) -> float:
        """Yaw (rad, counter-clockwise from +x) from the quaternion (REP-103)."""
        siny = 2.0 * (self.qw * self.qz + self.qx * self.qy)
        cosy = 1.0 - 2.0 * (self.qy * self.qy + self.qz * self.qz)
        return math.atan2(siny, cosy)

    @staticmethod
    def from_xy_yaw(x: float, y: float, yaw: float, z: float = 0.0) -> "Pose":
        return Pose(x, y, z, 0.0, 0.0, math.sin(yaw / 2), math.cos(yaw / 2))


@dataclass
class Twist:
    lx: float = 0.0
    ly: float = 0.0
    lz: float = 0.0
    ax: float = 0.0
    ay: float = 0.0
    az: float = 0.0


@dataclass
class PoseStamped:
    header: Header
    pose: Pose


@dataclass
class Odometry:
    header: Header
    child_frame_id: str
    pose: Pose
    twist: Twist


@dataclass
class BatteryState:
    header: Header
    voltage: float
    percentage: float  # 0..1, NaN when unknown
    present: bool


@dataclass
class SendGoalRequest:
    goal_id: bytes
    goal: PoseStamped
    behavior_tree: str = ""


@dataclass
class SendGoalResponse:
    accepted: bool
    stamp: Time


@dataclass
class CancelGoalResponse:
    return_code: int  # 0 none, 1 rejected, 2 unknown goal, 3 terminated
    goals: list = field(default_factory=list)  # list of (goal_id bytes, Time)


# ---- read helpers ----

def _time(r: Reader) -> Time:
    return Time(r.i32(), r.u32())


def _header(r: Reader) -> Header:
    return Header(_time(r), r.string())


def _pose(r: Reader) -> Pose:
    return Pose(*(r.f64() for _ in range(7)))


def _twist(r: Reader) -> Twist:
    return Twist(*(r.f64() for _ in range(6)))


# ---- write helpers ----

def _w_time(w: Writer, t: Time) -> None:
    w.i32(t.sec)
    w.u32(t.nanosec)


def _w_header(w: Writer, h: Header) -> None:
    _w_time(w, h.stamp)
    w.string(h.frame_id)


def _w_pose(w: Writer, p: Pose) -> None:
    for v in (p.x, p.y, p.z, p.qx, p.qy, p.qz, p.qw):
        w.f64(v)


def _w_twist(w: Writer, t: Twist) -> None:
    for v in (t.lx, t.ly, t.lz, t.ax, t.ay, t.az):
        w.f64(v)


# ---- public codecs ----

def decode_odometry(buf: bytes) -> Odometry:
    r = Reader(buf)
    h = _header(r)
    child = r.string()
    pose = _pose(r)
    r.f64_array(36)
    twist = _twist(r)
    r.f64_array(36)
    return Odometry(h, child, pose, twist)


def encode_odometry(o: Odometry) -> bytes:
    w = Writer()
    _w_header(w, o.header)
    w.string(o.child_frame_id)
    _w_pose(w, o.pose)
    w.f64_array([0.0] * 36)
    _w_twist(w, o.twist)
    w.f64_array([0.0] * 36)
    return w.bytes()


def decode_pose_with_covariance_stamped(buf: bytes) -> PoseStamped:
    r = Reader(buf)
    h = _header(r)
    pose = _pose(r)
    r.f64_array(36)
    return PoseStamped(h, pose)


def encode_pose_with_covariance_stamped(p: PoseStamped) -> bytes:
    w = Writer()
    _w_header(w, p.header)
    _w_pose(w, p.pose)
    w.f64_array([0.0] * 36)
    return w.bytes()


def decode_pose_stamped(buf: bytes) -> PoseStamped:
    r = Reader(buf)
    return PoseStamped(_header(r), _pose(r))


def encode_pose_stamped(p: PoseStamped) -> bytes:
    w = Writer()
    _w_header(w, p.header)
    _w_pose(w, p.pose)
    return w.bytes()


def decode_twist(buf: bytes) -> Twist:
    return _twist(Reader(buf))


def encode_twist(t: Twist) -> bytes:
    w = Writer()
    _w_twist(w, t)
    return w.bytes()


def decode_battery_state(buf: bytes) -> BatteryState:
    r = Reader(buf)
    h = _header(r)
    voltage = r.f32()
    r.f32()  # temperature
    r.f32()  # current
    r.f32()  # charge
    r.f32()  # capacity
    r.f32()  # design_capacity
    pct = r.f32()
    r.u8()  # power_supply_status
    r.u8()  # power_supply_health
    r.u8()  # power_supply_technology
    present = r.boolean()
    r.f32_seq()  # cell_voltage
    r.f32_seq()  # cell_temperature
    r.string()  # location
    r.string()  # serial_number
    return BatteryState(h, voltage, pct, present)


def encode_battery_state(b: BatteryState) -> bytes:
    w = Writer()
    _w_header(w, b.header)
    w.f32(b.voltage)
    for _ in range(5):
        w.f32(float("nan"))  # temperature, current, charge, capacity, design_capacity
    w.f32(b.percentage)
    w.u8(0)
    w.u8(0)
    w.u8(0)
    w.boolean(b.present)
    w.f32_seq([])
    w.f32_seq([])
    w.string("")
    w.string("")
    return w.bytes()


def encode_send_goal_request(req: SendGoalRequest) -> bytes:
    if len(req.goal_id) != 16:
        raise ValueError("goal_id must be 16 bytes")
    w = Writer()
    w.octets(req.goal_id)
    _w_header(w, req.goal.header)
    _w_pose(w, req.goal.pose)
    w.string(req.behavior_tree)
    return w.bytes()


def decode_send_goal_request(buf: bytes) -> SendGoalRequest:
    r = Reader(buf)
    gid = r.octets(16)
    goal = PoseStamped(_header(r), _pose(r))
    return SendGoalRequest(gid, goal, r.string())


def encode_send_goal_response(resp: SendGoalResponse) -> bytes:
    w = Writer()
    w.boolean(resp.accepted)
    _w_time(w, resp.stamp)
    return w.bytes()


def decode_send_goal_response(buf: bytes) -> SendGoalResponse:
    r = Reader(buf)
    return SendGoalResponse(r.boolean(), _time(r))


def encode_cancel_goal_request(goal_id: bytes = b"\x00" * 16, stamp: Time | None = None) -> bytes:
    """A zero goal_id and zero stamp cancel every goal (action_msgs/srv/CancelGoal)."""
    w = Writer()
    w.octets(goal_id)
    _w_time(w, stamp or Time())
    return w.bytes()


def decode_cancel_goal_request(buf: bytes) -> tuple[bytes, Time]:
    r = Reader(buf)
    return r.octets(16), _time(r)


def encode_cancel_goal_response(resp: CancelGoalResponse) -> bytes:
    w = Writer()
    w.i8(resp.return_code)
    w.u32(len(resp.goals))
    for gid, t in resp.goals:
        w.octets(gid)
        _w_time(w, t)
    return w.bytes()


def decode_cancel_goal_response(buf: bytes) -> CancelGoalResponse:
    r = Reader(buf)
    code = r.i8()
    n = r.seq_len()
    goals = []
    for _ in range(n):
        goals.append((r.octets(16), _time(r)))
    return CancelGoalResponse(code, goals)
