package policy

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"sync"
	"time"

	"github.com/nats-io/nats.go"
)

// Subjects served and consumed (BUS-DESIGN.md sections 3.3-3.4; asset-command.md section 1).
const (
	SubjCommandPropose = "svc.policy.command.propose"
	SubjCommandConfirm = "svc.policy.command.confirm"
	SubjTaskPropose    = "svc.policy.task.propose"
	SubjTaskConfirm    = "svc.policy.task.confirm"
	SubjActionGet      = "svc.policy.action.get"
	SubjCheck          = "svc.policy.check"
	auditPrefix        = "evt.so.action."
	queueGroup         = "policy"
)

// Config wires the service. Every dependency is an interface so tests run in-process.
type Config struct {
	Conn      *nats.Conn
	JS        nats.JetStreamContext
	Schemas   *Schemas
	Authz     Authorizer
	Entities  EntityResolver
	Zones     ZoneSource
	Directory *Directory
	Signer    *Signer // nil: commands go out unsigned (bridges in sig_required:false, logged)
	Log       *slog.Logger
	Now       func() time.Time

	Source            string        // envelope source, default service:policy
	ConfirmWindow     time.Duration // default 120s (asset-command.md step 3)
	CommandTTL        time.Duration // default 30s (expires_at)
	AckTimeout        time.Duration // default 10s (no_ack)
	TaskCreateTimeout time.Duration // default 5s (wait for evt.mission.{tsk}.created)
	SweepEvery        time.Duration // default 500ms
	CheckTimeout      time.Duration // default 2s per OpenFGA call
	Retain            time.Duration // default 24h: finished actions are dropped from memory after this
	ResultTimeout     time.Duration // default 30s: entity/mission/task command with no evt.entity/evt.mission result fails
	ExecuteTimeout    time.Duration // default 10m: acknowledged asset command with no executed/failed outcome fails
	MaxInFlight       int           // default 128 requests handled concurrently
	DenialAuditRate   int           // default 60: denial and refused-confirm audits per identity per minute
}

// entry holds one action. op serializes the operations that change it (confirm, dispatch,
// outcome, result, sweep) and may be held across I/O; it never blocks another action.
// cur is the committed state: read and replaced only under Service.mu, which is never held
// across I/O, so lookups never wait for an operation in progress.
type entry struct {
	op  sync.Mutex
	cur *Action
}

// Service is the policy and audit service. It is the only publisher of cmd.* (besides
// cmd.map.>) and the single writer of evt.so.action.*.
type Service struct {
	cfg   Config
	mu    sync.Mutex // guards byID, byRQ, byRes and every entry.cur; held for map and pointer work only
	byID  map[string]*entry
	byRQ  map[string]string // requester id + request_id -> action_id
	byRes map[string]string // command id, or "tsk:"+task id, awaiting an entity/task result -> action_id

	waitMu  sync.Mutex
	created map[string]chan struct{} // tsk id -> closed when evt.mission.{tsk}.created arrives

	limiter  *auditLimiter
	inflight chan struct{}
	subs     []*nats.Subscription
	cancel   context.CancelFunc
	done     chan struct{}
}

// New validates the config and fills defaults.
func New(cfg Config) (*Service, error) {
	if cfg.Conn == nil || cfg.JS == nil || cfg.Schemas == nil || cfg.Authz == nil || cfg.Entities == nil ||
		cfg.Zones == nil || cfg.Directory == nil {
		return nil, errors.New("policy: incomplete config")
	}
	def := func(d *time.Duration, v time.Duration) {
		if *d == 0 {
			*d = v
		}
	}
	def(&cfg.ConfirmWindow, 120*time.Second)
	def(&cfg.CommandTTL, 30*time.Second)
	def(&cfg.AckTimeout, 10*time.Second)
	def(&cfg.TaskCreateTimeout, 5*time.Second)
	def(&cfg.SweepEvery, 500*time.Millisecond)
	def(&cfg.CheckTimeout, 2*time.Second)
	def(&cfg.Retain, 24*time.Hour)
	def(&cfg.ResultTimeout, 30*time.Second)
	def(&cfg.ExecuteTimeout, 10*time.Minute)
	if cfg.MaxInFlight <= 0 {
		cfg.MaxInFlight = 128
	}
	if cfg.DenialAuditRate <= 0 {
		cfg.DenialAuditRate = 60
	}
	if cfg.Source == "" {
		cfg.Source = "service:policy"
	}
	if cfg.Now == nil {
		cfg.Now = time.Now
	}
	if cfg.Log == nil {
		cfg.Log = slog.Default()
	}
	return &Service{cfg: cfg, byID: map[string]*entry{}, byRQ: map[string]string{}, byRes: map[string]string{},
		created: map[string]chan struct{}{}, limiter: newAuditLimiter(cfg.DenialAuditRate, time.Minute),
		inflight: make(chan struct{}, cfg.MaxInFlight)}, nil
}

// lookup returns the entry of an action, or nil.
func (s *Service) lookup(id string) *entry {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.byID[id]
}

// load returns a private working copy of the committed state. The caller holds e.op when it
// will change the copy and save it.
func (s *Service) load(e *entry) *Action {
	s.mu.Lock()
	defer s.mu.Unlock()
	return e.cur.clone()
}

// save commits a working copy and keeps the result index in step. Caller holds e.op.
func (s *Service) save(e *entry, w *Action) {
	s.mu.Lock()
	defer s.mu.Unlock()
	e.cur = w.clone()
	for _, k := range resultKeys(w) {
		if w.awaitingResult() {
			s.byRes[k] = w.ActionID
		} else if s.byRes[k] == w.ActionID {
			delete(s.byRes, k)
		}
	}
}

// remove drops an action from every map. Caller holds s.mu.
func (s *Service) remove(a *Action) {
	delete(s.byID, a.ActionID)
	if s.byRQ[a.rqKey] == a.ActionID {
		delete(s.byRQ, a.rqKey)
	}
	for _, k := range resultKeys(a) {
		if s.byRes[k] == a.ActionID {
			delete(s.byRes, k)
		}
	}
}

// resultKeys are the index keys under which an entity/task result for the action arrives:
// the command id (the result's causation) and, for a task assignment to a person, team or
// agent, the task id (evt.mission.{tsk}.created).
func resultKeys(a *Action) []string {
	var ks []string
	if a.CommandID != "" {
		ks = append(ks, a.CommandID)
	}
	if a.ActionType == "task.assign" && a.TaskID != "" && !isAssetClass(a.Target.EntityType) {
		ks = append(ks, "tsk:"+a.TaskID)
	}
	return ks
}

func rqKey(requester, requestID string) string { return requester + "\x00" + requestID }

// Start subscribes to every endpoint and fact, and starts the expiry sweeper.
func (s *Service) Start() error {
	type h struct {
		subj  string
		fn    nats.MsgHandler
		queue bool
	}
	// Requests run concurrently (bounded by MaxInFlight): a slow authorization, entity lookup
	// or task wait on one request never holds up another, so a stop is never queued behind it.
	serve := func(fn func(*nats.Msg) any) nats.MsgHandler {
		return func(m *nats.Msg) {
			s.inflight <- struct{}{}
			go func() {
				defer func() { <-s.inflight }()
				s.respond(m, fn(m))
			}()
		}
	}
	hs := []h{
		{SubjCommandPropose, serve(func(m *nats.Msg) any { return s.handlePropose(m, FamilyCommand) }), true},
		{SubjTaskPropose, serve(func(m *nats.Msg) any { return s.handlePropose(m, FamilyTask) }), true},
		{SubjCommandConfirm, serve(func(m *nats.Msg) any { return s.handleConfirm(m, FamilyCommand) }), true},
		{SubjTaskConfirm, serve(func(m *nats.Msg) any { return s.handleConfirm(m, FamilyTask) }), true},
		{SubjActionGet, serve(s.handleGet), true},
		{SubjCheck, serve(s.handleCheck), true},
		{"evt.drone.*.command", s.handleOutcome, false},
		{"evt.robot.*.command", s.handleOutcome, false},
		{"evt.vehicle.*.command", s.handleOutcome, false},
		// Separate subscriptions: the task-created signal must never queue behind a result.
		{"evt.mission.*.created", s.handleCreated, false},
		{"evt.mission.>", s.handleResult, false},
		{"evt.entity.>", s.handleResult, false},
		{"evt.map.geofence.*", s.handleResult, false}, // zone command results (D1-P2)
	}
	for _, x := range hs {
		var sub *nats.Subscription
		var err error
		if x.queue {
			sub, err = s.cfg.Conn.QueueSubscribe(x.subj, queueGroup, x.fn)
		} else {
			sub, err = s.cfg.Conn.Subscribe(x.subj, x.fn)
		}
		if err != nil {
			s.Stop()
			return fmt.Errorf("subscribe %s: %w", x.subj, err)
		}
		s.subs = append(s.subs, sub)
	}
	if err := s.cfg.Conn.Flush(); err != nil {
		return err
	}
	ctx, cancel := context.WithCancel(context.Background())
	s.cancel, s.done = cancel, make(chan struct{})
	go s.sweep(ctx)
	return nil
}

// Stop unsubscribes and stops the sweeper.
func (s *Service) Stop() {
	for _, sub := range s.subs {
		_ = sub.Unsubscribe()
	}
	s.subs = nil
	// Wait for requests in flight: take every slot, then give them back.
	for i := 0; i < cap(s.inflight); i++ {
		s.inflight <- struct{}{}
	}
	for i := 0; i < cap(s.inflight); i++ {
		<-s.inflight
	}
	if s.cancel != nil {
		s.cancel()
		<-s.done
		s.cancel = nil
	}
}

func (s *Service) respond(m *nats.Msg, v any) {
	if m.Reply == "" {
		return
	}
	b, _ := json.Marshal(v)
	if err := m.Respond(b); err != nil {
		s.cfg.Log.Warn("reply failed", "err", err)
	}
}

func deny(code, msg string) *Denial { return &Denial{Code: code, Message: msg} }

func (s *Service) ctx() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), s.cfg.CheckTimeout)
}

// ---------------------------------------------------------------------------
// Propose

func (s *Service) handlePropose(m *nats.Msg, fam Family) ProposeReply {
	var req ActionRequest
	if err := json.Unmarshal(m.Data, &req); err != nil {
		return s.denyNew(&ActionRequest{}, Principal{}, fam, deny(CodeSchemaInvalid, "request is not valid JSON"))
	}
	p, known := s.cfg.Directory.Resolve(m)
	return s.propose(&req, p, known, fam)
}

// propose runs every check of asset-command.md section 4 steps 1-3. It holds no lock across
// the authorization, entity and geofence calls; s.mu is taken only to read or change maps.
func (s *Service) propose(req *ActionRequest, p Principal, known bool, fam Family) ProposeReply {
	// Identity: the verified NATS identity replaces the claimed one.
	verified := p.Verified
	if verified {
		if !known {
			return s.denyNew(req, p, fam, deny(CodeNotAuthorized, "unknown NATS identity "+p.Account+"/"+p.User))
		}
		claimed := req.RequestedBy
		req.RequestedBy = Actor{Kind: p.Kind, ID: p.ID}
		if p.Kind == "agent" {
			req.RequestedBy.OnBehalfOf = claimed.OnBehalfOf
		}
	}
	if d := s.validateEnvelope(req); d != nil {
		return s.denyNew(req, p, fam, d)
	}
	if hit, word := StopLine(req.ActionType); hit {
		return s.denyNew(req, p, fam, deny(CodeForbiddenAction, "action crosses the civilian stop line ("+word+")"))
	}
	spec, ok := Catalogue[req.ActionType]
	if !ok {
		return s.denyNew(req, p, fam, deny(CodeUnknownAction, "action type is not in the catalogue"))
	}
	if spec.Level == A5 {
		return s.denyNew(req, p, fam, deny(CodeForbiddenAction, "A5 actions are never exposed through OneMind actions"))
	}
	if spec.Family != fam {
		return s.denyNew(req, p, fam, deny(CodeConflict, fmt.Sprintf("%s is proposed through svc.policy.%s.propose", spec.Name, spec.Family)))
	}
	if !spec.supports(req.ActionVersion) {
		return s.denyNew(req, p, fam, deny(CodeSchemaInvalid, fmt.Sprintf("unsupported action_version %d", req.ActionVersion)))
	}
	if !spec.allowsTarget(req.Target.EntityType) {
		return s.denyNew(req, p, fam, deny(CodeSchemaInvalid, "target type "+req.Target.EntityType+" is not valid for "+spec.Name))
	}
	if !validToken(req.Target.EntityID) || !validToken(req.Target.EntityType) {
		return s.denyNew(req, p, fam, deny(CodeSchemaInvalid, "target id or type is not a valid subject token"))
	}
	isAsset := isAssetClass(req.Target.EntityType)
	if isAsset && !strings.HasPrefix(req.Target.EntityID, "sim-") {
		return s.denyNew(req, p, fam, deny(CodeNotSimulated, "only simulated assets (sim-*) can be commanded in this run"))
	}
	if len(req.Params) == 0 || string(req.Params) == "null" {
		req.Params = json.RawMessage(`{}`)
	}
	if err := s.cfg.Schemas.ValidateAction(req.ActionType, req.ActionVersion, req.Params); err != nil {
		return s.denyNew(req, p, fam, deny(CodeSchemaInvalid, "params: "+firstLine(err)))
	}
	if scanTK(req.Params) {
		return s.denyNew(req, p, fam, deny(CodeForbiddenAction, "CoT t-k (strike tasking) is never accepted"))
	}

	level, verb, capability, geoPoint, d := s.resolveAction(spec, req)
	if d != nil {
		return s.denyNew(req, p, fam, d)
	}

	hash, err := ParamsHash(req.ActionType, req.ActionVersion, req.Target, req.Params)
	if err != nil {
		return s.denyNew(req, p, fam, deny(CodeSchemaInvalid, "params cannot be canonicalized"))
	}
	// Idempotency is per requester: another principal reusing a request_id gets its own
	// checks, never the first requester's reply.
	key := rqKey(req.RequestedBy.ID, req.RequestID)
	if !req.DryRun {
		if r, d, ok := s.replay(key, req.ActionType, hash); ok {
			if d != nil {
				return s.denyNew(req, p, fam, d)
			}
			return r
		}
	}

	// Person references: never orbit, observe or task against a person (stop line SL2).
	if d := s.checkReferences(spec, req); d != nil {
		return s.denyNew(req, p, fam, d)
	}
	if isAsset && !spec.Unblockable {
		if d := s.checkAsset(req.Target, capability, geoPoint, fenced(spec, verb), spec.Geofence == "self"); d != nil {
			return s.denyNew(req, p, fam, d)
		}
	}

	pol, d := s.authorizePropose(spec, req, p, isAsset)
	if d != nil {
		return s.denyNewWithPolicy(req, p, fam, d, pol, level)
	}

	needed, _ := ApprovalsNeeded(level)
	now := s.cfg.Now()
	if req.DryRun {
		return ProposeReply{Status: StateAwaiting, Level: level, ParamsHash: hash, ApprovalsNeeded: needed, DryRun: true}
	}
	a := &Action{
		ActionID: "act-" + NewULID(now), Family: fam, RequestID: req.RequestID, ActionType: req.ActionType,
		ActionVersion: req.ActionVersion, Level: level, Target: req.Target, Params: req.Params, ParamsHash: hash,
		RequestedBy: req.RequestedBy, ActorVerified: verified, Reason: req.Reason, State: StateAwaiting,
		ApprovalsNeeded: needed, Approvals: []Approval{}, ConfirmBy: now.Add(s.cfg.ConfirmWindow),
		CreatedAt: now, UpdatedAt: now, TaskID: req.TaskID, Policy: pol, Corr: req.Corr,
		Simulated: isAsset, Agent: req, rqKey: key,
	}
	// Reserve the request_id with the operation lock held, so a concurrent duplicate gets
	// this action and nothing can act on it before the proposed event is recorded.
	e := &entry{cur: a.clone()}
	e.op.Lock()
	defer e.op.Unlock()
	s.mu.Lock()
	if _, taken := s.byRQ[key]; taken {
		s.mu.Unlock()
		if r, d, ok := s.replay(key, req.ActionType, hash); ok {
			if d != nil {
				return s.denyNew(req, p, fam, d)
			}
			return r
		}
		return s.denyNew(req, p, fam, deny(CodeConflict, "request_id is in use"))
	}
	s.byID[a.ActionID], s.byRQ[key] = e, a.ActionID
	s.mu.Unlock()
	if err := s.audit(a, "proposed", nil, ""); err != nil {
		s.mu.Lock()
		s.remove(a)
		s.mu.Unlock()
		return s.denyNewWithPolicy(req, p, fam, deny(CodePolicyUnavailable, "audit trail unavailable: "+err.Error()), pol, level)
	}
	if needed == 0 { // A0/A1: no gate, logged
		s.dispatch(a, "")
	}
	s.save(e, a)
	return replyFor(a)
}

// replay answers a repeated request_id from the same requester: the first action's reply
// when the body matches, a conflict denial when it does not. ok is false for a new key.
// A duplicate that arrives while the first is still being proposed waits for it.
func (s *Service) replay(key, actionType, hash string) (ProposeReply, *Denial, bool) {
	s.mu.Lock()
	id, ok := s.byRQ[key]
	e := s.byID[id]
	s.mu.Unlock()
	if !ok || e == nil {
		return ProposeReply{}, nil, false
	}
	e.op.Lock() // the first proposal holds it until its proposed event is recorded
	prev := s.load(e)
	e.op.Unlock()
	if s.lookup(id) == nil { // its audit failed and it was withdrawn
		return ProposeReply{}, nil, false
	}
	if prev.ParamsHash == hash && prev.ActionType == actionType {
		return replyFor(prev), nil, true
	}
	return ProposeReply{}, deny(CodeConflict, "request_id was already used with different parameters"), true
}

// fenced reports whether the action moves the asset, so an unevaluable zone denies it.
func fenced(spec ActionSpec, verb string) bool {
	if spec.Geofence != "" {
		return true
	}
	switch verb {
	case "goto", "land", "orbit", "takeoff":
		return spec.Exec == execTask
	}
	return false
}

func replyFor(a *Action) ProposeReply {
	r := ProposeReply{ActionID: a.ActionID, Status: a.State, Level: a.Level, ParamsHash: a.ParamsHash,
		ApprovalsNeeded: a.ApprovalsNeeded, Denial: a.Denial}
	if !a.ConfirmBy.IsZero() {
		r.ConfirmBy = ts(a.ConfirmBy)
	}
	if a.ApprovalsNeeded == 0 && (a.State == StateDispatched || a.State == StateAcked) {
		r.Status = StateExecuted
	}
	return r
}

func firstLine(err error) string {
	s := err.Error()
	if i := strings.IndexByte(s, '\n'); i > 0 {
		s = s[:i]
	}
	if len(s) > 300 {
		s = s[:300]
	}
	return s
}

func (s *Service) validateEnvelope(req *ActionRequest) *Denial {
	switch {
	case req.V != 1:
		return deny(CodeSchemaInvalid, "v must be 1")
	case !uuidRE.MatchString(req.RequestID):
		return deny(CodeSchemaInvalid, "request_id must be a lowercase UUID")
	case req.ActionType == "" || len(req.ActionType) > 100:
		return deny(CodeSchemaInvalid, "action_type is required")
	case req.ActionVersion < 1:
		return deny(CodeSchemaInvalid, "action_version is required")
	case req.Target.EntityType == "" || req.Target.EntityID == "":
		return deny(CodeSchemaInvalid, "target is required")
	case req.RequestedBy.ID == "" || (req.RequestedBy.Kind != "user" && req.RequestedBy.Kind != "agent" && req.RequestedBy.Kind != "system"):
		return deny(CodeSchemaInvalid, "requested_by {kind: user|agent|system, id} is required")
	case len(req.Reason) > 500:
		return deny(CodeSchemaInvalid, "reason is longer than 500 characters")
	case req.Corr != "" && !uuidRE.MatchString(req.Corr):
		return deny(CodeSchemaInvalid, "corr must be a UUID")
	}
	return nil
}

type taskParams struct {
	TaskID   string          `json:"task_id"`
	ParentID *string         `json:"parent_id"`
	Type     string          `json:"type"`
	Spec     json.RawMessage `json:"spec"`
}

type geoPoint struct {
	Lat, Lon float64
	Set      bool
}

// resolveAction decides level, fleet verb, capability and geofence point for the request.
func (s *Service) resolveAction(spec ActionSpec, req *ActionRequest) (Level, string, string, geoPoint, *Denial) {
	isAsset := isAssetClass(req.Target.EntityType)
	switch spec.Exec {
	case execTask:
		var tp taskParams
		_ = json.Unmarshal(req.Params, &tp)
		if hit, word := StopLine(tp.Type); hit {
			return "", "", "", geoPoint{}, deny(CodeForbiddenAction, "task type crosses the civilian stop line ("+word+")")
		}
		verb, ok := TaskTypes[tp.Type]
		if !ok || !s.cfg.Schemas.HasTask(tp.Type) {
			return "", "", "", geoPoint{}, deny(CodeUnknownAction, "task type "+tp.Type+" is not supported (composite types are decomposed by the task service)")
		}
		if err := s.cfg.Schemas.ValidateTask(tp.Type, tp.Spec); err != nil {
			return "", "", "", geoPoint{}, deny(CodeSchemaInvalid, "spec: "+firstLine(err))
		}
		if isAsset != (verb != "") {
			return "", "", "", geoPoint{}, deny(CodeSchemaInvalid, "task type "+tp.Type+" cannot be assigned to a "+req.Target.EntityType)
		}
		if scanTK(tp.Spec) {
			return "", "", "", geoPoint{}, deny(CodeForbiddenAction, "CoT t-k (strike tasking) is never accepted")
		}
		if !isAsset {
			return A2, "", "", geoPoint{}, nil
		}
		return A3, verb, tp.Type, pointFrom(verb, tp.Spec), nil
	case execAbort:
		if isAsset {
			return A3, "hold", "", geoPoint{}, nil
		}
		return A2, "", "", geoPoint{}, nil
	case execFleet:
		var gp geoPoint
		if spec.Geofence == "point" || spec.Geofence == "center" {
			gp = pointFrom(spec.Verb, req.Params)
		}
		return spec.Level, spec.Verb, spec.Capability, gp, nil
	}
	return spec.Level, spec.Verb, "", geoPoint{}, nil
}

// pointFrom extracts the point a verb moves to (goto, land: lat/lon; orbit: center).
func pointFrom(verb string, raw json.RawMessage) geoPoint {
	var p struct {
		Lat    *float64 `json:"lat"`
		Lon    *float64 `json:"lon"`
		Center *struct {
			Lat *float64 `json:"lat"`
			Lon *float64 `json:"lon"`
		} `json:"center"`
	}
	_ = json.Unmarshal(raw, &p)
	switch verb {
	case "goto", "land":
		if p.Lat != nil && p.Lon != nil {
			return geoPoint{*p.Lat, *p.Lon, true}
		}
	case "orbit":
		if p.Center != nil && p.Center.Lat != nil && p.Center.Lon != nil {
			return geoPoint{*p.Center.Lat, *p.Center.Lon, true}
		}
	}
	return geoPoint{}
}

// referencedEntities lists entity ids the action points at (orbit center, observe target,
// task spec target).
func referencedEntities(spec ActionSpec, req *ActionRequest) []string {
	raw := req.Params
	if spec.Exec == execTask {
		var tp taskParams
		_ = json.Unmarshal(req.Params, &tp)
		raw = tp.Spec
	}
	var p struct {
		Center *struct {
			EntityID string `json:"entity_id"`
		} `json:"center"`
		Target *struct {
			EntityID string `json:"entity_id"`
		} `json:"target"`
		TargetEntityID string `json:"target_entity_id"`
	}
	_ = json.Unmarshal(raw, &p)
	var out []string
	if p.Center != nil && p.Center.EntityID != "" {
		out = append(out, p.Center.EntityID)
	}
	if p.Target != nil && p.Target.EntityID != "" {
		out = append(out, p.Target.EntityID)
	}
	if p.TargetEntityID != "" {
		out = append(out, p.TargetEntityID)
	}
	return out
}

func (s *Service) checkReferences(spec ActionSpec, req *ActionRequest) *Denial {
	for _, id := range referencedEntities(spec, req) {
		ctx, cancel := s.ctx()
		e, err := s.cfg.Entities.Get(ctx, "", id)
		cancel()
		switch {
		case errors.Is(err, ErrEntityNotFound):
			return &Denial{Code: CodePreconditionFailed, Message: "referenced entity " + id + " is unknown", Detail: "unknown_entity"}
		case err != nil:
			return deny(CodePolicyUnavailable, "entity service unavailable: "+err.Error())
		case e.EntityType == "person" || strings.EqualFold(e.Kind, "person"):
			return deny(CodeForbiddenAction, "actions never target a person or a person-class track")
		}
	}
	return nil
}

// checkAsset runs the preconditions on the target asset: it exists, is simulated, is
// capable, is online, and the destination passes the geofence. A movement verb (fenced) is
// denied while any zone that applies to the asset cannot be evaluated.
func (s *Service) checkAsset(t Target, capability string, gp geoPoint, fenced, selfGeofence bool) *Denial {
	ctx, cancel := s.ctx()
	e, err := s.cfg.Entities.Get(ctx, t.EntityType, t.EntityID)
	cancel()
	if errors.Is(err, ErrEntityNotFound) {
		return &Denial{Code: CodePreconditionFailed, Message: "asset " + t.EntityID + " is unknown", Detail: "unknown_entity"}
	}
	if err != nil {
		return deny(CodePolicyUnavailable, "entity service unavailable: "+err.Error())
	}
	if e.Indicators != nil && e.Indicators.Simulated != nil && !*e.Indicators.Simulated {
		return deny(CodeNotSimulated, "asset record is not marked simulated")
	}
	if capability != "" {
		found := false
		for _, c := range e.Capabilities {
			if c == capability {
				found = true
			}
		}
		if !found {
			return deny(CodeNotCapable, "asset does not advertise "+capability)
		}
	}
	if e.Health == nil || e.Health.Connection != "online" {
		// No queued replay: a command is never held for an offline asset.
		return &Denial{Code: CodePreconditionFailed, Message: "asset is not online", Detail: "offline"}
	}
	if selfGeofence && e.Location != nil {
		gp = geoPoint{e.Location.Lat, e.Location.Lon, true}
	}
	if fenced || gp.Set {
		zones, err := s.cfg.Zones.Zones()
		if err != nil {
			return deny(CodePolicyUnavailable, "geofences unavailable: "+err.Error())
		}
		if z, bad := UnevaluableZone(zones, t.EntityID); bad {
			return &Denial{Code: CodePreconditionFailed, Detail: DetailUnevaluable,
				Message: "geofence " + z.ID + " cannot be evaluated (" + z.Problem + "); movement is denied until it is fixed"}
		}
		if !gp.Set {
			return nil
		}
		if v := CheckPoint(zones, t.EntityID, gp.Lat, gp.Lon); v != "" {
			return &Denial{Code: CodePreconditionFailed, Message: "point violates a geofence", Detail: v}
		}
	}
	return nil
}

func (s *Service) checkFGA(pol *PolicyInfo, user, relation, object string) (bool, error) {
	ctx, cancel := s.ctx()
	defer cancel()
	now := s.cfg.Now().UTC().Format(time.RFC3339Nano)
	ok, err := s.cfg.Authz.Check(ctx, CheckRequest{User: user, Relation: relation, Object: object,
		Context: map[string]any{"current_time": now}})
	if err != nil {
		return false, err
	}
	pol.Checks = append(pol.Checks, PolicyCheck{User: user, Relation: relation, Object: object, Allowed: ok})
	return ok, nil
}

// authorizePropose: can_propose for the requester; for an agent, the intersection with the
// on-behalf-of human (palantir-patterns.md section 7); operator on the asset.
func (s *Service) authorizePropose(spec ActionSpec, req *ActionRequest, p Principal, isAsset bool) (PolicyInfo, *Denial) {
	pol := PolicyInfo{Store: s.cfg.Authz.StoreID(), ModelID: s.cfg.Authz.ModelID(), Checks: []PolicyCheck{}}
	who := req.RequestedBy
	actionObj := "action_type:" + spec.Name
	var users []string
	human := ""
	switch who.Kind {
	case "agent":
		if !strings.HasPrefix(who.OnBehalfOf, "user:") {
			return pol, deny(CodeNotAuthorized, "an agent proposal needs on_behalf_of a human user")
		}
		users = []string{who.ID, who.OnBehalfOf}
		human = who.OnBehalfOf
	case "user":
		users = []string{who.ID}
		human = who.ID
	default:
		users = []string{who.ID}
	}
	for _, u := range users {
		ok, err := s.checkFGA(&pol, u, "proposer", actionObj)
		if err != nil {
			return pol, deny(CodePolicyUnavailable, "authorization unavailable (fail closed): "+err.Error())
		}
		if !ok {
			return pol, deny(CodeNotAuthorized, u+" may not propose "+spec.Name)
		}
	}
	if isAsset && human != "" {
		ok, err := s.checkFGA(&pol, human, "operator", "asset:"+req.Target.EntityID)
		if err != nil {
			return pol, deny(CodePolicyUnavailable, "authorization unavailable (fail closed): "+err.Error())
		}
		if !ok {
			return pol, deny(CodeNotAuthorized, human+" does not operate "+req.Target.EntityID)
		}
	}
	return pol, nil
}

// denyNew audits a denied proposal and builds the reply.
func (s *Service) denyNew(req *ActionRequest, p Principal, fam Family, d *Denial) ProposeReply {
	return s.denyNewWithPolicy(req, p, fam, d, PolicyInfo{Store: s.cfg.Authz.StoreID(), ModelID: s.cfg.Authz.ModelID(), Checks: []PolicyCheck{}}, "")
}

func (s *Service) denyNewWithPolicy(req *ActionRequest, p Principal, fam Family, d *Denial, pol PolicyInfo, level Level) ProposeReply {
	now := s.cfg.Now()
	if level == "" {
		if spec, ok := Catalogue[req.ActionType]; ok {
			level = spec.Level
		} else {
			level = A0
		}
	}
	a := &Action{ActionID: "act-" + NewULID(now), Family: fam, RequestID: req.RequestID, ActionType: req.ActionType,
		ActionVersion: req.ActionVersion, Level: level, Target: req.Target, Params: req.Params,
		RequestedBy: req.RequestedBy, ActorVerified: p.Verified, State: StateDenied, Approvals: []Approval{},
		CreatedAt: now, UpdatedAt: now, Denial: d, Policy: pol, Simulated: isAssetClass(req.Target.EntityType),
		Agent: req}
	if a.ActionType == "" {
		a.ActionType = "unknown"
	}
	if h, err := ParamsHash(a.ActionType, a.ActionVersion, a.Target, a.Params); err == nil {
		a.ParamsHash = h
	} else {
		a.ParamsHash = "sha256:invalid"
	}
	if req.DryRun {
		return ProposeReply{ActionID: "", Status: StateDenied, Denial: d, DryRun: true}
	}
	// A denial is not kept in memory: the audit event is its record. The audit is rate
	// limited per identity so a flood of bad requests cannot flood EVT.
	if s.limiter.allow(attemptBy(p), now) {
		if err := s.audit(a, "denied", d, ""); err != nil {
			s.cfg.Log.Error("audit of denial failed", "action_id", a.ActionID, "err", err)
		}
	} else {
		s.cfg.Log.Warn("denial audit rate limited", "action_id", a.ActionID, "by", attemptBy(p), "code", d.Code)
	}
	s.cfg.Log.Info("denied", "action_id", a.ActionID, "action_type", a.ActionType, "code", d.Code, "by", a.RequestedBy.ID)
	return ProposeReply{ActionID: a.ActionID, Status: StateDenied, Level: level, ParamsHash: a.ParamsHash, Denial: d}
}

// ---------------------------------------------------------------------------
// Confirm

func (s *Service) handleConfirm(m *nats.Msg, fam Family) ConfirmReply {
	var req ConfirmRequest
	if err := json.Unmarshal(m.Data, &req); err != nil || req.V != 1 || req.ActionID == "" ||
		!strings.HasPrefix(req.ParamsHash, "sha256:") || (req.Decision != "approve" && req.Decision != "reject") ||
		(req.Decision == "reject" && req.Reason == "") {
		return ConfirmReply{ActionID: req.ActionID, Status: StateDenied,
			Denial: deny(CodeSchemaInvalid, "confirm needs v:1, action_id, params_hash, decision approve|reject, and reason for reject")}
	}
	p, known := s.cfg.Directory.Resolve(m)
	// Only this action's operation lock is held across the checks, the audit and the
	// dispatch (which may wait for the task service); other actions are never blocked.
	var a *Action
	e := s.lookup(req.ActionID)
	if e != nil {
		e.op.Lock()
		defer e.op.Unlock()
		a = s.load(e)
	}
	refuse := func(d *Denial) ConfirmReply {
		r := ConfirmReply{ActionID: req.ActionID, Status: StateDenied, Denial: d}
		if a != nil {
			r.Status, r.ApprovalsNeeded, r.ApprovalsReceived = a.State, a.ApprovalsNeeded, len(a.Approvals)
			// A refused confirmation is audited (rate limited per identity); the action
			// itself keeps its state.
			if s.limiter.allow(attemptBy(p), s.cfg.Now()) {
				if err := s.audit(a, "denied", d, attemptBy(p)); err != nil {
					s.cfg.Log.Error("audit failed", "err", err)
				}
			} else {
				s.cfg.Log.Warn("refused-confirm audit rate limited", "action_id", a.ActionID, "by", attemptBy(p))
			}
		}
		s.cfg.Log.Info("confirm refused", "action_id", req.ActionID, "code", d.Code, "by", attemptBy(p))
		return r
	}
	switch {
	case p.Verified && p.Kind == "agent":
		return refuse(deny(CodeAgentCannotConfirm, "agents may propose but never confirm"))
	case !p.Verified:
		return refuse(deny(CodeNotAuthorized, "confirmation needs an authenticated human identity (service import with share)"))
	case !known || p.Kind != "user":
		return refuse(deny(CodeNotAuthorized, "caller is not a known human user"))
	case a == nil:
		return refuse(deny(CodeConflict, "no such action"))
	case a.Family != fam:
		return refuse(deny(CodeConflict, fmt.Sprintf("this action is confirmed through svc.policy.%s.confirm", a.Family)))
	case a.ParamsHash != req.ParamsHash:
		return refuse(deny(CodeConflict, "params_hash does not match the action that was proposed"))
	}
	switch a.State {
	case StateAwaiting:
	case StateExpired:
		return refuse(deny(CodeExpired, "the confirmation window has closed"))
	case StateApproved, StateDispatched, StateAcked, StateExecuted, StateFailed, StateRefused:
		return refuse(deny(CodeReplayed, "this action was already approved; an approval is used once"))
	default:
		return refuse(deny(CodeConflict, "action is "+a.State))
	}
	now := s.cfg.Now()
	if now.After(a.ConfirmBy) {
		s.expire(a)
		s.save(e, a)
		return ConfirmReply{ActionID: a.ActionID, Status: a.State, ApprovalsNeeded: a.ApprovalsNeeded,
			ApprovalsReceived: len(a.Approvals), Denial: a.Denial}
	}
	for _, ap := range a.Approvals {
		if ap.By == p.ID {
			return refuse(deny(CodeDuplicateConfirmation, "the same person cannot confirm twice"))
		}
	}
	if req.Decision == "reject" {
		a.State, a.UpdatedAt = StateRejected, now
		a.Denial = deny(CodeRejectedByHuman, req.Reason)
		if err := s.audit(a, "rejected", a.Denial, p.ID); err != nil {
			s.cfg.Log.Error("audit failed", "err", err)
		}
		s.save(e, a)
		return ConfirmReply{ActionID: a.ActionID, Status: a.State, ApprovalsNeeded: a.ApprovalsNeeded,
			ApprovalsReceived: len(a.Approvals), Denial: a.Denial}
	}
	step := len(a.Approvals) + 1
	relation := "approver"
	if step == 2 {
		relation = "second_approver"
	}
	ok, err := s.checkFGA(&a.Policy, p.ID, relation, "action_type:"+a.ActionType)
	if err == nil && ok && step == 1 && isAssetClass(a.Target.EntityType) {
		ok, err = s.checkFGA(&a.Policy, p.ID, "operator", "asset:"+a.Target.EntityID)
	}
	if err != nil {
		return refuse(deny(CodePolicyUnavailable, "authorization unavailable (fail closed): "+err.Error()))
	}
	if !ok {
		return refuse(deny(CodeNotAuthorized, p.ID+" may not confirm step "+fmt.Sprint(step)+" of "+a.ActionType))
	}
	a.Approvals = append(a.Approvals, Approval{By: p.ID, At: ts(now), Step: step})
	a.UpdatedAt = now
	if len(a.Approvals) < a.ApprovalsNeeded {
		s.save(e, a)
		return ConfirmReply{ActionID: a.ActionID, Status: a.State, ApprovalsNeeded: a.ApprovalsNeeded,
			ApprovalsReceived: len(a.Approvals)}
	}
	a.State = StateApproved
	approvedID := NewUUID()
	if err := s.auditWithID(a, "approved", nil, "", approvedID); err != nil {
		s.fail(a, deny(CodePolicyUnavailable, "audit trail unavailable; nothing dispatched: "+err.Error()))
	} else {
		s.save(e, a) // readers see approved while the dispatch runs
		s.dispatch(a, approvedID)
	}
	s.save(e, a)
	return ConfirmReply{ActionID: a.ActionID, Status: a.State, ApprovalsNeeded: a.ApprovalsNeeded,
		ApprovalsReceived: len(a.Approvals), CommandID: a.CommandID, Denial: a.Denial}
}

func attemptBy(p Principal) string {
	if p.ID != "" {
		return p.ID
	}
	if p.Verified {
		return "nats:" + p.Account + "/" + p.User
	}
	return "unverified"
}

func (s *Service) expire(a *Action) {
	a.State, a.UpdatedAt = StateExpired, s.cfg.Now()
	a.Denial = deny(CodeExpired, "no confirmation before confirm_by")
	if err := s.audit(a, "expired", a.Denial, ""); err != nil {
		s.cfg.Log.Error("audit failed", "err", err)
	}
}

func (s *Service) fail(a *Action, d *Denial) {
	a.State, a.UpdatedAt, a.Denial = StateFailed, s.cfg.Now(), d
	if err := s.audit(a, "failed", d, ""); err != nil {
		s.cfg.Log.Error("audit failed", "err", err)
	}
}

// ---------------------------------------------------------------------------
// Dispatch (asset-command.md step 5; mission-task.md section 6.1 steps 4-7)

func (s *Service) dispatch(a *Action, causation string) {
	isAsset := isAssetClass(a.Target.EntityType)
	spec := Catalogue[a.ActionType]
	// Re-check preconditions: the world may have changed since the proposal.
	if isAsset && !spec.Unblockable {
		req := a.Agent
		_, verb, capability, gp, d := s.resolveAction(spec, req)
		if d == nil {
			d = s.checkReferences(spec, req)
		}
		if d == nil {
			d = s.checkAsset(a.Target, capability, gp, fenced(spec, verb), spec.Geofence == "self")
		}
		if d != nil {
			if d.Code == CodePolicyUnavailable || d.Code == CodePreconditionFailed {
				a.State, a.UpdatedAt, a.Denial = StateDenied, s.cfg.Now(), d
				if err := s.audit(a, "denied", d, ""); err != nil {
					s.cfg.Log.Error("audit failed", "err", err)
				}
				return
			}
			s.fail(a, d)
			return
		}
	}
	switch spec.Exec {
	case execFleet:
		s.sendFleet(a, spec.Verb, a.Params, "", causation)
	case execEntity, execMission, execMap:
		subj, err := commandSubject(spec.Exec, "", a.Target.EntityID, spec.Verb)
		if err != nil {
			s.fail(a, deny(CodeSchemaInvalid, err.Error()))
			return
		}
		id := NewUUID()
		s.expectResult(a.ActionID, id)
		if err := s.publishCommandID(a, id, subj, a.Params, "", causation); err != nil {
			s.unexpect(a.ActionID, id)
			s.fail(a, deny(CodePolicyUnavailable, err.Error()))
			return
		}
		a.CommandID, a.CommandSubject = id, subj
		s.markDispatched(a, false)
	case execTask:
		var tp taskParams
		_ = json.Unmarshal(a.Params, &tp)
		tsk := tp.TaskID
		if tsk == "" {
			tsk = "tsk-" + NewULID(s.cfg.Now())
		}
		a.TaskID = tsk
		data := map[string]any{}
		_ = json.Unmarshal(a.Params, &data)
		data["task_id"] = tsk
		data["assignee"] = a.Target
		data["author"] = a.RequestedBy
		data["action_id"] = a.ActionID
		data["approval"] = map[string]any{"level": a.Level, "approvals": a.Approvals}
		b, _ := json.Marshal(data)
		createID, createSubj := NewUUID(), "cmd.mission."+tsk+".create"
		verb := TaskTypes[tp.Type]
		if verb == "" {
			// Person, team or agent task: closes on evt.mission.{tsk}.created.
			s.expectResult(a.ActionID, createID, "tsk:"+tsk)
			if err := s.publishCommandID(a, createID, createSubj, b, tsk, causation); err != nil {
				s.unexpect(a.ActionID, createID, "tsk:"+tsk)
				s.fail(a, deny(CodePolicyUnavailable, err.Error()))
				return
			}
			a.CommandID, a.CommandSubject = createID, createSubj
			s.markDispatched(a, false)
			return
		}
		ch := s.awaitCreated(tsk)
		if err := s.publishCommandID(a, createID, createSubj, b, tsk, causation); err != nil {
			s.dropWait(tsk)
			s.fail(a, deny(CodePolicyUnavailable, err.Error()))
			return
		}
		select {
		case <-ch:
		case <-time.After(s.cfg.TaskCreateTimeout):
			s.dropWait(tsk)
			s.fail(a, deny(CodeNoAck, "task service did not record "+tsk+"; no asset command sent"))
			return
		}
		s.sendFleet(a, verb, tp.Spec, tsk, causation)
	case execAbort:
		var ap struct {
			TaskID string `json:"task_id"`
			Reason string `json:"reason"`
		}
		_ = json.Unmarshal(a.Params, &ap)
		a.TaskID = ap.TaskID
		abortSubj, id := "cmd.mission."+ap.TaskID+".abort", NewUUID()
		if !isAsset {
			s.expectResult(a.ActionID, id)
		}
		if err := s.publishCommandID(a, id, abortSubj, a.Params, ap.TaskID, causation); err != nil {
			s.unexpect(a.ActionID, id)
			s.fail(a, deny(CodePolicyUnavailable, err.Error()))
			return
		}
		if isAsset {
			s.sendFleet(a, "hold", json.RawMessage(`{"reason":"abort"}`), ap.TaskID, causation)
			return
		}
		a.CommandID, a.CommandSubject = id, abortSubj
		s.markDispatched(a, false)
	default:
		s.fail(a, deny(CodeForbiddenAction, "action cannot execute"))
	}
}

func (s *Service) sendFleet(a *Action, verb string, data json.RawMessage, taskID, causation string) {
	subj, err := commandSubject(execFleet, a.Target.EntityType, a.Target.EntityID, verb)
	if err != nil {
		s.fail(a, deny(CodeSchemaInvalid, err.Error()))
		return
	}
	id, err := s.publishCommand(a, subj, data, taskID, causation)
	if err != nil {
		s.fail(a, deny(CodePolicyUnavailable, err.Error()))
		return
	}
	a.CommandID, a.CommandSubject = id, subj
	s.markDispatched(a, true)
}

// markDispatched records the dispatch. track: an asset command, closed by the bridge
// outcomes (no_ack after AckTimeout). Otherwise an entity, mission or task command, closed by
// the entity or task service result (asset-command.md section 2), failing after ResultTimeout.
func (s *Service) markDispatched(a *Action, track bool) {
	a.State, a.UpdatedAt = StateDispatched, s.cfg.Now()
	if track {
		a.ackDeadline = s.cfg.Now().Add(s.cfg.AckTimeout)
	} else {
		a.resultDeadline = s.cfg.Now().Add(s.cfg.ResultTimeout)
	}
	if err := s.audit(a, "dispatched", nil, ""); err != nil {
		s.cfg.Log.Error("audit of dispatch failed", "action_id", a.ActionID, "err", err)
	}
}

// publishCommand builds, validates, signs and publishes one command with a JetStream ack
// (the CMD PubAck is the proof of record). It returns the command id.
func (s *Service) publishCommand(a *Action, subject string, data json.RawMessage, taskID, causation string) (string, error) {
	id := NewUUID()
	return id, s.publishCommandID(a, id, subject, data, taskID, causation)
}

// expectResult indexes an entity/task command before it is published, so a result that
// arrives before the dispatch is saved is not lost: it waits on the action's operation lock.
func (s *Service) expectResult(actionID string, keys ...string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, k := range keys {
		s.byRes[k] = actionID
	}
}

func (s *Service) unexpect(actionID string, keys ...string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, k := range keys {
		if s.byRes[k] == actionID {
			delete(s.byRes, k)
		}
	}
}

// publishCommandID publishes one command with the given id (see publishCommand).
func (s *Service) publishCommandID(a *Action, id, subject string, data json.RawMessage, taskID, causation string) error {
	now := s.cfg.Now()
	var d any
	if err := json.Unmarshal(data, &d); err != nil || d == nil {
		d = map[string]any{}
	}
	approver := ""
	if n := len(a.Approvals); n > 0 {
		approver = a.Approvals[n-1].By
	}
	status := "approved"
	if a.ApprovalsNeeded == 0 {
		status = "not_required"
	}
	// approval.params_hash is the hash of the approved action itself (asset-command.md
	// sections 4-5, 6.3; mission-task.md section 6 step 6): asset.{verb} over the command
	// data for an asset action, but task.assign (or task.abort) over the task parameters for
	// a task-bearing command (task_id set). Bridges recompute it only when task_id is absent;
	// with task_id they cannot, and rely on the signature (decision HASH, 2026-10-29).
	env := map[string]any{
		"v": 1, "id": id, "ts": ts(now), "cmd": subject, "source": s.cfg.Source, "subject": subject, "data": d,
		"action_id":  a.ActionID,
		"expires_at": ts(now.Add(s.cfg.CommandTTL)),
		"simulated":  a.Simulated,
		"approval": map[string]any{"required": a.ApprovalsNeeded > 0, "status": status, "approver": approver,
			"reason": a.Reason, "level": a.Level, "params_hash": a.ParamsHash, "approvals": a.Approvals},
		"policy": map[string]any{"decision_id": "dec-" + strings.TrimPrefix(a.ActionID, "act-"), "allowed": true},
	}
	if taskID != "" {
		env["task_id"] = taskID
	}
	if uuidRE.MatchString(a.Corr) {
		env["corr"] = a.Corr
	}
	if uuidRE.MatchString(causation) {
		env["causation"] = causation
	}
	if s.cfg.Signer != nil {
		if err := s.cfg.Signer.Sign(env); err != nil {
			return err
		}
	} else {
		s.cfg.Log.Warn("command sent unsigned (no signing key configured)", "subject", subject)
	}
	b, err := json.Marshal(env)
	if err != nil {
		return err
	}
	if err := s.cfg.Schemas.ValidateCommand(b); err != nil {
		return fmt.Errorf("command envelope invalid: %s", firstLine(err))
	}
	if _, err := s.cfg.JS.Publish(subject, b, nats.MsgId(id)); err != nil {
		return fmt.Errorf("CMD stream did not ack %s: %w", subject, err)
	}
	s.cfg.Log.Info("command published", "subject", subject, "command_id", id, "action_id", a.ActionID)
	return nil
}

func (s *Service) awaitCreated(tsk string) chan struct{} {
	s.waitMu.Lock()
	defer s.waitMu.Unlock()
	ch := make(chan struct{})
	s.created[tsk] = ch
	return ch
}

func (s *Service) dropWait(tsk string) {
	s.waitMu.Lock()
	delete(s.created, tsk)
	s.waitMu.Unlock()
}

func (s *Service) handleCreated(m *nats.Msg) {
	parts := strings.Split(m.Subject, ".")
	if len(parts) != 4 {
		return
	}
	s.waitMu.Lock()
	if ch, ok := s.created[parts[2]]; ok {
		close(ch)
		delete(s.created, parts[2])
	}
	s.waitMu.Unlock()
}

// ---------------------------------------------------------------------------
// Outcomes from bridges (evt.{class}.{id}.command) are republished as evt.so.action.*.

func (s *Service) handleOutcome(m *nats.Msg) {
	parts := strings.Split(m.Subject, ".")
	if len(parts) != 4 {
		return
	}
	var env struct {
		ID   string `json:"id"`
		Data struct {
			ActionID  string          `json:"action_id"`
			CommandID string          `json:"command_id"`
			Outcome   string          `json:"outcome"`
			Code      string          `json:"code"`
			Message   string          `json:"message"`
			MAVLink   json.RawMessage `json:"mavlink"`
		} `json:"data"`
	}
	if json.Unmarshal(m.Data, &env) != nil {
		return
	}
	e := s.lookup(env.Data.ActionID)
	if e == nil {
		s.cfg.Log.Warn("outcome for unknown command dropped", "subject", m.Subject, "action_id", env.Data.ActionID)
		return
	}
	e.op.Lock() // waits only for an operation on this same action (its dispatch)
	defer e.op.Unlock()
	a := s.load(e)
	if a.Target.EntityType != parts[1] || a.Target.EntityID != parts[2] ||
		a.CommandID == "" || env.Data.CommandID != a.CommandID {
		s.cfg.Log.Warn("outcome for mismatched command dropped", "subject", m.Subject, "action_id", env.Data.ActionID)
		return
	}
	if a.terminal() {
		return
	}
	var d *Denial
	now := s.cfg.Now()
	switch env.Data.Outcome {
	case "acknowledged":
		a.State, a.acked = StateAcked, true
		a.execDeadline = now.Add(s.cfg.ExecuteTimeout)
	case "executed":
		a.State, a.acked = StateExecuted, true
	case "failed", "refused":
		a.State, a.acked = env.Data.Outcome, true
		code := env.Data.Code
		if code == "" {
			code = "vehicle_rejected"
		}
		d = &Denial{Code: code, Message: env.Data.Message}
		a.Denial = d
	default:
		return
	}
	a.UpdatedAt = now
	extra := map[string]any{}
	if len(env.Data.MAVLink) > 0 {
		extra["mavlink"] = env.Data.MAVLink
	}
	if err := s.auditFull(a, env.Data.Outcome, d, "", NewUUID(), env.ID, extra); err != nil {
		s.cfg.Log.Error("audit failed", "err", err)
	}
	s.save(e, a)
}

// handleResult closes an entity, mission or task command when its result arrives
// (asset-command.md section 2: the entity and task services report each cmd.entity.* /
// cmd.mission.* they apply as evt.entity.{id}.* / evt.mission.{id}.*; the policy service
// turns it into evt.so.action.executed or .failed). A result matches by its causation (the
// command id), by data.action_id, or, for a task assigned to a person, team or agent, by
// evt.mission.{tsk}.created. Its subject id must be the action's target or task.
func (s *Service) handleResult(m *nats.Msg) {
	parts := strings.Split(m.Subject, ".")
	if len(parts) != 4 {
		return
	}
	s.mu.Lock()
	idle := len(s.byRes) == 0
	s.mu.Unlock()
	if idle {
		return
	}
	var env struct {
		ID        string `json:"id"`
		Causation string `json:"causation"`
		Data      struct {
			ActionID string `json:"action_id"`
			ZoneID   string `json:"zone_id"`
			Code     string `json:"code"`
			Message  string `json:"message"`
		} `json:"data"`
	}
	_ = json.Unmarshal(m.Data, &env)
	s.mu.Lock()
	id := ""
	if env.Causation != "" {
		id = s.byRes[env.Causation]
	}
	if id == "" && env.Data.ActionID != "" {
		if e := s.byID[env.Data.ActionID]; e != nil && e.cur.awaitingResult() {
			id = env.Data.ActionID
		}
	}
	if id == "" && parts[1] == "mission" && parts[3] == "created" {
		id = s.byRes["tsk:"+parts[2]]
	}
	e := s.byID[id]
	s.mu.Unlock()
	if e == nil {
		return
	}
	e.op.Lock()
	defer e.op.Unlock()
	a := s.load(e)
	subjectID := parts[2]
	if parts[1] == "map" && parts[2] == "geofence" {
		// evt.map.geofence.{event}: four tokens, the zone id travels in data.zone_id (D1-P2).
		subjectID = env.Data.ZoneID
	}
	if !a.awaitingResult() || (subjectID != a.Target.EntityID && subjectID != a.TaskID) {
		return
	}
	var d *Denial
	event := "executed"
	switch parts[3] {
	case "failed", "rejected", "refused", "error":
		event = StateFailed
		code := env.Data.Code
		if code == "" {
			code = CodePreconditionFailed
		}
		d = &Denial{Code: code, Message: env.Data.Message, Detail: m.Subject}
		a.Denial = d
	}
	a.State, a.UpdatedAt = event, s.cfg.Now()
	if err := s.auditFull(a, event, d, "", NewUUID(), env.ID, map[string]any{"result_subject": m.Subject}); err != nil {
		s.cfg.Log.Error("audit failed", "err", err)
	}
	s.save(e, a)
}

func (s *Service) sweep(ctx context.Context) {
	defer close(s.done)
	t := time.NewTicker(s.cfg.SweepEvery)
	defer t.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-t.C:
			s.SweepOnce()
		}
	}
}

// SweepOnce expires unconfirmed actions, fails commands with no acknowledgement, no
// execution outcome or no entity/task result, and drops finished actions older than Retain
// (the audit trail in EVT keeps them). Every non-terminal state has a deadline, so every
// action becomes terminal and is dropped: memory stays bounded. An action whose operation
// is in progress is skipped until the next sweep.
func (s *Service) SweepOnce() {
	now := s.cfg.Now()
	s.mu.Lock()
	es := make([]*entry, 0, len(s.byID))
	for _, e := range s.byID {
		es = append(es, e)
	}
	s.mu.Unlock()
	for _, e := range es {
		if !e.op.TryLock() {
			continue
		}
		s.sweepOne(e, now)
		e.op.Unlock()
	}
	if s.limiter != nil {
		s.limiter.prune(now)
	}
}

func (s *Service) sweepOne(e *entry, now time.Time) {
	a := s.load(e)
	switch {
	case a.terminal() && now.Sub(a.UpdatedAt) > s.cfg.Retain:
		s.mu.Lock()
		s.remove(a)
		s.mu.Unlock()
		return
	case a.State == StateAwaiting && now.After(a.ConfirmBy):
		s.expire(a)
	case a.State == StateDispatched && !a.ackDeadline.IsZero() && !a.acked && now.After(a.ackDeadline):
		s.fail(a, deny(CodeNoAck, "no acknowledgement from the asset within "+s.cfg.AckTimeout.String()))
	case a.awaitingResult() && now.After(a.resultDeadline):
		s.fail(a, deny(CodeNoAck, "no result from the entity or task service within "+s.cfg.ResultTimeout.String()+
			"; the command is recorded in CMD"))
	case a.State == StateAcked && !a.execDeadline.IsZero() && now.After(a.execDeadline):
		s.fail(a, deny(CodeNoAck, "acknowledged, but no executed or failed outcome within "+s.cfg.ExecuteTimeout.String()))
	default:
		return
	}
	s.save(e, a)
}

// ---------------------------------------------------------------------------
// Lookups

func (s *Service) handleGet(m *nats.Msg) any {
	var req struct {
		V        int    `json:"v"`
		ActionID string `json:"action_id"`
	}
	if json.Unmarshal(m.Data, &req) != nil || req.ActionID == "" {
		return map[string]any{"action": nil, "denial": deny(CodeSchemaInvalid, "action_id is required")}
	}
	p, known := s.cfg.Directory.Resolve(m)
	if p.Verified && !known {
		return map[string]any{"action": nil, "denial": deny(CodeNotAuthorized, "unknown NATS identity")}
	}
	e := s.lookup(req.ActionID)
	if e == nil {
		return map[string]any{"action": nil, "denial": deny(CodeConflict, "no such action")}
	}
	a := s.load(e) // committed state; never waits for an operation in progress
	if p.Verified && p.Kind == "agent" && a.RequestedBy.ID != p.ID {
		return map[string]any{"action": nil, "denial": deny(CodeNotAuthorized, "agents may read only their own actions")}
	}
	return map[string]any{"action": a, "denial": nil}
}

func (s *Service) handleCheck(m *nats.Msg) any {
	var req struct {
		V        int            `json:"v"`
		User     string         `json:"user"`
		Relation string         `json:"relation"`
		Object   string         `json:"object"`
		Context  map[string]any `json:"context"`
		Action   *ActionRequest `json:"action"`
	}
	if json.Unmarshal(m.Data, &req) != nil {
		return map[string]any{"allowed": false, "denial": deny(CodeSchemaInvalid, "not JSON")}
	}
	p, known := s.cfg.Directory.Resolve(m)
	if req.Action != nil { // full dry run of a proposal: no record, no audit, no command
		req.Action.DryRun = true
		spec, ok := Catalogue[req.Action.ActionType]
		fam := FamilyCommand
		if ok {
			fam = spec.Family
		}
		r := s.propose(req.Action, p, known, fam)
		return map[string]any{"allowed": r.Denial == nil, "decision": r}
	}
	// The tuple pass-through could map the whole permission graph, so it answers only
	// callers in the policy service's own account (no service import: unverified). Clients
	// and agents use the action form, which checks their own identity.
	if p.Verified {
		return map[string]any{"allowed": false, "denial": deny(CodeNotAuthorized,
			"tuple checks are limited to the policy service's own account; send {v, action} for a dry run")}
	}
	if req.User == "" || req.Relation == "" || req.Object == "" {
		return map[string]any{"allowed": false, "denial": deny(CodeSchemaInvalid, "user, relation and object are required")}
	}
	ctx, cancel := s.ctx()
	defer cancel()
	ok, err := s.cfg.Authz.Check(ctx, CheckRequest{User: req.User, Relation: req.Relation, Object: req.Object, Context: req.Context})
	if err != nil {
		return map[string]any{"allowed": false, "denial": deny(CodePolicyUnavailable, "authorization unavailable (fail closed)")}
	}
	return map[string]any{"allowed": ok, "store": s.cfg.Authz.StoreID(), "model_id": s.cfg.Authz.ModelID(), "denial": nil}
}

// Snapshot returns a copy of every action's committed state (tests and diagnostics).
func (s *Service) Snapshot() []Action {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := make([]Action, 0, len(s.byID))
	for _, e := range s.byID {
		out = append(out, *e.cur.clone())
	}
	return out
}

// auditLimiter caps denial and refused-confirm audits per identity in a fixed window.
type auditLimiter struct {
	mu     sync.Mutex
	limit  int
	window time.Duration
	seen   map[string]*limitWindow
}

type limitWindow struct {
	start      time.Time
	n          int
	suppressed int
}

func newAuditLimiter(limit int, window time.Duration) *auditLimiter {
	return &auditLimiter{limit: limit, window: window, seen: map[string]*limitWindow{}}
}

func (l *auditLimiter) allow(who string, now time.Time) bool {
	if l == nil {
		return true
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	w := l.seen[who]
	if w == nil || now.Sub(w.start) >= l.window {
		w = &limitWindow{start: now}
		l.seen[who] = w
	}
	if w.n >= l.limit {
		w.suppressed++
		return false
	}
	w.n++
	return true
}

// prune forgets identities whose window has passed.
func (l *auditLimiter) prune(now time.Time) {
	l.mu.Lock()
	defer l.mu.Unlock()
	for k, w := range l.seen {
		if now.Sub(w.start) >= l.window {
			delete(l.seen, k)
		}
	}
}

// ---------------------------------------------------------------------------
// Audit (single writer of evt.so.action.*)

func (s *Service) audit(a *Action, event string, d *Denial, attempt string) error {
	return s.auditFull(a, event, d, attempt, NewUUID(), "", nil)
}

func (s *Service) auditWithID(a *Action, event string, d *Denial, attempt, id string) error {
	return s.auditFull(a, event, d, attempt, id, "", nil)
}

func (s *Service) auditFull(a *Action, event string, d *Denial, attempt, id, causation string, extra map[string]any) error {
	data := map[string]any{
		"action_id": a.ActionID, "action_type": a.ActionType, "action_version": a.ActionVersion,
		"level": a.Level, "target": a.Target, "params_hash": a.ParamsHash, "requested_by": a.RequestedBy,
		"actor_verified": a.ActorVerified, "policy": a.Policy, "simulated": a.Simulated, "status": a.State,
		"approvals": a.Approvals,
	}
	if a.CommandID != "" {
		data["command_id"] = a.CommandID
	}
	if a.TaskID != "" {
		data["task_id"] = a.TaskID
	}
	if d != nil {
		data["code"], data["message"] = d.Code, d.Message
		if d.Detail != "" {
			data["detail"] = d.Detail
		}
	}
	if attempt != "" {
		data["attempt_by"] = attempt
	}
	if r := a.Agent; r != nil && a.RequestedBy.Kind == "agent" {
		for k, v := range map[string]string{"agent_session": r.AgentSession, "model_alias": r.ModelAlias,
			"prompt_sha256": r.PromptSHA256, "tool_input_sha256": r.ToolInputSHA256} {
			if v != "" {
				data[k] = v
			}
		}
	}
	for k, v := range extra {
		data[k] = v
	}
	subject := auditPrefix + event
	env := map[string]any{"v": 1, "id": id, "ts": ts(s.cfg.Now()), "type": subject, "source": s.cfg.Source,
		"subject": subject, "data": data}
	if uuidRE.MatchString(a.Corr) {
		env["corr"] = a.Corr
	}
	if uuidRE.MatchString(causation) {
		env["causation"] = causation
	}
	b, err := json.Marshal(env)
	if err != nil {
		return err
	}
	if err := s.cfg.Schemas.ValidateEvent(b); err != nil {
		return fmt.Errorf("event envelope invalid: %s", firstLine(err))
	}
	if _, err := s.cfg.JS.Publish(subject, b, nats.MsgId(id)); err != nil {
		return fmt.Errorf("EVT stream did not ack %s: %w", subject, err)
	}
	return nil
}
