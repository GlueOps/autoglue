package dto

import "github.com/glueops/autoglue/internal/common"

type NodeRole string

const (
	NodeRoleMaster NodeRole = "master"
	NodeRoleWorker NodeRole = "worker"
)

type CreateNodePoolRequest struct {
	Name string   `json:"name"`
	Role NodeRole `json:"role" enums:"master,worker" swaggertype:"string"`
}

type UpdateNodePoolRequest struct {
	Name *string   `json:"name"`
	Role *NodeRole `json:"role" enums:"master,worker" swaggertype:"string"`
}

type NodePoolResponse struct {
	common.AuditFields
	Name        string               `json:"name"`
	Role        NodeRole             `json:"role" enums:"master,worker" swaggertype:"string"`
	Servers     []ServerResponse     `json:"servers"`
	Annotations []AnnotationResponse `json:"annotations"`
	Labels      []LabelResponse      `json:"labels"`
	Taints      []TaintResponse      `json:"taints"`
}

type AttachServersRequest struct {
	ServerIDs []string `json:"server_ids"`
}

type AttachTaintsRequest struct {
	TaintIDs []string `json:"taint_ids"`
}

type AttachLabelsRequest struct {
	LabelIDs []string `json:"label_ids"`
}

type AttachAnnotationsRequest struct {
	AnnotationIDs []string `json:"annotation_ids"`
}
