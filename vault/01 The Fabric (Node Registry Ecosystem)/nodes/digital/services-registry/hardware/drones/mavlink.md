---
title: "MAVLink"
type: reference
status: active
updated: 2026-06-06
tags: [platform, hardware, drones, protocol, communication, messaging]
---
# MAVLink

> Universal drone messaging protocol. The common language between autopilot, ground station, companion computer, and fleet controller. Lightweight binary protocol over serial/UDP/TCP.

## What It Is

MAVLink (Micro Air Vehicle Link) is a header-only message marshalling library. Every drone subsystem speaks it:
- Autopilot → Ground Control (telemetry, mission upload)
- Companion Computer → Autopilot (commands, offboard mode)
- Ground Station → Drone (waypoints, params)
- Fleet Controller → Multiple Drones

## Versions

| Version | Status | Notes |
|---------|--------|-------|
| MAVLink 1 | Legacy | 8-byte header, no signing |
| **MAVLink 2** | Current | Extended message IDs, packet signing, tunnel support |

## Key Concepts

| Term | Meaning |
|------|---------|
| System ID | Unique ID per vehicle (1 = autopilot default) |
| Component ID | Subsystem ID (autopilot, camera, companion, GCS) |
| Message | Typed binary struct (HEARTBEAT, SET_MODE, MISSION_ITEM, etc.) |
| Stream | Continuous telemetry messages at configurable Hz |
| Command | One-shot request/response (ARM, TAKEOFF, LAND) |

## Common Messages

- `HEARTBEAT` — alive signal (1 Hz)
- `GLOBAL_POSITION_INT` — lat/lon/alt/hdg
- `ATTITUDE` — roll/pitch/yaw
- `COMMAND_LONG` — send MAV_CMD (arm, takeoff, mode change)
- `MISSION_ITEM_INT` — waypoint upload
- `RC_CHANNELS` — raw RC input values
- `SET_POSITION_TARGET_LOCAL_NED` — offboard position commands

## Connections

```
Autopilot (UART/USB) ──── MAVLink ──── Companion Computer (MAVROS / MAVSDK)
                                        └── GCS (QGroundControl / Mission Planner)
```

- Serial: typically `/dev/ttyAMA0` at 921600 baud (Pixhawk ↔ Pi)
- UDP: port 14550 (GCS), 14540 (MAVSDK), 14545 (secondary)
- TCP: supported but not preferred (stateful)

## SDKs & Bridges

| Tool | Language | Notes |
|------|----------|-------|
| **MAVROS** | ROS2 topic bridge | Best for ROS2 integration |
| **MAVSDK** | Python, C++, Rust, Swift | Clean async, PX4-native |
| **DroneKit** | Python | Legacy, ArduPilot, deprecated-ish |
| **pymavlink** | Python | Low-level, direct message access |

## OneMind OS Angle

MAVLink is the NATS equivalent at the vehicle edge. Bridge pattern:
```
NATS subject: cmd.dimos.nav.goto
  → MAVSDK/MAVROS bridge
  → MAVLink COMMAND_LONG (MAV_CMD_DO_SET_POSITION)
  → Autopilot executes
```

## Related

- [[px4]] — autopilot that speaks MAVLink
- [[ardupilot]] — autopilot that speaks MAVLink
- [[mavsdk]] — modern MAVLink SDK
- [[mavros]] — ROS2 MAVLink bridge
- [[qgroundcontrol]] — GCS that speaks MAVLink
- [[ros2-drones]] — upper layer using MAVLink bridges
