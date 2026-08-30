package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bluenviron/gomavlib/v3"
	"github.com/bluenviron/gomavlib/v3/pkg/dialects/ardupilotmega"
	"github.com/nats-io/nats.go"
)

// -----------------------------------------------------------------------
// Standard One Mind message envelope — same across ALL omos-bridge bridges
// -----------------------------------------------------------------------

type BridgeMessage struct {
	ID      string          `json:"id"`
	TS      int64           `json:"ts"`
	Source  string          `json:"source"`
	Entity  string          `json:"entity"`
	Domain  string          `json:"domain"`
	Payload json.RawMessage `json:"payload"`
}

// -----------------------------------------------------------------------
// Bridge-defined telemetry payloads (MAVLink message types → Go structs)
// -----------------------------------------------------------------------

type HeartbeatPayload struct {
	SysID          uint8  `json:"sys_id"`
	CompID         uint8  `json:"comp_id"`
	Type           uint32 `json:"type"`
	Autopilot      uint32 `json:"autopilot"`
	BaseMode       uint32 `json:"base_mode"`
	SystemStatus   uint32 `json:"system_status"`
	MavlinkVersion uint32 `json:"mavlink_version"`
}

type PositionPayload struct {
	SysID     uint8   `json:"sys_id"`
	CompID    uint8   `json:"comp_id"`
	TimeBootMs uint32 `json:"time_boot_ms"`
	Lat       float64 `json:"lat"`
	Lon       float64 `json:"lon"`
	Alt       float64 `json:"alt_m"`
	RelAlt    float64 `json:"relative_alt_m"`
	Vx        float32 `json:"vx_ms"`
	Vy        float32 `json:"vy_ms"`
	Vz        float32 `json:"vz_ms"`
	Hdg       float32 `json:"heading_deg"`
}

type AttitudePayload struct {
	SysID      uint8   `json:"sys_id"`
	CompID     uint8   `json:"comp_id"`
	TimeBootMs uint32  `json:"time_boot_ms"`
	Roll       float32 `json:"roll_rad"`
	Pitch      float32 `json:"pitch_rad"`
	Yaw        float32 `json:"yaw_rad"`
	RollSpeed  float32 `json:"roll_speed_rads"`
	PitchSpeed float32 `json:"pitch_speed_rads"`
	YawSpeed   float32 `json:"yaw_speed_rads"`
}

type BatteryPayload struct {
	SysID          uint8   `json:"sys_id"`
	CompID         uint8   `json:"comp_id"`
	ID             uint8   `json:"battery_id"`
	BatteryFunction uint32 `json:"function"`
	Type           uint32  `json:"type"`
	Temperature    float32 `json:"temp_c"`
	Voltages       []int32 `json:"cell_voltages_mv"`
	CurrentBattery float32 `json:"current_a"`
	CurrentConsumed int32  `json:"current_consumed_mah"`
	EnergyConsumed  int32  `json:"energy_consumed_hj"`
	BatteryRemaining int8  `json:"remaining_pct"`
}

type VFRHUDPayload struct {
	SysID    uint8   `json:"sys_id"`
	CompID   uint8   `json:"comp_id"`
	Airspeed float32 `json:"airspeed_ms"`
	Groundspeed float32 `json:"groundspeed_ms"`
	Heading  int32   `json:"heading_deg"`
	Throttle uint32  `json:"throttle_pct"`
	Alt      float32 `json:"alt_m"`
	Climb    float32 `json:"climb_ms"`
}

// -----------------------------------------------------------------------
// Command envelope (fabric → drone)
// -----------------------------------------------------------------------

type Command struct {
	ID      string          `json:"id"`
	TS      int64           `json:"ts"`
	Entity  string          `json:"entity"`
	Action  string          `json:"action"`
	Payload json.RawMessage `json:"payload"`
}

// -----------------------------------------------------------------------
// Config — environment only
// -----------------------------------------------------------------------

type Config struct {
	NATSUrl    string
	BridgeName string
	EntityID   string
	DroneAddr  string // host:port UDP, e.g. "192.168.1.100:14550"
}

func loadConfig() Config {
	droneIP := getEnv("DRONE_IP", "")
	if droneIP == "" {
		log.Fatal("[mavlink-bridge] DRONE_IP is required")
	}
	return Config{
		NATSUrl:    getEnv("NATS_URL", "nats://localhost:4222"),
		BridgeName: getEnv("BRIDGE_NAME", "mavlink-bridge"),
		EntityID:   getEnv("BRIDGE_ENTITY_ID", "drone-alpha-1"),
		DroneAddr:  fmt.Sprintf("%s:%s", droneIP, getEnv("DRONE_PORT", "14550")),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

// -----------------------------------------------------------------------
// Main
// -----------------------------------------------------------------------

func main() {
	cfg := loadConfig()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// --- NATS connect ---
	nc, err := nats.Connect(cfg.NATSUrl,
		nats.MaxReconnects(-1),
		nats.ReconnectWait(2*time.Second),
		nats.DisconnectErrHandler(func(nc *nats.Conn, err error) {
			log.Printf("[%s] NATS disconnected: %v", cfg.BridgeName, err)
		}),
		nats.ReconnectHandler(func(nc *nats.Conn) {
			log.Printf("[%s] NATS reconnected: %s", cfg.BridgeName, nc.ConnectedUrl())
		}),
	)
	if err != nil {
		log.Fatalf("[%s] NATS connect failed: %v", cfg.BridgeName, err)
	}
	defer nc.Drain()
	log.Printf("[%s] NATS connected: %s", cfg.BridgeName, nc.ConnectedUrl())

	// --- Subscribe to inbound commands ---
	cmdSubject := fmt.Sprintf("omos.air.%s.command", cfg.EntityID)
	nc.Subscribe(cmdSubject, func(msg *nats.Msg) {
		var cmd Command
		if err := json.Unmarshal(msg.Data, &cmd); err != nil {
			publishError(nc, cfg.BridgeName, fmt.Sprintf("invalid command envelope: %v", err))
			return
		}
		log.Printf("[%s] command received: action=%s entity=%s", cfg.BridgeName, cmd.Action, cmd.Entity)
		// TODO v0.3: marshal cmd.Action → MAVLink COMMAND_LONG and send via node.WriteMessageTo
	})
	log.Printf("[%s] subscribed to commands: %s", cfg.BridgeName, cmdSubject)

	// --- Heartbeat goroutine ---
	go runHeartbeat(ctx, nc, cfg)

	// --- gomavlib node ---
	node, err := gomavlib.NewNode(gomavlib.NodeConf{
		Endpoints: []gomavlib.EndpointConf{
			gomavlib.EndpointUDPClient{
				Address: cfg.DroneAddr,
			},
		},
		Dialect:     ardupilotmega.Dialect,
		OutVersion:  gomavlib.V2,
		OutSystemID: 254, // GCS system ID
	})
	if err != nil {
		log.Fatalf("[%s] gomavlib node failed: %v", cfg.BridgeName, err)
	}
	defer node.Close()
	log.Printf("[%s] MAVLink node connected to %s", cfg.BridgeName, cfg.DroneAddr)

	// --- MAVLink event loop ---
	go func() {
		for evt := range node.Events() {
			select {
			case <-ctx.Done():
				return
			default:
			}

			frm, ok := evt.(*gomavlib.EventFrame)
			if !ok {
				continue
			}

			sysID := frm.SystemID()
			compID := frm.ComponentID()

			switch msg := frm.Message().(type) {

			case *ardupilotmega.MessageHeartbeat:
				pub(nc, cfg, "heartbeat", HeartbeatPayload{
					SysID:          sysID,
					CompID:         compID,
					Type:           uint32(msg.Type),
					Autopilot:      uint32(msg.Autopilot),
					BaseMode:       uint32(msg.BaseMode),
					SystemStatus:   uint32(msg.SystemStatus),
					MavlinkVersion: uint32(msg.MavlinkVersion),
				})

			case *ardupilotmega.MessageGlobalPositionInt:
				pub(nc, cfg, "position", PositionPayload{
					SysID:      sysID,
					CompID:     compID,
					TimeBootMs: msg.TimeBootMs,
					Lat:        float64(msg.Lat) / 1e7,
					Lon:        float64(msg.Lon) / 1e7,
					Alt:        float64(msg.Alt) / 1000.0,
					RelAlt:     float64(msg.RelativeAlt) / 1000.0,
					Vx:         float32(msg.Vx) / 100.0,
					Vy:         float32(msg.Vy) / 100.0,
					Vz:         float32(msg.Vz) / 100.0,
					Hdg:        float32(msg.Hdg) / 100.0,
				})

			case *ardupilotmega.MessageAttitude:
				pub(nc, cfg, "attitude", AttitudePayload{
					SysID:      sysID,
					CompID:     compID,
					TimeBootMs: msg.TimeBootMs,
					Roll:       msg.Roll,
					Pitch:      msg.Pitch,
					Yaw:        msg.Yaw,
					RollSpeed:  msg.Rollspeed,
					PitchSpeed: msg.Pitchspeed,
					YawSpeed:   msg.Yawspeed,
				})

			case *ardupilotmega.MessageBatteryStatus:
				voltages := make([]int32, len(msg.Voltages))
				for i, v := range msg.Voltages {
					voltages[i] = int32(v)
				}
				pub(nc, cfg, "battery", BatteryPayload{
					SysID:            sysID,
					CompID:           compID,
					ID:               msg.Id,
					BatteryFunction:  uint32(msg.BatteryFunction),
					Type:             uint32(msg.Type),
					Temperature:      float32(msg.Temperature) / 100.0,
					Voltages:         voltages,
					CurrentBattery:   float32(msg.CurrentBattery) / 100.0,
					CurrentConsumed:  msg.CurrentConsumed,
					EnergyConsumed:   msg.EnergyConsumed,
					BatteryRemaining: msg.BatteryRemaining,
				})

			case *ardupilotmega.MessageVfrHud:
				pub(nc, cfg, "vfrhud", VFRHUDPayload{
					SysID:       sysID,
					CompID:      compID,
					Airspeed:    msg.Airspeed,
					Groundspeed: msg.Groundspeed,
					Heading:     int32(msg.Heading),
					Throttle:    uint32(msg.Throttle),
					Alt:         msg.Alt,
					Climb:       msg.Climb,
				})
			}
		}
	}()

	log.Printf("[%s] bridge running — publishing to omos.air.%s.telemetry", cfg.BridgeName, cfg.EntityID)

	// --- Graceful shutdown ---
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit
	log.Printf("[%s] shutting down gracefully", cfg.BridgeName)
	cancel()
}

// -----------------------------------------------------------------------
// Helpers
// -----------------------------------------------------------------------

// pub marshals any payload into the standard One Mind envelope and publishes.
// Subject: omos.air.{entity_id}.telemetry
// The msg_type is embedded in the payload — consumers can filter by it.
func pub(nc *nats.Conn, cfg Config, msgType string, payload interface{}) {
	type wrappedPayload struct {
		MsgType string      `json:"msg_type"`
		Data    interface{} `json:"data"`
	}

	inner, err := json.Marshal(wrappedPayload{MsgType: msgType, Data: payload})
	if err != nil {
		publishError(nc, cfg.BridgeName, fmt.Sprintf("marshal error: %v", err))
		return
	}

	env := BridgeMessage{
		ID:      fmt.Sprintf("%d", time.Now().UnixNano()),
		TS:      time.Now().UnixMilli(),
		Source:  cfg.BridgeName,
		Entity:  fmt.Sprintf("air.%s", cfg.EntityID),
		Domain:  "air",
		Payload: inner,
	}

	data, err := json.Marshal(env)
	if err != nil {
		publishError(nc, cfg.BridgeName, fmt.Sprintf("envelope marshal error: %v", err))
		return
	}

	subject := fmt.Sprintf("omos.air.%s.telemetry", cfg.EntityID)
	if err := nc.Publish(subject, data); err != nil {
		log.Printf("[%s] publish error: %v", cfg.BridgeName, err)
	}
}

func runHeartbeat(ctx context.Context, nc *nats.Conn, cfg Config) {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			payload, _ := json.Marshal(map[string]interface{}{
				"status": "ok",
				"ts":     time.Now().UnixMilli(),
				"bridge": cfg.BridgeName,
			})
			nc.Publish(fmt.Sprintf("omos.bridge.%s.health", cfg.BridgeName), payload)
		}
	}
}

func publishError(nc *nats.Conn, bridgeName, errMsg string) {
	payload, _ := json.Marshal(map[string]interface{}{
		"ts":    time.Now().UnixMilli(),
		"error": errMsg,
	})
	nc.Publish(fmt.Sprintf("omos.bridge.%s.errors", bridgeName), payload)
}
