package projector

import (
	"encoding/json"
	"fmt"
	"time"
)

// fixture builds an entity.v1-shaped record (services/schemas/entity/entity.v1.schema.json) as
// the entity service writes it, with the defaults of cot-projection.md §3.1 for a world feed:
// live, active, locatable, CLEAR, egress_tak true. mods edit the map before encoding.
func fixture(typ, id string, lat, lon float64, mods ...func(map[string]any)) map[string]any {
	now := time.Now().UnixMilli()
	m := map[string]any{
		"uid": "OM." + typ + "." + id, "type": "", "lat": lat, "lon": lon, "t": now,
		"entity_id": id, "entity_type": typ, "confidence": 0.9, "properties": map[string]any{},
		"v": 1, "seq": 1, "live": true, "lifecycle": "active", "created_t": now, "expiry_t": nil,
		"pinned": false, "template": "track", "traits": []any{"locatable", "trackable"},
		"aliases":  map[string]any{"name": nil, "ids": []any{}},
		"location": map[string]any{"lat": lat, "lon": lon, "measured_t": now},
		"geo":      nil, "current_task": nil, "tracked": nil, "correlation": nil, "relationships": []any{},
		"route": nil, "cot_security": nil,
		"indicators": map[string]any{"simulated": false, "exercise": false, "emergency": false, "starred": false, "egress_tak": true},
		"marking":    map[string]any{"tlp": "CLEAR", "releasable_to": []any{}},
		"provenance": map[string]any{"primary": "feed-test", "how": "m-g", "source_t": now},
		"sources": []any{map[string]any{"integration": "feed-test", "source_id": id, "data_type": "geo",
			"subject": "evt.world.test", "stream": "WORLD", "seq": 1, "first_seen_t": now, "last_seen_t": now,
			"source_t": now, "components": []any{"location"}}},
		"confidence_parts": map[string]any{"source_prior": 0.9, "track_quality": 1, "corroboration": 0},
	}
	for _, f := range mods {
		f(m)
	}
	return m
}

func set(k string, v any) func(map[string]any) { return func(m map[string]any) { m[k] = v } }

func withKind(k string) func(map[string]any) { return set("kind", k) }

func ind(k string, v bool) func(map[string]any) {
	return func(m map[string]any) { m["indicators"].(map[string]any)[k] = v }
}

func marking(tlp string, rel ...string) func(map[string]any) {
	return func(m map[string]any) {
		r := []any{}
		for _, x := range rel {
			r = append(r, x)
		}
		m["marking"] = map[string]any{"tlp": tlp, "releasable_to": r}
	}
}

// ownAsset makes a sim-* drone/robot/vehicle record (robotics-cot-profile.md §2).
func ownAsset(m map[string]any) {
	m["template"] = "asset"
	m["traits"] = []any{"locatable", "trackable", "taskable", "observable"}
	m["disposition"] = "friendly"
	m["indicators"].(map[string]any)["simulated"] = true
}

func zoneGeo(kind string, ring [][2]float64) func(map[string]any) {
	return func(m map[string]any) {
		coords := [][]float64{}
		for _, p := range ring {
			coords = append(coords, []float64{p[0], p[1]})
		}
		m["template"] = "geo"
		m["pinned"] = true
		m["traits"] = []any{"locatable"}
		m["geo"] = map[string]any{"shape": "polygon", "kind": kind,
			"coords": map[string]any{"type": "Polygon", "coordinates": [][][]float64{coords}}}
	}
}

func mustJSON(v any) []byte {
	b, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return b
}

func decode(m map[string]any) *Entity {
	var e Entity
	if err := json.Unmarshal(mustJSON(m), &e); err != nil {
		panic(err)
	}
	return &e
}

// boxRing is a closed [lon,lat] ring around a centre.
func boxRing(lat, lon, half float64) [][2]float64 {
	return [][2]float64{{lon - half, lat - half}, {lon + half, lat - half}, {lon + half, lat + half}, {lon - half, lat + half}, {lon - half, lat - half}}
}

// geofence builds a KV geofences value (minimal form, A6-P2).
func geofence(id, kind string, ring [][2]float64, active bool) []byte {
	coords := [][]float64{}
	for _, p := range ring {
		coords = append(coords, []float64{p[0], p[1]})
	}
	return mustJSON(map[string]any{"zone_id": id, "active": active, "geo": map[string]any{
		"shape": "polygon", "kind": kind, "coords": map[string]any{"type": "Polygon", "coordinates": [][][]float64{coords}}}})
}

func idN(prefix string, i int) string { return fmt.Sprintf("%s-%05d", prefix, i) }
