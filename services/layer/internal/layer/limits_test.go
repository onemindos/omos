package layer

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/nats-io/nats.go"
)

// limitStore is a ready store with fake time, for watch-limit unit tests.
func limitStore(defs []Def, budgetBytes int64) (*Store, *time.Time) {
	now := time.Unix(1_800_000_000, 0)
	s := NewStore(defs, "hub", 5*time.Minute, 100000, func() time.Time { return now })
	s.MarkReady()
	s.SetWatchBudget(budgetBytes)
	return s, &now
}

var limitSeq uint64

func applyDrones(s *Store, from, n int, lat, lon float64) {
	for i := from; i < from+n; i++ {
		limitSeq++
		e := ent{typ: "drone", id: fmt.Sprintf("sim-%05d", i), seq: 1, lat: lat, lon: lon, live: true, sim: true}
		s.Apply("ent.drone."+e.id+".state", e.json(s.now().UnixMilli()), "", limitSeq)
	}
}

func openWatch(t *testing.T, s *Store, req map[string]any) (*Watch, string) {
	t.Helper()
	req["v"] = 1
	b, _ := json.Marshal(req)
	w, reply := s.Watch(b, "_INBOX.test", Caller{Trusted: true, Known: true, ID: "service"}, WatchOpts{MaxWatches: 100, ZoomsFor: func(int) []int { return nil }})
	var r struct {
		Error string `json:"error"`
	}
	_ = json.Unmarshal(reply, &r)
	return w, r.Error
}

func endCode(w *Watch) string {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.endCode
}

// Entities-mode watches hold at most the layer cap: over the cap at start is too_many, growth
// past the cap ends the watch with too_many (memory released), and so does a renew to an
// over-cap view.
func TestWatchCap(t *testing.T) {
	d := fleetLayer()
	d.Cap = 100
	s, now := limitStore([]Def{d}, 0)
	applyDrones(s, 0, 150, 40, -74)
	if _, code := openWatch(t, s, map[string]any{"layer": "fleet", "zoom": 12}); code != ErrTooMany {
		t.Fatalf("over cap at start: %q", code)
	}
	near := map[string]any{"bbox": []float64{-75, 39, -73, 41}}
	far := map[string]any{"bbox": []float64{10, 10, 11, 11}}
	w, code := openWatch(t, s, map[string]any{"layer": "fleet", "aoi": far, "zoom": 12})
	if w == nil || code != "" {
		t.Fatalf("watch on empty AOI: %q", code)
	}
	applyDrones(s, 1000, 60, 10.5, 10.5)
	if endCode(w) != "" || w.units.Load() != 120 {
		t.Fatalf("60 in view: end %q units %d", endCode(w), w.units.Load())
	}
	applyDrones(s, 2000, 41, 10.5, 10.5)
	if endCode(w) != ErrTooMany || w.units.Load() != 0 || s.units.Load() != 0 {
		t.Fatalf("growth past cap: end %q units %d total %d", endCode(w), w.units.Load(), s.units.Load())
	}
	s.mu.RLock()
	w.mu.Lock()
	b, alive := w.buildBatchLocked(s, *now, time.Second, time.Second)
	w.mu.Unlock()
	s.mu.RUnlock()
	if alive || b == nil || !strings.Contains(string(b.data), `"error":"too_many"`) {
		t.Fatalf("ended watch batch: alive %v %v", alive, b)
	}
	// Renew into an over-cap view.
	w2, code := openWatch(t, s, map[string]any{"layer": "fleet", "aoi": map[string]any{"bbox": []float64{20, 20, 21, 21}}, "zoom": 12})
	if w2 == nil {
		t.Fatalf("watch on empty AOI refused: %q", code)
	}
	rb, _ := json.Marshal(map[string]any{"v": 1, "watch_id": w2.id, "aoi": near, "zoom": 12})
	_, reply := s.Watch(rb, "_INBOX.test", Caller{Trusted: true, Known: true, ID: "service"}, WatchOpts{MaxWatches: 100, ZoomsFor: func(int) []int { return nil }})
	if !strings.Contains(string(reply), ErrTooMany) || endCode(w2) != ErrTooMany {
		t.Fatalf("renew over cap: %s end %q", reply, endCode(w2))
	}
}

// The service-wide budget ends the watch holding the most with resnapshot; the others keep
// running, and the total stays within the budget.
func TestWatchBudget(t *testing.T) {
	s, _ := limitStore([]Def{fleetLayer()}, 3000*UnitBytes)
	applyDrones(s, 0, 1000, 40.5, -74.5) // inside both AOIs
	applyDrones(s, 1000, 500, 45, -70)   // inside the big AOI only
	big, _ := openWatch(t, s, map[string]any{"layer": "fleet", "zoom": 12})
	small, _ := openWatch(t, s, map[string]any{"layer": "fleet", "aoi": map[string]any{"bbox": []float64{-75, 40, -74, 41}}, "zoom": 12})
	if big == nil || small == nil {
		t.Fatal("watches refused")
	}
	// big holds 1,500 + 1,500 queued enters = 3,000 units: with small's 2,000 the total is
	// over 3,000, so the largest (big) is ended.
	if endCode(big) != ErrResnapshot || endCode(small) != "" {
		t.Fatalf("budget: big %q small %q", endCode(big), endCode(small))
	}
	if got := s.units.Load(); got > 3000 || got != small.units.Load() {
		t.Fatalf("total units %d, small %d", got, small.units.Load())
	}
	// Growth while running: small passes the budget on its own and is ended too.
	applyDrones(s, 5000, 600, 40.5, -74.5)
	if endCode(small) != ErrResnapshot || s.units.Load() != 0 {
		t.Fatalf("growth: small %q total %d", endCode(small), s.units.Load())
	}
	// Pending entity ops are pointers to the shared entry, rendered at batch time.
	w, _ := openWatch(t, s, map[string]any{"layer": "fleet", "aoi": map[string]any{"bbox": []float64{10, 10, 11, 11}}, "zoom": 12})
	applyDrones(s, 9000, 1, 10.5, 10.5)
	w.mu.Lock()
	op := w.pend["drone.sim-09000"]
	w.mu.Unlock()
	if op.e == nil || op.e != s.ents["drone.sim-09000"] {
		t.Fatal("pending op does not share the store entry")
	}
}

// A cells watch takes a new diff only after the previous one was sent, so cell ops never
// queue up behind max_rate_hz.
func TestCellOpsDoNotAccumulate(t *testing.T) {
	s, now := limitStore([]Def{seaLayer()}, 0)
	b, _ := json.Marshal(map[string]any{"v": 1, "layer": "sea", "aoi": map[string]any{"bbox": []float64{-180, -85, 180, 85}}, "zoom": 3})
	w, _ := s.Watch(b, "_INBOX.t", Caller{Trusted: true, Known: true, ID: "service"}, WatchOpts{ZoomsFor: func(int) []int { return nil }})
	build := func() *batch {
		*now = now.Add(2 * time.Second)
		s.mu.RLock()
		defer s.mu.RUnlock()
		w.mu.Lock()
		defer w.mu.Unlock()
		bt, _ := w.buildBatchLocked(s, *now, time.Hour, time.Second)
		return bt
	}
	for i := 0; i < 5000; i++ {
		limitSeq++
		e := ent{typ: "vessel", id: fmt.Sprintf("mmsi%06d", i), seq: 1, lat: float64(i%140) - 70, lon: float64(i%360) - 180, live: true}
		s.Apply("ent.vessel."+e.id+".state", e.json(now.UnixMilli()), "", limitSeq)
	}
	w.mu.Lock()
	w.cellOps = [][]byte{[]byte(`{"op":"cells","z":2,"cells":[]}`)} // an unsent diff
	w.mu.Unlock()
	build()
	w.mu.Lock()
	queued, dirty := len(w.cellOps), w.cellsDirty
	w.mu.Unlock()
	if queued != 0 || !dirty {
		t.Fatalf("diff taken while one was queued: %d queued, dirty %v", queued, dirty)
	}
	if bt := build(); bt == nil || !strings.Contains(string(bt.data), `"op":"cells"`) {
		t.Fatal("no cells diff after the queue drained")
	}
	// Two watches with the same view share one grid per store version.
	w2, _ := s.Watch(b, "_INBOX.t2", Caller{Trusted: true, Known: true, ID: "service"}, WatchOpts{ZoomsFor: func(int) []int { return nil }})
	s.mu.RLock()
	w.mu.Lock()
	g1 := s.sharedCellsLocked(w)
	w.mu.Unlock()
	w2.mu.Lock()
	g2 := s.sharedCellsLocked(w2)
	w2.mu.Unlock()
	s.mu.RUnlock()
	if fmt.Sprintf("%p", g1) != fmt.Sprintf("%p", g2) {
		t.Fatal("identical views computed separate grids")
	}
}

// Site contact: bounded map, pruned by the timeout, allowlist honoured.
func TestSiteSeenBounded(t *testing.T) {
	now := time.Unix(1_800_000_000, 0)
	s := &Service{cfg: Config{MaxSites: 2, SiteTimeout: 30 * time.Second}, now: func() time.Time { return now }, siteSeen: map[string]time.Time{}}
	s.seeSite("site-a")
	s.seeSite("site-b")
	s.seeSite("site-c")
	if len(s.siteSeen) != 2 || !s.siteSeen["site-a"].Equal(now) {
		t.Fatalf("cap: %v", s.siteSeen)
	}
	now = now.Add(31 * time.Second)
	s.seeSite("site-c")
	if len(s.siteSeen) != 1 || s.siteSeen["site-c"].IsZero() {
		t.Fatalf("prune: %v", s.siteSeen)
	}
	s.cfg.SiteAllow = []string{"site-mac"}
	s.seeSite("site-x")
	s.seeSite("site-mac")
	if _, ok := s.siteSeen["site-x"]; ok || s.siteSeen["site-mac"].IsZero() {
		t.Fatalf("allowlist: %v", s.siteSeen)
	}
}

// A request without Nats-Request-Info whose reply was rewritten by a service import fails
// closed; an own-account request (own inbox) is still a trusted service.
func TestResolveImportWithoutShare(t *testing.T) {
	d, _ := LoadDirectory("")
	if c := d.Resolve(&nats.Msg{Reply: "_R_.abc.def"}); c.Known || c.Trusted || c.CanSee("CLEAR", nil) {
		t.Fatalf("import without share: %+v", c)
	}
	if c := d.Resolve(&nats.Msg{Reply: "_INBOX.abc"}); !c.Trusted {
		t.Fatalf("own account: %+v", c)
	}
}
