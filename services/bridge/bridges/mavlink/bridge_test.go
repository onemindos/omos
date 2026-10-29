package main

import (
	"context"
	"crypto/ed25519"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"net"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/bluenviron/gomavlib/v3"
	"github.com/bluenviron/gomavlib/v3/pkg/dialects/ardupilotmega"
	"github.com/bluenviron/gomavlib/v3/pkg/dialects/common"

	ct "github.com/onemindos/omos-bridge/bridges/mavlink/internal/commandtest"
	"github.com/onemindos/omos-bridge/bridges/mavlink/internal/fence"
	"github.com/onemindos/omos-bridge/bridges/mavlink/internal/geoid"
)

// memBus records publishes.
type memBus struct {
	mu   sync.Mutex
	msgs []busMsg
}

type busMsg struct {
	Subject, MsgID string
	Data           []byte
}

func (m *memBus) Publish(subject, msgID string, data []byte) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.msgs = append(m.msgs, busMsg{subject, msgID, append([]byte(nil), data...)})
	return nil
}

func (m *memBus) bySubject(subject string) []busMsg {
	m.mu.Lock()
	defer m.mu.Unlock()
	var out []busMsg
	for _, x := range m.msgs {
		if x.Subject == subject {
			out = append(out, x)
		}
	}
	return out
}

// fakeVehicle is an in-process MAVLink vehicle (sysid 1, ArduCopter quad) that streams
// HEARTBEAT, SYS_STATUS and GLOBAL_POSITION_INT and ACKs COMMAND_LONG / COMMAND_INT.
type fakeVehicle struct {
	node     *gomavlib.Node
	mu       sync.Mutex
	longs    []*common.MessageCommandLong
	ints     []*common.MessageCommandInt
	result   common.MAV_RESULT
	silent   bool
	noInt    bool  // drop COMMAND_INT unanswered (a vehicle that reads only COMMAND_LONG)
	intUnsup bool  // answer COMMAND_INT with MAV_RESULT_UNSUPPORTED
	ackTo    uint8 // COMMAND_ACK target_system (0 = unset)
	paused   bool  // stop telemetry (link loss) but still answer commands
	lat, lon float64
	altMSL   float64
	gpsN     float64 // if non-zero, also send GPS_RAW_INT with alt_ellipsoid = alt + gpsN
}

func startFakeVehicle(t *testing.T, ctx context.Context, addr string, lat, lon, alt float64) *fakeVehicle {
	t.Helper()
	n := &gomavlib.Node{
		Endpoints:        []gomavlib.EndpointConf{gomavlib.EndpointUDPClient{Address: addr}},
		Dialect:          ardupilotmega.Dialect,
		OutVersion:       gomavlib.V2,
		OutSystemID:      1,
		HeartbeatDisable: true,
	}
	if err := n.Initialize(); err != nil {
		t.Fatal(err)
	}
	f := &fakeVehicle{node: n, lat: lat, lon: lon, altMSL: alt, result: common.MAV_RESULT_ACCEPTED}
	go func() {
		tick := time.NewTicker(50 * time.Millisecond)
		defer tick.Stop()
		i := 0
		for {
			select {
			case <-ctx.Done():
				n.Close()
				return
			case <-tick.C:
				f.mu.Lock()
				lat, lon, alt, gpsN, paused := f.lat, f.lon, f.altMSL, f.gpsN, f.paused
				f.mu.Unlock()
				if paused {
					continue
				}
				if i%5 == 0 {
					_ = n.WriteMessageAll(&ardupilotmega.MessageHeartbeat{
						Type: ardupilotmega.MAV_TYPE_QUADROTOR, Autopilot: ardupilotmega.MAV_AUTOPILOT_ARDUPILOTMEGA,
						BaseMode: ardupilotmega.MAV_MODE_FLAG_CUSTOM_MODE_ENABLED, CustomMode: 4,
						SystemStatus: ardupilotmega.MAV_STATE_STANDBY, MavlinkVersion: 3,
					})
					_ = n.WriteMessageAll(&ardupilotmega.MessageSysStatus{BatteryRemaining: 87, VoltageBattery: 15800})
					if gpsN != 0 {
						_ = n.WriteMessageAll(&ardupilotmega.MessageGpsRawInt{
							FixType: ardupilotmega.GPS_FIX_TYPE_3D_FIX, Lat: int32(lat * 1e7), Lon: int32(lon * 1e7),
							Alt: int32(alt * 1000), AltEllipsoid: int32((alt + gpsN) * 1000),
						})
					}
				}
				i++
				_ = n.WriteMessageAll(&ardupilotmega.MessageGlobalPositionInt{
					Lat: int32(lat * 1e7), Lon: int32(lon * 1e7), Alt: int32(alt * 1000), RelativeAlt: 20000,
					Vx: 300, Vy: 400, Hdg: 9000,
				})
			case evt, ok := <-n.Events():
				if !ok {
					return
				}
				fr, ok := evt.(*gomavlib.EventFrame)
				if !ok {
					continue
				}
				f.mu.Lock()
				res, silent, ackTo := f.result, f.silent, f.ackTo
				var cmd common.MAV_CMD
				switch m := fr.Message().(type) {
				case *ardupilotmega.MessageCommandLong:
					f.longs = append(f.longs, m)
					cmd = m.Command
				case *ardupilotmega.MessageCommandInt:
					f.ints = append(f.ints, m)
					cmd = m.Command
					silent = silent || f.noInt
					if f.intUnsup {
						res = common.MAV_RESULT_UNSUPPORTED
					}
				default:
					f.mu.Unlock()
					continue
				}
				f.mu.Unlock()
				if !silent {
					_ = n.WriteMessageTo(fr.Channel, &ardupilotmega.MessageCommandAck{Command: cmd, Result: res, TargetSystem: ackTo})
				}
			}
		}
	}()
	return f
}

func freeUDPAddr(t *testing.T) string {
	t.Helper()
	c, err := net.ListenPacket("udp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()
	return c.LocalAddr().String()
}

const testFence = `{
  "keep_in": [{"id": "sim-range", "circle": {"lat": 0, "lon": 0, "radius_m": 2000}}],
  "keep_out": [{"id": "no-fly", "circle": {"lat": 0.005, "lon": 0.005, "radius_m": 150}}],
  "max_hae_m": 300
}`

type harness struct {
	t    *testing.T
	b    *Bridge
	bus  *memBus
	veh  *fakeVehicle
	priv ed25519.PrivateKey
	ctx  context.Context
}

func newHarness(t *testing.T) *harness {
	t.Helper()
	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)
	pub, priv := ct.NewKey()
	fc, err := fence.Parse([]byte(testFence))
	if err != nil {
		t.Fatal(err)
	}
	addr := freeUDPAddr(t)
	cfg := Config{
		Source: "bridge-mavlink", AssetID: "sim-x1", Callsign: "SIM-X1", SysID: 1, GCSSysID: 254,
		Keys: map[string]ed25519.PublicKey{"policy-test": pub}, SigRequired: true, Fence: fc,
		MaxTTL: time.Minute, ClockSkew: 5 * time.Second, PositionInterval: 100 * time.Millisecond,
		StatusInterval: time.Second, StaleS: 10, LinkTimeout: 3 * time.Second,
		AckTimeout: 300 * time.Millisecond, AckRetries: 2,
	}
	node := &gomavlib.Node{
		Endpoints: []gomavlib.EndpointConf{gomavlib.EndpointUDPServer{Address: addr}},
		Dialect:   ardupilotmega.Dialect, OutVersion: gomavlib.V2, OutSystemID: 254,
	}
	if err := node.Initialize(); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(node.Close)
	mv := newMavVehicle(node, 1, 254, cfg.AckTimeout, cfg.AckRetries)
	bus := &memBus{}
	b, err := NewBridge(cfg, bus, mv)
	if err != nil {
		t.Fatal(err)
	}
	go pump(ctx, node, mv, b)
	fv := startFakeVehicle(t, ctx, addr, 0.001, 0.001, 50)
	h := &harness{t: t, b: b, bus: bus, veh: fv, priv: priv, ctx: ctx}
	h.waitFor(func() bool {
		return len(bus.bySubject(b.positionSubject())) > 0 && b.state.online(time.Now(), time.Second)
	}, "telemetry")
	return h
}

func (h *harness) waitFor(cond func() bool, what string) {
	h.t.Helper()
	deadline := time.Now().Add(5 * time.Second)
	for time.Now().Before(deadline) {
		if cond() {
			return
		}
		time.Sleep(20 * time.Millisecond)
	}
	h.t.Fatalf("timeout waiting for %s", what)
}

// send builds a command, delivers it synchronously and returns the outcome data in order.
func (h *harness) send(o ct.Opts) []map[string]any {
	h.t.Helper()
	o.AssetID = "sim-x1"
	if o.Key == nil {
		o.Key = h.priv
	}
	subj, body, err := ct.Build(o)
	if err != nil {
		h.t.Fatal(err)
	}
	before := len(h.bus.bySubject(h.b.outcomeSubject()))
	h.b.HandleCommand(h.ctx, subj, body)
	var out []map[string]any
	for _, m := range h.bus.bySubject(h.b.outcomeSubject())[before:] {
		var env map[string]any
		if err := json.Unmarshal(m.Data, &env); err != nil {
			h.t.Fatal(err)
		}
		if m.MsgID == "" || m.MsgID != env["id"] {
			h.t.Fatalf("outcome without Nats-Msg-Id = envelope id: %+v", m)
		}
		for _, k := range []string{"v", "id", "ts", "type", "source", "subject", "data"} {
			if _, ok := env[k]; !ok {
				h.t.Fatalf("outcome envelope missing %s: %s", k, m.Data)
			}
		}
		out = append(out, env["data"].(map[string]any))
	}
	return out
}

func outcomes(d []map[string]any) string {
	s := ""
	for _, x := range d {
		s += fmt.Sprintf("%v/%v ", x["outcome"], x["code"])
	}
	return s
}

func TestTelemetryEnvelopes(t *testing.T) {
	h := newHarness(t)
	h.b.PublishStatus()
	var pos GeoEnvelope
	msgs := h.bus.bySubject("tel.fleet.drone.sim-x1.position")
	if err := json.Unmarshal(msgs[len(msgs)-1].Data, &pos); err != nil {
		t.Fatal(err)
	}
	n, _ := geoid.Undulation(0.001, 0.001)
	if pos.UID != "sim-x1" || pos.Type != "fleet.drone" || pos.Callsign != "SIM-X1" || pos.T == 0 || pos.Stale != 10 {
		t.Fatalf("envelope identity fields: %+v", pos)
	}
	if math.Abs(pos.Lat-0.001) > 1e-7 || math.Abs(pos.Lon-0.001) > 1e-7 {
		t.Fatalf("lat/lon %v %v", pos.Lat, pos.Lon)
	}
	if math.Abs(pos.Hae-(50+n)) > 0.02 || pos.Meta["hae_source"] != "egm96" || math.Abs(n-17.16) > 1 {
		t.Fatalf("hae %.2f want MSL 50 + N %.2f (%v)", pos.Hae, n, pos.Meta["hae_source"])
	}
	// vx=3 m/s north, vy=4 m/s east -> 5 m/s, course 53.1 deg.
	if math.Abs(pos.Speed-5) > 0.01 || math.Abs(pos.Course-53.1) > 0.1 {
		t.Fatalf("speed %.2f course %.1f", pos.Speed, pos.Course)
	}
	if pos.Meta["mode"] != "GUIDED" || pos.Meta["battery_pct"].(float64) != 87 || pos.Meta["simulated"] != true {
		t.Fatalf("meta %+v", pos.Meta)
	}
	st := h.bus.bySubject("tel.fleet.drone.sim-x1.status")
	if len(st) == 0 {
		t.Fatal("no status published")
	}
	var s GeoEnvelope
	_ = json.Unmarshal(st[len(st)-1].Data, &s)
	if s.Meta["connection"] != "online" || s.Meta["kind"] != "multirotor" || s.Meta["autopilot"] != "ARDUPILOTMEGA" {
		t.Fatalf("status meta %+v", s.Meta)
	}
	// Every publish uses only the locked-grammar subjects.
	h.bus.mu.Lock()
	defer h.bus.mu.Unlock()
	for _, m := range h.bus.msgs {
		switch m.Subject {
		case "tel.fleet.drone.sim-x1.position", "tel.fleet.drone.sim-x1.status", "evt.drone.sim-x1.command", "evt.drone.sim-x1.task":
		default:
			t.Fatalf("unexpected subject %s", m.Subject)
		}
	}
}

func TestHAEPrefersAutopilotEllipsoid(t *testing.T) {
	h := newHarness(t)
	h.veh.mu.Lock()
	h.veh.gpsN = 12.5
	h.veh.mu.Unlock()
	h.waitFor(func() bool {
		msgs := h.bus.bySubject(h.b.positionSubject())
		var p GeoEnvelope
		_ = json.Unmarshal(msgs[len(msgs)-1].Data, &p)
		return p.Meta["hae_source"] == "gps_raw_int" && math.Abs(p.Hae-62.5) < 0.01
	}, "hae from GPS_RAW_INT.alt_ellipsoid")
}

func TestApprovedCommandsReachVehicle(t *testing.T) {
	h := newHarness(t)
	cases := []struct {
		verb  string
		data  any
		cmd   common.MAV_CMD
		isInt bool
	}{
		{"arm", nil, common.MAV_CMD_COMPONENT_ARM_DISARM, false},
		{"takeoff", map[string]any{"hae_m": 77.0}, common.MAV_CMD_NAV_TAKEOFF, false},
		{"goto", map[string]any{"lat": 0.002, "lon": -0.003, "hae_m": 100.0, "speed_mps": 6.0}, common.MAV_CMD_DO_REPOSITION, true},
		{"rtl", nil, common.MAV_CMD_NAV_RETURN_TO_LAUNCH, false},
		{"land", nil, common.MAV_CMD_NAV_LAND, false},
		{"disarm", nil, common.MAV_CMD_COMPONENT_ARM_DISARM, false},
		{"estop", nil, common.MAV_CMD_DO_FLIGHTTERMINATION, false},
	}
	for _, c := range cases {
		got := h.send(ct.Opts{Verb: c.verb, Data: c.data, TaskID: "tsk-test-" + c.verb})
		if len(got) != 2 || got[0]["outcome"] != "acknowledged" || got[1]["outcome"] != "executed" {
			t.Fatalf("%s: outcomes %s", c.verb, outcomes(got))
		}
		mav := got[1]["mavlink"].(map[string]any)
		if mav["command"] != c.cmd.String() || mav["result"].(float64) != 0 || got[0]["sig_verified"] != true {
			t.Fatalf("%s: mavlink %+v", c.verb, got)
		}
	}
	h.veh.mu.Lock()
	defer h.veh.mu.Unlock()
	if len(h.veh.longs) != 6 || len(h.veh.ints) != 1 {
		t.Fatalf("vehicle got %d COMMAND_LONG, %d COMMAND_INT", len(h.veh.longs), len(h.veh.ints))
	}
	arm, tko, disarm, estop := h.veh.longs[0], h.veh.longs[1], h.veh.longs[4], h.veh.longs[5]
	if arm.Param1 != 1 || disarm.Param1 != 0 || estop.Param1 != 1 || arm.TargetSystem != 1 {
		t.Fatalf("arm/disarm/estop params %v %v %v", arm.Param1, disarm.Param1, estop.Param1)
	}
	nHere, _ := geoid.Undulation(0.001, 0.001)
	if math.Abs(float64(tko.Param7)-(77-nHere)) > 0.01 {
		t.Fatalf("takeoff param7 %.2f want HAE 77 - N %.2f", tko.Param7, nHere)
	}
	g := h.veh.ints[0]
	n, _ := geoid.Undulation(0.002, -0.003)
	if g.X != 20000 || g.Y != -30000 || g.Frame != common.MAV_FRAME_GLOBAL || math.Abs(float64(g.Z)-(100-n)) > 0.01 ||
		g.Param1 != 6 || g.Param2 != 1 || !math.IsNaN(float64(g.Param4)) {
		t.Fatalf("DO_REPOSITION fields %+v (want z=%.2f)", g, 100-n)
	}
	// Task-bearing commands also report on evt.drone.sim-x1.task.
	tasks := h.bus.bySubject("evt.drone.sim-x1.task")
	if len(tasks) != 14 {
		t.Fatalf("task reports %d, want 14", len(tasks))
	}
}

func TestRefusedCommandsNeverReachVehicle(t *testing.T) {
	h := newHarness(t)
	f := false
	_, otherKey := ct.NewKey()
	cases := []struct {
		name   string
		o      ct.Opts
		code   string
		detail string
	}{
		{"unapproved", ct.Opts{Verb: "rtl", Status: "pending"}, "not_authorized", ""},
		{"agent approved", ct.Opts{Verb: "goto", Data: map[string]any{"lat": 0.0, "lon": 0.0, "hae_m": 50.0},
			Approvers: []string{"agent:hermes/zeus/legacy"}}, "agent_cannot_confirm", ""},
		{"arm with one human", ct.Opts{Verb: "arm", Approvers: []string{"user:zeus"}}, "not_authorized", ""},
		{"not simulated", ct.Opts{Verb: "rtl", Simulated: &f}, "not_simulated", ""},
		{"forged signature", ct.Opts{Verb: "rtl", Key: otherKey}, "bad_signature", ""},
		{"unsigned", ct.Opts{Verb: "rtl", NoSig: true}, "bad_signature", ""},
		{"expired", ct.Opts{Verb: "rtl", Now: time.Now().Add(-2 * time.Minute)}, "expired", ""},
		{"unknown verb", ct.Opts{Verb: "orbit"}, "unknown_action", ""},
		{"goto keep-out", ct.Opts{Verb: "goto", Data: map[string]any{"lat": 0.005, "lon": 0.005, "hae_m": 50.0}}, "precondition_failed", "keep_out"},
		{"goto outside keep-in", ct.Opts{Verb: "goto", Data: map[string]any{"lat": 0.1, "lon": 0.0, "hae_m": 50.0}}, "precondition_failed", "keep_in"},
		{"goto above ceiling", ct.Opts{Verb: "goto", Data: map[string]any{"lat": 0.0, "lon": 0.0, "hae_m": 400.0}}, "precondition_failed", "keep_in"},
		{"land outside keep-in", ct.Opts{Verb: "land", Data: map[string]any{"lat": 0.1, "lon": 0.1}}, "precondition_failed", "keep_in"},
	}
	for _, c := range cases {
		got := h.send(c.o)
		if len(got) != 1 || got[0]["outcome"] != "refused" || got[0]["code"] != c.code {
			t.Fatalf("%s: outcomes %s (%v)", c.name, outcomes(got), got)
		}
		if c.detail != "" && got[0]["detail"] != c.detail {
			t.Fatalf("%s: detail %v want %s", c.name, got[0]["detail"], c.detail)
		}
	}
	// Replay: the same signed command twice.
	subj, body, _ := ct.Build(ct.Opts{AssetID: "sim-x1", Verb: "rtl", Key: h.priv})
	h.b.HandleCommand(h.ctx, subj, body)
	h.b.HandleCommand(h.ctx, subj, body)
	last := h.bus.bySubject(h.b.outcomeSubject())
	var env map[string]any
	_ = json.Unmarshal(last[len(last)-1].Data, &env)
	if d := env["data"].(map[string]any); d["code"] != "replayed" {
		t.Fatalf("replay outcome %v", d)
	}
	// Garbage payload.
	h.b.HandleCommand(h.ctx, "cmd.fleet.drone.sim-x1.rtl", []byte("not json"))
	last = h.bus.bySubject(h.b.outcomeSubject())
	_ = json.Unmarshal(last[len(last)-1].Data, &env)
	if d := env["data"].(map[string]any); d["code"] != "schema_invalid" {
		t.Fatalf("garbage outcome %v", d)
	}
	h.veh.mu.Lock()
	defer h.veh.mu.Unlock()
	if len(h.veh.longs) != 1 || len(h.veh.ints) != 0 {
		t.Fatalf("refused commands reached the vehicle: %d long, %d int (only the first replay copy may)", len(h.veh.longs), len(h.veh.ints))
	}
}

func TestVehicleRejectAndNoAck(t *testing.T) {
	h := newHarness(t)
	h.veh.mu.Lock()
	h.veh.result = common.MAV_RESULT_UNSUPPORTED
	h.veh.mu.Unlock()
	got := h.send(ct.Opts{Verb: "rtl"})
	if len(got) != 2 || got[1]["outcome"] != "failed" || got[1]["code"] != "vehicle_rejected" {
		t.Fatalf("unsupported: %s", outcomes(got))
	}
	h.veh.mu.Lock()
	h.veh.silent = true
	h.veh.mu.Unlock()
	got = h.send(ct.Opts{Verb: "land"})
	if len(got) != 2 || got[1]["outcome"] != "failed" || got[1]["code"] != "no_ack" {
		t.Fatalf("silent: %s", outcomes(got))
	}
}

func TestOfflineVehicleRefusedButEstopSent(t *testing.T) {
	h := newHarness(t)
	h.veh.mu.Lock()
	h.veh.paused = true
	h.veh.mu.Unlock()
	time.Sleep(100 * time.Millisecond)
	h.b.offset.Store(int64(10 * time.Second)) // last heartbeat now >= 10 s old
	got := h.send(ct.Opts{Verb: "rtl", Now: time.Now().Add(10 * time.Second)})
	if len(got) != 1 || got[0]["code"] != "precondition_failed" || got[0]["detail"] != "offline" {
		t.Fatalf("offline rtl: %s %v", outcomes(got), got)
	}
	got = h.send(ct.Opts{Verb: "estop", Now: time.Now().Add(10 * time.Second)})
	if len(got) != 2 || got[0]["outcome"] != "acknowledged" {
		t.Fatalf("estop must not be blocked by link state: %s", outcomes(got))
	}
}

func TestDryRunNotSent(t *testing.T) {
	h := newHarness(t)
	got := h.send(ct.Opts{Verb: "rtl", Mutate: func(m map[string]any) { m["dry_run"] = true }})
	if len(got) != 1 || got[0]["outcome"] != "acknowledged" {
		t.Fatalf("dry run: %s", outcomes(got))
	}
	h.veh.mu.Lock()
	defer h.veh.mu.Unlock()
	if len(h.veh.longs) != 0 {
		t.Fatal("dry run reached the vehicle")
	}
}

func TestConfigRefusesNonSim(t *testing.T) {
	t.Setenv("ASSET_ID", "drone-alpha-1")
	if _, err := loadConfig(); err == nil {
		t.Fatal("loadConfig accepted a non-sim asset id")
	}
	t.Setenv("ASSET_ID", "sim-x1")
	t.Setenv("POLICY_PUBKEYS", "k1=notbase64!")
	if _, err := loadConfig(); err == nil {
		t.Fatal("loadConfig accepted a bad key")
	}
	t.Setenv("POLICY_PUBKEYS", "")
	t.Setenv("MAVLINK_ENDPOINT", "carrier-pigeon:x")
	if _, err := loadConfig(); err == nil {
		t.Fatal("loadConfig accepted a bad endpoint")
	}
	t.Setenv("MAVLINK_ENDPOINT", "")
	t.Setenv("DRONE_IP", "10.0.0.5")
	c, err := loadConfig()
	if err != nil || c.EndpointS != "udpclient:10.0.0.5:14550" || !c.SigRequired {
		t.Fatalf("legacy DRONE_IP: %v %+v", err, c.EndpointS)
	}
	if _, err := NewBridge(c, &memBus{}, nil); err == nil {
		t.Fatal("bridge started with SIG_REQUIRED and no keys")
	}
}

func TestRoundDropsNegativeZero(t *testing.T) {
	if b, _ := json.Marshal(round(-0.0001, 2)); string(b) != "0" {
		t.Fatalf("round(-0.0001) marshals as %s", b)
	}
}

// A vehicle that reads only COMMAND_LONG (the scripted drone) still executes goto:
// DO_REPOSITION is re-sent as COMMAND_LONG, and later gotos go straight to COMMAND_LONG.
func TestRepositionFallsBackToCommandLong(t *testing.T) {
	for _, mode := range []string{"no ack", "unsupported"} {
		t.Run(mode, func(t *testing.T) {
			h := newHarness(t)
			h.veh.mu.Lock()
			h.veh.noInt, h.veh.intUnsup = mode == "no ack", mode == "unsupported"
			h.veh.mu.Unlock()
			goTo := map[string]any{"lat": 0.002, "lon": -0.003, "hae_m": 100.0}
			for i := 0; i < 2; i++ {
				got := h.send(ct.Opts{Verb: "goto", Data: goTo})
				if len(got) != 2 || got[1]["outcome"] != "executed" {
					t.Fatalf("goto %d: %s", i, outcomes(got))
				}
			}
			h.veh.mu.Lock()
			defer h.veh.mu.Unlock()
			wantInts := 2 // AckRetries in the harness
			if mode == "unsupported" {
				wantInts = 1
			}
			if len(h.veh.ints) != wantInts || len(h.veh.longs) != 2 {
				t.Fatalf("vehicle got %d COMMAND_INT (want %d), %d COMMAND_LONG (want 2)", len(h.veh.ints), wantInts, len(h.veh.longs))
			}
			l := h.veh.longs[0]
			n, _ := geoid.Undulation(0.002, -0.003)
			if l.Command != common.MAV_CMD_DO_REPOSITION || math.Abs(float64(l.Param5)-0.002) > 1e-6 ||
				math.Abs(float64(l.Param6)+0.003) > 1e-6 || math.Abs(float64(l.Param7)-(100-n)) > 0.01 || l.Param2 != 1 {
				t.Fatalf("COMMAND_LONG DO_REPOSITION %+v", l)
			}
		})
	}
}

// An ACK addressed to another GCS on a shared link must not complete our command.
func TestAckForOtherGCSIgnored(t *testing.T) {
	h := newHarness(t)
	h.veh.mu.Lock()
	h.veh.ackTo = 200
	h.veh.mu.Unlock()
	if got := h.send(ct.Opts{Verb: "rtl"}); len(got) != 2 || got[1]["code"] != "no_ack" {
		t.Fatalf("ACK for sysid 200: %s", outcomes(got))
	}
	h.veh.mu.Lock()
	h.veh.ackTo = 254
	h.veh.mu.Unlock()
	if got := h.send(ct.Opts{Verb: "rtl"}); len(got) != 2 || got[1]["outcome"] != "executed" {
		t.Fatalf("ACK for our sysid 254: %s", outcomes(got))
	}
}

// Task reports use Nats-Msg-Id {task_id}:{status_version}, and status_version keeps
// increasing across the commands of one task.
func TestTaskStatusVersions(t *testing.T) {
	h := newHarness(t)
	for _, c := range []ct.Opts{
		{Verb: "takeoff", Data: map[string]any{"hae_m": 70.0}},
		{Verb: "goto", Data: map[string]any{"lat": 0.002, "lon": 0.002, "hae_m": 90.0}},
	} {
		c.TaskID = "tsk-multi"
		if got := h.send(c); len(got) != 2 || got[1]["outcome"] != "executed" {
			t.Fatalf("%s: %s", c.Verb, outcomes(got))
		}
	}
	tasks := h.bus.bySubject("evt.drone.sim-x1.task")
	if len(tasks) != 4 {
		t.Fatalf("task reports %d, want 4", len(tasks))
	}
	for i, m := range tasks {
		var env struct {
			Data struct {
				TaskID string `json:"task_id"`
				Ver    int    `json:"status_version"`
			} `json:"data"`
		}
		if err := json.Unmarshal(m.Data, &env); err != nil {
			t.Fatal(err)
		}
		if env.Data.Ver != i+1 || m.MsgID != fmt.Sprintf("tsk-multi:%d", i+1) {
			t.Fatalf("report %d: status_version %d, Nats-Msg-Id %q", i, env.Data.Ver, m.MsgID)
		}
	}
}

type errVehicle struct{ err error }

func (v errVehicle) Send(context.Context, MavRequest) (int, error) { return -1, v.err }

// A command the vehicle never received is not reported as vehicle_rejected.
func TestSendErrorsAreNotVehicleRejected(t *testing.T) {
	h := newHarness(t)
	cases := []struct {
		err          error
		code, detail string
	}{
		{fmt.Errorf("%w: sysid 1", ErrNotHeard), "precondition_failed", "offline"},
		{fmt.Errorf("%w: MAV_CMD_NAV_RETURN_TO_LAUNCH", ErrBusy), "conflict", ""},
		{context.Canceled, "no_ack", "aborted"},
		{errors.New("write: broken pipe"), "no_ack", "not_sent"},
	}
	for _, c := range cases {
		h.b.veh = errVehicle{c.err}
		got := h.send(ct.Opts{Verb: "rtl"})
		if len(got) != 2 || got[1]["outcome"] != "failed" || got[1]["code"] != c.code || (c.detail != "" && got[1]["detail"] != c.detail) {
			t.Fatalf("%v: %s %v", c.err, outcomes(got), got)
		}
	}
}

// hold is MAV_CMD_DO_REPOSITION (COMMAND_INT, frame GLOBAL, CHANGE_MODE) to the latest
// GLOBAL_POSITION_INT, the hold that task.abort sends (mission-task.md §6.2).
func TestHoldAtCurrentPosition(t *testing.T) {
	h := newHarness(t)
	h.veh.mu.Lock()
	h.veh.lat, h.veh.lon, h.veh.altMSL = 0.0012345, -0.0006789, 61.5
	h.veh.mu.Unlock()
	h.waitFor(func() bool {
		lat, _, _, ok := h.b.state.position()
		return ok && math.Abs(lat-0.0012345) < 1e-9
	}, "moved position")

	taskReport := func() map[string]any {
		tasks := h.bus.bySubject("evt.drone.sim-x1.task")
		var env struct{ Data map[string]any }
		if err := json.Unmarshal(tasks[len(tasks)-1].Data, &env); err != nil {
			t.Fatal(err)
		}
		return env.Data
	}

	// task.abort: hold with reason abort, task-bound -> executed, task done_not_ok + cancelled.
	got := h.send(ct.Opts{Verb: "hold", Data: map[string]any{"reason": "abort"}, TaskID: "tsk-abort-1",
		ParamsHash: "sha256:" + strings.Repeat("ab", 32)})
	if len(got) != 2 || got[0]["outcome"] != "acknowledged" || got[1]["outcome"] != "executed" {
		t.Fatalf("abort hold: %s %v", outcomes(got), got)
	}
	if mav := got[1]["mavlink"].(map[string]any); mav["command"] != "MAV_CMD_DO_REPOSITION" {
		t.Fatalf("abort hold mavlink %v", mav)
	}
	if r := taskReport(); r["state"] != "done_not_ok" || r["error"].(map[string]any)["code"] != "cancelled" {
		t.Fatalf("abort task report %v", r)
	}

	// Plain hold task: the vehicle is holding, the task keeps executing.
	got = h.send(ct.Opts{Verb: "hold", Data: map[string]any{"duration_s": 60}, TaskID: "tsk-hold-1"})
	if len(got) != 2 || got[1]["outcome"] != "executed" {
		t.Fatalf("hold: %s", outcomes(got))
	}
	if r := taskReport(); r["state"] != "executing" || r["error"] != nil {
		t.Fatalf("hold task report %v", r)
	}

	h.veh.mu.Lock()
	defer h.veh.mu.Unlock()
	if len(h.veh.ints) != 2 || len(h.veh.longs) != 0 {
		t.Fatalf("vehicle got %d COMMAND_INT, %d COMMAND_LONG", len(h.veh.ints), len(h.veh.longs))
	}
	for _, g := range h.veh.ints {
		if g.Command != common.MAV_CMD_DO_REPOSITION || g.Frame != common.MAV_FRAME_GLOBAL ||
			g.X != 12345 || g.Y != -6789 || math.Abs(float64(g.Z)-61.5) > 0.01 ||
			g.Param1 != -1 || g.Param2 != float32(common.MAV_DO_REPOSITION_FLAGS_CHANGE_MODE) || !math.IsNaN(float64(g.Param4)) {
			t.Fatalf("hold DO_REPOSITION fields %+v", g)
		}
	}
}

// A vehicle that does not read COMMAND_INT gets the hold as COMMAND_LONG (param5/6/7).
func TestHoldFallsBackToCommandLong(t *testing.T) {
	h := newHarness(t)
	h.veh.mu.Lock()
	h.veh.intUnsup = true
	h.veh.mu.Unlock()
	if got := h.send(ct.Opts{Verb: "hold"}); len(got) != 2 || got[1]["outcome"] != "executed" {
		t.Fatalf("hold: %s", outcomes(got))
	}
	h.veh.mu.Lock()
	defer h.veh.mu.Unlock()
	if len(h.veh.longs) != 1 {
		t.Fatalf("COMMAND_LONG count %d", len(h.veh.longs))
	}
	l := h.veh.longs[0]
	if l.Command != common.MAV_CMD_DO_REPOSITION || math.Abs(float64(l.Param5)-0.001) > 1e-6 ||
		math.Abs(float64(l.Param6)-0.001) > 1e-6 || math.Abs(float64(l.Param7)-50) > 0.01 || l.Param2 != 1 {
		t.Fatalf("COMMAND_LONG hold %+v", l)
	}
}

// hold needs a link and a recent position; a stale position is refused, never sent.
func TestHoldRefusedWithoutFreshPosition(t *testing.T) {
	h := newHarness(t)
	h.veh.mu.Lock()
	h.veh.paused = true
	h.veh.mu.Unlock()
	time.Sleep(100 * time.Millisecond)
	h.b.offset.Store(int64(10 * time.Second))
	got := h.send(ct.Opts{Verb: "hold", Now: time.Now().Add(10 * time.Second)})
	if len(got) != 1 || got[0]["code"] != "precondition_failed" || got[0]["detail"] != "offline" {
		t.Fatalf("offline hold: %s %v", outcomes(got), got)
	}

	s := newVehicleState()
	now := time.Now()
	if _, _, _, ok := s.freshPosition(now, time.Second); ok {
		t.Fatal("no position yet but freshPosition ok")
	}
	s.update(&ardupilotmega.MessageGlobalPositionInt{Lat: 10, Lon: 20, Alt: 30000}, now.Add(-5*time.Second), time.Second)
	if _, _, _, ok := s.freshPosition(now, 3*time.Second); ok {
		t.Fatal("5 s old position accepted with a 3 s limit")
	}
	if lat, lon, alt, ok := s.freshPosition(now, 10*time.Second); !ok || lat != 1e-6 || lon != 2e-6 || alt != 30 {
		t.Fatalf("freshPosition %v %v %v %v", lat, lon, alt, ok)
	}
}
