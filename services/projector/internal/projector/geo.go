package projector

// Geometry on WGS-84 coordinates with planar tests (cot-projection.md §3.2: fine for AOIs of a
// few tens of kilometres; larger or polar AOIs need a geodesic test). The polygon and circle
// tests follow services/layer/internal/layer/geo.go.

import (
	"encoding/json"
	"errors"
	"math"
)

// pt is [lon, lat] (GeoJSON order).
type pt [2]float64

// geoPoints returns every vertex of a geo component (Point, MultiPoint, LineString, Polygon
// outer ring, MultiPolygon outer rings), in [lon, lat].
func geoPoints(g *Geo) []pt {
	if g == nil || len(g.Coords) == 0 {
		return nil
	}
	var gj struct {
		Type        string          `json:"type"`
		Coordinates json.RawMessage `json:"coordinates"`
	}
	if json.Unmarshal(g.Coords, &gj) != nil {
		return nil
	}
	conv := func(in [][]float64) []pt {
		out := make([]pt, 0, len(in))
		for _, p := range in {
			if len(p) >= 2 && validLL(p[1], p[0]) {
				out = append(out, pt{p[0], p[1]})
			}
		}
		return out
	}
	switch gj.Type {
	case "Point":
		var p []float64
		if json.Unmarshal(gj.Coordinates, &p) == nil && len(p) >= 2 && validLL(p[1], p[0]) {
			return []pt{{p[0], p[1]}}
		}
	case "MultiPoint", "LineString":
		var ps [][]float64
		if json.Unmarshal(gj.Coordinates, &ps) == nil {
			return conv(ps)
		}
	case "Polygon":
		var rings [][][]float64
		if json.Unmarshal(gj.Coordinates, &rings) == nil && len(rings) > 0 {
			return conv(rings[0])
		}
	case "MultiPolygon":
		var polys [][][][]float64
		if json.Unmarshal(gj.Coordinates, &polys) == nil {
			var out []pt
			for _, p := range polys {
				if len(p) > 0 {
					out = append(out, conv(p[0])...)
				}
			}
			return out
		}
	}
	return nil
}

func geoType(g *Geo) string {
	if g == nil {
		return ""
	}
	var gj struct {
		Type string `json:"type"`
	}
	_ = json.Unmarshal(g.Coords, &gj)
	return gj.Type
}

func centroid(ps []pt) (lat, lon float64, ok bool) {
	if len(ps) == 0 {
		return 0, 0, false
	}
	n := len(ps)
	if n > 1 && ps[0] == ps[n-1] {
		n-- // closed ring: do not count the first vertex twice
	}
	var sLat, sLon float64
	for _, p := range ps[:n] {
		sLat += p[1]
		sLon += p[0]
	}
	return sLat / float64(n), sLon / float64(n), true
}

// area is a resolved AOI: a polygon ring or a circle, with its bounding box.
type area struct {
	minLon, minLat, maxLon, maxLat float64
	ring                           []pt
	circle                         *[3]float64 // lon, lat, radius metres
	altMin, altMax                 *float64
	cLat, cLon                     float64 // centroid, for the cap's distance order
}

var errNotArea = errors.New("geo is not an area (want Polygon, MultiPolygon, Point with radius_m, or a corridor LineString with radius_m)")

// geoArea turns a zone's geo into an area.
func geoArea(g *Geo) (*area, error) {
	if g == nil {
		return nil, errNotArea
	}
	pts := geoPoints(g)
	var a *area
	switch geoType(g) {
	case "Polygon", "MultiPolygon":
		if len(pts) < 3 {
			return nil, errNotArea
		}
		a = ringArea(pts)
	case "Point":
		if len(pts) != 1 || g.RadiusM <= 0 {
			return nil, errNotArea
		}
		p := pts[0]
		dLat := g.RadiusM / 111320.0
		dLon := dLat / math.Max(math.Cos(p[1]*math.Pi/180), 1e-6)
		a = &area{minLon: p[0] - dLon, minLat: p[1] - dLat, maxLon: p[0] + dLon, maxLat: p[1] + dLat,
			circle: &[3]float64{p[0], p[1], g.RadiusM}, cLat: p[1], cLon: p[0]}
	case "LineString":
		if len(pts) < 2 || g.RadiusM <= 0 {
			return nil, errNotArea
		}
		a = ringArea(bufferLine(pts, g.RadiusM))
	default:
		return nil, errNotArea
	}
	a.altMin, a.altMax = g.AltMinM, g.AltMaxM
	return a, nil
}

func ringArea(ring []pt) *area {
	a := &area{minLon: 180, minLat: 90, maxLon: -180, maxLat: -90, ring: ring}
	for _, p := range ring {
		a.minLon, a.maxLon = math.Min(a.minLon, p[0]), math.Max(a.maxLon, p[0])
		a.minLat, a.maxLat = math.Min(a.minLat, p[1]), math.Max(a.maxLat, p[1])
	}
	a.cLat, a.cLon, _ = centroid(ring)
	return a
}

// contains tests a point; the altitude band applies only when the altitude is known.
func (a *area) contains(lat, lon float64, alt float64, hasAlt bool) bool {
	if lat < a.minLat || lat > a.maxLat || lon < a.minLon || lon > a.maxLon {
		return false
	}
	if hasAlt && ((a.altMin != nil && alt < *a.altMin) || (a.altMax != nil && alt > *a.altMax)) {
		return false
	}
	if a.circle != nil {
		return haversineM(lat, lon, a.circle[1], a.circle[0]) <= a.circle[2]
	}
	return pointInRing(lon, lat, a.ring)
}

func pointInRing(x, y float64, ring []pt) bool {
	in := false
	for i, j := 0, len(ring)-1; i < len(ring); j, i = i, i+1 {
		xi, yi, xj, yj := ring[i][0], ring[i][1], ring[j][0], ring[j][1]
		if (yi > y) != (yj > y) && x < (xj-xi)*(y-yi)/(yj-yi)+xi {
			in = !in
		}
	}
	return in
}

func haversineM(lat1, lon1, lat2, lon2 float64) float64 {
	const r = 6371008.8
	p1, p2 := lat1*math.Pi/180, lat2*math.Pi/180
	dp, dl := (lat2-lat1)*math.Pi/180, (lon2-lon1)*math.Pi/180
	h := math.Sin(dp/2)*math.Sin(dp/2) + math.Cos(p1)*math.Cos(p2)*math.Sin(dl/2)*math.Sin(dl/2)
	return 2 * r * math.Asin(math.Min(1, math.Sqrt(h)))
}

// bufferLine buffers a polyline by r metres into a closed polygon (left side forward, right
// side back; joins are mitred by the averaged normal). Used for corridors (§5.5).
func bufferLine(line []pt, r float64) []pt {
	n := len(line)
	left := make([]pt, n)
	right := make([]pt, n)
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
		ex, ey := dx*cosLat, dy // local metric-ish direction
		l := math.Hypot(ex, ey)
		if l == 0 {
			l = 1
		}
		nx, ny := -ey/l, ex/l // unit normal (left)
		dLat := ny * r / 111320.0
		dLon := nx * r / 111320.0 / cosLat
		left[i] = pt{line[i][0] + dLon, line[i][1] + dLat}
		right[i] = pt{line[i][0] - dLon, line[i][1] - dLat}
	}
	out := append([]pt{}, left...)
	for i := n - 1; i >= 0; i-- {
		out = append(out, right[i])
	}
	return append(out, left[0])
}

// simplify reduces a ring or line to at most max vertices with Douglas-Peucker, raising the
// tolerance until it fits (§5.5 MAX_VERTICES).
func simplify(ps []pt, max int) []pt {
	if max < 4 || len(ps) <= max {
		return ps
	}
	tol := 1e-6
	out := ps
	for i := 0; i < 40 && len(out) > max; i++ {
		out = dp(ps, tol)
		tol *= 2
	}
	if len(out) > max { // degenerate input: fall back to even sampling
		step := float64(len(ps)-1) / float64(max-1)
		out = make([]pt, 0, max)
		for i := 0; i < max; i++ {
			out = append(out, ps[int(math.Round(float64(i)*step))])
		}
	}
	return out
}

func dp(ps []pt, tol float64) []pt {
	if len(ps) < 3 {
		return ps
	}
	a, b := ps[0], ps[len(ps)-1]
	idx, dmax := 0, 0.0
	for i := 1; i < len(ps)-1; i++ {
		if d := segDist(ps[i], a, b); d > dmax {
			idx, dmax = i, d
		}
	}
	if dmax <= tol {
		return []pt{a, b}
	}
	l := dp(ps[:idx+1], tol)
	r := dp(ps[idx:], tol)
	return append(l[:len(l)-1], r...)
}

func segDist(p, a, b pt) float64 {
	dx, dy := b[0]-a[0], b[1]-a[1]
	if dx == 0 && dy == 0 {
		return math.Hypot(p[0]-a[0], p[1]-a[1])
	}
	t := ((p[0]-a[0])*dx + (p[1]-a[1])*dy) / (dx*dx + dy*dy)
	t = math.Max(0, math.Min(1, t))
	return math.Hypot(p[0]-(a[0]+t*dx), p[1]-(a[1]+t*dy))
}
