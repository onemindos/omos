package mission

import (
	"context"
	"encoding/json"
	"strings"
	"testing"
	"time"
)

var simX1 = Ref{EntityType: "drone", EntityID: "sim-x1"}

var gotoSpec = map[string]any{"lat": 30.2690, "lon": -97.7420, "hae_m": 180}

// newMission creates a mission through cmd.mission.{msn}.create and waits for the fact.
func newMission(e *env) string {
	e.t.Helper()
	msn := newID("msn")
	e.command(msn, "create", map[string]any{"name": "North gate overwatch", "intent": "watch the gate"})
	e.waitFact(msn, "created")
	return msn
}

// newSentTask creates a goto task for sim-x1 under msn and moves it to sent as the policy
// service does: create, wait for created, dispatch audit with the fleet command id.
func newSentTask(e *env, msn string) (tsk, fleetCmd string) {
	e.t.Helper()
	tsk = newID("tsk")
	e.command(tsk, "create", taskData(tsk, msn, "onemind.task.goto.v1", gotoSpec, simX1))
	e.waitFact(tsk, "created")
	fleetCmd = uuid4()
	e.dispatched(tsk, "task.assign", fleetCmd, simX1)
	e.waitState(tsk, StateSent)
	return tsk, fleetCmd
}

// Contract test 1 (H4), task-service half: a goto task for sim-x1 goes created, sent,
// received, executing, done_ok; every step is a fact on evt.mission.{tsk}.* with
// status_version = version.status; the mission follows (planned, active, completed).
// The full flow with the real policy service and bridge is test/e2e.sh.
func TestContract1Lifecycle(t *testing.T) {
	e := newEnv(t)
	msn := newMission(e)
	m := e.mission(msn)
	if m.State != MissionPlanned || m.Rollup != "requested" || m.CreatedBy != "service:policy" {
		t.Fatalf("mission record %+v", m)
	}

	tsk := newID("tsk")
	createCmd := e.command(tsk, "create", taskData(tsk, msn, "onemind.task.goto.v1", gotoSpec, simX1))
	created := e.waitFact(tsk, "created")
	// POL-14: the created fact carries causation = the create command id.
	if created["causation"] != createCmd {
		t.Fatalf("created causation %v, want %s", created["causation"], createCmd)
	}
	rec := e.task(tsk)
	if rec.Status.State != StateCreated || rec.Version.Status != 1 || rec.Rollup != "requested" || !rec.Simulated ||
		rec.Delivery.BeforeT != e.clock.Now().UnixMilli()+60000 || !rec.Delivery.RequireAck || rec.Author.Kind != "agent" {
		t.Fatalf("task record %+v", rec)
	}
	e.waitFact(msn, "updated")
	if m := e.mission(msn); len(m.TaskIDs) != 1 || m.TaskIDs[0] != tsk || len(m.Assignments) != 1 || m.Assignments[0].To != "drone.sim-x1" || !m.Simulated {
		t.Fatalf("mission not linked: %+v", m)
	}

	fleet := uuid4()
	e.dispatched(tsk, "task.assign", fleet, simX1)
	e.waitState(tsk, StateSent)
	e.report(tsk, StateReceived, 1, fleet, nil)
	e.waitState(tsk, StateReceived)
	e.report(tsk, StateExecuting, 2, fleet, nil)
	rec = e.waitState(tsk, StateExecuting)
	if rec.Status.StartT == nil || rec.Rollup != "running" || rec.UpdatedBy != "bridge-mavlink" {
		t.Fatalf("executing record %+v", rec)
	}
	e.settle()
	if m := e.mission(msn); m.State != MissionActive || m.Rollup != "running" {
		t.Fatalf("mission should be active: %+v", m)
	}
	e.report(tsk, StateDoneOK, 3, fleet, nil)
	rec = e.waitState(tsk, StateDoneOK)
	if rec.Version.Status != 5 || rec.Rollup != "completed" || rec.Status.EndT == nil {
		t.Fatalf("done record %+v", rec)
	}
	e.waitFact(msn, "completed")
	if m := e.mission(msn); m.State != MissionCompleted || m.Rollup != "completed" {
		t.Fatalf("mission should be completed: %+v", m)
	}
	if got := strings.Join(e.subjects(tsk), ","); got != "created,sent,received,executing,done_ok" {
		t.Fatalf("task facts %s", got)
	}
	for i, st := range []string{"created", "sent", "received", "executing", "done_ok"} {
		f := e.waitFact(tsk, st)
		d := f["data"].(map[string]any)
		ver := d["status_version"]
		if st == "created" {
			ver = d["version"].(map[string]any)["status"]
		}
		if ver != float64(i+1) {
			t.Errorf("%s status_version %v, want %d", st, ver, i+1)
		}
	}
	e.validateFacts()
	if n := len(e.cmds.all()); n != 2 { // mission create and task create (the test's); none from the service
		t.Fatalf("cmd.> carried %d messages; the task service must publish none", n)
	}
}

// Contract test 2: an out-of-order report (status_version 3 after 5) is dropped; a report
// after done_ok is dropped.
func TestContract2OutOfOrderAndTerminal(t *testing.T) {
	e := newEnv(t)
	msn := newMission(e)
	tsk, fleet := newSentTask(e, msn)
	e.report(tsk, StateReceived, 1, fleet, nil)
	e.report(tsk, StateExecuting, 5, fleet, nil)
	e.waitState(tsk, StateExecuting)
	before := e.task(tsk)
	e.report(tsk, StateWaiting, 3, fleet, nil) // older than 5: dropped
	e.settle()
	after := e.task(tsk)
	if after.Status.State != StateExecuting || after.Version.Status != before.Version.Status || e.svc.C.ReportsStale.Load() != 1 {
		t.Fatalf("stale report applied: %+v (stale=%d)", after.Status, e.svc.C.ReportsStale.Load())
	}
	e.report(tsk, StateDoneOK, 6, fleet, nil)
	done := e.waitState(tsk, StateDoneOK)
	e.report(tsk, StateExecuting, 7, fleet, nil) // after done_ok: dropped (rule 2)
	e.report(tsk, StateDoneNotOK, 8, fleet, map[string]string{"code": "failed", "message": "late"})
	e.settle()
	final := e.task(tsk)
	if final.Status.State != StateDoneOK || final.Version.Status != done.Version.Status || e.svc.C.ReportsTerminal.Load() != 2 {
		t.Fatalf("terminal task changed: %+v (terminal drops=%d)", final.Status, e.svc.C.ReportsTerminal.Load())
	}
	if strings.Count(strings.Join(e.subjects(tsk), ","), "waiting") != 0 {
		t.Fatal("dropped report produced a fact")
	}
	e.validateFacts()
}

// The CAS design tolerates a restarted bridge (omos #550 P7): its status_version starts
// again at 1 for a new command of the same task, and that report is accepted, while a
// late report from the older command is still dropped. A report that beats the dispatch
// audit is accepted too.
func TestBridgeRestartAndEarlyReport(t *testing.T) {
	e := newEnv(t)
	msn := newMission(e)
	tsk := newID("tsk")
	e.command(tsk, "create", taskData(tsk, msn, "onemind.task.goto.v1", gotoSpec, simX1))
	e.waitFact(tsk, "created")
	fleet := uuid4()
	// The bridge acknowledges before the policy service audits the dispatch.
	e.report(tsk, StateReceived, 1, fleet, nil)
	e.waitState(tsk, StateReceived)
	e.dispatched(tsk, "task.assign", fleet, simX1) // late audit: no sent fact, command known
	e.report(tsk, StateExecuting, 2, fleet, nil)
	e.waitState(tsk, StateExecuting)
	e.settle()
	if s := e.subjects(tsk); strings.Contains(strings.Join(s, ","), "sent") {
		t.Fatalf("sent after received: %v", s)
	}
	// Abort: the policy service sends cmd.mission.{tsk}.abort, then a hold to a bridge that
	// restarted in between, so the hold's reports start at status_version 1 again.
	e.command(tsk, "abort", map[string]any{"task_id": tsk, "reason": "test"})
	e.waitState(tsk, StateCancelRequested)
	hold := uuid4()
	e.dispatched(tsk, "task.abort", hold, simX1)
	e.report(tsk, StateWaiting, 3, fleet, nil) // pre-abort command progress: state kept
	e.report(tsk, StateReceived, 1, hold, nil) // restarted bridge, new command, version 1
	e.settle()
	if st := e.task(tsk).Status.State; st != StateCancelRequested {
		t.Fatalf("state %s, want cancel_requested kept", st)
	}
	e.report(tsk, StateExecuting, 2, hold, nil) // the hold executed: cancelled
	rec := e.waitState(tsk, StateDoneNotOK)
	if rec.Status.Error == nil || rec.Status.Error.Code != ErrCancelled || rec.Rollup != "cancelled" {
		t.Fatalf("abort result %+v", rec.Status)
	}
	if len(rec.Guard.Commands) != 2 || rec.Guard.Epoch != 1 || rec.Guard.LastVersion != 2 {
		t.Fatalf("guard %+v", rec.Guard)
	}
	e.validateFacts()
}

// Within one task, a report from an older command after a newer command's report is stale.
func TestOlderCommandReportDropped(t *testing.T) {
	e := newEnv(t)
	msn := newMission(e)
	tsk, fleet := newSentTask(e, msn)
	e.report(tsk, StateExecuting, 1, fleet, nil)
	e.waitState(tsk, StateExecuting)
	e.command(tsk, "abort", map[string]any{"task_id": tsk})
	e.waitState(tsk, StateCancelRequested)
	hold := uuid4()
	e.dispatched(tsk, "task.abort", hold, simX1)
	// The hold is refused by the asset: the task keeps its state and says rejected.
	e.report(tsk, StateDoneNotOK, 2, hold, map[string]string{"code": "rejected", "message": "precondition_failed: link lost"})
	rec := e.waitState(tsk, StateExecuting)
	if rec.Status.Error == nil || rec.Status.Error.Code != ErrRejected {
		t.Fatalf("refused abort: %+v", rec.Status)
	}
	e.report(tsk, StateWaiting, 9, fleet, nil) // older command, higher number: still stale
	e.settle()
	if e.task(tsk).Status.State != StateExecuting || e.svc.C.ReportsStale.Load() != 1 {
		t.Fatalf("older command report applied (stale=%d)", e.svc.C.ReportsStale.Load())
	}
}

// Contract test 6: a task not acknowledged by delivery.before_t ends done_not_ok + timeout
// and is not sent again (the task service never publishes a command).
func TestContract6DeliveryTimeout(t *testing.T) {
	e := newEnv(t)
	msn := newMission(e)
	tsk, fleet := newSentTask(e, msn)
	cmdsBefore := len(e.cmds.all())
	e.clock.Advance(59 * time.Second)
	e.svc.SweepOnce()
	if st := e.task(tsk).Status.State; st != StateSent {
		t.Fatalf("timed out early: %s", st)
	}
	e.clock.Advance(2 * time.Second)
	e.svc.SweepOnce()
	rec := e.waitState(tsk, StateDoneNotOK)
	if rec.Status.Error == nil || rec.Status.Error.Code != ErrTimeout || rec.Rollup != "failed" || rec.UpdatedBy != "service:task" {
		t.Fatalf("timeout record %+v", rec.Status)
	}
	f := e.waitFact(tsk, StateDoneNotOK)
	if d := f["data"].(map[string]any); d["error"].(map[string]any)["code"] != "timeout" {
		t.Fatalf("fact %v", d)
	}
	// A late receipt changes nothing, and nothing is ever re-sent.
	e.report(tsk, StateReceived, 1, fleet, nil)
	e.svc.SweepOnce()
	e.settle()
	time.Sleep(200 * time.Millisecond)
	if st := e.task(tsk).Status.State; st != StateDoneNotOK {
		t.Fatalf("late receipt applied: %s", st)
	}
	if n := len(e.cmds.all()); n != cmdsBefore {
		t.Fatalf("%d commands after the timeout; want none", n-cmdsBefore)
	}
	if strings.Count(strings.Join(e.subjects(tsk), ","), "done_not_ok") != 1 {
		t.Fatalf("facts %v", e.subjects(tsk))
	}
	e.validateFacts()
}

// stale_t ends a running task too.
func TestStaleTimeout(t *testing.T) {
	e := newEnv(t)
	tsk := newID("tsk")
	d := taskData(tsk, "", "onemind.task.goto.v1", gotoSpec, simX1)
	d["stale_t"] = e.clock.Now().Add(10 * time.Minute).UnixMilli()
	e.command(tsk, "create", d)
	e.waitFact(tsk, "created")
	fleet := uuid4()
	e.dispatched(tsk, "task.assign", fleet, simX1)
	e.report(tsk, StateExecuting, 1, fleet, nil)
	e.waitState(tsk, StateExecuting)
	e.clock.Advance(11 * time.Minute)
	e.svc.SweepOnce()
	if rec := e.waitState(tsk, StateDoneNotOK); rec.Status.Error.Message != "stale_t passed" {
		t.Fatalf("%+v", rec.Status)
	}
}

func TestCommandRejections(t *testing.T) {
	e := newEnv(t)
	msn := newMission(e)
	cases := []struct {
		name string
		id   string
		verb string
		data func(id string) any
		mod  func(map[string]any)
		code string
	}{
		{"expired", newID("tsk"), "create", func(id string) any { return taskData(id, msn, "onemind.task.goto.v1", gotoSpec, simX1) },
			func(m map[string]any) { m["expires_at"] = rfc3339ms(e.clock.Now().Add(-time.Second)) }, CodeExpired},
		{"no approval", newID("tsk"), "create", func(id string) any { return taskData(id, msn, "onemind.task.goto.v1", gotoSpec, simX1) },
			func(m map[string]any) { m["approval"] = map[string]any{"status": "pending"} }, CodeNotAuthorized},
		{"not simulated", newID("tsk"), "create", func(id string) any {
			return taskData(id, msn, "onemind.task.goto.v1", gotoSpec, Ref{EntityType: "drone", EntityID: "x47"})
		}, nil, CodeNotSimulated},
		{"unknown type", newID("tsk"), "create", func(id string) any { return taskData(id, msn, "onemind.task.engage.v1", gotoSpec, simX1) }, nil, CodeUnknownAction},
		{"type does not fit", newID("tsk"), "create", func(id string) any {
			return taskData(id, msn, "onemind.task.manual.v1", map[string]any{"instructions": "x"}, simX1)
		}, nil, CodeSchemaInvalid},
		{"missing parent", newID("tsk"), "create", func(id string) any { return taskData(id, newID("msn"), "onemind.task.goto.v1", gotoSpec, simX1) }, nil, CodePreconditionFailed},
		{"id mismatch", newID("tsk"), "create", func(id string) any { return taskData(newID("tsk"), msn, "onemind.task.goto.v1", gotoSpec, simX1) }, nil, CodeSchemaInvalid},
		{"waypoint on a task", newID("tsk"), "waypoint", func(id string) any { return map[string]any{"seq": 0} }, nil, CodeUnknownAction},
		{"cancel is not a verb", msn, "cancel", func(id string) any { return map[string]any{} }, nil, CodeUnknownAction},
		{"bad waypoint", msn, "waypoint", func(id string) any { return map[string]any{"seq": 1, "lat": 95, "lon": 0, "hae_m": 1} }, nil, CodeSchemaInvalid},
		{"mission without name", newID("msn"), "create", func(id string) any { return map[string]any{"intent": "x"} }, nil, CodeSchemaInvalid},
		{"abort unknown task", newID("tsk"), "abort", func(id string) any { return map[string]any{"task_id": id} }, nil, CodeNotFound},
	}
	for _, c := range cases {
		var mods []func(map[string]any)
		if c.mod != nil {
			mods = append(mods, c.mod)
		}
		cid := e.command(c.id, c.verb, c.data(c.id), mods...)
		f := e.waitFactCause(c.id, "rejected", cid)
		d := f["data"].(map[string]any)
		if d["code"] != c.code || f["causation"] != cid {
			t.Errorf("%s: got code %v causation %v, want %s / %s", c.name, d["code"], f["causation"], c.code, cid)
		}
		if strings.HasPrefix(c.id, "tsk-") && c.verb == "create" {
			if _, err := e.kv.Get(context.Background(), c.id); err == nil {
				t.Errorf("%s: a record was written", c.name)
			}
		}
	}
	e.validateFacts()
}

// A redelivered create for the same action reports created again; another action conflicts.
func TestCreateIdempotent(t *testing.T) {
	e := newEnv(t)
	tsk := newID("tsk")
	d := taskData(tsk, "", "onemind.task.goto.v1", gotoSpec, simX1)
	e.command(tsk, "create", d)
	e.waitFact(tsk, "created")
	e.command(tsk, "create", d) // same action_id (redelivery)
	e.settle()
	if strings.Contains(strings.Join(e.subjects(tsk), ","), "rejected") {
		t.Fatalf("redelivery rejected: %v", e.subjects(tsk))
	}
	d2 := taskData(tsk, "", "onemind.task.goto.v1", gotoSpec, simX1) // new action_id
	e.command(tsk, "create", d2)
	if f := e.waitFact(tsk, "rejected"); f["data"].(map[string]any)["code"] != CodeConflict {
		t.Fatalf("%v", f)
	}
	if rec := e.task(tsk); rec.ActionID != d["action_id"] || rec.Version.Status != 1 {
		t.Fatalf("record changed: %+v", rec)
	}
}

func TestMissionVerbs(t *testing.T) {
	e := newEnv(t)
	msn := newMission(e)
	for _, w := range []map[string]any{
		{"seq": 2, "lat": 30.27, "lon": -97.74, "hae_m": 180, "name": "B"},
		{"seq": 1, "lat": 30.26, "lon": -97.74, "hae_m": 180, "name": "A"},
		{"seq": 2, "lat": 30.28, "lon": -97.74, "hae_m": 190, "name": "B2"},
	} {
		e.command(msn, "waypoint", w)
	}
	e.settle()
	m := e.mission(msn)
	if len(m.Route) != 2 || m.Route[0].Name != "A" || m.Route[1].Name != "B2" || m.Route[1].HaeM != 190 {
		t.Fatalf("route %+v", m.Route)
	}
	// Abort: an unsent child is cancelled at once; a sent one goes to cancel_requested.
	unsent := newID("tsk")
	e.command(unsent, "create", taskData(unsent, msn, "onemind.task.manual.v1", map[string]any{"instructions": "check the gate"},
		Ref{EntityType: "person", EntityID: "zeus"}))
	e.waitFact(unsent, "created")
	sent, _ := newSentTask(e, msn)
	abortCmd := e.command(msn, "abort", map[string]any{"reason": "weather"})
	f := e.waitFact(msn, "aborted")
	if f["causation"] != abortCmd {
		t.Fatalf("aborted causation %v", f["causation"])
	}
	if rec := e.waitState(unsent, StateDoneNotOK); rec.Status.Error.Code != ErrCancelled {
		t.Fatalf("unsent child %+v", rec.Status)
	}
	e.waitState(sent, StateCancelRequested)
	if m := e.mission(msn); m.State != MissionAborted || m.Rollup != "cancelled" {
		t.Fatalf("%+v", m)
	}
	// A finished mission refuses further verbs and new tasks.
	e.command(msn, "waypoint", map[string]any{"seq": 3, "lat": 1, "lon": 1, "hae_m": 1})
	if f := e.waitFact(msn, "rejected"); f["data"].(map[string]any)["code"] != CodePreconditionFailed {
		t.Fatalf("%v", f)
	}
	late := newID("tsk")
	e.command(late, "create", taskData(late, msn, "onemind.task.goto.v1", gotoSpec, simX1))
	if f := e.waitFact(late, "rejected"); f["data"].(map[string]any)["code"] != CodePreconditionFailed {
		t.Fatalf("%v", f)
	}
	msn2 := newMission(e)
	e.command(msn2, "complete", map[string]any{"note": "done"})
	e.waitFact(msn2, "completed")
	if m := e.mission(msn2); m.State != MissionCompleted {
		t.Fatalf("%+v", m)
	}
	e.validateFacts()
}

func TestScheduledAndComplete(t *testing.T) {
	e := newEnv(t)
	tsk := newID("tsk")
	d := taskData(tsk, "", "onemind.task.manual.v1", map[string]any{"instructions": "x"}, Ref{EntityType: "person", EntityID: "zeus"})
	start := e.clock.Now().Add(time.Hour).UnixMilli()
	d["execution"] = map[string]any{"start_after_t": start}
	e.command(tsk, "create", d)
	e.waitState(tsk, StateScheduled)
	e.command(tsk, "complete", map[string]any{})
	e.waitState(tsk, StateCompleteRequested)
	if got := strings.Join(e.subjects(tsk), ","); got != "created,scheduled,complete_requested" {
		t.Fatalf("%s", got)
	}
}

// People and agents report through svc.mission.report with a verified identity.
func TestReportRPC(t *testing.T) {
	e := newEnv(t)
	msn := newMission(e)
	tsk := newID("tsk")
	e.command(tsk, "create", taskData(tsk, msn, "onemind.task.manual.v1", map[string]any{"instructions": "open the gate"},
		Ref{EntityType: "person", EntityID: "zeus"}))
	e.waitFact(tsk, "created")
	e.dispatched(tsk, "task.assign", uuid4(), Ref{EntityType: "person", EntityID: "zeus"})
	e.waitState(tsk, StateSent)

	r := func(user, state string, ver int) map[string]any {
		return e.request(user, "svc.mission.report", map[string]any{"task_id": tsk, "state": state, "status_version": ver,
			"at_t": e.clock.Now().UnixMilli(), "simulated": false, "progress": map[string]any{"pct": 50}})
	}
	if v := r("alice", StateAcknowledged, 1); errCode(v) != CodeNotAuthorized {
		t.Fatalf("alice: %v", v)
	}
	if v := r("admin", StateAcknowledged, 1); errCode(v) != CodeNotAuthorized { // same account: unverified
		t.Fatalf("same-account caller: %v", v)
	}
	if v := r("zeus", StateAcknowledged, 1); v["accepted"] != true {
		t.Fatalf("zeus: %v", v)
	}
	if v := r("zeus", StateExecuting, 1); v["accepted"] == true { // not newer
		t.Fatalf("stale accepted: %v", v)
	}
	if v := r("zeus", StateDoneOK, 2); v["accepted"] != true || v["state"] != StateDoneOK {
		t.Fatalf("done: %v", v)
	}
	if rec := e.task(tsk); string(rec.Status.Progress) != `{"pct":50}` || rec.UpdatedBy != "user:zeus" {
		t.Fatalf("%+v", rec.Status)
	}
	e.waitFact(msn, "completed")

	// An agent task, reported by Hermes from the AGENTS account.
	at := newID("tsk")
	hermes := Ref{EntityType: "agent", EntityID: "hermes"}
	e.command(at, "create", taskData(at, "", "onemind.task.agent_request.v1", map[string]any{"question": "what moved?", "context_entity_ids": []string{}}, hermes))
	e.waitFact(at, "created")
	v := e.request("hermes", "svc.mission.report", map[string]any{"task_id": at, "state": "done_ok", "status_version": 1,
		"at_t": 1, "simulated": false, "result": map[string]any{"answer": "nothing"}})
	if v["accepted"] != true {
		t.Fatalf("hermes: %v", v)
	}
	// An asset task cannot be reported through the RPC.
	asset := newID("tsk")
	e.command(asset, "create", taskData(asset, "", "onemind.task.goto.v1", gotoSpec, simX1))
	e.waitFact(asset, "created")
	if v := e.request("zeus", "svc.mission.report", map[string]any{"task_id": asset, "state": "done_ok", "status_version": 9, "at_t": 1}); errCode(v) != CodeNotAuthorized {
		t.Fatalf("asset via RPC: %v", v)
	}
	// Team tasks: a listed member reports.
	team := newID("tsk")
	e.command(team, "create", taskData(team, "", "onemind.task.manual.v1", map[string]any{"instructions": "x"}, Ref{EntityType: "team", EntityID: "team-north"}))
	e.waitFact(team, "created")
	if v := e.request("zeus", "svc.mission.report", map[string]any{"task_id": team, "state": "received", "status_version": 1, "at_t": 1}); v["accepted"] != true {
		t.Fatalf("team: %v", v)
	}
	e.validateFacts()
}

func TestApplyGetList(t *testing.T) {
	e := newEnv(t)
	msn := newMission(e)
	act := newID("act")
	v := e.request("policy", "svc.mission.apply", map[string]any{"v": 1, "action_id": act, "op": "assign", "mission_id": msn,
		"to": "person.zeus", "role": "lead", "by": "user:zeus"})
	if v["record"] == nil {
		t.Fatalf("apply: %v", v)
	}
	f := e.waitFact(msn, "assigned")
	if f["data"].(map[string]any)["action_id"] != act {
		t.Fatalf("assigned fact lacks action_id: %v", f)
	}
	if v := e.request("policy", "svc.mission.apply", map[string]any{"v": 1, "action_id": act, "op": "assign", "mission_id": msn,
		"to": "drone.x47"}); errCode(v) != CodeNotSimulated {
		t.Fatalf("non-sim asset: %v", v)
	}
	// A client reaching apply through an import is refused by the service itself.
	if v := e.request("zeus", "svc.mission.apply", map[string]any{"v": 1, "action_id": act, "op": "assign", "mission_id": msn, "to": "person.alice"}); errCode(v) != CodeNotAuthorized {
		t.Fatalf("client apply: %v", v)
	}
	tsk, _ := newSentTask(e, msn)

	if v := e.request("zeus", "svc.mission.get", map[string]any{"v": 1, "id": tsk}); v["record"].(map[string]any)["task_id"] != tsk {
		t.Fatalf("get: %v", v)
	}
	if v := e.request("hermes", "svc.mission.get", map[string]any{"v": 1, "id": newID("msn")}); errCode(v) != CodeNotFound {
		t.Fatalf("get missing: %v", v)
	}
	e.settle()
	time.Sleep(100 * time.Millisecond) // KV watch copy
	list := func(q map[string]any) int {
		v := e.request("zeus", "svc.mission.list", q)
		return int(v["count"].(float64))
	}
	if n := list(map[string]any{}); n != 2 {
		t.Fatalf("list all: %d", n)
	}
	if n := list(map[string]any{"kind": "task", "assignee": "drone.sim-x1", "state": "sent"}); n != 1 {
		t.Fatalf("list by assignee/state: %d", n)
	}
	if n := list(map[string]any{"parent_id": msn}); n != 1 {
		t.Fatalf("list by parent: %d", n)
	}
	if n := list(map[string]any{"kind": "task", "bbox": []float64{-98, 30, -97, 31}}); n != 1 {
		t.Fatalf("bbox in: %d", n)
	}
	if n := list(map[string]any{"kind": "task", "bbox": []float64{0, 0, 1, 1}}); n != 0 {
		t.Fatalf("bbox out: %d", n)
	}
	if n := list(map[string]any{"kind": "mission", "assignee": "person.zeus"}); n != 1 {
		t.Fatalf("mission by assignee: %d", n)
	}
	v = e.request("policy", "svc.mission.apply", map[string]any{"v": 1, "action_id": act, "op": "unassign", "mission_id": msn, "to": "person.zeus"})
	if v["record"] == nil {
		t.Fatalf("unassign: %v", v)
	}
	e.waitFact(msn, "unassigned")
	e.validateFacts()
}

// The asset reaching the goto point completes the task from fleet telemetry (D2-P2).
func TestArrivalFromTelemetry(t *testing.T) {
	e := newEnv(t)
	msn := newMission(e)
	tsk, fleet := newSentTask(e, msn)
	e.report(tsk, StateExecuting, 1, fleet, nil)
	e.waitState(tsk, StateExecuting)
	time.Sleep(100 * time.Millisecond)
	pub := func(lat, lon, hae float64) {
		b, _ := json.Marshal(map[string]any{"uid": "sim-x1", "type": "fleet.drone", "lat": lat, "lon": lon, "hae": hae})
		_ = e.conn("bridge").Publish("tel.fleet.drone.sim-x1.position", b)
		_ = e.conn("bridge").Flush()
	}
	pub(30.2600, -97.7420, 180) // 1 km away
	pub(30.2690, -97.7420, 160) // overhead, 20 m low
	time.Sleep(200 * time.Millisecond)
	if st := e.task(tsk).Status.State; st != StateExecuting {
		t.Fatalf("completed too early: %s", st)
	}
	pub(30.26902, -97.7420, 181.5) // 2.2 m away, 1.5 m high
	rec := e.waitState(tsk, StateDoneOK)
	var res map[string]any
	_ = json.Unmarshal(rec.Status.Result, &res)
	if res["arrived"] != true || res["source"] != "telemetry" || rec.UpdatedBy != "service:task" {
		t.Fatalf("result %v by %s", res, rec.UpdatedBy)
	}
	e.waitFact(msn, "completed")
	// A later bridge report is dropped (terminal).
	e.report(tsk, StateExecuting, 2, fleet, nil)
	e.settle()
	if e.task(tsk).Status.State != StateDoneOK {
		t.Fatal("terminal changed")
	}
}

// After a restart the service reloads mission-state, keeps the guard, and applies the
// commands its durable consumer kept for it (unless they expired).
func TestRestart(t *testing.T) {
	e := newEnv(t)
	msn := newMission(e)
	tsk, fleet := newSentTask(e, msn)
	e.report(tsk, StateExecuting, 4, fleet, nil)
	e.waitState(tsk, StateExecuting)
	e.settle()
	e.svc.Stop()
	e.svc = nil

	late := newID("tsk")
	e.command(late, "create", taskData(late, msn, "onemind.task.goto.v1", gotoSpec, simX1))
	stale := newID("tsk")
	e.command(stale, "create", taskData(stale, msn, "onemind.task.goto.v1", gotoSpec, simX1),
		func(m map[string]any) { m["expires_at"] = rfc3339ms(e.clock.Now().Add(time.Second)) })
	e.report(tsk, StateWaiting, 3, fleet, nil) // stale against the stored guard
	e.clock.Advance(2 * time.Second)
	e.start()
	e.waitFact(late, "created")
	if f := e.waitFact(stale, "rejected"); f["data"].(map[string]any)["code"] != CodeExpired {
		t.Fatalf("%v", f)
	}
	e.settle()
	if st := e.task(tsk).Status.State; st != StateExecuting {
		t.Fatalf("guard lost on restart: %s", st)
	}
	e.clock.Advance(2 * time.Minute)
	e.svc.SweepOnce() // the reloaded copy drives the sweep
	e.waitState(late, StateDoneNotOK)
}

// The compare-and-set write retries when another writer moved the revision.
func TestCASConflictRetries(t *testing.T) {
	e := newEnv(t)
	tsk := newID("tsk")
	e.command(tsk, "create", taskData(tsk, "", "onemind.task.goto.v1", gotoSpec, simX1))
	e.waitFact(tsk, "created")
	calls := 0
	ctx := context.Background()
	got, err := e.svc.st.mutateTask(ctx, tsk, func(t *Task) error {
		calls++
		if calls == 1 {
			// Another writer updates the record between our read and our write.
			cur, _ := e.kv.Get(ctx, tsk)
			var other Task
			_ = json.Unmarshal(cur.Value(), &other)
			other.Description = "other writer"
			b, _ := json.Marshal(other)
			if _, err := e.kv.Put(ctx, tsk, b); err != nil {
				panic(err)
			}
		}
		t.Description += "+ours"
		return nil
	})
	if err != nil || calls != 2 || got.Description != "other writer+ours" {
		t.Fatalf("calls=%d err=%v desc=%q", calls, err, got.Description)
	}
}

func TestRollupTable(t *testing.T) {
	for _, c := range []struct{ state, code, want string }{
		{"created", "", "requested"}, {"scheduled", "", "requested"}, {"sent", "", "requested"},
		{"received", "", "running"}, {"executing", "", "running"}, {"cancel_requested", "", "running"},
		{"done_ok", "", "completed"}, {"done_not_ok", "rejected", "failed"}, {"done_not_ok", "cancelled", "cancelled"},
		{"done_not_ok", "timeout", "failed"}, {"done_not_ok", "no_asset", "failed"}, {"replaced", "", "cancelled"},
	} {
		if got := Rollup(c.state, c.code); got != c.want {
			t.Errorf("Rollup(%s,%s)=%s want %s", c.state, c.code, got, c.want)
		}
	}
	if d := HaversineM(30.269, -97.742, 30.26902, -97.742); d < 2 || d > 2.5 {
		t.Fatalf("haversine %f", d)
	}
}

// The dispatch audit of a person task can arrive before the create is applied (the policy
// service audits it right after publishing the create): sent is still recorded.
func TestDispatchedBeforeCreate(t *testing.T) {
	e := newEnv(t)
	tsk := newID("tsk")
	zeus := Ref{EntityType: "person", EntityID: "zeus"}
	e.dispatched(tsk, "task.assign", uuid4(), zeus)
	time.Sleep(300 * time.Millisecond)
	e.command(tsk, "create", taskData(tsk, "", "onemind.task.manual.v1", map[string]any{"instructions": "x"}, zeus))
	e.waitState(tsk, StateSent)
	if got := strings.Join(e.subjects(tsk), ","); got != "created,sent" {
		t.Fatalf("facts %s", got)
	}
}
