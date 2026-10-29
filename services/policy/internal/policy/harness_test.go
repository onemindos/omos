package policy

import (
	"context"
	"crypto/ed25519"
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
)

// Paths into the omos tree (this package is services/policy/internal/policy).
const (
	schemaDir = "../../../schemas"
	modelFile = "../../../../platform/security/openfga/onemind-actions.json"
)

type fakeClock struct {
	mu sync.Mutex
	t  time.Time
}

func (c *fakeClock) Now() time.Time          { c.mu.Lock(); defer c.mu.Unlock(); return c.t }
func (c *fakeClock) Advance(d time.Duration) { c.mu.Lock(); c.t = c.t.Add(d); c.mu.Unlock() }

// swapAuthz lets a test replace the authorizer (OpenFGA down, HTTP stub) at runtime.
type swapAuthz struct {
	mu sync.RWMutex
	a  Authorizer
}

func (s *swapAuthz) get() Authorizer  { s.mu.RLock(); defer s.mu.RUnlock(); return s.a }
func (s *swapAuthz) set(a Authorizer) { s.mu.Lock(); s.a = a; s.mu.Unlock() }
func (s *swapAuthz) StoreID() string  { return s.get().StoreID() }
func (s *swapAuthz) ModelID() string  { return s.get().ModelID() }
func (s *swapAuthz) Check(ctx context.Context, req CheckRequest) (bool, error) {
	return s.get().Check(ctx, req)
}

// recorder collects every message on a subject filter in arrival order.
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

type bridgeStub struct {
	mu       sync.Mutex
	pub      ed25519.PublicKey
	clock    *fakeClock
	silent   map[string]bool
	seen     map[string]bool // nonce
	received []map[string]any
	refused  []string
}

type env struct {
	t        *testing.T
	srv      *server.Server
	url      string
	pw       map[string]string
	clock    *fakeClock
	svc      *Service
	authz    *swapAuthz
	eval     *ModelEvaluator
	schemas  *Schemas
	admin    *nats.Conn
	cmds     *recorder
	audit    *recorder
	bridge   *bridgeStub
	entities map[string]map[string]any
	entMu    sync.Mutex
	entMicro bool // not-found replies in the NATS micro form (empty body, error headers)
	taskDown bool
	geoDown  bool // geofence stub silent
	taskMu   sync.Mutex
	conns    map[string]*nats.Conn
	errs     map[string][]error
	errMu    sync.Mutex
}

func randPW() string {
	b := make([]byte, 18)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}

// natsConfig is the per-identity layout proposed in docs/proposals (stage 1), shrunk to
// the identities these tests need. Passwords are random per run; none is committed.
func natsConfig(storeDir string, pw map[string]string) string {
	return fmt.Sprintf(`
listen: "127.0.0.1:-1"
jetstream { store_dir: %q, max_mem: 64MB, max_file: 256MB }
accounts {
  ONEMIND {
    jetstream: enabled
    users: [
      { user: policy, password: %q, permissions: {
          publish: { allow: ["cmd.>", "evt.so.action.>", "svc.entity.get", "$JS.API.>", "_INBOX.>"] }
          subscribe: { allow: ["svc.policy.>", "evt.drone.>", "evt.robot.>", "evt.vehicle.>", "evt.mission.>", "evt.entity.>", "evt.map.geofence.*", "_INBOX.>"] }
          allow_responses: true } }
      { user: entity, password: %q, permissions: {
          publish: { allow: ["_INBOX.>"] }, subscribe: { allow: ["svc.entity.>"] }, allow_responses: true } }
      { user: bridge, password: %q, permissions: {
          publish: { allow: ["evt.drone.sim-x1.>", "evt.drone.sim-x2.>", "tel.fleet.drone.sim-x1.>", "tel.fleet.drone.sim-x2.>"] }
          subscribe: { allow: ["cmd.fleet.drone.sim-x1.>", "cmd.fleet.drone.sim-x2.>"] } } }
      { user: tasksvc, password: %q, permissions: {
          publish: { allow: ["evt.mission.>"] }, subscribe: { allow: ["cmd.mission.>"] } } }
      { user: geofence, password: %q, permissions: {
          publish: { allow: ["evt.map.geofence.*"] }, subscribe: { allow: ["cmd.map.geofence.>"] } } }
      { user: watch, password: %q, permissions: {
          publish: { allow: ["svc.policy.command.propose", "svc.policy.task.propose"] }, subscribe: { allow: ["_INBOX.>"] } } }
      { user: other, password: %q, permissions: {
          publish: { allow: [">"], deny: ["svc.policy.>", "cmd.>", "evt.so.action.>"] }, subscribe: { allow: ["_INBOX.>"] } } }
      { user: admin, password: %q }
    ]
    exports: [ { service: "svc.policy.>" } ]
  }
  CLIENT {
    users: [
      { user: zeus, password: %q, permissions: {
          publish: { allow: ["svc.policy.>", "cmd.map.>"] }, subscribe: { allow: ["_INBOX.>", "cmd.map.>"] } } }
      { user: alice, password: %q, permissions: {
          publish: { allow: ["svc.policy.>", "cmd.map.>"] }, subscribe: { allow: ["_INBOX.>", "cmd.map.>"] } } }
      { user: carol, password: %q, permissions: {
          publish: { allow: ["svc.policy.>"] }, subscribe: { allow: ["_INBOX.>"] } } }
      { user: mallory, password: %q, permissions: {
          publish: { allow: ["svc.policy.>"] }, subscribe: { allow: ["_INBOX.>"] } } }
    ]
    imports: [ { service: { account: ONEMIND, subject: "svc.policy.>" }, share: true } ]
  }
  AGENTS {
    users: [
      { user: hermes, password: %q, permissions: {
          publish: { allow: ["svc.policy.command.propose", "svc.policy.task.propose", "svc.policy.action.get", "agents.>", "cmd.map.>"] }
          subscribe: { allow: ["_INBOX.>", "agents.prompt.hermes.>"] } } }
      { user: rogue, password: %q, permissions: {
          publish: { allow: ["svc.policy.>"] }, subscribe: { allow: ["_INBOX.>"] } } }
    ]
    imports: [ { service: { account: ONEMIND, subject: "svc.policy.>" }, share: true } ]
  }
}
`, storeDir, pw["policy"], pw["entity"], pw["bridge"], pw["tasksvc"], pw["geofence"], pw["watch"], pw["other"], pw["admin"],
		pw["zeus"], pw["alice"], pw["carol"], pw["mallory"], pw["hermes"], pw["rogue"])
}

var testDirectory = Directory{
	OwnAccount:    "ONEMIND",
	AgentAccounts: []string{"AGENTS"},
	Entries: []DirectoryEntry{
		{Account: "CLIENT", User: "zeus", Kind: "user", ID: "user:zeus"},
		{Account: "CLIENT", User: "alice", Kind: "user", ID: "user:alice"},
		{Account: "CLIENT", User: "carol", Kind: "user", ID: "user:carol"},
		{Account: "AGENTS", User: "hermes", Kind: "agent", ID: "agent:hermes/zeus/legacy"},
		// Mislabelled on purpose: the AGENTS account must still resolve it to an agent.
		{Account: "AGENTS", User: "rogue", Kind: "user", ID: "user:rogue"},
	},
}

func testTuples() []Tuple {
	ts := []Tuple{
		{User: "user:zeus", Relation: "operator", Object: "asset:sim-x1"},
		{User: "user:zeus", Relation: "operator", Object: "asset:sim-x2"},
		{User: "user:zeus", Relation: "operator", Object: "asset:sim-x3"},
		{User: "site:home", Relation: "site", Object: "asset:sim-x1"},
		{User: "user:alice", Relation: "operator", Object: "site:home"},
		{User: "user:rogue", Relation: "operator", Object: "asset:sim-x1"},
	}
	for name, spec := range Catalogue {
		if spec.Level == A5 {
			continue
		}
		obj := "action_type:" + name
		for _, u := range []string{"user:zeus", "user:alice", "user:rogue"} {
			ts = append(ts, Tuple{User: u, Relation: "proposer", Object: obj},
				Tuple{User: u, Relation: "approver", Object: obj},
				Tuple{User: u, Relation: "second_approver", Object: obj})
		}
		ts = append(ts, Tuple{User: "agent:hermes/zeus/legacy", Relation: "proposer", Object: obj},
			Tuple{User: "service:watch", Relation: "proposer", Object: obj})
	}
	return ts
}

func newEnv(t *testing.T) *env {
	t.Helper()
	e := &env{t: t, pw: map[string]string{}, clock: &fakeClock{t: time.Date(2026, 10, 29, 12, 0, 0, 0, time.UTC)},
		conns: map[string]*nats.Conn{}, errs: map[string][]error{}, entities: map[string]map[string]any{}}
	for _, u := range []string{"policy", "entity", "bridge", "tasksvc", "geofence", "watch", "other", "admin", "zeus", "alice", "carol", "mallory", "hermes", "rogue"} {
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
		for _, c := range e.conns {
			c.Close()
		}
		srv.Shutdown()
		srv.WaitForShutdown()
	})

	e.admin = e.conn("admin")
	js, _ := e.admin.JetStream()
	for _, sc := range []*nats.StreamConfig{
		{Name: "CMD", Subjects: []string{"cmd.>"}, Storage: nats.FileStorage, Discard: nats.DiscardNew, Duplicates: 2 * time.Minute},
		{Name: "EVT", Subjects: []string{"evt.so.>", "evt.drone.>", "evt.robot.>"}, Storage: nats.FileStorage, Duplicates: 2 * time.Minute},
		{Name: "MISSION", Subjects: []string{"evt.mission.>"}, Storage: nats.FileStorage},
	} {
		if _, err := js.AddStream(sc); err != nil {
			t.Fatal(err)
		}
	}
	kv, err := js.CreateKeyValue(&nats.KeyValueConfig{Bucket: "geofences", History: 10})
	if err != nil {
		t.Fatal(err)
	}
	for k, v := range map[string]string{
		// keep_out: a small square (the pond) north-east of the homestead.
		"zone-pond": `{"entity_id":"zone-pond","entity_type":"zone","live":true,"geo":{"shape":"polygon","kind":"keep_out",
		  "coords":{"type":"Polygon","coordinates":[[[-97.705,30.295],[-97.695,30.295],[-97.695,30.305],[-97.705,30.305],[-97.705,30.295]]]}}}`,
		// keep_in: 20 km around the homestead for every asset.
		"zone-home": `{"entity_id":"zone-home","entity_type":"zone","live":true,"geo":{"shape":"circle","kind":"keep_in","radius_m":20000,
		  "coords":{"type":"Point","coordinates":[-97.7431,30.2672]}}}`,
	} {
		if _, err := kv.Put(k, []byte(v)); err != nil {
			t.Fatal(err)
		}
	}

	e.cmds, e.audit = &recorder{}, &recorder{}
	if _, err := e.admin.Subscribe("cmd.>", e.cmds.add); err != nil {
		t.Fatal(err)
	}
	if _, err := e.admin.Subscribe("evt.so.action.>", e.audit.add); err != nil {
		t.Fatal(err)
	}
	_ = e.admin.Flush()

	// Entity service stub (svc.entity.get).
	online := func(id string, caps ...string) map[string]any {
		return map[string]any{"entity_id": id, "entity_type": "drone", "kind": "multirotor",
			"location":   map[string]any{"lat": 30.2672, "lon": -97.7431, "hae_m": 152.4},
			"health":     map[string]any{"connection": "online"},
			"indicators": map[string]any{"simulated": true}, "capabilities": caps}
	}
	e.entities["sim-x1"] = online("sim-x1", "onemind.task.goto.v1", "onemind.task.orbit.v1", "onemind.task.rtl.v1",
		"onemind.task.land.v1", "onemind.task.observe.v1", "onemind.task.hold.v1")
	e.entities["sim-x2"] = online("sim-x2", "onemind.task.goto.v1")
	x3 := online("sim-x3", "onemind.task.goto.v1")
	x3["health"] = map[string]any{"connection": "offline"}
	e.entities["sim-x3"] = x3
	e.entities["person-1"] = map[string]any{"entity_id": "person-1", "entity_type": "person"}
	e.entities["trk-9"] = map[string]any{"entity_id": "trk-9", "entity_type": "track", "kind": "person"}
	e.entities["trk-car"] = map[string]any{"entity_id": "trk-car", "entity_type": "track", "kind": "car"}
	ec := e.conn("entity")
	if _, err := ec.Subscribe("svc.entity.get", func(m *nats.Msg) {
		var q struct {
			EntityID string `json:"entity_id"`
		}
		_ = json.Unmarshal(m.Data, &q)
		e.entMu.Lock()
		rec, ok := e.entities[q.EntityID]
		microForm := e.entMicro
		var b []byte
		if ok {
			b, _ = json.Marshal(map[string]any{"entity": rec})
		} else if microForm {
			e.entMu.Unlock()
			// What micro's Request.Error("404", "entity not found", nil) sends.
			r := nats.NewMsg(m.Reply)
			r.Header.Set("Nats-Service-Error-Code", "404")
			r.Header.Set("Nats-Service-Error", "entity not found")
			_ = m.RespondMsg(r)
			return
		} else {
			b = []byte(`{"error":{"code":"not_found"}}`)
		}
		e.entMu.Unlock()
		_ = m.Respond(b)
	}); err != nil {
		t.Fatal(err)
	}
	_ = ec.Flush()

	// Signing key, generated per run.
	pub, priv, _ := ed25519.GenerateKey(rand.Reader)
	e.bridge = &bridgeStub{pub: pub, clock: e.clock, silent: map[string]bool{}, seen: map[string]bool{}}
	bc := e.conn("bridge")
	for _, id := range []string{"sim-x1", "sim-x2"} { // the bridge may subscribe only its own ids
		if _, err := bc.Subscribe("cmd.fleet.drone."+id+".>", func(m *nats.Msg) { e.bridge.handle(bc, m) }); err != nil {
			t.Fatal(err)
		}
	}
	_ = bc.Flush()

	// Task service stub: records cmd.mission.{tsk}.create as evt.mission.{tsk}.created.
	tc := e.conn("tasksvc")
	if _, err := tc.Subscribe("cmd.mission.*.create", func(m *nats.Msg) {
		e.taskMu.Lock()
		down := e.taskDown
		e.taskMu.Unlock()
		if down {
			return
		}
		tsk := strings.Split(m.Subject, ".")[2]
		_ = tc.Publish("evt.mission."+tsk+".created", []byte(`{"v":1}`))
	}); err != nil {
		t.Fatal(err)
	}
	_ = tc.Flush()

	// Geofence service stub (D1-P2): answers cmd.map.geofence.{id}.{verb} on the four-token
	// evt.map.geofence.{created|updated|deleted|refused} with the zone id in data. A zone id
	// starting "bad-" is refused.
	gc := e.conn("geofence")
	if _, err := gc.Subscribe("cmd.map.geofence.*.*", func(m *nats.Msg) {
		e.taskMu.Lock()
		down := e.geoDown
		e.taskMu.Unlock()
		if down {
			return
		}
		var c struct {
			ID       string `json:"id"`
			ActionID string `json:"action_id"`
		}
		_ = json.Unmarshal(m.Data, &c)
		parts := strings.Split(m.Subject, ".")
		zone, verb := parts[3], parts[4]
		event := map[string]string{"create": "created", "update": "updated", "remove": "deleted"}[verb]
		data := map[string]any{"zone_id": zone, "action_id": c.ActionID, "outcome": "executed", "revision": 1}
		if strings.HasPrefix(zone, "bad-") {
			event = "refused"
			data["outcome"], data["code"], data["message"] = "refused", "schema_invalid", "zone refused"
		}
		b, _ := json.Marshal(map[string]any{"v": 1, "id": NewUUID(), "causation": c.ID, "data": data})
		_ = gc.Publish("evt.map.geofence."+event, b)
	}); err != nil {
		t.Fatal(err)
	}
	_ = gc.Flush()

	// The policy service.
	pc := e.conn("policy")
	pjs, _ := pc.JetStream()
	e.schemas, err = LoadSchemas(schemaDir)
	if err != nil {
		t.Fatal(err)
	}
	e.eval, err = LoadModel(modelFile)
	if err != nil {
		t.Fatal(err)
	}
	if err := e.eval.Write(testTuples()...); err != nil {
		t.Fatal(err)
	}
	e.authz = &swapAuthz{a: e.eval}
	zones, err := WatchZones(t.Context(), pjs, "geofences", nil)
	if err != nil {
		t.Fatal(err)
	}
	dirCopy := testDirectory
	e.svc, err = New(Config{Conn: pc, JS: pjs, Schemas: e.schemas, Authz: e.authz, Directory: &dirCopy,
		Entities: &NATSEntities{Conn: pc, Timeout: time.Second}, Zones: zones,
		Signer: &Signer{Key: priv, KID: "test-run"}, Now: e.clock.Now,
		Log:        slog.New(slog.NewTextHandler(io.Discard, nil)),
		SweepEvery: 50 * time.Millisecond, TaskCreateTimeout: 400 * time.Millisecond, CheckTimeout: time.Second})
	if err != nil {
		t.Fatal(err)
	}
	if err := e.svc.Start(); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(e.svc.Stop)
	// Wait until the geofences are loaded.
	waitFor(t, func() bool { _, err := zones.Zones(); return err == nil }, "geofences loaded")
	return e
}

func (e *env) conn(user string) *nats.Conn {
	if c, ok := e.conns[user]; ok {
		return c
	}
	c, err := nats.Connect(e.url, nats.UserInfo(user, e.pw[user]), nats.Name(user),
		nats.ErrorHandler(func(_ *nats.Conn, _ *nats.Subscription, err error) {
			e.errMu.Lock()
			e.errs[user] = append(e.errs[user], err)
			e.errMu.Unlock()
		}))
	if err != nil {
		e.t.Fatalf("connect %s: %v", user, err)
	}
	e.conns[user] = c
	return c
}

func (e *env) permErrors(user string) []error {
	e.errMu.Lock()
	defer e.errMu.Unlock()
	return append([]error(nil), e.errs[user]...)
}

func (e *env) setOnline(id string, online bool) {
	e.entMu.Lock()
	defer e.entMu.Unlock()
	c := "online"
	if !online {
		c = "offline"
	}
	e.entities[id]["health"] = map[string]any{"connection": c}
}

func waitFor(t *testing.T, cond func() bool, what string) {
	t.Helper()
	deadline := time.Now().Add(5 * time.Second)
	for time.Now().Before(deadline) {
		if cond() {
			return
		}
		time.Sleep(10 * time.Millisecond)
	}
	t.Fatalf("timed out waiting for %s", what)
}

// handle is a minimal H6-style bridge: verify signature, expiry, nonce, simulated, own id;
// then report acknowledged and executed (or nothing when silent).
func (b *bridgeStub) handle(nc *nats.Conn, m *nats.Msg) {
	var envl map[string]any
	_ = json.Unmarshal(m.Data, &envl)
	parts := strings.Split(m.Subject, ".")
	asset := parts[3]
	b.mu.Lock()
	defer b.mu.Unlock()
	report := func(outcome, code string) {
		data := map[string]any{"action_id": envl["action_id"], "command_id": envl["id"], "outcome": outcome, "simulated": true}
		if code != "" {
			data["code"] = code
		}
		out, _ := json.Marshal(map[string]any{"v": 1, "id": NewUUID(), "ts": ts(b.clock.Now()), "type": "evt.drone.command",
			"source": "bridge-stub", "subject": "evt.drone." + asset + ".command", "data": data})
		_ = nc.Publish("evt.drone."+asset+".command", out)
	}
	id, _ := envl["id"].(string)
	exp, _ := time.Parse(time.RFC3339Nano, fmt.Sprint(envl["expires_at"]))
	switch {
	case VerifyCommand(b.pub, m.Data) != nil:
		b.refused = append(b.refused, "bad_signature")
		report("refused", "bad_signature")
		return
	case b.seen[id]:
		b.refused = append(b.refused, CodeReplayed)
		report("refused", CodeReplayed)
		return
	case !b.clock.Now().Before(exp):
		b.refused = append(b.refused, CodeExpired)
		report("refused", CodeExpired)
		return
	case envl["simulated"] != true || !strings.HasPrefix(asset, "sim-"):
		b.refused = append(b.refused, CodeNotSimulated)
		report("refused", CodeNotSimulated)
		return
	}
	b.seen[id] = true
	b.received = append(b.received, envl)
	if b.silent[asset] {
		return
	}
	report("acknowledged", "")
	report("executed", "")
}

func (b *bridgeStub) count() int { b.mu.Lock(); defer b.mu.Unlock(); return len(b.received) }

// geofenceDown silences the geofence service stub.
func (e *env) geofenceDown(down bool) {
	e.taskMu.Lock()
	e.geoDown = down
	e.taskMu.Unlock()
}
