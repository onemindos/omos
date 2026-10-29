package engine

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/onemindos/omos/services/entity/internal/identity"
	"github.com/onemindos/omos/services/entity/internal/model"
)

func fleetMsg(uid, class string, t int64, meta string) string {
	return `{"uid":"` + uid + `","type":"fleet.` + class + `","lat":30.27,"lon":-97.74,"hae":150,"stale":60,"t":` +
		jsonInt(t) + `,"meta":` + meta + `}`
}

var droneCaps = []string{"onemind.task.goto.v1", "onemind.task.rtl.v1", "onemind.task.land.v1"}

// FEATURES H2/H3 + mission-task.md §4 capability gate: the record the policy service reads
// (omos #554 checkAsset) lists the task types the bridge advertises, and a later position
// (which carries no list) keeps them.
func TestSimDroneAndRobotAdvertiseCapabilities(t *testing.T) {
	c := &clock{t0 + 5000}
	e := newEngine(c)
	e.Apply(obs(t, "tel.fleet.drone.sim-x1.position", fleetMsg("sim-x1", "drone", t0, `{"battery_pct":80}`)))
	r := e.Apply(obs(t, "tel.fleet.drone.sim-x1.status", fleetMsg("sim-x1", "drone", t0+1000,
		`{"battery_pct":80,"connection":"online","commands_accepted":"goto,rtl,land,arm,disarm,estop,takeoff"}`)))
	if r.Publish == nil || !reflect.DeepEqual(r.Publish.Entity.Capabilities, droneCaps) {
		t.Fatalf("status publish %+v", r.Publish)
	}
	e.Apply(obs(t, "tel.fleet.drone.sim-x1.position", fleetMsg("sim-x1", "drone", t0+2000, `{"battery_pct":79}`)))
	got, _ := e.Get("sim-x1")
	if !reflect.DeepEqual(got.Capabilities, droneCaps) || got.Health["connection"] != "online" || !got.Indicators.Simulated {
		t.Fatalf("after position: caps %v health %v", got.Capabilities, got.Health)
	}

	e.Apply(obs(t, "tel.fleet.robot.sim-tb1.position", fleetMsg("sim-tb1", "robot", t0, `{"kind":"ugv"}`)))
	e.Apply(obs(t, "tel.fleet.robot.sim-tb1.status", fleetMsg("sim-tb1", "robot", t0+1000,
		`{"kind":"ugv","battery_pct":90,"connection":"online","commands_accepted":"goto,hold,rtl,estop,estop_release"}`)))
	tb, _ := e.Get("sim-tb1")
	if want := []string{"onemind.task.goto.v1", "onemind.task.hold.v1", "onemind.task.rtl.v1"}; !reflect.DeepEqual(tb.Capabilities, want) {
		t.Fatalf("robot caps %v", tb.Capabilities)
	}
	b, _ := json.Marshal(tb)
	var m map[string]any
	_ = json.Unmarshal(b, &m)
	if _, ok := m["capabilities"].([]any); !ok {
		t.Fatalf("capabilities not on the wire: %s", b)
	}
	if err := validate(compileSchema(t), tb); err != nil {
		t.Fatal(err)
	}
}

// A capability change, or the bridge reporting the link offline, publishes at once.
func TestCapabilityAndConnectionChangesPublish(t *testing.T) {
	e := newEngine(&clock{t0 + 5000})
	e.Apply(obs(t, "tel.fleet.drone.sim-x1.position", fleetMsg("sim-x1", "drone", t0, `{"battery_pct":80}`)))
	e.Apply(obs(t, "tel.fleet.drone.sim-x1.status", fleetMsg("sim-x1", "drone", t0+1000,
		`{"battery_pct":80,"connection":"online","commands_accepted":"goto"}`)))
	r := e.Apply(obs(t, "tel.fleet.drone.sim-x1.status", fleetMsg("sim-x1", "drone", t0+2000,
		`{"battery_pct":80,"connection":"online","commands_accepted":"goto,rtl"}`)))
	if r.Publish == nil || len(r.Publish.Entity.Capabilities) != 2 {
		t.Fatalf("capability change not published: %+v", r)
	}
	r = e.Apply(obs(t, "tel.fleet.drone.sim-x1.status", fleetMsg("sim-x1", "drone", t0+3000,
		`{"battery_pct":80,"connection":"offline","commands_accepted":"goto,rtl"}`)))
	if r.Publish == nil || r.Publish.Entity.Health["connection"] != "offline" {
		t.Fatalf("offline not published: %+v", r)
	}
	r = e.Apply(obs(t, "tel.fleet.drone.sim-x1.status", fleetMsg("sim-x1", "drone", t0+4000,
		`{"battery_pct":80,"connection":"online","commands_accepted":""}`)))
	if r.Publish == nil || len(r.Publish.Entity.Capabilities) != 0 {
		t.Fatalf("explicit empty list did not clear: %+v", r.Publish)
	}
}

// Never inferred: a non-sim fleet id, a feed, or a TAK/CoT drone gets no capabilities.
func TestNoCapabilitiesForNonSimFeedsOrCot(t *testing.T) {
	e := newEngine(&clock{t0 + 60000})
	e.Apply(obs(t, "tel.fleet.drone.x47.status", fleetMsg("x47", "drone", t0, `{"commands_accepted":"goto,rtl"}`)))
	if r, _ := e.Get("fleet-x47"); r.Capabilities != nil {
		t.Fatalf("non-sim fleet asset: %v", r.Capabilities)
	}
	last, _ := replay(t, e)
	for id, p := range last {
		if p.Entity.Capabilities != nil && id != "sim-x1" {
			t.Fatalf("%s got capabilities %v", id, p.Entity.Capabilities)
		}
	}
	// A hand-built observation claiming capabilities from a non-fleet integration is ignored.
	o := &model.Observation{Integration: "feed-aircraft", SourceID: "sim-a1", DataType: "geo",
		Keys: []model.Alias{{Type: identity.KeySim, ID: "sim-a1"}}, EntityType: "drone", Template: "asset",
		Traits: []string{"locatable", "taskable"}, Simulated: true, HasPos: true, Lat: 1, Lon: 1, SourceT: t0 + 59000, RecvT: t0 + 59000, StaleS: 60,
		Capabilities: droneCaps, HasCapabilities: true}
	res := e.Apply(o)
	if res.Publish == nil || res.Publish.Entity.EntityID != "sim-a1" || res.Publish.Entity.Capabilities != nil {
		t.Fatalf("feed integration set capabilities: %+v", res)
	}
}

// Restart (entity.md §5): capabilities survive Restore from the last ENT record.
func TestRestoreKeepsCapabilities(t *testing.T) {
	c := &clock{t0 + 5000}
	e1 := newEngine(c)
	e1.Apply(obs(t, "tel.fleet.drone.sim-x1.position", fleetMsg("sim-x1", "drone", t0, `{"battery_pct":80}`)))
	p := e1.Apply(obs(t, "tel.fleet.drone.sim-x1.status", fleetMsg("sim-x1", "drone", t0+1000,
		`{"battery_pct":80,"commands_accepted":"goto,rtl,land"}`))).Publish
	b, _ := json.Marshal(p.Entity)
	var r model.Entity
	_ = json.Unmarshal(b, &r)
	e2 := newEngine(c)
	e2.Restore(r)
	q := e2.Apply(obs(t, "tel.fleet.drone.sim-x1.position", fleetMsg("sim-x1", "drone", t0+2000, `{"battery_pct":79}`)))
	got, _ := e2.Get("sim-x1")
	if !reflect.DeepEqual(got.Capabilities, droneCaps) {
		t.Fatalf("restored caps %v (publish %v)", got.Capabilities, q.Publish != nil)
	}
}
