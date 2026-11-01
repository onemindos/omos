#!/usr/bin/env bash
# Run the TAK mesh gateway on the Mac next to the local NATS leaf
# (infra/edge/mac-leaf). Reads ./.env if present; the environment wins.
# Mesh keys and NATS credentials come from the environment (the host's secret
# store), never from a file in this repository.
#
#   MESH_INTERFACE=en0 ./run-mac.sh
set -euo pipefail

cd "$(dirname "$0")"

if [[ -f .env ]]; then
  # Keep values already exported in the environment.
  while IFS='=' read -r key value; do
    [[ -z "$key" || "$key" == \#* ]] && continue
    if [[ -z "${!key:-}" ]]; then
      export "$key=$value"
    fi
  done < .env
fi

export NATS_URL="${NATS_URL:-nats://localhost:4222}"
export MESH_TO_TAK_SERVER="${MESH_TO_TAK_SERVER:-false}"

# The leaf must be up: the gateway retries forever, but a typo should be loud.
hostport="${NATS_URL#*://}"
hostport="${hostport##*@}"
host="${hostport%%:*}"
port="${hostport##*:}"
if ! nc -z -w 2 "$host" "$port" 2>/dev/null; then
  echo "NATS leaf not answering at $NATS_URL; start it first (infra/edge/mac-leaf/README.md)" >&2
  exit 1
fi

if [[ -z "${MESH_INTERFACE:-}" ]]; then
  default_if="$(route -n get 239.2.3.1 2>/dev/null | awk '/interface:/ {print $2}')"
  echo "MESH_INTERFACE is unset: the system would use ${default_if:-unknown} for 239.2.3.1." >&2
  echo "Set MESH_INTERFACE to the interface on the ATAK LAN (ifconfig lists them)." >&2
  exit 1
fi
if ! ifconfig "$MESH_INTERFACE" >/dev/null 2>&1; then
  echo "MESH_INTERFACE=$MESH_INTERFACE does not exist" >&2
  exit 1
fi

go build -o tak .
echo "tak-mesh: NATS $NATS_URL, interface $MESH_INTERFACE, MESH_TO_TAK_SERVER=$MESH_TO_TAK_SERVER"
exec ./tak
