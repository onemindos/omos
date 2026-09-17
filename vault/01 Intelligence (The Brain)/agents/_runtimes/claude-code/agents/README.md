---
type: index
domain: SO-01
status: active
updated: 2026-05-14
owner: zeus
tags: [runtime/claude-code, domain/so, type/index]
ai_summary: "Claude Code agent profiles — sub-agents, specialized agents, and .claude/agents configs."
---

# Claude Code — Agents

Agent profiles and sub-agent definitions for Claude Code sessions.

## What Goes Here

- `.claude/agents/` YAML configs for sub-agents
- Specialized agent profiles (security reviewer, code architect, test writer, etc.)
- Agent permission sets and tool restrictions
- Per-project agent configs
- Agent customization instructions (`.instructions.md`, `.agent.md`)

## Agent Types

| Type | Purpose |
| ---- | ------- |
| **Sub-agents** | Specialized roles spun up mid-session (e.g. `security-reviewer`, `test-writer`) |
| **Project agents** | `.claude/agents/` in a specific repo — scoped to that project |
| **Workspace agents** | Global agents available in any Claude Code session |
| **Copilot agent** | VS Code Copilot agent mode — GitHub Copilot integration |

## Active Agents

| Agent | Scope | Purpose |
| ----- | ----- | ------- |
| *(none yet — add as created)* | — | — |

---
*Route new agent configs here. Global configs also go in `~/.claude/agents/` on the Mac.*
