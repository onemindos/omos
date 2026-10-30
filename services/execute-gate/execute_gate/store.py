"""Approval storage: a tiny revisioned key-value interface.

The in-memory store backs tests and single-process use; NatsKvStore keeps pending approvals across
restarts in the JetStream KV bucket `exec-approvals`. Both give compare-and-set so consuming an
approval is atomic.
"""

from __future__ import annotations

import asyncio
import threading
from typing import Optional


class MemoryStore:
    def __init__(self) -> None:
        self._d = {}
        self._rev = 0
        self._lock = threading.Lock()

    def get(self, key) -> Optional[tuple]:
        with self._lock:
            return self._d.get(key)

    def create(self, key, value) -> bool:
        with self._lock:
            if key in self._d:
                return False
            self._rev += 1
            self._d[key] = (value, self._rev)
            return True

    def update(self, key, value, revision) -> bool:
        with self._lock:
            cur = self._d.get(key)
            if cur is None or cur[1] != revision:
                return False
            self._rev += 1
            self._d[key] = (value, self._rev)
            return True

    def delete(self, key, revision) -> bool:
        with self._lock:
            cur = self._d.get(key)
            if cur is None or cur[1] != revision:
                return False
            del self._d[key]
            return True

    def healthy(self) -> bool:
        return True


class NatsKvStore:
    """JetStream KV store plus a core publisher, on a private asyncio loop thread.

    Not exercised by the unit tests (needs a live NATS); see README for the Phase 4 verification.
    """

    def __init__(self, url, user, password, bucket, timeout: float = 5.0):
        self._timeout = timeout
        self._bucket = bucket
        self._nc = None
        self._kv = None
        self._loop = asyncio.new_event_loop()
        self._thread = threading.Thread(target=self._loop.run_forever, name="nats-loop", daemon=True)
        self._thread.start()
        self._run(self._connect(url, user, password))

    def _run(self, coro):
        return asyncio.run_coroutine_threadsafe(coro, self._loop).result(self._timeout)

    async def _connect(self, url, user, password):
        import nats

        kwargs = {"name": "execute-gate", "max_reconnect_attempts": -1}
        if user:
            kwargs.update(user=user, password=password)
        self._nc = await nats.connect(url, **kwargs)
        self._kv = await self._nc.jetstream().key_value(self._bucket)

    def get(self, key):
        async def go():
            from nats.js.errors import KeyNotFoundError

            try:
                e = await self._kv.get(key)
            except KeyNotFoundError:
                return None
            if e is None or e.value is None:
                return None
            return bytes(e.value), int(e.revision)

        return self._run(go())

    def create(self, key, value):
        async def go():
            from nats.js.errors import KeyWrongLastSequenceError

            try:
                await self._kv.create(key, value)
                return True
            except KeyWrongLastSequenceError:
                return False

        return self._run(go())

    def update(self, key, value, revision):
        async def go():
            from nats.js.errors import KeyWrongLastSequenceError

            try:
                await self._kv.update(key, value, last=revision)
                return True
            except KeyWrongLastSequenceError:
                return False

        return self._run(go())

    def delete(self, key, revision):
        async def go():
            from nats.js.errors import KeyWrongLastSequenceError

            try:
                await self._kv.delete(key, last=revision)
                return True
            except KeyWrongLastSequenceError:
                return False

        return self._run(go())

    def publish(self, subject, payload) -> None:
        async def go():
            await self._nc.publish(subject, payload)

        self._run(go())

    def healthy(self) -> bool:
        return bool(self._nc and self._nc.is_connected)
