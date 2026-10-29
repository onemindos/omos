package policy

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/nats-io/nats.go"
)

// RequestInfoHeader is set by nats-server on every request that crosses a service import
// (server/accounts.go ClientInfoHdr). With `share: true` on the import it carries the
// requester's account and user; the server replaces any value the client sent
// (server/client.go processServiceImport), so it is the verified identity.
const RequestInfoHeader = "Nats-Request-Info"

// Principal is a resolved caller.
type Principal struct {
	Kind     string // user | agent | system
	ID       string // user:zeus, agent:hermes/zeus/legacy, service:watch
	Verified bool
	Account  string
	User     string
	// Groups are the OneMind role groups of a callout-issued user (FEATURES J2), taken from
	// the server-set user name. checkFGA sends them as contextual tuples.
	Groups []string
}

// DirectoryEntry maps one NATS identity to a OneMind principal.
type DirectoryEntry struct {
	Account string `json:"account"`
	User    string `json:"user"`
	Kind    string `json:"kind"`
	ID      string `json:"id"`
}

// CalloutAccount describes users that the NATS auth callout (services/auth-callout) places
// in an account. The callout names each connection "<name_prefix><username>;<group>,<group>"
// (auth-callout PrincipalName); nats-server copies that name into Nats-Request-Info and a
// client cannot change it. Only groups listed in Groups are kept.
type CalloutAccount struct {
	Account    string   `json:"account"`
	NamePrefix string   `json:"name_prefix"`
	Groups     []string `json:"groups"`
}

// Directory resolves NATS identities. Accounts listed in AgentAccounts always resolve to
// kind agent, whatever the entry says, so a mislabelled agent can never confirm.
type Directory struct {
	OwnAccount    string           `json:"own_account"`
	AgentAccounts []string         `json:"agent_accounts"`
	Entries       []DirectoryEntry `json:"principals"`
	Callout       []CalloutAccount `json:"callout_accounts,omitempty"`
}

var (
	calloutUserRE  = regexp.MustCompile(`^[a-z0-9][a-z0-9._@-]{0,63}$`)
	calloutGroupRE = regexp.MustCompile(`^[a-z0-9][a-z0-9-]{0,62}$`)
)

// resolveCallout parses a callout-issued name. ok is false when the name does not have the
// callout form for this account.
func (d *Directory) resolveCallout(account, user string) (p Principal, ok bool) {
	for _, c := range d.Callout {
		if c.Account != account || c.NamePrefix == "" || !strings.HasPrefix(user, c.NamePrefix) {
			continue
		}
		name, groups, found := strings.Cut(strings.TrimPrefix(user, c.NamePrefix), ";")
		if !found || !calloutUserRE.MatchString(name) {
			return Principal{}, false
		}
		p = Principal{Kind: "user", ID: "user:" + name, Verified: true, Account: account, User: user}
		for _, g := range strings.Split(groups, ",") {
			if g == "" || !calloutGroupRE.MatchString(g) {
				continue
			}
			for _, allowed := range c.Groups {
				if g == allowed {
					p.Groups = append(p.Groups, g)
				}
			}
		}
		return p, true
	}
	return Principal{}, false
}

// LoadDirectory reads the principals file (JSON; mounted from a ConfigMap, no secrets).
func LoadDirectory(path string) (*Directory, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var d Directory
	if err := json.Unmarshal(b, &d); err != nil {
		return nil, err
	}
	if d.OwnAccount == "" {
		return nil, fmt.Errorf("principals file: own_account is required")
	}
	return &d, nil
}

func (d *Directory) isAgentAccount(acc string) bool {
	for _, a := range d.AgentAccounts {
		if a == acc {
			return true
		}
	}
	return false
}

// Resolve returns the verified principal of a request, or an unverified placeholder when
// the request did not cross a service import (same-account caller). known is false when
// a verified identity is not in the directory.
func (d *Directory) Resolve(msg *nats.Msg) (p Principal, known bool) {
	if msg.Header == nil {
		return Principal{}, true
	}
	raw := msg.Header.Get(RequestInfoHeader)
	if raw == "" {
		return Principal{}, true
	}
	var ci struct {
		Account string `json:"acc"`
		User    string `json:"user"`
	}
	if json.Unmarshal([]byte(raw), &ci) != nil || ci.Account == "" || ci.Account == d.OwnAccount {
		// Own-account requests carry no server-set identity; a header there could be forged.
		return Principal{}, true
	}
	// Callout users are always humans; an agent account never hosts callout users.
	if !d.isAgentAccount(ci.Account) {
		if p, ok := d.resolveCallout(ci.Account, ci.User); ok {
			return p, true
		}
	}
	for _, e := range d.Entries {
		if e.Account == ci.Account && e.User == ci.User {
			kind := e.Kind
			if d.isAgentAccount(ci.Account) {
				kind = "agent"
			}
			return Principal{Kind: kind, ID: e.ID, Verified: true, Account: ci.Account, User: ci.User}, true
		}
	}
	kind := ""
	if d.isAgentAccount(ci.Account) {
		kind = "agent"
	}
	return Principal{Kind: kind, Verified: true, Account: ci.Account, User: ci.User}, false
}
