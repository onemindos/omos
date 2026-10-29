package projector

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

// Target is where the projector sends CoT (cot-projection.md §2).
type Target string

const (
	TargetServer Target = "server" // hub: tak.cot.inject (NatsBridgePlugin)
	TargetMesh   Target = "mesh"   // site: tak.cot.mesh.tx (mesh gateway)
)

// Config is the projector configuration (cot-projection.md §10). Durations are Go durations
// so tests can shorten them; the environment gives seconds as the contract names them.
type Config struct {
	Target           Target
	Instance         string // "hub" or the site id; OneMind-Projector header
	Subject          string // output subject
	OwnInterval      time.Duration
	OtherInterval    time.Duration
	MaxRatePerSec    int
	MaxProjected     int
	ReleaseGroups    []string
	AOIExitGrace     time.Duration
	MeshWorldFeeds   bool
	MeshMaxBytes     int
	MaxVertices      int
	HeartbeatMax     time.Duration // upper bound of the heartbeat period (300 s)
	EntStream        string        // ENT; "" or ENT_MODE=core uses a core subscription
	EntSubject       string
	GeofenceBucket   string
	MissionBucket    string
	TakInputSubjects []string // watched by the tasking guard only (SL1)
	ControlDisable   string
	ControlEnable    string
	Tick             time.Duration // scheduler tick
	CapEvery         time.Duration // how often the cap ranking is recomputed
	Send             SendConfig    // svc.cot.send (cotsend.go)
}

// SendConfig configures svc.cot.send, the request/reply endpoint for client-originated CoT
// (cot-projection.md §1 "Ephemeral client sends", §9; owner proposal WEB-P4 and B-P4).
type SendConfig struct {
	Subject         string        // "" disables the endpoint
	Queue           string        // queue group, so replicas answer once
	Accounts        []string      // accounts whose verified users may send (share: true import)
	OwnAccount      string        // the projector's account; its callers carry no verified identity
	AllowUnverified bool          // local tests only: accept same-account callers
	PerIdentityMin  int           // sustained sends per identity per minute
	Burst           int           // per-identity burst
	MaxPerSec       int           // all client sends on this target, per second
	MaxBytes        int           // largest accepted CoT XML
	MaxStale        time.Duration // furthest allowed stale
	Skew            time.Duration // allowed |time - now|
	ChatStale       time.Duration // stale of a GeoChat built from the structured body
	ChatDefaultTLP  string        // marking applied to a GeoChat body without one ("" = refuse)
	DedupTTL        time.Duration // how long a client_ref / message_uid reply is remembered
}

// Defaults returns the contract defaults for a target.
func Defaults(t Target) Config {
	c := Config{
		Target:           t,
		Instance:         "hub",
		Subject:          "tak.cot.inject",
		OwnInterval:      time.Second,
		OtherInterval:    5 * time.Second,
		MaxRatePerSec:    50,
		MaxProjected:     500,
		AOIExitGrace:     60 * time.Second,
		MeshMaxBytes:     1200,
		MaxVertices:      200,
		HeartbeatMax:     300 * time.Second,
		EntStream:        "ENT",
		EntSubject:       "ent.*.*.state",
		GeofenceBucket:   "geofences",
		MissionBucket:    "mission-state",
		TakInputSubjects: []string{"tak.cot.out", "tak.cot.mesh.rx"},
		ControlDisable:   "cmd.projector.disable", // live legacy subjects; CP-1 proposes a 4-token form
		ControlEnable:    "cmd.projector.enable",
		Tick:             20 * time.Millisecond,
		CapEvery:         time.Second,
		Send: SendConfig{
			Subject:        "svc.cot.send",
			Queue:          "omos-tak-projector",
			Accounts:       []string{"CLIENT"},
			OwnAccount:     "ONEMIND",
			PerIdentityMin: 30,
			Burst:          10,
			MaxPerSec:      10,
			MaxBytes:       64 * 1024,
			MaxStale:       24 * time.Hour,
			Skew:           10 * time.Minute,
			ChatStale:      10 * time.Minute,
			ChatDefaultTLP: "CLEAR",
			DedupTTL:       10 * time.Minute,
		},
	}
	if t == TargetMesh {
		c.Instance = "site"
		c.Subject = "tak.cot.mesh.tx"
		c.OwnInterval = 5 * time.Second
		c.OtherInterval = 15 * time.Second
		c.MaxRatePerSec = 10
		c.MaxProjected = 100
		c.Send.MaxPerSec = 2
		c.Send.MaxBytes = c.MeshMaxBytes
	}
	return c
}

func envList(k string, def []string) []string {
	v, ok := os.LookupEnv(k)
	if !ok {
		return def
	}
	var out []string
	for _, p := range strings.Split(v, ",") {
		if p = strings.TrimSpace(p); p != "" {
			out = append(out, p)
		}
	}
	return out
}

// sendFromEnv reads the COT_SEND_* variables (README "svc.cot.send").
func sendFromEnv(c *Config) error {
	s := &c.Send
	var err error
	if v, ok := os.LookupEnv("COT_SEND_SUBJECT"); ok {
		s.Subject = strings.TrimSpace(v)
	}
	if os.Getenv("COT_SEND_ENABLED") == "false" {
		s.Subject = ""
	}
	s.Accounts = envList("COT_SEND_ACCOUNTS", s.Accounts)
	if v := os.Getenv("COT_SEND_OWN_ACCOUNT"); v != "" {
		s.OwnAccount = v
	}
	s.AllowUnverified = os.Getenv("COT_SEND_ALLOW_UNVERIFIED") == "true"
	if s.PerIdentityMin, err = envInt("COT_SEND_RATE_PER_MIN", s.PerIdentityMin); err != nil {
		return err
	}
	if s.Burst, err = envInt("COT_SEND_BURST", s.Burst); err != nil {
		return err
	}
	if s.MaxPerSec, err = envInt("COT_SEND_MAX_PER_SEC", s.MaxPerSec); err != nil {
		return err
	}
	if s.MaxBytes, err = envInt("COT_SEND_MAX_BYTES", s.MaxBytes); err != nil {
		return err
	}
	if s.MaxStale, err = envSeconds("COT_SEND_MAX_STALE_SEC", s.MaxStale); err != nil {
		return err
	}
	if s.ChatStale, err = envSeconds("COT_SEND_CHAT_STALE_SEC", s.ChatStale); err != nil {
		return err
	}
	if v, ok := os.LookupEnv("COT_SEND_CHAT_DEFAULT_TLP"); ok {
		s.ChatDefaultTLP = strings.ToUpper(strings.TrimSpace(v))
	}
	if s.Subject != "" && (s.PerIdentityMin == 0 || s.Burst == 0 || s.MaxPerSec == 0 || s.MaxBytes == 0) {
		return fmt.Errorf("COT_SEND_RATE_PER_MIN, COT_SEND_BURST, COT_SEND_MAX_PER_SEC and COT_SEND_MAX_BYTES must be at least 1")
	}
	return nil
}

func envSeconds(k string, def time.Duration) (time.Duration, error) {
	v := os.Getenv(k)
	if v == "" {
		return def, nil
	}
	f, err := strconv.ParseFloat(v, 64)
	if err != nil || f < 0 {
		return 0, fmt.Errorf("%s: want seconds, got %q", k, v)
	}
	return time.Duration(f * float64(time.Second)), nil
}

func envInt(k string, def int) (int, error) {
	v := os.Getenv(k)
	if v == "" {
		return def, nil
	}
	n, err := strconv.Atoi(v)
	if err != nil || n < 0 {
		return 0, fmt.Errorf("%s: want a non-negative integer, got %q", k, v)
	}
	return n, nil
}

// FromEnv reads the §10 variables. TAK_GROUP is accepted and ignored (deprecated, §5.4).
func FromEnv() (Config, error) {
	t := Target(strings.ToLower(os.Getenv("PROJECTOR_TARGET")))
	if t == "" {
		t = TargetServer
	}
	if t != TargetServer && t != TargetMesh {
		return Config{}, fmt.Errorf("PROJECTOR_TARGET must be server or mesh, got %q", t)
	}
	c := Defaults(t)
	var err error
	if v := os.Getenv("PROJECTOR_INSTANCE"); v != "" {
		c.Instance = v
	}
	if t == TargetServer {
		if v := os.Getenv("TAK_INJECT_SUBJECT"); v != "" {
			c.Subject = v
		}
	} else if v := os.Getenv("MESH_TX_SUBJECT"); v != "" {
		c.Subject = v
	}
	if c.OtherInterval, err = envSeconds("MIN_INTERVAL_SEC", c.OtherInterval); err != nil {
		return c, err
	}
	if c.OwnInterval, err = envSeconds("OWN_MIN_INTERVAL_SEC", c.OwnInterval); err != nil {
		return c, err
	}
	if c.AOIExitGrace, err = envSeconds("AOI_EXIT_GRACE_S", c.AOIExitGrace); err != nil {
		return c, err
	}
	if c.MaxRatePerSec, err = envInt("MAX_RATE_PER_SEC", c.MaxRatePerSec); err != nil {
		return c, err
	}
	if c.MaxProjected, err = envInt("MAX_PROJECTED", c.MaxProjected); err != nil {
		return c, err
	}
	if c.MeshMaxBytes, err = envInt("MESH_MAX_BYTES", c.MeshMaxBytes); err != nil {
		return c, err
	}
	if c.MaxVertices, err = envInt("MAX_VERTICES", c.MaxVertices); err != nil {
		return c, err
	}
	c.MeshWorldFeeds = os.Getenv("MESH_WORLD_FEEDS") == "true"
	for _, g := range strings.Split(os.Getenv("TAK_RELEASE_GROUPS"), ",") {
		if g = strings.TrimSpace(g); g != "" {
			c.ReleaseGroups = append(c.ReleaseGroups, g)
		}
	}
	if v, ok := os.LookupEnv("ENT_STREAM"); ok {
		c.EntStream = v
	}
	if os.Getenv("ENT_MODE") == "core" {
		c.EntStream = ""
	}
	if c.MaxRatePerSec == 0 {
		return c, fmt.Errorf("MAX_RATE_PER_SEC must be at least 1")
	}
	if t == TargetMesh {
		c.Send.MaxBytes = c.MeshMaxBytes
	}
	if err = sendFromEnv(&c); err != nil {
		return c, err
	}
	return c, nil
}
