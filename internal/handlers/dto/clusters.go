package dto

import (
	"time"

	"github.com/google/uuid"
)

type ClusterResponse struct {
	ID                  uuid.UUID          `json:"id"`
	Name                string             `json:"name"`
	Provider            string             `json:"provider"`
	Region              string             `json:"region"`
	Status              string             `json:"status"`
	CaptainDomain       string             `json:"captain_domain"`
	ClusterLoadBalancer string             `json:"cluster_load_balancer"`
	RandomToken         string             `json:"random_token"`
	CertificateKey      string             `json:"certificate_key"`
	ControlLoadBalancer string             `json:"control_load_balancer"`
	NodePools           []NodePoolResponse `json:"node_pools,omitempty"`
	BastionServer       *ServerResponse    `json:"bastion_server,omitempty"`
	CreatedAt           time.Time          `json:"created_at"`
	UpdatedAt           time.Time          `json:"updated_at"`
}

type CreateClusterRequest struct {
	Name                string  `json:"name"`
	Provider            string  `json:"provider"`
	Region              string  `json:"region"`
	Status              string  `json:"status"`
	CaptainDomain       string  `json:"captain_domain"`
	ClusterLoadBalancer *string `json:"cluster_load_balancer"`
	ControlLoadBalancer *string `json:"control_load_balancer"`
}
