package bg

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/glueops/autoglue/internal/mapper"
	"github.com/glueops/autoglue/internal/models"
	"github.com/glueops/autoglue/internal/utils"
	"github.com/google/uuid"
	"github.com/riverqueue/river"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ClusterActionArgs struct {
	// RunID identifies the models.ClusterRun row this job reports into. Under
	// archer this rode on the job ID itself; River assigns its own int64 IDs, so
	// it travels in the args.
	RunID      uuid.UUID `json:"run_id"`
	OrgID      uuid.UUID `json:"org_id"`
	ClusterID  uuid.UUID `json:"cluster_id"`
	Action     string    `json:"action"`
	MakeTarget string    `json:"make_target"`
}

func (ClusterActionArgs) Kind() string { return "cluster_action" }

// ClusterActionResult is recorded on the job row via river.RecordOutput, which
// is River's replacement for the jobs.result column archer used to persist.
type ClusterActionResult struct {
	Status    string `json:"status"`
	Action    string `json:"action"`
	ClusterID string `json:"cluster_id"`
	ElapsedMs int    `json:"elapsed_ms"`
}

func (ClusterActionArgs) InsertOpts() river.InsertOpts {
	// A cluster action is not safely repeatable: it drives make targets against
	// a live bastion. Archer ran these with MaxRetries(0); River counts the
	// first attempt, so 1 is the equivalent.
	return river.InsertOpts{Queue: QueueClusters, MaxAttempts: 1}
}

type ClusterActionWorker struct {
	river.WorkerDefaults[ClusterActionArgs]
	db      *gorm.DB
	baseURL string
}

// Timeout covers the full prepare -> ping-servers -> bootstrap sequence, whose
// inner steps carry their own shorter budgets.
func (w *ClusterActionWorker) Timeout(*river.Job[ClusterActionArgs]) time.Duration {
	return 168 * time.Hour
}

func (w *ClusterActionWorker) Work(ctx context.Context, j *river.Job[ClusterActionArgs]) error {
	db, baseURL := w.db, w.baseURL
	start := time.Now()
	args := j.Args
	runID := args.RunID

	// Everything this action prints streams into job_logs under the ClusterRun,
	// so the cluster page can tail it live instead of waiting for the run to
	// finish (and, on success, losing the output entirely).
	sink := NewLogSink(db, j.ID, args.OrgID, models.JobLogSubjectClusterRun, runID)
	defer func() { _ = sink.Close() }()

	{
		updateRun := func(status string, errMsg string) {
			updates := map[string]any{
				"status": status,
				"error":  errMsg,
			}
			if status == "succeeded" || status == "failed" {
				updates["finished_at"] = time.Now().UTC().Format(time.RFC3339)
			}
			db.Model(&models.ClusterRun{}).Where("id = ?", runID).Updates(updates)
		}

		updateRun("running", "")

		logger := log.With().
			Int64("job", j.ID).
			Str("cluster_id", args.ClusterID.String()).
			Str("action", args.Action).
			Logger()

		// Atomically claim the cluster using SELECT FOR UPDATE SKIP LOCKED.
		// This prevents two concurrent cluster_action workers from processing the
		// same cluster simultaneously (e.g. duplicate API calls or a retried job).
		// The status guard also stops jobs from re-entering an already in-progress run.
		var claimedIDs []uuid.UUID
		if err := db.Transaction(func(tx *gorm.DB) error {
			inProgressStatuses := []string{clusterStatusBootstrapping, clusterStatusProvisioning}
			if err := tx.Model(&models.Cluster{}).
				Where("id = ? AND organization_id = ? AND status NOT IN ?",
					args.ClusterID, args.OrgID, inProgressStatuses).
				Clauses(clause.Locking{Strength: "UPDATE", Options: "SKIP LOCKED"}).
				Pluck("id", &claimedIDs).Error; err != nil {
				return err
			}
			if len(claimedIDs) == 0 {
				return nil
			}
			return tx.Model(&models.Cluster{}).
				Where("id = ?", args.ClusterID).
				Updates(map[string]any{
					"status":     clusterStatusBootstrapping,
					"updated_at": time.Now(),
				}).Error
		}); err != nil {
			updateRun("failed", err.Error())
			return fmt.Errorf("claim cluster: %w", err)
		}

		if len(claimedIDs) == 0 {
			msg := fmt.Sprintf("cluster %s is already being processed by another worker", args.ClusterID)
			logger.Warn().Msg(msg)
			//updateRun("failed", msg)
			//return nil, errors.New(msg)
		}

		var c models.Cluster
		if err := db.
			Preload("BastionServer.SshKey").
			Preload("CaptainDomain").
			Preload("ControlPlaneRecordSet").
			Preload("AppsLoadBalancer").
			Preload("GlueOpsLoadBalancer").
			Preload("NodePools").
			Preload("NodePools.Labels").
			Preload("NodePools.Annotations").
			Preload("NodePools.Taints").
			Preload("NodePools.Servers.SshKey").
			Preload("Metadata").
			Where("id = ? AND organization_id = ?", args.ClusterID, args.OrgID).
			First(&c).Error; err != nil {
			updateRun("failed", fmt.Errorf("load cluster: %w", err).Error())
			return fmt.Errorf("load cluster: %w", err)
		}

		// ---- Step 1: Prepare (mostly lifted from ClusterPrepareWorker)
		// Status was already set to bootstrapping atomically during the claim transaction.
		c.Status = clusterStatusBootstrapping

		if err := validateClusterForPrepare(&c); err != nil {
			_ = setClusterStatus(db, c.ID, clusterStatusFailed, err.Error())
			updateRun("failed", err.Error())
			return fmt.Errorf("validate: %w", err)
		}

		allServers := flattenClusterServers(&c)
		keyPayloads, sshConfig, err := buildSSHAssetsForCluster(db, &c, allServers)
		if err != nil {
			_ = setClusterStatus(db, c.ID, clusterStatusFailed, err.Error())
			updateRun("failed", err.Error())
			return fmt.Errorf("build ssh assets: %w", err)
		}

		dtoCluster := mapper.ClusterToDTO(c)

		if c.EncryptedKubeconfig != "" && c.KubeIV != "" && c.KubeTag != "" {
			kubeconfig, err := utils.DecryptForOrg(
				c.OrganizationID,
				c.EncryptedKubeconfig,
				c.KubeIV,
				c.KubeTag,
				db,
			)
			if err != nil {
				_ = setClusterStatus(db, c.ID, clusterStatusFailed, err.Error())
				return fmt.Errorf("decrypt kubeconfig: %w", err)
			}
			dtoCluster.Kubeconfig = &kubeconfig
		}

		orgKey, orgSecret, err := findOrCreateClusterAutomationKey(db, c.OrganizationID, c.ID, 24*time.Hour)
		if err != nil {
			_ = setClusterStatus(db, c.ID, clusterStatusFailed, err.Error())
			updateRun("failed", err.Error())
			return fmt.Errorf("org key: %w", err)
		}
		dtoCluster.OrgKey = &orgKey
		dtoCluster.OrgSecret = &orgSecret
		dtoCluster.BaseURL = baseURL

		payloadJSON, err := json.MarshalIndent(dtoCluster, "", "  ")
		if err != nil {
			_ = setClusterStatus(db, c.ID, clusterStatusFailed, err.Error())
			updateRun("failed", err.Error())
			return fmt.Errorf("marshal payload: %w", err)
		}

		{
			runCtx, cancel := context.WithTimeout(ctx, 8*time.Minute)
			err := pushAssetsToBastion(runCtx, db, &c, sshConfig, keyPayloads, payloadJSON)
			cancel()
			if err != nil {
				_ = setClusterStatus(db, c.ID, clusterStatusFailed, err.Error())
				updateRun("failed", err.Error())
				return fmt.Errorf("push assets: %w", err)
			}
		}

		if err := setClusterStatus(db, c.ID, clusterStatusPending, ""); err != nil {
			updateRun("failed", err.Error())
			return fmt.Errorf("mark pending: %w", err)
		}
		c.Status = clusterStatusPending

		// ---- Step 2: Setup (ping-servers)
		{
			runCtx, cancel := context.WithTimeout(ctx, 30*time.Minute)
			sink.System("running make ping-servers")
			out, err := runMakeOnBastion(runCtx, db, &c, "ping-servers", sink)
			cancel()
			if err != nil {
				logger.Error().Err(err).Str("output", out).Msg("ping-servers failed")
				_ = setClusterStatus(db, c.ID, clusterStatusFailed, fmt.Sprintf("make ping-servers: %v", err))
				updateRun("failed", err.Error())
				return fmt.Errorf("ping-servers: %w", err)
			}
		}

		if err := setClusterStatus(db, c.ID, clusterStatusProvisioning, ""); err != nil {
			updateRun("failed", err.Error())
			return fmt.Errorf("mark provisioning: %w", err)
		}
		c.Status = clusterStatusProvisioning

		// ---- Step 3: Bootstrap (parameterized target)
		{
			runCtx, cancel := context.WithTimeout(ctx, 60*time.Minute)
			sink.System("running make " + args.MakeTarget)
			out, err := runMakeOnBastion(runCtx, db, &c, args.MakeTarget, sink)
			cancel()
			if err != nil {
				logger.Error().Err(err).Str("output", out).Msg("bootstrap target failed")
				_ = setClusterStatus(db, c.ID, clusterStatusFailed, fmt.Sprintf("make %s: %v", args.MakeTarget, err))
				updateRun("failed", err.Error())
				return fmt.Errorf("make %s: %w", args.MakeTarget, err)
			}
		}

		if err := setClusterStatus(db, c.ID, clusterStatusReady, ""); err != nil {
			updateRun("failed", err.Error())
			return fmt.Errorf("mark ready: %w", err)
		}

		updateRun("succeeded", "")
		sink.System("completed")

		if err := river.RecordOutput(ctx, ClusterActionResult{
			Status:    "ok",
			Action:    args.Action,
			ClusterID: c.ID.String(),
			ElapsedMs: int(time.Since(start).Milliseconds()),
		}); err != nil {
			logger.Warn().Err(err).Msg("[cluster_action] could not record output")
		}

		logger.Info().
			Str("cluster_id", c.ID.String()).
			Int64("elapsed_ms", time.Since(start).Milliseconds()).
			Msg("[cluster_action] completed")
	}
	return nil
}
