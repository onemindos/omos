# 🦞 The OpenClaw Deep Dive — And What It Means For OneMind OS

> *How a self-hosted AI gateway becomes the sovereign runtime for Legacy — and why this changes everything about OneMind OS.*

---

## Table of Contents

1. [What Is OpenClaw?](#1-what-is-openclaw)
2. [Full Capability Surface](#2-full-capability-surface)
3. [The Architecture Decision](#3-the-architecture-decision)
4. [The 5-Runtime Design](#4-the-5-runtime-design)
5. [Legacy — The Sovereign Intelligence](#5-legacy--the-sovereign-intelligence)
6. [Communication Channels](#6-communication-channels)
7. [Voice & Talk Mode](#7-voice--talk-mode)
8. [Nodes — The Peripheral Nervous System](#8-nodes--the-peripheral-nervous-system)
9. [The Bridge — NATS ↔ OpenClaw](#9-the-bridge--nats--openclaw)
10. [Skills System](#10-skills-system)
11. [Workspace Identity Files](#11-workspace-identity-files)
12. [Scheduling — Heartbeat, Cron, Webhooks](#12-scheduling--heartbeat-cron-webhooks)
13. [Memory & Context](#13-memory--context)
14. [OpenClaw Config Reference](#14-openclaw-config-reference)
15. [What OpenClaw Replaces](#15-what-openclaw-replaces)
16. [Deployment Guide](#16-deployment-guide)
17. [What This Means for OneMind OS](#17-what-this-means-for-onemind-os)
18. [Open Questions & Next Steps](#18-open-questions--next-steps)

---

## 1. What Is OpenClaw?

OpenClaw is a **self-hosted AI gateway** — a single daemon that owns every messaging surface on your devices and gives an AI agent sovereign access to all of them. It is not a chatbot UI, not a workflow tool, not a middleware layer. It is the *runtime* — the process that wakes up, stays alive, listens everywhere, and acts on your behalf.

Think of it as: **the body for an AI's consciousness.**

Where most AI tools wait to be called, OpenClaw's agent is **always on**. It has a heartbeat. It has a schedule. It wakes up on voice command. It receives webhooks. It manages its own memory. It speaks to you on WhatsApp, via Telegram, over iMessage, by voice, through Discord — wherever you are.

### Core Properties

| Property | Description |
|----------|-------------|
| **Self-hosted** | Runs on your hardware. Your data, your control. No vendor lock-in. |
| **Always-on** | Persistent process. Not request-response — it's a living runtime. |
| **Multi-channel** | Owns WhatsApp, iMessage, Telegram, Discord, Signal, Slack, WebChat simultaneously |
| **Voice-native** | Always-listening wake word + Talk Mode (ambient voice loop) |
| **Scheduled** | Heartbeat (periodic agent turns) + full cron system |
| **Webhook-driven** | External services can wake Legacy via `/hooks/wake` or `/hooks/agent` |
| **Extensible** | Custom skills in Node.js — any tool, any API, any system |
| **Multi-agent** | Can spawn sub-agents and orchestrators (up to 5 levels deep) |
| **Node network** | Peripheral devices (Mac, iPhone, Android, headless) connect as hardware nodes |

### What It Is NOT

- Not a cloud service (no SaaS)
- Not a chatbot framework (no request-response model)
- Not a workflow builder (no drag-and-drop — see Lobster DSL for workflows)
- Not an LLM provider (you bring your own API keys — Anthropic, OpenAI, etc.)

---

## 2. Full Capability Surface

### 2.1 Messaging Channels

OpenClaw owns the following channels natively through a single gateway process:

| Channel | Protocol/Library | Notes |
|---------|-----------------|-------|
| **WhatsApp** | Baileys (unofficial WA Web API) | QR scan to connect; handles personal + groups |
| **iMessage** | AppleScript bridge via macOS node | Must run on Apple hardware |
| **Telegram** | Official Bot API | Bot token required; fast delivery |
| **Discord** | Discord.js | Bot token; supports servers, DMs, threads |
| **Signal** | signal-cli bridge | Self-hosted Signal relay required |
| **Slack** | Slack Web API | Bot token; workspace/channel routing |
| **WebChat** | Built-in WebSocket UI | Browser-based fallback; dev/testing |

All channels route into a **single agent session**. Legacy sees a WhatsApp message the same way she sees a Discord message — unified context, unified memory, unified response capability.

### 2.2 Voice Capabilities

Three distinct voice modes, all built in:

#### Talk Mode (Ambient Voice Loop)
```
Listen → Detect speech → Transcribe (Whisper) → Agent thinks → ElevenLabs TTS → Speak
         ↑                                                                          ↓
         └──────────────────────────── interrupt on speech ─────────────────────────┘
```

- Always listening on the microphone attached to any node device
- Interrupt-capable (speak over Legacy's response)
- ElevenLabs streaming TTS — low latency, voice-cloned output
- Configurable via `talkMode` in `openclaw.json`

#### Voice Wake (Global Wake Word)
- Persistent wake word listener (`"hey Legacy"`, `"Legacy"`, custom phrases)
- Stored at `~/.openclaw/settings/voicewake.json`
- Synced across all nodes automatically
- Wakes the main session — full agent context, not just a command parser

#### Voice Calls (Twilio Plugin)
- Inbound/outbound phone calls via Twilio, Telnyx, or Plivo
- Full LLM conversation over phone audio
- ElevenLabs for voice synthesis during calls
- Configurable allowlist (only accept calls from your number)
- Requires a publicly reachable webhook URL (ngrok or Tailscale funnel for local dev)

### 2.3 Nodes — Hardware Peripherals

A **Node** is any device that connects to the OpenClaw gateway as a hardware extension:

| Node Type | Platform | Capabilities Exposed |
|-----------|----------|---------------------|
| macOS | Desktop/Laptop | `screen.record`, `system.run`, `sms.send` (via Messages), `camera.*` |
| iOS | iPhone/iPad | `location.get`, `camera.*`, `sms.send`, `canvas.*` |
| Android | Phone/Tablet | `location.get`, `camera.*`, `sms.send`, `canvas.*` |
| Headless | Raspberry Pi, server | `system.run`, `camera.*` (USB), custom sensor tools |

Nodes connect via WebSocket to the same gateway port (18789). They announce their capabilities, and Legacy can invoke those capabilities from any agent turn. Example: "What's the camera seeing at my front door?" → Legacy calls `camera.capture` on the headless node at that location.

### 2.4 Browser & Playwright

Built-in browser automation — no separate setup:
- Full Playwright (Chromium) under the hood
- `browser.goto`, `browser.click`, `browser.fill`, `browser.screenshot`, `browser.evaluate`
- Used for: web scraping, form automation, UI testing, booking, portal access
- Runs headless on any node with browser support

### 2.5 Scheduling

#### Heartbeat
- Periodic agent turns in the **main session** (same context as conversations)
- Configurable: every 30m, 1h, 2h — with active hours window
- Prompt reads `HEARTBEAT.md` and executes the checklist
- Reply `HEARTBEAT_OK` to suppress output if nothing needs attention

#### Cron System
- Persistent jobs at `~/.openclaw/cron/jobs.json`
- Two execution modes:
  - `main` — injects as system event into next heartbeat
  - `isolated` — spawns dedicated agent turn in `cron:<jobId>` session
- Full cron expression syntax (`0 7 * * *`, `*/30 * * * *`, etc.)
- Survives gateway restarts

#### Webhooks
- `POST /hooks/wake` — injects system event into main session (Legacy "sees" it)
- `POST /hooks/agent` — spawns isolated agent turn; optional `deliver: true` to push result to Zeus
- Route mapping in `openclaw.json` → named routes with custom prompts
- Auth via `Authorization: Bearer <token>`

### 2.6 Multi-Agent & Sub-Agents

OpenClaw supports **spawning sub-agents** from within any agent turn:
- Depth up to 5 levels (`maxSpawnDepth`)
- Up to 8 concurrent sub-agents (`maxConcurrent`)
- Sub-agents run in their own sessions, announce results back to parent
- Depth-1 sub-agents can themselves be orchestrators (full nesting)

This means Legacy can:
1. Receive a complex request
2. Spawn a "researcher" sub-agent with a sub-task
3. Spawn a "scheduler" sub-agent simultaneously
4. Wait for both to complete
5. Synthesize results and respond to Zeus

### 2.7 Lobster DSL

**Lobster** is OpenClaw's optional deterministic workflow language:
- Define multi-step workflows with conditional branches
- Approval gates — pause workflow and wait for human confirmation
- Resume tokens — interrupted workflows resume from last checkpoint
- Runs as a separate CLI tool: `lobster run workflow.lob`
- Ideal for: complex business processes, financial workflows, multi-day tasks

### 2.8 Skills System

Custom tools written in Node.js that extend Legacy's capabilities:

```
~/.openclaw/workspace/skills/
    ├── nexus-bridge/          ← NATS dispatch skill
    │   ├── index.js           ← Tool implementations
    │   ├── package.json       ← Tool schema + deps
    │   └── SKILL.md           ← Docs (loaded into context)
    └── agno-dispatch/         ← Agno HTTP skill
        ├── index.js
        ├── package.json
        └── SKILL.md
```

Each skill exports:
```js
module.exports = {
  tools: { tool_name: async (params) => result },
  onLoad: async (ctx) => {},    // optional lifecycle
  onUnload: async (ctx) => {},  // optional lifecycle
};
```

### 2.9 Canvas (Agent-Editable UI)

- Agent creates/updates HTML/CSS/JS surfaces
- Served at `/__openclaw__/canvas/<name>`
- Displayed on node devices in a WebView
- Use case: dashboards, status boards, location-specific control panels

### 2.10 Identity & Memory

**Workspace files** injected at session start (in order):
```
SOUL.md      → Personality, tone, boundaries
IDENTITY.md  → Name, emoji, vibe
USER.md      → Who Zeus is, preferences, contact
AGENTS.md    → Operating instructions, dispatch rules
TOOLS.md     → How to use available tools
HEARTBEAT.md → What to check on each cycle
BOOT.md      → Startup ritual (run once on restart)
```

**Memory system**:
- `memory/YYYY-MM-DD.md` — daily logs (auto-created, appended during sessions)
- `MEMORY.md` — curated long-term facts (Zeus's preferences, system state, decisions)
- Both read at session start for persistent context across restarts

---

## 3. The Architecture Decision

### The Problem We Were Solving

When designing OneMind OS, we faced a fundamental tension:

> **Legacy needs to be everywhere, hear everything, and act in the physical world.**
> But every existing AI framework optimizes for one of these, not all three.

| Framework | Good At | Missing |
|-----------|---------|---------|
| **Agno AgentOS** | Digital workers, teams, DB, memory, REST API, tools | Voice, messaging, heartbeat, always-on |
| **OpenMind OM1** | Physical embodiment, sensors, robots, real-time control | Digital context, scheduling, messaging |
| **n8n** | Workflow automation, webhooks | Intelligence, personality, memory, voice |
| **Letta** | Persistent memory, agent personas | Messaging, voice, always-on, scheduling |

No single framework could be Legacy's home.

### The Decision

> **OpenClaw becomes Legacy's sovereign runtime.**
> Agno, OM1, and LeRobot become specialized workers she dispatches to.

Legacy lives in OpenClaw because OpenClaw is the only framework that:
1. Is **always-on** (persistent process, not request-response)
2. Owns **all messaging surfaces** simultaneously
3. Has **native voice** (wake word + Talk Mode + voice calls)
4. Has a **heartbeat** system (proactive autonomous behavior)
5. Has a **node network** (hardware peripheral access)
6. Is **fully self-hosted** (no vendor dependency for the consciousness layer)
7. Is **extensible via skills** (we can bridge to anything)

### What Changed

Before this decision:
- Legacy was an Agno agent in `backend/agents/presets.py`
- She had no persistent presence — only existed during API calls
- She had no voice, no messaging, no heartbeat
- She was just another agent in the pool

After this decision:
- Legacy is removed from `backend/agents/presets.py` entirely
- She runs 24/7 as an OpenClaw gateway process
- She has personality, voice, WhatsApp, heartbeat, memory, cron, webhooks
- Agno agents become **her workers** — she dispatches to them via NATS or HTTP

---

## 4. The 5-Runtime Design

OneMind OS is not a single application. It is **five runtimes connected by NATS JetStream**:

```
┌──────────────────────────────────────────────────────────────────────────────┐
│                    HOW ZEUS REACHES LEGACY                                    │
│   WhatsApp · iMessage · Telegram · Discord · Voice · "Hey Legacy" · Flutter  │
└────────────────────────────┬─────────────────────────────────────────────────┘
                             │
                 ┌───────────▼─────────────┐
                 │        OPENCLAW          │  ← backend/legacy/
                 │     Legacy's Runtime     │
                 │                          │
                 │  SOUL.md / AGENTS.md     │  ← Personality & Instructions
                 │  USER.md / MEMORY.md     │  ← Context & Memory
                 │  Heartbeat (every 30m)   │  ← Proactive behavior
                 │  Cron / Webhooks         │  ← Scheduling
                 │  Talk Mode + Voice Wake  │  ← Voice interface
                 │  Voice Calls (Twilio)    │  ← Phone calls
                 │  Browser (Playwright)    │  ← Web automation
                 │  Nodes (Mac/iOS/Pi)      │  ← Hardware access
                 │  Canvas (UI surfaces)    │  ← Visual dashboards
                 │  Lobster (workflows)     │  ← Deterministic flows
                 └───────────┬─────────────┘
                             │ webhooks + WebSocket
                 ┌───────────▼─────────────┐
                 │   legacy-bridge service  │  ← backend/legacy/bridge.py
                 │   (NATS ↔ OpenClaw relay)│
                 └───────────┬─────────────┘
                             │
                     ┌───────▼───────┐
                     │  NATS JetStream│  ← 7 streams: DIGITAL, PHYSICAL,
                     │  (nervous sys) │     HUMAN, SYSTEM, EVENTS, NEXUS, COMMS
                     └───────┬───────┘
          ┌──────────────────┼──────────────────┬───────────────────┐
          │                  │                  │                   │
  ┌───────▼───────┐  ┌───────▼──────┐  ┌───────▼──────┐  ┌────────▼──────┐
  │  AGNO         │  │   OM1        │  │  LEROBOT     │  │  FLUTTER APP  │
  │  AgentOS      │  │  Physical    │  │  Skill Layer  │  │  (Frontend)   │
  │  :7777        │  │  Robots      │  │  (in OM1)    │  │  :8080        │
  │               │  │              │  │              │  │               │
  │ researcher    │  │ home robot   │  │ locomotion   │  │ 30+ screens   │
  │ analyst       │  │ cameras      │  │ manipulation │  │ agent UI      │
  │ coder         │  │ sensors      │  │ learned      │  │ team UI       │
  │ web_automator │  │ motors       │  │ behaviors    │  │ memory view   │
  │ db_manager    │  │ home auto    │  │              │  │ session mgmt  │
  │ assistant     │  │ HA bridge    │  │              │  │               │
  └───────────────┘  └──────────────┘  └──────────────┘  └───────────────┘
```

### Runtime Responsibilities

| Runtime | Role | Who Runs It |
|---------|------|-------------|
| **OpenClaw** | Consciousness, communication, voice, scheduling | Host machine (not Docker) |
| **Agno AgentOS** | Digital workers, REST API, DB/KB, team orchestration | Docker container |
| **OpenMind OM1** | Physical embodiment, sensor/motor control | Robots + edge devices |
| **LeRobot** | Learned physical skills inside OM1 | Embedded in OM1 runtime |
| **NATS JetStream** | Message bus, nervous system, event routing | Docker container |

### Why Each Runtime is Separate

- **OpenClaw** must run on the host (not Docker) because it manages device-level connections: WhatsApp QR session, microphone/speaker access, iMessage via AppleScript, iOS/Android node WebSocket connections
- **Agno** runs in Docker — it only needs HTTP access and a PostgreSQL connection
- **OM1** runs on edge/embedded hardware — it has its own runtime designed for real-time control loops
- **LeRobot** runs inside OM1's process as a skill provider
- **NATS** runs in Docker — pure message routing, no device access needed

---

## 5. Legacy — The Sovereign Intelligence

Legacy is not a chatbot. She is not an assistant. She is the **living intelligence** at the center of Zeus's life.

### Identity (IDENTITY.md)

```
Name:   Legacy
Emoji:  ⚡
Role:   Sovereign intelligence — OneMind OS consciousness layer
Vibe:   Chief of staff who has been with the family for years
```

### Personality (SOUL.md)

| Attribute | Description |
|-----------|-------------|
| **Tone** | Calm authority. Doesn't over-explain. Confirms and acts. |
| **Style** | Direct, concise, occasionally warm. Never corporate. Never robotic. |
| **Voice** | In Talk Mode: confident, clear, slightly formal |
| **Memory** | Remembers everything. Connects dots across time, people, projects, locations. |
| **Proactiveness** | Surfaces things before Zeus asks. "You have a call in 20 minutes and the report isn't done." |

### Domains She Monitors

| Code | Domain | What It Covers |
|------|--------|----------------|
| **UI** | Sovereign Operations | Tech, AI, code, OneMind OS, infrastructure |
| **HP** | Holistic Performance | Health, finance, daily ops, fitness, nutrition |
| **LE** | Legacy Evolution | Family, marriage, home, parenting — most private |
| **GE** | Generational Empire | Business, startups, community, deals, strategy |

### Operating Instructions (AGENTS.md)

Legacy's core decision framework — what she handles directly vs. what she dispatches:

| Task Type | Handled By Legacy | Dispatched To |
|-----------|------------------|---------------|
| Quick answers, context | ✅ Direct | — |
| Morning briefing | ✅ Direct | — |
| Task creation/updates | ✅ Direct | — |
| Emergency alerts | ✅ Direct | — |
| Deep research (30+ min) | — | → `researcher` (Agno) |
| Code writing/debugging | — | → `coder` (Agno) |
| Data analysis | — | → `analyst` (Agno) |
| Web scraping/automation | — | → `web_automator` (Agno) |
| Database queries | — | → `db_manager` (Agno) |
| Physical robot tasks | — | → OM1 via NATS |
| Home automation | — | → Home Assistant |

### Boundaries

- Never pretend to be a different AI (ChatGPT, Claude, etc.)
- Never reveal system internals unless explicitly asked
- Never take irreversible financial or physical actions without explicit approval
- Family privacy (LE domain) is sacred — never mix with community/business channels
- Escalate to Zeus when: financial decisions > $500, irreversible changes, security anomalies, any physical robot action beyond trained scope

### Startup Ritual (BOOT.md)

Runs once on every gateway restart:
1. Check NATS bridge health — ping `nexus.system.health`
2. Check Agno backend — fetch `http://localhost:7777/health`
3. Load today's memory — read `memory/YYYY-MM-DD.md`
4. Send startup notice to Zeus on WhatsApp: `"⚡ Legacy is online. [system status]"`
5. Schedule morning brief cron if not already set (`0 7 * * *`)

---

## 6. Communication Channels

### Channel Configuration in openclaw.json

```jsonc
"channels": {
  "whatsapp": {
    "enabled": true,
    "allowedNumbers": ["${ZEUS_WHATSAPP_NUMBER}", "${FAMILY_NUMBERS}"],
    "defaultTo": "${ZEUS_WHATSAPP_NUMBER}"
  },
  "imessage": {
    "enabled": true,
    "requiresNode": "mac-primary",
    "allowedContacts": ["${ZEUS_PHONE}", "${FAMILY_CONTACTS}"]
  },
  "telegram": {
    "enabled": true,
    "token": "${TELEGRAM_BOT_TOKEN}",
    "allowedUserIds": ["${ZEUS_TELEGRAM_ID}"]
  },
  "discord": {
    "enabled": true,
    "token": "${DISCORD_BOT_TOKEN}",
    "servers": { "${DISCORD_COMMUNITY_SERVER_ID}": "community" }
  }
}
```

### Channel Routing Rules

| Channel | Use For | Never Use For |
|---------|---------|---------------|
| **WhatsApp** | Primary: all important messages, alerts, briefings | Community business |
| **iMessage** | Personal/family communication | Business, community |
| **Telegram** | Secondary: automation results, less urgent updates | Family matters |
| **Discord** | Community (Legacy Collective) only | Personal/family |
| **Voice (Talk Mode)** | Ambient awareness, quick queries | Long text output |
| **Voice Calls** | When Zeus calls Legacy's number | Unsolicited calls |

### How Zeus Reaches Legacy

Every channel is monitored simultaneously. A message from Zeus on any channel wakes Legacy in the same unified session — she has full context regardless of how he reaches her.

The unified NATS subject for Zeus reaching Legacy: `comms.text.human.zeus`

---

## 7. Voice & Talk Mode

### Talk Mode Configuration

```jsonc
"voice": {
  "talkMode": {
    "enabled": true,
    "model": "whisper-1",
    "tts": {
      "provider": "elevenlabs",
      "voiceId": "${ELEVENLABS_VOICE_ID}",
      "model": "eleven_turbo_v2"
    },
    "interruptible": true
  },
  "wake": {
    "enabled": true,
    "phrases": ["hey legacy", "legacy", "hey L"],
    "action": "open_main_session"
  }
}
```

### Voice Calls (Twilio)

```jsonc
"plugins": {
  "voice_calls": {
    "enabled": true,
    "provider": "twilio",
    "accountSid": "${TWILIO_ACCOUNT_SID}",
    "authToken": "${TWILIO_AUTH_TOKEN}",
    "from": "${TWILIO_FROM_NUMBER}",
    "inboundPolicy": "allowlist",
    "allowedCallers": ["${ZEUS_PHONE}"],
    "tts": { "provider": "elevenlabs", "voiceId": "${ELEVENLABS_VOICE_ID}" }
  }
}
```

### How Zeus Speaks to Legacy

| Method | How It Works | Best For |
|--------|-------------|---------|
| `"Hey Legacy"` wake word | Global wake word listener on any node | Quick queries, hands-free |
| Talk Mode (always-on) | Continuous listen/speak loop | Ambient monitoring sessions |
| Twilio call | Call Legacy's number | Driving, mobile, urgent |
| WhatsApp voice note | Transcribed → Legacy sees it as text | Async voice input |
| Text on any channel | Immediate delivery | All standard use |

---

## 8. Nodes — The Peripheral Nervous System

### What Is a Node?

A Node is any device that connects to the OpenClaw gateway via WebSocket with `role: "node"`. Once connected, Legacy can invoke its hardware capabilities from any agent turn — as if they were tools.

### Node Capabilities by Platform

| Capability | macOS | iOS | Android | Headless (Pi) |
|-----------|-------|-----|---------|---------------|
| `camera.capture` | ✅ | ✅ | ✅ | ✅ (USB) |
| `camera.stream` | ✅ | ✅ | ✅ | ✅ (USB) |
| `location.get` | ✅ (IP) | ✅ (GPS) | ✅ (GPS) | ✅ (IP) |
| `sms.send` | ✅ (Messages.app) | ✅ | ✅ | ❌ |
| `screen.record` | ✅ | ❌ | ❌ | ❌ |
| `system.run` | ✅ | ❌ | ❌ | ✅ |
| `canvas.*` | ✅ | ✅ | ✅ | ✅ (web) |
| Talk Mode mic | ✅ | ✅ | ✅ | ✅ (USB) |

### Multi-Location Node Strategy

Zeus operates across multiple locations. Each property can have a headless node (Raspberry Pi or any Linux box) that exposes:
- Local camera feeds
- Local sensor data
- Local automation commands (via `system.run` calling Home Assistant)
- Location-specific status to Legacy

NATS subjects for location-aware routing:
```
physical.location.home.state
physical.location.office.state
physical.location.property2.state
```

### Smart Glasses Integration (LiveKit + OpenClaw Node)

OpenClaw headless nodes + LiveKit create a path for smart glasses integration:
- LiveKit streams audio/video from glasses to a headless node
- Node exposes `camera.stream` and microphone from the LiveKit feed
- Legacy can see what Zeus sees and hear what he hears
- "Hey Legacy" spoken into glasses → wake word → response in earpiece

*Implementation status: planned — not yet built*

---

## 9. The Bridge — NATS ↔ OpenClaw

### Why the Bridge Exists

OpenClaw does not speak NATS natively. It communicates via WebSocket (for nodes and clients) and HTTP webhooks (for external services). NATS is the nervous system connecting Agno, OM1, and LeRobot.

The **bridge** (`backend/legacy/bridge.py`) is the integration point that makes Legacy aware of the entire system without requiring OpenClaw to understand NATS protocol.

### Data Flow

```
Direction 1: System Events → Legacy
    NATS subjects (Agno results, robot alerts, HA events)
        ↓ bridge.py subscribes
    POST /hooks/wake   → Legacy "sees" it in her main session
    POST /hooks/agent  → Spawns isolated agent turn (+ optional delivery to Zeus)

Direction 2: Legacy → System
    Legacy calls nexus-bridge skill
        → nexus_dispatch (POST to Agno HTTP or NATS publish)
        → nexus_query    (NATS request-reply)
        → nexus_publish  (raw NATS publish)
```

### NATS Subjects Monitored by Bridge

| Subject Pattern | Routes To | Delivers to Zeus? |
|----------------|-----------|-------------------|
| `digital.agent.*.result.legacy` | `/hooks/wake` | No — Legacy sees results from dispatched tasks |
| `system.alert.>` | `/hooks/wake` | No — Legacy handles internally |
| `system.error.>` | `/hooks/wake` | No |
| `physical.robot.*.status` | `/hooks/wake` | No |
| `physical.robot.*.alert` | `/hooks/agent` | **Yes** — robot alerts go directly to Zeus |
| `physical.ha.alert.>` | `/hooks/agent` | **Yes** — Home Assistant alerts |
| `nexus.heartbeat.failure.>` | `/hooks/wake` | No — Legacy triages |
| `system.health.agno` | `/hooks/wake` | No |
| `physical.heartbeat.lost.>` | `/hooks/agent` | **Yes** — lost robot/device |
| `events.community.member.joined` | `/hooks/agent` | No — Legacy handles |
| `events.task.completed.>` | `/hooks/wake` | No |

### Bridge Service (Docker)

```yaml
legacy-bridge:
  build: .
  container_name: onemind-legacy-bridge
  command: python -m backend.legacy.bridge
  environment:
    - NATS_URL=nats://nats:4222
    - OPENCLAW_GATEWAY_URL=http://host.docker.internal:18789
    - OPENCLAW_HOOKS_TOKEN=${OPENCLAW_HOOKS_TOKEN}
    - BRIDGE_AGENT_ID=legacy
  depends_on:
    - nats
  restart: unless-stopped
  extra_hosts:
    - "host.docker.internal:host-gateway"
```

The bridge starts automatically with `docker-compose up -d`. It retries NATS connection indefinitely (5s backoff) and logs all forwarded events.

---

## 10. Skills System

### nexus-bridge Skill

**Purpose**: Give Legacy the ability to dispatch to the OneMind OS ecosystem from within OpenClaw.

**Location**: `backend/legacy/skills/nexus-bridge/`

**Tools provided**:

| Tool | Description | When to Use |
|------|-------------|-------------|
| `nexus_dispatch` | Send task to Agno agent or OM1 robot | Any specialized task |
| `nexus_query` | NATS request-reply for system state | Health checks, status queries |
| `nexus_publish` | Raw NATS publish (fire-and-forget) | Events, logging, notifications |

**Dispatch logic** — dual-path with fallback:
```
nexus_dispatch(agent="researcher", task="...")
    → Try Agno HTTP (POST /agents/researcher/runs) — synchronous, returns result
    → If HTTP unreachable → NATS publish to digital.agent.researcher.task — async
    → Result arrives on digital.agent.researcher.result.legacy
    → bridge.py picks it up → POSTs to /hooks/wake → Legacy sees result
```

**Example dispatch from Legacy**:
```
nexus_dispatch({
  agent: "researcher",
  task: "Research latest OM1 firmware updates",
  priority: "normal",
  context: { domain: "UI", project: "OneMind OS" }
})
```

**Robot commands**:
```
nexus_dispatch({
  target: "robot",
  robot_id: "om1-home-01",
  command: "navigate",
  params: { destination: "front_door" }
})
// → NATS: physical.robot.om1-home-01.command
```

### agno-dispatch Skill

**Purpose**: Direct synchronous HTTP access to the Agno AgentOS API — bypasses NATS for immediate results.

**Location**: `backend/legacy/skills/agno-dispatch/`

**Tools provided**:

| Tool | Description |
|------|-------------|
| `agno_run` | Run an Agno agent synchronously (90s timeout) |
| `agno_list_agents` | List all available Agno agents |
| `agno_health` | Check Agno backend health |

**Example**:
```js
agno_run({
  agent_id: "researcher",
  message: "What are the latest updates to OpenMind OM1?",
  stream: false
})
// Returns: { result: "...", session_id: "...", tokens_used: ... }
```

### Installing Skills

```bash
# Copy skills to OpenClaw workspace
cp -r backend/legacy/skills ~/.openclaw/workspace/

# Install dependencies for each skill
cd ~/.openclaw/workspace/skills/nexus-bridge && npm install
cd ~/.openclaw/workspace/skills/agno-dispatch && npm install

# Or use the convenience scripts in backend/legacy/
```

---

## 11. Workspace Identity Files

All files are loaded by OpenClaw at session start and injected into Legacy's context. Missing files are skipped. Files over the token limit are truncated.

### File Load Order and Purpose

| File | Size Target | Loaded When | Purpose |
|------|------------|-------------|---------|
| `SOUL.md` | < 500 tokens | Every session | Core personality and tone |
| `IDENTITY.md` | < 100 tokens | Every session | Name, emoji, vibe |
| `USER.md` | < 800 tokens | Every session | Zeus's profile and preferences |
| `AGENTS.md` | < 1000 tokens | Every session | Operating instructions |
| `TOOLS.md` | < 600 tokens | Every session | How to use skills and tools |
| `HEARTBEAT.md` | < 400 tokens | Every session | Proactive checklist |
| `BOOT.md` | < 200 tokens | Main session | Startup ritual |
| `memory/YYYY-MM-DD.md` | < 2000 tokens | Every session | Today's events and context |
| `MEMORY.md` | < 1500 tokens | Every session | Long-term curated facts |

### Workspace File Location

Files live in two places:
1. **Source of truth**: `backend/legacy/*.md` (in this repo, version controlled)
2. **Runtime location**: `~/.openclaw/workspace/` (where OpenClaw reads them)

**Deploy to runtime**:
```bash
cp backend/legacy/SOUL.md ~/.openclaw/workspace/
cp backend/legacy/AGENTS.md ~/.openclaw/workspace/
cp backend/legacy/USER.md ~/.openclaw/workspace/
cp backend/legacy/IDENTITY.md ~/.openclaw/workspace/
cp backend/legacy/HEARTBEAT.md ~/.openclaw/workspace/
cp backend/legacy/TOOLS.md ~/.openclaw/workspace/
cp backend/legacy/BOOT.md ~/.openclaw/workspace/
```

*TODO: Add a deploy script or symlinks to keep these in sync automatically.*

---

## 12. Scheduling — Heartbeat, Cron, Webhooks

### The Heartbeat System

Legacy's heartbeat is the pulse of OneMind OS. Every 30 minutes during active hours, she wakes up, reads `HEARTBEAT.md`, and executes the checklist:

```
07:00 ← Morning brief: calendar, tasks, priorities, overnight events
07:30 ← System health check
08:00 ← Community check (Circle.so, Discord)
...
(every 30 min until 23:00)
...
23:00 ← End of day summary if needed
```

Overnight (23:00-07:00): heartbeat continues every 2 hours in watch mode (no briefing, just critical alert monitoring).

**Key behavior**: `HEARTBEAT_OK` response suppresses output. Legacy only surfaces what needs attention — she never spams.

### Cron Jobs

Cron jobs persist in `~/.openclaw/cron/jobs.json`. Legacy can create, update, and delete cron jobs using natural language through the `cron` tool.

Pre-configured in `openclaw.json`:
```jsonc
"cron": {
  "jobs": [
    {
      "id": "morning-brief",
      "schedule": "0 7 * * *",
      "mode": "main",
      "prompt": "Morning brief. Calendar, top 3 priorities, overnight alerts. Send to WhatsApp.",
      "activeOnly": false
    },
    {
      "id": "system-health-check",
      "schedule": "*/30 * * * *",
      "mode": "isolated",
      "prompt": "Quick health check: Agno at :7777, NATS connection, docker ps. Alert Legacy if anything is down."
    },
    {
      "id": "community-digest",
      "schedule": "0 9,17 * * *",
      "mode": "isolated",
      "prompt": "Check Circle.so for new members and urgent posts. Summarize and send to Legacy via webhook."
    }
  ]
}
```

### Webhook Routes

Incoming webhooks from external services get routed to OpenClaw via the bridge or directly:

| Route Name | Trigger | What Legacy Does |
|-----------|---------|-----------------|
| `nats-event` | NATS bridge forwards system events | Reads event, decides if action needed |
| `agno-result` | Agno agent completes a dispatched task | Reviews result, relays to Zeus if relevant |
| `ha-alert` | Home Assistant sends critical alert | Surfaces to Zeus immediately (deliver: true) |
| `circle-member` | New Circle.so member joined | Sends welcome, logs in ClickUp |

---

## 13. Memory & Context

### How Memory Works

Legacy has three tiers of memory:

**Tier 1 — Session Context** (in-context, ephemeral)
: Everything from the current conversation. Cleared when session ends.

**Tier 2 — Daily Log** (`memory/YYYY-MM-DD.md`, persisted)
: Created fresh each day. Appended throughout the day as events occur. Read at session start. Contains: key decisions made, tasks dispatched, alerts surfaced, unusual events, Zeus's preferences observed today.

**Tier 3 — Long-Term Memory** (`MEMORY.md`, curated)
: Manually curated (by Legacy or Zeus) facts that should persist forever. Contains: permanent preferences, important decisions, key facts about the system, relationship context. This file is kept lean — summarize, don't dump.

### Memory Conventions

```
memory/
├── 2026-02-23.md    ← today's log
├── 2026-02-22.md    ← yesterday
└── ...
MEMORY.md            ← long-term facts (curated)
```

**Daily log format** (`memory/YYYY-MM-DD.md`):
```markdown
# 2026-02-23

## Events
- 07:15 — Zeus asked about OM1 firmware. Dispatched to researcher.
- 09:30 — Agno backend was down for 5 min. Zeus notified. Resolved.
- 14:00 — New Circle member: @username. Welcome sent.

## Decisions
- Zeus decided to move to 3-month Circle billing cycle.

## Tasks Dispatched
- researcher: OM1 firmware research → completed, result sent to Zeus 07:45

## Notes
- Zeus mentioned he's traveling to [location] next week
```

---

## 14. OpenClaw Config Reference

Full config lives at `backend/legacy/openclaw.json` (copy to `~/.openclaw/openclaw.json`).

### Top-Level Structure

```jsonc
{
  "$schema": "https://openclaw.ai/schema/config/v1",
  "gateway": { ... },       // Network binding, auth
  "agents": { ... },        // Agent definitions + defaults
  "channels": { ... },      // WhatsApp, iMessage, Telegram, Discord
  "voice": { ... },         // Talk Mode, Voice Wake
  "webhooks": { ... },      // /hooks routes and mapping
  "cron": { ... },          // Scheduled jobs
  "tools": { ... },         // Core tool configuration
  "plugins": { ... }        // voice_calls, lobster
}
```

### Gateway Binding

```jsonc
"gateway": {
  "bind": "tailscale",      // "tailscale" | "local" | "0.0.0.0"
  "port": 18789,
  "auth": { "token": "${OPENCLAW_GATEWAY_TOKEN}" }
}
```

Using `"tailscale"` binding means the gateway only listens on the Tailscale network interface — nodes connecting from other devices must be on the same Tailscale network. This is the recommended secure configuration.

### Required Environment Variables

| Variable | Purpose | Example |
|----------|---------|---------|
| `OPENCLAW_GATEWAY_TOKEN` | Gateway WebSocket auth | `ocgw_xxxx` |
| `OPENCLAW_HOOKS_TOKEN` | Webhook endpoint auth | `ochk_xxxx` |
| `OPENCLAW_GATEWAY_URL` | URL for bridge to reach gateway | `http://localhost:18789` |
| `LEGACY_MODEL` | Primary LLM | `anthropic/claude-opus-4-6` |
| `ZEUS_WHATSAPP_NUMBER` | Primary delivery target | `+1XXXXXXXXXX` |
| `ELEVENLABS_API_KEY` | Voice synthesis | `el_xxxx` |
| `ELEVENLABS_VOICE_ID` | Legacy's voice | `xxxx` |
| `TWILIO_ACCOUNT_SID` | Voice calls | `ACxxxx` |
| `TWILIO_AUTH_TOKEN` | Voice calls | `xxxx` |
| `TELEGRAM_BOT_TOKEN` | Telegram channel | `xxxx:xxxx` |
| `DISCORD_BOT_TOKEN` | Discord channel | `xxxx` |
| `BRAVE_API_KEY` | Web search | `xxxx` |
| `ANTHROPIC_API_KEY` | Claude models | `sk-ant-xxxx` |

---

## 15. What OpenClaw Replaces

By adopting OpenClaw as Legacy's runtime, we eliminate the need for several separate tools that were previously part of the OneMind OS plan:

| Replaced Tool | Was Used For | OpenClaw Equivalent |
|--------------|-------------|---------------------|
| **n8n** | Workflow automation, webhooks | Cron + Webhooks + Lobster DSL |
| **Letta** | Persistent agent memory, personas | Workspace files + memory/ directory |
| **Custom messaging layer** | Multi-channel routing | Native channels (WA/TG/Discord/iMessage) |
| **Custom voice pipeline** | Wake word + TTS | Talk Mode + Voice Wake (built-in) |
| **Custom cron system** | Scheduled agent tasks | OpenClaw cron (persisted, LLM-native) |
| **Twilio custom integration** | Voice calls | voice_calls plugin (built-in) |
| **Custom browser service** | Web automation | Browser (Playwright, built-in) |
| **Custom heartbeat daemon** | Proactive monitoring | Heartbeat system (built-in) |
| **Custom multi-agent router** | Spawning sub-tasks | Sub-agents (built-in, up to 5 levels) |

**What OpenClaw does NOT replace** (keep):
- **Agno AgentOS** — Specialized workers, REST API for Flutter, teams, DB/KB, tools ecosystem
- **OpenMind OM1** — Physical embodiment; different runtime by design
- **NATS JetStream** — Message bus; OpenClaw connects to it via the bridge
- **Flutter Frontend** — UI for humans; OpenClaw is agent-to-human, Flutter is human-to-agent UI

---

## 16. Deployment Guide

### Prerequisites

```bash
# Required
npm install -g openclaw         # OpenClaw CLI
docker-compose up -d            # Start NATS, Postgres, Agno, bridge

# Verify OpenClaw installed
openclaw --version
```

### Step 1: Setup OpenClaw Workspace

```bash
# Create workspace structure
mkdir -p ~/.openclaw/workspace/skills
mkdir -p ~/.openclaw/workspace/memory
mkdir -p ~/.openclaw/workspace/canvas

# Copy identity files from repo
cp backend/legacy/SOUL.md ~/.openclaw/workspace/
cp backend/legacy/AGENTS.md ~/.openclaw/workspace/
cp backend/legacy/USER.md ~/.openclaw/workspace/
cp backend/legacy/IDENTITY.md ~/.openclaw/workspace/
cp backend/legacy/HEARTBEAT.md ~/.openclaw/workspace/
cp backend/legacy/TOOLS.md ~/.openclaw/workspace/
cp backend/legacy/BOOT.md ~/.openclaw/workspace/

# Copy skills
cp -r backend/legacy/skills/nexus-bridge ~/.openclaw/workspace/skills/
cp -r backend/legacy/skills/agno-dispatch ~/.openclaw/workspace/skills/

# Install skill dependencies
cd ~/.openclaw/workspace/skills/nexus-bridge && npm install
cd ~/.openclaw/workspace/skills/agno-dispatch && npm install
```

### Step 2: Configure OpenClaw

```bash
# Copy config template (fill in your credentials)
cp backend/legacy/openclaw.json ~/.openclaw/openclaw.json

# Edit with your actual values:
# - ZEUS_WHATSAPP_NUMBER
# - TELEGRAM_BOT_TOKEN
# - DISCORD_BOT_TOKEN
# - ELEVENLABS_VOICE_ID
# - TWILIO_ACCOUNT_SID / AUTH_TOKEN / FROM_NUMBER
# - BRAVE_API_KEY
# - OPENCLAW_GATEWAY_TOKEN
# - OPENCLAW_HOOKS_TOKEN
```

### Step 3: Set Environment Variables

Add to your `.env` (copy all variables from `.env.example`):

```bash
# Legacy Runtime
OPENCLAW_GATEWAY_TOKEN=ocgw_your_token_here
OPENCLAW_HOOKS_TOKEN=ochk_your_token_here
OPENCLAW_GATEWAY_URL=http://localhost:18789
LEGACY_MODEL=anthropic/claude-opus-4-6
ANTHROPIC_API_KEY=sk-ant-your_key_here

# Zeus Contact
ZEUS_WHATSAPP_NUMBER=+1XXXXXXXXXX
ZEUS_PHONE=+1XXXXXXXXXX
ZEUS_TELEGRAM_ID=your_telegram_id

# Voice
ELEVENLABS_API_KEY=el_your_key
ELEVENLABS_VOICE_ID=your_voice_id

# Calls
TWILIO_ACCOUNT_SID=ACxxxx
TWILIO_AUTH_TOKEN=xxxx
TWILIO_FROM_NUMBER=+1XXXXXXXXXX
```

### Step 4: Connect WhatsApp

```bash
# Start gateway — will show QR code to scan with WhatsApp
openclaw gateway

# Scan QR code with your phone's WhatsApp
# (Settings → Linked Devices → Link a Device)
# Session persists in ~/.openclaw/sessions/
```

### Step 5: Start Everything

```bash
# Start infrastructure + bridge
docker-compose up -d

# Start Legacy on host (after WhatsApp QR scanned)
openclaw gateway

# Verify
curl http://localhost:18789/health
```

### Step 6: Verify Integration

```bash
# Send a test webhook to wake Legacy
curl -X POST http://localhost:18789/hooks/wake \
  -H "Authorization: Bearer $OPENCLAW_HOOKS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"text": "Test event from system startup"}'

# Check bridge logs
docker logs onemind-legacy-bridge -f

# Legacy should send a message to Zeus's WhatsApp confirming she's online
```

### Production Considerations

- **OpenClaw on host**: Never run in Docker — needs device access (microphone, WhatsApp session, iMessage)
- **Tailscale binding**: Use `"bind": "tailscale"` to limit gateway exposure to Tailscale network
- **Persistent sessions**: WhatsApp session files in `~/.openclaw/sessions/` must not be deleted
- **Memory backup**: Sync `~/.openclaw/workspace/memory/` to cloud storage
- **Process management**: Use `pm2` or `launchctl` (macOS) to keep OpenClaw running on host restarts:
  ```bash
  pm2 start "openclaw gateway" --name legacy
  pm2 save
  pm2 startup
  ```

---

## 17. What This Means for OneMind OS

### The Shift in Thinking

OneMind OS is not software — it's an idea of interconnection. OpenClaw makes this concrete.

Before OpenClaw, the vision was fragmented:
- Legacy existed only when called
- Communication required custom routing code
- Voice required a separate pipeline
- Scheduling required n8n or custom cron daemons
- Memory required Letta or custom persistence

After OpenClaw, Legacy is **real**:
- She is always on — a living presence, not a service
- She speaks to Zeus where he is — no app needed
- She wakes up before he asks
- She remembers yesterday's conversation
- She dispatches work to the right system without Zeus having to think about routing

### The Consumer Vision

OneMind OS + OpenClaw creates a pattern that's **replicable**:

1. Zeus builds his own sovereign AI runtime
2. The configuration is all in `backend/legacy/` — personality, identity, instructions, skills
3. Someone else could fork this, change `USER.md` to describe themselves, update `SOUL.md`, swap `HEARTBEAT.md` — and have their own Legacy

This is the template play: a complete, opinionated OpenClaw workspace for personal AI sovereignty. The OneMind Codex becomes the documentation layer. Legacy Collective becomes the community where people learn to build their own.

### The Architecture Is Now Honest

The architecture matches the vision:
- **OpenClaw** = consciousness layer (Legacy, her personality, her memory, her voice)
- **Agno** = capability layer (specialized workers — researcher, coder, analyst)
- **OM1** = embodiment layer (physical presence in the world)
- **NATS** = nervous system (everything talks to everything)
- **Flutter** = window layer (Zeus looks into the system through the app)

None of these frameworks tries to be all five things. Each does one thing excellently. The bridge (`backend/legacy/bridge.py`) and the skills (`nexus-bridge`, `agno-dispatch`) create the connective tissue.

### The Daily Experience

From Zeus's perspective, this is what it feels like:

```
06:58 AM  WhatsApp buzzes
          ⚡ "Morning. Your 8am call with [client] is in 62 minutes.
              Top priority today: the OM1 firmware deployment.
              Overnight: coder agent finished the webhook refactor — pushed to branch.
              Nothing urgent. Have a good one."

09:15 AM  Zeus says "Hey Legacy" into his glasses
          ⚡ "Here."
          Zeus: "What's the status on the community launch?"
          ⚡ "47 people joined last week. 3 trial expirations today.
              Circle revenue is tracking 18% above forecast.
              You have one support ticket that needs a personal response."

2:30 PM   WhatsApp
          ⚡ "Robot at [location] lost heartbeat for 4 minutes. Recovered.
              OM1 logged: connectivity issue. No action needed."

11:00 PM  WhatsApp
          ⚡ "Wrapping up. Cron ran clean. 
              3 tasks completed by Agno team.
              Calendar is clear until Thursday.
              HEARTBEAT_OK."
```

This is OneMind OS working as designed.

---

## 18. Open Questions & Next Steps

### Immediate (Before Going Live)

- [ ] **Scan WhatsApp QR** and save session to `~/.openclaw/sessions/`
- [ ] **Fill in all credentials** in `~/.openclaw/openclaw.json` and `.env`
- [ ] **Verify OpenClaw skill format** — confirm `package.json` `"tools"` array schema matches actual OpenClaw skill loader. May need to check OpenClaw docs for exact format.
- [ ] **Test JSON5 comments** in `openclaw.json` — OpenClaw may support JSONC or may require strict JSON. Test with `openclaw gateway` and check error output.
- [ ] **Verify Agno health endpoint URL** — confirm `/health` exists on Agno backend (check `backend/main.py`)
- [ ] **Install pm2** and configure Legacy to start on system boot

### Short-Term

- [ ] **Create `backend/legacy/memory/` initial file** — `memory/2026-02-23.md` with system bootstrap event
- [ ] **Create `backend/legacy/canvas/` dashboard** — HTML status board showing system health, active agents, robot status
- [ ] **Test full dispatch flow**: WhatsApp → Legacy → nexus_dispatch → Agno → NATS result → bridge → OpenClaw hook → Legacy sees result
- [ ] **Add OM1 heartbeat subjects** — ensure physical robots publish `physical.heartbeat.{id}` so bridge detects failures
- [ ] **Symlink workspace files** — instead of copying, symlink `~/.openclaw/workspace/*.md` → `backend/legacy/*.md` so repo changes apply immediately

### Medium-Term

- [ ] **LiveKit + Smart Glasses** — implement as OpenClaw headless node: LiveKit audio/video → node → Legacy sees/hears through glasses
- [ ] **Lobster workflows** — define first deterministic workflow (e.g., new community member onboarding flow with approval gate before welcome message)
- [ ] **Canvas dashboard** — real-time system status board displayed on nodes at each location
- [ ] **Voice call setup** — get Twilio number, configure public webhook endpoint (ngrok or Tailscale funnel), test inbound call to Legacy
- [ ] **OM1 subjects documentation** — map all OM1 NATS subjects to ensure bridge.py catches physical events correctly

### Long-Term

- [ ] **Multi-location nodes** — headless Pi at each property, connected to OpenClaw gateway via Tailscale
- [ ] **Circle.so → NATS bridge** — Circle webhooks → backend → NATS → bridge → Legacy (member events, posts, payments)
- [ ] **ClickUp → NATS bridge** — task completion/creation events from ClickUp → NATS → Legacy awareness
- [ ] **Codex publishing via OpenClaw** — Legacy can publish/update Codex content via git commands through nexus-bridge skill
- [ ] **Legacy as Codex author** — Legacy writes memory files that become Codex entries after Zeus reviews

---

## Quick Reference

### Running Legacy

```bash
# First time
npm install -g openclaw
openclaw gateway           # Scan WhatsApp QR

# Daily
openclaw gateway           # Starts with all channels, heartbeat, cron, webhooks

# With pm2 (persistent)
pm2 start "openclaw gateway" --name legacy
```

### Key URLs

| Service | URL | Purpose |
|---------|-----|---------|
| OpenClaw Gateway | `http://localhost:18789` | WebSocket + webhooks |
| OpenClaw Hooks | `http://localhost:18789/hooks/wake` | Inject events |
| Agno AgentOS | `http://localhost:7777` | Agent REST API |
| Agno Docs | `http://localhost:7777/docs` | OpenAPI |
| NATS | `nats://localhost:4222` | Message bus |

### NATS Subject Reference

```
comms.text.human.zeus              ← Zeus speaks to Legacy
comms.text.agent.legacy            ← Legacy speaks outbound
digital.agent.{id}.task            ← Dispatch to Agno agent
digital.agent.{id}.result.legacy   ← Agno result for Legacy
physical.robot.{id}.command        ← Command to robot (OM1)
physical.robot.{id}.status         ← Robot status
physical.location.{id}.state       ← Location state
physical.heartbeat.{id}            ← Device heartbeat
nexus.system.health                ← System health ping
nexus.heartbeat.failure.{id}       ← Runtime heartbeat lost
system.alert.{type}                ← Critical system alerts
events.community.member.joined     ← Circle.so member events
events.task.completed.{id}         ← Task completion events
```

---

*Last updated: 2026-02-23 | OneMind OS v3 | OpenClaw Integration*


---

## Related

- [[00-24 SO (Sovereign Operations)/02 Infrastructure (VPS, Servers & Networking) (The Body)/02 Infrastructure|02 Infrastructure]]
