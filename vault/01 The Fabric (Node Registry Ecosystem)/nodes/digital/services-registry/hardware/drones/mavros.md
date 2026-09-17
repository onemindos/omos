---
title: "MAVROS"
type: reference
status: active
updated: 2026-06-06
tags: [platform, hardware, drones, robotics, ros2, mavlink, bridge, autopilot]
---
# MAVROS

> ROS 2 package that bridges MAVLink ↔ ROS topics/services. The glue between the ROS2 autonomy stack and any MAVLink autopilot (ArduPilot, PX4).

## What It Does

MAVROS exposes autopilot state as ROS2 topics and provides services for sending commands:
- Subscribes to ROS2 topics → translates → MAVLink messages → autopilot
- Receives MAVLink telemetry → publishes as ROS2 topics → ROS2 nodes consume

## Key Topics

| ROS2 Topic | Direction | MAVLink Message |
|------------|-----------|----------------|
| `/mavros/state` | ↓ read | HEARTBEAT |
| `/mavros/local_position/pose` | ↓ read | LOCAL_POSITION_NED |
| `/mavros/global_position/global` | ↓ read | GLOBAL_POSITION_INT |
| `/mavros/imu/data` | ↓ read | ATTITUDE / RAW_IMU |
| `/mavros/rc/in` | ↓ read | RC_CHANNELS |
| `/mavros/setpoint_position/local` | ↑ send | SET_POSITION_TARGET_LOCAL_NED |
| `/mavros/setpoint_velocity/cmd_vel` | ↑ send | SET_POSITION_TARGET_LOCAL_NED |

## Key Services

| Service | Purpose |
|---------|---------|
| `/mavros/cmd/arming` | Arm/disarm |
| `/mavros/cmd/takeoff` | Takeoff to altitude |
| `/mavros/cmd/land` | Land command |
| `/mavros/set_mode` | Switch flight mode (GUIDED, STABILIZE, etc.) |
| `/mavros/mission/push` | Upload waypoint mission |

## Install

```bash
sudo apt install ros-humble-mavros ros-humble-mavros-extras
# Install geographiclib (required)
wget https://raw.githubusercontent.com/mavlink/mavros/master/mavros/scripts/install_geographiclib_datasets.sh
chmod +x install_geographiclib_datasets.sh && sudo ./install_geographiclib_datasets.sh
```

## Launch

```bash
# Serial connection (Pixhawk via USB)
ros2 launch mavros apm.launch fcu_url:=/dev/ttyUSB0:921600

# UDP (SITL)
ros2 launch mavros apm.launch fcu_url:=udp://127.0.0.1:14550@14555
```

## PX4 vs ArduPilot

| Feature | ArduPilot | PX4 |
|---------|-----------|-----|
| Launch file | `apm.launch` | `px4.launch` |
| Mode naming | `GUIDED`, `STABILIZE` | `OFFBOARD`, `POSCTL` |
| Preferred bridge | MAVROS | MAVROS or native px4_ros_com |
| Native ROS2 | No (MAVLink only) | Yes (uXRCE-DDS via px4_ros_com) |

## OneMind OS Angle

MAVROS is the vehicle-side bridge when using ROS2 on a companion computer (Jetson, Pi):
```
NATS → agent → ROS2 node → MAVROS service → autopilot
```
For PX4-only deployments, `px4_ros_com` (uXRCE-DDS) can replace MAVROS entirely.

## Related

- [[mavlink]] — protocol MAVROS translates
- [[ros2-drones]] — upper autonomy layer
- [[ardupilot]] — primary target autopilot
- [[px4]] — alternative autopilot (has native ROS2 bridge)
- [[mavsdk]] — alternative SDK (no ROS dependency, async Python/Rust)
