package bg

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/dyaksa/archer"
	"github.com/dyaksa/archer/job"
	"github.com/glueops/autoglue/internal/models"
	"github.com/glueops/autoglue/internal/utils"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/ssh"
	"gorm.io/gorm"
)

type ClusterPrepareArgs struct {
	IntervalS int `json:"interval_seconds,omitempty"`
}

type ClusterPrepareFailure struct {
	ClusterID uuid.UUID `json:"cluster_id"`
	Step      string    `json:"step"`
	Reason    string    `json:"reason"`
}

type ClusterPrepareResult struct {
	Status        string                  `json:"status"`
	Processed     int                     `json:"processed"`
	MarkedPending int                     `json:"marked_pending"`
	Failed        int                     `json:"failed"`
	ElapsedMs     int                     `json:"elapsed_ms"`
	FailedIDs     []uuid.UUID             `json:"failed_cluster_ids"`
	Failures      []ClusterPrepareFailure `json:"failures"`
}

// Alias the status constants from models to avoid string drift.
const (
	clusterStatusPrePending   = models.ClusterStatusPrePending
	clusterStatusPending      = models.ClusterStatusPending
	clusterStatusProvisioning = models.ClusterStatusProvisioning
	clusterStatusReady        = models.ClusterStatusReady
	clusterStatusFailed       = models.ClusterStatusFailed
)

func ClusterPrepareWorker(db *gorm.DB, jobs *Jobs) archer.WorkerFn {
	return func(ctx context.Context, j job.Job) (any, error) {
		args := ClusterPrepareArgs{IntervalS: 120}
		jobID := j.ID
		start := time.Now()

		_ = j.ParseArguments(&args)
		if args.IntervalS <= 0 {
			args.IntervalS = 120
		}

		// Load all clusters that are pre_pending; we’ll filter for bastion.ready in memory.
		var clusters []models.Cluster
		if err := db.
			Preload("BastionServer.SshKey").
			Preload("CaptainDomain").
			Preload("ControlPlaneRecordSet").
			Preload("NodePools.Servers.SshKey").
			Where("status = ?", clusterStatusPrePending).
			Find(&clusters).Error; err != nil {
			log.Error().Err(err).Msg("[cluster_prepare] query clusters failed")
			return nil, err
		}

		proc, ok, fail := 0, 0, 0
		var failedIDs []uuid.UUID
		var failures []ClusterPrepareFailure

		perClusterTimeout := 8 * time.Minute

		for i := range clusters {
			c := &clusters[i]
			proc++

			// bastion must exist and be ready
			if c.BastionServer == nil || c.BastionServerID == nil || *c.BastionServerID == uuid.Nil || c.BastionServer.Status != "ready" {
				continue
			}

			clusterLog := log.With().
				Str("job", jobID).
				Str("cluster_id", c.ID.String()).
				Str("cluster_name", c.Name).
				Logger()

			clusterLog.Info().Msg("[cluster_prepare] starting")

			if err := validateClusterForPrepare(c); err != nil {
				fail++
				failedIDs = append(failedIDs, c.ID)
				failures = append(failures, ClusterPrepareFailure{
					ClusterID: c.ID,
					Step:      "validate",
					Reason:    err.Error(),
				})
				clusterLog.Error().Err(err).Msg("[cluster_prepare] validation failed")
				_ = setClusterStatus(db, c.ID, clusterStatusFailed, err.Error())
				continue
			}

			allServers := flattenClusterServers(c)
			keyPayloads, sshConfig, err := buildSSHAssetsForCluster(db, c, allServers)
			if err != nil {
				fail++
				failedIDs = append(failedIDs, c.ID)
				failures = append(failures, ClusterPrepareFailure{
					ClusterID: c.ID,
					Step:      "build_ssh_assets",
					Reason:    err.Error(),
				})
				clusterLog.Error().Err(err).Msg("[cluster_prepare] build ssh assets failed")
				_ = setClusterStatus(db, c.ID, clusterStatusFailed, err.Error())
				continue
			}

			payloadJSON, err := json.MarshalIndent(c, "", "  ")
			if err != nil {
				fail++
				failedIDs = append(failedIDs, c.ID)
				failures = append(failures, ClusterPrepareFailure{
					ClusterID: c.ID,
					Step:      "marshal_payload",
					Reason:    err.Error(),
				})
				clusterLog.Error().Err(err).Msg("[cluster_prepare] json marshal failed")
				_ = setClusterStatus(db, c.ID, clusterStatusFailed, err.Error())
				continue
			}

			runCtx, cancel := context.WithTimeout(ctx, perClusterTimeout)
			err = pushAssetsToBastion(runCtx, db, c, sshConfig, keyPayloads, payloadJSON)
			cancel()

			if err != nil {
				fail++
				failedIDs = append(failedIDs, c.ID)
				failures = append(failures, ClusterPrepareFailure{
					ClusterID: c.ID,
					Step:      "ssh_push",
					Reason:    err.Error(),
				})
				clusterLog.Error().Err(err).Msg("[cluster_prepare] failed to push assets to bastion")
				_ = setClusterStatus(db, c.ID, clusterStatusFailed, err.Error())
				continue
			}

			if err := setClusterStatus(db, c.ID, clusterStatusPending, ""); err != nil {
				fail++
				failedIDs = append(failedIDs, c.ID)
				failures = append(failures, ClusterPrepareFailure{
					ClusterID: c.ID,
					Step:      "set_pending",
					Reason:    err.Error(),
				})
				clusterLog.Error().Err(err).Msg("[cluster_prepare] failed to mark cluster pending")
				continue
			}

			ok++
			clusterLog.Info().Msg("[cluster_prepare] cluster marked pending")
		}

		res := ClusterPrepareResult{
			Status:        "ok",
			Processed:     proc,
			MarkedPending: ok,
			Failed:        fail,
			ElapsedMs:     int(time.Since(start).Milliseconds()),
			FailedIDs:     failedIDs,
			Failures:      failures,
		}

		log.Info().
			Int("processed", proc).
			Int("pending", ok).
			Int("failed", fail).
			Msg("[cluster_prepare] reconcile tick ok")

		next := time.Now().Add(time.Duration(args.IntervalS) * time.Second)
		_, _ = jobs.Enqueue(
			ctx,
			uuid.NewString(),
			"prepare_cluster",
			args,
			archer.WithScheduleTime(next),
			archer.WithMaxRetries(1),
		)
		return res, nil
	}
}

// ---------- helpers ----------

func validateClusterForPrepare(c *models.Cluster) error {
	if c.BastionServer == nil || c.BastionServerID == nil || *c.BastionServerID == uuid.Nil {
		return fmt.Errorf("missing bastion server")
	}
	if c.BastionServer.Status != "ready" {
		return fmt.Errorf("bastion server not ready (status=%s)", c.BastionServer.Status)
	}

	// CaptainDomain is a value type; presence is via *ID
	if c.CaptainDomainID == nil || *c.CaptainDomainID == uuid.Nil {
		return fmt.Errorf("missing captain domain for cluster")
	}

	// ControlPlaneRecordSet is a pointer; presence is via *ID + non-nil struct
	if c.ControlPlaneRecordSetID == nil || *c.ControlPlaneRecordSetID == uuid.Nil || c.ControlPlaneRecordSet == nil {
		return fmt.Errorf("missing control_plane_record_set for cluster")
	}

	if len(c.NodePools) == 0 {
		return fmt.Errorf("cluster has no node pools")
	}

	hasServer := false
	for i := range c.NodePools {
		if len(c.NodePools[i].Servers) > 0 {
			hasServer = true
			break
		}
	}
	if !hasServer {
		return fmt.Errorf("cluster has no servers attached to node pools")
	}

	return nil
}

func flattenClusterServers(c *models.Cluster) []*models.Server {
	var out []*models.Server
	for i := range c.NodePools {
		for j := range c.NodePools[i].Servers {
			s := &c.NodePools[i].Servers[j]
			out = append(out, s)
		}
	}
	return out
}

type keyPayload struct {
	FileName      string
	PrivateKeyB64 string
}

// build ssh-config for all servers + decrypt keys.
// ssh-config is intended to live on the bastion and connect via *private* IPs.
func buildSSHAssetsForCluster(db *gorm.DB, c *models.Cluster, servers []*models.Server) (map[uuid.UUID]keyPayload, string, error) {
	var sb strings.Builder
	keys := make(map[uuid.UUID]keyPayload)

	for _, s := range servers {
		// Defensive checks
		if strings.TrimSpace(s.PrivateIPAddress) == "" {
			return nil, "", fmt.Errorf("server %s missing private ip", s.ID)
		}
		if s.SshKeyID == uuid.Nil {
			return nil, "", fmt.Errorf("server %s missing ssh key relation", s.ID)
		}

		// de-dupe keys: many servers may share the same ssh key
		if _, ok := keys[s.SshKeyID]; !ok {
			priv, err := utils.DecryptForOrg(
				s.OrganizationID,
				s.SshKey.EncryptedPrivateKey,
				s.SshKey.PrivateIV,
				s.SshKey.PrivateTag,
				db,
			)
			if err != nil {
				return nil, "", fmt.Errorf("decrypt key for server %s: %w", s.ID, err)
			}

			fname := fmt.Sprintf("%s.pem", s.SshKeyID.String())
			keys[s.SshKeyID] = keyPayload{
				FileName:      fname,
				PrivateKeyB64: base64.StdEncoding.EncodeToString([]byte(priv)),
			}
		}

		// ssh config entry per server
		keyFile := keys[s.SshKeyID].FileName

		hostAlias := s.Hostname
		if hostAlias == "" {
			hostAlias = s.ID.String()
		}

		sb.WriteString(fmt.Sprintf("Host %s\n", hostAlias))
		sb.WriteString(fmt.Sprintf("  HostName %s\n", s.PrivateIPAddress))
		sb.WriteString(fmt.Sprintf("  User %s\n", s.SSHUser))
		sb.WriteString(fmt.Sprintf("  IdentityFile ~/.ssh/autoglue/keys/%s\n", keyFile))
		sb.WriteString("  IdentitiesOnly yes\n")
		sb.WriteString("  StrictHostKeyChecking accept-new\n\n")
	}

	return keys, sb.String(), nil
}

func pushAssetsToBastion(
	ctx context.Context,
	db *gorm.DB,
	c *models.Cluster,
	sshConfig string,
	keyPayloads map[uuid.UUID]keyPayload,
	payloadJSON []byte,
) error {
	bastion := c.BastionServer
	if bastion == nil {
		return fmt.Errorf("bastion server is nil")
	}

	if bastion.PublicIPAddress == nil || strings.TrimSpace(*bastion.PublicIPAddress) == "" {
		return fmt.Errorf("bastion server missing public ip")
	}

	privKey, err := utils.DecryptForOrg(
		bastion.OrganizationID,
		bastion.SshKey.EncryptedPrivateKey,
		bastion.SshKey.PrivateIV,
		bastion.SshKey.PrivateTag,
		db,
	)
	if err != nil {
		return fmt.Errorf("decrypt bastion key: %w", err)
	}

	signer, err := ssh.ParsePrivateKey([]byte(privKey))
	if err != nil {
		return fmt.Errorf("parse bastion private key: %w", err)
	}

	hkcb := makeDBHostKeyCallback(db, bastion)

	config := &ssh.ClientConfig{
		User:            bastion.SSHUser,
		Auth:            []ssh.AuthMethod{ssh.PublicKeys(signer)},
		HostKeyCallback: hkcb,
		Timeout:         30 * time.Second,
	}

	host := net.JoinHostPort(*bastion.PublicIPAddress, "22")

	dialer := &net.Dialer{}
	conn, err := dialer.DialContext(ctx, "tcp", host)
	if err != nil {
		return fmt.Errorf("dial bastion: %w", err)
	}
	defer conn.Close()

	cconn, chans, reqs, err := ssh.NewClientConn(conn, host, config)
	if err != nil {
		return fmt.Errorf("ssh handshake bastion: %w", err)
	}
	client := ssh.NewClient(cconn, chans, reqs)
	defer client.Close()

	sess, err := client.NewSession()
	if err != nil {
		return fmt.Errorf("ssh session: %w", err)
	}
	defer sess.Close()

	// build one shot script to:
	// - mkdir ~/.ssh/autoglue/keys
	// - write cluster-specific ssh-config
	// - write all private keys
	// - write payload.json
	clusterDir := fmt.Sprintf("$HOME/autoglue/clusters/%s", c.ID.String())
	configPath := fmt.Sprintf("$HOME/.ssh/autoglue/cluster-%s.config", c.ID.String())

	var script bytes.Buffer

	script.WriteString("set -euo pipefail\n")
	script.WriteString("mkdir -p \"$HOME/.ssh/autoglue/keys\"\n")
	script.WriteString("mkdir -p " + clusterDir + "\n")
	script.WriteString("chmod 700 \"$HOME/.ssh\" || true\n")

	// ssh-config
	script.WriteString("cat > " + configPath + " <<'EOF_CFG'\n")
	script.WriteString(sshConfig)
	script.WriteString("EOF_CFG\n")
	script.WriteString("chmod 600 " + configPath + "\n")

	// keys
	for id, kp := range keyPayloads {
		tag := "KEY_" + id.String()
		target := fmt.Sprintf("$HOME/.ssh/autoglue/keys/%s", kp.FileName)

		script.WriteString("cat <<'" + tag + "' | base64 -d > " + target + "\n")
		script.WriteString(kp.PrivateKeyB64 + "\n")
		script.WriteString(tag + "\n")
		script.WriteString("chmod 600 " + target + "\n")
	}

	// payload.json
	payloadPath := clusterDir + "/payload.json"
	script.WriteString("cat > " + payloadPath + " <<'EOF_PAYLOAD'\n")
	script.Write(payloadJSON)
	script.WriteString("\nEOF_PAYLOAD\n")
	script.WriteString("chmod 600 " + payloadPath + "\n")

	// If you later want to always include cluster configs automatically, you can
	// optionally manage ~/.ssh/config here (kept simple for now).

	sess.Stdin = strings.NewReader(script.String())
	out, runErr := sess.CombinedOutput("bash -s")

	if runErr != nil {
		return wrapSSHError(runErr, string(out))
	}
	return nil
}

func setClusterStatus(db *gorm.DB, id uuid.UUID, status, lastError string) error {
	updates := map[string]any{
		"status":     status,
		"updated_at": time.Now(),
	}
	if lastError != "" {
		updates["last_error"] = lastError
	}
	return db.Model(&models.Cluster{}).
		Where("id = ?", id).
		Updates(updates).Error
}

// runMakeOnBastion runs `make <target>` from the cluster's directory on the bastion.
func runMakeOnBastion(
	ctx context.Context,
	db *gorm.DB,
	c *models.Cluster,
	target string,
) (string, error) {
	bastion := c.BastionServer
	if bastion == nil {
		return "", fmt.Errorf("bastion server is nil")
	}

	if bastion.PublicIPAddress == nil || strings.TrimSpace(*bastion.PublicIPAddress) == "" {
		return "", fmt.Errorf("bastion server missing public ip")
	}

	privKey, err := utils.DecryptForOrg(
		bastion.OrganizationID,
		bastion.SshKey.EncryptedPrivateKey,
		bastion.SshKey.PrivateIV,
		bastion.SshKey.PrivateTag,
		db,
	)
	if err != nil {
		return "", fmt.Errorf("decrypt bastion key: %w", err)
	}

	signer, err := ssh.ParsePrivateKey([]byte(privKey))
	if err != nil {
		return "", fmt.Errorf("parse bastion private key: %w", err)
	}

	hkcb := makeDBHostKeyCallback(db, bastion)

	config := &ssh.ClientConfig{
		User:            bastion.SSHUser,
		Auth:            []ssh.AuthMethod{ssh.PublicKeys(signer)},
		HostKeyCallback: hkcb,
		Timeout:         30 * time.Second,
	}

	host := net.JoinHostPort(*bastion.PublicIPAddress, "22")

	dialer := &net.Dialer{}
	conn, err := dialer.DialContext(ctx, "tcp", host)
	if err != nil {
		return "", fmt.Errorf("dial bastion: %w", err)
	}
	defer conn.Close()

	cconn, chans, reqs, err := ssh.NewClientConn(conn, host, config)
	if err != nil {
		return "", fmt.Errorf("ssh handshake bastion: %w", err)
	}
	client := ssh.NewClient(cconn, chans, reqs)
	defer client.Close()

	sess, err := client.NewSession()
	if err != nil {
		return "", fmt.Errorf("ssh session: %w", err)
	}
	defer sess.Close()

	clusterDir := fmt.Sprintf("$HOME/autoglue/clusters/%s", c.ID.String())
	sshDir := fmt.Sprintf("$HOME/.ssh/autoglue")

	cmd := fmt.Sprintf("cd %s && docker run -it -v %s:/root/.ssh -v ./payload.json:/opt/gluekube/platform.json %s:%s make %s", clusterDir, sshDir, c.DockerImage, c.DockerTag, target)

	out, runErr := sess.CombinedOutput(cmd)
	if runErr != nil {
		return string(out), wrapSSHError(runErr, string(out))
	}
	return string(out), nil
}
