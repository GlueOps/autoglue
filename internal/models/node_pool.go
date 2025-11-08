package models

import (
	"github.com/glueops/autoglue/internal/common"
)

type NodePool struct {
	common.AuditFields
	Organization Organization `gorm:"foreignKey:OrganizationID;constraint:OnDelete:CASCADE" json:"organization"`
	Name         string       `gorm:"not null" json:"name"`
	Servers      []Server     `gorm:"many2many:node_servers;constraint:OnDelete:CASCADE" json:"servers,omitempty"`
	Annotations  []Annotation `gorm:"many2many:node_annotations;constraint:OnDelete:CASCADE" json:"annotations,omitempty"`
	Labels       []Label      `gorm:"many2many:node_labels;constraint:OnDelete:CASCADE" json:"labels,omitempty"`
	Taints       []Taint      `gorm:"many2many:node_taints;constraint:OnDelete:CASCADE" json:"taints,omitempty"`
	//Clusters       []Cluster    `gorm:"many2many:cluster_node_pools;constraint:OnDelete:CASCADE" json:"clusters,omitempty"`
	//Topology  string    `gorm:"not null,default:'stacked'" json:"topology,omitempty"` // stacked or external
	Role string `gorm:"not null,default:'worker'" json:"role,omitempty"` // master, worker, or etcd (etcd only if topology = external
}
