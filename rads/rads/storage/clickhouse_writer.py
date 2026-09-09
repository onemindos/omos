"""
ClickHouse Time Machine writer — drains RADS_LAND + RADS_ROBOTS streams
into onemind DB for 10-year historical soil/drone data.

Tables auto-created on first run.
"""
import asyncio, json, time
import nats
from clickhouse_driver import Client
from rads.config import (NATS_URL, HOMESTEAD_ID,
                         CLICKHOUSE_HOST, CLICKHOUSE_PORT, CLICKHOUSE_DB)
from rads.envelope import unwrap

CREATE_SOIL = """
CREATE TABLE IF NOT EXISTS rads_soil (
    id           UUID DEFAULT generateUUIDv4(),
    ts           DateTime64(3),
    homestead_id String,
    field        UInt8,
    moisture     Float32,
    temperature_c Float32,
    humidity_pct Float32
) ENGINE = MergeTree()
ORDER BY (homestead_id, field, ts)
TTL toDateTime(ts) + INTERVAL 10 YEAR
"""

CREATE_DRONE = """
CREATE TABLE IF NOT EXISTS rads_drone_telemetry (
    id           UUID DEFAULT generateUUIDv4(),
    ts           DateTime64(3),
    homestead_id String,
    drone_id     String,
    lat          Float64,
    lon          Float64,
    alt_m        Float32,
    alt_rel_m    Float32
) ENGINE = MergeTree()
ORDER BY (homestead_id, drone_id, ts)
TTL toDateTime(ts) + INTERVAL 10 YEAR
"""

class ClickHouseWriter:
    def __init__(self):
        self._ch = Client(host=CLICKHOUSE_HOST, port=CLICKHOUSE_PORT,
                          database=CLICKHOUSE_DB)
        self._nc = None

    def _init_tables(self):
        self._ch.execute(CREATE_SOIL)
        self._ch.execute(CREATE_DRONE)
        print("[ch-writer] tables ready", flush=True)

    async def connect(self):
        self._nc = await nats.connect(NATS_URL)
        self._init_tables()
        print(f"[ch-writer] connected → {NATS_URL} + ClickHouse {CLICKHOUSE_HOST}:{CLICKHOUSE_PORT}", flush=True)

    async def _handle_soil(self, msg):
        env  = unwrap(msg.data)
        d    = env["payload"]
        self._ch.execute(
            "INSERT INTO rads_soil (ts, homestead_id, field, moisture, temperature_c, humidity_pct) VALUES",
            [{"ts": d["ts"], "homestead_id": HOMESTEAD_ID, "field": d["field"],
              "moisture": d["moisture"], "temperature_c": d.get("temperature_c", 0),
              "humidity_pct": d.get("humidity_pct", 0)}]
        )
        await msg.ack()

    async def _handle_drone(self, msg):
        env = unwrap(msg.data)
        d   = env["payload"]
        self._ch.execute(
            "INSERT INTO rads_drone_telemetry (ts, homestead_id, drone_id, lat, lon, alt_m, alt_rel_m) VALUES",
            [{"ts": d["ts"], "homestead_id": HOMESTEAD_ID, "drone_id": d["drone_id"],
              "lat": d["lat"], "lon": d["lon"], "alt_m": d["alt_m"], "alt_rel_m": d["alt_rel_m"]}]
        )
        await msg.ack()

    async def run(self):
        await self.connect()
        js = self._nc.jetstream()
        await js.subscribe(f"rad.{HOMESTEAD_ID}.land.field.*.soil",
                           cb=self._handle_soil, durable="ch-writer-soil", stream="RADS_LAND")
        await js.subscribe(f"rad.{HOMESTEAD_ID}.robots.*.telemetry",
                           cb=self._handle_drone, durable="ch-writer-drone", stream="RADS_ROBOTS")
        print("[ch-writer] draining streams → ClickHouse", flush=True)
        await asyncio.Future()

if __name__ == "__main__":
    asyncio.run(ClickHouseWriter().run())
