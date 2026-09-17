---
title: "Pi Agent"
type: reference
status: active
tags: [platform, ai, agent, physical, raspberry-pi, iot]
---
# Pi Agent

> Physical node agent running on Raspberry Pi hardware. Bridges the digital fleet to physical reality via sensors, cameras, and actuators.

## Use Case in OneMind OS
- Physical presence at nodes (House, Barn, Office, Greenhouse)
- Sensor data collection → NATS telemetry
- Camera/vision processing (local inference)
- Voice satellite (wake word + STT + TTS)
- GPIO control for actuators (relays, motors, servos)

## Hardware
| Node | Pi Model | Purpose |
|------|----------|---------|
| House | Pi 5 | Central hub, Frigate NVR, voice |
| Barn | Pi 4 | Greenhouse sensors, camera |
| Office | Pi 4 | Development, testing |

## Stack
- **OS:** Raspberry Pi OS (64-bit)
- **Runtime:** Python 3.11+
- **Messaging:** NATS client (nats-py)
- **Voice:** Wyoming protocol (HA compatible)
- **Vision:** OpenCV, ONNX Runtime (local inference)
- **Sensors:** GPIO, I2C, SPI (BME280, soil moisture, etc.)

## NATS Subjects
```
telemetry.{node}.temperature
telemetry.{node}.humidity
telemetry.{node}.camera.motion
cmd.{node}.relay.{id}
cmd.{node}.voice.speak
```

## Integration
- Reports to Hermes via NATS
- Home Assistant via MQTT bridge on :1883
- Fleet citizen contract (heartbeat, service discovery)

## Cost
$0 (hardware already owned)
