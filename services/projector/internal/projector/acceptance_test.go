package projector

// Acceptance tests of onemind-xr docs/contracts/cot-projection.md §12 (tests 2 to 6 and the
// projector half of 8; test 1 is TestTypeTable) and JADC2-MAP §5.3 stop-line tests 8 and 10,
// on a disposable nats-server with fixture entities.

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"testing"
	"time"
)

const (
	subjInject = "tak.cot.inject"
	subjMeshTx = "tak.cot.mesh.tx"
)

func uids(ms []cotMsg, removal bool) map[string]int {
	out := map[string]int{}
	for _, m := range ms {
		if m.removal == removal {
			out[m.uid]++
		}
	}
	return out
}

// Test 2 (A6): 2,000 aircraft, one AOI around 20 of them: exactly 20 OM.aircraft.* uids on
// tak.cot.inject. Deleting the AOI gives 20 removals within AOI_EXIT_GRACE_S + 5 s, then
// nothing more. Also: a restart with the AOI in place projects the same 20 from ENT.
func TestAcceptance2AOI(t *testing.T) {
	b := setupBus(t)
	c := collect(t, b.nc, subjInject)
	cfg := testCfg(TargetServer)
	cfg.AOIExitGrace = 3 * time.Second
	p := startProjector(t, b, cfg)

	// 20 aircraft inside a 0.1 degree box at 30N 97W, 1,980 spread over 20..40N, 80..120W
	// but outside the box.
	inside := map[string]bool{}
	for i := 0; i < 2000; i++ {
		id := idN("icao", i)
		lat, lon := 20+float64(i%200)*0.1, -120+float64(i/200)*4.0
		if i < 20 {
			lat, lon = 30.01+float64(i)*0.003, -97.04+float64(i)*0.003
			inside["OM.aircraft."+id] = true
		} else if lat > 29.8 && lat < 30.2 && lon > -97.2 && lon < -96.8 {
			lat += 1
		}
		b.putEnt(t, fixture("aircraft", id, lat, lon, withKind("fixed_wing")))
	}
	time.Sleep(1500 * time.Millisecond)
	if n := len(c.all()); n != 0 {
		t.Fatalf("before any AOI: %d messages on %s, want 0", n, subjInject)
	}
	if _, err := b.geof.Put(context.Background(), "aoi-test", geofence("aoi-test", "aoi", boxRing(30.05, -97, 0.1), true)); err != nil {
		t.Fatal(err)
	}
	took := waitFor(t, 10*time.Second, "20 uids", func() bool { return len(uids(c.all(), false)) >= 20 })
	time.Sleep(2 * time.Second) // nothing else may arrive
	got := uids(c.all(), false)
	for u := range got {
		if !inside[u] {
			t.Errorf("projected outside the AOI: %s", u)
		}
	}
	if len(got) != 20 {
		t.Fatalf("got %d distinct uids, want exactly 20", len(got))
	}
	for _, m := range c.all() {
		if m.hdr.Get("OneMind-Projector") != "hub" || m.hdr.Get("OneMind-Marking") != "CLEAR" || !strings.HasPrefix(m.hdr.Get("Nats-Msg-Id"), m.uid+":") {
			t.Fatalf("headers: %v", m.hdr)
		}
	}
	evidence(t, "2,000 aircraft, AOI around 20: %d distinct OM.aircraft uids on %s (all inside), first 20 within %s", len(got), subjInject, took.Round(time.Millisecond))

	// Restart: ENT last-per-subject plus the KV watch project the same 20 at once.
	p.stop()
	c.reset()
	p = startProjector(t, b, cfg)
	waitFor(t, 10*time.Second, "20 uids after restart", func() bool { return len(uids(c.all(), false)) >= 20 })
	time.Sleep(time.Second)
	if n := len(uids(c.all(), false)); n != 20 {
		t.Fatalf("after restart: %d uids, want 20", n)
	}
	evidence(t, "restart with AOI in place: 20 uids projected again from ENT")

	// Delete the AOI: 20 removals within grace + 5 s, then quiet.
	c.reset()
	deleted := time.Now()
	if err := b.geof.Delete(context.Background(), "aoi-test"); err != nil {
		t.Fatal(err)
	}
	limit := cfg.AOIExitGrace + 5*time.Second
	waitFor(t, limit, "20 removals", func() bool { return len(uids(c.all(), true)) >= 20 })
	rem := uids(c.all(), true)
	var lastRem time.Time
	for _, m := range c.all() {
		if m.removal && m.at.After(lastRem) {
			lastRem = m.at
		}
	}
	for u := range rem {
		if !inside[u] {
			t.Errorf("removal for an entity that was never projected: %s", u)
		}
	}
	quietFrom := len(c.all())
	time.Sleep(5 * time.Second)
	if extra := c.all()[quietFrom:]; len(extra) != 0 {
		t.Fatalf("%d messages after the removals (first %s)", len(extra), extra[0].uid)
	}
	evidence(t, "AOI deleted: %d removals, last one %s after delete (limit %s = grace %s + 5 s); 0 messages in the next 5 s",
		len(rem), lastRem.Sub(deleted).Round(time.Millisecond), limit, cfg.AOIExitGrace)
	_ = p
}

// Test 3: AMBER+STRICT releasable to homestead reaches the target only when TAK_RELEASE_GROUPS
// contains homestead; RED never does.
func TestAcceptance3Marking(t *testing.T) {
	b := setupBus(t)
	c := collect(t, b.nc, subjInject)
	b.putEnt(t, fixture("drone", "sim-amber", 30, -97, ownAsset, marking("AMBER+STRICT", "homestead")))
	b.putEnt(t, fixture("drone", "sim-red", 30, -97, ownAsset, marking("RED", "homestead")))
	b.putEnt(t, fixture("drone", "sim-clear", 30, -97, ownAsset))

	p := startProjector(t, b, testCfg(TargetServer))
	waitFor(t, 5*time.Second, "control entity", func() bool { return uids(c.all(), false)["OM.drone.sim-clear"] > 0 })
	time.Sleep(2 * time.Second)
	got := uids(c.all(), false)
	if got["OM.drone.sim-amber"] > 0 || got["OM.drone.sim-red"] > 0 {
		t.Fatalf("without TAK_RELEASE_GROUPS: %v", got)
	}
	evidence(t, "TAK_RELEASE_GROUPS empty: CLEAR sent, AMBER+STRICT[homestead] and RED not sent (%v)", got)
	p.stop()

	c.reset()
	cfg := testCfg(TargetServer)
	cfg.ReleaseGroups = []string{"homestead"}
	startProjector(t, b, cfg)
	waitFor(t, 5*time.Second, "amber entity", func() bool { return uids(c.all(), false)["OM.drone.sim-amber"] > 0 })
	time.Sleep(2 * time.Second)
	got = uids(c.all(), false)
	if got["OM.drone.sim-red"] > 0 {
		t.Fatalf("RED projected: %v", got)
	}
	for _, m := range c.all() {
		if m.uid == "OM.drone.sim-amber" && (m.hdr.Get("OneMind-Marking") != "AMBER+STRICT" || m.hdr.Get("OneMind-Releasable-To") != "homestead") {
			t.Fatalf("marking headers: %v", m.hdr)
		}
	}
	// Marking tightened to RED: removal at once.
	c.reset()
	b.putEnt(t, fixture("drone", "sim-amber", 30, -97, ownAsset, marking("RED", "homestead")))
	waitFor(t, 3*time.Second, "removal after RED", func() bool { return uids(c.all(), true)["OM.drone.sim-amber"] > 0 })
	evidence(t, "TAK_RELEASE_GROUPS=homestead: AMBER+STRICT[homestead] sent with marking headers, RED never; marking tightened to RED -> removal")
}

// Test 4 and stop-line test 10 (SL3): a person-class track inside an AOI and a person entity
// never produce CoT on either target, even with MESH_WORLD_FEEDS on the mesh projector.
func TestAcceptance4PeopleNeverProjected(t *testing.T) {
	b := setupBus(t)
	c := collect(t, b.nc, subjInject, subjMeshTx)
	if _, err := b.geof.Put(context.Background(), "aoi-gate", geofence("aoi-gate", "aoi", boxRing(30, -97, 0.05), true)); err != nil {
		t.Fatal(err)
	}
	startProjector(t, b, testCfg(TargetServer))
	mesh := testCfg(TargetMesh)
	mesh.MeshWorldFeeds = true
	mesh.OtherInterval = time.Second
	startProjector(t, b, mesh)

	people := []map[string]any{
		// A person-class track as the A4 detection adapter writes it (adapt/intel.go): kind person.
		fixture("track", "det-person-1", 30.001, -97.001, withKind("person"), marking("AMBER+STRICT"), ind("egress_tak", false)),
		// The same with every other gate opened: still never projected.
		fixture("track", "det-person-2", 30.002, -97.002, withKind("person"), ind("starred", true),
			set("current_task", map[string]any{"task_id": "tsk-x", "state": "executing"})),
		// A person entity (TAK EUD PLI), even with egress forced on and CLEAR.
		fixture("person", "tak-android-1", 30.003, -97.003, set("template", "asset")),
		fixture("station", "rf-1", 30.004, -97.004),
	}
	for _, m := range people {
		b.putEnt(t, m)
	}
	b.putEnt(t, fixture("aircraft", "icao-control", 30.005, -97.005))
	waitFor(t, 5*time.Second, "control on both targets", func() bool {
		return len(c.on(subjInject)) > 0 && len(c.on(subjMeshTx)) > 0
	})
	for i := 0; i < 20; i++ { // keep updating for 4 s
		for _, m := range people {
			m["seq"] = i + 2
			b.putEnt(t, m)
		}
		time.Sleep(200 * time.Millisecond)
	}
	for _, m := range c.all() {
		if !strings.HasPrefix(m.uid, "OM.aircraft.") {
			t.Fatalf("person-class or station entity projected on %s: %s", m.subject, m.uid)
		}
	}
	evidence(t, "person-class tracks, person and station entities inside the AOI over 4 s of updates: 0 CoT on %s and %s; control aircraft sent %d/%d",
		subjInject, subjMeshTx, len(c.on(subjInject)), len(c.on(subjMeshTx)))
}

// G6: indicators.egress_tak false stops an entity inside an AOI; turning it off removes one
// that was on TAK.
func TestAcceptanceEgressTAK(t *testing.T) {
	b := setupBus(t)
	c := collect(t, b.nc, subjInject)
	if _, err := b.geof.Put(context.Background(), "aoi-e", geofence("aoi-e", "keep_out", boxRing(30, -97, 0.05), true)); err != nil {
		t.Fatal(err)
	}
	startProjector(t, b, testCfg(TargetServer))
	b.putEnt(t, fixture("vessel", "mmsi-off", 30, -97, ind("egress_tak", false)))
	b.putEnt(t, fixture("vessel", "mmsi-on", 30, -97))
	waitFor(t, 5*time.Second, "egress on", func() bool { return uids(c.all(), false)["OM.vessel.mmsi-on"] > 0 })
	time.Sleep(time.Second)
	if uids(c.all(), false)["OM.vessel.mmsi-off"] > 0 {
		t.Fatal("egress_tak false projected")
	}
	b.putEnt(t, fixture("vessel", "mmsi-on", 30, -97, ind("egress_tak", false)))
	waitFor(t, 3*time.Second, "removal", func() bool { return uids(c.all(), true)["OM.vessel.mmsi-on"] > 0 })
	evidence(t, "egress_tak=false inside an AOI: not sent; egress_tak turned off: removal")
}

// Tasked and targeted entities (KV mission-state) are projected outside any AOI; the AOI of an
// active mission acts as an AOI even when its zone kind does not.
func TestAcceptanceMissionTier(t *testing.T) {
	b := setupBus(t)
	c := collect(t, b.nc, subjInject)
	startProjector(t, b, testCfg(TargetServer))
	b.putEnt(t, fixture("vessel", "mmsi-target", 10, 10))
	b.putEnt(t, fixture("aircraft", "icao-in-mission", 40, 40))
	time.Sleep(time.Second)
	if len(c.all()) != 0 {
		t.Fatal("projected before any mission")
	}
	ctx := context.Background()
	_, _ = b.missions.Put(ctx, "tsk-01", mustJSON(map[string]any{"task_id": "tsk-01", "rollup": "running",
		"assignee": map[string]any{"entity_type": "drone", "entity_id": "sim-x1"}, "target_snapshots": []any{map[string]any{"entity_id": "mmsi-target"}}}))
	waitFor(t, 5*time.Second, "targeted vessel", func() bool { return uids(c.all(), false)["OM.vessel.mmsi-target"] > 0 })
	_, _ = b.geof.Put(ctx, "zone-fp", geofence("zone-fp", "footprint", boxRing(40, 40, 0.1), true))
	time.Sleep(time.Second)
	if uids(c.all(), false)["OM.aircraft.icao-in-mission"] > 0 {
		t.Fatal("footprint zone acted as an AOI without a mission")
	}
	_, _ = b.missions.Put(ctx, "msn-01", mustJSON(map[string]any{"mission_id": "msn-01", "state": "active", "aoi": "zone-fp"}))
	waitFor(t, 5*time.Second, "mission AOI", func() bool { return uids(c.all(), false)["OM.aircraft.icao-in-mission"] > 0 })
	evidence(t, "task target projected outside any AOI; active mission's aoi zone projects the aircraft inside it")
}

// maxInWindow is the largest number of messages whose send times fall in any window of w.
func maxInWindow(ms []cotMsg, w time.Duration) int {
	ms = append([]cotMsg(nil), ms...)
	sort.Slice(ms, func(i, j int) bool { return ms[i].sent.Before(ms[j].sent) })
	best := 0
	for i := range ms {
		j := i
		for j < len(ms) && ms[j].sent.Sub(ms[i].sent) < w {
			j++
		}
		if j-i > best {
			best = j - i
		}
	}
	return best
}

func minGap(ms []cotMsg) time.Duration {
	g := time.Hour
	for i := 1; i < len(ms); i++ {
		if d := ms[i].sent.Sub(ms[i-1].sent); d < g {
			g = d
		}
	}
	return g
}

// Test 6: a sim-* drone at 10 Hz yields at most 1 msg/s on tak.cot.inject and at most 1 per
// 5 s on tak.cot.mesh.tx; the global ceiling (50 msg/s) holds with 5,000 entities in AOIs.
func TestAcceptance6Rate(t *testing.T) {
	b := setupBus(t)
	c := collect(t, b.nc, subjInject, subjMeshTx)
	startProjector(t, b, testCfg(TargetServer))
	startProjector(t, b, testCfg(TargetMesh))
	start := time.Now()
	for i := 0; time.Since(start) < 12*time.Second; i++ {
		b.putEnt(t, fixture("drone", "sim-x1", 30+float64(i)*1e-5, -97, ownAsset, withKind("multirotor"), set("seq", i+1)))
		time.Sleep(100 * time.Millisecond)
	}
	time.Sleep(200 * time.Millisecond)
	inj, mesh := c.on(subjInject), c.on(subjMeshTx)
	if w := maxInWindow(inj, time.Second); w > 1 || len(inj) < 10 {
		t.Fatalf("inject: %d msgs, max %d per second", len(inj), w)
	}
	if w := maxInWindow(mesh, 5*time.Second); w > 1 || len(mesh) < 2 {
		t.Fatalf("mesh: %d msgs, max %d per 5 s", len(mesh), w)
	}
	for _, m := range mesh {
		if strings.Contains(m.xml, "<onemind") || len(m.xml) > 1200 {
			t.Fatalf("mesh profile: %s", m.xml)
		}
	}
	evidence(t, "sim-x1 at 10 Hz for 12 s (~120 updates): %d on %s (min gap %s, max 1 per s), %d on %s (min gap %s, max 1 per 5 s)",
		len(inj), subjInject, minGap(inj).Round(time.Millisecond), len(mesh), subjMeshTx, minGap(mesh).Round(time.Millisecond))

	// Ceiling: 5,000 aircraft inside AOIs, default cap 500 and ceiling 50 msg/s.
	c.reset()
	ctx := context.Background()
	_, _ = b.geof.Put(ctx, "aoi-big", geofence("aoi-big", "aoi", boxRing(0, 0, 5), true))
	for i := 0; i < 5000; i++ {
		b.putEnt(t, fixture("aircraft", idN("bulk", i), -4+float64(i%100)*0.08, -4+float64(i/100)*0.16))
	}
	t0 := time.Now()
	waitFor(t, 40*time.Second, "500 live and the queue drained", func() bool {
		l, _ := liveSet(c.on(subjInject)) // sim-x1 from the first half is still live and counts
		return len(l) == 500 && time.Since(t0) > 12*time.Second
	})
	time.Sleep(2 * time.Second)
	inj = c.on(subjInject)
	w := maxInWindow(inj, time.Second)
	if w > 50 {
		t.Fatalf("ceiling broken: %d msgs in one second", w)
	}
	live, swapped := liveSet(inj)
	if len(live) != 500 || !live["OM.drone.sim-x1"] {
		t.Fatalf("cap: %d live uids (sim-x1 live: %v), want 500 (MAX_PROJECTED) with the own asset kept first", len(live), live["OM.drone.sim-x1"])
	}
	evidence(t, "5,000 aircraft in an AOI: %d msgs on %s in %s, max %d in any one-second window by CoT send time (ceiling 50); %d live uids incl. sim-x1 (cap 500), %d replaced by nearer ones with a removal (cap order, section 3.3)",
		len(inj), subjInject, time.Since(t0).Round(time.Second), w, len(live), swapped)
}

// liveSet returns the uids whose last message is not a removal, and how many were removed.
func liveSet(ms []cotMsg) (map[string]bool, int) {
	last := map[string]bool{}
	for _, m := range ms {
		last[m.uid] = !m.removal
	}
	live, removed := map[string]bool{}, 0
	for u, on := range last {
		if on {
			live[u] = true
		} else {
			removed++
		}
	}
	return live, removed
}

// cotXML is a minimal CoT event for the reverse-direction tests.
func cotXML(uid, typ string, extra string) []byte {
	now := time.Now().UTC()
	return []byte(fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?><event version="2.0" uid="%s" type="%s" how="h-e" time="%s" start="%s" stale="%s"><point lat="30.0001" lon="-97.0001" hae="10" ce="5" le="5"/><detail>%s</detail></event>`,
		uid, typ, now.Format(cotTimeFmt), now.Format(cotTimeFmt), now.Add(5*time.Minute).Format(cotTimeFmt), extra))
}

// Test 5 (echo), stop-line test 8 (SL1) and the projector half of test 8, with the A4 entity
// service from this branch's base reading COT (tak.cot.>).
func TestAcceptance5EchoAndTasking(t *testing.T) {
	b := setupBus(t)
	entLog := startEntityService(t, b)
	c := collect(t, b.nc, subjInject)
	p := startProjector(t, b, testCfg(TargetServer))
	ctx := context.Background()
	cmdSub, err := b.nc.SubscribeSync("cmd.>")
	if err != nil {
		t.Fatal(err)
	}
	polSub, err := b.nc.SubscribeSync("svc.policy.>")
	if err != nil {
		t.Fatal(err)
	}
	_ = b.nc.Flush()

	// Wait for the entity service to hold its lease and follow COT.
	time.Sleep(2 * time.Second)

	// --- Test 5: echo. The projector sends sim-x1; its own output comes back on both TAK input
	// subjects. No new entity, no second CoT.
	b.putEnt(t, fixture("drone", "sim-x1", 30, -97, ownAsset, withKind("multirotor")))
	waitFor(t, 5*time.Second, "sim-x1 projected", func() bool { return len(c.all()) > 0 })
	out := c.all()[0]
	before := b.entSubjects(t)
	for _, s := range []string{"tak.cot.out", "tak.cot.mesh.rx"} {
		if err := b.nc.Publish(s, []byte(out.xml)); err != nil {
			t.Fatal(err)
		}
	}
	// A copy with a foreign uid but the projector's flow tag is dropped as well.
	tagged := strings.Replace(out.xml, `uid="OM.drone.sim-x1"`, `uid="relay-copy-1"`, 1)
	_ = b.nc.Publish("tak.cot.out", []byte(tagged))
	time.Sleep(3 * time.Second)
	if after := b.entSubjects(t); after != before {
		t.Fatalf("echo created entities: ENT subjects %d -> %d", before, after)
	}
	for _, m := range c.all() {
		if m.uid != "OM.drone.sim-x1" {
			t.Fatalf("second CoT from the echo: %s", m.uid)
		}
	}
	evidence(t, "echo: projector output fed back on tak.cot.out and tak.cot.mesh.rx (and a relayed copy with the flow tag): ENT subjects %d -> %d, only OM.drone.sim-x1 on %s", before, b.entSubjects(t), subjInject)

	// --- Stop-line test 8: an incoming strike / engage CoT creates no task.
	before = b.entSubjects(t)
	_ = b.nc.Publish("tak.cot.out", cotXML("ATAK-strike-1", "t-k-d", `<contact callsign="EUD1"/><link uid="OM.drone.sim-x1" relation="p-p"/>`))
	_ = b.nc.Publish("tak.cot.mesh.rx", cotXML("ATAK-engage-2", "t-e-x", `<contact callsign="EUD2"/>`))
	waitFor(t, 3*time.Second, "refusal logged", func() bool { return p.svc.Engine.C.Get("tak_tasking_refused") >= 2 })
	time.Sleep(2 * time.Second)
	if after := b.entSubjects(t); after != before {
		t.Fatalf("tasking CoT created entities: %d -> %d", before, after)
	}
	keys, _ := b.missions.ListKeys(ctx)
	n := 0
	for range keys.Keys() {
		n++
	}
	if n != 0 {
		t.Fatalf("mission-state has %d keys after tasking CoT", n)
	}
	if m, err := cmdSub.NextMsg(100 * time.Millisecond); err == nil {
		t.Fatalf("command published: %s", m.Subject)
	}
	if m, err := polSub.NextMsg(100 * time.Millisecond); err == nil {
		t.Fatalf("policy request published: %s", m.Subject)
	}
	for _, m := range c.all() {
		if strings.HasPrefix(m.typ, "t-") {
			t.Fatalf("tasking type sent to TAK: %s", m.typ)
		}
	}
	logs := p.logs.String()
	if !strings.Contains(logs, `"cot_type":"t-k-d"`) || !strings.Contains(logs, `"cot_type":"t-e-x"`) || !strings.Contains(logs, "refused TAK tasking") {
		t.Fatalf("refusal not logged:\n%s", logs)
	}
	evidence(t, "stop-line 8: t-k-d on tak.cot.out and t-e-x on tak.cot.mesh.rx -> 0 ENT entities, 0 mission-state keys, 0 cmd.> and 0 svc.policy.> messages; projector logged %d refusals", p.svc.Engine.C.Get("tak_tasking_refused"))

	// --- Test 8, projector half (MESH_TO_TAK_SERVER=false is the gateway default): an EUD PLI on
	// the mesh becomes ent.person.* with AMBER+STRICT and egress_tak false, and never reaches
	// tak.cot.inject through the projector.
	_ = b.nc.Publish("tak.cot.mesh.rx", cotXML("ANDROID-a6test01", "a-f-G-U-C",
		`<contact callsign="EUD-A6" endpoint="192.168.1.20:4242:tcp"/><takv device="Pixel" platform="ATAK-CIV" os="34" version="5.5"/><__group name="Cyan" role="Team Member"/>`))
	var person map[string]any
	waitFor(t, 10*time.Second, "ent.person from the mesh PLI", func() bool {
		person = readEnt(t, b, "ent.person.tak-android-a6test01.state")
		return person != nil
	})
	if jsonGet(person, "marking", "tlp") != "AMBER+STRICT" || jsonGet(person, "indicators", "egress_tak") != false {
		t.Fatalf("mesh person record: marking=%v egress=%v", jsonGet(person, "marking"), jsonGet(person, "indicators", "egress_tak"))
	}
	time.Sleep(2 * time.Second)
	for _, m := range c.all() {
		if m.uid != "OM.drone.sim-x1" {
			t.Fatalf("mesh PLI reached %s as %s", subjInject, m.uid)
		}
	}
	if p.svc.Engine.C.Get("blocked_"+ReasonPerson) == 0 {
		t.Fatal("person gate not counted")
	}
	evidence(t, "test 8 (projector half): mesh PLI -> ent.person.tak-android-a6test01 (tlp %v, egress_tak %v); 0 messages for it on %s",
		jsonGet(person, "marking", "tlp"), jsonGet(person, "indicators", "egress_tak"), subjInject)
	_ = entLog
}
