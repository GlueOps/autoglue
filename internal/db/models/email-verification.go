package models

import (
	"time"

	"github.com/google/uuid"
)

type EmailVerification struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	UserID    uuid.UUID `gorm:"type:uuid;not null;index" json:"user_id"`
	User      User      `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"-"`
	Token     string    `gorm:"type:char(43);uniqueIndex;not null" json:"-"`
	ExpiresAt time.Time `gorm:"not null;index" json:"expires_at"`
	Used      bool      `gorm:"not null;default:false;index" json:"used"`
	Timestamped
}

func (e EmailVerification) IsActive(now time.Time) bool {
	return !e.Used && now.Before(e.ExpiresAt)
}
