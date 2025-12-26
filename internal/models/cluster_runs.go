package models

import (
	"time"

	"github.com/google/uuid"
)

const (
	ClusterRunStatusQueued   = "queued"
	ClusterRunStatusRunning  = "running"
	ClusterRunStatusSuccess  = "success"
	ClusterRunStatusFailed   = "failed"
	ClusterRunStatusCanceled = "canceled"
)

type ClusterRun struct {
	ID             uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id" format:"uuid"`
	OrganizationID uuid.UUID `json:"organization_id" gorm:"type:uuid;index"`
	ClusterID      uuid.UUID `json:"cluster_id" gorm:"type:uuid;index"`
	Action         string    `json:"action" gorm:"type:text;not null"`
	Status         string    `json:"status" gorm:"type:text;not null"`
	Error          string    `json:"error" gorm:"type:text;not null"`
	CreatedAt      time.Time `json:"created_at,omitempty" gorm:"type:timestamptz;column:created_at;not null;default:now()" format:"date-time"`
	UpdatedAt      time.Time `json:"updated_at,omitempty" gorm:"type:timestamptz;autoUpdateTime;column:updated_at;not null;default:now()" format:"date-time"`
	FinishedAt     time.Time `json:"finished_at,omitempty" gorm:"type:timestamptz" format:"date-time"`
}
