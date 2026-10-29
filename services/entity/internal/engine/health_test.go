package engine

import (
	"fmt"
	"testing"
	"time"
)

// meshCot is a tak-mesh CoT report about uid at epoch ms t (stale 60 s), with an optional
// battery in <status>.
func meshCot(uid string, t int64, battery string) string {
	ts := func(ms int64) string { return time.UnixMilli(ms).UTC().Format("2006-01-02T15:04:05.000Z") }
	status := ""
	if battery != "" {
		status = fmt.Sprintf(`<status battery=%q/>`, battery)
	}
	return fmt.Sprintf(`<event version="2.0" uid=%q type="a-f-A-M-H-Q" how="m-s" time=%q start=%q stale=%q>`+
		`<point lat="30.2672" lon="-97.7417" hae="152.4" ce="5" le="3"/><detail><contact callsign="SIM-X1"/>%s</detail></event>`,
		uid, ts(t), ts(t), ts(t+60000), status)
}

// F-ENT: the pre-stale connection is restored only by a fleet-* report; a CoT source (the
// tak-mesh gateway relaying the asset) resuming does not prove the bridge link is back.
func TestConnectionNotRestoredByCotSource(t *testing.T) {
	c := &clock{t0 + 100}
	e := newEngine(c)
	e.Apply(obs(t, "tel.fleet.drone.sim-x1.position", fleetPos("drone", "sim-x1", `"battery_pct":90,"connection":"online"`, t0)))
	c.ms = t0 + 11000 // fleet stale is 10 s
	for _, p := range e.Sweep() {
		e.Published(p)
	}
	if got, _ := e.Get("sim-x1"); got.Health["connection"] != "offline" {
		t.Fatalf("stale: connection %v", got.Health["connection"])
	}
	e.Apply(obs(t, "tak.cot.mesh.rx", meshCot("sim-x1", t0+11500, "")))
	got, _ := e.Get("sim-x1")
	if got.Lifecycle != "active" {
		t.Fatalf("mesh report: lifecycle %s", got.Lifecycle)
	}
	if got.Health["connection"] != "offline" {
		t.Fatalf("mesh report restored connection to %v", got.Health["connection"])
	}
	if got.Health["overall"] != "failed" {
		t.Fatalf("mesh report: overall %v, want failed (bridge offline)", got.Health["overall"])
	}
	// The next fleet report without meta.connection restores the pre-stale value.
	c.ms = t0 + 12000
	e.Apply(obs(t, "tel.fleet.drone.sim-x1.position", fleetPos("drone", "sim-x1", "", t0+12000)))
	if got, _ = e.Get("sim-x1"); got.Health["connection"] != "online" || got.Health["overall"] != "healthy" {
		t.Fatalf("fleet resumed: connection %v overall %v", got.Health["connection"], got.Health["overall"])
	}
}

// F-ENT: health.overall (entity.md §4.1) for a fleet entity from connection, battery and
// staleness; thresholds in policy.go.
func TestFleetHealthOverall(t *testing.T) {
	cases := []struct {
		name, meta, want string
	}{
		{"online full", `"battery_pct":90,"connection":"online"`, "healthy"},
		{"battery above warning", `"battery_pct":31`, "healthy"},
		{"battery at warning", `"battery_pct":30`, "warning"},
		{"battery above failed", `"battery_pct":16`, "warning"},
		{"battery at failed", `"battery_pct":15`, "failed"},
		{"battery empty", `"battery_pct":0,"connection":"online"`, "failed"},
		{"degraded link", `"battery_pct":90,"connection":"degraded"`, "warning"},
		{"offline link", `"battery_pct":90,"connection":"offline"`, "failed"},
		{"offline beats low battery warning", `"battery_pct":25,"connection":"offline"`, "failed"},
		{"connection only", `"connection":"online"`, "healthy"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			e := newEngine(&clock{t0 + 100})
			r := e.Apply(obs(t, "tel.fleet.drone.sim-x1.position", fleetPos("drone", "sim-x1", tc.meta, t0)))
			if r.Publish == nil {
				t.Fatalf("no publish: %+v", r)
			}
			if got := r.Publish.Entity.Health["overall"]; got != tc.want {
				t.Fatalf("overall %v, want %s (health %v)", got, tc.want, r.Publish.Entity.Health)
			}
		})
	}
}

func TestFleetHealthOverallStaleness(t *testing.T) {
	c := &clock{t0 + 100}
	e := newEngine(c)
	e.Apply(obs(t, "tel.fleet.drone.sim-x1.position", fleetPos("drone", "sim-x1", `"battery_pct":90,"connection":"online"`, t0)))
	if got, _ := e.Get("sim-x1"); got.Health["overall"] != "healthy" {
		t.Fatalf("fresh: overall %v", got.Health["overall"])
	}
	// Stale also forces connection offline, so overall is failed; the published flip carries it.
	c.ms = t0 + 11000
	pubs := e.Sweep()
	if len(pubs) != 1 || pubs[0].Entity.Health["overall"] != "failed" {
		t.Fatalf("stale flip: %+v", pubs)
	}
	for _, p := range pubs {
		e.Published(p)
	}
	c.ms = t0 + 31000
	e.Sweep()
	if got, _ := e.Get("sim-x1"); got.Lifecycle != "lost" || got.Health["overall"] != "failed" {
		t.Fatalf("lost: lifecycle %s overall %v", got.Lifecycle, got.Health["overall"])
	}
}

// Staleness alone (a fleet entity that never stated a connection is not possible via the geo
// adapter, so exercise the mapping directly).
func TestFleetOverallMapping(t *testing.T) {
	cases := []struct {
		h    map[string]any
		lc   string
		want string
	}{
		{map[string]any{"connection": "online"}, "active", "healthy"},
		{map[string]any{"connection": "online"}, "stale", "warning"},
		{map[string]any{"connection": "online"}, "lost", "failed"},
		{map[string]any{"connection": "online"}, "expired", "failed"},
		{map[string]any{"battery_pct": 50}, "active", "warning"}, // no link state reported
		{map[string]any{"connection": "online", "battery_pct": int64(10)}, "active", "failed"},
		{map[string]any{"connection": "online", "battery_pct": "low"}, "active", "healthy"}, // non-numeric ignored
	}
	for i, tc := range cases {
		if got := fleetOverall(tc.h, tc.lc); got != tc.want {
			t.Errorf("case %d %v/%s: got %s want %s", i, tc.h, tc.lc, got, tc.want)
		}
	}
}

// Entities with no fleet source get no derived overall.
func TestNoOverallWithoutFleetSource(t *testing.T) {
	e := newEngine(&clock{t0 + 100})
	e.Apply(obs(t, "tak.cot.mesh.rx", meshCot("sim-y2", t0, "40")))
	got, ok := e.Get("sim-y2")
	if !ok {
		t.Fatal("sim-y2 missing")
	}
	if _, has := got.Health["overall"]; has {
		t.Fatalf("CoT-only entity got overall: %v", got.Health)
	}
}
