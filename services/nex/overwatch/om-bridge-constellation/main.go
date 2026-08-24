package main

import (
	"log"
	"os"
	"strings"
	"time"

	"github.com/nats-io/nats.go"
)

// om-bridge-constellation: Bridges Constellation Overwatch NATS (:4299) → OneMind NATS (:4222)

func main() {
	mainURL := os.Getenv("NATS_URL")
	if mainURL == "" {
		mainURL = "nats://127.0.0.1:4222"
	}
	constellationURL := os.Getenv("CONSTELLATION_NATS_URL")
	if constellationURL == "" {
		constellationURL = "nats://127.0.0.1:4299"
	}

	mainNC, err := nats.Connect(mainURL, nats.Name("om-bridge-constellation"))
	if err != nil {
		log.Fatal("main NATS connect:", err)
	}
	defer mainNC.Close()
	log.Printf("om-bridge-constellation: connected to OneMind bus")
	mainNC.Publish("agents.hb.bridge.constellation.zeus", []byte(`{"status":"bridge_up"}`))

	var constNC *nats.Conn
	for i := 0; i < 5; i++ {
		constNC, err = nats.Connect(constellationURL, nats.Name("om-bridge-constellation-remote"))
		if err == nil {
			break
		}
		log.Printf("om-bridge-constellation: not reachable at %s (attempt %d/5)...", constellationURL, i+1)
		time.Sleep(3 * time.Second)
	}

	if constNC != nil {
		defer constNC.Close()
		log.Printf("om-bridge-constellation: BRIDGED at %s", constellationURL)

		// Constellation → OneMind
		constNC.Subscribe(">", func(msg *nats.Msg) {
			mainNC.Publish("tel.entity.constellation."+msg.Subject, msg.Data)
		})

		// OneMind → Constellation (commands/tasking)
		mainNC.Subscribe("cmd.entity.constellation.>", func(msg *nats.Msg) {
			constNC.Publish(strings.TrimPrefix(msg.Subject, "cmd.entity.constellation."), msg.Data)
		})
		log.Println("om-bridge-constellation: bidirectional bridge active")
	} else {
		log.Println("om-bridge-constellation: standing by for Constellation to come online")
	}

	ticker := time.NewTicker(15 * time.Second)
	for range ticker.C {
		status := "standing_by"
		if constNC != nil && constNC.IsConnected() {
			status = "bridged"
		}
		mainNC.Publish("agents.hb.bridge.constellation.zeus", []byte(`{"t":"`+time.Now().Format(time.RFC3339)+`","status":"`+status+`"}`))
	}
}
