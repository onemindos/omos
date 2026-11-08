// Package parsers holds the upstream parsers that sources.yaml can name. Each parser is
// pure (bytes in, records out) so recorded fixtures drive golden tests.
package parsers

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/onemindos/omos-bridge/ingest/internal/ingest"
)

func init() {
	ingest.RegisterParser("nws-alerts", NWSAlerts)
	ingest.RegisterParser("rss", RSS)
	ingest.RegisterParser("geojson", GeoJSON)
	ingest.RegisterParser("csv", CSV)
	ingest.RegisterParser("json", JSONMessage)
}

// NWSSubject is the default subject for NWS alerts: intel.{domain}.{source}.{type} with
// source token "noaa", as the feed registry row F094 maps it (onemind-xr
// docs/sources/sources.yaml, f094-noaa-nws-alerts). The legacy feed-intel subjects
// intel.weather.noaa.alert.{state}.{event} keep running unchanged (BUS-P3).
const NWSSubject = "intel.weather.noaa.alert"

type nwsCollection struct {
	Type     string            `json:"type"`
	Features []json.RawMessage `json:"features"`
}

type nwsFeature struct {
	ID         string          `json:"id"`
	Geometry   json.RawMessage `json:"geometry"`
	Properties *struct {
		ID            string              `json:"id"`
		AreaDesc      string              `json:"areaDesc"`
		Geocode       map[string][]string `json:"geocode"`
		AffectedZones []string            `json:"affectedZones"`
		References    []struct {
			Identifier string `json:"identifier"`
		} `json:"references"`
		Sent        string `json:"sent"`
		Effective   string `json:"effective"`
		Onset       string `json:"onset"`
		Expires     string `json:"expires"`
		Ends        string `json:"ends"`
		Status      string `json:"status"`
		MessageType string `json:"messageType"`
		Category    string `json:"category"`
		Severity    string `json:"severity"`
		Certainty   string `json:"certainty"`
		Urgency     string `json:"urgency"`
		Event       string `json:"event"`
		SenderName  string `json:"senderName"`
		Headline    string `json:"headline"`
		Description string `json:"description"`
		Instruction string `json:"instruction"`
		Response    string `json:"response"`
		Web         string `json:"web"`
	} `json:"properties"`
}

// NWSAlerts parses api.weather.gov /alerts/active (GeoJSON FeatureCollection). Each alert
// becomes one record keyed "nws:{alert id}"; an update is a new alert id upstream, so
// updates are new records and re-listing an alert is not. Alerts without geometry (most
// are zone-based) carry affectedZones and UGC/SAME codes; resolving zone shapes is left
// to the entity service so the collector makes one upstream request per poll.
func NWSAlerts(body []byte, opts map[string]string) ([]ingest.Record, error) {
	var fc nwsCollection
	if err := json.Unmarshal(body, &fc); err != nil {
		return nil, err
	}
	if fc.Type != "FeatureCollection" {
		return nil, fmt.Errorf("not a FeatureCollection (type %q)", fc.Type)
	}
	subject := opts["subject"]
	if subject == "" {
		subject = NWSSubject
	}
	out := make([]ingest.Record, 0, len(fc.Features))
	for _, raw := range fc.Features {
		var f nwsFeature
		if err := json.Unmarshal(raw, &f); err != nil {
			out = append(out, ingest.Record{Raw: raw, Err: fmt.Errorf("feature: %v", err)})
			continue
		}
		if f.Properties == nil {
			out = append(out, ingest.Record{Raw: raw, Err: errors.New("feature has no properties")})
			continue
		}
		p := f.Properties
		id := p.ID
		if id == "" {
			out = append(out, ingest.Record{Raw: raw, Err: errors.New("alert has no properties.id")})
			continue
		}
		refs := make([]string, 0, len(p.References))
		for _, r := range p.References {
			refs = append(refs, r.Identifier)
		}
		d := map[string]any{
			"id":            id,
			"type":          "weather_alert",
			"event":         p.Event,
			"headline":      p.Headline,
			"description":   trunc(p.Description, 1000),
			"instruction":   trunc(p.Instruction, 1000),
			"severity":      p.Severity,
			"urgency":       p.Urgency,
			"certainty":     p.Certainty,
			"status":        p.Status,
			"messageType":   p.MessageType,
			"category":      p.Category,
			"response":      p.Response,
			"sent":          p.Sent,
			"effective":     p.Effective,
			"onset":         nullable(p.Onset),
			"expires":       p.Expires,
			"ends":          nullable(p.Ends),
			"areaDesc":      p.AreaDesc,
			"senderName":    p.SenderName,
			"affectedZones": nonNil(p.AffectedZones),
			"ugc":           nonNil(p.Geocode["UGC"]),
			"same":          nonNil(p.Geocode["SAME"]),
			"references":    refs,
			"web":           p.Web,
		}
		if g := parseGeometry(f.Geometry); g != nil {
			d["geometry"] = g.raw
			d["lat"], d["lon"] = round6(g.lat), round6(g.lon)
		}
		ts, _ := time.Parse(time.RFC3339, p.Sent)
		out = append(out, ingest.Record{Subject: subject, ID: "nws:" + id, TS: ts, Data: d, Raw: raw})
	}
	return out, nil
}

type geom struct {
	raw      map[string]any
	lat, lon float64
}

// parseGeometry returns the geometry and a representative point: the vertex average of a
// Polygon's outer ring (or of all outer rings of a MultiPolygon), or the Point itself.
func parseGeometry(b json.RawMessage) *geom {
	if len(b) == 0 || string(b) == "null" {
		return nil
	}
	var g struct {
		Type        string          `json:"type"`
		Coordinates json.RawMessage `json:"coordinates"`
	}
	if json.Unmarshal(b, &g) != nil {
		return nil
	}
	var raw map[string]any
	_ = json.Unmarshal(b, &raw)
	var pts [][]float64
	switch g.Type {
	case "Point":
		var p []float64
		if json.Unmarshal(g.Coordinates, &p) != nil || len(p) < 2 {
			return nil
		}
		pts = [][]float64{p}
	case "Polygon":
		var rings [][][]float64
		if json.Unmarshal(g.Coordinates, &rings) != nil || len(rings) == 0 {
			return nil
		}
		pts = openRing(rings[0])
	case "MultiPolygon":
		var polys [][][][]float64
		if json.Unmarshal(g.Coordinates, &polys) != nil {
			return nil
		}
		for _, poly := range polys {
			if len(poly) > 0 {
				pts = append(pts, openRing(poly[0])...)
			}
		}
	default:
		return nil
	}
	var sx, sy float64
	n := 0
	for _, p := range pts {
		if len(p) < 2 || math.IsNaN(p[0]) || math.IsNaN(p[1]) {
			continue
		}
		sx, sy, n = sx+p[0], sy+p[1], n+1
	}
	if n == 0 {
		return nil
	}
	return &geom{raw: raw, lon: sx / float64(n), lat: sy / float64(n)}
}

// openRing drops the closing vertex that repeats the first.
func openRing(r [][]float64) [][]float64 {
	if len(r) > 1 {
		a, b := r[0], r[len(r)-1]
		if len(a) >= 2 && len(b) >= 2 && a[0] == b[0] && a[1] == b[1] {
			return r[:len(r)-1]
		}
	}
	return r
}

func round6(v float64) float64 { return math.Round(v*1e6) / 1e6 }

func trunc(s string, n int) string {
	if len(s) <= n {
		return s
	}
	// cut on a rune boundary
	for n > 0 && (s[n]&0xC0) == 0x80 {
		n--
	}
	return s[:n]
}

func nullable(s string) any {
	if s == "" {
		return nil
	}
	return s
}

func nonNil(v []string) []string {
	if v == nil {
		return []string{}
	}
	return v
}
