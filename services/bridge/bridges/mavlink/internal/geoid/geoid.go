// Package geoid converts MAVLink MSL altitudes to WGS-84 ellipsoid heights (HAE).
//
// The doctrine geo-envelope field `hae` is metres above the WGS-84 ellipsoid
// (docs/NATS-DOCTRINE.md §9.2), while GLOBAL_POSITION_INT.alt is millimetres above
// mean sea level. HAE = MSL + N, where N is the geoid undulation.
//
// N comes from an embedded EGM96 grid at 0.5 degree spacing, derived from the
// GeographicLib 15-minute EGM96 grid by gen_egm96_grid.py. Bilinear interpolation of
// the 0.5 degree grid differs from the 15-minute grid by at most 1.95 m (RMS 0.09 m),
// measured by that script over every 15-minute cell centre.
package geoid

import (
	"bytes"
	"compress/gzip"
	_ "embed"
	"encoding/binary"
	"fmt"
	"io"
	"math"
	"sync"
)

//go:embed egm96-30min.bin.gz
var gridGz []byte

const (
	rows   = 361 // 90N .. 90S
	cols   = 720 // 0E .. 359.5E
	step   = 0.5
	offset = -108.0
	scale  = 0.003
)

var (
	once    sync.Once
	grid    []uint16
	loadErr error
)

func load() {
	zr, err := gzip.NewReader(bytes.NewReader(gridGz))
	if err != nil {
		loadErr = err
		return
	}
	raw, err := io.ReadAll(zr)
	if err != nil {
		loadErr = err
		return
	}
	if len(raw) != rows*cols*2 {
		loadErr = fmt.Errorf("geoid grid: %d bytes, want %d", len(raw), rows*cols*2)
		return
	}
	g := make([]uint16, rows*cols)
	for i := range g {
		g[i] = binary.BigEndian.Uint16(raw[i*2:])
	}
	grid = g
}

func at(r, c int) float64 {
	return offset + scale*float64(grid[r*cols+((c%cols)+cols)%cols])
}

// Undulation returns the EGM96 geoid height N in metres at a WGS-84 position.
func Undulation(lat, lon float64) (float64, error) {
	once.Do(load)
	if loadErr != nil {
		return 0, loadErr
	}
	if math.IsNaN(lat) || math.IsNaN(lon) || lat < -90 || lat > 90 {
		return 0, fmt.Errorf("geoid: invalid position %v,%v", lat, lon)
	}
	y := (90.0 - lat) / step
	x := math.Mod(math.Mod(lon, 360.0)+360.0, 360.0) / step
	r0 := int(math.Floor(y))
	if r0 >= rows-1 {
		r0 = rows - 2
	}
	c0 := int(math.Floor(x))
	fy, fx := y-float64(r0), x-float64(c0)
	a, b := at(r0, c0), at(r0, c0+1)
	c, d := at(r0+1, c0), at(r0+1, c0+1)
	return (1-fy)*((1-fx)*a+fx*b) + fy*((1-fx)*c+fx*d), nil
}

// MSLToHAE converts metres above mean sea level to metres above the ellipsoid.
func MSLToHAE(lat, lon, msl float64) (float64, error) {
	n, err := Undulation(lat, lon)
	return msl + n, err
}

// HAEToMSL converts metres above the ellipsoid to metres above mean sea level.
func HAEToMSL(lat, lon, hae float64) (float64, error) {
	n, err := Undulation(lat, lon)
	return hae - n, err
}
