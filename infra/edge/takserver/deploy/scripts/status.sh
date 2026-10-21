#!/usr/bin/env bash
set -euo pipefail

TAK_ADMIN_P12="${TAK_ADMIN_P12:-/Users/zeus/ATAK/server/standard/takserver-docker-5.7-RELEASE-43/tak/certs/files/admin.p12}"

echo "== Docker tactical containers =="
docker ps --format 'table {{.Names}}\t{{.Image}}\t{{.Status}}\t{{.Ports}}' | grep -E 'NAMES|takserver|tak-database|murmur|geoserver|omos-grafana|omos-loki|omos-jaeger' || true

echo
echo "== TAK bind source check =="
docker inspect takserver --format '{{range .Mounts}}{{println .Source "->" .Destination}}{{end}}' 2>/dev/null || echo "takserver not found"

echo
echo "== NATS =="
python3 - <<'PY'
import json, urllib.request
try:
    d=json.load(urllib.request.urlopen('http://127.0.0.1:8222/varz', timeout=2))
    print(f"{d.get('server_name')} v{d.get('version')} connections={d.get('connections')} subs={d.get('subscriptions')} jetstream={bool(d.get('jetstream'))}")
except Exception as e:
    print(f"NATS monitor unavailable: {e}")
PY

echo
echo "== TAK active clients =="
if [[ -n "${TAK_ADMIN_P12_PASS:-}" && -f "$TAK_ADMIN_P12" ]]; then
  /usr/bin/curl -sk --cert-type P12 --cert "$TAK_ADMIN_P12:$TAK_ADMIN_P12_PASS" \
    https://localhost:8446/Marti/api/clientEndPoints | python3 -m json.tool || true
else
  echo "Set TAK_ADMIN_P12_PASS to query TAK API. Port-only check:"
  for port in 8089 8444 8446 8448; do
    nc -z 127.0.0.1 "$port" >/dev/null 2>&1 && echo "TAK port $port open" || echo "TAK port $port CLOSED"
  done
fi

echo
echo "== Murmur registered users =="
MURMUR_DB=/Users/zeus/ATAK/voice/data/mumble-server.sqlite
[[ -f "$MURMUR_DB" ]] || MURMUR_DB=/Users/zeus/murmur/data/mumble-server.sqlite   # pre-cutover fallback
if command -v sqlite3 >/dev/null && [[ -f "$MURMUR_DB" ]]; then
  echo "db: $MURMUR_DB"
  sqlite3 -header -column "$MURMUR_DB" "select user_id,name,last_active,last_disconnect from users order by user_id;" || true
else
  echo "sqlite3 missing or Murmur DB not found"
fi
