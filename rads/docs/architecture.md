# RADS Architecture

## Overview

```
┌─────────────────────────────────────────────────────────────────┐
│                        RADS Digital Twin                        │
│                                                                 │
│  ┌──────────────┐    ┌──────────────┐    ┌──────────────────┐  │
│  │ Soil Sensors │    │ MAVLink Brdg │    │  ClickHouse TM   │  │
│  │  (4 fields)  │    │  PX4 SITL    │    │  10yr retention  │  │
│  └──────┬───────┘    └──────┬───────┘    └────────┬─────────┘  │
│         │                  │                      │             │
│         └──────────────────┴──────────────────────┘             │
│                            │                                    │
│                    ┌───────▼────────┐                           │
│                    │  NATS JetStream │                           │
│                    │  RADS_LAND      │                           │
│                    │  RADS_ROBOTS    │                           │
│                    └───────┬────────┘                           │
│                            │                                    │
│              ┌─────────────┴──────────────┐                     │
│              │                            │                     │
│     ┌────────▼────────┐       ┌──────────▼──────────┐          │
│     │   Land Agent    │       │   Legacy (Mother)    │          │
│     │  auto-irrigate  │◄──────│  AI override cmds    │          │
│     └─────────────────┘       └─────────────────────┘          │
└─────────────────────────────────────────────────────────────────┘
```

## NATS Subject Tree

```
rad.<homestead_id>.
  land.
    field.<n>.soil          ← sensor readings (JetStream)
    field.<n>.irrigate      ← irrigation commands
    decisions               ← agent decision log
  robots.
    <drone_id>.telemetry    ← MAVLink telemetry (JetStream)
    <drone_id>.cmd          ← drone commands
  legacy.
    override                ← Legacy AI override channel
```

## Components

| Component | File | Role |
|-----------|------|------|
| SoilSensorArray | `rads/sensors/soil.py` | Multi-field sensor sim, native NATS |
| LandAgent | `rads/agents/land_agent.py` | Irrigation decisions + Legacy hook |
| MAVLinkNATSBridge | `rads/bridges/mavlink_bridge.py` | PX4 SITL → NATS (Phase 2) |
| ClickHouseWriter | `rads/storage/clickhouse_writer.py` | Time Machine persistence |

## Phases

| Phase | Status |
|-------|--------|
| 1. Soil sensor + land agent (native NATS, JetStream) | ✅ |
| 2. MAVLink→NATS bridge (PX4 SITL) | 🔄 scaffold ready |
| 3. Legacy closed-loop commands | 🔄 override channel wired |
| 4. ClickHouse Time Machine | ✅ scaffold ready |
| 5. Multi-robot coordination | ⏳ |
| 6. 3D world (Gazebo Harmonic) | ⏳ |
