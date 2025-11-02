package dto

import "github.com/google/uuid"

type TaintResponse struct {
	ID     uuid.UUID `json:"id"`
	Key    string    `json:"key"`
	Value  string    `json:"value"`
	Effect string    `json:"effect"`
}

type CreateTaintRequest struct {
	Key    string `json:"key"`
	Value  string `json:"value"`
	Effect string `json:"effect"`
}

type UpdateTaintRequest struct {
	Key    *string `json:"key,omitempty"`
	Value  *string `json:"value,omitempty"`
	Effect *string `json:"effect,omitempty"`
}
