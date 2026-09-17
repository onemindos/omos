# VPS Services & Sync Architecture

> **VPS**: 134.122.8.135 (DigitalOcean, Ubuntu 24.04.4 LTS)
> **Last updated**: 2026-03-06
> **Monorepo**: `/home/onemind/onemind-os/` → `github.com/One-Mind-OS/onemindos-v0.4`

---

## Architecture Overview

```
┌─────────────────────────────────────────────────────────────────────────┐
│                        MAC (Obsidian + VS Code)                         │
│  /Users/zeuslegacy/Documents/onemindos-v0.4/                           │
│  Obsidian-git plugin auto-commits + pushes (1–5 min)                   │
└──────────────────────────┬──────────────────────────────────────────────┘
                           │ git push (HTTPS)
                    ┌──────▼──────┐
                    │   GitHub    │  One-Mind-OS/onemindos-v0.4
                    └──────┬──────┘
                           │ git pull (every 1 min, codex-pull.timer)
┌──────────────────────────▼──────────────────────────────────────────────┐
│                        VPS (134.122.8.135)                              │
│  /home/onemind/onemind-os/                                             │
│                                                                         │
│  openclaw/workspaces/legacy/ ←──symlink── ~/.openclaw/workspace        │
│  codex-watch.service: inotifywait → auto commit+push (3s debounce)     │
│  openclaw-reload.service: inotifywait → auto restart on skill changes  │
└─────────────────────────────────────────────────────────────────────────┘
```

**Data flow**:
1. **Mac → GitHub**: Obsidian-git commits + pushes within ~1–5 minutes
2. **GitHub → VPS**: `codex-pull.timer` runs `git pull --rebase --autostash` every 1 minute
3. **VPS → GitHub**: `codex-watch.service` detects file changes via inotifywait, commits + pushes within ~3 seconds
4. **Skill/config → OpenClaw restart**: `openclaw-reload.service` watches `skills/`, `lobster/`, `openclaw.json` and auto-restarts gateway

---

## Systemd Services

### 1. `openclaw.service` — Main Gateway

Legacy's AI runtime. Reads workspace MDs every turn (hot-reload), skills on startup (cold).

| Field | Value |
|-------|-------|
| **Status** | `active (running)` |
| **User** | onemind |
| **WorkingDirectory** | `/home/onemind/.openclaw` |
| **ExecStart** | `/usr/bin/openclaw gateway` |
| **Restart** | always (10s delay, max 5 in 300s) |
| **Memory limit** | 4 GB |
| **Env file** | `/opt/onemind/.env` |

```ini
[Unit]
Description=OneMind OS — OpenClaw Gateway (Legacy's Runtime)
After=network-online.target docker.service
Wants=network-online.target

[Service]
Type=simple
User=onemind
Group=onemind
WorkingDirectory=/home/onemind/.openclaw
EnvironmentFile=-/opt/onemind/.env
ExecStart=/usr/bin/openclaw gateway
Restart=always
RestartSec=10
StartLimitBurst=5
StartLimitIntervalSec=300
LimitNOFILE=65536
MemoryMax=4G
StandardOutput=journal
StandardError=journal
SyslogIdentifier=openclaw

[Install]
WantedBy=multi-user.target
```

**Key paths**:
- Config: `/home/onemind/.openclaw/openclaw.json`
- Workspace: `/home/onemind/.openclaw/workspace` → symlink to `/home/onemind/onemind-os/openclaw/workspaces/legacy`
- Logs: `/tmp/openclaw-1000/openclaw-2026-*.log` + journalctl
- Canvas: served at `http://127.0.0.1:18789/__openclaw__/canvas/`
- Browser control: `http://127.0.0.1:18791/`

**Connected channels**: Discord (OneMind OS bot), BlueBubbles (iMessage)

---

### 2. `codex-pull.timer` + `codex-pull.service` — Git Pull (GitHub → VPS)

Picks up changes pushed from Mac (Obsidian-git) every minute.

| Field | Value |
|-------|-------|
| **Status** | `active (waiting)` |
| **Interval** | Every 1 min (30s after boot) |
| **User** | onemind |
| **WorkingDirectory** | `/home/onemind/onemind-os` |
| **Command** | `git pull --rebase --autostash -q` |

```ini
# /etc/systemd/system/codex-pull.timer
[Unit]
Description=OneMind git pull timer (every 1 min)

[Timer]
OnBootSec=30s
OnUnitActiveSec=1min
Persistent=true

[Install]
WantedBy=timers.target
```

```ini
# /etc/systemd/system/codex-pull.service
[Unit]
Description=OneMind git pull (picks up Mac Obsidian edits)

[Service]
Type=oneshot
User=onemind
Group=onemind
Environment=HOME=/home/onemind
WorkingDirectory=/home/onemind/onemind-os
ExecStart=/usr/bin/git pull --rebase --autostash -q
```

---

### 3. `codex-watch.service` — Real-time Git Sync (VPS → GitHub)

Uses inotifywait to detect changes in `codex/` and `openclaw/workspaces/`, then commits + pushes within ~3 seconds.

| Field | Value |
|-------|-------|
| **Status** | `active (running)` |
| **User** | onemind |
| **Script** | `/home/onemind/codex-watch.sh` |
| **Watched dirs** | `$REPO/codex/`, `$REPO/openclaw/workspaces/` |
| **Debounce** | 3 seconds |
| **Commit author** | `LegacyxAI <zeus.ops@legacyhq.ai>` |

```bash
#!/bin/bash
# /home/onemind/codex-watch.sh
REPO=/home/onemind/onemind-os
LOCK=/tmp/codex-watch.lock

cd "$REPO" || exit 1
git config user.email "zeus.ops@legacyhq.ai"
git config user.name "LegacyxAI"

inotifywait -m -r \
  -e close_write -e create -e delete -e move \
  --exclude '\.git' \
  "$REPO/codex/" "$REPO/openclaw/workspaces/" |
while read -r dir event file; do
  sleep 3
  while read -r -t 0.1 _ 2>/dev/null; do :; done

  exec 9>"$LOCK"
  flock -n 9 || continue

  cd "$REPO"
  git pull --rebase --autostash -q 2>/dev/null
  git add codex/ openclaw/workspaces/ 2>/dev/null
  if ! git diff --cached --quiet 2>/dev/null; then
    git commit -m "vps: auto-sync $(date -u +%Y-%m-%dT%H:%M:%SZ)" -q
    git push -q 2>&1 | head -3
    echo "[$(date)] Pushed changes from: ${dir}${file}"
  fi
  exec 9>&-
done
```

---

### 4. `openclaw-reload.service` — Auto-restart on Skill/Config Changes

Watches skill and lobster directories + `openclaw.json`. When files change, restarts the gateway so new skills/configs are loaded.

| Field | Value |
|-------|-------|
| **Status** | `active (running)` |
| **User** | onemind |
| **Script** | `/home/onemind/openclaw-reload.sh` |
| **Watched paths** | `openclaw/workspaces/legacy/skills/`, `openclaw/workspaces/legacy/lobster/`, `~/.openclaw/openclaw.json` |
| **Cooldown** | 10 seconds |
| **Debounce** | 3 seconds |
| **Requires** | `openclaw.service` |

```bash
#!/bin/bash
# /home/onemind/openclaw-reload.sh
WORKSPACE=/home/onemind/onemind-os/openclaw/workspaces/legacy
CONFIG=/home/onemind/.openclaw/openclaw.json
COOLDOWN=10
last_restart=0

inotifywait -m -r \
  -e close_write -e create -e delete -e move \
  --exclude "node_modules|\.git|__pycache__" \
  "$WORKSPACE/skills/" "$WORKSPACE/lobster/" "$CONFIG" |
while read -r dir event file; do
  now=$(date +%s)
  elapsed=$((now - last_restart))
  if [ "$elapsed" -lt "$COOLDOWN" ]; then continue; fi

  sleep 3
  while read -r -t 0.1 _ 2>/dev/null; do :; done

  echo "[$(date -u +%Y-%m-%dT%H:%M:%SZ)] Change: ${dir}${file} (${event}) — restarting"
  sudo systemctl restart openclaw
  last_restart=$(date +%s)
  echo "[$(date -u +%Y-%m-%dT%H:%M:%SZ)] OpenClaw restarted"
done
```

---

## Cron Jobs

### `onemind` user crontab

```cron
# workspace-backup (monorepo tarball every 15 min)
*/15 * * * * /home/onemind/workspace-backup.sh
```

**workspace-backup.sh** — Safety net that catches any uncommitted changes missed by codex-watch:

```bash
#!/bin/bash
# /home/onemind/workspace-backup.sh
REPO=/home/onemind/onemind-os
cd "$REPO" || exit 1
git pull --rebase --autostash -q 2>/dev/null
git add codex/ openclaw/workspaces/ 2>/dev/null
if ! git diff --cached --quiet 2>/dev/null; then
  git commit -m "backup: $(date -u +%Y-%m-%dT%H:%M:%SZ)" -q
  git push -q >> /home/onemind/workspace-backup.log 2>&1 || true
fi
```

### `root` crontab

None.

---

## Key Symlink

```
/home/onemind/.openclaw/workspace → /home/onemind/onemind-os/openclaw/workspaces/legacy
```

This means OpenClaw reads its workspace files directly from the monorepo. Any git pull automatically updates what Legacy sees, and any changes Legacy makes are automatically tracked by codex-watch.service.

**Backup**: `/home/onemind/.openclaw/workspace-backup-20260306/` (original standalone workspace before migration)

---

## Git Authentication

VPS uses HTTPS + GitHub CLI credential helper:

```bash
gh auth git-credential   # handles auth for git push/pull
```

No SSH keys for git — authentication is via `gh auth login` token.

---

## OpenClaw Hot vs Cold Reload

| What changed | Reload needed? | How |
|-------------|---------------|-----|
| AGENTS.md, SOUL.md, USER.md, etc. | **No** — hot reload | Read every turn automatically |
| `openclaw.json` | **Yes** — cold | `openclaw-reload.service` auto-restarts |
| `skills/*` | **Yes** — cold | `openclaw-reload.service` auto-restarts |
| `lobster/*` | **Yes** — cold | `openclaw-reload.service` auto-restarts |
| `canvas/*` | **No** — served live | Static files served by canvas host |

---

## Monitoring Commands

```bash
# Service status
systemctl status openclaw
systemctl status codex-watch
systemctl status codex-pull.timer
systemctl status openclaw-reload

# OpenClaw logs (live)
journalctl -u openclaw -f

# Reload watcher logs
journalctl -u openclaw-reload -f

# Git sync logs
journalctl -u codex-watch -f

# Recent sync activity
journalctl -u codex-watch --since "1 hour ago" --no-pager

# Check symlink
ls -la /home/onemind/.openclaw/workspace

# Check backup log
tail -20 /home/onemind/workspace-backup.log
```

---

## Troubleshooting

| Problem | Check | Fix |
|---------|-------|-----|
| Legacy not responding | `systemctl status openclaw` | `systemctl restart openclaw` |
| Edits from Mac not showing | `systemctl status codex-pull.timer` | `systemctl start codex-pull.timer` |
| VPS edits not pushing | `systemctl status codex-watch` | `systemctl restart codex-watch` |
| Skills not updating | `systemctl status openclaw-reload` | `systemctl restart openclaw-reload` |
| Git merge conflict | Check `journalctl -u codex-pull` | `cd /home/onemind/onemind-os && git rebase --abort && git pull` |
| Symlink broken | `ls -la ~/.openclaw/workspace` | `ln -sfn /home/onemind/onemind-os/workspaces/legacy ~/.openclaw/workspace` |
