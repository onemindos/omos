// Package geo reads zones from KV `geofences` (value: a zone entity record with an entity.md
// `geo` block, the same format the policy service reads) and answers point-in-zone and
// bounding-box questions. Polygons (with holes) and circles are supported; any other shape is
// reported as unreadable so a watch on it fails loudly instead of silently never firing.
package geo

import (
	"encoding/json"
	"errors"
	"math"
)

// Zone is one geofence.
type Zone struct {
	ID      string
	Kind    string // keep_in | keep_out | aoi | hazard | control_area
	Shape   string // polygon | circle
	Name    string
	Ring    [][2]float64 // lon, lat
	Holes   [][][2]float64
	Center  [2]float64 // lon, lat
	RadiusM float64
	Live    bool
}

// ParseZone decodes one KV value.
func ParseZone(key string, value []byte) (Zone, error) {
	var rec struct {
		EntityID string `json:"entity_id"`
		Callsign string `json:"callsign"`
		Live     *bool  `json:"live"`
		Geo      *struct {
			Shape   string          `json:"shape"`
			Coords  json.RawMessage `json:"coords"`
			RadiusM float64         `json:"radius_m"`
			Kind    string          `json:"kind"`
		} `json:"geo"`
	}
	if err := json.Unmarshal(value, &rec); err != nil {
		return Zone{ID: key}, err
	}
	if rec.Geo == nil {
		return Zone{ID: key}, errors.New("zone record has no geo")
	}
	z := Zone{ID: rec.EntityID, Kind: rec.Geo.Kind, Shape: rec.Geo.Shape, Name: rec.Callsign,
		RadiusM: rec.Geo.RadiusM, Live: rec.Live == nil || *rec.Live}
	if z.ID == "" {
		z.ID = key
	}
	var g struct {
		Type        string          `json:"type"`
		Coordinates json.RawMessage `json:"coordinates"`
	}
	if err := json.Unmarshal(rec.Geo.Coords, &g); err != nil {
		return z, errors.New("coords are not GeoJSON")
	}
	switch {
	case z.Shape == "polygon" && g.Type == "Polygon":
		var rings [][][2]float64
		if json.Unmarshal(g.Coordinates, &rings) != nil || len(rings) == 0 || len(rings[0]) < 4 {
			return z, errors.New("polygon coordinates are malformed")
		}
		z.Ring, z.Holes = rings[0], rings[1:]
	case z.Shape == "circle" && g.Type == "Point" && z.RadiusM > 0:
		if json.Unmarshal(g.Coordinates, &z.Center) != nil {
			return z, errors.New("circle center is malformed")
		}
	default:
		return z, errors.New("shape " + z.Shape + " with geometry " + g.Type + " is not supported")
	}
	return z, nil
}

// Contains reports whether the point is inside the zone.
func (z Zone) Contains(lat, lon float64) bool {
	switch z.Shape {
	case "circle":
		return HaversineM(lat, lon, z.Center[1], z.Center[0]) <= z.RadiusM
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
	}
	return false
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

// HaversineM is the great-circle distance in metres.
func HaversineM(lat1, lon1, lat2, lon2 float64) float64 {
	const r = 6371008.8
	p1, p2 := lat1*math.Pi/180, lat2*math.Pi/180
	dp, dl := (lat2-lat1)*math.Pi/180, (lon2-lon1)*math.Pi/180
	a := math.Sin(dp/2)*math.Sin(dp/2) + math.Cos(p1)*math.Cos(p2)*math.Sin(dl/2)*math.Sin(dl/2)
	return 2 * r * math.Asin(math.Min(1, math.Sqrt(a)))
}

// BBox is [lon0, lat0, lon1, lat1].
type BBox [4]float64

// Valid reports whether the box is ordered and in range.
func (b BBox) Valid() bool {
	return b[0] >= -180 && b[2] <= 180 && b[1] >= -90 && b[3] <= 90 && b[0] < b[2] && b[1] < b[3]
}

// Contains reports whether the point is in the box.
func (b BBox) Contains(lat, lon float64) bool {
	return lon >= b[0] && lon <= b[2] && lat >= b[1] && lat <= b[3]
}

// Bounds returns the zone's bounding box.
func (z Zone) Bounds() BBox {
	if z.Shape == "circle" {
		dLat := z.RadiusM / 111320.0
		dLon := dLat / math.Max(math.Cos(z.Center[1]*math.Pi/180), 0.01)
		return BBox{z.Center[0] - dLon, z.Center[1] - dLat, z.Center[0] + dLon, z.Center[1] + dLat}
	}
	b := BBox{180, 90, -180, -90}
	for _, p := range z.Ring {
		b[0], b[1] = math.Min(b[0], p[0]), math.Min(b[1], p[1])
		b[2], b[3] = math.Max(b[2], p[0]), math.Max(b[3], p[1])
	}
	return b
}
