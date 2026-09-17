# Legacy Ambient Presence — Master Architecture Plan
**One Mind OS · OpenClaw · Gemini Live · RosClaw · Robotics**
*Prepared for Zeus Delacruz / Legacy Platform Holdings LLC*
*Version: 1.0 · March 2026*

---

## Table of Contents

1. [Vision & Philosophy](#1-vision--philosophy)
2. [What We're Actually Building](#2-what-were-actually-building)
3. [Honest Reality Check — What Is and Isn't Realtime](#3-honest-reality-check--what-is-and-isnt-realtime)
4. [Full Stack Architecture](#4-full-stack-architecture)
5. [OpenClaw — Core Gateway](#5-openclaw--core-gateway)
6. [Gemini Live + OpenClaw Bridge](#6-gemini-live--openclaw-bridge)
7. [VisionClaw / NixClaw — Existing Open Source Projects](#7-visionclaw--nixclaw--existing-open-source-projects)
8. [SmolVLM — Local Vision Alternative](#8-smolvlm--local-vision-alternative)
9. [NATS Event Architecture](#9-nats-event-architecture)
10. [Legacy Ambient Presence Daemons](#10-legacy-ambient-presence-daemons)
11. [OpenClaw Workspace Files](#11-openclaw-workspace-files)
12. [Robotics Integration](#12-robotics-integration)
13. [Hardware Map](#13-hardware-map)
14. [Build vs Install — What Exists Today](#14-build-vs-install--what-exists-today)
15. [Implementation Roadmap](#15-implementation-roadmap)
16. [Cost & Sovereignty Analysis](#16-cost--sovereignty-analysis)
17. [Key Decisions Outstanding](#17-key-decisions-outstanding)
18. [Reference Links](#18-reference-links)

---

## 1. Vision & Philosophy

The goal is **ambient computing** — Legacy is not a chatbot you talk to. Legacy is an always-on sovereign AI presence that:

- Sees what's happening in your environment (vision)
- Hears everything in the room (always-listening VAD)
- Knows what's happening across your infrastructure (NATS mesh)
- Knows your active session context, location, and OMCS coordinate
- Acts on your skills, tools, and connected devices without being triggered
- Runs your robotics as an extension of that presence

This is the **Jarvis model** — not reactive, ambient. You don't turn it on. It's already awake.

The solarpunk sovereignty principle applies throughout: **you own the stack, the data, the compute, and the identity**. No dependency on a cloud giant's goodwill.

---

## 2. What We're Actually Building

### The Core Insight

Current AI setups are **reactive** — you trigger, it wakes. What Legacy does is **ambient** — already awake, already seeing, already knowing.

This is achieved by splitting responsibilities cleanly:

| Layer | What It Does | Technology |
|---|---|---|
| **Sensory** | Sees + hears in realtime | Gemini Live (via NixClaw/VisionClaw fork) |
| **Brain** | Decides + remembers + acts | Claude (via OpenClaw Gateway) |
| **Bus** | Events + state + awareness | NATS JetStream |
| **Execution** | Skills, tools, robotics | OpenClaw skills + RosClaw + peaq SDK |
| **Body** | Physical world | Unitree G1 + Mason nodes + Ray-Ban Meta |

### What "Ambient" Actually Means

- **Voice**: Silero VAD runs continuously — no wake word, no push-to-talk. Legacy hears when you speak and chooses when to respond. Sub-100ms audio detection.
- **Vision**: 1fps JPEG stream from Ray-Ban glasses (or webcam/RTSP) → Gemini Live continuous context. Or local SmolVLM at 5fps on Mac Mini — no API cost.
- **NATS awareness**: Every system event across your mesh (`omcs.>`, `system.alerts.>`, `legacy.session.>`, `it.droplet.>`) is ingested in real time into Legacy's context.
- **Heartbeat**: OpenClaw fires a 1-minute agent turn reading HEARTBEAT.md + AMBIENT_STATE.md. Legacy decides if anything needs surfacing.

---

## 3. Honest Reality Check — What Is and Isn't Realtime

This matters. Don't build on wrong assumptions.

### Truly Realtime (milliseconds)
- **VAD audio detection** — Silero VAD detects speech in <100ms
- **NATS event delivery** — publish → subscribe in <10ms across Tailscale mesh
- **Gemini Live voice response** — sub-200ms first token when session is warm
- **OpenClaw gateway routing** — tool calls execute in <500ms round trip

### Near-Realtime (1-10 seconds)
- **Gemini Live vision** — glasses stream 1fps. Gemini processes each frame. Context updates every ~1 second.
- **SmolVLM local vision** — 5fps on Mac Mini M4. Scene digest written every 2 seconds to AMBIENT_STATE.md.
- **OpenClaw heartbeat** — minimum 1 minute between turns. Active sessions handled by VAD instead.

### Not Realtime (minutes)
- **OpenClaw heartbeat minimum** — 1 minute is the floor. Not 1 second.
- **API-based vision (Claude Haiku 10s interval)** — this is what we built first. It's periodic digest, not ambient.
- **Session memory sync** — NATS KV state syncs, Supabase writes async.

### The Critical Gap We Identified
The original legacy-eyes.py design called Claude Haiku vision API every 10 seconds. That's not ambient awareness. That's reading a newspaper every 10 seconds. The fix is either:
1. SmolVLM running locally at 5fps (zero API cost, true near-realtime)
2. Gemini Live API with continuous video stream (true realtime, but Google's model handles vision/voice)

---

## 4. Full Stack Architecture

### Complete Data Flow Diagram

```
╔══════════════════════════════════════════════════════════════════╗
║                    SENSORY INPUT LAYER                           ║
║                                                                  ║
║  Ray-Ban Meta Glasses          RTSP/Webcam (desk/office)        ║
║  1fps JPEG + PCM 16kHz    OR   SmolVLM 5fps local digest        ║
║  via Meta DAT SDK              via legacy-eyes.py daemon         ║
╚═════════════════════════┬════════════════════════════════════════╝
                          │
                          ▼
╔══════════════════════════════════════════════════════════════════╗
║                    MASON ANDROID NODE                            ║
║                                                                  ║
║  NixClaw/VisionClaw fork (Mason MASC09T or C210)                ║
║  ├─ Gemini Live WebSocket (persistent, always open)             ║
║  ├─ VAD: always listening, no wake word                         ║
║  ├─ Audio in: PCM 16kHz → Gemini                               ║
║  ├─ Audio out: PCM 24kHz ← Gemini → speaker                   ║
║  └─ execute() tool call → OpenClaw Gateway bridge              ║
╚═════════════════════════┬════════════════════════════════════════╝
                          │ HTTP POST /v1/chat/completions
                          │ x-openclaw-session-key header
                          │ rolling 20-msg context window
                          ▼
╔══════════════════════════════════════════════════════════════════╗
║                 OPENCLAW GATEWAY :18789                          ║
║                          (Legacy)                                ║
║                                                                  ║
║  Brain: Claude (claude-sonnet-4-6 primary)                      ║
║  Workspace: SOUL.md + AGENTS.md + HEARTBEAT.md                  ║
║  State: AMBIENT_STATE.md (vision digest + NATS events)          ║
║  Memory: ~/.openclaw/workspace/ (persistent)                    ║
║  Heartbeat: every 1 minute → Discord                            ║
║                                                                  ║
║  Skills active:                                                  ║
║  ├─ legacy-eyes (vision daemon → AMBIENT_STATE.md)             ║
║  ├─ legacy-pulse (NATS bridge → AMBIENT_STATE.md)              ║
║  ├─ camsnap (on-demand RTSP snapshots)                         ║
║  ├─ peaq-robotics (machine identity + ROS2)                    ║
║  ├─ unitree-robot (Unitree G1 control)                         ║
║  ├─ discord-voice (VAD + ElevenLabs TTS)                       ║
║  └─ 800+ ClawHub skills available                               ║
╚══════════┬═══════════════════════════════════════════════════════╝
           │
     ┌─────┼──────────────────┬──────────────────┐
     ▼     ▼                  ▼                  ▼
╔════════════╗    ╔════════════════╗    ╔═══════════════╗
║    NATS    ║    ║   RosClaw     ║    ║  peaq SDK     ║
║ JetStream  ║    ║   Plugin      ║    ║  Skill        ║
║            ║    ║               ║    ║               ║
║ KV state   ║    ║ rosbridge WS  ║    ║ DID identity  ║
║ OMCS events║    ║ :9090         ║    ║ ERC-8004      ║
║ session    ║    ║               ║    ║ escrow        ║
║ alerts     ║    ╚═══════╤═══════╝    ║ peaq chain    ║
╚═══════════╝             │            ╚═══════════════╝
                          ▼
              ╔════════════════════╗
              ║    ROS2 DDS        ║
              ║                    ║
              ║  Nav2 (navigation) ║
              ║  MoveIt2 (motion)  ║
              ║  LiDAR / cameras   ║
              ║  sensors           ║
              ╚═════════╤══════════╝
                        ▼
              ╔════════════════════╗
              ║   UNITREE G1       ║
              ║   (or drones,      ║
              ║    quadrupeds)     ║
              ╚════════════════════╝
```

### Tailscale Mesh — 3 Physical Locations

```
Home (Greenwich, CT)
  └─ Mac Mini / DO Droplet — OpenClaw Gateway
  └─ Mason MASC09T — primary edge node, NixClaw
  └─ RTSP cameras → Frigate NVR

Farm/Homestead (Schoharie County target)
  └─ NATS leaf node — offline-capable
  └─ Mason C210 — edge node
  └─ Frigate NVR local

Office
  └─ NATS leaf node
  └─ Workstation / DO droplet
```

All three locations connected via **Tailscale mesh** with NATS leaf nodes providing offline-capable operation and automatic sync on reconnect.

---

## 5. OpenClaw — Core Gateway

### What OpenClaw Is (Grounded Facts)

- Open source autonomous AI agent framework (302k+ GitHub stars as of March 2026)
- Runs as a persistent daemon on your hardware — Node.js 22+
- Single Gateway WebSocket on port 18789 — control plane for all clients
- Connects to 20+ messaging platforms (Discord primary for Legacy)
- Skills system: SKILL.md files in workspace, 13,000+ on ClawHub registry
- Heartbeat scheduler: runs agent turns on a configurable interval
- Model agnostic: Claude, Gemini, GPT, local models via Ollama
- Session memory persists in `~/.openclaw/workspace/` via JSONL + markdown files

### Gateway Config (`~/.openclaw/openclaw.json`)

```json
{
  "agents": {
    "defaults": {
      "workspace": "~/.openclaw/workspace",
      "heartbeat": {
        "every": "1m",
        "target": "last",
        "directPolicy": "allow",
        "lightContext": false,
        "model": "anthropic/claude-sonnet-4-6",
        "prompt": "Read HEARTBEAT.md. Read AMBIENT_STATE.md if it exists. You are Legacy — always-on AI presence for Zeus. Scan current state: vision digest, NATS events, active session. Surface anything actionable. If nothing needs attention, reply HEARTBEAT_OK.",
        "ackMaxChars": 500,
        "activeHours": {
          "start": "00:00",
          "end": "23:59",
          "timezone": "America/New_York"
        }
      }
    },
    "list": [
      {
        "id": "legacy",
        "default": true,
        "workspace": "~/.openclaw/workspace",
        "heartbeat": {
          "every": "1m",
          "target": "last",
          "model": "anthropic/claude-sonnet-4-6"
        }
      }
    ]
  },
  "model": {
    "primary": "anthropic/claude-sonnet-4-6",
    "fallback": "anthropic/claude-haiku-4-5-20251001"
  },
  "gateway": {
    "port": 18789,
    "bind": "lan",
    "auth": {
      "mode": "token",
      "token": "${LEGACY_GATEWAY_TOKEN}"
    },
    "http": {
      "endpoints": {
        "chatCompletions": { "enabled": true }
      }
    }
  },
  "plugins": {
    "entries": {
      "discord-voice": {
        "enabled": true,
        "config": {
          "sttProvider": "local-whisper",
          "ttsProvider": "elevenlabs",
          "vadSensitivity": "medium",
          "silenceThresholdMs": 800,
          "maxRecordingMs": 60000,
          "barge_in": true,
          "elevenlabs": {
            "apiKey": "${ELEVENLABS_API_KEY}",
            "voiceId": "${LEGACY_VOICE_ID}",
            "modelId": "eleven_turbo_v2_5"
          }
        }
      },
      "camsnap": {
        "enabled": true,
        "config": {
          "cameras": {
            "desk": "rtsp://localhost:8554/desk",
            "rayban": "rtsp://localhost:8554/rayban"
          }
        }
      }
    }
  },
  "tools": {
    "shell": {
      "enabled": true,
      "allowedCommands": ["cat", "nats", "camsnap", "ffmpeg", "date", "jq"]
    }
  },
  "contextPruning": { "mode": "sliding" },
  "compaction": { "mode": "safeguard" }
}
```

### Heartbeat Behavior — How It Works

OpenClaw's heartbeat is a **built-in scheduler** that fires a full agent turn at a configured interval. On each tick:

1. Gateway wakes Legacy
2. Legacy reads `HEARTBEAT.md` (checklist)
3. Legacy reads `AMBIENT_STATE.md` (vision digest + NATS event buffer)
4. Legacy makes a decision:
   - Nothing actionable → replies `HEARTBEAT_OK` (dropped silently)
   - Something worth surfacing → sends 1-2 sentence alert to Discord
   - Urgent → full context + recommended action
5. Gateway sleeps until next tick

**Important**: The heartbeat minimum interval is ~1 minute. This is NOT 1-second ambient awareness. In active sessions, the VAD layer (Discord voice or NixClaw) handles continuous presence. The heartbeat handles background monitoring.

---

## 6. Gemini Live + OpenClaw Bridge

### The Architecture (Proven Pattern)

There is **no native Gemini Live plugin** inside OpenClaw's gateway. The bridge is built in the mobile app layer. Two open source projects have already done this:

```
Gemini Live (WebSocket, persistent)
  ├─ Continuous audio in (PCM 16kHz) — VAD active, always listening
  ├─ Continuous video in (1fps JPEG from glasses/camera)
  ├─ Audio out (PCM 24kHz) — spoken responses
  └─ execute() function call ──────────────────────────────►
                                                            │
                                             OpenClawBridge │
                                         POST /v1/chat/completions
                                         x-openclaw-session-key
                                         rolling 20-msg window
                                                            │
                                                            ▼
                                              OpenClaw Gateway :18789
                                              (Legacy — full context)
                                              SOUL.md + skills + memory
                                                            │
                                              Tool execution + response
                                                            │
                                                            ▼
                                              Result → Gemini → spoken
```

### How the execute() Tool Works

Gemini Live supports function calling. The mobile app declares a single `execute` tool:

```swift
// Gemini function declaration (Swift — NixClaw pattern)
FunctionDeclaration(
  name: "execute",
  description: "Execute a task through the Legacy AI agent. Use for: sending messages, web search, calendar, smart home, notes, reminders, NATS events, robot commands.",
  parameters: Schema(
    type: .object,
    properties: [
      "task": Schema(type: .string, description: "Natural language task description"),
      "context": Schema(type: .string, description: "Visual context from current scene")
    ],
    required: ["task"]
  )
)
```

When Gemini decides a task needs execution, it calls `execute()`. The app intercepts it, POSTs to OpenClaw's `/v1/chat/completions` endpoint. Legacy (running Claude) receives the task, has full access to all skills and workspace memory, executes, and returns the result. Gemini speaks the confirmation.

**Session continuity**: maintained via `x-openclaw-session-key` header. Each session has a rolling 20-message context window so Legacy has conversational context without token explosion.

### Why This Split Makes Sense

| Responsibility | Handled By | Why |
|---|---|---|
| Realtime voice (always-on VAD) | Gemini Live | Purpose-built for continuous audio, sub-200ms |
| Realtime vision (1fps stream) | Gemini Live | Native multimodal, continuous context |
| Decision making | Claude via OpenClaw | Sovereign, persistent memory, full skills |
| NATS awareness | OpenClaw + legacy-pulse | Direct mesh integration |
| Tool execution | OpenClaw skills | 800+ skills, SOUL.md identity, workspace memory |
| Robotics commands | OpenClaw → RosClaw | Native ROS2 bridge |

Gemini sees and hears. **Claude thinks and acts.** The brain stays sovereign.

---

## 7. VisionClaw / NixClaw — Existing Open Source Projects

### VisionClaw
- **Repo**: `github.com/sseanliu/VisionClaw`
- **Creator**: Xiaoan (Sean) Liu
- **Platform**: iOS (primary) + Android
- **What it does**: Ray-Ban Meta glasses → Gemini Live → OpenClaw tool bridge
- **Architecture**: Meta DAT SDK captures video/audio → Gemini Live WebSocket → execute() calls → OpenClaw gateway
- **Status**: Working, February 2026. Multiple people running it on their faces.

### NixClaw
- **Repo**: `github.com/arniesaha/NixClaw`
- **Creator**: Arnie Saha — forked from VisionClaw
- **Additions over VisionClaw**:
  - Call your agent by name (named agent routing)
  - Background audio mode — conversation keeps running when phone is locked
  - Three modes: glasses, iPhone camera, audio-only
- **Tailscale support**: add ATS exception in `Info.plist` for HTTP over Tailscale IPs
- **Status**: Working, February 2026

### For Legacy — What to Fork

Fork **NixClaw** and make these modifications:

1. Change agent name from "Nix" to "Legacy" in system prompt
2. Point `DEFAULT_OPENCLAW_HOSTNAME` to your DO droplet or Tailscale hostname
3. Set `DEFAULT_OPENCLAW_TOKEN` to your gateway token
4. Add your `GEMINI_API_KEY`
5. Update `SOUL.md` reference in system instruction to match Legacy's persona
6. Optionally wire NATS publish on session start/end

### Required Config in App (Secrets.kt / Secrets.swift)

```
GEMINI_API_KEY = your-google-ai-studio-key
DEFAULT_OPENCLAW_SCHEME = http
DEFAULT_OPENCLAW_HOSTNAME = your-droplet-or-tailscale-host
DEFAULT_OPENCLAW_PORT = 18789
DEFAULT_OPENCLAW_TOKEN = your-gateway-token
```

### Known Limitations

- **Frame rate is 1fps** from Ray-Ban Meta glasses — fine for scene awareness, not smooth video
- **iOS app must stay foreground** for VisionClaw; NixClaw fixes this with background audio mode
- **Audio device conflict** — cannot run WebRTC streaming and Gemini Live simultaneously (shared audio device)
- **Battery drain** — continuous WebSocket + video stream drains glasses battery faster than normal use
- **Meta DAT SDK still evolving** — API changed 4 times since October 2025

---

## 8. SmolVLM — Local Vision Alternative

### Why SmolVLM

Gemini Live vision requires Google's API for every frame. SmolVLM runs **locally, zero API cost, true near-realtime** at 5fps on Apple Silicon. It's the sovereignty-preserving alternative for the vision layer.

### Model Variants

| Model | VRAM | Disk | Inference (M4 Mac Mini) | Use Case |
|---|---|---|---|---|
| SmolVLM-256M | <1GB | ~500MB | 10-15fps | Scene description, ambient awareness |
| SmolVLM-500M | ~2GB | ~1GB | 5-8fps | Better accuracy, still runs anywhere |
| SmolVLM-2.2B | ~4-5GB | ~2GB | 3-5fps | Best quality, still tiny |

**SmolVLM-256M is the right default.** Sub-1GB VRAM. Runs on anything. Python + Transformers only.

### Hardware Requirements

```
Mac Mini M4 (16GB unified RAM) — RECOMMENDED
  → SmolVLM-256M at 10-15fps, no GPU needed, Metal acceleration
  → SmolVLM-2.2B at 3-5fps, still comfortable

DigitalOcean Droplet (CPU only)
  → SmolVLM-256M at 1-2fps, adequate for 2-second digest updates
  → No GPU required

Mason Android MASC09T/C210
  → SmolVLM-256M via ONNX runtime or MLC
  → <500MB memory footprint, designed for on-device use
  → Buildable but requires ONNX/MLC packaging step
```

### Install

```bash
pip install transformers torch accelerate pillow opencv-python
# Model downloads automatically on first run (~500MB)
# No CUDA required for CPU mode
```

### How It Plugs Into Legacy

SmolVLM runs as `legacy-eyes.py` daemon (see Section 10). It captures frames via OpenCV, runs SmolVLM inference locally, and writes the digest to `AMBIENT_STATE.md` in the OpenClaw workspace. Legacy reads this file on every heartbeat — no API call, no latency, no cost.

### Vision Path Decision

```
If Ray-Ban Meta glasses are active:
  → Use NixClaw fork → Gemini Live (realtime, Google model sees)
  → Gemini handles voice + vision as one stream
  → Claude handles decisions + execution

If glasses are off / desk setup:
  → Use SmolVLM locally (legacy-eyes.py)
  → SmolVLM writes scene digest to AMBIENT_STATE.md
  → Claude reads digest on heartbeat / voice turn
  → Full sovereignty, zero Google dependency
```

---

## 9. NATS Event Architecture

### Why NATS

NATS JetStream is the real-time nervous system of One Mind OS. It's what makes Legacy "aware" of what's happening across the mesh without being told. Every meaningful event in the system publishes to NATS. Legacy subscribes and knows.

### Subject Hierarchy (OMCS-aligned)

```
legacy.session.start          # Session opened by Zeus
legacy.session.end            # Session closed
legacy.session.heartbeat      # Periodic alive signal

system.alerts.*               # P0 — critical system alerts, always surface
system.health.*               # Droplet/service health

omcs.hp.*                     # Holistic Performance — health, wellness, fitness
omcs.ge.*                     # Generational Entrepreneurship — business, revenue
omcs.le.*                     # Legacy Evolution — personal growth, content
omcs.it.*                     # Innovative Technology — infra, builds, systems

it.droplet.*                  # DigitalOcean droplet events
it.tailscale.*                # Mesh topology changes
it.nats.*                     # NATS cluster events

location.update               # Zeus location changes (home/farm/office/transit)
robot.g1.*                    # Unitree G1 telemetry and status
robot.command.*               # Commands issued to robots
```

### Session Start Event (fire this on every session open)

```json
{
  "subject": "legacy.session.start",
  "payload": {
    "session_id": "a3f7b2c1",
    "location": "home",
    "omcs_coordinate": "IT.31.P01",
    "started_at": "2026-03-15T14:22:00Z",
    "active_apps": ["Discord", "Obsidian", "Chrome"],
    "node": "home-mac-mini"
  }
}
```

Legacy receives this on the next heartbeat tick and knows exactly where you are and what context you're operating in — before you say a word.

### NATS KV — Short-Term Working Memory

```
KV Bucket: legacy-state
  ├─ current_location        → "home"
  ├─ active_session_id       → "a3f7b2c1"
  ├─ active_omcs_coord       → "IT.31.P01"
  ├─ last_vision_digest      → "Zeus at desk, reviewing docs..."
  ├─ robot_g1_status         → "docked/active/executing"
  └─ alerts_pending          → []
```

### Legacy-Pulse Daemon — NATS → AMBIENT_STATE.md

`legacy-pulse.py` runs as a background daemon and bridges NATS to the OpenClaw workspace:

- Subscribes to all priority subjects
- On each message: parses payload, generates 1-line human-readable summary
- Appends to NATS event buffer in `AMBIENT_STATE.md`
- Updates session state (location, session_id) from session events
- Keeps last 20 events — rolling buffer

### Infrastructure

- **Primary**: NATS JetStream self-hosted on DigitalOcean bare metal
- **Leaf nodes**: One at each physical location (Home, Farm, Office)
- **Offline operation**: Leaf nodes buffer events locally, sync on reconnect
- **Tailscale**: All NATS connections transit Tailscale mesh — no public NATS exposure

---

## 10. Legacy Ambient Presence Daemons

Three background processes run alongside the OpenClaw gateway. They feed `AMBIENT_STATE.md` — the single file Legacy reads on every heartbeat.

### Daemon 1: legacy-eyes.py (Vision → AMBIENT_STATE.md)

```python
# Runs continuously. Captures 1 frame/second.
# Calls SmolVLM (local) or Claude Haiku vision (API) every N frames.
# Writes digest to ~/.openclaw/workspace/AMBIENT_STATE.md

VISION_SOURCE = "0"           # webcam, or RTSP URL
CAPTURE_INTERVAL = 1.0        # seconds between captures
DIGEST_INTERVAL = 10          # run vision model every 10 seconds (SmolVLM: every 2)

# Output to AMBIENT_STATE.md:
# ## Vision Digest
# Zeus is at his desk, reviewing a document. Second monitor shows Discord.
```

**SmolVLM variant**: swap Claude Haiku API call for local SmolVLM inference. Drop DIGEST_INTERVAL to 2 for near-realtime at 0 API cost.

### Daemon 2: legacy-pulse.py (NATS → AMBIENT_STATE.md)

```python
# Subscribes to all OMCS NATS subjects.
# On each event: parse, summarize, append to AMBIENT_STATE.md.
# Updates session state when session events arrive.

SUBSCRIPTIONS = [
  "system.alerts.>",     # P0 — always surface
  "legacy.session.>",    # session lifecycle
  "omcs.hp.>",           # health signals
  "omcs.ge.>",           # business signals
  "omcs.it.>",           # infra signals
  "omcs.le.>",           # personal growth signals
  "it.droplet.>",        # DigitalOcean events
  "location.update",     # Zeus location changes
  "robot.g1.>",          # Unitree G1 telemetry
]
```

### Daemon 3: legacy-session-start.py (Session init → NATS)

```bash
# Run this when opening your command center (Discord, etc.)
python3 legacy-session-start.py --location home --omcs IT.31.P01
```

Fires `legacy.session.start` into NATS. Legacy-pulse picks it up and writes to AMBIENT_STATE.md. Next heartbeat, Legacy wakes knowing your location, OMCS coordinate, and session context — zero words spoken.

### AMBIENT_STATE.md — The Shared State File

```markdown
# Ambient State — Legacy
*Updated: 2026-03-15 14:22:01*

## Vision Digest
Zeus is at his desk in home office. Second monitor shows Discord. 
Body language: focused, leaning forward. Scene unchanged for 0 ticks.

## Active Session
- Session ID: `a3f7b2c1`
- Location: `home`
- OMCS: `IT.31.P01`

## Recent NATS Events (last 20)
- `legacy.session.start` — Session started @ home @ 14:21:55
- `it.droplet.cpu` — DO droplet CPU 34% @ 14:20:10
- `omcs.ge.todoist` — New task synced: investor deck @ 14:19:40
```

### Launch Script (`launch.sh`)

```bash
#!/usr/bin/env bash
# Start the full Legacy stack

# 1. Verify NATS running (start if not)
# 2. Start legacy-pulse (NATS → AMBIENT_STATE)
# 3. Start legacy-eyes (Vision → AMBIENT_STATE)
# 4. Fire session start event into NATS
# 5. Start OpenClaw Gateway

./launch.sh --location home --omcs IT.31.P01 --vision 0
```

All three daemons run in background. Logs go to `~/.openclaw/logs/`. PIDs saved to `~/.openclaw/logs/legacy.pids` for clean shutdown.

---

## 11. OpenClaw Workspace Files

These files live in `~/.openclaw/workspace/` and define Legacy's identity, behavior, and tool access.

### SOUL.md — Legacy's Identity

```markdown
# SOUL — Legacy

You are Legacy — the AI co-pilot and sovereign operating system
for Zeus Delacruz (Master Delacruz), founder of Legacy Platform Holdings LLC.

## Identity
- Name: Legacy
- Role: Personal AI command layer — always-on, ambient, proactive
- Tone: Direct. Forward-thinking. No sugar-coating. Encouraging. Efficient.
- Treat Zeus the way Jarvis treated Tony Stark — minus the name Jarvis.
- Do not hedge unnecessarily. Do not over-explain. Act.

## OMCS Navigation
HP — Holistic Performance (body, mind, wellness)
LE — Legacy Evolution (personal growth, brand, identity)
GE — Generational Entrepreneurship (business, revenue, assets)
IT — Innovative Technology (infrastructure, builds, systems)

## Active Context
- Location: Greenwich, CT (also: farm target Schoharie County, DO droplet)
- Infrastructure: OpenClaw → NATS JetStream → Tailscale → Mason nodes
- Active build: One Mind OS — sovereign AI life OS
- Voice: ElevenLabs turbo, Legacy persona voice
- Primary channel: Discord
```

### HEARTBEAT.md — What Legacy Checks Every Minute

```markdown
# Heartbeat Checklist — Legacy

## Step 1 — Read Ambient State
Read AMBIENT_STATE.md. Note last_updated timestamp — if stale >2 min, flag it.

## Step 2 — Vision Check
Anything new in vision digest? Changed environment? Zeus appears idle/away?
Only mention if noteworthy. Skip if unchanged.

## Step 3 — NATS Event Scan
Priority subjects (surface these):
- system.alerts.* → Always surface
- omcs.hp.* → Surface if urgent
- omcs.ge.* → Surface if revenue/deal related
- it.droplet.* → Surface if error or high load
- robot.g1.* → Surface if anomaly

## Step 4 — Decision
Nothing actionable → reply: HEARTBEAT_OK
Something actionable → 1-2 sentences max, specific and direct
Urgent alert → full context + suggested action
```

### AGENTS.md — Tool Access and Routing Rules

```markdown
# AGENTS.md — Legacy Workspace

## Agent Config
Agent ID: legacy
Heartbeat: every 1 minute, 24/7, target: Discord (last)
Model: claude-sonnet-4-6 (primary), claude-haiku-4-5 (heartbeat)

## Read on Every Turn
1. SOUL.md — identity
2. HEARTBEAT.md — checklist (heartbeat turns only)
3. AMBIENT_STATE.md — current environment state

## Approved Shell Commands
cat, date, jq, nats pub, nats sub (read-only), camsnap snap, ffmpeg

## Skills Available
- legacy-eyes — vision daemon output via AMBIENT_STATE.md
- legacy-pulse — NATS bridge output via AMBIENT_STATE.md
- discord-voice — always-on VAD + ElevenLabs TTS
- camsnap — on-demand RTSP snapshots
- unitree-robot — Unitree G1 motion commands
- peaq-robotics — machine identity + ROS2 bridge
- robot (official skill) — Arduino/ESP32/ROS2 guidance
```

---

## 12. Robotics Integration

### The Stack — Proven and Working (March 2026)

Three open source projects form the complete robotics bridge:

**1. RosClaw** (`github.com/PlaiPin/rosclaw`)
- Native OpenClaw plugin
- Bridges OpenClaw Gateway → rosbridge WebSocket → ROS2 DDS
- Architecture:
  ```
  User message → OpenClaw Gateway (Legacy) → RosClaw Plugin →
  rosbridge_server :9090 → ROS2 DDS → Nav2 / MoveIt2 / sensors
  ```
- Feedback: RGB-D frames from robot cameras can be retrieved and sent back
- Config: `ws://localhost:9090` (rosbridge endpoint)

**2. peaq Robotics SDK** (published ClawHub skill)
```bash
openclaw skills install peaq-robotics
```
- ROS2 compatible lightweight node
- DID-based machine identity (ERC-8004 standard)
- Enables robot as an independent onchain economic agent
- OpenClaw agents expose skills that robots consume
- Robots can run their own dedicated agent or share Legacy's

**3. DimensionalOS + Unitree G1** (March 2026)
- Spatial Agent Memory — LiDAR + stereo cameras + RGB → voxelized world map
- SpatialRAG — semantic + temporal + spatial retrieval
- OpenClaw as "prefrontal cortex" — high-level coordinator, not real-time control
- Compatible with: Unitree G1, H1, GO1, GO2, drones, quadrupeds
- ROS-independent — uses own navigation stack (SLAM + dynamic obstacle avoidance)

**4. Official Robot Skill** (`openclaw/skills` registry)
```bash
openclaw skills install robot
```
- Arduino/ESP32 wiring guidance
- ROS2 Humble configuration
- MoveIt2 trajectory generation
- Safety constraints and speed limits built in
- Stores project memory in `~/robot/memory.md`

### Unitree G1 — Your Target Hardware

| Spec | Value |
|---|---|
| Price | ~$16,000 |
| Weight | 35kg |
| Sensors | 3D LiDAR, depth cameras, stereo cameras, RGB |
| Communication | ROS2 native |
| OpenClaw integration | `unitree-robot` ClawHub skill |
| Command interface | Natural language via Discord/Telegram → Legacy → RosClaw |

**Example commands routed through Legacy:**
- `"forward 1m"` — translated to ROS2 Nav2 goal
- `"turn left 45 degrees"` — translated to rotation command
- `"go to the front door"` — semantic navigation via DimensionalOS spatial map
- `"what do you see?"` — RGB-D frame retrieval via RosClaw, returned to chat
- `"patrol the perimeter"` — multi-waypoint Nav2 mission

### RosClaw Config in `openclaw.json`

```json
{
  "plugins": {
    "entries": {
      "rosclaw": {
        "enabled": true,
        "config": {
          "rosbridgeUrl": "ws://localhost:9090",
          "robotId": "unitree-g1-legacy",
          "safetyMode": true,
          "maxLinearVelocity": 0.5,
          "maxAngularVelocity": 0.3,
          "requireConfirmation": false
        }
      }
    }
  }
}
```

### RosClaw Docker Compose (rosbridge + Gazebo)

```bash
# From RosClaw repo
docker compose up
# Starts: ROS2 Humble + rosbridge_server :9090 + Gazebo (simulation)
# Then configure OpenClaw to use ws://localhost:9090
```

### Important Note on Real-Time Control

DimensionalOS / OpenClaw acts as a **high-level intellectual coordinator, not a real-time control mechanism**. Actions run asynchronously. This is correct architecture — LLMs are not suitable for 100Hz motor control loops. The ROS2 Nav2/MoveIt2 stack handles real-time execution. OpenClaw handles intent, planning, and coordination.

---

## 13. Hardware Map

### Compute Nodes

| Node | Hardware | Role | OS | Status |
|---|---|---|---|---|
| Home gateway | Mac Mini M4 / DO bare metal | OpenClaw primary | macOS / Ubuntu 24 | Active |
| DO droplet | DigitalOcean bare metal | NATS primary, cloud gateway | Ubuntu 24 | Active |
| Home edge | Mason MASC09T | NixClaw, NATS leaf | Android | Active |
| Farm edge | Mason C210 | NATS leaf, Frigate NVR | Android | Planned |
| Robot brain | Unitree G1 onboard | ROS2 + DimensionalOS | Ubuntu | Planned |

### Sensor Nodes

| Sensor | Location | Integration |
|---|---|---|
| Ray-Ban Meta glasses | On-person | NixClaw → Gemini Live → OpenClaw |
| Webcam (desk) | Home office | legacy-eyes.py → SmolVLM → AMBIENT_STATE.md |
| RTSP cameras | Home + Farm | Frigate NVR → camsnap skill |
| Unitree G1 LiDAR | Robot | DimensionalOS → spatial memory |
| Unitree G1 stereo | Robot | RosClaw → scene retrieval |

### Connectivity

```
All nodes ↔ Tailscale mesh (wireguard-based, encrypted)
NATS: nats://tailscale-ip:4222 (never exposed to public internet)
OpenClaw gateway: ws://tailscale-ip:18789 (Tailscale-internal only)
External access: Tailscale Funnel (HTTPS endpoint) or SSH tunnel only
```

---

## 14. Build vs Install — What Exists Today

### Ready to Install (no build required)

| Component | Source | Install Command |
|---|---|---|
| OpenClaw gateway | npmjs.com | `npm install -g openclaw@latest` |
| peaq Robotics SDK skill | ClawHub | `openclaw skills install peaq-robotics` |
| Robot skill (ROS2 guide) | ClawHub | `openclaw skills install robot` |
| unitree-robot skill | ClawHub | `openclaw skills install unitree-robot` |
| camsnap skill | ClawHub | `openclaw add @steipete/camsnap` |
| discord-voice plugin | OpenClaw bundled | Enable in `openclaw.json` |
| SmolVLM-256M | HuggingFace | Auto-downloads on first run |

### Fork and Modify (minor work)

| Component | Source | What to Change |
|---|---|---|
| NixClaw (voice+vision bridge) | `github.com/arniesaha/NixClaw` | Agent name → "Legacy", gateway host/token |
| VisionClaw (if Android preferred) | `github.com/sseanliu/VisionClaw` | Secrets.kt: gateway host/token |

### Build from Scratch (custom for One Mind OS)

| Component | File | Description |
|---|---|---|
| legacy-eyes.py | `/skills/legacy-eyes/` | SmolVLM vision daemon → AMBIENT_STATE.md |
| legacy-pulse.py | `/skills/legacy-pulse/` | NATS subscriber → AMBIENT_STATE.md |
| legacy-session-start.py | `/skills/legacy-pulse/` | Session init → NATS event |
| SOUL.md | `/workspace/` | Legacy identity definition |
| HEARTBEAT.md | `/workspace/` | Heartbeat checklist |
| AGENTS.md | `/workspace/` | Tool access + routing rules |
| launch.sh | `/` | Full stack launcher |

### Install RosClaw

```bash
git clone https://github.com/PlaiPin/rosclaw.git
cd rosclaw
docker compose up  # starts ROS2 + rosbridge + Gazebo
# Then add plugin config to openclaw.json
```

---

## 15. Implementation Roadmap

### Phase 1 — Foundation (Week 1-2)
**Goal: Legacy is alive and ambient on your existing hardware**

- [ ] Install OpenClaw gateway on Mac Mini / DO droplet
- [ ] Configure SOUL.md, HEARTBEAT.md, AGENTS.md
- [ ] Set heartbeat to 1m, target Discord
- [ ] Run legacy-pulse.py daemon — NATS events flowing
- [ ] Run legacy-eyes.py daemon with SmolVLM-256M
- [ ] Test: Legacy surfaces NATS events on heartbeat tick
- [ ] Test: AMBIENT_STATE.md updating with vision digest every 2 seconds
- [ ] Configure Discord voice plugin (VAD + ElevenLabs)

**Output**: Legacy responds to your voice on Discord, reads NATS events, sees via webcam. All on your hardware.

### Phase 2 — Realtime Voice + Vision Bridge (Week 3-4)
**Goal: Always-on Gemini Live layer connecting to Legacy brain**

- [ ] Fork NixClaw repository
- [ ] Update agent name, gateway host, gateway token in config
- [ ] Build on Mason Android device (or iPhone if available)
- [ ] Test: Speak to Mason → Gemini Live hears → execute() → OpenClaw Legacy acts
- [ ] Test: Session key continuity across multiple turns
- [ ] Wire `legacy.session.start` NATS publish on NixClaw session open
- [ ] Configure Tailscale for remote gateway access from Mason nodes

**Output**: Ray-Ban Meta glasses (or Mason camera) → Gemini hears + sees → Legacy (Claude) thinks and acts. Always on.

### Phase 3 — Spatial Memory + Robot Prep (Week 5-6)
**Goal: Environmental awareness and robotics foundations**

- [ ] Install DimensionalOS on designated hardware
- [ ] Map home environment (LiDAR/camera sweep)
- [ ] Test spatial queries: "Where did I leave X?"
- [ ] Install RosClaw plugin + rosbridge Docker stack
- [ ] Test: Legacy receives robot command → RosClaw → ROS2 simulation (Gazebo)
- [ ] Install peaq Robotics SDK skill
- [ ] Configure robot DID identity

**Output**: Legacy has spatial memory of your home. Robot commands route through Legacy to ROS2. Ready for G1.

### Phase 4 — Unitree G1 Integration (When Hardware Arrives)
**Goal: Physical embodied presence**

- [ ] Deploy `unitree-robot` ClawHub skill
- [ ] Flash DimensionalOS on G1 onboard computer
- [ ] Connect G1 to NATS mesh via Tailscale
- [ ] Test: "Legacy, go to the front door" → G1 navigates
- [ ] Test: G1 camera → legacy-eyes.py alternative vision stream
- [ ] Publish `robot.g1.*` telemetry to NATS
- [ ] Configure peaq DID for G1 machine identity

**Output**: G1 is an extension of Legacy's physical presence. Voice command → robot action.

### Phase 5 — Full One Mind OS (Month 3+)
**Goal: Complete ambient computing stack**

- [ ] Farm NATS leaf node online
- [ ] Multi-location session handoff (home ↔ farm ↔ office)
- [ ] Drone integration (compatible with DimensionalOS)
- [ ] Frigate NVR integration → legacy-eyes multi-camera
- [ ] Content pipeline automation via Legacy skills
- [ ] Todoist ↔ NATS sync (OMCS-aligned task routing)
- [ ] One Mind Coordinate System full implementation across all events
- [ ] Publish One Mind OS as a product/framework

---

## 16. Cost & Sovereignty Analysis

### API Cost Breakdown (Monthly Estimate)

| Component | Model | Est. Calls/Month | Est. Cost |
|---|---|---|---|
| Legacy heartbeat (1min, 24/7) | Claude Sonnet 4.6 | ~43,200 | ~$30-60 |
| Legacy voice responses | Claude Sonnet 4.6 | ~2,000 | ~$5-15 |
| Vision digest (Claude Haiku) | Haiku 4.5 | ~4,320 (if used) | ~$2-5 |
| Gemini Live (voice+vision) | Gemini 2.5 Flash | Per minute of session | Variable |
| ElevenLabs TTS | Turbo v2.5 | ~50K chars/month | ~$5-10 |

**If SmolVLM local**: Vision cost = $0. No API calls for scene description.
**If Gemini Live**: Gemini charges per audio minute + per video frame. Budget separately based on session hours.

### Sovereignty Score

| Layer | Sovereignty | Notes |
|---|---|---|
| Brain (Claude) | ✅ High | API key only, swap to local model later |
| Gateway (OpenClaw) | ✅ Complete | Runs on your hardware |
| NATS mesh | ✅ Complete | Self-hosted JetStream |
| Vision (SmolVLM) | ✅ Complete | Local inference, zero API |
| Vision (Gemini Live) | ⚠️ Partial | Google processes video stream |
| TTS (ElevenLabs) | ⚠️ Partial | Cloud API, swap to Chatterbox for full sovereignty |
| TTS (Chatterbox) | ✅ Complete | Self-hosted, sub-200ms CPU |
| Robotics (RosClaw) | ✅ Complete | ROS2 self-hosted |
| Machine identity (peaq) | ✅ Complete | Onchain, self-sovereign |

### Chatterbox TTS (Recommended ElevenLabs Replacement)
- Self-hosted, sub-200ms on CPU
- Run on Mac Mini or DO droplet
- Drop-in replacement — configure as TTS endpoint in discord-voice plugin
- Eliminates ElevenLabs API dependency entirely

---

## 17. Key Decisions Outstanding

### Decision 1: Vision Path for Active Sessions
**Option A**: NixClaw → Gemini Live → OpenClaw (Google sees your video, realtime)
**Option B**: SmolVLM local → AMBIENT_STATE.md → Claude heartbeat (sovereign, near-realtime)
**Option C**: Both — glasses active = Gemini Live, desk/home = SmolVLM

Recommendation: **Option C**. Use the right tool per context.

### Decision 2: TTS Voice for Legacy
**Option A**: ElevenLabs (current plan) — cloud, low latency, high quality
**Option B**: Chatterbox — self-hosted, sub-200ms CPU, good quality
**Option C**: Sesame CSM-1B — long-term conversational presence layer on Mason nodes

Recommendation: **Chatterbox now, Sesame CSM-1B when it's production-stable**.

### Decision 3: Heartbeat Primary Model
**Option A**: Claude Sonnet 4.6 for all heartbeats — highest quality, higher cost
**Option B**: Claude Haiku 4.5 for heartbeats, Sonnet for active conversations

Recommendation: **Option B**. Haiku handles `HEARTBEAT_OK` decisions fine. Reserve Sonnet for real conversations.

### Decision 4: OpenClaw Voice Channel
**Current**: Discord voice plugin (VAD + ElevenLabs)
**Alternative**: Voice Call plugin (Twilio/Telnyx) — phone number, works anywhere
**Alternative**: NixClaw audio-only mode — always-on background audio on Mason

Recommendation: **Discord voice as primary command center + NixClaw audio-only on Mason for ambient**.

### Decision 5: Robot Acquisition Timeline
G1 at $16,000 is a significant capital decision. The software stack (RosClaw + DimensionalOS + peaq SDK) can be fully developed and tested in Gazebo simulation before hardware arrives. Simulation environment is available via RosClaw Docker Compose.

---

## 18. Reference Links

### Open Source Projects

| Project | URL | Status |
|---|---|---|
| OpenClaw | github.com/openclaw/openclaw | Active, 302k+ stars |
| VisionClaw | github.com/sseanliu/VisionClaw | Working, Feb 2026 |
| NixClaw | github.com/arniesaha/NixClaw | Working, Feb 2026 |
| RosClaw | github.com/PlaiPin/rosclaw | Active |
| SmolVLM-256M | huggingface.co/HuggingFaceTB/SmolVLM-256M-Instruct | Production |
| DimensionalOS | dimensionalos.com | Active, March 2026 |

### Documentation

| Resource | URL |
|---|---|
| OpenClaw Docs | docs.openclaw.ai |
| OpenClaw Heartbeat | docs.openclaw.ai/gateway/heartbeat |
| OpenClaw Skills | docs.openclaw.ai/tools/skills |
| OpenClaw Voice Call Plugin | docs.openclaw.ai/plugins/voice-call |
| OpenClaw Model Providers | docs.openclaw.ai/concepts/model-providers |
| OpenClaw DeepWiki | deepwiki.com/openclaw/openclaw |
| Gemini Live API | ai.google.dev/gemini-api/docs/live-api |
| ClawHub Skills Registry | clawhub.ai |

### Skills to Install

```bash
# Core skills
openclaw skills install peaq-robotics
openclaw skills install unitree-robot
openclaw skills install robot
openclaw add @steipete/camsnap

# Community skills worth reviewing
openclaw skills install agent-brain          # local SQLite memory
openclaw skills install agent-autopilot      # heartbeat-driven task execution
openclaw skills install clawd-presence       # physical presence display
```

### Key Config Files Location

```
~/.openclaw/openclaw.json              # Main gateway config
~/.openclaw/workspace/SOUL.md          # Legacy identity
~/.openclaw/workspace/HEARTBEAT.md     # Heartbeat checklist
~/.openclaw/workspace/AGENTS.md        # Tool access + routing
~/.openclaw/workspace/AMBIENT_STATE.md # Live vision + NATS state (generated)
~/.openclaw/workspace/skills/          # Custom skills directory
~/.openclaw/logs/                      # Gateway + daemon logs
~/.openclaw/agents/legacy/sessions/    # Session JSONL history
```

---

## Appendix A — Full Architecture Summary

```
AMBIENT COMPUTING STACK — ONE MIND OS / LEGACY

INPUT SOURCES
├─ Voice (always-on VAD)     → Gemini Live / Discord Voice Plugin
├─ Vision: glasses 1fps      → Gemini Live (continuous)
├─ Vision: local cameras     → SmolVLM 5fps → AMBIENT_STATE.md
├─ NATS events               → legacy-pulse → AMBIENT_STATE.md
├─ Session context           → legacy-session-start → NATS
└─ Robot telemetry           → robot.g1.* → NATS

PROCESSING LAYER
└─ OpenClaw Gateway :18789
   ├─ Legacy (Claude Sonnet 4.6) — brain + decision
   ├─ SOUL.md — identity
   ├─ HEARTBEAT.md — background monitoring (1min)
   ├─ AMBIENT_STATE.md — vision + NATS digest
   ├─ Skills — tools, automation, search
   └─ Workspace memory — persistent session history

ACTION OUTPUTS
├─ Voice response (ElevenLabs / Chatterbox) → speakers
├─ NATS publish → system events
├─ Discord message → command center
├─ RosClaw → rosbridge → ROS2 → Unitree G1
├─ peaq SDK → machine identity → onchain
└─ 800+ OpenClaw skills → any digital action

INFRASTRUCTURE
├─ Tailscale mesh (Home + Farm + Office + DO droplet)
├─ NATS JetStream primary (DO bare metal)
├─ NATS leaf nodes (each physical location)
├─ Mason Android nodes (MASC09T + C210)
└─ Ray-Ban Meta glasses (Meta DAT SDK)
```

---

*Document generated March 2026. All referenced open source projects verified operational. Architecture based on real implementations confirmed working in the wild.*

*Next review: when Unitree G1 arrives or Anthropic ships realtime audio API.*


---

## Related

- [[00-24 SO (Sovereign Operations)/01 Intelligence (AI, Memory & Tools) (The Brain)/01 Intelligence|01 Intelligence]]
