package cot

import (
	"encoding/xml"
	"strings"
	"testing"
	"time"
)

// atakDirect is the GeoChat layout of the mesh gateway tests (PR #551 gateway_test.go
// lines 146-148), addressed to an agent contact instead of "All Chat Rooms".
const atakDirect = `<event version="2.0" uid="GeoChat.ANDROID-1.OM.agent.agent-hermes-zeus-legacy.m1" type="b-t-f" how="h-g-i-g-o" time="2026-10-29T15:04:05.123Z" start="2026-10-29T15:04:05.123Z" stale="2026-10-30T15:04:05.123Z">` +
	`<point lat="30.1" lon="-97.7" hae="0" ce="9999999" le="9999999"/><detail>` +
	`<__chat parent="RootContactGroup" groupOwner="false" messageId="m1" chatroom="HERMES" id="OM.agent.agent-hermes-zeus-legacy" senderCallsign="ZEUS">` +
	`<chatgrp uid0="ANDROID-1" uid1="OM.agent.agent-hermes-zeus-legacy" id="OM.agent.agent-hermes-zeus-legacy"/></__chat>` +
	`<link uid="ANDROID-1" type="a-f-G-U-C" relation="p-p"/>` +
	`<remarks source="BAO.F.ATAK.ANDROID-1" to="OM.agent.agent-hermes-zeus-legacy" time="2026-10-29T15:04:05.123Z">what is sim-x1 doing?</remarks>` +
	`</detail></event>`

func TestParseDirectChat(t *testing.T) {
	c, err := ParseChat([]byte(atakDirect))
	if err != nil {
		t.Fatal(err)
	}
	if c.SenderUID != "ANDROID-1" || c.SenderCallsign != "ZEUS" || c.Text != "what is sim-x1 doing?" || c.MessageID != "m1" {
		t.Fatalf("parsed %+v", c)
	}
	if !c.AddressedTo("OM.agent.agent-hermes-zeus-legacy", "HERMES") {
		t.Fatal("direct chat not addressed to the agent")
	}
	if c.AddressedTo("OM.agent.other", "OTHER") {
		t.Fatal("addressed to the wrong agent")
	}
}

func TestGroupMention(t *testing.T) {
	g := strings.NewReplacer(`chatroom="HERMES" id="OM.agent.agent-hermes-zeus-legacy"`, `chatroom="All Chat Rooms" id="All Chat Rooms"`,
		`uid1="OM.agent.agent-hermes-zeus-legacy" id="OM.agent.agent-hermes-zeus-legacy"`, `uid1="All Chat Rooms" id="All Chat Rooms"`,
		`to="OM.agent.agent-hermes-zeus-legacy"`, `to="All Chat Rooms"`).Replace(atakDirect)
	c, err := ParseChat([]byte(g))
	if err != nil {
		t.Fatal(err)
	}
	if c.AddressedTo("OM.agent.agent-hermes-zeus-legacy", "HERMES") {
		t.Fatal("a group message without a mention must not reach the agent")
	}
	c.Text = "@hermes status?"
	if !c.AddressedTo("OM.agent.agent-hermes-zeus-legacy", "HERMES") {
		t.Fatal("@mention in a group chat must reach the agent")
	}
}

func TestRejects(t *testing.T) {
	if _, err := ParseChat([]byte(`<!DOCTYPE x [<!ENTITY a "b">]>` + atakDirect)); err == nil {
		t.Fatal("DOCTYPE accepted")
	}
	if _, err := ParseChat([]byte(strings.Replace(atakDirect, `type="b-t-f"`, `type="a-f-G"`, 1))); err != ErrNotChat {
		t.Fatalf("non-chat: %v", err)
	}
	if _, err := ParseChat([]byte(atakDirect + strings.Repeat(" ", MaxEventBytes))); err == nil {
		t.Fatal("oversize accepted")
	}
}

func TestRenderRoundTrip(t *testing.T) {
	from := Contact{UID: "OM.agent.agent-hermes-zeus-legacy", Callsign: "HERMES", Lat: 30.1, Lon: -97.7}
	now := time.Date(2026, 10, 29, 15, 5, 0, 0, time.UTC)
	out := RenderChat(Reply{From: from, ToUID: "ANDROID-1", ToCall: "ZEUS", Text: `sim-x1 is <orbiting> & fine`, MessageID: "r1", Server: true}, now)
	c, err := ParseChat(out)
	if err != nil {
		t.Fatalf("%v\n%s", err, out)
	}
	if c.SenderUID != from.UID || c.Text != `sim-x1 is <orbiting> & fine` || !c.AddressedTo("ANDROID-1", "ZEUS") {
		t.Fatalf("round trip %+v", c)
	}
	if !strings.Contains(string(out), `<marti><dest callsign="ZEUS"/></marti>`) {
		t.Fatal("server reply without marti dest")
	}
	pli := RenderContact(Contact{UID: from.UID, Callsign: `HER"MES`, Lat: 1, Lon: 2, Remarks: "idle", Version: "0.1.0"}, now, 3*time.Minute)
	var v struct {
		UID     string `xml:"uid,attr"`
		Type    string `xml:"type,attr"`
		Contact struct {
			Callsign string `xml:"callsign,attr"`
			Endpoint string `xml:"endpoint,attr"`
		} `xml:"detail>contact"`
	}
	if err := xml.Unmarshal(pli, &v); err != nil {
		t.Fatalf("%v\n%s", err, pli)
	}
	if v.UID != from.UID || v.Type != "a-f-G-U-C" || v.Contact.Callsign != `HER"MES` || v.Contact.Endpoint == "" {
		t.Fatalf("contact %+v", v)
	}
}
