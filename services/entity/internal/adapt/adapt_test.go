package adapt

import (
	"bufio"
	"encoding/json"
	"math"
	"os"
	"testing"
)

type fixture struct {
	Subject string `json:"subject"`
	Data    string `json:"data"`
}

func fixtures(t *testing.T) []fixture {
	t.Helper()
	f, err := os.Open("../../testdata/fixtures/replay.jsonl")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	var out []fixture
	sc := bufio.NewScanner(f)
	sc.Buffer(make([]byte, 1<<20), 1<<20)
	for sc.Scan() {
		var fx fixture
		_ = json.Unmarshal(sc.Bytes(), &fx)
		out = append(out, fx)
	}
	return out
}

func TestEveryFixtureParsesOrSkipsOnPurpose(t *testing.T) {
	var parsed, skipped, failed int
	for _, fx := range fixtures(t) {
		o, skip, err := Parse(fx.Subject, []byte(fx.Data))
		switch {
		case err != nil:
			failed++
		case o == nil:
			skipped++
			if skip == "" {
				t.Fatalf("%s: nil observation without a reason", fx.Subject)
			}
		default:
			parsed++
			if len(o.Keys) == 0 || o.SourceT == 0 || o.EntityType == "" {
				t.Fatalf("%s: incomplete observation %+v", fx.Subject, o)
			}
		}
	}
	// 1 malformed CoT; 2 skips (projector echo, intel without coordinates).
	if failed != 1 || skipped != 2 || parsed != 27 {
		t.Fatalf("parsed=%d skipped=%d failed=%d", parsed, skipped, failed)
	}
}

func TestCoTAircraftFromFeed(t *testing.T) {
	x := `<?xml version="1.0"?><event version="2.0" type="a-f-A-M-F-Q" uid="ADSB-A1B2C3" how="m-g" time="2026-10-29T12:00:00.000Z" start="2026-10-29T12:00:00.000Z" stale="2026-10-29T12:01:00.000Z"><point lat="40.64" lon="-73.78" hae="1200" ce="9999999.0" le="9999999.0"/><detail><contact callsign="UAL123"/><track speed="250.0" course="45.0"/><remarks>ICAO:a1b2c3 United States</remarks></detail></event>`
	o, _, err := Parse("tak.cot.inject", []byte(x))
	if err != nil {
		t.Fatal(err)
	}
	if o.Keys[0].Type != "icao" || o.Keys[0].ID != "a1b2c3" || o.EntityType != "aircraft" || o.Integration != "feed-aircraft" {
		t.Fatalf("%+v", o)
	}
	if o.CE != 0 || math.Abs(*o.Speed-128.611) > 0.01 || o.StaleS != 60 || o.Marking != "CLEAR" || *o.EgressTAK {
		t.Fatalf("ce=%v speed=%v stale=%d marking=%s", o.CE, *o.Speed, o.StaleS, o.Marking)
	}
}

func TestCoTEchoFilter(t *testing.T) {
	for _, x := range []string{
		`<event version="2.0" type="a-f-A" uid="OM.aircraft.icao-a1b2c3" time="2026-10-29T12:00:00Z"><point lat="1" lon="1"/></event>`,
		`<event version="2.0" type="a-f-A" uid="X" time="2026-10-29T12:00:00Z"><point lat="1" lon="1"/><detail><_flow-tags_ omos-tak-projector="2026-10-29T12:00:00Z"/></detail></event>`,
	} {
		if o, skip, err := Parse("tak.cot.inject", []byte(x)); o != nil || skip != SkipEcho || err != nil {
			t.Fatalf("echo not filtered: %v %q %v", o, skip, err)
		}
	}
}

func TestCoTEUDIsPersonAndHostileFlagged(t *testing.T) {
	x := `<event version="2.0" uid="ANDROID-1" type="a-h-G-U-C" how="m-g" time="2026-10-29T12:00:00Z" start="2026-10-29T12:00:00Z" stale="2026-10-29T12:06:00Z"><point lat="30" lon="-97" hae="0" ce="9.9" le="9999999"/><detail><takv platform="ATAK-CIV"/><contact callsign="X"/><__group name="Cyan" role="Team Member"/></detail></event>`
	o, _, err := Parse("tak.cot.out", []byte(x))
	if err != nil {
		t.Fatal(err)
	}
	if !o.Person || o.EntityType != "person" || o.Disposition != "hostile" || o.Marking != "AMBER+STRICT" || o.Integration != "tak-server" {
		t.Fatalf("%+v", o)
	}
}

func TestCoTRejectsGarbage(t *testing.T) {
	for _, x := range []string{"not xml", "<event uid='broken'", `<event type="a-f-A" time="2026-10-29T12:00:00Z"><point lat="1" lon="1"/></event>`,
		`<event type="a-f-A" uid="X" time="2026-10-29T12:00:00Z"><point lat="100" lon="1"/></event>`} {
		if _, _, err := Parse("tak.cot.inject", []byte(x)); err == nil {
			t.Fatalf("accepted %q", x)
		}
	}
	if o, skip, _ := Parse("tak.cot.inject", []byte(`<event type="b-m-p-s-m" uid="M" time="2026-10-29T12:00:00Z"><point lat="1" lon="1"/></event>`)); o != nil || skip != SkipUnsupported {
		t.Fatal("marker type not skipped")
	}
}

func TestWorldMilitaryIsOutOfScope(t *testing.T) {
	for _, s := range []string{"evt.world.air.mil.ae1234", "evt.world.military.base1"} {
		if o, skip, _ := Parse(s, []byte(`{"uid":"x","lat":1,"lon":1,"t":1}`)); o != nil || skip != SkipOutOfScope {
			t.Fatalf("%s not skipped", s)
		}
	}
}

func TestFleetStatusHasNoPosition(t *testing.T) {
	o, _, err := Parse("tel.fleet.drone.sim-x1.status", []byte(`{"uid":"sim-x1","lat":0,"lon":0,"t":5,"meta":{"battery_pct":50}}`))
	if err != nil || o.HasPos || o.Health["battery_pct"] != 50.0 || !o.Simulated || o.Keys[0].ID != "sim-x1" {
		t.Fatalf("%+v %v", o, err)
	}
}

func TestIntelUSGS(t *testing.T) {
	o, _, err := Parse("intel.seismic.usgs.m4.us7000abcd", []byte(`{"id":"us7000abcd","type":"earthquake","source":"usgs","lat":35.1,"lon":-117.5,"timestamp":"2026-10-29T11:50:00.000Z","expires":"","properties":{"magnitude":4.4,"place":"Ridgecrest"}}`))
	if err != nil || o.Keys[0].ID != "usgs-us7000abcd" || o.EntityType != "event" || o.Template != "geo" || o.Callsign != "Ridgecrest" {
		t.Fatalf("%+v %v", o, err)
	}
}

func TestDetPersonIsTrack(t *testing.T) {
	o, _, err := Parse("det.land.camera.42", []byte(`{"lat":30,"lon":-97,"t":10,"class":"person","sensor_id":"cam0"}`))
	if err != nil || o.EntityType != "track" || o.Kind != "person" || o.Marking != "AMBER+STRICT" || *o.EgressTAK {
		t.Fatalf("%+v %v", o, err)
	}
}
