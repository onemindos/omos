# OneIntel - Implementation Plan

> Unified item system for OneMind OS - tasks, projects, docs, habits in one place.

## Executive Summary

Build a **simple, unified item system** that gives Agno agents full context over tasks, projects, docs, and habits. Learn from Life OS v1's over-engineering failure.

---

## Part 1: Why Life OS v1 Failed

### What We Built (Archived at `.archive/features/lifeos/`)

```
20 MODULES BUILT SIMULTANEOUSLY:
├── activity/      ├── goals/        ├── pages/
├── assets/        ├── habits/       ├── planner/
├── bridge/        ├── inbox/        ├── pomodoro/
├── calendar/      ├── insights/     ├── routines/
├── consciousness/ ├── journal/      ├── scheduling/
├── meals/         ├── notifications/├── scores/
└── tasks/         └── platform/
```

### Critical Mistakes

| Mistake | Evidence | Impact |
|---------|----------|--------|
| **40+ field Task model** | `tasks/models.py` lines 93-161 | Complexity explosion |
| **Inbox = Everything** | 8 types in one table | Confused data model |
| **Traditional REST services** | `TaskService` at 933 lines | Agents can't call services |
| **No Agno integration** | Zero `@tool` decorators | Agents had no access |
| **Built services, not tools** | `service.py` pattern everywhere | Wrong architecture |

### Architecture Comparison

```
LIFE OS v1 (FAILED):
┌─────────────┐    ┌─────────────────┐    ┌─────────────┐
│   Flutter   │───>│  REST Services  │───>│  PostgreSQL │
│   Frontend  │    │  (933+ lines)   │    │  (20 tables)│
└─────────────┘    └─────────────────┘    └─────────────┘
                          ↑
                    Agents CAN'T
                    call services!

ONEINTEL (NEW):
┌─────────────┐    ┌─────────────────┐    ┌─────────────┐
│   Flutter   │───>│   Agno Agents   │───>│  oneintel() │───>│  SQLite/PG  │
│   Frontend  │    │  (@tool access) │    │   Tool      │    │  (1 table)  │
└─────────────┘    └─────────────────┘    └─────────────┘    └─────────────┘
```

---

## Part 2: OneIntel Architecture

### Design Principles

1. **One table** - `items` stores everything
2. **One tool** - `oneintel()` for all agents
3. **12 fields max** - Start minimal, add when needed
4. **Agent-first** - Built for Agno, not REST
5. **Hierarchy via parent_id** - Projects contain tasks contain subtasks

### System Architecture

```
┌────────────────────────────────────────────────────────────────────────────┐
│                              ONEMIND OS                                     │
├────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │                         AGNO AGENT LAYER                            │   │
│  │                                                                     │   │
│  │   @legacy    @coder    @researcher    @planner    @reviewer        │   │
│  │      │          │           │             │            │           │   │
│  │      └──────────┴───────────┴─────────────┴────────────┘           │   │
│  │                              │                                      │   │
│  │                        ALL AGENTS USE                               │   │
│  │                              ▼                                      │   │
│  │  ┌────────────────────────────────────────────────────────────┐    │   │
│  │  │                     oneintel() TOOL                         │    │   │
│  │  │                                                            │    │   │
│  │  │  Actions:                                                  │    │   │
│  │  │  • list    - Filter by type, status, parent, tags, due    │    │   │
│  │  │  • get     - Single item by ID                            │    │   │
│  │  │  • create  - New project/task/doc/habit                   │    │   │
│  │  │  • update  - Modify any field                             │    │   │
│  │  │  • complete- Mark done (with optional notes)              │    │   │
│  │  │  • delete  - Remove item                                  │    │   │
│  │  │  • search  - Full-text search across all items            │    │   │
│  │  │  • tree    - Get hierarchy (project → tasks → subtasks)   │    │   │
│  │  │  • today   - Due today + overdue + habits                 │    │   │
│  │  │  • stats   - Counts by type, status, completion rate      │    │   │
│  │  └────────────────────────────────────────────────────────────┘    │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                                    │                                        │
│                                    ▼                                        │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │                         STORAGE LAYER                               │   │
│  │                                                                     │   │
│  │   ┌─────────────────────────────────────────────────────────────┐  │   │
│  │   │                      items TABLE                            │  │   │
│  │   │                                                             │  │   │
│  │   │  id         │ type    │ title        │ content             │  │   │
│  │   │  parent_id  │ status  │ due_date     │ priority            │  │   │
│  │   │  tags       │ streak  │ last_done    │ created_at          │  │   │
│  │   │  updated_at │ user_id │              │                     │  │   │
│  │   └─────────────────────────────────────────────────────────────┘  │   │
│  │                                                                     │   │
│  │   Storage Options:                                                  │   │
│  │   • SQLite (local, fast, simple) ← Current                         │   │
│  │   • PostgreSQL (production, RAG via pgvector)                      │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
└────────────────────────────────────────────────────────────────────────────┘
```

### Data Model

```
┌──────────────────────────────────────────────────────────────────────────┐
│                              items TABLE                                  │
├──────────────────────────────────────────────────────────────────────────┤
│                                                                          │
│  CORE FIELDS (12 total):                                                 │
│  ─────────────────────────────────────────────────────────────────────   │
│                                                                          │
│  id          TEXT           Primary key (UUID)                           │
│  type        TEXT           project | task | doc | habit                 │
│  title       TEXT           Required, max 500 chars                      │
│  content     TEXT           Markdown body (nullable)                     │
│  parent_id   TEXT           Foreign key to items.id (nullable)           │
│  status      TEXT           todo | in_progress | done | archived         │
│  due_date    DATE           When item is due (nullable)                  │
│  priority    INT            1=urgent, 2=high, 3=normal, 4=low            │
│  tags        TEXT           JSON array of strings                        │
│  streak      INT            Habit completion streak (default 0)          │
│  last_done   TIMESTAMP      Last completion (for habits)                 │
│  created_at  TIMESTAMP      Auto-set on create                           │
│  updated_at  TIMESTAMP      Auto-set on update                           │
│  user_id     TEXT           Default "zeus"                               │
│                                                                          │
├──────────────────────────────────────────────────────────────────────────┤
│                                                                          │
│  HIERARCHY EXAMPLES:                                                     │
│  ─────────────────────────────────────────────────────────────────────   │
│                                                                          │
│  Project "Website Redesign"     (type=project, parent_id=NULL)           │
│    └── Task "Design homepage"   (type=task, parent_id=project.id)        │
│         └── Task "Hero section" (type=task, parent_id=task.id)           │
│         └── Task "Nav menu"     (type=task, parent_id=task.id)           │
│    └── Doc "Design specs"       (type=doc, parent_id=project.id)         │
│                                                                          │
│  Habit "Morning meditation"     (type=habit, parent_id=NULL)             │
│    streak=15, last_done=2025-01-31                                       │
│                                                                          │
└──────────────────────────────────────────────────────────────────────────┘
```

---

## Part 3: Item Types

### PROJECT
- Container for tasks and docs
- Has progress (computed from child task completion)
- Can be nested (sub-projects)
- Status: todo | in_progress | done | archived

### TASK
- Actionable work item
- Can have subtasks (parent_id = another task)
- Can belong to project (parent_id = project)
- Has due_date and priority
- Status: todo | in_progress | done

### DOC
- Markdown document
- Can belong to project
- Content is the full document body
- Searchable via full-text
- Status: draft | published | archived

### HABIT
- Recurring item tracked by streak
- No parent_id (top-level only)
- streak increments on daily completion
- last_done tracks last completion
- Resets if missed (streak = 0)

---

## Part 4: Implementation Phases

### Phase 1: Backend Core (COMPLETE)

```
✅ oneintel/__init__.py     - Package init
✅ oneintel/models.py       - Pydantic models
✅ oneintel/schema.sql      - SQLite schema with FTS
✅ oneintel/storage.py      - Storage layer
✅ oneintel/tool.py         - Standalone tool
✅ oneintel/api.py          - FastAPI for testing
✅ oneintel/README.md       - Documentation
✅ oneintel/PLAN.md         - This document
```

### Phase 2: Testing & Validation

```
[ ] Test all 10 tool actions
[ ] Test API endpoints via curl
[ ] Create sample data (projects, tasks, habits)
[ ] Verify hierarchy works
[ ] Verify habit streaks work
[ ] Test full-text search
```

### Phase 3: OneMind Integration (Later)

```
[ ] Move tool.py to backend/agno/tools/oneintel.py
[ ] Register in backend/agno/tools/__init__.py
[ ] Add to agent presets
[ ] Create Flutter frontend screens
[ ] Add tool cards for chat
```

### Phase 4: Advanced Features (Only When Needed)

```
[ ] Recurrence (RRULE) for habits and tasks
[ ] Time estimates (estimated_minutes)
[ ] Contexts (GTD @home, @work, @phone)
[ ] Dependencies (blocked_by)
[ ] Calendar integration (CalDAV)
[ ] External sync (Todoist bidirectional)
[ ] Offline-first (Isar/Drift)

RULE: Only add a field when 3+ use cases require it
```

---

## Part 5: Quick Reference

### Tool Actions

| Action | Required Params | Returns |
|--------|-----------------|---------|
| list | (optional filters) | items[] |
| get | item_id | item |
| create | type, title | item |
| update | item_id, (fields) | item |
| complete | item_id | item |
| delete | item_id | success |
| search | query | items[] |
| tree | item_id | nested tree |
| today | - | today view |
| stats | - | statistics |

### Item Types

| Type | Purpose | Has Children | Has Streak |
|------|---------|--------------|------------|
| project | Container | Yes (tasks, docs) | No |
| task | Work item | Yes (subtasks) | No |
| doc | Document | No | No |
| habit | Recurring | No | Yes |

### Status Values

| Status | Used By | Meaning |
|--------|---------|---------|
| todo | task, habit | Not started |
| in_progress | task, project | Being worked on |
| done | task, habit | Completed |
| archived | all | Hidden from lists |
| draft | doc | Not published |
| published | doc | Finalized |

### Priority Levels

| Value | Name | Color |
|-------|------|-------|
| 1 | Urgent | Red |
| 2 | High | Orange |
| 3 | Normal | Blue |
| 4 | Low | Gray |

---

## Part 6: Usage Examples

### Tool Examples

```python
from oneintel import oneintel

# List all projects
oneintel(action="list", type="project")

# Create a task in a project
oneintel(action="create", type="task", title="Review PR",
         parent_id="proj_123", due_date="2025-02-05", priority=2)

# Complete a task
oneintel(action="complete", item_id="task_456")

# Track habit completion
oneintel(action="complete", item_id="habit_789")

# Get project tree with all tasks
oneintel(action="tree", item_id="proj_123")

# Today's focus
oneintel(action="today")

# Search across everything
oneintel(action="search", query="authentication")
```

### API Examples

```bash
# Create a project
curl -X POST http://localhost:8081/items \
  -H "Content-Type: application/json" \
  -d '{"type": "project", "title": "My First Project"}'

# List all items
curl http://localhost:8081/items

# Get today's focus
curl http://localhost:8081/items/today

# Create a habit
curl -X POST http://localhost:8081/items \
  -H "Content-Type: application/json" \
  -d '{"type": "habit", "title": "Morning meditation"}'

# Complete a habit
curl -X POST http://localhost:8081/items/{habit_id}/complete
```

### Chat Examples (Future)

```
USER: @legacy create a project called "Mobile App v2"

LEGACY: I'll create that project for you.
[Tool: oneintel(action="create", type="project", title="Mobile App v2")]

Created project "Mobile App v2" (id: proj_abc123)

---

USER: What's on my plate today?

LEGACY: Let me check your tasks and habits.
[Tool: oneintel(action="today")]

Today's Focus:

OVERDUE (2):
• Review PR #45 (due Jan 30)
• Update docs (due Jan 29)

DUE TODAY (1):
• Deploy staging

HABITS (3/5 done):
✅ Morning meditation (15 day streak!)
✅ Exercise
⬜ Journal
```

---

## Part 7: What NOT to Build

Learn from Life OS v1:

| DON'T Build | Why |
|-------------|-----|
| 20 modules at once | Scope creep killed v1 |
| 40+ field models | Over-engineering |
| Inbox = Everything | Confused data model |
| Traditional services | Agents need tools, not services |
| Multiple tables | One table is enough |
| Complex enums | Keep types simple |
| Scheduling engine | Add later if needed |
| Notification system | Use existing |
| Meal tracking | Separate concern |
| Score calculations | Add later |

---

*Document Version: 1.0*
*Created: 2025-02-01*
*Last Updated: 2025-02-01*
