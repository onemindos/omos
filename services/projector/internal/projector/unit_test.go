package projector

import (
	"encoding/xml"
	"io"
	"log/slog"
	"math"
	"os"
	"regexp"
	"strings"
	"testing"
	"time"
)

var quiet = slog.New(slog.NewTextHandler(io.Discard, nil))

// typeRows is one fixture per row of cot-projection.md §4.1 and robotics-cot-profile.md §2.
var typeRows = []struct {
	name string
	m    map[string]any
	want string
}{
	{"drone multirotor", fixture("drone", "sim-x1", 30, -97, ownAsset, withKind("multirotor")), "a-f-A-M-H-Q"},
	{"drone fixed wing", fixture("drone", "sim-x2", 30, -97, ownAsset, withKind("fixed_wing")), "a-f-A-M-F-Q"},
	{"drone vtol", fixture("drone", "sim-x3", 30, -97, ownAsset, withKind("vtol")), "a-f-A-M-F-Q"},
	{"drone unknown", fixture("drone", "sim-x4", 30, -97, ownAsset), "a-f-A-M-F-Q"},
	{"robot ugv", fixture("robot", "sim-r1", 30, -97, ownAsset, withKind("ugv")), "a-f-G-E-V"},
	{"robot usv", fixture("robot", "sim-r2", 30, -97, ownAsset, set("environment", "surface")), "a-f-S-X"},
	{"robot uuv", fixture("robot", "sim-r3", 30, -97, ownAsset, set("environment", "subsurface")), "a-f-U-S-U"},
	{"own car", fixture("vehicle", "sim-v1", 30, -97, ownAsset, withKind("car")), "a-f-G-E-V-C"},
	{"own boat", fixture("vehicle", "sim-v2", 30, -97, ownAsset, set("environment", "surface")), "a-f-S-X"},
	{"aircraft fixed", fixture("aircraft", "icao-a1", 30, -97, withKind("fixed_wing")), "a-n-A-C-F"},
	{"aircraft rotary", fixture("aircraft", "icao-a2", 30, -97, withKind("rotary")), "a-n-A-C-H"},
	{"aircraft blimp", fixture("aircraft", "icao-a3", 30, -97, withKind("blimp")), "a-n-A-C-L"},
	{"aircraft unknown", fixture("aircraft", "icao-a4", 30, -97), "a-n-A"},
	{"vessel cargo", fixture("vessel", "mmsi-1", 29, -95, withKind("cargo")), "a-n-S-X-M-C"},
	{"vessel tanker", fixture("vessel", "mmsi-2", 29, -95, withKind("tanker")), "a-n-S-X-M-O"},
	{"vessel passenger", fixture("vessel", "mmsi-3", 29, -95, withKind("passenger")), "a-n-S-X-M-P"},
	{"vessel ferry", fixture("vessel", "mmsi-4", 29, -95, withKind("ferry")), "a-n-S-X-M-F"},
	{"vessel fishing", fixture("vessel", "mmsi-5", 29, -95, withKind("fishing")), "a-n-S-X-F"},
	{"vessel tug", fixture("vessel", "mmsi-6", 29, -95, withKind("tug")), "a-n-S-X-M-T-U"},
	{"vessel towing", fixture("vessel", "mmsi-7", 29, -95, withKind("towing")), "a-n-S-X-M-T-O"},
	{"vessel leisure", fixture("vessel", "mmsi-8", 29, -95, withKind("leisure")), "a-n-S-X-R"},
	{"vessel hazmat", fixture("vessel", "mmsi-9", 29, -95, withKind("hazmat")), "a-n-S-X-M-H"},
	{"vessel other", fixture("vessel", "mmsi-10", 29, -95), "a-n-S-X"},
	{"track air", fixture("track", "trk-1", 30, -97, set("environment", "air")), "a-u-A"},
	{"track surface", fixture("track", "trk-2", 30, -97, set("environment", "surface")), "a-u-S"},
	{"track subsurface", fixture("track", "trk-3", 30, -97, set("environment", "subsurface")), "a-u-U"},
	{"track land", fixture("track", "trk-4", 30, -97, set("environment", "land")), "a-u-G"},
	{"track none", fixture("track", "trk-5", 30, -97), "a-u-X"},
	{"sensor", fixture("sensor", "sen-1", 30, -97), "a-n-G-E-S"},
	{"sensor emplaced", fixture("sensor", "sen-2", 30, -97, withKind("emplaced")), "a-n-G-E-S-E"},
	{"camera", fixture("camera", "cam-1", 30, -97), "a-n-G-E-S-E"},
	{"structure", fixture("structure", "bld-1", 30, -97), "a-n-G-I-c"},
	{"place", fixture("place", "plc-1", 30, -97), "b-m-p-s-m"},
	{"zone polygon", fixture("zone", "z-1", 30, -97, zoneGeo("aoi", boxRing(30, -97, 0.1))), "u-d-f"},
	{"zone circle", fixture("zone", "z-2", 30, -97, set("geo", map[string]any{"shape": "circle", "kind": "hazard", "radius_m": 500,
		"coords": map[string]any{"type": "Point", "coordinates": []float64{-97, 30}}})), "u-d-c-c"},
	{"zone rectangle", fixture("zone", "z-3", 30, -97, zoneGeo("keep_in", boxRing(30, -97, 0.1)), func(m map[string]any) {
		m["geo"].(map[string]any)["shape"] = "rectangle"
	}), "u-d-r"},
	{"zone line", fixture("zone", "z-4", 30, -97, set("geo", map[string]any{"shape": "line", "kind": "route",
		"coords": map[string]any{"type": "LineString", "coordinates": [][]float64{{-97, 30}, {-97.1, 30.1}, {-97.2, 30.1}}}})), "b-m-r"},
	{"zone corridor", fixture("zone", "z-5", 30, -97, set("geo", map[string]any{"shape": "corridor", "kind": "keep_in", "radius_m": 100,
		"coords": map[string]any{"type": "LineString", "coordinates": [][]float64{{-97, 30}, {-97.1, 30.1}}}})), "u-d-f"},
	{"event earthquake", fixture("event", "usgs-1", 30, -97, withKind("earthquake")), "a-n-X-i-g-e"},
	{"event wildfire", fixture("event", "ev-2", 30, -97, withKind("wildfire")), "a-n-X-i-f-w"},
	{"event hotspot", fixture("event", "ev-3", 30, -97, withKind("fire_hotspot")), "a-n-X-i-f-h"},
	{"event flood", fixture("event", "ev-4", 30, -97, withKind("flood")), "a-n-X-i-m-f"},
	{"event tornado (NWS)", fixture("event", "ev-5", 30, -97, withKind("weather_alert"), set("meta", map[string]any{"event": "Tornado Warning"})), "a-n-X-i-m-t"},
	{"event hurricane", fixture("event", "ev-6", 30, -97, withKind("hurricane")), "a-n-X-i-m-c"},
	{"event thunderstorm", fixture("event", "ev-7", 30, -97, withKind("thunderstorm")), "a-n-X-i-m-u"},
	{"event tsunami", fixture("event", "ev-8", 30, -97, withKind("tsunami")), "a-n-X-i-m-n"},
	{"event winter", fixture("event", "ev-9", 30, -97, withKind("winter_storm")), "a-n-X-i-m-s"},
	{"event hail", fixture("event", "ev-10", 30, -97, withKind("hail")), "a-n-X-i-m-h"},
	{"event fog", fixture("event", "ev-11", 30, -97, withKind("fog")), "a-n-X-i-m-g"},
	{"event drought", fixture("event", "ev-12", 30, -97, withKind("drought")), "a-n-X-i-m-d"},
	{"event volcano", fixture("event", "ev-13", 30, -97, withKind("volcano")), "a-n-X-i-g-v"},
	{"event landslide", fixture("event", "ev-14", 30, -97, withKind("landslide")), "a-n-X-i-g-l"},
	{"event avalanche", fixture("event", "ev-15", 30, -97, withKind("avalanche")), "a-n-X-i-g-a"},
	{"event pollution", fixture("event", "ev-16", 30, -97, withKind("pollution")), "a-n-X-i-e"},
	{"event public health", fixture("event", "ev-17", 30, -97, withKind("public_health")), "a-n-X-i-h"},
	{"event transport", fixture("event", "ev-18", 30, -97, withKind("transport_accident")), "a-n-X-i-t"},
	{"event conflict (ACLED)", fixture("event", "acled-1", 30, -97, withKind("conflict_event")), "a-n-X-i-l"},
	{"event other", fixture("event", "ev-19", 30, -97, withKind("something")), "a-n-X-i"},
	{"satellite", fixture("satellite", "norad-25544", 30, -97), "a-n-P"},
}

// takTypes are the b-/u- TAK drawing types the contract lists (§4.1, §5.5).
var takTypes = map[string]bool{"b-m-p-s-m": true, "u-d-c-c": true, "u-d-f": true, "u-d-r": true, "b-m-r": true}

var forbiddenType = regexp.MustCompile(`^(a-[hjk]-|t-k-|t-e-|t-|a-.-G-E-W|a-.-U-W)`)

// cotTypesXML loads MITRE CoTtypes.xml from $COTTYPES_XML (the Apache-2.0 copy the contract
// cites); the check is skipped without it.
func cotTypesXML(t *testing.T) map[string]bool {
	p := os.Getenv("COTTYPES_XML")
	if p == "" {
		return nil
	}
	b, err := os.ReadFile(p)
	if err != nil {
		t.Fatal(err)
	}
	out := map[string]bool{}
	for _, m := range regexp.MustCompile(`cot="([^"]+)"`).FindAllSubmatch(b, -1) {
		out[string(m[1])] = true
	}
	return out
}

// Test 1 (cot-projection.md §12): every §4.1 row gives its type, every type is in CoTtypes.xml
// or a listed TAK drawing type, and none is hostile, strike, engage or a weapon.
func TestTypeTable(t *testing.T) {
	known := cotTypesXML(t)
	for _, r := range typeRows {
		e := decode(r.m)
		got, _, why := CotType(e)
		if why != "" || got != r.want {
			t.Errorf("%s: got %q (%s), want %q", r.name, got, why, r.want)
			continue
		}
		if forbiddenType.MatchString(got) {
			t.Errorf("%s: forbidden type %q", r.name, got)
		}
		if known != nil && !takTypes[got] && got != "a-n-P" { // a-n-P: §4.1 notes no a-.-P atom
			wild := got[:2] + "." + got[3:]
			if !known[wild] {
				t.Errorf("%s: %q (%s) not in CoTtypes.xml", r.name, got, wild)
			}
		}
		b, err := Build(e, got, time.Now(), buildOpts{target: TargetServer, interval: 5 * time.Second, maxVertices: 200})
		if err != nil {
			t.Errorf("%s: build: %v", r.name, err)
			continue
		}
		var ev struct {
			Type string `xml:"type,attr"`
			UID  string `xml:"uid,attr"`
		}
		if err := xml.Unmarshal(b.XML[strings.Index(string(b.XML), "<event"):], &ev); err != nil || ev.Type != got {
			t.Errorf("%s: xml %v type %q", r.name, err, ev.Type)
		}
	}
	if known == nil {
		t.Log("COTTYPES_XML not set: CoTtypes.xml membership not checked")
	} else {
		t.Logf("CoTtypes.xml membership checked for %d rows", len(typeRows))
	}
}

// The civilian check (§4.1 rule 4): no record can make the projector emit a hostile, joker,
// faker, strike, engage or weapon type, whatever its disposition or precomputed cot_type.
func TestCivilianCheck(t *testing.T) {
	asks := []string{"a-h-G", "a-h-A-M-F", "a-j-S", "a-k-G", "t-k-d", "t-e-x", "t-x-takp-v", "a-n-G-E-W-M", "a-u-U-W-T", "a-h-G-E-W", "y-x", "garbage"}
	dispositions := []string{"", "hostile", "joker", "faker", "friendly", "suspicious", "unknown", "neutral"}
	n := 0
	for _, typ := range []string{"track", "aircraft", "vessel", "sensor", "structure", "event", "camera"} {
		for _, d := range dispositions {
			for _, ask := range append(asks, "") {
				e := decode(fixture(typ, "x-1", 1, 1, set("disposition", d), set("cot_type", ask)))
				got, _, why := CotType(e)
				n++
				if why == "" && forbiddenType.MatchString(got) {
					t.Errorf("%s disposition=%q cot_type=%q -> %q", typ, d, ask, got)
				}
			}
		}
	}
	e := decode(fixture("track", "x-2", 1, 1, set("disposition", "hostile")))
	if got, down, _ := CotType(e); got != "a-s-X" || !down {
		t.Errorf("hostile track: got %q downgraded=%v, want a-s-X true", got, down)
	}
	t.Logf("%d combinations, none forbidden", n)
}

func TestNeverTypes(t *testing.T) {
	for _, typ := range []string{"person", "team", "agent", "node", "station", "alert"} {
		if got, _, why := CotType(decode(fixture(typ, "x-1", 1, 1))); why == "" {
			t.Errorf("%s projected as %q", typ, got)
		}
	}
	if got, _, why := CotType(decode(fixture("track", "x-1", 1, 1, withKind("person")))); why == "" {
		t.Errorf("person track projected as %q", got)
	}
	if got, _, why := CotType(decode(fixture("drone", "real-1", 1, 1, ownAsset))); why == "" {
		t.Errorf("non-sim drone projected as %q", got)
	}
	if _, _, why := CotType(decode(fixture("widget", "x-1", 1, 1))); why != NoTypeRow {
		t.Errorf("unknown entity_type: %q", why)
	}
}

func cfgFor(target Target) *Config { c := Defaults(target); return &c }

func TestMarkingGate(t *testing.T) {
	x := newAOIIndex()
	cases := []struct {
		tlp    string
		rel    []string
		groups []string
		want   bool
	}{
		{"CLEAR", nil, nil, true},
		{"GREEN", nil, nil, true},
		{"AMBER", []string{"homestead"}, nil, false},
		{"AMBER", []string{"homestead"}, []string{"homestead"}, true},
		{"AMBER+STRICT", []string{"homestead"}, nil, false},
		{"AMBER+STRICT", []string{"homestead"}, []string{"homestead"}, true},
		{"AMBER+STRICT", []string{"partner"}, []string{"homestead"}, false},
		{"RED", []string{"homestead"}, []string{"homestead"}, false},
		{"", nil, []string{"homestead"}, false},
		{"TOP", nil, nil, false},
	}
	for _, c := range cases {
		cfg := cfgFor(TargetServer)
		cfg.ReleaseGroups = c.groups
		e := decode(fixture("drone", "sim-x1", 1, 1, ownAsset, marking(c.tlp, c.rel...)))
		d := evaluate(e, cfg, x, true, time.Now())
		if d.ok != c.want {
			t.Errorf("tlp=%s rel=%v groups=%v: ok=%v (%s), want %v", c.tlp, c.rel, c.groups, d.ok, d.reason, c.want)
		}
	}
}

func TestGatesOrderAndReasons(t *testing.T) {
	x := newAOIIndex()
	x.setZone("z", geofence("z", "aoi", boxRing(10, 10, 1), true))
	now := time.Now()
	cfg := cfgFor(TargetServer)
	cases := []struct {
		name string
		m    map[string]any
		want string
	}{
		{"disabled", fixture("aircraft", "a", 10, 10), ReasonDisabled},
		{"not live", fixture("aircraft", "a", 10, 10, set("live", false)), ReasonNotAlive},
		{"expired", fixture("aircraft", "a", 10, 10, set("lifecycle", "expired")), ReasonNotAlive},
		{"not locatable", fixture("aircraft", "a", 10, 10, set("traits", []any{})), ReasonNotAlive},
		{"person", fixture("person", "p", 10, 10), ReasonPerson},
		{"person track", fixture("track", "t", 10, 10, withKind("person")), ReasonPerson},
		{"station", fixture("station", "s", 10, 10), ReasonPerson},
		{"no type", fixture("widget", "w", 10, 10), ReasonNoType},
		{"on tak", fixture("track", "t", 10, 10, set("sources", []any{map[string]any{"integration": "tak-mesh",
			"subject": "tak.cot.mesh.rx", "last_seen_t": now.UnixMilli()}})), ReasonOnTAK},
		{"egress off", fixture("aircraft", "a", 10, 10, ind("egress_tak", false)), ReasonEgressOff},
		{"red", fixture("aircraft", "a", 10, 10, marking("RED")), ReasonMarking},
		{"outside", fixture("aircraft", "a", 50, 50), ReasonNotRelevant},
		{"inside", fixture("aircraft", "a", 10, 10), ""},
		{"own outside", fixture("drone", "sim-x1", 50, 50, ownAsset), ""},
		{"starred outside", fixture("aircraft", "a", 50, 50, ind("starred", true)), ""},
		{"tasked outside", fixture("aircraft", "a", 50, 50, set("current_task", map[string]any{"task_id": "tsk-1", "state": "executing"})), ""},
		{"zone outside", fixture("zone", "z9", 50, 50, zoneGeo("hazard", boxRing(50, 50, 0.1))), ""},
		{"old tak source", fixture("track", "t", 10, 10, set("sources", []any{map[string]any{"integration": "tak",
			"subject": "tak.cot.out", "last_seen_t": now.Add(-time.Hour).UnixMilli()}})), ""},
	}
	for _, c := range cases {
		d := evaluate(decode(c.m), cfg, x, c.name != "disabled", now)
		if (c.want == "") != d.ok || d.reason != c.want {
			t.Errorf("%s: ok=%v reason=%q, want %q", c.name, d.ok, d.reason, c.want)
		}
	}
}

func TestAOIGeometry(t *testing.T) {
	x := newAOIIndex()
	alt := 100.0
	x.setZone("box", mustJSON(map[string]any{"zone_id": "box", "geo": map[string]any{"shape": "polygon", "kind": "aoi",
		"alt_min_m": 0.0, "alt_max_m": alt, "coords": map[string]any{"type": "Polygon", "coordinates": [][][]float64{{{0, 0}, {1, 0}, {1, 1}, {0, 1}, {0, 0}}}}}}))
	x.setZone("circle", mustJSON(map[string]any{"zone_id": "circle", "geo": map[string]any{"shape": "circle", "kind": "hazard", "radius_m": 1000,
		"coords": map[string]any{"type": "Point", "coordinates": []float64{20, 20}}}}))
	x.setZone("off", geofence("off", "aoi", boxRing(40, 40, 1), false))
	x.setZone("footprint", geofence("footprint", "footprint", boxRing(-40, -40, 1), true))
	in := func(lat, lon float64, mods ...func(map[string]any)) bool {
		ok, _ := x.inAOI(decode(fixture("aircraft", "a", lat, lon, mods...)))
		return ok
	}
	high := func(m map[string]any) { m["location"].(map[string]any)["hae_m"] = 500.0 }
	if !in(0.5, 0.5) || in(1.5, 0.5) || in(0.5, 0.5, high) {
		t.Error("polygon or altitude band wrong")
	}
	if !in(20.005, 20) || in(20.02, 20) {
		t.Error("circle wrong")
	}
	if in(40, 40) {
		t.Error("inactive zone used")
	}
	if in(-40, -40) {
		t.Error("footprint kind is not an AOI unless a mission names it")
	}
	x.setMission("msn-1", mustJSON(map[string]any{"mission_id": "msn-1", "state": "active", "aoi": "footprint"}))
	if !in(-40, -40) {
		t.Error("active mission aoi not used")
	}
	x.setMission("msn-1", mustJSON(map[string]any{"mission_id": "msn-1", "state": "completed", "aoi": "footprint"}))
	if in(-40, -40) {
		t.Error("completed mission aoi still used")
	}
	// A zone's own vertices count (§3.2 "any vertex of its geo").
	z := decode(fixture("event", "e", 5, 5, set("geo", map[string]any{"shape": "polygon", "kind": "hazard",
		"coords": map[string]any{"type": "Polygon", "coordinates": [][][]float64{{{5, 5}, {0.5, 0.5}, {5, 6}, {5, 5}}}}})))
	z.Location = nil
	z.Lat, z.Lon = nil, nil
	if ok, _ := x.inAOI(z); !ok {
		t.Error("geo vertex inside an AOI not counted")
	}
}

func TestMissionTasking(t *testing.T) {
	x := newAOIIndex()
	cfg := cfgFor(TargetServer)
	x.setMission("tsk-1", mustJSON(map[string]any{"task_id": "tsk-1", "rollup": "running",
		"assignee":         map[string]any{"entity_type": "drone", "entity_id": "sim-x1"},
		"target_snapshots": []any{map[string]any{"entity_id": "mmsi-1"}}}))
	if d := evaluate(decode(fixture("vessel", "mmsi-1", 50, 50)), cfg, x, true, time.Now()); !d.ok || d.rank != RankTasked {
		t.Errorf("targeted vessel: %+v", d)
	}
	x.setMission("tsk-1", mustJSON(map[string]any{"task_id": "tsk-1", "rollup": "completed",
		"target_snapshots": []any{map[string]any{"entity_id": "mmsi-1"}}}))
	if d := evaluate(decode(fixture("vessel", "mmsi-1", 50, 50)), cfg, x, true, time.Now()); d.ok {
		t.Errorf("completed task still targets: %+v", d)
	}
}

func TestMeshProfile(t *testing.T) {
	x := newAOIIndex()
	x.setZone("z", geofence("z", "aoi", boxRing(10, 10, 1), true))
	cfg := cfgFor(TargetMesh)
	if d := evaluate(decode(fixture("aircraft", "a", 10, 10)), cfg, x, true, time.Now()); d.ok {
		t.Error("mesh projected a world feed without MESH_WORLD_FEEDS")
	}
	cfg.MeshWorldFeeds = true
	if d := evaluate(decode(fixture("aircraft", "a", 10, 10)), cfg, x, true, time.Now()); !d.ok {
		t.Error("mesh with MESH_WORLD_FEEDS did not project an AOI entity")
	}
	e := decode(fixture("drone", "sim-x1", 10, 10, ownAsset))
	b, err := Build(e, "a-f-A-M-H-Q", time.Now(), buildOpts{target: TargetMesh, interval: 5 * time.Second, maxBytes: 1200, maxVertices: 200})
	if err != nil || strings.Contains(string(b.XML), "<onemind") {
		t.Fatalf("mesh build: %v %s", err, b.XML)
	}
	// A big zone is simplified to fit MESH_MAX_BYTES.
	ring := [][2]float64{}
	for i := 0; i <= 300; i++ {
		a := float64(i) / 300 * 6.283185307
		ring = append(ring, [2]float64{10 + 0.5*math.Cos(a), 10 + 0.5*math.Sin(a)})
	}
	z := decode(fixture("zone", "big", 10, 10, zoneGeo("aoi", ring)))
	b, err = Build(z, "u-d-f", time.Now(), buildOpts{target: TargetMesh, interval: 15 * time.Second, maxBytes: 1200, maxVertices: 200})
	if err != nil || len(b.XML) > 1200 {
		t.Fatalf("big zone on mesh: err=%v len=%d", err, len(b.XML))
	}
}

func TestCoTContent(t *testing.T) {
	e := decode(fixture("vessel", "mmsi-367000001", 29.71234, -95.012345, withKind("cargo"),
		set("callsign", `EXAMPLE "CARRIER" <&>`), set("stale", 120),
		func(m map[string]any) {
			l := m["location"].(map[string]any)
			l["course_deg"], l["speed_mps"] = 92.0, 6.1
			m["uncertainty"] = map[string]any{"ce_m": 10.0, "le_m": 0.0}
			m["provenance"].(map[string]any)["how"] = "m-p"
		}))
	now := time.Date(2026, 10, 29, 12, 0, 0, 0, time.UTC)
	b, err := Build(e, "a-n-S-X-M-C", now, buildOpts{target: TargetServer, interval: 5 * time.Second, maxVertices: 200})
	if err != nil {
		t.Fatal(err)
	}
	s := string(b.XML)
	for _, want := range []string{`uid="OM.vessel.mmsi-367000001"`, `type="a-n-S-X-M-C"`, `how="m-p"`,
		`time="2026-10-29T12:00:00.000Z"`, `stale="2026-10-29T12:02:00.000Z"`, `lat="29.712340"`, `lon="-95.012345"`,
		`ce="10.0"`, `le="9999999.0"`, `<track course="92.0" speed="6.10"/>`, `&lt;&amp;&gt;`, `omos-tak-projector=`,
		`<onemind v="1" id="mmsi-367000001" type="vessel"`, `tlp="CLEAR"`} {
		if !strings.Contains(s, want) {
			t.Errorf("missing %s in\n%s", want, s)
		}
	}
	for _, bad := range []string{"__group", "access=", "releasableTo", "caveat"} {
		if strings.Contains(s, bad) {
			t.Errorf("must not contain %s", bad)
		}
	}
	sim := decode(fixture("drone", "sim-x1", 1, 1, ownAsset))
	b, _ = Build(sim, "a-f-A-M-F-Q", now, buildOpts{target: TargetServer, interval: time.Second})
	if !strings.Contains(string(b.XML), `how="m-s"`) || !strings.Contains(string(b.XML), `opex="s"`) || !strings.Contains(string(b.XML), "SIMULATED") {
		t.Errorf("simulated asset not labelled (SL7): %s", b.XML)
	}
	r, _ := Build(sim, "a-f-A-M-F-Q", now, buildOpts{target: TargetServer, interval: time.Second, removal: true})
	if !strings.Contains(string(r.XML), `stale="2026-10-29T12:00:01.000Z"`) || !strings.Contains(string(r.XML), "<remarks>removed</remarks>") {
		t.Errorf("removal: %s", r.XML)
	}
}

func TestStaleClamp(t *testing.T) {
	e := decode(fixture("aircraft", "a", 1, 1, set("stale", 1)))
	if s := StaleSeconds(e, 5*time.Second); s != 15 {
		t.Errorf("lower clamp: %d", s)
	}
	e = decode(fixture("aircraft", "a", 1, 1, set("stale", 99999)))
	if s := StaleSeconds(e, 5*time.Second); s != 3600 {
		t.Errorf("upper clamp: %d", s)
	}
	z := decode(fixture("place", "p", 1, 1))
	if s := StaleSeconds(z, 5*time.Second); s != 86400 {
		t.Errorf("place default: %d", s)
	}
}

// Engine tests with a fake clock.

type sink struct{ msgs []Built }

func (s *sink) pub(b Built, _ *Entity) error { s.msgs = append(s.msgs, b); return nil }

func newTestEngine(cfg Config) (*Engine, *sink) {
	s := &sink{}
	return NewEngine(cfg, quiet, s.pub), s
}

func TestEngineHysteresisAndRemoval(t *testing.T) {
	cfg := Defaults(TargetServer)
	cfg.AOIExitGrace = 10 * time.Second
	g, s := newTestEngine(cfg)
	t0 := time.Unix(1_800_000_000, 0)
	g.SetZone("z", geofence("z", "aoi", boxRing(10, 10, 1), true), t0)
	g.Upsert(decode(fixture("aircraft", "a", 10, 10)), t0)
	g.Tick(t0)
	if len(s.msgs) != 1 {
		t.Fatalf("want 1 send, got %d", len(s.msgs))
	}
	// Leaves the AOI: still projected during the grace period, removed after it.
	g.Upsert(decode(fixture("aircraft", "a", 50, 50)), t0.Add(6*time.Second))
	for i := 6; i <= 15; i++ {
		g.Tick(t0.Add(time.Duration(i) * time.Second))
	}
	last := s.msgs[len(s.msgs)-1]
	if strings.Contains(string(last.XML), "removed") {
		t.Fatalf("removed before grace ended")
	}
	for i := 16; i <= 18; i++ {
		g.Tick(t0.Add(time.Duration(i) * time.Second))
	}
	last = s.msgs[len(s.msgs)-1]
	if !strings.Contains(string(last.XML), "<remarks>removed</remarks>") {
		t.Fatalf("no removal after grace: %s", last.XML)
	}
	n := len(s.msgs)
	for i := 19; i <= 400; i++ {
		g.Tick(t0.Add(time.Duration(i) * time.Second))
	}
	if len(s.msgs) != n {
		t.Fatalf("messages after removal: %d", len(s.msgs)-n)
	}
	// Egress turned off: removal at once, no grace.
	g.Upsert(decode(fixture("aircraft", "b", 10, 10)), t0.Add(500*time.Second))
	g.Tick(t0.Add(500 * time.Second))
	g.Upsert(decode(fixture("aircraft", "b", 10, 10, ind("egress_tak", false))), t0.Add(501*time.Second))
	g.Tick(t0.Add(501 * time.Second))
	if last = s.msgs[len(s.msgs)-1]; last.UID != "OM.aircraft.b" || !strings.Contains(string(last.XML), "removed") {
		t.Fatalf("egress off: no immediate removal")
	}
}

func TestEngineIntervalsAndCeiling(t *testing.T) {
	cfg := Defaults(TargetServer)
	cfg.MaxProjected = 10000
	g, s := newTestEngine(cfg)
	t0 := time.Unix(1_800_000_000, 0)
	g.SetZone("z", geofence("z", "aoi", boxRing(10, 10, 5), true), t0)
	// 10 Hz own drone for 10 s: at most one per second.
	for i := 0; i < 100; i++ {
		now := t0.Add(time.Duration(i) * 100 * time.Millisecond)
		g.Upsert(decode(fixture("drone", "sim-x1", 10+float64(i)*1e-4, 10, ownAsset)), now)
		g.Tick(now)
	}
	if n := len(s.msgs); n < 9 || n > 10 {
		t.Fatalf("own drone 10 Hz over 10 s: %d sends, want 9..10", n)
	}
	// 5,000 entities: never more than 50 in any one-second window.
	s.msgs = nil
	var times []time.Time
	g.pub = func(b Built, e *Entity) error { s.msgs = append(s.msgs, b); return nil }
	t1 := t0.Add(time.Minute)
	for i := 0; i < 5000; i++ {
		g.Upsert(decode(fixture("aircraft", idN("a", i), 10, 10)), t1)
	}
	for ms := int64(0); ms < 20000; ms += 20 {
		now := t1.Add(time.Duration(ms) * time.Millisecond)
		before := len(s.msgs)
		g.Tick(now)
		for j := before; j < len(s.msgs); j++ {
			times = append(times, now)
		}
	}
	maxWin := 0
	for i := range times {
		j := i
		for j < len(times) && times[j].Sub(times[i]) < time.Second {
			j++
		}
		if j-i > maxWin {
			maxWin = j - i
		}
	}
	if maxWin > 50 || len(times) < 900 {
		t.Fatalf("ceiling: max %d in one second (want <= 50), %d sent in 20 s", maxWin, len(times))
	}
	t.Logf("5,000 entities: %d sent in 20 s, max %d in any one-second window", len(times), maxWin)
}

func TestEngineCapOrder(t *testing.T) {
	cfg := Defaults(TargetServer)
	cfg.MaxProjected = 3
	g, s := newTestEngine(cfg)
	t0 := time.Unix(1_800_000_000, 0)
	g.SetZone("z", geofence("z", "aoi", boxRing(10, 10, 5), true), t0)
	g.Upsert(decode(fixture("aircraft", "far", 14, 14)), t0)
	g.Upsert(decode(fixture("aircraft", "near", 10.01, 10)), t0)
	g.Upsert(decode(fixture("aircraft", "mid", 12, 12)), t0)
	g.Upsert(decode(fixture("drone", "sim-x1", 50, 50, ownAsset)), t0)
	g.Upsert(decode(fixture("place", "home", 50, 50)), t0)
	for i := 0; i < 5; i++ {
		g.Tick(t0.Add(time.Duration(i) * time.Second))
	}
	live := map[string]bool{}
	for _, m := range s.msgs {
		live[m.UID] = !strings.Contains(string(m.XML), "removed")
	}
	want := map[string]bool{"OM.drone.sim-x1": true, "OM.place.home": true, "OM.aircraft.near": true}
	for uid, on := range live {
		if on != want[uid] {
			t.Errorf("%s live=%v want %v", uid, on, want[uid])
		}
	}
	if g.C.Get("cap_stopped") == 0 {
		t.Error("cap not counted")
	}
}

func TestDisableRemovesAndStops(t *testing.T) {
	g, s := newTestEngine(Defaults(TargetServer))
	t0 := time.Unix(1_800_000_000, 0)
	g.Upsert(decode(fixture("drone", "sim-x1", 1, 1, ownAsset)), t0)
	g.Tick(t0)
	g.SetEnabled(false, t0.Add(time.Second))
	g.Tick(t0.Add(time.Second))
	n := len(s.msgs)
	if n != 2 || !strings.Contains(string(s.msgs[1].XML), "removed") {
		t.Fatalf("disable: %d msgs", n)
	}
	for i := 2; i < 30; i++ {
		g.Upsert(decode(fixture("drone", "sim-x1", 1, 1, ownAsset)), t0.Add(time.Duration(i)*time.Second))
		g.Tick(t0.Add(time.Duration(i) * time.Second))
	}
	if len(s.msgs) != n {
		t.Fatal("sent while disabled")
	}
}

func TestSimplifyAndBuffer(t *testing.T) {
	var line []pt
	for i := 0; i < 1000; i++ {
		line = append(line, pt{float64(i) * 0.001, 0.0001 * float64(i%7)})
	}
	if out := simplify(line, 200); len(out) > 200 || len(out) < 2 {
		t.Errorf("simplify: %d", len(out))
	}
	poly := bufferLine([]pt{{0, 0}, {0.01, 0}}, 100)
	a := ringArea(poly)
	if !a.contains(0.0005, 0.005, 0, false) || a.contains(0.002, 0.005, 0, false) {
		t.Error("corridor buffer wrong")
	}
}
