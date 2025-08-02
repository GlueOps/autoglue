package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Member struct {
	ID             uuid.UUID    `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	UserID         uuid.UUID    `gorm:"type:uuid;not null" json:"user_id"`
	User           User         `gorm:"foreignKey:UserID" json:"user"`
	OrganizationID uuid.UUID    `gorm:"type:uuid;not null" json:"organization_id"`
	Organization   Organization `gorm:"foreignKey:OrganizationID" json:"organization"`
	Role           string       `gorm:"not null;default:member" json:"role"` // e.g. admin, member
	Timestamped
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
