---
type: agent-heartbeat
agent: haven
status: active
updated: 2026-05-03
---
# HEARTBEAT.md — Haven

## Cadence

| Schedule | Time (ET) | What |
| -------- | --------- | ---- |
| Heartbeat | Every 2h (7am-10pm) | Family and home status |
| Morning family brief | 7:00am | Today's calendar, appointments, reminders |
| Evening check | 9:00pm | Day summary, tomorrow preview |
| Weekly family review | Sunday 9:00am | Week ahead, birthdays, milestones |
| Property maintenance | Monthly 1st | Upcoming maintenance, seasonal tasks |
| Elder care check-in | Weekly Wed + Sun | Parent health, needs, visits |
| Heritage project | Monthly 15th | Photo archival, family history progress |

## Escalation Rules

| Trigger | Action | Level |
| ------- | ------ | ----- |
| Family emergency | Immediate alert to Zeus | L0 |
| Property damage | Alert + emergency repair coordination | L1 |
| Smart home security alert | Verify + alert Zeus | L1 |
| Elder health concern | Alert Zeus + coordinate care | L1 |
| Missed family appointment | Alert + reschedule | L2 |
| Maintenance overdue >30 days | Escalation reminder | L2 |
| Birthday/anniversary this week | Reminder to Zeus | L3 |
