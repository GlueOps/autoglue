package bg

import (
	"context"
	"errors"
	"fmt"
	"io"
	"sync"
	"time"

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

// runSSHStreaming runs cmd and copies its combined output to w as it arrives,
// rather than buffering until exit the way (*ssh.Session).CombinedOutput does.
//
// That difference is the whole point: a `make bootstrap` can run for hours, and
// with CombinedOutput nothing is observable until it finishes — and on success
// the output was discarded entirely.
func runSSHStreaming(sess *ssh.Session, cmd string, w io.Writer) error {
	stdout, err := sess.StdoutPipe()
	if err != nil {
		return fmt.Errorf("stdout pipe: %w", err)
	}
	stderr, err := sess.StderrPipe()
	if err != nil {
		return fmt.Errorf("stderr pipe: %w", err)
	}

	if err := sess.Start(cmd); err != nil {
		return fmt.Errorf("start remote command: %w", err)
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

	return sess.Wait()
}

// Keepalive cadence for superviseSSH. Vars only so tests can shrink them.
var (
	sshKeepaliveEvery  = 15 * time.Second
	sshKeepaliveMisses = 3
)

var errKeepaliveLost = errors.New("ssh keepalive: host stopped answering")

// superviseSSH closes client when ctx is done, or when the host misses
// sshKeepaliveMisses keepalives in a row. Closing it is what unblocks a
// session stuck reading from a peer that is gone without a FIN.
//
// The returned stop ends supervision and reports why it closed the client:
// nil if it did not, ctx.Err() or errKeepaliveLost if it did. It is safe to
// call more than once.
func superviseSSH(ctx context.Context, client *ssh.Client) (stop func() error) {
	done := make(chan struct{})
	var (
		once   sync.Once
		mu     sync.Mutex
		reason error
	)
	closeFor := func(err error) {
		mu.Lock()
		if reason == nil {
			reason = err
		}
		mu.Unlock()
		_ = client.Close()
	}

	go func() {
		t := time.NewTicker(sshKeepaliveEvery)
		defer t.Stop()
		misses := 0
		for {
			select {
			case <-done:
				return
			case <-ctx.Done():
				closeFor(ctx.Err())
				return
			case <-t.C:
			}

			// Any reply, including "unsupported", proves the host is there.
			// An error means the connection is already gone, which the
			// session will report on its own.
			reply := make(chan error, 1)
			go func() {
				_, _, err := client.SendRequest("keepalive@openssh.com", true, nil)
				reply <- err
			}()
			select {
			case <-done:
				return
			case <-ctx.Done():
				closeFor(ctx.Err())
				return
			case err := <-reply:
				if err != nil {
					return
				}
				misses = 0
			case <-time.After(sshKeepaliveEvery):
				if misses++; misses >= sshKeepaliveMisses {
					closeFor(errKeepaliveLost)
					return
				}
			}
		}
	}()

	return func() error {
		once.Do(func() { close(done) })
		mu.Lock()
		defer mu.Unlock()
		return reason
	}
}
