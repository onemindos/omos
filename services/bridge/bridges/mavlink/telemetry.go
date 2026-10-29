package main

import (
	"encoding/json"
	"fmt"
	"math"
	"strings"
	"sync"
	"time"

	"github.com/bluenviron/gomavlib/v3/pkg/dialects/ardupilotmega"
	"github.com/bluenviron/gomavlib/v3/pkg/message"

	"github.com/onemindos/omos-bridge/bridges/mavlink/internal/geoid"
)

// GeoEnvelope is the doctrine §9.2 geo-envelope (docs/NATS-DOCTRINE.md), carried by
// tel.fleet.drone.{id}.position and tel.fleet.drone.{id}.status.
type GeoEnvelope struct {
	UID      string         `json:"uid"`
	Type     string         `json:"type"`
	Lat      float64        `json:"lat"`
	Lon      float64        `json:"lon"`
	Hae      float64        `json:"hae"`
	Course   float64        `json:"course"`
	Speed    float64        `json:"speed"`
	Callsign string         `json:"callsign"`
	Stale    int            `json:"stale"`
	T        int64          `json:"t"`
	Meta     map[string]any `json:"meta"`
}

// vehicleState is the latest telemetry of one vehicle.
type vehicleState struct {
	mu sync.Mutex

	hasPos       bool
	lat, lon     float64
	altMSL       float64
	relAlt       float64
	vn, ve, vd   float64 // m/s
	hdgDeg       float64 // -1 unknown
	groundspeed  float64
	hasVfr       bool
	posAt        time.Time
	lastPosPub   time.Time
	gpsN         float64 // GPS_RAW_INT alt_ellipsoid - alt, metres
	gpsNAt       time.Time
	batteryPct   int // -1 unknown
	voltageV     float64
	hbAt         time.Time
	armed        bool
	baseMode     uint8
	customMode   uint32
	mavType      ardupilotmega.MAV_TYPE
	autopilot    ardupilotmega.MAV_AUTOPILOT
	systemStatus ardupilotmega.MAV_STATE
	mavVersion   uint8
	hbSeen       bool
}

func newVehicleState() *vehicleState { return &vehicleState{hdgDeg: -1, batteryPct: -1} }

// update folds one MAVLink message into the state. It reports whether a position
// message is due (GLOBAL_POSITION_INT received and the publish interval elapsed).
func (s *vehicleState) update(m message.Message, now time.Time, posInterval time.Duration) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	switch msg := m.(type) {
	case *ardupilotmega.MessageHeartbeat:
		s.hbSeen = true
		s.hbAt = now
		s.baseMode = uint8(msg.BaseMode)
		s.armed = msg.BaseMode&ardupilotmega.MAV_MODE_FLAG_SAFETY_ARMED != 0
		s.customMode = msg.CustomMode
		s.mavType = msg.Type
		s.autopilot = msg.Autopilot
		s.systemStatus = msg.SystemStatus
		s.mavVersion = msg.MavlinkVersion
	case *ardupilotmega.MessageSysStatus:
		if msg.BatteryRemaining >= 0 {
			s.batteryPct = int(msg.BatteryRemaining)
		}
		if msg.VoltageBattery != math.MaxUint16 {
			s.voltageV = float64(msg.VoltageBattery) / 1000
		}
	case *ardupilotmega.MessageBatteryStatus:
		if msg.BatteryRemaining >= 0 {
			s.batteryPct = int(msg.BatteryRemaining)
		}
	case *ardupilotmega.MessageVfrHud:
		s.groundspeed = float64(msg.Groundspeed)
		s.hasVfr = true
	case *ardupilotmega.MessageGpsRawInt:
		// alt and alt_ellipsoid are both mm; alt_ellipsoid is a MAVLink 2 extension
		// and 0 when the autopilot does not fill it.
		if msg.AltEllipsoid != 0 && msg.FixType >= ardupilotmega.GPS_FIX_TYPE_3D_FIX {
			s.gpsN = float64(msg.AltEllipsoid-msg.Alt) / 1000
			s.gpsNAt = now
		}
	case *ardupilotmega.MessageGlobalPositionInt:
		s.hasPos = true
		s.posAt = now
		s.lat = float64(msg.Lat) / 1e7
		s.lon = float64(msg.Lon) / 1e7
		s.altMSL = float64(msg.Alt) / 1000
		s.relAlt = float64(msg.RelativeAlt) / 1000
		s.vn, s.ve, s.vd = float64(msg.Vx)/100, float64(msg.Vy)/100, float64(msg.Vz)/100
		if msg.Hdg == math.MaxUint16 {
			s.hdgDeg = -1
		} else {
			s.hdgDeg = float64(msg.Hdg) / 100
		}
		if now.Sub(s.lastPosPub) >= posInterval {
			s.lastPosPub = now
			return true
		}
	}
	return false
}

func (s *vehicleState) online(now time.Time, timeout time.Duration) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.hbSeen && now.Sub(s.hbAt) <= timeout
}

func (s *vehicleState) position() (lat, lon, altMSL float64, ok bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.lat, s.lon, s.altMSL, s.hasPos
}

// haeLocked returns the ellipsoid height and how it was derived.
// Preferred: the autopilot's own GPS_RAW_INT undulation (alt_ellipsoid - alt) seen in
// the last 5 s; otherwise the embedded EGM96 model.
func (s *vehicleState) haeLocked(now time.Time) (float64, string, float64) {
	if !s.gpsNAt.IsZero() && now.Sub(s.gpsNAt) < 5*time.Second {
		return s.altMSL + s.gpsN, "gps_raw_int", s.gpsN
	}
	n, err := geoid.Undulation(s.lat, s.lon)
	if err != nil {
		return s.altMSL, "msl_uncorrected", 0
	}
	return s.altMSL + n, "egm96", n
}

func (s *vehicleState) modeLocked() string {
	if s.baseMode&uint8(ardupilotmega.MAV_MODE_FLAG_CUSTOM_MODE_ENABLED) != 0 &&
		s.autopilot == ardupilotmega.MAV_AUTOPILOT_ARDUPILOTMEGA && kindOf(s.mavType) == "multirotor" {
		return strings.TrimPrefix(ardupilotmega.COPTER_MODE(s.customMode).String(), "COPTER_MODE_")
	}
	return fmt.Sprintf("custom:%d", s.customMode)
}

// kindOf maps HEARTBEAT.type to the entity kind (robotics-cot-profile.md §2).
func kindOf(t ardupilotmega.MAV_TYPE) string {
	switch t {
	case ardupilotmega.MAV_TYPE_QUADROTOR, ardupilotmega.MAV_TYPE_HEXAROTOR, ardupilotmega.MAV_TYPE_OCTOROTOR,
		ardupilotmega.MAV_TYPE_TRICOPTER, ardupilotmega.MAV_TYPE_COAXIAL, ardupilotmega.MAV_TYPE_HELICOPTER,
		ardupilotmega.MAV_TYPE_DODECAROTOR:
		return "multirotor"
	case ardupilotmega.MAV_TYPE_FIXED_WING:
		return "fixed_wing"
	case ardupilotmega.MAV_TYPE_VTOL_TILTROTOR, ardupilotmega.MAV_TYPE_VTOL_FIXEDROTOR,
		ardupilotmega.MAV_TYPE_VTOL_TAILSITTER, ardupilotmega.MAV_TYPE_VTOL_TILTWING,
		ardupilotmega.MAV_TYPE_VTOL_TAILSITTER_DUOROTOR, ardupilotmega.MAV_TYPE_VTOL_TAILSITTER_QUADROTOR:
		return "vtol"
	}
	return ""
}

func round(v float64, places int) float64 {
	p := math.Pow(10, float64(places))
	return math.Round(v*p)/p + 0 // + 0 turns -0 into 0
}

// positionEnvelope builds tel.fleet.drone.{id}.position.
func (s *vehicleState) positionEnvelope(cfg *Config, now time.Time) GeoEnvelope {
	s.mu.Lock()
	defer s.mu.Unlock()
	hae, src, n := s.haeLocked(now)
	speed := math.Hypot(s.vn, s.ve)
	course := -1.0
	if speed >= 0.5 { // course over ground from the velocity vector; -1 when stationary
		course = math.Mod(math.Atan2(s.ve, s.vn)*180/math.Pi+360, 360)
	}
	return GeoEnvelope{
		UID: cfg.AssetID, Type: "fleet.drone", Lat: s.lat, Lon: s.lon, Hae: round(hae, 2),
		Course: round(course, 1), Speed: round(speed, 2), Callsign: cfg.Callsign, Stale: cfg.StaleS,
		T: now.UnixMilli(),
		Meta: map[string]any{
			"alt_msl_m":          round(s.altMSL, 2),
			"relative_alt_m":     round(s.relAlt, 2),
			"geoid_undulation_m": round(n, 2),
			"hae_source":         src,
			"heading_deg":        round(s.hdgDeg, 1),
			"climb_mps":          round(-s.vd, 2),
			"battery_pct":        s.batteryPct,
			"mode":               s.modeLocked(),
			"armed":              s.armed,
			"mav_sysid":          int(cfg.SysID),
			"simulated":          true,
			"source":             cfg.Source,
		},
	}
}

// statusEnvelope builds tel.fleet.drone.{id}.status (battery, mode, link health).
func (s *vehicleState) statusEnvelope(cfg *Config, now time.Time) (GeoEnvelope, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if !s.hasPos {
		return GeoEnvelope{}, false // §9.2 requires lat/lon
	}
	hae, _, _ := s.haeLocked(now)
	conn := "offline"
	if s.hbSeen && now.Sub(s.hbAt) <= cfg.LinkTimeout {
		conn = "online"
	}
	lastHB := int64(-1)
	if s.hbSeen {
		lastHB = s.hbAt.UnixMilli()
	}
	speed := s.groundspeed
	if !s.hasVfr {
		speed = math.Hypot(s.vn, s.ve)
	}
	return GeoEnvelope{
		UID: cfg.AssetID, Type: "fleet.drone", Lat: s.lat, Lon: s.lon, Hae: round(hae, 2),
		Course: -1, Speed: round(speed, 2), Callsign: cfg.Callsign, Stale: cfg.StaleS, T: now.UnixMilli(),
		Meta: map[string]any{
			"battery_pct":       s.batteryPct,
			"voltage_v":         round(s.voltageV, 2),
			"mode":              s.modeLocked(),
			"armed":             s.armed,
			"connection":        conn,
			"link_rssi":         -1, // not mapped: RADIO_STATUS rssi units are radio-specific
			"last_heartbeat_t":  lastHB,
			"system_status":     strings.TrimPrefix(s.systemStatus.String(), "MAV_STATE_"),
			"kind":              kindOf(s.mavType),
			"mav_type":          strings.TrimPrefix(s.mavType.String(), "MAV_TYPE_"),
			"autopilot":         strings.TrimPrefix(s.autopilot.String(), "MAV_AUTOPILOT_"),
			"mavlink_version":   int(s.mavVersion),
			"mav_sysid":         int(cfg.SysID),
			"simulated":         true,
			"source":            cfg.Source,
			"position_age_ms":   now.Sub(s.posAt).Milliseconds(),
			"commands_accepted": "goto,rtl,land,arm,disarm,estop,takeoff",
		},
	}, true
}

func mustJSON(v any) []byte {
	b, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return b
}
