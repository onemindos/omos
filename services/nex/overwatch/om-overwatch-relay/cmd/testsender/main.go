package main

// testsender — emits synthetic MAVLink (HEARTBEAT + GLOBAL_POSITION_INT) to a
// UDP endpoint, to prove om-overwatch-relay end-to-end without a real drone.
// Usage: go run ./cmd/testsender 127.0.0.1:14560

import (
	"log"
	"os"
	"time"

	"github.com/bluenviron/gomavlib/v3"
	"github.com/bluenviron/gomavlib/v3/pkg/dialects/common"
)

func main() {
	target := "127.0.0.1:14560"
	if len(os.Args) > 1 {
		target = os.Args[1]
	}

	node, err := gomavlib.NewNode(gomavlib.NodeConf{
		Endpoints:   []gomavlib.EndpointConf{gomavlib.EndpointUDPClient{Address: target}},
		Dialect:     common.Dialect,
		OutVersion:  gomavlib.V2,
		OutSystemID: 7, // this "drone" is system id 7 -> entity drone-7
	})
	if err != nil {
		log.Fatal(err)
	}
	defer node.Close()
	log.Printf("testsender: sending MAVLink to %s as system 7", target)

	go func() {
		for range node.Events() {
		}
	}()

	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()
	count := 0
	for range ticker.C {
		// HEARTBEAT — required for the relay to fingerprint the device.
		node.WriteMessageAll(&common.MessageHeartbeat{
			Type:           common.MAV_TYPE_QUADROTOR,
			Autopilot:      common.MAV_AUTOPILOT_PX4,
			BaseMode:       common.MAV_MODE_FLAG_SAFETY_ARMED,
			CustomMode:     0,
			SystemStatus:   common.MAV_STATE_ACTIVE,
			MavlinkVersion: 3,
		})
		// GLOBAL_POSITION_INT — position telemetry (Kyiv-ish, moving north).
		node.WriteMessageAll(&common.MessageGlobalPositionInt{
			TimeBootMs:  uint32(count * 500),
			Lat:         int32((50.4501 + float64(count)*0.0001) * 1e7),
			Lon:         int32(30.5234 * 1e7),
			Alt:         int32(100 * 1000),
			RelativeAlt: int32(95 * 1000),
			Vx:          500, // cm/s
			Vy:          0,
			Vz:          0,
			Hdg:         0, // centidegrees (north)
		})
		count++
		if count >= 20 {
			log.Println("testsender: done (20 cycles)")
			return
		}
	}
}
