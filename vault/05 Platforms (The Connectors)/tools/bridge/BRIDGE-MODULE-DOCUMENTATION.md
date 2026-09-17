# LifeOS Bridge Module - Complete Documentation

> **Purpose:** Connect all LifeOS features through Agno's Learning Machine, enabling AI to understand relationships, detect patterns, and make predictions.

---

## Table of Contents

1. [Overview](#overview)
2. [Usage Overview](#usage-overview)
   - [Backend Usage](#backend-usage)
   - [Frontend Usage (Flutter)](#frontend-usage-flutter)
3. [Architecture](#architecture)
4. [Module Structure](#module-structure)
5. [Models Reference](#models-reference)
6. [Service API](#service-api)
7. [REST API Endpoints](#rest-api-endpoints)
8. [AI Tools for Legacy Agent](#ai-tools-for-legacy-agent)
9. [Auto-Registration Triggers](#auto-registration-triggers)
10. [Intelligence Scores (PI/LI/WI)](#intelligence-scores-piliwi)
11. [Pattern Detection](#pattern-detection)
12. [Predictions](#predictions)
13. [Family Support](#family-support)
14. [Integration Guide](#integration-guide)
15. [Examples](#examples)
16. [Flutter Frontend Integration](#flutter-frontend-integration)
17. [Service Integration (Automatic Triggers)](#service-integration-automatic-triggers)

---

## Overview

### The Problem We Solved

Before the Bridge, OneMind OS had **11 isolated LifeOS modules** (~44,856 lines of code) that didn't communicate with each other:

```
Tasks  Habits  Goals  Journal  Routines  Pomodoro
Calendar  Meals  Planner  Scheduling  Notifications
```

And a **Learning Machine** with powerful knowledge graph capabilities that wasn't connected to LifeOS:

```
learning_entities          → Store ANY type
learning_entity_facts      → Facts with confidence scores
learning_entity_relations  → Knowledge graph connections
learning_profiles          → Structured user/family profiles
```

**The Gap:** These systems didn't talk to each other!

### The Solution: Bridge Layer

Instead of creating 43 new database tables (as originally planned), we built a **Bridge Layer** (~1,970 lines) that connects everything:

```
┌─────────────────────────────────────────────────────────────┐
│                    CONNECTED LIFE OS                         │
├─────────────────────────────────────────────────────────────┤
│  ┌─────────────────────────────────────────────────────┐    │
│  │            LEARNING MACHINE (Hub)                    │    │
│  │  ┌──────────┐  ┌──────────┐  ┌──────────────────┐   │    │
│  │  │ Entities │──│Relations │──│ Facts + Context  │   │    │
│  │  └────┬─────┘  └────┬─────┘  └────────┬─────────┘   │    │
│  └───────┼─────────────┼─────────────────┼─────────────┘    │
│          │             │                 │                   │
│    ┌─────┴─────────────┴─────────────────┴─────┐            │
│    │           BRIDGE LAYER (NEW)               │            │
│    └─────┬─────────────┬─────────────────┬─────┘            │
│          │             │                 │                   │
│  ┌───────▼───┐   ┌─────▼────┐    ┌──────▼─────┐            │
│  │   Tasks   │   │  Habits  │    │   Goals    │   ...      │
│  └───────────┘   └──────────┘    └────────────┘            │
└─────────────────────────────────────────────────────────────┘
```

### Key Capabilities

| Feature | Description |
|---------|-------------|
| **Unified Context** | Single API call returns LifeOS data + facts + relations + predictions |
| **Knowledge Graph** | Connect any entity to any other (task → goal → child → project) |
| **Pattern Detection** | Automatically find habits patterns, task patterns, correlations |
| **Predictions** | Predict task/habit completion likelihood based on patterns |
| **Intelligence Scores** | PI (Performance), LI (Legacy), WI (Wealth) calculated from real data |
| **Family Support** | Multi-user profiles with parent/child relationships |
| **Auto-Registration** | LifeOS items automatically become Learning entities |

---

## Usage Overview

The Bridge module can be used in three ways from the backend and four ways from the frontend.

### Backend Usage

#### 1. Using the REST API Directly

Make HTTP calls to `/api/bridge/*` endpoints. Best for external integrations or simple queries.

```bash
# Get intelligence scores
curl "http://localhost:8081/api/bridge/scores?user_id=zeus"

# Get entity context
curl "http://localhost:8081/api/bridge/context/task/abc-123?user_id=zeus"

# Get day predictions
curl "http://localhost:8081/api/bridge/predictions/today?user_id=zeus"

# Link two entities
curl -X POST "http://localhost:8081/api/bridge/link?user_id=zeus" \
  -H "Content-Type: application/json" \
  -d '{"from_type":"task","from_id":"abc","to_type":"goal","to_id":"xyz","relation_type":"supports"}'

# Add a fact to an entity
curl -X POST "http://localhost:8081/api/bridge/fact/habit/morning-routine?user_id=zeus" \
  -H "Content-Type: application/json" \
  -d '{"fact":"Works best before 6am","source":"user","confidence":1.0}'
```

**Full API reference:** See [REST API Endpoints](#rest-api-endpoints)

#### 2. Using the Service Directly in Python

Import and use `LifeOSBridgeService` for programmatic access within the backend.

```python
from backend.lifeos.bridge import get_bridge_service

# Get service instance
service = await get_bridge_service()

# Register an entity
entity_id = await service.register_entity(
    user_id="zeus",
    entity_type=BridgeEntityType.TASK,
    lifeos_table="inbox_items",
    lifeos_id=task_uuid,
    name="Review report",
    pillar=Pillar.GE,
)

# Link entities
await service.link_entities(
    user_id="zeus",
    from_type=BridgeEntityType.TASK,
    from_id="uuid-or-name",
    to_type=BridgeEntityType.GOAL,
    to_id="uuid-or-name",
    relation_type=RelationType.SUPPORTS,
)

# Get unified context
context = await service.get_full_context("zeus", BridgeEntityType.TASK, "uuid")

# Calculate intelligence scores
scores = await service.calculate_intelligence_scores("zeus")

# Detect patterns
patterns = await service.detect_patterns("zeus")

# Predict day
prediction = await service.predict_day("zeus", date.today())

# Get family overview
family = await service.get_family_overview("zeus")
```

**Full API reference:** See [Service API](#service-api)

#### 3. Using Auto-Registration Triggers

Let Bridge automatically register entities when LifeOS items are created/updated.

```python
# In your service file (tasks/habits/goals):
try:
    from backend.lifeos.bridge.triggers import BridgeTriggers
    BRIDGE_ENABLED = True
except ImportError:
    BRIDGE_ENABLED = False

# After creating an item
async def create(self, data, user_id):
    item = await self._create_item(data)

    if BRIDGE_ENABLED:
        try:
            await BridgeTriggers.on_task_created(item, user_id)
        except Exception as e:
            logger.warning(f"Bridge trigger failed: {e}")

    return item

# After completing an item
async def complete(self, item_id, user_id):
    item = await self._complete_item(item_id)

    if BRIDGE_ENABLED:
        try:
            await BridgeTriggers.on_task_completed(item, user_id, actual_minutes)
        except Exception as e:
            logger.warning(f"Bridge trigger failed: {e}")

    return item
```

**Available triggers:**
| Trigger | When to Call |
|---------|--------------|
| `on_task_created` | After task creation |
| `on_task_completed` | After task completion |
| `on_habit_created` | After habit creation |
| `on_habit_completed` | After habit completion |
| `on_habit_skipped` | After habit skip |
| `on_goal_created` | After goal creation |
| `on_goal_progress_updated` | After goal progress change |
| `on_journal_entry_created` | After journal entry |
| `on_pomodoro_completed` | After pomodoro session |

**Full reference:** See [Auto-Registration Triggers](#auto-registration-triggers)

---

### Frontend Usage (Flutter)

The Flutter frontend at `frontend/lib/lifeos/bridge/` provides a complete implementation.

#### 1. Create Flutter Models

All data models are in `bridge_models.dart`:

```dart
// Import models
import 'package:frontend/lifeos/bridge/models/bridge_models.dart';

// Key models available:
IntelligenceScores   // PI/LI/WI scores with components
PillarScore          // Individual pillar score with trend
DayPrediction        // Day predictions with tasks/habits/blockers
TaskPrediction       // Task completion likelihood
HabitPrediction      // Habit completion likelihood
UserPatterns         // Detected habit/task patterns
FamilyOverview       // Family members and tasks
EntityContext        // Unified entity context
```

**Full model reference:** See [Flutter Frontend Integration](#flutter-frontend-integration)

#### 2. Create a Bridge API Service

The HTTP client in `bridge_api.dart`:

```dart
// Import API service
import 'package:frontend/lifeos/bridge/services/bridge_api.dart';

// Create instance
final api = BridgeApi(baseUrl: 'http://localhost:8081');

// Make API calls
final scores = await api.getIntelligenceScores(userId: 'zeus');
final predictions = await api.getTodayPredictions(userId: 'zeus');
final patterns = await api.getUserPatterns(userId: 'zeus');
final family = await api.getFamilyOverview(userId: 'zeus');
final context = await api.getEntityContext(
  entityType: 'task',
  entityId: 'abc-123',
  userId: 'zeus',
);
```

#### 3. Create a Provider

Riverpod state management in `bridge_provider.dart`:

```dart
// Import providers
import 'package:frontend/lifeos/bridge/providers/bridge_provider.dart';

// In a ConsumerWidget:
class MyWidget extends ConsumerWidget {
  @override
  Widget build(BuildContext context, WidgetRef ref) {
    // Watch intelligence scores
    final scoresState = ref.watch(currentUserScoresProvider);

    // Watch predictions
    final predictionsState = ref.watch(currentUserPredictionsProvider);

    // Watch patterns
    final patternsState = ref.watch(currentUserPatternsProvider);

    // Watch family
    final familyState = ref.watch(currentFamilyOverviewProvider);

    // For specific user (not current)
    final scores = ref.watch(intelligenceScoresProvider('zeus'));

    // Combined dashboard state
    final dashboard = ref.watch(bridgeDashboardProvider);
  }
}
```

#### 4. Create UI Widgets

Pre-built widgets in `widgets/`:

```dart
// Import widgets
import 'package:frontend/lifeos/bridge/widgets/widgets.dart';

// Use pre-built cards
IntelligenceScoresCard(userId: 'zeus', compact: false)
DayPredictionsCard(userId: 'zeus', showBlockers: true, showSuggestions: true)

// Or use the full dashboard screen
import 'package:frontend/lifeos/bridge/screens/bridge_dashboard_screen.dart';

BridgeDashboardScreen(userId: 'zeus')
// Has 4 tabs: Overview, Predictions, Patterns, Family
```

**Navigation:**
- Route: `/bridge`
- Menu: "Intelligence" in UIO section of drawer

**Full Flutter reference:** See [Flutter Frontend Integration](#flutter-frontend-integration)

---

## Architecture

### Data Flow

```
┌────────────────────────────────────────────────────────────┐
│                    METRICS DATA FLOW                        │
├────────────────────────────────────────────────────────────┤
│                                                             │
│  LIFEOS TABLES (Source of Truth)                           │
│  ┌─────────────────────────────────────────────────────┐   │
│  │ tasks: status, priority, time_estimate, completed_at│   │
│  │ habits: streak, grade, completion_rate, frequency   │   │
│  │ goals: progress, target_date, goal_type            │   │
│  │ projects: task_count, completed_count               │   │
│  │ habit_completions: date, status, mood, energy      │   │
│  │ journal: mood, energy, gratitude, tags             │   │
│  │ pomodoro: duration, focus_rating, interruptions    │   │
│  └───────────────────────┬─────────────────────────────┘   │
│                          │                                  │
│                          ▼                                  │
│  BRIDGE LAYER (Aggregation + Registration)                 │
│  ┌─────────────────────────────────────────────────────┐   │
│  │ - Queries LifeOS tables for real metrics            │   │
│  │ - Registers entities in Learning Machine           │   │
│  │ - Calculates patterns from historical data         │   │
│  │ - Creates/maintains relations                      │   │
│  └───────────────────────┬─────────────────────────────┘   │
│                          │                                  │
│                          ▼                                  │
│  LEARNING MACHINE (Context + Patterns)                     │
│  ┌─────────────────────────────────────────────────────┐   │
│  │ entities: type, name, pillar, metadata (cached)    │   │
│  │ facts: patterns, correlations, inferences          │   │
│  │ relations: supports, requires, contains, assigned  │   │
│  └───────────────────────┬─────────────────────────────┘   │
│                          │                                  │
│                          ▼                                  │
│  AI CONTEXT API (Unified View)                             │
│  ┌─────────────────────────────────────────────────────┐   │
│  │ get_full_context() returns:                         │   │
│  │  - LIVE metrics from LifeOS (streak, progress, etc)│   │
│  │  - Relationships from Learning Machine              │   │
│  │  - Patterns/facts from Learning Machine            │   │
│  │  - Predictions computed on-demand                   │   │
│  └─────────────────────────────────────────────────────┘   │
│                                                             │
└────────────────────────────────────────────────────────────┘
```

### Database Tables Used

The Bridge does NOT create new tables. It uses existing ones:

| Table | Source | Purpose |
|-------|--------|---------|
| `learning_entities` | Learning Machine | Store entity references |
| `learning_entity_facts` | Learning Machine | Store facts with confidence |
| `learning_entity_relations` | Learning Machine | Knowledge graph edges |
| `learning_profiles` | Learning Machine | User/family profiles |
| `inbox_items` | LifeOS Tasks | Task data (where type='task') |
| `habits` | LifeOS Habits | Habit data |
| `habit_completions` | LifeOS Habits | Completion history |
| `goals` | LifeOS Goals | Goal data |
| `journal_entries` | LifeOS Journal | Mood/energy data |
| `pomodoro_sessions` | LifeOS Pomodoro | Focus session data |
| `calendar_events` | LifeOS Calendar | Meeting data |

---

## Module Structure

```
backend/lifeos/bridge/
├── __init__.py          # Module exports
├── models.py            # Pydantic schemas (~310 lines)
├── service.py           # Core bridge service (~770 lines)
├── api.py               # REST API endpoints (~410 lines)
└── triggers.py          # Auto-registration triggers (~450 lines)
```

### File Purposes

| File | Purpose |
|------|---------|
| `__init__.py` | Exports `LifeOSBridgeService`, `get_bridge_service`, `BridgeTriggers`, backfill functions |
| `models.py` | All Pydantic models for requests/responses, entity types, pillars, relations |
| `service.py` | Core business logic: registration, linking, context, scores, patterns, predictions |
| `api.py` | FastAPI router with all REST endpoints at `/api/bridge/*` |
| `triggers.py` | Event triggers for auto-registering LifeOS items, plus backfill utilities |

---

## Models Reference

### Entity Types (`BridgeEntityType`)

All entity types supported by the Bridge:

```python
class BridgeEntityType(str, Enum):
    # Learning Machine original types
    PERSON = "person"
    PROJECT = "project"
    COMPANY = "company"
    LOCATION = "location"
    CONCEPT = "concept"

    # LifeOS types
    TASK = "task"
    HABIT = "habit"
    GOAL = "goal"
    JOURNAL = "journal"
    ROUTINE = "routine"
    POMODORO = "pomodoro"
    CALENDAR_EVENT = "calendar_event"
    MEAL = "meal"
    RECIPE = "recipe"
    TIME_BLOCK = "time_block"
    SCHEDULED_TASK = "scheduled_task"

    # Codex types (future)
    CHILD = "child"
    CONTACT = "contact"
    PROPERTY = "property"
    ACCOUNT = "account"
    BUSINESS = "business"
    VEHICLE = "vehicle"
    MAINTENANCE = "maintenance"
```

### Pillars (`Pillar`)

The Four Pillars of OneMind OS:

```python
class Pillar(str, Enum):
    HP = "hp"  # Holistic Performance (Health)
    LE = "le"  # Legacy Evolution (Family/Education)
    GE = "ge"  # Generational Empire (Business/Wealth)
    IT = "it"  # Innovative Technology (Integration)
```

### Relation Types (`RelationType`)

Standard relationship types for the knowledge graph:

```python
class RelationType(str, Enum):
    # Hierarchical
    PARENT_OF = "parent_of"
    CHILD_OF = "child_of"
    PART_OF = "part_of"
    CONTAINS = "contains"

    # Goal/Task relationships
    SUPPORTS = "supports"
    REQUIRES = "requires"
    BLOCKS = "blocks"
    BLOCKED_BY = "blocked_by"
    CONTRIBUTES_TO = "contributes_to"

    # Assignment
    ASSIGNED_TO = "assigned_to"
    OWNED_BY = "owned_by"
    MANAGED_BY = "managed_by"

    # General
    RELATES_TO = "relates_to"
    FOLLOWS = "follows"
    PRECEDES = "precedes"
```

### Fact Sources (`FactSource`)

Where facts come from:

```python
class FactSource(str, Enum):
    USER = "user"          # Explicitly stated by user
    SYSTEM = "system"      # Recorded by system (completion, etc.)
    INFERRED = "inferred"  # AI-detected pattern
    IMPORTED = "imported"  # Imported from external source
```

### Key Response Models

#### EntityContext

The unified context response combining all data sources:

```python
class EntityContext(BaseModel):
    id: UUID                              # Entity UUID
    type: BridgeEntityType                # Entity type
    name: str                             # Entity name
    pillar: Optional[Pillar]              # Associated pillar
    lifeos_data: Dict[str, Any]           # Live data from LifeOS tables
    relations: List[RelatedEntity]        # Connected entities
    facts: List[Dict[str, Any]]           # Facts with confidence
    predictions: Dict[str, Any]           # Computed predictions
    markdown: str                         # AI-friendly markdown format
```

#### IntelligenceScores

Complete intelligence scores for a user:

```python
class IntelligenceScores(BaseModel):
    user_id: str
    calculated_at: datetime
    pi: PillarScore      # Performance Intelligence (HP)
    li: PillarScore      # Legacy Intelligence (LE)
    wi: PillarScore      # Wealth Intelligence (GE)
    unified_score: float # Average of all three (0-100)
    suggestions: List[Dict[str, Any]]  # Improvement suggestions
```

#### DayPrediction

Complete prediction for a day:

```python
class DayPrediction(BaseModel):
    user_id: str
    date: date
    predicted_energy: float        # 0-10 scale
    meeting_hours: float           # Calendar load
    tasks: List[TaskPrediction]    # Per-task predictions
    habits: List[HabitPrediction]  # Per-habit predictions
    blockers: List[str]            # Identified blockers
    suggestions: List[str]         # Optimization suggestions
```

#### UserPatterns

Detected patterns from historical data:

```python
class UserPatterns(BaseModel):
    user_id: str
    analyzed_at: datetime
    habit_patterns: List[HabitPattern]    # Per-habit patterns
    task_patterns: Optional[TaskPattern]  # Task patterns
```

---

## Service API

### LifeOSBridgeService

The core service class. Get an instance via:

```python
from backend.lifeos.bridge import get_bridge_service

service = await get_bridge_service()
```

### Methods

#### Entity Registration

```python
async def register_entity(
    user_id: str,
    entity_type: BridgeEntityType,
    lifeos_table: str,
    lifeos_id: UUID,
    name: str,
    description: Optional[str] = None,
    pillar: Optional[Pillar] = None,
    metadata: dict = {},
) -> Optional[UUID]:
    """
    Register a LifeOS item as a Learning entity.

    Returns entity UUID if successful, None if already exists.
    """
```

#### Entity Linking

```python
async def link_entities(
    user_id: str,
    from_type: BridgeEntityType,
    from_id: str,
    to_type: BridgeEntityType,
    to_id: str,
    relation_type: Union[RelationType, str],
    description: Optional[str] = None,
    metadata: dict = {},
) -> bool:
    """
    Create a relationship between two entities.

    from_id and to_id can be UUIDs or names (will search).
    """
```

#### Adding Facts

```python
async def add_fact(
    user_id: str,
    entity_type: BridgeEntityType,
    entity_id: str,
    fact: str,
    source: FactSource = FactSource.USER,
    confidence: float = 1.0,
) -> bool:
    """
    Add a fact to an entity.

    Facts accumulate over time for pattern detection.
    """
```

#### Getting Context

```python
async def get_full_context(
    user_id: str,
    entity_type: BridgeEntityType,
    entity_id: str,
) -> Optional[EntityContext]:
    """
    Get unified context for an entity.

    Combines:
    - Live LifeOS data
    - Learning Machine facts
    - Related entities
    - Predictions

    Returns both structured data and markdown.
    """
```

#### Intelligence Scores

```python
async def calculate_intelligence_scores(
    user_id: str,
) -> IntelligenceScores:
    """
    Calculate PI/LI/WI scores from LifeOS data.

    - PI = HP pillar metrics (habits, energy, goals)
    - LI = LE pillar metrics (family, kids, home)
    - WI = GE pillar metrics (business, tasks, goals)
    """
```

#### Pattern Detection

```python
async def detect_patterns(
    user_id: str,
) -> UserPatterns:
    """
    Analyze historical data to find patterns.

    Habit patterns:
    - Best days for completion
    - Best times
    - Energy/mood correlations
    - Skip reasons

    Task patterns:
    - Estimation accuracy
    - Completion by day
    - Meeting impact
    """
```

#### Day Predictions

```python
async def predict_day(
    user_id: str,
    target_date: date,
) -> DayPrediction:
    """
    Generate predictions for a specific day.

    Returns:
    - Predicted energy level
    - Task completion likelihoods
    - Habit completion likelihoods
    - Blockers and suggestions
    """
```

#### Family Overview

```python
async def get_family_overview(
    user_id: str,
) -> FamilyOverview:
    """
    Get overview of family members and their tasks.

    Returns members, today's tasks per member, shared goals.
    """
```

---

## REST API Endpoints

Base URL: `/api/bridge`

### Context Endpoints

| Method | Path | Description |
|--------|------|-------------|
| GET | `/context/{entity_type}/{entity_id}` | Get full context for an entity |
| GET | `/context/{entity_type}/{entity_id}/markdown` | Get context as markdown only |

**Query Parameters:**
- `user_id` (required): User ID

**Example:**
```bash
GET /api/bridge/context/task/abc-123?user_id=zeus
```

### Registration Endpoints

| Method | Path | Description |
|--------|------|-------------|
| POST | `/register` | Register a LifeOS item as Learning entity |

**Request Body:**
```json
{
  "entity_type": "task",
  "lifeos_table": "inbox_items",
  "lifeos_id": "uuid-here",
  "name": "Review quarterly report",
  "description": "Q1 financial review",
  "pillar": "ge",
  "metadata": {"priority": "high"}
}
```

### Linking Endpoints

| Method | Path | Description |
|--------|------|-------------|
| POST | `/link` | Create relationship between entities |

**Request Body:**
```json
{
  "from_type": "task",
  "from_id": "task-uuid-or-name",
  "to_type": "goal",
  "to_id": "goal-uuid-or-name",
  "relation_type": "supports",
  "description": "This task directly contributes to the goal"
}
```

### Fact Endpoints

| Method | Path | Description |
|--------|------|-------------|
| POST | `/fact/{entity_type}/{entity_id}` | Add a fact to an entity |

**Request Body:**
```json
{
  "fact": "Completed 30 minutes under estimate",
  "source": "system",
  "confidence": 1.0
}
```

### Intelligence Score Endpoints

| Method | Path | Description |
|--------|------|-------------|
| GET | `/scores` | Get PI/LI/WI intelligence scores |

**Response:**
```json
{
  "user_id": "zeus",
  "calculated_at": "2026-01-28T10:00:00",
  "pi": {"pillar": "hp", "score": 82, "components": {...}, "trend": "up"},
  "li": {"pillar": "le", "score": 71, "components": {...}, "trend": "down"},
  "wi": {"pillar": "ge", "score": 76, "components": {...}, "trend": "stable"},
  "unified_score": 76.3,
  "suggestions": [...]
}
```

### Pattern Endpoints

| Method | Path | Description |
|--------|------|-------------|
| GET | `/patterns` | Get detected user patterns |

### Prediction Endpoints

| Method | Path | Description |
|--------|------|-------------|
| GET | `/predictions/{target_date}` | Get predictions for a specific date |
| GET | `/predictions/today` | Get predictions for today |

### Family Endpoints

| Method | Path | Description |
|--------|------|-------------|
| GET | `/family` | Get family overview |

### Stats Endpoints

| Method | Path | Description |
|--------|------|-------------|
| GET | `/stats` | Get Bridge statistics |

### Backfill Endpoints

| Method | Path | Description |
|--------|------|-------------|
| POST | `/backfill` | Backfill existing LifeOS data |

**Response:**
```json
{
  "success": true,
  "tasks": 50,
  "habits": 12,
  "goals": 8,
  "total": 70,
  "message": "Backfilled 70 entities"
}
```

---

## AI Tools for Legacy Agent

The Bridge provides 7 tools for the Legacy agent to interact with the knowledge graph.

### Tool: `unified_context`

Get complete context for any entity.

```python
@tool
async def unified_context(
    entity_type: str,
    identifier: str,
    user_id: str = "zeus",
) -> str:
```

**Parameters:**
- `entity_type`: Type of entity (task, habit, goal, child, project, etc.)
- `identifier`: Entity ID (UUID) or name to search for
- `user_id`: User ID (defaults to "zeus")

**Returns:** Markdown-formatted context

**Examples:**
```python
unified_context(entity_type="task", identifier="Review quarterly report")
unified_context(entity_type="habit", identifier="Morning routine")
unified_context(entity_type="goal", identifier="Q1 Revenue Target")
unified_context(entity_type="child", identifier="Noah")
```

### Tool: `link_entities`

Create a relationship between two entities.

```python
@tool
async def link_entities(
    from_type: str,
    from_id: str,
    to_type: str,
    to_id: str,
    relation: str,
    description: str = "",
    user_id: str = "zeus",
) -> str:
```

**Parameters:**
- `from_type`: Source entity type
- `from_id`: Source entity ID or name
- `to_type`: Target entity type
- `to_id`: Target entity ID or name
- `relation`: Relationship type (supports, requires, part_of, assigned_to, etc.)
- `description`: Optional description
- `user_id`: User ID

**Examples:**
```python
link_entities(
    from_type="task", from_id="Review report",
    to_type="goal", to_id="Q1 Revenue",
    relation="supports"
)

link_entities(
    from_type="task", from_id="Math homework",
    to_type="child", to_id="Noah",
    relation="assigned_to"
)
```

### Tool: `add_entity_context`

Add a fact to an entity.

```python
@tool
async def add_entity_context(
    entity_type: str,
    entity_id: str,
    fact: str,
    source: str = "user",
    confidence: float = 1.0,
    user_id: str = "zeus",
) -> str:
```

**Parameters:**
- `entity_type`: Entity type
- `entity_id`: Entity ID or name
- `fact`: The fact to record
- `source`: Source of fact (user, system, inferred)
- `confidence`: Confidence level 0.0-1.0

**Examples:**
```python
add_entity_context(
    entity_type="habit", entity_id="Morning routine",
    fact="Works best when started before 6am"
)

add_entity_context(
    entity_type="child", entity_id="Noah",
    fact="Learns math better with visual examples"
)
```

### Tool: `get_intelligence_scores`

Get PI/LI/WI scores.

```python
@tool
async def get_intelligence_scores(user_id: str = "zeus") -> str:
```

**Returns:** Formatted scores with breakdown and suggestions

### Tool: `get_day_predictions`

Get predictions for a day.

```python
@tool
async def get_day_predictions(
    target_date: str = "",
    user_id: str = "zeus",
) -> str:
```

**Parameters:**
- `target_date`: Date in YYYY-MM-DD format (defaults to today)
- `user_id`: User ID

**Returns:** Predictions with task/habit likelihoods

### Tool: `get_user_patterns`

Get detected patterns.

```python
@tool
async def get_user_patterns(user_id: str = "zeus") -> str:
```

**Returns:** Habit patterns, task patterns, correlations

### Tool: `get_family_overview`

Get family overview.

```python
@tool
async def get_family_overview(user_id: str = "zeus") -> str:
```

**Returns:** Family members, today's tasks, shared goals

---

## Auto-Registration Triggers

The `BridgeTriggers` class provides event handlers for auto-registering LifeOS items.

### Usage in LifeOS Services

```python
from backend.lifeos.bridge.triggers import BridgeTriggers

# In TaskService.create():
task = await self._create_task_in_db(data)
await BridgeTriggers.on_task_created(task, user_id)

# In HabitService.complete():
await BridgeTriggers.on_habit_completed(habit, user_id, mood, energy)
```

### Available Triggers

| Trigger | When to Call | What it Does |
|---------|--------------|--------------|
| `on_task_created` | After task creation | Registers task entity, links to project |
| `on_task_completed` | After task completion | Adds completion fact with timing |
| `on_habit_created` | After habit creation | Registers habit entity |
| `on_habit_completed` | After habit completion | Adds completion fact with mood/energy |
| `on_habit_skipped` | After habit skip | Records skip for pattern analysis |
| `on_goal_created` | After goal creation | Registers goal entity, links to parent |
| `on_goal_progress_updated` | After progress change | Adds progress fact |
| `on_journal_entry_created` | After journal entry | Registers entry for correlations |
| `on_pomodoro_completed` | After focus session | Links focus data to task |

### Trigger Details

#### on_task_created

```python
@staticmethod
async def on_task_created(
    task: Dict[str, Any],
    user_id: str = "zeus",
) -> Optional[UUID]:
    """
    Register task as Learning entity.
    Auto-links to project if project_id present.

    Args:
        task: Dict with id, title, pillar, project_id, etc.
        user_id: User who created

    Returns:
        Entity UUID if registered
    """
```

#### on_habit_completed

```python
@staticmethod
async def on_habit_completed(
    habit: Dict[str, Any],
    user_id: str = "zeus",
    mood: Optional[int] = None,
    energy: Optional[int] = None,
    notes: Optional[str] = None,
) -> bool:
    """
    Record habit completion fact.

    Fact includes: day of week, mood, energy, notes
    Used for pattern detection.
    """
```

### Backfill Functions

For migrating existing data:

```python
from backend.lifeos.bridge.triggers import (
    backfill_tasks,
    backfill_habits,
    backfill_goals,
    backfill_all,
)

# Backfill all existing data
result = await backfill_all(user_id="zeus")
# Returns: {"tasks": 50, "habits": 12, "goals": 8, "total": 70}

# Or backfill specific types
tasks_count = await backfill_tasks(user_id="zeus", limit=100)
habits_count = await backfill_habits(user_id="zeus", limit=100)
goals_count = await backfill_goals(user_id="zeus", limit=100)
```

---

## Intelligence Scores (PI/LI/WI)

### Score Definitions

| Score | Pillar | What it Measures |
|-------|--------|------------------|
| **PI** (Performance Intelligence) | HP | Health habits, energy, mood, focus, HP goals |
| **LI** (Legacy Intelligence) | LE | Family tasks, kids progress, home maintenance, LE goals |
| **WI** (Wealth Intelligence) | GE | Business tasks, revenue goals, work habits, GE goals |
| **Unified Score** | All | Average of PI + LI + WI |

### How Scores Are Calculated

```python
async def _calculate_pi(self, user_id: str) -> PillarScore:
    """
    Performance Intelligence (HP Pillar)

    Components (weights):
    - health_habits (40%): HP habit completion rate
    - energy_mood (20%): Journal energy/mood averages
    - focus (20%): Pomodoro focus ratings
    - goals (20%): HP goal progress
    """

async def _calculate_li(self, user_id: str) -> PillarScore:
    """
    Legacy Intelligence (LE Pillar)

    Components (weights):
    - family (35%): LE task/habit completion
    - kids (25%): Children's task completion
    - home (20%): Home maintenance on-schedule
    - goals (20%): LE goal progress
    """

async def _calculate_wi(self, user_id: str) -> PillarScore:
    """
    Wealth Intelligence (GE Pillar)

    Components (weights):
    - business_tasks (30%): GE task completion
    - business_goals (30%): Revenue goal progress
    - work_habits (20%): GE habit completion
    - finance (20%): Future financial metrics
    """
```

### Example Score Output

```markdown
## Intelligence Scores - Zeus

**Calculated:** 2026-01-28T10:00:00

**PI (Performance)**: 82/100 ↑
  - health_habits: 87%
  - energy_mood: 78%
  - focus: 85%
  - goals: 75%

**LI (Legacy)**: 71/100 ↓
  - family: 80%
  - kids: 65%
  - home: 60%
  - goals: 78%

**WI (Wealth)**: 76/100 →
  - business_tasks: 82%
  - business_goals: 70%
  - work_habits: 75%
  - finance: N/A

**Unified Score**: 76/100

### Suggestions
1. Complete 2 overdue home maintenance tasks → +5 LI
2. Follow up with Client B → +3 WI
3. Help Aria catch up on reading log → +3 LI
```

---

## Pattern Detection

### Habit Patterns

The system detects these patterns for each habit:

```python
class HabitPattern(BaseModel):
    habit_id: UUID
    habit_name: str
    completion_rate: float       # Overall completion %
    best_days: Dict[str, float]  # Day -> completion rate
    best_times: Dict[str, float] # Time period -> completion rate
    energy_threshold: Optional[float]        # Energy below which skips occur
    skip_rate_below_threshold: Optional[float]
    correlations: List[str]      # "After meditation", "When mood > 7", etc.
```

**Example Pattern:**
```json
{
  "habit_name": "Morning Routine",
  "completion_rate": 0.87,
  "best_days": {
    "Monday": 0.95,
    "Tuesday": 0.92,
    "Wednesday": 0.90,
    "Saturday": 0.60
  },
  "best_times": {
    "morning": 0.92,
    "afternoon": 0.0
  },
  "energy_threshold": 5,
  "skip_rate_below_threshold": 0.4,
  "correlations": [
    "After good sleep (>7h)",
    "When mood >= 7"
  ]
}
```

### Task Patterns

```python
class TaskPattern(BaseModel):
    user_id: str
    estimation_accuracy: float              # Overall accuracy
    estimation_by_pillar: Dict[str, float]  # HP: 1.1, GE: 0.8
    completion_by_day: Dict[str, Dict]      # Day -> {rate, avg_tasks}
    productivity_by_hour: Dict[str, Dict]   # Hour range -> {focus, rate}
    meeting_impact: Dict[str, float]        # Meeting count -> task completion
```

**Example Pattern:**
```json
{
  "estimation_accuracy": 0.85,
  "estimation_by_pillar": {
    "HP": 1.1,
    "LE": 0.9,
    "GE": 0.8
  },
  "completion_by_day": {
    "Monday": {"rate": 0.92, "avg_tasks": 5},
    "Tuesday": {"rate": 0.88, "avg_tasks": 6},
    "Friday": {"rate": 0.75, "avg_tasks": 3}
  },
  "meeting_impact": {
    "0_meetings": 0.95,
    "1-2_meetings": 0.85,
    "3+_meetings": 0.60
  }
}
```

---

## Predictions

### How Predictions Work

```python
async def predict_day(self, user_id: str, target_date: date) -> DayPrediction:
    """
    1. Get calendar events for the day
    2. Get tasks due on that day
    3. Get habits due on that day
    4. Get user's historical patterns
    5. Calculate meeting load
    6. Predict energy level based on patterns
    7. For each task: calculate completion likelihood
    8. For each habit: calculate completion likelihood
    9. Identify blockers
    10. Generate suggestions
    """
```

### Task Prediction

```python
class TaskPrediction(BaseModel):
    task_id: UUID
    task_name: str
    completion_likelihood: float    # 0.0 - 1.0
    optimal_time: Optional[str]     # "9-11am"
    estimated_duration: Optional[str]
    blockers: List[str]            # ["Meeting at 11am", "Low energy predicted"]
    suggestions: List[str]         # ["Start before 10am", "Break into 2 sessions"]
```

### Habit Prediction

```python
class HabitPrediction(BaseModel):
    habit_id: UUID
    habit_name: str
    completion_likelihood: float    # 0.0 - 1.0
    risk_factors: List[str]        # ["3 meetings today", "Low energy predicted"]
    suggestions: List[str]         # ["Do shorter version", "Move to evening"]
```

### Example Prediction Output

```markdown
## Day Predictions: Wednesday, Jan 29, 2026

**Predicted Energy:** 6/10
**Meeting Hours:** 2.5h

### Task Predictions
| Task | Likelihood | Optimal Time | Risk |
|------|------------|--------------|------|
| Review quarterly report | 85% | 8-10am | Meeting at 10 |
| Client B follow-up | 90% | 2-3pm | - |
| Code review | 70% | 4-5pm | Low energy |

### Habit Predictions
| Habit | Likelihood | Risk Factors |
|-------|------------|--------------|
| Morning routine | 92% | - |
| Workout | 75% | Meeting cluster |
| Evening meditation | 88% | - |

### Blockers
- Meeting cluster 10-12:30 → Deep work blocked
- Energy dip predicted 2pm

### Suggestions
- Do quarterly report 8-10am (before meetings)
- Schedule easy tasks 2-3pm (post-meeting recovery)
- Consider 30-min workout instead of 60-min
```

---

## Family Support

### How Family Works

The Bridge supports multi-user families through Learning Profiles:

```python
# Parent profile
zeus_profile = {
    "user_id": "zeus",
    "name": "Zeus Delacruz",
    "role": "parent",
    "four_pillars": {
        "HP": {"weight": 0.3},
        "LE": {"weight": 0.3},
        "GE": {"weight": 0.25},
        "IT": {"weight": 0.15}
    }
}

# Child profile
noah_profile = {
    "user_id": "noah",
    "name": "Noah Delacruz",
    "role": "child",
    "parent_id": "zeus",
    "age": 10,
    "preferences": {
        "best_homework_time": "16:00",
        "attention_span_minutes": 30
    }
}
```

### Family Knowledge Graph

```
              ┌─────────────────┐
              │   Zeus (Parent)  │
              └────────┬────────┘
                       │
         ┌─────────────┼─────────────┐
         │ parent_of   │ parent_of   │
         ▼             ▼             ▼
   ┌──────────┐  ┌──────────┐  ┌──────────┐
   │   Noah   │  │   Aria   │  │   ...    │
   │  age: 10 │  │  age: 7  │  │          │
   └────┬─────┘  └────┬─────┘  └──────────┘
        │             │
   assigned_to   assigned_to
        │             │
        ▼             ▼
   ┌──────────┐  ┌──────────┐
   │Math Task │  │Reading   │
   │due: Fri  │  │Goal      │
   └──────────┘  └──────────┘
```

### Family Overview Response

```python
class FamilyOverview(BaseModel):
    parent_id: str
    members: List[FamilyMember]
    shared_goals: List[Dict[str, Any]]
    today_tasks: Dict[str, List[Dict[str, Any]]]  # member_id -> tasks
```

**Example:**
```markdown
## Family Overview

### Family Members
- **Zeus** - parent
- **Noah** (10) - child
- **Aria** (7) - child

### Today's Tasks

**Noah**
  ⬜ Math homework - Chapter 5
  ⬜ Clean room
  ✅ Reading log

**Aria**
  ⬜ Spelling practice
  ✅ Piano lesson

### Shared Family Goals
- Kids Education: 75%
- Family Chores: 60%
```

---

## Integration Guide

### Adding Bridge to a New LifeOS Service

1. **Import the trigger:**

```python
from backend.lifeos.bridge.triggers import BridgeTriggers
```

2. **Call trigger after create:**

```python
async def create(self, data: TaskCreate, user_id: str) -> Task:
    # Create the item
    task = await self._insert_task(data)

    # Auto-register in Bridge (non-blocking)
    await BridgeTriggers.on_task_created(
        task=task.dict(),
        user_id=user_id,
    )

    return task
```

3. **Call trigger after completion:**

```python
async def complete(self, task_id: UUID, user_id: str) -> Task:
    task = await self._mark_complete(task_id)

    await BridgeTriggers.on_task_completed(
        task=task.dict(),
        user_id=user_id,
        actual_minutes=task.actual_minutes,
    )

    return task
```

### Using Bridge Tools in an Agent

1. **Ensure tools are registered:**

The tools are auto-registered in `backend/agno/tools/registry.py`.

2. **Add tools to agent:**

```python
from backend.agno.tools.registry import resolve_tools

tools = resolve_tools([
    "unified_context",
    "link_entities",
    "add_entity_context",
    "get_intelligence_scores",
    "get_day_predictions",
    "get_user_patterns",
    "get_family_overview",
])

agent = Agent(
    name="Legacy",
    tools=tools,
    ...
)
```

3. **Use in agent instructions:**

```markdown
You have access to the Bridge tools:
- Use unified_context to get complete context about any entity
- Use link_entities to connect related items
- Use add_entity_context to record important facts
- Use get_intelligence_scores for PI/LI/WI scores
- Use get_day_predictions for task/habit predictions
```

### Calling Bridge API from Frontend

```dart
// Flutter example

// Get unified context
final response = await http.get(
  Uri.parse('$baseUrl/api/bridge/context/task/$taskId?user_id=zeus'),
);
final context = EntityContext.fromJson(jsonDecode(response.body));

// Get intelligence scores
final scoresResponse = await http.get(
  Uri.parse('$baseUrl/api/bridge/scores?user_id=zeus'),
);
final scores = IntelligenceScores.fromJson(jsonDecode(scoresResponse.body));

// Link entities
await http.post(
  Uri.parse('$baseUrl/api/bridge/link?user_id=zeus'),
  body: jsonEncode({
    'from_type': 'task',
    'from_id': taskId,
    'to_type': 'goal',
    'to_id': goalId,
    'relation_type': 'supports',
  }),
);
```

---

## Examples

### Example 1: Complete Task Flow

```python
# 1. Task is created in LifeOS
task = await task_service.create(TaskCreate(
    title="Review quarterly report",
    pillar="GE",
    priority="high",
    due_date="2026-02-01",
))

# 2. Bridge auto-registers it (via trigger)
# Creates: learning_entity with type="task"

# 3. User links it to a goal
await link_entities(
    from_type="task", from_id="Review quarterly report",
    to_type="goal", to_id="Q1 Revenue",
    relation="supports"
)
# Creates: learning_entity_relation

# 4. AI queries context later
context = await unified_context(
    entity_type="task",
    identifier="Review quarterly report"
)
# Returns:
# ## Task: Review quarterly report
# **Pillar:** GE
# ### Core Data
# - status: in_progress
# - priority: high
# - due_date: 2026-02-01
# ### Relationships
# - supports → Goal: Q1 Revenue (70% progress)
# ### Facts & Patterns
# - (none yet)

# 5. Task is completed
await task_service.complete(task.id)

# 6. Bridge records completion (via trigger)
# Adds fact: "Completed on time" or "Completed 30min under estimate"

# 7. Next time AI queries, facts are included
```

### Example 2: AI Interaction

```
User: "What should I focus on today?"

Legacy Agent (internal):
1. Call get_day_predictions()
2. Call get_intelligence_scores()
3. Analyze patterns and blockers

Legacy Agent (response):
"Good morning! Based on your patterns:

**High Priority (WI at 76, needs boost)**
- Review quarterly report (85% completion likelihood)
- Best time: 8-10am before your meeting cluster

**Health (PI at 82, strong)**
- Morning routine on track (92% likelihood)
- Consider shorter workout due to 3 meetings

**Family (LI at 71, needs attention)**
- Noah has math homework due Friday
- Block 30min tonight to help

Should I schedule the quarterly report for 8am?"
```

### Example 3: Recording Context

```
User: "Noah finished his math early today"

Legacy Agent (internal):
1. Find Noah entity
2. Find math task
3. Mark task complete
4. Add fact about early completion

await add_entity_context(
    entity_type="child", entity_id="Noah",
    fact="Completed math homework 15min early on Tuesday"
)

Legacy Agent (response):
"Nice work Noah! I've noted that he finished math early today.

I'm starting to see a pattern - he's been completing math faster
lately. His average completion time improved from 45min to 30min
over the last 2 weeks.

Should I suggest more challenging problems to keep building momentum?"
```

---

## Summary

The Bridge module connects LifeOS to the Learning Machine with:

- **~1,970 lines of code** (vs 20,000+ for 43 tables)
- **10 API endpoints** at `/api/bridge/*`
- **7 AI tools** for the Legacy agent
- **8 auto-registration triggers** for LifeOS services
- **3 intelligence scores** (PI/LI/WI)
- **Pattern detection** for habits and tasks
- **Day predictions** with completion likelihoods
- **Family support** for multi-user households

All without creating new database tables - just connecting what already exists.

---

## Flutter Frontend Integration

The Bridge module includes a complete Flutter frontend at `frontend/lib/lifeos/bridge/`.

### File Structure

```
frontend/lib/lifeos/bridge/
├── models/
│   └── bridge_models.dart        # All data models (700+ lines)
├── services/
│   └── bridge_api.dart           # HTTP client (250+ lines)
├── providers/
│   └── bridge_provider.dart      # Riverpod state management (400+ lines)
├── widgets/
│   ├── intelligence_scores_card.dart  # PI/LI/WI display widget
│   ├── day_predictions_card.dart      # Predictions display widget
│   └── widgets.dart                   # Barrel export
└── screens/
    └── bridge_dashboard_screen.dart   # Main dashboard (500+ lines)
```

### Models (`bridge_models.dart`)

```dart
// Intelligence Scores
class IntelligenceScores {
  final PillarScore pi;  // Performance Intelligence
  final PillarScore li;  // Legacy Intelligence
  final PillarScore wi;  // Wealth Intelligence
  final double unifiedScore;
  final List<ScoreSuggestion> suggestions;
}

class PillarScore {
  final double score;      // 0-100
  final String trend;      // 'up', 'down', 'stable'
  final Map<String, double> components;
}

// Day Predictions
class DayPrediction {
  final DateTime date;
  final double predictedEnergy;
  final double meetingHours;
  final List<TaskPrediction> tasks;
  final List<HabitPrediction> habits;
  final List<String> blockers;
  final List<String> suggestions;
}

class TaskPrediction {
  final String taskName;
  final double completionLikelihood;  // 0.0-1.0
  final String? optimalTime;
  final List<String> blockers;
}

// User Patterns
class UserPatterns {
  final List<HabitPattern> habitPatterns;
  final TaskPatterns? taskPatterns;
}

class HabitPattern {
  final String habitName;
  final double completionRate;
  final Map<String, double> bestDays;
  final List<String> correlations;
}

// Family
class FamilyOverview {
  final List<FamilyMember> members;
  final Map<String, List<FamilyTask>> todayTasks;
  final List<FamilyGoal> sharedGoals;
}

// Entity Context
class EntityContext {
  final String entityId;
  final String type;
  final String name;
  final Pillar? pillar;
  final Map<String, dynamic> lifeosData;
  final List<EntityRelation> relations;
  final List<EntityFact> facts;
}
```

### API Service (`bridge_api.dart`)

```dart
class BridgeApi {
  final String baseUrl;

  // Context
  Future<EntityContext> getEntityContext({
    required String entityType,
    required String entityId,
    String userId = 'zeus',
  });

  // Intelligence Scores
  Future<IntelligenceScores> getIntelligenceScores({String userId = 'zeus'});

  // Predictions
  Future<DayPrediction> getDayPredictions({DateTime? date, String userId = 'zeus'});
  Future<DayPrediction> getTodayPredictions({String userId = 'zeus'});

  // Patterns
  Future<UserPatterns> getUserPatterns({String userId = 'zeus'});

  // Family
  Future<FamilyOverview> getFamilyOverview({String userId = 'zeus'});

  // Entity Operations
  Future<LinkResponse> linkEntities({
    required String fromType,
    required String fromId,
    required String toType,
    required String toId,
    required String relationType,
  });

  Future<FactResponse> addEntityFact({
    required String entityType,
    required String entityId,
    required String fact,
    String source = 'user',
    double confidence = 1.0,
  });

  // Admin
  Future<BackfillResponse> backfillData({String userId = 'zeus'});
  Future<BridgeStats> getStats({String userId = 'zeus'});
}
```

### Riverpod Providers (`bridge_provider.dart`)

```dart
// API client
final bridgeApiProvider = Provider<BridgeApi>((ref) {
  final baseUrl = ref.watch(apiBaseUrlProvider);
  return BridgeApi(baseUrl: baseUrl);
});

// Intelligence scores
final intelligenceScoresProvider = StateNotifierProvider.family<
    IntelligenceScoresNotifier, IntelligenceScoresState, String>(...);

final currentUserScoresProvider = Provider<IntelligenceScoresState>((ref) {
  final userId = ref.watch(currentUserIdProvider);
  return ref.watch(intelligenceScoresProvider(userId));
});

// Day predictions
final todayPredictionsProvider = StateNotifierProvider.family<
    DayPredictionNotifier, DayPredictionState, String>(...);

// User patterns
final userPatternsProvider = StateNotifierProvider.family<
    UserPatternsNotifier, UserPatternsState, String>(...);

// Family overview
final familyOverviewProvider = StateNotifierProvider.family<
    FamilyOverviewNotifier, FamilyOverviewState, String>(...);

// Entity context
final entityContextProvider = StateNotifierProvider.family<
    EntityContextNotifier, EntityContextState, EntityContextKey>(...);

// Combined dashboard state
final bridgeDashboardProvider = Provider<BridgeDashboardState>((ref) {
  return BridgeDashboardState(
    scores: ref.watch(currentUserScoresProvider),
    predictions: ref.watch(currentUserPredictionsProvider),
    patterns: ref.watch(currentUserPatternsProvider),
    family: ref.watch(currentFamilyOverviewProvider),
  );
});
```

### Using in Flutter App

```dart
// In a widget
class MyDashboard extends ConsumerWidget {
  @override
  Widget build(BuildContext context, WidgetRef ref) {
    // Watch intelligence scores
    final scoresState = ref.watch(currentUserScoresProvider);

    if (scoresState.isLoading) {
      return CircularProgressIndicator();
    }

    final scores = scoresState.scores!;
    return Column(
      children: [
        Text('PI: ${scores.pi.score}'),
        Text('LI: ${scores.li.score}'),
        Text('WI: ${scores.wi.score}'),
        Text('Unified: ${scores.unifiedScore}'),
      ],
    );
  }
}

// Use pre-built widgets
IntelligenceScoresCard(userId: 'zeus', compact: false)
DayPredictionsCard(userId: 'zeus', showBlockers: true)

// Full dashboard screen
BridgeDashboardScreen(userId: 'zeus')
```

---

## Service Integration (Automatic Triggers)

Bridge triggers are now automatically called when tasks, habits, and goals are created or updated.

### Tasks Service Integration

```python
# backend/lifeos/tasks/service.py

# Import at top
try:
    from backend.lifeos.bridge.triggers import BridgeTriggers
    BRIDGE_ENABLED = True
except ImportError:
    BRIDGE_ENABLED = False

# After task creation
async def create(self, data: TaskCreate, actor: str = "zeus") -> Task:
    # ... create task ...
    task = self._row_to_task(row)

    # Register in Bridge knowledge graph
    if BRIDGE_ENABLED:
        try:
            await BridgeTriggers.on_task_created(task, data.user_id)
        except Exception as e:
            logger.warning(f"Bridge trigger failed: {e}")

    return task

# After task completion
async def complete(self, task_id: UUID, data: TaskComplete, user_id: str) -> Task:
    # ... complete task ...
    task = self._row_to_task(row)

    # Record completion in Bridge knowledge graph
    if BRIDGE_ENABLED:
        try:
            await BridgeTriggers.on_task_completed(task, user_id, actual_minutes)
        except Exception as e:
            logger.warning(f"Bridge trigger failed: {e}")

    return task
```

### Habits Service Integration

```python
# backend/lifeos/habits/service.py

# Import at top
try:
    from backend.lifeos.bridge.triggers import BridgeTriggers
    BRIDGE_ENABLED = True
except ImportError:
    BRIDGE_ENABLED = False

# After habit creation
async def create_habit(self, user_id: str, data: HabitCreate) -> Habit:
    # ... create habit ...
    habit = self._row_to_habit(row)

    # Register in Bridge knowledge graph
    if BRIDGE_ENABLED and habit:
        try:
            await BridgeTriggers.on_habit_created(habit, user_id)
        except Exception as e:
            logger.warning(f"Bridge trigger failed: {e}")

    return habit

# After habit completion/skip
async def complete_habit(self, habit_id: UUID, data: HabitCompletionCreate) -> HabitCompletion:
    # ... complete habit ...

    # Record in Bridge knowledge graph
    if BRIDGE_ENABLED and status == "completed":
        try:
            habit = await self.get_habit(habit_id)
            await BridgeTriggers.on_habit_completed(habit, user_id, mood, energy)
        except Exception as e:
            logger.warning(f"Bridge trigger failed: {e}")
    elif BRIDGE_ENABLED and status == "skipped":
        try:
            habit = await self.get_habit(habit_id)
            await BridgeTriggers.on_habit_skipped(habit, user_id, reason, energy)
        except Exception as e:
            logger.warning(f"Bridge trigger failed: {e}")

    return completion
```

### Goals Service Integration

```python
# backend/lifeos/goals/service.py

# Import at top
try:
    from backend.lifeos.bridge.triggers import BridgeTriggers
    BRIDGE_ENABLED = True
except ImportError:
    BRIDGE_ENABLED = False

# After goal creation
async def create_goal(self, user_id: str, data: GoalCreate) -> Goal:
    # ... create goal ...
    goal = self._row_to_goal(row)

    # Register in Bridge knowledge graph
    if BRIDGE_ENABLED and goal:
        try:
            await BridgeTriggers.on_goal_created(goal, user_id)
        except Exception as e:
            logger.warning(f"Bridge trigger failed: {e}")

    return goal

# After progress update
async def update_progress(self, goal_id: UUID, data: GoalProgress) -> GoalProgressResponse:
    # ... update progress ...

    # Record progress in Bridge knowledge graph
    if BRIDGE_ENABLED:
        try:
            goal = await self.get_goal(goal_id)
            await BridgeTriggers.on_goal_progress_updated(goal, user_id, old_progress, new_progress)
        except Exception as e:
            logger.warning(f"Bridge trigger failed: {e}")

    return response
```

### What Gets Recorded

| Service | Event | Entity Registered | Facts Added |
|---------|-------|-------------------|-------------|
| Tasks | create | task entity | metadata, pillar, project |
| Tasks | complete | - | completion time, efficiency |
| Habits | create | habit entity | frequency, pillar, steps |
| Habits | complete | - | mood, energy, streak |
| Habits | skip | - | reason, energy |
| Goals | create | goal entity | target, pillar, parent |
| Goals | progress | - | progress percentage, milestone |

### Verifying Integration

```bash
# Create a task and verify entity was registered
curl -X POST "http://localhost:8081/api/tasks" \
  -H "Content-Type: application/json" \
  -d '{"title":"Test task","pillar":"HP","user_id":"zeus"}'

# Check Bridge stats
curl "http://localhost:8081/api/bridge/stats?user_id=zeus"
# Should show: {"total_entities": N+1, ...}

# Get context for the new task
curl "http://localhost:8081/api/bridge/context/task/Test%20task?user_id=zeus"
```

---

## Updated Summary

The Bridge module now includes:

- **~2,500 lines backend code** (service, models, api, triggers)
- **~1,900 lines Flutter code** (models, service, providers, widgets, screens)
- **12 API endpoints** at `/api/bridge/*`
- **7 AI tools** for the Legacy agent
- **3 service integrations** (Tasks, Habits, Goals)
- **Auto-registration** for all new items
- **Full Flutter dashboard** with intelligence scores, predictions, patterns, family


---

## Related

- [[00-24 SO (Sovereign Operations)/02 Infrastructure (VPS, Servers & Networking) (The Body)/02 Infrastructure|02 Infrastructure]]
