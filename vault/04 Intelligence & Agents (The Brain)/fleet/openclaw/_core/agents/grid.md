---
title: "Agent — Grid"
type: agent-profile
agent: grid
status: active
updated: 2026-05-11
tags: [type/agent-profile, domain/SO-02, openclaw/agents]
---

# ⚙️ Grid — Agent Profile

**Live source:** `~/.openclaw/openclaw.json` → `agents.list[id=grid]`  
**Workspace on VM:** `~/.openclaw/workspace-grid/`

## Identity

| Field | Value |
|---|---|
| **Agent ID** | `grid` |
| **Name** | Grid |
| **Model** | `github-copilot/claude-sonnet-4.6` (gateway default) |
| **Pillar** | SO 02 |
| **Type** | Domain agent — infrastructure specialist |
| **Node** | omos-core |
| **Workspace** | `~/.openclaw/workspace-grid/` |
| **Vault mirror** | `workspaces (Agent Workspaces)/Grid (Grid Workspace)/` |

## Role

VPS management, infrastructure, networking, API integrations, system configuration. Daily infra report.

## Skills

| Category | Skills |
|---|---|
| Domain | `grid`, `coding-agent`, `github`, `diffs` |
| Base | `onemind-maintainer`, `taskflow`, `tavily`, `summarize`, `telegram`, `onemind-inbox`, `onemind-sync` |

## Channels

| Channel | ID |
|---|---|
| Discord | `#grid-chat` (1480445951805096007) |

## Cron

| Job | Schedule | Time (ET) |
|---|---|---|
| `grid-daily-report` | `0 22 * * *` | 10:00pm |

## Related

- [Workspace mirror](../workspaces%20(Agent%20Workspaces)/Grid%20(Grid%20Workspace)/) — IDENTITY, SOUL, HEARTBEAT, TOOLS, AGENTS, USER, MEMORY
- [[../agents/README|Agent Registry]]
