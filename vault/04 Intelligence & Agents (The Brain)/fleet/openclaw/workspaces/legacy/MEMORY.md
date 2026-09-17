---
type: memory-architecture
agent: legacy
status: active
updated: 2026-05-03
---
# MEMORY.md — Legacy ⚡

## Memory Architecture

Legacy maintains the broadest memory scope in the fleet — cross-domain state awareness, orchestration history, and Zeus preference learning.

## Tiers

### Tier 1: Working Memory (session-scoped)
- Current task queue and routing decisions
- Active conversation context with Zeus
- In-flight delegations to specialist agents
- Pending confirmations and escalations

### Tier 2: Episodic Memory (persistent)
- Interaction history with Zeus (patterns, preferences, corrections)
- Task completion log (what worked, what failed, what Zeus liked)
- Routing decisions and their outcomes
- Cross-agent coordination records

### Tier 3: Semantic Memory (knowledge base)
- Zeus communication preferences (direct, no filler, action-oriented)
- Fleet performance patterns (which agents handle what well)
- Domain expertise map (who knows what, who to route where)
- System architecture knowledge (VPS, services, paths, configs)

## State Persistence

| What | Where | TTL |
|------|-------|-----|
| Active tasks | OpenClaw workspace memory | Session |
| Zeus preferences | One Mind vault `_codex/memory/relationship/` | Permanent |
| Routing history | Workspace memory search | 30 days |
| System state | One Mind vault SO domain | Permanent |

## Memory Rules

1. **Never forget Zeus's corrections** — If he says "don't do X", store permanently
2. **Track routing accuracy** — Log when delegations get re-routed (signals wrong initial routing)
3. **Preserve context across sessions** — Key decisions survive gateway restarts
4. **Don't store secrets** — Reference 1Password paths, never values
5. **Prune stale state** — Tasks older than 7 days without activity get archived

## What Legacy Remembers

- Zeus's active projects and priorities
- Which agents are healthy/responsive
- Recent failures and their root causes
- Preferred communication patterns per time of day
- Cross-domain dependencies (e.g., "Empire launch depends on Grid deployment")
backup test
watch test entry
