// SPDX-License-Identifier: Apache-2.0
// omos-bridge — webhook bridge
// Receives HTTP webhooks and publishes them as One Mind envelopes to NATS.

package main

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/nats-io/nats.go"
)

// ---------------------------------------------------------------------------
// Config
// ---------------------------------------------------------------------------

type Config struct {
	NatsURL     string
	BridgeName  string
	Port        string
	Secret      string
}

func loadConfig() Config {
	return Config{
		NatsURL:    getEnv("NATS_URL", "nats://localhost:4222"),
		BridgeName: getEnv("BRIDGE_NAME", "webhook-bridge"),
		Port:       getEnv("WEBHOOK_PORT", "8090"),
		Secret:     getEnv("WEBHOOK_SECRET", ""),
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
// HMAC verification (GitHub-style X-Hub-Signature-256)
// ---------------------------------------------------------------------------

func verifySignature(secret string, body []byte, sigHeader string) bool {
	if !strings.HasPrefix(sigHeader, "sha256=") {
		return false
	}
	gotHex := strings.TrimPrefix(sigHeader, "sha256=")
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(body)
	expected := hex.EncodeToString(mac.Sum(nil))
	return hmac.Equal([]byte(gotHex), []byte(expected))
}

// ---------------------------------------------------------------------------
// Handler
// ---------------------------------------------------------------------------

type bridge struct {
	cfg Config
	nc  *nats.Conn
}

func (b *bridge) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Only allow supported methods
	switch r.Method {
	case http.MethodPost, http.MethodGet, http.MethodPut, http.MethodPatch:
		// ok
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract source from path: /{source}[/...]
	source := strings.TrimPrefix(r.URL.Path, "/")
	source = strings.SplitN(source, "/", 2)[0]
	if source == "" {
		http.Error(w, "missing source in path", http.StatusBadRequest)
		return
	}

	// Read body
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		b.reportError(fmt.Sprintf("failed to read body: %v", err))
		http.Error(w, "failed to read body", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	// Optional HMAC verification
	if b.cfg.Secret != "" {
		sig := r.Header.Get("X-Hub-Signature-256")
		if !verifySignature(b.cfg.Secret, bodyBytes, sig) {
			log.Printf("[webhook] HMAC verification failed for source=%s", source)
			http.Error(w, `{"error":"invalid signature"}`, http.StatusUnauthorized)
			return
		}
	}

	// Collect headers
	headers := make(map[string]string, len(r.Header))
	for k, vs := range r.Header {
		headers[k] = strings.Join(vs, ", ")
	}

	// Body as string
	bodyRaw := string(bodyBytes)

	// Try to parse body as JSON
	var bodyParsed interface{}
	if len(bodyBytes) > 0 {
		var tmp map[string]interface{}
		if json.Unmarshal(bodyBytes, &tmp) == nil {
			bodyParsed = tmp
		}
	}

	// Domain
	domain := r.Header.Get("X-Omos-Domain")
	if domain == "" {
		domain = "net"
	}

	// Full URL
	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}
	if fwdProto := r.Header.Get("X-Forwarded-Proto"); fwdProto != "" {
		scheme = fwdProto
	}
	sourceURL := fmt.Sprintf("%s://%s%s", scheme, r.Host, r.RequestURI)

	// Build envelope
	env := Envelope{
		Entity:    fmt.Sprintf("%s.%s", domain, source),
		Domain:    domain,
		MsgType:   "webhook",
		Timestamp: time.Now().UnixMilli(),
		Data: map[string]interface{}{
			"source_url":   sourceURL,
			"method":       r.Method,
			"headers":      headers,
			"body_raw":     bodyRaw,
			"body_parsed":  bodyParsed,
			"received_at":  time.Now().UnixMilli(),
		},
	}

	payload, err := json.Marshal(env)
	if err != nil {
		b.reportError(fmt.Sprintf("failed to marshal envelope: %v", err))
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	subject := fmt.Sprintf("omos.%s.%s.telemetry", domain, source)
	if err := b.nc.Publish(subject, payload); err != nil {
		b.reportError(fmt.Sprintf("NATS publish failed: %v", err))
		http.Error(w, "publish failed", http.StatusInternalServerError)
		return
	}

	log.Printf("[webhook] %s %s → %s", r.Method, r.URL.Path, subject)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"status":"ok","subject":%q}`, subject)
}

func (b *bridge) reportError(msg string) {
	log.Printf("[webhook] ERROR: %s", msg)
	payload, _ := json.Marshal(map[string]interface{}{
		"bridge":    b.cfg.BridgeName,
		"error":     msg,
		"timestamp": time.Now().UnixMilli(),
	})
	subject := fmt.Sprintf("omos.bridge.%s.errors", b.cfg.BridgeName)
	_ = b.nc.Publish(subject, payload)
}

// ---------------------------------------------------------------------------
// Heartbeat
// ---------------------------------------------------------------------------

func (b *bridge) startHeartbeat(ctx context.Context) {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			payload, _ := json.Marshal(map[string]interface{}{
				"bridge":    b.cfg.BridgeName,
				"status":    "ok",
				"timestamp": time.Now().UnixMilli(),
			})
			subject := fmt.Sprintf("omos.bridge.%s.health", b.cfg.BridgeName)
			if err := b.nc.Publish(subject, payload); err != nil {
				log.Printf("[webhook] heartbeat publish failed: %v", err)
			} else {
				log.Printf("[webhook] heartbeat → %s", subject)
			}
		}
	}
}

// ---------------------------------------------------------------------------
// Main
// ---------------------------------------------------------------------------

func main() {
	cfg := loadConfig()

	// Connect to NATS
	nc, err := nats.Connect(cfg.NatsURL,
		nats.MaxReconnects(-1),
		nats.ReconnectWait(2*time.Second),
		nats.DisconnectErrHandler(func(_ *nats.Conn, err error) {
			log.Printf("[webhook] NATS disconnected: %v", err)
		}),
		nats.ReconnectHandler(func(_ *nats.Conn) {
			log.Printf("[webhook] NATS reconnected")
		}),
	)
	if err != nil {
		log.Fatalf("[webhook] NATS connect failed: %v", err)
	}
	defer nc.Drain()

	log.Printf("[webhook] connected to NATS at %s", cfg.NatsURL)

	b := &bridge{cfg: cfg, nc: nc}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go b.startHeartbeat(ctx)

	mux := http.NewServeMux()
	mux.Handle("/", b)

	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: mux,
	}

	// Graceful shutdown on SIGTERM / SIGINT
	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
		sig := <-quit
		log.Printf("[webhook] received signal %s — shutting down", sig)
		cancel()
		shutCtx, shutCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer shutCancel()
		_ = srv.Shutdown(shutCtx)
	}()

	log.Printf("[webhook] listening on :%s — bridge=%s", cfg.Port, cfg.BridgeName)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("[webhook] server error: %v", err)
	}
	log.Printf("[webhook] shutdown complete")
}
