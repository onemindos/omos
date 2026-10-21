# QIKI DTMP — Digital Twin Simulation Platform

**Language:** Python (FastStream) + Go  
**Stack:** NATS + JetStream, FastStream bridge  
**Role in OneMind:** Reference for operator TUI (ORION) + simulation truth telemetry pipeline

Digital twin simulation platform with operator TUI called ORION. Bilingual codebase (English/Russian). No-mocks policy: if data is missing, UI shows honest N/A rather than fabricated state.

## Components

- `q-sim-service` — generates world state, publishes telemetry and radar truth to NATS/JetStream
- `qiki-faststream-bridge-phase1` — normalizes and routes messages
- ORION (`operator-console`) — operator TUI rendering semantic view
- `qiki-nats-phase1` — NATS + JetStream configuration

## Data flow

```
Simulation → NATS JetStream → FastStream bridge → ORION TUI
                                                 ↑
                                    Operator commands loop back to sim
```

## Why it matters for OneMind

Direct reference for building the OneMind COP operator interface. ORION's no-mocks doctrine matches OneMind's philosophy — agents report real state or nothing. FastStream bridge pattern is reusable for OneMind's normalization layer.
