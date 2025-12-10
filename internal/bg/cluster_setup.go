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

type ClusterSetupArgs struct {
	IntervalS int `json:"interval_seconds,omitempty"`
}

type ClusterSetupResult struct {
	Status        string      `json:"status"`
	Processed     int         `json:"processed"`
	Provisioning  int         `json:"provisioning"`
	Failed        int         `json:"failed"`
	ElapsedMs     int         `json:"elapsed_ms"`
	FailedCluster []uuid.UUID `json:"failed_cluster_ids"`
}

func ClusterSetupWorker(db *gorm.DB, jobs *Jobs) archer.WorkerFn {
	return func(ctx context.Context, j job.Job) (any, error) {
		args := ClusterSetupArgs{IntervalS: 120}
		jobID := j.ID
		start := time.Now()

		_ = j.ParseArguments(&args)
		if args.IntervalS <= 0 {
			args.IntervalS = 120
		}

		var clusters []models.Cluster
		if err := db.
			Preload("BastionServer.SshKey").
			Where("status = ?", clusterStatusPending).
			Find(&clusters).Error; err != nil {
			log.Error().Err(err).Msg("[cluster_setup] query clusters failed")
			return nil, err
		}

		proc, prov, failCount := 0, 0, 0
		var failedIDs []uuid.UUID

		perClusterTimeout := 30 * time.Minute

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

			logger.Info().Msg("[cluster_setup] running make setup")

			runCtx, cancel := context.WithTimeout(ctx, perClusterTimeout)
			out, err := runMakeOnBastion(runCtx, db, c, "ping-servers")
			cancel()

			if err != nil {
				failCount++
				failedIDs = append(failedIDs, c.ID)
				logger.Error().Err(err).Str("output", out).Msg("[cluster_setup] make setup failed")
				_ = setClusterStatus(db, c.ID, clusterStatusFailed, fmt.Sprintf("make setup: %v", err))
				continue
			}

			if err := setClusterStatus(db, c.ID, clusterStatusProvisioning, ""); err != nil {
				failCount++
				failedIDs = append(failedIDs, c.ID)
				logger.Error().Err(err).Msg("[cluster_setup] failed to mark cluster provisioning")
				continue
			}

			prov++
			logger.Info().Msg("[cluster_setup] cluster moved to provisioning")
		}

		res := ClusterSetupResult{
			Status:        "ok",
			Processed:     proc,
			Provisioning:  prov,
			Failed:        failCount,
			ElapsedMs:     int(time.Since(start).Milliseconds()),
			FailedCluster: failedIDs,
		}

		log.Info().
			Int("processed", proc).
			Int("provisioning", prov).
			Int("failed", failCount).
			Msg("[cluster_setup] reconcile tick ok")

		// self-reschedule
		next := time.Now().Add(time.Duration(args.IntervalS) * time.Second)
		_, _ = jobs.Enqueue(
			ctx,
			uuid.NewString(),
			"cluster_setup",
			args,
			archer.WithScheduleTime(next),
			archer.WithMaxRetries(1),
		)
		return res, nil
	}
}
