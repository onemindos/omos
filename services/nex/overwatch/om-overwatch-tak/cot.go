package main

import (
	"encoding/xml"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// cot.go — the Cursor-on-Target (CoT) XML layer for om-overwatch-tak.
//
// CoT is MITRE's situational-awareness event schema used by ATAK/iTAK/WinTAK.
// We exchange plain CoT XML over NATS with the TAK Server's NatsBridgePlugin:
//   tak.cot.out    — every CoT the TAK server sees (devices, operators) → us
//   tak.cot.inject — CoT we publish → pushed to every ATAK client
//
// This file defines typed CoT structs (compile-time schema safety), a parser,
// a renderer, and the affiliation/entity-type mappings in both directions.
// Ported from cot2constellation/bridge.py — preserved byte-identical alongside
// this organ as bridge.py.ref (source repo retired 2026-07-02 after full review).

// ── Typed CoT schema ─────────────────────────────────────────────────────────

// CoTEvent is a Cursor-on-Target <event>.
type CoTEvent struct {
	XMLName xml.Name   `xml:"event"`
	Version string     `xml:"version,attr"`
	UID     string     `xml:"uid,attr"`
	Type    string     `xml:"type,attr"`
	How     string     `xml:"how,attr,omitempty"`
	Time    string     `xml:"time,attr"`
	Start   string     `xml:"start,attr"`
	Stale   string     `xml:"stale,attr"`
	Point   CoTPoint   `xml:"point"`
	Detail  *CoTDetail `xml:"detail,omitempty"`
}

// CoTPoint is the geographic position. hae = height above ellipsoid (m),
// ce = circular error (m), le = linear error (m).
type CoTPoint struct {
	Lat float64 `xml:"lat,attr"`
	Lon float64 `xml:"lon,attr"`
	HAE float64 `xml:"hae,attr"`
	CE  float64 `xml:"ce,attr"`
	LE  float64 `xml:"le,attr"`
}

// ff formats a float in fixed-point (never scientific) notation — ATAK's CoT
// parser rejects attributes like "9.999999e+06".
func ff(f float64) string { return strconv.FormatFloat(f, 'f', -1, 64) }

// MarshalXML renders <point> with fixed-point attributes so ATAK can parse it.
func (p CoTPoint) MarshalXML(enc *xml.Encoder, start xml.StartElement) error {
	start.Name = xml.Name{Local: "point"}
	start.Attr = []xml.Attr{
		{Name: xml.Name{Local: "lat"}, Value: ff(p.Lat)},
		{Name: xml.Name{Local: "lon"}, Value: ff(p.Lon)},
		{Name: xml.Name{Local: "hae"}, Value: ff(p.HAE)},
		{Name: xml.Name{Local: "ce"}, Value: ff(p.CE)},
		{Name: xml.Name{Local: "le"}, Value: ff(p.LE)},
	}
	if err := enc.EncodeToken(start); err != nil {
		return err
	}
	return enc.EncodeToken(start.End())
}

// CoTDetail carries the optional sub-schema (callsign, motion, status, TAK client).
type CoTDetail struct {
	Contact *CoTContact `xml:"contact,omitempty"`
	Track   *CoTTrack   `xml:"track,omitempty"`
	Status  *CoTStatus  `xml:"status,omitempty"`
	Takv    *CoTTakv    `xml:"takv,omitempty"`
	Remarks string      `xml:"remarks,omitempty"`
	// Overwatch marks CoT it originates so the inbound path can ignore its own
	// re-broadcasts (belt-and-suspenders alongside the UID prefix guard).
	Overwatch *CoTOverwatch `xml:"__overwatch,omitempty"`
}

type CoTContact struct {
	Callsign string `xml:"callsign,attr"`
}

type CoTTrack struct {
	Course float64 `xml:"course,attr"` // degrees true
	Speed  float64 `xml:"speed,attr"`  // m/s
}

// MarshalXML renders <track> with fixed-point attributes.
func (t CoTTrack) MarshalXML(enc *xml.Encoder, start xml.StartElement) error {
	start.Name = xml.Name{Local: "track"}
	start.Attr = []xml.Attr{
		{Name: xml.Name{Local: "course"}, Value: ff(t.Course)},
		{Name: xml.Name{Local: "speed"}, Value: ff(t.Speed)},
	}
	if err := enc.EncodeToken(start); err != nil {
		return err
	}
	return enc.EncodeToken(start.End())
}

type CoTStatus struct {
	Battery string `xml:"battery,attr,omitempty"`
}

type CoTTakv struct {
	Platform string `xml:"platform,attr,omitempty"`
	Version  string `xml:"version,attr,omitempty"`
}

type CoTOverwatch struct {
	Fused string `xml:"fused,attr"`
}

// ── Parsing (inbound: CoT → contactReport) ───────────────────────────────────

var tokenSanitizer = regexp.MustCompile(`[^a-zA-Z0-9_-]`)

// sanitizeID turns a CoT uid into a valid NATS subject token.
func sanitizeID(uid string) string {
	s := tokenSanitizer.ReplaceAllString(uid, "_")
	if len(s) > 64 {
		s = s[:64]
	}
	return s
}

// ParseCoT unmarshals a CoT XML document.
func ParseCoT(data []byte) (*CoTEvent, error) {
	var ev CoTEvent
	if err := xml.Unmarshal(data, &ev); err != nil {
		return nil, err
	}
	if ev.XMLName.Local != "event" {
		return nil, fmt.Errorf("not a CoT event: <%s>", ev.XMLName.Local)
	}
	return &ev, nil
}

// IsOverwatchOrigin reports whether this CoT is one of our own outbound
// re-broadcasts (so the inbound path skips it and no feedback loop forms).
func (ev *CoTEvent) IsOverwatchOrigin() bool {
	if strings.HasPrefix(ev.UID, owUIDPrefix) {
		return true
	}
	if ev.Detail != nil && ev.Detail.Overwatch != nil {
		return true
	}
	return false
}

// ToContactReport maps a parsed CoT event into the spine's permissive report
// shape, published on det.cot.{id}. Affiliation is derived from the CoT type
// and treated as authoritative by the spine (human-in-the-loop).
func (ev *CoTEvent) ToContactReport() (id string, report contactReport) {
	id = sanitizeID(ev.UID)

	callsign := ev.UID
	var heading, speed *float64
	if ev.Detail != nil {
		if ev.Detail.Contact != nil && ev.Detail.Contact.Callsign != "" {
			callsign = ev.Detail.Contact.Callsign
		}
		if ev.Detail.Track != nil {
			c := ev.Detail.Track.Course
			s := ev.Detail.Track.Speed
			heading = &c
			speed = &s
		}
	}

	lat := ev.Point.Lat
	lon := ev.Point.Lon
	alt := ev.Point.HAE

	report = contactReport{
		ID:          id,
		EntityType:  cotTypeToEntityType(ev.Type),
		Affiliation: cotTypeToAffiliation(ev.Type),
		Name:        callsign,
		Latitude:    &lat,
		Longitude:   &lon,
		Altitude:    &alt,
		Heading:     heading,
		Velocity:    speed,
		Source:      "cot",
	}
	return id, report
}

// ── Rendering (outbound: Entity → CoT) ───────────────────────────────────────

const (
	owUIDPrefix = "OW." // outbound CoT uid prefix — marks Overwatch-fused tracks
	cotTimeFmt  = "2006-01-02T15:04:05.000Z"
)

// RenderCoT builds a CoT XML document broadcasting a fused entity to ATAK.
// staleFor is how long ATAK should keep the track before ageing it out.
func RenderCoT(e *Entity, staleFor time.Duration) ([]byte, error) {
	now := time.Now().UTC()
	ev := CoTEvent{
		Version: "2.0",
		UID:     owUIDPrefix + e.EntityID,
		Type:    entityToCoTType(e.EntityType, e.Affiliation),
		How:     "m-g", // machine-generated, GPS-derived
		Time:    now.Format(cotTimeFmt),
		Start:   now.Format(cotTimeFmt),
		Stale:   now.Add(staleFor).Format(cotTimeFmt),
		Point: CoTPoint{
			CE: 9999999,
			LE: 9999999,
		},
		Detail: &CoTDetail{
			Overwatch: &CoTOverwatch{Fused: "true"},
		},
	}
	if e.Position != nil {
		ev.Point.Lat = e.Position.Latitude
		ev.Point.Lon = e.Position.Longitude
		ev.Point.HAE = e.Position.Altitude
		if e.Position.Heading != 0 || e.Position.Velocity != 0 {
			ev.Detail.Track = &CoTTrack{Course: e.Position.Heading, Speed: e.Position.Velocity}
		}
	}
	name := e.Name
	if name == "" {
		name = e.EntityID
	}
	ev.Detail.Contact = &CoTContact{Callsign: name}
	if e.Status != "" || len(e.Sources) > 0 {
		ev.Detail.Remarks = fmt.Sprintf("overwatch:%s status=%s sources=%s",
			e.EntityType, e.Status, strings.Join(e.Sources, "+"))
	}

	out, err := xml.Marshal(&ev)
	if err != nil {
		return nil, err
	}
	return append([]byte(xml.Header), out...), nil
}

// RenderCoTDelete builds a "tombstone" CoT that removes a previously-injected
// track from ATAK/iTAK. Re-sending the same UID with a stale time already in the
// past causes TAK clients to age the marker out immediately, so a track we drop
// from the spine (e.g. a spent interceptor after intercept) disappears from the
// map promptly rather than lingering until its normal stale window expires.
func RenderCoTDelete(entityID string) ([]byte, error) {
	now := time.Now().UTC()
	ev := CoTEvent{
		Version: "2.0",
		UID:     owUIDPrefix + entityID,
		Type:    "a-u-G", // type required; removal is matched by UID
		How:     "m-g",
		Time:    now.Format(cotTimeFmt),
		Start:   now.Format(cotTimeFmt),
		Stale:   now.Add(-1 * time.Second).Format(cotTimeFmt), // already stale ⇒ remove
		Point:   CoTPoint{CE: 9999999, LE: 9999999},
		Detail:  &CoTDetail{Overwatch: &CoTOverwatch{Fused: "true"}},
	}
	out, err := xml.Marshal(&ev)
	if err != nil {
		return nil, err
	}
	return append([]byte(xml.Header), out...), nil
}

// ── Affiliation + entity-type mappings ───────────────────────────────────────

// cotTypeToAffiliation reads the affiliation letter (2nd token) of a CoT type
// like "a-f-A-M-F". f=friend h=hostile n=neutral; u/p/a/s/j/k → unknown.
func cotTypeToAffiliation(cotType string) string {
	parts := strings.Split(strings.ToLower(cotType), "-")
	if len(parts) < 2 {
		return AffiliationUnknown
	}
	switch parts[1] {
	case "f":
		return AffiliationFriendly
	case "h", "j", "k", "s": // hostile / joker / faker / suspect → treat as hostile
		return AffiliationHostile
	case "n":
		return AffiliationNeutral
	default: // u (unknown), p (pending), a (assumed friend) → unknown
		return AffiliationUnknown
	}
}

// cotTypeToEntityType reads the battle-dimension (3rd token) of a CoT type and
// maps it to an om-overwatch entity type.
func cotTypeToEntityType(cotType string) string {
	t := strings.ToLower(cotType)
	parts := strings.Split(t, "-")
	if len(parts) < 3 {
		return EntityTypeSensorPlatform
	}
	switch parts[2] {
	case "a": // air
		return EntityTypeAircraftUnknown
	case "g": // ground — distinguish unit (operator) vs equipment (vehicle/sensor)
		if len(parts) >= 5 && parts[3] == "u" {
			return EntityTypeOperatorStation // ground unit = a person/operator
		}
		if len(parts) >= 5 && parts[3] == "e" && strings.HasPrefix(joinFrom(parts, 4), "s") {
			return EntityTypeSensorPlatform // ground equipment sensor
		}
		return EntityTypeGroundVehicle
	case "s": // sea surface
		return EntityTypeSurfaceVessel
	default:
		return EntityTypeSensorPlatform
	}
}

func joinFrom(parts []string, i int) string {
	if i >= len(parts) {
		return ""
	}
	return strings.Join(parts[i:], "-")
}

// entityToCoTType is the reverse map (outbound): om-overwatch entity type +
// affiliation → a MIL-STD-2525-flavoured CoT type string.
func entityToCoTType(entityType, affiliation string) string {
	aff := "u"
	switch affiliation {
	case AffiliationFriendly:
		aff = "f"
	case AffiliationHostile:
		aff = "h"
	case AffiliationNeutral:
		aff = "n"
	}
	switch entityType {
	case EntityTypeAircraftFixedWing:
		return "a-" + aff + "-A-M-F"
	case EntityTypeAircraftMultirotor, EntityTypeAircraftVTOL, EntityTypeAircraftHelicopter:
		return "a-" + aff + "-A-M-H-Q" // rotary-wing
	case EntityTypeAircraftUnknown, EntityTypeInterceptor:
		return "a-" + aff + "-A" // generic air track
	case EntityTypeGroundVehicle:
		return "a-" + aff + "-G-E-V"
	case EntityTypeSurfaceVessel:
		return "a-" + aff + "-S"
	case EntityTypeOperatorStation:
		return "a-" + aff + "-G-U-C" // ground unit combat
	case EntityTypeSensorPlatform:
		return "a-" + aff + "-G-E-S" // ground equipment sensor
	default:
		return "a-" + aff + "-G"
	}
}
