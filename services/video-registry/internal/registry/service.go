package registry

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"sort"
	"sync"
	"sync/atomic"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

// Subjects, all in the locked grammar (BUS-DESIGN sections 3.2 and 3.4):
//
//	evt.video.{entity_id}.online|offline   evt.{domain}.{entity}.{event}; not stored by any
//	                                       stream until owner proposal RA1 binds evt.video.> in EVT
//	svc.registry.video.get|list            svc.{service}.{endpoint}; svc.registry.> is exported
//	                                       to the CLIENT account (BUS-DESIGN 7.1)
const (
	SubjectGet  = "svc.registry.video.get"
	SubjectList = "svc.registry.video.list"
	QueueGroup  = "video-registry"
	DefaultKV   = "OVERWATCH_VIDEO"
)

// EventSubject is evt.video.{entity_id}.{event}.
func EventSubject(entityID, event string) string { return "evt.video." + entityID + "." + event }

// Config of a Service.
type Config struct {
	Bucket       string
	Origin       string // site or hub name written into every record
	PollInterval time.Duration
	StaleAfter   time.Duration // MediaMTX unreachable this long -> every live stream goes offline
	Endpoints    Endpoints
	Media        *MediaConfig
	Logger       *slog.Logger
}

// Event is the payload of evt.video.{entity_id}.{online|offline}.
type Event struct {
	V             int       `json:"v"`
	Event         string    `json:"event"`
	StreamID      string    `json:"stream_id"`
	EntityType    string    `json:"entity_type"`
	EntityID      string    `json:"entity_id"`
	EntitySubject string    `json:"entity_subject"`
	Reason        string    `json:"reason,omitempty"`
	TS            time.Time `json:"ts"`
	Origin        string    `json:"origin"`
}

// GetRequest selects one entity (any form ParseEntityRef accepts) or a stream id.
type GetRequest struct {
	Entity string `json:"entity"`
}

// Reply of svc.registry.video.get / list.
type Reply struct {
	Record  *Record  `json:"record,omitempty"`
	Records []Record `json:"records,omitempty"`
	Error   string   `json:"error,omitempty"`
}

// Stats are exposed on /metrics.
type Stats struct {
	Polls, PollErrors, Puts, Events, Rejected, RPC atomic.Uint64
}

// Service reconciles MediaMTX paths into the registry.
type Service struct {
	cfg    Config
	nc     *nats.Conn
	kv     jetstream.KeyValue
	lister PathLister
	log    *slog.Logger

	mu       sync.Mutex
	state    map[string]Record // key -> last written record
	seq      uint64
	lastOK   time.Time
	polledOK atomic.Bool
	rejected map[string]bool // path names already logged as rejected

	Stats Stats
}

// New binds the service to an existing KV bucket (the NACK CR creates it in the cluster).
func New(ctx context.Context, nc *nats.Conn, lister PathLister, cfg Config) (*Service, error) {
	if cfg.Bucket == "" {
		cfg.Bucket = DefaultKV
	}
	if cfg.PollInterval <= 0 {
		cfg.PollInterval = 2 * time.Second
	}
	if cfg.StaleAfter <= 0 {
		cfg.StaleAfter = 30 * time.Second
	}
	if cfg.Origin == "" {
		cfg.Origin = "hub"
	}
	if cfg.Logger == nil {
		cfg.Logger = slog.Default()
	}
	js, err := jetstream.New(nc)
	if err != nil {
		return nil, err
	}
	kv, err := js.KeyValue(ctx, cfg.Bucket)
	if err != nil {
		return nil, fmt.Errorf("bind KV %s: %w", cfg.Bucket, err)
	}
	s := &Service{cfg: cfg, nc: nc, kv: kv, lister: lister, log: cfg.Logger, state: map[string]Record{}, rejected: map[string]bool{}}
	if err := s.load(ctx); err != nil {
		return nil, err
	}
	return s, nil
}

// load reads the bucket so a restart neither re-announces live streams nor forgets streams
// that went away while the service was down (they are marked offline on the first poll).
func (s *Service) load(ctx context.Context) error {
	lister, err := s.kv.ListKeys(ctx)
	if err != nil {
		return err
	}
	for k := range lister.Keys() {
		e, err := s.kv.Get(ctx, k)
		if err != nil {
			continue
		}
		var r Record
		if json.Unmarshal(e.Value(), &r) != nil || r.V != 1 {
			continue
		}
		s.state[k] = r
		if r.Seq > s.seq {
			s.seq = r.Seq
		}
	}
	return nil
}

// Ready reports whether one poll has succeeded.
func (s *Service) Ready() bool { return s.polledOK.Load() }

// Run polls until ctx ends.
func (s *Service) Run(ctx context.Context) error {
	subs, err := s.serveRPC()
	if err != nil {
		return err
	}
	defer func() {
		for _, sub := range subs {
			_ = sub.Drain()
		}
	}()
	s.PollOnce(ctx)
	t := time.NewTicker(s.cfg.PollInterval)
	defer t.Stop()
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-t.C:
			s.PollOnce(ctx)
		}
	}
}

// PollOnce lists MediaMTX paths and reconciles them.
func (s *Service) PollOnce(ctx context.Context) {
	s.Stats.Polls.Add(1)
	pctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	items, err := s.lister.ListPaths(pctx)
	cancel()
	now := time.Now()
	if err != nil {
		s.Stats.PollErrors.Add(1)
		s.log.Warn("mediamtx poll failed", "err", err)
		s.mu.Lock()
		stale := !s.lastOK.IsZero() && now.Sub(s.lastOK) > s.cfg.StaleAfter
		s.mu.Unlock()
		if stale {
			s.Reconcile(ctx, nil, now, "mediamtx_unreachable")
		}
		return
	}
	s.mu.Lock()
	s.lastOK = now
	s.mu.Unlock()
	s.Reconcile(ctx, items, now, "")
	s.polledOK.Store(true)
}

// Reconcile writes new or changed live streams and marks vanished ones offline. offlineReason
// is used for streams that vanished ("" means the publisher stopped).
func (s *Service) Reconcile(ctx context.Context, items []PathItem, now time.Time, offlineReason string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	seen := map[string]bool{}
	for _, it := range items {
		if !it.Ready {
			continue
		}
		ref, err := ParsePath(it.Name)
		if err != nil {
			s.Stats.Rejected.Add(1)
			if !s.rejected[it.Name] {
				s.rejected[it.Name] = true
				s.log.Warn("stream not linked to an entity", "path", it.Name, "err", err)
			}
			continue
		}
		key := ref.Key()
		seen[key] = true
		rec := Build(ref, it, s.cfg.Media, s.cfg.Endpoints, s.cfg.Origin, now)
		prev, had := s.state[key]
		if had && sameContent(prev, rec) {
			continue
		}
		if s.put(ctx, key, rec) && (!had || !prev.Live) {
			s.emit(rec, "online", "")
		}
	}
	keys := make([]string, 0, len(s.state))
	for k := range s.state {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, key := range keys {
		prev := s.state[key]
		if seen[key] || !prev.Live {
			continue
		}
		off := prev
		off.Live = false
		off.Reason = offlineReason
		if off.Reason == "" {
			off.Reason = "publisher_stopped"
		}
		off.UpdatedAt = now.UTC()
		if s.put(ctx, key, off) {
			s.emit(off, "offline", off.Reason)
		}
	}
}

func (s *Service) put(ctx context.Context, key string, r Record) bool {
	s.seq++
	r.Seq = s.seq
	b, err := json.Marshal(r)
	if err != nil {
		s.log.Error("marshal record", "key", key, "err", err)
		return false
	}
	if _, err := s.kv.Put(ctx, key, b); err != nil {
		s.log.Warn("KV put failed; retried next poll", "key", key, "err", err)
		return false
	}
	s.state[key] = r
	s.Stats.Puts.Add(1)
	return true
}

func (s *Service) emit(r Record, event, reason string) {
	ev := Event{V: 1, Event: event, StreamID: r.StreamID, EntityType: r.EntityType, EntityID: r.EntityID,
		EntitySubject: r.EntitySubject, Reason: reason, TS: time.Now().UTC(), Origin: s.cfg.Origin}
	b, _ := json.Marshal(ev)
	if err := s.nc.Publish(EventSubject(r.EntityID, event), b); err != nil {
		s.log.Warn("event publish failed", "event", event, "err", err)
		return
	}
	s.Stats.Events.Add(1)
	s.log.Info("stream "+event, "stream", r.StreamID, "entity", r.EntitySubject)
}

func (s *Service) serveRPC() ([]*nats.Subscription, error) {
	get, err := s.nc.QueueSubscribe(SubjectGet, QueueGroup, s.handleGet)
	if err != nil {
		return nil, err
	}
	list, err := s.nc.QueueSubscribe(SubjectList, QueueGroup, s.handleList)
	if err != nil {
		_ = get.Unsubscribe()
		return nil, err
	}
	return []*nats.Subscription{get, list}, nil
}

func respond(m *nats.Msg, r Reply) {
	b, _ := json.Marshal(r)
	_ = m.Respond(b)
}

func (s *Service) handleGet(m *nats.Msg) {
	s.Stats.RPC.Add(1)
	var req GetRequest
	data := bytes.TrimSpace(m.Data)
	if len(data) > 0 && data[0] == '{' {
		if err := json.Unmarshal(data, &req); err != nil {
			respond(m, Reply{Error: "bad_request"})
			return
		}
	} else {
		req.Entity = string(data)
	}
	ref, err := ParseEntityRef(req.Entity)
	if err != nil {
		respond(m, Reply{Error: "bad_request: " + err.Error()})
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	e, err := s.kv.Get(ctx, ref.Key())
	if errors.Is(err, jetstream.ErrKeyNotFound) {
		respond(m, Reply{Error: "not_found"})
		return
	}
	if err != nil {
		respond(m, Reply{Error: "unavailable"})
		return
	}
	var r Record
	if err := json.Unmarshal(e.Value(), &r); err != nil {
		respond(m, Reply{Error: "corrupt"})
		return
	}
	respond(m, Reply{Record: &r})
}

func (s *Service) handleList(m *nats.Msg) {
	s.Stats.RPC.Add(1)
	s.mu.Lock()
	out := make([]Record, 0, len(s.state))
	for _, r := range s.state {
		out = append(out, r)
	}
	s.mu.Unlock()
	sort.Slice(out, func(i, j int) bool { return out[i].StreamID < out[j].StreamID })
	respond(m, Reply{Records: out})
}

func equalJSON(a, b any) bool {
	x, err1 := json.Marshal(a)
	y, err2 := json.Marshal(b)
	return err1 == nil && err2 == nil && bytes.Equal(x, y)
}
