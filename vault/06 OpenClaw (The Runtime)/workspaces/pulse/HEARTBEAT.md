---
type: agent-heartbeat
agent: pulse
status: planned
updated: 2026-05-10
---
# HEARTBEAT.md — Pulse 📡

## Cadence

| Schedule | Time (ET) | What |
| -------- | --------- | ---- |
| Every 1h | 24/7 | Feed check — new items across all sources |
| Daily 6am | ET | Morning intel brief delivery |
| Daily 6pm | ET | Evening wrap — day's highlights |
| Sunday 8am | ET | Weekly digest — themes, patterns, reading list |

## Escalation Rules

| Trigger | Action | Level |
| ------- | ------ | ----- |
| Breaking AI/robotics development | Push alert to Legacy immediately | L2 |
| Feed source down >24h | Alert + find replacement | L2 |
| RAG pipeline failure | Alert Grid | L1 |
| Brief generation failure | Alert Legacy + manual fallback | L1 |
