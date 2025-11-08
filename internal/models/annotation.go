package models

import (
	"github.com/glueops/autoglue/internal/common"
)

type Annotation struct {
	common.AuditFields `gorm:"embedded"`
	Organization       Organization `gorm:"foreignKey:OrganizationID;constraint:OnDelete:CASCADE" json:"organization"`
	Key                string       `gorm:"not null" json:"key"`
	Value              string       `gorm:"not null" json:"value"`
	NodePools          []NodePool   `gorm:"many2many:node_annotations;constraint:OnDelete:CASCADE" json:"node_pools,omitempty"`
}
