---
name: OpenClaw CLI — Gateway
updated: 2026-05-11
---

# `openclaw gateway`

Controls the Gateway WebSocket server that manages channels, nodes, sessions, and hooks.

**Official docs:** [docs.openclaw.ai/cli/gateway](https://docs.openclaw.ai/cli/gateway.md)

---

## Commands

### Start / Run
```bash
openclaw gateway          # Start the gateway
openclaw gateway run      # Explicit start
openclaw gateway run --allow-unconfigured   # Dev/testing mode
```

### Service Management (systemd on Linux)
```bash
openclaw gateway install   # Install as system service
openclaw gateway start     # Start service
openclaw gateway stop      # Stop service
openclaw gateway restart   # Restart service
openclaw gateway restart --safe   # Wait for active work to drain first
openclaw gateway restart --force  # Override immediately
openclaw gateway uninstall # Remove service
```

### Monitoring & Diagnostics
```bash
openclaw gateway health         # Liveness/readiness via HTTP endpoints
openclaw gateway status         # Service state + connectivity
openclaw gateway probe          # Comprehensive debug of all reachable gateways
openclaw gateway probe --ssh user@host  # Remote SSH probe
openclaw gateway stability      # Recent diagnostic events + memory/perf metrics
openclaw gateway diagnostics export    # Privacy-respecting bug report zip
```

### Discovery
```bash
openclaw gateway discover       # Find beacons via mDNS (local.) or DNS-SD (wide-area)
```

---

## Key Behaviors

**Safe restart:** `--safe` waits for active replies, embedded runs, and task runs to drain before restarting. Use `--force` to override.

**Diagnostics privacy:** Records event names, counts, byte sizes, memory readings, queue/session state, channel/plugin names, and redacted session summaries. Does NOT include chat text, webhook bodies, tool outputs, tokens, secrets, or raw session IDs.

**Config requirement:** Requires `gateway.mode=local` in `~/.openclaw/openclaw.json`. Use `--allow-unconfigured` for dev runs.

---

## Our Setup

Gateway runs as a systemd service on omos-core and omos-sentinel.

```bash
# On VPS
ssh ubuntu@100.123.126.34
openclaw gateway status
openclaw gateway restart --safe
openclaw gateway health --verbose
```

Service files: [`install (Install & Deploy)/systemd/`](../install%20(Install%20%26%20Deploy)/systemd%20(Systemd%20Services)/)

> See also: [`daemon.md`](daemon.md) — legacy alias for gateway service commands
