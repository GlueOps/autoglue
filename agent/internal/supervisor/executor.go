package supervisor

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/glueops/autoglue/agent/internal/runner"
	"github.com/glueops/autoglue/agent/internal/store"
)

// TaskArgs is what the control plane sends alongside a task. It describes one
// container, because a task is one container — see the note on
// tasks.container_id in the store schema.
type TaskArgs struct {
	ClusterID string   `json:"cluster_id"`
	RunID     string   `json:"run_id"`
	Image     string   `json:"image"`
	Tag       string   `json:"tag"`
	Target    string   `json:"target"`
	WorkDir   string   `json:"work_dir"`
	Mounts    []string `json:"mounts"` // "source:target"
}

// execute drives whatever the agent is currently holding to a terminal state.
//
// The branch on container phase is the whole point of adoption. "The agent is
// unreachable" and "the work is dead" are separate facts, and a container that
// is still applying must never be killed and restarted: the make targets are
// supplied rather than authored here, so re-running one part-way through is not
// something that can be assumed safe.
func (s *Supervisor) execute(ctx context.Context) error {
	held, err := s.store.InFlightTask(ctx)
	if err != nil || held == nil {
		return err
	}
	if s.runner == nil {
		return nil // scaffold: nothing to execute against yet
	}

	c, err := s.runner.Find(ctx, held.ID)
	if err != nil {
		return err
	}

	switch {
	case c.Phase == runner.PhaseAbsent && held.State == store.TaskAssigned:
		return s.startAndFollow(ctx, held)

	case c.Phase == runner.PhaseAbsent:
		// Started, and the container is gone. Completed-but-unreported and
		// half-applied are indistinguishable from here, and there is no
		// evidence anywhere that separates them. Dead-letter: it releases the
		// slot without claiming the work happened, and it must never be retried
		// automatically.
		s.log.Warn("task container vanished; dead-lettering",
			"task_id", held.ID, "container_id", held.ContainerID)
		return s.store.FinishTask(ctx, held.ID, store.TaskDeadLettered, nil,
			"container for a started task could not be found")

	case c.Phase == runner.PhaseRunning:
		// Adopt. Started by this process or a previous one, it makes no
		// difference: following a container is the same operation either way.
		if held.State == store.TaskAssigned {
			if err := s.store.StartTask(ctx, held.ID, c.ID); err != nil {
				return err
			}
		}
		s.log.Info("adopting running container", "task_id", held.ID, "container_id", c.ID)
		return s.follow(ctx, held.ID, c.ID, held.LogsCursor)

	default: // PhaseExited
		// The outcome was produced while nobody was listening. Harvest it,
		// then remove — this is the one place removing a container is correct,
		// and it is also what stops them accumulating in the absence of --rm.
		s.log.Info("harvesting exited container", "task_id", held.ID, "container_id", c.ID, "exit_code", c.ExitCode)
		if err := s.drainLogs(ctx, held.ID, c.ID, held.LogsCursor); err != nil {
			return err
		}
		if err := s.finishByExitCode(ctx, held.ID, c.ExitCode); err != nil {
			return err
		}
		return s.runner.Remove(ctx, c.ID)
	}
}

func (s *Supervisor) startAndFollow(ctx context.Context, held *store.Task) error {
	var args TaskArgs
	if err := json.Unmarshal([]byte(held.Args), &args); err != nil {
		// Args the agent cannot parse are the version-skew case: a control
		// plane ahead of this binary. Dead-letter with the reason rather than
		// crash-looping or, worse, silently skipping.
		return s.store.FinishTask(ctx, held.ID, store.TaskDeadLettered, nil,
			fmt.Sprintf("unparsable task args: %v", err))
	}

	spec := runner.Spec{
		TaskID:    held.ID,
		ClusterID: args.ClusterID,
		RunID:     args.RunID,
		Image:     args.Image,
		Tag:       args.Tag,
		Target:    args.Target,
		WorkDir:   args.WorkDir,
	}
	for _, m := range args.Mounts {
		src, dst, ok := splitMount(m)
		if !ok {
			return s.store.FinishTask(ctx, held.ID, store.TaskDeadLettered, nil,
				"malformed mount specification: "+m)
		}
		spec.Mounts = append(spec.Mounts, runner.Mount{Source: src, Target: dst})
	}

	containerID, err := s.runner.Start(ctx, spec)
	if err != nil {
		return err
	}
	// Recorded before following. If the agent dies in the next instant, the
	// container id is what lets its successor adopt rather than start a second
	// one against the same cluster.
	if err := s.store.StartTask(ctx, held.ID, containerID); err != nil {
		return err
	}
	if err := s.client.StartTask(ctx, held.ID, containerID); err != nil {
		s.log.Warn("could not report task start", "task_id", held.ID, "error", err)
	}
	return s.follow(ctx, held.ID, containerID, "")
}

func (s *Supervisor) follow(ctx context.Context, taskID, containerID, since string) error {
	code, err := s.runner.Follow(ctx, containerID, since, s.ingestLog(ctx, taskID))
	if err != nil {
		return err
	}
	return s.finishByExitCode(ctx, taskID, code)
}

// drainLogs picks up whatever a finished container emitted after the last
// ingest, without waiting on it.
func (s *Supervisor) drainLogs(ctx context.Context, taskID, containerID, since string) error {
	_, err := s.runner.Follow(ctx, containerID, since, s.ingestLog(ctx, taskID))
	return err
}

func (s *Supervisor) ingestLog(ctx context.Context, taskID string) runner.LogFunc {
	return func(stream, timestamp, line string) error {
		if _, err := s.store.AppendLog(ctx, taskID, stream, []byte(line+"\n")); err != nil {
			return err
		}
		// Advance the resume point as we go, so a crash costs the lines since
		// the last write rather than replaying the container's whole history.
		return s.store.SetLogsCursor(ctx, taskID, timestamp)
	}
}

func (s *Supervisor) finishByExitCode(ctx context.Context, taskID string, code int) error {
	state := store.TaskSucceeded
	errMsg := ""
	if code != 0 {
		state = store.TaskFailed
		errMsg = fmt.Sprintf("container exited %d", code)
	}
	return s.store.FinishTask(ctx, taskID, state, &code, errMsg)
}

func splitMount(m string) (src, dst string, ok bool) {
	for i := len(m) - 1; i > 0; i-- {
		if m[i] == ':' {
			return m[:i], m[i+1:], m[:i] != "" && m[i+1:] != ""
		}
	}
	return "", "", false
}
