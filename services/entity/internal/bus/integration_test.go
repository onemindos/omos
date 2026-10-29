package bus

// Integration tests against a disposable nats-server (binary from $NATS_SERVER_BIN or PATH;
// skipped when absent). Streams and the lease bucket are created here with the settings of the
// proposed CRs (platform/nats/streams/ent.yaml, kv-entity-correlation.yaml). The service itself
// never creates streams or buckets.

import (
	"bufio"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log/slog"
	"math"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"

	"github.com/onemindos/omos/services/entity/internal/engine"
	"github.com/onemindos/omos/services/entity/internal/model"
)

var update = flag.Bool("update", false, "rewrite testdata/golden/entities.json")

const t0 = int64(1793275200000)

// fixedNow is the fake engine clock: one minute after the first fixture.
func fixedNow() int64 { return t0 + 60000 }

type server struct {
	url string
	cmd *exec.Cmd
}

func startServer(t *testing.T) *server {
	t.Helper()
	bin := os.Getenv("NATS_SERVER_BIN")
	if bin == "" {
		var err error
		if bin, err = exec.LookPath("nats-server"); err != nil {
			t.Skip("nats-server not found; set NATS_SERVER_BIN")
		}
	}
	root := os.Getenv("NATS_TEST_STORE_ROOT")
	var dir string
	if root != "" {
		_ = os.MkdirAll(root, 0o755)
		d, err := os.MkdirTemp(root, "js-")
		if err != nil {
			t.Fatal(err)
		}
		dir = d
		t.Cleanup(func() { _ = os.RemoveAll(d) })
	} else {
		dir = t.TempDir()
	}
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	port := l.Addr().(*net.TCPAddr).Port
	_ = l.Close()
	cmd := exec.Command(bin, "-js", "-a", "127.0.0.1", "-p", fmt.Sprint(port), "-sd", dir)
	cmd.Stdout, cmd.Stderr = io.Discard, io.Discard
	if err := cmd.Start(); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = cmd.Process.Kill(); _, _ = cmd.Process.Wait() })
	url := fmt.Sprintf("nats://127.0.0.1:%d", port)
	for i := 0; i < 100; i++ {
		if nc, err := nats.Connect(url); err == nil {
			nc.Close()
			return &server{url: url, cmd: cmd}
		}
		time.Sleep(50 * time.Millisecond)
	}
	t.Fatal("nats-server did not start")
	return nil
}

func connect(t *testing.T, s *server) (*nats.Conn, jetstream.JetStream) {
	t.Helper()
	nc, err := nats.Connect(s.url)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(nc.Close)
	js, _ := jetstream.New(nc)
	return nc, js
}

// setupStreams mirrors the CR settings: input streams as in platform/nats/streams/*.yaml and
// ENT as proposed in ent.yaml (BUS-DESIGN §4.1).
func setupStreams(t *testing.T, js jetstream.JetStream) {
	t.Helper()
	ctx := context.Background()
	for name, subj := range map[string]string{"COT": "tak.cot.>", "TEL": "tel.fleet.>", "WORLD": "evt.world.>",
		"INTEL": "intel.>", "SYS": "evt.sys.>"} {
		if _, err := js.CreateStream(ctx, jetstream.StreamConfig{Name: name, Subjects: []string{subj},
			Storage: jetstream.FileStorage, Duplicates: 2 * time.Minute}); err != nil {
			t.Fatal(err)
		}
	}
	if _, err := js.CreateStream(ctx, jetstream.StreamConfig{
		Name: "ENT", Subjects: []string{"ent.*.*.state"}, Storage: jetstream.FileStorage,
		Retention: jetstream.LimitsPolicy, Discard: jetstream.DiscardOld, MaxMsgsPerSubject: 1,
		MaxAge: 7 * 24 * time.Hour, MaxBytes: 1 << 30, Duplicates: 2 * time.Minute,
		AllowDirect: true, AllowMsgTTL: true, SubjectDeleteMarkerTTL: 15 * time.Minute,
	}); err != nil {
		t.Fatal(err)
	}
	if _, err := js.CreateKeyValue(ctx, jetstream.KeyValueConfig{Bucket: "entity-correlation", History: 5}); err != nil {
		t.Fatal(err)
	}
}

type fixture struct {
	Subject string `json:"subject"`
	Data    string `json:"data"`
}

func publishFixtures(t *testing.T, js jetstream.JetStream) int {
	t.Helper()
	f, err := os.Open("../../testdata/fixtures/replay.jsonl")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	sc := bufio.NewScanner(f)
	sc.Buffer(make([]byte, 1<<20), 1<<20)
	n := 0
	for sc.Scan() {
		var fx fixture
		if err := json.Unmarshal(sc.Bytes(), &fx); err != nil {
			t.Fatal(err)
		}
		if _, err := js.Publish(context.Background(), fx.Subject, []byte(fx.Data)); err != nil {
			t.Fatalf("publish %s: %v", fx.Subject, err)
		}
		n++
	}
	return n
}

func testConfig(holder string) Config {
	return Config{
		Inputs: DefaultInputs(), MsgTTL: true, StartLookback: 0, Idle: 20 * time.Millisecond,
		SweepEvery: time.Hour, AckWait: 2 * time.Second, LeaseBucket: "entity-correlation",
		LeaseKey: "lease.entity-service", LeaseTTL: 2 * time.Second, Holder: holder,
	}
}

func newService(t *testing.T, nc *nats.Conn, cfg Config) (*Service, *engine.Engine) {
	t.Helper()
	eng := engine.New(engine.DefaultConfig(), fixedNow)
	svc, err := New(nc, eng, cfg, slog.New(slog.NewTextHandler(io.Discard, nil)))
	if err != nil {
		t.Fatal(err)
	}
	return svc, eng
}

// drained waits until every input consumer has nothing pending and nothing unacked.
func drained(t *testing.T, js jetstream.JetStream, prefix string, timeout time.Duration) {
	t.Helper()
	ctx := context.Background()
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		done, seen := true, 0
		for _, in := range []string{"COT", "TEL", "WORLD", "INTEL"} {
			c, err := js.Consumer(ctx, in, prefix+"-"+strings.ToLower(in))
			if err != nil {
				done = false
				break
			}
			info, err := c.Info(ctx)
			if err != nil || info.NumPending > 0 || info.NumAckPending > 0 {
				done = false
				break
			}
			seen++
		}
		if done && seen == 4 {
			time.Sleep(100 * time.Millisecond)
			return
		}
		time.Sleep(50 * time.Millisecond)
	}
	t.Fatal("consumers did not drain")
}

// runUntilDrained starts the service, waits for the inputs to drain and stops it.
func runUntilDrained(t *testing.T, nc *nats.Conn, js jetstream.JetStream, cfg Config) *engine.Engine {
	t.Helper()
	svc, eng := newService(t, nc, cfg)
	ctx, cancel := context.WithCancel(context.Background())
	errc := make(chan error, 1)
	go func() { errc <- svc.Run(ctx) }()
	for i := 0; !svc.Ready(); i++ {
		if i > 400 {
			t.Fatal("service never became ready")
		}
		select {
		case err := <-errc:
			t.Fatalf("run ended early: %v", err)
		case <-time.After(25 * time.Millisecond):
		}
	}
	drained(t, js, "entity-svc", 20*time.Second)
	cancel()
	if err := <-errc; err != nil {
		t.Fatalf("run: %v", err)
	}
	return eng
}

// entSnapshot returns the last record of every ENT subject, sorted by subject.
func entSnapshot(t *testing.T, js jetstream.JetStream) ([]map[string]any, uint64) {
	t.Helper()
	ctx := context.Background()
	st, err := js.Stream(ctx, "ENT")
	if err != nil {
		t.Fatal(err)
	}
	info, err := st.Info(ctx, jetstream.WithSubjectFilter("ent.>"))
	if err != nil {
		t.Fatal(err)
	}
	subs := make([]string, 0, len(info.State.Subjects))
	for s := range info.State.Subjects {
		subs = append(subs, s)
	}
	sort.Strings(subs)
	var out []map[string]any
	for _, s := range subs {
		m, err := st.GetLastMsgForSubject(ctx, s)
		if err != nil {
			t.Fatal(err)
		}
		var rec map[string]any
		if err := json.Unmarshal(m.Data, &rec); err != nil {
			t.Fatal(err)
		}
		if m.Header.Get("Nats-Msg-Id") != fmt.Sprintf("%s:%v", rec["entity_id"], rec["seq"]) {
			t.Fatalf("%s: Nats-Msg-Id %q", s, m.Header.Get("Nats-Msg-Id"))
		}
		if m.Header.Get("Nats-TTL") == "" || m.Header.Get("OneMind-Schema") != "entity.v1" {
			t.Fatalf("%s: headers %v", s, m.Header)
		}
		out = append(out, rec)
	}
	return out, info.State.LastSeq
}

// approxEqual compares decoded JSON with a float tolerance (FMA on arm64 vs amd64).
func approxEqual(a, b any) bool {
	switch x := a.(type) {
	case float64:
		y, ok := b.(float64)
		return ok && math.Abs(x-y) <= 1e-6*math.Max(1, math.Abs(x))
	case map[string]any:
		y, ok := b.(map[string]any)
		if !ok || len(x) != len(y) {
			return false
		}
		for k, v := range x {
			if !approxEqual(v, y[k]) {
				return false
			}
		}
		return true
	case []any:
		y, ok := b.([]any)
		if !ok || len(x) != len(y) {
			return false
		}
		for i := range x {
			if !approxEqual(x[i], y[i]) {
				return false
			}
		}
		return true
	}
	return reflect.DeepEqual(a, b)
}

func normalise(v any) any {
	b, _ := json.Marshal(v)
	var out any
	_ = json.Unmarshal(b, &out)
	return out
}

type recorder struct {
	mu   sync.Mutex
	msgs []*nats.Msg
}

func record(t *testing.T, nc *nats.Conn, subject string) *recorder {
	t.Helper()
	r := &recorder{}
	sub, err := nc.Subscribe(subject, func(m *nats.Msg) { r.mu.Lock(); r.msgs = append(r.msgs, m); r.mu.Unlock() })
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = sub.Unsubscribe() })
	_ = nc.Flush()
	return r
}

func (r *recorder) all() []*nats.Msg {
	r.mu.Lock()
	defer r.mu.Unlock()
	return append([]*nats.Msg(nil), r.msgs...)
}

func byID(recs []map[string]any) map[string]map[string]any {
	out := map[string]map[string]any{}
	for _, r := range recs {
		out[r["entity_id"].(string)] = r
	}
	return out
}

// TestGoldenReplay replays the recorded-format fixtures through COT, TEL, WORLD and INTEL and
// compares the resulting ENT state with testdata/golden/entities.json. It then rebuilds from
// scratch (purge ENT, drop the durables, new process) and requires the same entities
// (FEATURES A4: "a replay of COT and INTEL from scratch rebuilds the same entities").
func TestGoldenReplay(t *testing.T) {
	s := startServer(t)
	nc, js := connect(t, s)
	setupStreams(t, js)
	publishFixtures(t, js)
	rejected := record(t, nc, "evt.sys.entity.rejected")

	runUntilDrained(t, nc, js, testConfig("golden-1"))
	snap, _ := entSnapshot(t, js)

	golden := filepath.Join("..", "..", "testdata", "golden", "entities.json")
	got, _ := json.MarshalIndent(snap, "", "  ")
	if *update {
		if err := os.WriteFile(golden, append(got, '\n'), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	want, err := os.ReadFile(golden)
	if err != nil {
		t.Fatalf("read golden (run with -update to create): %v", err)
	}
	var wantV any
	_ = json.Unmarshal(want, &wantV)
	if !approxEqual(normalise(snap), wantV) {
		t.Fatalf("ENT state differs from %s; got:\n%s", golden, got)
	}

	ents := byID(snap)
	// A4 test 1: the same aircraft from two feeds is one entity with two sources.
	ac := ents["icao-a1b2c3"]
	if ac == nil || len(ac["sources"].([]any)) != 2 {
		t.Fatalf("icao-a1b2c3: %v", ac)
	}
	// A4 test 2: a ship's position updates in place (one subject, seq advanced).
	v := ents["mmsi-367000001"]
	if v == nil || v["seq"].(float64) < 3 || math.Abs(v["lat"].(float64)-40.601) > 0.0005 {
		t.Fatalf("vessel: %v", v)
	}
	// A4 test 4 / SL9: the hostile person is dead-lettered and never stored.
	if _, ok := ents["tak-android-badf00d"]; ok {
		t.Fatal("hostile person in ENT")
	}
	time.Sleep(100 * time.Millisecond)
	reasons := map[string]bool{}
	for _, m := range rejected.all() {
		var rec map[string]any
		_ = json.Unmarshal(m.Data, &rec)
		reasons[fmt.Sprint(rec["reason"])] = true
		if rec["reason"] == engine.RejectPersonHostile && rec["source_id"] != "ANDROID-badf00d" {
			t.Fatalf("wrong reject %v", rec)
		}
	}
	if !reasons[engine.RejectPersonHostile] || len(reasons) != 2 {
		t.Fatalf("dead letters: %v", reasons)
	}
	// Echo filter: the projector's own CoT did not create an entity.
	for id := range ents {
		if strings.HasPrefix(id, "tak-om") {
			t.Fatalf("echo became %s", id)
		}
	}

	// Rebuild from scratch: a new ENT (a purge keeps the stream's Nats-Msg-Id window, so a
	// rebuild inside the 2 min duplicate window would be deduplicated away; see README).
	ctx := context.Background()
	if err := js.DeleteStream(ctx, "ENT"); err != nil {
		t.Fatal(err)
	}
	if err := js.DeleteKeyValue(ctx, "entity-correlation"); err != nil {
		t.Fatal(err)
	}
	for _, n := range []string{"COT", "TEL", "WORLD", "INTEL", "SYS"} {
		if err := js.DeleteStream(ctx, n); err != nil {
			t.Fatal(err)
		}
	}
	setupStreams(t, js)
	publishFixtures(t, js)
	runUntilDrained(t, nc, js, testConfig("golden-2"))
	rebuilt, _ := entSnapshot(t, js)
	if !approxEqual(normalise(rebuilt), normalise(snap)) {
		r, _ := json.MarshalIndent(rebuilt, "", "  ")
		t.Fatalf("rebuild differs:\n%s", r)
	}
}

type idSeq struct {
	id  string
	seq uint64
}

func coreSeqs(t *testing.T, msgs []*nats.Msg) []idSeq {
	t.Helper()
	var out []idSeq
	for _, m := range msgs {
		var e model.Entity
		if err := json.Unmarshal(m.Data, &e); err != nil {
			t.Fatal(err)
		}
		out = append(out, idSeq{e.EntityID, e.Seq})
	}
	return out
}

// TestRestartNeverDuplicates crashes the service right after an ENT publish and before the
// input ack, restarts it (lease takeover, restore from ENT, redelivery), restarts it again
// idle, then replays every input a second time (entity.md §10 test 2). No (entity_id, seq)
// is ever published twice, seq only increases, and the idle restart and the second replay add
// nothing to ENT.
func TestRestartNeverDuplicates(t *testing.T) {
	s := startServer(t)
	nc, js := connect(t, s)
	setupStreams(t, js)
	publishFixtures(t, js)
	core := record(t, nc, "ent.>")

	cfgA := testConfig("crash-a")
	cfgA.crashAfter = 5
	svcA, _ := newService(t, nc, cfgA)
	if err := svcA.Run(context.Background()); err != errCrash {
		t.Fatalf("want simulated crash, got %v", err)
	}
	_ = nc.Flush()
	time.Sleep(200 * time.Millisecond)
	afterCrash := len(core.all())

	runUntilDrained(t, nc, js, testConfig("restart-b")) // waits out the crashed lease
	time.Sleep(100 * time.Millisecond)
	seen := map[idSeq]bool{}
	last := map[string]uint64{}
	for _, x := range coreSeqs(t, core.all()) {
		if seen[x] {
			t.Fatalf("duplicate publish %s seq %d", x.id, x.seq)
		}
		seen[x] = true
		if x.seq <= last[x.id] {
			t.Fatalf("%s: seq %d after %d", x.id, x.seq, last[x.id])
		}
		last[x.id] = x.seq
	}
	snap1, lastSeq1 := entSnapshot(t, js)
	n1 := len(core.all())
	if afterCrash != 5 || n1 <= afterCrash {
		t.Fatalf("publishes: %d before crash, %d after restart", afterCrash, n1)
	}
	for id, rec := range byID(snap1) {
		if uint64(rec["seq"].(float64)) != last[id] {
			t.Fatalf("%s: ENT seq %v, last core seq %d", id, rec["seq"], last[id])
		}
	}

	runUntilDrained(t, nc, js, testConfig("restart-c")) // idle restart
	time.Sleep(100 * time.Millisecond)
	if n := len(core.all()); n != n1 {
		t.Fatalf("idle restart published %d records", n-n1)
	}

	publishFixtures(t, js) // the same inputs again, stored again by the input streams
	runUntilDrained(t, nc, js, testConfig("replay-d"))
	time.Sleep(100 * time.Millisecond)
	snap2, lastSeq2 := entSnapshot(t, js)
	if n := len(core.all()); n != n1 || lastSeq2 != lastSeq1 || !approxEqual(normalise(snap2), normalise(snap1)) {
		t.Fatalf("second replay changed ENT: core %d->%d, ENT last seq %d->%d", n1, len(core.all()), lastSeq1, lastSeq2)
	}
}

// TestLeaseSingleActive: two contenders, one winner; the other takes over only after expiry.
func TestLeaseSingleActive(t *testing.T) {
	s := startServer(t)
	_, js := connect(t, s)
	setupStreams(t, js)
	ctx := context.Background()
	a, err := NewLease(ctx, js, "entity-correlation", "lease.entity-service", "a", time.Second)
	if err != nil {
		t.Fatal(err)
	}
	b, _ := NewLease(ctx, js, "entity-correlation", "lease.entity-service", "b", time.Second)
	if ok, err := a.TryAcquire(ctx); !ok || err != nil {
		t.Fatalf("a: %v %v", ok, err)
	}
	if ok, _ := b.TryAcquire(ctx); ok {
		t.Fatal("b acquired a held lease")
	}
	if err := a.Renew(ctx); err != nil {
		t.Fatal(err)
	}
	time.Sleep(1200 * time.Millisecond) // a stops renewing (crash)
	if ok, err := b.TryAcquire(ctx); !ok || err != nil {
		t.Fatalf("b takeover: %v %v", ok, err)
	}
	if err := a.Renew(ctx); err == nil {
		t.Fatal("a renewed a lease it lost")
	}
	if _, err := NewLease(ctx, js, "no-such-bucket", "k", "c", time.Second); err == nil {
		t.Fatal("lease created a missing bucket")
	}
}

// TestRPC answers svc.entity.get and svc.entity.query from the live state.
func TestRPC(t *testing.T) {
	s := startServer(t)
	nc, js := connect(t, s)
	setupStreams(t, js)
	publishFixtures(t, js)
	cfg := testConfig("rpc")
	cfg.ServeRPC = true
	svc, _ := newService(t, nc, cfg)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() { _ = svc.Run(ctx) }()
	drained(t, js, "entity-svc", 20*time.Second)

	m, err := nc.Request("svc.entity.get", []byte(`{"entity_id":"sim-x1"}`), 2*time.Second)
	if err != nil {
		t.Fatal(err)
	}
	var e model.Entity
	if err := json.Unmarshal(m.Data, &e); err != nil || e.EntityID != "sim-x1" {
		t.Fatalf("get: %s %v", m.Data, err)
	}
	m, err = nc.Request("svc.entity.query", []byte(`{"types":["aircraft"],"bbox":[-75,40,-73,41]}`), 2*time.Second)
	if err != nil {
		t.Fatal(err)
	}
	var q struct {
		Entities []model.Entity `json:"entities"`
	}
	_ = json.Unmarshal(m.Data, &q)
	if len(q.Entities) != 2 {
		t.Fatalf("query: %s", m.Data)
	}
	m, _ = nc.Request("svc.entity.get", []byte(`{"entity_id":"nope"}`), 2*time.Second)
	if m.Header.Get("Nats-Service-Error-Code") != "404" {
		t.Fatalf("missing entity: %v", m.Header)
	}
	var eb struct {
		V      int    `json:"v"`
		Error  string `json:"error"`
		Detail string `json:"detail"`
	}
	if err := json.Unmarshal(m.Data, &eb); err != nil || eb.V != 1 || eb.Error != "not_found" || eb.Detail == "" {
		t.Fatalf("missing entity body = %s (%v), want {v:1,error:not_found,detail}", m.Data, err)
	}
}
