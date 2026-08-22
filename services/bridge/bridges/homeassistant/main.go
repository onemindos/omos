// SPDX-License-Identifier: Apache-2.0

// homeassistant-bridge connects a Home Assistant instance to the One Mind OMOS
// message bus over NATS.  It subscribes to HA state_changed events via the HA
// WebSocket API and publishes standardised telemetry and event envelopes onto
// NATS subjects.  It also listens for NATS command messages and forwards them
// to HA as call_service WebSocket requests.
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"
	"os/signal"
	"strings"
	"sync"
	"sync/atomic"
	"syscall"
	"time"

	nats "github.com/nats-io/nats.go"
	"golang.org/x/net/websocket"
)

// ---------------------------------------------------------------------------
// Standard One Mind message envelope — same across ALL omos-bridge bridges
// ---------------------------------------------------------------------------

type BridgeMessage struct {
	ID      string          `json:"id"`
	TS      int64           `json:"ts"`
	Source  string          `json:"source"`
	Entity  string          `json:"entity"`
	Domain  string          `json:"domain"`
	MsgType string          `json:"msg_type"`
	Payload json.RawMessage `json:"payload"`
}

// ---------------------------------------------------------------------------
// HA WebSocket message shapes
// ---------------------------------------------------------------------------

type haBaseMsg struct {
	Type string `json:"type"`
}

type haAuthMsg struct {
	Type        string `json:"type"`
	AccessToken string `json:"access_token"`
}

type haSubscribeMsg struct {
	ID        int    `json:"id"`
	Type      string `json:"type"`
	EventType string `json:"event_type"`
}

type haCallServiceMsg struct {
	ID          int             `json:"id"`
	Type        string          `json:"type"`
	Domain      string          `json:"domain"`
	Service     string          `json:"service"`
	ServiceData json.RawMessage `json:"service_data,omitempty"`
}

// Incoming event envelope from HA
type haEventEnvelope struct {
	ID    int `json:"id"`
	Type  string `json:"type"`
	Event struct {
		EventType string `json:"event_type"`
		Data      struct {
			EntityID string          `json:"entity_id"`
			NewState json.RawMessage `json:"new_state"`
			OldState json.RawMessage `json:"old_state"`
		} `json:"data"`
	} `json:"event"`
}

type haStateObj struct {
	State      string                 `json:"state"`
	Attributes map[string]interface{} `json:"attributes"`
	LastChanged string               `json:"last_changed"`
	Context    struct {
		ID string `json:"id"`
	} `json:"context"`
}

// ---------------------------------------------------------------------------
// Payloads published onto NATS
// ---------------------------------------------------------------------------

type telemetryPayload struct {
	EntityID     string                 `json:"entity_id"`
	FriendlyName string                 `json:"friendly_name,omitempty"`
	State        string                 `json:"state"`
	Unit         string                 `json:"unit,omitempty"`
	DeviceClass  string                 `json:"device_class,omitempty"`
	Attributes   map[string]interface{} `json:"attributes"`
}

type eventPayload struct {
	EventType string `json:"event_type"`
	EntityID  string `json:"entity_id"`
	OldState  string `json:"old_state,omitempty"`
	NewState  string `json:"new_state"`
	ContextID string `json:"context_id,omitempty"`
	OccurredAt int64 `json:"occurred_at"`
}

// ---------------------------------------------------------------------------
// Command message received on NATS omos.home.*.command
// ---------------------------------------------------------------------------

type commandPayload struct {
	Domain      string          `json:"domain"`
	Service     string          `json:"service"`
	ServiceData json.RawMessage `json:"service_data,omitempty"`
}

// ---------------------------------------------------------------------------
// Config
// ---------------------------------------------------------------------------

type config struct {
	natsURL    string
	bridgeName string
	haHost     string
	haPort     string
	haToken    string
}

func loadConfig() config {
	c := config{
		natsURL:    getEnv("NATS_URL", "nats://localhost:4222"),
		bridgeName: getEnv("BRIDGE_NAME", "homeassistant-bridge"),
		haHost:     getEnv("HA_HOST", "localhost"),
		haPort:     getEnv("HA_PORT", "8123"),
		haToken:    os.Getenv("HA_TOKEN"),
	}
	if c.haToken == "" {
		log.Fatal("[bridge] HA_TOKEN env var is required")
	}
	return c
}

func getEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

// ---------------------------------------------------------------------------
// Bridge
// ---------------------------------------------------------------------------

type bridge struct {
	cfg    config
	nc     *nats.Conn
	ws     *websocket.Conn
	mu     sync.Mutex   // guards ws + msgID
	msgID  int64        // atomic counter for HA message IDs
	done   chan struct{}
}

func newBridge(cfg config) *bridge {
	return &bridge{
		cfg:  cfg,
		done: make(chan struct{}),
	}
}

func (b *bridge) nextID() int {
	return int(atomic.AddInt64(&b.msgID, 1))
}

// ---------------------------------------------------------------------------
// NATS helpers
// ---------------------------------------------------------------------------

func (b *bridge) connectNATS() error {
	nc, err := nats.Connect(b.cfg.natsURL,
		nats.MaxReconnects(-1),
		nats.ReconnectWait(2*time.Second),
		nats.DisconnectErrHandler(func(_ *nats.Conn, err error) {
			log.Printf("[nats] disconnected: %v", err)
		}),
		nats.ReconnectHandler(func(_ *nats.Conn) {
			log.Printf("[nats] reconnected")
		}),
	)
	if err != nil {
		return fmt.Errorf("nats connect: %w", err)
	}
	b.nc = nc
	log.Printf("[nats] connected to %s", b.cfg.natsURL)
	return nil
}

func (b *bridge) publish(subject, domain, entity, msgType string, payload interface{}) {
	raw, err := json.Marshal(payload)
	if err != nil {
		log.Printf("[nats] marshal error: %v", err)
		return
	}
	env := BridgeMessage{
		ID:      fmt.Sprintf("%s-%d", b.cfg.bridgeName, time.Now().UnixNano()),
		TS:      time.Now().UnixMilli(),
		Source:  b.cfg.bridgeName,
		Entity:  entity,
		Domain:  domain,
		MsgType: msgType,
		Payload: raw,
	}
	data, err := json.Marshal(env)
	if err != nil {
		log.Printf("[nats] envelope marshal error: %v", err)
		return
	}
	if err := b.nc.Publish(subject, data); err != nil {
		log.Printf("[nats] publish %s error: %v", subject, err)
	}
}

func (b *bridge) publishError(msg string) {
	errPayload := map[string]interface{}{
		"bridge": b.cfg.bridgeName,
		"error":  msg,
		"ts":     time.Now().UnixMilli(),
	}
	raw, _ := json.Marshal(errPayload)
	_ = b.nc.Publish("omos.bridge."+b.cfg.bridgeName+".errors", raw)
}

// ---------------------------------------------------------------------------
// HA WebSocket helpers
// ---------------------------------------------------------------------------

func (b *bridge) haURL() string {
	return fmt.Sprintf("ws://%s:%s/api/websocket", b.cfg.haHost, b.cfg.haPort)
}

func (b *bridge) connectHA() (*websocket.Conn, error) {
	url := b.haURL()
	origin := fmt.Sprintf("http://%s:%s", b.cfg.haHost, b.cfg.haPort)
	ws, err := websocket.Dial(url, "", origin)
	if err != nil {
		return nil, fmt.Errorf("ws dial %s: %w", url, err)
	}
	return ws, nil
}

func (b *bridge) sendJSON(ws *websocket.Conn, v interface{}) error {
	data, err := json.Marshal(v)
	if err != nil {
		return err
	}
	return websocket.Message.Send(ws, string(data))
}

func (b *bridge) recvJSON(ws *websocket.Conn, v interface{}) error {
	var raw string
	if err := websocket.Message.Receive(ws, &raw); err != nil {
		return err
	}
	return json.Unmarshal([]byte(raw), v)
}

// authenticate performs the HA auth handshake.
func (b *bridge) authenticate(ws *websocket.Conn) error {
	// 1. Expect auth_required
	var m haBaseMsg
	if err := b.recvJSON(ws, &m); err != nil {
		return fmt.Errorf("recv auth_required: %w", err)
	}
	if m.Type != "auth_required" {
		return fmt.Errorf("expected auth_required, got %s", m.Type)
	}
	// 2. Send auth
	if err := b.sendJSON(ws, haAuthMsg{Type: "auth", AccessToken: b.cfg.haToken}); err != nil {
		return fmt.Errorf("send auth: %w", err)
	}
	// 3. Expect auth_ok
	if err := b.recvJSON(ws, &m); err != nil {
		return fmt.Errorf("recv auth_ok: %w", err)
	}
	if m.Type != "auth_ok" {
		return fmt.Errorf("expected auth_ok, got %s (invalid token?)", m.Type)
	}
	log.Printf("[ha] authenticated OK")
	return nil
}

// subscribeEvents sends the subscribe_events message.
func (b *bridge) subscribeEvents(ws *websocket.Conn) error {
	id := b.nextID()
	return b.sendJSON(ws, haSubscribeMsg{
		ID:        id,
		Type:      "subscribe_events",
		EventType: "state_changed",
	})
}

// ---------------------------------------------------------------------------
// Domain resolution
// ---------------------------------------------------------------------------

var homeDomains = map[string]bool{
	"light": true, "switch": true, "sensor": true, "binary_sensor": true,
	"climate": true, "cover": true, "fan": true, "media_player": true,
	"automation": true, "script": true, "scene": true,
}

func domainOf(entityID string) string {
	if idx := strings.Index(entityID, "."); idx >= 0 {
		return entityID[:idx]
	}
	return entityID
}

func natsEntity(entityID string) string {
	return strings.ReplaceAll(entityID, ".", "-")
}

func omosDomain(entityID string) string {
	d := domainOf(entityID)
	if homeDomains[d] {
		return "home"
	}
	return d
}

// ---------------------------------------------------------------------------
// Event processing
// ---------------------------------------------------------------------------

func (b *bridge) handleEvent(env haEventEnvelope) {
	if env.Event.EventType != "state_changed" {
		return
	}
	data := env.Event.Data
	entityID := data.EntityID

	// Parse new_state
	var newState haStateObj
	if err := json.Unmarshal(data.NewState, &newState); err != nil {
		log.Printf("[ha] unmarshal new_state for %s: %v", entityID, err)
		return
	}

	// Parse old_state (may be null)
	var oldState haStateObj
	oldStateStr := ""
	if len(data.OldState) > 0 && string(data.OldState) != "null" {
		if err := json.Unmarshal(data.OldState, &oldState); err == nil {
			oldStateStr = oldState.State
		}
	}

	// Parse last_changed to unix ms
	var occurredAt int64
	if t, err := time.Parse(time.RFC3339Nano, newState.LastChanged); err == nil {
		occurredAt = t.UnixMilli()
	} else if t, err := time.Parse(time.RFC3339, newState.LastChanged); err == nil {
		occurredAt = t.UnixMilli()
	} else {
		occurredAt = time.Now().UnixMilli()
	}

	entity := natsEntity(entityID)
	domain := omosDomain(entityID)

	// -- a) Telemetry --
	attrs := newState.Attributes
	friendlyName, _ := attrs["friendly_name"].(string)
	unit, _ := attrs["unit_of_measurement"].(string)
	deviceClass, _ := attrs["device_class"].(string)

	tel := telemetryPayload{
		EntityID:     entityID,
		FriendlyName: friendlyName,
		State:        newState.State,
		Unit:         unit,
		DeviceClass:  deviceClass,
		Attributes:   attrs,
	}
	b.publish(
		"omos.home."+entity+".telemetry",
		domain, entity, "sensor", tel,
	)

	// -- b) Event --
	evt := eventPayload{
		EventType:  "state_changed",
		EntityID:   entityID,
		OldState:   oldStateStr,
		NewState:   newState.State,
		ContextID:  newState.Context.ID,
		OccurredAt: occurredAt,
	}
	b.publish(
		"omos.home."+entity+".event",
		domain, entity, "event", evt,
	)
}

// ---------------------------------------------------------------------------
// NATS command subscription
// ---------------------------------------------------------------------------

func (b *bridge) subscribeCommands() (*nats.Subscription, error) {
	sub, err := b.nc.Subscribe("omos.home.*.command", func(msg *nats.Msg) {
		var env BridgeMessage
		if err := json.Unmarshal(msg.Data, &env); err != nil {
			log.Printf("[nats] command unmarshal error: %v", err)
			return
		}
		var cmd commandPayload
		if err := json.Unmarshal(env.Payload, &cmd); err != nil {
			// Try direct decode (no envelope wrapper)
			if err2 := json.Unmarshal(msg.Data, &cmd); err2 != nil {
				log.Printf("[nats] command payload decode error: %v / %v", err, err2)
				return
			}
		}
		if cmd.Domain == "" || cmd.Service == "" {
			log.Printf("[nats] command missing domain/service")
			return
		}
		id := b.nextID()
		svcMsg := haCallServiceMsg{
			ID:          id,
			Type:        "call_service",
			Domain:      cmd.Domain,
			Service:     cmd.Service,
			ServiceData: cmd.ServiceData,
		}
		b.mu.Lock()
		ws := b.ws
		b.mu.Unlock()
		if ws == nil {
			log.Printf("[ha] command dropped — not connected (domain=%s service=%s)", cmd.Domain, cmd.Service)
			return
		}
		if err := b.sendJSON(ws, svcMsg); err != nil {
			log.Printf("[ha] call_service send error: %v", err)
			b.publishError(fmt.Sprintf("call_service send: %v", err))
		} else {
			log.Printf("[ha] call_service id=%d domain=%s service=%s", id, cmd.Domain, cmd.Service)
		}
	})
	return sub, err
}

// ---------------------------------------------------------------------------
// Heartbeat
// ---------------------------------------------------------------------------

func (b *bridge) startHeartbeat(ctx context.Context) {
	go func() {
		t := time.NewTicker(30 * time.Second)
		defer t.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-t.C:
				health := map[string]interface{}{
					"bridge": b.cfg.bridgeName,
					"ts":     time.Now().UnixMilli(),
					"status": "ok",
				}
				raw, _ := json.Marshal(health)
				_ = b.nc.Publish("omos.bridge."+b.cfg.bridgeName+".health", raw)
			}
		}
	}()
}

// ---------------------------------------------------------------------------
// Main HA read loop with reconnect
// ---------------------------------------------------------------------------

func (b *bridge) runHA(ctx context.Context) {
	backoff := 1 * time.Second
	maxBackoff := 30 * time.Second

	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		log.Printf("[ha] connecting to %s", b.haURL())
		ws, err := b.connectHA()
		if err != nil {
			log.Printf("[ha] connect error: %v — retry in %v", err, backoff)
			b.publishError(fmt.Sprintf("connect: %v", err))
			select {
			case <-ctx.Done():
				return
			case <-time.After(backoff):
			}
			backoff = minDuration(time.Duration(float64(backoff)*1.5)+time.Duration(rand.Intn(1000))*time.Millisecond, maxBackoff)
			continue
		}
		backoff = 1 * time.Second // reset on success

		if err := b.authenticate(ws); err != nil {
			log.Printf("[ha] auth error: %v", err)
			b.publishError(fmt.Sprintf("auth: %v", err))
			_ = ws.Close()
			select {
			case <-ctx.Done():
				return
			case <-time.After(5 * time.Second):
			}
			continue
		}

		if err := b.subscribeEvents(ws); err != nil {
			log.Printf("[ha] subscribe_events error: %v", err)
			b.publishError(fmt.Sprintf("subscribe_events: %v", err))
			_ = ws.Close()
			continue
		}

		b.mu.Lock()
		b.ws = ws
		b.mu.Unlock()

		log.Printf("[ha] subscribed to state_changed events")

		// Read loop
		readErr := b.readLoop(ctx, ws)
		b.mu.Lock()
		b.ws = nil
		b.mu.Unlock()
		_ = ws.Close()

		if readErr != nil && ctx.Err() == nil {
			log.Printf("[ha] read loop ended: %v — reconnecting in %v", readErr, backoff)
			b.publishError(fmt.Sprintf("read loop: %v", readErr))
			select {
			case <-ctx.Done():
				return
			case <-time.After(backoff):
			}
			backoff = minDuration(time.Duration(float64(backoff)*1.5)+time.Duration(rand.Intn(1000))*time.Millisecond, maxBackoff)
		}
	}
}

func (b *bridge) readLoop(ctx context.Context, ws *websocket.Conn) error {
	// Use a channel so we can respect ctx cancellation
	type result struct {
		raw string
		err error
	}
	ch := make(chan result, 1)

	for {
		go func() {
			var raw string
			err := websocket.Message.Receive(ws, &raw)
			ch <- result{raw, err}
		}()

		select {
		case <-ctx.Done():
			return nil
		case r := <-ch:
			if r.err != nil {
				return r.err
			}
			b.processRaw(r.raw)
		}
	}
}

func (b *bridge) processRaw(raw string) {
	// Quick type sniff
	var base haBaseMsg
	if err := json.Unmarshal([]byte(raw), &base); err != nil {
		return
	}
	if base.Type != "event" {
		return // ignore result/pong/etc.
	}
	var env haEventEnvelope
	if err := json.Unmarshal([]byte(raw), &env); err != nil {
		log.Printf("[ha] unmarshal event: %v", err)
		return
	}
	b.handleEvent(env)
}

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

func minDuration(a, b time.Duration) time.Duration {
	if a < b {
		return a
	}
	return b
}

// suppress unused import warning if math is only used in backoff
var _ = math.MaxFloat64

// ---------------------------------------------------------------------------
// main
// ---------------------------------------------------------------------------

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds)
	log.SetPrefix("[homeassistant-bridge] ")

	cfg := loadConfig()
	log.Printf("starting bridge — HA=%s:%s NATS=%s", cfg.haHost, cfg.haPort, cfg.natsURL)

	b := newBridge(cfg)

	// Connect NATS
	if err := b.connectNATS(); err != nil {
		log.Fatalf("NATS: %v", err)
	}
	defer b.nc.Drain()

	// Context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// SIGTERM / SIGINT handling
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		s := <-sig
		log.Printf("received %s — shutting down", s)
		cancel()
		// Give drain time
		time.Sleep(2 * time.Second)
		close(b.done)
	}()

	// Command subscription
	sub, err := b.subscribeCommands()
	if err != nil {
		log.Fatalf("command subscription: %v", err)
	}
	defer sub.Unsubscribe()

	// Heartbeat
	b.startHeartbeat(ctx)

	// HA loop (blocks until ctx cancelled)
	b.runHA(ctx)

	// Wait for done signal
	select {
	case <-b.done:
	case <-time.After(5 * time.Second):
	}
	log.Printf("shutdown complete")
}
