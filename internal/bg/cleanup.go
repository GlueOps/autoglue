package bg

import (
	"context"
	"time"

	"github.com/dyaksa/archer"
	"github.com/dyaksa/archer/job"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CleanupArgs struct {
	RetainDays int    `json:"retain_days"`
	Table      string `json:"table"`
}

type JobRow struct {
	ID        string `gorm:"primaryKey"`
	Status    string
	UpdatedAt time.Time
}

func (JobRow) TableName() string { return "jobs" }

func CleanupWorker(gdb *gorm.DB, jobs *Jobs, retainDays int) archer.WorkerFn {
	return func(ctx context.Context, j job.Job) (any, error) {
		if err := CleanupJobs(gdb, retainDays); err != nil {
			return nil, err
		}

		// schedule tomorrow 03:30
		next := time.Now().Truncate(24 * time.Hour).Add(24 * time.Hour).Add(3*time.Hour + 30*time.Minute)

		_, _ = jobs.Enqueue(
			ctx,
			uuid.NewString(),
			"archer_cleanup",
			CleanupArgs{RetainDays: retainDays, Table: "jobs"},
			archer.WithScheduleTime(next),
			archer.WithMaxRetries(1),
		)
		return nil, nil
	}
}

func CleanupJobs(db *gorm.DB, retainDays int) error {
	cutoff := time.Now().AddDate(0, 0, -retainDays)
	return db.
		Where("status IN ?", []string{"success", "failed", "cancelled"}).
		Where("updated_at < ?", cutoff).
		Delete(&JobRow{}).Error
}
