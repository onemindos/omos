---
title: "OpenClaw — Workspace Backup Repository"
type: reference
status: active
updated: 2026-05-11
owner: zeus
tags: [type/reference, domain/SO-02, openclaw/workspaces]
---

# Workspace Backup Repository

All agent workspace files are automatically version-controlled and pushed to GitHub.

## Repository

| Field | Value |
| --- | --- |
| **Repo** | `github.com/OneMind-OS/openclaw-core-workspaces` |
| **Branch** | `main` |
| **Tracked files** | All `workspace*/` and `agents/*/workspace/` files (7 files × 23 agents) |
| **Watcher script** | `~/.openclaw/watch.sh` |
| **Systemd service** | `onemind-watch.service` |
| **Debounce** | 3 seconds (inotifywait) |
| **Commit message format** | `watch: {EVENT} {FILE} — {ISO_TIMESTAMP}` |
| **Watcher log** | `/tmp/openclaw-backup.log` |

## How It Works

1. `watch.sh` runs under `onemind-watch.service` watching all `workspace*/` and `agents/*/workspace/` directories via `inotifywait`
2. On any file change (close_write, create, delete, move): waits 3s debounce, drains queued events
3. Runs `git add workspace*/` then commits + pushes to `origin/main`
4. Logs every push to `/tmp/openclaw-backup.log`

## Verify Health

```bash
# Check watcher is running
ssh ubuntu@100.123.126.34 "ps aux | grep watch.sh | grep -v grep"

# Check last push timestamp
ssh ubuntu@100.123.126.34 "git -C ~/.openclaw log --oneline -5"

# Check sync status (should show 'up to date')
ssh ubuntu@100.123.126.34 "git -C ~/.openclaw fetch origin && git -C ~/.openclaw status"

# View live watcher log
ssh ubuntu@100.123.126.34 "tail -20 /tmp/openclaw-backup.log"
```

## Fix Broken Push (Diverged Repo)

When `watch.sh` can't push (happens after gateway restarts that modify config files):

```bash
ssh ubuntu@100.123.126.34
cd ~/.openclaw

# Step 1: See what diverged
git fetch origin
git log HEAD..origin/main --oneline    # what remote has that we don't
git log origin/main..HEAD --oneline    # what we have that remote doesn't

# Step 2: Rebase local commits on top of remote
git stash push -u -m "pre-rebase save"
git pull --rebase origin main
git stash pop

# Step 3: Push and verify
git push origin main
git status    # should show 'up to date'
```

## Restore a Workspace from Backup

```bash
ssh ubuntu@100.123.126.34

# Clone the repo to a temp location
git clone https://github.com/OneMind-OS/openclaw-core-workspaces.git /tmp/ws-restore

# Restore a specific agent (e.g. zenith — workspace-zenith/ path)
cp /tmp/ws-restore/workspace-zenith/* ~/.openclaw/workspace-zenith/

# For agents/* path agents (e.g. apex)
cp /tmp/ws-restore/workspace-apex/* ~/.openclaw/agents/apex/workspace/ 2>/dev/null || true

# Restart gateway to reload
sudo systemctl restart openclaw.service

## cleanup
rm -rf /tmp/ws-restore
```

## What's NOT Tracked

The backup repo tracks workspace files only. These are NOT tracked (by design):

- `openclaw.json` — tracked separately by `openclaw-config-sync.service`
- `credentials/` — never version-controlled
- `cron/` — job state files, runtime only
- `logs/`, `media/`, `completions/` — runtime data

## Agent Workspace Path Reference

| Agent | Repo path | VM path |
| --- | --- | --- |
| Legacy (main) | `workspace/` | `~/.openclaw/workspace/` |
| Zenith | `workspace-zenith/` | `~/.openclaw/workspace-zenith/` |
| Unity | `workspace-unity/` | `~/.openclaw/workspace-unity/` |
| Eternal | `workspace-eternal/` | `~/.openclaw/workspace-eternal/` |
| Sovereignty | `workspace-sovereignty/` | `~/.openclaw/workspace-sovereignty/` |
| Spartan | `workspace-spartan/` | `~/.openclaw/workspace-spartan/` |
| Oracle | `workspace-oracle/` | `~/.openclaw/workspace-oracle/` |
| Eden | `workspace-eden/` | `~/.openclaw/workspace-eden/` |
| Haven | `workspace-haven/` | `~/.openclaw/workspace-haven/` |
| Empire | `workspace-empire/` | `~/.openclaw/workspace-empire/` |
| Guardian | `workspace-guardian/` | `~/.openclaw/workspace-guardian/` |
| Grid | `workspace-grid/` | `~/.openclaw/workspace-grid/` |
| Trinity | `workspace-trinity/` | `~/.openclaw/workspace-trinity/` |
| Overwatch | `workspace-overwatch/` | `~/.openclaw/workspace-overwatch/` |
| Librarian | `workspace-librarian/` | `~/.openclaw/workspace-librarian/` |
| Sentinel | `agents/sentinel/workspace/` | `~/.openclaw/agents/sentinel/workspace/` |
| Pulse | `agents/pulse/workspace/` | `~/.openclaw/agents/pulse/workspace/` |
| Genesis | `agents/genesis/workspace/` | `~/.openclaw/agents/genesis/workspace/` |
| Helix | `agents/helix/workspace/` | `~/.openclaw/agents/helix/workspace/` |
| Nomad | `agents/nomad/workspace/` | `~/.openclaw/agents/nomad/workspace/` |
| Apex | `agents/apex/workspace/` | `~/.openclaw/agents/apex/workspace/` |
| Creed | `agents/creed/workspace/` | `~/.openclaw/agents/creed/workspace/` |
| EagleEye | `agents/eagleeye/workspace/` | `~/.openclaw/agents/eagleeye/workspace/` |

## Related

- [README.md](README.md) — Workspace overview
- [[../agents/README|Agent Registry]] — All 23 agents
- [[../help/README|Help & Runbooks]] — More operational runbooks
