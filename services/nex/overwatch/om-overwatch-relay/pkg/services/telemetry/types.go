package telemetry

import (
	"sync"
	"time"
)

// GlobalStats tracks overall statistics for telemetry processing.
type GlobalStats struct {
	Received  int64
	Published int64
	Errors    int64
	Routed    int64
	Mu        sync.Mutex
}

// TelemetryMessage represents a parsed MAVLink message ready for publishing.
type TelemetryMessage struct {
	Timestamp   time.Time              `json:"timestamp"`
	MessageType string                 `json:"message_type"`
	MessageID   uint32                 `json:"message_id"`
	SystemID    byte                   `json:"system_id"`
	ComponentID byte                   `json:"component_id"`
	Data        map[string]interface{} `json:"data"`
}
