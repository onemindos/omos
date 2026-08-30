package main

import "time"

// Entity types — the om-overwatch ontology. Derived from Constellation Overwatch
// (pkg/ontology/entity.go) but trimmed to the counter-UAS domain and made
// NATS-native. Air contacts are the primary concern; ground/surface/sensor
// types are retained for the broader COP.
const (
	EntityTypeAircraftFixedWing  = "aircraft_fixed_wing"
	EntityTypeAircraftMultirotor = "aircraft_multirotor"
	EntityTypeAircraftVTOL       = "aircraft_vtol"
	EntityTypeAircraftHelicopter = "aircraft_helicopter"
	EntityTypeAircraftUnknown    = "aircraft_unknown" // unclassified air contact (threat track)
	EntityTypeGroundVehicle      = "ground_vehicle"
	EntityTypeSurfaceVessel      = "surface_vessel"
	EntityTypeSensorPlatform     = "sensor_platform"
	EntityTypeOperatorStation    = "operator_station"
	EntityTypeInterceptor        = "interceptor" // our own kill vehicle
)

// Affiliation — the MIL-STD-2525 friend/hostile axis, simplified.
const (
	AffiliationFriendly = "friendly"
	AffiliationHostile  = "hostile"
	AffiliationUnknown  = "unknown"
	AffiliationNeutral  = "neutral"
)

// Status of the track.
const (
	StatusActive   = "active"
	StatusInactive = "inactive"
	StatusLost     = "lost" // no updates within staleness window
)

// Priority / threat level.
const (
	PriorityLow      = "low"
	PriorityNormal   = "normal"
	PriorityHigh     = "high"
	PriorityCritical = "critical"
)

// Position is a geographic point plus motion.
type Position struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Altitude  float64 `json:"altitude,omitempty"`  // meters
	Heading   float64 `json:"heading,omitempty"`   // degrees true
	Velocity  float64 `json:"velocity,omitempty"`  // m/s ground speed
}

// Entity is the fused world-state record for one tracked object. This is what
// lives in the OVERWATCH_STATE KV bucket, keyed by EntityID, and what gets
// published on ent.air.{id}.state after each fusion.
type Entity struct {
	EntityID    string                 `json:"entity_id"`
	Name        string                 `json:"name,omitempty"`
	EntityType  string                 `json:"entity_type"`
	Affiliation string                 `json:"affiliation"`
	Status      string                 `json:"status"`
	Priority    string                 `json:"priority"`
	IsLive      bool                   `json:"is_live"`
	Pinned      bool                   `json:"pinned,omitempty"` // manual entities (no-fly zones, waypoints) — exempt from staleness sweep
	// AffiliationLocked marks affiliation as human-asserted (via TAK CoT or the
	// operator API). Once locked, machine sensors (RF/acoustic/visual/telemetry)
	// may still update position and other fields but MUST NOT overwrite the
	// affiliation — human-in-the-loop is the final authority on friend/foe.
	AffiliationLocked bool `json:"affiliation_locked,omitempty"`
	Position    *Position              `json:"position,omitempty"`
	Sources     []string               `json:"sources,omitempty"` // which sensors have reported this track
	Confidence  float64                `json:"confidence,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
	FirstSeen   time.Time              `json:"first_seen"`
	LastSeen    time.Time              `json:"last_seen"`
	UpdatedAt   time.Time              `json:"updated_at"`
}

// hasSource reports whether a sensor has already contributed to this track.
func (e *Entity) hasSource(src string) bool {
	for _, s := range e.Sources {
		if s == src {
			return true
		}
	}
	return false
}

// addSource records a sensor as having contributed to this track (deduped).
func (e *Entity) addSource(src string) {
	if src == "" {
		return
	}
	for _, s := range e.Sources {
		if s == src {
			return
		}
	}
	e.Sources = append(e.Sources, src)
}
