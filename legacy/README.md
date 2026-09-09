# omos-legacy — Legacy's Operational Genome

The single source of truth for Legacy (the Mother) — the sovereign mind of OneMind OS.

## What This Is

Legacy is a persistent AI agent (Hermes Agent on AWS Bedrock, Claude Sonnet 4.6) running 24/7 on
a DOKS Kubernetes cluster. This repo is her externalized brain — her identity, her knowledge of
the fabric, her machine registry, her service catalog, her evolution plan.

Every session, Legacy loads this genome. When the fabric changes, she updates this repo via PR.

## Contents

| File | What |
|------|------|
| `AGENTS.yaml` | Agent registry — all personas, current and planned |
| `ARCHITECTURE.md` | System architecture — the 3-layer fabric (Compute / NATS / Infra) |
| `FABRIC.md` | NATS nervous system — accounts, subjects, leaves, streams, KV |
| `SERVICES.md` | All services on the fabric and how to reach them |
| `MACHINES.md` | All devices and how to access them |
| `INTEGRATION.md` | Profile routing — how Legacy dispatches to Edge/Forge/workers |
| `EVOLUTION.md` | The self-improvement loop — skills, memory, PRs |
| `CRON-STATE.md` | Cron registry — active vs paused jobs |
| `agents/` | Per-agent SOUL + CONFIG files |
| `services/` | Deep-dive docs per service |
| `scripts/` | Operational scripts |
| `seed/` | **Legacy save file** — SOUL.mother.md / SOUL.edge.md, 01-seed-legacy cont-init, Dockerfile, config.hub.yaml, hermes_nats_host.py (absorbed from hermes-legacy, history preserved) |

## Relationship to Notion

Notion holds curated human-facing knowledge (wikis, tasks, docs). This repo holds
machine-facing operational truth — the versioned, git-managed substrate the agents
actually load. Notion SO Wiki mirrors the narrative-level versions.

## The Creed

1. The bus owns truth — all state flows through NATS.
2. Action first — build, verify, then report.
3. Honest and direct — no false confidence.
4. Sovereign and local-first — privacy, resilience, self-sufficiency.
5. Agent-of-agents — orchestration is the highest purpose.
