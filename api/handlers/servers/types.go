package servers

import "github.com/google/uuid"

// swagger:model createServerRequest
type createServerRequest struct {
	// Optional hostname
	// example: worker-01
	Hostname string `json:"hostname,omitempty"`
	// IPv4/IPv6 address
	// required: true
	// example: 10.0.1.23
	IPAddress string `json:"ip_address"`
	// SSH login user
	// required: true
	// example: ubuntu
	SSHUser string `json:"ssh_user"`
	// SSH key ID to use (must belong to the same org)
	// required: true
	// example: 2a1b9a6e-6fda-4e4b-8f80-0a3f8e0b8e4e
	SshKeyID string `json:"ssh_key_id"`
	// Role for this server (e.g., master, worker, bastion)
	// required: true
	// example: worker
	Role string `json:"role" example:"master|worker|bastion"`
	// Optional initial status (defaults to "pending")
	// enum: pending,provisioning,ready,failed
	// example: pending
	Status string `json:"status,omitempty" example:"pending|provisioning|ready|failed"`
}

// swagger:model updateServerRequest
type updateServerRequest struct {
	Hostname  *string `json:"hostname,omitempty"`
	IPAddress *string `json:"ip_address,omitempty"`
	SSHUser   *string `json:"ssh_user,omitempty"`
	SshKeyID  *string `json:"ssh_key_id,omitempty"`
	Role      *string `json:"role,omitempty" example:"master|worker|bastion"`
	// enum: pending,provisioning,ready,failed
	Status *string `json:"status,omitempty" example:"pending|provisioning|ready|failed"`
}

// swagger:model serverResponse
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
