package bus

// Test harness: a disposable nats-server (binary from $NATS_SERVER_BIN or PATH; the tests skip
// when it is absent) with the streams and bucket created here with the settings of the CRs
// (platform/nats/streams/ent.yaml, cmd.yaml, mission.yaml, evt.yaml, kv-geofences.yaml). The
// service itself never creates streams or buckets. A test Ed25519 key stands in for the policy
// service and signs commands exactly as services/policy publishCommandID does (PR #554).

import (
	"context"
	"crypto/ed25519"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net"
	"os"
	"os/exec"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"

	"github.com/onemindos/omos/services/geofence/internal/engine"
	"github.com/onemindos/omos/services/geofence/internal/sig"
)

type server struct{ url string }

func startServer(t *testing.T) *server {
	t.Helper()
	bin := os.Getenv("NATS_SERVER_BIN")
	if bin == "" {
		var err error
		if bin, err = exec.LookPath("nats-server"); err != nil {
			t.Skip("nats-server not found; set NATS_SERVER_BIN")
		}
	}
	dir := t.TempDir()
	if root := os.Getenv("NATS_TEST_STORE_ROOT"); root != "" {
		_ = os.MkdirAll(root, 0o755)
		d, err := os.MkdirTemp(root, "js-")
		if err != nil {
			t.Fatal(err)
		}
		dir = d
		t.Cleanup(func() { _ = os.RemoveAll(d) })
	}
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	port := l.Addr().(*net.TCPAddr).Port
	_ = l.Close()
	cmd := exec.Command(bin, "-js", "-a", "127.0.0.1", "-p", fmt.Sprint(port), "-sd", dir)
	cmd.Stdout, cmd.Stderr = io.Discard, io.Discard
	if err := cmd.Start(); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = cmd.Process.Kill(); _, _ = cmd.Process.Wait() })
	url := fmt.Sprintf("nats://127.0.0.1:%d", port)
	for i := 0; i < 100; i++ {
		if nc, err := nats.Connect(url); err == nil {
			nc.Close()
			return &server{url: url}
		}
		time.Sleep(50 * time.Millisecond)
	}
	t.Fatal("nats-server did not start")
	return nil
}

func connect(t *testing.T, s *server) (*nats.Conn, jetstream.JetStream) {
	t.Helper()
	nc, err := nats.Connect(s.url)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(nc.Close)
	js, _ := jetstream.New(nc)
	return nc, js
}

func setupStreams(t *testing.T, js jetstream.JetStream) {
	t.Helper()
	ctx := context.Background()
	for _, c := range []jetstream.StreamConfig{
		{Name: "ENT", Subjects: []string{"ent.*.*.state"}, MaxMsgsPerSubject: 1, AllowDirect: true, AllowMsgTTL: true,
			SubjectDeleteMarkerTTL: 15 * time.Minute, Duplicates: 2 * time.Minute, Storage: jetstream.FileStorage},
		{Name: "CMD", Subjects: []string{"cmd.>"}, Duplicates: 2 * time.Minute, Storage: jetstream.FileStorage},
		{Name: "MISSION", Subjects: []string{"evt.mission.>", "evt.map.>"}, Duplicates: 2 * time.Minute, Storage: jetstream.FileStorage},
		{Name: "EVT", Subjects: []string{"evt.so.>", "evt.drone.>", "evt.robot.>"}, Duplicates: 2 * time.Minute, Storage: jetstream.FileStorage},
	} {
		if _, err := js.CreateStream(ctx, c); err != nil {
			t.Fatalf("stream %s: %v", c.Name, err)
		}
	}
	if _, err := js.CreateKeyValue(ctx, jetstream.KeyValueConfig{Bucket: "geofences", History: 10, Storage: jetstream.FileStorage}); err != nil {
		t.Fatal(err)
	}
}

type policyKey struct {
	kid  string
	priv ed25519.PrivateKey
	pub  ed25519.PublicKey
}

func newPolicyKey(kid string) policyKey {
	pub, priv, _ := ed25519.GenerateKey(rand.Reader)
	return policyKey{kid: kid, priv: priv, pub: pub}
}

type running struct {
	svc    *Service
	cancel context.CancelFunc
	done   chan error
}

func runService(t *testing.T, s *server, keys sig.Keys, mod func(*Config, *engine.Config)) *running {
	t.Helper()
	nc, err := nats.Connect(s.url)
	if err != nil {
		t.Fatal(err)
	}
	ecfg := engine.Config{FeedSilence: 90 * time.Second}
	cfg := Config{Keys: keys, RequireSignature: true, RequireApproval: true, TickEvery: 200 * time.Millisecond, ServeRPC: true,
		Log: slog.New(slog.NewTextHandler(io.Discard, nil))}
	if os.Getenv("GEOFENCE_TEST_LOG") != "" {
		cfg.Log = slog.New(slog.NewTextHandler(os.Stderr, nil))
	}
	if mod != nil {
		mod(&cfg, &ecfg)
	}
	svc, err := New(nc, engine.New(ecfg), cfg)
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	r := &running{svc: svc, cancel: cancel, done: make(chan error, 1)}
	go func() { r.done <- svc.Run(ctx) }()
	deadline := time.Now().Add(10 * time.Second)
	for !svc.Ready() {
		select {
		case err := <-r.done:
			t.Fatalf("service stopped: %v", err)
		default:
		}
		if time.Now().After(deadline) {
			t.Fatal("service not ready")
		}
		time.Sleep(20 * time.Millisecond)
	}
	var once sync.Once
	stop := func() {
		once.Do(func() {
			cancel()
			select {
			case <-r.done:
			case <-time.After(5 * time.Second):
				t.Error("service did not stop")
			}
			nc.Close()
		})
	}
	t.Cleanup(stop)
	r.cancel = stop
	return r
}

func uuid4() string {
	var b [16]byte
	_, _ = rand.Read(b[:])
	b[6] = (b[6] & 0x0f) | 0x40
	b[8] = (b[8] & 0x3f) | 0x80
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:16])
}

type cmdOpt struct {
	key       *policyKey // nil = unsigned
	expiresIn time.Duration
	approval  string // approved | not_required | ""
	id        string
}

// policyCommand builds the envelope the policy service publishes after approval
// (publishCommandID in services/policy/internal/policy/service.go, PR #554).
func policyCommand(t *testing.T, subject string, data any, o cmdOpt) (string, []byte) {
	t.Helper()
	if o.id == "" {
		o.id = uuid4()
	}
	if o.expiresIn == 0 {
		o.expiresIn = 30 * time.Second
	}
	now := time.Now()
	env := map[string]any{"v": 1, "id": o.id, "ts": now.UTC().Format("2006-01-02T15:04:05.000Z"), "cmd": subject,
		"source": "omos-policy", "subject": subject, "data": data, "action_id": "act-" + o.id[:8],
		"expires_at": now.Add(o.expiresIn).UTC().Format("2006-01-02T15:04:05.000Z"), "simulated": true,
		"policy": map[string]any{"decision_id": "dec-" + o.id[:8], "allowed": true}}
	if o.approval != "" {
		env["approval"] = map[string]any{"required": o.approval == "approved", "status": o.approval, "approver": "user:zeus",
			"level": "A2", "params_hash": "sha256:00", "approvals": []map[string]any{{"by": "user:zeus"}}}
	}
	if o.key != nil {
		// round-trip through JSON first so numbers sign as they will be read
		b, _ := json.Marshal(env)
		env = map[string]any{}
		_ = json.Unmarshal(b, &env)
		if err := sig.Sign(o.key.priv, o.key.kid, env); err != nil {
			t.Fatal(err)
		}
	}
	b, err := json.Marshal(env)
	if err != nil {
		t.Fatal(err)
	}
	return o.id, b
}

func publishCmd(t *testing.T, js jetstream.JetStream, subject string, b []byte, msgID string) {
	t.Helper()
	if _, err := js.Publish(context.Background(), subject, b, jetstream.WithMsgID(msgID)); err != nil {
		t.Fatal(err)
	}
}

type fact struct {
	Subject string
	Env     struct {
		ID        string          `json:"id"`
		Type      string          `json:"type"`
		Source    string          `json:"source"`
		Causation string          `json:"causation"`
		Data      json.RawMessage `json:"data"`
	}
}

// facts reads every message on a stream that matches the filter.
func facts(t *testing.T, js jetstream.JetStream, stream, filter string) []fact {
	t.Helper()
	ctx := context.Background()
	st, err := js.Stream(ctx, stream)
	if err != nil {
		t.Fatal(err)
	}
	info, err := st.Info(ctx, jetstream.WithSubjectFilter(filter))
	if err != nil {
		t.Fatal(err)
	}
	var n uint64
	for _, c := range info.State.Subjects {
		n += c
	}
	if n == 0 {
		return nil
	}
	oc, err := js.OrderedConsumer(ctx, stream, jetstream.OrderedConsumerConfig{FilterSubjects: []string{filter}})
	if err != nil {
		t.Fatal(err)
	}
	var out []fact
	for uint64(len(out)) < n {
		batch, err := oc.Fetch(int(n)-len(out), jetstream.FetchMaxWait(2*time.Second))
		if err != nil {
			t.Fatal(err)
		}
		got := 0
		for m := range batch.Messages() {
			got++
			f := fact{Subject: m.Subject()}
			if err := json.Unmarshal(m.Data(), &f.Env); err != nil {
				t.Fatalf("not an envelope on %s: %s", m.Subject(), m.Data())
			}
			out = append(out, f)
		}
		if got == 0 {
			t.Fatalf("read %d of %d messages on %s", len(out), n, filter)
		}
	}
	return out
}

func waitFor(t *testing.T, what string, timeout time.Duration, cond func() bool) {
	t.Helper()
	deadline := time.Now().Add(timeout)
	for !cond() {
		if time.Now().After(deadline) {
			t.Fatalf("timed out waiting for %s", what)
		}
		time.Sleep(50 * time.Millisecond)
	}
}

func countFacts(t *testing.T, js jetstream.JetStream, stream, filter string) int {
	return len(facts(t, js, stream, filter))
}

func lastFact(t *testing.T, js jetstream.JetStream, stream, filter string) (fact, bool) {
	fs := facts(t, js, stream, filter)
	if len(fs) == 0 {
		return fact{}, false
	}
	return fs[len(fs)-1], true
}

// entRecord is an `ent.drone.sim-x1.state` record shaped like the entity service output
// (services/entity, PR #552; entity.md §4.1), reduced to the fields the watcher reads plus the
// usual envelope fields.
func entRecord(id string, lat, lon float64, seq uint64, lifecycle string) []byte {
	b, _ := json.Marshal(map[string]any{
		"uid": id, "type": "fleet.drone", "lat": lat, "lon": lon, "hae": 152.4, "callsign": strings.ToUpper(id),
		"stale": 60, "t": time.Now().UnixMilli(), "entity_id": id, "entity_type": "drone", "confidence": 0.9,
		"v": 1, "seq": seq, "live": lifecycle != "expired", "lifecycle": lifecycle, "template": "asset",
		"traits": []string{"locatable", "taskable"}, "indicators": map[string]any{"simulated": true},
	})
	return b
}

func publishEnt(t *testing.T, js jetstream.JetStream, typ, id string, data []byte, seq uint64) {
	t.Helper()
	if _, err := js.Publish(context.Background(), "ent."+typ+"."+id+".state", data,
		jetstream.WithMsgID(fmt.Sprintf("%s:%d", id, seq))); err != nil {
		t.Fatal(err)
	}
}

func keepIn(lat, lon, h float64) map[string]any {
	ring := [][]float64{{lon - h, lat - h}, {lon + h, lat - h}, {lon + h, lat + h}, {lon - h, lat + h}, {lon - h, lat - h}}
	return map[string]any{"type": "Polygon", "coordinates": [][][]float64{ring}}
}
