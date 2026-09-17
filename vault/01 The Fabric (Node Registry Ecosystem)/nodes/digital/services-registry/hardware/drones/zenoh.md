---
title: "Zenoh"
type: reference
status: active
updated: 2026-06-06
tags: [platform, infrastructure, messaging, ros2, dds, robotics, drones, edge, nats]
---
# Zenoh

> Unified data-in-motion + data-at-rest protocol. Designed for constrained/edge environments. Replaces DDS under ROS2 as the real-time transport layer. Has a native NATS bridge.

## What It Is

Zenoh (pronounced "zeno") is a pub/sub + query + storage protocol from Eclipse Foundation. It targets IoT, robotics, and edge computing where DDS is too heavy.

Key properties:
- **Wire-efficient** — minimal overhead vs DDS/RTPS
- **Peer-to-peer or brokered** — no mandatory central broker
- **Routable** — geo-distributed, multi-hop
- **Storage API** — queryable state, not just streams
- **ROS2 RMW** — drop-in DDS replacement via `rmw_zenoh_cpp`

## Zenoh vs DDS vs NATS

| Feature | Zenoh | DDS/RTPS | NATS |
|---------|-------|---------|------|
| Edge/IoT optimized | ✅ | ❌ (heavy) | ⚠️ |
| ROS2 RMW plugin | ✅ | ✅ (default) | ❌ |
| NATS bridge | ✅ | ❌ | — |
| Geo-distributed | ✅ | Limited | ✅ |
| Storage/query | ✅ | Limited | Via JetStream |
| Resource usage | Low | High | Low |

## Zenoh as ROS2 Transport

```bash
# Install rmw_zenoh
sudo apt install ros-humble-rmw-zenoh-cpp

# Use Zenoh instead of DDS
export RMW_IMPLEMENTATION=rmw_zenoh_cpp
ros2 run demo_nodes_cpp talker
```

ROS2 topics/services/actions work identically — just lower latency and less CPU on Pi/Jetson.

## Zenoh ↔ NATS Bridge

Zenoh has a native NATS backend bridge:
```
ROS2 (rmw_zenoh) ←→ Zenoh Router ←→ NATS Server
```
This makes ROS2 topics addressable from NATS subjects — the key integration point for OneMind OS.

Subject mapping: Zenoh key expression → NATS subject (configurable transform).

## Architecture in OneMind

```
NATS (OneMind OS backbone)
  ↕ Zenoh-NATS bridge
Zenoh Router (omos-cloud or edge node)
  ↕ rmw_zenoh
ROS2 nodes (Pi companion computer)
  ↕ MAVROS / px4_ros_com
Autopilot
```

## Use Cases in Fleet

- Replace DDS on Pi companion computers — cuts memory ~40%
- Tunnel ROS2 topics over NATS to cloud agent for telemetry
- Multi-robot pub/sub without a centralized ROS2 master
- Edge camera streams into NATS via Zenoh storage

## Install

```bash
# Zenoh standalone router
cargo install zenohd

# Python SDK
pip install zenoh

# ROS2 RMW
sudo apt install ros-humble-rmw-zenoh-cpp
```

## Related

- [[ros2-drones]] — uses Zenoh as transport
- [[nats]] — Zenoh bridges into NATS
- [[mavros]] — sits above transport layer
- [[px4]] — downstream of the ROS2 → Zenoh pipeline
