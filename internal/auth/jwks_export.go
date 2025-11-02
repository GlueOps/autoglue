package auth

import (
	"crypto/ed25519"
	"crypto/rsa"
	"encoding/base64"
	"fmt"
	"math/big"
)

// base64url (no padding)
func b64url(b []byte) string {
	return base64.RawURLEncoding.EncodeToString(b)
}

// convert small int (RSA exponent) to big-endian bytes
func fromInt(i int) []byte {
	var x big.Int
	x.SetInt64(int64(i))
	return x.Bytes()
}

// --- public accessors for JWKS ---

// KeyMeta is a minimal metadata view exposed for JWKS rendering.
type KeyMeta struct {
	Alg string
}

// MetaFor returns minimal metadata (currently the alg) for a given kid.
// If not found, returns zero value (Alg == "").
func MetaFor(kid string) KeyMeta {
	kc.mu.RLock()
	defer kc.mu.RUnlock()
	if m, ok := kc.meta[kid]; ok {
		return KeyMeta{Alg: m.Alg}
	}
	return KeyMeta{}
}

// KcCopy invokes fn with a shallow copy of the public key map (kid -> public key instance).
// Useful to iterate without holding the lock during JSON building.
func KcCopy(fn func(map[string]interface{})) {
	kc.mu.RLock()
	defer kc.mu.RUnlock()
	out := make(map[string]interface{}, len(kc.pub))
	for kid, pk := range kc.pub {
		out[kid] = pk
	}
	fmt.Println(out)
	fn(out)
}

// PubToJWK converts a parsed public key into bare JWK parameters + kty.
// - RSA: returns n/e (base64url) and kty="RSA"
// - Ed25519: returns x (base64url) and kty="OKP"
func PubToJWK(_kid, _alg string, pub any) (map[string]string, string) {
	switch k := pub.(type) {
	case *rsa.PublicKey:
		return map[string]string{
			"n": b64url(k.N.Bytes()),
			"e": b64url(fromInt(k.E)),
		}, "RSA"
	case ed25519.PublicKey:
		return map[string]string{
			"x": b64url([]byte(k)),
		}, "OKP"
	default:
		return nil, ""
	}
}
