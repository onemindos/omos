---
title: "Drones — NATS Integration Architecture"
type: reference
status: planned
tags: [platform, hardware, drones, nats, fleet, architecture]
---
# Drones — NATS Integration Architecture

> NATS as the fleet fabric for drone operations. Real pattern, real use cases.

## How It Works

NATS isn't on the drone's flight controller — it lives on the companion computer bridging MAVLink to the fleet.

```
[PX4/ArduPilot] ←MAVLink 100Hz→ [Companion Pi]
                                      │
                              DroneKit/MAVSDK
                                      │
                              NATS Leaf Node (1-10Hz)
                                      │
                    ┌─────────────────┼──────────────────┐
                    │                 │                   │
              [NATS Central]    [JetStream]         [Hermes]
                    │                 │                   │
              [QGC/ATAK]     [Telemetry Store]    [Fleet Cmds]
```

## Real Use Cases

| Company/Project | What They Do |
|-----------------|-------------|
| **Synadia (NATS creators)** | Published drone fleet management demos — NATS as C2 backbone |
| **DroneUp / Zipline** | Fleet logistics — pub/sub architectures identical to NATS |
| **US DoD / DARPA** | Swarm programs (OFFSET) use pub/sub for multi-agent coordination |
| **Auterion** | PX4 enterprise — cloud platform uses event-driven messaging for fleet telemetry |
| **NVIDIA Isaac + drones** | Jetson-based autonomy publishes to message buses for fleet awareness |

## Why NATS Works for Drones

| Need | NATS Feature |
|------|-------------|
| Low latency telemetry | Core pub/sub (fire-and-forget, microsecond routing) |
| Command reliability | JetStream (at-least-once delivery for critical commands) |
| Fleet coordination | Subject hierarchy (`telemetry.drone.recon.gps`, `cmd.drone.*.rtl`) |
| Disconnection handling | JetStream catches up missed messages on reconnect |
| Edge → Cloud | Leaf nodes (drone Pi connects as leaf to central NATS) |
| Multi-consumer | Ground station + Hermes + ATAK all subscribe independently |

## NATS Subject Hierarchy

```
telemetry.drone.{name}.gps
telemetry.drone.{name}.battery
telemetry.drone.{name}.attitude
telemetry.drone.{name}.camera
telemetry.drone.{name}.signal
cmd.drone.{name}.goto
cmd.drone.{name}.rtl
cmd.drone.{name}.land
cmd.drone.{name}.arm
cmd.drone.*.rtl              ← recall ALL drones
event.drone.{name}.geofence
event.drone.{name}.low_battery
event.drone.{name}.mission_complete
```

## Limitations (Honest)

- **Not for the flight-critical loop** — MAVLink runs at 100Hz between autopilot ↔ companion. NATS is not in that path.
- **Fleet-level awareness only** — telemetry reporting (1-10Hz), commands from ground, coordination between drones, logging to JetStream.
- **Connectivity** — drones lose signal. NATS leaf nodes handle reconnection gracefully, but local autonomy (ArduPilot/PX4) must operate independently when disconnected.

## OneMind Integration

```
Drone Fleet → NATS → Hermes (orchestrator)
                  → JetStream (telemetry store → Parquet/DuckDB)
                  → ATAK (tactical map overlay)
                  → Home Assistant (geofence triggers)
                  → Flutter App (fleet dashboard)
```

## Hardware Stack (Per Drone)

```
Flight Controller:  Pixhawk 6C (PX4 or ArduPilot)
Companion:          Raspberry Pi 4/CM4 or Jetson Nano
Video:              OpenHD or RTSP camera
Comms:              WiFi (short) / 4G-LTE (long) / Meshtastic (fallback)
Software:           MAVSDK-Python + nats-py + local autonomy
```
