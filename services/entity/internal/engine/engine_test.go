package engine

import (
	"bufio"
	"encoding/json"
	"os"
	"strings"
	"testing"

	"github.com/onemindos/omos/services/entity/internal/adapt"
	"github.com/onemindos/omos/services/entity/internal/model"
)

const t0 = int64(1793275200000)

type fixture struct {
	Subject string `json:"subject"`
	Data    string `json:"data"`
}

func loadFixtures(t *testing.T) []fixture {
	t.Helper()
	f, err := os.Open("../../testdata/fixtures/replay.jsonl")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	var out []fixture
	sc := bufio.NewScanner(f)
	sc.Buffer(make([]byte, 1<<20), 1<<20)
	for sc.Scan() {
		var fx fixture
		if err := json.Unmarshal(sc.Bytes(), &fx); err != nil {
			t.Fatal(err)
		}
		out = append(out, fx)
	}
	return out
}

type clock struct{ ms int64 }

func (c *clock) now() int64 { return c.ms }

func newEngine(c *clock) *Engine { return New(DefaultConfig(), c.now) }

func obs(t *testing.T, subject, data string) *model.Observation {
	t.Helper()
	o, skip, err := adapt.Parse(subject, []byte(data))
	if err != nil || o == nil {
		t.Fatalf("parse %s: skip=%q err=%v", subject, skip, err)
	}
	return o
}

// replay applies every fixture and returns the last publish per entity and all results.
func replay(t *testing.T, e *Engine) (map[string]*Publish, []Result) {
	t.Helper()
	last := map[string]*Publish{}
	var all []Result
	for i, fx := range loadFixtures(t) {
		o, _, err := adapt.Parse(fx.Subject, []byte(fx.Data))
		if err != nil || o == nil {
			continue
		}
		o.Stream, o.StreamSeq = "TEST", uint64(i+1)
		r := e.Apply(o)
		all = append(all, r)
		if r.Publish != nil {
			last[r.Publish.Entity.EntityID] = r.Publish
		}
	}
	return last, all
}

func TestSameAircraftTwoFeedsIsOneEntityWithTwoSources(t *testing.T) {
	e := newEngine(&clock{t0 + 60000})
	last, _ := replay(t, e)
	p := last["icao-a1b2c3"]
	if p == nil {
		t.Fatal("no entity icao-a1b2c3")
	}
	if p.Subject != "ent.aircraft.icao-a1b2c3.state" {
		t.Fatalf("subject %s", p.Subject)
	}
	ints := map[string]bool{}
	for _, s := range p.Entity.Sources {
		ints[s.Integration] = true
	}
	if len(p.Entity.Sources) != 2 || !ints["feed-aircraft"] || !ints["world-air"] {
		t.Fatalf("want 2 sources feed-aircraft+world-air, got %+v", p.Entity.Sources)
	}
	if p.Entity.Provenance.How != "m-f" {
		t.Fatalf("fused how = %s", p.Entity.Provenance.How)
	}
	for id := range last {
		if strings.HasPrefix(id, "tak-adsb") {
			t.Fatalf("aircraft also published under its TAK uid: %s", id)
		}
	}
	single := last["icao-abc123"]
	if single == nil || p.Entity.Confidence <= single.Entity.Confidence {
		t.Fatalf("corroborated confidence should exceed single-source confidence")
	}
}

func TestShipUpdatesInPlace(t *testing.T) {
	c := &clock{t0 + 300000}
	e := newEngine(c)
	var pubs []*Publish
	for i := 0; i < 3; i++ {
		d := `{"uid":"mmsi-367000001","type":"sea.cargo","lat":` + []string{"40.6000", "40.6050", "40.6100"}[i] +
			`,"lon":-74.04,"speed":5.5,"course":0,"stale":600,"t":` + jsonInt(t0+int64(i)*100000) + `,"meta":{"mmsi":"367000001"}}`
		r := e.Apply(obs(t, "evt.world.sea.mmsi-367000001", d))
		if r.Publish == nil {
			t.Fatalf("update %d not published (moved 550 m)", i)
		}
		pubs = append(pubs, r.Publish)
	}
	for i, p := range pubs {
		if p.Entity.EntityID != "mmsi-367000001" || p.Entity.Seq != uint64(i+1) || p.MsgID != "mmsi-367000001:"+jsonInt(int64(i+1)) {
			t.Fatalf("publish %d: id=%s seq=%d msgid=%s", i, p.Entity.EntityID, p.Entity.Seq, p.MsgID)
		}
	}
	if pubs[2].Entity.Lat <= pubs[0].Entity.Lat {
		t.Fatalf("position did not move north: %v -> %v", pubs[0].Entity.Lat, pubs[2].Entity.Lat)
	}
	if pubs[2].TTL == "" || pubs[2].TTL == "never" {
		t.Fatalf("vessel TTL %q", pubs[2].TTL)
	}
}

func jsonInt(v int64) string { b, _ := json.Marshal(v); return string(b) }

func TestOlderSourceTimeIsDroppedAndCounted(t *testing.T) {
	e := newEngine(&clock{t0 + 60000})
	d := func(ts int64) string {
		return `{"uid":"icao-a1b2c3","lat":40.6,"lon":-73.7,"t":` + jsonInt(ts) + `,"meta":{"icao":"a1b2c3"}}`
	}
	if r := e.Apply(obs(t, "evt.world.air.icao-a1b2c3", d(t0+10000))); r.Publish == nil {
		t.Fatal("first not published")
	}
	for _, ts := range []int64{t0 + 5000, t0 + 10000} { // older, then a duplicate
		if r := e.Apply(obs(t, "evt.world.air.icao-a1b2c3", d(ts))); r.Drop != DropOlderSourceT {
			t.Fatalf("t=%d: want drop, got %+v", ts, r)
		}
	}
	if n := e.Counters()[DropOlderSourceT]; n != 2 {
		t.Fatalf("counter = %d", n)
	}
}

func TestFutureSourceTimeIsReplacedByReceiveTime(t *testing.T) {
	c := &clock{t0}
	e := newEngine(c)
	r := e.Apply(obs(t, "evt.world.air.icao-a1b2c3", `{"uid":"icao-a1b2c3","lat":40.6,"lon":-73.7,"t":`+jsonInt(t0+60000)+`}`))
	if r.Publish == nil || r.Publish.Entity.T != t0 || e.Counters()[CountClockSkew] != 1 {
		t.Fatalf("skew not replaced: %+v", r)
	}
}

func TestPersonHostileIsRejected(t *testing.T) {
	e := newEngine(&clock{t0 + 60000})
	_, all := replay(t, e)
	rejected := 0
	for _, r := range all {
		if r.Reject != nil {
			rejected++
			if r.Reject.Reason != RejectPersonHostile || r.Reject.Obs.SourceID != "ANDROID-badf00d" {
				t.Fatalf("unexpected reject %+v", r.Reject)
			}
		}
		if r.Publish != nil && r.Publish.Entity.Disposition == "hostile" {
			t.Fatal("published a hostile disposition")
		}
	}
	if rejected != 1 {
		t.Fatalf("rejected = %d", rejected)
	}
	if _, ok := e.Get(identityTak("ANDROID-badf00d")); ok {
		t.Fatal("hostile person stored")
	}
}

func identityTak(uid string) string { return "tak-" + strings.ToLower(uid) }

func TestHostileAssetBecomesSuspiciousWithConflict(t *testing.T) {
	e := newEngine(&clock{t0 + 60000})
	x := `<event version="2.0" uid="TRK-9" type="a-h-A-M-F" how="h-e" time="2026-10-29T12:00:00Z" start="2026-10-29T12:00:00Z" stale="2026-10-29T12:05:00Z"><point lat="30" lon="-97" hae="500" ce="20" le="20"/><detail><contact callsign="UNK"/></detail></event>`
	o := obs(t, "tak.cot.out", x)
	o.SourceT = t0
	r := e.Apply(o)
	if r.Publish == nil || r.Publish.Entity.Disposition != "suspicious" || len(r.Publish.Entity.Conflicts) == 0 {
		t.Fatalf("got %+v", r)
	}
}

func TestCameraPersonDetectionIsATrackNeverAPerson(t *testing.T) {
	e := newEngine(&clock{t0 + 5000})
	var p *Publish
	for i := 0; i < 3; i++ {
		d := `{"lat":30.27,"lon":-97.74,"t":` + jsonInt(t0+int64(i)*1000) + `,"class":"person","confidence":0.8,"sensor_id":"cam0"}`
		if r := e.Apply(obs(t, "det.land.camera.7", d)); r.Publish != nil {
			p = r.Publish
		}
	}
	if p == nil {
		t.Fatal("confirmed track not published after 3 hits")
	}
	en := p.Entity
	if en.EntityType != "track" || en.Kind != "person" || en.Marking.TLP != "AMBER+STRICT" || en.Indicators.EgressTAK {
		t.Fatalf("person track: type=%s kind=%s tlp=%s egress=%v", en.EntityType, en.Kind, en.Marking.TLP, en.Indicators.EgressTAK)
	}
	if strings.HasPrefix(p.Subject, "ent.person.") {
		t.Fatal("detection produced ent.person")
	}
	if *en.ExpiryT-en.T > 60000 {
		t.Fatalf("person track TTL > 60 s: %d", *en.ExpiryT-en.T)
	}
}

func TestTentativeDetectionsAreNotPublished(t *testing.T) {
	e := newEngine(&clock{t0 + 60000})
	d := `{"lat":30.27,"lon":-97.74,"t":` + jsonInt(t0) + `,"class":"vehicle"}`
	if r := e.Apply(obs(t, "det.land.camera.8", d)); r.Drop != DropTentative {
		t.Fatalf("want tentative, got %+v", r)
	}
}

func TestChangeOnlyPublishAndHeartbeat(t *testing.T) {
	c := &clock{t0}
	e := newEngine(c)
	d := func(ts int64, lon string) string {
		return `{"uid":"sim-x1","lat":30.2672,"lon":` + lon + `,"hae":150,"speed":0,"course":0,"stale":60,"t":` + jsonInt(ts) + `,"meta":{"battery_pct":80}}`
	}
	if r := e.Apply(obs(t, "tel.fleet.drone.sim-x1.position", d(t0, "-97.7431"))); r.Publish == nil {
		t.Fatal("first not published")
	}
	c.ms = t0 + 1000
	if r := e.Apply(obs(t, "tel.fleet.drone.sim-x1.position", d(t0+1000, "-97.74310001"))); r.Publish != nil {
		t.Fatal("sub-metre move published")
	}
	c.ms = t0 + 31000 // heartbeat = min(stale/2, 120 s) = 30 s
	c2 := e.Sweep()
	if len(c2) != 1 || c2[0].Entity.Seq != 2 {
		t.Fatalf("heartbeat: %+v", c2)
	}
}

func TestLifecycleStaleLostExpired(t *testing.T) {
	c := &clock{t0}
	e := newEngine(c)
	e.Apply(obs(t, "evt.world.air.icao-a1b2c3", `{"uid":"icao-a1b2c3","lat":40.6,"lon":-73.7,"stale":30,"t":`+jsonInt(t0)+`}`))
	steps := []struct {
		at        int64
		lifecycle string
		live      bool
	}{{t0 + 31000, "stale", true}, {t0 + 91000, "lost", true}, {t0 + 120000, "expired", false}}
	for _, s := range steps {
		c.ms = s.at
		ps := e.Sweep()
		if len(ps) != 1 || ps[0].Entity.Lifecycle != s.lifecycle || ps[0].Entity.Live != s.live {
			t.Fatalf("at +%ds: %+v", (s.at-t0)/1000, ps)
		}
		if !s.live && ps[0].TTL != "60s" {
			t.Fatalf("final record TTL %s", ps[0].TTL)
		}
	}
	if _, ok := e.Get("icao-a1b2c3"); ok {
		t.Fatal("expired entity still held")
	}
}

func TestDuplicateMMSISplitsInsteadOfAveraging(t *testing.T) {
	e := newEngine(&clock{t0 + 60000})
	last, _ := replay(t, e)
	a, b := last["mmsi-367000001"], last["mmsi-367000001-x2"]
	if a == nil || b == nil {
		t.Fatalf("want original and split entity, got %v %v", a != nil, b != nil)
	}
	if a.Entity.Lat > 41 || b.Entity.Lat < 42 {
		t.Fatalf("positions averaged: %v %v", a.Entity.Lat, b.Entity.Lat)
	}
}

func TestSimDroneMergesTelemetryAndMeshCoT(t *testing.T) {
	e := newEngine(&clock{t0 + 60000})
	last, _ := replay(t, e)
	p := last["sim-x1"]
	if p == nil || len(p.Entity.Sources) != 2 || !p.Entity.Indicators.Simulated || !p.Entity.Indicators.EgressTAK || p.Entity.Provenance.How != "m-s" {
		t.Fatalf("sim-x1: %+v", p)
	}
	if p.Entity.EntityType != "drone" || p.Entity.Health["battery_pct"] == nil {
		t.Fatalf("sim-x1 type/health: %s %v", p.Entity.EntityType, p.Entity.Health)
	}
}

func TestRestoreContinuesSeqAndDropsRedelivery(t *testing.T) {
	c := &clock{t0 + 60000}
	e1 := newEngine(c)
	last, _ := replay(t, e1)
	e2 := newEngine(c)
	for _, p := range last {
		b, _ := json.Marshal(p.Entity)
		var r model.Entity
		_ = json.Unmarshal(b, &r)
		e2.Restore(r)
	}
	_, all := replay(t, e2) // the whole input redelivered
	for _, r := range all {
		if r.Publish != nil {
			t.Fatalf("redelivery republished %s seq %d", r.Publish.Entity.EntityID, r.Publish.Entity.Seq)
		}
	}
	o := obs(t, "evt.world.sea.mmsi-367000001", `{"uid":"mmsi-367000001","lat":40.6012,"lon":-74.04,"stale":600,"t":`+jsonInt(t0+50000)+`}`)
	r := e2.Apply(o)
	if r.Publish == nil || r.Publish.Entity.Seq != last["mmsi-367000001"].Entity.Seq+1 {
		t.Fatalf("seq did not continue: %+v", r.Publish)
	}
}

func TestConfidencePartsAreBoundedAndExplained(t *testing.T) {
	e := newEngine(&clock{t0 + 60000})
	last, _ := replay(t, e)
	for id, p := range last {
		cp := p.Entity.ConfidenceParts
		for _, v := range []float64{p.Entity.Confidence, cp.SourcePrior, cp.TrackQuality, cp.Corroboration} {
			if v < 0 || v > 1 {
				t.Fatalf("%s: out of range %+v", id, cp)
			}
		}
		if p.Entity.Confidence >= 1 {
			t.Fatalf("%s: confidence must stay below 1", id)
		}
	}
}

func TestHumanAssertionBeatsSensorForCallsign(t *testing.T) {
	e := newEngine(&clock{t0 + 60000})
	o1 := &model.Observation{Integration: "det-a", SourceID: "1", DataType: "geo", Keys: []model.Alias{{Type: "fleet_uid", ID: "x-1"}},
		EntityType: "vehicle", Template: "asset", Traits: []string{"locatable"}, Tier: model.TierHuman, HasPos: true, Lat: 1, Lon: 1,
		SourceT: t0, Callsign: "OPERATOR-NAME"}
	o2 := *o1
	o2.Integration, o2.Tier, o2.SourceT, o2.Callsign = "det-b", model.TierSensor, t0+1000, "SENSOR-GUESS"
	e.Apply(o1)
	r := e.Apply(&o2)
	got, _ := e.Get("fleet-x-1")
	if got.Callsign != "OPERATOR-NAME" || len(got.Conflicts) != 1 {
		t.Fatalf("callsign %q conflicts %+v (%+v)", got.Callsign, got.Conflicts, r)
	}
}

func TestCoreModeSeqStartsAtWallClock(t *testing.T) {
	cfg := DefaultConfig()
	cfg.CoreMode = true
	e := New(cfg, func() int64 { return t0 + 60000 })
	r := e.Apply(obs(t, "evt.world.air.icao-a1b2c3", `{"uid":"icao-a1b2c3","lat":40.6,"lon":-73.7,"t":`+jsonInt(t0)+`}`))
	if r.Publish == nil || r.Publish.Entity.Seq != uint64(t0+60000) {
		t.Fatalf("core-mode seq %+v", r.Publish)
	}
}
