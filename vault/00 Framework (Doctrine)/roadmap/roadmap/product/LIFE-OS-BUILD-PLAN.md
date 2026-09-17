# OneMind Life OS - Build Plan

> **Status:** Ready for Implementation
> **Coordinator:** Lead Agent (Mother)
> **Agents:** 6 parallel workers
> **Estimated:** 5-6 days with parallel execution
> **Created:** 2026-01-26

---

## Overview

Build a complete internal Life OS that replaces external services like Todoist. Key principles:
- **Data ownership** - All data stored internally, no external dependencies
- **AI-native dispatch** - Tap any task → dispatch to agent/team/workflow
- **Offline-first** - Edge locations work offline, sync when connected
- **Single user** - Optimized for Zeus, not multi-tenant

### Core Features

| Category | Features |
|----------|----------|
| **Tasks** | Projects, subtasks, dependencies, priorities, tags, contexts, dispatch to AI |
| **Habits** | Streaks, A-F scoring, reminders, pillar tracking |
| **Routines** | Morning/evening flows, step-by-step guidance |
| **Planner** | Time blocks, drag-drop tasks, section views |
| **Pomodoro** | Focus sessions, work modes (sit/stand/walk), focus levels |
| **Calendar** | Week/month views, all events unified, CalDAV sync |
| **Goals** | Target tracking, progress visualization |
| **Journal** | Daily entries, mood/energy tracking |
| **Scheduling** | Agent self-scheduling, cron expressions, one-time triggers |
| **Meals** | Mealie integration for recipes and meal planning |
| **Offline Sync** | Edge SQLite + NATS sync to cloud PostgreSQL |

---

## Architecture

### System Diagram

```
┌─────────────────────────────────────────────────────────────────────────┐
│                           OneMind Life OS                               │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│  ┌─────────────┐    ┌─────────────┐    ┌─────────────┐                 │
│  │   Flutter   │    │   Flutter   │    │   Flutter   │                 │
│  │   Web/iOS   │    │   macOS     │    │   Edge      │                 │
│  │  (Online)   │    │  (Hybrid)   │    │  (Offline)  │                 │
│  └──────┬──────┘    └──────┬──────┘    └──────┬──────┘                 │
│         │                  │                   │                        │
│         │         ┌────────┴────────┐          │                        │
│         │         │  Local SQLite   │          │                        │
│         │         │  (offline-first)│          │                        │
│         │         └────────┬────────┘          │                        │
│         │                  │                   │                        │
│         ▼                  ▼                   ▼                        │
│  ┌─────────────────────────────────────────────────────────────────┐   │
│  │                        NATS JetStream                            │   │
│  │  (Cloud Hub + Edge Leafs - Offline Queue + Sync)                │   │
│  └─────────────────────────────────────────────────────────────────┘   │
│         │                                                               │
│         ▼                                                               │
│  ┌─────────────────────────────────────────────────────────────────┐   │
│  │                    FastAPI Orchestrator                          │   │
│  │  ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌──────────┐            │   │
│  │  │  Tasks   │ │  Habits  │ │  Planner │ │ Schedule │            │   │
│  │  │  Module  │ │  Module  │ │  Module  │ │  Module  │            │   │
│  │  └──────────┘ └──────────┘ └──────────┘ └──────────┘            │   │
│  │  ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌──────────┐            │   │
│  │  │ Routines │ │ Pomodoro │ │ Calendar │ │  Goals   │            │   │
│  │  │  Module  │ │  Module  │ │  Module  │ │  Module  │            │   │
│  │  └──────────┘ └──────────┘ └──────────┘ └──────────┘            │   │
│  │  ┌──────────┐ ┌──────────┐ ┌──────────┐                         │   │
│  │  │ Journal  │ │  Mealie  │ │  Sync    │                         │   │
│  │  │  Module  │ │ Integr.  │ │  Engine  │                         │   │
│  │  └──────────┘ └──────────┘ └──────────┘                         │   │
│  └─────────────────────────────────────────────────────────────────┘   │
│         │                                                               │
│         ▼                                                               │
│  ┌─────────────────────────────────────────────────────────────────┐   │
│  │                       PostgreSQL                                 │   │
│  │  (Cloud - Source of Truth)                                       │   │
│  └─────────────────────────────────────────────────────────────────┘   │
│         │                                                               │
│         ▼                                                               │
│  ┌─────────────────────────────────────────────────────────────────┐   │
│  │                    Agent Tools (Agno)                            │   │
│  │  create_task, complete_habit, start_routine, schedule_recurring │   │
│  │  dispatch_task, get_today_plan, log_pomodoro, suggest_meal      │   │
│  └─────────────────────────────────────────────────────────────────┘   │
│                                                                         │
└─────────────────────────────────────────────────────────────────────────┘

Edge Deployment (Offline-First):
┌───────────────────────────────────────────────────────────────┐
│  Edge Device (Jetson/RPi)                                     │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐        │
│  │ NATS Leaf    │  │ SQLite       │  │ Sync Agent   │        │
│  │ (queue while │  │ (local copy) │  │ (reconcile)  │        │
│  │  offline)    │  │              │  │              │        │
│  └──────────────┘  └──────────────┘  └──────────────┘        │
│         │                 │                 │                 │
│         └────────────────┬┴─────────────────┘                 │
│                          │                                    │
│                          ▼                                    │
│                   When online:                                │
│                   Sync to Cloud PostgreSQL                    │
│                   via NATS JetStream                          │
└───────────────────────────────────────────────────────────────┘
```

### Tech Stack

| Layer | Technology | Purpose |
|-------|------------|---------|
| **Cloud DB** | PostgreSQL + pgvector | Source of truth |
| **Edge DB** | SQLite | Offline storage |
| **Sync** | NATS JetStream | Event-driven sync with offline queue |
| **Backend** | FastAPI + APScheduler | APIs + scheduled execution |
| **Frontend** | Flutter | Cross-platform UI |
| **Calendar** | syncfusion_flutter_calendar | Week/month views |
| **Meals** | Mealie (Docker) | Recipe management |
| **Scheduling** | APScheduler | Cron + one-time triggers |

---

## Database Schema

### New Tables (Add to init-db.sql)

```sql
-- =============================================================================
-- LIFE OS TABLES
-- =============================================================================

-- Projects (group related tasks)
CREATE TABLE IF NOT EXISTS projects (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id VARCHAR(50) DEFAULT 'zeus',
    name VARCHAR(255) NOT NULL,
    description TEXT,
    emoji VARCHAR(10),
    color VARCHAR(7) DEFAULT '#6366f1',
    pillar VARCHAR(10),  -- HP, LE, GE, IT
    mode VARCHAR(20) DEFAULT 'parallel',  -- parallel, sequential
    status VARCHAR(20) DEFAULT 'active',  -- active, completed, archived
    sort_order INT DEFAULT 0,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- Habits with streaks and A-F scoring
CREATE TABLE IF NOT EXISTS habits (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id VARCHAR(50) DEFAULT 'zeus',
    name VARCHAR(255) NOT NULL,
    description TEXT,
    emoji VARCHAR(10),
    pillar VARCHAR(10),
    frequency VARCHAR(50) DEFAULT 'daily',  -- daily, weekdays, weekends, weekly, custom
    frequency_days INT[] DEFAULT '{1,2,3,4,5,6,7}',  -- days of week (1=Mon)
    target_time TIME,
    reminder_enabled BOOLEAN DEFAULT false,
    reminder_minutes_before INT DEFAULT 15,
    current_streak INT DEFAULT 0,
    best_streak INT DEFAULT 0,
    total_completions INT DEFAULT 0,
    list_id UUID,  -- for grouping habits
    sort_order INT DEFAULT 0,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- Habit completions (for streak calculation and A-F scoring)
CREATE TABLE IF NOT EXISTS habit_completions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    habit_id UUID REFERENCES habits(id) ON DELETE CASCADE,
    completed_date DATE NOT NULL,
    completed_at TIMESTAMPTZ DEFAULT NOW(),
    notes TEXT,
    UNIQUE(habit_id, completed_date)
);

-- Routines (morning, evening, workout, etc.)
CREATE TABLE IF NOT EXISTS routines (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id VARCHAR(50) DEFAULT 'zeus',
    name VARCHAR(255) NOT NULL,
    description TEXT,
    emoji VARCHAR(10),
    trigger_type VARCHAR(20) DEFAULT 'manual',  -- manual, time, location, event
    trigger_time TIME,
    trigger_location VARCHAR(100),
    estimated_minutes INT,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- Routine steps
CREATE TABLE IF NOT EXISTS routine_steps (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    routine_id UUID REFERENCES routines(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    duration_minutes INT,
    sort_order INT DEFAULT 0,
    is_optional BOOLEAN DEFAULT false
);

-- Time blocks for planner
CREATE TABLE IF NOT EXISTS time_blocks (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id VARCHAR(50) DEFAULT 'zeus',
    date DATE NOT NULL,
    title VARCHAR(255),
    section VARCHAR(20),  -- morning, afternoon, evening, night
    start_time TIME NOT NULL,
    end_time TIME NOT NULL,
    color VARCHAR(7) DEFAULT '#6366f1',
    is_recurring BOOLEAN DEFAULT false,
    recurrence_rule VARCHAR(255),  -- RRULE format
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- Task assignments to time blocks
CREATE TABLE IF NOT EXISTS task_assignments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    task_id UUID REFERENCES inbox_items(id) ON DELETE CASCADE,
    block_id UUID REFERENCES time_blocks(id) ON DELETE CASCADE,
    sort_order INT DEFAULT 0,
    assigned_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(task_id, block_id)
);

-- Pomodoro focus sessions
CREATE TABLE IF NOT EXISTS pomodoro_sessions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id VARCHAR(50) DEFAULT 'zeus',
    task_id UUID REFERENCES inbox_items(id),
    started_at TIMESTAMPTZ NOT NULL,
    ended_at TIMESTAMPTZ,
    planned_duration INT DEFAULT 25,  -- minutes
    actual_duration INT,  -- minutes
    category VARCHAR(50),  -- work, study, creative, admin
    focus_level INT CHECK (focus_level BETWEEN 1 AND 5),
    work_mode VARCHAR(20),  -- sit, stand, walk
    interruptions INT DEFAULT 0,
    notes TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- Scheduled tasks (agent self-scheduling)
CREATE TABLE IF NOT EXISTS scheduled_tasks (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id VARCHAR(50) DEFAULT 'zeus',
    name VARCHAR(255) NOT NULL,
    description TEXT,
    cron_expression VARCHAR(100),  -- for recurring
    one_time_at TIMESTAMPTZ,  -- for one-time
    execution_type VARCHAR(20) NOT NULL,  -- agent, team, workflow
    execution_target VARCHAR(255) NOT NULL,  -- agent_id, team_id, workflow_name
    execution_params JSONB DEFAULT '{}',
    context TEXT,  -- additional context for the agent
    created_by VARCHAR(50) DEFAULT 'zeus',  -- could be an agent name
    is_active BOOLEAN DEFAULT true,
    last_run_at TIMESTAMPTZ,
    next_run_at TIMESTAMPTZ,
    run_count INT DEFAULT 0,
    last_result JSONB,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- Goals with progress tracking
CREATE TABLE IF NOT EXISTS goals (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id VARCHAR(50) DEFAULT 'zeus',
    name VARCHAR(255) NOT NULL,
    description TEXT,
    emoji VARCHAR(10),
    pillar VARCHAR(10),
    goal_type VARCHAR(20) DEFAULT 'target',  -- target, habit, milestone
    target_date DATE,
    target_value DECIMAL,
    current_value DECIMAL DEFAULT 0,
    unit VARCHAR(50),
    status VARCHAR(20) DEFAULT 'active',  -- active, completed, abandoned
    parent_goal_id UUID REFERENCES goals(id),
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- Journal entries with mood/energy
CREATE TABLE IF NOT EXISTS journal_entries (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id VARCHAR(50) DEFAULT 'zeus',
    entry_date DATE NOT NULL,
    content TEXT,
    mood INT CHECK (mood BETWEEN 1 AND 9),
    energy INT CHECK (energy BETWEEN 1 AND 9),
    gratitude TEXT[],
    wins TEXT[],
    challenges TEXT[],
    tags TEXT[] DEFAULT '{}',
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(user_id, entry_date)
);

-- Calendar events (synced from CalDAV + internal)
CREATE TABLE IF NOT EXISTS calendar_events (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id VARCHAR(50) DEFAULT 'zeus',
    external_id VARCHAR(255),  -- CalDAV UID
    calendar_source VARCHAR(50) DEFAULT 'internal',  -- internal, google, apple, outlook
    title VARCHAR(255) NOT NULL,
    description TEXT,
    location VARCHAR(255),
    start_time TIMESTAMPTZ NOT NULL,
    end_time TIMESTAMPTZ NOT NULL,
    all_day BOOLEAN DEFAULT false,
    recurrence_rule VARCHAR(255),
    color VARCHAR(7),
    reminder_minutes INT[],
    attendees JSONB DEFAULT '[]',
    status VARCHAR(20) DEFAULT 'confirmed',  -- confirmed, tentative, cancelled
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    synced_at TIMESTAMPTZ,
    UNIQUE(calendar_source, external_id)
);

-- Sync log for offline reconciliation
CREATE TABLE IF NOT EXISTS sync_log (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    table_name VARCHAR(100) NOT NULL,
    record_id UUID NOT NULL,
    operation VARCHAR(20) NOT NULL,  -- insert, update, delete
    data JSONB NOT NULL,
    device_id VARCHAR(100),
    synced BOOLEAN DEFAULT false,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    synced_at TIMESTAMPTZ
);

-- =============================================================================
-- MEALS MODULE (Custom - AI-Powered)
-- =============================================================================

-- Recipes
CREATE TABLE IF NOT EXISTS recipes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id VARCHAR(50) DEFAULT 'zeus',
    name VARCHAR(255) NOT NULL,
    description TEXT,
    source_url VARCHAR(500),
    image_url VARCHAR(500),
    prep_time_minutes INT,
    cook_time_minutes INT,
    total_time_minutes INT,
    servings INT DEFAULT 4,
    difficulty VARCHAR(20),  -- easy, medium, hard
    cuisine VARCHAR(50),
    tags TEXT[] DEFAULT '{}',
    ingredients JSONB NOT NULL DEFAULT '[]',  -- [{name, amount, unit, notes}]
    instructions JSONB NOT NULL DEFAULT '[]', -- [{step, text, time_minutes}]
    nutrition JSONB,  -- {calories, protein, carbs, fat, fiber}
    notes TEXT,
    rating INT CHECK (rating BETWEEN 1 AND 5),
    times_made INT DEFAULT 0,
    last_made_at TIMESTAMPTZ,
    is_favorite BOOLEAN DEFAULT false,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- Meal plans
CREATE TABLE IF NOT EXISTS meal_plans (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id VARCHAR(50) DEFAULT 'zeus',
    date DATE NOT NULL,
    meal_type VARCHAR(20) NOT NULL,  -- breakfast, lunch, dinner, snack
    recipe_id UUID REFERENCES recipes(id) ON DELETE SET NULL,
    custom_meal VARCHAR(255),  -- if not from recipes
    notes TEXT,
    servings INT DEFAULT 1,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(user_id, date, meal_type)
);

-- Shopping lists
CREATE TABLE IF NOT EXISTS shopping_lists (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id VARCHAR(50) DEFAULT 'zeus',
    name VARCHAR(255),
    start_date DATE,
    end_date DATE,
    items JSONB NOT NULL DEFAULT '[]',  -- [{name, amount, unit, category, checked, recipe_id}]
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- Pantry (what you have)
CREATE TABLE IF NOT EXISTS pantry_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id VARCHAR(50) DEFAULT 'zeus',
    name VARCHAR(255) NOT NULL,
    category VARCHAR(50),  -- produce, dairy, meat, pantry, frozen
    quantity DECIMAL,
    unit VARCHAR(20),
    expiry_date DATE,
    location VARCHAR(50),  -- fridge, freezer, pantry
    notes TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_recipes_user ON recipes(user_id);
CREATE INDEX IF NOT EXISTS idx_recipes_tags ON recipes USING GIN(tags);
CREATE INDEX IF NOT EXISTS idx_meal_plans_date ON meal_plans(user_id, date);
CREATE INDEX IF NOT EXISTS idx_pantry_expiry ON pantry_items(expiry_date) WHERE expiry_date IS NOT NULL;

-- =============================================================================
-- INBOX_ITEMS ENHANCEMENTS (ALTER existing table)
-- =============================================================================

ALTER TABLE inbox_items ADD COLUMN IF NOT EXISTS project_id UUID REFERENCES projects(id);
ALTER TABLE inbox_items ADD COLUMN IF NOT EXISTS parent_id UUID REFERENCES inbox_items(id);
ALTER TABLE inbox_items ADD COLUMN IF NOT EXISTS start_date DATE;
ALTER TABLE inbox_items ADD COLUMN IF NOT EXISTS planned_date DATE;
ALTER TABLE inbox_items ADD COLUMN IF NOT EXISTS estimated_minutes INT;
ALTER TABLE inbox_items ADD COLUMN IF NOT EXISTS actual_minutes INT;
ALTER TABLE inbox_items ADD COLUMN IF NOT EXISTS depends_on UUID[];
ALTER TABLE inbox_items ADD COLUMN IF NOT EXISTS context VARCHAR(50) DEFAULT 'personal';  -- personal, work, errands, calls
ALTER TABLE inbox_items ADD COLUMN IF NOT EXISTS is_must_do BOOLEAN DEFAULT false;
ALTER TABLE inbox_items ADD COLUMN IF NOT EXISTS energy_required VARCHAR(20);  -- low, medium, high
ALTER TABLE inbox_items ADD COLUMN IF NOT EXISTS sort_order INT DEFAULT 0;
ALTER TABLE inbox_items ADD COLUMN IF NOT EXISTS tags TEXT[] DEFAULT '{}';
ALTER TABLE inbox_items ADD COLUMN IF NOT EXISTS recurrence_rule VARCHAR(255);
ALTER TABLE inbox_items ADD COLUMN IF NOT EXISTS completed_at TIMESTAMPTZ;
ALTER TABLE inbox_items ADD COLUMN IF NOT EXISTS dispatch_target_type VARCHAR(20);  -- agent, team, workflow
ALTER TABLE inbox_items ADD COLUMN IF NOT EXISTS dispatch_target_id VARCHAR(255);
ALTER TABLE inbox_items ADD COLUMN IF NOT EXISTS dispatch_status VARCHAR(20);  -- pending, running, completed, failed
ALTER TABLE inbox_items ADD COLUMN IF NOT EXISTS dispatch_run_id VARCHAR(255);

-- =============================================================================
-- INDEXES
-- =============================================================================

CREATE INDEX IF NOT EXISTS idx_habits_user_active ON habits(user_id, is_active);
CREATE INDEX IF NOT EXISTS idx_habit_completions_date ON habit_completions(completed_date);
CREATE INDEX IF NOT EXISTS idx_time_blocks_date ON time_blocks(user_id, date);
CREATE INDEX IF NOT EXISTS idx_pomodoro_sessions_date ON pomodoro_sessions(user_id, started_at);
CREATE INDEX IF NOT EXISTS idx_scheduled_tasks_next_run ON scheduled_tasks(next_run_at) WHERE is_active = true;
CREATE INDEX IF NOT EXISTS idx_calendar_events_time ON calendar_events(user_id, start_time, end_time);
CREATE INDEX IF NOT EXISTS idx_sync_log_unsynced ON sync_log(synced) WHERE synced = false;
CREATE INDEX IF NOT EXISTS idx_inbox_items_planned ON inbox_items(user_id, planned_date) WHERE status != 'completed';
CREATE INDEX IF NOT EXISTS idx_inbox_items_project ON inbox_items(project_id);

-- =============================================================================
-- FUNCTIONS
-- =============================================================================

-- Calculate habit grade (A-F based on 30-day completion rate)
CREATE OR REPLACE FUNCTION calculate_habit_grade(habit_uuid UUID)
RETURNS CHAR(1) AS $$
DECLARE
    completion_rate DECIMAL;
    expected_days INT;
    actual_days INT;
BEGIN
    -- Count expected days in last 30 days based on frequency
    SELECT COUNT(*) INTO expected_days
    FROM generate_series(
        CURRENT_DATE - INTERVAL '30 days',
        CURRENT_DATE - INTERVAL '1 day',
        INTERVAL '1 day'
    ) AS d(day)
    WHERE EXTRACT(ISODOW FROM d.day)::INT = ANY(
        SELECT unnest(frequency_days) FROM habits WHERE id = habit_uuid
    );

    -- Count actual completions
    SELECT COUNT(*) INTO actual_days
    FROM habit_completions
    WHERE habit_id = habit_uuid
    AND completed_date >= CURRENT_DATE - INTERVAL '30 days';

    IF expected_days = 0 THEN
        RETURN 'A';
    END IF;

    completion_rate := (actual_days::DECIMAL / expected_days) * 100;

    RETURN CASE
        WHEN completion_rate >= 90 THEN 'A'
        WHEN completion_rate >= 80 THEN 'B'
        WHEN completion_rate >= 70 THEN 'C'
        WHEN completion_rate >= 60 THEN 'D'
        ELSE 'F'
    END;
END;
$$ LANGUAGE plpgsql;

-- Update habit streak on completion
CREATE OR REPLACE FUNCTION update_habit_streak()
RETURNS TRIGGER AS $$
DECLARE
    yesterday_completed BOOLEAN;
    current_streak_val INT;
BEGIN
    -- Check if yesterday was completed
    SELECT EXISTS(
        SELECT 1 FROM habit_completions
        WHERE habit_id = NEW.habit_id
        AND completed_date = NEW.completed_date - INTERVAL '1 day'
    ) INTO yesterday_completed;

    -- Get current streak
    SELECT current_streak INTO current_streak_val
    FROM habits WHERE id = NEW.habit_id;

    IF yesterday_completed THEN
        -- Increment streak
        UPDATE habits SET
            current_streak = current_streak + 1,
            best_streak = GREATEST(best_streak, current_streak + 1),
            total_completions = total_completions + 1,
            updated_at = NOW()
        WHERE id = NEW.habit_id;
    ELSE
        -- Reset streak to 1
        UPDATE habits SET
            current_streak = 1,
            total_completions = total_completions + 1,
            updated_at = NOW()
        WHERE id = NEW.habit_id;
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_update_habit_streak
AFTER INSERT ON habit_completions
FOR EACH ROW
EXECUTE FUNCTION update_habit_streak();
```

---

## Phase Breakdown (6 Agents)

### Agent Assignment

| Agent | Focus Area | Linear Issues | Dependencies |
|-------|------------|---------------|--------------|
| **Agent 1** | Database + Sync Engine | OMOS-323, OMOS-324 | None |
| **Agent 2** | Tasks + Projects + Dispatch | OMOS-325, OMOS-326 | Agent 1 (schema) |
| **Agent 3** | Habits + Routines | OMOS-327, OMOS-328 | Agent 1 (schema) |
| **Agent 4** | Planner + Calendar | OMOS-329, OMOS-330 | Agent 1 (schema) |
| **Agent 5** | Pomodoro + Journal + Goals | OMOS-331, OMOS-332 | Agent 1 (schema) |
| **Agent 6** | Meals (AI) + Agent Tools + Scheduling | OMOS-333, OMOS-334 | Agent 1 (schema) |

### Phase 1: Foundation (Day 1)
**Lead by:** Agent 1

| Task | Agent | Hours |
|------|-------|-------|
| Apply database schema | Agent 1 | 2 |
| Create sync engine base | Agent 1 | 4 |
| Set up APScheduler | Agent 6 | 2 |
| Add Mealie to Docker | Agent 6 | 1 |

### Phase 2: Backend APIs (Days 2-3)
**All agents work in parallel**

| Module | Agent | Files | Hours |
|--------|-------|-------|-------|
| Tasks API + Dispatch | Agent 2 | `core/tasks/` | 8 |
| Projects API | Agent 2 | `core/tasks/` | 4 |
| Habits API | Agent 3 | `core/habits/` | 6 |
| Routines API | Agent 3 | `core/routines/` | 4 |
| Planner API | Agent 4 | `core/planner/` | 6 |
| Calendar API + CalDAV | Agent 4 | `core/calendar/` | 6 |
| Pomodoro API | Agent 5 | `core/pomodoro/` | 4 |
| Journal API | Agent 5 | `core/journal/` | 3 |
| Goals API | Agent 5 | `core/goals/` | 3 |
| Mealie Integration | Agent 6 | `core/meals/` | 4 |
| Agent Tools | Agent 6 | `core/tools/life_os.py` | 4 |
| Scheduling Service | Agent 6 | `core/scheduling/` | 6 |

### Phase 3: Flutter UI (Days 3-5)
**All agents work in parallel**

| Screen | Agent | Files | Hours |
|--------|-------|-------|-------|
| Tasks Screen + Dispatch | Agent 2 | `features/tasks/` | 8 |
| Projects Screen | Agent 2 | `features/tasks/` | 4 |
| Habits Screen + Streaks | Agent 3 | `features/habits/` | 6 |
| Routines Screen | Agent 3 | `features/routines/` | 4 |
| Planner Screen + Blocks | Agent 4 | `features/planner/` | 8 |
| Calendar Screen | Agent 4 | `features/calendar/` | 6 |
| Pomodoro Screen | Agent 5 | `features/pomodoro/` | 4 |
| Journal Screen | Agent 5 | `features/journal/` | 4 |
| Goals Screen | Agent 5 | `features/goals/` | 4 |
| Meals Screen | Agent 6 | `features/meals/` | 4 |

### Phase 4: Integration + Polish (Day 5-6)

| Task | Agent | Hours |
|------|-------|-------|
| Offline sync testing | Agent 1 | 4 |
| Task dispatch integration | Agent 2 | 2 |
| Habit notifications | Agent 3 | 2 |
| Calendar sync testing | Agent 4 | 2 |
| Cross-feature integration | All | 4 |
| Bug fixes + polish | All | 4 |

---

## Linear Issues (Create These)

### Epic: Life OS (OMOS-323 to OMOS-334)

```markdown
## OMOS-323: Database Schema + Sync Foundation
**Assignee:** Agent 1
**Priority:** P0 Critical
**Labels:** backend, database, sync

### Description
Apply Life OS database schema and create sync engine foundation.

### Tasks
- [ ] Add all Life OS tables to init-db.sql
- [ ] Create SQLite schema mirror for edge devices
- [ ] Implement sync_log trigger functions
- [ ] Create sync engine with conflict resolution
- [ ] Add NATS subjects for sync events
- [ ] Test offline → online sync flow

### Files
- `infra/configs/postgres/init-db.sql`
- `core/sync/engine.py`
- `core/sync/conflict.py`
- `core/sync/nats_bridge.py`

### Acceptance Criteria
- All tables created with indexes
- Sync engine handles create/update/delete
- Conflict resolution uses last-write-wins + merge for arrays
- NATS JetStream stores offline events
```

```markdown
## OMOS-324: Edge SQLite + Offline Support
**Assignee:** Agent 1
**Priority:** P0 Critical
**Labels:** backend, sync, edge
**Depends:** OMOS-323

### Description
Implement SQLite storage for edge devices with sync capability.

### Tasks
- [ ] Create SQLite schema matching PostgreSQL
- [ ] Implement local-first write path
- [ ] Queue changes in NATS JetStream
- [ ] Sync on reconnection
- [ ] Handle merge conflicts

### Files
- `core/sync/sqlite.py`
- `core/sync/queue.py`
- `edge/sync_agent.py`
```

```markdown
## OMOS-325: Tasks API + AI Dispatch
**Assignee:** Agent 2
**Priority:** P0 Critical
**Labels:** backend, tasks, dispatch

### Description
Complete task management API with AI dispatch capability.

### Tasks
- [ ] CRUD endpoints for inbox_items enhancements
- [ ] Task dispatch to agent/team/workflow
- [ ] Dependency management
- [ ] Recurrence handling
- [ ] Context and energy filtering

### Endpoints
- GET/POST /api/tasks
- GET/PATCH/DELETE /api/tasks/{id}
- POST /api/tasks/{id}/dispatch
- POST /api/tasks/{id}/complete
- GET /api/tasks/today
- GET /api/tasks/by-context/{context}

### Files
- `core/tasks/models.py`
- `core/tasks/service.py`
- `core/tasks/api.py`
- `core/tasks/dispatch.py`
```

```markdown
## OMOS-326: Projects API
**Assignee:** Agent 2
**Priority:** P1 High
**Labels:** backend, tasks

### Description
Project management for grouping tasks.

### Tasks
- [ ] CRUD endpoints for projects
- [ ] Sequential vs parallel mode
- [ ] Project progress calculation
- [ ] Archive/complete projects

### Endpoints
- GET/POST /api/projects
- GET/PATCH/DELETE /api/projects/{id}
- GET /api/projects/{id}/tasks
- POST /api/projects/{id}/reorder

### Files
- `core/tasks/projects.py`
```

```markdown
## OMOS-327: Habits API
**Assignee:** Agent 3
**Priority:** P0 Critical
**Labels:** backend, habits

### Description
Habit tracking with streaks and A-F scoring.

### Tasks
- [ ] CRUD endpoints for habits
- [ ] Completion logging
- [ ] Streak calculation (use DB trigger)
- [ ] A-F grade calculation (30-day rolling)
- [ ] Frequency handling (daily, weekdays, custom)

### Endpoints
- GET/POST /api/habits
- GET/PATCH/DELETE /api/habits/{id}
- POST /api/habits/{id}/complete
- GET /api/habits/{id}/history
- GET /api/habits/today

### Files
- `core/habits/models.py`
- `core/habits/service.py`
- `core/habits/api.py`
```

```markdown
## OMOS-328: Routines API
**Assignee:** Agent 3
**Priority:** P1 High
**Labels:** backend, routines

### Description
Routine management with step tracking.

### Tasks
- [ ] CRUD for routines and steps
- [ ] Start/complete routine tracking
- [ ] Trigger types (time, location, manual)
- [ ] Step duration tracking

### Endpoints
- GET/POST /api/routines
- GET/PATCH/DELETE /api/routines/{id}
- POST /api/routines/{id}/start
- POST /api/routines/{id}/steps/{step_id}/complete

### Files
- `core/routines/models.py`
- `core/routines/service.py`
- `core/routines/api.py`
```

```markdown
## OMOS-329: Planner API
**Assignee:** Agent 4
**Priority:** P0 Critical
**Labels:** backend, planner

### Description
Time block planner with drag-drop task assignment.

### Tasks
- [ ] CRUD for time blocks
- [ ] Task assignment to blocks
- [ ] Section management (morning/afternoon/evening/night)
- [ ] Recurring blocks (weekly templates)
- [ ] Drag-drop reordering

### Endpoints
- GET/POST /api/planner/blocks
- GET/PATCH/DELETE /api/planner/blocks/{id}
- POST /api/planner/blocks/{id}/assign-task
- DELETE /api/planner/blocks/{id}/tasks/{task_id}
- GET /api/planner/day/{date}
- GET /api/planner/week/{date}

### Files
- `core/planner/models.py`
- `core/planner/service.py`
- `core/planner/api.py`
```

```markdown
## OMOS-330: Calendar API + CalDAV Sync
**Assignee:** Agent 4
**Priority:** P0 Critical
**Labels:** backend, calendar, sync

### Description
Unified calendar with CalDAV sync support.

### Tasks
- [ ] Internal calendar events CRUD
- [ ] CalDAV sync (Google, Apple, Outlook)
- [ ] Unified calendar view combining all sources
- [ ] Conflict with time blocks detection
- [ ] Reminder integration

### Endpoints
- GET/POST /api/calendar/events
- GET/PATCH/DELETE /api/calendar/events/{id}
- GET /api/calendar/day/{date}
- GET /api/calendar/week/{date}
- GET /api/calendar/month/{year}/{month}
- POST /api/calendar/sync
- GET /api/calendar/sources

### Files
- `core/calendar/models.py`
- `core/calendar/service.py`
- `core/calendar/api.py`
- `core/calendar/caldav.py`
```

```markdown
## OMOS-331: Pomodoro + Journal API
**Assignee:** Agent 5
**Priority:** P1 High
**Labels:** backend, pomodoro, journal

### Description
Focus session tracking and daily journaling.

### Tasks
- [ ] Pomodoro session start/end
- [ ] Focus level and work mode tracking
- [ ] Daily journal CRUD
- [ ] Mood/energy tracking
- [ ] Gratitude and wins logging

### Endpoints
- POST /api/pomodoro/start
- POST /api/pomodoro/{id}/end
- GET /api/pomodoro/today
- GET /api/pomodoro/stats
- GET/POST /api/journal
- GET/PATCH /api/journal/{date}

### Files
- `core/pomodoro/`
- `core/journal/`
```

```markdown
## OMOS-332: Goals API
**Assignee:** Agent 5
**Priority:** P1 High
**Labels:** backend, goals

### Description
Goal tracking with progress visualization.

### Tasks
- [ ] CRUD for goals
- [ ] Progress tracking
- [ ] Sub-goals (parent-child)
- [ ] Pillar alignment
- [ ] Target date management

### Endpoints
- GET/POST /api/goals
- GET/PATCH/DELETE /api/goals/{id}
- POST /api/goals/{id}/progress
- GET /api/goals/by-pillar/{pillar}

### Files
- `core/goals/`
```

```markdown
## OMOS-333: Custom Meals Module (AI-Powered)
**Assignee:** Agent 6
**Priority:** P1 High
**Labels:** backend, meals, ai

### Description
Custom recipe and meal planning module. AI handles parsing, suggestions, and scaling.

### Tasks
- [ ] Create `core/meals/` module (models, service, api)
- [ ] Recipe CRUD with JSONB ingredients/instructions
- [ ] AI-powered recipe import from URL
- [ ] Meal plan calendar (week view)
- [ ] Shopping list generation (AI consolidates)
- [ ] Pantry tracking with expiry alerts
- [ ] Agent tools for Legacy

### AI Does the Heavy Lifting
- Parse any recipe URL → structured JSON
- "2 cups flour, sifted" → {amount: 2, unit: "cups", name: "flour", notes: "sifted"}
- Consolidate shopping list (combine duplicate ingredients)
- Scale recipes (2 → 6 servings)
- Suggest meals from pantry items

### Files
- `core/meals/models.py`
- `core/meals/service.py`
- `core/meals/api.py`
- `core/tools/meals.py`

### Acceptance Criteria
- Recipes stored with structured ingredients/instructions
- Import any recipe URL via AI
- Meal plan shows week view
- Shopping list consolidates duplicates
- Pantry tracks items + expiring soon
```

```markdown
## OMOS-334: Agent Tools + Scheduling
**Assignee:** Agent 6
**Priority:** P0 Critical
**Labels:** backend, tools, scheduling

### Description
Life OS agent tools and APScheduler integration.

### Tasks
- [ ] Integrate APScheduler
- [ ] create_task tool
- [ ] complete_habit tool
- [ ] start_routine tool
- [ ] schedule_recurring tool
- [ ] dispatch_task tool
- [ ] get_today_plan tool
- [ ] Scheduled task execution

### Files
- `core/scheduling/scheduler.py`
- `core/scheduling/api.py`
- `core/tools/life_os.py`
- `agents/legacy/agent.py` (add tools)
```

---

## Frontend API Convention

> **IMPORTANT:** All API calls from the Flutter frontend MUST include the `/api` prefix.

### Endpoint Pattern
```
Base URL: http://{host}:8080
API Path: /api/{module}/{endpoint}

Examples:
✅ CORRECT: /api/tasks
❌ WRONG:   /tasks

✅ CORRECT: /api/habits/123/complete
❌ WRONG:   /habits/123/complete
```

### Why This Matters
The backend mounts all Life OS routers under the `/api` prefix:
- Tasks: `/api/tasks`
- Habits: `/api/habits`
- Goals: `/api/goals`
- Calendar: `/api/calendar`
- Planner: `/api/planner`
- etc.

Calling endpoints without `/api` will return 404 Not Found.

### Provider Implementation
All Flutter providers must use the full path:
```dart
// ✅ CORRECT
final response = await _apiClient.get('/api/tasks', queryParams: params);

// ❌ WRONG - will return 404
final response = await _apiClient.get('/tasks', queryParams: params);
```

### Fixed Files (2026-01-28)
- `frontend/lib/lifeos/tasks/providers/tasks_provider.dart` - All `/tasks` → `/api/tasks`
- `frontend/lib/lifeos/scheduling/providers/scheduling_provider.dart` - All `/scheduled` → `/api/scheduled`
- `frontend/lib/lifeos/inbox/providers/inbox_provider.dart` - `/inbox` → `/api/inbox`

---

## API Endpoints Summary

### Tasks Module
```
GET    /api/tasks                    - List tasks with filters
POST   /api/tasks                    - Create task
GET    /api/tasks/{id}               - Get task
PATCH  /api/tasks/{id}               - Update task
DELETE /api/tasks/{id}               - Delete task
POST   /api/tasks/{id}/complete      - Complete task
POST   /api/tasks/{id}/dispatch      - Dispatch to agent/team/workflow
GET    /api/tasks/today              - Today's tasks
GET    /api/tasks/overdue            - Overdue tasks
```

### Projects Module
```
GET    /api/projects                 - List projects
POST   /api/projects                 - Create project
GET    /api/projects/{id}            - Get project with tasks
PATCH  /api/projects/{id}            - Update project
DELETE /api/projects/{id}            - Delete/archive project
POST   /api/projects/{id}/reorder    - Reorder tasks
```

### Habits Module
```
GET    /api/habits                   - List habits
POST   /api/habits                   - Create habit
GET    /api/habits/{id}              - Get habit with stats
PATCH  /api/habits/{id}              - Update habit
DELETE /api/habits/{id}              - Delete habit
POST   /api/habits/{id}/complete     - Log completion
GET    /api/habits/today             - Today's habits
GET    /api/habits/grades            - All habits with A-F grades
```

### Routines Module
```
GET    /api/routines                 - List routines
POST   /api/routines                 - Create routine
GET    /api/routines/{id}            - Get routine with steps
PATCH  /api/routines/{id}            - Update routine
DELETE /api/routines/{id}            - Delete routine
POST   /api/routines/{id}/start      - Start routine
POST   /api/routines/{id}/steps/{step_id}/complete - Complete step
```

### Planner Module
```
GET    /api/planner/blocks           - List time blocks
POST   /api/planner/blocks           - Create block
PATCH  /api/planner/blocks/{id}      - Update block
DELETE /api/planner/blocks/{id}      - Delete block
POST   /api/planner/blocks/{id}/assign - Assign task to block
GET    /api/planner/day/{date}       - Day view
GET    /api/planner/week/{date}      - Week view
```

### Calendar Module
```
GET    /api/calendar/events          - List events
POST   /api/calendar/events          - Create event
GET    /api/calendar/events/{id}     - Get event
PATCH  /api/calendar/events/{id}     - Update event
DELETE /api/calendar/events/{id}     - Delete event
GET    /api/calendar/day/{date}      - Day view (events + blocks + tasks)
GET    /api/calendar/week/{date}     - Week view
GET    /api/calendar/month/{y}/{m}   - Month view
POST   /api/calendar/sync            - Trigger CalDAV sync
```

### Pomodoro Module
```
POST   /api/pomodoro/start           - Start session
POST   /api/pomodoro/{id}/end        - End session
GET    /api/pomodoro/today           - Today's sessions
GET    /api/pomodoro/stats           - Statistics
GET    /api/pomodoro/active          - Current active session
```

### Journal Module
```
GET    /api/journal                  - List entries
POST   /api/journal                  - Create entry
GET    /api/journal/{date}           - Get entry by date
PATCH  /api/journal/{date}           - Update entry
```

### Goals Module
```
GET    /api/goals                    - List goals
POST   /api/goals                    - Create goal
GET    /api/goals/{id}               - Get goal
PATCH  /api/goals/{id}               - Update goal
DELETE /api/goals/{id}               - Delete goal
POST   /api/goals/{id}/progress      - Update progress
```

### Scheduling Module
```
GET    /api/scheduled                - List scheduled tasks
POST   /api/scheduled                - Create scheduled task
PATCH  /api/scheduled/{id}           - Update scheduled task
DELETE /api/scheduled/{id}           - Delete scheduled task
POST   /api/scheduled/{id}/run       - Run now
GET    /api/scheduled/upcoming       - Next 24h scheduled
```

### Meals Module (Custom AI-Powered)
```
GET    /api/meals/recipes            - List recipes
POST   /api/meals/recipes            - Create recipe (manual or AI-parsed)
GET    /api/meals/recipes/{id}       - Get recipe
PATCH  /api/meals/recipes/{id}       - Update recipe
DELETE /api/meals/recipes/{id}       - Delete recipe
POST   /api/meals/recipes/import     - AI imports from URL (scrapes + parses)
POST   /api/meals/recipes/{id}/made  - Log that you made it

GET    /api/meals/plan               - Get meal plan (week view)
POST   /api/meals/plan               - Add to meal plan
DELETE /api/meals/plan/{id}          - Remove from plan

GET    /api/meals/shopping           - Get active shopping list
POST   /api/meals/shopping/generate  - AI generates from meal plan
PATCH  /api/meals/shopping/{id}      - Update (check items)

GET    /api/meals/pantry             - List pantry items
POST   /api/meals/pantry             - Add pantry item
PATCH  /api/meals/pantry/{id}        - Update pantry item
DELETE /api/meals/pantry/{id}        - Remove pantry item
GET    /api/meals/pantry/expiring    - Items expiring soon
```

**AI handles:**
- Recipe import (scrape URL → structured JSON)
- Ingredient parsing ("2 cups flour" → {amount: 2, unit: "cups", name: "flour"})
- Meal suggestions based on pantry + preferences
- Shopping list generation + consolidation
- Recipe scaling (2 servings → 6 servings)
- "What can I make?" from pantry

### Sync Module
```
POST   /api/sync/push                - Push local changes
GET    /api/sync/pull                - Pull remote changes
GET    /api/sync/status              - Sync status
POST   /api/sync/resolve             - Resolve conflict
```

---

## Flutter Screens Structure

```
ui/lib/features/
├── tasks/
│   ├── screens/
│   │   ├── tasks_screen.dart         - Main task list
│   │   ├── task_detail_screen.dart   - Task detail + dispatch
│   │   └── projects_screen.dart      - Project list
│   ├── widgets/
│   │   ├── task_card.dart
│   │   ├── task_dispatch_sheet.dart  - Dispatch to agent/team/workflow
│   │   ├── project_card.dart
│   │   └── task_filters.dart
│   └── providers/
│       ├── tasks_provider.dart
│       └── projects_provider.dart
│
├── habits/
│   ├── screens/
│   │   ├── habits_screen.dart        - Habit list with grades
│   │   └── habit_detail_screen.dart  - Habit history + streaks
│   ├── widgets/
│   │   ├── habit_card.dart           - With streak + grade
│   │   ├── habit_calendar.dart       - Completion calendar
│   │   └── grade_badge.dart          - A-F badge
│   └── providers/
│       └── habits_provider.dart
│
├── routines/
│   ├── screens/
│   │   ├── routines_screen.dart      - Routine list
│   │   └── routine_runner_screen.dart - Step-by-step runner
│   ├── widgets/
│   │   ├── routine_card.dart
│   │   └── routine_step.dart
│   └── providers/
│       └── routines_provider.dart
│
├── planner/
│   ├── screens/
│   │   └── planner_screen.dart       - Day view with blocks
│   ├── widgets/
│   │   ├── time_block.dart
│   │   ├── section_header.dart       - Morning/Afternoon/etc
│   │   └── task_chip.dart            - Draggable task
│   └── providers/
│       └── planner_provider.dart
│
├── calendar/
│   ├── screens/
│   │   ├── calendar_screen.dart      - Month/week view
│   │   └── day_detail_screen.dart    - Full day detail
│   ├── widgets/
│   │   ├── calendar_day.dart
│   │   ├── event_card.dart
│   │   └── unified_day_view.dart     - Events + blocks + tasks
│   └── providers/
│       └── calendar_provider.dart
│
├── pomodoro/
│   ├── screens/
│   │   └── pomodoro_screen.dart      - Timer + stats
│   ├── widgets/
│   │   ├── pomodoro_timer.dart       - Countdown circle
│   │   ├── focus_level_picker.dart
│   │   └── work_mode_picker.dart     - Sit/Stand/Walk
│   └── providers/
│       └── pomodoro_provider.dart
│
├── journal/
│   ├── screens/
│   │   └── journal_screen.dart       - Entry list + editor
│   ├── widgets/
│   │   ├── mood_picker.dart          - 1-9 scale
│   │   ├── energy_picker.dart
│   │   └── gratitude_list.dart
│   └── providers/
│       └── journal_provider.dart
│
├── goals/
│   ├── screens/
│   │   └── goals_screen.dart         - Goal list + progress
│   ├── widgets/
│   │   ├── goal_card.dart
│   │   └── progress_ring.dart
│   └── providers/
│       └── goals_provider.dart
│
└── meals/
    ├── screens/
    │   ├── meals_screen.dart         - Meal plan view
    │   └── recipes_screen.dart       - Recipe browser
    ├── widgets/
    │   ├── meal_card.dart
    │   └── recipe_card.dart
    └── providers/
        └── meals_provider.dart
```

---

## Agent Tools (Legacy)

### core/tools/life_os.py

```python
"""Life OS agent tools for Legacy."""

from agno.tools import tool
from datetime import datetime, date
import httpx

API_BASE = "http://localhost:8081/api"


@tool
async def create_task(
    title: str,
    description: str = None,
    due_date: str = None,
    priority: str = "normal",
    project: str = None,
    pillar: str = None,
    estimated_minutes: int = None,
    context: str = "personal",
    tags: list[str] = None,
) -> str:
    """
    Create a new task in the Life OS.

    Args:
        title: Task title (required)
        description: Task description
        due_date: Due date in YYYY-MM-DD format
        priority: Priority level (low, normal, high, urgent)
        project: Project name to assign to
        pillar: Four pillar category (HP, LE, GE, IT)
        estimated_minutes: Estimated time to complete
        context: Context for filtering (personal, work, errands, calls)
        tags: List of tags

    Returns:
        Confirmation message with task ID
    """
    async with httpx.AsyncClient() as client:
        response = await client.post(
            f"{API_BASE}/tasks",
            json={
                "title": title,
                "description": description,
                "due_date": due_date,
                "priority": priority,
                "project_name": project,
                "pillar": pillar,
                "estimated_minutes": estimated_minutes,
                "context": context,
                "tags": tags or [],
            }
        )
        data = response.json()
        return f"Created task '{title}' with ID {data['id']}"


@tool
async def complete_habit(habit_name: str, notes: str = None) -> str:
    """
    Mark a habit as complete for today.

    Args:
        habit_name: Name of the habit to complete
        notes: Optional notes about the completion

    Returns:
        Confirmation with streak info
    """
    async with httpx.AsyncClient() as client:
        # Find habit by name
        response = await client.get(f"{API_BASE}/habits", params={"search": habit_name})
        habits = response.json()

        if not habits:
            return f"Habit '{habit_name}' not found"

        habit = habits[0]

        # Complete it
        response = await client.post(
            f"{API_BASE}/habits/{habit['id']}/complete",
            json={"notes": notes}
        )
        data = response.json()

        return f"Completed '{habit_name}'! Current streak: {data['current_streak']} days"


@tool
async def start_routine(routine_name: str) -> str:
    """
    Start a routine (morning, evening, workout, etc.).

    Args:
        routine_name: Name of the routine to start

    Returns:
        First step of the routine
    """
    async with httpx.AsyncClient() as client:
        response = await client.get(f"{API_BASE}/routines", params={"search": routine_name})
        routines = response.json()

        if not routines:
            return f"Routine '{routine_name}' not found"

        routine = routines[0]

        response = await client.post(f"{API_BASE}/routines/{routine['id']}/start")
        data = response.json()

        steps = data.get('steps', [])
        if steps:
            first_step = steps[0]
            return f"Started '{routine_name}'!\n\nFirst step: {first_step['name']} ({first_step.get('duration_minutes', '?')} min)"

        return f"Started '{routine_name}' - no steps defined"


@tool
async def schedule_recurring(
    name: str,
    schedule: str,
    execution_type: str,
    execution_target: str,
    description: str = None,
    params: dict = None,
) -> str:
    """
    Schedule a recurring task or workflow.

    Args:
        name: Name for the scheduled task
        schedule: Cron expression (e.g., "0 6 * * *" for 6 AM daily)
        execution_type: Type of execution (agent, team, workflow)
        execution_target: ID or name of target
        description: Description of what this does
        params: Parameters to pass to execution

    Returns:
        Confirmation with next run time

    Examples:
        - "0 6 * * *" = Every day at 6 AM
        - "0 9 * * 1-5" = Weekdays at 9 AM
        - "0 20 * * 0" = Sundays at 8 PM
    """
    async with httpx.AsyncClient() as client:
        response = await client.post(
            f"{API_BASE}/scheduled",
            json={
                "name": name,
                "description": description,
                "cron_expression": schedule,
                "execution_type": execution_type,
                "execution_target": execution_target,
                "execution_params": params or {},
            }
        )
        data = response.json()

        return f"Scheduled '{name}' to run {execution_type}:{execution_target}\nNext run: {data['next_run_at']}"


@tool
async def dispatch_task(
    task_id: str,
    target_type: str,
    target_id: str,
    schedule: str = "now",
) -> str:
    """
    Dispatch a task to an agent, team, or workflow for execution.

    Args:
        task_id: ID of the task to dispatch
        target_type: Type of target (agent, team, workflow)
        target_id: ID or name of the target
        schedule: When to run ("now" or ISO datetime)

    Returns:
        Dispatch confirmation with run ID
    """
    async with httpx.AsyncClient() as client:
        response = await client.post(
            f"{API_BASE}/tasks/{task_id}/dispatch",
            json={
                "target_type": target_type,
                "target_id": target_id,
                "schedule": schedule,
            }
        )
        data = response.json()

        return f"Dispatched task to {target_type}:{target_id}\nRun ID: {data['run_id']}"


@tool
async def get_today_plan() -> str:
    """
    Get today's tasks, habits, time blocks, and scheduled executions.

    Returns:
        Formatted view of today's plan
    """
    async with httpx.AsyncClient() as client:
        today = date.today().isoformat()

        # Fetch all data in parallel
        tasks_resp = await client.get(f"{API_BASE}/tasks/today")
        habits_resp = await client.get(f"{API_BASE}/habits/today")
        planner_resp = await client.get(f"{API_BASE}/planner/day/{today}")
        scheduled_resp = await client.get(f"{API_BASE}/scheduled/upcoming")

        tasks = tasks_resp.json()
        habits = habits_resp.json()
        planner = planner_resp.json()
        scheduled = scheduled_resp.json()

        output = []

        # Tasks
        output.append(f"## Tasks ({len(tasks)})")
        for t in tasks[:10]:
            status = "[ ]" if t['status'] != 'completed' else "[x]"
            output.append(f"  {status} {t['title']}")

        # Habits
        output.append(f"\n## Habits ({len(habits)})")
        for h in habits:
            done = "[x]" if h.get('completed_today') else "[ ]"
            grade = h.get('grade', '?')
            output.append(f"  {done} {h['name']} ({grade})")

        # Time blocks
        output.append(f"\n## Schedule")
        for block in planner.get('blocks', []):
            output.append(f"  {block['start_time']}-{block['end_time']}: {block['title']}")

        # Scheduled executions
        if scheduled:
            output.append(f"\n## Scheduled Runs")
            for s in scheduled[:5]:
                output.append(f"  {s['next_run_at']}: {s['name']}")

        return "\n".join(output)


@tool
async def log_pomodoro(
    duration_minutes: int = 25,
    task: str = None,
    focus_level: int = None,
    work_mode: str = None,
    notes: str = None,
) -> str:
    """
    Log a completed pomodoro focus session.

    Args:
        duration_minutes: Duration of the session (default 25)
        task: Task name or ID worked on
        focus_level: Focus quality 1-5 (1=distracted, 5=flow state)
        work_mode: Work mode used (sit, stand, walk)
        notes: Any notes about the session

    Returns:
        Confirmation with daily stats
    """
    async with httpx.AsyncClient() as client:
        response = await client.post(
            f"{API_BASE}/pomodoro/start",
            json={"planned_duration": duration_minutes, "task_search": task}
        )
        session = response.json()

        # Immediately end it (for logging completed session)
        response = await client.post(
            f"{API_BASE}/pomodoro/{session['id']}/end",
            json={
                "focus_level": focus_level,
                "work_mode": work_mode,
                "notes": notes,
            }
        )
        data = response.json()

        return f"Logged {duration_minutes}min pomodoro!\nToday: {data['today_count']} sessions, {data['today_minutes']} minutes"
```

---

## Coordination Protocol

### Communication
1. **Linear is source of truth** - All work tracked in Linear issues
2. **CLAUDE.md updates** - Agents update CLAUDE.md when completing major features
3. **This document** - Reference for API contracts, schema, structure

### Sync Points
- **End of Day 1:** Schema applied, sync engine base ready
- **End of Day 2:** All backend APIs functional
- **End of Day 3:** Backend complete, Flutter scaffolds ready
- **End of Day 4:** All Flutter screens functional
- **End of Day 5:** Integration complete, testing done
- **End of Day 6:** Bug fixes, polish, deploy

### Conflict Resolution
- **API contracts** - This document is authoritative
- **Code conflicts** - Last agent to push resolves
- **Design decisions** - Escalate to Lead (Mother) agent

### File Ownership

| Files | Owner Agent |
|-------|-------------|
| `core/sync/` | Agent 1 |
| `core/tasks/`, `core/projects/` | Agent 2 |
| `core/habits/`, `core/routines/` | Agent 3 |
| `core/planner/`, `core/calendar/` | Agent 4 |
| `core/pomodoro/`, `core/journal/`, `core/goals/` | Agent 5 |
| `core/meals/`, `core/scheduling/`, `core/tools/life_os.py` | Agent 6 |
| `ui/lib/features/tasks/` | Agent 2 |
| `ui/lib/features/habits/`, `ui/lib/features/routines/` | Agent 3 |
| `ui/lib/features/planner/`, `ui/lib/features/calendar/` | Agent 4 |
| `ui/lib/features/pomodoro/`, `ui/lib/features/journal/`, `ui/lib/features/goals/` | Agent 5 |
| `ui/lib/features/meals/` | Agent 6 |

---

## Success Criteria

> **Updated: 2026-01-27** - Progress tracked below

- [x] All 10 database tables created with indexes ✅
- [x] All 76+ API endpoints functional ✅
- [ ] All 10 Flutter screens complete (4/10 done - 40%)
- [x] Offline sync working (SQLite → NATS → PostgreSQL) ✅
- [x] AI task dispatch functional (tap task → agent/team/workflow) ✅
- [x] Agent tools working (Legacy can create tasks, complete habits, etc.) ✅
- [x] APScheduler integration (backend ready, UI pending)
- [x] Custom Meals module with AI parsing (replaced Mealie) ✅
- [x] Calendar showing unified view ✅
- [x] Habits showing A-F grades and streaks ✅
- [ ] No Todoist dependency (migration pending)

**Current Progress: 70% Complete**

---

## Implementation Status (2026-01-27)

### Backend Modules (~18,853 lines)

| Module | Lines | Status |
|--------|-------|--------|
| Tasks | 2,840 | ✅ Complete |
| Habits | 1,223 | ✅ Complete |
| Goals | 768 | ✅ Complete |
| Journal | 668 | ✅ Complete |
| Pomodoro | 573 | ✅ Complete |
| Calendar | 1,727 | ✅ Complete |
| Routines | 1,352 | ✅ Complete |
| Meals | 2,051 | ✅ Complete |
| Sync Engine | 2,992 | ✅ Complete |
| Agent Tools | 1,125 | ✅ Complete |
| Planner | 1,269 | 🟡 Partial |
| Scheduling | 1,265 | 🟡 Partial |

### Flutter Screens (~4,916 lines)

| Screen | Status |
|--------|--------|
| Tasks | ✅ Complete |
| Goals | ✅ Complete |
| Journal | ✅ Complete |
| Pomodoro | ✅ Complete |
| Planner | ❌ Not started |
| Calendar | ❌ Not started |
| Habits | ❌ Not started |
| Routines | ❌ Not started |
| Meals | ❌ Not started |
| Scheduling | ❌ Not started |

### Remaining Work

| Phase | Task | Effort |
|-------|------|--------|
| 1 | Complete Planner backend | 1 day |
| 2 | Complete Scheduling backend | 1 day |
| 3 | Build 6 Flutter screens | 3-4 days |
| 4 | Integration testing | 1-2 days |
| 5 | Deployment | 1 day |

**ETA to 100%:** 7-9 days

---

## Timeline Summary

| Day | Focus | Agents Working | Status |
|-----|-------|----------------|--------|
| 1 | Database + Foundation | 2 (Agent 1, 6) | ✅ Done |
| 2 | Backend APIs | 6 (All) | ✅ Done |
| 3 | Backend + Flutter Start | 6 (All) | ✅ Done |
| 4 | Flutter UI | 6 (All) | 🟡 40% |
| 5 | Flutter + Integration | 6 (All) | ⏳ Pending |
| 6 | Polish + Deploy | 6 (All) | ⏳ Pending |

**Total: ~120 agent-hours across 6 days**
