package mission

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"math"
	"sort"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

// Config configures the task service. Zero values take the defaults in New.
type Config struct {
	Conn   *nats.Conn
	JS     jetstream.JetStream
	Log    *slog.Logger
	Now    func() time.Time
	Source string // envelope source and updated_by of service-made changes (default service:task)

	Bucket        string        // KV bucket (default mission-state)
	CmdStream     string        // stream holding cmd.mission.> (default CMD)
	EvtStream     string        // stream holding evt.so.action.> and asset reports (default EVT)
	ReportClasses []string      // asset classes whose evt.{class}.*.task EvtStream binds (default drone, robot)
	ConsumerName  string        // durable consumer name prefix (default mission)
	SweepEvery    time.Duration // timeout sweep period (default 1s)

	DefaultDelivery time.Duration // delivery.before_t default after creation (default 60s, section 3.2)
	DefaultStale    time.Duration // stale_t default when neither stale_t nor complete_before_t is set (default 1h)

	// CacheTerminalFor is how long a terminal task or mission stays in the in-memory copy
	// (svc.mission.list) after its last change (default 24h). KV keeps the record; get
	// reads KV, so an evicted record is still served by svc.mission.get.
	CacheTerminalFor time.Duration

	// Arrival: complete an asset goto task from fleet telemetry when the asset is within
	// ArrivalRadiusM (horizontal) and ArrivalVerticalM of the goal (owner proposal D2-P2).
	Arrival          bool
	ArrivalRadiusM   float64 // default 5
	ArrivalVerticalM float64 // default 5

	Directory *Directory // principals for svc.mission.report (nil: every report is refused)
}

// Counters are exported on /metrics by the command.
type Counters struct {
	CommandsApplied  atomic.Int64
	CommandsRejected atomic.Int64
	ReportsAccepted  atomic.Int64
	ReportsStale     atomic.Int64 // dropped by the status_version guard (rule 1)
	ReportsTerminal  atomic.Int64 // dropped because the task is terminal (rule 2)
	ReportsRefused   atomic.Int64 // not from the assignee, malformed, unknown task
	Timeouts         atomic.Int64
	Arrivals         atomic.Int64
	FactErrors       atomic.Int64
	// CommandsRedelivered counts commands already applied (ack lost) and not applied again.
	CommandsRedelivered atomic.Int64
}

// Service is the task service.
type Service struct {
	cfg   Config
	st    *store
	C     Counters
	stop  context.CancelFunc
	ctx   context.Context
	subs  []*nats.Subscription
	cons  []jetstream.ConsumeContext
	watch jetstream.KeyWatcher
	wg    sync.WaitGroup
}

// New checks the configuration and fills defaults.
func New(cfg Config) (*Service, error) {
	if cfg.Conn == nil || cfg.JS == nil {
		return nil, errors.New("mission: Conn and JS are required")
	}
	if cfg.Log == nil {
		cfg.Log = slog.Default()
	}
	if cfg.Now == nil {
		cfg.Now = time.Now
	}
	def := func(s *string, v string) {
		if *s == "" {
			*s = v
		}
	}
	def(&cfg.Source, "service:task")
	def(&cfg.Bucket, "mission-state")
	def(&cfg.CmdStream, "CMD")
	def(&cfg.EvtStream, "EVT")
	def(&cfg.ConsumerName, "mission")
	if len(cfg.ReportClasses) == 0 {
		cfg.ReportClasses = []string{"drone", "robot"}
	}
	for _, c := range cfg.ReportClasses {
		if !AssetClasses[c] {
			return nil, fmt.Errorf("mission: report class %q is not an asset class", c)
		}
	}
	if cfg.SweepEvery == 0 {
		cfg.SweepEvery = time.Second
	}
	if cfg.DefaultDelivery == 0 {
		cfg.DefaultDelivery = 60 * time.Second
	}
	if cfg.DefaultStale == 0 {
		cfg.DefaultStale = time.Hour
	}
	if cfg.CacheTerminalFor == 0 {
		cfg.CacheTerminalFor = 24 * time.Hour
	}
	if cfg.ArrivalRadiusM == 0 {
		cfg.ArrivalRadiusM = 5
	}
	if cfg.ArrivalVerticalM == 0 {
		cfg.ArrivalVerticalM = 5
	}
	return &Service{cfg: cfg}, nil
}

// Start loads KV mission-state, then consumes commands, dispatch audits and asset reports,
// serves svc.mission.*, and runs the sweep.
func (s *Service) Start(ctx context.Context) error {
	s.ctx, s.stop = context.WithCancel(ctx)
	kv, err := s.cfg.JS.KeyValue(s.ctx, s.cfg.Bucket)
	if err != nil {
		return fmt.Errorf("KV %s: %w", s.cfg.Bucket, err)
	}
	s.st = newStore(kv)
	if s.watch, err = s.st.watch(s.ctx); err != nil {
		return fmt.Errorf("watch %s: %w", s.cfg.Bucket, err)
	}

	cmd, err := s.consumer(s.cfg.CmdStream, s.cfg.ConsumerName+"-commands", []string{"cmd.mission.>"})
	if err != nil {
		return err
	}
	evtSubjects := []string{"evt.so.action.dispatched"}
	for _, c := range s.cfg.ReportClasses {
		evtSubjects = append(evtSubjects, "evt."+c+".*.task")
	}
	evt, err := s.consumer(s.cfg.EvtStream, s.cfg.ConsumerName+"-events", evtSubjects)
	if err != nil {
		return err
	}
	for _, c := range []struct {
		cons jetstream.Consumer
		h    func(jetstream.Msg) error
	}{{cmd, s.handleCommandMsg}, {evt, s.handleEventMsg}} {
		h := c.h
		cc, err := c.cons.Consume(func(m jetstream.Msg) {
			if err := h(m); err != nil {
				s.cfg.Log.Warn("will retry", "subject", m.Subject(), "err", err)
				_ = m.NakWithDelay(time.Second)
				return
			}
			_ = m.Ack()
		})
		if err != nil {
			return err
		}
		s.cons = append(s.cons, cc)
	}

	for subj, h := range map[string]nats.MsgHandler{
		"svc.mission.get":    s.handleGet,
		"svc.mission.list":   s.handleList,
		"svc.mission.report": s.handleReport,
		"svc.mission.apply":  s.handleApply,
	} {
		sub, err := s.cfg.Conn.QueueSubscribe(subj, "mission", h)
		if err != nil {
			return err
		}
		s.subs = append(s.subs, sub)
	}
	if s.cfg.Arrival {
		sub, err := s.cfg.Conn.Subscribe("tel.fleet.*.*.position", s.handlePosition)
		if err != nil {
			return err
		}
		s.subs = append(s.subs, sub)
	}
	if err := s.cfg.Conn.Flush(); err != nil {
		return err
	}
	s.EvictOnce()
	s.wg.Add(1)
	go s.sweepLoop()
	s.cfg.Log.Info("task service ready", "bucket", s.cfg.Bucket, "records", len(s.st.cachedTasks())+len(s.st.cachedMissions()))
	return nil
}

// Stop ends every subscription and the sweep.
func (s *Service) Stop() {
	for _, c := range s.cons {
		c.Stop()
	}
	for _, sub := range s.subs {
		_ = sub.Unsubscribe()
	}
	if s.stop != nil {
		s.stop()
	}
	s.wg.Wait()
	if s.watch != nil {
		_ = s.watch.Stop()
	}
}

func (s *Service) consumer(stream, name string, subjects []string) (jetstream.Consumer, error) {
	cfg := jetstream.ConsumerConfig{
		Durable: name, DeliverPolicy: jetstream.DeliverNewPolicy, AckPolicy: jetstream.AckExplicitPolicy,
		AckWait: 30 * time.Second, MaxDeliver: 10, MaxAckPending: 256,
	}
	if len(subjects) == 1 {
		cfg.FilterSubject = subjects[0]
	} else {
		cfg.FilterSubjects = subjects
	}
	c, err := s.cfg.JS.CreateOrUpdateConsumer(s.ctx, stream, cfg)
	if err != nil {
		return nil, fmt.Errorf("consumer %s on %s: %w", name, stream, err)
	}
	return c, nil
}

func (s *Service) nowMS() int64 { return s.cfg.Now().UnixMilli() }

// ---------------------------------------------------------------------------
// Facts: evt.mission.{id}.{event} in the omos event envelope v1 (section 7).

func uuid4() string {
	var u [16]byte
	_, _ = rand.Read(u[:])
	u[6] = u[6]&0x0f | 0x40
	u[8] = u[8]&0x3f | 0x80
	return fmt.Sprintf("%x-%x-%x-%x-%x", u[0:4], u[4:6], u[6:8], u[8:10], u[10:16])
}

func rfc3339ms(t time.Time) string { return t.UTC().Format("2006-01-02T15:04:05.000Z07:00") }

type factOpts struct {
	msgID     string
	causation string
	corr      string
	trace     string
}

// publishFact publishes one fact with a JetStream ack from MISSION. The record is already
// written; a failed publish is retried and then counted, and the record stays correct.
func (s *Service) publishFact(id, event string, data any, o factOpts) {
	subject := "evt.mission." + id + "." + event
	kind := "task"
	if strings.HasPrefix(id, "msn-") {
		kind = "mission"
	}
	raw, _ := json.Marshal(data)
	env := Event{V: 1, ID: uuid4(), TS: rfc3339ms(s.cfg.Now()), Type: "evt.mission." + kind + "." + event,
		Source: s.cfg.Source, Subject: subject, Data: raw}
	if uuidRE.MatchString(o.causation) {
		env.Causation = o.causation
	}
	if uuidRE.MatchString(o.corr) {
		env.Corr = o.corr
	}
	env.Trace = o.trace
	msgID := o.msgID
	if msgID == "" {
		msgID = env.ID
	}
	b, _ := json.Marshal(env)
	var err error
	for attempt := 0; attempt < 3; attempt++ {
		ctx, cancel := context.WithTimeout(s.ctx, 5*time.Second)
		_, err = s.cfg.JS.Publish(ctx, subject, b, jetstream.WithMsgID(msgID))
		cancel()
		if err == nil {
			return
		}
	}
	s.C.FactErrors.Add(1)
	s.cfg.Log.Error("fact not stored", "subject", subject, "err", err)
}

// taskFact is the data of evt.mission.{tsk}.{state} after creation: the section 7 report
// with the task service's own status_version (= version.status) plus context.
func taskFact(t *Task) map[string]any {
	d := map[string]any{
		"task_id": t.TaskID, "state": t.Status.State, "status_version": t.Version.Status,
		"error": t.Status.Error, "at_t": t.Status.UpdatedT, "simulated": t.Simulated,
		"parent_id": t.ParentID, "assignee": t.Assignee, "rollup": t.Rollup, "updated_by": t.UpdatedBy,
		"action_id": t.ActionID,
	}
	for k, v := range map[string]json.RawMessage{"progress": t.Status.Progress, "estimate": t.Status.Estimate, "result": t.Status.Result} {
		if len(v) > 0 && string(v) != "null" {
			d[k] = v
		}
	}
	return d
}

func (s *Service) emitTask(t *Task, causation string) {
	s.publishFact(t.TaskID, t.Status.State, taskFact(t),
		factOpts{msgID: fmt.Sprintf("%s:%d", t.TaskID, t.Version.Status), causation: causation, corr: t.Corr, trace: t.Trace})
}

func (s *Service) emitMission(m *Mission, event, causation string, extra map[string]any) {
	var d map[string]any
	b, _ := json.Marshal(m)
	_ = json.Unmarshal(b, &d)
	d["event"] = event
	for k, v := range extra {
		d[k] = v
	}
	s.publishFact(m.MissionID, event, d, factOpts{causation: causation, corr: m.Corr})
}

// reject reports a refused command as evt.mission.{id}.rejected. Its causation is the
// command id, so the policy service closes the action as failed with data.code (POL-14).
func (s *Service) reject(id, verb string, c *Command, code, msg string) {
	s.C.CommandsRejected.Add(1)
	s.cfg.Log.Warn("command rejected", "id", id, "verb", verb, "code", code, "message", msg)
	if !IsMissionID(id) && !IsTaskID(id) {
		return
	}
	d := map[string]any{"verb": verb, "code": code, "message": msg, "at_t": s.nowMS()}
	o := factOpts{}
	if c != nil {
		d["action_id"], d["command_id"] = c.ActionID, c.ID
		o.causation, o.corr = c.ID, c.Corr
	}
	s.publishFact(id, "rejected", d, o)
}

// ---------------------------------------------------------------------------
// Commands: cmd.mission.{id}.{verb} from the policy service (sections 6.1-6.3).

// retryable marks a handler error worth a redelivery (KV or stream unavailable).
type retryable struct{ error }

func (s *Service) handleCommandMsg(m jetstream.Msg) error {
	return s.HandleCommand(m.Subject(), m.Data())
}

// HandleCommand applies one command. It returns a retryable error only when KV failed.
func (s *Service) HandleCommand(subject string, data []byte) error {
	parts := strings.Split(subject, ".")
	if len(parts) != 4 || parts[0] != "cmd" || parts[1] != "mission" {
		return nil
	}
	id, verb := parts[2], parts[3]
	var c Command
	if err := json.Unmarshal(data, &c); err != nil {
		s.reject(id, verb, nil, CodeSchemaInvalid, "command envelope: "+err.Error())
		return nil
	}
	if c.V != 1 || !uuidRE.MatchString(c.ID) || c.Subject != subject {
		s.reject(id, verb, &c, CodeSchemaInvalid, "command envelope: v must be 1, id a UUID and subject the delivery subject")
		return nil
	}
	if !IsMissionID(id) && !IsTaskID(id) {
		s.reject(id, verb, &c, CodeSchemaInvalid, "id must be msn-{ulid} or tsk-{ulid}")
		return nil
	}
	// A redelivery of a command already applied (its ack was lost) is recognised before the
	// expiry and precondition checks, so it never produces a spurious rejection.
	if done, err := s.alreadyApplied(id, verb, &c); err != nil {
		return retryable{err}
	} else if done {
		return nil
	}
	exp, err := time.Parse(time.RFC3339Nano, c.ExpiresAt)
	if err != nil {
		s.reject(id, verb, &c, CodeSchemaInvalid, "expires_at missing or not RFC 3339")
		return nil
	}
	if !s.cfg.Now().Before(exp) {
		// Never applied late (asset-command.md section 4 "No queued replay").
		s.reject(id, verb, &c, CodeExpired, "command expired at "+c.ExpiresAt)
		return nil
	}
	if c.Approval.Status != "approved" && c.Approval.Status != "not_required" {
		s.reject(id, verb, &c, CodeNotAuthorized, "command carries no approval")
		return nil
	}
	var err2 error
	switch {
	case IsMissionID(id) && verb == "create":
		err2 = s.missionCreate(id, &c)
	case IsMissionID(id) && verb == "waypoint":
		err2 = s.missionWaypoint(id, &c)
	case IsMissionID(id) && verb == "complete":
		err2 = s.missionComplete(id, &c)
	case IsMissionID(id) && verb == "abort":
		err2 = s.missionAbort(id, &c)
	case IsTaskID(id) && verb == "create":
		err2 = s.taskCreate(id, &c)
	case IsTaskID(id) && verb == "abort":
		err2 = s.taskAbort(id, &c, true)
	case IsTaskID(id) && verb == "complete":
		err2 = s.taskComplete(id, &c)
	default:
		// Doctrine section 9.1 verbs only; cancelling is abort.
		s.reject(id, verb, &c, CodeUnknownAction, "verb "+verb+" is not create, waypoint, abort or complete for this id")
		return nil
	}
	if err2 != nil {
		var r retryable
		if errors.As(err2, &r) {
			return err2
		}
		s.cfg.Log.Error("command failed", "subject", subject, "err", err2)
		return nil
	}
	s.C.CommandsApplied.Add(1)
	return nil
}

// alreadyApplied reports whether the record already shows this command as applied: its id
// is in applied_commands, or (task create) the task exists for the same action. A create is
// reported again so the action can close if the first fact was lost (same Nats-Msg-Id for a
// task, so MISSION keeps one copy inside the duplicate window); any other verb is a no-op,
// because its fact was published when it was applied.
func (s *Service) alreadyApplied(id, verb string, c *Command) (bool, error) {
	if IsTaskID(id) {
		t, _, err := s.st.getTask(s.ctx, id)
		if err != nil || t == nil {
			return false, err
		}
		if !contains(t.Applied, c.ID) {
			if verb != "create" {
				return false, nil
			}
			var p struct {
				ActionID string `json:"action_id"`
			}
			_ = json.Unmarshal(c.Data, &p)
			if p.ActionID == "" {
				p.ActionID = c.ActionID
			}
			if p.ActionID == "" || t.ActionID != p.ActionID {
				return false, nil
			}
		}
		if verb == "create" {
			s.publishFact(id, StateCreated, t, factOpts{msgID: id + ":1", causation: c.ID, corr: t.Corr})
		}
		s.C.CommandsRedelivered.Add(1)
		s.cfg.Log.Info("command already applied (redelivery)", "id", id, "verb", verb, "command_id", c.ID)
		return true, nil
	}
	if !IsMissionID(id) {
		return false, nil
	}
	m, _, err := s.st.getMission(s.ctx, id)
	if err != nil || m == nil || !contains(m.Applied, c.ID) {
		return false, err
	}
	if verb == "create" {
		s.emitMission(m, "created", c.ID, map[string]any{"action_id": c.ActionID})
	}
	s.C.CommandsRedelivered.Add(1)
	s.cfg.Log.Info("command already applied (redelivery)", "id", id, "verb", verb, "command_id", c.ID)
	return true, nil
}

func by(c *Command) string {
	if c.Approval.Approver != "" {
		return c.Approval.Approver
	}
	return c.Source
}

func (s *Service) missionCreate(id string, c *Command) error {
	var p struct {
		Name   string  `json:"name"`
		Intent string  `json:"intent"`
		AOI    *string `json:"aoi"`
		Window *Window `json:"window"`
	}
	if err := json.Unmarshal(c.Data, &p); err != nil || strings.TrimSpace(p.Name) == "" || len(p.Name) > 120 || len(p.Intent) > 2000 {
		s.reject(id, "create", c, CodeSchemaInvalid, "mission.create needs name (1-120 chars), intent up to 2000")
		return nil
	}
	now := s.nowMS()
	m := &Mission{V: 1, MissionID: id, Name: p.Name, Intent: p.Intent, State: MissionPlanned,
		Rollup: MissionRollup(MissionPlanned), AOI: p.AOI, Route: []Waypoint{}, Assignments: []Assignment{},
		TaskIDs: []string{}, Window: p.Window, Marking: Marking{TLP: "AMBER", ReleasableTo: []string{}},
		TAKMission: json.RawMessage("null"), CreatedT: now, CreatedBy: c.Source, UpdatedT: now, UpdatedBy: by(c), Corr: c.Corr,
		Applied: []string{c.ID}}
	if _, err := s.st.createMission(s.ctx, m); err != nil {
		if !errors.Is(err, jetstream.ErrKeyExists) {
			return retryable{err}
		}
		old, _, gerr := s.st.getMission(s.ctx, id)
		if gerr != nil {
			return retryable{gerr}
		}
		if old != nil && old.Name == p.Name && old.CreatedBy == c.Source {
			// Redelivery of the same create: report it again so the action can close.
			s.emitMission(old, "created", c.ID, map[string]any{"action_id": c.ActionID})
			return nil
		}
		s.reject(id, "create", c, CodeConflict, "mission "+id+" exists")
		return nil
	}
	s.emitMission(m, "created", c.ID, map[string]any{"action_id": c.ActionID})
	return nil
}

func (s *Service) missionWaypoint(id string, c *Command) error {
	var w Waypoint
	var raw map[string]json.RawMessage
	if json.Unmarshal(c.Data, &raw) != nil || raw["seq"] == nil || raw["lat"] == nil || raw["lon"] == nil || raw["hae_m"] == nil ||
		json.Unmarshal(c.Data, &w) != nil || w.Seq < 0 || math.Abs(w.Lat) > 90 || math.Abs(w.Lon) > 180 {
		s.reject(id, "waypoint", c, CodeSchemaInvalid, "mission.waypoint needs seq >= 0, lat, lon, hae_m")
		return nil
	}
	reason := ""
	m, err := s.st.mutateMission(s.ctx, id, func(m *Mission) error {
		if m.State == MissionCompleted || m.State == MissionAborted {
			reason = "mission is " + m.State
			return errNoChange
		}
		replaced := false
		for i := range m.Route {
			if m.Route[i].Seq == w.Seq {
				m.Route[i], replaced = w, true
			}
		}
		if !replaced {
			m.Route = append(m.Route, w)
			sort.SliceStable(m.Route, func(i, j int) bool { return m.Route[i].Seq < m.Route[j].Seq })
		}
		m.UpdatedT, m.UpdatedBy = s.nowMS(), by(c)
		m.Applied = addApplied(m.Applied, c.ID)
		return nil
	})
	return s.finishMission(id, "waypoint", "waypoint", c, m, err, reason, map[string]any{"action_id": c.ActionID, "waypoint": w})
}

func (s *Service) missionComplete(id string, c *Command) error {
	reason := ""
	m, err := s.st.mutateMission(s.ctx, id, func(m *Mission) error {
		if m.State == MissionCompleted || m.State == MissionAborted {
			reason = "mission is " + m.State
			return errNoChange
		}
		m.State, m.Rollup = MissionCompleted, MissionRollup(MissionCompleted)
		m.UpdatedT, m.UpdatedBy = s.nowMS(), by(c)
		m.Applied = addApplied(m.Applied, c.ID)
		return nil
	})
	return s.finishMission(id, "complete", "completed", c, m, err, reason, map[string]any{"action_id": c.ActionID})
}

// missionAbort aborts the mission and every non-terminal child task (section 3.1).
func (s *Service) missionAbort(id string, c *Command) error {
	reason := ""
	m, err := s.st.mutateMission(s.ctx, id, func(m *Mission) error {
		if m.State == MissionCompleted || m.State == MissionAborted {
			reason = "mission is " + m.State
			return errNoChange
		}
		m.State, m.Rollup = MissionAborted, MissionRollup(MissionAborted)
		m.UpdatedT, m.UpdatedBy = s.nowMS(), by(c)
		m.Applied = addApplied(m.Applied, c.ID)
		return nil
	})
	if ferr := s.finishMission(id, "abort", "aborted", c, m, err, reason, map[string]any{"action_id": c.ActionID}); ferr != nil || reason != "" || m == nil {
		return ferr
	}
	for _, tid := range m.TaskIDs {
		if err := s.taskAbort(tid, c, false); err != nil {
			s.cfg.Log.Error("child abort failed", "mission", id, "task", tid, "err", err)
		}
	}
	return nil
}

// finishMission turns a mission mutation result into a fact or a rejection.
func (s *Service) finishMission(id, verb, event string, c *Command, m *Mission, err error, reason string, extra map[string]any) error {
	switch {
	case errors.Is(err, errNoChange):
		s.reject(id, verb, c, CodePreconditionFailed, reason)
		return nil
	case err != nil:
		return retryable{err}
	case m == nil:
		s.reject(id, verb, c, CodeNotFound, "no mission "+id)
		return nil
	}
	s.emitMission(m, event, c.ID, extra)
	return nil
}

// createParams is the data of cmd.mission.{tsk}.create (section 6.1 step 4): the
// task.assign parameters plus the policy service's assignee, author, action and approval.
type createParams struct {
	TaskID          string            `json:"task_id"`
	ParentID        *string           `json:"parent_id"`
	Type            string            `json:"type"`
	Spec            json.RawMessage   `json:"spec"`
	Description     string            `json:"description"`
	Delivery        *Delivery         `json:"delivery"`
	Execution       *Execution        `json:"execution"`
	StaleT          int64             `json:"stale_t"`
	Assignee        Ref               `json:"assignee"`
	Author          Author            `json:"author"`
	ActionID        string            `json:"action_id"`
	Approval        json.RawMessage   `json:"approval"`
	TargetSnapshots []json.RawMessage `json:"target_snapshots"`
	Allocation      []Ref             `json:"allocation"`
}

func (s *Service) taskCreate(id string, c *Command) error {
	var p createParams
	if err := json.Unmarshal(c.Data, &p); err != nil {
		s.reject(id, "create", c, CodeSchemaInvalid, "task create: "+err.Error())
		return nil
	}
	if p.TaskID != id || (c.TaskID != "" && c.TaskID != id) {
		s.reject(id, "create", c, CodeSchemaInvalid, "task_id must equal the subject id")
		return nil
	}
	if _, ok := TaskTypes[p.Type]; !ok {
		s.reject(id, "create", c, CodeUnknownAction, "task type "+p.Type+" is not in the catalogue")
		return nil
	}
	if len(p.Spec) == 0 || p.Spec[0] != '{' {
		s.reject(id, "create", c, CodeSchemaInvalid, "spec must be an object")
		return nil
	}
	if !AssigneeTypes[p.Assignee.EntityType] || !tokenRE.MatchString(p.Assignee.EntityID) {
		s.reject(id, "create", c, CodeSchemaInvalid, "assignee must be {entity_type: drone|robot|vehicle|camera|person|team|agent, entity_id}")
		return nil
	}
	asset := AssetClasses[p.Assignee.EntityType]
	if asset && !p.Assignee.Simulated() {
		s.reject(id, "create", c, CodeNotSimulated, "asset tasks are for simulated sim-* assets only")
		return nil
	}
	if !typeFits(p.Type, p.Assignee.EntityType) {
		s.reject(id, "create", c, CodeSchemaInvalid, "task type "+p.Type+" does not fit assignee "+p.Assignee.EntityType)
		return nil
	}
	if why := personTarget(p.Spec, p.TargetSnapshots); why != "" {
		// Defence in depth behind the policy service (stop line SL2, contract section 10
		// test 4): a task is never pointed at a person or a person-class track.
		s.reject(id, "create", c, CodeForbiddenAction, why)
		return nil
	}
	if p.ActionID == "" {
		p.ActionID = c.ActionID
	}
	var parentMission string
	if p.ParentID != nil {
		pid := *p.ParentID
		switch {
		case IsMissionID(pid):
			m, _, err := s.st.getMission(s.ctx, pid)
			if err != nil {
				return retryable{err}
			}
			if m == nil || m.State == MissionCompleted || m.State == MissionAborted {
				s.reject(id, "create", c, CodePreconditionFailed, "parent mission "+pid+" does not exist or is finished")
				return nil
			}
			parentMission = pid
		case IsTaskID(pid):
			pt, _, err := s.st.getTask(s.ctx, pid)
			if err != nil {
				return retryable{err}
			}
			if pt == nil {
				s.reject(id, "create", c, CodePreconditionFailed, "parent task "+pid+" does not exist")
				return nil
			}
		default:
			s.reject(id, "create", c, CodeSchemaInvalid, "parent_id must be msn-… or tsk-…")
			return nil
		}
	}

	now := s.nowMS()
	delivery := Delivery{BeforeT: now + s.cfg.DefaultDelivery.Milliseconds(), RequireAck: true}
	if p.Delivery != nil && p.Delivery.BeforeT > 0 {
		delivery.BeforeT = p.Delivery.BeforeT
	}
	stale := p.StaleT
	if stale == 0 && p.Execution != nil && p.Execution.CompleteBeforeT != nil {
		stale = *p.Execution.CompleteBeforeT
	}
	if stale == 0 {
		stale = now + s.cfg.DefaultStale.Milliseconds()
	}
	t := &Task{V: 1, TaskID: id, ParentID: p.ParentID, Type: p.Type, Spec: p.Spec, Description: p.Description,
		Author: p.Author, Assignee: p.Assignee, TargetSnapshots: p.TargetSnapshots, Delivery: delivery,
		Execution: p.Execution, StaleT: stale, ActionID: p.ActionID, Approval: p.Approval,
		Version:    Version{Definition: 1, Status: 1},
		Status:     Status{State: StateCreated, UpdatedT: now},
		Allocation: p.Allocation, Rollup: Rollup(StateCreated, ""), Simulated: asset && p.Assignee.Simulated(),
		CreatedT: now, CreatedBy: c.Source, UpdatedT: now, UpdatedBy: c.Source, Corr: c.Corr, Trace: c.Trace,
		Guard: Guard{Commands: []string{}}, Applied: []string{c.ID}}
	if t.TargetSnapshots == nil {
		t.TargetSnapshots = []json.RawMessage{}
	}
	if t.Allocation == nil {
		t.Allocation = []Ref{}
	}
	if len(t.Approval) == 0 {
		ap, _ := json.Marshal(map[string]any{"level": c.Approval.Level, "approvals": c.Approval.Approvals})
		t.Approval = ap
	}
	if _, err := s.st.createTask(s.ctx, t); err != nil {
		if !errors.Is(err, jetstream.ErrKeyExists) {
			return retryable{err}
		}
		old, _, gerr := s.st.getTask(s.ctx, id)
		if gerr != nil {
			return retryable{gerr}
		}
		if old != nil && old.ActionID == t.ActionID {
			// Redelivery: the record exists for this action; report created again (same
			// Nats-Msg-Id {task}:1, so MISSION keeps one copy inside the window).
			s.publishFact(id, StateCreated, old, factOpts{msgID: id + ":1", causation: c.ID, corr: old.Corr})
			return nil
		}
		s.reject(id, "create", c, CodeConflict, "task "+id+" exists for another action")
		return nil
	}
	// Section 6.1 step 5: data = the task record. Its causation (the create command id)
	// closes person, team and agent task actions in the policy service (POL-14).
	s.publishFact(id, StateCreated, t, factOpts{msgID: id + ":1", causation: c.ID, corr: t.Corr, trace: t.Trace})

	if t.Execution != nil && t.Execution.StartAfterT != nil && *t.Execution.StartAfterT > now {
		if nt, err := s.st.mutateTask(s.ctx, id, func(x *Task) error {
			if x.Status.State != StateCreated {
				return errNoChange
			}
			s.setState(x, StateScheduled, nil, c.Source)
			return nil
		}); err == nil && nt != nil {
			s.emitTask(nt, c.ID)
		}
	}
	if parentMission != "" {
		s.linkToMission(parentMission, t, c)
	}
	return nil
}

// personTarget returns why a task spec or its target snapshots point at a person or a
// person-class track, or "". The task service has no entity lookup of its own (the policy
// service resolves spec entity references before approval); it checks the frozen
// target_snapshots and any entity_type or kind the spec's target or center carries.
func personTarget(spec json.RawMessage, snaps []json.RawMessage) string {
	isPerson := func(raw json.RawMessage) bool {
		var e struct {
			EntityType string `json:"entity_type"`
			Kind       string `json:"kind"`
		}
		return json.Unmarshal(raw, &e) == nil && (e.EntityType == "person" || strings.EqualFold(e.Kind, "person"))
	}
	var sp map[string]json.RawMessage
	if json.Unmarshal(spec, &sp) == nil {
		for _, k := range []string{"target", "center"} {
			if v, ok := sp[k]; ok && isPerson(v) {
				return "spec " + k + " is a person or a person-class track"
			}
		}
	}
	for _, raw := range snaps {
		var ts struct {
			EntityID string          `json:"entity_id"`
			Snapshot json.RawMessage `json:"snapshot"`
		}
		if json.Unmarshal(raw, &ts) == nil && (isPerson(ts.Snapshot) || isPerson(raw)) {
			return "target " + ts.EntityID + " is a person or a person-class track"
		}
	}
	return ""
}

// typeFits checks the Assignee column of section 4.
func typeFits(typ, assigneeType string) bool {
	switch typ {
	case "onemind.task.manual.v1":
		return assigneeType == "person" || assigneeType == "team"
	case "onemind.task.agent_request.v1":
		return assigneeType == "agent"
	case "onemind.task.survey_area.v1":
		return AssetClasses[assigneeType] || assigneeType == "team"
	}
	return AssetClasses[assigneeType]
}

// linkToMission adds the task to its mission's task_ids and its assignee to the mission's
// assignments (link assigned_to, role "tasked"), and publishes evt.mission.{msn}.updated.
func (s *Service) linkToMission(msn string, t *Task, c *Command) {
	m, err := s.st.mutateMission(s.ctx, msn, func(m *Mission) error {
		changed := false
		if !contains(m.TaskIDs, t.TaskID) {
			m.TaskIDs = append(m.TaskIDs, t.TaskID)
			changed = true
		}
		to := t.Assignee.String()
		found := false
		for _, a := range m.Assignments {
			if a.To == to {
				found = true
			}
		}
		if !found {
			m.Assignments = append(m.Assignments, Assignment{To: to, Role: "tasked", SinceT: s.nowMS(), By: by(c)})
			changed = true
		}
		if t.Simulated && !m.Simulated {
			m.Simulated, changed = true, true
		}
		if !changed {
			return errNoChange
		}
		m.UpdatedT, m.UpdatedBy = s.nowMS(), s.cfg.Source
		return nil
	})
	if err != nil && !errors.Is(err, errNoChange) {
		s.cfg.Log.Error("mission link failed", "mission", msn, "task", t.TaskID, "err", err)
		return
	}
	if err == nil && m != nil {
		s.emitMission(m, "updated", c.ID, map[string]any{"task_id": t.TaskID})
	}
}

// setState records a state change: version.status +1, rollup, times.
func (s *Service) setState(t *Task, state string, e *TaskError, who string) {
	now := s.nowMS()
	t.Status.State, t.Status.Error = state, e
	t.Status.UpdatedT = now
	if t.Status.StartT == nil && (state == StateExecuting || state == StateDoneOK) {
		t.Status.StartT = &now
	}
	if IsTerminal(state) {
		t.Status.EndT = &now
	}
	code := ""
	if e != nil {
		code = e.Code
	}
	t.Rollup = Rollup(state, code)
	t.Version.Status++
	t.UpdatedT, t.UpdatedBy = now, who
}

// taskAbort applies cmd.mission.{tsk}.abort (section 6.2), or a mission abort's child
// abort when direct is false (then a terminal child is skipped, not rejected).
func (s *Service) taskAbort(id string, c *Command, direct bool) error {
	reason := ""
	t, err := s.st.mutateTask(s.ctx, id, func(t *Task) error {
		switch st := t.Status.State; {
		case IsTerminal(st):
			reason = "task is " + st
			return errNoChange
		case st == StateCancelRequested:
			reason = "abort already requested"
			return errNoChange
		case st == StateCreated || st == StateScheduled:
			// Not yet sent: cancelled at once (rule 3).
			s.setState(t, StateDoneNotOK, &TaskError{Code: ErrCancelled, Message: "aborted before dispatch"}, by(c))
		default:
			t.Guard.PrevState = st
			s.setState(t, StateCancelRequested, nil, by(c))
		}
		t.Applied = addApplied(t.Applied, c.ID)
		return nil
	})
	switch {
	case errors.Is(err, errNoChange):
		if direct {
			s.reject(id, "abort", c, CodePreconditionFailed, reason)
		}
		return nil
	case err != nil:
		return retryable{err}
	case t == nil:
		if direct {
			s.reject(id, "abort", c, CodeNotFound, "no task "+id)
		}
		return nil
	}
	s.emitTask(t, c.ID)
	s.updateMissionRollup(t, c.ID)
	// The hold may have reported (or been audited) before this abort was applied.
	s.resolveHold(id, c.ID)
	return nil
}

func (s *Service) taskComplete(id string, c *Command) error {
	reason := ""
	t, err := s.st.mutateTask(s.ctx, id, func(t *Task) error {
		if IsTerminal(t.Status.State) || t.Status.State == StateCompleteRequested || t.Status.State == StateCancelRequested {
			reason = "task is " + t.Status.State
			return errNoChange
		}
		s.setState(t, StateCompleteRequested, nil, by(c))
		t.Applied = addApplied(t.Applied, c.ID)
		return nil
	})
	switch {
	case errors.Is(err, errNoChange):
		s.reject(id, "complete", c, CodePreconditionFailed, reason)
		return nil
	case err != nil:
		return retryable{err}
	case t == nil:
		s.reject(id, "complete", c, CodeNotFound, "no task "+id)
		return nil
	}
	s.emitTask(t, c.ID)
	return nil
}

// ---------------------------------------------------------------------------
// Events: evt.so.action.dispatched (policy) and evt.{class}.{id}.task (bridges).

func (s *Service) handleEventMsg(m jetstream.Msg) error {
	return s.HandleEvent(m.Subject(), m.Data())
}

// HandleEvent applies one dispatch audit or asset status report.
func (s *Service) HandleEvent(subject string, data []byte) error {
	var env Event
	if err := json.Unmarshal(data, &env); err != nil {
		s.C.ReportsRefused.Add(1)
		return nil
	}
	if subject == "evt.so.action.dispatched" {
		return s.handleDispatched(&env)
	}
	parts := strings.Split(subject, ".")
	if len(parts) != 4 || parts[3] != "task" || !AssetClasses[parts[1]] {
		return nil
	}
	var r Report
	if err := json.Unmarshal(env.Data, &r); err != nil {
		s.C.ReportsRefused.Add(1)
		return nil
	}
	_, err := s.ApplyReport(r, Ref{EntityType: parts[1], EntityID: parts[2]}, env.Source, env.Causation, env.ID, true)
	return err
}

// handleDispatched sets sent on evt.so.action.dispatched for a task.assign (section 5) and
// records the asset command id in the report guard (task.assign and the abort hold).
func (s *Service) handleDispatched(env *Event) error {
	var d struct {
		ActionType string `json:"action_type"`
		TaskID     string `json:"task_id"`
		CommandID  string `json:"command_id"`
		Target     Ref    `json:"target"`
	}
	if json.Unmarshal(env.Data, &d) != nil || !IsTaskID(d.TaskID) {
		return nil
	}
	if d.ActionType != "task.assign" && d.ActionType != "task.abort" {
		return nil
	}
	asset := AssetClasses[d.Target.EntityType]
	sent := false
	t, err := s.st.mutateTask(s.ctx, d.TaskID, func(t *Task) error {
		if IsTerminal(t.Status.State) {
			return errNoChange
		}
		changed := false
		if asset && uuidRE.MatchString(d.CommandID) && !contains(t.Guard.Commands, d.CommandID) {
			t.Guard.Commands = append(t.Guard.Commands, d.CommandID)
			changed = true
		}
		if asset && d.ActionType == "task.abort" && uuidRE.MatchString(d.CommandID) && !contains(t.Guard.AbortCommands, d.CommandID) {
			// The hold is known by its audit, whether or not the abort command is applied yet.
			t.Guard.AbortCommands = append(t.Guard.AbortCommands, d.CommandID)
			changed = true
		}
		if d.ActionType == "task.assign" && (t.Status.State == StateCreated || t.Status.State == StateScheduled) {
			s.setState(t, StateSent, nil, s.cfg.Source)
			sent = true
			return nil
		}
		if !changed {
			return errNoChange
		}
		// Only the guard changed: written without a fact (version.status unchanged).
		return nil
	})
	if errors.Is(err, errNoChange) {
		return nil
	}
	if err != nil {
		return retryable{err}
	}
	if t == nil {
		// For a person, team or agent task the policy service audits the dispatch right
		// after publishing the create, so this audit can beat the record: redeliver it
		// (bounded by the consumer's MaxDeliver) until the create has been applied.
		return retryable{fmt.Errorf("task %s not recorded yet", d.TaskID)}
	}
	if sent {
		s.emitTask(t, env.ID)
	}
	if d.ActionType == "task.abort" {
		s.resolveHold(d.TaskID, env.ID)
	}
	return nil
}

// holdOutcome maps a hold report to the state it gives a cancel_requested task (section
// 6.2): executing or done_ok means the hold took effect (cancelled); done_not_ok means the
// asset refused it (previous state, rejected). Anything else changes nothing yet.
func holdOutcome(g *Guard, state string, e *TaskError) (string, *TaskError, bool) {
	switch state {
	case StateExecuting, StateDoneOK:
		return StateDoneNotOK, &TaskError{Code: ErrCancelled, Message: "aborted: hold executed by the asset"}, true
	case StateDoneNotOK:
		prev := g.PrevState
		if prev == "" || IsTerminal(prev) || prev == StateCancelRequested {
			prev = StateExecuting
		}
		msg := "abort refused"
		if e != nil && e.Message != "" {
			msg += ": " + e.Message
		}
		return prev, &TaskError{Code: ErrRejected, Message: msg}, true
	}
	return "", nil, false
}

// resolveHold applies a hold report that arrived before the task entered cancel_requested.
func (s *Service) resolveHold(id, causation string) {
	t, err := s.st.mutateTask(s.ctx, id, func(t *Task) error {
		h := t.Guard.Hold
		if t.Status.State != StateCancelRequested || h == nil {
			return errNoChange
		}
		state, terr, ok := holdOutcome(&t.Guard, h.State, h.Error)
		if !ok {
			return errNoChange
		}
		t.Guard.Hold = nil
		s.setState(t, state, terr, h.Source)
		if len(h.Progress) > 0 {
			t.Status.Progress = h.Progress
		}
		return nil
	})
	if err != nil || t == nil {
		if err != nil && !errors.Is(err, errNoChange) {
			s.cfg.Log.Error("hold resolution failed", "task", id, "err", err)
		}
		return
	}
	s.emitTask(t, causation)
	s.updateMissionRollup(t, causation)
}

// ReportResult says what happened to a report.
type ReportResult struct {
	Accepted      bool   `json:"accepted"`
	Reason        string `json:"reason,omitempty"`
	State         string `json:"state,omitempty"`
	StatusVersion int    `json:"status_version,omitempty"`
}

var (
	errStale    = errors.New("stale")
	errTerminal = errors.New("terminal")
	errNotYours = errors.New("not the assignee")
)

// ApplyReport applies one assignee status report (sections 5 and 7). fromBridge marks an
// asset report from evt.{class}.{id}.task, whose causation is the asset command id.
func (s *Service) ApplyReport(r Report, reporter Ref, source, causation, envID string, fromBridge bool) (ReportResult, error) {
	if !IsTaskID(r.TaskID) || !assigneeStates[r.State] || r.StatusVersion < 1 {
		s.C.ReportsRefused.Add(1)
		return ReportResult{Reason: CodeSchemaInvalid + ": task_id, an assignee state and status_version >= 1 are required"}, nil
	}
	if r.State == StateDoneNotOK {
		if r.Error == nil || !errorCodes[r.Error.Code] {
			r.Error = &TaskError{Code: ErrFailed, Message: "done_not_ok without a valid error code"}
		}
	} else {
		r.Error = nil
	}
	var reason string
	guardOnly := false
	var sentSnap *Task
	t, err := s.st.mutateTask(s.ctx, r.TaskID, func(t *Task) error {
		guardOnly, reason, sentSnap = false, "", nil
		if t.Assignee != reporter {
			return errNotYours
		}
		if IsTerminal(t.Status.State) {
			return errTerminal
		}
		g := &t.Guard
		idx, added := g.Epoch, false
		if causation != "" && fromBridge {
			idx = indexOf(g.Commands, causation)
			if idx < 0 {
				// The report beat the dispatch audit: the command is the newest one.
				g.Commands = append(g.Commands, causation)
				idx, added = len(g.Commands)-1, true
			}
		}
		if !(idx > g.Epoch || (idx == g.Epoch && r.StatusVersion > g.LastVersion)) {
			if added {
				g.Commands = g.Commands[:len(g.Commands)-1]
			}
			reason = fmt.Sprintf("status_version %d of command %d is not newer than %d of command %d", r.StatusVersion, idx, g.LastVersion, g.Epoch)
			return errStale
		}
		g.Epoch, g.LastVersion, g.Reporter = idx, r.StatusVersion, source

		// A hold is an abort command: audited as task.abort, or any command after the task's
		// one task.assign command (its report can beat both the abort and the audit).
		isHold := fromBridge && causation != "" && (contains(g.AbortCommands, causation) || idx > 0)
		state, terr := r.State, r.Error
		switch {
		case isHold:
			g.Hold = &HoldReport{Command: causation, State: r.State, Error: r.Error, Progress: r.Progress, Source: source, EnvID: envID}
			if t.Status.State != StateCancelRequested {
				// The abort is not applied yet: keep the state; resolveHold applies this report
				// when the abort lands.
				guardOnly, reason = true, "abort hold reported before the abort was applied; state kept"
				return nil
			}
			var ok bool
			if state, terr, ok = holdOutcome(g, r.State, r.Error); !ok {
				guardOnly, reason = true, "abort pending; state kept"
				return nil
			}
			g.Hold = nil
		case t.Status.State == StateCancelRequested && !IsTerminal(state):
			// Progress from the original command while the abort is pending.
			guardOnly, reason = true, "abort pending; state kept"
			return nil
		}
		if t.Status.State == StateCreated || t.Status.State == StateScheduled {
			// The assignee's report beat the policy service's dispatch audit (they are about
			// 1 ms apart on different streams): a report proves the task was sent, so sent is
			// recorded first and the lifecycle facts stay created, sent, received, ...; the
			// late audit then finds the task past sent and changes nothing (section 5).
			s.setState(t, StateSent, nil, s.cfg.Source)
			snap := *t
			sentSnap = &snap
		}
		s.setState(t, state, terr, source)
		for dst, src := range map[*json.RawMessage]json.RawMessage{&t.Status.Progress: r.Progress, &t.Status.Estimate: r.Estimate, &t.Status.Result: r.Result} {
			if len(src) > 0 {
				*dst = src
			}
		}
		return nil
	})
	switch {
	case err == nil && t != nil && guardOnly:
		// Only the guard moved (written without a fact; version.status unchanged).
		s.C.ReportsAccepted.Add(1)
		return ReportResult{Accepted: true, Reason: reason, State: t.Status.State}, nil
	case errors.Is(err, errStale):
		s.C.ReportsStale.Add(1)
		s.cfg.Log.Info("report dropped (out of order)", "task", r.TaskID, "reason", reason)
		return ReportResult{Reason: "stale: " + reason}, nil
	case errors.Is(err, errTerminal):
		s.C.ReportsTerminal.Add(1)
		s.cfg.Log.Info("report dropped (task terminal)", "task", r.TaskID, "state", r.State)
		return ReportResult{Reason: "terminal: the task state never changes again"}, nil
	case errors.Is(err, errNotYours):
		s.C.ReportsRefused.Add(1)
		s.cfg.Log.Warn("report refused (not the assignee)", "task", r.TaskID, "reporter", reporter.String())
		return ReportResult{Reason: CodeNotAuthorized + ": reporter is not the assignee"}, nil
	case err != nil:
		return ReportResult{}, retryable{err}
	case t == nil:
		s.C.ReportsRefused.Add(1)
		return ReportResult{Reason: CodeNotFound + ": no task " + r.TaskID}, nil
	}
	s.C.ReportsAccepted.Add(1)
	if sentSnap != nil {
		cause := causation
		if cause == "" {
			cause = envID
		}
		s.emitTask(sentSnap, cause) // as handleDispatched: no mission rollup on sent
	}
	s.emitTask(t, envID)
	s.updateMissionRollup(t, envID)
	return ReportResult{Accepted: true, State: t.Status.State, StatusVersion: t.Version.Status}, nil
}

// conclude sets a service-decided state (timeout, arrival) if the task still allows it.
func (s *Service) conclude(id string, allow func(t *Task) bool, state string, e *TaskError, result json.RawMessage) *Task {
	t, err := s.st.mutateTask(s.ctx, id, func(t *Task) error {
		if IsTerminal(t.Status.State) || !allow(t) {
			return errNoChange
		}
		s.setState(t, state, e, s.cfg.Source)
		if result != nil {
			t.Status.Result = result
		}
		return nil
	})
	if err != nil || t == nil {
		if err != nil && !errors.Is(err, errNoChange) {
			s.cfg.Log.Error("conclude failed", "task", id, "err", err)
		}
		return nil
	}
	s.emitTask(t, "")
	s.updateMissionRollup(t, "")
	return t
}

// updateMissionRollup follows a child task change: planned -> active when a child first
// executes; completed when every child is done_ok (section 3.1).
func (s *Service) updateMissionRollup(t *Task, causation string) {
	if t.ParentID == nil || !IsMissionID(*t.ParentID) {
		return
	}
	msn := *t.ParentID
	event := ""
	m, err := s.st.mutateMission(s.ctx, msn, func(m *Mission) error {
		if m.State == MissionCompleted || m.State == MissionAborted {
			return errNoChange
		}
		allOK := len(m.TaskIDs) > 0
		for _, id := range m.TaskIDs {
			st := ""
			if id == t.TaskID {
				st = t.Status.State
			} else if ct, ok := s.st.cachedTask(id); ok {
				st = ct.Status.State
			} else if kt, _, err := s.st.getTask(s.ctx, id); err == nil && kt != nil {
				// Evicted from the copy (terminal for longer than CacheTerminalFor).
				st = kt.Status.State
			}
			if st != StateDoneOK {
				allOK = false
			}
		}
		switch {
		case allOK:
			m.State, event = MissionCompleted, "completed"
		case m.State == MissionPlanned && (t.Status.State == StateExecuting || t.Status.State == StateDoneOK):
			m.State, event = MissionActive, "updated"
		default:
			return errNoChange
		}
		m.Rollup = MissionRollup(m.State)
		m.UpdatedT, m.UpdatedBy = s.nowMS(), s.cfg.Source
		return nil
	})
	if err != nil || m == nil {
		if err != nil && !errors.Is(err, errNoChange) {
			s.cfg.Log.Error("mission rollup failed", "mission", msn, "err", err)
		}
		return
	}
	s.emitMission(m, event, causation, map[string]any{"task_id": t.TaskID})
}

// ---------------------------------------------------------------------------
// Sweep: delivery and stale timeouts (section 5 rule 4).

func (s *Service) sweepLoop() {
	defer s.wg.Done()
	tk := time.NewTicker(s.cfg.SweepEvery)
	defer tk.Stop()
	ev := time.NewTicker(time.Minute)
	defer ev.Stop()
	for {
		select {
		case <-s.ctx.Done():
			return
		case <-tk.C:
			s.SweepOnce()
		case <-ev.C:
			s.EvictOnce()
		}
	}
}

// EvictOnce drops terminal records older than CacheTerminalFor from the in-memory copy.
func (s *Service) EvictOnce() int {
	n := s.st.evictTerminal(s.nowMS() - s.cfg.CacheTerminalFor.Milliseconds())
	if n > 0 {
		s.cfg.Log.Info("terminal records dropped from the in-memory copy", "count", n)
	}
	return n
}

// SweepOnce times out every open task with no receipt by delivery.before_t, or past
// stale_t. It never re-sends anything: re-tasking is a new proposal.
func (s *Service) SweepOnce() {
	now := s.nowMS()
	for _, t := range s.st.openTasks() {
		if IsTerminal(t.Status.State) {
			continue
		}
		notReceived := t.Status.State == StateCreated || t.Status.State == StateScheduled || t.Status.State == StateSent
		var msg string
		switch {
		case notReceived && t.Delivery.BeforeT > 0 && now > t.Delivery.BeforeT:
			msg = "not received by delivery.before_t"
		case t.StaleT > 0 && now > t.StaleT:
			msg = "stale_t passed"
		default:
			continue
		}
		if s.conclude(t.TaskID, func(x *Task) bool {
			nr := x.Status.State == StateCreated || x.Status.State == StateScheduled || x.Status.State == StateSent
			return (nr && x.Delivery.BeforeT > 0 && now > x.Delivery.BeforeT) || (x.StaleT > 0 && now > x.StaleT)
		}, StateDoneNotOK, &TaskError{Code: ErrTimeout, Message: msg}, nil) != nil {
			s.C.Timeouts.Add(1)
		}
	}
}

// ---------------------------------------------------------------------------
// Arrival: tel.fleet.{class}.{id}.position against open goto tasks (owner proposal D2-P2).

var runningStates = map[string]bool{StateReceived: true, StateAcknowledged: true, StateWilco: true,
	StateExecuting: true, StateWaiting: true, StateCompleteRequested: true}

func (s *Service) handlePosition(m *nats.Msg) {
	parts := strings.Split(m.Subject, ".")
	if len(parts) != 5 || !AssetClasses[parts[2]] {
		return
	}
	ids := s.st.openGotos(parts[2] + "." + parts[3])
	if len(ids) == 0 {
		return
	}
	var pos struct {
		UID string   `json:"uid"`
		Lat *float64 `json:"lat"`
		Lon *float64 `json:"lon"`
		Hae *float64 `json:"hae"`
	}
	if json.Unmarshal(m.Data, &pos) != nil || pos.Lat == nil || pos.Lon == nil || pos.Hae == nil {
		return
	}
	for _, id := range ids {
		t, ok := s.st.cachedTask(id)
		if !ok || t.Status.State != StateExecuting {
			continue
		}
		var g struct {
			Lat  float64 `json:"lat"`
			Lon  float64 `json:"lon"`
			HaeM float64 `json:"hae_m"`
		}
		if json.Unmarshal(t.Spec, &g) != nil {
			continue
		}
		d := HaversineM(*pos.Lat, *pos.Lon, g.Lat, g.Lon)
		dv := math.Abs(*pos.Hae - g.HaeM)
		if d > s.cfg.ArrivalRadiusM || dv > s.cfg.ArrivalVerticalM {
			continue
		}
		res, _ := json.Marshal(map[string]any{"arrived": true, "source": "telemetry", "subject": m.Subject,
			"distance_m": math.Round(d*100) / 100, "vertical_m": math.Round(dv*100) / 100})
		if s.conclude(id, func(x *Task) bool { return runningStates[x.Status.State] }, StateDoneOK, nil, res) != nil {
			s.C.Arrivals.Add(1)
		}
	}
}

// HaversineM is the great-circle distance in metres.
func HaversineM(lat1, lon1, lat2, lon2 float64) float64 {
	const r = 6371008.8
	p1, p2 := lat1*math.Pi/180, lat2*math.Pi/180
	dp, dl := (lat2-lat1)*math.Pi/180, (lon2-lon1)*math.Pi/180
	a := math.Sin(dp/2)*math.Sin(dp/2) + math.Cos(p1)*math.Cos(p2)*math.Sin(dl/2)*math.Sin(dl/2)
	return 2 * r * math.Asin(math.Min(1, math.Sqrt(a)))
}

func contains(xs []string, x string) bool { return indexOf(xs, x) >= 0 }

func indexOf(xs []string, x string) int {
	for i, v := range xs {
		if v == x {
			return i
		}
	}
	return -1
}
