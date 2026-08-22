// Package bus wraps the NATS connection lifecycle: connect, subscribe, publish,
// heartbeat, and graceful drain/close.
//
// Connection policy (NATS-DOCTRINE §6):
//   - nats.go core client only; no orbit extensions yet.
//   - Reconnect handlers log connect/disconnect/reconnect events.
//   - Credentials: no-auth user (GLOBAL/onemind account) by default.
//     When --creds is supplied, a NATS .creds file is loaded for an
//     authenticated account. The switch to FABRIC-account creds is
//     documented in TODO items here and in cmd/omd/main.go.
package bus

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/nats-io/nats.go"
)

// Connect establishes a NATS connection with reconnect handlers.
// credsFile may be empty (no-auth path). nodeID is embedded in log messages.
func Connect(url, nodeID, credsFile string) (*nats.Conn, error) {
	opts := []nats.Option{
		nats.Name(fmt.Sprintf("omd.%s", nodeID)),
		nats.MaxReconnects(-1), // infinite reconnect attempts
		nats.ReconnectWait(2 * time.Second),
		nats.DisconnectErrHandler(func(nc *nats.Conn, err error) {
			if err != nil {
				slog.Warn("NATS disconnected", "node", nodeID, "err", err)
			} else {
				slog.Info("NATS disconnected cleanly", "node", nodeID)
			}
		}),
		nats.ReconnectHandler(func(nc *nats.Conn) {
			slog.Info("NATS reconnected", "node", nodeID, "url", nc.ConnectedUrl())
		}),
		nats.ClosedHandler(func(nc *nats.Conn) {
			slog.Info("NATS connection closed", "node", nodeID)
		}),
		nats.ErrorHandler(func(nc *nats.Conn, sub *nats.Subscription, err error) {
			slog.Error("NATS async error", "node", nodeID, "sub", sub.Subject, "err", err)
		}),
	}

	// TODO: switch to FABRIC-account creds when P1-A account migration is complete.
	// Until then omd connects via the no_auth_user path (GLOBAL/onemind account).
	if credsFile != "" {
		opts = append(opts, nats.UserCredentials(credsFile))
		slog.Info("NATS using creds file", "file", credsFile)
	}

	nc, err := nats.Connect(url, opts...)
	if err != nil {
		return nil, fmt.Errorf("nats connect %s: %w", url, err)
	}
	slog.Info("NATS connected", "node", nodeID, "url", nc.ConnectedUrl(),
		"server", nc.ConnectedServerId())
	return nc, nil
}
