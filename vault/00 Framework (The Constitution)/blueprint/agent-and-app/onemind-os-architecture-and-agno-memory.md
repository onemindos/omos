---
id: onemind-os-architecture-and-agno-memory
title: "OneMind OS — Architecture, OpenClaw Memory & Learning Machine"
type: note
source: apple-notes
imported: 2026-03-07
modified: 2026-04-30
status: staged
tags:
  - apple-notes-import
  - onemind-os
  - architecture
  - openclaw
  - memory
  - nats
  - livekit
  - learning-machine
  - hive-mind
---

# OneMind OS — Architecture, OpenClaw Memory & Learning Machine

---

## OneMind OS Architecture

```
┌────────────────────────────────────────────────────────────────────────┐
│                        ONEMIND OS ARCHITECTURE                         │
├────────────────────────────────────────────────────────────────────────┤
│                                                                        │
│  LAYER 1: SENSES (LiveKit)                                             │
│     └── Real-time voice via Deepgram STT + Cartesia TTS               │
│     └── WebRTC for video/screen sharing                                │
│                                                                        │
│  LAYER 2: BRAIN (OpenClaw / LEGACY)                                    │
│     └── Claude Sonnet 4.6 via AWS Bedrock                             │
│     └── Stateful agent memory (PostgreSQL + pgvector)                  │
│     └── Tools: Todoist, GitHub, Home Assistant, MCP actions            │
│                                                                        │
│  LAYER 3: NERVOUS SYSTEM (NATS JetStream)                              │
│     └── Sub-millisecond event streaming (200-400μs)                    │
│     └── Events: tasks, code, home, calendar, sensors                   │
│     └── Agent-to-Agent (A2A) messaging                                 │
│                                                                        │
│  LAYER 4: BODY (Integrations)                                          │
│     └── Home Assistant for smart home                                  │
│     └── FarmBot for agriculture                                        │
│     └── Frigate for vision/cameras                                     │
│                                                                        │
└────────────────────────────────────────────────────────────────────────┘
```

### Tech Stack

| Layer | Technology | Purpose |
|---|---|---|
| **Agent Runtime** | OpenClaw + Claude Sonnet 4.6 | AI brain with stateful memory |
| **Frontend** | Flutter (iOS, Android, Web, Desktop) | Multi-platform UI |
| **Events / A2A** | NATS JetStream | Real-time nervous system |
| **Voice** | LiveKit + Deepgram + Cartesia | Voice I/O |
| **Database** | PostgreSQL + pgvector | Memory storage |
| **Deploy** | Docker Compose + AWS Lightsail (omos-core) | Infrastructure |

---

## Agent Memory Vision

| Feature | What We Want | Why It Matters |
|---|---|---|
| Auto-extraction | "Call me Z" → saved automatically | No manual memory management |
| Persistent personality | Each agent has consistent character | Not random each session |
| Shared blocks | 9 agents share zeus_profile | Update once, all know |
| Hive mind sync | Change block → instant propagation | True shared consciousness |
| Archive/long-term | Years of learnings, searchable | Agent grows over time |

### OpenClaw Memory Capabilities

| Feature | Status | Notes |
|---|---|---|
| Memory storage | ✅ Built-in | PostgreSQL-backed |
| Auto-extraction | ✅ Enabled | Preference extraction from conversation |
| Persistent per agent | ✅ Via agent_id isolation | Each agent has own scope |
| Shared memory blocks | ✅ NATS propagation | zeus_profile shared across all agents |
| Hive mind sync | ✅ NATS JetStream | Change → instant publish to all agents |
| Self-editing tools | ✅ Agent memory tools | Agents update their own memory blocks |
| Embeddings/search | ✅ pgvector | Semantic search over agent memories |

---

## Memory Database Schema

```
onemind database
│
├── YOUR EXISTING TABLES
│   ├── unified_items
│   ├── approvals
│   ├── inbox_items
│   └── ...
│
└── AGENT MEMORY TABLES
    │
    ├── ai.agent_sessions        ← Chat history & session data
    ├── ai.user_memories         ← "Call me Z", preferences
    ├── ai.user_profiles         ← Structured: name, role, etc.
    ├── ai.session_summaries     ← Compressed conversation summaries
    ├── ai.entities              ← Knowledge about companies, people
    ├── ai.entity_facts          ← Facts linked to entities
    ├── ai.entity_events         ← Events linked to entities
    ├── ai.entity_relationships  ← Connections between entities
    ├── ai.learned_knowledge     ← Shared insights/patterns
    └── ai.cultural_knowledge    ← Agent personality/principles
```

### Who Sees What

| Table | What Goes In | Who Sees It |
|---|---|---|
| `ai.agent_sessions` | Chat messages | Per-session |
| `ai.user_memories` | "Call me Z" facts | All agents (same `user_id`) |
| `ai.user_profiles` | Name, role, preferences | All agents (same `user_id`) |
| `ai.cultural_knowledge` | Agent personality rules | All agents (or specific) |
| `ai.entities` | Companies, people, projects | All agents |
| `ai.learned_knowledge` | Tips, patterns, wisdom | All agents (via namespace) |

Everything is in **YOUR** `onemind-postgres` container. You can:

- **Query it:** `SELECT * FROM ai.user_memories WHERE user_id = 'zeus'`
- **Back it up:** Normal `pg_dump`
- **Export it:** Your data, your choice

---

## Data Longevity

| Question | Answer |
|---|---|
| How long does data last? | **Forever** |
| Where is it stored? | Your PostgreSQL (omos-core) |
| What deletes it? | Only you (or disasters) |
| Will it fill up? | Not for decades |
| Can AI access old memories? | Yes, via pgvector semantic retrieval |

> **Your AI literally never forgets. That's the point.**

### Memory Optimization

If memories pile up too much:

```python
# Combine many small memories into fewer summarized ones
memory_manager.optimize_memories(
    user_id="zeus",
    strategy="SUMMARIZE",
    apply=True,
)
```

**Before:**
- "Likes coffee"
- "Drinks 3 cups a day"
- "Prefers dark roast"

**After:**
- "Coffee enthusiast: drinks 3 cups of dark roast daily, starting at 6am, sourced from local shop"

### What Can Delete It

| Event | Data Lost? | Prevention |
|---|---|---|
| Server crash | ❌ No (persisted to disk) | — |
| Agent restart | ❌ No | — |
| New agent version | ❌ No | — |
| You delete the database | ✅ Yes | Don't do that |
| Server dies completely | ✅ Yes | Regular backups |

---

## Context Window vs Storage

```
Your database: 10 years of memories (unlimited)
                    │
                    ▼
            ┌───────────────┐
            │ OpenClaw picks│
            │ relevant ones │  ← pgvector semantic retrieval
            └───────────────┘
                    │
                    ▼
            ┌───────────────┐
            │ Context window│  ← Only ~200k tokens fit
            │ (what AI sees)│
            └───────────────┘
```

OpenClaw handles this. It stores everything, but only loads what's relevant per query.

### Real Example: 10 Years From Now

```
2036: Zeus asks LEGACY "What did I name my first child?"

LEGACY thinks:
1. Search ai.user_memories for "child" or "baby"
2. Find: "2028-06-12: Had a baby, name is Luna"
3. Answer: "Your first child is Luna, born June 2028"

The memory is 8 years old. Still there. Still works.
```

---

## What's Built ✅

| Component | Status | What It Does |
|---|---|---|
| Memory CRUD | ✅ Working | Create/read/update/delete memories |
| Memory UI | ✅ Working | Flutter screen with search/filter |
| Inbox System | ✅ Working | Real-time notifications + WebSocket |
| Activity Log | ✅ Working | Audit trail + real-time feed |
| Agent Memory | ✅ Enabled | Agentic memory on |
| Session Summaries | ✅ Enabled | Compressed conversation summaries |
| Context Injection | ✅ Working | Profile, pillars, location/time |
| HITL Approvals | ✅ Working | Pending approvals workflow |
| NATS A2A | ✅ Working | Agent-to-agent messaging on omos-core |

---

## Learning Machine — What's Missing 🔴

| Feature | Current State | What to Build |
|---|---|---|
| User Profile | Static text block | Structured fields agents can update |
| Entity Memory | None | Knowledge graph — facts about companies, people, projects |
| Learned Knowledge | None | Shared insights across all agents |
| Culture | None | Agent personality that persists |
| Learning Modes | Manual only | ALWAYS/AGENTIC/PROPOSE auto-extraction |

---

## Implementation Plan

```
PHASE 1: Backend (1-2 hours)
  1. Extend agent memory config in LEGACY
  2. Add API endpoints for profile, entities, culture
  3. Test: "Call me Z" → auto-saves to profile

PHASE 2: Profile Screen (1 hour)
  1. Create Flutter profile screen
  2. Wire to GET/PUT /api/profile
  3. Show: name, role, preferences, communication style

PHASE 3: Entity Explorer (2 hours)
  1. Create Flutter entities screen
  2. Graph visualization (people, companies, projects)
  3. Click entity → see facts, events, relationships

PHASE 4: Culture Editor (1 hour)
  1. Create Flutter culture screen
  2. List agent principles
  3. Edit: "Be direct", "Use bullets", "Four domains"
```

---

## Hive Mind Architecture

```
┌──────────────────── HIVE MIND ────────────────────┐
│                                                    │
│  ┌─────────┐   ┌─────────┐   ┌─────────┐         │
│  │ LEGACY  │   │ Spartan │   │  Grid   │  ...     │
│  │ (main)  │   │  (HP)   │   │  (SO)   │         │
│  └────┬────┘   └────┬────┘   └────┬────┘         │
│       │              │             │               │
│       └──────────────┼─────────────┘               │
│                      │                             │
│            NATS JetStream (A2A)                    │
│                      │                             │
│               ┌──────▼──────┐                      │
│               │  SHARED DB  │  ← Same PostgreSQL   │
│               │  (user_id)  │  ← Same user_id      │
│               └──────┬──────┘                      │
│                      │                             │
│    ┌─────────────────┼─────────────────┐           │
│    ▼                 ▼                 ▼           │
│ ┌──────┐       ┌──────────┐     ┌──────────┐      │
│ │Profile│      │ Entities │     │ Culture  │      │
│ │"name" │      │ "facts"  │     │"be direct"│     │
│ │"prefs"│      │ "events" │     │"use       │     │
│ └──────┘       └──────────┘     │ bullets"  │     │
│                                 └──────────┘      │
│               ┌──────────────┐                    │
│               │   Learned    │  ← Shared insights  │
│               │  Knowledge   │  ← "mornings = HP"  │
│               └──────────────┘                    │
└───────────────────────────────────────────────────┘
```

### Agent Roster (OpenClaw)

| Agent | Domain | Focus |
|-------|--------|-------|
| LEGACY | SO (master) | Orchestration, routing, full tool access |
| Spartan | HP | Physical performance, fitness |
| Coach | HP | Mental performance, habits |
| Oracle | GE | Financial intelligence, investments |
| Empire | GE | Business operations, revenue |
| Eden | LE | Homestead, family environment |
| Haven | LE | Relationships, heritage |
| Grid | SO | Infrastructure, DevOps |
| Guardian | SO | Security, compliance |

All agents share `zeus_profile` memory block via NATS JetStream. Update once → all 9 agents know instantly.


---

## Related

- [[00-24 SO (Sovereign Operations)/00 Framework (Doctrine) (The Constitution)/00 Framework|00 Framework]]
