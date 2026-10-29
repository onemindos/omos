package callout

import (
	"context"
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"strings"
	"sync"
	"time"
)

// OIDCConfig describes the Authentik OAuth2 provider whose tokens the callout accepts.
type OIDCConfig struct {
	// Issuer is the exact iss claim, e.g. https://auth.onemindos.dev/application/o/onemind-xr/
	Issuer string
	// Audience is the client_id the token must be issued to (aud claim).
	Audience string
	// JWKSURL is where the signing keys are fetched (in-cluster Authentik service URL).
	JWKSURL string
	// Leeway tolerates clock skew on exp, nbf and iat.
	Leeway time.Duration
	// MinRefresh limits JWKS refetches triggered by an unknown kid.
	MinRefresh time.Duration
	// CacheTTL is how long a fetched key set is used before a background refetch.
	CacheTTL time.Duration
	HTTP     *http.Client
	Now      func() time.Time
}

// Identity is what the callout takes from a verified token.
type Identity struct {
	Subject  string // sub (Authentik "hashed user id" mode: 64 hex chars)
	Username string // preferred_username
	Groups   []string
	Expires  time.Time
}

// Verifier validates Authentik-issued JWTs against the provider's JWKS.
type Verifier struct {
	cfg     OIDCConfig
	mu      sync.Mutex
	keys    map[string]crypto.PublicKey
	fetched time.Time
	tried   time.Time
}

// NewVerifier returns a verifier with defaults filled in.
func NewVerifier(cfg OIDCConfig) (*Verifier, error) {
	if cfg.Issuer == "" || cfg.Audience == "" || cfg.JWKSURL == "" {
		return nil, errors.New("oidc: issuer, audience and jwks url are required")
	}
	if cfg.Leeway == 0 {
		cfg.Leeway = 30 * time.Second
	}
	if cfg.MinRefresh == 0 {
		cfg.MinRefresh = 30 * time.Second
	}
	if cfg.CacheTTL == 0 {
		cfg.CacheTTL = time.Hour
	}
	if cfg.HTTP == nil {
		cfg.HTTP = &http.Client{Timeout: 5 * time.Second}
	}
	if cfg.Now == nil {
		cfg.Now = time.Now
	}
	return &Verifier{cfg: cfg, keys: map[string]crypto.PublicKey{}}, nil
}

// ErrToken is returned for every token problem; the detail is for logs only.
var ErrToken = errors.New("invalid token")

func tokenErr(format string, a ...any) error {
	return fmt.Errorf("%w: %s", ErrToken, fmt.Sprintf(format, a...))
}

var b64 = base64.RawURLEncoding

// Verify checks signature (RS256 or ES256 only), iss, aud, exp, nbf and iat, and returns
// the identity. Any failure returns an error wrapping ErrToken.
func (v *Verifier) Verify(ctx context.Context, token string) (Identity, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 || len(token) > 16384 {
		return Identity{}, tokenErr("not a compact JWS")
	}
	hb, err := b64.DecodeString(parts[0])
	if err != nil {
		return Identity{}, tokenErr("header encoding")
	}
	var hdr struct {
		Alg  string `json:"alg"`
		Kid  string `json:"kid"`
		Crit []any  `json:"crit"`
	}
	if err := json.Unmarshal(hb, &hdr); err != nil {
		return Identity{}, tokenErr("header json")
	}
	if hdr.Alg != "RS256" && hdr.Alg != "ES256" {
		return Identity{}, tokenErr("alg %q not accepted", hdr.Alg)
	}
	if len(hdr.Crit) > 0 {
		return Identity{}, tokenErr("crit header not supported")
	}
	sig, err := b64.DecodeString(parts[2])
	if err != nil {
		return Identity{}, tokenErr("signature encoding")
	}
	key, err := v.key(ctx, hdr.Kid)
	if err != nil {
		return Identity{}, err
	}
	digest := sha256.Sum256([]byte(parts[0] + "." + parts[1]))
	switch k := key.(type) {
	case *rsa.PublicKey:
		if hdr.Alg != "RS256" {
			return Identity{}, tokenErr("alg %s does not match an RSA key", hdr.Alg)
		}
		if rsa.VerifyPKCS1v15(k, crypto.SHA256, digest[:], sig) != nil {
			return Identity{}, tokenErr("bad signature")
		}
	case *ecdsa.PublicKey:
		if hdr.Alg != "ES256" || len(sig) != 64 {
			return Identity{}, tokenErr("alg %s does not match a P-256 key", hdr.Alg)
		}
		r, s := new(big.Int).SetBytes(sig[:32]), new(big.Int).SetBytes(sig[32:])
		if !ecdsa.Verify(k, digest[:], r, s) {
			return Identity{}, tokenErr("bad signature")
		}
	default:
		return Identity{}, tokenErr("unsupported key type")
	}

	pb, err := b64.DecodeString(parts[1])
	if err != nil {
		return Identity{}, tokenErr("payload encoding")
	}
	var c struct {
		Iss      string          `json:"iss"`
		Sub      string          `json:"sub"`
		Aud      json.RawMessage `json:"aud"`
		Exp      *float64        `json:"exp"`
		Nbf      *float64        `json:"nbf"`
		Iat      *float64        `json:"iat"`
		Username string          `json:"preferred_username"`
		Groups   []string        `json:"groups"`
	}
	if err := json.Unmarshal(pb, &c); err != nil {
		return Identity{}, tokenErr("payload json")
	}
	now := v.cfg.Now()
	if c.Iss != v.cfg.Issuer {
		return Identity{}, tokenErr("issuer %q", c.Iss)
	}
	if !audienceHas(c.Aud, v.cfg.Audience) {
		return Identity{}, tokenErr("audience")
	}
	if c.Exp == nil {
		return Identity{}, tokenErr("no exp")
	}
	exp := time.Unix(int64(*c.Exp), 0)
	if !now.Before(exp.Add(v.cfg.Leeway)) {
		return Identity{}, tokenErr("expired")
	}
	if c.Nbf != nil && now.Add(v.cfg.Leeway).Before(time.Unix(int64(*c.Nbf), 0)) {
		return Identity{}, tokenErr("not yet valid")
	}
	if c.Iat != nil && now.Add(v.cfg.Leeway).Before(time.Unix(int64(*c.Iat), 0)) {
		return Identity{}, tokenErr("issued in the future")
	}
	return Identity{Subject: c.Sub, Username: c.Username, Groups: c.Groups, Expires: exp}, nil
}

func audienceHas(raw json.RawMessage, want string) bool {
	var one string
	if json.Unmarshal(raw, &one) == nil {
		return one == want
	}
	var many []string
	if json.Unmarshal(raw, &many) == nil {
		for _, a := range many {
			if a == want {
				return true
			}
		}
	}
	return false
}

// key returns the key for kid, refetching the JWKS when the set is stale or the kid is
// unknown (at most once per MinRefresh).
func (v *Verifier) key(ctx context.Context, kid string) (crypto.PublicKey, error) {
	v.mu.Lock()
	defer v.mu.Unlock()
	now := v.cfg.Now()
	k, ok := v.keys[kid]
	stale := now.Sub(v.fetched) > v.cfg.CacheTTL
	if ok && !stale {
		return k, nil
	}
	if now.Sub(v.tried) >= v.cfg.MinRefresh || v.fetched.IsZero() {
		v.tried = now
		keys, err := v.fetch(ctx)
		if err == nil {
			v.keys, v.fetched = keys, now
		} else if v.fetched.IsZero() {
			return nil, fmt.Errorf("jwks unavailable: %v", err)
		}
	}
	if k, ok := v.keys[kid]; ok {
		return k, nil
	}
	return nil, tokenErr("unknown kid %q", kid)
}

func (v *Verifier) fetch(ctx context.Context) (map[string]crypto.PublicKey, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, v.cfg.JWKSURL, nil)
	if err != nil {
		return nil, err
	}
	resp, err := v.cfg.HTTP.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("jwks: HTTP %d", resp.StatusCode)
	}
	return ParseJWKS(io.LimitReader(resp.Body, 1<<20))
}

// ParseJWKS reads RSA and P-256 signing keys; other keys are skipped.
func ParseJWKS(r io.Reader) (map[string]crypto.PublicKey, error) {
	var set struct {
		Keys []struct {
			Kty string `json:"kty"`
			Kid string `json:"kid"`
			Use string `json:"use"`
			N   string `json:"n"`
			E   string `json:"e"`
			Crv string `json:"crv"`
			X   string `json:"x"`
			Y   string `json:"y"`
		} `json:"keys"`
	}
	if err := json.NewDecoder(r).Decode(&set); err != nil {
		return nil, err
	}
	out := map[string]crypto.PublicKey{}
	for _, k := range set.Keys {
		if k.Use != "" && k.Use != "sig" {
			continue
		}
		switch k.Kty {
		case "RSA":
			n, err1 := b64.DecodeString(k.N)
			e, err2 := b64.DecodeString(k.E)
			if err1 != nil || err2 != nil || len(e) > 4 {
				continue
			}
			pk := &rsa.PublicKey{N: new(big.Int).SetBytes(n), E: int(new(big.Int).SetBytes(e).Int64())}
			if pk.N.BitLen() < 2048 {
				continue
			}
			out[k.Kid] = pk
		case "EC":
			if k.Crv != "P-256" {
				continue
			}
			x, err1 := b64.DecodeString(k.X)
			y, err2 := b64.DecodeString(k.Y)
			if err1 != nil || err2 != nil || len(x) != 32 || len(y) != 32 {
				continue
			}
			// Rejects points that are not on the curve.
			pk, err := ecdsa.ParseUncompressedPublicKey(elliptic.P256(), append(append([]byte{4}, x...), y...))
			if err != nil {
				continue
			}
			out[k.Kid] = pk
		}
	}
	if len(out) == 0 {
		return nil, errors.New("jwks: no usable signing keys")
	}
	return out, nil
}
