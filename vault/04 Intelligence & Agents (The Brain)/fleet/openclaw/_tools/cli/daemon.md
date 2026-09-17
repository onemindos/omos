---
name: OpenClaw CLI — Daemon (Legacy)
updated: 2026-05-11
---

# `openclaw daemon`

Legacy alias for gateway service management. Use [`openclaw gateway`](gateway.md) instead.

**Official docs:** [docs.openclaw.ai/cli/daemon](https://docs.openclaw.ai/cli/daemon.md)

---

## Commands (legacy aliases for `gateway`)

```bash
openclaw daemon status     # → openclaw gateway status
openclaw daemon install    # → openclaw gateway install
openclaw daemon start      # → openclaw gateway start
openclaw daemon stop       # → openclaw gateway stop
openclaw daemon restart    # → openclaw gateway restart
openclaw daemon restart --safe   # Wait for active work to drain
openclaw daemon uninstall  # → openclaw gateway uninstall
```

---

## Maintained for backward compatibility only

Use [`gateway.md`](gateway.md) for all current documentation and examples.
