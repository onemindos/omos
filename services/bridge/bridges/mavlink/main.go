// Command bridge-mavlink connects one simulated MAVLink vehicle to the OneMind bus.
//
// Telemetry: tel.fleet.drone.{id}.position and .status (doctrine §9.2 geo-envelope,
// HAE converted from MSL). Commands: cmd.fleet.drone.{id}.{goto|rtl|land|arm|disarm|estop|takeoff},
// accepted only when signed by the policy service, approved by humans, unexpired, unseen,
// simulated and inside the geofence. Outcomes: evt.drone.{id}.command (and .task for tasks).
// Contracts: onemind-xr docs/contracts/asset-command.md, robotics-cot-profile.md.
package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bluenviron/gomavlib/v3"
	"github.com/bluenviron/gomavlib/v3/pkg/dialects/ardupilotmega"
	"github.com/nats-io/nats.go"
)

type natsBus struct{ nc *nats.Conn }

func (n natsBus) Publish(subject, msgID string, data []byte) error {
	m := nats.NewMsg(subject)
	m.Data = data
	if msgID != "" {
		m.Header.Set(nats.MsgIdHdr, msgID)
	}
	return n.nc.PublishMsg(m)
}

func main() {
	cfg, err := loadConfig()
	if err != nil {
		log.Fatalf("[bridge-mavlink] config: %v", err)
	}
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()
	if err := run(ctx, cfg, nil); err != nil {
		log.Fatalf("[%s] %v", cfg.Source, err)
	}
}

// run starts the bridge and blocks until ctx ends. ready, if set, is called once the
// command subscription is active.
func run(ctx context.Context, cfg Config, ready func()) error {
	opts := []nats.Option{
		nats.Name(cfg.Source + "/" + cfg.AssetID),
		nats.MaxReconnects(-1),
		nats.ReconnectWait(2 * time.Second),
		nats.DisconnectErrHandler(func(_ *nats.Conn, err error) { log.Printf("[%s] NATS disconnected: %v", cfg.Source, err) }),
		nats.ReconnectHandler(func(nc *nats.Conn) { log.Printf("[%s] NATS reconnected: %s", cfg.Source, nc.ConnectedUrl()) }),
		nats.ErrorHandler(func(_ *nats.Conn, _ *nats.Subscription, err error) {
			log.Printf("[%s] NATS error: %v", cfg.Source, err)
		}),
	}
	if cfg.NATSCreds != "" {
		opts = append(opts, nats.UserCredentials(cfg.NATSCreds))
	} else if cfg.NATSUser != "" {
		opts = append(opts, nats.UserInfo(cfg.NATSUser, cfg.NATSPassword))
	}
	nc, err := nats.Connect(cfg.NATSURL, opts...)
	if err != nil {
		return err
	}
	defer nc.Drain()
	log.Printf("[%s] NATS connected: %s", cfg.Source, nc.ConnectedUrl())

	node := &gomavlib.Node{
		Endpoints:   []gomavlib.EndpointConf{cfg.Endpoint},
		Dialect:     ardupilotmega.Dialect,
		OutVersion:  cfg.OutVersion,
		OutSystemID: cfg.GCSSysID,
	}
	if err := node.Initialize(); err != nil {
		return err
	}
	defer node.Close()
	log.Printf("[%s] MAVLink endpoint %s, vehicle sysid %d, out %v", cfg.Source, cfg.EndpointS, cfg.SysID, cfg.OutVersion)

	veh := newMavVehicle(node, cfg.SysID, cfg.AckTimeout, cfg.AckRetries)
	b, err := NewBridge(cfg, natsBus{nc}, veh)
	if err != nil {
		return err
	}

	// Core subscription, not a durable consumer: an offline bridge gets no stale command
	// later (asset-command.md §4 "No queued replay").
	sub, err := nc.Subscribe(b.commandSubject(), func(m *nats.Msg) {
		go b.HandleCommand(ctx, m.Subject, m.Data)
	})
	if err != nil {
		return err
	}
	defer func() { _ = sub.Unsubscribe() }()
	if err := nc.Flush(); err != nil {
		return err
	}
	log.Printf("[%s] subscribed %s; publishing %s, %s, %s", cfg.Source, b.commandSubject(), b.positionSubject(), b.statusSubject(), b.outcomeSubject())
	if ready != nil {
		ready()
	}

	go b.RunStatus(ctx)
	pump(ctx, node, veh, b)
	log.Printf("[%s] shutting down", cfg.Source)
	return nil
}

// pump routes MAVLink frames to the vehicle (ACKs) and the bridge (telemetry) until ctx ends.
func pump(ctx context.Context, node *gomavlib.Node, veh *mavVehicle, b *Bridge) {
	for {
		select {
		case <-ctx.Done():
			return
		case evt, ok := <-node.Events():
			if !ok {
				return
			}
			switch e := evt.(type) {
			case *gomavlib.EventFrame:
				if m := veh.onFrame(e); m != nil {
					b.OnMessage(m)
				}
			case *gomavlib.EventChannelOpen:
				log.Printf("[%s] MAVLink channel open: %s", b.cfg.Source, e.Channel)
			case *gomavlib.EventChannelClose:
				log.Printf("[%s] MAVLink channel closed: %s", b.cfg.Source, e.Channel)
			}
		}
	}
}

func init() {
	log.SetOutput(os.Stderr)
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)
}
