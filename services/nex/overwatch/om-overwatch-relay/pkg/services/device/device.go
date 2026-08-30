package device

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/onemindos/om-overwatch-relay/pkg/services/config"
	"github.com/onemindos/om-overwatch-relay/pkg/services/telemetry"
	"github.com/bluenviron/gomavlib/v3"
	"github.com/nats-io/nats.go/jetstream"
)

// FingerprintTimeout is the duration after which a device will be fingerprinted
// even without receiving a HEARTBEAT message. This handles devices like gimbals,
// mavp2p relays, or other components that don't emit HEARTBEAT.
const FingerprintTimeout = 30 * time.Second

// GetOrCreateDeviceFingerprint retrieves or creates a device fingerprint
// This function loads persistent fingerprints from KV store to ensure consistency across restarts
func GetOrCreateDeviceFingerprint(devices map[string]*DeviceFingerprint, msg *telemetry.TelemetryMessage, channel *gomavlib.Channel, cfg *config.Config, kv jetstream.KeyValue) *DeviceFingerprint {
	deviceKey := fmt.Sprintf("%d.%d", msg.SystemID, msg.ComponentID)

	// Check if device already exists in memory
	if device, exists := devices[deviceKey]; exists {
		device.LastSeen = msg.Timestamp

		// If not yet fingerprinted and this is a HEARTBEAT, fingerprint now
		if !device.Fingerprinted && msg.MessageType == "HEARTBEAT" {
			fingerprintDevice(device, msg, cfg, kv)
		}

		// Fallback: fingerprint after timeout even without HEARTBEAT
		// This handles gimbals, mavp2p relays, and other non-HEARTBEAT devices
		if !device.Fingerprinted && time.Since(device.FirstSeen) > FingerprintTimeout {
			log.Printf("[Device %d/%d] Fingerprinting without HEARTBEAT (timeout after %v)",
				device.SystemID, device.ComponentID, FingerprintTimeout)
			fingerprintDeviceWithoutHeartbeat(device, msg, cfg, kv)
		}

		return device
	}

	// Create new device fingerprint
	remoteAddr := "unknown"
	if channel != nil {
		remoteAddr = channel.String()
	}

	device := &DeviceFingerprint{
		SystemID:      msg.SystemID,
		ComponentID:   msg.ComponentID,
		Channel:       channel,
		RemoteAddr:    remoteAddr,
		Fingerprinted: false,
		FirstSeen:     msg.Timestamp,
		LastSeen:      msg.Timestamp,
	}

	// Try to load existing fingerprint from KV store (global state)
	if kv != nil {
		kvKey := fmt.Sprintf("device.%d.%d", msg.SystemID, msg.ComponentID)
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		entry, err := kv.Get(ctx, kvKey)
		cancel()

		if err == nil {
			// Existing device found in global state - restore fingerprint
			var state DeviceState
			if err := json.Unmarshal(entry.Value(), &state); err == nil && state.Fingerprinted {
				device.OrgID = state.OrgID
				device.EntityID = state.EntityID
				device.Subject = state.Subject
				device.Fingerprinted = true
				device.FirstSeen = state.FirstSeen

				log.Printf("[Device %d/%d] Loaded existing fingerprint from global state: %s (subject: %s)",
					device.SystemID, device.ComponentID, device.EntityID, device.Subject)

				devices[deviceKey] = device
				return device
			}
		}
	}

	// Fingerprint on first HEARTBEAT
	if msg.MessageType == "HEARTBEAT" {
		fingerprintDevice(device, msg, cfg, kv)
	}

	devices[deviceKey] = device
	return device
}

// omSubject builds the OneMind telemetry subject for a device.
// Scheme: tel.drone.{entity_id}.mavlink  — token[2] is the entity id, which
// om-overwatch-state (the fusion spine) reads to fuse friendly UAV telemetry.
// No org_id: OneMind is a single sovereign operator, not a multi-tenant SaaS.
func omSubject(entityID string) string {
	return fmt.Sprintf("tel.drone.%s.mavlink", entityID)
}

// resolveEntityID picks a stable entity id for a device. overwatch.yaml is
// optional; when absent we auto-generate from system_id.
// Priority: overwatch.yaml (component_id) → overwatch.yaml (system_id) →
// ENTITY_MAPPINGS env → auto "drone-{sysid}".
func resolveEntityID(msg *telemetry.TelemetryMessage, cfg *config.Config) (entityID, name string) {
	if cfg.Overwatch != nil {
		if ec := cfg.Overwatch.GetEntityByComponentID(msg.ComponentID); ec != nil {
			return ec.EntityID, ec.Name
		}
		if ec := cfg.Overwatch.GetEntityBySystemID(msg.SystemID); ec != nil {
			return ec.EntityID, ec.Name
		}
	}
	if m, ok := cfg.EntityMappings[fmt.Sprintf("%d", msg.SystemID)]; ok {
		return m.EntityID, ""
	}
	prefix := cfg.DefaultEntityID
	if prefix == "" || prefix == "device" {
		prefix = "drone"
	}
	return fmt.Sprintf("%s-%d", prefix, msg.SystemID), ""
}

// fingerprintDevice assigns a stable entity id + OneMind subject to a device.
// Each MAVLink stream is uniquely identified by system_id.component_id, so
// multiple drones forwarded through one UDP port stay separate.
func fingerprintDevice(device *DeviceFingerprint, msg *telemetry.TelemetryMessage, cfg *config.Config, kv jetstream.KeyValue) {
	entityID, name := resolveEntityID(msg, cfg)
	device.EntityID = entityID
	device.Name = name
	device.Subject = omSubject(entityID)
	device.Fingerprinted = true
	log.Printf("[Device %d/%d] Fingerprinted: %s (subject: %s)",
		device.SystemID, device.ComponentID, device.EntityID, device.Subject)
}

// fingerprintDeviceWithoutHeartbeat fingerprints devices that don't send
// HEARTBEAT (gimbals, mavp2p relays, etc.) after FingerprintTimeout.
func fingerprintDeviceWithoutHeartbeat(device *DeviceFingerprint, msg *telemetry.TelemetryMessage, cfg *config.Config, kv jetstream.KeyValue) {
	entityID, name := resolveEntityID(msg, cfg)
	device.EntityID = entityID
	device.Name = name
	device.Subject = omSubject(entityID)
	device.Fingerprinted = true
	log.Printf("[Device %d/%d] Fingerprinted (no-heartbeat): %s (subject: %s)",
		device.SystemID, device.ComponentID, device.EntityID, device.Subject)
}

// =============================================================================
// MODULAR KV SIGNAL TREE PUBLISHING
// =============================================================================
// Publishes MAVLink telemetry to a hierarchical KV structure:
//   {entity_id}                    - Core entity metadata (updated every 5s)
//   {entity_id}.mavlink.heartbeat  - Heartbeat/mode data
//   {entity_id}.mavlink.position   - GPS/position data
//   {entity_id}.mavlink.attitude   - Attitude/orientation data
//   {entity_id}.mavlink.power      - Battery/power data
//   {entity_id}.mavlink.flight     - VFR HUD flight data
//
// Benefits:
// - High-frequency updates don't overwrite slower data
// - Selective rendering - only dirty subsystems update
// - Service isolation - each publisher owns its namespace
// - Bandwidth efficiency - smaller individual updates
// =============================================================================

// PublishMAVLinkSignal publishes a MAVLink message to the appropriate KV signal key.
// This implements the modular signal tree pattern for efficient real-time updates.
func PublishMAVLinkSignal(ctx context.Context, kv jetstream.KeyValue, device *DeviceFingerprint, msg *telemetry.TelemetryMessage) error {
	if kv == nil || !device.Fingerprinted {
		return nil
	}

	entityID := device.EntityID

	// Route message to appropriate signal key
	switch msg.MessageType {
	case "HEARTBEAT":
		return publishHeartbeatSignal(ctx, kv, entityID, msg)

	case "GPS_RAW_INT":
		return publishPositionSignal(ctx, kv, entityID, msg, "gps")

	case "GLOBAL_POSITION_INT":
		return publishPositionSignal(ctx, kv, entityID, msg, "global")

	case "ATTITUDE", "ATTITUDE_QUATERNION":
		return publishAttitudeSignal(ctx, kv, entityID, msg)

	case "BATTERY_STATUS":
		return publishPowerSignal(ctx, kv, entityID, msg, "battery")

	case "SYS_STATUS":
		return publishPowerSignal(ctx, kv, entityID, msg, "sys")

	case "VFR_HUD":
		return publishFlightSignal(ctx, kv, entityID, msg)

	default:
		// Other message types don't update KV (stream-only)
		return nil
	}
}

// publishHeartbeatSignal writes heartbeat data to {entity_id}.mavlink.heartbeat
func publishHeartbeatSignal(ctx context.Context, kv jetstream.KeyValue, entityID string, msg *telemetry.TelemetryMessage) error {
	key := fmt.Sprintf("%s.mavlink.heartbeat", entityID)

	// Decode armed state from base_mode bitmask
	armed := false
	if baseMode, ok := msg.Data["base_mode"].(uint8); ok {
		armed = (baseMode & 128) != 0 // MAV_MODE_FLAG_SAFETY_ARMED = 128
	}

	data := map[string]interface{}{
		"timestamp": msg.Timestamp,
		"armed":     armed, // Derived: easy boolean check
	}

	addIfNotNil(data, "type", msg.Data["type"])                    // MAV_TYPE enum
	addIfNotNil(data, "autopilot", msg.Data["autopilot"])          // MAV_AUTOPILOT enum
	addIfNotNil(data, "base_mode", msg.Data["base_mode"])          // MAV_MODE_FLAG bitmask
	addIfNotNil(data, "custom_mode", msg.Data["custom_mode"])      // Autopilot-specific mode
	addIfNotNil(data, "system_status", msg.Data["system_status"])  // MAV_STATE enum
	addIfNotNil(data, "mavlink_version", msg.Data["mavlink_version"])

	return putKVData(ctx, kv, key, data)
}

// publishPositionSignal writes position data to {entity_id}.mavlink.position
func publishPositionSignal(ctx context.Context, kv jetstream.KeyValue, entityID string, msg *telemetry.TelemetryMessage, source string) error {
	key := fmt.Sprintf("%s.mavlink.position", entityID)

	data := map[string]interface{}{
		"timestamp": msg.Timestamp,
		"source":    source, // "gps" or "global"
	}

	addIfNotNil(data, "lat", msg.Data["lat"])
	addIfNotNil(data, "lon", msg.Data["lon"])
	addIfNotNil(data, "alt", msg.Data["alt"])

	// Add source-specific fields
	if source == "gps" {
		addIfNotNil(data, "fix_type", msg.Data["fix_type"])
		addIfNotNil(data, "satellites", msg.Data["satellites_visible"])
		addIfNotNil(data, "eph", msg.Data["eph"])
		addIfNotNil(data, "epv", msg.Data["epv"])
		addIfNotNil(data, "vel", msg.Data["vel"])
		addIfNotNil(data, "cog", msg.Data["cog"])
	} else {
		addIfNotNil(data, "relative_alt", msg.Data["relative_alt"])
		addIfNotNil(data, "vx", msg.Data["vx"])
		addIfNotNil(data, "vy", msg.Data["vy"])
		addIfNotNil(data, "vz", msg.Data["vz"])
		addIfNotNil(data, "hdg", msg.Data["hdg"])
	}

	return putKVData(ctx, kv, key, data)
}

// publishAttitudeSignal writes attitude data to {entity_id}.mavlink.attitude
// Handles both ATTITUDE and ATTITUDE_QUATERNION messages with different key structures
func publishAttitudeSignal(ctx context.Context, kv jetstream.KeyValue, entityID string, msg *telemetry.TelemetryMessage) error {
	key := fmt.Sprintf("%s.mavlink.attitude", entityID)

	data := map[string]interface{}{
		"timestamp": msg.Timestamp,
	}

	// ATTITUDE_QUATERNION has quaternion values but no roll/pitch/yaw euler angles
	// ATTITUDE has euler angles but no quaternion
	if msg.MessageType == "ATTITUDE_QUATERNION" {
		// Quaternion-based attitude (q1-q4 are the quaternion components)
		addIfNotNil(data, "q1", msg.Data["q1"])
		addIfNotNil(data, "q2", msg.Data["q2"])
		addIfNotNil(data, "q3", msg.Data["q3"])
		addIfNotNil(data, "q4", msg.Data["q4"])
		addIfNotNil(data, "rollspeed", msg.Data["rollspeed"])
		addIfNotNil(data, "pitchspeed", msg.Data["pitchspeed"])
		addIfNotNil(data, "yawspeed", msg.Data["yawspeed"])
		data["representation"] = "quaternion"
	} else {
		// Euler angle-based attitude
		addIfNotNil(data, "roll", msg.Data["roll"])
		addIfNotNil(data, "pitch", msg.Data["pitch"])
		addIfNotNil(data, "yaw", msg.Data["yaw"])
		addIfNotNil(data, "rollspeed", msg.Data["rollspeed"])
		addIfNotNil(data, "pitchspeed", msg.Data["pitchspeed"])
		addIfNotNil(data, "yawspeed", msg.Data["yawspeed"])
		data["representation"] = "euler"
	}

	return putKVData(ctx, kv, key, data)
}

// addIfNotNil adds a key-value pair to the map only if the value is not nil
func addIfNotNil(data map[string]interface{}, key string, value interface{}) {
	if value != nil {
		data[key] = value
	}
}

// publishPowerSignal writes power/battery data to {entity_id}.mavlink.power
func publishPowerSignal(ctx context.Context, kv jetstream.KeyValue, entityID string, msg *telemetry.TelemetryMessage, source string) error {
	key := fmt.Sprintf("%s.mavlink.power", entityID)

	data := map[string]interface{}{
		"timestamp": msg.Timestamp,
		"source":    source, // "battery" or "sys"
	}

	addIfNotNil(data, "battery_remaining", msg.Data["battery_remaining"])

	if source == "battery" {
		addIfNotNil(data, "id", msg.Data["id"])
		addIfNotNil(data, "temperature", msg.Data["temperature"])
		addIfNotNil(data, "voltages", msg.Data["voltages"])
		addIfNotNil(data, "current_battery", msg.Data["current_battery"])
		addIfNotNil(data, "current_consumed", msg.Data["current_consumed"])
	} else {
		addIfNotNil(data, "voltage_battery", msg.Data["voltage_battery"])
		addIfNotNil(data, "current_battery", msg.Data["current_battery"])
		addIfNotNil(data, "load", msg.Data["load"])
		addIfNotNil(data, "drop_rate_comm", msg.Data["drop_rate_comm"])
		addIfNotNil(data, "errors_comm", msg.Data["errors_comm"])
	}

	return putKVData(ctx, kv, key, data)
}

// publishFlightSignal writes VFR HUD data to {entity_id}.mavlink.flight
func publishFlightSignal(ctx context.Context, kv jetstream.KeyValue, entityID string, msg *telemetry.TelemetryMessage) error {
	key := fmt.Sprintf("%s.mavlink.flight", entityID)

	data := map[string]interface{}{
		"timestamp": msg.Timestamp,
	}

	addIfNotNil(data, "airspeed", msg.Data["airspeed"])
	addIfNotNil(data, "groundspeed", msg.Data["groundspeed"])
	addIfNotNil(data, "heading", msg.Data["heading"])
	addIfNotNil(data, "throttle", msg.Data["throttle"])
	addIfNotNil(data, "alt", msg.Data["alt"])
	addIfNotNil(data, "climb", msg.Data["climb"])

	return putKVData(ctx, kv, key, data)
}

// putKVData marshals data and puts it in the KV store
func putKVData(ctx context.Context, kv jetstream.KeyValue, key string, data map[string]interface{}) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal KV data for %s: %w", key, err)
	}

	_, err = kv.Put(ctx, key, jsonData)
	if err != nil {
		return fmt.Errorf("failed to put KV data for %s: %w", key, err)
	}

	return nil
}

// =============================================================================
// CORE ENTITY STATE (rate-limited updates)
// =============================================================================

// CoreEntityUpdateInterval controls how often the core entity record is updated
const CoreEntityUpdateInterval = 5 * time.Second

// UpdateCoreEntityState updates the core entity metadata record.
// This should be called periodically (every CoreEntityUpdateInterval) rather than on every message.
// Returns true if an update was performed, false if skipped due to rate limiting.
func UpdateCoreEntityState(ctx context.Context, kv jetstream.KeyValue, device *DeviceFingerprint, msg *telemetry.TelemetryMessage) (bool, error) {
	if kv == nil || !device.Fingerprinted {
		return false, nil
	}

	// Rate limit: only update core entity record every CoreEntityUpdateInterval
	device.Mu.Lock()
	if time.Since(device.LastCoreUpdate) < CoreEntityUpdateInterval {
		device.Mu.Unlock()
		return false, nil
	}
	device.LastCoreUpdate = time.Now()
	device.Mu.Unlock()

	// Core entity key (no suffix)
	deviceKey := device.EntityID

	// Get existing state or create new one
	var state DeviceState
	entry, err := kv.Get(ctx, deviceKey)
	if err == nil {
		if err := json.Unmarshal(entry.Value(), &state); err != nil {
			return false, fmt.Errorf("failed to unmarshal device state: %w", err)
		}
	} else {
		// New device - initialize core metadata
		state = DeviceState{
			SystemID:      msg.SystemID,
			ComponentID:   msg.ComponentID,
			OrgID:         device.OrgID,
			EntityID:      device.EntityID,
			Name:          device.Name,
			StreamPort:    device.RemoteAddr,
			Subject:       device.Subject,
			FirstSeen:     msg.Timestamp,
			Fingerprinted: device.Fingerprinted,
		}
	}

	// Update last seen timestamp
	state.LastSeen = msg.Timestamp

	// Marshal and put back in KV store
	stateJSON, err := json.Marshal(state)
	if err != nil {
		return false, fmt.Errorf("failed to marshal device state: %w", err)
	}

	_, err = kv.Put(ctx, deviceKey, stateJSON)
	if err != nil {
		return false, fmt.Errorf("failed to put device state: %w", err)
	}

	return true, nil
}

// =============================================================================
// LEGACY SUPPORT (deprecated - use PublishMAVLinkSignal instead)
// =============================================================================

// UpdateDeviceState updates the device state in the KV store (LEGACY - monolithic)
// Deprecated: Use PublishMAVLinkSignal for modular signal tree updates
func UpdateDeviceState(ctx context.Context, kv jetstream.KeyValue, device *DeviceFingerprint, msg *telemetry.TelemetryMessage) error {
	// Route to new modular signal publishing
	return PublishMAVLinkSignal(ctx, kv, device, msg)
}
