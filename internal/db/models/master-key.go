package models

import "github.com/google/uuid"

type MasterKey struct {
	ID       uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Key      string    `gorm:"not null"`
	IsActive bool      `gorm:"default:true"`
	Timestamped
}
