package models

import (
	"github.com/google/uuid"
)

type Cluster struct {
	ID                  uuid.UUID    `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	OrganizationID      uuid.UUID    `gorm:"type:uuid;not null" json:"organization_id"`
	Organization        Organization `gorm:"foreignKey:OrganizationID;constraint:OnDelete:CASCADE" json:"organization"`
	Name                string       `json:"name"`
	Provider            string       `json:"provider"`
	Region              string       `json:"region"`
	Status              string       `json:"status"`
	EncryptedKubeconfig string       `gorm:"type:text" json:"-"`
	KubeIV              string       `json:"-"`
	KubeTag             string       `json:"-"`
	NodeGroups          []NodeGroup  `gorm:"many2many:cluster_node_groups;constraint:OnDelete:CASCADE" json:"node_groups,omitempty"`
	Timestamped
}
