package mavlink

import (
	"fmt"
	"log"
	"time"

	"github.com/onemindos/om-overwatch-relay/pkg/services/config"
	"github.com/bluenviron/gomavlib/v3"
	"github.com/bluenviron/gomavlib/v3/pkg/dialects/common"
)

// CreateNode creates a new MAVLink router node based on configuration.
// It first attempts to build endpoints from overwatch.yaml entity definitions,
// then falls back to legacy environment variable configuration.
func CreateNode(cfg *config.Config) (*gomavlib.Node, error) {
	// Build endpoint configuration for unified router
	endpoints := make([]gomavlib.EndpointConf, 0)

	// Priority 1: Build endpoints from overwatch.yaml entities (PRD v1.0-beta)
	if cfg.Overwatch != nil {
		entityEndpoints := buildEndpointsFromEntities(cfg.Overwatch)
		endpoints = append(endpoints, entityEndpoints...)
	}

	// Priority 2: Legacy serial device configuration
	if cfg.SerialDevice != "" {
		endpoints = append(endpoints, gomavlib.EndpointSerial{
			Device: cfg.SerialDevice,
			Baud:   cfg.SerialBaud,
		})
		log.Printf("Added serial endpoint (legacy): %s @ %d baud", cfg.SerialDevice, cfg.SerialBaud)
	}

	// Priority 3: Multi-serial device configuration
	for i, device := range cfg.SerialDevices {
		baud := 57600 // Default
		if i < len(cfg.SerialBauds) {
			baud = cfg.SerialBauds[i]
		}
		endpoints = append(endpoints, gomavlib.EndpointSerial{
			Device: device,
			Baud:   baud,
		})
		log.Printf("Added serial endpoint (multi): %s @ %d baud", device, baud)
	}

	// Priority 4: Discovery mode endpoint (recommended for network deployments)
	if cfg.DiscoveryMode {
		discoveryAddress := fmt.Sprintf("%s:%s", cfg.DiscoveryHost, cfg.DiscoveryPort)
		endpoints = append(endpoints, gomavlib.EndpointUDPServer{Address: discoveryAddress})
		log.Printf("Discovery mode ENABLED: listening on %s for all MAVLink sources", discoveryAddress)
	} else {
		// Add UDP server endpoints (for receiving from network sources - static port mode)
		for _, port := range cfg.MAVLinkUDPPorts {
			address := fmt.Sprintf("%s:%s", cfg.MAVLinkUDPHost, port)
			endpoints = append(endpoints, gomavlib.EndpointUDPServer{Address: address})
			log.Printf("Added UDP server endpoint: %s", address)
		}
	}

	// Priority 5: TCP server endpoint
	if cfg.EnableTCPServer {
		tcpAddress := fmt.Sprintf("%s:%s", cfg.TCPServerHost, cfg.TCPServerPort)
		endpoints = append(endpoints, gomavlib.EndpointTCPServer{Address: tcpAddress})
		log.Printf("Added TCP server endpoint: %s", tcpAddress)
	}

	// Priority 6: UDP client forwards (for GCS, gimbal, etc.)
	for _, forward := range cfg.UDPForwards {
		endpoints = append(endpoints, gomavlib.EndpointUDPClient{Address: forward})
		log.Printf("Added UDP forward endpoint: %s", forward)
	}

	// Priority 7: TCP client forwards
	for _, forward := range cfg.TCPForwards {
		endpoints = append(endpoints, gomavlib.EndpointTCPClient{Address: forward})
		log.Printf("Added TCP forward endpoint: %s", forward)
	}

	// Validate we have at least one endpoint
	if len(endpoints) == 0 {
		return nil, fmt.Errorf("no MAVLink endpoints configured. Configure entities in overwatch.yaml or set SERIAL_DEVICE, MAVLINK_UDP_PORTS, or UDP_FORWARDS")
	}

	// Configure heartbeat period (0 means use default)
	heartbeatPeriod := cfg.MAVLinkHeartbeatPeriod
	if heartbeatPeriod == 0 {
		heartbeatPeriod = 5 * time.Second
	}

	// Configure idle timeout (0 means use default)
	idleTimeout := cfg.MAVLinkIdleTimeout
	if idleTimeout == 0 {
		idleTimeout = 60 * time.Second
	}

	// Create unified MAVLink router node (mavp2p-aligned configuration)
	node, err := gomavlib.NewNode(gomavlib.NodeConf{
		Endpoints:   endpoints,
		Dialect:     common.Dialect,
		OutVersion:  gomavlib.V2,
		OutSystemID: cfg.MAVLinkSystemID, // Configurable (default: 255 = GCS)
		OutComponentID:          cfg.MAVLinkComponentID,      // Configurable (default: 190 = MISSIONPLANNER)
		HeartbeatDisable:        cfg.MAVLinkHeartbeatDisable, // Disable heartbeat emission
		HeartbeatPeriod:         heartbeatPeriod,             // Period between heartbeats
		StreamRequestEnable:     cfg.MAVLinkStreamReqEnable,  // Enable stream requests for Ardupilot
		StreamRequestFrequency:  cfg.MAVLinkStreamReqFreq,    // Stream request frequency (Hz)
		IdleTimeout:             idleTimeout,                 // Cleanup inactive channels
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create MAVLink router node: %w", err)
	}

	log.Printf("MAVLink router node created with %d endpoint(s)", len(endpoints))
	return node, nil
}

// buildEndpointsFromEntities creates gomavlib endpoint configurations from
// overwatch.yaml entity definitions. This implements the PRD v1.0-beta
// endpoint-per-entity pattern.
func buildEndpointsFromEntities(ow *config.OverwatchConfig) []gomavlib.EndpointConf {
	if ow == nil {
		return nil
	}

	var endpoints []gomavlib.EndpointConf
	entitiesWithEndpoints := ow.GetEntitiesWithEndpoints()

	for _, entity := range entitiesWithEndpoints {
		ep := entity.Endpoint

		switch ep.Type {
		case "serial":
			if ep.Device != "" {
				baud := ep.Baud
				if baud == 0 {
					baud = 57600 // Default baud
				}
				endpoints = append(endpoints, gomavlib.EndpointSerial{
					Device: ep.Device,
					Baud:   baud,
				})
				log.Printf("Added serial endpoint for %s: %s @ %d baud", entity.EntityID, ep.Device, baud)
			}

		case "udp":
			if ep.Address != "" && ep.Port > 0 {
				address := fmt.Sprintf("%s:%d", ep.Address, ep.Port)
				if ep.Mode == "server" {
					endpoints = append(endpoints, gomavlib.EndpointUDPServer{Address: address})
					log.Printf("Added UDP server endpoint for %s: %s", entity.EntityID, address)
				} else {
					// Default to client mode for forwards (GCS, gimbal, etc.)
					endpoints = append(endpoints, gomavlib.EndpointUDPClient{Address: address})
					log.Printf("Added UDP client endpoint for %s: %s", entity.EntityID, address)
				}
			}

		case "tcp":
			if ep.Address != "" && ep.Port > 0 {
				address := fmt.Sprintf("%s:%d", ep.Address, ep.Port)
				if ep.Mode == "server" {
					endpoints = append(endpoints, gomavlib.EndpointTCPServer{Address: address})
					log.Printf("Added TCP server endpoint for %s: %s", entity.EntityID, address)
				} else {
					endpoints = append(endpoints, gomavlib.EndpointTCPClient{Address: address})
					log.Printf("Added TCP client endpoint for %s: %s", entity.EntityID, address)
				}
			}

		case "discovery":
			if ep.Port > 0 {
				address := fmt.Sprintf("%s:%d", ep.Address, ep.Port)
				if address == ":0" || ep.Address == "" {
					address = fmt.Sprintf("0.0.0.0:%d", ep.Port)
				}
				endpoints = append(endpoints, gomavlib.EndpointUDPServer{Address: address})
				log.Printf("Added discovery endpoint for %s: %s", entity.EntityID, address)
			}
		}
	}

	return endpoints
}
