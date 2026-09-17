---
type: index
domain: SO-01
status: active
updated: 2026-05-14
owner: zeus
tags: [runtime/claude-code, domain/so, type/index]
ai_summary: "Claude Code Teams — team members, shared instructions, permissions, and collaborative session configs."
---

# Claude Code — Teams

Claude Code Teams configuration and membership. Teams share instructions, settings, and context across members.

## What Goes Here

- Team member roster + permission levels
- Shared `.claude/` instruction files published to the team
- Team-wide MCP server configs
- Onboarding guide for new team members
- Session sharing and collaboration notes

## One Mind Team

| Member | Role | Access Level |
| ------ | ---- | ------------ |
| Zeus | Owner | Full — all repos, all agents, all tools |
| *(add members as onboarded)* | — | — |

## Shared Instructions

Files in this folder that get distributed to all team members:

| File | Purpose |
| ---- | ------- |
| `CLAUDE.md` (vault root) | Primary workspace context for all Claude Code sessions in One Mind |
| `AGENTS.md` (vault root) | Universal agent instructions — routing, naming, structure |

## Team Onboarding

When adding a new team member to Claude Code Teams:

1. Invite via `claude.ai/teams` → Members
2. Share the One Mind repo access (GitHub org: `OneMind-OS`)
3. Point them to `CLAUDE.md` and `AGENTS.md` as the starting context
4. Add their name and role to the roster above

---
*Manage team at: `claude.ai/teams` → Settings*
*Billing and seats: `05 Platforms/tools/ai/anthropic/platform/`*
