"""
OneMind Robot Gateway — NATS ↔ ROS2 bridge
Runs on the robot's companion computer (or simulation host).
Sim/live mode is controlled by the SIM env var — no code changes needed.

NATS subjects:
  Publishes: tel.robot.{id}.position  tel.drone.{id}.position
             tel.robot.{id}.status    ent.robot.{id}.state
  Subscribes: cmd.robot.{id}.move    cmd.robot.{id}.stop
              cmd.drone.{id}.goto    cmd.drone.{id}.land
"""

import asyncio
import json
import math
import os
import time
import threading

import nats
import rclpy
from rclpy.node import Node
from geometry_msgs.msg import Twist
from nav_msgs.msg import Odometry
from sensor_msgs.msg import NavSatFix, BatteryState
from std_msgs.msg import String

ROBOT_ID   = os.environ["ROBOT_ID"]                          # e.g. "ugv-01"
ROBOT_TYPE = os.environ.get("ROBOT_TYPE", "robot")           # robot | drone
NATS_URL   = os.environ.get("NATS_URL", "nats://localhost:4222")
NATS_USER  = os.environ.get("NATS_USER", "")
NATS_PASS  = os.environ.get("NATS_PASS", "")
SIM        = os.environ.get("SIM", "false").lower() == "true"

TEL_POSITION = f"tel.{ROBOT_TYPE}.{ROBOT_ID}.position"
TEL_STATUS   = f"tel.{ROBOT_TYPE}.{ROBOT_ID}.status"
ENT_STATE    = f"ent.{ROBOT_TYPE}.{ROBOT_ID}.state"
CMD_MOVE     = f"cmd.{ROBOT_TYPE}.{ROBOT_ID}.move"
CMD_STOP     = f"cmd.{ROBOT_TYPE}.{ROBOT_ID}.stop"
CMD_GOTO     = f"cmd.drone.{ROBOT_ID}.goto"
CMD_LAND     = f"cmd.drone.{ROBOT_ID}.land"

_nc: nats.aio.client.Client | None = None
_cmd_vel_pub  = None
_loop: asyncio.AbstractEventLoop | None = None


class RobotBridgeNode(Node):
    def __init__(self):
        super().__init__("nats_ros2_bridge")
        self.last_gps: dict | None = None

        self.create_subscription(NavSatFix, "/fix", self._on_gps, 10)
        self.create_subscription(Odometry, "/odom", self._on_odom, 10)
        self.create_subscription(BatteryState, "/battery_state", self._on_battery, 10)

        global _cmd_vel_pub
        _cmd_vel_pub = self.create_publisher(Twist, "/cmd_vel", 10)

        self.get_logger().info(
            f"bridge ready | id={ROBOT_ID} type={ROBOT_TYPE} sim={SIM}"
        )

    def _on_gps(self, msg: NavSatFix):
        self.last_gps = {
            "lat": msg.latitude,
            "lon": msg.longitude,
            "alt": msg.altitude,
        }
        payload = {
            "id":   ROBOT_ID,
            "type": ROBOT_TYPE,
            "sim":  SIM,
            "ts":   time.time(),
            **self.last_gps,
        }
        _publish_async(TEL_POSITION, payload)
        _publish_async(ENT_STATE, {**payload, "callsign": ROBOT_ID})

    def _on_odom(self, msg: Odometry):
        # Publish odometry-derived heading when no GPS fix
        q = msg.pose.pose.orientation
        heading = math.degrees(math.atan2(
            2 * (q.w * q.z + q.x * q.y),
            1 - 2 * (q.y ** 2 + q.z ** 2),
        ))
        _publish_async(TEL_STATUS, {
            "id":      ROBOT_ID,
            "heading": heading,
            "vx":      msg.twist.twist.linear.x,
            "vy":      msg.twist.twist.linear.y,
            "ts":      time.time(),
        })

    def _on_battery(self, msg: BatteryState):
        _publish_async(TEL_STATUS, {
            "id":      ROBOT_ID,
            "battery": round(msg.percentage * 100, 1),
            "voltage": round(msg.voltage, 2),
            "ts":      time.time(),
        })


def _publish_async(subject: str, payload: dict):
    if _nc and _loop:
        asyncio.run_coroutine_threadsafe(
            _nc.publish(subject, json.dumps(payload).encode()),
            _loop,
        )


async def _handle_move(msg):
    try:
        data = json.loads(msg.data.decode())
        twist = Twist()
        twist.linear.x  = float(data.get("vx", 0.5))
        twist.angular.z = float(data.get("wz", 0.0))
        _cmd_vel_pub.publish(twist)
    except Exception as e:
        print(f"[bridge] cmd.move error: {e}")


async def _handle_stop(msg):
    _cmd_vel_pub.publish(Twist())  # zero velocity


async def _handle_goto(msg):
    # For UAVs: publish nav goal or MAVLink waypoint depending on stack
    try:
        data = json.loads(msg.data.decode())
        print(f"[bridge] goto: lat={data.get('lat')} lon={data.get('lon')} alt={data.get('alt')}")
        # TODO: publish to /mavros/setpoint_position/global for PX4/ArduPilot
    except Exception as e:
        print(f"[bridge] cmd.goto error: {e}")


async def _nats_loop():
    global _nc
    opts = {"servers": [NATS_URL]}
    if NATS_USER:
        opts["user"] = NATS_USER
        opts["password"] = NATS_PASS

    _nc = await nats.connect(**opts)
    print(f"[bridge] connected to NATS at {NATS_URL}")

    await _nc.subscribe(CMD_MOVE, cb=_handle_move)
    await _nc.subscribe(CMD_STOP, cb=_handle_stop)
    if ROBOT_TYPE == "drone":
        await _nc.subscribe(CMD_GOTO, cb=_handle_goto)
        await _nc.subscribe(CMD_LAND, cb=_handle_stop)

    print(f"[bridge] subscribed | move={CMD_MOVE} stop={CMD_STOP}")
    await asyncio.Event().wait()  # run forever


def _ros_thread():
    rclpy.init()
    node = RobotBridgeNode()
    rclpy.spin(node)
    rclpy.shutdown()


async def main():
    global _loop
    _loop = asyncio.get_running_loop()

    t = threading.Thread(target=_ros_thread, daemon=True)
    t.start()

    await _nats_loop()


if __name__ == "__main__":
    asyncio.run(main())
