package models

import (
	"time"

	"github.com/google/uuid"
)

type UserEmail struct {
	// example: 3fa85f64-5717-4562-b3fc-2c963f66afa6
	ID         uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id" format:"uuid"`
	UserID     uuid.UUID `gorm:"index;not null" json:"user_id" format:"uuid"`
	User       User      `gorm:"foreignKey:UserID" json:"user"`
	Email      string    `gorm:"not null" json:"email"`
	IsVerified bool      `gorm:"not null;default:false" json:"is_verified"`
	IsPrimary  bool      `gorm:"not null;default:false" json:"is_primary"`
	CreatedAt  time.Time `gorm:"column:created_at;not null;default:now()" json:"created_at" format:"date-time"`
	UpdatedAt  time.Time `gorm:"autoUpdateTime;column:updated_at;not null;default:now()" json:"updated_at" format:"date-time"`
}
