---
title: "Agent — Sovereignty"
type: agent-profile
agent: sovereignty
status: active
updated: 2026-05-11
tags: [type/agent-profile, domain/SO-02, openclaw/agents]
---

# 🌐 Sovereignty — Agent Profile

**Live source:** `~/.openclaw/openclaw.json` → `agents.list[id=sovereignty]`  
**Workspace on VM:** `~/.openclaw/workspace-sovereignty/`

## Identity

| Field | Value |
|---|---|
| **Agent ID** | `sovereignty` |
| **Name** | Sovereignty |
| **Model** | `github-copilot/claude-sonnet-4.6` |
| **Pillar** | SO 00-24 — Sovereign Operations |
| **Type** | Super-agent — SO domain orchestrator |
| **Node** | omos-core |
| **Workspace** | `~/.openclaw/workspace-sovereignty/` |
| **Vault mirror** | `workspaces (Agent Workspaces)/sovereignty (Sovereignty Workspace)/` |

## Role

SO domain super-agent. Handles infrastructure, security, automation, development, and research. Loads domain skills on demand. The operational backbone agent — runs the systems that everything else depends on.

## Skills

| Category | Skills |
|---|---|
| Domain | `agent-sovereignty`, `grid`, `guardian`, `guardian-ops`, `overwatch`, `eagleeye` |
| Dev | `coding-agent`, `github`, `diffs` |
| Base | `onemind-maintainer`, `taskflow`, `taskflow-inbox-triage`, `tavily`, `summarize`, `telegram`, `onemind-inbox`, `onemind-sync` |

## Channels

| Channel | ID |
|---|---|
| Discord | `#sovereignty-chat` |
| Telegram | topic 321 |

## Cron

No dedicated cron job. SO reports run via domain agents (guardian weekly; grid, trinity, overwatch, eagleeye daily).

## Related

- [Workspace mirror](../workspaces%20(Agent%20Workspaces)/sovereignty%20(Sovereignty%20Workspace)/) — IDENTITY, SOUL, HEARTBEAT, TOOLS, AGENTS, USER, MEMORY
- [[../agents/README|Agent Registry]]
