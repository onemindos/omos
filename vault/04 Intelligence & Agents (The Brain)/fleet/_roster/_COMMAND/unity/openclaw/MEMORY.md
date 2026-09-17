---
type: memory-architecture
agent: unity
pillar: LE 50-74
status: active
updated: 2026-05-11
---
# MEMORY.md — 🌿 Unity

## Memory Focus

Family calendar patterns, homestead seasonal cycles, estate document dates, elder care needs, children's milestones.

## Tiers

### Tier 1: Working Memory (session-scoped)
**Working memory**: Active family calendar, pending appointments, current homestead tasks

### Tier 2: Episodic Memory (persistent)
**Episodic memory**: Milestone records, estate change history, elder care plan changes

### Tier 3: Semantic Memory (knowledge base)
**Semantic memory**: Family member profiles, property details, recurring seasonal tasks, food store baselines

## State Persistence

| What | Where | TTL |
|------|-------|-----|
| Active tasks | OpenClaw workspace memory | Session |
| Domain state | Codex vault LE 50-74/ | Permanent |
| Zeus preferences | Codex vault _codex/memory/relationship/ | Permanent |

## Memory Rules

1. Never forget Zeus's corrections — store permanently
2. Log when skill activations get redirected (signals wrong skill choice)
3. Preserve context across sessions — key decisions survive gateway restarts
4. Don't store secrets — reference 1Password paths, never values
5. Prune stale state — tasks older than 7 days without activity get archived
