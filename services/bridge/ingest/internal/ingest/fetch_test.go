package ingest

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"golang.org/x/time/rate"
)

func testFetcher(maxBody int64) (*Fetcher, *[]time.Duration) {
	var waits []time.Duration
	f := &Fetcher{Client: &http.Client{Timeout: 5 * time.Second}, Limiter: rate.NewLimiter(rate.Inf, 1),
		Headers: map[string]string{"User-Agent": "omos-ingest-test"}, MaxBody: maxBody, Attempts: 4,
		BaseBackoff: time.Second, MaxBackoff: time.Minute}
	f.sleep = func(_ context.Context, d time.Duration) error { waits = append(waits, d); return nil }
	return f, &waits
}

// Politeness: conditional headers are sent and a 304 is a no-change answer.
func TestFetchConditionalHeadersAnd304(t *testing.T) {
	var got http.Header
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		got = r.Header.Clone()
		if r.Header.Get("If-None-Match") == `"v1"` {
			w.WriteHeader(http.StatusNotModified)
			return
		}
		w.Header().Set("ETag", `"v1"`)
		w.Header().Set("Last-Modified", "Thu, 29 Oct 2026 10:00:00 GMT")
		_, _ = w.Write([]byte(`{}`))
	}))
	defer srv.Close()
	f, _ := testFetcher(1 << 20)
	r1, err := f.Get(context.Background(), srv.URL, Conditional{})
	if err != nil || r1.Code != 200 || r1.ETag != `"v1"` || r1.LastModified == "" {
		t.Fatalf("first fetch: %+v %v", r1, err)
	}
	if got.Get("User-Agent") != "omos-ingest-test" {
		t.Fatalf("User-Agent not sent: %v", got)
	}
	r2, err := f.Get(context.Background(), srv.URL, Conditional{ETag: r1.ETag, LastModified: r1.LastModified})
	if err != nil || !r2.NotModified() {
		t.Fatalf("second fetch should be 304: %+v %v", r2, err)
	}
	if got.Get("If-None-Match") != `"v1"` || got.Get("If-Modified-Since") != r1.LastModified {
		t.Fatalf("conditional headers missing: %v", got)
	}
	if r2.ETag != `"v1"` {
		t.Fatalf("304 must keep the validators: %+v", r2)
	}
}

// Politeness: 429 with Retry-After is honoured, then the request succeeds.
func TestFetch429HonoursRetryAfter(t *testing.T) {
	var n atomic.Int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if n.Add(1) == 1 {
			w.Header().Set("Retry-After", "7")
			w.WriteHeader(http.StatusTooManyRequests)
			return
		}
		_, _ = w.Write([]byte("ok"))
	}))
	defer srv.Close()
	f, waits := testFetcher(1 << 20)
	var codes []int
	f.OnResponse = func(c int) { codes = append(codes, c) }
	r, err := f.Get(context.Background(), srv.URL, Conditional{})
	if err != nil || r.Code != 200 {
		t.Fatalf("%+v %v", r, err)
	}
	if len(*waits) != 1 || (*waits)[0] != 7*time.Second {
		t.Fatalf("expected one 7s Retry-After wait, got %v", *waits)
	}
	if len(codes) != 2 || codes[0] != 429 || codes[1] != 200 {
		t.Fatalf("codes %v", codes)
	}
}

// Politeness: 5xx backs off exponentially (with jitter) and gives up after Attempts.
func TestFetch5xxBackoffThenGiveUp(t *testing.T) {
	var n atomic.Int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		n.Add(1)
		w.WriteHeader(http.StatusServiceUnavailable)
	}))
	defer srv.Close()
	f, waits := testFetcher(1 << 20)
	_, err := f.Get(context.Background(), srv.URL, Conditional{})
	var se *StatusError
	if !errors.As(err, &se) || se.Code != 503 {
		t.Fatalf("want http 503, got %v", err)
	}
	if n.Load() != 4 || len(*waits) != 3 {
		t.Fatalf("attempts %d waits %v", n.Load(), *waits)
	}
	// base 1s, 2s, 4s each jittered into [d/2, d]
	for i, w := range *waits {
		base := time.Second << i
		if w < base/2 || w > base {
			t.Fatalf("wait %d = %v outside [%v,%v]", i, w, base/2, base)
		}
	}
	if ErrorReason(err) != "http_5xx" {
		t.Fatal(ErrorReason(err))
	}
}

// Politeness: 403 (and other 4xx) is never retried.
func TestFetch403NoRetry(t *testing.T) {
	var n atomic.Int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		n.Add(1)
		w.WriteHeader(http.StatusForbidden)
	}))
	defer srv.Close()
	f, waits := testFetcher(1 << 20)
	_, err := f.Get(context.Background(), srv.URL, Conditional{})
	if n.Load() != 1 || len(*waits) != 0 || ErrorReason(err) != "http_4xx" {
		t.Fatalf("403 retried: n=%d waits=%v err=%v", n.Load(), *waits, err)
	}
}

// Transport errors are retried with backoff.
func TestFetchTransportErrorRetried(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	u := srv.URL
	srv.Close()
	f, waits := testFetcher(1 << 20)
	_, err := f.Get(context.Background(), u, Conditional{})
	if err == nil || len(*waits) != 3 || ErrorReason(err) != "transport" {
		t.Fatalf("err=%v waits=%v", err, *waits)
	}
}

// Bounded memory: a body over the limit is refused, not buffered.
func TestFetchBodyLimit(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(strings.Repeat("x", 2048)))
	}))
	defer srv.Close()
	f, _ := testFetcher(1024)
	if _, err := f.Get(context.Background(), srv.URL, Conditional{}); !errors.Is(err, ErrBodyTooLarge) {
		t.Fatalf("want ErrBodyTooLarge, got %v", err)
	}
}

// Politeness: the limiter spaces requests.
func TestFetchRateLimit(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	defer srv.Close()
	f, _ := testFetcher(1024)
	f.Limiter = rate.NewLimiter(rate.Every(100*time.Millisecond), 1)
	start := time.Now()
	for i := 0; i < 4; i++ {
		if _, err := f.Get(context.Background(), srv.URL, Conditional{}); err != nil {
			t.Fatal(err)
		}
	}
	if el := time.Since(start); el < 290*time.Millisecond {
		t.Fatalf("4 requests at 10/s took %v, limiter not applied", el)
	}
}

// Quotas: requests past the window budget fail with a QuotaError until the window resets.
func TestQuotaWindow(t *testing.T) {
	now := time.Date(2026, 10, 29, 0, 0, 0, 0, time.UTC)
	q := NewQuota(2, time.Hour, func() time.Time { return now })
	if q.Take() != nil || q.Take() != nil {
		t.Fatal("first two must pass")
	}
	err := q.Take()
	var qe *QuotaError
	if !errors.As(err, &qe) || !qe.ResetAt.Equal(now.Add(time.Hour)) || q.Remaining() != 0 {
		t.Fatalf("want quota error, got %v", err)
	}
	now = now.Add(time.Hour)
	if q.Take() != nil || q.Remaining() != 1 {
		t.Fatal("window did not reset")
	}
	if ErrorReason(err) != "quota" {
		t.Fatal(ErrorReason(err))
	}
	var nilQ *Quota
	if nilQ.Take() != nil || nilQ.Remaining() != -1 {
		t.Fatal("nil quota means unlimited")
	}
}

// The fetcher spends quota before contacting the upstream.
func TestFetchQuotaStopsRequests(t *testing.T) {
	var n atomic.Int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { n.Add(1) }))
	defer srv.Close()
	f, _ := testFetcher(1024)
	f.Quota = NewQuota(1, time.Hour, nil)
	_, _ = f.Get(context.Background(), srv.URL, Conditional{})
	_, err := f.Get(context.Background(), srv.URL, Conditional{})
	var qe *QuotaError
	if !errors.As(err, &qe) || n.Load() != 1 {
		t.Fatalf("n=%d err=%v", n.Load(), err)
	}
}

func TestRetryAfterParsing(t *testing.T) {
	now := time.Date(2026, 10, 29, 10, 0, 0, 0, time.UTC)
	if d, ok := retryAfter("120", now); !ok || d != 2*time.Minute {
		t.Fatal(d)
	}
	if d, ok := retryAfter("Thu, 29 Oct 2026 10:00:30 GMT", now); !ok || d != 30*time.Second {
		t.Fatal(d)
	}
	if _, ok := retryAfter("soon", now); ok {
		t.Fatal("garbage accepted")
	}
}

// Bulk files: file:// URLs use the file's mtime as Last-Modified, so an unchanged file is a 304.
func TestFetchFileURL(t *testing.T) {
	p := filepath.Join(t.TempDir(), "bulk.csv")
	if err := os.WriteFile(p, []byte("id\n1\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	f, _ := testFetcher(1024)
	r, err := f.Get(context.Background(), "file://"+p, Conditional{})
	if err != nil || r.Code != 200 || r.LastModified == "" {
		t.Fatalf("%+v %v", r, err)
	}
	r2, err := f.Get(context.Background(), "file://"+p, Conditional{LastModified: r.LastModified})
	if err != nil || !r2.NotModified() {
		t.Fatalf("%+v %v", r2, err)
	}
}

// Politeness: a Retry-After longer than the backoff cap is not cut short: the fetch gives
// up now and the poll loop waits the full Retry-After.
func TestFetchLongRetryAfterHonouredByPollLoop(t *testing.T) {
	var n atomic.Int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		n.Add(1)
		w.Header().Set("Retry-After", "3600")
		w.WriteHeader(http.StatusTooManyRequests)
	}))
	defer srv.Close()
	f, waits := testFetcher(1024)
	_, err := f.Get(context.Background(), srv.URL, Conditional{})
	var se *StatusError
	if !errors.As(err, &se) || se.RetryAfter != time.Hour || n.Load() != 1 || len(*waits) != 0 {
		t.Fatalf("err=%v n=%d waits=%v", err, n.Load(), *waits)
	}
	s := &Source{cfg: SourceConfig{SourceConfigBase: SourceConfigBase{Interval: Duration{time.Minute}, MaxBackoff: Duration{10 * time.Minute}}}}
	if d := s.nextDelay(1, err); d != time.Hour {
		t.Fatalf("next poll in %v, want 1h", d)
	}
}
