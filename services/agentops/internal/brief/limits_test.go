package brief

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/onemindos/omos/services/agentops/internal/bus"
	"github.com/onemindos/omos/services/agentops/internal/model"
	"github.com/onemindos/omos/services/agentops/internal/natstest"
)

// gate is a model that answers only after release is closed, and counts its calls.
type gate struct {
	release chan struct{}
	calls   atomic.Int64
}

func (g *gate) Name() string { return "gate" }
func (g *gate) Local() bool  { return true }
func (g *gate) Complete(ctx context.Context, _ model.Request) (string, error) {
	g.calls.Add(1)
	select {
	case <-g.release:
		return "Quiet AOI.", nil
	case <-ctx.Done():
		return "", ctx.Err()
	}
}

func errCode(raw []byte) string {
	var r struct {
		Error *struct {
			Code string `json:"code"`
		} `json:"error"`
	}
	_ = json.Unmarshal(raw, &r)
	if r.Error == nil {
		return ""
	}
	return r.Error.Code
}

// TestAdhocGenerateIsBounded: a flood of ad-hoc svc.brief.generate requests is answered busy
// beyond MaxPending, ad-hoc briefs stop at AdhocPerHour, model calls stop at ModelPerHour
// (the brief is then written from the facts), ad-hoc ids never enter the latest map, and a
// request whose context ends while it waits returns instead of running late.
func TestAdhocGenerateIsBounded(t *testing.T) {
	srv := natstest.Start(t)
	nc := srv.Connect(t)
	natstest.Streams(t, nc, true)
	stubEntities(t, srv.Connect(t))
	pub, _ := bus.NewPublisher(nc, bus.ModeAuto)
	g := &gate{release: make(chan struct{})}
	chain := model.Chain{Models: []model.Model{g}, Timeout: 30 * time.Second}
	b := New(Config{AOIs: []AOI{home}, MaxPending: 3, AdhocPerHour: 5, ModelPerHour: 2}, nc, pub, chain, nil, nil,
		slog.New(slog.NewTextHandler(io.Discard, nil)))
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() { _ = b.Serve(ctx) }()
	time.Sleep(100 * time.Millisecond)

	adhoc := func(i int) []byte {
		return []byte(fmt.Sprintf(`{"aoi":{"id":"flood-%d","bbox":[-97.80,30.20,-97.70,30.30]}}`, i))
	}
	// 1. Flood of 20 while the model is held: 3 accepted (1 running, 2 waiting), 17 busy at once.
	var wg sync.WaitGroup
	codes := make([]string, 20)
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			m, err := nc.Request(GenerateSubject, adhoc(i), 20*time.Second)
			if err != nil {
				codes[i] = "transport:" + err.Error()
				return
			}
			codes[i] = errCode(m.Data)
		}(i)
		time.Sleep(5 * time.Millisecond) // let the first requests take their slots
	}
	natstest.Eventually(t, 5*time.Second, "busy replies", func() bool { return b.C.Busy.Load() == 17 })
	close(g.release)
	wg.Wait()
	busy, ok := 0, 0
	for _, c := range codes {
		switch c {
		case "busy":
			busy++
		case "":
			ok++
		default:
			t.Fatalf("unexpected reply %q", c)
		}
	}
	if busy != 17 || ok != 3 {
		t.Fatalf("busy=%d ok=%d", busy, ok)
	}
	natstest.Evidence(t, "flood of 20 ad-hoc requests with MaxPending=3: %d answered, %d busy", ok, busy)

	// 2. Ad-hoc hourly cap: 3 used; 2 more pass, then rate_limited.
	for i := 20; i < 23; i++ {
		m, err := nc.Request(GenerateSubject, adhoc(i), 20*time.Second)
		if err != nil {
			t.Fatal(err)
		}
		want := ""
		if i == 22 {
			want = "rate_limited"
		}
		if c := errCode(m.Data); c != want {
			t.Fatalf("request %d: code %q, want %q", i, c, want)
		}
	}
	// 3. Model cap 2: five briefs written, two model calls; the rest say so.
	m, err := nc.Request(GenerateSubject, []byte(`{"aoi_id":"home"}`), 20*time.Second)
	if err != nil || errCode(m.Data) != "" {
		t.Fatalf("configured AOI refused after the ad-hoc cap: %v %s", err, m.Data)
	}
	var br Brief
	_ = json.Unmarshal(m.Data, &br)
	if g.calls.Load() != 2 || br.Model.Skipped != "model_hourly_cap" || br.Model.Used != "none" || len(br.Citations) < 3 {
		t.Fatalf("model calls %d, model %+v, cited %d", g.calls.Load(), br.Model, len(br.Citations))
	}
	// 4. Only the configured AOI is kept.
	b.mu.Lock()
	n := len(b.latest)
	b.mu.Unlock()
	if n != 1 || b.Latest("home") == nil || b.Latest("flood-0") != nil {
		t.Fatalf("latest holds %d entries", n)
	}
	natstest.Evidence(t, "after %d briefs: latest map size %d, model calls %d (cap 2), ad-hoc limited %d, model limited %d",
		b.C.Briefs.Load(), n, g.calls.Load(), b.C.AdhocLimited.Load(), b.C.ModelLimited.Load())

	// 5. Bad inline AOIs are refused before any work.
	for _, body := range []string{
		`{"aoi":{"id":"../../x","bbox":[-97.80,30.20,-97.70,30.30]}}`,
		`{"aoi":{"id":"x","bbox":[1,1,0,0]}}`,
		`{"aoi":{"id":"x","bbox":[-97.80,30.20,-97.70,30.30],"name":"` + strings.Repeat("a", 5000) + `"}}`,
	} {
		m, err := nc.Request(GenerateSubject, []byte(body), 5*time.Second)
		if err != nil || errCode(m.Data) != "bad_request" {
			t.Fatalf("accepted %.60s: %v %s", body, err, m.Data)
		}
	}

	// 6. A wait whose context ends returns at once with the context error.
	b.gen <- struct{}{}
	wctx, wcancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	start := time.Now()
	_, err = b.Generate(wctx, home)
	wcancel()
	<-b.gen
	if !errors.Is(err, context.DeadlineExceeded) || time.Since(start) > 2*time.Second {
		t.Fatalf("waiting generate: %v after %s", err, time.Since(start))
	}
}
