package bg

import (
	"context"
	"time"

	"github.com/glueops/autoglue/internal/models"
	"github.com/riverqueue/river"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

type JobLogsCleanupArgs struct {
	RetainDays int `json:"retain_days,omitempty"`
}

func (JobLogsCleanupArgs) Kind() string { return "job_logs_cleanup" }

func (JobLogsCleanupArgs) InsertOpts() river.InsertOpts {
	return river.InsertOpts{Queue: QueueMaintenance, MaxAttempts: 2}
}

// JobLogsCleanupResult is recorded on the job row via river.RecordOutput.
type JobLogsCleanupResult struct {
	Status  string `json:"status"`
	Deleted int64  `json:"deleted"`
}

// JobLogsCleanupWorker prunes job_logs. River expires its own job rows, but
// job_logs is our table and nothing else would ever remove a transcript.
type JobLogsCleanupWorker struct {
	river.WorkerDefaults[JobLogsCleanupArgs]
	db *gorm.DB
}

func (w *JobLogsCleanupWorker) Timeout(*river.Job[JobLogsCleanupArgs]) time.Duration {
	return 10 * time.Minute
}

func (w *JobLogsCleanupWorker) Work(ctx context.Context, job *river.Job[JobLogsCleanupArgs]) error {
	retainDays := job.Args.RetainDays
	if retainDays <= 0 {
		retainDays = jobLogRetainDays()
	}

	cutoff := time.Now().AddDate(0, 0, -retainDays)
	res := w.db.Where("created_at < ?", cutoff).Delete(&models.JobLog{})
	if res.Error != nil {
		return res.Error
	}

	if res.RowsAffected > 0 {
		log.Info().
			Int64("deleted", res.RowsAffected).
			Int("retain_days", retainDays).
			Msg("[job_logs_cleanup] pruned old job output")
	}

	if err := river.RecordOutput(ctx, JobLogsCleanupResult{
		Status:  "ok",
		Deleted: res.RowsAffected,
	}); err != nil {
		log.Warn().Err(err).Msg("[job_logs_cleanup] could not record output")
	}
	return nil
}
