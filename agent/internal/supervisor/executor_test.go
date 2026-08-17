package supervisor

import (
	"context"
	"path/filepath"
	"strings"
	"testing"

	"github.com/glueops/autoglue/agent/internal/runner"
	"github.com/glueops/autoglue/agent/internal/store"
)

type fakeRunner struct {
	found    runner.Container
	findErr  error
	started  []runner.Spec
	startID  string
	followed []string
	since    []string
	exitCode int
	lines    [][3]string
	removed  []string
}

func (f *fakeRunner) Find(context.Context, string) (runner.Container, error) {
	return f.found, f.findErr
}
func (f *fakeRunner) Start(_ context.Context, spec runner.Spec) (string, error) {
	f.started = append(f.started, spec)
	return f.startID, nil
}
func (f *fakeRunner) Follow(_ context.Context, id, since string, onLog runner.LogFunc) (int, error) {
	f.followed = append(f.followed, id)
	f.since = append(f.since, since)
	for _, l := range f.lines {
		if err := onLog(l[0], l[1], l[2]); err != nil {
			return 0, err
		}
	}
	return f.exitCode, nil
}
func (f *fakeRunner) Remove(_ context.Context, id string) error {
	f.removed = append(f.removed, id)
	return nil
}

func newExec(t *testing.T, r runner.Runner) (*Supervisor, *store.Store, context.Context) {
	t.Helper()
	ctx := context.Background()
	st, err := store.Open(ctx, filepath.Join(t.TempDir(), "agent.db"))
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	t.Cleanup(func() { _ = st.Close() })
	return New(st, &fakeClient{}, nil, Options{}).WithRunner(r), st, ctx
}

func seedTask(t *testing.T, st *store.Store, ctx context.Context, id, args string) {
	t.Helper()
	if err := st.AcceptTask(ctx, store.Task{ID: id, Command: "run_make", Args: args}); err != nil {
		t.Fatalf("seed: %v", err)
	}
}

const goodArgs = `{"cluster_id":"cl-1","run_id":"run-1","image":"img","tag":"v1",
	"target":"ping-servers","work_dir":"/opt/c","mounts":["/home/d/.ssh:/root/.ssh"]}`

// Absent container + assigned task: nothing has run yet, so start it.
func TestExecuteStartsWhenNothingIsRunning(t *testing.T) {
	r := &fakeRunner{found: runner.Container{Phase: runner.PhaseAbsent}, startID: "c1"}
	s, st, ctx := newExec(t, r)
	seedTask(t, st, ctx, "t1", goodArgs)

	if err := s.execute(ctx); err != nil {
		t.Fatalf("execute: %v", err)
	}
	if len(r.started) != 1 || r.started[0].Target != "ping-servers" {
		t.Fatalf("start specs = %+v", r.started)
	}
	if len(r.started[0].Mounts) != 1 || r.started[0].Mounts[0].Target != "/root/.ssh" {
		t.Fatalf("mounts not parsed: %+v", r.started[0].Mounts)
	}
	if len(r.followed) != 1 || r.followed[0] != "c1" {
		t.Fatalf("did not follow the started container: %v", r.followed)
	}
}

// The case the whole design turns on: the agent restarted, the work did not.
// Adopt it — never kill and restart, because the make targets are supplied and
// re-running one part-way is not safe to assume.
func TestExecuteAdoptsARunningContainerInsteadOfRestarting(t *testing.T) {
	r := &fakeRunner{found: runner.Container{ID: "survivor", Phase: runner.PhaseRunning}, exitCode: 0}
	s, st, ctx := newExec(t, r)
	seedTask(t, st, ctx, "t1", goodArgs)
	if err := st.StartTask(ctx, "t1", "survivor"); err != nil {
		t.Fatalf("mark started: %v", err)
	}
	if err := st.SetLogsCursor(ctx, "t1", "2026-08-13T10:00:00Z"); err != nil {
		t.Fatalf("cursor: %v", err)
	}

	if err := s.execute(ctx); err != nil {
		t.Fatalf("execute: %v", err)
	}
	if len(r.started) != 0 {
		t.Fatalf("adoption started a second container: %+v", r.started)
	}
	if len(r.followed) != 1 || r.followed[0] != "survivor" {
		t.Fatalf("did not adopt: %v", r.followed)
	}
	if r.since[0] != "2026-08-13T10:00:00Z" {
		t.Fatalf("adoption replayed from the start instead of resuming: %q", r.since[0])
	}
	if len(r.removed) != 0 {
		t.Fatalf("adoption removed a container it should have followed: %v", r.removed)
	}
}

// The outcome happened while nobody was listening. Harvest it rather than
// treating it as unknowable — this is the most valuable thing adoption recovers.
func TestExecuteHarvestsAnExitedContainer(t *testing.T) {
	r := &fakeRunner{found: runner.Container{ID: "done", Phase: runner.PhaseExited, ExitCode: 3}}
	s, st, ctx := newExec(t, r)
	seedTask(t, st, ctx, "t1", goodArgs)
	if err := st.StartTask(ctx, "t1", "done"); err != nil {
		t.Fatalf("mark started: %v", err)
	}

	if err := s.execute(ctx); err != nil {
		t.Fatalf("execute: %v", err)
	}

	pending, _ := st.PendingResults(ctx)
	if len(pending) != 1 || pending[0].State != store.TaskFailed {
		t.Fatalf("outcome not harvested: %+v", pending)
	}
	if pending[0].ExitCode == nil || *pending[0].ExitCode != 3 {
		t.Fatalf("exit code lost: %+v", pending[0])
	}
	// Removing is correct only here, and it is what stops containers piling up
	// in the absence of --rm.
	if len(r.removed) != 1 || r.removed[0] != "done" {
		t.Fatalf("exited container not removed: %v", r.removed)
	}
}

// Started, container gone. Completed-but-unreported and half-applied are
// indistinguishable, so this must dead-letter — never retry, and never guess a
// success.
func TestExecuteDeadLettersAVanishedContainer(t *testing.T) {
	r := &fakeRunner{found: runner.Container{Phase: runner.PhaseAbsent}}
	s, st, ctx := newExec(t, r)
	seedTask(t, st, ctx, "t1", goodArgs)
	if err := st.StartTask(ctx, "t1", "gone"); err != nil {
		t.Fatalf("mark started: %v", err)
	}

	if err := s.execute(ctx); err != nil {
		t.Fatalf("execute: %v", err)
	}
	if len(r.started) != 0 {
		t.Fatalf("a vanished container was retried: %+v", r.started)
	}

	pending, _ := st.PendingResults(ctx)
	if len(pending) != 1 || pending[0].State != store.TaskDeadLettered {
		t.Fatalf("not dead-lettered: %+v", pending)
	}
	// And the slot is free, which is the entire reason dead-lettering exists.
	if err := st.AcceptTask(ctx, store.Task{ID: "t2"}); err != nil {
		t.Fatalf("dead-letter did not release the slot: %v", err)
	}
}

// Version skew: args from a control plane ahead of this binary. Dead-letter
// with the reason rather than crash-looping or silently skipping.
func TestExecuteDeadLettersUnparsableArgs(t *testing.T) {
	r := &fakeRunner{found: runner.Container{Phase: runner.PhaseAbsent}, startID: "c1"}
	s, st, ctx := newExec(t, r)
	seedTask(t, st, ctx, "t1", `{"cluster_id": [this is not json`)

	if err := s.execute(ctx); err != nil {
		t.Fatalf("execute: %v", err)
	}
	if len(r.started) != 0 {
		t.Fatalf("started a container from unparsable args: %+v", r.started)
	}
	pending, _ := st.PendingResults(ctx)
	if len(pending) != 1 || pending[0].State != store.TaskDeadLettered {
		t.Fatalf("not dead-lettered: %+v", pending)
	}
	if !strings.Contains(pending[0].Error, "unparsable") {
		t.Fatalf("reason not recorded: %q", pending[0].Error)
	}
}

// Output must reach the outbox, and the cursor must advance with it so an
// adopted follow resumes instead of replaying.
func TestExecuteIngestsLogsAndAdvancesCursor(t *testing.T) {
	r := &fakeRunner{
		found:   runner.Container{Phase: runner.PhaseAbsent},
		startID: "c1",
		lines: [][3]string{
			{"stdout", "2026-08-13T10:00:00Z", "first"},
			{"stdout", "2026-08-13T10:00:05Z", "second"},
		},
	}
	s, st, ctx := newExec(t, r)
	seedTask(t, st, ctx, "t1", goodArgs)

	if err := s.execute(ctx); err != nil {
		t.Fatalf("execute: %v", err)
	}

	chunks, _ := st.PendingLogs(ctx, 10)
	if len(chunks) != 2 || !strings.Contains(string(chunks[0].Chunk), "first") {
		t.Fatalf("logs not ingested: %+v", chunks)
	}
	if chunks[0].Seq >= chunks[1].Seq {
		t.Fatalf("seq not monotonic across ingest: %+v", chunks)
	}
}

// The store's index guards one task at a time; this guards one executor at a
// time, so two docker follows never race to write the same terminal state.
func TestSpawnExecutorRunsOneAtATime(t *testing.T) {
	r := &fakeRunner{found: runner.Container{Phase: runner.PhaseAbsent}, startID: "c1"}
	s, _, _ := newExec(t, r)

	if !s.executing.CompareAndSwap(false, true) {
		t.Fatal("fresh supervisor already marked executing")
	}
	// With the flag held, a spawn must be a no-op rather than a second goroutine.
	s.spawnExecutor(context.Background())
	if len(r.started) != 0 {
		t.Fatalf("a second executor ran while one was in flight: %+v", r.started)
	}
}
