---
type: agent-heartbeat
agent: sentinel
status: planned
updated: 2026-05-10
---
# HEARTBEAT.md — Sentinel 👁️

## Cadence

| Schedule | Time (ET) | What |
| -------- | --------- | ---- |
| Every 5min | 24/7 | omos-core reachability check |
| Every 15min | 24/7 | Self-health check — memory, disk, process |
| Daily 6am | ET | Full gateway status report to Legacy |

## Escalation Rules

| Trigger | Action | Level |
| ------- | ------ | ----- |
| omos-core unreachable >2min | Alert Zeus + Legacy immediately | L1 |
| Disk >85% on omos-sentinel | Alert Grid | L2 |
| Gateway process crash | Restart + alert Legacy | L2 |
| omos-sentinel unreachable | Failsafe — depends on external monitor | L0 |
