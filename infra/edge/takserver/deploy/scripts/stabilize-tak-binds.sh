#!/usr/bin/env bash
set -euo pipefail

ATAK_STACK="${ATAK_STACK:-/Users/zeus/ATAK/server/standard/takserver-docker-5.7-RELEASE-43}"
EXPECTED_ROOT="$ATAK_STACK/tak"

if [[ ! -f "$ATAK_STACK/docker-compose.yml" ]]; then
  echo "Missing compose: $ATAK_STACK/docker-compose.yml" >&2
  exit 1
fi

actual_sources="$(docker inspect takserver --format '{{range .Mounts}}{{println .Source}}{{end}}' 2>/dev/null || true)"
if echo "$actual_sources" | grep -q "$EXPECTED_ROOT"; then
  echo "TAK binds already point at $EXPECTED_ROOT"
  exit 0
fi

echo "Current TAK bind sources do not point at $EXPECTED_ROOT"
echo "$actual_sources"
echo
echo "Recreating takserver from $ATAK_STACK/docker-compose.yml"
echo "Database volume is preserved. This will briefly restart TAK device/API service."
cd "$ATAK_STACK"
docker compose up -d --force-recreate takserver

echo
echo "New TAK bind sources:"
docker inspect takserver --format '{{range .Mounts}}{{println .Source "->" .Destination}}{{end}}'
