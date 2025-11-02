package auth

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"time"

	"github.com/alexedwards/argon2id"
)

func SHA256Hex(s string) string {
	sum := sha256.Sum256([]byte(s))
	return hex.EncodeToString(sum[:])
}

var argonParams = &argon2id.Params{
	Memory:      64 * 1024, // 64MB
	Iterations:  3,
	Parallelism: 2,
	SaltLength:  16,
	KeyLength:   32,
}

func HashSecretArgon2id(plain string) (string, error) {
	return argon2id.CreateHash(plain, argonParams)
}

func VerifySecretArgon2id(encodedHash, plain string) (bool, error) {
	if encodedHash == "" {
		return false, errors.New("empty hash")
	}
	return argon2id.ComparePasswordAndHash(plain, encodedHash)
}

func NotExpired(expiresAt *time.Time) bool {
	return expiresAt == nil || time.Now().Before(*expiresAt)
}
