---
title: "Blue Iris"
type: reference
status: planned
tags: [platform, iot, surveillance, nvr, windows]
---
# Blue Iris

> Professional NVR software for Windows. Supports 100+ cameras, AI alerts, robust recording.

## Key Features
| Feature | Notes |
|---------|-------|
| Camera support | 100+ simultaneous, any ONVIF/RTSP |
| AI integration | DeepStack, CodeProject.AI, Frigate |
| Recording | Continuous + motion-triggered |
| Alerts | Push, email, MQTT, webhook |
| Web UI | Remote access via browser/app |
| MQTT | Publishes events for HA/NATS integration |

## License
$70 one-time (no subscription)

## Integration
- MQTT → NATS bridge for fleet alerts
- Home Assistant via MQTT binary sensors
- DeepStack/CodeProject.AI for local object detection
