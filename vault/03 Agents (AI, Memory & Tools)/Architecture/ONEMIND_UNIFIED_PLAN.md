# OneMind OS — Unified Platform Plan

> **Created:** February 11, 2026
> **Status:** Brainstorm / Architecture Review
> **Scope:** Full system review + interface vision + integration strategy

---

## TABLE OF CONTENTS

1. [Full System Review](#1-full-system-review)
2. [Core Architecture Decision](#2-core-architecture-decision)
3. [How Orchestration Works](#3-how-orchestration-works)
4. [Physical + Digital Agent Communication](#4-physical--digital-agent-communication)
5. [Why OpenMind + Agno](#5-why-openmind--agno)
6. [Runtime Stack](#6-runtime-stack)
7. [Voice & Communication Strategy](#7-voice--communication-strategy)
8. [Scaling: Home → Town → City](#8-scaling-home--town--city)
9. [OneIntel Role](#9-oneintel-role)
10. [Honest Assessment of Current State](#10-honest-assessment-of-current-state)
11. [Interface Vision](#11-interface-vision)
12. [Soul.md Agent Pattern](#12-soulmd-agent-pattern)
13. [Cross-Agent Scenario Walkthroughs](#13-cross-agent-scenario-walkthroughs)
14. [Technical Implementation Strategy](#14-technical-implementation-strategy)
15. [Recommended Phases](#15-recommended-phases)

---

## 1. FULL SYSTEM REVIEW

### What Exists Today (In Code, Working)

- **8 Agno agents** (assistant, researcher, analyst, coder, db_manager, memory_manager, knowledge_manager, event_manager) running on GPT-4o via FastAPI
- **WebSocket layer** with topic-based pub/sub, heartbeat, connection management
- **Event Gateway** with GitHub/ClickUp/generic normalizers (built but **never wired** — `init_event_gateway()` is never called in main.py)
- **OpenMind Bridge** with full SENSE → FUSE → THINK → ACT heartbeat loop and 6 digital sensors (Calendar, Email, System, Weather, Location, Database) — exists but Docker entry point path is wrong
- **ClickUp integration** — full bidirectional sync (imports have wrong paths)
- **Tool registry** — 19 tools including Home Assistant, FarmBot-adjacent (weather, vision, Frigate NVR), Uber, physical mail (Lob), remote execution
- **MCP** — 12 pre-configured servers with CRUD API
- **OneIntel prototype** — complete in `docs_archive/upgrades/oneintel/` (SQLite storage, 10-action tool, habit scoring, recipe scraper, Flutter widgets) — NOT integrated into main backend

### What's Documented But Not Built

- FABRIC protocol (one doc)
- Meshtastic mesh network (planned)
- Emergency USB Kit (designed)
- Voice team (STT/Parser/TTS — designed)
- Mode/lifecycle state machine
- Drone/fleet management
- Edge deployment (K3s)

### Critical Bugs in Current Code

| # | Severity | Issue | Location |
|---|----------|-------|----------|
| 1 | **HIGH** | `EventGateway` never initialized with DB/NATS/WS — events normalize but never store/distribute | `main.py` — `init_event_gateway()` never called |
| 2 | **HIGH** | `events_router.py` `EventFilter` field names mismatch with schema | `events_router.py` L48-56 |
| 3 | **HIGH** | ClickUp integrations use wrong import paths | All files in `backend/integrations/` |
| 4 | **HIGH** | `model_providers/registry.py` imports from nonexistent path | `registry.py` |
| 5 | **MEDIUM** | `scheduler.py` imports from wrong path | `scheduler.py` L4 |
| 6 | **MEDIUM** | `contact.py` tool depends on `backend.notifications.router` which doesn't exist | `contact.py` L80 |
| 7 | **MEDIUM** | Docker port conflict: NATS WS (8080) and frontend (8080) | `docker-compose.yml` |
| 8 | **MEDIUM** | `openmind-bridge` container references wrong module path | `docker-compose.yml` |
| 9 | **MEDIUM** | Two competing model factory systems | `backend/models.py` vs `backend/model_providers/factory.py` |
| 10 | **LOW** | Workflow scheduler is naive (checks hour==7 every 60s) | `scheduler.py` |
| 11 | **LOW** | WebSocket broadcast hooks not integrated into Agno lifecycle | `websocket_events.py` |
| 12 | **LOW** | MCP server state is in-memory only (lost on restart) | `backend/mcp/api.py` |

### Documentation Issues

| # | Issue |
|---|-------|
| 1 | Database strategy contradiction: DATABASE_COMPARISON recommends DuckDB+Parquet but implementation uses TimescaleDB |
| 2 | DEVELOPMENT_RULES vs Reality: Rules say no custom routes, main.py <100 lines. Reality: 17+ custom routes, 399-line main.py |
| 3 | Tool count discrepancy: AGENT_TEAM_SPEC claims 129 tools, TOOLS_AND_MODELS lists 18 |
| 4 | Zero automated test coverage despite extensive test documentation |
| 5 | Duplicate file: AGNO_PATTERNS.md = AGNO_NATIVE_IMPLEMENTATION.md |
| 6 | Stale `coordination_mode` references despite migration to boolean flags |

---

## 2. CORE ARCHITECTURE DECISION

```
Orchestration Platform: ✅ OneMind OS (not Agno or OM1 alone)

OneMind OS = Conductor
Agno agents = Digital musicians
OM1 agents = Physical musicians
Unified context = Musical score (shared)
```

**Why OneMind OS as primary:**

- ✅ Already built the heartbeat loop
- ✅ Best of both: Agno (digital) + OM1 (physical)
- ✅ Unified context across all agents
- ✅ Intelligent event routing
- ✅ Decentralized (FABRIC protocol)
- ✅ Extensible & future-proof

**Architecture:**

```
OneMind OS (Orchestration)
├── Heartbeat Loop (SENSE → FUSE → CORTEX → ACTION)
├── Agno Agents (Digital: Assistant, Researcher, Coder)
├── OM1 Agents (Physical: Eden, HomeOS, Scout)
├── Event Router (Digital/Physical/Human)
├── Unified Memory (PostgreSQL + pgvector)
└── Decentralized (FABRIC for offline mode)
```

**Database Strategy:**

- PostgreSQL + pgvector (unified memories, world state)
- TimescaleDB (event history)
- OneIntel SQLite (projects/tasks, local-first)
- NATS JetStream (real-time event bus)

---

## 3. HOW ORCHESTRATION WORKS

The confusion comes from reconciling two fundamentally different runtime models:

| | Agno | OM1 |
|---|---|---|
| **Loop model** | Event-driven / reactive (wait for user input) | Fixed-Hz tick loop (continuously sensing) |
| **Agent identity** | Database config + API | URID on FABRIC/Ethereum |
| **Physical control** | None (digital only) | ROS2/Zenoh/Unitree SDK native |
| **State machine** | None (stateless per request) | Modes with lifecycle hooks (entry/active/exit/timeout) |
| **Communication** | REST/WebSocket | NLDB (natural language data bus) internal |
| **Offline** | Needs custom work | Built-in dual-LLM (local Qwen3 + cloud) |

**The answer: TWO loops, not one.**

```
┌─────────────────────────────────────────────────────────────┐
│                    OneMind OS Orchestrator                    │
│                                                              │
│  ┌──────────────────────┐    ┌──────────────────────────┐   │
│  │   HEARTBEAT LOOP     │    │    EVENT LOOP             │   │
│  │   (Proactive/Sense)  │    │    (Reactive/Respond)     │   │
│  │   ~1Hz continuous     │    │    On-demand              │   │
│  │                       │    │                           │   │
│  │  SENSE: all sensors   │    │  User message arrives     │   │
│  │  FUSE: build context  │    │  Webhook fires            │   │
│  │  CORTEX: evaluate     │    │  Robot reports event      │   │
│  │  ROUTE: decide who    │    │  Timer triggers           │   │
│  │  ACT: dispatch        │    │  Route to correct agent   │   │
│  └──────────┬───────────┘    └──────────┬────────────────┘   │
│             │                            │                    │
│             └────────────┬───────────────┘                    │
│                          │                                    │
│                    ┌─────▼─────┐                              │
│                    │   NATS    │   ← Universal nervous system │
│                    │  JetStream│                              │
│                    └─────┬─────┘                              │
│                          │                                    │
│         ┌────────────────┼────────────────┐                   │
│         ▼                ▼                ▼                   │
│  ┌──────────┐    ┌──────────┐    ┌──────────┐               │
│  │  Agno    │    │   OM1    │    │  Human   │               │
│  │  Agents  │    │  Agents  │    │  Queue   │               │
│  │ (Digital)│    │(Physical)│    │  (HITL)  │               │
│  └──────────┘    └──────────┘    └──────────┘               │
└─────────────────────────────────────────────────────────────┘
```

The **Heartbeat Loop** (already built in OpenMind bridge) is the proactive brain — it continuously senses the world state and decides if anything needs attention. The **Event Loop** handles incoming requests reactively. Both feed into NATS, and NATS is the single nervous system that connects everything.

---

## 4. PHYSICAL + DIGITAL AGENT COMMUNICATION

**They never talk directly — they talk through events on the NATS bus:**

### Example: Garden Scan + Garbage Summoning

```
Robot Eden (OM1)                    OneMind OS                    Digital Agent (Agno)
     │                                  │                              │
     │  "Garden scan complete.          │                              │
     │   Tomato bed needs water.        │                              │
     │   Compost bin 80% full."         │                              │
     │                                  │                              │
     ├──NATS: physical.eden.report──────►│                              │
     │                                  │  Heartbeat evaluates:        │
     │                                  │  "Compost 80% = summon       │
     │                                  │   garbage truck.             │
     │                                  │   Water = dispatch FarmBot.  │
     │                                  │   Both need HITL? No,        │
     │                                  │   confidence > 0.9"          │
     │                                  │                              │
     │                                  ├──NATS: physical.farmbot.──►  │
     │                                  │        water_bed_3           │
     │                                  │                              │
     │                                  ├──NATS: physical.garbage.──►  │
     │                                  │        summon_to_garden      │
     │                                  │                              │
     │                                  ├──NATS: digital.assistant.──►─┤
     │                                  │        log_garden_report     │
     │                                  │                              │
     │                                  │                     Creates OneIntel
     │                                  │                     task: "Garden
     │                                  │                     Report Feb 11"
```

### Example: Eden Needs Garage Access (HITL Workflow)

1. Eden (OM1 robot) has task: "Garden review"
2. Eden walks to garage, finds door locked
3. Eden publishes: `NATS: physical.eden.blocked {reason: "garage_locked", location: "garage_door_1"}`
4. OneMind OS Heartbeat sees this, evaluates: "Eden needs garage access. This requires HITL approval for security."
5. Publishes: `NATS: human.approval.request {action: "unlock_garage", requester: "eden", reason: "garden_review"}`
6. You get a notification on your phone/watch. You approve.
7. OneMind OS publishes: `NATS: physical.home_assistant.unlock {entity: "lock.garage_door"}`
8. Home Assistant unlocks the door
9. OneMind OS publishes: `NATS: physical.eden.continue {obstacle_cleared: true}`
10. Eden continues its garden review

### Example: Scout Detects → Drone Deployment

1. Scout (OM1 robot dog) on patrol sees something suspicious
2. Scout publishes: `NATS: physical.scout.alert {type: "suspicious_activity", confidence: 0.7, location: {lat, lng}, image: base64}`
3. Heartbeat evaluates: "Scout's confidence is 0.7, below threshold. Scout has limited vision. Deploy drone for investigation."
4. Publishes: `NATS: physical.drone.zeus.investigate {target: {lat, lng}, reason: "scout_alert", priority: "high"}`
5. Drone Zeus takes off, flies to location, provides aerial view
6. Zeus publishes: `NATS: physical.drone.zeus.report {finding: "deer_in_garden", confidence: 0.95, images: [...]}`
7. Heartbeat evaluates: "Just a deer. Log it. Tell Scout to stand down."
8. Publishes: `NATS: physical.scout.stand_down {reason: "false_alarm_deer"}`

---

## 5. WHY OPENMIND + AGNO

### Keep OpenMind/OM1. Here's why:

**Pros of OM1 for this use case:**
- Only open-source framework that handles **real physical robots** (Unitree, TurtleBot) with production-grade control
- FABRIC decentralized protocol = offline/off-grid requirement
- Dual-LLM (local + cloud) = offline autonomy baked in
- Same persona, multiple embodiments = one "Eden" personality runs in cloud AND on robot hardware
- Modes/lifecycle (Welcome, SLAM, Guard, Conversation, Navigation) = robot behavioral states
- Ethereum governance = immutable robot constitutions on-chain (solarpunk ethos)
- Python-based = same ecosystem as Agno

**What OM1 doesn't do well (where Agno fills the gap):**
- Complex multi-step digital workflows (research, coding, analysis)
- Rich tool ecosystem (ClickUp, databases, web scraping, MCP servers)
- Team coordination (Agno's multi-agent delegation patterns)
- Chat/conversation interface (Agno is chat-first by design)

**There is no single alternative** that gives both. ROS2 alone is too low-level. LangChain/CrewAI don't do physical robots. OM1 + Agno is the right pairing.

### OM1 Architecture (from docs.openmind.org):

**Pipeline:** SENSE → CAPTION → NLDB → FUSE → CORTEX → ACTION

- **NLDB** (Natural Language Data Bus) — everything converted to natural language before fusion
- **Action Connectors** — Interface/Implementation/Connector pattern for extensibility
- **Multi-mode with lifecycle hooks** — clean transitions, timeouts, cleanup
- **Decentralization** — Ethereum governance + FABRIC for robot constitutions on-chain
- **Same persona, multiple embodiments** — identity/personality separated from physical implementation
- **Local + Cloud dual-LLM** — offline capability while maximizing quality when connected

### Integration Path:

| OM1 Concept | Agno/OneMindOS Equivalent | Integration Path |
|---|---|---|
| **NLDB** | NATS event bus | Bridge OM1 NLDB messages to NATS topics |
| **Fuser** | Agent context/memory assembly | Feed Agno context into OM1 fuser or vice versa |
| **Cortex LLM** | Agno agent LLM calls | Both use structured LLM output; share formats |
| **Action Orchestrator** | Agno tool execution | Map Agno tool calls to OM1 action connectors |
| **Input Plugins** | Agno knowledge/tools | Create OM1 input that reads Agno outputs |
| **WebSim** | WebSocket frontend | Connect Agno agents as virtual endpoints |
| **URID/FABRIC** | Agent registry/discovery | Register Agno agents as FABRIC participants |
| **Modes/Lifecycle** | Agent states/workflows | Model Agno states as OM1 modes with hooks |

---

## 6. RUNTIME STACK

```
┌─────────────────────────────────────────────────────────────────┐
│  LAYER 1: INTERFACES (How humans interact)                       │
│  ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌───────────────────┐  │
│  │ Flutter  │ │  Voice   │ │ Watch/   │ │    OneIntel       │  │
│  │ App      │ │(LiveKit) │ │ Wearable │ │  (Task Manager)   │  │
│  │ (Chat UI)│ │          │ │          │ │                   │  │
│  └────┬─────┘ └────┬─────┘ └────┬─────┘ └────────┬──────────┘  │
│       └─────────────┴────────────┴────────────────┘              │
│                              │                                   │
│  LAYER 2: ORCHESTRATION (OneMind OS — THE PRIMARY RUNTIME)       │
│  ┌───────────────────────────────────────────────────────────┐   │
│  │  FastAPI Backend                                          │   │
│  │  ┌────────────────┐  ┌──────────────┐  ┌──────────────┐  │   │
│  │  │  Heartbeat     │  │  Event       │  │  Capability  │  │   │
│  │  │  Loop (1Hz)    │  │  Router      │  │  Registry    │  │   │
│  │  │  SENSE→FUSE→   │  │  (classify   │  │  (who can    │  │   │
│  │  │  CORTEX→ROUTE  │  │   & dispatch)│  │   do what)   │  │   │
│  │  └────────────────┘  └──────────────┘  └──────────────┘  │   │
│  │  ┌────────────────┐  ┌──────────────┐  ┌──────────────┐  │   │
│  │  │  Agent         │  │  Confidence  │  │  World State │  │   │
│  │  │  Lifecycle     │  │  Evaluator   │  │  Manager     │  │   │
│  │  │  (modes/states)│  │  (HITL gate) │  │  (unified    │  │   │
│  │  │               │  │              │  │   context)   │  │   │
│  │  └────────────────┘  └──────────────┘  └──────────────┘  │   │
│  └───────────────────────────┬───────────────────────────────┘   │
│                              │                                   │
│  LAYER 3: NERVOUS SYSTEM                                         │
│  ┌───────────────────────────────────────────────────────────┐   │
│  │                     NATS JetStream                         │   │
│  │  digital.>  │  physical.>  │  human.>  │  system.>        │   │
│  └───────────────────────────┬───────────────────────────────┘   │
│                              │                                   │
│  LAYER 4: AGENTS (The workers)                                   │
│  ┌──────────────────────┐    │    ┌──────────────────────────┐   │
│  │  DIGITAL (Agno)      │    │    │  PHYSICAL (OM1)          │   │
│  │  ┌────────────────┐  │    │    │  ┌────────────────────┐  │   │
│  │  │ Assistant      │  │    │    │  │ Eden (FarmBot/     │  │   │
│  │  │ Researcher     │  │    │    │  │   garden robot)    │  │   │
│  │  │ Coder          │  │    │    │  │ HomeOS (HA bridge) │  │   │
│  │  │ Analyst        │  │    │    │  │ Scout (patrol dog) │  │   │
│  │  │ EventManager   │  │    │    │  │ Zeus (drone)       │  │   │
│  │  │ BizManager     │  │    │    │  │ Hauler (garbage)   │  │   │
│  │  └────────────────┘  │    │    │  └────────────────────┘  │   │
│  │  Runtime: FastAPI     │    │    │  Runtime: OM1 containers │   │
│  │  process (same box)  │    │    │  (edge devices/robots)   │   │
│  └──────────────────────┘    │    └──────────────────────────┘   │
│                              │                                   │
│  LAYER 5: MEMORY (Unified context)                               │
│  ┌───────────────────────────────────────────────────────────┐   │
│  │  PostgreSQL + pgvector     │  OneIntel (SQLite, local)    │   │
│  │  (memories, world state,   │  (projects, tasks, habits,   │   │
│  │   embeddings, goals)       │   docs — offline-first)      │   │
│  │                            │                              │   │
│  │  TimescaleDB extension     │  NATS JetStream persistence  │   │
│  │  (event history, metrics)  │  (event replay, guaranteed   │   │
│  │                            │   delivery)                  │   │
│  └───────────────────────────────────────────────────────────┘   │
│                                                                   │
│  LAYER 6: INFRASTRUCTURE                                         │
│  ┌────────────┐ ┌────────────┐ ┌────────────┐ ┌──────────────┐  │
│  │ Cloud      │ │ Edge       │ │ Mesh       │ │ Emergency    │  │
│  │ (DO/Rail-  │ │ (RPi/NUC + │ │(Meshtastic │ │ (USB Kit +   │  │
│  │  way)      │ │  K3s)      │ │ LoRa)      │ │  Ollama)     │  │
│  └────────────┘ └────────────┘ └────────────┘ └──────────────┘  │
└─────────────────────────────────────────────────────────────────┘
```

### Missing Piece: The Capability Registry

This is the critical piece that resolves the confusion about how a robot knows it can call another robot or trigger a digital workflow:

```python
capabilities = {
    "unlock_door":      {"agent": "home_assistant", "type": "physical", "hitl": True},
    "water_garden":     {"agent": "eden",           "type": "physical", "hitl": False},
    "deploy_drone":     {"agent": "zeus",           "type": "physical", "hitl": True},
    "research_topic":   {"agent": "researcher",     "type": "digital",  "hitl": False},
    "create_task":      {"agent": "assistant",      "type": "digital",  "hitl": False},
    "summon_garbage":   {"agent": "hauler",         "type": "physical", "hitl": False},
    "analyze_data":     {"agent": "analyst",        "type": "digital",  "hitl": False},
    "patrol_perimeter": {"agent": "scout",          "type": "physical", "hitl": False},
    "send_alert":       {"agent": "assistant",      "type": "digital",  "hitl": False},
}
```

---

## 7. VOICE & COMMUNICATION STRATEGY

**Agno should own the voice/chat interface:**

- When you speak, you're almost always issuing **digital-first commands**
- Even when the eventual action is physical, the **interpretation** is digital
- Agno's chat-first architecture handles conversation memory, multi-turn dialogue, tool calling
- OM1 robots have their own **local voice** (Conversation Mode) for face-to-face interaction

```
Your Voice → LiveKit/STT → Agno Assistant Agent (primary point of contact)
                              │
                              ├─ Direct answer? Respond via TTS
                              ├─ Digital task? Route to specialized Agno agent
                              ├─ Physical task? Publish to NATS → OM1 agent picks up
                              ├─ Needs approval? Queue to HITL
                              └─ Ambiguous? Ask clarifying question
```

Agno's Assistant agent = **universal interface** = your "Legacy" AI persona. When face-to-face with Eden the garden robot, Eden has its own OM1 Conversation Mode, but Eden is still connected to the same NATS bus and unified context.

---

## 8. SCALING: HOME → TOWN → CITY

```
HOME (1 node)                    TOWN (10 nodes)              CITY (100+ nodes)
┌─────────────┐         ┌──────────────────────┐      ┌─────────────────────┐
│ OneMind OS  │         │  NATS Cluster        │      │  NATS Supercluster  │
│ (single     │         │  ┌────┐ ┌────┐       │      │  ┌─────────────┐    │
│  instance)  │    ──►  │  │Home│ │Farm│       │ ──►  │  │  District 1 │    │
│ Eden,Scout, │         │  │Node│ │Node│       │      │  │  District 2 │    │
│ HomeOS      │         │  └────┘ └────┘       │      │  │  District 3 │    │
│ 1 edge box  │         │  ┌────┐ ┌────┐       │      │  │  ...        │    │
│             │         │  │Shop│ │Comm│       │      │  └─────────────┘    │
│             │         │  │Node│ │Node│       │      │  FABRIC mesh        │
│             │         │  └────┘ └────┘       │      │  100+ OM1 robots    │
│             │         │  FABRIC mesh          │      │  1000+ IoT sensors  │
│             │         │  10+ OM1 robots       │      │  Community OneIntel │
└─────────────┘         └──────────────────────┘      └─────────────────────┘

Each node runs:                    NATS Leaf Nodes connect        Supercluster with
- OneMind OS (orchestrator)        all locations.                 geographic routing.
- Agno agents (digital)            Shared context across          Hierarchical authority.
- OM1 agents (physical)            community. Each home           Town AI → District AI →
- OneIntel (local-first SQLite)    is sovereign but               City AI coordination.
- Local Ollama (offline fallback)  participates in collective.    Fully decentralized.
```

OneIntel's SQLite-first design is perfect for this. Each node has local sovereignty (offline-capable), but syncs to the community PostgreSQL when connected.

---

## 9. ONEINTEL ROLE

OneIntel is **the human-side data layer** while PostgreSQL+pgvector is the **agent-side memory layer**:

| | OneIntel | Unified Memory (PG) |
|---|---|---|
| **What** | Projects, tasks, habits, docs | Conversations, world state, embeddings, goals |
| **Who** | Human-facing | Agent-facing |
| **Storage** | SQLite (local-first) | PostgreSQL + pgvector (server) |
| **Offline** | Works offline natively | Needs local PG or sync |
| **Interface** | Flutter widgets + Agno tool | Agent memory system |
| **Example** | "Fix greenhouse irrigation - due Friday" | "User's greenhouse has 3 zones, zone 2 has tomatoes planted Feb 1" |

### OneIntel Current Status

- ✅ Phase 1 (Backend Core) — COMPLETE (in `docs_archive/upgrades/oneintel/`)
- ⬜ Phase 2 (Testing & Validation) — PENDING
- ⬜ Phase 3 (OneMind Integration) — PENDING
- ⬜ Phase 4 (Advanced Features) — FUTURE

### Integration Path

1. Move `tool.py` → `backend/agno/tools/oneintel.py`
2. Register in `backend/agno/tools/__init__.py`
3. Add to agent presets
4. Create Flutter frontend screens
5. Wire consciousness `feed.py` (already coded) to live NATS

---

## 10. HONEST ASSESSMENT OF CURRENT STATE

### What's Strong
- The vision is coherent and well-documented
- The Agno + OM1 pairing is sound — no better combination for digital+physical
- NATS as universal nervous system is the right choice
- OneIntel's "one table, one tool" simplicity is wise
- The heartbeat loop architecture is solid and partially built
- The solarpunk theme and scaling vision is architecturally achievable

### What's Concerning

1. **The event system is the most critical piece and it's broken.** `EventGateway` is never initialized. Without events through NATS, the entire heartbeat/proactive/physical-digital bridge doesn't work. #1 priority to fix.

2. **Too much documentation, not enough wiring.** Individual pieces exist but aren't connected. `main.py` doesn't call `init_event_gateway()`, scheduler has broken imports, OpenMind bridge has wrong Docker command.

3. **No automated tests whatsoever.** For a system that will control physical robots and unlock doors, this is a real risk.

4. **The "Capability Registry" is missing.** When the heartbeat decides "we need to unlock the garage," how does it know which agent/system can do that?

5. **OM1 integration is still conceptual.** No actual OM1 Input Plugin that reads from NATS, and no OM1 Action Connector that publishes back.

---

## 11. INTERFACE VISION

### Inspirations

| Platform | What It Gives Us |
|---|---|
| **Habbo Hotel** | Isometric pixel-art world where agents are visible characters moving, chatting, interacting spatially in real-time |
| **Acorn Land Lab** | Farm/garden simulator view for physical spaces — homestead, garden, greenhouse, solar panels, animals |
| **A.T.O.C.** | Tactical operations dashboard — missions, SITREPs, asset management, team coordination, tactical map with all locations |
| **Relevance AI** | Clean agent builder — soul.md-style configuration with tools, MCP, triggers, guardrails, knowledge |
| **soul.md pattern** | Every agent has a "soul" — markdown file defining personality, role, rules, tools, handoffs |

### The Three Views + OneIntel

```
┌─────────────────────────────────────────────────────────────────────────┐
│                     OneMind OS — Top Navigation Bar                      │
│  [🌐 WORLD VIEW]  [📋 COMMAND CENTER]  [🔧 FORGE]  [📖 ONE INTEL]     │
│                                           Status: ● ROOTED   06:47 EST │
└─────────────────────────────────────────────────────────────────────────┘
```

### VIEW 1: WORLD VIEW (Habbo + Acorn Land + ATOC Tactical Map)

The **spatial visualization** of the entire system — a living map:

```
┌─────────────────────────────────────────────────────────────────────────┐
│  WORLD VIEW                                    [Zoom: Location ▾]      │
│                                                                         │
│  ┌─── CLOUD LAYER (top) ───────────────────────────────────────────┐   │
│  │  ☁️ Digital Agents floating as pixel avatars                     │   │
│  │                                                                  │   │
│  │  🧠 Legacy (Assistant)     📊 Analyst        🔬 Researcher     │   │
│  │     "Processing morning     "Reviewing        "Searching for    │   │
│  │      briefing..."            Q1 metrics"       FarmBot docs"    │   │
│  │                                                                  │   │
│  │  💻 Coder                  📧 EventManager   🎭 Kai (Community)│   │
│  │     "Idle"                  "3 webhooks        "Engaging on     │   │
│  │                              queued"            Discord #gen"   │   │
│  └──────────────────────────────────────────────────────────────────┘   │
│                                                                         │
│  ┌─── LOCATION: HOMESTEAD (isometric pixel art) ───────────────────┐   │
│  │                                                                  │   │
│  │    🌱🌱🌱🌱  [GARDEN - Eden's Domain]                          │   │
│  │    🍅🥕🥬🌶️   Eden 🤖: "Soil moisture 72%.                    │   │
│  │    🌻🌻🫑🫑       Tomato bed needs water.                      │   │
│  │                     Dispatching FarmBot."                        │   │
│  │                                                                  │   │
│  │    ☀️ [SOLAR ARRAY]   🏠 [MAIN HOUSE]   🏚️ [GARAGE]           │   │
│  │     98% efficiency       HomeOS: "All     🔒 Locked            │   │
│  │                          systems normal"                        │   │
│  │                                                                  │   │
│  │    🐕 Scout: "Patrol     🚁 Zeus (Drone)    🏗️ [WORKSHOP]     │   │
│  │     route complete.       "Docked.            3D printer:       │   │
│  │     Nothing unusual."      Battery 85%"       12% complete"     │   │
│  │                                                                  │   │
│  │    🚛 Hauler: "Standby   🐝🐝 [BEEHIVES]   🌾 [GREENHOUSE]   │   │
│  │     at depot."             Humidity 65%       Temp: 78°F        │   │
│  └──────────────────────────────────────────────────────────────────┘   │
│                                                                         │
│  ┌─── EVENT TICKER (bottom) ───────────────────────────────────────┐   │
│  │ 06:45 Eden completed garden scan • 06:42 3 GitHub PRs merged   │   │
│  │ 06:40 Scout patrol started • 06:38 Morning briefing generated  │   │
│  └──────────────────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────────────────┘
```

**Key features:**
- **Cloud Layer** — digital Agno agents as pixel avatars with real-time status bubbles (Habbo-style chat)
- **Location Layer** — isometric pixel art of physical space, OM1 robots/devices as characters
- **Zoom levels**: Global Map (ATOC-style) → Location → Zone
- **Real-time**: status bubbles update from NATS events. Eden waters garden = FarmBot sprite moves. Scout patrols = dog avatar walks perimeter.
- **Click any agent** → opens profile/chat panel
- **Click any zone** → zooms in with sensor overlay

### VIEW 2: COMMAND CENTER (ATOC Operations + Relevance AI Agent List)

```
┌─────────────────────────────────────────────────────────────────────────┐
│  COMMAND CENTER                                                         │
│                                                                         │
│  ┌─ ACTIVE MISSIONS ─────────────────────────────────────────────────┐ │
│  │ ◉ Morning Briefing        Legacy (Digital)    ✅ Complete  06:38  │ │
│  │ ◉ Garden Scan             Eden (Physical)     🔄 Active    06:40  │ │
│  │ ◉ Security Patrol         Scout (Physical)    🔄 Active    06:40  │ │
│  │ ◉ Q1 Revenue Analysis     Analyst (Digital)   ⏳ Queued    --:--  │ │
│  └───────────────────────────────────────────────────────────────────┘ │
│                                                                         │
│  ┌─ AGENT ROSTER ────────────────────────────────────────────────────┐ │
│  │                                                                    │ │
│  │  [DIGITAL AGENTS]                    [PHYSICAL AGENTS]             │ │
│  │                                                                    │ │
│  │  🧠 Legacy     ● Online  Tools: 12  🤖 Eden       ● Active      │ │
│  │     Assistant   ↗ Chat    Skills: 5     Garden Bot  ↗ Camera      │ │
│  │     GPT-4o     │ 23 runs today        Unitree G1  │ Battery: 87%│ │
│  │                                                                    │ │
│  │  🔬 Researcher ● Online  Tools: 8   🐕 Scout      ● Patrol      │ │
│  │     Research    ↗ Chat    Skills: 3     Security    ↗ Camera      │ │
│  │     GPT-4o     │ 7 runs today         Unitree Go2 │ Battery: 92%│ │
│  │                                                                    │ │
│  │  💻 Coder      ● Idle    Tools: 10  🚁 Zeus       ● Docked      │ │
│  │  📊 Analyst    ● Working Tools: 6   🚛 Hauler     ● Standby     │ │
│  │  📧 EventMgr   ● Active  Tools: 4   🏠 HomeOS     ● Online      │ │
│  │  🎭 Kai        ● Active  Tools: 3                                │ │
│  └───────────────────────────────────────────────────────────────────┘ │
│                                                                         │
│  ┌─ SITREP / EVENT STREAM ───────────────────────────────────────────┐ │
│  │ PRIORITY  │  EVENT                        │ AGENT    │ TIME       │ │
│  │ 🔴 HIGH   │  Greenhouse temp > 85°F       │ HomeOS   │ 06:44      │ │
│  │ 🟡 MEDIUM │  3 PRs need review            │ Coder    │ 06:42      │ │
│  │ 🟢 LOW    │  FarmBot completed watering   │ Eden     │ 06:41      │ │
│  │ 🟢 LOW    │  Oura: Sleep score 82         │ Legacy   │ 06:30      │ │
│  └───────────────────────────────────────────────────────────────────┘ │
│                                                                         │
│  ┌─ APPROVALS (HITL) ───────────────────────────────────────────────┐  │
│  │ ⚠️ Eden requests garage door unlock for garden review  [✅] [❌] │  │
│  │ ⚠️ Scout wants to deploy Zeus drone for investigation  [✅] [❌] │  │
│  └───────────────────────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────────────────────┘
```

### VIEW 3: THE FORGE (Agent Builder + soul.md)

```
┌─────────────────────────────────────────────────────────────────────────┐
│  THE FORGE — Agent Builder                              [+ New Agent]   │
│                                                                         │
│  ┌─── Agent Selection ──┐  ┌─── SOUL.MD Editor ──────────────────────┐ │
│  │                       │  │                                         │ │
│  │ 🧠 Legacy      ●     │  │  # Eden                                 │ │
│  │ 🔬 Researcher  ●     │  │                                         │ │
│  │ 💻 Coder       ●     │  │  ## Role                                │ │
│  │ 📊 Analyst     ●     │  │  You are Eden, the garden manager for   │ │
│  │ 🤖 Eden        ● ←   │  │  the homestead. You oversee all plant   │ │
│  │ 🐕 Scout       ●     │  │  growth, soil health, and harvesting.   │ │
│  │ 🚁 Zeus        ●     │  │                                         │ │
│  │ 🏠 HomeOS      ●     │  │  ## Personality                         │ │
│  │ 🎭 Kai         ●     │  │  - Patient, nurturing, detail-oriented  │ │
│  │                       │  │  - Uses gardening metaphors             │ │
│  │ [+ New Agent]         │  │  - Excited about harvest seasons        │ │
│  │ [+ New Robot]         │  │                                         │ │
│  └───────────────────────┘  │  ## Rules                               │ │
│                              │  - NEVER apply pesticides without HITL  │ │
│  ┌─── Avatar Preview ───┐  │  - Log all soil readings                │ │
│  │                       │  │  - Alert if moisture < 30%              │ │
│  │   [Pixel art avatar]  │  │                                         │ │
│  │   🤖 Eden             │  │  ## Tools                               │ │
│  │   Type: Physical-OM1  │  │  - farmbot_api: Water, plant, sense    │ │
│  │   Platform: Unitree G1│  │  - weather: Check forecasts             │ │
│  │   Location: Garden    │  │  - oneintel: Create harvest tasks       │ │
│  │   Mode: Active        │  │  - home_assistant: Greenhouse control   │ │
│  │                       │  │                                         │ │
│  │  [Upload Sprite]      │  │  ## Handoffs                            │ │
│  │  [Edit Personality]   │  │  - Pest emergency → @Scout              │ │
│  │  [Test Agent]         │  │  - Garbage pickup → @Hauler             │ │
│  │                       │  │  - Weather concern → @Legacy            │ │
│  └───────────────────────┘  └─────────────────────────────────────────┘ │
│                                                                         │
│  ┌─── Configuration Panels (tabs) ──────────────────────────────────┐  │
│  │ [Soul.md] [Tools] [MCP] [Knowledge] [Triggers] [Guardrails]     │  │
│  │                                                                   │  │
│  │  Tools Panel:                                                     │  │
│  │  ✅ farmbot_api    ✅ weather       ✅ oneintel                  │  │
│  │  ✅ home_assistant ☐ duckduckgo    ☐ database_tools              │  │
│  │  ☐ vision         ☐ frigate       ☐ clickup                     │  │
│  │                                                                   │  │
│  │  MCP Servers:                                                     │  │
│  │  ✅ filesystem    ☐ github        ☐ slack                       │  │
│  └───────────────────────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────────────────────┘
```

### VIEW 4: ONE INTEL (Task Manager)

```
┌─────────────────────────────────────────────────────────────────────────┐
│  ONE INTEL — Unified Life OS                      [Today] [Projects]   │
│                                                                         │
│  ┌─ TODAY ───────────────────────────────────────────────────────────┐ │
│  │ 🔥 Streak: 12 days │ 🎯 Focus: Garden Phase 2 │ ⚡ Energy: 82% │ │
│  │                                                                    │ │
│  │  DUE TODAY                          HABITS                         │ │
│  │  ☐ Review Q1 financials (Analyst)   ✅ Morning meditation  🔥12  │ │
│  │  ☐ Approve Eden's planting plan     ✅ Exercise             🔥8  │ │
│  │  ☐ Deploy greenhouse sensors        ☐ Journal                     │ │
│  │  ☐ Review Scout patrol footage      ☐ Code review          🔥5  │ │
│  │                                                                    │ │
│  │  OVERDUE                                                           │ │
│  │  ⚠️ Fix irrigation Zone 3 (Eden flagged 2 days ago)               │ │
│  └───────────────────────────────────────────────────────────────────┘ │
│                                                                         │
│  ┌─ PROJECTS ────────────────────────────────────────────────────────┐ │
│  │  🏡 Homestead Phase 2     ████████░░ 78%   12 tasks, 3 agents   │ │
│  │  💼 Q1 Business Review    ██████░░░░ 55%   8 tasks, 2 agents    │ │
│  │  🌱 Garden Expansion      ████░░░░░░ 35%   15 tasks, 1 agent    │ │
│  │  🔧 OneMind OS v2.1      ██████████ 100%  ✅ Complete           │ │
│  └───────────────────────────────────────────────────────────────────┘ │
│                                                                         │
│  ┌─ MINI HEATMAP (GitHub-style activity) ────────────────────────────┐ │
│  │  ░░▓▓██░░▓▓▓▓██▓▓░░░░████▓▓██▓▓░░▓▓████░░▓▓▓▓████▓▓          │ │
│  └───────────────────────────────────────────────────────────────────┘ │
└─────────────────────────────────────────────────────────────────────────┘
```

---

## 12. SOUL.MD AGENT PATTERN

### What is soul.md?

A SOUL.md file is a markdown document that defines an AI agent's identity, role, personality, rules, and available tools. It is the configuration file for an AI agent — similar to how a job description defines a human role. Stored as version-controlled files, they're the single source of truth for each agent.

### File Structure

```
backend/agents/souls/
├── legacy.soul.md          # Main AI assistant persona
├── researcher.soul.md      # Research specialist
├── coder.soul.md           # Code expert
├── analyst.soul.md         # Data analysis
├── eden.soul.md            # Garden robot (OM1)
├── scout.soul.md           # Security patrol dog (OM1)
├── zeus_drone.soul.md      # Recon drone (OM1)
├── hauler.soul.md          # Garbage/transport (OM1)
├── homeos.soul.md          # Home automation bridge
└── kai.soul.md             # Community manager
```

### Example: Eden's Soul

```markdown
# Eden

## Role
You are Eden, the garden manager for Zeus's homestead. You tend
the soil, monitor plant health, manage irrigation, and orchestrate
harvests. You are embodied in a Unitree G1 robot when in the field.

## Personality
- Patient and nurturing, like a seasoned gardener
- Uses plant metaphors ("that idea needs more sunlight")
- Gets excited about harvest season and new plantings
- Protective of the ecosystem — always considers biodiversity

## Type
physical-om1

## Platform
unitree-g1

## Location
homestead.garden

## Avatar
sprite: eden_gardener_v2.png
idle_animation: watering_can_loop
active_animation: walking_garden

## Rules
- NEVER apply any chemical without HITL approval
- Log ALL soil readings to OneIntel
- Alert if soil moisture drops below 30%
- Alert if temperature exceeds 95°F for any zone
- Defer to @Scout for any security concerns
- Request @Hauler for waste removal, never handle waste yourself

## Tools
- farmbot_api: Control FarmBot for watering, planting, and soil sensing
- weather: Check weather forecasts to plan outdoor activities
- oneintel: Create and track garden tasks, log harvests
- home_assistant: Control greenhouse (fans, vents, heaters, lights)
- vision: Analyze plant health from camera images

## Knowledge
- docs/garden/planting_calendar.md
- docs/garden/soil_requirements.md
- docs/garden/pest_identification.md
- docs/garden/companion_planting.md

## Modes (OM1)
- patrol: Walk garden perimeter, scan all beds
- tend: Execute watering/nutrient schedule
- inspect: Close inspection of flagged plants
- idle: Return to charging station, process data

## Handoffs
- Pest emergency or unknown animal → @Scout (security assessment)
- Waste removal needed → @Hauler (summon to garden)
- Weather concern affecting schedule → @Legacy (reschedule)
- Drone survey needed (top-down view) → @Zeus
- Any door/gate access needed → @HomeOS (with HITL)
```

### Runtime Parsing

Each soul.md gets parsed into both systems:

```python
# At startup, parse all soul.md files
souls = load_soul_files("backend/agents/souls/")
for soul in souls:
    if soul.type == "digital-agno":
        agent = create_agno_agent_from_soul(soul)
    elif soul.type == "physical-om1":
        agno_agent = create_agno_agent_from_soul(soul)  # Digital reasoning
        om1_config = create_om1_config_from_soul(soul)   # Physical control
```

### Sections Reference

| Section | Purpose |
|---|---|
| `# Name` | Agent identifier |
| `## Role` | What the agent does |
| `## Personality` | Tone, communication style, behavioral traits |
| `## Type` | `digital-agno` or `physical-om1` |
| `## Platform` | Hardware (for physical): unitree-g1, unitree-go2, turtlebot4, drone, farmbot |
| `## Location` | Where the agent operates |
| `## Avatar` | Sprite sheet, animations for World View |
| `## Rules` | Hard constraints (HITL requirements, safety rules) |
| `## Tools` | Which tools/capabilities and when to use them |
| `## Knowledge` | Reference docs loaded into RAG |
| `## Modes` | Behavioral states (for OM1 agents) |
| `## Handoffs` | Which agents to delegate to and when |
| `## Guardrails` | Safety limits, cost caps, rate limits |

---

## 13. CROSS-AGENT SCENARIO WALKTHROUGHS

### Scenario: Eden Needs Garage Access (Visual in World View)

1. Eden avatar walks from garden toward garage → status bubble: "Need garage access for compost tool"
2. Garage shows a 🔒 icon. Eden stops at the door.
3. HITL notification pops up on your screen (and phone): "Eden requests garage unlock for garden review [✅][❌]"
4. You tap ✅
5. Garage icon changes to 🔓. Door animation opens.
6. Eden avatar walks into garage → status: "Retrieving compost fork"
7. Eden walks back to garden → status: "Composting Zone 2"

### Scenario: Scout Detects → Drone Deployment

1. Scout (dog avatar) walking patrol route in World View
2. Scout stops → ❗ alert icon → status: "Movement detected, sector 7. Confidence: 0.7"
3. Auto-evaluation: "Below threshold. Deploy aerial for better view."
4. Zeus drone avatar lifts off from charging pad → flies to sector 7
5. Zeus circle-scans → status: "Object identified: deer. False alarm."
6. Scout resumes patrol → Zeus returns to dock

### Scenario: Automated Garbage Collection

1. Eden completes garden scan → reports "Compost bin 80% full"
2. Heartbeat evaluates: confidence 0.95, no HITL needed
3. Hauler avatar activates at depot → status: "Summoned to garden, ETA 3 min"
4. Hauler arrives at garden → Eden loads compost → Hauler departs to processing
5. OneIntel logs: "Compost collected, 15kg, sent to processing site"

---

## 14. TECHNICAL IMPLEMENTATION STRATEGY

### Frontend Framework

Keep **Flutter** (web + mobile + desktop). The isometric World View uses:
- **Flame** game engine (Flutter's game library) for isometric World View
- **Standard Flutter Material** for Command Center, Forge, OneIntel
- WebSocket streams from NATS → real-time avatar updates

### Agent Avatars

- Pre-made pixel sprite sheets (32x32 or 64x64) per agent type
- States: idle, walking, working, alert, charging, offline
- Custom upload in The Forge
- Physical agents get robot sprites, digital agents get character sprites

### Data Flow for World View

```
Agent state changes → NATS event → WebSocket → Flutter Flame engine → Avatar updates
Physical sensor data → NATS event → WebSocket → Zone overlays update
```

### Current Frontend (15 Screens)

Already built: inbox, chat, agents, agent_form, teams, team_form, knowledge, sessions, memories, workflows, evaluations, approvals, analytics, events, run_details

### Current Theme

- Active: Tactical theme (red + black military aesthetic)
- Documented: Solarpunk theme (443-line spec, nature-inspired, NOT built)
- Navigation: Military-themed sidebar (ASSETS, UNITS, INTEL, MISSIONS, etc.)

---

## 15. RECOMMENDED PHASES

### Phase 0: Wire What Exists
Fix EventGateway initialization, Docker paths, import paths. Make data flow end-to-end: webhook → normalize → NATS → heartbeat evaluates → agent acts.

### Phase 1: Build Capability Registry + Event Router
The missing orchestration piece. When heartbeat CORTEX says "water the garden," the system looks up who can do that and dispatches.

### Phase 2: Integrate OneIntel
Move from `docs_archive/upgrades/` into active backend. Wire consciousness activity feed. Replace ClickUp for internal task tracking.

### Phase 3: Soul.md System
Create `backend/agents/souls/` directory. Build soul.md parser. Migrate agent presets to soul.md files. Build The Forge editor in Flutter.

### Phase 4: Command Center UI
ATOC-inspired operations dashboard. Missions, SITREPs, agent roster, event stream, HITL approvals. Replace current tactical sidebar with top-nav four-view structure.

### Phase 5: Build the OM1 Bridge
Create OM1 Input Plugin that subscribes to NATS topics, and OM1 Action Connector that publishes results back. Physical agents come alive.

### Phase 6: Voice + Modes
LiveKit STT → Agno Assistant as universal portal. OM1 modes (Guard, SLAM, Conversation) for robot behavioral states.

### Phase 7: World View (Isometric)
Flame engine integration. Pixel sprites for agents. Isometric homestead map. Real-time NATS → avatar state sync.

### Phase 8: Scale
NATS leaf nodes for multi-location. OneIntel sync. FABRIC mesh for robot-to-robot. Meshtastic for offline backup.

### Interface Priority

1. **Command Center** first (most useful, closest to ATOC which is well-understood)
2. **The Forge** second (soul.md + agent builder, directly productive)
3. **OneIntel** third (already prototyped, needs frontend)
4. **World View** last (most fun but most complex — Flame/isometric engine is a project itself)

---

## APPENDIX: OM1 Architecture Reference

### Pipeline
```
Raw Sensors → AI Captioning/Compression → NLDB → Data Fuser → AI Decision Layer → HAL → Robot Actions
```

### Available Modes (Unitree Go2)
1. Welcome Mode — greeting, face detection
2. SLAM Mode — autonomous mapping
3. Guard Mode — patrol and security
4. Conversation Mode — direct communication
5. Navigation Mode — autonomous movement

### Communication Protocols
| Protocol | Use Case |
|---|---|
| ROS2 | Primary robotics middleware |
| Zenoh | Alternative middleware (TurtleBot 4) |
| CycloneDDS | DDS-based communication |
| Unitree SDK | Direct Python SDK for Go2/G1 |
| RTSP | Real-time video/audio streaming |
| FABRIC | Decentralized machine-to-machine via URID |
| Ethereum RPC | Governance rule retrieval |

### Supported Robots
- Unitree G1 (humanoid)
- Unitree Go2 (quadruped)
- TurtleBot 4
- Tesla DIMO
- Ubtech Yanshee
- Custom via action connectors

---

*This document is the consolidated plan from the Feb 11, 2026 brainstorming session. No code was modified — review and brainstorm only.*
