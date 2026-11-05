package geo

import (
	"encoding/json"
	"testing"
)

func raw(s string) json.RawMessage { return json.RawMessage(s) }

func TestPolygonWithHole(t *testing.T) {
	a, err := Parse(ShapePolygon, raw(`{"type":"Polygon","coordinates":[
		[[-97.75,30.26],[-97.73,30.26],[-97.73,30.28],[-97.75,30.28]],
		[[-97.745,30.265],[-97.735,30.265],[-97.735,30.275],[-97.745,30.275],[-97.745,30.265]]]}`), 0, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	cases := []struct {
		lat, lon float64
		want     bool
	}{
		{30.262, -97.748, true},    // inside outer ring
		{30.27, -97.74, false},     // inside the hole
		{30.29, -97.74, false},     // north of it
		{30.2672, -97.7431, false}, // hole again
	}
	for _, c := range cases {
		if got := a.Contains(c.lat, c.lon, 0, false); got != c.want {
			t.Errorf("Contains(%v,%v)=%v want %v", c.lat, c.lon, got, c.want)
		}
	}
}

func TestCircleAndAltitude(t *testing.T) {
	lo, hi := 10.0, 120.0
	a, err := Parse(ShapeCircle, raw(`{"type":"Point","coordinates":[-97.7431,30.2672]}`), 500, &lo, &hi)
	if err != nil {
		t.Fatal(err)
	}
	if !a.Contains(30.2672, -97.7431, 0, false) {
		t.Error("centre without altitude must be inside")
	}
	if a.Contains(30.2672, -97.7431, 152.4, true) {
		t.Error("above alt_max_m must be outside")
	}
	if !a.Contains(30.2672, -97.7431, 50, true) {
		t.Error("inside the band must be inside")
	}
	if a.Contains(30.2672, -97.7331, 0, false) { // about 960 m east
		t.Error("960 m from the centre of a 500 m circle must be outside")
	}
}

func TestCorridorAndMultiPolygon(t *testing.T) {
	c, err := Parse(ShapeCorridor, raw(`{"type":"LineString","coordinates":[[-97.75,30.27],[-97.73,30.27]]}`), 200, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	if !c.Contains(30.2705, -97.74, 0, false) || c.Contains(30.275, -97.74, 0, false) {
		t.Error("corridor: 55 m off the line inside, 550 m off outside")
	}
	m, err := Parse(ShapeMultiPolygon, raw(`{"type":"MultiPolygon","coordinates":[
		[[[0,0],[1,0],[1,1],[0,1],[0,0]]], [[[5,5],[6,5],[6,6],[5,6],[5,5]]]]}`), 0, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	if !m.Contains(0.5, 0.5, 0, false) || !m.Contains(5.5, 5.5, 0, false) || m.Contains(3, 3, 0, false) {
		t.Error("multipolygon membership wrong (the gap between parts must be outside)")
	}
}

func TestParseRejects(t *testing.T) {
	lo, hi := 100.0, 50.0
	bad := []struct {
		shape, coords string
		r             float64
		lo, hi        *float64
	}{
		{ShapePolygon, `{"type":"Point","coordinates":[0,0]}`, 0, nil, nil},
		{ShapePolygon, `{"type":"Polygon","coordinates":[[[0,0],[1,1],[0,0]]]}`, 0, nil, nil},
		{ShapePolygon, `{"type":"Polygon","coordinates":[[[0,0],[1,0],[1,95],[0,0]]]}`, 0, nil, nil},
		{ShapeCircle, `{"type":"Point","coordinates":[0,0]}`, 0, nil, nil},
		{ShapeCircle, `{"type":"Point","coordinates":[0,0]}`, 300000, nil, nil},
		{ShapeCorridor, `{"type":"LineString","coordinates":[[0,0]]}`, 10, nil, nil},
		{"hexagon", `{"type":"Polygon","coordinates":[]}`, 0, nil, nil},
		{ShapeCircle, `{"type":"Point","coordinates":[0,0]}`, 10, &lo, &hi},
		{ShapePolygon, `not json`, 0, nil, nil},
	}
	for i, b := range bad {
		if _, err := Parse(b.shape, raw(b.coords), b.r, b.lo, b.hi); err == nil {
			t.Errorf("case %d (%s) accepted", i, b.shape)
		}
	}
}

func TestOpenRingIsClosed(t *testing.T) {
	a, err := Parse(ShapePolygon, raw(`{"type":"Polygon","coordinates":[[[0,0],[1,0],[1,1],[0,1]]]}`), 0, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	if !a.Contains(0.5, 0.5, 0, false) {
		t.Error("open ring should be closed and contain its centre")
	}
	c := Closed([][][]float64{{{0, 0}, {1, 0}, {1, 1}}})
	if n := len(c[0]); n != 4 || c[0][3][0] != 0 || c[0][3][1] != 0 {
		t.Errorf("Closed did not close the ring: %v", c)
	}
}

func TestIndex(t *testing.T) {
	x := NewIndex(0.1)
	small, _ := Parse(ShapeCircle, raw(`{"type":"Point","coordinates":[-97.7431,30.2672]}`), 300, nil, nil)
	big, _ := Parse(ShapePolygon, raw(`{"type":"Polygon","coordinates":[[[-120,20],[-70,20],[-70,50],[-120,50],[-120,20]]]}`), 0, nil, nil)
	x.Put("small", small)
	x.Put("big", big)
	if got := x.Containing(30.2672, -97.7431, 0, false); len(got) != 2 || got[0] != "big" || got[1] != "small" {
		t.Fatalf("Containing = %v", got)
	}
	if got := x.Containing(30.30, -97.7431, 0, false); len(got) != 1 || got[0] != "big" {
		t.Fatalf("Containing far = %v", got)
	}
	x.Remove("small")
	x.Put("big", small) // replace geometry under the same id
	if got := x.Containing(40, -100, 0, false); len(got) != 0 {
		t.Fatalf("replaced zone still matches old geometry: %v", got)
	}
	if x.Len() != 1 {
		t.Fatalf("Len = %d", x.Len())
	}
}
