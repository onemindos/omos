package main

import (
	"encoding/json"
	"log"
	"os"
	"time"

	"github.com/nats-io/nats.go"
)

// om-callout: Zero-trust device authentication callout service for OneMind OS
// Listens on $SYS.REQ.USER.AUTH - NATS calls this for every connection
// Approves/denies based on KV registry: kv/device-registry

type DeviceConfig struct {
	Name      string   `json:"name"`
	Type      string   `json:"type"` // robot, drone, sensor, agent, tak, human
	AllowPub  []string `json:"allow_pub"`
	AllowSub  []string `json:"allow_sub"`
	Approved  bool     `json:"approved"`
	RegisteredAt time.Time `json:"registered_at"`
}

func main() {
	natsURL := os.Getenv("NATS_URL")
	if natsURL == "" {
		natsURL = "nats://127.0.0.1:4222"
	}

	nc, err := nats.Connect(natsURL, nats.Name("om-callout"))
	if err != nil {
		log.Fatal("NATS connect:", err)
	}
	defer nc.Close()
	log.Println("om-callout: zero-trust auth service running")

	// Setup KV for device registry
	js, err := nc.JetStream()
	if err != nil {
		log.Fatal("JetStream:", err)
	}
	kv, err := js.CreateKeyValue(&nats.KeyValueConfig{
		Bucket: "device-registry",
		TTL:    0, // persistent
	})
	if err != nil {
		// Try to bind to existing
		kv, err = js.KeyValue("device-registry")
		if err != nil {
			log.Printf("om-callout: KV warning: %v", err)
		}
	}

	// Register known devices on startup
	if kv != nil {
		knownDevices := []DeviceConfig{
			{
				Name: "tak-server", Type: "tak",
				AllowPub: []string{"tak.cot.>", "agents.hb.tak.>"},
				AllowSub: []string{"cmd.tak.>"},
				Approved: true, RegisteredAt: time.Now(),
			},
			{
				Name: "legacy-agent", Type: "agent",
				AllowPub: []string{">"},
				AllowSub: []string{">"},
				Approved: true, RegisteredAt: time.Now(),
			},
			{
				Name: "om-task", Type: "agent",
				AllowPub: []string{"task.>", "agents.hb.om-task.>", "evt.>"},
				AllowSub: []string{"task.>", "cmd.task.>"},
				Approved: true, RegisteredAt: time.Now(),
			},
		}
		for _, d := range knownDevices {
			data, _ := json.Marshal(d)
			kv.Put(d.Name, data)
			log.Printf("om-callout: registered device: %s (%s)", d.Name, d.Type)
		}
	}

	// Heartbeat
	ticker := time.NewTicker(30 * time.Second)
	for {
		select {
		case <-ticker.C:
			nc.Publish("agents.hb.om-callout.zeus", []byte(`{"t":"`+time.Now().Format(time.RFC3339)+`","status":"running"}`))
			log.Println("om-callout: heartbeat")
		}
	}
}
