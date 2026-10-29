package layer

import (
	"bytes"
	"encoding/json"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/nats-io/nuid"
)

// Watch limits (layer-delta.md §2).
const (
	MaxBatchBytes = 256 * 1024
	MaxBatchOps   = 1000
	maxCellsPerOp = 2000
)

// UnitBytes is the estimated memory one watch unit costs: one key in a watch's held set, one
// pending op (a key, an op name and a pointer to the shared entry; the op is rendered at batch
// time, not copied per watch), one queued cell op or one remembered cell. The service-wide
// watch budget (Config.WatchMemBudget) is counted in these units.
const UnitBytes = 160

// pendOp is one queued op. Entity ops point at the shared, immutable entry and are rendered
// when the batch is built; remove and exit ops carry only the reason.
type pendOp struct {
	op     string
	e      *entry
	reason string
}

// Watch is one `svc.layer.watch` subscription.
type Watch struct {
	id     string
	li     int
	def    *Def
	caller Caller
	s      *Store
	units  atomic.Int64 // held + pending + cell state, in UnitBytes (read without mu)

	mu         sync.Mutex
	reply      string
	area       *area
	aoiKey     string // canonical AOI request (shared cell grids)
	zoom       int
	cellsMode  bool
	cellZ      int
	full       bool
	rateHz     int
	leaseUntil time.Time
	held       map[string]struct{}
	order      []string
	pend       map[string]pendOp
	cellOps    [][]byte
	cellsDirty bool
	lastCells  map[[2]int]cell
	lastCellsT time.Time
	meta       json.RawMessage
	n          uint64
	nextFrom   uint64
	lastSent   time.Time
	closed     bool
	endCode    string // set when the service ends the watch (too_many, resnapshot)
	endDetail  string
	wake       chan struct{}
}

// visible reports whether entry e is shown to this watch.
func (w *Watch) visible(e *entry) bool {
	return member(w.def, e) && w.area.contains(e.r.Lat, e.r.Lon) && w.caller.CanSee(e.tlp, e.releasable)
}

func (w *Watch) entityJSON(e *entry) []byte {
	if w.full {
		return e.raw
	}
	return e.renderJSON
}

func opEntity(op string, e []byte) []byte {
	b := make([]byte, 0, len(e)+24)
	b = append(b, `{"op":"`...)
	b = append(b, op...)
	b = append(b, `","e":`...)
	b = append(b, e...)
	return append(b, '}')
}

func opID(op, id, reason string) []byte {
	b, _ := json.Marshal(struct {
		Op     string `json:"op"`
		ID     string `json:"id"`
		Reason string `json:"reason"`
	}{op, id, reason})
	return b
}

// queue stores the newest op for key (coalescing: one op per entity per batch).
func (w *Watch) queue(key string, op pendOp) {
	if _, ok := w.pend[key]; !ok {
		w.order = append(w.order, key)
	}
	w.pend[key] = op
}

func (w *Watch) render(key string, op pendOp) []byte {
	if op.e != nil {
		return opEntity(op.op, w.entityJSON(op.e))
	}
	return opID(op.op, idOf(key), op.reason)
}

// accountLocked updates this watch's units and the service-wide total. w.mu is held.
func (w *Watch) accountLocked() {
	n := int64(len(w.held) + len(w.pend) + len(w.cellOps) + len(w.lastCells))
	if d := n - w.units.Load(); d != 0 {
		w.units.Store(n)
		if w.s != nil {
			w.s.units.Add(d)
		}
	}
}

// clearLocked drops everything the watch holds. w.mu is held.
func (w *Watch) clearLocked() {
	w.held, w.pend, w.order, w.cellOps, w.lastCells = map[string]struct{}{}, map[string]pendOp{}, nil, nil, nil
	w.accountLocked()
}

// endLocked ends the watch from the service side: its memory is released now and the next
// batch is the error reply (layer-delta.md §2.1 error codes). w.mu is held.
func (w *Watch) endLocked(code, detail string) {
	if w.closed || w.endCode != "" {
		return
	}
	w.endCode, w.endDetail = code, detail
	w.clearLocked()
	w.signal()
}

// releaseLocked closes the watch and returns its units (idempotent). w.mu is held.
func (w *Watch) releaseLocked() {
	w.closed = true
	w.clearLocked()
}

func (w *Watch) capDetail() string {
	return "more than the layer cap (" + strconv.Itoa(w.def.Cap) + ") entities in the AOI; zoom out for cells or narrow the AOI"
}

// countVisibleLocked counts the entities this watch would hold in entities mode, stopping
// once the count passes limit. Store.mu is held.
func (w *Watch) countVisibleLocked(s *Store, limit int) int {
	n := 0
	for _, e := range s.perLayer[w.li] {
		if w.visible(e) {
			if n++; n > limit {
				break
			}
		}
	}
	return n
}

func idOf(key string) string {
	for i := 0; i < len(key); i++ {
		if key[i] == '.' {
			return key[i+1:]
		}
	}
	return key
}

// onChange is called with Store.mu held (write).
func (w *Watch) onChange(key string, old, nw *entry, reason string) {
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.closed || w.endCode != "" {
		return
	}
	if w.cellsMode {
		w.cellsDirty = true
		return
	}
	_, held := w.held[key]
	switch {
	case nw != nil && w.visible(nw):
		if !held && len(w.held) >= w.def.Cap {
			// The view grew past the cap: no layer sends more than its cap (FEATURES A14).
			w.endLocked(ErrTooMany, w.capDetail())
			return
		}
		op := "upsert"
		if !held && old != nil {
			op = "enter"
		}
		w.queue(key, pendOp{op: op, e: nw})
		w.held[key] = struct{}{}
	case held:
		switch {
		case nw == nil:
			w.queue(key, pendOp{op: "remove", reason: reason})
		case member(w.def, nw) && !w.caller.CanSee(nw.tlp, nw.releasable):
			w.queue(key, pendOp{op: "remove", reason: "released_out"})
		default:
			w.queue(key, pendOp{op: "exit", reason: "left_aoi"})
		}
		delete(w.held, key)
	}
	w.accountLocked()
	w.signal()
}

func (w *Watch) signal() {
	select {
	case w.wake <- struct{}{}:
	default:
	}
}

// resetViewLocked recomputes what the client holds after an AOI or zoom change, queueing
// enter/exit for the difference (layer-delta.md §2.1). wasCells means the client held cells,
// not entities, before; prevZ is the cell zoom it held. Store.mu (read) and w.mu are held.
func (w *Watch) resetViewLocked(s *Store, wasCells bool, prevZ int) {
	if w.cellsMode {
		if !wasCells {
			for key := range w.held {
				w.queue(key, pendOp{op: "exit", reason: "left_aoi"})
			}
			w.held = map[string]struct{}{}
			w.lastCells = nil
		} else if w.cellZ != prevZ {
			w.lastCells = nil // the client drops cells of another zoom
			w.cellOps = nil   // queued cells of the old zoom are moot
		}
		w.cellsDirty = true
		w.lastCellsT = time.Time{}
		w.accountLocked()
		return
	}
	if wasCells {
		w.cellOps, w.lastCells = nil, nil
	}
	now := map[string]struct{}{}
	for key, e := range s.perLayer[w.li] {
		if w.visible(e) {
			now[key] = struct{}{}
			if _, held := w.held[key]; !held || wasCells {
				w.queue(key, pendOp{op: "enter", e: e})
			}
		}
	}
	if !wasCells {
		for key := range w.held {
			if _, ok := now[key]; !ok {
				w.queue(key, pendOp{op: "exit", reason: "left_aoi"})
			}
		}
	}
	w.held = now
	w.accountLocked()
}

// computeCellsLocked queues changed cells for a cell-mode watch. Store.mu (read) and w.mu held.
func (w *Watch) computeCellsLocked(s *Store) {
	cur := s.sharedCellsLocked(w)
	var changed []cell
	for k, c := range cur {
		if p, ok := w.lastCells[k]; !ok || p != c {
			changed = append(changed, c)
		}
	}
	for k, p := range w.lastCells {
		if _, ok := cur[k]; !ok {
			changed = append(changed, cell{p[0], p[1], 0, p[3], p[4]})
		}
	}
	w.lastCells = cur
	for i := 0; i < len(changed); i += maxCellsPerOp {
		j := min(i+maxCellsPerOp, len(changed))
		b, _ := json.Marshal(struct {
			Op    string `json:"op"`
			Z     int    `json:"z"`
			Cells []cell `json:"cells"`
		}{"cells", w.cellZ, changed[i:j]})
		w.cellOps = append(w.cellOps, b)
	}
	w.cellsDirty = false
	w.accountLocked()
}

// batch is one built delta batch.
type batch struct {
	reply string
	data  []byte
	ops   int
}

// buildBatch takes at most one batch of pending ops (Store.mu read held by the caller).
func (w *Watch) buildBatchLocked(s *Store, now time.Time, heartbeat, cellEvery time.Duration) (*batch, bool) {
	if w.closed {
		return nil, false
	}
	if w.endCode != "" {
		w.releaseLocked()
		return &batch{reply: w.reply, data: errReply(w.endCode, w.endDetail)}, false
	}
	if now.After(w.leaseUntil) {
		w.releaseLocked()
		return nil, false
	}
	// Cells are diffed against the last computed state, so a new diff is taken only once the
	// previous one has been sent: queued cell ops never pile up behind a slow max_rate_hz.
	if w.cellsMode && w.cellsDirty && len(w.cellOps) == 0 && now.Sub(w.lastCellsT) >= cellEvery {
		w.computeCellsLocked(s)
		w.lastCellsT = now
	}
	if len(w.order) == 0 && len(w.cellOps) == 0 && w.meta == nil && now.Sub(w.lastSent) < heartbeat {
		return nil, true
	}
	to := s.lastSeq
	from := w.nextFrom
	if to >= w.nextFrom {
		w.nextFrom = to + 1
	} else {
		from = to
	}
	w.n++
	mode := "entities"
	if w.cellsMode {
		mode = "cells"
	}
	var head bytes.Buffer
	head.WriteString(`{"v":1,"watch_id":`)
	head.WriteString(strconv.Quote(w.id))
	head.WriteString(`,"layer":`)
	head.WriteString(strconv.Quote(w.def.ID))
	head.WriteString(`,"n":` + strconv.FormatUint(w.n, 10))
	head.WriteString(`,"from_seq":` + strconv.FormatUint(from, 10))
	head.WriteString(`,"to_seq":` + strconv.FormatUint(to, 10))
	head.WriteString(`,"as_of":` + strconv.FormatInt(now.UnixMilli(), 10))
	head.WriteString(`,"mode":"` + mode + `","ops":[`)
	tail := []byte(`],"meta":null}`)
	if w.meta != nil {
		tail = append(append([]byte(`],"meta":`), w.meta...), '}')
		w.meta = nil
	}
	budget := MaxBatchBytes - head.Len() - len(tail)
	body := make([]byte, 0, 4096)
	ops := 0
	add := func(b []byte) bool {
		extra := len(b)
		if ops > 0 {
			extra++
		}
		if ops >= MaxBatchOps || (ops > 0 && len(body)+extra > budget) {
			return false
		}
		if ops > 0 {
			body = append(body, ',')
		}
		body = append(body, b...)
		ops++
		return true
	}
	ci := 0
	for ci < len(w.cellOps) && add(w.cellOps[ci]) {
		ci++
	}
	w.cellOps = w.cellOps[ci:]
	oi := 0
	for oi < len(w.order) {
		key := w.order[oi]
		if !add(w.render(key, w.pend[key])) {
			break
		}
		delete(w.pend, key)
		oi++
	}
	w.order = w.order[oi:]
	if len(w.order) == 0 {
		w.order = nil
	} else if cap(w.order) > 4*len(w.order)+1024 {
		w.order = append([]string(nil), w.order...)
	}
	w.accountLocked()
	out := make([]byte, 0, head.Len()+len(body)+len(tail))
	out = append(out, head.Bytes()...)
	out = append(out, body...)
	out = append(out, tail...)
	w.lastSent = now
	return &batch{reply: w.reply, data: out, ops: ops}, true
}

// pendingLen reports queued ops (for tests and metrics).
func (w *Watch) pendingLen() int {
	w.mu.Lock()
	defer w.mu.Unlock()
	return len(w.order) + len(w.cellOps)
}

func newWatchID() string { return "w-" + nuid.Next() }
