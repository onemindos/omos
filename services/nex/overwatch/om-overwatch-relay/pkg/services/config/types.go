package config

import (
	"strconv"
	"time"
)

// Config holds all application configuration loaded from environment variables and overwatch.yaml.
type Config struct {
	// Multi-Serial Device Support (legacy - prefer overwatch.yaml entities)
	SerialDevices    []string // e.g., ["/dev/ttyUSB0", "/dev/ttyACM0"]
	SerialBauds      []int    // e.g., [57600, 115200]
	SerialEntityIDs  []string // e.g., ["uav001", "gimbal001"]
	SerialAutoDetect bool     // Auto-detect serial devices

	// Legacy single serial support (deprecated)
	SerialDevice string // e.g., "/dev/ttyUSB0"
	SerialBaud   int    // e.g., 57600

	// MAVLink Discovery Mode (recommended for network deployments)
	DiscoveryMode bool   // Enable UDP server discovery (listens for any MAVLink source)
	DiscoveryHost string // Host to bind for discovery (e.g., "0.0.0.0")
	DiscoveryPort string // Port to listen on for discovery (e.g., "14550")

	// TCP Server Support
	EnableTCPServer bool   // Enable TCP server for legacy MAVLink tools
	TCPServerHost   string // TCP server bind host
	TCPServerPort   string // TCP server port

	// MAVLink Node Configuration (mavp2p-aligned)
	MAVLinkSystemID         uint8         // System ID for outgoing frames (default: 255 = GCS)
	MAVLinkComponentID      uint8         // Component ID for outgoing frames (default: 190 = MISSIONPLANNER)
	MAVLinkHeartbeatDisable bool          // Disable heartbeat emission
	MAVLinkHeartbeatPeriod  time.Duration // Heartbeat period (default: 5s)
	MAVLinkStreamReqEnable  bool          // Enable automatic stream requests (for Ardupilot)
	MAVLinkStreamReqFreq    int           // Stream request frequency in Hz (default: 4)
	MAVLinkIdleTimeout      time.Duration // Idle timeout for channels (default: 60s)

	// MAVLink UDP Server Inputs (optional - for network-based MAVLink sources)
	MAVLinkUDPHost  string
	MAVLinkUDPPorts []string // Support multiple UDP ports (e.g., ["14580", "14581", "14582"])

	// MAVLink UDP Client Forwards (optional - for forwarding to GCS, gimbal, etc.)
	UDPForwards []string // e.g., ["127.0.0.1:14550", "127.0.0.1:14551"]
	TCPForwards []string // e.g., ["127.0.0.1:14553", "192.168.1.200:5760"]

	// NATS Configuration
	NATSUrl         string
	NATSAuthToken   string // Auth token for NATS connection (optional)
	NATSStreamName  string
	NATSRootSubject string // Root subject for hierarchical naming (e.g., "constellation")
	NATSKVBucket    string
	NATSMaxAge      time.Duration
	NATSMaxMsgs     int64
	NATSMaxBytes    int64
	NATSStorageType string

	// Multi-Entity Configuration (loaded from overwatch.yaml or env)
	DefaultOrgID       string                   // Default org_id if not mapped
	DefaultEntityID    string                   // Default entity_id prefix if not mapped
	EntityMappings     map[string]EntityMapping // system_id -> entity assignment
	StreamMappings     map[string]StreamMapping // Legacy port -> org/entity mapping
	KVKeyPrefix        string                   // KV key prefix (e.g., "entity")
	KVUpdateTimeout    time.Duration            // Timeout for KV updates
	EnableKVPublishing bool                     // Enable KV publishing to CONSTELLATION_GLOBAL_STATE

	// Overwatch Configuration (loaded from overwatch.yaml)
	Overwatch *OverwatchConfig // Full overwatch configuration

	// Performance Configuration
	EnableAsyncPublishing bool // Enable high-performance async publishing (now uses Bridge)
	AsyncBatchSize        int  // Batch size for async publishing
	BatchTimeoutMS        int  // Batch timeout in milliseconds
	MaxConcurrentEntities int  // Maximum concurrent entities
	RoutingWorkers        int  // Number of routing workers
	ConnectionPoolSize    int  // Connection pool size

	// Monitoring Configuration
	EnableMetrics bool   // Enable HTTP metrics endpoint
	MetricsPort   string // Port for metrics endpoint (default: 9090)

	// Application Configuration
	LogLevel          string
	MessageBufferSize int
}

// EntityMapping defines entity assignment for a specific system_id.
type EntityMapping struct {
	EntityID    string `json:"entity_id" yaml:"entity_id"`
	Description string `json:"description,omitempty" yaml:"description,omitempty"`
}

// StreamMapping defines the org/entity assignment for a MAVLink stream (legacy).
type StreamMapping struct {
	Port     string `json:"port" yaml:"port"`
	OrgID    string `json:"org_id" yaml:"org_id"`
	EntityID string `json:"entity_id" yaml:"entity_id"`
}

// OverwatchConfig represents the full overwatch.yaml configuration structure.
// This aligns with the PRD v1.0-beta specification for entity-first configuration.
type OverwatchConfig struct {
	Organization OrganizationConfig `yaml:"organization"`
	Entities     []EntityConfig     `yaml:"entities"` // Array of entities with endpoints
	Defaults     EntityDefaultsConfig    `yaml:"defaults"`
	Routing      RoutingConfig           `yaml:"routing"`
	Registration RegistrationConfig      `yaml:"registration"`

	// Internal lookup maps (built after loading)
	entitiesBySystemID    map[uint8]*EntityConfig
	entitiesByComponentID map[uint8]*EntityConfig
	entitiesByEntityID    map[string]*EntityConfig
	entitiesByName        map[string]*EntityConfig
}

// OrganizationConfig holds organization-level settings.
type OrganizationConfig struct {
	ID          string `yaml:"id"`
	OrgID       string `yaml:"org_id"` // Alternative field name per PRD
	Name        string `yaml:"name"`
	OrgType     string `yaml:"org_type"`
	Description string `yaml:"description"`
	APIEndpoint string `yaml:"api_endpoint"`
	APIKey      string `yaml:"api_key"`
}

// GetOrgID returns the organization ID, preferring org_id over id.
func (o *OrganizationConfig) GetOrgID() string {
	if o.OrgID != "" {
		return o.OrgID
	}
	return o.ID
}

// EntityConfig defines a single entity's configuration with physical endpoint.
// This is the PRD v1.0-beta aligned structure with endpoint-per-entity.
type EntityConfig struct {
	EntityID        string            `yaml:"entity_id"`
	SystemID        uint8             `yaml:"system_id"`
	ComponentID     uint8             `yaml:"component_id"`
	Type            string            `yaml:"type"`
	Name            string            `yaml:"name"`
	Status          string            `yaml:"status"`
	Priority        string            `yaml:"priority"`
	Description     string            `yaml:"description"`
	Capabilities    []string          `yaml:"capabilities"`
	Endpoint        EndpointConfig    `yaml:"endpoint"`
	Metadata        map[string]any    `yaml:"metadata"`
	CommandPriority int               `yaml:"command_priority"`
	AuthorizedCmds  []string          `yaml:"authorized_commands"`
}

// EndpointConfig represents the physical MAVLink connection for an entity.
type EndpointConfig struct {
	Type        string `yaml:"type"`         // serial, udp, tcp, discovery
	Device      string `yaml:"device"`       // Serial device path (e.g., /dev/ttyUSB0)
	Baud        int    `yaml:"baud"`         // Serial baud rate
	FlowControl bool   `yaml:"flow_control"` // Serial flow control
	Address     string `yaml:"address"`      // UDP/TCP address
	Port        int    `yaml:"port"`         // UDP/TCP port
	Mode        string `yaml:"mode"`         // client or server
}

// EntityDefaultsConfig holds default settings for unmapped entities.
type EntityDefaultsConfig struct {
	OrgID        string   `yaml:"org_id"`
	EntityPrefix string   `yaml:"entity_prefix"`
	Type         string   `yaml:"type"`
	Status       string   `yaml:"status"`
	Priority     string   `yaml:"priority"`
	AutoRegister bool     `yaml:"auto_register"`
	Capabilities []string `yaml:"capabilities"`
}

// RoutingConfig holds NATS routing configuration.
type RoutingConfig struct {
	RootSubject            string `yaml:"root_subject"`
	CommandSubjectSuffix   string `yaml:"command_subject_suffix"`
	TelemetrySubjectSuffix string `yaml:"telemetry_subject_suffix"`
}

// RegistrationConfig holds Constellation API registration settings.
type RegistrationConfig struct {
	Enabled                  bool   `yaml:"enabled"`
	VerifyOnStartup          bool   `yaml:"verify_on_startup"`
	AutoRegisterUnknown      bool   `yaml:"auto_register_unknown"`
	HeartbeatInterval        string `yaml:"heartbeat_interval"`
	FingerprintOnHeartbeat   bool   `yaml:"fingerprint_on_heartbeat"`
	IDStrategy               string `yaml:"id_strategy"` // system_id, endpoint, or overwatch_managed
	StateTTL                 string `yaml:"state_ttl"`
}

// BuildLookupMaps creates internal lookup maps for fast entity resolution.
func (c *OverwatchConfig) BuildLookupMaps() {
	if c == nil {
		return
	}
	c.entitiesBySystemID = make(map[uint8]*EntityConfig)
	c.entitiesByComponentID = make(map[uint8]*EntityConfig)
	c.entitiesByEntityID = make(map[string]*EntityConfig)
	c.entitiesByName = make(map[string]*EntityConfig)

	for i := range c.Entities {
		entity := &c.Entities[i]
		if entity.SystemID > 0 {
			c.entitiesBySystemID[entity.SystemID] = entity
		}
		if entity.ComponentID > 0 {
			c.entitiesByComponentID[entity.ComponentID] = entity
		}
		if entity.EntityID != "" {
			c.entitiesByEntityID[entity.EntityID] = entity
		}
		if entity.Name != "" {
			c.entitiesByName[entity.Name] = entity
		}
	}
}

// GetEntityBySystemID returns the entity configuration for a given system_id.
func (c *OverwatchConfig) GetEntityBySystemID(systemID uint8) *EntityConfig {
	if c == nil || c.entitiesBySystemID == nil {
		return nil
	}
	return c.entitiesBySystemID[systemID]
}

// GetEntityByComponentID returns the entity configuration for a given component_id.
// This allows mapping by the device's component_id (often correlates to hostname/device type).
func (c *OverwatchConfig) GetEntityByComponentID(componentID uint8) *EntityConfig {
	if c == nil || c.entitiesByComponentID == nil {
		return nil
	}
	return c.entitiesByComponentID[componentID]
}

// GetEntityByName returns the entity configuration for a given name.
func (c *OverwatchConfig) GetEntityByName(name string) *EntityConfig {
	if c == nil || c.entitiesByName == nil {
		return nil
	}
	return c.entitiesByName[name]
}

// GetEntityByID returns the entity configuration for a given entity_id.
func (c *OverwatchConfig) GetEntityByID(entityID string) *EntityConfig {
	if c == nil || c.entitiesByEntityID == nil {
		return nil
	}
	return c.entitiesByEntityID[entityID]
}

// GetEntitiesWithEndpoints returns all entities that have endpoint configurations.
func (c *OverwatchConfig) GetEntitiesWithEndpoints() []*EntityConfig {
	if c == nil {
		return nil
	}
	var result []*EntityConfig
	for i := range c.Entities {
		if c.Entities[i].Endpoint.Type != "" {
			result = append(result, &c.Entities[i])
		}
	}
	return result
}

// ToEntityMappings converts OverwatchConfig entities to the legacy EntityMappings format.
func (c *OverwatchConfig) ToEntityMappings() map[string]EntityMapping {
	mappings := make(map[string]EntityMapping)
	if c == nil {
		return mappings
	}
	for _, entity := range c.Entities {
		if entity.SystemID > 0 {
			systemIDStr := strconv.Itoa(int(entity.SystemID))
			mappings[systemIDStr] = EntityMapping{
				EntityID:    entity.EntityID,
				Description: entity.Description,
			}
		}
	}
	return mappings
}
