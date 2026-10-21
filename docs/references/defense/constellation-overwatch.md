# Constellation Overwatch — C4ISR Data Fabric

**Repo:** https://github.com/Constellation-Overwatch/constellation-overwatch
**Language:** Go
**License:** Community Copyleft
**Stars:** 34
**Topics:** c4isr, drones, nats, robotics, golang, turso

## What it is

Vendor-agnostic C4ISR Data Plane — a Go-based data fabric for drones, robots, sensors, and video feeds on the industrial edge. Built on NATS JetStream. Essentially the same problem space as OneMind OS.

## Architecture

```
Drones │ Robots │ Sensors │ Video Feeds
              ▼
       NATS JetStream (data fabric)
              ▼
    Constellation Overwatch Core
    - Entity tracking
    - Video feed management
    - Mission coordination
    - Web UI (pkg/services/web)
              ▼
         Turso (SQLite edge DB)
```

## Why it's relevant to OneMind

| Pattern | OneMind application |
|---------|---------------------|
| NATS JetStream as C4ISR fabric | Same architecture — swap subjects and it's OneMind |
| Drone + robot + sensor unified data plane | Exactly what omos-bridge + omd are building |
| MAVLink → NATS bridge | `mavlink2constellation` repo — directly usable in OneMind bridge layer |
| Entity tracking over NATS | Mirrors OneMind `ent.{type}.{id}.state` subject pattern |
| Go services on NATS | Same language, same bus — code patterns directly applicable |
| Web UI for C4ISR | Reference for omos-ui Integrations + Intel panels |

## Companion repos

| Repo | Purpose |
|------|---------|
| https://github.com/Constellation-Overwatch/mavlink2constellation | MAVLink UDP → NATS JetStream bridge — directly usable in OneMind |
| https://github.com/Constellation-Overwatch/toolbelt | IaC toolbelt for deploying the stack |

## Key differences from OneMind

- Constellation uses Turso (SQLite edge DB) — OneMind uses ClickHouse
- No TAK/CoT integration — OneMind has TAK Server + NatsBridgePlugin
- No AI agent layer — OneMind has Hermes/Legacy on the bus
- No account topology (CORE/FABRIC/INTERFACE) — OneMind has multi-tenant NATS accounts

## What to borrow

- MAVLink bridge implementation (`mavlink2constellation`)
- Entity state tracking patterns
- Go NATS service structure
- Web UI C4ISR layout patterns for omos-ui
