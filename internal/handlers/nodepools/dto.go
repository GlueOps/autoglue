package nodepools

import (
	"github.com/google/uuid"
)

type createNodePoolRequest struct {
	Name      string   `json:"name"`
	ServerIDs []string `json:"server_ids"`
}

type updateNodePoolRequest struct {
	Name *string `json:"name"`
}

type attachServersRequest struct {
	ServerIDs []string `json:"server_ids"`
}

type attachLabelsRequest struct {
	LabelIDs []string `json:"label_ids"`
}

type attachTaintsRequest struct {
	TaintIDs []string `json:"taint_ids"`
}

type nodePoolResponse struct {
	ID      uuid.UUID     `json:"id"`
	Name    string        `json:"name"`
	Servers []serverBrief `json:"servers,omitempty"`
}

type serverBrief struct {
	ID       uuid.UUID `json:"id"`
	Hostname string    `json:"hostname,omitempty"`
	IP       string    `json:"ip,omitempty"`
	Role     string    `json:"role,omitempty"`
	Status   string    `json:"status,omitempty"`
}

type labelBrief struct {
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

type annotationBrief struct {
	ID    uuid.UUID `json:"id"`
	Name  string    `json:"name"`
	Value string    `json:"value"`
}

type attachAnnotationsRequest struct {
	AnnotationIDs []string `json:"annotation_ids"`
}
