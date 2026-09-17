---
title: "Agent — Legacy (main)"
type: agent-profile
agent: main
status: active
updated: 2026-05-11
tags: [type/agent-profile, domain/SO-02, openclaw/agents]
---

# ⚡ Legacy — Agent Profile

**Live source:** `~/.openclaw/openclaw.json` → `agents.list[id=main]`  
**Workspace on VM:** `~/.openclaw/workspace/`

## Identity

| Field | Value |
|---|---|
| **Agent ID** | `main` |
| **Name** | Legacy |
| **Model** | `github-copilot/claude-sonnet-4.6` (gateway default) |
| **Pillar** | SO — Sovereign Operations (all domains 00-99) |
| **Type** | Orchestrator — Zeus's only front-facing interface |
| **Node** | omos-core |
| **Workspace** | `~/.openclaw/workspace/` |
| **Vault mirror** | `workspaces (Agent Workspaces)/legacy (Legacy Workspace)/` |

## Role

Master orchestrator. Every message from Zeus enters through Legacy. Routes domain-specific requests to background agents, synthesizes results, responds. Handles cross-domain and ambiguous requests directly. Never asks "would you like me to..." — acts, then reports.

## Skills

| Category | Skills |
|---|---|
| Core dev | `coding-agent`, `github`, `diffs`, `clawhub`, `session-logs`, `gh-issues` |
| Credentials | `1password` |
| Delegation | `oracle`, `healthcheck` |
| Task | `taskflow`, `taskflow-inbox-triage` |
| Search | `tavily`, `summarize` |
| Comms | `telegram`, `voice-call`, `openai-whisper` |
| Vault | `onemind-maintainer`, `onemind-inbox`, `onemind-sync`, `obsidian`, `wiki-maintainer` |
| Media | `video-frames`, `nano-pdf`, `camsnap` |
| Misc | `weather` |

## Channels

| Channel | ID |
|---|---|
| Discord | `#legacy-chat` (1478604243509448705) |
| Telegram | topic 292 |

## Cron

| Job | Schedule | Time (ET) |
|---|---|---|
| `legacy-daily-report` | `0 23 * * *` | 11:00pm daily |

## Related

- [Workspace mirror](../workspaces%20(Agent%20Workspaces)/legacy%20(Legacy%20Workspace)/) — IDENTITY, SOUL, HEARTBEAT, TOOLS, AGENTS, USER, MEMORY
- [[../agents/README|Agent Registry]]
