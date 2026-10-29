package engine

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/onemindos/omos/services/entity/internal/adapt"
)

func fleetPos(class, uid, meta string, t int64) string {
	return fmt.Sprintf(`{"uid":%q,"type":"fleet.%s","lat":30.2672,"lon":-97.7431,"hae":150,"stale":10,"t":%d,"meta":{%s}}`,
		uid, class, t, meta)
}

// E2E bug 2: an own sim-* fleet asset carries the configured release group, so the
// projector's G7 gate can release it (entity.md §3 example `["homestead"]`).
func TestOwnSimFleetAssetReleasableTo(t *testing.T) {
	e := newEngine(&clock{t0 + 1000})
	r := e.Apply(obs(t, "tel.fleet.drone.sim-x1.position", fleetPos("drone", "sim-x1", `"kind":"multirotor","connection":"online"`, t0)))
	if r.Publish == nil {
		t.Fatalf("no publish: %+v", r)
	}
	m := r.Publish.Entity.Marking
	if m.TLP != "AMBER+STRICT" || !reflect.DeepEqual(m.ReleasableTo, []string{"homestead"}) {
		t.Fatalf("sim fleet marking = %+v, want AMBER+STRICT [homestead]", m)
	}
}

// The group is configurable, and an empty config releases to no group.
func TestOwnReleasableToConfigurable(t *testing.T) {
	for _, tc := range []struct {
		cfg  []string
		want []string
	}{{[]string{"homestead", "ops"}, []string{"homestead", "ops"}}, {nil, []string{}}} {
		c := DefaultConfig()
		c.OwnReleasableTo = tc.cfg
		e := New(c, (&clock{t0 + 1000}).now)
		e.Apply(obs(t, "tel.fleet.drone.sim-x1.position", fleetPos("drone", "sim-x1", "", t0)))
		got, _ := e.Get("sim-x1")
		if !reflect.DeepEqual(got.Marking.ReleasableTo, tc.want) {
			t.Fatalf("cfg %v: releasable_to = %v, want %v", tc.cfg, got.Marking.ReleasableTo, tc.want)
		}
	}
}

// Never widened: a non-sim fleet id, a feed, and a sim asset that a CoT source also reports
// all keep `releasable_to: []`.
func TestReleasableToNeverWidenedForFeeds(t *testing.T) {
	e := newEngine(&clock{t0 + 60000})
	e.Apply(obs(t, "tel.fleet.drone.x47.position", fleetPos("drone", "x47", "", t0+30000)))
	e.Apply(obs(t, "tel.fleet.drone.sim-x1.position", fleetPos("drone", "sim-x1", "", t0+30000)))
	if got, _ := e.Get("sim-x1"); len(got.Marking.ReleasableTo) != 1 {
		t.Fatalf("precondition: fleet-only sim-x1 releasable_to = %v", got.Marking.ReleasableTo)
	}
	e.Apply(obs(t, "tak.cot.mesh.rx", `<event version="2.0" uid="sim-x1" type="a-f-A-M-H-Q" how="m-s" `+
		`time="2026-10-29T12:00:37.000Z" start="2026-10-29T12:00:37.000Z" stale="2026-10-29T12:01:37.000Z">`+
		`<point lat="30.2672" lon="-97.7417" hae="152.4" ce="5" le="3"/><detail><contact callsign="SIM-X1"/></detail></event>`))
	for _, fx := range loadFixtures(t) {
		if o, _, err := adapt.Parse(fx.Subject, []byte(fx.Data)); err == nil && o != nil &&
			strings.HasPrefix(o.Integration, "feed-") {
			e.Apply(o)
		}
	}
	for _, id := range []string{"fleet-x47", "sim-x1", "icao-abc123", "norad-25544"} {
		got, ok := e.Get(id)
		if !ok {
			t.Fatalf("%s missing", id)
		}
		if len(got.Marking.ReleasableTo) != 0 {
			t.Fatalf("%s releasable_to = %v, want []", id, got.Marking.ReleasableTo)
		}
	}
}

// E2E bug 6: fleet entities carry cot_type from their class (robotics-cot-profile.md §2).
func TestFleetCotTypeFromClass(t *testing.T) {
	cases := []struct{ class, uid, meta, want string }{
		{"drone", "sim-x1", `"kind":"multirotor"`, "a-f-A-M-H-Q"},
		{"drone", "sim-fw1", `"kind":"fixed_wing"`, "a-f-A-M-F-Q"},
		{"drone", "sim-v1", `"kind":"vtol"`, "a-f-A-M-F-Q"},
		{"drone", "sim-u1", "", "a-f-A-M-H-Q"}, // adapter defaults a drone with no airframe to multirotor
		{"robot", "sim-tb1", `"kind":"ugv"`, "a-f-G-E-V"},
		{"robot", "sim-dog1", `"kind":"quadruped"`, "a-f-G-E-V"},
		{"robot", "sim-r1", "", "a-f-G-E-V"},
		{"robot", "sim-boat1", `"kind":"usv"`, "a-f-S-X"},
		{"robot", "sim-sub1", `"kind":"uuv"`, "a-f-U-S-U"},
		{"vehicle", "sim-car1", `"kind":"truck"`, "a-f-G-E-V-C"},
	}
	e := newEngine(&clock{t0 + 1000})
	for _, c := range cases {
		e.Apply(obs(t, "tel.fleet."+c.class+"."+c.uid+".position", fleetPos(c.class, c.uid, c.meta, t0)))
		got, ok := e.Get(c.uid)
		if !ok || got.CotType != c.want {
			t.Errorf("%s/%s: cot_type = %q, want %q", c.class, c.uid, got.CotType, c.want)
		}
		if got.Disposition != "friendly" {
			t.Errorf("%s: disposition %q", c.uid, got.Disposition)
		}
	}
}

// CAP review: health.connection goes "offline" once telemetry stops and the entity is stale,
// stays offline while lost, and comes back with the next report.
func TestConnectionOfflineWhenStale(t *testing.T) {
	c := &clock{t0 + 100}
	e := newEngine(c)
	e.Apply(obs(t, "tel.fleet.drone.sim-x1.position", fleetPos("drone", "sim-x1", `"battery_pct":90,"connection":"online"`, t0)))
	if got, _ := e.Get("sim-x1"); got.Health["connection"] != "online" {
		t.Fatalf("start: connection %v", got.Health["connection"])
	}
	c.ms = t0 + 11000 // stale is 10 s
	pubs := e.Sweep()
	got, _ := e.Get("sim-x1")
	if got.Lifecycle != "stale" || got.Health["connection"] != "offline" {
		t.Fatalf("stale: lifecycle %s connection %v", got.Lifecycle, got.Health["connection"])
	}
	if len(pubs) != 1 || pubs[0].Entity.Health["connection"] != "offline" {
		t.Fatalf("stale flip not published: %+v", pubs)
	}
	for _, p := range pubs {
		e.Published(p)
	}
	c.ms = t0 + 31000
	e.Sweep()
	if got, _ = e.Get("sim-x1"); got.Lifecycle != "lost" || got.Health["connection"] != "offline" {
		t.Fatalf("lost: lifecycle %s connection %v", got.Lifecycle, got.Health["connection"])
	}
	// A position without meta.connection restores the previous value.
	e.Apply(obs(t, "tel.fleet.drone.sim-x1.position", fleetPos("drone", "sim-x1", "", t0+31500)))
	if got, _ = e.Get("sim-x1"); got.Lifecycle != "active" || got.Health["connection"] != "online" {
		t.Fatalf("resumed: lifecycle %s connection %v", got.Lifecycle, got.Health["connection"])
	}
	// A bridge that reports its own link state wins.
	c.ms = t0 + 45000
	e.Sweep()
	e.Apply(obs(t, "tel.fleet.drone.sim-x1.status", fleetPos("drone", "sim-x1", `"connection":"degraded"`, t0+45000)))
	if got, _ = e.Get("sim-x1"); got.Health["connection"] != "degraded" {
		t.Fatalf("bridge state: connection %v", got.Health["connection"])
	}
}
