package nats

import (
	"context"
	"sync"
	"time"

	"github.com/onemindos/om-overwatch-relay/pkg/services/config"
	"github.com/nats-io/nats.go/jetstream"
)

// MessageCategory represents the priority and batching strategy for different message types.
type MessageCategory int

const (
	CategoryCritical MessageCategory = iota // Immediate publish - commands, safety
	CategoryHigh                            // High priority, low batching - navigation, heartbeat
	CategoryMedium                          // Medium priority, moderate batching - telemetry
	CategoryLow                             // Low priority, high batching - sensors, simulation
)

// String returns the string representation of a MessageCategory.
func (c MessageCategory) String() string {
	switch c {
	case CategoryCritical:
		return "critical"
	case CategoryHigh:
		return "high"
	case CategoryMedium:
		return "medium"
	case CategoryLow:
		return "low"
	default:
		return "unknown"
	}
}

// MessageDirection indicates whether a message is being published or consumed.
type MessageDirection int

const (
	DirectionPublish MessageDirection = iota // Message being published to JetStream
	DirectionConsume                         // Message consumed from JetStream (command)
)

// BridgeMessage represents a message flowing through the bridge in either direction.
type BridgeMessage struct {
	Subject     string
	Payload     []byte
	Category    MessageCategory
	Timestamp   time.Time
	EntityID    string
	MessageType string
	Direction   MessageDirection
}

// BridgeStatistics holds per-category bridge statistics.
type BridgeStatistics struct {
	sync.RWMutex
	Published    map[MessageCategory]int64
	Batched      map[MessageCategory]int64
	Errors       map[MessageCategory]int64
	AvgBatchSize map[MessageCategory]float64
	Consumed     int64 // Commands consumed from JetStream
}

// Bridge handles bidirectional NATS JetStream communication.
// It publishes telemetry to JetStream and consumes commands from it.
type Bridge struct {
	js     jetstream.JetStream
	config *config.Config

	// Publish channels by category (telemetry → JetStream)
	criticalCh chan *BridgeMessage
	highCh     chan *BridgeMessage
	mediumCh   chan *BridgeMessage
	lowCh      chan *BridgeMessage

	// Command channel (JetStream → MAVLink)
	commandCh chan *BridgeMessage

	// Batch configuration
	batchSizes    map[MessageCategory]int
	batchTimeouts map[MessageCategory]time.Duration

	// Statistics
	stats BridgeStatistics

	// Control
	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup
}

// Legacy type aliases for backward compatibility
// Deprecated: Use BridgeMessage instead
type BatchedMessage = BridgeMessage

// Deprecated: Use BridgeStatistics instead
type PublisherStatistics = BridgeStatistics
