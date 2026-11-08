package parsers

import (
	"bytes"
	"encoding/json"
	"flag"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/onemindos/omos-bridge/ingest/internal/ingest"
)

var update = flag.Bool("update", false, "rewrite golden files from the recorded fixtures")

const (
	fixtures = "../../testdata/fixtures"
	golden   = "../../testdata/golden"
	schemas  = "../../../../schemas"
)

type goldenRec struct {
	Subject string         `json:"subject"`
	MsgID   string         `json:"msgId"`
	TS      string         `json:"ts"`
	Data    map[string]any `json:"data,omitempty"`
	Error   string         `json:"error,omitempty"`
}

func toGolden(recs []ingest.Record) []goldenRec {
	out := make([]goldenRec, 0, len(recs))
	for _, r := range recs {
		g := goldenRec{Subject: r.Subject, MsgID: r.ID, Data: r.Data}
		if !r.TS.IsZero() {
			g.TS = r.TS.UTC().Format(time.RFC3339)
		}
		if r.Err != nil {
			g.Error, g.Data = r.Err.Error(), nil
		}
		out = append(out, g)
	}
	return out
}

func checkGolden(t *testing.T, name string, recs []ingest.Record) {
	t.Helper()
	got, err := json.MarshalIndent(toGolden(recs), "", " ")
	if err != nil {
		t.Fatal(err)
	}
	got = append(got, '\n')
	path := filepath.Join(golden, name)
	if *update {
		if err := os.WriteFile(path, got, 0o644); err != nil {
			t.Fatal(err)
		}
	}
	want, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("%v (run go test ./internal/parsers -update)", err)
	}
	if !bytes.Equal(got, want) {
		t.Fatalf("%s differs from golden; diff the output of -update to review", name)
	}
}

func read(t *testing.T, p string) []byte {
	t.Helper()
	b, err := os.ReadFile(filepath.Join(fixtures, p))
	if err != nil {
		t.Fatal(err)
	}
	return b
}

// Replayable: a recorded api.weather.gov response (2026-10-29, trimmed to 12 alerts)
// always yields the same records, and every record passes the published schema.
func TestNWSGolden(t *testing.T) {
	sch, err := ingest.NewSchemas(schemas).Get("intel/weather-alert.v1.schema.json")
	if err != nil {
		t.Fatal(err)
	}
	for _, f := range []string{"active-2026-10-29.json", "active-2026-10-29-next.json"} {
		recs, err := NWSAlerts(read(t, "nws/"+f), nil)
		if err != nil {
			t.Fatal(err)
		}
		for _, r := range recs {
			if r.Err != nil {
				t.Fatalf("%s: unexpected bad record %v", f, r.Err)
			}
			if r.Subject != NWSSubject || !strings.HasPrefix(r.ID, "nws:urn:oid:") || r.TS.IsZero() {
				t.Fatalf("%s: %+v", f, r)
			}
			if err := ingest.Validate(sch, r.Data); err != nil {
				t.Fatalf("%s %s: schema: %v", f, r.ID, err)
			}
		}
		checkGolden(t, "nws-"+strings.TrimSuffix(f, ".json")+".golden.json", recs)
	}
	recs, _ := NWSAlerts(read(t, "nws/active-2026-10-29.json"), nil)
	if len(recs) != 12 {
		t.Fatalf("want 12 alerts, got %d", len(recs))
	}
	withPoint := 0
	for _, r := range recs {
		if _, ok := r.Data["lat"]; ok {
			withPoint++
		}
	}
	if withPoint != 3 {
		t.Fatalf("the fixture has 3 polygon alerts, got %d with a point", withPoint)
	}
}

// Malformed upstream data becomes a per-item error (dead letter), never a dropped item,
// and a broken body is a parse error.
func TestNWSMalformed(t *testing.T) {
	sch, _ := ingest.NewSchemas(schemas).Get("intel/weather-alert.v1.schema.json")
	body := read(t, "nws/malformed.json")
	recs, err := NWSAlerts(body, nil)
	if err != nil {
		t.Fatal(err)
	}
	var bad, invalid, good int
	for _, r := range recs {
		switch {
		case r.Err != nil:
			bad++
		case ingest.Validate(sch, r.Data) != nil:
			invalid++
		default:
			good++
		}
	}
	if bad != 2 || invalid != 3 || good != 1 {
		t.Fatalf("bad=%d invalid=%d good=%d", bad, invalid, good)
	}
	checkGolden(t, "nws-malformed.golden.json", recs)
	for _, b := range []string{"<html>503</html>", `{"type":"Feature"}`, `{"type":"FeatureCollection","features":`} {
		if _, err := NWSAlerts([]byte(b), nil); err == nil {
			t.Errorf("%q parsed", b)
		}
	}
}

// Clock skew: dedup keys and envelope ids come from upstream ids only.
func TestKeysIndependentOfClock(t *testing.T) {
	recs, _ := NWSAlerts(read(t, "nws/active-2026-10-29.json"), nil)
	e1 := ingest.NewEnvelope("nws-alerts", recs[0].Subject, recs[0].ID, recs[0].TS, time.Now(), recs[0].Data)
	e2 := ingest.NewEnvelope("nws-alerts", recs[0].Subject, recs[0].ID, recs[0].TS, time.Now().Add(10*365*24*time.Hour), recs[0].Data)
	if e1.ID != e2.ID || e1.TS != e2.TS || e1.MsgID != e2.MsgID {
		t.Fatalf("skewed clock changed identity: %+v %+v", e1, e2)
	}
}

func TestRSSAndAtom(t *testing.T) {
	recs, err := RSS(read(t, "rss/gdacs.xml"), map[string]string{"prefix": "gdacs", "subject": "intel.disaster.gdacs.event"})
	if err != nil {
		t.Fatal(err)
	}
	checkGolden(t, "rss-gdacs.golden.json", recs)
	if len(recs) != 3 || recs[0].Data["lat"] != 14.5 || recs[2].Err == nil {
		t.Fatalf("%+v", toGolden(recs))
	}
	atom, err := RSS(read(t, "rss/atom.xml"), map[string]string{"prefix": "atom"})
	if err != nil || len(atom) != 2 || atom[0].ID != "atom:urn:uuid:1" || atom[1].Data["lat"] != 40.7 {
		t.Fatalf("%+v %v", toGolden(atom), err)
	}
	if _, err := RSS([]byte("<html><body/></html>"), nil); err == nil {
		t.Fatal("html accepted as rss")
	}
}

func TestCSVAndGeoJSONAndJSON(t *testing.T) {
	recs, err := CSV(read(t, "file/firms-sample.csv"), map[string]string{"prefix": "firms", "idFields": "latitude,longitude,acq_date,acq_time", "latField": "latitude", "lonField": "longitude"})
	if err != nil {
		t.Fatal(err)
	}
	checkGolden(t, "csv-firms.golden.json", recs)
	if len(recs) != 4 || recs[0].ID != "firms:34.1234:-117.5678:2026-10-29:1030" || recs[3].Err == nil {
		t.Fatalf("%+v", toGolden(recs))
	}
	gj, err := GeoJSON(read(t, "geojson/usgs-sample.geojson"), map[string]string{"prefix": "usgs", "timeField": "time"})
	if err != nil {
		t.Fatal(err)
	}
	checkGolden(t, "geojson-usgs.golden.json", gj)
	if len(gj) != 2 || gj[0].ID != "usgs:us7000abcd" || gj[0].TS.IsZero() {
		t.Fatalf("%+v", toGolden(gj))
	}
	// A large integer id must not become float exponent notation (finding 6).
	js, err := JSONMessage([]byte(`{"mmsi":366000184,"ts":1761732000,"lat":40.1,"lon":-73.9}`), map[string]string{"prefix": "ais", "idFields": "mmsi,ts", "timeField": "ts"})
	if err != nil || js[0].ID != "ais:366000184:1761732000" || js[0].TS.Unix() != 1761732000 {
		t.Fatalf("%+v %v", js, err)
	}
	miss, _ := JSONMessage([]byte(`{"x":1}`), map[string]string{"idField": "id"})
	if miss[0].Err == nil {
		t.Fatal("record without id must be a dead letter")
	}
}
