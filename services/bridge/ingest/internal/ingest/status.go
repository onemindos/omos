package ingest

import (
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

// Status is one source's freshness record. It is published on tel.feed.{id}.status by
// the lease holder and backs the ingest_feed_* Prometheus metrics.
type Status struct {
	Source        string     `json:"source"`
	Kind          string     `json:"kind"`
	Collector     string     `json:"collector"`
	Instance      string     `json:"instance"`
	Leader        bool       `json:"leader"`
	FetchedAt     *time.Time `json:"fetchedAt"` // last successful fetch (or stream message)
	Count         int        `json:"count"`     // records in the last successful fetch
	Published     int        `json:"published"` // new records acked in the last fetch
	LastError     *string    `json:"lastError"` // null after a success
	ErrorAt       *time.Time `json:"errorAt"`
	StaleAfterSec float64    `json:"staleAfterSec"`
	AgeSec        *float64   `json:"ageSec"`
	Stale         bool       `json:"stale"`
	At            time.Time  `json:"at"`
}

// StatusBoard holds every source's status for one process. Safe for concurrent use.
type StatusBoard struct {
	mu       sync.Mutex
	now      func() time.Time
	instance string
	s        map[string]*Status
}

func NewStatusBoard(instance string, now func() time.Time) *StatusBoard {
	if now == nil {
		now = time.Now
	}
	return &StatusBoard{instance: instance, now: now, s: map[string]*Status{}}
}

func (b *StatusBoard) Add(id, kind string, staleAfter time.Duration) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.s[id] = &Status{Source: id, Kind: kind, Collector: "omos-ingest", Instance: b.instance, StaleAfterSec: staleAfter.Seconds()}
}

func (b *StatusBoard) OK(id string, count, published int) {
	b.mu.Lock()
	defer b.mu.Unlock()
	s := b.s[id]
	now := b.now()
	s.FetchedAt, s.Count, s.Published, s.LastError, s.ErrorAt = &now, count, published, nil, nil
}

func (b *StatusBoard) Fail(id string, err error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	s := b.s[id]
	now := b.now()
	e := err.Error()
	if len(e) > 300 {
		e = e[:300]
	}
	s.LastError, s.ErrorAt = &e, &now
}

func (b *StatusBoard) SetLeader(id string, v bool) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.s[id].Leader = v
}

// Snapshot returns a copy with age and stale computed now. A source that never
// fetched is stale.
func (b *StatusBoard) Snapshot() []Status {
	b.mu.Lock()
	defer b.mu.Unlock()
	now := b.now()
	out := make([]Status, 0, len(b.s))
	for _, s := range b.s {
		c := *s
		c.At = now
		if c.FetchedAt != nil {
			age := now.Sub(*c.FetchedAt).Seconds()
			c.AgeSec = &age
			c.Stale = age > c.StaleAfterSec
		} else {
			c.Stale = true
		}
		out = append(out, c)
	}
	return out
}

func (b *StatusBoard) Get(id string) Status {
	for _, s := range b.Snapshot() {
		if s.Source == id {
			return s
		}
	}
	return Status{}
}

// Metrics are the Prometheus series for every source (label "source").
type Metrics struct {
	Reg         *prometheus.Registry
	Published   *prometheus.CounterVec
	Duplicates  *prometheus.CounterVec
	Rejected    *prometheus.CounterVec
	Retries     *prometheus.CounterVec
	FetchErrors *prometheus.CounterVec
	Requests    *prometheus.CounterVec
	Dropped     *prometheus.CounterVec
	Checkpoints *prometheus.CounterVec
	QueueDepth  *prometheus.GaugeVec
	Quota       *prometheus.GaugeVec
	board       *StatusBoard
}

func NewMetrics(board *StatusBoard) *Metrics {
	cv := func(name, help string, labels ...string) *prometheus.CounterVec {
		return prometheus.NewCounterVec(prometheus.CounterOpts{Name: name, Help: help}, append([]string{"source"}, labels...))
	}
	m := &Metrics{
		Reg:         prometheus.NewRegistry(),
		Published:   cv("ingest_published_total", "Records stored by JetStream (PubAck, not duplicate)."),
		Duplicates:  cv("ingest_duplicates_total", "Publishes JetStream acked as duplicates of a Nats-Msg-Id."),
		Rejected:    cv("ingest_rejected_total", "Records sent to evt.sys.{source}.rejected."),
		Retries:     cv("ingest_publish_retries_total", "JetStream publish attempts that failed and were retried."),
		FetchErrors: cv("ingest_fetch_errors_total", "Failed upstream fetches or stream reads.", "reason"),
		Requests:    cv("ingest_upstream_requests_total", "Upstream HTTP requests by status code.", "code"),
		Dropped:     cv("ingest_dropped_total", "Stream messages dropped because the bounded queue was full."),
		Checkpoints: cv("ingest_checkpoints_total", "Checkpoints written to KV after acks.", "result"),
		QueueDepth:  prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: "ingest_queue_depth", Help: "Messages waiting in a stream source's bounded queue."}, []string{"source"}),
		Quota:       prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: "ingest_quota_remaining", Help: "Upstream requests left in the current quota window."}, []string{"source"}),
		board:       board,
	}
	m.Reg.MustRegister(m.Published, m.Duplicates, m.Rejected, m.Retries, m.FetchErrors, m.Requests,
		m.Dropped, m.Checkpoints, m.QueueDepth, m.Quota, freshness{board}, prometheus.NewGoCollector(), prometheus.NewProcessCollector(prometheus.ProcessCollectorOpts{}))
	return m
}

// freshness exports the status board at scrape time, so age is always current.
type freshness struct{ b *StatusBoard }

var (
	descLastSuccess = prometheus.NewDesc("ingest_feed_last_success_timestamp_seconds", "Unix time of the last successful fetch (fetchedAt). Absent until the first success.", []string{"source"}, nil)
	descAge         = prometheus.NewDesc("ingest_feed_age_seconds", "Seconds since the last successful fetch; -1 before the first.", []string{"source"}, nil)
	descStaleAfter  = prometheus.NewDesc("ingest_feed_stale_after_seconds", "The source's freshness limit from sources.yaml.", []string{"source"}, nil)
	descStale       = prometheus.NewDesc("ingest_feed_stale", "1 when the source is older than its limit or never fetched.", []string{"source"}, nil)
	descCount       = prometheus.NewDesc("ingest_feed_records", "Records in the last successful fetch.", []string{"source"}, nil)
	descUp          = prometheus.NewDesc("ingest_feed_up", "1 when the last fetch attempt succeeded (lastError is null).", []string{"source"}, nil)
	descLastErr     = prometheus.NewDesc("ingest_feed_last_error_timestamp_seconds", "Unix time of the last error.", []string{"source"}, nil)
	descLeader      = prometheus.NewDesc("ingest_feed_leader", "1 when this instance holds the source's lease.", []string{"source", "instance"}, nil)
)

func (f freshness) Describe(ch chan<- *prometheus.Desc) {
	for _, d := range []*prometheus.Desc{descLastSuccess, descAge, descStaleAfter, descStale, descCount, descUp, descLastErr, descLeader} {
		ch <- d
	}
}

func (f freshness) Collect(ch chan<- prometheus.Metric) {
	g := func(d *prometheus.Desc, v float64, l ...string) {
		ch <- prometheus.MustNewConstMetric(d, prometheus.GaugeValue, v, l...)
	}
	for _, s := range f.b.Snapshot() {
		// Only the lease holder speaks for a source's freshness; a standby has no data.
		g(descLeader, b2f(s.Leader), s.Source, s.Instance)
		if !s.Leader && s.FetchedAt == nil {
			continue
		}
		g(descStaleAfter, s.StaleAfterSec, s.Source)
		g(descStale, b2f(s.Stale), s.Source)
		g(descUp, b2f(s.LastError == nil && s.FetchedAt != nil), s.Source)
		age := -1.0
		if s.FetchedAt != nil {
			age = *s.AgeSec
			g(descLastSuccess, float64(s.FetchedAt.UnixMilli())/1000, s.Source)
			g(descCount, float64(s.Count), s.Source)
		}
		g(descAge, age, s.Source)
		if s.ErrorAt != nil {
			g(descLastErr, float64(s.ErrorAt.UnixMilli())/1000, s.Source)
		}
	}
}

func b2f(b bool) float64 {
	if b {
		return 1
	}
	return 0
}
