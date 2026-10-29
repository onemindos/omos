package layer

import (
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

// change is one entry of the bounded replay ring that `svc.layer.watch` uses to serve
// `from_seq` (layer-delta.md §2.1: a bounded recent window, proposed 5 minutes).
type change struct {
	seq    uint64 // ENT stream sequence (for a sweep: the last sequence seen at sweep time)
	key    string
	reason string // remove reason; "" for an upsert
	sweep  bool   // removal by the expiry sweep (no ENT message of its own)
	atMs   int64
}

// Store holds the latest render state of every ENT entity, the replay ring and the watches.
//
// Lock order: Store.mu, then Store.wmu, then Watch.mu, then Store.cellMu. Never take an
// earlier lock while holding a later one, and never hold two Watch.mu at once.
type Store struct {
	mu       sync.RWMutex
	defs     []*Def
	byType   map[string][]int
	ents     map[string]*entry
	perLayer []map[string]*entry
	lastSeq  uint64
	ready    bool
	down     bool
	bad      uint64 // ENT records that could not be read

	ring        []change
	floor       uint64 // every change after floor is in the ring
	floorStrict bool   // a sweep removal at sequence floor was evicted (client state at floor is not enough)
	window      time.Duration
	maxRing     int

	origin string
	now    func() time.Time

	version uint64 // bumped on every change of ents (Store.mu write)

	wmu     sync.RWMutex
	watches map[string]*Watch

	units       atomic.Int64 // sum of Watch.units
	budgetUnits int64        // service-wide watch budget in units (0: none)

	cellMu    sync.Mutex
	cellVer   uint64
	cellCache map[string]map[[2]int]cell
}

// SetWatchBudget sets the service-wide memory budget of all watches, in bytes (0: none).
// When the total passes it, the watch holding the most is ended with `resnapshot`.
func (s *Store) SetWatchBudget(bytes int64) {
	s.wmu.Lock()
	s.budgetUnits = bytes / UnitBytes
	s.wmu.Unlock()
}

// WatchBytes returns the estimated memory held by all watches, in bytes.
func (s *Store) WatchBytes() int64 { return s.units.Load() * UnitBytes }

// enforceBudgetLocked ends the largest watches until the total is within the budget.
// Store.wmu (read or write) is held; no Watch.mu is held.
func (s *Store) enforceBudgetLocked() {
	for s.budgetUnits > 0 && s.units.Load() > s.budgetUnits {
		var big *Watch
		var bu int64
		for _, w := range s.watches {
			if u := w.units.Load(); u > bu {
				big, bu = w, u
			}
		}
		if big == nil {
			return
		}
		big.mu.Lock()
		if big.closed || big.endCode != "" {
			big.clearLocked() // already ending: make sure its units are returned
		} else {
			big.endLocked(ErrResnapshot, "layer service watch memory budget reached; re-snapshot")
		}
		big.mu.Unlock()
	}
}

// sharedCellsLocked returns the cell grid of w's view. Watches with the same layer, cell
// zoom, AOI and caller share one grid per store version, so N identical views cost one scan.
// Store.mu (read) and w.mu are held; the returned map must not be modified.
func (s *Store) sharedCellsLocked(w *Watch) map[[2]int]cell {
	vis := w.caller.ID
	if w.caller.Trusted {
		vis = "*"
	}
	key := strconv.Itoa(w.li) + "|" + strconv.Itoa(w.cellZ) + "|" + w.aoiKey + "|" + vis
	s.cellMu.Lock()
	if s.cellVer != s.version || s.cellCache == nil {
		s.cellVer, s.cellCache = s.version, map[string]map[[2]int]cell{}
	}
	if m, ok := s.cellCache[key]; ok {
		s.cellMu.Unlock()
		return m
	}
	s.cellMu.Unlock()
	g := newCellGrid(w.cellZ)
	for _, e := range s.perLayer[w.li] {
		if w.visible(e) {
			g.add(e.r.Lat, e.r.Lon)
		}
	}
	m := g.asMap()
	s.cellMu.Lock()
	if s.cellVer == s.version {
		s.cellCache[key] = m
	}
	s.cellMu.Unlock()
	return m
}

// NewStore builds an empty store for the layer definitions.
func NewStore(defs []Def, origin string, window time.Duration, maxRing int, now func() time.Time) *Store {
	s := &Store{byType: map[string][]int{}, ents: map[string]*entry{}, window: window, maxRing: maxRing,
		origin: origin, now: now, watches: map[string]*Watch{}}
	for i := range defs {
		d := defs[i]
		s.defs = append(s.defs, &d)
		s.perLayer = append(s.perLayer, map[string]*entry{})
		for _, t := range d.EntityTypes {
			s.byType[t] = append(s.byType[t], i)
		}
		for _, t := range d.CotOnlyTypes {
			s.byType[t] = append(s.byType[t], i)
		}
	}
	return s
}

// layerIndex returns the index of layer id, or -1.
func (s *Store) layerIndex(id string) int {
	for i, d := range s.defs {
		if d.ID == id {
			return i
		}
	}
	return -1
}

// member reports whether e belongs to layer d.
func member(d *Def, e *entry) bool {
	if e == nil {
		return false
	}
	ok, cotOnly := d.HasType(e.r.EntityType)
	return ok && (!cotOnly || e.fromCot)
}

func parseSubject(subj string) (typ, id string, ok bool) {
	p := strings.Split(subj, ".")
	if len(p) != 4 || p[0] != "ent" || p[3] != "state" || !tokenRe.MatchString(p[1]) || !tokenRe.MatchString(p[2]) {
		return "", "", false
	}
	return p[1], p[2], true
}

// Apply processes one ENT message in stream order. marker is the Nats-Marker-Reason header.
func (s *Store) Apply(subject string, data []byte, marker string, seq uint64) {
	typ, id, ok := parseSubject(subject)
	s.mu.Lock()
	defer s.mu.Unlock()
	if seq != 0 && seq <= s.lastSeq {
		return // redelivery after an ordered-consumer reset
	}
	if seq != 0 {
		s.lastSeq = seq
	}
	if !ok {
		s.bad++
		return
	}
	key := typ + "." + id
	old := s.ents[key]
	if len(data) == 0 || marker != "" {
		if old != nil {
			s.removeLocked(old, markerReason(marker), false)
		}
		return
	}
	e, in, err := parseEntity(typ, id, data, s.origin)
	if err != nil {
		s.bad++
		return
	}
	if !in.Live {
		if old != nil {
			s.removeLocked(old, removalReason(in), false)
		}
		return
	}
	if e.r.ExpiryT != nil && *e.r.ExpiryT <= s.now().UnixMilli() {
		if old != nil {
			s.removeLocked(old, "expired", false)
		}
		return
	}
	s.ents[key] = e
	s.version++
	for _, i := range s.byType[typ] {
		if member(s.defs[i], e) {
			s.perLayer[i][key] = e
		} else {
			delete(s.perLayer[i], key)
		}
	}
	s.recordLocked(change{seq: s.lastSeq, key: key})
	s.dispatchLocked(key, old, e, "")
}

func (s *Store) removeLocked(old *entry, reason string, sweep bool) {
	delete(s.ents, old.key)
	s.version++
	for _, i := range s.byType[old.r.EntityType] {
		delete(s.perLayer[i], old.key)
	}
	s.recordLocked(change{seq: s.lastSeq, key: old.key, reason: reason, sweep: sweep})
	s.dispatchLocked(old.key, old, nil, reason)
}

// Sweep removes entities whose expiry_t has passed. It covers ENT without per-message TTL
// (no delete markers) and the gap between expiry and the marker.
func (s *Store) Sweep() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	now := s.now().UnixMilli()
	n := 0
	for _, e := range s.ents {
		if e.r.ExpiryT != nil && *e.r.ExpiryT <= now {
			s.removeLocked(e, "expired", true)
			n++
		}
	}
	s.evictLocked()
	return n
}

func (s *Store) recordLocked(c change) {
	if !s.ready {
		return
	}
	c.atMs = s.now().UnixMilli()
	s.ring = append(s.ring, c)
	s.evictLocked()
}

func (s *Store) evictLocked() {
	cut := s.now().Add(-s.window).UnixMilli()
	i := 0
	for i < len(s.ring) && (len(s.ring)-i > s.maxRing || s.ring[i].atMs < cut) {
		if c := s.ring[i]; c.seq > s.floor {
			s.floor, s.floorStrict = c.seq, c.sweep
		} else if c.seq == s.floor && c.sweep {
			s.floorStrict = true
		}
		i++
	}
	if i > 0 {
		s.ring = s.ring[i:]
		if cap(s.ring) > 4*len(s.ring)+1024 {
			s.ring = append([]change(nil), s.ring...)
		}
	}
}

// MarkReady ends the catch-up: the replay ring starts at the current sequence.
func (s *Store) MarkReady() {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.ready {
		return
	}
	s.ready, s.down = true, false
	s.floor, s.floorStrict = s.lastSeq, false
}

// SetLastSeq sets the sequence before catch-up (an empty ENT whose last sequence is known).
func (s *Store) SetLastSeq(seq uint64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if seq > s.lastSeq {
		s.lastSeq = seq
	}
}

// SetDown marks ENT readable or not.
func (s *Store) SetDown(down bool) {
	s.mu.Lock()
	s.down = down
	s.mu.Unlock()
}

// State returns readiness, down and the last ENT sequence.
func (s *Store) State() (ready, down bool, lastSeq uint64) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.ready, s.down, s.lastSeq
}

// LastSeq returns the last ENT sequence applied.
func (s *Store) LastSeq() uint64 {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.lastSeq
}

// Counts returns the live entity count per layer index and the last sequence.
func (s *Store) Counts() ([]int, uint64) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]int, len(s.defs))
	for i := range s.perLayer {
		out[i] = len(s.perLayer[i])
	}
	return out, s.lastSeq
}

// Bad returns the number of unreadable ENT records.
func (s *Store) Bad() uint64 {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.bad
}

// Cells aggregates the whole layer at each zoom (for KV `{layer}.cells.z{z}`).
func (s *Store) Cells(li int, zooms []int) (map[int][]cell, uint64) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	grids := make([]*cellGrid, len(zooms))
	for i, z := range zooms {
		grids[i] = newCellGrid(z)
	}
	for _, e := range s.perLayer[li] {
		for _, g := range grids {
			g.add(e.r.Lat, e.r.Lon)
		}
	}
	out := map[int][]cell{}
	for _, g := range grids {
		out[g.z] = g.list()
	}
	return out, s.lastSeq
}

// replayValid reports whether a watch can start at fromSeq (client state is fromSeq-1).
func (s *Store) replayValidLocked(fromSeq uint64) bool {
	at := fromSeq - 1
	if at > s.lastSeq {
		return false
	}
	if at > s.floor {
		return true
	}
	return at == s.floor && !s.floorStrict
}

// dispatchLocked hands one change to every watch of an affected layer.
func (s *Store) dispatchLocked(key string, old, nw *entry, reason string) {
	s.wmu.RLock()
	defer s.wmu.RUnlock()
	if len(s.watches) == 0 {
		return
	}
	for _, w := range s.watches {
		d := s.defs[w.li]
		if !member(d, old) && !member(d, nw) {
			continue
		}
		w.onChange(key, old, nw, reason)
	}
	s.enforceBudgetLocked()
}

// WatchCount returns the number of open watches.
func (s *Store) WatchCount() int {
	s.wmu.RLock()
	defer s.wmu.RUnlock()
	return len(s.watches)
}
