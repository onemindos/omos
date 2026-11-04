// Package fence is the bridge's own keep-in / keep-out check (asset-command.md §7 rule 3):
// it runs before any movement command, independent of the policy service.
//
// It fails closed: with no keep-in zone configured, every position is outside.
// Polygons use planar ray casting on lat/lon, which is adequate for site-sized zones and
// does not handle zones that cross the antimeridian or contain a pole.
package fence

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
)

// Point is a WGS-84 position in decimal degrees.
type Point struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}

// Circle is a zone of RadiusM metres around a centre.
type Circle struct {
	Lat     float64 `json:"lat"`
	Lon     float64 `json:"lon"`
	RadiusM float64 `json:"radius_m"`
}

// Zone is either a polygon (3 or more vertices) or a circle.
type Zone struct {
	ID      string  `json:"id"`
	Polygon []Point `json:"polygon,omitempty"`
	Circle  *Circle `json:"circle,omitempty"`
}

// Fence is the bridge geofence configuration.
type Fence struct {
	KeepIn  []Zone   `json:"keep_in"`
	KeepOut []Zone   `json:"keep_out"`
	MaxHaeM *float64 `json:"max_hae_m,omitempty"`
	MinHaeM *float64 `json:"min_hae_m,omitempty"`
}

// Violation says why a position failed the check. Detail uses the
// precondition_failed details of asset-command.md §6.4 (keep_in, keep_out).
type Violation struct {
	Detail string
	ZoneID string
	Reason string
}

func (v *Violation) Error() string { return v.Reason }

// Load reads a fence from a JSON file.
func Load(path string) (*Fence, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return Parse(b)
}

// Parse decodes and validates a fence document.
func Parse(b []byte) (*Fence, error) {
	var f Fence
	if err := json.Unmarshal(b, &f); err != nil {
		return nil, fmt.Errorf("geofence: %w", err)
	}
	for _, z := range append(append([]Zone{}, f.KeepIn...), f.KeepOut...) {
		if err := z.validate(); err != nil {
			return nil, err
		}
	}
	if f.MaxHaeM != nil && f.MinHaeM != nil && *f.MinHaeM > *f.MaxHaeM {
		return nil, fmt.Errorf("geofence: min_hae_m > max_hae_m")
	}
	return &f, nil
}

func validPoint(lat, lon float64) bool {
	return !math.IsNaN(lat) && !math.IsNaN(lon) && lat >= -90 && lat <= 90 && lon >= -180 && lon <= 180
}

func (z Zone) validate() error {
	switch {
	case z.ID == "":
		return fmt.Errorf("geofence: zone without id")
	case z.Circle != nil && len(z.Polygon) > 0:
		return fmt.Errorf("geofence: zone %s has both polygon and circle", z.ID)
	case z.Circle != nil:
		if !validPoint(z.Circle.Lat, z.Circle.Lon) || !(z.Circle.RadiusM > 0) {
			return fmt.Errorf("geofence: zone %s: bad circle", z.ID)
		}
	case len(z.Polygon) >= 3:
		for _, p := range z.Polygon {
			if !validPoint(p.Lat, p.Lon) {
				return fmt.Errorf("geofence: zone %s: bad vertex %v", z.ID, p)
			}
		}
	default:
		return fmt.Errorf("geofence: zone %s needs a circle or 3+ polygon vertices", z.ID)
	}
	return nil
}

// Contains reports whether p is inside the zone (boundary counts as inside for circles).
func (z Zone) Contains(p Point) bool {
	if z.Circle != nil {
		return Distance(Point{z.Circle.Lat, z.Circle.Lon}, p) <= z.Circle.RadiusM
	}
	in := false
	n := len(z.Polygon)
	for i, j := 0, n-1; i < n; j, i = i, i+1 {
		a, b := z.Polygon[i], z.Polygon[j]
		if (a.Lat > p.Lat) != (b.Lat > p.Lat) &&
			p.Lon < (b.Lon-a.Lon)*(p.Lat-a.Lat)/(b.Lat-a.Lat)+a.Lon {
			in = !in
		}
	}
	return in
}

// Distance is the haversine distance in metres.
func Distance(a, b Point) float64 {
	const r = 6371008.8
	la1, la2 := a.Lat*math.Pi/180, b.Lat*math.Pi/180
	dla, dlo := la2-la1, (b.Lon-a.Lon)*math.Pi/180
	h := math.Sin(dla/2)*math.Sin(dla/2) + math.Cos(la1)*math.Cos(la2)*math.Sin(dlo/2)*math.Sin(dlo/2)
	return 2 * r * math.Asin(math.Min(1, math.Sqrt(h)))
}

// Check returns nil when p (and hae, if given) is inside at least one keep-in zone,
// outside every keep-out zone and within the altitude band.
func (f *Fence) Check(p Point, hae *float64) *Violation {
	if f == nil || len(f.KeepIn) == 0 {
		return &Violation{Detail: "keep_in", Reason: "no keep-in zone configured (fail closed)"}
	}
	if !validPoint(p.Lat, p.Lon) {
		return &Violation{Detail: "keep_in", Reason: fmt.Sprintf("invalid position %v,%v", p.Lat, p.Lon)}
	}
	for _, z := range f.KeepOut {
		if z.Contains(p) {
			return &Violation{Detail: "keep_out", ZoneID: z.ID, Reason: "position inside keep-out zone " + z.ID}
		}
	}
	inside := false
	for _, z := range f.KeepIn {
		if z.Contains(p) {
			inside = true
			break
		}
	}
	if !inside {
		return &Violation{Detail: "keep_in", Reason: "position outside every keep-in zone"}
	}
	if hae != nil {
		if math.IsNaN(*hae) || math.IsInf(*hae, 0) {
			return &Violation{Detail: "keep_in", Reason: "invalid altitude"}
		}
		if f.MaxHaeM != nil && *hae > *f.MaxHaeM {
			return &Violation{Detail: "keep_in", Reason: fmt.Sprintf("hae %.1f m above ceiling %.1f m", *hae, *f.MaxHaeM)}
		}
		if f.MinHaeM != nil && *hae < *f.MinHaeM {
			return &Violation{Detail: "keep_in", Reason: fmt.Sprintf("hae %.1f m below floor %.1f m", *hae, *f.MinHaeM)}
		}
	}
	return nil
}
