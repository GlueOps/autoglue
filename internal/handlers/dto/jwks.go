package dto

// JWK represents a single JSON Web Key (public only).
// swagger:model JWK
type JWK struct {
	Kty string `json:"kty" example:"RSA" gorm:"-"`
	Use string `json:"use,omitempty" example:"sig" gorm:"-"`
	Kid string `json:"kid,omitempty" example:"7c6f1d0a-7a98-4e6a-9dbf-6b1af4b9f345" gorm:"-"`
	Alg string `json:"alg,omitempty" example:"RS256" gorm:"-"`
	N   string `json:"n,omitempty" gorm:"-"`
	E   string `json:"e,omitempty" example:"AQAB" gorm:"-"`
	X   string `json:"x,omitempty" gorm:"-"`
}

// JWKS is a JSON Web Key Set container.
// swagger:model JWKS
type JWKS struct {
	Keys []JWK `json:"keys" gorm:"-"`
}
