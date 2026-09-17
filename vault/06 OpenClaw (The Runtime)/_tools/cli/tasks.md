---
name: OpenClaw CLI — Tasks
updated: 2026-05-11
---

# `openclaw tasks`

Inspects durable background tasks and Task Flow state.

**Official docs:** [docs.openclaw.ai/cli/tasks](https://docs.openclaw.ai/cli/tasks.md)

---

## Commands

### Listing & Inspection
```bash
openclaw tasks                               # List all background tasks
openclaw tasks --runtime <type>             # Filter by runtime type
openclaw tasks --status <status>            # Filter by status
openclaw tasks --json                       # Machine-readable output
openclaw tasks show <id>                    # Show task by ID
openclaw tasks show --run-id <id>           # Show by run ID
openclaw tasks show --session-key <key>     # Show by session key
```

### Management
```bash
openclaw tasks cancel <id>                  # Cancel a running task
openclaw tasks notify --id <id> [--policy]  # Modify notification policy
```

### Maintenance & Auditing
```bash
openclaw tasks audit        # Surface stale, lost, delivery-failed, or inconsistent records
openclaw tasks maintenance  # Reconcile and clean up task records
```

### Task Flow
```bash
openclaw tasks flow list                    # List Task Flow state
openclaw tasks flow show <id>               # Inspect Task Flow
openclaw tasks flow cancel <id>             # Cancel Task Flow
```

---

## Runtime Types

| Type | Description |
|---|---|
| `subagent` | Spawned sub-agent tasks |
| `acp` | ACP protocol tasks |
| `cron` | Scheduled cron tasks |
| `cli` | CLI-initiated tasks |

## Status Values

`queued` · `running` · `succeeded` · `failed` · `timed_out` · `cancelled` · `lost`

---

## Notes

- `audit` handles cron tasks specially — checks persisted run logs before marking tasks lost
- `maintenance` performs reconciliation across all runtime types
