package models

import (
	"time"

	"github.com/google/uuid"
)

type RefreshToken struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	UserID    uuid.UUID
	Token     string `gorm:"uniqueIndex"`
	ExpiresAt time.Time
	Revoked   bool
}
