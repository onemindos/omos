---
title: "OpenHD"
type: reference
status: planned
tags: [platform, hardware, drones, video, fpv, digital-link]
---
# OpenHD

> Open-source HD digital video link for drones. Low-latency FPV over WiFi. Alternative to DJI O3/Walksnail.

## What It Is
Turns Raspberry Pi + WiFi adapters into a digital HD video transmitter. Ground station receives on another Pi or PC. Low latency, long range, fully customizable.

## Key Features
| Feature | Notes |
|---------|-------|
| Latency | 80-120ms (vs 30ms DJI, but open/hackable) |
| Resolution | Up to 1080p60 |
| Range | 5-20km (depending on antenna/power) |
| Telemetry | MAVLink OSD overlay |
| Recording | DVR on ground + air |
| Bidirectional | Control link + video on same hardware |
| Hardware | Raspberry Pi + RTL8812AU WiFi adapters |

## Hardware Stack
```
AIR UNIT:
  Raspberry Pi 4/CM4 + USB camera + WiFi adapter + antenna

GROUND UNIT:
  Raspberry Pi 4 + WiFi adapter + antenna + HDMI to goggles
```

## Integration
- MAVLink telemetry overlay (GPS, battery, altitude, speed)
- RTSP output for recording/streaming
- Can forward video to Frigate for AI analysis
- NATS: stream video metadata/events

## vs Proprietary
| | OpenHD | DJI O3 | HDZero |
|---|---|---|---|
| Cost | $100 | $180 | $150 |
| Latency | 80-120ms | 30ms | 25ms |
| Range | 5-20km | 13km | 5km |
| Hackable | Full | No | Partial |
| Resolution | 1080p | 1080p | 1080p90 |

## Resources
- https://openhdfpv.org
- https://github.com/OpenHD/OpenHD
