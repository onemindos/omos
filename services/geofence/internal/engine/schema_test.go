package engine

import (
	"bytes"
	"encoding/json"
	"path/filepath"
	"testing"
	"time"

	"github.com/santhosh-tekuri/jsonschema/v6"

	"github.com/onemindos/omos/services/geofence/internal/zone"
)

func compile(t *testing.T, rel string) *jsonschema.Schema {
	t.Helper()
	p, _ := filepath.Abs("../../../schemas/" + rel)
	s, err := jsonschema.NewCompiler().Compile(p)
	if err != nil {
		t.Fatal(err)
	}
	return s
}

func validate(s *jsonschema.Schema, v any) error {
	b, _ := json.Marshal(v)
	doc, err := jsonschema.UnmarshalJSON(bytes.NewReader(b))
	if err != nil {
		return err
	}
	return s.Validate(doc)
}

// TestFactsValidate: every alert the engine emits validates against
// services/schemas/alert/alert.v1.schema.json, across all three alert kinds.
func TestFactsValidate(t *testing.T) {
	as := compile(t, "alert/alert.v1.schema.json")
	g, c := newEngine()
	var outs []Output
	outs = append(outs, g.SetZone(mkZone(t, "home", zone.KindKeepIn, box(0, 0, 0.01), ""))...)
	outs = append(outs, g.SetEntity(drone("sim-x1", 0, 0, 1))...)
	outs = append(outs, g.SetEntity(drone("sim-x1", 1, 1, 2))...)
	st := drone("sim-x1", 1, 1, 3)
	st.Lifecycle = "lost"
	outs = append(outs, g.SetEntity(st)...)
	outs = append(outs, g.SetEntity(drone("sim-x1", 0, 0, 4))...)
	outs = append(outs, g.FeedStatus("nws", []byte(`{"stale":true,"lastError":"x"}`))...)
	c.t = c.t.Add(time.Hour)
	outs = append(outs, g.Tick()...)
	kinds := map[string]bool{}
	n := 0
	for _, o := range outs {
		a, ok := o.Data.(Alert)
		if !ok {
			continue
		}
		n++
		kinds[a.Kind+"/"+a.State] = true
		if err := validate(as, a); err != nil {
			t.Errorf("%s %s: %v", o.Subject, a.AlertID, err)
		}
	}
	if n < 6 || !kinds["geofence/open"] || !kinds["geofence/resolved"] || !kinds["stale_asset/open"] || !kinds["feed_outage/open"] {
		t.Fatalf("too few alert shapes covered: %d %v", n, kinds)
	}
}

// TestRecordsAndParamsValidate: records written by zone.Build validate against
// services/schemas/map/geofence.v1.schema.json, and the parameters used to build them
// validate against the zone.* action schemas the policy service would check.
func TestRecordsAndParamsValidate(t *testing.T) {
	rs := compile(t, "map/geofence.v1.schema.json")
	cs := compile(t, "actions/zone.create.v1.schema.json")
	us := compile(t, "actions/zone.update.v1.schema.json")
	rm := compile(t, "actions/zone.remove.v1.schema.json")
	lo, hi := 0.0, 100.0
	params := []map[string]any{
		{"name": "Home", "geo": map[string]any{"shape": "polygon", "kind": "keep_in", "coords": map[string]any{"type": "Polygon",
			"coordinates": [][][]float64{{{0, 0}, {1, 0}, {1, 1}, {0, 0}}}}}, "applies_to": []string{"sim-x1"}},
		{"geo": map[string]any{"shape": "circle", "kind": "keep_out", "coords": map[string]any{"type": "Point", "coordinates": []float64{0, 0}}, "radius_m": 300}},
		{"geo": map[string]any{"shape": "corridor", "kind": "aoi", "coords": map[string]any{"type": "LineString", "coordinates": [][]float64{{0, 0}, {1, 1}}},
			"radius_m": 100, "alt_min_m": lo, "alt_max_m": hi}, "watch": true, "watch_types": []string{"aircraft"}},
		{"polygon": [][]float64{{0, 0}, {1, 0}, {1, 1}}, "kind": "hazard"},
	}
	for i, p := range params {
		if err := validate(cs, p); err != nil {
			t.Errorf("params %d: %v", i, err)
		}
		b, _ := json.Marshal(p)
		var zp zone.Params
		_ = json.Unmarshal(b, &zp)
		r, rf := zone.Build("z", nil, zp, zone.Meta{Actor: "user:zeus", CommandID: "c", NowMs: 1})
		if rf != nil {
			t.Fatalf("params %d: %v", i, rf)
		}
		if err := validate(rs, r); err != nil {
			t.Errorf("record %d: %v", i, err)
		}
	}
	bad := []map[string]any{
		{"geo": map[string]any{"shape": "corridor", "kind": "keep_out", "coords": map[string]any{"type": "LineString", "coordinates": [][]float64{{0, 0}, {1, 1}}}, "radius_m": 5}},
		{"polygon": [][]float64{{0, 0}, {1, 0}, {1, 1}}, "kind": "aoi", "watch_types": []string{"person"}},
		{"name": "no shape"},
		{"polygon": [][]float64{{0, 0}, {1, 0}, {1, 1}}, "kind": "aoi", "surprise": 1},
	}
	for i, p := range bad {
		if validate(cs, p) == nil {
			t.Errorf("bad params %d accepted by zone.create", i)
		}
	}
	if err := validate(us, map[string]any{"active": false, "revision": 3}); err != nil {
		t.Error(err)
	}
	if err := validate(rm, map[string]any{"reason": "season over"}); err != nil {
		t.Error(err)
	}
}
