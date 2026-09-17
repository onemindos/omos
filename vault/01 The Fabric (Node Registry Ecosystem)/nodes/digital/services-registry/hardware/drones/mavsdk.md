---
title: "MAVSDK"
type: reference
status: active
updated: 2026-06-06
tags: [platform, hardware, drones, sdk, python, rust, c++, autopilot, px4, ardupilot]
---
# MAVSDK

> Modern drone SDK by Dronecode. Programmatic MAVLink control without ROS dependency. Python, C++, Rust, Swift. Built and optimized for PX4; also works with ArduPilot.

## What It Is

MAVSDK is an abstraction layer over MAVLink. Instead of crafting raw MAVLink messages, you call clean async APIs:
```python
await drone.action.takeoff()
await drone.action.goto_location(lat, lon, alt, heading)
await drone.action.land()
```

## Language Support

| Language | Package | Notes |
|----------|---------|-------|
| **Python** | `mavsdk` (pip) | Asyncio-based, most docs/examples |
| **C++** | mavsdk C++ | Low-level, maximum performance |
| **Rust** | mavsdk-rs | Async/await native |
| **Swift** | mavsdk-swift | iOS/macOS |

## Install

```bash
pip install mavsdk
```

## Core Usage (Python)

```python
import asyncio
from mavsdk import System

async def run():
    drone = System()
    await drone.connect(system_address="udp://:14540")

    # Wait for connection
    async for state in drone.core.connection_state():
        if state.is_connected:
            print("Connected")
            break

    # Arm and takeoff
    await drone.action.arm()
    await drone.action.takeoff()
    await asyncio.sleep(5)

    # Fly to position
    await drone.action.goto_location(
        latitude_deg=47.398155,
        longitude_deg=8.546756,
        absolute_altitude_m=505,
        yaw_deg=0
    )
    await asyncio.sleep(10)
    await drone.action.land()

asyncio.run(run())
```

## Plugin System

MAVSDK is plugin-based. Each domain has its own plugin:

| Plugin | Purpose |
|--------|---------|
| `action` | Arm, takeoff, land, goto, RTL |
| `telemetry` | Position, attitude, battery, GPS, mode |
| `mission` | Upload/download/execute waypoint missions |
| `offboard` | Direct velocity/attitude setpoints |
| `camera` | Trigger, settings, capture |
| `gimbal` | Gimbal control |
| `param` | Read/write autopilot parameters |
| `follow_me` | Follow a moving target |

## Connection URLs

```python
# UDP (SITL default)
"udp://:14540"

# Serial (hardware, USB)
"serial:///dev/ttyUSB0:921600"

# TCP
"tcp://192.168.1.10:5760"
```

## SITL Testing

```bash
# Start PX4 SITL
make px4_sitl gz_x500

# Run MAVSDK script
python3 my_mission.py  # connects to udp://:14540
```

## MAVSDK vs MAVROS vs DroneKit

| SDK | ROS Required | Language | Maintainability |
|-----|-------------|----------|----------------|
| **MAVSDK** | No | Python/C++/Rust/Swift | ✅ Actively maintained |
| **MAVROS** | Yes (ROS2) | ROS topics/services | ✅ Active |
| **DroneKit** | No | Python | ⚠️ Mostly maintenance mode |

## OneMind OS Angle

MAVSDK is the preferred non-ROS drone SDK. Pattern:
```
NATS subject received → agent callback → MAVSDK async command → drone executes
```
No ROS overhead for simple point-to-point missions or remote command execution.

## Related

- [[mavlink]] — protocol underneath MAVSDK
- [[px4]] — primary target autopilot
- [[ardupilot]] — also supported
- [[mavros]] — ROS2 alternative bridge
- [[ros2-drones]] — full autonomy stack with ROS2
