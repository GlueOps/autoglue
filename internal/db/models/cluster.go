package models

import "github.com/google/uuid"

type Cluster struct {
	ID                  uuid.UUID    `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	OrganizationID      uuid.UUID    `gorm:"type:uuid;not null" json:"organization_id"`
	Organization        Organization `gorm:"foreignKey:OrganizationID;constraint:OnDelete:CASCADE" json:"organization"`
	Name                string       `gorm:"not null" json:"name"`
	Provider            string       `json:"provider"`
	Region              string       `json:"region"`
	Status              string       `json:"status"`
	EncryptedKubeconfig string       `gorm:"type:text" json:"-"`
	KubeIV              string       `json:"-"`
	KubeTag             string       `json:"-"`
	NodePools           []NodePool   `gorm:"many2many:cluster_node_pools;constraint:OnDelete:CASCADE" json:"node_pools,omitempty"`
	BastionServerID     *uuid.UUID   `gorm:"type:uuid" json:"bastion_server_id,omitempty"`
	BastionServer       *Server      `gorm:"foreignKey:BastionServerID" json:"bastion_server,omitempty"`
	ClusterLoadBalancer string       `gorm:"type:text" json:"cluster_load_balancer"`
	ControlLoadBalancer string       `gorm:"type:text" json:"control_load_balancer"`
}
