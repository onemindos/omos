# mavlink-bridge

Connects MAVLink v2 drones and autopilots to the One Mind NATS Fabric.

Powered by [gomavlib](https://github.com/bluenviron/gomavlib) — parses real MAVLink frames and publishes decoded telemetry to `omos.air.{entity_id}.telemetry`. Subscribes to `omos.air.{entity_id}.command` for passthrough back to the vehicle.

Compatible with PX4, ArduPilot, and any MAVLink v2 autopilot (QGroundControl, Mission Planner, MAVSDK).

**Status:** v0.2 — full message parsing live.

---

## Parsed Message Types

| MAVLink Message | Subject suffix | Payload fields |
|---|---|---|
| `HEARTBEAT` | `.telemetry` (msg_type: `heartbeat`) | type, autopilot, base_mode, system_status |
| `GLOBAL_POSITION_INT` | `.telemetry` (msg_type: `position`) | lat, lon, alt_m, relative_alt_m, vx/vy/vz, heading_deg |
| `ATTITUDE` | `.telemetry` (msg_type: `attitude`) | roll/pitch/yaw (rad), roll/pitch/yaw speed |
| `BATTERY_STATUS` | `.telemetry` (msg_type: `battery`) | cell voltages, current_a, remaining_pct, temp_c |
| `VFR_HUD` | `.telemetry` (msg_type: `vfrhud`) | airspeed, groundspeed, heading, throttle_pct, alt, climb |

All messages share the standard One Mind envelope:
```json
{
  "id": "...",
  "ts": 1718123456789,
  "source": "mavlink-bridge",
  "entity": "air.drone-alpha-1",
  "domain": "air",
  "payload": {
    "msg_type": "position",
    "data": { "lat": 30.267, "lon": -97.743, "alt_m": 120.5, ... }
  }
}
```

---

## Setup

```bash
cp .env.example .env
# Edit .env — set DRONE_IP at minimum

go build -o omos-bridge-mavlink .
./omos-bridge-mavlink
```

## Environment Variables

| Variable | Required | Default | Description |
|---|---|---|---|
| `NATS_URL` | No | `nats://localhost:4222` | NATS server — comma-separate for cluster |
| `BRIDGE_NAME` | No | `mavlink-bridge` | Bridge identity on the fabric |
| `BRIDGE_ENTITY_ID` | No | `drone-alpha-1` | Entity ID — becomes `air.{id}` |
| `DRONE_IP` | **Yes** | — | Drone/GCS IP address |
| `DRONE_PORT` | No | `14550` | MAVLink UDP port |

## Subjects

```
omos.air.{entity_id}.telemetry    ← decoded MAVLink frames (all msg types)
omos.air.{entity_id}.command      ← inbound commands (passthrough v0.3)
omos.bridge.mavlink-bridge.health ← heartbeat every 30s
omos.bridge.mavlink-bridge.errors ← structured error reporting
```

## Watch it

```bash
# All telemetry
nats sub "omos.air.drone-alpha-1.telemetry"

# Bridge heartbeat
nats sub "omos.bridge.mavlink-bridge.health"

# Errors
nats sub "omos.bridge.mavlink-bridge.errors"
```

## Filter by message type (JetStream consumer)

```bash
# Position only
nats sub "omos.air.drone-alpha-1.telemetry" --filter-subject "omos.air.drone-alpha-1.telemetry"
# Then filter msg_type in your consumer
```

---

## Roadmap

- [x] v0.1 — raw UDP frame capture
- [x] v0.2 — full MAVLink parsing via gomavlib (HEARTBEAT, GLOBAL_POSITION_INT, ATTITUDE, BATTERY_STATUS, VFR_HUD)
- [ ] v0.3 — command passthrough (NATS → MAVLink COMMAND_LONG → vehicle)
- [ ] v0.3 — TAK/CoT position bridge (lat/lon → CoT XML → OTS via omos-bridge-tak)
- [ ] v0.4 — multi-drone fleet (one process per drone via env var, supervisor config)
- [ ] v0.4 — orbit.go integration (batch publish, partitioned consumers for fleet scale)
