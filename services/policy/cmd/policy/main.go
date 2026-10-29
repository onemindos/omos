// Command policy runs the OneMind policy and audit service (FEATURES I4).
//
// It serves svc.policy.{command,task}.{propose,confirm}, svc.policy.action.get and
// svc.policy.check, is the only publisher of cmd.* after a human confirm, and the single
// writer of evt.so.action.*. Configuration is environment only; credentials come from a
// Secret (NATS_PASSWORD, OPENFGA_TOKEN) and a mounted key file, never from Git.
package main

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/nats-io/nats.go"

	"github.com/onemindos/omos/services/policy/internal/policy"
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
		log.Error("policy service stopped", "err", err)
		os.Exit(1)
	}
}

func run(log *slog.Logger) error {
	opts := []nats.Option{nats.Name("policy-service"), nats.MaxReconnects(-1)}
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
	js, err := nc.JetStream()
	if err != nil {
		return err
	}
	schemas, err := policy.LoadSchemas(env("POLICY_SCHEMA_DIR", "/schemas"))
	if err != nil {
		return err
	}
	dir, err := policy.LoadDirectory(env("POLICY_PRINCIPALS_FILE", "/config/principals.json"))
	if err != nil {
		return err
	}

	var authz policy.Authorizer
	switch env("POLICY_AUTHZ", "openfga") {
	case "openfga":
		if os.Getenv("OPENFGA_URL") == "" || os.Getenv("OPENFGA_STORE_ID") == "" {
			return errors.New("OPENFGA_URL and OPENFGA_STORE_ID are required (POLICY_AUTHZ=openfga)")
		}
		authz = &policy.OpenFGAClient{BaseURL: os.Getenv("OPENFGA_URL"), Store: os.Getenv("OPENFGA_STORE_ID"),
			Model: os.Getenv("OPENFGA_MODEL_ID"), Token: os.Getenv("OPENFGA_TOKEN")}
	case "local": // local stack only (decision 0004 #11): model JSON plus a tuples file
		ev, err := policy.LoadModel(env("POLICY_MODEL_FILE", "/config/onemind-actions.json"))
		if err != nil {
			return err
		}
		b, err := os.ReadFile(env("POLICY_TUPLES_FILE", "/config/tuples.json"))
		if err != nil {
			return err
		}
		var tuples []policy.Tuple
		if err := json.Unmarshal(b, &tuples); err != nil {
			return err
		}
		if err := ev.Write(tuples...); err != nil {
			return err
		}
		authz = ev
	default:
		return errors.New("POLICY_AUTHZ must be openfga or local")
	}

	var signer *policy.Signer
	if kf := os.Getenv("POLICY_SIGNING_KEY_FILE"); kf != "" {
		if signer, err = policy.LoadSigner(kf, env("POLICY_SIGNING_KID", "policy-1")); err != nil {
			return err
		}
	} else {
		log.Warn("POLICY_SIGNING_KEY_FILE not set: commands are unsigned; bridges must run sig_required:false")
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	zones, err := policy.WatchZones(ctx, js, env("POLICY_GEOFENCE_BUCKET", "geofences"))
	if err != nil {
		return err
	}
	svc, err := policy.New(policy.Config{
		Conn: nc, JS: js, Schemas: schemas, Authz: authz, Directory: dir, Signer: signer, Log: log,
		Entities: &policy.NATSEntities{Conn: nc, Timeout: 2 * time.Second}, Zones: zones,
	})
	if err != nil {
		return err
	}
	if err := svc.Start(); err != nil {
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
	srv := &http.Server{Addr: env("HEALTH_ADDR", ":8080"), Handler: mux, ReadHeaderTimeout: 5 * time.Second}
	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error("health server", "err", err)
		}
	}()
	log.Info("policy service ready", "authz", env("POLICY_AUTHZ", "openfga"), "signed", signer != nil)
	<-ctx.Done()
	sctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return srv.Shutdown(sctx)
}
