---
name: OpenClaw — Thinking Levels
updated: 2026-05-11
---

# Thinking Levels

Per-message reasoning depth control. Match compute cost to task complexity.

**Official docs:** [docs.openclaw.ai/tools/thinking](https://docs.openclaw.ai/tools/thinking.md)

---

## Directives

```bash
/think <level>
/t <level>
/think:<level>
/thinking <level>
```

---

## Levels

| Level | Alias | Description |
|---|---|---|
| `off` | — | No reasoning |
| `minimal` | `think` | Minimal reasoning |
| `low` | `think hard` | Low reasoning |
| `medium` | `think harder` | Medium reasoning |
| `high` | `ultrathink` | Max budget (most models) |
| `xhigh` | `ultrathink+` | GPT-5.2+ / Claude Opus 4.7 effort |
| `adaptive` | — | Provider-managed adaptive (Claude 4.6 default) |
| `max` | — | Provider max reasoning (Opus 4.7, Ollama max) |

---

## Resolution Order

1. Inline directive on the message (one-shot)
2. Session override (`/think:medium` alone sets session default)
3. Per-agent config (`agents.list[].thinkingDefault`)
4. Global config (`agents.defaults.thinkingDefault`)
5. Provider default (Claude 4.6 → `adaptive`; others → `medium` or nearest supported)

---

## Session Control

```bash
/think:high          # Set session default to high
/think off           # Disable for this session
/think default       # Clear session override, inherit config/provider default
/think               # Show current level (no argument)
```

---

## Provider Mapping

| Provider | Notes |
|---|---|
| Claude 4.6 (Anthropic) | Default = `adaptive`. `/think xhigh` = extended thinking + xhigh effort |
| Claude Opus 4.7 | `/think xhigh` = adaptive + xhigh effort. `/think max` = provider max |
| DeepSeek V4 | `xhigh`/`max` → `reasoning_effort: "max"`. Lower levels → `high` |
| Google Gemini | `adaptive` → dynamic thinking. Others map to nearest `thinkingLevel` |
| Ollama | `max` → native `think: "high"`. Supports `low`/`medium`/`high`/`max` |
| OpenAI GPT | Maps to Responses API `reasoning.effort`. `/think off` sends `effort: "none"` |
| Z.AI | Binary only: any non-`off` level = `on` (mapped to `low`) |

---

## One Mind Usage

Since all agents run `claude-sonnet-4-6`, thinking defaults to `adaptive`. Practical guide:

| Task type | Recommended level |
|---|---|
| Routine cron dispatch, status checks | `off` |
| Standard task execution | `adaptive` (default) |
| Complex multi-step planning | `high` |
| Strategic analysis, deep research | `max` |

Per-agent override example (in `openclaw.json`):
```json
{
  "agents": {
    "list": [
      { "id": "legacy", "thinkingDefault": "adaptive" },
      { "id": "sovereignty", "thinkingDefault": "high" }
    ]
  }
}
```

---

## Fast Mode

```bash
/fast on       # Priority processing tier
/fast off      # Standard tier
/fast default  # Clear session override
/fast          # Show current state
```

For Anthropic: `on` = `service_tier=auto`, `off` = `service_tier=standard_only`.

---

## Related

- [cli/Models.md](../cli/Models.md)
- [tools/lobster.md](lobster.md) — thinking level applies to llm-task steps inside Lobster
