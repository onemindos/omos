---
title: "Agent — Overwatch"
type: agent-profile
agent: overwatch
status: active
updated: 2026-05-11
tags: [type/agent-profile, domain/SO-02, openclaw/agents]
---

# 👁️ Overwatch — Agent Profile

**Live source:** `~/.openclaw/openclaw.json` → `agents.list[id=overwatch]`  
**Workspace on VM:** `~/.openclaw/workspace-overwatch/`

## Identity

| Field | Value |
|---|---|
| **Agent ID** | `overwatch` |
| **Name** | Overwatch |
| **Model** | `github-copilot/claude-sonnet-4.6` (gateway default) |
| **Pillar** | SO 07 |
| **Type** | Domain agent — automation specialist |
| **Node** | omos-core |
| **Workspace** | `~/.openclaw/workspace-overwatch/` |
| **Vault mirror** | `workspaces (Agent Workspaces)/Overwatch (Overwatch Workspace)/` |

## Role

n8n workflow management, process automation, scheduled tasks, operations coordination. Daily automation report.

## Skills

| Category | Skills |
|---|---|
| Domain | `overwatch` |
| Base | `onemind-maintainer`, `taskflow`, `tavily`, `summarize`, `telegram`, `onemind-inbox`, `onemind-sync` |

## Channels

| Channel | ID |
|---|---|
| Discord | `#overwatch-chat` (1502954950764793870) |

## Cron

| Job | Schedule | Time (ET) |
|---|---|---|
| `overwatch-daily-report` | `45 22 * * *` | 10:45pm |

## Related

- [Workspace mirror](../workspaces%20(Agent%20Workspaces)/Overwatch%20(Overwatch%20Workspace)/) — IDENTITY, SOUL, HEARTBEAT, TOOLS, AGENTS, USER, MEMORY
- [[../agents/README|Agent Registry]]
