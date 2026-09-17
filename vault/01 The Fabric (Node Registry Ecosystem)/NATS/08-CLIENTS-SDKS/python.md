---
title: "Python Client (nats-py)"
section: "08-CLIENTS-SDKS"
type: note
tags: [nats, python, asyncio, nats-py, client]
source: https://github.com/nats-io/nats.py
created: 2026-06-14
---

# Python Client (`nats-py`)

asyncio-based client. Package: **`nats-py`** (import `nats`). OneMind's primary agent/runtime language for bus work.

See also: [[NATS/08-CLIENTS-SDKS/README|Clients Index]] · [[NATS/03-KV-OBJECT-STORE/key-value-store|KV Store]] · [[NATS/04-SERVICES-MICRO/services-micro|Micro]]

---

## Install & Connect

```shell
pip install nats-py
```

```python
import asyncio, nats

async def main():
    nc = await nats.connect(
        servers=["nats://localhost:4222"],
        name="my-service",
        user_credentials="zeus.creds",     # or token=, user=/password=, nkeys_seed=
        max_reconnect_attempts=-1,
        reconnect_time_wait=2,
    )
    # ... use nc ...
    await nc.drain()                        # graceful close

asyncio.run(main())
```

---

## Core Pub/Sub & Request-Reply

```python
# Publish
await nc.publish("orders.new", b'{"id":1}')

# Subscribe (async handler)
async def handler(msg):
    print(msg.subject, msg.data)
sub = await nc.subscribe("orders.*", cb=handler)

# Queue subscribe (load balanced)
await nc.subscribe("orders.new", queue="workers", cb=handler)

# Request-reply
resp = await nc.request("time.now", b"", timeout=1)
print(resp.data)

# Responder
async def time_svc(msg):
    await msg.respond(str(time.time()).encode())
await nc.subscribe("time.now", cb=time_svc)
```

---

## JetStream

```python
js = nc.jetstream()

# Stream
await js.add_stream(name="ORDERS", subjects=["orders.>"])

# Publish with dedup id
ack = await js.publish("orders.new", b"...", headers={"Nats-Msg-Id": "order-1"})

# Pull subscribe (preferred)
psub = await js.pull_subscribe("orders.new", durable="dispatch")
msgs = await psub.fetch(10, timeout=5)
for m in msgs:
    await m.ack()              # or m.nak(), m.term(), m.in_progress()
```

> **Pitfalls (verified):**
> - **Don't pass both `queue=` and `durable=`** to `js.subscribe()` — the queue name *is* the durable name; passing both errors.
> - `StreamConfig.max_age` is **seconds (int)**, not nanoseconds — `max_age=30*24*3600` = 30 days. Multiplying by 1e9 expires instantly.

---

## KV & Object Store

```python
kv = await js.create_key_value(bucket="CONFIG", history=5, ttl=3600)
await kv.put("version", b"2.14.2")
entry = await kv.get("version")
await kv.update("version", b"2.14.3", last=entry.revision)   # atomic CAS
watcher = await kv.watch("version")
async for e in watcher:
    print("changed:", e.key, e.value)

obs = await js.create_object_store(bucket="ASSETS")
await obs.put("model.bin", data)
```

> **KV bucket names cannot contain dots** (`coral-cache`, not `coral.cache`).

---

## Micro Service

```python
from nats.micro import add_service
from nats.micro.service import GroupConfig, EndpointConfig

async def add(req):
    await req.respond(b"result")

srv = await add_service(nc, name="calc", version="1.0.0")
grp = srv.add_group(GroupConfig(name="math"))       # config object, not a string
await grp.add_endpoint(EndpointConfig(name="add", handler=add))
```

> **Pitfall (verified, nats-py 2.14):** `add_group("name")` errors — pass `GroupConfig(name=...)`; endpoints take `EndpointConfig(name=..., handler=...)`. Import from `nats.micro.service`.

---

## KV Watch on a separate connection (OneMind pattern)

OneMind's live `zeus-context` injection runs a **KV watch in a background task on its own connection**, flipping `_context_changed=True` so the next LLM turn re-reads context. Keep watch loops off the main request connection to avoid head-of-line blocking.

## Sources
- nats.py — https://github.com/nats-io/nats.py
- nats.py docs — https://nats-io.github.io/nats.py/
- JetStream (py) — https://github.com/nats-io/nats.py/blob/main/nats/js/
- NATS by Example (Python) — https://natsbyexample.com
