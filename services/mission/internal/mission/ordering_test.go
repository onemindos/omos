package mission

import (
	"encoding/json"
	"strings"
	"sync"
	"testing"
	"time"
)

// executingTask returns a goto task for sim-x1 that is executing on its assign command.
func executingTask(e *env) (tsk, fleet string) {
	e.t.Helper()
	tsk, fleet = newSentTask(e, newMission(e))
	e.report(tsk, StateReceived, 1, fleet, nil)
	e.report(tsk, StateExecuting, 2, fleet, nil)
	e.waitState(tsk, StateExecuting)
	return tsk, fleet
}

func wantCancelled(e *env, tsk string) Task {
	e.t.Helper()
	rec := e.waitState(tsk, StateDoneNotOK)
	if rec.Status.Error == nil || rec.Status.Error.Code != ErrCancelled || rec.Rollup != "cancelled" {
		e.t.Fatalf("abort result %+v", rec.Status)
	}
	if g := e.task(tsk).Guard; g.Hold != nil { // fresh read: waitState reuses its struct
		e.t.Fatalf("hold report not cleared: %+v", g)
	}
	return rec
}

// The hold's dispatch audit (EVT) is applied before the abort command (CMD), and the hold
// reports arrive after the abort: the hold is known by its audit, so the task is cancelled.
func TestAbortAuditBeforeAbortCommand(t *testing.T) {
	e := newEnv(t)
	tsk, _ := executingTask(e)
	hold := uuid4()
	e.dispatched(tsk, "task.abort", hold, simX1)
	e.settle()
	if g := e.task(tsk).Guard; !contains(g.AbortCommands, hold) {
		t.Fatalf("abort command not recorded: %+v", g)
	}
	e.command(tsk, "abort", map[string]any{"task_id": tsk})
	e.waitState(tsk, StateCancelRequested)
	e.report(tsk, StateReceived, 1, hold, nil)
	e.report(tsk, StateExecuting, 2, hold, nil)
	wantCancelled(e, tsk)
	e.validateFacts()
}

// Audit first, then the hold reports, then the abort command: the hold reports keep the
// state (no regression to received) and take effect when the abort lands.
func TestAbortAuditAndHoldReportsBeforeAbortCommand(t *testing.T) {
	e := newEnv(t)
	tsk, _ := executingTask(e)
	hold := uuid4()
	e.dispatched(tsk, "task.abort", hold, simX1)
	e.report(tsk, StateReceived, 1, hold, nil)
	e.report(tsk, StateExecuting, 2, hold, nil)
	e.settle()
	if rec := e.task(tsk); rec.Status.State != StateExecuting || rec.Guard.Hold == nil || rec.Guard.Hold.State != StateExecuting {
		t.Fatalf("pre-abort hold reports must keep the state and be kept: %+v %+v", rec.Status, rec.Guard)
	}
	abort := e.command(tsk, "abort", map[string]any{"task_id": tsk})
	rec := wantCancelled(e, tsk)
	if got := strings.Join(e.subjects(tsk), ","); !strings.HasSuffix(got, "executing,cancel_requested,done_not_ok") {
		t.Fatalf("task facts %s", got)
	}
	if f := e.waitFact(tsk, StateDoneNotOK); f["causation"] != abort {
		t.Fatalf("cancelled fact causation %v, want the abort command %s", f["causation"], abort)
	}
	_ = rec
	e.validateFacts()
}

// The hold's first report beats both the abort command and the hold's audit: it does not
// move the state back to received, and the hold executing after the abort cancels the task.
func TestHoldReportBeforeAbortAndAudit(t *testing.T) {
	e := newEnv(t)
	tsk, _ := executingTask(e)
	hold := uuid4()
	e.report(tsk, StateReceived, 1, hold, nil)
	e.settle()
	if st := e.task(tsk).Status.State; st != StateExecuting {
		t.Fatalf("a hold report before the abort moved the state to %s", st)
	}
	e.command(tsk, "abort", map[string]any{"task_id": tsk})
	e.waitState(tsk, StateCancelRequested)
	e.dispatched(tsk, "task.abort", hold, simX1)
	e.report(tsk, StateExecuting, 2, hold, nil)
	wantCancelled(e, tsk)
	e.validateFacts()
}

// Every hold report beats the abort command and the audit: the abort resolves at once.
func TestHoldExecutedBeforeAbortAndAudit(t *testing.T) {
	e := newEnv(t)
	tsk, _ := executingTask(e)
	hold := uuid4()
	e.report(tsk, StateReceived, 1, hold, nil)
	e.report(tsk, StateExecuting, 2, hold, nil)
	e.settle()
	e.command(tsk, "abort", map[string]any{"task_id": tsk})
	wantCancelled(e, tsk)
	e.dispatched(tsk, "task.abort", hold, simX1) // late audit: task terminal, nothing changes
	e.settle()
	if rec := e.task(tsk); rec.Status.State != StateDoneNotOK {
		t.Fatalf("late audit changed the task: %+v", rec.Status)
	}
	e.validateFacts()
}

// The asset refuses the hold before the abort is applied: once applied, the task returns to
// its previous state with error rejected.
func TestHoldRefusedBeforeAbort(t *testing.T) {
	e := newEnv(t)
	tsk, _ := executingTask(e)
	hold := uuid4()
	e.dispatched(tsk, "task.abort", hold, simX1)
	e.report(tsk, StateDoneNotOK, 1, hold, map[string]string{"code": "rejected", "message": "link lost"})
	e.settle()
	if st := e.task(tsk).Status.State; st != StateExecuting {
		t.Fatalf("refused hold before abort changed the state to %s", st)
	}
	e.command(tsk, "abort", map[string]any{"task_id": tsk})
	e.waitFact(tsk, StateCancelRequested)
	e.settle()
	rec := e.task(tsk)
	if got := strings.Join(e.subjects(tsk), ","); !strings.HasSuffix(got, "executing,cancel_requested,executing") {
		t.Fatalf("task facts %s", got)
	}
	if rec.Status.State != StateExecuting || rec.Status.Error == nil || rec.Status.Error.Code != ErrRejected || !strings.Contains(rec.Status.Error.Message, "link lost") {
		t.Fatalf("refused abort: %+v", rec.Status)
	}
}

// cmdBytes returns the recorded cmd.mission envelope with the given id.
func cmdBytes(e *env, id string) (string, []byte) {
	e.t.Helper()
	for _, m := range e.cmds.all() {
		var v struct {
			ID string `json:"id"`
		}
		if json.Unmarshal(m.Data, &v) == nil && v.ID == id {
			return m.Subject, m.Data
		}
	}
	e.t.Fatalf("command %s not recorded", id)
	return "", nil
}

// A command applied whose ack was lost is redelivered after expires_at: it is recognised as
// applied and produces no rejection (abort, complete, waypoint, task and mission create).
func TestRedeliveryOfAppliedCommand(t *testing.T) {
	e := newEnv(t)
	msn := newMission(e)
	wp := e.command(msn, "waypoint", map[string]any{"seq": 0, "lat": 30.27, "lon": -97.74, "hae_m": 50})
	e.waitFact(msn, "waypoint")
	tsk, _ := executingTask(e)
	abort := e.command(tsk, "abort", map[string]any{"task_id": tsk})
	e.waitState(tsk, StateCancelRequested)
	t2 := newID("tsk")
	create := e.command(t2, "create", taskData(t2, "", "onemind.task.goto.v1", gotoSpec, simX1))
	e.waitFact(t2, "created")
	complete := e.command(t2, "complete", map[string]any{"task_id": t2})
	e.waitState(t2, StateCompleteRequested)
	e.settle()

	e.clock.Advance(time.Minute) // every command is past expires_at now
	before := e.svc.C.CommandsRejected.Load()
	for _, id := range []string{wp, abort, create, complete} {
		subj, data := cmdBytes(e, id)
		if err := e.svc.HandleCommand(subj, data); err != nil {
			t.Fatalf("redelivery of %s: %v", subj, err)
		}
	}
	e.settle()
	if n := e.svc.C.CommandsRejected.Load() - before; n != 0 {
		t.Fatalf("%d redeliveries were rejected", n)
	}
	for _, id := range []string{msn, tsk, t2} {
		if s := strings.Join(e.subjects(id), ","); strings.Contains(s, "rejected") {
			t.Fatalf("%s: spurious rejection in %s", id, s)
		}
	}
	if e.svc.C.CommandsRedelivered.Load() != 4 {
		t.Fatalf("redelivered counter %d, want 4", e.svc.C.CommandsRedelivered.Load())
	}
	// A new command past expiry is still refused.
	e2 := e.command(tsk, "complete", map[string]any{"task_id": tsk}, func(m map[string]any) {
		m["expires_at"] = rfc3339ms(e.clock.Now().Add(-time.Second))
	})
	if f := e.waitFactCause(tsk, "rejected", e2); f["data"].(map[string]any)["code"] != CodeExpired {
		t.Fatalf("expired new command: %v", f)
	}
}

// svc.mission.list caps one reply by size and pages with after/next.
func TestListSizeCapAndPaging(t *testing.T) {
	e := newEnv(t)
	big := strings.Repeat("x", 100*1024)
	want := map[string]bool{}
	for i := 0; i < 12; i++ {
		id := newID("tsk")
		d := taskData(id, "", "onemind.task.goto.v1", gotoSpec, simX1)
		d["description"] = big
		e.command(id, "create", d)
		want[id] = true
	}
	e.settle()
	time.Sleep(100 * time.Millisecond)
	seen := map[string]bool{}
	after, pages := "", 0
	for {
		q := map[string]any{"kind": "task"}
		if after != "" {
			q["after"] = after
		}
		v := e.request("zeus", "svc.mission.list", q)
		if v["error"] != nil {
			t.Fatalf("list: %v", v["error"])
		}
		pages++
		for _, it := range v["items"].([]any) {
			id := it.(map[string]any)["task_id"].(string)
			if seen[id] {
				t.Fatalf("%s listed twice", id)
			}
			seen[id] = true
		}
		if v["truncated"] != true {
			break
		}
		after = v["next"].(string)
		if pages > 5 {
			t.Fatal("paging does not end")
		}
	}
	if pages < 2 || len(seen) != len(want) {
		t.Fatalf("pages %d, listed %d of %d", pages, len(seen), len(want))
	}
	// Count paging works too.
	v := e.request("zeus", "svc.mission.list", map[string]any{"kind": "task", "limit": 5})
	if v["count"].(float64) != 5 || v["truncated"] != true || v["next"] == nil {
		t.Fatalf("limit 5: count %v truncated %v", v["count"], v["truncated"])
	}
}

// Terminal records leave the sweep set at once and the in-memory copy after
// CacheTerminalFor; svc.mission.get still reads them from KV.
func TestTerminalEviction(t *testing.T) {
	e := newEnv(t)
	msn := newMission(e)
	tsk, fleet := newSentTask(e, msn)
	e.report(tsk, StateDoneOK, 1, fleet, nil)
	e.waitState(tsk, StateDoneOK)
	e.waitFact(msn, "completed")
	e.settle()
	time.Sleep(100 * time.Millisecond)
	for _, o := range e.svc.st.openTasks() {
		if o.TaskID == tsk {
			t.Fatal("terminal task still in the sweep set")
		}
	}
	if n := e.svc.EvictOnce(); n != 0 {
		t.Fatalf("evicted %d records before CacheTerminalFor", n)
	}
	e.clock.Advance(25 * time.Hour)
	if n := e.svc.EvictOnce(); n != 2 {
		t.Fatalf("evicted %d records, want the task and the mission", n)
	}
	if v := e.request("zeus", "svc.mission.list", map[string]any{}); v["count"].(float64) != 0 {
		t.Fatalf("evicted records still listed: %v", v["count"])
	}
	if v := e.request("zeus", "svc.mission.get", map[string]any{"v": 1, "id": tsk}); v["record"] == nil {
		t.Fatalf("get of an evicted record: %v", v)
	}
}

// Defence in depth (contract section 10 test 4): an approved create whose target is a
// person or a person-class track is refused forbidden_action and no record is written.
func TestPersonTargetRefused(t *testing.T) {
	e := newEnv(t)
	observe := map[string]any{"target": map[string]any{"entity_id": "trk-p1"}, "sensor_id": "sim-x1.cam0"}
	cases := []func(d map[string]any){
		func(d map[string]any) {
			d["target_snapshots"] = []any{map[string]any{"entity_id": "trk-p1", "snapshot": map[string]any{"entity_type": "track", "kind": "person"}}}
		},
		func(d map[string]any) {
			d["spec"] = map[string]any{"target": map[string]any{"entity_id": "p-1", "entity_type": "person"}, "sensor_id": "sim-x1.cam0"}
		},
	}
	for i, mod := range cases {
		id := newID("tsk")
		d := taskData(id, "", "onemind.task.observe.v1", observe, simX1)
		mod(d)
		cid := e.command(id, "create", d)
		f := e.waitFactCause(id, "rejected", cid)
		if f["data"].(map[string]any)["code"] != CodeForbiddenAction {
			t.Fatalf("case %d: %v", i, f)
		}
		e.settle()
		if _, ok := e.svc.st.cachedTask(id); ok {
			t.Fatalf("case %d: a record was written", i)
		}
	}
	// A point target is fine.
	id := newID("tsk")
	e.command(id, "create", taskData(id, "", "onemind.task.observe.v1",
		map[string]any{"target": map[string]any{"lat": 30.27, "lon": -97.74, "hae_m": 30}, "sensor_id": "sim-x1.cam0"}, simX1))
	e.waitFact(id, "created")
}

// E2E-2 race: the bridge's received report (EVT, evt.drone.{id}.task) is applied before the
// policy service's dispatch audit (EVT, evt.so.action.dispatched), about 1 ms apart. The
// report proves the task was sent, so sent is recorded first; the late audit changes nothing.
func TestReceivedBeforeDispatchAudit(t *testing.T) {
	e := newEnv(t)
	msn := newMission(e)
	tsk := newID("tsk")
	e.command(tsk, "create", taskData(tsk, msn, "onemind.task.goto.v1", gotoSpec, simX1))
	e.waitFact(tsk, "created")
	fleet := uuid4()
	e.report(tsk, StateReceived, 1, fleet, nil)
	e.waitState(tsk, StateReceived)
	e.dispatched(tsk, "task.assign", fleet, simX1) // late audit
	e.report(tsk, StateExecuting, 2, fleet, nil)
	e.waitState(tsk, StateExecuting)
	e.report(tsk, StateDoneOK, 3, fleet, nil)
	rec := e.waitState(tsk, StateDoneOK)
	e.settle()
	if got := strings.Join(e.subjects(tsk), ","); got != "created,sent,received,executing,done_ok" {
		t.Fatalf("task facts %s", got)
	}
	for i, st := range []string{"sent", "received", "executing", "done_ok"} {
		if v := e.waitFact(tsk, st)["data"].(map[string]any)["status_version"]; v != float64(i+2) {
			t.Errorf("%s status_version %v, want %d", st, v, i+2)
		}
	}
	if f := e.waitFact(tsk, StateSent); f["causation"] != fleet || f["data"].(map[string]any)["updated_by"] != "service:task" {
		t.Fatalf("sent fact causation %v updated_by %v, want the fleet command and the task service", f["causation"], f["data"].(map[string]any)["updated_by"])
	}
	if rec.Version.Status != 5 || len(rec.Guard.Commands) != 1 || rec.Guard.Commands[0] != fleet {
		t.Fatalf("record %+v guard %+v", rec.Version, rec.Guard)
	}
	e.waitFact(msn, "completed")
	e.validateFacts()
}

// Both orders, concurrently: for each task the audit and the received report are published
// at the same time from two goroutines. Every task records sent exactly once, before received.
func TestDispatchAuditReportRace(t *testing.T) {
	e := newEnv(t)
	msn := newMission(e)
	const n = 12
	tasks, fleets := make([]string, n), make([]string, n)
	for i := range tasks {
		tasks[i], fleets[i] = newID("tsk"), uuid4()
		e.command(tasks[i], "create", taskData(tasks[i], msn, "onemind.task.goto.v1", gotoSpec, simX1))
		e.waitFact(tasks[i], "created")
	}
	e.conn("policy") // open both connections before the goroutines share the cache
	e.conn("bridge")
	var wg sync.WaitGroup
	for i := range tasks {
		wg.Add(2)
		go func(i int) { defer wg.Done(); e.dispatched(tasks[i], "task.assign", fleets[i], simX1) }(i)
		go func(i int) { defer wg.Done(); e.report(tasks[i], StateReceived, 1, fleets[i], nil) }(i)
	}
	wg.Wait()
	for _, tsk := range tasks {
		e.waitState(tsk, StateReceived)
	}
	e.settle()
	for i, tsk := range tasks {
		if got := strings.Join(e.subjects(tsk), ","); got != "created,sent,received" {
			t.Errorf("task %d facts %s", i, got)
		}
		if rec := e.task(tsk); rec.Version.Status != 3 || len(rec.Guard.Commands) != 1 {
			t.Errorf("task %d record %+v guard %+v", i, rec.Version, rec.Guard)
		}
	}
	e.validateFacts()
}
