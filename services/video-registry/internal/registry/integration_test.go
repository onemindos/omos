package registry

import (
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/nats-io/nats.go"
)

// fakeMTX serves /v3/paths/list from a mutable item list and checks the API credentials.
type fakeMTX struct {
	mu    sync.Mutex
	items []PathItem
	down  bool
}

func (f *fakeMTX) set(items ...PathItem) { f.mu.Lock(); f.items = items; f.mu.Unlock() }

func (f *fakeMTX) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.down {
		http.Error(w, "down", http.StatusBadGateway)
		return
	}
	if u, p, ok := r.BasicAuth(); !ok || u != "video-registry" || p != "test-only" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	_ = json.NewEncoder(w).Encode(map[string]any{"itemCount": len(f.items), "pageCount": 1, "items": f.items})
}

func live(name string) PathItem {
	now := time.Now()
	it := PathItem{Name: name, Ready: true, ReadyTime: &now, Tracks: []string{"H264"}}
	it.Source = &struct {
		Type string `json:"type"`
		ID   string `json:"id"`
	}{Type: "rtspSession"}
	return it
}

func collect(t *testing.T, nc *nats.Conn, subj string) (func(time.Duration) []Event, func()) {
	t.Helper()
	ch := make(chan Event, 32)
	sub, err := nc.Subscribe(subj, func(m *nats.Msg) {
		var e Event
		if json.Unmarshal(m.Data, &e) == nil {
			ch <- e
		}
	})
	if err != nil {
		t.Fatal(err)
	}
	_ = nc.Flush()
	next := func(wait time.Duration) []Event {
		var out []Event
		deadline := time.After(wait)
		for {
			select {
			case e := <-ch:
				out = append(out, e)
			case <-deadline:
				return out
			}
		}
	}
	return next, func() { _ = sub.Unsubscribe() }
}

func get(t *testing.T, nc *nats.Conn, entity string) Reply {
	t.Helper()
	b, _ := json.Marshal(GetRequest{Entity: entity})
	m, err := nc.Request(SubjectGet, b, 2*time.Second)
	if err != nil {
		t.Fatal(err)
	}
	var r Reply
	if err := json.Unmarshal(m.Data, &r); err != nil {
		t.Fatal(err)
	}
	return r
}

func TestRegistryLifecycle(t *testing.T) {
	nc := startNATS(t)
	kv := createBucket(t, nc)
	mtx := &fakeMTX{}
	srv := httptest.NewServer(mtx)
	defer srv.Close()
	media, err := LoadMediaConfig(repoFile(t, mediaYML))
	if err != nil {
		t.Fatal(err)
	}
	cfg := Config{Media: media, PollInterval: 50 * time.Millisecond, StaleAfter: 300 * time.Millisecond,
		Endpoints: Endpoints{Host: "media.example.ts.net"}, Logger: slog.New(slog.NewTextHandler(io.Discard, nil))}
	client := &MediaMTXClient{BaseURL: srv.URL, User: "video-registry", Password: "test-only"}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	svc, err := New(ctx, nc, client, cfg)
	if err != nil {
		t.Fatal(err)
	}
	events, stop := collect(t, nc, "evt.video.>")
	defer stop()
	done := make(chan error, 1)
	go func() { done <- svc.Run(ctx) }()

	// 1. sim-x1 goes live: one KV record linked to ent.drone.sim-x1.state and one online event.
	//    A rejected path (not a sim asset) produces nothing.
	mtx.set(live("drone/sim-x1"), live("drone/x47"))
	evs := events(600 * time.Millisecond)
	if len(evs) != 1 || evs[0].Event != "online" || evs[0].EntitySubject != "ent.drone.sim-x1.state" {
		t.Fatalf("events after publish = %+v", evs)
	}
	e, err := kv.Get(ctx, "drone.sim-x1")
	if err != nil {
		t.Fatal(err)
	}
	var rec Record
	if err := json.Unmarshal(e.Value(), &rec); err != nil {
		t.Fatal(err)
	}
	if !rec.Live || rec.EntitySubject != "ent.drone.sim-x1.state" || rec.TAKURL != "rtsp://media.example.ts.net:8554/drone/sim-x1" ||
		len(rec.Permissions.Read) == 0 || !rec.Recording.Enabled {
		t.Fatalf("record = %+v", rec)
	}
	if _, err := kv.Get(ctx, "drone.x47"); err == nil {
		t.Fatal("non-sim drone was registered")
	}
	if svc.Stats.Rejected.Load() == 0 {
		t.Error("rejected path not counted")
	}

	// 2. "Selecting the sim-x1 entity shows its live stream": the client asks by entity subject.
	r := get(t, nc, "ent.drone.sim-x1.state")
	if r.Record == nil || r.Record.StreamID != "drone/sim-x1" || r.Record.URLs.HLS == "" {
		t.Fatalf("svc.registry.video.get = %+v", r)
	}
	if r := get(t, nc, "ent.drone.sim-x9.state"); r.Error != "not_found" {
		t.Errorf("unknown entity: %+v", r)
	}
	if r := get(t, nc, "ent.ship.mmsi1.state"); r.Record != nil || r.Error == "" {
		t.Errorf("bad entity type: %+v", r)
	}
	m, err := nc.Request(SubjectList, nil, 2*time.Second)
	if err != nil {
		t.Fatal(err)
	}
	var lr Reply
	_ = json.Unmarshal(m.Data, &lr)
	if len(lr.Records) != 1 {
		t.Errorf("list = %+v", lr)
	}

	// 3. Steady state: repeated polls write nothing new.
	puts := svc.Stats.Puts.Load()
	time.Sleep(300 * time.Millisecond)
	if svc.Stats.Puts.Load() != puts {
		t.Errorf("idle polls wrote %d records", svc.Stats.Puts.Load()-puts)
	}

	// 4. The publisher stops: live:false and an offline event.
	mtx.set()
	evs = events(400 * time.Millisecond)
	if len(evs) != 1 || evs[0].Event != "offline" || evs[0].Reason != "publisher_stopped" {
		t.Fatalf("events after stop = %+v", evs)
	}
	if r := get(t, nc, "drone.sim-x1"); r.Record == nil || r.Record.Live {
		t.Fatalf("after stop = %+v", r)
	}

	// 5. Back online, then MediaMTX becomes unreachable: after StaleAfter the stream goes offline.
	mtx.set(live("drone/sim-x1"))
	if evs = events(300 * time.Millisecond); len(evs) != 1 || evs[0].Event != "online" {
		t.Fatalf("events after republish = %+v", evs)
	}
	mtx.mu.Lock()
	mtx.down = true
	mtx.mu.Unlock()
	evs = events(800 * time.Millisecond)
	if len(evs) != 1 || evs[0].Reason != "mediamtx_unreachable" {
		t.Fatalf("events while MediaMTX down = %+v", evs)
	}
	cancel()
	if err := <-done; err != nil {
		t.Fatal(err)
	}
}

// TestRestartMarksVanishedStreamsOffline: a stream that was live when the registry stopped and
// is gone when it starts again is marked offline on the first poll; one that is still live is not
// announced twice.
func TestRestartMarksVanishedStreamsOffline(t *testing.T) {
	nc := startNATS(t)
	kv := createBucket(t, nc)
	media, _ := LoadMediaConfig(repoFile(t, mediaYML))
	mtx := &fakeMTX{}
	srv := httptest.NewServer(mtx)
	defer srv.Close()
	client := &MediaMTXClient{BaseURL: srv.URL, User: "video-registry", Password: "test-only"}
	quiet := slog.New(slog.NewTextHandler(io.Discard, nil))
	cfg := Config{Media: media, Endpoints: Endpoints{Host: "media.example.ts.net"}, Logger: quiet}
	ctx := context.Background()

	mtx.set(live("drone/sim-x1"), live("cam/gate-north"))
	first, err := New(ctx, nc, client, cfg)
	if err != nil {
		t.Fatal(err)
	}
	first.PollOnce(ctx)

	events, stop := collect(t, nc, "evt.video.>")
	defer stop()
	mtx.set(live("drone/sim-x1"))
	second, err := New(ctx, nc, client, cfg)
	if err != nil {
		t.Fatal(err)
	}
	second.PollOnce(ctx)
	evs := events(300 * time.Millisecond)
	if len(evs) != 1 || evs[0].Event != "offline" || evs[0].EntitySubject != "ent.camera.gate-north.state" {
		t.Fatalf("restart events = %+v", evs)
	}
	e, err := kv.Get(ctx, "camera.gate-north")
	if err != nil {
		t.Fatal(err)
	}
	var r Record
	_ = json.Unmarshal(e.Value(), &r)
	if r.Live {
		t.Error("vanished stream still live")
	}
	if !second.Ready() {
		t.Error("not ready after a good poll")
	}
}

// TestMissingBucketFails: the service never creates OVERWATCH_VIDEO itself (the NACK CR owns it).
func TestMissingBucketFails(t *testing.T) {
	nc := startNATS(t)
	_, err := New(context.Background(), nc, &MediaMTXClient{}, Config{Logger: slog.New(slog.NewTextHandler(io.Discard, nil))})
	if err == nil {
		t.Fatal("New succeeded without the bucket")
	}
}
