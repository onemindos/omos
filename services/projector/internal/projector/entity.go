// Package projector is the omos TAK projector (FEATURES A6): it turns fused entities from
// `ent.{entity_type}.{entity_id}.state` into Cursor-on-Target for TAK, only when every gate of
// onemind-xr docs/contracts/cot-projection.md §3 says yes, and sends a removal when an entity
// that was on TAK stops passing them.
//
// The entity record is entity.v1 (services/schemas/entity/entity.v1.schema.json; the writer is
// services/entity). This package decodes only the fields the projector reads.
package projector

import (
	"encoding/json"
	"math"
	"strings"
)

// Entity is the subset of entity.v1 the projector reads (field names are the contract's).
type Entity struct {
	EntityID    string          `json:"entity_id"`
	EntityType  string          `json:"entity_type"`
	Kind        string          `json:"kind"`
	Template    string          `json:"template"`
	Traits      []string        `json:"traits"`
	Live        bool            `json:"live"`
	Lifecycle   string          `json:"lifecycle"`
	Lat         *float64        `json:"lat"`
	Lon         *float64        `json:"lon"`
	HAE         *float64        `json:"hae"`
	Course      *float64        `json:"course"`
	Speed       *float64        `json:"speed"`
	Callsign    string          `json:"callsign"`
	Stale       int             `json:"stale"`
	T           int64           `json:"t"`
	Seq         uint64          `json:"seq"`
	Confidence  float64         `json:"confidence"`
	Disposition string          `json:"disposition"`
	Environment string          `json:"environment"`
	CotType     string          `json:"cot_type"`
	Aliases     struct {
		Name *string `json:"name"`
	} `json:"aliases"`
	Location *struct {
		Lat       float64  `json:"lat"`
		Lon       float64  `json:"lon"`
		HAEM      *float64 `json:"hae_m"`
		CourseDeg *float64 `json:"course_deg"`
		SpeedMps  *float64 `json:"speed_mps"`
	} `json:"location"`
	Uncertainty *struct {
		CEM float64 `json:"ce_m"`
		LEM float64 `json:"le_m"`
	} `json:"uncertainty"`
	Geo         *Geo           `json:"geo"`
	Health      map[string]any `json:"health"`
	Meta        map[string]any `json:"meta"`
	CurrentTask map[string]any `json:"current_task"`
	Indicators  struct {
		Simulated bool `json:"simulated"`
		Exercise  bool `json:"exercise"`
		Starred   bool `json:"starred"`
		EgressTAK bool `json:"egress_tak"`
	} `json:"indicators"`
	Marking struct {
		TLP          string   `json:"tlp"`
		ReleasableTo []string `json:"releasable_to"`
	} `json:"marking"`
	Provenance struct {
		How string `json:"how"`
	} `json:"provenance"`
	Sources []struct {
		Integration string `json:"integration"`
		Subject     string `json:"subject"`
		LastSeenT   int64  `json:"last_seen_t"`
	} `json:"sources"`
}

// Geo is the `geo` component (entity.md §4.1): GeoJSON geometry in coords.
type Geo struct {
	Shape   string          `json:"shape"`
	Coords  json.RawMessage `json:"coords"`
	RadiusM float64         `json:"radius_m"`
	Kind    string          `json:"kind"`
	AltMinM *float64        `json:"alt_min_m"`
	AltMaxM *float64        `json:"alt_max_m"`
}

// Key is "{entity_type}.{entity_id}", the uid suffix and the map key.
func (e *Entity) Key() string { return e.EntityType + "." + e.EntityID }

// UID is the CoT uid (cot-projection.md §5.1).
func (e *Entity) UID() string { return "OM." + e.Key() }

func validLL(lat, lon float64) bool {
	return !math.IsNaN(lat) && !math.IsNaN(lon) && lat >= -90 && lat <= 90 && lon >= -180 && lon <= 180
}

// Position returns the entity point: location first, then the doctrine envelope lat/lon, then
// the centre of its geo. ok is false when there is no valid point.
func (e *Entity) Position() (lat, lon float64, ok bool) {
	if e.Location != nil && validLL(e.Location.Lat, e.Location.Lon) {
		return e.Location.Lat, e.Location.Lon, true
	}
	if e.Lat != nil && e.Lon != nil && validLL(*e.Lat, *e.Lon) {
		return *e.Lat, *e.Lon, true
	}
	if e.Geo != nil {
		if pts := geoPoints(e.Geo); len(pts) > 0 {
			return centroid(pts)
		}
	}
	return 0, 0, false
}

// Altitude returns metres HAE when known.
func (e *Entity) Altitude() (float64, bool) {
	if e.Location != nil && e.Location.HAEM != nil {
		return *e.Location.HAEM, true
	}
	if e.HAE != nil {
		return *e.HAE, true
	}
	return 0, false
}

func (e *Entity) hasTrait(t string) bool {
	for _, x := range e.Traits {
		if x == t {
			return true
		}
	}
	return false
}

// IsOwnAsset is the "own asset" tier reason (cot-projection.md §3.2; decision 0004 #13).
func (e *Entity) IsOwnAsset() bool {
	return e.Template == "asset" && e.hasTrait("taskable") && strings.HasPrefix(e.EntityID, "sim-")
}

// IsPinnedType is the "pinned map content" tier (place, zone, structure).
func (e *Entity) IsPinnedType() bool {
	switch e.EntityType {
	case "place", "zone", "structure":
		return true
	}
	return false
}

// personKinds are track kinds treated as a person (entity.md §5 "Person safety"; SL3).
var personKinds = map[string]bool{"person": true, "human": true, "pedestrian": true, "people": true, "crowd": true}

// IsPersonClass reports a person entity or a person-class track (gate G4).
func (e *Entity) IsPersonClass() bool {
	if e.EntityType == "person" {
		return true
	}
	return personKinds[strings.ToLower(e.Kind)]
}
