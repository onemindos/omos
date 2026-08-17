#!/usr/bin/env bash
set -euo pipefail
DB="${MURMUR_DB:-/Users/zeus/ATAK/voice/data/mumble-server.sqlite}"
if ! command -v sqlite3 >/dev/null; then
  echo "sqlite3 not found" >&2
  exit 1
fi
sqlite3 -header -column "$DB" "select user_id,name,last_active,last_disconnect from users order by user_id;"
