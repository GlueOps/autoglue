package models

import "github.com/google/uuid"

type SshKey struct {
	ID             uuid.UUID    `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	OrganizationID uuid.UUID    `gorm:"type:uuid;not null" json:"organization_id"`
	Organization   Organization `gorm:"foreignKey:OrganizationID;constraint:OnDelete:CASCADE" json:"organization"`
	PublicKey      string       `gorm:"not null"`
	PrivateKey     string       `gorm:"not null"`
	Timestamped
}
