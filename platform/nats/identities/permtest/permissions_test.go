package permtest

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

// streams mirrors the subjects of the omos stream CRs the services use (platform/nats/streams/
// plus ENT, DET and the KV buckets proposed in #552/#553/#554). Created by the shared `onemind`
// user, as NACK does today.
var streams = map[string][]string{
	"COT":   {"tak.cot.>"},
	"TEL":   {"tel.>"},
	"WORLD": {"evt.world.>"},
	"INTEL": {"intel.>"},
	"DET":   {"det.>"},
	"ENT":   {"ent.*.*.state"},
	"CMD":   {"cmd.fleet.>", "cmd.mission.>", "cmd.entity.>"},
	"EVT":   {"evt.so.>", "evt.drone.>", "evt.entity.>", "evt.mission.>"},
	"SYS":   {"evt.sys.>"},
}

var buckets = []string{"geofences", "mission-state", "entity-correlation", "ingest-locks", "ingest-state", "feed-config", "world-state", "fleet-state"}

func provision(t *testing.T, h *hub) {
	t.Helper()
	admin := h.conn(t, "onemind")
	js, _ := jetstream.New(admin.nc)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	for name, subj := range streams {
		if _, err := js.CreateStream(ctx, jetstream.StreamConfig{Name: name, Subjects: subj, Storage: jetstream.FileStorage}); err != nil {
			t.Fatalf("create stream %s: %v", name, err)
		}
	}
	for _, b := range buckets {
		kv, err := js.CreateKeyValue(ctx, jetstream.KeyValueConfig{Bucket: b, Storage: jetstream.FileStorage})
		if err != nil {
			t.Fatalf("create kv %s: %v", b, err)
		}
		if _, err := kv.Put(ctx, "seed", []byte(`{"seed":true}`)); err != nil {
			t.Fatal(err)
		}
	}
}

// observerFor picks a hub user that may subscribe subj, to prove an allowed publish is delivered
// (through the account imports where the subject crosses accounts).
func observerFor(t *testing.T, h *hub, publisherAccount, subj string, cache map[string]*client) *client {
	t.Helper()
	pick := "onemind"
	switch {
	case strings.HasPrefix(subj, "svc.policy."):
		pick = "policy"
	case strings.HasPrefix(subj, "svc.entity."):
		pick = "entity"
	case strings.HasPrefix(subj, "svc.layer."):
		pick = "layer"
	case strings.HasPrefix(subj, "svc.cot."):
		pick = "projector"
	case strings.HasPrefix(subj, "cmd.map."):
		pick = "zeus" // second CLIENT connection
	case strings.HasPrefix(subj, "agents.prompt.hermes."):
		return nil // served in AGENTS (leaf); covered by TestHermesLeaf
	}
	if publisherAccount == "CLIENT" && pick == "onemind" {
		return nil
	}
	if c, ok := cache[pick]; ok {
		return c
	}
	c := h.conn(t, pick)
	cache[pick] = c
	return c
}

type row struct {
	user, account                        string
	pubAllow, pubDeny, subAllow, subDeny []string
}

// matrix is the stage-B (target) allow/deny matrix for every directly connecting user.
var matrix = []row{
	{"policy", "ONEMIND",
		[]string{"cmd.fleet.drone.sim-x1.goto", "cmd.mission.m1.create", "evt.so.action.approved", "svc.entity.get", "$KV.geofences.z1"},
		[]string{"tak.cot.inject", "ent.drone.sim-x1.state", "intel.weather.noaa.alert", "_INBOX.spoof", "agents.prompt.hermes.zeus.legacy", "$JS.API.STREAM.CREATE.X", "$KV.entity-correlation.k"},
		[]string{"svc.policy.>", "evt.drone.>", "evt.mission.>", "_INBOX.>"},
		[]string{"tel.>", "intel.>", "cmd.>", "svc.entity.>"}},
	{"entity", "ONEMIND",
		[]string{"ent.drone.sim-x1.state", "evt.entity.e1.created", "evt.sys.entity.rejected", "$KV.entity-correlation.k"},
		[]string{"cmd.fleet.drone.sim-x1.goto", "cmd.entity.e1.override", "tak.cot.inject", "evt.so.action.approved", "intel.a.b.c", "svc.policy.command.confirm", "_INBOX.spoof", "$JS.API.STREAM.CREATE.X"},
		[]string{"svc.entity.>", "cmd.entity.>", "_INBOX.>", "$SRV.PING", "$SRV.INFO.entity", "$SRV.STATS.entity.abc123"},
		[]string{"svc.policy.>", "tak.cot.>", "cmd.fleet.>", "$SRV.PING.policy"}},
	{"ingest", "ONEMIND",
		[]string{"intel.weather.noaa.alert", "evt.world.sea.mmsi-367000001", "tel.feed.nws-alerts.status", "evt.sys.nws-alerts.rejected", "$KV.ingest-locks.nws-alerts"},
		[]string{"ent.ship.mmsi-1.state", "cmd.fleet.drone.sim-x1.goto", "tak.cot.inject", "evt.so.action.approved",
			"svc.policy.command.propose", "tel.fleet.drone.sim-x1.position",
			"$KV.feed-config.nws",
			"$JS.API.STREAM.DELETE.INTEL"},
		[]string{"_INBOX.>"},
		[]string{"intel.>", "svc.policy.>", "cmd.>"}},
	{"layer", "ONEMIND",
		[]string{"$KV.world-state.air", "$KV.fleet-state.sim-x1"},
		[]string{"ent.drone.sim-x1.state", "cmd.map.view.fly", "tak.cot.inject", "svc.layer.snapshot", "$JS.API.STREAM.MSG.DELETE.ENT"},
		[]string{"svc.layer.>", "tel.feed.>", "_INBOX.>"},
		[]string{"svc.policy.>", "cmd.>", "intel.>"}},
	{"projector", "ONEMIND",
		[]string{"tak.cot.inject", "tak.cot.mesh.tx"},
		[]string{"cmd.fleet.drone.sim-x1.goto", "ent.drone.sim-x1.state", "tak.cot.out", "$KV.geofences.z1"},
		[]string{"svc.cot.send", "ent.*.*.state", "cmd.projector.disable", "_INBOX.>"},
		[]string{"svc.policy.>", "cmd.>", "cmd.fleet.>"}},
	{"tak-client", "ONEMIND",
		[]string{"tak.cot.out"},
		[]string{"tak.cot.inject", "cmd.fleet.drone.sim-x1.goto", "ent.drone.sim-x1.state", "tak.cot.mesh.rx"},
		[]string{"_INBOX.>"},
		[]string{"tak.cot.inject", "tak.cot.>"}},
	{"bridge-mavlink-home", "ONEMIND",
		[]string{"tel.fleet.drone.sim-x1.position", "tel.fleet.drone.sim-x1.status", "evt.drone.sim-x1.command"},
		[]string{"tel.fleet.drone.sim-x2.position", "evt.drone.sim-x2.command", "cmd.fleet.drone.sim-x1.goto", "ent.drone.sim-x1.state", "tak.cot.inject"},
		[]string{"cmd.fleet.drone.sim-x1.>", "_INBOX.>"},
		[]string{"cmd.fleet.drone.sim-x2.>", "cmd.fleet.>", "cmd.>"}},
	{"mesh-gateway", "ONEMIND",
		[]string{"tak.cot.mesh.rx", "tel.fleet.edge.tak-mesh.status"},
		[]string{"tak.cot.inject", "tak.cot.mesh.tx", "tel.fleet.edge.other.status", "cmd.fleet.drone.sim-x1.goto"},
		[]string{"tak.cot.mesh.tx", "_INBOX.>"},
		[]string{"tak.cot.inject", "tak.cot.mesh.rx"}},
	{"agent-fabric", "ONEMIND",
		[]string{"agents.prompt.claude.zeus.x", "$SRV.PING"},
		[]string{"cmd.fleet.drone.sim-x1.goto", "cmd.map.view.fly", "tak.cot.inject", "svc.policy.command.confirm"},
		[]string{"agents.>", "_INBOX.>"},
		[]string{"cmd.>", "intel.>"}},
	{"zeus", "CLIENT",
		[]string{"svc.policy.command.propose", "svc.policy.command.confirm", "svc.entity.get", "svc.layer.snapshot", "svc.cot.send", "cmd.map.view.fly"},
		[]string{"cmd.fleet.drone.sim-x1.goto", "cmd.fleet.drone.sim-x1.land", "ent.drone.sim-x1.state", "ent.aircraft.icao-a1b2c3.state", "tak.cot.inject", "intel.a.b.c", "evt.so.action.approved", "$JS.API.INFO", "$KV.geofences.z1"},
		[]string{"cmd.map.>", "_INBOX.>"},
		[]string{"tel.>", "intel.>", "ent.>", "cmd.fleet.>", ">"}},
	{"onemind", "ONEMIND",
		[]string{"intel.legacy.x.y", "tel.drone.x47.position", "evt.world.air.icao-1"},
		[]string{"cmd.fleet.drone.sim-x1.goto", "cmd.drone.x47.land", "tak.cot.inject", "ent.drone.sim-x1.state", "evt.so.action.approved", "svc.policy.command.confirm"},
		[]string{"intel.>", "tel.>"},
		[]string{"svc.policy.>", "svc.entity.>", "svc.layer.>", "svc.cot.send"}},
}

func TestStageBMatrix(t *testing.T) {
	h := startHub(t, "B")
	provision(t, h)
	obs := map[string]*client{}
	parent := t
	for _, r := range matrix {
		r := r
		t.Run(r.user, func(t *testing.T) {
			c := h.conn(t, r.user)
			for _, s := range r.pubAllow {
				expectPub(t, c, s, true, observerFor(parent, h, r.account, s, obs))
			}
			for _, s := range r.pubDeny {
				var o *client
				if !strings.HasPrefix(s, "$JS.API") && !strings.HasPrefix(s, "_INBOX") {
					o = observerFor(parent, h, r.account, s, obs)
					if o != nil && (o.user == r.user) {
						o = nil
					}
				}
				expectPub(t, c, s, false, o)
			}
			for _, s := range r.subAllow {
				expectSub(t, c, s, true)
			}
			for _, s := range r.subDeny {
				expectSub(t, c, s, false)
			}
		})
	}
}

// A19 acceptance (FEATURES A19; asset-command.md §9 test 8): a client user cannot publish
// cmd.fleet.>, and cannot publish ent.*; each refusal is a server permission violation and
// nothing reaches the bridge or the ENT stream.
func TestClientCannotCommandOrForgeEntities(t *testing.T) {
	h := startHub(t, "B")
	provision(t, h)
	zeus := h.conn(t, "zeus")
	bridge := h.conn(t, "bridge-mavlink-home")
	entityObs := h.conn(t, "onemind")
	expectPub(t, zeus, "cmd.fleet.drone.sim-x1.goto", false, bridge)
	expectPub(t, zeus, "cmd.fleet.drone.sim-x1.takeoff", false, bridge)
	expectPub(t, zeus, "ent.drone.sim-x1.state", false, entityObs)
	// The ENT stream is empty: nothing was stored from the client.
	js, _ := jetstream.New(entityObs.nc)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	st, err := js.Stream(ctx, "ENT")
	if err != nil {
		t.Fatal(err)
	}
	if n := st.CachedInfo().State.Msgs; n != 0 {
		t.Fatalf("ENT has %d messages after denied client publishes", n)
	}
	pass(t, "ENT stream holds 0 messages after the denied client publishes")
	cst, _ := js.Stream(ctx, "CMD")
	if n := cst.CachedInfo().State.Msgs; n != 0 {
		t.Fatalf("CMD has %d messages after denied client publishes", n)
	}
	pass(t, "CMD stream holds 0 messages after the denied client publishes")
	// The client account has no JetStream at all.
	zjs, _ := jetstream.New(zeus.nc)
	if _, err := zjs.AccountInfo(ctx); err == nil {
		t.Fatal("CLIENT has JetStream access")
	} else {
		pass(t, "CLIENT JetStream AccountInfo refused: %v", err)
	}
}

// tctx is a fresh 3 s context per JetStream call, cancelled when the test ends.
func tctx(t *testing.T) context.Context {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	t.Cleanup(cancel)
	return ctx
}

func mustJS(t *testing.T, what string, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("%s: %v", what, err)
	}
	pass(t, "ALLOW js   %s", what)
}

func mustNotJS(t *testing.T, c *client, what string, err error, apiSubj string) {
	t.Helper()
	if err == nil {
		t.Fatalf("%s: expected refusal", what)
	}
	e, ok := c.sink.violation("publish", apiSubj, time.Second)
	if !ok {
		t.Fatalf("%s: failed (%v) but no permission violation on %s", what, err, apiSubj)
	}
	pass(t, "DENY  js   %s -> server: %s", what, e)
}

// JetStream grants are per stream and bucket: each service can use exactly what it needs.
func TestJetStreamGrants(t *testing.T) {
	h := startHub(t, "B")
	provision(t, h)
	admin := h.conn(t, "onemind")
	if err := admin.nc.Publish("tak.cot.out", []byte(`<event uid="a"/>`)); err != nil {
		t.Fatal(err)
	}
	admin.nc.Flush()

	t.Run("entity", func(t *testing.T) {
		c := h.conn(t, "entity")
		js, _ := jetstream.New(c.nc)
		st, err := js.Stream(tctx(t), "COT")
		mustJS(t, "entity: stream info COT", err)
		cons, err := st.CreateOrUpdateConsumer(tctx(t), jetstream.ConsumerConfig{Durable: "entity-svc-cot", FilterSubjects: []string{"tak.cot.>"}, AckPolicy: jetstream.AckExplicitPolicy})
		mustJS(t, "entity: durable consumer entity-svc-cot on COT", err)
		batch, err := cons.Fetch(1, jetstream.FetchMaxWait(2*time.Second))
		mustJS(t, "entity: fetch from COT", err)
		got := 0
		for m := range batch.Messages() {
			mustJS(t, "entity: ack COT message", m.Ack())
			got++
		}
		if got != 1 {
			t.Fatalf("entity fetched %d messages from COT", got)
		}
		ack, err := js.Publish(tctx(t), "ent.drone.sim-x1.state", []byte(`{"id":"sim-x1"}`))
		mustJS(t, "entity: acked publish into ENT", err)
		ent, _ := js.Stream(tctx(t), "ENT")
		mustJS(t, "entity: delete ENT message by sequence", ent.DeleteMsg(tctx(t), ack.Sequence))
		_, err = js.OrderedConsumer(tctx(t), "ENT", jetstream.OrderedConsumerConfig{FilterSubjects: []string{"ent.>"}, DeliverPolicy: jetstream.DeliverLastPerSubjectPolicy})
		mustJS(t, "entity: ordered consumer on ENT (restore)", err)
		kv, err := js.KeyValue(tctx(t), "entity-correlation")
		mustJS(t, "entity: bind KV entity-correlation", err)
		_, err = kv.Put(tctx(t), "lease.entity-service", []byte("holder"))
		mustJS(t, "entity: KV put entity-correlation", err)
		_, err = kv.Get(tctx(t), "lease.entity-service")
		mustJS(t, "entity: KV get entity-correlation", err)

		_, err = js.CreateStream(tctx(t), jetstream.StreamConfig{Name: "EVIL", Subjects: []string{"x.>"}})
		mustNotJS(t, c, "entity: create stream", err, "$JS.API.STREAM.CREATE.EVIL")
		_, err = js.Stream(tctx(t), "CMD")
		mustNotJS(t, c, "entity: stream info CMD", err, "$JS.API.STREAM.INFO.CMD")
		_, err = js.KeyValue(tctx(t), "geofences")
		mustNotJS(t, c, "entity: bind KV geofences", err, "$JS.API.STREAM.INFO.KV_geofences")
	})

	t.Run("ingest", func(t *testing.T) {
		c := h.conn(t, "ingest")
		js, _ := jetstream.New(c.nc)
		_, err := js.Publish(tctx(t), "intel.weather.noaa.alert", []byte(`{}`), jetstream.WithMsgID("nws-1"))
		mustJS(t, "ingest: acked publish into INTEL with Nats-Msg-Id", err)
		locks, err := js.KeyValue(tctx(t), "ingest-locks")
		mustJS(t, "ingest: bind KV ingest-locks", err)
		_, err = locks.Create(tctx(t), "nws-alerts", []byte("pod-1"))
		mustJS(t, "ingest: KV create (lease) ingest-locks", err)
		state, err := js.KeyValue(tctx(t), "ingest-state")
		mustJS(t, "ingest: bind KV ingest-state", err)
		_, err = state.Put(tctx(t), "nws-alerts", []byte(`{"etag":"x"}`))
		mustJS(t, "ingest: KV put ingest-state", err)
		fc, err := js.KeyValue(tctx(t), "feed-config")
		mustJS(t, "ingest: bind KV feed-config (read)", err)
		_, err = fc.Get(tctx(t), "seed")
		mustJS(t, "ingest: KV get feed-config", err)
		_, err = fc.Put(tctx(t), "nws", []byte("x"))
		mustNotJS(t, c, "ingest: KV put feed-config", err, "$KV.feed-config.nws")
		err = js.DeleteStream(tctx(t), "INTEL")
		mustNotJS(t, c, "ingest: delete stream INTEL", err, "$JS.API.STREAM.DELETE.INTEL")
		_, err = js.CreateKeyValue(tctx(t), jetstream.KeyValueConfig{Bucket: "evil"})
		mustNotJS(t, c, "ingest: create KV bucket", err, "$JS.API.STREAM.CREATE.KV_evil")
	})

	t.Run("policy", func(t *testing.T) {
		c := h.conn(t, "policy")
		js, _ := jetstream.New(c.nc)
		_, err := js.Publish(tctx(t), "cmd.fleet.drone.sim-x1.goto", []byte(`{}`), jetstream.WithMsgID("cmd-1"))
		mustJS(t, "policy: acked publish into CMD", err)
		_, err = js.Publish(tctx(t), "evt.so.action.approved", []byte(`{}`))
		mustJS(t, "policy: acked publish into EVT", err)
		gf, err := js.KeyValue(tctx(t), "geofences")
		mustJS(t, "policy: bind KV geofences", err)
		w, err := gf.WatchAll(tctx(t))
		mustJS(t, "policy: KV watch geofences", err)
		select {
		case e := <-w.Updates():
			if e != nil && e.Key() != "seed" {
				t.Fatalf("watch got %s", e.Key())
			}
			pass(t, "ALLOW js   policy: watch delivered geofences entry")
		case <-time.After(2 * time.Second):
			t.Fatal("policy watch delivered nothing")
		}
		w.Stop()
		_, err = gf.Put(tctx(t), "zone-1", []byte(`{}`))
		mustJS(t, "policy: KV put geofences", err)
		ms, err := js.KeyValue(tctx(t), "mission-state")
		mustJS(t, "policy: bind KV mission-state", err)
		_, err = ms.Put(tctx(t), "m1", []byte(`{}`))
		mustJS(t, "policy: KV put mission-state", err)
		_, err = js.Stream(tctx(t), "ENT")
		mustNotJS(t, c, "policy: stream info ENT", err, "$JS.API.STREAM.INFO.ENT")
	})

	t.Run("projector", func(t *testing.T) {
		c := h.conn(t, "projector")
		js, _ := jetstream.New(c.nc)
		st, err := js.Stream(tctx(t), "INTEL")
		mustJS(t, "projector: stream info INTEL", err)
		_, err = st.CreateOrUpdateConsumer(tctx(t), jetstream.ConsumerConfig{Durable: "sink-tak", AckPolicy: jetstream.AckExplicitPolicy})
		mustJS(t, "projector: durable consumer sink-tak on INTEL", err)
		_, err = js.StreamNameBySubject(tctx(t), "intel.>")
		mustJS(t, "projector: stream lookup by subject (nats-py js.subscribe)", err)
		gf, err := js.KeyValue(tctx(t), "geofences")
		mustJS(t, "projector: bind KV geofences (read)", err)
		_, err = gf.Get(tctx(t), "seed")
		mustJS(t, "projector: KV get geofences", err)
		_, err = gf.Put(tctx(t), "zone-x", []byte("x"))
		mustNotJS(t, c, "projector: KV put geofences", err, "$KV.geofences.zone-x")
	})

	t.Run("layer", func(t *testing.T) {
		c := h.conn(t, "layer")
		js, _ := jetstream.New(c.nc)
		_, err := js.OrderedConsumer(tctx(t), "ENT", jetstream.OrderedConsumerConfig{FilterSubjects: []string{"ent.aircraft.>"}, DeliverPolicy: jetstream.DeliverLastPerSubjectPolicy})
		mustJS(t, "layer: ordered consumer on ENT (snapshot path)", err)
		ws, err := js.KeyValue(tctx(t), "world-state")
		mustJS(t, "layer: bind KV world-state", err)
		_, err = ws.Put(tctx(t), "air", []byte(`{}`))
		mustJS(t, "layer: KV put world-state", err)
		ent, _ := js.Stream(tctx(t), "ENT")
		err = ent.DeleteMsg(tctx(t), 1)
		mustNotJS(t, c, "layer: delete ENT message", err, "$JS.API.STREAM.MSG.DELETE.ENT")
	})
}

// Cross-account requests: the server stamps the caller's identity (share: true) and replaces a
// forged header; streamed svc.layer.watch responses reach a CLIENT caller.
func TestCrossAccountRequests(t *testing.T) {
	h := startHub(t, "B")
	policy := h.conn(t, "policy")
	_, err := policy.nc.Subscribe("svc.policy.>", func(m *nats.Msg) {
		_ = m.Respond([]byte(m.Header.Get("Nats-Request-Info") + "|" + m.Subject))
	})
	if err != nil {
		t.Fatal(err)
	}
	policy.nc.Flush()

	zeus := h.conn(t, "zeus")
	req := nats.NewMsg("svc.policy.command.propose")
	req.Header.Set("Nats-Request-Info", `{"acc":"ONEMIND","user":"policy"}`)
	req.Data = []byte(`{}`)
	resp, err := zeus.nc.RequestMsg(req, 2*time.Second)
	if err != nil {
		t.Fatalf("CLIENT propose: %v", err)
	}
	var info struct {
		Acc  string `json:"acc"`
		User string `json:"user"`
	}
	hdr, _, _ := strings.Cut(string(resp.Data), "|")
	_ = json.Unmarshal([]byte(hdr), &info)
	if info.Acc != "CLIENT" || info.User != "zeus" {
		t.Fatalf("policy saw identity %q, want CLIENT/zeus", hdr)
	}
	pass(t, "CLIENT zeus -> svc.policy.command.propose: policy saw server-set identity acc=%s user=%s (forged header replaced)", info.Acc, info.User)
	resp, err = zeus.nc.Request("svc.policy.command.confirm", []byte(`{}`), 2*time.Second)
	if err != nil {
		t.Fatalf("CLIENT confirm: %v", err)
	}
	pass(t, "CLIENT zeus -> svc.policy.command.confirm reaches the policy service (humans may confirm): %s", resp.Data)

	// An ONEMIND service cannot reach svc.policy.*.confirm (only the policy service subscribes,
	// and stage B denies the shared user).
	onemind := h.conn(t, "onemind")
	expectPub(t, onemind, "svc.policy.command.confirm", false, nil)

	// Streamed responses through the CLIENT import (svc.layer.watch, layer-delta.md §3).
	layer := h.conn(t, "layer")
	_, err = layer.nc.Subscribe("svc.layer.watch", func(m *nats.Msg) {
		for i := 0; i < 5; i++ {
			_ = layer.nc.Publish(m.Reply, []byte{byte('0' + i)})
		}
	})
	if err != nil {
		t.Fatal(err)
	}
	layer.nc.Flush()
	inbox := nats.NewInbox()
	sub, _ := zeus.nc.SubscribeSync(inbox)
	if err := zeus.nc.PublishRequest("svc.layer.watch", inbox, []byte(`{"v":1}`)); err != nil {
		t.Fatal(err)
	}
	for i := 0; i < 5; i++ {
		if _, err := sub.NextMsg(2 * time.Second); err != nil {
			t.Fatalf("watch batch %d not received: %v", i, err)
		}
	}
	if e, bad := layer.sink.violation("publish", "", settle); bad {
		t.Fatalf("layer violation while streaming: %s", e)
	}
	pass(t, "CLIENT zeus <- svc.layer.watch: 5 streamed responses through the import (layer has no _INBOX publish grant)")

	// A responder without allow_responses for the caller's inbox cannot spoof replies.
	ingest := h.conn(t, "ingest")
	expectPub(t, ingest, inbox, false, nil)
}

// FEATURES A19 acceptance: Hermes (AGENTS, its own leaf credential) cannot publish
// tak.cot.inject or cmd.fleet.> (nor any cmd.* except cmd.map.>), cannot confirm, cannot read
// the raw firehoses, and each refusal is a server permission violation; it can propose, read
// entities, serve prompts and drive cmd.map.>.
//
// The leaf is configured as the sealed sidecar leaf.conf is (`nats_identities.py leaf-conf`):
// the agent logs in to a local user carrying the hub user's allow lists. The hub alone only
// filters (it advertises to a leaf just the subjects that user may publish, so a forbidden
// publish is dropped on the spoke with no error); the sidecar's local user turns each attempt
// into a `Permissions Violation` for the agent and a `Publish Violation` line in the sidecar log.
func TestHermesLeaf(t *testing.T) {
	h := startHub(t, "B")
	l := h.startLeafWith(t, "hermes", renderSidecar(t, "hermes"))
	if !waitLeafs(t, h.srv, 1, 5*time.Second) {
		t.Fatalf("hermes leaf did not connect: %v", h.log.grep("hermes"))
	}
	pass(t, "hermes leaf connected to the hub (account AGENTS)")
	hermes := l.conn(t, "hermes-agent")
	bridge := h.conn(t, "bridge-mavlink-home")
	tak := h.startLeaf(t, "leaf-tak")
	waitLeafs(t, h.srv, 2, 5*time.Second)
	takPlugin := tak.conn(t, "tak-plugin")
	onemind := h.conn(t, "onemind")
	zeus := h.conn(t, "zeus")

	// Denied: commands, TAK inject, firehoses.
	// Observers in ONEMIND: the bridge (cmd.fleet.drone.sim-x1.>), the TAK leaf (tak.cot.inject), onemind (cmd.ha.*).
	expectLeafPubTo := func(subj string, sub *client) {
		t.Helper()
		s, err := sub.nc.SubscribeSync(subj)
		if err != nil {
			t.Fatal(err)
		}
		defer s.Unsubscribe()
		sub.nc.Flush()
		time.Sleep(200 * time.Millisecond)
		hermes.nc.Publish(subj, []byte(`{"from":"hermes"}`))
		hermes.nc.Flush()
		e, ok := hermes.sink.violation("publish", subj, 2*time.Second)
		if !ok {
			t.Fatalf("hermes publish %s: expected a server permission violation at the sidecar", subj)
		}
		var logged string
		for _, line := range l.log.grep("Publish Violation") {
			if strings.Contains(line, fmt.Sprintf("%q", subj)) {
				logged = line
			}
		}
		if logged == "" {
			t.Fatalf("hermes publish %s: no Publish Violation line in the sidecar log", subj)
		}
		if m, err := s.NextMsg(700 * time.Millisecond); err == nil {
			t.Fatalf("hermes published %s and %s received %q", subj, sub.user, m.Data)
		}
		pass(t, "DENY  leaf-pub hermes            %-40s client: %s; sidecar log: %s; not delivered (%s listening)", subj, e, logged, sub.user)
	}
	expectLeafPubTo("cmd.fleet.drone.sim-x1.goto", bridge)
	expectLeafPubTo("cmd.fleet.drone.sim-x1.land", bridge)
	expectLeafPubTo("tak.cot.inject", takPlugin)
	expectLeafPubTo("cmd.ha.light.on", onemind)
	expectLeafPubTo("cmd.mission.m1.create", onemind)
	expectLeafPubTo("ent.drone.sim-x1.state", onemind)
	expectLeafPubTo("evt.so.action.approved", onemind)
	for _, subj := range []string{"intel.weather.noaa.alert", "tel.drone.x47.position", "tak.cot.out"} {
		expectLeafSub(t, hermes, subj, false, onemind)
		if e, ok := hermes.sink.violation("subscription", subj, 2*time.Second); !ok {
			t.Fatalf("hermes subscribe %s: expected a server permission violation at the sidecar", subj)
		} else {
			pass(t, "DENY  leaf-sub hermes            %-40s client: %s", subj, e)
		}
	}

	// The leaf stays up: the sidecar refused locally and nothing was forwarded, so the hub never
	// had to close the connection.
	if h.srv.NumLeafNodes() != 2 {
		t.Fatalf("leaf count %d after denied publishes", h.srv.NumLeafNodes())
	}
	if !hermes.nc.IsConnected() {
		t.Fatal("hermes agent connection closed after the violations")
	}
	pass(t, "hermes leaf and agent connection still up after the denied publishes (refused at the sidecar, no disconnect)")

	// Allowed: proposals with server-set AGENTS identity; confirm is unreachable.
	policy := h.conn(t, "policy")
	seen := make(chan string, 8)
	policy.nc.Subscribe("svc.policy.>", func(m *nats.Msg) {
		seen <- m.Subject
		_ = m.Respond([]byte(m.Header.Get("Nats-Request-Info")))
	})
	policy.nc.Flush()
	time.Sleep(200 * time.Millisecond)
	resp, err := hermes.nc.Request("svc.policy.command.propose", []byte(`{}`), 2*time.Second)
	if err != nil {
		t.Fatalf("hermes propose: %v", err)
	}
	<-seen
	pass(t, "hermes -> svc.policy.command.propose answered; policy saw Nats-Request-Info %s", resp.Data)
	if !strings.Contains(string(resp.Data), `"acc":"AGENTS"`) {
		t.Fatalf("hermes request identity %s, want account AGENTS", resp.Data)
	}
	pass(t, "policy sees Hermes as account AGENTS (policy maps every AGENTS identity to kind agent)")
	if _, err := hermes.nc.Request("svc.policy.command.confirm", []byte(`{}`), 700*time.Millisecond); err == nil {
		t.Fatal("hermes reached svc.policy.command.confirm")
	} else if e, ok := hermes.sink.violation("publish", "svc.policy.command.confirm", 2*time.Second); !ok {
		t.Fatalf("hermes svc.policy.command.confirm: %v, but no server permission violation", err)
	} else {
		pass(t, "DENY  hermes -> svc.policy.command.confirm: %v; client: %s", err, e)
	}
	if _, err := hermes.nc.Request("svc.policy.task.confirm", []byte(`{}`), 700*time.Millisecond); err == nil {
		t.Fatal("hermes reached svc.policy.task.confirm")
	} else {
		pass(t, "DENY  hermes -> svc.policy.task.confirm: %v", err)
	}
	select {
	case s := <-seen:
		t.Fatalf("policy service received %s from hermes", s)
	case <-time.After(200 * time.Millisecond):
		pass(t, "policy service received no confirm from hermes")
	}

	// Allowed: entity reads; not layer snapshots (only list/meta).
	entity := h.conn(t, "entity")
	entity.nc.Subscribe("svc.entity.>", func(m *nats.Msg) { _ = m.Respond([]byte(`{"ok":true}`)) })
	layer := h.conn(t, "layer")
	layer.nc.Subscribe("svc.layer.>", func(m *nats.Msg) { _ = m.Respond([]byte(m.Subject)) })
	entity.nc.Flush()
	layer.nc.Flush()
	time.Sleep(200 * time.Millisecond)
	if _, err := hermes.nc.Request("svc.entity.get", []byte(`{"id":"sim-x1"}`), 2*time.Second); err != nil {
		t.Fatalf("hermes svc.entity.get: %v", err)
	}
	pass(t, "ALLOW hermes -> svc.entity.get answered")
	if _, err := hermes.nc.Request("svc.layer.list", []byte(`{}`), 2*time.Second); err != nil {
		t.Fatalf("hermes svc.layer.list: %v", err)
	}
	pass(t, "ALLOW hermes -> svc.layer.list answered")
	if _, err := hermes.nc.Request("svc.layer.snapshot", []byte(`{}`), 700*time.Millisecond); err == nil {
		t.Fatal("hermes reached svc.layer.snapshot")
	} else {
		pass(t, "DENY  hermes -> svc.layer.snapshot: %v", err)
	}

	// Allowed: cmd.map.> reaches the clients; heartbeats reach ONEMIND (stream AGENTS).
	zsub, _ := zeus.nc.SubscribeSync("cmd.map.>")
	zeus.nc.Flush()
	hb, _ := onemind.nc.SubscribeSync("agents.hb.>")
	onemind.nc.Flush()
	time.Sleep(200 * time.Millisecond)
	hermes.nc.Publish("cmd.map.view.fly", []byte(`{"lat":1}`))
	hermes.nc.Publish("agents.hb.hermes.zeus.legacy", []byte(`{}`))
	hermes.nc.Flush()
	if _, err := zsub.NextMsg(2 * time.Second); err != nil {
		t.Fatal("cmd.map.view.fly from hermes did not reach the CLIENT")
	}
	pass(t, "ALLOW hermes cmd.map.view.fly delivered to CLIENT zeus (stream export AGENTS -> CLIENT)")
	if _, err := hb.NextMsg(2 * time.Second); err != nil {
		t.Fatal("hermes heartbeat did not reach ONEMIND")
	}
	pass(t, "ALLOW hermes agents.hb.hermes.zeus.legacy delivered to ONEMIND (stream export AGENTS -> ONEMIND)")

	// Hermes serves prompts, svc.bridge and $SRV discovery to ONEMIND and CLIENT callers.
	hermes.nc.Subscribe("agents.prompt.hermes.>", func(m *nats.Msg) {
		_ = m.Respond([]byte("pong " + m.Header.Get("Nats-Request-Info")))
	})
	hermes.nc.Subscribe("svc.bridge.>", func(m *nats.Msg) { _ = m.Respond([]byte("bridge")) })
	hermes.nc.Subscribe("$SRV.PING.agents", func(m *nats.Msg) { _ = m.Respond([]byte(`{"name":"agents"}`)) })
	hermes.nc.Flush()
	time.Sleep(300 * time.Millisecond)
	fabric := h.conn(t, "agent-fabric")
	if r, err := fabric.nc.Request("agents.prompt.hermes.zeus.legacy", []byte("hi"), 2*time.Second); err != nil {
		t.Fatalf("agent-fabric prompt: %v", err)
	} else {
		pass(t, "ALLOW agent-fabric (ONEMIND) -> agents.prompt.hermes.zeus.legacy: %s", r.Data)
	}
	if r, err := zeus.nc.Request("agents.prompt.hermes.zeus.legacy", []byte("hi"), 2*time.Second); err != nil {
		t.Fatalf("zeus prompt: %v", err)
	} else {
		pass(t, "ALLOW zeus (CLIENT) -> agents.prompt.hermes.zeus.legacy: %s", r.Data)
	}
	if _, err := onemind.nc.Request("svc.bridge.status", nil, 2*time.Second); err != nil {
		t.Fatalf("svc.bridge: %v", err)
	}
	pass(t, "ALLOW onemind -> svc.bridge.status served by hermes")
	if _, err := fabric.nc.Request("$SRV.PING.agents", nil, 2*time.Second); err != nil {
		t.Fatalf("$SRV.PING.agents: %v", err)
	}
	pass(t, "ALLOW agent-fabric -> $SRV.PING.agents served by hermes")
	t.Logf("hub log lines mentioning Violation: %q", h.log.grep("Violation"))
	t.Logf("hermes sidecar log lines mentioning Violation: %q", l.log.grep("Violation"))
}

// Leaf users are filtered by the hub in both directions (security-and-trust.md fact 16).
func TestLeafUsers(t *testing.T) {
	h := startHub(t, "B")
	provision(t, h)
	onemind := h.conn(t, "onemind")
	policy := h.conn(t, "policy")
	projector := h.conn(t, "projector")

	t.Run("leaf-tak", func(t *testing.T) {
		l := h.startLeaf(t, "leaf-tak")
		if !waitLeafs(t, h.srv, 1, 5*time.Second) {
			t.Fatal("leaf-tak did not connect")
		}
		p := l.conn(t, "tak-plugin")
		expectLeafSub(t, p, "tak.cot.inject", true, projector)
		expectLeafPub(t, p, "tak.cot.out", true, onemind)
		expectLeafPub(t, p, "cmd.fleet.drone.sim-x1.goto", false, h.conn(t, "bridge-mavlink-home"))
		expectLeafPub(t, p, "ent.drone.sim-x1.state", false, onemind)
		expectLeafPub(t, p, "tak.cot.inject", false, onemind)
		expectLeafSub(t, p, "intel.weather.noaa.alert", false, onemind)
		expectLeafSub(t, p, "cmd.fleet.drone.sim-x1.goto", false, policy)
	})

	t.Run("leaf-site-mac", func(t *testing.T) {
		before := h.srv.NumLeafNodes()
		l := h.startLeaf(t, "leaf-site-mac")
		if !waitLeafs(t, h.srv, before+1, 5*time.Second) {
			t.Fatal("leaf-site-mac did not connect")
		}
		b := l.conn(t, "mac-bridges")
		expectLeafSub(t, b, "cmd.fleet.drone.sim-x1.goto", true, policy)
		expectLeafSub(t, b, "cmd.fleet.drone.sim-x2.goto", false, policy)
		expectLeafSub(t, b, "tak.cot.mesh.tx", true, projector)
		expectLeafSub(t, b, "tak.cot.inject", false, projector)
		expectLeafPub(t, b, "tel.fleet.drone.sim-x1.position", true, onemind)
		expectLeafPub(t, b, "tel.fleet.drone.sim-x2.position", false, onemind)
		expectLeafPub(t, b, "evt.drone.sim-x1.command", true, onemind)
		expectLeafPub(t, b, "tak.cot.mesh.rx", true, onemind)
		expectLeafPub(t, b, "tak.cot.inject", false, onemind)
		expectLeafPub(t, b, "cmd.fleet.drone.sim-x1.goto", false, onemind)
		// The mesh gateway's acked JetStream publish into COT works through the leaf.
		js, _ := jetstream.New(b.nc)
		time.Sleep(200 * time.Millisecond)
		if _, err := js.Publish(tctx(t), "tak.cot.mesh.rx", []byte(`<event uid="m"/>`), jetstream.WithMsgID("mesh-1")); err != nil {
			t.Fatalf("JetStream publish through the site leaf: %v", err)
		}
		pass(t, "ALLOW leaf-site-mac: acked JetStream publish tak.cot.mesh.rx into COT through the leaf")
	})
}

// allowed_connection_types: leaf credentials cannot be used as client logins, service logins
// cannot use WebSocket, CLIENT users can.
func TestConnectionTypes(t *testing.T) {
	h := startHub(t, "B")
	for _, u := range []string{"hermes", "leaf-tak", "leaf-site-mac", "leaf_tak"} {
		if _, err := dial(t, h.url(), u); err == nil {
			t.Fatalf("%s connected as a STANDARD client", u)
		} else {
			pass(t, "DENY  conn %-20s as STANDARD client: %v", u, err)
		}
	}
	for _, u := range []string{"policy", "entity", "ingest"} {
		if _, err := dial(t, h.wsURL(), u); err == nil {
			t.Fatalf("%s connected over WebSocket", u)
		} else {
			pass(t, "DENY  conn %-20s over WebSocket: %v", u, err)
		}
	}
	for _, u := range []string{"zeus", "agent-fabric"} {
		c, err := dial(t, h.wsURL(), u)
		if err != nil {
			t.Fatalf("%s over WebSocket: %v", u, err)
		}
		pass(t, "ALLOW conn %-20s over WebSocket", u)
		if u == "zeus" {
			expectPub(t, c, "cmd.fleet.drone.sim-x1.goto", false, nil)
			expectPub(t, c, "ent.drone.sim-x1.state", false, nil)
		}
	}
	if _, err := dial(t, h.url(), "", nats.UserInfo("policy", "wrong-password")); err == nil {
		t.Fatal("wrong password accepted")
	} else {
		pass(t, "DENY  conn policy with a wrong password: %v", err)
	}
	// From A2 there is no leafnodes.authorization, so any account user without connection types
	// could open a leaf link. The shared (browser-exposed) onemind login and the SYS users cannot.
	for _, u := range []string{"onemind", "sys", "monitor"} {
		h.startLeaf(t, u)
		time.Sleep(1200 * time.Millisecond)
		if n := h.srv.NumLeafNodes(); n != 0 {
			t.Fatalf("%s logged in as a leaf (%d leafs)", u, n)
		}
		pass(t, "DENY  leaf %-20s cannot log in as a leaf node: %s", u, firstOr(h.log.grep("connection type"), h.log.grep("Authorization")))
	}
	for _, u := range []string{"onemind", "sys"} {
		if _, err := dial(t, h.url(), u); err != nil {
			t.Fatalf("%s as STANDARD client: %v", u, err)
		}
		pass(t, "ALLOW conn %-20s as STANDARD client", u)
	}
	if _, err := dial(t, h.wsURL(), "onemind"); err != nil {
		t.Fatalf("onemind over WebSocket (ws.onemindos.dev): %v", err)
	}
	pass(t, "ALLOW conn onemind              over WebSocket (ws.onemindos.dev clients keep working)")
}

// Stage A1 is additive: shared users keep working unchanged (so nothing breaks), new users are
// restricted from the first minute, and new leaf users cannot log in until A2.
func TestStageA1Additive(t *testing.T) {
	h := startHub(t, "A1")
	onemind := h.conn(t, "onemind")
	bridge := h.conn(t, "bridge-mavlink-home")
	expectPub(t, onemind, "cmd.fleet.drone.sim-x1.goto", true, bridge)
	t.Log("NOTE stage A1 still lets the shared onemind user publish commands; stage B closes it")
	zeus := h.conn(t, "zeus")
	expectPub(t, zeus, "cmd.fleet.drone.sim-x1.goto", false, bridge)
	ingest := h.conn(t, "ingest")
	expectPub(t, ingest, "cmd.fleet.drone.sim-x1.goto", false, bridge)
	h.startLeaf(t, "leaf_tak")
	if !waitLeafs(t, h.srv, 1, 5*time.Second) {
		t.Fatal("leaf_tak (leafnodes.authorization) did not connect at A1")
	}
	pass(t, "A1: leaf_tak still connects through leafnodes.authorization")
	h.startLeaf(t, "leaf-tak")
	time.Sleep(1500 * time.Millisecond)
	if h.srv.NumLeafNodes() != 1 {
		t.Fatal("leaf-tak connected at A1, but leafnodes.authorization should take precedence")
	}
	pass(t, "A1: new leaf user leaf-tak is refused while leafnodes.authorization exists (A2 needed): %v", firstOr(h.log.grep("leaf-tak"), h.log.grep("Authorization")))
}

func firstOr(a, b []string) string {
	if len(a) > 0 {
		return a[0]
	}
	if len(b) > 0 {
		return b[0]
	}
	return "(no log line)"
}

// The activation path of IDENTITIES.md on one server: live-like -> A1 by reload, A1 -> A2
// needs a restart (leafnode users are not reloadable), A2 -> B by reload (applies to open
// connections), and the rollbacks.
func TestActivationAndRollback(t *testing.T) {
	live := liveLike(t)
	h := newHub(t, "live", live)
	h.startLeaf(t, "leaf_tak")
	if !waitLeafs(t, h.srv, 1, 5*time.Second) {
		t.Fatal("leaf_tak did not connect to the live-like config")
	}
	onemind := h.conn(t, "onemind")
	pass(t, "live-like config: onemind and the leaf_tak leaf connected")

	// live -> A1 (ConfigMap change; the reloader sends a reload).
	if err := h.reload(renderStage(t, "A1")); err != nil {
		t.Fatalf("reload live -> A1: %v", err)
	}
	pass(t, "reload live -> A1 accepted")
	if h.srv.NumLeafNodes() != 1 || !onemind.nc.IsConnected() {
		t.Fatal("A1 reload dropped existing connections")
	}
	pass(t, "A1 reload kept the leaf_tak leaf and the onemind connection")
	if _, err := dial(t, h.url(), "policy"); err != nil {
		t.Fatalf("policy cannot log in after A1: %v", err)
	}
	pass(t, "A1: new user policy logs in")

	// A1 -> A2 cannot be reloaded.
	err := h.reload(renderStage(t, "A2"))
	if err == nil || !strings.Contains(err.Error(), "not supported") {
		t.Fatalf("reload A1 -> A2: want a 'not supported' error, got %v", err)
	}
	pass(t, "reload A1 -> A2 refused as expected: %v", err)
	if !onemind.nc.IsConnected() {
		t.Fatal("refused reload dropped clients")
	}
	pass(t, "refused reload left the server running on the A1 config")

	// Restart with A2 (one pod at a time in the cluster).
	h.stop()
	h.writeConf(renderStage(t, "A2"))
	h.start()
	if !waitLeafs(t, h.srv, 1, 10*time.Second) {
		t.Fatal("leaf_tak did not reconnect after the A2 restart")
	}
	pass(t, "A2 restart: leaf_tak leaf reconnected as an ONEMIND account user")
	h.startLeaf(t, "leaf-tak")
	if !waitLeafs(t, h.srv, 2, 5*time.Second) {
		t.Fatal("leaf-tak did not connect at A2")
	}
	pass(t, "A2: new leaf user leaf-tak connects")
	om := h.conn(t, "onemind")
	bridge := h.conn(t, "bridge-mavlink-home")
	expectPub(t, om, "cmd.fleet.drone.sim-x1.goto", true, bridge)

	// A2 -> B by reload, enforced on the already open onemind connection.
	if err := h.reload(renderStage(t, "B")); err != nil {
		t.Fatalf("reload A2 -> B: %v", err)
	}
	pass(t, "reload A2 -> B accepted")
	time.Sleep(200 * time.Millisecond)
	expectPub(t, om, "cmd.fleet.drone.sim-x1.goto", false, bridge)
	pass(t, "B applied to the open onemind connection without a reconnect")

	// Rollback B -> A2 (reload) and A2 -> live (restart).
	if err := h.reload(renderStage(t, "A2")); err != nil {
		t.Fatalf("rollback reload B -> A2: %v", err)
	}
	time.Sleep(200 * time.Millisecond)
	om2 := h.conn(t, "onemind")
	expectPub(t, om2, "cmd.fleet.drone.sim-x1.goto", true, bridge)
	pass(t, "rollback B -> A2 by reload restores the shared user")
	if err := h.reload(live); err == nil {
		t.Fatal("rollback A2 -> live by reload unexpectedly accepted")
	}
	h.stop()
	h.writeConf(live)
	h.start()
	if !waitLeafs(t, h.srv, 1, 10*time.Second) {
		t.Fatal("leaf_tak did not reconnect after rollback to the live config")
	}
	pass(t, "rollback A2 -> live needs a restart; after it the leaf_tak leaf reconnects")
	if _, err := dial(t, h.url(), "policy"); err == nil {
		t.Fatal("policy still logs in after rollback")
	}
	pass(t, "after rollback the new users are gone (policy refused)")
	_ = errors.New
}
