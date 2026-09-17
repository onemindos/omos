---
name: OpenClaw CLI — Agent (singular)
updated: 2026-05-11
---

# `openclaw agent`

Executes a single agent turn through the Gateway (or locally with `--local`).

**Official docs:** [docs.openclaw.ai/cli/agent](https://docs.openclaw.ai/cli/agent.md)

> Note: This is `openclaw agent` (singular — runs a turn). For managing agent workspaces, see [`agents.md`](agents.md).

---

## Commands

```bash
# Basic usage (requires at least one session selector + --message)
openclaw agent --to +15555550123 --message "status update" --deliver
openclaw agent --session-id <id> --message "Summarize logs"
openclaw agent --agent ops --message "What's the status?"

# With model override
openclaw agent --agent ops --model openai/gpt-5.4 --message "Summarize logs"

# With thinking level
openclaw agent --agent main --message "..." --thinking max

# Local (embedded, no gateway)
openclaw agent --local --message "Hello"

# JSON output (diagnostics → stderr, response → stdout)
openclaw agent --agent main --message "Hello" --json
```

---

## Session Selectors (at least one required)

| Flag | Description |
|---|---|
| `--to <dest>` | Recipient-derived session key |
| `--session-id <id>` | Explicit session ID |
| `--agent <id>` | Target a specific configured agent |

---

## Flags

| Flag | Description |
|---|---|
| `-m / --message <text>` | The agent prompt **(required)** |
| `--model <id>` | Override model selection |
| `--thinking <level>` | Reasoning depth: `off`, `low`, `medium`, `high`, `max` |
| `--local` | Embedded agent (no gateway required) |
| `--timeout <seconds>` | Agent timeout |
| `--deliver` | Send reply to channel/target |
| `--reply-channel` | Customize delivery channel |
| `--reply-to` | Customize delivery recipient |
| `--json` | Structured JSON to stdout |

---

## Behavior

- Gateway requests fall back to embedded execution on failure unless `--local` forces immediate local
- With `--json`: diagnostics go to stderr, response to stdout (pipeline-safe)
