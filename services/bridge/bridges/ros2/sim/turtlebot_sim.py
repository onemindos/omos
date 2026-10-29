#!/usr/bin/env python3
"""Scripted TurtleBot-like robot on Zenoh, standing in for ROS 2 + zenoh-bridge-ros2dds.

ROS 2 and Gazebo are not installed on the build machine, so this script emits what the
robot side looks like on Zenoh after zenoh-bridge-ros2dds (see ros2bridge/zenoh_transport.py):

  {ns}/odom                                 nav_msgs/msg/Odometry, CDR, 20 Hz (frame odom -> base_footprint)
  {ns}/battery_state                        sensor_msgs/msg/BatteryState, CDR, 1 Hz
  {ns}/cmd_vel                    (sub)     geometry_msgs/msg/Twist, CDR; overrides navigation for 0.5 s
  {ns}/navigate_to_pose/_action/send_goal   (queryable) NavigateToPose goal -> accepted
  {ns}/navigate_to_pose/_action/cancel_goal (queryable) cancel all

Kinematics are a unicycle limited to the TurtleBot3 Burger's 0.22 m/s and 2.84 rad/s.
The odom frame is treated as the map frame (AMCL not simulated). It is a test fixture,
not a robot: only sim-* ids are ever bridged.
"""

from __future__ import annotations

import argparse
import json
import math
import os
import sys
import threading
import time

sys.path.insert(0, os.path.join(os.path.dirname(os.path.abspath(__file__)), ".."))

from ros2bridge import msgs  # noqa: E402
from ros2bridge.zenoh_transport import key  # noqa: E402

MAX_V = 0.22
MAX_W = 2.84


def _wrap(a: float) -> float:
    return (a + math.pi) % (2 * math.pi) - math.pi


class TurtlebotSim:
    def __init__(self, namespace: str = "", x: float = 0.0, y: float = 0.0, yaw: float = 0.0,
                 arena_m: float = 50.0, battery: float = 0.95, log=None):
        self.ns, self.arena_m = namespace, arena_m
        self.x, self.y, self.yaw = x, y, yaw
        self.v = self.w = 0.0
        self.goal = None  # (x, y, goal_id)
        self.manual_until = 0.0
        self.manual = msgs.Twist()
        self.battery = battery
        self.lock = threading.Lock()
        self.log = log
        self.goals_received = 0
        self.cancels = 0
        self.twists = 0

    def _emit(self, **kv):
        if self.log:
            kv["t"] = round(time.time(), 3)
            self.log.write(json.dumps(kv) + "\n")
            self.log.flush()

    # ---- Zenoh handlers ----
    def on_send_goal(self, query):
        try:
            req = msgs.decode_send_goal_request(bytes(query.payload.to_bytes()))
        except Exception as e:  # malformed request: answer rejected
            self._emit(ev="goal_undecodable", err=str(e))
            query.reply(query.key_expr, msgs.encode_send_goal_response(msgs.SendGoalResponse(False, msgs.Time.now())))
            return
        gx, gy = req.goal.pose.x, req.goal.pose.y
        ok = math.hypot(gx, gy) <= self.arena_m and req.goal.header.frame_id in ("map", "odom")
        with self.lock:
            self.goals_received += 1
            if ok:
                self.goal = (gx, gy, req.goal_id)
        self._emit(ev="goal", accepted=ok, x=round(gx, 3), y=round(gy, 3), frame=req.goal.header.frame_id,
                   goal_id=req.goal_id.hex())
        query.reply(query.key_expr, msgs.encode_send_goal_response(msgs.SendGoalResponse(ok, msgs.Time.now())))

    def on_cancel(self, query):
        with self.lock:
            canceled = [(self.goal[2], msgs.Time.now())] if self.goal else []
            self.goal = None
            self.cancels += 1
        self._emit(ev="cancel", canceled=len(canceled))
        query.reply(query.key_expr, msgs.encode_cancel_goal_response(msgs.CancelGoalResponse(0, canceled)))

    def on_cmd_vel(self, sample):
        try:
            t = msgs.decode_twist(bytes(sample.payload.to_bytes()))
        except Exception:
            return
        with self.lock:
            self.manual, self.manual_until = t, time.monotonic() + 0.5
            self.twists += 1
            if t.lx == 0 and t.az == 0:
                self.goal = None  # a zero twist (e-stop) also drops the goal in this sim
        self._emit(ev="cmd_vel", lx=t.lx, az=t.az)

    # ---- physics ----
    def step(self, dt: float) -> None:
        with self.lock:
            if time.monotonic() < self.manual_until:
                v = max(-MAX_V, min(MAX_V, self.manual.lx))
                w = max(-MAX_W, min(MAX_W, self.manual.az))
            elif self.goal is not None:
                dx, dy = self.goal[0] - self.x, self.goal[1] - self.y
                dist = math.hypot(dx, dy)
                if dist < 0.05:
                    v = w = 0.0
                    self._emit(ev="arrived", x=round(self.x, 3), y=round(self.y, 3))
                    self.goal = None
                else:
                    e = _wrap(math.atan2(dy, dx) - self.yaw)
                    w = max(-MAX_W, min(MAX_W, 3.0 * e))
                    v = min(MAX_V, 0.8 * dist) * max(0.0, math.cos(e)) if abs(e) < 1.0 else 0.0
            else:
                v = w = 0.0
            self.v, self.w = v, w
            self.yaw = _wrap(self.yaw + w * dt)
            self.x += v * math.cos(self.yaw) * dt
            self.y += v * math.sin(self.yaw) * dt
            self.battery = max(0.0, self.battery - (0.00002 + 0.0001 * abs(v)) * dt)

    def odom(self) -> bytes:
        with self.lock:
            o = msgs.Odometry(msgs.Header(msgs.Time.now(), "odom"), "base_footprint",
                              msgs.Pose.from_xy_yaw(self.x, self.y, self.yaw), msgs.Twist(lx=self.v, az=self.w))
        return msgs.encode_odometry(o)

    def battery_msg(self) -> bytes:
        with self.lock:
            b = msgs.BatteryState(msgs.Header(msgs.Time.now(), "base_link"), 11.1 + self.battery, self.battery, True)
        return msgs.encode_battery_state(b)


def open_session(mode: str, listen: list, connect: list):
    import zenoh

    cfg = zenoh.Config()
    cfg.insert_json5("mode", json.dumps(mode))
    if listen:
        cfg.insert_json5("listen/endpoints", json.dumps(listen))
    if connect:
        cfg.insert_json5("connect/endpoints", json.dumps(connect))
    cfg.insert_json5("scouting/multicast/enabled", "false")
    return zenoh.open(cfg)


def serve(sim: TurtlebotSim, session, stop: threading.Event, rate_hz: float = 20.0) -> None:
    ns = sim.ns
    decls = [
        session.declare_queryable(key(ns, "navigate_to_pose/_action/send_goal"), sim.on_send_goal),
        session.declare_queryable(key(ns, "navigate_to_pose/_action/cancel_goal"), sim.on_cancel),
        session.declare_subscriber(key(ns, "cmd_vel"), sim.on_cmd_vel),
    ]
    odom_pub = session.declare_publisher(key(ns, "odom"))
    bat_pub = session.declare_publisher(key(ns, "battery_state"))
    dt = 1.0 / rate_hz
    n = 0
    next_t = time.monotonic()
    while not stop.is_set():
        sim.step(dt)
        odom_pub.put(sim.odom())
        if n % int(rate_hz) == 0:
            bat_pub.put(sim.battery_msg())
        n += 1
        next_t += dt
        stop.wait(max(0.0, next_t - time.monotonic()))
    for d in decls:
        d.undeclare()


def main() -> int:
    ap = argparse.ArgumentParser(description=__doc__.splitlines()[0])
    ap.add_argument("--mode", default="peer", choices=["peer", "router", "client"])
    ap.add_argument("--listen", default="tcp/127.0.0.1:7447", help="comma-separated zenoh endpoints")
    ap.add_argument("--connect", default="", help="comma-separated zenoh endpoints")
    ap.add_argument("--namespace", default="", help="key prefix, as zenoh-bridge-ros2dds 'namespace'")
    ap.add_argument("--x", type=float, default=0.0)
    ap.add_argument("--y", type=float, default=0.0)
    ap.add_argument("--yaw", type=float, default=0.0, help="initial yaw, rad")
    ap.add_argument("--arena-m", type=float, default=50.0, help="goals farther than this from the origin are rejected")
    ap.add_argument("--log", default="", help="JSON-lines log of goals, cancels and cmd_vel")
    args = ap.parse_args()
    split = lambda s: [x for x in s.split(",") if x]  # noqa: E731
    log = open(args.log, "a") if args.log else None
    sim = TurtlebotSim(args.namespace, args.x, args.y, args.yaw, args.arena_m, log=log)
    s = open_session(args.mode, split(args.listen), split(args.connect))
    stop = threading.Event()
    import signal

    signal.signal(signal.SIGTERM, lambda *_: stop.set())
    signal.signal(signal.SIGINT, lambda *_: stop.set())
    print(f"turtlebot sim on zenoh {args.mode} listen={args.listen} ns={args.namespace!r}", flush=True)
    try:
        serve(sim, s, stop)
    finally:
        s.close()
    return 0


if __name__ == "__main__":
    sys.exit(main())
