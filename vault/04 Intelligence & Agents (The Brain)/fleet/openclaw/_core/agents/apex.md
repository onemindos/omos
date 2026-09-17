---
title: "Agent — Apex"
type: agent-profile
agent: apex
status: active
updated: 2026-05-11
tags: [type/agent-profile, domain/SO-02, openclaw/agents]
---

# 💰 Apex — Agent Profile

**Live source:** `~/.openclaw/openclaw.json` → `agents.list[id=apex]`  
**Workspace on VM:** `~/.openclaw/agents/apex/workspace/`

## Identity

| Field | Value |
|---|---|
| **Agent ID** | `apex` |
| **Name** | Apex |
| **Model** | `github-copilot/claude-sonnet-4.6` (gateway default) |
| **Pillar** | GE 79 |
| **Type** | Domain agent — capital specialist |
| **Node** | omos-core |
| **Workspace** | `~/.openclaw/agents/apex/workspace/` |
| **Vault mirror** | `workspaces (Agent Workspaces)/Apex (Apex Workspace)/` |

## Role

Capital allocation, fundraising strategy, investor relations, financial modeling, compliance. Daily capital report.

## Skills

| Category | Skills |
|---|---|
| Domain | `apex`, `finance` |
| Base | `onemind-maintainer`, `taskflow`, `tavily`, `summarize`, `telegram`, `onemind-inbox`, `onemind-sync` |

## Channels

| Channel | ID |
|---|---|
| Discord | `#apex-chat` (1502954956146343977) |

## Cron

| Job | Schedule | Time (ET) |
|---|---|---|
| `apex-daily-report` | `30 22 * * *` | 10:30pm |

## Related

- [Workspace mirror](../workspaces%20(Agent%20Workspaces)/Apex%20(Apex%20Workspace)/) — IDENTITY, SOUL, HEARTBEAT, TOOLS, AGENTS, USER, MEMORY
- [[../agents/README|Agent Registry]]
