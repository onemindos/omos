package projector

// NATS wiring: ENT (ordered consumer, last per subject then every change; core subscription when
// ENT is absent), KV watches on `geofences` and `mission-state`, the legacy control subjects,
// the read-only tasking guard on TAK input subjects, and the publisher.

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

type evKind int

const (
	evEntity evKind = iota
	evDelete
	evZone
	evZoneDel
	evMission
	evMissionDel
	evEnable
	evDisable
)

type event struct {
	kind evKind
	key  string
	data []byte
	ent  *Entity
}

// Service runs one projector instance.
type Service struct {
	cfg    Config
	nc     *nats.Conn
	js     jetstream.JetStream
	log    *slog.Logger
	Engine *Engine
	events chan event
	ready  atomic.Bool
	entOK  atomic.Bool
	kvOK   atomic.Int32
	statMu sync.Mutex
	stats  map[string]any
	now    func() time.Time
	done   <-chan struct{}
}

// emit hands an event to the loop; it gives up once Run has ended.
func (s *Service) emit(ev event) {
	select {
	case s.events <- ev:
	case <-s.done:
	}
}

// New builds a service on an open connection.
func New(cfg Config, nc *nats.Conn, log *slog.Logger) (*Service, error) {
	js, err := jetstream.New(nc)
	if err != nil {
		return nil, err
	}
	s := &Service{cfg: cfg, nc: nc, js: js, log: log, events: make(chan event, 8192), now: time.Now}
	s.Engine = NewEngine(cfg, log, s.publish)
	return s, nil
}

// Ready reports that ENT and both KV watches finished their initial load.
func (s *Service) Ready() bool { return s.ready.Load() }

// Status is the /status body.
func (s *Service) Status() map[string]any {
	s.statMu.Lock()
	st := map[string]any{}
	for k, v := range s.stats {
		st[k] = v
	}
	s.statMu.Unlock()
	st["counters"] = s.Engine.C.Snapshot()
	st["target"] = s.cfg.Target
	st["subject"] = s.cfg.Subject
	st["ready"] = s.Ready()
	return st
}

func (s *Service) publish(b Built, e *Entity) error {
	m := nats.NewMsg(s.cfg.Subject)
	m.Data = b.XML
	now := s.now()
	m.Header.Set("Nats-Msg-Id", b.UID+":"+strconv.FormatInt(now.UnixMilli(), 10))
	m.Header.Set("OneMind-Marking", strings.ToUpper(e.Marking.TLP))
	if len(e.Marking.ReleasableTo) > 0 {
		m.Header.Set("OneMind-Releasable-To", strings.Join(e.Marking.ReleasableTo, ","))
	}
	m.Header.Set("OneMind-Projector", s.cfg.Instance)
	return s.nc.PublishMsg(m) // core publish, no reply subject (§6.3)
}

// Run blocks until ctx ends.
func (s *Service) Run(ctx context.Context) error {
	s.done = ctx.Done()
	var subs []*nats.Subscription
	defer func() {
		for _, sub := range subs {
			_ = sub.Unsubscribe()
		}
	}()
	for subj, kind := range map[string]evKind{s.cfg.ControlDisable: evDisable, s.cfg.ControlEnable: evEnable} {
		kind := kind
		sub, err := s.nc.Subscribe(subj, func(*nats.Msg) { s.emit(event{kind: kind}) })
		if err != nil {
			return err
		}
		subs = append(subs, sub)
	}
	for _, subj := range s.cfg.TakInputSubjects {
		sub, err := s.nc.Subscribe(subj, s.guard)
		if err != nil {
			return err
		}
		subs = append(subs, sub)
	}
	var wg sync.WaitGroup
	defer wg.Wait()
	wg.Add(3)
	go func() { defer wg.Done(); s.watchKV(ctx, s.cfg.GeofenceBucket, evZone, evZoneDel) }()
	go func() { defer wg.Done(); s.watchKV(ctx, s.cfg.MissionBucket, evMission, evMissionDel) }()
	go func() { defer wg.Done(); s.followEnt(ctx, &subs) }()
	return s.loop(ctx)
}

func (s *Service) loop(ctx context.Context) error {
	tick := time.NewTicker(s.cfg.Tick)
	defer tick.Stop()
	statEvery := time.NewTicker(time.Second)
	defer statEvery.Stop()
	for {
		select {
		case <-ctx.Done():
			return nil
		case ev := <-s.events:
			now := s.now()
			switch ev.kind {
			case evEntity:
				s.Engine.Upsert(ev.ent, now)
			case evDelete:
				s.Engine.Delete(ev.key, now)
			case evZone:
				s.Engine.SetZone(ev.key, ev.data, now)
			case evZoneDel:
				s.Engine.DelZone(ev.key, now)
			case evMission:
				s.Engine.SetMission(ev.key, ev.data, now)
			case evMissionDel:
				s.Engine.DelMission(ev.key, now)
			case evEnable:
				s.log.Info("projector enabled", "subject", s.cfg.ControlEnable)
				s.Engine.SetEnabled(true, now)
			case evDisable:
				s.log.Warn("projector disabled", "subject", s.cfg.ControlDisable)
				s.Engine.SetEnabled(false, now)
			}
		case <-tick.C:
			s.Engine.Tick(s.now())
		case <-statEvery.C:
			st := s.Engine.Stats()
			s.statMu.Lock()
			s.stats = st
			s.statMu.Unlock()
		}
	}
}

func (s *Service) markReady() {
	if s.entOK.Load() && s.kvOK.Load() >= 2 && !s.ready.Swap(true) {
		s.log.Info("projector ready", "target", s.cfg.Target, "subject", s.cfg.Subject)
	}
}

// watchKV follows a bucket; a missing bucket counts as empty and is retried.
func (s *Service) watchKV(ctx context.Context, bucket string, set, del evKind) {
	counted := false
	for ctx.Err() == nil {
		kv, err := s.js.KeyValue(ctx, bucket)
		if err == nil {
			var w jetstream.KeyWatcher
			if w, err = kv.WatchAll(ctx); err == nil {
				s.log.Info("watching kv", "bucket", bucket)
				for ent := range w.Updates() {
					if ent == nil { // initial values delivered
						if !counted {
							counted = true
							s.kvOK.Add(1)
							s.markReady()
						}
						continue
					}
					switch ent.Operation() {
					case jetstream.KeyValuePut:
						s.emit(event{kind: set, key: ent.Key(), data: ent.Value()})
					default:
						s.emit(event{kind: del, key: ent.Key()})
					}
				}
				_ = w.Stop()
			}
		}
		if err != nil {
			s.log.Warn("kv bucket unavailable; treated as empty, retrying", "bucket", bucket, "err", err)
			if !counted {
				counted = true
				s.kvOK.Add(1)
				s.markReady()
			}
		}
		select {
		case <-ctx.Done():
		case <-time.After(5 * time.Second):
		}
	}
}

// followEnt reads ENT, or falls back to a core subscription on ent.*.*.state.
func (s *Service) followEnt(ctx context.Context, subs *[]*nats.Subscription) {
	if s.cfg.EntStream != "" {
		for ctx.Err() == nil {
			err := s.consumeEnt(ctx)
			if errors.Is(err, jetstream.ErrStreamNotFound) {
				s.log.Warn("ENT stream not found; using a core subscription", "subject", s.cfg.EntSubject)
				break
			}
			if ctx.Err() != nil {
				return
			}
			s.log.Warn("ENT consumer ended; restarting", "err", err)
			select {
			case <-ctx.Done():
				return
			case <-time.After(2 * time.Second):
			}
		}
		if ctx.Err() != nil {
			return
		}
	}
	sub, err := s.nc.Subscribe(s.cfg.EntSubject, func(m *nats.Msg) {
		s.onEnt(m.Subject, m.Data, m.Header.Get("Nats-Marker-Reason"))
	})
	if err != nil {
		s.log.Error("core subscription failed", "err", err)
		return
	}
	_ = sub // unsubscribed when the connection closes
	s.entOK.Store(true)
	s.markReady()
	<-ctx.Done()
	_ = sub.Unsubscribe()
}

func (s *Service) consumeEnt(ctx context.Context) error {
	st, err := s.js.Stream(ctx, s.cfg.EntStream)
	if err != nil {
		return err
	}
	info, err := st.Info(ctx)
	if err != nil {
		return err
	}
	cfg := jetstream.OrderedConsumerConfig{FilterSubjects: []string{s.cfg.EntSubject},
		DeliverPolicy: jetstream.DeliverLastPerSubjectPolicy}
	if info.State.Msgs == 0 {
		cfg.DeliverPolicy, cfg.OptStartSeq = jetstream.DeliverByStartSequencePolicy, info.State.LastSeq+1
		s.entOK.Store(true)
		s.markReady()
	}
	oc, err := s.js.OrderedConsumer(ctx, s.cfg.EntStream, cfg)
	if err != nil {
		return err
	}
	it, err := oc.Messages(jetstream.PullMaxMessages(1000))
	if err != nil {
		return err
	}
	stop := make(chan struct{})
	defer close(stop)
	go func() {
		select {
		case <-ctx.Done():
			it.Stop()
		case <-stop:
		}
	}()
	defer it.Stop()
	s.log.Info("following ENT", "stream", s.cfg.EntStream)
	for {
		msg, err := it.Next()
		if err != nil {
			return err
		}
		s.onEnt(msg.Subject(), msg.Data(), msg.Headers().Get("Nats-Marker-Reason"))
		if meta, err := msg.Metadata(); err == nil && meta.NumPending == 0 {
			s.entOK.Store(true)
			s.markReady()
		}
	}
}

func (s *Service) onEnt(subject string, data []byte, marker string) {
	tok := strings.Split(subject, ".")
	if len(tok) != 4 || tok[0] != "ent" || tok[3] != "state" {
		s.Engine.C.Inc("ent_bad_subject")
		return
	}
	key := tok[1] + "." + tok[2]
	if len(data) == 0 || marker != "" {
		s.emit(event{kind: evDelete, key: key})
		return
	}
	var e Entity
	if err := json.Unmarshal(data, &e); err != nil {
		s.Engine.C.Inc("ent_unreadable")
		return
	}
	if e.EntityType == "" {
		e.EntityType = tok[1]
	}
	if e.EntityID == "" {
		e.EntityID = tok[2]
	}
	if e.EntityType != tok[1] || e.EntityID != tok[2] {
		s.Engine.C.Inc("ent_subject_mismatch")
		return
	}
	s.emit(event{kind: evEntity, ent: &e})
}

var (
	cotTypeAttr = regexp.MustCompile(`<event\b[^>]*?\btype\s*=\s*["']([^"']+)["']`)
	cotUIDAttr  = regexp.MustCompile(`<event\b[^>]*?\buid\s*=\s*["']([^"']+)["']`)
)

// guard watches CoT arriving from TAK (tak.cot.out, tak.cot.mesh.rx) for tasking. It never
// creates, forwards or executes anything: a strike or engage request (t-k-*, t-e-*) is refused
// and logged (stop line SL1, JADC2-MAP §5.3 test 8); other tasking is logged as not supported
// until the policy path of cot-projection.md §8.3 exists.
func (s *Service) guard(m *nats.Msg) {
	head := m.Data
	if len(head) > 4096 {
		head = head[:4096]
	}
	mt := cotTypeAttr.FindSubmatch(head)
	if mt == nil {
		return
	}
	t := string(mt[1])
	if !strings.HasPrefix(t, "t-") || strings.HasPrefix(t, "t-x-") {
		return
	}
	uid := ""
	if mu := cotUIDAttr.FindSubmatch(head); mu != nil {
		uid = string(mu[1])
	}
	if strings.HasPrefix(t, "t-k") || strings.HasPrefix(t, "t-e") {
		s.Engine.C.Inc("tak_tasking_refused")
		s.log.Warn("refused TAK tasking: strike or engage CoT creates no task (SL1)",
			"cot_type", t, "uid", uid, "subject", m.Subject, "gateway", m.Header.Get("Tak-Gateway"))
		return
	}
	s.Engine.C.Inc("tak_tasking_ignored")
	s.log.Info("TAK tasking not supported yet; ignored (cot-projection.md §8.3)", "cot_type", t, "uid", uid, "subject", m.Subject)
}
