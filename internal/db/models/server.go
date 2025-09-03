package models

import "github.com/google/uuid"

type Server struct {
	ID             uuid.UUID    `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	OrganizationID uuid.UUID    `gorm:"type:uuid;not null" json:"organization_id"`
	Organization   Organization `gorm:"foreignKey:OrganizationID;constraint:OnDelete:CASCADE" json:"organization"`
	Hostname       string       `json:"hostname"`
	IPAddress      string       `gorm:"not null"`
	SSHUser        string       `gorm:"not null"`
	SshKeyID       uuid.UUID    `gorm:"type:uuid;not null"`
	SshKey         SshKey       `gorm:"foreignKey:SshKeyID"`
	Role           string       `gorm:"not null"`          // e.g., "master", "worker", "bastion"
	Status         string       `gorm:"default:'pending'"` // pending, provisioning, ready, failed
	NodePools      []NodePool   `gorm:"many2many:node_servers;constraint:OnDelete:CASCADE" json:"servers,omitempty"`
	Timestamped
}
