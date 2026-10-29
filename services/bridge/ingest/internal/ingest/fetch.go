package ingest

import (
	"context"
	"errors"
	"fmt"
	"io"
	"math/rand/v2"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

// ErrBodyTooLarge keeps memory bounded: a response over MaxBodyBytes is refused, not buffered.
var ErrBodyTooLarge = errors.New("response body exceeds maxBodyBytes")

// StatusError is a non-retryable upstream answer (4xx other than 429) or a retryable one
// that kept failing after every attempt.
type StatusError struct{ Code int }

func (e *StatusError) Error() string { return fmt.Sprintf("http %d", e.Code) }

// QuotaError means the source used its request quota for the current window.
type QuotaError struct{ ResetAt time.Time }

func (e *QuotaError) Error() string {
	return "quota exhausted until " + e.ResetAt.UTC().Format(time.RFC3339)
}

// Quota is a fixed-window request budget (for example 1000 requests per day). It is kept
// in memory: a restart starts a fresh window, so set it below the upstream's real limit.
type Quota struct {
	mu    sync.Mutex
	max   int
	per   time.Duration
	used  int
	start time.Time
	now   func() time.Time
}

func NewQuota(maxReq int, per time.Duration, now func() time.Time) *Quota {
	if maxReq <= 0 {
		return nil
	}
	if now == nil {
		now = time.Now
	}
	return &Quota{max: maxReq, per: per, now: now, start: now()}
}

func (q *Quota) Take() error {
	if q == nil {
		return nil
	}
	q.mu.Lock()
	defer q.mu.Unlock()
	now := q.now()
	if now.Sub(q.start) >= q.per {
		q.start, q.used = now, 0
	}
	if q.used >= q.max {
		return &QuotaError{ResetAt: q.start.Add(q.per)}
	}
	q.used++
	return nil
}

func (q *Quota) Remaining() int {
	if q == nil {
		return -1
	}
	q.mu.Lock()
	defer q.mu.Unlock()
	if q.now().Sub(q.start) >= q.per {
		return q.max
	}
	return q.max - q.used
}

// Conditional carries the validators from the last acked fetch.
type Conditional struct {
	ETag         string
	LastModified string
}

type Response struct {
	Code         int
	Body         []byte
	ETag         string
	LastModified string
}

func (r Response) NotModified() bool { return r.Code == http.StatusNotModified }

// Fetcher does polite upstream requests for one source: a rate limiter, an optional
// quota, conditional headers, Retry-After on 429/503, capped exponential backoff with
// jitter on 429 and 5xx, no retry on other 4xx (403 means stop, per the CelesTrak rule in
// source-registry-method.md), and a hard body-size limit.
type Fetcher struct {
	Client      *http.Client
	Limiter     *rate.Limiter
	Quota       *Quota
	Headers     map[string]string
	MaxBody     int64
	Attempts    int
	BaseBackoff time.Duration
	MaxBackoff  time.Duration
	// OnResponse is called with every status code (0 for a transport error) for metrics.
	OnResponse func(code int)
	sleep      func(context.Context, time.Duration) error
}

func NewFetcher(s SourceConfig, onResponse func(int)) *Fetcher {
	return &Fetcher{
		Client:      &http.Client{Timeout: 30 * time.Second},
		Limiter:     rate.NewLimiter(rate.Every(s.RateLimit.Every.Duration), s.RateLimit.Burst),
		Quota:       NewQuota(s.Quota.Requests, s.Quota.Per.Duration, nil),
		Headers:     s.Headers,
		MaxBody:     s.MaxBodyBytes,
		Attempts:    4,
		BaseBackoff: 2 * time.Second,
		MaxBackoff:  min(s.MaxBackoff.Duration, 2*time.Minute),
		OnResponse:  onResponse,
	}
}

func sleepCtx(ctx context.Context, d time.Duration) error {
	t := time.NewTimer(d)
	defer t.Stop()
	select {
	case <-t.C:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// Get fetches rawURL. file:// URLs are read from disk (bulk sources mounted as files),
// with the file's modification time standing in for Last-Modified.
func (f *Fetcher) Get(ctx context.Context, rawURL string, cond Conditional) (Response, error) {
	if u, err := url.Parse(rawURL); err == nil && u.Scheme == "file" {
		return f.readFile(u.Path, cond)
	}
	return f.Do(ctx, http.MethodGet, rawURL, cond)
}

func (f *Fetcher) Do(ctx context.Context, method, rawURL string, cond Conditional) (Response, error) {
	sleep := f.sleep
	if sleep == nil {
		sleep = sleepCtx
	}
	backoff := f.BaseBackoff
	var last error
	for attempt := 1; attempt <= max(f.Attempts, 1); attempt++ {
		if err := f.Quota.Take(); err != nil {
			return Response{}, err
		}
		if err := f.Limiter.Wait(ctx); err != nil {
			return Response{}, err
		}
		req, err := http.NewRequestWithContext(ctx, method, rawURL, nil)
		if err != nil {
			return Response{}, err
		}
		for k, v := range f.Headers {
			req.Header.Set(k, v)
		}
		if cond.ETag != "" {
			req.Header.Set("If-None-Match", cond.ETag)
		}
		if cond.LastModified != "" {
			req.Header.Set("If-Modified-Since", cond.LastModified)
		}
		resp, err := f.Client.Do(req)
		if err != nil {
			f.observe(0)
			if ctx.Err() != nil {
				return Response{}, ctx.Err()
			}
			last = err
			if attempt < f.Attempts {
				if err := sleep(ctx, jitter(backoff)); err != nil {
					return Response{}, err
				}
				backoff = min(backoff*2, f.MaxBackoff)
			}
			continue
		}
		f.observe(resp.StatusCode)
		body, rerr := readLimited(resp.Body, f.MaxBody)
		resp.Body.Close()
		switch {
		case resp.StatusCode == http.StatusOK:
			if rerr != nil {
				return Response{}, rerr
			}
			return Response{Code: 200, Body: body, ETag: resp.Header.Get("ETag"), LastModified: resp.Header.Get("Last-Modified")}, nil
		case resp.StatusCode == http.StatusNotModified:
			return Response{Code: 304, ETag: cond.ETag, LastModified: cond.LastModified}, nil
		case resp.StatusCode == http.StatusTooManyRequests || resp.StatusCode >= 500:
			last = &StatusError{resp.StatusCode}
			if attempt < f.Attempts {
				wait := jitter(backoff)
				if ra, ok := retryAfter(resp.Header.Get("Retry-After"), time.Now()); ok {
					wait = min(ra, f.MaxBackoff)
				}
				if err := sleep(ctx, wait); err != nil {
					return Response{}, err
				}
				backoff = min(backoff*2, f.MaxBackoff)
			}
		default:
			return Response{}, &StatusError{resp.StatusCode} // 4xx: never hammer
		}
	}
	return Response{}, last
}

func (f *Fetcher) observe(code int) {
	if f.OnResponse != nil {
		f.OnResponse(code)
	}
}

func (f *Fetcher) readFile(path string, cond Conditional) (Response, error) {
	fi, err := os.Stat(path)
	if err != nil {
		return Response{}, err
	}
	lm := fi.ModTime().UTC().Format(http.TimeFormat)
	if cond.LastModified == lm {
		return Response{Code: 304, LastModified: lm}, nil
	}
	fh, err := os.Open(path)
	if err != nil {
		return Response{}, err
	}
	defer fh.Close()
	body, err := readLimited(fh, f.MaxBody)
	if err != nil {
		return Response{}, err
	}
	return Response{Code: 200, Body: body, LastModified: lm}, nil
}

func readLimited(r io.Reader, maxBytes int64) ([]byte, error) {
	if maxBytes <= 0 {
		maxBytes = 16 << 20
	}
	b, err := io.ReadAll(io.LimitReader(r, maxBytes+1))
	if err != nil {
		return nil, err
	}
	if int64(len(b)) > maxBytes {
		return nil, ErrBodyTooLarge
	}
	return b, nil
}

// retryAfter parses delta-seconds or an HTTP-date (RFC 9110 §10.2.3).
func retryAfter(v string, now time.Time) (time.Duration, bool) {
	if v == "" {
		return 0, false
	}
	if s, err := strconv.Atoi(v); err == nil && s >= 0 {
		return time.Duration(s) * time.Second, true
	}
	if t, err := http.ParseTime(v); err == nil {
		return max(t.Sub(now), 0), true
	}
	return 0, false
}

func jitter(d time.Duration) time.Duration {
	if d <= 0 {
		return 0
	}
	return d/2 + time.Duration(rand.Int64N(int64(d/2)+1))
}

// ErrorReason classifies an error for the ingest_fetch_errors_total "reason" label.
func ErrorReason(err error) string {
	var se *StatusError
	var qe *QuotaError
	switch {
	case errors.As(err, &se) && se.Code == 429:
		return "http_429"
	case errors.As(err, &se) && se.Code >= 500:
		return "http_5xx"
	case errors.As(err, &se):
		return "http_4xx"
	case errors.As(err, &qe):
		return "quota"
	case errors.Is(err, ErrBodyTooLarge):
		return "too_large"
	case errors.Is(err, errCheckpointRead):
		return "checkpoint"
	case errors.Is(err, errParse):
		return "parse"
	default:
		return "transport"
	}
}
