---
id: so-02-sync-services
title: "onemind Vault Sync Services"
type: runbook
status: active
node: onemindos-core
created: 2026-05-10
modified: 2026-05-10
owner: zeus
tags: [domain/so, infra/vps, infra/git, infra/sync]
---

# onemind Vault Sync Services

Three systemd units keep the vault in bidirectional sync between the VPS and GitHub. All run as the `ubuntu` user on `onemindos-core`.

```
Mac (Obsidian / Claude Code)
    ↓ git push
GitHub (OneMind-OS/onemind-codex) ← SOURCE OF TRUTH
    ↓ onemind-pull.timer (every 10s)
    ↑ onemind-watch.service (inotifywait → git commit + push)
/home/ubuntu/onemind/ (VPS)
```

---

## Services

### onemind-pull.timer
Fires `onemind-pull.service` every 10 seconds after boot.

```ini
[Unit]
Description=OneMind Pull Timer (every 10s)
After=network.target

[Timer]
OnBootSec=30s
OnUnitActiveSec=10s
Unit=onemind-pull.service

[Install]
WantedBy=timers.target
```

### onemind-pull.service
One-shot service: runs `git pull --rebase` on the vault.

```ini
[Unit]
Description=OneMind Auto-pull from GitHub
After=network.target

[Service]
Type=oneshot
WorkingDirectory=/home/ubuntu/onemind
ExecStart=/usr/bin/git -C /home/ubuntu/onemind pull --rebase
User=ubuntu
StandardOutput=journal
StandardError=journal
```

### onemind-watch.service
Always-running inotifywait loop. Watches for any file change in the vault, then auto-commits and pushes to GitHub within ~30s.

```ini
[Unit]
Description=OneMind Watch + Auto-commit
After=network.target

[Service]
Type=simple
WorkingDirectory=/home/ubuntu/onemind
ExecStart=/usr/local/bin/onemind-watch.sh
User=ubuntu
Restart=always
RestartSec=10
StandardOutput=journal
StandardError=journal

[Install]
WantedBy=multi-user.target
```

**Watch script** at `/usr/local/bin/onemind-watch.sh`:

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
    echo "[onemind-watch] Committed and pushed: $CHANGED"
  fi
done
```

---

## File Locations

| File | Path |
|------|------|
| Pull service | `/etc/systemd/system/onemind-pull.service` |
| Pull timer | `/etc/systemd/system/onemind-pull.timer` |
| Watch service | `/etc/systemd/system/onemind-watch.service` |
| Watch script | `/usr/local/bin/onemind-watch.sh` |
| Vault root | `/home/ubuntu/onemind/` |

---

## Common Commands

```bash
# Status check
systemctl status onemind-pull.timer onemind-watch.service

# View logs
journalctl -u onemind-watch.service -f
journalctl -u onemind-pull.service --since "10 min ago"

# Restart watch after manual vault work
systemctl restart onemind-watch.service

# Force immediate pull
systemctl start onemind-pull.service

# Check timer schedule
systemctl list-timers onemind-pull.timer
```

---

## Notes

- Pull runs as `ubuntu`, vault lives at `/home/ubuntu/onemind/`
- Watch script commits use `agent:` prefix — distinguishable in git log from Zeus/Claude commits
- If pull + watch conflict (simultaneous commit + pull), `--rebase` resolves it cleanly
- Deploy key for VPS→GitHub push: `~/.ssh/codex_deploy` (GitHub deploy key ID 148743097)
