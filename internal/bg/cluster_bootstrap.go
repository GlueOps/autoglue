package bg

import (
	"context"
	"fmt"
	"time"

	"github.com/dyaksa/archer"
	"github.com/dyaksa/archer/job"
	"github.com/glueops/autoglue/internal/models"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

type ClusterBootstrapArgs struct {
	IntervalS int `json:"interval_seconds,omitempty"`
}

type ClusterBootstrapResult struct {
	Status    string      `json:"status"`
	Processed int         `json:"processed"`
	Ready     int         `json:"ready"`
	Failed    int         `json:"failed"`
	ElapsedMs int         `json:"elapsed_ms"`
	FailedIDs []uuid.UUID `json:"failed_cluster_ids"`
}

func ClusterBootstrapWorker(db *gorm.DB, jobs *Jobs) archer.WorkerFn {
	return func(ctx context.Context, j job.Job) (any, error) {
		args := ClusterBootstrapArgs{IntervalS: 120}
		jobID := j.ID
		start := time.Now()

		_ = j.ParseArguments(&args)
		if args.IntervalS <= 0 {
			args.IntervalS = 120
		}

		var clusters []models.Cluster
		if err := db.
			Preload("BastionServer.SshKey").
			Where("status = ?", clusterStatusProvisioning).
			Find(&clusters).Error; err != nil {
			log.Error().Err(err).Msg("[cluster_bootstrap] query clusters failed")
			return nil, err
		}

		proc, ready, failCount := 0, 0, 0
		var failedIDs []uuid.UUID

		perClusterTimeout := 60 * time.Minute

		for i := range clusters {
			c := &clusters[i]
			proc++

			if c.BastionServer.ID == uuid.Nil || c.BastionServer.Status != "ready" {
				continue
			}

			logger := log.With().
				Str("job", jobID).
				Str("cluster_id", c.ID.String()).
				Str("cluster_name", c.Name).
				Logger()

			logger.Info().Msg("[cluster_bootstrap] running make bootstrap")

			runCtx, cancel := context.WithTimeout(ctx, perClusterTimeout)
			out, err := runMakeOnBastion(runCtx, db, c, "setup")
			cancel()

			if err != nil {
				failCount++
				failedIDs = append(failedIDs, c.ID)
				logger.Error().Err(err).Str("output", out).Msg("[cluster_bootstrap] make setup failed")
				_ = setClusterStatus(db, c.ID, clusterStatusFailed, fmt.Sprintf("make setup: %v", err))
				continue
			}

			// you can choose a different terminal status here if you like
			if err := setClusterStatus(db, c.ID, clusterStatusReady, ""); err != nil {
				failCount++
				failedIDs = append(failedIDs, c.ID)
				logger.Error().Err(err).Msg("[cluster_bootstrap] failed to mark cluster ready")
				continue
			}

			ready++
			logger.Info().Msg("[cluster_bootstrap] cluster marked ready")
		}

		res := ClusterBootstrapResult{
			Status:    "ok",
			Processed: proc,
			Ready:     ready,
			Failed:    failCount,
			ElapsedMs: int(time.Since(start).Milliseconds()),
			FailedIDs: failedIDs,
		}

		log.Info().
			Int("processed", proc).
			Int("ready", ready).
			Int("failed", failCount).
			Msg("[cluster_bootstrap] reconcile tick ok")

		// self-reschedule
		next := time.Now().Add(time.Duration(args.IntervalS) * time.Second)
		_, _ = jobs.Enqueue(
			ctx,
			uuid.NewString(),
			"cluster_bootstrap",
			args,
			archer.WithScheduleTime(next),
			archer.WithMaxRetries(1),
		)
		return res, nil
	}
}
