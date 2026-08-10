package bg

import (
	"context"
	"time"

	"github.com/riverqueue/river"
	"gorm.io/gorm"
)

type RefreshTokenRow struct {
	ID        string `gorm:"primaryKey"`
	RevokedAt *time.Time
	ExpiresAt time.Time
	UpdatedAt time.Time
}

func (RefreshTokenRow) TableName() string { return "refresh_tokens" }

type TokensCleanupArgs struct{}

func (TokensCleanupArgs) Kind() string { return "tokens_cleanup" }

func (TokensCleanupArgs) InsertOpts() river.InsertOpts {
	return river.InsertOpts{Queue: QueueMaintenance, MaxAttempts: 2}
}

type TokensCleanupWorker struct {
	river.WorkerDefaults[TokensCleanupArgs]
	db *gorm.DB
}

func (w *TokensCleanupWorker) Timeout(*river.Job[TokensCleanupArgs]) time.Duration {
	return 5 * time.Minute
}

func (w *TokensCleanupWorker) Work(_ context.Context, _ *river.Job[TokensCleanupArgs]) error {
	return CleanupRefreshTokens(w.db)
}

func CleanupRefreshTokens(db *gorm.DB) error {
	now := time.Now()
	return db.
		Where("revoked_at IS NOT NULL OR expires_at < ?", now).
		Delete(&RefreshTokenRow{}).Error
}
