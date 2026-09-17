---
name: OpenClaw CLI — Update
updated: 2026-05-11
---

# `openclaw update`

Updates OpenClaw and manages release channels.

**Official docs:** [docs.openclaw.ai/cli/update](https://docs.openclaw.ai/cli/update.md)

---

## Commands

```bash
openclaw update                          # Update to latest on current channel
openclaw update --channel stable         # Switch to stable channel + update
openclaw update --channel beta           # Switch to beta channel + update
openclaw update --channel dev            # Switch to dev channel (git checkout)
openclaw update --dry-run                # Preview without changes
openclaw update --no-restart             # Skip gateway service restart
openclaw update status                   # Show active channel + availability
openclaw update status --json
```

---

## Release Channels

| Channel | Description |
|---|---|
| `stable` | Production releases |
| `beta` | Pre-release features |
| `dev` | Git checkout, bleeding edge |

Channel setting persists in configuration.

---

## Update Workflow (git checkout)

1. Verify clean worktree
2. Switch to selected channel/tag
3. Fetch upstream (dev only)
4. Preflight build check (walks back up to 10 commits if needed)
5. Install dependencies
6. Build Control UI + run doctor checks
7. Sync plugins to active channel

---

## Notes

- Beta channel tries `@beta` plugin releases first, falls back to latest
- Nix mode (`OPENCLAW_NIX_MODE=1`) disables mutating updates
- Package-manager updates bypass restart cooldown to prevent old process lazy-loading
- Post-update plugin warnings don't fail the overall update — they're reported separately
