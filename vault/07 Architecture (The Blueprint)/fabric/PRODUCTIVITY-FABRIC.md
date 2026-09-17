---
id: productivity-fabric-20260430
title: "The Productivity Fabric"
type: architecture
status: active
created: 2026-04-30
modified: 2026-04-30
owner: zeus
tags: [type/architecture, domain/so, status/active, system/todoist, system/agents]
ai_summary: "Master architecture for the Todoist + Codex + Agent productivity fabric — heartbeat polling, self-scheduling, real-time webhooks, sub-agent delegation, and effectiveness metrics."
---

# The Productivity Fabric

> **"The nervous system that connects Zeus's brain (Codex) to his hands (Agents) through his task engine (Todoist)."**
>
> Status: Active | Version: 1.0 | Created: 2026-04-30

---

## Executive Summary

The Productivity Fabric is the **operational layer** that binds three systems into one coherent flow:

| Component | Role | Metaphor |
|-----------|------|----------|
| **Todoist** | Task Engine | The bloodstream — carries work items |
| **Codex** | Knowledge Brain | The cerebral cortex — context & memory |
| **Agents** | Execution Layer | The hands — do the actual work |

Without the fabric, these are three disconnected tools. With it, they become **ONE MIND** operating in real-time.

---

## Architecture Overview

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                         THE PRODUCTIVITY FABRIC                              │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  ┌──────────────┐                                    ┌──────────────┐      │
│  │   TODOIST    │◄──── Webhooks (real-time) ────────►│   WEBHOOK    │      │
│  │  Task Engine │                                    │   ROUTER     │      │
│  │              │◄──── Heartbeat (15 min) ──────────►│              │      │
│  └──────┬───────┘                                    └──────┬───────┘      │
│         │                                                    │              │
│         │  task sync                              route to   │              │
│         ▼                                         agents     ▼              │
│  ┌──────────────┐         ┌──────────────┐      ┌──────────────────┐      │
│  │    CODEX     │◄───────►│   SCHEDULER  │◄────►│  AGENT DISPATCH  │      │
│  │   (Brain)    │ context │  (Cron + On- │      │                  │      │
│  │              │─────────►│   Demand)    │      │  ┌────┐ ┌────┐  │      │
│  │  • Domain    │         └──────────────┘      │  │ L  │ │ S  │  │      │
│  │    context   │                                │  │ e  │ │ p  │  │      │
│  │  • Skills    │    metrics / logs              │  │ g  │ │ a  │  │      │
│  │  • Memory    │◄──────────────────────────────│  │ a  │ │ r  │  │      │
│  │  • Metrics   │                                │  │ c  │ │ t  │  │      │
│  └──────────────┘                                │  │ y  │ │ a  │  │      │
│                                                   │  │    │ │ n  │  │      │
│                                                   │  ├────┤ ├────┤  │      │
│  ┌──────────────┐                                │  │ O  │ │ Ed │  │      │
│  │   METRICS    │◄─── effectiveness data ───────│  │ r  │ │ en │  │      │
│  │   ENGINE     │                                │  │ a  │ │    │  │      │
│  │              │                                │  │ c  │ ├────┤  │      │
│  │  • Velocity  │                                │  │ l  │ │ Em │  │      │
│  │  • Hit rate  │                                │  │ e  │ │ p  │  │      │
│  │  • Autonomy  │                                │  └────┘ └────┘  │      │
│  └──────────────┘                                └──────────────────┘      │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## The Three Flows

### Flow 1: Heartbeat (Scheduled Polling)
Every 15 minutes, agents check Todoist and the Codex for work.

### Flow 2: Self-Scheduling (Cron Autonomy)
Agents create their own cron jobs to handle tasks at optimal times between heartbeats.

### Flow 3: Real-Time Webhooks (Event-Driven)
Immediate response to external events — Todoist changes, system alerts, user commands.

---

## How Zeus + Agents Work Together

```
ZEUS (Human)                    FABRIC                         AGENTS
─────────────                   ──────                         ──────
                                                               
Add task to Todoist ──────────► Webhook fires ───────────────► Agent picks up task
                                                               Agent loads Codex context
                                                               Agent executes
                                                               Agent marks complete in Todoist
                                                               Agent logs metrics to Codex
                                                               
Review metrics in Codex ◄────── Metrics aggregated ◄─────────── Agents report effectiveness
                                                               
Override / redirect ──────────► Priority change ─────────────► Agent reschedules
                                                               
Sleep / Away ─────────────────► Heartbeat continues ─────────► Agents work autonomously
                                                               Agents self-schedule crons
                                                               Agents pick up new tasks
                                                               Zeus reviews on return
```

---

## Integration with CODEX Method

| CODEX Step | Fabric Role |
|------------|-------------|
| **C — Capture** | Webhook router ingests new signals → creates Todoist tasks |
| **O — Organize** | Heartbeat classifies + routes tasks to correct domain/agent |
| **D — Direct** | Scheduler prioritizes based on Codex goals + deadlines |
| **E — Execute** | Agent Dispatch assigns sub-agents with skills |
| **X — × Multiply** | Metrics Engine measures velocity, logs learnings to Codex |

---

## File Index (This Fabric)

| File | Purpose |
|------|---------|
| `PRODUCTIVITY-FABRIC.md` | This master architecture doc |
| `fabric/HEARTBEAT-PROTOCOL.md` | 15-min polling spec |
| `fabric/SELF-SCHEDULING.md` | Agent cron autonomy rules |
| `fabric/WEBHOOK-ROUTER.md` | Real-time event handling |
| `fabric/METRICS-ENGINE.md` | Effectiveness measurement |
| `fabric/SUB-AGENT-DISPATCH.md` | Task delegation & skills |

---

## Principles

1. **Zeus is sovereign** — Agents never act outside their scope without permission
2. **Todoist is the single task queue** — No shadow task lists
3. **Codex is the single source of context** — Agents always load domain context before acting
4. **Metrics prove value** — If an agent can't show effectiveness, it gets deprecated
5. **Graceful degradation** — If webhooks fail, heartbeat catches it; if heartbeat fails, cron catches it
6. **No silent failures** — Every missed task gets flagged in the next heartbeat
7. **Sub-agents for complexity** — Main agents delegate specialized work to sub-agents with skills


---

## Related

- [[00-24 SO (Sovereign Operations)/03 Command (Operations & Navigation) (The Controls)/03 Command|03 Command]]
