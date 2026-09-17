---
type: memory-architecture
agent: sovereignty
pillar: SO 00-24
status: active
updated: 2026-05-11
---
# MEMORY.md — 🌐 Sovereignty

## Memory Focus

Infrastructure topology, service health baselines, credential rotation schedule, known vulnerabilities, automation workflow state.

## Tiers

### Tier 1: Working Memory (session-scoped)
**Working memory**: Current service health, active deployments, open PRs, running workflows

### Tier 2: Episodic Memory (persistent)
**Episodic memory**: Incident history, deployment outcomes, security anomalies, backup verification log

### Tier 3: Semantic Memory (knowledge base)
**Semantic memory**: Service topology, credential locations (1Password paths), cert expiry dates, n8n workflow inventory

## State Persistence

| What | Where | TTL |
|------|-------|-----|
| Active tasks | OpenClaw workspace memory | Session |
| Domain state | Codex vault SO 00-24/ | Permanent |
| Zeus preferences | Codex vault _codex/memory/relationship/ | Permanent |

## Memory Rules

1. Never forget Zeus's corrections — store permanently
2. Log when skill activations get redirected (signals wrong skill choice)
3. Preserve context across sessions — key decisions survive gateway restarts
4. Don't store secrets — reference 1Password paths, never values
5. Prune stale state — tasks older than 7 days without activity get archived
