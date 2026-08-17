// Package supervisor runs the agent's two planes.
//
// They are separate loops on separate timers, and that is the whole point. The
// config plane converges toward a desired generation; the task plane runs one
// imperative task at a time. Coupling them would mean a six-hour make target
// blocks credential rotation until it finishes — and an agent whose credential
// expires mid-task goes dark while its work is running perfectly.
//
// Decoupling also keeps liveness honest: the config loop keeps calling home
// while a task runs, so long work never looks like a dead agent.
package supervisor

import (
	"context"
	"errors"
	"log/slog"
	"sync/atomic"
	"time"

	"github.com/glueops/autoglue/agent/internal/api"
	"github.com/glueops/autoglue/agent/internal/runner"
	"github.com/glueops/autoglue/agent/internal/store"
)

// Reconciler applies one resource type. Registering none for a type is an
// expected condition, not a crash: agents are deployed software and will run
// behind the control plane, so an unknown type is what version skew looks like.
type Reconciler interface {
	Apply(ctx context.Context, r store.DesiredResource) error
}

type Options struct {
	ConfigInterval time.Duration
	TaskInterval   time.Duration
	OutboxInterval time.Duration
	Log            *slog.Logger
}

type Supervisor struct {
	store       *store.Store
	client      api.Client
	runner      runner.Runner
	reconcilers map[string]Reconciler
	opts        Options
	log         *slog.Logger

	// executing guards the long-running executor goroutine. Following a
	// container can take hours, and it must not run on the ticker's goroutine
	// or block the config loop — the config loop is the liveness heartbeat, and
	// an agent that stops calling home while a six-hour target runs looks dead
	// to the recovery sweeper.
	executing atomic.Bool
}

// WithRunner attaches the container runner. Separate from New because the
// scaffold is useful without one: the store and the config plane are complete
// and testable before any container is ever started.
func (s *Supervisor) WithRunner(r runner.Runner) *Supervisor {
	s.runner = r
	return s
}

func New(st *store.Store, client api.Client, reconcilers map[string]Reconciler, opts Options) *Supervisor {
	if opts.ConfigInterval <= 0 {
		opts.ConfigInterval = 30 * time.Second
	}
	if opts.TaskInterval <= 0 {
		opts.TaskInterval = 10 * time.Second
	}
	if opts.OutboxInterval <= 0 {
		opts.OutboxInterval = 5 * time.Second
	}
	if opts.Log == nil {
		opts.Log = slog.Default()
	}
	return &Supervisor{
		store:       st,
		client:      client,
		reconcilers: reconcilers,
		opts:        opts,
		log:         opts.Log,
	}
}

// Run drives all three loops until ctx is cancelled.
//
// A tick that errors is logged and the loop continues. The agent's job is to
// keep converging and keep draining; a control plane that is briefly
// unreachable is an ordinary condition, not a reason to exit and lose the
// outbox.
func (s *Supervisor) Run(ctx context.Context) error {
	done := make(chan struct{}, 3)

	go s.loop(ctx, "config", s.opts.ConfigInterval, s.ConfigTick, done)
	go s.loop(ctx, "task", s.opts.TaskInterval, s.TaskTick, done)
	go s.loop(ctx, "outbox", s.opts.OutboxInterval, s.OutboxTick, done)

	<-ctx.Done()
	for range 3 {
		<-done
	}
	return ctx.Err()
}

func (s *Supervisor) loop(ctx context.Context, name string, every time.Duration, tick func(context.Context) error, done chan<- struct{}) {
	defer func() { done <- struct{}{} }()

	t := time.NewTicker(every)
	defer t.Stop()

	for {
		if err := tick(ctx); err != nil && !errors.Is(err, context.Canceled) {
			s.log.Warn("tick failed", "loop", name, "error", err)
		}
		select {
		case <-ctx.Done():
			return
		case <-t.C:
		}
	}
}

// ConfigTick fetches any newer generation, reconciles toward it, and reports
// convergence. It also serves as the liveness heartbeat, which is why it runs
// on its own timer rather than between tasks.
func (s *Supervisor) ConfigTick(ctx context.Context) error {
	current, err := s.store.CurrentGeneration(ctx)
	if err != nil {
		return err
	}

	snap, err := s.client.Sync(ctx, current)
	if err != nil {
		return err
	}
	if snap != nil {
		if err := s.store.ReplaceDesiredSnapshot(ctx, snap.Generation, toDesired(snap.Resources)); err != nil {
			return err
		}
		current = snap.Generation
	}
	if current == 0 {
		return nil // nothing has ever been delivered
	}

	if err := s.reconcile(ctx, current); err != nil {
		return err
	}

	// AppliedGeneration advances only on real convergence. A pass that finished
	// is not a pass that succeeded — one that failed every resource also
	// finishes, and reporting that as applied would tell the control plane the
	// fleet is healthy while nothing landed.
	converged, err := s.store.GenerationConverged(ctx, current)
	if err != nil {
		return err
	}
	if converged {
		if err := s.store.SetAppliedGeneration(ctx, current); err != nil {
			return err
		}
		if err := s.store.PruneGenerationsBefore(ctx, current); err != nil {
			return err
		}
	}

	applied, err := s.store.AppliedGeneration(ctx)
	if err != nil {
		return err
	}
	statuses, err := s.store.ListReconcileStatus(ctx)
	if err != nil {
		return err
	}
	return s.client.ReportReconcile(ctx, api.ReconcileReport{
		CurrentGeneration: current,
		AppliedGeneration: applied,
		Healthy:           converged,
		Resources:         toStatuses(statuses),
	})
}

func (s *Supervisor) reconcile(ctx context.Context, generation int64) error {
	resources, err := s.store.ListDesiredResources(ctx, generation)
	if err != nil {
		return err
	}

	for _, r := range resources {
		if err := ctx.Err(); err != nil {
			return err
		}

		st := store.ReconcileStatus{
			ResourceType:      r.ResourceType,
			ResourceID:        r.ResourceID,
			DesiredGeneration: generation,
			State:             store.ReconcileApplied,
		}

		rec, ok := s.reconcilers[r.ResourceType]
		switch {
		case !ok:
			// Recorded as failed, never skipped. A silent skip would let the
			// agent report convergence to a generation it did not apply, and
			// the control plane would believe it.
			st.State = store.ReconcileFailed
			st.LastError = "no reconciler registered for " + r.ResourceType
		default:
			if err := rec.Apply(ctx, r); err != nil {
				st.State = store.ReconcileFailed
				st.LastError = err.Error()
			}
		}

		if err := s.store.UpsertReconcileStatus(ctx, st); err != nil {
			return err
		}
	}
	return nil
}

// TaskTick asks what the agent should be working on and accepts it if the slot
// is free. Execution itself is not wired up yet: the runner that starts a
// container and adopts a surviving one is the next piece.
func (s *Supervisor) TaskTick(ctx context.Context) error {
	assignment, err := s.client.Assignment(ctx)
	if err != nil {
		return err
	}
	if assignment == nil {
		return nil
	}

	held, err := s.store.InFlightTask(ctx)
	if err != nil {
		return err
	}
	if held != nil && held.ID != assignment.TaskID {
		// The control plane and the agent disagree about reality. Never guess:
		// resolving it by dropping one of them is how two make targets end up
		// running against one cluster.
		s.log.Warn("assignment conflicts with the task in flight",
			"holding", held.ID, "assigned", assignment.TaskID)
		return nil
	}

	// Redelivery of the held task is the normal case, and AcceptTask treats it
	// as a no-op.
	if err := s.store.AcceptTask(ctx, store.Task{
		ID:      assignment.TaskID,
		Command: assignment.Command,
		Args:    assignment.Args,
	}); err != nil {
		return err
	}

	s.spawnExecutor(ctx)
	return nil
}

// spawnExecutor runs execute in the background, at most one at a time.
//
// The compare-and-swap is the concurrency-1 rule at the process level, sitting
// above the store's partial unique index. Two executors on one task would mean
// two docker follows racing to write the same terminal state.
func (s *Supervisor) spawnExecutor(ctx context.Context) {
	if s.runner == nil {
		return
	}
	if !s.executing.CompareAndSwap(false, true) {
		return
	}
	go func() {
		defer s.executing.Store(false)
		if err := s.execute(ctx); err != nil && !errors.Is(err, context.Canceled) {
			s.log.Warn("task execution failed", "error", err)
		}
	}()
}

// AdoptOnStart reconciles whatever was in flight when the process last died.
// Called once before the loops begin, so a surviving container is picked up
// immediately rather than at the first task tick — and, more importantly, so a
// vanished one is dead-lettered before anything can be assigned on top of it.
func (s *Supervisor) AdoptOnStart(ctx context.Context) error {
	held, err := s.store.InFlightTask(ctx)
	if err != nil || held == nil {
		return err
	}
	s.log.Info("task was in flight at last shutdown",
		"task_id", held.ID, "state", held.State, "container_id", held.ContainerID)
	return s.execute(ctx)
}

// OutboxTick drains what the agent knows and the control plane does not.
//
// Results before logs: an outcome is the thing that cannot be recomputed, and
// the make targets are supplied rather than authored here, so re-running to
// rediscover an exit code is not available. Logs are worth having and worth
// losing last.
func (s *Supervisor) OutboxTick(ctx context.Context) error {
	results, err := s.store.PendingResults(ctx)
	if err != nil {
		return err
	}
	for _, r := range results {
		if err := s.client.FinishTask(ctx, api.TaskResult{
			TaskID:   r.TaskID,
			State:    r.State,
			ExitCode: r.ExitCode,
			Error:    r.Error,
		}); err != nil {
			return err
		}
		if err := s.store.AckResult(ctx, r.TaskID); err != nil {
			return err
		}
	}

	const logBatch = 64
	for {
		chunks, err := s.store.PendingLogs(ctx, logBatch)
		if err != nil {
			return err
		}
		if len(chunks) == 0 {
			return nil
		}
		if err := s.client.AppendLogs(ctx, toAPIChunks(chunks)); err != nil {
			return err
		}
		// Ack only what was accepted, and only after acceptance. The watermark
		// is what a restarted agent resumes from, so advancing it early loses
		// output permanently.
		if err := s.store.AckLogs(ctx, chunks[len(chunks)-1].Seq); err != nil {
			return err
		}
		if len(chunks) < logBatch {
			return nil
		}
	}
}

func toDesired(in []api.Resource) []store.DesiredResource {
	out := make([]store.DesiredResource, 0, len(in))
	for _, r := range in {
		out = append(out, store.DesiredResource{
			ResourceType: r.ResourceType,
			ResourceID:   r.ResourceID,
			ResourceName: r.ResourceName,
			Phase:        r.Phase,
			DependsOn:    r.DependsOn,
			Spec:         r.Spec,
			SpecHash:     r.SpecHash,
			Required:     r.Required,
		})
	}
	return out
}

func toStatuses(in []store.ReconcileStatus) []api.ResourceStatus {
	out := make([]api.ResourceStatus, 0, len(in))
	for _, st := range in {
		out = append(out, api.ResourceStatus{
			ResourceType:      st.ResourceType,
			ResourceID:        st.ResourceID,
			DesiredGeneration: st.DesiredGeneration,
			State:             st.State,
			LastError:         st.LastError,
		})
	}
	return out
}

func toAPIChunks(in []store.LogChunk) []api.LogChunk {
	out := make([]api.LogChunk, 0, len(in))
	for _, c := range in {
		out = append(out, api.LogChunk{
			TaskID: c.TaskID,
			Seq:    c.Seq,
			Stream: c.Stream,
			Chunk:  string(c.Chunk),
		})
	}
	return out
}
