// Command entity-service is the OneMind entity service (FEATURES A4): it fuses observations
// from COT, TEL, WORLD, INTEL (and DET once it exists) into `ent.{entity_type}.{entity_id}.state`
// records in the ENT stream. Configuration is environment only; credentials come from a
// SealedSecret-backed Secret, never from Git.
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"sort"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/nats-io/nats.go"

	"github.com/onemindos/omos/services/entity/internal/bus"
	"github.com/onemindos/omos/services/entity/internal/engine"
)

func env(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}

func envDur(k string, def time.Duration) time.Duration {
	if d, err := time.ParseDuration(os.Getenv(k)); err == nil {
		return d
	}
	return def
}

func envInt(k string, def int) int {
	if n, err := strconv.Atoi(os.Getenv(k)); err == nil {
		return n
	}
	return def
}

// parseInputs reads "COT=tak.cot.>;TEL=tel.fleet.>,tel.drone.>".
func parseInputs(s string) []bus.Input {
	if s == "" {
		return bus.DefaultInputs()
	}
	var out []bus.Input
	for _, part := range strings.Split(s, ";") {
		name, filters, ok := strings.Cut(strings.TrimSpace(part), "=")
		if !ok || name == "" {
			continue
		}
		out = append(out, bus.Input{Stream: name, Filters: strings.Split(filters, ",")})
	}
	return out
}

func main() {
	log := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	host, _ := os.Hostname()

	ecfg := engine.DefaultConfig()
	ecfg.CoreMode = env("ENT_MODE", "jetstream") == "core"
	if s := os.Getenv("ENTITY_LINKS"); s != "" { // JSON {"cot_uid:ABC": "sim-x1"}
		if err := json.Unmarshal([]byte(s), &ecfg.Links); err != nil {
			log.Error("ENTITY_LINKS is not a JSON object", "err", err)
			os.Exit(2)
		}
	}
	if s := os.Getenv("SOURCE_PRIORS"); s != "" { // JSON {"feed-aircraft": 0.8}
		if err := json.Unmarshal([]byte(s), &ecfg.Priors); err != nil {
			log.Error("SOURCE_PRIORS is not a JSON object", "err", err)
			os.Exit(2)
		}
	}
	eng := engine.New(ecfg, func() int64 { return time.Now().UnixMilli() })

	opts := []nats.Option{nats.Name("omos-entity-service"), nats.MaxReconnects(-1)}
	if u := os.Getenv("NATS_USER"); u != "" {
		opts = append(opts, nats.UserInfo(u, os.Getenv("NATS_PASSWORD")))
	}
	nc, err := nats.Connect(env("NATS_URL", "nats://127.0.0.1:4222"), opts...)
	if err != nil {
		log.Error("connect", "err", err)
		os.Exit(1)
	}
	defer nc.Close()

	cfg := bus.Config{
		Inputs:        parseInputs(os.Getenv("INPUTS")),
		EntStream:     env("ENT_STREAM", "ENT"),
		CoreMode:      ecfg.CoreMode,
		MsgTTL:        env("ENT_MSG_TTL", "true") == "true",
		DurablePrefix: env("DURABLE_PREFIX", "entity-svc"),
		StartLookback: envDur("START_LOOKBACK", 10*time.Minute),
		RejectSubject: env("REJECT_SUBJECT", "evt.sys.entity.rejected"),
		LeaseBucket:   env("LEASE_BUCKET", "entity-correlation"),
		LeaseKey:      env("LEASE_KEY", "lease.entity-service"),
		LeaseTTL:      envDur("LEASE_TTL", 15*time.Second),
		Holder:        env("POD_NAME", host) + "-" + strconv.Itoa(os.Getpid()),
		ServeRPC:      env("SERVE_RPC", "true") == "true",
		OrbitEvery:    envDur("ORBIT_EVERY", 30*time.Second),
		MaxOrbits:     envInt("MAX_ORBITS", 500),
	}
	if cfg.LeaseBucket == "none" { // local runs only
		cfg.LeaseBucket = ""
	}
	if os.Getenv("START_LOOKBACK") == "0" {
		cfg.StartLookback = 0
	}
	svc, err := bus.New(nc, eng, cfg, log)
	if err != nil {
		log.Error("init", "err", err)
		os.Exit(1)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, _ *http.Request) { fmt.Fprintln(w, "ok") })
	mux.HandleFunc("/readyz", func(w http.ResponseWriter, _ *http.Request) {
		if !svc.Ready() {
			http.Error(w, "not leader or not consuming", http.StatusServiceUnavailable)
			return
		}
		fmt.Fprintln(w, "ready")
	})
	mux.HandleFunc("/metrics", func(w http.ResponseWriter, _ *http.Request) {
		c := eng.Counters()
		keys := make([]string, 0, len(c))
		for k := range c {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			fmt.Fprintf(w, "omos_entity_%s %d\n", strings.ReplaceAll(k, "-", "_"), c[k])
		}
	})
	srv := &http.Server{Addr: env("HTTP_ADDR", ":8080"), Handler: mux, ReadHeaderTimeout: 5 * time.Second}
	go func() { _ = srv.ListenAndServe() }()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	go func() {
		t := time.NewTicker(time.Minute)
		defer t.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-t.C:
				log.Info("counters", "c", eng.Counters())
			}
		}
	}()
	err = svc.Run(ctx)
	_ = srv.Shutdown(context.Background())
	if err != nil {
		log.Error("stopped", "err", err)
		os.Exit(1)
	}
}
