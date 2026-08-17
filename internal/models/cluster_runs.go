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

// Who executes a run. Both the River ClusterActionWorker and the bastion agent
// read cluster_runs, and neither can tell from status alone whether the other is
// already working on a row.
//
// Without this marker the agent's assignment query claims any queued or running
// run for its cluster — including one River is part-way through executing over
// SSH — and the same make target runs twice, concurrently, against one Terraform
// state. The default is river because that is what every run created before this
// column existed was, and because claiming a run by accident must require an
// explicit opt-in rather than a silent one.
const (
	ClusterRunExecutorRiver = "river"
	ClusterRunExecutorAgent = "agent"
)

type ClusterRun struct {
	ID             uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id" format:"uuid"`
	OrganizationID uuid.UUID `json:"organization_id" gorm:"type:uuid;index"`
	ClusterID      uuid.UUID `json:"cluster_id" gorm:"type:uuid;index"`
	Action         string    `json:"action" gorm:"type:text;not null"`
	Status         string    `json:"status" gorm:"type:text;not null"`
	// Executor is set at creation and never changed. See the constants above:
	// it is the only thing stopping River and a bastion agent from both
	// executing the same run.
	Executor string `json:"executor" gorm:"type:text;not null;default:'river';index"`
	Error    string `json:"error" gorm:"type:text;not null"`
	// JobID is the River job executing this run, so logs can be correlated.
	JobID      *int64    `json:"job_id,omitempty" gorm:"index"`
	CreatedAt  time.Time `json:"created_at,omitempty" gorm:"type:timestamptz;column:created_at;not null;default:now()" format:"date-time"`
	UpdatedAt  time.Time `json:"updated_at,omitempty" gorm:"type:timestamptz;autoUpdateTime;column:updated_at;not null;default:now()" format:"date-time"`
	FinishedAt time.Time `json:"finished_at,omitempty" gorm:"type:timestamptz" format:"date-time"`
}
