package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type IssueOpts struct {
	Subject  string
	Issuer   string
	Audience string
	TTL      time.Duration
	Claims   map[string]any // extra app claims
}

func IssueAccessToken(opts IssueOpts) (string, error) {
	kc.mu.RLock()
	defer kc.mu.RUnlock()

	if kc.selPriv == nil || kc.selKid == "" || kc.selAlg == "" {
		return "", errors.New("no active signing key")
	}

	claims := jwt.MapClaims{
		"iss": opts.Issuer,
		"aud": opts.Audience,
		"sub": opts.Subject,
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(opts.TTL).Unix(),
	}
	for k, v := range opts.Claims {
		claims[k] = v
	}

	var method jwt.SigningMethod
	switch kc.selAlg {
	case "RS256":
		method = jwt.SigningMethodRS256
	case "RS384":
		method = jwt.SigningMethodRS384
	case "RS512":
		method = jwt.SigningMethodRS512
	case "EdDSA":
		method = jwt.SigningMethodEdDSA
	default:
		return "", errors.New("unsupported alg")
	}

	token := jwt.NewWithClaims(method, claims)
	token.Header["kid"] = kc.selKid

	return token.SignedString(kc.selPriv)
}
