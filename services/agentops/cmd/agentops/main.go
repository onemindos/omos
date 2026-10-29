// Command agentops runs the OneMind agent operations components (FEATURES I1, I3, D10):
// agent presence on the map, the ATAK GeoChat path to agents, the watch-job runner and the
// AOI briefer. Configuration is environment only; see services/agentops/README.md.
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/nats-io/nats.go"

	"github.com/onemindos/omos/services/agentops/internal/agentproto"
	"github.com/onemindos/omos/services/agentops/internal/brief"
	"github.com/onemindos/omos/services/agentops/internal/bus"
	"github.com/onemindos/omos/services/agentops/internal/chat"
	"github.com/onemindos/omos/services/agentops/internal/feeds"
	"github.com/onemindos/omos/services/agentops/internal/model"
	"github.com/onemindos/omos/services/agentops/internal/presence"
	"github.com/onemindos/omos/services/agentops/internal/watch"
)

func env(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}

func envDur(k string, def time.Duration) time.Duration {
	if v := os.Getenv(k); v != "" {
		if d, err := time.ParseDuration(v); err == nil {
			return d
		}
	}
	return def
}

func envJSON(k string, into any) error {
	v := os.Getenv(k)
	if v == "" {
		return nil
	}
	if err := json.Unmarshal([]byte(v), into); err != nil {
		return fmt.Errorf("%s: %w", k, err)
	}
	return nil
}

func splitRef(ref string) (string, string, string, error) {
	p := strings.Split(ref, "/")
	if len(p) != 3 || p[0] == "" || p[1] == "" || p[2] == "" {
		return "", "", "", fmt.Errorf("agent ref %q must be agent/owner/name", ref)
	}
	return p[0], p[1], p[2], nil
}

// parseModel reads "none", "stub", "agent:{prompt subject}" or "openai:{base url}|{model}".
func parseModel(spec string, nc *nats.Conn, allowRemote bool) (model.Model, error) {
	switch {
	case spec == "" || spec == "none":
		return nil, nil
	case spec == "stub":
		return model.Stub{}, nil
	case strings.HasPrefix(spec, "agent:"):
		subj := strings.TrimPrefix(spec, "agent:")
		if !strings.HasPrefix(subj, "agents.prompt.") {
			return nil, fmt.Errorf("agent model %q must be an agents.prompt.* subject", subj)
		}
		return &model.Agent{NC: nc, Subject: subj}, nil
	case strings.HasPrefix(spec, "openai:"):
		base, name, ok := strings.Cut(strings.TrimPrefix(spec, "openai:"), "|")
		if !ok || base == "" || name == "" {
			return nil, fmt.Errorf("openai model %q must be openai:{base url}|{model}", spec)
		}
		return &model.OpenAICompat{BaseURL: base, ModelName: name, AllowRemote: allowRemote}, nil
	}
	return nil, fmt.Errorf("unknown model spec %q", spec)
}

func main() {
	log := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	if err := run(log); err != nil {
		log.Error("agentops stopped", "err", err)
		os.Exit(1)
	}
}

func run(log *slog.Logger) error {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	comps := map[string]bool{}
	for _, c := range strings.Split(env("COMPONENTS", "presence,chat,watch,brief"), ",") {
		if c = strings.TrimSpace(c); c != "" {
			comps[c] = true
		}
	}
	nc, err := bus.Connect(env("NATS_URL", "nats://127.0.0.1:4222"), os.Getenv("NATS_USER"), os.Getenv("NATS_PASSWORD"), "omos-agentops")
	if err != nil {
		return err
	}
	defer nc.Drain()
	pub, err := bus.NewPublisher(nc, env("PUBLISH_MODE", bus.ModeAuto))
	if err != nil {
		return err
	}

	allowRemote := env("MODEL_ALLOW_REMOTE", "false") == "true"
	var models []model.Model
	for _, k := range []string{"MODEL_PRIMARY", "MODEL_LOCAL"} {
		m, err := parseModel(os.Getenv(k), nc, allowRemote)
		if err != nil {
			return err
		}
		if m != nil {
			models = append(models, m)
		}
	}
	models = append(models, model.Stub{}) // always last: a brief is written with the network off
	chain := model.Chain{Models: models, Timeout: envDur("MODEL_TIMEOUT", 60*time.Second)}

	expected := map[string]float64{}
	layers := map[string][]string{}
	if err := envJSON("EXPECTED_FEEDS", &expected); err != nil {
		return err
	}
	if err := envJSON("FEED_LAYERS", &layers); err != nil {
		return err
	}
	board := feeds.NewBoard(nil, expected, layers)

	errc := make(chan error, 8)
	metrics := map[string]func() int64{
		"agentops_publish_stored_total":   pub.Stored.Load,
		"agentops_publish_core_total":     pub.Core.Load,
		"agentops_publish_failures_total": pub.Failures.Load,
	}

	var tracker *presence.Tracker
	if comps["presence"] || comps["chat"] {
		var cfg presence.Config
		cfg.NodeID = env("NODE_ID", "")
		cfg.PublishEvery = envDur("PRESENCE_EVERY", 30*time.Second)
		if err := envJSON("AGENT_ANCHORS", &cfg.Anchors); err != nil {
			return err
		}
		if err := envJSON("AGENT_CALLSIGNS", &cfg.Callsigns); err != nil {
			return err
		}
		lat, e1 := strconv.ParseFloat(os.Getenv("SITE_LAT"), 64)
		lon, e2 := strconv.ParseFloat(os.Getenv("SITE_LON"), 64)
		if e1 != nil || e2 != nil {
			if comps["presence"] {
				return fmt.Errorf("presence needs SITE_LAT and SITE_LON (the anchor where agents are drawn)")
			}
		}
		cfg.Site = presence.Anchor{Lat: lat, Lon: lon}
		tracker = presence.New(cfg, nil)
		if comps["presence"] {
			go func() { errc <- tracker.Run(ctx, nc, pub, env("AGENTS_STREAM", "AGENTS")) }()
		} else {
			// chat without presence still needs the agent list
			go func() {
				_, _ = bus.Follow(ctx, nc, pub.JS, env("AGENTS_STREAM", "AGENTS"), []string{"agents.hb.>"}, true,
					func(s string, d []byte, _ nats.Header) { _, _ = tracker.OnHeartbeat(s, d) })
			}()
		}
	}

	if comps["chat"] {
		targets := func(k string) (bool, bool) {
			v := os.Getenv(k)
			return strings.Contains(v, "server"), strings.Contains(v, "mesh")
		}
		rs, rm := targets("REPLY_TARGETS")
		cs, cm := targets("CONTACT_TARGETS")
		var agents []string
		for _, a := range strings.Split(env("CHAT_AGENTS", "hermes/zeus/legacy"), ",") {
			if a = strings.TrimSpace(a); a != "" {
				if _, _, _, err := splitRef(a); err != nil {
					return err
				}
				agents = append(agents, a)
			}
		}
		gw := chat.New(chat.Config{ReplyServer: rs, ReplyMesh: rm, ContactServer: cs, ContactMesh: cm, Agents: agents,
			ContactEvery: envDur("CONTACT_EVERY", 60*time.Second), CoTType: env("AGENT_COT_TYPE", "a-f-G-U-C"),
			Team: env("AGENT_TAK_TEAM", "Cyan"), Role: env("AGENT_TAK_ROLE", "HQ"), Version: version()}, nc, tracker, log)
		metrics["agentops_chat_received_total"] = gw.C.Received.Load
		metrics["agentops_chat_forwarded_total"] = gw.C.Forwarded.Load
		metrics["agentops_chat_replied_total"] = gw.C.Replied.Load
		metrics["agentops_chat_rate_limited_total"] = gw.C.RateLimited.Load
		metrics["agentops_chat_agent_errors_total"] = gw.C.AgentErrors.Load
		metrics["agentops_chat_reply_disabled_total"] = gw.C.ReplyDisabled.Load
		go func() { errc <- gw.Run(ctx) }()
	}

	feedFollowed := false
	followFeeds := func() error {
		if feedFollowed {
			return nil
		}
		feedFollowed = true
		_, err := bus.Follow(ctx, nc, pub.JS, env("TEL_STREAM", "TEL"), []string{"tel.feed.*.status"}, true,
			func(s string, d []byte, _ nats.Header) { board.Update(s, d) })
		return err
	}

	if comps["watch"] {
		path := env("WATCH_FILE", "/config/watches.yaml")
		ws, err := watch.Load(path)
		if err != nil {
			return err
		}
		ref := env("WATCH_AGENT", "onemind-watch/zeus/hub")
		a, o, n, err := splitRef(ref)
		if err != nil {
			return err
		}
		beater := agentproto.NewBeater(nc, a, o, n, "Watch", 30*time.Second)
		go func() { _ = beater.Run(ctx) }()
		eng := watch.NewEngine(ws, board, nil)
		r := watch.NewRunner(watch.Config{AgentID: "agent:" + ref, ENTStream: env("ENT_STREAM", "ENT"),
			TELStream: env("TEL_STREAM", "TEL"), ZoneBucket: env("ZONE_BUCKET", "geofences")}, nc, pub, eng, chain, beater, log)
		metrics["agentops_watch_alerts_total"] = r.C.Alerts.Load
		metrics["agentops_watch_geofence_events_total"] = r.C.GeofenceEvents.Load
		metrics["agentops_watch_proposed_total"] = r.C.Proposed.Load
		metrics["agentops_watch_propose_denied_total"] = r.C.ProposeDenied.Load
		metrics["agentops_watch_propose_errors_total"] = r.C.ProposeErrors.Load
		metrics["agentops_watch_model_failed_total"] = r.C.ModelFailed.Load
		metrics["agentops_watch_llm_limited_total"] = r.C.LLMLimited.Load
		if err := followFeeds(); err != nil {
			return err
		}
		go func() { errc <- r.Run(ctx, nil) }()
		log.Info("watch runner", "jobs", len(eng.Watches()), "file", path)
	}

	if comps["brief"] {
		var aois []brief.AOI
		if err := envJSON("BRIEF_AOIS", &aois); err != nil {
			return err
		}
		ref := env("BRIEF_AGENT", "onemind-brief/zeus/hub")
		a, o, n, err := splitRef(ref)
		if err != nil {
			return err
		}
		beater := agentproto.NewBeater(nc, a, o, n, "Brief", 30*time.Second)
		go func() { _ = beater.Run(ctx) }()
		br := brief.New(brief.Config{AOIs: aois, Every: envDur("BRIEF_EVERY", 15*time.Minute)}, nc, pub, chain, board, beater, log)
		metrics["agentops_brief_total"] = br.C.Briefs.Load
		metrics["agentops_brief_offline_total"] = br.C.Offline.Load
		metrics["agentops_brief_repaired_total"] = br.C.Repaired.Load
		metrics["agentops_brief_unknown_ids_total"] = br.C.UnknownIDs.Load
		if err := followFeeds(); err != nil {
			return err
		}
		go func() { errc <- br.Serve(ctx) }()
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		if !nc.IsConnected() {
			http.Error(w, "nats disconnected", http.StatusServiceUnavailable)
			return
		}
		_, _ = w.Write([]byte("ok\n"))
	})
	mux.HandleFunc("/metrics", func(w http.ResponseWriter, _ *http.Request) {
		for k, f := range metrics {
			fmt.Fprintf(w, "%s %d\n", k, f())
		}
	})
	srv := &http.Server{Addr: env("HTTP_ADDR", ":8080"), Handler: mux, ReadHeaderTimeout: 5 * time.Second}
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errc <- err
		}
	}()
	log.Info("agentops ready", "components", env("COMPONENTS", "presence,chat,watch,brief"), "models", len(models))
	select {
	case <-ctx.Done():
	case err := <-errc:
		if err != nil {
			stop()
			_ = srv.Shutdown(context.Background())
			return err
		}
	}
	_ = srv.Shutdown(context.Background())
	return nil
}

func version() string {
	if b, err := os.ReadFile("/VERSION"); err == nil {
		return strings.TrimSpace(string(b))
	}
	return "dev"
}
