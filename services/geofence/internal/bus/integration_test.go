package bus

import (
	"context"
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/nats-io/nats.go/jetstream"

	"github.com/onemindos/omos/services/geofence/internal/engine"
	"github.com/onemindos/omos/services/geofence/internal/sig"
	"github.com/onemindos/omos/services/geofence/internal/zone"
)

const lat0, lon0 = 30.2672, -97.7431 // sim-x1 start (services/entity testdata replay)

func createHome(t *testing.T, js jetstream.JetStream, pk policyKey) string {
	t.Helper()
	data := map[string]any{"name": "Home keep-in", "geo": map[string]any{"shape": "polygon", "kind": "keep_in",
		"coords": keepIn(lat0, lon0, 0.003)}, "applies_to": []string{"sim-x1", "sim-x2"}}
	id, b := policyCommand(t, "cmd.map.geofence.home.create", data, cmdOpt{key: &pk, approval: "approved"})
	publishCmd(t, js, "cmd.map.geofence.home.create", b, id)
	waitFor(t, "zone created", 5*time.Second, func() bool {
		f, ok := lastFact(t, js, "MISSION", "evt.map.geofence.created")
		return ok && f.Env.Causation == id
	})
	return id
}

// TestSimX1LeavesGeofenceRaisesAlert is the D6 acceptance test (FEATURES: "A sim-* drone
// leaving the geofence raises an alert") on a disposable server: a policy-signed command draws
// a keep-in zone, the entity records of sim-x1 move it out, and the watcher publishes the exit
// crossing and a warning alert.
func TestSimX1LeavesGeofenceRaisesAlert(t *testing.T) {
	s := startServer(t)
	_, js := connect(t, s)
	setupStreams(t, js)
	pk := newPolicyKey("policy-test")
	runService(t, s, sig.Keys{pk.kid: pk.pub}, nil)
	createHome(t, js, pk)

	kv, _ := js.KeyValue(context.Background(), "geofences")
	e, err := kv.Get(context.Background(), "home")
	if err != nil {
		t.Fatal(err)
	}
	var rec zone.Record
	if err := json.Unmarshal(e.Value(), &rec); err != nil || rec.ZoneID != "home" || rec.Geo.Kind != "keep_in" || !rec.Watch ||
		rec.UpdatedBy != "user:zeus" || !strings.HasPrefix(rec.ActionID, "act-") {
		t.Fatalf("KV record %+v %v", rec, err)
	}

	// sim-x1 flies east at 12.5 m/s from inside the 330 m keep-in to well outside.
	lon := lon0
	for seq := uint64(1); seq <= 12; seq++ {
		publishEnt(t, js, "drone", "sim-x1", entRecord("sim-x1", lat0, lon, seq, "active"), seq)
		lon += 0.0006
	}
	waitFor(t, "alert raised", 5*time.Second, func() bool { return countFacts(t, js, "EVT", "evt.so.alert.raised") == 1 })
	cross := facts(t, js, "MISSION", "evt.map.geofence.home.*")
	var events []string
	for _, f := range cross {
		events = append(events, f.Subject)
	}
	if strings.Join(events, " ") != "evt.map.geofence.home.enter evt.map.geofence.home.exit" {
		t.Fatalf("MISSION facts %v", events)
	}
	var c engine.Crossing
	_ = json.Unmarshal(cross[1].Env.Data, &c)
	if c.EntityID != "sim-x1" || c.Event != "exit" || c.Cause != engine.CauseCrossing || c.ZoneKind != "keep_in" || !c.Simulated {
		t.Fatalf("exit crossing %+v", c)
	}
	f, _ := lastFact(t, js, "EVT", "evt.so.alert.raised")
	var a engine.Alert
	_ = json.Unmarshal(f.Env.Data, &a)
	if a.Kind != engine.KindGeofence || a.Level != engine.LevelWarning || a.About.EntityID != "sim-x1" || a.ZoneID != "home" ||
		a.Reason != "left_keep_in" || a.State != "open" || f.Env.Type != "evt.so.alert.raised" || f.Env.Source != "omos-geofence" {
		t.Fatalf("alert %+v env %+v", a, f.Env)
	}
	t.Logf("alert: %s", f.Env.Data)
}

// TestStaleFeedRaisesAlert: a `tel.feed.{source}.status` saying stale raises a feed outage
// alert, a fresh one resolves it, and silence beyond FEED_SILENCE raises it again.
func TestStaleFeedRaisesAlert(t *testing.T) {
	s := startServer(t)
	nc, js := connect(t, s)
	setupStreams(t, js)
	pk := newPolicyKey("policy-test")
	runService(t, s, sig.Keys{pk.kid: pk.pub}, func(_ *Config, e *engine.Config) { e.FeedSilence = 1500 * time.Millisecond })

	pub := func(body string) {
		if err := nc.Publish("tel.feed.nws-alerts.status", []byte(body)); err != nil {
			t.Fatal(err)
		}
		_ = nc.Flush()
	}
	pub(`{"source":"nws-alerts","kind":"rest_poll","collector":"omos-ingest","stale":false,"lastError":null,"staleAfterSec":900}`)
	pub(`{"source":"nws-alerts","kind":"rest_poll","collector":"omos-ingest","stale":true,"lastError":"HTTP 503 from api.weather.gov","staleAfterSec":900}`)
	waitFor(t, "feed alert", 5*time.Second, func() bool { return countFacts(t, js, "EVT", "evt.so.alert.raised") == 1 })
	f, _ := lastFact(t, js, "EVT", "evt.so.alert.raised")
	var a engine.Alert
	_ = json.Unmarshal(f.Env.Data, &a)
	if a.Kind != engine.KindFeedOutage || a.About.Feed != "nws-alerts" || a.Reason != "stale" || !strings.Contains(a.Message, "HTTP 503") {
		t.Fatalf("feed alert %+v", a)
	}
	pub(`{"source":"nws-alerts","stale":false,"lastError":null}`)
	waitFor(t, "feed alert resolved", 5*time.Second, func() bool { return countFacts(t, js, "EVT", "evt.so.alert.resolved") == 1 })
	// now the collector goes quiet
	waitFor(t, "silence alert", 5*time.Second, func() bool { return countFacts(t, js, "EVT", "evt.so.alert.raised") == 2 })
	f, _ = lastFact(t, js, "EVT", "evt.so.alert.raised")
	_ = json.Unmarshal(f.Env.Data, &a)
	if a.Reason != "silent" || a.AlertID != "alr-feed-nws-alerts" || a.Seq != 3 {
		t.Fatalf("silence alert %+v", a)
	}
	t.Logf("silence alert: %s", f.Env.Data)
}

// TestRestartIsIdempotent: after a restart the service rebuilds memberships and open alerts
// from MISSION and EVT, publishes nothing for an unchanged world, applies a redelivered
// command only once, and still detects the next crossing.
func TestRestartIsIdempotent(t *testing.T) {
	s := startServer(t)
	nc, js := connect(t, s)
	setupStreams(t, js)
	pk := newPolicyKey("policy-test")
	keys := sig.Keys{pk.kid: pk.pub}
	mod := func(_ *Config, e *engine.Config) { e.FeedSilence = time.Hour }
	r1 := runService(t, s, keys, mod)
	cmdID := createHome(t, js, pk)
	publishEnt(t, js, "drone", "sim-x1", entRecord("sim-x1", lat0, lon0, 1, "active"), 1)
	publishEnt(t, js, "drone", "sim-x2", entRecord("sim-x2", lat0, lon0, 1, "active"), 1)
	publishEnt(t, js, "drone", "sim-x2", entRecord("sim-x2", lat0, lon0+0.01, 2, "active"), 2) // x2 leaves
	_ = nc.Publish("tel.feed.usgs.status", []byte(`{"source":"usgs","stale":true}`))
	_ = nc.Flush()
	waitFor(t, "two alerts", 5*time.Second, func() bool { return countFacts(t, js, "EVT", "evt.so.alert.raised") == 2 })
	time.Sleep(500 * time.Millisecond)
	missionBefore := countFacts(t, js, "MISSION", "evt.map.>")
	evtBefore := countFacts(t, js, "EVT", "evt.so.alert.>")
	if missionBefore != 4 { // created, enter x1, enter x2, exit x2
		t.Fatalf("MISSION before restart: %d", missionBefore)
	}
	r1.cancel()

	r2 := runService(t, s, keys, mod)
	if n := len(r2.svc.eng.OpenAlerts()); n != 2 {
		t.Fatalf("open alerts after restart: %d", n)
	}
	// the same feed status and the same command again (new Nats-Msg-Id, so CMD stores it)
	_ = nc.Publish("tel.feed.usgs.status", []byte(`{"source":"usgs","stale":true}`))
	_ = nc.Flush()
	data := map[string]any{"name": "Home keep-in", "geo": map[string]any{"shape": "polygon", "kind": "keep_in",
		"coords": keepIn(lat0, lon0, 0.003)}, "applies_to": []string{"sim-x1", "sim-x2"}}
	_, b := policyCommand(t, "cmd.map.geofence.home.create", data, cmdOpt{key: &pk, approval: "approved", id: cmdID})
	publishCmd(t, js, "cmd.map.geofence.home.create", b, cmdID+"-redelivered")
	time.Sleep(1500 * time.Millisecond)
	if n := countFacts(t, js, "MISSION", "evt.map.>"); n != missionBefore {
		t.Fatalf("restart published %d new MISSION facts", n-missionBefore)
	}
	if n := countFacts(t, js, "EVT", "evt.so.alert.>"); n != evtBefore {
		t.Fatalf("restart published %d new alert facts", n-evtBefore)
	}
	kv, _ := js.KeyValue(context.Background(), "geofences")
	if e, _ := kv.Get(context.Background(), "home"); e == nil || e.Revision() != 1 {
		t.Fatalf("redelivered create rewrote the zone: %+v", e)
	}
	// the next crossing is detected from the replayed membership: x1 leaves, x2 returns
	publishEnt(t, js, "drone", "sim-x1", entRecord("sim-x1", lat0, lon0+0.01, 2, "active"), 2)
	publishEnt(t, js, "drone", "sim-x2", entRecord("sim-x2", lat0, lon0, 3, "active"), 3)
	waitFor(t, "post-restart facts", 5*time.Second, func() bool {
		return countFacts(t, js, "MISSION", "evt.map.>") == missionBefore+2 && countFacts(t, js, "EVT", "evt.so.alert.>") == evtBefore+2
	})
	var subj []string
	for _, f := range facts(t, js, "EVT", "evt.so.alert.>")[evtBefore:] {
		var a engine.Alert
		_ = json.Unmarshal(f.Env.Data, &a)
		subj = append(subj, f.Subject+":"+a.About.EntityID)
	}
	if strings.Join(subj, " ") != "evt.so.alert.raised:sim-x1 evt.so.alert.resolved:sim-x2" {
		t.Fatalf("post-restart alerts %v", subj)
	}
}

func result(t *testing.T, js jetstream.JetStream, zoneID, cmdID string) (string, Result) {
	t.Helper()
	var subj string
	var r Result
	waitFor(t, "result of "+cmdID, 5*time.Second, func() bool {
		for _, f := range facts(t, js, "MISSION", "evt.map.geofence.*") {
			if f.Env.Causation == cmdID {
				subj = f.Subject
				_ = json.Unmarshal(f.Env.Data, &r)
				return true
			}
		}
		return false
	})
	return subj, r
}

// TestCommandGate: only policy-signed, approved, unexpired commands write KV `geofences`;
// create, update (compare-and-set) and delete behave as specified.
func TestCommandGate(t *testing.T) {
	s := startServer(t)
	_, js := connect(t, s)
	setupStreams(t, js)
	pk := newPolicyKey("policy-test")
	rogue := newPolicyKey("policy-test") // right kid, wrong key
	other := newPolicyKey("agent-hermes")
	runService(t, s, sig.Keys{pk.kid: pk.pub}, nil)
	kv, _ := js.KeyValue(context.Background(), "geofences")
	poly := map[string]any{"shape": "polygon", "kind": "keep_out", "coords": keepIn(lat0, lon0, 0.001)}
	send := func(zoneID, verb string, data any, o cmdOpt) Result {
		subj := "cmd.map.geofence." + zoneID + "." + verb
		id, b := policyCommand(t, subj, data, o)
		publishCmd(t, js, subj, b, id)
		_, r := result(t, js, zoneID, id)
		return r
	}
	ok := cmdOpt{key: &pk, approval: "approved"}
	cases := []struct {
		name, zone, verb string
		data             any
		o                cmdOpt
		code             string
	}{
		{"direct publish by a client or Hermes (unsigned)", "barn", "create", map[string]any{"geo": poly}, cmdOpt{approval: "approved"}, zone.CodeUnsigned},
		{"signed with a key the service does not trust", "barn", "create", map[string]any{"geo": poly}, cmdOpt{key: &other, approval: "approved"}, zone.CodeBadSignature},
		{"forged signature under the policy kid", "barn", "create", map[string]any{"geo": poly}, cmdOpt{key: &rogue, approval: "approved"}, zone.CodeBadSignature},
		{"expired", "barn", "create", map[string]any{"geo": poly}, cmdOpt{key: &pk, approval: "approved", expiresIn: -time.Minute}, zone.CodeExpired},
		{"no human approval (A2)", "barn", "create", map[string]any{"geo": poly}, cmdOpt{key: &pk, approval: "not_required"}, zone.CodeNotApproved},
		{"keep-out corridor the policy cannot evaluate", "barn", "create", map[string]any{"geo": map[string]any{"shape": "corridor", "kind": "keep_out",
			"coords": map[string]any{"type": "LineString", "coordinates": [][]float64{{0, 0}, {1, 1}}}, "radius_m": 50}}, ok, zone.CodeUnsupportedShape},
		{"watch people", "barn", "create", map[string]any{"geo": poly, "watch_types": []string{"person"}}, ok, zone.CodeForbidden},
		{"update of a missing zone", "barn", "update", map[string]any{"name": "x"}, ok, zone.CodeNotFound},
	}
	for _, c := range cases {
		r := send(c.zone, c.verb, c.data, c.o)
		if r.Outcome != "refused" || r.Code != c.code {
			t.Errorf("%s: got %s/%s (%s), want refused/%s", c.name, r.Outcome, r.Code, r.Message, c.code)
		}
	}
	if _, err := kv.Get(context.Background(), "barn"); err == nil {
		t.Fatal("a refused command wrote the zone")
	}

	r := send("barn", "create", map[string]any{"name": "Barn", "geo": poly}, ok)
	if r.Outcome != "created" || r.Revision == 0 {
		t.Fatalf("create %+v", r)
	}
	if r := send("barn", "create", map[string]any{"geo": poly}, ok); r.Code != zone.CodeExists {
		t.Fatalf("second create %+v", r)
	}
	stale := r.Revision
	r2 := send("barn", "update", map[string]any{"name": "Barn (north)", "revision": stale}, ok)
	if r2.Outcome != "updated" || r2.Revision <= stale {
		t.Fatalf("update %+v", r2)
	}
	if r := send("barn", "update", map[string]any{"name": "lost edit", "revision": stale}, ok); r.Code != zone.CodeConflict {
		t.Fatalf("stale edit %+v", r)
	}
	e, _ := kv.Get(context.Background(), "barn")
	var rec zone.Record
	_ = json.Unmarshal(e.Value(), &rec)
	if rec.Name != "Barn (north)" || rec.Geo.Kind != "keep_out" {
		t.Fatalf("record after update %+v", rec)
	}
	if r := send("barn", "delete", nil, ok); r.Outcome != "deleted" {
		t.Fatalf("delete %+v", r)
	}
	if _, err := kv.Get(context.Background(), "barn"); err == nil {
		t.Fatal("zone still present after delete")
	}
	if r := send("barn", "remove", nil, ok); r.Code != zone.CodeNotFound {
		t.Fatalf("second delete %+v", r)
	}
	// the doctrine §9.7 shorthand, through policy
	if r := send("north-field", "create", map[string]any{"polygon": [][]float64{{-97.75, 30.27}, {-97.74, 30.27}, {-97.74, 30.28}}, "kind": "aoi"}, ok); r.Outcome != "created" {
		t.Fatalf("shorthand %+v", r)
	}
}

// TestStaleAssetAndRPC: an asset record going stale raises a stale_asset alert, the alert is
// served on svc.geofence.alerts, and a fresh record resolves it.
func TestStaleAssetAndRPC(t *testing.T) {
	s := startServer(t)
	nc, js := connect(t, s)
	setupStreams(t, js)
	pk := newPolicyKey("policy-test")
	runService(t, s, sig.Keys{pk.kid: pk.pub}, nil)
	publishEnt(t, js, "drone", "sim-x1", entRecord("sim-x1", lat0, lon0, 1, "active"), 1)
	publishEnt(t, js, "drone", "sim-x1", entRecord("sim-x1", lat0, lon0, 2, "stale"), 2)
	waitFor(t, "stale alert", 5*time.Second, func() bool { return countFacts(t, js, "EVT", "evt.so.alert.raised") == 1 })
	m, err := nc.Request("svc.geofence.alerts", []byte(`{"kind":"stale_asset"}`), 2*time.Second)
	if err != nil {
		t.Fatal(err)
	}
	var reply struct {
		Alerts []engine.Alert `json:"alerts"`
	}
	if err := json.Unmarshal(m.Data, &reply); err != nil || len(reply.Alerts) != 1 || reply.Alerts[0].AlertID != "alr-stale-sim-x1" {
		t.Fatalf("rpc %s", m.Data)
	}
	publishEnt(t, js, "drone", "sim-x1", entRecord("sim-x1", lat0, lon0, 3, "active"), 3)
	waitFor(t, "stale resolved", 5*time.Second, func() bool { return countFacts(t, js, "EVT", "evt.so.alert.resolved") == 1 })
}
