---
title: "IoT Hardware Roadmap — Edge Nodes & Vision AI"
type: reference
domain: SO/03 Infrastructure
topic: hardware, iot, edge-ai, robotics
status: reference
created: 2026-05-24
updated: 2026-05-26
tags: [type/reference, domain/so, topic/hardware, topic/iot, topic/vision]
---

# IoT Hardware Roadmap — Edge Nodes & Vision AI

Hardware evaluation for OneMind OS physical fleet expansion.

---

## Decision: Flock Safety Cameras

**Skip.** API is law enforcement only — no civilian access. Build sovereign LPR stack instead.

See `SO/08 Research/lpr-alpr-research.md` for full DIY LPR architecture.

---

## Node Type 1 — Raspberry Pi Zero 2 W (Sensor Nodes)

**Role:** Lightweight IoT sensor nodes — GPIO, sensors, messaging.
**NOT for vision AI** (512 MB RAM is insufficient for YOLOv8).

| Spec | Value |
|------|-------|
| CPU | 1 GHz quad-core ARM Cortex-A53 |
| RAM | 512 MB |
| Connectivity | Wi-Fi 2.4 GHz only, Bluetooth |
| GPIO | 40-pin header |
| Power | 5V micro-USB (low draw) |
| Price | ~$15 |

### Use Cases
- Temperature/humidity monitoring (DHT22)
- Motion detection (PIR sensor)
- Door/window sensors (reed switches)
- Drone telemetry (GPS, IMU, altimeter)
- NATS edge client (publish sensor data, subscribe to commands)

### Architecture
```
Pi Zero → Sensor data → MQTT/NATS → Mac M2 Max → Home Assistant
```

### Order List (~$175 for 5 nodes)
```
5x Pi Zero 2 W           $75
5x microSD 32GB          $30
5x power supplies        $20
5x DHT22 (temp/humidity) $25
5x PIR (motion)          $15
5x reed switches         $10
```

---

## Node Type 2 — Raspberry Pi 5 + AI HAT 2 (Edge Vision Nodes)

**Role:** On-device YOLOv8 inference. Deploy at cameras, drones, entry points.

| Spec | Value |
|------|-------|
| NPU | Hailo-8L AI accelerator (13 TOPS) |
| Inference | 30+ FPS (YOLOv8n) |
| Power | ~3W typical |
| Compatibility | Raspberry Pi 5 only |
| Price | ~$70–100 (HAT) + $80 (Pi 5 8GB) |
| Status | Q2 2026 — check availability |

### vs. Mac M2 Max

| | Pi 5 + AI HAT 2 | Mac M2 Max |
|-|-----------------|------------|
| Speed | 30 FPS YOLOv8n | 100 FPS YOLOv8n |
| Power | 3W | 15–30W |
| Cost | ~$195/node | Already owned |
| Role | Edge (distributed) | Central hub |

### Use Cases
- Security camera node (person/vehicle detection on-device)
- Smart doorbell (face/package/visitor detection)
- Drone onboard vision (object avoidance, landing pad)
- Home Assistant room occupancy (no cloud)

### Architecture
```
Pi 5 + AI HAT 2 → YOLOv8 (on-device) → NATS → HA Bridge → Home Assistant
```

### Order List (~$195/node, recommend 5 = $975)
```
Pi 5 8GB             $80
AI HAT 2             $70–100
Pi Camera Module v3  $25
Power supply 5V/5A   $10
microSD 64GB         $10
```

---

## Strategy: Hybrid Architecture

Mac M2 Max = central AI hub (high throughput, always-on)
Pi 5 + AI HAT 2 = edge vision nodes (distributed, low power, near cameras/drones)
Pi Zero 2 W = sensor nodes (ultra-low power, GPIO, telemetry)

```
Sensors (Pi Zero) ─────────────────────┐
Cameras (Pi 5 + HAT) ── NATS ──── Mac M2 Max ── Home Assistant
Drones (telemetry) ────────────────────┘
```

---

## World Labs / Marble (Exploratory)

**URL:** https://marble.worldlabs.ai
**API:** https://platform.worldlabs.ai
**What:** Text/image → 3D world generation (Spatial AI, 3D Gaussian Splatting)

Potential uses: home renovation visualization, drone simulation environments, digital twin of property, IoT device placement planning.

**Status:** Not critical path. Cloud-based (not sovereign). Test free tier before committing.

---

## Action Items

- [ ] Order Pi Zero 2 W starter kit when sensor deployment begins (~$175)
- [ ] Check Pi 5 + AI HAT 2 availability (Q2 2026 release)
- [ ] Design NATS event schema for sensor data (`events.iot.*`)
- [ ] Plan camera placement (driveway, entry points) before hardware order
- [ ] Test World Labs Marble (low priority, exploratory only)

---

## Related

- `SO/08 Research/lpr-alpr-research.md` — LPR/ALPR stack for camera nodes
- `SO/08 Research/obs-automation-research.md` — OBS for camera feeds
- `SO/03 Infrastructure/network/` — NATS fabric, Tailscale/NetBird mesh
- Hermes skill: `iot-fleet` — Cognitum Seed device fleet management
