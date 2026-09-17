# OneIntel

> Unified item system for OneMind OS - tasks, projects, docs, habits in one place.

## Overview

OneIntel is the all-in-one data layer for OneMind OS. It provides a single, unified `items` table that stores:

- **Projects** - Containers for tasks and docs
- **Tasks** - Actionable work items with hierarchy (subtasks)
- **Docs** - Markdown documents
- **Habits** - Recurring items with streak tracking

## Quick Start

### 1. Run Standalone API (for testing)

```bash
cd oneintel
uvicorn api:app --reload --port 8081
```

### 2. Test Endpoints

```bash
# Create a project
curl -X POST http://localhost:8081/items \
  -H "Content-Type: application/json" \
  -d '{"type": "project", "title": "My First Project"}'

# List all items
curl http://localhost:8081/items

# Get today's focus
curl http://localhost:8081/items/today

# Create a task
curl -X POST http://localhost:8081/items \
  -H "Content-Type: application/json" \
  -d '{"type": "task", "title": "Review PR", "due_date": "2025-02-05", "priority": 2}'

# Complete a task
curl -X POST http://localhost:8081/items/{item_id}/complete

# Search
curl "http://localhost:8081/items/search?q=review"
```

### 3. Test the Tool (standalone)

```python
from oneintel.tool import oneintel

# List all projects
result = oneintel(action="list", type="project")
print(result)

# Create a task
result = oneintel(action="create", type="task", title="Test task")
print(result)

# Get today's view
result = oneintel(action="today")
print(result)
```

## File Structure

```
oneintel/
├── README.md          # This file
├── PLAN.md            # Full implementation plan
├── __init__.py        # Package init
├── models.py          # Pydantic models (Item, ItemCreate, etc.)
├── schema.sql         # SQLite schema
├── storage.py         # Storage layer (SQLite)
├── api.py             # Standalone FastAPI for testing
├── tool.py            # Agno tool (standalone, not registered yet)
│
├── backend/           # Reusable backend utilities
│   ├── __init__.py
│   ├── scoring.py     # Habit scoring (exponential decay, A-F grades)
│   └── recipe_scraper.py  # Recipe extraction from URLs
│
└── frontend/          # Flutter widgets (standalone)
    └── widgets/
        ├── widgets.dart         # Barrel export
        ├── streak_indicator.dart  # Streak badges with tier levels
        ├── progress_ring.dart     # Circular progress rings
        ├── mini_heatmap.dart      # GitHub-style completion grids
        └── trend_chart.dart       # Line charts (requires fl_chart)
```

## Data Model

```
┌─────────────────────────────────────────────────────────────────┐
│                         items TABLE                              │
├─────────────────────────────────────────────────────────────────┤
│  id          TEXT PRIMARY KEY                                    │
│  type        TEXT (project | task | doc | habit)                │
│  title       TEXT NOT NULL                                       │
│  content     TEXT (markdown body)                                │
│  parent_id   TEXT (foreign key for hierarchy)                   │
│  status      TEXT (todo | in_progress | done | archived)        │
│  due_date    DATE                                                │
│  priority    INT (1=urgent, 2=high, 3=normal, 4=low)            │
│  tags        TEXT (JSON array)                                   │
│  streak      INT (for habits)                                    │
│  last_done   TIMESTAMP (for habits)                              │
│  created_at  TIMESTAMP                                           │
│  updated_at  TIMESTAMP                                           │
│  user_id     TEXT                                                │
└─────────────────────────────────────────────────────────────────┘
```

## Tool Actions

| Action | Description | Required Params |
|--------|-------------|-----------------|
| `list` | List items with filters | (optional) type, status, parent_id, tags |
| `get` | Get single item | item_id |
| `create` | Create new item | type, title |
| `update` | Update item fields | item_id, (fields to update) |
| `complete` | Mark task done / complete habit | item_id |
| `delete` | Remove item (cascades) | item_id |
| `search` | Full-text search | query |
| `tree` | Get item with children | item_id |
| `today` | Today's tasks + habits | - |
| `stats` | Get statistics | - |

## Design Principles

1. **One table** - Everything in `items`
2. **One tool** - `oneintel()` for all operations
3. **12 fields max** - Start minimal
4. **Agent-first** - Built for Agno tools, not REST services
5. **Hierarchy via parent_id** - Projects → Tasks → Subtasks

## Backend Utilities

### Habit Scoring (Loop Habits-inspired)

```python
from oneintel.backend import calculate_habit_score, calculate_grade, get_grade_info

# Exponential decay scoring (recent completions count more)
completion_dates = {date(2025, 1, 28), date(2025, 1, 29), date(2025, 1, 30)}
score = calculate_habit_score(completion_dates)
print(score.score)  # 0-100 score
print(score.trend)  # 'improving', 'stable', 'declining'
print(score.streak) # Current streak

# A-F grading
grade = calculate_grade(completions=9, expected=10)  # Grade.A (90%)

# Full grade info
grade_info = get_grade_info(completion_dates, frequency="daily", period_days=30)
print(grade_info.grade, grade_info.percentage)
```

### Recipe Scraper

```python
from oneintel.backend import scrape_recipe

# Extracts JSON-LD Schema.org or microdata from recipe sites
recipe = await scrape_recipe("https://example.com/recipe/chocolate-cake")
print(recipe.name)
print(recipe.ingredients)  # List of Ingredient objects
print(recipe.instructions) # List of Instruction objects
print(recipe.nutrition)    # Nutrition info if available
```

## Frontend Widgets (Flutter)

Standalone widgets for habit/task tracking UI. Copy to your Flutter project.

### Dependencies

Add to `pubspec.yaml`:
```yaml
dependencies:
  fl_chart: ^0.68.0  # Required for TrendChart
```

### Usage

```dart
import 'package:oneintel/frontend/widgets/widgets.dart';

// Streak indicator with tier colors
StreakIndicator(streak: 15)  // Shows flame icon with "15"
StreakBadge(streak: 45)      // "45 days - On Fire" badge

// Progress ring
ProgressRing(progress: 0.75, color: Colors.green)
LargeProgressRing(progress: 0.8, label: "Tasks")

// Mini heatmaps
MiniHeatmap(completionDates: {...}, days: 7)  // 7-day row
WeeklyHeatmap(completionDates: {...}, frequencyDays: [1,2,3,4,5])
MonthHeatmap(completionDates: {...})  // Full month grid

// Trend chart (requires fl_chart)
TrendChart(
  dataPoints: [DataPoint(date: DateTime.now(), value: 85.0), ...],
  periodDays: 30,
)
PeriodSelector(selectedDays: 30, onPeriodChanged: (days) => ...)
```

### Widget Features

| Widget | Description |
|--------|-------------|
| `StreakIndicator` | Flame icon with tier colors (Spark/Warm/Fire/Legendary) |
| `StreakBadge` | Badge with day count and tier label |
| `ScoreIndicator` | 0-100 score circle with trend arrow |
| `ProgressRing` | Circular progress with glow on 100% |
| `MiniHeatmap` | 7-day completion squares |
| `WeeklyHeatmap` | M-T-W-T-F-S-S grid with frequency awareness |
| `MonthHeatmap` | Full month calendar grid |
| `TrendChart` | Line chart with tooltip (uses fl_chart) |
| `CompletionBars` | Vertical bar chart for rates |

## Integration (Later)

When ready to integrate with OneMind OS:

1. Move `tool.py` to `backend/agno/tools/oneintel.py`
2. Register in `backend/agno/tools/__init__.py`
3. Add to agent presets
4. Create Flutter frontend screens

## Dependencies

### Python (Backend)

```bash
pip install fastapi uvicorn pydantic httpx
```

### Flutter (Frontend)

```yaml
# pubspec.yaml
dependencies:
  fl_chart: ^0.68.0
```

## Database

By default uses SQLite (`oneintel.db` in current directory).

For production, use PostgreSQL schema in `schema_pg.sql`.
