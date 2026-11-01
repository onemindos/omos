package engine

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/onemindos/omos/services/entity/internal/model"
)

func telDrone(uid string, lat float64, t int64) string {
	return fmt.Sprintf(`{"uid":%q,"type":"fleet.drone","lat":%v,"lon":-97.7431,"hae":150,"stale":60,"t":%d,`+
		`"meta":{"node_id":"build-mac","mavlink_sysid":1}}`, uid, lat, t)
}

// A shared MAVLink sysid (1 is the ArduPilot and PX4 SITL default) is a soft key: it never
// merges two drones (entity.md §5: auto-merge only on a shared hard key or a declared link).
func TestSameSysidTwoSimDronesAreTwoEntities(t *testing.T) {
	e := newEngine(&clock{t0 + 5000})
	r1 := e.Apply(obs(t, "tel.fleet.drone.sim-x1.position", telDrone("sim-x1", 30.0, t0)))
	r2 := e.Apply(obs(t, "tel.fleet.drone.sim-x2.position", telDrone("sim-x2", 31.0, t0+1000)))
	if r1.Publish == nil || r2.Publish == nil {
		t.Fatalf("want two publishes, got %+v %+v", r1, r2)
	}
	a, okA := e.Get("sim-x1")
	b, okB := e.Get("sim-x2")
	if !okA || !okB || a.Lat != 30.0 || b.Lat != 31.0 || len(a.Sources) != 1 || len(b.Sources) != 1 {
		t.Fatalf("merged: x1=%v %+v x2=%v %+v", okA, a.Location, okB, b.Location)
	}
	found := false
	for _, k := range b.Aliases.IDs {
		found = found || (k.Type == "mavlink_sysid" && k.ID == "1")
	}
	if !found {
		t.Fatalf("sysid not kept in aliases: %+v", b.Aliases.IDs)
	}
}

func eud(uid string, lat float64) string {
	return fmt.Sprintf(`<event version="2.0" uid=%q type="a-f-G-U-C" how="m-g" time="2026-10-29T12:00:03.000Z" `+
		`start="2026-10-29T12:00:03.000Z" stale="2026-10-29T12:06:03.000Z"><point lat="%v" lon="-97.7429" hae="210" ce="9.9" le="9999999.0"/>`+
		`<detail><takv device="Pixel 8" platform="ATAK-CIV" os="34" version="5.4.0"/><contact callsign="ALPHA"/>`+
		`<__group name="Cyan" role="Team Member"/></detail></event>`, uid, lat)
}

// A shared TAK callsign is a soft key: two devices called ALPHA stay two entities.
func TestSameCallsignTwoEUDsAreTwoEntities(t *testing.T) {
	e := newEngine(&clock{t0 + 60000})
	e.Apply(obs(t, "tak.cot.out", eud("ANDROID-aaaa", 30.0)))
	e.Apply(obs(t, "tak.cot.out", eud("ANDROID-bbbb", 31.0)))
	a, okA := e.Get("tak-android-aaaa")
	b, okB := e.Get("tak-android-bbbb")
	if !okA || !okB || a.Lat != 30.0 || b.Lat != 31.0 || len(a.Sources) != 1 {
		t.Fatalf("merged: a=%v lat %v sources %d, b=%v lat %v", okA, a.Lat, len(a.Sources), okB, b.Lat)
	}
}

// A failed final `live:false` publish is retried by the next Sweep, and the tombstone is
// deleted only once the final record is acked.
func TestFailedFinalPublishIsRetried(t *testing.T) {
	c := &clock{t0}
	e := newEngine(c)
	e.Apply(obs(t, "evt.world.air.icao-a1b2c3", `{"uid":"icao-a1b2c3","lat":40.6,"lon":-73.7,"stale":30,"t":`+jsonInt(t0)+`}`))
	c.ms = t0 + 120000
	ps := e.Sweep()
	if len(ps) != 1 || ps[0].Entity.Live {
		t.Fatalf("want final record, got %+v", ps)
	}
	final := ps[0]
	e.PublishFailed(final)
	c.ms += 1000
	ps = e.Sweep()
	if len(ps) != 1 || ps[0].MsgID != final.MsgID || ps[0].Entity.Live {
		t.Fatalf("final record not retried: %+v", ps)
	}
	if _, ok := e.Get("icao-a1b2c3"); ok {
		t.Fatal("tombstone visible to reads")
	}
	e.Published(ps[0])
	if n := e.Counters()["entities"]; n != 0 {
		t.Fatalf("tombstone kept after ack: %d entities", n)
	}
	// The entity comes back: seq continues past the final record.
	c.ms += 1000
	r := e.Apply(obs(t, "evt.world.air.icao-a1b2c3", `{"uid":"icao-a1b2c3","lat":40.6,"lon":-73.7,"stale":30,"t":`+jsonInt(c.ms)+`}`))
	if r.Publish == nil || r.Publish.Entity.Seq != 1 {
		t.Fatalf("new life after ack: %+v", r.Publish)
	}
}

// An entity that reappears while its final record is unacked continues seq and supersedes it.
func TestReappearanceSupersedesUnackedFinal(t *testing.T) {
	c := &clock{t0}
	e := newEngine(c)
	e.Apply(obs(t, "evt.world.air.icao-a1b2c3", `{"uid":"icao-a1b2c3","lat":40.6,"lon":-73.7,"stale":30,"t":`+jsonInt(t0)+`}`))
	c.ms = t0 + 120000
	final := e.Sweep()[0]
	e.PublishFailed(final)
	r := e.Apply(obs(t, "evt.world.air.icao-a1b2c3", `{"uid":"icao-a1b2c3","lat":40.6,"lon":-73.7,"stale":30,"t":`+jsonInt(c.ms)+`}`))
	if r.Publish == nil || !r.Publish.Entity.Live || r.Publish.Entity.Seq != final.Entity.Seq+1 {
		t.Fatalf("reappearance: %+v (final seq %d)", r.Publish, final.Entity.Seq)
	}
	for _, p := range e.Sweep() {
		if !p.Entity.Live {
			t.Fatalf("stale final retried after reappearance: %+v", p)
		}
	}
}

// Restore keeps each source's tier, so a sensor source is not promoted after a restart.
func TestRestoreKeepsSourceTier(t *testing.T) {
	c := &clock{t0 + 60000}
	e1 := newEngine(c)
	o1 := &model.Observation{Integration: "self", SourceID: "1", DataType: "geo", Keys: []model.Alias{{Type: "fleet_uid", ID: "v-1"}},
		EntityType: "vehicle", Template: "asset", Traits: []string{"locatable"}, Tier: model.TierSelfReport, HasPos: true,
		Lat: 1, Lon: 1, SourceT: t0, Callsign: "SELF"}
	o2 := *o1
	o2.Integration, o2.Tier, o2.SourceT, o2.Callsign = "radar", model.TierSensor, t0+1000, ""
	e1.Apply(o1)
	p := e1.Apply(&o2).Publish
	if p == nil {
		t.Fatal("no publish")
	}
	tiers := map[string]string{}
	for _, s := range p.Entity.Sources {
		tiers[s.Integration] = s.Tier
	}
	if tiers["self"] != "self_report" || tiers["radar"] != "sensor" {
		t.Fatalf("sources[].tier: %v", tiers)
	}
	b, _ := json.Marshal(p.Entity)
	var r model.Entity
	_ = json.Unmarshal(b, &r)
	e2 := newEngine(c)
	e2.Restore(r)
	tr := e2.tracks[r.EntityID]
	if tr.srcTier["radar|1"] != model.TierSensor || tr.srcTier["self|1"] != model.TierSelfReport {
		t.Fatalf("restored tiers: %v", tr.srcTier)
	}
}

func TestFleetUIDGetsFleetPrefix(t *testing.T) {
	e := newEngine(&clock{t0 + 5000})
	r := e.Apply(obs(t, "tel.fleet.drone.x47.position", `{"uid":"x47","type":"fleet.drone","lat":30,"lon":-97,"stale":60,"t":`+jsonInt(t0)+`}`))
	if r.Publish == nil || r.Publish.Entity.EntityID != "fleet-x47" {
		t.Fatalf("got %+v", r.Publish)
	}
}

func TestQueryPagesWithCursor(t *testing.T) {
	e := newEngine(&clock{t0 + 60000})
	for i := 0; i < 25; i++ {
		e.Apply(obs(t, "evt.world.air.x", fmt.Sprintf(`{"uid":"icao-%06x","lat":40,"lon":-73,"stale":60,"t":%d}`, i, t0)))
	}
	seen, cursor := 0, ""
	for page := 0; ; page++ {
		ents, more := e.Query(nil, nil, cursor, 10)
		seen += len(ents)
		if !more {
			break
		}
		cursor = ents[len(ents)-1].EntityID
		if page > 5 {
			t.Fatal("paging does not end")
		}
	}
	if seen != 25 {
		t.Fatalf("paged %d of 25", seen)
	}
}
