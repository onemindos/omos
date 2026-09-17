---
name: OpenClaw CLI — Logs
updated: 2026-05-10
---

# `openclaw logs`

Tails Gateway file logs over RPC.

**Official docs:** [docs.openclaw.ai/cli/logs](https://docs.openclaw.ai/cli/logs.md)

---

## Commands

```bash
openclaw logs                          # Last 200 lines (default)
openclaw logs --limit 500              # Custom line count
openclaw logs --follow                 # Stream continuously
openclaw logs --follow --interval 500  # Poll every 500ms (default 1000ms)
openclaw logs --json                   # Line-delimited JSON output
openclaw logs --local-time             # Timestamps in local timezone
openclaw logs --plain                  # No color formatting

# Custom gateway target
openclaw logs --url ws://127.0.0.1:18789 --token "$TOKEN"
```

---

## Gateway RPC Flags

| Flag | Description |
|---|---|
| `--url` | Gateway WebSocket URL |
| `--token` | Auth token |
| `--timeout` | Connection timeout |
| `--expect-final` | Wait for final response |

---

## Smart Reconnection

When `--follow` is active and the connection drops, automatic reconnection kicks in with exponential backoff — up to 8 retries, capped at 30s between attempts. Falls back to local loopback gateway when no `--url` is specified.

---

## Quick Reference

```bash
# Live tail from VPS
ssh ubuntu@100.123.126.34 'openclaw logs --follow'

# Last 100 lines as JSON
openclaw logs --limit 100 --json
```
