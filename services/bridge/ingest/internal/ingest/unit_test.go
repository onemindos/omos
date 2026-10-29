package ingest

import (
	"bufio"
	"context"
	"errors"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus/testutil"
)

const validCfg = `
defaults:
  headers: {User-Agent: "omos-ingest-test"}
sources:
  - id: nws-alerts
    kind: rest_poll
    parser: nws-alerts
    url: https://api.weather.gov/alerts/active
    subject: intel.weather.noaa.alert
    schema: intel/weather-alert.v1.schema.json
    interval: 60s
    staleAfter: 10m
    terms: {licence: "open data", url: "https://www.weather.gov/documentation/services-web-api"}
`

func TestConfigDefaultsAndValidation(t *testing.T) {
	t.Setenv("TEST_KEY", "k-123")
	c, err := ParseConfig([]byte(validCfg + `    headers: {X-Key: "${TEST_KEY}"}` + "\n"))
	if err != nil {
		t.Fatal(err)
	}
	s := c.Sources[0]
	if c.LockBucket != "ingest-locks" || c.StateBucket != "ingest-state" || c.LeaseTTL.Duration != 10*time.Second {
		t.Fatalf("bucket defaults: %+v", c)
	}
	if s.Headers["User-Agent"] != "omos-ingest-test" || s.Headers["X-Key"] != "k-123" {
		t.Fatalf("headers not merged/expanded: %v", s.Headers)
	}
	if s.SeenMax != 4096 || s.MaxBodyBytes != 16<<20 || s.RateLimit.Burst != 1 || s.StaleAfter.Duration != 10*time.Minute {
		t.Fatalf("source defaults: %+v", s)
	}
	bad := map[string]string{
		"no terms":          strings.Replace(validCfg, `    terms: {licence: "open data", url: "https://www.weather.gov/documentation/services-web-api"}`, "", 1),
		"no schema":         strings.Replace(validCfg, "    schema: intel/weather-alert.v1.schema.json\n", "", 1),
		"off-grammar":       strings.Replace(validCfg, "intel.weather.noaa.alert", "fabric.intel.weather", 1),
		"5-token intel":     strings.Replace(validCfg, "intel.weather.noaa.alert", "intel.weather.noaa.alert.ny", 1),
		"ent forbidden":     strings.Replace(validCfg, "intel.weather.noaa.alert", "ent.alert.x.state", 1),
		"bad id":            strings.Replace(validCfg, "id: nws-alerts", "id: nws.alerts", 1),
		"impolite":          strings.Replace(validCfg, "interval: 60s", "interval: 100ms", 1),
		"unknown kind":      strings.Replace(validCfg, "kind: rest_poll", "kind: carrier_pigeon", 1),
		"quota without per": strings.Replace(validCfg, "    interval: 60s", "    interval: 60s\n    quota: {requests: 10}", 1),
	}
	for name, y := range bad {
		if _, err := ParseConfig([]byte(y)); err == nil {
			t.Errorf("%s: accepted", name)
		}
	}
	dup := validCfg + strings.SplitN(validCfg, "sources:\n", 2)[1]
	if _, err := ParseConfig([]byte(dup)); err == nil {
		t.Error("duplicate id accepted")
	}
}

func TestCheckDataSubjectGrammar(t *testing.T) {
	ok := []string{"intel.weather.noaa.alert", "evt.world.sea.366000184", "tel.weather.ws1.reading", "det.air.visual.track42", "tel.fleet.drone.sim-x1.position"}
	for _, s := range ok {
		if err := CheckDataSubject(s); err != nil {
			t.Errorf("%s: %v", s, err)
		}
	}
	bad := []string{"intel.weather.noaa", "intel.weather.noaa.alert.ny.flood", "evt.feed.x.rejected", "evt.world.sea", "cmd.drone.x.goto",
		"tak.cot.inject", "ent.ship.1.state", "svc.x.y", "intel.weather.*.alert", "intel..noaa.alert", "fabric.intel.x.y", "evt.world.sea.3.66e+08"}
	for _, s := range bad {
		if err := CheckDataSubject(s); err == nil {
			t.Errorf("%s accepted", s)
		}
	}
	if StatusSubject("nws-alerts") != "tel.feed.nws-alerts.status" || RejectSubject("nws-alerts") != "evt.sys.nws-alerts.rejected" {
		t.Fatal("fixed subjects changed")
	}
}

// Bounded memory: the stream queue never exceeds its capacity and drops the oldest.
func TestDropQueueBoundedDropsOldest(t *testing.T) {
	q := NewDropQueue[int](3)
	drops := 0
	for i := 1; i <= 5; i++ {
		if q.Push(i) {
			drops++
		}
	}
	if drops != 2 || q.Len() != 3 || q.Cap() != 3 {
		t.Fatalf("drops=%d len=%d", drops, q.Len())
	}
	ctx := context.Background()
	for _, want := range []int{3, 4, 5} {
		if v, _ := q.Pop(ctx); v != want {
			t.Fatalf("got %d want %d", v, want)
		}
	}
	cctx, cancel := context.WithTimeout(ctx, 20*time.Millisecond)
	defer cancel()
	if _, ok := q.Pop(cctx); ok {
		t.Fatal("pop on empty queue must wait for ctx")
	}
}

func TestKeepCurrentLastKeepsSnapshotInRing(t *testing.T) {
	seen := []string{"old1", "a", "old2", "b", "c"}
	recs := []Record{{ID: "a"}, {ID: "b"}, {ID: "x", Err: errors.New("bad")}, {ID: "c"}}
	got := keepCurrentLast(seen, recs)
	want := "old1,old2,a,b,c"
	if strings.Join(got, ",") != want {
		t.Fatalf("got %v want %s", got, want)
	}
}

func TestStableUUIDDeterministicV4(t *testing.T) {
	a, b := StableUUID("nws:urn:oid:1"), StableUUID("nws:urn:oid:1")
	re := regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-4[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$`)
	if a != b || !re.MatchString(a) || a == StableUUID("nws:urn:oid:2") {
		t.Fatalf("%s %s", a, b)
	}
}

func TestReadSSE(t *testing.T) {
	in := ": comment\nid: 1\nevent: e\ndata: {\"a\":\ndata: 1}\n\nretry: 1500\ndata: two\n\nid: 3\ndata: three\n\n"
	var got []SSEEvent
	err := ReadSSE(bufio.NewScanner(strings.NewReader(in)), func(e SSEEvent) error { got = append(got, e); return nil })
	if err != nil || len(got) != 3 {
		t.Fatalf("%v %v", got, err)
	}
	if got[0].ID != "1" || got[0].Event != "e" || got[0].Data != "{\"a\":\n1}" {
		t.Fatalf("event 1: %+v", got[0])
	}
	if got[1].ID != "1" || got[1].Retry != 1500*time.Millisecond { // last event id persists
		t.Fatalf("event 2: %+v", got[1])
	}
	if got[2].ID != "3" {
		t.Fatalf("event 3: %+v", got[2])
	}
}

// Freshness: fetchedAt, count, lastError and stale come from the status board.
func TestStatusBoardFreshness(t *testing.T) {
	now := time.Date(2026, 10, 29, 10, 0, 0, 0, time.UTC)
	b := NewStatusBoard("a", func() time.Time { return now })
	b.Add("nws-alerts", KindRESTPoll, time.Minute)
	if s := b.Get("nws-alerts"); !s.Stale || s.FetchedAt != nil {
		t.Fatal("never fetched must be stale")
	}
	b.OK("nws-alerts", 12, 3)
	now = now.Add(30 * time.Second)
	s := b.Get("nws-alerts")
	if s.Stale || s.Count != 12 || s.Published != 3 || s.LastError != nil || *s.AgeSec != 30 {
		t.Fatalf("%+v", s)
	}
	b.Fail("nws-alerts", errors.New("http 503"))
	now = now.Add(31 * time.Second)
	s = b.Get("nws-alerts")
	if !s.Stale || s.LastError == nil || *s.LastError != "http 503" || s.Count != 12 {
		t.Fatalf("upstream down keeps last good count and marks stale: %+v", s)
	}
}

// Observable: freshness series are exported at scrape time; a standby exports none.
func TestFreshnessMetrics(t *testing.T) {
	now := time.Date(2026, 10, 29, 10, 0, 0, 0, time.UTC)
	b := NewStatusBoard("a", func() time.Time { return now })
	b.Add("nws-alerts", KindRESTPoll, time.Minute)
	m := NewMetrics(b)
	if n, _ := testutil.GatherAndCount(m.Reg, "ingest_feed_stale"); n != 0 {
		t.Fatal("standby must not export freshness")
	}
	b.SetLeader("nws-alerts", true)
	b.OK("nws-alerts", 5, 5)
	now = now.Add(2 * time.Minute)
	want := `
# HELP ingest_feed_stale 1 when the source is older than its limit or never fetched.
# TYPE ingest_feed_stale gauge
ingest_feed_stale{source="nws-alerts"} 1
# HELP ingest_feed_age_seconds Seconds since the last successful fetch; -1 before the first.
# TYPE ingest_feed_age_seconds gauge
ingest_feed_age_seconds{source="nws-alerts"} 120
# HELP ingest_feed_records Records in the last successful fetch.
# TYPE ingest_feed_records gauge
ingest_feed_records{source="nws-alerts"} 5
`
	if err := testutil.GatherAndCompare(m.Reg, strings.NewReader(want), "ingest_feed_stale", "ingest_feed_age_seconds", "ingest_feed_records"); err != nil {
		t.Fatal(err)
	}
	if n, _ := testutil.GatherAndCount(m.Reg, "ingest_feed_last_success_timestamp_seconds", "ingest_feed_stale_after_seconds", "ingest_feed_up", "ingest_feed_leader"); n != 4 {
		t.Fatalf("missing freshness series: %d", n)
	}
}

func TestNextDelayBackoff(t *testing.T) {
	s := &Source{cfg: SourceConfig{SourceConfigBase: SourceConfigBase{Interval: Duration{time.Minute}, MaxBackoff: Duration{10 * time.Minute}}}}
	cases := map[int]time.Duration{0: time.Minute, 1: 2 * time.Minute, 3: 8 * time.Minute, 4: 10 * time.Minute, 10: 10 * time.Minute}
	for fails, want := range cases {
		if got := s.nextDelay(fails, errors.New("x")); got != want {
			t.Errorf("fails=%d got %v want %v", fails, got, want)
		}
	}
	reset := time.Now().Add(30 * time.Minute)
	if got := s.nextDelay(0, &QuotaError{ResetAt: reset}); got < 29*time.Minute {
		t.Errorf("quota wait %v", got)
	}
}

// No loss: in block mode (the websocket default) a full queue makes the reader wait for
// room instead of dropping a frame.
func TestDropQueueBlockModeLosesNothing(t *testing.T) {
	q := NewDropQueue[int](2)
	q.Block = true
	ctx := context.Background()
	go func() {
		for i := 1; i <= 100; i++ {
			if dropped, err := q.PushCtx(ctx, i); dropped || err != nil {
				panic("dropped in block mode")
			}
		}
	}()
	for want := 1; want <= 100; want++ {
		if v, _ := q.Pop(ctx); v != want {
			t.Fatalf("got %d want %d", v, want)
		}
	}
	_, _ = q.PushCtx(ctx, 1)
	_, _ = q.PushCtx(ctx, 2)
	cctx, cancel := context.WithTimeout(ctx, 20*time.Millisecond)
	defer cancel()
	if _, err := q.PushCtx(cctx, 3); err == nil {
		t.Fatal("push on a full blocking queue must wait for ctx")
	}
}

func TestQueueFullConfig(t *testing.T) {
	c, err := ParseConfig([]byte(strings.Replace(validCfg, "kind: rest_poll", "kind: rest_poll\n    queueFull: drop-oldest", 1)))
	if err != nil || c.Sources[0].QueueFull != QueueDropOldest {
		t.Fatalf("%v", err)
	}
	if c, _ := ParseConfig([]byte(validCfg)); c.Sources[0].QueueFull != QueueBlock {
		t.Fatal("default must be block")
	}
	if _, err := ParseConfig([]byte(strings.Replace(validCfg, "kind: rest_poll", "kind: rest_poll\n    queueFull: maybe", 1))); err == nil {
		t.Fatal("bad queueFull accepted")
	}
}

// The seen index lives with the checkpoint and stays in step with trimming.
func TestSeenIndexFollowsTrim(t *testing.T) {
	var cp Checkpoint
	for _, k := range []string{"a", "b", "c", "d"} {
		cp.addSeen(k)
	}
	idx := cp.seenIndex()
	cp.trimSeen(2)
	if _, ok := idx["a"]; ok || len(idx) != 2 || len(cp.Seen) != 2 {
		t.Fatalf("index %v seen %v", idx, cp.Seen)
	}
	cp.addSeen("e")
	if _, ok := idx["e"]; !ok {
		t.Fatal("index was rebuilt instead of reused")
	}
}
