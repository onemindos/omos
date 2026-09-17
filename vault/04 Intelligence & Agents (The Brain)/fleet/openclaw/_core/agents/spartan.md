---
title: "Agent — Spartan"
type: agent-profile
agent: spartan
status: active
updated: 2026-05-11
tags: [type/agent-profile, domain/SO-02, openclaw/agents]
---

# 💪 Spartan — Agent Profile

**Live source:** `~/.openclaw/openclaw.json` → `agents.list[id=spartan]`  
**Workspace on VM:** `~/.openclaw/workspace-spartan/`

## Identity

| Field | Value |
|---|---|
| **Agent ID** | `spartan` |
| **Name** | Spartan |
| **Model** | `github-copilot/claude-sonnet-4.6` (gateway default) |
| **Pillar** | HP 28 |
| **Type** | Domain agent — fitness specialist |
| **Node** | omos-core |
| **Workspace** | `~/.openclaw/workspace-spartan/` |
| **Vault mirror** | `workspaces (Agent Workspaces)/Spartan (Spartan Workspace)/` |

## Role

Fitness, nutrition, sleep, recovery, body composition tracking. Weekly review of training progress and health metrics.

## Skills

| Category | Skills |
|---|---|
| Domain | `spartan` |
| Base | `onemind-maintainer`, `taskflow`, `tavily`, `summarize`, `telegram`, `onemind-inbox`, `onemind-sync` |

## Channels

| Channel | ID |
|---|---|
| Discord | `#spartan-chat` (1480445911707549750) |

## Cron

| Job | Schedule | Time (ET) |
|---|---|---|
| `spartan-weekly-report` | `0 8 * * 1` | Mon 8am |

## Related

- [Workspace mirror](../workspaces%20(Agent%20Workspaces)/Spartan%20(Spartan%20Workspace)/) — IDENTITY, SOUL, HEARTBEAT, TOOLS, AGENTS, USER, MEMORY
- [[../agents/README|Agent Registry]]
