package identity

import (
	"testing"

	"github.com/onemindos/omos/services/entity/internal/model"
)

func TestMakeID(t *testing.T) {
	cases := map[[2]string]string{
		{"icao", "A1B2C3"}:              "icao-a1b2c3",
		{"mmsi", "367000001"}:           "mmsi-367000001",
		{"tak", "ANDROID-a1b2c3d4e5f6"}: "tak-android-a1b2c3d4e5f6",
		{"gtfs", "mnr-1234_5678"}:       "gtfs-mnr-1234_5678",
		{"", "sim-x1"}:                  "sim-x1",
	}
	for in, want := range cases {
		if got := MakeID(in[0], in[1]); got != want {
			t.Errorf("MakeID(%q,%q) = %q, want %q", in[0], in[1], got, want)
		}
	}
	// A CoT uid with '.' is hashed; the result is a valid, stable token (entity.md §3).
	h := MakeID("tak", "S-1-5-21-942292099.WINTAK")
	if !ValidID(h) || len(h) != len("tak-")+26 || h != MakeID("tak", "S-1-5-21-942292099.WINTAK") {
		t.Fatalf("hashed id %q", h)
	}
	if ValidID("Upper") || ValidID("a.b") || ValidID("a*b") || ValidID("") {
		t.Fatal("ValidID accepted an invalid token")
	}
}

func TestResolverPrecedence(t *testing.T) {
	r := NewResolver(map[string]string{"cot_uid:SIM-X1-ATAK": "sim-x1"})
	icao := model.Alias{Type: KeyICAO, ID: "a1b2c3"}
	uid := model.Alias{Type: KeyCotUID, ID: "ADSB-A1B2C3"}
	if id, _ := r.Resolve([]model.Alias{icao, uid}); id != "icao-a1b2c3" {
		t.Fatalf("canonical: %s", id)
	}
	r.Register("icao-a1b2c3", []model.Alias{icao, uid})
	if id, _ := r.Resolve([]model.Alias{uid}); id != "icao-a1b2c3" {
		t.Fatalf("alias lookup: %s", id)
	}
	if id, _ := r.Resolve([]model.Alias{{Type: KeyCotUID, ID: "SIM-X1-ATAK"}}); id != "sim-x1" {
		t.Fatalf("declared link: %s", id)
	}
	r.Register("other", []model.Alias{uid}) // an alias never moves silently
	if id, _ := r.Resolve([]model.Alias{uid}); id != "icao-a1b2c3" {
		t.Fatalf("alias moved: %s", id)
	}
	r.Forget("icao-a1b2c3")
	if id, _ := r.Resolve([]model.Alias{uid}); id == "icao-a1b2c3" {
		t.Fatal("forget kept alias")
	}
	if _, ok := r.Resolve(nil); ok {
		t.Fatal("empty keys resolved")
	}
}
