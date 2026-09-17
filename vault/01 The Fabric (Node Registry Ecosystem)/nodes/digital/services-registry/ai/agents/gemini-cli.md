---
title: "Gemini CLI"
type: reference
status: active
tags: [platform, ai, agent, coding, google]
---
# Gemini CLI

> Google's coding agent. Large context window (1M+ tokens), good for codebase-wide analysis.

## Use Case in OneMind OS
- Large codebase analysis (massive context window)
- Reads `GEMINI.md` for project-specific instructions
- Alternative perspective on architecture decisions

## CLI
```bash
gemini                           # Interactive
gemini "analyze this codebase"   # One-shot
```

## Key Files
- **GEMINI.md** — Project instructions (root of repo)

## Integration with Hermes
- Can be spawned via terminal for second-opinion analysis
- Best for: large repo understanding, documentation generation, architecture review

## Cost
Google AI API — free tier available, Pro ~$20/mo
