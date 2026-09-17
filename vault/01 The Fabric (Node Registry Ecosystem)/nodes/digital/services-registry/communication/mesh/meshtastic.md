---
title: "Meshtastic"
type: reference
status: planned
tags: [platform, communication, mesh, radio, off-grid, lora]
---
# Meshtastic

> Open-source LoRa mesh networking. Off-grid encrypted communication with no internet, no cell towers, no infrastructure dependency.

## Use Case in OneMind OS
- Off-grid comms between nodes (house ↔ barn ↔ greenhouse)
- Disaster/fallback communication (no internet needed)
- Sensor telemetry over long range (LoRa reaches 1-10+ km)
- ATAK integration for tactical mapping
- Sovereign mesh — zero dependency on any provider

## Key Concepts
| Concept | Description |
|---------|-------------|
| LoRa | Long Range radio (868/915 MHz, low power) |
| Mesh | Messages hop node-to-node, self-healing |
| Channels | Encrypted groups (PSK-based) |
| MQTT | Bridge mesh to internet (optional) |
| Traceroute | See which nodes relayed your message |

## Hardware Options
| Device | Range | Notes |
|--------|-------|-------|
| Heltec V3 | ~5km | Cheap ($20), ESP32-S3, OLED |
| RAK WisBlock | ~10km | Modular, solar-friendly |
| LilyGo T-Beam | ~8km | GPS built-in, good for mobile |
| Station G2 | ~15km+ | Base station, high power |

## Software
```bash
# Python CLI
pip install meshtastic
meshtastic --info                 # Node info
meshtastic --sendtext "hello"    # Send message
meshtastic --port /dev/ttyUSB0   # Serial connection
meshtastic --set lora.region US  # Set region
```

## Integration Points
- MQTT bridge → NATS (mesh telemetry into fleet)
- ATAK plugin (tactical awareness)
- Python API for automation (sensor polling, alerts)
- Home Assistant integration (presence, geofence)

## Network Topology (Planned)
```
[House Node] ←→ [Barn Node] ←→ [Greenhouse Node]
      ↕                              ↕
[Office Node]                  [Mobile (Go2/Drone)]
```
