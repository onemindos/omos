---
name: OpenClaw CLI — Channels
updated: 2026-05-11
---

# `openclaw channels`

Manages chat channel accounts and their runtime status on the Gateway.

**Official docs:** [docs.openclaw.ai/cli/channels](https://docs.openclaw.ai/cli/channels.md)

---

## Commands

```bash
# List & Status
openclaw channels list                         # Configured accounts + status tags
openclaw channels status --probe               # Live health checks on reachable gateways
openclaw channels status --json

# Capabilities
openclaw channels capabilities --channel discord --target channel:<id>  # Permission audit
openclaw channels resolve --channel telegram --target @username         # Name → ID

# Account Management
openclaw channels add --channel telegram       # Add new account (interactive)
openclaw channels remove --channel telegram --account <id>

# Authentication
openclaw channels login --channel telegram
openclaw channels logout --channel telegram --account <id>
```

---

## Status Tags

| Tag | Meaning |
|---|---|
| `installed` | Channel plugin is installed |
| `configured` | Account credentials are configured |
| `enabled` | Channel is active and running |

---

## Important

> Do NOT use `openclaw sessions` as a channel health signal. Sessions shows stored conversation rows — not live connectivity. A connected channel can be healthy with zero session rows.
> Use `openclaw channels status --probe` for actual live health.

---

## Auth Note

`channels login/logout` must run **from the gateway host** — they don't work via remote `agent exec`.

---

## Our Channels

| Channel | Account | Notes |
|---|---|---|
| Telegram | `onemindos_bot` | Group `-1003905828404`, 5 forum topics (Legacy, Zenith, Unity, Eternal, Sovereignty) |
| Discord | `@OneMind OS` | Server `1453411552181883027`, 6 agent channels + operational categories |

Full map: [`channels/CHANNELS.md`](../channels/CHANNELS.md)
