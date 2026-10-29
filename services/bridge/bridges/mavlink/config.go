package main

import (
	"crypto/ed25519"
	"encoding/base64"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/bluenviron/gomavlib/v3"

	"github.com/onemindos/omos-bridge/bridges/mavlink/internal/fence"
)

// Config is read from the environment only (see README.md).
type Config struct {
	NATSURL      string
	NATSUser     string
	NATSPassword string
	NATSCreds    string
	Source       string // envelope source of outcome events

	AssetID  string // e.g. sim-x1; must start sim-
	Callsign string
	SysID    uint8 // MAVLink system id of the vehicle

	Endpoint   gomavlib.EndpointConf
	EndpointS  string
	OutVersion gomavlib.Version
	GCSSysID   uint8

	Keys        map[string]ed25519.PublicKey
	SigRequired bool
	Fence       *fence.Fence
	MaxTTL      time.Duration
	ClockSkew   time.Duration

	PositionInterval time.Duration
	StatusInterval   time.Duration
	StaleS           int
	LinkTimeout      time.Duration
	AckTimeout       time.Duration
	AckRetries       int
}

func getenv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}

func durEnv(k string, def time.Duration) (time.Duration, error) {
	v := os.Getenv(k)
	if v == "" {
		return def, nil
	}
	d, err := time.ParseDuration(v)
	if err != nil || d <= 0 {
		return 0, fmt.Errorf("%s: invalid duration %q", k, v)
	}
	return d, nil
}

func intEnv(k string, def, lo, hi int) (int, error) {
	v := os.Getenv(k)
	if v == "" {
		return def, nil
	}
	n, err := strconv.Atoi(v)
	if err != nil || n < lo || n > hi {
		return 0, fmt.Errorf("%s: want integer in [%d,%d], got %q", k, lo, hi, v)
	}
	return n, nil
}

// parseEndpoint accepts udpserver:host:port, udpclient:host:port, tcpclient:host:port,
// tcpserver:host:port and serial:/dev/x:baud.
func parseEndpoint(s string) (gomavlib.EndpointConf, error) {
	kind, rest, ok := strings.Cut(s, ":")
	if !ok || rest == "" {
		return nil, fmt.Errorf("MAVLINK_ENDPOINT %q: want kind:address", s)
	}
	switch kind {
	case "udpserver":
		return gomavlib.EndpointUDPServer{Address: rest}, nil
	case "udpclient":
		return gomavlib.EndpointUDPClient{Address: rest}, nil
	case "tcpclient":
		return gomavlib.EndpointTCPClient{Address: rest}, nil
	case "tcpserver":
		return gomavlib.EndpointTCPServer{Address: rest}, nil
	case "serial":
		i := strings.LastIndex(rest, ":")
		if i < 0 {
			return nil, fmt.Errorf("MAVLINK_ENDPOINT %q: want serial:/dev/x:baud", s)
		}
		baud, err := strconv.Atoi(rest[i+1:])
		if err != nil {
			return nil, fmt.Errorf("MAVLINK_ENDPOINT %q: bad baud", s)
		}
		return gomavlib.EndpointSerial{Device: rest[:i], Baud: baud}, nil
	}
	return nil, fmt.Errorf("MAVLINK_ENDPOINT %q: unknown kind %q", s, kind)
}

// parseKeys reads POLICY_PUBKEYS: comma-separated kid=base64(ed25519 public key).
func parseKeys(s string) (map[string]ed25519.PublicKey, error) {
	keys := map[string]ed25519.PublicKey{}
	for _, part := range strings.Split(s, ",") {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		kid, b64, ok := strings.Cut(part, "=")
		if !ok || kid == "" {
			return nil, fmt.Errorf("POLICY_PUBKEYS: want kid=base64, got %q", part)
		}
		raw, err := base64.StdEncoding.DecodeString(b64)
		if err != nil || len(raw) != ed25519.PublicKeySize {
			return nil, fmt.Errorf("POLICY_PUBKEYS: key %q is not a base64 ed25519 public key", kid)
		}
		keys[kid] = ed25519.PublicKey(raw)
	}
	return keys, nil
}

func loadConfig() (Config, error) {
	c := Config{
		NATSURL:      getenv("NATS_URL", "nats://127.0.0.1:4222"),
		NATSUser:     os.Getenv("NATS_USER"),
		NATSPassword: os.Getenv("NATS_PASSWORD"),
		NATSCreds:    os.Getenv("NATS_CREDS"),
		Source:       getenv("BRIDGE_NAME", "bridge-mavlink"),
		AssetID:      getenv("ASSET_ID", os.Getenv("BRIDGE_ENTITY_ID")),
	}
	if !strings.HasPrefix(c.AssetID, "sim-") {
		return c, fmt.Errorf("ASSET_ID %q must start with sim- (simulated assets only, decision 0004 #13)", c.AssetID)
	}
	c.Callsign = getenv("CALLSIGN", strings.ToUpper(c.AssetID))
	sys, err := intEnv("MAV_SYSID", 1, 1, 255)
	if err != nil {
		return c, err
	}
	c.SysID = uint8(sys)
	gcs, err := intEnv("MAV_GCS_SYSID", 254, 1, 255)
	if err != nil {
		return c, err
	}
	c.GCSSysID = uint8(gcs)

	// Legacy DRONE_IP/DRONE_PORT keep working as a UDP client endpoint.
	ep := os.Getenv("MAVLINK_ENDPOINT")
	if ep == "" && os.Getenv("DRONE_IP") != "" {
		ep = "udpclient:" + os.Getenv("DRONE_IP") + ":" + getenv("DRONE_PORT", "14550")
	}
	if ep == "" {
		ep = "udpserver:127.0.0.1:14550"
	}
	c.EndpointS = ep
	if c.Endpoint, err = parseEndpoint(ep); err != nil {
		return c, err
	}
	switch getenv("MAVLINK_OUT_VERSION", "2") {
	case "2":
		c.OutVersion = gomavlib.V2
	case "1":
		c.OutVersion = gomavlib.V1
	default:
		return c, fmt.Errorf("MAVLINK_OUT_VERSION must be 1 or 2")
	}

	if c.Keys, err = parseKeys(os.Getenv("POLICY_PUBKEYS")); err != nil {
		return c, err
	}
	c.SigRequired = getenv("SIG_REQUIRED", "true") != "false"
	if path := os.Getenv("GEOFENCE_FILE"); path != "" {
		if c.Fence, err = fence.Load(path); err != nil {
			return c, err
		}
	} else if js := os.Getenv("GEOFENCE_JSON"); js != "" {
		if c.Fence, err = fence.Parse([]byte(js)); err != nil {
			return c, err
		}
	}

	for _, d := range []struct {
		k   string
		p   *time.Duration
		def time.Duration
	}{
		{"CMD_MAX_TTL", &c.MaxTTL, 60 * time.Second},
		{"CLOCK_SKEW", &c.ClockSkew, 5 * time.Second},
		{"POSITION_INTERVAL", &c.PositionInterval, 200 * time.Millisecond},
		{"STATUS_INTERVAL", &c.StatusInterval, time.Second},
		{"LINK_TIMEOUT", &c.LinkTimeout, 3 * time.Second},
		{"ACK_TIMEOUT", &c.AckTimeout, 1500 * time.Millisecond},
	} {
		if *d.p, err = durEnv(d.k, d.def); err != nil {
			return c, err
		}
	}
	if c.StaleS, err = intEnv("STALE_S", 10, 1, 3600); err != nil {
		return c, err
	}
	if c.AckRetries, err = intEnv("ACK_RETRIES", 3, 1, 10); err != nil {
		return c, err
	}
	return c, nil
}
