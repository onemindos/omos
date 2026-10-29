// Package cot parses ATAK GeoChat events and renders the two CoT events agentops sends for
// an agent: its contact presence and a directed GeoChat reply.
//
// GeoChat layout: the `__chat`/`chatgrp`/`remarks` detail used by ATAK, as in the mesh
// gateway's tests (omos PR #551, services/bridge/bridges/tak/gateway_test.go lines 146-148)
// and onemind-xr docs/research/tak-protocols.md §1. The `b-t-f` type prefix and the
// directed-chat attributes are from general TAK usage and are *unverified* against an ATAK
// device in this run (onemind-xr docs/contracts/cot-projection.md §13); the owner checklist
// has the device check.
package cot

import (
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"strings"
	"time"
)

// MaxEventBytes bounds what Parse accepts.
const MaxEventBytes = 64 << 10

// Chat is one GeoChat message.
type Chat struct {
	EventUID       string
	Type           string
	Time           time.Time
	MessageID      string
	SenderUID      string
	SenderCallsign string
	Chatroom       string // room name, or the peer's callsign for a direct chat
	ChatID         string // room id, or the peer's uid for a direct chat
	ToUID          string // remarks@to
	GroupUIDs      []string
	DestCallsigns  []string // marti/dest@callsign
	Text           string
}

type xEvent struct {
	XMLName xml.Name `xml:"event"`
	UID     string   `xml:"uid,attr"`
	Type    string   `xml:"type,attr"`
	Time    string   `xml:"time,attr"`
	Detail  struct {
		Chat *struct {
			ID             string `xml:"id,attr"`
			Chatroom       string `xml:"chatroom,attr"`
			SenderCallsign string `xml:"senderCallsign,attr"`
			MessageID      string `xml:"messageId,attr"`
			Group          *struct {
				Attrs []xml.Attr `xml:",any,attr"`
			} `xml:"chatgrp"`
		} `xml:"__chat"`
		Link []struct {
			UID      string `xml:"uid,attr"`
			Relation string `xml:"relation,attr"`
		} `xml:"link"`
		Remarks *struct {
			Source string `xml:"source,attr"`
			To     string `xml:"to,attr"`
			Text   string `xml:",chardata"`
		} `xml:"remarks"`
		Marti *struct {
			Dest []struct {
				Callsign string `xml:"callsign,attr"`
				UID      string `xml:"uid,attr"`
			} `xml:"dest"`
		} `xml:"marti"`
	} `xml:"detail"`
}

// IsChatType reports whether a CoT type is GeoChat (`b-t-f` and `b-t-f-*`).
func IsChatType(t string) bool { return t == "b-t-f" || strings.HasPrefix(t, "b-t-f-") }

// ErrNotChat means the event is not GeoChat.
var ErrNotChat = errors.New("not a GeoChat event")

// ParseChat decodes one CoT XML event as GeoChat.
func ParseChat(data []byte) (*Chat, error) {
	if len(data) > MaxEventBytes {
		return nil, errors.New("event too large")
	}
	if bytes.Contains(data, []byte("<!DOCTYPE")) || bytes.Contains(data, []byte("<!ENTITY")) {
		return nil, errors.New("DOCTYPE or ENTITY not allowed")
	}
	var ev xEvent
	dec := xml.NewDecoder(bytes.NewReader(data))
	dec.Strict = true
	if err := dec.Decode(&ev); err != nil {
		return nil, err
	}
	if !IsChatType(ev.Type) {
		return nil, ErrNotChat
	}
	c := &Chat{EventUID: ev.UID, Type: ev.Type}
	if t, err := time.Parse(time.RFC3339Nano, ev.Time); err == nil {
		c.Time = t
	}
	if ch := ev.Detail.Chat; ch != nil {
		c.ChatID, c.Chatroom, c.SenderCallsign, c.MessageID = ch.ID, ch.Chatroom, ch.SenderCallsign, ch.MessageID
		if ch.Group != nil {
			for _, a := range ch.Group.Attrs {
				if strings.HasPrefix(a.Name.Local, "uid") {
					c.GroupUIDs = append(c.GroupUIDs, a.Value)
				}
			}
		}
	}
	for _, l := range ev.Detail.Link {
		if l.Relation == "p-p" && c.SenderUID == "" {
			c.SenderUID = l.UID
		}
	}
	if r := ev.Detail.Remarks; r != nil {
		c.Text = strings.TrimSpace(r.Text)
		c.ToUID = r.To
		if c.SenderUID == "" && strings.HasPrefix(r.Source, "BAO.F.") {
			if i := strings.LastIndex(r.Source, "."); i > 0 {
				c.SenderUID = r.Source[i+1:]
			}
		}
	}
	if c.SenderUID == "" && len(c.GroupUIDs) > 0 {
		c.SenderUID = c.GroupUIDs[0] // chatgrp uid0 is the sender in ATAK's layout
	}
	if m := ev.Detail.Marti; m != nil {
		for _, d := range m.Dest {
			if d.Callsign != "" {
				c.DestCallsigns = append(c.DestCallsigns, d.Callsign)
			}
		}
	}
	if c.SenderUID == "" {
		return nil, errors.New("GeoChat without a sender uid")
	}
	return c, nil
}

// AddressedTo reports whether the chat is for the contact with this uid or callsign. A
// group message counts only if it starts with "@CALLSIGN".
func (c *Chat) AddressedTo(uid, callsign string) bool {
	if uid != "" && (c.ChatID == uid || c.ToUID == uid) {
		return true
	}
	if uid != "" && len(c.GroupUIDs) == 2 && c.GroupUIDs[1] == uid {
		return true
	}
	for _, d := range c.DestCallsigns {
		if callsign != "" && strings.EqualFold(d, callsign) {
			return true
		}
	}
	if callsign != "" {
		t := strings.TrimSpace(c.Text)
		p := "@" + callsign
		if len(t) >= len(p) && strings.EqualFold(t[:len(p)], p) {
			return true
		}
	}
	return false
}

// Contact is an agent as TAK sees it.
type Contact struct {
	UID      string // OM.agent.{entity_id}: the OM. prefix makes the entity service drop it as an echo
	Callsign string
	CoTType  string // default a-f-G-U-C (unverified choice, see README)
	Lat, Lon float64
	Team     string // __group name, default Cyan
	Role     string // __group role, default HQ
	Remarks  string
	Version  string
}

func esc(s string) string {
	var b bytes.Buffer
	_ = xml.EscapeText(&b, []byte(s))
	return b.String()
}

func ts(t time.Time) string { return t.UTC().Format("2006-01-02T15:04:05.000Z") }

// RenderContact renders the agent's presence: a PLI with `<contact endpoint>` and `<takv>`
// so ATAK lists it as a contact that can be messaged. Point accuracy is unknown on purpose
// (ce/le 9999999): an agent has no real position.
func RenderContact(c Contact, now time.Time, stale time.Duration) []byte {
	typ := c.CoTType
	if typ == "" {
		typ = "a-f-G-U-C"
	}
	team, role := c.Team, c.Role
	if team == "" {
		team = "Cyan"
	}
	if role == "" {
		role = "HQ"
	}
	var b strings.Builder
	fmt.Fprintf(&b, `<event version="2.0" uid="%s" type="%s" how="m-g" time="%s" start="%s" stale="%s">`,
		esc(c.UID), esc(typ), ts(now), ts(now), ts(now.Add(stale)))
	fmt.Fprintf(&b, `<point lat="%.6f" lon="%.6f" hae="9999999" ce="9999999" le="9999999"/>`, c.Lat, c.Lon)
	b.WriteString(`<detail>`)
	fmt.Fprintf(&b, `<contact callsign="%s" endpoint="*:-1:stcp"/>`, esc(c.Callsign))
	fmt.Fprintf(&b, `<__group name="%s" role="%s"/>`, esc(team), esc(role))
	fmt.Fprintf(&b, `<takv device="omos-agentops" platform="OneMind Agent" os="omos" version="%s"/>`, esc(c.Version))
	fmt.Fprintf(&b, `<uid Droid="%s"/>`, esc(c.Callsign))
	if c.Remarks != "" {
		fmt.Fprintf(&b, `<remarks>%s</remarks>`, esc(c.Remarks))
	}
	b.WriteString(`<_flow-tags_ omos-agentops="` + ts(now) + `"/>`)
	b.WriteString(`</detail></event>`)
	return []byte(b.String())
}

// Reply is a directed GeoChat from an agent to one TAK user.
type Reply struct {
	From      Contact
	ToUID     string
	ToCall    string
	Text      string
	MessageID string
	Server    bool // add <marti><dest callsign>, so TAK Server delivers it to that user only
}

// RenderChat renders a directed GeoChat reply.
func RenderChat(r Reply, now time.Time) []byte {
	uid := "GeoChat." + r.From.UID + "." + r.ToUID + "." + r.MessageID
	var b strings.Builder
	fmt.Fprintf(&b, `<event version="2.0" uid="%s" type="b-t-f" how="h-g-i-g-o" time="%s" start="%s" stale="%s">`,
		esc(uid), ts(now), ts(now), ts(now.Add(24*time.Hour)))
	fmt.Fprintf(&b, `<point lat="%.6f" lon="%.6f" hae="9999999" ce="9999999" le="9999999"/>`, r.From.Lat, r.From.Lon)
	b.WriteString(`<detail>`)
	fmt.Fprintf(&b, `<__chat parent="RootContactGroup" groupOwner="false" messageId="%s" chatroom="%s" id="%s" senderCallsign="%s">`,
		esc(r.MessageID), esc(r.ToCall), esc(r.ToUID), esc(r.From.Callsign))
	fmt.Fprintf(&b, `<chatgrp uid0="%s" uid1="%s" id="%s"/></__chat>`, esc(r.From.UID), esc(r.ToUID), esc(r.ToUID))
	fmt.Fprintf(&b, `<link uid="%s" type="%s" relation="p-p"/>`, esc(r.From.UID), esc(orDefault(r.From.CoTType, "a-f-G-U-C")))
	fmt.Fprintf(&b, `<remarks source="BAO.F.OneMind.%s" to="%s" time="%s">%s</remarks>`, esc(r.From.UID), esc(r.ToUID), ts(now), esc(r.Text))
	if r.Server && r.ToCall != "" {
		fmt.Fprintf(&b, `<marti><dest callsign="%s"/></marti>`, esc(r.ToCall))
	}
	b.WriteString(`<_flow-tags_ omos-agentops="` + ts(now) + `"/>`)
	b.WriteString(`</detail></event>`)
	return []byte(b.String())
}

func orDefault(s, d string) string {
	if s == "" {
		return d
	}
	return s
}
