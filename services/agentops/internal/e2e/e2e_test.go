// Package e2e runs agentops against the real entity service (omos PR #552) and the real
// policy service (omos PR #554) built from their branches (run rule 9: unmerged services run
// locally). The tests skip unless ENTITY_SERVICE_BIN / POLICY_SERVICE_BIN point at binaries;
// see services/agentops/README.md "Tests".
package e2e

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/nats-io/nats.go"

	"github.com/onemindos/omos/services/agentops/internal/agentproto"
	"github.com/onemindos/omos/services/agentops/internal/brief"
	"github.com/onemindos/omos/services/agentops/internal/bus"
	"github.com/onemindos/omos/services/agentops/internal/feeds"
	"github.com/onemindos/omos/services/agentops/internal/geo"
	"github.com/onemindos/omos/services/agentops/internal/model"
	"github.com/onemindos/omos/services/agentops/internal/natstest"
	"github.com/onemindos/omos/services/agentops/internal/presence"
	"github.com/onemindos/omos/services/agentops/internal/watch"
)

func freeAddr(t *testing.T) string {
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer l.Close()
	return l.Addr().String()
}

func runBin(t *testing.T, bin string, env ...string) {
	t.Helper()
	cmd := exec.Command(bin)
	cmd.Env = append(os.Environ(), env...)
	logf, _ := os.Create(filepath.Join(t.TempDir(), filepath.Base(bin)+".log"))
	cmd.Stdout, cmd.Stderr = logf, logf
	if err := cmd.Start(); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		_ = cmd.Process.Signal(os.Interrupt)
		done := make(chan struct{})
		go func() { _, _ = cmd.Process.Wait(); close(done) }()
		select {
		case <-done:
		case <-time.After(5 * time.Second):
			_ = cmd.Process.Kill()
		}
		if t.Failed() {
			b, _ := os.ReadFile(logf.Name())
			t.Logf("%s log:\n%s", filepath.Base(bin), b)
		}
	})
}

func fleetPos(nc *nats.Conn, class, id, callsign string, lat, lon float64) {
	b, _ := json.Marshal(map[string]any{"uid": id, "type": "fleet." + class, "lat": lat, "lon": lon, "hae": 120, "course": 90, "speed": 5,
		"callsign": callsign, "stale": 60, "t": time.Now().UnixMilli(), "meta": map[string]any{"battery_pct": 80, "mode": "AUTO"}})
	_ = nc.Publish("tel.fleet."+class+"."+id+".position", b)
}

// TestRealEntityServiceShowsAgentAndBriefCitesIt: Hermes' heartbeat becomes
// ent.agent.agent-hermes-zeus-legacy.state in the real entity service, and an offline brief
// built from the real svc.entity.query cites it and two simulated drones.
func TestRealEntityServiceShowsAgentAndBriefCitesIt(t *testing.T) {
	bin := os.Getenv("ENTITY_SERVICE_BIN")
	if bin == "" {
		t.Skip("ENTITY_SERVICE_BIN not set")
	}
	srv := natstest.Start(t)
	nc := srv.Connect(t)
	js := natstest.Streams(t, nc, true)
	runBin(t, bin, "NATS_URL="+srv.URL, "ENT_MODE=jetstream", "ENT_MSG_TTL=false", "LEASE_BUCKET=none", "START_LOOKBACK=0",
		"INPUTS=TEL=tel.fleet.>", "HTTP_ADDR="+freeAddr(t), "SERVE_RPC=true",
		`ENTITY_LINKS={"fleet_uid:agent-hermes-zeus-legacy":"agent-hermes-zeus-legacy"}`)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	hermes := agentproto.NewBeater(srv.Connect(t), "hermes", "zeus", "legacy", "Legacy", 2*time.Second)
	hermes.Set(agentproto.OneMind{Mode: "busy", Activity: "answering a TAK chat", CurrentTask: &agentproto.Task{TaskID: "tsk-01jbx7k2m9q4r8s6t0v2w4y6z8", State: "executing"}})
	go func() { _ = hermes.Run(ctx) }()
	pub, _ := bus.NewPublisher(nc, bus.ModeAuto)
	tr := presence.New(presence.Config{Site: presence.Anchor{Lat: 30.2672, Lon: -97.7431}}, nil)
	go func() { _ = tr.Run(ctx, nc, pub, "AGENTS") }()
	fleetPos(nc, "drone", "sim-x1", "SIM-X1", 30.27, -97.74)
	fleetPos(nc, "drone", "sim-x2", "SIM-X2", 30.26, -97.75)

	subj := "ent.agent.agent-hermes-zeus-legacy.state"
	natstest.Eventually(t, 15*time.Second, subj, func() bool { return natstest.Count(t, js, "ENT", subj) == 1 })
	rec := map[string]any{}
	_ = json.Unmarshal(natstest.Messages(t, js, "ENT", subj)[0], &rec)
	meta, _ := rec["meta"].(map[string]any)
	ct, _ := meta["current_task"].(map[string]any)
	if rec["entity_type"] != "agent" || rec["callsign"] != "LEGACY" || rec["type"] != "fleet.agent" || ct["task_id"] != "tsk-01jbx7k2m9q4r8s6t0v2w4y6z8" || meta["mode"] != "busy" {
		t.Fatalf("agent entity: %v", rec)
	}
	ind, _ := rec["indicators"].(map[string]any)
	natstest.Evidence(t, "real entity service published %s: entity_type=%v callsign=%v type=%v lifecycle=%v meta.mode=%v meta.activity=%v meta.current_task=%v egress_tak=%v current_task(field)=%v",
		subj, rec["entity_type"], rec["callsign"], rec["type"], rec["lifecycle"], meta["mode"], meta["activity"], ct, ind["egress_tak"], rec["current_task"])

	natstest.Eventually(t, 10*time.Second, "drones in ENT", func() bool { return natstest.Count(t, js, "ENT", "ent.drone.*.state") == 2 })
	board := feeds.NewBoard(nil, map[string]float64{"nws-alerts": 300}, nil)
	chain := model.Chain{Models: []model.Model{&model.Agent{NC: nc, Subject: "agents.prompt.hermes.zeus.legacy", Inactivity: time.Second},
		&model.OpenAICompat{BaseURL: "http://" + freeAddr(t) + "/v1", ModelName: "qwen2.5:7b-instruct"}, model.Stub{}}, Timeout: 3 * time.Second}
	b := brief.New(brief.Config{}, nc, pub, chain, board, nil, slog.New(slog.NewTextHandler(io.Discard, nil)))
	br, err := b.Generate(ctx, brief.AOI{ID: "home", Name: "Homestead (test AOI)", BBox: geo.BBox{-97.80, 30.20, -97.70, 30.30}})
	if err != nil {
		t.Fatal(err)
	}
	ids := []string{}
	for _, c := range br.Citations {
		ids = append(ids, c.EntityID)
	}
	natstest.Evidence(t, "offline brief from the real svc.entity.query: model=%s offline=%v cited=%v gaps=%d text:\n%s", br.Model.Used, br.Model.Offline, ids, len(br.Gaps), br.Text)
	if !br.Model.Offline || len(ids) < 3 || !strings.Contains(strings.Join(ids, ","), "agent-hermes-zeus-legacy") || !strings.Contains(br.Text, "nws-alerts") {
		t.Fatalf("brief %+v", br)
	}
}

const policyWatch = `
watches:
  - id: sim-keep-in
    owner: user:zeus
    level: warning
    rule: {kind: zone, entity_types: [drone], entity_ids: ["sim-*"], zone: zone-home, alert: outside}
    inputs: [ent.drone.*.state, "kv:geofences"]
    on_trigger:
      explain: true
      propose: {task_types: [onemind.task.rtl.v1]}
`

// TestRealPolicyKeepsWatchProposalWaiting: the watch runner's proposal reaches the real policy
// service, is audited on evt.so.action.proposed, waits for a human, and nothing reaches cmd.>.
func TestRealPolicyKeepsWatchProposalWaiting(t *testing.T) {
	bin, schemas, fga := os.Getenv("POLICY_SERVICE_BIN"), os.Getenv("POLICY_SCHEMA_DIR"), os.Getenv("POLICY_MODEL_FILE")
	if bin == "" || schemas == "" || fga == "" {
		t.Skip("POLICY_SERVICE_BIN, POLICY_SCHEMA_DIR and POLICY_MODEL_FILE are required")
	}
	srv := natstest.Start(t)
	nc := srv.Connect(t)
	js := natstest.Streams(t, nc, true)
	dir := t.TempDir()
	tuples := []map[string]string{
		{"user": "user:zeus", "relation": "operator", "object": "asset:sim-x1"},
		{"user": "user:zeus", "relation": "proposer", "object": "action_type:task.assign"},
		{"user": "user:zeus", "relation": "approver", "object": "action_type:task.assign"},
		{"user": "agent:onemind-watch/zeus/hub", "relation": "proposer", "object": "action_type:task.assign"},
	}
	tb, _ := json.Marshal(tuples)
	_ = os.WriteFile(filepath.Join(dir, "tuples.json"), tb, 0o600)
	_ = os.WriteFile(filepath.Join(dir, "principals.json"), []byte(`{"own_account":"$G","agent_accounts":["AGENTS"],"principals":[]}`), 0o600)
	// svc.entity.get stand-in with the entity.md §4.2 example record (the real entity service
	// does not fill `capabilities` yet; see the PR description).
	var gets atomic.Int64
	_, _ = srv.Connect(t).Subscribe("svc.entity.get", func(m *nats.Msg) {
		gets.Add(1)
		b, _ := json.Marshal(map[string]any{"entity_id": "sim-x1", "entity_type": "drone", "kind": "multirotor",
			"location": map[string]any{"lat": 30.30, "lon": -97.74, "hae_m": 150}, "health": map[string]any{"connection": "online"},
			"capabilities": []string{"onemind.task.goto.v1", "onemind.task.rtl.v1", "onemind.task.hold.v1"}, "indicators": map[string]any{"simulated": true}})
		_ = m.Respond(b)
	})
	runBin(t, bin, "NATS_URL="+srv.URL, "POLICY_AUTHZ=local", "POLICY_MODEL_FILE="+fga, "POLICY_TUPLES_FILE="+filepath.Join(dir, "tuples.json"),
		"POLICY_PRINCIPALS_FILE="+filepath.Join(dir, "principals.json"), "POLICY_SCHEMA_DIR="+schemas, "HEALTH_ADDR="+freeAddr(t))
	natstest.Eventually(t, 15*time.Second, "policy service up", func() bool {
		_, err := nc.Request("svc.policy.action.get", []byte(`{"v":1,"action_id":"act-none"}`), time.Second)
		return err == nil
	})
	var cmds atomic.Int64
	_, _ = srv.Connect(t).Subscribe("cmd.>", func(*nats.Msg) { cmds.Add(1) })

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	kv, _ := js.KeyValue(ctx, "geofences")
	_, _ = kv.Put(ctx, "zone-home", []byte(`{"entity_id":"zone-home","geo":{"shape":"polygon","kind":"keep_in","coords":{"type":"Polygon","coordinates":[[[-97.75,30.26],[-97.73,30.26],[-97.73,30.28],[-97.75,30.28],[-97.75,30.26]]]}}}`))
	ws, err := watch.Parse([]byte(policyWatch))
	if err != nil {
		t.Fatal(err)
	}
	pub, _ := bus.NewPublisher(nc, bus.ModeAuto)
	r := watch.NewRunner(watch.Config{AgentID: "agent:onemind-watch/zeus/hub"}, nc, pub, watch.NewEngine(ws, nil, nil),
		model.Chain{Models: []model.Model{model.Stub{}}}, nil, slog.New(slog.NewTextHandler(io.Discard, nil)))
	go func() { _ = r.Run(ctx, nil) }()
	time.Sleep(500 * time.Millisecond)
	put := func(lat float64, seq int) {
		b, _ := json.Marshal(map[string]any{"entity_id": "sim-x1", "entity_type": "drone", "callsign": "SIM-X1", "t": time.Now().UnixMilli(),
			"seq": seq, "live": true, "location": map[string]any{"lat": lat, "lon": -97.74}})
		_, _ = js.Publish(ctx, "ent.drone.sim-x1.state", b)
	}
	put(30.27, 1)
	time.Sleep(200 * time.Millisecond)
	put(30.30, 2)
	var final watch.Alert
	natstest.Eventually(t, 15*time.Second, "alert with proposal", func() bool {
		for _, b := range natstest.Messages(t, js, "EVT", "evt.entity.sim-x1.alert") {
			var a watch.Alert
			if json.Unmarshal(b, &a) == nil && a.Proposal != nil {
				final = a
				return true
			}
		}
		return false
	})
	p := final.Proposal
	natstest.Evidence(t, "real policy answered: action_id=%s status=%s level=%v approvals_needed=%d denial=%s", p.ActionID, p.Status, p.Level, p.Approvals, string(p.Denial))
	if p.Status != "awaiting_confirmation" || p.ActionID == "" || len(p.Denial) > 0 && string(p.Denial) != "null" {
		t.Fatalf("proposal %+v", p)
	}
	natstest.Eventually(t, 5*time.Second, "evt.so.action.proposed", func() bool { return natstest.Count(t, js, "EVT", "evt.so.action.proposed") >= 1 })
	audit := natstest.Messages(t, js, "EVT", "evt.so.action.proposed")
	natstest.Evidence(t, "audit evt.so.action.proposed: %s", truncate(string(audit[0]), 900))
	time.Sleep(time.Second)
	if cmds.Load() != 0 || natstest.Count(t, js, "CMD", "cmd.>") != 0 {
		t.Fatalf("a command went out without a human: %d", cmds.Load())
	}
	// The action is still waiting: svc.policy.action.get says so.
	raw, err := nc.Request("svc.policy.action.get", []byte(fmt.Sprintf(`{"v":1,"action_id":%q}`, p.ActionID)), 3*time.Second)
	if err != nil {
		t.Fatal(err)
	}
	natstest.Evidence(t, "svc.policy.action.get after 1 s: %s; cmd.> messages=0", truncate(string(raw.Data), 400))
	if !strings.Contains(string(raw.Data), "awaiting_confirmation") {
		t.Fatalf("action not waiting: %s", raw.Data)
	}
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "..."
}
