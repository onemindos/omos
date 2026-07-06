#!/usr/bin/env python3
"""
NATS monitoring (:8222) -> ClickHouse onemind.metrics

Native metric history for the OneMind bus. No Prometheus, no system-account creds.
The NATS monitoring endpoint (varz/jsz) is reliable and auth-free; we scrape the
key gauges + per-stream JetStream state and land them in ClickHouse.

Runs via launchd (omos.observ.metrics.plist) every 15s.
Env overrides: NATS_MON (default http://localhost:8222), CH_PORT (default 9001).

If NATS Surveyor (:7777) ever gets system-account creds and starts surveying
(nats_survey_surveyed_count > 0), set SURVEYOR=http://localhost:7777/metrics to
also fold its Prometheus gauges in.
"""
import json
import os
import subprocess
import sys
import urllib.request

MON = os.environ.get("NATS_MON", "http://localhost:8222")
CH_PORT = os.environ.get("CH_PORT", "9001")
SURVEYOR = os.environ.get("SURVEYOR", "")  # e.g. http://localhost:7777/metrics


def get(path):
    with urllib.request.urlopen(MON + path, timeout=5) as r:
        return json.load(r)


def num(x):
    try:
        return float(x)
    except (TypeError, ValueError):
        return None


rows = []  # (name, labels, value)

# 1) Server vitals — /varz
try:
    v = get("/varz")
    for k in ("connections", "total_connections", "in_msgs", "out_msgs",
              "in_bytes", "out_bytes", "slow_consumers", "subscriptions",
              "mem", "cpu", "routes", "remotes"):
        val = num(v.get(k))
        if val is not None:
            rows.append((f"nats_varz_{k}", "", val))
except Exception as e:
    print(f"varz: {e}", file=sys.stderr)

# 2) JetStream totals + per-stream state — /jsz?streams=1
try:
    j = get("/jsz?streams=1")
    for k in ("streams", "consumers", "messages", "bytes", "memory", "storage",
              "api_total", "api_errors"):
        val = num(j.get(k))
        if val is not None:
            rows.append((f"nats_jsz_{k}", "", val))
    for acc in (j.get("account_details") or []):
        for st in (acc.get("stream_detail") or []):
            name = st.get("name", "?")
            state = st.get("state") or {}
            lbl = f"stream={name}"
            for mk, sk in (("messages", "messages"), ("bytes", "bytes"),
                           ("consumers", "consumer_count"),
                           ("first_seq", "first_seq"), ("last_seq", "last_seq")):
                val = num(state.get(sk))
                if val is not None:
                    rows.append((f"nats_stream_{mk}", lbl, val))
except Exception as e:
    print(f"jsz: {e}", file=sys.stderr)

# 3) Optional: fold in Surveyor Prometheus gauges if it's actually surveying
if SURVEYOR:
    try:
        with urllib.request.urlopen(SURVEYOR, timeout=5) as r:
            for line in r.read().decode().splitlines():
                if not line or line.startswith("#"):
                    continue
                parts = line.rsplit(" ", 1)
                if len(parts) != 2:
                    continue
                key, valstr = parts
                val = num(valstr)
                if val is None:
                    continue
                if "{" in key:
                    name = key[:key.index("{")]
                    labels = key[key.index("{") + 1:key.rindex("}")]
                else:
                    name, labels = key, ""
                rows.append((name, labels, val))
    except Exception as e:
        print(f"surveyor: {e}", file=sys.stderr)

if not rows:
    print("no metrics scraped", file=sys.stderr)
    sys.exit(1)

tsv = "\n".join(f"{n}\t{l}\t{v}" for n, l, v in rows).encode()
subprocess.run(
    ["clickhouse", "client", "--port", CH_PORT,
     "-q", "INSERT INTO onemind.metrics (name,labels,value) FORMAT TSV"],
    input=tsv, check=True,
)
print(f"inserted {len(rows)} metric rows")
