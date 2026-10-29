package policy

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"math"
	"sync"

	"github.com/nats-io/nats.go"
)

// Zone is a geofence from KV `geofences` (value: a zone entity record, entity.md `geo`).
type Zone struct {
	ID        string
	Kind      string // keep_in | keep_out | aoi | hazard | control_area
	Shape     string // polygon | circle
	Ring      [][2]float64
	Holes     [][][2]float64
	Center    [2]float64 // lon, lat
	RadiusM   float64
	AppliesTo []string // asset ids; empty = every asset
	Live      bool
	// Unevaluable marks a keep_in/keep_out record (or a record whose kind cannot be read)
	// whose shape the policy cannot evaluate. It fails closed: movement verbs are denied for
	// the assets it applies to (detail geofence_unevaluable) until the shape is supported.
	Unevaluable bool
	Problem     string
}

// DetailUnevaluable is the denial detail for a movement blocked by an unreadable zone.
const DetailUnevaluable = "geofence_unevaluable"

// ParseZone decodes a KV value. It returns false only for a zone that does not constrain
// movement (aoi, hazard, control_area) and whose shape cannot be read; an unreadable
// keep_in/keep_out record, or one whose kind cannot be read, comes back Unevaluable.
func ParseZone(key string, value []byte) (Zone, bool) {
	var rec struct {
		EntityID string `json:"entity_id"`
		Live     *bool  `json:"live"`
		Geo      *struct {
			Shape   string          `json:"shape"`
			Coords  json.RawMessage `json:"coords"`
			RadiusM float64         `json:"radius_m"`
			Kind    string          `json:"kind"`
		} `json:"geo"`
		Properties struct {
			AppliesTo []string `json:"applies_to"`
		} `json:"properties"`
	}
	if json.Unmarshal(value, &rec) != nil {
		return Zone{ID: key, Live: true, Unevaluable: true, Problem: "record is not valid JSON"}, true
	}
	if rec.Geo == nil {
		// The kind is unknown, so it may be a keep_out: fail closed for every asset it names.
		return Zone{ID: key, Live: rec.Live == nil || *rec.Live, AppliesTo: rec.Properties.AppliesTo,
			Unevaluable: true, Problem: "record has no geo"}, true
	}
	z := Zone{ID: rec.EntityID, Kind: rec.Geo.Kind, Shape: rec.Geo.Shape, RadiusM: rec.Geo.RadiusM,
		AppliesTo: rec.Properties.AppliesTo, Live: rec.Live == nil || *rec.Live}
	if z.ID == "" {
		z.ID = key
	}
	bad := func(problem string) (Zone, bool) {
		if z.Kind != "" && z.Kind != "keep_in" && z.Kind != "keep_out" {
			return Zone{}, false // does not constrain movement
		}
		z.Unevaluable, z.Problem = true, problem
		return z, true
	}
	if z.Kind == "" {
		return bad("zone kind is missing")
	}
	var geom struct {
		Type        string          `json:"type"`
		Coordinates json.RawMessage `json:"coordinates"`
	}
	if json.Unmarshal(rec.Geo.Coords, &geom) != nil {
		return bad("coords are not GeoJSON")
	}
	switch {
	case z.Shape == "polygon" && geom.Type == "Polygon":
		var rings [][][2]float64
		if json.Unmarshal(geom.Coordinates, &rings) != nil || len(rings) == 0 || len(rings[0]) < 4 {
			return bad("polygon coordinates are malformed")
		}
		z.Ring, z.Holes = rings[0], rings[1:]
	case z.Shape == "circle" && geom.Type == "Point" && z.RadiusM > 0:
		if json.Unmarshal(geom.Coordinates, &z.Center) != nil {
			return bad("circle center is malformed")
		}
	default:
		return bad("shape " + z.Shape + " with geometry " + geom.Type + " is not supported")
	}
	return z, true
}

// UnevaluableZone returns the first live unevaluable zone that applies to the asset.
func UnevaluableZone(zones []Zone, asset string) (Zone, bool) {
	for _, z := range zones {
		if z.Live && z.Unevaluable && z.appliesTo(asset) {
			return z, true
		}
	}
	return Zone{}, false
}

func pointInRing(lon, lat float64, ring [][2]float64) bool {
	in := false
	for i, j := 0, len(ring)-1; i < len(ring); j, i = i, i+1 {
		xi, yi, xj, yj := ring[i][0], ring[i][1], ring[j][0], ring[j][1]
		if (yi > lat) != (yj > lat) && lon < (xj-xi)*(lat-yi)/(yj-yi)+xi {
			in = !in
		}
	}
	return in
}

func haversineM(lat1, lon1, lat2, lon2 float64) float64 {
	const r = 6371008.8
	p1, p2 := lat1*math.Pi/180, lat2*math.Pi/180
	dp, dl := (lat2-lat1)*math.Pi/180, (lon2-lon1)*math.Pi/180
	a := math.Sin(dp/2)*math.Sin(dp/2) + math.Cos(p1)*math.Cos(p2)*math.Sin(dl/2)*math.Sin(dl/2)
	return 2 * r * math.Asin(math.Sqrt(a))
}

// Contains reports whether the point is inside the zone.
func (z Zone) Contains(lat, lon float64) bool {
	switch z.Shape {
	case "polygon":
		if !pointInRing(lon, lat, z.Ring) {
			return false
		}
		for _, h := range z.Holes {
			if pointInRing(lon, lat, h) {
				return false
			}
		}
		return true
	case "circle":
		return haversineM(lat, lon, z.Center[1], z.Center[0]) <= z.RadiusM
	}
	return false
}

func (z Zone) appliesTo(asset string) bool {
	if len(z.AppliesTo) == 0 {
		return true
	}
	for _, a := range z.AppliesTo {
		if a == asset {
			return true
		}
	}
	return false
}

// CheckPoint returns "" when the point is allowed for the asset, else keep_out, keep_in or
// geofence_unevaluable. A point must be outside every live keep-out zone and, when any
// keep-in zone applies to the asset, inside at least one of them; any unevaluable zone that
// applies to the asset denies the point (fail closed).
func CheckPoint(zones []Zone, asset string, lat, lon float64) string {
	if _, ok := UnevaluableZone(zones, asset); ok {
		return DetailUnevaluable
	}
	haveKeepIn, inKeepIn := false, false
	for _, z := range zones {
		if !z.Live || !z.appliesTo(asset) {
			continue
		}
		switch z.Kind {
		case "keep_out":
			if z.Contains(lat, lon) {
				return "keep_out"
			}
		case "keep_in":
			haveKeepIn = true
			if z.Contains(lat, lon) {
				inKeepIn = true
			}
		}
	}
	if haveKeepIn && !inKeepIn {
		return "keep_in"
	}
	return ""
}

// ZoneSource supplies the current geofences.
type ZoneSource interface {
	Zones() ([]Zone, error)
}

// KVZones keeps a live copy of KV `geofences` through a watcher.
type KVZones struct {
	mu    sync.RWMutex
	zones map[string]Zone
	ready bool
	err   error
}

// WatchZones starts watching the bucket. Zones() fails until the initial values are loaded.
// An unevaluable zone is kept (it denies movement) and logged as a warning.
func WatchZones(ctx context.Context, js nats.JetStreamContext, bucket string, log *slog.Logger) (*KVZones, error) {
	if log == nil {
		log = slog.Default()
	}
	kv, err := js.KeyValue(bucket)
	if err != nil {
		return nil, err
	}
	w, err := kv.WatchAll(nats.Context(ctx))
	if err != nil {
		return nil, err
	}
	k := &KVZones{zones: map[string]Zone{}}
	go func() {
		for {
			select {
			case <-ctx.Done():
				_ = w.Stop()
				return
			case e, ok := <-w.Updates():
				if !ok {
					k.mu.Lock()
					k.err = errors.New("geofence watcher stopped")
					k.mu.Unlock()
					return
				}
				k.mu.Lock()
				if e == nil {
					k.ready = true
				} else if e.Operation() == nats.KeyValuePut {
					if z, ok := ParseZone(e.Key(), e.Value()); ok {
						if z.Unevaluable {
							log.Warn("geofence cannot be evaluated; movement is denied for the assets it applies to",
								"key", e.Key(), "zone", z.ID, "kind", z.Kind, "problem", z.Problem, "applies_to", z.AppliesTo)
						}
						k.zones[e.Key()] = z
					} else {
						delete(k.zones, e.Key())
					}
				} else {
					delete(k.zones, e.Key())
				}
				k.mu.Unlock()
			}
		}
	}()
	return k, nil
}

func (k *KVZones) Zones() ([]Zone, error) {
	k.mu.RLock()
	defer k.mu.RUnlock()
	if k.err != nil {
		return nil, k.err
	}
	if !k.ready {
		return nil, errors.New("geofences not loaded yet")
	}
	out := make([]Zone, 0, len(k.zones))
	for _, z := range k.zones {
		out = append(out, z)
	}
	return out, nil
}

// StaticZones is a fixed zone list (tests, local runs without KV).
type StaticZones []Zone

func (s StaticZones) Zones() ([]Zone, error) { return s, nil }
