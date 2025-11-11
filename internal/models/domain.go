package models

import (
	"time"

	"github.com/google/uuid"
)

type Domain struct {
	ID             uuid.UUID    `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	OrganizationID uuid.UUID    `gorm:"type:uuid;not null;uniqueIndex:idx_credentials_org_provider" json:"organization_id"`
	Organization   Organization `gorm:"foreignKey:OrganizationID;constraint:OnDelete:CASCADE" json:"organization"`
	ClusterID      *uuid.UUID   `gorm:"type:uuid" json:"cluster_id,omitempty"`
	Cluster        *Cluster     `gorm:"foreignKey:ClusterID" json:"cluster,omitempty"`
	DomainName     string       `gorm:"not null;index" json:"domain_name,omitempty"`
	DomainID       string
	CredentialID   uuid.UUID  `gorm:"type:uuid;not null" json:"credential_id"`
	Credential     Credential `gorm:"foreignKey:CredentialID" json:"credential,omitempty"`
	CreatedAt      time.Time  `json:"created_at,omitempty" gorm:"type:timestamptz;column:created_at;not null;default:now()"`
	UpdatedAt      time.Time  `json:"updated_at,omitempty" gorm:"type:timestamptz;autoUpdateTime;column:updated_at;not null;default:now()"`
}
