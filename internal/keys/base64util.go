package keys

import (
	"encoding/base64"
	"errors"
	"strings"
)

func decode32ByteKey(s string) ([]byte, error) {
	try := func(enc *base64.Encoding, v string) ([]byte, bool) {
		if b, err := enc.DecodeString(v); err == nil && len(b) == 32 {
			return b, true
		}
		return nil, false
	}

	// Try raw (no padding) variants first
	if b, ok := try(base64.RawURLEncoding, s); ok {
		return b, nil
	}
	if b, ok := try(base64.RawStdEncoding, s); ok {
		return b, nil
	}

	// Try padded variants (add padding if missing)
	pad := func(v string) string { return v + strings.Repeat("=", (4-len(v)%4)%4) }
	if b, ok := try(base64.URLEncoding, pad(s)); ok {
		return b, nil
	}
	if b, ok := try(base64.StdEncoding, pad(s)); ok {
		return b, nil
	}

	return nil, errors.New("key must be 32 bytes in base64/base64url")
}
