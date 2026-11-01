package bus

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"math"
	"testing"
	"time"

	"github.com/onemindos/omos/services/entity/internal/adapt"
	"github.com/onemindos/omos/services/entity/internal/engine"
	"github.com/onemindos/omos/services/entity/internal/orbit"
)

const (
	issL1 = "1 25544U 98067A   26301.59601386  .00005724  00000+0  11161-3 0  9990"
	issL2 = "2 25544  51.6336 359.9178 0006994 320.4528  39.5949 15.49090321592436"
)

// TestSatellitePropagation (FEATURES A12): one TLE on intel.space.celestrak.tle keeps the ISS
// moving in ENT between CelesTrak polls, every position within 5 km of SGP4 at its measured
// time; after a restart the service resumes propagation from the restored entity alone.
func TestSatellitePropagation(t *testing.T) {
	s := startServer(t)
	nc, js := connect(t, s)
	setupStreams(t, js)

	base, start := time.Date(2026, 10, 29, 12, 0, 0, 0, time.UTC), time.Now()
	clock := func() time.Time { return base.Add(time.Since(start)) }
	oldNow := adapt.Now
	adapt.Now = clock
	defer func() { adapt.Now = oldNow }()

	rec := record(t, nc, "ent.vehicle.norad-25544.state")
	body := fmt.Sprintf(`{"id":"celestrak-25544","type":"tle","source":"celestrak","lat":null,"lon":null,"timestamp":"2026-10-28T14:18:15.597Z","properties":{"norad_id":"25544","name":"ISS (ZARYA)","group":"stations","line1":%q,"line2":%q}}`, issL1, issL2)
	if _, err := js.Publish(context.Background(), adapt.TLESubject, []byte(body)); err != nil {
		t.Fatal(err)
	}
	el, err := orbit.Parse(issL1, issL2)
	if err != nil {
		t.Fatal(err)
	}

	run := func(want int) []map[string]any {
		cfg := testConfig("orbit")
		cfg.OrbitEvery, cfg.clock = 150*time.Millisecond, clock
		eng := engine.New(engine.DefaultConfig(), func() int64 { return clock().UnixMilli() })
		svc, err := New(nc, eng, cfg, slog.New(slog.NewTextHandler(io.Discard, nil)))
		if err != nil {
			t.Fatal(err)
		}
		ctx, cancel := context.WithCancel(context.Background())
		errc := make(chan error, 1)
		go func() { errc <- svc.Run(ctx) }()
		deadline := time.Now().Add(15 * time.Second)
		for len(rec.all()) < want && time.Now().Before(deadline) {
			time.Sleep(50 * time.Millisecond)
		}
		cancel()
		if err := <-errc; err != nil {
			t.Fatalf("run: %v", err)
		}
		var out []map[string]any
		for _, m := range rec.all() {
			var r map[string]any
			if err := json.Unmarshal(m.Data, &r); err != nil {
				t.Fatal(err)
			}
			out = append(out, r)
		}
		if len(out) < want {
			t.Fatalf("%d ENT records, want %d", len(out), want)
		}
		return out
	}

	first := run(4)
	n1 := len(first)
	all := run(n1 + 3) // restart: no new TLE is published
	var lastSeq, lastT float64
	for i, r := range all {
		loc := r["location"].(map[string]any)
		mt := int64(loc["measured_t"].(float64))
		st, err := el.At(time.UnixMilli(mt))
		if err != nil {
			t.Fatal(err)
		}
		d := haversine(st.Lat, st.Lon, loc["lat"].(float64), loc["lon"].(float64))
		if d > 5000 {
			t.Errorf("record %d: %.0f m from SGP4 at %d", i, d, mt)
		}
		if r["seq"].(float64) <= lastSeq || float64(mt) < lastT {
			t.Errorf("record %d: seq %v t %d not increasing", i, r["seq"], mt)
		}
		lastSeq, lastT = r["seq"].(float64), float64(mt)
		if r["kind"] != "satellite" || r["environment"] != "space" || r["entity_type"] != "vehicle" {
			t.Fatalf("record %d: %v %v %v", i, r["kind"], r["environment"], r["entity_type"])
		}
	}
	t.Logf("%d records before restart, %d after; all within 5 km of SGP4", n1, len(all)-n1)
}

func haversine(lat1, lon1, lat2, lon2 float64) float64 {
	r := math.Pi / 180
	h := math.Pow(math.Sin((lat2-lat1)*r/2), 2) + math.Cos(lat1*r)*math.Cos(lat2*r)*math.Pow(math.Sin((lon2-lon1)*r/2), 2)
	return 2 * 6371008.8 * math.Asin(math.Sqrt(h))
}
