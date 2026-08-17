#!/usr/bin/env bash
set -euo pipefail

# Cutover: move Murmur runtime from /Users/zeus/murmur into the ATAK home (voice/),
# then recreate the container bound to the new path. Brief voice outage during recreate.
# Data (sqlite) is MOVED, not copied, so there is one source of truth afterward.

OLD_ROOT="${OLD_MURMUR_ROOT:-/Users/zeus/murmur}"
NEW_ROOT="${NEW_MURMUR_ROOT:-/Users/zeus/ATAK/voice}"
DEPLOY_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

if [[ ! -f "$DEPLOY_DIR/.env" ]]; then
  echo "Missing $DEPLOY_DIR/.env — copy .env.example to .env and set MUMBLE_SUPERUSER_PASSWORD first." >&2
  exit 1
fi

echo "== Stopping and removing current murmur container =="
docker stop murmur 2>/dev/null || true
docker rm murmur 2>/dev/null || true

echo "== Moving Murmur data into $NEW_ROOT =="
mkdir -p "$NEW_ROOT"
if [[ -d "$OLD_ROOT/data" && ! -e "$NEW_ROOT/data" ]]; then
  mv "$OLD_ROOT/data" "$NEW_ROOT/data"
else
  echo "Skipping data move (already at $NEW_ROOT/data or source missing)."
fi
# config already copied into voice/config; keep whichever exists
mkdir -p "$NEW_ROOT/config"
[[ -f "$OLD_ROOT/config/murmur.ini" && ! -f "$NEW_ROOT/config/murmur.ini" ]] && cp "$OLD_ROOT/config/murmur.ini" "$NEW_ROOT/config/murmur.ini" || true

echo "== Recreating murmur from deploy/docker-compose.yml =="
cd "$DEPLOY_DIR"
docker compose --env-file .env up -d murmur

echo
echo "Murmur now bound to $NEW_ROOT. Old path $OLD_ROOT can be archived once verified."
docker ps --format 'table {{.Names}}\t{{.Status}}\t{{.Ports}}' | grep -E 'NAMES|murmur' || true
