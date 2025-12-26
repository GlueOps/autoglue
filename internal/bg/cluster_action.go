package bg

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/dyaksa/archer"
	"github.com/dyaksa/archer/job"
	"github.com/glueops/autoglue/internal/mapper"
	"github.com/glueops/autoglue/internal/models"
	"github.com/glueops/autoglue/internal/utils"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

type ClusterActionArgs struct {
	OrgID      uuid.UUID `json:"org_id"`
	ClusterID  uuid.UUID `json:"cluster_id"`
	Action     string    `json:"action"`
	MakeTarget string    `json:"make_target"`
}

type ClusterActionResult struct {
	Status    string `json:"status"`
	Action    string `json:"action"`
	ClusterID string `json:"cluster_id"`
	ElapsedMs int    `json:"elapsed_ms"`
}

func ClusterActionWorker(db *gorm.DB) archer.WorkerFn {
	return func(ctx context.Context, j job.Job) (any, error) {
		start := time.Now()
		var args ClusterActionArgs
		_ = j.ParseArguments(&args)

		runID, _ := uuid.Parse(j.ID)

		updateRun := func(status string, errMsg string) {
			updates := map[string]any{
				"status": status,
				"error":  errMsg,
			}
			if status == "succeeded" || status == "failed" {
				updates["finised_at"] = time.Now().UTC()
			}
			db.Model(&models.ClusterRun{}).Where("id = ?", runID).Updates(updates)
		}

		updateRun("running", "")

		logger := log.With().
			Str("job", j.ID).
			Str("cluster_id", args.ClusterID.String()).
			Str("action", args.Action).
			Logger()

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
			Where("id = ? AND organization_id = ?", args.ClusterID, args.OrgID).
			First(&c).Error; err != nil {
			updateRun("failed", fmt.Errorf("load cluster: %w", err).Error())
			return nil, fmt.Errorf("load cluster: %w", err)
		}

		// ---- Step 1: Prepare (mostly lifted from ClusterPrepareWorker)
		if err := setClusterStatus(db, c.ID, clusterStatusBootstrapping, ""); err != nil {
			updateRun("failed", err.Error())
			return nil, fmt.Errorf("mark bootstrapping: %w", err)
		}
		c.Status = clusterStatusBootstrapping

		if err := validateClusterForPrepare(&c); err != nil {
			_ = setClusterStatus(db, c.ID, clusterStatusFailed, err.Error())
			updateRun("failed", err.Error())
			return nil, fmt.Errorf("validate: %w", err)
		}

		allServers := flattenClusterServers(&c)
		keyPayloads, sshConfig, err := buildSSHAssetsForCluster(db, &c, allServers)
		if err != nil {
			_ = setClusterStatus(db, c.ID, clusterStatusFailed, err.Error())
			updateRun("failed", err.Error())
			return nil, fmt.Errorf("build ssh assets: %w", err)
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
				return nil, fmt.Errorf("decrypt kubeconfig: %w", err)
			}
			dtoCluster.Kubeconfig = &kubeconfig
		}

		orgKey, orgSecret, err := findOrCreateClusterAutomationKey(db, c.OrganizationID, c.ID, 24*time.Hour)
		if err != nil {
			_ = setClusterStatus(db, c.ID, clusterStatusFailed, err.Error())
			updateRun("failed", err.Error())
			return nil, fmt.Errorf("org key: %w", err)
		}
		dtoCluster.OrgKey = &orgKey
		dtoCluster.OrgSecret = &orgSecret

		payloadJSON, err := json.MarshalIndent(dtoCluster, "", "  ")
		if err != nil {
			_ = setClusterStatus(db, c.ID, clusterStatusFailed, err.Error())
			updateRun("failed", err.Error())
			return nil, fmt.Errorf("marshal payload: %w", err)
		}

		{
			runCtx, cancel := context.WithTimeout(ctx, 8*time.Minute)
			err := pushAssetsToBastion(runCtx, db, &c, sshConfig, keyPayloads, payloadJSON)
			cancel()
			if err != nil {
				_ = setClusterStatus(db, c.ID, clusterStatusFailed, err.Error())
				updateRun("failed", err.Error())
				return nil, fmt.Errorf("push assets: %w", err)
			}
		}

		if err := setClusterStatus(db, c.ID, clusterStatusPending, ""); err != nil {
			updateRun("failed", err.Error())
			return nil, fmt.Errorf("mark pending: %w", err)
		}
		c.Status = clusterStatusPending

		// ---- Step 2: Setup (ping-servers)
		{
			runCtx, cancel := context.WithTimeout(ctx, 30*time.Minute)
			out, err := runMakeOnBastion(runCtx, db, &c, "ping-servers")
			cancel()
			if err != nil {
				logger.Error().Err(err).Str("output", out).Msg("ping-servers failed")
				_ = setClusterStatus(db, c.ID, clusterStatusFailed, fmt.Sprintf("make ping-servers: %v", err))
				updateRun("failed", err.Error())
				return nil, fmt.Errorf("ping-servers: %w", err)
			}
		}

		if err := setClusterStatus(db, c.ID, clusterStatusProvisioning, ""); err != nil {
			updateRun("failed", err.Error())
			return nil, fmt.Errorf("mark provisioning: %w", err)
		}
		c.Status = clusterStatusProvisioning

		// ---- Step 3: Bootstrap (parameterized target)
		{
			runCtx, cancel := context.WithTimeout(ctx, 60*time.Minute)
			out, err := runMakeOnBastion(runCtx, db, &c, args.MakeTarget)
			cancel()
			if err != nil {
				logger.Error().Err(err).Str("output", out).Msg("bootstrap target failed")
				_ = setClusterStatus(db, c.ID, clusterStatusFailed, fmt.Sprintf("make %s: %v", args.MakeTarget, err))
				updateRun("failed", err.Error())
				return nil, fmt.Errorf("make %s: %w", args.MakeTarget, err)
			}
		}

		if err := setClusterStatus(db, c.ID, clusterStatusReady, ""); err != nil {
			updateRun("failed", err.Error())
			return nil, fmt.Errorf("mark ready: %w", err)
		}

		updateRun("succeeded", "")

		return ClusterActionResult{
			Status:    "ok",
			Action:    args.Action,
			ClusterID: c.ID.String(),
			ElapsedMs: int(time.Since(start).Milliseconds()),
		}, nil
	}
}
