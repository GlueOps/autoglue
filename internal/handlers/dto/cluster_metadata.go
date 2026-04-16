package dto

import "github.com/glueops/autoglue/internal/common"

type ClusterMetadataResponse struct {
	common.AuditFields
	ClusterID string `json:"cluster_id"`
	Key       string `json:"key"`
	Value     string `json:"value"`
}

type CreateClusterMetadataRequest struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type UpdateClusterMetadataRequest struct {
	Key   *string `json:"key,omitempty"`
	Value *string `json:"value,omitempty"`
}
