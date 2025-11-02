package models

import (
	"time"

	"github.com/google/uuid"
)

type SshKey struct {
	ID                  uuid.UUID    `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	OrganizationID      uuid.UUID    `gorm:"type:uuid;not null" json:"organization_id"`
	Organization        Organization `gorm:"foreignKey:OrganizationID;constraint:OnDelete:CASCADE" json:"organization"`
	Name                string       `gorm:"not null" json:"name"`
	PublicKey           string       `gorm:"not null"`
	EncryptedPrivateKey string       `gorm:"not null"`
	PrivateIV           string       `gorm:"not null"`
	PrivateTag          string       `gorm:"not null"`
	Fingerprint         string       `gorm:"not null;index" json:"fingerprint"`
	CreatedAt           time.Time    `gorm:"not null;default:now()" json:"created_at" format:"date-time"`
	UpdatedAt           time.Time    `gorm:"not null;default:now()" json:"updated_at" format:"date-time"`
}
