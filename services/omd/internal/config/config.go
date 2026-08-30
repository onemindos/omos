// Package config loads omd runtime configuration from env vars and CLI flags.
// NO secrets are hardcoded here. The NATS connection currently uses the
// no_auth_user path (onemind account) that the local server already allows.
// When FABRIC-account creds are issued (future), pass them via --creds flag.
package config

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

// C holds the runtime configuration for omd.
type C struct {
	// NATSUrl is the NATS server URL. Default: nats://127.0.0.1:4222.
	// Override with OMD_NATS_URL env var or --nats-url flag.
	NATSUrl string

	// NodeID is the stable bus identity for this node.
	// Default: first DNS label of hostname, lowercased (e.g. "oneminds-macbook-pro").
	// Override with OMD_NODE_ID env var or --node-id flag.
	// MUST NOT contain dots — dots create extra NATS subject tokens in
	// cmd.node.{id}.> and svc.omd.{id}.info endpoints.
	NodeID string

	// CredsFile is the optional path to a NATS .creds file for authenticated accounts.
	// Leave empty to connect via the no_auth_user path (GLOBAL/onemind account).
	// TODO: switch to FABRIC-account creds when P1-A account migration is complete.
	CredsFile string
}

// Load parses flags and environment, returning a populated *C.
// Fatals if NodeID contains a dot (NATS subject token violation).
// Must be called before any flag.Xxx access.
func Load() *C {
	c := &C{}
	flag.StringVar(&c.NATSUrl, "nats-url",
		envOr("OMD_NATS_URL", "nats://127.0.0.1:4222"),
		"NATS server URL (env: OMD_NATS_URL)")
	flag.StringVar(&c.NodeID, "node-id",
		envOr("OMD_NODE_ID", sanitizedHostname()),
		"Stable node identity used as bus actor id (env: OMD_NODE_ID). Must not contain dots.")
	flag.StringVar(&c.CredsFile, "creds", "", "Path to NATS .creds file (optional; default: no-auth user)")
	flag.Parse()

	if strings.Contains(c.NodeID, ".") {
		fmt.Fprintf(os.Stderr,
			"fatal: node-id %q contains a dot, which creates extra NATS subject tokens.\n"+
				"Use a simple hostname label (e.g. \"oneminds-macbook-pro\").\n", c.NodeID)
		os.Exit(1)
	}

	return c
}

func envOr(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

// sanitizedHostname returns the first DNS label of the system hostname,
// lowercased. Example: "OneMinds-MacBook-Pro.local" -> "oneminds-macbook-pro".
// This ensures the node-id is safe as a NATS subject token with no extra dots.
func sanitizedHostname() string {
	h, err := os.Hostname()
	if err != nil {
		return "unknown-node"
	}
	// Take only the first label (before any dot), lowercase it.
	return strings.ToLower(strings.SplitN(h, ".", 2)[0])
}
