package adapt

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/onemindos/omos/services/entity/internal/identity"
	"github.com/onemindos/omos/services/entity/internal/model"
)

// intelRecord is the payload the live feed-intel collectors publish on `intel.>`
// (omos platform/omos/intel-scripts.yaml build_payload: id, type, source, lat, lon, timestamp,
// expires, properties).
type intelRecord struct {
	ID         string         `json:"id"`
	Type       string         `json:"type"`
	Source     string         `json:"source"`
	Lat        *float64       `json:"lat"`
	Lon        *float64       `json:"lon"`
	Timestamp  string         `json:"timestamp"`
	Expires    string         `json:"expires"`
	Properties map[string]any `json:"properties"`
}

// Intel parses an `intel.{domain}.{source}.…` point event into an `event` observation.
// Records without a position are skipped (they are not map entities).
func Intel(data []byte, subject string) (*model.Observation, string, error) {
	var r intelRecord
	if err := json.Unmarshal(data, &r); err != nil {
		return nil, "", fmt.Errorf("intel json: %w", err)
	}
	if r.Lat == nil || r.Lon == nil || (*r.Lat == 0 && *r.Lon == 0) {
		return nil, SkipNoPosition, nil
	}
	if r.ID == "" {
		return nil, "", fmt.Errorf("intel record without id")
	}
	tok := strings.Split(subject, ".")
	src := strings.ToLower(r.Source)
	if src == "" && len(tok) > 2 {
		src = tok[2]
	}
	srcT, ok := parseISO(r.Timestamp)
	if !ok {
		return nil, "", fmt.Errorf("intel timestamp %q", r.Timestamp)
	}
	native := strings.TrimPrefix(strings.ToLower(r.ID), src+"-")
	id := identity.MakeID(src, native)
	o := &model.Observation{
		Integration: "feed-intel-" + src, SourceID: r.ID, DataType: "intel", Subject: subject,
		Keys:       []model.Alias{{Type: identity.KeyIntelID, ID: id}},
		EntityType: "event", Kind: r.Type, DocType: "onemind.event", Template: "geo",
		Traits: []string{"locatable"}, Tier: model.TierSensor, How: "m-g",
		HasPos: true, Lat: *r.Lat, Lon: *r.Lon, SourceT: srcT, Marking: "CLEAR",
		Meta: r.Properties, Components: []string{"location", "geo"},
	}
	if p, ok := r.Properties["place"].(string); ok {
		o.Callsign = p
	}
	if exp, ok := parseISO(r.Expires); ok && exp > srcT {
		o.ExpireT = exp
	}
	return o, "", nil
}

type detRecord struct {
	Lat        *float64 `json:"lat"`
	Lon        *float64 `json:"lon"`
	HAE        *float64 `json:"hae"`
	T          int64    `json:"t"`
	Class      string   `json:"class"`
	Confidence float64  `json:"confidence"`
	SensorID   string   `json:"sensor_id"`
	CE         *float64 `json:"ce"`
}

// Det parses `det.{domain}.{type}.{id}` where {id} is the detector's own track id. Keyless
// detections that need association (Stage 3 GNN) are not handled in this slice. A detected
// person is always a `track` of kind person, never a `person` entity (entity.md §5 "Person
// safety"; ONTOLOGY §3.3).
func Det(data []byte, subject string) (*model.Observation, string, error) {
	tok := strings.Split(subject, ".")
	if len(tok) != 4 {
		return nil, SkipUnsupported, nil
	}
	var d detRecord
	if err := json.Unmarshal(data, &d); err != nil {
		return nil, "", fmt.Errorf("det json: %w", err)
	}
	if d.Lat == nil || d.Lon == nil || d.T <= 0 {
		return nil, "", fmt.Errorf("detection without lat/lon/t")
	}
	sensor := d.SensorID
	if sensor == "" {
		sensor = tok[2]
	}
	class := strings.ToLower(d.Class)
	o := &model.Observation{
		Integration: "det-" + tok[1] + "-" + tok[2], SourceID: sensor + ":" + tok[3], DataType: "det",
		Subject: subject, Keys: []model.Alias{{Type: identity.KeyDetTrack, ID: sensor + ":" + tok[1] + "." + tok[2] + "." + tok[3]}},
		EntityType: "track", Kind: class, DocType: "onemind.track", Template: "track",
		Traits: []string{"locatable", "trackable"}, Tier: model.TierSensor, How: "m-f",
		HasPos: true, Lat: *d.Lat, Lon: *d.Lon, HAE: d.HAE, SourceT: d.T, Marking: "GREEN",
		Disposition: "unknown", Components: []string{"location", "tracked"},
	}
	if d.CE != nil {
		o.CE = *d.CE
	}
	if class == "person" || class == "human" {
		f := false
		o.Kind, o.Marking, o.EgressTAK, o.StaleS = "person", "AMBER+STRICT", &f, 30
	}
	return o, "", nil
}
