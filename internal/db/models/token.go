package models

import "time"

type RefreshToken struct {
	ID        string `gorm:"primaryKey"`
	UserID    string
	Token     string `gorm:"uniqueIndex"`
	ExpiresAt time.Time
	Revoked   bool
}
