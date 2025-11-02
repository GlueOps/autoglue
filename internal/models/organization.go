package models

import (
	"time"

	"github.com/google/uuid"
)

type Organization struct {
	// example: 3fa85f64-5717-4562-b3fc-2c963f66afa6
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id" format:"uuid"`
	Name      string    `gorm:"not null" json:"name"`
	Domain    *string   `gorm:"index" json:"domain"`
	CreatedAt time.Time `gorm:"column:created_at;not null;default:now()" json:"created_at" format:"date-time"`
	UpdatedAt time.Time `gorm:"autoUpdateTime;column:updated_at;not null;default:now()" json:"updated_at" format:"date-time"`
}
