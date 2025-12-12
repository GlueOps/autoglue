package bg

import (
	"context"
	"time"

	"github.com/dyaksa/archer"
	"github.com/dyaksa/archer/job"
	"github.com/glueops/autoglue/internal/models"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

type OrgKeySweeperArgs struct {
	IntervalS     int `json:"interval_seconds,omitempty"`
	RetentionDays int `json:"retention_days,omitempty"`
}

type OrgKeySweeperResult struct {
	Status           string `json:"status"`
	MarkedRevoked    int    `json:"marked_revoked"`
	DeletedEphemeral int    `json:"deleted_ephemeral"`
	ElapsedMs        int    `json:"elapsed_ms"`
}

func OrgKeySweeperWorker(db *gorm.DB, jobs *Jobs) archer.WorkerFn {
	return func(ctx context.Context, j job.Job) (any, error) {
		args := OrgKeySweeperArgs{
			IntervalS:     3600,
			RetentionDays: 10,
		}
		start := time.Now()

		_ = j.ParseArguments(&args)
		if args.IntervalS <= 0 {
			args.IntervalS = 3600
		}
		if args.RetentionDays <= 0 {
			args.RetentionDays = 10
		}

		now := time.Now()

		// 1) Mark expired keys as revoked
		res1 := db.Model(&models.APIKey{}).
			Where("expires_at IS NOT NULL AND expires_at <= ? AND revoked = false", now).
			Updates(map[string]any{
				"revoked":    true,
				"updated_at": now,
			})

		if res1.Error != nil {
			log.Error().Err(res1.Error).Msg("[org_key_sweeper] mark expired revoked failed")
			return nil, res1.Error
		}
		markedRevoked := int(res1.RowsAffected)

		// 2) Hard-delete ephemeral keys that are revoked and older than retention
		cutoff := now.Add(-time.Duration(args.RetentionDays) * 24 * time.Hour)
		res2 := db.
			Where("is_ephemeral = ? AND revoked = ? AND updated_at <= ?", true, true, cutoff).
			Delete(&models.APIKey{})

		if res2.Error != nil {
			log.Error().Err(res2.Error).Msg("[org_key_sweeper] delete revoked ephemeral keys failed")
			return nil, res2.Error
		}
		deletedEphemeral := int(res2.RowsAffected)

		out := OrgKeySweeperResult{
			Status:           "ok",
			MarkedRevoked:    markedRevoked,
			DeletedEphemeral: deletedEphemeral,
			ElapsedMs:        int(time.Since(start).Milliseconds()),
		}

		log.Info().
			Int("marked_revoked", markedRevoked).
			Int("deleted_ephemeral", deletedEphemeral).
			Msg("[org_key_sweeper] cleanup tick ok")

		// Re-enqueue the sweeper
		next := time.Now().Add(time.Duration(args.IntervalS) * time.Second)
		_, _ = jobs.Enqueue(
			ctx,
			uuid.NewString(),
			"org_key_sweeper",
			args,
			archer.WithScheduleTime(next),
			archer.WithMaxRetries(1),
		)
		return out, nil
	}
}
