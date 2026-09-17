# OneMind OS: Single Agent → Hive Mind Transformation

> **Created:** 2026-01-24
> **Status:** Planning
> **Goal:** Transform OneMind OS from single Legacy agent to multi-agent hive mind with shared memory

---

## Executive Summary

OneMind OS currently runs a single "Legacy" agent. This document outlines the transformation to a **hive mind architecture** where multiple specialized agents share memory blocks, collaborate on tasks, and maintain a unified understanding of Zeus's world.

**Key insight:** Legacy is not just an agent - she IS the system. Legacy = PostgreSQL + NATS + Redis + LiveKit. She's the intelligence layer that coordinates all agents.

---

## Container Status: NO NEW CONTAINERS NEEDED

| Container | Current | Future | Notes |
|-----------|---------|--------|-------|
| **postgres** (pgvector) | ✅ Running | ✅ Keep | Already has pgvector enabled! |
| **redis** | ✅ Running | ✅ Keep | Hot state cache |
| **nats** | ✅ Running | ✅ Keep | Event backbone |
| **agno** | ✅ Running | ✅ Enhance | Add memory modules |
| **ui** | ✅ Running | ✅ Enhance | Add block management UI |
| **livekit** | ✅ Running | ✅ Keep | Voice interface |
| **timescaledb** | ✅ Running | ✅ Keep | Event history |
| **letta server** | ❌ | ❌ NOT NEEDED | We recreate patterns, not the service |

---

## Database Changes Required

| Table | Current | Future | Change Type |
|-------|---------|--------|-------------|
| `ai.user_memories` | Basic Agno memories | + embeddings, + labels, + sharing | **ALTER** |
| `ai.memory_blocks` | ❌ Doesn't exist | Core memory blocks | **CREATE** |
| `ai.block_assignments` | ❌ Doesn't exist | Block ↔ Agent junction | **CREATE** |
| `unified_items` | ✅ Exists | No change | - |
| `inbox_items` | ✅ Exists | No change | - |
| `approvals` | ✅ Exists | No change | - |

---

## Agent Architecture Comparison

| Aspect | Current | Future |
|--------|---------|--------|
| **Agent count** | 1 (Legacy) | Many (Legacy + Pillar agents + Specialists) |
| **Memory type** | `enable_agentic_memory=True` (flat) | Core blocks + Archival + Recall |
| **Memory sharing** | Via `user_id` only | Block-level sharing with `shared_with_agents` |
| **Self-editing** | ❌ No | ✅ memory_insert, memory_replace, memory_rethink |
| **Vector search** | ❌ No | ✅ archival_memory_search |
| **Block assignment** | ❌ No | ✅ Attach/detach via API ("button press") |

---

## The Transformation Diagram

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                           CURRENT (Single Agent)                            │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│   ┌──────────────────────────────────────────┐                             │
│   │               Legacy Agent               │                             │
│   │  ┌─────────┐ ┌─────────┐ ┌────────────┐ │                             │
│   │  │ todoist │ │ github  │ │  browser   │ │                             │
│   │  └─────────┘ └─────────┘ └────────────┘ │                             │
│   │                                          │                             │
│   │  Memory: enable_agentic_memory=True      │                             │
│   │  (Flat, no blocks, no sharing)           │                             │
│   └──────────────────────────────────────────┘                             │
│                         │                                                   │
│                         ▼                                                   │
│                    PostgreSQL                                               │
│               (ai.user_memories - flat)                                     │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘

                                    │
                                    │ TRANSFORMATION
                                    ▼

┌─────────────────────────────────────────────────────────────────────────────┐
│                        FUTURE (Hive Mind Architecture)                      │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│                    ┌──────────────────────────┐                             │
│                    │      LEGACY (Mother)     │                             │
│                    │   Intelligence Layer     │                             │
│                    │  Coordination & Routing  │                             │
│                    └────────────┬─────────────┘                             │
│                                 │                                           │
│        ┌────────────┬───────────┴───────────┬────────────┐                  │
│        ▼            ▼                       ▼            ▼                  │
│   ┌─────────┐  ┌─────────┐            ┌─────────┐  ┌─────────┐             │
│   │HP Agent │  │LE Agent │            │GE Agent │  │Grid Agent │             │
│   │ Health  │  │ Legacy  │            │Business │  │ Tech    │             │
│   │ Fitness │  │ Family  │            │ Revenue │  │ Infra   │             │
│   └────┬────┘  └────┬────┘            └────┬────┘  └────┬────┘             │
│        │            │                       │            │                  │
│        └────────────┴───────────┬───────────┴────────────┘                  │
│                                 │                                           │
│                    ┌────────────┴────────────┐                              │
│                    │   SHARED MEMORY LAYER   │                              │
│                    │                         │                              │
│  ┌─────────────────┴─────────────────────────┴─────────────────────┐       │
│  │                                                                  │       │
│  │  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐              │       │
│  │  │ CORE BLOCKS │  │  ARCHIVAL   │  │   RECALL    │              │       │
│  │  │ (In-Context)│  │ (Embeddings)│  │  (History)  │              │       │
│  │  ├─────────────┤  ├─────────────┤  ├─────────────┤              │       │
│  │  │ zeus_profile│  │ Decisions   │  │ Chat logs   │              │       │
│  │  │ pillar_hp   │  │ Learnings   │  │ Session     │              │       │
│  │  │ pillar_le   │  │ Insights    │  │ Summaries   │              │       │
│  │  │ pillar_ge   │  │ Facts       │  │             │              │       │
│  │  │ pillar_it   │  │ Patterns    │  │             │              │       │
│  │  │ onemind_sys │  │             │  │             │              │       │
│  │  └─────────────┘  └─────────────┘  └─────────────┘              │       │
│  │                                                                  │       │
│  │  shared_with_agents: ['*'] or ['legacy', 'hp_agent', ...]       │       │
│  └──────────────────────────────────────────────────────────────────┘       │
│                                 │                                           │
│                    ┌────────────┴────────────┐                              │
│                    │       PostgreSQL        │                              │
│                    │  ┌──────────────────┐   │                              │
│                    │  │ ai.memory_blocks │   │  Core memory                 │
│                    │  │ ai.user_memories │   │  + embeddings (archival)     │
│                    │  │ ai.block_assigns │   │  Block ↔ Agent junction      │
│                    │  └──────────────────┘   │                              │
│                    └─────────────────────────┘                              │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## Three-Tier Memory System

### Tier 1: Core Memory (Always In-Context)

Core memory blocks are **always loaded into the system prompt**. No tool call needed - the agent sees this information on every turn.

| Block Label | Content | Shared With |
|-------------|---------|-------------|
| `zeus_profile` | Identity, preferences, communication style | ALL agents (`*`) |
| `pillar_hp` | HP domain knowledge, active health goals | `legacy`, `hp_agent` |
| `pillar_le` | LE domain knowledge, family context | `legacy`, `le_agent` |
| `pillar_ge` | GE domain knowledge, business context | `legacy`, `ge_agent` |
| `pillar_it` | IT domain knowledge, infra status | `legacy`, `it_agent` |
| `onemind_system` | System status, active services | ALL agents (`*`) |

**Self-editing tools:**
- `memory_insert(label, content)` - Add to existing block
- `memory_replace(label, old, new)` - Replace specific text
- `memory_rethink(label, new_content)` - Full block rewrite

### Tier 2: Archival Memory (Vector Search)

Long-term storage with semantic search. Agent calls `archival_memory_search(query)` to retrieve relevant memories.

| Content Type | Example | When Used |
|--------------|---------|-----------|
| Decisions | "Chose PostgreSQL over MongoDB for..." | When making similar decisions |
| Learnings | "Zeus prefers bullet points over..." | When formatting responses |
| Facts | "Zeus's birthday is March 15" | When relevant context needed |
| Patterns | "Every Monday Zeus reviews tasks at 9am" | When planning schedules |
| Project context | "The inbox feature uses WebSocket..." | When working on projects |

**Tools:**
- `archival_memory_insert(content, tags, pillar)` - Store with embedding
- `archival_memory_search(query, pillar, tags, top_k)` - Semantic retrieval

### Tier 3: Recall Memory (Conversation History)

Recent conversation history. Managed automatically by Agno.

| Scope | Retention | Purpose |
|-------|-----------|---------|
| Session | Full history | Current conversation context |
| Cross-session | Last N runs | Continuity between sessions |
| Summaries | Compressed | Long-term conversation essence |

---

## Agent Hierarchy

```
                    ┌─────────────────────┐
                    │       LEGACY        │
                    │   (Mother Agent)    │
                    │                     │
                    │ - Coordination      │
                    │ - Cross-pillar      │
                    │ - Approvals (HITL)  │
                    │ - System ops        │
                    └──────────┬──────────┘
                               │
         ┌─────────────────────┼─────────────────────┐
         │                     │                     │
         ▼                     ▼                     ▼
┌─────────────────┐  ┌─────────────────┐  ┌─────────────────┐
│   HP Agent      │  │   LE Agent      │  │   GE Agent      │
│ (10-19 Tasks)   │  │ (20-29 Tasks)   │  │ (30-39 Tasks)   │
│                 │  │                 │  │                 │
│ - Health        │  │ - Family        │  │ - Business      │
│ - Fitness       │  │ - Home          │  │ - Revenue       │
│ - Habits        │  │ - Homestead     │  │ - Products      │
│ - Learning      │  │ - Heritage      │  │ - Operations    │
└─────────────────┘  └─────────────────┘  └─────────────────┘
         │
         ▼
┌─────────────────┐
│   Grid Agent      │
│ (40-49 Tasks)   │
│                 │
│ - Infra         │
│ - Automation    │
│ - AI systems    │
│ - Security      │
└─────────────────┘
```

---

## Memory Block Sharing Patterns

### Pattern 1: Universal Blocks (All Agents)

```python
# Zeus profile - every agent needs to know who they serve
await manager.create_block(
    label="zeus_profile",
    value=ZEUS_PROFILE_CONTENT,
    shared_with_agents=["*"],  # ALL agents
    read_only=False,  # Agents can suggest updates
)
```

### Pattern 2: Pillar-Specific Blocks

```python
# HP domain knowledge - only HP-related agents
await manager.create_block(
    label="pillar_hp",
    value=HP_DOMAIN_KNOWLEDGE,
    shared_with_agents=["legacy", "hp_agent"],
    pillar="HP",
)
```

### Pattern 3: Team Blocks

```python
# Dev team context - shared across dev team agents
await manager.create_block(
    label="team_dev",
    value=DEV_TEAM_CONTEXT,
    shared_with_teams=["dev_team"],
)
```

### Pattern 4: Dynamic Assignment (Button Press)

```python
# User clicks "Add to HP Agent" in UI
await manager.attach_block_to_agent(
    block_id="block-uuid",
    agent_id="hp_agent"
)

# User clicks "Remove from HP Agent"
await manager.detach_block_from_agent(
    block_id="block-uuid",
    agent_id="hp_agent"
)
```

---

## Code Modules to Create

| Module | Purpose | Lines | Priority |
|--------|---------|-------|----------|
| `core/memory/blocks.py` | MemoryBlockManager class | ~400 | P0 |
| `core/memory/tools.py` | SelfEditingMemoryTools toolkit | ~350 | P0 |
| `core/memory/embeddings.py` | OpenAI embedding service | ~150 | P0 |
| `core/orchestrator/memory_api.py` | REST endpoints for blocks | ~300 | P1 |
| `core/agent/pillar_agents.py` | HP, LE, GE, IT specialist agents | ~500 | P1 |

---

## Implementation Roadmap (5 Days)

| Day | Task | Files | LOC |
|-----|------|-------|-----|
| **Day 1** | Database schema + migrations | `init-db.sql` | ~100 |
| **Day 2** | MemoryBlockManager + Embeddings | `core/memory/blocks.py`, `embeddings.py` | ~550 |
| **Day 3** | SelfEditingMemoryTools | `core/memory/tools.py` | ~350 |
| **Day 4** | API endpoints + integrate with Legacy | `memory_api.py`, `legacy.py` | ~400 |
| **Day 5** | Pillar agents (HP, LE, GE, IT) | `pillar_agents.py` | ~500 |

**Total: ~1,900 lines of new code**

---

## What We Take From Letta (Patterns Only)

| Pattern | Source | Our Implementation |
|---------|--------|-------------------|
| Core memory blocks | Letta `blocks` table | `ai.memory_blocks` table |
| Archival embeddings | Letta `passages` table | `embedding` column in `ai.user_memories` |
| `memory_insert` | Letta tool | Our `SelfEditingMemoryTools` |
| `memory_replace` | Letta tool | Our `SelfEditingMemoryTools` |
| `memory_rethink` | Letta tool | Our `SelfEditingMemoryTools` |
| Block sharing | `shared_with_agents` array | Same pattern |
| Attach/detach API | Letta REST endpoints | Our `memory_api.py` |

**We don't take:**
- Letta server container (not needed)
- Letta's database schema (use our own)
- Letta's SDK (use Agno)
- Letta's ADE UI (use our Flutter UI)

---

## Cost Analysis

| Item | One-Time | Monthly |
|------|----------|---------|
| New containers | $0 | $0 |
| OpenAI embeddings (text-embedding-3-small) | $0 | ~$1-5 |
| Additional PostgreSQL storage | $0 | ~$1 |
| **Total** | **$0** | **~$2-6** |

---

## API Endpoints (Future)

### Block Management

| Method | Endpoint | Purpose |
|--------|----------|---------|
| `GET` | `/api/memory/blocks` | List all blocks for user |
| `POST` | `/api/memory/blocks` | Create new block |
| `GET` | `/api/memory/blocks/{id}` | Get block details |
| `PATCH` | `/api/memory/blocks/{id}` | Update block content |
| `DELETE` | `/api/memory/blocks/{id}` | Delete block |
| `POST` | `/api/memory/blocks/{id}/attach` | Attach block to agent |
| `POST` | `/api/memory/blocks/{id}/detach` | Detach block from agent |
| `GET` | `/api/memory/agents/{id}/blocks` | Get blocks for agent |

### Archival Memory

| Method | Endpoint | Purpose |
|--------|----------|---------|
| `POST` | `/api/memory/archival` | Insert with embedding |
| `GET` | `/api/memory/archival/search` | Semantic search |

---

## Flutter UI Changes

### Memory Screen Enhancement

```
┌─────────────────────────────────────────────────────────────┐
│  Memory                                          [+ Block]  │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│  CORE BLOCKS                                                │
│  ┌─────────────────────────────────────────────────────┐   │
│  │ zeus_profile                      [Edit] [Agents ▼] │   │
│  │ Identity, preferences, communication style...       │   │
│  │ Shared: ALL AGENTS                                  │   │
│  └─────────────────────────────────────────────────────┘   │
│  ┌─────────────────────────────────────────────────────┐   │
│  │ pillar_hp                         [Edit] [Agents ▼] │   │
│  │ Health domain knowledge...                          │   │
│  │ Shared: legacy, hp_agent                            │   │
│  └─────────────────────────────────────────────────────┘   │
│                                                             │
│  ARCHIVAL MEMORIES                           [Search 🔍]   │
│  ┌─────────────────────────────────────────────────────┐   │
│  │ "Chose PostgreSQL over MongoDB..."     [HP] [2d ago]│   │
│  │ "Zeus prefers bullet points..."        [--] [1w ago]│   │
│  │ "Morning briefing pattern: 9am..."     [GE] [2w ago]│   │
│  └─────────────────────────────────────────────────────┘   │
│                                                             │
└─────────────────────────────────────────────────────────────┘
```

### Agent Assignment Dialog

```
┌─────────────────────────────────────────┐
│  Assign Block: zeus_profile             │
├─────────────────────────────────────────┤
│                                         │
│  ☑ Legacy (Mother)                      │
│  ☑ HP Agent                             │
│  ☑ LE Agent                             │
│  ☑ GE Agent                             │
│  ☑ Grid Agent                             │
│  ☐ Research Agent                       │
│  ☐ Code Agent                           │
│                                         │
│  [Cancel]                    [Save]     │
└─────────────────────────────────────────┘
```

---

## Success Criteria

| Metric | Target |
|--------|--------|
| Agents can share memory blocks | ✅ |
| Agents can self-edit their memory | ✅ |
| UI can assign blocks to agents | ✅ |
| Semantic search works (<200ms) | ✅ |
| No new containers required | ✅ |
| Monthly cost < $10 | ✅ |

---

## References

- [LETTA-AGNO-MEMORY-ARCHITECTURE.md](../../LETTA-AGNO-MEMORY-ARCHITECTURE.md) - Full research document
- [LETTA-RESEARCH-FINDINGS.md](./LETTA-RESEARCH-FINDINGS.md) - Letta deep dive
- [LETTA-AGNO-TEAMS-WORKFLOWS.md](./LETTA-AGNO-TEAMS-WORKFLOWS.md) - Teams integration
- [Letta Documentation](https://docs.letta.com) - Official Letta docs
- [Agno Documentation](https://docs.agno.com) - Official Agno docs
