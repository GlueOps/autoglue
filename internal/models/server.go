package models

import (
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Server struct {
	ID               uuid.UUID    `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	OrganizationID   uuid.UUID    `gorm:"type:uuid;not null" json:"organization_id"`
	Organization     Organization `gorm:"foreignKey:OrganizationID;constraint:OnDelete:CASCADE" json:"organization"`
	Hostname         string       `json:"hostname"`
	PublicIPAddress  *string      `json:"public_ip_address,omitempty"`
	PrivateIPAddress string       `gorm:"not null" json:"private_ip_address"`
	SSHUser          string       `gorm:"not null" json:"ssh_user"`
	SshKeyID         uuid.UUID    `gorm:"type:uuid;not null" json:"ssh_key_id"`
	SshKey           SshKey       `gorm:"foreignKey:SshKeyID" json:"ssh_key"`
	Role             string       `gorm:"not null" json:"role" enums:"master,worker,bastion"`                           // e.g., "master", "worker", "bastion"
	Status           string       `gorm:"default:'pending'" json:"status" enums:"pending, provisioning, ready, failed"` // pending, provisioning, ready, failed
	NodePools        []NodePool   `gorm:"many2many:node_servers;constraint:OnDelete:CASCADE" json:"node_pools,omitempty"`
	SSHHostKey       string       `gorm:"column:ssh_host_key"`
	SSHHostKeyAlgo   string       `gorm:"column:ssh_host_key_algo"`
	CreatedAt        time.Time    `gorm:"not null;default:now()" json:"created_at" format:"date-time"`
	UpdatedAt        time.Time    `gorm:"not null;default:now()" json:"updated_at" format:"date-time"`
}

func (s *Server) BeforeSave(tx *gorm.DB) error {
	role := strings.ToLower(strings.TrimSpace(s.Role))
	if role == "bastion" {
		if s.PublicIPAddress == nil || strings.TrimSpace(*s.PublicIPAddress) == "" {
			return errors.New("public_ip_address is required for role=bastion")
		}
	}
	return nil
}
