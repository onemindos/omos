---
name: OpenClaw CLI — Sessions
updated: 2026-05-11
---

# `openclaw sessions`

Manages conversation session storage across messaging platforms.

**Official docs:** [docs.openclaw.ai/cli/sessions](https://docs.openclaw.ai/cli/sessions.md)

---

## Commands

```bash
openclaw sessions                         # List newest 100 sessions
openclaw sessions --limit 50             # Custom limit
openclaw sessions --limit all            # All sessions
openclaw sessions --agent <id>           # Scope to specific agent
openclaw sessions --json                 # Machine-readable output

# Export trajectory bundle for a session
openclaw sessions export --session-id <id>   # Outputs to .openclaw/trajectory-exports/
```

### Cleanup / Maintenance
```bash
openclaw sessions cleanup                # Run housekeeping per session.maintenance config
openclaw sessions cleanup --dry-run      # Preview without modifying
openclaw sessions cleanup --enforce      # Override config maintenance mode
openclaw sessions cleanup --fix-dm-scope # Retire obsolete DM entries from legacy routing
openclaw sessions cleanup --active-key <key>  # Protect a session from eviction
```

---

## What Cleanup Removes

- Orphaned transcripts
- Compaction checkpoints
- Trajectory sidecars exceeding `pruneAfter` threshold
- Preserves referenced files

---

## Important Note

**Do not use `sessions` as a channel health signal.** It shows stored conversation rows — not live channel connectivity. A connected account can be healthy with zero session rows. For live health: use `openclaw channels status --probe` or `openclaw health --verbose`.
