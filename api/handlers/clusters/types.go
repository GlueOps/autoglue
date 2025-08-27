package clusters

import (
	"github.com/google/uuid"
)

type nodeGroupBrief struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

type nodeGroupIds struct {
	NodeGroupIDs []string `json:"node_group_ids"`
}

type createClusterRequest struct {
	Name         string   `json:"name"`
	Provider     string   `json:"provider"`
	Region       string   `json:"region"`
	Status       string   `json:"status,omitempty"`
	Kubeconfig   string   `json:"kubeconfig,omitempty"`
	NodeGroupIDs []string `json:"node_group_ids,omitempty"` // CHANGED: node groups
}

type updateClusterRequest struct {
	Name       *string `json:"name,omitempty"`
	Provider   *string `json:"provider,omitempty"`
	Region     *string `json:"region,omitempty"`
	Status     *string `json:"status,omitempty"`
	Kubeconfig *string `json:"kubeconfig,omitempty"`
}

type clusterResponse struct {
	ID         uuid.UUID        `json:"id"`
	Name       string           `json:"name"`
	Provider   string           `json:"provider"`
	Region     string           `json:"region"`
	Status     string           `json:"status"`
	Kubeconfig string           `json:"kubeconfig,omitempty"`
	NodeGroups []nodeGroupBrief `json:"node_groups,omitempty"`
	CreatedAt  string           `json:"created_at"`
	UpdatedAt  string           `json:"updated_at"`
}
