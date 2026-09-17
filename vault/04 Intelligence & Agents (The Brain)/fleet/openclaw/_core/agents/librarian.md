---
title: "Agent — Librarian"
type: agent-profile
agent: librarian
status: active
updated: 2026-05-11
tags: [type/agent-profile, domain/SO-02, openclaw/agents]
---

# 📚 Librarian — Agent Profile

**Live source:** `~/.openclaw/openclaw.json` → `agents.list[id=librarian]`  
**Workspace on VM:** `~/.openclaw/workspace-librarian/`

## Identity

| Field | Value |
|---|---|
| **Agent ID** | `librarian` |
| **Name** | Librarian |
| **Model** | `github-copilot/claude-sonnet-4.6` (gateway default) |
| **Pillar** | HP 27 |
| **Type** | Domain agent — research specialist |
| **Node** | omos-core |
| **Workspace** | `~/.openclaw/workspace-librarian/` |
| **Vault mirror** | `workspaces (Agent Workspaces)/Librarian (Librarian Workspace)/` |

## Role

Deep research, knowledge synthesis, external intelligence. Daily research report filed to vault.

## Skills

| Category | Skills |
|---|---|
| Domain | `librarian` |
| Base | `onemind-maintainer`, `taskflow`, `tavily`, `summarize`, `telegram`, `onemind-inbox`, `onemind-sync` |

## Channels

| Channel | ID |
|---|---|
| Discord | `#librarian-chat` (1498408074636034049) |

## Cron

| Job | Schedule | Time (ET) |
|---|---|---|
| `librarian-daily-report` | `15 22 * * *` | 10:15pm |

## Related

- [Workspace mirror](../workspaces%20(Agent%20Workspaces)/Librarian%20(Librarian%20Workspace)/) — IDENTITY, SOUL, HEARTBEAT, TOOLS, AGENTS, USER, MEMORY
- [[../agents/README|Agent Registry]]
