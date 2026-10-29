package bus

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"sync/atomic"
	"testing"
	"time"

	"github.com/nats-io/nats.go/jetstream"

	"github.com/onemindos/omos/services/entity/internal/engine"
	"github.com/onemindos/omos/services/entity/internal/model"
)

func startRun(t *testing.T, svc *Service) context.CancelFunc {
	t.Helper()
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
	t.Cleanup(func() { cancel(); <-errc })
	return cancel
}

// svc.entity.query pages a result far larger than one reply: every page fits max_payload,
// the cursor walks all 1,200 aircraft exactly once, and an over-large limit is halved.
func TestQueryPagesLargeResult(t *testing.T) {
	s := startServer(t)
	nc, js := connect(t, s)
	setupStreams(t, js)
	const n = 1200
	for i := 0; i < n; i++ {
		uid := fmt.Sprintf("icao-%06x", i)
		body := fmt.Sprintf(`{"uid":%q,"lat":40.%04d,"lon":-73.7,"stale":600,"t":%d}`, uid, i, t0+1000)
		if _, err := js.Publish(context.Background(), "evt.world.air."+uid, []byte(body)); err != nil {
			t.Fatal(err)
		}
	}
	cfg := testConfig("query")
	cfg.ServeRPC, cfg.QueryLimit = true, 5000
	svc, _ := newService(t, nc, cfg)
	startRun(t, svc)
	drained(t, js, "entity-svc", 30*time.Second)

	seen := map[string]bool{}
	cursor, pages, maxBody := "", 0, 0
	for {
		req, _ := json.Marshal(map[string]any{"types": []string{"aircraft"}, "limit": 5000, "cursor": cursor})
		m, err := nc.Request("svc.entity.query", req, 5*time.Second)
		if err != nil {
			t.Fatalf("page %d: %v", pages, err)
		}
		if code := m.Header.Get("Nats-Service-Error-Code"); code != "" {
			t.Fatalf("page %d: error %s %s", pages, code, m.Header.Get("Nats-Service-Error"))
		}
		var q struct {
			Entities []model.Entity `json:"entities"`
			Cursor   *string        `json:"cursor"`
		}
		if err := json.Unmarshal(m.Data, &q); err != nil {
			t.Fatal(err)
		}
		if len(m.Data) > maxBody {
			maxBody = len(m.Data)
		}
		for _, e := range q.Entities {
			if seen[e.EntityID] {
				t.Fatalf("duplicate %s across pages", e.EntityID)
			}
			seen[e.EntityID] = true
		}
		pages++
		if q.Cursor == nil {
			break
		}
		cursor = *q.Cursor
		if pages > 50 {
			t.Fatal("paging does not end")
		}
	}
	if len(seen) != n || pages < 2 || int64(maxBody) > nc.MaxPayload() {
		t.Fatalf("seen %d of %d in %d pages, largest reply %d bytes (max_payload %d)", len(seen), n, pages, maxBody, nc.MaxPayload())
	}
	t.Logf("%d entities in %d pages, largest reply %d bytes, max_payload %d", len(seen), pages, maxBody, nc.MaxPayload())
}

const hostileEUD = `<event version="2.0" uid="ANDROID-hostile1" type="a-h-G-U-C" how="m-g" time="2026-10-29T12:00:03.000Z" ` +
	`start="2026-10-29T12:00:03.000Z" stale="2026-10-29T12:06:03.000Z"><point lat="30.2" lon="-97.7" hae="200" ce="9" le="9"/>` +
	`<detail><takv device="Pixel 8" platform="ATAK-CIV" os="34" version="5.4.0"/><contact callsign="H-1"/></detail></event>`

// A stop-line SL9 reject is stored in SYS with an acked publish. While SYS is unavailable
// the input is not acked (redelivered), and the redelivery is deduplicated once SYS exists.
func TestDeadLetterIsAckedOrInputRedelivered(t *testing.T) {
	s := startServer(t)
	nc, js := connect(t, s)
	setupStreams(t, js)
	ctx := context.Background()
	if err := js.DeleteStream(ctx, "SYS"); err != nil {
		t.Fatal(err)
	}
	if _, err := js.Publish(ctx, "tak.cot.out", []byte(hostileEUD)); err != nil {
		t.Fatal(err)
	}
	svc, eng := newService(t, nc, testConfig("dlq"))
	startRun(t, svc)
	deadline := time.Now().Add(10 * time.Second)
	for eng.Counters()["dead_letter_failed"] == 0 {
		if time.Now().After(deadline) {
			t.Fatalf("no failed dead letter counted: %v", eng.Counters())
		}
		time.Sleep(50 * time.Millisecond)
	}
	c, _ := js.Consumer(ctx, "COT", "entity-svc-cot")
	info, _ := c.Info(ctx)
	if info.AckFloor.Stream != 0 {
		t.Fatalf("input acked although its dead letter was lost: ack floor %d", info.AckFloor.Stream)
	}
	if _, err := js.CreateStream(ctx, jetstream.StreamConfig{Name: "SYS", Subjects: []string{"evt.sys.>"},
		Storage: jetstream.FileStorage, Duplicates: 2 * time.Minute}); err != nil {
		t.Fatal(err)
	}
	drained(t, js, "entity-svc", 20*time.Second)
	st, _ := js.Stream(ctx, "SYS")
	si, _ := st.Info(ctx)
	if si.State.Msgs != 1 {
		t.Fatalf("SYS holds %d dead letters, want 1", si.State.Msgs)
	}
	m, err := st.GetLastMsgForSubject(ctx, "evt.sys.entity.rejected")
	if err != nil {
		t.Fatal(err)
	}
	if m.Header.Get("Nats-Msg-Id") != "COT:1" {
		t.Fatalf("dead letter msg id %q", m.Header.Get("Nats-Msg-Id"))
	}
	t.Logf("dead letter stored once after SYS returned: %s", m.Data)
}

// With MsgTTL off, the final live:false record is deleted from ENT after FinalRetain.
func TestFinalRecordDeletedWhenTTLOff(t *testing.T) {
	s := startServer(t)
	nc, js := connect(t, s)
	setupStreams(t, js)
	ctx := context.Background()
	if _, err := js.Publish(ctx, "evt.world.air.icao-0000aa",
		[]byte(fmt.Sprintf(`{"uid":"icao-0000aa","lat":40.6,"lon":-73.7,"stale":30,"t":%d}`, t0))); err != nil {
		t.Fatal(err)
	}
	var now atomic.Int64
	now.Store(t0 + 5000)
	eng := engine.New(engine.DefaultConfig(), now.Load)
	cfg := testConfig("final")
	cfg.MsgTTL, cfg.SweepEvery, cfg.FinalRetain = false, 50*time.Millisecond, 300*time.Millisecond
	svc, err := New(nc, eng, cfg, slog.New(slog.NewTextHandler(io.Discard, nil)))
	if err != nil {
		t.Fatal(err)
	}
	startRun(t, svc)
	drained(t, js, "entity-svc", 20*time.Second)
	st, _ := js.Stream(ctx, "ENT")
	subj := "ent.aircraft.icao-0000aa.state"
	if _, err := st.GetLastMsgForSubject(ctx, subj); err != nil {
		t.Fatalf("live record missing: %v", err)
	}
	now.Store(t0 + 200000) // past expiry
	sawFinal := false
	deadline := time.Now().Add(10 * time.Second)
	for {
		m, err := st.GetLastMsgForSubject(ctx, subj)
		if errors.Is(err, jetstream.ErrMsgNotFound) {
			break
		}
		if err == nil {
			var e model.Entity
			_ = json.Unmarshal(m.Data, &e)
			sawFinal = sawFinal || !e.Live
		}
		if time.Now().After(deadline) {
			t.Fatalf("final record not deleted (saw final %v): %v", sawFinal, eng.Counters())
		}
		time.Sleep(20 * time.Millisecond)
	}
	if eng.Counters()["final_deleted"] != 1 || eng.Counters()["entities"] != 0 {
		t.Fatalf("counters %v (saw final %v)", eng.Counters(), sawFinal)
	}
}
