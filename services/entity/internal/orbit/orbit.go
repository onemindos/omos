// Package orbit propagates CelesTrak two-line element sets with SGP4 (FEATURES A12).
//
// It wraps github.com/joshuaferrara/go-satellite (BSD-2-Clause, a Go port of Vallado's SGP4)
// and works around three behaviours of that library:
//
//   - its TLE parser calls log.Fatal on a malformed number, so Parse validates every field the
//     library reads (length, checksum, line numbers, NORAD id, numeric columns) before calling it;
//   - it truncates the TLE epoch to whole seconds, so At only propagates to instants that are a
//     whole number of seconds from the true epoch and passes the library the same offset from
//     its truncated epoch (time-since-epoch is then exact);
//   - its LatLongDeg wraps longitudes wrongly and calls log.Fatal, so the TEME to geodetic
//     conversion (IAU-82 GMST, WGS84 ellipsoid) is done here.
package orbit

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	satellite "github.com/joshuaferrara/go-satellite"
)

// MaxAge is how far from its epoch a TLE is propagated. SGP4 error grows by kilometres per day
// for LEO objects; CelesTrak refreshes the stations group several times a day.
const MaxAge = 7 * 24 * time.Hour

// Elements is one parsed and initialised TLE.
type Elements struct {
	NORAD string
	Line1 string
	Line2 string
	Epoch time.Time // true epoch, millisecond precision

	sat        satellite.Satellite
	truncEpoch time.Time // the epoch the library uses (whole seconds)
}

// State is a propagated position.
type State struct {
	T         time.Time // the instant propagated to (a whole number of seconds from Epoch)
	Lat, Lon  float64   // geodetic degrees, WGS84
	AltM      float64   // height above the ellipsoid, metres
	SpeedMps  float64   // ground-track speed of the sub-satellite point
	CourseDeg float64   // ground-track course, degrees true
}

// Parse validates and initialises a TLE.
func Parse(line1, line2 string) (*Elements, error) {
	l1 := strings.TrimRight(line1, " \r\n")
	l2 := strings.TrimRight(line2, " \r\n")
	for i, l := range []string{l1, l2} {
		if len(l) != 69 {
			return nil, fmt.Errorf("tle line %d: length %d, want 69", i+1, len(l))
		}
		if l[0] != byte('1'+i) || l[1] != ' ' {
			return nil, fmt.Errorf("tle line %d: bad line number", i+1)
		}
		if !checksumOK(l) {
			return nil, fmt.Errorf("tle line %d: checksum", i+1)
		}
	}
	norad := strings.TrimSpace(l1[2:7])
	if norad != strings.TrimSpace(l2[2:7]) {
		return nil, errors.New("tle: NORAD id differs between lines")
	}
	if err := fieldsOK(l1, l2); err != nil {
		return nil, err
	}
	yy, _ := strconv.Atoi(l1[18:20])
	days, _ := strconv.ParseFloat(l1[20:32], 64)
	year := 1900 + yy
	if yy < 57 {
		year = 2000 + yy
	}
	if days < 1 || days >= 367 {
		return nil, fmt.Errorf("tle: epoch day %v", days)
	}
	jan1 := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)
	epoch := jan1.Add(time.Duration(math.Round((days-1)*86400e3)) * time.Millisecond)

	sat := satellite.TLEToSat(l1, l2, satellite.GravityWGS72)
	if sat.Error != 0 {
		return nil, fmt.Errorf("tle: sgp4 init error %d", sat.Error)
	}
	return &Elements{NORAD: norad, Line1: l1, Line2: l2, Epoch: epoch, sat: sat,
		truncEpoch: libraryEpoch(year, days)}, nil
}

// At propagates to the whole second from Epoch nearest to t.
func (e *Elements) At(t time.Time) (State, error) {
	off := t.Sub(e.Epoch).Round(time.Second)
	if off > MaxAge || off < -MaxAge {
		return State{}, fmt.Errorf("tle epoch %s is %s from %s (max %s)", e.Epoch.Format(time.RFC3339), off, t.UTC().Format(time.RFC3339), MaxAge)
	}
	p0, err := e.subpoint(off)
	if err != nil {
		return State{}, err
	}
	p1, err := e.subpoint(off + time.Second)
	if err != nil {
		return State{}, err
	}
	p0.SpeedMps = haversineM(p0.Lat, p0.Lon, p1.Lat, p1.Lon)
	p0.CourseDeg = bearingDeg(p0.Lat, p0.Lon, p1.Lat, p1.Lon)
	return p0, nil
}

func (e *Elements) subpoint(off time.Duration) (State, error) {
	tl := e.truncEpoch.Add(off)
	pos, _ := satellite.Propagate(e.sat, tl.Year(), int(tl.Month()), tl.Day(), tl.Hour(), tl.Minute(), tl.Second())
	r := math.Sqrt(pos.X*pos.X + pos.Y*pos.Y + pos.Z*pos.Z)
	if math.IsNaN(r) || math.IsInf(r, 0) || r < 6378.137 || r > 500000 {
		return State{}, fmt.Errorf("sgp4: no valid position (|r| = %v km); orbit decayed or elements invalid", r)
	}
	t := e.Epoch.Add(off)
	lat, lon, alt := temeToGeodetic(pos.X, pos.Y, pos.Z, gmst(julian(t)))
	return State{T: t, Lat: lat, Lon: lon, AltM: alt * 1000}, nil
}

// libraryEpoch reproduces the library's epoch: day-of-year to month, day, hour, minute and
// whole seconds (go-satellite days2mdhms + JDay drop the fractional second).
func libraryEpoch(year int, epochDays float64) time.Time {
	months := [12]int{31, 28, 31, 30, 31, 30, 31, 31, 30, 31, 30, 31}
	if year%4 == 0 {
		months[1] = 29
	}
	doy := math.Floor(epochDays)
	mon, acc := 1, 0.0
	for mon < 12 && doy > acc+float64(months[mon-1]) {
		acc += float64(months[mon-1])
		mon++
	}
	day := doy - acc
	h := (epochDays - doy) * 24
	hr := math.Floor(h)
	m := (h - hr) * 60
	mn := math.Floor(m)
	sec := (m - mn) * 60
	return time.Date(year, time.Month(mon), int(day), int(hr), int(mn), int(sec), 0, time.UTC)
}

// checksumOK is the TLE modulo-10 checksum: digits count their value, '-' counts 1.
func checksumOK(l string) bool {
	sum := 0
	for _, c := range l[:68] {
		switch {
		case c >= '0' && c <= '9':
			sum += int(c - '0')
		case c == '-':
			sum++
		}
	}
	return l[68] >= '0' && l[68] <= '9' && sum%10 == int(l[68]-'0')
}

// fieldsOK parses every column exactly as go-satellite's ParseTLE does, so that the library
// never reaches its log.Fatal path.
func fieldsOK(l1, l2 string) error {
	ints := []string{strings.TrimSpace(l1[2:7]), l1[18:20]}
	floats := []string{
		l1[20:32],
		strings.Replace(l1[33:43], " ", "", 2),
		strings.Replace(l1[44:45]+"."+l1[45:50]+"e"+l1[50:52], " ", "", 2),
		strings.Replace(l1[53:54]+"."+l1[54:59]+"e"+l1[59:61], " ", "", 2),
		strings.Replace(l2[8:16], " ", "", 2),
		strings.Replace(l2[17:25], " ", "", 2),
		"." + l2[26:33],
		strings.Replace(l2[34:42], " ", "", 2),
		strings.Replace(l2[43:51], " ", "", 2),
		strings.Replace(l2[52:63], " ", "", 2),
	}
	for _, s := range ints {
		if _, err := strconv.ParseInt(s, 10, 0); err != nil {
			return fmt.Errorf("tle field %q: %w", s, err)
		}
	}
	for _, s := range floats {
		if _, err := strconv.ParseFloat(s, 64); err != nil {
			return fmt.Errorf("tle field %q: %w", s, err)
		}
	}
	return nil
}

func julian(t time.Time) float64 {
	return float64(t.UnixMilli())/86400e3 + 2440587.5
}

// gmst is IAU-82 Greenwich mean sidereal time in radians (Vallado, Fundamentals of
// Astrodynamics, eq. 3-45), the frame SGP4's TEME output is rotated by.
func gmst(jd float64) float64 {
	tut1 := (jd - 2451545.0) / 36525.0
	sec := -6.2e-6*tut1*tut1*tut1 + 0.093104*tut1*tut1 + (876600.0*3600+8640184.812866)*tut1 + 67310.54841
	g := math.Mod(sec*math.Pi/180/240, 2*math.Pi)
	if g < 0 {
		g += 2 * math.Pi
	}
	return g
}

// temeToGeodetic rotates TEME (km) to Earth-fixed by GMST and converts to WGS84 geodetic.
func temeToGeodetic(x, y, z, theta float64) (latDeg, lonDeg, altKm float64) {
	const a, f = 6378.137, 1 / 298.257223563
	e2 := f * (2 - f)
	c, s := math.Cos(theta), math.Sin(theta)
	xe, ye := c*x+s*y, -s*x+c*y
	p := math.Hypot(xe, ye)
	lon := math.Atan2(ye, xe)
	lat := math.Atan2(z, p*(1-e2))
	var n float64
	for i := 0; i < 10; i++ {
		sl := math.Sin(lat)
		n = a / math.Sqrt(1-e2*sl*sl)
		lat = math.Atan2(z+e2*n*sl, p)
	}
	sl := math.Sin(lat)
	n = a / math.Sqrt(1-e2*sl*sl)
	if cl := math.Cos(lat); math.Abs(cl) > 1e-9 {
		altKm = p/cl - n
	} else {
		altKm = math.Abs(z) - n*(1-e2)
	}
	return lat * 180 / math.Pi, lon * 180 / math.Pi, altKm
}

const earthR = 6371008.8

func haversineM(lat1, lon1, lat2, lon2 float64) float64 {
	r := math.Pi / 180
	dl, dn := (lat2-lat1)*r, (lon2-lon1)*r
	h := math.Sin(dl/2)*math.Sin(dl/2) + math.Cos(lat1*r)*math.Cos(lat2*r)*math.Sin(dn/2)*math.Sin(dn/2)
	return 2 * earthR * math.Asin(math.Min(1, math.Sqrt(h)))
}

func bearingDeg(lat1, lon1, lat2, lon2 float64) float64 {
	r := math.Pi / 180
	y := math.Sin((lon2-lon1)*r) * math.Cos(lat2*r)
	x := math.Cos(lat1*r)*math.Sin(lat2*r) - math.Sin(lat1*r)*math.Cos(lat2*r)*math.Cos((lon2-lon1)*r)
	return math.Mod(math.Atan2(y, x)/r+360, 360)
}
