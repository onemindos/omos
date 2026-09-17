#!/usr/bin/env bash
# Stop all Legacy ambient presence daemons
set -euo pipefail

LOG_DIR="$HOME/.openclaw/logs"
PID_FILE="$LOG_DIR/legacy.pids"

if [[ ! -f "$PID_FILE" ]]; then
  echo "[stop] No PID file found at $PID_FILE"
  exit 0
fi

source "$PID_FILE"

echo "[stop] Stopping Legacy stack..."

for name_pid in "legacy-pulse:$PULSE_PID" "legacy-eyes:$EYES_PID" "gateway:$GATEWAY_PID"; do
  name="${name_pid%%:*}"
  pid="${name_pid##*:}"
  if kill -0 "$pid" 2>/dev/null; then
    kill "$pid"
    echo "[stop] Stopped $name (PID $pid)"
  else
    echo "[stop] $name (PID $pid) already stopped"
  fi
done

rm -f "$PID_FILE"
echo "[stop] Done."
