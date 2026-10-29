package registry

import (
	"fmt"
	"os"
	"regexp"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"
)

// MediaConfig is the part of mediamtx.yml the registry reads: the users (one per bus identity)
// and the per-class recording retention. The registry mounts the same ConfigMap MediaMTX uses,
// so the permissions it publishes cannot drift from what MediaMTX enforces.
type MediaConfig struct {
	AuthMethod string               `yaml:"authMethod"`
	Users      []AuthUser           `yaml:"authInternalUsers"`
	PathDefs   PathConf             `yaml:"pathDefaults"`
	Paths      map[string]*PathConf `yaml:"paths"`
}

// AuthUser is one entry of authInternalUsers.
type AuthUser struct {
	User        string       `yaml:"user"`
	Pass        string       `yaml:"pass"`
	IPs         []string     `yaml:"ips"`
	Permissions []Permission `yaml:"permissions"`
}

// Permission is one action grant, optionally limited to a path (exact, or ~regex).
type Permission struct {
	Action string `yaml:"action"`
	Path   string `yaml:"path"`
}

// PathConf holds the recording fields of a path entry.
type PathConf struct {
	Record            *bool  `yaml:"record"`
	RecordFormat      string `yaml:"recordFormat"`
	RecordDeleteAfter string `yaml:"recordDeleteAfter"`
}

// Permissions lists the bus identities allowed each action on one stream.
type Permissions struct {
	Publish  []string `json:"publish"`
	Read     []string `json:"read"`
	Playback []string `json:"playback"`
}

// Recording describes how MediaMTX records the stream.
type Recording struct {
	Enabled   bool   `json:"enabled"`
	Format    string `json:"format,omitempty"`
	Retention string `json:"retention,omitempty"`
}

// LoadMediaConfig parses mediamtx.yml.
func LoadMediaConfig(path string) (*MediaConfig, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var c MediaConfig
	if err := yaml.Unmarshal(b, &c); err != nil {
		return nil, fmt.Errorf("%s: %w", path, err)
	}
	if c.AuthMethod != "internal" {
		return nil, fmt.Errorf("%s: authMethod %q, want internal", path, c.AuthMethod)
	}
	return &c, nil
}

func matchPath(pattern, name string) bool {
	if pattern == "" {
		return true
	}
	if strings.HasPrefix(pattern, "~") {
		re, err := regexp.Compile(pattern[1:])
		return err == nil && re.MatchString(name)
	}
	return pattern == name
}

// PermissionsFor returns who may publish, read and play back a stream. Users without a password
// (only the localhost metrics user) are ignored: they are not bus identities.
func (c *MediaConfig) PermissionsFor(stream string) Permissions {
	set := map[string]map[string]bool{"publish": {}, "read": {}, "playback": {}}
	for _, u := range c.Users {
		if u.Pass == "" || u.User == "any" {
			continue
		}
		for _, p := range u.Permissions {
			if m, ok := set[p.Action]; ok && matchPath(p.Path, stream) {
				m[u.User] = true
			}
		}
	}
	list := func(m map[string]bool) []string {
		out := make([]string, 0, len(m))
		for k := range m {
			out = append(out, k)
		}
		sort.Strings(out)
		return out
	}
	return Permissions{Publish: list(set["publish"]), Read: list(set["read"]), Playback: list(set["playback"])}
}

// RecordingFor resolves the recording settings that apply to a stream (path entry over
// pathDefaults, first matching regex entry in sorted order, as MediaMTX matches one entry).
func (c *MediaConfig) RecordingFor(stream string) Recording {
	r := Recording{Format: c.PathDefs.RecordFormat, Retention: c.PathDefs.RecordDeleteAfter}
	if c.PathDefs.Record != nil {
		r.Enabled = *c.PathDefs.Record
	}
	names := make([]string, 0, len(c.Paths))
	for n := range c.Paths {
		names = append(names, n)
	}
	sort.Strings(names)
	for _, n := range names {
		if !matchPath(n, stream) {
			continue
		}
		if p := c.Paths[n]; p != nil {
			if p.Record != nil {
				r.Enabled = *p.Record
			}
			if p.RecordFormat != "" {
				r.Format = p.RecordFormat
			}
			if p.RecordDeleteAfter != "" {
				r.Retention = p.RecordDeleteAfter
			}
		}
		break
	}
	if !r.Enabled {
		return Recording{Enabled: false}
	}
	return r
}
