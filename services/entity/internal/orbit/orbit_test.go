package orbit

import (
	"encoding/json"
	"math"
	"os"
	"strings"
	"testing"
	"time"
)

func loadTLEs(t *testing.T) map[string][2]string {
	t.Helper()
	b, err := os.ReadFile("testdata/stations.tle")
	if err != nil {
		t.Fatal(err)
	}
	var lines []string
	for _, l := range strings.Split(string(b), "\n") {
		if l = strings.TrimRight(l, " \r"); l != "" {
			lines = append(lines, l)
		}
	}
	out := map[string][2]string{}
	for i := 0; i+1 < len(lines); i++ {
		if strings.HasPrefix(lines[i], "1 ") && strings.HasPrefix(lines[i+1], "2 ") {
			out[strings.TrimSpace(lines[i][2:7])] = [2]string{lines[i], lines[i+1]}
		}
	}
	return out
}

// TestMatchesIndependentSGP4 is the FEATURES A12 acceptance check: every propagated position
// (ISS included) is within 5 km of python-sgp4 (testdata/gen_reference.py), out to 6 days.
func TestMatchesIndependentSGP4(t *testing.T) {
	tles := loadTLEs(t)
	b, err := os.ReadFile("testdata/reference.json")
	if err != nil {
		t.Fatal(err)
	}
	var ref []struct {
		NORAD string  `json:"norad"`
		Query string  `json:"query"`
		TMs   int64   `json:"t_ms"`
		Lat   float64 `json:"lat"`
		Lon   float64 `json:"lon"`
		AltKm float64 `json:"alt_km"`
	}
	if err := json.Unmarshal(b, &ref); err != nil {
		t.Fatal(err)
	}
	if len(ref) < 100 {
		t.Fatalf("reference has %d rows", len(ref))
	}
	var worst float64
	iss, compared := 0, 0
	for _, r := range ref {
		l := tles[r.NORAD]
		el, err := Parse(l[0], l[1])
		if r.NORAD == "99999" { // the fixture's deliberate bad-checksum record (#558)
			if err == nil {
				t.Fatal("bad-checksum fixture TLE accepted")
			}
			continue
		}
		if err != nil {
			t.Fatalf("%s: %v", r.NORAD, err)
		}
		q, _ := time.Parse(time.RFC3339Nano, r.Query)
		s, err := el.At(q)
		if age := q.Sub(el.Epoch); age > MaxAge || age < -MaxAge {
			if err == nil {
				t.Fatalf("%s at %s: propagated past MaxAge", r.NORAD, r.Query)
			}
			continue
		}
		if err != nil {
			t.Fatalf("%s at %s: %v", r.NORAD, r.Query, err)
		}
		if s.T.UnixMilli() != r.TMs {
			t.Fatalf("%s: propagated to %d, reference %d", r.NORAD, s.T.UnixMilli(), r.TMs)
		}
		compared++
		h := haversineM(s.Lat, s.Lon, r.Lat, r.Lon)
		d := math.Hypot(h, s.AltM-r.AltKm*1000)
		worst = math.Max(worst, d)
		if d > 5000 {
			t.Errorf("%s at %s: %.0f m from reference", r.NORAD, r.Query, d)
		}
		if r.NORAD == "25544" {
			iss++
			if s.AltM < 380e3 || s.AltM > 450e3 || s.SpeedMps < 6500 || s.SpeedMps > 7700 {
				t.Errorf("ISS state implausible: %+v", s)
			}
		}
	}
	if iss == 0 {
		t.Fatal("no ISS rows")
	}
	t.Logf("%d positions compared, worst 3-D difference %.1f m", compared, worst)
}

func TestParseRejectsBadTLEWithoutExiting(t *testing.T) {
	good := loadTLEs(t)["25544"]
	flip := func(s string, i int, c byte) string { b := []byte(s); b[i] = c; return string(b) }
	cases := map[string][2]string{
		"short line 1":      {good[0][:60], good[1]},
		"checksum":          {flip(good[0], 68, '0'+(good[0][68]-'0'+1)%10), good[1]},
		"line number":       {flip(good[0], 0, '3'), good[1]},
		"norad mismatch":    {good[0], flip(flip(good[1], 6, '5'), 68, good[1][68]+1)},
		"non-numeric field": {fixSum(flip(good[0], 22, 'x')), good[1]},
		"bad eccentricity":  {good[0], fixSum(flip(good[1], 28, 'Z'))},
	}
	for name, c := range cases {
		if _, err := Parse(c[0], c[1]); err == nil {
			t.Errorf("%s: accepted", name)
		}
	}
	el, err := Parse(good[0], good[1])
	if err != nil {
		t.Fatal(err)
	}
	if _, err := el.At(el.Epoch.Add(8 * 24 * time.Hour)); err == nil {
		t.Error("propagated a TLE 8 days past its epoch")
	}
}

func TestEpochAndTruncation(t *testing.T) {
	el, err := Parse(loadTLEs(t)["25544"][0], loadTLEs(t)["25544"][1])
	if err != nil {
		t.Fatal(err)
	}
	// 26301.59601386 = 2026 day 301 (28 Oct) 14:18:15.597504 UTC.
	if got := el.Epoch.Format("2006-01-02T15:04:05.000Z"); got != "2026-10-28T14:18:15.598Z" {
		t.Fatalf("epoch %s", got)
	}
	if got := el.truncEpoch.Format(time.RFC3339); got != "2026-10-28T14:18:15Z" {
		t.Fatalf("library epoch %s", got)
	}
	s, _ := el.At(time.Date(2026, 10, 29, 12, 0, 0, 0, time.UTC))
	if s.T.Sub(el.Epoch)%time.Second != 0 {
		t.Fatalf("state time %s is not a whole second from epoch", s.T)
	}
}

func fixSum(l string) string {
	sum := 0
	for _, c := range l[:68] {
		if c >= '0' && c <= '9' {
			sum += int(c - '0')
		} else if c == '-' {
			sum++
		}
	}
	return l[:68] + string(rune('0'+sum%10))
}
