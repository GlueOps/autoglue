package dto

// swagger:model AuthStartResponse
type AuthStartResponse struct {
	AuthURL string `json:"auth_url" example:"https://accounts.google.com/o/oauth2/v2/auth?client_id=..."`
}

// swagger:model TokenPair
type TokenPair struct {
	AccessToken  string `json:"access_token" example:"eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCIsImtpZCI6Ij..."`
	RefreshToken string `json:"refresh_token" example:"m0l9o8rT3t0V8d3eFf...."`
	TokenType    string `json:"token_type" example:"Bearer"`
	ExpiresIn    int64  `json:"expires_in" example:"3600"`
}

// swagger:model RefreshRequest
type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" example:"m0l9o8rT3t0V8d3eFf..."`
}

// swagger:model LogoutRequest
type LogoutRequest struct {
	RefreshToken string `json:"refresh_token" example:"m0l9o8rT3t0V8d3eFf..."`
}
