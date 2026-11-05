package layer

import (
	"encoding/json"
	"fmt"
	"slices"
)

// Render is the render subset of an entity (layer-snapshot.md §4). Unknown course and speed
// are -1; unknown hae and the expiry of a pinned entity are null.
type Render struct {
	EntityID    string   `json:"entity_id"`
	EntityType  string   `json:"entity_type"`
	Lat         float64  `json:"lat"`
	Lon         float64  `json:"lon"`
	HAE         *float64 `json:"hae"`
	Course      float64  `json:"course"`
	Speed       float64  `json:"speed"`
	Callsign    string   `json:"callsign"`
	CotType     string   `json:"cot_type"`
	Disposition string   `json:"disposition"`
	T           int64    `json:"t"`
	Stale       int      `json:"stale"`
	ExpiryT     *int64   `json:"expiry_t"`
	Seq         uint64   `json:"seq"`
	Live        bool     `json:"live"`
	Origin      string   `json:"origin"`
	Simulated   bool     `json:"simulated"`
}

// entIn is the part of an entity.v1 record (entity.md §4.1) the layer service reads.
type entIn struct {
	EntityID    string   `json:"entity_id"`
	EntityType  string   `json:"entity_type"`
	Kind        string   `json:"kind"`
	Lat         float64  `json:"lat"`
	Lon         float64  `json:"lon"`
	HAE         *float64 `json:"hae"`
	Course      *float64 `json:"course"`
	Speed       *float64 `json:"speed"`
	Callsign    string   `json:"callsign"`
	CotType     string   `json:"cot_type"`
	Disposition string   `json:"disposition"`
	T           int64    `json:"t"`
	Stale       int      `json:"stale"`
	ExpiryT     *int64   `json:"expiry_t"`
	Seq         uint64   `json:"seq"`
	Live        bool     `json:"live"`
	Lifecycle   string   `json:"lifecycle"`
	Origin      string   `json:"origin"`
	Location    *struct {
		Lat       float64  `json:"lat"`
		Lon       float64  `json:"lon"`
		HAEM      *float64 `json:"hae_m"`
		CourseDeg *float64 `json:"course_deg"`
		SpeedMps  *float64 `json:"speed_mps"`
	} `json:"location"`
	Indicators struct {
		Simulated bool `json:"simulated"`
	} `json:"indicators"`
	Marking struct {
		TLP          string   `json:"tlp"`
		ReleasableTo []string `json:"releasable_to"`
	} `json:"marking"`
	Sources []struct {
		DataType string `json:"data_type"`
		Subject  string `json:"subject"`
	} `json:"sources"`
	Geo *GeoIn `json:"geo"`
}

// GeoIn is the `geo` component (entity.md: shape, GeoJSON geometry in coords, radius_m).
type GeoIn struct {
	Shape   string          `json:"shape"`
	Coords  json.RawMessage `json:"coords"`
	RadiusM float64         `json:"radius_m"`
}

// entry is one entity held in memory.
type entry struct {
	key        string // "{type}.{id}"
	r          Render
	tlp        string
	releasable []string
	fromCot    bool
	kind       string          // entity kind (layer kind rules)
	subjects   []string        // distinct source subjects (layer source-subject rules)
	geo        *GeoIn          // zones only (AOI by zone_id)
	raw        json.RawMessage // full entity (fields: "full")
	renderJSON []byte
}

// parseEntity turns one ENT message body into an entry. defOrigin fills `origin` when the
// record has none (entity.v1 has no origin field yet; the owning site is the writer's site).
func parseEntity(subjType, subjID string, data []byte, defOrigin string) (*entry, *entIn, error) {
	var in entIn
	if err := json.Unmarshal(data, &in); err != nil {
		return nil, nil, fmt.Errorf("entity json: %w", err)
	}
	if in.EntityID == "" {
		in.EntityID = subjID
	}
	if in.EntityType == "" {
		in.EntityType = subjType
	}
	if in.EntityID != subjID || in.EntityType != subjType {
		return nil, nil, fmt.Errorf("subject ent.%s.%s.state carries %s/%s", subjType, subjID, in.EntityType, in.EntityID)
	}
	r := Render{
		EntityID: in.EntityID, EntityType: in.EntityType, Lat: in.Lat, Lon: in.Lon, HAE: in.HAE,
		Course: -1, Speed: -1, Callsign: in.Callsign, CotType: in.CotType,
		Disposition: in.Disposition, T: in.T, Stale: in.Stale, ExpiryT: in.ExpiryT, Seq: in.Seq,
		Live: in.Live, Origin: in.Origin, Simulated: in.Indicators.Simulated,
	}
	if in.Course != nil {
		r.Course = *in.Course
	}
	if in.Speed != nil {
		r.Speed = *in.Speed
	}
	if l := in.Location; l != nil {
		if r.Lat == 0 && r.Lon == 0 {
			r.Lat, r.Lon = l.Lat, l.Lon
		}
		if r.HAE == nil {
			r.HAE = l.HAEM
		}
		if in.Course == nil && l.CourseDeg != nil {
			r.Course = *l.CourseDeg
		}
		if in.Speed == nil && l.SpeedMps != nil {
			r.Speed = *l.SpeedMps
		}
	}
	if r.Disposition == "" {
		r.Disposition = "unknown"
	}
	if r.Origin == "" {
		r.Origin = defOrigin
	}
	e := &entry{key: subjType + "." + subjID, r: r, tlp: in.Marking.TLP, releasable: in.Marking.ReleasableTo, raw: append(json.RawMessage(nil), data...)}
	if e.tlp == "" {
		e.tlp = "AMBER" // entity service default (engine.go) when no source sets a marking
	}
	e.kind = in.Kind
	for _, s := range in.Sources {
		if s.DataType == "cot" {
			e.fromCot = true
		}
		if s.Subject != "" && !slices.Contains(e.subjects, s.Subject) {
			e.subjects = append(e.subjects, s.Subject)
		}
	}
	if in.EntityType == "zone" {
		e.geo = in.Geo
	}
	e.renderJSON, _ = json.Marshal(r)
	return e, &in, nil
}

// removalReason maps a non-live record to a delta remove reason (layer-delta.md §2.2).
func removalReason(in *entIn) string {
	switch in.Lifecycle {
	case "deleted":
		return "deleted"
	case "expired":
		return "expired"
	}
	return "not_live"
}

// markerReason maps Nats-Marker-Reason (layer-delta.md §3) to a remove reason.
func markerReason(h string) string {
	if h == "MaxAge" {
		return "expired"
	}
	return "deleted"
}
