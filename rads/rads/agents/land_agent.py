"""
Land Agent — subscribes to soil readings, makes irrigation decisions,
publishes decisions to the bus, and accepts Legacy override commands.

Subjects consumed:  rad.<id>.land.field.*.soil
Subjects produced:  rad.<id>.land.field.<n>.irrigate
                    rad.<id>.land.decisions
Legacy override in: rad.<id>.legacy.override
"""
import asyncio, json
import nats
from rads.config import (NATS_URL, HOMESTEAD_ID, FIELD_COUNT,
                         MOISTURE_THRESHOLD, soil_subject, irrigate_subject,
                         decision_subject, legacy_override_subject)
from rads.envelope import wrap, unwrap

class LandAgent:
    def __init__(self):
        self._nc   = None
        self._js   = None
        self._cool = set()   # fields in irrigation cooldown

    async def connect(self):
        self._nc = await nats.connect(NATS_URL)
        self._js = self._nc.jetstream()
        print(f"[land-agent] connected → {NATS_URL}", flush=True)

    async def _handle_soil(self, msg):
        env  = unwrap(msg.data)
        data = env["payload"]
        field    = data["field"]
        moisture = data["moisture"]
        print(f"[land-agent] field={field} moisture={moisture} temp={data.get('temperature_c')} hum={data.get('humidity_pct')}", flush=True)
        await msg.ack()

        if moisture < MOISTURE_THRESHOLD and field not in self._cool:
            decision = {
                "action":   "irrigate",
                "field":    field,
                "reason":   f"moisture {moisture:.1f} < threshold {MOISTURE_THRESHOLD}",
                "moisture": moisture,
            }
            await self._js.publish(irrigate_subject(field), wrap("land.irrigate", decision, source="land-agent"))
            await self._js.publish(decision_subject(),      wrap("land.decision",  decision, source="land-agent"))
            print(f"[land-agent] IRRIGATE field={field}", flush=True)
            self._cool.add(field)

        elif moisture >= MOISTURE_THRESHOLD + 10:
            self._cool.discard(field)

    async def _handle_override(self, msg):
        """Legacy can send manual override commands: irrigate / pause / resume."""
        env  = unwrap(msg.data)
        cmd  = env["payload"]
        print(f"[land-agent] LEGACY OVERRIDE → {cmd}", flush=True)
        await msg.ack()
        action = cmd.get("action")
        field  = cmd.get("field")
        if action == "irrigate" and field:
            decision = {"action": "irrigate", "field": field, "reason": "legacy-override", "moisture": None}
            await self._js.publish(irrigate_subject(field), wrap("land.irrigate", decision, source="legacy-override"))
        elif action == "pause":
            self._cool = set(range(1, FIELD_COUNT + 1))
            print("[land-agent] irrigation PAUSED by Legacy", flush=True)
        elif action == "resume":
            self._cool.clear()
            print("[land-agent] irrigation RESUMED by Legacy", flush=True)

    async def run(self):
        await self.connect()
        soil_wildcard = f"rad.{HOMESTEAD_ID}.land.field.*.soil"
        override_subj = legacy_override_subject()

        # durable consumers
        await self._js.subscribe(soil_wildcard, cb=self._handle_soil,
                                  durable="land-agent-soil", stream="RADS_LAND")
        await self._nc.subscribe(override_subj, cb=self._handle_override)
        print(f"[land-agent] watching {soil_wildcard}", flush=True)
        print(f"[land-agent] override channel: {override_subj}", flush=True)
        await asyncio.Future()  # run forever

if __name__ == "__main__":
    asyncio.run(LandAgent().run())
