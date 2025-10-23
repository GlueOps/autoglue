package jobs

import (
	"time"

	"github.com/glueops/autoglue/internal/db/models"
	"gorm.io/gorm"
)

type QueueRollup struct {
	QueueName       string  `json:"queue_name" example:"emails"`
	Running         int64   `json:"running" example:"2"`
	QueuedDue       int64   `json:"queued_due" example:"7"`
	QueuedFuture    int64   `json:"queued_future" example:"3"`
	Success24h      int64   `json:"success_24h" example:"124"`
	Failed24h       int64   `json:"failed_24h" example:"1"`
	AvgDurationSecs float64 `json:"avg_duration_secs" swaggertype:"number" format:"double" example:"1.234"`
}

func LoadPerQueue(db *gorm.DB) ([]QueueRollup, error) {
	var queues []string
	if err := db.Model(&models.Job{}).Distinct().Pluck("queue_name", &queues).Error; err != nil {
		return nil, err
	}

	now := time.Now()
	dayAgo := now.Add(-24 * time.Hour)
	out := make([]QueueRollup, 0, len(queues))

	for _, q := range queues {
		var rr, qd, qf, s24, f24 int64
		var avgDur *float64

		_ = db.Model(&models.Job{}).
			Where("queue_name = ? AND status = 'running'", q).
			Count(&rr).Error

		_ = db.Model(&models.Job{}).
			Where("queue_name = ? AND status IN ('queued','scheduled','pending') AND scheduled_at <= ?", q, now).
			Count(&qd).Error

		_ = db.Model(&models.Job{}).
			Where("queue_name = ? AND status IN ('queued','scheduled','pending') AND scheduled_at >  ?", q, now).
			Count(&qf).Error

		// Sum result.ready / result.failed over successes in last 24h
		_ = db.Model(&models.Job{}).
			Select("COALESCE(SUM((result->>'ready')::int), 0)").
			Where("queue_name = ? AND status = 'success' AND updated_at >= ?", q, dayAgo).
			Scan(&s24).Error

		_ = db.Model(&models.Job{}).
			Select("COALESCE(SUM((result->>'failed')::int), 0)").
			Where("queue_name = ? AND status = 'success' AND updated_at >= ?", q, dayAgo).
			Scan(&f24).Error

		_ = db.
			Model(&models.Job{}).
			Select("AVG(EXTRACT(EPOCH FROM (updated_at - started_at)))").
			Where("queue_name = ? AND status = 'success' AND started_at IS NOT NULL AND updated_at >= ?", q, dayAgo).
			Scan(&avgDur).Error

		out = append(out, QueueRollup{
			QueueName:       q,
			Running:         rr,
			QueuedDue:       qd,
			QueuedFuture:    qf,
			Success24h:      s24,
			Failed24h:       f24,
			AvgDurationSecs: coalesceF64(avgDur, 0),
		})
	}
	return out, nil
}

func coalesceF64(p *float64, d float64) float64 {
	if p == nil {
		return d
	}
	return *p
}
