#!/bin/bash
# openclaw-config-sync.sh — auto-commit OpenClaw config changes to GitHub
# Mirrors /home/ubuntu/.openclaw/ → local workspace repo clone → git push
# Deployed to: /home/ubuntu/openclaw-config-sync.sh
# Service: openclaw-config-sync.service
#
# SETUP: The workspace repo must already be cloned at $REPO_DIR
# and the remote must be set to the correct GitHub repo.

set -euo pipefail

OPENCLAW_DIR=/home/ubuntu/.openclaw
REPO_DIR=/home/ubuntu/openclaw-workspaces
COOLDOWN=10
last_commit=0

if [ ! -d "$REPO_DIR/.git" ]; then
  echo "[openclaw-config-sync] ERROR: $REPO_DIR is not a git repo. Clone the workspace repo first."
  exit 1
fi

echo "[openclaw-config-sync] Watching $OPENCLAW_DIR for changes..."

inotifywait -m -r \
  -e close_write -e create -e delete -e move \
  --exclude '(\.git|node_modules|\.session|\.log|__pycache__|\.tmp)' \
  "$OPENCLAW_DIR" |
while read -r _dir _event _file; do
  now=$(date +%s)
  elapsed=$((now - last_commit))
  if [ "$elapsed" -lt "$COOLDOWN" ]; then continue; fi

  sleep 3
  while IFS= read -r -t 0.1 _; do :; done

  rsync -a --delete \
    --exclude '.git' \
    --exclude '*.log' \
    --exclude '*.session' \
    --exclude 'node_modules' \
    "$OPENCLAW_DIR/" "$REPO_DIR/"

  cd "$REPO_DIR"
  if git diff --quiet && git diff --cached --quiet; then
    continue
  fi

  git add -A
  git commit -m "config: auto-sync $(date -u +%Y-%m-%dT%H:%M:%SZ)"
  git push origin main
  last_commit=$(date +%s)
  echo "[openclaw-config-sync] Committed and pushed at $(date -u +%Y-%m-%dT%H:%M:%SZ)"
done
