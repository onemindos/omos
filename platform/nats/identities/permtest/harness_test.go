// Package permtest runs the per-service NATS identities (platform/nats/identities/identities.yaml)
// on disposable in-process nats-servers and checks that every user can use exactly its
// subjects. The config under test is rendered by scripts/nats_identities.py, the same renderer
// the owner uses for the live ConfigMap, with test-only random passwords written outside the
// repository (NATS_IDENTITY_WORKDIR, default a temp dir). nats-server is pinned in go.mod to
// the live version (2.14.2, platform/nats/nats-statefulset.yaml).
//
// NATS_IDENTITY_HASH=1 renders bcrypt hashes (as live) instead of plaintext passwords; it needs
// the nats CLI on PATH. PYTHON selects the interpreter (needs PyYAML).
package permtest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/nats-io/nats-server/v2/server"
	"github.com/nats-io/nats.go"
)

var checks atomic.Int64

func TestMain(m *testing.M) {
	code := m.Run()
	fmt.Printf("permtest: %d permission checks passed in this run (version %s)\n", checks.Load(), server.VERSION)
	os.Exit(code)
}

func pass(t *testing.T, format string, a ...any) {
	t.Helper()
	checks.Add(1)
	t.Logf("PASS "+format, a...)
}

// ---------------------------------------------------------------- config rendering

var (
	setupOnce sync.Once
	setupErr  error
	workDir   string
	passwords map[string]string
	repoRoot  string
)

func python() string {
	if p := os.Getenv("PYTHON"); p != "" {
		return p
	}
	return "python3"
}

func setup(t *testing.T) {
	t.Helper()
	setupOnce.Do(func() {
		wd, _ := os.Getwd()
		for d := wd; d != "/"; d = filepath.Dir(d) {
			if _, err := os.Stat(filepath.Join(d, "scripts/nats_identities.py")); err == nil {
				repoRoot = d
				break
			}
		}
		if repoRoot == "" {
			setupErr = fmt.Errorf("scripts/nats_identities.py not found above %s", wd)
			return
		}
		workDir = os.Getenv("NATS_IDENTITY_WORKDIR")
		if workDir == "" {
			workDir, setupErr = os.MkdirTemp("", "a19-permtest-")
			if setupErr != nil {
				return
			}
		}
		if setupErr = os.MkdirAll(workDir, 0o700); setupErr != nil {
			return
		}
		// Fresh passwords per run: drop configs rendered by an earlier run.
		old, _ := filepath.Glob(filepath.Join(workDir, "fragment-*.json"))
		for _, f := range old {
			os.Remove(f)
		}
		pwFile := filepath.Join(workDir, "passwords.json")
		if out, err := exec.Command(python(), filepath.Join(repoRoot, "scripts/nats_identities.py"),
			"gen-passwords", "--out", pwFile).CombinedOutput(); err != nil {
			setupErr = fmt.Errorf("gen-passwords: %v: %s", err, out)
			return
		}
		b, err := os.ReadFile(pwFile)
		if err != nil {
			setupErr = err
			return
		}
		setupErr = json.Unmarshal(b, &passwords)
	})
	if setupErr != nil {
		t.Fatal(setupErr)
	}
}

// renderStage returns the rendered nats.conf fragment for a stage (system_account, accounts and,
// at A1, leafnodes.authorization).
func renderStage(t *testing.T, stage string) map[string]any {
	t.Helper()
	setup(t)
	args := []string{filepath.Join(repoRoot, "scripts/nats_identities.py"), "conf",
		"--passwords", filepath.Join(workDir, "passwords.json"), "--stage", stage}
	name := "fragment-" + stage
	if os.Getenv("NATS_IDENTITY_HASH") == "1" {
		args = append(args, "--hash")
		name += "-bcrypt"
	}
	out := filepath.Join(workDir, name+".json")
	if _, err := os.Stat(out); err != nil {
		args = append(args, "--out", out)
		if b, err := exec.Command(python(), args...).CombinedOutput(); err != nil {
			t.Fatalf("render %s: %v: %s", stage, err, b)
		}
	}
	b, err := os.ReadFile(out)
	if err != nil {
		t.Fatal(err)
	}
	var m map[string]any
	if err := json.Unmarshal(b, &m); err != nil {
		t.Fatal(err)
	}
	return m
}

// liveLike is the shape of today's live config (nats-config ConfigMap read 2026-10-29 with
// passwords redacted): accounts ONEMIND{onemind} and SYS{sys, monitor}, leaf_tak in
// leafnodes.authorization, nothing restricted.
// The legacy users carry exactly the password values the A1 render gives them (plaintext, or
// bcrypt hashes with NATS_IDENTITY_HASH=1), as `nats_identities.py live` keeps the live hashes.
func liveLike(t *testing.T) map[string]any {
	setup(t)
	a1 := renderStage(t, "A1")
	pw := map[string]string{}
	for _, acc := range a1["accounts"].(map[string]any) {
		for _, u := range acc.(map[string]any)["users"].([]any) {
			m := u.(map[string]any)
			pw[m["user"].(string)] = m["password"].(string)
		}
	}
	for _, u := range a1["leafnodes"].(map[string]any)["authorization"].(map[string]any)["users"].([]any) {
		m := u.(map[string]any)
		pw[m["user"].(string)] = m["password"].(string)
	}
	passwords := pw
	return map[string]any{
		"system_account": "SYS",
		"accounts": map[string]any{
			"ONEMIND": map[string]any{
				"jetstream": map[string]any{"max_mem": -1, "max_file": -1, "max_streams": -1, "max_consumers": -1},
				"users":     []any{map[string]any{"user": "onemind", "password": passwords["onemind"]}},
				"imports":   []any{map[string]any{"service": map[string]any{"account": "SYS", "subject": "$SYS.REQ.>"}}},
			},
			"SYS": map[string]any{
				"users": []any{
					map[string]any{"user": "sys", "password": passwords["sys"]},
					map[string]any{"user": "monitor", "password": passwords["monitor"]},
				},
				"exports": []any{map[string]any{"service": "$SYS.REQ.>", "accounts": []any{"ONEMIND"}}},
			},
		},
		"leafnodes": map[string]any{"authorization": map[string]any{"users": []any{
			map[string]any{"user": "leaf_tak", "password": passwords["leaf_tak"], "account": "ONEMIND"},
		}}},
	}
}

// ---------------------------------------------------------------- servers

type logBuf struct {
	mu    sync.Mutex
	lines []string
}

func (l *logBuf) add(lvl, f string, a ...any) {
	l.mu.Lock()
	l.lines = append(l.lines, lvl+" "+fmt.Sprintf(f, a...))
	l.mu.Unlock()
}
func (l *logBuf) Noticef(f string, a ...any) { l.add("[INF]", f, a...) }
func (l *logBuf) Warnf(f string, a ...any)   { l.add("[WRN]", f, a...) }
func (l *logBuf) Fatalf(f string, a ...any)  { l.add("[FTL]", f, a...) }
func (l *logBuf) Errorf(f string, a ...any)  { l.add("[ERR]", f, a...) }
func (l *logBuf) Debugf(f string, a ...any)  { l.add("[DBG]", f, a...) }
func (l *logBuf) Tracef(f string, a ...any)  { l.add("[TRC]", f, a...) }
func (l *logBuf) grep(s string) []string {
	l.mu.Lock()
	defer l.mu.Unlock()
	var out []string
	for _, x := range l.lines {
		if strings.Contains(x, s) {
			out = append(out, x)
		}
	}
	return out
}

// marshalConf writes JSON the nats-server config parser accepts (no \u003e escapes).
func marshalConf(v any) []byte {
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", "  ")
	_ = enc.Encode(v)
	return buf.Bytes()
}

func freePort(t *testing.T) int {
	t.Helper()
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer l.Close()
	return l.Addr().(*net.TCPAddr).Port
}

type hub struct {
	t                      *testing.T
	stage                  string
	port, leafPort, wsPort int
	storeDir, confPath     string
	srv                    *server.Server
	log                    *logBuf
}

func (h *hub) url() string   { return fmt.Sprintf("nats://127.0.0.1:%d", h.port) }
func (h *hub) wsURL() string { return fmt.Sprintf("ws://127.0.0.1:%d", h.wsPort) }

// writeConf writes the hub config for a fragment: the rendered accounts plus local listeners.
func (h *hub) writeConf(frag map[string]any) {
	conf := map[string]any{
		"server_name": "hub",
		"host":        "127.0.0.1",
		"port":        h.port,
		"jetstream":   map[string]any{"store_dir": h.storeDir, "max_file_store": 256 << 20, "max_memory_store": 0},
		"websocket":   map[string]any{"host": "127.0.0.1", "port": h.wsPort, "no_tls": true},
	}
	leaf := map[string]any{"host": "127.0.0.1", "port": h.leafPort}
	for k, v := range frag {
		if k == "leafnodes" {
			for lk, lv := range v.(map[string]any) {
				leaf[lk] = lv
			}
			continue
		}
		conf[k] = v
	}
	conf["leafnodes"] = leaf
	b := marshalConf(conf)
	if err := os.WriteFile(h.confPath, b, 0o600); err != nil {
		h.t.Fatal(err)
	}
}

func (h *hub) start() {
	h.t.Helper()
	opts, err := server.ProcessConfigFile(h.confPath)
	if err != nil {
		h.t.Fatalf("config %s rejected: %v", h.stage, err)
	}
	opts.NoSigs = true
	srv, err := server.NewServer(opts)
	if err != nil {
		h.t.Fatalf("server: %v", err)
	}
	srv.SetLogger(h.log, false, false)
	go srv.Start()
	if !srv.ReadyForConnections(10 * time.Second) {
		h.t.Fatal("hub not ready")
	}
	h.srv = srv
}

func (h *hub) stop() {
	if h.srv != nil {
		h.srv.Shutdown()
		h.srv.WaitForShutdown()
		h.srv = nil
	}
}

// reload rewrites the config and asks the running server to reload it (what the reloader
// sidecar does after a ConfigMap change).
func (h *hub) reload(frag map[string]any) error {
	h.writeConf(frag)
	return h.srv.Reload()
}

func newHub(t *testing.T, stage string, frag map[string]any) *hub {
	t.Helper()
	setup(t)
	dir := t.TempDir()
	h := &hub{t: t, stage: stage, port: freePort(t), leafPort: freePort(t), wsPort: freePort(t),
		storeDir: filepath.Join(dir, "js"), confPath: filepath.Join(workDir, fmt.Sprintf("hub-%s-%d.conf", stage, time.Now().UnixNano())), log: &logBuf{}}
	h.writeConf(frag)
	h.start()
	t.Cleanup(func() { h.stop(); os.Remove(h.confPath) })
	return h
}

func startHub(t *testing.T, stage string) *hub {
	return newHub(t, stage, renderStage(t, stage))
}

// leaf is a local leaf server whose remote authenticates to the hub as one hub user; local
// clients connect without credentials (like the Hermes sidecar and the TAK leaf).
type leaf struct {
	srv  *server.Server
	port int
	log  *logBuf
}

func (l *leaf) url() string { return fmt.Sprintf("nats://127.0.0.1:%d", l.port) }

func (h *hub) startLeaf(t *testing.T, user string) *leaf {
	t.Helper()
	l := &leaf{port: freePort(t), log: &logBuf{}}
	conf := map[string]any{
		"server_name":  "leaf-" + user,
		"host":         "127.0.0.1",
		"port":         l.port,
		"accounts":     map[string]any{"LOCAL": map[string]any{"users": []any{map[string]any{"user": "local"}}}},
		"no_auth_user": "local",
		"leafnodes": map[string]any{
			"reconnect": "250ms",
			"remotes": []any{map[string]any{
				"url":     fmt.Sprintf("nats://%s:%s@127.0.0.1:%d", user, passwords[user], h.leafPort),
				"account": "LOCAL",
			}},
		},
	}
	path := filepath.Join(workDir, fmt.Sprintf("leaf-%s-%d.conf", user, time.Now().UnixNano()))
	b := marshalConf(conf)
	if err := os.WriteFile(path, b, 0o600); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { os.Remove(path) })
	opts, err := server.ProcessConfigFile(path)
	if err != nil {
		t.Fatal(err)
	}
	opts.NoSigs = true
	srv, err := server.NewServer(opts)
	if err != nil {
		t.Fatal(err)
	}
	srv.SetLogger(l.log, false, false)
	go srv.Start()
	if !srv.ReadyForConnections(10 * time.Second) {
		t.Fatal("leaf not ready")
	}
	l.srv = srv
	t.Cleanup(func() { srv.Shutdown(); srv.WaitForShutdown() })
	return l
}

func waitLeafs(t *testing.T, s *server.Server, n int, d time.Duration) bool {
	t.Helper()
	deadline := time.Now().Add(d)
	for time.Now().Before(deadline) {
		if s.NumLeafNodes() >= n {
			return true
		}
		time.Sleep(50 * time.Millisecond)
	}
	return false
}

// ---------------------------------------------------------------- clients

type sink struct {
	mu   sync.Mutex
	errs []string
}

func (s *sink) add(err error) {
	s.mu.Lock()
	s.errs = append(s.errs, err.Error())
	s.mu.Unlock()
}

// violation waits for a server permission violation for kind ("publish"/"subscription") on subj.
func (s *sink) violation(kind, subj string, d time.Duration) (string, bool) {
	want := strings.ToLower(fmt.Sprintf("permissions violation for %s to %q", kind, subj))
	deadline := time.Now().Add(d)
	for {
		s.mu.Lock()
		for _, e := range s.errs {
			if strings.Contains(strings.ToLower(e), want) {
				s.mu.Unlock()
				return e, true
			}
		}
		s.mu.Unlock()
		if time.Now().After(deadline) {
			return "", false
		}
		time.Sleep(20 * time.Millisecond)
	}
}

type client struct {
	user string
	nc   *nats.Conn
	sink *sink
}

func dial(t *testing.T, url, user string, extra ...nats.Option) (*client, error) {
	t.Helper()
	sk := &sink{}
	opts := append([]nats.Option{
		nats.Name("permtest-" + user),
		nats.NoReconnect(),
		nats.ErrorHandler(func(_ *nats.Conn, _ *nats.Subscription, err error) { sk.add(err) }),
	}, extra...)
	if user != "" {
		opts = append(opts, nats.UserInfo(user, passwords[user]))
	}
	nc, err := nats.Connect(url, opts...)
	if err != nil {
		return nil, err
	}
	t.Cleanup(nc.Close)
	return &client{user: user, nc: nc, sink: sk}, nil
}

func (h *hub) conn(t *testing.T, user string) *client {
	t.Helper()
	c, err := dial(t, h.url(), user)
	if err != nil {
		t.Fatalf("connect %s: %v", user, err)
	}
	return c
}

func (l *leaf) conn(t *testing.T, name string) *client {
	t.Helper()
	c, err := dial(t, l.url(), "", nats.Name("permtest-leafclient-"+name))
	if err != nil {
		t.Fatalf("connect leaf client: %v", err)
	}
	c.user = name
	return c
}

const settle = 150 * time.Millisecond

// expectPub publishes subj as c. If obs is set it must (allowed) or must not (denied) receive
// the message. A denied publish must produce a server permission violation.
func expectPub(t *testing.T, c *client, subj string, allowed bool, obs *client) {
	t.Helper()
	var sub *nats.Subscription
	if obs != nil {
		var err error
		if sub, err = obs.nc.SubscribeSync(subj); err != nil {
			t.Fatalf("observer %s subscribe %s: %v", obs.user, subj, err)
		}
		defer sub.Unsubscribe()
		if err := obs.nc.Flush(); err != nil {
			t.Fatal(err)
		}
		time.Sleep(20 * time.Millisecond)
	}
	if err := c.nc.Publish(subj, []byte(`{"permtest":true}`)); err != nil {
		t.Fatalf("%s publish %s: %v", c.user, subj, err)
	}
	if err := c.nc.Flush(); err != nil {
		t.Fatalf("%s flush: %v", c.user, err)
	}
	if allowed {
		if obs != nil {
			if _, err := sub.NextMsg(2 * time.Second); err != nil {
				t.Fatalf("%s publish %s: observer %s got nothing (%v)", c.user, subj, obs.user, err)
			}
		}
		if e, bad := c.sink.violation("publish", subj, settle); bad {
			t.Fatalf("%s publish %s: unexpected %s", c.user, subj, e)
		}
		via := "no violation"
		if obs != nil {
			via = "delivered to " + obs.user
		}
		pass(t, "ALLOW pub  %-20s %-40s %s", c.user, subj, via)
		return
	}
	e, ok := c.sink.violation("publish", subj, 2*time.Second)
	if !ok {
		t.Fatalf("%s publish %s: expected a server permission violation", c.user, subj)
	}
	if obs != nil {
		if m, err := sub.NextMsg(settle); err == nil {
			t.Fatalf("%s publish %s: denied but %s received %q", c.user, subj, obs.user, m.Data)
		}
	}
	pass(t, "DENY  pub  %-20s %-40s server: %s", c.user, subj, e)
}

func expectSub(t *testing.T, c *client, subj string, allowed bool) {
	t.Helper()
	sub, err := c.nc.SubscribeSync(subj)
	if err != nil {
		t.Fatalf("%s subscribe %s: %v", c.user, subj, err)
	}
	defer sub.Unsubscribe()
	if err := c.nc.Flush(); err != nil {
		t.Fatal(err)
	}
	if allowed {
		if e, bad := c.sink.violation("subscription", subj, settle); bad {
			t.Fatalf("%s subscribe %s: unexpected %s", c.user, subj, e)
		}
		pass(t, "ALLOW sub  %-20s %-40s no violation", c.user, subj)
		return
	}
	e, ok := c.sink.violation("subscription", subj, 2*time.Second)
	if !ok {
		t.Fatalf("%s subscribe %s: expected a server permission violation", c.user, subj)
	}
	pass(t, "DENY  sub  %-20s %-40s server: %s", c.user, subj, e)
}

// expectNoDelivery publishes from a leaf client (the spoke filters with the permissions the hub
// sent) and checks the hub-side observer receives nothing.
func expectLeafPub(t *testing.T, c *client, subj string, allowed bool, obs *client) {
	t.Helper()
	sub, err := obs.nc.SubscribeSync(subj)
	if err != nil {
		t.Fatal(err)
	}
	defer sub.Unsubscribe()
	obs.nc.Flush()
	time.Sleep(100 * time.Millisecond) // interest reaches the leaf
	if err := c.nc.Publish(subj, []byte(`{"permtest":true}`)); err != nil {
		t.Fatal(err)
	}
	c.nc.Flush()
	if allowed {
		if _, err := sub.NextMsg(2 * time.Second); err != nil {
			t.Fatalf("leaf %s publish %s: hub observer %s got nothing", c.user, subj, obs.user)
		}
		pass(t, "ALLOW leaf-pub %-17s %-40s delivered to hub %s", c.user, subj, obs.user)
		return
	}
	if m, err := sub.NextMsg(500 * time.Millisecond); err == nil {
		t.Fatalf("leaf %s publish %s: denied but hub %s received %q", c.user, subj, obs.user, m.Data)
	}
	pass(t, "DENY  leaf-pub %-17s %-40s not delivered to hub %s", c.user, subj, obs.user)
}

// expectLeafSub checks whether a message published on the hub by pub reaches a leaf subscriber.
func expectLeafSub(t *testing.T, c *client, subj string, allowed bool, pub *client) {
	t.Helper()
	sub, err := c.nc.SubscribeSync(subj)
	if err != nil {
		t.Fatal(err)
	}
	defer sub.Unsubscribe()
	c.nc.Flush()
	time.Sleep(150 * time.Millisecond)
	if err := pub.nc.Publish(subj, []byte(`{"permtest":true}`)); err != nil {
		t.Fatal(err)
	}
	pub.nc.Flush()
	if allowed {
		if _, err := sub.NextMsg(2 * time.Second); err != nil {
			t.Fatalf("leaf %s: %s published by %s did not arrive", c.user, subj, pub.user)
		}
		pass(t, "ALLOW leaf-sub %-17s %-40s from hub %s", c.user, subj, pub.user)
		return
	}
	if m, err := sub.NextMsg(500 * time.Millisecond); err == nil {
		t.Fatalf("leaf %s: %s published by %s arrived %q but must be filtered", c.user, subj, pub.user, m.Data)
	}
	pass(t, "DENY  leaf-sub %-17s %-40s from hub %s filtered", c.user, subj, pub.user)
}
