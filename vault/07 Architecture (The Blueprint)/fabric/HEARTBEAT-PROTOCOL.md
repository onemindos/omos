---
id: heartbeat-protocol-20260430
title: "Heartbeat Protocol"
type: protocol
status: active
created: 2026-04-30
modified: 2026-04-30
owner: zeus
tags: [type/protocol, domain/so, system/todoist, system/agents]
ai_summary: "15-minute polling protocol for agents to check Todoist, assess work state, and self-schedule cron jobs for coverage gaps."
---

# Heartbeat Protocol

> Every 15 minutes, agents wake up, check what's happening, and decide what to do next.

---

## Overview

The Heartbeat is the **baseline pulse** of the Productivity Fabric. It guarantees:
- No task sits unattended for more than 15 minutes
- Agents always have current awareness of the task landscape
- Gaps between heartbeats are covered by self-scheduled crons
- Failures in real-time webhooks are caught at the next heartbeat

---

## Heartbeat Cycle (Every 15 Minutes)

```text
HEARTBEAT TICK (every 15 min)
│
├─► 1. POLL TODOIST
│       • Fetch all tasks updated since last heartbeat
│       • Fetch tasks due within next 30 min
│       • Fetch tasks assigned to this agent's domain
│       • Check for new comments/instructions from Zeus
│
├─► 2. ASSESS STATE
│       • Compare current tasks vs. last known state
│       • Identify NEW tasks (added since last heartbeat)
│       • Identify CHANGED tasks (priority, due date, labels)
│       • Identify OVERDUE tasks (missed deadlines)
│       • Identify BLOCKED tasks (dependencies unmet)
│
├─► 3. LOAD CODEX CONTEXT
│       • Read domain-specific goals from Codex
│       • Check active priorities / weekly focus
│       • Load relevant skills for task types detected
│
├─► 4. DECIDE & ACT
│       • Execute immediately: tasks due NOW or overdue
│       • Schedule cron: tasks due in 1-4 hours (between heartbeats)
│       • Queue for next heartbeat: low priority, no deadline
│       • Escalate to Zeus: ambiguous, high-stakes, or out-of-scope
│
├─► 5. SELF-SCHEDULE
│       • Create cron jobs for tasks that need attention between heartbeats
│       • Example: "Deploy at 3:00 AM" → cron at 2:55 AM
│       • Example: "Send report by noon" → cron at 11:45 AM
│       • Register scheduled crons in `_system/crons/active.yaml`
│
└─► 6. LOG & REPORT
        • Log heartbeat result to metrics
        • Update `last_heartbeat` timestamp
        • If anomalies found → flag for Zeus
```

---

## Heartbeat Timing

| Parameter | Value | Rationale |
|-----------|-------|-----------|
| Interval | 15 minutes | Balanced response time vs. API rate limits |
| Offset | Staggered per agent | Prevents all agents hitting Todoist simultaneously |
| Timeout | 30 seconds | If heartbeat takes longer, abort and retry next cycle |
| Retry | Skip to next cycle | Never double-poll; wait for next heartbeat |
| Quiet hours | Configurable per agent | Agents can reduce to 30-min heartbeat overnight |

### Agent Stagger Schedule

```yaml
heartbeat_schedule:
  legacy:   "*/15 * * * *"       # :00, :15, :30, :45
  spartan:  "2-59/15 * * * *"   # :02, :17, :32, :47
  oracle:   "5-59/15 * * * *"   # :05, :20, :35, :50
  eden:     "8-59/15 * * * *"   # :08, :23, :38, :53
  empire:   "11-59/15 * * * *"  # :11, :26, :41, :56
```

---

## What Agents Check in Todoist

### Filters Per Agent (by Domain)

| Agent | Todoist Filter | Domain |
|-------|---------------|--------|
| **Legacy** | All projects (meta-agent, oversees everything) | SO (00-24) |
| **Spartan** | Projects tagged `@health`, `@fitness`, `@body` | HP (25-49) |
| **Oracle** | Projects tagged `@mind`, `@money`, `@learning` | HP (25-49) |
| **Eden** | Projects tagged `@family`, `@home`, `@community` | LE (50-74) |
| **Empire** | Projects tagged `@business`, `@career`, `@brand` | GE (75-99) |

### Task Priority Mapping

```yaml
todoist_priority_to_action:
  p1: "Execute immediately — do not wait for next heartbeat"
  p2: "Execute this heartbeat cycle if capacity allows"
  p3: "Schedule cron for optimal time, or next heartbeat"
  p4: "Queue — pick up when nothing higher is pending"
```

---

## Self-Scheduling (Cron Autonomy)

When an agent sees a task that needs attention **between heartbeats**, it creates a cron job:

### Rules for Self-Scheduling

1. **Only schedule within your domain** — never touch another agent's work
2. **Buffer time** — Schedule 5 min before actual deadline (not at deadline)
3. **Register all crons** — Write to `_system/crons/active.yaml` so other agents see it
4. **Clean up after** — Remove cron entry once task completes
5. **Max 10 active crons per agent** — Prevents runaway scheduling
6. **Zeus override** — Zeus can cancel any cron via Todoist comment `@cancel-cron`

### Cron Registration Format

```yaml
# _system/crons/active.yaml
crons:
  - id: cron-20260430-legacy-001
    agent: legacy
    task_id: "todoist-8234567890"
    task_title: "Deploy staging build"
    scheduled_at: "2026-04-30T14:55:00Z"
    reason: "Task due at 15:00, scheduling 5 min early"
    created_at: "2026-04-30T14:30:00Z"
    status: pending  # pending | running | completed | cancelled

  - id: cron-20260430-empire-001
    agent: empire
    task_id: "todoist-8234567891"
    task_title: "Send invoice to client"
    scheduled_at: "2026-04-30T08:55:00Z"
    reason: "Task due at 09:00, morning delivery preferred"
    created_at: "2026-04-30T02:15:00Z"
    status: completed
```

---

## Heartbeat Data Model

Each heartbeat produces a log entry:

```yaml
# Heartbeat log entry
heartbeat:
  agent: legacy
  timestamp: "2026-04-30T14:00:00Z"
  duration_ms: 2340
  tasks_found: 12
  tasks_new: 2
  tasks_changed: 1
  tasks_overdue: 0
  actions_taken:
    - executed: "todoist-123" # did it now
    - scheduled_cron: "todoist-456" # will do later
    - queued: "todoist-789" # next heartbeat
  crons_created: 1
  crons_completed: 0
  anomalies: []
  next_heartbeat: "2026-04-30T14:15:00Z"
```

---

## Failure Modes & Recovery

| Failure | Detection | Recovery |
|---------|-----------|----------|
| Todoist API down | HTTP 5xx / timeout | Skip this heartbeat, retry next cycle, flag Zeus after 3 consecutive |
| Agent crash during heartbeat | No heartbeat log after expected time | Supervisor restarts agent, next heartbeat catches up |
| Duplicate task execution | Two agents claim same task | Task-level locking via Todoist labels (`@claimed-by:legacy`) |
| Cron missed | `scheduled_at` passed but status still `pending` | Next heartbeat detects stale cron, executes immediately |
| Rate limit hit | HTTP 429 from Todoist | Back off, reduce to essential queries only |

---

## Implementation Stack

| Component | Technology | Notes |
|-----------|-----------|-------|
| Scheduler | Cron (VPS) + N8N | N8N for visual flow, cron for reliability |
| Todoist API | REST API v1 / Sync API v9 | Sync API for efficient polling |
| State store | `_system/heartbeat/` in Codex | YAML logs, human-readable |
| Alerting | Todoist task + webhook to Zeus | High-priority alert = new P1 task |
| Cron registry | `_system/crons/active.yaml` | Shared state, all agents read/write |


---

## Related

- [[00-24 SO (Sovereign Operations)/03 Command (Operations & Navigation) (The Controls)/03 Command|03 Command]]
