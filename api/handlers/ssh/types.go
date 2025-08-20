package ssh

import "github.com/google/uuid"

// swagger:model createSSHRequest
type createSSHRequest struct {
	// RSA key size in bits. Allowed: 2048, 3072, 4096. Default: 4096
	// example: 4096
	Bits *int `json:"bits,omitempty"`
	// Optional comment appended to the authorized_keys string
	// example: deploy@autoglue
	Comment string `json:"comment,omitempty"`
	// Optional immediate download: "none" (default), "public", "private", "both"
	// example: none
	Download string `json:"download,omitempty"`
}

// swagger:model sshResponse
type sshResponse struct {
	ID             uuid.UUID `json:"id"`
	OrganizationID uuid.UUID `json:"organization_id"`
	PublicKey      string    `json:"public_keys"`
	CreatedAt      string    `json:"created_at,omitempty"`
	UpdatedAt      string    `json:"updated_at,omitempty"`
}

// swagger:model sshRevealResponse
type sshRevealResponse struct {
	sshResponse
	// Private key in PEM format (revealed only when requested)
	PrivateKey string `json:"private_key"`
}
