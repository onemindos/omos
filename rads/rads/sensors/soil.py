"""
Multi-field soil sensor array — publishes native NATS with JetStream persistence.
Subjects: rad.<homestead_id>.land.field.<n>.soil
Payload (OMOS-wrapped):
  { field, moisture, temperature_c, humidity_pct, ts }
"""
import asyncio, random, time
import nats
from nats.js import JetStreamContext
from rads.config import (NATS_URL, HOMESTEAD_ID, FIELD_COUNT,
                         SENSOR_INTERVAL, soil_subject)
from rads.envelope import wrap

class SoilSensorArray:
    def __init__(self):
        self.fields = {
            i: {"moisture": 65.0, "temperature_c": 22.0, "humidity_pct": 55.0}
            for i in range(1, FIELD_COUNT + 1)
        }
        self._nc = None
        self._js: JetStreamContext = None

    async def connect(self):
        self._nc = await nats.connect(NATS_URL)
        self._js = self._nc.jetstream()
        # ensure stream exists
        try:
            await self._js.find_stream_name_by_subject(soil_subject(1))
        except Exception:
            await self._nc.jsm().add_stream(
                name="RADS_LAND",
                subjects=[f"rad.{HOMESTEAD_ID}.land.>"],
                retention="limits",
                max_age=10 * 365 * 24 * 3600,  # 10 years
            )
        print(f"[soil-array] connected → {NATS_URL}", flush=True)

    async def run(self):
        await self.connect()
        while True:
            for field_id, state in self.fields.items():
                # simulate drift
                state["moisture"]     = max(5.0,  min(99.0, state["moisture"]     - random.uniform(0.5, 2.0)))
                state["temperature_c"]= max(5.0,  min(40.0, state["temperature_c"]+ random.uniform(-0.2, 0.2)))
                state["humidity_pct"] = max(20.0, min(90.0, state["humidity_pct"] + random.uniform(-0.5, 0.5)))

                reading = {
                    "field":          field_id,
                    "moisture":       round(state["moisture"], 1),
                    "temperature_c":  round(state["temperature_c"], 1),
                    "humidity_pct":   round(state["humidity_pct"], 1),
                    "ts":             time.time(),
                }
                subj = soil_subject(field_id)
                await self._js.publish(subj, wrap("soil.reading", reading, source=f"field-{field_id}"))
                print(f"[soil] {subj} → {reading}", flush=True)

            await asyncio.sleep(SENSOR_INTERVAL)

    def irrigate(self, field_id: int, amount: float = 30.0):
        if field_id in self.fields:
            self.fields[field_id]["moisture"] = min(95.0, self.fields[field_id]["moisture"] + amount)
            print(f"[soil] field {field_id} irrigated → moisture={self.fields[field_id]['moisture']:.1f}", flush=True)

if __name__ == "__main__":
    asyncio.run(SoilSensorArray().run())
