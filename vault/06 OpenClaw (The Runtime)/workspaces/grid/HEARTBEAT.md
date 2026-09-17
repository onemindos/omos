---
type: agent-heartbeat
agent: grid
status: active
updated: 2026-05-03
---
# HEARTBEAT.md — Grid

## Cadence

| Schedule | Time (ET) | What |
| -------- | --------- | ---- |
| Service health check | Every 30m (24/7) | All systemd services, Docker containers |
| Disk usage check | Daily 3:00am | VPS storage, log rotation needed |
| Backup verification | Daily 3:30am | DB backup integrity, vault backup |
| SSL/cert expiry | Weekly Friday | Tailscale certs, domain certs |
| GitHub PR review | Daily 10:00am | Open PRs, stale branches |
| NATS stream health | Every 1h | Consumer lag, stream size |
| Dependency audit | Weekly Monday | Security patches, version updates |

## Escalation Rules

| Trigger | Action | Level |
| ------- | ------ | ----- |
| Service down | Restart + alert Zeus immediately | L0 |
| Disk >90% | Emergency cleanup + alert | L0 |
| Disk >80% | Alert Zeus with cleanup plan | L1 |
| Backup failure | Alert within 1 hour | L1 |
| Database connection issues | Diagnose + alert | L1 |
| Stale PR >7 days | Notify Zeus | L2 |
| Dependency update available | Queue for next maintenance window | L3 |
