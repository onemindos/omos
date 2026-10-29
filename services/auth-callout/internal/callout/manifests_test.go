package callout

import (
	"encoding/json"
	"os"
	"reflect"
	"strings"
	"testing"
)

// literalBlock returns the YAML literal block under `key: |` (4-space content indent below a
// 2-space key), enough for the two files checked here.
func literalBlock(t *testing.T, path, key string) string {
	t.Helper()
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	lines := strings.Split(string(b), "\n")
	for i, l := range lines {
		if strings.TrimSpace(l) != key+": |" {
			continue
		}
		indent := len(l) - len(strings.TrimLeft(l, " ")) + 2
		var out []string
		for _, c := range lines[i+1:] {
			if strings.TrimSpace(c) != "" && len(c)-len(strings.TrimLeft(c, " ")) < indent {
				break
			}
			if len(c) >= indent {
				out = append(out, c[indent:])
			} else {
				out = append(out, "")
			}
		}
		return strings.Join(out, "\n")
	}
	t.Fatalf("%s: no literal block %q", path, key)
	return ""
}

// The ConfigMap the callout reads in the cluster is the source role map.
func TestConfigMapMatchesSource(t *testing.T) {
	var a, b any
	src, err := os.ReadFile(roleSource)
	if err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal(src, &a); err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal([]byte(literalBlock(t, "../../../../platform/auth-callout/configmap.yaml", "roles.json")), &b); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(a, b) {
		t.Fatal("platform/auth-callout/configmap.yaml roles.json differs from platform/security/roles/onemind-roles.json")
	}
}

// Every role group exists in the Authentik blueprint and can obtain a token.
func TestBlueprintHasEveryRoleGroup(t *testing.T) {
	m, err := LoadRoleMap(roleSource)
	if err != nil {
		t.Fatal(err)
	}
	bp := literalBlock(t, "../../../../platform/auth-callout/authentik-blueprint.yaml", "onemind-xr.yaml")
	for _, r := range m.Roles {
		if !strings.Contains(bp, "identifiers: {name: "+r.Group+"}") {
			t.Errorf("blueprint lacks group %s", r.Group)
		}
		if !strings.Contains(bp, "[name, "+r.Group+"]], order:") {
			t.Errorf("blueprint does not bind group %s to the application", r.Group)
		}
	}
	for _, want := range []string{"client_type: public", "client_id: onemind-xr", "sub_mode: hashed_user_id", "signing_key:"} {
		if !strings.Contains(bp, want) {
			t.Errorf("blueprint lacks %q", want)
		}
	}
	if strings.Contains(bp, "client_secret") {
		t.Error("blueprint carries a client secret")
	}
}
