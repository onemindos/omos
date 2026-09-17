---
domain: 04
name: OM1 vs DimOS vs OpenClaw — Platform Comparison
status: reference
owner: zeus
tags: [robotics, om1, dimos, openclaw, drone, comparison, research]
created: 2026-05-12
updated: 2026-05-12
---

# OM1 vs DimOS vs OpenClaw — Platform Comparison

Evaluated May 2026 for Zeus's drone-at-scale use case with real-time commands.

---

## What Each Actually Is

| | **OpenMind OM1** | **DimOS (Dimensional)** | **OpenClaw** |
| --- | --- | --- | --- |
| **Category** | Modular AI runtime for robots | Agentive OS for physical space | AI agent gateway (chat/task orchestration) |
| **Core strength** | Ground robots (humanoids, quadrupeds) | Multi-form robots including **drones** | Conversational agents, workflows, channels |
| **Drone support** | **None documented** | **Yes — MAVLink + DJI Mavic** (alpha) | **None** (not a robotics framework) |
| **Stars / contributors** | 2.8k / 122 | 3.2k / 29 | Private |
| **License** | MIT | Open source | Commercial |
| **Requires API key** | Yes (OMCU billing, 50 free/mo) | No (uses your own LLM keys) | N/A |
| **Source** | [github.com/OpenMind/OM1](https://github.com/OpenMind/OM1) | [github.com/dimensionalOS/dimos](https://github.com/dimensionalOS/dimos) | Private |

---

## Drone-Specific Capabilities

| Feature | OM1 | DimOS | OpenClaw |
| --- | --- | --- | --- |
| MAVLink protocol | No | Yes | No |
| DJI drone support | No | Yes (via RosettaDrone bridge) | No |
| Takeoff / land / arm / disarm | No | Yes (native skills) | No |
| GPS fly-to coordinates | No | Yes (`fly_to(lat, lon, alt)`) | No |
| Visual servoing / object tracking | No | Yes (PID-based, Qwen VLM detection) | No |
| LLM agent control of drone | No | Yes (MCP server + client, GPT-4o default) | No |
| Google Maps / OSM integration | No | Yes (GPS nav + OSM queries) | No |
| Real-time velocity commands | No | Yes (LCM pub/sub, `move(x, y, z, duration)`) | No |
| Simulation / replay without hardware | Yes (Gazebo, Isaac Sim — ground only) | Yes (recorded MAVLink + video replay) | No |
| CLI real-time commands | No drone CLI | `dimos agent-send "fly to the park"` / `dimos mcp call fly_to` | No |

---

## Stack Layer Model — Complementary, Not Competing

```
┌─────────────────────────────────────────────┐
│  OpenClaw — ORCHESTRATION LAYER             │
│  Mission planning, human commands via       │
│  Telegram, multi-agent coordination         │
├─────────────────────────────────────────────┤
│  DimOS — ROBOTICS RUNTIME LAYER            │
│  MAVLink, sensor streams, motor control,    │
│  MCP skills (takeoff, fly_to, track, land)  │
├─────────────────────────────────────────────┤
│  Hardware — PHYSICAL LAYER                  │
│  DJI Mavic + RosettaDrone + Android device  │
└─────────────────────────────────────────────┘
```

- **OpenClaw** = command center (Zeus says "patrol the north field" in Telegram)
- **DimOS** = drone's brain (translates commands into MAVLink flight instructions)
- **OM1** = irrelevant for drones (built for ground robots)

---

## DimOS Drone Module Details

### Architecture

```
DJI Drone → RosettaDrone → MAVLink UDP → connection_module → LCM Topics
                         → Video UDP   → dji_video_stream → tracking_module
```

### Available Drone Skills (MCP-exposed)

**Movement and Control:**
- `move(x, y, z, duration)` — Move with velocity (m/s)
- `takeoff(altitude)` — Takeoff to altitude
- `land()` — Land at current position
- `arm()` / `disarm()` — Arm/disarm motors
- `set_mode(mode)` — Set flight mode (GUIDED, LOITER, etc.)
- `fly_to(lat, lon, alt)` — Fly to GPS coordinates

**Perception:**
- `observe()` — Get current camera frame
- `follow_object(description, duration)` — Follow object with visual servoing
- `is_flying_to_target()` — Check if navigating to GPS target

### Blueprints

| Blueprint | What It Includes |
| --- | --- |
| `drone-basic` | Connection + camera + visualization (foundation) |
| `drone-agentic` | Basic + tracking + Google Maps + OSM + MCP server + LLM agent + web input |

### Network Ports

| Port | Protocol | Purpose |
| --- | --- | --- |
| 14550 | UDP | MAVLink commands/telemetry |
| 5600 | UDP | Video stream |
| 7779 | WebSocket | DimOS web visualization |
| 8765 | WebSocket | Foxglove bridge |
| 7667 | UDP | LCM messaging |

### Hardware Connection

```
DJI Drone ← Wireless → DJI Controller ← USB → Android Device ← WiFi → DimOS Computer
```

Requires RosettaDrone APK on Android as MAVLink bridge.

### Multi-Drone at Scale
- DimOS supports `--robot-ips` for multiple targets
- Each drone runs its own DimOS instance
- OpenClaw coordinates centrally via MCP tool calls

---

## Caveats

- DimOS drone support is **alpha** — expect rough edges
- DJI integration requires Android phone running RosettaDrone
- Real-time latency depends on WiFi between Android device and DimOS host
- For fleet/swarm: one DimOS instance per drone + coordination layer (OpenClaw + custom MCP aggregator)

---

## Recommendation

**Use DimOS as the robotics layer, keep OpenClaw as the orchestration layer.**

DimOS is the only option with actual drone support. It explicitly names OpenClaw as a compatible agent in their docs. The MCP bridge is the glue — DimOS exposes drone skills via MCP server, OpenClaw agents call those tools.

OM1's architecture patterns (dual LLM, parallel LLM, mode transitions, blockchain governance) are worth studying for enhancing the DimOS + OpenClaw stack, but OM1 itself has zero drone capability.

---

## Next Steps

1. Get DJI Mavic + Android device + RosettaDrone APK
2. Install DimOS on VPS or local machine with GPU (`pip install 'dimos[drone]'`)
3. Test with replay first: `dimos --replay run drone-agentic`
4. Add MCP server config to OpenClaw pointing at DimOS MCP endpoint
5. Wire EagleEye agent to call drone skills via MCP tools

See also: [[openmind-om1-deep-analysis]]
