"""Publish committed mission facts through the existing NATS event spine."""
import asyncio
import json
import os
from . import store
from . import bus


def envelope(row):
    return {"v": 1, "id": row["id"], "ts": row["created_at"].isoformat(timespec="milliseconds"),
            "type": row["subject"], "source": "service.onemind-core", "subject": row["subject"],
            "corr": row["payload"]["mission_id"], "data": row["payload"]}


async def publish_outbox():
    while True:
        nc = None
        try:
            nc = await bus.connect()
            js = nc.jetstream()
            while True:
                for row in await asyncio.to_thread(store.pending_events):
                    body = envelope(row)
                    await js.publish(row["subject"], json.dumps(body).encode(),
                                     headers={"Nats-Msg-Id": row["id"]}, timeout=10)
                    await asyncio.to_thread(store.event_sent, row["id"])
                await asyncio.sleep(2)
        except asyncio.CancelledError:
            raise
        except Exception:
            # Committed events remain in the database until JetStream acknowledges.
            await asyncio.sleep(5)
        finally:
            if nc:
                await nc.close()
