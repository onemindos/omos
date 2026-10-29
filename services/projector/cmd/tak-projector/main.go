// Command tak-projector projects fused OneMind entities to TAK (FEATURES A6;
// onemind-xr docs/contracts/cot-projection.md). Configuration is the environment of §10.
package main

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/nats-io/nats.go"

	"github.com/onemindos/omos/services/projector/internal/projector"
)

func env(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}

func main() {
	level := slog.LevelInfo
	if os.Getenv("LOG_LEVEL") == "debug" {
		level = slog.LevelDebug
	}
	log := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: level}))
	cfg, err := projector.FromEnv()
	if err != nil {
		log.Error("config", "err", err)
		os.Exit(2)
	}
	url := env("NATS_URL", "nats://127.0.0.1:4222")
	opts := []nats.Option{nats.Name("omos-tak-projector-" + cfg.Instance), nats.MaxReconnects(-1),
		nats.ReconnectWait(2 * time.Second),
		nats.DisconnectErrHandler(func(_ *nats.Conn, err error) { log.Warn("nats disconnected", "err", err) }),
		nats.ReconnectHandler(func(*nats.Conn) { log.Info("nats reconnected") })}
	if u := os.Getenv("NATS_USER"); u != "" {
		opts = append(opts, nats.UserInfo(u, os.Getenv("NATS_PASSWORD")))
	}
	if strings.HasPrefix(url, "tls://") && os.Getenv("NATS_TLS_INSECURE") == "true" {
		opts = append(opts, nats.Secure(&tls.Config{InsecureSkipVerify: true})) //nolint:gosec // matches the legacy projector; opt-in only
	}
	nc, err := nats.Connect(url, opts...)
	if err != nil {
		log.Error("nats connect", "err", err)
		os.Exit(1)
	}
	defer nc.Drain()
	svc, err := projector.New(cfg, nc, log)
	if err != nil {
		log.Error("service", "err", err)
		os.Exit(1)
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		if !nc.IsConnected() {
			http.Error(w, "nats not connected", http.StatusServiceUnavailable)
			return
		}
		_, _ = w.Write([]byte("ok"))
	})
	mux.HandleFunc("/readyz", func(w http.ResponseWriter, _ *http.Request) {
		if !svc.Ready() || !nc.IsConnected() {
			http.Error(w, "not ready", http.StatusServiceUnavailable)
			return
		}
		_, _ = w.Write([]byte("ready"))
	})
	mux.HandleFunc("/status", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(svc.Status())
	})
	srv := &http.Server{Addr: env("HTTP_ADDR", ":8080"), Handler: mux, ReadHeaderTimeout: 5 * time.Second}
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error("http", "err", err)
		}
	}()
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	log.Info("omos-tak-projector starting", "target", cfg.Target, "subject", cfg.Subject, "instance", cfg.Instance,
		"max_projected", cfg.MaxProjected, "max_rate", cfg.MaxRatePerSec, "release_groups", cfg.ReleaseGroups)
	if err := svc.Run(ctx); err != nil {
		log.Error("run", "err", err)
	}
	sctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_ = srv.Shutdown(sctx)
}
