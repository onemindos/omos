# Consciousness Layer — Architecture Blueprint

> The awareness system that makes Legacy "alive"

**Origin**: Consolidated from `reference/consciousness/` (79 files, ~30K lines).
Reference implementations preserved in `reference/` subfolder.

---

## Table of Contents

1. [Core Concept](#core-concept)
2. [Architecture Overview](#architecture-overview)
3. [Three-Layer Design](#three-layer-design)
4. [NATS Event Bus](#nats-event-bus)
5. [Awareness Levels](#awareness-levels)
6. [Context Injection](#context-injection)
7. [Event Processing Pipeline](#event-processing-pipeline)
8. [Activity Feed & Triage](#activity-feed--triage)
9. [Temporal Infrastructure](#temporal-infrastructure)
10. [Deep Links](#deep-links)
11. [Data Models](#data-models)
12. [NATS Best Practices](#nats-best-practices)
13. [Build Sequence](#build-sequence)
14. [Reference Files](#reference-files)

---

## Core Concept

OneMind OS is fundamentally different from chat-based AI assistants. Rather than waiting for user queries, it **continuously receives events, maintains state, and acts proactively**.

```
Traditional (Reactive):  User asks → System searches → Returns answer
OneMind OS (Proactive):  Events flow → Context builds → Agent acts autonomously
```

**Without consciousness**: Agent only knows what user tells it — "I don't know your schedule."
**With consciousness**: Agent gets real-time context injected — "You've been coding for 45 minutes. Your project proposal is due in 2.5 hours. I'd suggest taking a break, then focusing on it."

---

## Architecture Overview

```
┌──────────────────────────────────────────────────────────────┐
│                     EXTERNAL WORLD                            │
│   Todoist  GitHub  Zoho  HomeAssistant  Composio  Oura       │
└─────────────────────────┬────────────────────────────────────┘
                          │ webhooks / polling
                          ▼
┌──────────────────────────────────────────────────────────────┐
│                 NATS JETSTREAM (The Heart)                    │
│                                                               │
│   events.>     action.>     legacy.>     urgent.>             │
│   (Inbound)    (Outbound)   (Inbox)      (Priority)           │
│                                                               │
│   KV Store (State)  │  Object Store (Files)  │  $SRV (Health) │
└─────────────────────────┬────────────────────────────────────┘
                          │
        ┌─────────────────┼─────────────────┐
        │                 │                 │
        ▼                 ▼                 ▼
┌───────────────┐ ┌───────────────┐ ┌───────────────┐
│   PROCESSOR   │ │  CORRELATOR   │ │   ARCHIVER    │
│   (Filter)    │ │  (Linking)    │ │  (Storage)    │
└───────┬───────┘ └───────┬───────┘ └───────┬───────┘
        │                 │                 │
        └─────────────────┼─────────────────┘
                          │
         ┌────────────────┼────────────────┐
         │                │                │
         ▼                ▼                ▼
   ┌───────────┐   ┌───────────┐   ┌───────────┐
   │ ACTIVITY  │   │   Agent   │   │  Mobile   │
   │ SERVICE   │   │ (action)  │   │  (push)   │
   └───────────┘   └───────────┘   └───────────┘
```

### End-to-End Event Flow

1. External services emit webhooks/events
2. NATS JetStream receives and persists events (~200-400μs latency)
3. **Processor** filters by awareness level (noise filtering, batching)
4. **Correlator** links related events across services (Issue → Task → PR → Merge = one story)
5. **Archiver** stores to TimescaleDB for historical queries
6. **Activity Service** creates actionable items, broadcasts to UI
7. **Context Injector** wraps agent calls with live awareness

---

## Three-Layer Design

### Layer 1: The Nervous System (NATS EventBus)

All external events flow through NATS JetStream. Push model — events delivered instantly, no polling. 7-day retention, auto-reconnect.

### Layer 2: The Gatekeeper (LightweightProcessor)

Filters events based on awareness level. See [Awareness Levels](#awareness-levels).

Features:
- Noise filtering (typing indicators, heartbeats)
- Time-aware thresholds (night vs day)
- Location-aware (per-room awareness)
- Batch grouping (50 motion events → "Motion in 5 areas")
- Critical event bypass (emergencies skip all filters)
- Configurable rules via KV Watch (real-time, no polling)

### Layer 3: The Context Injector (AwareRunner)

Automatically injects real-time context into agent prompts. The core pattern:

```python
# Instead of:
response = await agent.arun(message)

# Use:
response = await run_with_awareness(agent, message, session_id=session_id)
```

Injected context includes:
```markdown
# Real-Time Awareness (2024-01-31 14:30 UTC)

## Current State
📍 Home Office | 🎯 Working (45min) | 💓 HR:72 | Stress:23%

## Upcoming Deadlines
🔴 Project proposal (due: 17:00 - 2.5h left)
🟡 Weekly report (due: 2024-02-01 09:00)

## Recent Events
⭐ Task completed: Review PR #123 (30min ago)
📅 Meeting starting in 30min
🏠 Motion detected in garage
```

| Context Source | Provider | Update Speed |
|---------------|----------|-------------|
| Recent events | `EventBus.get_recent()` | Real-time |
| User state | `StateManager.get_context_summary()` | Real-time |
| Deadlines | `TemporalService.get_upcoming_deadlines()` | Real-time |
| Config | `KVStore.get_state_snapshot()` | Real-time (KV Watch) |

---

## NATS Event Bus

### Subject Taxonomy

```python
class EventType(str, Enum):
    # Todoist
    TASK_CREATED = "todoist.task.created"
    TASK_COMPLETED = "todoist.task.completed"
    TASK_UPDATED = "todoist.task.updated"
    TASK_DUE_SOON = "todoist.task.due_soon"

    # GitHub
    PR_OPENED = "github.pr.opened"
    PR_MERGED = "github.pr.merged"
    PR_REVIEW_REQUESTED = "github.pr.review_requested"
    ISSUE_CREATED = "github.issue.created"
    COMMIT_PUSHED = "github.commit.pushed"

    # Calendar
    MEETING_SOON = "calendar.meeting.soon"
    MEETING_STARTED = "calendar.meeting.started"
    SCHEDULE_CONFLICT = "calendar.schedule.conflict"

    # Home Assistant
    MOTION_DETECTED = "home.motion.detected"
    DOOR_OPENED = "home.door.opened"
    DEVICE_STATE_CHANGED = "home.device.changed"

    # System
    AGENT_STARTED = "system.agent.started"
    AGENT_ERROR = "system.agent.error"

    # Biometric (Oura, Apple Watch, etc.)
    BIOMETRIC_HEART = "biometric.heart"
    BIOMETRIC_STRESS_HIGH = "biometric.stress.high"
    BIOMETRIC_SLEEP_STARTED = "biometric.sleep.started"
    BIOMETRIC_SLEEP_ENDED = "biometric.sleep.ended"
    BIOMETRIC_ALERT_FALL = "biometric.alert.fall"          # Critical
    BIOMETRIC_ALERT_HR_CRITICAL = "biometric.alert.hr_critical"  # Critical

    # Environment (IoT sensors)
    ENVIRONMENT_TEMPERATURE = "environment.temperature"
    ENVIRONMENT_AIR_QUALITY = "environment.air_quality"

    # Presence
    PRESENCE_MOTION = "presence.motion"
    PRESENCE_LOCATION_CHANGED = "presence.location.changed"
```

### Stream Configuration

| Stream | Subjects | Purpose |
|--------|----------|---------|
| `events` | `events.>` | All inbound events |
| `action` | `action.>` | Outbound agent actions |
| `legacy` | `legacy.>` | Filtered inbox for Legacy |
| `urgent` | `urgent.>` | Priority bypass channel |

### KV Store (Replaces Redis)

Three buckets:
- `legacy-state`: Awareness level, last_seen_seq (critical state)
- `legacy-config`: Filter rules, preferences
- `legacy-cache`: Temporary data with TTL

Uses **KV Watch** for reactive updates — no polling.

### Why NATS?

NATS replaces multiple services:

| Need | Traditional | NATS |
|------|-------------|------|
| Event streaming | Kafka + Zookeeper | JetStream |
| State cache | Redis | KV Store |
| File storage | S3/MinIO | Object Store |
| Service discovery | Consul | $SRV Protocol |
| Retry logic | Custom code | Native backoff |
| Dead letter queue | Custom stream | Native advisories |

Single binary (~10MB), sub-millisecond latency (~100μs), no operational complexity.

---

## Awareness Levels

Four levels control what reaches Legacy:

| Level | Value | Threshold | Use Case |
|-------|-------|-----------|----------|
| `DORMANT` | 0 | Priority ≥ 4 | Sleeping — only critical alerts |
| `AWARE` | 1 | Priority ≥ 3 | Watching — high priority + approvals |
| `PRESENT` | 2 | Priority ≥ 1 | Engaged — most events (**DEFAULT**) |
| `OMNIPRESENT` | 3 | Priority ≥ 0 | Full presence — everything |

### Priority Scale

| Priority | Description | Examples |
|----------|-------------|----------|
| 0 | Debug | Internal metrics, verbose logs |
| 1 | Low | FYI events, non-urgent updates |
| 2 | Normal | Standard events (default) |
| 3 | High | IoT alerts, @mentions |
| 4 | Urgent | Deadlines, conflicts |
| 5 | Critical | Security events, failures, fall detection |

### Smart Filtering

The processor goes beyond simple threshold matching:
- **Time-aware**: Night mode filters more aggressively
- **Location-aware**: Per-room/per-site awareness levels
- **Batch grouping**: High-volume events get summarized (50 motion events → "Motion in 5 areas")
- **Critical bypass**: Priority ≥ 4 always gets through regardless of awareness level
- **Noise filtering**: Typing indicators, heartbeats, status pings automatically dropped

---

## Context Injection

### The Core Pattern: `run_with_awareness()`

```python
async def run_with_awareness(
    agent, message, *,
    session_id=None,
    include_events=True,
    include_state=True,
    include_deadlines=True,
    event_count=10,
    additional_context=None,
):
    # 1. Build awareness context from all sources
    awareness_context = await get_agent_context(
        include_events=include_events,
        include_state=include_state,
        include_deadlines=include_deadlines,
        event_count=event_count,
    )

    # 2. Inject into agent's additional_context parameter
    response = await agent.arun(
        message=message,
        session_id=session_id,
        additional_context=awareness_context,
    )
    return response
```

### Context Assembly (`get_agent_context()`)

Builds formatted markdown from multiple sources:
1. **State summary**: Location, activity, focus task, biometrics
2. **Deadlines**: Top 5. Emoji-coded by urgency (🔴 critical → 🟢 low)
3. **Recent events**: Last N events with priority stars

### Usage Options

1. **Explicit injection**: `await run_with_awareness(agent, message)`
2. **Wrapped agent**: `AwareAgent(agent)` — wraps `.run()` and `.stream()`
3. **Just get context**: `await get_agent_context()` or `get_agent_context_dict()`

See: `reference/aware_runner.py`, `reference/awareness_provider.py`

---

## Event Processing Pipeline

### Processor Pipeline (7 steps)

1. **Noise filter** — Drop typing indicators, heartbeats
2. **Critical check** — Bypass all filters for emergencies
3. **Time-aware filtering** — Check schedule/time overrides
4. **Location-aware filtering** — Per-site awareness levels
5. **Awareness filter** — Respect Legacy's awareness state
6. **Batch check** — Group high-volume events
7. **Enrich and forward** — Add human-readable context, publish to `legacy.inbox`

See: `reference/processor.py`

### Correlation Engine

Links related events across services to build complete stories:

```
10:00  GitHub Issue #100 created
10:05  Todoist task created mentioning "#100"
10:30  PR #42 created fixing issue #100
11:00  PR merged
11:05  Todoist task completed
→ Correlator links these as ONE workflow
```

**Entity extraction** via regex:
- GitHub: issue/PR numbers, repo names, commit SHAs
- Todoist: task IDs, project IDs
- Cross-references in text (e.g., "#100" in Todoist task title)

**Relationship graph**: Entity → Correlation mapping stored in NATS KV.
**Story queries**: `correlator.get_story(entity)` returns timeline narrative.

See: `reference/correlator.py`

### Archiver (TimescaleDB)

- Hypertable with automatic partitioning
- 90-day retention, 7-day compression
- GIN index for JSONB full-text search
- Continuous aggregates: `events_hourly` (15-min refresh), `events_daily` (hourly refresh)

### Dead Letter Queue

Uses native NATS advisories (`$JS.EVENT.ADVISORY.CONSUMER.MAX_DELIVERIES.>`).
Retry schedule: 5s → 30s → 2min → 10min → 1hr (max 5 retries).

---

## Activity Feed & Triage

### Event → Activity Item Mapping

```python
EVENT_TYPE_MAPPING = {
    "github.pull_request.review_requested" → APPROVAL (HIGH),
    "todoist.item.added" → TASK (NORMAL),
    "homeassistant.alarm_triggered" → ALERT (CRITICAL),
    "hitl.approval_required" → APPROVAL (HIGH),
}
```

### Auto-Assignment Rules

| Type | Assignment |
|------|------------|
| `APPROVAL` | Always HUMAN |
| `ALERT` | Always HUMAN |
| `TASK` (low priority) + OMNIPRESENT | Agent can handle |
| `SUGGESTION` + PRESENT+ | Agent can auto-handle |

### AI Triage Engine (Three Paths + GTD)

LLM-powered classification using Claude Haiku:

**Three Paths (Pillars)**:
- **HP** (Holistic Performance): Health, fitness, wellness, self-development
- **LE** (Legacy Evolution): Family, relationships, home, community
- **GE** (Generational Endeavor): Business, career, finances, investments

**GTD Action Types**:
- **do**: < 2 minutes, actionable immediately
- **delegate**: Assign to agent or someone else
- **defer**: Schedule for later
- **delete**: Not actionable, dismiss
- **reference**: Useful info to save

**Triage prompt** classifies items → returns JSON with pillar, priority, action_type, suggested_due, confidence, reasoning.

**Heuristic shortcuts** for common cases (keyword matching before LLM call):
- Critical keywords: urgent, emergency, asap, 911
- HP keywords: health, fitness, workout, sleep, meditation
- LE keywords: family, kids, home, school, birthday
- GE keywords: business, project, meeting, revenue, investment

See: `reference/triage.py`, `reference/activity_models.py`

---

## Temporal Infrastructure

> Time is "baked in" to the OS. It's permanent infrastructure, not an agent.

### Why Not a "Chronos Agent"?

| Aspect | Agent Approach | Infrastructure Approach |
|--------|---------------|------------------------|
| Lifespan | Per-conversation | Always running |
| Access | Via @mention | Injected everywhere |
| Storage | Queries existing data | Owns time-series data |
| Pattern Detection | On-demand | Continuous background |
| Scheduling | Manual | Automatic triggers |

### TemporalContext (injected into every agent)

```python
class TemporalContext:
    # Current moment
    now: datetime
    timezone: str
    time_of_day: TimeOfDay      # early_morning|morning|afternoon|evening|night|late_night
    day_type: DayType           # workday|weekend|holiday|vacation

    # Calendar
    current_event: Optional[CalendarEvent]
    next_event: Optional[CalendarEvent]
    minutes_until_next: Optional[int]
    events_count_today: int

    # Deadlines
    urgent_deadlines: List[Deadline]
    deadline_pressure: float    # 0.0–1.0
    overdue_count: int

    # Patterns (from TimescaleDB)
    productivity_score: float   # 0.0–1.0 for current hour
    energy_prediction: float    # 0.0–1.0
    typical_load: str           # lighter|normal|heavier

    # Focus
    is_focus_time: bool
    focus_windows_today: List[TimeRange]
    next_focus_window: Optional[TimeRange]
```

### Time Classifications

```python
class TimeOfDay(str, Enum):
    EARLY_MORNING = "early_morning"  # 5am-7am
    MORNING = "morning"              # 7am-12pm
    AFTERNOON = "afternoon"          # 12pm-5pm
    EVENING = "evening"              # 5pm-9pm
    NIGHT = "night"                  # 9pm-12am
    LATE_NIGHT = "late_night"        # 12am-5am

class DeadlineUrgency(IntEnum):
    OVERDUE = 0       # Past due
    CRITICAL = 1      # < 1 hour
    URGENT = 2        # < 4 hours
    SOON = 3          # < 24 hours
    UPCOMING = 4      # < 1 week
    LATER = 5         # > 1 week
```

### Scheduler

Supports one-time + recurring (RRULE) schedules:
- Reminders, task deadlines, meeting triggers
- Daily/weekly recurring (morning briefing, weekly review)
- Action handlers: workflow triggers, notification dispatch
- Snooze, cancel, query pending

### Pattern Detection

Analyzes TimescaleDB aggregates to detect:
- Productivity peaks/dips by hour and day
- Optimal deep work windows
- Energy level patterns
- Stress patterns from biometrics
- Deadline behavior patterns

See: `reference/temporal_models.py`, `reference/temporal_context.py`

---

## Deep Links

Native app URL schemes for 12+ apps. Used by Canvas and Notifier for direct content access.

| App | Pattern | Example |
|-----|---------|---------|
| Gmail | `googlegmail:///co?messageId={id}` | Open specific email |
| Outlook | `ms-outlook://emails/{id}` | Open email |
| GitHub | `github://pull-request/{owner}/{repo}/{num}` | Open PR |
| GitHub Issue | `github://issue/{owner}/{repo}/{num}` | Open issue |
| Todoist Task | `todoist://task?id={id}` | Open task |
| Todoist Add | `todoist://addtask?content={text}` | Quick add |
| Apple Calendar | `calshow:{eventId}` | Open event |
| Google Calendar | `googlecalendar://event?eventId={id}` | Open event |
| Slack Channel | `slack://channel?team={t}&id={c}` | Open channel |
| Discord | `discord://discord.com/channels/{g}/{c}` | Open channel |
| Home Assistant | `homeassistant://navigate/lovelace/default_view?entity_id={id}` | Open entity |
| Zoho CRM | `https://crm.zoho.com/crm/tab/{module}/{id}` | Open record |
| OneMind | `onemind://canvas/{item_id}` | Open canvas item |

See: `reference/deep_links.py`

---

## Data Models

### Activity Item

```python
class ActivityItemType(str, Enum):
    MESSAGE, NOTIFICATION, EVENT, APPROVAL, TASK, ALERT, REMINDER, SUGGESTION

class Priority(str, Enum):
    CRITICAL, HIGH, NORMAL, LOW

class ActivityStatus(str, Enum):
    NEW, READING, HANDLED, ARCHIVED, SNOOZED, FILTERED

class AssigneeType(str, Enum):
    HUMAN, AGENT, TEAM, AUTO

class Pillar(str, Enum):
    HP, LE, GE  # Three Paths
```

Key fields on `ActivityItem`:
- `id`, `type`, `priority`, `status`
- `title`, `body`, `source`
- `assignee_type`, `assignee_id`
- `approval_id`, `approval_status` (HITL items)
- `task_due_at`, `task_pillar` (triage results)
- `metadata`, `action_url`
- `session_id`, `run_id`, `workflow_id`, `event_id` (linkage)
- `awareness_mode` (what awareness level was active)
- `snoozed_until` (snooze support)

See: `reference/activity_models.py`

### Event

```python
@dataclass
class Event:
    type: Union[EventType, str]
    source: str
    data: dict
    timestamp: str      # ISO 8601, auto-set
    priority: int       # 0=low → 5=critical
    sequence: int = 0   # Set by NATS JetStream
```

### Scheduled Event

```python
class ScheduledEvent:
    type: ScheduleEventType  # reminder|task_due|meeting|deadline|recurring|trigger|digest
    status: ScheduleStatus   # pending|triggered|completed|failed|cancelled|snoozed
    scheduled_at: datetime
    is_recurring: bool
    recurrence_rule: str     # RRULE format
    action_type: str         # workflow|agent|tool|notification
    action_target: str       # workflow_id|agent_id|tool_name
    action_params: dict
```

---

## NATS Best Practices

### Native Consumer Backoff (no custom retry code)
```python
ConsumerConfig(
    durable_name="legacy-agent",
    max_deliver=5,
    backoff=[1000, 5000, 30000, 120000, 600000],  # 1s → 10min
)
```

### Dead Letter via Advisories (no custom DLQ stream)
```python
await nc.subscribe(
    "$JS.EVENT.ADVISORY.CONSUMER.MAX_DELIVERIES.>",
    cb=handle_dead_letter,
)
```

### Service Discovery ($SRV protocol)
```python
from nats.micro import add_service, ServiceConfig
svc = await add_service(nc, ServiceConfig(
    name="onemind-agent",
    version="2.0.0",
    queue_group="consciousness-workers",
))
```

### Queue Groups for Scaling
```python
await js.subscribe(
    "events.>",
    queue="consciousness-processors",
    cb=process_event,
)
# Run N instances — NATS distributes automatically
```

### Graceful Shutdown (drain mode)
```python
if _event_bus:
    await _event_bus.disconnect()  # Drains before closing
```

---

## Build Sequence

### Phase 1: Core Event Infrastructure
- [ ] NATS JetStream setup (streams: events, action, legacy, urgent)
- [ ] KV Store buckets (legacy-state, legacy-config, legacy-cache)
- [ ] Event Bus wrapper (publish, subscribe, get_recent)
- [ ] Event data model + serialization

### Phase 2: Processing Pipeline
- [ ] LightweightProcessor (noise filter + awareness filter)
- [ ] Awareness levels (4 levels, KV-backed)
- [ ] Critical event bypass
- [ ] Event batching (motion, commits, task updates)

### Phase 3: Context Injection
- [ ] `get_agent_context()` — assemble from NATS + KV
- [ ] `run_with_awareness()` — wrap OpenClaw agent calls
- [ ] Temporal context injection (time-of-day, deadlines)

### Phase 4: Correlation
- [ ] Entity extraction (GitHub, Todoist, cross-references)
- [ ] Correlation graph (KV-backed)
- [ ] Story queries

### Phase 5: Activity Feed
- [ ] Event → ActivityItem conversion
- [ ] AI Triage (Three Paths + GTD)
- [ ] Auto-assignment rules
- [ ] Activity stats

### Phase 6: Temporal Core
- [ ] TemporalContext generation
- [ ] Scheduler (one-time + RRULE recurring)
- [ ] Calendar provider interfaces
- [ ] Pattern detection (TimescaleDB aggregates)

### Phase 7: Storage & History
- [ ] TimescaleDB archiver (hypertables, retention, compression)
- [ ] Continuous aggregates (hourly, daily)
- [ ] Replay (catch-up after restart)
- [ ] DLQ via native advisories

### Phase 8: Deep Links + Notifications
- [ ] Deep link generator (12+ apps)
- [ ] Unified notifier (route by presence + priority)

---

## Reference Files

Preserved implementation files in `reference/`:

| File | Lines | What It Contains |
|------|-------|-----------------|
| `aware_runner.py` | 213 | Core context injection wrapper (`run_with_awareness`, `AwareAgent` class) |
| `awareness_provider.py` | 271 | Context assembly from NATS events, KV state, deadlines |
| `processor.py` | 570 | LightweightProcessor — full 7-step filter pipeline with time/location awareness |
| `correlator.py` | 465 | Entity extraction regex patterns, correlation graph, story queries |
| `triage.py` | 374 | AI Triage Engine — Three Paths prompt, GTD classification, heuristic shortcuts |
| `activity_models.py` | 343 | Complete Pydantic models for activity feed (items, stats, preferences, runs) |
| `bus.py` | 653 | NATS JetStream wrapper — EventType enum (all subjects), Event dataclass, stream config |
| `temporal_models.py` | 480 | Temporal data models (TimeOfDay, DeadlineUrgency, TemporalContext, ScheduledEvent) |
| `temporal_context.py` | 401 | Temporal context injection middleware, decorators, enrichers |
| `deep_links.py` | 446 | Native app URL scheme generator for 12+ apps |

Total preserved: **4,216 lines** of implementation-spec code.
