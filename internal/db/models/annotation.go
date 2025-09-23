package models

import "github.com/google/uuid"

type Annotation struct {
	ID             uuid.UUID    `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	OrganizationID uuid.UUID    `gorm:"type:uuid;not null" json:"organization_id"`
	Organization   Organization `gorm:"foreignKey:OrganizationID;constraint:OnDelete:CASCADE" json:"organization"`
	Key            string       `gorm:"not null" json:"key"`
	Value          string       `gorm:"not null" json:"value"`
	NodePools      []NodePool   `gorm:"many2many:node_annotations;constraint:OnDelete:CASCADE" json:"servers,omitempty"`
	Timestamped
}
