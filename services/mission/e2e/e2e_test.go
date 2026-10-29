//go:build e2e

// End-to-end test of the task service with the real policy service (omos #554), the real
// MAVLink bridge (omos #550) and the onemind-xr scripted drone, on a disposable
// nats-server with the A19 account layout (ONEMIND, CLIENT, AGENTS). Run through
// test/e2e.sh, which builds the binaries and sets the E2E_* variables.
//
// Flow (mission-task.md section 6.1 and section 10 tests 1, 3 and 6):
//  1. two humans arm and take off sim-x1 (A4) so it can fly;
//  2. zeus creates a mission and two waypoints (A2), each closed by the task service's
//     result (POL-14: causation = command id);
//  3. Hermes proposes a goto task for sim-x1; nothing happens until zeus confirms; the
//     drone flies; the task goes created, sent, received, executing, done_ok in KV
//     mission-state and on evt.mission.{tsk}.*; the mission completes;
//  4. an orbit task for an asset without that capability is denied not_capable;
//  5. a goto for sim-x2 (no bridge) times out at delivery.before_t and is not sent again.
package e2e

import (
	"context"
	"crypto/ed25519"
	"crypto/rand"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

const ulidChars = "0123456789abcdefghjkmnpqrstvwxyz"

func newID(prefix string) string {
	b := make([]byte, 26)
	_, _ = rand.Read(b)
	for i := range b {
		b[i] = ulidChars[int(b[i])%32]
	}
	return prefix + "-" + string(b)
}

func uuid4() string {
	var u [16]byte
	_, _ = rand.Read(u[:])
	u[6] = u[6]&0x0f | 0x40
	u[8] = u[8]&0x3f | 0x80
	return fmt.Sprintf("%x-%x-%x-%x-%x", u[0:4], u[4:6], u[6:8], u[8:10], u[10:16])
}

// envPair builds one KEY=value entry for a child process environment.
func envPair(k, v string) string { return k + "=" + v }

func mustEnv(t *testing.T, k string) string {
	v := os.Getenv(k)
	if v == "" {
		t.Skipf("%s not set; run test/e2e.sh", k)
	}
	return v
}

func freePort(t *testing.T, network string) int {
	t.Helper()
	if network == "udp" {
		c, err := net.ListenPacket("udp", "127.0.0.1:0")
		if err != nil {
			t.Fatal(err)
		}
		defer c.Close()
		return c.LocalAddr().(*net.UDPAddr).Port
	}
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer l.Close()
	return l.Addr().(*net.TCPAddr).Port
}

type transcript struct {
	mu sync.Mutex
	f  *os.File
	t  *testing.T
}

func (tr *transcript) log(kind, subject string, data any) {
	tr.mu.Lock()
	defer tr.mu.Unlock()
	line, _ := json.Marshal(map[string]any{"at": time.Now().UTC().Format(time.RFC3339Nano), "kind": kind, "subject": subject, "data": data})
	if tr.f != nil {
		fmt.Fprintln(tr.f, string(line))
	}
	if kind == "step" || kind == "check" {
		tr.t.Log(string(line))
	}
}

type bus struct {
	mu   sync.Mutex
	msgs []*nats.Msg
}

func (b *bus) add(m *nats.Msg) { b.mu.Lock(); b.msgs = append(b.msgs, m); b.mu.Unlock() }
func (b *bus) all() []*nats.Msg {
	b.mu.Lock()
	defer b.mu.Unlock()
	return append([]*nats.Msg(nil), b.msgs...)
}

func decode(m *nats.Msg) map[string]any {
	var v map[string]any
	_ = json.Unmarshal(m.Data, &v)
	return v
}

func waitFor(t *testing.T, d time.Duration, what string, cond func() bool) {
	t.Helper()
	deadline := time.Now().Add(d)
	for time.Now().Before(deadline) {
		if cond() {
			return
		}
		time.Sleep(50 * time.Millisecond)
	}
	t.Fatalf("timed out after %s waiting for %s", d, what)
}

func natsConf(port int, store string, pw map[string]string) string {
	return fmt.Sprintf(`
listen: "127.0.0.1:%d"
jetstream { store_dir: %q, max_mem: 64MB, max_file: 512MB }
accounts {
  ONEMIND {
    jetstream: enabled
    users: [
      { user: admin, password: %q }
      { user: policy, password: %q, permissions: {
          publish: { allow: ["cmd.>", "evt.so.action.>", "svc.entity.get", "svc.mission.apply", "$JS.API.>", "_INBOX.>"] }
          subscribe: { allow: ["svc.policy.>", "evt.drone.>", "evt.robot.>", "evt.vehicle.>", "evt.mission.>", "evt.entity.>", "_INBOX.>"] }
          allow_responses: true } }
      { user: mission, password: %q, permissions: {
          publish: { allow: ["evt.mission.>", "$JS.API.INFO", "$JS.API.STREAM.INFO.*", "$JS.API.CONSUMER.CREATE.CMD.>", "$JS.API.CONSUMER.CREATE.EVT.>", "$JS.API.CONSUMER.INFO.CMD.mission-commands", "$JS.API.CONSUMER.INFO.EVT.mission-events", "$JS.API.CONSUMER.MSG.NEXT.CMD.mission-commands", "$JS.API.CONSUMER.MSG.NEXT.EVT.mission-events", "$JS.ACK.CMD.mission-commands.>", "$JS.ACK.EVT.mission-events.>", "$JS.FC.>", "$JS.API.DIRECT.GET.KV_mission-state.>", "$JS.API.STREAM.MSG.GET.KV_mission-state", "$JS.API.CONSUMER.CREATE.KV_mission-state", "$JS.API.CONSUMER.CREATE.KV_mission-state.>", "$JS.API.CONSUMER.DELETE.KV_mission-state.>", "$KV.mission-state.>"] }
          subscribe: { allow: ["svc.mission.>", "tel.fleet.*.*.position", "_INBOX.>"] }
          allow_responses: true } }
      { user: bridge, password: %q, permissions: {
          publish: { allow: ["tel.fleet.drone.sim-x1.>", "evt.drone.sim-x1.>"] }
          subscribe: { allow: ["cmd.fleet.drone.sim-x1.>", "_INBOX.>"] } } }
      { user: entity, password: %q, permissions: {
          publish: { allow: ["_INBOX.>"] }, subscribe: { allow: ["svc.entity.>"] }, allow_responses: true } }
    ]
    exports: [
      { service: "svc.policy.>" }
      { service: "svc.mission.get" }, { service: "svc.mission.list" }, { service: "svc.mission.report" }
    ]
  }
  CLIENT {
    users: [
      { user: zeus, password: %q, permissions: { publish: { allow: ["svc.policy.>", "svc.mission.>"] }, subscribe: { allow: ["_INBOX.>"] } } }
      { user: alice, password: %q, permissions: { publish: { allow: ["svc.policy.>", "svc.mission.>"] }, subscribe: { allow: ["_INBOX.>"] } } }
    ]
    imports: [
      { service: { account: ONEMIND, subject: "svc.policy.>" }, share: true }
      { service: { account: ONEMIND, subject: "svc.mission.get" }, share: true }
      { service: { account: ONEMIND, subject: "svc.mission.list" }, share: true }
      { service: { account: ONEMIND, subject: "svc.mission.report" }, share: true }
    ]
  }
  AGENTS {
    users: [
      { user: hermes, password: %q, permissions: {
          publish: { allow: ["svc.policy.command.propose", "svc.policy.task.propose", "svc.policy.action.get", "svc.mission.get", "svc.mission.list"] }
          subscribe: { allow: ["_INBOX.>"] } } }
    ]
    imports: [
      { service: { account: ONEMIND, subject: "svc.policy.command.propose" }, share: true }
      { service: { account: ONEMIND, subject: "svc.policy.task.propose" }, share: true }
      { service: { account: ONEMIND, subject: "svc.policy.action.get" }, share: true }
      { service: { account: ONEMIND, subject: "svc.mission.get" }, share: true }
      { service: { account: ONEMIND, subject: "svc.mission.list" }, share: true }
    ]
  }
}
`, port, store, pw["admin"], pw["policy"], pw["mission"], pw["bridge"], pw["entity"], pw["zeus"], pw["alice"], pw["hermes"])
}

type harness struct {
	t     *testing.T
	url   string
	pw    map[string]string
	conns map[string]*nats.Conn
	tr    *transcript
	all   *bus
}

func (h *harness) conn(user string) *nats.Conn {
	if c, ok := h.conns[user]; ok {
		return c
	}
	c, err := nats.Connect(h.url, nats.UserInfo(user, h.pw[user]), nats.Name("e2e-"+user))
	if err != nil {
		h.t.Fatalf("connect %s: %v", user, err)
	}
	h.conns[user] = c
	return c
}

func (h *harness) request(user, subject string, body any) map[string]any {
	h.t.Helper()
	b, _ := json.Marshal(body)
	m, err := h.conn(user).Request(subject, b, 5*time.Second)
	if err != nil {
		h.t.Fatalf("%s %s: %v", user, subject, err)
	}
	v := decode(m)
	h.tr.log("rpc", subject, map[string]any{"user": user, "request": body, "reply": v})
	return v
}

func actionReq(actionType, entityType, id string, params map[string]any) map[string]any {
	return map[string]any{"v": 1, "request_id": uuid4(), "action_type": actionType, "action_version": 1,
		"target": map[string]any{"entity_type": entityType, "entity_id": id}, "params": params,
		"requested_by": map[string]any{"kind": "user", "id": "user:claimed"}, "reason": "e2e"}
}

func family(actionType string) string {
	if strings.HasPrefix(actionType, "task.") || strings.HasPrefix(actionType, "mission.") {
		return "task"
	}
	return "command"
}

// approve proposes as proposer and confirms with each approver; it returns the action id.
func (h *harness) approve(proposer string, body map[string]any, approvers ...string) string {
	h.t.Helper()
	fam := family(body["action_type"].(string))
	r := h.request(proposer, "svc.policy."+fam+".propose", body)
	if r["status"] != "awaiting_confirmation" {
		h.t.Fatalf("propose %s: %v", body["action_type"], r)
	}
	id := r["action_id"].(string)
	for _, a := range approvers {
		c := h.request(a, "svc.policy."+fam+".confirm", map[string]any{"v": 1, "action_id": id, "params_hash": r["params_hash"], "decision": "approve"})
		if c["denial"] != nil {
			h.t.Fatalf("confirm %s by %s: %v", body["action_type"], a, c)
		}
	}
	return id
}

func (h *harness) audit(actionID, event string) bool {
	for _, m := range h.all.all() {
		if m.Subject == "evt.so.action."+event {
			if d, _ := decode(m)["data"].(map[string]any); d["action_id"] == actionID {
				return true
			}
		}
	}
	return false
}

func (h *harness) auditCode(actionID, event string) string {
	for _, m := range h.all.all() {
		if m.Subject == "evt.so.action."+event {
			if d, _ := decode(m)["data"].(map[string]any); d["action_id"] == actionID {
				c, _ := d["code"].(string)
				return c
			}
		}
	}
	return ""
}

func (h *harness) facts(id string) []string {
	var out []string
	for _, m := range h.all.all() {
		if strings.HasPrefix(m.Subject, "evt.mission."+id+".") {
			out = append(out, strings.TrimPrefix(m.Subject, "evt.mission."+id+"."))
		}
	}
	return out
}

func (h *harness) count(prefix string) int {
	n := 0
	for _, m := range h.all.all() {
		if strings.HasPrefix(m.Subject, prefix) {
			n++
		}
	}
	return n
}

func TestE2E(t *testing.T) {
	natsBin := mustEnv(t, "E2E_NATS_SERVER")
	policyBin := mustEnv(t, "E2E_POLICY_BIN")
	bridgeBin := mustEnv(t, "E2E_BRIDGE_BIN")
	missionBin := mustEnv(t, "E2E_MISSION_BIN")
	schemas := mustEnv(t, "E2E_POLICY_SCHEMAS")
	model := mustEnv(t, "E2E_POLICY_MODEL")
	sim := mustEnv(t, "E2E_SIM_DRONE")
	python := mustEnv(t, "E2E_PYTHON")
	work := mustEnv(t, "E2E_WORKDIR")

	tr := &transcript{t: t}
	if p := os.Getenv("E2E_TRANSCRIPT"); p != "" {
		f, err := os.Create(p)
		if err != nil {
			t.Fatal(err)
		}
		defer f.Close()
		tr.f = f
	}
	h := &harness{t: t, pw: map[string]string{}, conns: map[string]*nats.Conn{}, tr: tr, all: &bus{}}
	for _, u := range []string{"admin", "policy", "mission", "bridge", "entity", "zeus", "alice", "hermes"} {
		b := make([]byte, 18)
		_, _ = rand.Read(b)
		h.pw[u] = hex.EncodeToString(b)
	}

	// Processes, stopped in reverse order at the end.
	var procs []*exec.Cmd
	start := func(name string, env []string, bin string, args ...string) {
		t.Helper()
		logf, err := os.Create(filepath.Join(work, name+".log"))
		if err != nil {
			t.Fatal(err)
		}
		c := exec.Command(bin, args...)
		c.Env, c.Stdout, c.Stderr = append(os.Environ(), env...), logf, logf
		if err := c.Start(); err != nil {
			t.Fatalf("start %s: %v", name, err)
		}
		procs = append(procs, c)
		tr.log("step", "", map[string]any{"started": name, "pid": c.Process.Pid})
	}
	defer func() {
		for i := len(procs) - 1; i >= 0; i-- {
			_ = procs[i].Process.Signal(os.Interrupt)
			done := make(chan struct{})
			go func(c *exec.Cmd) { _ = c.Wait(); close(done) }(procs[i])
			select {
			case <-done:
			case <-time.After(5 * time.Second):
				_ = procs[i].Process.Kill()
			}
		}
		for _, c := range h.conns {
			c.Close()
		}
	}()

	// 1. nats-server with the A19 account layout; passwords for this run only.
	port := freePort(t, "tcp")
	confPath := filepath.Join(work, "nats.conf")
	if err := os.WriteFile(confPath, []byte(natsConf(port, filepath.Join(work, "js"), h.pw)), 0o600); err != nil {
		t.Fatal(err)
	}
	start("nats-server", nil, natsBin, "-c", confPath)
	h.url = fmt.Sprintf("nats://127.0.0.1:%d", port)
	waitFor(t, 10*time.Second, "nats-server", func() bool {
		c, err := nats.Connect(h.url, nats.UserInfo("admin", h.pw["admin"]))
		if err == nil {
			c.Close()
		}
		return err == nil
	})

	// 2. Streams and buckets as platform/nats/streams (CMD, EVT, MISSION, TEL, KV mission-state, KV geofences).
	ctx := context.Background()
	admin := h.conn("admin")
	js, _ := jetstream.New(admin)
	for name, subj := range map[string][]string{"CMD": {"cmd.>"}, "EVT": {"evt.so.>", "evt.drone.>", "evt.robot.>"},
		"MISSION": {"evt.mission.>"}, "TEL": {"tel.fleet.>"}} {
		if _, err := js.CreateStream(ctx, jetstream.StreamConfig{Name: name, Subjects: subj, Storage: jetstream.FileStorage, Duplicates: 2 * time.Minute}); err != nil {
			t.Fatal(err)
		}
	}
	if _, err := js.CreateKeyValue(ctx, jetstream.KeyValueConfig{Bucket: "mission-state", History: 5}); err != nil {
		t.Fatal(err)
	}
	geo, err := js.CreateKeyValue(ctx, jetstream.KeyValueConfig{Bucket: "geofences", History: 10})
	if err != nil {
		t.Fatal(err)
	}
	// keep_in: 20 km around the scripted drone's home point.
	if _, err := geo.Put(ctx, "zone-home", []byte(`{"entity_id":"zone-home","entity_type":"zone","live":true,"geo":{"shape":"circle","kind":"keep_in","radius_m":20000,"coords":{"type":"Point","coordinates":[-97.7431,30.2672]}}}`)); err != nil {
		t.Fatal(err)
	}
	mstate, _ := js.KeyValue(ctx, "mission-state")

	// Recorder: every fact, audit, command, outcome and task report (positions sampled).
	var posMu sync.Mutex
	var lastPos map[string]any
	var posN int
	for _, s := range []string{"evt.mission.>", "evt.so.action.>", "cmd.>", "evt.drone.>"} {
		if _, err := admin.Subscribe(s, func(m *nats.Msg) { h.all.add(m); tr.log("bus", m.Subject, decode(m)) }); err != nil {
			t.Fatal(err)
		}
	}
	if _, err := admin.Subscribe("tel.fleet.drone.sim-x1.position", func(m *nats.Msg) {
		posMu.Lock()
		lastPos = decode(m)
		posN++
		if posN%10 == 0 {
			tr.log("tel", m.Subject, lastPos)
		}
		posMu.Unlock()
	}); err != nil {
		t.Fatal(err)
	}
	_ = admin.Flush()

	// 3. Entity service stand-in (svc.entity.get; the A4 entity service is #552). sim-x1
	// advertises goto, hold, rtl and land, not orbit (test 3).
	home := map[string]any{"lat": 30.2672, "lon": -97.7431, "hae_m": 34.0}
	ents := map[string]map[string]any{
		"sim-x1": {"entity_id": "sim-x1", "entity_type": "drone", "kind": "multirotor", "location": home,
			"health": map[string]any{"connection": "online"}, "indicators": map[string]any{"simulated": true},
			"capabilities": []string{"onemind.task.goto.v1", "onemind.task.hold.v1", "onemind.task.rtl.v1", "onemind.task.land.v1"}},
		"sim-x2": {"entity_id": "sim-x2", "entity_type": "drone", "kind": "multirotor", "location": home,
			"health": map[string]any{"connection": "online"}, "indicators": map[string]any{"simulated": true},
			"capabilities": []string{"onemind.task.goto.v1"}},
	}
	ec := h.conn("entity")
	if _, err := ec.Subscribe("svc.entity.get", func(m *nats.Msg) {
		var q struct {
			EntityID string `json:"entity_id"`
		}
		_ = json.Unmarshal(m.Data, &q)
		if e, ok := ents[q.EntityID]; ok {
			b, _ := json.Marshal(map[string]any{"entity": e})
			_ = m.Respond(b)
			return
		}
		_ = m.Respond([]byte(`{"error":{"code":"not_found"}}`))
	}); err != nil {
		t.Fatal(err)
	}
	_ = ec.Flush()

	// 4. Policy signing key (Ed25519, generated for this run), principals and tuples.
	pub, priv, _ := ed25519.GenerateKey(rand.Reader)
	der, _ := x509.MarshalPKCS8PrivateKey(priv)
	keyPath := filepath.Join(work, "policy-signing.pem")
	if err := os.WriteFile(keyPath, pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: der}), 0o600); err != nil {
		t.Fatal(err)
	}
	principals := map[string]any{"own_account": "ONEMIND", "agent_accounts": []string{"AGENTS"}, "principals": []map[string]string{
		{"account": "CLIENT", "user": "zeus", "kind": "user", "id": "user:zeus"},
		{"account": "CLIENT", "user": "alice", "kind": "user", "id": "user:alice"},
		{"account": "AGENTS", "user": "hermes", "kind": "agent", "id": "agent:hermes/zeus/legacy"},
	}}
	prPath := filepath.Join(work, "principals.json")
	b, _ := json.Marshal(principals)
	_ = os.WriteFile(prPath, b, 0o600)
	tuples := []map[string]string{
		{"user": "user:zeus", "relation": "operator", "object": "asset:sim-x1"},
		{"user": "user:zeus", "relation": "operator", "object": "asset:sim-x2"},
		{"user": "user:alice", "relation": "operator", "object": "asset:sim-x1"},
	}
	for _, a := range []string{"asset.arm", "asset.takeoff", "asset.rtl", "mission.create", "mission.waypoint", "mission.complete",
		"mission.abort", "task.assign", "task.abort"} {
		obj := "action_type:" + a
		for _, u := range []string{"user:zeus", "user:alice"} {
			for _, r := range []string{"proposer", "approver", "second_approver"} {
				tuples = append(tuples, map[string]string{"user": u, "relation": r, "object": obj})
			}
		}
		tuples = append(tuples, map[string]string{"user": "agent:hermes/zeus/legacy", "relation": "proposer", "object": obj})
	}
	tuPath := filepath.Join(work, "tuples.json")
	b, _ = json.Marshal(tuples)
	_ = os.WriteFile(tuPath, b, 0o600)

	// 5. Scripted drone and the real MAVLink bridge (signatures required).
	udp := freePort(t, "udp")
	start("sim-drone", []string{"MAVLINK20=1"}, python, "-I", sim, "--out", fmt.Sprintf("udpout:127.0.0.1:%d", udp),
		"--lat", "30.2672", "--lon", "-97.7431", "--radius", "40", "--alt", "30", "--speed", "12")
	start("bridge-mavlink", []string{"NATS_URL=" + h.url, "NATS_USER=bridge", "NATS_PASSWORD=" + h.pw["bridge"],
		"ASSET_ID=sim-x1", fmt.Sprintf("MAVLINK_ENDPOINT=udpserver:127.0.0.1:%d", udp),
		"POLICY_PUBKEYS=e2e=" + base64.StdEncoding.EncodeToString(pub), "SIG_REQUIRED=true",
		// The bridge's own fence (H6, it fails closed without one): 2 km around home.
		`GEOFENCE_JSON={"keep_in":[{"id":"sim-range","circle":{"lat":30.2672,"lon":-97.7431,"radius_m":2000}}],"max_hae_m":300}`}, bridgeBin)

	// 6. Policy service (local OpenFGA evaluator) and the task service.
	polHealth, misHealth := freePort(t, "tcp"), freePort(t, "tcp")
	start("policy", []string{"NATS_URL=" + h.url, "NATS_USER=policy", "NATS_PASSWORD=" + h.pw["policy"],
		"POLICY_AUTHZ=local", "POLICY_MODEL_FILE=" + model, "POLICY_TUPLES_FILE=" + tuPath, "POLICY_SCHEMA_DIR=" + schemas,
		"POLICY_PRINCIPALS_FILE=" + prPath, envPair("POLICY_SIGNING_KEY_FILE", keyPath), "POLICY_SIGNING_KID=e2e",
		fmt.Sprintf("HEALTH_ADDR=127.0.0.1:%d", polHealth)}, policyBin)
	start("mission", []string{"NATS_URL=" + h.url, "NATS_USER=mission", "NATS_PASSWORD=" + h.pw["mission"],
		"MISSION_PRINCIPALS_FILE=" + prPath, "MISSION_SWEEP_EVERY=200ms", fmt.Sprintf("HEALTH_ADDR=127.0.0.1:%d", misHealth)}, missionBin)
	for _, p := range []int{polHealth, misHealth} {
		url := fmt.Sprintf("http://127.0.0.1:%d/healthz", p)
		waitFor(t, 20*time.Second, url, func() bool {
			r, err := http.Get(url)
			if err == nil {
				r.Body.Close()
			}
			return err == nil && r.StatusCode == 200
		})
	}
	waitFor(t, 20*time.Second, "sim-x1 telemetry through the bridge", func() bool { posMu.Lock(); defer posMu.Unlock(); return posN > 5 })
	tr.log("step", "", "stack up: nats-server, sim drone, bridge, policy, task service")

	// Step 1: arm and take off (A4: zeus and alice).
	arm := h.approve("zeus", actionReq("asset.arm", "drone", "sim-x1", map[string]any{}), "zeus", "alice")
	waitFor(t, 20*time.Second, "arm executed", func() bool { return h.audit(arm, "executed") })
	takeoff := h.approve("zeus", actionReq("asset.takeoff", "drone", "sim-x1", map[string]any{"hae_m": 4.0}), "zeus", "alice")
	waitFor(t, 20*time.Second, "takeoff executed", func() bool { return h.audit(takeoff, "executed") })
	tr.log("step", "", "sim-x1 armed and airborne")

	// Step 2: mission and waypoints (A2), closed by the task service result (POL-14).
	msn := newID("msn")
	mc := h.approve("zeus", actionReq("mission.create", "mission", msn, map[string]any{"name": "North field check", "intent": "e2e D2/H4"}), "zeus")
	waitFor(t, 15*time.Second, "mission.create executed", func() bool { return h.audit(mc, "executed") })
	for i, wp := range [][2]float64{{30.2690, -97.7431}, {30.2700, -97.7420}} {
		w := h.approve("zeus", actionReq("mission.waypoint", "mission", msn, map[string]any{"seq": i, "lat": wp[0], "lon": wp[1], "hae_m": 60.0, "name": fmt.Sprintf("WP%d", i)}), "zeus")
		waitFor(t, 15*time.Second, "waypoint executed", func() bool { return h.audit(w, "executed") })
	}
	g := h.request("zeus", "svc.mission.get", map[string]any{"v": 1, "id": msn})
	rec, _ := g["record"].(map[string]any)
	if rec == nil || rec["state"] != "planned" || len(rec["route"].([]any)) != 2 {
		t.Fatalf("mission record: %v", g)
	}
	tr.log("check", "", map[string]any{"mission": msn, "state": rec["state"], "route_points": 2, "policy_closed_on_causation": true})

	// Step 3 (contract test 1): Hermes proposes; nothing happens until a human confirms.
	tsk := newID("tsk")
	goal := map[string]any{"lat": 30.2690, "lon": -97.7431, "hae_m": 30.0}
	body := actionReq("task.assign", "drone", "sim-x1", map[string]any{"task_id": tsk, "parent_id": msn,
		"type": "onemind.task.goto.v1", "spec": goal, "description": "fly to WP0"})
	body["requested_by"] = map[string]any{"kind": "agent", "id": "agent:hermes/zeus/legacy", "on_behalf_of": "user:zeus"}
	p := h.request("hermes", "svc.policy.task.propose", body)
	if p["status"] != "awaiting_confirmation" {
		t.Fatalf("hermes propose: %v", p)
	}
	act := p["action_id"].(string)
	time.Sleep(1500 * time.Millisecond)
	if n := h.count("cmd.fleet.drone.sim-x1.goto"); n != 0 {
		t.Fatalf("a goto went out before the human confirm (%d)", n)
	}
	if _, err := mstate.Get(ctx, tsk); err == nil {
		t.Fatal("a task record exists before approval")
	}
	// Hermes cannot reach confirm at all (no publish permission; agent_cannot_confirm in policy).
	if _, err := h.conn("hermes").Request("svc.policy.task.confirm", []byte(`{"v":1}`), time.Second); err == nil {
		t.Fatal("hermes reached svc.policy.task.confirm")
	}
	tr.log("check", "", map[string]any{"waiting_for_human": true, "task_record_before_approval": false, "hermes_confirm": "refused by the server"})
	c := h.request("zeus", "svc.policy.task.confirm", map[string]any{"v": 1, "action_id": act, "params_hash": p["params_hash"], "decision": "approve"})
	if c["status"] != "dispatched" {
		t.Fatalf("confirm: %v", c)
	}
	// With the branches as they are, the bridge (#550) refuses every task-bearing command:
	// it recomputes params_hash as asset.{verb} over the command data, while the policy
	// service (#554) stamps the task.assign hash (owner proposal D2-P1). Then the task must
	// end done_not_ok + rejected, faithfully. With test/interop/policy-task-fleet-params-hash.patch
	// applied to the policy build (E2E_POLICY_PATCHED=1), the drone flies and the task completes.
	patched := os.Getenv("E2E_POLICY_PATCHED") == "1"
	want := "done_ok"
	if !patched {
		want = "done_not_ok"
	}
	var final map[string]any
	waitFor(t, 90*time.Second, "task "+want+" in mission-state", func() bool {
		e, err := mstate.Get(ctx, tsk)
		if err != nil {
			return false
		}
		_ = json.Unmarshal(e.Value(), &final)
		st, _ := final["status"].(map[string]any)
		return st["state"] == want
	})
	posMu.Lock()
	tr.log("check", "", map[string]any{"task": tsk, "final_record": final, "last_position": lastPos})
	posMu.Unlock()
	facts := h.facts(tsk)
	if !patched {
		st := final["status"].(map[string]any)
		er, _ := st["error"].(map[string]any)
		if er["code"] != "rejected" || !strings.Contains(fmt.Sprint(er["message"]), "params_hash") {
			t.Fatalf("unpatched run: want done_not_ok rejected (params_hash), got %v", st)
		}
		tr.log("check", "", map[string]any{"contract_test": 1, "result": "interop blocker reproduced: bridge refused the task goto (params_hash); task service recorded done_not_ok + rejected",
			"task_facts": facts, "blocker": "D2-P1", "pass": "partial"})
	}
	for _, want := range map[bool][]string{true: {"created", "received", "executing", "done_ok"}, false: {"created", "done_not_ok"}}[patched] {
		found := false
		for _, f := range facts {
			found = found || f == want
		}
		if !found {
			t.Fatalf("no evt.mission.%s.%s; facts %v", tsk, want, facts)
		}
	}
	for _, ev := range []string{"proposed", "approved", "dispatched"} {
		if !h.audit(act, ev) {
			t.Fatalf("no evt.so.action.%s for %s", ev, act)
		}
	}
	if patched {
		waitFor(t, 10*time.Second, "mission completed", func() bool {
			for _, f := range h.facts(msn) {
				if f == "completed" {
					return true
				}
			}
			return false
		})
		tr.log("check", "", map[string]any{"contract_test": 1, "task_facts": facts, "mission_facts": h.facts(msn), "pass": true})
	}

	// Step 4 (contract test 3): orbit for an asset without the capability: not_capable.
	t3 := newID("tsk")
	r3 := h.request("zeus", "svc.policy.task.propose", actionReq("task.assign", "drone", "sim-x1", map[string]any{"task_id": t3, "parent_id": nil,
		"type": "onemind.task.orbit.v1", "spec": map[string]any{"center": map[string]any{"lat": 30.27, "lon": -97.74}, "radius_m": 60, "hae_m": 60, "duration_s": 60}}))
	den, _ := r3["denial"].(map[string]any)
	if r3["status"] != "denied" || den["code"] != "not_capable" {
		t.Fatalf("orbit without capability: %v", r3)
	}
	time.Sleep(500 * time.Millisecond)
	if _, err := mstate.Get(ctx, t3); err == nil || h.count("evt.mission."+t3) != 0 {
		t.Fatal("a denied task reached the task service")
	}
	tr.log("check", "", map[string]any{"contract_test": 3, "denial": den, "pass": true})

	// Step 5 (contract test 6): no receipt by delivery.before_t -> done_not_ok timeout, never re-sent.
	t6 := newID("tsk")
	before := time.Now().Add(4 * time.Second).UnixMilli()
	a6 := h.approve("zeus", actionReq("task.assign", "drone", "sim-x2", map[string]any{"task_id": t6, "parent_id": nil,
		"type": "onemind.task.goto.v1", "spec": map[string]any{"lat": 30.2680, "lon": -97.7431, "hae_m": 40.0},
		"delivery": map[string]any{"before_t": before, "require_ack": true}}), "zeus")
	var r6 map[string]any
	waitFor(t, 20*time.Second, "sim-x2 task timeout", func() bool {
		e, err := mstate.Get(ctx, t6)
		if err != nil {
			return false
		}
		_ = json.Unmarshal(e.Value(), &r6)
		st, _ := r6["status"].(map[string]any)
		return st["state"] == "done_not_ok"
	})
	st6 := r6["status"].(map[string]any)
	if st6["error"].(map[string]any)["code"] != "timeout" {
		t.Fatalf("sim-x2 task: %v", st6)
	}
	time.Sleep(12 * time.Second) // past the policy service's 10 s ack window
	if n := h.count("cmd.fleet.drone.sim-x2."); n != 1 {
		t.Fatalf("sim-x2 got %d commands; want exactly 1 (no re-send)", n)
	}
	tr.log("check", "", map[string]any{"contract_test": 6, "status": st6, "fleet_commands_to_sim_x2": 1,
		"policy_outcome": h.auditCode(a6, "failed"), "pass": true})

	// Whole-run invariant (stop-line test 4): no cmd.fleet.> without a human approval.
	approved := map[string]bool{}
	for _, m := range h.all.all() {
		if m.Subject == "evt.so.action.approved" {
			d := decode(m)["data"].(map[string]any)
			approved[d["action_id"].(string)] = true
		}
	}
	for _, m := range h.all.all() {
		if strings.HasPrefix(m.Subject, "cmd.fleet.") {
			if a, _ := decode(m)["action_id"].(string); !approved[a] {
				t.Fatalf("%s for %s without evt.so.action.approved", m.Subject, a)
			}
		}
	}
	tr.log("check", "", map[string]any{"invariant": "every cmd.fleet.> has a human evt.so.action.approved", "fleet_commands": h.count("cmd.fleet."), "pass": true})
}
