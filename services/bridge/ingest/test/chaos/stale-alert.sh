#!/usr/bin/env bash
# Proves the stale-feed alert fires end to end: a real Prometheus scrapes a live omos-ingest
# with the unchanged platform/ingest/prometheusrule.yaml rules, the upstream goes down for
# 3 minutes (longer than the alert's for: 2m), and the alert goes pending, firing, then
# resolves after the upstream returns.
# Needs: go, nats-server, nats, prometheus, python3 with PyYAML. Usage: stale-alert.sh [out-dir]
set -euo pipefail
HERE="$(cd "$(dirname "$0")" && pwd)"
MOD="$(cd "$HERE/../.." && pwd)"
REPO="$(cd "$MOD/../../.." && pwd)"
OUT="${1:-$MOD/test/chaos/out-stale}"
WORK="${CHAOS_WORK:-$(mktemp -d)}"
NPORT=34323; NURL="nats://127.0.0.1:$NPORT"; PPORT=39090
mkdir -p "$OUT" "$WORK/bin" "$WORK/js" "$WORK/prom"
T0=$(date +%s)
log() { printf '%s +%03ds %s\n' "$(date -u +%H:%M:%S)" "$(( $(date +%s) - T0 ))" "$*" | tee -a "$OUT/timeline.txt"; }
at() { while (( $(date +%s) - T0 < $1 )); do sleep 0.5; done; }
pids=()
# shellcheck disable=SC2329  # invoked by the EXIT trap
cleanup() { for p in "${pids[@]}"; do kill "$p" 2>/dev/null || true; done; }
trap cleanup EXIT
alerts() { curl -s "http://127.0.0.1:$PPORT/api/v1/alerts" > "$OUT/alerts-$1.json"; log "alerts at $1: $(python3 -c 'import json,sys; d=json.load(open(sys.argv[1])); print(sorted((a["labels"]["alertname"], a["labels"].get("source",""), a["state"]) for a in d["data"]["alerts"]))' "$OUT/alerts-$1.json")"; }

: > "$OUT/timeline.txt"
(cd "$MOD" && go build -o "$WORK/bin/omos-ingest" ./cmd/omos-ingest && go build -o "$WORK/bin/chaos-stub" ./test/chaos/stub)
python3 - "$REPO/platform/ingest/prometheusrule.yaml" "$WORK/rules.yaml" <<'PY'
import sys, yaml
yaml.safe_dump(yaml.safe_load(open(sys.argv[1]))["spec"], open(sys.argv[2], "w"), sort_keys=False)
PY
cat > "$WORK/prometheus.yml" <<YML
global: {scrape_interval: 5s, evaluation_interval: 5s}
rule_files: ["$WORK/rules.yaml"]
scrape_configs:
  - job_name: omos-ingest
    static_configs: [{targets: ["127.0.0.1:19466"]}]
YML
nats-server -a 127.0.0.1 -p "$NPORT" -js -sd "$WORK/js" >>"$OUT/nats-server.log" 2>&1 & pids+=("$!")
sleep 1
for s in "INTEL:intel.>" "SYS:evt.sys.>" "TEL:tel.feed.>"; do
  nats --server "$NURL" stream add "${s%%:*}" --subjects "${s#*:}" --storage file --dupe-window 2m --max-bytes 1GB --replicas 1 --defaults >/dev/null
done
"$WORK/bin/chaos-stub" -listen 127.0.0.1:38180 -rate 5 -window 60s -sse-rate 5 -log "$WORK/emitted.jsonl" >>"$OUT/stub.log" 2>&1 & pids+=("$!")
"$WORK/bin/omos-ingest" -config "$HERE/sources.chaos.yaml" -schemas "$REPO/services/schemas" -nats "$NURL" -instance a \
  -listen 127.0.0.1:19466 -create-buckets >>"$OUT/collector.log" 2>&1 & pids+=("$!")
prometheus --config.file="$WORK/prometheus.yml" --storage.tsdb.path="$WORK/prom" --web.listen-address="127.0.0.1:$PPORT" >>"$OUT/prometheus.log" 2>&1 & pids+=("$!")
log "nats, stub, collector a and prometheus started"
at 45; alerts "45s-healthy"
curl -s "http://127.0.0.1:38180/ctl/outage?on=1"; log "upstream outage on (503) for 180 s"
at 80; alerts "80s-outage-35s"
at 200; alerts "200s-outage-155s"
curl -s "http://127.0.0.1:19466/metrics" | grep '^ingest_feed' > "$OUT/metrics-at-200s.txt"
at 225; curl -s "http://127.0.0.1:38180/ctl/outage?on=0"; log "upstream outage off"
at 275; alerts "275s-recovered"
python3 - "$OUT" <<'PY'
import json, sys
out = sys.argv[1]
st = lambda f: {(a["labels"]["alertname"], a["labels"].get("source")): a["state"] for a in json.load(open(f"{out}/alerts-{f}.json"))["data"]["alerts"]}
h, f, r = st("45s-healthy"), st("200s-outage-155s"), st("275s-recovered")
ok = (not any(k[0] == "OneMindFeedStale" for k in h)
      and f.get(("OneMindFeedStale", "chaos-nws")) == "firing" and f.get(("OneMindFeedStale", "chaos-sse")) == "firing"
      and not any(k[0] == "OneMindFeedStale" for k in r))
print("STALE ALERT", "PASS" if ok else "FAIL", json.dumps({"healthy": list(map(list, h)), "outage": [list(k)+[v] for k, v in f.items()], "recovered": list(map(list, r))}))
sys.exit(0 if ok else 1)
PY
