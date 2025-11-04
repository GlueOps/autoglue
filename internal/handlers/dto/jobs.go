package dto

import (
	"encoding/json"
	"time"
)

type JobStatus string

const (
	StatusQueued    JobStatus = "queued"
	StatusRunning   JobStatus = "running"
	StatusSucceeded JobStatus = "succeeded"
	StatusFailed    JobStatus = "failed"
	StatusCanceled  JobStatus = "canceled"
	StatusRetrying  JobStatus = "retrying"
	StatusScheduled JobStatus = "scheduled"
)

// Job represents a background job managed by Archer.
// swagger:model Job
type Job struct {
	ID          string     `json:"id" example:"01HF7SZK8Z8WG1M3J7S2Z8M2N6"`
	Type        string     `json:"type" example:"email.send"`
	Queue       string     `json:"queue" example:"default"`
	Status      JobStatus  `json:"status" example:"queued" enums:"queued|running|succeeded|failed|canceled|retrying|scheduled"`
	Attempts    int        `json:"attempts" example:"0"`
	MaxAttempts int        `json:"max_attempts,omitempty" example:"3"`
	CreatedAt   time.Time  `json:"created_at" example:"2025-11-04T09:30:00Z"`
	UpdatedAt   *time.Time `json:"updated_at,omitempty" example:"2025-11-04T09:30:00Z"`
	LastError   *string    `json:"last_error,omitempty" example:"error message"`
	RunAt       *time.Time `json:"run_at,omitempty" example:"2025-11-04T09:30:00Z"`
	Payload     any        `json:"payload,omitempty"`
}

// QueueInfo holds queue-level counts.
// swagger:model QueueInfo
type QueueInfo struct {
	Name      string `json:"name" example:"default"`
	Pending   int    `json:"pending" example:"42"`
	Running   int    `json:"running" example:"3"`
	Failed    int    `json:"failed" example:"5"`
	Scheduled int    `json:"scheduled" example:"7"`
}

// PageJob is a concrete paginated response for Job (generics not supported by swag).
// swagger:model PageJob
type PageJob struct {
	Items    []Job `json:"items"`
	Total    int   `json:"total" example:"120"`
	Page     int   `json:"page" example:"1"`
	PageSize int   `json:"page_size" example:"25"`
}

// EnqueueRequest is the POST body for creating a job.
// swagger:model EnqueueRequest
type EnqueueRequest struct {
	Queue   string          `json:"queue" example:"default"`
	Type    string          `json:"type" example:"email.send"`
	Payload json.RawMessage `json:"payload"`
	RunAt   *time.Time      `json:"run_at" example:"2025-11-05T08:00:00Z"`
}
