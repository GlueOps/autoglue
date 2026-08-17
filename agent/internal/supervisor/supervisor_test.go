package supervisor

import (
	"context"
	"errors"
	"path/filepath"
	"testing"

	"github.com/glueops/autoglue/agent/internal/api"
	"github.com/glueops/autoglue/agent/internal/store"
)

type fakeClient struct {
	snapshot   *api.Snapshot
	assignment *api.Assignment

	reports  []api.ReconcileReport
	logs     [][]api.LogChunk
	results  []api.TaskResult
	logErr   error
	finishOK bool
}

func (f *fakeClient) Sync(context.Context, int64) (*api.Snapshot, error) {
	s := f.snapshot
	f.snapshot = nil // deliver a generation once, like a real conditional fetch
	return s, nil
}
func (f *fakeClient) Assignment(context.Context) (*api.Assignment, error) {
	return f.assignment, nil
}
func (f *fakeClient) ReportReconcile(_ context.Context, r api.ReconcileReport) error {
	f.reports = append(f.reports, r)
	return nil
}
func (f *fakeClient) StartTask(context.Context, string, string) error { return nil }
func (f *fakeClient) AppendLogs(_ context.Context, c []api.LogChunk) error {
	if f.logErr != nil {
		return f.logErr
	}
	f.logs = append(f.logs, c)
	return nil
}
func (f *fakeClient) FinishTask(_ context.Context, r api.TaskResult) error {
	f.results = append(f.results, r)
	f.finishOK = true
	return nil
}

type okReconciler struct{ applied int }

func (r *okReconciler) Apply(context.Context, store.DesiredResource) error {
	r.applied++
	return nil
}

func newSup(t *testing.T, c api.Client, recs map[string]Reconciler) (*Supervisor, *store.Store, context.Context) {
	t.Helper()
	ctx := context.Background()
	st, err := store.Open(ctx, filepath.Join(t.TempDir(), "agent.db"))
	if err != nil {
		t.Fatalf("open store: %v", err)
	}
	t.Cleanup(func() { _ = st.Close() })
	return New(st, c, recs, Options{}), st, ctx
}

func snap(gen int64, types ...string) *api.Snapshot {
	s := &api.Snapshot{Generation: gen, FullSnapshot: true}
	for i, ty := range types {
		s.Resources = append(s.Resources, api.Resource{
			ResourceType: ty,
			ResourceID:   string(rune('a' + i)),
			Spec:         "{}",
			SpecHash:     "h",
			Required:     true,
		})
	}
	return s
}

func TestConfigTickAppliesAndReportsConvergence(t *testing.T) {
	rec := &okReconciler{}
	c := &fakeClient{snapshot: snap(4, "cluster_shape")}
	s, st, ctx := newSup(t, c, map[string]Reconciler{"cluster_shape": rec})

	if err := s.ConfigTick(ctx); err != nil {
		t.Fatalf("tick: %v", err)
	}
	if rec.applied != 1 {
		t.Fatalf("reconciler applied %d times, want 1", rec.applied)
	}

	applied, _ := st.AppliedGeneration(ctx)
	if applied != 4 {
		t.Fatalf("applied generation = %d, want 4", applied)
	}
	if len(c.reports) != 1 || !c.reports[0].Healthy || c.reports[0].AppliedGeneration != 4 {
		t.Fatalf("report wrong: %+v", c.reports)
	}
}

// Version skew: an agent running behind the control plane meets a resource type
// it has no reconciler for. It must report that generation as NOT applied,
// rather than skipping the resource and claiming convergence it never reached.
func TestUnknownResourceTypeBlocksConvergence(t *testing.T) {
	c := &fakeClient{snapshot: snap(9, "cluster_shape", "something_from_the_future")}
	s, st, ctx := newSup(t, c, map[string]Reconciler{"cluster_shape": &okReconciler{}})

	if err := s.ConfigTick(ctx); err != nil {
		t.Fatalf("tick: %v", err)
	}

	if applied, _ := st.AppliedGeneration(ctx); applied != 0 {
		t.Fatalf("applied generation advanced to %d despite an unhandled resource", applied)
	}
	if c.reports[0].Healthy {
		t.Fatal("reported healthy with an unhandled resource type")
	}

	var found bool
	for _, rs := range c.reports[0].Resources {
		if rs.ResourceType == "something_from_the_future" {
			found = true
			if rs.State != store.ReconcileFailed || rs.LastError == "" {
				t.Fatalf("unknown type not reported as a failure: %+v", rs)
			}
		}
	}
	if !found {
		t.Fatal("unknown resource type was silently dropped from the report")
	}
}

// An assignment naming a different task than the one in flight is a
// disagreement about reality. The agent must hold, not switch — swapping tasks
// is how two make targets end up running against one cluster.
func TestConflictingAssignmentIsRefused(t *testing.T) {
	c := &fakeClient{assignment: &api.Assignment{TaskID: "t2", Command: "make"}}
	s, st, ctx := newSup(t, c, nil)

	if err := st.AcceptTask(ctx, store.Task{ID: "t1", Command: "make"}); err != nil {
		t.Fatalf("seed: %v", err)
	}
	if err := s.TaskTick(ctx); err != nil {
		t.Fatalf("tick returned an error rather than holding: %v", err)
	}

	held, _ := st.InFlightTask(ctx)
	if held == nil || held.ID != "t1" {
		t.Fatalf("in-flight task changed to %+v", held)
	}
}

func TestRedeliveredAssignmentIsAccepted(t *testing.T) {
	c := &fakeClient{assignment: &api.Assignment{TaskID: "t1", Command: "make"}}
	s, st, ctx := newSup(t, c, nil)

	if err := s.TaskTick(ctx); err != nil {
		t.Fatalf("first tick: %v", err)
	}
	if err := s.TaskTick(ctx); err != nil {
		t.Fatalf("redelivery should be a no-op: %v", err)
	}
	held, _ := st.InFlightTask(ctx)
	if held == nil || held.ID != "t1" {
		t.Fatalf("held = %+v, want t1", held)
	}
}

// Results drain before logs, and the result is the thing that must not be lost:
// an exit code cannot be recomputed by re-running a non-idempotent target.
func TestOutboxDrainsResultsBeforeLogs(t *testing.T) {
	c := &fakeClient{logErr: errors.New("api unreachable")}
	s, st, ctx := newSup(t, c, nil)

	if err := st.AcceptTask(ctx, store.Task{ID: "t1"}); err != nil {
		t.Fatalf("accept: %v", err)
	}
	if _, err := st.AppendLog(ctx, "t1", "stdout", []byte("output")); err != nil {
		t.Fatalf("append: %v", err)
	}
	code := 0
	if err := st.FinishTask(ctx, "t1", store.TaskSucceeded, &code, ""); err != nil {
		t.Fatalf("finish: %v", err)
	}

	// Log delivery fails, so the tick errors — but the result must already be
	// gone, because it was drained first.
	if err := s.OutboxTick(ctx); err == nil {
		t.Fatal("expected the log failure to surface")
	}
	if !c.finishOK {
		t.Fatal("result was not delivered before logs were attempted")
	}
	if pending, _ := st.PendingResults(ctx); len(pending) != 0 {
		t.Fatalf("result still pending after successful delivery: %+v", pending)
	}

	// The log survives for the next drain rather than being acked on failure.
	if chunks, _ := st.PendingLogs(ctx, 10); len(chunks) != 1 {
		t.Fatalf("failed log delivery lost the chunk: %+v", chunks)
	}
}
