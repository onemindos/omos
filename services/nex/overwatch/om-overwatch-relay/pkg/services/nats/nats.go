package nats

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/onemindos/om-overwatch-relay/pkg/services/config"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

// Connect establishes a connection to the NATS server
func Connect(cfg *config.Config) (*nats.Conn, error) {
	opts := []nats.Option{
		nats.Name("mavlink-nats-bridge"),
		nats.MaxReconnects(-1),
		nats.ReconnectWait(2 * time.Second),
		nats.DisconnectErrHandler(func(nc *nats.Conn, err error) {
			if err != nil {
				log.Printf("NATS disconnected: %v", err)
			}
		}),
		nats.ReconnectHandler(func(nc *nats.Conn) {
			log.Printf("NATS reconnected to %s", nc.ConnectedUrl())
		}),
	}

	// Add auth token if configured
	if cfg.NATSAuthToken != "" {
		opts = append(opts, nats.Token(cfg.NATSAuthToken))
	}

	nc, err := nats.Connect(cfg.NATSUrl, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to NATS: %w", err)
	}
	log.Println("Connected to NATS")
	return nc, nil
}

// SetupJetStream initializes JetStream context and creates/updates the stream
func SetupJetStream(nc *nats.Conn, cfg *config.Config) (jetstream.JetStream, jetstream.Stream, error) {
	// Create JetStream context
	js, err := jetstream.New(nc)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create JetStream context: %w", err)
	}

	// Setup JetStream stream
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	storageType := jetstream.FileStorage
	if cfg.NATSStorageType == "memory" {
		storageType = jetstream.MemoryStorage
	}

	// Stream supports hierarchical subjects: constellation.telemetry.{org_id}.{entity_id}
	// Use wildcard to match all telemetry: constellation.telemetry.>
	streamConfig := jetstream.StreamConfig{
		Name:        cfg.NATSStreamName,
		Description: "om-overwatch MAVLink telemetry — friendly UAV streams (tel.drone.>)",
		Subjects:    []string{"tel.drone.>"}, // OneMind telemetry taxonomy
		Retention:   jetstream.LimitsPolicy,
		MaxAge:      cfg.NATSMaxAge,
		MaxMsgs:     cfg.NATSMaxMsgs,
		MaxBytes:    cfg.NATSMaxBytes,
		Storage:     storageType,
		Replicas:    1,
		Discard:     jetstream.DiscardOld,
	}

	// Try to create or update stream
	stream, err := js.CreateStream(ctx, streamConfig)
	if err != nil {
		// If stream already exists, try to get it
		log.Printf("Stream creation returned error (may already exist): %v", err)
		stream, err = js.Stream(ctx, cfg.NATSStreamName)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to create/get stream '%s': %w\nEnsure NATS server has JetStream enabled (nats-server -js)", cfg.NATSStreamName, err)
		}
		log.Printf("Using existing JetStream stream '%s'", cfg.NATSStreamName)
	} else {
		log.Printf("Created new JetStream stream '%s'", cfg.NATSStreamName)
	}

	// Verify stream configuration
	streamInfo, err := stream.Info(ctx)
	if err != nil {
		log.Printf("Warning: Could not get stream info: %v", err)
	} else {
		log.Printf("Stream info: subjects=%v, messages=%d, bytes=%d, storage=%s",
			streamInfo.Config.Subjects, streamInfo.State.Msgs, streamInfo.State.Bytes, streamInfo.Config.Storage)
	}
	log.Printf("JetStream stream '%s' ready", cfg.NATSStreamName)

	return js, stream, nil
}

// TestConnectivity performs a test publish to verify JetStream connectivity
func TestConnectivity(js jetstream.JetStream, cfg *config.Config) {
	testSubject := "tel.drone.test.connectivity"
	testData := []byte(`{"test":"connectivity","timestamp":"` + time.Now().Format(time.RFC3339) + `"}`)
	testCtx, testCancel := context.WithTimeout(context.Background(), 5*time.Second)
	ack, err := js.Publish(testCtx, testSubject, testData)
	testCancel()
	if err != nil {
		log.Printf("⚠ WARNING: JetStream test publish FAILED: %v", err)
		log.Printf("  Test subject: %s", testSubject)
		log.Printf("  This usually means:")
		log.Printf("    1. JetStream is not enabled on NATS server (start with: nats-server -js)")
		log.Printf("    2. The stream subjects pattern doesn't match")
		log.Printf("    3. NATS server is overloaded or unreachable")
		log.Printf("  Bridge will continue but publishing will likely fail!")
	} else {
		log.Printf("✓ JetStream connectivity test PASSED (stream: %s, seq: %d)", ack.Stream, ack.Sequence)
	}
}

// GetKVBucket connects to the KV bucket for device state
func GetKVBucket(js jetstream.JetStream, cfg *config.Config) (jetstream.KeyValue, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	kv, err := js.KeyValue(ctx, cfg.NATSKVBucket)
	if err != nil {
		return nil, fmt.Errorf("KV bucket '%s' not available: %w", cfg.NATSKVBucket, err)
	}
	log.Printf("Connected to KV bucket '%s'", cfg.NATSKVBucket)
	return kv, nil
}
