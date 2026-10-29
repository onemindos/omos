// SPDX-License-Identifier: Apache-2.0

package main

import (
	"bytes"
	"encoding/hex"
	"net"
	"strings"
	"testing"
	"time"
)

// A typical ATAK SA report.
const atakSA = `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>` +
	`<event version="2.0" uid="ANDROID-589520ccfcd20f01" type="a-f-G-U-C" how="h-e" ` +
	`time="2026-10-29T15:04:05.123Z" start="2026-10-29T15:04:05.123Z" stale="2026-10-29T15:10:05.123Z">` +
	`<point lat="30.2672" lon="-97.7431" hae="150.5" ce="9.9" le="9999999.0"/>` +
	`<detail>` +
	`<takv os="34" version="5.6.0.12 (a1b2c3d4)[playstore].1700000000-CIV" device="SAMSUNG SM-G998U" platform="ATAK-CIV"/>` +
	`<contact endpoint="192.168.1.20:4242:tcp" callsign="ZEUS"/>` +
	`<uid Droid="ZEUS"/>` +
	`<precisionlocation altsrc="GPS" geopointsrc="GPS"/>` +
	`<__group role="Team Lead" name="Cyan"/>` +
	`<status battery="87"/>` +
	`<track course="271.5" speed="1.25"/>` +
	`</detail></event>`

func TestSARoundTripThroughProtobuf(t *testing.T) {
	ev, err := ParseCotXML([]byte(atakSA))
	if err != nil {
		t.Fatal(err)
	}
	d := ev.Detail
	if d.Contact == nil || d.Contact.Callsign != "ZEUS" || d.Contact.Endpoint != "192.168.1.20:4242:tcp" {
		t.Fatalf("contact not typed: %+v", d.Contact)
	}
	if d.Group == nil || d.Group.Name != "Cyan" || d.Status == nil || d.Status.Battery != 87 ||
		d.Track == nil || d.Track.Course != 271.5 || d.Takv == nil || d.PrecisionLocation == nil {
		t.Fatalf("typed detail missing: %+v", d)
	}
	if d.XML != `<uid Droid="ZEUS"/>` {
		t.Fatalf("xmlDetail = %q, want only the untyped <uid> element", d.XML)
	}

	wire := EncodeMesh(&TakMessage{Event: ev})
	if !bytes.HasPrefix(wire, []byte{0xbf, 0x01, 0xbf}) {
		t.Fatalf("mesh header = % x", wire[:3])
	}
	if len(wire) >= len(atakSA) {
		t.Errorf("protobuf %d bytes is not smaller than XML %d bytes", len(wire), len(atakSA))
	}

	back, ctl, v, err := Decode(wire)
	if err != nil || ctl != nil || v != 1 {
		t.Fatalf("decode: v=%d ctl=%v err=%v", v, ctl, err)
	}
	if back.UID != ev.UID || back.Lat != 30.2672 || back.Lon != -97.7431 || back.HAE != 150.5 ||
		back.SendTime != ev.SendTime || back.StaleTime != ev.StaleTime {
		t.Fatalf("core fields changed: %+v", back)
	}

	xmlOut := string(back.MarshalCotXML(false))
	for _, want := range []string{
		`uid="ANDROID-589520ccfcd20f01"`, `time="2026-10-29T15:04:05.123Z"`, `stale="2026-10-29T15:10:05.123Z"`,
		`<contact endpoint="192.168.1.20:4242:tcp" callsign="ZEUS"/>`, `<__group name="Cyan" role="Team Lead"/>`,
		`<status battery="87"/>`, `<track speed="1.25" course="271.5"/>`, `<uid Droid="ZEUS"/>`,
	} {
		if !strings.Contains(xmlOut, want) {
			t.Errorf("XML missing %s\n%s", want, xmlOut)
		}
	}
	if _, err := ParseCotXML([]byte(xmlOut)); err != nil {
		t.Fatalf("re-rendered XML does not parse: %v", err)
	}
}

func TestWholeElementRule(t *testing.T) {
	// contact with an extra attribute, a duplicated element and an element with
	// children must stay in xmlDetail untouched.
	raw := `<event version="2.0" uid="u1" type="a-f-G" how="m-g" time="2026-10-29T15:04:05Z" start="2026-10-29T15:04:05Z" stale="2026-10-29T15:05:05Z">` +
		`<point lat="1" lon="2" hae="0" ce="1" le="1"/><detail>` +
		`<contact callsign="A" phone="555"/>` +
		`<status battery="10"/><status battery="20"/>` +
		`<track speed="1" course="2"><extra/></track>` +
		`</detail></event>`
	ev, err := ParseCotXML([]byte(raw))
	if err != nil {
		t.Fatal(err)
	}
	if ev.Detail.Contact != nil || ev.Detail.Status != nil || ev.Detail.Track != nil {
		t.Fatalf("partial elements were typed: %+v", ev.Detail)
	}
	want := `<contact callsign="A" phone="555"/><status battery="10"/><status battery="20"/><track speed="1" course="2"><extra/></track>`
	if ev.Detail.XML != want {
		t.Fatalf("xmlDetail =\n%s\nwant\n%s", ev.Detail.XML, want)
	}
}

func TestTypedCopyIgnoredWhenXMLHasElement(t *testing.T) {
	ev := &Event{UID: "u", Type: "a-f-G", Detail: Detail{
		XML:     `<contact callsign="FROM-XML"/>`,
		Contact: &Contact{Callsign: "FROM-PROTO"},
	}}
	out := string(ev.MarshalCotXML(false))
	if strings.Contains(out, "FROM-PROTO") || !strings.Contains(out, "FROM-XML") {
		t.Fatalf("receiver rule 5 not applied: %s", out)
	}
}

func TestTakControlWire(t *testing.T) {
	wire := EncodeMesh(&TakMessage{Control: &TakControl{MinProtoVersion: 1, MaxProtoVersion: 1, ContactUID: "gw"}})
	// header, field 1 (len 8): min=1, max=1, uid="gw"
	if got, want := hex.EncodeToString(wire), "bf01bf0a08080110011a026777"; got != want {
		t.Fatalf("wire = %s, want %s", got, want)
	}
	ev, ctl, _, err := Decode(wire)
	if err != nil || ev != nil || ctl == nil || ctl.ContactUID != "gw" || ctl.MaxProtoVersion != 1 {
		t.Fatalf("decode control: ev=%v ctl=%+v err=%v", ev, ctl, err)
	}
}

func TestDecodeRejectsGarbage(t *testing.T) {
	for _, b := range [][]byte{{0xbf}, {0xbf, 0x02, 0xbf, 0x00}, {0xbf, 0x01, 0x00}, []byte("<event"), {0xbf, 0x01, 0xbf, 0x12, 0x50}} {
		if _, _, _, err := Decode(b); err == nil {
			t.Errorf("Decode(% x) succeeded", b)
		}
	}
}

func newTestMesh(forceXML bool) *Mesh {
	m := &Mesh{cfg: MeshConfig{ContactUID: "gw", ForceXML: forceXML}, contacts: map[string]*contact{}, sent: map[[32]byte]time.Time{}}
	m.version = m.computeVersion(time.Now())
	return m
}

func sa(uid string) *Event {
	return &Event{UID: uid, Type: "a-f-G-U-C", Detail: Detail{Contact: &Contact{Endpoint: "10.0.0.2:4242:tcp", Callsign: uid}}}
}

func TestNegotiation(t *testing.T) {
	now := time.Now()
	m := newTestMesh(false)
	if v := m.computeVersion(now); v != 1 {
		t.Fatalf("alone on the mesh: version %d, want 1", v)
	}

	// A legacy client heard over XML pulls the mesh down to XML (rule 3a, 4).
	m.observe(sa("legacy"), nil, 0, now)
	if v := m.computeVersion(now); v != 0 {
		t.Fatalf("with legacy client: version %d, want 0", v)
	}

	// The same client then advertises v1 (rule 3b).
	m.observe(nil, &TakControl{MinProtoVersion: 1, MaxProtoVersion: 1, ContactUID: "legacy"}, 1, now)
	if v := m.computeVersion(now); v != 1 {
		t.Fatalf("after TakControl: version %d, want 1", v)
	}

	// No TakControl for over two minutes and its last message was XML (rule 3c).
	m.observe(sa("legacy"), nil, 0, now.Add(3*time.Minute))
	if v := m.computeVersion(now.Add(3 * time.Minute)); v != 0 {
		t.Fatalf("after control timeout: version %d, want 0", v)
	}

	// Contacts we stop hearing are forgotten.
	if v := m.computeVersion(now.Add(10 * time.Minute)); v != 1 || len(m.contacts) != 0 {
		t.Fatalf("after contact timeout: version %d, %d contacts", v, len(m.contacts))
	}

	// Markers without a contact endpoint do not register a client.
	m.observe(&Event{UID: "marker-1", Type: "b-m-p-s-m"}, nil, 0, now)
	if len(m.contacts) != 0 {
		t.Fatal("marker was tracked as a contact")
	}

	if v := newTestMesh(true).computeVersion(now); v != 0 {
		t.Fatalf("MESH_PROTOCOL=xml: version %d, want 0", v)
	}
}

func TestDedupIDStableAcrossProtocols(t *testing.T) {
	ev, _ := ParseCotXML([]byte(atakSA))
	back, _, _, _ := Decode(EncodeMesh(&TakMessage{Event: ev}))
	if dedupID(ev) != dedupID(back) {
		t.Fatal("same report has different dedup IDs over XML and protobuf")
	}
}

// TestMulticastLoopback joins a real multicast group on this host, sends from
// a second socket and checks the gateway hears it, but not its own echo.
func TestMulticastLoopback(t *testing.T) {
	group := &net.UDPAddr{IP: net.IPv4(239, 2, 3, 1), Port: 16969}
	m, err := NewMesh(MeshConfig{Group: group, ContactUID: "gw-test"})
	if err != nil {
		t.Skipf("multicast unavailable here: %v", err)
	}
	defer m.Close()

	got := make(chan Received, 4)
	go m.Run(func(r Received) { got <- r }, func(err error) { t.Logf("mesh: %v", err) })

	// Our own transmission must be suppressed.
	if err := m.Send(sa("gw-self")); err != nil {
		t.Fatal(err)
	}

	peer, err := net.DialUDP("udp4", nil, group)
	if err != nil {
		t.Fatal(err)
	}
	defer peer.Close()
	ev, _ := ParseCotXML([]byte(atakSA))
	if _, err := peer.Write(EncodeMesh(&TakMessage{Event: ev})); err != nil {
		t.Fatal(err)
	}

	select {
	case r := <-got:
		if r.Event.UID == "gw-self" {
			t.Fatal("gateway published its own echo")
		}
		if r.Event.UID != ev.UID || r.Version != 1 || r.Delivery != DeliveryMesh {
			t.Fatalf("unexpected message: %+v", r)
		}
	case <-time.After(3 * time.Second):
		t.Skip("no multicast loopback delivery on this host")
	}
}

func TestDirectTCP(t *testing.T) {
	ln, err := net.Listen("tcp4", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer ln.Close()
	got := make(chan Received, 1)
	go ServeDirect(ln, nil, func(r Received) { got <- r }, func(err error) { t.Errorf("direct: %v", err) })

	c, err := net.Dial("tcp4", ln.Addr().String())
	if err != nil {
		t.Fatal(err)
	}
	c.Write([]byte(atakSA))
	c.Close()

	select {
	case r := <-got:
		if r.Event.Detail.Contact.Callsign != "ZEUS" || r.Delivery != DeliveryDirect || r.Version != 0 {
			t.Fatalf("unexpected: %+v", r)
		}
	case <-time.After(3 * time.Second):
		t.Fatal("no directed message")
	}
}
