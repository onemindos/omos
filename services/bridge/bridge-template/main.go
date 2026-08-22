package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/nats-io/nats.go"
)

// BridgeMessage is the standard One Mind message envelope.
// Every bridge publishes this structure — no exceptions.
type BridgeMessage struct {
	ID      string          `json:"id"`      // ULID — unique, monotonic
	TS      int64           `json:"ts"`      // Unix ms
	Source  string          `json:"source"`  // bridge name
	Entity  string          `json:"entity"`  // domain.entity-id
	Domain  string          `json:"domain"`  // air / ground / home / field / etc.
	Payload json.RawMessage `json:"payload"` // protocol-specific, bridge-defined schema
}

// BridgeAdapter is the interface every bridge must implement.
type BridgeAdapter interface {
	Connect(ctx context.Context) error
	Subscribe(ctx context.Context, out chan<- BridgeMessage) error
	SendCommand(ctx context.Context, cmd Command) error
	Disconnect() error
}

// Command is the standard inbound command envelope.
type Command struct {
	ID      string          `json:"id"`
	TS      int64           `json:"ts"`
	Entity  string          `json:"entity"`
	Action  string          `json:"action"`
	Payload json.RawMessage `json:"payload"`
}

// Config holds all bridge config — loaded from environment only.
type Config struct {
	NATSUrl    string
	BridgeName string
	EntityID   string
	Domain     string
}

func loadConfig() Config {
	cfg := Config{
		NATSUrl:    getEnv("NATS_URL", "nats://localhost:4222"),
		BridgeName: getEnv("BRIDGE_NAME", "template-bridge"),
		EntityID:   getEnv("BRIDGE_ENTITY_ID", "template.device-1"),
		Domain:     getEnv("BRIDGE_DOMAIN", "home"),
	}
	return cfg
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func main() {
	cfg := loadConfig()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Connect to NATS with reconnect options
	nc, err := nats.Connect(cfg.NATSUrl,
		nats.MaxReconnects(-1),
		nats.ReconnectWait(2*time.Second),
		nats.DisconnectErrHandler(func(nc *nats.Conn, err error) {
			log.Printf("[%s] disconnected: %v", cfg.BridgeName, err)
		}),
		nats.ReconnectHandler(func(nc *nats.Conn) {
			log.Printf("[%s] reconnected to %s", cfg.BridgeName, nc.ConnectedUrl())
		}),
	)
	if err != nil {
		log.Fatalf("[%s] failed to connect to NATS: %v", cfg.BridgeName, err)
	}
	defer nc.Close()

	log.Printf("[%s] connected to NATS at %s", cfg.BridgeName, nc.ConnectedUrl())

	// Subscribe to inbound commands
	commandSubject := fmt.Sprintf("omos.%s.%s.command", cfg.Domain, cfg.EntityID)
	_, err = nc.Subscribe(commandSubject, func(msg *nats.Msg) {
		var cmd Command
		if err := json.Unmarshal(msg.Data, &cmd); err != nil {
			publishError(nc, cfg.BridgeName, fmt.Sprintf("invalid command envelope: %v", err))
			return
		}
		log.Printf("[%s] received command: %s", cfg.BridgeName, cmd.Action)
		// TODO: implement command dispatch to your protocol
	})
	if err != nil {
		log.Fatalf("[%s] failed to subscribe to commands: %v", cfg.BridgeName, err)
	}

	// Start heartbeat goroutine
	go runHeartbeat(ctx, nc, cfg)

	// TODO: Start your protocol listener here and publish telemetry:
	// publishTelemetry(nc, cfg, yourPayload)

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit
	log.Printf("[%s] shutting down gracefully", cfg.BridgeName)
	nc.Drain()
}

// publishTelemetry publishes a standard envelope to the fabric.
func publishTelemetry(nc *nats.Conn, cfg Config, payload interface{}) error {
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	msg := BridgeMessage{
		ID:      newULID(),
		TS:      time.Now().UnixMilli(),
		Source:  cfg.BridgeName,
		Entity:  cfg.EntityID,
		Domain:  cfg.Domain,
		Payload: payloadBytes,
	}

	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	subject := fmt.Sprintf("omos.%s.%s.telemetry", cfg.Domain, cfg.EntityID)
	return nc.Publish(subject, data)
}

// publishError reports a structured error to the bridge error subject.
func publishError(nc *nats.Conn, bridgeName, errMsg string) {
	payload, _ := json.Marshal(map[string]interface{}{
		"ts":    time.Now().UnixMilli(),
		"error": errMsg,
	})
	nc.Publish(fmt.Sprintf("omos.bridge.%s.errors", bridgeName), payload)
}

// runHeartbeat publishes a health signal every 30 seconds.
func runHeartbeat(ctx context.Context, nc *nats.Conn, cfg Config) {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			payload, _ := json.Marshal(map[string]interface{}{
				"status": "ok",
				"ts":     time.Now().UnixMilli(),
				"bridge": cfg.BridgeName,
			})
			nc.Publish(fmt.Sprintf("omos.bridge.%s.health", cfg.BridgeName), payload)
		}
	}
}

// newULID returns a simple timestamp-based ID.
// Replace with a real ULID library (github.com/oklog/ulid) in production.
func newULID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}
