package presence

import (
	"encoding/json"
	"testing"
	"time"
)

func TestUIDFor(t *testing.T) {
	if got := UIDFor("hermes", "zeus", "legacy-edge"); got != "agent-hermes-zeus-legacy-edge" {
		t.Fatal(got)
	}
	if got := UIDFor("Claude Code", "A.B", "x*y"); got != "agent-claude-code-a-b-x-y" {
		t.Fatal(got)
	}
	long := UIDFor("a", "b", string(make([]byte, 100)))
	if len(long) > 64 {
		t.Fatal(long)
	}
}

func TestHeartbeatToEnvelope(t *testing.T) {
	now := time.Date(2026, 10, 29, 12, 0, 0, 0, time.UTC)
	tr := New(Config{Site: Anchor{30.2672, -97.7431}, Anchors: map[string]Anchor{"hermes/zeus/legacy": {30.27, -97.74}}}, func() time.Time { return now })
	hb := `{"agent":"hermes","owner":"zeus","session":"legacy","instance_id":"i1","ts":"2026-10-29T12:00:00Z","interval_s":30,"name":"Legacy",
	  "onemind":{"mode":"busy","activity":"answering a TAK chat","current_task":{"task_id":"tsk-01","state":"executing","summary":"brief"}}}`
	a, err := tr.OnHeartbeat("agents.hb.hermes.zeus.legacy", []byte(hb))
	if err != nil {
		t.Fatal(err)
	}
	env := tr.EnvelopeFor(*a, now)
	if env.UID != "agent-hermes-zeus-legacy" || env.Type != "fleet.agent" || env.Callsign != "LEGACY" || env.Lat != 30.27 || env.Stale != 90 {
		t.Fatalf("%+v", env)
	}
	ct, _ := env.Meta["current_task"].(map[string]any)
	if env.Meta["mode"] != "busy" || env.Meta["position_kind"] != "anchor" || ct["task_id"] != "tsk-01" || env.Meta["prompt_subject"] != "agents.prompt.hermes.zeus.legacy" {
		t.Fatalf("meta %+v", env.Meta)
	}
	b, _ := json.Marshal(env)
	var back map[string]any
	_ = json.Unmarshal(b, &back)
	for _, k := range []string{"uid", "type", "lat", "lon", "t"} { // doctrine §9.2 required fields
		if _, ok := back[k]; !ok {
			t.Errorf("missing %s", k)
		}
	}
	// Due: first call publishes, an unchanged state does not, a change does.
	if len(tr.Due(now, false)) != 1 || len(tr.Due(now.Add(time.Second), false)) != 0 {
		t.Fatal("due on change only")
	}
	// Offline after 3 intervals, forgotten after ForgetAfter.
	later := now.Add(91 * time.Second)
	d := tr.Due(later, false)
	if len(d) != 1 || d[0].Meta["mode"] != "offline" || d[0].Meta["connection"] != "offline" {
		t.Fatalf("offline %+v", d)
	}
	if len(tr.Due(now.Add(11*time.Minute), false)) != 0 || len(tr.List()) != 0 {
		t.Fatal("not forgotten")
	}
}

func TestHeartbeatMismatchRefused(t *testing.T) {
	tr := New(Config{}, nil)
	if _, err := tr.OnHeartbeat("agents.hb.hermes.zeus.legacy", []byte(`{"agent":"evil","owner":"zeus","instance_id":"x","interval_s":30}`)); err == nil {
		t.Fatal("body/subject mismatch accepted")
	}
	if _, err := tr.OnHeartbeat("agents.prompt.hermes.zeus.legacy", []byte(`{}`)); err == nil {
		t.Fatal("non-heartbeat subject accepted")
	}
}
