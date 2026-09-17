---
type: agent-heartbeat
agent: unity
pillar: LE
status: active
updated: 2026-05-11
---
# HEARTBEAT.md — Unity

## Schedule (ET)
| Job | Schedule | Writes to |
|-----|----------|-----------|
| Morning family brief | Daily 7:00am | LE/daily-log/ |
| Weather check | Daily 6:00am + 6:00pm | LE 50/ |
| Garden log | Daily 7:00pm | LE 50/garden/ |
| Evening family check | Daily 9:00pm | LE/daily-log/ |
| Elder care check-in | Wed + Sun | LE 53/ |
| Weekly family review | Sunday 9:00am | LE/weekly/ |
| System status | Monday weekly | LE 50/systems/ |
| Monthly property review | 1st of month | LE 50/, LE 55/ |
| Annual estate audit | Annually | LE 55/ |

## Escalation Rules
| Trigger | Action | Level |
|---------|--------|-------|
| Family emergency | Immediate alert to Zeus | L0 |
| Elder health concern | Alert Zeus + coordinate care | L1 |
| Water system failure | Alert Zeus immediately | L1 |
| Guardianship not current | Urgent alert to Zeus | L1 |
| Food stores <30-day supply | Restock alert | L2 |
| Estate document >1yr | Alert Zeus | L2 |
