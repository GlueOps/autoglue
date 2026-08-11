package models

import (
	"time"

	"github.com/google/uuid"
)

// Log subjects. A job log is attached to whatever the user would go looking
// for it under, not to the job: a bastion bootstrap has no ClusterRun, and a
// single tick bootstraps several servers.
const (
	JobLogSubjectServer     = "server"
	JobLogSubjectClusterRun = "cluster_run"
)

// Log streams.
const (
	JobLogStreamStdout = "stdout"
	JobLogStreamSystem = "system"
)

// JobLog is one chunk of output produced by a background job.
//
// Rows are append-only and the autoincrement ID doubles as the polling cursor:
// a reader asks for everything with id > after. Chunks are batched by the
// writer rather than written per line, so a chatty `make` target does not turn
// into thousands of INSERTs.
type JobLog struct {
	ID int64 `gorm:"primaryKey;autoIncrement" json:"id"`

	// JobID is the River job that produced this output. Zero when the writer
	// had no job context.
	JobID int64 `gorm:"not null;default:0;index" json:"job_id"`

	// OrganizationID scopes reads. Carried on the row so the log endpoints do
	// not have to join back through the subject to authorize.
	OrganizationID uuid.UUID `gorm:"type:uuid;not null;index:idx_job_logs_subject,priority:1" json:"organization_id"`

	SubjectType string    `gorm:"type:text;not null;index:idx_job_logs_subject,priority:2" json:"subject_type"`
	SubjectID   uuid.UUID `gorm:"type:uuid;not null;index:idx_job_logs_subject,priority:3" json:"subject_id"`

	Stream string `gorm:"type:text;not null;default:'stdout'" json:"stream"`
	Chunk  string `gorm:"type:text;not null" json:"chunk"`

	CreatedAt time.Time `gorm:"type:timestamptz;not null;default:now();index" json:"created_at"`
}

func (JobLog) TableName() string { return "job_logs" }
