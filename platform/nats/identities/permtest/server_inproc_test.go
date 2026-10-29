//go:build !natsbin

package permtest

import (
	"fmt"
	"time"

	"github.com/nats-io/nats-server/v2/server"
)

const serverMode = "in-process"

func serverVersion() string { return server.VERSION }

type inprocServer struct{ *server.Server }

func (s inprocServer) Shutdown() {
	s.Server.Shutdown()
	s.Server.WaitForShutdown()
}

// startServer runs the config in this process (no signal handlers, logs into lb).
func startServer(confPath string, _ int, lb *logBuf) (natsServer, error) {
	opts, err := server.ProcessConfigFile(confPath)
	if err != nil {
		return nil, fmt.Errorf("config rejected: %w", err)
	}
	opts.NoSigs = true
	srv, err := server.NewServer(opts)
	if err != nil {
		return nil, err
	}
	srv.SetLogger(lb, false, false)
	go srv.Start()
	if !srv.ReadyForConnections(10 * time.Second) {
		srv.Shutdown()
		return nil, fmt.Errorf("not ready")
	}
	return inprocServer{srv}, nil
}
