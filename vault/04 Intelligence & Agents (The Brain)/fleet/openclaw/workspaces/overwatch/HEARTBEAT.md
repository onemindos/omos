---
type: agent-heartbeat
agent: overwatch
status: active
updated: 2026-05-03
---
# HEARTBEAT.md — Overwatch

## Cadence

| Schedule | Time (ET) | What |
| -------- | --------- | ---- |
| Heartbeat | Every 1h (7am-10pm) | Automation status check |
| Workflow health | Every 30m | n8n execution status, failed runs |
| Daily automation report | 8:00am | Overnight runs, failures, metrics |
| Weather conditions | Every 3h | Field conditions for future ops |
| Weekly automation review | Saturday 9:00am | Performance, optimization opportunities |

## Escalation Rules

| Trigger | Action | Level |
| ------- | ------ | ----- |
| Workflow failure — critical path | Retry + alert Zeus | L1 |
| Workflow failure — non-critical | Retry + log | L2 |
| Multiple workflows failing | Alert Zeus + Grid | L1 |
| Field conditions dangerous (future) | Cancel mission + alert | L0 |
| Hardware malfunction (future) | Ground all ops + alert | L0 |
| New automation request | Queue for design review | L2 |
