// SPDX-License-Identifier: Apache-2.0

package main

// Conversion between CoT XML (TAK Protocol version 0) and Event, following the
// xmlDetail rules in takproto/detail.proto: a detail child becomes a typed
// field only when the whole element maps onto it; everything else stays XML.

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"
)

const cotTimeLayout = "2006-01-02T15:04:05.000Z"

type xmlEvent struct {
	XMLName      xml.Name `xml:"event"`
	UID          string   `xml:"uid,attr"`
	Type         string   `xml:"type,attr"`
	How          string   `xml:"how,attr"`
	Time         string   `xml:"time,attr"`
	Start        string   `xml:"start,attr"`
	Stale        string   `xml:"stale,attr"`
	Access       string   `xml:"access,attr"`
	Caveat       string   `xml:"caveat,attr"`
	ReleasableTo string   `xml:"releasableTo,attr"`
	QoS          string   `xml:"qos,attr"`
	Opex         string   `xml:"opex,attr"`
	Point        struct {
		Lat string `xml:"lat,attr"`
		Lon string `xml:"lon,attr"`
		HAE string `xml:"hae,attr"`
		CE  string `xml:"ce,attr"`
		LE  string `xml:"le,attr"`
	} `xml:"point"`
	Detail *struct {
		Inner []byte `xml:",innerxml"`
	} `xml:"detail"`
}

func parseCotTime(s string) (uint64, error) {
	if s == "" {
		return 0, nil
	}
	t, err := time.Parse(time.RFC3339Nano, strings.TrimSpace(s))
	if err != nil {
		return 0, err
	}
	return uint64(t.UnixMilli()), nil
}

func formatCotTime(ms uint64) string {
	return time.UnixMilli(int64(ms)).UTC().Format(cotTimeLayout)
}

func parseNum(name, s string) (float64, error) {
	if strings.TrimSpace(s) == "" {
		return 0, nil
	}
	v, err := strconv.ParseFloat(strings.TrimSpace(s), 64)
	if err != nil {
		return 0, fmt.Errorf("point %s: %w", name, err)
	}
	return v, nil
}

// ParseCotXML decodes one CoT <event> into an Event.
func ParseCotXML(raw []byte) (*Event, error) {
	var x xmlEvent
	if err := xml.Unmarshal(raw, &x); err != nil {
		return nil, fmt.Errorf("xml decode: %w", err)
	}
	if x.UID == "" || x.Type == "" {
		return nil, fmt.Errorf("event missing uid or type")
	}
	e := &Event{
		Type: x.Type, Access: x.Access, Caveat: x.Caveat, ReleasableTo: x.ReleasableTo,
		QoS: x.QoS, Opex: x.Opex, UID: x.UID, How: x.How,
	}
	var err error
	if e.SendTime, err = parseCotTime(x.Time); err != nil {
		return nil, fmt.Errorf("time: %w", err)
	}
	if e.StartTime, err = parseCotTime(x.Start); err != nil {
		return nil, fmt.Errorf("start: %w", err)
	}
	if e.StaleTime, err = parseCotTime(x.Stale); err != nil {
		return nil, fmt.Errorf("stale: %w", err)
	}
	for _, p := range []struct {
		name string
		src  string
		dst  *float64
	}{
		{"lat", x.Point.Lat, &e.Lat}, {"lon", x.Point.Lon, &e.Lon}, {"hae", x.Point.HAE, &e.HAE},
		{"ce", x.Point.CE, &e.CE}, {"le", x.Point.LE, &e.LE},
	} {
		if *p.dst, err = parseNum(p.name, p.src); err != nil {
			return nil, err
		}
	}
	if x.Detail != nil {
		d, err := splitDetail(x.Detail.Inner)
		if err != nil {
			return nil, fmt.Errorf("detail: %w", err)
		}
		e.Detail = d
	}
	return e, nil
}

// detailChild is one top-level element under <detail>.
type detailChild struct {
	start    xml.StartElement
	raw      []byte
	hasInner bool // child elements or non-space text
}

// splitDetail moves whole, convertible elements into typed fields and keeps
// the rest, byte for byte, as xmlDetail.
func splitDetail(inner []byte) (Detail, error) {
	dec := xml.NewDecoder(bytes.NewReader(inner))
	var children []detailChild
	depth := 0
	var cur *detailChild
	var startOff int64
	for {
		off := dec.InputOffset()
		tok, err := dec.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return Detail{}, err
		}
		switch t := tok.(type) {
		case xml.StartElement:
			if depth == 0 {
				cur = &detailChild{start: t.Copy()}
				startOff = off
			} else {
				cur.hasInner = true
			}
			depth++
		case xml.EndElement:
			depth--
			if depth == 0 {
				cur.raw = inner[startOff:dec.InputOffset()]
				children = append(children, *cur)
				cur = nil
			}
		case xml.CharData:
			if depth > 0 && len(bytes.TrimSpace(t)) > 0 {
				cur.hasInner = true
			}
		}
	}

	counts := map[string]int{}
	for _, c := range children {
		counts[c.start.Name.Local]++
	}

	var d Detail
	var rest bytes.Buffer
	for _, c := range children {
		if c.start.Name.Space != "" || c.hasInner || counts[c.start.Name.Local] != 1 || !d.take(c.start) {
			rest.Write(c.raw)
		}
	}
	d.XML = rest.String()
	return d, nil
}

// attrsOnly returns the attributes as a map if every one is in allowed.
func attrsOnly(se xml.StartElement, allowed ...string) (map[string]string, bool) {
	m := map[string]string{}
	for _, a := range se.Attr {
		ok := false
		for _, name := range allowed {
			if a.Name.Space == "" && a.Name.Local == name {
				ok = true
			}
		}
		if !ok {
			return nil, false
		}
		m[a.Name.Local] = a.Value
	}
	return m, true
}

// take converts se into its typed field. It returns false, leaving the
// element for xmlDetail, if the element does not map exactly.
func (d *Detail) take(se xml.StartElement) bool {
	switch se.Name.Local {
	case "contact":
		if a, ok := attrsOnly(se, "endpoint", "callsign"); ok {
			d.Contact = &Contact{Endpoint: a["endpoint"], Callsign: a["callsign"]}
			return true
		}
	case "__group":
		if a, ok := attrsOnly(se, "name", "role"); ok {
			d.Group = &Group{Name: a["name"], Role: a["role"]}
			return true
		}
	case "precisionlocation":
		if a, ok := attrsOnly(se, "geopointsrc", "altsrc"); ok {
			d.PrecisionLocation = &PrecisionLocation{GeoPointSrc: a["geopointsrc"], AltSrc: a["altsrc"]}
			return true
		}
	case "status":
		if a, ok := attrsOnly(se, "battery"); ok {
			if v, err := strconv.ParseUint(a["battery"], 10, 32); err == nil {
				d.Status = &Status{Battery: uint32(v)}
				return true
			}
		}
	case "takv":
		if a, ok := attrsOnly(se, "device", "platform", "os", "version"); ok {
			d.Takv = &Takv{Device: a["device"], Platform: a["platform"], OS: a["os"], Version: a["version"]}
			return true
		}
	case "track":
		if a, ok := attrsOnly(se, "speed", "course"); ok {
			sp, err1 := strconv.ParseFloat(a["speed"], 64)
			co, err2 := strconv.ParseFloat(a["course"], 64)
			if err1 == nil && err2 == nil {
				d.Track = &Track{Speed: sp, Course: co}
				return true
			}
		}
	}
	return false
}

// ── Event → XML ──────────────────────────────────────────────────────────────

type attrWriter struct{ b *bytes.Buffer }

func (w attrWriter) str(name, v string) {
	w.b.WriteString(" " + name + `="`)
	_ = xml.EscapeText(w.b, []byte(v))
	w.b.WriteByte('"')
}

func (w attrWriter) opt(name, v string) {
	if v != "" {
		w.str(name, v)
	}
}

func (w attrWriter) num(name string, v float64) {
	w.str(name, strconv.FormatFloat(v, 'f', -1, 64))
}

// MarshalCotXML renders the event as CoT XML. withHeader adds the
// <?xml?> declaration that mesh and streaming senders put first.
func (e *Event) MarshalCotXML(withHeader bool) []byte {
	var b bytes.Buffer
	if withHeader {
		b.WriteString(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>`)
	}
	w := attrWriter{&b}
	b.WriteString("<event")
	w.str("version", "2.0")
	w.str("uid", e.UID)
	w.str("type", e.Type)
	w.opt("access", e.Access)
	w.opt("caveat", e.Caveat)
	w.opt("releasableTo", e.ReleasableTo)
	w.opt("qos", e.QoS)
	w.opt("opex", e.Opex)
	w.str("how", e.How)
	w.str("time", formatCotTime(e.SendTime))
	w.str("start", formatCotTime(e.StartTime))
	w.str("stale", formatCotTime(e.StaleTime))
	b.WriteString("><point")
	w.num("lat", e.Lat)
	w.num("lon", e.Lon)
	w.num("hae", e.HAE)
	w.num("ce", e.CE)
	w.num("le", e.LE)
	b.WriteString("/><detail>")

	d := e.Detail
	b.WriteString(d.XML)
	// detail.proto receiver rule 5: if xmlDetail already carries an element,
	// the typed copy is ignored.
	inXML := func(name string) bool {
		for _, end := range []string{" ", "/", ">", "\t", "\n", "\r"} {
			if strings.Contains(d.XML, "<"+name+end) {
				return true
			}
		}
		return false
	}
	if c := d.Contact; c != nil && !inXML("contact") {
		b.WriteString("<contact")
		w.opt("endpoint", c.Endpoint)
		w.opt("callsign", c.Callsign)
		b.WriteString("/>")
	}
	if g := d.Group; g != nil && !inXML("__group") {
		b.WriteString("<__group")
		w.opt("name", g.Name)
		w.opt("role", g.Role)
		b.WriteString("/>")
	}
	if p := d.PrecisionLocation; p != nil && !inXML("precisionlocation") {
		b.WriteString("<precisionlocation")
		w.opt("geopointsrc", p.GeoPointSrc)
		w.opt("altsrc", p.AltSrc)
		b.WriteString("/>")
	}
	if s := d.Status; s != nil && !inXML("status") {
		b.WriteString("<status")
		w.str("battery", strconv.FormatUint(uint64(s.Battery), 10))
		b.WriteString("/>")
	}
	if t := d.Takv; t != nil && !inXML("takv") {
		b.WriteString("<takv")
		w.opt("device", t.Device)
		w.opt("platform", t.Platform)
		w.opt("os", t.OS)
		w.opt("version", t.Version)
		b.WriteString("/>")
	}
	if t := d.Track; t != nil && !inXML("track") {
		b.WriteString("<track")
		w.num("speed", t.Speed)
		w.num("course", t.Course)
		b.WriteString("/>")
	}
	b.WriteString("</detail></event>")
	return b.Bytes()
}
