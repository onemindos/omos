// Command geofence-service is the OneMind geofence and alert service (FEATURES D1, D6): it
// applies policy-approved `cmd.map.geofence.{id}.{verb}` commands to KV `geofences`, watches
// `ent.*.*.state` for zone crossings (`evt.map.geofence.{id}.{enter|exit}`) and raises alerts
// for crossings, stale assets and feed outages on `evt.so.alert.*`. Configuration is
// environment only; NATS credentials come from a SealedSecret-backed Secret, the policy public
// keys from a mounted directory of PEM files.
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

	"github.com/onemindos/omos/services/geofence/internal/bus"
	"github.com/onemindos/omos/services/geofence/internal/engine"
	"github.com/onemindos/omos/services/geofence/internal/sig"
)

func env(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}

func envDur(k string, def time.Duration) time.Duration {
	v := os.Getenv(k)
	if v == "0" {
		return 0
	}
	if d, err := time.ParseDuration(v); err == nil {
		return d
	}
	return def
}

func main() {
	log := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	assets := engine.DefaultAssetTypes()
	if s := os.Getenv("ASSET_TYPES"); s != "" {
		assets = map[string]bool{}
		for _, t := range strings.Split(s, ",") {
			if t = strings.TrimSpace(t); t != "" {
				assets[t] = true
			}
		}
	}
	staleDefault, _ := strconv.Atoi(env("STALE_DEFAULT_S", "60"))
	eng := engine.New(engine.Config{
		Source:       env("SOURCE", "omos-geofence"),
		AssetTypes:   assets,
		FeedSilence:  envDur("FEED_SILENCE", 90*time.Second),
		StaleDefault: staleDefault,
	})

	cfg := bus.Config{
		Source:           env("SOURCE", "omos-geofence"),
		GeofenceBucket:   env("GEOFENCE_BUCKET", "geofences"),
		EntStream:        env("ENT_STREAM", "ENT"),
		EntCore:          env("ENT_MODE", "jetstream") == "core",
		CmdStream:        env("CMD_STREAM", "CMD"),
		CmdDurable:       env("CMD_DURABLE", "geofence-svc-cmd"),
		MissionStream:    env("MISSION_STREAM", "MISSION"),
		EvtStream:        env("EVT_STREAM", "EVT"),
		RestoreLookback:  envDur("RESTORE_LOOKBACK", 720*time.Hour),
		FeedSubject:      env("FEED_SUBJECT", "tel.feed.*.status"),
		RequireSignature: env("REQUIRE_SIGNATURE", "true") == "true",
		RequireApproval:  env("REQUIRE_APPROVAL", "true") == "true",
		TickEvery:        envDur("TICK", 5*time.Second),
		ServeRPC:         env("SERVE_RPC", "true") == "true",
		Log:              log,
	}
	if dir := os.Getenv("POLICY_KEYS_DIR"); dir != "" {
		keys, err := sig.LoadDir(dir)
		if err != nil {
			log.Error("policy public keys", "dir", dir, "err", err)
			os.Exit(2)
		}
		cfg.Keys = keys
	}
	if !cfg.RequireSignature {
		log.Warn("REQUIRE_SIGNATURE=false: zone writes are accepted without a policy signature (local runs only)")
	}

	opts := []nats.Option{nats.Name("omos-geofence-service"), nats.MaxReconnects(-1)}
	if u := os.Getenv("NATS_USER"); u != "" {
		opts = append(opts, nats.UserInfo(u, os.Getenv("NATS_PASSWORD")))
	}
	nc, err := nats.Connect(env("NATS_URL", "nats://127.0.0.1:4222"), opts...)
	if err != nil {
		log.Error("connect", "err", err)
		os.Exit(1)
	}
	defer nc.Close()

	svc, err := bus.New(nc, eng, cfg)
	if err != nil {
		log.Error("init", "err", err)
		os.Exit(2)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, _ *http.Request) { fmt.Fprintln(w, "ok") })
	mux.HandleFunc("/readyz", func(w http.ResponseWriter, _ *http.Request) {
		if !svc.Ready() {
			http.Error(w, "loading zones and entities", http.StatusServiceUnavailable)
			return
		}
		fmt.Fprintln(w, "ready")
	})
	mux.HandleFunc("/metrics", func(w http.ResponseWriter, _ *http.Request) {
		c := svc.Counters()
		keys := make([]string, 0, len(c))
		for k := range c {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			fmt.Fprintf(w, "omos_geofence_%s %d\n", strings.ReplaceAll(k, "-", "_"), c[k])
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
				log.Info("counters", "c", svc.Counters())
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
