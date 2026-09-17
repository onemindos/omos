---
id: onemind-os-architecture-and-agno-memory
title: "OneMind OS — Architecture, Agno Memory & Learning Machine"
type: note
source: apple-notes
imported: 2026-03-07
status: staged
tags:
  - apple-notes-import
  - onemind-os
  - architecture
  - agno
  - memory
  - letta
  - postgresql
  - pgvector
  - livekit
  - nats
  - learning-machine
  - hive-mind
---

# OneMind OS — Architecture, Agno Memory & Learning Machine

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
│  LAYER 2: BRAIN (Agno)                                                 │
│     └── Claude Opus 4.5 via AWS Bedrock                                │
│     └── Self-editing memory (PostgreSQL + pgvector)                    │
│     └── Tools: Todoist, GitHub, Home Assistant, Actions                │
│                                                                        │
│  LAYER 3: NERVOUS SYSTEM (NATS JetStream)                              │
│     └── Sub-millisecond event streaming (200-400μs)                    │
│     └── Events: tasks, code, home, calendar, sensors                   │
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
| **Agent** | Agno + Claude Opus 4.5 | AI brain with self-editing memory |
| **Frontend** | Flutter (iOS, Android, Web, Desktop) | Multi-platform UI |
| **Events** | NATS JetStream | Real-time nervous system |
| **Voice** | LiveKit + Deepgram + Cartesia | Voice I/O |
| **Database** | PostgreSQL + pgvector | Memory storage |
| **Deploy** | Docker Compose + AWS Lightsail | Infrastructure |

---

## Your Vision

| Feature | What You Want | Why It Matters |
|---|---|---|
| Auto-extraction | "Call me Z" → saved automatically | No manual memory management |
| Persistent personality | Each agent has consistent character | Not random each session |
| Shared blocks | 5 agents share zeus_profile | Update once, all know |
| Hive mind sync | Change block → instant propagation | True shared consciousness |
| Archive/long-term | Years of learnings, searchable | Agent grows over time |

### Can Agno Do This Without Letta Patterns?

| Feature | Agno Built-in | Needs Custom |
|---|---|---|
| Memory storage | ✅ `enable_agentic_memory` | — |
| Auto-extraction | ✅ `enable_user_memories` | — |
| Persistent per agent | ⚠️ Via `user_id` isolation | — |
| Shared blocks | ❌ No | Need to build |
| Hive mind sync | ❌ No | Need to build |
| Self-editing tools | ❌ No | Need to build |
| Embeddings/search | ❌ No | Need to build |

---

## My Revised Opinion: Yes, Implement It

The difference:

- **Chatbot** → overkill
- **Cognitive OS with evolving personality** → essential

> Yes. If you're building a cognitive OS that evolves over years, stateful memory is foundational. Every interaction makes the system smarter. That's the whole point.

- Pin your Agno version (`agno==2.1.10`)
- `add_culture_to_context` → ✅ Yes — works, just pin version
- `learning=LearningMachine(...)` → ⚠️ Maybe — newer, test first

---

## What Agno Adds (When You Enable Features)

```
onemind database
│
├── YOUR EXISTING TABLES
│   ├── unified_items
│   ├── approvals
│   ├── inbox_items
│   └── ...
│
└── AGNO TABLES (auto-created)
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
| Where is it stored? | Your PostgreSQL |
| What deletes it? | Only you (or disasters) |
| Will it fill up? | Not for decades |
| Can AI access old memories? | Yes, via smart retrieval |

> **Your AI literally never forgets. That's the point.**

### Memory Optimization (Optional)

If memories pile up too much, Agno has optimization:

```python
# Combine 100 small memories into 10 summarized ones
memory_manager.optimize_memories(
    user_id="zeus",
    strategy="SUMMARIZE",  # Compress without losing info
    apply=True,
)
```

**Before:**

- "Likes coffee"
- "Drinks 3 cups a day"
- "Prefers dark roast"
- "Gets coffee from local shop"
- "Morning coffee at 6am"

**After:**

- "Coffee enthusiast: drinks 3 cups of dark roast daily, starting at 6am, sourced from local shop"

Same info, fewer tokens.

### What Can Delete It

| Event | Data Lost? | Prevention |
|---|---|---|
| Server crash | ❌ No (persisted to disk) | — |
| Agent restart | ❌ No | — |
| New agent version | ❌ No | — |
| You delete the database | ✅ Yes | Don't do that |
| You call `db.clear_memories()` | ✅ Yes | Don't do that |
| Server dies completely | ✅ Yes | Backups |

---

## Context Window vs Storage

Important distinction:

| Concept | Limit | What It Means |
|---|---|---|
| **Storage** | Unlimited | Database holds everything forever |
| **Context window** | ~200k tokens | What the AI can "see" at once |

```
Your database: 10 years of memories (unlimited)
                    │
                    ▼
            ┌───────────────┐
            │ Agno picks    │
            │ relevant ones │  ← Smart retrieval
            └───────────────┘
                    │
                    ▼
            ┌───────────────┐
            │ Context window│  ← Only ~200k tokens fit
            │ (what AI sees)│
            └───────────────┘
```

Agno handles this. It stores everything, but only loads what's relevant.

### Real Example: 10 Years From Now

```
2036: Zeus asks Legacy "What did I name my first child?"

Legacy thinks:
1. Search ai.user_memories for "child" or "baby"
2. Find: "2028-06-12: Had a baby, name is Luna"
3. Answer: "Your first child is Luna, born June 2028"

The memory is 8 years old. Still there. Still works.
```

---

## What You Already Have ✅

| Component | Status | What It Does |
|---|---|---|
| Memory CRUD | ✅ Working | Create/read/update/delete memories |
| Memory UI | ✅ Working | 690-line Flutter screen with search/filter |
| Inbox System | ✅ Working | Real-time notifications + WebSocket |
| Activity Log | ✅ Working | Audit trail + real-time feed |
| Agent Memory | ✅ Enabled | `enable_agentic_memory=True` |
| Session Summaries | ✅ Enabled | `enable_session_summaries=True` |
| Context Injection | ✅ Working | Profile, pillars, location/time |
| HITL Approvals | ✅ Working | Pending approvals workflow |

---

## What's Missing for Agno Learning Machine 🔴

| Feature | Current State | What Agno Provides |
|---|---|---|
| User Profile | Static text block | Structured fields (name, preferences, etc.) that agent can update |
| Entity Memory | None | Knowledge graph — facts about companies, people, projects |
| Learned Knowledge | None | Shared insights across all agents ("morning = HP focus") |
| Culture | None | Agent personality that persists ("be direct", "use bullets") |
| Learning Modes | Manual only | ALWAYS/AGENTIC/PROPOSE automatic extraction |
| Session Context | Basic history | Goals, plans, progress tracking |

---

## What We Need to Build

### Backend Changes (~200 lines)

**1. Update `legacy.py`** — Add Learning Machine config:

```python
# Current: enable_agentic_memory=True (basic)
# Needed: Full LearningMachine configuration
```

**2. New API Endpoints (~150 lines):**

```
GET  /api/profile              → User profile fields
PUT  /api/profile              → Update profile
GET  /api/entities             → Knowledge graph entities
POST /api/entities             → Create entity
GET  /api/entities/{id}/facts  → Entity facts
GET  /api/culture              → Agent culture/principles
PUT  /api/culture              → Update culture
GET  /api/learning/stats       → Learning statistics
```

### Frontend Screens (~1,500 lines)

| Screen | Purpose | Priority |
|---|---|---|
| Profile Screen | View/edit user profile fields | High |
| Entity Explorer | Browse knowledge graph (people, companies, projects) | Medium |
| Culture Editor | View/edit agent personality principles | Medium |
| Learning Dashboard | Stats + control learning modes | Low |
| Hive Mind Status | All agents + shared knowledge | Low |

### Database Changes (Automatic)

Agno auto-creates these tables when you enable Learning Machine:

```sql
ai.user_profiles        -- NEW: Structured user fields
ai.entities             -- NEW: Knowledge graph nodes
ai.entity_facts         -- NEW: Facts about entities
ai.entity_events        -- NEW: Events for entities
ai.entity_relationships -- NEW: Connections between entities
ai.learned_knowledge    -- NEW: Shared insights
ai.cultural_knowledge   -- NEW: Agent personality
```

---

## Implementation Plan

```
┌─────────────────────────────────────────────────────────────────┐
│  PHASE 1: Backend (1-2 hours)                                   │
│  ┌─────────────────────────────────────────────────────────┐    │
│  │ 1. Update core/agent/legacy.py with LearningMachine     │    │
│  │ 2. Add API endpoints for profile, entities, culture     │    │
│  │ 3. Test: "Call me Z" → auto-saves to profile            │    │
│  └─────────────────────────────────────────────────────────┘    │
├─────────────────────────────────────────────────────────────────┤
│  PHASE 2: Profile Screen (1 hour)                               │
│  ┌─────────────────────────────────────────────────────────┐    │
│  │ 1. Create ui/lib/features/profile/ screen               │    │
│  │ 2. Wire to GET/PUT /api/profile                         │    │
│  │ 3. Show: name, role, preferences, communication style   │    │
│  └─────────────────────────────────────────────────────────┘    │
├─────────────────────────────────────────────────────────────────┤
│  PHASE 3: Entity Explorer (2 hours)                             │
│  ┌─────────────────────────────────────────────────────────┐    │
│  │ 1. Create ui/lib/features/entities/ screen              │    │
│  │ 2. Graph visualization (people, companies, projects)    │    │
│  │ 3. Click entity → see facts, events, relationships      │    │
│  └─────────────────────────────────────────────────────────┘    │
├─────────────────────────────────────────────────────────────────┤
│  PHASE 4: Culture Editor (1 hour)                               │
│  ┌─────────────────────────────────────────────────────────┐    │
│  │ 1. Create ui/lib/features/culture/ screen               │    │
│  │ 2. List agent principles                                │    │
│  │ 3. Edit: "Be direct", "Use bullets", "Four pillars"     │    │
│  └─────────────────────────────────────────────────────────┘    │
└─────────────────────────────────────────────────────────────────┘
```

---

## Visual Architecture — Hive Mind

```
┌──────────────────── HIVE MIND ────────────────────┐
│                                                    │
│  ┌─────────┐   ┌─────────┐   ┌─────────┐         │
│  │ Legacy  │   │HP Agent │   │Grid Agent │  ...     │
│  │ (main)  │   │(health) │   │(infra)  │         │
│  └────┬────┘   └────┬────┘   └────┬────┘         │
│       │              │             │               │
│       └──────────────┼─────────────┘               │
│                      │                             │
│               ┌──────▼──────┐                      │
│               │  SHARED DB  │  ← Same PostgreSQL   │
│               │  (user_id)  │  ← Same user_id      │
│               └──────┬──────┘                      │
│                      │                             │
│    ┌─────────────────┼─────────────────┐           │
│    │                 │                 │           │
│    ▼                 ▼                 ▼           │
│ ┌──────┐       ┌──────────┐     ┌──────────┐     │
│ │Profile│      │ Entities │     │ Culture  │     │
│ │"name" │      │ "facts"  │     │"be direct"│    │
│ │"prefs"│      │ "events" │     │"use       │    │
│ └──────┘       └──────────┘     │ bullets"  │    │
│                                 └──────────┘     │
│               ┌──────────────┐                    │
│               │   Learned    │  ← Shared insights │
│               │  Knowledge   │  ← "mornings = HP" │
│               └──────────────┘                    │
└───────────────────────────────────────────────────┘
```

### Quick Win: Start with Backend

Update `core/agent/legacy.py` with the full Learning Machine configuration (~50 lines). That's the foundation — once the backend is learning, we can build the UI to view it.

**The change enables:**

- Automatic profile extraction ("Call me Z" → saved)
- Entity memory (knowledge about people, companies)
- Cultural knowledge (agent personality)
- Shared learning across all agents
