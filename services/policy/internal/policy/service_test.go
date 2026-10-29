package policy

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/nats-io/nats.go"
)

func gotoReq(target string, lat, lon float64) map[string]any {
	return map[string]any{
		"v": 1, "request_id": NewUUID(), "action_type": "asset.goto", "action_version": 1,
		"target":       map[string]any{"entity_type": "drone", "entity_id": target},
		"params":       map[string]any{"lat": lat, "lon": lon, "hae_m": 160.0, "speed_mps": 8.0},
		"requested_by": map[string]any{"kind": "user", "id": "user:claimed"},
		"reason":       "test",
	}
}

func req(actionType, entityType, id string, params map[string]any) map[string]any {
	return map[string]any{
		"v": 1, "request_id": NewUUID(), "action_type": actionType, "action_version": 1,
		"target":       map[string]any{"entity_type": entityType, "entity_id": id},
		"params":       params,
		"requested_by": map[string]any{"kind": "user", "id": "user:claimed"},
	}
}

func (e *env) request(user, subject string, body any, hdr nats.Header) (*nats.Msg, error) {
	b, _ := json.Marshal(body)
	return e.conn(user).RequestMsg(&nats.Msg{Subject: subject, Data: b, Header: hdr}, 3*time.Second)
}

func (e *env) propose(user string, fam Family, body map[string]any) ProposeReply {
	e.t.Helper()
	m, err := e.request(user, "svc.policy."+string(fam)+".propose", body, nil)
	if err != nil {
		e.t.Fatalf("propose as %s: %v", user, err)
	}
	var r ProposeReply
	if err := json.Unmarshal(m.Data, &r); err != nil {
		e.t.Fatal(err)
	}
	return r
}

func (e *env) confirm(user string, fam Family, actionID, hash, decision string) ConfirmReply {
	e.t.Helper()
	body := map[string]any{"v": 1, "action_id": actionID, "params_hash": hash, "decision": decision}
	if decision == "reject" {
		body["reason"] = "not now"
	}
	m, err := e.request(user, "svc.policy."+string(fam)+".confirm", body, nil)
	if err != nil {
		e.t.Fatalf("confirm as %s: %v", user, err)
	}
	var r ConfirmReply
	if err := json.Unmarshal(m.Data, &r); err != nil {
		e.t.Fatal(err)
	}
	return r
}

type auditEvt struct {
	Subject string
	Env     struct {
		ID   string         `json:"id"`
		Data map[string]any `json:"data"`
	}
}

func (e *env) auditFor(actionID string) []auditEvt {
	var out []auditEvt
	for _, m := range e.audit.all() {
		var a auditEvt
		a.Subject = m.Subject
		_ = json.Unmarshal(m.Data, &a.Env)
		if a.Env.Data["action_id"] == actionID {
			out = append(out, a)
		}
	}
	return out
}

func (e *env) hasAudit(actionID, event, code string) bool {
	for _, a := range e.auditFor(actionID) {
		if a.Subject == "evt.so.action."+event && (code == "" || a.Env.Data["code"] == code) {
			return true
		}
	}
	return false
}

func (e *env) cmdsFor(actionID string) []map[string]any {
	var out []map[string]any
	for _, m := range e.cmds.all() {
		var c map[string]any
		_ = json.Unmarshal(m.Data, &c)
		if c["action_id"] == actionID {
			out = append(out, c)
		}
	}
	return out
}

func expectDenied(t *testing.T, r ProposeReply, code string) {
	t.Helper()
	if r.Status != StateDenied || r.Denial == nil || r.Denial.Code != code {
		t.Fatalf("want denied %s, got %+v (denial %+v)", code, r, r.Denial)
	}
}

func TestPolicyService(t *testing.T) {
	e := newEnv(t)

	t.Run("human goto: propose waits, confirm dispatches one signed command, outcomes audited", func(t *testing.T) {
		r := e.propose("zeus", FamilyCommand, gotoReq("sim-x1", 30.27, -97.742))
		if r.Status != StateAwaiting || r.Level != A3 || r.ApprovalsNeeded != 1 || !strings.HasPrefix(r.ParamsHash, "sha256:") {
			t.Fatalf("propose: %+v", r)
		}
		if n := len(e.cmdsFor(r.ActionID)); n != 0 {
			t.Fatalf("command before confirm: %d", n)
		}
		waitFor(t, func() bool { return e.hasAudit(r.ActionID, "proposed", "") }, "proposed event")
		ev := e.auditFor(r.ActionID)[0].Env.Data
		if rb := ev["requested_by"].(map[string]any); rb["id"] != "user:zeus" || ev["actor_verified"] != true {
			t.Fatalf("identity not taken from the connection: %v", ev["requested_by"])
		}
		c := e.confirm("zeus", FamilyCommand, r.ActionID, r.ParamsHash, "approve")
		if c.Status != StateDispatched || c.CommandID == "" {
			t.Fatalf("confirm: %+v (denial %+v)", c, c.Denial)
		}
		waitFor(t, func() bool { return e.hasAudit(r.ActionID, "executed", "") }, "executed republished")
		cmds := e.cmdsFor(r.ActionID)
		if len(cmds) != 1 || cmds[0]["subject"] != "cmd.fleet.drone.sim-x1.goto" || cmds[0]["simulated"] != true {
			t.Fatalf("commands: %v", cmds)
		}
		raw, _ := json.Marshal(cmds[0])
		if err := e.schemas.ValidateCommand(raw); err != nil {
			t.Fatalf("command envelope: %v", err)
		}
		for _, ev := range []string{"proposed", "approved", "dispatched", "acknowledged", "executed"} {
			if !e.hasAudit(r.ActionID, ev, "") {
				t.Fatalf("missing evt.so.action.%s", ev)
			}
		}
		for _, m := range e.audit.all() {
			if err := e.schemas.ValidateEvent(m.Data); err != nil {
				t.Fatalf("audit envelope on %s: %v", m.Subject, err)
			}
		}
	})

	t.Run("agent proposes but can never confirm", func(t *testing.T) {
		body := gotoReq("sim-x1", 30.268, -97.741)
		body["requested_by"] = map[string]any{"kind": "user", "id": "user:zeus", "on_behalf_of": "user:zeus"} // claims to be zeus
		body["prompt_sha256"] = "sha256:abc"
		r := e.propose("hermes", FamilyCommand, body)
		if r.Status != StateAwaiting {
			t.Fatalf("agent propose: %+v %+v", r, r.Denial)
		}
		waitFor(t, func() bool { return e.hasAudit(r.ActionID, "proposed", "") }, "proposed")
		ev := e.auditFor(r.ActionID)[0].Env.Data
		if rb := ev["requested_by"].(map[string]any); rb["kind"] != "agent" || rb["id"] != "agent:hermes/zeus/legacy" || ev["prompt_sha256"] != "sha256:abc" {
			t.Fatalf("agent identity or hashes wrong: %v", ev)
		}
		// Layer 1, the server: Hermes has no publish permission on svc.policy.*.confirm.
		_, err := e.request("hermes", SubjCommandConfirm, map[string]any{"v": 1, "action_id": r.ActionID, "params_hash": r.ParamsHash, "decision": "approve"}, nil)
		if err == nil {
			t.Fatal("hermes reached svc.policy.command.confirm")
		}
		waitFor(t, func() bool { return len(e.permErrors("hermes")) > 0 }, "permission violation for hermes")
		// Layer 2, the policy: an agent account that can reach confirm (misconfigured, and
		// mislabelled kind user in the directory) is still refused.
		c := e.confirm("rogue", FamilyCommand, r.ActionID, r.ParamsHash, "approve")
		if c.Denial == nil || c.Denial.Code != CodeAgentCannotConfirm || c.Status != StateAwaiting {
			t.Fatalf("rogue confirm: %+v %+v", c, c.Denial)
		}
		waitFor(t, func() bool { return e.hasAudit(r.ActionID, "denied", CodeAgentCannotConfirm) }, "agent confirm audited")
		if len(e.cmdsFor(r.ActionID)) != 0 {
			t.Fatal("a command went out after an agent confirm")
		}
		// The human still can.
		c = e.confirm("zeus", FamilyCommand, r.ActionID, r.ParamsHash, "approve")
		if c.Status != StateDispatched {
			t.Fatalf("human confirm after agent proposal: %+v %+v", c, c.Denial)
		}
	})

	t.Run("agent without a human on_behalf_of is refused", func(t *testing.T) {
		body := gotoReq("sim-x1", 30.268, -97.741)
		body["requested_by"] = map[string]any{"kind": "agent", "id": "x"}
		expectDenied(t, e.propose("hermes", FamilyCommand, body), CodeNotAuthorized)
	})

	t.Run("stop-line verbs are rejected by name and nothing reaches cmd.fleet", func(t *testing.T) {
		before := len(e.cmds.all())
		for _, at := range []string{"asset.strike", "asset.engage", "asset.intercept", "asset.designate",
			"asset.designate_target", "asset.follow_person", "asset.attack", "asset.weapon_release"} {
			r := e.propose("zeus", FamilyCommand, req(at, "drone", "sim-x1", map[string]any{}))
			expectDenied(t, r, CodeForbiddenAction)
			waitFor(t, func() bool { return e.hasAudit(r.ActionID, "denied", CodeForbiddenAction) }, at+" denial audited")
		}
		// Hermes asked to "follow that person": the proposal it could make is denied and recorded.
		body := req("asset.orbit", "drone", "sim-x1", map[string]any{"center": map[string]any{"entity_id": "person-1"},
			"radius_m": 30, "hae_m": 120, "duration_s": 600})
		body["requested_by"] = map[string]any{"kind": "agent", "id": "a", "on_behalf_of": "user:zeus"}
		r := e.propose("hermes", FamilyCommand, body)
		expectDenied(t, r, CodeForbiddenAction)
		waitFor(t, func() bool { return e.hasAudit(r.ActionID, "denied", CodeForbiddenAction) }, "follow-person denial audited")
		// A CoT strike-tasking type smuggled into a task is refused too.
		tr := req("task.assign", "drone", "sim-x1", map[string]any{"parent_id": nil, "type": "onemind.task.engage.v1", "spec": map[string]any{}})
		expectDenied(t, e.propose("zeus", FamilyTask, tr), CodeForbiddenAction)
		if len(e.cmds.all()) != before {
			t.Fatal("a stop-line proposal produced a command")
		}
	})

	t.Run("unknown verbs are refused (closed allowlist)", func(t *testing.T) {
		expectDenied(t, e.propose("zeus", FamilyCommand, req("asset.dance", "drone", "sim-x1", map[string]any{})), CodeUnknownAction)
	})

	t.Run("approval levels A1 A2 A3 A4 A5", func(t *testing.T) {
		// A1: runs without a confirm, still logged.
		r := e.propose("zeus", FamilyCommand, req("entity.star", "drone", "sim-x1", map[string]any{"starred": true}))
		if r.Status != StateExecuted || r.ApprovalsNeeded != 0 || len(e.cmdsFor(r.ActionID)) == 0 && !waitCmd(e, r.ActionID) {
			t.Fatalf("A1: %+v %+v", r, r.Denial)
		}
		// A2: one human.
		r = e.propose("zeus", FamilyCommand, req("entity.override", "drone", "sim-x2", map[string]any{"field": "name", "value": "Scout"}))
		if r.Level != A2 || r.ApprovalsNeeded != 1 || r.Status != StateAwaiting {
			t.Fatalf("A2: %+v", r)
		}
		// A4: two distinct humans; the same human twice is refused.
		r = e.propose("zeus", FamilyCommand, req("asset.arm", "drone", "sim-x1", map[string]any{}))
		if r.Level != A4 || r.ApprovalsNeeded != 2 {
			t.Fatalf("A4: %+v %+v", r, r.Denial)
		}
		c := e.confirm("zeus", FamilyCommand, r.ActionID, r.ParamsHash, "approve")
		if c.Status != StateAwaiting || c.ApprovalsReceived != 1 || len(e.cmdsFor(r.ActionID)) != 0 {
			t.Fatalf("A4 step 1: %+v", c)
		}
		c = e.confirm("zeus", FamilyCommand, r.ActionID, r.ParamsHash, "approve")
		if c.Denial == nil || c.Denial.Code != CodeDuplicateConfirmation {
			t.Fatalf("A4 same human twice: %+v", c)
		}
		c = e.confirm("alice", FamilyCommand, r.ActionID, r.ParamsHash, "approve")
		if c.Status != StateDispatched || c.ApprovalsReceived != 2 {
			t.Fatalf("A4 step 2: %+v %+v", c, c.Denial)
		}
		// A5: never exposed.
		expectDenied(t, e.propose("zeus", FamilyCommand, req("system.credentials", "drone", "sim-x1", map[string]any{})), CodeForbiddenAction)
	})

	t.Run("approvals are bound to action id and params hash", func(t *testing.T) {
		a := e.propose("zeus", FamilyCommand, gotoReq("sim-x1", 30.27, -97.74))
		b := e.propose("zeus", FamilyCommand, gotoReq("sim-x1", 30.26, -97.75))
		c := e.confirm("zeus", FamilyCommand, a.ActionID, b.ParamsHash, "approve")
		if c.Denial == nil || c.Denial.Code != CodeConflict || len(e.cmdsFor(a.ActionID)) != 0 {
			t.Fatalf("hash of another action accepted: %+v", c)
		}
		// Task actions are confirmed only on the task endpoint.
		c = e.confirm("zeus", FamilyTask, a.ActionID, a.ParamsHash, "approve")
		if c.Denial == nil || c.Denial.Code != CodeConflict {
			t.Fatalf("family mismatch accepted: %+v", c)
		}
		_ = e.confirm("zeus", FamilyCommand, b.ActionID, b.ParamsHash, "reject")
		waitFor(t, func() bool { return e.hasAudit(b.ActionID, "rejected", CodeRejectedByHuman) }, "reject audited")
	})

	t.Run("a replayed approval is refused and sends nothing", func(t *testing.T) {
		r := e.propose("zeus", FamilyCommand, gotoReq("sim-x1", 30.269, -97.742))
		c := e.confirm("zeus", FamilyCommand, r.ActionID, r.ParamsHash, "approve")
		if c.Status != StateDispatched {
			t.Fatalf("first confirm: %+v", c)
		}
		for _, who := range []string{"zeus", "alice"} {
			c = e.confirm(who, FamilyCommand, r.ActionID, r.ParamsHash, "approve")
			if c.Denial == nil || c.Denial.Code != CodeReplayed {
				t.Fatalf("replay by %s: %+v", who, c)
			}
		}
		time.Sleep(100 * time.Millisecond)
		if n := len(e.cmdsFor(r.ActionID)); n != 1 {
			t.Fatalf("replay produced %d commands", n)
		}
		// Re-publishing the recorded command itself (bypassing policy) is refused by the
		// bridge as a replayed nonce; after expires_at it is refused as expired.
		raw := e.cmds.all()
		var last *nats.Msg
		for _, m := range raw {
			var c map[string]any
			_ = json.Unmarshal(m.Data, &c)
			if c["action_id"] == r.ActionID {
				last = m
			}
		}
		before := e.bridge.count()
		e.bridge.handle(e.conn("bridge"), last)
		e.bridge.mu.Lock()
		e.bridge.seen = map[string]bool{}
		e.bridge.mu.Unlock()
		e.clock.Advance(31 * time.Second)
		e.bridge.handle(e.conn("bridge"), last)
		e.bridge.mu.Lock()
		got := append([]string(nil), e.bridge.refused...)
		e.bridge.mu.Unlock()
		if e.bridge.count() != before || len(got) < 2 || got[len(got)-2] != CodeReplayed || got[len(got)-1] != CodeExpired {
			t.Fatalf("bridge refusals: %v", got)
		}
	})

	t.Run("an unconfirmed action expires and cannot be confirmed later", func(t *testing.T) {
		r := e.propose("zeus", FamilyCommand, gotoReq("sim-x1", 30.27, -97.743))
		e.clock.Advance(121 * time.Second)
		waitFor(t, func() bool { return e.hasAudit(r.ActionID, "expired", CodeExpired) }, "expired by the sweeper")
		c := e.confirm("zeus", FamilyCommand, r.ActionID, r.ParamsHash, "approve")
		if c.Denial == nil || c.Denial.Code != CodeExpired || len(e.cmdsFor(r.ActionID)) != 0 {
			t.Fatalf("late confirm: %+v", c)
		}
	})

	t.Run("no replay to offline assets", func(t *testing.T) {
		expectDenied(t, e.propose("zeus", FamilyCommand, gotoReq("sim-x3", 30.27, -97.74)), CodePreconditionFailed)
		r := e.propose("zeus", FamilyCommand, gotoReq("sim-x2", 30.27, -97.74))
		if r.Status != StateAwaiting {
			t.Fatalf("%+v %+v", r, r.Denial)
		}
		e.setOnline("sim-x2", false)
		c := e.confirm("zeus", FamilyCommand, r.ActionID, r.ParamsHash, "approve")
		if c.Status != StateDenied || c.Denial == nil || c.Denial.Code != CodePreconditionFailed {
			t.Fatalf("confirm to offline asset: %+v %+v", c, c.Denial)
		}
		e.setOnline("sim-x2", true)
		time.Sleep(100 * time.Millisecond)
		if len(e.cmdsFor(r.ActionID)) != 0 {
			t.Fatal("command sent to an offline asset, or queued and sent on reconnect")
		}
		waitFor(t, func() bool { return e.hasAudit(r.ActionID, "denied", CodePreconditionFailed) }, "offline denial audited")
	})

	t.Run("no acknowledgement fails the action", func(t *testing.T) {
		e.bridge.mu.Lock()
		e.bridge.silent["sim-x2"] = true
		e.bridge.mu.Unlock()
		r := e.propose("zeus", FamilyCommand, gotoReq("sim-x2", 30.27, -97.74))
		c := e.confirm("zeus", FamilyCommand, r.ActionID, r.ParamsHash, "approve")
		if c.Status != StateDispatched {
			t.Fatalf("%+v %+v", c, c.Denial)
		}
		e.clock.Advance(11 * time.Second)
		waitFor(t, func() bool { return e.hasAudit(r.ActionID, "failed", CodeNoAck) }, "no_ack")
	})

	t.Run("simulated only, schema, capability, person targets, geofence", func(t *testing.T) {
		expectDenied(t, e.propose("zeus", FamilyCommand, gotoReq("real-1", 30.27, -97.74)), CodeNotSimulated)
		bad := gotoReq("sim-x1", 30.27, -97.74)
		delete(bad["params"].(map[string]any), "lat")
		expectDenied(t, e.propose("zeus", FamilyCommand, bad), CodeSchemaInvalid)
		bad = gotoReq("sim-x1", 30.27, -97.74)
		bad["params"].(map[string]any)["payload"] = "x"
		expectDenied(t, e.propose("zeus", FamilyCommand, bad), CodeSchemaInvalid)
		bad = gotoReq("sim-x1", 30.27, -97.74)
		bad["action_version"] = 2
		expectDenied(t, e.propose("zeus", FamilyCommand, bad), CodeSchemaInvalid)
		orbit := req("asset.orbit", "drone", "sim-x2", map[string]any{"center": map[string]any{"lat": 30.27, "lon": -97.74},
			"radius_m": 50, "hae_m": 120, "duration_s": 60})
		expectDenied(t, e.propose("zeus", FamilyCommand, orbit), CodeNotCapable)
		obs := req("asset.observe", "drone", "sim-x1", map[string]any{"target": map[string]any{"entity_id": "trk-9"}, "sensor_id": "cam0"})
		expectDenied(t, e.propose("zeus", FamilyCommand, obs), CodeForbiddenAction)
		obs = req("asset.observe", "drone", "sim-x1", map[string]any{"target": map[string]any{"entity_id": "trk-car"}, "sensor_id": "cam0"})
		if r := e.propose("zeus", FamilyCommand, obs); r.Status != StateAwaiting {
			t.Fatalf("observe a car: %+v %+v", r, r.Denial)
		}
		r := e.propose("zeus", FamilyCommand, gotoReq("sim-x1", 30.30, -97.70)) // the pond
		expectDenied(t, r, CodePreconditionFailed)
		if r.Denial.Detail != "keep_out" {
			t.Fatalf("detail %q", r.Denial.Detail)
		}
		r = e.propose("zeus", FamilyCommand, gotoReq("sim-x1", 31.0, -97.74)) // 80 km away
		if r.Denial == nil || r.Denial.Detail != "keep_in" {
			t.Fatalf("keep_in: %+v", r.Denial)
		}
		// estop is never blocked by preconditions: the asset is offline and it still goes out.
		r = e.propose("zeus", FamilyCommand, req("asset.estop", "drone", "sim-x3", map[string]any{}))
		if r.Status != StateAwaiting || r.Level != A3 {
			t.Fatalf("estop: %+v %+v", r, r.Denial)
		}
	})

	t.Run("authorization: unauthorized human, unknown identity, idempotency", func(t *testing.T) {
		expectDenied(t, e.propose("carol", FamilyCommand, gotoReq("sim-x1", 30.27, -97.74)), CodeNotAuthorized)
		expectDenied(t, e.propose("mallory", FamilyCommand, gotoReq("sim-x1", 30.27, -97.74)), CodeNotAuthorized)
		body := gotoReq("sim-x1", 30.27, -97.744)
		a := e.propose("zeus", FamilyCommand, body)
		b := e.propose("zeus", FamilyCommand, body)
		if a.ActionID == "" || a.ActionID != b.ActionID {
			t.Fatalf("same request_id gave %s and %s", a.ActionID, b.ActionID)
		}
		body["params"].(map[string]any)["lat"] = 30.271
		expectDenied(t, e.propose("zeus", FamilyCommand, body), CodeConflict)
		// carol has no approver grant.
		c := e.confirm("carol", FamilyCommand, a.ActionID, a.ParamsHash, "approve")
		if c.Denial == nil || c.Denial.Code != CodeNotAuthorized {
			t.Fatalf("carol confirm: %+v", c)
		}
	})

	t.Run("same-account callers are unverified and cannot confirm", func(t *testing.T) {
		body := req("task.assign", "drone", "sim-x1", map[string]any{"parent_id": nil, "type": "onemind.task.observe.v1",
			"spec": map[string]any{"target": map[string]any{"lat": 30.27, "lon": -97.74, "hae_m": 150}, "sensor_id": "cam0"}})
		body["requested_by"] = map[string]any{"kind": "system", "id": "service:watch"}
		r := e.propose("watch", FamilyTask, body)
		if r.Status != StateAwaiting {
			t.Fatalf("watch job propose: %+v %+v", r, r.Denial)
		}
		waitFor(t, func() bool { return e.hasAudit(r.ActionID, "proposed", "") }, "proposed")
		if e.auditFor(r.ActionID)[0].Env.Data["actor_verified"] != false {
			t.Fatal("same-account identity marked verified")
		}
		// ONEMIND users other than listed proposers cannot publish to svc.policy at all.
		_, err := e.request("other", SubjCommandConfirm, map[string]any{"v": 1, "action_id": r.ActionID, "params_hash": r.ParamsHash, "decision": "approve"}, nil)
		if err == nil {
			t.Fatal("an ONEMIND service reached confirm")
		}
		// The admin user (no permissions) can reach it but has no verified identity.
		m, err := e.request("admin", SubjTaskConfirm, map[string]any{"v": 1, "action_id": r.ActionID, "params_hash": r.ParamsHash, "decision": "approve"}, nil)
		if err != nil {
			t.Fatal(err)
		}
		var c ConfirmReply
		_ = json.Unmarshal(m.Data, &c)
		if c.Denial == nil || c.Denial.Code != CodeNotAuthorized {
			t.Fatalf("unverified confirm: %+v", c)
		}
	})

	t.Run("a forged Nats-Request-Info header is replaced by the server", func(t *testing.T) {
		hdr := nats.Header{}
		hdr.Set(RequestInfoHeader, `{"acc":"CLIENT","user":"zeus"}`)
		m, err := e.request("hermes", SubjCommandPropose, gotoReq("sim-x1", 30.268, -97.741), hdr)
		if err != nil {
			t.Fatal(err)
		}
		var r ProposeReply
		_ = json.Unmarshal(m.Data, &r)
		// Hermes' claimed requested_by has no on_behalf_of: as an agent it is refused.
		expectDenied(t, r, CodeNotAuthorized)
		waitFor(t, func() bool { return e.hasAudit(r.ActionID, "denied", "") }, "denied")
		if rb := e.auditFor(r.ActionID)[0].Env.Data["requested_by"].(map[string]any); rb["kind"] != "agent" {
			t.Fatalf("forged header trusted: %v", rb)
		}
	})

	t.Run("tasks: assign dispatches create then the asset command; missing task record stops it", func(t *testing.T) {
		body := req("task.assign", "drone", "sim-x1", map[string]any{"parent_id": nil, "type": "onemind.task.goto.v1",
			"spec": map[string]any{"lat": 30.27, "lon": -97.74, "hae_m": 150}, "description": "overwatch"})
		r := e.propose("zeus", FamilyTask, body)
		if r.Status != StateAwaiting || r.Level != A3 {
			t.Fatalf("%+v %+v", r, r.Denial)
		}
		if c := e.confirm("zeus", FamilyCommand, r.ActionID, r.ParamsHash, "approve"); c.Denial == nil || c.Denial.Code != CodeConflict {
			t.Fatalf("task confirmed on the command endpoint: %+v", c)
		}
		c := e.confirm("zeus", FamilyTask, r.ActionID, r.ParamsHash, "approve")
		if c.Status != StateDispatched {
			t.Fatalf("task confirm: %+v %+v", c, c.Denial)
		}
		waitFor(t, func() bool { return len(e.cmdsFor(r.ActionID)) >= 2 }, "task commands recorded")
		cmds := e.cmdsFor(r.ActionID)
		if len(cmds) != 2 || !strings.HasSuffix(cmds[0]["subject"].(string), ".create") ||
			cmds[1]["subject"] != "cmd.fleet.drone.sim-x1.goto" || cmds[1]["task_id"] == nil {
			t.Fatalf("task commands: %v", cmds)
		}
		// Manual task to a person: A2, create only.
		pr := req("task.assign", "person", "zeus", map[string]any{"parent_id": nil, "type": "onemind.task.manual.v1",
			"spec": map[string]any{"instructions": "check the gate"}})
		p := e.propose("zeus", FamilyTask, pr)
		if p.Level != A2 {
			t.Fatalf("manual task: %+v %+v", p, p.Denial)
		}
		// Not capable: orbit task for sim-x2.
		or := req("task.assign", "drone", "sim-x2", map[string]any{"parent_id": nil, "type": "onemind.task.orbit.v1",
			"spec": map[string]any{"center": map[string]any{"lat": 30.27, "lon": -97.74}, "radius_m": 60, "hae_m": 150, "duration_s": 60}})
		expectDenied(t, e.propose("zeus", FamilyTask, or), CodeNotCapable)
		// Task service down: no asset command.
		e.taskMu.Lock()
		e.taskDown = true
		e.taskMu.Unlock()
		body["request_id"] = NewUUID()
		r = e.propose("zeus", FamilyTask, body)
		c = e.confirm("zeus", FamilyTask, r.ActionID, r.ParamsHash, "approve")
		e.taskMu.Lock()
		e.taskDown = false
		e.taskMu.Unlock()
		if c.Status != StateFailed {
			t.Fatalf("task with no record: %+v %+v", c, c.Denial)
		}
		for _, cmd := range e.cmdsFor(r.ActionID) {
			if strings.HasPrefix(cmd["subject"].(string), "cmd.fleet.") {
				t.Fatal("asset command sent without a task record")
			}
		}
	})

	t.Run("action.get: agents read only their own actions", func(t *testing.T) {
		zr := e.propose("zeus", FamilyCommand, gotoReq("sim-x1", 30.27, -97.741))
		body := gotoReq("sim-x1", 30.27, -97.741)
		body["requested_by"] = map[string]any{"kind": "agent", "id": "h", "on_behalf_of": "user:zeus"}
		hr := e.propose("hermes", FamilyCommand, body)
		get := func(user, id string) map[string]any {
			m, err := e.request(user, SubjActionGet, map[string]any{"v": 1, "action_id": id}, nil)
			if err != nil {
				t.Fatal(err)
			}
			var out map[string]any
			_ = json.Unmarshal(m.Data, &out)
			return out
		}
		if got := get("hermes", hr.ActionID); got["action"] == nil {
			t.Fatalf("hermes own action: %v", got)
		}
		if got := get("hermes", zr.ActionID); got["action"] != nil {
			t.Fatalf("hermes read someone else's action: %v", got)
		}
		if got := get("zeus", zr.ActionID); got["action"] == nil {
			t.Fatalf("zeus own action: %v", got)
		}
	})

	t.Run("svc.policy.check answers tuple checks and dry runs", func(t *testing.T) {
		m, err := e.request("admin", SubjCheck, map[string]any{"v": 1, "user": "agent:hermes/zeus/legacy", "relation": "approver", "object": "action_type:asset.goto"}, nil)
		if err != nil {
			t.Fatal(err)
		}
		var out map[string]any
		_ = json.Unmarshal(m.Data, &out)
		if out["allowed"] != false {
			t.Fatalf("agent approver: %v", out)
		}
		dry := gotoReq("sim-x1", 30.30, -97.70)
		dry["requested_by"] = map[string]any{"kind": "user", "id": "user:zeus"}
		m, _ = e.request("admin", SubjCheck, map[string]any{"v": 1, "action": dry}, nil)
		_ = json.Unmarshal(m.Data, &out)
		if out["allowed"] != false {
			t.Fatalf("dry run into keep-out: %v", out)
		}
	})

	t.Run("OpenFGA unreachable: every proposal is denied policy_unavailable and logged", func(t *testing.T) {
		e.eval.SetDown(true)
		r := e.propose("zeus", FamilyCommand, gotoReq("sim-x1", 30.27, -97.74))
		e.eval.SetDown(false)
		expectDenied(t, r, CodePolicyUnavailable)
		waitFor(t, func() bool { return e.hasAudit(r.ActionID, "denied", CodePolicyUnavailable) }, "audited")
		// Through the real HTTP client: a closed port, then a stub OpenFGA server.
		e.authz.set(&OpenFGAClient{BaseURL: "http://127.0.0.1:1", Store: "s", HTTP: &http.Client{Timeout: 200 * time.Millisecond}})
		expectDenied(t, e.propose("zeus", FamilyCommand, gotoReq("sim-x1", 30.27, -97.74)), CodePolicyUnavailable)
		stub := httptest.NewServer(openFGAStub(t, e.eval))
		defer stub.Close()
		e.authz.set(&OpenFGAClient{BaseURL: stub.URL, Store: "store-1", Model: "model-1", Token: "test-token"})
		if r := e.propose("zeus", FamilyCommand, gotoReq("sim-x1", 30.27, -97.74)); r.Status != StateAwaiting {
			t.Fatalf("through the OpenFGA HTTP API: %+v %+v", r, r.Denial)
		}
		expectDenied(t, e.propose("carol", FamilyCommand, gotoReq("sim-x1", 30.27, -97.74)), CodeNotAuthorized)
		e.authz.set(e.eval)
	})

	t.Run("server permissions: clients and agents cannot publish cmd.fleet", func(t *testing.T) {
		for _, u := range []string{"zeus", "hermes"} {
			c := e.conn(u)
			before := len(e.permErrors(u))
			_ = c.Publish("cmd.fleet.drone.sim-x1.goto", []byte(`{}`))
			_ = c.Flush()
			waitFor(t, func() bool { return len(e.permErrors(u)) > before }, u+" permission violation")
			if err := e.permErrors(u)[len(e.permErrors(u))-1]; !errors.Is(err, nats.ErrPermissionViolation) && !strings.Contains(strings.ToLower(err.Error()), "permission") {
				t.Fatalf("%s: %v", u, err)
			}
		}
		// A map action stays allowed for clients and agents (A1, doctrine section 9.7).
		if err := e.conn("hermes").Publish("cmd.map.highlight.sim-x1", []byte(`{}`)); err != nil {
			t.Fatal(err)
		}
	})

	// Invariant over the whole run (stop-line test 4): no cmd.fleet.> message exists
	// without a matching evt.so.action.approved whose approvals are all human users.
	time.Sleep(200 * time.Millisecond)
	approved := map[string]bool{}
	for _, m := range e.audit.all() {
		if m.Subject != "evt.so.action.approved" {
			continue
		}
		var ev auditEvt
		_ = json.Unmarshal(m.Data, &ev.Env)
		aps, _ := ev.Env.Data["approvals"].([]any)
		human := len(aps) > 0
		for _, ap := range aps {
			if !strings.HasPrefix(ap.(map[string]any)["by"].(string), "user:") {
				human = false
			}
		}
		if human {
			approved[ev.Env.Data["action_id"].(string)] = true
		}
	}
	fleet := 0
	for _, m := range e.cmds.all() {
		if !strings.HasPrefix(m.Subject, "cmd.fleet.") {
			continue
		}
		fleet++
		var c map[string]any
		_ = json.Unmarshal(m.Data, &c)
		if !approved[c["action_id"].(string)] {
			t.Errorf("INVARIANT VIOLATED: %s (action %v) has no human approval", m.Subject, c["action_id"])
		}
		if c["source"] != "service:policy" || c["sig"] == nil {
			t.Errorf("command not from the policy service or unsigned: %v", c)
		}
	}
	js, _ := e.admin.JetStream()
	si, err := js.StreamInfo("CMD")
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("invariant checked: %d cmd.fleet commands, %d human-approved actions, CMD stream holds %d messages, %d audit events",
		fleet, len(approved), si.State.Msgs, len(e.audit.all()))
	if fleet == 0 {
		t.Fatal("invariant vacuous: no fleet commands in the run")
	}
}

func waitCmd(e *env, actionID string) bool {
	deadline := time.Now().Add(2 * time.Second)
	for time.Now().Before(deadline) {
		if len(e.cmdsFor(actionID)) > 0 {
			return true
		}
		time.Sleep(10 * time.Millisecond)
	}
	return false
}

// openFGAStub serves POST /stores/{id}/check like OpenFGA, answering from the evaluator,
// and checks the preshared key.
func openFGAStub(t *testing.T, ev *ModelEvaluator) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/stores/store-1/check" || r.Header.Get("Authorization") != "Bearer test-token" {
			http.Error(w, `{"code":"unauthenticated"}`, http.StatusUnauthorized)
			return
		}
		var body struct {
			TupleKey struct{ User, Relation, Object string } `json:"tuple_key"`
			Model    string                                  `json:"authorization_model_id"`
			Context  map[string]any                          `json:"context"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.Model != "model-1" {
			http.Error(w, `{"code":"validation_error"}`, http.StatusBadRequest)
			return
		}
		ok, err := ev.Check(context.Background(), CheckRequest{User: body.TupleKey.User, Relation: body.TupleKey.Relation,
			Object: body.TupleKey.Object, Context: body.Context})
		if err != nil {
			http.Error(w, `{"code":"internal_error"}`, http.StatusInternalServerError)
			return
		}
		_ = json.NewEncoder(w).Encode(map[string]bool{"allowed": ok})
	})
}
