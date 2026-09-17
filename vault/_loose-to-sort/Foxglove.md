---
id: app-foxglove
title: Foxglove Studio
category: Applications
status: documented
created: 2026-06-17
updated: 2026-06-17
tags: [robotics, ros, data-visualization, debugging, telemetry]
---

# Foxglove Studio

## Application Overview

**Foxglove Studio** is an open-source robotics visualization and debugging tool designed to help developers understand, monitor, and debug autonomous robots and embedded systems. It connects to ROS (Robot Operating System) networks and provides real-time 3D visualization, message inspection, and data recording capabilities.

| Property | Value |
|----------|-------|
| **Developer** | Foxglove (FKA Cruise) |
| **Bundle Identifier** | `dev.foxglove.studio` |
| **Version** | 2.53.1 |
| **Release Date** | Active (weekly releases) |
| **Architecture** | arm64 (Apple Silicon native) |
| **Install Location** | `/Applications/Foxglove.app` |
| **App Size** | 322 MB |
| **License** | AGPL v3 (Open Source) |
| **Website** | https://foxglove.dev |
| **Source Code** | https://github.com/foxglove |

## Copyright & Legal

```
Foxglove Technologies, Inc.
Licensed under GNU Affero General Public License v3
```

## Key Features

### Data Visualization
- **3D Visualization**: Point clouds, meshes, trajectories, coordinate frames
- **2D Plotting**: Real-time charts and time-series data
- **Map Visualization**: GPS, occupancy grids, laser scans
- **Image Streams**: Display camera feeds and vision outputs
- **3D Viewer**: Interactive pan/zoom/rotate with multiple cameras

### ROS Integration
- **ROS 1 & ROS 2 Support**: Full compatibility with both ecosystems
- **Message Inspection**: Subscribe to any ROS topic and inspect message structure
- **Message Publishing**: Publish messages to ROS topics for testing
- **Recording & Playback**: Record bag files and replay for debugging

### Debugging & Analysis
- **Live Streaming**: Connect to robots over network or serial
- **Local Playback**: Analyze pre-recorded bag files
- **Message Timeline**: Temporal inspection of message sequences
- **Parameter Inspection**: Monitor and set ROS parameters
- **Service Calling**: Invoke ROS services directly from UI

### Data Management
- **Multi-Layout Support**: Save and switch between dashboard configurations
- **Custom Panels**: Extensible plugin system for custom visualizations
- **Data Recording**: Record all topics in real-time
- **Performance Profiling**: Monitor CPU/memory usage of ROS nodes

## User Data Locations

| Path | Size | Purpose |
|------|------|---------|
| `~/Library/Application Support/Foxglove/` | Variable | User dashboards, layouts, preferences |
| `~/Library/Preferences/dev.foxglove.studio.plist` | KB | App preferences and settings |
| `~/Library/Containers/dev.foxglove.studio.quicklook/` | Variable | Quick Look thumbnail rendering |
| `~/Library/Application Scripts/dev.foxglove.studio.quicklook` | KB | Script support sandbox |
| Browser IndexedDB (Brave/Chrome): `https_app.foxglove.dev_0.indexeddb.leveldb` | MB | Web version cache if used |

## Supported Data Formats

| Format | Purpose |
|--------|---------|
| `.bag` | ROS 1 bag files (recording format) |
| `.mcap` | Modern Container And Protocol (next-gen recording) |
| `.rosbag2` | ROS 2 bag format |
| `.json` | Configuration and layout exports |
| `.yaml` | Sensor calibration and configuration |

## Connectivity Options

### Network
- **Local Network**: Connect to ROS master/ROS 2 daemon on LAN
- **Remote SSH**: SSH tunneling for remote robot access
- **WebSocket**: Real-time streaming over internet

### Serial
- **USB/Serial**: Direct serial connection for embedded systems
- **Baud Rate Configuration**: 9600 – 921600 bps

### Cloud
- **Foxglove Cloud**: Optional integration with hosted Foxglove dashboard
- **Data Upload**: Send recordings to cloud for team collaboration

## Performance Characteristics

- **Memory**: 300 MB – 2 GB depending on data rate and visualization complexity
- **CPU**: Single-thread UI; multi-threaded data processing
- **Network**: Efficient for 10+ Hz telemetry streams
- **Disk**: Minimal; bag file storage depends on recording duration
- **GPU**: Optional hardware acceleration for 3D viewer (Metal on macOS)

## ROS Ecosystem Integration

- **Supported Sensor Types**:
  - LiDAR point clouds (velodyne, sick, robosense, livox)
  - Camera images (single, stereo, depth)
  - IMU/odometry
  - GPS/GNSS
  - Radar
  - Custom message types

- **Supported Robot Types**:
  - Wheeled platforms (differential drive, Ackermann)
  - Manipulator arms
  - Drones/quadrotors
  - Humanoid robots
  - Custom kinematics

## Extension & Customization

- **Web Technology Stack**: Built with TypeScript, React, Three.js
- **Plugin API**: Create custom panels in TypeScript
- **Message Handlers**: Support for custom ROS message types
- **Theme Customization**: Light/dark mode support

## Integration with One Mind System

**Use Case**: Robotics development and autonomous systems monitoring
- **Pillar**: SO (Sovereign Operations) — robotics infrastructure, real-time systems
- **Pillar**: GE (Generational Empire) — if building autonomous products
- **Relevant Domains**: Robotics, hardware integration, real-time systems

## System Requirements

- **macOS**: 10.13+ (Intel) / 11.0+ (Apple Silicon)
- **RAM**: 2 GB minimum; 8+ GB recommended for large datasets
- **Disk**: 500 MB for app; storage for recordings scales with data
- **Network**: Ethernet or Wi-Fi (1 Mbps+ recommended)

## Common Workflows

1. **Local Development**: Connect to ROS master on dev machine
2. **Real Robot Debugging**: SSH tunnel to robot and inspect live data
3. **Offline Analysis**: Load recorded bag file and replay
4. **Sensor Calibration**: Visualize camera/LiDAR and adjust parameters
5. **Performance Profiling**: Monitor node CPU/memory over time
6. **Team Collaboration**: Share recordings and dashboards via Foxglove Cloud

## Removal Instructions

To completely uninstall Foxglove Studio:

```bash
rm -rf /Applications/Foxglove.app
rm -rf ~/Library/Application\ Support/Foxglove
rm -rf ~/Library/Preferences/dev.foxglove.studio.plist
rm -rf ~/Library/Containers/dev.foxglove.studio*
rm -rf ~/Library/Application\ Scripts/dev.foxglove.studio*
rm -rf ~/Library/Caches/dev.foxglove.studio*
```

## Alternatives

- **RViz** (ROS native, barebones)
- **RViz 2** (ROS 2 replacement)
- **Matlab/Simulink** (Commercial, simulation-focused)
- **CARLA Simulator** (Autonomous driving scenario tool)
- **Gazebo** (Physics-based simulation)

## Notes

- Free and open-source; no licensing costs
- Active development; updates weekly
- Growing industry adoption in robotics and autonomous vehicle teams
- Strong community support and documentation
- Can be self-hosted; no mandatory cloud dependency
- Designed for both hobbyists and professional teams

## Common Use Cases

- **Autonomous Vehicles**: Monitor perception, planning, control stacks
- **Mobile Robots**: Debug navigation and SLAM algorithms
- **Manipulators**: Visualize joint states and end-effector trajectories
- **Drones**: Real-time flight telemetry and sensor fusion
- **Research**: Data analysis and algorithm prototyping

---

*Last verified: June 17, 2026 (v2.53.1)*
