package bg

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/glueops/autoglue/internal/models"
	"github.com/glueops/autoglue/internal/utils"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/riverqueue/river"
	"github.com/riverqueue/river/rivertype"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"golang.org/x/crypto/ssh"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// ----- Public types -----

// BastionSweepArgs drives the periodic claim tick. It does no SSH work itself:
// it moves servers from pending to provisioning and fans out one bootstrap job
// per server.
type BastionSweepArgs struct{}

func (BastionSweepArgs) Kind() string { return "bastion_sweep" }

func (BastionSweepArgs) InsertOpts() river.InsertOpts {
	return river.InsertOpts{Queue: QueueMaintenance, MaxAttempts: 2}
}

// BastionBootstrapArgs bootstraps exactly one bastion.
//
// One job per server, rather than one job per tick, because the claim flips a
// server to provisioning up front: if a single job owned every claimed server
// and then hit its timeout, the ones it never reached would sit in
// provisioning forever with nothing to reclaim them.
type BastionBootstrapArgs struct {
	ServerID uuid.UUID `json:"server_id"`
}

func (BastionBootstrapArgs) Kind() string { return "bootstrap_bastion" }

func (BastionBootstrapArgs) InsertOpts() river.InsertOpts {
	return river.InsertOpts{
		Queue:       QueueClusters,
		MaxAttempts: 1,
		// Never bootstrap the same host twice concurrently: the remote script
		// takes apt/dpkg locks and two runs would fight over them.
		UniqueOpts: river.UniqueOpts{
			ByArgs: true,
			ByState: []rivertype.JobState{
				rivertype.JobStateAvailable, rivertype.JobStateScheduled,
				rivertype.JobStateRunning, rivertype.JobStateRetryable,
				rivertype.JobStatePending,
			},
		},
	}
}

type BastionBootstrapFailure struct {
	ID     uuid.UUID `json:"id"`
	Step   string    `json:"step"`
	Reason string    `json:"reason"`
}

type BastionSweepResult struct {
	Status     string      `json:"status"`
	Claimed    int         `json:"claimed"`
	Dispatched int         `json:"dispatched"`
	Duplicates int         `json:"duplicates"`
	ServerIDs  []uuid.UUID `json:"server_ids"`
}

type BastionBootstrapResult struct {
	Status    string    `json:"status"`
	ServerID  uuid.UUID `json:"server_id"`
	ElapsedMs int       `json:"elapsed_ms"`
}

// ----- Sweep (claim + dispatch) -----

type BastionSweepWorker struct {
	river.WorkerDefaults[BastionSweepArgs]
	db *gorm.DB
}

func (w *BastionSweepWorker) Timeout(*river.Job[BastionSweepArgs]) time.Duration {
	return time.Minute
}

func (w *BastionSweepWorker) Work(ctx context.Context, j *river.Job[BastionSweepArgs]) error {
	db := w.db

	// Atomically claim pending bastion servers using SELECT FOR UPDATE SKIP LOCKED
	// so that concurrent sweeps never claim the same server. `pending` is the
	// enqueue signal; flipping to `provisioning` is what removes a server from
	// this tick's queue.
	var claimedIDs []uuid.UUID
	if err := db.Transaction(func(tx *gorm.DB) error {
		if err := tx.
			Model(&models.Server{}).
			Where("role = ? AND status = ?", "bastion", "pending").
			Clauses(clause.Locking{Strength: "UPDATE", Options: "SKIP LOCKED"}).
			Pluck("id", &claimedIDs).Error; err != nil {
			return err
		}
		if len(claimedIDs) == 0 {
			return nil
		}
		return tx.Model(&models.Server{}).
			Where("id IN ?", claimedIDs).
			Updates(map[string]any{
				"status":     "provisioning",
				"updated_at": time.Now(),
			}).Error
	}); err != nil {
		log.Printf("[bastion] level=ERROR job=%d step=claim_servers msg=%q", j.ID, err)
		return err
	}

	if len(claimedIDs) == 0 {
		return nil
	}

	dispatched, duplicates := dispatchBastionBootstraps(ctx, river.ClientFromContext[pgx.Tx](ctx), db, claimedIDs)

	log.Info().Int("claimed", len(claimedIDs)).Int("dispatched", dispatched).
		Int("duplicates", duplicates).Msg("[bastion] sweep dispatched bootstrap jobs")

	if err := river.RecordOutput(ctx, BastionSweepResult{
		Status:     "ok",
		Claimed:    len(claimedIDs),
		Dispatched: dispatched,
		Duplicates: duplicates,
		ServerIDs:  claimedIDs,
	}); err != nil {
		log.Warn().Err(err).Msg("[bastion] could not record sweep output")
	}
	return nil
}

// dispatchBastionBootstraps inserts one bootstrap job per claimed server and
// reports how many were new and how many deduplicated against a live job.
func dispatchBastionBootstraps(ctx context.Context, client *Client, db *gorm.DB, claimedIDs []uuid.UUID) (dispatched, duplicates int) {
	for _, id := range claimedIDs {
		res, err := client.Insert(ctx, BastionBootstrapArgs{ServerID: id}, nil)
		if err != nil {
			// Hand the server back so a later tick retries it, rather than
			// leaving it stranded in provisioning.
			log.Error().Err(err).Str("server_id", id.String()).
				Msg("[bastion] could not dispatch bootstrap; returning server to pending")
			_ = setServerStatus(db, id, "pending")
			continue
		}
		if res.UniqueSkippedAsDuplicate {
			// Usually a server set back to pending while its bootstrap is
			// snoozed waiting for sshd. That job is still live and will set the
			// final status, so there is nothing to hand back.
			log.Info().Str("server_id", id.String()).Int64("existing_job_id", res.Job.ID).
				Msg("[bastion] bootstrap already queued for server; skipped duplicate dispatch")
			duplicates++
			continue
		}
		dispatched++
	}

	return dispatched, duplicates
}

// ----- Bootstrap (one server) -----

type BastionBootstrapWorker struct {
	river.WorkerDefaults[BastionBootstrapArgs]
	db *gorm.DB
}

// Timeout bounds a single host. The remote script installs packages over a
// network that may be slow, so this is generous, but it is per server rather
// than per tick.
func (w *BastionBootstrapWorker) Timeout(*river.Job[BastionBootstrapArgs]) time.Duration {
	return 30 * time.Minute
}

func (w *BastionBootstrapWorker) Work(ctx context.Context, j *river.Job[BastionBootstrapArgs]) error {
	db := w.db
	jobID := strconv.FormatInt(j.ID, 10)
	start := time.Now()

	var s models.Server
	if err := db.Preload("SshKey").
		Where("id = ?", j.Args.ServerID).
		First(&s).Error; err != nil {
		// The server was deleted between claim and execution. Nothing to do,
		// and failing would only requeue a job that can never succeed.
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Warn().Str("server_id", j.Args.ServerID.String()).
				Msg("[bastion] server vanished before bootstrap")
			return nil
		}
		return err
	}

	// Open the sink before the first thing that can fail. A bastion with no
	// public IP, or an org key that will not decrypt, are the two most likely
	// early failures — and previously neither wrote anything to job_logs, so
	// the UI showed "no output recorded" and the reason stayed in worker
	// stdout, which is exactly what this is meant to replace.
	sink := NewLogSink(db, j.ID, s.OrganizationID, models.JobLogSubjectServer, s.ID)
	defer func() { _ = sink.Close() }()

	sink.System(fmt.Sprintf("claimed bastion %s (%s)", s.ID, s.Hostname))

	// Whether this host had booted and been reached before this job started,
	// decided on the first attempt and carried across snoozes: TOFU records the
	// host key during that first handshake, so from the second attempt on the
	// server row can no longer tell us.
	meta := parseBastionMeta(j.Metadata)
	if meta.KnownHost == nil {
		known := s.SSHHostKey != ""
		meta.KnownHost = &known
		_ = river.MetadataSet(ctx, bastionMetaKnownHost, known)
	}

	fail := func(step string, err error) error {
		sink.System("bootstrap failed at " + step + ": " + err.Error())
		logHostErr(jobID, &s, step, err)
		_ = setServerStatus(db, s.ID, "failed")
		// The failure is recorded on the server row and in job_logs; returning
		// nil keeps it out of River's retry path, since the operator re-runs by
		// setting the server back to pending.
		return nil
	}

	// 1) Defensive IP check
	if s.PublicIPAddress == nil || *s.PublicIPAddress == "" {
		return fail("ip_check", fmt.Errorf("missing public ip"))
	}

	// 2) Decrypt private key for org
	privKey, err := utils.DecryptForOrg(
		s.OrganizationID,
		s.SshKey.EncryptedPrivateKey,
		s.SshKey.PrivateIV,
		s.SshKey.PrivateTag,
		db,
	)
	if err != nil {
		return fail("decrypt_key", err)
	}

	// 3) SSH + install docker. Output streams into job_logs under this server,
	// so a failed bootstrap can be read back from the API instead of hunting
	// through worker pod stdout for a truncated tail.
	host := net.JoinHostPort(*s.PublicIPAddress, bastionSSHPort)
	sink.System(fmt.Sprintf("connecting to %s as %s", host, s.SSHUser))

	out, err := sshInstallDockerWithOutput(ctx, db, &s, host, s.SSHUser, []byte(privKey), sink)
	if err != nil {
		// The worker is shutting down. The script is safe to re-run, so hand
		// the job back to run again as soon as a worker is up, rather than
		// failing a bastion that did nothing wrong. A job cancelled on purpose
		// (River UI, JobCancel) is not handed back; it fails below.
		if errors.Is(err, context.Canceled) && ctx.Err() != nil &&
			!errors.Is(context.Cause(ctx), river.ErrJobCancelledRemotely) {
			sink.System("worker shutting down; bootstrap will restart")
			return river.JobSnooze(0)
		}

		// A bastion is usually claimed seconds after its VM was created, well
		// before sshd is listening, and cloud-init may reboot it once more
		// while the script waits on it. Snooze rather than fail so the server
		// stays in provisioning and the bootstrap picks up once the host has
		// booted. The wait is measured from the job's creation, which snoozing
		// does not reset, or from the last lost connection, since a reboot
		// starts the boot over.
		since := j.CreatedAt
		if meta.SessionLostAt > 0 {
			if t := time.Unix(meta.SessionLostAt, 0); t.After(since) {
				since = t
			}
		}
		waited := time.Since(since)
		limits := bastionLimits()
		if delay, why, ok := bastionRetryDelay(err, waited, meta, limits); ok {
			if isSessionLost(err) {
				_ = river.MetadataSet(ctx, bastionMetaSessionLost, meta.SessionLost+1)
				_ = river.MetadataSet(ctx, bastionMetaSessionLostAt, time.Now().Unix())
			}
			sink.System(fmt.Sprintf("%s (%v); retrying in %s (waited %s of %s)",
				why, err, delay, waited.Round(time.Second), limits.window))
			logHostInfo(jobID, &s, "ssh_wait", why+", snoozing",
				"delay", delay, "waited", waited.Round(time.Second), "reason", err.Error())
			return river.JobSnooze(delay)
		}
		switch {
		case isSessionLost(err):
			err = fmt.Errorf("connection lost %d times, giving up: %w", meta.SessionLost+1, err)
		case isHostNotReady(err):
			err = fmt.Errorf("host still unreachable after %s: %w", waited.Round(time.Second), err)
		case isAuthRejected(err):
			err = fmt.Errorf("auth still rejected after %s: %w", waited.Round(time.Second), err)
		}
		tail := out
		if len(tail) > 800 {
			tail = tail[len(tail)-800:]
		}
		return fail("ssh_install", fmt.Errorf("%v | tail=%q", err, tail))
	}

	// 4) Mark ready
	if err := setServerStatus(db, s.ID, "ready"); err != nil {
		return fail("set_ready", err)
	}

	sink.System("bastion ready")

	log.Info().Str("server_id", s.ID.String()).
		Int64("elapsed_ms", time.Since(start).Milliseconds()).
		Msg("[bastion] bootstrap ok")

	if err := river.RecordOutput(ctx, BastionBootstrapResult{
		Status:    "ok",
		ServerID:  s.ID,
		ElapsedMs: int(time.Since(start).Milliseconds()),
	}); err != nil {
		log.Warn().Err(err).Msg("[bastion] could not record output")
	}
	return nil
}

// ----- Helpers -----

// bastionSSHPort is a var only so tests can point the bootstrap at a local
// fake sshd.
var bastionSSHPort = "22"

const (
	bastionRetryInitial = 5 * time.Second
	bastionRetryMax     = time.Minute
)

// bastionSSHWait is how long a bootstrap keeps waiting for a freshly created
// host to start accepting SSH before giving up and marking it failed.
func bastionSSHWait() time.Duration {
	return interval("bastion.ssh_wait_seconds", 10*time.Minute)
}

// bastionAuthGrace is how long a rejected key is still put down to the host
// booting. sshd can be up before it will let the user in: pam_nologin refuses
// non-root logins until systemd-user-sessions has run, and user-data that
// creates the user or its key late (runcmd, write_files) leaves it missing for
// all of cloud-final. Past this, a rejected key is taken to be the wrong key.
func bastionAuthGrace() time.Duration {
	return interval("bastion.ssh_auth_grace_seconds", 3*time.Minute)
}

// bastionSessionLostMax is how many times a bootstrap may lose its connection
// mid-script and start over. Counted rather than timed: a reboot at the end of
// cloud-init's package upgrade routinely lands well after the wait window
// since job creation has closed. Each attempt is still bounded by Timeout.
func bastionSessionLostMax() int {
	if n := viper.GetInt("bastion.session_lost_retries"); n > 0 {
		return n
	}
	return 3
}

type bastionRetryLimits struct {
	window         time.Duration
	authGrace      time.Duration
	sessionLostMax int
}

func bastionLimits() bastionRetryLimits {
	return bastionRetryLimits{
		window:         bastionSSHWait(),
		authGrace:      bastionAuthGrace(),
		sessionLostMax: bastionSessionLostMax(),
	}
}

// Job metadata keys a bootstrap carries across snoozes.
const (
	bastionMetaKnownHost     = "bastion_known_host"
	bastionMetaSessionLost   = "bastion_session_lost"
	bastionMetaSessionLostAt = "bastion_session_lost_at"
)

// bastionJobMeta is the job metadata a bootstrap reads back on each attempt.
// Snoozes is River's own count; snoozing does not advance Attempt, so it is
// the only counter River moves.
type bastionJobMeta struct {
	Snoozes       int   `json:"snoozes"`
	KnownHost     *bool `json:"bastion_known_host"`
	SessionLost   int   `json:"bastion_session_lost"`
	SessionLostAt int64 `json:"bastion_session_lost_at"`
}

func parseBastionMeta(metadata []byte) bastionJobMeta {
	var m bastionJobMeta
	_ = json.Unmarshal(metadata, &m)
	return m
}

// bastionRetryDelay decides whether a failed bootstrap should be snoozed, for
// how long, and why. "Host not up yet" errors are retried inside the window,
// a lost connection up to sessionLostMax times, and a rejected key inside the
// auth grace period, but only for a host never reached before: on one that has
// been, fail2ban would ban us long before the grace ran out and turn a wrong
// key into what looks like an unreachable host. Everything else (bad key, host
// key mismatch, the remote script exiting non-zero) fails immediately, since
// waiting would not change the outcome.
func bastionRetryDelay(err error, waited time.Duration, meta bastionJobMeta, l bastionRetryLimits) (time.Duration, string, bool) {
	var why string
	switch {
	case isSessionLost(err):
		if meta.SessionLost >= l.sessionLostMax {
			return 0, "", false
		}
		why = "connection lost, host may be rebooting"
	case waited >= l.window:
		return 0, "", false
	case isHostNotReady(err):
		why = "host not reachable yet"
	case isAuthRejected(err) && waited < l.authGrace && (meta.KnownHost == nil || !*meta.KnownHost):
		why = "auth rejected, host may still be booting"
	default:
		return 0, "", false
	}
	delay := bastionRetryInitial
	for i := 0; i < meta.Snoozes && delay < bastionRetryMax; i++ {
		delay *= 2
	}
	if delay > bastionRetryMax {
		delay = bastionRetryMax
	}
	return delay, why, true
}

// sshConnectError marks a failure to reach sshd at all (dial or handshake),
// as opposed to a failure of the remote script once connected.
type sshConnectError struct{ err error }

func (e *sshConnectError) Error() string { return e.err.Error() }
func (e *sshConnectError) Unwrap() error { return e.err }

// isHostNotReady reports whether err looks like a host that is still booting:
// nothing listening on 22 yet, no route while the network comes up, or sshd
// dropping the connection before the handshake completes.
func isHostNotReady(err error) bool {
	var ce *sshConnectError
	if !errors.As(err, &ce) {
		return false
	}
	for _, errno := range []syscall.Errno{
		syscall.ECONNREFUSED, syscall.EHOSTUNREACH, syscall.ENETUNREACH,
		syscall.ETIMEDOUT, syscall.ECONNRESET, syscall.ECONNABORTED,
	} {
		if errors.Is(err, errno) {
			return true
		}
	}
	var ne net.Error
	if errors.As(err, &ne) && ne.Timeout() {
		return true
	}
	return errors.Is(err, io.EOF) || errors.Is(err, io.ErrUnexpectedEOF)
}

// isAuthRejected reports whether sshd answered but refused our key. x/crypto
// returns this as an unwrapped fmt.Errorf, so it can only be matched on text.
func isAuthRejected(err error) bool {
	var ce *sshConnectError
	return errors.As(err, &ce) && strings.Contains(err.Error(), "ssh: unable to authenticate")
}

// sshSessionError marks a failure after the handshake: opening the session or
// running the remote script.
type sshSessionError struct{ err error }

func (e *sshSessionError) Error() string { return e.err.Error() }
func (e *sshSessionError) Unwrap() error { return e.err }

// isSessionLost reports whether the remote script was cut off by the
// connection going away rather than finishing. The usual cause is the host
// rebooting under it (cloud-init's package_reboot_if_required lands while the
// script sits in "cloud-init status --wait"), and the script is safe to re-run
// from the top: every step is guarded or overwrites its own file.
//
// A script that ran to completion and exited non-zero is not this, and is
// never retried. A script killed by TERM, HUP or KILL is: that is what a
// shutdown does to the session before the connection drops.
func isSessionLost(err error) bool {
	var se *sshSessionError
	if !errors.As(err, &se) {
		return false
	}
	if errors.Is(err, errKeepaliveLost) {
		return true
	}
	var missing *ssh.ExitMissingError
	if errors.As(err, &missing) {
		return true
	}
	var exit *ssh.ExitError
	if errors.As(err, &exit) {
		switch exit.Signal() {
		case "TERM", "HUP", "KILL":
			return true
		}
		return false
	}
	for _, errno := range []syscall.Errno{syscall.ECONNRESET, syscall.ECONNABORTED, syscall.EPIPE} {
		if errors.Is(err, errno) {
			return true
		}
	}
	return errors.Is(err, io.EOF) || errors.Is(err, io.ErrUnexpectedEOF)
}

// sessionErr tags a post-handshake failure, unless the supervisor closed the
// connection: then its reason is the real cause. A cancelled context is
// deliberately left untagged so it is never mistaken for a lost session.
func sessionErr(closedBecause, err error) error {
	switch {
	case closedBecause == nil:
		return &sshSessionError{err}
	case errors.Is(closedBecause, errKeepaliveLost):
		return &sshSessionError{fmt.Errorf("%w: %v", closedBecause, err)}
	default:
		return fmt.Errorf("aborted: %w (%v)", closedBecause, err)
	}
}

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
// sshInstallDockerWithOutput bootstraps a bastion over SSH, streaming the
// remote script's combined output to sink as it arrives. sink may be nil.
//
// The returned string is the trailing logMaxTailBytes only; the full transcript
// is in job_logs when a sink is supplied. That matters here because the script
// runs under `set -euxo pipefail`, so the useful evidence is the last `+` trace
// line before it died — which the old 800-byte stdout tail routinely cut off.
func sshInstallDockerWithOutput(
	ctx context.Context,
	db *gorm.DB,
	s *models.Server,
	host, user string,
	privateKeyPEM []byte,
	sink io.Writer,
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

	// context-aware dial. config.Timeout only applies to ssh.Dial, so bound
	// the dial and the handshake here: a host that is still booting can accept
	// the TCP connection and then never speak.
	dialer := &net.Dialer{Timeout: config.Timeout}
	conn, err := dialer.DialContext(ctx, "tcp", host)
	if err != nil {
		return "", &sshConnectError{fmt.Errorf("dial: %w", err)}
	}
	defer conn.Close()

	_ = conn.SetDeadline(time.Now().Add(config.Timeout))
	c, chans, reqs, err := ssh.NewClientConn(conn, host, config)
	if err != nil {
		return "", &sshConnectError{fmt.Errorf("ssh handshake: %w", err)}
	}
	_ = conn.SetDeadline(time.Time{})
	client := ssh.NewClient(c, chans, reqs)
	defer client.Close()

	// Without this a host that vanishes without a FIN (a hard reset, a dropped
	// route) leaves the session blocked until TCP gives up, and a cancelled
	// job keeps its script running remotely.
	stopSupervise := superviseSSH(ctx, client)
	defer func() { _ = stopSupervise() }()

	sess, err := client.NewSession()
	if err != nil {
		return "", sessionErr(stopSupervise(), fmt.Errorf("session: %w", err))
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
: "${APT_LOCK_WAIT_SECS:=300}"

# ----------- helpers -----------
have() { command -v "$1" >/dev/null 2>&1; }
docker_works() { have docker && sudo docker info >/dev/null 2>&1; }

# Wait for dpkg/apt locks to be released (handles cloud-init, unattended-upgrades, etc.)
apt_wait_lock() {
  local max_wait="$APT_LOCK_WAIT_SECS" waited=0
  local lock_files="/var/lib/dpkg/lock-frontend /var/lib/dpkg/lock /var/lib/apt/lists/lock /var/cache/apt/archives/lock"

  # Determine which tool to check lock holders
  local check_cmd=""
  if have fuser; then
    check_cmd="fuser"
  elif have lsof; then
    check_cmd="lsof"
  else
    echo "WARNING: neither fuser nor lsof available, skipping apt lock wait" >&2
    return 0
  fi

  while [ $waited -lt $max_wait ]; do
    local locked=false
    if [ "$check_cmd" = "fuser" ]; then
      sudo fuser $lock_files >/dev/null 2>&1 && locked=true
    else
      sudo lsof $lock_files >/dev/null 2>&1 && locked=true
    fi
    if ! $locked; then
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

# ----------- wait for cloud-init before any package activity -----------
# Strict about fatal errors (exit 1), tolerant of "degraded done" (exit 2).
#
# Degraded means one or more modules hit recoverable errors and the boot still
# completed. A transient dpkg lock collision inside cloud-init's own apt
# activity lands exactly there - which is the very race this wait exists to
# survive. Aborting on it would move the failure from get.docker.com to here
# and dress it in a message that points at cloud-init instead of apt, so the
# bootstrap would fail about as often and be harder to trace.
#
# --long on the degraded path names the module that failed, and that lands in
# job_logs where whoever is triaging will actually see it.
#
# </dev/null because the script arrives on stdin via "bash -s": a command that
# reads stdin would otherwise consume the rest of its own source.
if have cloud-init; then
  echo "waiting for cloud-init to finish..."
  ci_rc=0
  cloud-init status --wait </dev/null || ci_rc=$?
  if [ "$ci_rc" -eq 2 ]; then
    echo "WARNING: cloud-init completed degraded, continuing" >&2
    cloud-init status --long >&2 || true
  elif [ "$ci_rc" -ne 0 ]; then
    exit "$ci_rc"
  fi
fi

# ----------- apt lock patience -----------
# apt's default on a busy dpkg lock is to fail instantly (status 100). A
# timeout makes every apt invocation on the box wait for the lock instead -
# including the apt-get calls inside get.docker.com, which apt_wait_lock
# cannot reach.
#
# Bounded, and sharing apt_wait_lock's budget rather than waiting forever.
# An unbounded value would make that helper's deliberate give-up path dead
# code, since the apt call right after it would block regardless. This file
# also persists on the bastion, so "wait forever" would mean a wedged lock
# hangs every later apt run by a human, with nothing saying why.
if [ "$pm" = "apt" ]; then
  sudo mkdir -p /etc/apt/apt.conf.d
  printf 'DPkg::Lock::Timeout "%s";\n' "$APT_LOCK_WAIT_SECS" \
    | sudo tee /etc/apt/apt.conf.d/90autoglue-lock-timeout >/dev/null

  # This script is re-run from the top after a reboot cuts it off, and a
  # reboot mid-install leaves packages unpacked but unconfigured. apt refuses
  # to do anything until that is finished, so finish it first.
  apt_wait_lock
  sudo DEBIAN_FRONTEND=noninteractive dpkg --configure -a || true
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

# find_sshd locates the daemon binary. It lives in /usr/sbin, which is not on a
# non-root user's PATH, so "command -v sshd" alone finds nothing on exactly the
# hosts this script runs on.
find_sshd() {
  for p in /usr/sbin/sshd /sbin/sshd /usr/local/sbin/sshd; do
    if sudo test -x "$p"; then printf '%s\n' "$p"; return 0; fi
  done
  p="$(command -v sshd 2>/dev/null || true)"
  if [ -n "$p" ]; then printf '%s\n' "$p"; return 0; fi
  return 1
}

# sshd_reload validates the config before applying it, and no longer swallows
# the outcome. A systemctl reload against a config sshd rejects leaves the old
# one running, so the previous "|| true" turned "hardening was rejected" into a
# silent no-op that surfaced at the next reboot.
sshd_reload() {
  sshd_bin="$(find_sshd || true)"
  if [ -z "$sshd_bin" ]; then
    echo "FATAL: cannot locate the sshd binary to validate its config" >&2
    return 1
  fi
  sudo "$sshd_bin" -t

  if ! have systemctl; then
    echo "FATAL: no systemctl available to reload sshd" >&2
    return 1
  fi
  if sudo systemctl reload ssh 2>/dev/null || sudo systemctl reload sshd 2>/dev/null; then
    return 0
  fi
  # Socket-activated sshd (Ubuntu 24.04+) has no long-running process to
  # signal: each connection spawns its own sshd, which reads the config fresh.
  # Nothing to reload is not the same as a failed reload.
  if sudo systemctl is-active --quiet ssh.socket 2>/dev/null; then
    echo "sshd is socket-activated; config applies to new connections without a reload"
    return 0
  fi
  echo "FATAL: could not reload sshd, and it is not socket-activated" >&2
  return 1
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
  # An existing install may just not be started yet after a reboot.
  if have docker && have systemctl; then
    sudo systemctl enable --now docker || true
  fi

  # "have docker" alone is not enough: a reboot during get.docker.com can
  # leave the CLI installed with no working daemon behind it.
  if ! docker_works; then
    if [ "$pm" = "apt" ]; then apt_wait_lock; fi
    curl -fsSL https://get.docker.com | sh
  fi

  if have systemctl; then
    sudo systemctl enable --now docker
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

  # "sudo test -d", not "[ -w ]". The write below runs under sudo, so the
  # calling user's own permission on a root-owned 0755 directory says nothing
  # about whether it will succeed -- and testing the wrong user is why every
  # bastion with a non-root ssh_user reported a clean bootstrap while sshd kept
  # the image defaults. The banner block below only ever tested -d, which is how
  # that one landed on the same hosts where this one silently did not.
  if ! sudo test -d "$confd"; then
    echo "FATAL: $confd does not exist, so drop-in SSH hardening cannot apply." >&2
    echo "       Refusing to report a bootstrapped bastion with sshd unhardened." >&2
    exit 1
  fi

  # ChallengeResponseAuthentication is deliberately absent: it has been a
  # deprecated alias for KbdInteractiveAuthentication since OpenSSH 8.7, and an
  # sshd that rejects it would invalidate this whole file -- taking
  # PasswordAuthentication with it.
  sudo tee "$confd/10-bastion.conf" >/dev/null <<'EOF'
# Bastion hardening. Managed by autoglue; local edits are overwritten.
PasswordAuthentication no
KbdInteractiveAuthentication no
UsePAM yes
PermitEmptyPasswords no
PubkeyAuthentication yes
ClientAliveInterval 300
ClientAliveCountMax 2
LoginGraceTime 20
MaxAuthTries 3
MaxSessions 10
X11Forwarding no
EOF
  sudo chmod 0644 "$confd/10-bastion.conf"
  sshd_reload

  # Assert the effective config, rather than trusting the file just written.
  # sshd takes the FIRST value it obtains across Include'd drop-ins, so an image
  # shipping its own 00-*.conf beats 10-bastion.conf, and a main sshd_config
  # carrying no Include at all reads none of them. Neither is visible from the
  # file on disk, and both look exactly like success.
  sshd_bin="$(find_sshd || true)"
  if [ -z "$sshd_bin" ]; then
    echo "FATAL: cannot locate sshd to verify the hardening took effect" >&2
    exit 1
  fi
  eff="$(sudo "$sshd_bin" -T)"
  for want in "passwordauthentication no" "permitemptypasswords no" \
              "kbdinteractiveauthentication no"; do
    if ! printf '%s\n' "$eff" | grep -qix "$want"; then
      echo "FATAL: SSH hardening did not take effect; expected: $want" >&2
      echo "       got: $(printf '%s\n' "$eff" | grep -i "^${want%% *} " || echo '<unset>')" >&2
      exit 1
    fi
  done
  echo "SSH hardening verified in effect: password auth and agent forwarding are off"

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

# Never report a bastion bootstrapped with docker broken.
if [ "$INSTALL_DOCKER" = "1" ] && ! docker_works; then
  echo "FATAL: docker is installed but not working (sudo docker info failed)" >&2
  exit 1
fi

echo "Bootstrap complete. If you were added to the docker group, log out and back in to apply."
`

	// Send script via stdin to avoid quoting/escaping issues
	sess.Stdin = strings.NewReader(script)

	// Stream combined stdout+stderr rather than buffering to exit, so a
	// bootstrap that hangs on (say) an apt lock is visible while it hangs.
	tail := &tailBuffer{max: logMaxTailBytes}
	var w io.Writer = tail
	if sink != nil {
		w = io.MultiWriter(tail, sink)
	}

	runErr := runSSHStreaming(sess, "bash -s", w)
	if runErr != nil {
		return tail.String(), sessionErr(stopSupervise(), wrapSSHError(runErr, tail.String()))
	}
	return tail.String(), nil
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
