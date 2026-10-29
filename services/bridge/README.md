# omos-bridge

> The universal protocol bridge layer for the One Mind Fabric.

**NATS is the nervous system. omos-bridge is the peripheral network — every device, protocol, and system in the physical and digital world, connected to one fabric.**

[![CI](https://github.com/onemindos/omos-bridge/actions/workflows/ci.yml/badge.svg)](https://github.com/onemindos/omos-bridge/actions/workflows/ci.yml)
[![License](https://img.shields.io/badge/license-Apache%202.0-blue.svg)](./LICENSE)

---

## What This Is

`omos-bridge` is an open-source collection of protocol adapters that connect anything — drones, robots, cars, satellites, sensors, smart homes, SaaS apps, cameras, industrial equipment — to the One Mind NATS Fabric.

Every bridge follows the **same standard contract**: one subject pattern, one message envelope, one process interface. Plug in a new protocol — the fabric routes it. No custom consumer logic per device. No per-bridge maintenance hell.

**The stack:**

```
[Device / Protocol]
        ↓
[omos-bridge-{protocol}]     ← this repo
        ↓
[NATS Fabric — JetStream]    ← omos-nats (your hub)
        ↓
[AI Agents / Fleet UI / Automation / TAK Map]
```

---

## Quick Start

```bash
# Start a local NATS server
docker compose up -d

# Watch the entire fabric
omos fabric watch

# Run a bridge
cd bridges/mavlink
cp .env.example .env && go run .
```

---

## Repository Structure

```
omos-bridge/
├── bridge-template/          ← copy this to start a new bridge
├── bridges/
│   ├── mavlink/              ← MAVLink v2 drones/autopilots (gomavlib)
│   ├── tak/                  ← TAK / CoT XML ↔ NATS (ATAK, OTS)
│   ├── tak-client/           ← TAK Server streaming client (8089 mTLS) ↔ tak.cot.out / tak.cot.inject
│   ├── homeassistant/        ← Home Assistant WebSocket API ↔ NATS
│   ├── webhook/              ← HTTP POST → NATS (GitHub, Stripe, Twilio, etc.)
│   └── mqtt/                 ← MQTT 3.1/5.0 ↔ NATS (IoT, sensors, HA devices)
├── internal/
│   └── dashboard/            ← Bridge health dashboard (TUI + HTTP /health)
├── cmd/
│   └── omos/                 ← Developer CLI: fabric watch, bridge list, entity list
├── schemas/                  ← JSON Schema (draft 2020-12) for all message types
│   ├── envelope.json         ← Standard One Mind message envelope
│   ├── air.position.json
│   ├── air.attitude.json
│   ├── air.battery.json
│   ├── air.heartbeat.json
│   ├── air.vfrhud.json
│   ├── home.sensor.json
│   ├── home.event.json
│   ├── net.webhook.json
│   ├── field.cot.json
│   └── ground.telemetry.json
├── docs/
│   └── standard-contract.md
├── docker-compose.yml        ← One-command dev environment (NATS + nats-box)
├── LICENSE                   ← Apache 2.0
└── CONTRIBUTING.md
```

---

## Bridge Catalog

### Available

| Bridge | Protocol | Domain | Language | Status |
|---|---|---|---|---|
| [mavlink](./bridges/mavlink/) | MAVLink v2 | `air` | Go | ✅ v0.2 |
| [tak](./bridges/tak/) | TAK / CoT XML | `field` | Go | ✅ v0.1 |
| [homeassistant](./bridges/homeassistant/) | HA WebSocket API | `home` | Go | ✅ v0.1 |
| [webhook](./bridges/webhook/) | HTTP POST/GET/PUT | `net` | Go | ✅ v0.1 |
| [mqtt](./bridges/mqtt/) | MQTT 3.1 / 5.0 | `home` | Go | ✅ v0.1 |

### Planned — Tier 2 (Physical AI & Robotics)

| Bridge | Protocol | Domain | Notes |
|---|---|---|---|
| ros2 | ROS2 / DDS | `ground` | Every serious robot speaks ROS2 |
| can | CAN Bus | `ground` | Every vehicle (post-1996) has CAN |
| serial | UART / RS232 / RS485 | `infra` | Arduino, ESP32, Modbus RTU |
| grpc | gRPC / Protobuf | `net` | Modern microservices + robotics APIs |

### Planned — Tier 3 (Industrial & Infrastructure)

| Bridge | Protocol | Domain | Notes |
|---|---|---|---|
| modbus | Modbus TCP/RTU | `infra` | 60% of industrial equipment |
| opcua | OPC-UA | `infra` | Modern industrial standard |
| bacnet | BACnet | `infra` | Building automation (HVAC, lighting) |
| snmp | SNMP | `net` | Network device telemetry |
| kafka | Kafka / Confluent | `net` | Enterprise on-ramp |

### Planned — Tier 4 (Human & Mobile)

| Bridge | Protocol | Domain | Notes |
|---|---|---|---|
| location | iOS/Android background | `human` | Field team tracking → TAK map |
| wearable | BLE GATT / Apple Health | `human` | HR, HRV, sleep → AI agents |
| rtsp | RTSP camera streams | `home` | Security/drone FPV → fabric |
| webrtc | WebRTC data channels | `human` | Real-time field comms |

### Planned — Tier 5 (Aerospace & Space)

| Bridge | Protocol | Domain | Notes |
|---|---|---|---|
| adsb | ADS-B / SBS-1 | `air` | Aircraft transponder via RTL-SDR |
| nmea | NMEA 0183/2000 | `sea` | Marine GPS, AIS vessel tracking |
| ccsds | CCSDS | `space` | NASA/ESA satellite protocol |
| sdr | GNU Radio / RTL-SDR | `space` | Software-defined radio raw streams |

> Want to add a bridge? Read [CONTRIBUTING.md](./CONTRIBUTING.md).

---

## The Standard Contract

Every bridge follows the same three rules. No exceptions.

### 1. Subject Naming

```
omos.{domain}.{entity_id}.telemetry   ← device → fabric
omos.{domain}.{entity_id}.command     ← fabric → device
omos.{domain}.{entity_id}.status
omos.bridge.{bridge_name}.health      ← heartbeat every 30s
omos.bridge.{bridge_name}.errors
```

**Domains:** `air` · `ground` · `home` · `field` · `sea` · `infra` · `net` · `human` · `space`

### 2. Message Envelope

```json
{
  "id": "01HXYZ...",
  "ts": 1718123456789,
  "source": "mavlink-bridge",
  "entity": "air.drone-alpha-1",
  "domain": "air",
  "payload": {
    "msg_type": "position",
    "data": { ... }
  }
}
```

### 3. Process Contract

| Requirement | Detail |
|---|---|
| Config via env vars | No hardcoded values |
| Heartbeat | `omos.bridge.{name}.health` every 30s |
| Error reporting | `omos.bridge.{name}.errors` structured JSON |
| Graceful shutdown | SIGTERM → drain → exit |
| Reconnect | NATS `MaxReconnects(-1)` — never crash on disconnect |
| NEX deploy | `bridge.yaml` Nexfile in every bridge |

Full spec: [docs/standard-contract.md](./docs/standard-contract.md)

---

## Schema Registry

All message payload shapes are defined as JSON Schema (draft 2020-12) in `/schemas/`.

```
schemas/
├── envelope.json          ← the outer envelope (all bridges)
├── air.position.json      ← MAVLink GLOBAL_POSITION_INT
├── air.attitude.json      ← MAVLink ATTITUDE
├── air.battery.json       ← MAVLink BATTERY_STATUS
├── air.heartbeat.json     ← MAVLink HEARTBEAT
├── air.vfrhud.json        ← MAVLink VFR_HUD
├── home.sensor.json       ← HA state update
├── home.event.json        ← HA state_changed event
├── net.webhook.json       ← inbound HTTP webhook
├── field.cot.json         ← TAK / CoT XML
└── ground.telemetry.json  ← generic ground vehicle
```

---

## omos CLI

```bash
# Install
cd cmd/omos && go install .

# Commands
omos bridge list                    # Bridge catalog
omos bridge health                  # Live bridge health (35s window)
omos fabric watch                   # Watch all fabric messages
omos fabric watch --domain air      # Watch air domain only
omos entity list                    # Discover entities (5s window)
```

---

## Dashboard

```bash
cd internal/dashboard && go run .
# Terminal table refreshes every 10s
# HTTP: GET http://localhost:8099/health
```

```
=== omos-bridge Dashboard ===  2026-06-17 14:32:00
BRIDGE                  STATUS    LAST HB        ERRORS  MSGS
mavlink-bridge          ONLINE    2s ago         0       142
tak-bridge              ONLINE    28s ago        0       89
homeassistant-bridge    ONLINE    15s ago        0       445
webhook-bridge          ONLINE    8s ago         0       23
mqtt-bridge             ONLINE    22s ago        0       312

DOMAIN   MSGS
air      142
field    89
home     757
net      23
Total: 1011 messages
```

---

## Development

```bash
# One-command dev environment
docker compose up -d       # NATS + nats-box

# Build everything
for d in bridge-template bridges/* internal/dashboard cmd/omos; do
  (cd $d && go build ./...)
done

# Watch the fabric while developing
omos fabric watch
```

---

## NEX Integration

Every bridge ships a `bridge.yaml` Nexfile for deployment via [NEX](https://github.com/synadia-io/nex) — the NATS Execution Engine. Run bridges directly on NATS leaf nodes at the edge.

```bash
nex node up
nex deploy bridges/mavlink/bridge.yaml
```

---

## SDK Foundation: Synadia Orbit

Go bridges are built on [`orbit.go`](https://github.com/synadia-io/orbit.go) (where applicable) for partitioned consumers, batch publish, scheduled messages, and KV codecs on top of the official NATS Go client.

---

## Architecture

```
┌──────────────────────────────────────────────────────────────────┐
│                      omos-bridge ecosystem                        │
├─────────────────┬──────────────────┬─────────────────────────────┤
│  Physical Layer │  Digital Layer   │  Human Layer                 │
│  Drones/MAVLink │  SaaS/Webhooks   │  Phones / Wearables          │
│  Robots/ROS2    │  MQTT brokers    │  TAK/ATAK field teams        │
│  Vehicles/CAN   │  HTTP APIs       │  Location tracking           │
│  Sensors/Modbus │  Kafka streams   │  WebRTC voice/video          │
│  Satellites     │  HA WebSocket    │  Apple Health                │
└────────┬────────┴────────┬─────────┴────────────┬────────────────┘
         └─────────────────┼──────────────────────┘
                           ↓
         ┌─────────────────────────────────────────┐
         │         NATS Fabric (JetStream)          │
         │   omos.{domain}.{entity}.telemetry       │
         │   omos.{domain}.{entity}.command         │
         │   omos.{domain}.{entity}.status          │
         └─────────────────────────────────────────┘
                           ↓
         ┌─────────────────────────────────────────┐
         │             Consumers                    │
         │  AI Agents · Fleet UI · TAK Map · n8n   │
         │  Home Assistant · Grafana · Your App     │
         └─────────────────────────────────────────┘
```

---

## Relation to One Mind OS

```
One Mind OS (Brain)
      ↕
NATS Fabric (Spinal Cord)
      ↕
omos-bridge (Peripheral Nervous System)
      ↕
Everything in the world
```

---

## Roadmap

### ✅ Scaffold & Infrastructure
- [x] Bridge template (`bridge-template/`)
- [x] Standard contract — subject pattern, message envelope, process interface
- [x] Schema registry — 11 JSON Schema files (`schemas/`)
- [x] `docker-compose.yml` — one-command dev environment
- [x] GitHub Actions CI — build, vet, test, schema validate, license check
- [x] `LICENSE` — Apache 2.0
- [x] `internal/dashboard` — bridge health TUI + HTTP `/health`
- [x] `cmd/omos` — developer CLI (fabric watch, bridge list, entity list)

### ✅ Tier 1 — Live Stack (shipped)
- [x] `bridges/mavlink` — MAVLink v2 drones/autopilots (gomavlib, UDP)
- [x] `bridges/tak` — TAK / CoT XML ↔ NATS (ATAK, OTS, UDP)
- [x] `bridges/homeassistant` — Home Assistant WebSocket API ↔ NATS
- [x] `bridges/webhook` — HTTP POST/GET/PUT → NATS (GitHub, Stripe, Twilio, etc.)
- [x] `bridges/mqtt` — MQTT 3.1/5.0 ↔ NATS bidirectional

### 🔧 Tier 2 — Physical AI & Robotics
- [ ] `bridges/ros2` — ROS2 / DDS (every serious robot)
- [ ] `bridges/can` — CAN Bus / OBD-II (every vehicle post-1996)
- [ ] `bridges/serial` — UART / RS232 / RS485 (Arduino, ESP32, Modbus RTU)
- [ ] `bridges/grpc` — gRPC / Protobuf (modern microservices + robotics APIs)

### 🏭 Tier 3 — Industrial & Infrastructure
- [ ] `bridges/modbus` — Modbus TCP/RTU (60% of industrial equipment)
- [ ] `bridges/opcua` — OPC-UA (modern industrial standard, EU mandated)
- [ ] `bridges/bacnet` — BACnet (HVAC, lighting, access control)
- [ ] `bridges/snmp` — SNMP (network device telemetry)
- [ ] `bridges/kafka` — Kafka / Confluent (enterprise on-ramp)

### 👤 Tier 4 — Human & Mobile Layer
- [ ] `bridges/location` — iOS/Android background location → fabric → TAK map
- [ ] `bridges/wearable` — BLE GATT / Apple Health / Garmin (HR, HRV, sleep)
- [ ] `bridges/rtsp` — RTSP camera streams → fabric (security, drone FPV)
- [ ] `bridges/webrtc` — WebRTC data channels (real-time field comms)

### 🛰 Tier 5 — Aerospace & Space
- [ ] `bridges/adsb` — ADS-B / SBS-1 via RTL-SDR (aircraft transponders)
- [ ] `bridges/nmea` — NMEA 0183/2000 (marine GPS, AIS vessel tracking)
- [ ] `bridges/ccsds` — CCSDS (NASA/ESA satellite protocol)
- [ ] `bridges/sdr` — GNU Radio / RTL-SDR raw streams

### 🔩 Platform (ongoing)
- [ ] `orbit.go` integration — partitioned consumers, batch publish in all bridges
- [ ] `cmd/omos` v0.2 — `bridge run`, `bridge test`, synthetic frame injection
- [ ] Schema versioning — `v1`/`v2` migration support
- [ ] NEX deployment docs — run bridges at edge nodes via NATS Execution Engine
- [ ] Grafana dashboard template — plug fabric metrics into existing observability stack

> Want to claim a bridge? Open an issue. Read [CONTRIBUTING.md](./CONTRIBUTING.md) first — the one rule is the standard contract.

---

## Contributing

Read [CONTRIBUTING.md](./CONTRIBUTING.md). The one rule: implement the standard contract.

---

## License

Apache 2.0 — same as NATS. Open forever.

Copyright 2026 Zeus Delacruz / One Mind OS

---

> *"The internet routes packets. One Mind routes meaning."*
