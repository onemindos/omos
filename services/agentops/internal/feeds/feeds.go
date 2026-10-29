// Package feeds keeps the latest `tel.feed.{source}.status` heartbeat of every collector
// (omos-ingest feed-status.v1, services/schemas/ingest/feed-status.v1.schema.json in PR #553)
// and turns it into intelligence gaps: feeds that are stale, failing, or silent.
package feeds

import (
	"encoding/json"
	"sort"
	"strings"
	"sync"
	"time"
)

// Status is the subset of feed-status.v1 this package reads.
type Status struct {
	Source        string   `json:"source"`
	FetchedAt     *string  `json:"fetchedAt"`
	Count         int      `json:"count"`
	LastError     *string  `json:"lastError"`
	StaleAfterSec float64  `json:"staleAfterSec"`
	AgeSec        *float64 `json:"ageSec"`
	Stale         bool     `json:"stale"`
	At            string   `json:"at"`

	received time.Time
}

// Gap is one degraded feed.
type Gap struct {
	Source    string   `json:"source"`
	Reason    string   `json:"reason"` // stale | failing | silent
	Detail    string   `json:"detail"`
	AgeS      *float64 `json:"age_s,omitempty"`
	LastError string   `json:"last_error,omitempty"`
	FetchedAt string   `json:"fetched_at,omitempty"`
	Layers    []string `json:"layers,omitempty"`
}

// Board holds feed statuses. Safe for concurrent use.
type Board struct {
	mu       sync.Mutex
	now      func() time.Time
	s        map[string]*Status
	expected map[string]float64 // source -> staleAfterSec, for feeds that must report
	layers   map[string][]string
	silentK  float64
}

// NewBoard returns a board. expected lists sources that must report (missing ones are
// "silent" gaps); layers maps a source to the layers it feeds (for "weakens" in a brief).
func NewBoard(now func() time.Time, expected map[string]float64, layers map[string][]string) *Board {
	if now == nil {
		now = time.Now
	}
	return &Board{now: now, s: map[string]*Status{}, expected: expected, layers: layers, silentK: 3}
}

// SourceOf returns the source token of `tel.feed.{source}.status`, or "".
func SourceOf(subject string) string {
	t := strings.Split(subject, ".")
	if len(t) == 4 && t[0] == "tel" && t[1] == "feed" && t[3] == "status" {
		return t[2]
	}
	return ""
}

// Update records one status message.
func (b *Board) Update(subject string, data []byte) bool {
	src := SourceOf(subject)
	if src == "" {
		return false
	}
	var st Status
	if json.Unmarshal(data, &st) != nil {
		return false
	}
	st.Source = src
	st.received = b.now()
	b.mu.Lock()
	b.s[src] = &st
	b.mu.Unlock()
	return true
}

// Gaps returns every degraded feed, sorted by source.
func (b *Board) Gaps() []Gap {
	b.mu.Lock()
	defer b.mu.Unlock()
	now := b.now()
	out := []Gap{}
	seen := map[string]bool{}
	for src, st := range b.s {
		seen[src] = true
		limit := st.StaleAfterSec
		if limit <= 0 {
			limit = 300
		}
		since := now.Sub(st.received).Seconds()
		g := Gap{Source: src, AgeS: st.AgeSec, Layers: b.layers[src]}
		if st.FetchedAt != nil {
			g.FetchedAt = *st.FetchedAt
		}
		if st.LastError != nil {
			g.LastError = *st.LastError
		}
		switch {
		case since > b.silentK*maxf(limit, 10):
			// The collector itself stopped reporting: its last word may say "fresh".
			g.Reason, g.Detail = "silent", "no status heartbeat for "+fmtSecs(since)
		case st.Stale:
			g.Reason, g.Detail = "stale", "stale beyond its "+fmtSecs(limit)+" freshness limit"
			if st.AgeSec != nil {
				g.Detail = "last good data " + fmtSecs(*st.AgeSec) + " old, limit " + fmtSecs(limit)
			}
		case st.LastError != nil:
			g.Reason, g.Detail = "failing", "last fetch failed: "+*st.LastError
		default:
			continue
		}
		out = append(out, g)
	}
	for src, limit := range b.expected {
		if !seen[src] {
			out = append(out, Gap{Source: src, Reason: "silent", Detail: "never reported (expected, limit " + fmtSecs(limit) + ")", Layers: b.layers[src]})
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Source < out[j].Source })
	return out
}

func maxf(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}

func fmtSecs(s float64) string {
	return (time.Duration(s) * time.Second).Round(time.Second).String()
}
