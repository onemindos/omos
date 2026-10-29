package callout

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"errors"
	"math/big"
	"strings"
	"testing"
	"time"
)

func newTestVerifier(t *testing.T, p *idp) *Verifier {
	t.Helper()
	v, err := NewVerifier(OIDCConfig{Issuer: testIss, Audience: testAud, JWKSURL: p.srv.URL})
	if err != nil {
		t.Fatal(err)
	}
	return v
}

func TestVerifierClaims(t *testing.T) {
	p := newIDP(t)
	v := newTestVerifier(t, p)
	ctx := context.Background()

	id, err := v.Verify(ctx, p.token(t, "olga", "onemind-operators"))
	if err != nil || id.Username != "olga" || id.Subject != subOf("olga") || len(id.Groups) != 1 {
		t.Fatalf("valid token: %+v %v", id, err)
	}
	arr := claims("olga", time.Minute, "onemind-operators")
	arr["aud"] = []string{"other", testAud}
	if _, err := v.Verify(ctx, p.sign(t, "RS256", p.rsaKid, arr)); err != nil {
		t.Fatalf("aud array: %v", err)
	}
	for name, mut := range map[string]func(map[string]any){
		"nbf in the future": func(c map[string]any) { c["nbf"] = time.Now().Add(5 * time.Minute).Unix() },
		"iat in the future": func(c map[string]any) { c["iat"] = time.Now().Add(5 * time.Minute).Unix() },
		"no exp":            func(c map[string]any) { delete(c, "exp") },
		"aud array without": func(c map[string]any) { c["aud"] = []string{"argocd"} },
		"exp past leeway":   func(c map[string]any) { c["exp"] = time.Now().Add(-31 * time.Second).Unix() },
	} {
		c := claims("olga", time.Minute, "onemind-operators")
		mut(c)
		if _, err := v.Verify(ctx, p.sign(t, "RS256", p.rsaKid, c)); !errors.Is(err, ErrToken) {
			t.Errorf("%s: %v", name, err)
		}
	}
	// An ES256 header on the RSA key (algorithm confusion) is refused.
	if _, err := v.Verify(ctx, p.sign(t, "ES256", p.rsaKid, claims("olga", time.Minute))); !errors.Is(err, ErrToken) {
		t.Errorf("alg/key mismatch: %v", err)
	}
	// crit headers are not understood, so refused.
	tok := p.token(t, "olga")
	parts := strings.Split(tok, ".")
	parts[0] = b64.EncodeToString([]byte(`{"alg":"RS256","kid":"rsa-1","crit":["exp"]}`))
	if _, err := v.Verify(ctx, strings.Join(parts, ".")); !errors.Is(err, ErrToken) {
		t.Errorf("crit: %v", err)
	}
}

func TestVerifierJWKSUnavailableAtStart(t *testing.T) {
	p := newIDP(t)
	p.setDown(true)
	v := newTestVerifier(t, p)
	if _, err := v.Verify(context.Background(), p.token(t, "olga", "onemind-operators")); err == nil {
		t.Fatal("verified with no key set")
	}
}

func TestParseJWKSSkipsWeakAndForeignKeys(t *testing.T) {
	small, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		t.Fatal(err)
	}
	doc := `{"keys":[
	  {"kty":"RSA","kid":"weak","n":"` + b64.EncodeToString(small.N.Bytes()) + `","e":"` + b64.EncodeToString(big.NewInt(65537).Bytes()) + `"},
	  {"kty":"RSA","kid":"enc","use":"enc","n":"x","e":"AQAB"},
	  {"kty":"EC","kid":"p384","crv":"P-384","x":"x","y":"y"},
	  {"kty":"EC","kid":"offcurve","crv":"P-256","x":"` + b64.EncodeToString(make([]byte, 32)) + `","y":"` + b64.EncodeToString(make([]byte, 32)) + `"},
	  {"kty":"oct","kid":"hmac","k":"c2VjcmV0"}]}`
	if _, err := ParseJWKS(strings.NewReader(doc)); err == nil {
		t.Fatal("accepted a key set with only weak or foreign keys")
	}
}
