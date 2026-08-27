#!/bin/bash
# OneMind OS — Subscriber launchd launcher
# Usage: ./launch.sh <name> [subject_prefix]
#
# Writes ~/Library/LaunchAgents/com.onemind.subscriber.<name>.plist
# and loads it so the subscriber starts now and on every login.
#
# Subscriber names: home | security | health | comms | agents

NAME="${1:-home}"
PREFIX="${2:-events.${NAME}.>}"
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PLIST_PATH="$HOME/Library/LaunchAgents/com.onemind.subscriber.${NAME}.plist"
LOG_PATH="${SCRIPT_DIR}/logs/${NAME}.log"
PYTHON="/opt/homebrew/bin/python3"

mkdir -p "${SCRIPT_DIR}/logs"

cat > "${PLIST_PATH}" <<EOF
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>Label</key>
    <string>com.onemind.subscriber.${NAME}</string>
    <key>ProgramArguments</key>
    <array>
        <string>${PYTHON}</string>
        <string>${SCRIPT_DIR}/subscriber.py</string>
        <string>${NAME}</string>
        <string>${PREFIX}</string>
    </array>
    <key>RunAtLoad</key>
    <true/>
    <key>KeepAlive</key>
    <true/>
    <key>StandardOutPath</key>
    <string>${LOG_PATH}</string>
    <key>StandardErrorPath</key>
    <string>${LOG_PATH}</string>
    <key>WorkingDirectory</key>
    <string>${SCRIPT_DIR}</string>
    <key>EnvironmentVariables</key>
    <dict>
        <key>NATS_URL</key>
        <string>nats://127.0.0.1:4222</string>
        <key>PATH</key>
        <string>/opt/homebrew/bin:/usr/local/bin:/usr/bin:/bin</string>
    </dict>
</dict>
</plist>
EOF

# Unload existing (ignore error if not loaded)
launchctl unload "${PLIST_PATH}" 2>/dev/null

# Load new
launchctl load "${PLIST_PATH}"

echo "Loaded:  com.onemind.subscriber.${NAME}"
echo "Subject: ${PREFIX}"
echo "Log:     ${LOG_PATH}"
echo "Plist:   ${PLIST_PATH}"
echo ""
echo "Commands:"
echo "  Stop:   launchctl unload ${PLIST_PATH}"
echo "  Reload: launchctl unload ${PLIST_PATH} && launchctl load ${PLIST_PATH}"
echo "  Log:    tail -f ${LOG_PATH}"
