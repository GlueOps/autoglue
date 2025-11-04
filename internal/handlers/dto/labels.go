package dto

import (
	"github.com/glueops/autoglue/internal/common"
)

type LabelResponse struct {
	common.AuditFields
	Key   string `json:"key"`
	Value string `json:"value"`
}

type CreateLabelRequest struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type UpdateLabelRequest struct {
	Key   *string `json:"key"`
	Value *string `json:"value"`
}
