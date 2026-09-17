---
id: so-01-openclaw-framework-structures-20260430
title: OpenClaw Framework Structures
type: architecture
domain: SO
area: "01"
agent: Legacy
status: active
tags:
  - domain/so
  - area/01
  - type/reference
  - status/active
  - openclaw
  - nats
  - models
  - memory
  - scoring
related:
  - "[[AGENTS]]"
  - "[[AGENT-ASSETS]]"
  - "[[HIVE-AGENT-ARCHITECTURE]]"
  - "[[00-24 SO (Sovereign Operations)/00 Framework (Doctrine)/architecture/diagrams/memory-architecture]]"
  - "[[LLM-STRATEGY]]"
  - "[[LIVEKIT-NATS-BOUNDARIES]]"
todoist_project_id: 6g6P9rCGG89q6qPr
updated: 2026-04-30
source: extracted from 00-Framework/architecture/archive — architecture-synthesis.md, LEARNING-MACHINE-ARCHITECTURE.md, UNIFIED_ARCHITECTURE_PLAN.md, unified-physical-digital-architecture.md
---

# OpenClaw Framework Structures

> *The foundational design structures of OneMind OS — how the system is architected to think, route, score, store memory, and assign intelligence across layers. Extracted and consolidated from multiple architecture cycles. Use this as the structural guide when building or extending OpenClaw.*

---

## 1. LLM Assignment Map

Every agent has an assigned model. Don't over-provision — pick the right tier for the job.

| Agent / Layer | Model | Why |
|---|---|---|
| **Legacy** (OpenClaw consciousness) | Claude Opus 4.6 | The consciousness deserves the best — all reasoning, routing, voice |
| **Researcher** (Agno) | Claude Opus / GPT-4.1 | Deep research and complex reasoning — cost justified by quality |
| **Analyst** (Agno) | Claude Sonnet 4.6 | Strong analysis, faster than Opus |
| **Coder** (Agno) | Claude Sonnet 4.6 | Best code model currently |
| **Assistant** (Agno, quick tasks) | GPT-4.1-mini | Fast, cheap, sufficient for general tasks |
| **Web Automator** (Agno) | GPT-4o | Vision required for browser screenshots |
| **Heartbeat** (Agno background) | GPT-4.1-mini | Cheap — runs 48x/day, only escalates if something found |
| **OM1 Cortex** (physical reflexes) | Qwen3-30B local / GPT-4.1-nano fallback | 200ms latency, runs 3600x/hr, offline-capable |
| **LiveKit Voice Sessions** | Gemini 2.5 Pro (Gemini Realtime) | Best for real-time voice + vision (Meta Ray-Bans, mobile) |
| **Spartan / Coach / Haven / Empire** (OpenClaw domain agents) | Claude Sonnet 4.6 | Domain specialists — balanced cost/capability |

---

## 2. LLM Decision Tree

Use this when deciding which model to invoke for a given task:

```
Is this a real-time physical decision? (< 500ms needed)
  YES → OM1 local Qwen3-30B
  NO  → continue...

Is this a voice/vision session via WebRTC (LiveKit)?
  YES → Gemini Realtime
  NO  → continue...

Is this a consciousness/communication/routing decision?
  YES → Claude Opus 4.6 (Legacy in OpenClaw)
  NO  → continue...

Is this deep research, complex reasoning, or multi-step analysis?
  YES → Claude Opus 4.6 or GPT-4.1 (Agno researcher/analyst)
  NO  → continue...

Is this a quick task, heartbeat check, or routine operation?
  YES → GPT-4.1-mini (cheap, fast, sufficient)
```

---

## 3. Layered Autonomy — Proactive Architecture

The full stack is neither always-reactive nor always-polling. Each layer has the right frequency for its purpose:

| Layer | Frequency | Model | Trigger |
|---|---|---|---|
| OM1 Physical Loop | 1 Hz | Qwen3-30B local | Always-on (physical safety, reflexes) |
| OpenClaw Heartbeat | Every 30 min | Claude Opus 4.6 | Time + events + directives |
| Agno Background Loop | Every 15 min | GPT-4.1-mini | Gathered data digest |
| Todoist/ClickUp Webhooks | On-event | GPT-4.1-mini | Task create/complete/update |
| Home Assistant Webhooks | On-event | Claude Opus 4.6 | Critical home events |
| NATS Robot Alerts | On-event | Claude Opus 4.6 | Physical world alerts |
| User Conversation | On-demand | Full model (Legacy) | Explicit message |

**Key principle**: Legacy is always **prepared** — the heartbeat has already checked context before Zeus sends a message. She's been watching.

---

## 4. OpenClaw vs LiveKit — Channel Split

These are **complementary**, not competing. Never remove one for the other.

| Capability | OpenClaw | LiveKit |
|---|---|---|
| WhatsApp | ✅ Native | ❌ |
| iMessage | ✅ Native | ❌ |
| Telegram | ✅ Native | ❌ |
| Discord | ✅ Native | ❌ |
| Local wake word / Talk Mode | ✅ Native | ⚠️ Needs custom |
| Heartbeat / Cron | ✅ Native | ❌ |
| Webhooks | ✅ Native | ❌ |
| Skills / Node peripherals | ✅ Native | ❌ |
| Meta Ray-Bans (Wearables API) | ❌ | ✅ Native |
| Mobile WebRTC (cellular voice) | ❌ | ✅ Native |
| Remote robot video streaming | ❌ | ✅ Native |
| Real-time multi-user room | ❌ | ✅ Native |
| Gemini Live Vision | ❌ | ✅ Native |

**Rule**: OpenClaw owns async channels. LiveKit owns real-time WebRTC sessions.

---

## 5. Universal Chat Routing Model

When any message arrives, Legacy routes it through this model:

```
mode: direct | agent | team
target_id: optional (agent_id or team_id)
model_id: optional override (defaults to Legacy's assigned model)
tools: [] optional list

direct  → model responds with optional tools, no agent wrapper
agent   → route to named OpenClaw agent (Sentinel, Coach, Haven, Empire)
team    → route to Agno team (researcher+analyst+coder etc.)
```

**Operations vs Assets distinction** (always applies):
- **Operations** (Command Center): Inbox, Tasks, Calendar, Chat
- **Assets** (Entities): Agents, Humans, Machines, Locations

Asset types recognized by the awareness layer:
- `LOCATION` — buildings, zones, outdoor areas, farms
- `DEVICE` — phones, sensors, cameras, beacons
- `HUMAN` — with biometrics (heart_rate, blood_oxygen, steps, stress, mood)
- `MACHINE` — robots, drones, vehicles, FarmBot, 3D printers

---

## 6. PI / LI / WI Domain Scoring Model

Three scores map to the three pillars. Used for agent context building, daily brief, and goal tracking.

| Score | Pillar | Domain | Formula |
|---|---|---|---|
| **PI** — Performance Index | HP | 25-49 | health_habits (40%) + energy_mood (20%) + focus (20%) + goals (20%) |
| **LI** — Legacy Index | LE | 50-74 | family_habits (35%) + home (20%) + goals (20%) + legacy (25%) |
| **WI** — Wealth Index | GE | 75-99 | business_tasks (30%) + business_goals (30%) + work_habits (20%) + finance (20%) |

All three scores live in the OpenClaw heartbeat context. Legacy references them during morning brief and daily check-ins.

---

## 7. Knowledge Graph Memory Model

This is how OpenClaw memory should be structured. Three-layer graph:

### Layer 1 — Entities (Nodes)
```
entity:
  id: uuid
  type: person | project | company | location | concept | task | habit | goal
  name: string
  description: text
  pillar: hp | le | ge | so
  metadata: {}   # flexible — stores lifeos_id, source_table, priority, etc.
```

### Layer 2 — Facts (What AI learns about entities)
```
fact:
  entity_id: uuid → entity
  fact: string          # "Completed 15 min under estimate"
  source: user | system | inferred | imported
  confidence: 0.0–1.0
```

Example facts:
- `"Completed on Monday mood=7/9 energy=8/9"`
- `"Skipped 3x in a row — pattern: low energy days"`
- `"Progress increased from 40% → 60% in one sprint"`

### Layer 3 — Relations (Graph Edges)
```
relation:
  from_entity_id: uuid
  to_entity_id: uuid
  relation_type: part_of | supports | requires | child_of | assigned_to | blocks
  metadata: {}
```

Relation types:
- `part_of` — Task is part of Project
- `child_of` — Sub-goal under parent Goal
- `supports` — Habit supports Goal
- `requires` — Task requires another Task to complete first
- `assigned_to` — Task/project assigned to Person/Agent
- `blocks` — Task blocks another

**This graph is the memory layer Legacy traverses for context.** When Zeus asks about a project, Legacy walks: project → tasks (part_of) → habits (supports) → facts → recent completions.

---

## 8. NATS Subject Master Map

The canonical namespace for all NATS topics in OneMind OS.

```
# OpenClaw ↔ Agno (via bridge)
comms.text.human.zeus                  ← Zeus's inbound messages
comms.text.agent.legacy                ← Legacy's outbound
digital.agent.{id}.task               ← Dispatch to Agno worker
digital.agent.{id}.result.legacy      ← Agno result back to Legacy
system.alert.{type}                   ← Critical system alerts
nexus.heartbeat.failure.{id}          ← Runtime heartbeat lost

# OM1 Physical Layer
onemind.commands.{robot_id}           ← Agno → OM1 (send command)
onemind.vision.{robot_id}            ← OM1 → Agno (camera feed caption)
onemind.audio.{robot_id}             ← OM1 → Agno (mic input)
onemind.status.{robot_id}            ← OM1 → Agno (battery, position)
onemind.alerts.{robot_id}            ← OM1 → Agno (urgent physical events)
onemind.personality                  ← Agno → OM1 (soul/context sync every 5min)
onemind.fleet.worldstate             ← Aggregated fleet view (all robots)

# Location-Scoped
physical.location.{id}.state         ← Query location state
physical.location.{id}.heartbeat     ← Location alive check

# FarmBot
physical.farmbot.{id}.sensor         ← Soil moisture, temperature
physical.farmbot.{id}.status         ← Current position, job, sequence
physical.farmbot.{id}.command        ← Trigger sequence or move arm
events.farmbot.{id}.alert            ← Plant health alerts, low water

# Meshtastic (Off-grid mesh)
physical.mesh.{location}.{node}.sensor   ← LoRa sensor data
physical.mesh.{location}.alert           ← Urgent offline alerts
physical.mesh.{location}.heartbeat       ← Node still alive

# Community
events.community.member.joined       ← Circle.so new member
events.community.{event_type}        ← General community events
events.task.completed.{id}           ← Todoist/ClickUp task done
events.task.created.{id}             ← New task created
```

---

## 9. Eden — Physical Agent Spec (FarmBot / Garden Robot)

The vision for the garden physical agent. When OM1 integration is live, Eden is the first physical agent to deploy.

**Role**: Manages farm/garden operations — watering, plant health monitoring, soil sensors, harvest planning.

**Capabilities**:
- Move FarmBot arm to coordinates → water specific plant
- Query soil sensor via FarmBot MQTT → moisture level
- Camera vision via OM1 → VLM captions plant health ("Leaves droopy, possible underwatering")
- Publish results to NATS → Legacy synthesizes → notifies Zeus

**NATS subjects Eden owns**:
```
physical.farmbot.{id}.sensor        ← publishes soil data
physical.farmbot.{id}.command       ← subscribes to execute
events.farmbot.{id}.alert           ← publishes plant health alerts
```

**Trigger example**:
```
Zeus → "Check if tomatoes in zone 3 need water"
  → Legacy dispatches to physical_agent via NATS
  → OM1 + FarmBot MQTT query soil sensor (returns 22% moisture)
  → OM1 camera: "Leaves slightly droopy"
  → Eden publishes: "Zone 3 soil at 22%, recommend irrigation"
  → Legacy: "Zone 3 is dry — triggering FarmBot irrigation"
  → FarmBot executes watering sequence
```

---

## See Also

- [[HIVE-AGENT-ARCHITECTURE]] — The 9-agent hive structure built on these models
- [[00-24 SO (Sovereign Operations)/00 Framework (Doctrine)/architecture/diagrams/memory-architecture]] — Full memory layer implementation
- [[LIVEKIT-NATS-BOUNDARIES]] — LiveKit vs OpenClaw boundary in detail
- [[LLM-STRATEGY]] — Model routing and selection strategy
- [[MULTI-AGENT]] — 4-agent-per-domain OpenClaw design
- [[HOMESTEAD-EDGE-ARCHITECTURE]] — 6-building edge intelligence grid (Eden's home)
