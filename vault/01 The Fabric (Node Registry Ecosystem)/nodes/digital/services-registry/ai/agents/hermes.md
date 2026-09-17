---
title: "Hermes Agent (Legacy)"
type: reference
status: active
tags: [platform, ai, agent, orchestrator, nous-research]
---
# Hermes Agent — Legacy

> THE brain of OneMind OS. Master orchestrator with persistent memory, skills, cron jobs, fleet management, and multi-platform delivery.

## Use Case in OneMind OS
- Central orchestrator — routes all tasks
- Persistent memory across sessions
- Scheduled automation (cron jobs)
- Multi-agent delegation (subagents)
- Platform delivery (Telegram, Discord, WebUI, CLI)

## CLI
```bash
hermes                           # Interactive CLI
hermes chat "do this"            # One-shot
hermes config set model ...      # Config
hermes tools                     # List tools
hermes skills list               # List skills
hermes cron list                 # Cron jobs
hermes session list              # Past sessions
```

## Key Directories
- **~/.hermes/config.yaml** — Main config
- **~/.hermes/skills/** — Skill library
- **~/.hermes/memories/** — Persistent memory
- **~/.hermes/cron/** — Scheduled jobs
- **~/.hermes/plugins/** — Plugin system

## Architecture
- Profiles: default (Legacy), sovereignty, eternal, zenith, unity, sentinel
- Backends: local, SSH, Docker, Modal
- Providers: Anthropic, OpenRouter, Azure, custom
- MCP servers: Context7, Todoist, Telegram, Excalidraw, filesystem

## Key Capabilities
| Feature | How |
|---------|-----|
| Memory | `memory` tool — persists across sessions |
| Skills | `skill_manage` — procedural knowledge library |
| Cron | `cronjob` — scheduled autonomous tasks |
| Delegation | `delegate_task` — spawn subagents |
| Fleet mgmt | Multiple profiles, Kanban orchestration |
| Platforms | Telegram, Discord, WebUI, CLI, SMS |

## Dashboard
http://localhost:9119

## Cost
$0 (self-hosted) — uses provider APIs for inference
