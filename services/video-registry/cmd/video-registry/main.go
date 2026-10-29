// video-registry links MediaMTX streams to OneMind entities (FEATURES K2). It polls the
// MediaMTX control API, writes one record per stream to KV OVERWATCH_VIDEO (key
// {entity_type}.{entity_id}), publishes evt.video.{entity_id}.online|offline and serves
// svc.registry.video.get|list.
//
// Usage:
//
//	video-registry                 run the service (configuration from the environment)
//	video-registry select ENTITY   ask a running registry for an entity's stream, e.g.
//	                               video-registry select ent.drone.sim-x1.state
package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/nats-io/nats.go"

	"github.com/onemindos/omos/services/video-registry/internal/registry"
)

func env(k, d string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return d
}

func envInt(k string, d int) int {
	if n, err := strconv.Atoi(os.Getenv(k)); err == nil {
		return n
	}
	return d
}

func envDur(k string, d time.Duration) time.Duration {
	if v, err := time.ParseDuration(os.Getenv(k)); err == nil {
		return v
	}
	return d
}

func connect(name string) (*nats.Conn, error) {
	opts := []nats.Option{nats.Name(name), nats.MaxReconnects(-1), nats.ReconnectWait(2 * time.Second)}
	if u := os.Getenv("NATS_USER"); u != "" {
		opts = append(opts, nats.UserInfo(u, os.Getenv("NATS_PASSWORD")))
	}
	if c := os.Getenv("NATS_CREDS"); c != "" {
		opts = append(opts, nats.UserCredentials(c))
	}
	return nats.Connect(env("NATS_URL", nats.DefaultURL), opts...)
}

func main() {
	log := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	if len(os.Args) == 3 && os.Args[1] == "select" {
		os.Exit(selectEntity(os.Args[2]))
	}
	if err := run(log); err != nil {
		log.Error("video-registry stopped", "err", err)
		os.Exit(1)
	}
}

func selectEntity(entity string) int {
	nc, err := connect("video-registry-select")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
	defer nc.Close()
	req, _ := json.Marshal(registry.GetRequest{Entity: entity})
	msg, err := nc.Request(registry.SubjectGet, req, 3*time.Second)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
	fmt.Println(string(msg.Data))
	var r registry.Reply
	if json.Unmarshal(msg.Data, &r) != nil || r.Record == nil {
		return 2
	}
	return 0
}

func run(log *slog.Logger) error {
	cfgPath := env("MEDIAMTX_CONFIG", "/config/mediamtx.yml")
	media, err := registry.LoadMediaConfig(cfgPath)
	if err != nil {
		return err
	}
	pass := os.Getenv("MEDIAMTX_PASSWORD")
	if pass == "" {
		return errors.New("MEDIAMTX_PASSWORD is required (onemind-media-users Secret, key video-registry)")
	}
	nc, err := connect("video-registry")
	if err != nil {
		return err
	}
	defer nc.Close()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	svc, err := registry.New(ctx, nc, &registry.MediaMTXClient{
		BaseURL: env("MEDIAMTX_API", "http://onemind-media.onetak.svc.cluster.local:9997"),
		User:    env("MEDIAMTX_USER", "video-registry"), Password: pass,
	}, registry.Config{
		Bucket:       env("KV_BUCKET", registry.DefaultKV),
		Origin:       env("ORIGIN", "hub"),
		PollInterval: envDur("POLL_INTERVAL", 2*time.Second),
		StaleAfter:   envDur("STALE_AFTER", 30*time.Second),
		Media:        media,
		Logger:       log,
		Endpoints: registry.Endpoints{
			Host: os.Getenv("MEDIA_HOST"), SiteHost: os.Getenv("SITE_MEDIA_HOST"),
			RTSPPort: envInt("MEDIA_RTSP_PORT", 8554), HLSPort: envInt("MEDIA_HLS_PORT", 8888),
			WebRTC: envInt("MEDIA_WEBRTC_PORT", 8889), SRTPort: envInt("MEDIA_SRT_PORT", 8890),
			Playback: envInt("MEDIA_PLAYBACK_PORT", 9996),
		},
	})
	if err != nil {
		return err
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		if !nc.IsConnected() {
			http.Error(w, "nats disconnected", http.StatusServiceUnavailable)
			return
		}
		_, _ = w.Write([]byte("ok"))
	})
	mux.HandleFunc("/readyz", func(w http.ResponseWriter, _ *http.Request) {
		if !nc.IsConnected() || !svc.Ready() {
			http.Error(w, "not ready", http.StatusServiceUnavailable)
			return
		}
		_, _ = w.Write([]byte("ok"))
	})
	mux.HandleFunc("/metrics", func(w http.ResponseWriter, _ *http.Request) {
		st := &svc.Stats
		fmt.Fprintf(w, "video_registry_polls_total %d\nvideo_registry_poll_errors_total %d\n", st.Polls.Load(), st.PollErrors.Load())
		fmt.Fprintf(w, "video_registry_kv_puts_total %d\nvideo_registry_events_total %d\n", st.Puts.Load(), st.Events.Load())
		fmt.Fprintf(w, "video_registry_rejected_paths_total %d\nvideo_registry_rpc_total %d\n", st.Rejected.Load(), st.RPC.Load())
	})
	srv := &http.Server{Addr: env("HTTP_ADDR", ":8080"), Handler: mux, ReadHeaderTimeout: 5 * time.Second}
	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error("http server", "err", err)
		}
	}()
	defer func() {
		sctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		_ = srv.Shutdown(sctx)
	}()
	log.Info("video-registry started", "bucket", env("KV_BUCKET", registry.DefaultKV))
	return svc.Run(ctx)
}
