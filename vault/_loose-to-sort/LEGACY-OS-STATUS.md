---
title: Legacy OS — Master Status & Architecture Document
status: active
updated: 2026-05-26
tags: [legacy-ai, onemind-os, architecture, status, nats, hermes, livekit, frontend]
---

# Legacy OS — Master Status Document

> **Core Equation:** `HUMAN (Zeus) + AI (Legacy) = ONE MIND`
> **Architecture:** Hermes (brain) + LiveKit (transport) + NATS (nervous system) + ~700 lines of glue

---

## Table of Contents

1. [System Overview](#1-system-overview)
2. [What Is Running Right Now](#2-what-is-running-right-now)
3. [NATS Nervous System](#3-nats-nervous-system)
4. [Home Assistant Integration](#4-home-assistant-integration)
5. [The Fabric Architecture](#5-the-fabric-architecture)
6. [Frontend — Legacy UI](#6-frontend--legacy-ui)
7. [Key Architecture Decisions](#7-key-architecture-decisions)
8. [Benchmarks](#8-benchmarks)
9. [File & Code Inventory](#9-file--code-inventory)
10. [Vault Notes Index](#10-vault-notes-index)
11. [Next Steps](#11-next-steps)
12. [Blocked Items](#12-blocked-items)

---

## 1. System Overview

Legacy AI is an omnipresent AI system — one identity, one memory, one conversation across all devices and locations. It is not a chatbot. It is an always-on operating intelligence for OneMind OS.

### The Stack

```
┌─────────────────────────────────────────────────────────────────────┐
│                        LEGACY AI — THE FABRIC                        │
├─────────────────────────────────────────────────────────────────────┤
│                                                                       │
│  ┌──────────────┐   ┌──────────────────┐   ┌──────────────────────┐ │
│  │   HERMES     │   │    LIVEKIT       │   │       NATS           │ │
│  │  (The Brain) │   │  (Transport)     │   │  (Nervous System)    │ │
│  │              │   │                  │   │                      │ │
│  │ Memory       │   │ WebRTC voice     │   │ HA commands          │ │
│  │ Skills       │   │ Video streams    │   │ Drone telemetry      │ │
│  │ Sessions     │   │ Data channels    │   │ Watch events         │ │
│  │ Tools/MCP    │   │ Console mode     │   │ IoT devices          │ │
│  │ Soul/identity│   │ (local solo)     │   │ Services fabric      │ │
│  └──────┬───────┘   └────────┬─────────┘   └──────────┬───────────┘ │
│         │                    │                          │             │
│         └────────────────────┴──────────────────────────┘            │
│                           ~700 lines of glue                          │
│              hermes_livekit_transport.py (300)                        │
│              hermes_watch_daemon.py (200)                             │
│              hermes_perception.py (200)                               │
│                                                                       │
└─────────────────────────────────────────────────────────────────────┘
```

### Devices In The Fleet

| Device | Role | Status |
|---|---|---|
| Zeuss-MBP (M2 Max 32GB) | Primary compute, always-on | ✅ Online |
| iPhone | Voice-first mobile interface | 🔜 NATS swift app exists |
| Drones | Autonomous patrol + vision | 🔜 Planned |
| Eden (garden robot) | Field agent | 🔜 Planned |
| Office display | Ambient dashboard | 🔜 Planned |
| Glasses | Audio HUD | 🔜 Future |

---

## 2. What Is Running Right Now

### Active Processes

| Service | PID | Port | Status |
|---|---|---|---|
| NATS Server | 98142 | 4222 | ✅ Running |
| HA-Bridge (ha-bridge.js) | 90791 | — | ✅ Running |
| Legacy UI (Next.js dev) | 51788 | 3000 | ✅ Running |
| Hermes Agent | 7043 | 9119 (dash) / 8787 (WebUI) | ✅ Running |
| Composio SAAS worker | — | — | ✅ Running (launchd) |

### NATS KV Buckets

| Bucket | Size | Values | Purpose |
|---|---|---|---|
| `connections` | 0B | 0 | Agent connections |
| `schedules` | 0B | 0 | Scheduled tasks |
| `triggers` | 0B | 0 | Event triggers |
| `config` | 120B | 2 | System config |
| `agent-state` | 126B | 2 | Agent state tracking |
| `zeus-context` | 356B | 5 | Zeus context |
| `rules` | 4.2KB | 13 | Automation rules |

### NATS JetStream Streams

| Stream | Messages | Purpose |
|---|---|---|
| AGENT_EVENTS | 0 | Agent lifecycle events |
| AGENT_PROMPTS | 0 | Prompts to agents |
| COMMANDS | 0 | System commands |
| COMMS | 0 | Communications |
| HEALTH_EVENTS | 0 | Health monitoring |
| SECURITY_EVENTS | 0 | Security alerts |
| SYSTEM | 0 | System events |
| TASKS | 0 | Task queue |
| ZEUS_EVENTS | 2 | Zeus activity events |
| SUBSCRIBER_EVENTS | 2 | Subscriber notifications |
| SAAS_EVENTS | 8 | SaaS action results |

### Hermes Toolsets

| Toolset | Status |
|---|---|
| web, browser, terminal, file | ✅ enabled |
| code_execution, vision, tts | ✅ enabled |
| skills, todo, memory, session_search | ✅ enabled |
| clarify, delegation, cronjob | ✅ enabled |
| nats | ✅ enabled |
| homeassistant (MCP) | ❌ disabled (removed) |
| video, video_gen, x_search, moa | ❌ disabled |

---

## 3. NATS Nervous System

### Config

```
Config:    ~/.nats/config/onemind.conf
Port:      4222
Domain:    onemind (JetStream)
Server:    nats://127.0.0.1:4222
```

### Subject Hierarchy

```
cmd.ha.service.call          → Call HA service (lights, switches, etc.)
cmd.ha.state.get             → Get HA entity state
cmd.ha.states.list           → List all HA entities
events.ha.*                  → HA state change events (published by ha-bridge)
events.ha.state              → All state changes
events.ha.{domain}           → Domain-specific (light, switch, sensor, etc.)
events.legacy.watch.*        → Watch system events
agent.status                 → Agent state (LISTENING/SPEAKING/THINKING)
perception.*                 → Vision/MediaPipe context feeds
cmd.drone.*.telemetry        → Drone position + battery
saas.execute                 → SaaS action requests (Composio)
saas.event.>                 → SaaS action results
```

### Services Under ~/.nats/services/

| Service | File | Purpose |
|---|---|---|
| ha-bridge | `ha-bridge/bridge.js` | HA ↔ NATS bridge (Node.js, PID 90791) |
| vikunja-bridge | `vikunja-bridge/` | Vikunja task manager bridge |
| vikunja-mcp | `vikunja-mcp/` | Vikunja MCP server |

### NATS WebSocket

NATS runs on TCP port 4222. For browser/frontend access, a WebSocket proxy is needed at port 9222. This is the `nats.ws` connection used by Legacy UI's dashboard panels.

```
nats-server --config ~/.nats/config/onemind.conf
# Add websocket block to config:
websocket {
  port: 9222
  no_tls: true
}
```

---

## 4. Home Assistant Integration

### Architecture

```
Before (slow):  Legacy → Nabu Casa MCP → Cloud → HA   (avg 348ms)
After (fast):   Legacy → NATS → ha-bridge → HA LAN     (avg 0.8-5ms)
```

### Connection Details

```
HA Host:       http://100.64.25.12:8123  (Tailscale)
HA LAN:        http://192.168.1.180:8123
Entities:      185 (cached in ha-bridge)
Bridge:        ~/.nats/services/ha-bridge/bridge.js
Bridge PID:    90791
Token:         ~/.hermes/.env → HASS_TOKEN
```

### Benchmark Results

| Method | Avg RTT | Notes |
|---|---|---|
| Nabu Casa MCP | 348ms | Cloud round-trip, removed |
| NATS ha-bridge state.get | 0.8ms | LAN direct |
| NATS ha-bridge states.list | 2.9ms | 185 entities |
| NATS ha-bridge service.call | 4.3ms | e.g. turn on light |
| Pure NATS (no bridge) | 0.42ms | Baseline |

### Known Entity Status

| Entity | Status |
|---|---|
| switch.plug_mini_ef5a | ✅ Available (tested on/off) |
| switch.plug_mini_6826 | ❌ Unavailable (device offline) |
| alarm_control_panel | ❌ Unavailable (device offline) |
| binary_sensor.bedroom_motion | ❌ Unavailable |
| binary_sensor.keypad | ❌ Unavailable |
| HA backup sensors (all) | ⚠️ Unknown |

### NATS Commands (Usage)

```python
import asyncio
from nats.aio.client import Client as NATS

async def ha_call():
    nc = await NATS().connect("nats://127.0.0.1:4222")
    
    # Get state
    r = await nc.request("cmd.ha.state.get",
        b'{"entity_id": "light.living_room"}', timeout=5)
    
    # Call service
    r = await nc.request("cmd.ha.service.call",
        b'{"domain": "light", "service": "turn_on",
           "entity_id": "light.living_room"}', timeout=5)
    
    # List all entities
    r = await nc.request("cmd.ha.states.list", b'{}', timeout=10)
```

---

## 5. The Fabric Architecture

### Philosophy

> One identity. One memory. One conversation. Across all devices, all locations, all time.

Legacy is NOT a chatbot. It is an omnipresent operating intelligence. The "Fabric" is the architecture that makes this possible without forking or rebuilding anything from scratch.

### The Three Layers

**Layer 1 — Hermes (The Brain)**
- Persistent memory (cross-session, cross-device)
- Skills library (reusable procedures)
- Session management
- 100+ MCP tools
- Soul/identity (SOUL.md injected every session)
- Provider system (swap models without changing code)

**Layer 2 — LiveKit (The Transport)**
- WebRTC audio/video — full duplex, always-on
- Data channels — structured JSON messages to frontend UI
- Console mode — no server needed for local solo use
- Server mode — when phone/glasses/remote devices join
- Model connects DIRECTLY — no STT→text→Hermes→text→TTS wrapping

**Layer 3 — NATS (The Nervous System)**
- Home Assistant commands (sub-5ms)
- Drone telemetry
- IoT/sensor events
- Watch system
- SaaS actions (Composio)
- All services communicate here

### The ~700 Lines of Glue

Three files connect everything. None fork or replace upstream projects.

```
hermes_livekit_transport.py   (~300 lines)  ← NOT WRITTEN YET
  - LiveKit session source plugin for Hermes
  - Connects to LiveKit room (console or server mode)
  - Re-registers Hermes tools as LiveKit @function_tools
  - Loads Hermes memory/soul into system prompt
  - Bridges voice audio ↔ local model pipeline

hermes_watch_daemon.py        (~200 lines)  ← NOT WRITTEN YET
  - NATS subscriber loop
  - Reads watch registry from KV bucket 'watches'
  - LLM evaluates each incoming event against active watches
  - Injects triggers into active LiveKit session via data channel
  - Survives restarts (state in NATS KV, not memory)

hermes_perception.py          (~200 lines)  ← NOT WRITTEN YET
  - MediaPipe hand/pose/face detection
  - Camera feed processing
  - Publishes context to perception.* NATS subjects
  - Feeds the model context, NOT a transport
  - Additive: model sees more, doesn't change how it speaks
```

### What Was Rejected and Why

| Option | Rejected Because |
|---|---|
| OpenClaw voice | Locked to OpenAI Realtime + Google Gemini only. No local model support. |
| Mark-XXXIX (PyQt6) | Desktop only, Gemini only, no NATS, no persistent sessions, 11K line toy |
| Fork Hermes | Too much maintenance for solo developer. Upstream moves fast. |
| Build from scratch | Same problem. Hermes has 3 years of memory/session/tool infrastructure. |
| Hermes push-to-talk voice | It's Ctrl+B record → silence → STT → turn-based. NOT streaming, NOT WebRTC. |
| Nabu Casa MCP for HA | 348ms cloud round-trip. Replaced with direct LAN NATS bridge. |

### Watch System

Watches are persistent background monitors — LLM-supervised event evaluations.

```
Storage:  NATS KV bucket 'watches'  ← BUCKET NOT CREATED YET
Lifecycle: CREATE → ACTIVE → TRIGGERED → DISMISSED/EXPIRED
Evaluation: Every incoming event runs against all active watches
Judgment:   LLM applies context, not just pattern matching

Example watches:
  - "Alert me if Sarah comes home after midnight"
  - "Watch drone battery — warn below 20%"
  - "Tell me if the bar front door opens"
  - "Monitor office temperature — alert if >80°F"
```

### Data Flow (Full)

```
[Zeus speaks] → [LiveKit mic] → [Whisper.cpp STT]
     → [Hermes/LLM processes] → [tool calls via @function_tool]
     → [NATS cmd.ha.* / saas.execute / etc.]
     → [Legacy responds] → [Piper TTS] → [LiveKit speaker]
     → [Data channel] → [Legacy UI frontend updates]

[NATS event arrives] → [Watch daemon evaluates]
     → [Triggered?] → [Inject into LiveKit session]
     → [Legacy speaks proactively]
```

---

## 6. Frontend — Legacy UI

### Location & Status

```
Path:     ~/workspace/legacy-ui/
Status:   ✅ Running at http://localhost:3000 (PID 51788)
Build:    ✅ npm run build passes (6 routes, 146KB first load)
```

### Tech Stack

| Layer | Tech | Version |
|---|---|---|
| Framework | Next.js App Router | 15 |
| Language | TypeScript | 5 |
| Styling | Tailwind CSS | v4 |
| Voice | livekit-client + @livekit/components-react | 2.19 / 2.9 |
| Animation | framer-motion | 12 |
| Realtime data | nats.ws | 1.30 |
| Future desktop | Tauri | — |

### File Structure

```
legacy-ui/
├── app/
│   ├── layout.tsx          ← Root layout, JetBrains Mono font
│   ├── page.tsx            ← 3-column dashboard
│   ├── globals.css         ← Tailwind v4 + scan/pulse keyframes
│   └── api/
│       ├── metrics/route.ts ← Real CPU/MEM/NET via Node os
│       └── agent/route.ts   ← Hermes agent stub (POST text, GET status)
├── components/
│   ├── HudOrb.tsx          ← Animated canvas centerpiece (arcs, pulses, particles)
│   ├── SystemPanel.tsx     ← CPU/MEM/GPU/NET/TEMP bars (polls every 2s)
│   ├── MetricBar.tsx       ← Single metric bar component
│   ├── ActivityLog.tsx     ← Auto-scroll color-coded log
│   ├── WatchPanel.tsx      ← Active watches (●/◉/○ status)
│   ├── VoiceControls.tsx   ← LiveKit connect/mute/VAD (stubbed)
│   ├── TextInput.tsx       ← Cyan-bordered command input
│   ├── Header.tsx          ← L.E.G.A.C.Y title + live clock
│   └── Footer.tsx          ← Status bar + hotkeys
├── lib/
│   ├── types.ts            ← AgentState, Metric, LogEntry, Watch types
│   └── nats.ts             ← NATS WebSocket singleton + subscribe helper
└── public/manifest.json    ← PWA manifest
```

### Color Palette (Mark-XXXIX Inherited)

```
BG:       #00060a   PRI:      #00d4ff
PANEL:    #010d14   PRI_DIM:  #007a99
BORDER:   #0d3347   ACC:      #ff6b00
TEXT:     #8ffcff   ACC2:     #ffcc00
TEXT_DIM: #3a8a9a   GREEN:    #00ff88
                    RED:      #ff3355
```

### Layout

```
┌──────────────────────────────────────────────────────────────┐
│  L.E.G.A.C.Y                                    HH:MM:SS     │
├────────┬────────────────────────────────┬────────────────────┤
│ SYS    │                                │  ACTIVITY LOG      │
│ MONITOR│         HUD ORB                │  (auto-scroll)     │
│        │    (60fps canvas animation)    │                    │
│ CPU ▓▓ │                                │  WATCHES           │
│ MEM ▓▓ │      ● LISTENING              │  ● bar door idle   │
│ GPU ▓▓ │                                │  ◉ drone triggered │
│ NET ▓▓ │  [connect voice] [mute]        │                    │
│        │  [text input ▸]                │  [send ▸]          │
├────────┴────────────────────────────────┴────────────────────┤
│  [F4] Mute  ·  [F11] Fullscreen  ·  HERMES CONNECTED        │
└──────────────────────────────────────────────────────────────┘
```

### What's Wired vs Stubbed

| Feature | Status |
|---|---|
| HudOrb animation (canvas) | ✅ Working |
| System metrics (CPU/MEM/NET) | ✅ Working (via /api/metrics) |
| Activity log | ✅ Working |
| Watch panel | ✅ Working (static demo data) |
| Header clock | ✅ Working |
| LiveKit voice connection | ⏳ Stubbed — needs URL + token |
| NATS WebSocket data | ⏳ Stubbed — needs WS port 9222 |
| Agent text input → Hermes | ⏳ Stubbed — /api/agent not wired |
| Camera/drone video tiles | ⏳ Not built |
| Map panel (drone positions) | ⏳ Not built |

---

## 7. Key Architecture Decisions

| # | Decision | Rationale |
|---|---|---|
| 1 | Pure NATS for HA | 0.8ms vs 348ms. No cloud dependency. |
| 2 | MCP removed | Replaced by native NATS tool. Less latency, less complexity. |
| 3 | Tailscale IP over Nabu Casa | LAN-local. Stable. No internet required. |
| 4 | LiveKit = transport, not brain | Hermes IS the brain. LiveKit is how audio gets in/out. |
| 5 | Model → LiveKit directly | STT→text→Hermes→TTS wrapping is the fragmentation anti-pattern. |
| 6 | Hermes tools as @function_tools | Same capabilities, re-registered. Same result, different caller. |
| 7 | DO NOT fork Hermes | Solo developer. Upstream maintenance cost too high. |
| 8 | DO NOT rebuild from scratch | Hermes has persistent memory, sessions, 100+ tools. Use them. |
| 9 | OpenClaw rejected for voice | Locked to OpenAI Realtime + Gemini. No local model support. |
| 10 | LiveKit console mode for local | No server needed. Direct mic/speaker. Server only for remote. |
| 11 | Watches = NATS KV, not memory | Persistent across restarts. LLM-supervised. Not just patterns. |
| 12 | MediaPipe = perception plugins | Feed context to model. NOT a transport. Additive. |
| 13 | Voice/vision pipeline is local | NATS carries commands/context. Not audio. |
| 14 | PWA (Next.js 15) for frontend | One codebase: desktop/mobile/display. Tauri wrapper if needed. |
| 15 | Mark-XXXIX = UI reference only | Aesthetic only. No NATS, no sessions, 11K lines for one desktop. |

---

## 8. Benchmarks

### NATS Server (Pure)

```
Messages/sec:    2,600,000
Throughput:      1.3 GiB/sec
Latency P50:     0.08 µs
Latency avg:     1.20 µs
```

### HA Integration

```
Nabu Casa MCP:            avg 348ms   ← REMOVED
NATS ha-bridge state.get: avg 0.8ms   (430x faster)
NATS ha-bridge state.lst: avg 2.9ms   (185 entities)
NATS ha-bridge svc.call:  avg 4.3ms
Pure NATS baseline:       avg 0.42ms
```

### Model (Target — Not Yet Built)

```
Target platform:   MBP M2 Max 32GB
STT:               Whisper.cpp (base.en or small.en)
Model:             Qwen2.5-7B-Instruct (Q4_K_M GGUF via llama.cpp)
TTS:               Piper TTS
Estimated TTFT:    ~500ms
Estimated tokens:  ~15-25 tok/sec
RAM budget:        ~10GB (model) + 22GB (system)
```

---

## 9. File & Code Inventory

### Config Files

| File | Purpose |
|---|---|
| `~/.nats/config/onemind.conf` | NATS server config, port 4222, domain onemind |
| `~/.hermes/config.yaml` | Hermes config — toolsets, providers, MCP servers |
| `~/.hermes/.env` | HASS_URL=http://100.64.25.12:8123, HASS_TOKEN |

### NATS Services

| File | Purpose |
|---|---|
| `~/.nats/services/ha-bridge/bridge.js` | HA↔NATS bridge (PID 90791) |
| `~/.nats/services/vikunja-bridge/` | Vikunja task bridge |
| `~/.nats/services/vikunja-mcp/` | Vikunja MCP server |

### Hermes Agent

| File | Purpose |
|---|---|
| `~/.hermes/hermes-agent/` | Hermes agent runtime |
| `~/.hermes/hermes-agent/venv/` | Python venv |
| `~/.hermes/hermes-agent/tools/nats_native_tool.py` | Native NATS toolset |
| `~/.hermes/memories/SOUL.md` | Legacy AI soul/identity |
| `~/.hermes/memories/MEMORY.md` | Persistent agent memory |
| `~/.hermes/memories/USER.md` | Zeus profile |

### Frontend

| File | Purpose |
|---|---|
| `~/workspace/legacy-ui/` | Next.js 15 PWA (21 files) |
| `~/workspace/legacy-ui/components/HudOrb.tsx` | Animated canvas orb |
| `~/workspace/legacy-ui/lib/nats.ts` | NATS WebSocket (stubbed) |
| `~/workspace/legacy-ui/components/VoiceControls.tsx` | LiveKit voice (stubbed) |

### Glue Files (Not Yet Written)

| File | Location | Lines |
|---|---|---|
| `hermes_livekit_transport.py` | `~/workspace/legacy-fabric/` | ~300 |
| `hermes_watch_daemon.py` | `~/workspace/legacy-fabric/` | ~200 |
| `hermes_perception.py` | `~/workspace/legacy-fabric/` | ~200 |

---

## 10. Vault Notes Index

All notes at vault root: `/Users/onemind/Documents/onemind/`

Realtime core notes in: `00-24 SO (Sovereign Operations)/02 Infrastructure (The Body)/services/realtime-core/`

| Note | Size | Contents |
|---|---|---|
| `HA-NATS-Bridge.md` | 6985 bytes | Bridge setup, benchmarks, entity list, usage |
| `OneMind-Realtime-Core.md` | ~37KB, 801 lines | Full architecture, 40 sections, all decisions |
| `Unified-Agent-Runtime.md` | ~20KB | OpenClaw vs Hermes vs LiveKit research, Mark-XXXIX |
| `Legacy-UI-Frontend.md` | 8210 bytes | Frontend architecture, component docs, next steps |
| `LEGACY-OS-STATUS.md` | This file | Master status doc |

---

## 11. Next Steps

### Phase 0 — Immediate (This Session)

- [ ] Create NATS KV bucket `watches`
      `nats kv add watches`
- [ ] Add WebSocket block to NATS config (port 9222)
      Wire `nats.ts` in legacy-ui to live data
- [ ] Wire `/api/agent` in legacy-ui to Hermes HTTP API

### Phase 1 — POC (1-2 days)

- [ ] Write `hermes_livekit_transport.py` (~300 lines)
  - Import LiveKit agents SDK
  - Console mode: `python agent.py console`
  - Load Hermes SOUL.md into system prompt
  - Register 5 key Hermes tools as @function_tools
  - Test: speak → STT → LLM → TTS → hear response

- [ ] Wire VoiceControls in legacy-ui
  - LiveKit console mode URL/token
  - Mute/unmute, VAD state display
  - Agent status dot (LISTENING/SPEAKING/THINKING)

### Phase 2 — Watch System (2-3 days)

- [ ] Write `hermes_watch_daemon.py` (~200 lines)
  - Subscribe to `events.ha.*` and `events.legacy.*`
  - Read active watches from KV bucket `watches`
  - LLM eval each event against watches
  - Inject triggers into active LiveKit session via data channel

- [ ] Wire WatchPanel in legacy-ui to NATS
  - Subscribe to `events.legacy.watch.*`
  - Real-time watch status updates

### Phase 3 — Local Model (1 week)

- [ ] Install Whisper.cpp (STT)
  - `brew install whisper-cpp` or build from source
  - Test: `whisper-cpp --model base.en --file test.wav`

- [ ] Download Qwen2.5-7B-Instruct GGUF
  - `huggingface-cli download Qwen/Qwen2.5-7B-Instruct-GGUF`
  - Q4_K_M quantization (~5GB)

- [ ] Install Piper TTS
  - English neural voice (en_US-lessac-medium)
  - Test: `echo 'System online.' | piper --model en_US-lessac-medium --output_file out.wav`

- [ ] Connect pipeline: Whisper → Qwen → Piper → LiveKit

### Phase 4 — Perception (Optional, add-on)

- [ ] Write `hermes_perception.py` (~200 lines)
  - MediaPipe hand/pose detection
  - Publish to `perception.hands`, `perception.pose`, `perception.face`
  - Test: gesture triggers a watch

### Phase 5 — Remote Devices

- [ ] LiveKit server (Go binary)
  - `livekit-server --dev` for testing
  - Caddy + UDP 50000-60000 for production
- [ ] iPhone NATS swift app connects
- [ ] Legacy speaks to Zeus on phone

---

## 12. Blocked Items

| Item | Blocker | Resolution |
|---|---|---|
| switch.plug_mini_6826 | Device offline | Power cycle the device |
| alarm_control_panel | Device offline | Check hardware |
| binary_sensor.bedroom_motion | Device offline | Check hardware |
| HA backup sensors | Returning `unknown` | HA backup config issue |
| KV bucket `watches` | Not created | Run: `nats kv add watches` |
| NATS WebSocket (9222) | Not configured | Add WS block to onemind.conf |
| LiveKit voice in UI | No transport.py | Write hermes_livekit_transport.py |
| Watch daemon | No daemon file | Write hermes_watch_daemon.py |
| Local model pipeline | Not installed | Phase 3 work |
| Remote device support | No LiveKit server | Phase 5 work |

---

*Generated by Legacy AI — OneMind OS*
*Last updated: 2026-05-26*

