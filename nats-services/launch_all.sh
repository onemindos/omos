#!/bin/bash
# OneMind OS — Launch all subscribers
# Usage: ./launch_all.sh
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

echo "=== OneMind OS Subscribers ==="
echo ""

# home — presence, plugs, context
bash "${SCRIPT_DIR}/launch.sh" home "events.zeus.>"
echo ""

# security — motion, cameras, alert modes
bash "${SCRIPT_DIR}/launch.sh" security "events.security.>"
echo ""

# health — biometrics, battery, device health
bash "${SCRIPT_DIR}/launch.sh" health "events.health.>"
echo ""

# comms — Telegram triggers, mode switching
bash "${SCRIPT_DIR}/launch.sh" comms "events.comms.>"
echo ""

# agents — agent invocation bus (prompt.agent.> and cmd.mode.*)
# Note: agents subscriber watches both prompt.agent.> AND cmd.mode.*
bash "${SCRIPT_DIR}/launch.sh" agents "prompt.agent.>"
echo ""

echo "=== All subscribers launched ==="
echo ""
echo "Check status:   python3 ${SCRIPT_DIR}/ctl.py status"
echo "Show rules:     python3 ${SCRIPT_DIR}/ctl.py rules show"
echo "Test event:     python3 ${SCRIPT_DIR}/ctl.py test events.zeus.left_home '{}'"
echo "Prompt Legacy:  python3 ${SCRIPT_DIR}/ctl.py agent legacy 'What is the current status of OneMind OS?'"
