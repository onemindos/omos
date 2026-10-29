// Package layerclient is a reference client for the layer contracts (onemind-xr
// docs/contracts/layer-snapshot.md §3.1 and layer-delta.md §2): one `svc.layer.snapshot`
// read, then `svc.layer.watch` from `ent_seq + 1`, applying ops with the client rules (seq
// order, gap detection on `n`, dead-watch timeout, re-snapshot on `resnapshot`). The service
// tests use it; the Unreal and L0 clients implement the same rules.
package layerclient

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/nats-io/nats.go"
)

// Entity is the render entity (layer-snapshot.md §4) as the client holds it.
type Entity struct {
	EntityID   string   `json:"entity_id"`
	EntityType string   `json:"entity_type"`
	Lat        float64  `json:"lat"`
	Lon        float64  `json:"lon"`
	HAE        *float64 `json:"hae"`
	Seq        uint64   `json:"seq"`
	T          int64    `json:"t"`
	ExpiryT    *int64   `json:"expiry_t"`
	Live       bool     `json:"live"`
	Simulated  bool     `json:"simulated"`
	Origin     string   `json:"origin"`
}

// Chunk is one snapshot reply message.
type Chunk struct {
	V         int               `json:"v"`
	Layer     string            `json:"layer"`
	Part      int               `json:"part"`
	Last      bool              `json:"last"`
	Mode      string            `json:"mode"`
	EntSeq    uint64            `json:"ent_seq"`
	AsOf      int64             `json:"as_of"`
	Truncated bool              `json:"truncated"`
	Items     []json.RawMessage `json:"items"`
	Error     string            `json:"error"`
	Detail    string            `json:"detail"`
	Bytes     int               `json:"-"`
}

// Snapshot is a whole snapshot read.
type Snapshot struct {
	Chunks    []Chunk
	Items     []json.RawMessage
	Mode      string
	EntSeq    uint64
	AsOf      int64
	Truncated bool
}

// ErrService is an error reply from the layer service.
type ErrService struct{ Code, Detail string }

func (e *ErrService) Error() string { return e.Code + ": " + e.Detail }

// IsCode reports whether err is a service error with code.
func IsCode(err error, code string) bool {
	var es *ErrService
	return errors.As(err, &es) && es.Code == code
}

// GetSnapshot sends one svc.layer.snapshot request and collects chunks until `last`.
func GetSnapshot(nc *nats.Conn, req map[string]any, timeout time.Duration) (*Snapshot, error) {
	body, _ := json.Marshal(req)
	inbox := nc.NewInbox()
	sub, err := nc.SubscribeSync(inbox)
	if err != nil {
		return nil, err
	}
	defer sub.Unsubscribe()
	if err := nc.PublishRequest("svc.layer.snapshot", inbox, body); err != nil {
		return nil, err
	}
	out := &Snapshot{}
	for {
		m, err := sub.NextMsg(timeout)
		if err != nil {
			return nil, fmt.Errorf("snapshot: %w", err)
		}
		var c Chunk
		if err := json.Unmarshal(m.Data, &c); err != nil {
			return nil, err
		}
		c.Bytes = len(m.Data)
		if c.Error != "" {
			return nil, &ErrService{c.Error, c.Detail}
		}
		if c.Part != len(out.Chunks) {
			return nil, fmt.Errorf("snapshot chunk %d arrived as number %d", c.Part, len(out.Chunks))
		}
		out.Chunks = append(out.Chunks, c)
		out.Items = append(out.Items, c.Items...)
		out.Mode, out.EntSeq, out.AsOf = c.Mode, c.EntSeq, c.AsOf
		out.Truncated = out.Truncated || c.Truncated
		if c.Last {
			return out, nil
		}
	}
}

// Op is one delta op.
type Op struct {
	Op     string          `json:"op"`
	E      json.RawMessage `json:"e"`
	ID     string          `json:"id"`
	Reason string          `json:"reason"`
	Z      int             `json:"z"`
	Cells  [][]float64     `json:"cells"`
}

// Batch is one delta batch (or the watch ack / an error).
type Batch struct {
	V       int             `json:"v"`
	WatchID string          `json:"watch_id"`
	Layer   string          `json:"layer"`
	N       uint64          `json:"n"`
	FromSeq uint64          `json:"from_seq"`
	ToSeq   uint64          `json:"to_seq"`
	AsOf    int64           `json:"as_of"`
	Mode    string          `json:"mode"`
	Ops     []Op            `json:"ops"`
	Meta    json.RawMessage `json:"meta"`
	LeaseS  int             `json:"lease_s"`
	EntSeq  uint64          `json:"ent_seq"`
	Error   string          `json:"error"`
	Detail  string          `json:"detail"`
	Bytes   int             `json:"-"`
	At      time.Time       `json:"-"`
}

// Follower keeps one layer view current: snapshot, then watch, re-snapshotting on any gap.
type Follower struct {
	NC        *nats.Conn
	Req       map[string]any // layer, aoi, zoom, fields (+ lease_s, max_rate_hz for the watch)
	DeadAfter time.Duration  // no batch for this long: re-snapshot (contract: 30 s)
	Timeout   time.Duration  // per-message wait
	// DropBatch, when set, makes the client skip a batch (to test gap recovery).
	DropBatch func(b *Batch) bool

	mu          sync.Mutex
	ents        map[string]Entity
	cells       map[[3]int][]float64
	watchID     string
	n           uint64
	batches     []Batch
	resnapshots int
	reasons     []string
	exits       map[string]int
	enters      map[string]int
	removes     map[string]string
	inbox       string
	snapshotSeq uint64
	lastMeta    json.RawMessage
}

// Entities returns a copy of the held entities by id.
func (f *Follower) Entities() map[string]Entity {
	f.mu.Lock()
	defer f.mu.Unlock()
	out := make(map[string]Entity, len(f.ents))
	for k, v := range f.ents {
		out[k] = v
	}
	return out
}

// Cells returns a copy of the held cells keyed by z, x, y.
func (f *Follower) Cells() map[[3]int][]float64 {
	f.mu.Lock()
	defer f.mu.Unlock()
	out := make(map[[3]int][]float64, len(f.cells))
	for k, v := range f.cells {
		out[k] = v
	}
	return out
}

// Watching reports whether a watch is acknowledged and running.
func (f *Follower) Watching() bool {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.watchID != ""
}

// Batches returns the batches received so far.
func (f *Follower) Batches() []Batch {
	f.mu.Lock()
	defer f.mu.Unlock()
	return append([]Batch(nil), f.batches...)
}

// Resnapshots returns how many times the client re-snapshotted, with the reasons.
func (f *Follower) Resnapshots() (int, []string) {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.resnapshots, append([]string(nil), f.reasons...)
}

// Counts returns how many enter, exit and remove ops arrived per id.
func (f *Follower) Counts() (enters, exits map[string]int, removes map[string]string) {
	f.mu.Lock()
	defer f.mu.Unlock()
	cp := func(m map[string]int) map[string]int {
		o := map[string]int{}
		for k, v := range m {
			o[k] = v
		}
		return o
	}
	r := map[string]string{}
	for k, v := range f.removes {
		r[k] = v
	}
	return cp(f.enters), cp(f.exits), r
}

// SetAOI changes the AOI and zoom of the running watch (renew with the same watch_id).
func (f *Follower) SetAOI(aoi any, zoom int) error {
	f.mu.Lock()
	f.Req["aoi"], f.Req["zoom"] = aoi, zoom
	id, inbox := f.watchID, f.inbox
	f.mu.Unlock()
	if id == "" {
		return errors.New("no watch running yet")
	}
	return f.renew(id, inbox, true)
}

func (f *Follower) req() map[string]any {
	f.mu.Lock()
	defer f.mu.Unlock()
	m := make(map[string]any, len(f.Req))
	for k, v := range f.Req {
		m[k] = v
	}
	return m
}

func (f *Follower) watchReq(extra map[string]any) []byte {
	m := f.req()
	m["v"] = 1
	for k, v := range extra {
		m[k] = v
	}
	b, _ := json.Marshal(m)
	return b
}

func (f *Follower) renew(id, inbox string, withAOI bool) error {
	f.mu.Lock()
	m := map[string]any{"v": 1, "watch_id": id}
	for k, v := range f.Req {
		m[k] = v
	}
	if !withAOI {
		delete(m, "aoi")
		delete(m, "zoom")
	}
	f.mu.Unlock()
	body, _ := json.Marshal(m)
	// The renew reuses the watch inbox as its reply subject, so batches keep arriving there
	// and a cross-account response mapping is refreshed (layer-delta.md §2.1).
	return f.NC.PublishRequest("svc.layer.watch", inbox, body)
}

// Run follows the layer until ctx ends.
func (f *Follower) Run(ctx context.Context) error {
	if f.DeadAfter == 0 {
		f.DeadAfter = 30 * time.Second
	}
	if f.Timeout == 0 {
		f.Timeout = 5 * time.Second
	}
	f.mu.Lock()
	f.ents, f.cells = map[string]Entity{}, map[[3]int][]float64{}
	f.exits, f.enters, f.removes = map[string]int{}, map[string]int{}, map[string]string{}
	f.mu.Unlock()
	first := true
	for ctx.Err() == nil {
		reason, err := f.session(ctx, first)
		first = false
		if ctx.Err() != nil {
			return nil
		}
		f.mu.Lock()
		f.resnapshots++
		f.reasons = append(f.reasons, reason)
		f.mu.Unlock()
		if err != nil || strings.HasPrefix(reason, "watch error") {
			// Back off after a failure or an error reply (too_many, resnapshot, layer_down)
			// instead of retrying in a tight loop.
			select {
			case <-ctx.Done():
				return nil
			case <-time.After(200 * time.Millisecond):
			}
		}
	}
	return nil
}

// session: snapshot, watch, apply until a gap; returns why it ended.
func (f *Follower) session(ctx context.Context, first bool) (string, error) {
	snap, err := GetSnapshot(f.NC, f.req(), f.Timeout)
	if err != nil {
		return "snapshot failed: " + err.Error(), err
	}
	f.mu.Lock()
	f.watchID = ""
	f.ents, f.cells = map[string]Entity{}, map[[3]int][]float64{}
	for _, it := range snap.Items {
		if snap.Mode == "cells" {
			var c []float64
			_ = json.Unmarshal(it, &c)
			if len(c) == 5 {
				f.cells[[3]int{-1, int(c[0]), int(c[1])}] = c
			}
			continue
		}
		var e Entity
		if json.Unmarshal(it, &e) == nil {
			f.ents[e.EntityID] = e
		}
	}
	f.snapshotSeq = snap.EntSeq
	f.inbox = f.NC.NewInbox()
	inbox := f.inbox
	f.mu.Unlock()

	sub, err := f.NC.SubscribeSync(inbox)
	if err != nil {
		return "subscribe failed", err
	}
	defer sub.Unsubscribe()
	if err := f.NC.PublishRequest("svc.layer.watch", inbox, f.watchReq(map[string]any{"from_seq": snap.EntSeq + 1})); err != nil {
		return "watch failed", err
	}
	m, err := sub.NextMsg(f.Timeout)
	if err != nil {
		return "watch ack timeout", err
	}
	var ack Batch
	if err := json.Unmarshal(m.Data, &ack); err != nil {
		return "watch ack unreadable", err
	}
	if ack.Error != "" {
		return "watch error " + ack.Error, nil
	}
	lease := time.Duration(ack.LeaseS) * time.Second
	f.mu.Lock()
	f.watchID, f.n = ack.WatchID, 0
	f.mu.Unlock()
	renewAt := time.Now().Add(lease / 2)
	last := time.Now()
	for ctx.Err() == nil {
		if time.Now().After(renewAt) {
			_ = f.renew(ack.WatchID, inbox, false)
			renewAt = time.Now().Add(lease / 2)
		}
		m, err := sub.NextMsg(100 * time.Millisecond)
		if errors.Is(err, nats.ErrTimeout) {
			if time.Since(last) > f.DeadAfter {
				return "dead watch: no batch for " + f.DeadAfter.String(), nil
			}
			continue
		}
		if errors.Is(err, nats.ErrNoResponders) {
			// A renew sent with the watch inbox as reply found no layer service.
			return "dead watch: no responders to the renew", nil
		}
		if err != nil {
			return "inbox error: " + err.Error(), err
		}
		var b Batch
		if err := json.Unmarshal(m.Data, &b); err != nil {
			return "batch unreadable", err
		}
		b.Bytes, b.At = len(m.Data), time.Now()
		if b.Error != "" {
			return "watch error " + b.Error, nil
		}
		if b.N == 0 { // renew ack
			continue
		}
		last = time.Now()
		if f.DropBatch != nil && f.DropBatch(&b) {
			continue
		}
		f.mu.Lock()
		if b.N != f.n+1 {
			f.mu.Unlock()
			return fmt.Sprintf("gap: batch %d after %d", b.N, f.n), nil
		}
		f.n = b.N
		f.batches = append(f.batches, b)
		f.applyLocked(&b)
		f.mu.Unlock()
	}
	return "stopped", nil
}

func (f *Follower) applyLocked(b *Batch) {
	if b.Meta != nil && string(b.Meta) != "null" {
		f.lastMeta = b.Meta
	}
	for _, op := range b.Ops {
		switch op.Op {
		case "upsert", "enter":
			var e Entity
			if json.Unmarshal(op.E, &e) != nil {
				continue
			}
			if op.Op == "enter" {
				f.enters[e.EntityID]++
			}
			if old, ok := f.ents[e.EntityID]; ok && op.Op == "upsert" && e.Seq < old.Seq {
				continue // entity revision rule: never go back
			}
			f.ents[e.EntityID] = e
		case "exit":
			f.exits[op.ID]++
			delete(f.ents, op.ID)
		case "remove":
			f.removes[op.ID] = op.Reason
			delete(f.ents, op.ID)
		case "cells":
			for k := range f.cells {
				if k[0] != op.Z {
					delete(f.cells, k)
				}
			}
			for _, c := range op.Cells {
				if len(c) != 5 {
					continue
				}
				k := [3]int{op.Z, int(c[0]), int(c[1])}
				if c[2] == 0 {
					delete(f.cells, k)
				} else {
					f.cells[k] = c
				}
			}
		}
	}
}
