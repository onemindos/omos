package projector

// Integration harness: a disposable nats-server (binary from $NATS_SERVER_BIN or PATH; tests
// skip when absent), streams and buckets created with the settings of their CRs
// (platform/nats/streams/ent.yaml, cot.yaml, kv-geofences.yaml, kv-mission-state.yaml,
// kv-entity-correlation.yaml), projectors running in-process, and collectors on the two
// output subjects. The projector never creates streams or buckets.

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

type testServer struct {
	url string
	cmd *exec.Cmd
}

func storeRoot(t *testing.T) string {
	t.Helper()
	if root := os.Getenv("NATS_TEST_STORE_ROOT"); root != "" {
		_ = os.MkdirAll(root, 0o755)
		d, err := os.MkdirTemp(root, "js-")
		if err != nil {
			t.Fatal(err)
		}
		t.Cleanup(func() { _ = os.RemoveAll(d) })
		return d
	}
	return t.TempDir()
}

func freePort(t *testing.T) int {
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer l.Close()
	return l.Addr().(*net.TCPAddr).Port
}

func natsServerBin(t *testing.T) string {
	if bin := os.Getenv("NATS_SERVER_BIN"); bin != "" {
		return bin
	}
	bin, err := exec.LookPath("nats-server")
	if err != nil {
		t.Skip("nats-server not found; set NATS_SERVER_BIN")
	}
	return bin
}

func startServer(t *testing.T) *testServer {
	t.Helper()
	bin := natsServerBin(t)
	port := freePort(t)
	s := &testServer{url: fmt.Sprintf("nats://127.0.0.1:%d", port)}
	s.cmd = exec.Command(bin, "-js", "-a", "127.0.0.1", "-p", fmt.Sprint(port), "-sd", storeRoot(t))
	s.cmd.Stdout, s.cmd.Stderr = io.Discard, io.Discard
	if err := s.cmd.Start(); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = s.cmd.Process.Kill(); _, _ = s.cmd.Process.Wait() })
	for i := 0; i < 100; i++ {
		if c, err := net.Dial("tcp", fmt.Sprintf("127.0.0.1:%d", port)); err == nil {
			c.Close()
			return s
		}
		time.Sleep(50 * time.Millisecond)
	}
	t.Fatal("nats-server did not start")
	return nil
}

func connect(t *testing.T, url string) (*nats.Conn, jetstream.JetStream) {
	t.Helper()
	nc, err := nats.Connect(url)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(nc.Close)
	js, _ := jetstream.New(nc)
	return nc, js
}

type bus struct {
	url      string
	nc       *nats.Conn
	js       jetstream.JetStream
	geof     jetstream.KeyValue
	missions jetstream.KeyValue
}

func setupBus(t *testing.T) *bus {
	t.Helper()
	s := startServer(t)
	nc, js := connect(t, s.url)
	ctx := context.Background()
	if _, err := js.CreateStream(ctx, jetstream.StreamConfig{ // ent.yaml
		Name: "ENT", Subjects: []string{"ent.*.*.state"}, Storage: jetstream.FileStorage,
		Retention: jetstream.LimitsPolicy, Discard: jetstream.DiscardOld, MaxMsgsPerSubject: 1,
		MaxAge: 7 * 24 * time.Hour, MaxBytes: 1 << 30, Duplicates: 2 * time.Minute,
		AllowDirect: true, AllowMsgTTL: true, SubjectDeleteMarkerTTL: 15 * time.Minute,
	}); err != nil {
		t.Fatal(err)
	}
	if _, err := js.CreateStream(ctx, jetstream.StreamConfig{ // cot.yaml
		Name: "COT", Subjects: []string{"tak.cot.>"}, Storage: jetstream.FileStorage,
		MaxAge: 720 * time.Hour, Duplicates: 2 * time.Minute,
	}); err != nil {
		t.Fatal(err)
	}
	b := &bus{url: s.url, nc: nc, js: js}
	var err error
	if b.geof, err = js.CreateKeyValue(ctx, jetstream.KeyValueConfig{Bucket: "geofences", History: 10, Storage: jetstream.FileStorage}); err != nil {
		t.Fatal(err)
	}
	if b.missions, err = js.CreateKeyValue(ctx, jetstream.KeyValueConfig{Bucket: "mission-state", History: 5, Storage: jetstream.FileStorage}); err != nil {
		t.Fatal(err)
	}
	if _, err = js.CreateKeyValue(ctx, jetstream.KeyValueConfig{Bucket: "entity-correlation", History: 5}); err != nil {
		t.Fatal(err)
	}
	return b
}

// putEnt writes one entity record to ENT as the entity service does.
func (b *bus) putEnt(t *testing.T, m map[string]any) {
	t.Helper()
	subj := fmt.Sprintf("ent.%s.%s.state", m["entity_type"], m["entity_id"])
	if _, err := b.js.Publish(context.Background(), subj, mustJSON(m)); err != nil {
		t.Fatal(err)
	}
}

func (b *bus) entSubjects(t *testing.T) uint64 {
	st, err := b.js.Stream(context.Background(), "ENT")
	if err != nil {
		t.Fatal(err)
	}
	info, err := st.Info(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	return info.State.NumSubjects
}

// logBuf captures the projector's JSON log lines.
type logBuf struct {
	mu sync.Mutex
	b  bytes.Buffer
}

func (l *logBuf) Write(p []byte) (int, error) { l.mu.Lock(); defer l.mu.Unlock(); return l.b.Write(p) }
func (l *logBuf) String() string              { l.mu.Lock(); defer l.mu.Unlock(); return l.b.String() }

type runningProjector struct {
	svc  *Service
	logs *logBuf
	stop func()
}

// startProjector runs a projector in-process with test timers and waits until it is ready.
func startProjector(t *testing.T, b *bus, cfg Config) *runningProjector {
	t.Helper()
	nc, err := nats.Connect(b.url)
	if err != nil {
		t.Fatal(err)
	}
	lb := &logBuf{}
	log := slog.New(slog.NewJSONHandler(io.MultiWriter(lb), &slog.HandlerOptions{Level: slog.LevelInfo}))
	svc, err := New(cfg, nc, log)
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan struct{})
	go func() { defer close(done); _ = svc.Run(ctx) }()
	stopped := false
	stop := func() {
		if stopped {
			return
		}
		stopped = true
		cancel()
		<-done
		nc.Close()
	}
	t.Cleanup(stop)
	deadline := time.Now().Add(10 * time.Second)
	for !svc.Ready() {
		if time.Now().After(deadline) {
			t.Fatal("projector not ready")
		}
		time.Sleep(20 * time.Millisecond)
	}
	return &runningProjector{svc: svc, logs: lb, stop: stop}
}

// cotMsg is one message seen on an output subject.
type cotMsg struct {
	at      time.Time // receive time
	sent    time.Time // the CoT `time` attribute (the projector's send time)
	subject string
	uid     string
	typ     string
	removal bool
	xml     string
	hdr     nats.Header
}

type collector struct {
	mu   sync.Mutex
	msgs []cotMsg
}

var reUID = regexp.MustCompile(`\buid="([^"]+)"`)
var reType = regexp.MustCompile(`<event[^>]*\btype="([^"]+)"`)
var reTime = regexp.MustCompile(`<event[^>]*\btime="([^"]+)"`)

// collect records the output subjects on its own connection, so a busy publisher in the test
// does not delay delivery.
func collect(t *testing.T, b0 *nats.Conn, subjects ...string) *collector {
	t.Helper()
	nc, err := nats.Connect(b0.ConnectedUrl())
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(nc.Close)
	c := &collector{}
	for _, s := range subjects {
		sub, err := nc.Subscribe(s, func(m *nats.Msg) {
			x := string(m.Data)
			cm := cotMsg{at: time.Now(), subject: m.Subject, xml: x, hdr: m.Header,
				removal: strings.Contains(x, "<remarks>removed</remarks>")}
			if u := reUID.FindStringSubmatch(x); u != nil {
				cm.uid = u[1]
			}
			if ty := reType.FindStringSubmatch(x); ty != nil {
				cm.typ = ty[1]
			}
			if tm := reTime.FindStringSubmatch(x); tm != nil {
				cm.sent, _ = time.Parse(cotTimeFmt, tm[1])
			}
			c.mu.Lock()
			c.msgs = append(c.msgs, cm)
			c.mu.Unlock()
		})
		if err != nil {
			t.Fatal(err)
		}
		t.Cleanup(func() { _ = sub.Unsubscribe() })
	}
	_ = nc.Flush()
	return c
}

func (c *collector) all() []cotMsg {
	c.mu.Lock()
	defer c.mu.Unlock()
	return append([]cotMsg(nil), c.msgs...)
}

func (c *collector) reset() {
	c.mu.Lock()
	c.msgs = nil
	c.mu.Unlock()
}

func (c *collector) on(subject string) []cotMsg {
	var out []cotMsg
	for _, m := range c.all() {
		if m.subject == subject {
			out = append(out, m)
		}
	}
	return out
}

// waitFor polls cond until it is true or the timeout passes.
func waitFor(t *testing.T, d time.Duration, what string, cond func() bool) time.Duration {
	t.Helper()
	start := time.Now()
	for !cond() {
		if time.Since(start) > d {
			t.Fatalf("timeout after %s waiting for %s", d, what)
		}
		time.Sleep(20 * time.Millisecond)
	}
	return time.Since(start)
}

// testCfg is a target's defaults with test-friendly timers.
func testCfg(target Target) Config {
	c := Defaults(target)
	c.CapEvery = 200 * time.Millisecond
	return c
}

// evidence appends a line to $A6_EVIDENCE_LOG when set (the onemind-xr evidence dir).
func evidence(t *testing.T, format string, args ...any) {
	t.Logf(format, args...)
	p := os.Getenv("A6_EVIDENCE_LOG")
	if p == "" {
		return
	}
	f, err := os.OpenFile(p, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return
	}
	defer f.Close()
	fmt.Fprintf(f, "%s %s: %s\n", time.Now().UTC().Format(time.RFC3339), t.Name(), fmt.Sprintf(format, args...))
}

// entityServiceBin builds services/entity (the A4 writer of ENT, on this branch's base) once,
// or uses $ENTITY_SERVICE_BIN.
var (
	entBinOnce sync.Once
	entBin     string
	entBinErr  error
)

func entityServiceBin(t *testing.T) string {
	t.Helper()
	if p := os.Getenv("ENTITY_SERVICE_BIN"); p != "" {
		return p
	}
	entBinOnce.Do(func() {
		dir, err := os.MkdirTemp("", "a6-entity-bin-")
		if err != nil {
			entBinErr = err
			return
		}
		entBin = filepath.Join(dir, "entity-service")
		src, _ := filepath.Abs("../../../entity")
		cmd := exec.Command("go", "build", "-o", entBin, "./cmd/entity-service")
		cmd.Dir = src
		out, err := cmd.CombinedOutput()
		if err != nil {
			entBinErr = fmt.Errorf("build entity-service: %v\n%s", err, out)
		}
	})
	if entBinErr != nil {
		t.Skipf("entity service not buildable here: %v", entBinErr)
	}
	return entBin
}

// startEntityService runs the A4 entity service against the bus (inputs: COT only).
func startEntityService(t *testing.T, b *bus) *logBuf {
	t.Helper()
	bin := entityServiceBin(t)
	lb := &logBuf{}
	cmd := exec.Command(bin)
	cmd.Env = append(os.Environ(), "NATS_URL="+b.url, "INPUTS=COT=tak.cot.>", "SERVE_RPC=false",
		fmt.Sprintf("HTTP_ADDR=127.0.0.1:%d", freePort(t)), "POD_NAME=a6-test")
	cmd.Stdout, cmd.Stderr = lb, lb
	if err := cmd.Start(); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = cmd.Process.Kill(); _, _ = cmd.Process.Wait() })
	return lb
}

func jsonGet(m map[string]any, path ...string) any {
	var cur any = m
	for _, p := range path {
		mm, ok := cur.(map[string]any)
		if !ok {
			return nil
		}
		cur = mm[p]
	}
	return cur
}

func readEnt(t *testing.T, b *bus, subject string) map[string]any {
	t.Helper()
	st, err := b.js.Stream(context.Background(), "ENT")
	if err != nil {
		t.Fatal(err)
	}
	msg, err := st.GetLastMsgForSubject(context.Background(), subject)
	if err != nil {
		return nil
	}
	var m map[string]any
	_ = json.Unmarshal(msg.Data, &m)
	return m
}
