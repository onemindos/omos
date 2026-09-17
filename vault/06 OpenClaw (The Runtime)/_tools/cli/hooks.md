---
name: OpenClaw CLI — Hooks
updated: 2026-05-10
---

# `openclaw hooks`

Manages event-driven automations that trigger on actions like `/new`, `/reset`, and gateway startup.

**Official docs:** [docs.openclaw.ai/cli/hooks](https://docs.openclaw.ai/cli/hooks.md)

---

## Commands

```bash
openclaw hooks list               # All discovered hooks + status
openclaw hooks list --verbose     # Include missing requirements
openclaw hooks list --json        # Machine-readable output

openclaw hooks info <name>        # Details: source, handler path, triggered events
openclaw hooks check              # Summary: ready vs unavailable count

openclaw hooks enable <name>      # Enable a hook in config
openclaw hooks disable <name>     # Disable a hook in config
```

> After enabling or disabling, restart the gateway for changes to take effect.

---

## Bundled Hooks

| Hook | Trigger | What it does |
|---|---|---|
| `session-memory` | `/new` or `/reset` | Saves context before session reset |
| `bootstrap-extra-files` | Agent init | Injects extra files during bootstrap |
| `command-logger` | Every command | Writes all command events to audit log |
| `boot-md` | Gateway startup | Executes `BOOT.md` |

---

## Installing Hook Packs

```bash
# From npm
openclaw plugins install <pack-name>

# From local path
openclaw plugins install -l /path/to/pack

# With pinned version
openclaw plugins install <pack-name> --pin

# Update installed packs
openclaw plugins update <id>
openclaw plugins update --all
```

> Note: Plugin-managed hooks cannot be toggled via `hooks enable/disable` — manage the plugin instead.

---

## Our Hooks

Hooks live in `_system/openclaw/hooks/`:

| Hook | File | Purpose |
|---|---|---|
| `onemind-session-logger` | `handler.js` | Logs sessions to vault memory |
| Stop handler | `stop.js` | Graceful stop |

Install guide: [`_system/openclaw/hooks/INSTALL.md`](../../../../_system/openclaw/hooks/INSTALL.md)

---

## Related

- [automation/README.md](../automation/README.md) — Cron + hooks overview
- [Official hooks concept](https://docs.openclaw.ai/automation/hooks.md)
- [Plugin hooks SDK](https://docs.openclaw.ai/plugins/hooks.md)
