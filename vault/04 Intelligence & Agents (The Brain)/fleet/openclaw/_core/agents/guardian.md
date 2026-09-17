---
title: "Agent — Guardian"
type: agent-profile
agent: guardian
status: active
updated: 2026-05-11
tags: [type/agent-profile, domain/SO-02, openclaw/agents]
---

# 🛡️ Guardian — Agent Profile

**Live source:** `~/.openclaw/openclaw.json` → `agents.list[id=guardian]`  
**Workspace on VM:** `~/.openclaw/workspace-guardian/`

## Identity

| Field | Value |
|---|---|
| **Agent ID** | `guardian` |
| **Name** | Guardian |
| **Model** | `github-copilot/claude-sonnet-4.6` (gateway default) |
| **Pillar** | SO 04 |
| **Type** | Domain agent — cybersecurity specialist |
| **Node** | omos-core |
| **Workspace** | `~/.openclaw/workspace-guardian/` |
| **Vault mirror** | `workspaces (Agent Workspaces)/Guardian (Guardian Workspace)/` |

## Role

Cybersecurity, access control, privacy, threat analysis, credential management. Weekly security report.

## Skills

| Category | Skills |
|---|---|
| Domain | `guardian`, `guardian-ops` |
| Base | `onemind-maintainer`, `taskflow`, `tavily`, `summarize`, `telegram`, `onemind-inbox`, `onemind-sync` |

## Channels

| Channel | ID |
|---|---|
| Discord | `#guardian-chat` (1480445943789785180) |

## Cron

| Job | Schedule | Time (ET) |
|---|---|---|
| `guardian-weekly-report` | `0 9 * * 3` | Wed 9am |

## Related

- [Workspace mirror](../workspaces%20(Agent%20Workspaces)/Guardian%20(Guardian%20Workspace)/) — IDENTITY, SOUL, HEARTBEAT, TOOLS, AGENTS, USER, MEMORY
- [[../agents/README|Agent Registry]]
