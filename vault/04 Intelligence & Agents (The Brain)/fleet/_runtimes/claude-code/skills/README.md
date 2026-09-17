---
type: index
domain: SO-01
status: active
updated: 2026-05-14
owner: zeus
tags: [runtime/claude-code, domain/so, type/index]
ai_summary: "Claude Code skill packs — .claude/skills/ loaded into sessions to give domain-specific expertise."
---

# Claude Code — Skills

Skill packs loaded into Claude Code sessions via `.claude/skills/`. Each skill is a markdown file with structured instructions that gets injected when the topic is relevant.

## What Goes Here

- `.claude/skills/` SKILL.md files mirrored here for vault tracking
- Domain skill packs (One Mind domains, programming languages, frameworks)
- Skill install/activation notes
- Custom One Mind skill definitions

## Skill Sources

| Source | Path | What |
| ------ | ---- | ---- |
| **One Mind skills (vault)** | `01 Intelligence/skills/` | Domain-specific skills — spartan, oracle, empire, guardian, etc. |
| **VS Code built-in** | `~/.vscode/extensions/copilot/assets/prompts/skills/` | VS Code agent skills |
| **Custom skills** | `~/.claude/skills/` | Personal runtime skills |

## Active Skills in `~/.claude/skills/`

| Skill | Domain | Purpose |
| ----- | ------ | ------- |
| `om-openclaw` | SO-06 | OpenClaw runtime context |
| `sentinel` | SO | Sentinel node runbook |
| `composio-cli` | SO | Composio CLI operations |
| *(add as installed)* | — | — |

---
*Full skill library: `01 Intelligence/skills/`. Runtime copies live at `~/.claude/skills/` on Mac.*
