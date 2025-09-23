package labels

import "github.com/google/uuid"

type labelResponse struct {
	ID         uuid.UUID       `json:"id"`
	Key        string          `json:"key"`
	Value      string          `json:"value"`
	NodeGroups []nodePoolBrief `json:"node_groups,omitempty"`
}

type nodePoolBrief struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

type createLabelRequest struct {
	Key         string   `json:"key"`
	Value       string   `json:"value"`
	NodePoolIDs []string `json:"node_pool_ids,omitempty"`
}

type updateLabelRequest struct {
	Key   *string `json:"key"`
	Value *string `json:"value"`
}

type addLabelToPoolRequest struct {
	NodePoolIDs []string `json:"node_pool_ids"`
}
