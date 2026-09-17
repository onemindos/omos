---
type: quick-reference
domain: SO/02
tags: [openclaw, lobster, agents, commands]
created: 2026-05-12
---

# Agent Task Workflow — Quick Reference

## Run Workflow

```bash
# Zenith (HP)
lobster(action="run", pipeline="~/.openclaw/workflows/agent-tasks.lobster", argsJson='{"agent_id":"zenith"}')

# Unity (LE)
lobster(action="run", pipeline="~/.openclaw/workflows/agent-tasks.lobster", argsJson='{"agent_id":"unity"}')

# Eternal (GE)
lobster(action="run", pipeline="~/.openclaw/workflows/agent-tasks.lobster", argsJson='{"agent_id":"eternal"}')

# Sovereignty (SO)
lobster(action="run", pipeline="~/.openclaw/workflows/agent-tasks.lobster", argsJson='{"agent_id":"sovereignty"}')
```

## Test (Mock Data)

```bash
bash ~/.openclaw/workflows/test-agent-workflow.sh zenith
```

## Schedule Daily

```bash
openclaw cron add \
  --name "zenith-tasks" \
  --cron "0 7 * * *" \
  --tz "America/New_York" \
  --session session:zenith \
  --message "Run agent-tasks.lobster for zenith"
```

## Files

- Workflow: `~/.openclaw/workflows/agent-tasks.lobster`
- Config: `~/.openclaw/workflows/agent-config.json`
- Test: `~/.openclaw/workflows/test-agent-workflow.sh`
- Docs: `~/.openclaw/workflows/README-agent-tasks.md`

## Update Project IDs

```bash
vim ~/.openclaw/workflows/agent-config.json
```

## Architecture

**One workflow, all agents** — only difference is `agent_id` arg
