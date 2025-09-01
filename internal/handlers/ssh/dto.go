package ssh

import "github.com/google/uuid"

type createSSHRequest struct {
	Bits     *int   `json:"bits,omitempty" example:"4096"`
	Comment  string `json:"comment,omitempty" example:"deploy@autoglue"`
	Download string `json:"download,omitempty" example:"both"`
	Name     string `json:"name"`
}

type sshResponse struct {
	ID             uuid.UUID `json:"id"`
	OrganizationID uuid.UUID `json:"organization_id"`
	Name           string    `json:"name"`
	PublicKey      string    `json:"public_keys"`
	Fingerprint    string    `json:"fingerprint"`
	CreatedAt      string    `json:"created_at,omitempty"`
	UpdatedAt      string    `json:"updated_at,omitempty"`
}

type sshRevealResponse struct {
	sshResponse
	PrivateKey string `json:"private_key"`
}
