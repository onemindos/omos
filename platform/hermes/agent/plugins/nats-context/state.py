"""LiveState — thread-safe in-memory operational picture built from NATS streams."""
from __future__ import annotations

import copy
import threading
from collections import deque
from dataclasses import dataclass, field
from datetime import datetime, timezone
from typing import Any, Deque, Dict, List, Optional


def _now() -> float:
    return datetime.now(timezone.utc).timestamp()


@dataclass
class LiveState:
    _lock: threading.Lock = field(default_factory=threading.Lock, repr=False)

    # TAK / CoT tracks: callsign → latest payload
    tak: Dict[str, Dict[str, Any]] = field(default_factory=dict)

    # Entities by type → id → latest payload + _ts
    # types: drone, robot, vehicle, node, edge, camera, sensor, mission, task,
    #        orp, region, vessel, aircraft
    entities: Dict[str, Dict[str, Any]] = field(default_factory=dict)

    # Telemetry: subject → deque(maxlen=5) of recent readings
    telemetry: Dict[str, Deque[Any]] = field(default_factory=dict)

    # Intel items: deque(maxlen=50), newest first
    intel: Deque[Any] = field(default_factory=lambda: deque(maxlen=50))

    # Active threats: id → payload (auto-cleared on ack)
    threats: Dict[str, Any] = field(default_factory=dict)

    # Recent detections: deque(maxlen=20)
    detections: Deque[Any] = field(default_factory=lambda: deque(maxlen=20))

    # Agent heartbeats: "agent.owner.session" → payload + _ts
    agents: Dict[str, Any] = field(default_factory=dict)

    # Bridge status: name → payload + _ts
    bridges: Dict[str, Any] = field(default_factory=dict)

    # NATS server stats (from $SYS.REQ.SERVER.PING or sys.nats.stats)
    nats_stats: Optional[Dict[str, Any]] = None

    # Mission events ring buffer
    mission_events: Deque[Any] = field(default_factory=lambda: deque(maxlen=20))

    # Task events ring buffer
    task_events: Deque[Any] = field(default_factory=lambda: deque(maxlen=20))

    def update_tak(self, callsign: str, payload: Any) -> None:
        with self._lock:
            self.tak[callsign] = {**payload, "_ts": _now()}

    def update_entity(self, etype: str, eid: str, payload: Any) -> None:
        with self._lock:
            if etype not in self.entities:
                self.entities[etype] = {}
            self.entities[etype][eid] = {**payload, "_ts": _now()}

    def update_telemetry(self, subject: str, payload: Any) -> None:
        with self._lock:
            if subject not in self.telemetry:
                self.telemetry[subject] = deque(maxlen=5)
            self.telemetry[subject].appendleft({**payload, "_ts": _now()})

    def add_intel(self, payload: Any) -> None:
        with self._lock:
            self.intel.appendleft({**payload, "_ts": _now()})

    def update_threat(self, tid: str, payload: Any) -> None:
        with self._lock:
            if payload.get("acknowledged") or payload.get("acked"):
                self.threats.pop(tid, None)
            else:
                self.threats[tid] = {**payload, "_ts": _now()}

    def add_detection(self, subject: str, payload: Any) -> None:
        with self._lock:
            self.detections.appendleft({**payload, "subject": subject, "_ts": _now()})

    def update_agent(self, key: str, payload: Any) -> None:
        with self._lock:
            self.agents[key] = {**payload, "_ts": _now()}

    def update_bridge(self, name: str, payload: Any) -> None:
        with self._lock:
            self.bridges[name] = {**payload, "_ts": _now()}

    def update_nats_stats(self, payload: Any) -> None:
        with self._lock:
            self.nats_stats = {**payload, "_ts": _now()}

    def add_mission_event(self, payload: Any) -> None:
        with self._lock:
            self.mission_events.appendleft({**payload, "_ts": _now()})

    def add_task_event(self, payload: Any) -> None:
        with self._lock:
            self.task_events.appendleft({**payload, "_ts": _now()})

    def snapshot(self) -> dict:
        """Return a deep-copied thread-safe snapshot of all state."""
        with self._lock:
            return {
                "tak": copy.deepcopy(self.tak),
                "entities": copy.deepcopy(self.entities),
                "telemetry": {k: list(v) for k, v in self.telemetry.items()},
                "intel": list(self.intel),
                "threats": copy.deepcopy(self.threats),
                "detections": list(self.detections),
                "agents": copy.deepcopy(self.agents),
                "bridges": copy.deepcopy(self.bridges),
                "nats_stats": copy.deepcopy(self.nats_stats),
                "mission_events": list(self.mission_events),
                "task_events": list(self.task_events),
            }

    def prune_stale(
        self,
        entity_ttl: float = 300.0,
        agent_ttl: float = 120.0,
        bridge_ttl: float = 60.0,
        tak_ttl: float = 120.0,
        threat_ack_ttl: float = 300.0,
    ) -> None:
        """Remove entries that haven't been updated within their TTL."""
        now = _now()
        with self._lock:
            # TAK tracks
            for cs in list(self.tak):
                if now - self.tak[cs].get("_ts", 0) > tak_ttl:
                    del self.tak[cs]
            # Entities
            for etype in list(self.entities):
                for eid in list(self.entities[etype]):
                    if now - self.entities[etype][eid].get("_ts", 0) > entity_ttl:
                        del self.entities[etype][eid]
            # Agents
            for key in list(self.agents):
                if now - self.agents[key].get("_ts", 0) > agent_ttl:
                    del self.agents[key]
            # Bridges — mark offline instead of deleting
            for name in self.bridges:
                if now - self.bridges[name].get("_ts", 0) > bridge_ttl:
                    self.bridges[name]["_stale"] = True
