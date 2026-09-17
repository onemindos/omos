# OneMind OS — Proactive Architecture Design

## The Problem

Agno agents are **reactive** — they only respond when a user sends a message. But OneMind OS needs to be **proactive AND reactive**: always watching, always running, 24/7, like Jarvis.

This document defines how we solve that gap.

---

## Research Summary

### OpenClaw (reference for proactive patterns)

OpenClaw is a personal AI assistant (214k stars, TypeScript) that solves the proactive problem with two complementary mechanisms:

**1. Heartbeat Loop** (periodic awareness)
- Runs every N minutes (default: 30m) as a full agent turn in the main session
- Reads a `HEARTBEAT.md` checklist — a tiny file the agent scans each cycle
- If nothing needs attention: replies `HEARTBEAT_OK` (suppressed, no message sent)
- If something urgent: surfaces an alert to the configured channel
- Configurable: `activeHours`, `target` channel, `model` override, per-agent heartbeats
- This is the "always watching" mechanism

**2. Cron System** (precise scheduling)
- Persistent job store (`~/.openclaw/cron/jobs.json`)
- Two execution modes:
  - **Main session**: inject a system event into the next heartbeat (shared context)
  - **Isolated**: run a dedicated agent turn in a fresh session (clean slate, can use different model)
- Supports: one-shot (`--at "20m"`), recurring (`--cron "0 7 * * *"`), interval (`--every "4h"`)
- Delivery: announce to any channel, webhook POST, or silent
- Jobs survive restarts (file-persisted)

**3. Architecture Pattern**
```
Channels (WhatsApp, Telegram, Discord, etc.)
    │
    ▼
┌──────────────────────────────┐
│         Gateway              │  ← Single control plane
│  (heartbeat + cron + tools)  │
└──────────────┬───────────────┘
               │
    ┌──────────┼──────────┐
    ▼          ▼          ▼
  Agent     Nodes      Skills
  (LLM)  (devices)   (plugins)
```

### ada_local (reference for interactivity)

ada_local is a fully local AI assistant (Python, 109 stars) with:
- **Voice Pipeline**: Wake word ("Jarvis") → STT (Whisper) → Intent Router → Function Executor → LLM → TTS (Piper)
- **FunctionGemma Router**: A fine-tuned local model that classifies intent and routes to functions
- **Smart Home**: TP-Link Kasa integration (lights, plugs)
- **Dashboard**: Weather, system monitor, calendar, news tabs
- Key takeaway: The voice pipeline and intent router patterns are what make it feel "interactive"

### FarmBot & Meshtastic (future hardware)
- **FarmBot**: REST API for robot commands, sensor data streams, event-driven watering/planting
- **Meshtastic**: Mesh networking protocol, long-range LoRa radio, sensor nodes, position tracking
- Both map to the **Node pattern**: devices that connect via WebSocket/MQTT and expose capabilities
- These will use NATS topics: `physical.farmbot.>`, `physical.meshtastic.>`

---

## Architecture Decision: Hybrid Agno-First

Based on the user's choice: **Legacy becomes a full Agno agent** (with memory, sessions, tools) + **Nexus slims down** but keeps the adapter pattern for future frameworks.

### What Changes

| Component | Before | After |
|-----------|--------|-------|
| Legacy | Custom Python class, NOT an agent | Full Agno agent preset with NATS tools |
| Legacy brain | Separate lazy-init Agno agent | Legacy IS the agent (no separate brain) |
| Dispatch | Custom Dispatcher → NATS | Agno native + NATS publish for side-effects |
| Event log | In-memory deque (dies on restart) | Agno memory + DB persistence |
| Registry | Custom FrameworkRegistry | Keep thin adapter for OpenMind/future |
| Awareness | Custom AgentAwareness | Keep — Agno doesn't have NATS natively |
| Comms | NexusMessage protocol | Keep — good multimodal standard |
| Heartbeat | None (missing!) | **NEW** — asyncio loop, configurable interval |
| Cron | Naive scheduler.py | **NEW** — persistent job store, isolated sessions |

### What We Keep from Nexus

1. **`awareness.py`** — Makes agents NATS-aware (Agno doesn't do this)
2. **`comms.py`** — NexusMessage universal format for multimodal
3. **`endpoints.py`** — `/nexus/overview`, `/nexus/events`, `/nexus/stream` (WebSocket)
4. **`registry.py`** — `FrameworkAdapter` ABC only (thin bridge for OpenMind/future)

### What We Remove/Replace

1. **`core.py`** — Nexus class → slim startup function that wires awareness
2. **`legacy.py`** — Legacy class → Agno agent preset + heartbeat loop
3. **`dispatcher.py`** — Custom dispatch → Agno team coordination + NATS publish

---

## New Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                    OneMind OS Control Plane                      │
│                                                                 │
│  ┌─────────────┐  ┌──────────────┐  ┌────────────────────────┐ │
│  │  Heartbeat   │  │  Cron/Sched  │  │  NATS Event Bus        │ │
│  │  Loop        │  │  (persistent)│  │  (7 JetStream streams) │ │
│  │  (30m cycle) │  │              │  │                        │ │
│  └──────┬───────┘  └──────┬───────┘  └───────────┬────────────┘ │
│         │                 │                       │              │
│         ▼                 ▼                       ▼              │
│  ┌─────────────────────────────────────────────────────────────┐ │
│  │                     LEGACY (Agno Agent)                      │ │
│  │  • Full memory (user + agentic + session summaries)         │ │
│  │  • NATS tools (observe, publish, fan_out, system_status)    │ │
│  │  • All standard tools (search, code, browser, etc.)         │ │
│  │  • Team coordination (can delegate to any agent/team)       │ │
│  │  • Voice/Vision/Audio pipeline                              │ │
│  └──────────────────────────┬──────────────────────────────────┘ │
│                             │                                    │
│  ┌──────────┬───────────────┼───────────────┬──────────────┐    │
│  ▼          ▼               ▼               ▼              ▼    │
│ Agents    Teams         Nodes           Adapters        Skills   │
│ (Agno)   (Agno)     (WebSocket)     (OpenMind/etc)   (plugins) │
│                                                                  │
│  researcher  research-team   FarmBot     OpenMind OM1    ...    │
│  analyst     dev-team        Meshtastic  Custom bots            │
│  coder       ops-team        iOS/Android                        │
│  ...                         Camera/Mic                         │
└─────────────────────────────────────────────────────────────────┘
```

---

## Implementation Plan

### Phase 1: Legacy as Agno Agent

**File: `backend/agents/presets.py`** — Add Legacy preset:

```python
"legacy": {
    "name": "Legacy",
    "model_id": "gpt-4o",
    "tools": [
        # Standard tools
        "duckduckgo", "clickup", "shell", "file", "python",
        # NATS tools (new)
        "nats_observe", "nats_publish", "nats_fan_out",
        # System tools (new)
        "system_status", "system_overview", "entity_list",
    ],
    "instructions": [
        "You are Legacy — the super-orchestrator of OneMind OS.",
        "You are the user's direct presence in the system, like Jarvis.",
        "You see everything happening across all agents, teams, and devices.",
        "You can dispatch tasks to any agent or team.",
        "You have access to the full NATS event bus for real-time awareness.",
        "",
        "PROACTIVE BEHAVIOR:",
        "During heartbeat cycles, read HEARTBEAT.md and act on it.",
        "Surface urgent items. If nothing needs attention, say HEARTBEAT_OK.",
        "",
        "CAPABILITIES:",
        "- Voice, vision, audio (multimodal)",
        "- All agent tools (search, code, browse, etc.)",
        "- NATS bus (observe events, publish commands, fan out to teams)",
        "- System overview (see all entities, their status, capabilities)",
        "- Team delegation (route complex tasks to specialized teams)",
    ],
    "add_history_to_context": True,
    "num_history_runs": 20,  # Longer context for orchestrator
    "markdown": True,
    "enable_user_memories": True,
    "enable_agentic_memory": True,
    "enable_session_summaries": True,
    "enable_reasoning": True,
    "reasoning_max_steps": 15,
    "run_parallel_tools": True,
}
```

### Phase 2: Heartbeat Loop

**File: `backend/heartbeat.py`** — New proactive engine:

```python
"""
Heartbeat Loop — The Proactive Engine
=======================================

Inspired by OpenClaw's heartbeat pattern. Runs periodic agent turns
so Legacy can surface anything that needs attention without being asked.

How it works:
1. Every N minutes (configurable, default 30m)
2. Read HEARTBEAT.md checklist from workspace
3. Run a Legacy agent turn with the heartbeat prompt
4. If nothing needs attention: HEARTBEAT_OK (suppressed)
5. If urgent: deliver alert to configured channel (WebSocket, NATS)

Config (environment variables):
    HEARTBEAT_ENABLED=true
    HEARTBEAT_INTERVAL_MINUTES=30
    HEARTBEAT_ACTIVE_HOURS_START=08:00
    HEARTBEAT_ACTIVE_HOURS_END=24:00
"""

class HeartbeatLoop:
    def __init__(self, legacy_agent, interval_minutes=30):
        ...

    async def start(self):
        """Start the heartbeat loop as an asyncio task."""
        ...

    async def _tick(self):
        """Single heartbeat cycle."""
        # 1. Check active hours
        # 2. Read HEARTBEAT.md
        # 3. Run Legacy agent turn with heartbeat prompt
        # 4. Process response (HEARTBEAT_OK vs alert)
        # 5. Deliver to configured targets (WebSocket, NATS)
        ...
```

### Phase 3: Enhanced Scheduler (Cron)

**File: `backend/scheduler.py`** — Replace naive scheduler with persistent cron:

```python
"""
Cron Scheduler — Persistent Job Store
=======================================

Inspired by OpenClaw's cron system. Two execution modes:
- Main session: inject system event into next heartbeat
- Isolated: run dedicated agent turn in fresh session

Jobs persist to database (survive restarts).
"""

class CronJob:
    id: str
    name: str
    schedule: ScheduleConfig  # at, every, cron
    session_target: str       # "main" or "isolated"
    payload: PayloadConfig    # systemEvent or agentTurn
    delivery: DeliveryConfig  # announce, webhook, none
    enabled: bool
    agent_id: str = "legacy"  # which agent runs it

class CronScheduler:
    def __init__(self, db, agents_dict):
        self.jobs: Dict[str, CronJob] = {}
        ...

    async def start(self):
        """Load jobs from DB, start the scheduler loop."""
        ...

    async def add_job(self, job: CronJob) -> str:
        """Add a new cron job (persisted)."""
        ...

    async def _run_job(self, job: CronJob):
        """Execute a cron job."""
        ...
```

### Phase 4: NATS Tools for Legacy

**File: `backend/tools/nats_tools.py`** — Tools that give Legacy NATS superpowers:

```python
"""
NATS Tools — Give Legacy Eyes and Hands on the Event Bus
=========================================================

These tools let the Legacy agent:
1. Observe: See what's happening on NATS topics
2. Publish: Send commands/events to NATS
3. Fan-out: Broadcast to multiple targets
4. Status: Get system-wide health snapshot
"""

def nats_observe(topic: str, limit: int = 10) -> str:
    """Observe recent messages on a NATS topic."""
    ...

def nats_publish(topic: str, message: dict) -> str:
    """Publish a message to a NATS topic."""
    ...

def nats_fan_out(targets: list, command: str) -> str:
    """Fan out a command to multiple NATS targets."""
    ...

def system_overview() -> str:
    """Get current system overview: all agents, teams, devices, their status."""
    ...
```

### Phase 5: Slim Nexus

Refactor Nexus to keep only the unique value:

```
backend/nexus/
    ├── awareness.py     → KEEP (NATS-aware agents)
    ├── comms.py         → KEEP (NexusMessage multimodal format)
    ├── endpoints.py     → KEEP (overview + events + stream)
    ├── registry.py      → SLIM (FrameworkAdapter ABC only, remove EntityDescriptor bloat)
    ├── core.py          → REPLACE with slim startup function
    ├── legacy.py        → REMOVE (Legacy is now an Agno agent preset)
    └── dispatcher.py    → REMOVE (Agno handles coordination)
```

### Phase 6: HEARTBEAT.md

**File: `backend/workspace/HEARTBEAT.md`** — Legacy's proactive checklist:

```markdown
# Heartbeat Checklist

- Check for high-priority events on NATS bus
- Scan for any failed agent runs or errors
- Review pending tasks or approvals
- If daytime: light check-in if nothing else is pending
- If a scheduled workflow completed, summarize results
- Monitor system health (memory, CPU, disk)
```

### Phase 7: Voice/Vision Pipeline (future)

Based on ada_local's pattern:
```
Wake Word → STT (Whisper) → Legacy Agent → TTS (ElevenLabs/Piper)
              ↓
         Intent Router (Agno reasoning)
              ↓
         Tool Execution + Team Delegation
```

This will be implemented as a **Node** that connects via WebSocket, similar to OpenClaw's macOS/iOS nodes.

---

## Sidebar Consolidation

Combine COMMAND + SYSTEM + NEXUS + GAME into a single **COMMAND CENTER** section:

```
COMMAND CENTER (always expanded)
├── Chat (Legacy)          — main interaction
├── System Overview        — agents, teams, devices
├── Event Stream           — real-time NATS events
├── Heartbeat Status       — last heartbeat, next tick
├── Cron Jobs              — scheduled tasks
├── Game Dashboard         — XP, streaks, achievements
└── System Health          — CPU, memory, uptime

WORKFORCE
├── Agents
├── Teams
└── Models

WORKSPACE
├── Sessions
├── Memories
├── Knowledge
└── Documents

AUTOMATION
├── Workflows
├── Triggers
└── Webhooks

SETTINGS
├── Settings
└── Integrations
```

---

## Deployment Topology

```
LOCAL (dev)                   CLOUD (prod)              MOBILE
─────────────               ─────────────             ─────────
Docker Compose              DigitalOcean              Flutter App
├── postgres                ├── Managed Postgres       ├── WebSocket → Cloud
├── nats                    ├── NATS (container)       ├── Push notifications
├── backend                 ├── Backend (container)    └── Voice node
├── frontend                ├── Frontend (CDN)
└── openmind-bridge         └── openmind-bridge

All connected via NATS for real-time sync.
```

---

## Migration Path

1. **Now**: Add Legacy preset, heartbeat loop, NATS tools, sidebar consolidation
2. **Next**: Replace scheduler.py with persistent cron, slim Nexus
3. **Later**: Voice pipeline node, FarmBot/Meshtastic node adapters
4. **Future**: Multi-gateway federation, community-scale management

---

## Key Principles

1. **Agno-first**: Legacy is an Agno agent. All agent features (memory, sessions, tools, teams) come from Agno.
2. **NATS for awareness**: NATS is the nervous system. Everything that happens publishes to NATS. Legacy subscribes.
3. **Heartbeat for proactivity**: The heartbeat loop makes the system proactive. It runs whether or not anyone is chatting.
4. **Cron for precision**: Exact-time tasks use cron. Background monitoring uses heartbeat.
5. **Nodes for hardware**: Devices (FarmBot, Meshtastic, phones) connect as "nodes" via WebSocket/MQTT.
6. **Adapters for future**: The thin FrameworkAdapter keeps the door open for OpenMind, custom bots, etc.
