"""Tool-level tests of the onemind-bus plugin against stub services on a disposable nats-server."""
from __future__ import annotations

import json
import re
import time

import pytest

from conftest import FakeCtx, load_plugin

P = load_plugin()
ALLOWED_PREFIXES = ("svc.entity.", "svc.layer.", "svc.policy.command.propose", "svc.policy.task.propose",
                    "svc.policy.action.get", "cmd.map.")


# -- no bus needed ------------------------------------------------------------------------------

def test_snapshot_tool_is_off_and_unregistered_by_default():
    # I2-P1 is unruled and the A19 hermes user cannot publish svc.layer.snapshot
    ctx = FakeCtx()
    P.register(ctx, env={"ONEMIND_BUS_URL": "nats://127.0.0.1:1"})
    assert "onemind_layer_snapshot" not in ctx.tools
    assert set(ctx.tools) == set(P.tools.TOOL_NAMES) - {"onemind_layer_snapshot"}
    ctx2 = FakeCtx()
    P.register(ctx2, env={"ONEMIND_BUS_URL": "nats://127.0.0.1:1", "ONEMIND_BUS_LAYER_SNAPSHOT": "0"})
    assert "onemind_layer_snapshot" not in ctx2.tools


def test_registers_the_expected_tools_and_nothing_that_confirms_or_publishes():
    ctx = FakeCtx()
    P.register(ctx, env={"ONEMIND_BUS_URL": "nats://127.0.0.1:1", "ONEMIND_BUS_LAYER_SNAPSHOT": "1"})
    names = set(ctx.tools)
    assert names == set(P.tools.TOOL_NAMES)
    assert len(names) == 16
    assert all(ctx.tools[n]["toolset"] == "onemind-bus" for n in names)
    for n in names:
        assert "confirm" not in n and "publish" not in n
        props = ctx.tools[n]["schema"]["parameters"].get("properties", {})
        assert "subject" not in props, f"{n} lets the model choose a subject"
    manifest = (P.tools.__file__.rsplit("/", 1)[0]) + "/plugin.yaml"
    import yaml

    assert set(yaml.safe_load(open(manifest))["provides_tools"]) == names


def test_subject_allowlist_refuses_everything_else():
    bus = P.bus
    for s in ("svc.policy.command.confirm", "svc.policy.task.confirm", "svc.policy.check", "cmd.fleet.drone.sim-x1.goto",
              "cmd.mission.t1.create", "cmd.entity.x.override", "tak.cot.inject", "ent.drone.sim-x1.state",
              "svc.layer.watch", "svc.entity.history", "evt.so.action.approved", ">", "svc.entity.>"):
        with pytest.raises(bus.SubjectNotAllowed):
            bus.check_request_subject(s)
    for s in ("cmd.fleet.drone.sim-x1.goto", "cmd.map.geofence.z1.create", "cmd.map.mission.draw", "cmd.map.>",
              "cmd.map.highlight.a.b", "cmd.map.annotation.x.update", "cmd.map.viewport.set.x", "cmd.map.highlight.*"):
        with pytest.raises(bus.SubjectNotAllowed):
            bus.check_map_subject(s)
    for s in ("cmd.map.viewport.set", "cmd.map.filter.set", "cmd.map.highlight.sim-x1",
              "cmd.map.annotation.hermes-1.create", "cmd.map.annotation.hermes-1.remove"):
        bus.check_map_subject(s)
    # geofence edits are A2 (policy), so they are not a direct map action for the agent
    assert not bus.MAP_SUBJECT.match("cmd.map.geofence.zone-a.create")


def test_retire_flag_blocks_nats_publish_by_default_and_can_be_turned_off():
    ctx = FakeCtx()
    P.register(ctx, env={})
    msg = ctx.blocked("nats_publish", {"subject": "cmd.fleet.drone.sim-x1.goto", "payload": "{}"})
    assert msg and "retired" in msg
    assert ctx.blocked("onemind_entity_get") is None
    ctx2 = FakeCtx()
    P.register(ctx2, env={"ONEMIND_BUS_RETIRE_NATS_PUBLISH": "0"})
    assert ctx2.blocked("nats_publish") is None
    ctx3 = FakeCtx()
    P.register(ctx3, env={"ONEMIND_BUS_RETIRED_TOOLS": "tak_publish_cot"})
    assert ctx3.blocked("tak_publish_cot") and ctx3.blocked("nats_publish")


def test_geo_helpers():
    g = P.geo
    d = g.haversine_m(30.2700, -97.7400, 30.2740, -97.7445)
    assert 600 < d < 640
    a = g.polygon_area_m2([(0, 0), (0, 0.001), (0.001, 0.001), (0.001, 0)])
    assert 12_000 < a < 12_500  # ~111 m x 111 m
    b = g.bbox_around(30.0, -97.0, 1000)
    assert b[0] < -97.0 < b[2] and b[1] < 30.0 < b[3]
    assert g.valid_bbox([-1, -1, 1, 1]) and not g.valid_bbox([1, 1, -1, -1]) and not g.valid_bbox("x")


def test_bus_down_is_reported_not_raised(plugin):
    ctx = FakeCtx()
    P.register(ctx, env={"ONEMIND_BUS_URL": "nats://127.0.0.1:9", "ONEMIND_BUS_TIMEOUT_S": "1"})
    r = ctx.dispatch("onemind_entity_get", {"entity_id": "sim-x1"})
    assert r["ok"] is False and r["error"] in ("bus_unavailable", "internal_error")


# -- reads --------------------------------------------------------------------------------------

def test_entity_get(hermes):
    r = hermes.dispatch("onemind_entity_get", {"entity_id": "sim-x1"})
    e = r["entity"]
    assert e["entity_id"] == "sim-x1" and e["entity_type"] == "drone" and e["simulated"] is True
    assert e["lat"] == 30.27 and e["sources"] == ["mavlink-sim"] and "note" in r
    assert hermes.dispatch("onemind_entity_get", {"entity_id": "nope"})["error"] == "not_found"
    assert hermes.dispatch("onemind_entity_get", {"entity_id": "bad id.>"})["error"] == "invalid_argument"


def test_entity_query_near_an_area_moving_only(hermes):
    r = hermes.dispatch("onemind_entity_query", {"near": {"entity_id": "zone-north-field", "radius_m": 500},
                                                 "moving_only": True})
    ids = [e["entity_id"] for e in r["entities"]]
    assert "sim-tractor-1" in ids and "trk-p1" in ids
    assert "sim-x2" not in ids  # not moving
    assert "icao-a1b2c3" not in ids  # 20 km away
    assert "sim-truck-9" not in ids  # parked
    assert all("distance_m" in e for e in r["entities"])
    assert [e["distance_m"] for e in r["entities"]] == sorted(e["distance_m"] for e in r["entities"])


def test_entity_query_by_name_and_type(hermes):
    r = hermes.dispatch("onemind_entity_query", {"name": "barn"})
    assert [e["entity_id"] for e in r["entities"]] == ["structure-barn"]
    r = hermes.dispatch("onemind_entity_query", {"types": ["drone"]})
    assert sorted(e["entity_id"] for e in r["entities"]) == ["sim-x1", "sim-x2", "sim-x3"]
    assert hermes.dispatch("onemind_entity_query", {"bbox": [1, 2]})["error"] == "invalid_argument"
    assert hermes.dispatch("onemind_entity_query", {"limit": 500})["error"] == "invalid_argument"


def test_layer_list_meta_snapshot(hermes):
    layers = {l["layer"] for l in hermes.dispatch("onemind_layer_list", {})["layers"]}
    assert {"fleet", "air", "zones"} <= layers
    m = hermes.dispatch("onemind_layer_meta", {"layer": "fleet"})["meta"]
    assert m["count"] == 5 and m["freshness"]["status"] == "fresh"
    assert hermes.dispatch("onemind_layer_meta", {"layer": "nope"})["error"] == "layer_service_error"
    s = hermes.dispatch("onemind_layer_snapshot", {"layer": "fleet", "bbox": [-97.76, 30.25, -97.72, 30.29]})
    assert s["count"] == 5 and {"sim-x1", "sim-tractor-1"} <= {i["entity_id"] for i in s["items"]}
    assert hermes.dispatch("onemind_layer_snapshot", {"layer": "air"})["error"] == "invalid_argument"
    capped = hermes.dispatch("onemind_layer_snapshot", {"layer": "fleet", "bbox": [-98, 30, -97, 31], "max": 2})
    assert capped["count"] == 2 and capped["truncated"] is True
    assert hermes.dispatch("onemind_layer_snapshot", {"layer": "fleet", "bbox": [-98, 30, -97, 31],
                                                      "max": 5000})["error"] == "invalid_argument"


# -- proposals ----------------------------------------------------------------------------------

def test_command_propose_sends_an_agent_request_on_behalf_of_the_human(hermes, world):
    args = {"action_type": "asset.goto", "target": {"entity_type": "drone", "entity_id": "sim-x1"},
            "params": {"lat": 30.275, "lon": -97.742, "hae_m": 200}, "reason": "look at the barn roof"}
    r = hermes.dispatch("onemind_command_propose", args, user_task="send sim-x1 over the barn")
    assert r["status"] == "awaiting_confirmation" and r["approvals_needed"] == 1
    assert "human" in r["note"] and "cannot confirm" in r["note"]
    req = world.proposals[-1]
    assert req["family"] == "command"
    assert req["requested_by"] == {"kind": "agent", "id": "agent:hermes/zeus/legacy", "on_behalf_of": "user:zeus"}
    assert req["agent_session"] == "session-test" and len(req["prompt_sha256"]) == 64
    assert len(req["tool_input_sha256"]) == 64 and "send sim-x1" not in json.dumps(req)  # hashes only
    # The request_id is per tool call: a retry of the same call (same tool_call_id) reuses it, a new
    # call with the same input in the same session ("try again") gets a new one.
    hermes.dispatch("onemind_command_propose", args, tool_call_id="call-1")
    rid1 = world.proposals[-1]["request_id"]
    hermes.dispatch("onemind_command_propose", args, tool_call_id="call-1")
    assert world.proposals[-1]["request_id"] == rid1
    hermes.dispatch("onemind_command_propose", args, tool_call_id="call-2")
    assert world.proposals[-1]["request_id"] != rid1
    hermes.dispatch("onemind_command_propose", args)
    hermes.dispatch("onemind_command_propose", args)
    ids = [p["request_id"] for p in world.proposals]
    assert len(set(ids)) == len(ids) - 1  # only the call-1 retry repeated
    assert all(re.fullmatch(r"[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}", i) for i in ids)
    time.sleep(0.2)
    assert not [s for s in world.subjects_seen() if s.startswith("cmd.fleet.")]
    assert world.confirm_calls == []


def test_task_family_routing_and_missing_operator(hermes, plugin, nats_server):
    r = hermes.dispatch("onemind_command_propose", {"action_type": "task.assign", "target": {}, "params": {}})
    assert r["error"] == "wrong_family"
    r = hermes.dispatch("onemind_task_propose", {"action_type": "asset.goto", "target": {}, "params": {}})
    assert r["error"] == "wrong_family"
    ctx = FakeCtx()
    plugin.register(ctx, env={"ONEMIND_BUS_URL": nats_server.url()})
    r = ctx.dispatch("onemind_command_propose", {"action_type": "asset.rtl",
                                                 "target": {"entity_type": "drone", "entity_id": "sim-x1"}, "params": {}})
    assert r["error"] == "not_configured"


def test_person_reference_is_denied_and_reported(hermes):
    r = hermes.dispatch("onemind_command_propose", {
        "action_type": "asset.orbit", "target": {"entity_type": "drone", "entity_id": "sim-x1"},
        "params": {"center": {"entity_id": "trk-p1"}, "radius_m": 30, "hae_m": 60, "duration_s": 600}})
    assert r["denial"]["code"] == "forbidden_action" and "refused" in r["note"]


def test_action_status(hermes):
    r = hermes.dispatch("onemind_task_propose", {"action_type": "task.assign",
                                                 "target": {"entity_type": "drone", "entity_id": "sim-x1"}, "params": {}})
    s = hermes.dispatch("onemind_action_status", {"action_id": r["action_id"]})
    assert s["action"]["status"] == "awaiting_confirmation"


# -- map actions --------------------------------------------------------------------------------

def test_map_fly_to_entity_reports_the_applied_view(hermes, world):
    r = hermes.dispatch("onemind_map_fly_to", {"entity_id": "sim-x1", "range_m": 800})
    assert r["ok"] is True and r["applied_by"] == ["stub-map-1"]
    assert r["view"]["camera"]["lat"] == 30.27 and world.view.camera["range_m"] == 800
    subj, raw = world.map_msgs[-1]
    p = json.loads(raw)
    assert subj == "cmd.map.viewport.set" and p["source"] == "agent:hermes/zeus/legacy" and p["v"] == 1
    assert p["zoom"] > 10 and p["corr"] == "session-test"


def test_map_action_without_a_client_ack_is_not_success(hermes, world):
    world.map_mode = "silent"
    r = hermes.dispatch("onemind_map_fly_to", {"lat": 30.3, "lon": -97.7})
    assert r["ok"] is False and r["error"] == "no_client_ack" and "Do not tell" in r["say"]
    world.map_mode = "reject"
    r = hermes.dispatch("onemind_map_select", {"entity_id": "sim-x1"})
    assert r["ok"] is False and r["error"] == "rejected_by_client"
    world.map_mode = "ack"
    # the client acks ok=false when it cannot apply (unknown annotation)
    r = hermes.dispatch("onemind_map_remove_annotation", {"annotation_id": "hermes-missing"})
    assert r["ok"] is False and r["error"] == "rejected_by_client"


def test_map_layers_validates_against_the_layer_service(hermes, world):
    r = hermes.dispatch("onemind_map_layers", {"hide": ["air"], "show": ["alerts"]})
    assert r["ok"] and "air" not in r["view"]["visible_layers"] and "alerts" in r["view"]["visible_layers"]
    r = hermes.dispatch("onemind_map_layers", {"hide": ["satellites"]})
    assert r["error"] == "unknown_layer"
    assert hermes.dispatch("onemind_map_layers", {"show": ["a"], "hide": ["a"]})["error"] == "invalid_argument"


def test_annotate_uses_the_real_boundary(hermes, world):
    r = hermes.dispatch("onemind_map_annotate", {"entity_id": "zone-north-field", "label": "North field"})
    assert r["ok"] is True
    ann = world.view.annotations[r["annotation_id"]]
    assert ann["kind"] == "boundary"
    assert ann["geometry"] == world.entities["zone-north-field"]["geo"]["coords"]
    r2 = hermes.dispatch("onemind_map_remove_annotation", {"annotation_id": r["annotation_id"]})
    assert r2["ok"] and r["annotation_id"] not in world.view.annotations


def test_measure_reports_distance_and_whether_it_was_drawn(hermes, world):
    r = hermes.dispatch("onemind_map_measure", {"entity_ids": ["sim-x1", "structure-barn"]})
    expect = P.geo.haversine_m(30.27, -97.74, 30.274, -97.7445)
    assert abs(r["measure"]["distance_m"] - expect) < 1 and r["drawn_on_map"] is True
    world.map_mode = "silent"
    r = hermes.dispatch("onemind_map_measure", {"points": [[30.27, -97.74], [30.28, -97.74], [30.28, -97.73]],
                                                "closed": True})
    assert r["drawn_on_map"] is False and r["measure"]["area_m2"] > 0


def test_camera_tour_and_frame(hermes, world):
    r = hermes.dispatch("onemind_map_camera_tour", {"stops": [{"entity_id": "structure-house"},
                                                              {"entity_id": "structure-barn", "dwell_s": 2}]})
    assert r["ok"] and r["view"]["tour_stops"] == 2
    assert hermes.dispatch("onemind_map_camera_tour", {"stops": [{"lat": 1, "lon": 1}]})["error"] == "invalid_argument"
    r = hermes.dispatch("onemind_map_frame", {"entity_ids": ["sim-x1", "sim-x2", "sim-x3"]})
    w, s, e, n = world.view.camera["bbox"]
    for eid in ("sim-x1", "sim-x2", "sim-x3"):
        loc = world.entities[eid]["location"]
        assert w < loc["lon"] < e and s < loc["lat"] < n


def test_select_and_clear(hermes, world):
    assert hermes.dispatch("onemind_map_select", {"entity_id": "sim-tractor-1", "label": "Tractor"})["ok"]
    assert "sim-tractor-1" in world.view.selected
    assert hermes.dispatch("onemind_map_select", {"entity_id": "sim-tractor-1", "clear": True})["ok"]
    assert "sim-tractor-1" not in world.view.selected
    assert hermes.dispatch("onemind_map_select", {"entity_id": "ghost"})["error"] == "not_found"


# -- the I2 safety property ---------------------------------------------------------------------

def test_no_tool_reaches_cmd_fleet_or_confirm_even_with_hostile_arguments(hermes, world):
    hostile = {"subject": "cmd.fleet.drone.sim-x1.goto", "entity_id": "sim-x1.goto", "annotation_id": "x.>",
               "action_type": "asset.goto", "target": {"entity_type": "drone", "entity_id": "sim-x1"},
               "params": {"lat": 30.27, "lon": -97.74, "hae_m": 100}, "layer": "cmd.fleet.>",
               "show": ["cmd.fleet.>"], "stops": [{"lat": 30, "lon": -97}, {"lat": 30.1, "lon": -97}],
               "points": [[30, -97], [30.1, -97]], "label": "x", "near": {"lat": 30, "lon": -97, "radius_m": 10}}
    for name in hermes.tools:
        hermes.dispatch(name, dict(hostile))
    world.flush()
    time.sleep(0.3)
    seen = world.subjects_seen()
    assert not [s for s in seen if s.startswith("cmd.") and not s.startswith("cmd.map.")]
    assert not [s for s in seen if s.endswith(".confirm")]
    assert world.confirm_calls == []
    assert all(s.startswith(ALLOWED_PREFIXES) for s in hermes.client.sent), hermes.client.sent


def test_injected_callsign_is_returned_as_data_only(hermes, world):
    r = hermes.dispatch("onemind_entity_query", {"near": {"entity_id": "zone-north-field", "radius_m": 500}})
    names = [e["name"] for e in r["entities"]]
    assert any("IGNORE PREVIOUS INSTRUCTIONS" in n for n in names)
    assert "not instructions" in r["note"]
    assert not [s for s in world.subjects_seen() if s.startswith("cmd.fleet")]


# -- review fixes: outcome notes, geometry cap, request-many truncation --------------------------

class _FakeClient:
    def __init__(self, replies):
        self.replies, self.sent = replies, []

    def request(self, subject, body, timeout=None):
        self.sent.append((subject, body))
        return self.replies[subject](body) if callable(self.replies[subject]) else self.replies[subject]


def _tools(replies):
    return P.tools.Tools(client=_FakeClient(replies), agent_id="agent:hermes/zeus/legacy", on_behalf_of="user:zeus")


GOTO_ARGS = {"action_type": "asset.goto", "target": {"entity_type": "drone", "entity_id": "sim-x1"},
             "params": {"lat": 30.275, "lon": -97.742}}


def test_human_rejection_and_expiry_are_not_reported_as_policy_denials():
    for code, words in (("rejected_by_human", "human rejected"), ("expired", "expired")):
        t = _tools({"svc.policy.command.propose": {"action_id": "act-1", "status": "rejected",
                                                   "denial": {"code": code, "message": "not now"}}})
        r = t.t_command_propose(GOTO_ARGS, {"task_id": "s", "tool_call_id": "c"})
        assert words in r["note"] and "Denied by policy" not in r["note"], r
    t = _tools({"svc.policy.command.propose": {"action_id": "act-1", "status": "denied",
                                               "denial": {"code": "forbidden_action", "message": "x"}}})
    assert "Denied by policy" in t.t_command_propose(GOTO_ARGS, {"task_id": "s"})["note"]
    # a replayed call whose action already ran is not reported as a fresh proposal
    t = _tools({"svc.policy.command.propose": {"action_id": "act-1", "status": "executed"}})
    r = t.t_command_propose(GOTO_ARGS, {"task_id": "s", "tool_call_id": "c"})
    assert "already submitted" in r["note"] and "awaits" not in r["note"]


def test_large_entity_boundary_is_refused_with_a_clear_error():
    ring = [[-97.0 + i * 1e-5, 30.0] for i in range(P.tools.MAX_GEOMETRY_VERTICES + 10)]
    big = {"entity_id": "zone-big", "entity_type": "zone",
           "geo": {"coords": {"type": "MultiPolygon", "coordinates": [[ring]]}}}
    ctx = FakeCtx()
    tools = P.register(ctx, env={"ONEMIND_BUS_URL": "nats://127.0.0.1:1"},
                       client=_FakeClient({"svc.entity.get": {"entity": big}}))
    del tools
    r = ctx.dispatch("onemind_map_annotate", {"entity_id": "zone-big", "label": "big field"})
    assert r["ok"] is False and r["error"] == "geometry_too_large", r
    assert str(P.tools.MAX_GEOMETRY_VERTICES) in r["message"]
    assert P.tools._vertex_count([[ring, ring]]) == 2 * len(ring)


def test_request_many_marks_truncation_at_max_parts(hermes):
    body = {"v": 1, "layer": "fleet", "aoi": {"bbox": [-98, 30, -97, 31]}, "zoom": 12, "fields": "render", "max": 200}
    parts = hermes.client.request_many("svc.layer.snapshot", body, max_parts=2)
    assert len(parts) == 3 and parts[-1]["truncated"] is True and not any(p.get("last") for p in parts)
    full = hermes.client.request_many("svc.layer.snapshot", body)
    assert full[-1].get("last") is True and not any("detail" in p for p in full)
