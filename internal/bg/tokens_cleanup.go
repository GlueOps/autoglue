package bg

import (
	"context"
	"time"

	"github.com/dyaksa/archer"
	"github.com/dyaksa/archer/job"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RefreshTokenRow struct {
	ID        string `gorm:"primaryKey"`
	RevokedAt *time.Time
	ExpiresAt time.Time
	UpdatedAt time.Time
}

func (RefreshTokenRow) TableName() string { return "refresh_tokens" }

type TokensCleanupArgs struct {
	// kept in case you want to change retention or add dry-run later
}

func TokensCleanupWorker(db *gorm.DB, jobs *Jobs) archer.WorkerFn {
	return func(ctx context.Context, j job.Job) (any, error) {
		if err := CleanupRefreshTokens(db); err != nil {
			return nil, err
		}

		// schedule tomorrow 03:45
		next := time.Now().Truncate(24 * time.Hour).Add(24 * time.Hour).Add(3*time.Hour + 45*time.Minute)
		_, _ = jobs.Enqueue(
			ctx,
			uuid.NewString(),
			"tokens_cleanup",
			TokensCleanupArgs{},
			archer.WithScheduleTime(next),
			archer.WithMaxRetries(1),
		)
		return nil, nil
	}
}

func CleanupRefreshTokens(db *gorm.DB) error {
	now := time.Now()
	return db.
		Where("revoked_at IS NOT NULL OR expires_at < ?", now).
		Delete(&RefreshTokenRow{}).Error
}
