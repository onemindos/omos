---
id: app-qgroundcontrol
title: QGroundControl
category: Applications
status: documented
created: 2026-06-17
updated: 2026-06-17
tags: [uav, drone, autopilot, mission-planning, ros]
---

# QGroundControl

## Application Overview

**QGroundControl** is the open-source ground control station for autonomous vehicles and drones. It provides mission planning, vehicle telemetry monitoring, parameter tuning, and real-time control for aerial, surface, and underwater vehicles using the MAVLink protocol.

| Property | Value |
|----------|-------|
| **Developer** | QGroundControl Project (Auterion, PX4 Community) |
| **Bundle Identifier** | `org.qgroundcontrol.QGroundControl` |
| **Version** | 5.0 |
| **Release Date** | Stable (quarterly releases) |
| **Architecture** | arm64, x86_64 (Universal Binary) |
| **Install Location** | `/Applications/QGroundControl.app` |
| **App Size** | 831 MB |
| **License** | Apache License 2.0 (Open Source) |
| **Website** | http://qgroundcontrol.com |
| **Source Code** | https://github.com/mavlink/qgroundcontrol |

## Copyright & Legal

```
QGroundControl Contributors
Licensed under Apache License 2.0
Maintained by Auterion and PX4 Community
```

## Key Features

### Mission Planning
- **Waypoint Planning**: Create multi-waypoint missions with speed/altitude control
- **Survey Missions**: Automated grid-based mapping missions for drones
- **Corridor Scanning**: Corridor-based survey patterns
- **Geofencing**: Set no-fly zones and altitude limits
- **Rally Points**: Emergency return-to-base waypoints
- **Action Sequencing**: Trigger payloads, camera actions at waypoints

### Vehicle Telemetry
- **Real-Time Data**: Altitude, speed, heading, battery status, GPS
- **Multi-Vehicle Support**: Monitor multiple drones simultaneously
- **Attitude Indicator**: Pitch, roll, yaw visualization
- **System Health**: Check motor status, sensor health, compass calibration
- **Event Logging**: Record all mission events and anomalies

### Parameter Tuning
- **Firmware Configuration**: Adjust flight controller parameters
- **PID Tuning**: Optimize attitude and velocity control loops
- **Sensor Calibration**: Compass, accelerometer, gyro calibration
- **RC Channel Mapping**: Configure radio controller inputs
- **Failsafe Configuration**: Set behavior on signal loss or low battery

### Vehicle Communication
- **MAVLink Protocol**: Support for MAVLink v1 and v2
- **Wireless Telemetry**: Radio modem connection (3DR, other vendors)
- **UDP/TCP Connection**: Network-based vehicle connectivity
- **Serial Connection**: Direct USB/UART link to autopilot
- **Cloud Integration**: Optional cloud connectivity for remote operations

### Flight Analysis
- **Flight Logs**: Download and analyze flight data
- **Log Viewer**: Visualize sensor readings over time
- **Replay**: Replay mission execution step-by-step
- **Performance Metrics**: Analyze stability, accuracy, efficiency

## User Data Locations

| Path | Size | Purpose |
|------|------|---------|
| `~/Library/Application Support/QGroundControl/` | MB-GB | Missions, logs, vehicle profiles |
| `~/Library/Application Support/QGroundControl/Missions/` | Variable | Saved mission files (.plan) |
| `~/Library/Application Support/QGroundControl/Logs/` | MB-GB | Flight logs (binary) |
| `~/Library/Application Support/QGroundControl/Vehicle Cache/` | Variable | Vehicle configuration cache |
| `~/Library/Preferences/org.qgroundcontrol.QGroundControl.plist` | KB | User preferences |
| `~/Library/Caches/org.qgroundcontrol.QGroundControl/` | Variable | Map tiles, media cache |

## Supported Autopilots

| Autopilot | Firmware | Support |
|-----------|----------|---------|
| **PX4** | PX4 Autopilot | Full support |
| **ArduPilot** | Copter, Plane, Rover, Sub | Full support |
| **Pixhawk** | PX4-based | Full support |
| **Cube** | ArduPilot-compatible | Full support |
| **Intelligent Flight Control System (IFCS)** | Custom | Limited support |

## Supported Vehicle Types

- **Multicopter** (Quadcopter, Hexacopter, Octocopter)
- **Fixed-Wing Aircraft**
- **VTOL** (Vertical Takeoff and Landing)
- **Ground Rovers** (Autonomous vehicles)
- **Boats/USVs** (Autonomous surface vehicles)
- **Submarines/AUVs** (Autonomous underwater vehicles)

## File Format Support

| Format | Type | Purpose |
|--------|------|---------|
| `.plan` | QGC Mission | Mission files (JSON-based) |
| `.log` | Binary Flight Log | Telemetry and sensor data |
| `.bin` | Autopilot Firmware | Firmware images for flashing |
| `.param` | Parameter File | Vehicle configuration backup |
| `.kml` | Keyhole Markup Language | Export missions for Google Earth |

## Connectivity Options

### Vehicle Link
- **Serial Port**: Direct USB connection to autopilot (CH340, PL2303, FTDI)
- **3DR Radio**: 915 MHz or 433 MHz wireless telemetry module
- **WiFi Bridge**: Connect via router or access point
- **Network (UDP/TCP)**: IP-based connection for testing/simulation
- **Cloud Link**: Via Auterion cloud for remote operations

### External Devices
- **Camera**: USB webcam for vehicle FPV monitoring
- **Joystick/Gamepad**: Control vehicle with game controller
- **RC Radio**: Passthrough RC input to vehicle
- **GPS**: External GPS source for ground station position

## Performance Characteristics

- **Memory**: 200-800 MB depending on active missions and logs
- **CPU**: Low-overhead; efficient for older hardware
- **Network**: Low bandwidth requirement (1-4 kbps for telemetry)
- **Disk**: Mission files 10-100 KB; logs 1-100+ MB per flight
- **GUI**: Responsive; 30 FPS video streams supported

## Mission Planning Workflow

1. **Create New Plan**: Start blank or template-based mission
2. **Add Waypoints**: Click on map to define flight path
3. **Configure Parameters**: Set altitude, speed, camera actions
4. **Set Geofence**: Define safe flight boundaries
5. **Upload to Vehicle**: Transfer mission to autopilot
6. **Arm & Fly**: Execute mission automatically
7. **Monitor**: Watch telemetry in real-time
8. **Analyze**: Download logs post-flight for review

## Advanced Features

- **Survey (Mapping) Missions**: Generate optimal camera flight paths
- **Structure Scanning**: Spiral and tangential scanning patterns
- **ROI (Region of Interest)**: Point camera at specific location
- **Gimbal Control**: Stabilized camera control
- **Payload Trigger**: Drop packages, fire sprayers on command
- **Custom Commands**: Execute arbitrary MAVLink commands

## Integration with One Mind System

**Use Case**: Autonomous vehicle and drone operations
- **Pillar**: SO (Sovereign Operations) — UAV/robotics fleet management
- **Pillar**: GE (Generational Empire) — if operating commercial drone services
- **Relevant Domains**: Robotics, autonomous systems, infrastructure monitoring

## System Requirements

- **macOS**: 10.11+ (Intel) / 11.0+ (Apple Silicon)
- **RAM**: 2 GB minimum; 4+ GB recommended
- **Disk**: 500 MB for app; 1+ GB for logs
- **Network**: Internet for map tiles; optional for vehicle operation

## Common Use Cases

- **Aerial Mapping**: Survey and orthomosaic generation
- **Agricultural Monitoring**: Crop health assessment flights
- **Infrastructure Inspection**: Building/bridge/power line surveys
- **Search & Rescue**: Coordinated multi-drone search operations
- **Environmental Monitoring**: Wildlife tracking, environmental surveys
- **Delivery Systems**: Autonomous package delivery routes

## Removal Instructions

To completely uninstall QGroundControl:

```bash
rm -rf /Applications/QGroundControl.app
rm -rf ~/Library/Application\ Support/QGroundControl
rm -rf ~/Library/Preferences/org.qgroundcontrol.QGroundControl.plist
rm -rf ~/Library/Caches/org.qgroundcontrol.QGroundControl*
rm -rf ~/Library/Saved\ Application\ State/org.qgroundcontrol.QGroundControl.savedState
```

## Alternatives

- **Mission Planner** (ArduPilot ecosystem, Windows primary)
- **MAVProxy** (Command-line based, scriptable)
- **Auterion Skynode Manager** (Commercial Auterion ecosystem)
- **DJI Flightsafe** (Proprietary DJI drones)
- **Custom Ground Stations** (Python/ROS-based alternatives)

## Notes

- 100% open-source; community-driven development
- Cross-platform (macOS, Windows, Linux, mobile)
- Supports multi-vehicle coordination
- Active development; releases align with PX4/ArduPilot updates
- Extensive documentation and community support
- Can be extended via plugins for custom functionality
- Used by research institutions, commercial operators, and hobbyists
- No mandatory cloud dependency; fully functional offline

## MAVLink Protocol

- **Communication Standard**: Universal protocol for autopilots and ground stations
- **Ubiquitous**: Supported by 100+ autopilot systems
- **Extensible**: Custom message types for specialized payloads
- **Efficient**: Binary encoding for low bandwidth

---

*Last verified: June 17, 2026 (v5.0)*
