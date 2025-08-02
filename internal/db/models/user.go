package models

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	ID            string `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Name          string `gorm:"type:varchar(255);not null" json:"name"`
	Email         string `gorm:"uniqueIndex" json:"email"`
	EmailVerified bool   `gorm:"default:false" json:"email_verified"`
	Role          string `gorm:"type:varchar(255);not null" json:"role"`
	Password      string
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}
