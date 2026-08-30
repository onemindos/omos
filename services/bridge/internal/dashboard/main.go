// SPDX-License-Identifier: Apache-2.0

// Package main implements the omos-bridge health dashboard.
// It subscribes to NATS subjects published by omos-bridge bridges,
// tracks their liveness, error counts, and message volumes,
// and exposes an HTTP /health endpoint plus a periodic terminal table.
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sort"
	"strings"
	"sync"
	"syscall"
	"time"

	nats "github.com/nats-io/nats.go"
)

// ────────────────────────────────────────────────────────────────────────────
// Data model
// ────────────────────────────────────────────────────────────────────────────

// BridgeState holds per-bridge tracking data.
type BridgeState struct {
	Name          string
	LastHeartbeat time.Time
	ErrorCount    int64
	MsgCount      int64
}

// Status returns "online" or "offline" based on the 90 s threshold.
func (b *BridgeState) Status() string {
	if time.Since(b.LastHeartbeat) > 90*time.Second {
		return "offline"
	}
	return "online"
}

// Store is the central in-memory store protected by a single mutex.
type Store struct {
	mu      sync.RWMutex
	bridges map[string]*BridgeState // keyed by bridge name
	domains map[string]int64        // keyed by domain (e.g. "air", "home")
}

func newStore() *Store {
	return &Store{
		bridges: make(map[string]*BridgeState),
		domains: make(map[string]int64),
	}
}

// recordHeartbeat marks a bridge as seen right now.
func (s *Store) recordHeartbeat(name string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	bs := s.ensureBridge(name)
	bs.LastHeartbeat = time.Now()
}

// recordError increments a bridge's error counter.
func (s *Store) recordError(name string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	bs := s.ensureBridge(name)
	bs.ErrorCount++
}

// recordTelemetry increments the domain message counter.
// subject looks like "omos.air.sensor.telemetry" → domain = "air"
func (s *Store) recordTelemetry(subject string) {
	parts := strings.SplitN(subject, ".", 3) // ["omos","air","..."]
	if len(parts) < 2 {
		return
	}
	domain := parts[1]
	s.mu.Lock()
	defer s.mu.Unlock()
	s.domains[domain]++
}

// recordBridgeTelemetry increments per-bridge msg count when a bridge name
// can be extracted from the subject.  Currently only called from health/errors
// handlers as a secondary side-effect; for generic telemetry subjects the
// bridge name is not in the subject, so only domain counts are updated there.
func (s *Store) recordBridgeMsgCount(name string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	bs := s.ensureBridge(name)
	bs.MsgCount++
}

// ensureBridge returns or creates a BridgeState; caller must hold s.mu.
func (s *Store) ensureBridge(name string) *BridgeState {
	if _, ok := s.bridges[name]; !ok {
		s.bridges[name] = &BridgeState{Name: name}
	}
	return s.bridges[name]
}

// ────────────────────────────────────────────────────────────────────────────
// JSON snapshot types
// ────────────────────────────────────────────────────────────────────────────

type BridgeSnapshot struct {
	Name          string `json:"name"`
	Status        string `json:"status"`
	LastHeartbeat int64  `json:"last_heartbeat"` // unix ms
	ErrorCount    int64  `json:"error_count"`
	MsgCount      int64  `json:"msg_count"`
}

type DomainSnapshot struct {
	MsgCount int64 `json:"msg_count"`
}

type HealthResponse struct {
	TS           int64                     `json:"ts"` // unix ms
	Bridges      []BridgeSnapshot          `json:"bridges"`
	Domains      map[string]DomainSnapshot `json:"domains"`
	TotalMsgCount int64                    `json:"total_msg_count"`
}

// snapshot builds a HealthResponse from the current store state.
func (s *Store) snapshot() HealthResponse {
	s.mu.RLock()
	defer s.mu.RUnlock()

	resp := HealthResponse{
		TS:      time.Now().UnixMilli(),
		Domains: make(map[string]DomainSnapshot),
	}

	// bridges — sorted for deterministic output
	names := make([]string, 0, len(s.bridges))
	for n := range s.bridges {
		names = append(names, n)
	}
	sort.Strings(names)
	for _, n := range names {
		b := s.bridges[n]
		resp.Bridges = append(resp.Bridges, BridgeSnapshot{
			Name:          b.Name,
			Status:        b.Status(),
			LastHeartbeat: b.LastHeartbeat.UnixMilli(),
			ErrorCount:    b.ErrorCount,
			MsgCount:      b.MsgCount,
		})
	}

	// domains
	var total int64
	for d, c := range s.domains {
		resp.Domains[d] = DomainSnapshot{MsgCount: c}
		total += c
	}
	resp.TotalMsgCount = total

	return resp
}

// ────────────────────────────────────────────────────────────────────────────
// Terminal table
// ────────────────────────────────────────────────────────────────────────────

func formatAge(t time.Time) string {
	if t.IsZero() {
		return "never"
	}
	d := time.Since(t)
	switch {
	case d < time.Minute:
		return fmt.Sprintf("%ds ago", int(d.Seconds()))
	case d < time.Hour:
		return fmt.Sprintf("%dm ago", int(d.Minutes()))
	default:
		return fmt.Sprintf("%dh ago", int(d.Hours()))
	}
}

func printTable(s *Store) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// clear screen
	fmt.Print("\033[2J\033[H")

	now := time.Now().Format("2006-01-02 15:04:05")
	fmt.Printf("=== omos-bridge Dashboard ===  %s\n", now)
	fmt.Printf("%-24s %-9s %-14s %-7s %s\n", "BRIDGE", "STATUS", "LAST HB", "ERRORS", "MSGS")
	fmt.Println(strings.Repeat("-", 65))

	names := make([]string, 0, len(s.bridges))
	for n := range s.bridges {
		names = append(names, n)
	}
	sort.Strings(names)
	for _, n := range names {
		b := s.bridges[n]
		status := strings.ToUpper(b.Status())
		fmt.Printf("%-24s %-9s %-14s %-7d %d\n",
			b.Name, status, formatAge(b.LastHeartbeat), b.ErrorCount, b.MsgCount)
	}

	fmt.Println()
	fmt.Printf("%-9s %s\n", "DOMAIN", "MSGS")
	fmt.Println(strings.Repeat("-", 20))

	domains := make([]string, 0, len(s.domains))
	for d := range s.domains {
		domains = append(domains, d)
	}
	sort.Strings(domains)
	var total int64
	for _, d := range domains {
		c := s.domains[d]
		fmt.Printf("%-9s %d\n", d, c)
		total += c
	}
	fmt.Printf("\nTotal: %d messages\n", total)
}

// ────────────────────────────────────────────────────────────────────────────
// main
// ────────────────────────────────────────────────────────────────────────────

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func main() {
	natsURL := getEnv("NATS_URL", nats.DefaultURL)
	dashPort := getEnv("DASHBOARD_PORT", "8099")

	store := newStore()

	// ── NATS connection ────────────────────────────────────────────────────
	nc, err := nats.Connect(natsURL,
		nats.MaxReconnects(-1),
		nats.ReconnectWait(2*time.Second),
		nats.DisconnectErrHandler(func(nc *nats.Conn, err error) {
			log.Printf("[NATS] disconnected: %v", err)
		}),
		nats.ReconnectHandler(func(nc *nats.Conn) {
			log.Printf("[NATS] reconnected to %s", nc.ConnectedUrl())
		}),
		nats.ErrorHandler(func(nc *nats.Conn, sub *nats.Subscription, err error) {
			log.Printf("[NATS] error on %s: %v", sub.Subject, err)
		}),
	)
	if err != nil {
		log.Fatalf("NATS connect to %s failed: %v", natsURL, err)
	}
	defer nc.Drain()
	log.Printf("[NATS] connected to %s", nc.ConnectedUrl())

	// ── Subscription: omos.bridge.*.health ────────────────────────────────
	// Subject format: omos.bridge.<bridge-name>.health
	_, err = nc.Subscribe("omos.bridge.*.health", func(msg *nats.Msg) {
		parts := strings.Split(msg.Subject, ".")
		if len(parts) < 4 {
			return
		}
		bridgeName := parts[2]
		store.recordHeartbeat(bridgeName)
		store.recordBridgeMsgCount(bridgeName)
	})
	if err != nil {
		log.Fatalf("subscribe omos.bridge.*.health: %v", err)
	}

	// ── Subscription: omos.bridge.*.errors ────────────────────────────────
	_, err = nc.Subscribe("omos.bridge.*.errors", func(msg *nats.Msg) {
		parts := strings.Split(msg.Subject, ".")
		if len(parts) < 4 {
			return
		}
		bridgeName := parts[2]
		store.recordError(bridgeName)
		store.recordBridgeMsgCount(bridgeName)
	})
	if err != nil {
		log.Fatalf("subscribe omos.bridge.*.errors: %v", err)
	}

	// ── Subscription: omos.>.telemetry ────────────────────────────────────
	// Matches any subject under omos.* that ends with a token "telemetry",
	// e.g. omos.air.sensor.telemetry, omos.home.device.telemetry
	_, err = nc.Subscribe("omos.>", func(msg *nats.Msg) {
		// Only count subjects ending in .telemetry to avoid double-counting
		// health / errors subjects (those are handled above).
		if !strings.HasSuffix(msg.Subject, ".telemetry") {
			return
		}
		store.recordTelemetry(msg.Subject)
	})
	if err != nil {
		log.Fatalf("subscribe omos.>: %v", err)
	}

	// ── HTTP server ────────────────────────────────────────────────────────
	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		snap := store.snapshot()
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(snap); err != nil {
			log.Printf("[HTTP] encode error: %v", err)
		}
	})

	srv := &http.Server{
		Addr:    ":" + dashPort,
		Handler: mux,
	}
	go func() {
		log.Printf("[HTTP] /health listening on :%s", dashPort)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("[HTTP] server error: %v", err)
		}
	}()

	// ── Terminal table ticker (every 10 s) ────────────────────────────────
	tableTicker := time.NewTicker(10 * time.Second)
	defer tableTicker.Stop()
	// Print once immediately on startup
	printTable(store)

	// ── Graceful shutdown ─────────────────────────────────────────────────
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	for {
		select {
		case <-tableTicker.C:
			printTable(store)
		case sig := <-quit:
			log.Printf("[main] received %s — shutting down", sig)
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			if err := srv.Shutdown(ctx); err != nil {
				log.Printf("[HTTP] shutdown error: %v", err)
			}
			return
		}
	}
}
