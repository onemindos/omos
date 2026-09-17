---
title: "ROS 2 (Drone & Robot Autonomy Stack)"
type: reference
status: active
updated: 2026-06-06
tags: [platform, hardware, drones, robotics, autonomy, ros, slam, middleware, zenoh, dds]
---
# ROS 2 — Autonomy Stack

> Full autonomous stack: SLAM, obstacle avoidance, computer vision, swarm coordination. The "brain" on the companion computer. Applies to both drones and ground robots (DIMOS, OM1).

## Architecture

```
ROS 2 Node (Companion Computer)
  ↕ DDS / Zenoh transport (rmw layer)
  ↕ MAVROS or MAVSDK bridge
  ↕ MAVLink protocol
  Autopilot (PX4 / ArduPilot)
  ↕ PWM / CAN
  Motors, ESCs, Sensors
```

## ROS 2 ↔ Autopilot Bridges

| Bridge | Protocol | Best For |
|--------|----------|---------|
| **MAVROS** | MAVLink ↔ ROS topics | ArduPilot + ROS2, full telemetry access |
| **MAVSDK** | MAVLink SDK (Python/C++/Rust) | PX4, clean async API, no ROS required |
| **px4_ros_com** | uXRCE-DDS (direct) | PX4 native ROS2 messaging, zero-bridge |

## DDS / Zenoh Transport Layer

- Default ROS 2 transport: **DDS** (Data Distribution Service)
- Drop-in RT replacement: **Zenoh** — lower latency, resource-constrained friendly, NATS bridge available
- Switch: `RMW_IMPLEMENTATION=rmw_zenoh_cpp`
- Zenoh collapses the DDS overhead — critical for edge nodes (Pi, Jetson)

## Key Packages

| Package | Purpose |
|---------|---------|
| `nav2` | Navigation stack — path planning, costmaps, behaviors |
| `slam_toolbox` | 2D/3D SLAM |
| `mavros` | MAVLink ↔ ROS bridge |
| `px4_ros_com` | PX4 native uXRCE-DDS bridge |
| `micro_ros_agent` | microROS for MCUs (Arduino, STM32) |
| `robot_localization` | EKF/UKF sensor fusion |
| `octomap` | 3D occupancy mapping |

## Simulation Backends

| Sim | Use Case |
|-----|---------|
| **Gazebo (Ignition)** | Primary sim, tight ROS2 integration, sensor plugins |
| **PyBullet** | Fast physics, ML training, RL environments |
| **MuJoCo** | High-fidelity physics, DeepMind RL research |

## SITL Integration

All three sims support SITL (Software in the Loop) with ROS2:
- Gazebo: native via `ros_gz_bridge`
- ArduPilot SITL: connect via MAVROS to localhost MAVLink port
- PX4 SITL: connect via `px4_ros_com` or MAVSDK to sim instance

## OneMind OS Angle

ROS2's pub/sub model maps directly to NATS subjects:
- `robot.dimos.telemetry.>` ← ROS2 topic bridge → NATS
- `cmd.dimos.nav.goto` → NATS → ROS2 `/goal_pose`
- Zenoh's NATS bridge makes this native

## Related

- [[mavlink]] — protocol layer between ROS2 and autopilot
- [[mavros]] — ROS ↔ autopilot MAVLink bridge
- [[mavsdk]] — programmatic SDK (no ROS dependency)
- [[px4]] — autopilot
- [[ardupilot]] — autopilot
- [[zenoh]] — RT transport replacing DDS under ROS2
- [[gazebo]] — primary simulation
- [[foxglove-studio]] — visualization & debugging IDE for ROS2
