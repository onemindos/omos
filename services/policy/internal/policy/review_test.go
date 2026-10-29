package policy

import (
	"encoding/json"
	"testing"
	"time"
)

// A task.assign confirm waiting on a silent task service must not delay an estop on the
// same asset: no lock is shared between actions.
func TestStopNotBlockedByTaskWait(t *testing.T) {
	e := newEnv(t)
	zc := e.conn("zeus")
	e.taskMu.Lock()
	e.taskDown = true
	e.taskMu.Unlock()

	body := req("task.assign", "drone", "sim-x1", map[string]any{"parent_id": nil, "type": "onemind.task.goto.v1",
		"spec": map[string]any{"lat": 30.27, "lon": -97.74, "hae_m": 150}})
	r := e.propose("zeus", FamilyTask, body)
	if r.Status != StateAwaiting {
		t.Fatalf("task propose: %+v %+v", r, r.Denial)
	}
	done := make(chan ConfirmReply, 1)
	go func() {
		b, _ := json.Marshal(map[string]any{"v": 1, "action_id": r.ActionID, "params_hash": r.ParamsHash, "decision": "approve"})
		var c ConfirmReply
		if m, err := zc.Request(SubjTaskConfirm, b, 3*time.Second); err == nil {
			_ = json.Unmarshal(m.Data, &c)
		}
		done <- c
	}()
	// The confirm has published cmd.mission.{tsk}.create and now waits for .created.
	waitFor(t, func() bool { return len(e.cmdsFor(r.ActionID)) == 1 }, "task create published")

	start := time.Now()
	sr := e.propose("zeus", FamilyCommand, req("asset.estop", "drone", "sim-x1", map[string]any{}))
	if sr.Status != StateAwaiting {
		t.Fatalf("estop propose: %+v %+v", sr, sr.Denial)
	}
	sc := e.confirm("zeus", FamilyCommand, sr.ActionID, sr.ParamsHash, "approve")
	elapsed := time.Since(start)
	if sc.Status != StateDispatched || sc.CommandID == "" {
		t.Fatalf("estop confirm: %+v %+v", sc, sc.Denial)
	}
	select {
	case c := <-done:
		t.Fatalf("the task confirm finished (%s) before the estop; nothing was tested concurrently", c.Status)
	default:
	}
	t.Logf("estop proposed and confirmed in %s while the task confirm waited (TaskCreateTimeout %s)",
		elapsed, e.svc.cfg.TaskCreateTimeout)
	if elapsed >= e.svc.cfg.TaskCreateTimeout {
		t.Fatalf("estop took %s: it waited for the task service", elapsed)
	}
	if c := <-done; c.Status != StateFailed {
		t.Fatalf("task confirm with a silent task service: %+v %+v", c, c.Denial)
	}
	waitFor(t, func() bool { return e.hasAudit(sr.ActionID, "executed", "") }, "estop executed")
}

func (e *env) result(subject, causation string, data map[string]any) {
	e.t.Helper()
	env := map[string]any{"v": 1, "id": NewUUID(), "ts": ts(e.clock.Now()), "type": subject, "source": "service-stub",
		"subject": subject, "data": data}
	if causation != "" {
		env["causation"] = causation
	}
	b, _ := json.Marshal(env)
	if err := e.admin.Publish(subject, b); err != nil {
		e.t.Fatal(err)
	}
}

func (e *env) held(actionID string) bool {
	for _, a := range e.svc.Snapshot() {
		if a.ActionID == actionID {
			return true
		}
	}
	return false
}

// Every action reaches a terminal state (by its result or a deadline) and is dropped after
// Retain; denials are never held in memory.
func TestActionsCloseAndAreDropped(t *testing.T) {
	e := newEnv(t)

	// A1 entity.star: closed by the entity service result (causation = command id).
	r1 := e.propose("zeus", FamilyCommand, req("entity.star", "drone", "sim-x1", map[string]any{"starred": true}))
	if r1.Status != StateExecuted || !waitCmd(e, r1.ActionID) {
		t.Fatalf("A1: %+v %+v", r1, r1.Denial)
	}
	e.result("evt.entity.sim-x1.overridden", e.cmdsFor(r1.ActionID)[0]["id"].(string), map[string]any{})
	waitFor(t, func() bool { return e.hasAudit(r1.ActionID, "executed", "") }, "A1 closed by the entity result")

	// A1 with no result: fails after ResultTimeout.
	r2 := e.propose("zeus", FamilyCommand, req("entity.star", "drone", "sim-x2", map[string]any{"starred": true}))
	if !waitCmd(e, r2.ActionID) {
		t.Fatal("A1 without result: no command")
	}

	// Mission actions: closed by evt.mission results, executed and failed.
	m1 := e.propose("zeus", FamilyTask, req("mission.complete", "mission", "msn-1", map[string]any{}))
	if c := e.confirm("zeus", FamilyTask, m1.ActionID, m1.ParamsHash, "approve"); c.Status != StateDispatched {
		t.Fatalf("mission confirm: %+v %+v", c, c.Denial)
	}
	e.result("evt.mission.msn-1.completed", e.cmdsFor(m1.ActionID)[0]["id"].(string), map[string]any{})
	waitFor(t, func() bool { return e.hasAudit(m1.ActionID, "executed", "") }, "mission closed by its result")
	m2 := e.propose("zeus", FamilyTask, req("mission.complete", "mission", "msn-2", map[string]any{}))
	if c := e.confirm("zeus", FamilyTask, m2.ActionID, m2.ParamsHash, "approve"); c.Status != StateDispatched {
		t.Fatalf("mission confirm: %+v %+v", c, c.Denial)
	}
	e.result("evt.mission.msn-2.rejected", "", map[string]any{"action_id": m2.ActionID, "message": "mission is not active"})
	waitFor(t, func() bool { return e.hasAudit(m2.ActionID, "failed", CodePreconditionFailed) }, "mission failed by its result")
	// A result naming another entity does not close an action.
	m3 := e.propose("zeus", FamilyTask, req("mission.complete", "mission", "msn-3", map[string]any{}))
	_ = e.confirm("zeus", FamilyTask, m3.ActionID, m3.ParamsHash, "approve")
	e.result("evt.mission.msn-9.completed", e.cmdsFor(m3.ActionID)[0]["id"].(string), map[string]any{})

	// Task to a person: closed by evt.mission.{tsk}.created from the task service.
	pt := e.propose("zeus", FamilyTask, req("task.assign", "person", "zeus", map[string]any{"parent_id": nil,
		"type": "onemind.task.manual.v1", "spec": map[string]any{"instructions": "check the gate"}}))
	if c := e.confirm("zeus", FamilyTask, pt.ActionID, pt.ParamsHash, "approve"); c.Status != StateDispatched {
		t.Fatalf("person task confirm: %+v %+v", c, c.Denial)
	}
	waitFor(t, func() bool { return e.hasAudit(pt.ActionID, "executed", "") }, "person task closed by created")

	// Asset command acknowledged but never executed: fails after ExecuteTimeout.
	e.bridge.mu.Lock()
	e.bridge.silent["sim-x2"] = true
	e.bridge.mu.Unlock()
	g := e.propose("zeus", FamilyCommand, gotoReq("sim-x2", 30.27, -97.74))
	gc := e.confirm("zeus", FamilyCommand, g.ActionID, g.ParamsHash, "approve")
	if gc.Status != StateDispatched {
		t.Fatalf("goto confirm: %+v %+v", gc, gc.Denial)
	}
	e.result("evt.drone.sim-x2.command", "", map[string]any{"action_id": g.ActionID, "command_id": gc.CommandID,
		"outcome": "acknowledged", "simulated": true})
	waitFor(t, func() bool { return e.hasAudit(g.ActionID, "acknowledged", "") }, "acknowledged")

	// Denials are audited but not held.
	d := e.propose("zeus", FamilyCommand, req("asset.dance", "drone", "sim-x1", map[string]any{}))
	expectDenied(t, d, CodeUnknownAction)
	waitFor(t, func() bool { return e.hasAudit(d.ActionID, "denied", CodeUnknownAction) }, "denial audited")
	if e.held(d.ActionID) {
		t.Fatal("denied proposal kept in memory")
	}

	e.clock.Advance(e.svc.cfg.ResultTimeout + time.Second)
	waitFor(t, func() bool { return e.hasAudit(r2.ActionID, "failed", CodeNoAck) }, "A1 without result failed")
	waitFor(t, func() bool { return e.hasAudit(m3.ActionID, "failed", CodeNoAck) }, "mission with a mismatched result failed")
	if e.hasAudit(m3.ActionID, "executed", "") {
		t.Fatal("a result for another mission closed the action")
	}
	e.clock.Advance(e.svc.cfg.ExecuteTimeout)
	waitFor(t, func() bool { return e.hasAudit(g.ActionID, "failed", CodeNoAck) }, "acknowledged-only command failed")

	if n := len(e.svc.Snapshot()); n == 0 {
		t.Fatal("finished actions dropped before Retain")
	}
	e.clock.Advance(e.svc.cfg.Retain + time.Minute)
	waitFor(t, func() bool { return len(e.svc.Snapshot()) == 0 }, "every action dropped after Retain")
	e.svc.mu.Lock()
	rq, res := len(e.svc.byRQ), len(e.svc.byRes)
	e.svc.mu.Unlock()
	if rq != 0 || res != 0 {
		t.Fatalf("indexes not emptied: byRQ %d byRes %d", rq, res)
	}
}

// An unreadable keep_out zone denies movement for the assets it names; stops still go out.
func TestUnevaluableGeofenceFailsClosed(t *testing.T) {
	e := newEnv(t)
	js, _ := e.admin.JetStream()
	kv, err := js.KeyValue("geofences")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := kv.Put("zone-corridor", []byte(`{"entity_id":"zone-corridor","entity_type":"zone","live":true,
	  "geo":{"shape":"corridor","kind":"keep_out","coords":{"type":"LineString","coordinates":[[-97.75,30.26],[-97.73,30.28]]},"width_m":50},
	  "properties":{"applies_to":["sim-x1"]}}`)); err != nil {
		t.Fatal(err)
	}
	var r ProposeReply
	waitFor(t, func() bool {
		r = e.propose("zeus", FamilyCommand, gotoReq("sim-x1", 30.27, -97.742))
		return r.Denial != nil && r.Denial.Detail == DetailUnevaluable
	}, "goto denied by the unevaluable zone")
	expectDenied(t, r, CodePreconditionFailed)
	task := req("task.assign", "drone", "sim-x1", map[string]any{"parent_id": nil, "type": "onemind.task.goto.v1",
		"spec": map[string]any{"lat": 30.27, "lon": -97.74, "hae_m": 150}})
	if tr := e.propose("zeus", FamilyTask, task); tr.Denial == nil || tr.Denial.Detail != DetailUnevaluable {
		t.Fatalf("goto task past an unevaluable zone: %+v %+v", tr, tr.Denial)
	}
	for _, at := range []string{"asset.estop", "asset.rtl", "asset.hold"} {
		p := map[string]any{}
		if at == "asset.hold" {
			p = map[string]any{"duration_s": 30}
		}
		if s := e.propose("zeus", FamilyCommand, req(at, "drone", "sim-x1", p)); s.Status != StateAwaiting {
			t.Fatalf("%s blocked by an unevaluable zone: %+v %+v", at, s, s.Denial)
		}
	}
	if o := e.propose("zeus", FamilyCommand, gotoReq("sim-x2", 30.27, -97.742)); o.Status != StateAwaiting {
		t.Fatalf("zone applied to an asset it does not name: %+v %+v", o, o.Denial)
	}
	if err := kv.Delete("zone-corridor"); err != nil {
		t.Fatal(err)
	}
	waitFor(t, func() bool {
		return e.propose("zeus", FamilyCommand, gotoReq("sim-x1", 30.27, -97.742)).Status == StateAwaiting
	}, "goto allowed once the zone is removed")
}
