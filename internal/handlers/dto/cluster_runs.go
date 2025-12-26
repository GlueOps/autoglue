package dto

import (
	"time"

	"github.com/google/uuid"
)

type ClusterRunResponse struct {
	ID             uuid.UUID  `json:"id" format:"uuid"`
	OrganizationID uuid.UUID  `json:"organization_id" format:"uuid"`
	ClusterID      uuid.UUID  `json:"cluster_id" format:"uuid"`
	Action         string     `json:"action"`
	Status         string     `json:"status"`
	Error          string     `json:"error"`
	CreatedAt      time.Time  `json:"created_at" format:"date-time"`
	UpdatedAt      time.Time  `json:"updated_at" format:"date-time"`
	FinishedAt     *time.Time `json:"finished_at,omitempty" format:"date-time"`
}
