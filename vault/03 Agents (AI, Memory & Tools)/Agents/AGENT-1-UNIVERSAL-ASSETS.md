# Agent 1: Universal Assets (Opus 4.5)

> **Model:** Opus 4.5
> **Focus:** Backend Architecture - Universal Asset System
> **Issues:** 3 | **Points:** 13
> **Dependencies:** None - Start immediately

---

## Your Mission

You are Agent 1 of the 8-agent swarm building OneMind's mobile-first evolution. Your critical job is creating the **Universal Asset** system that unifies all trackable items (tasks, habits, goals, events, meals, journal entries) into a single schema with one API.

This is FOUNDATIONAL work. Agents 2-5 depend on you.

---

## Your Assigned Issues

| ID | Title | Points | Priority |
|----|-------|--------|----------|
| **OMOS-467** | Universal Asset Schema | 5 | P0-Critical |
| **OMOS-468** | Asset Unified API | 5 | P0-Critical |
| **OMOS-460** | Planner Drag-Drop Backend | 3 | P0-Critical |

---

## Existing Codebase Context

### What Already Exists (DO NOT REWRITE)

The backend has mature modules you should build ON TOP OF:

```
backend/lifeos/
├── tasks/          # ✅ Full CRUD, projects, subtasks, recurrence
│   ├── models.py   # Task, TaskCreate, TaskUpdate
│   ├── service.py  # TaskService
│   └── api.py      # /api/tasks endpoints
├── habits/         # ✅ Streak tracking, A-F grading
│   ├── models.py   # Habit, HabitCompletion
│   ├── service.py  # HabitService
│   └── api.py      # /api/habits endpoints
├── goals/          # ✅ Progress tracking, sub-goals
│   ├── models.py   # Goal, GoalProgress
│   ├── service.py  # GoalService
│   └── api.py      # /api/goals endpoints
├── calendar/       # ✅ CalDAV integration
├── inbox/          # ✅ GTD-style processing
├── journal/        # ✅ Daily reflection
├── meals/          # ✅ Nutritional tracking
├── planner/        # ✅ Time block management
└── scheduling/     # ✅ Event scheduling
```

### Database
- **PostgreSQL** via asyncpg
- Existing tables: tasks, habits, goals, calendar_events, inbox_items, journal_entries, meals
- Use Alembic for migrations

### Patterns to Follow

```python
# Service pattern from existing code
class TaskService:
    def __init__(self, db: Database):
        self.db = db

    async def create(self, task: TaskCreate) -> Task:
        ...

    async def get(self, task_id: str) -> Task:
        ...

    async def list(self, filters: TaskFilters) -> List[Task]:
        ...

# API pattern
from fastapi import APIRouter, Depends
router = APIRouter(prefix="/api/tasks", tags=["tasks"])

@router.post("/", response_model=Task)
async def create_task(task: TaskCreate, db: Database = Depends(get_db)):
    return await TaskService(db).create(task)
```

---

## Implementation Requirements

### OMOS-467: Universal Asset Schema

**Create:** `backend/lifeos/assets/models.py`

```python
from enum import Enum
from typing import Optional, List, Dict, Any
from uuid import UUID
from datetime import datetime
from pydantic import BaseModel, Field

class AssetType(str, Enum):
    TASK = "task"
    HABIT = "habit"
    GOAL = "goal"
    EVENT = "event"
    MEAL = "meal"
    JOURNAL = "journal"
    PROJECT = "project"
    ROUTINE = "routine"
    NOTE = "note"
    REMINDER = "reminder"

class Pillar(str, Enum):
    HP = "hp"  # Holistic Performance
    LE = "le"  # Legacy Evolution
    GE = "ge"  # Generational Endeavor

class Priority(str, Enum):
    P0 = "p0"  # Critical
    P1 = "p1"  # High
    P2 = "p2"  # Medium
    P3 = "p3"  # Low

class EnergyLevel(str, Enum):
    HIGH = "high"
    MEDIUM = "medium"
    LOW = "low"

class AssetStatus(str, Enum):
    PENDING = "pending"
    IN_PROGRESS = "in_progress"
    COMPLETED = "completed"
    DEFERRED = "deferred"
    ARCHIVED = "archived"

class RecurrenceRule(BaseModel):
    frequency: str  # daily, weekly, monthly, yearly
    interval: int = 1
    days_of_week: Optional[List[int]] = None
    day_of_month: Optional[int] = None
    end_date: Optional[datetime] = None

class Asset(BaseModel):
    """Universal asset that can represent any trackable item"""

    # Identity
    id: UUID
    asset_type: AssetType
    title: str
    description: Optional[str] = None

    # Classification
    pillar: Optional[Pillar] = None
    priority: Priority = Priority.P2
    energy_level: Optional[EnergyLevel] = None

    # Timing
    do_date: Optional[datetime] = None  # When to DO it (Motion-style)
    due_date: Optional[datetime] = None  # Hard deadline
    duration_minutes: Optional[int] = None
    recurring: Optional[RecurrenceRule] = None

    # Status
    status: AssetStatus = AssetStatus.PENDING
    progress: float = 0.0  # 0.0 to 1.0
    streak: Optional[int] = None  # For habits

    # Relations
    parent_id: Optional[UUID] = None  # For subtasks/sub-goals
    project_id: Optional[UUID] = None
    linked_assets: List[UUID] = []

    # Metadata
    tags: List[str] = []
    context: Dict[str, Any] = {}
    source: str = "manual"  # manual, capture, agent, sync

    # Timestamps
    created_at: datetime
    updated_at: datetime
    completed_at: Optional[datetime] = None

    class Config:
        from_attributes = True

class AssetCreate(BaseModel):
    """Schema for creating new assets"""
    asset_type: AssetType
    title: str
    description: Optional[str] = None
    pillar: Optional[Pillar] = None
    priority: Priority = Priority.P2
    energy_level: Optional[EnergyLevel] = None
    do_date: Optional[datetime] = None
    due_date: Optional[datetime] = None
    duration_minutes: Optional[int] = None
    recurring: Optional[RecurrenceRule] = None
    parent_id: Optional[UUID] = None
    project_id: Optional[UUID] = None
    tags: List[str] = []
    context: Dict[str, Any] = {}
    source: str = "manual"

class AssetUpdate(BaseModel):
    """Schema for updating assets"""
    title: Optional[str] = None
    description: Optional[str] = None
    pillar: Optional[Pillar] = None
    priority: Optional[Priority] = None
    energy_level: Optional[EnergyLevel] = None
    do_date: Optional[datetime] = None
    due_date: Optional[datetime] = None
    duration_minutes: Optional[int] = None
    recurring: Optional[RecurrenceRule] = None
    status: Optional[AssetStatus] = None
    progress: Optional[float] = None
    parent_id: Optional[UUID] = None
    project_id: Optional[UUID] = None
    tags: Optional[List[str]] = None
    context: Optional[Dict[str, Any]] = None

class AssetFilters(BaseModel):
    """Filters for listing assets"""
    asset_type: Optional[AssetType] = None
    pillar: Optional[Pillar] = None
    status: Optional[AssetStatus] = None
    priority: Optional[Priority] = None
    do_date_from: Optional[datetime] = None
    do_date_to: Optional[datetime] = None
    due_date_from: Optional[datetime] = None
    due_date_to: Optional[datetime] = None
    project_id: Optional[UUID] = None
    parent_id: Optional[UUID] = None
    tags: Optional[List[str]] = None
    search: Optional[str] = None
    limit: int = 100
    offset: int = 0
```

### OMOS-468: Asset Unified API

**Create:** `backend/lifeos/assets/service.py`

```python
from typing import List, Optional
from uuid import UUID, uuid4
from datetime import datetime
from database import Database
from .models import Asset, AssetCreate, AssetUpdate, AssetFilters, AssetStatus

class AssetService:
    def __init__(self, db: Database):
        self.db = db

    async def create(self, asset: AssetCreate) -> Asset:
        """Create a new asset of any type."""
        now = datetime.utcnow()
        asset_dict = asset.model_dump()
        asset_dict.update({
            "id": uuid4(),
            "status": AssetStatus.PENDING,
            "progress": 0.0,
            "created_at": now,
            "updated_at": now,
        })

        # Insert into database
        query = """
            INSERT INTO assets (id, asset_type, title, description, pillar, priority,
                energy_level, do_date, due_date, duration_minutes, recurring, status,
                progress, streak, parent_id, project_id, linked_assets, tags, context,
                source, created_at, updated_at, completed_at)
            VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15,
                $16, $17, $18, $19, $20, $21, $22, $23)
            RETURNING *
        """
        row = await self.db.fetchrow(query, *asset_dict.values())
        return Asset(**row)

    async def get(self, asset_id: UUID) -> Optional[Asset]:
        """Get a single asset by ID."""
        query = "SELECT * FROM assets WHERE id = $1"
        row = await self.db.fetchrow(query, asset_id)
        return Asset(**row) if row else None

    async def list(self, filters: AssetFilters) -> List[Asset]:
        """List assets with filtering."""
        conditions = []
        params = []
        param_idx = 1

        if filters.asset_type:
            conditions.append(f"asset_type = ${param_idx}")
            params.append(filters.asset_type)
            param_idx += 1

        if filters.pillar:
            conditions.append(f"pillar = ${param_idx}")
            params.append(filters.pillar)
            param_idx += 1

        if filters.status:
            conditions.append(f"status = ${param_idx}")
            params.append(filters.status)
            param_idx += 1

        # Add more filter conditions...

        where_clause = " AND ".join(conditions) if conditions else "1=1"
        query = f"""
            SELECT * FROM assets
            WHERE {where_clause}
            ORDER BY COALESCE(do_date, due_date, created_at) ASC
            LIMIT ${param_idx} OFFSET ${param_idx + 1}
        """
        params.extend([filters.limit, filters.offset])

        rows = await self.db.fetch(query, *params)
        return [Asset(**row) for row in rows]

    async def update(self, asset_id: UUID, update: AssetUpdate) -> Optional[Asset]:
        """Update an asset."""
        update_dict = update.model_dump(exclude_unset=True)
        update_dict["updated_at"] = datetime.utcnow()

        # Build SET clause
        set_parts = [f"{k} = ${i+1}" for i, k in enumerate(update_dict.keys())]
        set_clause = ", ".join(set_parts)

        query = f"""
            UPDATE assets SET {set_clause}
            WHERE id = ${len(update_dict) + 1}
            RETURNING *
        """
        row = await self.db.fetchrow(query, *update_dict.values(), asset_id)
        return Asset(**row) if row else None

    async def delete(self, asset_id: UUID) -> bool:
        """Delete an asset."""
        query = "DELETE FROM assets WHERE id = $1 RETURNING id"
        row = await self.db.fetchrow(query, asset_id)
        return row is not None

    async def complete(self, asset_id: UUID) -> Optional[Asset]:
        """Mark an asset as completed."""
        now = datetime.utcnow()
        query = """
            UPDATE assets
            SET status = $1, completed_at = $2, updated_at = $2, progress = 1.0
            WHERE id = $3
            RETURNING *
        """
        row = await self.db.fetchrow(query, AssetStatus.COMPLETED, now, asset_id)
        return Asset(**row) if row else None

    async def get_today(self) -> List[Asset]:
        """Get all assets for today (by do_date or due_date)."""
        today_start = datetime.utcnow().replace(hour=0, minute=0, second=0, microsecond=0)
        today_end = today_start.replace(hour=23, minute=59, second=59)

        query = """
            SELECT * FROM assets
            WHERE status NOT IN ('completed', 'archived')
            AND (
                (do_date >= $1 AND do_date <= $2) OR
                (due_date >= $1 AND due_date <= $2) OR
                (asset_type = 'habit' AND recurring IS NOT NULL)
            )
            ORDER BY COALESCE(do_date, due_date) ASC
        """
        rows = await self.db.fetch(query, today_start, today_end)
        return [Asset(**row) for row in rows]

    async def search(self, query_text: str, limit: int = 20) -> List[Asset]:
        """Full-text search across assets."""
        query = """
            SELECT * FROM assets
            WHERE to_tsvector('english', title || ' ' || COALESCE(description, ''))
                @@ plainto_tsquery('english', $1)
            ORDER BY ts_rank(to_tsvector('english', title || ' ' || COALESCE(description, '')),
                plainto_tsquery('english', $1)) DESC
            LIMIT $2
        """
        rows = await self.db.fetch(query, query_text, limit)
        return [Asset(**row) for row in rows]
```

**Create:** `backend/lifeos/assets/api.py`

```python
from fastapi import APIRouter, Depends, HTTPException, Query
from typing import List, Optional
from uuid import UUID
from database import get_db, Database
from .models import Asset, AssetCreate, AssetUpdate, AssetFilters, AssetType, Pillar, AssetStatus
from .service import AssetService

router = APIRouter(prefix="/api/v1/assets", tags=["assets"])

@router.get("/", response_model=List[Asset])
async def list_assets(
    asset_type: Optional[AssetType] = None,
    pillar: Optional[Pillar] = None,
    status: Optional[AssetStatus] = None,
    limit: int = Query(default=100, le=500),
    offset: int = 0,
    db: Database = Depends(get_db)
):
    """List assets with optional filters."""
    filters = AssetFilters(
        asset_type=asset_type,
        pillar=pillar,
        status=status,
        limit=limit,
        offset=offset
    )
    return await AssetService(db).list(filters)

@router.post("/", response_model=Asset, status_code=201)
async def create_asset(asset: AssetCreate, db: Database = Depends(get_db)):
    """Create a new asset of any type."""
    return await AssetService(db).create(asset)

@router.get("/today", response_model=List[Asset])
async def get_today_assets(db: Database = Depends(get_db)):
    """Get all assets for today."""
    return await AssetService(db).get_today()

@router.get("/search", response_model=List[Asset])
async def search_assets(
    q: str = Query(..., min_length=1),
    limit: int = Query(default=20, le=100),
    db: Database = Depends(get_db)
):
    """Full-text search across assets."""
    return await AssetService(db).search(q, limit)

@router.get("/{asset_id}", response_model=Asset)
async def get_asset(asset_id: UUID, db: Database = Depends(get_db)):
    """Get a single asset by ID."""
    asset = await AssetService(db).get(asset_id)
    if not asset:
        raise HTTPException(status_code=404, detail="Asset not found")
    return asset

@router.patch("/{asset_id}", response_model=Asset)
async def update_asset(asset_id: UUID, update: AssetUpdate, db: Database = Depends(get_db)):
    """Update an asset."""
    asset = await AssetService(db).update(asset_id, update)
    if not asset:
        raise HTTPException(status_code=404, detail="Asset not found")
    return asset

@router.delete("/{asset_id}", status_code=204)
async def delete_asset(asset_id: UUID, db: Database = Depends(get_db)):
    """Delete an asset."""
    deleted = await AssetService(db).delete(asset_id)
    if not deleted:
        raise HTTPException(status_code=404, detail="Asset not found")

@router.post("/{asset_id}/complete", response_model=Asset)
async def complete_asset(asset_id: UUID, db: Database = Depends(get_db)):
    """Mark an asset as completed."""
    asset = await AssetService(db).complete(asset_id)
    if not asset:
        raise HTTPException(status_code=404, detail="Asset not found")
    return asset
```

### OMOS-460: Planner Drag-Drop Backend

**Update:** `backend/lifeos/planner/service.py`

```python
# Add these methods to existing PlannerService

async def move_block(self, block_id: str, new_start: datetime, new_end: datetime) -> TimeBlock:
    """Move a time block to a new time slot."""
    # Check for conflicts
    conflicts = await self._check_conflicts(new_start, new_end, exclude_id=block_id)
    if conflicts:
        raise ConflictError(f"Conflicts with: {[c.title for c in conflicts]}")

    query = """
        UPDATE time_blocks
        SET start_time = $1, end_time = $2, updated_at = NOW()
        WHERE id = $3
        RETURNING *
    """
    row = await self.db.fetchrow(query, new_start, new_end, block_id)
    return TimeBlock(**row)

async def resize_block(self, block_id: str, new_end: datetime) -> TimeBlock:
    """Resize a time block by changing end time."""
    block = await self.get_block(block_id)
    if not block:
        raise NotFoundError("Block not found")

    # Check for conflicts with new duration
    conflicts = await self._check_conflicts(block.start_time, new_end, exclude_id=block_id)
    if conflicts:
        raise ConflictError(f"Conflicts with: {[c.title for c in conflicts]}")

    query = """
        UPDATE time_blocks SET end_time = $1, updated_at = NOW()
        WHERE id = $2
        RETURNING *
    """
    row = await self.db.fetchrow(query, new_end, block_id)
    return TimeBlock(**row)

async def swap_blocks(self, block_id_1: str, block_id_2: str) -> tuple[TimeBlock, TimeBlock]:
    """Swap two time blocks."""
    block1 = await self.get_block(block_id_1)
    block2 = await self.get_block(block_id_2)

    if not block1 or not block2:
        raise NotFoundError("One or both blocks not found")

    # Swap times
    async with self.db.transaction():
        query = """
            UPDATE time_blocks SET start_time = $1, end_time = $2, updated_at = NOW()
            WHERE id = $3
            RETURNING *
        """
        new_block1 = await self.db.fetchrow(query, block2.start_time, block2.end_time, block_id_1)
        new_block2 = await self.db.fetchrow(query, block1.start_time, block1.end_time, block_id_2)

    return TimeBlock(**new_block1), TimeBlock(**new_block2)

async def assign_task_to_block(self, block_id: str, task_id: str) -> TimeBlock:
    """Assign an unscheduled task to a time block."""
    query = """
        UPDATE time_blocks
        SET linked_asset_id = $1, updated_at = NOW()
        WHERE id = $2
        RETURNING *
    """
    row = await self.db.fetchrow(query, task_id, block_id)
    return TimeBlock(**row)

async def _check_conflicts(self, start: datetime, end: datetime, exclude_id: str = None) -> List[TimeBlock]:
    """Check for overlapping time blocks."""
    query = """
        SELECT * FROM time_blocks
        WHERE start_time < $2 AND end_time > $1
        AND ($3::uuid IS NULL OR id != $3)
    """
    rows = await self.db.fetch(query, start, end, exclude_id)
    return [TimeBlock(**row) for row in rows]
```

### Database Migration

**Create:** `backend/lifeos/assets/migrations/001_create_assets.py`

```python
"""Create assets table for universal asset system."""

up = """
CREATE TABLE IF NOT EXISTS assets (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    asset_type VARCHAR(20) NOT NULL,
    title VARCHAR(500) NOT NULL,
    description TEXT,

    -- Classification
    pillar VARCHAR(10),
    priority VARCHAR(10) DEFAULT 'p2',
    energy_level VARCHAR(10),

    -- Timing
    do_date TIMESTAMP WITH TIME ZONE,
    due_date TIMESTAMP WITH TIME ZONE,
    duration_minutes INTEGER,
    recurring JSONB,

    -- Status
    status VARCHAR(20) DEFAULT 'pending',
    progress FLOAT DEFAULT 0.0,
    streak INTEGER,

    -- Relations
    parent_id UUID REFERENCES assets(id) ON DELETE SET NULL,
    project_id UUID,
    linked_assets UUID[] DEFAULT '{}',

    -- Metadata
    tags TEXT[] DEFAULT '{}',
    context JSONB DEFAULT '{}',
    source VARCHAR(50) DEFAULT 'manual',

    -- Timestamps
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    completed_at TIMESTAMP WITH TIME ZONE
);

-- Indexes for common queries
CREATE INDEX idx_assets_type ON assets(asset_type);
CREATE INDEX idx_assets_pillar ON assets(pillar);
CREATE INDEX idx_assets_status ON assets(status);
CREATE INDEX idx_assets_do_date ON assets(do_date);
CREATE INDEX idx_assets_due_date ON assets(due_date);
CREATE INDEX idx_assets_project ON assets(project_id);

-- Full-text search index
CREATE INDEX idx_assets_search ON assets USING gin(
    to_tsvector('english', title || ' ' || COALESCE(description, ''))
);
"""

down = """
DROP TABLE IF EXISTS assets CASCADE;
"""
```

### Compatibility Layer

**Create:** `backend/lifeos/assets/compatibility.py`

```python
"""Backward compatibility layer - maps old endpoints to new Asset system."""

from typing import List
from .models import Asset, AssetType, AssetCreate
from .service import AssetService

class TaskCompatibility:
    """Makes old /api/tasks endpoints work with new Asset system."""

    def __init__(self, asset_service: AssetService):
        self.assets = asset_service

    async def create_task(self, task_data: dict) -> Asset:
        """Create task via Asset system."""
        asset = AssetCreate(
            asset_type=AssetType.TASK,
            title=task_data["title"],
            description=task_data.get("description"),
            pillar=task_data.get("pillar"),
            priority=task_data.get("priority", "p2"),
            due_date=task_data.get("due_date"),
            tags=task_data.get("tags", []),
        )
        return await self.assets.create(asset)

    async def list_tasks(self) -> List[Asset]:
        """List tasks via Asset system."""
        from .models import AssetFilters
        filters = AssetFilters(asset_type=AssetType.TASK)
        return await self.assets.list(filters)


class HabitCompatibility:
    """Makes old /api/habits endpoints work with new Asset system."""

    def __init__(self, asset_service: AssetService):
        self.assets = asset_service

    async def create_habit(self, habit_data: dict) -> Asset:
        """Create habit via Asset system."""
        asset = AssetCreate(
            asset_type=AssetType.HABIT,
            title=habit_data["title"],
            description=habit_data.get("description"),
            pillar=habit_data.get("pillar"),
            recurring=habit_data.get("recurring"),
        )
        return await self.assets.create(asset)

    async def complete_habit(self, habit_id: str) -> Asset:
        """Record habit completion (increments streak)."""
        # This is a daily check-in, not a one-time completion
        asset = await self.assets.get(habit_id)
        if asset:
            new_streak = (asset.streak or 0) + 1
            from .models import AssetUpdate
            update = AssetUpdate(streak=new_streak)
            return await self.assets.update(habit_id, update)
        return None


# Similar for GoalCompatibility...
```

---

## When Done

### 1. Update Linear Issues to "In Review"

Use this GraphQL mutation:

```graphql
mutation UpdateIssue($id: String!, $stateId: String!) {
  issueUpdate(id: $id, input: { stateId: $stateId }) {
    success
    issue {
      id
      identifier
      state {
        name
      }
    }
  }
}
```

Get the "In Review" state ID first:
```graphql
query GetStates {
  workflowStates(filter: { team: { id: { eq: "5dcfd649-4f48-47be-8419-5ef5550e5b76" } } }) {
    nodes { id, name, type }
  }
}
```

### 2. Post Handoff Message

```markdown
## Handoff: Agent 1 → Agents 2-5

### Completed
- ✅ OMOS-467: Universal Asset schema in `backend/lifeos/assets/models.py`
- ✅ OMOS-468: Asset Unified API in `backend/lifeos/assets/api.py`
- ✅ OMOS-460: Planner drag-drop in `backend/lifeos/planner/service.py`
- ✅ Database migration applied

### For Dependent Agents
- Import: `from backend.lifeos.assets.models import Asset, AssetType, AssetCreate`
- Service: `from backend.lifeos.assets.service import AssetService`
- API Base: `/api/v1/assets`

### API Endpoints
- `GET /api/v1/assets` - List with filters
- `POST /api/v1/assets` - Create any type
- `GET /api/v1/assets/today` - Today's items
- `GET /api/v1/assets/{id}` - Single asset
- `PATCH /api/v1/assets/{id}` - Update
- `POST /api/v1/assets/{id}/complete` - Mark complete

### Notes
- Backward compatibility layer in `compatibility.py`
- Existing `/api/tasks`, `/api/habits`, `/api/goals` still work
```

---

## Testing Checklist

Before marking complete:

- [ ] All models pass Pydantic validation
- [ ] Migration runs without errors
- [ ] CRUD operations work for all asset types
- [ ] Filters work correctly
- [ ] Search returns relevant results
- [ ] `/api/v1/assets/today` returns today's items
- [ ] Compatibility layer works with old endpoints
- [ ] Planner drag-drop operations functional
- [ ] No regressions in existing functionality

---

**You are critical to this sprint. Other agents are waiting on you. Build solid, test thoroughly, document clearly.**
