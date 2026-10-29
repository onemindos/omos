package geoid

import (
	"math"
	"testing"
)

// Expected values are bilinear interpolations of the GeographicLib 15-minute EGM96 grid,
// printed by the generator run on 2026-10-29 (0,0 -> 17.163; Austin -> -26.904;
// New York -> -32.760). The 0.5 degree grid must agree within its measured 1.95 m bound.
func TestUndulationKnownPoints(t *testing.T) {
	cases := []struct {
		name     string
		lat, lon float64
		want     float64
	}{
		{"null island", 0, 0, 17.163},
		{"austin", 30.269, -97.742, -26.904},
		{"new york", 40.7128, -74.006, -32.760},
	}
	for _, c := range cases {
		got, err := Undulation(c.lat, c.lon)
		if err != nil {
			t.Fatal(err)
		}
		if math.Abs(got-c.want) > 1.95 {
			t.Errorf("%s: N=%.3f want %.3f +/-1.95", c.name, got, c.want)
		}
	}
}

func TestRoundTripAndWrap(t *testing.T) {
	hae, err := MSLToHAE(10, 179.9, 100)
	if err != nil {
		t.Fatal(err)
	}
	msl, _ := HAEToMSL(10, 179.9, hae)
	if math.Abs(msl-100) > 1e-9 {
		t.Fatalf("round trip %.6f", msl)
	}
	a, _ := Undulation(10, -180)
	b, _ := Undulation(10, 180)
	if math.Abs(a-b) > 1e-9 {
		t.Fatalf("antimeridian mismatch %v %v", a, b)
	}
	for _, lat := range []float64{90, -90} {
		if _, err := Undulation(lat, 0); err != nil {
			t.Fatalf("pole %v: %v", lat, err)
		}
	}
	if _, err := Undulation(91, 0); err == nil {
		t.Fatal("expected error for lat 91")
	}
	if _, err := Undulation(math.NaN(), 0); err == nil {
		t.Fatal("expected error for NaN")
	}
}
