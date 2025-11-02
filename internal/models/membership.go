package models

import (
	"time"

	"github.com/google/uuid"
)

type Membership struct {
	ID             uuid.UUID    `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id" format:"uuid"`
	UserID         uuid.UUID    `gorm:"index;not null" json:"user_id" format:"uuid"`
	User           User         `gorm:"foreignKey:UserID" json:"-"`
	OrganizationID uuid.UUID    `gorm:"index;not null" json:"org_id" format:"uuid"`
	Organization   Organization `gorm:"foreignKey:OrganizationID;constraint:OnDelete:CASCADE" json:"-"`
	Role           string       `gorm:"not null;default:'member'" json:"role"`
	CreatedAt      time.Time    `gorm:"not null;default:now()" json:"created_at" format:"date-time"`
	UpdatedAt      time.Time    `gorm:"not null;default:now()" json:"updated_at" format:"date-time"`
}
