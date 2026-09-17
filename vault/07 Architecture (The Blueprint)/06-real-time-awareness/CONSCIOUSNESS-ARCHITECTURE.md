# Legacy Consciousness Architecture

> The awareness layer that makes Legacy "alive"

## Overview

The consciousness layer is Legacy's nervous system - the continuous flow of all experiences, actions, and context in real-time. Without this stream, Legacy is just another chatbot. With it, she's aware of everything happening around you.

```
                    ┌──────────────────────────────────────────────────────────┐
                    │                     EXTERNAL WORLD                        │
                    │   Todoist  GitHub  Calendar  Slack  HomeAssistant  etc.   │
                    └─────────────────────────┬────────────────────────────────┘
                                              │
                                              ▼
                    ┌──────────────────────────────────────────────────────────┐
                    │              REDPANDA CONNECT (HTTP → NATS)               │
                    │   Webhook detection, transformation, routing              │
                    └─────────────────────────┬────────────────────────────────┘
                                              │
                                              ▼
    ┌─────────────────────────────────────────────────────────────────────────────────┐
    │                           NATS JETSTREAM (The Heart)                             │
    │                                                                                  │
    │   ┌─────────────┐   ┌─────────────┐   ┌─────────────┐   ┌─────────────┐         │
    │   │ events.>    │   │ action.>    │   │ legacy.>    │   │ urgent.>    │         │
    │   │ Inbound     │   │ Outbound    │   │ Legacy      │   │ High        │         │
    │   │ Events      │   │ Actions     │   │ Inbox       │   │ Priority    │         │
    │   └─────────────┘   └─────────────┘   └─────────────┘   └─────────────┘         │
    │                                                                                  │
    │   ┌─────────────────────────────────────────────────────────────────────────┐   │
    │   │                         NATS FEATURES                                    │   │
    │   │   KV Store (State)  │  Object Store (Files)  │  $SRV (Health/Stats)     │   │
    │   └─────────────────────────────────────────────────────────────────────────┘   │
    └─────────────────────────────────────────────────────────────────────────────────┘
                                              │
              ┌───────────────────────────────┼───────────────────────────────┐
              │                               │                               │
              ▼                               ▼                               ▼
    ┌─────────────────┐             ┌─────────────────┐             ┌─────────────────┐
    │   PROCESSOR     │             │   CORRELATOR    │             │   ARCHIVER      │
    │   (Awareness    │             │   (Entity       │             │   (TimescaleDB  │
    │    Filtering)   │             │    Linking)     │             │    Storage)     │
    └────────┬────────┘             └────────┬────────┘             └────────┬────────┘
             │                               │                               │
             └───────────────────────────────┼───────────────────────────────┘
                                             │
                                             ▼
                               ┌─────────────────────────┐
                               │     Legacy AGENT        │
                               │     (The Brain)         │
                               └────────────┬────────────┘
                                            │
                                            ▼
                               ┌─────────────────────────┐
                               │    ACTION GATEWAY       │
                               │    (The Hands)          │
                               └─────────────────────────┘
```

## Components

### 1. Event Bus (`bus.py`)

The nervous system - NATS JetStream event streaming.

```python
from agents.legacy.consciousness import EventBus, Event, EventType

# Create and connect
bus = await create_event_bus("nats://localhost:4222")

# Publish events
await bus.publish(Event(
    type=EventType.TASK_CREATED,
    source="todoist",
    data={"title": "Call investor", "due": "tomorrow"},
    priority=2
))

# Subscribe to events (push model - events come to you!)
async def on_task_event(event: Event):
    print(f"Got event: {event}")

await bus.subscribe("todoist.*", on_task_event)
await bus.listen()
```

**Features:**
- Sub-millisecond latency (200-400μs vs 1-5ms for Redis Streams)
- Push model (events come to you, no polling)
- Persistence (events survive restarts)
- Stream subjects: `events.>`, `action.>`, `legacy.>`, `urgent.>`

### 2. KV Store (`kv_store.py`)

State management - replaces Redis entirely.

```python
from agents.legacy.consciousness import create_kv_store, AwarenessLevel

kv = await create_kv_store(js)

# Awareness state
await kv.set_awareness_level(AwarenessLevel.ACTIVE)
level = await kv.get_awareness_level()

# Event tracking
await kv.set_last_seen_seq(12345)
seq = await kv.get_last_seen_seq()

# Cache operations
await kv.cache_set("user.preferences", {"theme": "dark"}, ttl_seconds=3600)
prefs = await kv.cache_get("user.preferences")

# Watch for changes (reactive)
async for key, value in kv.watch_state("awareness.*"):
    print(f"State changed: {key} = {value}")
```

**Buckets:**
- `legacy-state`: Critical state (awareness, last_seen_seq)
- `legacy-config`: Configuration (filters, preferences)
- `legacy-cache`: Temporary data with TTL

### 3. Object Store (`object_store.py`)

File and memory storage with automatic chunking.

```python
from agents.legacy.consciousness import create_memory_store

store = await create_memory_store(js)

# Save and load memories
await store.save_memory("conversation_2024_01_15", memory_data)
memory = await store.load_memory("conversation_2024_01_15")

# Daily summaries
await store.save_daily_summary(date.today(), summary_text)

# Context snapshots
await store.save_context_snapshot("session_123", context_data)
```

**Buckets:**
- `legacy-memories`: Conversation summaries, learned facts
- `legacy-attachments`: Files, images, documents

### 4. Event Processor (`processor.py`)

The gatekeeper - filters events by awareness level.

```python
from agents.legacy.consciousness import EventProcessor

processor = EventProcessor(nats_url="nats://localhost:4222")
await processor.start()

# Events are filtered by awareness level:
# - DORMANT: Only critical (priority >= 4)
# - AMBIENT: High priority (priority >= 3)
# - ACTIVE: Normal events (priority >= 2)
# - ENGAGED: Most events (priority >= 1)
# - OMNIPRESENT: Everything (priority >= 0)
```

### 5. Correlation Engine (`correlator.py`)

Links related events across services.

```python
from agents.legacy.consciousness import create_correlator, Entity

correlator = create_correlator(kv_store)

# Process events - entities extracted automatically
correlation = await correlator.process_event(event)

# Query related events
events = await correlator.get_related_events(Entity(
    type="issue",
    id="100",
    service="github",
    display="Issue #100"
))

# Get the complete story
story = await correlator.get_story(entity)
# Returns:
# "Story for GitHub Issue #100:
#   2024-01-15 10:00 | Issue #100: Fix login bug
#   2024-01-15 10:30 | Task: Fix issue #100
#   2024-01-15 14:00 | PR #42: Fix login validation
#   2024-01-15 16:00 | PR merged"
```

**Entity Extraction:**
- GitHub: issues, PRs, repos, commits
- Todoist: tasks, projects
- Slack: channels, users, mentions
- Cross-references in text (e.g., "#100" in Todoist task)

### 6. Event Replay (`replay.py`)

Catches up on missed events after restart.

```python
from agents.legacy.consciousness import EventReplay, ReplayConfig

config = ReplayConfig(
    max_replay_events=1000,    # Summarize if more
    max_replay_age_hours=24,   # Don't replay ancient events
    auto_summarize=True
)

replay = EventReplay(js, kv_store, config)
stats = await replay.replay_missed_events(handler=process_event)

# Returns ReplayStats:
# - events_replayed: 42
# - events_skipped: 0
# - duration_seconds: 1.5
```

### 7. Event Archiver (`archiver.py`)

TimescaleDB storage for historical queries.

```python
from agents.legacy.consciousness import create_archiver

archiver = await create_archiver(js, db_url="postgresql://...")
await archiver.start()

# Query historical events
result = await archiver.query(
    start=datetime(2024, 1, 1),
    event_type="github.pr.*",
    limit=100
)

# Activity summary
summary = await archiver.get_activity_summary(
    start=datetime.now() - timedelta(days=7)
)
# Returns:
# {
#   "total_events": 1234,
#   "inbound_events": 1000,
#   "outbound_actions": 234,
#   "by_source": {"github": 500, "todoist": 300, ...},
#   "top_types": {"pr.opened": 50, "task.completed": 45, ...}
# }

# Full-text search
result = await archiver.search_text("investor meeting")
```

**TimescaleDB Features:**
- Hypertable with automatic partitioning
- 90-day retention policy
- 7-day compression policy
- GIN index for JSONB search

### 8. Dead Letter Queue (`dlq.py`)

Handles failed event processing with retries.

```python
from agents.legacy.consciousness import create_dlq

dlq = await create_dlq(js, kv_store)
await dlq.start()

# Failed events go here automatically
# Retry schedule: 5s, 30s, 2min, 10min, 1hr (max 5 retries)

# Manual operations
await dlq.retry_entry(entry_id)
await dlq.dismiss_entry(entry_id)

# Statistics
stats = dlq.get_stats()
# {
#   "total_entries": 5,
#   "by_error_type": {"timeout": 3, "rate_limit": 2},
#   "retry_queue_size": 2,
#   "permanent_failures": 1
# }
```

### 9. Service Framework (`service.py`)

NATS $SRV protocol for health and discovery.

```python
from agents.legacy.consciousness import create_service, record_request

service = await create_service(
    nc,
    name="legacy-consciousness",
    version="1.0.0",
    description="Legacy's awareness layer"
)

# Record requests for stats
record_request("process_event", 0.05, success=True)

# Responds to:
# - $SRV.PING - Service discovery
# - $SRV.STATS.legacy-consciousness - Service statistics
# - $SRV.INFO.legacy-consciousness - Service details
```

### 10. Action Logger (`action_logger.py`)

Logs outbound actions for bidirectional awareness.

```python
from agents.legacy.consciousness import set_action_bus, action_aware

set_action_bus(bus)

# Decorator for async functions
@action_aware
async def create_task(title: str, due: str):
    # Creates task in Todoist
    return task

# Decorator for sync functions
@sync_action_aware
def send_message(channel: str, text: str):
    # Sends Slack message
    return response

# Actions are automatically published to NATS:
# Subject: action.todoist.create_task
# Data: {"action": "create_task", "service": "todoist", "params": {...}, "result": {...}}
```

## Webhook Integration (Redpanda Connect)

The `connect.yaml` configuration handles webhook transformation:

**Supported Sources:**
- Todoist (tasks: created, completed, updated)
- GitHub (PRs, issues, commits, releases)
- Google Calendar (events, meetings)
- Home Assistant (state changes, automations)
- Slack (messages, mentions, reactions)
- Discord (messages, reactions)
- Zoho CRM (contacts, deals)
- Linear (issues, projects)
- Notion (pages, databases)
- Composio (any integrated service)

**Webhook Endpoint:**
```
POST http://your-server:4196/webhooks
```

Add header `X-Webhook-Source: service_name` or let Redpanda Connect auto-detect.

## Event Flow Example

### Task Lifecycle

```
1. User creates task in Todoist
   └─> Webhook → Redpanda Connect → events.todoist.task.added

2. Event Processor filters by awareness
   └─> If ACTIVE mode: passes through
   └─> Updates KV Store (last_seen_seq)

3. Correlator extracts entities
   └─> Creates new correlation if task mentions "#123" (issue)

4. Archiver stores in TimescaleDB
   └─> Available for "What did I do yesterday?"

5. Legacy receives via legacy.inbox
   └─> May decide to take action

6. Legacy creates GitHub issue
   └─> Action Logger publishes to action.github.create_issue

7. Correlator links the new issue to original task
   └─> Full story now spans both services
```

## Environment Variables

```env
# Required
NATS_URL=nats://nats:4222

# Optional
TIMESCALE_URL=postgresql://onemind:password@timescaledb:5432/events
```

## API Reference

### Event Types

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
```

### Awareness Levels

```python
class AwarenessLevel(IntEnum):
    DORMANT = 0      # Only critical (priority >= 4)
    AMBIENT = 1      # High priority (priority >= 3)
    ACTIVE = 2       # Normal (priority >= 2) - DEFAULT
    ENGAGED = 3      # Most events (priority >= 1)
    OMNIPRESENT = 4  # Everything (priority >= 0)
```

### Priority Levels

| Priority | Description | Examples |
|----------|-------------|----------|
| 0 | Debug | Internal metrics, verbose logs |
| 1 | Low | FYI events, non-urgent updates |
| 2 | Normal | Standard events (default) |
| 3 | High | IoT alerts, @mentions |
| 4 | Urgent | Calendar conflicts, deadlines |
| 5 | Critical | Security events, failures |

## Deployment

### Docker Compose Services

```yaml
services:
  nats:
    image: nats:alpine
    ports:
      - "4222:4222"  # Client
      - "8222:8222"  # Monitoring

  redpanda-connect:
    image: redpandadata/connect:latest
    ports:
      - "4196:4196"  # Webhook ingestion

  timescaledb:
    image: timescale/timescaledb:latest-pg15
    ports:
      - "5433:5432"

  event-processor:
    environment:
      NATS_URL: nats://nats:4222
      TIMESCALE_URL: postgresql://...
```

### Health Checks

All components respond to NATS $SRV protocol:

```bash
# Ping all services
nats request '$SRV.PING' ''

# Get stats for consciousness
nats request '$SRV.STATS.legacy-consciousness' ''
```

## Migration from Redis

The KV Store completely replaces Redis for:

| Redis Usage | KV Store Equivalent |
|-------------|---------------------|
| `GET/SET awareness_level` | `kv.get/set_awareness_level()` |
| `GET/SET last_seen_seq` | `kv.get/set_last_seen_seq()` |
| `HSET/HGET config` | `kv.config_get/set()` |
| `SETEX cache` | `kv.cache_set(ttl_seconds=...)` |
| `SUBSCRIBE channel` | `kv.watch_state(pattern)` |

**Benefits:**
- Single infrastructure (NATS handles everything)
- Built-in persistence
- Watch/reactive updates
- No separate Redis container needed


---

## Related

- [[00-24 SO (Sovereign Operations)/01 Intelligence (AI, Memory & Tools) (The Brain)/01 Intelligence|01 Intelligence]]
