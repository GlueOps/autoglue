package models

import (
	"time"

	"github.com/google/uuid"
)

const (
	ClusterStatusPrePending   = "pre_pending" // needs validation
	ClusterStatusIncomplete   = "incomplete"  // invalid/missing shape
	ClusterStatusPending      = "pending"     // valid shape, waiting for provisioning
	ClusterStatusProvisioning = "provisioning"
	ClusterStatusReady        = "ready"
	ClusterStatusFailed       = "failed" // provisioning/runtime failure
)

type Cluster struct {
	ID             uuid.UUID    `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	OrganizationID uuid.UUID    `gorm:"type:uuid;not null" json:"organization_id"`
	Organization   Organization `gorm:"foreignKey:OrganizationID;constraint:OnDelete:CASCADE" json:"organization"`

	Name     string `gorm:"not null" json:"name"`
	Provider string `json:"provider"`
	Region   string `json:"region"`

	Status    string `gorm:"type:varchar(20);not null;default:'pre_pending'" json:"status"`
	LastError string `gorm:"type:text;not null;default:''" json:"last_error"`

	CaptainDomainID         *uuid.UUID    `gorm:"type:uuid" json:"captain_domain_id"`
	CaptainDomain           Domain        `gorm:"foreignKey:CaptainDomainID" json:"captain_domain"`
	ControlPlaneRecordSetID *uuid.UUID    `gorm:"type:uuid" json:"control_plane_record_set_id,omitempty"`
	ControlPlaneRecordSet   *RecordSet    `gorm:"foreignKey:ControlPlaneRecordSetID" json:"control_plane_record_set,omitempty"`
	AppsLoadBalancerID      *uuid.UUID    `gorm:"type:uuid" json:"apps_load_balancer_id,omitempty"`
	AppsLoadBalancer        *LoadBalancer `gorm:"foreignKey:AppsLoadBalancerID" json:"apps_load_balancer,omitempty"`
	GlueOpsLoadBalancerID   *uuid.UUID    `gorm:"type:uuid" json:"glueops_load_balancer_id,omitempty"`
	GlueOpsLoadBalancer     *LoadBalancer `gorm:"foreignKey:GlueOpsLoadBalancerID" json:"glueops_load_balancer,omitempty"`
	BastionServerID         *uuid.UUID    `gorm:"type:uuid" json:"bastion_server_id,omitempty"`
	BastionServer           *Server       `gorm:"foreignKey:BastionServerID" json:"bastion_server,omitempty"`

	NodePools []NodePool `gorm:"many2many:cluster_node_pools;constraint:OnDelete:CASCADE" json:"node_pools,omitempty"`

	RandomToken    string `json:"random_token"`
	CertificateKey string `json:"certificate_key"`

	EncryptedKubeconfig string `gorm:"type:text" json:"-"`
	KubeIV              string `json:"-"`
	KubeTag             string `json:"-"`

	CreatedAt time.Time `json:"created_at,omitempty" gorm:"type:timestamptz;column:created_at;not null;default:now()"`
	UpdatedAt time.Time `json:"updated_at,omitempty" gorm:"type:timestamptz;autoUpdateTime;column:updated_at;not null;default:now()"`
}
