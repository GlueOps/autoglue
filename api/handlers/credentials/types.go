package credentials

import "github.com/google/uuid"

// swagger:model createCredentialRequest
type createCredentialRequest struct {
	// Provider name (e.g., "aws", "gitlab")
	// required: true
	Provider string `json:"provider"`
	// Secret material in plaintext; will be encrypted at rest
	// required: true
	Plaintext string `json:"plaintext"`
}

// swagger:model updateCredentialRequest
type updateCredentialRequest struct {
	Provider  *string `json:"provider,omitempty"`
	Plaintext *string `json:"plaintext,omitempty"`
}

// swagger:model credentialResponse
type credentialResponse struct {
	ID             uuid.UUID `json:"id"`
	OrganizationID uuid.UUID `json:"organization_id"`
	Provider       string    `json:"provider"`
	CreatedAt      string    `json:"created_at,omitempty"`
	UpdatedAt      string    `json:"updated_at,omitempty"`
}

// swagger:model credentialRevealResponse
type credentialRevealResponse struct {
	credentialResponse
	// Decrypted plaintext (only when ?reveal=true)
	Decrypted string `json:"decrypted"`
}

// ---- Crypto helpers (replace with your actual crypto) -----
type encryptedPayload struct {
	Data string
	IV   string
	Tag  string
}
