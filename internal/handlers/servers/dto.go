package servers

import "github.com/google/uuid"

type createServerRequest struct {
	Hostname  string `json:"hostname,omitempty"`
	IPAddress string `json:"ip_address"`
	SSHUser   string `json:"ssh_user"`
	SshKeyID  string `json:"ssh_key_id"`
	Role      string `json:"role" example:"master|worker|bastion"`
	Status    string `json:"status,omitempty" example:"pending|provisioning|ready|failed"`
}

type updateServerRequest struct {
	Hostname  *string `json:"hostname,omitempty"`
	IPAddress *string `json:"ip_address,omitempty"`
	SSHUser   *string `json:"ssh_user,omitempty"`
	SshKeyID  *string `json:"ssh_key_id,omitempty"`
	Role      *string `json:"role,omitempty" example:"master|worker|bastion"`
	// enum: pending,provisioning,ready,failed
	Status *string `json:"status,omitempty" example:"pending|provisioning|ready|failed"`
}

type serverResponse struct {
	ID             uuid.UUID `json:"id"`
	OrganizationID uuid.UUID `json:"organization_id"`
	Hostname       string    `json:"hostname"`
	IPAddress      string    `json:"ip_address"`
	SSHUser        string    `json:"ssh_user"`
	SshKeyID       uuid.UUID `json:"ssh_key_id"`
	Role           string    `json:"role"`
	Status         string    `json:"status"`
	CreatedAt      string    `json:"created_at,omitempty"`
	UpdatedAt      string    `json:"updated_at,omitempty"`
}
