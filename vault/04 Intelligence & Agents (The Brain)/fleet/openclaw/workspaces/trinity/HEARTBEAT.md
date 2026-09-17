---
type: agent-heartbeat
agent: trinity
status: active
updated: 2026-05-03
---
# HEARTBEAT.md — Trinity

## Cadence

| Schedule | Time (ET) | What |
| -------- | --------- | ---- |
| Heartbeat | Every 2h (9am-10pm) | Community status, engagement check |
| Community digest | Daily 9:00am | New members, trending posts, engagement |
| Content schedule check | Daily 10:00am | What's due, what's overdue |
| Discord health | Every 2h | Unanswered questions, engagement |
| Member retention scan | Weekly Friday | At-risk members, churn signals |
| Email campaign status | Mon/Thu 11:00am | Open rates, click rates |
| Partnership follow-up | Weekly Wednesday | Outreach status, responses |

## Escalation Rules

| Trigger | Action | Level |
| ------- | ------ | ----- |
| Abusive community member | Warn once, escalate to Zeus if repeat | L1 |
| Member churn spike >10% | Alert Zeus with analysis | L1 |
| Brand reputation issue | Immediate alert to Zeus | L0 |
| Content publishing failure | Retry + alert if persistent | L2 |
| Partnership response received | Notify Zeus with summary | L2 |
| Missed scheduled content | Reschedule + notify Zeus | L2 |
