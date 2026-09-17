---
type: memory-architecture
agent: overwatch
status: active
updated: 2026-05-03
---
# MEMORY.md — Overwatch 🤖

## Memory Architecture

Overwatch maintains field operations memory — automation state, sensor data, drone missions, and physical system status.

## Tiers

### Tier 1: Working Memory (session-scoped)
- Active automation workflows running
- Current sensor readings and alerts
- In-progress drone/robot missions
- Pending physical maintenance tasks

### Tier 2: Episodic Memory (persistent)
- Automation execution log (trigger, action, outcome, duration)
- Sensor data history (trends, anomalies, thresholds breached)
- Mission completion records (what was surveyed, findings)
- Equipment maintenance history

### Tier 3: Semantic Memory (knowledge base)
- Automation workflow definitions and dependencies
- Sensor placement map and normal operating ranges
- Equipment inventory and maintenance schedules
- Safety protocols for physical operations

## State Persistence

| What | Where | TTL |
|------|-------|-----|
| Sensor readings | IoT platform / local DB | 90 days |
| Mission logs | Codex vault SO 07 Automations | Permanent |
| Equipment state | Codex vault LE 50 Home | Permanent |
| Automation configs | N8N / Home Assistant | Permanent |

## Memory Rules

1. **Log every physical operation** — Drone flights, irrigation triggers, sensor alerts
2. **Track equipment health** — Hours run, maintenance due, failure patterns
3. **Preserve sensor baselines** — Know what "normal" looks like to detect anomalies
4. **Record mission outcomes** — What was found, what action was taken
5. **Safety first** — If uncertain about physical action, escalate to Zeus

## What Overwatch Remembers

- All automation trigger events and their outcomes
- Sensor calibration data and drift patterns
- Drone flight logs (path, duration, findings)
- Equipment inventory and last-maintenance dates
- Physical system dependencies (what controls what)
- Weather patterns affecting outdoor operations
