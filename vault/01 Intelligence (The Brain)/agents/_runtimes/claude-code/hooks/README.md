---
type: index
domain: SO-01
status: active
updated: 2026-05-14
owner: zeus
tags: [runtime/claude-code, domain/so, type/index]
ai_summary: "Claude Code lifecycle hooks — PreToolUse, PostToolUse, Stop, Notification. Event-driven automation wired to Claude Code sessions."
---

# Claude Code — Hooks

Lifecycle hooks that fire during Claude Code sessions. Defined in `.claude/settings.json` under the `hooks` key.

## Hook Types

| Hook | When It Fires | Common Uses |
| ---- | ------------- | ----------- |
| `PreToolUse` | Before any tool call | Security validation, input sanitization, rate limiting |
| `PostToolUse` | After any tool call | Logging, audit trail, output capture |
| `Stop` | When Claude finishes a turn | Session logging, auto-commit to vault, notifications |
| `Notification` | On permission requests or alerts | Telegram/Slack alerts, approval routing |

## Config Location

```
~/.claude/settings.json        ← global (all projects)
<project>/.claude/settings.json ← project-scoped
```

## Hook Definition Format

```json
{
  "hooks": {
    "Stop": [
      {
        "matcher": "",
        "hooks": [
          {
            "type": "command",
            "command": "~/.claude/hooks/session-log.sh"
          }
        ]
      }
    ]
  }
}
```

## Active Hooks

| Hook | Trigger | Script | Purpose |
| ---- | ------- | ------ | ------- |
| *(none yet — add as wired)* | — | — | — |

## Hook Scripts Location

Scripts live at `~/.claude/hooks/` on Mac. Document them here when created.

---
*System hooks reference: `_system/hooks/` in vault root. Claude Code hooks are a subset of the broader One Mind hook system.*
