"""Background NATS subscription thread — continuously feeds LiveState."""
from __future__ import annotations

import asyncio
import json
import logging
import threading
import time
from typing import TYPE_CHECKING

import nats

from .conn import connect_options

if TYPE_CHECKING:
    from .state import LiveState

logger = logging.getLogger("nats-context.subscriptions")

# Full subject taxonomy — every wildcard we care about
SUBJECTS = [
    "tak.cot.>",          # TAK / CoT position reports
    "ent.>",              # Entity state (all types)
    "tel.>",              # Telemetry feeds
    "intel.>",            # Intelligence reports
    "evt.threat.>",       # Threat events
    "det.>",              # Detections
    "cmd.alert.>",        # Alert acknowledgements
    "agents.hb.>",        # Agent heartbeats
    "svc.bridge.>",       # Bridge health
    "sys.>",              # System / NATS stats
    "evt.mission.>",      # Mission events
    "evt.task.>",         # Task events
    "evt.>",              # All other events (catch-all, after specifics)
]

PRUNE_INTERVAL = 30.0  # seconds between stale-entry sweeps
STATS_INTERVAL = 60.0  # seconds between $SYS.REQ.SERVER.PING polls


class SubscriptionManager:
    """Runs a background asyncio loop and feeds LiveState from NATS messages."""

    def __init__(self, state: "LiveState", nats_url: str = "nats://localhost:4222"):
        self._state = state
        self._nats_url = nats_url
        self._thread: threading.Thread | None = None
        self._loop: asyncio.AbstractEventLoop | None = None
        self._nc: nats.aio.client.Client | None = None
        self._running = False

    def start(self) -> None:
        if self._running:
            return
        self._running = True
        self._thread = threading.Thread(
            target=self._run_loop,
            name="nats-context-bg",
            daemon=True,
        )
        self._thread.start()
        logger.info("nats-context: background subscription thread started")

    def stop(self) -> None:
        self._running = False
        if self._loop and not self._loop.is_closed():
            asyncio.run_coroutine_threadsafe(self._shutdown(), self._loop)

    def _run_loop(self) -> None:
        self._loop = asyncio.new_event_loop()
        asyncio.set_event_loop(self._loop)
        try:
            self._loop.run_until_complete(self._main())
        except Exception as exc:
            logger.error("nats-context: loop crashed: %s", exc, exc_info=True)
        finally:
            self._loop.close()

    async def _main(self) -> None:
        retry_delay = 5.0
        while self._running:
            try:
                await self._connect_and_subscribe()
            except Exception as exc:
                logger.warning("nats-context: connection error: %s — retry in %.0fs", exc, retry_delay)
                await asyncio.sleep(retry_delay)
                retry_delay = min(retry_delay * 1.5, 60.0)
            else:
                retry_delay = 5.0

    async def _connect_and_subscribe(self) -> None:
        self._nc = await nats.connect(
            self._nats_url,
            reconnect_time_wait=5,
            max_reconnect_attempts=-1,
            error_cb=self._on_error,
            disconnected_cb=self._on_disconnect,
            reconnected_cb=self._on_reconnect,
            **connect_options(),  # OneMind: credentials, CA and client name from the environment
        )
        logger.info("nats-context: connected to %s", self._nats_url)

        # Subscribe to all subjects
        for subj in SUBJECTS:
            await self._nc.subscribe(subj, cb=self._make_handler(subj))

        # Start background tasks
        prune_task = asyncio.create_task(self._prune_loop())
        stats_task = asyncio.create_task(self._stats_loop())

        # Wait until stopped
        try:
            while self._running:
                await asyncio.sleep(1.0)
        finally:
            prune_task.cancel()
            stats_task.cancel()
            await self._nc.drain()

    def _make_handler(self, pattern: str):
        state = self._state

        async def handler(msg):
            try:
                payload = json.loads(msg.data.decode("utf-8", errors="replace"))
            except json.JSONDecodeError:
                payload = {"raw": msg.data.decode("utf-8", errors="replace")}
            _route_message(state, msg.subject, payload)

        return handler

    async def _prune_loop(self) -> None:
        while True:
            await asyncio.sleep(PRUNE_INTERVAL)
            try:
                self._state.prune_stale()
            except Exception as exc:
                logger.debug("prune error: %s", exc)

    async def _stats_loop(self) -> None:
        while True:
            await asyncio.sleep(STATS_INTERVAL)
            if not self._nc or self._nc.is_closed:
                continue
            try:
                # Request server stats from $SYS subject
                resp = await self._nc.request("$SYS.REQ.SERVER.PING", b"{}", timeout=3.0)
                payload = json.loads(resp.data.decode())
                self._state.update_nats_stats(payload)
            except Exception as exc:
                logger.debug("stats poll error: %s", exc)

    async def _shutdown(self) -> None:
        if self._nc and not self._nc.is_closed:
            await self._nc.drain()

    async def _on_error(self, exc: Exception) -> None:
        logger.warning("nats-context: nats error: %s", exc)

    async def _on_disconnect(self) -> None:
        logger.warning("nats-context: disconnected from NATS")

    async def _on_reconnect(self) -> None:
        logger.info("nats-context: reconnected to NATS")


# ── Message routing ──────────────────────────────────────────────────────────

def _route_message(state: "LiveState", subject: str, payload: dict) -> None:
    parts = subject.split(".")

    # tak.cot.<callsign> or tak.cot.<callsign>.<suffix>
    if parts[0] == "tak" and len(parts) >= 3:
        callsign = parts[2]
        state.update_tak(callsign, payload)
        return

    # ent.<type>.<id>[.<action>]
    if parts[0] == "ent" and len(parts) >= 3:
        etype, eid = parts[1], parts[2]
        action = parts[3] if len(parts) > 3 else "update"
        if etype == "orp":
            state.update_entity("orp", eid, {**payload, "action": action})
        elif etype == "region":
            state.update_entity("region", eid, {**payload, "action": action})
        else:
            state.update_entity(etype, eid, {**payload, "action": action})
        return

    # tel.<platform>.<id>
    if parts[0] == "tel":
        state.update_telemetry(subject, payload)
        return

    # intel.*.*
    if parts[0] == "intel":
        state.add_intel({**payload, "subject": subject})
        return

    # evt.threat.*
    if parts[0] == "evt" and len(parts) >= 2 and parts[1] == "threat":
        tid = parts[2] if len(parts) > 2 else "unknown"
        state.update_threat(tid, payload)
        return

    # evt.mission.*
    if parts[0] == "evt" and len(parts) >= 2 and parts[1] == "mission":
        state.add_mission_event({**payload, "subject": subject})
        return

    # evt.task.*
    if parts[0] == "evt" and len(parts) >= 2 and parts[1] == "task":
        state.add_task_event({**payload, "subject": subject})
        return

    # det.*.*
    if parts[0] == "det":
        state.add_detection(subject, payload)
        return

    # cmd.alert.<id>.ack → clear threat
    if parts[0] == "cmd" and len(parts) >= 4 and parts[1] == "alert" and parts[3] == "ack":
        state.update_threat(parts[2], {"acknowledged": True})
        return

    # agents.hb.<agent>.<owner>.<session>
    if parts[0] == "agents" and len(parts) >= 2 and parts[1] == "hb":
        key = ".".join(parts[2:]) if len(parts) > 2 else "unknown"
        state.update_agent(key, payload)
        return

    # svc.bridge.<name>
    if parts[0] == "svc" and len(parts) >= 3 and parts[1] == "bridge":
        state.update_bridge(parts[2], payload)
        return

    # sys.nats.stats or any sys.> stats
    if parts[0] == "sys":
        state.update_nats_stats({**payload, "subject": subject})
        return
