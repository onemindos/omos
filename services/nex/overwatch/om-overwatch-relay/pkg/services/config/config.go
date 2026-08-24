package config

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"gopkg.in/yaml.v3"
)

// LoadConfig loads all configuration from environment variables and overwatch.yaml.
func LoadConfig() (*Config, error) {
	// Load .env file if it exists (ignore error if file doesn't exist)
	_ = godotenv.Load()

	// Parse multi-serial device configuration
	var serialDevices []string
	var serialBauds []int
	var serialEntityIDs []string

	// Multi-device support
	devicesStr := getEnv("SERIAL_DEVICES", "")
	if devicesStr != "" {
		for _, device := range strings.Split(devicesStr, ",") {
			serialDevices = append(serialDevices, strings.TrimSpace(device))
		}
	}

	baudsStr := getEnv("SERIAL_BAUDS", "57600")
	for _, baudStr := range strings.Split(baudsStr, ",") {
		baud := getEnvInt("", 57600) // Parse individual baud rate
		if _, err := fmt.Sscanf(strings.TrimSpace(baudStr), "%d", &baud); err == nil {
			serialBauds = append(serialBauds, baud)
		} else {
			serialBauds = append(serialBauds, 57600) // Default
		}
	}

	entitiesStr := getEnv("SERIAL_ENTITY_IDS", "")
	if entitiesStr != "" {
		for _, entity := range strings.Split(entitiesStr, ",") {
			serialEntityIDs = append(serialEntityIDs, strings.TrimSpace(entity))
		}
	}

	// Legacy single serial device support
	legacySerialDevice := getEnv("SERIAL_DEVICE", "")
	legacySerialBaud := getEnvInt("SERIAL_BAUD", 57600)

	// Auto-detect serial devices if enabled
	serialAutoDetect := getEnv("SERIAL_AUTO_DETECT", "true") == "true"
	if serialAutoDetect && len(serialDevices) == 0 && legacySerialDevice == "" {
		if detected := autoDetectSerialDevice(); detected != "" {
			legacySerialDevice = detected
		}
	}

	// Parse UDP server ports (comma-separated)
	portsStr := getEnv("MAVLINK_UDP_PORTS", "")
	var ports []string
	if portsStr != "" {
		ports = strings.Split(portsStr, ",")
		for i := range ports {
			ports[i] = strings.TrimSpace(ports[i])
		}
	}

	// Parse UDP forwards (comma-separated list of address:port)
	forwardsStr := getEnv("UDP_FORWARDS", "")
	var udpForwards []string
	if forwardsStr != "" {
		udpForwards = strings.Split(forwardsStr, ",")
		for i := range udpForwards {
			udpForwards[i] = strings.TrimSpace(udpForwards[i])
		}
	}

	// Parse TCP forwards (comma-separated list of address:port)
	tcpForwardsStr := getEnv("TCP_FORWARDS", "")
	var tcpForwards []string
	if tcpForwardsStr != "" {
		tcpForwards = strings.Split(tcpForwardsStr, ",")
		for i := range tcpForwards {
			tcpForwards[i] = strings.TrimSpace(tcpForwards[i])
		}
	}

	// Load overwatch.yaml configuration
	overwatchConfig, err := loadOverwatchConfig()
	if err != nil {
		log.Printf("Warning: Could not load overwatch.yaml: %v", err)
	} else if overwatchConfig != nil {
		// Build lookup maps for fast entity resolution
		overwatchConfig.BuildLookupMaps()
		log.Printf("Loaded %d entities from overwatch.yaml", len(overwatchConfig.Entities))
	}

	// Parse entity mappings from env (JSON format) or use overwatch.yaml
	entityMappings := make(map[string]EntityMapping)
	entityMappingsStr := getEnv("ENTITY_MAPPINGS", "")
	if entityMappingsStr != "" {
		// Environment variable takes precedence
		var rawMappings map[string]EntityMapping
		if err := json.Unmarshal([]byte(entityMappingsStr), &rawMappings); err != nil {
			log.Printf("Warning: Failed to parse ENTITY_MAPPINGS, using defaults: %v", err)
		} else {
			entityMappings = rawMappings
		}
	} else if overwatchConfig != nil {
		// Use overwatch.yaml entity mappings (derived from entities array)
		entityMappings = overwatchConfig.ToEntityMappings()
		log.Printf("Built %d entity mappings from overwatch.yaml entities", len(entityMappings))
	}

	// Parse legacy stream mappings (JSON format)
	mappingsStr := getEnv("STREAM_MAPPINGS", "")
	streamMappings := make(map[string]StreamMapping)
	if mappingsStr != "" {
		var rawMappings map[string]struct {
			OrgID    string `json:"org_id"`
			EntityID string `json:"entity_id"`
		}
		if err := json.Unmarshal([]byte(mappingsStr), &rawMappings); err != nil {
			log.Printf("Warning: Failed to parse STREAM_MAPPINGS, using defaults: %v", err)
		} else {
			for port, mapping := range rawMappings {
				streamMappings[port] = StreamMapping{
					Port:     port,
					OrgID:    mapping.OrgID,
					EntityID: mapping.EntityID,
				}
			}
		}
	}

	// Determine default org ID (overwatch.yaml takes precedence if no env var)
	defaultOrgID := getEnv("CONSTELLATION_ORG_ID", "")
	if defaultOrgID == "" && overwatchConfig != nil && overwatchConfig.Organization.ID != "" {
		defaultOrgID = overwatchConfig.Organization.ID
	}
	if defaultOrgID == "" {
		defaultOrgID = "org001"
	}

	// Determine default entity prefix (overwatch.yaml takes precedence if no env var)
	defaultEntityID := getEnv("CONSTELLATION_ENTITY_ID", "")
	if defaultEntityID == "" && overwatchConfig != nil && overwatchConfig.Defaults.EntityPrefix != "" {
		defaultEntityID = overwatchConfig.Defaults.EntityPrefix
	}
	if defaultEntityID == "" {
		defaultEntityID = "device"
	}

	// Determine NATS root subject (overwatch.yaml routing takes precedence if no env var)
	natsRootSubject := getEnv("NATS_ROOT_SUBJECT", "")
	if natsRootSubject == "" && overwatchConfig != nil && overwatchConfig.Routing.RootSubject != "" {
		natsRootSubject = overwatchConfig.Routing.RootSubject
	}
	if natsRootSubject == "" {
		natsRootSubject = "constellation"
	}

	cfg := &Config{
		// Multi-Serial device configuration
		SerialDevices:    serialDevices,
		SerialBauds:      serialBauds,
		SerialEntityIDs:  serialEntityIDs,
		SerialAutoDetect: serialAutoDetect,

		// Legacy single serial device (deprecated but maintained for compatibility)
		SerialDevice: legacySerialDevice,
		SerialBaud:   legacySerialBaud,

		// MAVLink Discovery mode defaults
		DiscoveryMode: getEnv("MAVLINK_DISCOVERY_MODE", "false") == "true",
		DiscoveryHost: getEnv("MAVLINK_DISCOVERY_HOST", "0.0.0.0"),
		DiscoveryPort: getEnv("MAVLINK_DISCOVERY_PORT", "14550"),

		// TCP Server configuration
		EnableTCPServer: getEnv("ENABLE_TCP_SERVER", "false") == "true",
		TCPServerHost:   getEnv("TCP_SERVER_HOST", "0.0.0.0"),
		TCPServerPort:   getEnv("TCP_SERVER_PORT", "5760"),

		// MAVLink Node Configuration (mavp2p-aligned)
		MAVLinkSystemID:         uint8(getEnvInt("MAVLINK_SYSTEM_ID", 255)),         // GCS default
		MAVLinkComponentID:      uint8(getEnvInt("MAVLINK_COMPONENT_ID", 190)),      // MISSIONPLANNER default
		MAVLinkHeartbeatDisable: getEnv("MAVLINK_HEARTBEAT_DISABLE", "true") == "true",
		MAVLinkHeartbeatPeriod:  getEnvDuration("MAVLINK_HEARTBEAT_PERIOD", 5*time.Second),
		MAVLinkStreamReqEnable:  getEnv("MAVLINK_STREAM_REQ_ENABLE", "false") == "true",
		MAVLinkStreamReqFreq:    getEnvInt("MAVLINK_STREAM_REQ_FREQ", 4),
		MAVLinkIdleTimeout:      getEnvDuration("MAVLINK_IDLE_TIMEOUT", 60*time.Second),

		// MAVLink UDP server defaults
		MAVLinkUDPHost:  getEnv("MAVLINK_UDP_HOST", "127.0.0.1"),
		MAVLinkUDPPorts: ports,

		// MAVLink forwards
		UDPForwards: udpForwards,
		TCPForwards: tcpForwards,

		// NATS defaults — OneMind sovereign bus
		NATSUrl:         getEnv("NATS_URL", "nats://127.0.0.1:4222"),
		NATSAuthToken:   getEnv("NATS_AUTH_TOKEN", ""),
		NATSStreamName:  getEnv("NATS_STREAM_NAME", "TEL"),
		NATSRootSubject: natsRootSubject,
		NATSKVBucket:    getEnv("NATS_KV_BUCKET", "OVERWATCH_STATE"),
		NATSMaxAge:      getEnvDuration("NATS_MAX_AGE", 24*time.Hour),
		NATSMaxMsgs:     getEnvInt64("NATS_MAX_MSGS", 100000),
		NATSMaxBytes:    getEnvInt64("NATS_MAX_BYTES", 1024*1024*1024), // 1GB
		NATSStorageType: getEnv("NATS_STORAGE_TYPE", "file"),

		// Multi-Entity configuration
		DefaultOrgID:       defaultOrgID,
		DefaultEntityID:    defaultEntityID,
		EntityMappings:     entityMappings,
		StreamMappings:     streamMappings, // Legacy support
		KVKeyPrefix:        getEnv("KV_KEY_PREFIX", "entity"),
		KVUpdateTimeout:    getEnvDuration("KV_UPDATE_TIMEOUT", 2*time.Second),
		// KV publishing off by default — om-overwatch-state (the spine) owns KV fusion.
		EnableKVPublishing: getEnv("ENABLE_KV_PUBLISHING", "false") == "true",

		// Overwatch configuration
		Overwatch: overwatchConfig,

		// Performance configuration
		EnableAsyncPublishing: getEnv("ENABLE_ASYNC_PUBLISHING", "true") == "true",
		AsyncBatchSize:        getEnvInt("ASYNC_BATCH_SIZE", 100),
		BatchTimeoutMS:        getEnvInt("BATCH_TIMEOUT_MS", 50),
		MaxConcurrentEntities: getEnvInt("MAX_CONCURRENT_ENTITIES", 20),
		RoutingWorkers:        getEnvInt("ROUTING_WORKERS", 4),
		ConnectionPoolSize:    getEnvInt("CONNECTION_POOL_SIZE", 20),

		// Monitoring configuration
		EnableMetrics: getEnv("ENABLE_METRICS", "true") == "true",
		MetricsPort:   getEnv("METRICS_PORT", "9091"),

		// Application defaults
		LogLevel:          getEnv("LOG_LEVEL", "info"),
		MessageBufferSize: getEnvInt("MESSAGE_BUFFER_SIZE", 1000),
	}

	return cfg, nil
}

// loadOverwatchConfig loads the overwatch.yaml configuration file.
func loadOverwatchConfig() (*OverwatchConfig, error) {
	// Check for overwatch.yaml in current directory and common locations
	configPaths := []string{
		"overwatch.yaml",
		"config/overwatch.yaml",
		"/etc/mavlink2constellation/overwatch.yaml",
	}

	// Also check OVERWATCH_CONFIG env var
	if customPath := getEnv("OVERWATCH_CONFIG", ""); customPath != "" {
		configPaths = append([]string{customPath}, configPaths...)
	}

	var configPath string
	for _, path := range configPaths {
		if _, err := os.Stat(path); err == nil {
			configPath = path
			break
		}
	}

	if configPath == "" {
		return nil, fmt.Errorf("overwatch.yaml not found in any search path")
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read %s: %w", configPath, err)
	}

	var config OverwatchConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse %s: %w", configPath, err)
	}

	log.Printf("Loaded overwatch configuration from %s", configPath)
	return &config, nil
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	if value, exists := os.LookupEnv(key); exists {
		var intVal int
		if _, err := fmt.Sscanf(value, "%d", &intVal); err == nil {
			return intVal
		}
	}
	return fallback
}

func getEnvInt64(key string, fallback int64) int64 {
	if value, exists := os.LookupEnv(key); exists {
		var intVal int64
		if _, err := fmt.Sscanf(value, "%d", &intVal); err == nil {
			return intVal
		}
	}
	return fallback
}

func getEnvDuration(key string, fallback time.Duration) time.Duration {
	if value, exists := os.LookupEnv(key); exists {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return fallback
}

// autoDetectSerialDevice attempts to find a MAVLink-capable serial device
// Returns the first available serial device path, or empty string if none found
// Supports both macOS and Linux device naming conventions
func autoDetectSerialDevice() string {
	// Define search patterns for different platforms
	patterns := []string{
		// macOS patterns (FTDI chips used in RFD900, 3DR radios, etc.)
		"/dev/cu.usbserial-*",      // FTDI USB-to-Serial (call-out device - preferred for macOS)
		"/dev/tty.usbserial-*",     // FTDI USB-to-Serial (dial-in device)
		"/dev/cu.SLAB_USBtoUART*",  // Silicon Labs CP210x (call-out)
		"/dev/tty.SLAB_USBtoUART*", // Silicon Labs CP210x (dial-in)
		"/dev/cu.usbmodem*",        // Pixhawk/ArduPilot native USB
		"/dev/tty.usbmodem*",       // Pixhawk/ArduPilot native USB

		// Linux patterns
		"/dev/ttyUSB*",            // FTDI, Prolific, etc.
		"/dev/ttyACM*",            // CDC ACM devices (Pixhawk native USB)
		"/dev/serial/by-id/usb-*", // Persistent by USB ID (most reliable on Linux)
	}

	log.Println("Searching for serial devices...")

	// Collect all matches
	var allMatches []string
	for _, pattern := range patterns {
		matches, err := filepath.Glob(pattern)
		if err != nil {
			continue
		}
		allMatches = append(allMatches, matches...)
	}

	if len(allMatches) == 0 {
		log.Println("⚠ No serial devices found")
		return ""
	}

	// Sort to get consistent results
	sort.Strings(allMatches)

	// On macOS, prefer /dev/cu.* over /dev/tty.* (recommended for serial communications)
	// cu.* devices are "call-out" devices that don't wait for carrier detect
	for _, match := range allMatches {
		base := filepath.Base(match)
		if len(base) >= 3 && base[:3] == "cu." {
			log.Printf("✓ Auto-detected serial device: %s (macOS call-out device)", match)
			return match
		}
	}

	// Fallback to first match
	log.Printf("✓ Auto-detected serial device: %s", allMatches[0])
	if len(allMatches) > 1 {
		log.Printf("  Note: %d other serial devices available", len(allMatches)-1)
	}

	return allMatches[0]
}
