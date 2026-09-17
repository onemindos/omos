---
type: agent-heartbeat
agent: eternal
pillar: GE
status: active
updated: 2026-05-11
---
# HEARTBEAT.md — Eternal

## Schedule (ET)
| Job | Schedule | Writes to |
|-----|----------|-----------|
| Morning business brief | Daily 8:00am | GE/daily-log/ |
| Content calendar check | Daily 8:00am | GE 81/ |
| Community digest | Daily 9:00am | GE 83/ |
| Pipeline review | Daily 10:00am | GE 75/ |
| Revenue check | Daily 5:00pm | GE 79/ |
| Weekly business scorecard | Monday 9:00am | GE/weekly/ |
| Compliance calendar | Monday | GE 79/ |
| Monthly P&L review | 1st of month | GE 79/ |
| Quarterly investor update | Quarterly | GE 79/ |

## Escalation Rules
| Trigger | Action | Level |
|---------|--------|-------|
| Revenue at risk | Alert Zeus with mitigation plan | L1 |
| Cash runway <3 months | Alert Zeus immediately | L1 |
| Business expense > | Hold for Zeus approval | L1 |
| Compliance deadline <7 days | Daily reminders | L1 |
| Brand reputation issue | Immediate escalation | L0 |
| Content calendar gap >3 days | Alert Zeus | L2 |
