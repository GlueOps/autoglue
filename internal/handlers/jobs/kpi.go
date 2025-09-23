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

	if err := db.Model(&models.Job{}).
		Where("status = ?", "running").
		Count(&k.RunningNow).Error; err != nil {
		return k, err
	}

	if err := db.Model(&models.Job{}).
		Where("status IN ?", []string{"queued", "scheduled", "pending"}).
		Where("scheduled_at > ?", now).
		Count(&k.ScheduledFuture).Error; err != nil {
		return k, err
	}

	if err := db.Model(&models.Job{}).
		Where("status IN ?", []string{"queued", "scheduled", "pending"}).
		Where("scheduled_at <= ?", now).
		Count(&k.DueNow).Error; err != nil {
		return k, err
	}

	if err := db.Model(&models.Job{}).
		Where("status = ?", "success").
		Where("updated_at >= ?", dayAgo).
		Count(&k.Succeeded24h).Error; err != nil {
		return k, err
	}

	if err := db.Model(&models.Job{}).
		Where("status = ?", "failed").
		Where("updated_at >= ?", dayAgo).
		Count(&k.Failed24h).Error; err != nil {
		return k, err
	}

	if err := db.Model(&models.Job{}).
		Where("status = ?", "failed").
		Where("retry_count < max_retry").
		Count(&k.Retryable).Error; err != nil {
		return k, err
	}

	return k, nil
}
