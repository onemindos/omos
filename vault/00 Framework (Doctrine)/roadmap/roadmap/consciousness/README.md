# Legacy Consciousness Layer

> The awareness system that makes Legacy "alive"

## Overview

The consciousness layer is Legacy's nervous system - a continuous flow of events, actions, and context in real-time. Without this layer, Legacy is just another chatbot. With it, she's aware of everything happening around you.

```
┌──────────────────────────────────────────────────────────────┐
│                     EXTERNAL WORLD                            │
│   Todoist  GitHub  Zoho  HomeAssistant  Composio              │
└─────────────────────────┬────────────────────────────────────┘
                          │
                          ▼
┌──────────────────────────────────────────────────────────────┐
│              REDPANDA CONNECT (HTTP → NATS)                   │
│   Webhook detection, transformation, routing                  │
└─────────────────────────┬────────────────────────────────────┘
                          │
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
                          ▼
                ┌─────────────────┐
                │  Legacy AGENT   │
                │   (The Brain)   │
                └────────┬────────┘
                         │
                         ▼
                ┌─────────────────┐
                │ ACTION LOGGER   │
                │  (The Hands)    │
                └─────────────────┘
```

## How Agents Become Aware

The key question: **How do AI agents actually see events from NATS?**

### The Problem (Without This)

```python
# Agent only knows what user tells it
agent = Agent(name="legacy", instructions="Be helpful")
response = await agent.arun("What should I work on?")
# Agent: "I don't know your schedule. What would you like to do?"
```

### The Solution (With Awareness)

```python
from backend.consciousness.aware_runner import run_with_awareness

# Agent gets real-time context injected!
response = await run_with_awareness(agent, "What should I work on?")
# Agent: "You've been coding for 45 minutes. Your project proposal is due
#         in 2.5 hours. I'd suggest taking a break, then focusing on it."
```

### What Gets Injected

When you call `run_with_awareness()`, the agent receives:

```markdown
# Real-Time Awareness (as of 2024-01-31 14:30 UTC)

## Current State
📍 Location: Home Office | 🎯 Activity: Working (45min)
💻 Task: Code review | 💓 HR:72 | Stress:23%

## Upcoming Deadlines
- 🔴 Project proposal (due: 2024-01-31 17:00)
- 🟡 Weekly report (due: 2024-02-01 09:00)

## Recent Events
- ⭐ Task completed: Review PR #123 (todoist.task.completed)
- Meeting starting in 30min (calendar.meeting.soon)
- Motion in garage (ha.motion.detected)
```

### Usage Options

**Option 1: Explicit injection (recommended)**
```python
from backend.consciousness.aware_runner import run_with_awareness
from backend.agno.agents.factory import get_agent

agent = get_agent("legacy")
response = await run_with_awareness(
    agent,
    "What's happening?",
    session_id="session_123",
)
```

**Option 2: Wrapped agent**
```python
from backend.consciousness.aware_runner import AwareAgent

aware_agent = AwareAgent(get_agent("legacy"))
response = await aware_agent.run("What should I focus on?")
```

**Option 3: Just get context**
```python
from backend.consciousness.awareness import get_agent_context

context = await get_agent_context()  # Formatted string
data = await get_agent_context_dict()  # Structured dict
```

### Context Sources

| Source | Provider | Update Speed |
|--------|----------|--------------|
| Recent events | `EventBus.get_recent()` | Real-time |
| User state | `StateManager.get_context_summary()` | Real-time |
| Deadlines | `TemporalAwareness.get_upcoming_deadlines()` | Real-time |
| Config | `KVStore.get_state_snapshot()` | Real-time (KV Watch) |

### Related Files

- `backend/consciousness/awareness.py` - Context provider functions
- `backend/consciousness/aware_runner.py` - Agent wrapper with awareness
- `backend/agno/onemindos.py` - Initializes StateManager + TemporalAwareness

---

## Modules

### Core Event System

| Module | File | Purpose |
|--------|------|---------|
| **Event Bus** | `bus.py` | NATS JetStream streaming - the nervous system |
| **KV Store** | `kv_store.py` | State management (replaces Redis) |
| **Object Store** | `object_store.py` | Memory and file storage |

### Event Processing

| Module | File | Purpose |
|--------|------|---------|
| **Processor** | `processor.py` | Filters events by awareness level |
| **Correlator** | `correlator.py` | Links related events across services |
| **Archiver** | `archiver.py` | TimescaleDB historical storage |
| **Replay** | `replay.py` | Catches up on missed events after restart |
| **DLQ** | `dlq.py` | Dead letter queue for failed events |

### Integration

| Module | File | Purpose |
|--------|------|---------|
| **Gateway** | `gateway.py` | Webhook receiver (sensory input) |
| **Composio Bridge** | `composio_bridge.py` | External service integration |
| **Action Logger** | `action_logger.py` | Logs outbound actions for awareness |
| **Service** | `service.py` | NATS $SRV protocol (health/stats) |

### Awareness & Utilities

| Module | File | Purpose |
|--------|------|---------|
| **State** | `awareness/state.py` | User context (location, activity, devices) |
| **Temporal** | `awareness/temporal.py` | Time awareness + deadline tracking |
| **Notifier** | `utils/notifier.py` | Unified notifications (Agno HITL + ntfy.sh) |
| **Deep Links** | `utils/deep_links.py` | Native app link generation |

> **Note:** Some modules (presence, canvas, patterns) have been moved to `.paused/consciousness-unused/` during simplification. They can be restored if needed.

---

## Module Details

### 1. Event Bus (`bus.py`)

The nervous system - NATS JetStream event streaming with sub-millisecond latency.

```python
from backend.consciousness import EventBus, Event, EventType

bus = await create_event_bus("nats://localhost:4222")

# Publish events
await bus.publish(Event(
    type=EventType.TASK_CREATED,
    source="todoist",
    data={"title": "Call investor", "due": "tomorrow"},
    priority=2
))

# Subscribe to events (push model)
async def on_task_event(event: Event):
    print(f"Got event: {event}")

await bus.subscribe("todoist.*", on_task_event)
await bus.listen()
```

**Features:**
- Sub-millisecond latency (200-400μs)
- Push model (events come to you, no polling)
- Persistence (events survive restarts)
- Stream subjects: `events.>`, `action.>`, `legacy.>`, `urgent.>`

---

### 2. KV Store (`kv_store.py`)

State management using NATS KV - completely replaces Redis.

```python
from backend.consciousness import create_kv_store, AwarenessLevel

kv = await create_kv_store(js)

# Awareness state
await kv.set_awareness_level(AwarenessLevel.ACTIVE)
level = await kv.get_awareness_level()

# Event tracking
await kv.set_last_seen_seq(12345)
seq = await kv.get_last_seen_seq()

# Cache with TTL
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

---

### 3. Object Store (`object_store.py`)

File and memory storage with automatic chunking.

```python
from backend.consciousness import create_memory_store

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

---

### 4. Event Processor (`processor.py`)

The gatekeeper - filters events based on Legacy's current awareness level.

```python
from backend.consciousness import EventProcessor

processor = EventProcessor(nats_url="nats://localhost:4222")
await processor.start()

# Events are filtered by awareness level (4 levels):
# - DORMANT: Only critical (priority >= 4)
# - AWARE: High priority (priority >= 3)
# - PRESENT: Most events (priority >= 1)
# - OMNIPRESENT: Everything (priority >= 0)
```

---

### 5. Correlation Engine (`correlator.py`)

Links related events across different services.

```python
from backend.consciousness import create_correlator, Entity

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
- Cross-references in text (e.g., "#100" in Todoist task)

---

### 6. Event Archiver (`archiver.py`)

TimescaleDB storage for historical queries like "What happened last week?"

```python
from backend.consciousness import create_archiver

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

# Full-text search
result = await archiver.search_text("investor meeting")
```

**TimescaleDB Features:**
- Hypertable with automatic partitioning
- 90-day retention policy
- 7-day compression policy
- GIN index for JSONB search

---

### 7. Event Replay (`replay.py`)

Catches up on missed events after restart.

```python
from backend.consciousness import EventReplay, ReplayConfig

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

---

### 8. Dead Letter Queue (`dlq.py`)

Handles failed event processing with automatic retries.

```python
from backend.consciousness import create_dlq

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

---

### 9. Service Framework (`service.py`)

NATS $SRV protocol for health checks and service discovery.

```python
from backend.consciousness import create_service, record_request

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

---

### 10. Action Logger (`action_logger.py`)

Logs outbound actions for bidirectional awareness.

```python
from backend.consciousness import set_action_bus, action_aware

set_action_bus(bus)

# Decorator for async functions
@action_aware
async def create_task(title: str, due: str):
    # Creates task in Todoist
    return task

# Decorator for sync functions
@sync_action_aware
def send_message(channel: str, text: str):
    # Sends message
    return response

# Actions automatically published to NATS:
# Subject: action.todoist.create_task
# Data: {"action": "create_task", "service": "todoist", "params": {...}, "result": {...}}
```

---

### 11. Gateway (`gateway.py`)

Webhook receiver for external services (alternative to Redpanda Connect).

```python
from backend.consciousness import Gateway

gateway = Gateway(bus)
await gateway.start(port=8000)

# Endpoints:
# POST /webhooks/todoist - Todoist events
# POST /webhooks/github - GitHub events
# POST /webhooks/zoho - Zoho CRM events
# POST /webhooks/homeassistant - Home Assistant events
# POST /webhooks/composio - Composio events
```

---

### 12. Composio Bridge (`composio_bridge.py`)

Integration with Composio for external service connections.

```python
from backend.consciousness import initialize_composio, setup_recommended_triggers

bridge = await initialize_composio(bus, api_key="...")

# Set up recommended triggers
await setup_recommended_triggers(bridge)

# Events from Composio flow through the bus as:
# events.composio.{trigger_type}
```

---

## Archived Modules (in `.paused/consciousness-unused/`)

The following modules have been moved to `.paused/` during simplification. They can be restored if needed.

<details>
<summary>Click to expand archived module documentation</summary>

### Presence Detection (`presence.py`)

Device registry and activity tracking - knows where you are.

```python
from backend.consciousness import (
    create_presence_manager,
    Device,
    DeviceType,
    UserStatus,
)

presence = await create_presence_manager(kv_store, bus)

# Register a device
device = Device(
    id="ipad-123",
    name="Zeus's iPad Pro",
    type=DeviceType.IPAD,
)
await presence.register_device(device)

# Check if user is active
if await presence.is_user_active():
    # Send via UI
else:
    # Send push notification

# Get current state
state = await presence.get_state()
# PresenceState(
#   status=UserStatus.ACTIVE,
#   primary_device=Device(...),
#   connected_devices=[...],
#   last_activity=datetime(...)
# )

# Track device connections
await presence.on_device_connected("ipad-123")
await presence.on_activity("ipad-123")
```

---

### 14. Shared Canvas (`canvas.py`)

Legacy's shared "screen" - what she's looking at, broadcast to all devices.

```python
from backend.consciousness import (
    create_canvas,
    CanvasItem,
    ContentType,
)

canvas = await create_canvas(kv_store, presence, bus)

# Show an email to user
await canvas.show(CanvasItem(
    type=ContentType.EMAIL,
    title="Re: Investor Meeting",
    content={"from": "john@sequoia.com", "body": "..."},
    context="The email you asked about",
    source="gmail",
))

# Helper methods
await canvas.show_email(
    email_id="123",
    subject="Re: Investor Meeting",
    sender="john@sequoia.com",
    body="Looking forward to our call...",
    provider="gmail",
)

await canvas.show_github_pr(
    owner="Zeus-Delacruz",
    repo="onemind-os",
    number=42,
    title="Add consciousness layer",
)

# Show on specific device
await canvas.show_on_device(item, device_id="ipad-123")
```

</details>

---

## Active Utility Modules

### Unified Notifier (`utils/notifier.py`)

Routes notifications to the right channel based on presence.

```python
from backend.consciousness import (
    create_notifier,
    NotificationPriority,
)

notifier = await create_notifier(presence, kv_store)

# Non-blocking alert
await notifier.alert("Task is overdue", priority=3)

# Alert with actions (ntfy buttons)
await notifier.alert_with_actions(
    "PR #42 needs review",
    actions=[
        {"label": "Approve", "action": "approve_pr", "data": {"pr": 42}},
        {"label": "Open", "url": "https://github.com/..."},
    ]
)

# Blocking input request (HITL)
response = await notifier.request_input(
    "Approve this PR?",
    options=["Approve", "Reject", "Later"],
)

# Confirmation
if await notifier.request_confirmation("Deploy to production?"):
    # Deploy...

# Convenience methods
await notifier.info("Build complete")
await notifier.success("Deployed successfully")
await notifier.warning("API rate limit at 80%")
await notifier.error("Build failed")
await notifier.urgent("Server down!")
```

---

### 16. Pattern Engine (`patterns.py`)

Detects patterns in event streams for proactive insights.

```python
from backend.consciousness import (
    create_pattern_engine,
    InsightPriority,
)

patterns = await create_pattern_engine(archiver, kv_store, bus)

# Patterns are detected automatically
# Query current insights
insights = await patterns.get_insights()
# [
#   Insight(message="You usually review PRs on Monday mornings - 3 pending"),
#   Insight(message="Task deadline in 2 hours, you haven't started"),
# ]

# Get as text
messages = await patterns.get_insights_text()
# ["You usually review PRs...", "Task deadline..."]

# Register callback for high-priority insights
def on_insight(insight):
    print(f"New insight: {insight.message}")

patterns.on_insight(on_insight)
```

> **Note:** Pattern Engine is archived in `.paused/consciousness-unused/patterns.py`

---

## Active Awareness Modules

### Temporal Awareness (`awareness/temporal.py`)

Time intelligence - deadlines, calendar, optimal timing.

```python
from backend.consciousness import (
    create_temporal_awareness,
    Deadline,
    TimeContext,
)

temporal = await create_temporal_awareness(kv_store, bus, notifier)

# Deadlines are tracked automatically from events
# Or add manually
await temporal.add_deadline(Deadline(
    id="task-123",
    title="Investor pitch",
    due=datetime(2024, 1, 20, 14, 0),
    source="todoist",
    source_id="123",
))

# Get upcoming deadlines
deadlines = await temporal.get_upcoming_deadlines(hours=24)
# Alerts sent automatically when deadlines approach

# Calendar integration
events = await temporal.get_upcoming_events(hours=24)
next_meeting = await temporal.get_next_meeting()

# Check for conflicts
conflicts = await temporal.check_conflicts(
    start=datetime.now(),
    end=datetime.now() + timedelta(hours=2),
)

# Find optimal focus time
slots = await temporal.suggest_focus_time(duration_minutes=60)
# [TimeSlot(start=..., end=..., quality=0.9), ...]

# Get time context
context = temporal.get_time_context()
# TimeContext.MORNING / AFTERNOON / EVENING / NIGHT / WEEKEND
```

---

### Deep Links (`utils/deep_links.py`)

Generate native app deep links for iOS/macOS/Android.

```python
from backend.consciousness import DeepLinks

# Email
DeepLinks.email("gmail", message_id="123")
# "googlegmail:///co?messageId=123"

DeepLinks.email("outlook", to="user@example.com", subject="Hello")
# "ms-outlook://compose?to=user%40example.com&subject=Hello"

# GitHub
DeepLinks.github("pr", "owner", "repo", 42)
# "github://pull-request/owner/repo/42"

DeepLinks.github("issue", "owner", "repo", 100)
# "github://issue/owner/repo/100"

# Todoist
DeepLinks.todoist_task("123456")
# "todoist://task?id=123456"

# Calendar
DeepLinks.calendar_apple(event_id="abc")
# "calshow:abc"

# Slack
DeepLinks.slack_channel("team123", "channel456")
# "slack://channel?team=team123&id=channel456"

# Home Assistant
DeepLinks.home_assistant(entity_id="light.living_room")
# "homeassistant://navigate/lovelace/default_view?entity_id=light.living_room"

# OneMind (future app)
DeepLinks.onemind_canvas("item-123")
# "onemind://canvas/item-123"
```

---

## Event Types

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

    # Home Assistant
    MOTION_DETECTED = "home.motion.detected"
    DOOR_OPENED = "home.door.opened"
    DEVICE_STATE_CHANGED = "home.device.changed"

    # System
    AGENT_STARTED = "system.agent.started"
    AGENT_ERROR = "system.agent.error"
```

---

## Awareness Levels (4 levels)

```python
class AwarenessLevel(IntEnum):
    DORMANT = 0      # Sleeping - only critical (priority >= 4)
    AWARE = 1        # Watching passively - high priority (priority >= 3)
    PRESENT = 2      # Engaged and ready - most events (priority >= 1) - DEFAULT
    OMNIPRESENT = 3  # Full presence - everything (priority >= 0)
```

---

## Priority Levels

| Priority | Description | Examples |
|----------|-------------|----------|
| 0 | Debug | Internal metrics, verbose logs |
| 1 | Low | FYI events, non-urgent updates |
| 2 | Normal | Standard events (default) |
| 3 | High | IoT alerts, @mentions |
| 4 | Urgent | Deadlines, conflicts |
| 5 | Critical | Security events, failures |

---

## Environment Variables

```env
# Required
NATS_URL=nats://nats:4222

# Optional (for archiver)
TIMESCALE_URL=postgresql://onemind:password@timescaledb:5432/events
```

---

## Quick Start

```python
import asyncio
from backend.consciousness import (
    create_event_bus,
    create_kv_store,
    create_memory_store,
    EventProcessor,
    Event,
    EventType,
)

async def main():
    # Connect to NATS
    bus = await create_event_bus("nats://localhost:4222")

    # Initialize stores
    kv = await create_kv_store(bus.js)
    memory = await create_memory_store(bus.js)

    # Start processor
    processor = EventProcessor(nats_url="nats://localhost:4222")
    await processor.start()

    # Publish an event
    await bus.publish(Event(
        type=EventType.TASK_CREATED,
        source="todoist",
        data={"title": "Review PR #42"},
        priority=2
    ))

asyncio.run(main())
```

---

## NATS Best Practices Implemented

This codebase follows NATS best practices from the official documentation:

### 1. Native Consumer Backoff

```python
# In core/bus.py - automatic exponential backoff
ConsumerConfig(
    durable_name="legacy-agent",
    max_deliver=5,
    backoff=[1000, 5000, 30000, 120000, 600000],  # 1s → 10min
)
# No custom retry logic needed!
```

### 2. Dead Letter Queue via Advisories

```python
# In processing/dlq.py - subscribe to native advisories
await nc.subscribe(
    "$JS.EVENT.ADVISORY.CONSUMER.MAX_DELIVERIES.>",
    cb=handle_dead_letter,
)
# No custom DLQ stream needed - NATS handles it!
```

### 3. Error Callbacks for Slow Consumers

```python
# In core/bus.py - handle slow consumers gracefully
async def error_cb(e):
    if isinstance(e, nats.errors.SlowConsumerError):
        logger.warning(f"Slow consumer detected: {e}")
    else:
        logger.error(f"NATS error: {e}")

await nats.connect(
    nats_url,
    error_cb=error_cb,
    disconnected_cb=disconnected_cb,
    reconnected_cb=reconnected_cb,
    max_reconnect_attempts=-1,  # Infinite reconnects
)
```

### 4. Native Service Discovery

```python
# In bridges/service.py - uses nats.micro framework
from nats.micro import add_service, ServiceConfig

svc = await add_service(nc, ServiceConfig(
    name="onemind-agent",
    version="2.0.0",
    queue_group="consciousness-workers",  # Load balancing!
))
# PING/STATS/INFO handled automatically!
```

### 5. Queue Groups for Scaling

```python
# Load balancing across multiple instances
await js.subscribe(
    "events.>",
    queue="consciousness-processors",  # Queue group
    cb=process_event,
)
# Run 3 instances - NATS distributes automatically
```

### 6. Drain Mode for Graceful Shutdown

```python
# Graceful shutdown in agentos_setup.py
if _event_bus:
    await _event_bus.disconnect()  # Drains before closing
```

---

## Startup Initialization Order

When the app starts (`backend/agno/agentos_setup.py`):

```
✓ Event Bus connected              # NATS JetStream
✓ KV Store initialized             # NATS KV buckets
✓ State Manager initialized        # User context tracking
✓ Temporal Awareness initialized   # Time/deadline tracking
✓ NATS Service registered ($SRV)   # Service discovery
✓ DLQ advisory handler started     # Dead letter handling
✓ Processor started                # Event filtering
✓ Archiver started (optional)      # TimescaleDB storage
```

---

## Why NATS?

NATS replaces multiple services with one:

| Need | Traditional | NATS |
|------|-------------|------|
| Event streaming | Kafka + Zookeeper | JetStream |
| State cache | Redis | KV Store |
| File storage | S3/MinIO | Object Store |
| Service discovery | Consul | $SRV Protocol |
| Retry logic | Custom code | Native backoff |
| Dead letter queue | Custom stream | Native advisories |

**Benefits:**
- **Single binary** (~10MB)
- **Sub-millisecond latency** (~100μs)
- **No operational complexity** (no Zookeeper, no Erlang)
- **Built-in resilience** (auto-reconnect, replay, backoff)
