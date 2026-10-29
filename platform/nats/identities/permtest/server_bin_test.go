//go:build natsbin

package permtest

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"time"
)

const serverMode = "external binary"

func serverBin() string {
	if b := os.Getenv("NATS_SERVER_BIN"); b != "" {
		return b
	}
	return "nats-server"
}

func serverVersion() string {
	out, err := exec.Command(serverBin(), "--version").Output()
	if err != nil {
		return "unknown (" + err.Error() + ")"
	}
	return strings.TrimSpace(string(out))
}

type binServer struct {
	cmd      *exec.Cmd
	httpPort int
	log      *logBuf
	done     chan struct{}
}

// startServer runs `nats-server -c confPath` as a process. Its log (stderr) goes into lb line
// by line; readiness and the leaf count come from the monitoring port.
func startServer(confPath string, httpPort int, lb *logBuf) (natsServer, error) {
	cmd := exec.Command(serverBin(), "-c", confPath)
	pr, pw := io.Pipe()
	cmd.Stdout, cmd.Stderr = pw, pw
	if err := cmd.Start(); err != nil {
		return nil, err
	}
	s := &binServer{cmd: cmd, httpPort: httpPort, log: lb, done: make(chan struct{})}
	go func() {
		sc := bufio.NewScanner(pr)
		for sc.Scan() {
			lb.raw(sc.Text())
		}
	}()
	go func() { _ = cmd.Wait(); pw.Close(); close(s.done) }()
	deadline := time.Now().Add(10 * time.Second)
	for time.Now().Before(deadline) {
		select {
		case <-s.done:
			time.Sleep(50 * time.Millisecond) // let the scanner drain the error
			return nil, fmt.Errorf("config rejected: nats-server exited: %s", strings.Join(lb.since(0), " | "))
		default:
		}
		if r, err := http.Get(s.mon("/varz")); err == nil {
			r.Body.Close()
			if r.StatusCode == 200 {
				return s, nil
			}
		}
		time.Sleep(50 * time.Millisecond)
	}
	s.Shutdown()
	return nil, errors.New("not ready")
}

func (s *binServer) mon(path string) string {
	return fmt.Sprintf("http://127.0.0.1:%d%s", s.httpPort, path)
}

func (s *binServer) NumLeafNodes() int {
	r, err := http.Get(s.mon("/leafz"))
	if err != nil {
		return 0
	}
	defer r.Body.Close()
	var v struct {
		N int `json:"leafnodes"`
	}
	if json.NewDecoder(r.Body).Decode(&v) != nil {
		return 0
	}
	return v.N
}

// Reload sends SIGHUP (what the reloader sidecar does) and reads the server's verdict from its log.
func (s *binServer) Reload() error {
	n := s.log.len()
	if err := s.cmd.Process.Signal(syscall.SIGHUP); err != nil {
		return err
	}
	deadline := time.Now().Add(5 * time.Second)
	for time.Now().Before(deadline) {
		for _, l := range s.log.since(n) {
			if strings.Contains(l, "Failed to reload server configuration") {
				return errors.New(l)
			}
			if strings.Contains(l, "Reloaded server configuration") {
				return nil
			}
		}
		time.Sleep(25 * time.Millisecond)
	}
	return errors.New("no reload verdict in the server log within 5s")
}

func (s *binServer) Shutdown() {
	select {
	case <-s.done:
		return
	default:
	}
	_ = s.cmd.Process.Signal(syscall.SIGTERM)
	select {
	case <-s.done:
	case <-time.After(5 * time.Second):
		_ = s.cmd.Process.Kill()
		<-s.done
	}
}
