package layer

import (
	"encoding/json"
	"strings"
	"sync"
	"time"
)

// MaxCellRecord is the largest `{layer}.cells.z{z}` value (layer-snapshot.md §2); a zoom whose
// record would be larger is dropped from the manifest's `cells.zooms`.
const MaxCellRecord = 512 * 1024

// FeedStatus is the part of a `tel.feed.{source}.status` message the manifest copies
// (services/schemas/ingest/feed-status.v1.schema.json in omos PR #553).
type FeedStatus struct {
	Source        string
	FetchedAtMs   int64 // 0 = never
	Count         int64
	LastError     *string
	Stale         bool
	StaleAfterSec float64
	RecvAt        time.Time
}

// ParseFeedStatus reads one status message. fetchedAt is RFC 3339 (omos-ingest) or epoch ms.
func ParseFeedStatus(subject string, data []byte, recv time.Time) (FeedStatus, bool) {
	p := strings.Split(subject, ".")
	if len(p) != 4 || p[0] != "tel" || p[1] != "feed" || p[3] != "status" {
		return FeedStatus{}, false
	}
	var m struct {
		FetchedAt     any     `json:"fetchedAt"`
		Count         int64   `json:"count"`
		LastError     *string `json:"lastError"`
		Stale         bool    `json:"stale"`
		StaleAfterSec float64 `json:"staleAfterSec"`
	}
	if json.Unmarshal(data, &m) != nil {
		return FeedStatus{}, false
	}
	fs := FeedStatus{Source: p[2], Count: m.Count, LastError: m.LastError, Stale: m.Stale, StaleAfterSec: m.StaleAfterSec, RecvAt: recv}
	switch v := m.FetchedAt.(type) {
	case string:
		if t, err := time.Parse(time.RFC3339Nano, v); err == nil {
			fs.FetchedAtMs = t.UnixMilli()
		}
	case float64:
		fs.FetchedAtMs = int64(v)
	}
	return fs, true
}

// feedOut is one `freshness.feeds[]` entry.
type feedOut struct {
	Source    string  `json:"source"`
	FetchedAt *int64  `json:"fetchedAt"`
	Count     int64   `json:"count"`
	LastError *string `json:"lastError"`
	Stale     bool    `json:"stale"`
}

type freshness struct {
	Status string    `json:"status"`
	LimitS int       `json:"limit_s"`
	Feeds  []feedOut `json:"feeds"`
}

// Manifest is the `{layer}.meta` value (layer-snapshot.md §2). `filters` is additive: the
// exact ENT filter set for layers with more than one entity type.
type Manifest struct {
	V             int        `json:"v"`
	Layer         string     `json:"layer"`
	Title         string     `json:"title"`
	Tier          string     `json:"tier"`
	EntityTypes   []string   `json:"entity_types"`
	Filter        string     `json:"filter"`
	Filters       []string   `json:"filters"`
	Site          string     `json:"site"`
	Count         int        `json:"count"`
	AsOf          int64      `json:"as_of"`
	EntSeq        uint64     `json:"ent_seq"`
	Cap           int        `json:"cap"`
	Cells         *CellsCfg  `json:"cells"`
	Object        *ObjectRef `json:"object"`
	Freshness     freshness  `json:"freshness"`
	DefaultStaleS int        `json:"default_stale_s"`
	Marking       string     `json:"marking"`
	Style         string     `json:"style"`
}

// computeFreshness applies layer-snapshot.md §2 "Field rules": fresh when every feed is within
// its limit, degraded when some are, stale when none are, down when ENT is not readable. A
// feed is within its limit when its last successful fetch is no older than the limit (the
// feed's staleAfterSec, else the layer's limit_s) and the collector does not flag it stale.
func computeFreshness(d *Def, feeds map[string]FeedStatus, entDown bool, now time.Time) freshness {
	f := freshness{LimitS: d.LimitS, Feeds: []feedOut{}}
	fresh := 0
	for _, src := range d.Feeds {
		st, ok := feeds[src]
		out := feedOut{Source: src, Stale: true}
		if ok {
			out.Count, out.LastError = st.Count, st.LastError
			limit := float64(d.LimitS)
			if st.StaleAfterSec > 0 {
				limit = st.StaleAfterSec
			}
			if st.FetchedAtMs > 0 {
				fa := st.FetchedAtMs
				out.FetchedAt = &fa
				age := now.Sub(time.UnixMilli(fa)).Seconds()
				out.Stale = st.Stale || age > limit
			}
		} else {
			msg := "no status received"
			out.LastError = &msg
		}
		if !out.Stale {
			fresh++
		}
		f.Feeds = append(f.Feeds, out)
	}
	switch {
	case entDown:
		f.Status = "down"
	case fresh == len(d.Feeds):
		f.Status = "fresh"
	case fresh == 0:
		f.Status = "stale"
	default:
		f.Status = "degraded"
	}
	return f
}

// manifestState is what was last written for a layer.
type manifestState struct {
	written   time.Time
	cellsAt   time.Time
	zooms     []int
	signature string
	count     int
	json      []byte
}

// signature is the part of a manifest whose change is pushed to watches at once.
func (m *Manifest) signature() string {
	b, _ := json.Marshal(struct {
		F freshness
		O *ObjectRef
		C *CellsCfg
	}{m.Freshness, m.Object, m.Cells})
	return string(b)
}

// Manifests keeps the current manifest of every layer for svc.layer.list/meta and the watch path.
type Manifests struct {
	mu    sync.RWMutex
	state []manifestState
	objs  []*ObjectRef
	feeds map[string]FeedStatus
}

// NewManifests builds the cache.
func NewManifests(defs []Def) *Manifests {
	m := &Manifests{state: make([]manifestState, len(defs)), objs: make([]*ObjectRef, len(defs)), feeds: map[string]FeedStatus{}}
	for i := range defs {
		m.objs[i] = defs[i].Object
		if defs[i].Cells != nil {
			m.state[i].zooms = append([]int(nil), defs[i].Cells.Zooms...)
		}
	}
	return m
}

// SetFeed records a feed status.
func (m *Manifests) SetFeed(fs FeedStatus) {
	m.mu.Lock()
	m.feeds[fs.Source] = fs
	m.mu.Unlock()
}

// Zooms returns the current cell zooms of layer li.
func (m *Manifests) Zooms(li int) []int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return append([]int(nil), m.state[li].zooms...)
}

// Object returns the current object of static layer li.
func (m *Manifests) Object(li int) *ObjectRef {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.objs[li]
}

// JSON returns the last written manifest of layer li (nil before the first write).
func (m *Manifests) JSON(li int) []byte {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.state[li].json
}

// build makes the manifest of layer li.
func (m *Manifests) build(d *Def, li int, site string, count int, seq uint64, entDown bool, now time.Time) *Manifest {
	m.mu.RLock()
	defer m.mu.RUnlock()
	man := &Manifest{V: 1, Layer: d.ID, Title: d.Title, Tier: d.Tier, EntityTypes: append(append([]string{}, d.EntityTypes...), d.CotOnlyTypes...),
		Filter: d.Filter(), Filters: d.Filters(), Site: site, Count: count, AsOf: now.UnixMilli(), EntSeq: seq,
		Cap: d.Cap, Object: m.objs[li], Freshness: computeFreshness(d, m.feeds, entDown, now),
		DefaultStaleS: d.DefaultStaleS, Marking: d.Marking, Style: d.Style}
	if man.EntityTypes == nil {
		man.EntityTypes = []string{}
	}
	if man.Filters == nil {
		man.Filters = []string{}
	}
	if d.Cells != nil {
		man.Cells = &CellsCfg{Zooms: append([]int{}, m.state[li].zooms...), BelowZoom: d.Cells.BelowZoom}
	}
	return man
}

// cellRecord is the `{layer}.cells.z{z}` value.
type cellRecord struct {
	V      int    `json:"v"`
	Layer  string `json:"layer"`
	Z      int    `json:"z"`
	AsOf   int64  `json:"as_of"`
	EntSeq uint64 `json:"ent_seq"`
	Scheme string `json:"scheme"`
	Cells  []cell `json:"cells"`
}

// encodeCells returns the records per zoom that fit MaxCellRecord, and the zooms kept.
func encodeCells(layer string, cells map[int][]cell, zooms []int, seq uint64, now time.Time, limit int) (map[int][]byte, []int) {
	out := map[int][]byte{}
	var kept []int
	for _, z := range zooms {
		b, _ := json.Marshal(cellRecord{V: 1, Layer: layer, Z: z, AsOf: now.UnixMilli(), EntSeq: seq, Scheme: "webmercator-tile", Cells: append([]cell{}, cells[z]...)})
		if len(b) > limit {
			continue
		}
		out[z] = b
		kept = append(kept, z)
	}
	if kept == nil {
		kept = []int{}
	}
	return out, kept
}

// SiteStatus is the `fleet-state` `{site}` value (layer-snapshot.md §2).
type SiteStatus struct {
	V           int                `json:"v"`
	Site        string             `json:"site"`
	Connected   bool               `json:"connected"`
	LastContact int64              `json:"last_contact"`
	SourceLagS  map[string]float64 `json:"source_lag_s"`
}
