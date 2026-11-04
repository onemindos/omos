// SPDX-License-Identifier: Apache-2.0

package main

// Mesh network encryption as Common Commo does it (commoncommo
// core/impl/cryptoutil.cpp, MeshNetCrypto): two distinct 256-bit keys; the
// whole message (TAK protocol header included) is encrypted with AES-256-CBC
// under a random IV with PKCS#7 padding, and HMAC-SHA256 under the auth key
// covers IV || ciphertext. Wire layout: IV (16) | ciphertext | HMAC (32).
// Once keys are set, Commo drops every inbound mesh or directed message that
// does not decrypt; the gateway does the same.

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
)

const (
	meshKeyLen  = 32
	meshIVLen   = aes.BlockSize
	meshHMACLen = sha256.Size
)

// ErrDecrypt is returned for any message that fails authentication or padding.
// The reason is not distinguished, so the error is no padding oracle.
var ErrDecrypt = errors.New("mesh decrypt failed")

// MeshCrypto holds the two mesh keys.
type MeshCrypto struct {
	cipherKey [meshKeyLen]byte
	authKey   [meshKeyLen]byte
}

// NewMeshCrypto parses two hex-encoded 256-bit keys. Like Commo's
// setCryptoKeys, both must be given and they must differ.
func NewMeshCrypto(authHex, cipherHex string) (*MeshCrypto, error) {
	parse := func(name, s string) ([]byte, error) {
		b, err := hex.DecodeString(strings.TrimSpace(s))
		if err != nil {
			return nil, fmt.Errorf("%s: not hex: %w", name, err)
		}
		if len(b) != meshKeyLen {
			return nil, fmt.Errorf("%s: %d bytes, want %d (64 hex characters)", name, len(b), meshKeyLen)
		}
		return b, nil
	}
	a, err := parse("auth key", authHex)
	if err != nil {
		return nil, err
	}
	c, err := parse("cipher key", cipherHex)
	if err != nil {
		return nil, err
	}
	if bytes.Equal(a, c) {
		return nil, errors.New("auth and cipher keys must differ")
	}
	m := &MeshCrypto{}
	copy(m.authKey[:], a)
	copy(m.cipherKey[:], c)
	return m, nil
}

// Encrypt returns IV | AES-256-CBC(plain) | HMAC-SHA256(IV | ciphertext).
func (m *MeshCrypto) Encrypt(plain []byte) ([]byte, error) {
	iv := make([]byte, meshIVLen)
	if _, err := rand.Read(iv); err != nil {
		return nil, err
	}
	return m.encryptWithIV(plain, iv), nil
}

func (m *MeshCrypto) encryptWithIV(plain, iv []byte) []byte {
	block, _ := aes.NewCipher(m.cipherKey[:]) // key length is fixed, cannot fail
	pad := aes.BlockSize - len(plain)%aes.BlockSize
	padded := make([]byte, len(plain)+pad)
	copy(padded, plain)
	for i := len(plain); i < len(padded); i++ {
		padded[i] = byte(pad)
	}
	out := make([]byte, meshIVLen+len(padded), meshIVLen+len(padded)+meshHMACLen)
	copy(out, iv)
	cipher.NewCBCEncrypter(block, iv).CryptBlocks(out[meshIVLen:], padded)
	mac := hmac.New(sha256.New, m.authKey[:])
	mac.Write(out)
	return mac.Sum(out)
}

// Decrypt checks the HMAC first, then decrypts and removes the padding.
func (m *MeshCrypto) Decrypt(msg []byte) ([]byte, error) {
	if len(msg) < meshIVLen+aes.BlockSize+meshHMACLen {
		return nil, ErrDecrypt
	}
	body, tag := msg[:len(msg)-meshHMACLen], msg[len(msg)-meshHMACLen:]
	mac := hmac.New(sha256.New, m.authKey[:])
	mac.Write(body)
	if !hmac.Equal(mac.Sum(nil), tag) {
		return nil, ErrDecrypt
	}
	ct := body[meshIVLen:]
	if len(ct)%aes.BlockSize != 0 {
		return nil, ErrDecrypt
	}
	block, _ := aes.NewCipher(m.cipherKey[:])
	plain := make([]byte, len(ct))
	cipher.NewCBCDecrypter(block, body[:meshIVLen]).CryptBlocks(plain, ct)
	pad := int(plain[len(plain)-1])
	if pad == 0 || pad > aes.BlockSize {
		return nil, ErrDecrypt
	}
	for _, b := range plain[len(plain)-pad:] {
		if int(b) != pad {
			return nil, ErrDecrypt
		}
	}
	return plain[:len(plain)-pad], nil
}
