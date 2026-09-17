---
name: OpenClaw CLI — Status
updated: 2026-05-10
---

# `openclaw status`

Full system diagnostics for channels, sessions, memory, and security.

**Official docs:** [docs.openclaw.ai/cli/status](https://docs.openclaw.ai/cli/status.md)

---

## Commands

```bash
openclaw status             # Fast read-only diagnostics
openclaw status --all       # Comprehensive: security + memory checks
openclaw status --deep      # Live probes across all messaging platforms
openclaw status --usage     # Remaining provider quota as percentages
```

---

## Key Output Sections

- **Gateway:** uptime, version, update channel
- **Execution:** sandbox path (e.g. docker)
- **Runtime:** actual runtime environment
- **Channels:** connection status per platform
- **Agents:** per-agent session store status
- **Memory:** index health (skipped unless `--all`)
- **Token tracking:** backfills from transcript usage logs when snapshots lack data

---

## Notes

- Plain `status` avoids heavy probes — use `--all` or `--deep` for intensive checks
- `--deep` runs live probes against Telegram, Discord, Slack, Signal, WhatsApp
- Secret credentials degrade gracefully — shows warnings rather than failing
- Per-agent session stores shown when multiple agents are active
