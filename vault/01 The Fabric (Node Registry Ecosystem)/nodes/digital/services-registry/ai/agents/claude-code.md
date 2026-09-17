---
title: "Claude Code"
type: reference
status: active
tags: [platform, ai, agent, coding, anthropic]
---
# Claude Code

> Autonomous coding agent from Anthropic. Delegates full feature implementation, PR creation, and code review.

## Use Case in OneMind OS
- Feature implementation from spec → PR
- Code review and refactoring
- Bug fixing with full repo context
- Reads `CLAUDE.md` for project-specific instructions

## CLI
```bash
claude                           # Interactive REPL
claude "fix the auth bug"        # One-shot task
claude --model opus              # Specify model
claude --allowedTools bash,read,write  # Restrict tools
claude --resume                  # Resume last session
```

## Key Files
- **CLAUDE.md** — Project instructions (root of repo)
- **~/.claude/settings.json** — Global config
- **~/.claude/projects/** — Per-project memory

## Integration with Hermes
- Hermes skill: `claude-code` — delegates coding tasks
- Spawned via `delegate_task` or direct terminal
- Best for: multi-file features, large refactors, PR workflows

## Cost
Uses Anthropic API (Claude Sonnet/Opus) — included in ~$50/mo usage
