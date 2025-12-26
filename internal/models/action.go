package models

import (
	"time"

	"github.com/google/uuid"
)

type Action struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id" format:"uuid"`
	Label       string    `gorm:"type:varchar(255);not null;uniqueIndex" json:"label"`
	Description string    `gorm:"type:text;not null" json:"description"`
	MakeTarget  string    `gorm:"type:varchar(255);not null;uniqueIndex" json:"make_target"`
	CreatedAt   time.Time `json:"created_at,omitempty" gorm:"type:timestamptz;column:created_at;not null;default:now()" format:"date-time"`
	UpdatedAt   time.Time `json:"updated_at,omitempty" gorm:"type:timestamptz;autoUpdateTime;column:updated_at;not null;default:now()" format:"date-time"`
}
