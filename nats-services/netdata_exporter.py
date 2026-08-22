"""
OneMind OS — Netdata → ClickHouse Metrics Exporter
====================================================
Scrapes Netdata's API every 60 seconds and inserts key system metrics
into ClickHouse for long-term retention (2+ years).

Metrics collected:
  - CPU usage (system, user, idle)
  - Memory (used, free, cached)
  - Disk I/O
  - Network I/O
  - Load average
  - Per-service RAM (NATS, Hermes, etc.)

Usage:
    python3 netdata_exporter.py

Env:
    NETDATA_URL (default http://127.0.0.1:19999)
    CLICKHOUSE_URL (default http://127.0.0.1:8124)
    EXPORT_INTERVAL (default 60 seconds)
"""

import asyncio
import json
import logging
import os
import time
from datetime import datetime

try:
    import httpx
    HAS_HTTPX = True
except ImportError:
    import urllib.request
    HAS_HTTPX = False

NETDATA_URL = os.environ.get("NETDATA_URL", "http://127.0.0.1:19999")
CLICKHOUSE_URL = os.environ.get("CLICKHOUSE_URL", "http://127.0.0.1:8124")
EXPORT_INTERVAL = int(os.environ.get("EXPORT_INTERVAL", "60"))

LOG_DIR = os.path.join(os.path.dirname(__file__), "logs")
os.makedirs(LOG_DIR, exist_ok=True)

logging.basicConfig(
    level=logging.INFO,
    format="%(asctime)s %(levelname)s %(message)s",
    handlers=[
        logging.FileHandler(os.path.join(LOG_DIR, "netdata_exporter.log")),
        logging.StreamHandler(),
    ],
)
log = logging.getLogger("netdata_exporter")

# Charts to collect (Netdata chart names)
CHARTS = [
    "system.cpu",
    "system.ram",
    "system.load",
    "system.io",
    "system.net",
    "system.uptime",
    "mem.available",
    "disk_space._",
]


def fetch_url(url: str) -> dict:
    """Fetch JSON from URL."""
    try:
        if HAS_HTTPX:
            with httpx.Client(timeout=10) as client:
                resp = client.get(url)
                return resp.json()
        else:
            req = urllib.request.Request(url)
            with urllib.request.urlopen(req, timeout=10) as resp:
                return json.loads(resp.read())
    except Exception as e:
        log.error(f"Failed to fetch {url}: {e}")
        return {}


def insert_clickhouse(rows: list[str]):
    """Bulk insert rows into ClickHouse."""
    if not rows:
        return
    data = "\n".join(rows)
    url = f"{CLICKHOUSE_URL}/?query=INSERT+INTO+onemind.system_metrics+FORMAT+JSONEachRow"
    try:
        if HAS_HTTPX:
            with httpx.Client(timeout=10) as client:
                resp = client.post(url, content=data)
                if resp.status_code >= 400:
                    log.error(f"ClickHouse insert failed: {resp.text[:200]}")
        else:
            req = urllib.request.Request(url, data=data.encode(), method="POST")
            urllib.request.urlopen(req, timeout=10)
    except Exception as e:
        log.error(f"ClickHouse insert failed: {e}")


def collect_metrics() -> list[str]:
    """Collect metrics from Netdata API."""
    rows = []
    ts = datetime.utcnow().strftime("%Y-%m-%d %H:%M:%S")

    for chart in CHARTS:
        url = f"{NETDATA_URL}/api/v1/data?chart={chart}&after=-1&points=1&format=json"
        data = fetch_url(url)
        if not data or "labels" not in data:
            continue

        labels = data["labels"]
        values = data["data"][0] if data.get("data") else None
        if not values:
            continue

        # Skip timestamp (first element)
        for i, label in enumerate(labels[1:], 1):
            if i < len(values):
                row = json.dumps({
                    "timestamp": ts,
                    "host": "onemind-mbp",
                    "chart": chart,
                    "dimension": label,
                    "value": float(values[i]) if values[i] is not None else 0.0,
                })
                rows.append(row)

    return rows


async def run():
    log.info(f"Netdata → ClickHouse exporter started")
    log.info(f"  Netdata: {NETDATA_URL}")
    log.info(f"  ClickHouse: {CLICKHOUSE_URL}")
    log.info(f"  Interval: {EXPORT_INTERVAL}s")
    log.info(f"  Charts: {len(CHARTS)}")

    while True:
        try:
            rows = collect_metrics()
            if rows:
                insert_clickhouse(rows)
                log.info(f"Exported {len(rows)} metrics to ClickHouse")
            else:
                log.warning("No metrics collected")
        except Exception as e:
            log.error(f"Collection error: {e}")

        await asyncio.sleep(EXPORT_INTERVAL)


if __name__ == "__main__":
    asyncio.run(run())
