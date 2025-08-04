package models

import "github.com/google/uuid"

type OrganizationKey struct {
	ID             uuid.UUID    `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	OrganizationID uuid.UUID    `gorm:"type:uuid;not null" json:"organization_id"`
	Organization   Organization `gorm:"foreignKey:OrganizationID;constraint:OnDelete:CASCADE" json:"organization"`
	EncryptedKey   string       `gorm:"not null"`
	IV             string       `gorm:"not null"`
	Tag            string       `gorm:"not null"`
	Timestamped
}

type Credential struct {
	ID             uuid.UUID    `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	OrganizationID uuid.UUID    `gorm:"type:uuid;not null" json:"organization_id"`
	Organization   Organization `gorm:"foreignKey:OrganizationID;constraint:OnDelete:CASCADE" json:"organization"`
	Provider       string       `gorm:"type:varchar(50);not null"`
	EncryptedData  string       `gorm:"not null"`
	IV             string       `gorm:"not null"`
	Tag            string       `gorm:"not null"`
	Timestamped
}

type SshKey struct {
	ID             uuid.UUID    `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	OrganizationID uuid.UUID    `gorm:"type:uuid;not null" json:"organization_id"`
	Organization   Organization `gorm:"foreignKey:OrganizationID;constraint:OnDelete:CASCADE" json:"organization"`
	PublicKey      string       `gorm:"not null"`
	PrivateKey     string       `gorm:"not null"`
	Timestamped
}
