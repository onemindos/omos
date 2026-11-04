// SPDX-License-Identifier: Apache-2.0

package main

// A scripted TAK Server streaming peer for tests: mutual TLS, XML streaming,
// and the t-x-takp-v / -q / -r negotiation to TAK Protocol v1 streaming. It
// follows atak-civ takproto/README.txt, not TAK Server source.

import (
	"bufio"
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"io"
	"math/big"
	"net"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"
)

// ── Test PKI (generated per test run; no key material is committed) ─────────

type testPKI struct {
	caPEM                 []byte
	caPool                *x509.CertPool
	server                tls.Certificate
	clientCert, clientKey []byte // PEM
	dir                   string
}

func newCert(t *testing.T, tmpl, parent *x509.Certificate, parentKey *ecdsa.PrivateKey) (*x509.Certificate, *ecdsa.PrivateKey, []byte, []byte) {
	t.Helper()
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	if parent == nil {
		parent, parentKey = tmpl, key
	}
	der, err := x509.CreateCertificate(rand.Reader, tmpl, parent, &key.PublicKey, parentKey)
	if err != nil {
		t.Fatal(err)
	}
	cert, err := x509.ParseCertificate(der)
	if err != nil {
		t.Fatal(err)
	}
	kder, err := x509.MarshalECPrivateKey(key)
	if err != nil {
		t.Fatal(err)
	}
	return cert, key,
		pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: der}),
		pem.EncodeToMemory(&pem.Block{Type: "EC PRIVATE KEY", Bytes: kder})
}

func newTestPKI(t *testing.T) *testPKI {
	t.Helper()
	now := time.Now()
	caT := &x509.Certificate{SerialNumber: big.NewInt(1), Subject: pkix.Name{CommonName: "test-tak-ca"},
		NotBefore: now.Add(-time.Hour), NotAfter: now.Add(time.Hour), IsCA: true, BasicConstraintsValid: true,
		KeyUsage: x509.KeyUsageCertSign | x509.KeyUsageDigitalSignature}
	ca, caKey, caPEM, _ := newCert(t, caT, nil, nil)
	// The server certificate names "takserver", not 127.0.0.1: like a real TAK
	// Server certificate behind a Service name, so the chain-only check is exercised.
	srvT := &x509.Certificate{SerialNumber: big.NewInt(2), Subject: pkix.Name{CommonName: "takserver"},
		DNSNames: []string{"takserver"}, NotBefore: now.Add(-time.Hour), NotAfter: now.Add(time.Hour),
		KeyUsage: x509.KeyUsageDigitalSignature, ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth}}
	_, _, sPEM, sKey := newCert(t, srvT, ca, caKey)
	server, err := tls.X509KeyPair(sPEM, sKey)
	if err != nil {
		t.Fatal(err)
	}
	cliT := &x509.Certificate{SerialNumber: big.NewInt(3), Subject: pkix.Name{CommonName: "omos-tak-client"},
		NotBefore: now.Add(-time.Hour), NotAfter: now.Add(time.Hour),
		KeyUsage: x509.KeyUsageDigitalSignature, ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth}}
	_, _, cPEM, cKey := newCert(t, cliT, ca, caKey)
	pool := x509.NewCertPool()
	pool.AddCert(ca)
	p := &testPKI{caPEM: caPEM, caPool: pool, server: server, clientCert: cPEM, clientKey: cKey, dir: t.TempDir()}
	for name, b := range map[string][]byte{"ca.pem": caPEM, "client.pem": cPEM, "client.key": cKey} {
		if err := os.WriteFile(filepath.Join(p.dir, name), b, 0o600); err != nil {
			t.Fatal(err)
		}
	}
	return p
}

func (p *testPKI) file(name string) string { return filepath.Join(p.dir, name) }

// clientTLS is the tak-client TLS config built the production way.
func (p *testPKI) clientTLS(t *testing.T, serverName string) *tls.Config {
	t.Helper()
	tc, err := tlsConfig(config{certFile: p.file("client.pem"), keyFile: p.file("client.key"),
		caFile: p.file("ca.pem"), serverName: serverName})
	if err != nil {
		t.Fatal(err)
	}
	return tc
}

// ── The peer ────────────────────────────────────────────────────────────────

type peerMode int

const (
	peerXMLOnly      peerMode = iota // never offers v1
	peerOfferAccept                  // offers v1 and accepts the request
	peerOfferDecline                 // offers v1 and declines
	peerOfferSilent                  // offers v1 and never answers
)

type peerConn struct {
	conn  net.Conn
	mu    sync.Mutex
	proto int
}

func (pc *peerConn) send(ev *Event) error {
	pc.mu.Lock()
	defer pc.mu.Unlock()
	var b []byte
	if pc.proto == 1 {
		b = EncodeStream(&TakMessage{Event: ev})
	} else {
		b = streamXML(ev)
	}
	_, err := pc.conn.Write(b)
	return err
}

type takPeer struct {
	t          *testing.T
	ln         net.Listener
	mode       peerMode
	answerPing bool

	mu       sync.Mutex
	conns    []*peerConn
	got      []*Event // everything the client sent (including pings and requests)
	gotProto []int
	clientCN []string
	accepted chan *peerConn
}

func startPeer(t *testing.T, pki *testPKI, mode peerMode) *takPeer {
	t.Helper()
	ln, err := tls.Listen("tcp", "127.0.0.1:0", &tls.Config{
		Certificates: []tls.Certificate{pki.server},
		ClientAuth:   tls.RequireAndVerifyClientCert, ClientCAs: pki.caPool, MinVersion: tls.VersionTLS12,
	})
	if err != nil {
		t.Fatal(err)
	}
	p := &takPeer{t: t, ln: ln, mode: mode, answerPing: true, accepted: make(chan *peerConn, 16)}
	go p.serve()
	t.Cleanup(func() { p.close() })
	return p
}

func (p *takPeer) addr() string { return p.ln.Addr().String() }

func (p *takPeer) close() {
	p.ln.Close()
	p.mu.Lock()
	defer p.mu.Unlock()
	for _, c := range p.conns {
		c.conn.Close()
	}
}

// dropConnections closes every open client connection (the listener stays).
func (p *takPeer) dropConnections() {
	p.mu.Lock()
	defer p.mu.Unlock()
	for _, c := range p.conns {
		c.conn.Close()
	}
	p.conns = nil
}

func (p *takPeer) serve() {
	for {
		c, err := p.ln.Accept()
		if err != nil {
			return
		}
		go p.handle(c)
	}
}

func (p *takPeer) handle(c net.Conn) {
	tc := c.(*tls.Conn)
	if err := tc.Handshake(); err != nil {
		c.Close()
		return
	}
	pc := &peerConn{conn: c}
	p.mu.Lock()
	p.conns = append(p.conns, pc)
	if certs := tc.ConnectionState().PeerCertificates; len(certs) > 0 {
		p.clientCN = append(p.clientCN, certs[0].Subject.CommonName)
	}
	p.mu.Unlock()

	if p.mode != peerXMLOnly {
		offer := controlEvent("protouid", typeProtoOffer,
			`<TakControl><TakProtocolSupport version="1"/><TakServerVersionInfo serverVersion="5.8-RELEASE-69-test" apiVersion="3"/></TakControl>`,
			time.Now(), time.Minute)
		_ = pc.send(offer)
	}
	p.accepted <- pc

	br := bufio.NewReader(c)
	xr := newXMLReader(br)
	for {
		raw, err := xr.Next()
		if err != nil {
			return
		}
		ev, err := ParseCotXML(raw)
		if err != nil {
			p.t.Logf("peer: bad XML from client: %v", err)
			continue
		}
		p.record(ev, 0)
		if ev.Type == typePing && p.answerPing {
			_ = pc.send(controlEvent("takPong", typePong, "", time.Now(), 20*time.Second))
		}
		if ev.Type == typeProtoRequest {
			switch p.mode {
			case peerOfferAccept:
				pc.mu.Lock()
				_, _ = c.Write(streamXML(controlEvent("protouid", typeProtoResponse,
					`<TakControl><TakResponse status="true"/></TakControl>`, time.Now(), time.Minute)))
				pc.proto = 1
				pc.mu.Unlock()
				p.readV1(pc, bufio.NewReader(io.MultiReader(bytes.NewReader(xr.Leftover()), br)))
				return
			case peerOfferDecline:
				_ = pc.send(controlEvent("protouid", typeProtoResponse,
					`<TakControl><TakResponse status="false"/></TakControl>`, time.Now(), time.Minute))
			case peerOfferSilent:
			}
		}
	}
}

func (p *takPeer) readV1(pc *peerConn, br *bufio.Reader) {
	for {
		m, err := readStreamMessage(br)
		if err != nil {
			return
		}
		if m.Event == nil {
			continue
		}
		p.record(m.Event, 1)
		if m.Event.Type == typePing && p.answerPing {
			_ = pc.send(controlEvent("takPong", typePong, "", time.Now(), 20*time.Second))
		}
	}
}

func (p *takPeer) record(ev *Event, proto int) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.got = append(p.got, ev)
	p.gotProto = append(p.gotProto, proto)
}

// received returns client-sent events of one type (or all types when "").
func (p *takPeer) received(cotType string) ([]*Event, []int) {
	p.mu.Lock()
	defer p.mu.Unlock()
	var evs []*Event
	var protos []int
	for i, ev := range p.got {
		if cotType == "" || ev.Type == cotType {
			evs = append(evs, ev)
			protos = append(protos, p.gotProto[i])
		}
	}
	return evs, protos
}

func (p *takPeer) waitConn(t *testing.T) *peerConn {
	t.Helper()
	select {
	case pc := <-p.accepted:
		return pc
	case <-time.After(10 * time.Second):
		t.Fatal("client never connected to the peer")
		return nil
	}
}

func waitFor(t *testing.T, what string, timeout time.Duration, cond func() bool) {
	t.Helper()
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		if cond() {
			return
		}
		time.Sleep(20 * time.Millisecond)
	}
	t.Fatalf("timed out waiting for %s", what)
}

// pli is an EUD position report like ATAK sends (the takproto README example shape).
func pli(uid, callsign string, ms uint64) *Event {
	return &Event{
		Type: "a-f-G-U-C", UID: uid, How: "m-g", SendTime: ms, StartTime: ms, StaleTime: ms + 120000,
		Lat: 30.2672, Lon: -97.7431, HAE: 152.4, CE: 9.9, LE: 9999999,
		Detail: Detail{
			XML:     `<remarks>test</remarks>`,
			Contact: &Contact{Endpoint: "*:-1:stcp", Callsign: callsign},
			Group:   &Group{Name: "Cyan", Role: "Team Member"},
			Takv:    &Takv{Device: "Pixel", Platform: "ATAK-CIV", OS: "34", Version: "5.6.0"},
			Track:   &Track{Speed: 1.5, Course: 90},
		},
	}
}

func nowMS() uint64 { return uint64(time.Now().UnixMilli()) }
