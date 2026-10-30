#!/usr/bin/env bash
# Provision TAK Server Data Sync channels for OneMind feed layers.
# Run once after a fresh TAK Server deployment, or to restore channels.
# Idempotent — existing channels return 200, new ones return 201.
#
# Usage: ./provision-channels.sh [context]
# Default context: tailscale-operator.tail717ba0.ts.net

set -euo pipefail

CONTEXT="${1:-tailscale-operator.tail717ba0.ts.net}"
POD=$(kubectl get pod -n takserver -l app=takserver-api --context "$CONTEXT" \
      -o jsonpath='{.items[0].metadata.name}')

echo "Using pod: $POD"

# TODO(step-2, PR tak-ca-password-rotation-step1): replace the stock keystore password with a Secret reference; see docs/operations/tak-ca-password-rotation-step1-2026-10-29.md
kubectl exec -n takserver "$POD" --context "$CONTEXT" -- sh -c '
  openssl pkcs12 -legacy -in /certs-configmap/admin.p12 -passin pass:atakatak \
    -nokeys -clcerts -out /tmp/admin.crt 2>/dev/null
  openssl pkcs12 -legacy -in /certs-configmap/admin.p12 -passin pass:atakatak \
    -nocerts -nodes -out /tmp/admin.key 2>/dev/null

  # Channel name → description
  # Air, Sea, Land = feed layers (aircraft/AIS/APRS+ground)
  # Drones = MAVLink drone telemetry from NATS bridge
  # Hazards = earthquakes + NWS weather alerts
  # Team = team member SA (reserved — feeds do not route here; client SA does)
  # Intel = ACLED / intelligence feeds (future)

  for CH in Air Sea Land Drones Hazards Team Intel; do
    STATUS=$(curl -sk --cert /tmp/admin.crt --key /tmp/admin.key \
      -X PUT "https://127.0.0.1:8443/Marti/api/missions/${CH}?creatorUid=omos-infra&defaultRole=MISSION_SUBSCRIBER" \
      -w "%{http_code}" -o /tmp/resp.json)
    if [ "$STATUS" = "201" ]; then
      echo "CREATED  $CH"
    elif [ "$STATUS" = "200" ]; then
      echo "EXISTS   $CH"
    else
      echo "ERROR    $CH HTTP $STATUS: $(cat /tmp/resp.json)"
    fi
  done
'

echo "Done."
