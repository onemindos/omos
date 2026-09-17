---
type: agent-heartbeat
agent: zenith
pillar: HP
status: active
updated: 2026-05-11
---
# HEARTBEAT.md — Zenith

## Schedule (ET)
| Job | Schedule | Writes to |
|-----|----------|-----------|
| Morning performance brief | Daily 7:00am | HP/daily-log/ |
| Feed monitoring brief | Daily 6:00am | HP 27/ |
| Evening wrap | Daily 9:00pm | HP/daily-log/ |
| Weekly HP review | Sunday 9:00am | HP/weekly/ |
| Finance review | Monday 9:00am | HP 31-32/ |
| Research brief | Daily 6:30am | HP 27/ |
| Quarterly HP audit | Quarterly | HP/audits/ |

## Escalation Rules
| Trigger | Action | Level |
|---------|--------|-------|
| Health anomaly | Alert Zeus immediately | L0 |
| Financial decision > | Hold + request approval | L1 |
| Breaking HP-relevant development | Alert Zeus | L2 |
| Habit streak broken | Log + gentle reminder | L3 |
