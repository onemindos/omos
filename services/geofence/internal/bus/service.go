// Package bus connects the geofence engine to NATS:
//
//	in   cmd.map.geofence.{id}.{create|update|delete|remove}  CMD durable consumer (policy-signed)
//	in   KV geofences (watch)                                 zones, including hand edits
//	in   ent.*.*.state                                        ENT ordered consumer (last per subject)
//	in   tel.feed.*.status                                    core subscription (TEL once it binds it)
//	out  KV geofences                                         the only writer
//	out  evt.map.geofence.{id}.{enter|exit}                   MISSION (JetStream ack, Nats-Msg-Id)
//	out  evt.map.geofence.{created|updated|deleted|refused}     command results, zone id in data (proposal D1-P2)
//	out  evt.so.alert.{raised|updated|resolved}               EVT (proposal D1-P3)
//	rpc  svc.geofence.alerts                                  open alerts
//
// On start it replays MISSION crossings and EVT alerts to rebuild membership and the open
// alert set before it reads zones and entities, so a restart publishes nothing new for a
// world that did not change.
package bus

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"regexp"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"

	"github.com/onemindos/omos/services/geofence/internal/engine"
	"github.com/onemindos/omos/services/geofence/internal/sig"
	"github.com/onemindos/omos/services/geofence/internal/zone"
)

// Config is the service configuration.
type Config struct {
	Source           string
	GeofenceBucket   string
	EntStream        string
	EntCore          bool // interim: core subscription on ent.*.*.state, no ENT stream
	CmdStream        string
	CmdDurable       string
	MissionStream    string
	EvtStream        string
	RestoreLookback  time.Duration // 0 = replay everything
	FeedSubject      string
	Keys             sig.Keys
	RequireSignature bool
	RequireApproval  bool
	ClockSkew        time.Duration
	TickEvery        time.Duration
	ServeRPC         bool
	Log              *slog.Logger
	Now              func() time.Time
}

// Defaults fills unset fields.
func (c *Config) Defaults() {
	set := func(p *string, v string) {
		if *p == "" {
			*p = v
		}
	}
	set(&c.Source, "omos-geofence")
	set(&c.GeofenceBucket, "geofences")
	set(&c.EntStream, "ENT")
	set(&c.CmdStream, "CMD")
	set(&c.CmdDurable, "geofence-svc-cmd")
	set(&c.MissionStream, "MISSION")
	set(&c.EvtStream, "EVT")
	set(&c.FeedSubject, "tel.feed.*.status")
	if c.ClockSkew <= 0 {
		c.ClockSkew = 5 * time.Second
	}
	if c.TickEvery <= 0 {
		c.TickEvery = 5 * time.Second
	}
	if c.Log == nil {
		c.Log = slog.Default()
	}
	if c.Now == nil {
		c.Now = time.Now
	}
}

// Service runs the engine on the bus.
type Service struct {
	cfg   Config
	nc    *nats.Conn
	js    jetstream.JetStream
	eng   *engine.Engine
	kv    jetstream.KeyValue
	out   chan engine.Output
	ready atomic.Bool

	seenMu sync.Mutex
	seen   map[string]bool
	ring   []string

	cntMu sync.Mutex
	cnt   map[string]int64
}

// New builds a service.
func New(nc *nats.Conn, eng *engine.Engine, cfg Config) (*Service, error) {
	cfg.Defaults()
	if cfg.RequireSignature && len(cfg.Keys) == 0 {
		return nil, errors.New("signature required but no policy public keys configured")
	}
	js, err := jetstream.New(nc)
	if err != nil {
		return nil, err
	}
	return &Service{cfg: cfg, nc: nc, js: js, eng: eng, out: make(chan engine.Output, 10000),
		seen: map[string]bool{}, cnt: map[string]int64{}}, nil
}

// Ready reports whether zones and entities have been loaded.
func (s *Service) Ready() bool { return s.ready.Load() }

func (s *Service) inc(k string) {
	s.cntMu.Lock()
	s.cnt[k]++
	s.cntMu.Unlock()
}

// Counters merges service and engine counters.
func (s *Service) Counters() map[string]int64 {
	out := s.eng.Counters()
	s.cntMu.Lock()
	for k, v := range s.cnt {
		out[k] = v
	}
	s.cntMu.Unlock()
	return out
}

// Run blocks until ctx ends or a fatal error.
func (s *Service) Run(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	kv, err := s.js.KeyValue(ctx, s.cfg.GeofenceBucket)
	if err != nil {
		return fmt.Errorf("KV %s: %w (the service never creates buckets)", s.cfg.GeofenceBucket, err)
	}
	s.kv = kv
	if err := s.restore(ctx); err != nil {
		return err
	}
	errc := make(chan error, 8)
	var wg sync.WaitGroup
	start := func(name string, f func(context.Context) error) {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := f(ctx); err != nil && ctx.Err() == nil {
				errc <- fmt.Errorf("%s: %w", name, err)
			}
		}()
	}
	start("publisher", s.publisher)
	zonesLoaded := make(chan struct{})
	start("zones", func(ctx context.Context) error { return s.watchZones(ctx, zonesLoaded) })
	select {
	case <-zonesLoaded:
	case err := <-errc:
		return err
	case <-ctx.Done():
		return nil
	}
	entsLoaded := make(chan struct{})
	start("entities", func(ctx context.Context) error { return s.consumeEntities(ctx, entsLoaded) })
	select {
	case <-entsLoaded:
	case err := <-errc:
		return err
	case <-ctx.Done():
		return nil
	}
	sub, err := s.nc.Subscribe(s.cfg.FeedSubject, func(m *nats.Msg) {
		parts := strings.Split(m.Subject, ".")
		if len(parts) != 4 {
			return
		}
		s.emit(s.eng.FeedStatus(parts[2], m.Data))
	})
	if err != nil {
		return err
	}
	defer func() { _ = sub.Unsubscribe() }()
	start("commands", s.consumeCommands)
	start("ticker", s.ticker)
	if s.cfg.ServeRPC {
		rpc, err := s.nc.Subscribe("svc.geofence.alerts", s.serveAlerts)
		if err != nil {
			return err
		}
		defer func() { _ = rpc.Unsubscribe() }()
	}
	s.ready.Store(true)
	s.cfg.Log.Info("geofence service ready", "zones", s.eng.Counters()["zones"], "entities", s.eng.Counters()["entities"],
		"alerts_open", len(s.eng.OpenAlerts()), "signature_required", s.cfg.RequireSignature, "keys", s.cfg.Keys.IDs())
	select {
	case err := <-errc:
		cancel()
		wg.Wait()
		return err
	case <-ctx.Done():
		wg.Wait()
		return nil
	}
}

func (s *Service) emit(outs []engine.Output) {
	for _, o := range outs {
		s.out <- o
	}
}

// ---- restore ----

func (s *Service) replay(ctx context.Context, stream string, filters []string, fn func(jetstream.Msg)) (int, error) {
	st, err := s.js.Stream(ctx, stream)
	if err != nil {
		return 0, fmt.Errorf("stream %s: %w", stream, err)
	}
	info, err := st.Info(ctx, jetstream.WithSubjectFilter(filters[0]))
	if err != nil {
		return 0, err
	}
	if info.State.Msgs == 0 {
		return 0, nil
	}
	occ := jetstream.OrderedConsumerConfig{FilterSubjects: filters, DeliverPolicy: jetstream.DeliverAllPolicy}
	if s.cfg.RestoreLookback > 0 {
		t := s.cfg.Now().Add(-s.cfg.RestoreLookback)
		occ.DeliverPolicy, occ.OptStartTime = jetstream.DeliverByStartTimePolicy, &t
	}
	oc, err := s.js.OrderedConsumer(ctx, stream, occ)
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
		for m := range batch.Messages() {
			got++
			n++
			fn(m)
			if meta, err := m.Metadata(); err == nil && meta.NumPending == 0 {
				done = true
			}
		}
		if err := batch.Error(); err != nil && !errors.Is(err, nats.ErrTimeout) && !errors.Is(err, context.DeadlineExceeded) {
			return n, err
		}
		if done || got == 0 {
			return n, nil
		}
	}
}

type envelope struct {
	Data json.RawMessage `json:"data"`
}

func (s *Service) restore(ctx context.Context) error {
	nc, err := s.replay(ctx, s.cfg.MissionStream, []string{"evt.map.geofence.*.enter", "evt.map.geofence.*.exit"}, func(m jetstream.Msg) {
		var env envelope
		var c engine.Crossing
		if json.Unmarshal(m.Data(), &env) == nil && json.Unmarshal(env.Data, &c) == nil {
			s.eng.RestoreCrossing(c)
		}
	})
	if err != nil {
		return fmt.Errorf("replay crossings: %w", err)
	}
	na, err := s.replay(ctx, s.cfg.EvtStream, []string{"evt.so.alert.>"}, func(m jetstream.Msg) {
		var env envelope
		var a engine.Alert
		if json.Unmarshal(m.Data(), &env) == nil && json.Unmarshal(env.Data, &a) == nil {
			s.eng.RestoreAlert(a)
		}
	})
	if err != nil {
		return fmt.Errorf("replay alerts: %w", err)
	}
	s.eng.Restored()
	s.cfg.Log.Info("restored", "crossings", nc, "alert_facts", na, "alerts_open", len(s.eng.OpenAlerts()))
	return nil
}

// ---- inputs ----

func (s *Service) watchZones(ctx context.Context, loaded chan<- struct{}) error {
	w, err := s.kv.WatchAll(ctx)
	if err != nil {
		return err
	}
	defer func() { _ = w.Stop() }()
	initial := true
	for {
		select {
		case <-ctx.Done():
			return nil
		case e, ok := <-w.Updates():
			if !ok {
				return errors.New("geofences watcher stopped")
			}
			if e == nil {
				if initial {
					initial = false
					s.emit(s.eng.ReconcileZones())
					close(loaded)
				}
				continue
			}
			switch e.Operation() {
			case jetstream.KeyValuePut:
				z, err := zone.Parse(e.Key(), e.Value(), e.Revision())
				if err != nil {
					// An unreadable zone is dropped (and its crossings closed), never guessed at.
					s.inc("zone_unreadable")
					s.cfg.Log.Warn("unreadable geofence ignored", "key", e.Key(), "err", err)
					s.emit(s.eng.DeleteZone(e.Key()))
					continue
				}
				s.emit(s.eng.SetZone(z))
			default:
				s.emit(s.eng.DeleteZone(e.Key()))
			}
		}
	}
}

func (s *Service) consumeEntities(ctx context.Context, loaded chan<- struct{}) error {
	handle := func(subject string, data []byte, hdr nats.Header) {
		parts := strings.Split(subject, ".")
		if len(parts) != 4 {
			return
		}
		if len(data) == 0 || (hdr != nil && hdr.Get("Nats-Marker-Reason") != "") {
			s.emit(s.eng.RemoveEntity(parts[2]))
			return
		}
		e, err := engine.ParseEntity(data)
		if err != nil {
			s.inc("entity_unreadable")
			return
		}
		s.emit(s.eng.SetEntity(e))
	}
	if s.cfg.EntCore {
		sub, err := s.nc.Subscribe("ent.*.*.state", func(m *nats.Msg) { handle(m.Subject, m.Data, m.Header) })
		if err != nil {
			return err
		}
		defer func() { _ = sub.Unsubscribe() }()
		close(loaded)
		<-ctx.Done()
		return nil
	}
	st, err := s.js.Stream(ctx, s.cfg.EntStream)
	if err != nil {
		return fmt.Errorf("stream %s: %w", s.cfg.EntStream, err)
	}
	info, err := st.Info(ctx)
	if err != nil {
		return err
	}
	oc, err := s.js.OrderedConsumer(ctx, s.cfg.EntStream, jetstream.OrderedConsumerConfig{
		FilterSubjects: []string{"ent.*.*.state"}, DeliverPolicy: jetstream.DeliverLastPerSubjectPolicy})
	if err != nil {
		return err
	}
	initial := info.State.Msgs > 0
	if !initial {
		s.emit(s.eng.ReconcileEntities())
		close(loaded)
	}
	for ctx.Err() == nil {
		batch, err := oc.Fetch(500, jetstream.FetchMaxWait(time.Second))
		if err != nil {
			if ctx.Err() != nil {
				return nil
			}
			return err
		}
		got := 0
		for m := range batch.Messages() {
			got++
			handle(m.Subject(), m.Data(), m.Headers())
			if initial {
				if meta, err := m.Metadata(); err == nil && meta.NumPending == 0 {
					initial = false
					s.emit(s.eng.ReconcileEntities())
					close(loaded)
				}
			}
		}
		if initial && got == 0 {
			initial = false
			s.emit(s.eng.ReconcileEntities())
			close(loaded)
		}
	}
	return nil
}

func (s *Service) ticker(ctx context.Context) error {
	t := time.NewTicker(s.cfg.TickEvery)
	defer t.Stop()
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-t.C:
			s.emit(s.eng.Tick())
		}
	}
}

func (s *Service) serveAlerts(m *nats.Msg) {
	var q struct {
		Kind     string `json:"kind"`
		EntityID string `json:"entity_id"`
		Feed     string `json:"feed"`
	}
	_ = json.Unmarshal(m.Data, &q)
	alerts := []engine.Alert{}
	for _, a := range s.eng.OpenAlerts() {
		if (q.Kind == "" || a.Kind == q.Kind) && (q.EntityID == "" || a.About.EntityID == q.EntityID) &&
			(q.Feed == "" || a.About.Feed == q.Feed) {
			alerts = append(alerts, a)
		}
	}
	b, _ := json.Marshal(map[string]any{"v": 1, "alerts": alerts})
	_ = m.Respond(b)
}

// ---- output ----

// UUIDFromKey derives a stable UUID v4-shaped id from a message key, so a republished fact
// carries the same envelope id as the first attempt.
func UUIDFromKey(key string) string {
	h := sha256.Sum256([]byte(key))
	h[6] = (h[6] & 0x0f) | 0x40
	h[8] = (h[8] & 0x3f) | 0x80
	return fmt.Sprintf("%x-%x-%x-%x-%x", h[0:4], h[4:6], h[6:8], h[8:10], h[10:16])
}

func ts(t time.Time) string { return t.UTC().Format("2006-01-02T15:04:05.000Z") }

var idTokenRe = regexp.MustCompile(`^evt\.map\.geofence\.[^.]+\.`)

func eventType(subject string) string {
	if idTokenRe.MatchString(subject) {
		p := strings.Split(subject, ".")
		return "evt.map.geofence." + p[len(p)-1]
	}
	return subject
}

func (s *Service) envelope(subject, key string, data any, corr, causation string) ([]byte, error) {
	env := map[string]any{"v": 1, "id": UUIDFromKey(key), "ts": ts(s.cfg.Now()), "type": eventType(subject),
		"source": s.cfg.Source, "subject": subject, "data": data}
	if corr != "" {
		env["corr"] = corr
	}
	if causation != "" {
		env["causation"] = causation
	}
	return json.Marshal(env)
}

// publisher sends outputs in order with a JetStream ack. A failed publish is retried with the
// same Nats-Msg-Id until it is acked (the stream deduplicates a retry that did land), so
// nothing is skipped or reordered; a full queue blocks the inputs (backpressure).
func (s *Service) publisher(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		case o := <-s.out:
			b, err := s.envelope(o.Subject, o.MsgID, o.Data, "", "")
			if err != nil {
				s.cfg.Log.Error("cannot encode output", "subject", o.Subject, "err", err)
				continue
			}
			backoff := 100 * time.Millisecond
			for {
				pctx, cancel := context.WithTimeout(ctx, 5*time.Second)
				_, err := s.js.Publish(pctx, o.Subject, b, jetstream.WithMsgID(o.MsgID))
				cancel()
				if err == nil {
					s.inc("published")
					break
				}
				if ctx.Err() != nil {
					return nil
				}
				s.inc("publish_retries")
				s.cfg.Log.Warn("publish failed, retrying", "subject", o.Subject, "err", err)
				select {
				case <-time.After(backoff):
				case <-ctx.Done():
					return nil
				}
				if backoff < 5*time.Second {
					backoff *= 2
				}
			}
		}
	}
}

// ---- commands ----

func (s *Service) consumeCommands(ctx context.Context) error {
	cons, err := s.js.CreateOrUpdateConsumer(ctx, s.cfg.CmdStream, jetstream.ConsumerConfig{
		Durable: s.cfg.CmdDurable, FilterSubject: "cmd.map.geofence.>", AckPolicy: jetstream.AckExplicitPolicy,
		DeliverPolicy: jetstream.DeliverNewPolicy, AckWait: 30 * time.Second, MaxDeliver: 10,
		MaxAckPending: 64, Description: "geofence service: policy-approved zone writes",
	})
	if err != nil {
		return fmt.Errorf("consumer on %s: %w", s.cfg.CmdStream, err)
	}
	cc, err := cons.Consume(func(m jetstream.Msg) { s.handleCommand(ctx, m) })
	if err != nil {
		return err
	}
	defer cc.Stop()
	<-ctx.Done()
	return nil
}

func (s *Service) remember(id string) bool {
	s.seenMu.Lock()
	defer s.seenMu.Unlock()
	if s.seen[id] {
		return false
	}
	s.seen[id] = true
	s.ring = append(s.ring, id)
	if len(s.ring) > 4096 {
		delete(s.seen, s.ring[0])
		s.ring = s.ring[1:]
	}
	return true
}

func (s *Service) forget(id string) {
	s.seenMu.Lock()
	delete(s.seen, id)
	s.seenMu.Unlock()
}

// Command is the part of the policy command envelope (asset-command.md §6.3) read here.
type Command struct {
	V         int             `json:"v"`
	ID        string          `json:"id"`
	Cmd       string          `json:"cmd"`
	Subject   string          `json:"subject"`
	Source    string          `json:"source"`
	Data      json.RawMessage `json:"data"`
	ActionID  string          `json:"action_id"`
	ExpiresAt string          `json:"expires_at"`
	Simulated bool            `json:"simulated"`
	Corr      string          `json:"corr"`
	Approval  *struct {
		Status    string `json:"status"`
		Approver  string `json:"approver"`
		Approvals []struct {
			By string `json:"by"`
		} `json:"approvals"`
	} `json:"approval"`
	Policy *struct {
		Allowed bool `json:"allowed"`
	} `json:"policy"`
}

// Result is the data of evt.map.geofence.{created|updated|deleted|refused} (four tokens, as
// the grammar evt.{domain}.{entity}.{event}; the zone id is in the data).
type Result struct {
	V         int    `json:"v"`
	ZoneID    string `json:"zone_id"`
	Verb      string `json:"verb"`
	Outcome   string `json:"outcome"`
	Code      string `json:"code,omitempty"`
	Message   string `json:"message,omitempty"`
	Revision  uint64 `json:"revision,omitempty"`
	ActionID  string `json:"action_id,omitempty"`
	CommandID string `json:"command_id,omitempty"`
	Actor     string `json:"actor,omitempty"`
}

var uuidRe = regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)

var verbs = map[string]string{"create": "created", "update": "updated", "delete": "deleted", "remove": "deleted"}

func (s *Service) handleCommand(ctx context.Context, m jetstream.Msg) {
	subj := m.Subject()
	p := strings.Split(subj, ".")
	if len(p) != 5 || verbs[p[4]] == "" || !zone.IDRe.MatchString(p[3]) {
		s.inc("cmd_bad_subject")
		s.cfg.Log.Warn("command on an unusable subject dropped", "subject", subj)
		_ = m.Term()
		return
	}
	id, verb := p[3], p[4]
	var c Command
	res := Result{V: 1, ZoneID: id, Verb: verb}
	refuse := func(code, msg string) {
		res.Outcome, res.Code, res.Message = "refused", code, msg
		s.finish(ctx, m, c, res)
	}
	if err := json.Unmarshal(m.Data(), &c); err != nil {
		refuse(zone.CodeSchemaInvalid, "command is not a JSON envelope")
		return
	}
	res.ActionID, res.CommandID = c.ActionID, c.ID
	if c.Subject != subj || c.Cmd != subj || !uuidRe.MatchString(c.ID) {
		refuse(zone.CodeSchemaInvalid, "envelope subject, cmd and id must match the subject and be a UUID")
		return
	}
	if s.cfg.RequireSignature {
		if _, err := s.cfg.Keys.Verify(m.Data()); err != nil {
			code := zone.CodeBadSignature
			if errors.Is(err, sig.ErrUnsigned) {
				code = zone.CodeUnsigned
			}
			refuse(code, "zone writes need a policy-signed command (zone.* actions, level A2): "+err.Error())
			return
		}
		if c.ExpiresAt == "" {
			refuse(zone.CodeSchemaInvalid, "signed command has no expires_at")
			return
		}
	}
	if c.ExpiresAt != "" {
		exp, err := time.Parse(time.RFC3339Nano, c.ExpiresAt)
		if err != nil {
			refuse(zone.CodeSchemaInvalid, "expires_at is not RFC 3339")
			return
		}
		if s.cfg.Now().After(exp.Add(s.cfg.ClockSkew)) {
			refuse(zone.CodeExpired, "command expired at "+c.ExpiresAt)
			return
		}
	}
	actor := c.Source
	if s.cfg.RequireApproval {
		if c.Approval == nil || c.Approval.Status != "approved" || c.Policy == nil || !c.Policy.Allowed {
			refuse(zone.CodeNotApproved, "zone writes are level A2: the command must carry a human approval")
			return
		}
	}
	if c.Approval != nil && c.Approval.Approver != "" {
		actor = c.Approval.Approver
	}
	res.Actor = actor
	if !s.remember(c.ID) {
		s.inc("cmd_duplicate")
		_ = m.Ack()
		return
	}
	var params zone.Params
	if len(c.Data) > 0 && string(c.Data) != "null" {
		if err := json.Unmarshal(c.Data, &params); err != nil {
			s.forget(c.ID)
			refuse(zone.CodeSchemaInvalid, "data: "+err.Error())
			return
		}
	}
	meta := zone.Meta{Actor: actor, ActionID: c.ActionID, CommandID: c.ID, Simulated: c.Simulated, NowMs: s.cfg.Now().UnixMilli()}
	var rf *zone.Refusal
	var rev uint64
	var err error
	switch verb {
	case "create":
		rev, rf, err = s.create(ctx, id, params, meta)
	case "update":
		rev, rf, err = s.update(ctx, id, params, meta)
	default:
		rev, rf, err = s.remove(ctx, id, meta)
	}
	if err != nil {
		// Bus trouble, not a decision: retry the command later (it may expire meanwhile).
		s.forget(c.ID)
		s.inc("cmd_retry")
		s.cfg.Log.Warn("zone write failed, will retry", "subject", subj, "err", err)
		_ = m.NakWithDelay(2 * time.Second)
		return
	}
	if rf != nil {
		refuse(rf.Code, rf.Message)
		return
	}
	res.Outcome, res.Revision = verbs[verb], rev
	s.finish(ctx, m, c, res)
}

// finish publishes the result fact, then acks the command. If the result cannot be stored the
// command is redelivered; the write itself is idempotent by command id.
func (s *Service) finish(ctx context.Context, m jetstream.Msg, c Command, res Result) {
	subj := "evt.map.geofence." + res.Outcome
	key := "res:" + res.CommandID + ":" + res.Outcome
	if res.CommandID == "" {
		meta, _ := m.Metadata()
		if meta != nil {
			key = fmt.Sprintf("res:cmdseq:%d", meta.Sequence.Stream)
		}
	}
	causation := ""
	if uuidRe.MatchString(c.ID) {
		causation = c.ID
	}
	corr := ""
	if uuidRe.MatchString(c.Corr) {
		corr = c.Corr
	}
	b, err := s.envelope(subj, key, res, corr, causation)
	if err == nil {
		pctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		_, err = s.js.Publish(pctx, subj, b, jetstream.WithMsgID(key))
		cancel()
	}
	if err != nil {
		s.forget(c.ID)
		s.cfg.Log.Warn("result not stored, command will be redelivered", "subject", subj, "err", err)
		_ = m.NakWithDelay(2 * time.Second)
		return
	}
	s.inc("cmd_" + res.Outcome)
	s.cfg.Log.Info("zone command", "zone", res.ZoneID, "verb", res.Verb, "outcome", res.Outcome, "code", res.Code,
		"action_id", res.ActionID, "command_id", res.CommandID, "actor", res.Actor)
	_ = m.Ack()
}

func (s *Service) getRecord(ctx context.Context, id string) (*zone.Record, uint64, error) {
	e, err := s.kv.Get(ctx, id)
	if errors.Is(err, jetstream.ErrKeyNotFound) {
		return nil, 0, nil
	}
	if err != nil {
		return nil, 0, err
	}
	var r zone.Record
	if err := json.Unmarshal(e.Value(), &r); err != nil {
		return nil, e.Revision(), fmt.Errorf("existing record %s unreadable: %w", id, err)
	}
	return &r, e.Revision(), nil
}

func (s *Service) create(ctx context.Context, id string, p zone.Params, m zone.Meta) (uint64, *zone.Refusal, error) {
	prev, rev, err := s.getRecord(ctx, id)
	if err != nil {
		return 0, nil, err
	}
	if prev != nil {
		if prev.CommandID == m.CommandID {
			return rev, nil, nil // redelivered: already applied
		}
		return 0, &zone.Refusal{Code: zone.CodeExists, Message: "zone " + id + " exists; use update"}, nil
	}
	rec, rf := zone.Build(id, nil, p, m)
	if rf != nil {
		return 0, rf, nil
	}
	b, _ := json.Marshal(rec)
	rev, err = s.kv.Create(ctx, id, b)
	if errors.Is(err, jetstream.ErrKeyExists) {
		return 0, &zone.Refusal{Code: zone.CodeExists, Message: "zone " + id + " exists; use update"}, nil
	}
	return rev, nil, err
}

func (s *Service) update(ctx context.Context, id string, p zone.Params, m zone.Meta) (uint64, *zone.Refusal, error) {
	prev, rev, err := s.getRecord(ctx, id)
	if err != nil {
		return 0, nil, err
	}
	if prev == nil {
		return 0, &zone.Refusal{Code: zone.CodeNotFound, Message: "zone " + id + " does not exist"}, nil
	}
	if prev.CommandID == m.CommandID {
		return rev, nil, nil
	}
	if p.Revision != nil && *p.Revision != rev {
		return 0, &zone.Refusal{Code: zone.CodeConflict,
			Message: fmt.Sprintf("zone %s is at revision %d, the edit was made on %d", id, rev, *p.Revision)}, nil
	}
	rec, rf := zone.Build(id, prev, p, m)
	if rf != nil {
		return 0, rf, nil
	}
	b, _ := json.Marshal(rec)
	nrev, err := s.kv.Update(ctx, id, b, rev)
	if err != nil && isWrongSeq(err) {
		return 0, &zone.Refusal{Code: zone.CodeConflict, Message: "zone " + id + " changed during the edit"}, nil
	}
	return nrev, nil, err
}

func isWrongSeq(err error) bool {
	// R1 streams answer 10071; replicated (R3, as KV geofences is) answer 10164.
	var apiErr *jetstream.APIError
	return errors.As(err, &apiErr) && (apiErr.ErrorCode == jetstream.JSErrCodeStreamWrongLastSequence ||
		apiErr.ErrorCode == jetstream.JSErrCodeStreamWrongLastSequenceConstant)
}

// remove marks the record inactive with this command id (so a redelivery is recognised),
// then deletes the key; KV history keeps the last record.
func (s *Service) remove(ctx context.Context, id string, m zone.Meta) (uint64, *zone.Refusal, error) {
	prev, rev, err := s.getRecord(ctx, id)
	if err != nil {
		return 0, nil, err
	}
	if prev == nil {
		hist, herr := s.kv.History(ctx, id)
		if herr == nil {
			for i := len(hist) - 1; i >= 0; i-- {
				var r zone.Record
				if hist[i].Operation() == jetstream.KeyValuePut && json.Unmarshal(hist[i].Value(), &r) == nil {
					if r.CommandID == m.CommandID {
						return hist[i].Revision(), nil, nil
					}
					break
				}
			}
		}
		return 0, &zone.Refusal{Code: zone.CodeNotFound, Message: "zone " + id + " does not exist"}, nil
	}
	if prev.CommandID != m.CommandID {
		t := *prev
		t.Active, t.Live, t.Lifecycle = false, false, "removed"
		t.UpdatedT, t.UpdatedBy, t.ActionID, t.CommandID = m.NowMs, m.Actor, m.ActionID, m.CommandID
		b, _ := json.Marshal(t)
		rev, err = s.kv.Update(ctx, id, b, rev)
		if err != nil {
			if isWrongSeq(err) {
				return 0, &zone.Refusal{Code: zone.CodeConflict, Message: "zone " + id + " changed during the delete"}, nil
			}
			return 0, nil, err
		}
	}
	if err := s.kv.Delete(ctx, id, jetstream.LastRevision(rev)); err != nil {
		if isWrongSeq(err) {
			return 0, &zone.Refusal{Code: zone.CodeConflict, Message: "zone " + id + " changed during the delete"}, nil
		}
		return 0, nil, err
	}
	return rev, nil, nil
}
