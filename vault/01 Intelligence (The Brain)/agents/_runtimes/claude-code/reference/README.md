---
type: index
domain: SO-01
status: active
updated: 2026-05-14
owner: zeus
tags: [runtime/claude-code, domain/so, type/index]
ai_summary: "Claude Code reference — system prompt, CLAUDE.md conventions, behavior notes, version history, and configuration reference."
---

# Claude Code — Reference

System-level reference docs for Claude Code. How it works, how it's configured in One Mind, and version notes.

## What Goes Here

- System prompt snapshots (versioned)
- `CLAUDE.md` writing conventions and best practices
- Behavior notes — what Claude Code does and doesn't do by default
- Version changelog (capabilities added per version)
- Settings reference (`~/.claude/settings.json` schema)
- Known quirks and workarounds

## Key Files

| File | What |
| ---- | ---- |
| `system-prompt.md` | Current Claude Code system prompt snapshot |
| `settings-reference.md` | Full `~/.claude/settings.json` schema with One Mind config |
| `claude-md-conventions.md` | How to write effective `CLAUDE.md` files |

## Current Config State

| Setting | Value |
| ------- | ----- |
| **Model** | `claude-sonnet-4.6` (via GitHub Copilot) |
| **MCP servers** | context7, github, filesystem, desktop-commander |
| **Skills path** | `~/.claude/skills/` |
| **Agents path** | `~/.claude/agents/` |
| **Hooks** | See `hooks/` folder |
| **CLAUDE.md** | Vault root + per-repo |

## Version Notes

| Version | Notable Capabilities |
| ------- | -------------------- |
| Current (May 2026) | Sub-agents, hooks, MCP, skills, Claude Code Teams, Co-Work |
| *(track changes here as Claude Code updates)* | — |

---
*System prompt archive: `01 Intelligence/prompts/system-prompts-reference/anthropic/claude.txt`*
*Platform info: `05 Platforms/tools/ai/anthropic/claude-code/`*
