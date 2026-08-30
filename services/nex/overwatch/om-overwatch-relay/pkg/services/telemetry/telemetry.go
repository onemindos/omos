package telemetry

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"math"
	"time"

	"github.com/onemindos/om-overwatch-relay/pkg/services/shared"
	"github.com/bluenviron/gomavlib/v3"
	"github.com/bluenviron/gomavlib/v3/pkg/dialects/common"
	"github.com/bluenviron/gomavlib/v3/pkg/message"
)

// GetMessageTypeName returns a human-readable name for a MAVLink message ID.
// Uses the shared MAVLinkMessageNames map.
func GetMessageTypeName(msgID uint32) string {
	if name, ok := shared.MAVLinkMessageNames[msgID]; ok {
		return name
	}
	return fmt.Sprintf("MSG_%d", msgID)
}

// decodeRawMessage manually decodes a MessageRaw payload into structured fields
// This is needed because gomavlib sometimes returns MessageRaw instead of typed messages
func decodeRawMessage(rawMsg *message.MessageRaw, messageID uint32) (map[string]interface{}, error) {
	payload := rawMsg.Payload
	buf := bytes.NewReader(payload)
	data := make(map[string]interface{})

	// Helper function to read from buffer with error handling
	readUint8 := func() uint8 {
		var v uint8
		binary.Read(buf, binary.LittleEndian, &v)
		return v
	}
	readUint16 := func() uint16 {
		var v uint16
		binary.Read(buf, binary.LittleEndian, &v)
		return v
	}
	readUint32 := func() uint32 {
		var v uint32
		binary.Read(buf, binary.LittleEndian, &v)
		return v
	}
	readUint64 := func() uint64 {
		var v uint64
		binary.Read(buf, binary.LittleEndian, &v)
		return v
	}
	readInt16 := func() int16 {
		var v int16
		binary.Read(buf, binary.LittleEndian, &v)
		return v
	}
	readInt32 := func() int32 {
		var v int32
		binary.Read(buf, binary.LittleEndian, &v)
		return v
	}
	readFloat32 := func() float32 {
		var v float32
		binary.Read(buf, binary.LittleEndian, &v)
		return v
	}

	// Decode based on message ID
	switch messageID {
	case 0: // HEARTBEAT
		data["custom_mode"] = readUint32()
		data["type"] = readUint8()
		data["autopilot"] = readUint8()
		data["base_mode"] = readUint8()
		data["system_status"] = readUint8()
		data["mavlink_version"] = readUint8()

	case 1: // SYS_STATUS
		data["onboard_control_sensors_present"] = readUint32()
		data["onboard_control_sensors_enabled"] = readUint32()
		data["onboard_control_sensors_health"] = readUint32()
		data["load"] = readUint16()
		voltage := readUint16()
		data["voltage_battery"] = float64(voltage) / 1000.0 // mV to V
		current := int16(readUint16())
		data["current_battery"] = float64(current) / 100.0 // cA to A
		data["battery_remaining"] = int8(readUint8())
		data["drop_rate_comm"] = readUint16()
		data["errors_comm"] = readUint16()

	case 24: // GPS_RAW_INT
		data["time_usec"] = readUint64()
		data["fix_type"] = readUint8()
		lat := readInt32()
		data["lat"] = float64(lat) / 1e7 // degE7 to degrees
		lon := readInt32()
		data["lon"] = float64(lon) / 1e7 // degE7 to degrees
		alt := readInt32()
		data["alt"] = float64(alt) / 1000.0 // mm to m
		data["eph"] = readUint16()
		data["epv"] = readUint16()
		data["vel"] = readUint16()
		data["cog"] = readUint16()
		data["satellites_visible"] = readUint8()

	case 29: // SCALED_PRESSURE
		data["time_boot_ms"] = readUint32()
		data["press_abs"] = readFloat32()
		data["press_diff"] = readFloat32()
		data["temperature"] = readInt16()

	case 30: // ATTITUDE
		data["time_boot_ms"] = readUint32()
		data["roll"] = readFloat32()
		data["pitch"] = readFloat32()
		data["yaw"] = readFloat32()
		data["rollspeed"] = readFloat32()
		data["pitchspeed"] = readFloat32()
		data["yawspeed"] = readFloat32()

	case 31: // ATTITUDE_QUATERNION
		data["time_boot_ms"] = readUint32()
		data["q1"] = readFloat32()
		data["q2"] = readFloat32()
		data["q3"] = readFloat32()
		data["q4"] = readFloat32()
		data["rollspeed"] = readFloat32()
		data["pitchspeed"] = readFloat32()
		data["yawspeed"] = readFloat32()

	case 32: // LOCAL_POSITION_NED
		data["time_boot_ms"] = readUint32()
		data["x"] = readFloat32()
		data["y"] = readFloat32()
		data["z"] = readFloat32()
		data["vx"] = readFloat32()
		data["vy"] = readFloat32()
		data["vz"] = readFloat32()

	case 33: // GLOBAL_POSITION_INT
		data["time_boot_ms"] = readUint32()
		lat := readInt32()
		data["lat"] = float64(lat) / 1e7
		lon := readInt32()
		data["lon"] = float64(lon) / 1e7
		alt := readInt32()
		data["alt"] = float64(alt) / 1000.0
		relAlt := readInt32()
		data["relative_alt"] = float64(relAlt) / 1000.0
		data["vx"] = readInt16()
		data["vy"] = readInt16()
		data["vz"] = readInt16()
		data["hdg"] = readUint16()

	case 36: // SERVO_OUTPUT_RAW
		data["time_usec"] = readUint32()
		data["servo1_raw"] = readUint16()
		data["servo2_raw"] = readUint16()
		data["servo3_raw"] = readUint16()
		data["servo4_raw"] = readUint16()
		data["servo5_raw"] = readUint16()
		data["servo6_raw"] = readUint16()
		data["servo7_raw"] = readUint16()
		data["servo8_raw"] = readUint16()
		data["port"] = readUint8()

	case 42: // MISSION_CURRENT
		data["seq"] = readUint16()

	case 74: // VFR_HUD
		data["airspeed"] = readFloat32()
		data["groundspeed"] = readFloat32()
		data["heading"] = readInt16()
		data["throttle"] = readUint16()
		data["alt"] = readFloat32()
		data["climb"] = readFloat32()

	case 83: // ATTITUDE_TARGET
		data["time_boot_ms"] = readUint32()
		data["type_mask"] = readUint8()
		q := make([]float32, 4)
		for i := 0; i < 4; i++ {
			q[i] = readFloat32()
		}
		data["q"] = q
		data["body_roll_rate"] = readFloat32()
		data["body_pitch_rate"] = readFloat32()
		data["body_yaw_rate"] = readFloat32()
		data["thrust"] = readFloat32()

	case 85: // POSITION_TARGET_LOCAL_NED
		data["time_boot_ms"] = readUint32()
		data["type_mask"] = readUint16()
		data["x"] = readFloat32()
		data["y"] = readFloat32()
		data["z"] = readFloat32()
		data["vx"] = readFloat32()
		data["vy"] = readFloat32()
		data["vz"] = readFloat32()
		data["afx"] = readFloat32()
		data["afy"] = readFloat32()
		data["afz"] = readFloat32()
		data["yaw"] = readFloat32()
		data["yaw_rate"] = readFloat32()
		data["coordinate_frame"] = readUint8()

	case 141: // ALTITUDE
		data["time_usec"] = readUint64()
		data["altitude_monotonic"] = readFloat32()
		data["altitude_amsl"] = readFloat32()
		data["altitude_local"] = readFloat32()
		data["altitude_relative"] = readFloat32()
		data["altitude_terrain"] = readFloat32()
		data["bottom_clearance"] = readFloat32()

	case 147: // BATTERY_STATUS
		data["id"] = readUint8()
		data["battery_function"] = readUint8()
		data["type"] = readUint8()
		data["temperature"] = readInt16()
		voltages := make([]uint16, 10)
		for i := 0; i < 10; i++ {
			voltages[i] = readUint16()
		}
		data["voltages"] = voltages
		data["current_battery"] = readInt16()
		data["current_consumed"] = readInt32()
		data["energy_consumed"] = readInt32()
		data["battery_remaining"] = int8(readUint8())

	case 245: // EXTENDED_SYS_STATE
		data["vtol_state"] = readUint8()
		data["landed_state"] = readUint8()

	default:
		// Unknown message - return empty to fall back to generic handling
		return nil, fmt.Errorf("unknown message ID: %d", messageID)
	}

	// Filter out NaN and Inf values that might cause JSON marshaling issues
	for key, val := range data {
		if f, ok := val.(float32); ok {
			if math.IsNaN(float64(f)) || math.IsInf(float64(f), 0) {
				data[key] = nil
			}
		} else if f, ok := val.(float64); ok {
			if math.IsNaN(f) || math.IsInf(f, 0) {
				data[key] = nil
			}
		}
	}

	return data, nil
}

// ParseMavlinkMessage converts a MAVLink message to our TelemetryMessage format
// Uses specific type handlers for known messages, with JSON fallback for unknown types
func ParseMavlinkMessage(frame *gomavlib.EventFrame) (*TelemetryMessage, error) {
	msg := &TelemetryMessage{
		Timestamp:   time.Now().UTC(),
		MessageID:   frame.Message().GetID(),
		SystemID:    frame.SystemID(),
		ComponentID: frame.ComponentID(),
		Data:        make(map[string]interface{}),
	}

	// Get human-readable message type name
	msg.MessageType = GetMessageTypeName(msg.MessageID)

	// Debug: Log the actual type being received
	frameMsg := frame.Message()
	msgType := fmt.Sprintf("%T", frameMsg)

	// HANDLE RAW MESSAGES: gomavlib returns MessageRaw when messages aren't pre-parsed
	// We need to manually decode the binary payload into structured fields
	if rawMsg, ok := frameMsg.(*message.MessageRaw); ok {
		decodedData, err := decodeRawMessage(rawMsg, msg.MessageID)
		if err == nil && len(decodedData) > 0 {
			msg.Data = decodedData
			return msg, nil
		}
		// If decoding failed, fall through to generic handling below
	}

	// SPECIFIC MESSAGE TYPE HANDLERS
	// Use type assertions to extract specific message fields (for pre-parsed messages)
	switch m := frameMsg.(type) {
	case *common.MessageHeartbeat:
		msg.Data = map[string]interface{}{
			"type":            m.Type,
			"autopilot":       m.Autopilot,
			"base_mode":       m.BaseMode,
			"custom_mode":     m.CustomMode,
			"system_status":   m.SystemStatus,
			"mavlink_version": m.MavlinkVersion,
		}

	case *common.MessageSysStatus:
		msg.Data = map[string]interface{}{
			"onboard_control_sensors_present": m.OnboardControlSensorsPresent,
			"onboard_control_sensors_enabled": m.OnboardControlSensorsEnabled,
			"onboard_control_sensors_health":  m.OnboardControlSensorsHealth,
			"load":                            m.Load,
			"voltage_battery":                 float64(m.VoltageBattery) / 1000.0, // mV to V
			"current_battery":                 float64(m.CurrentBattery) / 100.0,  // cA to A
			"battery_remaining":               m.BatteryRemaining,
			"drop_rate_comm":                  m.DropRateComm,
			"errors_comm":                     m.ErrorsComm,
		}

	case *common.MessageGpsRawInt:
		msg.Data = map[string]interface{}{
			"time_usec":          m.TimeUsec,
			"fix_type":           m.FixType,
			"lat":                float64(m.Lat) / 1e7,    // degE7 to degrees
			"lon":                float64(m.Lon) / 1e7,    // degE7 to degrees
			"alt":                float64(m.Alt) / 1000.0, // mm to m
			"eph":                m.Eph,
			"epv":                m.Epv,
			"vel":                m.Vel,
			"cog":                m.Cog,
			"satellites_visible": m.SatellitesVisible,
		}

	case *common.MessageAttitude:
		msg.Data = map[string]interface{}{
			"time_boot_ms": m.TimeBootMs,
			"roll":         m.Roll,
			"pitch":        m.Pitch,
			"yaw":          m.Yaw,
			"rollspeed":    m.Rollspeed,
			"pitchspeed":   m.Pitchspeed,
			"yawspeed":     m.Yawspeed,
		}

	case *common.MessageAttitudeQuaternion:
		msg.Data = map[string]interface{}{
			"time_boot_ms": m.TimeBootMs,
			"q1":           m.Q1,
			"q2":           m.Q2,
			"q3":           m.Q3,
			"q4":           m.Q4,
			"rollspeed":    m.Rollspeed,
			"pitchspeed":   m.Pitchspeed,
			"yawspeed":     m.Yawspeed,
		}

	case *common.MessageLocalPositionNed:
		msg.Data = map[string]interface{}{
			"time_boot_ms": m.TimeBootMs,
			"x":            m.X,
			"y":            m.Y,
			"z":            m.Z,
			"vx":           m.Vx,
			"vy":           m.Vy,
			"vz":           m.Vz,
		}

	case *common.MessageGlobalPositionInt:
		msg.Data = map[string]interface{}{
			"time_boot_ms": m.TimeBootMs,
			"lat":          float64(m.Lat) / 1e7,            // degE7 to degrees
			"lon":          float64(m.Lon) / 1e7,            // degE7 to degrees
			"alt":          float64(m.Alt) / 1000.0,         // mm to m
			"relative_alt": float64(m.RelativeAlt) / 1000.0, // mm to m
			"vx":           m.Vx,
			"vy":           m.Vy,
			"vz":           m.Vz,
			"hdg":          m.Hdg,
		}

	case *common.MessageRcChannels:
		msg.Data = map[string]interface{}{
			"time_boot_ms": m.TimeBootMs,
			"chan1_raw":    m.Chan1Raw,
			"chan2_raw":    m.Chan2Raw,
			"chan3_raw":    m.Chan3Raw,
			"chan4_raw":    m.Chan4Raw,
			"chan5_raw":    m.Chan5Raw,
			"chan6_raw":    m.Chan6Raw,
			"chan7_raw":    m.Chan7Raw,
			"chan8_raw":    m.Chan8Raw,
			"chan9_raw":    m.Chan9Raw,
			"chan10_raw":   m.Chan10Raw,
			"chan11_raw":   m.Chan11Raw,
			"chan12_raw":   m.Chan12Raw,
			"chan13_raw":   m.Chan13Raw,
			"chan14_raw":   m.Chan14Raw,
			"chan15_raw":   m.Chan15Raw,
			"chan16_raw":   m.Chan16Raw,
			"chan17_raw":   m.Chan17Raw,
			"chan18_raw":   m.Chan18Raw,
			"chancount":    m.Chancount,
			"rssi":         m.Rssi,
		}

	case *common.MessageVfrHud:
		msg.Data = map[string]interface{}{
			"airspeed":    m.Airspeed,
			"groundspeed": m.Groundspeed,
			"heading":     m.Heading,
			"throttle":    m.Throttle,
			"alt":         m.Alt,
			"climb":       m.Climb,
		}

	case *common.MessageBatteryStatus:
		msg.Data = map[string]interface{}{
			"id":                m.Id,
			"battery_function":  m.BatteryFunction,
			"type":              m.Type,
			"temperature":       m.Temperature,
			"voltages":          m.Voltages,
			"current_battery":   m.CurrentBattery,
			"current_consumed":  m.CurrentConsumed,
			"energy_consumed":   m.EnergyConsumed,
			"battery_remaining": m.BatteryRemaining,
			"time_remaining":    m.TimeRemaining,
			"charge_state":      m.ChargeState,
		}

	case *common.MessageScaledPressure:
		msg.Data = map[string]interface{}{
			"time_boot_ms": m.TimeBootMs,
			"press_abs":    m.PressAbs,
			"press_diff":   m.PressDiff,
			"temperature":  m.Temperature,
		}

	case *common.MessageServoOutputRaw:
		msg.Data = map[string]interface{}{
			"time_usec":  m.TimeUsec,
			"port":       m.Port,
			"servo1_raw": m.Servo1Raw,
			"servo2_raw": m.Servo2Raw,
			"servo3_raw": m.Servo3Raw,
			"servo4_raw": m.Servo4Raw,
			"servo5_raw": m.Servo5Raw,
			"servo6_raw": m.Servo6Raw,
			"servo7_raw": m.Servo7Raw,
			"servo8_raw": m.Servo8Raw,
		}

	case *common.MessageMissionCurrent:
		msg.Data = map[string]interface{}{
			"seq": m.Seq,
		}

	case *common.MessageAttitudeTarget:
		msg.Data = map[string]interface{}{
			"time_boot_ms":    m.TimeBootMs,
			"type_mask":       m.TypeMask,
			"q":               []float32{m.Q[0], m.Q[1], m.Q[2], m.Q[3]},
			"body_roll_rate":  m.BodyRollRate,
			"body_pitch_rate": m.BodyPitchRate,
			"body_yaw_rate":   m.BodyYawRate,
			"thrust":          m.Thrust,
		}

	case *common.MessagePositionTargetLocalNed:
		msg.Data = map[string]interface{}{
			"time_boot_ms":     m.TimeBootMs,
			"coordinate_frame": m.CoordinateFrame,
			"type_mask":        m.TypeMask,
			"x":                m.X,
			"y":                m.Y,
			"z":                m.Z,
			"vx":               m.Vx,
			"vy":               m.Vy,
			"vz":               m.Vz,
			"afx":              m.Afx,
			"afy":              m.Afy,
			"afz":              m.Afz,
			"yaw":              m.Yaw,
			"yaw_rate":         m.YawRate,
		}

	case *common.MessageAltitude:
		msg.Data = map[string]interface{}{
			"time_usec":          m.TimeUsec,
			"altitude_monotonic": m.AltitudeMonotonic,
			"altitude_amsl":      m.AltitudeAmsl,
			"altitude_local":     m.AltitudeLocal,
			"altitude_relative":  m.AltitudeRelative,
			"altitude_terrain":   m.AltitudeTerrain,
			"bottom_clearance":   m.BottomClearance,
		}

	case *common.MessageExtendedSysState:
		msg.Data = map[string]interface{}{
			"vtol_state":   m.VtolState,
			"landed_state": m.LandedState,
		}

	default:
		// GENERIC FALLBACK: Use JSON marshaling for unknown message types
		// This should rarely be hit now that we handle MessageRaw directly
		_ = msgType // silence unused variable warning

		rawJSON, err := json.Marshal(frame.Message())
		if err != nil {
			msg.Data = map[string]interface{}{
				"raw_message": fmt.Sprintf("%+v", frame.Message()),
				"error":       err.Error(),
			}
			return msg, nil
		}

		if err := json.Unmarshal(rawJSON, &msg.Data); err != nil {
			msg.Data = map[string]interface{}{
				"raw_json": string(rawJSON),
			}
			return msg, nil
		}

		// Clean up byte array fields for unknown messages
		cleanByteArrayFields(msg.Data)
	}

	return msg, nil
}

// cleanByteArrayFields removes or converts byte array fields from parsed data
// Byte arrays in MAVLink messages (like payload fields) marshal as large integer arrays
// which are not useful in JSON and take up space. This function cleans them up.
func cleanByteArrayFields(data map[string]interface{}) {
	// Fields that commonly contain byte buffers or large arrays
	fieldsToClean := []string{
		"payload", "data", "value", "param_value",
		"file_buf", "image_data", "chunk_data",
		"covariance", // Covariance matrices (can be 36-45 floats)
	}

	for _, field := range fieldsToClean {
		if val, exists := data[field]; exists {
			// Check if it's an array (byte arrays unmarshal as []interface{})
			if arr, ok := val.([]interface{}); ok {
				// If array is large (>20 elements), replace with metadata
				// This catches covariance matrices and large buffers
				if len(arr) > 20 {
					data[field+"_length"] = len(arr)
					data[field+"_type"] = "array_omitted"
					delete(data, field)
				} else if len(arr) > 8 {
					// Medium arrays (8-20 elements): just provide length
					data[field+"_length"] = len(arr)
					delete(data, field)
				}
			}
		}
	}

	// Also clean up UID fields and other byte arrays that might be all zeros
	for key, val := range data {
		if arr, ok := val.([]interface{}); ok {
			// Skip if we already handled it above
			skip := false
			for _, field := range fieldsToClean {
				if key == field {
					skip = true
					break
				}
			}
			if skip {
				continue
			}

			// Check if it's a byte array with all zeros (common in unused fields)
			if len(arr) > 0 && len(arr) <= 20 {
				allZeros := true
				for _, v := range arr {
					if num, ok := v.(float64); ok {
						if num != 0 {
							allZeros = false
							break
						}
					}
				}
				// Remove zero-filled byte arrays to reduce JSON clutter
				if allZeros {
					delete(data, key)
				}
			}
		}
	}
}
