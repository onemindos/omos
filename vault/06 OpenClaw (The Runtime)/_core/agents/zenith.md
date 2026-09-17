---
title: "Agent — Zenith"
type: agent-profile
agent: zenith
status: active
updated: 2026-05-11
tags: [type/agent-profile, domain/SO-02, openclaw/agents]
---

# 🏔 Zenith — Agent Profile

**Live source:** `~/.openclaw/openclaw.json` → `agents.list[id=zenith]`  
**Workspace on VM:** `~/.openclaw/workspace-zenith/`

## Identity

| Field | Value |
|---|---|
| **Agent ID** | `zenith` |
| **Name** | Zenith |
| **Model** | `github-copilot/claude-sonnet-4.6` |
| **Pillar** | HP 25-49 — Holistic Performance |
| **Type** | Super-agent — HP domain orchestrator |
| **Node** | omos-core |
| **Workspace** | `~/.openclaw/workspace-zenith/` |
| **Vault mirror** | `workspaces (Agent Workspaces)/zenith (Zenith Workspace)/` |

## Role

HP domain super-agent. Handles every request Legacy routes into the HP pillar. Covers body, fitness, nutrition, sleep, recovery, mindset, habits, learning, personal finance, investments, purpose, identity, and intelligence. Loads domain skills on demand for deep expertise.

## Skills

| Category | Skills |
|---|---|
| Domain | `agent-zenith`, `spartan`, `oracle`, `identity`, `librarian`, `pulse` |
| Base | `onemind-maintainer`, `taskflow`, `taskflow-inbox-triage`, `tavily`, `summarize`, `telegram`, `onemind-inbox`, `onemind-sync` |

## Channels

| Channel | ID |
|---|---|
| Discord | `#zenith-chat` |
| Telegram | topic 318 |

## Cron

No dedicated cron job. HP reports run via domain agents (spartan, oracle weekly; genesis, pulse, librarian daily).

## Related

- [Workspace mirror](../workspaces%20(Agent%20Workspaces)/zenith%20(Zenith%20Workspace)/) — IDENTITY, SOUL, HEARTBEAT, TOOLS, AGENTS, USER, MEMORY
- [[../agents/README|Agent Registry]]
