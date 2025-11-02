package models

import (
	"time"

	"github.com/google/uuid"
)

type OrganizationKey struct {
	ID             uuid.UUID    `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	OrganizationID uuid.UUID    `gorm:"type:uuid;not null" json:"organization_id"`
	Organization   Organization `gorm:"foreignKey:OrganizationID;constraint:OnDelete:CASCADE" json:"organization"`
	MasterKeyID    uuid.UUID    `gorm:"type:uuid;not null"`
	MasterKey      MasterKey    `gorm:"foreignKey:MasterKeyID;constraint:OnDelete:CASCADE" json:"master_key"`
	EncryptedKey   string       `gorm:"not null"`
	IV             string       `gorm:"not null"`
	Tag            string       `gorm:"not null"`
	CreatedAt      time.Time    `gorm:"not null;default:now()" json:"created_at" format:"date-time"`
	UpdatedAt      time.Time    `gorm:"not null;default:now()" json:"updated_at" format:"date-time"`
}
