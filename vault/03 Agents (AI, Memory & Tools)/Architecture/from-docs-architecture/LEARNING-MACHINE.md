# OneMind OS: Learning Machine Architecture

> **Status:** Phase 1 Complete, Phase 2 In Progress
> **Created:** 2026-01-24
> **Last Updated:** 2026-01-24
> **Goal:** Transform OneMind OS from single agent to multi-agent hive mind with shared memory

---

## Table of Contents

1. [Executive Summary](#executive-summary)
2. [The Journey: How We Got Here](#the-journey-how-we-got-here)
3. [Key Research Findings](#key-research-findings)
4. [Architecture Decision: Agno + Letta Patterns](#architecture-decision-agno--letta-patterns)
5. [Learning Store Types](#learning-store-types)
6. [Implementation Status](#implementation-status)
7. [Database Schema](#database-schema)
8. [Backend API Reference](#backend-api-reference)
9. [Agent Integration](#agent-integration)
10. [Flutter Integration](#flutter-integration)
11. [Four Pillars Framework](#four-pillars-framework)
12. [Teams & Workflows Integration](#teams--workflows-integration)
13. [Hive Mind Vision](#hive-mind-vision)
14. [Future Roadmap](#future-roadmap)
15. [Testing & Verification](#testing--verification)
16. [Cost Analysis](#cost-analysis)
17. [Files Reference](#files-reference)
18. [External References](#external-references)

---

## Executive Summary

The Learning Machine enables Legacy to **continuously integrate information and improve over time**—across users, sessions, and tasks. Rather than static memory (facts about users), this is **dynamic learning** (evolving knowledge that compounds and sharpens).

### Key Insight

**Agno does NOT have a built-in "LearningMachine" class with 5 stores.** We built custom extensions using Agno's primitives:

- `enable_agentic_memory=True` - Agno's native memory
- `enable_session_summaries=True` - Agno's session tracking
- `dependencies` - Custom context injection (culture, profiles)
- PostgreSQL tables - Custom learning stores

### What This Document Covers

- **Letta (MemGPT) Research** - Tiered memory architecture patterns
- **Agno Learning Machine** - Blog concepts and our implementation
- **Hive Mind Transformation** - Multi-agent shared memory vision
- **Phase 1 Implementation** - What we've built
- **Phase 2 Roadmap** - What's coming next

### Container Status: NO NEW CONTAINERS NEEDED

| Container | Current | Future | Notes |
|-----------|---------|--------|-------|
| **postgres** (pgvector) | Running | Keep | Already has pgvector enabled |
| **redis** | Running | Keep | Hot state cache |
| **nats** | Running | Keep | Event backbone |
| **agno** | Running | Enhance | Add memory modules |
| **ui** | Running | Enhance | Add block management UI |
| **livekit** | Running | Keep | Voice interface |
| **letta server** | - | NOT NEEDED | We recreate patterns, not the service |

---

## The Journey: How We Got Here

### Starting Point: Basic Agno Memory

OneMind OS started with basic Agno memory:

```python
# What we had (basic)
legacy_agent = Agent(
    enable_agentic_memory=True,      # Flat memory in ai.user_memories
    enable_session_summaries=True,    # Session tracking in ai.sessions
)
```

**Limitations:**
- No persistent context about user in every conversation
- No semantic memory retrieval (vector search)
- No memory hierarchy (short-term vs long-term)
- Requires manual memory search by agent
- No shared memory between agents
- Static personality (hardcoded in system prompt)

### The Question

> "How do we make Legacy truly learn and remember everything, share knowledge across agents, and evolve over time?"

### Research Phase

We researched two key sources:

**1. Letta (MemGPT)** - Paper: https://arxiv.org/abs/2310.08560
- Tiered memory architecture (core, archival, recall)
- Core memory always in context (human profile + agent persona)
- Archival memory with vector search
- Self-editing memory tools

**2. Agno Learning Machine Blog** - https://www.agno.com/blog/turning-your-agents-into-learning-machines
- Memory Manager patterns
- Learning store types (7 categories)
- Memory capture instructions
- Dependencies for context injection

### Key Decision

**We decided to build custom extensions using Agno's primitives, incorporating Letta's patterns as inspiration.**

No new containers. No external services. Just PostgreSQL + pgvector + custom code.

---

## Key Research Findings

### From Letta (MemGPT) Research

#### Three-Tier Memory Architecture

```
┌─────────────────────────────────────────────────────────────────────────────┐
│              Agent Context Window (200k tokens)                              │
│  ┌────────────────────────────────────────────────────────────────────────┐ │
│  │ System Prompt + Instructions                                            │ │
│  │                                                                          │ │
│  │  ┌────────────────────────────────────────────────────────────────────┐ │ │
│  │  │ CORE MEMORY (Always Present) - From database, injected every run   │ │ │
│  │  │ ┌────────────────────────┐  ┌──────────────────────────────────┐  │ │ │
│  │  │ │ Human Profile          │  │ Agent Persona                    │  │ │ │
│  │  │ │ - Name: Zeus           │  │ - Name: Legacy                   │  │ │ │
│  │  │ │ - Location: Haven      │  │ - Model: Claude Opus 4.5         │  │ │ │
│  │  │ │ - Focus: 5:30-9AM      │  │ - Personality: Proactive, warm   │  │ │ │
│  │  │ │ - Four Pillars         │  │ - Capabilities: 138+ actions     │  │ │ │
│  │  │ │ - Preferences          │  │ - Guardrails: HITL approvals     │  │ │ │
│  │  │ └────────────────────────┘  └──────────────────────────────────┘  │ │ │
│  │  │ (~2k tokens)                 (~1k tokens)                         │ │ │
│  │  └────────────────────────────────────────────────────────────────────┘ │ │
│  │                                                                          │ │
│  │  ┌────────────────────────────────────────────────────────────────────┐ │ │
│  │  │ PAGED ARCHIVAL MEMORIES (auto-retrieved based on conversation)     │ │ │
│  │  │ 1. [95% match] "Flutter migration: 24 screens, Riverpod, 4 weeks" │ │ │
│  │  │ 2. [88% match] "Linear is Zeus's preferred PM tool"               │ │ │
│  │  │ 3. [82% match] "Multi-agent sprint planned"                       │ │ │
│  │  └────────────────────────────────────────────────────────────────────┘ │ │
│  │                                                                          │ │
│  │  Recent Conversation (Last 10 messages)                                  │ │
│  └────────────────────────────────────────────────────────────────────────┘ │
│                                                                              │
│  Tools: Todoist, GitHub, Comms, + LETTA MEMORY TOOLS                        │
│         - core_memory_append, core_memory_replace                           │
│         - archival_memory_insert, archival_memory_search                    │
│         - conversation_search                                               │
└─────────────────────────────────────────────────────────────────────────────┘
                                   ↕️
┌─────────────────────────────────────────────────────────────────────────────┐
│               POSTGRESQL + PGVECTOR (Memory Persistence)                     │
│                                                                              │
│  ai.core_memory        - Always loaded into agent context                   │
│  ai.archival_memory    - Vector search with VECTOR(1536) column             │
│  ai.recall_memory      - Searchable conversation history                    │
└─────────────────────────────────────────────────────────────────────────────┘
```

#### Letta Database Schema (What They Use)

```sql
-- Letta's blocks table (core memory)
CREATE TABLE blocks (
    id UUID PRIMARY KEY,
    label VARCHAR NOT NULL,      -- 'human', 'persona'
    value TEXT NOT NULL,
    limit INTEGER DEFAULT 2000   -- Character limit
);

-- Letta's passages table (archival memory)
CREATE TABLE passages (
    id UUID PRIMARY KEY,
    text TEXT NOT NULL,
    embedding VECTOR(1536),      -- pgvector column
    metadata_ JSONB DEFAULT '{}'
);

CREATE INDEX ON passages USING ivfflat (embedding vector_cosine_ops);
```

#### Letta Memory Tools

```python
# Core memory tools (always-in-context editing)
core_memory_append(label, content)   # Add to existing block
core_memory_replace(label, old, new) # Replace specific text

# Archival memory tools (long-term storage)
archival_memory_insert(content, tags)  # Store with embedding
archival_memory_search(query, top_k)   # Semantic retrieval

# Recall memory tools (conversation history)
conversation_search(query)             # Search past conversations
```

#### Key Insights from Letta

1. **NOT a service** - Letta is Python code + PostgreSQL + pgvector. No container needed.
2. **Core memory = always in context** - No tool call needed, agent sees it every turn.
3. **Archival memory = vector search** - Semantic similarity via embeddings.
4. **Recall memory has NO embeddings** - Just sequential message storage (Agno already does this).
5. **pgvector IVFFlat index** - Fast approximate nearest neighbor search.
6. **Cost: ~$1/year for 100k memories** - OpenAI embeddings + PostgreSQL storage.

### From Agno Learning Machine Blog

The [Agno blog post](https://www.agno.com/blog/turning-your-agents-into-learning-machines) describes:

#### Memory Manager Pattern

```python
from agno.memory import MemoryManager

memory_manager = MemoryManager(
    model=model,
    db=db,
    memory_capture_instructions="""
    Capture memories tagged by Four Pillars (HP/LE/GE/IT).
    Remember: preferences, decisions, entity relationships.
    """,
)
```

#### Seven Learning Store Types

| Store Type | Purpose | Example |
|------------|---------|---------|
| **User Profile** | Preferences, personal context | Name, timezone, communication style |
| **Session Context** | Goals, plans, progress | Current project, session summary |
| **Entity Memory** | Facts, events, relationships | People, projects, companies |
| **Learned Knowledge** | Insights, patterns | Decision rationale, workflow patterns |
| **Decision Logs** | Rationale behind decisions | Why PostgreSQL over MongoDB |
| **Behavioral Feedback** | Successes and failures | What worked, what didn't |
| **Self-Improvement** | Evolving instructions | Updated personality rules |

#### Key Insight from Agno Blog

> "The Learning Machine enables agents to continuously integrate information and improve over time—across users, sessions, and tasks."

---

## Architecture Decision: Agno + Letta Patterns

### What We Take from Letta (Patterns Only)

| Pattern | Source | Our Implementation |
|---------|--------|-------------------|
| Core memory blocks | Letta `blocks` table | `learning_profiles` + `learning_agent_culture` |
| Archival embeddings | Letta `passages` table | `embedding` column in `ai.user_memories` (Phase 2) |
| Memory self-editing | Letta tools | Culture principles with toggle/priority |
| Block sharing | `shared_with_agents` array | `learning_culture_principles.agent_id` |
| Persona blocks | Letta persona | `learning_agent_culture.communication_style` |

### What We DON'T Take

- ❌ Letta server container (not needed)
- ❌ Letta's database schema (use our own)
- ❌ Letta's SDK (use Agno)
- ❌ Letta's ADE UI (use our Flutter UI)

### Our Hybrid Architecture

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                           ONEMIND LEARNING MACHINE                           │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│   ┌──────────────────┐    ┌──────────────────┐    ┌──────────────────┐     │
│   │   USER PROFILE   │    │  ENTITY MEMORY   │    │  AGENT CULTURE   │     │
│   │   (Letta-style)  │    │  (Knowledge Graph)│    │  (Letta-style)   │     │
│   │                  │    │                  │    │                  │     │
│   │  • Name/Email    │    │  • People        │    │  • Principles    │     │
│   │  • Timezone      │    │  • Projects      │    │  • Comm Style    │     │
│   │  • Preferences   │    │  • Companies     │    │  • Priority      │     │
│   │  • Four Pillars  │    │  • Locations     │    │  • Active/Toggle │     │
│   │  • Custom Fields │    │  • Concepts      │    │                  │     │
│   │                  │    │                  │    │                  │     │
│   │  ┌────────────┐  │    │  ┌────────────┐  │    │  ┌────────────┐  │     │
│   │  │HP│LE│GE│IT │  │    │  │   Facts    │  │    │  │ Priority   │  │     │
│   │  └────────────┘  │    │  │ Relations  │  │    │  │ 10: Comms  │  │     │
│   │  Pillar-tagged   │    │  └────────────┘  │    │  │  9: Pillars│  │     │
│   └──────────────────┘    └──────────────────┘    └──────────────────┘     │
│                                                                              │
│   ┌─────────────────────────────────────────────────────────────────────┐   │
│   │                        AGNO NATIVE MEMORY                            │   │
│   │                                                                      │   │
│   │    ai.user_memories (agentic memory)    ai.sessions (summaries)     │   │
│   │    enable_agentic_memory=True           enable_session_summaries    │   │
│   └─────────────────────────────────────────────────────────────────────┘   │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
                                     │
                                     ▼
                          ┌──────────────────┐
                          │   LEGACY AGENT   │
                          │                  │
                          │  dependencies:   │
                          │  • user_profile  │
                          │  • four_pillars  │
                          │  • culture       │◄─── Injected at runtime
                          └──────────────────┘
```

---

## Learning Store Types

### Phase 1 (Implemented)

| Store | Purpose | Table | Status |
|-------|---------|-------|--------|
| **User Profile** | Preferences, personal context | `learning_profiles` + `_fields` | ✅ Done |
| **Entity Memory** | Knowledge graph (people, projects) | `learning_entities` + `_facts` + `_relations` | ✅ Done |
| **Agent Culture** | Personality, behavioral principles | `learning_agent_culture` + `_principles` | ✅ Done |
| **Session Context** | Goals, plans, progress | `ai.sessions` (Agno native) | ✅ Via Agno |
| **Agentic Memory** | Learned insights | `ai.user_memories` (Agno native) | ✅ Via Agno |

### Phase 2 (Planned)

| Store | Purpose | Table | Status |
|-------|---------|-------|--------|
| **Vector Embeddings** | Semantic search | `embedding VECTOR(1536)` column | ⏳ Planned |
| **Decision Logs** | Rationale behind decisions | `learning_decisions` | ⏳ Planned |
| **Behavioral Feedback** | Successes and failures | `learning_feedback` | ⏳ Planned |
| **Workflow Learnings** | Workflow patterns | `learning_workflow_patterns` | ⏳ Planned |

---

## Implementation Status

### Phase 1: Complete ✅

| Component | Files | Lines | Status |
|-----------|-------|-------|--------|
| Database Schema | `init-db.sql` | ~160 | ✅ Done |
| Backend Models | `core/learning/models.py` | ~250 | ✅ Done |
| Backend Service | `core/learning/service.py` | ~780 | ✅ Done |
| Backend API | `core/learning/api.py` | ~370 | ✅ Done |
| Agent Enhancement | `core/agent/legacy.py` | +50 | ✅ Done |
| Flutter Models | `entity.dart`, `culture.dart` | ~150 | ✅ Done |
| Flutter API Client | `agno_client.dart` | +290 | ✅ Done |
| Flutter Providers | `app_providers.dart` | +30 | ✅ Done |
| Entity Explorer Screen | `entity_explorer_screen.dart` | ~900 | ✅ Done |
| Culture Editor Screen | `culture_editor_screen.dart` | ~600 | ✅ Done |
| Router + Navigation | `app_router.dart`, `app_drawer.dart` | +30 | ✅ Done |
| **Total** | **17 files** | **~3,600** | **✅** |

---

## Database Schema

### Tables Overview

```
┌─────────────────────────┐     ┌─────────────────────────┐
│   learning_profiles     │     │   learning_agent_culture│
│   (User Profile)        │     │   (Agent Personality)   │
├─────────────────────────┤     ├─────────────────────────┤
│ user_id (PK)            │     │ agent_id (PK)           │
│ name                    │     │ communication_style     │
│ email                   │     │ updated_at              │
│ timezone                │     └───────────┬─────────────┘
│ four_pillars            │                 │
│ preferences             │                 ▼
└───────────┬─────────────┘     ┌─────────────────────────┐
            │                   │learning_culture_principles
            ▼                   ├─────────────────────────┤
┌─────────────────────────┐     │ id (PK)                 │
│learning_profile_fields  │     │ agent_id (FK)           │
├─────────────────────────┤     │ key                     │
│ id (PK)                 │     │ value                   │
│ user_id (FK)            │     │ priority                │
│ key                     │     │ active                  │
│ value (JSONB)           │     └─────────────────────────┘
│ pillar                  │
│ source                  │
└─────────────────────────┘

┌─────────────────────────┐
│   learning_entities     │
│   (Knowledge Graph)     │
├─────────────────────────┤
│ id (PK)                 │
│ user_id                 │
│ type                    │───────┐
│ name                    │       │
│ description             │       │
│ pillar                  │       │
│ metadata                │       │
└───────────┬─────────────┘       │
            │                     │
            ▼                     │
┌─────────────────────────┐       │
│ learning_entity_facts   │       │
├─────────────────────────┤       │
│ id (PK)                 │       │
│ entity_id (FK)          │       │
│ fact                    │       │
│ source                  │       │
│ confidence              │       │
└─────────────────────────┘       │
                                  │
            ┌─────────────────────┘
            ▼
┌─────────────────────────┐
│learning_entity_relations│
├─────────────────────────┤
│ id (PK)                 │
│ from_entity_id (FK)     │
│ to_entity_id (FK)       │
│ relation_type           │
│ metadata                │
└─────────────────────────┘
```

### SQL Schema

```sql
-- =============================================================================
-- LEARNING MACHINE: User Profile
-- =============================================================================

CREATE TABLE learning_profiles (
    user_id VARCHAR(255) PRIMARY KEY,
    name VARCHAR(255) NOT NULL DEFAULT 'Zeus',
    email VARCHAR(255),
    timezone VARCHAR(100) DEFAULT 'America/New_York',
    four_pillars JSONB DEFAULT '["HP", "LE", "GE", "IT"]',
    preferences JSONB DEFAULT '{"communication": "concise", "use_bullets": true}',
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE learning_profile_fields (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id VARCHAR(255) REFERENCES learning_profiles(user_id) ON DELETE CASCADE,
    key VARCHAR(255) NOT NULL,
    value JSONB NOT NULL,
    pillar VARCHAR(10),  -- hp, le, ge, it
    source VARCHAR(50) DEFAULT 'manual',  -- manual, conversation, inferred
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(user_id, key)
);

-- =============================================================================
-- LEARNING MACHINE: Entity Memory (Knowledge Graph)
-- =============================================================================

CREATE TABLE learning_entities (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id VARCHAR(255) NOT NULL,
    type VARCHAR(50) NOT NULL,  -- person, project, company, location, concept
    name VARCHAR(255) NOT NULL,
    description TEXT,
    pillar VARCHAR(10),  -- hp, le, ge, it
    metadata JSONB DEFAULT '{}',
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE learning_entity_facts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    entity_id UUID REFERENCES learning_entities(id) ON DELETE CASCADE,
    fact TEXT NOT NULL,
    source VARCHAR(50) DEFAULT 'conversation',
    confidence FLOAT DEFAULT 1.0,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE learning_entity_relations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    from_entity_id UUID REFERENCES learning_entities(id) ON DELETE CASCADE,
    to_entity_id UUID REFERENCES learning_entities(id) ON DELETE CASCADE,
    relation_type VARCHAR(100) NOT NULL,  -- works_at, knows, owns, manages, part_of
    metadata JSONB DEFAULT '{}',
    created_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(from_entity_id, to_entity_id, relation_type)
);

-- =============================================================================
-- LEARNING MACHINE: Agent Culture (Personality)
-- =============================================================================

CREATE TABLE learning_agent_culture (
    agent_id VARCHAR(255) PRIMARY KEY DEFAULT 'legacy',
    communication_style JSONB DEFAULT '{"format": "concise", "use_bullets": true}',
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE learning_culture_principles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    agent_id VARCHAR(255) REFERENCES learning_agent_culture(agent_id) ON DELETE CASCADE,
    key VARCHAR(255) NOT NULL,
    value TEXT NOT NULL,
    priority INT DEFAULT 0,  -- Higher = more important
    active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(agent_id, key)
);
```

### Default Data

```sql
-- Default profile for Zeus
INSERT INTO learning_profiles (user_id, name, timezone, preferences)
VALUES (
    'zeus', 'Zeus Delacruz', 'America/New_York',
    '{"communication": "concise", "use_bullets": true, "morning_focus": "5:30-9:00 AM"}'
);

-- Default profile fields with pillar tagging
INSERT INTO learning_profile_fields (user_id, key, value, pillar, source) VALUES
    ('zeus', 'role', '"One Mind Operator"', NULL, 'manual'),
    ('zeus', 'focus_time', '"5:30-9:00 AM"', 'hp', 'manual'),
    ('zeus', 'family_day', '"Sunday"', 'le', 'manual'),
    ('zeus', 'decision_style', '"data-driven"', NULL, 'manual');

-- Default culture for Legacy (priority 10 = highest)
INSERT INTO learning_culture_principles (agent_id, key, value, priority) VALUES
    ('legacy', 'communication', 'Be direct and efficient, no fluff or pleasantries', 10),
    ('legacy', 'approval_required', 'Request approval for: deletes, finances >$100, external comms', 10),
    ('legacy', 'four_pillars', 'Always consider HP/LE/GE/IT context when responding', 9),
    ('legacy', 'proactive', 'Anticipate needs and proactively suggest actions', 8),
    ('legacy', 'data_driven', 'Base decisions on data when available', 7),
    ('legacy', 'solar_punk', 'Maintain solar-punk optimism - technology serves humanity', 5);
```

---

## Backend API Reference

### Base URL

```
/api/learning
```

### User Profile Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/profile?user_id=zeus` | Get user profile with fields |
| PUT | `/profile?user_id=zeus` | Update profile |
| POST | `/profile/fields?user_id=zeus` | Set/upsert a profile field |

```bash
# Get profile
curl http://localhost:8081/api/learning/profile?user_id=zeus

# Update profile
curl -X PUT http://localhost:8081/api/learning/profile?user_id=zeus \
  -H "Content-Type: application/json" \
  -d '{"name": "Zeus D", "preferences": {"theme": "dark"}}'

# Add profile field
curl -X POST http://localhost:8081/api/learning/profile/fields?user_id=zeus \
  -H "Content-Type: application/json" \
  -d '{"key": "coffee", "value": "black", "pillar": "hp"}'
```

### Entity Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/entities?user_id=zeus` | List entities (supports type, pillar, search filters) |
| POST | `/entities?user_id=zeus` | Create entity |
| GET | `/entities/{id}` | Get entity details |
| PATCH | `/entities/{id}` | Update entity |
| DELETE | `/entities/{id}` | Delete entity |
| GET | `/entities/{id}/facts` | Get facts about entity |
| POST | `/entities/{id}/facts` | Add fact to entity |
| DELETE | `/entities/facts/{id}` | Delete fact |
| GET | `/entities/{id}/relations` | Get entity relationships |
| POST | `/relations` | Create relationship |

```bash
# List entities with filters
curl "http://localhost:8081/api/learning/entities?user_id=zeus&type=person&pillar=ge"

# Create entity
curl -X POST http://localhost:8081/api/learning/entities?user_id=zeus \
  -H "Content-Type: application/json" \
  -d '{"type": "person", "name": "John Smith", "description": "CTO at TechCorp", "pillar": "ge"}'

# Add fact to entity
curl -X POST http://localhost:8081/api/learning/entities/{id}/facts \
  -H "Content-Type: application/json" \
  -d '{"fact": "Prefers async communication", "source": "conversation", "confidence": 0.9}'
```

### Culture Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/culture?agent_id=legacy` | Get agent culture |
| PUT | `/culture?agent_id=legacy` | Update communication style |
| GET | `/culture/context?agent_id=legacy` | Get formatted culture for injection |
| POST | `/culture/principles?agent_id=legacy` | Add/update principle |
| PATCH | `/culture/principles/{id}` | Update principle |
| DELETE | `/culture/principles/{id}` | Delete principle |

```bash
# Get culture
curl http://localhost:8081/api/learning/culture?agent_id=legacy

# Get formatted culture context (for agent injection)
curl http://localhost:8081/api/learning/culture/context?agent_id=legacy

# Add principle
curl -X POST http://localhost:8081/api/learning/culture/principles?agent_id=legacy \
  -H "Content-Type: application/json" \
  -d '{"key": "morning_greeting", "value": "Greet with status update", "priority": 6}'

# Toggle principle
curl -X PATCH http://localhost:8081/api/learning/culture/principles/{id} \
  -H "Content-Type: application/json" \
  -d '{"active": false}'
```

### Stats Endpoint

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/stats?user_id=zeus` | Get learning statistics |

```bash
curl http://localhost:8081/api/learning/stats?user_id=zeus
```

**Response:**
```json
{
  "total_memories": 42,
  "memories_by_pillar": {"hp": 10, "le": 8, "ge": 15, "it": 9},
  "total_entities": 15,
  "entities_by_type": {"person": 8, "project": 4, "company": 3},
  "total_facts": 67,
  "total_relations": 12,
  "culture_principles": 6,
  "session_summaries": 23
}
```

---

## Agent Integration

### Culture Injection via Dependencies

Culture is dynamically injected into the agent context:

```python
# core/agent/legacy.py

from core.learning.service import get_learning_service_async

async def get_culture_context() -> str:
    """Fetch agent culture and format for system prompt injection."""
    try:
        service = await get_learning_service_async()
        culture = await service.get_culture("legacy")

        lines = ["# Agent Culture & Principles"]
        for p in culture.principles:
            if p.active:
                lines.append(f"- **{p.key}**: {p.value}")

        if culture.communication_style:
            lines.append("\n## Communication Style")
            for k, v in culture.communication_style.items():
                lines.append(f"- {k}: {v}")

        return "\n".join(lines)
    except Exception as e:
        logger.warning(f"Failed to fetch culture: {e}")
        return ""

# In create_legacy_agent():
dependencies = {
    "user_profile": get_user_profile,
    "four_pillars": get_four_pillars_state,
    "current_context": get_current_context,
    "culture": get_culture_context,  # Dynamic culture from database
}
```

### Flow Diagram

```
1. User sends message
      │
      ▼
2. Agent created with dependencies
      │
      ▼
3. get_culture_context() called
      │
      ▼
4. Culture fetched from PostgreSQL
      │
      ▼
5. Formatted into system prompt:

   # Agent Culture & Principles
   - **communication**: Be direct and efficient
   - **four_pillars**: Always consider HP/LE/GE/IT context
   - **proactive**: Anticipate needs and suggest actions

   ## Communication Style
   - format: concise
   - use_bullets: true
      │
      ▼
6. Agent responds with culture awareness
```

---

## Flutter Integration

### Models (Freezed)

**entity.dart:**
```dart
enum EntityType { person, project, company, location, concept }
enum Pillar { hp, le, ge, it }

@freezed
sealed class Entity with _$Entity { ... }

@freezed
sealed class EntityFact with _$EntityFact { ... }

@freezed
sealed class EntityRelation with _$EntityRelation { ... }
```

**culture.dart:**
```dart
@freezed
sealed class CulturePrinciple with _$CulturePrinciple { ... }

@freezed
sealed class AgentCulture with _$AgentCulture { ... }

@freezed
sealed class LearningStats with _$LearningStats { ... }
```

### API Client (agno_client.dart)

```dart
// Entities
Future<EntityListResponse> getEntities({EntityType? type, Pillar? pillar, ...})
Future<Entity?> createEntity({required EntityType type, required String name, ...})
Future<List<EntityFact>> getEntityFacts(String entityId, {int limit = 50})
Future<EntityFact?> addEntityFact(String entityId, {required String fact, ...})

// Culture
Future<AgentCulture?> getAgentCulture({String agentId = 'legacy'})
Future<CulturePrinciple?> addCulturePrinciple({required String key, ...})
Future<CulturePrinciple?> updateCulturePrinciple(String id, {bool? active, ...})

// Stats
Future<LearningStats?> getLearningStats({String userId = 'zeus'})
```

### Providers (Riverpod)

```dart
@riverpod
Future<EntityListResponse> learningEntities(ref) async { ... }

@riverpod
Future<AgentCulture?> agentCulture(ref) async { ... }

@riverpod
Future<LearningStats?> learningStats(ref) async { ... }
```

### Screens

| Screen | Route | Features |
|--------|-------|----------|
| Entity Explorer | `/entities` | Stats card, search, type/pillar filters, CRUD dialogs |
| Culture Editor | `/culture` | Communication style editor, principles list with toggle |

---

## Four Pillars Framework

Every entity, profile field, and fact can be tagged with a pillar:

| Pillar | Code | Description |
|--------|------|-------------|
| HP | `hp` | Holistic Performance - Health, wellness, fitness, habits |
| LE | `le` | Legacy Evolution - Family, relationships, home, heritage |
| GE | `ge` | Generational Endeavor - Business, revenue, ventures |
| IT | `it` | Innovative Technology - Infrastructure, automation, AI |

**Example Usage:**
```sql
-- Profile field with pillar
INSERT INTO learning_profile_fields (user_id, key, value, pillar)
VALUES ('zeus', 'focus_time', '"5:30-9:00 AM"', 'hp');

-- Entity with pillar
INSERT INTO learning_entities (user_id, type, name, pillar)
VALUES ('zeus', 'project', 'OneMind OS', 'it');
```

**API Filtering:**
```bash
curl http://localhost:8081/api/learning/entities?pillar=ge&user_id=zeus
```

---

## Teams & Workflows Integration

### Shared Memory Blocks (Phase 2)

| Block Label | Content | Shared With |
|-------------|---------|-------------|
| `zeus_profile` | Identity, preferences | ALL agents (`*`) |
| `pillar_hp` | HP domain knowledge | `legacy`, `hp_agent` |
| `pillar_le` | LE domain knowledge | `legacy`, `le_agent` |
| `pillar_ge` | GE domain knowledge | `legacy`, `ge_agent` |
| `pillar_it` | IT domain knowledge | `legacy`, `it_agent` |
| `onemind_system` | System status | ALL agents (`*`) |

### Team Memory Patterns

```python
# Universal block (all agents)
await manager.create_block(
    label="zeus_profile",
    value=ZEUS_PROFILE,
    shared_with_agents=["*"],
)

# Pillar-specific block
await manager.create_block(
    label="pillar_hp",
    value=HP_KNOWLEDGE,
    shared_with_agents=["legacy", "hp_agent"],
    pillar="HP",
)

# Team block
await manager.create_block(
    label="team_dev",
    value=DEV_CONTEXT,
    shared_with_teams=["dev_team"],
)
```

### Workflow Memory Integration

```python
# Workflow stores learnings
async def morning_briefing_workflow():
    # ... do briefing ...

    # Store pattern in archival memory
    await archival_memory_insert(
        content=f"Morning briefing pattern: {summary}",
        metadata={"type": "workflow_learning", "workflow": "morning_briefing"}
    )
```

---

## Hive Mind Vision

### Current State → Target State

```
CURRENT (Single Agent)                    TARGET (Hive Mind)
┌─────────────────────┐                   ┌─────────────────────────────────┐
│   Legacy Agent      │                   │        LEGACY (Mother)          │
│                     │                   │   Coordination & Routing        │
│   Flat memory       │         →         └──────────────┬──────────────────┘
│   No sharing        │                                  │
│   Static personality│                   ┌──────┬───────┼───────┬──────┐
└─────────────────────┘                   ▼      ▼       ▼       ▼      ▼
                                        HP    LE     GE     IT    ...
                                       Agent  Agent  Agent  Agent  Agents
                                                     │
                                          ┌──────────┴──────────┐
                                          │  SHARED MEMORY LAYER │
                                          │  Core + Archival +   │
                                          │  Entity + Culture    │
                                          └──────────────────────┘
```

### Memory Self-Editing Tools (Phase 2)

```python
# Core memory tools
memory_insert(label, content)      # Add to existing block
memory_replace(label, old, new)    # Replace specific text
memory_rethink(label, new_content) # Full block rewrite

# Archival memory tools
archival_memory_insert(content, tags, pillar)  # Store with embedding
archival_memory_search(query, pillar, top_k)   # Semantic retrieval
```

---

## Future Roadmap

### Phase 2: Vector Embeddings (~6 hours)

| Task | Description |
|------|-------------|
| Enable pgvector | `CREATE EXTENSION vector;` |
| Embedding service | OpenAI text-embedding-3-small |
| Add embedding column | `ALTER TABLE ai.user_memories ADD COLUMN embedding VECTOR(1536);` |
| Vector search API | `/api/memory/search?query=...` |

### Phase 3: Archival Memory (~7 hours)

| Task | Description |
|------|-------------|
| Archival memory table | Long-term storage with vectors |
| Insert endpoint | Store with auto-embedding |
| Search endpoint | Semantic similarity search |
| Auto-paging hook | Retrieve relevant memories before runs |

### Phase 4: Entity Auto-Extraction (~10 hours)

| Task | Description |
|------|-------------|
| Entity extraction prompt | Instructions for agent to extract entities |
| Auto-create entities | Agent creates entities from conversations |
| Auto-create facts | Agent extracts facts about entities |
| Auto-create relations | Agent identifies relationships |

### Phase 5: Multi-Agent Sharing (~12 hours)

| Task | Description |
|------|-------------|
| Pillar agents | HP, LE, GE, IT specialist agents |
| Shared blocks | Block ↔ Agent junction table |
| Block attach/detach API | Dynamic block assignment |
| Team handoffs | Context transfer via handoff blocks |

---

## Testing & Verification

```bash
# Profile
curl http://localhost:8081/api/learning/profile?user_id=zeus

# Entities
curl http://localhost:8081/api/learning/entities?user_id=zeus

# Culture
curl http://localhost:8081/api/learning/culture?agent_id=legacy
curl http://localhost:8081/api/learning/culture/context?agent_id=legacy

# Stats
curl http://localhost:8081/api/learning/stats?user_id=zeus

# Flutter
cd ui && flutter run -d chrome
# Navigate to /entities and /culture
```

---

## Cost Analysis

### Current (Phase 1)

| Item | Cost |
|------|------|
| PostgreSQL storage | $0 (existing) |
| New containers | $0 (none) |
| **Total** | **$0** |

### Future (Phase 2+)

| Item | Monthly |
|------|---------|
| OpenAI embeddings (text-embedding-3-small) | ~$1-5 |
| Additional PostgreSQL storage | ~$1 |
| **Total** | **~$2-6** |

---

## Files Reference

### Backend

| File | Lines | Purpose |
|------|-------|---------|
| `core/learning/__init__.py` | ~15 | Module exports |
| `core/learning/models.py` | ~250 | Pydantic models |
| `core/learning/service.py` | ~780 | Database CRUD |
| `core/learning/api.py` | ~370 | FastAPI router |
| `core/agent/legacy.py` | +50 | Culture injection |

### Frontend

| File | Lines | Purpose |
|------|-------|---------|
| `ui/lib/core/models/entity.dart` | ~95 | Freezed models |
| `ui/lib/core/models/culture.dart` | ~55 | Freezed models |
| `ui/lib/core/api/agno_client.dart` | +290 | API client |
| `ui/lib/features/learning/screens/entity_explorer_screen.dart` | ~900 | Entity UI |
| `ui/lib/features/learning/screens/culture_editor_screen.dart` | ~600 | Culture UI |

### Database

| File | Purpose |
|------|---------|
| `infra/configs/postgres/init-db.sql` | Schema + defaults (lines 289-451) |

---

## External References

### Research Sources

- **Letta (MemGPT) Paper:** https://arxiv.org/abs/2310.08560
- **Letta GitHub:** https://github.com/letta-ai/letta
- **Letta Docs:** https://docs.letta.com
- **Agno Learning Machine Blog:** https://www.agno.com/blog/turning-your-agents-into-learning-machines
- **pgvector:** https://github.com/pgvector/pgvector
- **OpenAI Embeddings:** https://platform.openai.com/docs/guides/embeddings

### OneMind OS Documentation

- [HIVE-MIND-TRANSFORMATION.md](./HIVE-MIND-TRANSFORMATION.md)
- [LETTA-RESEARCH-FINDINGS.md](./LETTA-RESEARCH-FINDINGS.md)
- [LETTA-AGNO-TEAMS-WORKFLOWS.md](./LETTA-AGNO-TEAMS-WORKFLOWS.md)
- [LETTA-MEMORY-INTEGRATION.md](./LETTA-MEMORY-INTEGRATION.md)
- [LETTA-IMPLEMENTATION-CHECKLIST.md](./LETTA-IMPLEMENTATION-CHECKLIST.md)
- [LETTA-ARCHITECTURE-VISUAL.md](./LETTA-ARCHITECTURE-VISUAL.md)
- [MEMORY-FEATURE-COMPLETE.md](../../MEMORY-FEATURE-COMPLETE.md)

---

## Summary

The Learning Machine transforms OneMind OS from a single-agent system with basic memory into a multi-agent hive mind with:

1. **User Profile** - Structured preferences with pillar tagging
2. **Entity Memory** - Knowledge graph of people, projects, companies
3. **Agent Culture** - Dynamic personality principles (self-editing)
4. **Shared Memory** - Blocks shared across agents (Phase 2)
5. **Vector Search** - Semantic retrieval (Phase 2)

**Phase 1 is complete.** We've built the foundation for long-term intelligence.

**Agno Performance + Letta Patterns = Super Long-Term Intelligence**
