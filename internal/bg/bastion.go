package bg

import (
	"context"
	"fmt"
	"log"
	"net"
	"strings"
	"time"

	"github.com/dyaksa/archer"
	"github.com/dyaksa/archer/job"
	"github.com/glueops/autoglue/internal/models"
	"github.com/glueops/autoglue/internal/utils"
	"github.com/google/uuid"
	"golang.org/x/crypto/ssh"
	"gorm.io/gorm"
)

// ----- Public types -----

type BastionBootstrapArgs struct{}

type BastionBootstrapFailure struct {
	ID     uuid.UUID `json:"id"`
	Step   string    `json:"step"`
	Reason string    `json:"reason"`
}

type BastionBootstrapResult struct {
	Status       string                    `json:"status"`
	Processed    int                       `json:"processed"`
	Ready        int                       `json:"ready"`
	Failed       int                       `json:"failed"`
	ElapsedMs    int                       `json:"elapsed_ms"`
	FailedServer []uuid.UUID               `json:"failed_server_ids"`
	Failures     []BastionBootstrapFailure `json:"failures"`
}

// ----- Worker -----

func BastionBootstrapWorker(db *gorm.DB) archer.WorkerFn {
	return func(ctx context.Context, j job.Job) (any, error) {
		jobID := j.ID
		start := time.Now()

		var servers []models.Server
		if err := db.
			Preload("SshKey").
			Where("role = ? AND status = ?", "bastion", "pending").
			Find(&servers).Error; err != nil {
			log.Printf("[bastion] level=ERROR job=%s step=query msg=%q", jobID, err)
			return nil, err
		}

		log.Printf("[bastion] level=INFO job=%s step=start count=%d", jobID, len(servers))

		proc, ok, fail := 0, 0, 0
		var failedIDs []uuid.UUID
		var failures []BastionBootstrapFailure

		perHostTimeout := 8 * time.Minute

		for i := range servers {
			s := &servers[i]
			hostStart := time.Now()
			proc++

			// 1) Defensive IP check
			if s.PublicIPAddress == nil || *s.PublicIPAddress == "" {
				fail++
				failedIDs = append(failedIDs, s.ID)
				failures = append(failures, BastionBootstrapFailure{ID: s.ID, Step: "ip_check", Reason: "missing public ip"})
				logHostErr(jobID, s, "ip_check", fmt.Errorf("missing public ip"))
				_ = setServerStatus(db, s.ID, "failed")
				continue
			}

			// 2) Move to provisioning
			if err := setServerStatus(db, s.ID, "provisioning"); err != nil {
				fail++
				failedIDs = append(failedIDs, s.ID)
				failures = append(failures, BastionBootstrapFailure{ID: s.ID, Step: "set_provisioning", Reason: err.Error()})
				logHostErr(jobID, s, "set_provisioning", err)
				continue
			}

			// 3) Decrypt private key for org
			privKey, err := utils.DecryptForOrg(
				s.OrganizationID,
				s.SshKey.EncryptedPrivateKey,
				s.SshKey.PrivateIV,
				s.SshKey.PrivateTag,
				db,
			)
			if err != nil {
				fail++
				failedIDs = append(failedIDs, s.ID)
				failures = append(failures, BastionBootstrapFailure{ID: s.ID, Step: "decrypt_key", Reason: err.Error()})
				logHostErr(jobID, s, "decrypt_key", err)
				_ = setServerStatus(db, s.ID, "failed")
				continue
			}

			// 4) SSH + install docker
			host := net.JoinHostPort(*s.PublicIPAddress, "22")
			runCtx, cancel := context.WithTimeout(ctx, perHostTimeout)
			out, err := sshInstallDockerWithOutput(runCtx, host, s.SSHUser, []byte(privKey))
			cancel()

			if err != nil {
				fail++
				failedIDs = append(failedIDs, s.ID)
				failures = append(failures, BastionBootstrapFailure{ID: s.ID, Step: "ssh_install", Reason: err.Error()})
				// include a short tail of output to speed debugging without flooding logs
				tail := out
				if len(tail) > 800 {
					tail = tail[len(tail)-800:]
				}
				logHostErr(jobID, s, "ssh_install", fmt.Errorf("%v | tail=%q", err, tail))
				_ = setServerStatus(db, s.ID, "failed")
				continue
			}

			// 5) Mark ready
			if err := setServerStatus(db, s.ID, "ready"); err != nil {
				fail++
				failedIDs = append(failedIDs, s.ID)
				failures = append(failures, BastionBootstrapFailure{ID: s.ID, Step: "set_ready", Reason: err.Error()})
				logHostErr(jobID, s, "set_ready", err)
				_ = setServerStatus(db, s.ID, "failed")
				continue
			}

			ok++
			logHostInfo(jobID, s, "done", "host completed",
				"elapsed_ms", time.Since(hostStart).Milliseconds())
		}

		res := BastionBootstrapResult{
			Status:       "ok",
			Processed:    proc,
			Ready:        ok,
			Failed:       fail,
			ElapsedMs:    int(time.Since(start).Milliseconds()),
			FailedServer: failedIDs,
			Failures:     failures,
		}

		log.Printf("[bastion] level=INFO job=%s step=finish processed=%d ready=%d failed=%d elapsed_ms=%d",
			jobID, proc, ok, fail, res.ElapsedMs)

		return res, nil
	}
}

// ----- Helpers -----

func setServerStatus(db *gorm.DB, id uuid.UUID, status string) error {
	return db.Model(&models.Server{}).
		Where("id = ?", id).
		Updates(map[string]any{
			"status":     status,
			"updated_at": time.Now(),
		}).Error
}

// uniform log helpers for consistent, greppable output
func logHostErr(jobID string, s *models.Server, step string, err error) {
	ip := ""
	if s.PublicIPAddress != nil {
		ip = *s.PublicIPAddress
	}
	log.Printf("[bastion] level=ERROR job=%s server_id=%s host=%s step=%s msg=%q",
		jobID, s.ID, ip, step, err)
}

func logHostInfo(jobID string, s *models.Server, step, msg string, kv ...any) {
	ip := ""
	if s.PublicIPAddress != nil {
		ip = *s.PublicIPAddress
	}
	log.Printf("[bastion] level=INFO job=%s server_id=%s host=%s step=%s %s kv=%v",
		jobID, s.ID, ip, step, msg, kv)
}

// ----- SSH & command execution -----

// returns combined stdout/stderr so caller can log it on error
// returns combined stdout/stderr so caller can log it on error
func sshInstallDockerWithOutput(ctx context.Context, host, user string, privateKeyPEM []byte) (string, error) {
	signer, err := ssh.ParsePrivateKey(privateKeyPEM)
	if err != nil {
		return "", fmt.Errorf("parse private key: %w", err)
	}

	config := &ssh.ClientConfig{
		User:            user,
		Auth:            []ssh.AuthMethod{ssh.PublicKeys(signer)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // TODO: known_hosts verification
		Timeout:         30 * time.Second,
	}

	// context-aware dial
	dialer := &net.Dialer{}
	conn, err := dialer.DialContext(ctx, "tcp", host)
	if err != nil {
		return "", fmt.Errorf("dial: %w", err)
	}
	defer conn.Close()

	c, chans, reqs, err := ssh.NewClientConn(conn, host, config)
	if err != nil {
		return "", fmt.Errorf("ssh handshake: %w", err)
	}
	client := ssh.NewClient(c, chans, reqs)
	defer client.Close()

	sess, err := client.NewSession()
	if err != nil {
		return "", fmt.Errorf("session: %w", err)
	}
	defer sess.Close()

	// --- script to run remotely (no extra quoting) ---
	script := `
set -euxo pipefail

if ! command -v docker >/dev/null 2>&1; then
  curl -fsSL https://get.docker.com | sh
fi

# try to enable/start (handles distros with systemd)
if command -v systemctl >/dev/null 2>&1; then
  sudo systemctl enable --now docker || true
fi

# add current ssh user to docker group if exists
if getent group docker >/dev/null 2>&1; then
  sudo usermod -aG docker "$(id -un)" || true
fi
`

	// Send script via stdin to avoid quoting/escaping issues
	sess.Stdin = strings.NewReader(script)

	// Capture combined stdout+stderr
	out, runErr := sess.CombinedOutput("bash -s")
	return string(out), wrapSSHError(runErr, string(out))
}

// annotate common SSH/remote failure modes to speed triage
func wrapSSHError(err error, output string) error {
	if err == nil {
		return nil
	}
	switch {
	case strings.Contains(output, "Could not resolve host"):
		return fmt.Errorf("remote run: name resolution failed: %w", err)
	case strings.Contains(output, "Permission denied"):
		return fmt.Errorf("remote run: permission denied (check user/key/authorized_keys): %w", err)
	case strings.Contains(output, "apt-get"):
		return fmt.Errorf("remote run: apt failed: %w", err)
	case strings.Contains(output, "yum"):
		return fmt.Errorf("remote run: yum failed: %w", err)
	default:
		return fmt.Errorf("remote run: %w", err)
	}
}

// super simple escaping for a here-string; avoids quoting hell
func sshEscape(s string) string {
	return fmt.Sprintf("%q", s)
}
