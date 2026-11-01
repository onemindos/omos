package chat

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/nats-io/nats.go"

	"github.com/onemindos/omos/services/agentops/internal/agentproto"
	"github.com/onemindos/omos/services/agentops/internal/bus"
	"github.com/onemindos/omos/services/agentops/internal/cot"
	"github.com/onemindos/omos/services/agentops/internal/natstest"
	"github.com/onemindos/omos/services/agentops/internal/presence"
)

func geochat(sender, callsign, toUID, toCall, text string, at time.Time) []byte {
	ts := at.UTC().Format("2006-01-02T15:04:05.000Z")
	return []byte(fmt.Sprintf(`<event version="2.0" uid="GeoChat.%s.%s.%d" type="b-t-f" how="h-g-i-g-o" time="%s" start="%s" stale="%s">`+
		`<point lat="30.1" lon="-97.7" hae="0" ce="9999999" le="9999999"/><detail>`+
		`<__chat parent="RootContactGroup" groupOwner="false" messageId="m%d" chatroom="%s" id="%s" senderCallsign="%s">`+
		`<chatgrp uid0="%s" uid1="%s" id="%s"/></__chat><link uid="%s" type="a-f-G-U-C" relation="p-p"/>`+
		`<remarks source="BAO.F.ATAK.%s" to="%s" time="%s">%s</remarks></detail></event>`,
		sender, toUID, at.UnixNano(), ts, ts, ts, at.UnixNano(), toCall, toUID, callsign, sender, toUID, toUID, sender, sender, toUID, ts, text))
}

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

// chatsTo returns the GeoChat replies addressed to uid.
func (r *recorder) chatsTo(uid string) []*cot.Chat {
	var out []*cot.Chat
	for _, m := range r.all() {
		if c, err := cot.ParseChat(m.Data); err == nil && c.AddressedTo(uid, "") {
			out = append(out, c)
		}
	}
	return out
}

// TestAgentIsAContactAndATAKChatReachesIt is the I1 test on this machine: Hermes (a stub
// agent speaking the Synadia Agent Protocol) beats with its current task; presence publishes
// it as a fleet.agent geo-envelope; the gateway announces it to TAK as a contact; a GeoChat
// from an ATAK user (TAK Server path and mesh path) reaches it and the answer goes back.
func TestAgentIsAContactAndATAKChatReachesIt(t *testing.T) {
	srv := natstest.Start(t)
	nc := srv.Connect(t)
	js := natstest.Streams(t, nc, true)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Stub Hermes: answers prompts and beats with its state.
	hermesNC := srv.Connect(t)
	var mu sync.Mutex
	var prompts []map[string]any
	var promptText []string
	sub, err := agentproto.Serve(hermesNC, "agents.prompt.hermes.zeus.legacy", func(_ context.Context, p string, env map[string]any) (string, error) {
		mu.Lock()
		prompts, promptText = append(prompts, env), append(promptText, p)
		mu.Unlock()
		return "SIM-X1 is inside zone-home, battery 82%. I can propose an RTL; a human must confirm it in OneMind.", nil
	})
	if err != nil {
		t.Fatal(err)
	}
	beater := agentproto.NewBeater(hermesNC, "hermes", "zeus", "legacy", "Legacy", time.Second)
	beater.Set(agentproto.OneMind{Mode: "busy", Activity: "watching the homestead", CurrentTask: &agentproto.Task{TaskID: "tsk-01jbx7k2m9q4r8s6t0v2w4y6z8", State: "executing", Summary: "AOI brief"}})
	go func() { _ = beater.Run(ctx) }()

	pub, _ := bus.NewPublisher(nc, bus.ModeAuto)
	tr := presence.New(presence.Config{Site: presence.Anchor{Lat: 30.2672, Lon: -97.7431}, PublishEvery: time.Second}, nil)
	go func() { _ = tr.Run(ctx, nc, pub, "AGENTS") }()

	natstest.Eventually(t, 5*time.Second, "agent geo-envelope in TEL", func() bool {
		return natstest.Count(t, js, "TEL", presence.Subject("agent-hermes-zeus-legacy")) > 0
	})
	msgs := natstest.Messages(t, js, "TEL", presence.Subject("agent-hermes-zeus-legacy"))
	var env presence.Envelope
	_ = json.Unmarshal(msgs[len(msgs)-1], &env)
	ct, _ := env.Meta["current_task"].(map[string]any)
	if env.Type != "fleet.agent" || env.Callsign != "LEGACY" || env.Meta["mode"] != "busy" || ct["task_id"] != "tsk-01jbx7k2m9q4r8s6t0v2w4y6z8" {
		t.Fatalf("envelope %+v", env)
	}
	natstest.Evidence(t, "TEL %s: type=%s callsign=%s lat=%v lon=%v mode=%v activity=%v current_task=%v position_kind=%v",
		presence.Subject(env.UID), env.Type, env.Callsign, env.Lat, env.Lon, env.Meta["mode"], env.Meta["activity"], ct, env.Meta["position_kind"])

	inject, mesh := &recorder{}, &recorder{}
	obs := srv.Connect(t)
	_, _ = obs.Subscribe(OutServer, inject.add)
	_, _ = obs.Subscribe(OutMesh, mesh.add)
	_ = obs.Flush()
	gw := New(Config{ReplyServer: true, ReplyMesh: true, ContactServer: true, ContactEvery: time.Second, Agents: []string{"hermes/zeus/legacy"},
		PerSenderPerMin: 3, Version: "test"}, nc, tr, slog.New(slog.NewTextHandler(io.Discard, nil)))
	go func() { _ = gw.Run(ctx) }()

	agentUID := ContactUID("agent-hermes-zeus-legacy")
	natstest.Eventually(t, 5*time.Second, "contact PLI on tak.cot.inject", func() bool {
		for _, m := range inject.all() {
			if strings.Contains(string(m.Data), `uid="`+agentUID+`"`) && strings.Contains(string(m.Data), `<contact callsign="LEGACY" endpoint=`) {
				return true
			}
		}
		return false
	})
	for _, m := range inject.all() {
		if strings.Contains(string(m.Data), `uid="`+agentUID+`"`) {
			natstest.Evidence(t, "contact PLI on %s: %s", OutServer, m.Data)
			break
		}
	}

	// TAK Server path: an ATAK user messages LEGACY directly.
	_ = nc.Publish(InServer, geochat("ANDROID-1", "ZEUS", agentUID, "LEGACY", "what is sim-x1 doing?", time.Now()))
	natstest.Eventually(t, 5*time.Second, "reply on tak.cot.inject", func() bool { return len(inject.chatsTo("ANDROID-1")) == 1 })
	r := inject.chatsTo("ANDROID-1")[0]
	if r.SenderUID != agentUID || !strings.Contains(r.Text, "SIM-X1 is inside zone-home") || len(r.DestCallsigns) != 1 || r.DestCallsigns[0] != "ZEUS" {
		t.Fatalf("server reply %+v", r)
	}
	mu.Lock()
	p0, pt0 := prompts[0], promptText[0]
	mu.Unlock()
	om, _ := p0["onemind"].(map[string]any)
	if om["trust"] != "claimed" || om["sender_uid"] != "ANDROID-1" || !strings.Contains(pt0, "identity claimed, not verified") || !strings.Contains(pt0, "what is sim-x1 doing?") {
		t.Fatalf("prompt %v %q", p0, pt0)
	}
	natstest.Evidence(t, "ATAK->Hermes via tak.cot.out: prompt envelope onemind=%v", om)
	natstest.Evidence(t, "Hermes->ATAK reply on tak.cot.inject: from=%s to=ANDROID-1 marti dest=%v text=%q", r.SenderUID, r.DestCallsigns, r.Text)

	// Mesh path: reply goes to tak.cot.mesh.tx, directed with Tak-Dest-Uid.
	_ = nc.Publish(InMesh, geochat("ANDROID-2", "ALPHA", agentUID, "LEGACY", "status?", time.Now()))
	natstest.Eventually(t, 5*time.Second, "reply on tak.cot.mesh.tx", func() bool { return len(mesh.chatsTo("ANDROID-2")) == 1 })
	var dest string
	for _, m := range mesh.all() {
		if c, err := cot.ParseChat(m.Data); err == nil && c.AddressedTo("ANDROID-2", "") {
			dest = m.Header.Get("Tak-Dest-Uid")
		}
	}
	if dest != "ANDROID-2" {
		t.Fatalf("mesh reply Tak-Dest-Uid=%q", dest)
	}
	natstest.Evidence(t, "mesh path: reply on %s with Tak-Dest-Uid=%s", OutMesh, dest)

	// Our own reply heard back from TAK is an echo, never a new prompt.
	mu.Lock()
	before := len(prompts)
	mu.Unlock()
	_ = nc.Publish(InServer, inject.all()[len(inject.all())-1].Data)
	// A message for nobody, and a stale one, are not forwarded.
	_ = nc.Publish(InServer, geochat("ANDROID-1", "ZEUS", "ANDROID-9", "BRAVO", "hi bravo", time.Now()))
	_ = nc.Publish(InServer, geochat("ANDROID-1", "ZEUS", agentUID, "LEGACY", "old message", time.Now().Add(-10*time.Minute)))
	time.Sleep(500 * time.Millisecond)
	mu.Lock()
	after := len(prompts)
	mu.Unlock()
	if after != before || gw.C.Echo.Load() < 1 || gw.C.TooOld.Load() != 1 || gw.C.NotForAgent.Load() < 1 {
		t.Fatalf("echo/stale/foreign forwarded: prompts %d->%d echo=%d old=%d foreign=%d", before, after, gw.C.Echo.Load(), gw.C.TooOld.Load(), gw.C.NotForAgent.Load())
	}
	natstest.Evidence(t, "echo, stale and foreign chats not forwarded (echo=%d too_old=%d not_for_agent=%d)", gw.C.Echo.Load(), gw.C.TooOld.Load(), gw.C.NotForAgent.Load())

	// Rate limit: 3 per minute per sender; ANDROID-1 already used 1.
	for i := 0; i < 4; i++ {
		_ = nc.Publish(InServer, geochat("ANDROID-1", "ZEUS", agentUID, "LEGACY", fmt.Sprintf("ping %d", i), time.Now()))
	}
	natstest.Eventually(t, 5*time.Second, "rate-limit reply", func() bool {
		for _, c := range inject.chatsTo("ANDROID-1") {
			if strings.HasPrefix(c.Text, "Rate limit") {
				return true
			}
		}
		return false
	})
	natstest.Evidence(t, "rate limit: %d messages limited", gw.C.RateLimited.Load())

	// Agent offline: the user is told, nothing is acted on.
	_ = sub.Unsubscribe()
	_ = hermesNC.Flush()
	_ = nc.Publish(InMesh, geochat("ANDROID-3", "CHARLIE", agentUID, "LEGACY", "are you there?", time.Now()))
	natstest.Eventually(t, 5*time.Second, "offline reply", func() bool {
		cs := mesh.chatsTo("ANDROID-3")
		return len(cs) == 1 && strings.Contains(cs[0].Text, "could not answer right now (agent offline)")
	})
	natstest.Evidence(t, "agent offline: user told %q", mesh.chatsTo("ANDROID-3")[0].Text)
}

// TestRateLimitNoticeOncePerSender: a flooding sender gets one rate-limit notice per minute,
// the notice uses the global budget, and no notice is sent once that budget is spent.
func TestRateLimitNoticeOncePerSender(t *testing.T) {
	g := New(Config{PerSenderPerMin: 2, GlobalPerMin: 6}, nil, nil, nil)
	now := time.Date(2026, 10, 29, 12, 0, 0, 0, time.UTC)
	ok, notices := 0, 0
	for i := 0; i < 20; i++ {
		a, n := g.allow("ANDROID-1", now)
		if a {
			ok++
		}
		if n {
			notices++
		}
	}
	if ok != 2 || notices != 1 || len(g.global) != 3 {
		t.Fatalf("ok=%d notices=%d global=%d", ok, notices, len(g.global))
	}
	// Next minute: one more notice is allowed.
	later := now.Add(61 * time.Second)
	_, _ = g.allow("ANDROID-1", later)
	_, _ = g.allow("ANDROID-1", later)
	if _, n := g.allow("ANDROID-1", later); !n {
		t.Fatal("no notice in the next window")
	}
	// Global budget spent: other senders are refused without a notice.
	for i := 0; i < 10; i++ {
		_, _ = g.allow(fmt.Sprintf("U-%d", i), later)
	}
	if a, n := g.allow("U-new", later); a || n {
		t.Fatalf("global full: allowed=%v notice=%v", a, n)
	}
	if len(g.global) > 6 {
		t.Fatalf("global budget exceeded: %d", len(g.global))
	}
	if IsChat([]byte(`<event type="a-f-G-U-C"/>`)) || !IsChat([]byte(`<event type="b-t-f"/>`)) {
		t.Fatal("prefilter")
	}
}
