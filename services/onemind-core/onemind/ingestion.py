"""Fresh operational projections, with external services retaining authority."""
import asyncio
import json
import os
import hashlib
import nats
from . import store
from .capabilities import cluster_observations


def project_cluster():
    data = cluster_observations()
    for group, rows in data["sources"].items():
        for row in rows:
            identity = (row.get("namespace") or "cluster") + "/" + row["name"]
            store.upsert_object({"id": "source:kubernetes/" + group + "/" + identity,
                                 "kind": "kubernetes-" + group.rstrip("s"),
                                 "name": row["name"], "source": "kubernetes/" + group,
                                 "properties": row})


async def run():
    await asyncio.to_thread(store.pool.open, wait=True)
    while True:
        try:
            await asyncio.to_thread(store.objects, 1)
            break
        except Exception:
            await asyncio.sleep(5)

    async def refresh():
        while True:
            try:
                await asyncio.to_thread(project_cluster)
            except Exception:
                pass
            await asyncio.sleep(30)

    async def observe_bus():
        while True:
            nc = None
            try:
                nc = await nats.connect(servers=os.environ["NATS_URL"],
                                        user=os.environ["NATS_USER"], password=os.environ["NATS_PASSWORD"],
                                        connect_timeout=5, max_reconnect_attempts=3)

                async def receive(msg):
                    if len(msg.data) > 32000:
                        return
                    try:
                        value = json.loads(msg.data)
                        if not isinstance(value, dict):
                            return
                        identity = hashlib.sha256(msg.subject.encode()).hexdigest()[:32]
                        name = str(value.get("name") or value.get("callsign") or msg.subject)[:200]
                        await asyncio.to_thread(store.upsert_object, {
                            "id": "source:nats/" + identity,
                            "kind": "fleet-observation" if msg.subject.startswith("tel.fleet.") else "world-observation",
                            "name": name, "source": msg.subject,
                            "properties": {"subject": msg.subject, "observation": value},
                        })
                    except Exception:
                        pass

                for subject in ["tel.fleet.>", "evt.world.>", "evt.tak.cot.>"]:
                    await nc.subscribe(subject, cb=receive, pending_msgs_limit=1000,
                                       pending_bytes_limit=16 * 1024 * 1024)
                while nc.is_connected:
                    await asyncio.sleep(5)
            except Exception:
                await asyncio.sleep(5)
            finally:
                if nc:
                    await nc.close()

    await asyncio.gather(refresh(), observe_bus())


if __name__ == "__main__":
    asyncio.run(run())
