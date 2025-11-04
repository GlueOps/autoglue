package models

import (
	"github.com/glueops/autoglue/internal/common"
)

type Label struct {
	common.AuditFields
	Organization Organization `gorm:"foreignKey:OrganizationID;constraint:OnDelete:CASCADE" json:"organization"`
	Key          string       `gorm:"not null" json:"key"`
	Value        string       `gorm:"not null" json:"value"`
	NodePools    []NodePool   `gorm:"many2many:node_labels;constraint:OnDelete:CASCADE" json:"servers,omitempty"`
}
