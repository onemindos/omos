package layer

import (
	"bytes"
	"encoding/json"
	"strconv"
	"time"
)

// Snapshot limits (layer-snapshot.md §3.1).
const (
	MaxChunkBytes = 512 * 1024
	MaxChunkItems = 2000
	fullMaxItems  = 200
	defaultZoom   = 22
)

// Error codes (layer-snapshot.md §3.1, layer-delta.md §2.1).
const (
	ErrUnknownLayer = "unknown_layer"
	ErrAOIRequired  = "aoi_required"
	ErrTooMany      = "too_many"
	ErrLayerDown    = "layer_down"
	ErrForbidden    = "forbidden"
	ErrResnapshot   = "resnapshot"
	ErrBadRequest   = "bad_request"
)

// errReply is an error answer; `last: true` ends a request-many read.
func errReply(code, detail string) []byte {
	b, _ := json.Marshal(struct {
		V      int    `json:"v"`
		Error  string `json:"error"`
		Detail string `json:"detail,omitempty"`
		Last   bool   `json:"last"`
	}{1, code, detail, true})
	return b
}

type snapReq struct {
	V      int     `json:"v"`
	Layer  string  `json:"layer"`
	AOI    *AOIReq `json:"aoi"`
	Zoom   *int    `json:"zoom"`
	Fields string  `json:"fields"`
	Max    int     `json:"max"`
}

// resolveAreaLocked turns an AOI request into an area (Store.mu held).
func (s *Store) resolveAreaLocked(a *AOIReq) (*area, string) {
	if a == nil {
		return nil, ""
	}
	switch {
	case len(a.BBox) > 0:
		ar, err := bboxArea(a.BBox)
		if err != nil {
			return nil, err.Error()
		}
		return ar, ""
	case len(a.Polygon) > 0:
		ar, err := polyArea(a.Polygon)
		if err != nil {
			return nil, err.Error()
		}
		return ar, ""
	case a.ZoneID != "":
		if !tokenRe.MatchString(a.ZoneID) {
			return nil, "zone_id is not a valid entity id"
		}
		z := s.ents["zone."+a.ZoneID]
		if z == nil {
			return nil, "zone " + a.ZoneID + " is not in ENT"
		}
		ar, err := zoneArea(z.geo)
		if err != nil {
			return nil, err.Error()
		}
		return ar, ""
	}
	return nil, errBadAOI.Error()
}

// isEmptyAOI reports an aoi of null or {}.
func isEmptyAOI(a *AOIReq) bool {
	return a == nil || (len(a.BBox) == 0 && len(a.Polygon) == 0 && a.ZoneID == "")
}

// SnapshotResult is a snapshot ready to send.
type SnapshotResult struct {
	Chunks [][]byte
	Items  int
	Mode   string
	EntSeq uint64
}

// Snapshot answers one `svc.layer.snapshot` request (layer-snapshot.md §3.1). zooms are the
// manifest's current cell zooms for the layer (nil: the definition's).
func (s *Store) Snapshot(raw []byte, caller Caller, zoomsFor func(li int) []int, object func(li int) *ObjectRef) (*SnapshotResult, []byte) {
	var req snapReq
	if err := json.Unmarshal(raw, &req); err != nil {
		return nil, errReply(ErrBadRequest, "request is not JSON")
	}
	if !caller.Known {
		return nil, errReply(ErrForbidden, "identity not in the layer service directory")
	}
	s.mu.RLock()
	li := s.layerIndex(req.Layer)
	if li < 0 {
		s.mu.RUnlock()
		return nil, errReply(ErrUnknownLayer, req.Layer)
	}
	d := s.defs[li]
	if !s.ready || s.down {
		s.mu.RUnlock()
		return nil, errReply(ErrLayerDown, "ENT not readable")
	}
	asOf := s.now().UnixMilli()
	if d.Tier == TierStatic {
		s.mu.RUnlock()
		obj := object(li)
		if obj == nil {
			return nil, errReply(ErrLayerDown, "no object published for static layer")
		}
		item, _ := json.Marshal(obj)
		return &SnapshotResult{Chunks: [][]byte{chunk(d.ID, 0, true, "object", 0, asOf, false, [][]byte{item})}, Items: 1, Mode: "object"}, nil
	}
	if d.Tier == TierLarge && isEmptyAOI(req.AOI) {
		s.mu.RUnlock()
		return nil, errReply(ErrAOIRequired, "large layer")
	}
	full := req.Fields == "full"
	max := req.Max
	if max <= 0 || max > d.Cap {
		max = d.Cap
	}
	if full && d.Tier != TierSmall && max > fullMaxItems {
		s.mu.RUnlock()
		return nil, errReply(ErrTooMany, "fields full needs a small layer or max <= 200")
	}
	var ar *area
	if !isEmptyAOI(req.AOI) {
		var msg string
		if ar, msg = s.resolveAreaLocked(req.AOI); msg != "" {
			s.mu.RUnlock()
			return nil, errReply(ErrBadRequest, msg)
		}
	}
	zoom := defaultZoom
	if req.Zoom != nil {
		zoom = *req.Zoom
	}
	seq := s.lastSeq
	if d.Cells != nil && zoom < d.Cells.BelowZoom {
		zooms := zoomsFor(li)
		if len(zooms) == 0 {
			zooms = d.Cells.Zooms
		}
		g := newCellGrid(d.CellZoom(zoom, zooms))
		for _, e := range s.perLayer[li] {
			if ar.contains(e.r.Lat, e.r.Lon) && caller.CanSee(e.tlp, e.releasable) {
				g.add(e.r.Lat, e.r.Lon)
			}
		}
		s.mu.RUnlock()
		cells := g.list()
		truncated := len(cells) > max
		if truncated {
			cells = cells[:max]
		}
		items := make([][]byte, len(cells))
		for i, c := range cells {
			items[i], _ = json.Marshal(c)
		}
		ch := chunkAll(d.ID, "cells", seq, asOf, truncated, items)
		return &SnapshotResult{Chunks: ch, Items: len(items), Mode: "cells", EntSeq: seq}, nil
	}
	items := make([][]byte, 0, min(len(s.perLayer[li]), max))
	truncated := false
	for _, e := range s.perLayer[li] {
		if !ar.contains(e.r.Lat, e.r.Lon) || !caller.CanSee(e.tlp, e.releasable) {
			continue
		}
		if len(items) >= max {
			truncated = true
			break
		}
		if full {
			items = append(items, e.raw)
		} else {
			items = append(items, e.renderJSON)
		}
	}
	s.mu.RUnlock()
	ch := chunkAll(d.ID, "entities", seq, asOf, truncated, items)
	return &SnapshotResult{Chunks: ch, Items: len(items), Mode: "entities", EntSeq: seq}, nil
}

func chunk(layer string, part int, last bool, mode string, seq uint64, asOf int64, truncated bool, items [][]byte) []byte {
	var b bytes.Buffer
	b.WriteString(`{"v":1,"layer":` + strconv.Quote(layer) + `,"part":` + strconv.Itoa(part) + `,"last":` + strconv.FormatBool(last))
	b.WriteString(`,"mode":"` + mode + `","ent_seq":` + strconv.FormatUint(seq, 10) + `,"as_of":` + strconv.FormatInt(asOf, 10))
	b.WriteString(`,"truncated":` + strconv.FormatBool(truncated) + `,"items":[`)
	for i, it := range items {
		if i > 0 {
			b.WriteByte(',')
		}
		b.Write(it)
	}
	b.WriteString(`]}`)
	return b.Bytes()
}

// chunkAll splits items into chunks of at most MaxChunkBytes and MaxChunkItems.
func chunkAll(layer, mode string, seq uint64, asOf int64, truncated bool, items [][]byte) [][]byte {
	const envelope = 256 // generous upper bound of the chunk envelope without items
	var groups [][][]byte
	cur, size := [][]byte{}, 0
	for _, it := range items {
		if len(cur) > 0 && (len(cur) >= MaxChunkItems || size+len(it)+1+envelope > MaxChunkBytes) {
			groups = append(groups, cur)
			cur, size = [][]byte{}, 0
		}
		cur = append(cur, it)
		size += len(it) + 1
	}
	groups = append(groups, cur)
	out := make([][]byte, len(groups))
	for i, g := range groups {
		out[i] = chunk(layer, i, i == len(groups)-1, mode, seq, asOf, truncated, g)
	}
	return out
}

// ---------------------------------------------------------------- watch

type watchReq struct {
	V         int     `json:"v"`
	Layer     string  `json:"layer"`
	AOI       *AOIReq `json:"aoi"`
	Zoom      *int    `json:"zoom"`
	Fields    string  `json:"fields"`
	FromSeq   uint64  `json:"from_seq"`
	LeaseS    int     `json:"lease_s"`
	MaxRateHz int     `json:"max_rate_hz"`
	WatchID   string  `json:"watch_id"`
}

// watchAck is the first reply of a watch and the reply to a renew.
type watchAck struct {
	V       int    `json:"v"`
	WatchID string `json:"watch_id"`
	LeaseS  int    `json:"lease_s"`
	EntSeq  uint64 `json:"ent_seq"`
}

// WatchOpts are the limits of the watch path.
type WatchOpts struct {
	MaxWatches int
	ZoomsFor   func(li int) []int
}

func clamp(v, lo, hi, def int) int {
	if v == 0 {
		return def
	}
	return max(lo, min(hi, v))
}

// Watch starts or renews a watch (layer-delta.md §2.1). It returns the watch to run (nil on
// a renew or an error) and the reply to send to msg.Reply.
func (s *Store) Watch(raw []byte, reply string, caller Caller, o WatchOpts) (*Watch, []byte) {
	var req watchReq
	if err := json.Unmarshal(raw, &req); err != nil {
		return nil, errReply(ErrBadRequest, "request is not JSON")
	}
	if !caller.Known {
		return nil, errReply(ErrForbidden, "identity not in the layer service directory")
	}
	if reply == "" {
		return nil, nil
	}
	leaseS := clamp(req.LeaseS, 10, 300, 60)
	rate := clamp(req.MaxRateHz, 1, 10, 4)
	s.mu.RLock()
	defer s.mu.RUnlock()
	if !s.ready || s.down {
		return nil, errReply(ErrLayerDown, "ENT not readable")
	}
	if req.WatchID != "" {
		return nil, s.renewLocked(req, reply, caller, leaseS, rate, o)
	}
	li := s.layerIndex(req.Layer)
	if li < 0 {
		return nil, errReply(ErrUnknownLayer, req.Layer)
	}
	d := s.defs[li]
	if d.Tier == TierStatic {
		return nil, errReply(ErrBadRequest, "static layers change only through the manifest (watch the meta)")
	}
	if d.Tier == TierLarge && isEmptyAOI(req.AOI) {
		return nil, errReply(ErrAOIRequired, "large layer")
	}
	full := req.Fields == "full"
	if full && d.Tier != TierSmall {
		return nil, errReply(ErrTooMany, "fields full on the watch path needs a small layer")
	}
	var ar *area
	if !isEmptyAOI(req.AOI) {
		var msg string
		if ar, msg = s.resolveAreaLocked(req.AOI); msg != "" {
			return nil, errReply(ErrBadRequest, msg)
		}
	}
	if req.FromSeq != 0 && !s.replayValidLocked(req.FromSeq) {
		return nil, errReply(ErrResnapshot, "from_seq is outside the replay window")
	}
	if s.WatchCountLockedOK(o.MaxWatches) {
		return nil, errReply(ErrTooMany, "watch limit reached")
	}
	now := s.now()
	w := &Watch{id: newWatchID(), li: li, def: d, caller: caller, reply: reply, area: ar, full: full,
		rateHz: rate, leaseUntil: now.Add(time.Duration(leaseS) * time.Second), held: map[string]struct{}{},
		pend: map[string][]byte{}, lastSent: now, wake: make(chan struct{}, 1)}
	w.setZoomLocked(req.Zoom, o.ZoomsFor(li))
	if req.FromSeq == 0 {
		w.nextFrom = s.lastSeq + 1
		w.resetViewLocked(s, true, -1) // no snapshot: the whole current view goes out as enter ops
	} else {
		w.nextFrom = req.FromSeq
		w.replayLocked(s, req.FromSeq)
	}
	s.wmu.Lock()
	s.watches[w.id] = w
	s.wmu.Unlock()
	ack, _ := json.Marshal(watchAck{V: 1, WatchID: w.id, LeaseS: leaseS, EntSeq: s.lastSeq})
	return w, ack
}

// WatchCountLockedOK reports whether the watch limit is reached.
func (s *Store) WatchCountLockedOK(limit int) bool {
	s.wmu.RLock()
	defer s.wmu.RUnlock()
	return limit > 0 && len(s.watches) >= limit
}

func (w *Watch) setZoomLocked(z *int, zooms []int) {
	w.zoom = defaultZoom
	if z != nil {
		w.zoom = *z
	}
	w.cellsMode = false
	if w.def.Cells != nil && w.zoom < w.def.Cells.BelowZoom {
		if len(zooms) == 0 {
			zooms = w.def.Cells.Zooms
		}
		w.cellsMode = true
		w.cellZ = w.def.CellZoom(w.zoom, zooms)
	}
}

// replayLocked brings a client from fromSeq-1 to now: for each entity changed since then
// (sweep removals from fromSeq-1, because they carry the sequence seen at sweep time), send
// the current state (upsert), or remove/exit when it is gone or no longer visible. Then the
// held set is the current visible set. Store.mu is held; w is not yet shared.
func (w *Watch) replayLocked(s *Store, fromSeq uint64) {
	if w.cellsMode {
		w.cellsDirty = true
		return
	}
	reasons := map[string]string{}
	var keys []string
	for _, c := range s.ring {
		if c.seq < fromSeq && !(c.sweep && c.seq+1 >= fromSeq) {
			continue
		}
		if _, seen := reasons[c.key]; !seen {
			keys = append(keys, c.key)
		}
		reasons[c.key] = c.reason
	}
	for _, key := range keys {
		e := s.ents[key]
		switch {
		case e != nil && w.visible(e):
			w.queue(key, opEntity("upsert", w.entityJSON(e)))
		case e == nil:
			r := reasons[key]
			if r == "" {
				r = "deleted"
			}
			w.queue(key, opID("remove", idOf(key), r))
		case member(w.def, e) && !w.caller.CanSee(e.tlp, e.releasable):
			w.queue(key, opID("remove", idOf(key), "released_out"))
		default:
			w.queue(key, opID("exit", idOf(key), "left_aoi"))
		}
	}
	for key, e := range s.perLayer[w.li] {
		if w.visible(e) {
			w.held[key] = struct{}{}
		}
	}
}

func (s *Store) renewLocked(req watchReq, reply string, caller Caller, leaseS, rate int, o WatchOpts) []byte {
	s.wmu.RLock()
	w := s.watches[req.WatchID]
	s.wmu.RUnlock()
	if w == nil {
		return errReply(ErrResnapshot, "unknown or ended watch")
	}
	if w.caller.ID != caller.ID {
		return errReply(ErrForbidden, "watch belongs to another caller")
	}
	var ar *area
	changeAOI := req.AOI != nil
	if changeAOI && !isEmptyAOI(req.AOI) {
		var msg string
		if ar, msg = s.resolveAreaLocked(req.AOI); msg != "" {
			return errReply(ErrBadRequest, msg)
		}
	}
	if changeAOI && ar == nil && w.def.Tier == TierLarge {
		return errReply(ErrAOIRequired, "large layer")
	}
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.closed {
		return errReply(ErrResnapshot, "watch ended")
	}
	w.leaseUntil = s.now().Add(time.Duration(leaseS) * time.Second)
	w.rateHz = rate
	w.reply = reply // a renew's reply subject keeps a cross-account response mapping fresh
	zoomChanged := req.Zoom != nil && *req.Zoom != w.zoom
	if changeAOI || zoomChanged {
		wasCells, prevZ := w.cellsMode, w.cellZ
		if changeAOI {
			w.area = ar
		}
		z := &w.zoom
		if req.Zoom != nil {
			z = req.Zoom
		}
		w.setZoomLocked(z, o.ZoomsFor(w.li))
		w.resetViewLocked(s, wasCells, prevZ)
		w.signal()
	}
	ack, _ := json.Marshal(watchAck{V: 1, WatchID: w.id, LeaseS: leaseS, EntSeq: s.lastSeq})
	return ack
}

// Unwatch ends a watch.
func (s *Store) Unwatch(raw []byte, caller Caller) []byte {
	var req struct {
		WatchID string `json:"watch_id"`
	}
	if json.Unmarshal(raw, &req) != nil || req.WatchID == "" {
		return errReply(ErrBadRequest, "watch_id required")
	}
	s.wmu.Lock()
	w := s.watches[req.WatchID]
	if w != nil && w.caller.ID == caller.ID {
		delete(s.watches, req.WatchID)
	}
	s.wmu.Unlock()
	if w == nil {
		return errReply(ErrResnapshot, "unknown or ended watch")
	}
	if w.caller.ID != caller.ID {
		return errReply(ErrForbidden, "watch belongs to another caller")
	}
	w.mu.Lock()
	w.closed = true
	w.mu.Unlock()
	w.signal()
	return []byte(`{"v":1,"ok":true}`)
}

// removeWatch drops a closed watch.
func (s *Store) removeWatch(id string) {
	s.wmu.Lock()
	delete(s.watches, id)
	s.wmu.Unlock()
}

// SetMeta hands a changed manifest to every watch of layer li.
func (s *Store) SetMeta(li int, meta []byte) {
	s.wmu.RLock()
	defer s.wmu.RUnlock()
	for _, w := range s.watches {
		if w.li == li {
			w.mu.Lock()
			w.meta = meta
			w.mu.Unlock()
			w.signal()
		}
	}
}
