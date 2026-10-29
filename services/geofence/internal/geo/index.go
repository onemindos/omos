package geo

import (
	"math"
	"sort"
)

// Index is the in-process spatial index: a uniform lon/lat grid that maps each cell to the
// zones whose bounding box touches it. A lookup reads one cell, then runs the exact test on
// those zones only. Zones whose box would cover more than maxCells go to a short list that
// every lookup checks. Zone counts are small (layer-snapshot.md §1: the `zones` layer stays
// under 1,000), so this is enough; it needs no extra process, unlike Tile38.
type Index struct {
	cellDeg  float64
	maxCells int
	cells    map[[2]int32][]string
	big      map[string]bool
	areas    map[string]*Area
	keys     map[string][][2]int32
}

// NewIndex returns an index with cells of cellDeg degrees (default 0.1, about 11 km).
func NewIndex(cellDeg float64) *Index {
	if cellDeg <= 0 {
		cellDeg = 0.1
	}
	return &Index{cellDeg: cellDeg, maxCells: 4096, cells: map[[2]int32][]string{},
		big: map[string]bool{}, areas: map[string]*Area{}, keys: map[string][][2]int32{}}
}

func (x *Index) cell(lat, lon float64) [2]int32 {
	return [2]int32{int32(math.Floor(lon / x.cellDeg)), int32(math.Floor(lat / x.cellDeg))}
}

// Put adds or replaces a zone.
func (x *Index) Put(id string, a *Area) {
	x.Remove(id)
	x.areas[id] = a
	lo, hi := x.cell(a.MinLat, a.MinLon), x.cell(a.MaxLat, a.MaxLon)
	n := int64(hi[0]-lo[0]+1) * int64(hi[1]-lo[1]+1)
	if n > int64(x.maxCells) {
		x.big[id] = true
		return
	}
	var ks [][2]int32
	for cx := lo[0]; cx <= hi[0]; cx++ {
		for cy := lo[1]; cy <= hi[1]; cy++ {
			k := [2]int32{cx, cy}
			x.cells[k] = append(x.cells[k], id)
			ks = append(ks, k)
		}
	}
	x.keys[id] = ks
}

// Remove deletes a zone.
func (x *Index) Remove(id string) {
	if _, ok := x.areas[id]; !ok {
		return
	}
	delete(x.areas, id)
	delete(x.big, id)
	for _, k := range x.keys[id] {
		ids := x.cells[k]
		for i, v := range ids {
			if v == id {
				ids = append(ids[:i], ids[i+1:]...)
				break
			}
		}
		if len(ids) == 0 {
			delete(x.cells, k)
		} else {
			x.cells[k] = ids
		}
	}
	delete(x.keys, id)
}

// Containing returns the ids of the zones that contain the point, sorted.
func (x *Index) Containing(lat, lon, alt float64, hasAlt bool) []string {
	var out []string
	for _, id := range x.cells[x.cell(lat, lon)] {
		if x.areas[id].Contains(lat, lon, alt, hasAlt) {
			out = append(out, id)
		}
	}
	for id := range x.big {
		if x.areas[id].Contains(lat, lon, alt, hasAlt) {
			out = append(out, id)
		}
	}
	sort.Strings(out)
	return out
}

// Len is the number of zones.
func (x *Index) Len() int { return len(x.areas) }
