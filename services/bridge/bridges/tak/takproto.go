// SPDX-License-Identifier: Apache-2.0

package main

// TAK Protocol version 1 codec, written from the published wire spec
// (atak-civ takproto/README.txt and the field numbers in its .proto files).
// No GPL source is copied; protowire does the low-level encoding.
//
// Mesh message:      0xbf <varint version> 0xbf <TakMessage protobuf>
// Streaming message: 0xbf <varint length>  <TakMessage protobuf>

import (
	"errors"
	"fmt"
	"math"

	"google.golang.org/protobuf/encoding/protowire"
)

const takMagic = 0xbf

// TakControl advertises the protocol versions a mesh participant can decode.
type TakControl struct {
	MinProtoVersion uint32
	MaxProtoVersion uint32
	ContactUID      string
}

// Event is the protocol-neutral form of one CoT event. XML and protobuf
// both convert to and from it. Times are unix milliseconds.
type Event struct {
	Type, Access, Caveat, ReleasableTo, QoS, Opex string
	UID, How                                      string
	SendTime, StartTime, StaleTime                uint64
	Lat, Lon, HAE, CE, LE                         float64
	Detail                                        Detail
}

// Detail holds the typed detail elements plus whatever XML is left over.
// A typed field is set only when the whole XML element maps onto it.
type Detail struct {
	XML               string // inner XML of <detail>, typed elements removed
	Contact           *Contact
	Group             *Group
	PrecisionLocation *PrecisionLocation
	Status            *Status
	Takv              *Takv
	Track             *Track
	// Extensions are registered binary encodings we cannot decode. They are
	// kept so the message can be relayed, and counted so callers can log it.
	Extensions int
}

type Contact struct{ Endpoint, Callsign string }
type Group struct{ Name, Role string }
type PrecisionLocation struct{ GeoPointSrc, AltSrc string }
type Status struct{ Battery uint32 }
type Takv struct{ Device, Platform, OS, Version string }
type Track struct{ Speed, Course float64 }

// TakMessage is the version 1 payload. Either part may be nil.
type TakMessage struct {
	Control *TakControl
	Event   *Event
}

// ── Framing ──────────────────────────────────────────────────────────────────

// IsTakProtocol reports whether a datagram starts with the TAK Protocol magic
// byte. Legacy XML always starts with '<' (or whitespace), never 0xbf.
func IsTakProtocol(b []byte) bool { return len(b) > 0 && b[0] == takMagic }

// DecodeMeshHeader strips the mesh header and returns the protocol version and payload.
func DecodeMeshHeader(b []byte) (uint64, []byte, error) {
	if len(b) < 3 || b[0] != takMagic {
		return 0, nil, errors.New("missing TAK protocol magic byte")
	}
	v, n := protowire.ConsumeVarint(b[1:])
	if n < 0 {
		return 0, nil, errors.New("bad protocol version varint")
	}
	if len(b) < 1+n+1 || b[1+n] != takMagic {
		return 0, nil, errors.New("missing closing magic byte")
	}
	return v, b[2+n:], nil
}

// EncodeMesh frames a version 1 payload for a mesh datagram or TCP send-once.
func EncodeMesh(m *TakMessage) []byte {
	b := []byte{takMagic}
	b = protowire.AppendVarint(b, 1)
	b = append(b, takMagic)
	return append(b, m.Marshal()...)
}

// ── Encoding ─────────────────────────────────────────────────────────────────

func appendString(b []byte, num protowire.Number, s string) []byte {
	if s == "" {
		return b
	}
	b = protowire.AppendTag(b, num, protowire.BytesType)
	return protowire.AppendString(b, s)
}

func appendUint(b []byte, num protowire.Number, v uint64) []byte {
	if v == 0 {
		return b
	}
	b = protowire.AppendTag(b, num, protowire.VarintType)
	return protowire.AppendVarint(b, v)
}

func appendDouble(b []byte, num protowire.Number, v float64) []byte {
	if v == 0 {
		return b
	}
	b = protowire.AppendTag(b, num, protowire.Fixed64Type)
	return protowire.AppendFixed64(b, math.Float64bits(v))
}

func appendMessage(b []byte, num protowire.Number, msg []byte) []byte {
	b = protowire.AppendTag(b, num, protowire.BytesType)
	return protowire.AppendBytes(b, msg)
}

// Marshal serializes the TakMessage protobuf (no framing).
func (m *TakMessage) Marshal() []byte {
	var b []byte
	if c := m.Control; c != nil {
		var cb []byte
		cb = appendUint(cb, 1, uint64(c.MinProtoVersion))
		cb = appendUint(cb, 2, uint64(c.MaxProtoVersion))
		cb = appendString(cb, 3, c.ContactUID)
		b = appendMessage(b, 1, cb)
	}
	if e := m.Event; e != nil {
		b = appendMessage(b, 2, e.marshal())
	}
	return b
}

func (e *Event) marshal() []byte {
	var b []byte
	b = appendString(b, 1, e.Type)
	b = appendString(b, 2, e.Access)
	b = appendString(b, 3, e.QoS)
	b = appendString(b, 4, e.Opex)
	b = appendString(b, 5, e.UID)
	b = appendUint(b, 6, e.SendTime)
	b = appendUint(b, 7, e.StartTime)
	b = appendUint(b, 8, e.StaleTime)
	b = appendString(b, 9, e.How)
	b = appendDouble(b, 10, e.Lat)
	b = appendDouble(b, 11, e.Lon)
	b = appendDouble(b, 12, e.HAE)
	b = appendDouble(b, 13, e.CE)
	b = appendDouble(b, 14, e.LE)
	if d := e.Detail.marshal(); len(d) > 0 {
		b = appendMessage(b, 15, d)
	}
	b = appendString(b, 16, e.Caveat)
	b = appendString(b, 17, e.ReleasableTo)
	return b
}

func (d *Detail) marshal() []byte {
	var b []byte
	b = appendString(b, 1, d.XML)
	if c := d.Contact; c != nil {
		b = appendMessage(b, 2, appendString(appendString(nil, 1, c.Endpoint), 2, c.Callsign))
	}
	if g := d.Group; g != nil {
		b = appendMessage(b, 3, appendString(appendString(nil, 1, g.Name), 2, g.Role))
	}
	if p := d.PrecisionLocation; p != nil {
		b = appendMessage(b, 4, appendString(appendString(nil, 1, p.GeoPointSrc), 2, p.AltSrc))
	}
	if s := d.Status; s != nil {
		b = appendMessage(b, 5, appendUint(nil, 1, uint64(s.Battery)))
	}
	if t := d.Takv; t != nil {
		var tb []byte
		tb = appendString(tb, 1, t.Device)
		tb = appendString(tb, 2, t.Platform)
		tb = appendString(tb, 3, t.OS)
		tb = appendString(tb, 4, t.Version)
		b = appendMessage(b, 6, tb)
	}
	if t := d.Track; t != nil {
		b = appendMessage(b, 7, appendDouble(appendDouble(nil, 1, t.Speed), 2, t.Course))
	}
	return b
}

// ── Decoding ─────────────────────────────────────────────────────────────────

// field walks every field of a protobuf message and calls fn with its number,
// wire type and raw value. Unknown fields are skipped, as proto3 requires.
func fields(b []byte, fn func(num protowire.Number, typ protowire.Type, v []byte) error) error {
	for len(b) > 0 {
		num, typ, n := protowire.ConsumeTag(b)
		if n < 0 {
			return protowire.ParseError(n)
		}
		b = b[n:]
		m := protowire.ConsumeFieldValue(num, typ, b)
		if m < 0 {
			return protowire.ParseError(m)
		}
		if err := fn(num, typ, b[:m]); err != nil {
			return err
		}
		b = b[m:]
	}
	return nil
}

func asString(typ protowire.Type, v []byte) (string, error) {
	if typ != protowire.BytesType {
		return "", fmt.Errorf("want bytes, got wire type %d", typ)
	}
	s, n := protowire.ConsumeString(v)
	if n < 0 {
		return "", protowire.ParseError(n)
	}
	return s, nil
}

func asBytes(typ protowire.Type, v []byte) ([]byte, error) {
	if typ != protowire.BytesType {
		return nil, fmt.Errorf("want bytes, got wire type %d", typ)
	}
	bs, n := protowire.ConsumeBytes(v)
	if n < 0 {
		return nil, protowire.ParseError(n)
	}
	return bs, nil
}

func asUint(typ protowire.Type, v []byte) (uint64, error) {
	if typ != protowire.VarintType {
		return 0, fmt.Errorf("want varint, got wire type %d", typ)
	}
	u, n := protowire.ConsumeVarint(v)
	if n < 0 {
		return 0, protowire.ParseError(n)
	}
	return u, nil
}

func asDouble(typ protowire.Type, v []byte) (float64, error) {
	if typ != protowire.Fixed64Type {
		return 0, fmt.Errorf("want fixed64, got wire type %d", typ)
	}
	u, n := protowire.ConsumeFixed64(v)
	if n < 0 {
		return 0, protowire.ParseError(n)
	}
	return math.Float64frombits(u), nil
}

// stringPair decodes the two-string messages (Contact, Group, PrecisionLocation).
func stringPair(b []byte) (a, c string, err error) {
	err = fields(b, func(num protowire.Number, typ protowire.Type, v []byte) error {
		s, err := asString(typ, v)
		switch num {
		case 1:
			a = s
		case 2:
			c = s
		}
		return err
	})
	return
}

// UnmarshalTakMessage parses a version 1 payload (framing already removed).
func UnmarshalTakMessage(b []byte) (*TakMessage, error) {
	m := &TakMessage{}
	err := fields(b, func(num protowire.Number, typ protowire.Type, v []byte) error {
		switch num {
		case 1:
			raw, err := asBytes(typ, v)
			if err != nil {
				return err
			}
			c := &TakControl{}
			m.Control = c
			return fields(raw, func(num protowire.Number, typ protowire.Type, v []byte) error {
				var err error
				var u uint64
				switch num {
				case 1:
					u, err = asUint(typ, v)
					c.MinProtoVersion = uint32(u)
				case 2:
					u, err = asUint(typ, v)
					c.MaxProtoVersion = uint32(u)
				case 3:
					c.ContactUID, err = asString(typ, v)
				}
				return err
			})
		case 2:
			raw, err := asBytes(typ, v)
			if err != nil {
				return err
			}
			m.Event, err = unmarshalEvent(raw)
			return err
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("takmessage: %w", err)
	}
	return m, nil
}

func unmarshalEvent(b []byte) (*Event, error) {
	e := &Event{}
	err := fields(b, func(num protowire.Number, typ protowire.Type, v []byte) error {
		var err error
		switch num {
		case 1:
			e.Type, err = asString(typ, v)
		case 2:
			e.Access, err = asString(typ, v)
		case 3:
			e.QoS, err = asString(typ, v)
		case 4:
			e.Opex, err = asString(typ, v)
		case 5:
			e.UID, err = asString(typ, v)
		case 6:
			e.SendTime, err = asUint(typ, v)
		case 7:
			e.StartTime, err = asUint(typ, v)
		case 8:
			e.StaleTime, err = asUint(typ, v)
		case 9:
			e.How, err = asString(typ, v)
		case 10:
			e.Lat, err = asDouble(typ, v)
		case 11:
			e.Lon, err = asDouble(typ, v)
		case 12:
			e.HAE, err = asDouble(typ, v)
		case 13:
			e.CE, err = asDouble(typ, v)
		case 14:
			e.LE, err = asDouble(typ, v)
		case 15:
			var raw []byte
			if raw, err = asBytes(typ, v); err == nil {
				err = e.Detail.unmarshal(raw)
			}
		case 16:
			e.Caveat, err = asString(typ, v)
		case 17:
			e.ReleasableTo, err = asString(typ, v)
		}
		return err
	})
	return e, err
}

func (d *Detail) unmarshal(b []byte) error {
	return fields(b, func(num protowire.Number, typ protowire.Type, v []byte) error {
		if num == 1 {
			s, err := asString(typ, v)
			d.XML = s
			return err
		}
		raw, err := asBytes(typ, v)
		if err != nil {
			return err
		}
		switch num {
		case 2:
			c := &Contact{}
			c.Endpoint, c.Callsign, err = stringPair(raw)
			d.Contact = c
		case 3:
			g := &Group{}
			g.Name, g.Role, err = stringPair(raw)
			d.Group = g
		case 4:
			p := &PrecisionLocation{}
			p.GeoPointSrc, p.AltSrc, err = stringPair(raw)
			d.PrecisionLocation = p
		case 5:
			s := &Status{}
			d.Status = s
			err = fields(raw, func(num protowire.Number, typ protowire.Type, v []byte) error {
				if num != 1 {
					return nil
				}
				u, err := asUint(typ, v)
				s.Battery = uint32(u)
				return err
			})
		case 6:
			t := &Takv{}
			d.Takv = t
			err = fields(raw, func(num protowire.Number, typ protowire.Type, v []byte) error {
				s, err := asString(typ, v)
				switch num {
				case 1:
					t.Device = s
				case 2:
					t.Platform = s
				case 3:
					t.OS = s
				case 4:
					t.Version = s
				}
				return err
			})
		case 7:
			t := &Track{}
			d.Track = t
			err = fields(raw, func(num protowire.Number, typ protowire.Type, v []byte) error {
				f, err := asDouble(typ, v)
				switch num {
				case 1:
					t.Speed = f
				case 2:
					t.Course = f
				}
				return err
			})
		case 8:
			d.Extensions++
		}
		return err
	})
}
