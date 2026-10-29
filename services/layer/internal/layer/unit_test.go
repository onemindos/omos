package layer

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/nats-io/nats.go"
)

func TestDefaultLayersValid(t *testing.T) {
	if err := ValidateLayers(DefaultLayers()); err != nil {
		t.Fatal(err)
	}
	d := DefaultLayers()[0]
	if d.Filter() != "ent.*.*.state" || len(d.Filters()) != 5 || d.Filters()[0] != "ent.drone.>" {
		t.Fatalf("fleet filters %v %v", d.Filter(), d.Filters())
	}
	if (&Def{EntityTypes: []string{"aircraft"}}).Filter() != "ent.aircraft.>" {
		t.Fatal("single filter")
	}
	bad := []Def{{ID: "Bad", Tier: TierSmall, EntityTypes: []string{"x"}, Cap: 1, LimitS: 1},
		{ID: "l", Tier: TierLarge, EntityTypes: []string{"x"}, Cap: 1, LimitS: 1},
		{ID: "t", Tier: TierSmall, EntityTypes: []string{"a.b"}, Cap: 1, LimitS: 1}}
	for _, b := range bad {
		if b.Validate() == nil {
			t.Fatalf("%+v should fail", b)
		}
	}
	if ValidateLayers([]Def{DefaultLayers()[0], DefaultLayers()[0]}) == nil {
		t.Fatal("duplicate ids accepted")
	}
}

func TestCellZoom(t *testing.T) {
	d := &Def{}
	z := []int{2, 4, 6}
	for in, want := range map[int]int{0: 2, 3: 2, 4: 4, 5: 4, 6: 6, 9: 6} {
		if got := d.CellZoom(in, z); got != want {
			t.Fatalf("zoom %d -> %d, want %d", in, got, want)
		}
	}
}

func TestAreas(t *testing.T) {
	a, err := bboxArea([]float64{-74.3, 40.5, -73.6, 41.0})
	if err != nil || !a.contains(40.7, -74) || a.contains(42, -74) || a.contains(40.7, -75) {
		t.Fatal("bbox")
	}
	w, _ := bboxArea([]float64{170, -10, -170, 10}) // antimeridian
	if !w.contains(0, 179) || !w.contains(0, -179) || w.contains(0, 0) {
		t.Fatal("antimeridian bbox")
	}
	if _, err := bboxArea([]float64{0, 10, 1, 5}); err == nil {
		t.Fatal("inverted lat accepted")
	}
	p, _ := polyArea([][]float64{{0, 0}, {10, 0}, {0, 10}})
	if !p.contains(1, 1) || p.contains(9, 9) {
		t.Fatal("polygon")
	}
	c, err := zoneArea(&GeoIn{Coords: json.RawMessage(`{"type":"Point","coordinates":[-74,40.7]}`), RadiusM: 1000})
	if err != nil || !c.contains(40.7, -74) || !c.contains(40.705, -74) || c.contains(40.72, -74) {
		t.Fatalf("circle %v", err)
	}
	if _, err := zoneArea(&GeoIn{Coords: json.RawMessage(`{"type":"LineString","coordinates":[[0,0],[1,1]]}`)}); err == nil {
		t.Fatal("line accepted as area")
	}
	var nilArea *area
	if !nilArea.contains(89, 179) {
		t.Fatal("nil area")
	}
}

func TestTileXY(t *testing.T) {
	// z1: four quadrants.
	for _, c := range []struct {
		lat, lon float64
		x, y     int
	}{{45, -90, 0, 0}, {45, 90, 1, 0}, {-45, -90, 0, 1}, {-45, 90, 1, 1}} {
		if x, y := tileXY(c.lat, c.lon, 1); x != c.x || y != c.y {
			t.Fatalf("%v -> %d,%d", c, x, y)
		}
	}
	if x, y := tileXY(90, 180, 3); x != 7 || y != 0 {
		t.Fatalf("edge clamp %d %d", x, y)
	}
	g := newCellGrid(4)
	g.add(40.7, -74)
	g.add(40.8, -74.1)
	l := g.list()
	if len(l) != 1 || l[0][2] != 2 || l[0][3] != 40.75 {
		t.Fatalf("cells %v", l)
	}
}

func TestCanSee(t *testing.T) {
	d := &Directory{OwnAccount: "ONEMIND", OrgGroup: "onemind", Principals: []Principal{
		{Account: "CLIENT", User: "alice", Groups: []string{"team-a", "onemind"}},
		{Account: "CLIENT", User: "bob", Groups: []string{"team-b"}}}}
	msg := func(acc, user string) *nats.Msg {
		m := nats.NewMsg("svc.layer.snapshot")
		if acc != "" {
			m.Header.Set(RequestInfoHeader, fmt.Sprintf(`{"acc":%q,"user":%q}`, acc, user))
		}
		return m
	}
	svc := d.Resolve(msg("", ""))
	forged := d.Resolve(msg("ONEMIND", "x"))
	alice, bob, eve := d.Resolve(msg("CLIENT", "alice")), d.Resolve(msg("CLIENT", "bob")), d.Resolve(msg("CLIENT", "eve"))
	if !svc.Trusted || !forged.Trusted || alice.Trusted || eve.Known {
		t.Fatal("resolve")
	}
	cases := []struct {
		c    Caller
		tlp  string
		rel  []string
		want bool
	}{
		{svc, "RED", nil, true},
		{bob, "CLEAR", nil, true}, {bob, "GREEN", nil, true},
		{bob, "AMBER+STRICT", []string{"team-a"}, false}, {alice, "AMBER+STRICT", []string{"team-a"}, true},
		{bob, "AMBER", nil, false}, {alice, "AMBER", nil, true},
		{alice, "RED", nil, false}, {bob, "RED", []string{"team-b"}, true},
		{eve, "CLEAR", nil, false},
	}
	for i, c := range cases {
		if got := c.c.CanSee(c.tlp, c.rel); got != c.want {
			t.Fatalf("case %d (%s %v %v): %v", i, c.c.ID, c.tlp, c.rel, got)
		}
	}
}

func TestFreshness(t *testing.T) {
	now := time.Unix(1793275200, 0)
	d := &Def{LimitS: 60, Feeds: []string{"a", "b"}}
	mk := func(age time.Duration, stale bool) FeedStatus {
		return FeedStatus{FetchedAtMs: now.Add(-age).UnixMilli(), Count: 3, Stale: stale}
	}
	for _, c := range []struct {
		feeds map[string]FeedStatus
		down  bool
		want  string
	}{
		{map[string]FeedStatus{"a": mk(time.Second, false), "b": mk(time.Second, false)}, false, "fresh"},
		{map[string]FeedStatus{"a": mk(time.Second, false), "b": mk(2*time.Minute, false)}, false, "degraded"},
		{map[string]FeedStatus{"a": mk(time.Second, true)}, false, "stale"},
		{map[string]FeedStatus{"a": mk(time.Second, false), "b": mk(time.Second, false)}, true, "down"},
	} {
		if got := computeFreshness(d, c.feeds, c.down, now).Status; got != c.want {
			t.Fatalf("%v: %s, want %s", c, got, c.want)
		}
	}
	f := computeFreshness(&Def{LimitS: 60}, nil, false, now)
	if f.Status != "fresh" || len(f.Feeds) != 0 {
		t.Fatal("no feeds")
	}
	fs, ok := ParseFeedStatus("tel.feed.nws-alerts.status", []byte(`{"fetchedAt":"2026-10-29T12:00:00Z","count":12,"lastError":null,"stale":false,"staleAfterSec":900}`), now)
	if !ok || fs.Source != "nws-alerts" || fs.FetchedAtMs != time.Date(2026, 10, 29, 12, 0, 0, 0, time.UTC).UnixMilli() || fs.StaleAfterSec != 900 {
		t.Fatalf("parse %+v", fs)
	}
	if _, ok := ParseFeedStatus("tel.feed.x", nil, now); ok {
		t.Fatal("bad subject accepted")
	}
}

func TestChunkAll(t *testing.T) {
	item := []byte(`{"entity_id":"` + strings.Repeat("x", 300) + `"}`)
	items := make([][]byte, 4500)
	for i := range items {
		items[i] = item
	}
	ch := chunkAll("fleet", "entities", 5, 1, false, items)
	total := 0
	for i, c := range ch {
		var m struct {
			Part  int               `json:"part"`
			Last  bool              `json:"last"`
			Items []json.RawMessage `json:"items"`
		}
		if err := json.Unmarshal(c, &m); err != nil {
			t.Fatal(err)
		}
		if len(c) > MaxChunkBytes || len(m.Items) > MaxChunkItems || m.Part != i || m.Last != (i == len(ch)-1) {
			t.Fatalf("chunk %d: %d bytes %d items last %v", i, len(c), len(m.Items), m.Last)
		}
		total += len(m.Items)
	}
	if total != 4500 || len(ch) != 3 {
		t.Fatalf("total %d in %d chunks", total, len(ch))
	}
	big := make([][]byte, 300)
	for i := range big {
		big[i] = []byte(`"` + strings.Repeat("y", 5000) + `"`)
	}
	for _, c := range chunkAll("fleet", "entities", 5, 1, false, big) {
		if len(c) > MaxChunkBytes {
			t.Fatalf("byte limit: %d", len(c))
		}
	}
	if ch := chunkAll("fleet", "entities", 5, 1, false, nil); len(ch) != 1 || !strings.Contains(string(ch[0]), `"last":true`) {
		t.Fatal("empty snapshot must still send one last chunk")
	}
}

func TestEncodeCellsDropsOversizeZoom(t *testing.T) {
	cells := map[int][]cell{2: {{1, 1, 3, 1, 1}}, 6: make([]cell, 5000)}
	recs, kept := encodeCells("sea", cells, []int{2, 6}, 9, time.Unix(0, 0), 32*1024)
	if fmt.Sprint(kept) != "[2]" || recs[6] != nil || recs[2] == nil {
		t.Fatalf("kept %v", kept)
	}
}

func TestReplayWindow(t *testing.T) {
	now := time.Unix(1793275200, 0)
	clock := func() time.Time { return now }
	s := NewStore([]Def{fleetLayer()}, "hub", time.Minute, 3, clock)
	ent := func(id string, seq uint64) []byte {
		return []byte(fmt.Sprintf(`{"entity_id":%q,"entity_type":"drone","lat":1,"lon":2,"seq":%d,"live":true,"marking":{"tlp":"CLEAR"}}`, id, seq))
	}
	s.Apply("ent.drone.a.state", ent("a", 1), "", 10)
	s.MarkReady()
	if !s.replayValidLocked(11) || s.replayValidLocked(10) || s.replayValidLocked(12) {
		t.Fatal("initial window")
	}
	for i := uint64(11); i <= 14; i++ {
		s.Apply("ent.drone.b.state", ent("b", i), "", i)
	}
	// ring keeps 3 (12, 13, 14): floor 11
	if s.replayValidLocked(11) || !s.replayValidLocked(12) || !s.replayValidLocked(15) {
		t.Fatalf("after eviction floor %d", s.floor)
	}
	// A sweep at seq 14 then evicted: a client at 14 would miss it.
	exp := now.Add(-time.Second).UnixMilli()
	s.Apply("ent.drone.c.state", []byte(fmt.Sprintf(`{"entity_id":"c","entity_type":"drone","lat":1,"lon":2,"seq":1,"live":true,"expiry_t":%d}`, exp+10000)), "", 15)
	now = now.Add(11 * time.Second)
	s.Sweep()
	for i := 0; i < 3; i++ {
		s.Apply("ent.drone.d.state", ent("d", uint64(16+i)), "", uint64(16+i))
	}
	if s.floor != 15 || !s.floorStrict || s.replayValidLocked(16) || !s.replayValidLocked(17) {
		t.Fatalf("sweep floor %d strict %v", s.floor, s.floorStrict)
	}
	// time window eviction
	now = now.Add(2 * time.Minute)
	s.Sweep()
	if len(s.ring) != 0 || s.floor != 18 {
		t.Fatalf("time eviction: ring %d floor %d", len(s.ring), s.floor)
	}
}

func TestParseEntityMapping(t *testing.T) {
	e, in, err := parseEntity("drone", "sim-x1", []byte(`{"entity_id":"sim-x1","entity_type":"drone","location":{"lat":1.5,"lon":2.5,"hae_m":30,"course_deg":45,"speed_mps":7},"t":5,"stale":30,"seq":3,"live":true,"indicators":{"simulated":true},"marking":{"tlp":"AMBER+STRICT","releasable_to":["homestead"]},"sources":[{"data_type":"cot"}]}`), "site-mac")
	if err != nil || !in.Live {
		t.Fatal(err)
	}
	r := e.r
	if r.Lat != 1.5 || r.Lon != 2.5 || *r.HAE != 30 || r.Course != 45 || r.Speed != 7 || r.Origin != "site-mac" || !r.Simulated || r.Disposition != "unknown" || !e.fromCot || e.tlp != "AMBER+STRICT" {
		t.Fatalf("render %+v", r)
	}
	e2, _, _ := parseEntity("aircraft", "a1", []byte(`{"lat":1,"lon":2,"live":true}`), "hub")
	if e2.r.Course != -1 || e2.r.Speed != -1 || e2.r.HAE != nil || e2.tlp != "AMBER" {
		t.Fatalf("defaults %+v", e2.r)
	}
	if _, _, err := parseEntity("drone", "x", []byte(`{"entity_id":"y","entity_type":"drone"}`), "hub"); err == nil {
		t.Fatal("subject/id mismatch accepted")
	}
	if domain, ok := domainOf("$JS.site-mac.API"); !ok || domain != "site-mac" {
		t.Fatal("domainOf")
	}
	if _, ok := domainOf("$JS.API"); ok {
		t.Fatal("hub prefix taken as a site")
	}
}
