package jobs

import (
	"time"

	"github.com/glueops/autoglue/internal/db/models"
	"gorm.io/gorm"
)

type KPI struct {
	RunningNow      int64
	DueNow          int64
	ScheduledFuture int64
	Succeeded24h    int64
	Failed24h       int64
	Retryable       int64
}

func LoadKPI(db *gorm.DB) (KPI, error) {
	var k KPI
	now := time.Now()
	dayAgo := now.Add(-24 * time.Hour)

	// Running now
	if err := db.Model(&models.Job{}).
		Where("status = ?", "running").
		Count(&k.RunningNow).Error; err != nil {
		return k, err
	}

	// Scheduled in the future
	if err := db.Model(&models.Job{}).
		Where("status IN ?", []string{"queued", "scheduled", "pending"}).
		Where("scheduled_at > ?", now).
		Count(&k.ScheduledFuture).Error; err != nil {
		return k, err
	}

	// Due now (queued/scheduled/pending with scheduled_at <= now)
	if err := db.Model(&models.Job{}).
		Where("status IN ?", []string{"queued", "scheduled", "pending"}).
		Where("scheduled_at <= ?", now).
		Count(&k.DueNow).Error; err != nil {
		return k, err
	}

	// Sum of 'ready' over successful jobs in last 24h
	if err := db.Model(&models.Job{}).
		Select("COALESCE(SUM((result->>'ready')::int), 0)").
		Where("status = 'success' AND updated_at >= ?", dayAgo).
		Scan(&k.Succeeded24h).Error; err != nil {
		return k, err
	}

	// Sum of 'failed' over successful jobs in last 24h (failures within a “success” run)
	if err := db.Model(&models.Job{}).
		Select("COALESCE(SUM((result->>'failed')::int), 0)").
		Where("status = 'success' AND updated_at >= ?", dayAgo).
		Scan(&k.Failed24h).Error; err != nil {
		return k, err
	}

	// Retryable failed job rows (same as before)
	if err := db.Model(&models.Job{}).
		Where("status = 'failed'").
		Where("retry_count < max_retry").
		Count(&k.Retryable).Error; err != nil {
		return k, err
	}

	return k, nil
}
