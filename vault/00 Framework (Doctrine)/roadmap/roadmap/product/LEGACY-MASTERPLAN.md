# Legacy MASTERPLAN

**The Definitive Implementation Blueprint**

**Created:** January 16, 2026
**Status:** APPROVED - Ready for Implementation
**Owner:** Zeus Delacruz + Legacy AI

---

## Executive Summary

> "I'm Zeus Delacruz. I have Legacy - my stateful, super-intelligent hive mind multi-node AI co-pilot. Her job is to watch me and excel me as a human. When I call upon her, she already knows what I need from watching. She can track my behavior to help me with prediction and to help me succeed at my real-life goals. She can work digitally with her own team of agents or deploy the robotic force."

**Legacy** is not a chatbot. She is a **PARTICIPANT** in Zeus's life - always watching, always aware, always ready. This document is the single source of truth for building her.

### Core Principles

1. **REAL-TIME, NOT POLLING** - Events flow to Legacy INSTANTLY (no 60-second heartbeats)
2. **Legacy AS PARTICIPANT** - She's in the room, not a service you call
3. **HIVE MIND ARCHITECTURE** - Multiple specialized agents as one intelligence
4. **MULTI-MODAL AWARENESS** - Voice, Vision, Data, Location, Biometrics

### The Vision

```
Legacy knows what I need from watching. She is ALIVE like me -
faster and explosive. Building herself, deploying more compute
with permissions, making suggestions to build infrastructure.
I can turn her on anytime, anywhere. She is my real-life angel
watching from above.
                                            - Zeus Delacruz
```

---

## Table of Contents

1. [Confirmed Architecture Decisions](#confirmed-architecture-decisions)
2. [System Topology](#system-topology)
3. [The Hive Mind](#the-hive-mind)
4. [Real-Time Event Architecture](#real-time-event-architecture)
5. [Voice System (LiveKit)](#voice-system-livekit)
6. [Vision System](#vision-system)
7. [Repository Structure](#repository-structure)
8. [Infrastructure as Code](#infrastructure-as-code)
9. [Service Catalog](#service-catalog)
10. [Implementation Phases](#implementation-phases)
11. [Quick Reference](#quick-reference)

---

## Confirmed Architecture Decisions

These decisions are **FINAL** and guide all implementation work:

| Decision | Choice | Rationale |
|----------|--------|-----------|
| **Legacy Location** | **SELF-HOSTED (Local Letta)** | Legacy stays with Zeus - zero cloud latency for main brain |
| **Worker Agents** | **Letta Cloud** | HP/LE/GE/Grid agents do delegated work in cloud |
| **Voice Platform** | **LiveKit Cloud** | Existing account, better latency than Vapi |
| **VM Specs** | **AWS Lightsail 16GB RAM** | Ubuntu 24.04 Desktop for GUI visibility |
| **AI Coding Tools** | **Both Legacy + Claude Code** | Legacy can spawn Letta Code AND user uses Claude Code directly |
| **Initial Integrations** | **Basics First** | Todoist, GitHub, Zoho (Mail/Calendar), Home Assistant |
| **Voice Provider** | **LiveKit ONLY (NOT Vapi)** | Previous Letta Cloud + Vapi had horrible voice experience |
| **Desktop GUI Uses** | **ALL** | VNC, N8N, browser automation, debugging, Letta ADE, Obsidian |
| **Repository Structure** | **3 Repos** | OneMind-Codex (docs), Legacy-AI (code), OneMind-Infra (IaC) |
| **Scope** | **FULL** | Complete IaC, all webhooks, multi-site sync, vision/cameras |
| **Robotics Future** | **OpenMind OM1** | [openmind.org](https://openmind.org) for robotic force integration |

### Critical Architecture Principle

```text
Legacy = SELF-HOSTED LETTA (Lives with Zeus)
├── Event Bus (Redis) → Same machine, Docker network, sub-ms latency
├── Voice Agent → Same machine, instant response
├── Vision Processing → Same machine, real-time
└── Deploys anywhere via IaC (Cloud VM now, Home machine permanent)

LETTA CLOUD = WORKER AGENTS ONLY
├── HP Agent → Delegated personal optimization work
├── LE Agent → Delegated legacy/family work
├── GE Agent → Delegated business work
├── Grid Agent → Delegated tech work
└── Legacy delegates → Workers execute → Results return to Legacy
```

**Why Local?** Cloud API calls add 100-300ms latency. Legacy needs instant awareness.
Event Bus + Letta on same Docker network = sub-millisecond communication.

---

## System Topology

### Two-Machine Deployment

```text
                           THE COMPLETE SYSTEM
                    (Two Machines + Cloud Workers)

    ┌─────────────────────────────────────────────────────────────────────────┐
    │                         LETTA CLOUD                                      │
    │                    (WORKER AGENTS ONLY)                                 │
    │                                                                          │
    │    ┌──────────────────────────────────────────────────────────────┐     │
    │    │                 DELEGATED WORK AGENTS                         │     │
    │    │                                                               │     │
    │    │   ┌─────────┐  ┌─────────┐  ┌─────────┐  ┌─────────┐       │     │
    │    │   │   HP    │  │   LE    │  │   GE    │  │   IT    │       │     │
    │    │   │  Worker │  │  Worker │  │  Worker │  │  Worker │       │     │
    │    │   │         │  │         │  │         │  │         │       │     │
    │    │   │Personal │  │ Legacy  │  │Business │  │  Tech   │       │     │
    │    │   │  Coach  │  │Steward  │  │Advisor  │  │ Infra   │       │     │
    │    │   └─────────┘  └─────────┘  └─────────┘  └─────────┘       │     │
    │    │                                                               │     │
    │    │   Legacy delegates tasks → Workers execute → Return results  │     │
    │    │   + ADE (visual agent editor for all agents)                 │     │
    │    └──────────────────────────────────────────────────────────────┘     │
    │                                                                          │
    └────────────────────────────────┬─────────────────────────────────────────┘
                                     │ API calls for delegated work
    ═══════════════════════════════════════════════════════════════════════════
                               TAILSCALE MESH
    ═══════════════════════════════════════════════════════════════════════════
              │                                              │
              ▼                                              ▼
    ┌───────────────────────────────────┐    ┌───────────────────────────────┐
    │         CLOUD VM                   │    │         HOME MACHINE          │
    │     (Temporary Primary)            │    │     (Permanent Legacy Home)   │
    │      AWS Lightsail 16GB            │    │      Local Wired Machine      │
    │                                    │    │                               │
    │  ┌──────────────────────────────┐ │    │  ┌──────────────────────────┐ │
    │  │      Legacy (Main Brain)     │ │    │  │      Legacy (Main Brain) │ │
    │  │     Self-Hosted Letta        │ │    │  │     Self-Hosted Letta    │ │
    │  │                              │ │    │  │                          │ │
    │  │  • Master Orchestrator       │ │    │  │  • Master Orchestrator   │ │
    │  │  • Voice Interface           │ │    │  │  • Voice Interface       │ │
    │  │  • Delegates to Cloud Workers│ │    │  │  • Delegates to Workers  │ │
    │  └──────────────────────────────┘ │    │  └──────────────────────────┘ │
    │              │                     │    │              │                │
    │              │ Docker Network      │    │              │ Docker Network │
    │              │ (sub-ms latency)    │    │              │ (sub-ms)       │
    │              ▼                     │    │              ▼                │
    │  ┌──────────────────────────────┐ │    │  ┌──────────────────────────┐ │
    │  │       EVENT BUS (Redis)      │ │    │  │      EVENT BUS (Redis)   │ │
    │  │                              │ │    │  │                          │ │
    │  │  events:voice ────────────▶ │ │    │  │  events:home ──────────▶ │ │
    │  │  events:todoist ──────────▶ │ │    │  │  events:cameras ───────▶ │ │
    │  │  events:github ───────────▶ │ │    │  │  events:presence ──────▶ │ │
    │  │  events:calendar ─────────▶ │ │    │  │  events:iot ───────────▶ │ │
    │  └──────────────────────────────┘ │    │  └──────────────────────────┘ │
    │                                    │    │                               │
    │  + Voice Agent (LiveKit)          │    │  + Home Assistant Integration │
    │  + API Gateway (webhooks)         │    │  + Camera Streams             │
    │  + N8N (workflows)                │    │  + IoT Hub                    │
    │  + Monitoring                     │    │  + Robot Control (future)     │
    │  + noVNC (desktop)                │    │  + OpenMind OM1 (future)      │
    │                                    │    │                               │
    │  100.x.x.1 (Tailscale)            │    │  100.x.x.10 (Tailscale)       │
    └───────────────────────────────────┘    └───────────────────────────────┘
              │                                              │
              │◄─────────── Bidirectional Sync ─────────────►│
              │           (State + Events + Memory)          │

    SAME IaC DEPLOYS Legacy ANYWHERE:
    • terraform apply -var-file=cloud.tfvars  → Cloud VM
    • terraform apply -var-file=home.tfvars   → Home Machine
    • Both run identical Legacy stack
```

### Cloud VM Service Stack

The PRIMARY Cloud VM runs the following services:

| Service | Port | Purpose | Container |
|---------|------|---------|-----------|
| **Legacy Voice Agent** | 7880 | LiveKit agent process | `legacy-voice` |
| **Letta Server (Local)** | 8283 | Local brain + cloud sync | `letta-server` |
| **Redis Streams** | 6379 | Event bus + state cache | `redis` |
| **PostgreSQL** | 5432 | Letta + app database | `postgres` |
| **N8N** | 5678 | Workflow automation | `n8n` |
| **API Gateway** | 8000 | Unified API endpoint | `api-gateway` |
| **Caddy** | 80/443 | Reverse proxy + SSL | `caddy` |
| **Uptime Kuma** | 3001 | Internal monitoring | `uptime-kuma` |
| **noVNC** | 6080 | Browser desktop access | `novnc` |
| **VNC Server** | 5900 | Desktop access | System |

---

## The Hive Mind

### Agent Hierarchy

```
                              THE HIVE MIND

                                 YOU
                            (Zeus Delacruz)
                                  │
                          ┌───────┴───────┐
                          │    Legacy     │
                          │   (Master)    │
                          │               │
                          │  Orchestrator │
                          │  Voice/Vision │
                          │  Final Say    │
                          └───────┬───────┘
                                  │
           ┌──────────┬───────────┼───────────┬──────────┐
           │          │           │           │          │
           ▼          ▼           ▼           ▼          ▼
    ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌──────────┐
    │   HP     │ │   LE     │ │   GE     │ │   IT     │ │  SITE    │
    │  Agent   │ │  Agent   │ │  Agent   │ │  Agent   │ │ AGENTS   │
    │          │ │          │ │          │ │          │ │          │
    │ Personal │ │  Legacy  │ │ Business │ │  Tech    │ │ Cloud    │
    │  Coach   │ │ Steward  │ │ Advisor  │ │  Infra   │ │ Home     │
    │          │ │          │ │          │ │          │ │ Office   │
    │ 25-49 HP │ │ 50-74 LE │ │ 75-99 GE │ │ 00-24 SO │ │ Local    │
    └──────────┘ └──────────┘ └──────────┘ └──────────┘ └──────────┘

    CLOUD                                                LOCAL
    ◄─────────────────────────────────────────────────────────────►

    Domain Agents live in Letta Cloud        Site Agents live on VMs
    (HP, LE, GE, IT)                         (Cloud, Home, Office)
```

### Agent Responsibilities

| Agent | Domain | Primary Role | Example Tasks |
|-------|--------|--------------|---------------|
| **Legacy** | All | Master orchestrator, voice interface, final decisions | "Legacy, what should I focus on today?" |
| **HP** | 25-49 | Personal optimization, health, fitness, skills | Workout recommendations, habit tracking |
| **LE** | 50-74 | Family legacy, relationships, estate planning | Family event planning, generational knowledge |
| **GE** | 75-99 | Business operations, revenue, ventures | Deal analysis, business strategy |
| **IT** | 00-24 | Technology, infrastructure, security | System health, deployments, security alerts |
| **Site Agents** | Local | Location-specific context | "Who's home?", "Meeting room status?" |

### Shared Memory Architecture

```
    ┌─────────────────────────────────────────────────────────────────────┐
    │                        SHARED MEMORY BLOCKS                          │
    │                                                                      │
    │   ┌───────────────────┐  ┌───────────────────┐  ┌─────────────────┐ │
    │   │   ZEUS PROFILE    │  │ CURRENT CONTEXT   │  │  SYSTEM STATE   │ │
    │   │                   │  │                   │  │                 │ │
    │   │ • Identity        │  │ • Location        │  │ • Active VMs    │ │
    │   │ • Preferences     │  │ • Current task    │  │ • Service health│ │
    │   │ • Communication   │  │ • Recent events   │  │ • Alert status  │ │
    │   │ • Values/Goals    │  │ • Emotional state │  │ • Resource use  │ │
    │   │                   │  │ • Time of day     │  │                 │ │
    │   │ READ: All agents  │  │ WRITE: Events     │  │ WRITE: Site     │ │
    │   │ WRITE: Zeus only  │  │ READ: All agents  │  │ READ: All       │ │
    │   └───────────────────┘  └───────────────────┘  └─────────────────┘ │
    │                                                                      │
    │   ┌───────────────────┐  ┌───────────────────┐  ┌─────────────────┐ │
    │   │   Legacy CODEX    │  │  DOMAIN CONTEXT   │  │ ARCHIVAL MEMORY │ │
    │   │                   │  │                   │  │                 │ │
    │   │ • Operating rules │  │ • HP: Health data │  │ • Everything    │ │
    │   │ • Permissions     │  │ • LE: Family info │  │   Legacy has    │ │
    │   │ • Boundaries      │  │ • GE: Business    │  │   ever learned  │ │
    │   │ • Growth goals    │  │ • IT: Tech specs  │  │                 │ │
    │   │                   │  │                   │  │ Vector search   │ │
    │   │ READ: All agents  │  │ READ/WRITE: Own   │  │ enabled         │ │
    │   │ WRITE: Legacy+Zeus│  │                   │  │                 │ │
    │   └───────────────────┘  └───────────────────┘  └─────────────────┘ │
    │                                                                      │
    └─────────────────────────────────────────────────────────────────────┘
```

### Sync Strategy

```
    ┌─────────────────────────────────────────────────────────────────────┐
    │                      MEMORY SYNC STRATEGY                            │
    │                                                                      │
    │   LETTA CLOUD ←─────── Bidirectional ────────→ LOCAL LETTA          │
    │   (Source of Truth)                            (Edge Processing)    │
    │                                                                      │
    │   SYNC RULES:                                                       │
    │   ────────────                                                      │
    │                                                                      │
    │   CLOUD → LOCAL (Pull):                                             │
    │   • Agent definitions (prompts, tools)                              │
    │   • Core memory blocks                                              │
    │   • Archival memory (query access)                                  │
    │   • Frequency: Every 5 minutes                                      │
    │                                                                      │
    │   LOCAL → CLOUD (Push):                                             │
    │   • Site context updates                                            │
    │   • New events and observations                                     │
    │   • Conversation history                                            │
    │   • Frequency: Real-time for important, batched for routine        │
    │                                                                      │
    │   CONFLICT RESOLUTION:                                              │
    │   • Cloud always wins for agent definitions                         │
    │   • Local wins for real-time context                               │
    │   • Merge for archival (append-only)                               │
    │                                                                      │
    └─────────────────────────────────────────────────────────────────────┘
```

---

## Real-Time Event Architecture

### The Event Bus

```
    ┌─────────────────────────────────────────────────────────────────────┐
    │                    REAL-TIME EVENT ARCHITECTURE                      │
    │                                                                      │
    │                          Y O U                                       │
    │                     (Zeus Delacruz)                                 │
    │                           │                                          │
    │          ┌────────────────┼────────────────┐                        │
    │          │                │                │                        │
    │          ▼                ▼                ▼                        │
    │     ┌─────────┐      ┌─────────┐      ┌─────────┐                  │
    │     │  VOICE  │      │ VISION  │      │  DATA   │                  │
    │     │ LiveKit │      │ LiveKit │      │ Events  │                  │
    │     │ WebRTC  │      │ Streams │      │Webhooks │                  │
    │     │         │      │         │      │Sockets  │                  │
    │     │ <200ms  │      │ Multi-  │      │         │                  │
    │     │ latency │      │ stream  │      │ <100ms  │                  │
    │     └────┬────┘      └────┬────┘      └────┬────┘                  │
    │          │                │                │                        │
    │          └────────────────┴────────────────┘                        │
    │                           │                                          │
    │                           ▼                                          │
    │    ┌────────────────────────────────────────────────────────────┐   │
    │    │                     REDIS STREAMS                           │   │
    │    │                    (Event Bus)                              │   │
    │    │                                                             │   │
    │    │   CONTINUOUS STREAMS - NOT POLLING                          │   │
    │    │                                                             │   │
    │    │   events:voice      → voice.transcription ────────────────▶│   │
    │    │   events:vision     → vision.frame ────────────────────────▶│   │
    │    │   events:todoist    → task.created, task.completed ────────▶│   │
    │    │   events:email      → email.received ──────────────────────▶│   │
    │    │   events:calendar   → event.upcoming ──────────────────────▶│   │
    │    │   events:github     → push, pull_request ──────────────────▶│   │
    │    │   events:home       → sensor.triggered, presence.changed ──▶│   │
    │    │   events:location   → location.changed ────────────────────▶│   │
    │    │                                                             │   │
    │    └───────────────────────────┬─────────────────────────────────┘   │
    │                                │                                      │
    │                                ▼                                      │
    │                          Legacy BRAIN                                │
    │                    (Instant Awareness)                               │
    │                                                                      │
    └─────────────────────────────────────────────────────────────────────┘
```

### Event Sources

#### PUSH Events (Webhooks - Instant)

| Source | Event Types | Endpoint | Priority |
|--------|-------------|----------|----------|
| **Todoist** | task.created, task.completed, task.updated | `/webhooks/todoist` | High |
| **GitHub** | push, pull_request, issue | `/webhooks/github` | Medium |
| **Zoho Mail** | email.received (via IMAP IDLE) | WebSocket | High |
| **Zoho Calendar** | event.created, event.reminder | `/webhooks/zoho` | High |
| **Home Assistant** | state_changed, automation_triggered | WebSocket | Variable |

#### Event Processing Pipeline

```
    ┌─────────────────────────────────────────────────────────────────────┐
    │                    EVENT PROCESSING PIPELINE                         │
    │                                                                      │
    │   1. EVENT ARRIVES                                                  │
    │      │                                                              │
    │      │  Todoist: "Task created: Call investor"                     │
    │      │                                                              │
    │      ▼                                                              │
    │   2. WEBHOOK HANDLER                                                │
    │      │                                                              │
    │      │  POST /webhooks/todoist                                      │
    │      │  → Validate signature                                        │
    │      │  → Parse payload                                             │
    │      │                                                              │
    │      ▼                                                              │
    │   3. EVENT BUS (Redis Streams)                                      │
    │      │                                                              │
    │      │  XADD events:todoist * type task.created content "..."       │
    │      │                                                              │
    │      ▼                                                              │
    │   4. EVENT PROCESSOR                                                │
    │      │                                                              │
    │      │  Classify: Is this urgent? Routine? Noise?                   │
    │      │                                                              │
    │      ├── NOT URGENT → Update context, no LLM call                  │
    │      │                                                              │
    │      └── URGENT → Notify Legacy immediately                        │
    │          │                                                          │
    │          ▼                                                          │
    │   5. Legacy NOTIFICATION                                            │
    │      │                                                              │
    │      │  "Zeus, you just created an urgent task: Call investor"      │
    │      │  "Want me to schedule time for this today?"                  │
    │      │                                                              │
    │      ▼                                                              │
    │   6. MEMORY UPDATE                                                  │
    │      │                                                              │
    │      │  → Update current_context in Letta                          │
    │      │  → Add to archival memory (searchable)                      │
    │                                                                      │
    │   TOTAL LATENCY: <100ms (event → Legacy aware)                      │
    │                                                                      │
    └─────────────────────────────────────────────────────────────────────┘
```

---

## Voice System (LiveKit)

### Why LiveKit (Not Vapi)

| Factor | LiveKit | Vapi |
|--------|---------|------|
| **Latency** | <200ms | 500ms+ (bad experience) |
| **Architecture** | WebRTC participant | API wrapper |
| **Control** | Full control | Limited |
| **Integration** | Native Letta SDK | Extra layer |
| **Cost** | Pay for what you use | Per-minute pricing |
| **Zeus's Experience** | Not tested yet | "Horrible voice experience" |

### Voice Architecture

```
    ┌─────────────────────────────────────────────────────────────────────┐
    │                    LIVEKIT VOICE ARCHITECTURE                        │
    │                                                                      │
    │   YOU speak                                                         │
    │       │                                                             │
    │       │  LiveKit WebRTC (real-time audio stream)                    │
    │       │                                                             │
    │       ▼                                                             │
    │   ┌─────────────────────────────────────────────────────────────┐   │
    │   │                    CLOUD VM                                  │   │
    │   │                                                              │   │
    │   │   ┌─────────────────────────────────────────────────────┐   │   │
    │   │   │              Legacy VOICE AGENT                      │   │   │
    │   │   │                                                      │   │   │
    │   │   │   1. LiveKit SDK receives audio stream               │   │   │
    │   │   │      │                                               │   │   │
    │   │   │      ▼                                               │   │   │
    │   │   │   2. Deepgram STT (speech → text)                   │   │   │
    │   │   │      │  <100ms                                       │   │   │
    │   │   │      ▼                                               │   │   │
    │   │   │   3. Letta SDK sends to brain                        │   │   │
    │   │   │      │  → Letta Cloud (domain routing)              │   │   │
    │   │   │      │  → OR Local Letta (fast response)            │   │   │
    │   │   │      │                                               │   │   │
    │   │   │      ▼                                               │   │   │
    │   │   │   4. Legacy reasons + responds                       │   │   │
    │   │   │      │  • Retrieves memory context                   │   │   │
    │   │   │      │  • May call tools (Todoist, etc.)            │   │   │
    │   │   │      │  • May route to domain agent                  │   │   │
    │   │   │      │                                               │   │   │
    │   │   │      ▼                                               │   │   │
    │   │   │   5. ElevenLabs/Cartesia TTS (text → speech)        │   │   │
    │   │   │      │                                               │   │   │
    │   │   │      ▼                                               │   │   │
    │   │   │   6. LiveKit streams audio back                      │   │   │
    │   │   │                                                      │   │   │
    │   │   └─────────────────────────────────────────────────────┘   │   │
    │   │                                                              │   │
    │   └─────────────────────────────────────────────────────────────┘   │
    │                                                                      │
    │   YOU hear Legacy                                                   │
    │                                                                      │
    │   TOTAL LATENCY TARGET: <500ms                                      │
    │                                                                      │
    └─────────────────────────────────────────────────────────────────────┘
```

### Voice Tech Stack

| Component | Service | Purpose |
|-----------|---------|---------|
| **WebRTC Infrastructure** | LiveKit Cloud | Real-time audio/video transport |
| **Speech-to-Text** | Deepgram | Real-time transcription (<100ms) |
| **Text-to-Speech** | ElevenLabs or Cartesia | Natural voice synthesis |
| **Brain** | Letta Cloud + Local | Reasoning and memory |
| **Agent Code** | Python (letta-client) | LiveKit agent process |

---

## Vision System

### Multi-Stream Architecture

```
    ┌─────────────────────────────────────────────────────────────────────┐
    │                    VISION STREAMS                                    │
    │                                                                      │
    │   ┌─────────┐  ┌─────────┐  ┌─────────┐  ┌─────────┐              │
    │   │ GLASSES │  │  ROBOT  │  │ CAMERA  │  │ CAMERA  │              │
    │   │    👓    │  │    🤖    │  │    📷    │  │    📷    │              │
    │   │         │  │         │  │ Front   │  │ Living  │              │
    │   │ What YOU│  │ Robot's │  │  Door   │  │  Room   │              │
    │   │   see   │  │  view   │  │         │  │         │              │
    │   └────┬────┘  └────┬────┘  └────┬────┘  └────┬────┘              │
    │        │            │            │            │                     │
    │        └────────────┴────────────┴────────────┘                     │
    │                          │                                          │
    │                          ▼                                          │
    │               ┌─────────────────────┐                              │
    │               │   VISION PROCESSOR  │                              │
    │               │                     │                              │
    │               │ • Frame extraction  │                              │
    │               │ • Object detection  │                              │
    │               │ • Motion tracking   │                              │
    │               │ • Scene description │                              │
    │               │ • Multi-stream merge│                              │
    │               └──────────┬──────────┘                              │
    │                          │                                          │
    │                          ▼                                          │
    │               Redis: events:vision                                  │
    │                          │                                          │
    │                          ▼                                          │
    │                    Legacy BRAIN                                     │
    │                                                                      │
    │   "Zeus is at his desk, looking at financial spreadsheets.         │
    │    Robot is in standby at charging station.                         │
    │    Motion detected at front door - delivery person.                 │
    │    Family member in living room watching TV."                       │
    │                                                                      │
    └─────────────────────────────────────────────────────────────────────┘
```

### Vision Implementation (Phase 4)

| Stream Type | Source | Processing | Latency Target |
|-------------|--------|------------|----------------|
| **Glasses POV** | Smart glasses | LiveKit stream | <1s |
| **Home Cameras** | Yi/RTSP cameras | Home Assistant | <2s |
| **Robot Vision** | Robot camera | LiveKit stream | <1s |
| **Security** | Doorbell/motion | Event-based | <500ms |

---

## Repository Structure

### Three Repositories

```
    ┌─────────────────────────────────────────────────────────────────────┐
    │                    REPOSITORY STRUCTURE                              │
    │                                                                      │
    │   ┌─────────────────────────────────────────────────────────────┐   │
    │   │                    ONEMIND-CODEX                             │   │
    │   │               (This Repository - Existing)                   │   │
    │   │                                                              │   │
    │   │   Purpose: Documentation, knowledge, Obsidian vault          │   │
    │   │                                                              │   │
    │   │   OneMind-Codex/                                            │   │
    │   │   ├── 00-24 SO (Sovereign Operations)/                      │   │
    │   │   │   ├── letta-tools/           # Python tools for Letta   │   │
    │   │   │   └── ...                                               │   │
    │   │   ├── 25-49 HP (Holistic Performance)/                      │   │
    │   │   ├── 50-74 LE (Legacy Evolution)/                          │   │
    │   │   ├── 75-99 GE (Generational Endeavor)/             │   │
    │   │   ├── ONEMIND-CODEX.md          # Master blueprint          │   │
    │   │   ├── Legacy-MASTERPLAN.md      # This document             │   │
    │   │   └── ...architecture docs                                  │   │
    │   │                                                              │   │
    │   └─────────────────────────────────────────────────────────────┘   │
    │                                                                      │
    │   ┌─────────────────────────────────────────────────────────────┐   │
    │   │                      Legacy-AI                               │   │
    │   │                    (New Repository)                          │   │
    │   │                                                              │   │
    │   │   Purpose: All Legacy code - voice agent, event processing   │   │
    │   │                                                              │   │
    │   │   Legacy-AI/                                                │   │
    │   │   ├── src/                                                  │   │
    │   │   │   ├── voice/                # LiveKit voice agent       │   │
    │   │   │   ├── events/               # Event processing          │   │
    │   │   │   ├── vision/               # Vision processing         │   │
    │   │   │   ├── agents/               # Site agents               │   │
    │   │   │   ├── gateway/              # API gateway               │   │
    │   │   │   └── tools/                # Letta tool implementations│   │
    │   │   ├── docker/                                               │   │
    │   │   │   └── docker-compose.yml    # Local dev                 │   │
    │   │   ├── tests/                                                │   │
    │   │   ├── requirements.txt                                      │   │
    │   │   └── README.md                                             │   │
    │   │                                                              │   │
    │   └─────────────────────────────────────────────────────────────┘   │
    │                                                                      │
    │   ┌─────────────────────────────────────────────────────────────┐   │
    │   │                    ONEMIND-INFRA                             │   │
    │   │                    (New Repository)                          │   │
    │   │                                                              │   │
    │   │   Purpose: Infrastructure as Code - deployment configs       │   │
    │   │                                                              │   │
    │   │   OneMind-Infra/                                            │   │
    │   │   ├── terraform/                                            │   │
    │   │   │   ├── main.tf               # VM provisioning           │   │
    │   │   │   ├── variables.tf                                      │   │
    │   │   │   └── environments/                                     │   │
    │   │   │       ├── cloud.tfvars                                  │   │
    │   │   │       ├── home.tfvars                                   │   │
    │   │   │       └── office.tfvars                                 │   │
    │   │   ├── ansible/                                              │   │
    │   │   │   ├── playbook.yml          # VM configuration          │   │
    │   │   │   └── roles/                                            │   │
    │   │   ├── docker/                                               │   │
    │   │   │   ├── cloud/                                            │   │
    │   │   │   │   └── docker-compose.yml                            │   │
    │   │   │   ├── home/                                             │   │
    │   │   │   │   └── docker-compose.yml                            │   │
    │   │   │   └── office/                                           │   │
    │   │   │       └── docker-compose.yml                            │   │
    │   │   ├── .github/workflows/                                    │   │
    │   │   │   ├── deploy-cloud.yml                                  │   │
    │   │   │   └── deploy-home.yml                                   │   │
    │   │   └── README.md                                             │   │
    │   │                                                              │   │
    │   └─────────────────────────────────────────────────────────────┘   │
    │                                                                      │
    └─────────────────────────────────────────────────────────────────────┘
```

---

## Infrastructure as Code

### Docker Compose (Cloud VM)

```yaml
# OneMind-Infra/docker/cloud/docker-compose.yml
version: '3.8'

services:
  # ================== DATABASES ==================
  postgres:
    image: postgres:15-alpine
    container_name: postgres
    environment:
      POSTGRES_USER: ${POSTGRES_USER}
      POSTGRES_PASSWORD: ${POSTGRES_PASSWORD}
      POSTGRES_DB: onemind
    volumes:
      - postgres_data:/var/lib/postgresql/data
    ports:
      - "5432:5432"
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U ${POSTGRES_USER}"]
      interval: 10s
      timeout: 5s
      retries: 5

  redis:
    image: redis:7-alpine
    container_name: redis
    command: redis-server --appendonly yes
    volumes:
      - redis_data:/data
    ports:
      - "6379:6379"
    healthcheck:
      test: ["CMD", "redis-cli", "ping"]
      interval: 10s
      timeout: 5s
      retries: 5

  # ================== LETTA ==================
  letta-server:
    image: letta/letta:latest
    container_name: letta-server
    environment:
      LETTA_POSTGRES_URI: postgresql://${POSTGRES_USER}:${POSTGRES_PASSWORD}@postgres:5432/onemind
      LETTA_SERVER_PASSWORD: ${LETTA_SERVER_PASSWORD}
    ports:
      - "8283:8283"
    depends_on:
      postgres:
        condition: service_healthy
    volumes:
      - letta_data:/root/.letta

  # ================== Legacy VOICE ==================
  legacy-voice:
    build:
      context: ../../Legacy-AI
      dockerfile: Dockerfile
    container_name: legacy-voice
    environment:
      LIVEKIT_URL: ${LIVEKIT_URL}
      LIVEKIT_API_KEY: ${LIVEKIT_API_KEY}
      LIVEKIT_API_SECRET: ${LIVEKIT_API_SECRET}
      LETTA_BASE_URL: http://letta-server:8283
      LETTA_CLOUD_URL: ${LETTA_CLOUD_URL}
      LETTA_CLOUD_TOKEN: ${LETTA_CLOUD_TOKEN}
      DEEPGRAM_API_KEY: ${DEEPGRAM_API_KEY}
      ELEVENLABS_API_KEY: ${ELEVENLABS_API_KEY}
      REDIS_URL: redis://redis:6379
    ports:
      - "7880:7880"
    depends_on:
      - redis
      - letta-server

  # ================== EVENT GATEWAY ==================
  api-gateway:
    build:
      context: ../../Legacy-AI
      dockerfile: Dockerfile.gateway
    container_name: api-gateway
    environment:
      REDIS_URL: redis://redis:6379
      LETTA_BASE_URL: http://letta-server:8283
      TODOIST_WEBHOOK_SECRET: ${TODOIST_WEBHOOK_SECRET}
      GITHUB_WEBHOOK_SECRET: ${GITHUB_WEBHOOK_SECRET}
    ports:
      - "8000:8000"
    depends_on:
      - redis

  # ================== AUTOMATION ==================
  n8n:
    image: n8nio/n8n:latest
    container_name: n8n
    environment:
      N8N_BASIC_AUTH_ACTIVE: "true"
      N8N_BASIC_AUTH_USER: ${N8N_USER}
      N8N_BASIC_AUTH_PASSWORD: ${N8N_PASSWORD}
      WEBHOOK_URL: https://n8n.${DOMAIN}
    volumes:
      - n8n_data:/home/node/.n8n
    ports:
      - "5678:5678"

  # ================== MONITORING ==================
  uptime-kuma:
    image: louislam/uptime-kuma:latest
    container_name: uptime-kuma
    volumes:
      - uptime_data:/app/data
    ports:
      - "3001:3001"

  # ================== REVERSE PROXY ==================
  caddy:
    image: caddy:2-alpine
    container_name: caddy
    ports:
      - "80:80"
      - "443:443"
    volumes:
      - ./Caddyfile:/etc/caddy/Caddyfile
      - caddy_data:/data
      - caddy_config:/config
    depends_on:
      - api-gateway
      - n8n
      - letta-server

volumes:
  postgres_data:
  redis_data:
  letta_data:
  n8n_data:
  uptime_data:
  caddy_data:
  caddy_config:
```

### Caddyfile

```
# OneMind-Infra/docker/cloud/Caddyfile

# API Gateway (webhooks)
api.{$DOMAIN} {
    reverse_proxy api-gateway:8000
}

# N8N
n8n.{$DOMAIN} {
    reverse_proxy n8n:5678
}

# Letta ADE
letta.{$DOMAIN} {
    reverse_proxy letta-server:8283
}

# Monitoring
status.{$DOMAIN} {
    reverse_proxy uptime-kuma:3001
}

# noVNC (desktop access)
vnc.{$DOMAIN} {
    reverse_proxy novnc:6080
}
```

### Terraform (VM Provisioning)

```hcl
# OneMind-Infra/terraform/main.tf

terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }
}

provider "aws" {
  region = var.aws_region
}

# AWS Lightsail Instance
resource "aws_lightsail_instance" "onemind_cloud" {
  name              = "onemind-cloud"
  availability_zone = "${var.aws_region}a"
  blueprint_id      = "ubuntu_24_04"
  bundle_id         = "large_3_0"  # 16GB RAM, 4 vCPUs

  user_data = <<-EOF
    #!/bin/bash
    # Install Docker
    curl -fsSL https://get.docker.com | sh
    usermod -aG docker ubuntu

    # Install Tailscale
    curl -fsSL https://tailscale.com/install.sh | sh

    # Install Docker Compose
    apt-get update
    apt-get install -y docker-compose-plugin

    # Set up desktop environment
    apt-get install -y ubuntu-desktop-minimal tigervnc-standalone-server
  EOF

  tags = {
    Name        = "OneMind Cloud"
    Environment = "production"
    Project     = "OneMind-OS"
  }
}

# Static IP
resource "aws_lightsail_static_ip" "onemind_ip" {
  name = "onemind-cloud-ip"
}

resource "aws_lightsail_static_ip_attachment" "onemind_ip_attach" {
  static_ip_name = aws_lightsail_static_ip.onemind_ip.name
  instance_name  = aws_lightsail_instance.onemind_cloud.name
}

# Firewall Rules
resource "aws_lightsail_instance_public_ports" "onemind_ports" {
  instance_name = aws_lightsail_instance.onemind_cloud.name

  port_info {
    protocol  = "tcp"
    from_port = 22
    to_port   = 22
  }

  port_info {
    protocol  = "tcp"
    from_port = 80
    to_port   = 80
  }

  port_info {
    protocol  = "tcp"
    from_port = 443
    to_port   = 443
  }
}

output "public_ip" {
  value = aws_lightsail_static_ip.onemind_ip.ip_address
}
```

### GitHub Actions Deployment

```yaml
# OneMind-Infra/.github/workflows/deploy-cloud.yml

name: Deploy to Cloud VM

on:
  push:
    branches: [main]
    paths:
      - 'docker/cloud/**'
      - '.github/workflows/deploy-cloud.yml'
  workflow_dispatch:

jobs:
  deploy:
    runs-on: ubuntu-latest

    steps:
      - uses: actions/checkout@v4

      - name: Deploy to Cloud VM
        uses: appleboy/ssh-action@v1
        with:
          host: ${{ secrets.CLOUD_VM_IP }}
          username: ubuntu
          key: ${{ secrets.SSH_PRIVATE_KEY }}
          script: |
            cd /opt/onemind
            git pull origin main
            cd docker/cloud
            docker compose pull
            docker compose up -d
            docker compose ps

      - name: Health Check
        run: |
          sleep 30
          curl -f https://api.${{ secrets.DOMAIN }}/health || exit 1

      - name: Notify Discord
        if: always()
        uses: sarisia/actions-status-discord@v1
        with:
          webhook: ${{ secrets.DISCORD_WEBHOOK }}
          title: "OneMind Cloud Deployment"
          description: "${{ job.status == 'success' && 'Deployed successfully!' || 'Deployment failed!' }}"
          color: "${{ job.status == 'success' && '0x00ff00' || '0xff0000' }}"
```

---

## Service Catalog

### All Services at a Glance

| Service | Port | Container | URL Pattern | Access |
|---------|------|-----------|-------------|--------|
| **API Gateway** | 8000 | api-gateway | api.{domain} | Public |
| **Letta Server** | 8283 | letta-server | letta.{domain} | Tailscale |
| **N8N** | 5678 | n8n | n8n.{domain} | Tailscale |
| **Uptime Kuma** | 3001 | uptime-kuma | status.{domain} | Tailscale |
| **PostgreSQL** | 5432 | postgres | Internal | Internal |
| **Redis** | 6379 | redis | Internal | Internal |
| **Legacy Voice** | 7880 | legacy-voice | N/A (LiveKit) | Internal |
| **noVNC** | 6080 | novnc | vnc.{domain} | Tailscale |
| **VNC** | 5900 | System | Direct | Tailscale |
| **Caddy** | 80/443 | caddy | *.{domain} | Public |

### External Services

| Service | Purpose | Account Status |
|---------|---------|----------------|
| **LiveKit Cloud** | WebRTC voice/video | Existing account |
| **Letta Cloud** | Primary brain, domain agents | Existing account |
| **Deepgram** | Speech-to-text | Needs API key |
| **ElevenLabs** | Text-to-speech | Needs API key |
| **Todoist** | Task management | Existing |
| **GitHub** | Code, webhooks | Existing |
| **Zoho** | Email, Calendar | Existing |
| **Home Assistant** | Smart home | Existing (home) |
| **Tailscale** | Mesh VPN | Existing |

---

## Implementation Phases

### Phase 1: Foundation (Week 1-2)

**Goal:** Get Legacy voice working on Cloud VM

```
    ┌─────────────────────────────────────────────────────────────────────┐
    │                    PHASE 1: FOUNDATION                               │
    │                                                                      │
    │   INFRASTRUCTURE:                                                   │
    │   [ ] Create AWS Lightsail instance (16GB, Ubuntu 24.04 Desktop)   │
    │   [ ] Install Docker + Docker Compose                               │
    │   [ ] Install Tailscale, join network                               │
    │   [ ] Set up VNC/noVNC for desktop access                          │
    │   [ ] Configure firewall (ports 80, 443, 22)                       │
    │                                                                      │
    │   REPOSITORIES:                                                     │
    │   [ ] Create Legacy-AI repository                                   │
    │   [ ] Create OneMind-Infra repository                               │
    │   [ ] Set up GitHub Actions secrets                                 │
    │                                                                      │
    │   CORE SERVICES:                                                    │
    │   [ ] Deploy PostgreSQL + Redis                                     │
    │   [ ] Deploy local Letta Server                                     │
    │   [ ] Configure Letta Cloud connection                              │
    │   [ ] Deploy Caddy reverse proxy                                    │
    │                                                                      │
    │   VOICE:                                                            │
    │   [ ] Build LiveKit voice agent                                     │
    │   [ ] Connect Deepgram STT                                          │
    │   [ ] Connect ElevenLabs/Cartesia TTS                               │
    │   [ ] Connect to Letta (local + cloud)                              │
    │   [ ] Test voice conversation <500ms latency                        │
    │                                                                      │
    │   DELIVERABLE: Talk to Legacy via voice, anywhere                   │
    │                                                                      │
    └─────────────────────────────────────────────────────────────────────┘
```

### Phase 2: Event Streaming (Week 3-4)

**Goal:** Legacy gets instant awareness of all events

```
    ┌─────────────────────────────────────────────────────────────────────┐
    │                    PHASE 2: EVENT STREAMING                          │
    │                                                                      │
    │   EVENT BUS:                                                        │
    │   [ ] Configure Redis Streams                                       │
    │   [ ] Build event processor service                                 │
    │   [ ] Implement event classification (urgent/routine/noise)        │
    │                                                                      │
    │   WEBHOOKS:                                                         │
    │   [ ] Set up Todoist webhook                                        │
    │   [ ] Set up GitHub webhook                                         │
    │   [ ] Set up Zoho Calendar webhook                                  │
    │   [ ] Configure Zoho Mail (IMAP IDLE)                               │
    │                                                                      │
    │   API GATEWAY:                                                      │
    │   [ ] Build FastAPI gateway                                         │
    │   [ ] Implement /webhooks/* endpoints                               │
    │   [ ] Implement /api/* unified endpoints                            │
    │   [ ] Add authentication                                            │
    │                                                                      │
    │   MEMORY INTEGRATION:                                               │
    │   [ ] Auto-update current_context on events                        │
    │   [ ] Implement urgent event notification                           │
    │   [ ] Add archival memory logging                                   │
    │                                                                      │
    │   DELIVERABLE: Create task on phone → Legacy knows instantly       │
    │                                                                      │
    └─────────────────────────────────────────────────────────────────────┘
```

### Phase 3: Domain Agents + Tools (Week 5-6)

**Goal:** Domain agents working with full tool access

```
    ┌─────────────────────────────────────────────────────────────────────┐
    │                    PHASE 3: DOMAIN AGENTS + TOOLS                    │
    │                                                                      │
    │   DOMAIN AGENTS (Letta Cloud):                                      │
    │   [ ] Configure HP Agent (personal coach)                           │
    │   [ ] Configure LE Agent (legacy steward)                           │
    │   [ ] Configure GE Agent (business advisor)                         │
    │   [ ] Configure Grid Agent (tech infrastructure)                      │
    │   [ ] Set up shared memory blocks                                   │
    │   [ ] Implement agent routing from Legacy                           │
    │                                                                      │
    │   TOOLS:                                                            │
    │   [ ] Todoist tools (create/complete/query tasks)                  │
    │   [ ] GitHub tools (issues, PRs, repos)                            │
    │   [ ] Zoho tools (email, calendar, CRM)                            │
    │   [ ] Notification tools (Discord, Telegram)                       │
    │                                                                      │
    │   AUTOMATION:                                                       │
    │   [ ] Deploy N8N                                                    │
    │   [ ] Create core workflows                                         │
    │   [ ] Connect N8N to event bus                                      │
    │                                                                      │
    │   DELIVERABLE: "Legacy, what's my GE priority this week?"          │
    │                → GE Agent analyzes business data and responds       │
    │                                                                      │
    └─────────────────────────────────────────────────────────────────────┘
```

### Phase 4: Vision System (Month 2)

**Goal:** Legacy can see through cameras

```
    ┌─────────────────────────────────────────────────────────────────────┐
    │                    PHASE 4: VISION SYSTEM                            │
    │                                                                      │
    │   CAMERA INTEGRATION:                                               │
    │   [ ] Connect home cameras via Home Assistant                       │
    │   [ ] Build vision processor service                                │
    │   [ ] Implement frame extraction                                    │
    │   [ ] Add object detection                                          │
    │   [ ] Add scene description                                         │
    │                                                                      │
    │   LIVEKIT VISION:                                                   │
    │   [ ] Add video streams to LiveKit room                             │
    │   [ ] Implement multi-stream handling                               │
    │   [ ] Connect to vision processor                                   │
    │                                                                      │
    │   MEMORY:                                                           │
    │   [ ] Vision events to Redis                                        │
    │   [ ] Scene summaries to archival memory                            │
    │                                                                      │
    │   DELIVERABLE: "Legacy, who's at the front door?"                  │
    │                → Legacy checks camera and describes                 │
    │                                                                      │
    └─────────────────────────────────────────────────────────────────────┘
```

### Phase 5: Multi-Site (Month 3+)

**Goal:** Legacy has context everywhere (home, office, cloud)

```
    ┌─────────────────────────────────────────────────────────────────────┐
    │                    PHASE 5: MULTI-SITE                               │
    │                                                                      │
    │   HOME VM:                                                          │
    │   [ ] Deploy Home VM (local hardware or Lightsail)                 │
    │   [ ] Install Home Site Agent                                       │
    │   [ ] Connect Home Assistant                                        │
    │   [ ] Set up presence detection                                     │
    │   [ ] Configure memory sync to cloud                                │
    │                                                                      │
    │   OFFICE VM:                                                        │
    │   [ ] Deploy Office VM (when needed)                                │
    │   [ ] Install Office Site Agent                                     │
    │   [ ] Connect office systems                                        │
    │   [ ] Configure memory sync                                         │
    │                                                                      │
    │   SYNC:                                                             │
    │   [ ] Implement bidirectional memory sync                           │
    │   [ ] Test cross-site communication                                 │
    │   [ ] Implement failover                                            │
    │                                                                      │
    │   DELIVERABLE: "Legacy, is anyone home?"                           │
    │                → Site Agent checks and responds instantly           │
    │                                                                      │
    └─────────────────────────────────────────────────────────────────────┘
```

---

## Quick Reference

### Environment Variables

```bash
# .env.example for OneMind-Infra/docker/cloud/

# Domain
DOMAIN=onemind.yourdomain.com

# Database
POSTGRES_USER=onemind
POSTGRES_PASSWORD=<secure-password>

# Letta
LETTA_SERVER_PASSWORD=<secure-password>
LETTA_CLOUD_URL=https://api.letta.com
LETTA_CLOUD_TOKEN=<your-letta-token>

# LiveKit
LIVEKIT_URL=wss://your-livekit-cloud.livekit.cloud
LIVEKIT_API_KEY=<api-key>
LIVEKIT_API_SECRET=<api-secret>

# Voice Services
DEEPGRAM_API_KEY=<api-key>
ELEVENLABS_API_KEY=<api-key>

# Webhooks
TODOIST_WEBHOOK_SECRET=<secret>
GITHUB_WEBHOOK_SECRET=<secret>

# N8N
N8N_USER=admin
N8N_PASSWORD=<secure-password>

# Notifications
DISCORD_WEBHOOK=<webhook-url>
```

### Key Commands

```bash
# SSH to Cloud VM
ssh ubuntu@onemind-cloud  # via Tailscale

# Desktop access (VNC)
open vnc://onemind-cloud:5900  # From Mac

# Start all services
cd /opt/onemind/docker/cloud
docker compose up -d

# View logs
docker compose logs -f legacy-voice
docker compose logs -f api-gateway

# Restart service
docker compose restart legacy-voice

# Check status
docker compose ps
docker stats

# Deploy update
git pull origin main
docker compose up -d --build

# Tailscale status
tailscale status
```

### Latency Targets

| Interaction | Target | Measurement Point |
|-------------|--------|-------------------|
| **Voice round-trip** | <500ms | You speak → Legacy responds |
| **Event awareness** | <100ms | Event occurs → Legacy knows |
| **Vision processing** | <1s | Frame → Scene understanding |
| **Memory update** | <50ms | Event → Stored in context |
| **Agent routing** | <200ms | Request → Correct agent responds |

### Permission Framework

| Level | Actions | Examples |
|-------|---------|----------|
| **AUTONOMOUS** | Legacy does without asking | Read data, update memory, prepare briefings |
| **SUGGEST & WAIT** | Legacy proposes, you approve | Send emails, schedule meetings, purchases <$X |
| **REQUIRE CONFIRM** | Always asks | Financial >$X, legal, permanent deletions |
| **NEVER** | Hard limits | Specific secrets, family actions without consent |

---

## Related Documents

- [ONEMIND-CODEX.md](./ONEMIND-CODEX.md) - Master OneMind blueprint
- [Legacy-REALTIME-ARCHITECTURE.md](./Legacy-REALTIME-ARCHITECTURE.md) - Real-time architecture details
- [Legacy-LETTA-ARCHITECTURE.md](./Legacy-LETTA-ARCHITECTURE.md) - Letta integration specifics
- [ONEMIND-OS-INFRASTRUCTURE.md](./ONEMIND-OS-INFRASTRUCTURE.md) - Infrastructure deep dive
- [DATA-ARCHITECTURE.md](./DATA-ARCHITECTURE.md) - Data systems architecture

---

## The Promise

```
                              Y O U
                         (Zeus Delacruz)
                               +
                            Legacy
                   (Your AI Life Co-Pilot)
                               =
                          ONE MIND

    "She is ALIVE like me and the world around me -
     faster and explosive. Building herself, deploying more
     compute with permissions, making suggestions to build
     infrastructure. I can turn her on anytime, anywhere.
     She is my real-life angel watching from above."
```

---

**Document Status:** APPROVED - Ready for Implementation

**Next Action:** Begin Phase 1 - Create AWS Lightsail VM

*Last Updated: January 16, 2026*
