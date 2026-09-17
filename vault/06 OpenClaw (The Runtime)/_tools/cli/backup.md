---
name: OpenClaw CLI — Backup
updated: 2026-05-11
---

# `openclaw backup`

Creates compressed archives of the OpenClaw installation state, configuration, credentials, and workspaces.

**Official docs:** [docs.openclaw.ai/cli/backup](https://docs.openclaw.ai/cli/backup.md)

---

## Commands

```bash
openclaw backup create                          # Full backup (timestamped .tar.gz)
openclaw backup create --only-config            # Config only (no workspaces)
openclaw backup create --no-include-workspace   # Skip workspace dirs
openclaw backup create --verify                 # Verify archive integrity after creation
openclaw backup verify <archive-path>           # Verify an existing archive
```

---

## What Gets Backed Up

| Content | Included |
|---|---|
| State directory (`~/.openclaw`) | Yes |
| Active config file | Yes |
| Credentials directory (if external) | Yes |
| Workspace directories | Yes (use `--no-include-workspace` to skip) |
| Model auth profiles | Yes (part of `agents/<id>/agent/auth-profiles.json`) |
| Plugin manifests | Yes |
| Plugin `node_modules` | No (rebuildable via `openclaw plugins update`) |

---

## Verification

`--verify` checks the manifest and validates all declared files exist in the archive.

---

## Our Backup Snapshots

Config snapshots are stored in [`install (Install & Deploy)/backups/`](../install%20(Install%20%26%20Deploy)/backups/):
- `omos-core-openclaw-2026-05-04.json`
- `omos-sentinel-openclaw-2026-05-04.json`

See also: [`install (Install & Deploy)/scripts/workspace-backup.sh`](../install%20(Install%20%26%20Deploy)/scripts%20(Deploy%20Scripts)/workspace-backup.sh)
