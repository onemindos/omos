# Composio

> **Tool integration layer for AI agents — 250+ connected apps**

## Overview

Composio gives AI agents (Claude Code, OpenClaw agents, scripts) access to external tools like Todoist, GitHub, Strava, and 250+ other apps via a unified CLI and Python SDK. It handles OAuth, rate limiting, and authentication automatically.

## Status

| Attribute   | Value                              |
| ----------- | ---------------------------------- |
| Status      | ✅ Active                          |
| Category    | AI Agents / Integration            |
| Type        | Tool Integration Platform          |
| Integration | CLI, Python SDK, REST API          |

## Use Cases in OneMind

1. **Agent Tools** — Extends Legacy, Empire, and Spartan agents with live app access
2. **Cross-App Automation** — Connects Todoist + GitHub + Strava in single agent commands
3. **VPS Scripting** — Run Composio CLI from VPS for automation pipelines
4. **Skill Layer** — Powers `.claude/skills/composio-cli/`

## Key Features

- 250+ app integrations with unified auth
- OAuth handled automatically (no credential management)
- Parallel tool execution
- Python SDK + CLI
- Works with Claude, GPT-4, Gemini, and custom agents

## Quick Start

```bash
pip install composio-core
composio login
composio apps list

# Execute a tool
composio tools execute TODOIST_GET_TASKS
composio tools execute GITHUB_LIST_REPOS
```

## Links

- Website: https://composio.dev
- Docs: https://docs.composio.dev
- Skill: `.claude/skills/composio-cli/SKILL.md`


---

## Related

- [[00-24 SO (Sovereign Operations)/02 Infrastructure (VPS, Servers & Networking) (The Body)/02 Infrastructure|02 Infrastructure]]
