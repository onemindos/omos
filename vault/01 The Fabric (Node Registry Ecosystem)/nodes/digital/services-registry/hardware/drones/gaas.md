---
title: "GAAS"
type: reference
status: planned
tags: [platform, hardware, drones, autonomy, full-stack]
---
# GAAS — Generalized Autonomy Aviation System

> Full autonomous flight stack: perception → planning → control. End-to-end drone autonomy.

## What It Is
Complete open-source autonomous drone software stack. Combines SLAM, object detection, path planning, and flight control into one integrated system.

## Architecture
```
Perception Layer:
  - Stereo vision (depth estimation)
  - LiDAR point cloud processing
  - Object detection (YOLO)
  - Semantic segmentation

Planning Layer:
  - Global path planning (A*, RRT)
  - Local obstacle avoidance
  - Mission state machine

Control Layer:
  - PX4 offboard mode
  - Velocity/position commands
  - Emergency behaviors (RTL, hover, land)
```

## Key Features
| Feature | Notes |
|---------|-------|
| Full stack | Perception + planning + control |
| Simulation | Gazebo + PX4 SITL |
| Hardware | Jetson + PX4 + stereo camera + LiDAR |
| Languages | C++, Python |
| ROS | Built on ROS (original) |

## Hardware Requirements
- **Compute:** NVIDIA Jetson Nano/Xavier/Orin
- **Flight controller:** PX4 on Pixhawk
- **Cameras:** Stereo pair (Intel RealSense D435)
- **LiDAR:** Optional (RPLiDAR, Livox)

## Resources
- https://github.com/generalized-intelligence/GAAS
