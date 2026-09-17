# OpenClaw

> **Self-hosted AI agent gateway and runtime for One Mind**

## Overview

OpenClaw is the agent platform running all One Mind domain agents on `onemind-main` (VPS). It provides isolated workspaces, memory search, WebSocket streaming, and multi-agent orchestration. The primary interface between Zeus (via Telegram) and the 9 specialized agents.

## Status

| Attribute   | Value                               |
| ----------- | ----------------------------------- |
| Status      | ✅ Active                           |
| Category    | AI Agents                           |
| Type        | Agent Runtime / Gateway             |
| Integration | REST API, WebSocket, systemd        |

## Use Cases in OneMind

1. **Agent Execution** — Runs Legacy, Spartan, Oracle, Haven, Empire, Guardian, Grid, Trinity, Eden
2. **Telegram Gateway** — Routes messages to the correct domain agent
3. **Memory Indexing** — Embeds vault content into pgvector via `agents.defaults.memorySearch.extraPaths`
4. **Workspace Isolation** — Each agent has its own IDENTITY.md, SOUL.md, shared CODEX.md + USER.md

## Key Features

- 9 domain agents with isolated workspaces
- Shared context via symlinked files in `/home/onemind/onemind-os/openclaw/shared/`
- Universal vault search across full onemind-codex
- systemd services: `openclaw.service`, `openclaw-reload.service`
- WebSocket streaming + REST API

## Operations

```bash
# VPS access
ssh root@104.248.4.212

# Service management
systemctl status openclaw
systemctl restart openclaw
journalctl -u openclaw -f

# Agent workspaces
/home/onemind/onemind-os/openclaw/workspaces/{agent}/
```

## Links

- VPS: `onemind-main` — `104.248.4.212`
- Shared context: `/home/onemind/onemind-os/openclaw/shared/`
- Codex docs: `00-24 SO/03 Agents/`


---

## Related

- [[00-24 SO (Sovereign Operations)/02 Infrastructure (VPS, Servers & Networking) (The Body)/02 Infrastructure|02 Infrastructure]]
