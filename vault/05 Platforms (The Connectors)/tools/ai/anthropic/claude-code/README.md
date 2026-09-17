---
type: platform
domain: SO-05
status: active
updated: 2026-05-14
owner: zeus
tags: [platform/anthropic, product/claude-code, domain/so]
ai_summary: "Claude Code — Anthropic's agentic coding tool. Runs in VS Code, terminal, and CI. Deep codebase context, sub-agents, hooks, MCP, and Teams."
---

# Claude Code — Agentic Coding Tool

Anthropic's agentic coding tool. Runs inside VS Code (via GitHub Copilot), terminal (`claude` CLI), and CI pipelines.

## Access

- VS Code: GitHub Copilot extension (Claude Sonnet 4.6 backend)
- Terminal: `claude` CLI — install via `npm install -g @anthropic-ai/claude-code`
- Docs: `docs.anthropic.com/claude-code`

## Current Setup in One Mind

| Setting | Value |
| ------- | ----- |
| VS Code model | `claude-sonnet-4.6` (via GitHub Copilot) |
| CLI model | `claude-sonnet-4.6` |
| CLAUDE.md | Vault root + per-repo |
| MCP servers | context7, github, filesystem, desktop-commander |
| Skills | `~/.claude/skills/` |
| Sub-agents | `~/.claude/agents/` |
| Hooks | `~/.claude/hooks/` |

## Key Capabilities

| Capability | What |
| ---------- | ---- |
| **Sub-agents** | Specialized roles mid-session — security reviewer, test writer, etc. |
| **MCP plugins** | Tool extensions — file access, GitHub, docs, etc. |
| **Hooks** | Lifecycle events — PreToolUse, PostToolUse, Stop, Notification |
| **Skills** | Domain packs injected as context when topics match |
| **Claude Code Teams** | Shared instructions and settings across team members |
| **Co-Work** | Collaborative sessions (see `claude-cowork/`) |
| **Extended thinking** | Deep reasoning mode for complex problems |
| **CLAUDE.md** | Project-aware context — automatically loaded per repo |

## Operational Config

Full runtime config → `01 Intelligence/agents/_runtimes/claude-code/`

## Pricing

- Included with Claude.ai Pro/Team via GitHub Copilot
- Direct API usage billed to Anthropic Console
- Track costs: `platform/` folder
