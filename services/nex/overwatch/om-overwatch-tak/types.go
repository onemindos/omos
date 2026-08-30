package main

import "time"

// types.go — the om-overwatch data contract, mirrored from om-overwatch-state
// (ontology.go). The tak organ speaks the same entity/report shapes as the
// spine so CoT translation is a pure mapping with no schema drift.

// Affiliation — MIL-STD-2525 friend/foe axis (simplified).
const (
	AffiliationFriendly = "friendly"
	AffiliationHostile  = "hostile"
	AffiliationUnknown  = "unknown"
	AffiliationNeutral  = "neutral"
)

// Entity types (subset used for CoT mapping).
const (
	EntityTypeAircraftFixedWing  = "aircraft_fixed_wing"
	EntityTypeAircraftMultirotor = "aircraft_multirotor"
	EntityTypeAircraftVTOL       = "aircraft_vtol"
	EntityTypeAircraftHelicopter = "aircraft_helicopter"
	EntityTypeAircraftUnknown    = "aircraft_unknown"
	EntityTypeGroundVehicle      = "ground_vehicle"
	EntityTypeSurfaceVessel      = "surface_vessel"
	EntityTypeSensorPlatform     = "sensor_platform"
	EntityTypeOperatorStation    = "operator_station"
	EntityTypeInterceptor        = "interceptor"
)

// contactReport is the permissive report shape the spine accepts on det.cot.{id}.
type contactReport struct {
	ID          string   `json:"id"`
	EntityType  string   `json:"entity_type"`
	Affiliation string   `json:"affiliation"`
	Name        string   `json:"name"`
	Latitude    *float64 `json:"latitude"`
	Longitude   *float64 `json:"longitude"`
	Altitude    *float64 `json:"altitude"`
	Heading     *float64 `json:"heading"`
	Velocity    *float64 `json:"velocity"`
	Confidence  float64  `json:"confidence"`
	Source      string   `json:"source"`
}

// Position mirrors the spine's Position.
type Position struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Altitude  float64 `json:"altitude,omitempty"`
	Heading   float64 `json:"heading,omitempty"`
	Velocity  float64 `json:"velocity,omitempty"`
}

// Entity mirrors the spine's fused world-state record (fields we render to CoT).
type Entity struct {
	EntityID    string    `json:"entity_id"`
	Name        string    `json:"name,omitempty"`
	EntityType  string    `json:"entity_type"`
	Affiliation string    `json:"affiliation"`
	Status      string    `json:"status"`
	Priority    string    `json:"priority"`
	IsLive      bool      `json:"is_live"`
	Position    *Position `json:"position,omitempty"`
	Sources     []string  `json:"sources,omitempty"`
	LastSeen    time.Time `json:"last_seen"`
}

// hasSource reports whether a given source contributed to this entity.
func (e *Entity) hasSource(src string) bool {
	for _, s := range e.Sources {
		if s == src {
			return true
		}
	}
	return false
}
