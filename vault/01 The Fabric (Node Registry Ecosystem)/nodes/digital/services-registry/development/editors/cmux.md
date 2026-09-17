---
title: "cmux"
type: reference
status: active
tags: [platform, development, terminal, multiplexer]
---
# cmux

> Terminal multiplexer for managing multiple coding agent sessions. Orchestrate Claude Code, Codex, OpenCode simultaneously.

## Use Case in OneMind OS
- Run multiple coding agents in parallel panes
- Monitor agent output side-by-side
- Session persistence (detach/reattach)

## CLI
```bash
cmux                    # Launch
cmux new -s dev         # Named session
cmux attach -t dev      # Reattach
cmux ls                 # List sessions
```
