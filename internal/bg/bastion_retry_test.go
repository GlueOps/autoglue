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
	window := 10 * time.Minute
	authErr := &sshConnectError{errors.New("ssh: handshake failed: ssh: unable to authenticate")}
	cases := []struct {
		name      string
		err       error
		waited    time.Duration
		snoozes   int
		wantDelay time.Duration
		wantRetry bool
	}{
		{"first refusal", refusedDialErr(), 2 * time.Second, 0, 5 * time.Second, true},
		{"backs off", refusedDialErr(), time.Minute, 2, 20 * time.Second, true},
		{"capped", refusedDialErr(), 5 * time.Minute, 4, time.Minute, true},
		{"capped with large count", refusedDialErr(), 9 * time.Minute, 1000, time.Minute, true},
		{"window exceeded", refusedDialErr(), window, 3, 0, false},
		{"auth failure in window", authErr, time.Second, 0, 0, false},
		{"script failure in window", errors.New("remote run: exit status 1"), time.Second, 0, 0, false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			delay, retry := bastionRetryDelay(tc.err, tc.waited, tc.snoozes, window)
			if retry != tc.wantRetry || delay != tc.wantDelay {
				t.Errorf("bastionRetryDelay = (%s, %v), want (%s, %v)", delay, retry, tc.wantDelay, tc.wantRetry)
			}
		})
	}
}

func TestBastionSnoozes(t *testing.T) {
	cases := map[string]int{
		``:                          0,
		`{}`:                        0,
		`{"snoozes":3}`:             3,
		`{"snoozes":2,"other":"x"}`: 2,
		`not json`:                  0,
	}
	for in, want := range cases {
		if got := bastionSnoozes([]byte(in)); got != want {
			t.Errorf("bastionSnoozes(%q) = %d, want %d", in, got, want)
		}
	}
}

// ----- end to end against a fake sshd -----

// fakeSSHD accepts only authorized, runs nothing, and exits 0 for any exec.
func fakeSSHD(t *testing.T, ln net.Listener, authorized ssh.PublicKey) {
	t.Helper()
	_, hostPriv, _ := ed25519.GenerateKey(rand.Reader)
	hostSigner, err := ssh.NewSignerFromKey(hostPriv)
	if err != nil {
		t.Fatalf("host signer: %v", err)
	}
	cfg := &ssh.ServerConfig{
		PublicKeyCallback: func(_ ssh.ConnMetadata, k ssh.PublicKey) (*ssh.Permissions, error) {
			if string(k.Marshal()) == string(authorized.Marshal()) {
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
			go serveFakeSSH(c, cfg)
		}
	}()
}

func serveFakeSSH(c net.Conn, cfg *ssh.ServerConfig) {
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
				status := make([]byte, 4)
				binary.BigEndian.PutUint32(status, 0)
				_, _ = ch.SendRequest("exit-status", false, status)
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

func (f bastionFixture) work(t *testing.T, createdAt time.Time, snoozes int) error {
	t.Helper()
	w := &BastionBootstrapWorker{db: f.db}
	j := &river.Job[BastionBootstrapArgs]{
		JobRow: &rivertype.JobRow{
			ID:        time.Now().UnixNano(),
			CreatedAt: createdAt,
			Metadata:  []byte(fmt.Sprintf(`{"snoozes":%d}`, snoozes)),
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
	fakeSSHD(t, ln, f.pub)

	if err := f.work(t, created, 2); err != nil {
		t.Fatalf("Work once sshd is up = %v, want nil", err)
	}
	if got := f.status(t); got != "ready" {
		t.Fatalf("status = %q, want ready", got)
	}
}

func TestBastionBootstrapAuthFailureFailsImmediately(t *testing.T) {
	f := newBastionFixture(t)
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	t.Cleanup(func() { _ = ln.Close() })
	other, _, _ := ed25519.GenerateKey(rand.Reader)
	otherPub, _ := ssh.NewPublicKey(other)
	fakeSSHD(t, ln, otherPub) // sshd up, but our key is not authorized
	usePort(t, strconv.Itoa(ln.Addr().(*net.TCPAddr).Port))

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
