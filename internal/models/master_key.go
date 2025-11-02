package models

import (
	"time"

	"github.com/google/uuid"
)

type MasterKey struct {
	ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Key       string    `gorm:"not null"`
	IsActive  bool      `gorm:"default:true"`
	CreatedAt time.Time `gorm:"column:created_at;not null;default:now()" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime;column:updated_at;not null;default:now()" json:"updated_at"`
}
