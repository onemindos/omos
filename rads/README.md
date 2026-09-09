# RADS — Remote Autonomous Deployment System

> **OneMind OS · Homestead Digital Twin**
> Status: Phase 1 ✅ · Phase 2 (MAVLink bridge) 🔄

A full **CPU-only digital twin** of the homestead — drones, ground robots, sensors,
and Hermes AI agents all running in simulation, talking over NATS with
**identical message schemas to the real hardware**.

## Architecture

```
3D Homestead Model (House + Barn + Office)
        ↓
Simulation Engine (Gazebo Harmonic / PX4 SITL)
        ↓
MAVLink→NATS Bridge  +  Soil Sensor Array
        ↓
NATS JetStream (RADS_LAND · RADS_ROBOTS) — 10yr retention
        ↓
Land Agent + Legacy AI Override
        ↓
ClickHouse Time Machine (onemind DB)
```

## Quick Start

```bash
cp .env.example .env   # edit as needed
pip install -r requirements.txt

# Run all components:
python run_sensors.py       # soil sensor array (4 fields)
python run_agent.py         # land irrigation agent
python run_mavlink.py       # MAVLink→NATS bridge (Phase 2)
python run_clickhouse.py    # ClickHouse Time Machine writer

# Or via Docker Compose:
docker compose up
```

## NATS Subject Tree

```
rad.<homestead_id>.
  land.field.<n>.soil          ← sensor readings
  land.field.<n>.irrigate      ← irrigation commands
  land.decisions               ← agent decision log
  robots.<drone_id>.telemetry  ← MAVLink telemetry
  robots.<drone_id>.cmd        ← drone commands
  legacy.override              ← Legacy AI override channel
```

## Stack

| Layer | Tool |
|-------|------|
| Sim engine | Gazebo Harmonic |
| Drone firmware | PX4 SITL |
| Drone control | MAVSDK-Python |
| Mission planner | QGroundControl |
| Message bus | NATS JetStream |
| AI agents | Hermes (Legacy) |
| Time series | ClickHouse (onemind DB) |
| Telemetry dashboard | Foxglove Studio |

## Build Sequence

| # | Step | Status |
|---|------|--------|
| 1 | Soil sensor + land agent (native NATS, JetStream) | ✅ Done |
| 2 | PX4 SITL drone flying in Gazebo | ✅ Done |
| 3 | MAVLink→NATS bridge (versioned schema) | 🔄 Scaffold ready |
| 4 | ClickHouse Time Machine writer | 🔄 Scaffold ready |
| 5 | Legacy closed-loop drone commands | 🔄 Override channel wired |
| 6 | Multi-robot coordination | ⏳ |
| 7 | 3D world (barn/field meshes) | ⏳ |

## Infrastructure

- **Lightsail VM:** `onemind-rads-worldsim-01` · Ubuntu 24.04 · 8 vCPU / 30 GB
  - Public IP: `13.223.102.62` · Tailscale: `100.75.87.51`
- **Mac M2 Max:** Tailscale `100.98.178.23`
- **NATS Hub:** `100.74.182.117:4222`
- **ClickHouse:** `100.74.182.117:9001` · database `onemind`

## Docs

- [Architecture](docs/architecture.md)
- [NATS Subject Reference](docs/nats-subjects.md)
- [MAVLink Bridge](docs/mavlink-bridge.md)

## Related

- Vault: `onemind/RADS-Virtual-Homestead-Digital-Twin.md`
- One Mind pillar: **LE (Legacy Evolution)** — Home & Homestead
