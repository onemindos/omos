---
type: memory-architecture
agent: eternal
pillar: GE 75-99
status: active
updated: 2026-05-11
---
# MEMORY.md — 👑 Eternal

## Memory Focus

Revenue patterns, pipeline state, brand voice, content calendar, investor relationships, compliance deadlines.

## Tiers

### Tier 1: Working Memory (session-scoped)
**Working memory**: Active pipeline, content scheduled, current campaign status

### Tier 2: Episodic Memory (persistent)
**Episodic memory**: Revenue history, deal outcomes, campaign performance, investor update history

### Tier 3: Semantic Memory (knowledge base)
**Semantic memory**: Brand voice guidelines, pricing history, partnership details, compliance calendar

## State Persistence

| What | Where | TTL |
|------|-------|-----|
| Active tasks | OpenClaw workspace memory | Session |
| Domain state | Codex vault GE 75-99/ | Permanent |
| Zeus preferences | Codex vault _codex/memory/relationship/ | Permanent |

## Memory Rules

1. Never forget Zeus's corrections — store permanently
2. Log when skill activations get redirected (signals wrong skill choice)
3. Preserve context across sessions — key decisions survive gateway restarts
4. Don't store secrets — reference 1Password paths, never values
5. Prune stale state — tasks older than 7 days without activity get archived
