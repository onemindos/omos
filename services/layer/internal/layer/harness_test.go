package layer

// Test harness: a disposable nats-server (binary from $NATS_SERVER_BIN or PATH; tests skip when
// absent), the ENT stream and the world-state / fleet-state buckets created with the settings
// of their CRs (platform/nats/streams/ent.yaml, kv-world-state.yaml, kv-fleet-state.yaml), and
// a layer service with short timers. The service never creates streams or buckets.

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

type testServer struct {
	url  string
	cmd  *exec.Cmd
	port int
	dir  string
	conf string
}

func storeRoot(t *testing.T) string {
	t.Helper()
	if root := os.Getenv("NATS_TEST_STORE_ROOT"); root != "" {
		_ = os.MkdirAll(root, 0o755)
		d, err := os.MkdirTemp(root, "js-")
		if err != nil {
			t.Fatal(err)
		}
		t.Cleanup(func() { _ = os.RemoveAll(d) })
		return d
	}
	return t.TempDir()
}

func freePort(t *testing.T) int {
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer l.Close()
	return l.Addr().(*net.TCPAddr).Port
}

// startServer starts nats-server with JetStream. conf, when set, is a config file body; the
// port, store dir and listen address are always given on the command line.
func startServer(t *testing.T, conf string) *testServer {
	t.Helper()
	bin := os.Getenv("NATS_SERVER_BIN")
	if bin == "" {
		var err error
		if bin, err = exec.LookPath("nats-server"); err != nil {
			t.Skip("nats-server not found; set NATS_SERVER_BIN")
		}
	}
	s := &testServer{port: freePort(t), dir: storeRoot(t)}
	args := []string{"-js", "-a", "127.0.0.1", "-p", fmt.Sprint(s.port), "-sd", s.dir}
	if conf != "" {
		s.conf = filepath.Join(s.dir, "server.conf")
		if err := os.WriteFile(s.conf, []byte(conf), 0o600); err != nil {
			t.Fatal(err)
		}
		args = append([]string{"-c", s.conf}, args...)
	}
	s.cmd = exec.Command(bin, args...)
	s.cmd.Stdout, s.cmd.Stderr = io.Discard, io.Discard
	if err := s.cmd.Start(); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = s.cmd.Process.Kill(); _, _ = s.cmd.Process.Wait() })
	s.url = fmt.Sprintf("nats://127.0.0.1:%d", s.port)
	for i := 0; i < 100; i++ {
		c, err := net.Dial("tcp", fmt.Sprintf("127.0.0.1:%d", s.port))
		if err == nil {
			c.Close()
			return s
		}
		time.Sleep(50 * time.Millisecond)
	}
	t.Fatal("nats-server did not start")
	return nil
}

func connect(t *testing.T, url string, opts ...nats.Option) (*nats.Conn, jetstream.JetStream) {
	t.Helper()
	nc, err := nats.Connect(url, opts...)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(nc.Close)
	js, _ := jetstream.New(nc)
	return nc, js
}

type bucketTTL struct{ world, fleet time.Duration }

// setupBus creates ENT, TEL and the two buckets as their CRs define them. ttl overrides the
// bucket TTLs (the layer-down and site tests shorten them to keep the run short).
func setupBus(t *testing.T, js jetstream.JetStream, ttl bucketTTL) {
	t.Helper()
	ctx := context.Background()
	if _, err := js.CreateStream(ctx, jetstream.StreamConfig{
		Name: "ENT", Subjects: []string{"ent.*.*.state"}, Storage: jetstream.FileStorage,
		Retention: jetstream.LimitsPolicy, Discard: jetstream.DiscardOld, MaxMsgsPerSubject: 1,
		MaxAge: 7 * 24 * time.Hour, MaxBytes: 1 << 30, Duplicates: 2 * time.Minute,
		AllowDirect: true, AllowMsgTTL: true, SubjectDeleteMarkerTTL: 15 * time.Minute,
	}); err != nil {
		t.Fatal(err)
	}
	if _, err := js.CreateStream(ctx, jetstream.StreamConfig{Name: "TEL", Subjects: []string{"tel.feed.>", "tel.fleet.>"},
		Storage: jetstream.FileStorage, MaxMsgsPerSubject: 2000, AllowDirect: true}); err != nil {
		t.Fatal(err)
	}
	if ttl.world == 0 {
		ttl.world = 5 * time.Minute // kv-world-state.yaml
	}
	if ttl.fleet == 0 {
		ttl.fleet = 30 * time.Second // kv-fleet-state.yaml
	}
	if _, err := js.CreateKeyValue(ctx, jetstream.KeyValueConfig{Bucket: "world-state", History: 1, TTL: ttl.world, Storage: jetstream.FileStorage}); err != nil {
		t.Fatal(err)
	}
	if _, err := js.CreateKeyValue(ctx, jetstream.KeyValueConfig{Bucket: "fleet-state", History: 3, TTL: ttl.fleet, Storage: jetstream.FileStorage}); err != nil {
		t.Fatal(err)
	}
}

// testConfig is the service config with short timers.
func testConfig(layers []Def) Config {
	return Config{
		Layers: layers, LeaseBucket: "world-state", LeaseTTL: 1500 * time.Millisecond, Holder: "test-" + randHex(4),
		Heartbeat: 500 * time.Millisecond, CellEvery: 200 * time.Millisecond, ManifestTick: 100 * time.Millisecond,
		ManifestEvery: time.Second, CountEvery: 200 * time.Millisecond, CellsEvery: 500 * time.Millisecond,
		SiteEvery: 300 * time.Millisecond, SiteTimeout: 2 * time.Second, SweepEvery: 200 * time.Millisecond,
	}
}

type runningSvc struct {
	svc    *Service
	cancel context.CancelFunc
	done   chan error
	once   sync.Once
}

func startService(t *testing.T, url string, cfg Config, opts ...nats.Option) *runningSvc {
	t.Helper()
	nc, err := nats.Connect(url, opts...)
	if err != nil {
		t.Fatal(err)
	}
	svc, err := New(nc, cfg, slog.New(slog.NewTextHandler(io.Discard, nil)))
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	r := &runningSvc{svc: svc, cancel: cancel, done: make(chan error, 1)}
	go func() { r.done <- svc.Run(ctx); nc.Close() }()
	t.Cleanup(r.stop)
	return r
}

func (r *runningSvc) stop() { r.once.Do(r.stopOnce) }

func (r *runningSvc) stopOnce() {
	r.cancel()
	select {
	case <-r.done:
	case <-time.After(10 * time.Second):
		buf := make([]byte, 1<<20)
		n := runtime.Stack(buf, true)
		panic("service did not stop:\n" + string(buf[:n]))
	}
}

func (r *runningSvc) waitReady(t *testing.T) {
	t.Helper()
	waitFor(t, 20*time.Second, "service ready", r.svc.Ready)
}

func waitFor(t *testing.T, d time.Duration, what string, f func() bool) {
	t.Helper()
	end := time.Now().Add(d)
	for time.Now().Before(end) {
		if f() {
			return
		}
		time.Sleep(20 * time.Millisecond)
	}
	t.Fatalf("timed out waiting for %s", what)
}

func randHex(n int) string {
	b := make([]byte, n)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}

// ent is a fixture entity (entity.v1 fields the layer service reads, plus the envelope).
type ent struct {
	typ, id    string
	seq        uint64
	lat, lon   float64
	live       bool
	expiryT    *int64
	tlp        string
	releasable []string
	sim        bool
	cot        bool
	pad        int // extra bytes in meta (paging tests)
	geo        map[string]any
}

func (e ent) json(now int64) []byte {
	tlp := e.tlp
	if tlp == "" {
		tlp = "CLEAR"
	}
	rel := e.releasable
	if rel == nil {
		rel = []string{}
	}
	dt := "geo"
	if e.cot {
		dt = "cot"
	}
	m := map[string]any{
		"v": 1, "uid": e.id, "type": "a-f-A", "entity_id": e.id, "entity_type": e.typ,
		"lat": e.lat, "lon": e.lon, "hae": 120.5, "course": 90.0, "speed": 12.0, "callsign": e.id,
		"t": now, "stale": 30, "seq": e.seq, "live": e.live, "lifecycle": map[bool]string{true: "active", false: "deleted"}[e.live],
		"expiry_t": e.expiryT, "cot_type": "a-f-A", "disposition": "friendly",
		"indicators": map[string]any{"simulated": e.sim}, "marking": map[string]any{"tlp": tlp, "releasable_to": rel},
		"sources": []map[string]any{{"data_type": dt}},
	}
	if e.pad > 0 {
		b := make([]byte, e.pad)
		for i := range b {
			b[i] = 'x'
		}
		m["meta"] = map[string]any{"pad": string(b)}
	}
	if e.geo != nil {
		m["geo"] = e.geo
		m["template"] = "geo"
	}
	b, _ := json.Marshal(m)
	return b
}

var pubSeq atomic.Uint64

// publishEnts writes entities to ENT with a JetStream ack. ttl, when set, is sent as Nats-TTL.
func publishEnts(t *testing.T, js jetstream.JetStream, ents []ent, ttl string) {
	t.Helper()
	ctx := context.Background()
	now := time.Now().UnixMilli()
	var futs []jetstream.PubAckFuture
	for _, e := range ents {
		m := nats.NewMsg("ent." + e.typ + "." + e.id + ".state")
		m.Data = e.json(now)
		m.Header.Set("Nats-Msg-Id", fmt.Sprintf("%s:%d:%d", e.id, e.seq, pubSeq.Add(1)))
		if ttl != "" {
			m.Header.Set("Nats-TTL", ttl)
		}
		f, err := js.PublishMsgAsync(m)
		if err != nil {
			t.Fatal(err)
		}
		futs = append(futs, f)
		if len(futs) >= 4000 {
			waitAcks(t, ctx, js, futs)
			futs = futs[:0]
		}
	}
	waitAcks(t, ctx, js, futs)
}

func waitAcks(t *testing.T, ctx context.Context, js jetstream.JetStream, futs []jetstream.PubAckFuture) {
	t.Helper()
	select {
	case <-js.PublishAsyncComplete():
	case <-time.After(30 * time.Second):
		t.Fatal("publish acks timed out")
	}
	for _, f := range futs {
		select {
		case <-f.Ok():
		case err := <-f.Err():
			t.Fatal(err)
		}
	}
}

func getManifest(t *testing.T, js jetstream.JetStream, layer string) (*Manifest, error) {
	t.Helper()
	kv, err := js.KeyValue(context.Background(), "world-state")
	if err != nil {
		t.Fatal(err)
	}
	e, err := kv.Get(context.Background(), layer+".meta")
	if err != nil {
		return nil, err
	}
	var m Manifest
	if err := json.Unmarshal(e.Value(), &m); err != nil {
		t.Fatal(err)
	}
	return &m, nil
}

// referenceState reads ENT from the start (last message per subject) and builds the entity
// set by the same rules a client uses: live records are held, live:false and markers remove.
func referenceState(t *testing.T, js jetstream.JetStream, typ string) map[string]uint64 {
	t.Helper()
	ctx := context.Background()
	oc, err := js.OrderedConsumer(ctx, "ENT", jetstream.OrderedConsumerConfig{FilterSubjects: []string{"ent." + typ + ".>"}, DeliverPolicy: jetstream.DeliverLastPerSubjectPolicy})
	if err != nil {
		t.Fatal(err)
	}
	out := map[string]uint64{}
	st, _ := js.Stream(ctx, "ENT")
	info, _ := st.Info(ctx, jetstream.WithSubjectFilter("ent."+typ+".>"))
	if info.State.Msgs == 0 {
		return out
	}
	for {
		b, err := oc.Fetch(1000, jetstream.FetchMaxWait(time.Second))
		if err != nil {
			t.Fatal(err)
		}
		got, done := 0, false
		for m := range b.Messages() {
			got++
			meta, _ := m.Metadata()
			if len(m.Data()) > 0 && m.Headers().Get("Nats-Marker-Reason") == "" {
				var e struct {
					ID      string `json:"entity_id"`
					Seq     uint64 `json:"seq"`
					Live    bool   `json:"live"`
					ExpiryT *int64 `json:"expiry_t"`
				}
				_ = json.Unmarshal(m.Data(), &e)
				if e.Live && (e.ExpiryT == nil || *e.ExpiryT > time.Now().UnixMilli()) {
					out[e.ID] = e.Seq
				}
			}
			if meta.NumPending == 0 {
				done = true
			}
		}
		if done || got == 0 {
			return out
		}
	}
}

func fleetLayer() Def {
	return Def{ID: "fleet", Title: "Fleet", Tier: TierSmall, EntityTypes: []string{"drone", "robot", "vehicle", "node", "edge"},
		Cap: 20000, LimitS: 60, DefaultStaleS: 60, Marking: "AMBER+STRICT", Style: "fleet.default.v1"}
}

func seaLayer() Def {
	return Def{ID: "sea", Title: "Vessels", Tier: TierLarge, EntityTypes: []string{"vessel"}, Cap: 5000,
		Cells: &CellsCfg{Zooms: []int{2, 4, 6}, BelowZoom: 7}, LimitS: 1800, DefaultStaleS: 1800, Marking: "CLEAR", Style: "sea.default.v1"}
}

func airLayer() Def {
	return Def{ID: "air", Title: "Aircraft", Tier: TierLarge, EntityTypes: []string{"aircraft"}, Cap: 5000,
		Cells: &CellsCfg{Zooms: []int{2, 4, 6}, BelowZoom: 7}, LimitS: 60, DefaultStaleS: 30, Marking: "CLEAR", Style: "air.default.v1"}
}

func ptr[T any](v T) *T { return &v }
