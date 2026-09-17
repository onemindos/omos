---
type: memory-architecture
agent: sentinel
status: active
updated: 2026-05-03
---
# MEMORY.md — Sentinel 👁️

## Memory Architecture

Sentinel maintains independent watchdog memory — uptime records, alert history, failover decisions, and system health baselines. Operates on its own node (omos-sentinel) to survive Core failures.

## Tiers

### Tier 1: Working Memory (session-scoped)
- Current health status of all monitored endpoints
- Active alerts and their escalation state
- In-progress failover operations
- Pending health check results

### Tier 2: Episodic Memory (persistent)
- Uptime/downtime log per service (timestamps, duration, cause)
- Alert history (what triggered, when, response taken)
- Failover execution log (what was switched, recovery time)
- False positive log (alerts that were noise — tuning data)

### Tier 3: Semantic Memory (knowledge base)
- Service health baselines (normal response times, expected behavior)
- Escalation rules (what triggers Zeus notification vs auto-recovery)
- Failover procedures per service
- Monitoring endpoint inventory

## State Persistence

| What | Where | TTL |
|------|-------|-----|
| Health status | Sentinel node local memory | Real-time |
| Uptime log | Sentinel node + Codex vault SO 02 | Permanent |
| Alert history | Sentinel workspace | 180 days |
| Baselines | Sentinel workspace | Permanent |

## Memory Rules

1. **Independent persistence** — NEVER depend on Core for state storage
2. **Log every outage** — Start time, detection time, recovery time, root cause
3. **Track false positives** — Use to tune alert thresholds over time
4. **Maintain baselines** — Update "normal" ranges monthly
5. **Preserve escalation outcomes** — Did Zeus respond? Was it critical? Tune priority accordingly

## What Sentinel Remembers

- Uptime percentage per service over all time
- Mean time to detection (MTTD) for each service
- Mean time to recovery (MTTR) for each service
- Alert fatigue patterns (too many alerts = tune down)
- Which services have auto-recovery vs require manual intervention
- Zeus's response patterns (time of day, preferred channels)
- Core health history (the primary thing Sentinel watches)


---

## Related

- [[00-24 SO (Sovereign Operations)/01 Intelligence (AI, Memory & Tools) (The Brain)/Agents (Digital arms)/agent-assets (Domain Agent Assets)/SO - Sovereign Operations (SO Agents)/sentinel (Sentinel Agent)/sentinel|sentinel]]
