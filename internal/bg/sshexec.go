package bg

import (
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
