package permtest

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

// expectDelivered subscribes subj as sub, publishes it as pub and requires delivery (the CLIENT
// stream imports of A19-P11: onemind-xr WEB-P2 and UE-P3).
func expectDelivered(t *testing.T, sub *client, subj string, pub *client, payload string) {
	t.Helper()
	s, err := sub.nc.SubscribeSync(subj)
	if err != nil {
		t.Fatal(err)
	}
	defer s.Unsubscribe()
	sub.nc.Flush()
	time.Sleep(50 * time.Millisecond)
	if err := pub.nc.Publish(subj, []byte(payload)); err != nil {
		t.Fatal(err)
	}
	pub.nc.Flush()
	m, err := s.NextMsg(2 * time.Second)
	if err != nil {
		t.Fatalf("%s published %s but %s received nothing: %v", pub.user, subj, sub.user, err)
	}
	if e, bad := sub.sink.violation("subscription", subj, settle); bad {
		t.Fatalf("%s subscribe %s: %s", sub.user, subj, e)
	}
	pass(t, "ALLOW recv %-20s %-40s from %s (%d bytes)", sub.user, subj, pub.user, len(m.Data))
}

// The browser and Unreal clients (onemind-xr M4 WEB and UE) read live action status, mission
// facts, zone crossings, feed freshness and GeoChat through CLIENT stream imports, and call the
// task, alert, video and TAK mission services through service imports with a server-set identity.
func TestClientDisplaysAndServices(t *testing.T) {
	h := startHub(t, "B")
	provision(t, h)
	zeus := h.conn(t, "zeus")
	geochat := `<event version="2.0" uid="GeoChat.ANDROID-1.All Chat Rooms.1" type="b-t-f" how="h-g-i-g-o"><detail><__chat chatroom="All Chat Rooms"/><remarks>copy that</remarks></detail></event>`

	t.Run("displays", func(t *testing.T) {
		expectDelivered(t, zeus, "evt.so.action.executed", h.conn(t, "policy"), `{"action_id":"a1"}`)
		expectDelivered(t, zeus, "evt.mission.m1.completed", h.conn(t, "mission"), `{"mission_id":"m1"}`)
		expectDelivered(t, zeus, "evt.map.geofence.home.enter", h.conn(t, "geofence"), `{"zone":"home"}`)
		expectDelivered(t, zeus, "tel.feed.nws-alerts.status", h.conn(t, "ingest"), `{"ok":true}`)
		expectDelivered(t, zeus, "tak.cot.out", h.conn(t, "tak-client"), geochat)
		expectDelivered(t, zeus, "tak.cot.mesh.rx", h.conn(t, "mesh-gateway"), geochat)
		// Wildcard subscriptions the clients use (J1 main.ts, UE FeedGaps).
		expectDelivered(t, zeus, "evt.map.geofence.*.*", h.conn(t, "geofence"), `{}`)
		expectDelivered(t, zeus, "tel.feed.*.status", h.conn(t, "ingest"), `{}`)
		// Not imported: other evt.so.* and the TAK inject path stay invisible.
		expectSub(t, zeus, "evt.so.alert.>", false)
		expectSub(t, zeus, "tak.cot.inject", false)
		expectSub(t, zeus, "tel.fleet.>", false)
	})

	t.Run("services", func(t *testing.T) {
		serve := func(user, subj string) *client {
			c := h.conn(t, user)
			if _, err := c.nc.Subscribe(subj, func(m *nats.Msg) {
				_ = m.Respond([]byte(m.Header.Get("Nats-Request-Info") + "|" + m.Subject))
			}); err != nil {
				t.Fatal(err)
			}
			c.nc.Flush()
			return c
		}
		serve("mission", "svc.mission.>")
		serve("geofence", "svc.geofence.alerts")
		serve("video-registry", "svc.registry.video.>")
		time.Sleep(100 * time.Millisecond)
		for _, subj := range []string{"svc.mission.get", "svc.mission.list", "svc.mission.report", "svc.geofence.alerts", "svc.registry.video.get", "svc.registry.video.list"} {
			r, err := zeus.nc.Request(subj, []byte(`{}`), 2*time.Second)
			if err != nil {
				t.Fatalf("zeus %s: %v", subj, err)
			}
			if !strings.Contains(string(r.Data), `"acc":"CLIENT"`) || !strings.Contains(string(r.Data), `"user":"zeus"`) {
				t.Fatalf("zeus %s: service saw %q, want the server-set CLIENT/zeus identity", subj, r.Data)
			}
			pass(t, "ALLOW req  zeus -> %-28s answered; service saw CLIENT/zeus", subj)
		}
		// svc.mission.apply is policy-only: no CLIENT grant, no export, and the shared user is denied.
		expectPub(t, zeus, "svc.mission.apply", false, nil)
		expectPub(t, h.conn(t, "onemind"), "svc.mission.apply", false, nil)
		if r, err := h.conn(t, "policy").nc.Request("svc.mission.apply", []byte(`{}`), 2*time.Second); err != nil {
			t.Fatalf("policy svc.mission.apply: %v", err)
		} else {
			pass(t, "ALLOW req  policy -> svc.mission.apply answered (same account, no Nats-Request-Info: %q)", r.Data)
		}
		// svc.tak.mission.> is imported for UE-P2 but nobody serves it yet.
		if _, err := zeus.nc.Request("svc.tak.mission.list", []byte(`{}`), time.Second); !errors.Is(err, nats.ErrNoResponders) {
			t.Fatalf("zeus svc.tak.mission.list: want no responders, got %v", err)
		}
		if e, bad := zeus.sink.violation("publish", "svc.tak.mission.list", settle); bad {
			t.Fatalf("zeus svc.tak.mission.list: %s", e)
		}
		pass(t, "ALLOW req  zeus -> svc.tak.mission.list permitted; no responders until tak-client serves UE-P2")
	})

	t.Run("cmd.map stays in CLIENT", func(t *testing.T) {
		// A client may publish cmd.map.> (reversible UI) but a zone write never reaches the geofence
		// service or CMD: the subject lives in CLIENT, which exports it only to AGENTS as a service.
		om := h.conn(t, "onemind")
		s, _ := om.nc.SubscribeSync("cmd.map.>")
		om.nc.Flush()
		time.Sleep(50 * time.Millisecond)
		zeus.nc.Publish("cmd.map.geofence.home.create", []byte(`{"forged":true}`))
		zeus.nc.Flush()
		if m, err := s.NextMsg(500 * time.Millisecond); err == nil {
			t.Fatalf("CLIENT cmd.map.geofence.home.create reached ONEMIND: %q", m.Data)
		}
		js, _ := jetstream.New(om.nc)
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		st, err := js.Stream(ctx, "CMD")
		if err != nil {
			t.Fatal(err)
		}
		if n := st.CachedInfo().State.Msgs; n != 0 {
			t.Fatalf("CMD holds %d messages after a CLIENT cmd.map publish", n)
		}
		pass(t, "DENY  route zeus cmd.map.geofence.home.create: not in ONEMIND, CMD holds 0 messages")
	})
}

type roleFile struct {
	Roles []struct {
		Name              string   `json:"name"`
		Publish           []string `json:"publish"`
		Subscribe         []string `json:"subscribe"`
		ProposedPublish   []string `json:"proposed_publish"`
		ProposedSubscribe []string `json:"proposed_subscribe"`
	} `json:"roles"`
}

// One CLIENT user per #565 role mode, with exactly that role's allow lists plus the A19-P11
// additions (testdata/client-roles.json), as the auth callout would issue them in its user JWT.
// newcomer and viewer can never command, forge entities, inject CoT or confirm; each role can
// read what its displays need through the CLIENT imports.
func TestClientRoles(t *testing.T) {
	setup(t)
	b, err := os.ReadFile(filepath.Join("testdata", "client-roles.json"))
	if err != nil {
		t.Fatal(err)
	}
	var rf roleFile
	if err := json.Unmarshal(b, &rf); err != nil {
		t.Fatal(err)
	}
	frag := renderStage(t, "B")
	var cp map[string]any
	_ = json.Unmarshal(marshalConf(frag), &cp)
	clientAcc := cp["accounts"].(map[string]any)["CLIENT"].(map[string]any)
	users := clientAcc["users"].([]any)
	for _, r := range rf.Roles {
		name := "role-" + r.Name
		raw := make([]byte, 18)
		_, _ = rand.Read(raw)
		passwords[name] = hex.EncodeToString(raw)
		pub := append(append([]string{}, r.Publish...), r.ProposedPublish...)
		sub := append(append([]string{"_INBOX.>"}, r.Subscribe...), r.ProposedSubscribe...)
		if len(pub) == 0 {
			pub = []string{"_NO_PUBLISH_"}
		}
		users = append(users, map[string]any{
			"user": name, "password": passwords[name],
			"allowed_connection_types": []any{"STANDARD", "WEBSOCKET"},
			"permissions": map[string]any{
				"publish":         map[string]any{"allow": pub},
				"subscribe":       map[string]any{"allow": sub},
				"allow_responses": map[string]any{"max": -1, "expires": "5s"},
			},
		})
	}
	clientAcc["users"] = users
	h := newHub(t, "B-roles", cp)
	provision(t, h)
	policy := h.conn(t, "policy")
	if _, err := policy.nc.Subscribe("svc.policy.>", func(m *nats.Msg) { _ = m.Respond([]byte(m.Header.Get("Nats-Request-Info"))) }); err != nil {
		t.Fatal(err)
	}
	policy.nc.Flush()
	ingest := h.conn(t, "ingest")
	tak := h.conn(t, "tak-client")
	bridge := h.conn(t, "bridge-mavlink-home")
	entObs := h.conn(t, "onemind")

	low := map[string]bool{"newcomer": true, "viewer": true}
	for _, r := range rf.Roles {
		r := r
		t.Run(r.Name, func(t *testing.T) {
			c := h.conn(t, "role-"+r.Name)
			expectPub(t, c, "cmd.fleet.drone.sim-x1.goto", false, bridge)
			expectPub(t, c, "ent.drone.sim-x1.state", false, entObs)
			expectPub(t, c, "tak.cot.inject", false, nil)
			expectPub(t, c, "evt.so.action.approved", false, nil)
			expectPub(t, c, "svc.mission.apply", false, nil)
			expectPub(t, c, "$JS.API.INFO", false, nil)
			expectDelivered(t, c, "tel.feed.*.status", ingest, `{"ok":true}`)
			if low[r.Name] {
				expectPub(t, c, "svc.policy.command.confirm", false, nil)
				expectPub(t, c, "svc.policy.task.confirm", false, nil)
				expectPub(t, c, "svc.policy.command.propose", false, nil)
				expectPub(t, c, "cmd.map.view.fly", false, nil)
				expectSub(t, c, "tak.cot.out", false)
			} else {
				if _, err := c.nc.Request("svc.policy.command.confirm", []byte(`{}`), 2*time.Second); err != nil {
					t.Fatalf("%s confirm: %v", r.Name, err)
				}
				pass(t, "ALLOW req  role-%s -> svc.policy.command.confirm reaches policy (humans confirm)", r.Name)
				expectDelivered(t, c, "tak.cot.out", tak, `<event uid="chat"/>`)
			}
			if r.Name == "newcomer" {
				expectSub(t, c, "evt.so.action.>", false)
				expectSub(t, c, "evt.mission.>", false)
				expectPub(t, c, "svc.policy.action.get", false, nil)
			} else {
				expectDelivered(t, c, "evt.so.action.executed", policy, `{}`)
				expectSub(t, c, "evt.mission.>", true)
			}
		})
	}
}
