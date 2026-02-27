package bg

import (
	"context"
	"encoding/base64"
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

// ----- Public types -----

type BastionBootstrapArgs struct {
	IntervalS int `json:"interval_seconds,omitempty"`
}

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

func BastionBootstrapWorker(db *gorm.DB, jobs *Jobs) archer.WorkerFn {
	return func(ctx context.Context, j job.Job) (any, error) {
		args := BastionBootstrapArgs{IntervalS: 120}
		jobID := j.ID
		start := time.Now()

		_ = j.ParseArguments(&args)
		if args.IntervalS <= 0 {
			args.IntervalS = 120
		}

		var servers []models.Server
		if err := db.
			Preload("SshKey").
			Where("role = ? AND status = ?", "bastion", "pending").
			Find(&servers).Error; err != nil {
			log.Printf("[bastion] level=ERROR job=%s step=query msg=%q", jobID, err)
			return nil, err
		}

		// log.Printf("[bastion] level=INFO job=%s step=start count=%d", jobID, len(servers))

		proc, ok, fail := 0, 0, 0
		var failedIDs []uuid.UUID
		var failures []BastionBootstrapFailure

		perHostTimeout := 8 * time.Minute

		for i := range servers {
			s := &servers[i]
			// hostStart := time.Now()
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
			out, err := sshInstallDockerWithOutput(runCtx, db, s, host, s.SSHUser, []byte(privKey))
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

		log.Debug().Int("processed", proc).Int("ready", ok).Int("failed", fail).Msg("[bastion] reconcile tick ok")

		next := time.Now().Add(time.Duration(args.IntervalS) * time.Second)
		_, _ = jobs.Enqueue(
			ctx,
			uuid.NewString(),
			"bootstrap_bastion",
			args,
			archer.WithScheduleTime(next),
			archer.WithMaxRetries(1),
		)
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
func sshInstallDockerWithOutput(
	ctx context.Context,
	db *gorm.DB,
	s *models.Server,
	host, user string,
	privateKeyPEM []byte,
) (string, error) {
	signer, err := ssh.ParsePrivateKey(privateKeyPEM)
	if err != nil {
		return "", fmt.Errorf("parse private key: %w", err)
	}

	hkcb := makeDBHostKeyCallback(db, s)

	config := &ssh.ClientConfig{
		User:            user,
		Auth:            []ssh.AuthMethod{ssh.PublicKeys(signer)},
		HostKeyCallback: hkcb,
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

# ----------- toggles (set to 0 to skip) -----------
: "${BASELINE_PKGS:=1}"
: "${INSTALL_DOCKER:=1}"
: "${SSH_HARDEN:=1}"
: "${FIREWALL:=1}"
: "${AUTO_UPDATES:=1}"
: "${TIME_SYNC:=1}"
: "${FAIL2BAN:=1}"
: "${BANNER:=1}"

# ----------- helpers -----------
have() { command -v "$1" >/dev/null 2>&1; }

# Wait for dpkg/apt locks to be released (handles cloud-init, unattended-upgrades, etc.)
apt_wait_lock() {
  local max_wait=300 waited=0
  while [ $waited -lt $max_wait ]; do
    if ! sudo fuser /var/lib/dpkg/lock-frontend /var/lib/dpkg/lock /var/lib/apt/lists/lock /var/cache/apt/archives/lock >/dev/null 2>&1; then
      return 0
    fi
    echo "Waiting for apt/dpkg lock to be released... (${waited}s/${max_wait}s)"
    sleep 5
    waited=$((waited + 5))
  done
  echo "WARNING: apt/dpkg lock still held after ${max_wait}s, proceeding anyway" >&2
}

pm=""
if have apt-get; then pm="apt"
elif have dnf; then pm="dnf"
elif have yum; then pm="yum"
elif have zypper; then pm="zypper"
elif have apk; then pm="apk"
fi

pm_update_install() {
  case "$pm" in
    apt)
      apt_wait_lock
      sudo apt-get update -y
      sudo DEBIAN_FRONTEND=noninteractive apt-get install -y --no-install-recommends "$@"
      ;;
    dnf)    sudo dnf install -y "$@" ;;
    yum)    sudo yum install -y "$@" ;;
    zypper) sudo zypper --non-interactive install -y "$@" || true ;;
    apk)    sudo apk add --no-cache "$@" ;;
    *)
      echo "Unsupported distro: couldn't detect package manager" >&2
      return 1
      ;;
  esac
}

systemd_enable_now() {
  if have systemctl; then
    sudo systemctl enable --now "$1" || true
  fi
}

sshd_reload() {
  if have systemctl && systemctl is-enabled ssh >/dev/null 2>&1; then
    sudo systemctl reload ssh || true
  elif have systemctl && systemctl is-enabled sshd >/dev/null 2>&1; then
    sudo systemctl reload sshd || true
  fi
}

# ----------- baseline packages -----------
if [ "$BASELINE_PKGS" = "1" ] && [ -n "$pm" ]; then
  pkgs_common="curl ca-certificates gnupg git jq unzip tar vim tmux htop net-tools"
  case "$pm" in
    apt)   pkgs="$pkgs_common ufw openssh-client" ;;
    dnf|yum) pkgs="$pkgs_common firewalld openssh-clients" ;;
    zypper)  pkgs="$pkgs_common firewalld openssh" ;;
    apk)     pkgs="$pkgs_common openssh-client" ;;
  esac
  pm_update_install $pkgs || true
fi

# ----------- docker & compose v2 -----------
if [ "$INSTALL_DOCKER" = "1" ]; then
  if ! have docker; then
    if [ "$pm" = "apt" ]; then apt_wait_lock; fi
    curl -fsSL https://get.docker.com | sh
  fi

  # try to enable/start (handles distros with systemd)
  if have systemctl; then
    sudo systemctl enable --now docker || true
  fi

  # add current ssh user to docker group if exists
  if getent group docker >/dev/null 2>&1; then
    sudo usermod -aG docker "$(id -un)" || true
  fi

  # docker compose v2 (plugin) if missing
  if ! docker compose version >/dev/null 2>&1; then
    # Try package first (Debian/Ubuntu name)
    if [ "$pm" = "apt" ]; then
      apt_wait_lock
      sudo apt-get update -y
      sudo apt-get install -y docker-compose-plugin || true
    fi

    # Fallback: install static plugin binary under ~/.docker/cli-plugins
    if ! docker compose version >/dev/null 2>&1; then
      mkdir -p ~/.docker/cli-plugins
      arch="$(uname -m)"
      case "$arch" in
        x86_64|amd64) arch="x86_64" ;;
        aarch64|arm64) arch="aarch64" ;;
      esac
      curl -fsSL -o ~/.docker/cli-plugins/docker-compose \
        "https://github.com/docker/compose/releases/download/v2.29.7/docker-compose-$(uname -s)-$arch"
      chmod +x ~/.docker/cli-plugins/docker-compose
    fi
  fi
fi

# ----------- SSH hardening (non-destructive: separate conf file) -----------
if [ "$SSH_HARDEN" = "1" ]; then
  confd="/etc/ssh/sshd_config.d"
  if [ -d "$confd" ] && [ -w "$confd" ]; then
    sudo tee "$confd/10-bastion.conf" >/dev/null <<'EOF'
# Bastion hardening
PasswordAuthentication no
ChallengeResponseAuthentication no
KbdInteractiveAuthentication no
UsePAM yes
PermitEmptyPasswords no
PubkeyAuthentication yes
ClientAliveInterval 300
ClientAliveCountMax 2
LoginGraceTime 20
MaxAuthTries 3
MaxSessions 10
AllowAgentForwarding no
X11Forwarding no
EOF
    sshd_reload
  else
    echo "Skipping SSH hardening: $confd not present or not writable" >&2
  fi

  # lock root password (no effect if already locked)
  if have passwd; then
    sudo passwd -l root || true
  fi
fi

# ----------- firewall -----------
if [ "$FIREWALL" = "1" ]; then
  if have ufw; then
    # Keep it minimal: allow SSH and rate-limit
    sudo ufw --force reset || true
    sudo ufw default deny incoming
    sudo ufw default allow outgoing
    sudo ufw allow OpenSSH || sudo ufw allow 22/tcp
    sudo ufw limit OpenSSH || true
    sudo ufw --force enable
  elif have firewall-cmd; then
    systemd_enable_now firewalld
    sudo firewall-cmd --permanent --add-service=ssh || sudo firewall-cmd --permanent --add-port=22/tcp
    sudo firewall-cmd --reload || true
  else
    echo "No supported firewall tool detected; skipping." >&2
  fi
fi

# ----------- unattended / automatic updates -----------
if [ "$AUTO_UPDATES" = "1" ] && [ -n "$pm" ]; then
  case "$pm" in
    apt)
      pm_update_install unattended-upgrades apt-listchanges || true
      sudo dpkg-reconfigure -f noninteractive unattended-upgrades || true
      sudo tee /etc/apt/apt.conf.d/20auto-upgrades >/dev/null <<'EOF'
APT::Periodic::Update-Package-Lists "1";
APT::Periodic::Unattended-Upgrade "1";
APT::Periodic::AutocleanInterval "7";
EOF
      ;;
    dnf)
      pm_update_install dnf-automatic || true
      sudo sed -i 's/^apply_updates = .*/apply_updates = yes/' /etc/dnf/automatic.conf || true
      systemd_enable_now dnf-automatic.timer
      ;;
    yum)
      pm_update_install yum-cron || true
      sudo sed -i 's/apply_updates = no/apply_updates = yes/' /etc/yum/yum-cron.conf || true
      systemd_enable_now yum-cron
      ;;
    zypper)
      pm_update_install pkgconf-pkg-config || true
      # SUSE has automatic updates via transactional-update / yast2-online-update; skipping heavy config.
      ;;
    apk)
      # Alpine: no official unattended updater; consider periodic 'apk upgrade' via cron (skipped by default).
      ;;
  esac
fi

# ----------- time sync -----------
if [ "$TIME_SYNC" = "1" ]; then
  if have timedatectl; then
    # Prefer systemd-timesyncd if available; else install/enable chrony
    if [ -f /lib/systemd/system/systemd-timesyncd.service ] || [ -f /usr/lib/systemd/system/systemd-timesyncd.service ]; then
      systemd_enable_now systemd-timesyncd
    else
      pm_update_install chrony || true
      systemd_enable_now chronyd || systemd_enable_now chrony || true
    fi
    timedatectl set-ntp true || true
  else
    pm_update_install chrony || true
    systemd_enable_now chronyd || systemd_enable_now chrony || true
  fi
fi

# ----------- fail2ban (basic sshd jail) -----------
if [ "$FAIL2BAN" = "1" ]; then
  pm_update_install fail2ban || true
  if [ -d /etc/fail2ban ]; then
    sudo tee /etc/fail2ban/jail.d/sshd.local >/dev/null <<'EOF'
[sshd]
enabled = true
port    = ssh
logpath = %(sshd_log)s
maxretry = 4
bantime = 1h
findtime = 10m
EOF
    systemd_enable_now fail2ban
  fi
fi

# ----------- SSH banner / MOTD -----------
if [ "$BANNER" = "1" ]; then
  if [ -w /etc/issue.net ] || sudo test -w /etc/issue.net; then
    sudo tee /etc/issue.net >/dev/null <<'EOF'
NOTICE: Authorized use only. Activity may be monitored and reported.
EOF
    # Ensure banner is enabled via our bastion conf
    if [ -d /etc/ssh/sshd_config.d ]; then
      if ! grep -q '^Banner ' /etc/ssh/sshd_config.d/10-bastion.conf 2>/dev/null; then
        echo 'Banner /etc/issue.net' | sudo tee -a /etc/ssh/sshd_config.d/10-bastion.conf >/dev/null
        sshd_reload
      fi
    fi
  fi
fi

echo "Bootstrap complete. If you were added to the docker group, log out and back in to apply."
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

// makeDBHostKeyCallback returns a HostKeyCallback bound to a specific server row.
// TOFU semantics:
//   - If s.SSHHostKey is empty: store the current key in DB and accept.
//   - If s.SSHHostKey is set: require exact match, else error (possible MITM/reinstall).
func makeDBHostKeyCallback(db *gorm.DB, s *models.Server) ssh.HostKeyCallback {
	return func(hostname string, remote net.Addr, key ssh.PublicKey) error {
		algo := key.Type()
		enc := base64.StdEncoding.EncodeToString(key.Marshal())

		// First-time connect: persist key (TOFU).
		if s.SSHHostKey == "" {
			if err := db.Model(&models.Server{}).
				Where("id = ? AND (ssh_host_key IS NULL or ssh_host_key = '')", s.ID).
				Updates(map[string]any{
					"ssh_host_key":      enc,
					"ssh_host_key_algo": algo,
				}).Error; err != nil {
				return fmt.Errorf("store new host key for %s (%s): %w", hostname, s.ID, err)
			}

			s.SSHHostKey = enc
			s.SSHHostKeyAlgo = algo
			return nil
		}

		if s.SSHHostKeyAlgo != algo || s.SSHHostKey != enc {
			return fmt.Errorf(
				"host key mismatch for %s (server_id=%s, stored=%s/%s, got=%s/%s) - POSSIBLE MITM or host reinstalled",
				hostname, s.ID, s.SSHHostKeyAlgo, s.SSHHostKey, algo, enc,
			)
		}
		return nil
	}
}
