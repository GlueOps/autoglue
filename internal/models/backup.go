package models

import (
	"time"

	"github.com/google/uuid"
)

type Backup struct {
	ID             uuid.UUID    `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	OrganizationID uuid.UUID    `gorm:"type:uuid;not null;index;uniqueIndex:uniq_org_credential,priority:1"`
	Organization   Organization `gorm:"foreignKey:OrganizationID;constraint:OnDelete:CASCADE" json:"organization"`
	Enabled        bool         `gorm:"not null;default:false" json:"enabled"`
	CredentialID   uuid.UUID    `gorm:"type:uuid;not null;uniqueIndex:uniq_org_credential,priority:2" json:"credential_id"`
	Credential     Credential   `gorm:"foreignKey:CredentialID" json:"credential,omitempty"`
	CreatedAt      time.Time    `json:"created_at,omitempty" gorm:"type:timestamptz;column:created_at;not null;default:now()"`
	UpdatedAt      time.Time    `json:"updated_at,omitempty" gorm:"type:timestamptz;autoUpdateTime;column:updated_at;not null;default:now()"`
}
