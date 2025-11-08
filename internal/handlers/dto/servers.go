package dto

import "github.com/google/uuid"

type CreateServerRequest struct {
	Hostname         string `json:"hostname,omitempty"`
	PublicIPAddress  string `json:"public_ip_address,omitempty"`
	PrivateIPAddress string `json:"private_ip_address"`
	SSHUser          string `json:"ssh_user"`
	SshKeyID         string `json:"ssh_key_id"`
	Role             string `json:"role" example:"master|worker|bastion" enums:"master,worker,bastion"`
	Status           string `json:"status,omitempty" example:"pending|provisioning|ready|failed" enums:"pending,provisioning,ready,failed"`
}

type UpdateServerRequest struct {
	Hostname         *string `json:"hostname,omitempty"`
	PublicIPAddress  *string `json:"public_ip_address,omitempty"`
	PrivateIPAddress *string `json:"private_ip_address,omitempty"`
	SSHUser          *string `json:"ssh_user,omitempty"`
	SshKeyID         *string `json:"ssh_key_id,omitempty"`
	Role             *string `json:"role" example:"master|worker|bastion" enums:"master,worker,bastion"`
	Status           *string `json:"status,omitempty" example:"pending|provisioning|ready|failed" enums:"pending,provisioning,ready,failed"`
}

type ServerResponse struct {
	ID               uuid.UUID `json:"id"`
	OrganizationID   uuid.UUID `json:"organization_id"`
	Hostname         string    `json:"hostname"`
	PublicIPAddress  *string   `json:"public_ip_address,omitempty"`
	PrivateIPAddress string    `json:"private_ip_address"`
	SSHUser          string    `json:"ssh_user"`
	SshKeyID         uuid.UUID `json:"ssh_key_id"`
	Role             string    `json:"role" example:"master|worker|bastion" enums:"master,worker,bastion"`
	Status           string    `json:"status,omitempty" example:"pending|provisioning|ready|failed" enums:"pending,provisioning,ready,failed"`
	CreatedAt        string    `json:"created_at,omitempty"`
	UpdatedAt        string    `json:"updated_at,omitempty"`
}
