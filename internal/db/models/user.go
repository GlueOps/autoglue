package models

import (
	"time"

	"github.com/google/uuid"
)

type Role string

const (
	RoleAdmin Role = "admin"
	RoleUser  Role = "user"
)

type User struct {
	ID              uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Name            string    `gorm:"type:varchar(255);not null" json:"name"`
	Email           string    `gorm:"uniqueIndex" json:"email"`
	EmailVerified   bool      `gorm:"default:false" json:"email_verified"`
	EmailVerifiedAt time.Time `gorm:"default:null" json:"email_verified_at"`
	Password        string
	Role            Role
	Timestamped
}

func (r Role) IsValid() bool {
	switch r {
	case RoleAdmin, RoleUser:
		return true
	}
	return false
}
