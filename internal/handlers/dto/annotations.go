package dto

import "github.com/glueops/autoglue/internal/common"

type AnnotationResponse struct {
	common.AuditFields
	Key   string `json:"key"`
	Value string `json:"value"`
}

type CreateAnnotationRequest struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type UpdateAnnotationRequest struct {
	Key   *string `json:"key,omitempty"`
	Value *string `json:"value,omitempty"`
}
