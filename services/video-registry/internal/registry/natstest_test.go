package registry

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

// freePort returns a TCP port that is free right now.
func freePort(t *testing.T) int {
	t.Helper()
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer l.Close()
	return l.Addr().(*net.TCPAddr).Port
}

// freeUDPPort returns a UDP port that is free right now.
func freeUDPPort(t *testing.T) int {
	t.Helper()
	c, err := net.ListenPacket("udp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()
	return c.LocalAddr().(*net.UDPAddr).Port
}

func lookBin(t *testing.T, envName, name string) string {
	t.Helper()
	if p := os.Getenv(envName); p != "" {
		return p
	}
	p, err := exec.LookPath(name)
	if err != nil {
		t.Skipf("%s not found (set %s): %v", name, envName, err)
	}
	return p
}

// startNATS runs a disposable JetStream server and returns a connection to it.
func startNATS(t *testing.T) *nats.Conn {
	t.Helper()
	bin := lookBin(t, "NATS_SERVER_BIN", "nats-server")
	port := freePort(t)
	dir := t.TempDir()
	cmd := exec.Command(bin, "-a", "127.0.0.1", "-p", fmt.Sprint(port), "-js", "-sd", filepath.Join(dir, "js"))
	cmd.Stdout, cmd.Stderr = nil, nil
	if err := cmd.Start(); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = cmd.Process.Kill(); _, _ = cmd.Process.Wait() })
	var nc *nats.Conn
	var err error
	for i := 0; i < 100; i++ {
		if nc, err = nats.Connect(fmt.Sprintf("nats://127.0.0.1:%d", port)); err == nil {
			break
		}
		time.Sleep(50 * time.Millisecond)
	}
	if err != nil {
		t.Fatalf("nats-server did not start: %v", err)
	}
	t.Cleanup(nc.Close)
	return nc
}

// createBucket creates OVERWATCH_VIDEO as the NACK CR does (history 1, file storage).
func createBucket(t *testing.T, nc *nats.Conn) jetstream.KeyValue {
	t.Helper()
	js, err := jetstream.New(nc)
	if err != nil {
		t.Fatal(err)
	}
	kv, err := js.CreateKeyValue(context.Background(), jetstream.KeyValueConfig{Bucket: DefaultKV, History: 1, Storage: jetstream.FileStorage})
	if err != nil {
		t.Fatal(err)
	}
	return kv
}
