"""Stub OneMind services on a disposable nats-server, for the onemind-bus tool tests.

* entity service: ``svc.entity.get`` (NATS micro error headers on 404, as omos PR #552) and
  ``svc.entity.query`` (``{types, bbox, limit, cursor}`` pages, cursor = last entity_id)
* layer service: ``svc.layer.list``, ``svc.layer.meta``, ``svc.layer.snapshot`` (request-many,
  last chunk ``"last": true``, as omos PR #561)
* policy service: ``svc.policy.{command,task}.propose`` and ``svc.policy.action.get`` answering
  in the PR #554 reply shape; ``svc.policy.*.confirm`` is subscribed only to prove nobody calls it
* map client: applies ``cmd.map.>`` to a view state and acks on the reply subject
* spy: records every message on ``>`` so tests can assert what was (not) published

The real policy service (PR #554) is exercised in test_policy_integration.py.
"""
from __future__ import annotations

import asyncio
import json
import math
import threading
from pathlib import Path
from typing import Any, Dict, List, Optional

import nats


class LoopThread:
    def __init__(self):
        self.loop = asyncio.new_event_loop()
        self.t = threading.Thread(target=self.loop.run_forever, daemon=True)
        self.t.start()

    def run(self, coro, timeout=10):
        return asyncio.run_coroutine_threadsafe(coro, self.loop).result(timeout)

    def stop(self):
        self.loop.call_soon_threadsafe(self.loop.stop)
        self.t.join(5)


def _pos(e):
    loc = e.get("location") or {}
    if "lat" in loc:
        return loc["lat"], loc["lon"]
    ring = (((e.get("geo") or {}).get("coords") or {}).get("coordinates") or [[]])[0]
    if ring:
        return sum(p[1] for p in ring[:-1]) / (len(ring) - 1), sum(p[0] for p in ring[:-1]) / (len(ring) - 1)
    return None


def _in_bbox(e, b):
    p = _pos(e)
    if not p:
        return False
    if (e.get("geo") or {}).get("coords"):
        ring = e["geo"]["coords"]["coordinates"][0]
        return any(b[0] <= x <= b[2] and b[1] <= y <= b[3] for x, y in ring) or (b[1] <= p[0] <= b[3] and b[0] <= p[1] <= b[2])
    return b[1] <= p[0] <= b[3] and b[0] <= p[1] <= b[2]


class MapView:
    """What a map client shows; the D11 assertions read this state."""

    def __init__(self, visible: List[str]):
        self.camera: Dict[str, Any] = {"lat": 30.2672, "lon": -97.7431, "range_m": 5000.0, "heading_deg": 0, "pitch_deg": -45}
        self.visible = set(visible)
        self.selected: set = set()
        self.annotations: Dict[str, Dict[str, Any]] = {}
        self.tour: Optional[List[Dict[str, Any]]] = None
        self.applied: List[str] = []

    def apply(self, subject: str, p: Dict[str, Any]) -> Dict[str, Any]:
        parts = subject.split(".")
        if subject == "cmd.map.viewport.set":
            mode = p.get("mode", "fly_to")
            if mode == "fly_to":
                self.camera = {k: p[k] for k in ("lat", "lon", "range_m", "heading_deg", "pitch_deg")}
                self.tour = None
            elif mode == "frame":
                w, s, e, n = p["bbox"]
                diag = math.hypot((e - w) * 111320 * math.cos(math.radians((s + n) / 2)), (n - s) * 110540)
                self.camera = {"lat": (s + n) / 2, "lon": (w + e) / 2, "range_m": max(200.0, diag * 1.2),
                               "heading_deg": 0, "pitch_deg": -90, "bbox": p["bbox"]}
                self.tour = None
            elif mode == "tour":
                self.tour = p["stops"]
                last = p["stops"][-1]
                self.camera = {k: last[k] for k in ("lat", "lon", "range_m", "heading_deg", "pitch_deg")}
            else:
                raise ValueError(f"unknown viewport mode {mode}")
            out = {"camera": self.camera, "tour_stops": len(self.tour) if self.tour else 0}
        elif subject == "cmd.map.filter.set":
            self.visible |= set(p.get("show") or [])
            self.visible -= set(p.get("hide") or [])
            if "layers" in p:  # doctrine §9.7 full list form
                self.visible = set(p["layers"])
            out = {"visible_layers": sorted(self.visible)}
        elif parts[2] == "highlight":
            uid = parts[3]
            if p.get("clear"):
                self.selected.discard(uid)
            else:
                self.selected.add(uid)
            out = {"selected": sorted(self.selected)}
        elif parts[2] == "annotation":
            aid, verb = parts[3], parts[4]
            if verb == "create":
                self.annotations[aid] = {k: p.get(k) for k in ("kind", "label", "note", "geometry", "measure")}
            elif verb == "remove":
                if aid not in self.annotations:
                    raise KeyError(f"no annotation {aid}")
                del self.annotations[aid]
            out = {"annotations": sorted(self.annotations)}
        else:
            raise ValueError(f"unsupported {subject}")
        self.applied.append(subject)
        return out


class World:
    """All stubs on one background loop."""

    def __init__(self, url: str, fixture: Path, map_mode: str = "ack", client_id: str = "stub-map-1"):
        data = json.loads(Path(fixture).read_text())
        self.url = url
        self.entities: Dict[str, Dict[str, Any]] = {e["entity_id"]: e for e in data["entities"]}
        self.layers = data["layers"]
        self.view = MapView(data["initial_view"]["visible_layers"])
        self.map_mode = map_mode  # ack | silent | reject
        self.client_id = client_id
        self.proposals: List[Dict[str, Any]] = []
        self.confirm_calls: List[str] = []
        self.spy: List[tuple] = []
        self.map_msgs: List[tuple] = []
        self.lt = LoopThread()
        self.nc = None
        self._actions: Dict[str, Dict[str, Any]] = {}

    def start(self):
        self.lt.run(self._start())

    def stop(self):
        try:
            self.lt.run(self.nc.close())
        finally:
            self.lt.stop()

    def subjects_seen(self) -> List[str]:
        return [s for s, _ in self.spy]

    def flush(self):
        self.lt.run(self.nc.flush())

    async def _start(self):
        self.nc = await nats.connect(self.url, name="stub-world")
        nc = self.nc

        async def spy(m):
            self.spy.append((m.subject, m.data))

        await nc.subscribe(">", cb=spy)

        async def entity_get(m):
            q = json.loads(m.data or b"{}")
            e = self.entities.get(q.get("entity_id", ""))
            if not e:
                await nc.publish(m.reply, b"", headers={"Nats-Service-Error-Code": "404", "Nats-Service-Error": "entity not found"})
                return
            await m.respond(json.dumps(e).encode())

        async def entity_query(m):
            q = json.loads(m.data or b"{}")
            ents = sorted(self.entities.values(), key=lambda e: e["entity_id"])
            if q.get("types"):
                ents = [e for e in ents if e["entity_type"] in q["types"]]
            if q.get("bbox"):
                ents = [e for e in ents if _in_bbox(e, q["bbox"])]
            if q.get("cursor"):
                ents = [e for e in ents if e["entity_id"] > q["cursor"]]
            limit = min(int(q.get("limit") or 200), 200)
            page, more = ents[:limit], len(ents) > limit
            await m.respond(json.dumps({"entities": page, "cursor": page[-1]["entity_id"] if more else None}).encode())

        def counts(l):
            return sum(1 for e in self.entities.values() if e["entity_type"] in l["entity_types"])

        async def layer_list(m):
            await m.respond(json.dumps({"v": 1, "layers": [
                {"layer": l["layer"], "title": l["title"], "tier": l["tier"], "count": counts(l), "status": "fresh"}
                for l in self.layers]}).encode())

        async def layer_meta(m):
            q = json.loads(m.data or b"{}")
            l = next((l for l in self.layers if l["layer"] == q.get("layer")), None)
            if not l:
                await m.respond(json.dumps({"v": 1, "error": "unknown_layer", "detail": str(q.get("layer"))}).encode())
                return
            await m.respond(json.dumps({"v": 1, **l, "count": counts(l), "as_of": 1793275200000, "cap": 5000,
                                        "freshness": {"status": "fresh", "limit_s": 60, "feeds": []}}).encode())

        async def layer_snapshot(m):
            q = json.loads(m.data or b"{}")
            l = next((l for l in self.layers if l["layer"] == q.get("layer")), None)
            if not l:
                await m.respond(json.dumps({"v": 1, "error": "unknown_layer"}).encode())
                return
            aoi = q.get("aoi") or {}
            if not aoi and l["tier"] == "large":
                await m.respond(json.dumps({"v": 1, "error": "aoi_required"}).encode())
                return
            items = [e for e in self.entities.values() if e["entity_type"] in l["entity_types"]]
            if aoi.get("bbox"):
                items = [e for e in items if _in_bbox(e, aoi["bbox"])]
            render = [{"entity_id": e["entity_id"], "entity_type": e["entity_type"], "callsign": e.get("callsign"),
                       "lat": (_pos(e) or (None, None))[0], "lon": (_pos(e) or (None, None))[1], "t": e["t"],
                       "simulated": e.get("indicators", {}).get("simulated")} for e in items]
            mx = int(q.get("max") or 5000)
            trunc = len(render) > mx
            render = render[:mx]
            chunks = [render[i:i + 2] for i in range(0, len(render), 2)] or [[]]
            for i, ch in enumerate(chunks):
                await nc.publish(m.reply, json.dumps({"v": 1, "layer": l["layer"], "part": i, "last": i == len(chunks) - 1,
                                                      "mode": "entities", "ent_seq": 42, "as_of": 1793275200000,
                                                      "truncated": trunc, "items": ch}).encode())

        async def propose(m):
            fam = m.subject.split(".")[2]
            req = json.loads(m.data)
            self.proposals.append({"family": fam, **req})
            refs = [req.get("target", {}).get("entity_id")]
            c = (req.get("params") or {}).get("center")
            if isinstance(c, dict) and c.get("entity_id"):
                refs.append(c["entity_id"])
            aid = f"act-{len(self.proposals):04d}"
            if any((self.entities.get(r) or {}).get("kind") == "person" for r in refs):
                rep = {"action_id": aid, "status": "denied", "approvals_needed": 0,
                       "denial": {"code": "forbidden_action", "message": "actions never target a person"}}
            else:
                rep = {"action_id": aid, "status": "awaiting_confirmation", "level": "A3", "params_hash": "sha256:x",
                       "confirm_by": "2026-10-29T12:02:00Z", "approvals_needed": 1, "denial": None}
            self._actions[aid] = {"action_id": aid, "action_type": req["action_type"], **rep}
            await m.respond(json.dumps(rep).encode())

        async def action_get(m):
            q = json.loads(m.data)
            a = self._actions.get(q.get("action_id"))
            await m.respond(json.dumps({"action": a, "denial": None} if a else
                                       {"action": None, "denial": {"code": "not_authorized", "message": "unknown"}}).encode())

        async def confirm(m):
            self.confirm_calls.append(m.subject)
            await m.respond(json.dumps({"denial": {"code": "agent_cannot_confirm"}}).encode())

        async def map_client(m):
            self.map_msgs.append((m.subject, m.data))
            if self.map_mode == "silent" or not m.reply:
                return
            p = json.loads(m.data)
            if self.map_mode == "reject":
                await nc.publish(m.reply, json.dumps({"v": 1, "id": p.get("id"), "ok": False, "client": self.client_id,
                                                      "error": "unsupported"}).encode())
                return
            try:
                applied = self.view.apply(m.subject, p)
                rep = {"v": 1, "id": p.get("id"), "ok": True, "client": self.client_id, "applied": applied}
            except Exception as e:  # the client could not apply it: it says so
                rep = {"v": 1, "id": p.get("id"), "ok": False, "client": self.client_id, "error": str(e)}
            await nc.publish(m.reply, json.dumps(rep).encode())

        for subj, cb in (("svc.entity.get", entity_get), ("svc.entity.query", entity_query),
                         ("svc.layer.list", layer_list), ("svc.layer.meta", layer_meta),
                         ("svc.layer.snapshot", layer_snapshot), ("svc.policy.command.propose", propose),
                         ("svc.policy.task.propose", propose), ("svc.policy.action.get", action_get),
                         ("svc.policy.*.confirm", confirm), ("cmd.map.>", map_client)):
            await nc.subscribe(subj, cb=cb)
        await nc.flush()


def RecordingClient(config):
    """A BusClient that records every subject it sends to (built lazily from the loaded plugin)."""
    import sys

    BusClient = sys.modules["hermes_plugins.onemind_bus.bus"].BusClient

    class _Recording(BusClient):
        def __init__(self, cfg):
            super().__init__(cfg)
            self.sent: List[str] = []

        def request(self, subject, body, timeout=None):
            self.sent.append(subject)
            return super().request(subject, body, timeout)

        def request_many(self, subject, body, max_parts=16, timeout=None):
            self.sent.append(subject)
            return super().request_many(subject, body, max_parts, timeout)

        def map_action(self, subject, payload, window=None):
            self.sent.append(subject)
            return super().map_action(subject, payload, window)

    return _Recording(config)
