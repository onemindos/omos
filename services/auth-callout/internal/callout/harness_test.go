package callout

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/nats-io/nats-server/v2/server"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nkeys"
)

// Every key and password in these tests is generated per run; none is committed.

const (
	roleSource = "../../../../platform/security/roles/onemind-roles.json"
	testIss    = "https://auth.example.test/application/o/onemind-xr/"
	testAud    = "onemind-xr"
)

// idp is a stand-in for Authentik: an RSA and a P-256 signing key behind a JWKS endpoint.
type idp struct {
	rsa     *rsa.PrivateKey
	ec      *ecdsa.PrivateKey
	rsaKid  string
	ecKid   string
	srv     *httptest.Server
	fetches atomic.Int64
	mu      sync.Mutex
	down    bool
	extra   []map[string]string
}

func newIDP(t *testing.T) *idp {
	t.Helper()
	rk, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatal(err)
	}
	ek, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	p := &idp{rsa: rk, ec: ek, rsaKid: "rsa-1", ecKid: "ec-1"}
	p.srv = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		p.fetches.Add(1)
		p.mu.Lock()
		defer p.mu.Unlock()
		if p.down {
			http.Error(w, "down", http.StatusServiceUnavailable)
			return
		}
		_ = json.NewEncoder(w).Encode(map[string]any{"keys": p.jwks()})
	}))
	t.Cleanup(p.srv.Close)
	return p
}

func (p *idp) jwks() []map[string]string {
	ecb, _ := p.ec.PublicKey.Bytes() // uncompressed point 04||X||Y
	keys := []map[string]string{
		{"kty": "RSA", "kid": p.rsaKid, "use": "sig", "alg": "RS256",
			"n": b64.EncodeToString(p.rsa.N.Bytes()), "e": b64.EncodeToString(big.NewInt(int64(p.rsa.E)).Bytes())},
		{"kty": "EC", "kid": p.ecKid, "use": "sig", "crv": "P-256",
			"x": b64.EncodeToString(ecb[1:33]), "y": b64.EncodeToString(ecb[33:])},
	}
	return append(keys, p.extra...)
}

func (p *idp) setDown(d bool) { p.mu.Lock(); p.down = d; p.mu.Unlock() }

// claims returns a valid claim set for user with groups, valid for ttl.
func claims(user string, ttl time.Duration, groups ...string) map[string]any {
	now := time.Now()
	sum := sha256.Sum256([]byte(user))
	return map[string]any{"iss": testIss, "aud": testAud, "sub": hex.EncodeToString(sum[:]),
		"preferred_username": user, "groups": groups, "iat": now.Unix(), "exp": now.Add(ttl).Unix()}
}

func (p *idp) sign(t *testing.T, alg, kid string, c map[string]any) string {
	t.Helper()
	hb, _ := json.Marshal(map[string]string{"alg": alg, "kid": kid, "typ": "JWT"})
	pb, _ := json.Marshal(c)
	in := b64.EncodeToString(hb) + "." + b64.EncodeToString(pb)
	d := sha256.Sum256([]byte(in))
	var sig []byte
	switch alg {
	case "RS256":
		s, err := rsa.SignPKCS1v15(rand.Reader, p.rsa, 5, d[:]) // 5 = crypto.SHA256
		if err != nil {
			t.Fatal(err)
		}
		sig = s
	case "ES256":
		r, s, err := ecdsa.Sign(rand.Reader, p.ec, d[:])
		if err != nil {
			t.Fatal(err)
		}
		sig = append(r.FillBytes(make([]byte, 32)), s.FillBytes(make([]byte, 32))...)
	default:
		sig = []byte("x")
	}
	return in + "." + b64.EncodeToString(sig)
}

func (p *idp) token(t *testing.T, user string, groups ...string) string {
	return p.sign(t, "RS256", p.rsaKid, claims(user, 10*time.Minute, groups...))
}

func randPW() string {
	b := make([]byte, 18)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}

// world is one disposable nats-server with the stage-1 account layout (A19, #562) plus the
// AUTH account and the auth_callout block from platform/auth-callout/nats-callout.conf.
type world struct {
	t       *testing.T
	srv     *server.Server
	url     string
	wsURL   string
	leafURL string
	pw      map[string]string
	issuer  nkeys.KeyPair
	xkey    nkeys.KeyPair
	idp     *idp
	svc     *Service
	roles   *RoleMap
	conns   []*nats.Conn
}

type worldOpts struct {
	xkey     bool
	leafAuth bool // leafnodes { authorization { users } } as today and at stage A1
}

func natsConf(dir string, w *world, o worldOpts) string {
	ipub, _ := w.issuer.PublicKey()
	xline := ""
	if w.xkey != nil {
		xp, _ := w.xkey.PublicKey()
		xline = "xkey: " + xp
	}
	leafAuth := ""
	if o.leafAuth {
		leafAuth = fmt.Sprintf(`authorization { users: [ { user: leaf_tak, password: %q, account: ONEMIND } ] }`, w.pw["leaf_tak"])
	}
	return fmt.Sprintf(`
listen: "127.0.0.1:-1"
server_name: j2-test
websocket { listen: "127.0.0.1:-1", no_tls: true }
leafnodes {
  listen: "127.0.0.1:-1"
  %s
}
jetstream { store_dir: %q, max_mem: 32MB, max_file: 64MB }
system_account: SYS
accounts {
  SYS { users: [ { user: sys, password: %q } ] }
  # Least privilege for the callout (platform/auth-callout/identities.fragment.yaml):
  # it can only receive authorization requests and answer them.
  AUTH { users: [ { user: auth-callout, password: %q, allowed_connection_types: [STANDARD], permissions: {
      publish: { deny: [">"] }, subscribe: { allow: ["$SYS.REQ.USER.AUTH"] }, allow_responses: { max: 1, expires: "5s" } } } ] }
  ONEMIND {
    jetstream: enabled
    users: [
      { user: policy, password: %q }
      { user: leaf-site, password: %q, allowed_connection_types: [LEAFNODE] }
    ]
    exports: [
      { service: "svc.policy.>" }
      { service: "svc.entity.>" }
      { service: "svc.layer.>", response_type: stream }
      { service: "svc.cot.send" }
      { stream: "cmd.map.>" }
      { stream: "evt.so.action.>" }
    ]
  }
  CLIENT {
    users: [
      { user: zeus, password: %q, permissions: {
          publish: { allow: ["svc.policy.>", "cmd.map.>"] }, subscribe: { allow: ["_INBOX.>", "cmd.map.>"] } } }
    ]
    imports: [
      { service: { account: ONEMIND, subject: "svc.policy.>" }, share: true }
      { service: { account: ONEMIND, subject: "svc.entity.>" }, share: true }
      { service: { account: ONEMIND, subject: "svc.layer.>" }, share: true }
      { service: { account: ONEMIND, subject: "svc.cot.send" }, share: true }
      { stream: { account: ONEMIND, subject: "cmd.map.>" } }
      { stream: { account: ONEMIND, subject: "evt.so.action.>" } }
    ]
  }
}
authorization {
  auth_callout {
    issuer: %s
    auth_users: [ auth-callout ]
    account: AUTH
    allowed_accounts: [ AUTH ]
    %s
  }
}
`, leafAuth, filepath.Join(dir, "js"), w.pw["sys"], w.pw["auth-callout"], w.pw["policy"], w.pw["leaf-site"], w.pw["zeus"], ipub, xline)
}

func newWorld(t *testing.T, o worldOpts) *world {
	t.Helper()
	w := &world{t: t, pw: map[string]string{}}
	for _, u := range []string{"sys", "auth-callout", "policy", "leaf-site", "leaf_tak", "zeus"} {
		w.pw[u] = randPW()
	}
	var err error
	if w.issuer, err = nkeys.CreateAccount(); err != nil {
		t.Fatal(err)
	}
	if o.xkey {
		if w.xkey, err = nkeys.CreateCurveKeys(); err != nil {
			t.Fatal(err)
		}
	}
	dir := t.TempDir()
	cfg := filepath.Join(dir, "nats.conf")
	if err := os.WriteFile(cfg, []byte(natsConf(dir, w, o)), 0o600); err != nil {
		t.Fatal(err)
	}
	opts, err := server.ProcessConfigFile(cfg)
	if err != nil {
		t.Fatal(err)
	}
	opts.NoLog, opts.NoSigs = true, true
	opts.AuthTimeout = 2
	w.srv, err = server.NewServer(opts)
	if err != nil {
		t.Fatal(err)
	}
	go w.srv.Start()
	if !w.srv.ReadyForConnections(10 * time.Second) {
		t.Fatal("nats-server not ready")
	}
	w.url = w.srv.ClientURL()
	w.wsURL = w.srv.WebsocketURL()
	w.leafURL = fmt.Sprintf("nats://127.0.0.1:%d", opts.LeafNode.Port)
	t.Cleanup(func() {
		for _, c := range w.conns {
			c.Close()
		}
		w.srv.Shutdown()
		w.srv.WaitForShutdown()
	})
	w.idp = newIDP(t)
	if w.roles, err = LoadRoleMap(roleSource); err != nil {
		t.Fatal(err)
	}
	w.startCallout()
	return w
}

func (w *world) startCallout() {
	w.t.Helper()
	nc, err := nats.Connect(w.url, nats.UserInfo("auth-callout", w.pw["auth-callout"]), nats.Name("auth-callout"))
	if err != nil {
		w.t.Fatal(err)
	}
	w.conns = append(w.conns, nc)
	ver, err := NewVerifier(OIDCConfig{Issuer: testIss, Audience: testAud, JWKSURL: w.idp.srv.URL})
	if err != nil {
		w.t.Fatal(err)
	}
	w.svc, err = New(Config{Conn: nc, Issuer: w.issuer, XKey: w.xkey, Verifier: ver, Roles: w.roles})
	if err != nil {
		w.t.Fatal(err)
	}
	if err := w.svc.Start(); err != nil {
		w.t.Fatal(err)
	}
	w.t.Cleanup(w.svc.Stop)
}

type client struct {
	nc   *nats.Conn
	mu   sync.Mutex
	errs []error
}

func (c *client) permErrors() []error {
	c.mu.Lock()
	defer c.mu.Unlock()
	return append([]error(nil), c.errs...)
}

// connect logs in with an OIDC token (as a browser does over WebSocket when ws is true).
func (w *world) connect(token, inboxSub string, ws bool, extra ...nats.Option) (*client, error) {
	c := &client{}
	url := w.url
	if ws {
		url = w.wsURL
	}
	opts := append([]nats.Option{nats.Token(token), nats.NoReconnect(), nats.Timeout(5 * time.Second),
		nats.ErrorHandler(func(_ *nats.Conn, _ *nats.Subscription, err error) {
			c.mu.Lock()
			c.errs = append(c.errs, err)
			c.mu.Unlock()
		})}, extra...)
	if inboxSub != "" {
		opts = append(opts, nats.CustomInboxPrefix(w.roles.InboxPrefix+"."+inboxSub))
	}
	nc, err := nats.Connect(url, opts...)
	if err != nil {
		return nil, err
	}
	w.conns = append(w.conns, nc)
	c.nc = nc
	return c, nil
}

func subOf(user string) string {
	sum := sha256.Sum256([]byte(user))
	return hex.EncodeToString(sum[:])
}

func (w *world) password(user string) *nats.Conn {
	w.t.Helper()
	nc, err := nats.Connect(w.url, nats.UserInfo(user, w.pw[user]), nats.Name(user))
	if err != nil {
		w.t.Fatalf("connect %s: %v", user, err)
	}
	w.conns = append(w.conns, nc)
	return nc
}

// responder serves a subject in ONEMIND as the policy user and records the verified caller.
type seen struct {
	mu   sync.Mutex
	info []map[string]any
	subj []string
}

func (s *seen) list() ([]map[string]any, []string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return append([]map[string]any(nil), s.info...), append([]string(nil), s.subj...)
}

func (w *world) serve(nc *nats.Conn, subject string) *seen {
	w.t.Helper()
	s := &seen{}
	if _, err := nc.Subscribe(subject, func(m *nats.Msg) {
		var ci map[string]any
		if m.Header != nil {
			_ = json.Unmarshal([]byte(m.Header.Get("Nats-Request-Info")), &ci)
		}
		s.mu.Lock()
		s.info = append(s.info, ci)
		s.subj = append(s.subj, m.Subject)
		s.mu.Unlock()
		_ = m.Respond([]byte(`{"ok":true}`))
	}); err != nil {
		w.t.Fatal(err)
	}
	_ = nc.Flush()
	return s
}

func waitFor(t *testing.T, cond func() bool, what string) {
	t.Helper()
	deadline := time.Now().Add(5 * time.Second)
	for time.Now().Before(deadline) {
		if cond() {
			return
		}
		time.Sleep(10 * time.Millisecond)
	}
	t.Fatalf("timed out waiting for %s", what)
}

func isAuthErr(err error) bool {
	return err != nil && (strings.Contains(strings.ToLower(err.Error()), "authorization violation") ||
		strings.Contains(strings.ToLower(err.Error()), "authentication"))
}

var _ = base64.RawURLEncoding
