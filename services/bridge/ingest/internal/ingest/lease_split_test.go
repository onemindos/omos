package ingest_test

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/nats-io/nats-server/v2/server"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"

	"github.com/onemindos/omos-bridge/ingest/internal/ingest"
	"github.com/onemindos/omos-bridge/ingest/internal/natstest"
)

func (e *env) runtimeOn(t *testing.T, instance string, nc *nats.Conn, js jetstream.JetStream) *ingest.Runtime {
	t.Helper()
	board := ingest.NewStatusBoard(instance, nil)
	rt := &ingest.Runtime{NC: nc, JS: js, Board: board, Metrics: ingest.NewMetrics(board), Schemas: ingest.NewSchemas(schemaDir),
		Log: slog.New(slog.NewJSONHandler(&logBuf{}, nil)).With("instance", instance), Instance: instance, LeaseTTL: e.cfg.LeaseTTL.Duration,
		AckTimeout: time.Second, MaxPublishBackoff: 200 * time.Millisecond}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := rt.BindBuckets(ctx, e.cfg, true); err != nil {
		t.Fatal(err)
	}
	return rt
}

type tick struct {
	I string `json:"i"`
	N int    `json:"n"`
}

// ticker publishes a numbered record every 10 ms until ctx ends.
func ticker(rt *ingest.Runtime, inst string, started, stopped *atomic.Int64) func(context.Context) {
	n := 0
	return func(ctx context.Context) {
		started.Store(time.Now().UnixNano())
		defer func() { stopped.Store(time.Now().UnixNano()) }()
		for ctx.Err() == nil {
			n++
			b, _ := json.Marshal(tick{inst, n})
			if _, err := rt.Publish(ctx, "src", "intel.test.lease.tick", fmt.Sprintf("%s-%d", inst, n), b); err != nil {
				return
			}
			select {
			case <-ctx.Done():
			case <-time.After(10 * time.Millisecond):
			}
		}
	}
}

// Reliability (A0 "two collector instances never both publish"): the leader alone is cut
// off from the bus (a network split, a node partition) for longer than the lease TTL while
// the standby stays connected. The old leader must stop before the standby takes over,
// must not replay buffered publishes on reconnect, and must not overwrite the new owner's
// checkpoint afterwards.
func TestLeaderCutOffAloneNeverPublishesAfterTakeover(t *testing.T) {
	e := newEnv(t, 2*time.Minute)
	proxy := e.srv.Proxy(t)
	ncA, jsA := natstest.ConnectURL(t, proxy.URL(), "omos-ingest-a")
	rtA := e.runtimeOn(t, "a", ncA, jsA)
	ncB, jsB := e.srv.Connect(t, "omos-ingest-b")
	rtB := e.runtimeOn(t, "b", ncB, jsB)
	rtA.Board.Add("src", ingest.KindRESTPoll, time.Minute)
	rtB.Board.Add("src", ingest.KindRESTPoll, time.Minute)

	var aStart, aStop, bStart, bStop atomic.Int64
	ctxA, cancelA := context.WithCancel(context.Background())
	defer cancelA()
	ctxB, cancelB := context.WithCancel(context.Background())
	defer cancelB()
	go rtA.RunOwned(ctxA, "src", ticker(rtA, "a", &aStart, &aStop))
	waitFor(t, 3*time.Second, "a to publish", func() bool { return count(t, e.js, "INTEL", "intel.>") >= 20 })
	go rtB.RunOwned(ctxB, "src", ticker(rtB, "b", &bStart, &bStop))
	time.Sleep(e.cfg.LeaseTTL.Duration)
	if bStart.Load() != 0 {
		t.Fatal("standby ran while the leader held the lease")
	}

	proxy.Cut()
	cut := time.Now()
	waitFor(t, 3*e.cfg.LeaseTTL.Duration+2*time.Second, "b to take over", func() bool { return bStart.Load() != 0 })
	if aStop.Load() == 0 || aStop.Load() >= bStart.Load() {
		t.Fatalf("old leader still running when the standby started (a stopped %d, b started %d)", aStop.Load(), bStart.Load())
	}
	t.Logf("a stopped %v after the cut; b started %v after the cut",
		time.Duration(aStop.Load()-cut.UnixNano()).Round(time.Millisecond), time.Duration(bStart.Load()-cut.UnixNano()).Round(time.Millisecond))

	// Keep a cut off well past the TTL, then let it back in and give it time to misbehave.
	time.Sleep(3*e.cfg.LeaseTTL.Duration - time.Since(cut))
	proxy.Restore()
	waitFor(t, 5*time.Second, "a to reconnect", func() bool { return ncA.IsConnected() })
	time.Sleep(3 * e.cfg.LeaseTTL.Duration)
	if rtA.Board.Get("src").Leader || !rtB.Board.Get("src").Leader {
		t.Fatal("lease moved back to the old leader")
	}

	// Checkpoint fencing: b owns the checkpoint; a (lease lost) cannot overwrite it, even
	// with the revision it read before the split.
	ctx := context.Background()
	stale, err := rtA.LoadCheckpoint(ctx, "src")
	if err != nil {
		t.Fatal(err)
	}
	cpB, _ := rtB.LoadCheckpoint(ctx, "src")
	cpB.Cursor = "b-newer"
	if err := rtB.SaveCheckpoint(ctx, "src", &cpB, 0); err != nil {
		t.Fatalf("owner checkpoint: %v", err)
	}
	stale.Cursor = "a-older"
	if err := rtA.SaveCheckpoint(ctx, "src", &stale, 0); !errors.Is(err, ingest.ErrCheckpointFenced) {
		t.Fatalf("old leader wrote the checkpoint: %v", err)
	}
	if got, _ := rtB.LoadCheckpoint(ctx, "src"); got.Cursor != "b-newer" {
		t.Fatalf("checkpoint overwritten: %q", got.Cursor)
	}

	cancelB()
	msgs := natstest.Msgs(t, e.js, "INTEL", "intel.>")
	firstB, lastA, nA, nB := uint64(0), uint64(0), 0, 0
	for _, m := range msgs {
		var tk tick
		_ = json.Unmarshal(m.Data(), &tk)
		md, _ := m.Metadata()
		switch tk.I {
		case "a":
			nA++
			lastA = md.Sequence.Stream
		case "b":
			nB++
			if firstB == 0 {
				firstB = md.Sequence.Stream
			}
		}
	}
	t.Logf("a stored %d, b stored %d; last a seq %d, first b seq %d", nA, nB, lastA, firstB)
	if nB == 0 {
		t.Fatal("standby never published")
	}
	if lastA > firstB {
		t.Fatalf("old leader's message stored at seq %d after the new leader's first at %d", lastA, firstB)
	}
}

// Fencing: checkpoint writes are compare-and-set on the revision read, so two writers
// cannot silently overwrite each other.
func TestCheckpointCompareAndSet(t *testing.T) {
	e := newEnv(t, 2*time.Minute)
	rtA, _ := e.runtime(t, "a")
	rtB, _ := e.runtime(t, "b")
	ctx := context.Background()
	a, _ := rtA.LoadCheckpoint(ctx, "src")
	b, _ := rtB.LoadCheckpoint(ctx, "src")
	b.Cursor = "b1"
	if err := rtB.SaveCheckpoint(ctx, "src", &b, 0); err != nil {
		t.Fatal(err)
	}
	a.Cursor = "a1"
	if err := rtA.SaveCheckpoint(ctx, "src", &a, 0); !errors.Is(err, ingest.ErrCheckpointFenced) {
		t.Fatalf("stale create not fenced: %v", err)
	}
	b.Cursor = "b2"
	if err := rtB.SaveCheckpoint(ctx, "src", &b, 0); err != nil {
		t.Fatalf("writer on its own revision: %v", err)
	}
	if got, _ := rtA.LoadCheckpoint(ctx, "src"); got.Cursor != "b2" {
		t.Fatalf("got %q", got.Cursor)
	}
}

// Errors no retry can fix do not block a feed: a message over max_payload fails at once,
// and records over the stream's max_msg_size are dead-lettered with the reason.
func TestPermanentPublishErrorsDeadLettered(t *testing.T) {
	e := newEnv(t, 2*time.Minute)
	rt, _ := e.runtime(t, "a")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	t0 := time.Now()
	_, err := rt.Publish(ctx, "src", "intel.test.big.record", "big-1", make([]byte, 2<<20))
	if !errors.Is(err, ingest.ErrPermanentPublish) || time.Since(t0) > 2*time.Second {
		t.Fatalf("max payload: err %v after %v", err, time.Since(t0))
	}

	st, err := e.js.Stream(ctx, "INTEL")
	if err != nil {
		t.Fatal(err)
	}
	cfg := st.CachedInfo().Config
	cfg.MaxMsgSize = 512
	if _, err := e.js.UpdateStream(ctx, cfg); err != nil {
		t.Fatal(err)
	}
	up := newUpstream(t, fixture(t, "nws/active-2026-10-29.json"))
	src := nwsSource(t, rt, up.URL)
	t0 = time.Now()
	recs, n, err := src.PollOnce(ctx)
	if err != nil {
		t.Fatalf("poll blocked or failed: %v", err)
	}
	if n != 0 || count(t, e.js, "INTEL", "intel.>") != 0 {
		t.Fatalf("stored %d oversize records", n)
	}
	if got := count(t, e.js, "SYS", "evt.sys.nws-alerts.rejected"); got != recs {
		t.Fatalf("dead letters %d, want %d", got, recs)
	}
	m := natstest.Msgs(t, e.js, "SYS", "evt.sys.nws-alerts.rejected")[0]
	if !strings.Contains(string(m.Data()), "publish cannot succeed") {
		t.Fatalf("reason missing: %s", m.Data())
	}
	t.Logf("%d oversize records dead-lettered in %v", recs, time.Since(t0).Round(time.Millisecond))
}

// A publish permission violation (per-service users) is also permanent: the record is
// dead-lettered instead of retried forever.
func TestPermissionViolationDeadLettered(t *testing.T) {
	s := natstest.Start(t, func(o *server.Options) {
		o.Users = []*server.User{
			{Username: "admin", Password: "test-only"},
			{Username: "collector", Password: "test-only", Permissions: &server.Permissions{
				Publish: &server.SubjectPermission{Allow: []string{"$JS.API.>", "$KV.>", "_INBOX.>", "evt.sys.>"}},
			}},
		}
	})
	_, adminJS := natstest.ConnectURL(t, s.URL(), "admin", nats.UserInfo("admin", "test-only"))
	natstest.Streams(t, adminJS, 2*time.Minute)
	nc, js := natstest.ConnectURL(t, s.URL(), "collector", nats.UserInfo("collector", "test-only"))
	cfg, _ := ingest.ParseConfig([]byte("leaseTTL: 1500ms\n"))
	e := &env{srv: s, js: adminJS, cfg: cfg}
	rt := e.runtimeOn(t, "a", nc, js)
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	src := nwsSource(t, rt, newUpstream(t, fixture(t, "nws/active-2026-10-29.json")).URL)
	if _, _, err := src.PollOnce(ctx); err != nil {
		t.Fatalf("poll blocked or failed: %v", err)
	}
	if count(t, adminJS, "INTEL", "intel.>") != 0 {
		t.Fatal("denied subject stored")
	}
	msgs := natstest.Msgs(t, adminJS, "SYS", "evt.sys.nws-alerts.rejected")
	if len(msgs) == 0 || !strings.Contains(strings.ToLower(string(msgs[0].Data())), "permission") {
		t.Fatalf("no permission dead letter (%d)", len(msgs))
	}
}
