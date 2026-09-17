---
title: "QGroundControl"
type: reference
status: active
updated: 2026-06-06
tags: [platform, hardware, drones, ground-station, gcs, gui, mavlink, px4, ardupilot]
---
# QGroundControl

> Cross-platform ground control station (GCS). Mission planning, live telemetry, firmware upload, parameter tuning. Works with both ArduPilot and PX4 via MAVLink.

## What It Does

QGroundControl (QGC) is the primary human-in-the-loop interface for drone operations:
- **Mission Planner** — draw waypoint routes on map, upload to autopilot
- **Live Telemetry** — attitude, position, battery, GPS, mode, RC signal
- **Vehicle Setup** — calibrate sensors, configure params, flash firmware
- **Video Stream** — RTSP/UDP video overlay on map
- **Log Analysis** — download and review .bin / .ulg flight logs

## Platforms

- macOS, Windows, Linux, Android, iOS
- Download: https://qgroundcontrol.com

## Connection

```
Launch QGC → Settings → Comm Links → Add
  Type: UDP
  Port: 14550 (standard MAVLink GCS port)
  Auto-connect on startup ✅
```

For hardware via USB: QGC auto-detects Pixhawk on most systems.

## Mission Planning

1. Click **Plan** view
2. Click map to place waypoints
3. Set altitude per waypoint (AGL or AMSL)
4. Add survey grid, corridor scan, or structure scan patterns
5. **Upload** mission to vehicle
6. Switch to **Fly** view, arm, start mission

## Key Panels (Fly View)

| Panel | Data |
|-------|------|
| Instrument panel | Speed, altitude, heading, vertical speed |
| GPS indicator | Fix type, satellite count, HDOP |
| Flight mode | Current mode + arm state |
| Battery | Voltage, percentage, current |
| RC signal | Link quality |
| Messages | MAVLink status messages |

## Parameter Tuning

- **Vehicle Setup → Parameters** — searchable param tree
- Filter by component, subsystem, or keyword
- Change values live (some require reboot)
- Download/upload full param files (.params)

## Firmware Flash

- **Vehicle Setup → Firmware** — flash PX4 or ArduPilot over USB
- Stable, Beta, or Custom build
- Sensor calibration wizard follows

## SITL Connection

```bash
# Run PX4/ArduPilot SITL → QGC auto-connects via UDP 14550
# No extra config needed
```

## Video Integration

- Add video stream: **Settings → Video → RTSP URL**
- Displays in fly view overlay
- Works with any RTSP source (IP cameras, companion computer streams)

## Alternatives

| Tool | Best For |
|------|---------|
| **QGroundControl** | Cross-platform, PX4 + ArduPilot, clean UI |
| **Mission Planner** | Windows, ArduPilot-specific, most features |
| **APM Planner 2** | Linux/Mac ArduPilot (older) |
| **Foxglove Studio** | ROS2 visualization + debugging (not a GCS) |

## OneMind OS Angle

QGC is the manual oversight layer when Legacy agents are in autonomous flight mode. Pattern:
- NATS-driven missions run autonomously
- QGC open on tablet as safety monitor — override/abort capability
- Live telemetry feeds into OneMind dashboard via MAVLink → NATS bridge

## Related

- [[mavlink]] — protocol QGC speaks
- [[px4]] — primary autopilot
- [[ardupilot]] — also fully supported
- [[foxglove-studio]] — complementary ROS2 viz tool
- [[ros2-drones]] — autonomy stack running above the GCS layer
