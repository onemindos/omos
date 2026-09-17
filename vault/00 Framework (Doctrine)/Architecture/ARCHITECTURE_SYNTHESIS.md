# 🧠 OneMind OS — Architecture Synthesis & Route Forward

> **Status**: Research only — no code changes. This document answers every architecture question and lays out exactly which route to go and why.

**Date**: 2026-02-23 | Based on: full codebase review + Agno/OM1/LiveKit/OpenClaw docs

---

## Table of Contents

1. [TL;DR — The Answers Upfront](#1-tldr--the-answers-upfront)
2. [The Core Question: One Soul or Two?](#2-the-core-question-one-soul-or-two)
3. [Framework Deep Dive: What Each Does Best](#3-framework-deep-dive-what-each-does-best)
4. [Agno + OM1: How They Actually Connect](#4-agno--om1-how-they-actually-connect)
5. [Can OM1 Join an Agno Team?](#5-can-om1-join-an-agno-team)
6. [LiveKit: The Missing Real-Time Layer](#6-livekit-the-missing-real-time-layer)
7. [The LLM Question: Which Model, When](#7-the-llm-question-which-model-when)
8. [The Full Stack: Everything Together](#8-the-full-stack-everything-together)
9. [Communication Hierarchy](#9-communication-hierarchy)
10. [The Proactive + Reactive Balance](#10-the-proactive--reactive-balance)
11. [FarmBot Integration Path](#11-farmbot-integration-path)
12. [Meshtastic + Offline-First Design](#12-meshtastic--offline-first-design)
13. [Multi-Location Community Architecture](#13-multi-location-community-architecture)
14. [What the Codebase Already Has (Confirmed)](#14-what-the-codebase-already-has-confirmed)
15. [The Recommended Route](#15-the-recommended-route)
16. [Phase Plan: Cloud-First, Hybrid-Ready](#16-phase-plan-cloud-first-hybrid-ready)

---

## 1. TL;DR — The Answers Upfront

Before the deep dive — here are the direct answers to every question asked:

| Question | Answer |
|----------|--------|
| **One soul or two?** | **One soul. Two runtimes.** Legacy has one personality projected to both Agno (digital) and OM1 (physical) via NATS. |
| **Should Agno + OM1 be combined?** | **Keep separate runtimes, deeply integrate via NATS.** They're both Python — the integration is already half-built in your codebase. |
| **Can a robot be in an Agno team?** | **Yes.** Add a "physical_agent" to Agno teams that routes to OM1 via NATS. Robot becomes a team member. |
| **Can I chat with a robot?** | **Yes.** All channels route to the same session. When you text Legacy, the robot hears it too via NATS personality sync. |
| **Should Legacy use Agno as its brain?** | **No — Legacy lives in OpenClaw.** OpenClaw is her runtime. Agno agents are her workers she dispatches to. |
| **Do we need LiveKit?** | **Yes, but specifically for:** Meta Ray-Bans, remote robot video, real-time WebRTC sessions. OpenClaw Talk Mode covers local voice. |
| **Keep or remove Agno?** | **Keep Agno.** Its knowledge base (RAG), memory, MCP tools, team coordination, and 50+ APIs are irreplaceable for digital intelligence. |
| **FarmBot integration?** | **Agno MCP tool + OM1 action plugin.** Legacy commands FarmBot via Agno tools digitally; OM1 can do the same locally on-site. |
| **Meshtastic?** | **NATS leaf node at each Meshtastic gateway.** Offline mesh → local OM1 → gateway → NATS when connected. True offline-first. |
| **Multi-location?** | **NATS JetStream leaf nodes.** One NATS cluster, location-scoped subjects. Legacy sees everything from one interface. |
| **Which LLM?** | **Model factory (Agno has this).** Different agents use different models. OM1 uses local Qwen3-30B for physical reflexes. |
| **Autonomous physical or not?** | **Both autonomous AND reactive.** OM1 runs continuously at 1Hz. Agno is reactive + proactive via heartbeat. Together = full coverage. |

---

## 2. The Core Question: One Soul or Two?

### The Question

> Is it one personality/soul that has digital and physical capabilities, or should they be different souls?

### The Answer: **One Soul. Two Runtimes. One Mind.**

This is the defining principle of OneMind OS and it's already in your architecture docs:

> *"One Personality — Legacy speaks with one voice across all channels and modalities. The personality is defined once in Agno and projected to OM1."*

Here's why one soul works and how it maps technically:

```
                    ┌────────────────────────────────────┐
                    │           THE SOUL                  │
                    │         (One Definition)            │
                    │                                     │
                    │  Name: Legacy                       │
                    │  Tone: Calm authority               │
                    │  Values: Zeus's 4 domains           │
                    │  Memory: Unified across sessions    │
                    │  Boundaries: Same everywhere        │
                    └──────────────┬─────────────────────┘
                                   │ personality projection
                    ┌──────────────▼─────────────────────┐
                    │           NATS JetStream            │
                    │   subject: onemind.personality      │
                    └──────┬─────────────────────┬────────┘
                           │                     │
            ┌──────────────▼──────┐   ┌──────────▼──────────────┐
            │    AGNO RUNTIME     │   │      OM1 RUNTIME         │
            │  (Digital mind)     │   │   (Physical mind)        │
            │                     │   │                          │
            │  Same personality   │   │  Same personality        │
            │  in system prompt   │   │  in system_prompt_base   │
            │                     │   │                          │
            │  "I am Legacy,      │   │  "I am Legacy,           │
            │  your AI chief      │   │  your AI companion       │
            │  of staff..."       │   │  with a physical body.   │
            │                     │   │  I see the world around  │
            │  Digital context:   │   │  me through my sensors   │
            │  email, calendar,   │   │  and can take physical   │
            │  tasks, research    │   │  action in the world."   │
            └─────────────────────┘   └──────────────────────────┘
```

### Why NOT two souls?

Two souls would create confusion:
- Zeus asks Legacy something in WhatsApp — which Legacy responds?
- Robot Legacy and digital Legacy would have different memories = inconsistency
- Every conversation would need routing logic
- The robot wouldn't know what was discussed digitally

### Personality Projection (How It Works)

Agno is the **source of truth** for the personality. Every time something significant happens (context shift, new directive from Zeus, personality update), Agno publishes to NATS:

```
NATS subject: onemind.personality

Payload:
{
  "soul": "You are Legacy, a unified AI...",
  "current_directive": "Focus on the farm build this week",
  "zeus_context": "Zeus is at the Austin property",
  "memory_summary": "Recent: discussed irrigation schedule...",
  "timestamp": "2026-02-23T11:32:18Z"
}
```

OM1's `AgnoNATSInput` plugin receives this and updates its `system_prompt_base` — the robot immediately knows what Agno knows.

---

## 3. Framework Deep Dive: What Each Does Best

### Agno AgentOS — The Digital Mind

**What it IS**: A production-grade Python framework for cognitive agents. It's not a chatbot wrapper — it's a full agent runtime with its own database, API, memory system, and knowledge infrastructure.

**Irreplaceable capabilities**:

| Capability | Why Agno Does It Best | Alternative? |
|-----------|----------------------|-------------|
| **Knowledge Base (RAG)** | PgVector + hybrid search (vector + BM25), agentic retrieval | Would need to build from scratch |
| **3-Layer Memory** | Session → persistent → user memories, all backed by PostgreSQL | Letta (removed), custom (expensive) |
| **MCP Tools** | Native MCP support — connect ANY MCP server (Gmail, Calendar, Notion, etc.) | Only framework with this natively |
| **Teams** | Coordinate, Route, Broadcast, Tasks modes — real multi-agent coordination | OpenClaw sub-agents (simpler) |
| **Workflows** | Deterministic sequential/parallel/conditional pipelines with audit trails | Lobster DSL in OpenClaw (simpler) |
| **50+ REST APIs** | AgentOS auto-generates all endpoints for Flutter frontend | Would need custom FastAPI routes |
| **Approval Flows** | Human-in-the-loop, paused runs, approval enforcement | Nothing else has this natively |
| **Model Factory** | Per-agent model selection — GPT-4.1, Claude, Gemini, local Ollama | Yes, Agno has this |
| **Observability** | Full trace storage in your DB, no third-party egress | LangSmith (cloud, costs money) |

**What it is NOT**:
- Not always-on (reactive, not proactive by default — heartbeat solves this)
- Not a communication layer (no WhatsApp, voice, Telegram natively)
- Not a real-time WebRTC system (use LiveKit for that)
- Not a physical runtime (OM1 handles physical embodiment)

### OpenMind OM1 — The Physical Mind

**What it IS**: A modular Python runtime for physical AI agents. It runs a continuous perception-action loop (1Hz default), converts ALL sensor data to natural language via NLDB, and executes structured action commands through a Hardware Abstraction Layer.

**Irreplaceable capabilities**:

| Capability | Why OM1 Does It Best | Alternative? |
|-----------|---------------------|-------------|
| **NLDB** | ALL sensors become natural language — camera, LiDAR, mic, IMU. Hardware-agnostic. | Custom (months of work) |
| **HAL** | Abstracts robot hardware — same code for Unitree G1, Go2, TurtleBot, ROS2, Zenoh | ROS2 directly (much harder) |
| **1Hz heartbeat loop** | Never stops. Built into the core runtime. | Custom (error-prone) |
| **Dual LLM inference** | Local + cloud run in parallel, fastest good response wins | Custom (complex) |
| **Plugin system** | inputs + actions as Python plugins — trivially extensible | Custom (months) |
| **VLM integration** | Camera → text description via Vision Language Models, built-in | Custom |
| **Multi-robot support** | Unitree G1, Go2, TurtleBot 4, UbTech Yanshee, LeRobot via Zenoh | Custom per robot |

**Key insight**: OM1 ALREADY has a digital cortex layer. The existing codebase's `OpenMindLoop` (`backend/openmind/bridge.py`) IS an OM1-pattern digital agent — it runs the Sense → Fuse → Think (Agno agent) → Act pipeline. This is real OM1 architecture running RIGHT NOW on digital sensors.

**What it is NOT**:
- Not a knowledge/memory system (Agno handles this)
- Not a communication layer (OpenClaw handles this)
- Not a team orchestrator (Agno handles this)
- Not for digital tasks (email, calendar, APIs — Agno handles this)

### OpenClaw — The Consciousness Layer

**What it IS**: The always-on communication and consciousness runtime. Legacy's sovereign home. It owns every messaging surface, has native voice, heartbeat, cron, webhooks, skills, and node peripherals.

**What it is NOT**: A replacement for Agno's cognitive depth or OM1's physical capabilities.

### LiveKit — The Real-Time Transport

**What it IS**: WebRTC infrastructure for real-time audio/video/data. NOT a replacement for any of the above — it's the transport layer they all need for real-time multimodal sessions.

**2026 key features** (from LiveKit docs):
- **Gemini Live Vision** — real-time video + Gemini Realtime model, proactivity enabled
- **Meta Ray-Bans via Wearables API** — Gemini Live through glasses
- **Robot brain in the cloud** — robot's perception streams via WebRTC → cloud LLM → actions back
- **Multi-agent handoff** — pass sessions between agents seamlessly
- **Self-hosted** — run your own LiveKit server (open source)

---

## 4. Agno + OM1: How They Actually Connect

Both are Python. The connection is NATS. This is already half-built in your codebase.

### The Integration Architecture

```
AGNO (Python)                          OM1 (Python)
┌──────────────────────┐               ┌──────────────────────────┐
│  AgentOS backend     │               │  OM1 Runtime             │
│  FastAPI :7777       │               │  cortex_rate_hz: 1.0     │
│                      │               │                          │
│  Subscribe to:       │               │  Custom Input Plugins:   │
│  onemind.vision.>    │               │  └─ AgnoNATSInput        │
│  onemind.audio.>     │◄──────────────┤     subscribes to:       │
│  onemind.status.>    │    NATS       │     onemind.commands.>   │
│  onemind.alerts.>    │               │     onemind.personality  │
│                      │               │                          │
│  Publish to:         │               │  Custom Action Plugins:  │
│  onemind.commands.>  ├──────────────►│  └─ NATSPublisherAction  │
│  onemind.personality │               │     publishes to:        │
│                      │               │     onemind.vision.>     │
│  Team coordination:  │               │     onemind.status.>     │
│  physical_agent      │               │     onemind.alerts.>     │
│  → routes via NATS   │               │                          │
└──────────────────────┘               │  LLM #2 (local Qwen)    │
                                       │  decides physical actions│
                                       │  HAL → ROS2/Zenoh        │
                                       └──────────────────────────┘
```

### The Bi-Directional Flow

**Direction 1: Agno → OM1 (Command)**
```
Zeus says: "Go check the front door"
    → OpenClaw (WhatsApp) receives
    → bridge.py forwards to OpenClaw /hooks/wake
    → Legacy in OpenClaw reasons
    → calls nexus_dispatch → Agno HTTP
    → Agno team assigns to "physical_agent"
    → physical_agent publishes to NATS:
      onemind.commands.robot1 = {"action": "navigate", "destination": "front_door"}
    → OM1 subscribes (AgnoNATSInput plugin)
    → OM1 fuser receives directive as text input
    → LLM #2 decides: navigate sequence
    → HAL executes motor commands
```

**Direction 2: OM1 → Agno (Sensor data / Alert)**
```
Robot camera sees unknown person
    → VLM captions: "Unknown adult male at front door"
    → OM1 NATSPublisherAction publishes:
      onemind.alerts.robot1 = {"alert": "Unknown person at front door", "priority": "high"}
    → bridge.py subscribes to physical.robot.*.alert
    → Forwards to OpenClaw /hooks/agent with deliver: true
    → Legacy on OpenClaw reasons, notifies Zeus on WhatsApp
    → Legacy also dispatches to Agno researcher to run face recognition
    → Agno processes → result back to NATS → bridge → Legacy → Zeus
```

### Shared Session Context

The bridge between Agno's session memory and OM1's physical context:

```
onemind.personality → synced every 5 min (Agno → OM1)
    Contains: soul, current directive, recent memory summary

onemind.fleet.worldstate → aggregated view of all robots
    Contains: all robot statuses, positions, battery levels, current tasks

Agno subscribes to fleet worldstate → injects into heartbeat prompt
    → Legacy knows where every robot is at every heartbeat
```

---

## 5. Can OM1 Join an Agno Team?

**Yes.** This is the most powerful integration point.

### Physical Agent as Team Member

Add a `physical_agent` to Agno's team presets that doesn't run locally — it dispatches to OM1 via NATS:

```python
# How it works conceptually (in team presets):
"physical-team": {
    "name": "Physical + Digital Team",
    "members": ["researcher", "analyst", "physical_agent"],
    "mode": "coordinate",
}

# physical_agent is a custom Agno Agent with a NATS dispatch tool:
# When the team leader assigns it a task:
#   "physical_agent: go to the kitchen and check the FarmBot water level"
# physical_agent's tool publishes to NATS → OM1 executes → result back
```

### What This Unlocks

You can have ONE conversation in Flutter/WhatsApp/OpenClaw that orchestrates:

```
Zeus: "Check on the garden, send me a report, and book a call with the vendor"
         │
         ▼
  Agno Team Leader (coordinator)
         │
    ┌────┼──────────────┐
    │    │              │
    ▼    ▼              ▼
researcher  physical_agent  assistant
(find vendor)  (go to garden,    (book call
               look at plants,   in calendar)
               take photo)
    │    │              │
    └────┴──────────────┘
         │
    Synthesize results
         │
    "Garden looks good. Vendor meeting booked for Thursday.
     FarmBot water level at 78%."
```

This is ONE conversation, ONE session, ONE Legacy — doing physical and digital simultaneously.

### Chat With a Robot Directly

Via any channel (WhatsApp, Flutter chat screen, Telegram):
```
Zeus: "Hey robot, what do you see right now?"
Legacy: routes to physical_agent (or directly via nexus_dispatch to OM1)
OM1 vision: "Living room — couch, TV, coffee on table, window open"
Legacy: "The living room looks clear. Coffee's still on the table."
```

The robot is just another agent you chat with through Legacy's unified session.

---

## 6. LiveKit: The Missing Real-Time Layer

### The Original Architecture Said "No LiveKit"

The existing architecture docs say: *"No LiveKit — The Interaction Gateway + Agno's native session system replaces LiveKit rooms entirely. ~200 lines of FastAPI vs. a full media server cluster."*

**This was correct for a single-user, local-first setup.** But the scope has expanded dramatically:

- Multi-location communities around the world
- Meta Ray-Bans as an input device
- Remote robot video streaming
- Real-time voice sessions from mobile (not just local Talk Mode)
- Family members and community managers needing voice access

**LiveKit is now essential for the following specific use cases:**

### Use Case 1: Meta Ray-Bans + Gemini Live Vision

```
Meta Ray-Bans
  → LiveKit Wearables API (Meta approved)
  → LiveKit Room (audio + video stream)
  → Gemini Live Vision agent (sees what Zeus sees)
  → Proactivity enabled: agent speaks when relevant
  → ElevenLabs TTS → back to glasses speaker

Zeus is walking through the farm.
Glasses are streaming. Legacy sees everything.
"That plant in row 4 looks dry — FarmBot irrigation is off in that zone."
```

This is LIVE right now. LiveKit has a demo of this working with Meta Ray-Bans.

### Use Case 2: Remote Robot Video Monitoring

```
Robot camera → OM1 processes locally
             → streams H.264 via LiveKit Room
             → Community managers can view robot POV remotely
             → Voice command back to robot via LiveKit
             → OM1 receives audio → ASR → processes command
```

### Use Case 3: Real-Time Voice Sessions (Mobile)

```
Zeus on iPhone in the field
  → LiveKit mobile SDK
  → Voice session with Legacy
  → STT → Agno processes → TTS → back to phone
  → Full conversation, full context, no latency

Benefit over OpenClaw Talk Mode:
  - Works over cellular, not just local WiFi
  - Multiple people can join the same room
  - Video sharing enabled
  - Works from any location
```

### Use Case 4: Community Members Accessing Legacy

```
Legacy Collective member needs help
  → WebRTC session via Flutter/web
  → LiveKit room joins
  → Community-scoped Legacy instance responds
  → (Not Zeus's personal Legacy — a community variant)
```

### LiveKit Integration Architecture

```
LiveKit Server (self-hosted on your infrastructure)
    │
    ├── Room: zeus-personal
    │     ├── Legacy agent (Gemini Realtime or custom STT→Agno→TTS)
    │     ├── Zeus's glasses (Meta Ray-Bans via Wearables API)
    │     └── Zeus's phone (LiveKit mobile SDK)
    │
    ├── Room: robot-robot1-feed
    │     ├── OM1 robot (publishes video/audio track)
    │     ├── Legacy agent (monitors, can speak to robot)
    │     └── Zeus's dashboard (Flutter app, view-only)
    │
    └── Room: community-support
          ├── Community Legacy agent (different soul, public)
          └── Community member (Flutter web app)
```

### LiveKit Does NOT Replace OpenClaw

| Capability | OpenClaw | LiveKit |
|-----------|----------|---------|
| WhatsApp | ✅ Native | ❌ No |
| iMessage | ✅ Native | ❌ No |
| Telegram | ✅ Native | ❌ No |
| Discord | ✅ Native | ❌ No |
| Local wake word | ✅ Native | ❌ No |
| Local Talk Mode | ✅ ElevenLabs | ⚠️ Needs custom |
| Heartbeat/Cron | ✅ Native | ❌ No |
| Webhooks | ✅ Native | ❌ No |
| Lobster workflows | ✅ Native | ❌ No |
| Node peripherals | ✅ Native | ❌ No |
| Meta Ray-Bans | ❌ No | ✅ Wearables API |
| Mobile WebRTC | ❌ No | ✅ Native |
| Remote robot video | ❌ No | ✅ Native |
| Real-time multi-user room | ❌ No | ✅ Native |
| Gemini Live Vision | ❌ No | ✅ Native |

**They're complementary.** OpenClaw owns async channels. LiveKit owns real-time WebRTC sessions.

---

## 7. The LLM Question: Which Model, When

### The Three-Tier LLM Strategy

Agno has a model factory. You can — and should — use different models per agent. Here's the full strategy:

```
┌─────────────────────────────────────────────────────────────────┐
│                     LLM ASSIGNMENT MAP                           │
├──────────────────┬──────────────────────┬───────────────────────┤
│  Agent           │  Model               │  Why                  │
├──────────────────┼──────────────────────┼───────────────────────┤
│ researcher       │ GPT-4.1 / Claude Opus│ Deep reasoning, cost  │
│                  │                      │ justified by quality  │
├──────────────────┼──────────────────────┼───────────────────────┤
│ analyst          │ Claude Sonnet        │ Strong analysis,      │
│                  │                      │ faster than Opus      │
├──────────────────┼──────────────────────┼───────────────────────┤
│ coder            │ Claude Sonnet        │ Best code model       │
│                  │                      │ currently             │
├──────────────────┼──────────────────────┼───────────────────────┤
│ assistant        │ GPT-4.1-mini         │ Quick, cheap, good    │
│                  │                      │ enough for general    │
├──────────────────┼──────────────────────┼───────────────────────┤
│ web_automator    │ GPT-4o               │ Vision for browser    │
│                  │                      │ screenshots           │
├──────────────────┼──────────────────────┼───────────────────────┤
│ heartbeat (Agno) │ GPT-4.1-mini         │ Cheap, runs 48x/day   │
│                  │                      │ Only escalates if     │
│                  │                      │ something found       │
├──────────────────┼──────────────────────┼───────────────────────┤
│ OM1 cortex       │ Qwen3-30B (local)    │ 200ms, physical       │
│ (LLM #2)         │ GPT-4.1-nano fallback│ reflexes, runs 3600x  │
│                  │                      │ per hour              │
├──────────────────┼──────────────────────┼───────────────────────┤
│ Legacy (OpenClaw)│ Claude Opus 4.6      │ The consciousness —   │
│                  │                      │ deserves the best     │
├──────────────────┼──────────────────────┼───────────────────────┤
│ LiveKit agent    │ Gemini 2.5 Pro       │ Best for real-time    │
│ (voice sessions) │ (via Gemini Realtime)│ voice + vision        │
└──────────────────┴──────────────────────┴───────────────────────┘
```

### Agno Model Factory

This is already in your codebase (`backend/models.py`). Each agent in `presets.py` specifies `model_id`. You can:

```python
# Per-agent model selection already works:
"researcher": {
    "model_id": "anthropic/claude-opus-4-6",  # Best for deep research
}
"assistant": {
    "model_id": "openai/gpt-4.1-mini",  # Cheap for quick tasks
}
```

No new infrastructure needed — just update the `model_id` values.

### The LLM Decision Tree

```
Is this a real-time physical decision? (< 500ms needed)
  YES → LLM #2 (OM1 local Qwen3-30B)
  NO → continue...

Is this a voice/vision session via WebRTC?
  YES → Gemini Realtime (via LiveKit)
  NO → continue...

Is this a consciousness/communication decision?
  YES → Claude Opus 4.6 (Legacy in OpenClaw)
  NO → continue...

Is this deep research or complex reasoning?
  YES → Claude Opus or GPT-4.1 (Agno researcher/analyst)
  NO → continue...

Is this a quick task, heartbeat, or routine check?
  YES → GPT-4.1-mini (cheap, fast, sufficient)
```

---

## 8. The Full Stack: Everything Together

This is the complete OneMind OS architecture as it should be:

```
╔════════════════════════════════════════════════════════════════════════════╗
║                    HOW THE WORLD REACHES ONEMIND OS                        ║
╠════════════════════════════════════════════════════════════════════════════╣
║                                                                            ║
║  ASYNC CHANNELS (OpenClaw)        REAL-TIME SESSIONS (LiveKit)             ║
║  WhatsApp · iMessage · Telegram   Meta Ray-Bans · Mobile · Robot Video    ║
║  Discord · Slack · WebChat        Community WebRTC · Remote Monitoring     ║
║                                                                            ║
╚═════════════════╤══════════════════════╤═════════════════════════════════╝
                  │                      │
      ┌───────────▼──────────┐  ┌────────▼─────────────┐
      │      OPENCLAW         │  │       LIVEKIT         │
      │   Legacy's Runtime    │  │  Real-Time Transport  │
      │                       │  │                       │
      │  SOUL · AGENTS.md     │  │  Gemini Live Vision   │
      │  Heartbeat (30m)      │  │  WebRTC Rooms         │
      │  Cron · Webhooks      │  │  Multi-agent handoff  │
      │  Talk Mode · Voice    │  │  Wearables API        │
      │  Nodes · Browser      │  │  Robot video streams  │
      │  Lobster workflows    │  │                       │
      └───────────┬───────────┘  └──────────┬────────────┘
                  │                          │
                  │ webhooks + WS            │ routes to Agno
                  │                          │
      ┌───────────▼──────────────────────────▼────────────┐
      │              LEGACY BRIDGE SERVICE                  │
      │           backend/legacy/bridge.py                  │
      │     (NATS ↔ OpenClaw ↔ LiveKit relay layer)        │
      └───────────────────────────┬────────────────────────┘
                                  │
                     ┌────────────▼───────────┐
                     │    NATS JetStream       │
                     │   (Nervous System)      │
                     │                        │
                     │  7 streams:            │
                     │  DIGITAL · PHYSICAL    │
                     │  HUMAN · SYSTEM        │
                     │  EVENTS · NEXUS        │
                     │  COMMS                 │
                     └──────┬─────────────────┘
          ┌─────────────────┼──────────────────┬──────────────────┐
          │                 │                  │                  │
  ┌───────▼───────┐ ┌───────▼───────┐  ┌──────▼──────┐  ┌───────▼───────┐
  │  AGNO         │ │     OM1       │  │  LEROBOT    │  │  FLUTTER APP  │
  │  AgentOS :7777│ │  Physical     │  │  Skills     │  │  Frontend     │
  │               │ │  Runtime      │  │  (in OM1)   │  │  :8080        │
  │  researcher   │ │               │  │             │  │               │
  │  analyst      │ │  1Hz loop:    │  │  locomotion │  │  30+ screens  │
  │  coder        │ │  SENSE→FUSE   │  │  grasping   │  │  agent UI     │
  │  web_automator│ │  →THINK→ACT   │  │  learned    │  │  team UI      │
  │  db_manager   │ │               │  │  behaviors  │  │  robot status │
  │  physical_agnt│ │  LLM #2 local │  │             │  │  community    │
  │               │ │  Qwen3-30B    │  │             │  │  dashboards   │
  │  Knowledge DB │ │  Camera/LiDAR │  │             │  │               │
  │  Memory (3-layer) │  Mic/IMU   │  │             │  │               │
  │  MCP Tools    │ │  HAL→ROS2     │  │             │  │               │
  │  Teams        │ │               │  │             │  │               │
  │  Workflows    │ │  NATS plugins │  │             │  │               │
  └───────────────┘ └───────────────┘  └─────────────┘  └───────────────┘
          │                 │
          └────────────┬────┘
                       │
               ┌───────▼───────────────────┐
               │   EXTERNAL INTEGRATIONS    │
               │                           │
               │  FarmBot (REST/MQTT)       │
               │  Meshtastic (LoRa mesh)    │
               │  Home Assistant (local)    │
               │  Circle.so (community)     │
               │  ClickUp (tasks)           │
               │  Zoho One (business)       │
               │  Webflow (brand site)      │
               │  Obsidian Git (codex)      │
               └───────────────────────────┘
```

---

## 9. Communication Hierarchy

This is the routing map for how any communication reaches Legacy and how Legacy responds:

```
INPUT → WHO HANDLES IT → HOW IT ROUTES

WhatsApp message from Zeus
  → OpenClaw native channel
  → Legacy in OpenClaw reasons
  → If simple: responds directly via WhatsApp
  → If complex: dispatches via nexus_dispatch → Agno
  → Agno processes → result via NATS → bridge → OpenClaw → WhatsApp

"Hey Legacy" wake word at home
  → OpenClaw Voice Wake (local microphone, any node)
  → Triggers Talk Mode session
  → Same session as WhatsApp (same context)

Meta Ray-Bans: "Legacy, what's that plant?"
  → LiveKit Wearables API
  → Video frame → Gemini Live Vision
  → Gemini sees plant, identifies it
  → Response → ElevenLabs → glasses speaker
  → Context synced to main Legacy session via NATS

Robot sends vision alert (unknown person at door)
  → OM1 NATSPublisherAction
  → NATS: physical.robot.robot1.alert
  → bridge.py picks up
  → POST /hooks/agent with deliver: true
  → Legacy spawns agent turn
  → WhatsApp: "Someone at the front door. Do you recognize them? [photo]"

Legacy heartbeat fires (every 30 min)
  → OpenClaw reads HEARTBEAT.md
  → Checks: calendar, tasks, system health, robot status, community
  → If nothing: HEARTBEAT_OK (silent)
  → If something: WhatsApp message to Zeus

Flutter chat screen message
  → POST to Agno AgentOS /agents/{id}/runs
  → Agno processes with full knowledge+memory
  → SSE streaming response → Flutter renders live
  → Important events also published to NATS → Legacy aware

Cron job fires (e.g., daily farm report at 8am)
  → OpenClaw cron (isolated session)
  → Dispatches to Agno via nexus_dispatch or agno_dispatch skill
  → researcher + physical_agent both run
  → researcher: FarmBot API data analysis
  → physical_agent: OM1 robot does a walkthrough
  → Results combine → WhatsApp morning brief to Zeus
```

---

## 10. The Proactive + Reactive Balance

### The Problem with Pure Reactive Systems

Agno agents are reactive — they respond when called. If you only use Agno, Legacy goes silent the moment no one is talking to her. That's not OneMind OS.

### The Problem with Pure Proactive Systems

OM1 runs 1Hz always. If you make the digital side the same, it's expensive — GPT-4.1 at 1Hz costs ~$100+/day just on tokens.

### The Solution: Layered Autonomy

```
AUTONOMY LEVEL    FREQUENCY   MODEL       TRIGGERS
─────────────────────────────────────────────────────
OM1 Physical Loop   1 Hz      Qwen3-30B   Always (physical safety)
OpenClaw Heartbeat  30 min    Claude Opus Events, time, directives
Agno Heartbeat      15 min    GPT-4.1-mini Gathered data analysis
ClickUp Webhooks    On-event  Mini        Task updates
Circle.so Webhooks  On-event  Mini        Member events
HA Webhooks         On-event  Claude Opus Critical home events
NATS Robot Alerts   On-event  Claude Opus Physical world events
User conversation   On-demand Full model  Explicit chat
```

### How Context Flows Between Proactive and Reactive

When Legacy's heartbeat fires AND a user message arrives during the same cycle:
- They share the same session context (same `session_id`)
- The heartbeat result is injected as system context
- The user sees a response that incorporates both

This means Legacy is **always prepared** when Zeus asks a question — she's been watching even when not talking.

---

## 11. FarmBot Integration Path

FarmBot is a precision farming robot with a full REST API and MQTT broker. It fits in two places:

### Path 1: Agno MCP Tool (Digital Control)

```
Legacy → agno_dispatch skill → Agno coder agent
  → FarmBot MCP tool (or custom REST tool)
  → FarmBot REST API: /api/sequences, /api/plants, /api/sensors
  → Read sensor data, trigger sequences, move to coordinates
  → Result back via NATS → Legacy
```

**Example**: "Schedule watering for zone 3 at 6am tomorrow"
- Legacy → Agno coder → FarmBot API: `POST /api/schedules`
- Done. No robot needed.

### Path 2: OM1 Action Plugin (Physical Integration)

```
OM1 robot is physically near the FarmBot installation
  → Custom OM1 action plugin: FarmBotPlugin
  → Connects to FarmBot's local MQTT broker
  → Can: move FarmBot arm, check plant health via OM1 camera
  → OM1 sees what FarmBot is doing + nearby plants
  → Reports to NATS → Agno → Legacy
```

**Example**: "Check if tomatoes in zone 3 need water"
- physical_agent dispatch → OM1 at farm location
- OM1 camera looks at plants (VLM: "Leaves slightly droopy, possible underwatering")
- OM1 FarmBotPlugin: query soil sensor via MQTT (returns 22% moisture)
- OM1 publishes to NATS: "Zone 3 soil moisture at 22%, recommend irrigation"
- Legacy: "Zone 3 is a bit dry. Triggering FarmBot irrigation sequence."
- Agno FarmBot tool: POST irrigation sequence

**FarmBot subjects on NATS**:
```
physical.farmbot.{id}.sensor    → Soil moisture, temperature
physical.farmbot.{id}.status    → Position, current job, sequence running
physical.farmbot.{id}.command   → Trigger sequences, move arm
events.farmbot.{id}.alert       → Low water, plant health alerts
```

---

## 12. Meshtastic + Offline-First Design

### What Meshtastic Is

Meshtastic is a LoRa mesh radio protocol — long-range, low-power, no internet required. It creates a local mesh network between devices. This is your off-grid backbone.

### How It Integrates

```
OFF-GRID SCENARIO: Farm has no internet
  │
  ├── OM1 robot on farm has Meshtastic radio module
  │     → Publishes sensor data to local mesh
  │
  ├── FarmBot has Meshtastic gateway (Pi + LoRa radio)
  │     → Receives mesh messages, stores locally
  │
  ├── Zeus has Meshtastic on phone/walkie
  │     → Can receive alerts even without cellular
  │
  └── Gateway device (has satellite/cellular occasionally)
        → When connected: syncs all queued NATS messages
        → Bridge: Meshtastic ↔ NATS ↔ Legacy
        → When disconnected: local OM1 operates autonomously
```

### The NATS Leaf Node Pattern

When connectivity exists, each location runs a NATS leaf node that connects to the main NATS cluster:

```
[Zeus's Server / Cloud] ←──── Main NATS Cluster
                                      │
                      ┌───────────────┼───────────────┐
                      │               │               │
                [Farm NATS      [City NATS      [Rental NATS
                 Leaf Node]      Leaf Node]      Leaf Node]
                      │               │               │
                [Local OM1]    [Local OM1]    [Home Assistant]
                [FarmBot]      [Security]     [Sensors]
                [Meshtastic]
```

**When offline**: OM1 operates autonomously using local LLM #2 (Qwen3-30B on-device). No cloud needed for physical reflexes.

**When online**: NATS leaf node syncs all events. Legacy gets caught up on everything that happened offline.

### Meshtastic NATS Bridge

```python
# Conceptual bridge: Meshtastic → NATS
# Runs on the gateway Pi at each location
# When internet available: publishes queued Meshtastic messages to NATS
# When offline: stores to local queue, flushes when reconnected

meshtastic_subjects:
  physical.mesh.{location}.{node_id}.sensor    → Sensor readings
  physical.mesh.{location}.alert               → Urgent alerts (low battery, intrusion)
  physical.mesh.{location}.heartbeat           → Node still alive
```

### Offline-First Principle for OM1

OM1 is already offline-capable at the physical level:
- LLM #2 (Qwen3-30B) runs 100% locally
- Camera, sensors, motors work without internet
- Robot continues patrolling, monitoring, acting
- Queues events locally → pushes to NATS when reconnected

**What degrades offline**:
- Legacy (OpenClaw) can't respond via WhatsApp
- Agno research/analysis not available
- New directives from Zeus can't arrive

**What keeps working offline**:
- OM1 physical loop at 1Hz
- FarmBot autonomous schedules (already programmed)
- Meshtastic local alerts (Zeus's phone gets text alerts)
- Home Assistant local automations

---

## 13. Multi-Location Community Architecture

### The Vision: Managing Everything From One Interface

Zeus runs multiple locations — farms, communities, properties. Each needs:
- Physical monitoring (OM1)
- Digital management (Agno)
- Communication (OpenClaw)
- Local automation (Home Assistant)
- Farm automation (FarmBot if applicable)

### Location Architecture

```
LEGACY (Single Legacy, One Interface, Sees All)
          │
          └── NATS JetStream (clustered across locations)
                    │
        ┌───────────┼───────────┬────────────────┐
        │           │           │                │
  [Austin TX]  [Property 2] [Farm Location] [Future]
        │           │           │
  ┌─────▼─────┐  ┌──▼──────┐  ┌─▼──────────┐
  │ OM1 robot │  │ OM1 node│  │ OM1 robot  │
  │ Home Auto │  │ Cameras │  │ FarmBot    │
  │ Security  │  │ HA auto │  │ Meshtastic │
  │ NATS leaf │  │ NATS lf │  │ NATS leaf  │
  └───────────┘  └─────────┘  └────────────┘

NATS subjects scoped by location:
  physical.location.austin.>
  physical.location.property2.>
  physical.location.farm.>

Legacy can ask:
  "What's the status at all locations?"
  → Queries all location subjects
  → Synthesizes fleet worldstate
  → "Austin: all clear. Property 2: motion detected at 2am. Farm: FarmBot running irrigation."
```

### Community Management (Circle.so Integration)

```
Community events → Circle.so webhook → backend webhooks.py
  → NATS: events.community.{event_type}
  → bridge.py → OpenClaw /hooks/agent
  → Legacy handles:
      - New member: welcome message, ClickUp task, Discord intro
      - Support request: route to appropriate agent
      - Payment: Zoho update, tier change, thank you message
      - Post alert: urgent community issue surfaced to Zeus

Community agents (separate Agno team):
  "community-team": {
      "members": ["researcher", "assistant"],
      "purpose": "Community management for Legacy Collective",
      "scope": "Circle.so, Discord, community support"
  }
```

### Multi-User Access (Later Phase)

When other team members need access:
```
Community manager → Flutter app (their own login)
  → Agno AgentOS JWT auth (per-user isolation)
  → Their own session, their own context
  → Specific agent scope (not Legacy's full access)
  → Zeus's dashboard shows their activity

Family member → Flutter app
  → Limited scope agents (LE domain only)
  → Full voice via LiveKit if needed
```

---

## 14. What the Codebase Already Has (Confirmed)

After reviewing every relevant file, here's what's actually built and working:

### ✅ Already Built

| Component | File | Status |
|-----------|------|--------|
| OpenClaw workspace (7 identity files) | `backend/legacy/*.md` | ✅ Complete |
| OpenClaw config | `backend/legacy/openclaw.json` | ✅ Complete |
| NATS→OpenClaw bridge | `backend/legacy/bridge.py` | ✅ Complete |
| nexus-bridge skill | `backend/legacy/skills/nexus-bridge/` | ✅ Complete |
| agno-dispatch skill | `backend/legacy/skills/agno-dispatch/` | ✅ Complete |
| Legacy removed from Agno | `backend/agents/presets.py` | ✅ Complete |
| docker-compose with bridge | `docker-compose.yml` | ✅ Complete |
| 9 Agno worker agents | `backend/agents/presets.py` | ✅ Complete |
| 3 Agno teams | `backend/agents/teams/presets.py` | ✅ Complete |
| OM1 heartbeat loop (digital) | `backend/openmind/bridge.py` | ✅ Complete |
| OM1 digital sensors | `backend/openmind/sensors/digital.py` | ✅ Complete |
| NATS 7-stream bus | `backend/bus.py` | ✅ Complete |
| Nexus slim coordinator | `backend/nexus/core.py` | ✅ Complete |
| Flutter frontend (30+ screens) | `frontend/lib/` | ✅ Complete |
| Agno AgentOS backend | `backend/main.py` | ✅ Complete |
| PostgreSQL + pgvector | `docker-compose.yml` | ✅ Complete |

### 🔲 Needs to Be Built (Confirmed Gaps)

| Component | Priority | Complexity | Route |
|-----------|---------|-----------|-------|
| LiveKit server + agent | **HIGH** | Medium | Add LiveKit Python agent |
| Meta Ray-Bans integration | HIGH | Low | LiveKit Wearables API |
| `physical_agent` Agno preset | HIGH | Low | Add to `presets.py` (routes to OM1) |
| FarmBot Agno tool | MEDIUM | Medium | REST tool in `backend/tools/` |
| FarmBot OM1 plugin | MEDIUM | Medium | OM1 action plugin (separate repo) |
| Meshtastic NATS bridge | LOW | Medium | Pi service at each location |
| NATS leaf nodes per location | LOW | Low (config) | NATS server config |
| LiveKit OM1 video streaming | MEDIUM | High | OM1 + LiveKit WebRTC |
| Community LiveKit rooms | LOW | Low | LiveKit server config |
| OM1 real robot config | MEDIUM | High | OM1 plugin for Unitree/TurtleBot |

---

## 15. The Recommended Route

### The Verdict: Keep Everything, Add LiveKit

Do NOT collapse or remove any framework. Each has an irreplaceable role. The architecture is correct. Here's the definitive assignment:

```
FRAMEWORK     ROLE                    IRREPLACEABLE BECAUSE
──────────────────────────────────────────────────────────────────
OpenClaw    Legacy's consciousness   Always-on, multi-channel,
            + communication          heartbeat, voice, native
                                     messaging. Nothing else
                                     does ALL of this.

Agno        Digital intelligence     Knowledge (RAG), 3-layer
            + worker agents          memory, MCP, teams,
                                     workflows, 50+ APIs for
                                     Flutter. Cannot be replaced.

OM1         Physical embodiment      HAL, NLDB, 1Hz loop,
                                     multi-robot support.
                                     Both are Python → native
                                     integration via NATS.

LiveKit     Real-time WebRTC         Meta Ray-Bans, mobile
            transport                voice sessions, remote
                                     robot video. The gap that
                                     OpenClaw Talk Mode can't
                                     fill for mobile/global use.

NATS        Nervous system           Connects everything.
            (already built)          Sub-ms routing. Never remove.
```

### The One Soul Answer (Final)

**ONE SOUL. Dual runtime embodiment.**

Legacy has one personality, one memory, one identity. She is instantiated in:
1. **OpenClaw** — consciousness, communication, always-on awareness
2. **Agno** — cognitive depth, digital capabilities, team coordination
3. **OM1** — physical presence, sensor awareness, motor control

All three read from the same personality definition. All three publish to the same NATS bus. All three share context via NATS personality projection. Legacy doesn't split into three — she extends into all three simultaneously.

The robot doesn't have its own personality. It has Legacy's personality expressed through a body.

---

## 16. Phase Plan: Cloud-First, Hybrid-Ready

### Phase 1 — Get It Live (Now)

Focus: OpenClaw + Agno working together as one system.

```
Week 1-2:
  ☐ Install OpenClaw: npm install -g openclaw
  ☐ Scan WhatsApp QR, connect to Legacy workspace
  ☐ Start docker-compose (NATS, Postgres, Agno, bridge)
  ☐ Test: WhatsApp → Legacy → Agno researcher → result back to WhatsApp
  ☐ Connect Telegram bot
  ☐ Set up morning brief cron
  ☐ Verify heartbeat: Legacy wakes every 30min

Week 3-4:
  ☐ Add physical_agent to Agno presets (routes via NATS to OM1 digital cortex)
  ☐ Test digital-only physical agent (OM1 bridge.py already running)
  ☐ Add FarmBot REST tool to Agno toolset
  ☐ Connect Circle.so webhook → NATS → bridge → Legacy
  ☐ Connect ClickUp webhook → NATS → bridge → Legacy
  ☐ Set up LiveKit self-hosted server
  ☐ Test basic voice session: Flutter → LiveKit → STT → Agno → TTS → Flutter
```

### Phase 2 — Real-Time + Physical (Next 60 Days)

```
  ☐ Meta Ray-Bans: LiveKit Wearables API + Gemini Live Vision agent
  ☐ OM1 with real robot (TurtleBot 4 is easiest entry)
  ☐ OM1 NATS plugins live (real robot commands + sensor data to NATS)
  ☐ physical_agent in Agno team receives real robot data
  ☐ LiveKit room for remote robot video monitoring (Flutter dashboard)
  ☐ Home Assistant → NATS webhook bridge
  ☐ Meshtastic at one location (farm or property)
```

### Phase 3 — Multi-Location + Community (90+ Days)

```
  ☐ NATS leaf nodes at each location
  ☐ Multi-robot fleet (location-scoped NATS subjects)
  ☐ FarmBot at farm location (OM1 + FarmBot plugin)
  ☐ Community LiveKit rooms (Legacy Collective voice sessions)
  ☐ Multi-user access (community managers via Flutter)
  ☐ Offline-first fully working (OM1 autonomous + Meshtastic alerts)
  ☐ LeRobot skill learning deployed into OM1
```

### Start Here Today

The single highest-ROI action: **get OpenClaw running with WhatsApp connected**.

Everything else builds on Legacy being alive. Once she's live on WhatsApp with the heartbeat running, dispatching to Agno workers, and receiving NATS events from the bridge — the architecture is real. Everything else is additive.

```bash
# Step 1: Install OpenClaw
npm install -g openclaw

# Step 2: Copy workspace
cp backend/legacy/*.md ~/.openclaw/workspace/
cp -r backend/legacy/skills ~/.openclaw/workspace/
cd ~/.openclaw/workspace/skills/nexus-bridge && npm install
cd ~/.openclaw/workspace/skills/agno-dispatch && npm install

# Step 3: Configure credentials
cp backend/legacy/openclaw.json ~/.openclaw/openclaw.json
# Edit: fill in WhatsApp number, Telegram token, ElevenLabs voice ID

# Step 4: Start infrastructure
docker-compose up -d  # NATS, Postgres, Agno, bridge

# Step 5: Start Legacy
openclaw gateway  # Scan WhatsApp QR

# Step 6: Verify
# Send "hey" to WhatsApp → Legacy responds
# Send "research the latest LiveKit features" → Legacy dispatches to researcher → result
```

---

## Appendix A: Quick Reference — Framework Decisions

| Decision | Choice | Rationale |
|----------|--------|-----------|
| Remove Agno? | **No** | Knowledge, memory, MCP, teams irreplaceable |
| Remove OM1? | **No** | HAL, NLDB, 1Hz loop irreplaceable |
| Remove OpenClaw? | **No** | Always-on consciousness, messaging, heartbeat |
| Add LiveKit? | **Yes** | Ray-Bans, mobile WebRTC, robot video |
| One soul? | **Yes** | Personality projected via NATS to all runtimes |
| OM1 in Agno team? | **Yes** | physical_agent preset routes to OM1 via NATS |
| Which LLM? | **Model factory** | Per-agent assignment, OM1 uses local Qwen3-30B |
| FarmBot? | **Both Agno tool + OM1 plugin** | Digital control + physical co-location |
| Meshtastic? | **NATS gateway bridge** | Offline-first, syncs when connected |
| Multi-location? | **NATS leaf nodes** | One cluster, location-scoped subjects |
| Offline-first? | **Hybrid** | OM1 always works, cloud optional |
| Legacy's home? | **OpenClaw** | Consciousness runtime, not inside Agno |

## Appendix B: NATS Subject Master Map

```
# OpenClaw ↔ Agno (via bridge.py)
comms.text.human.zeus              ← Zeus's messages
comms.text.agent.legacy            ← Legacy's outbound
digital.agent.{id}.task            ← Dispatch to Agno worker
digital.agent.{id}.result.legacy   ← Agno result for Legacy
system.alert.{type}                ← Critical system alerts
nexus.heartbeat.failure.{id}       ← Runtime heartbeat lost

# OM1 Integration
onemind.commands.{robot_id}        ← Agno → OM1 (commands)
onemind.vision.{robot_id}          ← OM1 → Agno (camera)
onemind.audio.{robot_id}           ← OM1 → Agno (mic)
onemind.status.{robot_id}          ← OM1 → Agno (battery, pos)
onemind.alerts.{robot_id}          ← OM1 → Agno (urgent events)
onemind.personality                ← Agno → OM1 (soul sync)
onemind.fleet.worldstate           ← Aggregated fleet view

# Location-Scoped
physical.location.{id}.state       ← Location query
physical.location.{id}.heartbeat   ← Location alive check

# FarmBot
physical.farmbot.{id}.sensor       ← Soil, temp readings
physical.farmbot.{id}.command      ← Trigger sequence
events.farmbot.{id}.alert          ← Plant/water alerts

# Meshtastic
physical.mesh.{location}.{node}.sensor  ← Mesh sensor data
physical.mesh.{location}.alert         ← Mesh urgent alerts
physical.mesh.{location}.heartbeat     ← Node alive

# Community
events.community.member.joined     ← Circle.so new member
events.community.{event_type}      ← General community events
events.task.completed.{id}         ← ClickUp task done
```

---

*Last updated: 2026-02-23 | Research document — no code changes | OneMind OS v3*
