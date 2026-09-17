# OneMind OS — Skills, Integration Patterns & Architecture Comparison

> How Legacy connects to the world, and why OpenClaw bare metal replaced Agno containerized as the brain.

---

## Part 1: How OpenClaw Skills Work

### What Is a Skill?

A skill is a **folder** installed into OpenClaw's workspace that gives an agent new abilities. Each skill has up to three components:

```text
skills/nexus-bridge/
├── SKILL.md       ← Markdown instructions injected into LLM system prompt
├── index.js       ← JavaScript code that registers callable tool functions
└── package.json   ← npm dependencies (nats, axios, etc.)
```

### Two Things Happen When a Skill Loads

**1. SKILL.md → Injected into the prompt**

The markdown content becomes part of the agent's system instructions. The LLM reads it like documentation — it learns when to use the skill, what parameters to pass, and what output to expect.

Example from our `nexus-bridge` skill:

```text
# SKILL.md — nexus-bridge

## What This Skill Does
Connects Legacy to the OneMind OS NATS nervous system. Use this to:
- Dispatch tasks to Agno digital agents (researcher, coder, analyst, etc.)
- Send commands to OpenMind OM1 physical robots
- Query system health across all runtimes

## When To Use It
Use nexus_dispatch when you need a specialized agent to handle a task:
- Research tasks → dispatch to "researcher"
- Code tasks → dispatch to "coder"
- Physical actions → dispatch to OM1 via robot commands
```

**2. index.js → Registers callable tools**

The JavaScript file exports functions that OpenClaw registers as tools the LLM can invoke. When the LLM decides to call `nexus_dispatch({agent: "researcher", task: "..."})`, OpenClaw executes the actual JavaScript:

```javascript
// nexus-bridge/index.js (simplified)
async function nexus_dispatch({ agent, robot_id, task, priority }) {
  const nc = await getNats();

  if (robot_id) {
    // Physical dispatch → OM1 robot via NATS
    await nc.publish(`physical.robot.${robot_id}.command`, encode(payload));
    return { status: "dispatched", target: "robot" };
  }

  if (agent) {
    // Digital dispatch → Agno agent via NATS or HTTP
    await nc.publish(`digital.agent.${agent}.task`, encode(payload));
    return { status: "dispatched", target: "agent" };
  }
}
```

### The Full Flow

```text
1. User sends WhatsApp message: "Research the best solar panels for the farm"

2. OpenClaw receives message via WhatsApp channel

3. LLM sees the SKILL.md instructions in its system prompt:
   "Use nexus_dispatch when you need a specialized agent..."

4. LLM decides: this is a research task → call nexus_dispatch

5. LLM generates tool call:
   nexus_dispatch({ agent: "researcher", task: "Research best solar panels..." })

6. OpenClaw executes the JavaScript function

7. Function publishes to NATS: digital.agent.researcher.task

8. Agno's researcher agent processes the task

9. Result flows back: Agno → NATS → bridge → OpenClaw → WhatsApp → Zeus
```

### Why So Many Skills? (5,700+ on ClawHub)

Same reason npm has 2 million packages — community-contributed, hyper-specific:

| Skill Type | Examples | Has Code? |
| --- | --- | --- |
| **Instruction-only** | daily-briefing, writing, problem-solving, code-review | No — just SKILL.md |
| **API integration** | github-api, home-assistant, clickup, obsidian | Yes — index.js calls APIs |
| **System bridge** | nexus-bridge (NATS), agno-dispatch | Yes — publishes to event bus |
| **Security** | clawsec (NVD alerts, integrity checks) | Yes — runs audits |
| **Hardware** | wireclaw (ESP32), visionclaw (smart glasses) | Yes — hardware protocols |

Most skills are **just a SKILL.md** — pure prompt engineering packaged as a module. The LLM reads instructions and behaves differently. Only integration skills have executable code.

**You don't load all 5,700.** Each agent in `openclaw.json` specifies which skills it gets. Legacy might load 15-20 skills.

---

## Part 2: Four Integration Patterns

There are exactly four ways to connect something to OneMind OS. Each has a specific purpose:

### The Decision Tree

```text
Does the AI need to DECIDE when to use it?
├── YES → SKILL
│
└── NO → Is it event-driven / real-time?
    ├── YES → NATS SUBJECT
    │
    └── NO → Is it deterministic / background?
        ├── YES → DIRECT INTEGRATION
        │
        └── Does it modify the OpenClaw runtime itself?
            └── YES → PLUGIN / HOOK
```

### Pattern 1: Skill — LLM Decides When To Call

```text
LLM reads SKILL.md → LLM decides "I should unlock the door" → calls tool → JS executes
```

**Use when**: The AI needs judgment about when to act. Context-dependent decisions.

**How it works**: SKILL.md tells the LLM what the tool does. index.js provides the executable function. The LLM chooses when to invoke it based on conversation context.

**OneMind examples**:

| Skill | What It Does | Why Skill (Not NATS) |
| --- | --- | --- |
| home-assistant | Control lights, locks, thermostat | LLM decides WHEN to unlock door |
| nexus-bridge | Dispatch to Agno agents / OM1 robots | LLM decides WHICH agent to use |
| github-api | Create issues, PRs, manage repos | LLM decides WHEN to create issue |
| obsidian | Read/write vault files | LLM decides WHICH file to read |
| todoist | Create/complete tasks | LLM decides task priority/timing |
| agno-dispatch | Call Agno HTTP API directly | LLM picks agent and frames task |

### Pattern 2: NATS Subject — Event-Driven, No LLM Trigger

```text
Sensor publishes event → NATS delivers → Subscriber processes automatically
```

**Use when**: Things happen regardless of what the LLM is doing. Telemetry, sensor data, heartbeats, cross-runtime communication.

**How it works**: Producers publish to NATS subjects. Consumers subscribe. No AI judgment involved in the delivery — though the data may eventually reach the LLM via the Fuser.

**OneMind examples**:

| Subject Pattern | What Flows | Publisher → Subscriber |
| --- | --- | --- |
| `physical.camera.*.motion` | Motion detection events | Camera → Legacy (via bridge) |
| `physical.door.*.state` | Door open/close/lock state | Sensor → OmniDB (via ingester) |
| `physical.robot.*.status` | Robot battery, position, mode | OM1 → Legacy |
| `digital.agent.*.result.*` | Agent task completions | Agno → Legacy (via bridge) |
| `nexus.heartbeat.>` | Runtime health pings | All runtimes → Monitoring |
| `comms.text.human.zeus` | Zeus's inbound messages | OpenClaw → All subscribers |

### Pattern 3: Direct Integration — No AI, No Events

```text
Python/JS code → HTTP call → External service (deterministic, scheduled)
```

**Use when**: Deterministic operations. No AI judgment needed. Background jobs, scheduled tasks, data pipelines.

**How it works**: A script runs on a schedule or trigger. It directly calls APIs or processes data. No LLM, no NATS — just code.

**OneMind examples**:

| Integration | What It Does | Trigger |
| --- | --- | --- |
| nats_ingester.py | Archives NATS events → OmniDB | Continuous (NATS subscriber) |
| git cron | `git pull` Codex vault every 5 min | Cron job |
| Health checks | Ping services, log status | Cron / systemd timer |
| Log rotation | Clean old logs, compress | Cron job |
| Cert renewal | Let's Encrypt renewal | Certbot timer |

### Pattern 4: Plugin / Hook — Extends the OpenClaw Runtime

```text
Hooks into OpenClaw lifecycle: onMessage, onBoot, onHeartbeat, etc.
```

**Use when**: You need to modify how OpenClaw **behaves as a runtime**, not what the agent can do. New channels, message routing, storage backends.

**OneMind examples**:

| Plugin | What It Modifies | Why Not a Skill |
| --- | --- | --- |
| Meta glasses channel | New input source for OpenClaw | It's a transport, not a tool |
| Mode switching | Changes which skills/prompts load | Modifies runtime config |
| Custom message router | Routes messages to specific agents | Changes routing logic |
| Storage backend hook | Where memory gets persisted | Infrastructure, not agent behavior |

### How They Work Together

```text
Camera motion fires → NATS subject (event-driven)
   → bridge.py forwards to OpenClaw (direct integration)
   → Reflex heartbeat catches it (plugin/hook)
   → Fuser merges camera + door state + calendar
   → LLM uses Home Assistant skill to unlock door (skill)
   → LLM uses Discord skill to notify Zeus (skill)
   → Event archived to OmniDB (direct integration via ingester)
```

**NATS delivers the signal. Skills are the muscles the brain uses to act. Direct integrations handle the plumbing. Plugins modify the brain itself.**

---

## Part 3: Agno Containerized vs OpenClaw Bare Metal — Deep Comparison

### What Was Agno (The Old Architecture)

Agno ran Legacy as a **Python agent inside a Docker container**. Agent definitions lived in `presets.py`. A factory instantiated them. FastAPI served them over HTTP at `:7777`. The agents were ephemeral — they existed within the container lifecycle.

```text
OLD: Agno as the Brain
┌─────────────────────────────────────┐
│ Docker Container (Agno AgentOS)     │
│                                     │
│  FastAPI :7777                      │
│  ├── Legacy agent (main)            │
│  ├── researcher agent               │
│  ├── analyst agent                  │
│  ├── coder agent                    │
│  ├── web_automator agent            │
│  └── assistant agent                │
│                                     │
│  PostgreSQL (separate container)    │
│  ├── Session memory                 │
│  ├── Knowledge base (pgvector)      │
│  └── Agent state                    │
│                                     │
│  No native channels                 │
│  └── Flutter app calls HTTP API     │
└─────────────────────────────────────┘
```

### What Is OpenClaw (The Current Architecture)

OpenClaw runs Legacy as a **persistent Node.js daemon on bare metal**. Agno still exists but is **demoted to a worker pool** — Legacy dispatches tasks to Agno agents, not the other way around.

```text
NOW: OpenClaw as the Brain, Agno as Workers
┌────────────────────────────────────────────────────┐
│ Bare Metal (OpenClaw Gateway)                      │
│                                                    │
│  Legacy daemon (:18789 via Tailscale)              │
│  ├── SOUL.md, AGENTS.md, USER.md (on disk)         │
│  ├── SQLite memory (BM25 + vector search)          │
│  ├── Skills (npm packages, persistent)             │
│  ├── Codex vault (git-synced markdown)             │
│  ├── Cron jobs (heartbeat, git sync)               │
│  └── All channels: WhatsApp, Discord, Telegram,    │
│      iMessage, Voice, Web                          │
│                                                    │
│  When Legacy needs specialist work:                │
│  └── nexus_dispatch → NATS → Agno container        │
│      or agno_dispatch → HTTP → Agno :7777          │
└────────────────────────────────────────────────────┘
         │
         │ dispatches to
         ▼
┌────────────────────────────────────────────────────┐
│ Docker Containers (Worker Pool)                    │
│                                                    │
│  Agno AgentOS :7777                                │
│  ├── researcher (deep research)                    │
│  ├── analyst (data analysis)                       │
│  ├── coder (code generation)                       │
│  ├── web_automator (browser tasks)                 │
│  └── physical_agent (→ NATS → OM1 robot)           │
│                                                    │
│  OmniDB (PostgreSQL + TimescaleDB + pgvector)       │
│  NATS JetStream (event bus)                        │
│  Legacy Bridge (NATS → OpenClaw relay)             │
└────────────────────────────────────────────────────┘
```

---

### Table 1: Runtime Comparison

| Dimension | Agno Container (Old) | OpenClaw Bare Metal (Now) |
| --- | --- | --- |
| **Language** | Python | Node.js |
| **Process model** | Container restart = cold start | nohup/systemd = persistent daemon |
| **Memory** | PostgreSQL in separate container | Native SQLite with hybrid search |
| **Filesystem** | Ephemeral (needs volumes) | Persistent — owns ~/.openclaw/ |
| **Skills/Tools** | Python functions in registry.py | npm packages installed to disk |
| **Identity** | System prompt string in presets.py | SOUL.md, AGENTS.md, USER.md files |
| **Personality updates** | Redeploy container or API call | Edit markdown file, hot-reload |
| **Channels** | HTTP API only (needs proxy for WhatsApp) | Native: WhatsApp, Discord, Telegram, iMessage, Voice, Web |
| **Multi-agent** | Multiple agents in same container | Multi-agent in openclaw.json |
| **Scaling** | Horizontal (more containers) | Vertical (one powerful daemon) |
| **Updates** | Rebuild Docker image | `openclaw update`, restart |
| **Host access** | Sandboxed — no direct host access | Full: disk, cron, git, npm, network |
| **Offline** | Only if container + deps running | Always-on, full offline capability |
| **Voice** | Need LiveKit or custom integration | Built-in Talk Mode + ElevenLabs |
| **Heartbeat** | Custom cron + FastAPI endpoint | Native HEARTBEAT.md + cron |
| **Workflows** | Agno Workflows (Python, deterministic) | Lobster DSL (composable pipelines) |

---

### Table 2: Agent Capabilities

| Capability | Agno (Old Brain) | Agno (Worker Pool) | OpenClaw (Brain) |
| --- | --- | --- | --- |
| **Conversational AI** | Via HTTP API | Same | Via all channels natively |
| **Knowledge Base (RAG)** | Built-in pgvector | Still available | Via Codex + OmniDB pgvector |
| **3-Layer Memory** | Session → persistent → user | Still available | SQLite (session) + OmniDB (long-term) |
| **MCP Tools** | Native support | Still available | Via skill wrappers |
| **Team Coordination** | Route, Broadcast, Coordinate | Still available (Legacy dispatches) | Legacy IS the coordinator |
| **50+ REST APIs** | Auto-generated by AgentOS | Still available for Flutter | Webhooks + skills |
| **Approval Flows (HITL)** | Built-in | Still available | Via Discord/Telegram approval skill |
| **Model Factory** | Per-agent model selection | Still available | Per-agent in openclaw.json |
| **Proactive Behavior** | Needs custom heartbeat | Heartbeat agent (cheap) | Native heartbeat + cron |
| **Communication** | None native | None (worker pool) | WhatsApp, Discord, Telegram, Voice... |
| **Physical Dispatch** | Team member → NATS | physical_agent → OM1 | nexus_dispatch skill → NATS → OM1 |

---

### Table 3: What You Gain by Moving Brain to OpenClaw

| Gain | Why It Matters |
| --- | --- |
| **Always-on presence** | Legacy runs 24/7 as a daemon, not waiting for HTTP requests |
| **Native multi-channel** | One process handles WhatsApp, Discord, Telegram, iMessage, Voice — no proxy |
| **Persistent memory** | SQLite memory survives restarts, no external DB dependency for fast recall |
| **Codex on disk** | Git-synced markdown vault at /home/openclaw/codex/ — direct file access |
| **Sub-millisecond local memory** | SQLite BM25 + vector search, no Docker network hop |
| **Skill ecosystem** | 5,700+ community skills on ClawHub, `npm install` |
| **Identity as files** | SOUL.md, AGENTS.md, USER.md — edit them, behavior changes immediately |
| **Native heartbeat** | HEARTBEAT.md runs every 30 min — proactive without custom infra |
| **Voice built-in** | Talk Mode + wake word + ElevenLabs — no LiveKit needed for local |
| **Cron jobs** | Git sync, health checks, scheduled reports — native, not hacked |
| **Offline-first** | Internet drops? Legacy still thinks, remembers, and responds locally |
| **Node peripherals** | ESP32 via WireClaw, smart glasses via VisionClaw — native edge support |

### Table 4: What You Lose (and How It's Covered)

| Loss from Agno-as-Brain | How OpenClaw Covers It |
| --- | --- |
| Agno's 50+ auto-generated REST APIs | Flutter talks to Agno worker pool for structured data; Legacy via WebSocket/webhook for conversations |
| Agno's built-in team coordination | Legacy IS the coordinator — uses nexus_dispatch to assign tasks to Agno agents |
| Agno's native pgvector RAG | OmniDB has pgvector — same capability, different location |
| Agno's approval flows | Discord/Telegram approval skill — human approval via existing channels |
| Agno's Python tool ecosystem | Skills in JS — or dispatch to Agno for Python-specific tools |
| Agno's model factory | openclaw.json per-agent model config — same idea, different format |
| Container isolation | OpenClaw runs under dedicated `openclaw` user — process isolation, not container isolation |
| Horizontal scaling | Not needed for single-user system — vertical on one Droplet |

---

### Table 5: When Each Pattern Wins

| Scenario | Best Architecture | Why |
| --- | --- | --- |
| **Personal AI assistant** (1 user, multi-channel) | OpenClaw bare metal | Always-on, all channels, persistent identity |
| **SaaS multi-tenant** (1000 users) | Agno containerized | Horizontal scaling, isolation per tenant |
| **Enterprise chatbot** (1 channel, stateless) | Agno containerized | Simple HTTP API, no persistence needed |
| **Human enhancement framework** (OneMind) | OpenClaw brain + Agno workers | Best of both: persistent mind + scalable workers |
| **Physical + digital hybrid** | OpenClaw + NATS + OM1 | Bare metal for consciousness, NATS for nervous system, OM1 for body |
| **Research-heavy workflow** | Agno team coordination | Team leader assigns researcher + analyst + coder in parallel |
| **IoT / smart home integration** | OpenClaw + NATS/MQTT | Native MQTT on :1883, direct sensor access |

---

### Table 6: The Separation of Concerns

| Role | Who Does It | Why |
| --- | --- | --- |
| **Consciousness** (personality, identity, memory) | OpenClaw (bare metal) | Persistent daemon, owns soul files, always-on |
| **Communication** (all channels) | OpenClaw (bare metal) | Native WhatsApp, Discord, etc. — no proxy |
| **Specialist Work** (research, code, analysis) | Agno (containerized) | Best tools: RAG, memory, teams, MCP |
| **Physical Actions** (robots, sensors) | OM1 (containerized or local) | HAL, NLDB, 1Hz cortex loop |
| **Event Bus** (cross-runtime comms) | NATS JetStream (containerized) | Universal nervous system, MQTT bridge |
| **Long-term Memory** (time-series + semantic) | OmniDB (containerized) | TimescaleDB + pgvector in one instance |
| **Real-time AV** (voice/video/glasses) | LiveKit (containerized) | WebRTC, Gemini Live, Meta Ray-Bans |
| **Coordination** (who does what) | Legacy via OpenClaw | Legacy decides, dispatches to workers |

---

## Part 4: Integration Map for OneMind Services

### Current Skills (Installed or Planned)

| Skill | Pattern | Status | What It Does |
| --- | --- | --- | --- |
| nexus-bridge | Skill (JS + NATS) | Installed | Dispatch to Agno agents and OM1 robots |
| github-api | Skill (JS + HTTP) | Installed | Create issues, PRs, manage repos |
| obsidian | Skill (JS + filesystem) | Installed | Read/write Codex vault files |
| home-assistant | Skill (JS + REST/MQTT) | Planned | Control smart home devices |
| todoist | Skill (JS + HTTP) | Planned | Task management |
| discord-notify | Skill (JS + webhook) | Planned | Send alerts, request approvals |
| daily-briefing | Skill (SKILL.md only) | Reference | Morning/evening briefing format |
| code-review | Skill (SKILL.md + scripts) | Reference | Code review instructions |
| clawsec | Skill (JS) | Planned | Security audits, NVD alerts |

### Current NATS Subjects

| Subject | Direction | Purpose |
| --- | --- | --- |
| `digital.agent.{id}.task` | OpenClaw → Agno | Dispatch task to agent |
| `digital.agent.{id}.result.legacy` | Agno → OpenClaw | Agent result back |
| `physical.robot.{id}.command` | OpenClaw → OM1 | Robot command |
| `physical.robot.{id}.status` | OM1 → OpenClaw | Robot state |
| `physical.location.{id}.state` | Sensors → All | Location state |
| `physical.camera.{id}.motion` | Camera → All | Motion detection |
| `nexus.system.health` | All → Monitoring | Health ping |
| `nexus.heartbeat.>` | All runtimes | Heartbeat signals |
| `comms.text.human.zeus` | OpenClaw → All | Zeus's messages |
| `comms.text.agent.legacy` | OpenClaw → All | Legacy's responses |

### Current Direct Integrations

| Integration | File | Purpose |
| --- | --- | --- |
| NATS Ingester | `openclaw/nats/nats_ingester.py` | Archive all NATS events → OmniDB |
| Legacy Bridge | `backend/legacy/bridge.py` | Forward NATS events → OpenClaw webhooks |
| Git Cron | crontab | Pull Codex vault every 5 min |
| Droplet Setup | `.do/setup-droplet.sh` | Initial server provisioning |

---

## Summary

```text
SKILL       = The AI decides when to act     (judgment required)
NATS        = Events flow automatically       (no judgment, real-time)
DIRECT      = Code runs deterministically     (no AI, no events)
PLUGIN      = Modifies the runtime itself     (infrastructure change)

OpenClaw    = The brain (bare metal, always-on, all channels)
Agno        = The workers (containerized, specialist agents)
NATS        = The nervous system (event bus connecting everything)
OmniDB      = The memory (time-series + semantic in one DB)
OM1         = The body (physical robots, sensors, HAL)
LiveKit     = The senses (real-time audio/video/glasses)
```


---

## Related

- [[00-24 SO (Sovereign Operations)/01 Intelligence (AI, Memory & Tools) (The Brain)/01 Intelligence|01 Intelligence]]
