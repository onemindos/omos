"""Scripted-model runs (no LLM, no paid API): FEATURES D11 and I2 acceptance shapes.

A scripted model stands in for the LLM. For each text request it calls the plugin's tools the way
a model would (through the Hermes-style dispatch, so the pre_tool_call hook applies), and writes
its answer only from what the tools returned. Each request is checked against the view state of
the stub map client, which applied the ``cmd.map.*`` action and acknowledged it.

The real-model run and the screenshot check need a map client that implements ``cmd.map.*``
and the owner's Hermes; they are in the evidence folder's OWNER-STEPS.md.

Set ONEMIND_BUS_EVIDENCE_DIR to write the transcript (JSONL and Markdown).
"""
from __future__ import annotations

import copy
import json
import os
from pathlib import Path

from conftest import load_plugin

P = load_plugin()


class ScriptedModel:
    def __init__(self, ctx, world):
        self.ctx, self.world = ctx, world
        self.calls = []

    def call(self, tool, args):
        r = self.ctx.dispatch(tool, args, task_id="scripted-session", user_task=self.request)
        self.calls.append({"tool": tool, "args": args, "result": r})
        return r

    def find(self, name, types=None):
        args = {"name": name}
        if types:
            args["types"] = types
        rows = self.call("onemind_entity_query", args)["entities"]
        return rows[0]["entity_id"] if rows else None

    def done(self, r, what):
        """Claim an action only if the tool says the map applied it."""
        if r.get("ok"):
            return f"Done: {what}."
        return f"I could not {what}: the map did not confirm it ({r.get('error')})."

    def run(self, request, plan):
        self.request, self.calls = request, []
        answer = plan(self)
        return {"request": request, "tool_calls": self.calls, "answer": answer}


def view_state(world):
    v = world.view
    return {"camera": copy.deepcopy(v.camera), "visible_layers": sorted(v.visible), "selected": sorted(v.selected),
            "annotations": copy.deepcopy(v.annotations), "tour_stops": len(v.tour) if v.tour else 0}


def pos(world, eid):
    loc = world.entities[eid].get("location")
    return loc["lat"], loc["lon"]


# Ten D11 requests: (text, plan, check(world, before, after, run) -> None)
def _r1(m):
    return m.done(m.call("onemind_map_fly_to", {"entity_id": "sim-x1", "range_m": 600}), "fly to SIM-X1")


def _c1(w, b, a, run):
    assert (round(a["camera"]["lat"], 4), round(a["camera"]["lon"], 4)) == pos(w, "sim-x1")
    assert a["camera"]["range_m"] == 600


def _r2(m):
    ids = [e["entity_id"] for e in m.call("onemind_entity_query", {"types": ["drone"]})["entities"]]
    return m.done(m.call("onemind_map_frame", {"entity_ids": ids}), f"frame {len(ids)} drones")


def _c2(w, b, a, rr):
    west, south, east, north = a["camera"]["bbox"]
    for eid in ("sim-x1", "sim-x2", "sim-x3"):
        lat, lon = pos(w, eid)
        assert west < lon < east and south < lat < north


def _r3(m):
    return m.done(m.call("onemind_map_layers", {"hide": ["air"]}), "hide the aircraft layer")


def _c3(w, b, a, rr):
    assert "air" in b["visible_layers"] and "air" not in a["visible_layers"]


def _r4(m):
    return m.done(m.call("onemind_map_layers", {"show": ["alerts"]}), "show the alerts layer")


def _c4(w, b, a, rr):
    assert "alerts" not in b["visible_layers"] and "alerts" in a["visible_layers"]


def _r5(m):
    eid = m.find("tractor")
    return m.done(m.call("onemind_map_select", {"entity_id": eid, "label": "Tractor"}), f"select {eid}")


def _c5(w, b, a, rr):
    assert a["selected"] == ["sim-tractor-1"]


def _r6(m):
    eid = m.find("north field", ["zone"])
    r = m.call("onemind_map_annotate", {"entity_id": eid, "label": "North field"})
    m.last_annotation = r.get("annotation_id")
    return m.done(r, "outline the North Field boundary")


def _c6(w, b, a, rr):
    new = set(a["annotations"]) - set(b["annotations"])
    assert len(new) == 1
    ann = a["annotations"][new.pop()]
    assert ann["kind"] == "boundary" and ann["geometry"] == w.entities["zone-north-field"]["geo"]["coords"]


def _r7(m):
    barn = m.find("barn")
    r = m.call("onemind_map_measure", {"entity_ids": ["sim-x1", barn]})
    d = r["measure"]["distance_m"]
    shown = "and drew it on the map" if r.get("drawn_on_map") else "but the map did not show the line"
    return f"SIM-X1 is {d:.0f} m from the barn ({shown})."


def _c7(w, b, a, rr):
    new = [a["annotations"][k] for k in set(a["annotations"]) - set(b["annotations"])]
    assert len(new) == 1 and new[0]["kind"] == "measure"
    expect = P.geo.haversine_m(*pos(w, "sim-x1"), *pos(w, "structure-barn"))
    assert abs(new[0]["measure"]["distance_m"] - expect) < 1
    assert f"{expect:.0f} m" in rr["answer"]


def _r8(m):
    stops = [{"entity_id": m.find(n), "dwell_s": 3} for n in ("house", "barn", "pond")]
    return m.done(m.call("onemind_map_camera_tour", {"stops": stops}), "run a camera tour of the house, barn and pond")


def _c8(w, b, a, rr):
    assert a["tour_stops"] == 3
    pond = w.entities["zone-pond"]["geo"]["coords"]["coordinates"][0][:-1]
    plat = sum(p[1] for p in pond) / len(pond)
    assert abs(a["camera"]["lat"] - plat) < 1e-5


def _r9(m):
    gate = m.find("gate")
    return m.done(m.call("onemind_map_annotate", {"entity_id": gate, "label": "Check lock", "note": "check lock tonight"}),
                  "mark the front gate")


def _c9(w, b, a, rr):
    new = [a["annotations"][k] for k in set(a["annotations"]) - set(b["annotations"])]
    lat, lon = pos(w, "place-gate")
    assert new and new[0]["geometry"] == {"type": "Point", "coordinates": [lon, lat]} and new[0]["label"] == "Check lock"


def _r10(m):
    return m.done(m.call("onemind_map_remove_annotation", {"annotation_id": m.last_annotation}),
                  "remove the North Field outline")


def _c10(w, b, a, rr):
    gone = set(b["annotations"]) - set(a["annotations"])
    assert len(gone) == 1 and b["annotations"][gone.pop()]["kind"] == "boundary"


D11 = [
    ("Fly to SIM-X1.", _r1, _c1),
    ("Frame all the drones.", _r2, _c2),
    ("Hide the aircraft layer.", _r3, _c3),
    ("Show the alerts layer.", _r4, _c4),
    ("Select the tractor.", _r5, _c5),
    ("Outline the north field.", _r6, _c6),
    ("How far is SIM-X1 from the barn? Show it on the map.", _r7, _c7),
    ("Give me a camera tour of the house, the barn and the pond.", _r8, _c8),
    ("Put a note on the front gate: check lock.", _r9, _c9),
    ("Remove the north field outline.", _r10, _c10),
]


def _write_evidence(name, runs):
    d = os.environ.get("ONEMIND_BUS_EVIDENCE_DIR")
    if not d:
        return
    out = Path(d)
    out.mkdir(parents=True, exist_ok=True)
    with open(out / f"{name}.jsonl", "w") as f:
        for r in runs:
            f.write(json.dumps(r, default=str) + "\n")
    lines = [f"# {name}", ""]
    for i, r in enumerate(runs, 1):
        lines += [f"## {i}. {r['request']}", ""]
        for c in r["tool_calls"]:
            res = c["result"]
            brief = {k: res.get(k) for k in ("ok", "error", "subject", "applied_by", "view", "measure", "count",
                                             "status", "action_id", "denial") if k in res}
            if "entities" in res:
                brief["entities"] = [e["entity_id"] for e in res["entities"]]
            lines.append(f"- `{c['tool']}` {json.dumps(c['args'])} -> `{json.dumps(brief, default=str)}`")
        lines += ["", f"Answer: {r['answer']}", ""]
        if "check" in r:
            lines += [f"View check: **{r['check']}**", ""]
    (out / f"{name}.md").write_text("\n".join(lines))


def test_d11_ten_scripted_requests_each_change_the_view(hermes, world):
    m = ScriptedModel(hermes, world)
    runs = []
    for text, plan, check in D11:
        before = view_state(world)
        run = m.run(text, plan)
        after = view_state(world)
        run.update(view_before=before, view_after=after)
        try:
            check(world, before, after, run)
            run["check"] = "pass"
        except AssertionError as e:
            run["check"] = f"FAIL {e}"
        assert not run["answer"].startswith("I could not"), run
        runs.append(run)
    _write_evidence("d11-scripted-run", runs)
    failed = [r["request"] for r in runs if r["check"] != "pass"]
    assert not failed, failed
    assert len(runs) == 10
    tools_used = {c["tool"] for r in runs for c in r["tool_calls"]}
    assert all(t.startswith("onemind_") for t in tools_used)


def test_d11_without_a_map_client_ack_the_model_reports_failure(hermes, world):
    world.map_mode = "silent"
    m = ScriptedModel(hermes, world)
    runs = [m.run(t, p) for t, p, _ in D11[:5]]
    assert all(r["answer"].startswith("I could not") for r in runs), [r["answer"] for r in runs]
    assert world.view.applied == []
    _write_evidence("d11-no-ack-run", runs)


def test_i2_spatial_question_answered_with_bus_tools_only(hermes, world):
    """'What's moving near the north field?' (FEATURES D4/I2) from svc.entity.* only."""
    m = ScriptedModel(hermes, world)

    def plan(mm):
        field = mm.find("north field", ["zone"])
        r = mm.call("onemind_entity_query", {"near": {"entity_id": field, "radius_m": 500}, "moving_only": True,
                                             "exclude_ids": [field]})
        ids = [e["entity_id"] for e in r["entities"]]
        if ids:
            mm.call("onemind_map_frame", {"entity_ids": ids + [field]})
            for e in r["entities"]:
                if e["entity_type"] in ("drone", "vehicle"):
                    mm.call("onemind_map_select", {"entity_id": e["entity_id"]})
        parts = [f"{e['name']} ({e['entity_id']}, {e['entity_type']}, {e['speed_mps']} m/s, "
                 f"{e['distance_m']} m away, source {','.join(e.get('sources', []))}, observed_t {e.get('observed_t')})"
                 for e in r["entities"]]
        return "Moving near the North Field: " + "; ".join(parts) if parts else "Nothing is moving near the North Field."

    run = m.run("What's moving near the north field?", plan)
    run["check"] = "pass"
    assert "sim-tractor-1" in run["answer"] and "trk-p1" in run["answer"]
    assert "sim-x2" not in run["answer"] and "icao-a1b2c3" not in run["answer"]
    assert {c["tool"] for c in run["tool_calls"]} <= {"onemind_entity_query", "onemind_map_frame", "onemind_map_select"}
    # every subject the plugin sent is a bus tool subject; nothing reached cmd.fleet
    assert all(s.startswith(("svc.entity.", "cmd.map.")) for s in hermes.client.sent)
    assert not [s for s in world.subjects_seen() if s.startswith("cmd.fleet")]
    # the person-class track is answered as a track, not selected or followed
    assert "trk-p1" not in world.view.selected
    _write_evidence("i2-spatial-question", [run])


def test_d5_scripted_refusal_follow_that_person(hermes, world):
    """Asked to follow a person, the script refuses; the one proposal it tries is denied and recorded."""
    m = ScriptedModel(hermes, world)

    def plan(mm):
        r = mm.call("onemind_command_propose", {
            "action_type": "asset.orbit", "target": {"entity_type": "drone", "entity_id": "sim-x1"},
            "params": {"center": {"entity_id": "trk-p1"}, "radius_m": 30, "hae_m": 60, "duration_s": 600},
            "reason": "operator asked to follow that person"})
        if r.get("denial"):
            return f"I won't do that: following a person is refused ({r['denial']['code']}), and the attempt is logged."
        return "Proposed; waiting for a human."

    run = m.run("Follow that person with SIM-X1.", plan)
    assert "refused (forbidden_action)" in run["answer"]
    assert not [s for s in world.subjects_seen() if s.startswith("cmd.fleet")]
    _write_evidence("d5-follow-person-stub", [run])
