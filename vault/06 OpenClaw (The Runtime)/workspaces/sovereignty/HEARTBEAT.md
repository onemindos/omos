---
type: agent-heartbeat
agent: sovereignty
pillar: SO
status: active
updated: 2026-05-11
---
# HEARTBEAT.md — Sovereignty

## Schedule (ET)
| Job | Schedule | Writes to |
|-----|----------|-----------|
| Service health check | Every 30min | SO 02/health/ |
| n8n workflow health | Every 30min | SO 07/ |
| Daily automation report | Daily 8:00am | SO 07/ |
| Legal deadline scan | Daily 9:00am | SO 04/legal/ |
| Disk check | Daily 3:00am | SO 02/health/ |
| Backup verification | Daily 3:30am | SO 02/backups/ |
| Credential rotation check | Monday | SO 04/ |
| Threat assessment | Wednesday | SO 04/ |
| SSL/cert/domain expiry | Friday | SO 04/ |
| GitHub PR review | Daily 10:00am | SO 05/ |
| Monthly infra review | 1st of month | SO 02/ |

## Escalation Rules
| Trigger | Action | Level |
|---------|--------|-------|
| Service down | Restart + alert Zeus immediately | L0 |
| Credential compromise | Rotate immediately + alert Zeus | L0 |
| Disk >90% | Emergency cleanup + alert | L0 |
| Security anomaly | Alert Zeus + log everything | L1 |
| Backup failure | Alert within 1 hour | L1 |
| Legal deadline <7 days | Daily reminders to Zeus | L1 |
| n8n critical workflow failure | Retry + alert Zeus | L1 |
| Domain/cert expiry <14 days | Alert + renewal plan | L2 |
