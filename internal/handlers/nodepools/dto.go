package nodepools

import (
	"github.com/google/uuid"
)

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

type createNodePoolRequest struct {
	Name      string   `json:"name"`
	ServerIDs []string `json:"server_ids,omitempty"` // optional initial servers
}

type updateNodePoolRequest struct {
	Name *string `json:"name,omitempty"`
}

type attachServersRequest struct {
	ServerIDs []string `json:"server_ids"`
}

type taintBrief struct {
	ID     uuid.UUID `json:"id"`
	Key    string    `json:"key"`
	Value  string    `json:"value"`
	Effect string    `json:"effect"`
}

type attachTaintsRequest struct {
	TaintIDs []string `json:"taint_ids"`
}

type labelBrief struct {
	ID    uuid.UUID `json:"id"`
	Key   string    `json:"key"`
	Value string    `json:"value"`
}

type attachLabelsRequest struct {
	LabelIDs []string `json:"label_ids"`
}

type annotationBrief struct {
	ID    uuid.UUID `json:"id"`
	Key   string    `json:"key"`
	Value string    `json:"value"`
}

type attachAnnotationsRequest struct {
	AnnotationIDs []string `json:"annotation_ids"`
}
