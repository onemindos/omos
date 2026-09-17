---
title: "Foxglove Studio"
type: reference
status: active
updated: 2026-06-06
tags: [platform, hardware, drones, robotics, visualization, debugging, ros2, telemetry, ide]
---
# Foxglove Studio

> Visualization and debugging IDE for robotics. The "DevTools" for ROS2. Live and recorded data playback. Works with ROS2 bags, MCAP, WebSocket streams. Cloud + desktop.

## What It Is

Foxglove Studio is an open-source robotics observability platform:
- **Live data** — connect to running ROS2 system via Foxglove WebSocket
- **Recorded data** — load `.bag` (ROS2) or `.mcap` files for post-flight analysis
- **Panels** — modular dashboard: 3D viewer, plots, images, maps, raw messages, diagnostics
- **Web + Desktop** — browser or Electron app

## Install

```bash
# Electron desktop app
brew install --cask foxglove-studio

# Or download from https://foxglove.dev/download
# Or run in browser: https://studio.foxglove.dev
```

## Connecting to Live ROS2

```bash
# Install Foxglove WebSocket bridge on robot
pip install foxglove-websocket

# Or via ROS2 package
sudo apt install ros-humble-foxglove-bridge

# Launch bridge
ros2 launch foxglove_bridge foxglove_bridge_launch.xml

# In Foxglove Studio → Open Connection → Foxglove WebSocket
# ws://robot-ip:8765
```

## Key Panels

| Panel | Purpose |
|-------|---------|
| **3D** | Point clouds, robot models (URDF), transforms, markers |
| **Image** | Camera feeds, computer vision overlays |
| **Plot** | Time-series data (battery, attitude, velocities) |
| **Map** | GPS tracks, waypoints on satellite imagery |
| **Raw Messages** | Inspect any ROS2 topic message contents |
| **State Transitions** | Mode/state changes over time |
| **Diagnostics** | ROS diagnostics aggregator |
| **Teleop** | Send twist commands to robot |

## MCAP Format

Foxglove's recommended recording format (replaces ROS2 `.bag`):
```bash
# Record to MCAP
ros2 bag record -a --storage mcap -o my_flight

# Convert existing bag to MCAP
mcap convert input.db3 output.mcap
```
MCAP is self-indexing, seekable, works with web viewer without download.

## Cloud Platform

Foxglove cloud (foxglove.dev):
- Upload MCAP recordings for team sharing
- Stream live robots to cloud dashboard
- Annotations and event marking

## Post-Flight Analysis Workflow

```
Flight completes
  → ros2 bag recorded (MCAP)
  → Download to laptop
  → Open in Foxglove Studio
  → Scrub timeline, inspect anomalies
  → Plot altitude vs battery vs GPS accuracy
  → Export clips for documentation
```

## Foxglove vs RViz vs PlotJuggler

| Tool | Live | Playback | 3D | Plots | Web |
|------|------|---------|-----|-------|-----|
| **Foxglove** | ✅ | ✅ | ✅ | ✅ | ✅ |
| **RViz2** | ✅ | ✅ | ✅ | ❌ | ❌ |
| **PlotJuggler** | ✅ | ✅ | ❌ | ✅ Best | ❌ |

Use Foxglove as the primary. Use PlotJuggler when deep time-series analysis is needed.

## OneMind OS Angle

Foxglove is the pilot/engineer view into fleet behavior:
- Connect to drone companion computer via WebSocket during flight
- Ground operator monitors 3D view, telemetry plots, camera feeds simultaneously
- Post-mission: upload MCAP to cloud, Legacy auto-analyzes anomalies

## Related

- [[ros2-drones]] — the stack Foxglove visualizes
- [[qgroundcontrol]] — the GCS (MAVLink layer); Foxglove covers ROS2 layer above it
- [[gazebo]] — Foxglove can visualize Gazebo sim data via ROS2 bridge
- [[mavlink]] — below Foxglove's visibility layer
