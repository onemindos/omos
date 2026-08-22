# One Mind — omos-bridge docs

## Standard Contract Reference

Full contract spec lives in the root [README.md](../README.md#the-standard-contract).

Quick reference:

### Subject Pattern
```
omos.{domain}.{entity_id}.telemetry
omos.{domain}.{entity_id}.command
omos.{domain}.{entity_id}.status
omos.bridge.{bridge_name}.health
omos.bridge.{bridge_name}.errors
```

### Domains
| Domain | Covers |
|---|---|
| `air` | Drones, UAVs, aircraft |
| `ground` | Vehicles, rovers, robots |
| `home` | Smart home, IoT, HVAC |
| `field` | Agricultural, weather |
| `sea` | Marine, buoys |
| `infra` | Industrial, grid, pipelines |
| `net` | SaaS, APIs, webhooks |
| `human` | Phones, wearables, location |
| `space` | Satellites, ground stations |

### Message Envelope
```json
{
  "id": "ULID",
  "ts": 1718123456789,
  "source": "bridge-name",
  "entity": "domain.entity-id",
  "domain": "domain",
  "payload": {}
}
```

## SDK References

- [orbit.go](https://github.com/synadia-io/orbit.go) — Go bridge SDK foundation
- [orbit.rs](https://github.com/synadia-io/orbit.rs) — Rust
- [orbit.js](https://github.com/synadia-io/orbit.js) — TypeScript
- [nats.go](https://github.com/nats-io/nats.go) — official Go client
- [nats.py](https://github.com/nats-io/nats.py) — official Python client
- [NEX](https://github.com/synadia-io/nex) — NATS Execution Engine (edge deploy)

## Architecture Overview

```
[Protocol / Device]
        ↓
[omos-bridge-{protocol}]
        ↓                         ← standard envelope on every message
[NATS Fabric — omos.{domain}.{entity}.telemetry]
        ↓
[AI Agents | Fleet UI | TAK Map | n8n | Home Assistant]
```
