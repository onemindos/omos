// Command layer-service is the OneMind layer service (FEATURES A5): it follows ENT, writes
// layer manifests and density cells to KV `world-state` and site status to KV `fleet-state`,
// and serves `svc.layer.{list,meta,snapshot,watch,unwatch}`. Configuration is environment
// only; credentials come from a SealedSecret-backed Secret, never from Git.
package main

import (
	"context"
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

	"github.com/onemindos/omos/services/layer/internal/layer"
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
	if v, err := strconv.Atoi(os.Getenv(k)); err == nil {
		return v
	}
	return def
}

func main() {
	log := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	host, _ := os.Hostname()

	defs, err := layer.LoadLayers(os.Getenv("LAYERS_FILE"))
	if err != nil {
		log.Error("layers", "err", err)
		os.Exit(2)
	}
	dir, err := layer.LoadDirectory(os.Getenv("PRINCIPALS_FILE"))
	if err != nil {
		log.Error("principals", "err", err)
		os.Exit(2)
	}
	opts := []nats.Option{nats.Name("omos-layer-service"), nats.MaxReconnects(-1)}
	if u := os.Getenv("NATS_USER"); u != "" {
		opts = append(opts, nats.UserInfo(u, os.Getenv("NATS_PASSWORD")))
	}
	nc, err := nats.Connect(env("NATS_URL", "nats://127.0.0.1:4222"), opts...)
	if err != nil {
		log.Error("connect", "err", err)
		os.Exit(1)
	}
	defer nc.Close()

	cfg := layer.Config{
		Layers:        defs,
		Directory:     dir,
		EntStream:     env("ENT_STREAM", "ENT"),
		TelStream:     env("TEL_STREAM", "TEL"),
		WorldBucket:   env("WORLD_BUCKET", "world-state"),
		FleetBucket:   env("FLEET_BUCKET", "fleet-state"),
		Site:          env("SITE", "hub"),
		LeaseBucket:   env("LEASE_BUCKET", "world-state"),
		LeaseKey:      env("LEASE_KEY", "_lease.layer-service"),
		LeaseTTL:      envDur("LEASE_TTL", 15*time.Second),
		Holder:        env("POD_NAME", host) + "-" + strconv.Itoa(os.Getpid()),
		ReplayWindow:  envDur("REPLAY_WINDOW", 5*time.Minute),
		ReplayMax:     envInt("REPLAY_MAX", 200000),
		MaxWatches:    envInt("MAX_WATCHES", layer.DefaultMaxWatches),
		MaxSites:      envInt("MAX_SITES", 64),
		Heartbeat:     envDur("WATCH_HEARTBEAT", 10*time.Second),
		ManifestEvery: envDur("MANIFEST_EVERY", 30*time.Second),
		CellsEvery:    envDur("CELLS_EVERY", 30*time.Second),
		SiteEvery:     envDur("SITE_EVERY", 10*time.Second),
		SiteSubject:   env("SITE_SUBJECT", "tel.fleet.node.*.status"),
	}
	if v, err := strconv.ParseInt(os.Getenv("WATCH_MEM_BUDGET_MB"), 10, 64); err == nil && v > 0 {
		cfg.WatchMemBudget = v << 20
	}
	if v := os.Getenv("SITE_ALLOWLIST"); v != "" {
		cfg.SiteAllow = strings.Split(v, ",")
	}
	if v := os.Getenv("SITE_STREAMS"); v != "" {
		cfg.SiteStreams = strings.Split(v, ",")
	}
	if cfg.LeaseBucket == "none" { // local runs only
		cfg.LeaseBucket = ""
	}
	svc, err := layer.New(nc, cfg, log)
	if err != nil {
		log.Error("init", "err", err)
		os.Exit(1)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, _ *http.Request) { fmt.Fprintln(w, "ok") })
	mux.HandleFunc("/readyz", func(w http.ResponseWriter, _ *http.Request) {
		if !svc.Ready() {
			http.Error(w, "not leader or not caught up with ENT", http.StatusServiceUnavailable)
			return
		}
		fmt.Fprintln(w, "ready")
	})
	mux.HandleFunc("/metrics", func(w http.ResponseWriter, _ *http.Request) {
		m := svc.Metrics()
		keys := make([]string, 0, len(m))
		for k := range m {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			fmt.Fprintf(w, "omos_layer_%s %d\n", k, m[k])
		}
	})
	srv := &http.Server{Addr: env("HTTP_ADDR", ":8080"), Handler: mux, ReadHeaderTimeout: 5 * time.Second}
	go func() { _ = srv.ListenAndServe() }()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	err = svc.Run(ctx)
	_ = srv.Shutdown(context.Background())
	if err != nil {
		log.Error("stopped", "err", err)
		os.Exit(1)
	}
}
