package jobs

import (
	"context"

	"github.com/glueops/autoglue/internal/db/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type claimedJob struct {
	models.Job
}

func ClaimBatch(ctx context.Context, db *gorm.DB, workerID string, limit int) ([]models.Job, error) {
	sql := `
WITH cte AS (
  SELECT id
  FROM jobs
  WHERE status = 'queued'
    AND scheduled_at <= NOW()
  ORDER BY priority DESC, scheduled_at ASC
  FOR UPDATE SKIP LOCKED
  LIMIT @limit
)
UPDATE jobs j
SET status     = 'running',
    locked_at  = NOW(),
    locked_by  = @worker,
    started_at = NOW()
FROM cte
WHERE j.id = cte.id
RETURNING j.*;
`
	rows, err := db.WithContext(ctx).Raw(sql, gorm.Named("limit", limit), gorm.Named("worker", workerID)).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []models.Job
	for rows.Next() {
		var j models.Job
		if err := db.ScanRows(rows, &j); err != nil {
			return nil, err
		}
		out = append(out, j)
	}
	return out, nil
}

func FinishSuccess(ctx context.Context, db *gorm.DB, id uuid.UUID) error {
	q := db.DB
}

func FinishRetry(ctx context.Context, db *gorm.DB, id uuid.UUID, attempts, maxAttempts int, lastErr string) error {
}

func FinishFailed(ctx context.Context, db *gorm.DB, id uuid.UUID, lastErr string) error {
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n]
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
