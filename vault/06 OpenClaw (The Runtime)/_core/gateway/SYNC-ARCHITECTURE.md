---
title: "OneMind — Sync Architecture"
type: reference
status: active
created: 2026-05-11
modified: 2026-05-11
owner: zeus
tags: [type/reference, domain/SO-02, infra/sync, infra/git]
---

# Sync Architecture

Two independent git-based sync systems keep One Mind operational. Both run on omos-core as `ubuntu`.

---

## System 1: Vault Sync (One Mind ↔ GitHub)

Keeps the Obsidian vault in bidirectional sync between Mac, GitHub, and VPS.

```
Mac (Obsidian / Claude Code)
    ↓ git push (manual or Obsidian Git plugin)
GitHub: OneMind-OS/One-Mind  ← SOURCE OF TRUTH
    ↓ onemind-pull.timer (every 10s)
/home/ubuntu/onemind/ on VPS
    ↓ onemind-watch.service (inotifywait → commit → push)
GitHub ← updated
    ↓ Mac: git pull (manual or Obsidian Git plugin)
Mac ← updated
```

### Services

| Service | Type | Does |
|---------|------|------|
| `onemind-pull.timer` | Timer (10s) | Triggers pull service |
| `onemind-pull.service` | Oneshot | `git -C /home/ubuntu/onemind pull --rebase` |
| `onemind-watch.service` | Always-on | Watches vault, auto-commits + pushes changes |

### Auth

The VPS authenticates to GitHub via a **token embedded in the HTTPS remote URL**:

```
origin https://gho_XXXX@github.com/OneMind-OS/One-Mind.git
```

This is a GitHub Organization token (`gho_` prefix). If sync stops working, the token may have expired.

**To check:** `git -C /home/ubuntu/onemind remote get-url origin`  
**To update:** `git -C /home/ubuntu/onemind remote set-url origin https://NEW_TOKEN@github.com/OneMind-OS/One-Mind.git`

### Files on VPS

| What | Path |
|------|------|
| Pull timer unit | `/etc/systemd/system/onemind-pull.timer` |
| Pull service unit | `/etc/systemd/system/onemind-pull.service` |
| Watch service unit | `/etc/systemd/system/onemind-watch.service` |
| Watch script | `/usr/local/bin/onemind-watch.sh` |
| Vault directory | `/home/ubuntu/onemind/` |
| Git user | `LegacyxAI <legacy@onemind.os>` |

### Watch Script Source

```bash
#!/bin/bash
set -e
VAULT=/home/ubuntu/onemind

while true; do
  inotifywait -r -e modify,create,delete,move --timeout 30 "$VAULT" \
    --exclude '.git' 2>/dev/null || true

  if git -C "$VAULT" status --porcelain 2>/dev/null | grep -q .; then
    cd "$VAULT"
    git add -A
    CHANGED=$(git diff --cached --name-only 2>/dev/null | head -5 | tr '\n' ' ' | sed 's/ $//')
    git commit -m "agent: ${CHANGED:-update}" 2>/dev/null || true
    git push 2>&1 || true
  fi
done
```

Commit messages from the VPS use `agent:` prefix — distinguishable from Zeus/Claude commits.

---

## System 2: OpenClaw Workspace Sync (Workspaces ↔ GitHub)

Keeps all agent workspace files (`IDENTITY.md`, `SOUL.md`, etc.) backed up to a separate repo.

```
Agent writes to ~/.openclaw/workspace-*/
    ↓ watch.sh (inotifywait, 3s debounce)
    ↓ git commit + push
GitHub: OneMind-OS/openclaw-core-workspaces
```

### Service

| What | Path |
|------|------|
| Watch script | `/home/ubuntu/.openclaw/watch.sh` |
| Process | Background bash (not systemd) — started manually or at boot |
| Repo root | `/home/ubuntu/.openclaw/` |
| Remote | `https://github.com/OneMind-OS/openclaw-core-workspaces.git` |

### Auth

Same pattern — GitHub token in HTTPS URL. Separate token from the vault sync.

**To check:** `git -C ~/.openclaw remote get-url origin`

### What's Tracked

- All `workspace*/` directories (15 agent workspaces, 7 files each)
- `watch.sh`, `backup.sh`
- `hooks/`, `plugin-skills/` (select files)
- `openclaw.json` (config snapshot)
- `wiki/`

See [[../workspaces (Agent Workspaces)/backup-repo|backup-repo.md]] for full details and health checks.

---

## Key Differences Between the Two Systems

| | Vault Sync | Workspace Sync |
|---|---|---|
| **Repo** | `OneMind-OS/One-Mind` | `OneMind-OS/openclaw-core-workspaces` |
| **Direction** | Bidirectional (Mac ↔ VPS) | One-way (VPS → GitHub) |
| **Pull** | Every 10s via timer | None (push only) |
| **Push** | On any vault file change | On workspace file change |
| **Mac involved?** | Yes — Mac pushes, VPS pulls | No — VPS only |
| **Managed by** | systemd services | Background bash process |
| **Commit prefix** | `agent:` | `watch:` |

---

## Troubleshooting

### Mac → VPS not syncing

```bash
# Step 1: Did Mac push?
git -C ~/Documents/onemind status
# If "ahead of origin" → push first:
git -C ~/Documents/onemind push origin main

# Step 2: Is VPS pull timer running?
ssh omos-core "systemctl status onemind-pull.timer"
# Should show "active (waiting)" — if dead:
ssh omos-core "systemctl start onemind-pull.timer"

# Step 3: Force immediate pull
ssh omos-core "systemctl start onemind-pull.service"

# Step 4: Verify
ssh omos-core "git -C /home/ubuntu/onemind log --oneline -1"
```

### VPS → Mac not syncing

```bash
# Step 1: Did VPS push?
ssh omos-core "git -C /home/ubuntu/onemind status"
# If changes unpushed, watch may be down:
ssh omos-core "systemctl status onemind-watch"
# Restart if needed:
ssh omos-core "systemctl restart onemind-watch"

# Step 2: Pull on Mac
git -C ~/Documents/onemind pull
```

### Workspace backup not pushing

```bash
# Check watch.sh process
ssh omos-core "ps aux | grep watch.sh | grep -v grep"
# If dead, restart:
ssh omos-core "nohup /home/ubuntu/.openclaw/watch.sh &"

# Check for diverged state
ssh omos-core "git -C ~/.openclaw status"
# If diverged, fix:
ssh omos-core "cd ~/.openclaw && git stash push -u -m save && git pull --rebase origin main && git stash pop && git push"
```

### Token expired (sync silently stops)

```bash
# Check if push/pull is failing
ssh omos-core "journalctl -u onemind-pull --since '5 min ago' | tail -5"
# Look for "Authentication failed" or "403"

# Update token (get new one from GitHub → Settings → Tokens)
ssh omos-core "git -C /home/ubuntu/onemind remote set-url origin https://NEW_TOKEN@github.com/OneMind-OS/One-Mind.git"

# Same for workspace repo if needed
ssh omos-core "git -C ~/.openclaw remote set-url origin https://NEW_TOKEN@github.com/OneMind-OS/openclaw-core-workspaces.git"
```

### Merge conflict on pull

```bash
ssh omos-core
cd /home/ubuntu/onemind

# See what's conflicting
git status

# Most conflicts: accept VPS version (agent wrote it)
git checkout --theirs .
git add -A
git commit -m "resolve: accept VPS changes"
git push

# OR accept Mac version
git checkout --ours .
git add -A
git commit -m "resolve: accept Mac changes"
git push
```

### Nuclear option (full resync from GitHub)

```bash
# WARNING: destroys any uncommitted VPS changes
ssh omos-core
cd /home/ubuntu
mv onemind onemind.bak
git clone https://TOKEN@github.com/OneMind-OS/One-Mind.git onemind
systemctl restart onemind-watch onemind-pull.timer
```

---

## Obsidian Git Plugin Settings (Mac)

For seamless Mac → VPS flow, ensure these are configured in Obsidian:

| Setting | Value |
|---------|-------|
| Auto commit interval | 5 minutes |
| **Auto push after commit** | **ON** (critical — without this, VPS never sees changes) |
| Auto pull interval | 5 minutes |
| Pull on startup | ON |
| Push on close | ON |

Without auto-push enabled, every commit stays local on Mac and the VPS will never receive it until you manually run `git push`.

---

## Health Dashboard Commands

Run these from Mac to check everything at once:

```bash
# Full health check
echo "=== MAC ===" && \
git -C ~/Documents/onemind status | head -3 && \
echo "=== VPS VAULT ===" && \
ssh omos-core "systemctl is-active onemind-pull.timer onemind-watch && git -C /home/ubuntu/onemind log --oneline -1" && \
echo "=== VPS WORKSPACES ===" && \
ssh omos-core "ps aux | grep '[w]atch.sh' | wc -l && git -C ~/.openclaw log --oneline -1"
```

Expected output:
- MAC: "up to date" or "nothing to commit"
- VPS VAULT: "active active" + recent commit
- VPS WORKSPACES: "2" (watch processes) + recent commit

---

## Related

- [[VPS-SERVICES|VPS-SERVICES]] — Full service inventory and architecture
- [[../workspaces (Agent Workspaces)/backup-repo|Backup Repo]] — Workspace backup details
- [[DISASTER-RECOVERY|Disaster Recovery]] — Full VPS rebuild runbook
