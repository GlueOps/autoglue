package models

import (
	"time"

	"github.com/google/uuid"
)

type APIKey struct {
	ID         uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id" format:"uuid"`
	Name       string     `gorm:"not null;default:''" json:"name"`
	KeyHash    string     `gorm:"uniqueIndex;not null" json:"-"`
	Scope      string     `gorm:"not null;default:''" json:"scope"`
	UserID     *uuid.UUID `json:"user_id,omitempty" format:"uuid"`
	OrgID      *uuid.UUID `json:"org_id,omitempty" format:"uuid"`
	SecretHash *string    `json:"-"`
	ExpiresAt  *time.Time `json:"expires_at,omitempty" format:"date-time"`
	Revoked    bool       `gorm:"not null;default:false" json:"revoked"`
	Prefix     *string    `json:"prefix,omitempty"`
	LastUsedAt *time.Time `json:"last_used_at,omitempty" format:"date-time"`
	CreatedAt  time.Time  `gorm:"not null;default:now()" json:"created_at" format:"date-time"`
	UpdatedAt  time.Time  `gorm:"not null;default:now()" json:"updated_at" format:"date-time"`
}
