package fence

import "testing"

const doc = `{
  "keep_in": [{"id": "sim-box", "polygon": [
    {"lat": -0.01, "lon": -0.01}, {"lat": -0.01, "lon": 0.01},
    {"lat": 0.01, "lon": 0.01}, {"lat": 0.01, "lon": -0.01}]}],
  "keep_out": [{"id": "no-fly", "circle": {"lat": 0.005, "lon": 0.005, "radius_m": 100}}],
  "max_hae_m": 200
}`

func TestCheck(t *testing.T) {
	f, err := Parse([]byte(doc))
	if err != nil {
		t.Fatal(err)
	}
	h := func(v float64) *float64 { return &v }
	cases := []struct {
		name   string
		p      Point
		hae    *float64
		detail string
	}{
		{"inside", Point{0, 0}, h(100), ""},
		{"inside no alt", Point{0.002, -0.002}, nil, ""},
		{"outside keep-in", Point{0.02, 0}, h(100), "keep_in"},
		{"in keep-out", Point{0.005, 0.0051}, h(100), "keep_out"},
		{"above ceiling", Point{0, 0}, h(250), "keep_in"},
		{"bad lat", Point{95, 0}, nil, "keep_in"},
	}
	for _, c := range cases {
		v := f.Check(c.p, c.hae)
		got := ""
		if v != nil {
			got = v.Detail
		}
		if got != c.detail {
			t.Errorf("%s: detail %q want %q (%v)", c.name, got, c.detail, v)
		}
	}
}

func TestFailClosed(t *testing.T) {
	var f *Fence
	if v := f.Check(Point{0, 0}, nil); v == nil || v.Detail != "keep_in" {
		t.Fatalf("nil fence must fail closed, got %v", v)
	}
	empty, err := Parse([]byte(`{}`))
	if err != nil {
		t.Fatal(err)
	}
	if v := empty.Check(Point{0, 0}, nil); v == nil {
		t.Fatal("empty fence must fail closed")
	}
}

func TestParseRejects(t *testing.T) {
	for _, bad := range []string{
		`{"keep_in":[{"id":"a","polygon":[{"lat":0,"lon":0},{"lat":1,"lon":1}]}]}`,
		`{"keep_in":[{"polygon":[{"lat":0,"lon":0},{"lat":1,"lon":1},{"lat":1,"lon":0}]}]}`,
		`{"keep_in":[{"id":"c","circle":{"lat":0,"lon":0,"radius_m":0}}]}`,
		`{"keep_in":[{"id":"d","polygon":[{"lat":99,"lon":0},{"lat":1,"lon":1},{"lat":1,"lon":0}]}]}`,
		`{"keep_in":[],"max_hae_m":10,"min_hae_m":20}`,
		`not json`,
	} {
		if _, err := Parse([]byte(bad)); err == nil {
			t.Errorf("expected error for %s", bad)
		}
	}
}
