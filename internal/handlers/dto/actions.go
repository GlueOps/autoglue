package dto

import (
	"time"

	"github.com/google/uuid"
)

type ActionResponse struct {
	ID          uuid.UUID `json:"id" format:"uuid"`
	Label       string    `json:"label"`
	Description string    `json:"description"`
	MakeTarget  string    `json:"make_target"`
	CreatedAt   time.Time `json:"created_at" format:"date-time"`
	UpdatedAt   time.Time `json:"updated_at" format:"date-time"`
}

type CreateActionRequest struct {
	Label       string `json:"label"`
	Description string `json:"description"`
	MakeTarget  string `json:"make_target"`
}

type UpdateActionRequest struct {
	Label       *string `json:"label,omitempty"`
	Description *string `json:"description,omitempty"`
	MakeTarget  *string `json:"make_target,omitempty"`
}
