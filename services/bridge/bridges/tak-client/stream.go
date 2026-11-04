// SPDX-License-Identifier: Apache-2.0

package main

// TAK Server streaming framing and protocol negotiation, written from the
// published spec (atak-civ takproto/README.txt, "Streaming" sections):
//
//   - Version 0 (XML): events back to back, each starting with an XML
//     declaration; a reader splits on </event>.
//   - The server MAY offer TAK Protocol v1 with a t-x-takp-v event whose
//     <TakControl> holds <TakProtocolSupport version="1"/>. The client asks with
//     t-x-takp-q (<TakRequest version="1"/>) and stops sending XML; the server
//     answers t-x-takp-r (<TakResponse status="true|false"/>). On true both
//     sides switch to streaming framing: 0xbf <payload length varint> <TakMessage>.
//   - The client waits at least one minute for the answer, then disconnects.

import (
	"bufio"
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"strings"
	"time"

	"google.golang.org/protobuf/encoding/protowire"
)

const (
	typeProtoOffer    = "t-x-takp-v"
	typeProtoRequest  = "t-x-takp-q"
	typeProtoResponse = "t-x-takp-r"
	typePing          = "t-x-c-t"
	typePong          = "t-x-c-t-r"

	// maxEventBytes bounds one CoT event in either framing, so a broken or
	// hostile peer cannot grow the read buffer without limit.
	maxEventBytes = 1 << 20
)

var (
	errEventTooLarge = errors.New("CoT event exceeds 1 MiB")
	errBadFrame      = errors.New("bad TAK Protocol v1 stream frame")
	eventEnd         = []byte("</event>")
	eventStart       = []byte("<event")
)

// EncodeStream frames a version 1 payload for a streaming (server) connection.
func EncodeStream(m *TakMessage) []byte {
	p := m.Marshal()
	b := make([]byte, 0, len(p)+6)
	b = append(b, takMagic)
	b = protowire.AppendVarint(b, uint64(len(p)))
	return append(b, p...)
}

// xmlReader splits a version 0 stream into events.
type xmlReader struct {
	r   *bufio.Reader
	buf []byte
}

func newXMLReader(r *bufio.Reader) *xmlReader { return &xmlReader{r: r} }

// Next returns the next <event>…</event> (anything before <event, such as the
// XML declaration or whitespace, is dropped).
func (x *xmlReader) Next() ([]byte, error) {
	chunk := make([]byte, 8192)
	for {
		if i := bytes.Index(x.buf, eventEnd); i >= 0 {
			end := i + len(eventEnd)
			raw := x.buf[:end]
			x.buf = append([]byte(nil), x.buf[end:]...)
			if s := bytes.Index(raw, eventStart); s >= 0 {
				return raw[s:], nil
			}
			continue // stray </event> with no start: skip it
		}
		if len(x.buf) > maxEventBytes {
			return nil, errEventTooLarge
		}
		n, err := x.r.Read(chunk)
		x.buf = append(x.buf, chunk[:n]...)
		if err != nil {
			if n > 0 && err == io.EOF {
				continue
			}
			return nil, err
		}
	}
}

// Leftover returns bytes read past the last event. After a successful
// negotiation they are the start of the version 1 stream.
func (x *xmlReader) Leftover() []byte { return x.buf }

// readStreamMessage reads one version 1 streaming frame.
func readStreamMessage(r *bufio.Reader) (*TakMessage, error) {
	magic, err := r.ReadByte()
	if err != nil {
		return nil, err
	}
	if magic != takMagic {
		return nil, fmt.Errorf("%w: starts with 0x%02x, want 0xbf", errBadFrame, magic)
	}
	var n uint64
	for shift := uint(0); ; shift += 7 {
		if shift >= 70 {
			return nil, fmt.Errorf("%w: length varint longer than 10 bytes", errBadFrame)
		}
		c, err := r.ReadByte()
		if err != nil {
			return nil, err
		}
		n |= uint64(c&0x7f) << shift
		if c < 0x80 {
			break
		}
	}
	if n > maxEventBytes {
		return nil, errEventTooLarge
	}
	p := make([]byte, n)
	if _, err := io.ReadFull(r, p); err != nil {
		return nil, err
	}
	return UnmarshalTakMessage(p)
}

// takControlXML is the <TakControl> element in negotiation events.
type takControlXML struct {
	XMLName xml.Name `xml:"TakControl"`
	Support []struct {
		Version string `xml:"version,attr"`
	} `xml:"TakProtocolSupport"`
	Response *struct {
		Status string `xml:"status,attr"`
	} `xml:"TakResponse"`
	ServerVersion *struct {
		ServerVersion string `xml:"serverVersion,attr"`
	} `xml:"TakServerVersionInfo"`
}

func parseTakControl(detailXML string) (*takControlXML, error) {
	i := strings.Index(detailXML, "<TakControl")
	if i < 0 {
		return nil, errors.New("no <TakControl> in detail")
	}
	var tc takControlXML
	if err := xml.NewDecoder(strings.NewReader(detailXML[i:])).Decode(&tc); err != nil {
		return nil, fmt.Errorf("TakControl: %w", err)
	}
	return &tc, nil
}

// offersV1 reports whether a t-x-takp-v event offers version 1, and the
// server version string it carries (TAK Server sends it; empty if absent).
func offersV1(ev *Event) (bool, string) {
	tc, err := parseTakControl(ev.Detail.XML)
	if err != nil {
		return false, ""
	}
	sv := ""
	if tc.ServerVersion != nil {
		sv = tc.ServerVersion.ServerVersion
	}
	for _, s := range tc.Support {
		if strings.TrimSpace(s.Version) == "1" {
			return true, sv
		}
	}
	return false, sv
}

// responseAccepted reports the status of a t-x-takp-r event.
func responseAccepted(ev *Event) bool {
	tc, err := parseTakControl(ev.Detail.XML)
	return err == nil && tc.Response != nil && strings.EqualFold(strings.TrimSpace(tc.Response.Status), "true")
}

// controlEvent builds a protocol or keep-alive event from this client.
func controlEvent(uid, cotType, detailXML string, now time.Time, stale time.Duration) *Event {
	ms := uint64(now.UnixMilli())
	return &Event{
		Type: cotType, UID: uid, How: "m-g",
		SendTime: ms, StartTime: ms, StaleTime: ms + uint64(stale.Milliseconds()),
		CE: 999999, LE: 999999, HAE: 999999,
		Detail: Detail{XML: detailXML},
	}
}

// xmlDecl starts every version 0 streaming event, followed by a newline (spec).
const xmlDecl = `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>` + "\n"

// streamXML renders one event for a version 0 stream.
func streamXML(ev *Event) []byte {
	return append([]byte(xmlDecl), ev.MarshalCotXML(false)...)
}

func protoRequest(uid string, now time.Time) []byte {
	return streamXML(controlEvent(uid, typeProtoRequest, `<TakControl><TakRequest version="1"/></TakControl>`, now, time.Minute))
}

// isControlType reports CoT that is stream plumbing, never published or injected.
func isControlType(cotType string) bool {
	return strings.HasPrefix(cotType, "t-x-takp") || cotType == typePing || cotType == typePong
}
