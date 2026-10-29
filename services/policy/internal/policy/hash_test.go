package policy

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"
)

// vectorFile is the cross-implementation params_hash vector file shared with the MAVLink
// and ROS 2 bridges (asset-command.md section 5; decision HASH, 2026-10-29).
var vectorFile = filepath.Join(schemaDir, "command", "testdata", "params-hash-vectors.json")

type hashVector struct {
	Name          string `json:"name"`
	ActionType    string `json:"action_type"`
	ActionVersion int    `json:"action_version"`
	Target        Target `json:"target"`
	ParamsJSON    string `json:"params_json"`
	Canonical     string `json:"canonical"`
	ParamsHash    string `json:"params_hash"`
}

func loadVectors(t *testing.T) map[string]hashVector {
	t.Helper()
	b, err := os.ReadFile(vectorFile)
	if err != nil {
		t.Fatal(err)
	}
	var doc struct {
		Vectors []hashVector `json:"vectors"`
	}
	if err := json.Unmarshal(b, &doc); err != nil {
		t.Fatal(err)
	}
	out := map[string]hashVector{}
	for _, v := range doc.Vectors {
		out[v.Name] = v
	}
	if len(out) < 10 {
		t.Fatalf("only %d vectors", len(out))
	}
	return out
}

func TestParamsHashVectors(t *testing.T) {
	for _, v := range loadVectors(t) {
		h, err := ParamsHash(v.ActionType, v.ActionVersion, v.Target, json.RawMessage(v.ParamsJSON))
		if err != nil || h != v.ParamsHash {
			t.Errorf("%s: got %s %v, want %s", v.Name, h, err, v.ParamsHash)
		}
		p := v.ParamsJSON
		if p == "" {
			p = "{}"
		}
		c, err := Canonicalize(map[string]any{"action_type": v.ActionType, "action_version": v.ActionVersion,
			"target": v.Target, "params": json.RawMessage(p)})
		if err != nil || string(c) != v.Canonical {
			t.Errorf("%s: canonical %s %v, want %s", v.Name, c, err, v.Canonical)
		}
	}
}

// A command's approval.params_hash is the approved action's own hash: asset.{verb} over the
// command data for an asset action; the task.assign hash for a task-bearing command.
func TestCommandCarriesActionHash(t *testing.T) {
	e := newEnv(t)
	r := e.propose("zeus", FamilyCommand, gotoReq("sim-x1", 30.27, -97.742))
	if c := e.confirm("zeus", FamilyCommand, r.ActionID, r.ParamsHash, "approve"); c.Status != StateDispatched {
		t.Fatalf("goto confirm: %+v %+v", c, c.Denial)
	}
	waitFor(t, func() bool { return len(e.cmdsFor(r.ActionID)) == 1 }, "goto command")
	cmd := e.cmdsFor(r.ActionID)[0]
	data, _ := json.Marshal(cmd["data"])
	want, _ := ParamsHash("asset.goto", 1, Target{"drone", "sim-x1"}, data)
	if got := cmd["approval"].(map[string]any)["params_hash"]; got != r.ParamsHash || got != want || cmd["task_id"] != nil {
		t.Fatalf("asset command hash %v, action %s, recomputed %s", got, r.ParamsHash, want)
	}

	body := req("task.assign", "drone", "sim-x1", map[string]any{"parent_id": nil, "type": "onemind.task.goto.v1",
		"spec": map[string]any{"lat": 30.27, "lon": -97.74, "hae_m": 150}, "description": "hash binding"})
	tr := e.propose("zeus", FamilyTask, body)
	if c := e.confirm("zeus", FamilyTask, tr.ActionID, tr.ParamsHash, "approve"); c.Status != StateDispatched {
		t.Fatalf("task confirm: %+v %+v", c, c.Denial)
	}
	waitFor(t, func() bool { return len(e.cmdsFor(tr.ActionID)) == 2 }, "task commands")
	fleet := e.cmdsFor(tr.ActionID)[1]
	if fleet["subject"] != "cmd.fleet.drone.sim-x1.goto" || fleet["task_id"] == nil {
		t.Fatalf("fleet command: %v", fleet)
	}
	if got := fleet["approval"].(map[string]any)["params_hash"]; got != tr.ParamsHash {
		t.Fatalf("task-bearing command carries %v, want the task.assign hash %s", got, tr.ParamsHash)
	}
}

// D1-P4: zone.create|update|remove are A2, execute as cmd.map.geofence.{id}.{verb} and close
// on the geofence service's evt.map.geofence.{created|updated|deleted|refused}.
func TestZoneActions(t *testing.T) {
	e := newEnv(t)
	geo := map[string]any{"shape": "circle", "kind": "keep_in",
		"coords": map[string]any{"type": "Point", "coordinates": []float64{-97.74, 30.27}}, "radius_m": 500}
	cases := []struct {
		action, zone, verb, event string
		params                    map[string]any
	}{
		{"zone.create", "north-gate", "create", "executed", map[string]any{"name": "North gate", "geo": geo}},
		{"zone.update", "north-gate", "update", "executed", map[string]any{"name": "North gate 2"}},
		{"zone.remove", "north-gate", "remove", "executed", map[string]any{"reason": "done"}},
		{"zone.create", "bad-zone", "create", StateFailed, map[string]any{"geo": geo}},
	}
	for _, c := range cases {
		r := e.propose("zeus", FamilyCommand, req(c.action, "zone", c.zone, c.params))
		if r.Status != StateAwaiting || r.Level != A2 || r.ApprovalsNeeded != 1 {
			t.Fatalf("%s propose: %+v %+v", c.action, r, r.Denial)
		}
		cf := e.confirm("zeus", FamilyCommand, r.ActionID, r.ParamsHash, "approve")
		if cf.Status != StateDispatched {
			t.Fatalf("%s confirm: %+v %+v", c.action, cf, cf.Denial)
		}
		waitFor(t, func() bool { return e.hasAudit(r.ActionID, "dispatched", "") }, c.action+" dispatched")
		waitFor(t, func() bool { return e.hasAudit(r.ActionID, c.event, "") }, c.action+" "+c.event)
		cmds := e.cmdsFor(r.ActionID)
		if len(cmds) != 1 || cmds[0]["subject"] != "cmd.map.geofence."+c.zone+"."+c.verb {
			t.Fatalf("%s commands: %v", c.action, cmds)
		}
		if c.event == StateFailed && !e.hasAudit(r.ActionID, StateFailed, "schema_invalid") {
			t.Fatalf("refusal code not carried: %v", e.auditFor(r.ActionID))
		}
	}
	// A zone action never targets an asset, and an unknown zone verb is not in the catalogue.
	expectDenied(t, e.propose("zeus", FamilyCommand, req("zone.create", "drone", "sim-x1", map[string]any{})), CodeSchemaInvalid)
	expectDenied(t, e.propose("zeus", FamilyCommand, req("zone.delete", "zone", "north-gate", map[string]any{})), CodeUnknownAction)
	// A result for another zone does not close the action.
	r := e.propose("zeus", FamilyCommand, req("zone.update", "zone", "south", map[string]any{"name": "S"}))
	e.geofenceDown(true)
	_ = e.confirm("zeus", FamilyCommand, r.ActionID, r.ParamsHash, "approve")
	waitFor(t, func() bool { return len(e.cmdsFor(r.ActionID)) == 1 }, "zone command")
	b, _ := json.Marshal(map[string]any{"v": 1, "id": NewUUID(), "causation": e.cmdsFor(r.ActionID)[0]["id"],
		"data": map[string]any{"zone_id": "other", "action_id": r.ActionID, "outcome": "executed"}})
	_ = e.conn("geofence").Publish("evt.map.geofence.updated", b)
	time.Sleep(200 * time.Millisecond)
	if e.hasAudit(r.ActionID, "executed", "") {
		t.Fatal("a result naming another zone closed the action")
	}
	e.geofenceDown(false)
}
