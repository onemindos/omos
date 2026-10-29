package engine

import (
	"math"
	"strings"
)

// The numbers below are starting points from the design documents, not values measured
// against feed cadence (entity.md §11; fusion-and-tracking.md "Unverified items").

// ttlSeconds is the per-type Nats-TTL default of entity.md §4.3 and BUS-DESIGN §4.1:
// aircraft 120 s, vessel 30 min, fleet asset 60 s, TAK EUD = CoT stale. A person track is
// capped at 60 s (entity.md §5 "Person safety").
func ttlSeconds(entityType, kind string, sourceStale int) int {
	switch entityType {
	case "aircraft":
		return 120
	case "vessel":
		return 1800
	case "drone", "robot", "node":
		return 60
	case "event":
		return 86400
	case "track":
		if kind == "person" {
			if sourceStale > 0 && sourceStale < 60 {
				return sourceStale
			}
			return 60
		}
	}
	if sourceStale > 0 {
		return sourceStale
	}
	switch entityType {
	case "person", "vehicle", "agent":
		return 300
	}
	return 120
}

// defaultStale is the `stale` (seconds) used when the source does not say.
func defaultStale(entityType string) int {
	switch entityType {
	case "aircraft", "track":
		return 30
	case "drone", "robot":
		return 30
	case "vessel":
		return 600
	case "event":
		return 3600
	}
	return 120
}

// defaultSigma is the 1-sigma horizontal error (metres) used when a source gives none
// (CoT ce 9999999 means unknown). Risk 7 of fusion-and-tracking.md: tune on fixtures.
func defaultSigma(entityType, kind string) float64 {
	if kind == "satellite" {
		return 1000
	}
	switch entityType {
	case "aircraft", "vehicle":
		return 50
	case "vessel":
		return 20
	case "drone", "robot":
		return 5
	case "person":
		return 10
	case "event":
		return 1000
	}
	return 25
}

// maxSpeed is the plausibility limit (m/s) for the AIS duplicate-MMSI split (Stage 1, F8).
func maxSpeed(entityType string) float64 {
	if entityType == "vessel" {
		return 50 // about 97 knots
	}
	return math.Inf(1)
}

// sourcePrior is the per-feed reliability. The feed registry (A2) will own these values.
func sourcePrior(integration string, overrides map[string]float64) float64 {
	if v, ok := overrides[integration]; ok {
		return v
	}
	exact := map[string]float64{
		"tak-server": 0.9, "tak-mesh": 0.8, "tak-inject": 0.6,
		"feed-aircraft": 0.85, "world-air": 0.85, "feed-ships": 0.8, "world-sea": 0.8,
		"feed-satellites": 0.9, "feed-mta": 0.85,
	}
	if v, ok := exact[integration]; ok {
		return v
	}
	for prefix, v := range map[string]float64{"fleet-": 0.95, "feed-intel-": 0.7, "world-": 0.7, "det-": 0.5} {
		if strings.HasPrefix(integration, prefix) {
			return v
		}
	}
	return 0.6
}

var tlpRank = map[string]int{"CLEAR": 0, "GREEN": 1, "AMBER": 2, "AMBER+STRICT": 3, "RED": 4}

// stricterTLP returns the stricter of two TLP markings.
func stricterTLP(a, b string) string {
	if a == "" {
		return b
	}
	if tlpRank[b] > tlpRank[a] {
		return b
	}
	return a
}

func round(v float64, places int) float64 {
	p := math.Pow(10, float64(places))
	return math.Round(v*p) / p
}

func ptr(v float64) *float64 { return &v }
