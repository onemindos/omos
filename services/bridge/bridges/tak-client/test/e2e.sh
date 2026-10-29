#!/usr/bin/env bash
# End-to-end check of the tak-client binary against the scripted TAK peer (test/peer.py)
# and a disposable nats-server with JetStream, in XML and in TAK Protocol v1 mode.
# Needs: go, nats-server, nats (CLI), openssl, python3 with the takproto package.
#   PYTHON=/path/to/venv/bin/python ./test/e2e.sh
# Everything (certificates, JetStream store) lives in a temp dir that is removed at exit.
set -euo pipefail
cd "$(dirname "$0")/.."
PY="${PYTHON:-python3}"
WORK="$(mktemp -d "${TMPDIR:-/tmp}/tak-client-e2e.XXXXXX")"
PIDS=()
cleanup() { for p in "${PIDS[@]}"; do kill "$p" 2>/dev/null || true; done; rm -rf "$WORK"; }
trap cleanup EXIT

free_port() { "$PY" -I -c 'import socket; s=socket.socket(); s.bind(("127.0.0.1",0)); print(s.getsockname()[1])'; }

echo "== test PKI (throwaway, in $WORK)"
cd "$WORK"
openssl req -x509 -newkey ec -pkeyopt ec_paramgen_curve:P-256 -nodes -days 1 -subj "/CN=e2e-tak-ca" \
  -keyout ca.key -out ca.pem 2>/dev/null
mk() { # name cn eku
  openssl req -newkey ec -pkeyopt ec_paramgen_curve:P-256 -nodes -subj "/CN=$2" -keyout "$1.key" -out "$1.csr" 2>/dev/null
  printf 'extendedKeyUsage=%s\nsubjectAltName=DNS:takserver\n' "$3" > "$1.ext"
  openssl x509 -req -in "$1.csr" -CA ca.pem -CAkey ca.key -CAcreateserial -days 1 -extfile "$1.ext" -out "$1.pem" 2>/dev/null
}
mk server takserver serverAuth
mk client omos-tak-client clientAuth
cd - >/dev/null

go build -o "$WORK/tak-client" .

NPORT="$(free_port)"
nats-server -a 127.0.0.1 -p "$NPORT" -js -sd "$WORK/js" >"$WORK/nats.log" 2>&1 &
PIDS+=($!)
sleep 1
NATS=(nats -s "nats://127.0.0.1:$NPORT")
"${NATS[@]}" stream add COT --subjects 'tak.cot.>' --storage file --replicas 1 --dupe-window 2m --defaults >/dev/null
echo "== nats-server :$NPORT, stream COT (tak.cot.>)"

run_mode() {
  local mode="$1" tport before after
  tport="$(free_port)"
  echo "== mode $mode: peer on :$tport"
  before="$("${NATS[@]}" stream subjects COT tak.cot.out --json 2>/dev/null | "$PY" -I -c 'import sys,json; d=json.load(sys.stdin) or {}; print(d.get("tak.cot.out",0))')"
  "$PY" -I test/peer.py --port "$tport" --cert "$WORK/server.pem" --key "$WORK/server.key" --ca "$WORK/ca.pem" \
    --mode "$mode" --send 2 --expect 1 --timeout 20 >"$WORK/peer-$mode.log" 2>&1 &
  local peer=$!
  sleep 0.5
  NATS_URL="nats://127.0.0.1:$NPORT" TAK_ADDR="127.0.0.1:$tport" TAK_CERT_FILE="$WORK/client.pem" \
    TAK_KEY_FILE="$WORK/client.key" TAK_CA_FILE="$WORK/ca.pem" HTTP_ADDR="127.0.0.1:0" \
    "$WORK/tak-client" >"$WORK/client-$mode.log" 2>&1 &
  local client=$!
  PIDS+=("$client")
  for _ in $(seq 1 50); do
    after="$("${NATS[@]}" stream subjects COT tak.cot.out --json 2>/dev/null | "$PY" -I -c 'import sys,json; d=json.load(sys.stdin) or {}; print(d.get("tak.cot.out",0))')"
    [ "$after" -ge $((before + 2)) ] && break
    sleep 0.2
  done
  echo "   tak.cot.out messages: $before -> $after"
  "${NATS[@]}" stream get COT --last-for tak.cot.out 2>/dev/null | grep -E 'Nats-Msg-Id|Tak-Proto|Tak-Client|callsign' | sed 's/^/   /' | cut -c1-200
  "${NATS[@]}" pub tak.cot.inject "<event version=\"2.0\" uid=\"OM.drone.sim-x1\" type=\"a-f-A-M-H-Q\" how=\"m-g\" time=\"$(date -u +%Y-%m-%dT%H:%M:%SZ)\" start=\"$(date -u +%Y-%m-%dT%H:%M:%SZ)\" stale=\"2099-01-01T00:00:00Z\"><point lat=\"30.2\" lon=\"-97.7\" hae=\"150\" ce=\"10\" le=\"10\"/><detail><contact callsign=\"SIM-X1\"/></detail></event>" >/dev/null 2>&1
  local rc=0
  wait "$peer" || rc=$?
  kill "$client" 2>/dev/null || true
  echo "   peer transcript:"; sed 's/^/     /' "$WORK/peer-$mode.log"
  echo "   tak-client log:"; sed 's/^[0-9/]* [0-9:]* /     /' "$WORK/client-$mode.log"
  if [ "$rc" -ne 0 ] || [ "$after" -lt $((before + 2)) ]; then echo "FAIL mode $mode"; return 1; fi
  echo "PASS mode $mode"
}

run_mode xml
run_mode v1
echo "== all modes passed"
