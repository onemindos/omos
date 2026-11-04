package watch

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/onemindos/omos/services/agentops/internal/feeds"
)

type clock struct{ t time.Time }

func (c *clock) now() time.Time      { return c.t }
func (c *clock) add(d time.Duration) { c.t = c.t.Add(d) }
func zoneRec(id, kind string) []byte {
	return []byte(fmt.Sprintf(`{"entity_id":%q,"geo":{"shape":"polygon","kind":%q,"coords":{"type":"Polygon","coordinates":[[[-97.75,30.26],[-97.73,30.26],[-97.73,30.28],[-97.75,30.28],[-97.75,30.26]]]}}}`, id, kind))
}

func entRec(id string, lat, lon float64, t time.Time, extra string) []byte {
	return []byte(fmt.Sprintf(`{"entity_id":%q,"entity_type":"drone","callsign":"SIM-X1","lat":%v,"lon":%v,"t":%d,"seq":1,"live":true,"lifecycle":"active","location":{"lat":%v,"lon":%v}%s}`, id, lat, lon, t.UnixMilli(), lat, lon, extra))
}

func mustParse(t *testing.T, y string) []Watch {
	t.Helper()
	ws, err := Parse([]byte(y))
	if err != nil {
		t.Fatal(err)
	}
	return ws
}

func TestZoneEdgeCooldownAndGeofenceEvents(t *testing.T) {
	ck := &clock{time.Date(2026, 10, 29, 12, 0, 0, 0, time.UTC)}
	ws := mustParse(t, `
watches:
  - id: keep-in
    owner: user:zeus
    rule: {kind: zone, entity_types: [drone], entity_ids: ["sim-*"], zone: zone-home}
    inputs: [ent.drone.*.state, "kv:geofences"]
    rate: {cooldown_s: 60}
`)
	e := NewEngine(ws, nil, ck.now)
	if err := e.OnZone("zone-home", zoneRec("zone-home", "keep_in")); err != nil {
		t.Fatal(err)
	}
	subj := "ent.drone.sim-x1.state"
	trig, gev, _ := e.OnEntity(subj, entRec("sim-x1", 30.27, -97.74, ck.t, ""))
	if len(trig) != 0 || len(gev) != 0 {
		t.Fatalf("inside: %v %v", trig, gev)
	}
	ck.add(time.Second)
	trig, gev, _ = e.OnEntity(subj, entRec("sim-x1", 30.30, -97.74, ck.t, ""))
	if len(trig) != 1 || trig[0].Kind != "geofence_exit" || trig[0].AboutID != "sim-x1" || len(gev) != 1 || gev[0].Transition != "exit" {
		t.Fatalf("exit: %+v %+v", trig, gev)
	}
	// Still outside: no second alert.
	ck.add(time.Second)
	if trig, _, _ = e.OnEntity(subj, entRec("sim-x1", 30.31, -97.74, ck.t, "")); len(trig) != 0 {
		t.Fatal("repeated alert while still outside")
	}
	// Back in, out again inside the cooldown: crossing events, but no alert.
	ck.add(time.Second)
	_, gev, _ = e.OnEntity(subj, entRec("sim-x1", 30.27, -97.74, ck.t, ""))
	ck.add(time.Second)
	trig, gev2, _ := e.OnEntity(subj, entRec("sim-x1", 30.30, -97.74, ck.t, ""))
	if len(gev) != 1 || gev[0].Transition != "enter" || len(gev2) != 1 || len(trig) != 0 || e.Suppressed["cooldown"] != 1 {
		t.Fatalf("cooldown: %v %v %v %v", gev, gev2, trig, e.Suppressed)
	}
	// After the cooldown a new exit alerts again.
	ck.add(2 * time.Minute)
	_, _, _ = e.OnEntity(subj, entRec("sim-x1", 30.27, -97.74, ck.t, ""))
	ck.add(time.Second)
	if trig, _, _ = e.OnEntity(subj, entRec("sim-x1", 30.30, -97.74, ck.t, "")); len(trig) != 1 {
		t.Fatal("no alert after cooldown")
	}
	// Not selected: another drone id.
	if trig, _, _ = e.OnEntity("ent.drone.icao-abc.state", entRec("icao-abc", 30.30, -97.74, ck.t, "")); len(trig) != 0 {
		t.Fatal("unselected entity alerted")
	}
}

func TestMissingZoneIsReported(t *testing.T) {
	ws := mustParse(t, `
watches:
  - id: keep-in
    owner: user:zeus
    rule: {kind: zone, entity_types: [drone], zone: zone-gone}
    inputs: [ent.drone.*.state, "kv:geofences"]
`)
	e := NewEngine(ws, nil, nil)
	_ = e.OnZone("zone-gone", []byte(`{"entity_id":"zone-gone","geo":{"shape":"corridor","kind":"keep_in","coords":{"type":"LineString","coordinates":[]}}}`))
	if p := e.ZoneProblems(); !strings.Contains(p["zone-gone"], "not supported") {
		t.Fatalf("problems %v", p)
	}
}

func TestThresholdHysteresis(t *testing.T) {
	ck := &clock{time.Date(2026, 10, 29, 12, 0, 0, 0, time.UTC)}
	ws := mustParse(t, `
watches:
  - id: battery
    owner: user:zeus
    rule: {kind: threshold, entity_types: [drone], field: health.battery_pct, op: lt, value: 25, clear: 30}
    inputs: [ent.drone.*.state]
    rate: {cooldown_s: 10}
`)
	e := NewEngine(ws, nil, ck.now)
	b := func(p int) []Trigger {
		ck.add(time.Minute)
		tr, _, err := e.OnEntity("ent.drone.sim-x1.state", entRec("sim-x1", 30.27, -97.74, ck.t, fmt.Sprintf(`,"health":{"battery_pct":%d}`, p)))
		if err != nil {
			t.Fatal(err)
		}
		return tr
	}
	if len(b(40)) != 0 || len(b(24)) != 1 || len(b(27)) != 0 || len(b(24)) != 0 {
		t.Fatal("hysteresis: 24 after 27 must not re-alert (still below clear 30)")
	}
	if len(b(31)) != 0 || len(b(20)) != 1 {
		t.Fatal("after clearing above 30, a new drop must alert")
	}
}

func TestStaleAndLost(t *testing.T) {
	ck := &clock{time.Date(2026, 10, 29, 12, 0, 0, 0, time.UTC)}
	ws := mustParse(t, `
watches:
  - id: stale
    owner: user:zeus
    rule: {kind: stale, entity_types: [drone], max_age_s: 30}
    inputs: [ent.drone.*.state]
`)
	e := NewEngine(ws, nil, ck.now)
	_, _, _ = e.OnEntity("ent.drone.sim-x1.state", entRec("sim-x1", 30.27, -97.74, ck.t, ""))
	ck.add(20 * time.Second)
	if len(e.Sweep()) != 0 {
		t.Fatal("fresh entity alerted")
	}
	ck.add(20 * time.Second)
	tr := e.Sweep()
	if len(tr) != 1 || tr[0].Kind != "asset_stale" {
		t.Fatalf("stale: %+v", tr)
	}
	// A final live:false record is lost; it is a new condition only after a fresh record.
	ck.add(10 * time.Minute)
	_, _, _ = e.OnEntity("ent.drone.sim-x2.state", entRec("sim-x2", 30.27, -97.74, ck.t, ""))
	tr, _, _ = e.OnEntity("ent.drone.sim-x2.state", []byte(strings.Replace(string(entRec("sim-x2", 30.27, -97.74, ck.t, "")), `"live":true`, `"live":false`, 1)))
	if len(tr) != 1 || tr[0].Kind != "asset_lost" {
		t.Fatalf("lost: %+v", tr)
	}
}

func TestFeedStaleAndHourlyCap(t *testing.T) {
	ck := &clock{time.Date(2026, 10, 29, 12, 0, 0, 0, time.UTC)}
	board := feeds.NewBoard(ck.now, map[string]float64{"usgs-quakes": 300}, map[string][]string{"nws-alerts": {"weather"}})
	ws := mustParse(t, `
watches:
  - id: feeds
    owner: user:zeus
    rule: {kind: feed_stale, reasons: [stale, silent]}
    inputs: [tel.feed.*.status]
    rate: {cooldown_s: 10, max_alerts_per_hour: 2}
`)
	e := NewEngine(ws, board, ck.now)
	st := func(stale bool) []byte {
		b, _ := json.Marshal(map[string]any{"source": "nws-alerts", "stale": stale, "staleAfterSec": 120, "count": 3, "lastError": nil})
		return b
	}
	board.Update("tel.feed.nws-alerts.status", st(false))
	tr := e.Sweep()
	if len(tr) != 1 || tr[0].AboutID != "feed-usgs-quakes" || tr[0].Kind != "feed_silent" {
		t.Fatalf("expected the silent expected feed first: %+v", tr)
	}
	board.Update("tel.feed.nws-alerts.status", st(true))
	tr = e.Sweep()
	if len(tr) != 1 || tr[0].AboutID != "feed-nws-alerts" || tr[0].Kind != "feed_stale" {
		t.Fatalf("stale: %+v", tr)
	}
	// Hourly cap 2 is reached: clear and re-stale after the cooldown is suppressed.
	board.Update("tel.feed.nws-alerts.status", st(false))
	ck.add(time.Minute)
	e.Sweep()
	board.Update("tel.feed.nws-alerts.status", st(true))
	if tr = e.Sweep(); len(tr) != 0 || e.Suppressed["hourly_cap"] != 1 {
		t.Fatalf("cap: %+v %v", tr, e.Suppressed)
	}
}

func TestNewID(t *testing.T) {
	a, b := NewID("alr"), NewID("alr")
	if a == b || len(a) != 30 || strings.ToLower(a) != a {
		t.Fatalf("ids %s %s", a, b)
	}
}

// TestAlertIDDeterministic: the same condition in the same cooldown window (an ENT replay
// after a restart) gives the same alert id; the next window or another entity does not.
func TestAlertIDDeterministic(t *testing.T) {
	w := &Watch{ID: "w1", Rate: Rate{CooldownS: 300}}
	t0 := time.Date(2026, 10, 29, 12, 1, 0, 0, time.UTC)
	a := AlertIDFor(w.ID, "sim-x1", "threshold", cooldownBucket(w, t0))
	b := AlertIDFor(w.ID, "sim-x1", "threshold", cooldownBucket(w, t0.Add(3*time.Minute)))
	c := AlertIDFor(w.ID, "sim-x1", "threshold", cooldownBucket(w, t0.Add(5*time.Minute)))
	d := AlertIDFor(w.ID, "sim-x2", "threshold", cooldownBucket(w, t0))
	if a != b || a == c || a == d || len(a) != 30 || !strings.HasPrefix(a, "alr-") || strings.ToLower(a) != a {
		t.Fatalf("ids %s %s %s %s", a, b, c, d)
	}
}

// TestRestartReplayKeepsAlertID: a fresh engine (a pod restart) that gets the ENT
// last-per-subject replay of a drone still outside its zone raises the same alert id as
// the first engine did, so the proposal request_id repeats and the policy service dedups it.
func TestRestartReplayKeepsAlertID(t *testing.T) {
	ws := mustParse(t, `
watches:
  - id: keep-in
    owner: user:zeus
    rule: {kind: zone, entity_types: [drone], entity_ids: ["sim-*"], zone: zone-home}
    inputs: [ent.drone.*.state, "kv:geofences"]
    rate: {cooldown_s: 300}
`)
	first := func(at time.Time) string {
		ck := &clock{at}
		e := NewEngine(ws, nil, ck.now)
		_ = e.OnZone("zone-home", zoneRec("zone-home", "keep_in"))
		trig, _, _ := e.OnEntity("ent.drone.sim-x1.state", entRec("sim-x1", 30.30, -97.74, at, ""))
		if len(trig) != 1 {
			t.Fatalf("no alert: %+v", trig)
		}
		return trig[0].AlertID
	}
	t0 := time.Date(2026, 10, 29, 12, 0, 30, 0, time.UTC)
	a, b := first(t0), first(t0.Add(2*time.Minute))
	if a != b || requestID(a) != requestID(b) {
		t.Fatalf("restart changed alert id: %s %s", a, b)
	}
}
