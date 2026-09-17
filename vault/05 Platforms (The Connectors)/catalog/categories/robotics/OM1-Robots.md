# OM1 (OpenMind Robots)

> **Open-source framework for embodied AI robotics**

## Overview

OM1 is an open-source robotics framework by OpenMind that enables robots to understand and interact with the physical world. Combines AI perception, reasoning, and control for embodied intelligence.

## Status

| Attribute   | Value                 |
| ----------- | --------------------- |
| Status      | 🔄 Experimental       |
| Category    | Robotics              |
| Type        | Embodied AI Framework |
| Integration | Python, ROS2          |

## Use Cases in OneMind

1. **Home Robots** - AI-controlled mobile robots
2. **Manipulation** - Pick and place tasks
3. **Navigation** - Autonomous movement
4. **Human Interaction** - Voice and gesture control
5. **Task Execution** - Real-world agent actions

## Key Features

- Multimodal perception (vision, audio, touch)
- LLM-based reasoning
- Motion planning
- ROS2 integration
- Simulation support
- Real-time control
- Open-source

## Architecture

```
┌─────────────────────────────────────────────────────────┐
│                      OM1 Stack                          │
├─────────────────────────────────────────────────────────┤
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────────┐  │
│  │  Perception │  │  Reasoning  │  │     Action      │  │
│  │  (Vision,   │──│  (LLM,      │──│  (Motion,       │  │
│  │   Audio)    │  │   Planning) │  │   Control)      │  │
│  └─────────────┘  └─────────────┘  └─────────────────┘  │
├─────────────────────────────────────────────────────────┤
│                   Hardware Abstraction                   │
├─────────────────────────────────────────────────────────┤
│  Robot Base  │  Manipulator  │  Sensors  │  Actuators   │
└─────────────────────────────────────────────────────────┘
```

## Installation

```bash
# Clone OM1 repository
git clone https://github.com/OpenmindAGI/OM1.git
cd OM1

# Create virtual environment
python -m venv venv
source venv/bin/activate

# Install dependencies
pip install -e .

# Install ROS2 packages (optional)
sudo apt install ros-humble-desktop
```

## Basic Usage

### Initialize Robot

```python
from om1 import Robot, Perception, Planner, Controller

# Create robot instance
robot = Robot(
    name="onemind-bot",
    config={
        "base_type": "differential",
        "has_arm": True,
        "cameras": ["rgb", "depth"],
        "microphone": True
    }
)

# Initialize subsystems
perception = Perception(robot)
planner = Planner(robot, model="claude-3-opus")
controller = Controller(robot)

await robot.start()
```

### Vision and Understanding

```python
async def describe_scene():
    """Get AI description of what robot sees."""
    image = await perception.capture_rgb()

    description = await planner.analyze_image(
        image,
        prompt="Describe what you see. Identify any objects I could pick up."
    )

    return description

async def find_object(object_name: str):
    """Find and locate an object in view."""
    image = await perception.capture_rgb()
    depth = await perception.capture_depth()

    # Use vision model to find object
    detection = await planner.find_object(
        image,
        object_name,
        return_bbox=True
    )

    if detection:
        # Convert to 3D coordinates
        x, y = detection["center"]
        z = depth[y, x]

        return {
            "found": True,
            "position": (x, y, z),
            "confidence": detection["confidence"]
        }

    return {"found": False}
```

### Movement and Navigation

```python
async def move_to_location(location: str):
    """Navigate to a named location."""
    # Use LLM to understand location
    target = await planner.resolve_location(location)

    # Plan path
    path = await planner.plan_path(
        start=robot.current_position,
        goal=target
    )

    # Execute movement
    await controller.follow_path(path)

    return f"Arrived at {location}"

async def approach_object(object_name: str):
    """Move towards a detected object."""
    detection = await find_object(object_name)

    if detection["found"]:
        await controller.approach(detection["position"])
        return True

    return False
```

### Manipulation

```python
async def pick_up_object(object_name: str):
    """Pick up an object with the robot arm."""
    # Find object
    detection = await find_object(object_name)

    if not detection["found"]:
        return f"Could not find {object_name}"

    # Plan grasp
    grasp = await planner.plan_grasp(
        object_position=detection["position"],
        object_type=object_name
    )

    # Execute pick
    await controller.execute_grasp(grasp)
    await controller.lift(height=0.1)  # 10cm lift

    return f"Picked up {object_name}"

async def place_object(location: str):
    """Place held object at location."""
    target = await planner.resolve_location(location)

    await controller.move_arm(target)
    await controller.open_gripper()

    return f"Placed object at {location}"
```

## LLM Integration

```python
from om1 import LLMPlanner

planner = LLMPlanner(
    model="claude-3-opus",
    system_prompt="""You are controlling a home robot.
    You can see through cameras and control movement/manipulation.
    Plan actions step by step to complete tasks safely."""
)

async def execute_task(task: str):
    """Use LLM to plan and execute complex task."""
    # Get current state
    image = await perception.capture_rgb()
    position = robot.current_position

    # Ask LLM for plan
    plan = await planner.plan_task(
        task=task,
        image=image,
        current_position=position,
        available_actions=controller.available_actions
    )

    # Execute each step
    results = []
    for step in plan["steps"]:
        result = await controller.execute(step)
        results.append(result)

        if not result["success"]:
            break

    return {
        "task": task,
        "plan": plan,
        "results": results
    }
```

## Legacy AI Integration

```python
from agno import Tool

@Tool
async def command_robot(
    command: str,
    robot_name: str = "onemind-bot"
) -> dict:
    """Send command to OM1 robot.

    Args:
        command: Natural language command
        robot_name: Which robot to command
    """
    robot = get_robot(robot_name)
    result = await robot.execute_task(command)

    return {
        "command": command,
        "success": result["results"][-1]["success"] if result["results"] else False,
        "steps_completed": len([r for r in result["results"] if r["success"]])
    }

@Tool
async def robot_look_for(
    item: str,
    robot_name: str = "onemind-bot"
) -> dict:
    """Have robot look for an item.

    Args:
        item: What to look for
        robot_name: Which robot
    """
    robot = get_robot(robot_name)
    result = await robot.find_object(item)

    return result

@Tool
async def robot_move_to(
    location: str,
    robot_name: str = "onemind-bot"
) -> str:
    """Move robot to location.

    Args:
        location: Where to go (e.g., "kitchen", "living room")
        robot_name: Which robot
    """
    robot = get_robot(robot_name)
    await robot.move_to_location(location)

    return f"Robot moved to {location}"

@Tool
async def robot_pick_and_place(
    item: str,
    destination: str,
    robot_name: str = "onemind-bot"
) -> str:
    """Have robot pick up item and place it somewhere.

    Args:
        item: What to pick up
        destination: Where to place it
        robot_name: Which robot
    """
    robot = get_robot(robot_name)

    await robot.pick_up_object(item)
    await robot.move_to_location(destination)
    await robot.place_object(destination)

    return f"Moved {item} to {destination}"
```

## ROS2 Integration

```python
import rclpy
from rclpy.node import Node
from geometry_msgs.msg import Twist
from sensor_msgs.msg import Image

class OM1Bridge(Node):
    def __init__(self):
        super().__init__('om1_bridge')

        self.cmd_vel_pub = self.create_publisher(
            Twist, '/cmd_vel', 10
        )

        self.image_sub = self.create_subscription(
            Image, '/camera/rgb/image_raw',
            self.image_callback, 10
        )

    def move(self, linear: float, angular: float):
        msg = Twist()
        msg.linear.x = linear
        msg.angular.z = angular
        self.cmd_vel_pub.publish(msg)

    def image_callback(self, msg):
        # Process incoming camera image
        pass
```

## Simulation

```python
from om1.sim import GazeboSimulator

# Start in simulation
sim = GazeboSimulator(
    world="home_environment",
    robot="turtlebot3_with_arm"
)

robot = Robot(
    name="sim-bot",
    simulator=sim
)

# Same API works in sim or real
await robot.move_to_location("kitchen")
```

## Environment Variables

```bash
OM1_ROBOT_IP=192.168.1.100
OM1_ROBOT_PORT=5000
OM1_MODEL=claude-3-opus
ANTHROPIC_API_KEY=your_key

# ROS2
ROS_DOMAIN_ID=42
RMW_IMPLEMENTATION=rmw_fastrtps_cpp
```

## Best Practices

1. **Safety first** - Always have emergency stop
2. **Start in simulation** - Test before real robot
3. **Incremental tasks** - Small steps, verify each
4. **Error recovery** - Handle failures gracefully
5. **Human oversight** - Approval for risky actions

## Resources

- [OpenMind OM1 GitHub](https://github.com/OpenmindAGI/OM1)
- [ROS2 Documentation](https://docs.ros.org/en/humble/)
- [Embodied AI Research](https://embodied-ai.org/)

---

Last Updated: January 25, 2026


---

## Related

- [[00-24 SO (Sovereign Operations)/02 Infrastructure (VPS, Servers & Networking) (The Body)/02 Infrastructure|02 Infrastructure]]
