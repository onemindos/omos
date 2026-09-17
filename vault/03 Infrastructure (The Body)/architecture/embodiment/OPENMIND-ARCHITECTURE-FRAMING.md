# One Mind OS v2 — System Architecture

> Two systems. One personality. Agno is the brain. OM1 is the body. NATS is the nervous system.

## Overview

One Mind OS is a four-layer architecture that unifies digital AI agent capabilities (Agno) with physical robot embodiment (OpenMind OM1) into a single coherent personality called **Legacy**.

```
═══════════════════════════════════════════════════════════
  LAYER 0: INTERACTION GATEWAY
  Voice, Chat, Channels — how you talk to Legacy
═══════════════════════════════════════════════════════════
  🎙️ Voice (Gemini Live API / ElevenLabs TTS)
  💬 Chat (Flutter app / Web UI via WebSocket)
  📱 Channels (Telegram, WhatsApp, Slack, iMessage)
  🤖 Robot mic (OM1 speech-to-text)
  🕶️ AR Glasses (Brilliant Labs Frame)
     │
     │  ALL inputs → text → POST to Agno AgentOS
     │  ALL outputs → TTS / text → back to channel
     ▼
═══════════════════════════════════════════════════════════
  LAYER 1: AGNO COGNITIVE LAYER (The Brain)
  Thinks, remembers, coordinates, digital tools
═══════════════════════════════════════════════════════════
  🧠 LLM #1: GPT-4.1 / Claude (powerful, 1-5 sec)
  📚 Knowledge (PgVector RAG)
  💾 Memory (PostgreSQL, 3-layer, cross-session)
  🔌 MCP Tools (Gmail, Calendar, ClickUp, etc.)
  👥 Teams: UI, GE, HP, LE + Physical Agent
  💓 Heartbeat: every 15 min (OpenClaw pattern)
  ⚡ Reactive: instant response to chat/events
     │
     │ publish commands / subscribe events
     ▼
═══════════════════════════════════════════════════════════
  LAYER 2: NATS MESSAGE BUS (Nervous System)
  Connects brain and body — sub-millisecond
═══════════════════════════════════════════════════════════
  onemind.commands.{robot_id}    → Agno to OM1
  onemind.vision.{robot_id}      → OM1 to Agno
  onemind.audio.{robot_id}       → OM1 to Agno
  onemind.status.{robot_id}      → OM1 to Agno
  onemind.personality             → Agno to OM1
  onemind.fleet.worldstate        → Aggregated view
     │
     │ subscribe commands / publish sensors
     ▼
═══════════════════════════════════════════════════════════
  LAYER 3: OM1 EMBODIMENT LAYER × N ROBOTS (The Body)
  Sees, hears, moves, reacts fast
═══════════════════════════════════════════════════════════
  ⚡ LLM #2: Qwen3-30B local (fast, 200ms)
  👁️ Camera → VLM → text descriptions
  👂 Mic → ASR → text transcripts
  📡 LiDAR / IMU → spatial awareness
  🦿 HAL → ROS2 → motors
  💚 Heartbeat: every 0.5-2 sec (never stops)
```

## Core Principles

1. **One Personality** — Legacy speaks with one voice across all channels and modalities. The personality is defined once in Agno and projected to OM1.

2. **Two Separate Programs** — Agno and OM1 run as independent processes. If one crashes, the other continues. They communicate only through NATS.

3. **Two LLMs** — A powerful slow model for thinking (LLM #1) and a fast local model for physical reflexes (LLM #2). Neither replaces the other.

4. **Two Heartbeats** — Agno's slow heartbeat (15 min, digital awareness) and OM1's fast heartbeat (0.5 sec, physical survival) run independently but share context through NATS.

5. **One Session** — All interaction channels (voice, text, Telegram, robot mic) share the same Agno `session_id`. One conversation, one memory, regardless of input source.

6. **No LiveKit** — The Interaction Gateway + Agno's native session system replaces LiveKit rooms entirely. ~200 lines of FastAPI vs. a full media server cluster.

## Documentation Index

| Document | Description |
|----------|-------------|
| [INTERACTION-GATEWAY.md](INTERACTION-GATEWAY.md) | Voice interface, chat, channel adapters, STT/TTS |
| [AGNO-LAYER.md](./AGNO-LAYER.md) | Cognitive layer — knowledge, memory, teams, MCP, tools |
| [OM1-LAYER.md](./OM1-LAYER.md) | Embodiment layer — sensors, motors, HAL, robot control |
| [NATS-BUS.md](./NATS-BUS.md) | Message bus — subjects, payloads, bridge plugins |
| [HEARTBEAT.md](00-24%20SO%20(Sovereign%20Operations)/01%20Intelligence%20(AI,%20Memory%20&%20Tools)%20(The%20brain)/architecture/HEARTBEAT.md) | Dual heartbeat system — Agno slow + OM1 fast |
| [UNIFIED-SESSION.md](UNIFIED-SESSION.md) | Session management replacing LiveKit rooms |
| [LLM-STRATEGY.md](./LLM-STRATEGY.md) | Two-tier LLM architecture — thinker + reflex |
| [PERSONALITY.md](./PERSONALITY.md) | Unified soul/identity across digital and physical |
| [SHARED-VISION.md](./SHARED-VISION.md) | Multi-robot vision sharing and world state |
| [INTEGRATION-GUIDE.md](INTEGRATION-GUIDE.md) | Step-by-step implementation plan |

## Tech Stack

| Component | Technology |
|-----------|------------|
| Cognitive Framework | Agno (Python) |
| Embodiment Runtime | OpenMind OM1 (Python) |
| Message Bus | NATS |
| Database | PostgreSQL + PgVector |
| Frontend | Flutter (mobile/desktop) |
| Voice STT/TTS | Gemini Live API + ElevenLabs |
| LLM #1 (Thinker) | GPT-4.1 / Claude / Gemini 2.5 Pro |
| LLM #2 (Reflex) | Qwen3-30B local / GPT-4.1-nano |
| Robot Hardware | Unitree G1, Go2, TurtleBot 4 |
| Robot Middleware | ROS2 / CycloneDDS / Zenoh |
| AR Glasses | Brilliant Labs Frame |
| Deployment | Railway (Agno) + On-device (OM1) |


---

## Related

- [[00-24 SO (Sovereign Operations)/01 Intelligence (AI, Memory & Tools) (The Brain)/01 Intelligence|01 Intelligence]]
