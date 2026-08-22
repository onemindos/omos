// omd — OneMind daemon skeleton.
//
// This binary is the L3 Compute layer (ARCHITECTURE.md §2 / §3):
//   - Connects to NATS
//   - Registers the node as a polymorphic Actor in the om-config KV bucket
//   - Publishes evt.sys.node.online on startup
//   - Subscribes to cmd.node.{id}.ping for request/reply
//   - Registers a NATS micro service ("omd") for $SRV discovery
//   - Sends a heartbeat KV update every 15 s
//   - On SIGINT/SIGTERM: sets actor offline, publishes evt.sys.node.offline,
//     drains, and closes.
//
// SKELETON ONLY — NO fs.read/fs.write/shell.exec. Capabilities = ["host.info"].
// Any local-execution capability must be added in a later, capability-gated card.
//
// Connection notes:
//   - Currently uses the no_auth_user path (GLOBAL/onemind account) that the
//     local NATS server already allows without credentials.
//   - When FABRIC-account creds are issued (future card), pass --creds to
//     load a NATS .creds file. The switch does NOT require a code change,
//     only a flag/env value. No secrets are hardcoded here.
//
// Envelope note:
//   - evt.sys.node.online / offline payloads use a minimal {id, ts} shape today.
//     TODO: once P0-E (event envelope standard) is defined and written to om-schema,
//     wrap these events in the versioned envelope.
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/micro"

	"github.com/onemindos/onemindos-core/omd/internal/actor"
	"github.com/onemindos/onemindos-core/omd/internal/bus"
	"github.com/onemindos/onemindos-core/omd/internal/config"
)

func main() {
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})))

	cfg := config.Load()

	hostname, _ := os.Hostname()
	nodeActor := actor.New(cfg.NodeID, hostname)

	slog.Info("omd starting", "node", cfg.NodeID, "nats", cfg.NATSUrl)

	nc, err := bus.Connect(cfg.NATSUrl, cfg.NodeID, cfg.CredsFile)
	if err != nil {
		slog.Error("failed to connect to NATS", "err", err)
		os.Exit(1)
	}
	defer nc.Close()

	js, err := nc.JetStream()
	if err != nil {
		slog.Error("failed to get JetStream context", "err", err)
		os.Exit(1)
	}

	// -----------------------------------------------------------------
	// REGISTRATION — write actor record to om-config KV bucket
	// -----------------------------------------------------------------
	kv, err := js.KeyValue("om-config")
	if err != nil {
		slog.Error("failed to bind om-config KV", "err", err)
		os.Exit(1)
	}

	actorBytes, err := nodeActor.Marshal()
	if err != nil {
		slog.Error("failed to marshal actor", "err", err)
		os.Exit(1)
	}

	kvKey := actor.KVKey(cfg.NodeID)
	if _, err := kv.Put(kvKey, actorBytes); err != nil {
		slog.Error("failed to register actor in KV", "key", kvKey, "err", err)
		os.Exit(1)
	}
	slog.Info("actor registered", "key", kvKey, "capabilities", nodeActor.Capabilities)

	// Clean up any legacy keys created before node-id sanitization was enforced.
	// Previously, os.Hostname() returned the full FQDN (e.g. "OneMinds-MacBook-Pro.local")
	// which created a malformed subject-token key. Delete it now if present.
	if hostname != "" && hostname != cfg.NodeID {
		legacyKey := actor.KVKey(hostname)
		if err := kv.Delete(legacyKey); err == nil {
			slog.Info("deleted legacy actor key with dots", "key", legacyKey)
		}
	}

	// -----------------------------------------------------------------
	// PRESENCE EVENT — evt.sys.node.online
	//
	// TODO (P0-E): wrap this payload in the versioned event envelope once
	// P0-E is defined and published to om-schema.
	// -----------------------------------------------------------------
	onlineEvent := presenceEvent(cfg.NodeID, "online")
	if err := nc.Publish("evt.sys.node.online", onlineEvent); err != nil {
		slog.Warn("failed to publish online event", "err", err)
	} else {
		slog.Info("published evt.sys.node.online", "node", cfg.NodeID)
	}

	// -----------------------------------------------------------------
	// CMD SUBSCRIPTION — cmd.node.{id}.ping  (request/reply)
	// -----------------------------------------------------------------
	pingSub, err := nc.Subscribe(
		fmt.Sprintf("cmd.node.%s.ping", cfg.NodeID),
		func(msg *nats.Msg) {
			slog.Info("ping received", "reply", msg.Reply)
			if msg.Reply == "" {
				return // fire-and-forget, no reply needed
			}
			// Refresh last_seen before replying
			reply, _ := nodeActor.SetStatus("online").Marshal()
			if err := msg.Respond(reply); err != nil {
				slog.Warn("ping reply failed", "err", err)
			}
		},
	)
	if err != nil {
		slog.Error("failed to subscribe to ping", "err", err)
		os.Exit(1)
	}
	defer pingSub.Unsubscribe()
	slog.Info("subscribed", "subject", fmt.Sprintf("cmd.node.%s.ping", cfg.NodeID))

	// -----------------------------------------------------------------
	// MICRO SERVICE — $SRV "omd"
	// Registers with NATS micro so `nats micro ls` / $SRV.PING discovery
	// see this daemon. One endpoint: "info" returns the actor JSON.
	// Subject follows NATS-DOCTRINE §2 grammar: svc.{service}.{node-id}.{verb}
	// -----------------------------------------------------------------
	svc, err := micro.AddService(nc, micro.Config{
		Name:        "omd",
		Version:     "0.1.0",
		Description: "OneMind daemon — node actor registration and command gateway",
		Endpoint: &micro.EndpointConfig{
			Subject: fmt.Sprintf("svc.omd.%s.info", cfg.NodeID),
			Handler: micro.HandlerFunc(func(req micro.Request) {
				data, _ := nodeActor.SetStatus("online").Marshal()
				if err := req.Respond(data); err != nil {
					slog.Warn("micro info respond failed", "err", err)
				}
			}),
		},
	})
	if err != nil {
		slog.Error("failed to register micro service", "err", err)
		os.Exit(1)
	}
	defer svc.Stop()
	slog.Info("micro service registered", "service", svc.Info().Name, "id", svc.Info().ID)

	// -----------------------------------------------------------------
	// HEARTBEAT — every 15 s, update the actor KV entry with latest
	// last_seen timestamp and confirm status:online.
	// Choice: KV put on the actor key (not a short-TTL separate key).
	// Rationale: keeps one canonical record; consumers always read the
	// actor key; no secondary presence-key cleanup needed.
	// -----------------------------------------------------------------
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		ticker := time.NewTicker(15 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				data, _ := nodeActor.SetStatus("online").Marshal()
				if _, err := kv.Put(kvKey, data); err != nil {
					slog.Warn("heartbeat KV put failed", "err", err)
				} else {
					slog.Debug("heartbeat sent", "node", cfg.NodeID)
				}
			}
		}
	}()

	// -----------------------------------------------------------------
	// SHUTDOWN — wait for SIGINT or SIGTERM
	// -----------------------------------------------------------------
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	sig := <-quit
	slog.Info("shutdown signal received", "signal", sig)

	cancel() // stop heartbeat

	// Mark offline in KV
	offlineData, _ := nodeActor.SetStatus("offline").Marshal()
	if _, err := kv.Put(kvKey, offlineData); err != nil {
		slog.Warn("failed to set actor offline in KV", "err", err)
	} else {
		slog.Info("actor status set offline in KV", "key", kvKey)
	}

	// Publish presence-offline event
	// TODO (P0-E): wrap in versioned envelope when P0-E is defined.
	offlineEvent := presenceEvent(cfg.NodeID, "offline")
	if err := nc.Publish("evt.sys.node.offline", offlineEvent); err != nil {
		slog.Warn("failed to publish offline event", "err", err)
	} else {
		slog.Info("published evt.sys.node.offline", "node", cfg.NodeID)
	}

	// Drain — gracefully complete in-flight messages before close
	slog.Info("draining NATS connection...")
	if err := nc.Drain(); err != nil {
		slog.Warn("NATS drain error", "err", err)
	}
	slog.Info("omd stopped cleanly", "node", cfg.NodeID)
}

// presenceEvent builds a minimal presence payload.
// Shape: {"id": "<node-id>", "ts": "<rfc3339>", "event": "online|offline"}
// TODO (P0-E): replace with versioned envelope once P0-E defines the schema.
func presenceEvent(nodeID, event string) []byte {
	payload := map[string]string{
		"id":    nodeID,
		"ts":    time.Now().UTC().Format(time.RFC3339),
		"event": event,
	}
	b, _ := json.Marshal(payload)
	return b
}
