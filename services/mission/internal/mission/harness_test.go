package mission

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/nats-io/nats-server/v2/server"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/santhosh-tekuri/jsonschema/v6"
)

// Path into the omos tree (this package is services/mission/internal/mission).
const envelopeSchema = "../../../schemas/envelope/envelope.schema.json"

type fakeClock struct {
	mu sync.Mutex
	t  time.Time
}

func (c *fakeClock) Now() time.Time          { c.mu.Lock(); defer c.mu.Unlock(); return c.t }
func (c *fakeClock) Advance(d time.Duration) { c.mu.Lock(); c.t = c.t.Add(d); c.mu.Unlock() }

type recorder struct {
	mu   sync.Mutex
	msgs []*nats.Msg
}

func (r *recorder) add(m *nats.Msg) { r.mu.Lock(); r.msgs = append(r.msgs, m); r.mu.Unlock() }
func (r *recorder) all() []*nats.Msg {
	r.mu.Lock()
	defer r.mu.Unlock()
	return append([]*nats.Msg(nil), r.msgs...)
}

type env struct {
	t      *testing.T
	srv    *server.Server
	url    string
	pw     map[string]string
	clock  *fakeClock
	svc    *Service
	cfg    Config
	conns  map[string]*nats.Conn
	facts  *recorder // evt.mission.>
	cmds   *recorder // cmd.> (to prove the service never sends a command)
	schema *jsonschema.Schema
	kv     jetstream.KeyValue
	errs   []error
	errMu  sync.Mutex
}

func randPW() string {
	b := make([]byte, 18)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}

// natsConfig mirrors the A19 identity layout (omos #562) for the users these tests need.
// The `mission` user is the least-privilege identity proposed in
// docs/proposals/2026-10-29-mission-service-owner-proposals.md (D2-P5). Passwords are
// random per run; none is committed.
func natsConfig(storeDir string, pw map[string]string) string {
	return fmt.Sprintf(`
listen: "127.0.0.1:-1"
jetstream { store_dir: %q, max_mem: 64MB, max_file: 256MB }
accounts {
  ONEMIND {
    jetstream: enabled
    users: [
      { user: mission, password: %q, permissions: {
          publish: { allow: ["evt.mission.>", "$JS.API.INFO", "$JS.API.STREAM.INFO.*", "$JS.API.CONSUMER.CREATE.CMD.>", "$JS.API.CONSUMER.CREATE.EVT.>", "$JS.API.CONSUMER.INFO.CMD.mission-commands", "$JS.API.CONSUMER.INFO.EVT.mission-events", "$JS.API.CONSUMER.MSG.NEXT.CMD.mission-commands", "$JS.API.CONSUMER.MSG.NEXT.EVT.mission-events", "$JS.ACK.CMD.mission-commands.>", "$JS.ACK.EVT.mission-events.>", "$JS.FC.>", "$JS.API.STREAM.INFO.KV_mission-state", "$JS.API.DIRECT.GET.KV_mission-state.>", "$JS.API.STREAM.MSG.GET.KV_mission-state", "$JS.API.CONSUMER.CREATE.KV_mission-state", "$JS.API.CONSUMER.CREATE.KV_mission-state.>", "$JS.API.CONSUMER.DELETE.KV_mission-state.>", "$KV.mission-state.>"] }
          subscribe: { allow: ["svc.mission.>", "tel.fleet.*.*.position", "_INBOX.>"] }
          allow_responses: true } }
      { user: policy, password: %q, permissions: {
          publish: { allow: ["cmd.>", "evt.so.action.>", "svc.mission.apply", "$JS.API.>", "_INBOX.>"] }
          subscribe: { allow: ["evt.mission.>", "_INBOX.>"] } } }
      { user: bridge, password: %q, permissions: {
          publish: { allow: ["evt.drone.sim-x1.>", "tel.fleet.drone.sim-x1.>", "$JS.API.>"] }
          subscribe: { allow: ["cmd.fleet.drone.sim-x1.>", "_INBOX.>"] } } }
      { user: admin, password: %q }
    ]
    exports: [
      { service: "svc.mission.get" }, { service: "svc.mission.list" }, { service: "svc.mission.report" },
      { service: "svc.mission.apply", accounts: [CLIENT] }
    ]
  }
  CLIENT {
    users: [
      { user: zeus, password: %q, permissions: { publish: { allow: ["svc.mission.>"] }, subscribe: { allow: ["_INBOX.>"] } } }
      { user: alice, password: %q, permissions: { publish: { allow: ["svc.mission.>"] }, subscribe: { allow: ["_INBOX.>"] } } }
    ]
    imports: [
      { service: { account: ONEMIND, subject: "svc.mission.get" }, share: true }
      { service: { account: ONEMIND, subject: "svc.mission.list" }, share: true }
      { service: { account: ONEMIND, subject: "svc.mission.report" }, share: true }
      { service: { account: ONEMIND, subject: "svc.mission.apply" }, share: true }
    ]
  }
  AGENTS {
    users: [
      { user: hermes, password: %q, permissions: { publish: { allow: ["svc.mission.get", "svc.mission.list", "svc.mission.report"] }, subscribe: { allow: ["_INBOX.>"] } } }
    ]
    imports: [
      { service: { account: ONEMIND, subject: "svc.mission.get" }, share: true }
      { service: { account: ONEMIND, subject: "svc.mission.list" }, share: true }
      { service: { account: ONEMIND, subject: "svc.mission.report" }, share: true }
    ]
  }
}
`, storeDir, pw["mission"], pw["policy"], pw["bridge"], pw["admin"], pw["zeus"], pw["alice"], pw["hermes"])
}

var testDirectory = &Directory{
	OwnAccount:    "ONEMIND",
	AgentAccounts: []string{"AGENTS"},
	Entries: []DirectoryEntry{
		{Account: "CLIENT", User: "zeus", Kind: "user", ID: "user:zeus", Teams: []string{"team-north"}},
		{Account: "CLIENT", User: "alice", Kind: "user", ID: "user:alice"},
		{Account: "AGENTS", User: "hermes", Kind: "agent", ID: "agent:hermes/zeus/legacy"},
	},
}

func newEnv(t *testing.T) *env {
	t.Helper()
	e := &env{t: t, pw: map[string]string{}, conns: map[string]*nats.Conn{}, facts: &recorder{}, cmds: &recorder{},
		clock: &fakeClock{t: time.Date(2026, 10, 29, 12, 0, 0, 0, time.UTC)}}
	for _, u := range []string{"mission", "policy", "bridge", "admin", "zeus", "alice", "hermes"} {
		e.pw[u] = randPW()
	}
	dir := t.TempDir()
	cfgPath := filepath.Join(dir, "nats.conf")
	if err := os.WriteFile(cfgPath, []byte(natsConfig(filepath.Join(dir, "js"), e.pw)), 0o600); err != nil {
		t.Fatal(err)
	}
	opts, err := server.ProcessConfigFile(cfgPath)
	if err != nil {
		t.Fatal(err)
	}
	opts.NoLog, opts.NoSigs = true, true
	srv, err := server.NewServer(opts)
	if err != nil {
		t.Fatal(err)
	}
	go srv.Start()
	if !srv.ReadyForConnections(10 * time.Second) {
		t.Fatal("nats-server not ready")
	}
	e.srv, e.url = srv, srv.ClientURL()
	t.Cleanup(func() {
		if e.svc != nil {
			e.svc.Stop()
		}
		for _, c := range e.conns {
			c.Close()
		}
		srv.Shutdown()
		srv.WaitForShutdown()
		e.errMu.Lock()
		defer e.errMu.Unlock()
		for _, err := range e.errs {
			if strings.Contains(err.Error(), "mission") {
				t.Errorf("permission error on the mission connection: %v", err)
			}
		}
	})

	sc := jsonschema.NewCompiler()
	f, err := os.Open(envelopeSchema)
	if err != nil {
		t.Fatal(err)
	}
	doc, err := jsonschema.UnmarshalJSON(f)
	f.Close()
	if err != nil {
		t.Fatal(err)
	}
	if err := sc.AddResource("envelope.json", doc); err != nil {
		t.Fatal(err)
	}
	if e.schema, err = sc.Compile("envelope.json"); err != nil {
		t.Fatal(err)
	}

	// Streams and bucket equivalent to platform/nats/streams (CMD, EVT, MISSION, kv-mission-state).
	ctx := context.Background()
	js, _ := jetstream.New(e.conn("admin"))
	for name, subj := range map[string][]string{"CMD": {"cmd.>"}, "EVT": {"evt.so.>", "evt.drone.>", "evt.robot.>"}, "MISSION": {"evt.mission.>"}} {
		if _, err := js.CreateStream(ctx, jetstream.StreamConfig{Name: name, Subjects: subj, Storage: jetstream.FileStorage, Duplicates: 2 * time.Minute}); err != nil {
			t.Fatal(err)
		}
	}
	if e.kv, err = js.CreateKeyValue(ctx, jetstream.KeyValueConfig{Bucket: "mission-state", History: 5, Storage: jetstream.FileStorage}); err != nil {
		t.Fatal(err)
	}
	admin := e.conn("admin")
	if _, err := admin.Subscribe("evt.mission.>", e.facts.add); err != nil {
		t.Fatal(err)
	}
	if _, err := admin.Subscribe("cmd.>", e.cmds.add); err != nil {
		t.Fatal(err)
	}
	_ = admin.Flush()
	e.start()
	return e
}

// start runs the service as the `mission` user (again after a stop: restart test).
func (e *env) start() {
	e.t.Helper()
	nc := e.connFresh("mission")
	js, _ := jetstream.New(nc)
	e.cfg = Config{Conn: nc, JS: js, Now: e.clock.Now, Log: slog.New(slog.NewTextHandler(io.Discard, nil)),
		Arrival: true, Directory: testDirectory, SweepEvery: time.Hour}
	if testing.Verbose() {
		e.cfg.Log = slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}
	svc, err := New(e.cfg)
	if err != nil {
		e.t.Fatal(err)
	}
	if err := svc.Start(context.Background()); err != nil {
		e.t.Fatal(err)
	}
	e.svc = svc
}

func (e *env) connFresh(user string) *nats.Conn {
	e.t.Helper()
	nc, err := nats.Connect(e.url, nats.UserInfo(user, e.pw[user]), nats.ErrorHandler(func(_ *nats.Conn, _ *nats.Subscription, err error) {
		e.errMu.Lock()
		e.errs = append(e.errs, fmt.Errorf("%s: %w", user, err))
		e.errMu.Unlock()
	}))
	if err != nil {
		e.t.Fatalf("connect %s: %v", user, err)
	}
	e.conns[user+fmt.Sprint(len(e.conns))] = nc
	return nc
}

func (e *env) conn(user string) *nats.Conn {
	e.t.Helper()
	if c, ok := e.conns[user]; ok {
		return c
	}
	nc, err := nats.Connect(e.url, nats.UserInfo(user, e.pw[user]))
	if err != nil {
		e.t.Fatalf("connect %s: %v", user, err)
	}
	e.conns[user] = nc
	return nc
}

const ulidChars = "0123456789abcdefghjkmnpqrstvwxyz"

func newID(prefix string) string {
	b := make([]byte, 26)
	_, _ = rand.Read(b)
	for i := range b {
		b[i] = ulidChars[int(b[i])%32]
	}
	return prefix + "-" + string(b)
}

// command publishes cmd.mission.{id}.{verb} the way services/policy does
// (publishCommandID: envelope, approval, expires_at 30 s, JetStream ack). It returns the id.
func (e *env) command(id, verb string, data any, mods ...func(map[string]any)) string {
	e.t.Helper()
	cid := uuid4()
	subj := "cmd.mission." + id + "." + verb
	now := e.clock.Now()
	env := map[string]any{
		"v": 1, "id": cid, "ts": rfc3339ms(now), "cmd": subj, "source": "service:policy", "subject": subj, "data": data,
		"action_id": newID("act"), "expires_at": rfc3339ms(now.Add(30 * time.Second)), "simulated": true,
		"approval": map[string]any{"required": true, "status": "approved", "approver": "user:zeus", "level": "A3",
			"params_hash": strings.Repeat("ab", 32), "approvals": []map[string]any{{"by": "user:zeus", "at_t": now.UnixMilli(), "step": 1}}},
		"policy": map[string]any{"decision_id": "dec-1", "allowed": true},
		"corr":   "1f9a6a52-8d2b-4b8e-b7e4-5f0c2a9d3e11",
	}
	if strings.HasPrefix(id, "tsk-") {
		env["task_id"] = id
	}
	for _, m := range mods {
		m(env)
	}
	b, _ := json.Marshal(env)
	js, _ := jetstream.New(e.conn("policy"))
	if _, err := js.Publish(context.Background(), subj, b, jetstream.WithMsgID(cid)); err != nil {
		e.t.Fatalf("publish %s: %v", subj, err)
	}
	return cid
}

// taskData is cmd.mission.{tsk}.create data as services/policy builds it (dispatch, execTask).
func taskData(tsk, parent, typ string, spec any, assignee Ref) map[string]any {
	d := map[string]any{"task_id": tsk, "type": typ, "spec": spec, "description": "test",
		"assignee": assignee, "author": map[string]any{"kind": "agent", "id": "agent:hermes/zeus/legacy", "on_behalf_of": "user:zeus"},
		"action_id": newID("act"), "approval": map[string]any{"level": "A3", "approvals": []map[string]any{{"by": "user:zeus", "at_t": 1, "step": 1}}}}
	if parent != "" {
		d["parent_id"] = parent
	} else {
		d["parent_id"] = nil
	}
	return d
}

// dispatched publishes evt.so.action.dispatched as services/policy audits it.
func (e *env) dispatched(tsk, actionType, commandID string, target Ref) string {
	e.t.Helper()
	id := uuid4()
	b, _ := json.Marshal(map[string]any{"v": 1, "id": id, "ts": rfc3339ms(e.clock.Now()), "type": "evt.so.action.dispatched",
		"source": "service:policy", "subject": "evt.so.action.dispatched", "data": map[string]any{
			"action_id": newID("act"), "action_type": actionType, "target": target, "task_id": tsk,
			"command_id": commandID, "status": "dispatched"}})
	js, _ := jetstream.New(e.conn("policy"))
	if _, err := js.Publish(context.Background(), "evt.so.action.dispatched", b, jetstream.WithMsgID(id)); err != nil {
		e.t.Fatal(err)
	}
	return id
}

// report publishes evt.drone.sim-x1.task as the MAVLink bridge does (omos #550 reportTask):
// causation = the asset command id, Nats-Msg-Id = {task_id}:{status_version}.
func (e *env) report(tsk, state string, ver int, causation string, errObj map[string]string) {
	e.t.Helper()
	data := map[string]any{"task_id": tsk, "state": state, "status_version": ver, "error": errObj,
		"at_t": e.clock.Now().UnixMilli(), "simulated": true}
	env := map[string]any{"v": 1, "id": uuid4(), "ts": rfc3339ms(e.clock.Now()), "type": "evt.drone.task",
		"source": "bridge-mavlink", "subject": "evt.drone.sim-x1.task", "data": data}
	if causation != "" {
		env["causation"] = causation
	}
	b, _ := json.Marshal(env)
	js, _ := jetstream.New(e.conn("bridge"))
	if _, err := js.Publish(context.Background(), "evt.drone.sim-x1.task", b, jetstream.WithMsgID(fmt.Sprintf("%s:%d:%s", tsk, ver, causation))); err != nil {
		e.t.Fatal(err)
	}
}

// waitFact waits for a fact on evt.mission.{id}.{event} and returns its envelope.
func (e *env) waitFact(id, event string) map[string]any {
	e.t.Helper()
	return e.waitFactCause(id, event, "")
}

// waitFactCause waits for the fact with the given causation ("" = any).
func (e *env) waitFactCause(id, event, causation string) map[string]any {
	e.t.Helper()
	subj := "evt.mission." + id + "." + event
	deadline := time.Now().Add(5 * time.Second)
	for time.Now().Before(deadline) {
		for _, m := range e.facts.all() {
			if m.Subject == subj {
				var v map[string]any
				_ = json.Unmarshal(m.Data, &v)
				if causation == "" || v["causation"] == causation {
					return v
				}
			}
		}
		time.Sleep(10 * time.Millisecond)
	}
	e.t.Fatalf("no fact on %s; got %v", subj, e.subjects(id))
	return nil
}

func (e *env) subjects(id string) []string {
	var out []string
	for _, m := range e.facts.all() {
		if strings.Contains(m.Subject, id) {
			out = append(out, strings.TrimPrefix(m.Subject, "evt.mission."+id+"."))
		}
	}
	return out
}

// task reads the task record from KV.
func (e *env) task(id string) Task {
	e.t.Helper()
	ent, err := e.kv.Get(context.Background(), id)
	if err != nil {
		e.t.Fatalf("kv get %s: %v", id, err)
	}
	var t Task
	if err := json.Unmarshal(ent.Value(), &t); err != nil {
		e.t.Fatal(err)
	}
	return t
}

func (e *env) mission(id string) Mission {
	e.t.Helper()
	ent, err := e.kv.Get(context.Background(), id)
	if err != nil {
		e.t.Fatalf("kv get %s: %v", id, err)
	}
	var m Mission
	if err := json.Unmarshal(ent.Value(), &m); err != nil {
		e.t.Fatal(err)
	}
	return m
}

// waitState polls KV until the task reaches state.
func (e *env) waitState(id, state string) Task {
	e.t.Helper()
	deadline := time.Now().Add(5 * time.Second)
	var t Task
	for time.Now().Before(deadline) {
		if ent, err := e.kv.Get(context.Background(), id); err == nil {
			_ = json.Unmarshal(ent.Value(), &t)
			if t.Status.State == state {
				return t
			}
		}
		time.Sleep(10 * time.Millisecond)
	}
	e.t.Fatalf("task %s did not reach %s (at %s, facts %v)", id, state, t.Status.State, e.subjects(id))
	return t
}

// settle waits until the service has consumed everything published so far.
func (e *env) settle() {
	e.t.Helper()
	ctx := context.Background()
	js, _ := jetstream.New(e.conn("admin"))
	deadline := time.Now().Add(5 * time.Second)
	for time.Now().Before(deadline) {
		done := true
		for _, sc := range [][2]string{{"CMD", "mission-commands"}, {"EVT", "mission-events"}} {
			c, err := js.Consumer(ctx, sc[0], sc[1])
			if err != nil {
				done = false
				continue
			}
			info, err := c.Info(ctx)
			if err != nil || info.NumPending > 0 || info.NumAckPending > 0 {
				done = false
			}
		}
		if done {
			time.Sleep(20 * time.Millisecond)
			return
		}
		time.Sleep(10 * time.Millisecond)
	}
	e.t.Fatal("service did not settle")
}

// validateFacts checks every fact against the omos event envelope v1 schema and the
// section 7 rules for task facts.
func (e *env) validateFacts() {
	e.t.Helper()
	for _, m := range e.facts.all() {
		v, err := jsonschema.UnmarshalJSON(strings.NewReader(string(m.Data)))
		if err != nil {
			e.t.Fatalf("%s: %v", m.Subject, err)
		}
		if err := e.schema.Validate(v); err != nil {
			e.t.Errorf("%s does not validate against the envelope schema: %v", m.Subject, err)
		}
		var env Event
		_ = json.Unmarshal(m.Data, &env)
		parts := strings.Split(m.Subject, ".")
		kind := "task"
		if strings.HasPrefix(parts[2], "msn-") {
			kind = "mission"
		}
		if want := "evt.mission." + kind + "." + parts[3]; env.Type != want || env.Subject != m.Subject {
			e.t.Errorf("%s: type %q subject %q, want type %q", m.Subject, env.Type, env.Subject, want)
		}
		if kind == "task" && parts[3] != "rejected" {
			var d struct {
				StatusVersion int `json:"status_version"`
				Version       struct {
					Status int `json:"status"`
				} `json:"version"`
			}
			_ = json.Unmarshal(env.Data, &d)
			ver := d.StatusVersion
			if parts[3] == "created" {
				ver = d.Version.Status
			}
			if want := fmt.Sprintf("%s:%d", parts[2], ver); m.Header.Get("Nats-Msg-Id") != want {
				e.t.Errorf("%s: Nats-Msg-Id %q, want %q", m.Subject, m.Header.Get("Nats-Msg-Id"), want)
			}
		}
	}
}

// request sends a request as user and decodes the reply.
func (e *env) request(user, subject string, body any) map[string]any {
	e.t.Helper()
	b, _ := json.Marshal(body)
	msg, err := e.conn(user).Request(subject, b, 3*time.Second)
	if err != nil {
		e.t.Fatalf("%s %s: %v", user, subject, err)
	}
	var v map[string]any
	_ = json.Unmarshal(msg.Data, &v)
	return v
}

func errCode(v map[string]any) string {
	if e, ok := v["error"].(map[string]any); ok {
		s, _ := e["code"].(string)
		return s
	}
	return ""
}
