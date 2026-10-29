package adapt

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/onemindos/omos/services/entity/internal/identity"
	"github.com/onemindos/omos/services/entity/internal/model"
	"github.com/onemindos/omos/services/entity/internal/orbit"
)

// intelRecord is the flat payload the live feed-intel collectors publish on `intel.>`
// (omos platform/omos/intel-scripts.yaml build_payload: id, type, source, lat, lon, timestamp,
// expires, properties). The interim Redpanda CelesTrak pipeline (#558) uses it too.
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

// ingestEnvelope is the omos bus envelope omos-ingest publishes (#553 record.go NewEnvelope;
// services/schemas/envelope/envelope.schema.json): {v, id, ts, type: <subject>, source:
// "ingest.<source id>", subject, data, fetchedAt, msgId}. The record is in data.
type ingestEnvelope struct {
	V      int             `json:"v"`
	ID     string          `json:"id"`
	TS     string          `json:"ts"`
	Type   string          `json:"type"`
	Source string          `json:"source"`
	Data   json.RawMessage `json:"data"`
	MsgID  string          `json:"msgId"`
}

// TLESubject carries one CelesTrak element set per message (BUS-DESIGN §3.1).
const TLESubject = "intel.space.celestrak.tle"

// Now is the clock TLE propagation uses; tests replace it.
var Now = time.Now

// Intel parses an `intel.{domain}.{source}.…` point event into an `event` observation. It
// accepts the flat feed-intel record and the omos-ingest envelope. TLE records become
// satellite observations propagated with SGP4. Records without a position are skipped.
func Intel(data []byte, subject string) (*model.Observation, string, error) {
	r, err := decodeIntel(data, subject)
	if err != nil {
		return nil, "", err
	}
	tok := strings.Split(subject, ".")
	src := strings.ToLower(r.Source)
	if src == "" && len(tok) > 2 {
		src = tok[2]
	}
	src = canonSource(src)
	if subject == TLESubject || strings.EqualFold(r.Type, "tle") {
		return tleObservation(r.Properties, subject, src, Now())
	}
	if r.Lat == nil || r.Lon == nil || (*r.Lat == 0 && *r.Lon == 0) {
		return nil, SkipNoPosition, nil
	}
	if *r.Lat < -90 || *r.Lat > 90 || *r.Lon < -180 || *r.Lon > 180 {
		return nil, "", fmt.Errorf("intel position out of range (%v, %v)", *r.Lat, *r.Lon)
	}
	if r.ID == "" {
		return nil, "", fmt.Errorf("intel record without id")
	}
	srcT, ok := parseISO(r.Timestamp)
	if !ok {
		return nil, "", fmt.Errorf("intel timestamp %q", r.Timestamp)
	}
	native := strings.TrimPrefix(nativeID(r.ID), src+"-")
	id := identity.MakeID(src, native)
	o := &model.Observation{
		Integration: "feed-intel-" + src, SourceID: r.ID, DataType: "intel", Subject: subject,
		Keys:       []model.Alias{{Type: identity.KeyIntelID, ID: id}},
		EntityType: "event", Kind: r.Type, DocType: "onemind.event", Template: "geo",
		Traits: []string{"locatable"}, Tier: model.TierSensor, How: "m-g",
		HasPos: true, Lat: *r.Lat, Lon: *r.Lon, SourceT: srcT, Marking: "CLEAR",
		Meta: r.Properties, Components: []string{"location", "geo"},
	}
	if p, ok := r.Properties["place"].(string); ok && p != "" {
		o.Callsign = p
	} else if ev, ok := r.Properties["event"].(string); ok {
		o.Callsign = ev // NWS alerts: "Flood Warning"
	}
	if exp, ok := parseISO(r.Expires); ok && exp > srcT {
		o.ExpireT = exp
	}
	return o, "", nil
}

// sourceAliases maps legacy feed-intel `source` values to the grammar source token omos-ingest
// uses, so both input forms give the same entity id prefix and integration (FEATURES A4).
var sourceAliases = map[string]string{
	"nasa_firms": "firms", "nasa-firms": "firms",
	"nasa-eonet": "eonet", "nasa_eonet": "eonet",
}

func canonSource(src string) string {
	if a, ok := sourceAliases[src]; ok {
		return a
	}
	return src
}

// nativeID lowercases a record id and reduces a URL id to its last path segment. The legacy
// NWS collector sends feature.id ("https://api.weather.gov/alerts/urn:oid:...") where
// omos-ingest sends properties.id ("urn:oid:..."); both then name the same alert.
func nativeID(id string) string {
	n := strings.ToLower(strings.TrimSpace(id))
	for _, scheme := range []string{"https://", "http://"} {
		if !strings.HasPrefix(n, scheme) {
			continue
		}
		path := strings.TrimRight(n[len(scheme):], "/")
		if i := strings.LastIndex(path, "/"); i >= 0 {
			return path[i+1:] // host/.../segment: keep the last segment
		}
	}
	return n
}

// eventTimeMs reads a numeric epoch-millisecond event time (USGS `time`). Values outside
// 2001..2286 are not milliseconds and are ignored.
func eventTimeMs(v any) (int64, bool) {
	n, ok := v.(json.Number)
	if !ok {
		return 0, false
	}
	ms, err := n.Int64()
	if err != nil || ms < 1e12 || ms >= 1e13 {
		return 0, false
	}
	return ms, true
}

// envelopeDropKeys are data fields the envelope form maps to observation fields, or that are
// too large for entity meta (the alert polygon; the entity carries a point geo for now).
var envelopeDropKeys = map[string]bool{"id": true, "lat": true, "lon": true, "type": true, "geometry": true}

// decodeIntel normalises either input form to an intelRecord.
func decodeIntel(data []byte, subject string) (intelRecord, error) {
	var probe struct {
		V    *int            `json:"v"`
		Data json.RawMessage `json:"data"`
	}
	if err := json.Unmarshal(data, &probe); err != nil {
		return intelRecord{}, fmt.Errorf("intel json: %w", err)
	}
	if probe.V == nil || len(probe.Data) == 0 || probe.Data[0] != '{' {
		var r intelRecord
		if err := json.Unmarshal(data, &r); err != nil {
			return intelRecord{}, fmt.Errorf("intel json: %w", err)
		}
		return r, nil
	}
	var env ingestEnvelope
	if err := json.Unmarshal(data, &env); err != nil {
		return intelRecord{}, fmt.Errorf("intel envelope: %w", err)
	}
	if env.V != 1 {
		return intelRecord{}, fmt.Errorf("intel envelope version %d", env.V)
	}
	if env.Type != "" && env.Type != subject {
		return intelRecord{}, fmt.Errorf("intel envelope type %q does not match subject %q", env.Type, subject)
	}
	var d map[string]any
	dec := json.NewDecoder(bytes.NewReader(env.Data))
	dec.UseNumber()
	if err := dec.Decode(&d); err != nil {
		return intelRecord{}, fmt.Errorf("intel envelope data: %w", err)
	}
	r := intelRecord{Timestamp: env.TS, Properties: map[string]any{}}
	// The grammar source token, not the ingest source id: the same feed keeps the same
	// entity ids when it moves from feed-intel to omos-ingest.
	if tok := strings.Split(subject, "."); len(tok) > 2 {
		r.Source = tok[2]
	}
	if s, ok := d["id"].(string); ok && s != "" {
		r.ID = s
	} else {
		r.ID = env.ID
	}
	if s, ok := d["type"].(string); ok && s != "" {
		r.Type = s
	} else if tok := strings.Split(subject, "."); len(tok) > 3 {
		r.Type = tok[3]
	}
	r.Lat, r.Lon = jsonFloat(d["lat"]), jsonFloat(d["lon"])
	// Prefer the record's own event time: omos-ingest falls back to fetchedAt for `ts` when
	// its timeField does not parse (USGS epoch ms), and the entity's source_t is the event.
	if ms, ok := eventTimeMs(d["time"]); ok {
		r.Timestamp = time.UnixMilli(ms).UTC().Format(time.RFC3339Nano)
	}
	if s, ok := d["expires"].(string); ok {
		r.Expires = s
	}
	for k, v := range d {
		if envelopeDropKeys[k] {
			continue
		}
		if n, ok := v.(json.Number); ok {
			if f, err := n.Float64(); err == nil {
				v = f
			}
		}
		r.Properties[k] = v
	}
	if r.Timestamp == "" {
		return intelRecord{}, errors.New("intel envelope without ts")
	}
	return r, nil
}

func jsonFloat(v any) *float64 {
	n, ok := v.(json.Number)
	if !ok {
		return nil
	}
	f, err := n.Float64()
	if err != nil {
		return nil
	}
	return &f
}

// tleObservation propagates a TLE record (properties line1, line2, norad_id, name) to now.
// The satellite is a `vehicle` of kind `satellite` in space, keyed on its NORAD id: the same
// entity the legacy feed-satellites CoT makes (cot.go, ONT-P8).
func tleObservation(props map[string]any, subject, src string, now time.Time) (*model.Observation, string, error) {
	l1, _ := props["line1"].(string)
	l2, _ := props["line2"].(string)
	el, err := orbit.Parse(l1, l2)
	if err != nil {
		return nil, "", err
	}
	if n, ok := props["norad_id"].(string); ok && strings.TrimSpace(n) != "" && normNORAD(n) != normNORAD(el.NORAD) {
		return nil, "", fmt.Errorf("tle norad_id %q does not match the element set (%s)", n, el.NORAD)
	}
	if src == "" {
		src = "celestrak"
	}
	name, _ := props["name"].(string)
	name = strings.TrimSpace(name)
	if name == "" {
		name = "NORAD " + el.NORAD
	}
	meta := map[string]any{}
	for _, k := range []string{"norad_id", "name", "group", "classification", "intl_designator", "epoch", "epoch_ms"} {
		if v, ok := props[k]; ok {
			meta[k] = v
		}
	}
	meta["norad_id"], meta["name"], meta["line1"], meta["line2"] = el.NORAD, name, el.Line1, el.Line2
	tmpl := &model.Observation{
		Integration: "feed-intel-" + src, SourceID: el.NORAD, DataType: "intel", Subject: subject,
		Keys:       []model.Alias{{Type: identity.KeyNORAD, ID: normNORAD(el.NORAD)}},
		EntityType: "vehicle", Kind: "satellite", DocType: "onemind.vehicle", Template: "asset",
		Environment: "space", Traits: []string{"locatable", "trackable"}, Tier: model.TierSensor,
		How: "m-p", Callsign: name, Marking: "CLEAR", StaleS: SatelliteStaleS, Meta: meta,
		Components: []string{"location"}, Orbit: el,
	}
	o, err := Propagate(tmpl, now)
	if err != nil {
		return nil, "", err
	}
	return o, "", nil
}

// normNORAD drops leading zeros so "05" and "5" name the same object.
func normNORAD(n string) string {
	if t := strings.TrimLeft(strings.TrimSpace(n), "0"); t != "" {
		return t
	}
	return "0"
}

// SatelliteStaleS is how long one propagated position stays fresh. The bus layer
// re-propagates more often than this (Config.OrbitEvery, 30 s like the legacy feed).
const SatelliteStaleS = 120

// Propagate returns a copy of a TLE observation positioned at now with SGP4. `m-p` is the
// CoT how for a predicted (machine-computed) position.
func Propagate(tmpl *model.Observation, now time.Time) (*model.Observation, error) {
	if tmpl.Orbit == nil {
		return nil, errors.New("observation has no orbit")
	}
	s, err := tmpl.Orbit.At(now)
	if err != nil {
		return nil, err
	}
	o := *tmpl
	alt, sp, crs := s.AltM, s.SpeedMps, s.CourseDeg
	o.HasPos, o.Lat, o.Lon, o.HAE, o.Speed, o.Course = true, s.Lat, s.Lon, &alt, &sp, &crs
	o.SourceT = s.T.UnixMilli()
	o.RecvT, o.Stream, o.StreamSeq = 0, "", 0
	return &o, nil
}

// TLEFromMeta rebuilds a TLE observation template from a restored satellite entity's meta
// (line1, line2, name), so propagation resumes after a restart without waiting for CelesTrak.
func TLEFromMeta(meta map[string]any, now time.Time) (*model.Observation, error) {
	o, _, err := tleObservation(meta, TLESubject, "celestrak", now)
	return o, err
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
