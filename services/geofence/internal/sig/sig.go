// Package sig verifies policy-signed command envelopes (asset-command.md §6.3 `sig`: Ed25519
// over the RFC 8785 JCS form of the envelope without `sig`), with the same rule as the policy
// service's VerifyCommand (services/policy/internal/policy/envelope.go, PR #554).
package sig

import (
	"crypto/ed25519"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// Errors returned by Verify; the service maps them to refusal codes.
var (
	ErrUnsigned   = errors.New("unsigned")
	ErrUnknownKey = errors.New("unknown signing key")
	ErrBadSig     = errors.New("bad signature")
)

// Keys are the trusted policy public keys by key id (`sig.kid`).
type Keys map[string]ed25519.PublicKey

// Verify checks the envelope signature. A key id that is not in keys is refused, so a
// rotated-out key stops working as soon as it is removed from the mounted key set.
func (k Keys) Verify(raw []byte) (kid string, err error) {
	var env map[string]any
	if err := json.Unmarshal(raw, &env); err != nil {
		return "", fmt.Errorf("envelope is not JSON: %w", err)
	}
	s, ok := env["sig"].(map[string]any)
	if !ok {
		return "", ErrUnsigned
	}
	kid, _ = s["kid"].(string)
	if alg, _ := s["alg"].(string); alg != "ed25519" {
		return kid, fmt.Errorf("%w: alg %q", ErrBadSig, alg)
	}
	pub, ok := k[kid]
	if !ok {
		return kid, fmt.Errorf("%w %q", ErrUnknownKey, kid)
	}
	val, _ := s["value"].(string)
	b, err := base64.StdEncoding.DecodeString(val)
	if err != nil {
		return kid, fmt.Errorf("%w: %v", ErrBadSig, err)
	}
	delete(env, "sig")
	c, err := Canonicalize(env)
	if err != nil {
		return kid, err
	}
	if !ed25519.Verify(pub, c, b) {
		return kid, ErrBadSig
	}
	return kid, nil
}

// Sign adds `sig` to an envelope. The service never signs; tests and the evidence harness use
// it to stand in for the policy service.
func Sign(priv ed25519.PrivateKey, kid string, env map[string]any) error {
	delete(env, "sig")
	c, err := Canonicalize(env)
	if err != nil {
		return err
	}
	env["sig"] = map[string]any{"alg": "ed25519", "kid": kid,
		"value": base64.StdEncoding.EncodeToString(ed25519.Sign(priv, c))}
	return nil
}

// LoadDir reads every `{kid}.pem` (PKIX "PUBLIC KEY" PEM, Ed25519) in dir. The directory is a
// mounted ConfigMap of public keys; nothing secret is read.
func LoadDir(dir string) (Keys, error) {
	ents, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	keys := Keys{}
	for _, e := range ents {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".pem") {
			continue
		}
		b, err := os.ReadFile(filepath.Join(dir, e.Name()))
		if err != nil {
			return nil, err
		}
		pub, err := ParsePublicPEM(b)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", e.Name(), err)
		}
		keys[strings.TrimSuffix(e.Name(), ".pem")] = pub
	}
	if len(keys) == 0 {
		return nil, fmt.Errorf("no *.pem public keys in %s", dir)
	}
	return keys, nil
}

// ParsePublicPEM decodes a PKIX Ed25519 public key.
func ParsePublicPEM(b []byte) (ed25519.PublicKey, error) {
	blk, _ := pem.Decode(b)
	if blk == nil {
		return nil, errors.New("no PEM block")
	}
	k, err := x509.ParsePKIXPublicKey(blk.Bytes)
	if err != nil {
		return nil, err
	}
	pub, ok := k.(ed25519.PublicKey)
	if !ok {
		return nil, errors.New("not an Ed25519 public key")
	}
	return pub, nil
}

// IDs returns the key ids, sorted (for logs).
func (k Keys) IDs() []string {
	out := make([]string, 0, len(k))
	for id := range k {
		out = append(out, id)
	}
	sort.Strings(out)
	return out
}
