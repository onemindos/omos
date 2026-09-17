# Bridge v2 Architecture - Simplified

## Overview

Bridge v2 simplifies the LifeOS-to-Learning Machine connection by removing the entity wrapper layer. **Entity operations now go directly to the Learning Machine**, while Bridge focuses purely on analytics and intelligence.

### What Changed

| Before (v1) | After (v2) |
|-------------|------------|
| Bridge wraps entity registration | Learning Machine handles directly |
| Bridge wraps entity linking | Learning Machine handles directly |
| Bridge wraps fact addition | Learning Machine handles directly |
| ~1,100 line service | ~450 line service |
| Duplicate entity storage logic | Single source of truth |

### What Bridge v2 Keeps

1. **Intelligence Scores** - PI/LI/WI calculation from LifeOS data
2. **Pattern Detection** - Habit/task pattern analysis
3. **Day Predictions** - Energy, completion likelihood, blockers
4. **Unified Context** - Aggregates LifeOS + Learning data for AI
5. **Family Features** - Family member overview

---

## Architecture Diagram

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                              ONEMIND OS                                      │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  ┌────────────────────────────────────────────────────────────────────────┐ │
│  │                        FLUTTER FRONTEND                                │ │
│  │  ┌────────────┐  ┌────────────┐  ┌────────────┐  ┌────────────────┐  │ │
│  │  │   Tasks    │  │   Habits   │  │   Goals    │  │   Dashboard    │  │ │
│  │  └─────┬──────┘  └─────┬──────┘  └─────┬──────┘  └───────┬────────┘  │ │
│  └────────┼───────────────┼───────────────┼─────────────────┼───────────┘ │
│           │               │               │                 │              │
│           ▼               ▼               ▼                 ▼              │
│  ┌────────────────────────────────────────────────────────────────────────┐ │
│  │                          REST API LAYER                                │ │
│  │                                                                        │ │
│  │  ┌─────────────────────────┐     ┌──────────────────────────────────┐ │ │
│  │  │   /api/lifeos/*         │     │   /api/learning/*                │ │ │
│  │  │   - Tasks CRUD          │     │   - Entity CRUD                  │ │ │
│  │  │   - Habits CRUD         │     │   - Relations CRUD               │ │ │
│  │  │   - Goals CRUD          │     │   - Facts CRUD                   │ │ │
│  │  │   - Journal CRUD        │     │   - Search/Query                 │ │ │
│  │  └───────────┬─────────────┘     └───────────┬──────────────────────┘ │ │
│  │              │                               │                         │ │
│  │              │     ┌─────────────────────────┴─────────┐              │ │
│  │              │     │   /api/bridge/*  (ANALYTICS)      │              │ │
│  │              │     │   - Intelligence Scores           │              │ │
│  │              │     │   - Pattern Detection             │              │ │
│  │              │     │   - Day Predictions               │              │ │
│  │              │     │   - Unified Context               │              │ │
│  │              │     └───────────────────────────────────┘              │ │
│  └──────────────┼──────────────────────────────────────────────────────┘ │
│                 │                                                         │
│                 ▼                                                         │
│  ┌────────────────────────────────────────────────────────────────────────┐ │
│  │                        BACKEND SERVICES                                │ │
│  │                                                                        │ │
│  │  ┌──────────────────┐  ┌──────────────────┐  ┌───────────────────┐   │ │
│  │  │  LifeOS Services │  │ Learning Machine │  │  Bridge Service   │   │ │
│  │  │  ────────────────│  │  ───────────────│  │  ───────────────  │   │ │
│  │  │  • TaskService   │  │  • Entities     │  │  • Scores (PI/LI/ │   │ │
│  │  │  • HabitService  │──▶  • Relations    │◀─│    WI)            │   │ │
│  │  │  • GoalService   │  │  • Facts        │  │  • Patterns       │   │ │
│  │  │  • JournalService│  │  • Search       │  │  • Predictions    │   │ │
│  │  └────────┬─────────┘  └────────┬────────┘  └─────────┬─────────┘   │ │
│  │           │                     │                     │              │ │
│  │           │   ┌─────────────────┴─────────────────────┘              │ │
│  │           │   │                                                      │ │
│  │           │   ▼  TRIGGERS (on create/update/complete)                │ │
│  │           │   ┌───────────────────────────────────────────────────┐  │ │
│  │           │   │  LifeOSTriggers                                   │  │ │
│  │           │   │  • on_task_created() ──▶ learning.create_entity() │  │ │
│  │           │   │  • on_habit_completed() ──▶ learning.add_fact()   │  │ │
│  │           │   │  • on_goal_progress() ──▶ learning.add_fact()     │  │ │
│  │           └───┤  (Calls Learning Machine DIRECTLY)                │  │ │
│  │               └───────────────────────────────────────────────────┘  │ │
│  └───────────────────────────────────────────────────────────────────────┘ │
│                                                                             │
│  ┌───────────────────────────────────────────────────────────────────────┐ │
│  │                          POSTGRESQL                                   │ │
│  │  ┌─────────────────┐  ┌─────────────────┐  ┌─────────────────────┐   │ │
│  │  │  LIFEOS TABLES  │  │ LEARNING TABLES │  │   AI TABLES         │   │ │
│  │  │  • inbox_items  │  │ • learning_     │  │   • ai.user_        │   │ │
│  │  │  • habits       │  │   entities      │  │     memories        │   │ │
│  │  │  • goals        │  │ • learning_     │  │   • letta_          │   │ │
│  │  │  • projects     │  │   entity_facts  │  │     messages        │   │ │
│  │  │  • journal_     │  │ • learning_     │  │                     │   │ │
│  │  │    entries      │  │   entity_       │  │                     │   │ │
│  │  │                 │  │   relations     │  │                     │   │ │
│  │  └─────────────────┘  └─────────────────┘  └─────────────────────┘   │ │
│  └───────────────────────────────────────────────────────────────────────┘ │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## PostgreSQL Schema

### LifeOS Core Tables

```sql
-- =============================================================================
-- TASKS (Universal Inbox)
-- =============================================================================
CREATE TABLE inbox_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    -- Core fields
    type VARCHAR(50) NOT NULL,              -- task, message, notification, approval, etc.
    priority VARCHAR(20) DEFAULT 'normal',  -- critical, high, normal, low
    status VARCHAR(20) DEFAULT 'new',       -- new, in_progress, done, archived

    -- Content
    title VARCHAR(500) NOT NULL,
    body TEXT,
    source VARCHAR(100),                    -- github, todoist, system, chat, manual

    -- Assignment & Dispatch
    assignee_type VARCHAR(20) DEFAULT 'auto',   -- human, agent, team
    assignee_id VARCHAR(100),
    assigned_at TIMESTAMPTZ,
    assigned_by VARCHAR(100),

    -- Task-specific
    task_pillar VARCHAR(10),                -- hp, le, ge (Three Pillars)
    start_date DATE,                        -- Due/planned date
    completed_at TIMESTAMPTZ,
    estimated_minutes INT,
    context VARCHAR(20),                    -- GTD context: home, work, errands, etc.

    -- Organization
    project_id UUID REFERENCES projects(id),
    user_id VARCHAR(50) DEFAULT 'zeus',

    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- Indexes for common queries
CREATE INDEX idx_inbox_user_status ON inbox_items(user_id, status);
CREATE INDEX idx_inbox_pillar ON inbox_items(task_pillar);
CREATE INDEX idx_inbox_project ON inbox_items(project_id);
CREATE INDEX idx_inbox_date ON inbox_items(start_date);


-- =============================================================================
-- HABITS
-- =============================================================================
CREATE TABLE habits (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id VARCHAR(50) DEFAULT 'zeus',
    name VARCHAR(255) NOT NULL,
    description TEXT,
    emoji VARCHAR(10),
    pillar VARCHAR(10),                         -- hp, le, ge

    -- Scheduling
    frequency VARCHAR(50) DEFAULT 'daily',      -- daily, weekdays, weekly, custom
    frequency_days INT[] DEFAULT '{1,2,3,4,5,6,7}',
    target_time TIME,

    -- Tracking
    current_streak INT DEFAULT 0,
    best_streak INT DEFAULT 0,
    total_completions INT DEFAULT 0,
    completion_rate DECIMAL(5,4),               -- 0.0000 to 1.0000
    grade CHAR(1),                              -- A, B, C, D, F

    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE habit_completions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    habit_id UUID REFERENCES habits(id) ON DELETE CASCADE,
    user_id VARCHAR(50) DEFAULT 'zeus',
    completion_date DATE NOT NULL,
    status VARCHAR(20) DEFAULT 'completed',     -- completed, skipped, missed
    completed_at TIMESTAMPTZ,

    -- Context for pattern analysis
    mood INT CHECK (mood BETWEEN 1 AND 9),
    energy INT CHECK (energy BETWEEN 1 AND 9),
    notes TEXT,

    UNIQUE(habit_id, completion_date)
);


-- =============================================================================
-- GOALS
-- =============================================================================
CREATE TABLE goals (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id VARCHAR(50) DEFAULT 'zeus',
    name VARCHAR(255) NOT NULL,
    description TEXT,
    emoji VARCHAR(10),
    pillar VARCHAR(10),

    -- Type and tracking
    goal_type VARCHAR(20) DEFAULT 'target',     -- target, habit, milestone
    target_date DATE,
    target_value DECIMAL,
    current_value DECIMAL DEFAULT 0,
    progress INT DEFAULT 0,                     -- 0-100%
    unit VARCHAR(50),

    -- Hierarchy
    parent_goal_id UUID REFERENCES goals(id),
    status VARCHAR(20) DEFAULT 'active',

    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);


-- =============================================================================
-- PROJECTS
-- =============================================================================
CREATE TABLE projects (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id VARCHAR(50) DEFAULT 'zeus',
    name VARCHAR(255) NOT NULL,
    description TEXT,
    emoji VARCHAR(10),
    color VARCHAR(7) DEFAULT '#6366f1',
    pillar VARCHAR(10),

    -- Tracking
    status VARCHAR(20) DEFAULT 'active',
    task_count INT DEFAULT 0,
    completed_count INT DEFAULT 0,

    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);


-- =============================================================================
-- JOURNAL
-- =============================================================================
CREATE TABLE journal_entries (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id VARCHAR(50) DEFAULT 'zeus',
    entry_date DATE NOT NULL UNIQUE,
    content TEXT,

    -- Tracking
    mood INT CHECK (mood BETWEEN 1 AND 9),
    energy INT CHECK (energy BETWEEN 1 AND 9),
    tags TEXT[],

    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);
```

### Learning Machine Tables

```sql
-- =============================================================================
-- LEARNING ENTITIES (Knowledge Graph Nodes)
-- =============================================================================
CREATE TABLE learning_entities (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id VARCHAR(255) NOT NULL,
    type VARCHAR(50) NOT NULL,                  -- person, project, company, location, concept
    name VARCHAR(255) NOT NULL,
    description TEXT,
    pillar VARCHAR(10),                         -- hp, le, ge
    metadata JSONB DEFAULT '{}',                -- Flexible storage for lifeos_id, etc.

    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- Bridge stores LifeOS references in metadata:
-- {
--   "bridge_entity_type": "task",              -- task, habit, goal, project, etc.
--   "lifeos_id": "uuid-of-task",
--   "lifeos_table": "inbox_items",
--   "priority": "high",
--   "status": "in_progress"
-- }

CREATE INDEX idx_entities_user ON learning_entities(user_id);
CREATE INDEX idx_entities_type ON learning_entities(type);
CREATE INDEX idx_entities_metadata_lifeos ON learning_entities((metadata->>'lifeos_id'));


-- =============================================================================
-- LEARNING ENTITY FACTS (What AI learns about entities)
-- =============================================================================
CREATE TABLE learning_entity_facts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    entity_id UUID REFERENCES learning_entities(id) ON DELETE CASCADE,
    fact TEXT NOT NULL,
    source VARCHAR(50) DEFAULT 'system',        -- user, system, inferred, imported
    confidence FLOAT DEFAULT 1.0,               -- 0.0 to 1.0

    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- Example facts:
-- "Completed on Monday mood=7/9 energy=8/9"
-- "Completed 15 minutes under estimate"
-- "Progress increased from 40% to 60% (+20%)"
-- "Skipped on Saturday reason: traveling"

CREATE INDEX idx_facts_entity ON learning_entity_facts(entity_id);


-- =============================================================================
-- LEARNING ENTITY RELATIONS (Knowledge Graph Edges)
-- =============================================================================
CREATE TABLE learning_entity_relations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    from_entity_id UUID REFERENCES learning_entities(id) ON DELETE CASCADE,
    to_entity_id UUID REFERENCES learning_entities(id) ON DELETE CASCADE,
    relation_type VARCHAR(100) NOT NULL,        -- part_of, supports, requires, child_of, etc.
    metadata JSONB DEFAULT '{}',

    created_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(from_entity_id, to_entity_id, relation_type)
);

-- Relation types:
-- part_of: Task is part of Project
-- child_of: Sub-goal is child of parent Goal
-- supports: Habit supports Goal
-- requires: Task requires another Task
-- assigned_to: Task assigned to Person

CREATE INDEX idx_relations_from ON learning_entity_relations(from_entity_id);
CREATE INDEX idx_relations_to ON learning_entity_relations(to_entity_id);
```

---

## API Endpoints

### Bridge v2 (Analytics Only)

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/api/bridge/context/{type}/{id}` | GET | Unified context for AI (LifeOS + facts + relations) |
| `/api/bridge/scores` | GET | Intelligence scores (PI/LI/WI) |
| `/api/bridge/patterns` | GET | Pattern detection (habit/task analysis) |
| `/api/bridge/predictions/{date}` | GET | Day predictions (energy, likelihood, blockers) |
| `/api/bridge/family` | GET | Family member overview |
| `/api/bridge/stats` | GET | Bridge statistics |
| `/api/bridge/backfill` | POST | Bulk registration of existing LifeOS items |

### Learning Machine (Entity Operations)

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/api/learning/entities` | POST | Create entity |
| `/api/learning/entities/{id}` | GET | Get entity |
| `/api/learning/entities/{id}` | PUT | Update entity |
| `/api/learning/entities/{id}` | DELETE | Delete entity |
| `/api/learning/entities/{id}/facts` | POST | Add fact |
| `/api/learning/relations` | POST | Create relation |
| `/api/learning/search` | GET | Search entities |

---

## Data Flow

### 1. Task Created → Learning Entity

```python
# In TaskService.create()
task = await self.create(data)

# Trigger registers in Learning Machine
await LifeOSTriggers.on_task_created(task, user_id)
```

```sql
-- Creates Learning entity with metadata pointing to LifeOS item
INSERT INTO learning_entities (user_id, type, name, description, pillar, metadata)
VALUES (
    'zeus',
    'concept',
    'Fix login bug',
    'Users cannot log in on mobile',
    'ge',
    '{
        "bridge_entity_type": "task",
        "lifeos_id": "550e8400-e29b-41d4-a716-446655440000",
        "lifeos_table": "inbox_items",
        "priority": "high",
        "status": "in_progress"
    }'
);
```

### 2. Task Completed → Add Fact

```python
# In TaskService.complete()
await self.update_status(task_id, "done")

# Trigger adds completion fact
await LifeOSTriggers.on_task_completed(task, user_id, actual_minutes=45)
```

```sql
-- Adds fact to the entity
INSERT INTO learning_entity_facts (entity_id, fact, source, confidence)
VALUES (
    '..entity-uuid..',
    'Completed 15 minutes under estimate',
    'system',
    1.0
);
```

### 3. AI Requests Context

```python
# Agent calls Bridge for rich context
context = await bridge.get_entity_context(user_id, "task", task_id)
```

Returns:
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "type": "task",
  "name": "Fix login bug",
  "pillar": "ge",
  "lifeos_data": {
    "status": "in_progress",
    "priority": "high",
    "due_date": "2024-01-30",
    "estimated_minutes": 60
  },
  "relations": [
    {"type": "project", "name": "Mobile App v2", "relation": "part_of"}
  ],
  "facts": [
    {"fact": "Similar bug fixed last week", "source": "system"}
  ],
  "markdown": "# Fix login bug\n**Type:** task | **Pillar:** GE\n..."
}
```

---

## Files Changed

| File | Lines | Change |
|------|-------|--------|
| `service_v2.py` | ~450 | NEW - Simplified service (analytics only) |
| `triggers_v2.py` | ~350 | NEW - Direct Learning Machine calls |
| `api_v2.py` | ~150 | NEW - Analytics-only endpoints |
| `service.py` | 1,119 | OLD - Can be deprecated |
| `triggers.py` | 732 | OLD - Can be deprecated |
| `api.py` | 409 | OLD - Can be deprecated |

**Total reduction: ~1,800 lines removed, ~950 lines new = ~850 lines saved (47% smaller)**

---

## Migration Steps

1. **Import new modules**:
   ```python
   from backend.lifeos.bridge.service_v2 import get_bridge_service_v2
   from backend.lifeos.bridge.triggers_v2 import LifeOSTriggers
   ```

2. **Update LifeOS services to use new triggers**:
   ```python
   # TaskService
   task = await self.create(data)
   await LifeOSTriggers.on_task_created(task, user_id)
   ```

3. **Update API router** (if using new API):
   ```python
   from backend.lifeos.bridge.api_v2 import router as bridge_router
   app.include_router(bridge_router)
   ```

4. **Entity operations go to Learning Machine**:
   ```python
   # Instead of bridge.register_entity()
   from backend.agno.learning.service import get_learning_service_async
   learning = await get_learning_service_async()
   await learning.create_entity(user_id, EntityCreate(...))
   ```

---

## Intelligence Scores Reference

| Score | Pillar | Components |
|-------|--------|------------|
| **PI** (Performance) | HP | health_habits (40%), energy_mood (20%), focus (20%), goals (20%) |
| **LI** (Legacy) | LE | family_habits (35%), home (20%), goals (20%), legacy (25%) |
| **WI** (Wealth) | GE | business_tasks (30%), business_goals (30%), work_habits (20%), finance (20%) |

---

## Three Pillars

| Pillar | Code | Focus | Examples |
|--------|------|-------|----------|
| **Holistic Performance** | HP | Personal health & productivity | Exercise, sleep, meditation, focus |
| **Legacy Evolution** | LE | Family & long-term impact | Kids, home, family time, legacy projects |
| **Generational Empire** | GE | Business & wealth | Work tasks, business goals, finances |


---

## Related

- [[00-24 SO (Sovereign Operations)/02 Infrastructure (VPS, Servers & Networking) (The Body)/02 Infrastructure|02 Infrastructure]]
