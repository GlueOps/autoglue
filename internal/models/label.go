package models

import (
	"time"

	"github.com/google/uuid"
)

type Label struct {
	ID             uuid.UUID    `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	OrganizationID uuid.UUID    `gorm:"type:uuid;not null" json:"organization_id"`
	Organization   Organization `gorm:"foreignKey:OrganizationID;constraint:OnDelete:CASCADE" json:"organization"`
	Key            string       `gorm:"not null" json:"key"`
	Value          string       `gorm:"not null" json:"value"`
	NodePools      []NodePool   `gorm:"many2many:node_labels;constraint:OnDelete:CASCADE" json:"servers,omitempty"`
	CreatedAt      time.Time    `gorm:"column:created_at;not null;default:now()" json:"created_at"`
	UpdatedAt      time.Time    `gorm:"autoUpdateTime;column:updated_at;not null;default:now()" json:"updated_at"`
}
