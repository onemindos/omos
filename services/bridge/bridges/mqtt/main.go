// SPDX-License-Identifier: Apache-2.0
// omos-bridge — MQTT bridge
// Bidirectional bridge: MQTT ↔ NATS using the One Mind envelope schema.

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/nats-io/nats.go"
)

// ---------------------------------------------------------------------------
// Config
// ---------------------------------------------------------------------------

type Config struct {
	NatsURL      string
	BridgeName   string
	MQTTHost     string
	MQTTPort     string
	MQTTUsername string
	MQTTPassword string
	MQTTClientID string
	TopicFilter  string
}

func loadConfig() Config {
	return Config{
		NatsURL:      getEnv("NATS_URL", "nats://localhost:4222"),
		BridgeName:   getEnv("BRIDGE_NAME", "mqtt-bridge"),
		MQTTHost:     getEnv("MQTT_HOST", "localhost"),
		MQTTPort:     getEnv("MQTT_PORT", "1883"),
		MQTTUsername: getEnv("MQTT_USERNAME", ""),
		MQTTPassword: getEnv("MQTT_PASSWORD", ""),
		MQTTClientID: getEnv("MQTT_CLIENT_ID", "omos-mqtt-bridge"),
		TopicFilter:  getEnv("MQTT_TOPIC_FILTER", "#"),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

// ---------------------------------------------------------------------------
// Envelope
// ---------------------------------------------------------------------------

type Envelope struct {
	Entity    string                 `json:"entity"`
	Domain    string                 `json:"domain"`
	MsgType   string                 `json:"msg_type"`
	Timestamp int64                  `json:"timestamp"`
	Data      map[string]interface{} `json:"data"`
}

// ---------------------------------------------------------------------------
// Known domains for first-segment detection
// ---------------------------------------------------------------------------

var knownDomains = map[string]bool{
	"home":   true,
	"air":    true,
	"ground": true,
	"field":  true,
	"infra":  true,
}

// sanitizeTopic replaces / with - and trims leading/trailing dashes.
func sanitizeTopic(topic string) string {
	s := strings.ReplaceAll(topic, "/", "-")
	return strings.Trim(s, "-")
}

// domainFromTopic inspects the first segment of the MQTT topic.
// Returns the segment if it is a known domain, otherwise "home".
func domainFromTopic(topic string) string {
	parts := strings.SplitN(topic, "/", 2)
	if len(parts) > 0 && knownDomains[parts[0]] {
		return parts[0]
	}
	return "home"
}

// ---------------------------------------------------------------------------
// Bridge
// ---------------------------------------------------------------------------

type bridge struct {
	cfg  Config
	nc   *nats.Conn
	mc   mqtt.Client
}

// ---------------------------------------------------------------------------
// MQTT → NATS
// ---------------------------------------------------------------------------

func (b *bridge) onMQTTMessage(_ mqtt.Client, msg mqtt.Message) {
	topic := msg.Topic()
	payloadBytes := msg.Payload()

	domain := domainFromTopic(topic)
	sanitized := sanitizeTopic(topic)

	// Try JSON parse
	var payloadParsed interface{}
	var tmp map[string]interface{}
	if json.Unmarshal(payloadBytes, &tmp) == nil {
		payloadParsed = tmp
	}

	env := Envelope{
		Entity:    fmt.Sprintf("%s.%s", domain, sanitized),
		Domain:    domain,
		MsgType:   "sensor",
		Timestamp: time.Now().UnixMilli(),
		Data: map[string]interface{}{
			"mqtt_topic":     topic,
			"payload_raw":    string(payloadBytes),
			"payload_parsed": payloadParsed,
			"qos":            int(msg.Qos()),
			"retained":       msg.Retained(),
		},
	}

	data, err := json.Marshal(env)
	if err != nil {
		b.reportError(fmt.Sprintf("marshal error for topic %s: %v", topic, err))
		return
	}

	subject := fmt.Sprintf("omos.%s.%s.telemetry", domain, sanitized)
	if err := b.nc.Publish(subject, data); err != nil {
		b.reportError(fmt.Sprintf("NATS publish error for %s: %v", subject, err))
		return
	}

	log.Printf("[mqtt] MQTT %s → NATS %s", topic, subject)
}

// ---------------------------------------------------------------------------
// NATS → MQTT (commands)
// ---------------------------------------------------------------------------

func (b *bridge) subscribeNATSCommands() (*nats.Subscription, error) {
	return b.nc.Subscribe("omos.*.*.command", func(msg *nats.Msg) {
		// subject format: omos.{domain}.{entity}.command
		parts := strings.Split(msg.Subject, ".")
		if len(parts) < 4 {
			log.Printf("[mqtt] unexpected command subject: %s", msg.Subject)
			return
		}
		domain := parts[1]
		entity := parts[2]

		mqttTopic := fmt.Sprintf("omos/%s/%s/command", domain, entity)
		token := b.mc.Publish(mqttTopic, 1, false, msg.Data)
		token.Wait()
		if err := token.Error(); err != nil {
			b.reportError(fmt.Sprintf("MQTT publish error for %s: %v", mqttTopic, err))
			return
		}
		log.Printf("[mqtt] NATS %s → MQTT %s", msg.Subject, mqttTopic)
	})
}

// ---------------------------------------------------------------------------
// Error reporting & heartbeat
// ---------------------------------------------------------------------------

func (b *bridge) reportError(msg string) {
	log.Printf("[mqtt] ERROR: %s", msg)
	payload, _ := json.Marshal(map[string]interface{}{
		"bridge":    b.cfg.BridgeName,
		"error":     msg,
		"timestamp": time.Now().UnixMilli(),
	})
	_ = b.nc.Publish(fmt.Sprintf("omos.bridge.%s.errors", b.cfg.BridgeName), payload)
}

func (b *bridge) startHeartbeat(ctx context.Context) {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			payload, _ := json.Marshal(map[string]interface{}{
				"bridge":      b.cfg.BridgeName,
				"status":      "ok",
				"mqtt_connected": b.mc.IsConnected(),
				"timestamp":   time.Now().UnixMilli(),
			})
			subject := fmt.Sprintf("omos.bridge.%s.health", b.cfg.BridgeName)
			if err := b.nc.Publish(subject, payload); err != nil {
				log.Printf("[mqtt] heartbeat publish failed: %v", err)
			} else {
				log.Printf("[mqtt] heartbeat → %s", subject)
			}
		}
	}
}

// ---------------------------------------------------------------------------
// Main
// ---------------------------------------------------------------------------

func main() {
	cfg := loadConfig()

	// ── NATS ─────────────────────────────────────────────────────────────────
	nc, err := nats.Connect(cfg.NatsURL,
		nats.MaxReconnects(-1),
		nats.ReconnectWait(2*time.Second),
		nats.DisconnectErrHandler(func(_ *nats.Conn, err error) {
			log.Printf("[mqtt] NATS disconnected: %v", err)
		}),
		nats.ReconnectHandler(func(_ *nats.Conn) {
			log.Printf("[mqtt] NATS reconnected")
		}),
	)
	if err != nil {
		log.Fatalf("[mqtt] NATS connect failed: %v", err)
	}
	defer nc.Drain()
	log.Printf("[mqtt] connected to NATS at %s", cfg.NatsURL)

	// ── MQTT ─────────────────────────────────────────────────────────────────
	brokerURL := fmt.Sprintf("tcp://%s:%s", cfg.MQTTHost, cfg.MQTTPort)
	opts := mqtt.NewClientOptions().
		AddBroker(brokerURL).
		SetClientID(cfg.MQTTClientID).
		SetAutoReconnect(true).
		SetConnectRetry(true).
		SetConnectRetryInterval(3 * time.Second).
		SetOnConnectHandler(func(_ mqtt.Client) {
			log.Printf("[mqtt] connected to MQTT broker at %s", brokerURL)
		}).
		SetConnectionLostHandler(func(_ mqtt.Client, err error) {
			log.Printf("[mqtt] MQTT connection lost: %v — auto-reconnecting", err)
		})

	if cfg.MQTTUsername != "" {
		opts.SetUsername(cfg.MQTTUsername)
	}
	if cfg.MQTTPassword != "" {
		opts.SetPassword(cfg.MQTTPassword)
	}

	mc := mqtt.NewClient(opts)
	if token := mc.Connect(); token.Wait() && token.Error() != nil {
		log.Fatalf("[mqtt] MQTT connect failed: %v", token.Error())
	}

	b := &bridge{cfg: cfg, nc: nc, mc: mc}

	// Subscribe to MQTT wildcard
	if token := mc.Subscribe(cfg.TopicFilter, 1, b.onMQTTMessage); token.Wait() && token.Error() != nil {
		log.Fatalf("[mqtt] MQTT subscribe failed: %v", token.Error())
	}
	log.Printf("[mqtt] subscribed to MQTT topic filter: %s", cfg.TopicFilter)

	// Subscribe to NATS commands
	sub, err := b.subscribeNATSCommands()
	if err != nil {
		log.Fatalf("[mqtt] NATS command subscribe failed: %v", err)
	}
	defer sub.Unsubscribe()
	log.Printf("[mqtt] subscribed to NATS omos.*.*.command")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go b.startHeartbeat(ctx)

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	sig := <-quit
	log.Printf("[mqtt] received signal %s — shutting down", sig)

	cancel()
	mc.Disconnect(500)
	log.Printf("[mqtt] shutdown complete")
}
