package bg

import (
	"fmt"
	"strings"
	"sync"
	"testing"
)

func TestTailBufferKeepsOnlyTheEnd(t *testing.T) {
	tb := &tailBuffer{max: 10}

	for i := 0; i < 100; i++ {
		if _, err := tb.Write([]byte("0123456789")); err != nil {
			t.Fatalf("write: %v", err)
		}
	}

	got := tb.String()
	if len(got) != 10 {
		t.Fatalf("tail length = %d, want 10", len(got))
	}
	if got != "0123456789" {
		t.Errorf("tail = %q, want the most recent bytes", got)
	}
}

func TestTailBufferUnboundedWhenMaxIsZero(t *testing.T) {
	tb := &tailBuffer{}
	for i := 0; i < 5; i++ {
		_, _ = tb.Write([]byte("ab"))
	}
	if got := tb.String(); got != "ababababab" {
		t.Errorf("tail = %q, want everything retained", got)
	}
}

func TestTailBufferConcurrentWritesDoNotRace(t *testing.T) {
	// stdout and stderr are pumped by separate goroutines into the same writer,
	// so this is the real access pattern. Run with -race to mean anything.
	tb := &tailBuffer{max: 1 << 16}

	var wg sync.WaitGroup
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				_, _ = tb.Write([]byte(fmt.Sprintf("line-%d-%d\n", n, j)))
			}
		}(i)
	}
	wg.Wait()

	if got := strings.Count(tb.String(), "\n"); got != 800 {
		t.Errorf("line count = %d, want 800", got)
	}
}
