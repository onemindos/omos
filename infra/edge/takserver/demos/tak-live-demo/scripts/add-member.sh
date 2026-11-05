#!/usr/bin/env bash
# Provision a new TAK team member end-to-end: cert -> auth-file (Red group) ->
# data package ZIP -> server restart so the group actually loads.
#
# THE GROUP FIX: TAKServer only reads UserAuthenticationFile.xml at startup.
# If you add a user while the server is running and they connect BEFORE a restart,
# the server accepts their valid cert but drops them in the __ANON__ group instead
# of Red — so they can't see (and aren't seen by) the Red team. This script
# restarts the server AFTER registering the user, so they land in Red on first
# connect. (This is the bug we hit with carku and fixed for paulino.)
#
# Usage: add-member.sh <username> [public_host]
#   e.g. add-member.sh paulino 74.108.23.135
set -euo pipefail

USER="${1:?usage: add-member.sh <username> [public_host]}"
HOST="${2:-74.108.23.135}"
PASS="${P12_PASS:?set P12_PASS (password for the new member .p12)}"

TAK=/Users/zeus/ATAK/server/standard/takserver-docker-5.7-RELEASE-43/tak
CERTS="$TAK/certs/files"
AUTH="$TAK/UserAuthenticationFile.xml"
BUILD=/Users/zeus/ATAK/devices/scripts/build-itak-datapackage.sh
OUTDIR=/Users/zeus/ATAK/demos/tak-live-demo/packages
CONTAINER=takserver

cd "$CERTS"

echo "[1/5] Generating cert for $USER (signed by TAK-CA)..."
openssl genrsa -out "${USER}.key" 2048 2>/dev/null
openssl req -new -key "${USER}.key" -out "${USER}.csr" \
  -subj "/C=US/ST=TX/L=Austin/O=TAK/OU=TAK-OPS/CN=${USER}" 2>/dev/null
openssl x509 -req -in "${USER}.csr" \
  -CA ca.pem -CAkey ca-plain.key -CAcreateserial \
  -out "${USER}.pem" -days 730 -sha256 2>/dev/null
openssl pkcs12 -legacy -export \
  -in "${USER}.pem" -inkey "${USER}.key" -certfile ca.pem \
  -out "${USER}.p12" -name "${USER}" -passout pass:${PASS} 2>/dev/null
rm -f "${USER}.csr"

FP=$(openssl x509 -in "${USER}.pem" -noout -fingerprint -sha256 | sed 's/sha256 Fingerprint=//')
echo "      fingerprint: $FP"

echo "[2/5] Registering $USER in auth file with Red group..."
if grep -q "identifier=\"${USER}\"" "$AUTH"; then
  echo "      already present — skipping insert"
else
  # insert before the closing tag
  python3 - "$AUTH" "$USER" "$FP" <<'PY'
import sys
auth, user, fp = sys.argv[1], sys.argv[2], sys.argv[3]
with open(auth) as f: x = f.read()
entry = (f'    <User identifier="{user}" fingerprint="{fp}" passwordHashed="false">\n'
         f'        <groupList>Red</groupList>\n    </User>\n')
x = x.replace('</UserAuthenticationFile>', entry + '</UserAuthenticationFile>')
with open(auth, 'w') as f: f.write(x)
PY
fi

echo "[3/5] Building data package ZIP..."
mkdir -p "$OUTDIR"
PKG_DESC="oneMind TAK" P12_PASS="$PASS" bash "$BUILD" \
  "$USER" "$CERTS/${USER}.p12" "$CERTS/truststore-root.p12" \
  "$HOST" "$OUTDIR/OneMind-${USER^}.zip"

echo "[4/5] Restarting TAKServer so the Red group loads (the group fix)..."
docker restart "$CONTAINER" >/dev/null

echo "[5/5] Waiting for server to come back up..."
for i in $(seq 1 40); do
  if nc -z localhost 8089 2>/dev/null; then echo "      up (~$((i*5))s)"; break; fi
  sleep 5
done

echo ""
echo "DONE. Send $OUTDIR/OneMind-${USER^}.zip to $USER."
echo "In iTAK: import ZIP -> Settings -> Team = Red, callsign = ${USER} -> connect."
echo "They will land directly in Red (not __ANON__)."
