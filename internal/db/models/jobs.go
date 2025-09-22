package models

import (
	"time"

	"gorm.io/datatypes"
)

type Job struct {
	ID            string         `gorm:"type:varchar;primaryKey" json:"id"` // no default; supply from app
	QueueName     string         `gorm:"type:varchar;not null" json:"queue_name"`
	Status        string         `gorm:"type:varchar;not null" json:"status"`
	Arguments     datatypes.JSON `gorm:"type:jsonb;not null;default:'{}'"`
	Result        datatypes.JSON `gorm:"type:jsonb;not null;default:'{}'"`
	LastError     *string        `gorm:"type:varchar"`
	RetryCount    int            `gorm:"not null;default:0"`
	MaxRetry      int            `gorm:"not null;default:0"`
	RetryInterval int            `gorm:"not null;default:0"`
	ScheduledAt   time.Time      `gorm:"type:timestamptz;default:now();index"`
	StartedAt     *time.Time     `gorm:"type:timestamptz;index"`
	Timestamped
}
