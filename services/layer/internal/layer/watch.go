package layer

import (
	"bytes"
	"encoding/json"
	"strconv"
	"sync"
	"time"

	"github.com/nats-io/nuid"
)

// Watch limits (layer-delta.md §2).
const (
	MaxBatchBytes = 256 * 1024
	MaxBatchOps   = 1000
	maxCellsPerOp = 2000
)

// Watch is one `svc.layer.watch` subscription.
type Watch struct {
	id     string
	li     int
	def    *Def
	caller Caller

	mu         sync.Mutex
	reply      string
	area       *area
	zoom       int
	cellsMode  bool
	cellZ      int
	full       bool
	rateHz     int
	leaseUntil time.Time
	held       map[string]struct{}
	order      []string
	pend       map[string][]byte
	cellOps    [][]byte
	cellsDirty bool
	lastCells  map[[2]int]cell
	lastCellsT time.Time
	meta       json.RawMessage
	n          uint64
	nextFrom   uint64
	lastSent   time.Time
	closed     bool
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
func (w *Watch) queue(key string, data []byte) {
	if _, ok := w.pend[key]; !ok {
		w.order = append(w.order, key)
	}
	w.pend[key] = data
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
	if w.closed {
		return
	}
	if w.cellsMode {
		w.cellsDirty = true
		return
	}
	_, held := w.held[key]
	switch {
	case nw != nil && w.visible(nw):
		op := "upsert"
		if !held && old != nil {
			op = "enter"
		}
		w.queue(key, opEntity(op, w.entityJSON(nw)))
		w.held[key] = struct{}{}
	case held:
		switch {
		case nw == nil:
			w.queue(key, opID("remove", idOf(key), reason))
		case member(w.def, nw) && !w.caller.CanSee(nw.tlp, nw.releasable):
			w.queue(key, opID("remove", idOf(key), "released_out"))
		default:
			w.queue(key, opID("exit", idOf(key), "left_aoi"))
		}
		delete(w.held, key)
	}
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
				w.queue(key, opID("exit", idOf(key), "left_aoi"))
			}
			w.held = map[string]struct{}{}
			w.lastCells = nil
		} else if w.cellZ != prevZ {
			w.lastCells = nil // the client drops cells of another zoom
		}
		w.cellsDirty = true
		w.lastCellsT = time.Time{}
		return
	}
	now := map[string]struct{}{}
	for key, e := range s.perLayer[w.li] {
		if w.visible(e) {
			now[key] = struct{}{}
			if _, held := w.held[key]; !held || wasCells {
				w.queue(key, opEntity("enter", w.entityJSON(e)))
			}
		}
	}
	if !wasCells {
		for key := range w.held {
			if _, ok := now[key]; !ok {
				w.queue(key, opID("exit", idOf(key), "left_aoi"))
			}
		}
	}
	w.held = now
}

// computeCellsLocked queues changed cells for a cell-mode watch. Store.mu (read) and w.mu held.
func (w *Watch) computeCellsLocked(s *Store) {
	g := newCellGrid(w.cellZ)
	for _, e := range s.perLayer[w.li] {
		if w.visible(e) {
			g.add(e.r.Lat, e.r.Lon)
		}
	}
	cur := g.asMap()
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
	if now.After(w.leaseUntil) {
		w.closed = true
		return nil, false
	}
	if w.cellsMode && w.cellsDirty && now.Sub(w.lastCellsT) >= cellEvery {
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
		if !add(w.pend[key]) {
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
