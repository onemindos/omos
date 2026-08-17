#!/usr/bin/env bash
set -euo pipefail

if [[ $# -lt 2 ]]; then
  echo "Usage: $0 <cert-common-name> <username> [group] [role-note]" >&2
  echo "Example: $0 zeus-phone zeus-phone Red 'Team Member'" >&2
  exit 1
fi

CERT_NAME="$1"
USERNAME="$2"
GROUP="${3:-Red}"
ROLE_NOTE="${4:-Team Member}"

for value in "$CERT_NAME" "$USERNAME" "$GROUP"; do
  if [[ ! "$value" =~ ^[A-Za-z0-9._-]+$ ]]; then
    echo "Invalid TAK identifier: $value" >&2
    echo "Use letters, numbers, dot, underscore, or dash only." >&2
    exit 1
  fi
done

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
STACK_DIR="$(cd "$SCRIPT_DIR/.." && pwd)"
ATAK_STACK="${ATAK_STACK:-/Users/zeus/ATAK/server/standard/takserver-docker-5.7-RELEASE-43}"
ATAK_ROOT="$ATAK_STACK/tak"
CERT_DIR="$ATAK_ROOT/certs/files"

actual_sources="$(docker inspect takserver --format '{{range .Mounts}}{{println .Source}}{{end}}' 2>/dev/null || true)"
if ! echo "$actual_sources" | grep -q "$ATAK_ROOT"; then
  echo "REFUSING: live takserver is not bound to $ATAK_ROOT." >&2
  echo "Run deploy/scripts/stabilize-tak-binds.sh first." >&2
  exit 2
fi

if [[ ! -x "$ATAK_ROOT/certs/makeCert.sh" ]]; then
  echo "Missing cert generator: $ATAK_ROOT/certs/makeCert.sh" >&2
  exit 1
fi

if [[ -e "$CERT_DIR/$CERT_NAME.p12" ]]; then
  echo "Cert already exists: $CERT_DIR/$CERT_NAME.p12" >&2
  exit 1
fi

# certs/files is mounted read-only into takserver, and makeCert.sh is Linux-only.
# Generate in a throwaway container from the same image with the tak dir mounted rw.
# Cert metadata matches the existing CA (C=US, ST=TX, L=Austin, O=TAK, OU=TAK-OPS).
# Client .p12 password = CAPASS (default atakatak) unless overridden.
docker run --rm --entrypoint bash \
  -e STATE="${STATE:-TX}" -e CITY="${CITY:-Austin}" \
  -e ORGANIZATION="${ORGANIZATION:-TAK}" -e ORGANIZATIONAL_UNIT="${ORGANIZATIONAL_UNIT:-TAK-OPS}" \
  -e CAPASS="${CAPASS:-atakatak}" \
  -v "$ATAK_ROOT:/opt/tak" \
  "${TAK_IMAGE:-takserver:latest}" -lc "cd /opt/tak/certs && ./makeCert.sh client '$CERT_NAME'"

# Register the cert-based user and set group (in+out). Cert users use certmod, not usermod.
docker exec takserver bash -lc "cd /opt/tak && java -jar utils/UserManager.jar certmod -g '$GROUP' '/opt/tak/certs/files/$CERT_NAME.pem'"

echo "== user status =="
docker exec takserver bash -lc "cd /opt/tak && java -jar utils/UserManager.jar certmod -s '/opt/tak/certs/files/$CERT_NAME.pem'" || true

mkdir -p "$STACK_DIR/exports"
cp "$CERT_DIR/$CERT_NAME.p12" "$STACK_DIR/exports/$CERT_NAME.p12"
cp "$CERT_DIR/ca.pem" "$STACK_DIR/exports/ca.pem" 2>/dev/null || true
chmod 600 "$STACK_DIR/exports/$CERT_NAME.p12"

echo
echo "Issued TAK device cert: $STACK_DIR/exports/$CERT_NAME.p12 (pass: ${CAPASS:-atakatak})"
echo "TAK username: $CERT_NAME    Group: $GROUP    Role note: $ROLE_NOTE"
echo "Connect ATAK/iTAK to 100.74.182.117:8089 (Tailscale) or 192.168.1.183:8089 (LAN), SSL/TLS."
echo "Register the device in devices/registry.md and its devices/<id>.yaml."
