package zone

import (
	"encoding/json"
	"strings"
	"testing"
)

func ptr[T any](v T) *T { return &v }

const square = `{"type":"Polygon","coordinates":[[[-97.75,30.26],[-97.73,30.26],[-97.73,30.28],[-97.75,30.28]]]}`

var meta = Meta{Actor: "user:zeus", ActionID: "act-1", CommandID: "c1", Simulated: true, NowMs: 1000}

func TestCreateDefaultsAndShapes(t *testing.T) {
	r, rf := Build("home-keepin", nil, Params{Name: ptr("Home"), Geo: &Geo{Shape: "polygon", Kind: KindKeepIn, Coords: json.RawMessage(square)}}, meta)
	if rf != nil {
		t.Fatal(rf)
	}
	if !r.Watch || !r.Active || !r.Live || r.Lifecycle != "active" || r.EntityType != "zone" || r.EntityID != "home-keepin" || r.ZoneID != "home-keepin" {
		t.Fatalf("defaults wrong: %+v", r)
	}
	if !strings.Contains(string(r.Geo.Coords), "[-97.75,30.26]]]") {
		t.Fatalf("polygon ring was not closed: %s", r.Geo.Coords)
	}
	r2, rf := Build("aoi1", nil, Params{Polygon: [][]float64{{0, 0}, {1, 0}, {1, 1}}, Kind: KindAOI}, meta)
	if rf != nil || r2.Watch || r2.Geo.Shape != "polygon" {
		t.Fatalf("doctrine shorthand: %+v %v", r2, rf)
	}
	r3, rf := Build("aoi2", nil, Params{Polygon: [][]float64{{0, 0}, {1, 0}, {1, 1}}, Kind: KindAOI, Watch: ptr(true)}, meta)
	if rf != nil || !r3.Watch {
		t.Fatalf("explicit watch: %+v %v", r3, rf)
	}
}

func TestRefusals(t *testing.T) {
	line := json.RawMessage(`{"type":"LineString","coordinates":[[0,0],[1,1]]}`)
	cases := []struct {
		id   string
		p    Params
		code string
	}{
		{"Bad.ID", Params{Polygon: [][]float64{{0, 0}, {1, 0}, {1, 1}}, Kind: KindAOI}, CodeBadSubject},
		{"z", Params{}, CodeSchemaInvalid},
		{"z", Params{Polygon: [][]float64{{0, 0}, {1, 0}, {1, 1}}, Kind: "nofly"}, CodeSchemaInvalid},
		{"z", Params{Geo: &Geo{Shape: "corridor", Kind: KindKeepOut, Coords: line, RadiusM: 50}}, CodeUnsupportedShape},
		{"z", Params{Geo: &Geo{Shape: "corridor", Kind: KindAOI, Coords: line}}, CodeSchemaInvalid},
		{"z", Params{Polygon: [][]float64{{0, 0}, {1, 0}, {1, 1}}, Kind: KindAOI, WatchTypes: &[]string{"aircraft", "person"}}, CodeForbidden},
		{"z", Params{Polygon: [][]float64{{0, 0}, {1, 0}, {1, 1}}, Kind: KindAOI, AppliesTo: &[]string{"bad id"}}, CodeSchemaInvalid},
		{"z", Params{ZoneID: "other", Polygon: [][]float64{{0, 0}, {1, 0}, {1, 1}}, Kind: KindAOI}, CodeSchemaInvalid},
	}
	for i, c := range cases {
		_, rf := Build(c.id, nil, c.p, meta)
		if rf == nil || rf.Code != c.code {
			t.Errorf("case %d: got %v want %s", i, rf, c.code)
		}
	}
	// a corridor is fine for an AOI
	if _, rf := Build("z", nil, Params{Geo: &Geo{Shape: "corridor", Kind: KindAOI, Coords: line, RadiusM: 50}}, meta); rf != nil {
		t.Errorf("AOI corridor refused: %v", rf)
	}
}

func TestUpdateMerges(t *testing.T) {
	r, _ := Build("z1", nil, Params{Name: ptr("A"), Geo: &Geo{Shape: "polygon", Kind: KindKeepOut, Coords: json.RawMessage(square)},
		AppliesTo: &[]string{"sim-x2", "sim-x1"}}, meta)
	m2 := meta
	m2.CommandID, m2.NowMs = "c2", 2000
	u, rf := Build("z1", &r, Params{Active: ptr(false)}, m2)
	if rf != nil {
		t.Fatal(rf)
	}
	if u.Name != "A" || u.Geo.Kind != KindKeepOut || !u.Watch || u.Active || u.Live || u.Lifecycle != "inactive" ||
		u.CreatedT != 1000 || u.UpdatedT != 2000 || u.CommandID != "c2" || strings.Join(u.Properties.AppliesTo, ",") != "sim-x1,sim-x2" {
		t.Fatalf("merge wrong: %+v", u)
	}
	// changing the kind resets watch to the new kind's default unless given
	u2, rf := Build("z1", &r, Params{Kind: KindAOI}, m2)
	if rf != nil || u2.Watch || u2.Geo.Kind != KindAOI {
		t.Fatalf("kind change: %+v %v", u2, rf)
	}
	if r.Properties.AppliesTo[0] != "sim-x1" {
		t.Fatal("update modified the previous record")
	}
}

func TestParseActivity(t *testing.T) {
	for _, c := range []struct {
		v    string
		want bool
	}{
		{`{"geo":{"shape":"polygon","kind":"aoi","coords":` + square + `}}`, true},
		{`{"active":false,"geo":{"shape":"polygon","kind":"aoi","coords":` + square + `}}`, false},
		{`{"live":false,"geo":{"shape":"polygon","kind":"aoi","coords":` + square + `}}`, false},
		{`{"lifecycle":"removed","geo":{"shape":"polygon","kind":"aoi","coords":` + square + `}}`, false},
	} {
		z, err := Parse("k", []byte(c.v), 3)
		if err != nil {
			t.Fatal(err)
		}
		if z.Active != c.want || z.ZoneID != "k" || z.Rev != 3 {
			t.Errorf("%s: active=%v id=%s", c.v, z.Active, z.ZoneID)
		}
	}
	if _, err := Parse("k", []byte(`{"geo":{"shape":"polygon","kind":"moat","coords":`+square+`}}`), 1); err == nil {
		t.Error("unknown kind accepted")
	}
}

// TestRecordReadableByProjectorAndPolicy decodes a written record with the field sets the two
// other readers use (projector gates.go geofenceRec, PR #563; policy geo.go ParseZone, PR #554).
func TestRecordReadableByProjectorAndPolicy(t *testing.T) {
	r, rf := Build("home", nil, Params{Geo: &Geo{Shape: "polygon", Kind: KindKeepIn, Coords: json.RawMessage(square)},
		AppliesTo: &[]string{"sim-x1"}}, meta)
	if rf != nil {
		t.Fatal(rf)
	}
	b, _ := json.Marshal(r)
	var proj struct {
		ZoneID    string `json:"zone_id"`
		Active    *bool  `json:"active"`
		Live      *bool  `json:"live"`
		Lifecycle string `json:"lifecycle"`
		Geo       *struct {
			Shape  string          `json:"shape"`
			Kind   string          `json:"kind"`
			Coords json.RawMessage `json:"coords"`
		} `json:"geo"`
	}
	if err := json.Unmarshal(b, &proj); err != nil || proj.ZoneID != "home" || proj.Active == nil || !*proj.Active ||
		proj.Live == nil || !*proj.Live || proj.Geo == nil || proj.Geo.Kind != "keep_in" {
		t.Fatalf("projector view: %+v %v", proj, err)
	}
	var pol struct {
		EntityID string `json:"entity_id"`
		Live     *bool  `json:"live"`
		Geo      *struct {
			Shape  string          `json:"shape"`
			Coords json.RawMessage `json:"coords"`
			Kind   string          `json:"kind"`
		} `json:"geo"`
		Properties struct {
			AppliesTo []string `json:"applies_to"`
		} `json:"properties"`
	}
	if err := json.Unmarshal(b, &pol); err != nil || pol.EntityID != "home" || pol.Geo.Shape != "polygon" ||
		len(pol.Properties.AppliesTo) != 1 {
		t.Fatalf("policy view: %+v %v", pol, err)
	}
	var gj struct {
		Type        string         `json:"type"`
		Coordinates [][][2]float64 `json:"coordinates"`
	}
	if err := json.Unmarshal(pol.Geo.Coords, &gj); err != nil || gj.Type != "Polygon" || len(gj.Coordinates[0]) < 4 {
		t.Fatalf("policy needs a closed Polygon ring of 4+ points: %s", pol.Geo.Coords)
	}
}
