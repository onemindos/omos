#!/bin/bash
# workspace-backup.sh — Safety net git sync (cron every 15 min)
# Catches any changes missed by onemind-watch.service.
# Deployed to: /home/ubuntu/workspace-backup.sh
# Crontab: */15 * * * * /home/ubuntu/workspace-backup.sh

VAULT=/home/ubuntu/onemind
cd "$VAULT" || exit 1
git pull --rebase --autostash -q 2>/dev/null
git add -A 2>/dev/null
if ! git diff --cached --quiet 2>/dev/null; then
  git commit -m "backup: $(date -u +%Y-%m-%dT%H:%M:%SZ)" -q
  git push -q >> /home/ubuntu/workspace-backup.log 2>&1 || true
fi
