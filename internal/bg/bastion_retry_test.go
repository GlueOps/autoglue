package bg

import (
	"context"
	"crypto/ed25519"
	"crypto/rand"
	"encoding/base64"
	"encoding/binary"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"sync/atomic"
	"syscall"
	"testing"
	"time"

	"github.com/glueops/autoglue/internal/common"
	"github.com/glueops/autoglue/internal/models"
	"github.com/glueops/autoglue/internal/testutil/pgtest"
	"github.com/glueops/autoglue/internal/utils"
	"github.com/riverqueue/river"
	"github.com/riverqueue/river/rivertype"
	"golang.org/x/crypto/ssh"
	"gorm.io/gorm"
)

func refusedDialErr() error {
	return &sshConnectError{fmt.Errorf("dial: %w", &net.OpError{
		Op: "dial", Net: "tcp",
		Err: &os.SyscallError{Syscall: "connect", Err: syscall.ECONNREFUSED},
	})}
}

type timeoutErr struct{}

func (timeoutErr) Error() string   { return "i/o timeout" }
func (timeoutErr) Timeout() bool   { return true }
func (timeoutErr) Temporary() bool { return true }

func TestIsHostNotReady(t *testing.T) {
	cases := []struct {
		name string
		err  error
		want bool
	}{
		{"nil", nil, false},
		{"connection refused", refusedDialErr(), true},
		{"no route to host", &sshConnectError{fmt.Errorf("dial: %w", &net.OpError{Op: "dial", Err: &os.SyscallError{Syscall: "connect", Err: syscall.EHOSTUNREACH}})}, true},
		{"network unreachable", &sshConnectError{fmt.Errorf("dial: %w", syscall.ENETUNREACH)}, true},
		{"connect timed out", &sshConnectError{fmt.Errorf("dial: %w", syscall.ETIMEDOUT)}, true},
		{"dial i/o timeout", &sshConnectError{fmt.Errorf("dial: %w", &net.OpError{Op: "dial", Err: timeoutErr{}})}, true},
		{"handshake EOF", &sshConnectError{fmt.Errorf("ssh handshake: %w", fmt.Errorf("ssh: handshake failed: %w", io.EOF))}, true},
		{"handshake reset", &sshConnectError{fmt.Errorf("ssh handshake: %w", syscall.ECONNRESET)}, true},
		{"auth failure", &sshConnectError{fmt.Errorf("ssh handshake: %w", errors.New("ssh: handshake failed: ssh: unable to authenticate, attempted methods [none publickey]"))}, false},
		{"host key mismatch", &sshConnectError{fmt.Errorf("ssh handshake: %w", errors.New("host key mismatch for x - POSSIBLE MITM or host reinstalled"))}, false},
		{"bad private key", fmt.Errorf("parse private key: %w", errors.New("ssh: no key found")), false},
		// The remote script failing is never "not up yet", even if its error
		// happens to wrap a network errno.
		{"remote script failure", fmt.Errorf("remote run: %w", syscall.ECONNREFUSED), false},
		{"context canceled", &sshConnectError{fmt.Errorf("dial: %w", context.Canceled)}, false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := isHostNotReady(tc.err); got != tc.want {
				t.Errorf("isHostNotReady(%v) = %v, want %v", tc.err, got, tc.want)
			}
		})
	}
}

func TestBastionRetryDelay(t *testing.T) {
	window, sessionLostMax := 10*time.Minute, 3
	authErr := &sshConnectError{errors.New("ssh: handshake failed: ssh: unable to authenticate")}
	lost := &sshSessionError{&ssh.ExitMissingError{}}
	m := func(snoozes, sessionLost int) bastionJobMeta {
		return bastionJobMeta{Snoozes: snoozes, SessionLost: sessionLost}
	}
	cases := []struct {
		name      string
		err       error
		waited    time.Duration
		meta      bastionJobMeta
		wantDelay time.Duration
		wantRetry bool
	}{
		{"first refusal", refusedDialErr(), 2 * time.Second, m(0, 0), 5 * time.Second, true},
		{"backs off", refusedDialErr(), time.Minute, m(2, 0), 20 * time.Second, true},
		{"capped", refusedDialErr(), 5 * time.Minute, m(4, 0), time.Minute, true},
		{"capped with large count", refusedDialErr(), 9 * time.Minute, m(1000, 0), time.Minute, true},
		{"window exceeded", refusedDialErr(), window, m(3, 0), 0, false},
		{"auth failure in window", authErr, time.Second, m(0, 0), 0, false},
		{"script failure in window", errors.New("remote run: exit status 1"), time.Second, m(0, 0), 0, false},
		{"session lost in window", lost, 4 * time.Minute, m(1, 0), 10 * time.Second, true},
		// The reboot after cloud-init's package upgrade lands late.
		{"session lost at minute 11", lost, 11 * time.Minute, m(1, 0), 10 * time.Second, true},
		{"third session loss", lost, time.Minute, m(2, 2), 20 * time.Second, true},
		{"fourth session loss", lost, time.Minute, m(3, 3), 0, false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			delay, retry := bastionRetryDelay(tc.err, tc.waited, tc.meta, window, sessionLostMax)
			if retry != tc.wantRetry || delay != tc.wantDelay {
				t.Errorf("bastionRetryDelay = (%s, %v), want (%s, %v)", delay, retry, tc.wantDelay, tc.wantRetry)
			}
		})
	}
}

func TestIsSessionLost(t *testing.T) {
	sess := func(err error) error { return &sshSessionError{fmt.Errorf("remote run: %w", err)} }
	cases := []struct {
		name string
		err  error
		want bool
	}{
		{"nil", nil, false},
		{"exit status missing", sess(&ssh.ExitMissingError{}), true},
		{"session EOF", &sshSessionError{fmt.Errorf("session: %w", io.EOF)}, true},
		{"unexpected EOF", sess(io.ErrUnexpectedEOF), true},
		{"connection reset", sess(&net.OpError{Op: "read", Err: &os.SyscallError{Syscall: "read", Err: syscall.ECONNRESET}}), true},
		{"broken pipe", sess(syscall.EPIPE), true},
		{"other run error", sess(errors.New("start remote command: boom")), false},
		// Missing exit status is only a lost session after the handshake.
		{"connect-phase EOF", &sshConnectError{io.EOF}, false},
		{"untagged", &ssh.ExitMissingError{}, false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := isSessionLost(tc.err); got != tc.want {
				t.Errorf("isSessionLost(%v) = %v, want %v", tc.err, got, tc.want)
			}
		})
	}
}

func TestParseBastionMeta(t *testing.T) {
	cases := map[string]bastionJobMeta{
		``:              {},
		`{}`:            {},
		`{"snoozes":3}`: {Snoozes: 3},
		`{"snoozes":2,"other":"x","bastion_session_lost":1,"bastion_session_lost_at":42}`: {
			Snoozes: 2, SessionLost: 1, SessionLostAt: 42,
		},
		`not json`: {},
	}
	for in, want := range cases {
		got := parseBastionMeta([]byte(in))
		if got != want {
			t.Errorf("parseBastionMeta(%q) = %+v, want %+v", in, got, want)
		}
	}
}

// ----- end to end against a fake sshd -----

// How a fake sshd ends an exec once it has read the script.
const (
	execOK     = "ok"     // exit-status 0
	execExit1  = "exit1"  // exit-status 1: the script ran and failed
	execDrop   = "drop"   // close the TCP connection: the host went away
	execNoExit = "noexit" // close the channel with no exit-status
)

// fakeSSHD accepts keys for which accept returns true, runs nothing, and ends
// each exec as exit says (nil means execOK).
func fakeSSHD(t *testing.T, ln net.Listener, accept func(ssh.PublicKey) bool, exit func() string) {
	t.Helper()
	_, hostPriv, _ := ed25519.GenerateKey(rand.Reader)
	hostSigner, err := ssh.NewSignerFromKey(hostPriv)
	if err != nil {
		t.Fatalf("host signer: %v", err)
	}
	cfg := &ssh.ServerConfig{
		PublicKeyCallback: func(_ ssh.ConnMetadata, k ssh.PublicKey) (*ssh.Permissions, error) {
			if accept(k) {
				return nil, nil
			}
			return nil, errors.New("unauthorized")
		},
	}
	cfg.AddHostKey(hostSigner)

	go func() {
		for {
			c, err := ln.Accept()
			if err != nil {
				return
			}
			go serveFakeSSH(c, cfg, exit)
		}
	}()
}

func onlyKey(want ssh.PublicKey) func(ssh.PublicKey) bool {
	return func(k ssh.PublicKey) bool { return string(k.Marshal()) == string(want.Marshal()) }
}

func serveFakeSSH(c net.Conn, cfg *ssh.ServerConfig, exit func() string) {
	defer func() { _ = c.Close() }()
	_, chans, reqs, err := ssh.NewServerConn(c, cfg)
	if err != nil {
		return
	}
	go ssh.DiscardRequests(reqs)
	for nc := range chans {
		ch, creqs, err := nc.Accept()
		if err != nil {
			continue
		}
		go func() {
			defer func() { _ = ch.Close() }()
			for r := range creqs {
				_ = r.Reply(r.Type == "exec", nil)
				if r.Type != "exec" {
					continue
				}
				_, _ = io.Copy(io.Discard, ch) // the script, until stdin closes
				mode := execOK
				if exit != nil {
					mode = exit()
				}
				switch mode {
				case execDrop:
					_ = c.Close()
				case execNoExit:
				default:
					status := make([]byte, 4)
					if mode == execExit1 {
						binary.BigEndian.PutUint32(status, 1)
					}
					_, _ = ch.SendRequest("exit-status", false, status)
				}
				return
			}
		}()
	}
}

type bastionFixture struct {
	db     *gorm.DB
	server models.Server
	pub    ssh.PublicKey
}

func newBastionFixture(t *testing.T) bastionFixture {
	t.Helper()
	db := pgtest.DB(t)

	var n int64
	db.Model(&models.MasterKey{}).Where("is_active = ?", true).Count(&n)
	if n == 0 {
		mk := make([]byte, 32)
		_, _ = rand.Read(mk)
		if err := db.Create(&models.MasterKey{Key: base64.StdEncoding.EncodeToString(mk), IsActive: true}).Error; err != nil {
			t.Fatalf("create master key: %v", err)
		}
	}

	org := models.Organization{Name: "bastion-retry-" + t.Name()}
	if err := db.Create(&org).Error; err != nil {
		t.Fatalf("create org: %v", err)
	}

	pubRaw, priv, _ := ed25519.GenerateKey(rand.Reader)
	block, err := ssh.MarshalPrivateKey(priv, "")
	if err != nil {
		t.Fatalf("marshal key: %v", err)
	}
	ct, iv, tag, err := utils.EncryptForOrg(org.ID, pem.EncodeToMemory(block), db)
	if err != nil {
		t.Fatalf("encrypt key: %v", err)
	}
	pub, _ := ssh.NewPublicKey(pubRaw)
	key := models.SshKey{
		AuditFields:         common.AuditFields{OrganizationID: org.ID},
		Name:                "k-" + t.Name(),
		PublicKey:           string(ssh.MarshalAuthorizedKey(pub)),
		EncryptedPrivateKey: ct,
		PrivateIV:           iv,
		PrivateTag:          tag,
		Fingerprint:         ssh.FingerprintSHA256(pub) + t.Name(),
	}
	if err := db.Create(&key).Error; err != nil {
		t.Fatalf("create ssh key: %v", err)
	}

	ip := "127.0.0.1"
	s := models.Server{
		OrganizationID:   org.ID,
		Hostname:         "bastion",
		PublicIPAddress:  &ip,
		PrivateIPAddress: "10.0.0.1",
		SSHUser:          "ubuntu",
		SshKeyID:         key.ID,
		Role:             "bastion",
		Status:           "provisioning",
	}
	if err := db.Create(&s).Error; err != nil {
		t.Fatalf("create server: %v", err)
	}
	return bastionFixture{db: db, server: s, pub: pub}
}

// work runs one attempt directly, with snoozes as the job's River metadata.
func (f bastionFixture) work(t *testing.T, createdAt time.Time, snoozes int) error {
	t.Helper()
	return f.workMeta(t, createdAt, fmt.Sprintf(`{"snoozes":%d}`, snoozes))
}

func (f bastionFixture) workMeta(t *testing.T, createdAt time.Time, meta string) error {
	t.Helper()
	w := &BastionBootstrapWorker{db: f.db}
	j := &river.Job[BastionBootstrapArgs]{
		JobRow: &rivertype.JobRow{
			ID:        time.Now().UnixNano(),
			CreatedAt: createdAt,
			Metadata:  []byte(meta),
		},
		Args: BastionBootstrapArgs{ServerID: f.server.ID},
	}
	return w.Work(context.Background(), j)
}

func (f bastionFixture) status(t *testing.T) string {
	t.Helper()
	var s models.Server
	if err := f.db.First(&s, "id = ?", f.server.ID).Error; err != nil {
		t.Fatalf("load server: %v", err)
	}
	return s.Status
}

// freePort returns a port with nothing listening on it, so a dial is refused.
func freePort(t *testing.T) string {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	port := strconv.Itoa(ln.Addr().(*net.TCPAddr).Port)
	_ = ln.Close()
	return port
}

func usePort(t *testing.T, port string) {
	old := bastionSSHPort
	bastionSSHPort = port
	t.Cleanup(func() { bastionSSHPort = old })
}

func wantSnooze(t *testing.T, err error) time.Duration {
	t.Helper()
	var se *river.JobSnoozeError
	if !errors.As(err, &se) {
		t.Fatalf("Work returned %v, want a snooze", err)
	}
	return se.Duration
}

// The bug this guards: a bastion claimed before its sshd is up used to be
// marked failed on the first refused connection and never retried.
func TestBastionBootstrapWaitsForSSHThenSucceeds(t *testing.T) {
	f := newBastionFixture(t)
	port := freePort(t)
	usePort(t, port)
	created := time.Now()

	for i := 0; i < 2; i++ {
		d := wantSnooze(t, f.work(t, created, i))
		if want := bastionRetryInitial << i; d != want {
			t.Errorf("snooze %d = %s, want %s", i, d, want)
		}
		if got := f.status(t); got != "provisioning" {
			t.Fatalf("status after refused attempt %d = %q, want provisioning", i, got)
		}
	}

	ln, err := net.Listen("tcp", "127.0.0.1:"+port)
	if err != nil {
		t.Fatalf("listen on %s: %v", port, err)
	}
	t.Cleanup(func() { _ = ln.Close() })
	fakeSSHD(t, ln, onlyKey(f.pub), nil)

	if err := f.work(t, created, 2); err != nil {
		t.Fatalf("Work once sshd is up = %v, want nil", err)
	}
	if got := f.status(t); got != "ready" {
		t.Fatalf("status = %q, want ready", got)
	}
}

// listenFakeSSHD starts a fake sshd on a free port and points the bootstrap
// at it.
func listenFakeSSHD(t *testing.T, accept func(ssh.PublicKey) bool, exit func() string) {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	t.Cleanup(func() { _ = ln.Close() })
	fakeSSHD(t, ln, accept, exit)
	usePort(t, strconv.Itoa(ln.Addr().(*net.TCPAddr).Port))
}

func TestBastionBootstrapAuthFailureFailsImmediately(t *testing.T) {
	f := newBastionFixture(t)
	other, _, _ := ed25519.GenerateKey(rand.Reader)
	otherPub, _ := ssh.NewPublicKey(other)
	listenFakeSSHD(t, onlyKey(otherPub), nil) // sshd up, but our key is not authorized

	if err := f.work(t, time.Now(), 0); err != nil {
		t.Fatalf("Work = %v, want nil (failure recorded on the server, not retried)", err)
	}
	if got := f.status(t); got != "failed" {
		t.Fatalf("status = %q, want failed", got)
	}
}

func TestBastionBootstrapFailsAfterWindow(t *testing.T) {
	f := newBastionFixture(t)
	usePort(t, freePort(t))

	if err := f.work(t, time.Now().Add(-bastionSSHWait()-time.Second), 7); err != nil {
		t.Fatalf("Work = %v, want nil (failure recorded on the server, not retried)", err)
	}
	if got := f.status(t); got != "failed" {
		t.Fatalf("status = %q, want failed", got)
	}
}

// TestIsSessionLostAgainstRealExits checks the classification against what
// x/crypto really returns for each way a session can end, rather than against
// hand-built errors. One fake sshd serves every case so the host key stays the
// one TOFU recorded first.
func TestIsSessionLostAgainstRealExits(t *testing.T) {
	f := newBastionFixture(t)
	var mode atomic.Value
	listenFakeSSHD(t, onlyKey(f.pub), func() string { return mode.Load().(string) })

	var s models.Server
	if err := f.db.Preload("SshKey").First(&s, "id = ?", f.server.ID).Error; err != nil {
		t.Fatalf("load server: %v", err)
	}
	priv, err := utils.DecryptForOrg(s.OrganizationID, s.SshKey.EncryptedPrivateKey, s.SshKey.PrivateIV, s.SshKey.PrivateTag, f.db)
	if err != nil {
		t.Fatalf("decrypt: %v", err)
	}
	run := func() error {
		_, err := sshInstallDockerWithOutput(context.Background(), f.db, &s,
			net.JoinHostPort("127.0.0.1", bastionSSHPort), s.SSHUser, []byte(priv), nil)
		return err
	}

	cases := []struct {
		mode string
		want bool
	}{
		{execExit1, false},
		{execDrop, true},
		{execNoExit, true},
	}
	mode.Store(execOK)
	if err := run(); err != nil {
		t.Fatalf("clean exit returned %v", err)
	}
	for _, tc := range cases {
		t.Run(tc.mode, func(t *testing.T) {
			mode.Store(tc.mode)
			err := run()
			var se *sshSessionError
			if !errors.As(err, &se) {
				t.Fatalf("mode %s returned %v, want a session error", tc.mode, err)
			}
			if got := isSessionLost(err); got != tc.want {
				t.Errorf("isSessionLost(%v) = %v, want %v", err, got, tc.want)
			}
		})
	}
}

// The host reboots under the script (cloud-init's package_reboot_if_required
// while it sits in "cloud-init status --wait"); the next attempt finds it back.
func TestBastionBootstrapRetriesAfterConnectionLost(t *testing.T) {
	f := newBastionFixture(t)
	var rebooted atomic.Bool
	listenFakeSSHD(t, onlyKey(f.pub), func() string {
		if rebooted.CompareAndSwap(false, true) {
			return execDrop
		}
		return execOK
	})
	created := time.Now()

	wantSnooze(t, f.work(t, created, 0))
	if got := f.status(t); got != "provisioning" {
		t.Fatalf("status after lost connection = %q, want provisioning", got)
	}

	if err := f.work(t, created, 1); err != nil {
		t.Fatalf("Work after reboot = %v, want nil", err)
	}
	if got := f.status(t); got != "ready" {
		t.Fatalf("status = %q, want ready", got)
	}
}

// A script that ran and exited non-zero is a real failure, not a reboot.
func TestBastionBootstrapScriptFailureFailsImmediately(t *testing.T) {
	f := newBastionFixture(t)
	listenFakeSSHD(t, onlyKey(f.pub), func() string { return execExit1 })

	if err := f.work(t, time.Now(), 0); err != nil {
		t.Fatalf("Work = %v, want nil (failure recorded on the server, not retried)", err)
	}
	if got := f.status(t); got != "failed" {
		t.Fatalf("status = %q, want failed", got)
	}
}

// The Proxmox reboot after cloud-init's package upgrade lands after the
// script has sat in "cloud-init status --wait" for the whole upgrade, often
// past the wait window. It must still be retried.
func TestBastionBootstrapSessionLostPastWindowIsRetried(t *testing.T) {
	f := newBastionFixture(t)
	listenFakeSSHD(t, onlyKey(f.pub), func() string { return execDrop })

	wantSnooze(t, f.work(t, time.Now().Add(-11*time.Minute), 0))
	if got := f.status(t); got != "provisioning" {
		t.Fatalf("status = %q, want provisioning", got)
	}
}

func TestBastionBootstrapFourthSessionLossFails(t *testing.T) {
	f := newBastionFixture(t)
	listenFakeSSHD(t, onlyKey(f.pub), func() string { return execDrop })

	err := f.workMeta(t, time.Now(), `{"snoozes":3,"bastion_session_lost":3}`)
	if err != nil {
		t.Fatalf("Work = %v, want nil (failure recorded on the server, not retried)", err)
	}
	if got := f.status(t); got != "failed" {
		t.Fatalf("status = %q, want failed", got)
	}
}
