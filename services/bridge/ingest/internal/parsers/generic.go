package parsers

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/onemindos/omos-bridge/ingest/internal/ingest"
)

// Generic parsers take their field names from the source's options map:
//   prefix     dedup key prefix (default: the parser name)
//   idField    property holding the upstream id (geojson: feature id if empty; csv/json required)
//   idFields   comma-separated fields joined with ":" (csv/json; for sources without one id)
//   latField, lonField, timeField, timeLayout (Go layout, default RFC 3339)
//   subject    overrides the source subject
//   subjectTemplate  per-record subject, e.g. "evt.world.sea.{mmsi}" (json parser)

func opt(o map[string]string, k, def string) string {
	if v, ok := o[k]; ok && v != "" {
		return v
	}
	return def
}

func parseTime(v any, layout string) time.Time {
	switch t := v.(type) {
	case string:
		if ts, err := time.Parse(layout, t); err == nil {
			return ts
		}
	case float64:
		return parseTime(int64(t), layout)
	case int64:
		if t > 1e12 {
			return time.UnixMilli(t).UTC()
		}
		return time.Unix(t, 0).UTC()
	}
	return time.Time{}
}

// keyOf builds the dedup key from one or more fields. Numbers are formatted as integers
// when they are whole, never in float exponent form (ingestion-platform.md finding 6).
func keyOf(m map[string]any, o map[string]string) (string, error) {
	fields := strings.Split(opt(o, "idFields", opt(o, "idField", "id")), ",")
	parts := make([]string, 0, len(fields))
	for _, f := range fields {
		v, ok := m[strings.TrimSpace(f)]
		if !ok || v == nil || v == "" {
			return "", fmt.Errorf("missing id field %q", f)
		}
		parts = append(parts, scalar(v))
	}
	return strings.Join(parts, ":"), nil
}

func scalar(v any) string {
	switch t := v.(type) {
	case float64:
		if t == float64(int64(t)) {
			return strconv.FormatInt(int64(t), 10)
		}
		return strconv.FormatFloat(t, 'f', -1, 64)
	case json.Number:
		return t.String()
	default:
		return fmt.Sprint(t)
	}
}

// JSONMessage parses one JSON object (an SSE event or a WebSocket frame).
func JSONMessage(body []byte, o map[string]string) ([]ingest.Record, error) {
	var m map[string]any
	dec := json.NewDecoder(bytes.NewReader(body))
	dec.UseNumber()
	if err := dec.Decode(&m); err != nil {
		return nil, err
	}
	normalizeNumbers(m)
	k, err := keyOf(m, o)
	if err != nil {
		return []ingest.Record{{Raw: body, Err: err}}, nil
	}
	if _, ok := m["id"]; !ok {
		m["id"] = k
	}
	return []ingest.Record{{Subject: subjectFor(m, o), ID: opt(o, "prefix", "json") + ":" + k,
		TS: parseTime(m[opt(o, "timeField", "time")], opt(o, "timeLayout", time.RFC3339)), Data: m, Raw: body}}, nil
}

// normalizeNumbers turns json.Number into int64 when whole, else float64, so ids keep
// their digits and schemas see numbers.
func normalizeNumbers(m map[string]any) {
	for k, v := range m {
		switch t := v.(type) {
		case json.Number:
			if i, err := t.Int64(); err == nil {
				m[k] = i
			} else if f, err := t.Float64(); err == nil {
				m[k] = f
			}
		case map[string]any:
			normalizeNumbers(t)
		}
	}
}

// GeoJSON parses a FeatureCollection; each feature's properties become the record data,
// with lat/lon from a Point geometry.
func GeoJSON(body []byte, o map[string]string) ([]ingest.Record, error) {
	var fc struct {
		Type     string `json:"type"`
		Features []json.RawMessage
	}
	if err := json.Unmarshal(body, &fc); err != nil {
		return nil, err
	}
	if fc.Type != "FeatureCollection" {
		return nil, fmt.Errorf("not a FeatureCollection (type %q)", fc.Type)
	}
	var out []ingest.Record
	for _, raw := range fc.Features {
		var f struct {
			ID         any            `json:"id"`
			Properties map[string]any `json:"properties"`
			Geometry   json.RawMessage
		}
		if err := json.Unmarshal(raw, &f); err != nil || f.Properties == nil {
			out = append(out, ingest.Record{Raw: raw, Err: errors.New("feature without properties")})
			continue
		}
		d := f.Properties
		var id string
		if fld := o["idField"]; fld != "" {
			if v, ok := d[fld]; ok && v != nil {
				id = scalar(v)
			}
		} else if f.ID != nil {
			id = scalar(f.ID)
		}
		if id == "" {
			out = append(out, ingest.Record{Raw: raw, Err: errors.New("feature has no id")})
			continue
		}
		d["id"] = id
		if g := parseGeometry(f.Geometry); g != nil {
			d["lat"], d["lon"] = round6(g.lat), round6(g.lon)
		}
		out = append(out, ingest.Record{Subject: o["subject"], ID: opt(o, "prefix", "geojson") + ":" + id,
			TS: parseTime(d[opt(o, "timeField", "time")], opt(o, "timeLayout", time.RFC3339)), Data: d, Raw: raw})
	}
	return out, nil
}

// CSV parses a header row plus records (bulk files such as NASA FIRMS). Numeric-looking
// cells become numbers.
func CSV(body []byte, o map[string]string) ([]ingest.Record, error) {
	r := csv.NewReader(bytes.NewReader(body))
	r.FieldsPerRecord = -1
	r.ReuseRecord = false
	head, err := r.Read()
	if err != nil {
		return nil, err
	}
	var out []ingest.Record
	for line := 2; ; line++ {
		row, err := r.Read()
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			out = append(out, ingest.Record{Raw: []byte(fmt.Sprintf("line %d", line)), Err: err})
			break
		}
		raw := []byte(strings.Join(row, ","))
		if len(row) != len(head) {
			out = append(out, ingest.Record{Raw: raw, Err: fmt.Errorf("line %d: %d fields, header has %d", line, len(row), len(head))})
			continue
		}
		d := make(map[string]any, len(head))
		for i, h := range head {
			d[h] = cell(row[i])
		}
		k, err := keyOf(d, o)
		if err != nil {
			out = append(out, ingest.Record{Raw: raw, Err: err})
			continue
		}
		if _, ok := d["id"]; !ok {
			d["id"] = k
		}
		if lf, lo := o["latField"], o["lonField"]; lf != "" && lo != "" {
			d["lat"], d["lon"] = d[lf], d[lo]
		}
		out = append(out, ingest.Record{Subject: o["subject"], ID: opt(o, "prefix", "csv") + ":" + k,
			TS: parseTime(d[opt(o, "timeField", "time")], opt(o, "timeLayout", time.RFC3339)), Data: d, Raw: raw})
	}
	return out, nil
}

func cell(s string) any {
	if i, err := strconv.ParseInt(s, 10, 64); err == nil {
		return i
	}
	if f, err := strconv.ParseFloat(s, 64); err == nil {
		return f
	}
	return s
}

// RSS parses RSS 2.0 items and Atom entries, with W3C geo or GeoRSS points.
func RSS(body []byte, o map[string]string) ([]ingest.Record, error) {
	var doc struct {
		XMLName xml.Name
		Items   []rssItem `xml:"channel>item"`
		Entries []rssItem `xml:"entry"`
	}
	if err := xml.Unmarshal(body, &doc); err != nil {
		return nil, err
	}
	items := append(doc.Items, doc.Entries...)
	if doc.XMLName.Local != "rss" && doc.XMLName.Local != "feed" && doc.XMLName.Local != "RDF" {
		return nil, fmt.Errorf("not an RSS or Atom document (root %q)", doc.XMLName.Local)
	}
	var out []ingest.Record
	for _, it := range items {
		raw, _ := xml.Marshal(it)
		id := strings.TrimSpace(it.GUID)
		if id == "" {
			id = strings.TrimSpace(it.ID)
		}
		if id == "" {
			id = strings.TrimSpace(it.Link.Href + it.Link.Text)
		}
		if id == "" {
			out = append(out, ingest.Record{Raw: raw, Err: errors.New("item has no guid, id or link")})
			continue
		}
		d := map[string]any{"id": id, "title": strings.TrimSpace(it.Title), "link": strings.TrimSpace(it.Link.Href + it.Link.Text),
			"summary": trunc(strings.TrimSpace(it.Description+it.Summary), 1000)}
		published := strings.TrimSpace(it.PubDate + it.Updated)
		var ts time.Time
		for _, layout := range []string{time.RFC1123Z, time.RFC1123, time.RFC3339} {
			if t, err := time.Parse(layout, published); err == nil {
				ts = t
				break
			}
		}
		if published != "" {
			d["published"] = published
		}
		lat, lon, ok := it.point()
		if ok {
			d["lat"], d["lon"] = lat, lon
		}
		out = append(out, ingest.Record{Subject: o["subject"], ID: opt(o, "prefix", "rss") + ":" + id, TS: ts, Data: d, Raw: raw})
	}
	return out, nil
}

type rssItem struct {
	Title string `xml:"title"`
	GUID  string `xml:"guid"`
	ID    string `xml:"id"`
	Link  struct {
		Href string `xml:"href,attr"`
		Text string `xml:",chardata"`
	} `xml:"link"`
	Description string `xml:"description"`
	Summary     string `xml:"summary"`
	PubDate     string `xml:"pubDate"`
	Updated     string `xml:"updated"`
	GeoLat      string `xml:"http://www.w3.org/2003/01/geo/wgs84_pos# lat"`
	GeoLong     string `xml:"http://www.w3.org/2003/01/geo/wgs84_pos# long"`
	GeoRSSPoint string `xml:"http://www.georss.org/georss point"`
}

func (it rssItem) point() (float64, float64, bool) {
	if it.GeoLat != "" && it.GeoLong != "" {
		la, e1 := strconv.ParseFloat(strings.TrimSpace(it.GeoLat), 64)
		lo, e2 := strconv.ParseFloat(strings.TrimSpace(it.GeoLong), 64)
		return la, lo, e1 == nil && e2 == nil
	}
	if f := strings.Fields(it.GeoRSSPoint); len(f) == 2 {
		la, e1 := strconv.ParseFloat(f[0], 64)
		lo, e2 := strconv.ParseFloat(f[1], 64)
		return la, lo, e1 == nil && e2 == nil
	}
	return 0, 0, false
}

// subjectFor fills {field} placeholders in options.subjectTemplate. A value that would
// add or break tokens (dots, spaces, wildcards) yields an invalid subject, which the
// runtime dead-letters instead of publishing.
func subjectFor(m map[string]any, o map[string]string) string {
	tpl := o["subjectTemplate"]
	if tpl == "" {
		return o["subject"]
	}
	out := tpl
	for {
		i := strings.Index(out, "{")
		j := strings.Index(out, "}")
		if i < 0 || j < i {
			return out
		}
		v, ok := m[out[i+1:j]]
		val := ""
		if ok && v != nil {
			val = scalar(v)
		}
		if val == "" {
			val = "_missing_."
		}
		out = out[:i] + val + out[j+1:]
	}
}
