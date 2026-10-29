package callout

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/nats-io/nats-server/v2/server"
	"github.com/nats-io/nats.go"
)

// FEATURES J2 acceptance test: a newcomer cannot see asset commands. Layer 1 (this file):
// the bus refuses every svc.policy.* request and every subscription that would show one.
func TestNewcomerCannotSeeOrProposeAssetCommands(t *testing.T) {
	w := newWorld(t, worldOpts{})
	pol := w.password("policy")
	policy := w.serve(pol, "svc.policy.>")
	layer := w.serve(pol, "svc.layer.>")

	tok := w.idp.token(t, "nina", "onemind-newcomers")
	c, err := w.connect(tok, subOf("nina"), true) // over WebSocket, as the L0 web client
	if err != nil {
		t.Fatalf("newcomer connect: %v", err)
	}
	// Layers and entities are open to newcomers.
	if _, err := c.nc.Request("svc.layer.list", []byte(`{"v":1}`), 2*time.Second); err != nil {
		t.Fatalf("newcomer svc.layer.list: %v", err)
	}
	// Every policy endpoint is refused by the server before it reaches the service.
	for _, s := range []string{"svc.policy.command.propose", "svc.policy.command.confirm", "svc.policy.task.propose",
		"svc.policy.action.get", "svc.policy.check"} {
		if _, err := c.nc.Request(s, []byte(`{"v":1}`), 300*time.Millisecond); err == nil {
			t.Fatalf("newcomer reached %s", s)
		}
	}
	// No way to watch commands: map actions, the action audit, or other users' replies.
	for _, s := range []string{"cmd.map.>", "evt.so.action.>", "_INBOX.>", "cmd.>", ">"} {
		_, _ = c.nc.SubscribeSync(s)
	}
	// Nor to publish one.
	for _, s := range []string{"cmd.fleet.drone.sim-x1.goto", "cmd.map.highlight.x", "tak.cot.inject", "ent.drone.sim-x1.state", "$JS.API.INFO"} {
		_ = c.nc.Publish(s, []byte(`{}`))
	}
	_ = c.nc.Flush()
	waitFor(t, func() bool { return len(c.permErrors()) >= 15 }, "15 permission violations")
	for _, e := range c.permErrors() {
		if !errors.Is(e, nats.ErrPermissionViolation) && !strings.Contains(e.Error(), "Permissions Violation") {
			t.Fatalf("unexpected error %v", e)
		}
	}
	if infos, _ := policy.list(); len(infos) != 0 {
		t.Fatalf("policy service saw %d newcomer requests", len(infos))
	}
	if infos, _ := layer.list(); len(infos) != 1 {
		t.Fatalf("layer service saw %d requests, want 1", len(infos))
	}
	t.Logf("newcomer: %d permission violations, 0 policy requests delivered", len(c.permErrors()))
}

func TestViewerSeesActionsButCannotProposeOrConfirm(t *testing.T) {
	w := newWorld(t, worldOpts{})
	pol := w.password("policy")
	policy := w.serve(pol, "svc.policy.>")
	c, err := w.connect(w.idp.token(t, "vera", "onemind-viewers"), subOf("vera"), false)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := c.nc.Request("svc.policy.action.get", []byte(`{"v":1,"action_id":"x"}`), 2*time.Second); err != nil {
		t.Fatalf("viewer action.get: %v", err)
	}
	for _, s := range []string{"svc.policy.command.propose", "svc.policy.command.confirm", "svc.policy.task.propose", "svc.policy.task.confirm", "svc.policy.check"} {
		if _, err := c.nc.Request(s, []byte(`{"v":1}`), 300*time.Millisecond); err == nil {
			t.Fatalf("viewer reached %s", s)
		}
	}
	if _, subj := policy.list(); len(subj) != 1 || subj[0] != "svc.policy.action.get" {
		t.Fatalf("policy saw %v", subj)
	}
	// Viewers receive map actions but cannot send them.
	if _, err := c.nc.SubscribeSync("cmd.map.>"); err != nil {
		t.Fatal(err)
	}
	_ = c.nc.Publish("cmd.map.highlight.x", nil)
	_ = c.nc.Flush()
	waitFor(t, func() bool {
		for _, e := range c.permErrors() {
			if strings.Contains(e.Error(), `Publish to "cmd.map.highlight.x"`) {
				return true
			}
		}
		return false
	}, "viewer cmd.map publish refused")
}

// The policy service resolves the caller from the server-set Nats-Request-Info header:
// name "ak:<username>", tags role:<role> and group:<group> (services/policy identity.go).
func TestOperatorReachesPolicyWithVerifiedIdentity(t *testing.T) {
	w := newWorld(t, worldOpts{})
	pol := w.password("policy")
	policy := w.serve(pol, "svc.policy.>")
	tok := w.idp.token(t, "olga", "onemind-operators", "onemind-viewers", "argocd-admins")
	c, err := w.connect(tok, subOf("olga"), true)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := c.nc.Request("svc.policy.command.propose", []byte(`{"v":1}`), 2*time.Second); err != nil {
		t.Fatalf("operator propose: %v", err)
	}
	infos, _ := policy.list()
	if len(infos) != 1 {
		t.Fatalf("policy saw %d", len(infos))
	}
	ci := infos[0]
	// Groups travel in the server-set user name (tags are dropped in config mode, see
	// PrincipalName); non-OneMind groups are not included.
	if ci["acc"] != "CLIENT" || ci["user"] != "ak:olga;onemind-operators,onemind-viewers" {
		t.Fatalf("request info acc/user: %v", ci)
	}
	if _, ok := ci["tags"]; ok {
		t.Logf("server now forwards tags in config mode: %v", ci["tags"])
	}
	// A client cannot forge the header: the server replaces it on the import.
	m := nats.NewMsg("svc.policy.command.propose")
	m.Header.Set("Nats-Request-Info", `{"acc":"CLIENT","user":"ak:zeus","tags":["role:admin"]}`)
	m.Data = []byte(`{"v":1}`)
	if _, err := c.nc.RequestMsg(m, 2*time.Second); err != nil {
		t.Fatal(err)
	}
	infos, _ = policy.list()
	if infos[1]["user"] != ci["user"] || strings.Contains(fmt.Sprint(infos[1]), "role:admin") {
		t.Fatalf("forged Nats-Request-Info was kept: %v", infos[1])
	}
}

func TestInboxIsolationBetweenClients(t *testing.T) {
	w := newWorld(t, worldOpts{})
	pol := w.password("policy")
	w.serve(pol, "svc.policy.>")
	op, err := w.connect(w.idp.token(t, "olga", "onemind-operators"), subOf("olga"), false)
	if err != nil {
		t.Fatal(err)
	}
	nw, err := w.connect(w.idp.token(t, "nina", "onemind-newcomers"), subOf("nina"), false)
	if err != nil {
		t.Fatal(err)
	}
	// The newcomer tries to read the operator's replies.
	_, _ = nw.nc.SubscribeSync(w.roles.InboxPrefix + "." + subOf("olga") + ".>")
	_ = nw.nc.Flush()
	waitFor(t, func() bool { return len(nw.permErrors()) == 1 }, "subscription refused")
	if _, err := op.nc.Request("svc.policy.action.get", []byte(`{}`), 2*time.Second); err != nil {
		t.Fatal(err)
	}
	// A client that does not use its own inbox prefix gets no replies.
	bad, err := w.connect(w.idp.token(t, "olga", "onemind-operators"), "", false)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := bad.nc.Request("svc.policy.action.get", []byte(`{}`), 300*time.Millisecond); err == nil {
		t.Fatal("a reply arrived on the shared _INBOX prefix")
	}
}

func TestRejectedTokens(t *testing.T) {
	w := newWorld(t, worldOpts{})
	p := w.idp
	exp := claims("ed", -time.Minute, "onemind-operators")
	exp["exp"] = time.Now().Add(-2 * time.Minute).Unix()
	wrongIss := claims("ed", time.Minute, "onemind-operators")
	wrongIss["iss"] = "https://evil.test/"
	wrongAud := claims("ed", time.Minute, "onemind-operators")
	wrongAud["aud"] = "argocd"
	badUser := claims("Ed Smith", time.Minute, "onemind-operators")
	noGroup := claims("ed", time.Minute, "argocd-admins")
	good := p.token(t, "ed", "onemind-operators")
	parts := strings.Split(good, ".")
	tampered := parts[0] + "." + b64.EncodeToString([]byte(strings.Replace(string(mustB64(t, parts[1])), "onemind-operators", "onemind-admins", 1))) + "." + parts[2]
	cases := map[string]string{
		"expired":        p.sign(t, "RS256", p.rsaKid, exp),
		"wrong issuer":   p.sign(t, "RS256", p.rsaKid, wrongIss),
		"wrong audience": p.sign(t, "RS256", p.rsaKid, wrongAud),
		"alg none":       p.sign(t, "none", p.rsaKid, claims("ed", time.Minute, "onemind-operators")),
		"HS256":          p.sign(t, "HS256", p.rsaKid, claims("ed", time.Minute, "onemind-operators")),
		"unknown kid":    p.sign(t, "RS256", "nope", claims("ed", time.Minute, "onemind-operators")),
		"bad username":   p.sign(t, "RS256", p.rsaKid, badUser),
		"no role group":  p.sign(t, "RS256", p.rsaKid, noGroup),
		"tampered":       tampered,
		"garbage":        "not-a-jwt",
		"empty":          "",
	}
	for name, tok := range cases {
		_, err := w.connect(tok, "", name == "expired")
		if !isAuthErr(err) {
			t.Errorf("%s: want authorization violation, got %v", name, err)
		}
	}
	if w.svc.Allowed.Load() != 0 {
		t.Fatalf("allowed %d", w.svc.Allowed.Load())
	}
	// ES256 tokens are accepted too.
	if _, err := w.connect(p.sign(t, "ES256", p.ecKid, claims("ed", time.Minute, "onemind-viewers")), subOf("ed"), false); err != nil {
		t.Fatalf("ES256: %v", err)
	}
}

func mustB64(t *testing.T, s string) []byte {
	b, err := b64.DecodeString(s)
	if err != nil {
		t.Fatal(err)
	}
	return b
}

// Service and leaf logins with passwords never reach the callout, so turning the callout on
// does not touch them (allowed_accounts: [AUTH]).
func TestPasswordUsersAreNotCalledOut(t *testing.T) {
	w := newWorld(t, worldOpts{})
	w.password("policy")
	zeus := w.password("zeus")
	if err := zeus.Publish("cmd.map.viewport.set", []byte(`{}`)); err != nil {
		t.Fatal(err)
	}
	_ = zeus.Flush()
	if n := w.svc.Requests.Load(); n != 0 {
		t.Fatalf("callout received %d requests for password users", n)
	}
	// A wrong password falls through to the callout, which refuses it (no token).
	_, err := nats.Connect(w.url, nats.UserInfo("policy", "wrong"), nats.NoReconnect())
	if !isAuthErr(err) {
		t.Fatalf("wrong password: %v", err)
	}
	if w.svc.Requests.Load() != 1 || w.svc.Denied.Load() != 1 {
		t.Fatalf("requests %d denied %d", w.svc.Requests.Load(), w.svc.Denied.Load())
	}
}

// Leaf connections (TAK leaf, Hermes sidecar, Mac site) and the callout.
func TestLeafNodes(t *testing.T) {
	// Pins nats-server v2.14.2 behaviour (server/auth.go): a leaf that logs in through
	// leafnodes.authorization never reaches the callout; a leaf that logs in as an account user
	// is always delegated (the allowed_accounts filter excludes LEAF), and the callout refuses
	// it. So the callout is compatible with today's config and A19 stage A1, not with stage A2
	// (proposal J2-P1).
	for _, tc := range []struct {
		name      string
		leafAuth  bool
		user      string
		connected bool
	}{
		{"leafnodes.authorization user (today, stage A1)", true, "leaf_tak", true},
		{"account user with LEAFNODE type (stage A2)", false, "leaf-site", false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			w := newWorld(t, worldOpts{leafAuth: tc.leafAuth})
			u, _ := url.Parse(w.leafURL)
			u.User = url.UserPassword(tc.user, w.pw[tc.user])
			dir := t.TempDir()
			cfg := filepath.Join(dir, "leaf.conf")
			conf := fmt.Sprintf("listen: \"127.0.0.1:-1\"\nserver_name: spoke\nleafnodes { remotes: [ { url: %q } ] }\n", u.String())
			if err := os.WriteFile(cfg, []byte(conf), 0o600); err != nil {
				t.Fatal(err)
			}
			opts, err := server.ProcessConfigFile(cfg)
			if err != nil {
				t.Fatal(err)
			}
			opts.NoLog, opts.NoSigs = true, true
			spoke, err := server.NewServer(opts)
			if err != nil {
				t.Fatal(err)
			}
			go spoke.Start()
			t.Cleanup(func() { spoke.Shutdown(); spoke.WaitForShutdown() })
			deadline := time.Now().Add(3 * time.Second)
			for time.Now().Before(deadline) && w.srv.NumLeafNodes() == 0 {
				time.Sleep(20 * time.Millisecond)
			}
			t.Logf("leaf connected=%v, callout requests=%d allowed=%d denied=%d",
				w.srv.NumLeafNodes() == 1, w.svc.Requests.Load(), w.svc.Allowed.Load(), w.svc.Denied.Load())
			if got := w.srv.NumLeafNodes() == 1; got != tc.connected {
				t.Fatalf("leaf connected=%v, want %v", got, tc.connected)
			}
			if tc.connected && w.svc.Requests.Load() != 0 {
				t.Fatalf("callout consulted for a leafnodes.authorization leaf")
			}
			if !tc.connected && (w.svc.Requests.Load() == 0 || w.svc.Allowed.Load() != 0) {
				t.Fatalf("stage A2 leaf: requests %d allowed %d", w.svc.Requests.Load(), w.svc.Allowed.Load())
			}
		})
	}
}

func TestTokenExpiryDisconnects(t *testing.T) {
	w := newWorld(t, worldOpts{})
	c := claims("tim", 0, "onemind-viewers")
	c["exp"] = time.Now().Add(2 * time.Second).Unix()
	// The verifier allows 30 s of clock skew on exp; the user JWT itself carries the real exp.
	tok := w.idp.sign(t, "RS256", w.idp.rsaKid, c)
	closed := make(chan struct{})
	cl, err := w.connect(tok, subOf("tim"), false, nats.ClosedHandler(func(*nats.Conn) { close(closed) }))
	if err != nil {
		t.Fatal(err)
	}
	select {
	case <-closed:
	case <-time.After(8 * time.Second):
		t.Fatal("connection outlived its token")
	}
	if !strings.Contains(strings.ToLower(fmt.Sprint(cl.nc.LastError())), "expired") {
		t.Logf("close reason: %v", cl.nc.LastError())
	}
}

func TestCalloutDownFailsClosed(t *testing.T) {
	w := newWorld(t, worldOpts{})
	w.svc.Stop()
	_ = w.svc.cfg.Conn.Flush()
	start := time.Now()
	_, err := w.connect(w.idp.token(t, "olga", "onemind-operators"), subOf("olga"), false)
	if err == nil {
		t.Fatal("connected with no callout service")
	}
	t.Logf("refused after %v: %v", time.Since(start).Round(time.Millisecond), err)
	// Password users are unaffected.
	w.password("policy")
}

func TestJWKSOutageAndRotation(t *testing.T) {
	w := newWorld(t, worldOpts{})
	// First login fetches the key set.
	if _, err := w.connect(w.idp.token(t, "olga", "onemind-operators"), subOf("olga"), false); err != nil {
		t.Fatal(err)
	}
	// Authentik down: cached keys keep working.
	w.idp.setDown(true)
	if _, err := w.connect(w.idp.token(t, "olga", "onemind-operators"), subOf("olga"), false); err != nil {
		t.Fatalf("cached keys: %v", err)
	}
	w.idp.setDown(false)
	// Key rotation: a token signed by a new kid triggers one refetch.
	w.idp.mu.Lock()
	w.idp.rsaKid = "rsa-2"
	w.idp.mu.Unlock()
	before := w.idp.fetches.Load()
	w.svc.cfg.Verifier.mu.Lock()
	w.svc.cfg.Verifier.tried = time.Time{}.Add(time.Hour) // allow an immediate refetch
	w.svc.cfg.Verifier.mu.Unlock()
	if _, err := w.connect(w.idp.token(t, "olga", "onemind-operators"), subOf("olga"), false); err != nil {
		t.Fatalf("rotated key: %v", err)
	}
	if w.idp.fetches.Load() != before+1 {
		t.Fatalf("fetches %d -> %d", before, w.idp.fetches.Load())
	}
}

func TestEncryptedCallout(t *testing.T) {
	w := newWorld(t, worldOpts{xkey: true})
	pol := w.password("policy")
	policy := w.serve(pol, "svc.policy.>")
	c, err := w.connect(w.idp.token(t, "cora", "onemind-commanders"), subOf("cora"), true)
	if err != nil {
		t.Fatalf("encrypted callout: %v", err)
	}
	if _, err := c.nc.Request("svc.policy.task.confirm", []byte(`{}`), 2*time.Second); err != nil {
		t.Fatal(err)
	}
	infos, _ := policy.list()
	if infos[0]["user"] != "ak:cora;onemind-commanders" {
		t.Fatalf("info %v", infos[0])
	}
}

func TestUserJWTShape(t *testing.T) {
	w := newWorld(t, worldOpts{})
	rm := w.roles
	for _, r := range rm.Roles {
		t.Run(r.Name, func(t *testing.T) {
			c, err := w.connect(w.idp.token(t, "u-"+r.Name, r.Group), subOf("u-"+r.Name), false)
			if err != nil {
				t.Fatal(err)
			}
			// Connection-level view from the server: account and name.
			cz, err := w.srv.Connz(&server.ConnzOptions{Username: true, Subscriptions: true})
			if err != nil {
				t.Fatal(err)
			}
			found := false
			for _, ci := range cz.Conns {
				if ci.AuthorizedUser == "ak:u-"+r.Name+";"+r.Group {
					found = true
					if ci.Account != "CLIENT" {
						t.Fatalf("account %s", ci.Account)
					}
				}
			}
			if !found {
				b, _ := json.Marshal(cz.Conns)
				t.Fatalf("connection not found as ak:u-%s: %s", r.Name, b)
			}
			c.nc.Close()
		})
	}
}
