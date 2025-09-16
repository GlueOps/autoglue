package clusters

import "github.com/google/uuid"

// clusterResponse describes a cluster with optional expansions.
// swagger:model clusters.clusterResponse
type clusterResponse struct {
	ID                  uuid.UUID       `json:"id"`
	Name                string          `json:"name"`
	Provider            string          `json:"provider"`
	Region              string          `json:"region"`
	Status              string          `json:"status"`
	ClusterLoadBalancer string          `json:"cluster_load_balancer"`
	ControlLoadBalancer string          `json:"control_load_balancer"`
	NodePools           []nodePoolBrief `json:"node_pools,omitempty"`
	BastionServer       *serverBrief    `json:"bastion_server,omitempty"`
}

type serverBrief struct {
	ID       uuid.UUID `json:"id"`
	Hostname string    `json:"hostname"`
	IP       string    `json:"ip"`
	Role     string    `json:"role"`
	Status   string    `json:"status"`
}

type nodePoolBrief struct {
	ID          uuid.UUID         `json:"id"`
	Name        string            `json:"name"`
	Labels      []labelBrief      `json:"labels,omitempty"`
	Annotations []annotationBrief `json:"annotations,omitempty"`
	Taints      []taintBrief      `json:"taints,omitempty"`
	Servers     []serverBrief     `json:"servers,omitempty"`
}

type labelBrief struct {
	ID    uuid.UUID `json:"id"`
	Key   string    `json:"key"`
	Value string    `json:"value"`
}

type annotationBrief struct {
	ID    uuid.UUID `json:"id"`
	Key   string    `json:"key"`
	Value string    `json:"value"`
}

type taintBrief struct {
	ID     uuid.UUID `json:"id"`
	Key    string    `json:"key"`
	Value  string    `json:"value"`
	Effect string    `json:"effect"`
}

// swagger:model clusters.updateClusterRequest
type updateClusterRequest struct {
	Name                *string `json:"name"`
	Provider            *string `json:"provider"`
	Region              *string `json:"region"`
	Status              *string `json:"status"`
	BastionServerID     *string `json:"bastion_server_id"`
	ClusterLoadBalancer *string `json:"cluster_load_balancer"`
	ControlLoadBalancer *string `json:"control_load_balancer"`
	Kubeconfig          *string `json:"kubeconfig"`
}

// swagger:model clusters.attachNodePoolsRequest
type attachNodePoolsRequest struct {
	NodePoolIDs []string `json:"node_pool_ids"`
}

// swagger:model clusters.setBastionRequest
type setBastionRequest struct {
	ServerID string `json:"server_id"`
}

// swagger:model clusters.createClusterRequest
type createClusterRequest struct {
	Name                string   `json:"name"`
	Provider            string   `json:"provider"`
	Region              string   `json:"region"`
	NodePoolIDs         []string `json:"node_pool_ids"`
	BastionServerID     *string  `json:"bastion_server_id"`
	ClusterLoadBalancer *string  `json:"cluster_load_balancer"`
	ControlLoadBalancer *string  `json:"control_load_balancer"`
	Kubeconfig          *string  `json:"kubeconfig"`
}
