---
title: "Codex CLI"
type: reference
status: active
tags: [platform, ai, agent, coding, openai]
---
# Codex CLI (OpenAI)

> OpenAI's autonomous coding agent. Terminal-based, sandboxed execution.

## Use Case in OneMind OS
- Feature implementation with sandboxed shell access
- Reads `AGENTS.md` for project-specific instructions
- Good for tasks needing shell commands during development

## CLI
```bash
codex                            # Interactive
codex "implement the auth flow"  # One-shot
codex --model o4-mini            # Specify model
codex --approval auto            # No confirmation needed
codex --full-auto                # Full autonomy
```

## Key Files
- **AGENTS.md** — Project instructions (root of repo)
- **~/.codex/config.yaml** — Global config
- **~/.codex/instructions.md** — Global instructions

## Integration with Hermes
- Hermes skill: `codex` — delegates coding tasks
- Hermes skill: `kanban-codex-lane` — Kanban worker mode
- Spawned via terminal with `--full-auto` for autonomous work

## Cost
Uses OpenAI API (o4-mini, GPT-4.1) — included in Azure ~$30/mo
