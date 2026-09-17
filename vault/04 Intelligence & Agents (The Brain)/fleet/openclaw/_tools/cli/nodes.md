---
name: OpenClaw CLI — Nodes
updated: 2026-05-11
---

# `openclaw nodes`

Manages paired devices and invokes node capabilities remotely.

**Official docs:** [docs.openclaw.ai/cli/nodes](https://docs.openclaw.ai/cli/nodes.md)

---

## Commands

### Listing & Status
```bash
openclaw nodes list                            # All pending + paired nodes
openclaw nodes list --connected                # Only currently active nodes
openclaw nodes list --last-connected 24h       # Connected within timeframe
openclaw nodes list --json
```

### Pairing Management
```bash
openclaw nodes pending                         # List awaiting approval requests
openclaw nodes approve <requestId>             # Authorize a node
openclaw nodes reject <requestId>              # Deny a pairing request
openclaw nodes remove --node <id|name|ip>      # Delete obsolete pairing
```

### Node Operations
```bash
openclaw nodes rename --node <id|name|ip> --name <displayName>

# Invoke a node capability
openclaw nodes invoke \
  --node <id|name|ip> \
  --command <command> \
  --params '{"key":"value"}' \
  --invoke-timeout 15000 \
  --idempotency-key <key>
```

---

## Notes

- `system.run` and shell commands are NOT available via `invoke` — use the `exec` tool with `host=node` instead
- Approval scope requirements vary by command type (basic pairing = minimal, admin commands = elevated)
- `--idempotency-key` enables request deduplication

---

## Related

- [`node.md`](00-24%20SO%20(Sovereign%20Operations)/01%20Intelligence%20&%20Agents%20(The%20Brain)/fleet/openclaw/_tools/cli/node.md) — singular `openclaw node` for deploying headless node agents
- [`devices.md`](devices.md) — device pairing and token management
