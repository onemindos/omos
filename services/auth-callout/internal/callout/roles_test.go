package callout

import (
	"os"
	"slices"
	"strings"
	"testing"
)

// The committed role map satisfies the J2 rules.
func TestRoleMapSourceRules(t *testing.T) {
	m, err := LoadRoleMap(roleSource)
	if err != nil {
		t.Fatal(err)
	}
	if m.ClientAccount != "CLIENT" {
		t.Fatalf("client account %s", m.ClientAccount)
	}
	want := []string{"newcomer", "viewer", "operator", "commander", "admin"}
	for i, name := range want {
		r, ok := m.Role(name)
		if !ok || r.Rank != i+1 || !strings.HasPrefix(r.Group, "onemind-") {
			t.Fatalf("role %s: %+v", name, r)
		}
	}
	nw, _ := m.Role("newcomer")
	for _, s := range append(append([]string{}, nw.Publish...), nw.Subscribe...) {
		if strings.Contains(s, "policy") || strings.HasPrefix(s, "cmd.") || strings.HasPrefix(s, "agents.") {
			t.Fatalf("newcomer may reach %s", s)
		}
	}
	vw, _ := m.Role("viewer")
	for _, s := range vw.Publish {
		if strings.HasSuffix(s, ".propose") || strings.HasSuffix(s, ".confirm") || strings.HasPrefix(s, "cmd.") {
			t.Fatalf("viewer may publish %s", s)
		}
	}
	// Each higher role keeps every right of the role below it (modes only add).
	for i := 1; i < len(m.Roles); i++ {
		lo, hi := m.Roles[i-1], m.Roles[i]
		for _, s := range lo.Publish {
			if !slices.Contains(hi.Publish, s) {
				t.Errorf("%s lacks %s's publish %s", hi.Name, lo.Name, s)
			}
		}
		for _, s := range lo.Subscribe {
			if !slices.Contains(hi.Subscribe, s) {
				t.Errorf("%s lacks %s's subscribe %s", hi.Name, lo.Name, s)
			}
		}
		for _, s := range lo.UIModes {
			if !slices.Contains(hi.UIModes, s) {
				t.Errorf("%s lacks %s's ui mode %s", hi.Name, lo.Name, s)
			}
		}
	}
	// Only operator and above may propose or confirm.
	for _, r := range m.Roles {
		can := slices.Contains(r.Publish, "svc.policy.command.propose")
		if can != (r.Rank >= 3) {
			t.Errorf("%s propose=%v", r.Name, can)
		}
	}
}

func TestParseRoleMapRejectsUnsafeGrants(t *testing.T) {
	base, err := os.ReadFile(roleSource)
	if err != nil {
		t.Fatal(err)
	}
	for _, bad := range []struct{ from, to string }{
		{`"svc.layer.list",`, `"ent.>",`},
		{`"svc.layer.list",`, `"cmd.fleet.drone.sim-x1.goto",`},
		{`"svc.layer.list",`, `"tak.cot.inject",`},
		{`"svc.layer.list",`, `">",`},
		{`"svc.layer.list",`, `"svc.>",`},
		{`"svc.layer.list",`, `"$JS.API.>",`},
		{`"svc.layer.list",`, `"_INBOX.>",`},
		{`"subscribe": ["cmd.map.>"]`, `"subscribe": ["_INBOX.>"]`},
		{`"subscribe": ["cmd.map.>"]`, `"subscribe": ["evt.so.action.>"]`},
		{`"connection_types": ["WEBSOCKET", "STANDARD"]`, `"connection_types": ["LEAFNODE"]`},
		{`"rank": 2,`, `"rank": 1,`},
		{`"group": "onemind-viewers"`, `"group": "onemind-newcomers"`},
		{`"max_ttl_seconds": 3600`, `"max_ttl_seconds": 0`},
	} {
		s := strings.Replace(string(base), bad.from, bad.to, 1)
		if s == string(base) {
			t.Fatalf("pattern %q not found", bad.from)
		}
		if _, err := ParseRoleMap([]byte(s)); err == nil {
			t.Errorf("accepted %s", bad.to)
		}
	}
}

func TestResolveHighestRole(t *testing.T) {
	m, err := LoadRoleMap(roleSource)
	if err != nil {
		t.Fatal(err)
	}
	r, held, ok := m.Resolve([]string{"onemind-viewers", "onemind-commanders", "other"})
	if !ok || r.Name != "commander" || strings.Join(held, ",") != "onemind-commanders,onemind-viewers" {
		t.Fatalf("%v %v %v", r.Name, held, ok)
	}
	if _, _, ok := m.Resolve([]string{"argocd-admins"}); ok {
		t.Fatal("non-OneMind group granted a role")
	}
	if _, _, ok := m.Resolve(nil); ok {
		t.Fatal("no groups granted a role")
	}
}
