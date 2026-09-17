#!/bin/bash
# openclaw-reload.sh — Auto-restart gateway on skill/config changes
# Watches skills/, lobster/, and openclaw.json. Restarts after 10s cooldown.
# Deployed to: /home/onemind/openclaw-reload.sh
# Service: openclaw-reload.service

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
  # Drain queued events (3s debounce)
  while read -r -t 0.1 _ 2>/dev/null; do :; done

  echo "[$(date -u +%Y-%m-%dT%H:%M:%SZ)] Change: ${dir}${file} (${event}) — restarting"
  sudo systemctl restart openclaw
  last_restart=$(date +%s)
  echo "[$(date -u +%Y-%m-%dT%H:%M:%SZ)] OpenClaw restarted"
done
