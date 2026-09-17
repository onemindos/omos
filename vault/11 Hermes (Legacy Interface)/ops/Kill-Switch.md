# Kill Switch — Agent Process Manager

**Location:** `~/.local/bin/kill-agents`
**Created:** 2026-05-21
**Domain:** SO/11 Hermes (Legacy Interface)

---

## What It Does

Emergency kill switch for all AI agent processes running on omos-mbp. Supports targeted or nuclear shutdown of hermes, openclaw gateway, kanban workers, and MCP servers.

---

## Usage

```bash
kill-agents           # list all running agent processes (safe, no kill)
kill-agents all       # nuclear: SIGTERM everything, then SIGKILL if needed
kill-agents workers   # kanban/background workers only
kill-agents gateway   # openclaw or hermes gateway processes only
kill-agents <name>    # kill by partial name match (e.g. "codex", "hermes", "mcp")
```

---

## What Gets Killed

| Mode | Targets |
|------|---------|
| `all` | hermes CLI, hermes gateway, openclaw gateway, hermes dashboard, all MCP servers |
| `workers` | hermes kanban worker processes only |
| `gateway` | openclaw gateway (port 18789) + hermes gateway (port varies) |
| `<name>` | any process whose cmdline contains the string |

Shutdown sequence: SIGTERM → 3s wait → SIGKILL if still alive.

---

## Currently Running (as of 2026-05-21)

| PID | Process | Notes |
|-----|---------|-------|
| 87120 | hermes CLI | terminal s002 |
| 82569 | hermes CLI | terminal s003 |
| 87342 | hermes CLI | terminal s006 |
| 6388 | hermes gateway | background, port varies |
| 17851 | hermes dashboard | port 9119 |
| 46596 | openclaw gateway | port 18789 |
| — | MCP servers | context7, filesystem, todoist, homeassistant, wwv |

---

## Related

- [[Hermes-Overview]] — full agent architecture
- [[Hermes-API-Server]] — gateway + API server config
- Tailscale fleet: `tailscale status` — verify omos-mbp (100.98.178.23) is online

---

## Notes

- `kill-agents list` is always safe — it only reads, never kills
- Gateway restart: `hermes gateway run` (foreground) or as a launchd service
- Dashboard restart: `hermes dashboard --tui --no-open --host 127.0.0.1 --port 9119`
- MCP servers auto-restart when hermes CLI starts
