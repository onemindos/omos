// SPDX-License-Identifier: Apache-2.0

package main

// Integration tests against a disposable nats-server with JetStream (binary
// from NATS_SERVER_BIN or PATH; the tests skip without one). Stream COT is
// created as platform/nats/streams/cot.yaml defines it: tak.cot.>, 2m
// duplicate window.

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/nats-io/nats.go"
)

type natsServer struct {
	t    *testing.T
	bin  string
	port int
	dir  string
	cmd  *exec.Cmd
}

func freePort(t *testing.T) int {
	t.Helper()
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer l.Close()
	return l.Addr().(*net.TCPAddr).Port
}

func startNATS(t *testing.T) *natsServer {
	t.Helper()
	bin := os.Getenv("NATS_SERVER_BIN")
	if bin == "" {
		var err error
		if bin, err = exec.LookPath("nats-server"); err != nil {
			t.Skip("nats-server not found (set NATS_SERVER_BIN)")
		}
	}
	s := &natsServer{t: t, bin: bin, port: freePort(t), dir: t.TempDir()}
	s.start()
	t.Cleanup(s.stop)
	return s
}

func (s *natsServer) url() string { return "nats://127.0.0.1:" + strconv.Itoa(s.port) }

func (s *natsServer) start() {
	s.t.Helper()
	s.cmd = exec.Command(s.bin, "-a", "127.0.0.1", "-p", strconv.Itoa(s.port), "-js", "-sd", s.dir)
	if err := s.cmd.Start(); err != nil {
		s.t.Fatal(err)
	}
	waitFor(s.t, "nats-server listening", 10*time.Second, func() bool {
		c, err := net.DialTimeout("tcp", fmt.Sprintf("127.0.0.1:%d", s.port), 200*time.Millisecond)
		if err == nil {
			c.Close()
		}
		return err == nil
	})
}

func (s *natsServer) stop() {
	if s.cmd != nil && s.cmd.Process != nil {
		_ = s.cmd.Process.Kill()
		_, _ = s.cmd.Process.Wait()
		s.cmd = nil
	}
}

func connect(t *testing.T, url string) *nats.Conn {
	t.Helper()
	nc, err := nats.Connect(url, nats.MaxReconnects(-1), nats.ReconnectWait(100*time.Millisecond), nats.RetryOnFailedConnect(true))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(nc.Close)
	return nc
}

func addCOT(t *testing.T, nc *nats.Conn) nats.JetStreamContext {
	t.Helper()
	js, err := nc.JetStream()
	if err != nil {
		t.Fatal(err)
	}
	if _, err := js.AddStream(&nats.StreamConfig{Name: "COT", Subjects: []string{"tak.cot.>"},
		Storage: nats.FileStorage, Duplicates: 2 * time.Minute}); err != nil {
		t.Fatal(err)
	}
	return js
}

func outCount(t *testing.T, js nats.JetStreamContext) uint64 {
	t.Helper()
	si, err := js.StreamInfo("COT", &nats.StreamInfoRequest{SubjectsFilter: subjectOut})
	if err != nil {
		return 0
	}
	return si.State.Subjects[subjectOut]
}

// startBridge runs the production bridge (newBridge + start + Client.Run).
func startBridge(t *testing.T, ns *natsServer, pki *testPKI, addr string, tweak func(*config)) *bridge {
	t.Helper()
	cfg := config{natsURL: ns.url(), clientID: "omos-tak-client", takAddr: addr, protocol: "auto",
		injectEnabled: true, injectQueue: "tak-cot-inject", echoSuppress: true, echoTTL: 2 * time.Minute, outJetStream: true}
	if tweak != nil {
		tweak(&cfg)
	}
	nc := connect(t, ns.url())
	b := newBridge(cfg, nc, pki.clientTLS(t, ""))
	b.client.cfg.BackoffMin, b.client.cfg.BackoffMax = 50*time.Millisecond, 100*time.Millisecond
	if err := b.start(); err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan struct{})
	go func() { b.client.Run(ctx); close(done) }()
	t.Cleanup(func() { cancel(); <-done })
	return b
}

func TestBridgeEndToEnd(t *testing.T) {
	for _, tc := range []struct {
		name  string
		mode  peerMode
		proto string
	}{{"xml", peerXMLOnly, "0"}, {"takproto-v1", peerOfferAccept, "1"}} {
		t.Run(tc.name, func(t *testing.T) {
			ns := startNATS(t)
			js := addCOT(t, connect(t, ns.url()))
			pki := newTestPKI(t)
			peer := startPeer(t, pki, tc.mode)
			b := startBridge(t, ns, pki, peer.addr(), nil)
			pc := peer.waitConn(t)
			if tc.mode == peerOfferAccept {
				waitFor(t, "v1", 5*time.Second, func() bool { return b.client.Stats.Proto.Load() == 1 })
			}
			waitFor(t, "ready", 5*time.Second, b.ready)

			// TAK → bus: one EUD PLI lands once in COT on tak.cot.out.
			ev := pli("ANDROID-e2e", "ALPHA", nowMS())
			if err := pc.send(ev); err != nil {
				t.Fatal(err)
			}
			waitFor(t, "tak.cot.out stored", 5*time.Second, func() bool { return outCount(t, js) == 1 })
			m, err := js.GetLastMsg("COT", subjectOut)
			if err != nil {
				t.Fatal(err)
			}
			if m.Header.Get(nats.MsgIdHdr) != dedupID(ev) || m.Header.Get("Tak-Client") != "omos-tak-client" || m.Header.Get("Tak-Proto") != tc.proto {
				t.Fatalf("headers %v", m.Header)
			}
			got, err := ParseCotXML(m.Data)
			if err != nil || got.UID != "ANDROID-e2e" || !strings.Contains(string(m.Data), `callsign="ALPHA"`) {
				t.Fatalf("stored CoT %q: %v", m.Data, err)
			}

			// The same event again (as after a reconnect replay) is deduplicated.
			_ = pc.send(ev)
			time.Sleep(300 * time.Millisecond)
			if n := outCount(t, js); n != 1 {
				t.Fatalf("replayed event stored %d times", n)
			}

			// Bus → TAK: tak.cot.inject reaches the server.
			bus := connect(t, ns.url())
			inj := pli("OM.drone.sim-x1", "SIM-X1", nowMS())
			inj.Type = "a-f-A-M-H-Q"
			if err := bus.Publish(subjectInject, inj.MarshalCotXML(true)); err != nil {
				t.Fatal(err)
			}
			waitFor(t, "inject at TAK", 5*time.Second, func() bool { e, _ := peer.received("a-f-A-M-H-Q"); return len(e) == 1 })
			if e, p := peer.received("a-f-A-M-H-Q"); e[0].UID != "OM.drone.sim-x1" || strconv.Itoa(p[0]) != tc.proto {
				t.Fatalf("TAK got %s over proto %d", e[0].UID, p[0])
			}

			// TAK relays the injected event back (as it does to every other
			// client in the group): it must not reappear on tak.cot.out.
			_ = pc.send(inj)
			time.Sleep(300 * time.Millisecond)
			if n := outCount(t, js); n != 1 || b.echoDropped.Load() != 1 {
				t.Fatalf("echo: out=%d echoDropped=%d", n, b.echoDropped.Load())
			}
			// A real update from a TAK user still flows.
			_ = pc.send(pli("ANDROID-e2e", "ALPHA", nowMS()+1000))
			waitFor(t, "second PLI", 5*time.Second, func() bool { return outCount(t, js) == 2 })
		})
	}
}

// NatsBridgePlugin and tak-client share the queue group "tak-cot-inject", so
// each tak.cot.inject message is injected into TAK Server exactly once; while
// tak-client is disconnected from TAK it leaves the group and the plugin gets all.
func TestInjectQueueSharedWithPlugin(t *testing.T) {
	ns := startNATS(t)
	pki := newTestPKI(t)
	peer := startPeer(t, pki, peerXMLOnly)
	b := startBridge(t, ns, pki, peer.addr(), func(c *config) { c.outJetStream = false })
	peer.waitConn(t)
	waitFor(t, "queue sub", 5*time.Second, func() bool { b.mu.Lock(); defer b.mu.Unlock(); return b.queueSub != nil })

	plugin := connect(t, ns.url())
	var pluginGot atomic.Int64
	if _, err := plugin.QueueSubscribe(subjectInject, "tak-cot-inject", func(*nats.Msg) { pluginGot.Add(1) }); err != nil {
		t.Fatal(err)
	}
	_ = plugin.Flush()
	bus := connect(t, ns.url())
	time.Sleep(100 * time.Millisecond)
	const n = 40
	for i := 0; i < n; i++ {
		ev := pli(fmt.Sprintf("OM.track.q%02d", i), "Q", nowMS())
		ev.Type = "a-n-G"
		_ = bus.Publish(subjectInject, ev.MarshalCotXML(false))
	}
	waitFor(t, "all 40 injected once", 5*time.Second, func() bool {
		e, _ := peer.received("a-n-G")
		return int64(len(e))+pluginGot.Load() == n
	})
	time.Sleep(200 * time.Millisecond)
	e, _ := peer.received("a-n-G")
	seen := map[string]bool{}
	for _, ev := range e {
		if seen[ev.UID] {
			t.Fatalf("%s injected twice", ev.UID)
		}
		seen[ev.UID] = true
	}
	if len(e) == 0 || pluginGot.Load() == 0 || int64(len(e))+pluginGot.Load() != n {
		t.Fatalf("split tak-client=%d plugin=%d, want both > 0 and sum %d", len(e), pluginGot.Load(), n)
	}
	t.Logf("queue split: tak-client %d, plugin %d", len(e), pluginGot.Load())

	// TAK Server goes away: tak-client leaves the queue group.
	peer.close()
	waitFor(t, "queue sub dropped", 5*time.Second, func() bool { b.mu.Lock(); defer b.mu.Unlock(); return b.queueSub == nil })
	before := pluginGot.Load()
	for i := 0; i < 10; i++ {
		_ = bus.Publish(subjectInject, pli(fmt.Sprintf("OM.track.d%02d", i), "D", nowMS()).MarshalCotXML(false))
	}
	waitFor(t, "plugin takes all while TAK is down", 5*time.Second, func() bool { return pluginGot.Load()-before == 10 })
}

func TestInjectValidation(t *testing.T) {
	ns := startNATS(t)
	pki := newTestPKI(t)
	peer := startPeer(t, pki, peerXMLOnly)
	b := startBridge(t, ns, pki, peer.addr(), func(c *config) { c.outJetStream = false })
	peer.waitConn(t)
	waitFor(t, "queue sub", 5*time.Second, func() bool { b.mu.Lock(); defer b.mu.Unlock(); return b.queueSub != nil })
	bus := connect(t, ns.url())

	_ = bus.Publish(subjectInject, []byte("not xml at all"))
	_ = bus.Publish(subjectInject, protoRequest("evil", time.Now()))
	// One message that tries to carry a second event after the first.
	first := pli("OM.ok", "OK", nowMS())
	smuggled := pli("evil-second", "EVIL", nowMS())
	_ = bus.Publish(subjectInject, append(first.MarshalCotXML(true), smuggled.MarshalCotXML(true)...))
	waitFor(t, "valid part injected", 5*time.Second, func() bool { e, _ := peer.received("a-f-G-U-C"); return len(e) == 1 })
	time.Sleep(200 * time.Millisecond)
	e, _ := peer.received("")
	for _, ev := range e {
		if ev.UID == "evil-second" || ev.UID == "evil" || ev.Type == typeProtoRequest {
			t.Fatalf("TAK received %s %s", ev.UID, ev.Type)
		}
	}
	if b.injectRejected.Load() != 2 {
		t.Fatalf("rejected %d, want 2", b.injectRejected.Load())
	}
}

// NATS restarts while TAK keeps sending: acked publishes are retried and the
// Nats-Msg-Id keeps the retries from storing duplicates.
func TestNATSRestartDuringStream(t *testing.T) {
	ns := startNATS(t)
	js := addCOT(t, connect(t, ns.url()))
	pki := newTestPKI(t)
	peer := startPeer(t, pki, peerXMLOnly)
	b := startBridge(t, ns, pki, peer.addr(), nil)
	pc := peer.waitConn(t)
	waitFor(t, "ready", 5*time.Second, b.ready)

	_ = pc.send(pli("before", "B", nowMS()))
	waitFor(t, "first stored", 5*time.Second, func() bool { return outCount(t, js) == 1 })

	ns.stop()
	waitFor(t, "bridge sees NATS down", 5*time.Second, func() bool { return !b.nc.IsConnected() })
	for i := 0; i < 3; i++ {
		_ = pc.send(pli(fmt.Sprintf("during-%d", i), "D", nowMS()))
	}
	time.Sleep(300 * time.Millisecond)
	ns.start()
	js = addCOTIfMissing(t, connect(t, ns.url()))
	waitFor(t, "events sent during the outage stored once each", 20*time.Second, func() bool { return outCount(t, js) == 4 })
	time.Sleep(time.Second)
	if n := outCount(t, js); n != 4 {
		t.Fatalf("stored %d, want 4 (no duplicates)", n)
	}
}

func addCOTIfMissing(t *testing.T, nc *nats.Conn) nats.JetStreamContext {
	t.Helper()
	js, err := nc.JetStream()
	if err != nil {
		t.Fatal(err)
	}
	waitFor(t, "COT stream after restart", 10*time.Second, func() bool { _, err := js.StreamInfo("COT"); return err == nil })
	return js
}
