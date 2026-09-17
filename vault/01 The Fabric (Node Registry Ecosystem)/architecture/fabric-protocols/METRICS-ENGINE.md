---
id: metrics-engine-20260430
title: "Metrics Engine"
type: protocol
status: active
created: 2026-04-30
modified: 2026-04-30
owner: zeus
tags: [type/protocol, domain/so, system/metrics, system/agents]
ai_summary: "Effectiveness measurement system — tracks agent velocity, task hit rate, autonomy levels, and system health so Zeus can see how well the fabric is performing."
---

# Metrics Engine

> If you can't measure it, you can't improve it. Every agent proves its value through numbers.

---

## Purpose

The Metrics Engine answers ONE question for Zeus:

**"Is the system working? Should I trust it more, or intervene?"**

It measures:
- **Are tasks getting done?** (velocity & completion rate)
- **Are they done on time?** (hit rate & SLA compliance)
- **Can agents work alone?** (autonomy score)
- **Is the system healthy?** (uptime, error rate, heartbeat consistency)
- **Is Legacy (Zeus) being freed up?** (delegation ratio)

---

## Core Metrics

### 1. Task Velocity

How fast work moves through the system.

```yaml
velocity:
  tasks_completed_today: 14
  tasks_completed_this_week: 67
  avg_completion_time_hours: 3.2    # from creation to done
  throughput_per_agent:
    legacy: 5
    spartan: 3
    oracle: 2
    eden: 2
    empire: 2
  trend: "up"  # up | stable | down (vs last week)
```

### 2. Hit Rate (On-Time Delivery)

Percentage of tasks completed before their due date.

```yaml
hit_rate:
  overall: 0.89        # 89% on time
  by_priority:
    p1: 0.95           # P1s almost never missed
    p2: 0.88
    p3: 0.82
    p4: 0.75           # Low priority gets deferred most
  by_agent:
    legacy: 0.92
    spartan: 0.90
    oracle: 0.85
    eden: 0.88
    empire: 0.87
  missed_tasks_today: 2
  missed_reasons:
    - "todoist-123: Blocked by dependency (client reply needed)"
    - "todoist-456: Agent capacity exceeded, rescheduled"
```

### 3. Autonomy Score

How much work agents handle WITHOUT Zeus intervention.

```yaml
autonomy:
  overall_score: 0.73   # 73% of tasks handled fully autonomously
  levels:
    full_auto: 47       # Agent did everything, Zeus never touched it
    assisted: 12        # Agent did work, Zeus reviewed/approved
    manual: 5           # Zeus had to step in and do it
    escalated: 3        # Agent flagged, Zeus decided
  trend: "up"           # Moving toward more autonomy
  target: 0.85          # Goal: 85% full autonomy by end of Q2
```

### 4. System Health

Infrastructure reliability metrics.

```yaml
system_health:
  heartbeat_reliability: 0.99    # % of heartbeats that fired on time
  webhook_uptime: 0.97           # % of time webhook endpoint was reachable
  avg_response_time_ms: 340      # From event to agent acknowledgment
  crons_fired_on_time: 0.96      # % of self-scheduled crons that executed
  errors_today: 2
  last_error: "2026-04-30T09:12:00Z — Todoist rate limit (recovered)"
  agent_uptime:
    legacy: "99.8%"
    spartan: "99.2%"
    oracle: "98.5%"
    eden: "99.1%"
    empire: "99.5%"
```

### 5. Delegation Ratio

How much of Zeus's life is being handled by the system vs. manually.

```yaml
delegation:
  total_tasks_this_week: 80
  handled_by_agents: 67          # 84% delegated
  handled_by_zeus: 13            # 16% manual
  zeus_time_saved_hours: 12.5    # Estimated based on avg task time
  categories_most_delegated:
    - "system maintenance (100%)"
    - "fitness tracking (95%)"
    - "email triage (90%)"
    - "business ops (80%)"
  categories_least_delegated:
    - "creative decisions (20%)"
    - "relationship tasks (30%)"
    - "financial decisions (40%)"
```

---

## Metric Collection Points

| Where | What's Collected | Frequency |
|-------|-----------------|-----------|
| Heartbeat logs | Tasks found, actions taken, duration | Every 15 min |
| Webhook events | Event count, response time, errors | Real-time |
| Todoist sync | Completion times, overdue count | Every heartbeat |
| Cron registry | Scheduled vs. fired vs. missed | On each cron event |
| Agent self-report | Autonomy level per task, blockers | Per task completion |
| Zeus interaction log | Which tasks Zeus touched manually | Passive tracking |

---

## Reporting

### Daily Pulse (Auto-generated at 22:00)

A brief daily summary added as a Todoist task comment or Codex entry:

```markdown
## Daily Pulse — 2026-04-30

**Velocity:** 14 tasks completed (↑ from 12 yesterday)
**Hit Rate:** 89% on-time (1 P2 missed — dependency)
**Autonomy:** 73% full-auto | 3 escalations to Zeus
**Health:** All systems green | 1 rate limit (auto-recovered)
**Delegation:** Zeus touched 4 tasks manually today

### Highlights
- Empire closed Acme proposal (P1, on time)
- Spartan logged all workouts without prompting
- Oracle surfaced investment opportunity (queued for Zeus review)

### Issues
- Eden missed grocery order (API timeout with delivery service)
- 2 crons fired late by ~30 sec (VPS load spike at 14:00)
```

### Weekly Review (Auto-generated Sunday 20:00)

Comprehensive weekly report saved to Codex:

```yaml
location: "_system/metrics/weekly/2026-W18.md"
contains:
  - velocity_chart (tasks/day trend)
  - hit_rate_by_agent
  - autonomy_progression
  - top_5_achievements
  - top_3_failures_and_fixes
  - delegation_ratio_trend
  - recommendations_for_next_week
```

### Monthly Scorecard

```yaml
location: "_system/metrics/monthly/2026-04.md"
contains:
  - month_over_month_comparison
  - agent_ranking (by effectiveness)
  - system_cost_vs_time_saved
  - autonomy_growth_chart
  - decision: "Should any agent be deprecated or upgraded?"
```

---

## Effectiveness Thresholds

When metrics drop below thresholds, actions trigger automatically:

| Metric | Green | Yellow | Red | Red Action |
|--------|-------|--------|-----|------------|
| Hit rate | > 85% | 70-85% | < 70% | Alert Zeus + reduce agent scope |
| Heartbeat reliability | > 98% | 95-98% | < 95% | Auto-restart + alert |
| Autonomy score | > 70% | 50-70% | < 50% | Review agent skills/context |
| Webhook response time | < 500ms | 500ms-2s | > 2s | Scale infrastructure |
| Crons on-time | > 95% | 90-95% | < 90% | Investigate VPS performance |
| Agent uptime | > 99% | 97-99% | < 97% | Restart + post-mortem |

---

## Zeus Takeover Protocol

The metrics exist so Zeus can **trust and verify**. The takeover flow:

```text
GREEN (Autonomous Mode)
│  Zeus reviews weekly report, makes minor adjustments
│  Agents work independently
│
├─► YELLOW (Advisory Mode)
│   Metrics dipping — agents flag issues proactively
│   Zeus reviews daily pulse, provides course corrections
│   Agents request approval for high-impact tasks
│
└─► RED (Manual Mode)
    System is underperforming — Zeus takes direct control
    Agents pause autonomous execution
    All tasks require Zeus approval before execution
    Post-mortem required before returning to Green
```

---

## Storage

```text
_system/metrics/
├── daily/
│   ├── 2026-04-30.yaml
│   ├── 2026-04-29.yaml
│   └── ...
├── weekly/
│   ├── 2026-W18.md
│   └── ...
├── monthly/
│   ├── 2026-04.md
│   └── ...
├── heartbeat-logs/
│   └── (rolling 7-day window)
└── thresholds.yaml        ← configurable alert thresholds
```

---

## Integration with Codex × Multiply Step

This Metrics Engine IS the **× (Multiply)** step of the CODEX method:

```text
CODEX:  C → O → D → E → ×
                          ▲
                          │
              Metrics Engine lives here
              
× Multiply asks: "Did it work? What did we learn? How do we improve?"
The Metrics Engine answers with data, not guesses.
```


---

## Related

- [[00-24 SO (Sovereign Operations)/03 Command (Operations & Navigation) (The Controls)/03 Command|03 Command]]
