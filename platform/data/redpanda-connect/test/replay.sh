#!/usr/bin/env bash
# Replay test for the grammar feed pipelines (FEATURES A1).
#
# Runs the pipelines exactly as they are in the redpanda-connect-pipelines ConfigMap
# (deployment.yaml) against a throwaway nats-server and recorded fixtures served over HTTP,
# captures everything with `nats sub '>'`, and checks that only grammar subjects come out and
# nothing lands on tak.cot.inject. A negative control then runs the legacy feed-mta pipeline
# (URLs and NATS address rewritten to the local ones) and checks that the harness does see its
# tak.cot.inject traffic, so a pass is not an artefact of a deaf capture.
#
# Needs on PATH (or via env): nats-server, nats, redpanda-connect (or REDPANDA_CONNECT), and
# python3 with PyYAML (or PYTHON). Nothing touches a cluster.
#
#   platform/data/redpanda-connect/test/replay.sh [evidence-dir]
set -euo pipefail

HERE="$(cd "$(dirname "$0")" && pwd)"
RPC_DIR="$(cd "$HERE/.." && pwd)"
RPC="${REDPANDA_CONNECT:-redpanda-connect}"
NATS_SERVER="${NATS_SERVER:-nats-server}"
NATS="$(command -v "${NATS_CLI:-nats}")"
PY="${PYTHON:-python3}"
OUT="${1:-}"

WORK="$(mktemp -d "${TMPDIR:-/tmp}/feed-replay.XXXXXX")"
PIDS=()
cleanup() {
  for p in "${PIDS[@]:-}"; do [ -n "$p" ] && kill "$p" 2>/dev/null || true; done
  wait 2>/dev/null || true
  if [ -n "$OUT" ]; then mkdir -p "$OUT"; cp -R "$WORK/logs" "$OUT/" 2>/dev/null || true; fi
  rm -rf "$WORK"
}
trap cleanup EXIT

free_port() { "$PY" -c 'import socket; s=socket.socket(); s.bind(("127.0.0.1",0)); print(s.getsockname()[1]); s.close()'; }
mkdir -p "$WORK/logs" "$WORK/store" "$WORK/www" "$WORK/pipelines" "$WORK/proto"

NPORT="$(free_port)"; HPORT="$(free_port)"
# Throwaway credentials for the throwaway server, generated per run (the outputs refuse empty auth).
NATS_USER="replay"; NATS_PASSWORD="$("$PY" -c 'import secrets; print(secrets.token_hex(16))')"
export NATS_URL="nats://127.0.0.1:$NPORT" NATS_TLS_ENABLED=false NATS_USER NATS_PASSWORD
natsc() { "$NATS" -s "$NATS_URL" --user "$NATS_USER" --password "$NATS_PASSWORD" "$@"; }

echo "== nats-server on $NPORT"
"$NATS_SERVER" -a 127.0.0.1 -p "$NPORT" --user "$NATS_USER" --pass "$NATS_PASSWORD" -js -sd "$WORK/store" >"$WORK/logs/nats-server.log" 2>&1 &
PIDS+=($!)
for _ in $(seq 1 50); do natsc rtt >/dev/null 2>&1 && break; sleep 0.1; done

# The streams the pipelines and the legacy output land in (names and filters as in
# platform/nats/streams/, duplicate window 2m as there).
for s in "COT tak.cot.>" "WORLD evt.world.>" "INTEL intel.>" "SYS evt.sys.>"; do
  read -r name subj <<<"$s"
  natsc stream add "$name" --subjects "$subj" --storage file --replicas 1 \
    --retention limits --discard old --max-age 24h --dupe-window 2m --max-msgs=-1 \
    --max-bytes=-1 --max-msg-size=-1 --max-msgs-per-subject=-1 --no-allow-rollup \
    --no-deny-delete --no-deny-purge --allow-direct --defaults >/dev/null
done

echo "== extract pipelines from the ConfigMap"
"$PY" "$HERE/replay_check.py" extract "$RPC_DIR/deployment.yaml" "$WORK/pipelines"
cp "$RPC_DIR/gtfs-realtime.proto" "$WORK/proto/"

echo "== encode GTFS-RT fixtures and serve them on $HPORT"
mkdir -p "$WORK/www/Dataservice/mtagtfsfeeds/mnr" "$WORK/www/Dataservice/mtagtfsfeeds/lirr" "$WORK/www/NORAD/elements"
for f in mnr lirr; do
  cat >"$WORK/encode-$f.yaml" <<EOF
input:
  file:
    paths: ["$HERE/fixtures/$f.json"]
    scanner: {to_the_end: {}}
pipeline:
  processors:
    - protobuf: {operator: from_json, message: transit_realtime.FeedMessage, import_paths: ["$WORK/proto"]}
output:
  file: {path: "$WORK/www/Dataservice/mtagtfsfeeds/$f/gtfs-$f", codec: all-bytes}
logger: {level: WARN}
EOF
  "$RPC" run "$WORK/encode-$f.yaml"
done
cp "$HERE/fixtures/stations.tle" "$WORK/www/NORAD/elements/gp.php"
"$PY" -m http.server "$HPORT" --bind 127.0.0.1 --directory "$WORK/www" >"$WORK/logs/http.log" 2>&1 &
PIDS+=($!)
sleep 0.5

run_capture() { # name seconds pipeline-files...
  local name="$1" secs="$2"; shift 2
  mkdir -p "$WORK/logs/$name/dump"
  "$NATS" -s "$NATS_URL" --user "$NATS_USER" --password "$NATS_PASSWORD" sub ">" --dump="$WORK/logs/$name/dump" >"$WORK/logs/$name/sub.log" 2>&1 &
  local sub=$!
  sleep 0.5
  MTA_BASE_URL="http://127.0.0.1:$HPORT" CELESTRAK_BASE_URL="http://127.0.0.1:$HPORT" \
    MTA_PROTO_DIR="$WORK/proto" MTA_POLL_INTERVAL=3s CELESTRAK_POLL_INTERVAL=3s \
    "$RPC" streams --chilled "$@" >"$WORK/logs/$name/redpanda-connect.log" 2>&1 &
  local rp=$!
  sleep "$secs"
  kill "$rp" 2>/dev/null || true; wait "$rp" 2>/dev/null || true
  sleep 0.5
  kill "$sub" 2>/dev/null || true; wait "$sub" 2>/dev/null || true
  if grep -q 'level=error msg="failed to create stream' "$WORK/logs/$name/redpanda-connect.log"; then
    echo "FAIL: a pipeline did not start ($name):"; grep 'failed to create stream' "$WORK/logs/$name/redpanda-connect.log"; exit 1
  fi
  natsc stream ls --json >"$WORK/logs/$name/streams.json" 2>/dev/null || true
  for s in COT WORLD INTEL SYS; do
    natsc stream info "$s" --json >"$WORK/logs/$name/stream-$s.json"
  done
}

echo "== grammar pipelines (about 3 polls each)"
run_capture grammar 10 "$WORK/pipelines/feed-mta-grammar.yaml" "$WORK/pipelines/feed-satellites-grammar.yaml"

echo "== negative control: legacy feed-mta rewritten to the local fixture and NATS"
sed -e "s#https://api-endpoint.mta.info#http://127.0.0.1:$HPORT#g" \
    -e "s#nats://nats.nats.svc.cluster.local:4222#$NATS_URL#g" \
    -e "s#enabled: true#enabled: false#" \
    -e "s#/proto/#$WORK/proto/#g" \
    -e 's#interval: "30s"#interval: "3s"#' \
    "$WORK/pipelines/feed-mta.yaml" >"$WORK/pipelines/legacy-feed-mta-local.yaml"
run_capture legacy 5 "$WORK/pipelines/legacy-feed-mta-local.yaml"

echo "== check"
"$PY" "$HERE/replay_check.py" check "$HERE/fixtures" "$WORK/logs" | tee "$WORK/logs/check.txt"
