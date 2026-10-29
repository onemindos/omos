// Command mission runs the OneMind task service (FEATURES D2 and H4).
//
// It applies cmd.mission.{id}.{verb} from the policy service, is the only writer of KV
// mission-state, republishes every change as evt.mission.{id}.{event}, and serves
// svc.mission.{get,list,report,apply}. Configuration is environment only; the NATS
// password comes from a Secret, never from Git.
package main

import (
	"context"
	"errors"
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
	"github.com/nats-io/nats.go/jetstream"

	"github.com/onemindos/omos/services/mission/internal/mission"
)

func env(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}

func main() {
	log := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	if err := run(log); err != nil {
		log.Error("task service stopped", "err", err)
		os.Exit(1)
	}
}

func run(log *slog.Logger) error {
	opts := []nats.Option{nats.Name("mission-service"), nats.MaxReconnects(-1)}
	if u := os.Getenv("NATS_USER"); u != "" {
		opts = append(opts, nats.UserInfo(u, os.Getenv("NATS_PASSWORD")))
	}
	if ca := os.Getenv("NATS_CA_FILE"); ca != "" {
		opts = append(opts, nats.RootCAs(ca))
	}
	nc, err := nats.Connect(env("NATS_URL", "nats://127.0.0.1:4222"), opts...)
	if err != nil {
		return err
	}
	defer nc.Drain()
	js, err := jetstream.New(nc)
	if err != nil {
		return err
	}
	var dir *mission.Directory
	if p := env("MISSION_PRINCIPALS_FILE", "/config/principals.json"); p != "" {
		if dir, err = mission.LoadDirectory(p); err != nil {
			if !errors.Is(err, os.ErrNotExist) {
				return err
			}
			log.Warn("no principals file: svc.mission.report refuses every caller", "path", p)
			dir = nil
		}
	}
	radius, err := strconv.ParseFloat(env("MISSION_ARRIVAL_RADIUS_M", "5"), 64)
	if err != nil {
		return fmt.Errorf("MISSION_ARRIVAL_RADIUS_M: %w", err)
	}
	vert, err := strconv.ParseFloat(env("MISSION_ARRIVAL_VERTICAL_M", "5"), 64)
	if err != nil {
		return fmt.Errorf("MISSION_ARRIVAL_VERTICAL_M: %w", err)
	}
	sweep, err := time.ParseDuration(env("MISSION_SWEEP_EVERY", "1s"))
	if err != nil {
		return fmt.Errorf("MISSION_SWEEP_EVERY: %w", err)
	}
	svc, err := mission.New(mission.Config{
		Conn: nc, JS: js, Log: log, Source: env("MISSION_SOURCE", "service:task"),
		Bucket: env("MISSION_BUCKET", "mission-state"), CmdStream: env("MISSION_CMD_STREAM", "CMD"),
		EvtStream: env("MISSION_EVT_STREAM", "EVT"), ConsumerName: env("MISSION_CONSUMER", "mission"),
		ReportClasses: strings.Split(env("MISSION_REPORT_CLASSES", "drone,robot"), ","),
		SweepEvery:    sweep, Arrival: env("MISSION_ARRIVAL", "true") == "true",
		ArrivalRadiusM: radius, ArrivalVerticalM: vert, Directory: dir,
	})
	if err != nil {
		return err
	}
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	if err := svc.Start(ctx); err != nil {
		return err
	}
	defer svc.Stop()

	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		if !nc.IsConnected() {
			http.Error(w, "nats disconnected", http.StatusServiceUnavailable)
			return
		}
		_, _ = w.Write([]byte("ok"))
	})
	mux.HandleFunc("/metrics", func(w http.ResponseWriter, _ *http.Request) {
		c := &svc.C
		for _, m := range []struct {
			n string
			v int64
		}{
			{"mission_commands_applied_total", c.CommandsApplied.Load()},
			{"mission_commands_rejected_total", c.CommandsRejected.Load()},
			{"mission_reports_accepted_total", c.ReportsAccepted.Load()},
			{"mission_reports_dropped_stale_total", c.ReportsStale.Load()},
			{"mission_reports_dropped_terminal_total", c.ReportsTerminal.Load()},
			{"mission_reports_refused_total", c.ReportsRefused.Load()},
			{"mission_task_timeouts_total", c.Timeouts.Load()},
			{"mission_task_arrivals_total", c.Arrivals.Load()},
			{"mission_fact_publish_errors_total", c.FactErrors.Load()},
		} {
			fmt.Fprintf(w, "# TYPE %s counter\n%s %d\n", m.n, m.n, m.v)
		}
	})
	srv := &http.Server{Addr: env("HEALTH_ADDR", ":8080"), Handler: mux, ReadHeaderTimeout: 5 * time.Second}
	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error("health server", "err", err)
		}
	}()
	<-ctx.Done()
	sctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return srv.Shutdown(sctx)
}
