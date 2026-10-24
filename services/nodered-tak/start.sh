#!/usr/bin/env bash
set -euo pipefail

# Fail before touching persistent state if the operator secrets are absent.
node -e 'require("/opt/onemind/settings.cjs")'
node /opt/onemind/bootstrap.cjs
# Safe mode is deliberate on every restart; an authenticated operator must Deploy.
exec node /usr/src/node-red/node_modules/node-red/red.js \
  --userDir /data --settings /opt/onemind/settings.cjs --safe flows.json
