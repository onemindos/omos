// SPDX-License-Identifier: Apache-2.0

package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync/atomic"
	"testing"
	"time"
)

var testPort atomic.Int32

func init() { testPort.Store(int32(27000 + os.Getpid()%1000*10)) }

// nextGroup returns a multicast address on a port no other test uses.
func nextGroup(a, b, c, d byte) *net.UDPAddr {
	return &net.UDPAddr{IP: net.IPv4(a, b, c, d), Port: int(testPort.Add(1))}
}

// loopbackMulticast returns the loopback interface if it supports multicast
// (lo0 on macOS, lo on Linux), or skips the test.
func loopbackMulticast(t *testing.T) *net.Interface {
	t.Helper()
	ifs, err := net.Interfaces()
	if err != nil {
		t.Skipf("interfaces: %v", err)
	}
	for i := range ifs {
		f := ifs[i].Flags
		if f&net.FlagLoopback != 0 && f&net.FlagMulticast != 0 && f&net.FlagUp != 0 {
			return &ifs[i]
		}
	}
	t.Skip("no multicast-capable loopback interface on this host")
	return nil
}

// listen joins group on iface and returns datagrams received on it.
func listen(t *testing.T, iface *net.Interface, group *net.UDPAddr) <-chan []byte {
	t.Helper()
	c, err := net.ListenMulticastUDP("udp4", iface, group)
	if err != nil {
		t.Skipf("join %s on %s: %v", group, iface.Name, err)
	}
	t.Cleanup(func() { c.Close() })
	out := make(chan []byte, 16)
	go func() {
		buf := make([]byte, 65536)
		for {
			n, _, err := c.ReadFromUDP(buf)
			if err != nil {
				return
			}
			out <- append([]byte(nil), buf[:n]...)
		}
	}()
	return out
}

// peerOn returns a send function for a scripted mesh peer transmitting on iface.
func peerOn(t *testing.T, iface *net.Interface, group *net.UDPAddr) func([]byte) {
	t.Helper()
	c, err := openTransmit(iface, 1)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { c.Close() })
	return func(b []byte) {
		if _, err := c.WriteToUDP(b, group); err != nil {
			t.Fatal(err)
		}
	}
}

func recvWithin(ch <-chan []byte, d time.Duration) ([]byte, bool) {
	select {
	case b := <-ch:
		return b, true
	case <-time.After(d):
		return nil, false
	}
}

// TestTransmitOnNamedInterface is the mesh.go:84 regression: MESH_INTERFACE
// must apply to the transmit socket (IP_MULTICAST_IF), not only to receive.
// Before the fix, a peer joined on lo0 never heard the gateway, because the
// kernel sent on the default route's interface (en0).
func TestTransmitOnNamedInterface(t *testing.T) {
	lo := loopbackMulticast(t)
	group := nextGroup(239, 2, 3, 1)
	peer := listen(t, lo, group)

	m, err := NewMesh(MeshConfig{Group: group, Interface: lo, ContactUID: "gw-if"})
	if err != nil {
		t.Fatal(err)
	}
	defer m.Close()

	ti, err := m.TransmitInterface()
	if err != nil {
		t.Fatalf("read IP_MULTICAST_IF: %v", err)
	}
	if ti == nil || ti.Name != lo.Name {
		t.Fatalf("transmit interface = %v, want %s", ti, lo.Name)
	}

	ev, _ := ParseCotXML([]byte(atakSA))
	ev.UID = "omos-if-test"
	if err := m.Send(ev); err != nil {
		t.Fatal(err)
	}
	b, ok := recvWithin(peer, 3*time.Second)
	if !ok {
		t.Fatalf("peer joined on %s heard nothing from a gateway transmitting on %s", lo.Name, lo.Name)
	}
	got, _, v, err := Decode(b)
	if err != nil || got.UID != "omos-if-test" || v != 1 {
		t.Fatalf("peer decoded uid=%v v=%d err=%v", got, v, err)
	}
}

func TestTransmitDefaultInterfaceWhenUnset(t *testing.T) {
	m, err := NewMesh(MeshConfig{Group: nextGroup(239, 2, 3, 1), ContactUID: "gw-def"})
	if err != nil {
		t.Skipf("multicast unavailable: %v", err)
	}
	defer m.Close()
	if ti, err := m.TransmitInterface(); err == nil && ti != nil {
		t.Fatalf("no MESH_INTERFACE but transmit pinned to %s", ti.Name)
	}
}

func chatEvent(uid, text string) *Event {
	return &Event{
		Type: "b-t-f", UID: uid, How: "h-g-i-g-o", SendTime: 1, StartTime: 1, StaleTime: 2,
		Detail: Detail{XML: `<__chat parent="RootContactGroup" groupOwner="false" chatroom="All Chat Rooms" id="All Chat Rooms" senderCallsign="ZEUS"><chatgrp uid0="ANDROID-1" uid1="All Chat Rooms" id="All Chat Rooms"/></__chat>` +
			`<remarks source="BAO.F.ATAK.ANDROID-1" to="All Chat Rooms" time="2026-10-29T15:04:05.123Z">` + text + `</remarks>`},
	}
}

func TestGeoChatGroup(t *testing.T) {
	lo := loopbackMulticast(t)
	sa, chat := nextGroup(239, 2, 3, 1), nextGroup(224, 10, 10, 1)
	saPeer, chatPeer := listen(t, lo, sa), listen(t, lo, chat)

	m, err := NewMesh(MeshConfig{Group: sa, ChatGroup: chat, Interface: lo, ContactUID: "gw-chat"})
	if err != nil {
		t.Fatal(err)
	}
	defer m.Close()
	got := make(chan Received, 4)
	go m.Run(func(r Received) { got <- r }, func(err error) { t.Logf("mesh: %v", err) })

	// Bus -> mesh: chat goes to the chat group only.
	if err := m.Send(chatEvent("GeoChat.omos.All Chat Rooms.1", "hello from OneMind")); err != nil {
		t.Fatal(err)
	}
	b, ok := recvWithin(chatPeer, 3*time.Second)
	if !ok {
		t.Fatal("chat not sent to the chat group")
	}
	ev, _, _, err := Decode(b)
	if err != nil || !IsChat(ev) || !strings.Contains(ev.Detail.XML, "hello from OneMind") {
		t.Fatalf("chat peer got %+v err=%v", ev, err)
	}
	for {
		b, ok := recvWithin(saPeer, 300*time.Millisecond)
		if !ok {
			break
		}
		if ev, _, _, _ := Decode(b); ev != nil && IsChat(ev) {
			t.Fatal("chat also went to the SA group")
		}
	}

	// Mesh -> bus: chat heard on the chat group is delivered as "chat".
	peerOn(t, lo, chat)(chatEvent("GeoChat.ANDROID-1.All Chat Rooms.2", "hi Hermes").MarshalCotXML(true))
	select {
	case r := <-got:
		if r.Delivery != DeliveryChat || !strings.Contains(r.Event.Detail.XML, "hi Hermes") {
			t.Fatalf("unexpected: %+v", r)
		}
	case <-time.After(3 * time.Second):
		t.Fatal("chat from the mesh not delivered")
	}
}

func testCrypto(t *testing.T) *MeshCrypto {
	t.Helper()
	c, err := NewMeshCrypto(strings.Repeat("22", 32), strings.Repeat("11", 32))
	if err != nil {
		t.Fatal(err)
	}
	return c
}

// TestCryptoKnownAnswer checks the wire layout against OpenSSL, the library
// Commo uses:
//
//	openssl enc -aes-256-cbc -K 11..11 -iv 000102..0f   (PKCS#7)
//	openssl dgst -sha256 -mac HMAC -macopt hexkey:22..22 over IV||ciphertext
func TestCryptoKnownAnswer(t *testing.T) {
	plain, _ := hex.DecodeString("bf01bf0a08080110011a026777")
	iv, _ := hex.DecodeString("000102030405060708090a0b0c0d0e0f")
	want := "000102030405060708090a0b0c0d0e0f" + "beb649ed74f12888b0be176ae422cdce" +
		"f16fe7ef8a321eadddb0d4a66dcc348b2d118b98c6485ca5812cc207f5fe3e69"
	c := testCrypto(t)
	if got := hex.EncodeToString(c.encryptWithIV(plain, iv)); got != want {
		t.Fatalf("wire =\n%s\nwant\n%s", got, want)
	}
	w, _ := hex.DecodeString(want)
	back, err := c.Decrypt(w)
	if err != nil || !bytes.Equal(back, plain) {
		t.Fatalf("decrypt: %x %v", back, err)
	}
}

func TestCryptoRejects(t *testing.T) {
	c := testCrypto(t)
	msg, err := c.Encrypt([]byte(atakSA))
	if err != nil {
		t.Fatal(err)
	}
	if back, err := c.Decrypt(msg); err != nil || string(back) != atakSA {
		t.Fatalf("round trip: %v", err)
	}
	if bytes.Contains(msg, []byte("ANDROID")) {
		t.Fatal("ciphertext contains plaintext")
	}
	for _, i := range []int{0, 20, len(msg) - 1} {
		bad := append([]byte(nil), msg...)
		bad[i] ^= 1
		if _, err := c.Decrypt(bad); err == nil {
			t.Errorf("tampered byte %d accepted", i)
		}
	}
	other, _ := NewMeshCrypto(strings.Repeat("33", 32), strings.Repeat("11", 32))
	if _, err := other.Decrypt(msg); err == nil {
		t.Error("wrong auth key accepted")
	}
	for _, b := range [][]byte{nil, []byte(atakSA), msg[:40]} {
		if _, err := c.Decrypt(b); err == nil {
			t.Errorf("Decrypt(%d bytes) accepted", len(b))
		}
	}
	for _, k := range [][2]string{
		{strings.Repeat("11", 32), strings.Repeat("11", 32)}, // equal
		{strings.Repeat("11", 31), strings.Repeat("22", 32)}, // short
		{strings.Repeat("zz", 32), strings.Repeat("22", 32)}, // not hex
		{"", strings.Repeat("22", 32)},                       // one missing
	} {
		if _, err := NewMeshCrypto(k[0], k[1]); err == nil {
			t.Errorf("NewMeshCrypto(%q, %q) accepted", k[0], k[1])
		}
	}
}

func TestEncryptedMesh(t *testing.T) {
	lo := loopbackMulticast(t)
	group := nextGroup(239, 2, 3, 1)
	c := testCrypto(t)
	m, err := NewMesh(MeshConfig{Group: group, Interface: lo, ContactUID: "gw-enc", Crypto: c})
	if err != nil {
		t.Fatal(err)
	}
	defer m.Close()
	peerRx := listen(t, lo, group)
	got := make(chan Received, 4)
	errs := make(chan error, 8)
	go m.Run(func(r Received) { got <- r }, func(err error) { errs <- err })

	// Outbound is encrypted: the peer cannot read it without the keys.
	if err := m.Send(sa("omos-enc-out")); err != nil {
		t.Fatal(err)
	}
	for {
		b, ok := recvWithin(peerRx, 3*time.Second)
		if !ok {
			t.Fatal("peer heard nothing")
		}
		plain, err := c.Decrypt(b)
		if err != nil {
			t.Fatalf("peer got undecryptable datagram: %v", err)
		}
		if ev, _, _, _ := Decode(plain); ev != nil && ev.UID == "omos-enc-out" {
			if _, _, _, err := Decode(b); err == nil {
				t.Fatal("encrypted datagram decoded as plaintext")
			}
			break
		}
	}

	send := peerOn(t, lo, group)
	// A plaintext datagram is dropped with an error.
	ev, _ := ParseCotXML([]byte(atakSA))
	send(EncodeMesh(&TakMessage{Event: ev}))
	select {
	case err := <-errs:
		if !strings.Contains(err.Error(), "decrypt") {
			t.Fatalf("unexpected error: %v", err)
		}
	case r := <-got:
		t.Fatalf("plaintext accepted with encryption on: %+v", r)
	case <-time.After(3 * time.Second):
		t.Fatal("plaintext datagram neither rejected nor delivered")
	}
	// An encrypted one is accepted.
	sealed, _ := c.Encrypt(EncodeMesh(&TakMessage{Event: ev}))
	send(sealed)
	select {
	case r := <-got:
		if r.Event.UID != ev.UID || r.Version != 1 {
			t.Fatalf("unexpected: %+v", r)
		}
	case <-time.After(3 * time.Second):
		t.Fatal("encrypted datagram not delivered")
	}
}

func TestInjectAllowed(t *testing.T) {
	for uid, want := range map[string]bool{
		"ANDROID-589520ccfcd20f01":   true,
		"OM.drone.sim-x1":            false,
		"OM.vessel.mmsi-367000001":   false,
		"INTEL.acled.1":              false,
		"omos-gw-tak-mesh":           false,
		"omos-gw-other-site":         false,
		"GeoChat.ANDROID-1.All.1":    true,
		"OMEGA-not-a-projector-uid":  true,
		"intel-lowercase-is-a-phone": true,
	} {
		if got := injectAllowed(&Event{UID: uid}); got != want {
			t.Errorf("injectAllowed(%s) = %v, want %v", uid, got, want)
		}
	}
}

func TestDirectAddr(t *testing.T) {
	for in, want := range map[string]string{
		"192.168.1.20:4242:tcp": "192.168.1.20:4242",
		"10.0.0.5:4242":         "10.0.0.5:4242",
		"*:-1:stcp":             "",
		"10.0.0.5:4242:udp":     "",
		"10.0.0.5:0:tcp":        "",
		"bogus":                 "",
	} {
		got, err := directAddr(in)
		if (err != nil) != (want == "") || got != want {
			t.Errorf("directAddr(%q) = %q, %v; want %q", in, got, err, want)
		}
	}
}

func TestEndpointTrackedAndDirectEncrypted(t *testing.T) {
	c := testCrypto(t)
	m := newTestMesh(false)
	m.cfg.Crypto = c
	ln, err := net.Listen("tcp4", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer ln.Close()
	m.observe(&Event{UID: "ANDROID-2", Type: "a-f-G-U-C", Detail: Detail{Contact: &Contact{Endpoint: ln.Addr().String() + ":tcp", Callsign: "B"}}}, nil, 0, time.Now())
	ep := m.Endpoint("ANDROID-2")
	if ep != ln.Addr().String()+":tcp" {
		t.Fatalf("endpoint = %q", ep)
	}
	got := make(chan Received, 1)
	go ServeDirect(ln, m.DecodeInbound, nil, func(r Received) { got <- r }, func(err error) { t.Errorf("direct: %v", err) })
	if err := m.SendDirect(sa("omos-direct"), ep); err != nil {
		t.Fatal(err)
	}
	select {
	case r := <-got:
		if r.Event.UID != "omos-direct" || r.Delivery != DeliveryDirect {
			t.Fatalf("unexpected: %+v", r)
		}
	case <-time.After(3 * time.Second):
		t.Fatal("no directed message")
	}
}

// ── Data packages ────────────────────────────────────────────────────────────

func fileRequest(url, sum string, size int, ack bool) *Event {
	x := fmt.Sprintf(`<fileshare filename="mission.zip" senderUrl="%s" sizeInBytes="%d" sha256="%s" senderUid="ANDROID-1" senderCallsign="ZEUS" name="Mission 1"/>`, url, size, sum)
	if ack {
		x += `<ackrequest uid="ack-123" ackrequested="true" tag="Mission 1"/>`
	}
	return &Event{Type: typeFileRequest, UID: "req-1", How: "h-e", Detail: Detail{XML: x}}
}

func TestParseFileShare(t *testing.T) {
	sum := strings.Repeat("Ab", 32)
	fs, err := ParseFileShare(fileRequest("http://10.0.0.2:8080/getfile?file=abc", sum, 42, true))
	if err != nil {
		t.Fatal(err)
	}
	if fs.Size != 42 || fs.SHA256 != strings.ToLower(sum) || fs.AckUID != "ack-123" || fs.Name != "Mission 1" ||
		fs.SenderURL != "http://10.0.0.2:8080/getfile?file=abc" || fs.SenderUID != "ANDROID-1" {
		t.Fatalf("parsed %+v", fs)
	}
	for _, ev := range []*Event{
		fileRequest("http://x/getfile", "nothex", 1, false),
		fileRequest("http://x/getfile", sum, -1, false),
		{Type: typeFileRequest, Detail: Detail{XML: `<remarks/>`}},
		{Type: "a-f-G", Detail: Detail{XML: `<fileshare/>`}},
	} {
		if _, err := ParseFileShare(ev); err == nil {
			t.Errorf("accepted %+v", ev.Detail.XML)
		}
	}
}

func TestDataPackageReceive(t *testing.T) {
	payload := bytes.Repeat([]byte("PK\x03\x04 OneMind data package "), 100)
	h := sha256.Sum256(payload)
	sum := hex.EncodeToString(h[:])
	var hits atomic.Int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hits.Add(1)
		if r.URL.Path != "/getfile" || r.URL.Query().Get("file") != "abc" {
			http.NotFound(w, r)
			return
		}
		w.Write(payload)
	}))
	defer srv.Close()
	port := srv.URL[strings.LastIndex(srv.URL, ":")+1:]
	// The URL names another host; the gateway must fetch from the request's source.
	senderURL := "http://203.0.113.9:" + port + "/getfile?file=abc"
	from := "127.0.0.1:50123"

	acks := make(chan *Event, 4)
	d := &DataPackages{Dir: t.TempDir(), MaxBytes: 1 << 20, SelfUID: "omos-gw-test",
		OnErr: func(err error) { t.Logf("datapkg: %v", err) },
		SendAck: func(ack *Event, senderUID, f string) error {
			if senderUID != "ANDROID-1" || f != from {
				t.Errorf("ack routed to %s via %s", senderUID, f)
			}
			acks <- ack
			return nil
		}}
	if err := d.Init(); err != nil {
		t.Fatal(err)
	}

	// Heard on multicast, not addressed to the gateway: ignored.
	d.Handle(Received{Event: fileRequest(senderURL, sum, len(payload), true), Delivery: DeliveryMesh, From: from})
	// Directed to the gateway: fetched, verified, stored, acked.
	d.Handle(Received{Event: fileRequest(senderURL, sum, len(payload), true), Delivery: DeliveryDirect, From: from})
	ack := waitAck(t, acks)
	if ack.Type != typeFileAck || !strings.Contains(ack.Detail.XML, `uid="ack-123"`) ||
		!strings.Contains(ack.Detail.XML, `success="true"`) || !strings.Contains(ack.Detail.XML, `senderUid="omos-gw-test"`) ||
		!strings.Contains(ack.Detail.XML, `sha256="`+sum+`"`) {
		t.Fatalf("ack = %s", ack.Detail.XML)
	}
	stored, err := os.ReadFile(filepath.Join(d.Dir, "inbox", sum+"-mission.zip"))
	if err != nil || !bytes.Equal(stored, payload) {
		t.Fatalf("stored file: %v", err)
	}
	if hits.Load() != 1 || d.Received.Load() != 1 {
		t.Fatalf("hits=%d received=%d", hits.Load(), d.Received.Load())
	}
	// The ack renders as valid CoT.
	if _, err := ParseCotXML(ack.MarshalCotXML(true)); err != nil {
		t.Fatalf("ack XML: %v", err)
	}

	// Wrong hash: nothing stored, nack sent.
	bad := strings.Repeat("0", 64)
	d.Handle(Received{Event: fileRequest(senderURL, bad, len(payload), true), Delivery: DeliveryDirect, From: from})
	if ack := waitAck(t, acks); !strings.Contains(ack.Detail.XML, `success="false"`) {
		t.Fatalf("nack = %s", ack.Detail.XML)
	}
	if _, err := os.Stat(filepath.Join(d.Dir, "inbox", bad+"-mission.zip")); err == nil {
		t.Fatal("file with wrong hash stored")
	}

	// Over the size limit: refused before any request.
	d.MaxBytes = 10
	before := hits.Load()
	if _, err := d.Fetch(&FileShare{SenderURL: senderURL, SHA256: sum, Size: int64(len(payload))}, from); err == nil {
		t.Fatal("oversize accepted")
	}
	if hits.Load() != before {
		t.Fatal("oversize request was fetched")
	}
	entries, _ := os.ReadDir(filepath.Join(d.Dir, "inbox"))
	for _, e := range entries {
		if strings.HasPrefix(e.Name(), ".partial") {
			t.Fatalf("partial file left behind: %s", e.Name())
		}
	}
}

func waitAck(t *testing.T, ch <-chan *Event) *Event {
	t.Helper()
	select {
	case a := <-ch:
		return a
	case <-time.After(5 * time.Second):
		t.Fatal("no ack")
		return nil
	}
}

func TestFetchURL(t *testing.T) {
	for _, c := range []struct{ in, from, want string }{
		{"http://10.0.0.9:8080/getfile?file=x", "192.168.1.20:5000", "http://192.168.1.20:8080/getfile?file=x"},
		{"http://user:pw@evil.example/getfile?file=x", "192.168.1.20:5000", "http://192.168.1.20:80/getfile?file=x"},
		{"https://10.0.0.9:8443/Marti/sync/content?hash=x", "192.168.1.20:5000", ""},
		{"file:///etc/passwd", "192.168.1.20:5000", ""},
	} {
		got, err := fetchURL(c.in, c.from)
		if (err != nil) != (c.want == "") || got != c.want {
			t.Errorf("fetchURL(%q) = %q, %v; want %q", c.in, got, err, c.want)
		}
	}
	if s := safeFilename(`..\..\evil/../.bashrc`); strings.ContainsAny(s, `/\`) || strings.HasPrefix(s, ".") {
		t.Errorf("safeFilename = %q", s)
	}
}

func TestGetfileServe(t *testing.T) {
	d := &DataPackages{Dir: t.TempDir()}
	if err := d.Init(); err != nil {
		t.Fatal(err)
	}
	os.WriteFile(filepath.Join(d.Dir, "outbox", "pkg.zip"), []byte("zipdata"), 0o640)
	os.WriteFile(filepath.Join(d.Dir, "secret.txt"), []byte("nope"), 0o640)
	os.WriteFile(filepath.Join(d.Dir, "outbox", ".hidden"), []byte("nope"), 0o640)
	os.Mkdir(filepath.Join(d.Dir, "outbox", "sub"), 0o750)
	srv := httptest.NewServer(d)
	defer srv.Close()

	get := func(path string) (int, string) {
		resp, err := http.Get(srv.URL + path)
		if err != nil {
			t.Fatal(err)
		}
		defer resp.Body.Close()
		b := new(bytes.Buffer)
		b.ReadFrom(resp.Body)
		return resp.StatusCode, b.String()
	}
	if code, body := get("/getfile?file=pkg.zip"); code != 200 || body != "zipdata" {
		t.Fatalf("pkg.zip: %d %q", code, body)
	}
	for _, p := range []string{"/getfile?file=" + url.QueryEscape("../secret.txt"), "/getfile?file=.hidden",
		"/getfile?file=sub", "/getfile?file=", "/getfile?file=missing.zip", "/other?file=pkg.zip",
		"/getfile?file=" + url.QueryEscape(`..\secret.txt`)} {
		if code, _ := get(p); code != 404 {
			t.Errorf("%s: %d, want 404", p, code)
		}
	}
}
