package models

import "github.com/google/uuid"

type NodeGroup struct {
	ID             uuid.UUID    `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	OrganizationID uuid.UUID    `gorm:"type:uuid;not null" json:"organization_id"`
	Organization   Organization `gorm:"foreignKey:OrganizationID;constraint:OnDelete:CASCADE" json:"organization"`
	Name           string       `gorm:"not null" json:"name"`
	Servers        []Server     `gorm:"many2many:cluster_servers;constraint:OnDelete:CASCADE" json:"servers,omitempty"`
	Labels         []NodeLabel  `gorm:"many2many:node_group_labels;constraint:OnDelete:CASCADE" json:"labels,omitempty"`
	Taints         []NodeTaint  `gorm:"many2many:node_group_taints;constraint:OnDelete:CASCADE" json:"taints,omitempty"`
	Timestamped
}
