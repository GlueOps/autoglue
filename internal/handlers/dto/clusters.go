package dto

import (
	"time"

	"github.com/google/uuid"
)

type ClusterResponse struct {
	ID                    uuid.UUID             `json:"id"`
	Name                  string                `json:"name"`
	CaptainDomain         *DomainResponse       `json:"captain_domain,omitempty"`
	ControlPlaneRecordSet *RecordSetResponse    `json:"control_plane_record_set,omitempty"`
	AppsLoadBalancer      *LoadBalancerResponse `json:"apps_load_balancer,omitempty"`
	GlueOpsLoadBalancer   *LoadBalancerResponse `json:"glueops_load_balancer,omitempty"`
	BastionServer         *ServerResponse       `json:"bastion_server,omitempty"`
	Provider              string                `json:"provider"`
	Region                string                `json:"region"`
	Status                string                `json:"status"`
	LastError             string                `json:"last_error"`
	RandomToken           string                `json:"random_token"`
	CertificateKey        string                `json:"certificate_key"`
	NodePools             []NodePoolResponse    `json:"node_pools,omitempty"`
	CreatedAt             time.Time             `json:"created_at"`
	UpdatedAt             time.Time             `json:"updated_at"`
}

type CreateClusterRequest struct {
	Name     string `json:"name"`
	Provider string `json:"provider"`
	Region   string `json:"region"`
}

type UpdateClusterRequest struct {
	Name     *string `json:"name,omitempty"`
	Provider *string `json:"provider,omitempty"`
	Region   *string `json:"region,omitempty"`
}

type AttachCaptainDomainRequest struct {
	DomainID uuid.UUID `json:"domain_id"`
}

type AttachRecordSetRequest struct {
	RecordSetID uuid.UUID `json:"record_set_id"`
}

type AttachLoadBalancerRequest struct {
	LoadBalancerID uuid.UUID `json:"load_balancer_id"`
}

type AttachBastionRequest struct {
	ServerID uuid.UUID `json:"server_id"`
}

type SetKubeconfigRequest struct {
	Kubeconfig string `json:"kubeconfig"`
}
