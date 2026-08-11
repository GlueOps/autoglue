package bg

import (
	"context"
	"time"

	"github.com/glueops/autoglue/internal/models"
	"github.com/riverqueue/river"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

type OrgKeySweeperArgs struct {
	RetentionDays int `json:"retention_days,omitempty"`
}

func (OrgKeySweeperArgs) Kind() string { return "org_key_sweeper" }

// OrgKeySweeperResult is recorded on the job row via river.RecordOutput.
type OrgKeySweeperResult struct {
	Status           string `json:"status"`
	MarkedRevoked    int    `json:"marked_revoked"`
	DeletedEphemeral int    `json:"deleted_ephemeral"`
	ElapsedMs        int    `json:"elapsed_ms"`
}

func (OrgKeySweeperArgs) InsertOpts() river.InsertOpts {
	return river.InsertOpts{Queue: QueueMaintenance, MaxAttempts: 2}
}

type OrgKeySweeperWorker struct {
	river.WorkerDefaults[OrgKeySweeperArgs]
	db *gorm.DB
}

func (w *OrgKeySweeperWorker) Timeout(*river.Job[OrgKeySweeperArgs]) time.Duration {
	return 5 * time.Minute
}

func (w *OrgKeySweeperWorker) Work(ctx context.Context, job *river.Job[OrgKeySweeperArgs]) error {
	start := time.Now()
	retentionDays := job.Args.RetentionDays
	if retentionDays <= 0 {
		retentionDays = 10
	}

	now := time.Now()

	// 1) Mark expired keys as revoked
	res1 := w.db.Model(&models.APIKey{}).
		Where("expires_at IS NOT NULL AND expires_at <= ? AND revoked = false", now).
		Updates(map[string]any{
			"revoked":    true,
			"updated_at": now,
		})

	if res1.Error != nil {
		log.Error().Err(res1.Error).Msg("[org_key_sweeper] mark expired revoked failed")
		return res1.Error
	}
	markedRevoked := int(res1.RowsAffected)

	// 2) Hard-delete ephemeral keys that are revoked and older than retention
	cutoff := now.Add(-time.Duration(retentionDays) * 24 * time.Hour)
	res2 := w.db.
		Where("is_ephemeral = ? AND revoked = ? AND updated_at <= ?", true, true, cutoff).
		Delete(&models.APIKey{})

	if res2.Error != nil {
		log.Error().Err(res2.Error).Msg("[org_key_sweeper] delete revoked ephemeral keys failed")
		return res2.Error
	}
	deletedEphemeral := int(res2.RowsAffected)

	log.Info().
		Int("marked_revoked", markedRevoked).
		Int("deleted_ephemeral", deletedEphemeral).
		Msg("[org_key_sweeper] cleanup tick ok")

	if err := river.RecordOutput(ctx, OrgKeySweeperResult{
		Status:           "ok",
		MarkedRevoked:    markedRevoked,
		DeletedEphemeral: deletedEphemeral,
		ElapsedMs:        int(time.Since(start).Milliseconds()),
	}); err != nil {
		log.Warn().Err(err).Msg("[org_key_sweeper] could not record output")
	}
	return nil
}
