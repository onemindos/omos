package device

import (
	"sync"
	"time"

	"github.com/bluenviron/gomavlib/v3"
)

// DeviceFingerprint tracks per-device identification and routing info.
// This is the runtime representation used during message processing.
type DeviceFingerprint struct {
	SystemID       byte
	ComponentID    byte
	Channel        *gomavlib.Channel // Source channel for this device
	RemoteAddr     string            // UDP source address (for logging/debugging)
	OrgID          string            // Organization ID from overwatch.yaml (required)
	EntityID       string            // Entity ID from overwatch.yaml or auto-generated
	Name           string            // Human-readable name from overwatch.yaml
	Subject        string            // Computed hierarchical subject (e.g., "constellation.telemetry.org.entity")
	Fingerprinted  bool
	FirstSeen      time.Time
	LastSeen       time.Time
	LastCoreUpdate time.Time // Last time core entity record was updated (for rate-limiting)
	Mu             sync.Mutex
}

// DeviceState represents the CORE entity record in the KV store.
// This is a lean metadata record stored at {entity_id} key.
// Individual telemetry signals are stored in separate keys:
//   - {entity_id}.mavlink.heartbeat
//   - {entity_id}.mavlink.position
//   - {entity_id}.mavlink.attitude
//   - {entity_id}.mavlink.power
//   - {entity_id}.mavlink.flight
type DeviceState struct {
	SystemID      byte      `json:"system_id"`
	ComponentID   byte      `json:"component_id"`
	OrgID         string    `json:"org_id"`         // Organization ID from overwatch.yaml (required)
	EntityID      string    `json:"entity_id"`      // Entity ID (e.g., gcs001)
	Name          string    `json:"name,omitempty"` // Human-readable name from overwatch.yaml
	StreamPort    string    `json:"stream_port"`    // Source port/channel
	Subject       string    `json:"subject"`        // Hierarchical NATS subject
	FirstSeen     time.Time `json:"first_seen"`
	LastSeen      time.Time `json:"last_seen"`
	Fingerprinted bool      `json:"fingerprinted"` // Whether stream has been identified
}
