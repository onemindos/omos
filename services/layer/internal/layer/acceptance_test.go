package layer

// Acceptance tests of onemind-xr docs/contracts/layer-snapshot.md §7 and layer-delta.md §7,
// each against a disposable nats-server. Timers are shortened (heartbeat 0.5 s instead of
// 10 s, dead-watch 3 s instead of 30 s, bucket TTLs) where the contract value would only make
// the run longer; each test says which. Set A5_LONG=1 to also run the real-timer variants.

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"math/rand"
	"os"
	"runtime"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"

	lc "github.com/onemindos/omos/services/layer/internal/layerclient"
)

func drones(n int) []ent {
	out := make([]ent, n)
	for i := range out {
		out[i] = ent{typ: "drone", id: fmt.Sprintf("sim-d%05d", i), seq: 1, lat: 40.5 + float64(i%100)*0.005,
			lon: -74.3 + float64(i/100)*0.007, live: true, sim: true}
	}
	return out
}

// A5 one read + Sim only: 10,000-entity fleet layer in one svc.layer.snapshot call < 2 s,
// items == manifest count, every chunk <= 512 KiB and <= 2,000 items, every entity sim-*.
func TestA5OneRead(t *testing.T) {
	srv := startServer(t, "")
	nc, js := connect(t, srv.url)
	setupBus(t, js, bucketTTL{})
	publishEnts(t, js, drones(10000), "")
	r := startService(t, srv.url, testConfig([]Def{fleetLayer()}))
	r.waitReady(t)
	var man *Manifest
	waitFor(t, 10*time.Second, "manifest count 10000", func() bool {
		m, err := getManifest(t, js, "fleet")
		man = m
		return err == nil && m.Count == 10000
	})
	start := time.Now()
	snap, err := lc.GetSnapshot(nc, map[string]any{"v": 1, "layer": "fleet", "aoi": nil, "zoom": 12, "fields": "render"}, 5*time.Second)
	elapsed := time.Since(start)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("snapshot: %d items in %d chunks, %v, ent_seq %d, manifest count %d", len(snap.Items), len(snap.Chunks), elapsed, snap.EntSeq, man.Count)
	if elapsed > 2*time.Second {
		t.Fatalf("snapshot took %v (> 2 s)", elapsed)
	}
	if len(snap.Items) != man.Count {
		t.Fatalf("items %d != manifest count %d", len(snap.Items), man.Count)
	}
	for _, c := range snap.Chunks {
		if c.Bytes > MaxChunkBytes || len(c.Items) > MaxChunkItems {
			t.Fatalf("chunk %d: %d bytes, %d items", c.Part, c.Bytes, len(c.Items))
		}
	}
	for _, it := range snap.Items {
		var e lc.Entity
		_ = json.Unmarshal(it, &e)
		if !e.Simulated || !strings.HasPrefix(e.EntityID, "sim-") {
			t.Fatalf("non-simulated fleet entity %s", it)
		}
	}
	if man.Tier != "small" || man.Freshness.Status != "fresh" || man.Site != "hub" || man.EntSeq == 0 {
		t.Fatalf("manifest %+v", man)
	}
}

func vessels(n int, rng *rand.Rand, nycShare int) []ent {
	out := make([]ent, n)
	for i := range out {
		lat, lon := rng.Float64()*140-70, rng.Float64()*360-180
		if i < nycShare {
			lat, lon = 40.7+(rng.Float64()*2-1), -74.0+(rng.Float64()*2-1)
		}
		out[i] = ent{typ: "vessel", id: fmt.Sprintf("mmsi%06d", i), seq: 1, lat: lat, lon: lon, live: true}
	}
	return out
}

// AOI (A14): 50,000 vessels -> <= 5,000 cells at zoom 3; individual vessels inside a 50 km bbox
// at zoom 10; large layers refuse a request without an AOI; KV cells per zoom <= 512 KiB.
func TestA5AOICellsAndVessels(t *testing.T) {
	srv := startServer(t, "")
	nc, js := connect(t, srv.url)
	setupBus(t, js, bucketTTL{})
	fx := vessels(50000, rand.New(rand.NewSource(7)), 5000)
	zone := ent{typ: "zone", id: "zone-home", seq: 1, live: true, geo: map[string]any{"shape": "polygon", "kind": "aoi",
		"coords": map[string]any{"type": "Polygon", "coordinates": [][][]float64{{{-74.1, 40.6}, {-73.9, 40.6}, {-73.9, 40.8}, {-74.1, 40.8}, {-74.1, 40.6}}}}}}
	publishEnts(t, js, append(fx, zone), "")
	r := startService(t, srv.url, testConfig([]Def{seaLayer(), {ID: "zones", Title: "Zones", Tier: TierSmall, EntityTypes: []string{"zone"}, Cap: 1000, LimitS: 60}}))
	r.waitReady(t)
	waitFor(t, 15*time.Second, "sea manifest 50000", func() bool {
		m, err := getManifest(t, js, "sea")
		return err == nil && m.Count == 50000
	})

	world := map[string]any{"bbox": []float64{-180, -85, 180, 85}}
	snap, err := lc.GetSnapshot(nc, map[string]any{"v": 1, "layer": "sea", "aoi": world, "zoom": 3}, 5*time.Second)
	if err != nil {
		t.Fatal(err)
	}
	sum := 0.0
	for _, it := range snap.Items {
		var c []float64
		_ = json.Unmarshal(it, &c)
		sum += c[2]
	}
	t.Logf("zoom 3: mode %s, %d cells, %d vessels counted", snap.Mode, len(snap.Items), int(sum))
	if snap.Mode != "cells" || len(snap.Items) > 5000 || int(sum) != 50000 {
		t.Fatalf("zoom 3: mode %s, %d cells, sum %v", snap.Mode, len(snap.Items), sum)
	}

	dLat, dLon := 25000/111320.0, 25000/(111320.0*math.Cos(40.7*math.Pi/180))
	bbox := []float64{-74.0 - dLon, 40.7 - dLat, -74.0 + dLon, 40.7 + dLat}
	want := 0
	for _, v := range fx {
		if v.lon >= bbox[0] && v.lon <= bbox[2] && v.lat >= bbox[1] && v.lat <= bbox[3] {
			want++
		}
	}
	snap, err = lc.GetSnapshot(nc, map[string]any{"v": 1, "layer": "sea", "aoi": map[string]any{"bbox": bbox}, "zoom": 10}, 5*time.Second)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("zoom 10, 50 km bbox: mode %s, %d vessels (expected %d), truncated %v", snap.Mode, len(snap.Items), want, snap.Truncated)
	if snap.Mode != "entities" || len(snap.Items) != want || snap.Truncated || want == 0 {
		t.Fatalf("zoom 10: mode %s, %d items, want %d", snap.Mode, len(snap.Items), want)
	}
	for _, it := range snap.Items {
		var e lc.Entity
		_ = json.Unmarshal(it, &e)
		if e.Lon < bbox[0] || e.Lon > bbox[2] || e.Lat < bbox[1] || e.Lat > bbox[3] {
			t.Fatalf("vessel outside bbox: %s", it)
		}
	}

	// Zone AOI (zone_id resolves the zone entity's GeoJSON polygon).
	wantZ := 0
	for _, v := range fx {
		if v.lon >= -74.1 && v.lon <= -73.9 && v.lat >= 40.6 && v.lat <= 40.8 {
			wantZ++
		}
	}
	snap, err = lc.GetSnapshot(nc, map[string]any{"v": 1, "layer": "sea", "aoi": map[string]any{"zone_id": "zone-home"}, "zoom": 12}, 5*time.Second)
	if err != nil || len(snap.Items) != wantZ {
		t.Fatalf("zone aoi: err %v, %d items, want %d", err, len(snap.Items), wantZ)
	}
	t.Logf("zone_id zone-home: %d vessels (expected %d)", len(snap.Items), wantZ)

	// Truncation at max.
	snap, err = lc.GetSnapshot(nc, map[string]any{"v": 1, "layer": "sea", "aoi": world, "zoom": 12, "max": 100}, 5*time.Second)
	if err != nil || len(snap.Items) != 100 || !snap.Truncated {
		t.Fatalf("max 100: err %v, %d items, truncated %v", err, len(snap.Items), snap != nil && snap.Truncated)
	}
	// Large layers need an AOI.
	if _, err := lc.GetSnapshot(nc, map[string]any{"v": 1, "layer": "sea", "zoom": 10}, 5*time.Second); !lc.IsCode(err, ErrAOIRequired) {
		t.Fatalf("no aoi: %v", err)
	}
	if _, err := lc.GetSnapshot(nc, map[string]any{"v": 1, "layer": "nope"}, 5*time.Second); !lc.IsCode(err, ErrUnknownLayer) {
		t.Fatalf("unknown layer: %v", err)
	}
	// KV density cells.
	kv, _ := js.KeyValue(context.Background(), "world-state")
	for _, z := range []int{2, 4, 6} {
		e, err := kv.Get(context.Background(), fmt.Sprintf("sea.cells.z%d", z))
		if err != nil {
			t.Fatalf("cells z%d: %v", z, err)
		}
		var rec cellRecord
		_ = json.Unmarshal(e.Value(), &rec)
		total := 0.0
		for _, c := range rec.Cells {
			total += c[2]
		}
		t.Logf("KV sea.cells.z%d: %d cells, %d bytes, %d vessels", z, len(rec.Cells), len(e.Value()), int(total))
		if len(e.Value()) > MaxCellRecord || int(total) != 50000 || rec.Scheme != "webmercator-tile" {
			t.Fatalf("cells z%d: %d bytes, total %v", z, len(e.Value()), total)
		}
	}
	m, _ := getManifest(t, js, "sea")
	if m.Cells == nil || fmt.Sprint(m.Cells.Zooms) != "[2 4 6]" || m.Cells.BelowZoom != 7 || m.Cap != 5000 {
		t.Fatalf("manifest cells %+v", m.Cells)
	}
}

// Layer down: stopping the service makes `{layer}.meta` vanish by the bucket TTL, and the
// client shows "layer down". Bucket TTL is 2 s here (5 min in the CR); A5_LONG=1 runs 5 min.
func TestA5LayerDown(t *testing.T) {
	ttl := 2 * time.Second
	if os.Getenv("A5_LONG") == "1" {
		ttl = 5 * time.Minute
	}
	srv := startServer(t, "")
	_, js := connect(t, srv.url)
	setupBus(t, js, bucketTTL{world: ttl})
	publishEnts(t, js, drones(10), "")
	cfg := testConfig([]Def{fleetLayer()})
	cfg.ManifestEvery = 500 * time.Millisecond
	r := startService(t, srv.url, cfg)
	r.waitReady(t)
	waitFor(t, 5*time.Second, "manifest", func() bool { _, err := getManifest(t, js, "fleet"); return err == nil })
	if got := clientLayerBadge(t, js, "fleet"); got != "fresh" {
		t.Fatalf("badge %q", got)
	}
	r.stop()
	stopped := time.Now()
	waitFor(t, ttl+30*time.Second, "manifest to vanish", func() bool {
		_, err := getManifest(t, js, "fleet")
		return errors.Is(err, jetstream.ErrKeyNotFound)
	})
	t.Logf("manifest vanished %v after stop (bucket TTL %v); client badge %q", time.Since(stopped).Round(100*time.Millisecond), ttl, clientLayerBadge(t, js, "fleet"))
	if clientLayerBadge(t, js, "fleet") != "layer down" || time.Since(stopped) > ttl+30*time.Second {
		t.Fatal("client does not show layer down")
	}
}

// clientLayerBadge applies layer-snapshot.md §5 "Layer" rule.
func clientLayerBadge(t *testing.T, js jetstream.JetStream, layer string) string {
	m, err := getManifest(t, js, layer)
	if err != nil {
		return "layer down"
	}
	if m.Freshness.Status != "fresh" {
		return "stale"
	}
	return "fresh"
}

// Stale feed: stopping a collector moves freshness.status off fresh within its limit_s.
func TestA5StaleFeed(t *testing.T) {
	srv := startServer(t, "")
	nc, js := connect(t, srv.url)
	setupBus(t, js, bucketTTL{})
	hz := Def{ID: "hazards", Title: "Hazards", Tier: TierMedium, EntityTypes: []string{"hazard"}, Cap: 5000, LimitS: 2,
		Feeds: []string{"feed-a", "feed-b"}, DefaultStaleS: 3600, Marking: "CLEAR"}
	pub := func(src string, fetched time.Time, lastErr any) {
		b, _ := json.Marshal(map[string]any{"source": src, "kind": "rest_poll", "collector": "omos-ingest", "instance": "t", "leader": true,
			"fetchedAt": fetched.UTC().Format(time.RFC3339Nano), "count": 12, "lastError": lastErr, "staleAfterSec": 0, "stale": false, "at": time.Now().UTC().Format(time.RFC3339)})
		_ = nc.Publish("tel.feed."+src+".status", b)
	}
	// feed-a's last status is already in TEL before the service starts (read at start).
	pub("feed-a", time.Now(), nil)
	_ = nc.Flush()
	r := startService(t, srv.url, testConfig([]Def{hz}))
	r.waitReady(t)
	stopA, stopB := make(chan struct{}), make(chan struct{})
	collector := func(src string, stop chan struct{}) {
		tk := time.NewTicker(200 * time.Millisecond)
		defer tk.Stop()
		for {
			select {
			case <-stop:
				return
			case <-tk.C:
				pub(src, time.Now(), nil)
			}
		}
	}
	go collector("feed-a", stopA)
	go collector("feed-b", stopB)
	status := func() string {
		m, err := getManifest(t, js, "hazards")
		if err != nil {
			return ""
		}
		return m.Freshness.Status
	}
	waitFor(t, 5*time.Second, "fresh", func() bool { return status() == "fresh" })
	close(stopB)
	stopped := time.Now()
	waitFor(t, 5*time.Second, "degraded", func() bool { return status() == "degraded" })
	tDegraded := time.Since(stopped)
	close(stopA)
	stoppedA := time.Now()
	waitFor(t, 5*time.Second, "stale", func() bool { return status() == "stale" })
	tStale := time.Since(stoppedA)
	m, _ := getManifest(t, js, "hazards")
	t.Logf("limit_s 2: feed-b stopped -> degraded after %v; feed-a stopped -> stale after %v; feeds %+v", tDegraded.Round(10*time.Millisecond), tStale.Round(10*time.Millisecond), m.Freshness.Feeds)
	if tDegraded > 2*time.Second+500*time.Millisecond || tStale > 2*time.Second+500*time.Millisecond {
		t.Fatalf("status left fresh too late: %v / %v (limit 2 s + 0.5 s evaluation margin)", tDegraded, tStale)
	}
	for _, f := range m.Freshness.Feeds {
		if !f.Stale || f.FetchedAt == nil || f.Count != 12 {
			t.Fatalf("feed %+v", f)
		}
	}
	// lastError is copied through.
	pub("feed-a", time.Now(), "http 503")
	waitFor(t, 3*time.Second, "lastError copied", func() bool {
		m, _ := getManifest(t, js, "hazards")
		return m != nil && m.Freshness.Feeds[0].LastError != nil && *m.Freshness.Feeds[0].LastError == "http 503"
	})
}

// accountsConf: ONEMIND exports svc.layer.> as a streamed service; CLIENT imports it with
// share: true (BUS-DESIGN §7.1). Passwords are random per run (nothing in Git).
func accountsConf(pw map[string]string, threshold string) string {
	return fmt.Sprintf(`
accounts {
  ONEMIND {
    jetstream: enabled
    users: [ { user: svc, password: "%s" } ]
    exports: [ { service: "svc.layer.>", response_type: stream, response_threshold: "%s" } ]
  }
  CLIENT {
    users: [ { user: alice, password: "%s" }, { user: bob, password: "%s" }, { user: eve, password: "%s" } ]
    imports: [ { service: { account: ONEMIND, subject: "svc.layer.>" }, share: true } ]
  }
  NOSHARE {
    users: [ { user: mallory, password: "%s" } ]
    imports: [ { service: { account: ONEMIND, subject: "svc.layer.>" }, share: false } ]
  }
}
`, pw["svc"], threshold, pw["alice"], pw["bob"], pw["eve"], pw["mallory"])
}

// An import without share: true carries no Nats-Request-Info header; the layer service must
// not mistake such a request for a trusted same-account service (fail closed).
func TestA5ImportWithoutShare(t *testing.T) {
	pw := map[string]string{"svc": randHex(12), "alice": randHex(12), "bob": randHex(12), "eve": randHex(12), "mallory": randHex(12)}
	srv := startServer(t, accountsConf(pw, "2m"))
	_, js := connect(t, srv.url, nats.UserInfo("svc", pw["svc"]))
	setupBus(t, js, bucketTTL{})
	publishEnts(t, js, []ent{{typ: "drone", id: "sim-r", seq: 1, lat: 40.7, lon: -74, live: true, sim: true, tlp: "RED"}}, "")
	r := startService(t, srv.url, testConfig([]Def{fleetLayer()}), nats.UserInfo("svc", pw["svc"]))
	r.waitReady(t)
	mallory, _ := connect(t, srv.url, nats.UserInfo("mallory", pw["mallory"]))
	req := map[string]any{"v": 1, "layer": "fleet", "aoi": nil, "zoom": 12}
	if _, err := lc.GetSnapshot(mallory, req, 5*time.Second); !lc.IsCode(err, ErrForbidden) {
		t.Fatalf("share:false snapshot: %v", err)
	}
	rw := startRawWatch(t, mallory, map[string]any{"layer": "fleet", "aoi": nil, "zoom": 12})
	waitFor(t, 5*time.Second, "watch reply", func() bool { return rw.Err() != "" })
	if rw.Err() != ErrForbidden {
		t.Fatalf("share:false watch: %q", rw.Err())
	}
	list, err := mallory.Request("svc.layer.list", []byte(`{"v":1}`), 5*time.Second)
	if err != nil || !strings.Contains(string(list.Data), ErrForbidden) {
		t.Fatalf("share:false list: %v %s", err, list.Data)
	}
	t.Log("CLIENT import with share:false: snapshot, watch and list all refused with forbidden")
}

// Releasability: an entity marked AMBER+STRICT for team-a is absent from a CLIENT snapshot
// for team-b (and from its watch), present for team-a; an identity not in the directory is
// refused. Also proves streamed multi-chunk replies and watch batches across the account
// boundary (BUS-DESIGN §12 "unverified": streamed responses to a CLIENT inbox).
func TestA5Releasability(t *testing.T) {
	pw := map[string]string{"svc": randHex(12), "alice": randHex(12), "bob": randHex(12), "eve": randHex(12), "mallory": randHex(12)}
	threshold := "15s"
	if v := os.Getenv("A5_RESPONSE_THRESHOLD"); v != "" { // experiment knob, see README
		threshold = v
	}
	srv := startServer(t, accountsConf(pw, threshold))
	_, js := connect(t, srv.url, nats.UserInfo("svc", pw["svc"]))
	setupBus(t, js, bucketTTL{})
	fx := drones(2500)
	fx = append(fx,
		ent{typ: "drone", id: "sim-a", seq: 1, lat: 40.7, lon: -74, live: true, sim: true, tlp: "AMBER+STRICT", releasable: []string{"team-a"}},
		ent{typ: "drone", id: "sim-c", seq: 1, lat: 40.7, lon: -74, live: true, sim: true, tlp: "AMBER"},
		ent{typ: "drone", id: "sim-r", seq: 1, lat: 40.7, lon: -74, live: true, sim: true, tlp: "RED"})
	publishEnts(t, js, fx, "")
	cfg := testConfig([]Def{fleetLayer()})
	cfg.Directory = &Directory{OwnAccount: "ONEMIND", OrgGroup: "onemind", Principals: []Principal{
		{Account: "CLIENT", User: "alice", ID: "user:alice", Groups: []string{"team-a", "onemind"}},
		{Account: "CLIENT", User: "bob", ID: "user:bob", Groups: []string{"team-b"}},
	}}
	r := startService(t, srv.url, cfg, nats.UserInfo("svc", pw["svc"]))
	r.waitReady(t)
	alice, _ := connect(t, srv.url, nats.UserInfo("alice", pw["alice"]))
	bob, _ := connect(t, srv.url, nats.UserInfo("bob", pw["bob"]))
	eve, _ := connect(t, srv.url, nats.UserInfo("eve", pw["eve"]))
	req := map[string]any{"v": 1, "layer": "fleet", "aoi": nil, "zoom": 12}
	ids := func(s *lc.Snapshot) map[string]bool {
		out := map[string]bool{}
		for _, it := range s.Items {
			var e lc.Entity
			_ = json.Unmarshal(it, &e)
			out[e.EntityID] = true
		}
		return out
	}
	sb, err := lc.GetSnapshot(bob, req, 5*time.Second)
	if err != nil {
		t.Fatal(err)
	}
	sa, err := lc.GetSnapshot(alice, req, 5*time.Second)
	if err != nil {
		t.Fatal(err)
	}
	ib, ia := ids(sb), ids(sa)
	t.Logf("CLIENT bob (team-b): %d items in %d chunks; alice (team-a, org): %d items in %d chunks", len(sb.Items), len(sb.Chunks), len(sa.Items), len(sa.Chunks))
	if ib["sim-a"] || ib["sim-c"] || ib["sim-r"] || len(sb.Items) != 2500 {
		t.Fatalf("team-b sees restricted entities: a=%v c=%v r=%v n=%d", ib["sim-a"], ib["sim-c"], ib["sim-r"], len(sb.Items))
	}
	if !ia["sim-a"] || !ia["sim-c"] || ia["sim-r"] || len(sa.Items) != 2502 {
		t.Fatalf("team-a view wrong: a=%v c=%v r=%v n=%d", ia["sim-a"], ia["sim-c"], ia["sim-r"], len(sa.Items))
	}
	if len(sb.Chunks) < 2 {
		t.Fatalf("expected a multi-chunk streamed reply across accounts, got %d chunks", len(sb.Chunks))
	}
	if _, err := lc.GetSnapshot(eve, req, 5*time.Second); !lc.IsCode(err, ErrForbidden) {
		t.Fatalf("unknown identity: %v", err)
	}

	// Watch path across accounts, longer than the response threshold, with renewals.
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	fa := &lc.Follower{NC: alice, Req: map[string]any{"layer": "fleet", "aoi": nil, "zoom": 12, "lease_s": 10, "max_rate_hz": 4}, DeadAfter: 3 * time.Second}
	fb := &lc.Follower{NC: bob, Req: map[string]any{"layer": "fleet", "aoi": nil, "zoom": 12, "lease_s": 10, "max_rate_hz": 4}, DeadAfter: 3 * time.Second}
	go func() { _ = fa.Run(ctx) }()
	go func() { _ = fb.Run(ctx) }()
	waitFor(t, 10*time.Second, "both followers loaded", func() bool { return len(fa.Entities()) == 2502 && len(fb.Entities()) == 2500 })
	// Move sim-a for 25 s (longer than lease 10 s and threshold 15 s).
	end := time.Now().Add(25 * time.Second)
	seq := uint64(1)
	for time.Now().Before(end) {
		seq++
		publishEnts(t, js, []ent{{typ: "drone", id: "sim-a", seq: seq, lat: 40.7 + float64(seq)*1e-4, lon: -74, live: true, sim: true, tlp: "AMBER+STRICT", releasable: []string{"team-a"}}}, "")
		time.Sleep(200 * time.Millisecond)
	}
	waitFor(t, 5*time.Second, "alice has the last sim-a", func() bool { return fa.Entities()["sim-a"].Seq == seq })
	if _, ok := fb.Entities()["sim-a"]; ok {
		t.Fatal("team-b watch received sim-a")
	}
	na, ra := fa.Resnapshots()
	nb, rb := fb.Resnapshots()
	t.Logf("cross-account watches over 25 s (lease 10 s, response_threshold %s): alice %d batches, %d resnapshots %v; bob %d batches, %d resnapshots %v",
		threshold, len(fa.Batches()), na, ra, len(fb.Batches()), nb, rb)
	if na != 0 || nb != 0 {
		t.Fatal("cross-account watch had to re-snapshot")
	}
	// Re-mark sim-a for team-b: bob gets it (enter), alice loses it (remove released_out).
	seq++
	publishEnts(t, js, []ent{{typ: "drone", id: "sim-a", seq: seq, lat: 40.7, lon: -74, live: true, sim: true, tlp: "AMBER+STRICT", releasable: []string{"team-b"}}}, "")
	waitFor(t, 5*time.Second, "release change", func() bool {
		_, inA := fa.Entities()["sim-a"]
		_, inB := fb.Entities()["sim-a"]
		return !inA && inB
	})
	_, _, rem := fa.Counts()
	if rem["sim-a"] != "released_out" {
		t.Fatalf("alice remove reason %q", rem["sim-a"])
	}
}

// fleetChurn publishes rate msg/s of upserts and removals over n sim drones for d.
func fleetChurn(t *testing.T, js jetstream.JetStream, n, rate int, d time.Duration, seed int64) {
	rng := rand.New(rand.NewSource(seed))
	seqs := make([]uint64, n)
	live := make([]bool, n)
	tk := time.NewTicker(time.Second / time.Duration(rate))
	defer tk.Stop()
	end := time.Now().Add(d)
	for time.Now().Before(end) {
		<-tk.C
		i := rng.Intn(n)
		seqs[i]++
		l := rng.Intn(20) != 0 // 5 % removals
		live[i] = l
		publishEnts(t, js, []ent{{typ: "drone", id: fmt.Sprintf("sim-c%04d", i), seq: seqs[i], lat: 40 + rng.Float64(), lon: -74 + rng.Float64(), live: l, sim: true}}, "")
	}
}

func compareToReference(t *testing.T, js jetstream.JetStream, f *lc.Follower, typ string) (int, bool) {
	ref := referenceState(t, js, typ)
	got := f.Entities()
	if len(ref) != len(got) {
		return len(ref), false
	}
	for id, s := range ref {
		if g, ok := got[id]; !ok || g.Seq != s {
			return len(ref), false
		}
	}
	return len(ref), true
}

// No gap: snapshot then watch from ent_seq + 1 during a 200 msg/s replay; the client's entity
// set equals the reference built from the whole stream.
func TestA5NoGap(t *testing.T) {
	srv := startServer(t, "")
	nc, js := connect(t, srv.url)
	setupBus(t, js, bucketTTL{})
	r := startService(t, srv.url, testConfig([]Def{fleetLayer()}))
	r.waitReady(t)
	done := make(chan struct{})
	go func() { fleetChurn(t, js, 300, 200, 6*time.Second, 1); close(done) }()
	time.Sleep(2 * time.Second) // the client joins mid-stream
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	f := &lc.Follower{NC: nc, Req: map[string]any{"layer": "fleet", "aoi": nil, "zoom": 12, "lease_s": 30, "max_rate_hz": 10}, DeadAfter: 3 * time.Second}
	go func() { _ = f.Run(ctx) }()
	<-done
	var n int
	waitFor(t, 10*time.Second, "client equals reference", func() bool {
		var ok bool
		n, ok = compareToReference(t, js, f, "drone")
		return ok
	})
	resn, why := f.Resnapshots()
	b := f.Batches()
	t.Logf("200 msg/s for 6 s, client joined at 2 s: %d entities match the stream reference; %d batches; resnapshots %d %v", n, len(b), resn, why)
	if resn != 0 {
		t.Fatal("unexpected re-snapshot")
	}
}

// Resync after a gap: (1) a dropped batch (n skips) makes the client re-snapshot and
// converge; (2) a from_seq older than the replay window gets {"error":"resnapshot"}.
func TestA5ResyncAfterGap(t *testing.T) {
	srv := startServer(t, "")
	nc, js := connect(t, srv.url)
	setupBus(t, js, bucketTTL{})
	cfg := testConfig([]Def{fleetLayer()})
	cfg.ReplayMax = 50
	r := startService(t, srv.url, cfg)
	r.waitReady(t)
	publishEnts(t, js, drones(100), "")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	dropped := false
	var mu sync.Mutex
	f := &lc.Follower{NC: nc, Req: map[string]any{"layer": "fleet", "aoi": nil, "zoom": 12, "lease_s": 30, "max_rate_hz": 10}, DeadAfter: 3 * time.Second,
		DropBatch: func(b *lc.Batch) bool {
			mu.Lock()
			defer mu.Unlock()
			if !dropped && b.N == 3 {
				dropped = true
				return true
			}
			return false
		}}
	go func() { _ = f.Run(ctx) }()
	done := make(chan struct{})
	go func() { fleetChurn(t, js, 50, 100, 3*time.Second, 2); close(done) }()
	<-done
	waitFor(t, 10*time.Second, "client re-snapshotted and equals reference", func() bool {
		n, _ := f.Resnapshots()
		_, ok := compareToReference(t, js, f, "drone")
		return n >= 1 && ok
	})
	n, why := f.Resnapshots()
	t.Logf("dropped batch 3: client re-snapshotted %d time(s) %v and converged", n, why)
	if !strings.HasPrefix(why[0], "gap: batch 4 after 2") {
		t.Fatalf("reason %q", why[0])
	}

	// from_seq outside the window (ReplayMax 50 changes).
	snap, err := lc.GetSnapshot(nc, map[string]any{"v": 1, "layer": "fleet"}, 5*time.Second)
	if err != nil {
		t.Fatal(err)
	}
	fleetChurn(t, js, 50, 500, 400*time.Millisecond, 3) // about 200 changes > 50
	body, _ := json.Marshal(map[string]any{"v": 1, "layer": "fleet", "from_seq": snap.EntSeq + 1, "lease_s": 10})
	m, err := nc.Request("svc.layer.watch", body, 5*time.Second)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("watch from_seq %d after ~200 changes with a 50-change window: %s", snap.EntSeq+1, m.Data)
	if !strings.Contains(string(m.Data), `"error":"resnapshot"`) {
		t.Fatalf("expected resnapshot, got %s", m.Data)
	}
	// From the current sequence the watch starts.
	body, _ = json.Marshal(map[string]any{"v": 1, "layer": "fleet", "from_seq": r.svc.Store().LastSeq() + 1, "lease_s": 10})
	m, err = nc.Request("svc.layer.watch", body, 5*time.Second)
	if err != nil || !strings.Contains(string(m.Data), `"watch_id":"w-`) {
		t.Fatalf("watch from current seq: %v %s", err, m.Data)
	}
}

// Expiry: an aircraft with no new observation disappears within 5 s after its TTL on both
// paths (TTL 3 s here instead of 120 s). Also covers ENT without per-message TTL: the
// service's expiry sweep removes it at expiry_t.
func TestA5Expiry(t *testing.T) {
	srv := startServer(t, "")
	nc, js := connect(t, srv.url)
	setupBus(t, js, bucketTTL{})
	r := startService(t, srv.url, testConfig([]Def{airLayer()}))
	r.waitReady(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	aoi := map[string]any{"bbox": []float64{-75, 40, -73, 42}}
	f := &lc.Follower{NC: nc, Req: map[string]any{"layer": "air", "aoi": aoi, "zoom": 10, "lease_s": 30, "max_rate_hz": 10}, DeadAfter: 3 * time.Second}
	go func() { _ = f.Run(ctx) }()
	waitFor(t, 5*time.Second, "watch running", func() bool { return len(f.Batches()) > 0 })

	// Direct path: an ordered consumer on ENT sees the MaxAge delete marker.
	markerAt := make(chan time.Time, 1)
	oc, err := js.OrderedConsumer(ctx, "ENT", jetstream.OrderedConsumerConfig{FilterSubjects: []string{"ent.aircraft.sim-air1.state"}})
	if err != nil {
		t.Fatal(err)
	}
	cc, err := oc.Consume(func(m jetstream.Msg) {
		if m.Headers().Get("Nats-Marker-Reason") == "MaxAge" {
			select {
			case markerAt <- time.Now():
			default:
			}
		}
	})
	if err != nil {
		t.Fatal(err)
	}
	defer cc.Stop()

	exp1 := time.Now().Add(3 * time.Second)
	publishEnts(t, js, []ent{{typ: "aircraft", id: "sim-air1", seq: 1, lat: 41, lon: -74, live: true, sim: true, expiryT: ptr(exp1.UnixMilli())}}, "3s")
	exp2 := time.Now().Add(3 * time.Second)
	publishEnts(t, js, []ent{{typ: "aircraft", id: "sim-air2", seq: 1, lat: 41, lon: -74, live: true, sim: true, expiryT: ptr(exp2.UnixMilli())}}, "") // no Nats-TTL
	waitFor(t, 5*time.Second, "both aircraft held", func() bool { return len(f.Entities()) == 2 })
	waitFor(t, 10*time.Second, "both removed on the watch path", func() bool { return len(f.Entities()) == 0 })
	gone := time.Now()
	var mAt time.Time
	select {
	case mAt = <-markerAt:
	case <-time.After(8 * time.Second):
		t.Fatal("no MaxAge marker on the direct path")
	}
	_, _, rem := f.Counts()
	t.Logf("watch path: both removed %v after expiry (reasons %v); direct path: MaxAge marker %v after expiry",
		gone.Sub(exp2).Round(10*time.Millisecond), rem, mAt.Sub(exp1).Round(10*time.Millisecond))
	if gone.Sub(exp2) > 5*time.Second || mAt.Sub(exp1) > 5*time.Second {
		t.Fatal("removal later than 5 s after TTL")
	}
	if rem["sim-air1"] != "expired" || rem["sim-air2"] != "expired" {
		t.Fatalf("reasons %v", rem)
	}
}

func grid(n int, lat0, lon0, span float64) []ent {
	side := int(math.Ceil(math.Sqrt(float64(n))))
	out := make([]ent, 0, n)
	for i := 0; i < n; i++ {
		out = append(out, ent{typ: "vessel", id: fmt.Sprintf("mmsi9%05d", i), seq: 1, lat: lat0 + span*float64(i/side)/float64(side),
			lon: lon0 + span*float64(i%side)/float64(side), live: true})
	}
	return out
}

func inBox(lat, lon float64, b []float64) bool {
	return lon >= b[0] && lon <= b[2] && lat >= b[1] && lat <= b[3]
}

// AOI move: panning produces exit for entities that left and enter for new ones; no entity
// outside the AOI stays.
func TestA5AOIMove(t *testing.T) {
	srv := startServer(t, "")
	nc, js := connect(t, srv.url)
	setupBus(t, js, bucketTTL{})
	fx := grid(2000, 40, -75, 2)
	publishEnts(t, js, fx, "")
	r := startService(t, srv.url, testConfig([]Def{seaLayer()}))
	r.waitReady(t)
	A := []float64{-75, 40, -74, 41}
	B := []float64{-74.5, 40.5, -73.5, 41.5}
	want := func(b []float64) map[string]bool {
		o := map[string]bool{}
		for _, v := range fx {
			if inBox(v.lat, v.lon, b) {
				o[v.id] = true
			}
		}
		return o
	}
	wa, wb := want(A), want(B)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	f := &lc.Follower{NC: nc, Req: map[string]any{"layer": "sea", "aoi": map[string]any{"bbox": A}, "zoom": 10, "lease_s": 30, "max_rate_hz": 10}, DeadAfter: 3 * time.Second}
	go func() { _ = f.Run(ctx) }()
	same := func(w map[string]bool) bool {
		got := f.Entities()
		if len(got) != len(w) {
			return false
		}
		for id := range got {
			if !w[id] {
				return false
			}
		}
		return true
	}
	waitFor(t, 5*time.Second, "AOI A loaded", func() bool { return f.Watching() && same(wa) })
	if err := f.SetAOI(map[string]any{"bbox": B}, 10); err != nil {
		t.Fatal(err)
	}
	waitFor(t, 10*time.Second, "AOI B view", func() bool { return same(wb) })
	en, ex, _ := f.Counts()
	left, entered := 0, 0
	for id := range wa {
		if !wb[id] {
			left++
			if ex[id] != 1 {
				t.Fatalf("%s left the AOI without exactly one exit (%d)", id, ex[id])
			}
		}
	}
	for id := range wb {
		if !wa[id] {
			entered++
			if en[id] != 1 {
				t.Fatalf("%s entered without exactly one enter (%d)", id, en[id])
			}
		}
	}
	for id, e := range f.Entities() {
		if !inBox(e.Lat, e.Lon, B) {
			t.Fatalf("%s outside AOI B stays", id)
		}
	}
	t.Logf("AOI A %d vessels -> AOI B %d vessels: %d exit ops, %d enter ops, none outside B", len(wa), len(wb), left, entered)
	resn, _ := f.Resnapshots()
	if resn != 0 {
		t.Fatal("unexpected re-snapshot")
	}
	// Zooming out switches the watch to cells; zooming in again returns entities.
	if err := f.SetAOI(map[string]any{"bbox": B}, 3); err != nil {
		t.Fatal(err)
	}
	waitFor(t, 5*time.Second, "cells mode", func() bool { return len(f.Entities()) == 0 && len(f.Cells()) > 0 })
	total := 0.0
	for _, c := range f.Cells() {
		total += c[2]
	}
	if int(total) != len(wb) {
		t.Fatalf("cells count %v, want %d", total, len(wb))
	}
	if err := f.SetAOI(map[string]any{"bbox": B}, 10); err != nil {
		t.Fatal(err)
	}
	waitFor(t, 5*time.Second, "entities again", func() bool { return same(wb) })
}

// rawWatch opens a watch with from_seq 0 and records its batches and the first error.
type rawWatch struct {
	mu      sync.Mutex
	batches int
	err     string
}

func startRawWatch(t *testing.T, nc *nats.Conn, req map[string]any) *rawWatch {
	t.Helper()
	rw := &rawWatch{}
	inbox := nc.NewInbox()
	if _, err := nc.Subscribe(inbox, func(m *nats.Msg) {
		var b struct {
			N     uint64 `json:"n"`
			Error string `json:"error"`
		}
		_ = json.Unmarshal(m.Data, &b)
		rw.mu.Lock()
		defer rw.mu.Unlock()
		if b.Error != "" && rw.err == "" {
			rw.err = b.Error
		}
		if b.N > 0 {
			rw.batches++
		}
	}); err != nil {
		t.Fatal(err)
	}
	req["v"] = 1
	body, _ := json.Marshal(req)
	if err := nc.PublishRequest("svc.layer.watch", inbox, body); err != nil {
		t.Fatal(err)
	}
	return rw
}

func (rw *rawWatch) Err() string {
	rw.mu.Lock()
	defer rw.mu.Unlock()
	return rw.err
}

// Backpressure: a 50,000-vessel burst into a large layer watched by an AOI entity watch, four
// world cell watches and eight world entity watches (zoom 10, which would hold every vessel).
// Every batch stays <= 256 KiB and <= 1,000 ops at most max_rate_hz per second; the world
// entity watches end with too_many once their view passes the layer cap (5,000) instead of
// holding 50,000 entities each; the heap stays well under the 512Mi pod limit.
func TestA5Backpressure(t *testing.T) {
	srv := startServer(t, "")
	nc, js := connect(t, srv.url)
	setupBus(t, js, bucketTTL{})
	r := startService(t, srv.url, testConfig([]Def{seaLayer()}))
	r.waitReady(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	fx := vessels(50000, rand.New(rand.NewSource(9)), 0)
	box := []float64{-180, -70, -120, 5}
	want := 0
	for _, v := range fx {
		if inBox(v.lat, v.lon, box) {
			want++
		}
	}
	if want == 0 || want > 5000 {
		t.Fatalf("fixture: %d vessels in the AOI, want 1..5000", want)
	}
	world := map[string]any{"bbox": []float64{-180, -85, 180, 85}}
	f := &lc.Follower{NC: nc, Req: map[string]any{"layer": "sea", "aoi": map[string]any{"bbox": box}, "zoom": 10, "lease_s": 60, "max_rate_hz": 4}, DeadAfter: 5 * time.Second}
	go func() { _ = f.Run(ctx) }()
	var cells []*lc.Follower
	for i := 0; i < 4; i++ {
		fc := &lc.Follower{NC: nc, Req: map[string]any{"layer": "sea", "aoi": world, "zoom": 5, "lease_s": 60, "max_rate_hz": 4}, DeadAfter: 5 * time.Second}
		cells = append(cells, fc)
		go func() { _ = fc.Run(ctx) }()
	}
	var raws []*rawWatch
	for i := 0; i < 8; i++ {
		raws = append(raws, startRawWatch(t, nc, map[string]any{"layer": "sea", "aoi": world, "zoom": 10, "lease_s": 60, "max_rate_hz": 4}))
	}
	waitFor(t, 5*time.Second, "watches running", func() bool {
		for _, fc := range cells {
			if len(fc.Batches()) == 0 {
				return false
			}
		}
		return len(f.Batches()) > 0 && r.svc.Store().WatchCount() == 13
	})
	var before runtime.MemStats
	runtime.GC()
	runtime.ReadMemStats(&before)
	publishEnts(t, js, fx, "")
	maxPend, peakHeap, peakWatchBytes := 0, uint64(0), int64(0)
	sample := func() {
		var ms runtime.MemStats
		runtime.ReadMemStats(&ms)
		peakHeap = max(peakHeap, ms.HeapAlloc)
		peakWatchBytes = max(peakWatchBytes, r.svc.Store().WatchBytes())
		r.svc.Store().wmu.RLock()
		for _, w := range r.svc.Store().watches {
			maxPend = max(maxPend, w.pendingLen())
		}
		r.svc.Store().wmu.RUnlock()
	}
	waitFor(t, 60*time.Second, "AOI watch has its vessels; world entity watches ended", func() bool {
		sample()
		for _, rw := range raws {
			if rw.Err() == "" {
				return false
			}
		}
		return len(f.Entities()) == want
	})
	for i, rw := range raws {
		if rw.Err() != ErrTooMany {
			t.Fatalf("world entity watch %d ended with %q, want too_many", i, rw.Err())
		}
	}
	check := func(name string, bs []lc.Batch) (int, int, time.Duration) {
		maxB, maxOps, minGap := 0, 0, time.Hour
		for i, b := range bs {
			maxB, maxOps = max(maxB, b.Bytes), max(maxOps, len(b.Ops))
			if b.Bytes > MaxBatchBytes || len(b.Ops) > MaxBatchOps {
				t.Fatalf("%s batch %d: %d bytes, %d ops", name, b.N, b.Bytes, len(b.Ops))
			}
			if i > 0 {
				minGap = min(minGap, b.At.Sub(bs[i-1].At))
			}
		}
		// Within any 1 s window at most 4 batches (+1 for receive jitter at the window edge).
		for i := range bs {
			n := 0
			for j := i; j < len(bs) && bs[j].At.Sub(bs[i].At) < time.Second; j++ {
				n++
			}
			if n > 5 {
				t.Fatalf("%s: %d batches within 1 s at max_rate_hz 4", name, n)
			}
		}
		return maxB, maxOps, minGap
	}
	bs := f.Batches()
	maxB, maxOps, minGap := check("entities", bs)
	for i, fc := range cells {
		total := 0.0
		waitFor(t, 10*time.Second, "cells total 50000", func() bool {
			total = 0
			for _, c := range fc.Cells() {
				total += c[2]
			}
			return int(total) == 50000
		})
		check(fmt.Sprintf("cells %d", i), fc.Batches())
	}
	cb, cops, cgap := check("cells 0", cells[0].Batches())
	sample()
	heapMB := float64(peakHeap) / (1 << 20)
	t.Logf("50,000-vessel burst, max_rate_hz 4, 13 watches: AOI entities watch %d vessels in %d batches, max %d bytes, max %d ops, min gap %v; cells watch 0 %d batches, max %d bytes, max %d ops, min gap %v; 8 world entity watches ended with too_many; max queued ops %d; peak watch memory %.1f MiB (budget %d MiB); peak test-process heap %.1f MiB (before %.1f MiB)",
		want, len(bs), maxB, maxOps, minGap.Round(time.Millisecond), len(cells[0].Batches()), cb, cops, cgap.Round(time.Millisecond), maxPend,
		float64(peakWatchBytes)/(1<<20), DefaultWatchMemBudget>>20, heapMB, float64(before.HeapAlloc)/(1<<20))
	if maxPend > 5000 {
		t.Fatalf("queue %d exceeds the layer cap", maxPend)
	}
	if peakWatchBytes > DefaultWatchMemBudget {
		t.Fatalf("watch memory %d over the budget", peakWatchBytes)
	}
	if heapMB > 256 { // the pod limit is 512Mi; this process also holds 5 client followers
		t.Fatalf("heap %.1f MiB", heapMB)
	}
}

// Dead watch: killing the layer service makes the client re-snapshot within the dead-watch
// time of the last batch (heartbeat 0.5 s and dead time 3 s here; 10 s and 30 s in the
// contract; A5_LONG=1 runs the contract timers). The restarted service serves it again.
func TestA5DeadWatch(t *testing.T) {
	hb, dead := 500*time.Millisecond, 3*time.Second
	if os.Getenv("A5_LONG") == "1" {
		hb, dead = 10*time.Second, 30*time.Second
	}
	srv := startServer(t, "")
	nc, js := connect(t, srv.url)
	setupBus(t, js, bucketTTL{})
	publishEnts(t, js, drones(50), "")
	cfg := testConfig([]Def{fleetLayer()})
	cfg.Heartbeat = hb
	r := startService(t, srv.url, cfg)
	r.waitReady(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	f := &lc.Follower{NC: nc, Req: map[string]any{"layer": "fleet", "aoi": nil, "zoom": 12, "lease_s": 60, "max_rate_hz": 4}, DeadAfter: dead}
	go func() { _ = f.Run(ctx) }()
	waitFor(t, 5*time.Second, "loaded", func() bool { return len(f.Entities()) == 50 })
	waitFor(t, hb*3, "a heartbeat", func() bool {
		for _, b := range f.Batches() {
			if len(b.Ops) == 0 {
				return true
			}
		}
		return false
	})
	r.stop()
	waitFor(t, dead+5*time.Second, "re-snapshot attempt", func() bool { n, _ := f.Resnapshots(); return n >= 1 })
	bs := f.Batches()
	since := time.Since(bs[len(bs)-1].At)
	n, why := f.Resnapshots()
	t.Logf("service killed; client gave up the watch %v after the last batch (dead time %v): %v", since.Round(10*time.Millisecond), dead, why[:1])
	if since > dead+time.Second || !strings.HasPrefix(why[0], "dead watch") {
		t.Fatalf("detection too late or wrong reason (%v, %v)", since, why)
	}
	cfg.Holder = "test-restart"
	r2 := startService(t, srv.url, cfg)
	r2.waitReady(t)
	publishEnts(t, js, []ent{{typ: "drone", id: "sim-new", seq: 1, lat: 40.6, lon: -74, live: true, sim: true}}, "")
	waitFor(t, 20*time.Second, "client recovered after restart", func() bool { _, ok := compareToReference(t, js, f, "drone"); return ok && len(f.Entities()) == 51 })
	_ = n
}

// Site status in fleet-state: the own site is written; a site that publishes
// tel.fleet.node.{site}.status appears connected and vanishes when it stops.
func TestA5FleetState(t *testing.T) {
	srv := startServer(t, "")
	nc, js := connect(t, srv.url)
	setupBus(t, js, bucketTTL{fleet: time.Second})
	r := startService(t, srv.url, testConfig([]Def{fleetLayer()}))
	r.waitReady(t)
	kv, _ := js.KeyValue(context.Background(), "fleet-state")
	get := func(site string) *SiteStatus {
		e, err := kv.Get(context.Background(), site)
		if err != nil {
			return nil
		}
		var s SiteStatus
		_ = json.Unmarshal(e.Value(), &s)
		return &s
	}
	waitFor(t, 3*time.Second, "hub site", func() bool { s := get("hub"); return s != nil && s.Connected })
	stop := make(chan struct{})
	go func() {
		tk := time.NewTicker(200 * time.Millisecond)
		defer tk.Stop()
		for {
			select {
			case <-stop:
				return
			case <-tk.C:
				_ = nc.Publish("tel.fleet.node.site-mac.status", []byte(`{"v":1}`))
			}
		}
	}()
	waitFor(t, 3*time.Second, "site-mac connected", func() bool { s := get("site-mac"); return s != nil && s.Connected && s.V == 1 })
	close(stop)
	stopped := time.Now()
	waitFor(t, 10*time.Second, "site-mac key gone", func() bool { return get("site-mac") == nil })
	t.Logf("site-mac stopped reporting; fleet-state key gone after %v (site timeout 2 s, bucket TTL 1 s; 30 s in the CR)", time.Since(stopped).Round(100*time.Millisecond))
	if get("hub") == nil {
		t.Fatal("own site key missing")
	}
}

// Single writer: two instances, one lease; only the holder serves and writes; the standby
// takes over when the holder stops.
func TestA5SingleWriter(t *testing.T) {
	srv := startServer(t, "")
	nc, js := connect(t, srv.url)
	setupBus(t, js, bucketTTL{})
	publishEnts(t, js, drones(20), "")
	c1 := testConfig([]Def{fleetLayer()})
	c1.Holder = "one"
	c2 := c1
	c2.Holder = "two"
	r1 := startService(t, srv.url, c1)
	r1.waitReady(t)
	r2 := startService(t, srv.url, c2)
	time.Sleep(time.Second)
	if r2.svc.Ready() {
		t.Fatal("standby became ready while the lease is held")
	}
	kv, _ := js.KeyValue(context.Background(), "world-state")
	e, _ := kv.Get(context.Background(), "_lease.layer-service")
	t.Logf("lease held: %s", e.Value())
	r1.stop()
	took := time.Now()
	r2.waitReady(t)
	t.Logf("standby took over %v after the holder stopped", time.Since(took).Round(10*time.Millisecond))
	snap, err := lc.GetSnapshot(nc, map[string]any{"v": 1, "layer": "fleet"}, 5*time.Second)
	if err != nil || len(snap.Items) != 20 {
		t.Fatalf("after takeover: %v", err)
	}
}

// Paging under the 1 MiB max_payload: full entities of about 3 KiB each come back in chunks of
// at most 512 KiB; large layers refuse `fields: full` above 200 items.
func TestA5PagingFullFields(t *testing.T) {
	srv := startServer(t, "")
	nc, js := connect(t, srv.url)
	setupBus(t, js, bucketTTL{})
	fx := drones(600)
	for i := range fx {
		fx[i].pad = 3000
	}
	publishEnts(t, js, fx, "")
	r := startService(t, srv.url, testConfig([]Def{fleetLayer(), seaLayer()}))
	r.waitReady(t)
	snap, err := lc.GetSnapshot(nc, map[string]any{"v": 1, "layer": "fleet", "fields": "full"}, 5*time.Second)
	if err != nil {
		t.Fatal(err)
	}
	maxB := 0
	for _, c := range snap.Chunks {
		maxB = max(maxB, c.Bytes)
		if c.Bytes > MaxChunkBytes || int64(c.Bytes) >= nc.MaxPayload() {
			t.Fatalf("chunk %d: %d bytes (max_payload %d)", c.Part, c.Bytes, nc.MaxPayload())
		}
	}
	var full map[string]any
	_ = json.Unmarshal(snap.Items[0], &full)
	t.Logf("fields full: %d entities of ~%d bytes in %d chunks, largest %d bytes, server max_payload %d", len(snap.Items), len(snap.Items[0]), len(snap.Chunks), maxB, nc.MaxPayload())
	if len(snap.Items) != 600 || len(snap.Chunks) < 4 || full["meta"] == nil {
		t.Fatalf("items %d chunks %d", len(snap.Items), len(snap.Chunks))
	}
	_, err = lc.GetSnapshot(nc, map[string]any{"v": 1, "layer": "sea", "aoi": map[string]any{"bbox": []float64{-1, -1, 1, 1}}, "fields": "full", "max": 500}, 5*time.Second)
	if !lc.IsCode(err, ErrTooMany) {
		t.Fatalf("full on large layer: %v", err)
	}
	// list and meta (Hermes endpoints).
	m, err := nc.Request("svc.layer.list", []byte(`{"v":1}`), 2*time.Second)
	if err != nil || !strings.Contains(string(m.Data), `"layer":"fleet"`) {
		t.Fatalf("list: %v", err)
	}
	waitFor(t, 3*time.Second, "meta", func() bool {
		m, err := nc.Request("svc.layer.meta", []byte(`{"v":1,"layer":"fleet"}`), 2*time.Second)
		return err == nil && strings.Contains(string(m.Data), `"count":600`)
	})
}
