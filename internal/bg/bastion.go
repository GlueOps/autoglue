package bg

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/dyaksa/archer/job"
	"github.com/glueops/autoglue/internal/db"
	"github.com/glueops/autoglue/internal/db/models"
	"github.com/glueops/autoglue/internal/utils"
	"github.com/google/uuid"
	"golang.org/x/crypto/ssh"
)

type BastionBootstrapArgs struct {
}

type BastionBootstrapResult struct {
	Status       string      `json:"status"`
	Processed    int         `json:"processed"`
	Ready        int         `json:"ready"`
	Failed       int         `json:"failed"`
	ElapsedMs    int         `json:"elapsed_ms"`
	FailedServer []uuid.UUID `json:"failed_server_ids"`
}

func BastionBootstrap(ctx context.Context, j job.Job) (any, error) {
	start := time.Now()
	log.Printf("[bastion] scan for pending bastions...")

	var bastions []models.Server
	if err := db.DB.
		Preload("SshKey").
		Where("role = ? AND status = ?", "bastion", "pending").
		Find(&bastions).Error; err != nil {
		return nil, err
	}

	if len(bastions) == 0 {
		log.Printf("[bastion] nothing to do")
		return BastionBootstrapResult{
			Status:    "ok",
			Processed: 0,
			Ready:     0,
			Failed:    0,
			ElapsedMs: int(time.Since(start) / time.Millisecond),
		}, nil
	}

	res := BastionBootstrapResult{Status: "ok", Processed: 0}
	for _, s := range bastions {
		select {
		case <-ctx.Done():
			return res, ctx.Err()
		default:
		}

		// 2) Claim atomically so only one worker does it
		claimed, err := claimPendingBastion(s.ID)
		if err != nil {
			log.Printf("[bastion] claim %s error: %v", s.ID, err)
			res.Failed++
			res.FailedServer = append(res.FailedServer, s.ID)
			continue
		}
		if !claimed {
			continue // someone else took it
		}

		// 3) Decrypt the private key
		privPEM, err := utils.DecryptForOrg(s.OrganizationID, s.SshKey.EncryptedPrivateKey, s.SshKey.PrivateIV, s.SshKey.PrivateTag)
		if err != nil {
			log.Printf("[bastion] %s decrypt key: %v", s.ID, err)
			_ = markFailed(s.ID, fmt.Errorf("decrypt key: %w", err))
			res.Failed++
			res.FailedServer = append(res.FailedServer, s.ID)
			continue
		}

		// 4) Provision over SSH
		addr := s.IPAddress
		if !strings.Contains(addr, ":") {
			addr = addr + ":22"
		}
		err = provisionDocker(ctx, addr, s.SSHUser, []byte(privPEM))
		if err != nil {
			log.Printf("[bastion] %s provision failed: %v", s.ID, err)
			_ = markFailed(s.ID, err)
			res.Failed++
			res.FailedServer = append(res.FailedServer, s.ID)
			continue
		}

		// 5) Mark ready
		if err := markReady(s.ID); err != nil {
			log.Printf("[bastion] %s mark ready err: %v", s.ID, err)
			res.Failed++
			res.FailedServer = append(res.FailedServer, s.ID)
			continue
		}

		res.Ready++
		res.Processed++
	}

	res.ElapsedMs = int(time.Since(start) / time.Millisecond)
	log.Printf("[bastion] done: processed=%d ready=%d failed=%d", res.Processed, res.Ready, res.Failed)
	return res, nil
}

func claimPendingBastion(id uuid.UUID) (bool, error) {
	tx := db.DB.Model(&models.Server{}).
		Where("id = ? AND status = ?", id, "pending").
		Update("status", "provisioning")
	if tx.Error != nil {
		return false, tx.Error
	}
	return tx.RowsAffected == 1, nil
}

func markReady(id uuid.UUID) error {
	return db.DB.Model(&models.Server{}).
		Where("id = ?", id).
		Update("status", "ready").Error
}

func markFailed(id uuid.UUID, cause error) error {
	msg := cause.Error()
	if len(msg) > 1000 {
		msg = msg[:1000]
	}
	// You can also add a separate column for error details if you want to preserve more text.
	return db.DB.Model(&models.Server{}).
		Where("id = ?", id).
		Updates(map[string]any{
			"status": "failed",
		}).Error
}

func provisionDocker(ctx context.Context, addr, user string, privKeyPEM []byte) error {
	signer, err := ssh.ParsePrivateKey(privKeyPEM)
	if err != nil {
		return fmt.Errorf("parse private key: %w", err)
	}

	cfg := &ssh.ClientConfig{
		User:            user,
		Auth:            []ssh.AuthMethod{ssh.PublicKeys(signer)},
		Timeout:         20 * time.Second,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // TODO: replace with known_hosts verification
	}

	client, err := ssh.Dial("tcp", addr, cfg)
	if err != nil {
		return fmt.Errorf("dial ssh %s: %w", addr, err)
	}
	defer client.Close()

	script := `
set -euo pipefail

if command -v docker >/dev/null 2>&1; then
  echo "docker already installed"
else
  if command -v apt-get >/dev/null 2>&1; then
    export DEBIAN_FRONTEND=noninteractive
    sudo apt-get update -y
    sudo apt-get install -y ca-certificates curl gnupg lsb-release
    curl -fsSL https://get.docker.com | sh
  elif command -v dnf >/dev/null 2>&1; then
    sudo dnf -y install dnf-plugins-core || true
    curl -fsSL https://get.docker.com | sh
  elif command -v yum >/dev/null 2>&1; then
    sudo yum -y install yum-utils || true
    curl -fsSL https://get.docker.com | sh
  else
    curl -fsSL https://get.docker.com | sh
  fi
fi

# Ensure service is enabled and running
if command -v systemctl >/dev/null 2>&1; then
  sudo systemctl enable docker || true
  sudo systemctl restart docker || sudo systemctl start docker || true
fi

docker --version
`

	return runRemote(ctx, client, "bash -lc "+quoteForShell(script))
}

func runRemote(ctx context.Context, client *ssh.Client, cmd string) error {
	sess, err := client.NewSession()
	if err != nil {
		return fmt.Errorf("new session: %w", err)
	}
	defer sess.Close()

	var out, stderr bytes.Buffer
	sess.Stdout = &out
	sess.Stderr = &stderr

	done := make(chan error, 1)
	go func() { done <- sess.Run(cmd) }()

	select {
	case <-ctx.Done():
		_ = sess.Signal(ssh.SIGKILL) // best-effort
		return ctx.Err()
	case err := <-done:
		if err != nil {
			if stderr.Len() > 0 {
				return errors.New(strings.TrimSpace(stderr.String()))
			}
			return err
		}
	}
	return nil
}

func quoteForShell(s string) string {
	// naive single-quote wrapper; escape single quotes for bash
	return fmt.Sprintf("'%s'", strings.ReplaceAll(s, "'", `'\''`))
}
