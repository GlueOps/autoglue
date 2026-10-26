package bg

import (
	"context"
	"errors"
	"fmt"
	"io"
	"sync"

	"golang.org/x/crypto/ssh"
)

// tailBuffer keeps only the last max bytes written to it.
//
// Remote commands here can run for hours and emit far more output than is
// worth holding, but the recent tail is exactly what an error message needs.
type tailBuffer struct {
	mu  sync.Mutex
	b   []byte
	max int
}

func (t *tailBuffer) Write(p []byte) (int, error) {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.b = append(t.b, p...)
	if t.max > 0 && len(t.b) > t.max {
		t.b = t.b[len(t.b)-t.max:]
	}
	return len(p), nil
}

func (t *tailBuffer) String() string {
	t.mu.Lock()
	defer t.mu.Unlock()
	return string(t.b)
}

// closeOnDone closes c once ctx is done, and returns a func that cancels this.
//
// Nothing in x/crypto/ssh takes a context: once a handshake or a remote command
// is under way, the only way to stop waiting on it is to close what it is
// reading from. Bind the TCP connection to the step's context right after the
// dial, so the step's deadline bounds the handshake and the command too, even
// against a peer that has stopped answering altogether.
func closeOnDone(ctx context.Context, c io.Closer) (stop func() bool) {
	return context.AfterFunc(ctx, func() { _ = c.Close() })
}

// stepErr reports a remote step that its context cut short as exactly that,
// instead of as the transport error that closing the connection produces. The
// result does not wrap the original error, so it is never mistaken for a
// dropped connection worth retrying.
func stepErr(ctx context.Context, err error) error {
	if err == nil || ctx.Err() == nil {
		return err
	}
	if errors.Is(ctx.Err(), context.DeadlineExceeded) {
		return fmt.Errorf("remote step timed out: %w", context.Cause(ctx))
	}
	return fmt.Errorf("remote step cancelled: %w", context.Cause(ctx))
}

// runSSHStreaming runs cmd and copies its combined output to w as it arrives,
// rather than buffering until exit the way (*ssh.Session).CombinedOutput does.
//
// That difference is the whole point: a `make bootstrap` can run for hours, and
// with CombinedOutput nothing is observable until it finishes — and on success
// the output was discarded entirely.
//
// When ctx is done the session is closed, which ends the wait while the peer
// still answers. A peer that has gone silent will not answer that close either;
// callers bind the connection itself with closeOnDone for that case. This
// stops autoglue waiting on the command, not the command: a container started
// with --sig-proxy=false keeps running on the bastion by design.
func runSSHStreaming(ctx context.Context, sess *ssh.Session, cmd string, w io.Writer) error {
	defer closeOnDone(ctx, sess)()

	stdout, err := sess.StdoutPipe()
	if err != nil {
		return fmt.Errorf("stdout pipe: %w", err)
	}
	stderr, err := sess.StderrPipe()
	if err != nil {
		return fmt.Errorf("stderr pipe: %w", err)
	}

	// Start waits for the peer to acknowledge the exec request, so it can be
	// the call a deadline cuts short.
	if err := sess.Start(cmd); err != nil {
		return stepErr(ctx, fmt.Errorf("start remote command: %w", err))
	}

	// Both pipes must be drained fully before Wait, or Wait can block and the
	// remote side can stall on a full window.
	var wg sync.WaitGroup
	wg.Add(2)

	pump := func(r io.Reader) {
		defer wg.Done()
		_, _ = io.Copy(w, r)
	}
	go pump(stdout)
	go pump(stderr)

	wg.Wait()

	return stepErr(ctx, sess.Wait())
}
