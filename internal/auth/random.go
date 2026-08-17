package auth

import (
	"crypto/rand"
	"encoding/base64"
)

// RandomB64URL returns n bytes of cryptographic randomness, base64url encoded
// without padding.
//
// This is the canonical home for a helper that internal/bg and
// internal/handlers each grew a private copy of; new callers use this one
// rather than adding a third. URL-safe and unpadded matters because these
// values travel in HTTP headers and JSON and end up pasted into shell on a
// bastion — '+', '/' and '=' survive none of those reliably.
func RandomB64URL(n int) (string, error) {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}
