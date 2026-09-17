#!/bin/bash
# codex-watch.sh — Real-time git sync (VPS → GitHub)
# Watches codex/ and openclaw/workspaces/ for changes, commits + pushes within ~3s.
# Deployed to: /home/onemind/codex-watch.sh
# Service: codex-watch.service

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
  # Drain queued events (3s debounce)
  while read -r -t 0.1 _ 2>/dev/null; do :; done

  exec 9>"$LOCK"
  flock -n 9 || continue

  cd "$REPO"

  # Guard: recover from detached HEAD (e.g., after failed rebase)
  BRANCH=$(git symbolic-ref --short HEAD 2>/dev/null)
  if [[ -z "$BRANCH" ]]; then
    echo "[$(date)] WARNING: detached HEAD detected, checking out main"
    git checkout main 2>/dev/null
  fi

  git pull --rebase --autostash -q 2>/dev/null
  git add codex/ openclaw/workspaces/ 2>/dev/null
  if ! git diff --cached --quiet 2>/dev/null; then
    CHANGED=$(git diff --cached --name-only | head -5 | tr '\n' ', ')
    git commit -m "agent: update ${CHANGED%,}" -q
    git push -q 2>&1 | head -3
    echo "[$(date)] Pushed changes: ${CHANGED%,}"
  fi
  exec 9>&-
done
