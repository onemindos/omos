// Package natstest starts a disposable nats-server (JetStream on, temporary store, free
// port) for integration tests. The binary comes from NATS_SERVER_BIN or PATH; without one the
// test is skipped.
package natstest

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

// Server is a running nats-server.
type Server struct {
	URL string
	cmd *exec.Cmd
	dir string
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

// Start runs a server for the test and stops it on cleanup.
func Start(t *testing.T) *Server {
	t.Helper()
	bin := os.Getenv("NATS_SERVER_BIN")
	if bin == "" {
		var err error
		if bin, err = exec.LookPath("nats-server"); err != nil {
			t.Skip("nats-server not found (set NATS_SERVER_BIN)")
		}
	}
	dir := t.TempDir()
	port := freePort(t)
	conf := filepath.Join(dir, "server.conf")
	if err := os.WriteFile(conf, []byte(fmt.Sprintf("listen: 127.0.0.1:%d\nserver_name: agentops-test\njetstream { store_dir: %q }\n", port, filepath.Join(dir, "js"))), 0o600); err != nil {
		t.Fatal(err)
	}
	cmd := exec.Command(bin, "-c", conf)
	cmd.Stdout, cmd.Stderr = nil, nil
	if err := cmd.Start(); err != nil {
		t.Fatal(err)
	}
	s := &Server{URL: fmt.Sprintf("nats://127.0.0.1:%d", port), cmd: cmd, dir: dir}
	t.Cleanup(func() {
		_ = cmd.Process.Kill()
		_, _ = cmd.Process.Wait()
	})
	deadline := time.Now().Add(10 * time.Second)
	for time.Now().Before(deadline) {
		nc, err := nats.Connect(s.URL)
		if err == nil {
			nc.Close()
			return s
		}
		time.Sleep(50 * time.Millisecond)
	}
	t.Fatal("nats-server did not start")
	return nil
}

// Connect opens a client connection closed on cleanup.
func (s *Server) Connect(t *testing.T) *nats.Conn {
	t.Helper()
	nc, err := nats.Connect(s.URL)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(nc.Close)
	return nc
}

// Streams creates the streams and buckets agentops touches, with the subjects of the omos
// CRs plus the owner-proposed bindings (ENT, `tel.feed.>` in TEL, `evt.entity.>` in EVT).
// Pass withProposed=false to leave the proposed bindings out and test the core fallback.
func Streams(t *testing.T, nc *nats.Conn, withProposed bool) jetstream.JetStream {
	t.Helper()
	js, err := jetstream.New(nc)
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	tel := []string{"tel.fleet.>"}
	evt := []string{"evt.so.>", "evt.drone.>", "evt.robot.>"}
	cfgs := []jetstream.StreamConfig{
		{Name: "MISSION", Subjects: []string{"evt.mission.>", "evt.map.>"}},
		{Name: "AGENTS", Subjects: []string{"agents.hb.>"}},
		{Name: "COT", Subjects: []string{"tak.cot.>"}},
		{Name: "CMD", Subjects: []string{"cmd.>"}, Discard: jetstream.DiscardNew},
	}
	if withProposed {
		tel = append(tel, "tel.feed.>")
		evt = append(evt, "evt.entity.>")
		cfgs = append(cfgs, jetstream.StreamConfig{Name: "ENT", Subjects: []string{"ent.*.*.state"}, MaxMsgsPerSubject: 1, AllowDirect: true})
	}
	cfgs = append(cfgs, jetstream.StreamConfig{Name: "TEL", Subjects: tel}, jetstream.StreamConfig{Name: "EVT", Subjects: evt})
	for _, c := range cfgs {
		c.Storage = jetstream.FileStorage
		c.Duplicates = 2 * time.Minute
		if _, err := js.CreateStream(ctx, c); err != nil {
			t.Fatalf("stream %s: %v", c.Name, err)
		}
	}
	if _, err := js.CreateKeyValue(ctx, jetstream.KeyValueConfig{Bucket: "geofences", History: 5}); err != nil {
		t.Fatal(err)
	}
	return js
}

// Count returns the number of messages in stream matching subject.
func Count(t *testing.T, js jetstream.JetStream, stream, subject string) uint64 {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	s, err := js.Stream(ctx, stream)
	if err != nil {
		t.Fatal(err)
	}
	info, err := s.Info(ctx, jetstream.WithSubjectFilter(subject))
	if err != nil {
		t.Fatal(err)
	}
	var n uint64
	for _, c := range info.State.Subjects {
		n += c
	}
	return n
}

// Messages returns the data of every message in stream matching subject, in order.
func Messages(t *testing.T, js jetstream.JetStream, stream, subject string) [][]byte {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	s, err := js.Stream(ctx, stream)
	if err != nil {
		t.Fatal(err)
	}
	c, err := s.OrderedConsumer(ctx, jetstream.OrderedConsumerConfig{FilterSubjects: []string{subject}})
	if err != nil {
		t.Fatal(err)
	}
	var out [][]byte
	for {
		b, err := c.FetchNoWait(100)
		if err != nil {
			t.Fatal(err)
		}
		n := 0
		for m := range b.Messages() {
			out = append(out, m.Data())
			n++
		}
		if n == 0 {
			return out
		}
	}
}

// Eventually polls cond until it holds or the timeout passes.
func Eventually(t *testing.T, timeout time.Duration, what string, cond func() bool) {
	t.Helper()
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		if cond() {
			return
		}
		time.Sleep(50 * time.Millisecond)
	}
	t.Fatalf("timed out waiting for %s", what)
}

// Evidence logs a line and, when AGENTOPS_EVIDENCE names a file, appends it there too.
func Evidence(t *testing.T, format string, args ...any) {
	t.Helper()
	line := fmt.Sprintf("[%s] %s", t.Name(), fmt.Sprintf(format, args...))
	t.Log(line)
	if p := os.Getenv("AGENTOPS_EVIDENCE"); p != "" {
		f, err := os.OpenFile(p, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
		if err == nil {
			_, _ = f.WriteString(line + "\n")
			_ = f.Close()
		}
	}
}
