package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
)

var (
	ErrNoActiveMasterKey   = errors.New("no active master key found")
	ErrInvalidOrgID        = errors.New("invalid organization ID")
	ErrCredentialNotFound  = errors.New("credential not found")
	ErrInvalidMasterKeyLen = errors.New("invalid master key length")
)

func randomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return nil, fmt.Errorf("rand: %w", err)
	}
	return b, nil
}

func encryptAESGCM(plaintext, key []byte) (cipherNoTag, iv, tag []byte, _ error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("gcm: %w", err)
	}
	if gcm.NonceSize() != 12 {
		return nil, nil, nil, fmt.Errorf("unexpected nonce size: %d", gcm.NonceSize())
	}

	iv, err = randomBytes(gcm.NonceSize())
	if err != nil {
		return nil, nil, nil, err
	}

	// Go’s GCM returns ciphertext||tag, with 16-byte tag.
	cipherWithTag := gcm.Seal(nil, iv, plaintext, nil)
	if len(cipherWithTag) < 16 {
		return nil, nil, nil, errors.New("ciphertext too short")
	}
	tagLen := 16
	cipherNoTag = cipherWithTag[:len(cipherWithTag)-tagLen]
	tag = cipherWithTag[len(cipherWithTag)-tagLen:]
	return cipherNoTag, iv, tag, nil
}

func decryptAESGCM(cipherNoTag, key, iv, tag []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("cipher: %w", err)
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("gcm: %w", err)
	}
	if gcm.NonceSize() != len(iv) {
		return nil, fmt.Errorf("bad nonce size: %d", len(iv))
	}
	// Reattach tag
	cipherWithTag := append(append([]byte{}, cipherNoTag...), tag...)
	plain, err := gcm.Open(nil, iv, cipherWithTag, nil)
	if err != nil {
		return nil, fmt.Errorf("gcm open: %w", err)
	}
	return plain, nil
}

func encodeB64(b []byte) string {
	return base64.StdEncoding.EncodeToString(b)
}

func decodeB64(s string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(s)
}
