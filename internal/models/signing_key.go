package models

import (
	"time"

	"github.com/google/uuid"
)

type SigningKey struct {
	ID          uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Kid         string     `gorm:"uniqueIndex;not null" json:"kid"`   // key id (header 'kid')
	Alg         string     `gorm:"not null" json:"alg"`               // RS256|RS384|RS512|EdDSA
	Use         string     `gorm:"not null;default:'sig'" json:"use"` // "sig"
	IsActive    bool       `gorm:"not null;default:true" json:"is_active"`
	PublicPEM   string     `gorm:"type:text;not null" json:"-"`
	PrivatePEM  string     `gorm:"type:text;not null" json:"-"`
	NotBefore   *time.Time `json:"-"`
	ExpiresAt   *time.Time `json:"-"`
	CreatedAt   time.Time  `gorm:"not null;default:now()" json:"created_at"`
	UpdatedAt   time.Time  `gorm:"not null;default:now()" json:"updated_at"`
	RotatedFrom *uuid.UUID `json:"-"` // previous key id, if any
}
