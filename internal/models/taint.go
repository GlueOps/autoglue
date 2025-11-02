package models

import (
	"time"

	"github.com/google/uuid"
)

type Taint struct {
	ID             uuid.UUID    `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	OrganizationID uuid.UUID    `gorm:"type:uuid;not null" json:"organization_id"`
	Organization   Organization `gorm:"foreignKey:OrganizationID;constraint:OnDelete:CASCADE" json:"organization"`
	Key            string       `gorm:"not null" json:"key"`
	Value          string       `gorm:"not null" json:"value"`
	Effect         string       `gorm:"not null" json:"effect"`
	CreatedAt      time.Time    `gorm:"column:created_at;not null;default:now()" json:"created_at" format:"date-time"`
	UpdatedAt      time.Time    `gorm:"autoUpdateTime;column:updated_at;not null;default:now()" json:"updated_at" format:"date-time"`
}
