package taints

import "github.com/google/uuid"

type taintResponse struct {
	ID         uuid.UUID       `json:"id"`
	Key        string          `json:"key"`
	Value      string          `json:"value"`
	Effect     string          `json:"effect"`
	NodeGroups []nodePoolBrief `json:"node_groups,omitempty"`
}

type nodePoolBrief struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

type createTaintRequest struct {
	Key         string   `json:"key"`
	Value       string   `json:"value"`
	Effect      string   `json:"effect"`
	NodePoolIDs []string `json:"node_pool_ids,omitempty"`
}

type updateTaintRequest struct {
	Key    *string `json:"key,omitempty"`
	Value  *string `json:"value,omitempty"`
	Effect *string `json:"effect,omitempty"`
}

type addTaintToPoolRequest struct {
	NodePoolIDs []string `json:"node_pool_ids"`
}

type nodePoolResponse struct {
	ID      uuid.UUID     `json:"id"`
	Name    string        `json:"name"`
	Servers []serverBrief `json:"servers,omitempty"`
}

type serverBrief struct {
	ID       uuid.UUID `json:"id"`
	Hostname string    `json:"hostname"`
	IP       string    `json:"ip"`
	Role     string    `json:"role"`
	Status   string    `json:"status"`
}
