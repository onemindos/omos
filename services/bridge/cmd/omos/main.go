// SPDX-License-Identifier: Apache-2.0
// omos - One Mind Fabric CLI
// Developer CLI for the omos-bridge fabric. Connects to NATS to inspect
// live message traffic, entity presence, and bridge health.

package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"sort"
	"strings"
	"sync"
	"syscall"
	"time"

	nats "github.com/nats-io/nats.go"
)

// ─────────────────────────────────────────────
// Bridge catalog (static, no NATS needed)
// ─────────────────────────────────────────────

type bridgeEntry struct {
	Name     string
	Protocol string
	Domain   string
	Status   string
}

var bridgeCatalog = []bridgeEntry{
	{"mavlink", "MAVLink v2", "air", "available"},
	{"tak", "TAK/CoT", "field", "available"},
	{"homeassistant", "HA WS API", "home", "available"},
	{"webhook", "HTTP POST", "net", "available"},
	{"mqtt", "MQTT 3.1/5", "home", "available"},
	{"ros2", "ROS2/DDS", "ground", "planned"},
	{"can", "CAN Bus", "ground", "planned"},
	{"modbus", "Modbus TCP", "infra", "planned"},
	{"opcua", "OPC-UA", "infra", "planned"},
	{"rtsp", "RTSP", "home", "planned"},
}

// ─────────────────────────────────────────────
// Helpers
// ─────────────────────────────────────────────

func natsURL() string {
	if u := os.Getenv("NATS_URL"); u != "" {
		return u
	}
	return "nats://localhost:4222"
}

func connectNATS() (*nats.Conn, error) {
	return nats.Connect(natsURL(),
		nats.Name("omos-cli"),
		nats.Timeout(5*time.Second),
	)
}

func padRight(s string, n int) string {
	if len(s) >= n {
		return s
	}
	return s + strings.Repeat(" ", n-len(s))
}

func humanAgo(d time.Duration) string {
	secs := int(d.Seconds())
	if secs < 1 {
		return "now"
	}
	return fmt.Sprintf("%ds ago", secs)
}

// ─────────────────────────────────────────────
// Usage
// ─────────────────────────────────────────────

func printUsage() {
	fmt.Print(`omos - One Mind Fabric CLI

Usage:
  omos bridge list           List available bridges
  omos bridge health         Check bridge health (35s window)
  omos fabric watch          Watch all fabric messages
  omos fabric watch --domain [domain]  Watch a domain
  omos entity list           Discover entities (5s window)

Environment:
  NATS_URL   NATS server URL (default: nats://localhost:4222)
`)
}

// ─────────────────────────────────────────────
// omos bridge list
// ─────────────────────────────────────────────

func cmdBridgeList() {
	fmt.Printf("%-20s %-12s %-10s %s\n", "BRIDGE", "PROTOCOL", "DOMAIN", "STATUS")
	for _, b := range bridgeCatalog {
		fmt.Printf("%-20s %-12s %-10s %s\n", b.Name, b.Protocol, b.Domain, b.Status)
	}
}

// ─────────────────────────────────────────────
// omos fabric watch [--domain <d>] [subject]
// ─────────────────────────────────────────────

func cmdFabricWatch(args []string) {
	subject := "omos.>"

	// Parse flags from remaining args
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--domain":
			if i+1 < len(args) {
				i++
				subject = fmt.Sprintf("omos.%s.>", args[i])
			} else {
				fmt.Fprintln(os.Stderr, "error: --domain requires an argument")
				os.Exit(1)
			}
		default:
			// positional subject override
			if !strings.HasPrefix(args[i], "-") {
				subject = args[i]
			}
		}
	}

	nc, err := connectNATS()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: cannot connect to NATS at %s: %v\n", natsURL(), err)
		os.Exit(1)
	}
	defer nc.Close()

	fmt.Printf("Watching %s (Ctrl+C to exit)\n\n", subject)

	_, err = nc.Subscribe(subject, func(msg *nats.Msg) {
		ts := time.Now().Format("15:04:05.000")
		fmt.Printf("[%s] %s\n%s\n", ts, msg.Subject, string(msg.Data))
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: subscribe failed: %v\n", err)
		os.Exit(1)
	}

	// Block until Ctrl+C
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh
	fmt.Println("\nexiting.")
}

// ─────────────────────────────────────────────
// omos entity list
// ─────────────────────────────────────────────

type entityStat struct {
	domain      string
	lastMsgType string
	msgs        int
	lastSeen    time.Time
}

func cmdEntityList() {
	nc, err := connectNATS()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: cannot connect to NATS at %s: %v\n", natsURL(), err)
		os.Exit(1)
	}
	defer nc.Close()

	var mu sync.Mutex
	entities := make(map[string]*entityStat)

	fmt.Println("Collecting entities for 5 seconds…")

	_, err = nc.Subscribe("omos.>.telemetry", func(msg *nats.Msg) {
		// subject format: omos.<domain>.<entity-id>.telemetry
		parts := strings.Split(msg.Subject, ".")
		if len(parts) < 4 {
			return
		}
		domain := parts[1]
		entityID := strings.Join(parts[2:len(parts)-1], ".")

		// Try to extract msg_type from JSON payload
		var msgType string
		var raw map[string]interface{}
		if json.Unmarshal(msg.Data, &raw) == nil {
			if p, ok := raw["payload"].(map[string]interface{}); ok {
				if mt, ok := p["msg_type"].(string); ok {
					msgType = mt
				}
			}
			// fallback: top-level type field
			if msgType == "" {
				if mt, ok := raw["type"].(string); ok {
					msgType = mt
				}
			}
		}

		key := domain + "." + entityID

		mu.Lock()
		defer mu.Unlock()
		if e, exists := entities[key]; exists {
			e.msgs++
			e.lastSeen = time.Now()
			if msgType != "" {
				e.lastMsgType = msgType
			}
		} else {
			entities[key] = &entityStat{
				domain:      domain,
				lastMsgType: msgType,
				msgs:        1,
				lastSeen:    time.Now(),
			}
		}
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: subscribe failed: %v\n", err)
		os.Exit(1)
	}

	time.Sleep(5 * time.Second)

	mu.Lock()
	defer mu.Unlock()

	if len(entities) == 0 {
		fmt.Println("No entities observed in 5-second window.")
		return
	}

	// Sort by key for stable output
	keys := make([]string, 0, len(entities))
	for k := range entities {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	fmt.Printf("\n%-30s %-10s %-17s %s\n", "ENTITY", "DOMAIN", "LAST MSG TYPE", "MSGS")
	for _, k := range keys {
		e := entities[k]
		mt := e.lastMsgType
		if mt == "" {
			mt = "—"
		}
		fmt.Printf("%-30s %-10s %-17s %d\n", k, e.domain, mt, e.msgs)
	}
}

// ─────────────────────────────────────────────
// omos bridge health
// ─────────────────────────────────────────────

type bridgeHealth struct {
	status   string
	lastSeen time.Time
}

func cmdBridgeHealth() {
	nc, err := connectNATS()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: cannot connect to NATS at %s: %v\n", natsURL(), err)
		os.Exit(1)
	}
	defer nc.Close()

	var mu sync.Mutex
	bridges := make(map[string]*bridgeHealth)

	fmt.Println("Collecting bridge health for 35 seconds…")

	_, err = nc.Subscribe("omos.bridge.*.health", func(msg *nats.Msg) {
		// subject: omos.bridge.<bridge-name>.health
		parts := strings.Split(msg.Subject, ".")
		if len(parts) < 4 {
			return
		}
		bridgeName := parts[2]

		// Try to decode status from JSON
		status := "ONLINE"
		var raw map[string]interface{}
		if json.Unmarshal(msg.Data, &raw) == nil {
			if s, ok := raw["status"].(string); ok {
				status = strings.ToUpper(s)
			}
		}

		mu.Lock()
		defer mu.Unlock()
		bridges[bridgeName] = &bridgeHealth{
			status:   status,
			lastSeen: time.Now(),
		}
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: subscribe failed: %v\n", err)
		os.Exit(1)
	}

	// Show a countdown
	deadline := time.Now().Add(35 * time.Second)
	for time.Now().Before(deadline) {
		remaining := time.Until(deadline).Round(time.Second)
		fmt.Printf("\r  waiting… %2ds remaining", int(remaining.Seconds()))
		time.Sleep(500 * time.Millisecond)
	}
	fmt.Print("\r                              \r")

	mu.Lock()
	defer mu.Unlock()

	if len(bridges) == 0 {
		fmt.Println("No bridge health messages observed in 35-second window.")
		return
	}

	names := make([]string, 0, len(bridges))
	for n := range bridges {
		names = append(names, n)
	}
	sort.Strings(names)

	fmt.Printf("%-24s %-10s %s\n", "BRIDGE", "STATUS", "LAST SEEN")
	for _, n := range names {
		bh := bridges[n]
		ago := humanAgo(time.Since(bh.lastSeen))
		fmt.Printf("%-24s %-10s %s\n", n, bh.status, ago)
	}
}

// ─────────────────────────────────────────────
// Entry point
// ─────────────────────────────────────────────

func main() {
	args := os.Args[1:]

	if len(args) == 0 {
		printUsage()
		return
	}

	switch args[0] {
	case "bridge":
		if len(args) < 2 {
			printUsage()
			return
		}
		switch args[1] {
		case "list":
			cmdBridgeList()
		case "health":
			cmdBridgeHealth()
		default:
			fmt.Fprintf(os.Stderr, "unknown bridge subcommand: %s\n\n", args[1])
			printUsage()
			os.Exit(1)
		}

	case "fabric":
		if len(args) < 2 {
			printUsage()
			return
		}
		switch args[1] {
		case "watch":
			cmdFabricWatch(args[2:])
		default:
			fmt.Fprintf(os.Stderr, "unknown fabric subcommand: %s\n\n", args[1])
			printUsage()
			os.Exit(1)
		}

	case "entity":
		if len(args) < 2 {
			printUsage()
			return
		}
		switch args[1] {
		case "list":
			cmdEntityList()
		default:
			fmt.Fprintf(os.Stderr, "unknown entity subcommand: %s\n\n", args[1])
			printUsage()
			os.Exit(1)
		}

	default:
		fmt.Fprintf(os.Stderr, "unknown command: %s\n\n", args[0])
		printUsage()
		os.Exit(1)
	}
}
