package adapt

import (
	"reflect"
	"testing"
)

// Status envelopes exactly as the bridges publish them (meta trimmed to the fields that matter).
const (
	// omos #550 MAVLink bridge statusEnvelope.
	mavStatus = `{"uid":"sim-x1","type":"fleet.drone","lat":30.2672,"lon":-97.7431,"hae":152.4,"course":-1,"speed":0,` +
		`"callsign":"SIM-X1","stale":60,"t":1793275236000,"meta":{"battery_pct":79,"mode":"GUIDED","armed":true,` +
		`"connection":"online","link_rssi":-1,"kind":"multirotor","simulated":true,` +
		`"commands_accepted":"goto,rtl,land,arm,disarm,estop,takeoff"}}`
	// omos #566 ROS 2 bridge status_envelope.
	rosStatus = `{"uid":"sim-tb1","type":"fleet.robot","lat":30.27,"lon":-97.74,"hae":150,"course":-1,"speed":0,` +
		`"callsign":"SIM-TB1","stale":30,"t":1793275236000,"meta":{"kind":"ugv","battery_pct":90,"mode":"idle",` +
		`"connection":"online","link_rssi":-1,"simulated":true,` +
		`"commands_accepted":"goto,hold,rtl,estop,estop_release"}}`
)

func TestMavlinkStatusGivesDroneTaskTypes(t *testing.T) {
	o, _, err := Fleet([]byte(mavStatus), "tel.fleet.drone.sim-x1.status")
	if err != nil {
		t.Fatal(err)
	}
	want := []string{"onemind.task.goto.v1", "onemind.task.rtl.v1", "onemind.task.land.v1"}
	if !o.HasCapabilities || !reflect.DeepEqual(o.Capabilities, want) {
		t.Fatalf("capabilities %v %v, want %v", o.HasCapabilities, o.Capabilities, want)
	}
	if o.Kind != "multirotor" || o.Health["connection"] != "online" {
		t.Fatalf("kind %q health %v", o.Kind, o.Health)
	}
}

func TestRosStatusGivesRobotTaskTypes(t *testing.T) {
	o, _, err := Fleet([]byte(rosStatus), "tel.fleet.robot.sim-tb1.status")
	if err != nil {
		t.Fatal(err)
	}
	want := []string{"onemind.task.goto.v1", "onemind.task.hold.v1", "onemind.task.rtl.v1"}
	if !o.HasCapabilities || !reflect.DeepEqual(o.Capabilities, want) {
		t.Fatalf("capabilities %v, want %v", o.Capabilities, want)
	}
	if o.Kind != "ugv" {
		t.Fatalf("robot kind %q", o.Kind)
	}
}

func TestClassBoundsAdvertisedVerbs(t *testing.T) {
	meta := map[string]any{"commands_accepted": []any{"goto", "land", "orbit", "observe", "onemind.task.hold.v1", "bogus"}}
	robot, _ := FleetCapabilities("robot", "sim-tb1", []string{"locatable", "trackable", "taskable"}, meta)
	if want := []string{"onemind.task.goto.v1", "onemind.task.hold.v1"}; !reflect.DeepEqual(robot, want) {
		t.Fatalf("robot %v, want %v (no land/orbit for ground, no observe without a sensor)", robot, want)
	}
	drone, _ := FleetCapabilities("drone", "sim-x1", []string{"locatable", "trackable", "taskable", "observable"}, meta)
	if want := []string{"onemind.task.goto.v1", "onemind.task.hold.v1", "onemind.task.land.v1",
		"onemind.task.orbit.v1", "onemind.task.observe.v1"}; !reflect.DeepEqual(drone, want) {
		t.Fatalf("drone %v, want %v", drone, want)
	}
	if node, ok := FleetCapabilities("node", "sim-n1", []string{"locatable"}, meta); !ok || len(node) != 0 {
		t.Fatalf("node %v %v, want explicit empty", node, ok)
	}
}

func TestNoCapabilitiesWithoutSimIDOrAdvertisement(t *testing.T) {
	real := `{"uid":"x47","type":"fleet.drone","lat":30,"lon":-97,"stale":60,"t":1793275236000,` +
		`"meta":{"commands_accepted":"goto,rtl,land"}}`
	o, _, err := Fleet([]byte(real), "tel.fleet.drone.x47.status")
	if err != nil {
		t.Fatal(err)
	}
	if o.HasCapabilities || o.Capabilities != nil {
		t.Fatalf("non-sim fleet id got capabilities %v", o.Capabilities)
	}
	pos := `{"uid":"sim-x1","type":"fleet.drone","lat":30,"lon":-97,"stale":60,"t":1793275236000,"meta":{"battery_pct":80}}`
	o, _, _ = Fleet([]byte(pos), "tel.fleet.drone.sim-x1.position")
	if o.HasCapabilities {
		t.Fatal("a position without commands_accepted must not touch capabilities")
	}
}

// Feeds and CoT never get capabilities, even for a sim uid or a drone CoT type.
func TestFeedsAndCotNeverGetCapabilities(t *testing.T) {
	for _, fx := range fixtures(t) {
		o, _, err := Parse(fx.Subject, []byte(fx.Data))
		if err != nil || o == nil {
			continue
		}
		if o.HasCapabilities && (len(fx.Subject) < 10 || fx.Subject[:10] != "tel.fleet.") {
			t.Fatalf("%s got capabilities %v", fx.Subject, o.Capabilities)
		}
	}
}

func TestBridgeConnectionIsHonoured(t *testing.T) {
	off := `{"uid":"sim-x1","type":"fleet.drone","lat":30,"lon":-97,"stale":60,"t":1793275236000,` +
		`"meta":{"battery_pct":80,"connection":"offline","commands_accepted":"goto"}}`
	o, _, _ := Fleet([]byte(off), "tel.fleet.drone.sim-x1.status")
	if o.Health["connection"] != "offline" {
		t.Fatalf("health %v", o.Health)
	}
}
