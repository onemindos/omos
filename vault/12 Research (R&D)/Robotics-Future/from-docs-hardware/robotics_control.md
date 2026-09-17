# Robotics Control

## Gesture → Robot Mapping for OneMind OS

This document covers how to translate gesture input into robot/drone control commands.

---

## Safety First ⚠️

**Critical Safety Rules:**

1. **Emergency Stop**: Open palm = STOP everything immediately
2. **Confirmation Required**: Dangerous actions require explicit confirmation
3. **Geofencing**: Drones stay within defined boundaries
4. **Rate Limiting**: Limit command frequency to prevent jerky motion
5. **Watchdog**: If gesture input stops, robot defaults to hover/stop

```python
@dataclass
class SafetyLimits:
    max_velocity: float = 2.0  # m/s
    max_acceleration: float = 1.0  # m/s²
    max_altitude: float = 50.0  # meters
    min_altitude: float = 0.5  # meters
    geofence_radius: float = 100.0  # meters from origin
    command_timeout: float = 1.0  # seconds before auto-stop
    confirmation_timeout: float = 5.0  # seconds to confirm actions
```

---

## Gesture Mappings

### Universal Gestures

| Gesture             | Action              | Works For             |
| ------------------- | ------------------- | --------------------- |
| **Open Palm**       | Emergency STOP      | All robots            |
| **Closed Fist**     | Hold position       | Drones, arms          |
| **Thumbs Up**       | Confirm action      | All                   |
| **Thumbs Down**     | Cancel action       | All                   |
| **Point**           | Direct movement     | Drones, mobile robots |
| **Pinch**           | Fine control / Grab | Robot arms            |
| **Two-hand spread** | Zoom / Expand       | Visualization         |

---

## Drone Control

### Control Modes

**1. Point-to-Fly**

```
User points → Drone flies in that direction
Point up → Ascend
Point down → Descend
Point forward → Move forward
```

**2. Follow Mode**

```
Drone maintains fixed offset from user
User moves → Drone follows
Good for filming, inspection
```

**3. Virtual Joystick**

```
Left hand: Throttle + Yaw
Right hand: Pitch + Roll
Like a traditional controller, but with hands
```

### Implementation:

```python
class DroneGestureController:
    def __init__(self, drone, safety: SafetyLimits):
        self.drone = drone
        self.safety = safety
        self.mode = "point_to_fly"
        self.last_command_time = time.time()

    async def handle_gesture(self, gesture: NormalizedGesture):
        # Emergency stop
        if gesture.gesture_type == "open_palm":
            await self.drone.emergency_stop()
            return

        # Reset watchdog
        self.last_command_time = time.time()

        if self.mode == "point_to_fly":
            if gesture.gesture_type == "point":
                direction = gesture.direction

                # Map gesture direction to world coordinates
                velocity = self._calculate_velocity(direction)

                # Apply safety limits
                velocity = self._clamp_velocity(velocity)

                await self.drone.set_velocity(velocity)

            elif gesture.gesture_type == "closed_fist":
                await self.drone.hover()

    def _calculate_velocity(self, direction):
        # Convert normalized direction to velocity
        # Scale by maximum velocity
        return {
            "vx": direction["x"] * self.safety.max_velocity,
            "vy": direction["y"] * self.safety.max_velocity,
            "vz": -direction["z"] * self.safety.max_velocity,  # Z inverted
        }

    def _clamp_velocity(self, velocity):
        speed = np.sqrt(velocity["vx"]**2 + velocity["vy"]**2 + velocity["vz"]**2)
        if speed > self.safety.max_velocity:
            scale = self.safety.max_velocity / speed
            return {k: v * scale for k, v in velocity.items()}
        return velocity

    async def watchdog(self):
        """Stop drone if no commands received."""
        while True:
            if time.time() - self.last_command_time > self.safety.command_timeout:
                await self.drone.hover()
            await asyncio.sleep(0.1)
```

### Supported Drone Platforms

| Drone                | SDK             | Difficulty | Notes              |
| -------------------- | --------------- | ---------- | ------------------ |
| **DJI Tello**        | djitellopy      | Easy       | Great for learning |
| **DJI Mini series**  | Mobile SDK only | Hard       | No computer SDK    |
| **DJI Mavic**        | DJI Mobile SDK  | Medium     | Via mobile bridge  |
| **Parrot Bebop**     | olympe          | Medium     | Good Python SDK    |
| **Custom (Pixhawk)** | MAVLink         | Hard       | Full control       |

---

## Robot Arm Control

### Control Modes

**1. Direct Mapping**

```
User hand position → End effector position
Hand rotation → Gripper rotation
Pinch → Close gripper
Open → Release
```

**2. Relative Mode**

```
Hand movement = Delta movement
Smaller, more precise control
Good for fine manipulation
```

**3. Waypoint Mode**

```
Point at location → Set waypoint
Thumbs up → Execute movement
Closed fist → Clear waypoints
```

### Implementation:

```python
class ArmGestureController:
    def __init__(self, arm, workspace_bounds):
        self.arm = arm
        self.bounds = workspace_bounds
        self.mode = "direct"
        self.scaling = 0.5  # Hand movement to arm movement ratio

    async def handle_gesture(self, gesture: NormalizedGesture):
        if gesture.gesture_type == "open_palm":
            await self.arm.stop()
            return

        if self.mode == "direct":
            if gesture.hand == "right":
                # Map hand position to arm workspace
                target = self._map_to_workspace(gesture.position)

                # Clamp to bounds
                target = self._clamp_to_bounds(target)

                # Move arm
                await self.arm.move_to(target)

            # Gripper control
            if gesture.gesture_type == "pinch":
                await self.arm.close_gripper()
            elif gesture.gesture_type == "open_palm" and gesture.hand == "left":
                await self.arm.open_gripper()

    def _map_to_workspace(self, hand_pos):
        # Map normalized 0-1 hand position to arm workspace
        return {
            "x": self.bounds["x_min"] + hand_pos["x"] * (self.bounds["x_max"] - self.bounds["x_min"]),
            "y": self.bounds["y_min"] + hand_pos["y"] * (self.bounds["y_max"] - self.bounds["y_min"]),
            "z": self.bounds["z_min"] + hand_pos["z"] * (self.bounds["z_max"] - self.bounds["z_min"]),
        }
```

### Supported Platforms

| Robot                     | Controller   | SDK        | Notes           |
| ------------------------- | ------------ | ---------- | --------------- |
| **Servo-based**           | Arduino/RPi  | Serial     | Custom builds   |
| **UR (Universal Robots)** | URScript     | Python     | Industrial      |
| **Dobot**                 | Dobot Studio | Python     | Affordable      |
| **xArm**                  | xArm SDK     | Python     | Good Python API |
| **Franka Emika**          | Franka CI    | C++/Python | Research grade  |

---

## NATS Topics for Robotics

```
robotics.{robot_id}.command    # Outgoing commands
robotics.{robot_id}.status     # Current state
robotics.{robot_id}.telemetry  # Position, battery, etc.
robotics.{robot_id}.emergency  # Emergency signals
robotics.{robot_id}.geofence   # Boundary violations
```

### Command Schema:

```json
{
  "type": "velocity",
  "robot_id": "drone_01",
  "values": {
    "vx": 0.5,
    "vy": 0.0,
    "vz": 0.2,
    "yaw_rate": 0.0
  },
  "timestamp": 1706648400000,
  "source": "gesture",
  "priority": "normal"
}
```

### Status Schema:

```json
{
  "robot_id": "drone_01",
  "state": "flying",
  "position": { "x": 1.2, "y": 0.5, "z": 2.0 },
  "orientation": { "roll": 0.0, "pitch": 0.02, "yaw": 1.57 },
  "battery_percent": 85,
  "gps": { "lat": 41.8781, "lon": -87.6298 },
  "timestamp": 1706648400000
}
```

---

## Example: Full Pipeline

```python
# Complete gesture-to-drone example

import asyncio
import nats
from djitellopy import Tello

async def main():
    # Connect to NATS
    nc = await nats.connect("nats://localhost:4222")

    # Initialize drone
    drone = Tello()
    drone.connect()
    drone.takeoff()

    # Gesture handler
    async def handle_gesture(msg):
        data = json.loads(msg.data)

        if data["gesture_type"] == "open_palm":
            drone.land()

        elif data["gesture_type"] == "point":
            direction = data["direction"]
            # Map to RC commands (-100 to 100)
            rc_x = int(direction["x"] * 50)
            rc_y = int(direction["y"] * 50)
            rc_z = int(-direction["z"] * 30)
            drone.send_rc_control(rc_x, rc_y, rc_z, 0)

        elif data["gesture_type"] == "closed_fist":
            drone.send_rc_control(0, 0, 0, 0)  # Hover

    # Subscribe
    await nc.subscribe("gesture.classified", cb=handle_gesture)

    # Keep alive
    try:
        while True:
            await asyncio.sleep(0.1)
    finally:
        drone.land()
        await nc.close()

asyncio.run(main())
```

---

## Future: Swarm Control

Gesture control for multiple robots:

```
Two-hand gestures for swarm formation control:

┌─ Left Hand ─┐     ┌─ Right Hand ─┐
│             │     │              │
│  Formation  │     │  Movement    │
│  Control    │     │  Direction   │
│             │     │              │
└─────────────┘     └──────────────┘

Spread hands apart → Expand formation
Bring together → Contract formation
Rotate left hand → Rotate formation
Right hand point → All move in direction
```

This is a longer-term goal requiring coordination protocols and more sophisticated gesture recognition.
