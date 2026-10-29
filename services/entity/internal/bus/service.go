// Package bus connects the engine to NATS: durable pull consumers on the input streams,
// JetStream publish into `ENT` with Nats-Msg-Id and Nats-TTL, restore from `ENT` on start,
// dead letters, the KV lease and the `svc.entity.*` read endpoints.
package bus

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"sync/atomic"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/nats-io/nats.go/micro"

	"github.com/onemindos/omos/services/entity/internal/adapt"
	"github.com/onemindos/omos/services/entity/internal/engine"
	"github.com/onemindos/omos/services/entity/internal/model"
)

// Input is one source stream and the filters the service reads from it.
type Input struct {
	Stream  string
	Filters []string
}

// DefaultInputs are the streams of entity.md §1 "Inputs": COT (tak.cot.>), TEL (fleet
// telemetry), WORLD (raw air/sea), INTEL (feed-intel), DET (proposed; skipped while absent).
func DefaultInputs() []Input {
	return []Input{
		{Stream: "COT", Filters: []string{"tak.cot.>"}},
		{Stream: "TEL", Filters: []string{"tel.fleet.>"}},
		{Stream: "WORLD", Filters: []string{"evt.world.>"}},
		{Stream: "INTEL", Filters: []string{"intel.>"}},
		{Stream: "DET", Filters: []string{"det.>"}},
	}
}

// Config configures the service.
type Config struct {
	Inputs        []Input
	EntStream     string        // "ENT"
	CoreMode      bool          // interim: core publish, no ENT, no restore (entity.md §1)
	MsgTTL        bool          // send Nats-TTL (needs allowMsgTtl on ENT; BUS-DESIGN §4.1 prereq 2)
	DurablePrefix string        // "entity-svc"
	StartLookback time.Duration // where a NEW durable starts; 0 = deliver all
	Batch         int
	Idle          time.Duration
	SweepEvery    time.Duration
	RejectSubject string // `evt.sys.{source}.rejected` (BUS-P9)
	LeaseBucket   string
	LeaseKey      string
	LeaseTTL      time.Duration
	Holder        string
	ServeRPC      bool
	AckWait       time.Duration // redelivery delay for unacked input; default 30 s

	// crashAfter, when > 0, makes Run return errCrash right after the n-th ENT publish and
	// before acking the input message. Tests use it to prove a restart never duplicates.
	crashAfter int
}

var errCrash = errors.New("simulated crash after publish, before ack")

// Service is one running entity service.
type Service struct {
	cfg   Config
	nc    *nats.Conn
	js    jetstream.JetStream
	eng   *engine.Engine
	log   *slog.Logger
	ready atomic.Bool
	pubs  int
}

// New returns a service on an open connection.
func New(nc *nats.Conn, eng *engine.Engine, cfg Config, log *slog.Logger) (*Service, error) {
	js, err := jetstream.New(nc)
	if err != nil {
		return nil, err
	}
	if cfg.EntStream == "" {
		cfg.EntStream = "ENT"
	}
	if cfg.DurablePrefix == "" {
		cfg.DurablePrefix = "entity-svc"
	}
	if cfg.Batch == 0 {
		cfg.Batch = 256
	}
	if cfg.Idle == 0 {
		cfg.Idle = 100 * time.Millisecond
	}
	if cfg.SweepEvery == 0 {
		cfg.SweepEvery = time.Second
	}
	if cfg.RejectSubject == "" {
		cfg.RejectSubject = "evt.sys.entity.rejected"
	}
	if cfg.LeaseTTL == 0 {
		cfg.LeaseTTL = 15 * time.Second
	}
	if cfg.AckWait == 0 {
		cfg.AckWait = 30 * time.Second
	}
	if cfg.Inputs == nil {
		cfg.Inputs = DefaultInputs()
	}
	return &Service{cfg: cfg, nc: nc, js: js, eng: eng, log: log}, nil
}

// Ready reports whether the service holds the lease and is consuming.
func (s *Service) Ready() bool { return s.ready.Load() }

// Run holds the lease, restores, consumes and publishes until ctx ends or the lease is lost.
func (s *Service) Run(ctx context.Context) (retErr error) {
	ctx, cancel := context.WithCancelCause(ctx)
	defer cancel(nil)

	if s.cfg.LeaseBucket != "" {
		lease, err := NewLease(ctx, s.js, s.cfg.LeaseBucket, s.cfg.LeaseKey, s.cfg.Holder, s.cfg.LeaseTTL)
		if err != nil {
			return err
		}
		s.log.Info("waiting for lease", "bucket", s.cfg.LeaseBucket, "key", s.cfg.LeaseKey, "holder", s.cfg.Holder)
		if err := lease.Acquire(ctx); err != nil {
			return err
		}
		s.log.Info("lease acquired")
		defer func() {
			if !errors.Is(retErr, errCrash) { // a crash leaves the lease to expire
				lease.Release(context.Background())
			}
		}()
		go lease.Keep(ctx, func(err error) { cancel(err) })
	}

	if !s.cfg.CoreMode {
		if _, err := s.js.Stream(ctx, s.cfg.EntStream); err != nil {
			return fmt.Errorf("stream %s: %w (apply platform/nats/streams/ent.yaml, or run with ENT_MODE=core)", s.cfg.EntStream, err)
		}
		n, err := s.restore(ctx)
		if err != nil {
			return fmt.Errorf("restore from %s: %w", s.cfg.EntStream, err)
		}
		s.log.Info("restored entities", "count", n)
	}

	cons, err := s.consumers(ctx)
	if errors.Is(err, context.Canceled) {
		return nil
	}
	if err != nil {
		return err
	}
	if len(cons) == 0 {
		return errors.New("no input stream exists")
	}
	if s.cfg.ServeRPC {
		svc, err := s.serveRPC()
		if err != nil {
			return err
		}
		defer func() { _ = svc.Stop() }()
	}
	s.ready.Store(true)
	defer s.ready.Store(false)

	nextSweep := time.Now().Add(s.cfg.SweepEvery)
	for {
		if err := context.Cause(ctx); err != nil {
			if errors.Is(err, context.Canceled) {
				return nil
			}
			return err
		}
		got := 0
		for _, c := range cons {
			batch, err := c.FetchNoWait(s.cfg.Batch)
			if err != nil {
				return fmt.Errorf("fetch: %w", err)
			}
			for msg := range batch.Messages() {
				got++
				if err := s.handle(ctx, msg); err != nil {
					return err
				}
			}
		}
		if time.Now().After(nextSweep) {
			for _, p := range s.eng.Sweep() {
				if err := s.publish(ctx, p); err != nil {
					s.log.Warn("sweep publish failed", "subject", p.Subject, "err", err)
					s.eng.PublishFailed(p)
				}
			}
			nextSweep = time.Now().Add(s.cfg.SweepEvery)
		}
		if got == 0 {
			select {
			case <-ctx.Done():
			case <-time.After(s.cfg.Idle):
			}
		}
	}
}

func (s *Service) handle(ctx context.Context, msg jetstream.Msg) error {
	meta, err := msg.Metadata()
	if err != nil {
		return fmt.Errorf("metadata: %w", err)
	}
	obs, skip, perr := adapt.Parse(msg.Subject(), msg.Data())
	switch {
	case perr != nil:
		s.deadLetter(msg.Subject(), meta, "parse_error: "+perr.Error(), nil)
		s.eng.Count("parse_error")
		return msg.Ack()
	case obs == nil:
		s.eng.Count("skip_" + skip)
		return msg.Ack()
	}
	obs.Stream, obs.StreamSeq = meta.Stream, meta.Sequence.Stream
	res := s.eng.Apply(obs)
	if res.Reject != nil {
		s.deadLetter(msg.Subject(), meta, res.Reject.Reason, obs)
		return msg.Ack()
	}
	if res.Publish != nil {
		if err := s.publish(ctx, res.Publish); err != nil {
			s.log.Warn("publish failed; will retry", "subject", res.Publish.Subject, "err", err)
			s.eng.PublishFailed(res.Publish)
			_ = msg.NakWithDelay(2 * time.Second)
			time.Sleep(time.Second)
			return nil
		}
		s.pubs++
		if s.cfg.crashAfter > 0 && s.pubs >= s.cfg.crashAfter {
			return errCrash
		}
	}
	return msg.Ack()
}

// publish sends one entity record. In JetStream mode the PubAck confirms storage in ENT and a
// duplicate Nats-Msg-Id is reported, not stored twice.
func (s *Service) publish(ctx context.Context, p *engine.Publish) error {
	body, err := json.Marshal(p.Entity)
	if err != nil {
		return err
	}
	m := nats.NewMsg(p.Subject)
	m.Data = body
	m.Header.Set("Nats-Msg-Id", p.MsgID)
	if s.cfg.MsgTTL {
		m.Header.Set("Nats-TTL", p.TTL)
	}
	m.Header.Set("OneMind-Marking", p.Entity.Marking.TLP)
	m.Header.Set("OneMind-Releasable-To", strings.Join(p.Entity.Marking.ReleasableTo, ","))
	m.Header.Set("OneMind-Schema", "entity.v1")
	if s.cfg.CoreMode {
		return s.nc.PublishMsg(m)
	}
	ack, err := s.js.PublishMsg(ctx, m)
	if err != nil {
		return err
	}
	if ack.Duplicate {
		s.eng.Count("ent_duplicate_ack")
	}
	return nil
}

func (s *Service) deadLetter(subject string, meta *jetstream.MsgMetadata, reason string, obs *model.Observation) {
	rec := map[string]any{
		"reason": reason, "service": "entity-service", "subject": subject,
		"stream": meta.Stream, "seq": meta.Sequence.Stream, "t": time.Now().UnixMilli(),
	}
	if obs != nil {
		rec["integration"], rec["source_id"] = obs.Integration, obs.SourceID
		rec["entity_type"], rec["disposition"] = obs.EntityType, obs.Disposition
	}
	b, _ := json.Marshal(rec)
	if err := s.nc.Publish(s.cfg.RejectSubject, b); err != nil {
		s.log.Warn("dead letter publish failed", "err", err)
	}
	s.eng.Count("rejected")
}

// consumers opens (or creates once) one durable pull consumer per existing input stream.
// An existing durable is reused as is, so its position survives restarts and redeploys.
func (s *Service) consumers(ctx context.Context) ([]jetstream.Consumer, error) {
	var out []jetstream.Consumer
	for _, in := range s.cfg.Inputs {
		st, err := s.js.Stream(ctx, in.Stream)
		if ctx.Err() != nil {
			return nil, ctx.Err()
		}
		if err != nil {
			s.log.Warn("input stream missing; skipped", "stream", in.Stream, "err", err)
			continue
		}
		name := s.cfg.DurablePrefix + "-" + strings.ToLower(in.Stream)
		c, err := st.Consumer(ctx, name)
		if errors.Is(err, jetstream.ErrConsumerNotFound) {
			cc := jetstream.ConsumerConfig{
				Durable: name, FilterSubjects: in.Filters, AckPolicy: jetstream.AckExplicitPolicy,
				AckWait: s.cfg.AckWait, MaxAckPending: 1000, DeliverPolicy: jetstream.DeliverAllPolicy,
				Description: "omos entity service input",
			}
			if s.cfg.StartLookback > 0 {
				start := time.Now().Add(-s.cfg.StartLookback)
				cc.DeliverPolicy, cc.OptStartTime = jetstream.DeliverByStartTimePolicy, &start
			}
			c, err = st.CreateConsumer(ctx, cc)
		}
		if err != nil {
			return nil, fmt.Errorf("consumer %s on %s: %w", name, in.Stream, err)
		}
		out = append(out, c)
	}
	return out, nil
}

// restore reads the last record per subject from ENT with an ordered consumer.
func (s *Service) restore(ctx context.Context) (int, error) {
	st, err := s.js.Stream(ctx, s.cfg.EntStream)
	if err != nil {
		return 0, err
	}
	info, err := st.Info(ctx)
	if err != nil {
		return 0, err
	}
	if info.State.Msgs == 0 {
		return 0, nil
	}
	oc, err := s.js.OrderedConsumer(ctx, s.cfg.EntStream, jetstream.OrderedConsumerConfig{
		FilterSubjects: []string{"ent.>"}, DeliverPolicy: jetstream.DeliverLastPerSubjectPolicy,
	})
	if err != nil {
		return 0, err
	}
	n := 0
	for {
		batch, err := oc.Fetch(500, jetstream.FetchMaxWait(time.Second))
		if err != nil {
			return n, err
		}
		got, done := 0, false
		for msg := range batch.Messages() {
			got++
			meta, _ := msg.Metadata()
			if len(msg.Data()) > 0 && msg.Headers().Get("Nats-Marker-Reason") == "" {
				var e model.Entity
				if err := json.Unmarshal(msg.Data(), &e); err == nil {
					s.eng.Restore(e)
					n++
				}
			}
			if meta != nil && meta.NumPending == 0 {
				done = true
			}
		}
		if done || got == 0 {
			return n, nil
		}
	}
}

type queryReq struct {
	EntityID string    `json:"entity_id"`
	Types    []string  `json:"types"`
	BBox     []float64 `json:"bbox"`
	Limit    int       `json:"limit"`
}

// serveRPC registers `svc.entity.get` and `svc.entity.query` (entity.md §8, interim reads).
func (s *Service) serveRPC() (micro.Service, error) {
	svc, err := micro.AddService(s.nc, micro.Config{Name: "entity", Version: "0.1.0",
		Description: "OneMind entity service reads (entity.md §8)"})
	if err != nil {
		return nil, err
	}
	g := svc.AddGroup("svc.entity")
	_ = g.AddEndpoint("get", micro.HandlerFunc(func(r micro.Request) {
		var q queryReq
		if err := json.Unmarshal(r.Data(), &q); err != nil || q.EntityID == "" {
			_ = r.Error("400", "body must be {\"entity_id\": \"...\"}", nil)
			return
		}
		e, ok := s.eng.Get(q.EntityID)
		if !ok {
			_ = r.Error("404", "entity not found", nil)
			return
		}
		_ = r.RespondJSON(e)
	}))
	_ = g.AddEndpoint("query", micro.HandlerFunc(func(r micro.Request) {
		var q queryReq
		if len(r.Data()) > 0 {
			if err := json.Unmarshal(r.Data(), &q); err != nil {
				_ = r.Error("400", "invalid query", nil)
				return
			}
		}
		if q.Limit <= 0 || q.Limit > 5000 {
			q.Limit = 5000
		}
		_ = r.RespondJSON(map[string]any{"entities": s.eng.Query(q.Types, q.BBox, q.Limit)})
	}))
	return svc, nil
}
