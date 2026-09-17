# Unitree

> **Quadruped robots for research, inspection, and automation**

## Overview

Unitree Robotics develops high-performance quadruped (four-legged) robots. Their robot dogs combine agility, durability, and AI capabilities for navigation, inspection, and embodied AI applications.

## Status

| Attribute   | Value           |
| ----------- | --------------- |
| Status      | 🔄 Planned      |
| Category    | Robotics        |
| Type        | Quadruped Robot |
| Integration | SDK, ROS2       |

## Use Cases in OneMind

1. **Property Patrol** - Autonomous security rounds
2. **Inspection** - Check areas, read sensors
3. **Companion** - Interactive robot pet
4. **Delivery** - Carry items around property
5. **Embodied AI** - Physical agent presence

## Robot Models

| Model | Type       | Use Case                  |
| ----- | ---------- | ------------------------- |
| Go2   | Consumer   | Personal, education, R&D  |
| B2    | Industrial | Inspection, patrol        |
| H1    | Humanoid   | Research, general purpose |
| G1    | Humanoid   | Compact humanoid          |

## Key Features

- Dynamic locomotion (walk, run, climb)
- Obstacle avoidance
- SLAM navigation
- Voice control
- Camera & LiDAR
- SDK for custom development
- ROS2 support

## Installation

### Unitree SDK

```bash
# Clone Unitree SDK
git clone https://github.com/unitreerobotics/unitree_sdk2.git
cd unitree_sdk2

# Build
mkdir build && cd build
cmake ..
make -j4

# Python bindings
pip install unitree-sdk2-python
```

### ROS2 Package

```bash
# In ROS2 workspace
cd ~/ros2_ws/src
git clone https://github.com/unitreerobotics/unitree_ros2.git

cd ~/ros2_ws
colcon build --packages-select unitree_ros2
source install/setup.bash
```

## Python SDK Integration

### Basic Connection

```python
from unitree_sdk2py.core.channel import ChannelFactory
from unitree_sdk2py.go2.sport.sport_client import SportClient

# Connect to robot
channel = ChannelFactory.create_sender("rt/go2/cmd")
sport_client = SportClient()
sport_client.init()

# Get robot state
state = sport_client.get_motion_state()
print(f"Position: {state.position}")
print(f"Velocity: {state.velocity}")
```

### Movement Control

```python
from unitree_sdk2py.go2.sport.sport_client import SportClient
import time

sport_client = SportClient()
sport_client.init()

async def move_forward(distance: float, speed: float = 0.5):
    """Move robot forward by distance in meters."""
    duration = distance / speed

    sport_client.move(speed, 0, 0)  # vx, vy, vyaw
    await asyncio.sleep(duration)
    sport_client.stop_move()

async def turn(degrees: float, speed: float = 0.5):
    """Turn robot by degrees."""
    import math
    radians = math.radians(degrees)
    duration = abs(radians) / speed

    direction = 1 if degrees > 0 else -1
    sport_client.move(0, 0, direction * speed)
    await asyncio.sleep(duration)
    sport_client.stop_move()

async def stand_up():
    """Make robot stand up."""
    sport_client.recovery_stand()
    await asyncio.sleep(2)

async def sit_down():
    """Make robot sit/lay down."""
    sport_client.stand_down()
    await asyncio.sleep(1)
```

### Poses and Actions

```python
async def do_action(action_name: str):
    """Execute predefined action."""
    actions = {
        "sit": lambda: sport_client.stand_down(),
        "stand": lambda: sport_client.recovery_stand(),
        "stretch": lambda: sport_client.stretch(),
        "wave": lambda: sport_client.hello(),
        "dance": lambda: sport_client.dance1(),
        "jump": lambda: sport_client.jump_yaw(0),
    }

    if action_name in actions:
        actions[action_name]()
        await asyncio.sleep(3)
        return f"Executed: {action_name}"

    return f"Unknown action: {action_name}"
```

### Camera Access

```python
from unitree_sdk2py.go2.video.video_client import VideoClient
import cv2
import numpy as np

video_client = VideoClient()
video_client.init()

async def capture_image() -> np.ndarray:
    """Capture image from robot camera."""
    frame = video_client.get_image_sample()
    image = np.frombuffer(frame.data, dtype=np.uint8)
    image = cv2.imdecode(image, cv2.IMREAD_COLOR)
    return image

async def stream_video():
    """Stream video from robot."""
    while True:
        frame = await capture_image()
        cv2.imshow("Unitree Camera", frame)
        if cv2.waitKey(1) & 0xFF == ord('q'):
            break
    cv2.destroyAllWindows()
```

### Navigation

```python
from unitree_sdk2py.go2.obstacles.obstacles_avoid_client import ObstaclesAvoidClient

obstacles_client = ObstaclesAvoidClient()
obstacles_client.init()

async def navigate_to(x: float, y: float):
    """Navigate to coordinates with obstacle avoidance."""
    # Enable obstacle avoidance
    obstacles_client.switch_set(True)

    # Set target position
    sport_client.move_to_position(x, y, 0)

    # Wait for arrival
    while True:
        state = sport_client.get_motion_state()
        dist = ((state.position[0] - x)**2 + (state.position[1] - y)**2)**0.5
        if dist < 0.1:  # Within 10cm
            break
        await asyncio.sleep(0.5)

    sport_client.stop_move()
    return f"Arrived at ({x}, {y})"
```

## Legacy AI Integration

```python
from agno import Tool

@Tool
async def unitree_move(
    direction: str,
    distance: float = 1.0
) -> str:
    """Move the Unitree robot dog.

    Args:
        direction: Direction (forward, backward, left, right)
        distance: Distance in meters
    """
    if direction == "forward":
        await move_forward(distance)
    elif direction == "backward":
        await move_forward(-distance)
    elif direction == "left":
        await turn(90)
        await move_forward(distance)
    elif direction == "right":
        await turn(-90)
        await move_forward(distance)

    return f"Moved {direction} {distance}m"

@Tool
async def unitree_action(action: str) -> str:
    """Make robot perform an action.

    Args:
        action: Action name (sit, stand, wave, dance, stretch)
    """
    return await do_action(action)

@Tool
async def unitree_patrol(waypoints: list[tuple]) -> str:
    """Send robot on patrol route.

    Args:
        waypoints: List of (x, y) coordinates to visit
    """
    await stand_up()

    for i, (x, y) in enumerate(waypoints):
        await navigate_to(x, y)
        # Capture image at each point
        image = await capture_image()
        # Analyze with vision
        await analyze_patrol_point(image, i)

    return f"Patrol complete: {len(waypoints)} points visited"

@Tool
async def unitree_look_around() -> dict:
    """Have robot look around and describe surroundings."""
    await stand_up()

    images = []
    for angle in [0, 90, 180, 270]:
        await turn(90)
        image = await capture_image()
        images.append(image)

    # Analyze images with LLM Vision
    descriptions = []
    for i, img in enumerate(images):
        desc = await analyze_image(img, "Describe what you see")
        descriptions.append(f"Direction {i*90}°: {desc}")

    return {
        "scan_complete": True,
        "descriptions": descriptions
    }

@Tool
async def unitree_follow_person() -> str:
    """Enable person following mode."""
    # Use built-in tracking
    sport_client.follow_mode(True)
    return "Following mode enabled"
```

## ROS2 Integration

```python
import rclpy
from rclpy.node import Node
from geometry_msgs.msg import Twist
from sensor_msgs.msg import Image, Imu

class UnitreeController(Node):
    def __init__(self):
        super().__init__('unitree_controller')

        # Publishers
        self.cmd_pub = self.create_publisher(
            Twist, '/cmd_vel', 10
        )

        # Subscribers
        self.imu_sub = self.create_subscription(
            Imu, '/imu/data',
            self.imu_callback, 10
        )

        self.camera_sub = self.create_subscription(
            Image, '/camera/image_raw',
            self.camera_callback, 10
        )

    def move(self, linear_x: float, angular_z: float):
        msg = Twist()
        msg.linear.x = linear_x
        msg.angular.z = angular_z
        self.cmd_pub.publish(msg)

    def imu_callback(self, msg):
        # Process IMU data
        orientation = msg.orientation
        self.get_logger().info(f"Orientation: {orientation}")

    def camera_callback(self, msg):
        # Process camera image
        pass
```

## Home Assistant Integration

```yaml
# configuration.yaml
rest_command:
  unitree_move:
    url: "http://unitree-api.local:8080/move"
    method: POST
    payload: '{"direction": "{{ direction }}", "distance": {{ distance }}}'
    content_type: "application/json"

  unitree_action:
    url: "http://unitree-api.local:8080/action"
    method: POST
    payload: '{"action": "{{ action }}"}'
    content_type: "application/json"

automation:
  - alias: "Unitree Patrol on Motion"
    trigger:
      - platform: state
        entity_id: binary_sensor.front_door_motion
        to: "on"
    action:
      - service: rest_command.unitree_action
        data:
          action: "patrol"
```

## Docker Deployment

```yaml
version: "3.8"

services:
  unitree-api:
    build: ./unitree-controller
    container_name: unitree-api
    network_mode: host # Required for UDP communication
    environment:
      - ROBOT_IP=192.168.123.161
      - ROBOT_TYPE=go2
    volumes:
      - ./config:/app/config
    restart: unless-stopped
```

## Environment Variables

```bash
UNITREE_ROBOT_IP=192.168.123.161
UNITREE_ROBOT_TYPE=go2
UNITREE_SDK_PATH=/opt/unitree_sdk2
ROS_DOMAIN_ID=42
```

## Safety Considerations

1. **Emergency stop** - Always have remote kill switch
2. **Clear area** - Remove obstacles before movement
3. **Leash mode** - Start with tethered operation
4. **Speed limits** - Start slow, increase gradually
5. **Battery monitoring** - Return to base before low

## Best Practices

1. **Calibration** - Regular sensor calibration
2. **Terrain awareness** - Know robot's limits
3. **Supervision** - Human oversight for new tasks
4. **Logging** - Record all operations
5. **Fallback behaviors** - Handle failures gracefully

## Resources

- [Unitree SDK2](https://github.com/unitreerobotics/unitree_sdk2)
- [Unitree ROS2](https://github.com/unitreerobotics/unitree_ros2)
- [Go2 Documentation](https://support.unitree.com/home/en/Go2_developer)
- [Unitree Developer](https://dev.unitree.com/)

---

Last Updated: January 25, 2026


---

## Related

- [[00-24 SO (Sovereign Operations)/02 Infrastructure (VPS, Servers & Networking) (The Body)/02 Infrastructure|02 Infrastructure]]
