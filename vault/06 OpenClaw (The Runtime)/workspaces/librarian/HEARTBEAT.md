---
type: agent-heartbeat
agent: librarian
status: active
updated: 2026-05-03
---
# HEARTBEAT.md — Librarian

## Cadence

| Schedule | Time (ET) | What |
| -------- | --------- | ---- |
| Morning intel brief | 6:30am | Top developments across all feed categories |
| Feed ingestion cycle | Every 2h | Scan all sources, categorize, queue |
| Trend analysis | Daily 2:00pm | Cross-domain pattern detection |
| Weekly deep dive | Friday 3:00pm | One-topic deep research report |
| RAG pipeline refresh | Daily 4:00am | Index new content, prune stale entries |
| Monthly source audit | 1st of month | Review feed quality, add/remove sources |

## Escalation Rules

| Trigger | Action | Level |
| ------- | ------ | ----- |
| Critical development in active domain | Immediate alert to Zeus via Legacy | L1 |
| Trend affects multiple OneMind domains | Flag in morning brief + separate alert | L2 |
| Source goes offline or compromised | Log + find alternative | L3 |
| Deep research exceeds 2h estimate | Notify Zeus with partial results | L2 |
| Unverified but potentially critical intel | Flag as UNVERIFIED, include anyway | L2 |
| RAG pipeline error | Alert Grid for infrastructure support | L2 |
