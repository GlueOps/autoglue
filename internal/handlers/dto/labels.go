package dto

import "github.com/google/uuid"

type LabelResponse struct {
	ID             uuid.UUID `json:"id"`
	OrganizationID uuid.UUID `json:"organization_id"`
	Key            string    `json:"key"`
	Value          string    `json:"value"`
	CreatedAt      string    `json:"created_at,omitempty"`
	UpdatedAt      string    `json:"updated_at,omitempty"`
}

type CreateLabelRequest struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type UpdateLabelRequest struct {
	Key   *string `json:"key"`
	Value *string `json:"value"`
}
