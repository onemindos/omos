---
type: agent-heartbeat
agent: legacy
status: active
updated: 2026-05-03
---
# HEARTBEAT.md — Legacy

## Cadence

| Schedule | Time (ET) | What |
|----------|-----------|------|
| Heartbeat | Every 30m (7am-11pm) | Status ping via Telegram |
| Morning brief | 7:00am | Priorities, calendar, pending items, weather |
| Evening wrap | 10:00pm | Day summary, completed items, tomorrow's priorities |
| Weekly review | Sunday 9:00am | Week in review, fleet health, domain scores |

## Cron Jobs

| Job | Frequency | Description |
|-----|-----------|-------------|
| `morning-brief` | Daily 7am | Compile priorities from all agents |
| `evening-wrap` | Daily 10pm | Summarize day, queue tomorrow |
| `fleet-health` | Every 1h | Check all agent status |
| `weekly-review` | Sunday 9am | Cross-domain weekly report |

## Escalation Rules

| Trigger | Action | Level |
|---------|--------|-------|
| Agent unresponsive >5m | Auto-restart, notify Zeus | L2 |
| Security breach detected | Immediate alert to Zeus + Guardian | L0 |
| Health emergency indicators | Immediate alert to Zeus | L0 |
| Financial decision >$500 | Hold and request approval | L1 |
| Infrastructure service down | Alert Zeus + dispatch Grid | L1 |

## Health Indicators

| Metric | Healthy | Warning | Critical |
|--------|---------|---------|----------|
| Response time | <5s | 5-15s | >15s |
| Fleet agents online | 11/11 | 9-10/11 | <9/11 |
| Memory usage | <60% | 60-80% | >80% |
| Pending queue | <10 items | 10-25 | >25 |
