// Command auth-callout is the NATS auth callout for browser and remote clients (FEATURES J2):
// Authentik OIDC token in, short-lived CLIENT-account user with role permissions out.
//
// Configuration is environment only. The issuer and xkey seeds are read from files mounted
// from the SealedSecret auth-callout-runtime (platform/auth-callout/README.md), never from Git.
package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nkeys"

	"github.com/onemindos/omos/services/auth-callout/internal/callout"
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
		log.Error("auth-callout stopped", "err", err)
		os.Exit(1)
	}
}

func seedFile(k string) (nkeys.KeyPair, error) {
	p := os.Getenv(k)
	if p == "" {
		return nil, nil
	}
	b, err := os.ReadFile(p)
	if err != nil {
		return nil, err
	}
	return nkeys.FromSeed([]byte(strings.TrimSpace(string(b))))
}

func run(log *slog.Logger) error {
	issuer, err := seedFile("CALLOUT_ISSUER_SEED_FILE")
	if err != nil {
		return err
	}
	if issuer == nil {
		return errors.New("CALLOUT_ISSUER_SEED_FILE is required")
	}
	xkey, err := seedFile("CALLOUT_XKEY_SEED_FILE")
	if err != nil {
		return err
	}
	roles, err := callout.LoadRoleMap(env("ROLES_FILE", "/config/roles.json"))
	if err != nil {
		return err
	}
	ver, err := callout.NewVerifier(callout.OIDCConfig{
		Issuer: os.Getenv("OIDC_ISSUER"), Audience: os.Getenv("OIDC_AUDIENCE"), JWKSURL: os.Getenv("OIDC_JWKS_URL"),
	})
	if err != nil {
		return err
	}

	opts := []nats.Option{nats.Name("auth-callout"), nats.MaxReconnects(-1)}
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

	svc, err := callout.New(callout.Config{Conn: nc, Issuer: issuer, XKey: xkey, Verifier: ver, Roles: roles, Log: log})
	if err != nil {
		return err
	}
	if err := svc.Start(); err != nil {
		return err
	}
	defer svc.Stop()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
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
	log.Info("auth-callout ready", "client_account", roles.ClientAccount, "roles", len(roles.Roles), "xkey", xkey != nil)
	<-ctx.Done()
	sctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return srv.Shutdown(sctx)
}
