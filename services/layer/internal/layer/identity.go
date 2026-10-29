package layer

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/nats-io/nats.go"
)

// RequestInfoHeader is set by nats-server on a request that crosses a service import with
// `share: true`, and the server replaces any value the client sent (the same mechanism the
// policy service relies on, services/policy/internal/policy/identity.go in omos PR #554).
const RequestInfoHeader = "Nats-Request-Info"

// importReplyPrefix starts the reply subject nats-server substitutes on a request that crossed
// a service import (the response mapping). A same-account requester replies on its own inbox.
const importReplyPrefix = "_R_."

// Principal entry: one verified NATS identity and its OneMind groups.
type Principal struct {
	Account string   `json:"account"`
	User    string   `json:"user"`
	ID      string   `json:"id"`
	Groups  []string `json:"groups"`
}

// Directory maps verified identities to groups. OwnAccount is the account the layer service
// runs in: a request from it (no server-set identity) is a trusted service and sees every
// entity, like a direct ENT reader would. OrgGroup is the group whose members may see AMBER
// and AMBER+STRICT entities that carry no explicit releasable_to list.
type Directory struct {
	OwnAccount string      `json:"own_account"`
	OrgGroup   string      `json:"org_group"`
	Principals []Principal `json:"principals"`
}

// LoadDirectory reads the principals file (JSON from a ConfigMap; no secrets). An empty path
// gives a directory with no cross-account principals: every request that crosses a service
// import is then refused with `forbidden`.
func LoadDirectory(path string) (*Directory, error) {
	d := &Directory{OwnAccount: "ONEMIND", OrgGroup: "onemind"}
	if path == "" {
		return d, nil
	}
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(b, d); err != nil {
		return nil, fmt.Errorf("principals file: %w", err)
	}
	if d.OwnAccount == "" {
		return nil, fmt.Errorf("principals file: own_account is required")
	}
	return d, nil
}

// Caller is a resolved requester.
type Caller struct {
	Trusted bool // same-account service
	Known   bool // a verified identity found in the directory (or trusted)
	ID      string
	Groups  map[string]bool
	org     string
}

// Resolve returns the caller of msg. A request with no server-set identity is a trusted
// same-account service only when its reply subject was not rewritten by a service import.
func (d *Directory) Resolve(msg *nats.Msg) Caller {
	raw := ""
	if msg.Header != nil {
		raw = msg.Header.Get(RequestInfoHeader)
	}
	if raw == "" {
		if strings.HasPrefix(msg.Reply, importReplyPrefix) {
			// Crossed a service import that does not share the requester's identity
			// (`share: false`): fail closed instead of trusting it as a same-account service.
			return Caller{ID: "import-without-share"}
		}
		return Caller{Trusted: true, Known: true, ID: "service"}
	}
	var ci struct {
		Account string `json:"acc"`
		User    string `json:"user"`
	}
	if json.Unmarshal([]byte(raw), &ci) != nil || ci.Account == "" || ci.Account == d.OwnAccount {
		// Own-account requests carry no server-set identity; a header there could be forged,
		// and own-account callers may read ENT directly anyway.
		return Caller{Trusted: true, Known: true, ID: "service"}
	}
	for _, p := range d.Principals {
		if p.Account == ci.Account && p.User == ci.User {
			g := map[string]bool{}
			for _, x := range p.Groups {
				g[x] = true
			}
			return Caller{Known: true, ID: p.ID, Groups: g, org: d.OrgGroup}
		}
	}
	return Caller{ID: ci.Account + "/" + ci.User}
}

// CanSee applies the release rule (security-and-trust.md §5, entity.md `marking`):
//   - trusted same-account callers see everything;
//   - CLEAR and GREEN are visible to every known caller;
//   - AMBER, AMBER+STRICT and RED need a caller group in releasable_to; AMBER and
//     AMBER+STRICT with an empty releasable_to are visible to the organisation group only;
//     RED with an empty list is visible to no cross-account caller.
func (c Caller) CanSee(tlp string, releasable []string) bool {
	if c.Trusted {
		return true
	}
	if !c.Known {
		return false
	}
	switch tlp {
	case "CLEAR", "GREEN":
		return true
	}
	for _, g := range releasable {
		if c.Groups[g] {
			return true
		}
	}
	if len(releasable) == 0 && (tlp == "AMBER" || tlp == "AMBER+STRICT") && c.org != "" {
		return c.Groups[c.org]
	}
	return false
}
