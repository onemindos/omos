#!/usr/bin/env bash
# Chaos test for omos-ingest (feature A0): under load, kill -9 nats-server, take the
# upstream down, kill -9 the leading collector, and serve malformed records; then verify
# zero duplicates, zero losses and dead letters for every malformed record.
#
# Needs: go, nats-server and nats (CLI) on PATH. Uses only 127.0.0.1 and a temp store.
# Usage: test/chaos/chaos.sh [out-dir]
set -euo pipefail

HERE="$(cd "$(dirname "$0")" && pwd)"
MOD="$(cd "$HERE/../.." && pwd)"
SCHEMAS="$(cd "$MOD/../../schemas" && pwd)"
OUT="${1:-$MOD/test/chaos/out}"
WORK="${CHAOS_WORK:-$(mktemp -d)}"
NPORT="${CHAOS_NATS_PORT:-34322}"
NURL="nats://127.0.0.1:$NPORT"
mkdir -p "$OUT" "$WORK/bin" "$WORK/js"
T0=$(date +%s)

log() { printf '%s +%03ds %s\n' "$(date -u +%H:%M:%S)" "$(( $(date +%s) - T0 ))" "$*" | tee -a "$OUT/timeline.txt"; }
at() { while (( $(date +%s) - T0 < $1 )); do sleep 0.2; done; }
pids=()
PID_a=""
PID_b=""
# shellcheck disable=SC2329  # invoked by the EXIT trap
cleanup() { for p in "${pids[@]}"; do kill "$p" 2>/dev/null || true; done; [[ -n "${NATS_PID:-}" ]] && kill "$NATS_PID" 2>/dev/null || true; }
trap cleanup EXIT

: > "$OUT/timeline.txt"
log "build"
(cd "$MOD" && go build -o "$WORK/bin/omos-ingest" ./cmd/omos-ingest && go build -o "$WORK/bin/chaos-stub" ./test/chaos/stub && go build -o "$WORK/bin/chaos-verify" ./test/chaos/verify)

start_nats() {
  nats-server -a 127.0.0.1 -p "$NPORT" -js -sd "$WORK/js" >>"$OUT/nats-server.log" 2>&1 &
  NATS_PID=$!
  for _ in $(seq 50); do nats --server "$NURL" rtt >/dev/null 2>&1 && return 0; sleep 0.1; done
  echo "nats-server did not start" >&2; return 1
}
start_collector() { # $1 = instance
  "$WORK/bin/omos-ingest" -config "$HERE/sources.chaos.yaml" -schemas "$SCHEMAS" -nats "$NURL" -instance "$1" \
    -listen "127.0.0.1:$2" -create-buckets >>"$OUT/collector-$1.log" 2>&1 &
  if [[ $1 == a ]]; then PID_a=$!; else PID_b=$!; fi
  pids+=("$!")
}
leader_of() { # prints the instance that leads chaos-nws
  for i in a:19464 b:19465; do
    if curl -s "http://127.0.0.1:${i#*:}/metrics" | grep -q "ingest_feed_leader{instance=\"${i%%:*}\",source=\"chaos-nws\"} 1"; then echo "${i%%:*}"; return; fi
  done
}

start_nats
log "nats-server up on $NPORT (pid $NATS_PID)"
for s in "INTEL:intel.>" "SYS:evt.sys.>" "TEL:tel.feed.>"; do
  nats --server "$NURL" stream add "${s%%:*}" --subjects "${s#*:}" --storage file --retention limits --discard old \
    --dupe-window 2m --max-bytes 1GB --max-age 24h --replicas 1 --defaults >/dev/null
done
log "streams INTEL, SYS, TEL created (dupe window 2m, max-bytes 1GB)"

"$WORK/bin/chaos-stub" -listen 127.0.0.1:38180 -rate 50 -window 30s -sse-rate 50 -log "$OUT/emitted.jsonl" >>"$OUT/stub.log" 2>&1 &
pids+=("$!")
start_collector a 19464
sleep 2
start_collector b 19465
log "stub (50 alerts/s, 50 sse/s) and collectors a, b started"

at 20; log "CHAOS 1: kill -9 nats-server"; kill -9 "$NATS_PID"; wait "$NATS_PID" 2>/dev/null || true
at 35; start_nats; log "nats-server restarted (pid $NATS_PID)"

at 50; curl -s "http://127.0.0.1:38180/ctl/outage?on=1"; log "CHAOS 2: upstream outage (503) on"
at 60; L=$(leader_of); log "leader during outage: ${L:-none}"
curl -s "http://127.0.0.1:$([[ $L == b ]] && echo 19465 || echo 19464)/metrics" | grep '^ingest_' > "$OUT/metrics-during-upstream-outage.txt" || true
nats --server "$NURL" sub 'tel.feed.>' --count 2 --raw > "$OUT/status-during-upstream-outage.json" 2>/dev/null || true
at 65; curl -s "http://127.0.0.1:38180/ctl/outage?on=0"; log "upstream outage off"

at 80; L=$(leader_of); log "CHAOS 3: kill -9 leading collector ${L:-?}"
if [[ $L == b ]]; then kill -9 "$PID_b"; else L=a; kill -9 "$PID_a"; fi
at 95; start_collector "$L" "$([[ $L == b ]] && echo 19465 || echo 19464)"; log "collector $L restarted; new leader: $(leader_of)"

at 110; curl -s "http://127.0.0.1:38180/ctl/malformed?on=1"; log "CHAOS 4: malformed upstream records on"
at 120; curl -s "http://127.0.0.1:38180/ctl/malformed?on=0"; log "malformed off"

at 150; log "collect metrics and verify"
for i in a:19464 b:19465; do curl -s "http://127.0.0.1:${i#*:}/metrics" | grep '^ingest_' > "$OUT/metrics-final-${i%%:*}.txt" || true; done
nats --server "$NURL" stream ls -a > "$OUT/streams.txt" 2>&1 || true
"$WORK/bin/chaos-verify" -nats "$NURL" -log "$OUT/emitted.jsonl" -tail 15s -out "$OUT/verify.json" > /dev/null && RC=0 || RC=$?
log "verify exit $RC ($( [[ $RC == 0 ]] && echo PASS || echo FAIL ))"
grep -h '"lease\|lost\|deferred\|acquired' "$OUT"/collector-*.log | cut -c1-220 > "$OUT/lease-events.txt" || true
gzip -f "$OUT"/collector-*.log "$OUT/emitted.jsonl"
exit "$RC"
