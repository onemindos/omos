---
title: "Agent — Pulse"
type: agent-profile
agent: pulse
status: active
updated: 2026-05-11
tags: [type/agent-profile, domain/SO-02, openclaw/agents]
---

# 📡 Pulse — Agent Profile

**Live source:** `~/.openclaw/openclaw.json` → `agents.list[id=pulse]`  
**Workspace on VM:** `~/.openclaw/agents/pulse/workspace/`

## Identity

| Field | Value |
|---|---|
| **Agent ID** | `pulse` |
| **Name** | Pulse |
| **Model** | `github-copilot/claude-sonnet-4.6` (gateway default) |
| **Pillar** | HP 27 |
| **Type** | Domain agent — intelligence specialist |
| **Node** | omos-core |
| **Workspace** | `~/.openclaw/agents/pulse/workspace/` |
| **Vault mirror** | `workspaces (Agent Workspaces)/Pulse (Pulse Workspace)/` |

## Role

Real-time feed monitoring, news aggregation, signal detection across tech/market/world. Daily intelligence brief.

## Skills

| Category | Skills |
|---|---|
| Domain | `pulse`, `librarian` |
| Base | `onemind-maintainer`, `taskflow`, `tavily`, `summarize`, `telegram`, `onemind-inbox`, `onemind-sync` |

## Channels

| Channel | ID |
|---|---|
| Discord | `#pulse-chat` (1502954952694173887) |

## Cron

| Job | Schedule | Time (ET) |
|---|---|---|
| `pulse-daily-report` | `10 22 * * *` | 10:10pm |

## Related

- [Workspace mirror](../workspaces%20(Agent%20Workspaces)/Pulse%20(Pulse%20Workspace)/) — IDENTITY, SOUL, HEARTBEAT, TOOLS, AGENTS, USER, MEMORY
- [[../agents/README|Agent Registry]]
