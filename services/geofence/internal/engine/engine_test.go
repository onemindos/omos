package engine

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/onemindos/omos/services/geofence/internal/zone"
)

type clock struct{ t time.Time }

func (c *clock) now() time.Time { return c.t }

func newEngine() (*Engine, *clock) {
	c := &clock{t: time.Date(2026, 10, 29, 12, 0, 0, 0, time.UTC)}
	return New(Config{Now: c.now, FeedSilence: 30 * time.Second}), c
}

// box is a closed ring around (lat, lon) with half-size h degrees.
func box(lat, lon, h float64) string {
	return fmt.Sprintf(`{"type":"Polygon","coordinates":[[[%v,%v],[%v,%v],[%v,%v],[%v,%v],[%v,%v]]]}`,
		lon-h, lat-h, lon+h, lat-h, lon+h, lat+h, lon-h, lat+h, lon-h, lat-h)
}

func mkZone(t *testing.T, id, kind string, coords string, extra string) *zone.Zone {
	t.Helper()
	v := fmt.Sprintf(`{"zone_id":%q,"active":true,"live":true,"watch":true,"geo":{"shape":"polygon","kind":%q,"coords":%s}%s}`, id, kind, coords, extra)
	z, err := zone.Parse(id, []byte(v), 1)
	if err != nil {
		t.Fatal(err)
	}
	return z
}

func drone(id string, lat, lon float64, seq uint64) Entity {
	return Entity{EntityID: id, EntityType: "drone", Callsign: strings.ToUpper(id), Lat: lat, Lon: lon, HasPos: true,
		Live: true, Lifecycle: "active", Seq: seq, T: int64(seq) * 1000, Stale: 60, Simulated: true}
}

func subjects(outs []Output) []string {
	var s []string
	for _, o := range outs {
		s = append(s, o.Subject)
	}
	return s
}

func want(t *testing.T, outs []Output, subj ...string) {
	t.Helper()
	got := subjects(outs)
	if strings.Join(got, " ") != strings.Join(subj, " ") {
		t.Fatalf("outputs %v, want %v", got, subj)
	}
}

func TestSimX1LeavesKeepIn(t *testing.T) {
	g, _ := newEngine()
	want(t, g.SetZone(mkZone(t, "home", zone.KindKeepIn, box(30.2672, -97.7431, 0.003), "")))
	outs := g.SetEntity(drone("sim-x1", 30.2672, -97.7431, 1))
	want(t, outs, "evt.map.geofence.home.enter")
	if c := outs[0].Data.(Crossing); c.Cause != CauseFirstFix && c.Cause != CauseCrossing {
		t.Fatalf("first enter cause %q", c.Cause)
	}
	want(t, g.SetEntity(drone("sim-x1", 30.2672, -97.7420, 2))) // still inside
	outs = g.SetEntity(drone("sim-x1", 30.2672, -97.7380, 3))   // about 500 m east: out
	want(t, outs, "evt.map.geofence.home.exit", "evt.so.alert.raised")
	a := outs[1].Data.(Alert)
	if a.Kind != KindGeofence || a.Level != LevelWarning || a.Reason != "left_keep_in" || a.About.EntityID != "sim-x1" ||
		a.ZoneID != "home" || a.State != StateOpen || a.AlertID != "alr-gf-home-sim-x1" || a.Seq != 1 || !a.Simulated {
		t.Fatalf("alert %+v", a)
	}
	if outs[1].MsgID != "alr-gf-home-sim-x1:1" {
		t.Fatalf("alert msg id %q", outs[1].MsgID)
	}
	want(t, g.SetEntity(drone("sim-x1", 30.2672, -97.7370, 4))) // still outside: nothing new
	outs = g.SetEntity(drone("sim-x1", 30.2672, -97.7431, 5))
	want(t, outs, "evt.map.geofence.home.enter", "evt.so.alert.resolved")
	if r := outs[1].Data.(Alert); r.Reason != "returned" || r.State != StateResolved || r.Seq != 2 || r.ResolvedT == nil {
		t.Fatalf("resolved %+v", r)
	}
	if len(g.OpenAlerts()) != 0 {
		t.Fatal("alert still open")
	}
}

func TestKeepOutHazardAOIAndPeople(t *testing.T) {
	g, _ := newEngine()
	g.SetZone(mkZone(t, "barn", zone.KindKeepOut, box(30.27, -97.74, 0.001), ""))
	g.SetZone(mkZone(t, "pond", zone.KindHazard, box(30.27, -97.74, 0.002), ""))
	g.SetZone(mkZone(t, "farm", zone.KindAOI, box(30.27, -97.74, 0.01), ""))
	ac := Entity{EntityID: "icao-a1b2c3", EntityType: "aircraft", Lat: 30.27, Lon: -97.74, HasPos: true, Live: true, Lifecycle: "active", Seq: 1}
	outs := g.SetEntity(ac)
	want(t, outs, "evt.map.geofence.barn.enter", "evt.so.alert.raised", "evt.map.geofence.farm.enter",
		"evt.so.alert.raised", "evt.map.geofence.pond.enter", "evt.so.alert.raised")
	levels := []string{outs[1].Data.(Alert).Level, outs[3].Data.(Alert).Level, outs[5].Data.(Alert).Level}
	if strings.Join(levels, ",") != "warning,advisory,caution" {
		t.Fatalf("levels %v", levels)
	}
	// aircraft are not own assets: a keep-in without applies_to does not watch them
	g.SetZone(mkZone(t, "home", zone.KindKeepIn, box(10, 10, 0.01), ""))
	want(t, g.SetEntity(Entity{EntityID: "icao-ffffff", EntityType: "aircraft", Lat: 50, Lon: 50, HasPos: true, Live: true, Lifecycle: "active", Seq: 1}))
	// people are never watched, whatever the zone
	want(t, g.SetEntity(Entity{EntityID: "tak-ann", EntityType: "person", Lat: 30.27, Lon: -97.74, HasPos: true, Live: true, Lifecycle: "active", Seq: 1}))
	want(t, g.SetEntity(Entity{EntityID: "trk-9", EntityType: "track", Kind: "person", Lat: 30.27, Lon: -97.74, HasPos: true, Live: true, Lifecycle: "active", Seq: 1}))
}

func TestZoneChangesAndRemoval(t *testing.T) {
	g, _ := newEngine()
	g.SetEntity(drone("sim-x1", 30.2672, -97.7431, 1))
	// a keep-in drawn away from the drone: immediately outside, so the alert is raised
	outs := g.SetZone(mkZone(t, "field", zone.KindKeepIn, box(31, -97, 0.01), ""))
	want(t, outs, "evt.so.alert.raised")
	// redraw it around the drone: enter (cause zone_change) and resolve
	outs = g.SetZone(mkZone(t, "field", zone.KindKeepIn, box(30.2672, -97.7431, 0.01), ""))
	want(t, outs, "evt.map.geofence.field.enter", "evt.so.alert.resolved")
	if c := outs[0].Data.(Crossing); c.Cause != CauseZoneChange {
		t.Fatalf("cause %q", c.Cause)
	}
	// deactivate: exit, nothing to resolve
	outs = g.SetZone(mkZone(t, "field", zone.KindKeepIn, box(30.2672, -97.7431, 0.01), `,"lifecycle":"inactive"`))
	want(t, outs, "evt.map.geofence.field.exit")
	g.SetZone(mkZone(t, "barn", zone.KindKeepOut, box(30.2672, -97.7431, 0.01), ""))
	if len(g.OpenAlerts()) != 1 {
		t.Fatalf("keep-out alert missing: %v", g.OpenAlerts())
	}
	outs = g.DeleteZone("barn")
	want(t, outs, "evt.map.geofence.barn.exit", "evt.so.alert.resolved")
	if c := outs[0].Data.(Crossing); c.Cause != CauseZoneRemoved || c.ZoneKind != zone.KindKeepOut {
		t.Fatalf("crossing %+v", c)
	}
	if r := outs[1].Data.(Alert); r.Reason != "zone_removed" {
		t.Fatalf("reason %q", r.Reason)
	}
}

func TestAppliesToAndWatchTypes(t *testing.T) {
	g, _ := newEngine()
	g.SetZone(mkZone(t, "only-x2", zone.KindKeepOut, box(0, 0, 1), `,"properties":{"applies_to":["sim-x2"]}`))
	g.SetZone(mkZone(t, "vessels", zone.KindAOI, box(0, 0, 1), `,"watch_types":["vessel"]`))
	want(t, g.SetEntity(drone("sim-x1", 0, 0, 1)))
	want(t, g.SetEntity(drone("sim-x2", 0, 0, 1)), "evt.map.geofence.only-x2.enter", "evt.so.alert.raised")
	want(t, g.SetEntity(Entity{EntityID: "mmsi-1", EntityType: "vessel", Lat: 0, Lon: 0.5, HasPos: true, Live: true, Lifecycle: "active"}),
		"evt.map.geofence.vessels.enter", "evt.so.alert.raised")
}

func TestStaleAsset(t *testing.T) {
	g, c := newEngine()
	g.SetEntity(drone("sim-x1", 0, 0, 1))
	e := drone("sim-x1", 0, 0, 2)
	e.Lifecycle = "stale"
	outs := g.SetEntity(e)
	want(t, outs, "evt.so.alert.raised")
	if a := outs[0].Data.(Alert); a.Kind != KindStaleAsset || a.Level != LevelCaution || a.AlertID != "alr-stale-sim-x1" {
		t.Fatalf("stale alert %+v", a)
	}
	e.Lifecycle = "lost"
	outs = g.SetEntity(e)
	want(t, outs, "evt.so.alert.updated")
	if a := outs[0].Data.(Alert); a.Level != LevelWarning || a.Reason != "lost" {
		t.Fatalf("lost %+v", a)
	}
	want(t, g.SetEntity(drone("sim-x1", 0, 0, 3)), "evt.so.alert.resolved")
	// silence: no record for 1.5 x stale + 5 s while active
	c.t = c.t.Add(80 * time.Second)
	want(t, g.Tick())
	c.t = c.t.Add(20 * time.Second)
	outs = g.Tick()
	want(t, outs, "evt.so.alert.raised")
	if outs[0].Data.(Alert).Reason != "no_update" {
		t.Fatal("silence reason")
	}
	want(t, g.Tick())
	// expired by TTL
	outs = g.RemoveEntity("sim-x1")
	want(t, outs, "evt.so.alert.updated")
	if a := outs[0].Data.(Alert); a.Reason != "expired" || a.Level != LevelWarning {
		t.Fatalf("expired %+v", a)
	}
	// non-asset types never get stale alerts
	ac := Entity{EntityID: "icao-1", EntityType: "aircraft", HasPos: true, Lat: 1, Lon: 1, Live: true, Lifecycle: "stale"}
	want(t, g.SetEntity(ac))
}

func TestFeedOutage(t *testing.T) {
	g, c := newEngine()
	want(t, g.FeedStatus("nws", []byte(`{"source":"nws","stale":false,"lastError":null}`)))
	outs := g.FeedStatus("nws", []byte(`{"source":"nws","stale":true,"lastError":"HTTP 503"}`))
	want(t, outs, "evt.so.alert.raised")
	if a := outs[0].Data.(Alert); a.Kind != KindFeedOutage || a.About.Feed != "nws" || !strings.Contains(a.Message, "HTTP 503") || a.Reason != "stale" {
		t.Fatalf("feed alert %+v", a)
	}
	want(t, g.FeedStatus("nws", []byte(`{"stale":true}`)))
	want(t, g.FeedStatus("nws", []byte(`{"stale":false}`)), "evt.so.alert.resolved")
	c.t = c.t.Add(31 * time.Second)
	outs = g.Tick()
	want(t, outs, "evt.so.alert.raised")
	if outs[0].Data.(Alert).Reason != "silent" {
		t.Fatal("silent reason")
	}
	want(t, g.FeedStatus("nws", []byte(`{"stale":false}`)), "evt.so.alert.resolved")
	// legacy feed-intel status: poll 300 s, last poll 20 min ago -> stale
	old := c.t.Add(-20 * time.Minute).Unix()
	want(t, g.FeedStatus("noaa", []byte(fmt.Sprintf(`{"feed_id":"noaa","enabled":true,"poll_interval":300,"last_poll_ts":%d,"last_error":"timeout"}`, old))),
		"evt.so.alert.raised")
	want(t, g.FeedStatus("noaa", []byte(`{"feed_id":"noaa","enabled":false,"poll_interval":300}`)), "evt.so.alert.resolved")
	c.t = c.t.Add(10 * time.Minute)
	for _, o := range g.Tick() {
		if o.Data.(Alert).About.Feed == "noaa" {
			t.Fatal("disabled feed raised a silence alert")
		}
	}
	want(t, g.FeedStatus("x", []byte(`not json`)))
}

// TestRestoreIsIdempotent replays the facts of one engine into a fresh one: the same world
// produces nothing new, and a crossing after the restart is still detected.
func TestRestoreIsIdempotent(t *testing.T) {
	g, _ := newEngine()
	home := mkZone(t, "home", zone.KindKeepIn, box(30.2672, -97.7431, 0.003), "")
	var log []Output
	log = append(log, g.SetZone(home)...)
	log = append(log, g.SetEntity(drone("sim-x1", 30.2672, -97.7431, 1))...)
	log = append(log, g.SetEntity(drone("sim-x2", 30.2672, -97.7431, 1))...)
	log = append(log, g.SetEntity(drone("sim-x2", 30.2672, -97.7300, 2))...) // x2 leaves: alert
	log = append(log, g.FeedStatus("nws", []byte(`{"stale":true}`))...)
	want(t, log, "evt.map.geofence.home.enter", "evt.map.geofence.home.enter", "evt.map.geofence.home.exit",
		"evt.so.alert.raised", "evt.so.alert.raised")

	h, _ := newEngine()
	for _, o := range log {
		b, _ := json.Marshal(o.Data) // through JSON, as the bus does
		if strings.HasPrefix(o.Subject, "evt.map.geofence.") {
			var c Crossing
			_ = json.Unmarshal(b, &c)
			h.RestoreCrossing(c)
		} else {
			var a Alert
			_ = json.Unmarshal(b, &a)
			h.RestoreAlert(a)
		}
	}
	h.Restored()
	want(t, h.SetZone(home))
	want(t, h.ReconcileZones())
	want(t, h.SetEntity(drone("sim-x1", 30.2672, -97.7431, 1)))
	want(t, h.SetEntity(drone("sim-x2", 30.2672, -97.7300, 2)))
	want(t, h.ReconcileEntities())
	want(t, h.FeedStatus("nws", []byte(`{"stale":true}`)))
	// a new crossing after the restart: x1 leaves, x2 comes back
	outs := h.SetEntity(drone("sim-x1", 30.2672, -97.7300, 2))
	want(t, outs, "evt.map.geofence.home.exit", "evt.so.alert.raised")
	if outs[1].Data.(Alert).Seq != 1 {
		t.Fatal("seq")
	}
	outs = h.SetEntity(drone("sim-x2", 30.2672, -97.7431, 3))
	want(t, outs, "evt.map.geofence.home.enter", "evt.so.alert.resolved")
	if outs[1].Data.(Alert).Seq != 2 {
		t.Fatalf("restored seq not continued: %d", outs[1].Data.(Alert).Seq)
	}
}

func TestReconcileAfterDowntime(t *testing.T) {
	g, _ := newEngine()
	home := mkZone(t, "home", zone.KindKeepOut, box(0, 0, 1), "")
	var log []Output
	log = append(log, g.SetZone(home)...)
	log = append(log, g.SetEntity(drone("sim-x1", 0, 0, 1))...)
	st := drone("sim-x3", 5, 5, 1)
	st.Lifecycle = "stale"
	log = append(log, g.SetEntity(st)...)
	h, _ := newEngine()
	for _, o := range log {
		b, _ := json.Marshal(o.Data)
		if strings.HasPrefix(o.Subject, "evt.map.") {
			var c Crossing
			_ = json.Unmarshal(b, &c)
			h.RestoreCrossing(c)
		} else {
			var a Alert
			_ = json.Unmarshal(b, &a)
			h.RestoreAlert(a)
		}
	}
	// while down, the zone was deleted and sim-x3 expired
	outs := h.ReconcileZones()
	want(t, outs, "evt.map.geofence.home.exit", "evt.so.alert.resolved")
	outs = h.ReconcileEntities()
	want(t, outs, "evt.so.alert.updated")
	if a := outs[0].Data.(Alert); a.AlertID != "alr-stale-sim-x3" || a.Reason != "expired" {
		t.Fatalf("%+v", a)
	}
}

func TestParseEntity(t *testing.T) {
	e, err := ParseEntity([]byte(`{"entity_id":"sim-x1","entity_type":"drone","lat":30.1,"lon":-97.2,"hae":150,
		"live":true,"lifecycle":"active","seq":7,"t":5,"stale":60,"traits":["locatable"],"indicators":{"simulated":true}}`))
	if err != nil || !e.HasPos || e.HAE == nil || *e.HAE != 150 || !e.Simulated || e.Seq != 7 {
		t.Fatalf("%+v %v", e, err)
	}
	e, _ = ParseEntity([]byte(`{"entity_id":"team-a","entity_type":"team","lat":0,"lon":0,"live":true}`))
	if e.HasPos {
		t.Fatal("non-located record has a position")
	}
	if _, err := ParseEntity([]byte(`{"lat":1}`)); err == nil {
		t.Fatal("record without id accepted")
	}
}

func TestIDs(t *testing.T) {
	if GeofenceAlertID("home", "sim-x1") != "alr-gf-home-sim-x1" {
		t.Fatal(GeofenceAlertID("home", "sim-x1"))
	}
	long := GeofenceAlertID("home", "tak-ANDROID-0123456789abcdef0123456789abcdef0123456789")
	if !strings.HasPrefix(long, "alr-gf-home-h") || len(long) > 64 {
		t.Fatalf("long id %q", long)
	}
}
