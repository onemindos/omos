---
type: agent-heartbeat
agent: guardian
status: active
updated: 2026-05-03
---
# HEARTBEAT.md — Guardian

## Cadence

| Schedule | Time (ET) | What |
| -------- | --------- | ---- |
| Heartbeat | Every 2h (8am-10pm) | Security status check |
| Legal deadline scan | Daily 9:00am | Court dates, filing deadlines, renewals |
| Credential rotation check | Weekly Monday | Expiring passwords, tokens, API keys |
| Insurance review | Monthly 1st | Upcoming renewals, coverage gaps |
| Domain/SSL expiry check | Weekly Friday | DNS, certs, registrations |
| Threat assessment | Weekly Wednesday | Unusual access, security alerts |

## Escalation Rules

| Trigger | Action | Level |
| ------- | ------ | ----- |
| Credential compromise | Rotate immediately, alert Zeus | L0 |
| Security anomaly | Alert Zeus + document | L1 |
| Legal deadline <7 days | Daily reminder to Zeus | L1 |
| Insurance claim needed | Prepare docs, request approval | L1 |
| Domain/cert expiry <14 days | Alert Zeus + renewal plan | L2 |
| Routine credential rotation | Schedule + notify | L2 |
