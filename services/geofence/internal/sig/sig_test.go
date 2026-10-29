package sig

import (
	"crypto/ed25519"
	"crypto/rand"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func TestSignVerify(t *testing.T) {
	pub, priv, _ := ed25519.GenerateKey(rand.Reader)
	keys := Keys{"policy-1": pub}
	env := map[string]any{"v": 1, "id": "x", "cmd": "cmd.map.geofence.z.create", "data": map[string]any{"b": 2.5, "a": "é"}}
	if err := Sign(priv, "policy-1", env); err != nil {
		t.Fatal(err)
	}
	b, _ := json.Marshal(env)
	if kid, err := keys.Verify(b); err != nil || kid != "policy-1" {
		t.Fatalf("valid signature refused: %v", err)
	}
	// field order and whitespace do not matter (JCS)
	var m map[string]any
	_ = json.Unmarshal(b, &m)
	b2, _ := json.MarshalIndent(m, "", "  ")
	if _, err := keys.Verify(b2); err != nil {
		t.Fatalf("re-encoded envelope refused: %v", err)
	}
	m["data"].(map[string]any)["b"] = 3
	b3, _ := json.Marshal(m)
	if _, err := keys.Verify(b3); !errors.Is(err, ErrBadSig) {
		t.Fatalf("tampered envelope: %v", err)
	}
	delete(m, "sig")
	b4, _ := json.Marshal(m)
	if _, err := keys.Verify(b4); !errors.Is(err, ErrUnsigned) {
		t.Fatalf("unsigned: %v", err)
	}
	_, other, _ := ed25519.GenerateKey(rand.Reader)
	env2 := map[string]any{"v": 1}
	_ = Sign(other, "agent", env2)
	b5, _ := json.Marshal(env2)
	if _, err := keys.Verify(b5); !errors.Is(err, ErrUnknownKey) {
		t.Fatalf("unknown kid: %v", err)
	}
	_ = Sign(other, "policy-1", env2) // right kid, wrong key
	b6, _ := json.Marshal(env2)
	if _, err := keys.Verify(b6); !errors.Is(err, ErrBadSig) {
		t.Fatalf("wrong key: %v", err)
	}
}

func TestLoadDir(t *testing.T) {
	dir := t.TempDir()
	pub, _, _ := ed25519.GenerateKey(rand.Reader)
	der, _ := x509.MarshalPKIXPublicKey(pub)
	if err := os.WriteFile(filepath.Join(dir, "policy-2026-10.pem"), pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: der}), 0o600); err != nil {
		t.Fatal(err)
	}
	_ = os.WriteFile(filepath.Join(dir, "README"), []byte("ignored"), 0o600)
	keys, err := LoadDir(dir)
	if err != nil || len(keys) != 1 || !keys["policy-2026-10"].Equal(pub) {
		t.Fatalf("LoadDir: %v %v", keys.IDs(), err)
	}
	if _, err := LoadDir(t.TempDir()); err == nil {
		t.Fatal("empty dir accepted")
	}
}
