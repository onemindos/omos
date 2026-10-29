// Package callout is the NATS auth callout for browser and remote clients (FEATURES J2).
//
// A client connects (WebSocket or plain NATS) with its Authentik OIDC token as the CONNECT
// auth_token. nats-server (config mode, `authorization { auth_callout { ... } }`) sends an
// authorization request on $SYS.REQ.USER.AUTH in the AUTH account; this service verifies the
// token against Authentik's JWKS, maps the user's Authentik groups to one OneMind role
// (platform/security/roles/onemind-roles.json) and answers with a short-lived user JWT that
// places the connection in the CLIENT account with that role's allow lists, the name
// "ak:<username>" and the tags "role:<role>" and "group:<group>". The policy service reads name
// and tags from the server-set Nats-Request-Info header to resolve the principal and its
// OpenFGA group memberships. Server behaviour cited: nats-server v2.14.2 server/auth.go
// (processClientOrLeafAuthentication) and server/auth_callout.go.
package callout

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"regexp"
	"strings"
	"sync/atomic"
	"time"

	"github.com/nats-io/jwt/v2"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nkeys"
)

// AuthSubject is the server's authorization request subject (server/auth_callout.go
// AuthCalloutSubject). It is a NATS system subject, not an application subject.
const AuthSubject = "$SYS.REQ.USER.AUTH"

// xkeyHeader carries the server's curve public key when requests are encrypted
// (server/auth_callout.go AuthRequestXKeyHeader).
const xkeyHeader = "Nats-Server-Xkey"

// Config wires the service.
type Config struct {
	Conn     *nats.Conn
	Issuer   nkeys.KeyPair // account key pair named as `issuer` in auth_callout
	XKey     nkeys.KeyPair // optional curve key pair named as `xkey` in auth_callout
	Verifier *Verifier
	Roles    *RoleMap
	Log      *slog.Logger
	Now      func() time.Time
	// Timeout bounds token verification (JWKS fetch included); the server's own auth
	// timeout is the outer limit.
	Timeout time.Duration
}

// Service answers authorization requests.
type Service struct {
	cfg      Config
	sub      *nats.Subscription
	Allowed  atomic.Int64
	Denied   atomic.Int64
	Requests atomic.Int64
}

// New validates the configuration.
func New(cfg Config) (*Service, error) {
	if cfg.Conn == nil || cfg.Issuer == nil || cfg.Verifier == nil || cfg.Roles == nil {
		return nil, errors.New("callout: conn, issuer, verifier and roles are required")
	}
	pub, err := cfg.Issuer.PublicKey()
	if err != nil || !nkeys.IsValidPublicAccountKey(pub) {
		return nil, errors.New("callout: issuer must be an account key pair (seed SA...)")
	}
	if cfg.XKey != nil {
		if xp, err := cfg.XKey.PublicKey(); err != nil || !nkeys.IsValidPublicCurveKey(xp) {
			return nil, errors.New("callout: xkey must be a curve key pair (seed SX...)")
		}
	}
	if cfg.Log == nil {
		cfg.Log = slog.New(slog.DiscardHandler)
	}
	if cfg.Now == nil {
		cfg.Now = time.Now
	}
	if cfg.Timeout == 0 {
		cfg.Timeout = 1500 * time.Millisecond
	}
	return &Service{cfg: cfg}, nil
}

// Start subscribes to the authorization subject. Requests are answered one at a time per
// connection; run more replicas (queue group) for more throughput.
func (s *Service) Start() error {
	sub, err := s.cfg.Conn.QueueSubscribe(AuthSubject, "auth-callout", s.handle)
	if err != nil {
		return err
	}
	s.sub = sub
	return s.cfg.Conn.Flush()
}

// Stop unsubscribes.
func (s *Service) Stop() {
	if s.sub != nil {
		_ = s.sub.Unsubscribe()
	}
}

var usernameRE = regexp.MustCompile(`^[a-z0-9][a-z0-9._@-]{0,63}$`)

// Decision is the outcome for one request, for logs and tests.
type Decision struct {
	UserNkey string
	Username string
	Role     string
	Groups   []string
	Expires  time.Time
	Reason   string // empty when allowed
}

func (s *Service) handle(m *nats.Msg) {
	s.Requests.Add(1)
	data := m.Data
	serverX := ""
	if m.Header != nil {
		serverX = m.Header.Get(xkeyHeader)
	}
	if serverX != "" {
		if s.cfg.XKey == nil {
			s.cfg.Log.Error("encrypted authorization request but no xkey configured")
			return // the server times out and refuses the client
		}
		var err error
		if data, err = s.cfg.XKey.Open(m.Data, serverX); err != nil {
			s.cfg.Log.Error("cannot decrypt authorization request", "err", err)
			return
		}
	}
	rc, err := jwt.DecodeAuthorizationRequestClaims(string(data))
	if err != nil {
		s.cfg.Log.Error("bad authorization request", "err", err)
		return
	}
	d, userJWT := s.decide(rc)
	if d.Reason != "" {
		s.Denied.Add(1)
		s.cfg.Log.Info("denied", "user_nkey", rc.UserNkey, "username", d.Username, "reason", d.Reason,
			"client_type", rc.ClientInformation.Type, "host", rc.ClientInformation.Host)
	} else {
		s.Allowed.Add(1)
		s.cfg.Log.Info("allowed", "username", d.Username, "role", d.Role, "groups", d.Groups,
			"expires", d.Expires.UTC().Format(time.RFC3339), "client_type", rc.ClientInformation.Type)
	}
	resp := jwt.NewAuthorizationResponseClaims(rc.UserNkey)
	resp.Audience = rc.Server.ID
	if d.Reason != "" {
		resp.Error = "not authorized" // no detail to the client: no oracle for token probing
	} else {
		resp.Jwt = userJWT
	}
	out, err := resp.Encode(s.cfg.Issuer)
	if err != nil {
		s.cfg.Log.Error("cannot sign authorization response", "err", err)
		return
	}
	body := []byte(out)
	if serverX != "" {
		if body, err = s.cfg.XKey.Seal(body, serverX); err != nil {
			s.cfg.Log.Error("cannot encrypt authorization response", "err", err)
			return
		}
	}
	_ = m.Respond(body)
}

// decide verifies the token and builds the user JWT. A non-empty Reason means deny.
func (s *Service) decide(rc *jwt.AuthorizationRequestClaims) (Decision, string) {
	d := Decision{UserNkey: rc.UserNkey}
	deny := func(format string, a ...any) (Decision, string) {
		d.Reason = fmt.Sprintf(format, a...)
		return d, ""
	}
	if !nkeys.IsValidPublicUserKey(rc.UserNkey) {
		return deny("request has no valid user nkey")
	}
	switch rc.ClientInformation.Kind {
	case "Client":
	default:
		return deny("connection kind %q is not delegated to the callout", rc.ClientInformation.Kind)
	}
	token := rc.ConnectOptions.Token
	if token == "" {
		return deny("no auth_token (password or nkey logins are not handled by the callout)")
	}
	ctx, cancel := context.WithTimeout(context.Background(), s.cfg.Timeout)
	defer cancel()
	id, err := s.cfg.Verifier.Verify(ctx, token)
	if err != nil {
		return deny("%v", err)
	}
	d.Username = id.Username
	if !usernameRE.MatchString(id.Username) {
		return deny("preferred_username %q is not usable as a principal id", id.Username)
	}
	if !tokenRE.MatchString(id.Subject) || len(id.Subject) > 128 {
		return deny("sub is not usable as an inbox token")
	}
	role, held, ok := s.cfg.Roles.Resolve(id.Groups)
	if !ok {
		return deny("no OneMind role group in token")
	}
	d.Role, d.Groups = role.Name, held

	now := s.cfg.Now()
	exp := id.Expires
	if max := now.Add(time.Duration(s.cfg.Roles.MaxTTLSeconds) * time.Second); exp.After(max) {
		exp = max
	}
	if !exp.After(now) {
		return deny("token expires now")
	}
	d.Expires = exp

	uc := jwt.NewUserClaims(rc.UserNkey)
	uc.Name = PrincipalName(s.cfg.Roles, id.Username, held)
	uc.Audience = s.cfg.Roles.ClientAccount // config mode: placement account
	uc.IssuedAt = now.Unix()
	uc.Expires = exp.Unix()
	uc.Tags.Add(s.cfg.Roles.RoleTagPrefix + role.Name)
	for _, g := range held {
		uc.Tags.Add(s.cfg.Roles.GroupTagPrefix + g)
	}
	uc.AllowedConnectionTypes = jwt.StringList(append([]string(nil), s.cfg.Roles.ConnectionTypes...))
	uc.Pub.Allow.Add(role.Publish...)
	if len(role.Publish) == 0 {
		uc.Pub.Deny.Add(">")
	}
	uc.Sub.Allow.Add(InboxFor(s.cfg.Roles, id.Subject) + ".>")
	uc.Sub.Allow.Add(role.Subscribe...)
	// Defence in depth on top of the allow lists.
	uc.Pub.Deny.Add("$JS.>", "$KV.>", "$SYS.>", "_INBOX.>")
	vr := jwt.CreateValidationResults()
	uc.Validate(vr)
	if vr.IsBlocking(true) {
		return deny("user claims invalid: %v", vr.Errors())
	}
	enc, err := uc.Encode(s.cfg.Issuer)
	if err != nil {
		return deny("cannot sign user JWT: %v", err)
	}
	return d, enc
}

// PrincipalName is the NATS user name the callout gives a connection:
// "<name_prefix><username>;<group>,<group>", for example "ak:olga;onemind-operators".
//
// In config mode nats-server copies only the user name into Nats-Request-Info: user JWT tags
// are applied in operator mode only (nats-server v2.14.2 server/auth.go, `c.tags = juc.Tags`
// on the JWT path; the callout path builds an internal nkey user without them). So the
// OneMind role groups travel in the server-set name, which a client cannot forge. Neither
// the username nor a group name can contain ';' or ',' (usernameRE, groupRE).
func PrincipalName(m *RoleMap, username string, groups []string) string {
	return m.NamePrefix + username + ";" + strings.Join(groups, ",")
}

// InboxFor is the inbox prefix a client of this subject must use (nats.go
// CustomInboxPrefix, nats.ws `inboxPrefix`). Per-user prefixes stop one CLIENT user from
// reading another user's replies (A19-P5).
func InboxFor(m *RoleMap, subject string) string {
	return m.InboxPrefix + "." + subject
}
