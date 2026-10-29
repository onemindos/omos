// SPDX-License-Identifier: Apache-2.0

package main

import (
	"bufio"
	"bytes"
	"context"
	"crypto/tls"
	"errors"
	"strings"
	"sync"
	"testing"
	"time"
)

// ── Unit tests ──────────────────────────────────────────────────────────────

func TestStreamFramingRoundTrip(t *testing.T) {
	ev := pli("ANDROID-1", "ALPHA", 1793275200000)
	frame := EncodeStream(&TakMessage{Event: ev})
	if frame[0] != 0xbf {
		t.Fatalf("frame starts 0x%02x", frame[0])
	}
	// Two frames back to back must come out as two messages.
	m1, err := readStreamMessage(bufio.NewReader(bytes.NewReader(append(frame, frame...))))
	if err != nil {
		t.Fatal(err)
	}
	if m1.Event.UID != "ANDROID-1" || m1.Event.Detail.Contact.Callsign != "ALPHA" || m1.Event.Detail.Takv.Platform != "ATAK-CIV" {
		t.Fatalf("decoded %+v", m1.Event)
	}
}

func TestStreamFramingRejectsBadInput(t *testing.T) {
	if _, err := readStreamMessage(bufio.NewReader(bytes.NewReader([]byte{'<', 'e'}))); !errors.Is(err, errBadFrame) {
		t.Fatalf("XML byte as v1 frame: %v", err)
	}
	huge := []byte{0xbf, 0x80, 0x80, 0x80, 0x10} // length 32 MiB
	if _, err := readStreamMessage(bufio.NewReader(bytes.NewReader(huge))); !errors.Is(err, errEventTooLarge) {
		t.Fatalf("oversize frame: %v", err)
	}
	long := append([]byte{0xbf}, bytes.Repeat([]byte{0xff}, 11)...)
	if _, err := readStreamMessage(bufio.NewReader(bytes.NewReader(long))); !errors.Is(err, errBadFrame) {
		t.Fatalf("11-byte varint: %v", err)
	}
}

// slowReader returns at most n bytes per Read, to split events across reads.
type slowReader struct {
	b []byte
	n int
}

func (s *slowReader) Read(p []byte) (int, error) {
	if len(s.b) == 0 {
		return 0, errors.New("EOF")
	}
	k := min(s.n, len(p), len(s.b))
	copy(p, s.b[:k])
	s.b = s.b[k:]
	return k, nil
}

func TestXMLReaderSplitsEvents(t *testing.T) {
	a := streamXML(pli("A", "A", 1))
	b := streamXML(pli("B", "B", 2))
	stream := append(append(append([]byte("\n  "), a...), b...), []byte(`<?xml version="1.0"?>`)...)
	xr := newXMLReader(bufio.NewReaderSize(&slowReader{b: stream, n: 7}, 16))
	for _, want := range []string{"A", "B"} {
		raw, err := xr.Next()
		if err != nil {
			t.Fatal(err)
		}
		if !bytes.HasPrefix(raw, []byte("<event")) {
			t.Fatalf("event does not start at <event: %q", raw[:20])
		}
		ev, err := ParseCotXML(raw)
		if err != nil || ev.UID != want {
			t.Fatalf("got %v %v, want uid %s", ev, err, want)
		}
	}
	// Leftover is whatever was already read past the last event: a prefix of
	// the bytes that follow it (after a v1 agreement, the first v1 frame).
	if !strings.HasPrefix(`<?xml version="1.0"?>`, string(xr.Leftover())) {
		t.Fatalf("leftover %q", xr.Leftover())
	}
}

func TestXMLReaderBoundsBuffer(t *testing.T) {
	junk := bytes.Repeat([]byte("x"), maxEventBytes+10000)
	xr := newXMLReader(bufio.NewReader(bytes.NewReader(junk)))
	if _, err := xr.Next(); !errors.Is(err, errEventTooLarge) {
		t.Fatalf("got %v", err)
	}
}

func TestNegotiationParsing(t *testing.T) {
	offer, err := ParseCotXML(streamXML(controlEvent("protouid", typeProtoOffer,
		`<TakControl><TakProtocolSupport version="1"/><TakServerVersionInfo serverVersion="5.8-RELEASE-69" apiVersion="3"/></TakControl>`,
		time.Now(), time.Minute))[len(xmlDecl):])
	if err != nil {
		t.Fatal(err)
	}
	ok, sv := offersV1(offer)
	if !ok || sv != "5.8-RELEASE-69" {
		t.Fatalf("offersV1 = %v %q", ok, sv)
	}
	v2only, _ := ParseCotXML(controlEvent("p", typeProtoOffer, `<TakControl><TakProtocolSupport version="2"/></TakControl>`, time.Now(), time.Minute).MarshalCotXML(false))
	if ok, _ := offersV1(v2only); ok {
		t.Fatal("a v2-only offer was taken as v1")
	}
	yes, _ := ParseCotXML(controlEvent("p", typeProtoResponse, `<TakControl><TakResponse status="true"/></TakControl>`, time.Now(), time.Minute).MarshalCotXML(false))
	no, _ := ParseCotXML(controlEvent("p", typeProtoResponse, `<TakControl><TakResponse status="false"/></TakControl>`, time.Now(), time.Minute).MarshalCotXML(false))
	if !responseAccepted(yes) || responseAccepted(no) {
		t.Fatal("responseAccepted wrong")
	}
	req, err := ParseCotXML(protoRequest("me", time.Now())[len(xmlDecl):])
	if err != nil || req.Type != typeProtoRequest || !strings.Contains(req.Detail.XML, `<TakRequest version="1"/>`) {
		t.Fatalf("request %+v %v", req, err)
	}
	if !bytes.HasPrefix(protoRequest("me", time.Now()), []byte("<?xml")) || !bytes.Contains(protoRequest("me", time.Now()), []byte("?>\n<event")) {
		t.Fatal("v0 stream events must start with an XML declaration and a newline")
	}
}

func TestControlTypes(t *testing.T) {
	for _, ty := range []string{"t-x-takp-v", "t-x-takp-q", "t-x-takp-r", "t-x-c-t", "t-x-c-t-r"} {
		if !isControlType(ty) {
			t.Errorf("%s should be control", ty)
		}
	}
	for _, ty := range []string{"a-f-G-U-C", "b-t-f", "t-x-m-c", "t-x-d-d"} {
		if isControlType(ty) {
			t.Errorf("%s should not be control", ty)
		}
	}
}

func TestEchoFilter(t *testing.T) {
	f := newEchoFilter(time.Minute, 3)
	now := time.Now()
	a := pli("OM.aircraft.abc", "A", 1000)
	f.Remember(a, now)
	if !f.IsEcho(a, now.Add(30*time.Second)) {
		t.Fatal("echo within TTL not detected")
	}
	if f.IsEcho(a, now.Add(2*time.Minute)) {
		t.Fatal("echo after TTL still detected")
	}
	b := pli("OM.aircraft.abc", "A", 2000) // same uid, newer time: a real update
	if f.IsEcho(b, now) {
		t.Fatal("a newer event with the same uid was taken as an echo")
	}
	for i := 0; i < 10; i++ {
		f.Remember(pli("x", "x", uint64(i)), now)
	}
	if f.Len() > 3 {
		t.Fatalf("filter grew to %d past its bound", f.Len())
	}
}

func TestDedupIDStable(t *testing.T) {
	a, b := pli("ANDROID-1", "A", 5), pli("ANDROID-1", "A", 5)
	b.Detail.XML = "<remarks>different detail, same observation</remarks>"
	if dedupID(a) != dedupID(b) || !strings.HasPrefix(dedupID(a), "takc-") {
		t.Fatalf("dedupID %s %s", dedupID(a), dedupID(b))
	}
	if dedupID(a) == dedupID(pli("ANDROID-1", "A", 6)) {
		t.Fatal("different times gave the same dedup id")
	}
}

func TestLoadConfigRequiresCertAndCA(t *testing.T) {
	t.Setenv("TAK_CERT_FILE", "")
	if _, err := loadConfig(); err == nil {
		t.Fatal("missing certificate accepted")
	}
	t.Setenv("TAK_CERT_FILE", "c.pem")
	t.Setenv("TAK_KEY_FILE", "c.key")
	if _, err := loadConfig(); err == nil || !strings.Contains(err.Error(), "TAK_CA_FILE") {
		t.Fatalf("missing CA accepted: %v", err)
	}
	t.Setenv("TAK_CA_FILE", "ca.pem")
	t.Setenv("TAK_PROTOCOL", "v2")
	if _, err := loadConfig(); err == nil {
		t.Fatal("bad TAK_PROTOCOL accepted")
	}
	t.Setenv("TAK_PROTOCOL", "")
	c, err := loadConfig()
	if err != nil || !c.injectEnabled || !c.echoSuppress || !c.outJetStream || c.protocol != "auto" || c.injectQueue != "tak-cot-inject" {
		t.Fatalf("defaults %+v %v", c, err)
	}
}

// ── Client against the scripted peer ────────────────────────────────────────

type sink struct {
	mu     sync.Mutex
	evs    []*Event
	raw    [][]byte
	protos []int
}

func (s *sink) add(ev *Event, raw []byte, proto int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.evs = append(s.evs, ev)
	s.raw = append(s.raw, append([]byte(nil), raw...))
	s.protos = append(s.protos, proto)
}

func (s *sink) len() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return len(s.evs)
}

func runClient(t *testing.T, cfg ClientConfig) (*Client, *sink) {
	t.Helper()
	c := NewClient(cfg)
	s := &sink{}
	c.OnEvent = s.add
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan struct{})
	go func() { c.Run(ctx); close(done) }()
	t.Cleanup(func() { cancel(); <-done })
	return c, s
}

func TestClientXMLStream(t *testing.T) {
	pki := newTestPKI(t)
	peer := startPeer(t, pki, peerXMLOnly)
	c, s := runClient(t, ClientConfig{Addr: peer.addr(), TLS: pki.clientTLS(t, ""), UID: "omos-tak-client"})
	pc := peer.waitConn(t)

	if err := pc.send(pli("ANDROID-1", "ALPHA", nowMS())); err != nil {
		t.Fatal(err)
	}
	// Plumbing from the server is never handed on.
	_ = pc.send(controlEvent("takPong", typePong, "", time.Now(), time.Minute))
	if err := pc.send(pli("ANDROID-2", "BRAVO", nowMS())); err != nil {
		t.Fatal(err)
	}
	waitFor(t, "two PLIs", 5*time.Second, func() bool { return s.len() == 2 })
	if s.evs[0].UID != "ANDROID-1" || s.protos[0] != 0 || !bytes.HasPrefix(s.raw[0], []byte("<event")) {
		t.Fatalf("first event %+v proto %d", s.evs[0], s.protos[0])
	}
	if c.Stats.Proto.Load() != 0 {
		t.Fatal("switched protocol without an offer")
	}

	// Inject goes out as one XML event with a declaration.
	c.Inject(pli("OM.drone.sim-x1", "SIM-X1", nowMS()))
	waitFor(t, "inject at peer", 5*time.Second, func() bool { evs, _ := peer.received("a-f-G-U-C"); return len(evs) == 1 })
	evs, protos := peer.received("a-f-G-U-C")
	if evs[0].UID != "OM.drone.sim-x1" || protos[0] != 0 {
		t.Fatalf("peer got %+v over proto %d", evs[0], protos[0])
	}
	if peer.clientCN[0] != "omos-tak-client" {
		t.Fatalf("client certificate CN %q", peer.clientCN[0])
	}
}

func TestClientNegotiatesV1(t *testing.T) {
	pki := newTestPKI(t)
	peer := startPeer(t, pki, peerOfferAccept)
	c, s := runClient(t, ClientConfig{Addr: peer.addr(), TLS: pki.clientTLS(t, ""), UID: "omos-tak-client"})
	pc := peer.waitConn(t)

	waitFor(t, "v1 agreed", 5*time.Second, func() bool {
		pc.mu.Lock()
		defer pc.mu.Unlock()
		return pc.proto == 1 && c.Stats.Proto.Load() == 1
	})
	if sv, _ := c.Stats.ServerVersion.Load().(string); sv != "5.8-RELEASE-69-test" {
		t.Fatalf("server version %q", sv)
	}
	reqs, _ := peer.received(typeProtoRequest)
	if len(reqs) != 1 {
		t.Fatalf("%d t-x-takp-q sent, want 1", len(reqs))
	}

	if err := pc.send(pli("ANDROID-1", "ALPHA", nowMS())); err != nil {
		t.Fatal(err)
	}
	waitFor(t, "v1 PLI", 5*time.Second, func() bool { return s.len() == 1 })
	ev := s.evs[0]
	if s.protos[0] != 1 || ev.Detail.Contact.Callsign != "ALPHA" || ev.Detail.Group.Name != "Cyan" {
		t.Fatalf("v1 event %+v proto %d", ev, s.protos[0])
	}
	// The XML handed on is complete CoT that parses back to the same event.
	back, err := ParseCotXML(s.raw[0])
	if err != nil || back.UID != "ANDROID-1" || back.Detail.Takv == nil {
		t.Fatalf("published XML %q: %v", s.raw[0], err)
	}

	c.Inject(pli("OM.drone.sim-x1", "SIM-X1", nowMS()))
	waitFor(t, "v1 inject at peer", 5*time.Second, func() bool { evs, _ := peer.received("a-f-G-U-C"); return len(evs) == 1 })
	_, protos := peer.received("a-f-G-U-C")
	if protos[0] != 1 {
		t.Fatal("inject after agreement was not sent as v1")
	}
}

func TestClientStaysXMLWhenDeclinedOrDisabled(t *testing.T) {
	pki := newTestPKI(t)
	t.Run("declined", func(t *testing.T) {
		peer := startPeer(t, pki, peerOfferDecline)
		c, s := runClient(t, ClientConfig{Addr: peer.addr(), TLS: pki.clientTLS(t, ""), UID: "u"})
		pc := peer.waitConn(t)
		waitFor(t, "request", 5*time.Second, func() bool { r, _ := peer.received(typeProtoRequest); return len(r) == 1 })
		time.Sleep(100 * time.Millisecond)
		_ = pc.send(pli("A", "A", nowMS()))
		waitFor(t, "event", 5*time.Second, func() bool { return s.len() == 1 })
		c.Inject(pli("OM.x", "X", nowMS()))
		waitFor(t, "inject", 5*time.Second, func() bool { e, _ := peer.received("a-f-G-U-C"); return len(e) == 1 })
		if _, p := peer.received("a-f-G-U-C"); p[0] != 0 || c.Stats.Proto.Load() != 0 {
			t.Fatal("not XML after a decline")
		}
	})
	t.Run("TAK_PROTOCOL=xml", func(t *testing.T) {
		peer := startPeer(t, pki, peerOfferAccept)
		_, s := runClient(t, ClientConfig{Addr: peer.addr(), TLS: pki.clientTLS(t, ""), UID: "u", Protocol: "xml"})
		pc := peer.waitConn(t)
		_ = pc.send(pli("A", "A", nowMS()))
		waitFor(t, "event", 5*time.Second, func() bool { return s.len() == 1 })
		if r, _ := peer.received(typeProtoRequest); len(r) != 0 {
			t.Fatal("requested v1 although TAK_PROTOCOL=xml")
		}
	})
}

func TestClientNegotiationTimeoutFallsBackToXML(t *testing.T) {
	pki := newTestPKI(t)
	peer := startPeer(t, pki, peerOfferSilent)
	c, _ := runClient(t, ClientConfig{Addr: peer.addr(), TLS: pki.clientTLS(t, ""), UID: "u",
		NegotiateTimeout: 300 * time.Millisecond, BackoffMin: 50 * time.Millisecond, BackoffMax: 100 * time.Millisecond})
	peer.waitConn(t)
	peer.waitConn(t) // reconnected after the timeout
	time.Sleep(300 * time.Millisecond)
	if r, _ := peer.received(typeProtoRequest); len(r) != 1 {
		t.Fatalf("%d requests; the second session must stay on XML", len(r))
	}
	if !c.forceXML.Load() || c.Stats.Connects.Load() < 2 {
		t.Fatal("no sticky XML fallback after the timeout")
	}
	// Injects still flow on the XML session.
	c.Inject(pli("OM.x", "X", nowMS()))
	waitFor(t, "inject after fallback", 5*time.Second, func() bool { e, _ := peer.received("a-f-G-U-C"); return len(e) == 1 })
}

func TestClientKeepAliveAndReset(t *testing.T) {
	pki := newTestPKI(t)
	t.Run("ping answered keeps the stream", func(t *testing.T) {
		peer := startPeer(t, pki, peerXMLOnly)
		c, _ := runClient(t, ClientConfig{Addr: peer.addr(), TLS: pki.clientTLS(t, ""), UID: "omos-tak-client",
			IdleBeforePing: 200 * time.Millisecond, PingEvery: 100 * time.Millisecond, ResetAfter: 900 * time.Millisecond})
		peer.waitConn(t)
		waitFor(t, "pings", 5*time.Second, func() bool { p, _ := peer.received(typePing); return len(p) >= 3 })
		pings, _ := peer.received(typePing)
		if pings[0].UID != "omos-tak-client-ping" {
			t.Fatalf("ping uid %q", pings[0].UID)
		}
		time.Sleep(1200 * time.Millisecond)
		if n := c.Stats.Connects.Load(); n != 1 {
			t.Fatalf("answered pings must keep the connection; connects=%d", n)
		}
	})
	t.Run("silent server is reset", func(t *testing.T) {
		peer := startPeer(t, pki, peerXMLOnly)
		peer.answerPing = false
		c, _ := runClient(t, ClientConfig{Addr: peer.addr(), TLS: pki.clientTLS(t, ""), UID: "u",
			IdleBeforePing: 100 * time.Millisecond, PingEvery: 100 * time.Millisecond, ResetAfter: 400 * time.Millisecond,
			BackoffMin: 50 * time.Millisecond, BackoffMax: 100 * time.Millisecond})
		peer.waitConn(t)
		peer.waitConn(t)
		if c.Stats.Connects.Load() < 2 {
			t.Fatal("no reconnect after silence")
		}
	})
}

func TestClientReconnectsAfterServerDrop(t *testing.T) {
	pki := newTestPKI(t)
	peer := startPeer(t, pki, peerXMLOnly)
	c, s := runClient(t, ClientConfig{Addr: peer.addr(), TLS: pki.clientTLS(t, ""), UID: "u",
		BackoffMin: 50 * time.Millisecond, BackoffMax: 100 * time.Millisecond})
	peer.waitConn(t)
	peer.dropConnections()
	pc := peer.waitConn(t)
	_ = pc.send(pli("A", "A", nowMS()))
	waitFor(t, "event after reconnect", 5*time.Second, func() bool { return s.len() == 1 })
	if c.Stats.Connects.Load() != 2 || c.Stats.Disconnects.Load() != 1 {
		t.Fatalf("connects=%d disconnects=%d", c.Stats.Connects.Load(), c.Stats.Disconnects.Load())
	}
}

func TestClientTLSVerification(t *testing.T) {
	pki := newTestPKI(t)
	peer := startPeer(t, pki, peerXMLOnly)
	other := newTestPKI(t) // a different CA

	dial := func(tc *tls.Config) error {
		c := NewClient(ClientConfig{Addr: peer.addr(), TLS: tc, UID: "u"})
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		return c.session(ctx)
	}
	t.Run("server from an unknown CA is refused", func(t *testing.T) {
		tc, err := tlsConfig(config{certFile: pki.file("client.pem"), keyFile: pki.file("client.key"), caFile: other.file("ca.pem")})
		if err != nil {
			t.Fatal(err)
		}
		if err := dial(tc); err == nil || !strings.Contains(err.Error(), "TLS handshake") {
			t.Fatalf("got %v", err)
		}
	})
	t.Run("hostname is checked when TAK_SERVER_NAME is set", func(t *testing.T) {
		if err := dial(pki.clientTLS(t, "wrong-name")); err == nil {
			t.Fatal("wrong server name accepted")
		}
		c := NewClient(ClientConfig{Addr: peer.addr(), TLS: pki.clientTLS(t, "takserver"), UID: "u"})
		s := &sink{}
		c.OnEvent = s.add
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		go c.Run(ctx)
		pc := peer.waitConn(t)
		_ = pc.send(pli("A", "A", nowMS()))
		waitFor(t, "event with hostname check", 5*time.Second, func() bool { return s.len() == 1 })
	})
	t.Run("client without a TAK certificate is refused by the server", func(t *testing.T) {
		tc := pki.clientTLS(t, "")
		tc.Certificates = nil
		err := dial(tc)
		// TLS 1.3 reports the server's refusal on the first read, so the
		// session either fails the handshake or ends at once.
		if err == nil {
			t.Fatal("session without a client certificate succeeded")
		}
	})
}

func TestInjectQueueBoundedAndAged(t *testing.T) {
	c := NewClient(ClientConfig{Addr: "127.0.0.1:1", InjectQueue: 3, InjectMaxAge: 50 * time.Millisecond})
	for i := 0; i < 5; i++ {
		c.Inject(pli("x", "x", uint64(i)))
	}
	if c.Stats.InjectDropped.Load() != 2 || len(c.inject) != 3 {
		t.Fatalf("dropped=%d queued=%d", c.Stats.InjectDropped.Load(), len(c.inject))
	}
	first := <-c.inject
	if first.ev.SendTime != 2 {
		t.Fatalf("oldest kept is %d; the two oldest should have been dropped", first.ev.SendTime)
	}
}
