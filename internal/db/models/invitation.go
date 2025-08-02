package models

import (
	"github.com/google/uuid"
	"time"
)

type Invitation struct {
	ID             uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	OrganizationID uuid.UUID `gorm:"type:uuid;not null"`
	Email          string    `gorm:"type:text;not null"`
	Role           string    `gorm:"type:text;default:'member';not null"`
	Status         string    `gorm:"type:text;default:'pending';not null"` // pending, accepted, revoked
	ExpiresAt      time.Time `gorm:"not null"`
	InviterID      uuid.UUID `gorm:"type:uuid;not null"`
	Timestamped
}
