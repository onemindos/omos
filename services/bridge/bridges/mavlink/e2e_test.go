//go:build e2e

// End-to-end contract test: a real nats-server (per-role users, JetStream) and a real
// MAVLink vehicle (the onemind-xr scripted drone or SITL). Run through test/e2e.sh,
// which starts both and sets the E2E_* variables. Transcript lines go to E2E_TRANSCRIPT.
package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"

	ct "github.com/onemindos/omos-bridge/bridges/mavlink/internal/commandtest"
	"github.com/onemindos/omos-bridge/bridges/mavlink/internal/geoid"
)

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
	if kind != "tel" {
		tr.t.Log(string(line))
	}
}

func mustEnv(t *testing.T, k string) string {
	v := os.Getenv(k)
	if v == "" {
		t.Skipf("%s not set; run test/e2e.sh", k)
	}
	return v
}

func connect(t *testing.T, url, user, pass string, errs chan<- error) *nats.Conn {
	t.Helper()
	nc, err := nats.Connect(url, nats.UserInfo(user, pass), nats.ErrorHandler(func(_ *nats.Conn, _ *nats.Subscription, err error) {
		if errs != nil {
			select {
			case errs <- err:
			default:
			}
		}
	}))
	if err != nil {
		t.Fatalf("connect %s: %v", user, err)
	}
	t.Cleanup(nc.Close)
	return nc
}

func TestE2E(t *testing.T) {
	url := mustEnv(t, "E2E_NATS_URL")
	endpoint := mustEnv(t, "E2E_MAVLINK_ENDPOINT")
	pw := func(role string) string { return mustEnv(t, "E2E_PASS_"+strings.ToUpper(role)) }
	tr := &transcript{t: t}
	if p := os.Getenv("E2E_TRANSCRIPT"); p != "" {
		f, err := os.Create(p)
		if err != nil {
			t.Fatal(err)
		}
		defer f.Close()
		tr.f = f
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()

	// 1. JetStream streams equivalent to omos TEL, CMD and EVT (platform/nats/streams).
	admin := connect(t, url, "admin", pw("admin"), nil)
	js, err := jetstream.New(admin)
	if err != nil {
		t.Fatal(err)
	}
	streams := map[string][]string{"TEL": {"tel.fleet.>"}, "CMD": {"cmd.>"}, "EVT": {"evt.drone.>"}}
	for name, subj := range streams {
		_ = js.DeleteStream(ctx, name)
		if _, err := js.CreateStream(ctx, jetstream.StreamConfig{Name: name, Subjects: subj, Storage: jetstream.FileStorage, Duplicates: 2 * time.Minute}); err != nil {
			t.Fatalf("stream %s: %v", name, err)
		}
	}
	tr.log("setup", "", map[string]any{"streams": streams})

	// 2. A client-role observer records every fleet message it sees.
	client := connect(t, url, "client", pw("client"), nil)
	var mu sync.Mutex
	telSeen := 0
	var outcomes []map[string]any
	if _, err := client.Subscribe("tel.fleet.drone.sim-x1.>", func(m *nats.Msg) {
		var v any
		_ = json.Unmarshal(m.Data, &v)
		mu.Lock()
		telSeen++
		mu.Unlock()
		tr.log("tel", m.Subject, v)
	}); err != nil {
		t.Fatal(err)
	}
	if _, err := client.Subscribe("evt.drone.sim-x1.>", func(m *nats.Msg) {
		var v map[string]any
		_ = json.Unmarshal(m.Data, &v)
		tr.log("evt", m.Subject, v)
		if m.Subject == "evt.drone.sim-x1.command" {
			mu.Lock()
			outcomes = append(outcomes, v["data"].(map[string]any))
			mu.Unlock()
		}
	}); err != nil {
		t.Fatal(err)
	}
	_ = client.Flush()

	// 3. The bridge, with the bridge role and this run's policy key.
	pub, priv := ct.NewKey()
	t.Setenv("ASSET_ID", "sim-x1")
	t.Setenv("MAVLINK_ENDPOINT", endpoint)
	t.Setenv("MAVLINK_OUT_VERSION", os.Getenv("E2E_OUT_VERSION"))
	t.Setenv("POLICY_PUBKEYS", "policy-e2e="+base64.StdEncoding.EncodeToString(pub))
	t.Setenv("GEOFENCE_JSON", `{"keep_in":[{"id":"sim-range","circle":{"lat":0,"lon":0,"radius_m":1000}}],
		"keep_out":[{"id":"no-fly","circle":{"lat":0.005,"lon":0.005,"radius_m":150}}],"max_hae_m":300}`)
	t.Setenv("NATS_URL", url)
	t.Setenv("NATS_USER", "bridge")
	t.Setenv("NATS_PASSWORD", pw("bridge"))
	cfg, err := loadConfig()
	if err != nil {
		t.Fatal(err)
	}
	ready := make(chan struct{})
	bridgeDone := make(chan error, 1)
	bctx, bcancel := context.WithCancel(ctx)
	go func() { bridgeDone <- run(bctx, cfg, func() { close(ready) }) }()
	select {
	case <-ready:
	case err := <-bridgeDone:
		t.Fatalf("bridge exited: %v", err)
	case <-time.After(10 * time.Second):
		t.Fatal("bridge not ready")
	}
	tr.log("setup", "", map[string]any{"bridge": "running", "endpoint": endpoint, "out_version": cfg.OutVersion})

	// 4. Telemetry appears.
	deadline := time.Now().Add(15 * time.Second)
	for {
		mu.Lock()
		n := telSeen
		mu.Unlock()
		if n >= 10 {
			break
		}
		if time.Now().After(deadline) {
			t.Fatalf("telemetry: saw %d messages", n)
		}
		time.Sleep(100 * time.Millisecond)
	}

	// 5. Commands from the policy role, published through JetStream (the PubAck is the proof of record).
	policy := connect(t, url, "policy", pw("policy"), nil)
	pjs, err := jetstream.New(policy)
	if err != nil {
		t.Fatal(err)
	}
	send := func(name string, o ct.Opts, subjOverride string) []map[string]any {
		t.Helper()
		if o.AssetID == "" {
			o.AssetID = "sim-x1"
		}
		if o.Key == nil {
			o.Key = priv
		}
		if o.Kid == "" {
			o.Kid = "policy-e2e"
		}
		subj, body, err := ct.Build(o)
		if err != nil {
			t.Fatal(err)
		}
		if subjOverride != "" {
			subj = subjOverride
		}
		var env map[string]any
		_ = json.Unmarshal(body, &env)
		mu.Lock()
		before := len(outcomes)
		mu.Unlock()
		ack, err := pjs.Publish(ctx, subj, body)
		if err != nil {
			t.Fatalf("%s: publish: %v", name, err)
		}
		tr.log("cmd", subj, map[string]any{"case": name, "stream_seq": ack.Sequence, "envelope": env})
		var got []map[string]any
		wait := time.Now().Add(6 * time.Second)
		for time.Now().Before(wait) {
			mu.Lock()
			got = append([]map[string]any(nil), outcomes[before:]...)
			mu.Unlock()
			if n := len(got); n > 0 {
				last := got[n-1]["outcome"]
				if last == "executed" || last == "failed" || last == "refused" {
					break
				}
			}
			time.Sleep(50 * time.Millisecond)
		}
		return got
	}
	expect := func(name string, got []map[string]any, want ...string) {
		t.Helper()
		var have []string
		for _, o := range got {
			s := fmt.Sprint(o["outcome"])
			if c, ok := o["code"]; ok {
				s += "/" + fmt.Sprint(c)
			}
			have = append(have, s)
		}
		if strings.Join(have, ",") != strings.Join(want, ",") {
			t.Errorf("%s: outcomes %v, want %v", name, have, want)
		}
		tr.log("result", "", map[string]any{"case": name, "outcomes": have, "want": want, "pass": strings.Join(have, ",") == strings.Join(want, ",")})
	}

	n0, _ := geoid.Undulation(0, 0)
	f := false
	_, forged := ct.NewKey()

	// Approved sim commands: ACKed by the vehicle.
	expect("arm (A4, two humans)", send("arm (A4, two humans)", ct.Opts{Verb: "arm"}, ""), "acknowledged", "executed")
	expect("takeoff (A4)", send("takeoff (A4)", ct.Opts{Verb: "takeoff", Data: map[string]any{"hae_m": 60 + n0}}, ""), "acknowledged", "executed")
	expect("goto", send("goto", ct.Opts{Verb: "goto", Data: map[string]any{"lat": 0.001, "lon": 0.001, "hae_m": 80 + n0, "speed_mps": 6.0}}, ""), "acknowledged", "executed")
	// hold (task.abort's verb): DO_REPOSITION to the current position while airborne.
	expect("hold", send("hold", ct.Opts{Verb: "hold"}, ""), "acknowledged", "executed")
	expect("hold (reason abort)", send("hold (reason abort)", ct.Opts{Verb: "hold", Data: map[string]any{"reason": "abort"}}, ""), "acknowledged", "executed")
	expect("rtl", send("rtl", ct.Opts{Verb: "rtl"}, ""), "acknowledged", "executed")
	expect("land", send("land", ct.Opts{Verb: "land"}, ""), "acknowledged", "executed")
	// On the ground the scripted drone answers DO_REPOSITION TEMPORARILY_REJECTED.
	expect("hold while landed", send("hold while landed", ct.Opts{Verb: "hold"}, ""), "acknowledged", "failed/vehicle_rejected")
	expect("disarm", send("disarm", ct.Opts{Verb: "disarm"}, ""), "acknowledged", "executed")

	// Refused: nothing is sent to the vehicle.
	expect("unapproved (status pending)", send("unapproved (status pending)", ct.Opts{Verb: "rtl", Status: "pending"}, ""), "refused/not_authorized")
	expect("arm with one human", send("arm with one human", ct.Opts{Verb: "arm", Approvers: []string{"user:zeus"}}, ""), "refused/not_authorized")
	expect("approved by an agent", send("approved by an agent", ct.Opts{Verb: "rtl", Approvers: []string{"agent:hermes/zeus/legacy"}}, ""), "refused/agent_cannot_confirm")
	expect("simulated=false", send("simulated=false", ct.Opts{Verb: "rtl", Simulated: &f}, ""), "refused/not_simulated")
	expect("forged signature", send("forged signature", ct.Opts{Verb: "rtl", Key: forged}, ""), "refused/bad_signature")
	expect("unsigned", send("unsigned", ct.Opts{Verb: "rtl", NoSig: true}, ""), "refused/bad_signature")
	expect("expired", send("expired", ct.Opts{Verb: "rtl", Now: time.Now().Add(-2 * time.Minute)}, ""), "refused/expired")
	expect("goto into keep-out", send("goto into keep-out", ct.Opts{Verb: "goto", Data: map[string]any{"lat": 0.005, "lon": 0.005, "hae_m": 80 + n0}}, ""), "refused/precondition_failed")
	expect("verb not in catalogue", send("verb not in catalogue", ct.Opts{Verb: "strike"}, ""), "refused/unknown_action")

	// Replay of a recorded, already-executed command.
	replayID := ct.UUID()
	expect("disarm (to be replayed)", send("disarm (to be replayed)", ct.Opts{Verb: "disarm", ID: replayID}, ""), "acknowledged", "executed")
	expect("replay of the same command", send("replay of the same command", ct.Opts{Verb: "disarm", ID: replayID}, ""), "refused/replayed")

	// Non-sim target: the bridge serves only its own sim- id, so nobody answers.
	got := send("non-sim target cmd.fleet.drone.real-1.arm", ct.Opts{Verb: "arm", AssetID: "real-1"}, "")
	expect("non-sim target cmd.fleet.drone.real-1.arm", got)

	// The client role cannot publish cmd.fleet.> at all (server permission violation).
	errs := make(chan error, 4)
	c2 := connect(t, url, "client", pw("client"), errs)
	_, cbody, _ := ct.Build(ct.Opts{AssetID: "sim-x1", Verb: "rtl", Key: priv, Kid: "policy-e2e"})
	_ = c2.Publish("cmd.fleet.drone.sim-x1.rtl", cbody)
	_ = c2.Flush()
	select {
	case err := <-errs:
		ok := strings.Contains(strings.ToLower(err.Error()), "permissions violation")
		tr.log("result", "", map[string]any{"case": "client publishes cmd.fleet.drone.sim-x1.rtl", "server_error": err.Error(), "pass": ok})
		if !ok {
			t.Errorf("client publish: unexpected error %v", err)
		}
	case <-time.After(3 * time.Second):
		t.Error("client publish to cmd.fleet.> was not rejected by the server")
	}

	// 6. Stop the bridge, then compare what the streams stored with what was observed.
	time.Sleep(500 * time.Millisecond)
	bcancel()
	<-bridgeDone
	time.Sleep(300 * time.Millisecond)
	counts := map[string]uint64{}
	for name := range streams {
		s, err := js.Stream(ctx, name)
		if err != nil {
			t.Fatal(err)
		}
		info, err := s.Info(ctx)
		if err != nil {
			t.Fatal(err)
		}
		counts[name] = info.State.Msgs
	}
	mu.Lock()
	tel, outs := telSeen, len(outcomes)
	mu.Unlock()
	tr.log("result", "", map[string]any{"case": "streams stored every fleet message", "stream_msgs": counts, "observed_tel": tel, "observed_outcomes": outs,
		"pass": counts["TEL"] == uint64(tel) && counts["EVT"] >= uint64(outs)})
	if counts["TEL"] != uint64(tel) {
		t.Errorf("TEL stored %d, observer saw %d", counts["TEL"], tel)
	}
	if counts["EVT"] < uint64(outs) {
		t.Errorf("EVT stored %d < %d outcomes", counts["EVT"], outs)
	}
}
