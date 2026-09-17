---
title: "Active Memory Plugin"
type: reference
status: active
created: 2026-04-12
updated: 2026-04-12
owner: zeus
tags: [type/reference, domain/so, so/01-intelligence, openclaw, memory]
ai_summary: "Configuration reference for OpenClaw's Active Memory plugin — a blocking memory sub-agent that automatically surfaces relevant vault context before each reply."
---

# Active Memory Plugin

> Enabled: **2026-04-12** on OpenClaw `2026.4.11`
> Docs: https://docs.openclaw.ai/concepts/active-memory

## What It Does

Active Memory runs a bounded blocking memory sub-agent **before** the main reply on every eligible chat turn. It automatically pulls relevant context, preferences, and past details from vault memory — without requiring you to say "search memory" or explicitly trigger recall.

The sub-agent can only use `memory_search` and `memory_get`. It returns either a compact user-fact summary or `NONE`. This summary is injected as hidden system context for the main model.

**Why it matters for One Mind:** Legacy and all 9 agents now proactively surface your habits, preferences, recurring facts, and vault knowledge — without you having to prompt them to look.

---

## Current Setup (as of 2026-04-12)

Enabled via `openclaw plugins enable active-memory` — running with OpenClaw defaults:

```json
{
  "plugins": {
    "entries": {
      "active-memory": {
        "enabled": true,
        "config": {
          "agents": ["main"],
          "allowedChatTypes": ["direct"],
          "modelFallback": "google/gemini-3-flash",
          "queryMode": "recent",
          "promptStyle": "balanced",
          "timeoutMs": 15000,
          "maxSummaryChars": 220,
          "persistTranscripts": false,
          "logging": true
        }
      }
    }
  }
}
```

---

## Query Modes (`queryMode`)

Controls how much conversation context the memory sub-agent sees.

| Mode | What It Sees | Best For | Recommended Timeout |
|------|-------------|---------|-------------------|
| `message` | Latest user message only | Fastest; stable preference recall; when follow-ups don't need context | 3–5s |
| `recent` | Latest message + last few turns **(default)** | Best balance of speed and conversational grounding | ~15s |
| `full` | Entire conversation history | Strongest recall; topics mentioned early in a long thread | 15s+ |

**Rule:** `message` < `recent` < `full` in both latency and recall depth.

---

## Prompt Styles (`promptStyle`)

Controls how eager or strict the sub-agent is when deciding whether to surface memory.

| Style | Behavior | Default For |
|-------|---------|------------|
| `balanced` | General-purpose default | `recent` mode |
| `strict` | Least eager — returns NONE unless very confident | `message` mode |
| `contextual` | Most continuity-friendly; conversation history matters more | `full` mode |
| `recall-heavy` | Surfaces softer/looser matches | — |
| `precision-heavy` | Aggressively returns NONE unless match is obvious | — |
| `preference-only` | **Optimized for favorites, habits, routines, recurring personal facts** | Recommended for One Mind |

> **Recommendation:** Switch to `preference-only` to make Legacy surface personal facts (your habits, food preferences, decision patterns) more naturally.

---

## Advanced Options

### Prompt Overrides
- `promptAppend` — adds extra instructions after the default prompt before conversation context
  ```json
  "promptAppend": "Prefer stable long-term preferences over one-off events."
  ```
- `promptOverride` — replaces the entire default prompt (not recommended unless deliberately testing a different recall contract)

### Thinking Level (`thinking`)
Overrides the memory sub-agent thinking depth. Default is `"off"` — do not enable without reason as it directly adds user-visible latency.
Options: `off` | `minimal` | `low` | `medium` | `high` | `xhigh` | `adaptive`

### Transcript Persistence
By default, sub-agent transcripts are temporary (written to temp dir, deleted after run). To keep them for debugging:
```json
"persistTranscripts": true,
"transcriptDir": "active-memory"
```
Stored at: `agents/<agent>/sessions/active-memory/<session-id>.jsonl`

---

## In-Session Commands

```
/active-memory status        # Show current session state
/active-memory off           # Pause for this session only
/active-memory on            # Resume for this session

/active-memory status --global   # Global state
/active-memory off --global      # Pause all sessions (writes config)
/active-memory on --global       # Resume all sessions (writes config)

/verbose on                  # See active memory debug output inline
```

With `/verbose on`, replies include lines like:
```
🧩 Active Memory: ok 842ms recent 34 chars
🔎 Active Memory Debug: Lemon pepper wings with blue cheese.
```

---

## Where It Runs

| Surface | Active Memory? |
|---------|---------------|
| Control UI / webchat persistent sessions | ✅ Yes |
| Other interactive channel sessions (Telegram, etc.) | ✅ Yes |
| Headless one-shot runs | ❌ No |
| Heartbeat/background runs | ❌ No |
| Sub-agent / internal helper runs | ❌ No |

---

## Debugging

If active memory isn't triggering:
1. Confirm `plugins.entries.active-memory.enabled: true`
2. Confirm current agent id is in `config.agents`
3. Confirm you're in an interactive persistent chat session
4. Turn on `config.logging: true` and watch gateway logs
5. Run `openclaw memory status --deep` to verify memory search works

If hits are too noisy → tighten `maxSummaryChars`
If too slow → lower `queryMode`, reduce `timeoutMs`, reduce recent turn counts

---

## Related

- [[NATIVE-MEMORY-ARCHITECTURE]] — How vault memory search works
- [[VPS-GATEWAY]] — OpenClaw gateway setup and services
- [[OPENCLAW_CAPABILITIES]] — Full OpenClaw feature inventory
