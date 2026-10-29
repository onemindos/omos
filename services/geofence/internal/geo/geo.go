// Package geo holds the zone geometry tests of the geofence service: polygons with holes,
// multipolygons, circles and corridors on WGS-84 coordinates, with optional altitude bands.
//
// The tests are planar on lon/lat (ray casting), as in the policy service
// (services/policy/internal/policy/geo.go, PR #554) and the TAK projector
// (services/projector/internal/projector/geo.go, PR #563), so the three readers of KV
// `geofences` agree on what is inside. Circles use the haversine distance. Corridors are
// buffered into a polygon with the projector's mitred-normal method. Zones that cross the
// antimeridian or contain a pole are not supported (cot-projection.md §3.2: fine for zones of
// a few tens of kilometres).
package geo

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
)

// Pt is [lon, lat] (GeoJSON order).
type Pt [2]float64

// Shape names accepted in a geofence record (`geo.shape`).
const (
	ShapePolygon      = "polygon"
	ShapeMultiPolygon = "multipolygon"
	ShapeCircle       = "circle"
	ShapeCorridor     = "corridor"
)

// MaxVertices bounds one zone (all rings together).
const MaxVertices = 5000

// MaxRadiusM bounds circles and corridor half-widths.
const MaxRadiusM = 200000.0

// Area is a resolved zone shape with its bounding box.
type Area struct {
	MinLon, MinLat, MaxLon, MaxLat float64
	polys                          [][][]Pt // polygons: outer ring then holes
	circle                         *[3]float64
	AltMin, AltMax                 *float64
}

var errShape = errors.New("geometry does not match the shape")

// Parse resolves a GeoJSON geometry for a shape. radiusM is used by circle and corridor.
func Parse(shape string, coords json.RawMessage, radiusM float64, altMin, altMax *float64) (*Area, error) {
	var gj struct {
		Type        string          `json:"type"`
		Coordinates json.RawMessage `json:"coordinates"`
	}
	if err := json.Unmarshal(coords, &gj); err != nil {
		return nil, fmt.Errorf("coords are not GeoJSON: %w", err)
	}
	if altMin != nil && altMax != nil && *altMin >= *altMax {
		return nil, errors.New("alt_min_m must be below alt_max_m")
	}
	a := &Area{AltMin: altMin, AltMax: altMax}
	switch shape {
	case ShapePolygon:
		if gj.Type != "Polygon" {
			return nil, fmt.Errorf("%w: polygon needs a GeoJSON Polygon, got %q", errShape, gj.Type)
		}
		var rings [][][]float64
		if err := json.Unmarshal(gj.Coordinates, &rings); err != nil {
			return nil, fmt.Errorf("polygon coordinates: %w", err)
		}
		p, err := polygon(rings)
		if err != nil {
			return nil, err
		}
		a.polys = [][][]Pt{p}
	case ShapeMultiPolygon:
		if gj.Type != "MultiPolygon" {
			return nil, fmt.Errorf("%w: multipolygon needs a GeoJSON MultiPolygon, got %q", errShape, gj.Type)
		}
		var polys [][][][]float64
		if err := json.Unmarshal(gj.Coordinates, &polys); err != nil {
			return nil, fmt.Errorf("multipolygon coordinates: %w", err)
		}
		if len(polys) == 0 {
			return nil, errors.New("multipolygon has no polygons")
		}
		for _, rings := range polys {
			p, err := polygon(rings)
			if err != nil {
				return nil, err
			}
			a.polys = append(a.polys, p)
		}
	case ShapeCircle:
		if gj.Type != "Point" {
			return nil, fmt.Errorf("%w: circle needs a GeoJSON Point, got %q", errShape, gj.Type)
		}
		var p []float64
		if err := json.Unmarshal(gj.Coordinates, &p); err != nil || len(p) < 2 || !validLL(p[1], p[0]) {
			return nil, errors.New("circle centre is malformed")
		}
		if radiusM <= 0 || radiusM > MaxRadiusM {
			return nil, fmt.Errorf("circle radius_m must be in (0, %.0f]", MaxRadiusM)
		}
		a.circle = &[3]float64{p[0], p[1], radiusM}
	case ShapeCorridor:
		if gj.Type != "LineString" {
			return nil, fmt.Errorf("%w: corridor needs a GeoJSON LineString, got %q", errShape, gj.Type)
		}
		var ps [][]float64
		if err := json.Unmarshal(gj.Coordinates, &ps); err != nil {
			return nil, fmt.Errorf("corridor coordinates: %w", err)
		}
		line, err := points(ps)
		if err != nil {
			return nil, err
		}
		if len(line) < 2 {
			return nil, errors.New("corridor needs at least two points")
		}
		if radiusM <= 0 || radiusM > MaxRadiusM {
			return nil, fmt.Errorf("corridor radius_m must be in (0, %.0f]", MaxRadiusM)
		}
		a.polys = [][][]Pt{{bufferLine(line, radiusM)}}
	default:
		return nil, fmt.Errorf("unknown shape %q (want polygon, multipolygon, circle or corridor)", shape)
	}
	n := 0
	for _, p := range a.polys {
		for _, r := range p {
			n += len(r)
		}
	}
	if n > MaxVertices {
		return nil, fmt.Errorf("zone has %d vertices, more than %d", n, MaxVertices)
	}
	a.bbox()
	return a, nil
}

func validLL(lat, lon float64) bool {
	return lat >= -90 && lat <= 90 && lon >= -180 && lon <= 180 && !math.IsNaN(lat) && !math.IsNaN(lon)
}

func points(in [][]float64) ([]Pt, error) {
	out := make([]Pt, 0, len(in))
	for _, p := range in {
		if len(p) < 2 || !validLL(p[1], p[0]) {
			return nil, errors.New("coordinate out of range or malformed")
		}
		out = append(out, Pt{p[0], p[1]})
	}
	return out, nil
}

// polygon reads GeoJSON rings: each ring needs at least three distinct vertices and is closed
// if the writer left it open.
func polygon(rings [][][]float64) ([][]Pt, error) {
	if len(rings) == 0 {
		return nil, errors.New("polygon has no rings")
	}
	out := make([][]Pt, 0, len(rings))
	for i, r := range rings {
		ring, err := points(r)
		if err != nil {
			return nil, err
		}
		if len(ring) > 0 && ring[0] != ring[len(ring)-1] {
			ring = append(ring, ring[0])
		}
		if len(ring) < 4 {
			what := "outer ring"
			if i > 0 {
				what = "hole"
			}
			return nil, fmt.Errorf("polygon %s needs at least three distinct vertices", what)
		}
		out = append(out, ring)
	}
	return out, nil
}

// Closed returns the GeoJSON rings of a polygon with every ring closed (used to normalise
// what the service writes, so readers that require closed rings accept it).
func Closed(rings [][][]float64) [][][]float64 {
	out := make([][][]float64, 0, len(rings))
	for _, r := range rings {
		c := append([][]float64{}, r...)
		if n := len(c); n > 0 && (c[0][0] != c[n-1][0] || c[0][1] != c[n-1][1]) {
			c = append(c, []float64{c[0][0], c[0][1]})
		}
		out = append(out, c)
	}
	return out
}

func (a *Area) bbox() {
	a.MinLon, a.MinLat, a.MaxLon, a.MaxLat = 180, 90, -180, -90
	grow := func(lon, lat float64) {
		a.MinLon, a.MaxLon = math.Min(a.MinLon, lon), math.Max(a.MaxLon, lon)
		a.MinLat, a.MaxLat = math.Min(a.MinLat, lat), math.Max(a.MaxLat, lat)
	}
	if a.circle != nil {
		lon, lat, r := a.circle[0], a.circle[1], a.circle[2]
		dLat := r / 111320.0
		dLon := dLat / math.Max(math.Cos(lat*math.Pi/180), 1e-6)
		grow(math.Max(lon-dLon, -180), math.Max(lat-dLat, -90))
		grow(math.Min(lon+dLon, 180), math.Min(lat+dLat, 90))
		return
	}
	for _, p := range a.polys {
		for _, q := range p[0] {
			grow(q[0], q[1])
		}
	}
}

// Contains reports whether the point is inside the area. The altitude band applies only when
// the altitude is known (hasAlt), as in the projector.
func (a *Area) Contains(lat, lon, alt float64, hasAlt bool) bool {
	if lat < a.MinLat || lat > a.MaxLat || lon < a.MinLon || lon > a.MaxLon {
		return false
	}
	if hasAlt && ((a.AltMin != nil && alt < *a.AltMin) || (a.AltMax != nil && alt > *a.AltMax)) {
		return false
	}
	if a.circle != nil {
		return HaversineM(lat, lon, a.circle[1], a.circle[0]) <= a.circle[2]
	}
	for _, p := range a.polys {
		if !pointInRing(lon, lat, p[0]) {
			continue
		}
		inHole := false
		for _, h := range p[1:] {
			if pointInRing(lon, lat, h) {
				inHole = true
				break
			}
		}
		if !inHole {
			return true
		}
	}
	return false
}

func pointInRing(x, y float64, ring []Pt) bool {
	in := false
	for i, j := 0, len(ring)-1; i < len(ring); j, i = i, i+1 {
		xi, yi, xj, yj := ring[i][0], ring[i][1], ring[j][0], ring[j][1]
		if (yi > y) != (yj > y) && x < (xj-xi)*(y-yi)/(yj-yi)+xi {
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
	h := math.Sin(dp/2)*math.Sin(dp/2) + math.Cos(p1)*math.Cos(p2)*math.Sin(dl/2)*math.Sin(dl/2)
	return 2 * r * math.Asin(math.Min(1, math.Sqrt(h)))
}

// bufferLine buffers a polyline by r metres into a closed polygon, the same construction as
// the projector's bufferLine (left side forward, right side back, mitred by the averaged
// normal), so a corridor means the same area to both services.
func bufferLine(line []Pt, r float64) []Pt {
	n := len(line)
	left := make([]Pt, n)
	right := make([]Pt, n)
	for i := range line {
		var dx, dy float64
		if i > 0 {
			dx += line[i][0] - line[i-1][0]
			dy += line[i][1] - line[i-1][1]
		}
		if i < n-1 {
			dx += line[i+1][0] - line[i][0]
			dy += line[i+1][1] - line[i][1]
		}
		cosLat := math.Max(math.Cos(line[i][1]*math.Pi/180), 1e-6)
		ex, ey := dx*cosLat, dy
		l := math.Hypot(ex, ey)
		if l == 0 {
			l = 1
		}
		nx, ny := -ey/l, ex/l
		dLat := ny * r / 111320.0
		dLon := nx * r / 111320.0 / cosLat
		left[i] = Pt{line[i][0] + dLon, line[i][1] + dLat}
		right[i] = Pt{line[i][0] - dLon, line[i][1] - dLat}
	}
	out := append([]Pt{}, left...)
	for i := n - 1; i >= 0; i-- {
		out = append(out, right[i])
	}
	return append(out, out[0])
}
