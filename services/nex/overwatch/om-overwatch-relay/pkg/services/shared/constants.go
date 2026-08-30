// Package shared provides common constants, types, and utilities used across services.
package shared

import "time"

// Default timeouts used across services
const (
	DefaultKVTimeout      = 2 * time.Second
	DefaultPublishTimeout = 5 * time.Second
	DefaultConnectTimeout = 10 * time.Second
)

// Default buffer sizes
const (
	DefaultMessageBufferSize = 1000
	DefaultBatchSize         = 100
)

// MAVLink message priority categories
const (
	PriorityCritical = iota // Immediate publish - commands, safety
	PriorityHigh            // High priority, low batching - navigation, heartbeat
	PriorityMedium          // Medium priority, moderate batching - telemetry
	PriorityLow             // Low priority, high batching - sensors, simulation
)

// MAVLinkMessageNames maps MAVLink message IDs to human-readable names.
// Source: MAVLink common.xml dialect
var MAVLinkMessageNames = map[uint32]string{
	0:   "HEARTBEAT",
	1:   "SYS_STATUS",
	2:   "SYSTEM_TIME",
	4:   "PING",
	5:   "CHANGE_OPERATOR_CONTROL",
	6:   "CHANGE_OPERATOR_CONTROL_ACK",
	7:   "AUTH_KEY",
	11:  "SET_MODE",
	20:  "PARAM_REQUEST_READ",
	21:  "PARAM_REQUEST_LIST",
	22:  "PARAM_VALUE",
	23:  "PARAM_SET",
	24:  "GPS_RAW_INT",
	25:  "GPS_STATUS",
	26:  "SCALED_IMU",
	27:  "RAW_IMU",
	28:  "RAW_PRESSURE",
	29:  "SCALED_PRESSURE",
	30:  "ATTITUDE",
	31:  "ATTITUDE_QUATERNION",
	32:  "LOCAL_POSITION_NED",
	33:  "GLOBAL_POSITION_INT",
	34:  "RC_CHANNELS_SCALED",
	35:  "RC_CHANNELS_RAW",
	36:  "SERVO_OUTPUT_RAW",
	37:  "MISSION_REQUEST_PARTIAL_LIST",
	38:  "MISSION_WRITE_PARTIAL_LIST",
	39:  "MISSION_ITEM",
	40:  "MISSION_REQUEST",
	41:  "MISSION_SET_CURRENT",
	42:  "MISSION_CURRENT",
	43:  "MISSION_REQUEST_LIST",
	44:  "MISSION_COUNT",
	45:  "MISSION_CLEAR_ALL",
	46:  "MISSION_ITEM_REACHED",
	47:  "MISSION_ACK",
	48:  "SET_GPS_GLOBAL_ORIGIN",
	49:  "GPS_GLOBAL_ORIGIN",
	50:  "PARAM_MAP_RC",
	51:  "MISSION_REQUEST_INT",
	54:  "SAFETY_SET_ALLOWED_AREA",
	55:  "SAFETY_ALLOWED_AREA",
	61:  "ATTITUDE_QUATERNION_COV",
	62:  "NAV_CONTROLLER_OUTPUT",
	63:  "GLOBAL_POSITION_INT_COV",
	64:  "LOCAL_POSITION_NED_COV",
	65:  "RC_CHANNELS",
	66:  "REQUEST_DATA_STREAM",
	67:  "DATA_STREAM",
	69:  "MANUAL_CONTROL",
	70:  "RC_CHANNELS_OVERRIDE",
	73:  "MISSION_ITEM_INT",
	74:  "VFR_HUD",
	75:  "COMMAND_INT",
	76:  "COMMAND_LONG",
	77:  "COMMAND_ACK",
	81:  "MANUAL_SETPOINT",
	82:  "SET_ATTITUDE_TARGET",
	83:  "ATTITUDE_TARGET",
	84:  "SET_POSITION_TARGET_LOCAL_NED",
	85:  "POSITION_TARGET_LOCAL_NED",
	86:  "SET_POSITION_TARGET_GLOBAL_INT",
	87:  "POSITION_TARGET_GLOBAL_INT",
	89:  "LOCAL_POSITION_NED_SYSTEM_GLOBAL_OFFSET",
	90:  "HIL_STATE",
	91:  "HIL_CONTROLS",
	92:  "HIL_RC_INPUTS_RAW",
	100: "OPTICAL_FLOW",
	101: "GLOBAL_VISION_POSITION_ESTIMATE",
	102: "VISION_POSITION_ESTIMATE",
	103: "VISION_SPEED_ESTIMATE",
	104: "VICON_POSITION_ESTIMATE",
	105: "HIGHRES_IMU",
	106: "OPTICAL_FLOW_RAD",
	107: "HIL_SENSOR",
	108: "SIM_STATE",
	109: "RADIO_STATUS",
	110: "FILE_TRANSFER_PROTOCOL",
	111: "TIMESYNC",
	112: "CAMERA_TRIGGER",
	113: "HIL_GPS",
	114: "HIL_OPTICAL_FLOW",
	115: "HIL_STATE_QUATERNION",
	116: "SCALED_IMU2",
	117: "LOG_REQUEST_LIST",
	118: "LOG_ENTRY",
	119: "LOG_REQUEST_DATA",
	120: "LOG_DATA",
	121: "LOG_ERASE",
	122: "LOG_REQUEST_END",
	123: "GPS_INJECT_DATA",
	124: "GPS2_RAW",
	125: "POWER_STATUS",
	126: "SERIAL_CONTROL",
	127: "GPS_RTK",
	128: "GPS2_RTK",
	129: "SCALED_IMU3",
	130: "DATA_TRANSMISSION_HANDSHAKE",
	131: "ENCAPSULATED_DATA",
	132: "DISTANCE_SENSOR",
	133: "TERRAIN_REQUEST",
	134: "TERRAIN_DATA",
	135: "TERRAIN_CHECK",
	136: "TERRAIN_REPORT",
	137: "SCALED_PRESSURE2",
	138: "ATT_POS_MOCAP",
	139: "SET_ACTUATOR_CONTROL_TARGET",
	140: "ACTUATOR_CONTROL_TARGET",
	141: "ALTITUDE",
	142: "RESOURCE_REQUEST",
	143: "SCALED_PRESSURE3",
	144: "FOLLOW_TARGET",
	146: "CONTROL_SYSTEM_STATE",
	147: "BATTERY_STATUS",
	148: "AUTOPILOT_VERSION",
	149: "LANDING_TARGET",
	230: "ESTIMATOR_STATUS",
	231: "WIND_COV",
	232: "GPS_INPUT",
	233: "GPS_RTCM_DATA",
	234: "HIGH_LATENCY",
	235: "HIGH_LATENCY2",
	241: "VIBRATION",
	242: "HOME_POSITION",
	243: "SET_HOME_POSITION",
	244: "MESSAGE_INTERVAL",
	245: "EXTENDED_SYS_STATE",
	246: "ADSB_VEHICLE",
	247: "COLLISION",
	248: "V2_EXTENSION",
	249: "MEMORY_VECT",
	250: "DEBUG_VECT",
	251: "NAMED_VALUE_FLOAT",
	252: "NAMED_VALUE_INT",
	253: "STATUSTEXT",
	254: "DEBUG",
	256: "SETUP_SIGNING",
	257: "BUTTON_CHANGE",
	258: "PLAY_TUNE",
	259: "CAMERA_INFORMATION",
	260: "CAMERA_SETTINGS",
	261: "STORAGE_INFORMATION",
	262: "CAMERA_CAPTURE_STATUS",
	263: "CAMERA_IMAGE_CAPTURED",
	264: "FLIGHT_INFORMATION",
	265: "MOUNT_ORIENTATION",
	266: "LOGGING_DATA",
	267: "LOGGING_DATA_ACKED",
	268: "LOGGING_ACK",
	269: "VIDEO_STREAM_INFORMATION",
	270: "VIDEO_STREAM_STATUS",
	299: "WIFI_CONFIG_AP",
	300: "PROTOCOL_VERSION",
	310: "UAVCAN_NODE_STATUS",
	311: "UAVCAN_NODE_INFO",
	320: "PARAM_EXT_REQUEST_READ",
	321: "PARAM_EXT_REQUEST_LIST",
	322: "PARAM_EXT_VALUE",
	323: "PARAM_EXT_SET",
	324: "PARAM_EXT_ACK",
	330: "OBSTACLE_DISTANCE",
	331: "ODOMETRY",
}

// MessagePriorityMap defines the publishing priority for MAVLink message types.
// Used by the async publisher for batching decisions.
var MessagePriorityMap = map[string]int{
	// Critical messages - immediate publish
	"COMMAND_LONG":                   PriorityCritical,
	"COMMAND_INT":                    PriorityCritical,
	"SET_MODE":                       PriorityCritical,
	"ARM_DISARM_COMPONENT":           PriorityCritical,
	"SET_POSITION_TARGET_GLOBAL_INT": PriorityCritical,
	"SET_ATTITUDE_TARGET":            PriorityCritical,

	// High priority - low latency, small batches
	"HEARTBEAT":           PriorityHigh,
	"MISSION_ITEM_INT":    PriorityHigh,
	"MISSION_REQUEST_INT": PriorityHigh,
	"GPS_RAW_INT":         PriorityHigh,
	"GLOBAL_POSITION_INT": PriorityHigh,

	// Medium priority - moderate batching
	"SYS_STATUS":     PriorityMedium,
	"ATTITUDE":       PriorityMedium,
	"VFR_HUD":        PriorityMedium,
	"BATTERY_STATUS": PriorityMedium,
	"RC_CHANNELS":    PriorityMedium,

	// Low priority - high batching
	"SCALED_IMU":  PriorityLow,
	"SCALED_IMU2": PriorityLow,
	"SCALED_IMU3": PriorityLow,
	"HIGHRES_IMU": PriorityLow,
	"RAW_IMU":     PriorityLow,
	"LOG_ENTRY":   PriorityLow,
	"LOG_DATA":    PriorityLow,
}

// KV state message types that trigger state updates
var KVStateMessageTypes = []string{
	"HEARTBEAT",
	"GPS_RAW_INT",
	"GLOBAL_POSITION_INT",
	"ATTITUDE",
	"BATTERY_STATUS",
	"SYS_STATUS",
}

// IsKVStateMessage returns true if the message type should trigger a KV state update
func IsKVStateMessage(msgType string) bool {
	for _, t := range KVStateMessageTypes {
		if t == msgType {
			return true
		}
	}
	return false
}
