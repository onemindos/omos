---
title: "PX4"
type: reference
status: active
updated: 2026-06-06
tags: [platform, hardware, drones, autopilot, open-source, px4, dronecode, pixhawk]
---
# PX4

> Enterprise-grade open-source autopilot. Dronecode Foundation backed. Clean modular architecture. Native ROS2 integration via uXRCE-DDS. MAVSDK is its primary programmatic SDK.

## Architecture

PX4 is built on a publish-subscribe uORB message bus internally — all modules communicate via topics, no direct coupling. Maps well to NATS thinking.

```
uXRCE-DDS Agent (Companion Computer)
  ↕ UDP/Serial
PX4 Firmware (Pixhawk)
  ├── Estimator (EKF2)
  ├── Navigator (mission, RTL, loiter)
  ├── Commander (state machine, arming)
  ├── Mixer (output to ESCs)
  └── uORB (internal pub/sub)
```

## ROS2 Native Bridge (px4_ros_com)

PX4 has a **native** ROS2 bridge via **uXRCE-DDS** — no MAVROS needed:
```bash
# On companion computer
MicroXRCEAgent udp4 -p 8888

# On PX4 SITL or hardware
param set XRCE_DDS_EN 1
```

This exposes PX4 uORB topics directly as ROS2 topics — lower latency than MAVROS.

## SITL (Software In The Loop)

```bash
# Docker (easiest)
docker run --rm -it jonasvautherin/px4-gazebo-headless:latest

# Native
git clone https://github.com/PX4/PX4-Autopilot.git --recursive
cd PX4-Autopilot
make px4_sitl gz_x500          # Gazebo Harmonic + X500 quadcopter

# Connect MAVSDK
python3 mavsdk_example.py       # targets udp://:14540
```

## Flight Modes

| Mode | Description |
|------|-------------|
| STABILIZED | Manual + attitude stabilization |
| POSITION | GPS position hold |
| OFFBOARD | Full programmatic control (MAVSDK/MAVROS) |
| MISSION | Execute waypoint mission |
| RETURN | Return to launch |
| TAKEOFF / LAND | Automated vertical |

## Key Params

- `SYS_AUTOSTART` — airframe type (X500, S500, etc.)
- `EKF2_AID_MASK` — which sensor sources EKF uses
- `XRCE_DDS_EN` — enable native ROS2 bridge
- `MAV_SYS_ID` — system ID for multi-vehicle

## PX4 vs ArduPilot

| Feature | PX4 | ArduPilot |
|---------|-----|-----------|
| ROS2 native | ✅ uXRCE-DDS | Via MAVROS |
| MAVSDK | ✅ Primary | ✅ Supported |
| Rover/Sub | Limited | ✅ Full |
| Lua scripting | ❌ | ✅ |
| Architecture | uORB modular | Scheduler-based |
| Dronecode backed | ✅ | Community |

## Gazebo Integration

- `make px4_sitl gz_x500` — spawns X500 in Gazebo Harmonic
- `make px4_sitl gazebo-classic_iris` — older Gazebo Classic
- Sensor plugins: IMU, GPS, mag, optical flow, depth camera available

## OneMind OS Angle

PX4 + uXRCE-DDS + ROS2 is the cleanest path to NATS integration:
```
NATS subject: cmd.drone1.nav.goto
  → ROS2 node
  → px4_ros_com publisher
  → uXRCE-DDS
  → PX4 OFFBOARD mode
  → vehicle motion
```

## Related

- [[mavlink]] — protocol layer
- [[mavsdk]] — primary SDK (Python/C++/Rust)
- [[mavros]] — alternative ROS2 bridge (if px4_ros_com not used)
- [[ardupilot]] — competing autopilot (better for rovers/subs)
- [[qgroundcontrol]] — GCS
- [[gazebo]] — primary SITL sim
- [[ros2-drones]] — upper autonomy stack
