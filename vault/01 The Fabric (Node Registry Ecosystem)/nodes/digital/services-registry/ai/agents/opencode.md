---
title: "OpenCode"
type: reference
status: active
tags: [platform, ai, agent, coding, open-source]
---
# OpenCode

> Open-source terminal coding agent. Model-agnostic — works with any provider.

## Use Case in OneMind OS
- Coding agent that works with any LLM provider
- Reads `AGENTS.md` for project-specific instructions
- Good for tasks where you want a specific model (local, OpenRouter, etc.)

## CLI
```bash
opencode                         # Interactive TUI
opencode --provider openrouter   # Specify provider
opencode --model deepseek/deepseek-coder  # Any model
```

## Key Files
- **AGENTS.md** — Project instructions (root of repo)
- **~/.opencode/config.json** — Global config

## Integration with Hermes
- Hermes skill: `opencode` — delegates coding tasks
- Flexible provider choice — use for cost optimization or model-specific strengths

## Cost
Depends on provider — can use free/local models
