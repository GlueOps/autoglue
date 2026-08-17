package store

import (
	"context"
	"errors"
	"path/filepath"
	"testing"
)

func open(t *testing.T) (*Store, context.Context) {
	t.Helper()
	ctx := context.Background()
	s, err := Open(ctx, filepath.Join(t.TempDir(), "agent.db"))
	if err != nil {
		t.Fatalf("open store: %v", err)
	}
	t.Cleanup(func() { _ = s.Close() })
	return s, ctx
}

// The invariant the whole task plane rests on. Enforced by a partial unique
// index rather than by the code that calls AcceptTask, so a future caller that
// forgets to check still cannot break it.
func TestOnlyOneTaskInFlight(t *testing.T) {
	s, ctx := open(t)

	if err := s.AcceptTask(ctx, Task{ID: "t1", Command: "make"}); err != nil {
		t.Fatalf("accept first: %v", err)
	}

	err := s.AcceptTask(ctx, Task{ID: "t2", Command: "make"})
	if !errors.Is(err, ErrTaskInFlight) {
		t.Fatalf("accepting a second task = %v, want ErrTaskInFlight", err)
	}

	// The slot only frees on a terminal outcome. Starting is not finishing.
	if err := s.StartTask(ctx, "t1", "container-abc"); err != nil {
		t.Fatalf("start: %v", err)
	}
	if err := s.AcceptTask(ctx, Task{ID: "t2"}); !errors.Is(err, ErrTaskInFlight) {
		t.Fatalf("accepting while started = %v, want ErrTaskInFlight", err)
	}

	code := 0
	if err := s.FinishTask(ctx, "t1", TaskSucceeded, &code, ""); err != nil {
		t.Fatalf("finish: %v", err)
	}
	if err := s.AcceptTask(ctx, Task{ID: "t2", Command: "make"}); err != nil {
		t.Fatalf("accept after terminal: %v", err)
	}
}

// The call-home restates the current assignment on every poll, so the agent is
// told about its own running task repeatedly. That must be a no-op, not an error.
func TestRedeliveryOfTheHeldTaskIsANoOp(t *testing.T) {
	s, ctx := open(t)

	if err := s.AcceptTask(ctx, Task{ID: "t1", Command: "make"}); err != nil {
		t.Fatalf("accept: %v", err)
	}
	if err := s.StartTask(ctx, "t1", "c1"); err != nil {
		t.Fatalf("start: %v", err)
	}
	if err := s.AcceptTask(ctx, Task{ID: "t1", Command: "make"}); err != nil {
		t.Fatalf("redelivery should be a no-op, got %v", err)
	}

	got, err := s.InFlightTask(ctx)
	if err != nil {
		t.Fatalf("in flight: %v", err)
	}
	if got == nil || got.State != TaskStarted || got.ContainerID != "c1" {
		t.Fatalf("redelivery clobbered state: %+v", got)
	}
}

// Dead-lettering is the third exit. Without it, an ambiguous outcome can never
// complete and never retry, so it would wedge the agent forever.
func TestDeadLetterReleasesTheSlot(t *testing.T) {
	s, ctx := open(t)

	if err := s.AcceptTask(ctx, Task{ID: "t1"}); err != nil {
		t.Fatalf("accept: %v", err)
	}
	if err := s.FinishTask(ctx, "t1", TaskDeadLettered, nil, "container vanished"); err != nil {
		t.Fatalf("dead letter: %v", err)
	}
	if err := s.AcceptTask(ctx, Task{ID: "t2"}); err != nil {
		t.Fatalf("slot not released by dead-letter: %v", err)
	}
}

func TestFinishRejectsNonTerminalState(t *testing.T) {
	s, ctx := open(t)
	if err := s.AcceptTask(ctx, Task{ID: "t1"}); err != nil {
		t.Fatalf("accept: %v", err)
	}
	if err := s.FinishTask(ctx, "t1", TaskStarted, nil, ""); err == nil {
		t.Fatal("finishing into a non-terminal state should be refused")
	}
}

// Finishing must queue the outcome in the same transaction that frees the slot.
// Otherwise an agent can free its slot having permanently forgotten a result it
// cannot recompute — the make targets are supplied, so re-running to rediscover
// an exit code is not an option.
func TestFinishQueuesResultAtomically(t *testing.T) {
	s, ctx := open(t)

	if err := s.AcceptTask(ctx, Task{ID: "t1"}); err != nil {
		t.Fatalf("accept: %v", err)
	}
	code := 2
	if err := s.FinishTask(ctx, "t1", TaskFailed, &code, "boom"); err != nil {
		t.Fatalf("finish: %v", err)
	}

	pending, err := s.PendingResults(ctx)
	if err != nil {
		t.Fatalf("pending results: %v", err)
	}
	if len(pending) != 1 {
		t.Fatalf("got %d pending results, want 1", len(pending))
	}
	if pending[0].TaskID != "t1" || pending[0].State != TaskFailed ||
		pending[0].ExitCode == nil || *pending[0].ExitCode != 2 {
		t.Fatalf("result not preserved: %+v", pending[0])
	}

	if err := s.AckResult(ctx, "t1"); err != nil {
		t.Fatalf("ack: %v", err)
	}
	pending, _ = s.PendingResults(ctx)
	if len(pending) != 0 {
		t.Fatalf("ack left %d results", len(pending))
	}
}

// Log order is the contract: the server appends under an autoincrement the
// readers page with `id > after`, so a chunk arriving out of order leaves a hole
// no reader revisits.
func TestPendingLogsAreOrderedAndResumeFromWatermark(t *testing.T) {
	s, ctx := open(t)

	var seqs []int64
	for _, msg := range []string{"one", "two", "three", "four"} {
		seq, err := s.AppendLog(ctx, "t1", "stdout", []byte(msg))
		if err != nil {
			t.Fatalf("append %s: %v", msg, err)
		}
		seqs = append(seqs, seq)
	}
	for i := 1; i < len(seqs); i++ {
		if seqs[i] <= seqs[i-1] {
			t.Fatalf("seq not monotonic: %v", seqs)
		}
	}

	batch, err := s.PendingLogs(ctx, 2)
	if err != nil {
		t.Fatalf("pending: %v", err)
	}
	if len(batch) != 2 || string(batch[0].Chunk) != "one" || string(batch[1].Chunk) != "two" {
		t.Fatalf("first batch wrong: %+v", batch)
	}

	if err := s.AckLogs(ctx, batch[1].Seq); err != nil {
		t.Fatalf("ack: %v", err)
	}

	batch, err = s.PendingLogs(ctx, 10)
	if err != nil {
		t.Fatalf("pending after ack: %v", err)
	}
	if len(batch) != 2 || string(batch[0].Chunk) != "three" {
		t.Fatalf("resume wrong: %+v", batch)
	}
}

// The watermark is durable, so an agent that dies mid-drain resumes rather than
// replaying from zero. Reopening the same file is the closest honest test.
func TestLogWatermarkSurvivesReopen(t *testing.T) {
	ctx := context.Background()
	path := filepath.Join(t.TempDir(), "agent.db")

	s, err := Open(ctx, path)
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	for _, m := range []string{"a", "b", "c"} {
		if _, err := s.AppendLog(ctx, "t1", "stdout", []byte(m)); err != nil {
			t.Fatalf("append: %v", err)
		}
	}
	batch, _ := s.PendingLogs(ctx, 2)
	if err := s.AckLogs(ctx, batch[1].Seq); err != nil {
		t.Fatalf("ack: %v", err)
	}
	_ = s.Close()

	reopened, err := Open(ctx, path)
	if err != nil {
		t.Fatalf("reopen: %v", err)
	}
	defer func() { _ = reopened.Close() }()

	batch, err = reopened.PendingLogs(ctx, 10)
	if err != nil {
		t.Fatalf("pending after reopen: %v", err)
	}
	if len(batch) != 1 || string(batch[0].Chunk) != "c" {
		t.Fatalf("watermark lost across restart: %+v", batch)
	}
}

// A late or duplicated ack must not rewind the watermark and re-serve chunks
// the server already has.
func TestAckDoesNotRewindWatermark(t *testing.T) {
	s, ctx := open(t)
	for _, m := range []string{"a", "b", "c"} {
		if _, err := s.AppendLog(ctx, "t1", "stdout", []byte(m)); err != nil {
			t.Fatalf("append: %v", err)
		}
	}
	all, _ := s.PendingLogs(ctx, 10)
	if err := s.AckLogs(ctx, all[2].Seq); err != nil {
		t.Fatalf("ack all: %v", err)
	}
	if err := s.AckLogs(ctx, all[0].Seq); err != nil {
		t.Fatalf("stale ack: %v", err)
	}
	if got, _ := s.PendingLogs(ctx, 10); len(got) != 0 {
		t.Fatalf("stale ack rewound the watermark: %+v", got)
	}
}
