---
name: OpenClaw CLI — TUI
updated: 2026-05-11
---

# `openclaw tui`

Terminal UI for interacting with the Gateway or running in local embedded mode.

**Official docs:** [docs.openclaw.ai/cli/tui](https://docs.openclaw.ai/cli/tui.md)

---

## Commands

```bash
openclaw tui                                    # Connect to local gateway
openclaw tui --url ws://127.0.0.1:18789 --token <token>  # Explicit gateway
openclaw tui --session main --deliver           # Target session + deliver responses
openclaw chat                                   # Alias for openclaw tui --local
openclaw terminal                               # Alias for openclaw tui --local
```

---

## Modes

| Mode | Description |
|---|---|
| **Gateway mode** (default) | Connects to a running gateway |
| **Local mode** (`--local`) | Embedded agent runtime, no gateway required |

---

## Auto-Detection

When launched from an agent workspace directory, the TUI automatically selects that agent unless you explicitly specify `--session`.

---

## Local Mode Features

- Extra `/auth [provider]` command surface
- Good for config repair without a running gateway

In-TUI diagnostics:
```
!openclaw config validate
!openclaw doctor
```

---

## Constraint

Local mode cannot be combined with `--url`, `--token`, or `--password`.
