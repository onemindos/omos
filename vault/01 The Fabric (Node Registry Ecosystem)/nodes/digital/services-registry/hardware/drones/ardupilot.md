---
title: "ArduPilot"
type: reference
status: active
updated: 2026-06-06
tags: [platform, hardware, drones, robotics, autopilot, open-source, ardupilot, pixhawk]
---
# ArduPilot

> Most mature open-source autopilot firmware. Supports copter, plane, rover, submarine, boat. Runs on Pixhawk hardware family. Industry standard for research and custom builds.

## Variants

| Firmware | Vehicle |
|----------|---------|
| **ArduCopter** | Multirotors (quad, hex, octo, VTOL) |
| **ArduPlane** | Fixed-wing |
| **ArduRover** | Ground vehicles (UGV) |
| **ArduSub** | Underwater ROV |
| **ArduBlimp** | Lighter-than-air |

All variants share the same MAVLink interface and parameter system.

## Hardware

Runs on Pixhawk hardware family:
- **Pixhawk 6C / 6X** — current gen, CAN bus, redundant sensors
- **Cube Orange/Orange+** — industrial, vibration-isolated IMU
- **Pixhawk 4 Mini** — smaller form factor
- **Matek, Holybro, mRo** — compatible Pixhawk-clone boards

## Key Features

- **Mission Planner** (Windows) / **QGroundControl** (cross-platform) — GCS
- **MAVROS / MAVSDK** — ROS2 / programmatic bridge
- **SITL** — Software In The Loop sim (no hardware required)
- **EKF3** — Extended Kalman Filter for sensor fusion
- **Optical Flow** — GPS-denied indoor hover
- **Terrain Following** — digital elevation model awareness
- **Geofence** — enforced no-fly boundaries

## SITL (Software In The Loop)

```bash
# Clone and build
git clone https://github.com/ArduPilot/ardupilot.git
cd ardupilot && git submodule update --init --recursive
./Tools/environment_install/install-prereqs-ubuntu.sh -y

# Launch SITL (copter)
cd ArduCopter && sim_vehicle.py -v ArduCopter --console --map

# Connect MAVROS
ros2 launch mavros apm.launch fcu_url:=udp://127.0.0.1:14550@14555

# Connect QGroundControl → UDP 14550
```

## Flight Modes (Copter)

| Mode | Description |
|------|-------------|
| STABILIZE | Manual with attitude stabilization |
| ALT_HOLD | Altitude hold, manual lat/lon |
| LOITER | GPS hold, manual |
| AUTO | Execute waypoint mission |
| GUIDED | Programmatic commands (MAVROS/MAVSDK) |
| GUIDED_NOGPS | Indoor offboard without GPS |
| POSHOLD | GPS position hold, comfortable manual |
| RTL | Return to launch |

## Configuration

- Parameters stored in `EEPROM` on flight controller
- Upload via QGroundControl param editor or `MAV_PARAM` messages
- Key params: `ARMING_CHECK`, `INS_ACCEL_FILTER`, `COMPASS_DEC`, `FS_*` (failsafes)

## Companion Computer Integration

```
Pi / Jetson ──serial/USB──► Pixhawk (ArduPilot)
    ↕
MAVROS (ROS2) ↔ ROS2 nodes (nav2, slam_toolbox, vision)
```

## ArduPilot vs PX4

| Feature | ArduPilot | PX4 |
|---------|-----------|-----|
| Maturity | Higher (15+ yrs) | High (10+ yrs) |
| Rover/Sub support | ✅ Full | Limited |
| ROS2 native | Via MAVROS | Native uXRCE-DDS |
| MAVSDK support | ✅ | ✅ (primary) |
| Community size | Larger | Large |
| Scripting (Lua) | ✅ Built-in | ❌ |

## OneMind OS Angle

ArduPilot runs on OM1 rover and any future drone builds. GUIDED mode + MAVROS = full programmatic control from NATS-triggered ROS2 commands. Lua scripting enables onboard autonomous behaviors without companion computer.

## Related

- [[mavlink]] — protocol ArduPilot speaks
- [[mavros]] — ROS2 bridge
- [[mavsdk]] — SDK (PX4-preferred but works with ArduPilot)
- [[px4]] — competing autopilot
- [[qgroundcontrol]] — GCS
- [[ros2-drones]] — upper autonomy stack
