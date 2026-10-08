"""Fresh operational projections, with external services retaining authority."""
import asyncio
import json
import os
import hashlib
import time
import threading
from http.server import BaseHTTPRequestHandler, HTTPServer
import nats
from . import store
from .capabilities import cluster_observations

state = {"database": False, "bus": False, "cluster": 0.0}


class Health(BaseHTTPRequestHandler):
    def do_GET(self):
        ready = state["database"] and state["bus"] and time.monotonic() - state["cluster"] < 120
        self.send_response(200 if ready else 503)
        self.send_header("Content-Type", "application/json")
        self.end_headers()
        self.wfile.write(json.dumps({"ready": ready}).encode())

    def log_message(self, *args):
        pass


def project_cluster():
    data = cluster_observations()
    for group, rows in data["sources"].items():
        for row in rows:
            identity = (row.get("namespace") or "cluster") + "/" + row["name"]
            store.upsert_object({"id": "source:kubernetes/" + group + "/" + identity,
                                 "kind": "kubernetes-" + group.rstrip("s"),
                                 "name": row["name"], "source": "kubernetes/" + group,
                                 "properties": row})
    if not data["unavailable"]:
        state["cluster"] = time.monotonic()


async def run():
    server = HTTPServer(("0.0.0.0", 8081), Health)
    threading.Thread(target=server.serve_forever, daemon=True).start()
    await asyncio.to_thread(store.pool.open, wait=True)
    while True:
        try:
            await asyncio.to_thread(store.objects, 1)
            state["database"] = True
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
                state["bus"] = True
                while nc.is_connected:
                    await asyncio.sleep(5)
            except Exception:
                await asyncio.sleep(5)
            finally:
                state["bus"] = False
                if nc:
                    await nc.close()

    await asyncio.gather(refresh(), observe_bus())


if __name__ == "__main__":
    asyncio.run(run())
