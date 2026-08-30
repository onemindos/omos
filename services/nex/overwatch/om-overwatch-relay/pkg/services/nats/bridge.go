package nats

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/onemindos/om-overwatch-relay/pkg/services/config"
	"github.com/onemindos/om-overwatch-relay/pkg/services/shared"
	"github.com/nats-io/nats.go/jetstream"
)

// NewBridge creates a new bidirectional NATS JetStream bridge.
// The bridge handles both publishing telemetry to JetStream and consuming commands from it.
func NewBridge(js jetstream.JetStream, cfg *config.Config) *Bridge {
	ctx, cancel := context.WithCancel(context.Background())

	b := &Bridge{
		js:     js,
		config: cfg,
		ctx:    ctx,
		cancel: cancel,

		// Create buffered channels for each priority level (publish direction)
		criticalCh: make(chan *BridgeMessage, 100),  // Small buffer - immediate processing
		highCh:     make(chan *BridgeMessage, 500),  // Medium buffer
		mediumCh:   make(chan *BridgeMessage, 2000), // Large buffer for batching
		lowCh:      make(chan *BridgeMessage, 5000), // Largest buffer

		// Command channel (consume direction)
		commandCh: make(chan *BridgeMessage, 100), // Commands from Constellation

		// Configure batch sizes and timeouts per category
		batchSizes: map[MessageCategory]int{
			CategoryCritical: 1,   // No batching - immediate
			CategoryHigh:     10,  // Small batches
			CategoryMedium:   50,  // Medium batches
			CategoryLow:      100, // Large batches
		},
		batchTimeouts: map[MessageCategory]time.Duration{
			CategoryCritical: 0,                      // Immediate
			CategoryHigh:     5 * time.Millisecond,   // Very low latency
			CategoryMedium:   20 * time.Millisecond,  // Low latency
			CategoryLow:      100 * time.Millisecond, // Higher latency OK
		},
	}

	// Initialize statistics
	b.stats.Published = make(map[MessageCategory]int64)
	b.stats.Batched = make(map[MessageCategory]int64)
	b.stats.Errors = make(map[MessageCategory]int64)
	b.stats.AvgBatchSize = make(map[MessageCategory]float64)

	// Start background processors for each category (publish direction)
	b.wg.Add(4)
	go b.processCritical()
	go b.processHigh()
	go b.processMedium()
	go b.processLow()

	// Start statistics logger
	b.wg.Add(1)
	go b.logStatistics()

	log.Println("NATS Bridge started (publish + command channels ready)")
	return b
}

// PublishMessage queues a message for publishing to JetStream based on its priority.
func (b *Bridge) PublishMessage(subject string, payload []byte, messageType string, entityID string) error {
	// Determine message category from shared priority map
	priority, exists := shared.MessagePriorityMap[messageType]
	category := MessageCategory(priority)
	if !exists {
		// Default to medium priority for unknown message types
		category = CategoryMedium
	}

	msg := &BridgeMessage{
		Subject:     subject,
		Payload:     payload,
		Category:    category,
		Timestamp:   time.Now(),
		EntityID:    entityID,
		MessageType: messageType,
		Direction:   DirectionPublish,
	}

	// Route to appropriate channel
	select {
	case b.getCategoryChannel(category) <- msg:
		return nil
	case <-b.ctx.Done():
		return fmt.Errorf("bridge shutdown")
	default:
		// Channel full - this is a performance issue
		b.updateErrorStats(category)
		return fmt.Errorf("publish queue full for category %s", category.String())
	}
}

// GetCommandChannel returns the channel for receiving commands from Constellation.
// Commands are consumed from JetStream and forwarded to this channel for processing.
func (b *Bridge) GetCommandChannel() <-chan *BridgeMessage {
	return b.commandCh
}

// StartCommandConsumer starts consuming commands from the specified subject pattern.
// Commands received will be sent to the command channel for processing by the MAVLink router.
func (b *Bridge) StartCommandConsumer(subjectPattern string) error {
	// TODO: Implement command consumer
	// This will subscribe to command subjects like "constellation.{org_id}.{entity_id}.cmd.>"
	// and forward commands to the commandCh for the MAVLink router to process
	log.Printf("Command consumer ready for pattern: %s (not yet implemented)", subjectPattern)
	return nil
}

// getCategoryChannel returns the appropriate channel for a message category
func (b *Bridge) getCategoryChannel(category MessageCategory) chan<- *BridgeMessage {
	switch category {
	case CategoryCritical:
		return b.criticalCh
	case CategoryHigh:
		return b.highCh
	case CategoryMedium:
		return b.mediumCh
	case CategoryLow:
		return b.lowCh
	default:
		return b.mediumCh
	}
}

// processCritical handles critical messages with no batching
func (b *Bridge) processCritical() {
	defer b.wg.Done()

	for {
		select {
		case msg := <-b.criticalCh:
			b.publishSingle(msg)
		case <-b.ctx.Done():
			return
		}
	}
}

// processHigh handles high-priority messages with small batches
func (b *Bridge) processHigh() {
	defer b.wg.Done()
	b.processBatched(b.highCh, CategoryHigh)
}

// processMedium handles medium-priority messages with moderate batches
func (b *Bridge) processMedium() {
	defer b.wg.Done()
	b.processBatched(b.mediumCh, CategoryMedium)
}

// processLow handles low-priority messages with large batches
func (b *Bridge) processLow() {
	defer b.wg.Done()
	b.processBatched(b.lowCh, CategoryLow)
}

// processBatched handles batched publishing for a specific category
func (b *Bridge) processBatched(ch <-chan *BridgeMessage, category MessageCategory) {
	batchSize := b.batchSizes[category]
	timeout := b.batchTimeouts[category]

	var batch []*BridgeMessage
	ticker := time.NewTicker(timeout)
	defer ticker.Stop()

	for {
		select {
		case msg := <-ch:
			batch = append(batch, msg)

			// Publish when batch is full
			if len(batch) >= batchSize {
				b.publishBatch(batch, category)
				batch = batch[:0] // Reset batch
				ticker.Reset(timeout)
			}

		case <-ticker.C:
			// Publish partial batch on timeout
			if len(batch) > 0 {
				b.publishBatch(batch, category)
				batch = batch[:0]
			}

		case <-b.ctx.Done():
			// Publish any remaining messages before shutdown
			if len(batch) > 0 {
				b.publishBatch(batch, category)
			}
			return
		}
	}
}

// publishSingle publishes a single message immediately
func (b *Bridge) publishSingle(msg *BridgeMessage) {
	ctx, cancel := context.WithTimeout(b.ctx, 5*time.Second)
	defer cancel()

	_, err := b.js.Publish(ctx, msg.Subject, msg.Payload)
	if err != nil {
		b.updateErrorStats(msg.Category)
		if b.config.LogLevel == "debug" {
			log.Printf("Failed to publish %s message to %s: %v", msg.MessageType, msg.Subject, err)
		}
	} else {
		b.updatePublishStats(msg.Category, 1)
	}
}

// publishBatch publishes a batch of messages
func (b *Bridge) publishBatch(batch []*BridgeMessage, category MessageCategory) {
	if len(batch) == 0 {
		return
	}

	// Group messages by subject for efficient publishing
	subjectGroups := make(map[string][]*BridgeMessage)
	for _, msg := range batch {
		subjectGroups[msg.Subject] = append(subjectGroups[msg.Subject], msg)
	}

	// Publish each group
	published := 0
	for subject, msgs := range subjectGroups {
		for _, msg := range msgs {
			ctx, cancel := context.WithTimeout(b.ctx, 2*time.Second)

			_, err := b.js.Publish(ctx, subject, msg.Payload)
			cancel()

			if err != nil {
				b.updateErrorStats(category)
				if b.config.LogLevel == "debug" {
					log.Printf("Failed to publish %s: %v", msg.MessageType, err)
				}
				continue
			}

			published++
		}
	}

	// Update statistics
	b.updatePublishStats(category, published)
	b.updateBatchStats(category, len(batch))

	if b.config.LogLevel == "debug" {
		log.Printf("Published batch: %d messages, %d subjects, category: %s",
			len(batch), len(subjectGroups), category.String())
	}
}

// updatePublishStats updates publish statistics
func (b *Bridge) updatePublishStats(category MessageCategory, count int) {
	b.stats.Lock()
	b.stats.Published[category] += int64(count)
	b.stats.Unlock()
}

// updateErrorStats updates error statistics
func (b *Bridge) updateErrorStats(category MessageCategory) {
	b.stats.Lock()
	b.stats.Errors[category]++
	b.stats.Unlock()
}

// updateBatchStats updates batch statistics
func (b *Bridge) updateBatchStats(category MessageCategory, batchSize int) {
	b.stats.Lock()
	b.stats.Batched[category]++
	// Update average batch size (simple moving average)
	current := b.stats.AvgBatchSize[category]
	batches := b.stats.Batched[category]
	b.stats.AvgBatchSize[category] = (current*float64(batches-1) + float64(batchSize)) / float64(batches)
	b.stats.Unlock()
}

// logStatistics periodically logs bridge statistics
func (b *Bridge) logStatistics() {
	defer b.wg.Done()

	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			b.stats.RLock()
			log.Println("=== NATS Bridge Statistics ===")
			for category := CategoryCritical; category <= CategoryLow; category++ {
				published := b.stats.Published[category]
				batches := b.stats.Batched[category]
				errors := b.stats.Errors[category]
				avgBatch := b.stats.AvgBatchSize[category]

				if published > 0 || errors > 0 {
					log.Printf("  %s: %d published, %d batches, %.1f avg batch size, %d errors",
						category.String(), published, batches, avgBatch, errors)
				}
			}
			b.stats.RUnlock()

		case <-b.ctx.Done():
			return
		}
	}
}

// GetStatistics returns current bridge statistics
func (b *Bridge) GetStatistics() map[string]interface{} {
	b.stats.RLock()
	defer b.stats.RUnlock()

	stats := make(map[string]interface{})
	for category := CategoryCritical; category <= CategoryLow; category++ {
		categoryStats := map[string]interface{}{
			"published":      b.stats.Published[category],
			"batches":        b.stats.Batched[category],
			"errors":         b.stats.Errors[category],
			"avg_batch_size": b.stats.AvgBatchSize[category],
		}
		stats[category.String()] = categoryStats
	}

	return stats
}

// Close gracefully shuts down the bridge
func (b *Bridge) Close() error {
	log.Println("Shutting down NATS bridge...")

	// Signal shutdown
	b.cancel()

	// Close channels to trigger batch flushing
	close(b.criticalCh)
	close(b.highCh)
	close(b.mediumCh)
	close(b.lowCh)
	close(b.commandCh)

	// Wait for all goroutines to finish
	done := make(chan struct{})
	go func() {
		b.wg.Wait()
		close(done)
	}()

	// Wait with timeout
	select {
	case <-done:
		log.Println("NATS bridge shutdown completed")
		return nil
	case <-time.After(10 * time.Second):
		log.Println("NATS bridge shutdown timed out")
		return fmt.Errorf("shutdown timeout")
	}
}

// Legacy aliases for backward compatibility
// Deprecated: Use NewBridge instead
func NewAsyncPublisher(js jetstream.JetStream, cfg *config.Config) *Bridge {
	return NewBridge(js, cfg)
}

// AsyncPublisher is an alias for Bridge for backward compatibility
// Deprecated: Use Bridge instead
type AsyncPublisher = Bridge
