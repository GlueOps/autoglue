package dto

import "github.com/google/uuid"

type LabelResponse struct {
	ID    uuid.UUID `json:"id"`
	Key   string    `json:"key"`
	Value string    `json:"value"`
}

type CreateLabelRequest struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type UpdateLabelRequest struct {
	Key   *string `json:"key"`
	Value *string `json:"value"`
}
