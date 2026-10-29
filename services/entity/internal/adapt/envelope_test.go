package adapt

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/onemindos/omos/services/entity/internal/identity"
)

var update = flag.Bool("update", false, "rewrite testdata/golden/intel-envelope.json")

// ingestRow is one omos-ingest golden parser output (testdata/fixtures/ingest/README.md).
type ingestRow struct {
	Subject string         `json:"subject"`
	MsgID   string         `json:"msgId"`
	TS      string         `json:"ts"`
	Data    map[string]any `json:"data"`
	Error   string         `json:"error"`
}

// stableUUID and envelope mirror omos-ingest record.go (StableUUID, NewEnvelope).
func stableUUID(key string) string {
	h := sha256.Sum256([]byte(key))
	h[6] = (h[6] & 0x0f) | 0x40
	h[8] = (h[8] & 0x3f) | 0x80
	x := hex.EncodeToString(h[:16])
	return x[0:8] + "-" + x[8:12] + "-" + x[12:16] + "-" + x[16:20] + "-" + x[20:32]
}

func envelope(t *testing.T, sourceID, subject string, r ingestRow) []byte {
	t.Helper()
	ts := r.TS
	if pt, err := time.Parse(time.RFC3339, ts); err == nil {
		ts = pt.UTC().Format("2006-01-02T15:04:05.000Z07:00")
	} else {
		ts = "2026-10-29T12:00:00.000Z" // NewEnvelope falls back to fetchedAt
	}
	b, err := json.Marshal(map[string]any{
		"v": 1, "id": stableUUID(r.MsgID), "ts": ts, "type": subject, "source": "ingest." + sourceID,
		"subject": subject, "data": r.Data, "fetchedAt": "2026-10-29T12:00:00.000Z", "msgId": r.MsgID,
	})
	if err != nil {
		t.Fatal(err)
	}
	return b
}

type goldenObs struct {
	File        string   `json:"file"`
	MsgID       string   `json:"msgId"`
	Result      string   `json:"result"` // "entity", "skip:<reason>", "error:<message>"
	EntityID    string   `json:"entity_id,omitempty"`
	EntityType  string   `json:"entity_type,omitempty"`
	Kind        string   `json:"kind,omitempty"`
	Callsign    string   `json:"callsign,omitempty"`
	Integration string   `json:"integration,omitempty"`
	Lat         float64  `json:"lat,omitempty"`
	Lon         float64  `json:"lon,omitempty"`
	SourceT     int64    `json:"source_t,omitempty"`
	ExpireT     int64    `json:"expire_t,omitempty"`
	MetaKeys    []string `json:"meta_keys,omitempty"`
}

// TestEnvelopeGolden runs omos-ingest's real NWS and USGS outputs (#553), wrapped in the bus
// envelope, through the intel adapter.
func TestEnvelopeGolden(t *testing.T) {
	inputs := []struct{ file, source, subject string }{
		{"nws-active-2026-10-29.golden.json", "nws-alerts", "intel.weather.noaa.alert"},
		{"nws-active-2026-10-29-next.golden.json", "nws-alerts", "intel.weather.noaa.alert"},
		{"nws-malformed.golden.json", "nws-alerts", "intel.weather.noaa.alert"},
		// usgs rows carry subject "" (the source's configured subject); ingestion-platform.md
		// pilot table: intel.seismic.usgs.event.
		{"geojson-usgs.golden.json", "usgs", "intel.seismic.usgs.event"},
	}
	var got []goldenObs
	counts := map[string]int{}
	for _, in := range inputs {
		b, err := os.ReadFile("../../testdata/fixtures/ingest/" + in.file)
		if err != nil {
			t.Fatal(err)
		}
		var rows []ingestRow
		if err := json.Unmarshal(b, &rows); err != nil {
			t.Fatal(err)
		}
		for _, r := range rows {
			if r.Error != "" || r.Data == nil {
				continue // dead-lettered by omos-ingest, never published
			}
			subject := r.Subject
			if subject == "" {
				subject = in.subject
			}
			o, skip, err := Parse(subject, envelope(t, in.source, subject, r))
			g := goldenObs{File: in.file, MsgID: r.MsgID}
			switch {
			case err != nil:
				g.Result = "error:" + err.Error()
			case o == nil:
				g.Result = "skip:" + skip
			default:
				g.Result = "entity"
				g.EntityID = identity.Canonical(o.Keys[0])
				g.EntityType, g.Kind, g.Callsign, g.Integration = o.EntityType, o.Kind, o.Callsign, o.Integration
				g.Lat, g.Lon, g.SourceT, g.ExpireT = o.Lat, o.Lon, o.SourceT, o.ExpireT
				for k := range o.Meta {
					g.MetaKeys = append(g.MetaKeys, k)
				}
				sort.Strings(g.MetaKeys)
			}
			counts[strings.SplitN(g.Result, ":", 2)[0]]++
			got = append(got, g)
		}
	}
	out, _ := json.MarshalIndent(got, "", " ")
	out = append(out, '\n')
	const path = "../../testdata/golden/intel-envelope.json"
	if *update {
		if err := os.WriteFile(path, out, 0o644); err != nil {
			t.Fatal(err)
		}
	}
	want, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("%v (run with -update)", err)
	}
	if !bytes.Equal(want, out) {
		t.Fatalf("intel envelope output differs from %s (run with -update and review the diff)", path)
	}
	// NWS: 3 alerts with a centroid in the first poll and 1 in the second; zone-only alerts
	// have no geometry and are skipped. The malformed file adds one out-of-range latitude
	// (error) and three skips. USGS: 2 quakes.
	if counts["entity"] != 6 || counts["error"] != 1 || counts["skip"] != 22 {
		t.Fatalf("counts %v", counts)
	}
	t.Logf("envelope golden: %v", counts)
}

func TestEnvelopeFieldMapping(t *testing.T) {
	r := ingestRow{MsgID: "usgs:us7000abcd", TS: "2025-10-29T09:00:00Z", Data: map[string]any{
		"id": "us7000abcd", "lat": 35.7, "lon": -117.5, "mag": 4.6, "place": "Synthetic place A", "time": 1761728400000,
	}}
	o, _, err := Parse("intel.seismic.usgs.event", envelope(t, "usgs", "intel.seismic.usgs.event", r))
	if err != nil {
		t.Fatal(err)
	}
	// Same entity id as the flat feed-intel record for the same quake (TestIntelUSGS).
	if o.Keys[0].ID != "usgs-us7000abcd" || o.Integration != "feed-intel-usgs" || o.Kind != "event" ||
		o.Callsign != "Synthetic place A" || o.SourceT != 1761728400000 || o.Meta["mag"] != 4.6 {
		t.Fatalf("%+v", o)
	}
	if _, ok := o.Meta["lat"]; ok {
		t.Fatal("lat copied into meta")
	}
}

func TestEnvelopeRejects(t *testing.T) {
	base := ingestRow{MsgID: "x", TS: "2026-10-29T09:00:00Z", Data: map[string]any{"id": "a", "lat": 1.0, "lon": 1.0}}
	good := envelope(t, "usgs", "intel.seismic.usgs.event", base)
	cases := map[string][]byte{
		"type differs from subject": bytes.Replace(good, []byte(`"type":"intel.seismic.usgs.event"`), []byte(`"type":"intel.weather.noaa.alert"`), 1),
		"version 2":                 bytes.Replace(good, []byte(`"v":1`), []byte(`"v":2`), 1),
		"no ts":                     bytes.Replace(good, []byte(`"ts":"2026-10-29T09:00:00.000Z"`), []byte(`"ts":""`), 1),
		"latitude out of range":     bytes.Replace(good, []byte(`"lat":1`), []byte(`"lat":201`), 1),
	}
	for name, b := range cases {
		if bytes.Equal(b, good) {
			t.Fatalf("%s: case did not change the input", name)
		}
		if o, _, err := Parse("intel.seismic.usgs.event", b); err == nil {
			t.Errorf("%s: accepted %+v", name, o)
		}
	}
}

func TestWorldRail(t *testing.T) {
	// The #558 feed-mta-grammar envelope for one Metro-North train.
	x := `{"uid":"mta-mnr-1234567","type":"ground.rail","lat":41.0,"lon":-73.8,"hae":0,"course":-1,"speed":12.5,"callsign":"MNR-1","stale":120,"t":1793275200000,"meta":{"agency":"MTA","feed":"MNR","trip_id":"1234567","route_id":"1"}}`
	o, skip, err := Parse("evt.world.rail.mta-mnr-1234567", []byte(x))
	if err != nil || o == nil {
		t.Fatalf("%v %q %v", o, skip, err)
	}
	if o.EntityType != "vehicle" || o.Kind != "train" || o.Keys[0].Type != identity.KeyGTFS || o.Keys[0].ID != "mnr-1234567" ||
		o.Integration != "world-rail" || o.DocType != "ground.rail" || o.Course != nil || *o.Speed != 12.5 || o.StaleS != 120 {
		t.Fatalf("%+v", o)
	}
	// The legacy MTA CoT for the same train resolves to the same entity id.
	cot := `<event version="2.0" type="a-n-G-E-V-T" uid="MTA-MNR-1234567" how="m-g" time="2026-10-29T12:00:00Z" start="2026-10-29T12:00:00Z" stale="2026-10-29T12:02:00Z"><point lat="41" lon="-73.8" hae="0" ce="50" le="50"/><detail><contact callsign="MNR-1"/></detail></event>`
	c, _, err := Parse("tak.cot.inject", []byte(cot))
	if err != nil {
		t.Fatal(err)
	}
	if identity.Canonical(c.Keys[0]) != identity.Canonical(o.Keys[0]) {
		t.Fatalf("CoT %s vs world %s", identity.Canonical(c.Keys[0]), identity.Canonical(o.Keys[0]))
	}
	// A rail uid without the mta- shape still keys on the uid.
	o2, _, err := Parse("evt.world.rail.njt-1", []byte(`{"uid":"njt-1","lat":40.7,"lon":-74.1,"t":1793275200000}`))
	if err != nil || o2.Keys[0].ID != "njt-1" || o2.DocType != "ground.rail" {
		t.Fatalf("%+v %v", o2, err)
	}
}

func TestIntelEchoPrefix(t *testing.T) {
	x := `<event version="2.0" type="a-u-G" uid="INTEL.usgs-us7000abcd" time="2026-10-29T12:00:00Z"><point lat="35" lon="-117"/></event>`
	if o, skip, err := Parse("tak.cot.out", []byte(x)); o != nil || skip != SkipEcho || err != nil {
		t.Fatalf("INTEL. uid not dropped: %v %q %v", o, skip, err)
	}
	// Only the prefix: a uid that merely contains INTEL is not an echo.
	y := `<event version="2.0" type="a-f-G" uid="X-INTEL.1" time="2026-10-29T12:00:00Z"><point lat="35" lon="-117"/></event>`
	if _, skip, _ := Parse("tak.cot.out", []byte(y)); skip == SkipEcho {
		t.Fatal("X-INTEL.1 dropped as echo")
	}
}

// issTLE is the ISS element set from #558's stations.tle fixture.
const issL1 = "1 25544U 98067A   26301.59601386  .00005724  00000+0  11161-3 0  9990"
const issL2 = "2 25544  51.6336 359.9178 0006994 320.4528  39.5949 15.49090321592436"

func tleRecord(l1, l2 string) []byte {
	// The body #558's feed-satellites-grammar publishes on intel.space.celestrak.tle.
	return []byte(fmt.Sprintf(`{"id":"celestrak-25544","type":"tle","source":"celestrak","lat":null,"lon":null,"timestamp":"2026-10-28T14:18:15.597Z","properties":{"norad_id":"25544","name":"ISS (ZARYA)","group":"stations","classification":"U","intl_designator":"98067A","epoch":"26301.59601386","epoch_ms":1793197095597,"line1":%q,"line2":%q}}`, l1, l2))
}

func TestTLEPropagates(t *testing.T) {
	defer func(f func() time.Time) { Now = f }(Now)
	Now = func() time.Time { return time.Date(2026, 10, 29, 12, 0, 0, 0, time.UTC) }
	o, skip, err := Parse(TLESubject, tleRecord(issL1, issL2))
	if err != nil || o == nil {
		t.Fatalf("%v %q %v", o, skip, err)
	}
	// python-sgp4 reference for this instant (internal/orbit/testdata/reference.json):
	// 2026-10-29T11:59:59.598Z lat 3.6335 lon 140.476406 alt 416.103 km.
	if o.SourceT != 1793275199598 || o.Lat < 3.62 || o.Lat > 3.65 || o.Lon < 140.46 || o.Lon > 140.49 ||
		*o.HAE < 415e3 || *o.HAE > 417e3 {
		t.Fatalf("position %v %v %v at %d", o.Lat, o.Lon, *o.HAE, o.SourceT)
	}
	if o.EntityType != "vehicle" || o.Kind != "satellite" || o.Environment != "space" || o.Keys[0].Type != identity.KeyNORAD ||
		identity.Canonical(o.Keys[0]) != "norad-25544" || o.Callsign != "ISS (ZARYA)" || o.Orbit == nil ||
		o.Meta["line1"] != issL1 || o.How != "m-p" {
		t.Fatalf("%+v", o)
	}
	// Re-propagation 30 s later moves it about 220 km along track.
	o2, err := Propagate(o, Now().Add(30*time.Second))
	if err != nil || o2.SourceT != o.SourceT+30000 || o2.Lat == o.Lat {
		t.Fatalf("%+v %v", o2, err)
	}
	// After a restart the template comes back from the restored entity's meta.
	o3, err := TLEFromMeta(o.Meta, Now())
	if err != nil || o3.Lat != o.Lat || o3.Keys[0] != o.Keys[0] {
		t.Fatalf("%+v %v", o3, err)
	}
}

func TestTLERejects(t *testing.T) {
	defer func(f func() time.Time) { Now = f }(Now)
	Now = func() time.Time { return time.Date(2026, 10, 29, 12, 0, 0, 0, time.UTC) }
	bad := issL1[:68] + "1" // checksum
	if _, _, err := Parse(TLESubject, tleRecord(bad, issL2)); err == nil {
		t.Error("bad checksum accepted")
	}
	if _, _, err := Parse(TLESubject, bytes.Replace(tleRecord(issL1, issL2), []byte(`"norad_id":"25544"`), []byte(`"norad_id":"36086"`), 1)); err == nil {
		t.Error("norad_id mismatch accepted")
	}
	Now = func() time.Time { return time.Date(2026, 11, 20, 0, 0, 0, 0, time.UTC) }
	if _, _, err := Parse(TLESubject, tleRecord(issL1, issL2)); err == nil {
		t.Error("23-day-old TLE propagated")
	}
}
