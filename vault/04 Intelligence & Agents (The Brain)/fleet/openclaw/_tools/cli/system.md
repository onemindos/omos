---
name: OpenClaw CLI — System
updated: 2026-05-11
---

# `openclaw system`

Manages system-level events, heartbeat scheduling, and presence monitoring.

**Official docs:** [docs.openclaw.ai/cli/system](https://docs.openclaw.ai/cli/system.md)

---

## Commands

```bash
# System events
openclaw system events                        # List recent system events
openclaw system events --type heartbeat       # Filter by event type
openclaw system events --since 1h            # Events in last hour
openclaw system events --json                 # Structured output

# Enqueue a notification
openclaw system enqueue --type alert --message "Disk usage high"
openclaw system enqueue --agent legacy --message "Run health check" --deliver

# Heartbeat management
openclaw system heartbeat                     # Show current heartbeat status
openclaw system heartbeat --enable            # Enable periodic heartbeat
openclaw system heartbeat --disable           # Disable heartbeat
openclaw system heartbeat --interval 60       # Set interval in seconds

# Presence monitoring
openclaw system presence                      # Show which agents are active
openclaw system presence --agent legacy       # Single agent presence
openclaw system presence --json
```

---

## Event Types

| Type | Description |
|---|---|
| `heartbeat` | Periodic liveness signal from gateway |
| `alert` | Operator-triggered notification |
| `agent-start` | Agent session began |
| `agent-stop` | Agent session ended |
| `channel-connect` | Channel plugin connected |
| `channel-disconnect` | Channel plugin disconnected |

---

## Notes

- Heartbeat events are visible in logs and Discord if a channel binding is configured
- `enqueue` with `--deliver` routes the message through the configured reply channel
- Presence tracks active session counts per agent, not uptime of the gateway process itself
