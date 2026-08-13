package bg

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"net"
	"strings"
	"time"

	"github.com/glueops/autoglue/internal/auth"
	"github.com/glueops/autoglue/internal/models"
	"github.com/glueops/autoglue/internal/utils"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/ssh"
	"gorm.io/gorm"
)

// Alias the status constants from models to avoid string drift.
const (
	clusterStatusPrePending    = models.ClusterStatusPrePending
	clusterStatusPending       = models.ClusterStatusPending
	clusterStatusProvisioning  = models.ClusterStatusProvisioning
	clusterStatusReady         = models.ClusterStatusReady
	clusterStatusFailed        = models.ClusterStatusFailed
	clusterStatusBootstrapping = models.ClusterStatusBootstrapping
)

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
// runMakeOnBastion runs `make <target>` on the cluster's bastion, streaming the
// combined output to sink as it arrives. sink may be nil, in which case output
// is still captured for the returned tail but nothing is persisted.
//
// The returned string is the trailing logMaxTailBytes of output, not the whole
// run: the full transcript lives in job_logs when a sink is supplied.
func runMakeOnBastion(
	ctx context.Context,
	db *gorm.DB,
	c *models.Cluster,
	runID uuid.UUID,
	target string,
	sink io.Writer,
) (string, error) {
	logger := log.With().
		Str("cluster_id", c.ID.String()).
		Str("cluster_name", c.Name).
		Logger()

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
	sshDir := fmt.Sprintf("$HOME/.ssh")

	// Labels rather than --name. --sig-proxy=false means the container now
	// outlives a dropped SSH connection, so it has to be findable afterwards to
	// be reaped or read — but a name is a uniqueness constraint, and Docker
	// holds a name until the container is *removed*, not until it exits. With no
	// --rm, any per-cluster name would collide with its own previous run, and
	// even a per-run name would collide between the two steps of a single run
	// (ping-servers, then the target). Labels carry the same identity with none
	// of that:
	//
	//   docker ps -a --filter label=autoglue.cluster=<id>
	//   docker logs $(docker ps -aq --filter label=autoglue.run=<id>)
	//
	// Both values are UUIDs, so nothing here widens the existing interpolation
	// surface the way a free-text label would.
	labels := fmt.Sprintf("--label autoglue.cluster=%s --label autoglue.run=%s", c.ID.String(), runID.String())

	cmd := fmt.Sprintf("cd %s && docker run --sig-proxy=false %s -v %s:/root/.ssh -v ./payload.json:/opt/gluekube/platform.json %s:%s make %s", clusterDir, labels, sshDir, c.DockerImage, c.DockerTag, target)

	logger.Info().
		Str("cmd", cmd).
		Msg("[runMakeOnBastion] executing remote command")

	tail := &tailBuffer{max: logMaxTailBytes}
	var w io.Writer = tail
	if sink != nil {
		w = io.MultiWriter(tail, sink)
	}

	if runErr := runSSHStreaming(sess, cmd, w); runErr != nil {
		return tail.String(), wrapSSHError(runErr, tail.String())
	}
	return tail.String(), nil
}

func randomB64URL(n int) (string, error) {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

func findOrCreateClusterAutomationKey(
	db *gorm.DB,
	orgID uuid.UUID,
	clusterID uuid.UUID,
	ttl time.Duration,
) (orgKey string, orgSecret string, err error) {
	now := time.Now()
	name := fmt.Sprintf("cluster-%s-bastion", clusterID.String())

	// 1) Delete any existing ephemeral cluster-bastion key for this org+cluster
	if err := db.Where(
		"org_id = ? AND scope = ? AND purpose = ? AND cluster_id = ? AND is_ephemeral = ?",
		orgID, "org", "cluster_bastion", clusterID, true,
	).Delete(&models.APIKey{}).Error; err != nil {
		return "", "", fmt.Errorf("delete existing cluster key: %w", err)
	}

	// 2) Mint a fresh keypair
	keySuffix, err := randomB64URL(16)
	if err != nil {
		return "", "", fmt.Errorf("entropy_error: %w", err)
	}
	sec, err := randomB64URL(32)
	if err != nil {
		return "", "", fmt.Errorf("entropy_error: %w", err)
	}

	orgKey = "org_" + keySuffix
	orgSecret = sec

	keyHash := auth.SHA256Hex(orgKey)
	secretHash, err := auth.HashSecretArgon2id(orgSecret)
	if err != nil {
		return "", "", fmt.Errorf("hash_error: %w", err)
	}

	exp := now.Add(ttl)

	prefix := orgKey
	if len(prefix) > 12 {
		prefix = prefix[:12]
	}

	rec := models.APIKey{
		OrgID:       &orgID,
		Scope:       "org",
		Purpose:     "cluster_bastion",
		ClusterID:   &clusterID,
		IsEphemeral: true,
		Name:        name,
		KeyHash:     keyHash,
		SecretHash:  &secretHash,
		ExpiresAt:   &exp,
		Revoked:     false,
		Prefix:      &prefix,
	}

	if err := db.Create(&rec).Error; err != nil {
		return "", "", fmt.Errorf("db_error: %w", err)
	}

	return orgKey, orgSecret, nil
}
