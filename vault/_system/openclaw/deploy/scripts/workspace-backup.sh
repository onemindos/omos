#!/bin/bash
# workspace-backup.sh — Safety net git sync (cron every 15 min)
# Catches any changes missed by codex-watch.service.
# Deployed to: /home/onemind/workspace-backup.sh
# Crontab: */15 * * * * /home/onemind/workspace-backup.sh

REPO=/home/onemind/onemind-os
cd "$REPO" || exit 1
git pull --rebase --autostash -q 2>/dev/null
git add codex/ openclaw/workspaces/ 2>/dev/null
if ! git diff --cached --quiet 2>/dev/null; then
  git commit -m "backup: $(date -u +%Y-%m-%dT%H:%M:%SZ)" -q
  git push -q >> /home/onemind/workspace-backup.log 2>&1 || true
fi
