package handlers

import (
	"net/http"

	"github.com/glueops/autoglue/internal/auth"
	"github.com/glueops/autoglue/internal/handlers/dto"
	"github.com/glueops/autoglue/internal/utils"
)

type jwk struct {
	Kty string `json:"kty"`
	Use string `json:"use,omitempty"`
	Kid string `json:"kid,omitempty"`
	Alg string `json:"alg,omitempty"`
	N   string `json:"n,omitempty"` // RSA modulus (base64url)
	E   string `json:"e,omitempty"` // RSA exponent (base64url)
	X   string `json:"x,omitempty"` // Ed25519 public key (base64url)
}

type jwks struct {
	Keys []jwk `json:"keys"`
}

// JWKSHandler godoc
//
//	@ID				getJWKS
//	@Summary		Get JWKS
//	@Description	Returns the JSON Web Key Set for token verification
//	@Tags			Auth
//	@Produce		json
//	@Success		200	{object}	dto.JWKS
//	@Router			/.well-known/jwks.json [get]
func JWKSHandler(w http.ResponseWriter, _ *http.Request) {
	out := dto.JWKS{Keys: make([]dto.JWK, 0)}

	auth.KcCopy(func(pub map[string]interface{}) {
		for kid, pk := range pub {
			meta := auth.MetaFor(kid)
			params, kty := auth.PubToJWK(kid, meta.Alg, pk)
			if kty == "" {
				continue
			}
			j := dto.JWK{
				Kty: kty,
				Use: "sig",
				Kid: kid,
				Alg: meta.Alg,
				N:   params["n"],
				E:   params["e"],
				X:   params["x"],
			}
			out.Keys = append(out.Keys, j)
		}
	})
	utils.WriteJSON(w, http.StatusOK, out)
}
