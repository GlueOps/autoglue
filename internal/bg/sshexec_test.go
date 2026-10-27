package bg

import (
	"context"
	"crypto/ed25519"
	"crypto/rand"
	"errors"
	"io"
	"net"
	"strings"
	"sync"
	"testing"
	"time"

	"golang.org/x/crypto/ssh"
)

// hangingSSHD is an in-process sshd whose exec requests never finish: it
// acknowledges the command, writes a line, and then neither exits nor closes
// the channel — a remote `make bootstrap` wedged on a lock, as seen from the
// client. With finish set, the command instead exits 0 after its line.
type hangingSSHD struct {
	ln     net.Listener
	finish bool

	mu    sync.Mutex
	conns []*freezableConn
}

func startHangingSSHD(t *testing.T) *hangingSSHD {
	t.Helper()
	return startTestSSHD(t, false)
}

func startTestSSHD(t *testing.T, finish bool) *hangingSSHD {
	t.Helper()

	_, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	signer, err := ssh.NewSignerFromKey(priv)
	if err != nil {
		t.Fatal(err)
	}
	cfg := &ssh.ServerConfig{NoClientAuth: true}
	cfg.AddHostKey(signer)

	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	d := &hangingSSHD{ln: ln, finish: finish}
	t.Cleanup(d.close)

	go func() {
		for {
			raw, err := ln.Accept()
			if err != nil {
				return
			}
			fc := newFreezableConn(raw)
			d.mu.Lock()
			d.conns = append(d.conns, fc)
			d.mu.Unlock()
			go d.serve(fc, cfg)
		}
	}()
	return d
}

func (d *hangingSSHD) serve(c net.Conn, cfg *ssh.ServerConfig) {
	_, chans, reqs, err := ssh.NewServerConn(c, cfg)
	if err != nil {
		return
	}
	go ssh.DiscardRequests(reqs)
	for nc := range chans {
		if nc.ChannelType() != "session" {
			_ = nc.Reject(ssh.UnknownChannelType, "session only")
			continue
		}
		ch, chReqs, err := nc.Accept()
		if err != nil {
			continue
		}
		go func() {
			for req := range chReqs {
				if req.Type != "exec" {
					_ = req.Reply(false, nil)
					continue
				}
				_ = req.Reply(true, nil)
				_, _ = io.WriteString(ch, "started\n")
				if d.finish {
					_, _ = ch.SendRequest("exit-status", false, ssh.Marshal(struct{ Status uint32 }{0}))
					_ = ch.Close()
				}
				// Otherwise no exit-status, no close: the command never finishes.
			}
		}()
	}
}

// freeze makes the server stop reading from every connection, so it no longer
// answers anything, including a channel close: a peer that has gone silent
// with the TCP connection still up.
func (d *hangingSSHD) freeze() {
	d.mu.Lock()
	defer d.mu.Unlock()
	for _, c := range d.conns {
		c.freeze()
	}
}

func (d *hangingSSHD) close() {
	_ = d.ln.Close()
	d.mu.Lock()
	defer d.mu.Unlock()
	for _, c := range d.conns {
		_ = c.Close()
	}
}

type freezableConn struct {
	net.Conn
	frozen    chan struct{}
	closed    chan struct{}
	freezeOne sync.Once
	closeOne  sync.Once
}

func newFreezableConn(c net.Conn) *freezableConn {
	return &freezableConn{Conn: c, frozen: make(chan struct{}), closed: make(chan struct{})}
}

func (c *freezableConn) freeze() { c.freezeOne.Do(func() { close(c.frozen) }) }

// Read checks for a freeze after the underlying read returns, not before: the
// server's read loop is normally already blocked inside a Read when freeze is
// called, and whatever that read brings back (the client's channel close, say)
// must be dropped unseen, or the "silent" peer would answer it.
func (c *freezableConn) Read(p []byte) (int, error) {
	n, err := c.Conn.Read(p)
	select {
	case <-c.frozen:
		<-c.closed
		return 0, net.ErrClosed
	default:
		return n, err
	}
}

func (c *freezableConn) Close() error {
	c.closeOne.Do(func() { close(c.closed) })
	return c.Conn.Close()
}

// dialTestSSH connects the way the remote-step callers do: a context-bounded
// TCP dial followed by an explicit handshake. With bindConn the connection is
// also tied to the step's context, as the callers now do with closeOnDone.
func dialTestSSH(t *testing.T, ctx context.Context, addr string, bindConn bool) *ssh.Session {
	t.Helper()

	conn, err := (&net.Dialer{}).DialContext(ctx, "tcp", addr)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = conn.Close() })
	if bindConn {
		stop := closeOnDone(ctx, conn)
		t.Cleanup(func() { stop() })
	}

	cc, chans, reqs, err := ssh.NewClientConn(conn, addr, &ssh.ClientConfig{
		User:            "test",
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	})
	if err != nil {
		t.Fatal(err)
	}
	client := ssh.NewClient(cc, chans, reqs)
	t.Cleanup(func() { _ = client.Close() })

	sess, err := client.NewSession()
	if err != nil {
		t.Fatal(err)
	}
	return sess
}

// waitStep runs fn and fails the test if it is still blocked well after the
// step deadline.
func waitStep(t *testing.T, fn func() error) error {
	t.Helper()
	done := make(chan error, 1)
	go func() { done <- fn() }()
	select {
	case err := <-done:
		return err
	case <-time.After(3 * time.Second):
		t.Fatal("still blocked 3s after the step deadline passed")
		return nil
	}
}

func assertTimedOut(t *testing.T, err error) {
	t.Helper()
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("want an error wrapping context.DeadlineExceeded, got %v", err)
	}
	if !strings.Contains(err.Error(), "timed out") {
		t.Fatalf("want the error to say the step timed out, got %q", err)
	}
}

// A remote step whose command hangs must give up once its deadline passes
// (#1080). Before the fix the context reached the dial and stopped there, so
// the call blocked until the remote side exited — for a wedged command, never.
//
// The connection is deliberately not bound here: against a peer that still
// answers, runSSHStreaming has to stop on its own.
func TestRunSSHStreamingStopsAtStepDeadline(t *testing.T) {
	d := startHangingSSHD(t)

	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()
	sess := dialTestSSH(t, ctx, d.ln.Addr().String(), false)

	var out strings.Builder
	err := waitStep(t, func() error { return runSSHStreaming(ctx, sess, "make bootstrap", &out) })
	assertTimedOut(t, err)
	if !strings.Contains(out.String(), "started") {
		t.Fatalf("output before the deadline should still be streamed, got %q", out.String())
	}
}

// A peer that stops answering mid-command does not acknowledge the session
// close either. Binding the connection to the step's context is what ends the
// wait there.
func TestRunSSHStreamingStopsWhenPeerGoesSilent(t *testing.T) {
	d := startHangingSSHD(t)

	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()
	sess := dialTestSSH(t, ctx, d.ln.Addr().String(), true)

	err := waitStep(t, func() error {
		// Freeze once the command is running, not before the handshake.
		go func() { time.Sleep(100 * time.Millisecond); d.freeze() }()
		return runSSHStreaming(ctx, sess, "make bootstrap", io.Discard)
	})
	assertTimedOut(t, err)
}

// A peer that goes silent before acknowledging the command leaves Start, not
// Wait, as the call that blocks. That too must end at the deadline, and be
// reported as a timeout rather than as the EOF that closing the connection
// produces.
func TestRunSSHStreamingStopsWhenPeerGoesSilentBeforeStart(t *testing.T) {
	d := startHangingSSHD(t)

	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()
	sess := dialTestSSH(t, ctx, d.ln.Addr().String(), true)
	d.freeze()

	err := waitStep(t, func() error { return runSSHStreaming(ctx, sess, "make bootstrap", io.Discard) })
	assertTimedOut(t, err)
}

// The handshake runs outside any context in x/crypto/ssh. A host that accepts
// TCP and then never speaks must not outlive the step's deadline.
func TestStepDeadlineBoundsHandshake(t *testing.T) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = ln.Close() })
	go func() {
		for {
			c, err := ln.Accept()
			if err != nil {
				return
			}
			t.Cleanup(func() { _ = c.Close() })
		}
	}()

	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()

	conn, err := (&net.Dialer{}).DialContext(ctx, "tcp", ln.Addr().String())
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()
	defer closeOnDone(ctx, conn)()

	err = waitStep(t, func() error {
		_, _, _, err := ssh.NewClientConn(conn, ln.Addr().String(), &ssh.ClientConfig{
			User:            "test",
			HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		})
		return err
	})
	if err == nil {
		t.Fatal("handshake against a silent host should fail")
	}
}

// A command that finishes before its deadline is unaffected: no error, full
// output, and the connection is not closed out from under it afterwards.
func TestRunSSHStreamingCompletesWithinDeadline(t *testing.T) {
	d := startTestSSHD(t, true)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	sess := dialTestSSH(t, ctx, d.ln.Addr().String(), true)

	var out strings.Builder
	if err := runSSHStreaming(ctx, sess, "true", &out); err != nil {
		t.Fatalf("want success, got %v", err)
	}
	if out.String() != "started\n" {
		t.Fatalf("got output %q", out.String())
	}
}

// A timed-out step is reported as one, not by whatever its output last
// mentioned.
func TestWrapSSHErrorReportsTimeoutOverOutput(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Nanosecond)
	defer cancel()
	<-ctx.Done()

	err := wrapSSHError(stepErr(ctx, io.EOF), "E: Could not get lock /var/lib/dpkg/lock (apt-get)")
	assertTimedOut(t, err)
	if strings.Contains(err.Error(), "apt failed") {
		t.Fatalf("timeout misreported as an apt failure: %q", err)
	}
	if errors.Is(err, io.EOF) {
		t.Fatal("a timed-out step should not wrap the transport error it caused")
	}
}

// A step cut short by its deadline must not be retried as a lost session or a
// host still booting (#1245): closing the connection produces exactly the EOF
// both of those look for.
func TestTimedOutStepIsNotRetried(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Nanosecond)
	defer cancel()
	<-ctx.Done()

	for name, err := range map[string]error{
		"run":     &sshSessionError{wrapSSHError(stepErr(ctx, io.EOF), "")},
		"session": stepErr(ctx, io.EOF),
	} {
		if isSessionLost(err) || isHostNotReady(err) {
			t.Errorf("%s: timed-out step classified as retryable: %v", name, err)
		}
	}
}
