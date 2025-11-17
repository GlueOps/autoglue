package models

import (
	"time"

	"github.com/google/uuid"
)

type LoadBalancer struct {
	ID               uuid.UUID    `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	OrganizationID   uuid.UUID    `json:"organization_id" gorm:"type:uuid;index"`
	Organization     Organization `json:"organization" gorm:"foreignKey:OrganizationID;constraint:OnDelete:CASCADE"`
	Name             string       `json:"name" gorm:"not null"`
	Kind             string       `json:"kind" gorm:"not null"`
	PublicIPAddress  string       `json:"public_ip_address" gorm:"not null"`
	PrivateIPAddress string       `json:"private_ip_address" gorm:"not null"`
	CreatedAt        time.Time    `json:"created_at,omitempty" gorm:"type:timestamptz;column:created_at;not null;default:now()"`
	UpdatedAt        time.Time    `json:"updated_at,omitempty" gorm:"type:timestamptz;autoUpdateTime;column:updated_at;not null;default:now()"`
}
