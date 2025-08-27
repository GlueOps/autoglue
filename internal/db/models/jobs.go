package models

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type JobStatus string

const (
	StatusQueued    JobStatus = "queued"
	StatusRunning   JobStatus = "running"
	StatusCompleted JobStatus = "completed"
	StatusFailed    JobStatus = "failed"
)

type Job struct {
	ID          uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Type        string         `gorm:"not null"`
	Payload     datatypes.JSON `gorm:"type:jsonb;not null"`
	Status      JobStatus      `gorm:"type:job_status;not null;default:'queued'"`
	Priority    int            `gorm:"not null;default:0"`
	Attempts    int            `gorm:"not null;default:0"`
	MaxAttempts int            `gorm:"not null;default:8"`
	ScheduledAt time.Time      `gorm:"not null;default:now()"`
	LockedAt    *time.Time
	LockedBy    *string
	StartedAt   *time.Time
	FinishedAt  *time.Time
	LastError   *string
	DedupeKey   *string
}

type EnqueueOptions struct {
	Priority    int
	RunAfter    time.Duration // schedule delay
	DedupeKey   *string
	MaxAttempts int // optional override
}

func Enqueue(ctx context.Context, db *gorm.DB, typ string, payload any, opt *EnqueueOptions) (*Job, error) {
	j := &Job{
		Type:   typ,
		Status: StatusQueued,
	}
	b, _ := json.Marshal(payload)
	j.Payload = datatypes.JSON(b)

	if opt != nil {
		j.Priority = opt.Priority
		if opt.RunAfter > 0 {
			j.ScheduledAt = time.Now().Add(opt.RunAfter)
		}
		if opt.DedupeKey != nil {
			j.DedupeKey = opt.DedupeKey
		}
		if opt.MaxAttempts > 0 {
			j.MaxAttempts = opt.MaxAttempts
		}
	}

	if err := db.WithContext(ctx).Create(j).Error; err != nil {
		return nil, err
	}

	return j, nil
}
