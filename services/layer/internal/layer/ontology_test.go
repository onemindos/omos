package layer

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"slices"
	"sort"
	"testing"
	"time"

	"github.com/nats-io/nats.go"

	lc "github.com/onemindos/omos/services/layer/internal/layerclient"
)

// entity-service records, shaped exactly as omos #552 (feat/entity-service 73823045) writes
// them to ENT: testdata/entity-golden.json is its services/entity/testdata/golden/entities.json
// and testdata/entity-nws-alert.json is the ent.event.noaa-… record its engine wrote for an
// omos-ingest NWS alert in the whole-loop run (onemind-xr docs/evidence/INT/E2E/run2).
func entityServiceRecords(t *testing.T) []map[string]any {
	t.Helper()
	var recs []map[string]any
	b, err := os.ReadFile("testdata/entity-golden.json")
	if err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal(b, &recs); err != nil {
		t.Fatal(err)
	}
	var nws map[string]any
	b, err = os.ReadFile("testdata/entity-nws-alert.json")
	if err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal(b, &nws); err != nil {
		t.Fatal(err)
	}
	return append(recs, nws)
}

// derived copies a record under a new id, kind and source subject. Used only for the two
// shapes the fixtures lack: a non-hazard event and a fleet ground vehicle.
func derived(t *testing.T, src map[string]any, id, kind, subject string) map[string]any {
	t.Helper()
	b, _ := json.Marshal(src)
	var m map[string]any
	_ = json.Unmarshal(b, &m)
	m["entity_id"], m["kind"] = id, kind
	srcs := m["sources"].([]any)
	srcs[0].(map[string]any)["subject"] = subject
	m["sources"] = srcs[:1]
	return m
}

func byID(recs []map[string]any, id string) map[string]any {
	for _, r := range recs {
		if r["entity_id"] == id {
			return r
		}
	}
	return nil
}

func subjectOf(r map[string]any) string {
	return "ent." + r["entity_type"].(string) + "." + r["entity_id"].(string) + ".state"
}

// wantLayers is the membership the default layers must give the entity-service records.
var wantLayers = map[string][]string{
	"fleet":   {"sim-x1", "sim-r1"},
	"tak":     {"tak-android-a1b2c3d4e5f6", "tak-pbe3cesvp2esvfbrw3eze2we2h"},
	"hazards": {"emsc-20261029_0001", "noaa-j4yqrxihy3cjsuswvrwidyxpfb", "usgs-us7000abcd"},
	"space":   {"norad-25544"},
	"transit": {"gtfs-lirr-go103_25_2611", "gtfs-mnr-1234_5678"},
	"air":     {"icao-a1b2c3", "icao-abc123"},
	"sea":     {"mmsi-367000001", "mmsi-367000001-x2"},
}

func layerIDs(s *Store, layer string) []string {
	s.mu.Lock()
	defer s.mu.Unlock()
	var out []string
	for _, e := range s.perLayer[s.layerIndex(layer)] {
		out = append(out, e.r.EntityID)
	}
	sort.Strings(out)
	return out
}

func sorted(v []string) []string {
	out := append([]string{}, v...)
	sort.Strings(out)
	return out
}

func fixtureSet(t *testing.T) []map[string]any {
	recs := entityServiceRecords(t)
	quake := byID(recs, "usgs-us7000abcd")
	// A conflict event (evt.world.conflict, engine kind "conflict") is an event, not a hazard.
	conflict := derived(t, quake, "evt-conflict-x1", "conflict", "evt.world.conflict.x1")
	// A fleet UGV is a `vehicle` of kind ugv (geo.go robot/vehicle branch).
	ugv := derived(t, byID(recs, "norad-25544"), "sim-r1", "ugv", "tel.fleet.vehicle.sim-r1.status")
	ugv["environment"] = "land"
	return append(recs, conflict, ugv)
}

// Bug 4 of the whole-loop run: hazards/space/transit/rf filtered on types the entity service
// never emits. The default layers must place #552's records by ONTOLOGY type and kind.
func TestDefaultLayersPlaceEntityServiceRecords(t *testing.T) {
	recs := fixtureSet(t)
	now := time.UnixMilli(1793275260000) // the golden records' created_t: all still live
	s := NewStore(DefaultLayers(), "hub", time.Minute, 1000, func() time.Time { return now })
	for i, r := range recs {
		b, _ := json.Marshal(r)
		s.Apply(subjectOf(r), b, "", uint64(i+1))
	}
	if s.bad != 0 {
		t.Fatalf("%d records rejected", s.bad)
	}
	for layer, want := range wantLayers {
		if got := layerIDs(s, layer); !slices.Equal(got, sorted(want)) {
			t.Errorf("layer %s: got %v want %v", layer, got, sorted(want))
		}
	}
	for _, layer := range []string{"zones", "tracks", "rf", "cameras"} {
		if got := layerIDs(s, layer); len(got) != 0 {
			t.Errorf("layer %s: got %v want empty", layer, got)
		}
	}
	// A kind change moves the entity between layers (fleet -> transit) on the next record.
	ugv := byID(recs, "sim-r1")
	ugv["kind"] = "train"
	b, _ := json.Marshal(ugv)
	s.Apply(subjectOf(ugv), b, "", uint64(len(recs)+1))
	if f, tr := layerIDs(s, "fleet"), layerIDs(s, "transit"); slices.Contains(f, "sim-r1") || !slices.Contains(tr, "sim-r1") {
		t.Fatalf("after kind change: fleet %v transit %v", f, tr)
	}
}

func TestQualifiesAndValidate(t *testing.T) {
	hz := DefaultLayers()[3]
	if hz.ID != "hazards" {
		t.Fatal("layer order")
	}
	cases := []struct {
		kind string
		subj []string
		want bool
	}{
		{"weather_alert", nil, true},
		{"event", []string{"intel.seismic.usgs.m4.x"}, true}, // USGS record without its own type
		{"wildfire", []string{"evt.world.fire.f1"}, true},
		{"conflict", []string{"evt.world.conflict.c1"}, false},
		{"event", []string{"intel.conflict.acled.event"}, false},
		{"", nil, false},
	}
	for _, c := range cases {
		if got := hz.Qualifies(c.kind, c.subj); got != c.want {
			t.Errorf("hazards.Qualifies(%q, %v) = %v", c.kind, c.subj, got)
		}
	}
	d := Def{ID: "x", Tier: TierSmall, EntityTypes: []string{"vehicle"}, Cap: 1, LimitS: 1, SourceSubjects: []string{"intel.*."}}
	if d.Validate() == nil {
		t.Fatal("wildcard source subject accepted")
	}
	d.SourceSubjects, d.Kinds = []string{"intel.weather"}, nil
	if d.Validate() == nil {
		t.Fatal("source subject without trailing dot accepted")
	}
	d.SourceSubjects, d.Kinds = nil, []string{"a b"}
	if d.Validate() == nil {
		t.Fatal("bad kind accepted")
	}
}

// End to end on the bus: #552-shaped records on ENT reach svc.layer.snapshot in the right
// layers, and the manifest says which kinds narrow the filter.
func TestA5OntologyLayersOnBus(t *testing.T) {
	srv := startServer(t, "")
	nc, js := connect(t, srv.url)
	setupBus(t, js, bucketTTL{})
	exp := time.Now().Add(time.Hour).UnixMilli()
	for i, r := range fixtureSet(t) {
		r["expiry_t"] = exp // the golden expiries are in the past on the wall clock
		m := nats.NewMsg(subjectOf(r))
		m.Data, _ = json.Marshal(r)
		m.Header.Set("Nats-Msg-Id", fmt.Sprintf("ont-%s-%d", r["entity_id"], i))
		if _, err := js.PublishMsg(context.Background(), m); err != nil {
			t.Fatal(err)
		}
	}
	r := startService(t, srv.url, testConfig(DefaultLayers()))
	r.waitReady(t)
	for layer, want := range wantLayers {
		if layer == "air" || layer == "sea" {
			continue // large tiers need an AOI request; the store test above covers them
		}
		var got []string
		waitFor(t, 10*time.Second, "snapshot "+layer, func() bool {
			snap, err := lc.GetSnapshot(nc, map[string]any{"v": 1, "layer": layer, "aoi": nil, "zoom": 12, "fields": "render"}, 5*time.Second)
			if err != nil {
				return false
			}
			got = got[:0]
			for _, it := range snap.Items {
				var e lc.Entity
				_ = json.Unmarshal(it, &e)
				got = append(got, e.EntityID)
			}
			sort.Strings(got)
			return slices.Equal(got, sorted(want))
		})
		t.Logf("snapshot %s: %v", layer, got)
	}
	var man *Manifest
	waitFor(t, 10*time.Second, "hazards manifest count 3", func() bool {
		m, err := getManifest(t, js, "hazards")
		man = m
		return err == nil && m.Count == 3
	})
	if !slices.Equal(man.Filters, []string{"ent.event.>"}) || !slices.Contains(man.Kinds, "weather_alert") || !slices.Contains(man.SourceSubjects, "intel.weather.") {
		t.Fatalf("hazards manifest %+v", man)
	}
	sp, _ := getManifest(t, js, "space")
	if sp == nil || sp.Count != 1 || !slices.Equal(sp.Kinds, []string{"satellite"}) {
		t.Fatalf("space manifest %+v", sp)
	}
	fl, _ := getManifest(t, js, "fleet")
	if fl == nil || fl.Count != 2 || !slices.Equal(fl.ExcludeKinds, []string{"satellite", "train"}) {
		t.Fatalf("fleet manifest %+v", fl)
	}
}
