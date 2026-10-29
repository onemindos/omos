package callout

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"sort"
	"strings"
)

// RoleMap is platform/security/roles/onemind-roles.json: OneMind roles, the Authentik group
// that grants each one, and the NATS permissions a callout-issued user of that role gets.
type RoleMap struct {
	Version         int      `json:"version"`
	ClientAccount   string   `json:"client_account"`
	NamePrefix      string   `json:"name_prefix"`
	GroupTagPrefix  string   `json:"group_tag_prefix"`
	RoleTagPrefix   string   `json:"role_tag_prefix"`
	InboxPrefix     string   `json:"inbox_prefix"`
	MaxTTLSeconds   int64    `json:"max_ttl_seconds"`
	ConnectionTypes []string `json:"connection_types"`
	Roles           []Role   `json:"roles"`
}

// Role is one mode. Publish and Subscribe are allow lists; everything else is denied.
type Role struct {
	Name      string          `json:"name"`
	Rank      int             `json:"rank"`
	Group     string          `json:"group"`
	UIModes   []string        `json:"ui_modes"`
	Publish   []string        `json:"publish"`
	Subscribe []string        `json:"subscribe"`
	OpenFGA   json.RawMessage `json:"openfga"`
}

// LoadRoleMap reads and validates the role map.
func LoadRoleMap(path string) (*RoleMap, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return ParseRoleMap(b)
}

var (
	groupRE = regexp.MustCompile(`^[a-z0-9][a-z0-9-]{0,62}$`)
	tokenRE = regexp.MustCompile(`^[A-Za-z0-9_-]+$`)
)

// ParseRoleMap parses and validates. It refuses any grant that would let a client publish
// outside the request/reply surface (BUS-DESIGN §7.1: CLIENT has no publish to any data
// subject and no JetStream).
func ParseRoleMap(b []byte) (*RoleMap, error) {
	var m RoleMap
	if err := json.Unmarshal(b, &m); err != nil {
		return nil, err
	}
	if m.Version != 1 {
		return nil, fmt.Errorf("roles: version %d not supported", m.Version)
	}
	if m.ClientAccount == "" || m.NamePrefix == "" || m.GroupTagPrefix == "" || m.RoleTagPrefix == "" {
		return nil, fmt.Errorf("roles: client_account, name_prefix, group_tag_prefix and role_tag_prefix are required")
	}
	if !tokenRE.MatchString(m.InboxPrefix) {
		return nil, fmt.Errorf("roles: inbox_prefix %q must be one subject token", m.InboxPrefix)
	}
	if m.MaxTTLSeconds <= 0 || m.MaxTTLSeconds > 86400 {
		return nil, fmt.Errorf("roles: max_ttl_seconds must be 1..86400")
	}
	for _, ct := range m.ConnectionTypes {
		if ct != "STANDARD" && ct != "WEBSOCKET" {
			return nil, fmt.Errorf("roles: connection type %q not allowed (STANDARD, WEBSOCKET only)", ct)
		}
	}
	if len(m.ConnectionTypes) == 0 {
		return nil, fmt.Errorf("roles: connection_types is required")
	}
	seenRank, seenGroup, seenName := map[int]bool{}, map[string]bool{}, map[string]bool{}
	for _, r := range m.Roles {
		if r.Name == "" || seenName[r.Name] {
			return nil, fmt.Errorf("roles: empty or duplicate role name %q", r.Name)
		}
		if r.Rank <= 0 || seenRank[r.Rank] {
			return nil, fmt.Errorf("roles: role %s: rank must be positive and unique", r.Name)
		}
		if !groupRE.MatchString(r.Group) || seenGroup[r.Group] {
			return nil, fmt.Errorf("roles: role %s: invalid or duplicate group %q", r.Name, r.Group)
		}
		seenName[r.Name], seenRank[r.Rank], seenGroup[r.Group] = true, true, true
		for _, s := range r.Publish {
			if err := clientPublishAllowed(s); err != nil {
				return nil, fmt.Errorf("roles: role %s: publish %q: %v", r.Name, s, err)
			}
		}
		for _, s := range r.Subscribe {
			if err := clientSubscribeAllowed(s); err != nil {
				return nil, fmt.Errorf("roles: role %s: subscribe %q: %v", r.Name, s, err)
			}
		}
	}
	if len(m.Roles) == 0 {
		return nil, fmt.Errorf("roles: no roles")
	}
	return &m, nil
}

// clientPublishAllowed: requests to the imported services, Hermes prompts and reversible
// map actions only.
func clientPublishAllowed(s string) error {
	switch {
	case s == "" || strings.ContainsAny(s, " \t"):
		return fmt.Errorf("empty or blank")
	case s == ">" || strings.HasPrefix(s, "*") || strings.HasPrefix(s, "$") || strings.HasPrefix(s, "_"):
		return fmt.Errorf("wildcard root, system or inbox subject")
	case strings.HasPrefix(s, "svc.") && strings.Count(s, ".") >= 2 && !strings.HasPrefix(s, "svc.*") && !strings.HasPrefix(s, "svc.>"):
		return nil
	case s == "agents.prompt.hermes.>":
		return nil
	case s == "cmd.map.>" || strings.HasPrefix(s, "cmd.map."):
		return nil
	}
	return fmt.Errorf("clients may publish only svc.{service}.*, agents.prompt.hermes.> and cmd.map.>")
}

func clientSubscribeAllowed(s string) error {
	if s == "cmd.map.>" || strings.HasPrefix(s, "cmd.map.") {
		return nil
	}
	return fmt.Errorf("clients may subscribe only cmd.map.> (replies use the per-user inbox)")
}

// Resolve returns the highest-ranked role granted by any of the groups, and the OneMind role
// groups the user holds (sorted). ok is false when no group grants a role.
func (m *RoleMap) Resolve(groups []string) (role Role, held []string, ok bool) {
	have := map[string]bool{}
	for _, g := range groups {
		have[g] = true
	}
	best := -1
	for i, r := range m.Roles {
		if !have[r.Group] {
			continue
		}
		held = append(held, r.Group)
		if best < 0 || r.Rank > m.Roles[best].Rank {
			best = i
		}
	}
	if best < 0 {
		return Role{}, nil, false
	}
	sort.Strings(held)
	return m.Roles[best], held, true
}

// Role returns a role by name.
func (m *RoleMap) Role(name string) (Role, bool) {
	for _, r := range m.Roles {
		if r.Name == name {
			return r, true
		}
	}
	return Role{}, false
}
