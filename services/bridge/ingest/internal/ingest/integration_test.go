package ingest_test

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/coder/websocket"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/prometheus/client_golang/prometheus/testutil"

	"github.com/onemindos/omos-bridge/ingest/internal/ingest"
	"github.com/onemindos/omos-bridge/ingest/internal/natstest"
	"github.com/onemindos/omos-bridge/ingest/internal/parsers"
)

const (
	schemaDir  = "../../../../schemas"
	fixtureDir = "../../testdata/fixtures"
)

func fixture(t testing.TB, p string) []byte {
	b, err := os.ReadFile(filepath.Join(fixtureDir, p))
	if err != nil {
		t.Fatal(err)
	}
	return b
}

type logBuf struct {
	mu sync.Mutex
	b  strings.Builder
}

func (l *logBuf) Write(p []byte) (int, error) { l.mu.Lock(); defer l.mu.Unlock(); return l.b.Write(p) }
func (l *logBuf) String() string              { l.mu.Lock(); defer l.mu.Unlock(); return l.b.String() }

type env struct {
	srv *natstest.Server
	js  jetstream.JetStream
	cfg *ingest.Config
}

func newEnv(t *testing.T, dup time.Duration) *env {
	t.Helper()
	s := natstest.Start(t)
	_, js := s.Connect(t, "test-admin")
	natstest.Streams(t, js, dup)
	cfg, err := ingest.ParseConfig([]byte("leaseTTL: 1500ms\nstatusInterval: 200ms\n"))
	if err != nil {
		t.Fatal(err)
	}
	return &env{srv: s, js: js, cfg: cfg}
}

func (e *env) runtime(t *testing.T, instance string) (*ingest.Runtime, *logBuf) {
	t.Helper()
	nc, js := e.srv.Connect(t, "omos-ingest-"+instance)
	board := ingest.NewStatusBoard(instance, nil)
	lb := &logBuf{}
	rt := &ingest.Runtime{NC: nc, JS: js, Board: board, Metrics: ingest.NewMetrics(board), Schemas: ingest.NewSchemas(schemaDir),
		Log: slog.New(slog.NewJSONHandler(lb, nil)).With("instance", instance), Instance: instance, LeaseTTL: e.cfg.LeaseTTL.Duration,
		AckTimeout: time.Second, MaxPublishBackoff: 500 * time.Millisecond}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := rt.BindBuckets(ctx, e.cfg, true); err != nil {
		t.Fatal(err)
	}
	return rt, lb
}

func sourceCfg(id, kind, parser, url, subject, schema string) ingest.SourceConfig {
	y := fmt.Sprintf(`
sources:
  - id: %s
    kind: %s
    parser: %s
    url: %s
    subject: %s
    schema: %s
    interval: 1s
    staleAfter: 1s
    maxBackoff: 2s
    rateLimit: {every: 1ms, burst: 10}
    headers: {User-Agent: "omos-ingest-test"}
    terms: {licence: "test", url: "https://example.org/terms"}
`, id, kind, parser, url, subject, schema)
	c, err := ingest.ParseConfig([]byte(y))
	if err != nil {
		panic(err)
	}
	return c.Sources[0]
}

func nwsSource(t *testing.T, rt *ingest.Runtime, url string) *ingest.Source {
	t.Helper()
	s, err := ingest.NewSource(rt, sourceCfg("nws-alerts", ingest.KindRESTPoll, "nws-alerts", url, parsers.NWSSubject, "intel/weather-alert.v1.schema.json"))
	if err != nil {
		t.Fatal(err)
	}
	s.Fetcher().Attempts = 1
	return s
}

// upstream is an NWS-like stub: switchable body and status, ETag support, request count.
type upstream struct {
	*httptest.Server
	mu       sync.Mutex
	body     []byte
	code     int
	etag     bool
	requests atomic.Int32
}

func newUpstream(t *testing.T, body []byte) *upstream {
	u := &upstream{body: body, code: 200, etag: true}
	u.Server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		u.requests.Add(1)
		u.mu.Lock()
		defer u.mu.Unlock()
		if u.code != 200 {
			w.WriteHeader(u.code)
			return
		}
		tag := fmt.Sprintf(`W/"%x"`, len(u.body))
		if u.etag && r.Header.Get("If-None-Match") == tag {
			w.WriteHeader(http.StatusNotModified)
			return
		}
		if u.etag {
			w.Header().Set("ETag", tag)
		}
		w.Header().Set("Content-Type", "application/geo+json")
		_, _ = w.Write(u.body)
	}))
	t.Cleanup(u.Close)
	return u
}

func (u *upstream) set(body []byte, code int) {
	u.mu.Lock()
	defer u.mu.Unlock()
	if body != nil {
		u.body = body
	}
	u.code = code
}

func count(t *testing.T, js jetstream.JetStream, stream, filter string) int {
	return len(natstest.Msgs(t, js, stream, filter))
}

func waitFor(t *testing.T, d time.Duration, what string, cond func() bool) {
	t.Helper()
	deadline := time.Now().Add(d)
	for time.Now().Before(deadline) {
		if cond() {
			return
		}
		time.Sleep(20 * time.Millisecond)
	}
	t.Fatalf("timed out waiting for %s", what)
}

// ---------------------------------------------------------------- publish with ack and dedup

// Reliability: JetStream publish with server ack and Nats-Msg-Id dedup.
func TestPublishAckAndMsgIDDedup(t *testing.T) {
	e := newEnv(t, 2*time.Minute)
	rt, _ := e.runtime(t, "a")
	ctx := context.Background()
	dup1, err := rt.Publish(ctx, "x", "intel.test.x.event", "x:1", []byte(`{}`))
	if err != nil || dup1 {
		t.Fatal(err, dup1)
	}
	dup2, err := rt.Publish(ctx, "x", "intel.test.x.event", "x:1", []byte(`{}`))
	if err != nil || !dup2 {
		t.Fatal("second publish must be acked as duplicate", err)
	}
	if n := count(t, e.js, "INTEL", "intel.>"); n != 1 {
		t.Fatalf("stored %d, want 1", n)
	}
	m := natstest.Msgs(t, e.js, "INTEL", "intel.>")[0]
	if m.Headers().Get(jetstream.MsgIDHeader) != "x:1" {
		t.Fatal("Nats-Msg-Id not set")
	}
	if testutil.ToFloat64(rt.Metrics.Duplicates.WithLabelValues("x")) != 1 || testutil.ToFloat64(rt.Metrics.Published.WithLabelValues("x")) != 1 {
		t.Fatal("metrics")
	}
}

// Reliability: a publish to a subject no stream captures (or a stream not created yet)
// retries with backoff until it is acked; nothing is dropped.
func TestPublishRetriesUntilStreamExists(t *testing.T) {
	e := newEnv(t, 2*time.Minute)
	rt, _ := e.runtime(t, "a")
	done := make(chan error, 1)
	go func() {
		_, err := rt.Publish(context.Background(), "x", "det.air.visual.track1", "d:1", []byte(`{}`))
		done <- err
	}()
	time.Sleep(700 * time.Millisecond)
	if _, err := e.js.CreateStream(context.Background(), jetstream.StreamConfig{Name: "DET", Subjects: []string{"det.>"}}); err != nil {
		t.Fatal(err)
	}
	select {
	case err := <-done:
		if err != nil {
			t.Fatal(err)
		}
	case <-time.After(10 * time.Second):
		t.Fatal("publish never completed")
	}
	if testutil.ToFloat64(rt.Metrics.Retries.WithLabelValues("x")) < 1 {
		t.Fatal("expected retries")
	}
	if count(t, e.js, "DET", "det.>") != 1 {
		t.Fatal("not stored")
	}
}

// Reliability: NATS down → the publisher blocks and retries (no loss), then catches up.
func TestPublishBlocksThroughNATSRestart(t *testing.T) {
	e := newEnv(t, 2*time.Minute)
	rt, _ := e.runtime(t, "a")
	e.srv.Stop()
	done := make(chan error, 1)
	go func() {
		_, err := rt.Publish(context.Background(), "x", "intel.test.x.event", "x:outage", []byte(`{}`))
		done <- err
	}()
	time.Sleep(time.Second)
	select {
	case <-done:
		t.Fatal("publish returned while NATS was down")
	default:
	}
	e.srv.Restart()
	select {
	case err := <-done:
		if err != nil {
			t.Fatal(err)
		}
	case <-time.After(15 * time.Second):
		t.Fatal("publish did not complete after restart")
	}
	_, js := e.srv.Connect(t, "check")
	if count(t, js, "INTEL", "intel.>") != 1 {
		t.Fatal("message lost across restart")
	}
}

// Reliability: a full stream (discard new, standing in for a full disk) refuses the
// publish; the collector blocks without loss and without checkpointing, then completes.
func TestFullStreamBlocksWithoutLossOrCheckpoint(t *testing.T) {
	e := newEnv(t, 2*time.Minute)
	ctx := context.Background()
	if _, err := e.js.UpdateStream(ctx, jetstream.StreamConfig{Name: "INTEL", Subjects: []string{"intel.>"}, MaxMsgs: 2, Discard: jetstream.DiscardNew, Duplicates: 2 * time.Minute}); err != nil {
		t.Fatal(err)
	}
	up := newUpstream(t, fixture(t, "nws/active-2026-10-29.json"))
	rt, _ := e.runtime(t, "a")
	src := nwsSource(t, rt, up.URL)
	pctx, cancel := context.WithTimeout(ctx, 1500*time.Millisecond)
	_, n, err := src.PollOnce(pctx)
	cancel()
	if err == nil || n != 2 {
		t.Fatalf("want blocked after 2 stored, got n=%d err=%v", n, err)
	}
	cp, _ := rt.LoadCheckpoint(ctx, "nws-alerts")
	if cp.ETag != "" || len(cp.Seen) != 0 {
		t.Fatalf("checkpoint advanced past un-acked records: %+v", cp)
	}
	if _, err := e.js.UpdateStream(ctx, jetstream.StreamConfig{Name: "INTEL", Subjects: []string{"intel.>"}, MaxMsgs: -1, Discard: jetstream.DiscardNew, Duplicates: 2 * time.Minute}); err != nil {
		t.Fatal(err)
	}
	if _, _, err := src.PollOnce(ctx); err != nil {
		t.Fatal(err)
	}
	if got := count(t, e.js, "INTEL", "intel.>"); got != 12 {
		t.Fatalf("stored %d, want all 12 exactly once", got)
	}
}

// ---------------------------------------------------------------- checkpoints

// Reliability: REST poll end to end: publish all, then ETag gives a 304 and nothing new.
func TestRESTPollEndToEndWithETag(t *testing.T) {
	e := newEnv(t, 2*time.Minute)
	up := newUpstream(t, fixture(t, "nws/active-2026-10-29.json"))
	rt, _ := e.runtime(t, "a")
	src := nwsSource(t, rt, up.URL)
	ctx := context.Background()
	recs, n, err := src.PollOnce(ctx)
	if err != nil || recs != 12 || n != 12 {
		t.Fatalf("recs=%d new=%d err=%v", recs, n, err)
	}
	recs, n, err = src.PollOnce(ctx)
	if err != nil || recs != 0 || n != 0 {
		t.Fatalf("second poll should be 304: recs=%d new=%d err=%v", recs, n, err)
	}
	cp, _ := rt.LoadCheckpoint(ctx, "nws-alerts")
	if cp.ETag == "" || len(cp.Seen) != 12 {
		t.Fatalf("checkpoint %+v", cp)
	}
	// The next snapshot drops 2 alerts and adds 1 update: exactly one new record.
	up.set(fixture(t, "nws/active-2026-10-29-next.json"), 200)
	if _, n, err = src.PollOnce(ctx); err != nil || n != 1 {
		t.Fatalf("next snapshot new=%d err=%v", n, err)
	}
	msgs := natstest.Msgs(t, e.js, "INTEL", "intel.weather.noaa.alert")
	if len(msgs) != 13 {
		t.Fatalf("stored %d", len(msgs))
	}
	var env ingest.Envelope
	if err := json.Unmarshal(msgs[0].Data(), &env); err != nil {
		t.Fatal(err)
	}
	envSch, err := ingest.NewSchemas(schemaDir).Get("envelope/envelope.schema.json")
	if err != nil {
		t.Fatal(err)
	}
	var asMap map[string]any
	_ = json.Unmarshal(msgs[0].Data(), &asMap)
	if err := ingest.Validate(envSch, asMap); err != nil {
		t.Fatalf("published message is not a valid bus envelope: %v", err)
	}
	if env.Source != "ingest.nws-alerts" || env.Subject != parsers.NWSSubject || env.Data["type"] != "weather_alert" || env.FetchedAt == "" {
		t.Fatalf("envelope %+v", env)
	}
}

// Reliability: a restart (new process, same KV state) republishes nothing even after
// the stream's duplicate window has passed and the upstream ignores ETags.
func TestRestartBeyondDedupWindowDoesNotDuplicate(t *testing.T) {
	e := newEnv(t, 500*time.Millisecond)
	up := newUpstream(t, fixture(t, "nws/active-2026-10-29.json"))
	up.etag = false
	rt1, _ := e.runtime(t, "a")
	if _, n, err := nwsSource(t, rt1, up.URL).PollOnce(context.Background()); err != nil || n != 12 {
		t.Fatal(n, err)
	}
	time.Sleep(time.Second) // past the 500 ms duplicate window
	rt2, _ := e.runtime(t, "a-restarted")
	if _, n, err := nwsSource(t, rt2, up.URL).PollOnce(context.Background()); err != nil || n != 0 {
		t.Fatalf("restart republished %d (err %v)", n, err)
	}
	if got := count(t, e.js, "INTEL", "intel.>"); got != 12 {
		t.Fatalf("stored %d", got)
	}
}

// Reliability: checkpoint only after ack. A poll that cannot get acks leaves no checkpoint.
func TestCheckpointOnlyAfterAck(t *testing.T) {
	e := newEnv(t, 2*time.Minute)
	ctx := context.Background()
	if err := e.js.DeleteStream(ctx, "INTEL"); err != nil {
		t.Fatal(err)
	}
	up := newUpstream(t, fixture(t, "nws/active-2026-10-29.json"))
	rt, _ := e.runtime(t, "a")
	src := nwsSource(t, rt, up.URL)
	pctx, cancel := context.WithTimeout(ctx, time.Second)
	_, n, err := src.PollOnce(pctx)
	cancel()
	if err == nil || n != 0 {
		t.Fatalf("n=%d err=%v", n, err)
	}
	cp, err := rt.LoadCheckpoint(ctx, "nws-alerts")
	if err != nil || cp.ETag != "" || len(cp.Seen) != 0 {
		t.Fatalf("checkpoint written without acks: %+v %v", cp, err)
	}
	natstest.Streams(t, e.js, 2*time.Minute, "INTEL")
	if _, n, err := src.PollOnce(ctx); err != nil || n != 12 {
		t.Fatalf("after stream exists: n=%d err=%v", n, err)
	}
}

// Reliability: a failed checkpoint read is not "no state": the poll is skipped and the
// upstream is not even contacted (ingestion-platform.md finding 4).
func TestCheckpointReadErrorSkipsPoll(t *testing.T) {
	e := newEnv(t, 2*time.Minute)
	up := newUpstream(t, fixture(t, "nws/active-2026-10-29.json"))
	rt, _ := e.runtime(t, "a")
	src := nwsSource(t, rt, up.URL)
	e.srv.Stop()
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_, _, err := src.PollOnce(ctx)
	if !ingest.IsCheckpointReadError(err) {
		t.Fatalf("want checkpoint read error, got %v", err)
	}
	if up.requests.Load() != 0 {
		t.Fatal("upstream contacted without a checkpoint")
	}
}

// ---------------------------------------------------------------- schema and dead letters

// Reliability: bad records go to evt.sys.{source}.rejected with the reason; valid ones
// are published; a snapshot that keeps listing a bad record dead-letters it once.
func TestMalformedRecordsToDeadLetter(t *testing.T) {
	e := newEnv(t, 2*time.Minute)
	up := newUpstream(t, fixture(t, "nws/malformed.json"))
	up.etag = false
	rt, _ := e.runtime(t, "a")
	src := nwsSource(t, rt, up.URL)
	if _, n, err := src.PollOnce(context.Background()); err != nil || n != 1 {
		t.Fatalf("n=%d err=%v", n, err)
	}
	dlq := natstest.Msgs(t, e.js, "SYS", "evt.sys.nws-alerts.rejected")
	if len(dlq) != 5 {
		t.Fatalf("dead letters %d, want 5", len(dlq))
	}
	rejSch, _ := ingest.NewSchemas(schemaDir).Get("ingest/rejected.v1.schema.json")
	reasons := []string{}
	for _, m := range dlq {
		var env ingest.Envelope
		if err := json.Unmarshal(m.Data(), &env); err != nil {
			t.Fatal(err)
		}
		if err := ingest.Validate(rejSch, env.Data); err != nil {
			t.Fatal(err)
		}
		reasons = append(reasons, env.Data["reason"].(string))
	}
	all := strings.Join(reasons, " | ")
	for _, want := range []string{"no properties", "no properties.id", "severity", "lat", "sent"} {
		if !strings.Contains(all, want) {
			t.Errorf("no dead letter mentions %q: %s", want, all)
		}
	}
	if _, _, err := src.PollOnce(context.Background()); err != nil {
		t.Fatal(err)
	}
	if got := count(t, e.js, "SYS", "evt.sys.>"); got != 5 {
		t.Fatalf("re-listed bad records were dead-lettered again: %d", got)
	}
	if testutil.ToFloat64(rt.Metrics.Rejected.WithLabelValues("nws-alerts")) != 5 {
		t.Fatal("rejected metric")
	}
}

// A body that is not the expected document is one dead letter plus a parse error, and
// the checkpoint stays where it was.
func TestWholeBodyParseErrorDeadLettered(t *testing.T) {
	e := newEnv(t, 2*time.Minute)
	up := newUpstream(t, []byte("<html>upstream maintenance</html>"))
	rt, _ := e.runtime(t, "a")
	_, _, err := nwsSource(t, rt, up.URL).PollOnce(context.Background())
	if !ingest.IsParseError(err) || ingest.ErrorReason(err) != "parse" {
		t.Fatalf("want parse error, got %v", err)
	}
	if count(t, e.js, "SYS", "evt.sys.nws-alerts.rejected") != 1 {
		t.Fatal("no dead letter for unparseable body")
	}
	if cp, _ := rt.LoadCheckpoint(context.Background(), "nws-alerts"); cp.ETag != "" {
		t.Fatal("checkpoint moved")
	}
}

// ---------------------------------------------------------------- one owner per source

// Reliability: two instances for one source never both run it; after the leader dies
// (no clean release), the standby takes over within the lease TTL and nothing is lost
// or duplicated.
func TestLeaseSingleOwnerAndTakeover(t *testing.T) {
	e := newEnv(t, 2*time.Minute)
	up := newUpstream(t, fixture(t, "nws/active-2026-10-29.json"))
	rtA, _ := e.runtime(t, "a")
	rtB, _ := e.runtime(t, "b")
	var active, maxActive atomic.Int32
	var ranB atomic.Bool
	wrap := func(name string, s *ingest.Source) func(context.Context) {
		return func(ctx context.Context) {
			if v := active.Add(1); v > maxActive.Load() {
				maxActive.Store(v)
			}
			if name == "b" {
				ranB.Store(true)
			}
			s.Run(ctx)
			active.Add(-1)
		}
	}
	srcA, srcB := nwsSource(t, rtA, up.URL), nwsSource(t, rtB, up.URL)
	ctxA, cancelA := context.WithCancel(context.Background())
	ctxB, cancelB := context.WithCancel(context.Background())
	defer cancelB()
	go rtA.RunOwned(ctxA, "nws-alerts", wrap("a", srcA))
	waitFor(t, 3*time.Second, "a to lead", func() bool { return rtA.Board.Get("nws-alerts").Leader })
	go rtB.RunOwned(ctxB, "nws-alerts", wrap("b", srcB))
	waitFor(t, 5*time.Second, "a to publish", func() bool { return count(t, e.js, "INTEL", "intel.>") == 12 })
	time.Sleep(2 * e.cfg.LeaseTTL.Duration)
	if ranB.Load() || rtB.Board.Get("nws-alerts").Leader {
		t.Fatal("standby ran while the leader held the lease")
	}
	// Crash a: drop its connection first so it cannot release the lease, then stop it.
	rtA.NC.Close()
	cancelA()
	crashed := time.Now()
	waitFor(t, 3*e.cfg.LeaseTTL.Duration+2*time.Second, "b to take over", func() bool { return rtB.Board.Get("nws-alerts").Leader })
	took := time.Since(crashed)
	up.set(fixture(t, "nws/active-2026-10-29-next.json"), 200)
	waitFor(t, 5*time.Second, "b to publish the update", func() bool { return count(t, e.js, "INTEL", "intel.>") == 13 })
	if maxActive.Load() != 1 {
		t.Fatalf("%d instances ran the source at once", maxActive.Load())
	}
	ids := map[string]bool{}
	for _, m := range natstest.Msgs(t, e.js, "INTEL", "intel.>") {
		id := m.Headers().Get(jetstream.MsgIDHeader)
		if ids[id] {
			t.Fatalf("duplicate %s", id)
		}
		ids[id] = true
	}
	t.Logf("takeover after crash: %v (lease TTL %v)", took.Round(10*time.Millisecond), e.cfg.LeaseTTL.Duration)
}

// Reliability: a lease survives a bus outage longer than its TTL: the source keeps
// running (its context is never cancelled) and the lease is re-created afterwards.
func TestLeaseSurvivesBusOutage(t *testing.T) {
	e := newEnv(t, 2*time.Minute)
	rt, lb := e.runtime(t, "a")
	rt.Board.Add("src", ingest.KindRESTPoll, time.Minute)
	var cancelled atomic.Bool
	started := make(chan struct{})
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go rt.RunOwned(ctx, "src", func(c context.Context) {
		close(started)
		<-c.Done()
		cancelled.Store(true)
	})
	<-started
	e.srv.Stop()
	time.Sleep(3 * e.cfg.LeaseTTL.Duration)
	e.srv.Restart()
	waitFor(t, 10*time.Second, "lease back", func() bool {
		kv, err := e.srvJS(t).KeyValue(context.Background(), e.cfg.LockBucket)
		if err != nil {
			return false
		}
		v, err := kv.Get(context.Background(), "src")
		return err == nil && string(v.Value()) == "a"
	})
	if cancelled.Load() {
		t.Fatal("source stopped during a bus outage")
	}
	if !strings.Contains(lb.String(), "lease renew deferred") {
		t.Fatal("expected deferred renewals in the log")
	}
}

func (e *env) srvJS(t *testing.T) jetstream.JetStream {
	_, js := e.srv.Connect(t, "check")
	return js
}

// Reliability: when another instance holds the lease, the source stops.
func TestLeaseLostStopsSource(t *testing.T) {
	e := newEnv(t, 2*time.Minute)
	rt, _ := e.runtime(t, "a")
	rt.Board.Add("src", ingest.KindRESTPoll, time.Minute)
	stopped := make(chan struct{})
	started := make(chan struct{})
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go rt.RunOwned(ctx, "src", func(c context.Context) { close(started); <-c.Done(); close(stopped) })
	<-started
	if _, err := rt.Locks.Put(context.Background(), "src", []byte("b")); err != nil {
		t.Fatal(err)
	}
	select {
	case <-stopped:
	case <-time.After(3 * time.Second):
		t.Fatal("source kept running after losing the lease")
	}
}

// ---------------------------------------------------------------- freshness and observability

// Reliability: tel.feed.{source}.status carries fetchedAt, count and lastError; an
// upstream outage keeps the last good data, sets lastError and goes stale; metrics and
// structured logs carry the source.
func TestFreshnessStatusAndUpstreamOutage(t *testing.T) {
	e := newEnv(t, 2*time.Minute)
	up := newUpstream(t, fixture(t, "nws/active-2026-10-29.json"))
	rt, lb := e.runtime(t, "a")
	src := nwsSource(t, rt, up.URL)
	statusSch, _ := ingest.NewSchemas(schemaDir).Get("ingest/feed-status.v1.schema.json")
	var mu sync.Mutex
	var last map[string]any
	sub, err := rt.NC.Subscribe("tel.feed.nws-alerts.status", func(m *nats.Msg) {
		var v map[string]any
		_ = json.Unmarshal(m.Data, &v)
		mu.Lock()
		last = v
		mu.Unlock()
	})
	if err != nil {
		t.Fatal(err)
	}
	defer sub.Unsubscribe()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go rt.Heartbeat(ctx, 100*time.Millisecond)
	go rt.RunOwned(ctx, "nws-alerts", src.Run)
	get := func() map[string]any { mu.Lock(); defer mu.Unlock(); return last }
	waitFor(t, 5*time.Second, "fresh status", func() bool {
		s := get()
		return s != nil && s["fetchedAt"] != nil && s["count"] == float64(12) && s["lastError"] == nil && s["stale"] == false
	})
	if err := ingest.Validate(statusSch, get()); err != nil {
		t.Fatalf("status does not match schema: %v", err)
	}
	up.set(nil, 503)
	waitFor(t, 6*time.Second, "stale status with lastError", func() bool {
		s := get()
		return s["stale"] == true && s["lastError"] == "http 503" && s["count"] == float64(12)
	})
	if testutil.ToFloat64(rt.Metrics.FetchErrors.WithLabelValues("nws-alerts", "http_5xx")) < 1 {
		t.Fatal("fetch error metric")
	}
	if n, _ := testutil.GatherAndCount(rt.Metrics.Reg, "ingest_feed_stale"); n != 1 {
		t.Fatal("stale metric missing")
	}
	up.set(nil, 200)
	waitFor(t, 6*time.Second, "recovery", func() bool { s := get(); return s["stale"] == false && s["lastError"] == nil })
	logs := lb.String()
	if !strings.Contains(logs, `"source":"nws-alerts"`) || !strings.Contains(logs, `"level":"WARN"`) {
		t.Fatalf("structured logs missing source/level: %.300s", logs)
	}
	for _, line := range strings.Split(strings.TrimSpace(logs), "\n") {
		if !json.Valid([]byte(line)) {
			t.Fatalf("non-JSON log line: %s", line)
		}
	}
}

// ---------------------------------------------------------------- delivery-mode adapters

// SSE: resumes with Last-Event-ID across disconnects and a restart, with no gap and no
// duplicate.
func TestSSEResumeNoGapNoDuplicate(t *testing.T) {
	e := newEnv(t, 2*time.Minute)
	const total = 40
	var lastIDs []string
	var mu sync.Mutex
	stub := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		from := 0
		fmt.Sscan(r.Header.Get("Last-Event-ID"), &from)
		mu.Lock()
		lastIDs = append(lastIDs, r.Header.Get("Last-Event-ID"))
		mu.Unlock()
		w.Header().Set("Content-Type", "text/event-stream")
		fl := w.(http.Flusher)
		for i := from + 1; i <= min(from+7, total); i++ { // drop the connection every 7 events
			fmt.Fprintf(w, "id: %d\ndata: {\"id\":\"e%d\",\"lat\":1.5,\"lon\":2.5}\n\n", i, i)
			fl.Flush()
		}
	}))
	defer stub.Close()
	mk := func(rt *ingest.Runtime) *ingest.Source {
		c := sourceCfg("sse-test", ingest.KindSSE, "json", stub.URL, "intel.test.sse.event", "intel/generic-record.v1.schema.json")
		c.Options = map[string]string{"prefix": "sse", "idField": "id"}
		c.MaxBackoff = ingest.Duration{Duration: 100 * time.Millisecond}
		s, err := ingest.NewSource(rt, c)
		if err != nil {
			t.Fatal(err)
		}
		return s
	}
	rt1, _ := e.runtime(t, "a")
	ctx1, cancel1 := context.WithCancel(context.Background())
	go mk(rt1).Run(ctx1)
	waitFor(t, 10*time.Second, "first half", func() bool { return count(t, e.js, "INTEL", "intel.test.sse.event") >= 14 })
	cancel1() // restart: a new process resumes from the KV cursor
	time.Sleep(200 * time.Millisecond)
	rt2, _ := e.runtime(t, "a2")
	ctx2, cancel2 := context.WithCancel(context.Background())
	defer cancel2()
	go mk(rt2).Run(ctx2)
	waitFor(t, 15*time.Second, "all events", func() bool { return count(t, e.js, "INTEL", "intel.test.sse.event") >= total })
	time.Sleep(300 * time.Millisecond)
	ids := map[string]int{}
	for _, m := range natstest.Msgs(t, e.js, "INTEL", "intel.test.sse.event") {
		ids[m.Headers().Get(jetstream.MsgIDHeader)]++
	}
	for i := 1; i <= total; i++ {
		if ids[fmt.Sprintf("sse:e%d", i)] != 1 {
			t.Fatalf("event e%d stored %d times", i, ids[fmt.Sprintf("sse:e%d", i)])
		}
	}
	mu.Lock()
	defer mu.Unlock()
	if len(lastIDs) < 6 || lastIDs[0] != "" || lastIDs[1] != "7" {
		t.Fatalf("Last-Event-ID not sent on reconnect: %v", lastIDs)
	}
}

// WebSocket: frames are published, a bad frame is dead-lettered, and the subscribe
// message (env-expanded, so keys stay out of Git) is sent after connect.
func TestWebSocketAdapter(t *testing.T) {
	e := newEnv(t, 2*time.Minute)
	t.Setenv("TEST_WS_KEY", "k-1")
	got := make(chan string, 1)
	stub := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, err := websocket.Accept(w, r, nil)
		if err != nil {
			return
		}
		defer c.CloseNow()
		_, sub, _ := c.Read(r.Context())
		got <- string(sub)
		for i := 0; i < 50; i++ {
			_ = c.Write(r.Context(), websocket.MessageText, []byte(fmt.Sprintf(`{"mmsi":%d,"ts":1761732000,"lat":40.1,"lon":-73.9}`, 366000000+i)))
		}
		_ = c.Write(r.Context(), websocket.MessageText, []byte(`{"mmsi":1,"ts":1,"lat":999,"lon":0}`))
		_ = c.Write(r.Context(), websocket.MessageText, []byte(`not json`))
		time.Sleep(2 * time.Second)
	}))
	defer stub.Close()
	y := fmt.Sprintf(`
sources:
  - id: ais-test
    kind: websocket
    parser: json
    url: ws%s
    subject: evt.world.sea.placeholder
    schema: intel/generic-record.v1.schema.json
    subscribe: '{"key":"${TEST_WS_KEY}"}'
    queueMax: 1000
    options: {prefix: ais, idFields: "mmsi,ts", subjectTemplate: "evt.world.sea.{mmsi}"}
    terms: {licence: test, url: "https://example.org"}
`, strings.TrimPrefix(stub.URL, "http"))
	c, err := ingest.ParseConfig([]byte(y))
	if err != nil {
		t.Fatal(err)
	}
	rt, _ := e.runtime(t, "a")
	src, err := ingest.NewSource(rt, c.Sources[0])
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go src.Run(ctx)
	if s := <-got; s != `{"key":"k-1"}` {
		t.Fatalf("subscribe message %q", s)
	}
	waitFor(t, 10*time.Second, "50 positions + 2 dead letters", func() bool {
		return count(t, e.js, "WORLD", "evt.world.>") == 50 && count(t, e.js, "SYS", "evt.sys.ais-test.rejected") == 2
	})
	if count(t, e.js, "WORLD", "evt.world.sea.366000049") != 1 {
		t.Fatal("per-entity subject from subjectTemplate not used")
	}
}

// RSS: items become records on the configured subject.
func TestRSSAdapter(t *testing.T) {
	e := newEnv(t, 2*time.Minute)
	up := newUpstream(t, fixture(t, "rss/gdacs.xml"))
	rt, _ := e.runtime(t, "a")
	c := sourceCfg("gdacs-test", ingest.KindRSS, "rss", up.URL, "intel.disaster.gdacs.event", "intel/generic-record.v1.schema.json")
	c.Options = map[string]string{"prefix": "gdacs"}
	src, err := ingest.NewSource(rt, c)
	if err != nil {
		t.Fatal(err)
	}
	if _, n, err := src.PollOnce(context.Background()); err != nil || n != 2 {
		t.Fatalf("n=%d err=%v", n, err)
	}
	if count(t, e.js, "INTEL", "intel.disaster.gdacs.event") != 2 || count(t, e.js, "SYS", "evt.sys.gdacs-test.rejected") != 1 {
		t.Fatal("rss records or dead letter missing")
	}
}

// File/bulk: a file:// CSV in one-shot (CronJob) mode; an unchanged file is not re-read.
func TestFileBulkOnce(t *testing.T) {
	e := newEnv(t, 2*time.Minute)
	dir := t.TempDir()
	p := filepath.Join(dir, "firms.csv")
	if err := os.WriteFile(p, fixture(t, "file/firms-sample.csv"), 0o600); err != nil {
		t.Fatal(err)
	}
	rt, _ := e.runtime(t, "a")
	c := sourceCfg("firms-test", ingest.KindFile, "csv", "file://"+p, "intel.fire.firms.hotspot", "intel/generic-record.v1.schema.json")
	c.Options = map[string]string{"prefix": "firms", "idFields": "latitude,longitude,acq_date,acq_time", "latField": "latitude", "lonField": "longitude"}
	src, err := ingest.NewSource(rt, c)
	if err != nil {
		t.Fatal(err)
	}
	src.Once = true
	done := make(chan struct{})
	go func() { src.Run(context.Background()); close(done) }()
	select {
	case <-done:
	case <-time.After(5 * time.Second):
		t.Fatal("once mode did not exit")
	}
	if count(t, e.js, "INTEL", "intel.fire.firms.hotspot") != 3 || count(t, e.js, "SYS", "evt.sys.firms-test.rejected") != 1 {
		t.Fatal("bulk records or dead letter missing")
	}
	if recs, n, err := src.PollOnce(context.Background()); err != nil || recs != 0 || n != 0 {
		t.Fatalf("unchanged file re-read: %d %d %v", recs, n, err)
	}
}

// ---------------------------------------------------------------- replay

// Replayable: reading the stream from the start rebuilds exactly the parser output
// (same keys, same data), so the entity service can rebuild from scratch.
func TestReplayRebuildsRecords(t *testing.T) {
	e := newEnv(t, 2*time.Minute)
	body := fixture(t, "nws/active-2026-10-29.json")
	up := newUpstream(t, body)
	rt, _ := e.runtime(t, "a")
	if _, _, err := nwsSource(t, rt, up.URL).PollOnce(context.Background()); err != nil {
		t.Fatal(err)
	}
	want, _ := parsers.NWSAlerts(body, nil)
	msgs := natstest.Msgs(t, e.js, "INTEL", "intel.>")
	if len(msgs) != len(want) {
		t.Fatalf("%d vs %d", len(msgs), len(want))
	}
	for i, m := range msgs {
		var env ingest.Envelope
		_ = json.Unmarshal(m.Data(), &env)
		wb, _ := json.Marshal(want[i].Data)
		gb, _ := json.Marshal(env.Data)
		if env.MsgID != want[i].ID || string(wb) != string(gb) || env.ID != ingest.StableUUID(want[i].ID) {
			t.Fatalf("record %d differs on replay", i)
		}
	}
}
