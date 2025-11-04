package models

import (
	"github.com/glueops/autoglue/internal/common"
)

type SshKey struct {
	common.AuditFields
	Organization        Organization `gorm:"foreignKey:OrganizationID;constraint:OnDelete:CASCADE" json:"organization"`
	Name                string       `gorm:"not null" json:"name"`
	PublicKey           string       `gorm:"not null"`
	EncryptedPrivateKey string       `gorm:"not null"`
	PrivateIV           string       `gorm:"not null"`
	PrivateTag          string       `gorm:"not null"`
	Fingerprint         string       `gorm:"not null;index" json:"fingerprint"`
}
