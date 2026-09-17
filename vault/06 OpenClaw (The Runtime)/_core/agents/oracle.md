---
title: "Agent — Oracle"
type: agent-profile
agent: oracle
status: active
updated: 2026-05-11
tags: [type/agent-profile, domain/SO-02, openclaw/agents]
---

# 🔮 Oracle — Agent Profile

**Live source:** `~/.openclaw/openclaw.json` → `agents.list[id=oracle]`  
**Workspace on VM:** `~/.openclaw/workspace-oracle/`

## Identity

| Field | Value |
|---|---|
| **Agent ID** | `oracle` |
| **Name** | Oracle |
| **Model** | `github-copilot/claude-sonnet-4.6` (gateway default) |
| **Pillar** | HP 25-34 |
| **Type** | Domain agent — mindset specialist |
| **Node** | omos-core |
| **Workspace** | `~/.openclaw/workspace-oracle/` |
| **Vault mirror** | `workspaces (Agent Workspaces)/Oracle (Oracle Workspace)/` |

## Role

Mindset, habits, learning, personal finance, investments, mental performance. Weekly review of HP 25-34 domains.

## Skills

| Category | Skills |
|---|---|
| Domain | `oracle` |
| Base | `onemind-maintainer`, `taskflow`, `tavily`, `summarize`, `telegram`, `onemind-inbox`, `onemind-sync` |

## Channels

| Channel | ID |
|---|---|
| Discord | `#oracle-chat` (1480450417522643127) |

## Cron

| Job | Schedule | Time (ET) |
|---|---|---|
| `oracle-weekly-report` | `0 8 * * 1` | Mon 8am |

## Related

- [Workspace mirror](../workspaces%20(Agent%20Workspaces)/Oracle%20(Oracle%20Workspace)/) — IDENTITY, SOUL, HEARTBEAT, TOOLS, AGENTS, USER, MEMORY
- [[../agents/README|Agent Registry]]
