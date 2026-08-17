#!/usr/bin/env bash
# Show who is connected and — critically — which authorization GROUP they landed
# in. If someone shows __ANON__ instead of Red, they were added while the server
# was running and connected before a restart. Fix: restart takserver, reconnect.
set -euo pipefail

echo "=== Connected clients (last 2 min) — callsign / group ==="
docker exec tak-database bash -c "PGPASSWORD=martipassword psql -U martiuser -d cot -h localhost -c \"SELECT DISTINCT ON (uid) substring(detail from 'callsign=\\\"([^\\\"]*)\\\"') AS callsign, substring(detail from 'name=\\\"([^\\\"]*)\\\" role') AS grp, time FROM cot_router WHERE time > NOW() - INTERVAL '2 minutes' AND detail LIKE '%callsign%' ORDER BY uid, time DESC;\"" 2>&1

echo ""
echo "=== Last subscription assignments ==="
docker exec takserver bash -c "grep 'Set client' /opt/tak/logs/takserver-messaging.log 2>/dev/null | tail -5" 2>&1
