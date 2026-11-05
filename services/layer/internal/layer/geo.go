package layer

import (
	"encoding/json"
	"errors"
	"math"
	"sort"
)

// AOIReq is the `aoi` of a snapshot or watch request: one of bbox, polygon or zone_id.
type AOIReq struct {
	BBox    []float64   `json:"bbox,omitempty"`    // [minLon, minLat, maxLon, maxLat]
	Polygon [][]float64 `json:"polygon,omitempty"` // [[lon,lat],...]
	ZoneID  string      `json:"zone_id,omitempty"`
}

// area is a resolved AOI. A nil *area means "whole world".
type area struct {
	minLon, minLat, maxLon, maxLat float64
	wrap                           bool         // bbox crosses the antimeridian (minLon > maxLon)
	poly                           [][2]float64 // optional polygon ring, tested after the bbox
	circle                         *[3]float64  // lon, lat, radius metres
}

// MaxPolygonPoints bounds the per-entity cost of a polygon AOI.
const MaxPolygonPoints = 1000

var errBadAOI = errors.New("aoi must be {bbox:[minLon,minLat,maxLon,maxLat]}, {polygon:[[lon,lat],...] with 3 to 1000 points} or {zone_id}")

func validLL(lon, lat float64) bool {
	return lon >= -180 && lon <= 180 && lat >= -90 && lat <= 90 && !math.IsNaN(lon) && !math.IsNaN(lat)
}

// bboxArea builds an area from a bbox.
func bboxArea(b []float64) (*area, error) {
	if len(b) != 4 || !validLL(b[0], b[1]) || !validLL(b[2], b[3]) || b[1] > b[3] {
		return nil, errBadAOI
	}
	return &area{minLon: b[0], minLat: b[1], maxLon: b[2], maxLat: b[3], wrap: b[0] > b[2]}, nil
}

// polyArea builds an area from a ring of [lon,lat] points (bbox of the ring plus the ring).
func polyArea(pts [][]float64) (*area, error) {
	if len(pts) < 3 || len(pts) > MaxPolygonPoints {
		return nil, errBadAOI
	}
	a := &area{minLon: 180, minLat: 90, maxLon: -180, maxLat: -90}
	for _, p := range pts {
		if len(p) < 2 || !validLL(p[0], p[1]) {
			return nil, errBadAOI
		}
		a.poly = append(a.poly, [2]float64{p[0], p[1]})
		a.minLon, a.maxLon = math.Min(a.minLon, p[0]), math.Max(a.maxLon, p[0])
		a.minLat, a.maxLat = math.Min(a.minLat, p[1]), math.Max(a.maxLat, p[1])
	}
	return a, nil
}

// zoneArea builds an area from a zone entity's geo component (GeoJSON Polygon, or Point with
// radius_m for a circle).
func zoneArea(g *GeoIn) (*area, error) {
	if g == nil || len(g.Coords) == 0 {
		return nil, errors.New("zone has no geo")
	}
	var gj struct {
		Type        string          `json:"type"`
		Coordinates json.RawMessage `json:"coordinates"`
	}
	if err := json.Unmarshal(g.Coords, &gj); err != nil {
		return nil, err
	}
	switch gj.Type {
	case "Polygon":
		var rings [][][]float64
		if err := json.Unmarshal(gj.Coordinates, &rings); err != nil || len(rings) == 0 {
			return nil, errors.New("zone polygon unreadable")
		}
		return polyArea(rings[0])
	case "Point":
		var p []float64
		if err := json.Unmarshal(gj.Coordinates, &p); err != nil || len(p) < 2 || g.RadiusM <= 0 {
			return nil, errors.New("zone circle needs a Point and radius_m")
		}
		dLat := g.RadiusM / 111320.0
		dLon := dLat / math.Max(math.Cos(p[1]*math.Pi/180), 1e-6)
		return &area{minLon: math.Max(-180, p[0]-dLon), minLat: math.Max(-90, p[1]-dLat),
			maxLon: math.Min(180, p[0]+dLon), maxLat: math.Min(90, p[1]+dLat),
			circle: &[3]float64{p[0], p[1], g.RadiusM}}, nil
	}
	return nil, errors.New("zone geo type " + gj.Type + " is not an area")
}

// contains reports whether lat/lon is inside the area. A nil area contains everything.
func (a *area) contains(lat, lon float64) bool {
	if a == nil {
		return true
	}
	if lat < a.minLat || lat > a.maxLat {
		return false
	}
	if a.wrap {
		if lon < a.minLon && lon > a.maxLon {
			return false
		}
	} else if lon < a.minLon || lon > a.maxLon {
		return false
	}
	if a.circle != nil {
		return haversineM(lat, lon, a.circle[1], a.circle[0]) <= a.circle[2]
	}
	if len(a.poly) > 0 {
		return pointInRing(lon, lat, a.poly)
	}
	return true
}

func pointInRing(x, y float64, ring [][2]float64) bool {
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

// ---------------------------------------------------------------- density cells

// tileXY returns the web-mercator XYZ tile of lat/lon at zoom z.
func tileXY(lat, lon float64, z int) (int, int) {
	const maxLat = 85.05112878
	lat = math.Max(-maxLat, math.Min(maxLat, lat))
	n := float64(int(1) << uint(z))
	x := int(math.Floor((lon + 180) / 360 * n))
	lr := lat * math.Pi / 180
	y := int(math.Floor((1 - math.Log(math.Tan(lr)+1/math.Cos(lr))/math.Pi) / 2 * n))
	maxI := int(n) - 1
	if x < 0 {
		x = 0
	} else if x > maxI {
		x = maxI
	}
	if y < 0 {
		y = 0
	} else if y > maxI {
		y = maxI
	}
	return x, y
}

type cellAcc struct {
	n          int
	sLat, sLon float64
}

// cellGrid accumulates entities into cells at one zoom.
type cellGrid struct {
	z     int
	cells map[[2]int]*cellAcc
}

func newCellGrid(z int) *cellGrid { return &cellGrid{z: z, cells: map[[2]int]*cellAcc{}} }

func (g *cellGrid) add(lat, lon float64) {
	x, y := tileXY(lat, lon, g.z)
	c := g.cells[[2]int{x, y}]
	if c == nil {
		c = &cellAcc{}
		g.cells[[2]int{x, y}] = c
	}
	c.n++
	c.sLat += lat
	c.sLon += lon
}

func round5(v float64) float64 { return math.Round(v*1e5) / 1e5 }

// cell is one [x, y, count, lat_centroid, lon_centroid] record.
type cell [5]float64

// list returns the cells sorted by x, y.
func (g *cellGrid) list() []cell {
	out := make([]cell, 0, len(g.cells))
	for k, c := range g.cells {
		out = append(out, cell{float64(k[0]), float64(k[1]), float64(c.n), round5(c.sLat / float64(c.n)), round5(c.sLon / float64(c.n))})
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i][0] != out[j][0] {
			return out[i][0] < out[j][0]
		}
		return out[i][1] < out[j][1]
	})
	return out
}

// asMap returns the cells keyed by x, y.
func (g *cellGrid) asMap() map[[2]int]cell {
	out := make(map[[2]int]cell, len(g.cells))
	for _, c := range g.list() {
		out[[2]int{int(c[0]), int(c[1])}] = c
	}
	return out
}
