package policy

import (
	"crypto/ed25519"
	"crypto/rand"
	"crypto/x509"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"os"
	"regexp"
	"time"
)

var uuidRE = regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)

// NewUUID returns a lowercase UUID v4.
func NewUUID() string {
	var b [16]byte
	_, _ = rand.Read(b[:])
	b[6] = (b[6] & 0x0f) | 0x40
	b[8] = (b[8] & 0x3f) | 0x80
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:16])
}

const crockford = "0123456789abcdefghjkmnpqrstvwxyz"

// NewULID returns a lowercase ULID (48-bit ms time + 80 random bits, Crockford base32).
func NewULID(t time.Time) string {
	var b [16]byte
	ms := uint64(t.UnixMilli())
	binary.BigEndian.PutUint16(b[0:2], uint16(ms>>32))
	binary.BigEndian.PutUint32(b[2:6], uint32(ms))
	_, _ = rand.Read(b[6:])
	// 128 bits -> 26 chars of 5 bits (first char carries 3 bits).
	out := make([]byte, 26)
	var acc uint64
	var bits uint
	idx := 25
	for i := 15; i >= 0; i-- {
		acc |= uint64(b[i]) << bits
		bits += 8
		for bits >= 5 && idx >= 0 {
			out[idx] = crockford[acc&31]
			acc >>= 5
			bits -= 5
			idx--
		}
	}
	if idx == 0 {
		out[0] = crockford[acc&31]
	}
	return string(out)
}

func ts(t time.Time) string { return t.UTC().Format("2006-01-02T15:04:05.000Z") }

// Signer signs command envelopes with Ed25519 (asset-command.md section 6.3 `sig`).
type Signer struct {
	Key ed25519.PrivateKey
	KID string
}

// LoadSigner reads a PKCS#8 PEM Ed25519 private key (mounted from a SealedSecret).
func LoadSigner(path, kid string) (*Signer, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	blk, _ := pem.Decode(b)
	if blk == nil {
		return nil, errors.New("signing key: no PEM block")
	}
	k, err := x509.ParsePKCS8PrivateKey(blk.Bytes)
	if err != nil {
		return nil, err
	}
	ek, ok := k.(ed25519.PrivateKey)
	if !ok {
		return nil, errors.New("signing key is not Ed25519")
	}
	return &Signer{Key: ek, KID: kid}, nil
}

// Sign adds `sig` over the JCS form of the envelope without `sig`.
func (s *Signer) Sign(env map[string]any) error {
	delete(env, "sig")
	c, err := Canonicalize(env)
	if err != nil {
		return err
	}
	env["sig"] = map[string]any{"alg": "ed25519", "kid": s.KID,
		"value": base64.StdEncoding.EncodeToString(ed25519.Sign(s.Key, c))}
	return nil
}

// VerifyCommand checks a signed command envelope; bridges use the same rule.
func VerifyCommand(pub ed25519.PublicKey, raw []byte) error {
	var env map[string]any
	if err := json.Unmarshal(raw, &env); err != nil {
		return err
	}
	sig, ok := env["sig"].(map[string]any)
	if !ok {
		return errors.New("unsigned")
	}
	val, _ := sig["value"].(string)
	b, err := base64.StdEncoding.DecodeString(val)
	if err != nil {
		return err
	}
	delete(env, "sig")
	c, err := Canonicalize(env)
	if err != nil {
		return err
	}
	if !ed25519.Verify(pub, c, b) {
		return errors.New("bad signature")
	}
	return nil
}
