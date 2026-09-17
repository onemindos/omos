# Hybrid Mind Architecture — Digital Robotics for OneMind OS

> **The insight**: OpenClaw isn't just a chatbot — it's a bare-metal AI runtime
> orchestrating containerized services. This is the same pattern used in
> physical robotics (OM1). The only difference is **what the HAL controls** —
> digital actions (APIs, smart home, messages) vs physical actions (motors, grippers, legs).

---

## The Pattern: Bare Metal Orchestrator + Containerized Services

```
┌─────────────────────────────────────────────────────────────────────────┐
│                        DO DROPLET (One Machine)                        │
│                                                                        │
│  ┌──────────────────────────────────────────────────────────────────┐  │
│  │              BARE METAL — THE CORTEX                             │  │
│  │                                                                  │  │
│  │   OpenClaw Gateway (Node.js process)                             │  │
│  │   ├── Cortex Loop (heartbeat tick @ configurable Hz)             │  │
│  │   ├── Fuser (merge all inputs into unified context)              │  │
│  │   ├── LLM Decision Layer (Opus/Sonnet reasoning)                 │  │
│  │   ├── Action Dispatcher (execute decided actions)                │  │
│  │   ├── Memory (SQLite + hybrid BM25/vector search)                │  │
│  │   ├── Skills (npm packages — nexus-bridge, agno-dispatch, etc.)  │  │
│  │   ├── Codex (git-synced markdown vault @ /home/onemind/codex/)  │  │
│  │   └── Cron (heartbeat checks, git pull, scheduled tasks)         │  │
│  │                                                                  │  │
│  │   Listening on :18789 via Tailscale                              │  │
│  └────────────┬───────────────────────┬─────────────────────────────┘  │
│               │ localhost:4222        │ localhost:5432                  │
│  ┌────────────▼───────────────────────▼─────────────────────────────┐  │
│  │              DOCKER — THE NERVOUS SYSTEM                         │  │
│  │                                                                  │  │
│  │   NATS JetStream    — Event bus (:4222, :7422, :1883, :8222)     │  │
│  │   OmniDB            — PostgreSQL+TimescaleDB+pgvector (:5432)    │  │
│  │   NATS Ingester     — Event archival pipeline                    │  │
│  │   Legacy Bridge     — NATS → OpenClaw webhook relay              │  │
│  │                                                                  │  │
│  │   All on onemind-network (Docker bridge)                         │  │
│  └──────────────────────────────────────────────────────────────────┘  │
│                                                                        │
│  ┌──────────────────────────────────────────────────────────────────┐  │
│  │              EDGE CONNECTIONS — THE PERIPHERAL NERVES            │  │
│  │                                                                  │  │
│  │   Tailscale mesh → Mac (leaf node) → Home (leaf) → Office(leaf) │  │
│  │   MQTT :1883 → Home Assistant → Zigbee/Z-Wave/WiFi devices      │  │
│  │   WebSocket :9222 → Flutter UI → Browser dashboards              │  │
│  │   Leaf :7422 → ESP32 WireClaw nodes → Physical sensors          │  │
│  └──────────────────────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────────────────────┘
```

---

## Agno (Containerized) vs OpenClaw (Bare Metal) — Full Comparison

### What They Are

**Agno** (Python framework) runs as a **containerized service** inside Docker.
You define agent presets, a factory instantiates them, and FastAPI serves them
over HTTP. The agents are ephemeral — they exist inside the container's lifecycle.

**OpenClaw** (Node.js runtime) runs as a **bare-metal daemon** on the host OS.
It's a long-lived process that owns its filesystem, maintains persistent memory
in SQLite, runs npm-installed skills, and manages its own cron schedules.

### The Comparison

| Dimension | Agno (Containerized) | OpenClaw (Bare Metal) |
|-----------|---------------------|----------------------|
| **Runtime** | Python inside Docker | Node.js on host OS |
| **Lifecycle** | Container restart = cold start | nohup/systemd = persistent daemon |
| **Memory** | In-memory or external DB | Native SQLite + BM25/vector hybrid |
| **Filesystem** | Ephemeral (volumes for persistence) | Persistent — owns ~/.openclaw/ |
| **Skills/Tools** | Python functions registered in registry | npm packages installed to disk |
| **Identity** | Defined in presets.py config dict | SOUL.md, AGENTS.md, USER.md files |
| **Channels** | HTTP API endpoints | WhatsApp, Discord, Telegram, iMessage, Voice, Web |
| **Multi-agent** | Multiple agent instances in same container | Multi-agent config in openclaw.json |
| **Scaling** | Horizontal (more containers) | Vertical (one powerful instance) |
| **Updates** | Rebuild image, redeploy | `openclaw update`, restart daemon |
| **Access to host** | No — sandboxed in Docker network | Full — disk, cron, git, npm, network |
| **Offline** | Only if container is running | Always running — full offline capability |

### Pros: Agno Containerized

1. **Isolation** — Agents can't affect host system. Clean sandboxing.
2. **Reproducible** — Same Docker image = same behavior everywhere.
3. **Horizontal scaling** — Spin up 10 containers, load balance.
4. **Microservice pattern** — Each agent type gets its own container, own restart policy.
5. **Easy rollback** — Point to old image tag.
6. **Familiar DevOps** — Kubernetes, Helm, docker-compose tooling.

### Pros: OpenClaw Bare Metal

1. **Persistent memory** — SQLite memory survives restarts, no external DB needed.
2. **Filesystem ownership** — Skills install, Codex syncs, memory writes — all persistent.
3. **Channel multiplexing** — One process handles ALL channels (WhatsApp, Discord, voice, etc.).
4. **Sub-millisecond local calls** — No Docker network hop for memory queries.
5. **Cron & scheduling** — Native heartbeat, cron jobs, git pull sync.
6. **Personality persistence** — SOUL.md, AGENTS.md are stable files on disk.
7. **Offline-first** — Works completely disconnected from internet.
8. **Community ecosystem** — 5,700+ skills on ClawHub, plug and play.

### Why OpenClaw Wins for OneMind

The Agno pattern treats agents as **stateless workers** receiving tasks via HTTP.
The OpenClaw pattern treats agents as **persistent minds** that live on the machine.

For a human enhancement framework that needs to:
- Remember conversations across years (persistent SQLite memory)
- Respond on WhatsApp at 3am (always-on daemon)
- Pull Codex updates via git cron (filesystem access)
- Run skills that install npm dependencies (disk persistence)
- Operate offline when internet drops (no cloud dependency)

**OpenClaw bare metal is the clear winner.**

But the key insight is: **they're not either/or.** The best architecture is:
- **OpenClaw = the persistent mind (bare metal)**
- **Everything else = containerized services that the mind uses**

The mind doesn't live in a container. The mind uses containers as organs.

---

## What OM1 Teaches Us — The Cortex Loop

OM1's architecture is brilliant. Let me extract what matters for OneMind:

### OM1's Core Loop (ModeCortexRuntime)

```python
# OM1's cortex.py — simplified
async def _run_cortex_loop(self):
    while True:
        await self.sleep_ticker_provider.sleep(1 / self.config.hertz)  # e.g., 2Hz = every 0.5s
        await self._tick()

async def _tick(self):
    # 1. Flush completed action promises
    finished_promises = await self.action_orchestrator.flush_promises()

    # 2. FUSE all inputs into one prompt
    prompt = await self.fuser.fuse(self.config.agent_inputs, finished_promises)

    # 3. Check for mode transitions
    transition = await self.mode_manager.process_tick(last_input)
    if transition:
        self._pending_mode_transition = new_mode
        return

    # 4. LLM decides actions
    output = await self.config.cortex_llm.ask(prompt)

    # 5. Dispatch actions
    await self.action_orchestrator.promise(output.actions)
```

### OM1's Key Concepts to Adopt

| OM1 Concept | What It Is | OneMind Equivalent |
|-------------|-----------|-------------------|
| **Cortex Loop** | Continuous tick at configurable Hz | OpenClaw's heartbeat/cron — extend to continuous loop |
| **Hertz (tick rate)** | How often the mind processes (0.5-10 Hz) | OpenClaw HEARTBEAT.md runs every 30min — too slow for awareness |
| **Fuser** | Merges ALL inputs into one unified prompt | OpenClaw's context window — enhance with NATS event fusion |
| **Modes** | Hot-swappable operating configurations | OpenClaw multi-agent (legacy, hp, le, ge) — add mode switching |
| **NLDB** | Natural Language Data Bus — sensor→text | NATS subjects are this — sensor data arrives as JSON, fuse to NL |
| **HAL** | Hardware Abstraction Layer (LLM→robot actions) | OpenClaw skills ARE the HAL — LLM decides, skills execute |
| **Multi-LLM** | Fast Action + Core Cognition + Coach/Mentor | OpenClaw uses Opus for reasoning — add fast Haiku for reflex |
| **Hot Reload** | Change config, runtime reloads live | OpenClaw's config can reload — extend to mode transitions |

### The Hybrid Mind: OM1 Patterns + OpenClaw Runtime

```
┌─────────────────────────────────────────────────────────────────┐
│                    LEGACY — HYBRID DIGITAL ROBOT                │
│                                                                 │
│  ┌───────────────────────────────────────────────────────────┐  │
│  │                PERCEPTION LAYER (Inputs)                  │  │
│  │                                                           │  │
│  │  WhatsApp ──┐                                             │  │
│  │  Discord ───┤                                             │  │
│  │  Telegram ──┤    ┌──────────┐                             │  │
│  │  Voice ─────┤───→│  FUSER   │──→ Unified Context          │  │
│  │  iMessage ──┤    │ (merge)  │    "Zeus said 'unlock door' │  │
│  │  NATS ──────┤    └──────────┘     via WhatsApp. Front     │  │
│  │  MQTT/IoT ──┤                     door is locked. Camera  │  │
│  │  Sensors ───┤                     shows delivery person.  │  │
│  │  Camera ────┘                     Eden is in guard mode." │  │
│  └───────────────────────────────────────────────────────────┘  │
│                              │                                  │
│  ┌───────────────────────────▼───────────────────────────────┐  │
│  │              COGNITION LAYER (Decision)                   │  │
│  │                                                           │  │
│  │  ┌─────────────┐  ┌──────────────────┐  ┌─────────────┐  │  │
│  │  │ Fast Action  │  │  Core Cognition  │  │   Coach     │  │  │
│  │  │ (Haiku/local)│  │  (Opus/Sonnet)   │  │  (Mentor)   │  │  │
│  │  │              │  │                  │  │             │  │  │
│  │  │ <300ms       │  │ Deep reasoning   │  │ Periodic    │  │  │
│  │  │ Reflex       │  │ Planning         │  │ critique    │  │  │
│  │  │ Safety stops │  │ Context-aware    │  │ of agent    │  │  │
│  │  │ Quick acks   │  │ decisions        │  │ behavior    │  │  │
│  │  └─────────────┘  └──────────────────┘  └─────────────┘  │  │
│  └───────────────────────────┬───────────────────────────────┘  │
│                              │                                  │
│  ┌───────────────────────────▼───────────────────────────────┐  │
│  │              ACTION LAYER (Execution — HAL)               │  │
│  │                                                           │  │
│  │  Digital Actions:              Physical Actions:           │  │
│  │  ├── Reply to WhatsApp         ├── Unlock front door       │  │
│  │  ├── Create Todoist task       ├── Turn on porch light     │  │
│  │  ├── Send Discord alert        ├── Change thermostat       │  │
│  │  ├── Update CRM deal          ├── Arm/disarm security     │  │
│  │  ├── Schedule calendar         ├── Start robot patrol      │  │
│  │  ├── Query database           ├── Open garage door        │  │
│  │  └── Publish NATS event        └── Activate camera feed    │  │
│  └───────────────────────────────────────────────────────────┘  │
│                                                                 │
│  ┌───────────────────────────────────────────────────────────┐  │
│  │                    MEMORY LAYER                           │  │
│  │                                                           │  │
│  │  SQLite (OpenClaw native)  │  OmniDB (PostgreSQL)         │  │
│  │  ├── Conversation memory   │  ├── embeddings (pgvector)   │  │
│  │  ├── BM25 text search      │  ├── events (TimescaleDB)    │  │
│  │  ├── Vector search         │  ├── knowledge_chunks (RAG)  │  │
│  │  └── Codex docs index      │  ├── cortex_ticks (replay)   │  │
│  │                            │  └── daily aggregates        │  │
│  └───────────────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────────────┘
```

---

## The Cortex Heartbeat — From 30min to 2Hz

OpenClaw currently has a HEARTBEAT.md that runs every **30 minutes** as a cron job.
OM1 runs a cortex loop at **2 Hz** (every 0.5 seconds).

Neither extreme is right for a digital enhancement framework. Here's the design:

### Three Heartbeat Tiers

| Tier | Frequency | Purpose | Implementation |
|------|-----------|---------|----------------|
| **Reflex** | 2 Hz (0.5s) | Safety alerts, emergency responses, physical actions | NATS subscriber with instant callback |
| **Awareness** | 0.1 Hz (10s) | Sensor fusion, context updates, state monitoring | OpenClaw cron or background skill |
| **Reflection** | Every 30 min | Deep analysis, planning, memory consolidation | HEARTBEAT.md (existing) |

```
Reflex (2 Hz):
  "Smoke detector triggered" → IMMEDIATELY kill HVAC, alert Zeus, call 911
  "Front door opened" → Check if expected, notify if not
  "Baby crying detected" → Alert Yocelin, check camera

Awareness (10s):
  Fuse: NATS events + sensor states + active conversations
  Update: mode context, location awareness, energy levels
  Check: any pending approvals, overdue tasks, anomalies

Reflection (30 min):
  "Am I being helpful? What did I learn? What patterns do I see?"
  Memory consolidation → OmniDB embeddings
  Cortex tick logging → cortex_ticks table
```

---

## Modes — Operating Configurations

Inspired by OM1's mode system with hot-reload:

| Mode | Trigger | Config | Behavior |
|------|---------|--------|----------|
| **default** | Normal state | All channels, all skills | Full awareness, all domains |
| **guard** | Night / away | Security sensors prioritized | Camera monitoring, door alerts, minimize noise |
| **focus** | Zeus working | Reduce interruptions | Batch notifications, only urgent alerts pass through |
| **family** | Family home | Family-friendly, LE domain | Kid-safe, household management, family calendar |
| **emergency** | Safety trigger | All channels max priority | Bypass all filters, contact emergency services |
| **sleep** | Nighttime | Minimal processing | Log only, suppress notifications, track sleep data |
| **travel** | GPS away from home | Remote monitoring, mobile-first | Home security elevated, location-aware, travel assist |

Mode transitions happen via:
1. **Explicit command**: "Legacy, switch to guard mode"
2. **Sensor-triggered**: Motion sensor + night + nobody home → guard mode
3. **Time-based**: 11pm → sleep mode, 7am → default mode
4. **Location-based**: Zeus leaves home (Tailscale node disconnect) → guard mode

---

## The Eden Robot Scenario — Full Digital-Physical Hybrid

Your example of a robot walking to the door, unlocking it via smart home, and
having full sensor awareness — here's how it works with this architecture:

```
1. Camera detects person at door
   └→ NATS: physical.camera.front.detection {person: true, confidence: 0.94}

2. Reflex heartbeat (2Hz) catches it immediately
   └→ Fuser merges: camera feed + door lock state + expected visitors + time

3. Cognition decides:
   - Fast Action (Haiku): "Person at door, check if expected"
   - Core (Opus): Cross-references calendar, recent messages
   - Decision: "This is the Amazon delivery. Unlock door."

4. Action dispatch (parallel):
   ├── NATS → home.lock.front.unlock      (smart lock opens)
   ├── NATS → home.light.porch.on         (lights the porch)
   ├── NATS → home.camera.front.record    (record interaction)
   ├── NATS → physical.robot.eden.move    (Eden walks to door)
   ├── WhatsApp → Zeus: "Delivery at door, unlocked for pickup"
   └── OmniDB → log event with full context

5. Robot (Eden) via OM1 runtime:
   - Receives NATS command on physical.robot.eden.move
   - OM1 cortex processes: navigate to front door
   - ROS2/Zenoh HAL translates to motor commands
   - Eden's own sensors feed back state via NATS
   - Eden's camera confirms: package collected

6. Post-action:
   ├── NATS → home.lock.front.lock        (re-lock door)
   ├── NATS → home.light.porch.off        (after 5 min)
   └── OmniDB → cortex_ticks: full tick history for replay
```

The key: **Legacy (OpenClaw) is the brain. Eden (OM1) is a body. They share
the same nervous system (NATS). The brain makes decisions, the body executes
physical actions, and both feed perception back into the loop.**

---

## OmniDB — Why One Database Wins

Previously: TimescaleDB for time-series only.
Now: PostgreSQL 16 + TimescaleDB + pgvector = **OmniDB**.

### Three Extensions, One Instance

```
┌─────────────────────────────────────────────────────────────┐
│                    OmniDB (PostgreSQL 16)                    │
│                                                             │
│  ┌──────────────────┐  ┌──────────────┐  ┌──────────────┐  │
│  │   TimescaleDB    │  │   pgvector   │  │  PostgreSQL  │  │
│  │                  │  │              │  │   Core       │  │
│  │  events          │  │  embeddings  │  │  Full SQL    │  │
│  │  cortex_ticks    │  │  knowledge_  │  │  JSONB       │  │
│  │  daily_counts    │  │    chunks    │  │  Triggers    │  │
│  │  daily_biometrics│  │              │  │  Functions   │  │
│  │                  │  │  semantic_   │  │  Foreign     │  │
│  │  Compression     │  │    search()  │  │    keys      │  │
│  │  Cont. aggs      │  │  knowledge_  │  │  JSON ops    │  │
│  │  time_bucket()   │  │    search()  │  │              │  │
│  └──────────────────┘  └──────────────┘  └──────────────┘  │
│                                                             │
│  One connection string. One backup. One container.          │
│  One volume. One healthcheck. Zero coordination overhead.   │
└─────────────────────────────────────────────────────────────┘
```

### What Each Extension Gives You

| Extension | Tables | Superpower | Example Query |
|-----------|--------|------------|---------------|
| **TimescaleDB** | events, cortex_ticks | Time partitioning, compression, `time_bucket()` | "What happened every Tuesday at 3pm for the last year?" |
| **pgvector** | embeddings, knowledge_chunks | Cosine similarity, IVFFlat/HNSW indexes | "Find memories similar to 'that conversation about moving'" |
| **PostgreSQL** | All | Full SQL, JSONB, triggers, joins | "JOIN events + embeddings WHERE similarity > 0.8 AND time > last week" |

The killer feature: **you can JOIN time-series data with vector similarity in
one query**. No separate databases, no sync, no ETL.

```sql
-- "What events are semantically similar to 'security concern' from last week?"
SELECT e.timestamp, e.subject, e.data,
       1 - (emb.embedding <=> $1) as similarity
FROM events e
JOIN embeddings emb ON emb.metadata->>'event_subject' = e.subject
WHERE e.timestamp > NOW() - INTERVAL '7 days'
  AND 1 - (emb.embedding <=> $1) > 0.7
ORDER BY similarity DESC
LIMIT 10;
```

---

## Ecosystem Integrations from the Inventory

From the attached inventory, these are the highest-value integrations:

### Immediate (Phase 1 — This Sprint)

| Project | Why | Integration Path |
|---------|-----|-----------------|
| **WireClaw** | $5 ESP32 agents with NATS, 30ms latency | Direct NATS publish/subscribe — already compatible |
| **HomeAssistant skill** | Smart home control via REST + MQTT | MQTT bridge already in hub.conf on :1883 |
| **OpenClaw-voice** | Browser-based voice with local Whisper STT | OpenClaw has native voice — use this for web UI |
| **clawsec** | Security suite: NVD alerts, audit, integrity | Install as skill on OpenClaw gateway |

### Near-term (Phase 2 — Next Month)

| Project | Why | Integration Path |
|---------|-----|-----------------|
| **VisionClaw** | Smart glasses (Meta Ray-Ban) → agent eyes | Gemini Live + OpenClaw skill bridge |
| **openclaw-n8n-stack** | N8N for deterministic workflows | Docker-compose already supports — add n8n service |
| **Supermemory** | Auto-recall/capture for conversations | Skill install + OmniDB embeddings table |
| **Lobster** | Composable workflow pipelines | Install as shell — "GitHub Actions for OpenClaw" |

### Future (Phase 3 — OM1 Bridge)

| Project | Why | Integration Path |
|---------|-----|-----------------|
| **OM1** | Physical robot runtime (Eden) | NATS bridge: Legacy publishes physical.robot.eden.*, OM1 subscribes |
| **ExoClaw** | WASM-sandboxed agents via NATS | NATS subject pattern: exoclaw.channel.account.peer |
| **Virtual Protocol ACP** | Agent-to-agent commerce | ACP integration via @agentclientprotocol/sdk |
| **Network-AI** | 12-framework orchestrator | Bridge for specialized workloads (research, coding) |

---

## Summary: What Changed

| Before | After |
|--------|-------|
| TimescaleDB only (time-series) | **OmniDB** = PostgreSQL + TimescaleDB + pgvector (3-in-1) |
| No vector search in DB | pgvector embeddings + `semantic_search()` function |
| No cortex tick logging | `cortex_ticks` hypertable for agent state replay |
| No knowledge base | `knowledge_chunks` table for RAG from Codex |
| 30-min heartbeat only | 3-tier: Reflex (2Hz) + Awareness (10s) + Reflection (30min) |
| Agents as chatbots | Agents as **digital robots** with perception → cognition → action |
| No mode system | 7 modes with sensor/time/location triggers |
| Separate concerns | **Hybrid Mind** — digital and physical share one nervous system (NATS) |
