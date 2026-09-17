# Proactive Intelligence in OneMind OS

> How State and Events Create an AI That Acts Without Being Asked

## Overview

OneMind OS is fundamentally different from chat-based AI assistants. Rather than waiting for user queries, it **continuously receives events, maintains state, and acts proactively**.

```
Traditional (Reactive):  User asks → System searches → Returns answer
OneMind OS (Proactive):  Events flow → Context builds → Agent acts autonomously
```

---

## The Three-Layer Architecture

### Layer 1: The Nervous System (NATS EventBus)

**File:** `backend/consciousness/core/bus.py`

All external events flow through NATS JetStream:

```
┌────────────────────────────────────────────────────────────┐
│ EXTERNAL SOURCES                                           │
│ Todoist • GitHub • Home Assistant • Calendar • Biometrics │
└──────────────────────────┬─────────────────────────────────┘
                           │ webhooks / polling
                           ▼
              ┌────────────────────────┐
              │    NATS JETSTREAM      │
              │    (200-400μs latency) │
              │                        │
              │  events.>              │
              │  action.>              │
              │  legacy.inbox          │
              └────────────────────────┘
```

**Key Properties:**
- **Push Model**: Events delivered instantly (no polling)
- **Persistent**: JetStream survives restarts (7-day retention)
- **Priority Routing**: Urgent events bypass normal queues
- **Sub-millisecond**: 200-400μs end-to-end latency

### Layer 2: The Gatekeeper (LightweightProcessor)

**File:** `backend/consciousness/processing/processor.py`

Filters events based on awareness level:

| Level | Threshold | Use Case |
|-------|-----------|----------|
| `DORMANT` (0) | Priority ≥ 4 | Sleeping - only critical alerts |
| `AWARE` (1) | Priority ≥ 3 | Watching - high priority + approvals |
| `PRESENT` (2) | Priority ≥ 1 | Engaged - most events (DEFAULT) |
| `OMNIPRESENT` (3) | Priority ≥ 0 | Full presence - everything |

**Intelligent Filtering:**
- Noise filtering (typing indicators, heartbeats)
- Time-aware thresholds (night vs day)
- Location-aware (per-room awareness)
- Batch grouping (50 motion events → "Motion in 5 areas")

### Layer 3: The Context Injector (AwarenessProvider)

**File:** `backend/consciousness/aware_runner.py`

Automatically injects real-time context into agent prompts:

```python
# Without awareness - agent is blind
response = await agent.arun("What should I do?")
# Agent: "I don't have enough information"

# With awareness - agent has full context
response = await run_with_awareness(agent, "What should I do?")
# Agent: "You've been coding 45min, deadline in 2.5h. Take a break."
```

**Injected Context:**
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

---

## State Management

### NATS KV Store (Replaces Redis)

**File:** `backend/consciousness/core/kv_store.py`

State is stored in NATS KV buckets with **watch-based updates**:

```python
# Set state
await kv.set_awareness_level(AwarenessLevel.PRESENT)

# Watch for changes (no polling!)
async for key, value in kv.watch_state("awareness.*"):
    print(f"State changed: {key} = {value}")
```

### State Types

| Bucket | Contents |
|--------|----------|
| `awareness-level` | Current awareness mode (0-3) |
| `user-state` | Location, activity, health metrics |
| `filter-config` | Processing thresholds |
| `last-seen-seq` | Event replay position |

---

## Event Flow: End-to-End

```
┌─────────────────────────────────────────────────────────────────────┐
│ EXTERNAL WORLD                                                      │
│ (Todoist, GitHub, Home Assistant, Calendar, Oura)                  │
└────────────────────────────┬────────────────────────────────────────┘
                             │
                             ▼ webhooks / APIs
┌─────────────────────────────────────────────────────────────────────┐
│ NATS JETSTREAM                                                      │
│ ┌─────────────┐  ┌─────────────┐  ┌────────────────┐               │
│ │  events.>   │  │  action.>   │  │  legacy.inbox  │               │
│ └─────────────┘  └─────────────┘  └────────────────┘               │
└────────────────────────────┬────────────────────────────────────────┘
                             │
         ┌───────────────────┼───────────────────┐
         │                   │                   │
         ▼                   ▼                   ▼
   ┌───────────┐      ┌───────────┐      ┌───────────┐
   │ PROCESSOR │      │ CORRELATOR│      │  ARCHIVER │
   │ (filter)  │      │ (link)    │      │ (store)   │
   └─────┬─────┘      └─────┬─────┘      └───────────┘
         │                  │
         └────────┬─────────┘
                  │
                  ▼
         ┌────────────────┐
         │ ACTIVITY       │
         │ SERVICE        │
         │                │
         │ Creates items  │
         │ Broadcasts WS  │
         └────────┬───────┘
                  │
     ┌────────────┼────────────┐
     │            │            │
     ▼            ▼            ▼
┌─────────┐ ┌──────────┐ ┌──────────┐
│ Frontend│ │  Agent   │ │  Mobile  │
│ (UI)    │ │ (action) │ │  (push)  │
└─────────┘ └──────────┘ └──────────┘
```

---

## Activity Feed: Events → Actions

**Files:** `backend/consciousness/activity/`

The activity feed converts raw events into actionable items:

### Event Type Mapping

```python
EVENT_TYPE_MAPPING = {
    "github.pull_request.review_requested" → APPROVAL (HIGH),
    "todoist.item.added" → TASK (NORMAL),
    "homeassistant.alarm_triggered" → ALERT (CRITICAL),
    "hitl.approval_required" → APPROVAL (HIGH),
}
```

### Auto-Assignment

Based on item type and awareness mode:

| Type | Assignment |
|------|------------|
| `APPROVAL` | Always HUMAN |
| `ALERT` | Always HUMAN |
| `TASK` (low priority) + OMNIPRESENT | Agent can handle |
| `SUGGESTION` + PRESENT+ | Agent can auto-handle |

---

## Temporal Scheduling

**File:** `backend/consciousness/temporal/scheduler.py`

The temporal layer handles **when** actions should occur:

```python
# Schedule recurring briefing
await scheduler.schedule(ScheduledEventCreate(
    type=ScheduleEventType.TRIGGER,
    scheduled_at=tomorrow_9am,
    title="Morning briefing",
    action_type="workflow",
    action_target="morning-briefing",
    is_recurring=True,
    recurrence_rule="FREQ=DAILY;BYHOUR=9",
))
```

**Capabilities:**
- One-time triggers
- Recurring schedules (cron-like)
- Deadline monitoring
- Focus window suggestions

---

## What Makes This Proactive

### 1. Push Not Pull
Events delivered instantly via NATS. No "check for updates" loops.

### 2. Always Aware
Agents receive context injection automatically. They don't ask for context.

### 3. Awareness Modes
Different filtering based on user state. Night mode filters more aggressively.

### 4. Bidirectional Loop
Agent actions are logged back to the event bus, creating closed feedback loops.

### 5. Temporal Intelligence
Scheduler handles future events and suggests optimal timing for actions.

### 6. Correlation
Related events are linked across services (PR + CI + Deploy = one narrative).

---

## Comparison: OneMind OS vs ChatGPT

| Aspect | ChatGPT | OneMind OS |
|--------|---------|------------|
| **Awareness** | Only user's chat input | Real-time events + state |
| **Initiative** | Waits for messages | Acts on schedule/events |
| **Context** | Per-conversation | Continuous, cross-service |
| **Actions** | Text output only | Modifies external systems |
| **Memory** | Session-based | Persistent across sessions |
| **Patterns** | None | Correlates events over time |

---

## Quick Integration

```python
from backend.consciousness.aware_runner import run_with_awareness
from backend.agno.agents.factory import get_agent

# Get any agent
agent = get_agent("legacy")

# Run with automatic context injection
response = await run_with_awareness(
    agent,
    user_message="What should I focus on?",
    session_id="session_123",
    include_events=True,      # Recent NATS events
    include_state=True,       # Current user state
    include_deadlines=True,   # Upcoming deadlines
    event_count=10,           # Last 10 events
)
```

---

## Key Files Reference

### Event Infrastructure
| File | Purpose |
|------|---------|
| `core/bus.py` | NATS JetStream event streaming |
| `core/kv_store.py` | State management (replaces Redis) |
| `core/object_store.py` | Memory and file storage |

### Processing Pipeline
| File | Purpose |
|------|---------|
| `processing/processor.py` | Filter by awareness level |
| `processing/correlator.py` | Link related events |
| `processing/archiver.py` | TimescaleDB storage |

### Agent Awareness
| File | Purpose |
|------|---------|
| `awareness/state.py` | User state tracking |
| `awareness/temporal.py` | Time & deadline awareness |
| `aware_runner.py` | Agent wrapper with injection |

### Activity System
| File | Purpose |
|------|---------|
| `activity/models.py` | Pydantic models |
| `activity/service.py` | CRUD + broadcasting |
| `activity/feed.py` | Event→ActivityItem conversion |
| `activity/api.py` | FastAPI endpoints |

### Temporal
| File | Purpose |
|------|---------|
| `temporal/scheduler.py` | Schedule future events |
| `temporal/service.py` | Time awareness |
| `temporal/context.py` | Temporal context injection |

---

## Summary

OneMind OS achieves proactive intelligence through:

1. **NATS EventBus** - Real-time event delivery (sub-millisecond)
2. **LightweightProcessor** - Intelligent filtering by awareness
3. **AwarenessProvider** - Automatic context injection to agents
4. **Temporal Scheduler** - When to act, not just how
5. **Activity Feed** - Events become actionable items
6. **NATS KV State** - Watch-based updates, no polling

The key insight: **Agents don't wait for questions. They receive context continuously and act based on the situation.**
