#!/bin/bash
# onemind-watch.sh — Real-time git sync (VPS → GitHub)
# Watches the One Mind vault for changes, commits + pushes within ~3s.
# Deployed to: /usr/local/bin/onemind-watch.sh
# Service: onemind-watch.service
#
# NOTE: This file is DEPRECATED — the live script is at /usr/local/bin/onemind-watch.sh on VPS.
# See gateway/SYNC-ARCHITECTURE.md for the current watch script source.

VAULT=/home/ubuntu/onemind
LOCK=/tmp/onemind-watch.lock

cd "$VAULT" || exit 1
git config user.email "legacy@onemind.os"
git config user.name "LegacyxAI"

inotifywait -m -r \
  -e close_write -e create -e delete -e move \
  --exclude '\.git' \
  "$VAULT" |
while read -r dir event file; do
  sleep 3
  # Drain queued events (3s debounce)
  while read -r -t 0.1 _ 2>/dev/null; do :; done

  exec 9>"$LOCK"
  flock -n 9 || continue

  cd "$VAULT"

  # Guard: recover from detached HEAD (e.g., after failed rebase)
  BRANCH=$(git symbolic-ref --short HEAD 2>/dev/null)
  if [[ -z "$BRANCH" ]]; then
    echo "[$(date)] WARNING: detached HEAD detected, checking out main"
    git checkout main 2>/dev/null
  fi

  git add -A 2>/dev/null
  if ! git diff --cached --quiet 2>/dev/null; then
    CHANGED=$(git diff --cached --name-only | head -5 | tr '\n' ' ' | sed 's/ $//')
    git commit -m "agent: ${CHANGED:-update}" -q
    git push -q 2>&1 | head -3
    echo "[$(date)] Pushed changes: ${CHANGED}"
  fi
  exec 9>&-
done
