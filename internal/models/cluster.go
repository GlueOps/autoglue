package models

import (
	"time"

	"github.com/google/uuid"
)

type Cluster struct {
	ID                  uuid.UUID    `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	OrganizationID      uuid.UUID    `gorm:"type:uuid;not null" json:"organization_id"`
	Organization        Organization `gorm:"foreignKey:OrganizationID;constraint:OnDelete:CASCADE" json:"organization"`
	Name                string       `gorm:"not null" json:"name"`
	Provider            string       `json:"provider"`
	Region              string       `json:"region"`
	Status              string       `json:"status"`
	CaptainDomain       string       `gorm:"not null" json:"captain_domain"`
	ClusterLoadBalancer string       `json:"cluster_load_balancer"`
	RandomToken         string       `json:"random_token"`
	CertificateKey      string       `json:"certificate_key"`
	EncryptedKubeconfig string       `gorm:"type:text" json:"-"`
	KubeIV              string       `json:"-"`
	KubeTag             string       `json:"-"`
	NodePools           []NodePool   `gorm:"many2many:cluster_node_pools;constraint:OnDelete:CASCADE" json:"node_pools,omitempty"`
	BastionServerID     *uuid.UUID   `gorm:"type:uuid" json:"bastion_server_id,omitempty"`
	BastionServer       *Server      `gorm:"foreignKey:BastionServerID" json:"bastion_server,omitempty"`
	CreatedAt           time.Time    `json:"created_at,omitempty" gorm:"type:timestamptz;column:created_at;not null;default:now()"`
	UpdatedAt           time.Time    `json:"updated_at,omitempty" gorm:"type:timestamptz;autoUpdateTime;column:updated_at;not null;default:now()"`
}
