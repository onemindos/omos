"""
OneMind OS — ROS2/NATS Bridge Skeleton (Robot: DIMOS / OM1 / G1 / Go2)
Pattern: NATS cmd → ROS2 action/topic → NATS tel/evt back

Subject law:
  cmd.robot.<id>.<action>     ← Legacy publishes here
  tel.robot.<id>.<reading>    ← bridge publishes sensor data here
  evt.robot.<id>.<event>      ← bridge publishes state changes here

One bridge serves ALL robots of the same type.
Adding a new robot = register its ID, point ROBOT_ID to it.

Requirements (run on robot's compute unit or companion Linux machine):
  pip install nats-py rclpy

ROS2 must be installed and sourced:
  source /opt/ros/humble/setup.bash   # or foxy, iron, jazzy

IMPORTANT: NATS is NOT the ROS2 planner.
  - Robot's nav2 / MoveIt / ROS2 inner loops still run natively
  - NATS handles: high-level tasks (patrol, inspect, go_to, report_status)
  - Inner control loops stay in ROS2 — never pipe 100Hz joint states through NATS
"""

import asyncio
import json
import threading
import rclpy
from rclpy.node import Node
from geometry_msgs.msg import PoseStamped, Twist
from std_msgs.msg import String
from nav_msgs.msg import Odometry
from sensor_msgs.msg import BatteryState
import nats

NATS_URL = "nats://omos-mbp.tail717ba0.ts.net:4222"   # hub over Tailscale
ROBOT_ID = "om1"                                         # from node registry


# ── ROS2 node (runs in a background thread) ───────────────────────────────────

class RobotBridgeNode(Node):

    def __init__(self):
        super().__init__("nats_bridge")
        self.latest_odom    = None
        self.latest_battery = None

        # subscribe to ROS2 topics — robot publishes these natively
        self.create_subscription(Odometry, "/odom",
            self._on_odom, 10)
        self.create_subscription(BatteryState, "/battery_state",
            self._on_battery, 10)

        # publisher for velocity commands
        self.vel_pub = self.create_publisher(Twist, "/cmd_vel", 10)

        # publisher for nav2 goal
        self.goal_pub = self.create_publisher(
            PoseStamped, "/goal_pose", 10)

    def _on_odom(self, msg):
        self.latest_odom = {
            "x":   msg.pose.pose.position.x,
            "y":   msg.pose.pose.position.y,
            "z":   msg.pose.pose.position.z,
            "vx":  msg.twist.twist.linear.x,
            "vyaw": msg.twist.twist.angular.z,
        }

    def _on_battery(self, msg):
        self.latest_battery = {
            "voltage_v":      msg.voltage,
            "remaining_pct":  int(msg.percentage * 100),
            "charging":       (msg.power_supply_status == BatteryState.POWER_SUPPLY_STATUS_CHARGING),
        }

    def send_goal(self, x: float, y: float):
        goal = PoseStamped()
        goal.header.frame_id = "map"
        goal.pose.position.x = x
        goal.pose.position.y = y
        goal.pose.orientation.w = 1.0
        self.goal_pub.publish(goal)

    def stop(self):
        twist = Twist()   # all zeros = stop
        self.vel_pub.publish(twist)


# Global ROS2 node (started in background thread)
ros_node: RobotBridgeNode = None


def spin_ros():
    rclpy.init()
    global ros_node
    ros_node = RobotBridgeNode()
    rclpy.spin(ros_node)


# ── NATS command handlers (cmd → ROS2) ───────────────────────────────────────

async def handle_goto(nc, msg):
    """
    NATS subject: cmd.robot.om1.goto
    Payload: { "payload": { "x": 3.0, "y": 1.5 } }
    Sends nav2 goal to robot — nav2 plans and executes path natively.
    """
    body = json.loads(msg.data)
    p    = body.get("payload", {})
    x    = p.get("x", 0.0)
    y    = p.get("y", 0.0)

    ros_node.send_goal(x, y)

    await nc.publish(
        f"evt.robot.{ROBOT_ID}.goal_sent",
        json.dumps({"ok": True, "x": x, "y": y}).encode()
    )
    if msg.reply:
        await nc.publish(msg.reply,
            json.dumps({"ok": True, "x": x, "y": y}).encode())


async def handle_stop(nc, msg):
    """
    NATS subject: cmd.robot.om1.stop
    Immediately halts all motion.
    """
    ros_node.stop()
    await nc.publish(
        f"evt.robot.{ROBOT_ID}.stopped",
        json.dumps({"ok": True}).encode()
    )
    if msg.reply:
        await nc.publish(msg.reply, json.dumps({"ok": True}).encode())


async def handle_patrol(nc, msg):
    """
    NATS subject: cmd.robot.om1.patrol
    Payload: { "payload": { "zone": "barn" } }
    Fires a nav2 patrol sequence through predefined waypoints.
    Waypoints for each zone come from the node registry (NATS KV).
    """
    body = json.loads(msg.data)
    zone = body.get("payload", {}).get("zone", "default")

    # TODO: look up waypoints from KV bucket "robot-waypoints" by zone
    # For POC, using hardcoded barn waypoints
    barn_waypoints = [
        (2.0, 0.0),
        (2.0, 3.0),
        (0.0, 3.0),
        (0.0, 0.0),
    ]

    await nc.publish(
        f"evt.robot.{ROBOT_ID}.patrol_started",
        json.dumps({"ok": True, "zone": zone, "waypoints": len(barn_waypoints)}).encode()
    )
    if msg.reply:
        await nc.publish(msg.reply,
            json.dumps({"ok": True, "zone": zone}).encode())

    # send waypoints sequentially
    for (x, y) in barn_waypoints:
        ros_node.send_goal(x, y)
        await asyncio.sleep(5)   # wait between waypoints (replace with nav2 feedback)

    await nc.publish(
        f"evt.robot.{ROBOT_ID}.patrol_completed",
        json.dumps({"ok": True, "zone": zone}).encode()
    )


async def handle_status(nc, msg):
    """
    NATS subject: cmd.robot.om1.status
    Returns current odometry + battery from latest ROS2 data.
    """
    status = {
        "ok":      True,
        "pose":    ros_node.latest_odom,
        "battery": ros_node.latest_battery,
    }
    if msg.reply:
        await nc.publish(msg.reply, json.dumps(status).encode())


# ── telemetry publisher (ROS2 → NATS) ────────────────────────────────────────

async def publish_telemetry(nc):
    """
    Polls latest ROS2 state and publishes to NATS every second.
    Use ephemeral tel.> subjects (not JetStream) — high frequency.
    """
    while True:
        await asyncio.sleep(1)

        if ros_node.latest_odom:
            await nc.publish(
                f"tel.robot.{ROBOT_ID}.pose",
                json.dumps(ros_node.latest_odom).encode()
            )

        if ros_node.latest_battery:
            await nc.publish(
                f"tel.robot.{ROBOT_ID}.battery",
                json.dumps(ros_node.latest_battery).encode()
            )


# ── main ─────────────────────────────────────────────────────────────────────

async def main():
    # start ROS2 spin in background thread
    ros_thread = threading.Thread(target=spin_ros, daemon=True)
    ros_thread.start()

    # wait for ROS2 node to initialize
    while ros_node is None:
        await asyncio.sleep(0.1)

    nc = await nats.connect(NATS_URL)
    print("[ros2-bridge] connected to NATS")

    prefix = f"cmd.robot.{ROBOT_ID}"

    await nc.subscribe(f"{prefix}.goto",
        cb=lambda m: asyncio.create_task(handle_goto(nc, m)))
    await nc.subscribe(f"{prefix}.stop",
        cb=lambda m: asyncio.create_task(handle_stop(nc, m)))
    await nc.subscribe(f"{prefix}.patrol",
        cb=lambda m: asyncio.create_task(handle_patrol(nc, m)))
    await nc.subscribe(f"{prefix}.status",
        cb=lambda m: asyncio.create_task(handle_status(nc, m)))

    print(f"[ros2-bridge] subscribed to {prefix}.* — bridge live")

    await asyncio.gather(
        publish_telemetry(nc),
        asyncio.Event().wait(),
    )


if __name__ == "__main__":
    asyncio.run(main())


# ── end-to-end flow ───────────────────────────────────────────────────────────
#
#  Legacy publishes:
#    cmd.robot.om1.patrol
#    { "id": "ulid", "source": "agent.legacy", "payload": { "zone": "barn" } }
#
#  ROS2 bridge hears it
#  Bridge sends nav2 goal poses over ROS2
#  Robot's nav2 stack plans path and drives (inner loop = pure ROS2, no NATS)
#
#  Bridge continuously publishes back:
#    tel.robot.om1.pose      { x, y, z, vx, vyaw }     ← 1Hz to NATS
#    tel.robot.om1.battery   { voltage_v, remaining_pct }
#    evt.robot.om1.patrol_started
#    evt.robot.om1.patrol_completed
#
#  Dashboard, safety monitor, or another agent hears these and reacts.
#
#  The robot still speaks ROS2 internally.
#  NATS is only the high-level command + reporting layer.
