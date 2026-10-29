// omos-ingest runs the collectors named in a feed registry (sources.yaml) and delivers
// their records into NATS JetStream with the OneMind reliability list: one owner per
// source (KV lease), acked publishes with Nats-Msg-Id, checkpoints after acks, JSON Schema
// validation with dead letters, polite fetching, and freshness on tel.feed.{id}.status and
// Prometheus. See ../../README.md.
package main

import (
	"context"
	"crypto/tls"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/onemindos/omos-bridge/ingest/internal/ingest"
	_ "github.com/onemindos/omos-bridge/ingest/internal/parsers"
)

var version = "dev"

func main() {
	cfgPath := flag.String("config", envOr("INGEST_CONFIG", "/etc/omos-ingest/sources.yaml"), "feed registry file")
	schemas := flag.String("schemas", envOr("INGEST_SCHEMAS", "/schemas"), "JSON Schema directory (services/schemas)")
	natsURL := flag.String("nats", envOr("NATS_URL", nats.DefaultURL), "NATS URL")
	instance := flag.String("instance", envOr("POD_NAME", hostname()), "instance id used in leases")
	addr := flag.String("listen", envOr("INGEST_LISTEN", ":9464"), "metrics and health listen address")
	only := flag.String("only", "", "comma-separated source ids to run (default: all enabled)")
	once := flag.Bool("once", false, "poll sources: stop after one successful fetch (bulk CronJob mode)")
	createBuckets := flag.Bool("create-buckets", false, "create the lease and checkpoint KV buckets if missing (local runs only; the cluster uses NACK CRs)")
	check := flag.Bool("check", false, "load and validate the registry and schemas, then exit")
	flag.Parse()

	log := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})).With("app", "omos-ingest", "instance", *instance)
	cfg, err := ingest.LoadConfig(*cfgPath)
	if err != nil {
		log.Error("config", "err", err.Error())
		os.Exit(2)
	}

	board := ingest.NewStatusBoard(*instance, nil)
	rt := &ingest.Runtime{Board: board, Metrics: ingest.NewMetrics(board), Schemas: ingest.NewSchemas(*schemas),
		Log: log, Instance: *instance, LeaseTTL: cfg.LeaseTTL.Duration}

	var sources []*ingest.Source
	for _, sc := range cfg.Sources {
		if !sc.IsEnabled() || (*only != "" && !strings.Contains(","+*only+",", ","+sc.ID+",")) {
			continue
		}
		s, err := ingest.NewSource(rt, sc)
		if err != nil {
			log.Error("source", "err", err.Error())
			os.Exit(2)
		}
		s.Once = *once
		sources = append(sources, s)
	}
	if *check {
		if len(sources) == 0 {
			fmt.Fprintf(os.Stderr, "check failed: no enabled sources in %s (%d entries)\n", *cfgPath, len(cfg.Sources))
			os.Exit(1)
		}
		fmt.Printf("ok: %d sources valid\n", len(sources))
		return
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	ready := make(chan struct{})
	go serve(*addr, rt, ready, log)

	opts := []nats.Option{nats.Name("omos-ingest-" + *instance), nats.RetryOnFailedConnect(true),
		nats.MaxReconnects(-1), nats.ReconnectWait(time.Second),
		// No reconnect buffer: a publish while disconnected fails at once and Publish retries
		// it. A buffered copy would be sent on reconnect even if the lease was lost meanwhile.
		nats.ReconnectBufSize(-1),
		nats.DisconnectErrHandler(func(_ *nats.Conn, err error) { log.Warn("nats disconnected", "err", fmt.Sprint(err)) }),
		nats.ReconnectHandler(func(_ *nats.Conn) { log.Info("nats reconnected") })}
	if u, p := os.Getenv("NATS_USER"), os.Getenv("NATS_PASSWORD"); u != "" {
		opts = append(opts, nats.UserInfo(u, p))
	}
	if c := os.Getenv("NATS_CREDS"); c != "" {
		opts = append(opts, nats.UserCredentials(c))
	}
	// The hub's client port requires TLS (platform/nats/values.template.yaml config.nats.tls).
	// NATS_CA verifies it; NATS_TLS_INSECURE=true matches the current omos Python feeds
	// (ssl.CERT_NONE) until the NATS CA is available in the workload's namespace.
	if ca := os.Getenv("NATS_CA"); ca != "" {
		opts = append(opts, nats.RootCAs(ca))
	} else if os.Getenv("NATS_TLS_INSECURE") == "true" {
		log.Warn("NATS TLS certificate verification disabled (NATS_TLS_INSECURE=true)")
		opts = append(opts, nats.Secure(&tls.Config{InsecureSkipVerify: true, MinVersion: tls.VersionTLS12})) //nolint:gosec // documented interim posture
	}
	nc, err := nats.Connect(*natsURL, opts...)
	if err != nil {
		log.Error("nats connect", "err", err.Error())
		os.Exit(1)
	}
	defer nc.Close()
	rt.NC = nc
	if rt.JS, err = jetstream.New(nc); err != nil {
		log.Error("jetstream", "err", err.Error())
		os.Exit(1)
	}
	// Wait for the buckets (NATS may still be down at start: keep trying, never exit).
	for {
		bctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		err = rt.BindBuckets(bctx, cfg, *createBuckets)
		cancel()
		if err == nil {
			break
		}
		log.Warn("waiting for KV buckets", "err", err.Error())
		if ctx.Err() != nil {
			return
		}
		time.Sleep(2 * time.Second)
	}
	close(ready)
	log.Info("omos-ingest up", "version", version, "sources", len(sources))

	go rt.Heartbeat(ctx, cfg.Status.Duration)
	var wg sync.WaitGroup
	for _, s := range sources {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if *once {
				// CronJob: one owner is guaranteed by concurrencyPolicy: Forbid.
				rt.Board.SetLeader(s.Config().ID, true)
				s.Run(ctx)
				return
			}
			rt.RunOwned(ctx, s.Config().ID, s.Run)
		}()
	}
	wg.Wait()
	rt.PublishStatus()
	_ = nc.Drain()
	log.Info("omos-ingest stopped")
}

func serve(addr string, rt *ingest.Runtime, ready <-chan struct{}, log *slog.Logger) {
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.HandlerFor(rt.Metrics.Reg, promhttp.HandlerOpts{}))
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, _ *http.Request) { _, _ = w.Write([]byte("ok\n")) })
	mux.HandleFunc("/readyz", func(w http.ResponseWriter, _ *http.Request) {
		select {
		case <-ready:
			if rt.NC != nil && rt.NC.IsConnected() {
				_, _ = w.Write([]byte("ok\n"))
				return
			}
		default:
		}
		http.Error(w, "not ready", http.StatusServiceUnavailable)
	})
	srv := &http.Server{Addr: addr, Handler: mux, ReadHeaderTimeout: 5 * time.Second}
	if err := srv.ListenAndServe(); err != nil {
		log.Error("listen", "err", err.Error())
	}
}

func envOr(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}

func hostname() string {
	h, err := os.Hostname()
	if err != nil {
		return "omos-ingest"
	}
	return h
}
