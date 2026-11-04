"""NATS tools registered with the Hermes plugin context."""
from __future__ import annotations

import asyncio
import json
import logging
from concurrent.futures import ThreadPoolExecutor
from typing import Any

logger = logging.getLogger("nats-context.tools")

# Shared thread pool for offloading asyncio.run() calls from the sync tool handlers
_pool = ThreadPoolExecutor(max_workers=4, thread_name_prefix="nats-tool")

from .conn import connect_options, nats_url  # OneMind: URL, credentials and CA from the environment


def _run(coro) -> Any:
    """Run a coroutine to completion from a synchronous context."""
    return asyncio.get_event_loop().run_until_complete(coro) if asyncio.get_event_loop().is_running() \
        else asyncio.run(coro)


async def _publish(subject: str, payload: str) -> dict:
    import nats
    nc = await nats.connect(nats_url(), **connect_options())
    try:
        await nc.publish(subject, payload.encode())
        await nc.flush()
        return {"ok": True, "subject": subject}
    finally:
        await nc.drain()


async def _request(subject: str, payload: str, timeout: float = 5.0) -> dict:
    import nats
    nc = await nats.connect(nats_url(), **connect_options())
    try:
        resp = await nc.request(subject, payload.encode(), timeout=timeout)
        try:
            data = json.loads(resp.data)
        except Exception:
            data = resp.data.decode("utf-8", errors="replace")
        return {"ok": True, "subject": resp.subject, "data": data}
    finally:
        await nc.drain()


async def _peek(subject: str, n: int = 10) -> dict:
    """Subscribe for a short window and collect up to n messages."""
    import nats
    msgs: list[dict] = []
    nc = await nats.connect(nats_url(), **connect_options())

    async def _cb(msg):
        try:
            data = json.loads(msg.data)
        except Exception:
            data = msg.data.decode("utf-8", errors="replace")
        msgs.append({"subject": msg.subject, "data": data})

    try:
        sub = await nc.subscribe(subject, cb=_cb)
        await asyncio.sleep(2.0)
        await sub.unsubscribe()
        return {"ok": True, "count": len(msgs), "messages": msgs[:n]}
    finally:
        await nc.drain()


async def _jetstream_fetch(stream: str, subject_filter: str | None = None, n: int = 10) -> dict:
    import nats
    nc = await nats.connect(nats_url(), **connect_options())
    try:
        js = nc.jetstream()
        consumer_config = nats.js.api.ConsumerConfig(
            deliver_policy=nats.js.api.DeliverPolicy.LAST_PER_SUBJECT
            if not subject_filter else nats.js.api.DeliverPolicy.NEW,
            filter_subject=subject_filter,
            max_deliver=1,
            ack_policy=nats.js.api.AckPolicy.NONE,
        )
        sub = await js.subscribe(
            subject_filter or f"{stream}.>",
            config=consumer_config,
            stream=stream,
        )
        msgs = []
        for _ in range(n):
            try:
                msg = await asyncio.wait_for(sub.next_msg(), timeout=1.0)
                try:
                    data = json.loads(msg.data)
                except Exception:
                    data = msg.data.decode("utf-8", errors="replace")
                msgs.append({"subject": msg.subject, "data": data})
            except asyncio.TimeoutError:
                break
        await sub.unsubscribe()
        return {"ok": True, "stream": stream, "count": len(msgs), "messages": msgs}
    finally:
        await nc.drain()


async def _kv_get(bucket: str, key: str) -> dict:
    import nats
    nc = await nats.connect(nats_url(), **connect_options())
    try:
        js = nc.jetstream()
        kv = await js.key_value(bucket)
        entry = await kv.get(key)
        try:
            value = json.loads(entry.value)
        except Exception:
            value = entry.value.decode("utf-8", errors="replace")
        return {"ok": True, "bucket": bucket, "key": key, "value": value, "revision": entry.revision}
    except Exception as exc:
        return {"ok": False, "error": str(exc)}
    finally:
        await nc.drain()


async def _kv_put(bucket: str, key: str, value: Any) -> dict:
    import nats
    nc = await nats.connect(nats_url(), **connect_options())
    try:
        js = nc.jetstream()
        kv = await js.key_value(bucket)
        data = json.dumps(value).encode() if not isinstance(value, (bytes, str)) else \
               value.encode() if isinstance(value, str) else value
        revision = await kv.put(key, data)
        return {"ok": True, "bucket": bucket, "key": key, "revision": revision}
    except Exception as exc:
        return {"ok": False, "error": str(exc)}
    finally:
        await nc.drain()


async def _kv_list(bucket: str, keys_filter: str = ">") -> dict:
    import nats
    nc = await nats.connect(nats_url(), **connect_options())
    try:
        js = nc.jetstream()
        kv = await js.key_value(bucket)
        keys = []
        async for entry in kv.keys():
            keys.append(entry)
        return {"ok": True, "bucket": bucket, "keys": keys, "count": len(keys)}
    except Exception as exc:
        return {"ok": False, "error": str(exc)}
    finally:
        await nc.drain()


async def _server_info() -> dict:
    import nats
    nc = await nats.connect(nats_url(), **connect_options())
    try:
        info = nc._server_info  # type: ignore[attr-defined]
        return {"ok": True, "info": info}
    finally:
        await nc.drain()


async def _list_streams() -> dict:
    import nats
    nc = await nats.connect(nats_url(), **connect_options())
    try:
        js = nc.jetstream()
        streams = []
        async for info in js.streams_info():
            streams.append({
                "name": info.config.name,
                "subjects": info.config.subjects,
                "messages": info.state.messages,
                "bytes": info.state.bytes,
            })
        return {"ok": True, "streams": streams, "count": len(streams)}
    except Exception as exc:
        return {"ok": False, "error": str(exc)}
    finally:
        await nc.drain()


# ── Synchronous wrappers (called by Hermes tool runner) ──────────────────────

def nats_publish(subject: str, payload: str = "{}") -> str:
    """Publish a message to a NATS subject. payload is a JSON string."""
    result = _pool.submit(asyncio.run, _publish(subject, payload)).result(timeout=10)
    return json.dumps(result)


def nats_request(subject: str, payload: str = "{}", timeout: float = 5.0) -> str:
    """Send a request to a NATS subject and return the reply."""
    result = _pool.submit(asyncio.run, _request(subject, payload, timeout)).result(timeout=timeout + 5)
    return json.dumps(result)


def nats_peek(subject: str, n: int = 10) -> str:
    """Subscribe to a subject for 2 seconds and return up to n received messages."""
    result = _pool.submit(asyncio.run, _peek(subject, n)).result(timeout=15)
    return json.dumps(result)


def nats_jetstream_fetch(stream: str, subject_filter: str | None = None, n: int = 10) -> str:
    """Fetch up to n recent messages from a JetStream stream."""
    result = _pool.submit(asyncio.run, _jetstream_fetch(stream, subject_filter, n)).result(timeout=20)
    return json.dumps(result)


def nats_kv_get(bucket: str, key: str) -> str:
    """Read a key from a NATS KV bucket."""
    result = _pool.submit(asyncio.run, _kv_get(bucket, key)).result(timeout=10)
    return json.dumps(result)


def nats_kv_put(bucket: str, key: str, value: Any = None) -> str:
    """Write a value to a NATS KV bucket key."""
    result = _pool.submit(asyncio.run, _kv_put(bucket, key, value)).result(timeout=10)
    return json.dumps(result)


def nats_kv_list(bucket: str) -> str:
    """List all keys in a NATS KV bucket."""
    result = _pool.submit(asyncio.run, _kv_list(bucket)).result(timeout=15)
    return json.dumps(result)


def nats_server_info() -> str:
    """Return NATS server info (version, connections, cluster, etc.)."""
    result = _pool.submit(asyncio.run, _server_info()).result(timeout=10)
    return json.dumps(result)


def nats_list_streams() -> str:
    """List all JetStream streams with subject filters and message counts."""
    result = _pool.submit(asyncio.run, _list_streams()).result(timeout=15)
    return json.dumps(result)
