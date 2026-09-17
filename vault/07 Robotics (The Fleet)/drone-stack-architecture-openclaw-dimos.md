---
domain: 04
name: Drone Stack Architecture — OpenClaw + DimOS Integration
status: reference
owner: zeus
tags: [robotics, drones, openclaw, dimos, architecture, latency, fleet, research]
created: 2026-05-12
updated: 2026-05-13
---

# Drone Stack Architecture — OpenClaw + DimOS Integration

How OpenClaw, DimOS, and hardware interconnect for drone operations at scale.

---

## 1. The Latency Stack — Why Layers Matter

OpenClaw cannot and should not be in the real-time control loop. Each layer operates at a different frequency with different latency tolerance.

| Layer | Frequency | Latency | Who Handles It |
| --- | --- | --- | --- |
| Motor/stabilize | 100-400 Hz | 2-10ms | Flight controller (hardware) |
| Obstacle avoid | 10-30 Hz | 33-100ms | DimOS (local, no LLM) |
| Sensor fusion | 2-10 Hz | 100-500ms | DimOS (local VLM/PID) |
| LLM decisions | 0.5-2 Hz | 500ms-3s | DimOS Cortex (local/cloud LLM) |
| Mission commands | ~0.01 Hz | 1-10s OK | OpenClaw (VPS/cloud) |
| Fleet coordination | minutes | 30s+ OK | OpenClaw (VPS) |

OpenClaw lives at the top two layers. It sends **mission-level commands**, not motor commands. When Zeus says "patrol zone A" in Telegram, that's a mission command. DimOS handles the 10,000 decisions per second it takes to actually fly there without crashing.

---

## 2. What "Real-Time" Means for Each Layer

| Zeus Says in Telegram | OpenClaw Does | DimOS Does (Autonomously) |
| --- | --- | --- |
| "Patrol the north field" | Sends `fly_to` waypoint sequence via MCP | Plans route, avoids obstacles, manages altitude, handles wind, streams video |
| "Follow that red truck" | Sends `follow_object("red truck", 300)` via MCP | VLM detection → CSRT tracker → PID controller → velocity commands → MAVLink, all at 10-30 Hz |
| "Emergency land all drones" | Broadcasts `land()` to all DimOS instances | Controlled descent, motor disarm, telemetry report back |
| "What's drone 3 seeing?" | Calls `observe()` via MCP, returns frame | Captures camera frame, returns to agent |

OpenClaw never touches motor commands. It's the general giving orders, not the pilot's hands on the stick.

---

## 3. Wiring Options

### Option A: Tailscane Mesh (Current Infra)

```
Zeus (Telegram/Mac)
    │ message (~200ms)
    ▼
OpenClaw (omos-core VPS, 100.98.25.10)
    │ MCP tool call over Tailscale (~20-50ms)
    ▼
DimOS Host (Mac or Jetson, on Tailscale)
    │ localhost MCP server (0ms)
    │ MAVLink UDP 14550 (~5ms WiFi)
    ▼
RosettaDrone (Android) → DJI Drone
```

Total mission command latency: ~300ms — fine for "go to waypoint X."
DimOS internal loop: 2-30 Hz locally — no network involved, pure local processing.

### Option B: Full Local (No VPS in the Loop)

For sub-second response for everything:

```
Zeus (Mac, same network as drone)
    │ OpenClaw agent running LOCALLY (not VPS)
    │ or: DimOS humancli → direct commands
    ▼
DimOS (same Mac, localhost:9990 MCP)
    │ MAVLink UDP (~5ms)
    ▼
RosettaDrone → DJI Drone
```

Can run a local OpenClaw gateway on the same machine as DimOS. Or skip OpenClaw entirely for real-time and use DimOS CLI:

```bash
# Direct real-time commands (no OpenClaw needed)
dimos agent-send "fly to 37.7749, -122.4194, altitude 50"
dimos mcp call fly_to --arg lat=37.7749 --arg lon=-122.4194 --arg alt=50
dimos mcp call follow_object --arg description="red truck" --arg duration=120
```

### Option C: Hybrid (Recommended)

```
┌─────────────────────────────────────────────────┐
│ VPS (omos-core)                                  │
│                                                   │
│ OpenClaw Gateway                                  │
│   ├─ Legacy (orchestrator)                       │
│   └─ EagleEye (drone agent)                     │
│       └─ MCP client → calls DimOS MCP servers   │
│                                                   │
│ Responsibilities:                                │
│   • Mission planning                             │
│   • Fleet status dashboard                       │
│   • Multi-drone coordination                    │
│   • Telegram interface for Zeus                  │
│   • Approval workflows                          │
│   • Logging & replay                            │
└───────────┬─────────────────────────────────────┘
            │ Tailscale (~20-50ms)
            ▼
┌─────────────────────────────────────────────────┐
│ Field Machine (Mac/Jetson/NUC, on Tailscale)    │
│                                                   │
│ DimOS Instance per Drone                         │
│   ├─ Drone 1: localhost:9990 (MCP server)       │
│   ├─ Drone 2: localhost:9991 (MCP server)       │
│   └─ Drone 3: localhost:9992 (MCP server)       │
│                                                   │
│ Responsibilities:                                │
│   • ALL real-time flight control (10-400 Hz)    │
│   • Obstacle avoidance (local, no LLM)          │
│   • Visual servoing & tracking (local VLM)      │
│   • MAVLink protocol                            │
│   • Autonomous waypoint following               │
│   • Emergency behaviors (RTH, land)             │
│                                                   │
│ Each DimOS runs AUTONOMOUSLY once given mission  │
│ OpenClaw only sends new missions or overrides    │
└───────────┬─────────────────────────────────────┘
            │ WiFi/UDP (~5ms)
            ▼
┌─────────────────────────────────────────────────┐
│ RosettaDrone (Android per drone)                 │
│   → MAVLink bridge → DJI Drone hardware         │
└─────────────────────────────────────────────────┘
```

---

## 4. Autonomy Levels

DimOS runs autonomously once given a mission. OpenClaw doesn't need to send commands every 100ms. It sends a mission, then monitors.

| Level | What OpenClaw Does | What DimOS Does | Latency Requirement |
| --- | --- | --- | --- |
| **L0 — Manual** | Every command explicit | Execute single commands | ~1s (fine) |
| **L1 — Assisted** | Set waypoints, DimOS flies | Autonomous nav + avoidance | Mission: seconds. Flight: local |
| **L2 — Supervised** | Monitor + intervene if needed | Full autonomous patrol | Mission: minutes. Override: ~1s |
| **L3 — Autonomous** | Review logs after the fact | Full autonomous everything | Reports: minutes/hours |

For patrol mode at hertz 0.5 — that's one cognitive cycle every 2 seconds inside DimOS locally. OpenClaw checks in every 30-60 seconds to see if patrol is on track.

---

## 5. What OpenClaw Actually Does in the Drone Stack

```python
# EagleEye agent receives from Zeus via Telegram:
# "Start patrol pattern Alpha on drone 1 and 2"

# EagleEye calls DimOS MCP servers:
await mcp_call("drone1:9990", "fly_to", {lat: 37.77, lon: -122.41, alt: 50})
await mcp_call("drone2:9991", "fly_to", {lat: 37.78, lon: -122.42, alt: 50})

# Then monitors (every 30s polling or event-driven):
status1 = await mcp_call("drone1:9990", "observe")  # camera frame
status2 = await mcp_call("drone2:9991", "is_flying_to_target")

# Reports back to Zeus via Telegram:
# "Drone 1: On station, visual clear. Drone 2: En route, ETA 45s."

# Zeus: "Drone 1, follow that vehicle"
await mcp_call("drone1:9990", "follow_object",
    {description: "white pickup truck", duration: 300})
```

---

## 6. The Hybrid OM1 + DimOS Architecture (Concepts to Adopt)

Cherry-pick OM1's architecture patterns and apply on top of DimOS drone runtime, with OpenClaw as the command layer.

### Patterns worth adopting from OM1:

**Dual LLM (local + cloud race) for flight decisions:**
DimOS currently uses single LLM (GPT-4o). OM1's dual approach — send to both a local Qwen on edge AND cloud model, use whichever responds first/best — critical for drones where latency = crash risk. Run Qwen3-30B on Jetson for ~300ms local decisions while cloud model provides higher-quality planning in parallel.

**Parallel LLM with action filters for drones:**
- Fast local model → `move`, `avoid_obstacle`, `emergency_land` (latency-critical)
- Cloud model → `plan_route`, `analyze_terrain`, `identify_target` (quality-critical)
- Coach model → mission critique every 30s ("drifting off patrol zone")

**Mode system with automatic transitions:**
- Patrol Mode → hertz: 0.5, inputs: GPS + camera, actions: fly_to waypoints
- Follow Mode → hertz: 2, inputs: VLM tracking, actions: visual servoing
- Emergency Mode → triggered by battery < 20% or signal loss → auto RTH
- Idle/Hover Mode → waiting for commands

Transitions with cooldowns prevent oscillation.

**Natural Language Data Bus (NLDB):**
Caption all sensors into text:
```
"GPS: 37.7749, -122.4194, alt 50m. Battery: 73%. Wind: 12mph NW.
Camera: Open field, 2 vehicles parked, 1 person walking north.
Telemetry: Stable hover, heading 045°."
```
Gives LLM unified, compressed view of all sensor data in one prompt.

**Blockchain governance for drone flight rules:**
Immutable geofencing rules and flight constitutions on-chain:
- Maximum altitude limits
- No-fly zone coordinates
- Operating hours
- Privacy rules (no recording in zones X, Y, Z)
- Cannot be overridden by prompt injection or rogue command

### Combined Stack

```
┌──────────────────────────────────────────────────────────┐
│  OpenClaw (Orchestrator)                                  │
│  Zeus via Telegram → EagleEye agent → mission planning   │
│  Multi-drone coordination, fleet status, approvals       │
├──────────────────────────────────────────────────────────┤
│  OM1 Patterns (Cognitive Layer)                          │
│  Dual LLM (edge + cloud race), Parallel LLM specialists │
│  Mode system (patrol/follow/emergency/idle)              │
│  NLDB (all sensors → natural language → single prompt)   │
│  Governance (on-chain flight rules)                      │
├──────────────────────────────────────────────────────────┤
│  DimOS (Robotics Runtime)                                │
│  MAVLink, LCM pub/sub, MCP skill server                 │
│  DroneConnectionModule, DroneCameraModule                │
│  Visual servoing, GPS nav, object tracking               │
├──────────────────────────────────────────────────────────┤
│  Hardware                                                 │
│  DJI Mavic + RosettaDrone + Android + Jetson (optional)  │
└──────────────────────────────────────────────────────────┘
```

What this gives that neither platform alone provides:
- **DimOS alone**: Single LLM, no mode system, no governance, no orchestrator integration
- **OM1 alone**: Zero drone support, can't fly anything
- **Combined**: Drone hardware control (DimOS) + intelligent multi-LLM cognition (OM1 patterns) + fleet command (OpenClaw) + immutable safety rules (governance)

---

## 7. CRITICAL: Combat / Life-Safety Assessment

**OpenClaw + DimOS is NOT suitable for life-or-death autonomous deployment.**

See [[drone-combat-readiness-assessment]] for the full analysis of why this stack fails under fire (6 failure points: LLM latency, internet dependency, alpha software, single points of failure, no hardware failsafes, no voice/gesture control).

**The solution: a multi-mode system that degrades gracefully from full cloud intelligence to fully autonomous edge operation.**

See [[wartime-weapon-protocol]] for the complete tactical operating procedures — 5 threat levels (SOVEREIGN → VIGILANT → GUARDED → TACTICAL → WARFIGHTER), hardware loadouts, communication failover cascade, PX4 pre-loaded mission library, autonomous formation behaviors, voice command dictionary, EagleEye HUD integration, Meshtastic mesh command protocol, and training drill schedules.

### How the Wartime Protocol Extends This Architecture

This document (Sections 1-6) covers the **peacetime stack**: OpenClaw → DimOS → Drone. The Wartime Protocol adds layers below and around it:

| Layer | Peacetime (this doc) | Wartime (protocol doc) |
| --- | --- | --- |
| **Command** | Telegram → OpenClaw → EagleEye agent | Voice → Jetson → MAVLink / Physical buttons → PX4 |
| **Intelligence** | Cloud LLM (GPT-4o, Claude) via DimOS | Edge AI (Whisper, YOLO, Qwen3-4B) on Jetson |
| **Comms** | Tailscale VPN over internet | Meshtastic LoRa mesh (AES-256, no internet) |
| **Flight control** | DimOS MCP → MAVLink | PX4 pre-loaded autonomous missions (M01-M08) |
| **Display** | Telegram text / dashboard | EagleEye HUD + ATAK tactical map |
| **Failsafe** | DimOS RTH on signal loss | Dead-man switch → auto WARFIGHTER formation |
| **Input** | Keyboard / phone | Panic button / throat mic / RC controller |

The system operates as one unified stack with automatic mode switching:

```text
INTERNET UP:
  OpenClaw (cloud) → DimOS (edge) → Drone (hardware)
  Full intelligence, fleet coordination, Telegram interface

INTERNET DOWN:
  Meshtastic (mesh) → PX4 Companion (edge) → Drone (hardware)
  Mesh commands, edge AI, pre-loaded missions

UNDER FIRE:
  Panic button → PX4 (autonomous) → Drone (hardware)
  Zero human input needed, pre-programmed behaviors execute
```

---

## Related Documents

- [[wartime-weapon-protocol]] — Full tactical operating procedures (5 threat levels, hardware loadouts, training drills)
- [[drone-combat-readiness-assessment]] — Why OpenClaw alone fails for combat (6 failure points)
- [[openmind-om1-deep-analysis]] — OM1 feature analysis (patterns adopted in hybrid architecture)
- [[om1-vs-dimos-vs-openclaw-comparison]] — Platform comparison
- [[OMOS-EagleEye]] — Wearable sensor fusion HUD (tactical display layer)
- [[meshtastic-mesh-integration]] — Mesh network integration (comms backbone)
- [[multi-drone-fleet-architecture]] — Fleet architecture design
- [[OMOS-Drone-Swarm]] — Swarm concepts
