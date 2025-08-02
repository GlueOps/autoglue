package models

import "time"

type Timestamped struct {
	CreatedAt time.Time `gorm:"column:created_at;not null;default:now()" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime;column:updated_at;not null;default:now()" json:"updated_at"`
}
