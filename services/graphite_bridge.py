"""
OneMind OS — Graphite to ClickHouse Bridge
===========================================
Receives Netdata metrics via Graphite plaintext protocol (TCP :2003)
and inserts them into ClickHouse via HTTP interface.

Netdata sends: "metric.path value timestamp\n"
We batch-insert into onemind.metrics every 10 seconds.
"""

import asyncio
import logging
import os
import time
from collections import deque

import httpx

CH_URL = os.environ.get("CLICKHOUSE_URL", "http://127.0.0.1:8124")
CH_DB = "onemind"
CH_TABLE = "metrics"
LISTEN_PORT = 2003
BATCH_INTERVAL = 10  # seconds
LOG_DIR = os.path.join(os.path.dirname(__file__), "logs")

os.makedirs(LOG_DIR, exist_ok=True)
logging.basicConfig(
    level=logging.INFO,
    format="%(asctime)s %(levelname)s %(message)s",
    handlers=[
        logging.FileHandler(os.path.join(LOG_DIR, "graphite-bridge.log")),
        logging.StreamHandler(),
    ],
)
log = logging.getLogger("graphite-bridge")

buffer = deque(maxlen=500_000)  # ~50MB max buffer


class GraphiteProtocol(asyncio.Protocol):
    """Receives Graphite plaintext lines and buffers them."""

    def __init__(self):
        self._buf = b""

    def data_received(self, data: bytes):
        self._buf += data
        while b"\n" in self._buf:
            line, self._buf = self._buf.split(b"\n", 1)
            line = line.strip()
            if not line:
                continue
            try:
                decoded = line.decode()
                parts = decoded.split(" ")
                if len(parts) >= 3:
                    # Graphite tagged format: "metric.path;tag=val value timestamp"
                    # Strip tags from path for cleaner storage
                    raw_path = parts[0]
                    path = raw_path.split(";")[0]  # strip tags
                    value = parts[1]
                    ts = parts[2]
                    buffer.append(f"{path}\t{value}\t{ts}")
            except Exception:
                pass


async def flush_to_clickhouse():
    """Periodically flush buffered metrics to ClickHouse."""
    async with httpx.AsyncClient(timeout=30) as client:
        while True:
            await asyncio.sleep(BATCH_INTERVAL)
            if not buffer:
                continue

            # Drain buffer
            rows = []
            while buffer and len(rows) < 100_000:
                rows.append(buffer.popleft())

            tsv_data = "\n".join(rows)
            try:
                resp = await client.post(
                    f"{CH_URL}/?database={CH_DB}&query=INSERT%20INTO%20{CH_TABLE}%20(Path,Value,Time)%20FORMAT%20TabSeparated",
                    content=tsv_data.encode(),
                )
                if resp.status_code == 200:
                    log.info(f"Flushed {len(rows)} metrics to ClickHouse")
                else:
                    log.error(f"ClickHouse error {resp.status_code}: {resp.text[:200]}")
                    # Re-queue failed rows (put back at front)
                    for row in reversed(rows):
                        buffer.appendleft(row)
            except Exception as e:
                log.error(f"ClickHouse connection error: {e}")
                for row in reversed(rows):
                    buffer.appendleft(row)


async def main():
    loop = asyncio.get_event_loop()

    # Start Graphite TCP server
    server = await loop.create_server(GraphiteProtocol, "127.0.0.1", LISTEN_PORT)
    log.info(f"Graphite receiver listening on 127.0.0.1:{LISTEN_PORT}")
    log.info(f"Flushing to ClickHouse at {CH_URL} every {BATCH_INTERVAL}s")

    # Start flush task
    flush_task = asyncio.create_task(flush_to_clickhouse())

    try:
        await asyncio.gather(server.serve_forever(), flush_task)
    except (KeyboardInterrupt, asyncio.CancelledError):
        pass
    finally:
        server.close()


if __name__ == "__main__":
    asyncio.run(main())
