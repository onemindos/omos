# Temporal Core

> Permanent time infrastructure for OneMind OS

---

## Table of Contents

1. [Overview](#overview)
2. [Quick Start](#quick-start)
3. [Installation](#installation)
4. [Architecture](#architecture)
5. [Components](#components)
6. [Usage Guide](#usage-guide)
7. [API Reference](#api-reference)
8. [Database Setup](#database-setup)
9. [Integration Examples](#integration-examples)
10. [Troubleshooting](#troubleshooting)

---

## Overview

Temporal Core is **foundational infrastructure**, not an agent. It provides time awareness to every component in the system - agents, tools, APIs, and services all have access to rich temporal context.

**Key Principle**: Time is "baked in" to the OS itself. Every agent automatically knows:
- What time it is (with timezone awareness)
- What's on the calendar
- What deadlines are pressing
- When productivity typically peaks
- When focus windows occur

### Why Not a "Chronos Agent"?

An agent is ephemeral - it exists during a conversation. Temporal Core is **permanent infrastructure**:

| Aspect | Agent Approach | Infrastructure Approach |
|--------|---------------|------------------------|
| Lifespan | Per-conversation | Always running |
| Access | Via @mention | Injected everywhere |
| Storage | Queries existing data | Owns time-series data |
| Pattern Detection | On-demand | Continuous background |
| Scheduling | Manual | Automatic triggers |

---

## Quick Start

### 1. Basic Usage (No Database Required)

```python
from backend.consciousness.temporal import (
    get_temporal_service,
    TEMPORAL_TOOLS,
)

# Get temporal context immediately
service = get_temporal_service()
context = await service.get_context(user_id="zeus")

print(f"Time: {context.now}")
print(f"Time of day: {context.time_of_day.value}")
print(f"Is focus time: {context.is_focus_time}")
print(context.to_prompt_context())
```

### 2. Add Time Awareness to an Agent

```python
from agno import Agent
from backend.consciousness.temporal import TEMPORAL_TOOLS

agent = Agent(
    name="time-aware-agent",
    tools=TEMPORAL_TOOLS,  # Add all 11 temporal tools
    system_prompt="You are a time-aware assistant. Use get_time_context for awareness.",
)
```

### 3. Mount the API

```python
from fastapi import FastAPI
from backend.consciousness.temporal import temporal_router

app = FastAPI()
app.include_router(temporal_router, prefix="/api/temporal")

# Now available:
# GET /api/temporal/context
# GET /api/temporal/now
# GET /api/temporal/calendar
# etc.
```

---

## Installation

### Prerequisites

```bash
# Required
Python 3.11+
FastAPI
Pydantic v2
python-dateutil
zoneinfo (stdlib)

# Optional (for full functionality)
TimescaleDB 2.0+  # Time-series storage
NATS JetStream    # Event streaming
```

### Install Dependencies

```bash
# In your project
pip install python-dateutil

# For TimescaleDB (optional)
pip install asyncpg psycopg2-binary
```

### Verify Installation

```python
# Test import
from backend.consciousness.temporal import (
    TemporalService,
    TemporalContext,
    TEMPORAL_TOOLS,
    temporal_router,
)
print("✓ Temporal Core installed successfully")
```

---

## Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                    Temporal Core                            │
├─────────────────────────────────────────────────────────────┤
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────────────┐ │
│  │  Service    │  │  Scheduler  │  │  Pattern Engine     │ │
│  │             │  │             │  │                     │ │
│  └─────────────┘  └─────────────┘  └─────────────────────┘ │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────────────┐ │
│  │  Calendar   │  │  Deadlines  │  │  Context Generator  │ │
│  │  Providers  │  │             │  │                     │ │
│  └─────────────┘  └─────────────┘  └─────────────────────┘ │
├─────────────────────────────────────────────────────────────┤
│  TimescaleDB (event_ledger, continuous aggregates)         │
└─────────────────────────────────────────────────────────────┘
                          │
                          ▼
          Injected into every Agent/Tool/Service
```

### File Structure

```
backend/consciousness/temporal/
├── __init__.py          # Package exports (start here)
├── models.py            # Data models (enums, Pydantic classes)
├── schema.sql           # TimescaleDB schema (run this first)
├── service.py           # Core TemporalService
├── scheduler.py         # Event scheduling
├── patterns.py          # Pattern detection engine
├── calendar.py          # Calendar provider interfaces
├── context.py           # Context injection utilities
├── tools.py             # 11 Agno-native agent tools
├── api.py               # FastAPI REST endpoints
└── TEMPORAL.md          # This documentation
```

---

## Components

### 1. TemporalService (`service.py`)

The core service coordinating all temporal functionality.

**Key Methods:**

| Method | Description | Returns |
|--------|-------------|---------|
| `get_context(user_id, timezone)` | Complete temporal awareness | `TemporalContext` |
| `schedule(event)` | Schedule future action | `ScheduledEvent` |
| `get_deadlines(user_id)` | Get tracked deadlines | `List[Deadline]` |
| `add_deadline(title, due, ...)` | Track a new deadline | `Deadline` |
| `get_productivity_by_hour()` | Hourly productivity scores | `Dict[int, float]` |
| `get_patterns(user_id)` | Detected temporal patterns | `Dict[str, Pattern]` |

**Example:**

```python
from backend.consciousness.temporal import get_temporal_service

service = get_temporal_service()

# Get complete context
context = await service.get_context(user_id="zeus")

# Access context fields
print(f"Now: {context.now}")
print(f"Time of day: {context.time_of_day}")  # morning, afternoon, etc.
print(f"Day type: {context.day_type}")        # workday, weekend
print(f"Next event: {context.next_event}")
print(f"Deadline pressure: {context.deadline_pressure:.0%}")
print(f"Productivity: {context.productivity_score:.0%}")
print(f"Is focus time: {context.is_focus_time}")

# Generate prompt context string
print(context.to_prompt_context())
```

### 2. TemporalContext (`models.py`)

Complete time awareness - this is what gets injected into agents.

```python
class TemporalContext:
    # Current moment
    now: datetime
    timezone: str
    time_of_day: TimeOfDay      # early_morning, morning, afternoon, evening, night, late_night
    day_type: DayType           # workday, weekend, holiday, vacation

    # Calendar awareness
    current_event: Optional[CalendarEvent]
    next_event: Optional[CalendarEvent]
    minutes_until_next: Optional[int]
    today_events: List[CalendarEvent]
    events_count_today: int

    # Deadlines
    urgent_deadlines: List[Deadline]
    deadline_pressure: float    # 0.0-1.0 (higher = more pressure)
    overdue_count: int
    due_today_count: int

    # Patterns (from TimescaleDB analysis)
    productivity_score: float   # 0.0-1.0 for current hour
    typical_load: str           # "lighter", "normal", "heavier"
    energy_prediction: float    # 0.0-1.0

    # Focus windows
    is_focus_time: bool
    focus_windows_today: List[TimeRange]
    next_focus_window: Optional[TimeRange]

    # Activity
    events_last_hour: int
    events_last_24h: int
```

### 3. TemporalScheduler (`scheduler.py`)

Schedule future events and actions.

**Instructions:**

```python
from backend.consciousness.temporal import (
    TemporalScheduler,
    ScheduledEventCreate,
    ScheduleEventType,
    create_daily_schedule,
    create_weekly_schedule,
)
from datetime import datetime, timedelta

# Initialize scheduler
scheduler = TemporalScheduler(
    default_timezone="America/Los_Angeles",
    check_interval=60,  # Check every 60 seconds
)

# 1. Schedule a one-time reminder
await scheduler.schedule(ScheduledEventCreate(
    type=ScheduleEventType.REMINDER,
    scheduled_at=datetime.now() + timedelta(hours=1),
    title="Stand up and stretch",
    description="Take a 5 minute break",
))

# 2. Schedule a workflow trigger
await scheduler.schedule(ScheduledEventCreate(
    type=ScheduleEventType.TRIGGER,
    scheduled_at=datetime.now() + timedelta(days=1),
    title="Morning briefing",
    action_type="workflow",
    action_target="morning-briefing-workflow",
    action_params={"include_weather": True},
))

# 3. Create daily recurring schedule (helper function)
daily_reminder = create_daily_schedule(
    title="Daily standup reminder",
    hour=9,
    minute=30,
    timezone="America/Los_Angeles",
)
await scheduler.schedule(daily_reminder)

# 4. Create weekly recurring schedule
weekly_review = create_weekly_schedule(
    title="Weekly review",
    day_of_week=4,  # Friday (0=Monday, 6=Sunday)
    hour=16,
    minute=0,
)
await scheduler.schedule(weekly_review)

# 5. Register action handlers
async def handle_workflow_trigger(event: ScheduledEvent):
    # Called when a workflow trigger fires
    workflow_id = event.action_target
    params = event.action_params
    # ... execute workflow

async def handle_notification(event: ScheduledEvent):
    # Called for notification triggers
    # ... send notification

scheduler.register_handler("workflow", handle_workflow_trigger)
scheduler.register_handler("notification", handle_notification)

# 6. Start the scheduler (run in background)
await scheduler.start()

# 7. Query scheduled events
pending = await scheduler.get_pending(
    until=datetime.now() + timedelta(days=7),
    user_id="zeus",
)

# 8. Cancel a schedule
await scheduler.cancel(event_id)

# 9. Snooze an event
await scheduler.snooze(event_id, timedelta(hours=1))

# 10. Stop scheduler on shutdown
await scheduler.stop()
```

### 4. Pattern Detection (`patterns.py`)

Analyzes historical data to detect patterns.

**Available Pattern Types:**

| Pattern | What It Detects |
|---------|-----------------|
| `PRODUCTIVITY` | Task completion rates by hour/day |
| `ENERGY` | Energy level predictions |
| `FOCUS` | Optimal deep work windows |
| `STRESS` | Stress level patterns |
| `ACTIVITY` | General activity patterns |

**Instructions:**

```python
from backend.consciousness.temporal import TemporalPatternEngine

engine = TemporalPatternEngine()

# Detect all patterns from aggregated metrics
# (metrics come from TimescaleDB continuous aggregates)
patterns = await engine.detect_patterns(
    metrics=aggregated_metrics,  # List[AggregatedMetric]
    user_id="zeus",
)

# Access specific pattern
for pattern in patterns:
    print(f"Pattern: {pattern.type.value}")
    print(f"Confidence: {pattern.confidence:.0%}")
    print(f"Peak hours: {pattern.peak_hours}")
    print(f"Peak days: {pattern.peak_days}")
    print(f"Insights: {pattern.insights}")

# Get cached pattern
productivity = engine.get_cached_pattern(PatternType.PRODUCTIVITY, "zeus")
if productivity:
    print(f"Best hours: {productivity.optimal_deep_work_hours}")
    print(f"Meeting hours: {productivity.optimal_meeting_hours}")
```

### 5. Calendar Integration (`calendar.py`)

Abstract interface for calendar providers.

**Instructions:**

```python
from backend.consciousness.temporal import (
    CompositeCalendarProvider,
    InMemoryCalendarProvider,
    CalendarEvent,
)
from datetime import datetime, timedelta

# Option 1: In-memory provider (for testing/development)
calendar = InMemoryCalendarProvider(default_timezone="America/Los_Angeles")

# Add test events
calendar.add_event(CalendarEvent(
    id="event-1",
    title="Team Standup",
    start=datetime.now() + timedelta(hours=1),
    end=datetime.now() + timedelta(hours=1, minutes=30),
    source="test",
))

# Option 2: Composite provider (combine multiple sources)
composite = CompositeCalendarProvider()
composite.add_provider(InMemoryCalendarProvider())
# composite.add_provider(GoogleCalendarProvider())  # When implemented
# composite.add_provider(AppleCalendarProvider())   # When implemented

# Get events in range
events = await calendar.get_events(
    start=datetime.now(),
    end=datetime.now() + timedelta(days=7),
    user_id="zeus",
)

# Get current event (if in a meeting now)
current = await calendar.get_current_event(user_id="zeus")
if current:
    print(f"Currently in: {current.title}")

# Get next upcoming event
next_event = await calendar.get_next_event(user_id="zeus")
if next_event:
    minutes = next_event.minutes_until()
    print(f"Next: {next_event.title} in {minutes} minutes")

# Find free time slots
free_slots = await calendar.get_free_slots(
    date=datetime.now(),
    duration_minutes=30,
    work_hours=(9, 17),
    user_id="zeus",
)
for start, end in free_slots:
    print(f"Free: {start.strftime('%H:%M')} - {end.strftime('%H:%M')}")
```

### 6. Context Injection (`context.py`)

Utilities for injecting temporal context into agents and prompts.

**Instructions:**

```python
from backend.consciousness.temporal import (
    inject_temporal_context,
    inject_minimal_time,
    with_temporal_context,
    get_context_generator,
    get_context_enricher,
    TEMPORAL_SYSTEM_PROMPT_TEMPLATE,
)

# Method 1: Inject full context into system prompt
service = get_temporal_service()
context = await service.get_context(user_id="zeus")

enhanced_prompt = inject_temporal_context(
    system_prompt="You are a helpful assistant.",
    context=context,
)
# Result includes full time context section

# Method 2: Minimal injection (no async, no database)
# Good for quick context without external lookups
enhanced_prompt = inject_minimal_time(
    system_prompt="You are a helpful assistant.",
    timezone="America/Los_Angeles",
)
# Result includes basic time info only

# Method 3: Decorator for automatic context injection
@with_temporal_context(user_id="zeus")
async def my_agent_function(prompt: str, temporal_context: TemporalContext):
    # temporal_context is automatically provided
    print(f"Current time: {temporal_context.now}")
    print(f"Is focus time: {temporal_context.is_focus_time}")

# Method 4: Context generator for custom use
generator = get_context_generator()

# Get basic context (no async)
basic = generator.get_basic_context()
print(basic)  # {"now": "...", "time_of_day": "morning", ...}

# Get brief string
brief = generator.to_brief_string()
print(brief)  # "Morning (workday)"

# Method 5: Enrich agent parameters
enricher = get_context_enricher()
params = {"message": "Hello"}
enriched = await enricher.enrich_agent_params(params, user_id="zeus")
# params now includes temporal_context
```

### 7. Agent Tools (`tools.py`)

11 Agno-native tools for agent use.

**Available Tools:**

| Tool | Description | Key Parameters |
|------|-------------|----------------|
| `get_time_context` | Complete temporal awareness | `timezone`, `user_id` |
| `get_current_time` | Simple time lookup | `timezone`, `format` |
| `get_calendar_events` | Events in time range | `start`, `end`, `days_ahead` |
| `find_free_time` | Available time slots | `date`, `duration_minutes` |
| `get_deadlines` | Deadlines by urgency | `urgency`, `days_ahead` |
| `add_deadline` | Track new deadline | `title`, `due`, `priority` |
| `schedule_reminder` | Future reminder | `title`, `when` |
| `get_scheduled` | Pending schedules | `days_ahead`, `event_type` |
| `cancel_scheduled` | Cancel schedule | `event_id` |
| `get_productivity_patterns` | Hourly productivity | `user_id` |
| `suggest_optimal_time` | Best time for task | `task_type`, `duration` |

**Instructions:**

```python
from agno import Agent
from backend.consciousness.temporal import TEMPORAL_TOOLS

# Add all temporal tools to an agent
agent = Agent(
    name="time-aware-agent",
    model="claude-3-5-sonnet",
    tools=TEMPORAL_TOOLS,
    system_prompt="""You are a time-aware assistant.

    ALWAYS use get_time_context at the start of conversations to understand:
    - Current time and day type
    - Upcoming calendar events
    - Deadline pressure
    - Whether it's a good time for focused work

    Use suggest_optimal_time when the user needs to schedule something.
    Use find_free_time to check calendar availability.
    Use get_deadlines to understand time pressure.
    """,
)

# Example: Add only specific tools
from backend.consciousness.temporal.tools import (
    get_time_context,
    get_deadlines,
    schedule_reminder,
)

agent = Agent(
    name="minimal-time-agent",
    tools=[get_time_context, get_deadlines, schedule_reminder],
)
```

### 8. API Routes (`api.py`)

FastAPI router for REST access.

**Instructions:**

```python
from fastapi import FastAPI
from backend.consciousness.temporal import temporal_router

app = FastAPI()

# Mount the temporal router
app.include_router(temporal_router, prefix="/api/temporal")

# Now you have these endpoints:
# GET  /api/temporal/context              - Complete temporal context
# GET  /api/temporal/now                  - Current time
# GET  /api/temporal/calendar             - Calendar events
# GET  /api/temporal/calendar/current     - Current event
# GET  /api/temporal/calendar/next        - Next event
# GET  /api/temporal/deadlines            - Deadlines by urgency
# POST /api/temporal/deadlines            - Create deadline
# GET  /api/temporal/schedules            - Pending schedules
# POST /api/temporal/schedules            - Create schedule
# DELETE /api/temporal/schedules/{id}     - Cancel schedule
# GET  /api/temporal/patterns             - All patterns
# GET  /api/temporal/patterns/productivity - Hourly productivity
```

**Example API Calls:**

```bash
# Get temporal context
curl http://localhost:8080/api/temporal/context?timezone=America/Los_Angeles

# Get current time
curl http://localhost:8080/api/temporal/now?timezone=UTC

# Get deadlines
curl "http://localhost:8080/api/temporal/deadlines?days=7&urgency=urgent"

# Create a deadline
curl -X POST http://localhost:8080/api/temporal/deadlines \
  -H "Content-Type: application/json" \
  -d '{"title": "Submit report", "due": "2025-02-05T17:00:00", "priority": 2}'

# Schedule a reminder
curl -X POST http://localhost:8080/api/temporal/schedules \
  -H "Content-Type: application/json" \
  -d '{
    "type": "reminder",
    "scheduled_at": "2025-02-04T09:00:00",
    "title": "Morning standup",
    "timezone": "America/Los_Angeles"
  }'

# Get productivity patterns
curl http://localhost:8080/api/temporal/patterns/productivity
```

---

## Usage Guide

### Scenario 1: Make Any Agent Time-Aware

```python
from agno import Agent
from backend.consciousness.temporal import (
    TEMPORAL_TOOLS,
    get_temporal_service,
    inject_temporal_context,
)

# Option A: Add tools (agent can query time)
agent = Agent(
    name="my-agent",
    tools=TEMPORAL_TOOLS,
)

# Option B: Inject context into prompt (always aware)
service = get_temporal_service()
context = await service.get_context()

agent = Agent(
    name="my-agent",
    system_prompt=inject_temporal_context(
        "You are a helpful assistant.",
        context,
    ),
)

# Option C: Both (recommended)
agent = Agent(
    name="my-agent",
    tools=TEMPORAL_TOOLS,
    system_prompt=inject_temporal_context(
        "You are a time-aware assistant. Use temporal tools for detailed queries.",
        context,
    ),
)
```

### Scenario 2: Schedule Recurring Actions

```python
from backend.consciousness.temporal import (
    TemporalScheduler,
    ScheduledEventCreate,
    ScheduleEventType,
    create_daily_schedule,
)

scheduler = TemporalScheduler()

# Daily morning briefing
await scheduler.schedule(create_daily_schedule(
    title="Generate morning briefing",
    hour=7,
    minute=0,
    action_type="workflow",
    action_target="morning-briefing",
))

# Register the handler
async def run_morning_briefing(event):
    # Generate and send briefing
    pass

scheduler.register_handler("workflow", run_morning_briefing)
await scheduler.start()
```

### Scenario 3: Track Deadlines from Multiple Sources

```python
from backend.consciousness.temporal import get_temporal_service

service = get_temporal_service()

# Add deadlines from different sources
await service.add_deadline(
    title="PR Review: Feature X",
    due=datetime(2025, 2, 5, 17, 0),
    source="github",
    priority=2,
)

await service.add_deadline(
    title="Submit quarterly report",
    due=datetime(2025, 2, 7, 23, 59),
    source="manual",
    priority=3,  # High priority
)

# Get deadline pressure
context = await service.get_context()
print(f"Deadline pressure: {context.deadline_pressure:.0%}")
# Output: "Deadline pressure: 45%"
```

### Scenario 4: Optimize Task Scheduling

```python
from backend.consciousness.temporal.tools import suggest_optimal_time

# Find best time for focused work
result = await suggest_optimal_time(
    task_type="focus_work",
    duration_minutes=120,
    user_id="zeus",
)

print(result)
# {
#   "task_type": "focus_work",
#   "duration_minutes": 120,
#   "suggestions": [
#     {"start": "2025-02-04T09:00:00", "productivity_score": "85%", "day": "Tuesday"},
#     {"start": "2025-02-04T14:00:00", "productivity_score": "78%", "day": "Tuesday"},
#   ],
#   "recommendation": {"start": "2025-02-04T09:00:00", ...}
# }
```

---

## API Reference

### Endpoints

| Endpoint | Method | Description | Parameters |
|----------|--------|-------------|------------|
| `/context` | GET | Complete temporal context | `timezone`, `user_id` |
| `/now` | GET | Current time | `timezone` |
| `/calendar` | GET | Calendar events | `start`, `end`, `days`, `user_id` |
| `/calendar/current` | GET | Currently active event | `user_id` |
| `/calendar/next` | GET | Next upcoming event | `user_id` |
| `/deadlines` | GET | Deadlines by urgency | `urgency`, `days`, `user_id` |
| `/deadlines` | POST | Create deadline | Body: `title`, `due`, `source`, `priority` |
| `/schedules` | GET | Pending schedules | `days`, `status`, `type`, `user_id` |
| `/schedules` | POST | Create schedule | Body: `type`, `scheduled_at`, `title`, ... |
| `/schedules/{id}` | DELETE | Cancel schedule | Path: `schedule_id` |
| `/patterns` | GET | All patterns | `user_id`, `refresh` |
| `/patterns/productivity` | GET | Hourly productivity | `user_id` |

### Response Examples

**GET /context**
```json
{
  "now": "2025-02-03T14:30:00-08:00",
  "timezone": "America/Los_Angeles",
  "time_of_day": "afternoon",
  "day_type": "workday",
  "current_event": null,
  "next_event": "Team Sync",
  "minutes_until_next": 30,
  "events_today": 5,
  "urgent_deadlines": 2,
  "deadline_pressure": 0.45,
  "productivity_score": 0.72,
  "energy_prediction": 0.65,
  "is_focus_time": true
}
```

**GET /deadlines**
```json
{
  "total": 5,
  "deadline_pressure": 0.45,
  "overdue": [],
  "critical": [
    {"id": "d1", "title": "Submit report", "due": "2025-02-03T17:00:00", "urgency": "critical"}
  ],
  "urgent": [
    {"id": "d2", "title": "Review PR", "due": "2025-02-03T20:00:00", "urgency": "urgent"}
  ],
  "soon": [],
  "upcoming": []
}
```

---

## Database Setup

### TimescaleDB Installation

```bash
# macOS
brew install timescaledb

# Ubuntu/Debian
sudo apt install timescaledb-2-postgresql-15

# Docker
docker run -d --name timescaledb \
  -p 5432:5432 \
  -e POSTGRES_PASSWORD=password \
  timescale/timescaledb:latest-pg15
```

### Create Database and Run Schema

```bash
# Connect to PostgreSQL
psql -U postgres

# Create database
CREATE DATABASE onemind;
\c onemind

# Enable TimescaleDB
CREATE EXTENSION IF NOT EXISTS timescaledb;

# Run schema (from project root)
\i backend/consciousness/temporal/schema.sql
```

### Schema Overview

**Hypertables (time-partitioned):**
- `event_ledger` - All events (90 day retention, compressed after 7 days)
- `biometric_series` - Wearable data (30 day retention)
- `productivity_metrics` - Productivity scores
- `temporal_anomalies` - Detected anomalies

**Continuous Aggregates (auto-refreshed):**
- `events_hourly` - Hourly event counts (refreshed every 15 min)
- `events_daily` - Daily summaries (refreshed hourly)

**Regular Tables:**
- `calendar_events` - Calendar cache
- `deadlines` - Deadline tracking
- `scheduled_events` - Future schedules
- `temporal_patterns` - Detected patterns

**Helper Functions:**
- `get_temporal_context(timestamp, user_id)` - SQL-based context
- `calculate_deadline_pressure(user_id)` - Returns 0.0-1.0
- `get_hourly_productivity(user_id, days_back)` - Productivity by hour

### Connect Service to TimescaleDB

```python
# Implement TimeSeriesStore protocol
import asyncpg
from backend.consciousness.temporal import (
    configure_temporal_service,
    TimeSeriesStore,
    EventMetric,
    TimeSeriesQuery,
    TimeSeriesResult,
)

class TimescaleStore:
    def __init__(self, dsn: str):
        self.dsn = dsn
        self.pool = None

    async def connect(self):
        self.pool = await asyncpg.create_pool(self.dsn)

    async def write_event(self, event: EventMetric):
        await self.pool.execute("""
            INSERT INTO event_ledger (time, event_type, source, data)
            VALUES ($1, $2, $3, $4)
        """, event.time, event.event_type, event.source, event.metadata)

    async def query(self, query: TimeSeriesQuery) -> TimeSeriesResult:
        # Implement query logic
        pass

# Configure service with store
store = TimescaleStore("postgresql://user:pass@localhost:5432/onemind")
await store.connect()

configure_temporal_service(timeseries_store=store)
```

---

## Integration Examples

### Complete FastAPI Integration

```python
from fastapi import FastAPI
from contextlib import asynccontextmanager

from backend.consciousness.temporal import (
    configure_temporal_service,
    temporal_router,
    TemporalScheduler,
    InMemoryCalendarProvider,
    TemporalPatternEngine,
)

scheduler = None

@asynccontextmanager
async def lifespan(app: FastAPI):
    global scheduler

    # 1. Configure temporal service
    configure_temporal_service(
        calendar_provider=InMemoryCalendarProvider(),
        pattern_detector=TemporalPatternEngine(),
    )

    # 2. Start scheduler
    scheduler = TemporalScheduler()

    # Register action handlers
    async def handle_reminder(event):
        print(f"Reminder: {event.title}")

    scheduler.register_handler("reminder", handle_reminder)
    await scheduler.start()

    yield

    # 3. Cleanup on shutdown
    await scheduler.stop()

app = FastAPI(lifespan=lifespan)

# 4. Mount temporal API
app.include_router(temporal_router, prefix="/api/temporal")
```

### NATS Event Integration

```python
from backend.consciousness.core.bus import EventBus, Event
from backend.consciousness.temporal import get_temporal_service

service = get_temporal_service()

async def record_to_timescale(event: Event):
    """Record NATS events to TimescaleDB."""
    await service.record_event(
        event_type=str(event.type),
        source=event.source,
        metadata=event.data,
    )

# Subscribe to all events
bus = EventBus()
await bus.connect()
await bus.subscribe("events.>", record_to_timescale)
await bus.listen()
```

### Agent with Full Time Awareness

```python
from agno import Agent
from backend.consciousness.temporal import (
    TEMPORAL_TOOLS,
    get_temporal_service,
    inject_temporal_context,
)

async def create_time_aware_agent():
    service = get_temporal_service()
    context = await service.get_context()

    return Agent(
        name="temporal-agent",
        model="claude-3-5-sonnet",
        tools=TEMPORAL_TOOLS,
        system_prompt=f"""You are a highly time-aware assistant.

{context.to_prompt_context()}

INSTRUCTIONS:
1. Always consider current time context in your responses
2. Use get_deadlines to check urgent items before suggesting new tasks
3. Use suggest_optimal_time when scheduling anything
4. Warn user if they're trying to schedule during low-productivity hours
5. Respect focus windows - minimize interruptions during peak hours
""",
    )
```

---

## Troubleshooting

### Common Issues

**1. Import Errors**

```python
# Wrong
from consciousness.temporal import ...

# Correct
from backend.consciousness.temporal import ...
```

**2. Async Context**

```python
# Wrong - calling async function without await
context = service.get_context()

# Correct
context = await service.get_context()
```

**3. Timezone Issues**

```python
from zoneinfo import ZoneInfo

# Always use timezone-aware datetimes
from datetime import datetime
now = datetime.now(ZoneInfo("America/Los_Angeles"))

# Not this
now = datetime.now()  # Naive datetime
```

**4. TimescaleDB Not Connected**

The service works without TimescaleDB - patterns and time-series queries will return defaults:
- `productivity_score` defaults to 0.5
- `patterns` returns empty dict
- `events_last_hour` returns 0

**5. Calendar Provider Not Configured**

Without a calendar provider, calendar fields return None/empty:
```python
context.current_event  # None
context.next_event     # None
context.today_events   # []
```

### Debug Mode

```python
import logging
logging.getLogger("consciousness.temporal").setLevel(logging.DEBUG)
```

---

## Summary

| Component | Purpose | Status |
|-----------|---------|--------|
| `models.py` | Data structures | ✅ Complete |
| `schema.sql` | TimescaleDB schema | ✅ Complete |
| `service.py` | Core service | ✅ Structure ready |
| `scheduler.py` | Event scheduling | ✅ Structure ready |
| `patterns.py` | Pattern detection | ✅ Structure ready |
| `calendar.py` | Calendar providers | ✅ Interfaces ready |
| `context.py` | Context injection | ✅ Complete |
| `tools.py` | Agent tools (11) | ✅ Complete |
| `api.py` | REST endpoints | ✅ Complete |

### Wiring Checklist

- [ ] Run `schema.sql` on TimescaleDB
- [ ] Implement `TimeSeriesStore` protocol
- [ ] Implement Google/Apple calendar providers
- [ ] Configure OAuth credentials for calendars
- [ ] Wire NATS events → `event_ledger`
- [ ] Start scheduler in app lifespan
- [ ] Mount `temporal_router` on FastAPI app
- [ ] Add `TEMPORAL_TOOLS` to agents

### Next Steps

1. **Connect TimescaleDB** - Run schema, implement store
2. **Implement Calendar OAuth** - Google, Apple Calendar APIs
3. **Wire NATS → TimescaleDB** - Event persistence
4. **Deploy Scheduler** - Background task for triggers
5. **Pattern Detection** - Schedule periodic analysis
