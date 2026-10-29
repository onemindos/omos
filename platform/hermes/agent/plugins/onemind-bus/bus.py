"""NATS access for the onemind-bus Hermes plugin.

One connection, owned by a private asyncio loop on a daemon thread, so the synchronous Hermes
tool handlers can call the bus whether or not the caller already runs an event loop.

The client can only reach the subjects in ``ALLOWED_REQUESTS`` (request/reply) and
``MAP_SUBJECT`` (reversible UI actions). Everything else raises ``SubjectNotAllowed`` before
anything is sent. There is no generic publish, no ``cmd.*`` other than ``cmd.map.*``, and no
``svc.policy.*.confirm``: an agent proposes, a human confirms in the OneMind UI.
"""
from __future__ import annotations

import asyncio
import json
import os
import re
import threading
from dataclasses import dataclass, field
from typing import Any, Dict, List, Optional

# Request/reply endpoints Hermes may call (onemind-xr docs/design/BUS-DESIGN.md section 3.4 and
# the AGENTS imports in section 7). svc.layer.snapshot is proposal I2-P1 (bounded read).
ALLOWED_REQUESTS = frozenset({
    "svc.entity.get",
    "svc.entity.query",
    "svc.layer.list",
    "svc.layer.meta",
    "svc.layer.snapshot",
    "svc.policy.command.propose",
    "svc.policy.task.propose",
    "svc.policy.action.get",
})

# Reversible UI actions (omos docs/NATS-DOCTRINE.md section 9.1 and 9.7). Only these forms:
#   cmd.map.viewport.set | cmd.map.filter.set | cmd.map.highlight.{uid}
#   cmd.map.annotation.{id}.{create|remove}
_TOKEN = r"[A-Za-z0-9_-]{1,64}"
MAP_SUBJECT = re.compile(
    rf"^cmd\.map\.(viewport\.set|filter\.set|highlight\.{_TOKEN}|annotation\.{_TOKEN}\.(create|remove))$"
)


class SubjectNotAllowed(Exception):
    """Raised before sending to a subject outside the plugin allowlist."""


class BusError(Exception):
    """The bus could not answer (no responders, timeout, disconnected)."""


def check_request_subject(subject: str) -> None:
    if subject not in ALLOWED_REQUESTS:
        raise SubjectNotAllowed(f"request subject not allowed: {subject}")


def check_map_subject(subject: str) -> None:
    if not MAP_SUBJECT.match(subject):
        raise SubjectNotAllowed(f"map subject not allowed: {subject}")


@dataclass
class BusConfig:
    url: str = "nats://localhost:4222"
    user: Optional[str] = None
    password: Optional[str] = None
    creds: Optional[str] = None
    timeout_s: float = 5.0
    map_ack_window_s: float = 1.5
    name: str = "hermes-onemind-bus"

    @classmethod
    def from_env(cls, env: Optional[Dict[str, str]] = None) -> "BusConfig":
        e = os.environ if env is None else env
        return cls(
            url=e.get("ONEMIND_BUS_URL") or e.get("NATS_URL") or "nats://localhost:4222",
            user=e.get("ONEMIND_BUS_USER") or None,
            password=e.get("ONEMIND_BUS_PASSWORD") or None,
            creds=e.get("ONEMIND_BUS_CREDS") or None,
            timeout_s=float(e.get("ONEMIND_BUS_TIMEOUT_S") or 5.0),
            map_ack_window_s=float(e.get("ONEMIND_BUS_MAP_ACK_WINDOW_S") or 1.5),
        )


@dataclass
class MapResult:
    acks: List[Dict[str, Any]] = field(default_factory=list)       # ok acks from clients
    rejections: List[Dict[str, Any]] = field(default_factory=list)  # client replies with ok=false
    stored: Optional[Dict[str, Any]] = None                         # JetStream PubAck, if a stream captured it


class BusClient:
    """Thread-safe synchronous facade over one nats-py connection."""

    def __init__(self, config: BusConfig):
        self.config = config
        self._loop: Optional[asyncio.AbstractEventLoop] = None
        self._thread: Optional[threading.Thread] = None
        self._nc = None
        self._lock = threading.Lock()

    # -- lifecycle ---------------------------------------------------------------------------

    def _ensure(self):
        with self._lock:
            if self._loop is None:
                loop = asyncio.new_event_loop()
                t = threading.Thread(target=loop.run_forever, name="onemind-bus", daemon=True)
                t.start()
                self._loop, self._thread = loop, t
            if self._nc is None or self._nc.is_closed:
                self._nc = self._run(self._connect(), self.config.timeout_s + 5)
            return self._nc

    async def _connect(self):
        import nats  # nats-py, installed into /gateway-deps by the Hermes init container

        opts: Dict[str, Any] = {
            "servers": [self.config.url],
            "name": self.config.name,
            "connect_timeout": max(1, int(self.config.timeout_s)),
            "max_reconnect_attempts": -1,
            "allow_reconnect": True,
        }
        if self.config.creds:
            opts["user_credentials"] = self.config.creds
        elif self.config.user:
            opts["user"] = self.config.user
            opts["password"] = self.config.password or ""
        return await nats.connect(**opts)

    def _run(self, coro, timeout: float):
        assert self._loop is not None
        fut = asyncio.run_coroutine_threadsafe(coro, self._loop)
        try:
            return fut.result(timeout)
        except TimeoutError:
            fut.cancel()
            raise BusError("timeout")

    def close(self) -> None:
        with self._lock:
            if self._nc is not None and self._loop is not None:
                try:
                    self._run(self._nc.close(), 5)
                except Exception:
                    pass
            self._nc = None
            if self._loop is not None:
                self._loop.call_soon_threadsafe(self._loop.stop)
                if self._thread is not None:
                    self._thread.join(5)
                self._loop.close()
            self._loop, self._thread = None, None

    # -- request/reply -----------------------------------------------------------------------

    def request(self, subject: str, body: Dict[str, Any], timeout: Optional[float] = None) -> Dict[str, Any]:
        """One request, one JSON reply. NATS micro errors come back as ``{"error": {...}}``."""
        check_request_subject(subject)
        nc = self._ensure()
        t = timeout or self.config.timeout_s
        return self._run(self._request(nc, subject, body, t), t + 2)

    async def _request(self, nc, subject, body, t):
        from nats.errors import NoRespondersError, TimeoutError as NatsTimeout

        try:
            msg = await nc.request(subject, json.dumps(body).encode(), timeout=t)
        except NoRespondersError:
            raise BusError(f"no responders on {subject}")
        except NatsTimeout:
            raise BusError(f"timeout on {subject}")
        return _decode_reply(msg)

    def request_many(self, subject: str, body: Dict[str, Any], max_parts: int = 16,
                     timeout: Optional[float] = None) -> List[Dict[str, Any]]:
        """Request-many reply terminated by a chunk with ``"last": true`` (or an error)."""
        check_request_subject(subject)
        nc = self._ensure()
        t = timeout or self.config.timeout_s
        return self._run(self._request_many(nc, subject, body, t, max_parts), t + 2)

    async def _request_many(self, nc, subject, body, t, max_parts):
        inbox = nc.new_inbox()
        q: asyncio.Queue = asyncio.Queue()

        async def cb(m):
            await q.put(m)

        sub = await nc.subscribe(inbox, cb=cb)
        try:
            await nc.publish(subject, json.dumps(body).encode(), reply=inbox)
            parts: List[Dict[str, Any]] = []
            while len(parts) < max_parts:
                try:
                    m = await asyncio.wait_for(q.get(), t)
                except asyncio.TimeoutError:
                    if not parts:
                        raise BusError(f"timeout on {subject}")
                    parts.append({"error": "incomplete", "detail": "stream ended without last chunk"})
                    break
                if m.headers and m.headers.get("Status") == "503":
                    raise BusError(f"no responders on {subject}")
                d = _decode_reply(m)
                parts.append(d)
                if d.get("last") or d.get("error"):
                    break
            else:
                # max_parts reached without a last:true chunk: say so instead of returning a silent prefix.
                parts.append({"truncated": True, "detail": f"stopped after {max_parts} parts without last chunk"})
            return parts
        finally:
            await sub.unsubscribe()

    # -- reversible map actions --------------------------------------------------------------

    def map_action(self, subject: str, payload: Dict[str, Any], window: Optional[float] = None) -> MapResult:
        """Publish one ``cmd.map.*`` action with a reply inbox and collect client acks.

        A client that applied the action replies ``{"v":1,"id":..,"ok":true,"client":..,
        "applied":{..}}``. A JetStream PubAck (``{"stream":..,"seq":..}``) from a stream that
        captures ``cmd.>`` is recorded as the audit pointer, never counted as an ack.
        """
        check_map_subject(subject)
        nc = self._ensure()
        w = window if window is not None else self.config.map_ack_window_s
        return self._run(self._map_action(nc, subject, payload, w), w + self.config.timeout_s)

    async def _map_action(self, nc, subject, payload, window):
        inbox = nc.new_inbox()
        res = MapResult()
        first_ok = asyncio.Event()

        async def cb(m):
            try:
                d = json.loads(m.data or b"{}")
            except ValueError:
                return
            if not isinstance(d, dict):
                return
            if "stream" in d and "seq" in d and "ok" not in d:
                res.stored = {"stream": d.get("stream"), "seq": d.get("seq")}
                return
            if d.get("id") not in (None, payload.get("id")):
                return
            if d.get("ok") is True:
                res.acks.append(d)
                first_ok.set()
            else:
                res.rejections.append(d)

        sub = await nc.subscribe(inbox, cb=cb)
        try:
            headers = {"Nats-Msg-Id": str(payload.get("id", ""))}
            await nc.publish(subject, json.dumps(payload).encode(), reply=inbox, headers=headers)
            await nc.flush()
            try:
                await asyncio.wait_for(first_ok.wait(), window)
                await asyncio.sleep(min(0.25, window))  # let other clients' acks arrive
            except asyncio.TimeoutError:
                pass
            return res
        finally:
            await sub.unsubscribe()


def _decode_reply(msg) -> Dict[str, Any]:
    hdr = msg.headers or {}
    code = hdr.get("Nats-Service-Error-Code")
    if code:
        return {"error": {"code": code, "message": hdr.get("Nats-Service-Error", "")}}
    if hdr.get("Status") == "503":
        return {"error": {"code": "no_responders", "message": "no responders"}}
    try:
        d = json.loads(msg.data or b"{}")
    except ValueError:
        return {"error": {"code": "bad_reply", "message": "reply is not JSON"}}
    return d if isinstance(d, dict) else {"value": d}
