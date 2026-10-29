package ingest

import (
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

// Kinds of delivery the framework handles. Webhook push and MQTT are covered by the
// existing bridges (services/bridge/bridges/{webhook,mqtt}) until they adopt this core.
const (
	KindRESTPoll  = "rest_poll"
	KindSSE       = "sse"
	KindWebSocket = "websocket"
	KindRSS       = "rss"
	KindFile      = "file"
)

// Duration is a time.Duration that reads "30s" / "5m" from YAML.
type Duration struct{ time.Duration }

func (d *Duration) UnmarshalYAML(n *yaml.Node) error {
	v, err := time.ParseDuration(n.Value)
	if err != nil {
		return fmt.Errorf("line %d: %w", n.Line, err)
	}
	d.Duration = v
	return nil
}

// Config is the feed registry for one omos-ingest process (sources.yaml).
type Config struct {
	// LockBucket holds one lease per source (KV ingest-locks, TTL 10 s).
	LockBucket string `yaml:"lockBucket"`
	// StateBucket holds checkpoints, written only after JetStream acks.
	StateBucket string `yaml:"stateBucket"`
	// LeaseTTL must match the lock bucket TTL; renewals run every LeaseTTL/3.
	LeaseTTL Duration         `yaml:"leaseTTL"`
	Status   Duration         `yaml:"statusInterval"`
	Sources  []SourceConfig   `yaml:"sources"`
	Defaults SourceConfigBase `yaml:"defaults"`
}

// SourceConfigBase holds the knobs every source has, so defaults can set them once.
type SourceConfigBase struct {
	Interval     Duration          `yaml:"interval"`
	StaleAfter   Duration          `yaml:"staleAfter"`
	MaxBackoff   Duration          `yaml:"maxBackoff"`
	MaxBodyBytes int64             `yaml:"maxBodyBytes"`
	SeenMax      int               `yaml:"seenMax"`
	QueueMax     int               `yaml:"queueMax"`
	Headers      map[string]string `yaml:"headers"`
	RateLimit    RateLimitConfig   `yaml:"rateLimit"`
}

type RateLimitConfig struct {
	// Every is the minimum spacing between upstream requests; Burst allows short bursts.
	Every Duration `yaml:"every"`
	Burst int      `yaml:"burst"`
}

type QuotaConfig struct {
	// Requests is the most upstream requests allowed per Per window (0 = no quota).
	Requests int      `yaml:"requests"`
	Per      Duration `yaml:"per"`
}

// Terms records the upstream's terms of use; the loader refuses a source without them.
type Terms struct {
	Licence string `yaml:"licence"`
	URL     string `yaml:"url"`
	Notes   string `yaml:"notes"`
}

type SourceConfig struct {
	SourceConfigBase `yaml:",inline"`

	// ID is the feed id: one subject token, used in tel.feed.{id}.status,
	// evt.sys.{id}.rejected, the lease key and the checkpoint key.
	ID      string `yaml:"id"`
	Kind    string `yaml:"kind"`
	Parser  string `yaml:"parser"`
	URL     string `yaml:"url"`
	Subject string `yaml:"subject"`
	// Schema is a JSON Schema file relative to the schemas directory (services/schemas).
	Schema    string            `yaml:"schema"`
	Quota     QuotaConfig       `yaml:"quota"`
	Terms     Terms             `yaml:"terms"`
	Options   map[string]string `yaml:"options"`
	Subscribe string            `yaml:"subscribe"` // websocket: message sent after connect (env-expanded)
	Enabled   *bool             `yaml:"enabled"`
	// Registry is the onemind-xr docs/sources/sources.yaml id, for traceability.
	Registry string `yaml:"registry"`
}

func (s SourceConfig) IsEnabled() bool { return s.Enabled == nil || *s.Enabled }

var tokenRe = regexp.MustCompile(`^[a-z0-9][a-z0-9_-]*$`)

// LoadConfig reads and validates a registry file. ${VAR} in header values and the
// websocket subscribe message are expanded from the environment, so keys never sit in Git.
func LoadConfig(path string) (*Config, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return ParseConfig(b)
}

func ParseConfig(b []byte) (*Config, error) {
	var c Config
	if err := yaml.Unmarshal(b, &c); err != nil {
		return nil, err
	}
	c.applyDefaults()
	if err := c.Validate(); err != nil {
		return nil, err
	}
	return &c, nil
}

func (c *Config) applyDefaults() {
	if c.LockBucket == "" {
		c.LockBucket = "ingest-locks"
	}
	if c.StateBucket == "" {
		c.StateBucket = "ingest-state"
	}
	if c.LeaseTTL.Duration == 0 {
		c.LeaseTTL.Duration = 10 * time.Second
	}
	if c.Status.Duration == 0 {
		c.Status.Duration = 10 * time.Second
	}
	d := c.Defaults
	for i := range c.Sources {
		s := &c.Sources[i]
		if s.Interval.Duration == 0 {
			s.Interval = or(d.Interval, time.Minute)
		}
		if s.StaleAfter.Duration == 0 {
			s.StaleAfter = or(d.StaleAfter, 5*s.Interval.Duration)
		}
		if s.MaxBackoff.Duration == 0 {
			s.MaxBackoff = or(d.MaxBackoff, 10*time.Minute)
		}
		if s.MaxBodyBytes == 0 {
			s.MaxBodyBytes = d.MaxBodyBytes
			if s.MaxBodyBytes == 0 {
				s.MaxBodyBytes = 16 << 20
				if s.Kind == KindFile {
					s.MaxBodyBytes = 64 << 20
				}
			}
		}
		if s.SeenMax == 0 {
			s.SeenMax = d.SeenMax
			if s.SeenMax == 0 {
				s.SeenMax = 4096
			}
		}
		if s.QueueMax == 0 {
			s.QueueMax = d.QueueMax
			if s.QueueMax == 0 {
				s.QueueMax = 10000
			}
		}
		if s.RateLimit.Every.Duration == 0 {
			s.RateLimit.Every = or(d.RateLimit.Every, time.Second)
		}
		if s.RateLimit.Burst == 0 {
			s.RateLimit.Burst = max(d.RateLimit.Burst, 1)
		}
		h := map[string]string{}
		for k, v := range d.Headers {
			h[k] = v
		}
		for k, v := range s.Headers {
			h[k] = v
		}
		for k, v := range h {
			h[k] = os.ExpandEnv(v)
		}
		s.Headers = h
		s.Subscribe = os.ExpandEnv(s.Subscribe)
	}
}

func or(d Duration, def time.Duration) Duration {
	if d.Duration != 0 {
		return d
	}
	return Duration{def}
}

func (c *Config) Validate() error {
	seen := map[string]bool{}
	for _, s := range c.Sources {
		if !tokenRe.MatchString(s.ID) {
			return fmt.Errorf("source %q: id must be one lower-case subject token", s.ID)
		}
		if seen[s.ID] {
			return fmt.Errorf("source %q: duplicate id", s.ID)
		}
		seen[s.ID] = true
		switch s.Kind {
		case KindRESTPoll, KindSSE, KindWebSocket, KindRSS, KindFile:
		default:
			return fmt.Errorf("source %q: unknown kind %q", s.ID, s.Kind)
		}
		if s.URL == "" {
			return fmt.Errorf("source %q: url is required", s.ID)
		}
		if s.Schema == "" {
			return fmt.Errorf("source %q: schema is required (schema validation at the edge of the bus)", s.ID)
		}
		if s.Terms.Licence == "" || s.Terms.URL == "" {
			return fmt.Errorf("source %q: terms.licence and terms.url are required", s.ID)
		}
		if s.Subject != "" {
			if err := CheckDataSubject(s.Subject); err != nil {
				return fmt.Errorf("source %q: %w", s.ID, err)
			}
		}
		if s.Quota.Requests > 0 && s.Quota.Per.Duration <= 0 {
			return fmt.Errorf("source %q: quota.per is required with quota.requests", s.ID)
		}
		if s.Kind != KindWebSocket && s.Kind != KindSSE && s.Interval.Duration < time.Second {
			return fmt.Errorf("source %q: interval below 1s is not polite", s.ID)
		}
	}
	return nil
}

// CheckDataSubject enforces the locked omos grammar for what a collector may publish
// (omos CLAUDE.md; onemind-xr docs/design/BUS-DESIGN.md §3.1): intel.{domain}.{source}.{type},
// evt.world.{class}.{id}, tel.{...}, det.{domain}.{type}.{id}. Collectors never publish
// ent.*, cmd.*, tak.cot.* or svc.*.
func CheckDataSubject(subj string) error {
	toks := strings.Split(subj, ".")
	for _, t := range toks {
		if t == "" || strings.ContainsAny(t, " *>\t\r\n") {
			return fmt.Errorf("subject %q: empty or wildcard token", subj)
		}
	}
	switch toks[0] {
	case "intel":
		if len(toks) != 4 {
			return fmt.Errorf("subject %q: new collectors publish intel.{domain}.{source}.{type} (4 tokens; BUS-DESIGN §3.1)", subj)
		}
	case "evt":
		if len(toks) != 4 || toks[1] != "world" {
			return fmt.Errorf("subject %q: collectors publish only evt.world.{class}.{id}", subj)
		}
	case "det":
		if len(toks) != 4 {
			return fmt.Errorf("subject %q: det.{domain}.{type}.{id} has 4 tokens", subj)
		}
	case "tel":
		if len(toks) < 4 {
			return fmt.Errorf("subject %q: tel subjects have at least 4 tokens", subj)
		}
	default:
		return fmt.Errorf("subject %q: collectors may publish only intel.*, evt.world.*, det.*, tel.*", subj)
	}
	return nil
}

// StatusSubject and RejectSubject are the fixed health and dead-letter subjects
// (BUS-DESIGN §3.1 tel.feed.{source}.status; §3.3 evt.sys.{source}.rejected, BUS-P9).
func StatusSubject(id string) string { return "tel.feed." + id + ".status" }
func RejectSubject(id string) string { return "evt.sys." + id + ".rejected" }
