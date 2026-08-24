package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/onemindos/om-overwatch-relay/pkg/services/config"
	"github.com/onemindos/om-overwatch-relay/pkg/services/device"
	"github.com/onemindos/om-overwatch-relay/pkg/services/mavlink"
	"github.com/onemindos/om-overwatch-relay/pkg/services/nats"
	"github.com/onemindos/om-overwatch-relay/pkg/services/telemetry"
	"github.com/bluenviron/gomavlib/v3"
	"github.com/nats-io/nats.go/jetstream"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Println("Starting MAVLink2Constellation - Multi-Entity Router + NATS Bridge...")

	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Enhanced configuration logging
	log.Printf("Configuration loaded:")
	if len(cfg.SerialDevices) > 0 {
		log.Printf("  Multi-Serial Devices: %v (bauds: %v, entities: %v)",
			cfg.SerialDevices, cfg.SerialBauds, cfg.SerialEntityIDs)
	} else if cfg.SerialDevice != "" {
		log.Printf("  Serial Device: %s (baud: %d)", cfg.SerialDevice, cfg.SerialBaud)
	}
	log.Printf("  Discovery Mode: %v", cfg.DiscoveryMode)
	if cfg.DiscoveryMode {
		log.Printf("  Discovery Address: %s:%s", cfg.DiscoveryHost, cfg.DiscoveryPort)
	} else {
		log.Printf("  MAVLink UDP Servers: %v", cfg.MAVLinkUDPPorts)
	}
	if cfg.EnableTCPServer {
		log.Printf("  TCP Server: %s:%s", cfg.TCPServerHost, cfg.TCPServerPort)
	}
	log.Printf("  UDP Forwards: %v", cfg.UDPForwards)
	if len(cfg.TCPForwards) > 0 {
		log.Printf("  TCP Forwards: %v", cfg.TCPForwards)
	}
	log.Printf("  NATS URL: %s", cfg.NATSUrl)
	log.Printf("  NATS Stream: %s", cfg.NATSStreamName)
	log.Printf("  NATS Root Subject: %s", cfg.NATSRootSubject)
	log.Printf("  NATS KV Bucket: %s", cfg.NATSKVBucket)
	log.Printf("  Default Org ID: %s", cfg.DefaultOrgID)
	log.Printf("  Entity Mappings: %v", cfg.EntityMappings)
	log.Printf("  Async Publishing: %v", cfg.EnableAsyncPublishing)
	log.Printf("  Metrics Enabled: %v", cfg.EnableMetrics)

	// Connect to NATS
	nc, err := nats.Connect(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to NATS: %v", err)
	}
	defer nc.Close()

	// Setup JetStream
	js, _, err := nats.SetupJetStream(nc, cfg)
	if err != nil {
		log.Fatalf("Failed to setup JetStream: %v", err)
	}

	// Test NATS JetStream connectivity
	nats.TestConnectivity(js, cfg)

	// Connect to existing KV bucket for device state (managed by Overwatch)
	kv, err := nats.GetKVBucket(js, cfg)
	if err != nil {
		log.Printf("Warning: %v", err)
		log.Println("Device state updates will be disabled. Ensure Overwatch is running.")
		kv = nil // Continue without KV updates
	}

	// Choose publisher based on configuration
	var asyncPublisher *nats.AsyncPublisher
	if cfg.EnableAsyncPublishing {
		asyncPublisher = nats.NewAsyncPublisher(js, cfg)
		defer func() {
			if err := asyncPublisher.Close(); err != nil {
				log.Printf("Error closing async publisher: %v", err)
			}
		}()
		log.Println("Enhanced async publisher enabled with batching")
	}

	// Create unified MAVLink router node
	node, err := mavlink.CreateNode(cfg)
	if err != nil {
		log.Fatalf("Failed to create MAVLink router node: %v", err)
	}
	defer node.Close()

	// Enhanced device tracking with performance metrics
	devices := make(map[string]*device.DeviceFingerprint)
	var devicesMu sync.RWMutex

	// Enhanced statistics tracking
	stats := &EnhancedStats{
		StartTime:  time.Now(),
		PerMessage: make(map[string]*MessageStats),
	}

	// Setup graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Enhanced statistics logger
	statsTicker := time.NewTicker(30 * time.Second)
	defer statsTicker.Stop()

	// HTTP metrics endpoint for monitoring
	if cfg.EnableMetrics {
		go startMetricsServer(cfg, stats, asyncPublisher, devices, &devicesMu)
		log.Printf("Metrics available at http://localhost:%s/metrics", cfg.MetricsPort)
	}

	log.Println("Multi-Entity Router + Bridge is running. Press Ctrl+C to stop.")

	// Enhanced statistics logging
	go func() {
		for range statsTicker.C {
			if cfg.EnableAsyncPublishing && asyncPublisher != nil {
				logEnhancedStatistics(stats, devices, &devicesMu, asyncPublisher)
			} else {
				logBasicStatistics(stats, devices, &devicesMu, cfg)
			}
		}
	}()

	// Main enhanced event loop
	go func() {
		for evt := range node.Events() {
			switch e := evt.(type) {
			case *gomavlib.EventChannelOpen:
				handleChannelOpen(e, cfg)

			case *gomavlib.EventChannelClose:
				handleChannelClose(e, cfg, devices, &devicesMu)

			case *gomavlib.EventFrame:
				if cfg.EnableAsyncPublishing && asyncPublisher != nil {
					handleFrameEnhanced(e, node, cfg, asyncPublisher, devices, &devicesMu, kv, stats)
				} else {
					handleFrameBasic(e, node, cfg, js, devices, &devicesMu, kv, stats)
				}

			case *gomavlib.EventParseError:
				handleParseError(e, stats)
			}
		}
	}()

	// Wait for shutdown signal
	<-sigChan
	log.Println("Shutdown signal received, closing router...")

	// Print final enhanced statistics
	logFinalStatistics(stats, devices, &devicesMu, asyncPublisher)
	log.Println("Shutdown complete")
}

// Enhanced statistics structure
type EnhancedStats struct {
	mu          sync.RWMutex
	StartTime   time.Time
	Received    int64
	Routed      int64
	Published   int64
	Errors      int64
	ParseErrors int64
	PerMessage  map[string]*MessageStats
}

type MessageStats struct {
	Count     int64
	LastSeen  time.Time
	AvgSize   float64
	TotalSize int64
}

// Enhanced frame handling with async publishing
func handleFrameEnhanced(e *gomavlib.EventFrame, node *gomavlib.Node, cfg *config.Config,
	asyncPub *nats.AsyncPublisher, devices map[string]*device.DeviceFingerprint,
	devicesMu *sync.RWMutex, kv jetstream.KeyValue, stats *EnhancedStats) {

	stats.mu.Lock()
	stats.Received++
	stats.mu.Unlock()

	// ROUTE: Forward frame to all endpoints except source
	node.WriteFrameExcept(e.Channel, e.Frame)
	stats.mu.Lock()
	stats.Routed++
	stats.mu.Unlock()

	// PARSE: Convert to telemetry message
	telemetryMsg, err := telemetry.ParseMavlinkMessage(e)
	if err != nil {
		stats.mu.Lock()
		stats.Errors++
		stats.mu.Unlock()
		if cfg.LogLevel == "debug" {
			log.Printf("Failed to parse message: %v", err)
		}
		return
	}

	// Update per-message statistics
	updateMessageStats(stats, telemetryMsg)

	// FINGERPRINT: Get or create device fingerprint
	devicesMu.Lock()
	dev := device.GetOrCreateDeviceFingerprint(devices, telemetryMsg, e.Channel, cfg, kv)
	devicesMu.Unlock()

	// Skip publishing until device is fingerprinted
	if !dev.Fingerprinted {
		if cfg.LogLevel == "debug" {
			log.Printf("[%d/%d] %s (NOT fingerprinted - waiting for HEARTBEAT)",
				telemetryMsg.SystemID, telemetryMsg.ComponentID, telemetryMsg.MessageType)
		}
		return
	}

	// ASYNC PUBLISH: Send to NATS JetStream using high-performance async publisher
	jsonData, err := json.Marshal(telemetryMsg)
	if err != nil {
		stats.mu.Lock()
		stats.Errors++
		stats.mu.Unlock()
		if cfg.LogLevel == "debug" {
			log.Printf("Failed to marshal message: %v", err)
		}
		return
	}

	// Use async publisher with automatic batching and priority routing
	err = asyncPub.PublishMessage(dev.Subject, jsonData, telemetryMsg.MessageType, dev.EntityID)
	if err != nil {
		stats.mu.Lock()
		stats.Errors++
		stats.mu.Unlock()
		// Only log publish queue full errors (critical performance issue)
		if cfg.LogLevel == "debug" || err.Error() == "publish queue full" {
			log.Printf("Async publish failed [subject: %s, msg: %s]: %v",
				dev.Subject, telemetryMsg.MessageType, err)
		}
		return
	}

	stats.mu.Lock()
	stats.Published++
	stats.mu.Unlock()

	// UPDATE KV: Update device state if available
	if kv != nil && cfg.EnableKVPublishing {
		kvCtx, kvCancel := context.WithTimeout(context.Background(), cfg.KVUpdateTimeout)
		if err := device.UpdateDeviceState(kvCtx, kv, dev, telemetryMsg); err != nil {
			if cfg.LogLevel == "debug" {
				log.Printf("Failed to update device state: %v", err)
			}
		}
		kvCancel()
	}

	// Log important messages
	if cfg.LogLevel == "debug" || telemetryMsg.MessageType == "HEARTBEAT" {
		log.Printf("[Device %d/%d → %s] %s",
			telemetryMsg.SystemID, telemetryMsg.ComponentID,
			dev.Subject, telemetryMsg.MessageType)
	}
}

// Basic frame handling with sync publishing (fallback)
func handleFrameBasic(e *gomavlib.EventFrame, node *gomavlib.Node, cfg *config.Config,
	js jetstream.JetStream, devices map[string]*device.DeviceFingerprint,
	devicesMu *sync.RWMutex, kv jetstream.KeyValue, stats *EnhancedStats) {

	stats.mu.Lock()
	stats.Received++
	stats.mu.Unlock()

	// ROUTE: Forward frame to all endpoints except source
	node.WriteFrameExcept(e.Channel, e.Frame)
	stats.mu.Lock()
	stats.Routed++
	stats.mu.Unlock()

	// PARSE: Convert to telemetry message
	telemetryMsg, err := telemetry.ParseMavlinkMessage(e)
	if err != nil {
		stats.mu.Lock()
		stats.Errors++
		stats.mu.Unlock()
		if cfg.LogLevel == "debug" {
			log.Printf("Failed to parse message: %v", err)
		}
		return
	}

	// Update per-message statistics
	updateMessageStats(stats, telemetryMsg)

	// FINGERPRINT: Get or create device fingerprint
	devicesMu.Lock()
	dev := device.GetOrCreateDeviceFingerprint(devices, telemetryMsg, e.Channel, cfg, kv)
	devicesMu.Unlock()

	// Skip publishing until device is fingerprinted
	if !dev.Fingerprinted {
		if cfg.LogLevel == "debug" {
			log.Printf("[%d/%d] %s (NOT fingerprinted - waiting for HEARTBEAT)",
				telemetryMsg.SystemID, telemetryMsg.ComponentID, telemetryMsg.MessageType)
		}
		return
	}

	// PUBLISH: Send to NATS JetStream (synchronous)
	jsonData, err := json.Marshal(telemetryMsg)
	if err != nil {
		stats.mu.Lock()
		stats.Errors++
		stats.mu.Unlock()
		if cfg.LogLevel == "debug" {
			log.Printf("Failed to marshal message: %v", err)
		}
		return
	}

	// Use sync publish with timeout
	pubCtx, pubCancel := context.WithTimeout(context.Background(), 10*time.Second)
	_, err = js.Publish(pubCtx, dev.Subject, jsonData)
	pubCancel()

	if err != nil {
		stats.mu.Lock()
		stats.Errors++
		stats.mu.Unlock()
		log.Printf("NATS publish failed [subject: %s, msg: %s]: %v",
			dev.Subject, telemetryMsg.MessageType, err)
		return
	}

	stats.mu.Lock()
	stats.Published++
	stats.mu.Unlock()

	// UPDATE KV: Update device state if available
	if kv != nil && cfg.EnableKVPublishing {
		kvCtx, kvCancel := context.WithTimeout(context.Background(), cfg.KVUpdateTimeout)
		if err := device.UpdateDeviceState(kvCtx, kv, dev, telemetryMsg); err != nil {
			if cfg.LogLevel == "debug" {
				log.Printf("Failed to update device state: %v", err)
			}
		}
		kvCancel()
	}

	// Log important messages
	if cfg.LogLevel == "debug" || telemetryMsg.MessageType == "HEARTBEAT" {
		log.Printf("[Device %d/%d → %s] %s",
			telemetryMsg.SystemID, telemetryMsg.ComponentID,
			dev.Subject, telemetryMsg.MessageType)
	}
}

func handleChannelOpen(e *gomavlib.EventChannelOpen, cfg *config.Config) {
	if cfg.DiscoveryMode {
		log.Printf("✓ MAVLink source discovered: %s (waiting for HEARTBEAT...)", e.Channel)
	} else {
		log.Printf("MAVLink channel opened: %s", e.Channel)
	}
}

func handleChannelClose(e *gomavlib.EventChannelClose, cfg *config.Config,
	devices map[string]*device.DeviceFingerprint, devicesMu *sync.RWMutex) {

	// Check if we had a device on this channel
	devicesMu.RLock()
	var closedDevices []string
	for key, dev := range devices {
		if dev.Channel == e.Channel {
			closedDevices = append(closedDevices, key)
		}
	}
	devicesMu.RUnlock()

	if cfg.DiscoveryMode && len(closedDevices) > 0 {
		log.Printf("✗ MAVLink source disconnected: %s (devices: %v)", e.Channel, closedDevices)
	} else {
		log.Printf("MAVLink channel closed: %s", e.Channel)
	}
}

func handleParseError(e *gomavlib.EventParseError, stats *EnhancedStats) {
	stats.mu.Lock()
	stats.Errors++
	stats.ParseErrors++

	// Simple rate limiting: only log first 10 parse errors
	if stats.ParseErrors <= 10 {
		log.Printf("MAVLink parse error #%d: %v", stats.ParseErrors, e.Error)
	} else if stats.ParseErrors == 11 {
		log.Printf("MAVLink parse errors suppressed after 10 (total: %d) - this is normal during startup", stats.ParseErrors)
	}
	stats.mu.Unlock()
}

func updateMessageStats(stats *EnhancedStats, msg *telemetry.TelemetryMessage) {
	stats.mu.Lock()
	defer stats.mu.Unlock()

	msgStats, exists := stats.PerMessage[msg.MessageType]
	if !exists {
		msgStats = &MessageStats{}
		stats.PerMessage[msg.MessageType] = msgStats
	}

	// Estimate message size (rough approximation)
	estimatedSize := int64(len(msg.MessageType) + 100) // Base overhead + approximate data size

	msgStats.Count++
	msgStats.LastSeen = time.Now()
	msgStats.TotalSize += estimatedSize
	msgStats.AvgSize = float64(msgStats.TotalSize) / float64(msgStats.Count)
}

func logEnhancedStatistics(stats *EnhancedStats, devices map[string]*device.DeviceFingerprint,
	devicesMu *sync.RWMutex, asyncPub *nats.AsyncPublisher) {

	stats.mu.RLock()
	devicesMu.RLock()

	uptime := time.Since(stats.StartTime)
	rateRx := float64(stats.Received) / uptime.Seconds()
	rateTx := float64(stats.Published) / uptime.Seconds()

	log.Println("=== Enhanced Router Statistics ===")
	log.Printf("  Uptime: %v", uptime.Truncate(time.Second))
	log.Printf("  Messages: %d rx, %d routed, %d published, %d errors",
		stats.Received, stats.Routed, stats.Published, stats.Errors)
	log.Printf("  Rates: %.1f rx/s, %.1f pub/s", rateRx, rateTx)
	log.Printf("  Parse errors: %d", stats.ParseErrors)
	log.Printf("  Devices tracked: %d", len(devices))

	// Show device details
	for key, dev := range devices {
		if dev.Fingerprinted {
			log.Printf("    [%s] → %s (subject: %s)", key, dev.EntityID, dev.Subject)
		} else {
			log.Printf("    [%s] → waiting for HEARTBEAT...", key)
		}
	}

	// Show top message types
	type msgCount struct {
		msgType string
		count   int64
		avgSize float64
	}
	var topMsgs []msgCount
	for msgType, msgStats := range stats.PerMessage {
		topMsgs = append(topMsgs, msgCount{msgType, msgStats.Count, msgStats.AvgSize})
	}

	// Sort by count (simple bubble sort for small data)
	for i := 0; i < len(topMsgs)-1; i++ {
		for j := 0; j < len(topMsgs)-i-1; j++ {
			if topMsgs[j].count < topMsgs[j+1].count {
				topMsgs[j], topMsgs[j+1] = topMsgs[j+1], topMsgs[j]
			}
		}
	}

	log.Printf("  Top message types:")
	for i, msg := range topMsgs {
		if i >= 5 { // Show top 5
			break
		}
		log.Printf("    %s: %d (avg %.0f bytes)", msg.msgType, msg.count, msg.avgSize)
	}

	devicesMu.RUnlock()
	stats.mu.RUnlock()

	// Show async publisher statistics
	if asyncPub != nil {
		pubStats := asyncPub.GetStatistics()
		log.Printf("  Async Publisher Stats:")
		for category, categoryStats := range pubStats {
			if statsMap, ok := categoryStats.(map[string]interface{}); ok {
				if published, ok := statsMap["published"].(int64); ok && published > 0 {
					log.Printf("    %s: %d published, %.1f avg batch",
						category, published, statsMap["avg_batch_size"].(float64))
				}
			}
		}
	}
}

func logBasicStatistics(stats *EnhancedStats, devices map[string]*device.DeviceFingerprint,
	devicesMu *sync.RWMutex, cfg *config.Config) {

	stats.mu.RLock()
	devicesMu.RLock()
	defer stats.mu.RUnlock()
	defer devicesMu.RUnlock()

	log.Println("=== Router Statistics ===")
	log.Printf("  Received: %d, Routed: %d, Published: %d, Errors: %d",
		stats.Received, stats.Routed, stats.Published, stats.Errors)
	log.Printf("  Devices tracked: %d", len(devices))
	for key, dev := range devices {
		if dev.Fingerprinted {
			if cfg.DiscoveryMode {
				log.Printf("    [%s] from %s → %s (subject: %s)", key, dev.RemoteAddr, dev.EntityID, dev.Subject)
			} else {
				log.Printf("    [%s] → %s (subject: %s)", key, dev.EntityID, dev.Subject)
			}
		} else {
			if cfg.DiscoveryMode {
				log.Printf("    [%s] from %s → waiting for HEARTBEAT...", key, dev.RemoteAddr)
			} else {
				log.Printf("    [%s] → waiting for HEARTBEAT...", key)
			}
		}
	}
}

func logFinalStatistics(stats *EnhancedStats, devices map[string]*device.DeviceFingerprint,
	devicesMu *sync.RWMutex, asyncPub *nats.AsyncPublisher) {

	stats.mu.RLock()
	devicesMu.RLock()
	defer stats.mu.RUnlock()
	defer devicesMu.RUnlock()

	uptime := time.Since(stats.StartTime)

	log.Println("=== Final Statistics ===")
	log.Printf("  Total uptime: %v", uptime.Truncate(time.Second))
	log.Printf("  Total messages: %d received, %d routed, %d published, %d errors",
		stats.Received, stats.Routed, stats.Published, stats.Errors)
	if stats.Received > 0 {
		log.Printf("  Final rates: %.1f rx/s, %.1f pub/s",
			float64(stats.Received)/uptime.Seconds(),
			float64(stats.Published)/uptime.Seconds())
		log.Printf("  Success rate: %.2f%%",
			float64(stats.Published)*100.0/float64(stats.Received))
	}
	log.Printf("  Total devices: %d", len(devices))
	log.Printf("  Message types processed: %d", len(stats.PerMessage))
}

// HTTP metrics endpoint for monitoring
func startMetricsServer(cfg *config.Config, stats *EnhancedStats, asyncPub *nats.AsyncPublisher,
	devices map[string]*device.DeviceFingerprint, devicesMu *sync.RWMutex) {

	http.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		stats.mu.RLock()
		devicesMu.RLock()

		uptime := time.Since(stats.StartTime)

		metrics := map[string]interface{}{
			"uptime_seconds":     uptime.Seconds(),
			"messages_received":  stats.Received,
			"messages_routed":    stats.Routed,
			"messages_published": stats.Published,
			"messages_errors":    stats.Errors,
			"parse_errors":       stats.ParseErrors,
			"devices_tracked":    len(devices),
			"message_types":      len(stats.PerMessage),
			"rates": map[string]float64{
				"rx_per_second":  float64(stats.Received) / uptime.Seconds(),
				"pub_per_second": float64(stats.Published) / uptime.Seconds(),
			},
		}

		if asyncPub != nil {
			metrics["async_publisher"] = asyncPub.GetStatistics()
		}

		devicesMu.RUnlock()
		stats.mu.RUnlock()

		json.NewEncoder(w).Encode(metrics)
	})

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	http.HandleFunc("/entities", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		devicesMu.RLock()
		defer devicesMu.RUnlock()

		entities := make(map[string]interface{})
		for key, dev := range devices {
			entities[key] = map[string]interface{}{
				"entity_id":     dev.EntityID,
				"subject":       dev.Subject,
				"fingerprinted": dev.Fingerprinted,
				"remote_addr":   dev.RemoteAddr,
			}
		}

		json.NewEncoder(w).Encode(entities)
	})

	log.Printf("Starting metrics server on :%s", cfg.MetricsPort)
	if err := http.ListenAndServe(":"+cfg.MetricsPort, nil); err != nil {
		log.Printf("Metrics server error: %v", err)
	}
}
