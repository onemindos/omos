#!/usr/bin/env python3
"""
NATS monitoring (:8222) -> ClickHouse onemind.metrics

Scrapes NATS varz/jsz and optional Surveyor Prometheus metrics, then inserts the
results into ClickHouse. Designed to run continuously in Kubernetes.

Environment variables:
- NATS_MONITORING_URL (or legacy NATS_MON): NATS monitoring endpoint
- SURVEYOR_URL (or legacy SURVEYOR): optional Prometheus metrics endpoint
- CLICKHOUSE_HOST: ClickHouse HTTP host
- CLICKHOUSE_PORT (or legacy CH_PORT): ClickHouse HTTP port
- CLICKHOUSE_USER: ClickHouse username
- CLICKHOUSE_PASSWORD: ClickHouse password
- CLICKHOUSE_DB: ClickHouse database
- SCRAPE_INTERVAL: seconds between scrapes
"""

import json
import os
import sys
import time
import urllib.request

import clickhouse_connect

NATS_MONITORING_URL = os.environ.get("NATS_MONITORING_URL", os.environ.get("NATS_MON", "http://localhost:8222"))
SURVEYOR_URL = os.environ.get("SURVEYOR_URL", os.environ.get("SURVEYOR", ""))
CLICKHOUSE_HOST = os.environ.get("CLICKHOUSE_HOST", "localhost")
CLICKHOUSE_PORT = int(os.environ.get("CLICKHOUSE_PORT", os.environ.get("CH_PORT", "8123")))
CLICKHOUSE_USER = os.environ.get("CLICKHOUSE_USER", "default")
CLICKHOUSE_PASSWORD = os.environ.get("CLICKHOUSE_PASSWORD", "")
CLICKHOUSE_DB = os.environ.get("CLICKHOUSE_DB", "onemind")
SCRAPE_INTERVAL = float(os.environ.get("SCRAPE_INTERVAL", "15"))


def get(path):
    with urllib.request.urlopen(NATS_MONITORING_URL + path, timeout=5) as response:
        return json.load(response)


def num(value):
    try:
        return float(value)
    except (TypeError, ValueError):
        return None


def scrape_rows():
    rows = []  # (name, labels, value)

    # 1) Server vitals — /varz
    try:
        v = get("/varz")
        for key in (
            "connections",
            "total_connections",
            "in_msgs",
            "out_msgs",
            "in_bytes",
            "out_bytes",
            "slow_consumers",
            "subscriptions",
            "mem",
            "cpu",
            "routes",
            "remotes",
        ):
            val = num(v.get(key))
            if val is not None:
                rows.append((f"nats_varz_{key}", "", val))
    except Exception as exc:
        print(f"varz: {exc}", file=sys.stderr, flush=True)

    # 2) JetStream totals + per-stream state — /jsz?streams=1
    try:
        j = get("/jsz?streams=1")
        for key in ("streams", "consumers", "messages", "bytes", "memory", "storage", "api_total", "api_errors"):
            val = num(j.get(key))
            if val is not None:
                rows.append((f"nats_jsz_{key}", "", val))
        for acc in (j.get("account_details") or []):
            for st in (acc.get("stream_detail") or []):
                name = st.get("name", "?")
                state = st.get("state") or {}
                labels = f"stream={name}"
                for metric_name, state_key in (
                    ("messages", "messages"),
                    ("bytes", "bytes"),
                    ("consumers", "consumer_count"),
                    ("first_seq", "first_seq"),
                    ("last_seq", "last_seq"),
                ):
                    val = num(state.get(state_key))
                    if val is not None:
                        rows.append((f"nats_stream_{metric_name}", labels, val))
    except Exception as exc:
        print(f"jsz: {exc}", file=sys.stderr, flush=True)

    # 3) Optional: fold in Surveyor Prometheus gauges if it's actually surveying
    if SURVEYOR_URL:
        try:
            with urllib.request.urlopen(SURVEYOR_URL, timeout=5) as response:
                for line in response.read().decode().splitlines():
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
                        name = key[: key.index("{")]
                        labels = key[key.index("{") + 1 : key.rindex("}")]
                    else:
                        name, labels = key, ""
                    rows.append((name, labels, val))
        except Exception as exc:
            print(f"surveyor: {exc}", file=sys.stderr, flush=True)

    return rows


def make_clickhouse_client():
    return clickhouse_connect.get_client(
        host=CLICKHOUSE_HOST,
        port=CLICKHOUSE_PORT,
        username=CLICKHOUSE_USER,
        password=CLICKHOUSE_PASSWORD,
        database=CLICKHOUSE_DB,
    )


def ensure_schema(client):
    client.command(f"CREATE DATABASE IF NOT EXISTS `{CLICKHOUSE_DB}`")
    client.command(
        f"""
        CREATE TABLE IF NOT EXISTS `{CLICKHOUSE_DB}`.metrics (
            timestamp DateTime64(3) DEFAULT now64(3),
            name String,
            labels String,
            value Float64
        ) ENGINE = MergeTree
        PARTITION BY toYYYYMM(timestamp)
        ORDER BY (name, timestamp)
        TTL timestamp + toIntervalYear(1)
        SETTINGS index_granularity = 8192
        """
    )


def insert_rows(rows):
    client = make_clickhouse_client()
    ensure_schema(client)
    client.insert(
        f"{CLICKHOUSE_DB}.metrics",
        rows,
        column_names=["name", "labels", "value"],
    )


def main():
    while True:
        rows = scrape_rows()
        if not rows:
            print("no metrics scraped", file=sys.stderr, flush=True)
        else:
            try:
                insert_rows(rows)
                print(f"inserted {len(rows)} metric rows", flush=True)
            except Exception as exc:
                print(f"clickhouse: {exc}", file=sys.stderr, flush=True)
        time.sleep(SCRAPE_INTERVAL)


if __name__ == "__main__":
    main()
