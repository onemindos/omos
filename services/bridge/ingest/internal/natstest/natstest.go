// Package natstest runs a disposable JetStream server inside a test, on a fixed port and
// store directory so the test can kill and restart it (NATS outage tests).
package natstest

import (
	"context"
	"fmt"
	"net"
	"testing"
	"time"

	"github.com/nats-io/nats-server/v2/server"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

type Server struct {
	t    testing.TB
	Port int
	Dir  string
	S    *server.Server
}

func freePort(t testing.TB) int {
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer l.Close()
	return l.Addr().(*net.TCPAddr).Port
}

// Start runs a JetStream server; it is shut down when the test ends.
func Start(t testing.TB) *Server {
	t.Helper()
	s := &Server{t: t, Port: freePort(t), Dir: t.TempDir()}
	s.start()
	t.Cleanup(s.Stop)
	return s
}

func (s *Server) start() {
	opts := &server.Options{Host: "127.0.0.1", Port: s.Port, JetStream: true, StoreDir: s.Dir, NoLog: true, NoSigs: true}
	srv, err := server.NewServer(opts)
	if err != nil {
		s.t.Fatal(err)
	}
	go srv.Start()
	if !srv.ReadyForConnections(10 * time.Second) {
		s.t.Fatal("nats-server not ready")
	}
	s.S = srv
}

func (s *Server) URL() string { return fmt.Sprintf("nats://127.0.0.1:%d", s.Port) }

// Stop shuts the server down (an outage). Safe to call twice.
func (s *Server) Stop() {
	if s.S != nil {
		s.S.Shutdown()
		s.S.WaitForShutdown()
		s.S = nil
	}
}

// Restart brings the server back on the same port and store.
func (s *Server) Restart() { s.Stop(); s.start() }

// Connect returns a client that reconnects forever, as omos-ingest does.
func (s *Server) Connect(t testing.TB, name string) (*nats.Conn, jetstream.JetStream) {
	t.Helper()
	nc, err := nats.Connect(s.URL(), nats.Name(name), nats.RetryOnFailedConnect(true), nats.MaxReconnects(-1),
		nats.ReconnectWait(100*time.Millisecond))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(nc.Close)
	js, err := jetstream.New(nc)
	if err != nil {
		t.Fatal(err)
	}
	return nc, js
}

// Streams creates the streams a collector publishes to, mirroring omos
// platform/nats/streams (INTEL intel.>, SYS evt.sys.>, TEL tel.feed.>, WORLD evt.world.>).
func Streams(t testing.TB, js jetstream.JetStream, dupWindow time.Duration, names ...string) {
	t.Helper()
	all := map[string][]string{"INTEL": {"intel.>"}, "SYS": {"evt.sys.>"}, "TEL": {"tel.feed.>"}, "WORLD": {"evt.world.>"}}
	if len(names) == 0 {
		names = []string{"INTEL", "SYS", "TEL", "WORLD"}
	}
	for _, n := range names {
		_, err := js.CreateOrUpdateStream(context.Background(), jetstream.StreamConfig{Name: n, Subjects: all[n],
			Storage: jetstream.FileStorage, Duplicates: dupWindow, MaxBytes: 256 << 20, Discard: jetstream.DiscardOld})
		if err != nil {
			t.Fatal(err)
		}
	}
}

// Msgs reads every message stored on a stream subject filter, in order.
func Msgs(t testing.TB, js jetstream.JetStream, stream, filter string) []jetstream.Msg {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	st, err := js.Stream(ctx, stream)
	if err != nil {
		t.Fatal(err)
	}
	info, err := st.Info(ctx)
	if err != nil {
		t.Fatal(err)
	}
	var out []jetstream.Msg
	if info.State.Msgs == 0 {
		return out
	}
	c, err := st.OrderedConsumer(ctx, jetstream.OrderedConsumerConfig{FilterSubjects: []string{filter}})
	if err != nil {
		t.Fatal(err)
	}
	for {
		b, err := c.FetchNoWait(1000)
		if err != nil {
			t.Fatal(err)
		}
		n := 0
		for m := range b.Messages() {
			out = append(out, m)
			n++
		}
		if n == 0 {
			// FetchNoWait can return early; confirm nothing is pending.
			ci, err := c.Info(ctx)
			if err != nil || ci.NumPending == 0 {
				return out
			}
		}
	}
}
