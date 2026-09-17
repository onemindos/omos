#!/usr/bin/env bash
# ─────────────────────────────────────────────────────────────────────────────
# Legacy Stack Launcher
# Brings up the full ambient presence stack alongside OpenClaw gateway
#
# Usage:
#   ./launch.sh                         # home, webcam
#   ./launch.sh --location office       # office, webcam
#   ./launch.sh --vision rtsp://...     # custom RTSP source
#   ./launch.sh --location farm --omcs HP.71.P01
#
# Requires: OpenClaw installed, NATS running, Python deps installed
# ─────────────────────────────────────────────────────────────────────────────

set -euo pipefail

# ── Defaults ──────────────────────────────────────────────────────────────────
LOCATION="${LOCATION:-home}"
OMCS_COORD="${OMCS_COORD:-IT.00.P00}"
VISION_SOURCE="${VISION_SOURCE:-0}"
NATS_URL="${NATS_URL:-nats://localhost:4222}"
WORKSPACE="$HOME/.openclaw/workspace"
LOG_DIR="$HOME/.openclaw/logs"
SKILLS_DIR="$WORKSPACE/skills"

# ── Parse args ────────────────────────────────────────────────────────────────
while [[ $# -gt 0 ]]; do
  case $1 in
    --location) LOCATION="$2"; shift 2 ;;
    --omcs)     OMCS_COORD="$2"; shift 2 ;;
    --vision)   VISION_SOURCE="$2"; shift 2 ;;
    *) shift ;;
  esac
done

mkdir -p "$LOG_DIR"

echo "══════════════════════════════════════════════"
echo "  LEGACY STACK — $(date '+%Y-%m-%d %H:%M:%S')"
echo "  Location : $LOCATION"
echo "  OMCS     : $OMCS_COORD"
echo "  Vision   : $VISION_SOURCE"
echo "  NATS     : $NATS_URL"
echo "══════════════════════════════════════════════"

# ── 1. Check NATS ─────────────────────────────────────────────────────────────
echo "[launcher] Checking NATS..."
if ! nc -z localhost 4222 2>/dev/null; then
  echo "[launcher] NATS not running. Starting nats-server..."
  nats-server -js &
  sleep 2
fi
echo "[launcher] NATS OK"

# ── 2. Start legacy-pulse (NATS → AMBIENT_STATE) ─────────────────────────────
echo "[launcher] Starting legacy-pulse..."
NATS_URL="$NATS_URL" OPENCLAW_WORKSPACE="$WORKSPACE" \
  python3 "$SKILLS_DIR/legacy-pulse/legacy-pulse.py" \
  > "$LOG_DIR/legacy-pulse.log" 2>&1 &
PULSE_PID=$!
echo "[launcher] legacy-pulse PID: $PULSE_PID"

# ── 3. Start legacy-eyes (Vision → AMBIENT_STATE) ────────────────────────────
echo "[launcher] Starting legacy-eyes..."
LEGACY_VISION_SOURCE="$VISION_SOURCE" \
NATS_URL="$NATS_URL" \
OPENCLAW_WORKSPACE="$WORKSPACE" \
  python3 "$SKILLS_DIR/legacy-eyes/legacy-eyes.py" \
  > "$LOG_DIR/legacy-eyes.log" 2>&1 &
EYES_PID=$!
echo "[launcher] legacy-eyes PID: $EYES_PID"

# Give daemons 2 seconds to initialize
sleep 2

# ── 4. Fire session start event ───────────────────────────────────────────────
echo "[launcher] Firing session start event..."
NATS_URL="$NATS_URL" python3 "$SKILLS_DIR/legacy-pulse/legacy-session-start.py" \
  --location "$LOCATION" \
  --omcs "$OMCS_COORD"

# ── 5. Start OpenClaw Gateway ─────────────────────────────────────────────────
echo "[launcher] Starting OpenClaw Gateway..."
openclaw gateway \
  > "$LOG_DIR/gateway.log" 2>&1 &
GATEWAY_PID=$!
echo "[launcher] Gateway PID: $GATEWAY_PID"

# ── Save PIDs for shutdown ─────────────────────────────────────────────────────
cat > "$LOG_DIR/legacy.pids" <<EOF
PULSE_PID=$PULSE_PID
EYES_PID=$EYES_PID
GATEWAY_PID=$GATEWAY_PID
EOF

echo ""
echo "══════════════════════════════════════════════"
echo "  LEGACY IS LIVE"
echo "  Gateway   : ws://127.0.0.1:18789"
echo "  Eyes      : 1fps @ $VISION_SOURCE"
echo "  Pulse     : NATS subscriptions active"
echo "  Heartbeat : every 5m → Discord"
echo ""
echo "  To stop: ./stop.sh"
echo "  Logs: $LOG_DIR/"
echo "══════════════════════════════════════════════"

# ── Keep alive / tail logs ────────────────────────────────────────────────────
trap 'echo "[launcher] Caught interrupt. Stopping..."; kill $PULSE_PID $EYES_PID $GATEWAY_PID 2>/dev/null; exit' INT TERM

tail -f "$LOG_DIR/gateway.log" &
wait $GATEWAY_PID
