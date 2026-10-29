package adapt

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/onemindos/omos/services/entity/internal/identity"
	"github.com/onemindos/omos/services/entity/internal/model"
)

// geoEnvelope is omos docs/NATS-DOCTRINE.md §9.2.
type geoEnvelope struct {
	UID      string         `json:"uid"`
	Type     string         `json:"type"`
	Lat      *float64       `json:"lat"`
	Lon      *float64       `json:"lon"`
	HAE      *float64       `json:"hae"`
	Course   *float64       `json:"course"`
	Speed    *float64       `json:"speed"`
	Callsign string         `json:"callsign"`
	Stale    int            `json:"stale"`
	T        int64          `json:"t"`
	Meta     map[string]any `json:"meta"`
	CE       *float64       `json:"ce"`
	LE       *float64       `json:"le"`
}

// fleetClass maps the doctrine §9.1 `tel.fleet.{class}` vocabulary to entity_type.
var fleetClass = map[string]string{
	"drone": "drone", "robot": "robot", "vehicle": "vehicle", "node": "node", "edge": "node",
	"agent": "agent", "aircraft": "aircraft", "vessel": "vessel",
}

func neg1(p *float64) *float64 {
	if p == nil || *p < 0 {
		return nil // doctrine §9.2: -1 is unknown
	}
	return p
}

func metaString(m map[string]any, k string) string {
	switch v := m[k].(type) {
	case string:
		return v
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64)
	}
	return ""
}

// Fleet parses `tel.fleet.{class}.{id}.{position|status}` (doctrine §9.1, BUS-DESIGN §3.1).
func Fleet(data []byte, subject string) (*model.Observation, string, error) {
	tok := strings.Split(subject, ".")
	if len(tok) != 5 {
		return nil, SkipUnsupported, nil
	}
	class, id, metric := tok[2], tok[3], tok[4]
	et, ok := fleetClass[class]
	if !ok {
		return nil, SkipUnsupported, nil
	}
	var g geoEnvelope
	if err := json.Unmarshal(data, &g); err != nil {
		return nil, "", fmt.Errorf("geo envelope: %w", err)
	}
	if g.UID == "" {
		g.UID = id
	}
	if g.T <= 0 {
		return nil, "", fmt.Errorf("geo envelope without t")
	}
	o := &model.Observation{
		Integration: "fleet-" + class, SourceID: g.UID, DataType: "geo", Subject: subject,
		EntityType: et, DocType: "fleet." + class, Callsign: g.Callsign, SourceT: g.T,
		StaleS: g.Stale, Meta: g.Meta, Tier: model.TierSelfReport, Template: "asset",
		Disposition: "friendly", Marking: "AMBER", How: "m-g",
	}
	if node := metaString(g.Meta, "node_id"); node != "" {
		o.Integration = "fleet-" + class + "-" + node
	}
	low := strings.ToLower(g.UID)
	if strings.HasPrefix(low, "sim-") {
		o.Keys = []model.Alias{{Type: identity.KeySim, ID: low}}
		o.Simulated, o.How, o.Marking = true, "m-s", "AMBER+STRICT"
	} else {
		o.Keys = []model.Alias{{Type: identity.KeyFleetUID, ID: g.UID}}
	}
	if sys := metaString(g.Meta, "mavlink_sysid"); sys != "" {
		o.Keys = append(o.Keys, model.Alias{Type: identity.KeyMavlinkSys, ID: sys})
	}
	switch et {
	case "drone":
		o.Kind, o.Environment = metaString(g.Meta, "airframe"), "air"
		if o.Kind == "" {
			o.Kind = metaString(g.Meta, "kind") // omos #550 bridge: multirotor|fixed_wing|vtol
		}
		o.Traits = []string{"locatable", "trackable", "taskable", "observable"}
		if o.Kind == "" || o.Kind == "unknown" {
			o.Kind = "multirotor"
		}
	case "robot", "vehicle":
		o.Kind, o.Environment = metaString(g.Meta, "kind"), "land" // omos #566 bridge: ugv
		switch o.Kind {                                            // robotics-cot-profile.md §1: a boat or submersible is class robot
		case "usv":
			o.Environment = "surface"
		case "uuv":
			o.Environment = "subsurface"
		}
		o.Traits = []string{"locatable", "trackable", "taskable"}
	default:
		o.Traits = []string{"locatable"}
	}
	h := map[string]any{}
	for _, k := range []string{"battery_pct", "link_rssi"} {
		if v, ok := g.Meta[k]; ok {
			h[k] = v
		}
	}
	// A bridge that tracks its link (omos #550 heartbeat, #566 last_heard) says so in
	// meta.connection; the policy service refuses commands unless it reads "online".
	switch conn := metaString(g.Meta, "connection"); conn {
	case "online", "degraded", "offline":
		h["connection"] = conn
	default:
		if len(h) > 0 {
			h["connection"] = "online"
		}
	}
	if len(h) > 0 {
		o.Health = h
	}
	if caps, ok := FleetCapabilities(et, g.UID, o.Traits, g.Meta); ok && o.Simulated {
		o.Capabilities, o.HasCapabilities = caps, true
	}
	if mode := metaString(g.Meta, "mode"); mode != "" {
		o.Status = map[string]any{"mode": mode}
	}
	if metric == "position" {
		if g.Lat == nil || g.Lon == nil {
			return nil, "", fmt.Errorf("position without lat/lon")
		}
		o.HasPos, o.Lat, o.Lon = true, *g.Lat, *g.Lon
		o.HAE, o.Course, o.Speed = g.HAE, neg1(g.Course), neg1(g.Speed)
		if g.CE != nil {
			o.CE = *g.CE
		}
		if g.LE != nil {
			o.LE = *g.LE
		}
		o.Components = []string{"location", "health", "status"}
	} else {
		o.Components = []string{"health", "status"}
	}
	return o, "", nil
}

// worldEventKinds are the non-moving doctrine §9.1 world classes the service turns into
// `event` entities. Military surge subjects are out of civilian scope (BUS-DESIGN D9).
var worldEventKinds = map[string]string{
	"fire": "wildfire", "seismic": "earthquake", "conflict": "conflict", "chokepoint": "chokepoint",
}

// World parses `evt.world.{class}.{id}` (doctrine §9.1, §9.2).
func World(data []byte, subject string) (*model.Observation, string, error) {
	tok := strings.Split(subject, ".")
	if len(tok) < 4 {
		return nil, SkipUnsupported, nil
	}
	class := tok[2]
	if class == "military" || (class == "air" && len(tok) > 3 && tok[3] == "mil") {
		return nil, SkipOutOfScope, nil
	}
	var g geoEnvelope
	if err := json.Unmarshal(data, &g); err != nil {
		return nil, "", fmt.Errorf("geo envelope: %w", err)
	}
	if g.Lat == nil || g.Lon == nil || g.T <= 0 {
		return nil, "", fmt.Errorf("world record without lat/lon/t")
	}
	if g.UID == "" {
		g.UID = tok[len(tok)-1]
	}
	o := &model.Observation{
		Integration: "world-" + class, SourceID: g.UID, DataType: "geo", Subject: subject,
		HasPos: true, Lat: *g.Lat, Lon: *g.Lon, HAE: g.HAE, Course: neg1(g.Course),
		Speed: neg1(g.Speed), Callsign: g.Callsign, SourceT: g.T, StaleS: g.Stale, Meta: g.Meta,
		DocType: g.Type, Marking: "CLEAR", How: "m-g", Tier: model.TierSelfReport,
		Components: []string{"location"},
	}
	if g.CE != nil {
		o.CE = *g.CE
	}
	switch class {
	case "air":
		icao := strings.ToLower(metaString(g.Meta, "icao"))
		if icao == "" {
			icao = strings.TrimPrefix(strings.ToLower(g.UID), "icao-")
		}
		if !reHex6.MatchString(icao) {
			return nil, "", fmt.Errorf("aircraft without a valid ICAO hex (uid %q)", g.UID)
		}
		o.Keys = []model.Alias{{Type: identity.KeyICAO, ID: icao}}
		o.EntityType, o.Template, o.Environment = "aircraft", "asset", "air"
		o.Traits = []string{"locatable", "trackable"}
		if o.DocType == "" {
			o.DocType = "air.unknown"
		}
		if metaString(g.Meta, "squawk") == "7700" {
			o.Emergency = true
		}
	case "sea":
		mmsi := metaString(g.Meta, "mmsi")
		if mmsi == "" {
			mmsi = strings.TrimPrefix(strings.ToLower(g.UID), "mmsi-")
		}
		if !reDigit.MatchString(mmsi) {
			return nil, "", fmt.Errorf("vessel without a valid MMSI (uid %q)", g.UID)
		}
		o.Keys = []model.Alias{{Type: identity.KeyMMSI, ID: mmsi}}
		o.EntityType, o.Template, o.Environment = "vessel", "asset", "surface"
		o.Traits = []string{"locatable", "trackable"}
		if o.DocType == "" {
			o.DocType = "sea.unknown"
		}
	case "rail":
		// MTA trains from the grammar pipeline (#558 A1-P1/A1-P4): uid `mta-{feed}-{trip}`.
		// Keyed as gtfs `{feed}-{trip_id}`, the key the legacy MTA- CoT gives the same train
		// (cot.go), so both paths fuse into one entity during the parallel run.
		feed, trip := "", metaString(g.Meta, "trip_id")
		if rest, ok := strings.CutPrefix(strings.ToLower(g.UID), "mta-"); ok {
			f, tok, _ := strings.Cut(rest, "-")
			feed = f
			if trip == "" {
				trip = tok
			}
		}
		key := strings.ToLower(g.UID)
		if feed != "" && trip != "" {
			key = feed + "-" + trip
		}
		o.Keys = []model.Alias{{Type: identity.KeyGTFS, ID: key}}
		o.EntityType, o.Kind, o.Template, o.Environment = "vehicle", "train", "asset", "land"
		o.Traits = []string{"locatable", "trackable"}
		if o.DocType == "" {
			o.DocType = "ground.rail"
		}
	default:
		kind, ok := worldEventKinds[class]
		if !ok {
			return nil, SkipUnsupported, nil
		}
		o.Keys = []model.Alias{{Type: identity.KeyIntelID, ID: identity.MakeID("evt", class+"-"+g.UID)}}
		o.EntityType, o.Kind, o.Template, o.Tier = "event", kind, "geo", model.TierSensor
		o.Traits = []string{"locatable"}
		o.Course, o.Speed = nil, nil
	}
	return o, "", nil
}

// parseISO accepts the ISO timestamps the intel collectors write.
func parseISO(s string) (int64, bool) {
	for _, layout := range []string{time.RFC3339Nano, "2006-01-02T15:04:05", "2006-01-02 15:04:05"} {
		if t, err := time.Parse(layout, s); err == nil {
			return t.UnixMilli(), true
		}
	}
	return 0, false
}
