// Package runner starts and adopts the containers that carry tasks.
//
// The central decision here is that containers are started **detached**, and
// the agent then attaches for logs. Start and adopt therefore converge on one
// implementation: both are "follow the logs of a container that is already
// running, then wait for its exit code". A container the agent did not start is
// no different from one it did.
//
// It also removes the failure that motivated all of this. With a detached
// container nothing is proxying signals into it and no attached CLI can be
// killed out from under it, so the agent dying, being upgraded or losing its
// network never touches work in flight.
package runner

import (
	"context"
	"errors"
)

// Phase is a container's lifecycle state as the agent finds it.
type Phase string

const (
	// PhaseRunning: adopt it. Do not start anything, do not dispatch.
	PhaseRunning Phase = "running"
	// PhaseExited: harvest the exit code and logs, report, then remove.
	PhaseExited Phase = "exited"
	// PhaseAbsent: nothing ran, or the host was rebuilt.
	PhaseAbsent Phase = "absent"
)

// Container is what a lookup by task id found.
type Container struct {
	ID       string
	Phase    Phase
	ExitCode int
}

// ErrNoContainer is returned when a task has no container at all. For a task
// the agent believes it started, this is the unknowable case: completed but
// unreported and half-applied are indistinguishable, so it dead-letters rather
// than retrying.
var ErrNoContainer = errors.New("no container found for task")

// Spec is everything needed to launch one task's container.
type Spec struct {
	TaskID    string
	ClusterID string
	RunID     string
	Image     string
	Tag       string
	Target    string
	WorkDir   string
	Mounts    []Mount
}

type Mount struct {
	Source string
	Target string
}

// LogFunc receives one line of container output. Returning an error stops the
// stream: the caller uses that to stop pumping into a full outbox rather than
// growing it without bound.
type LogFunc func(stream, timestamp, line string) error

// Runner is the container lifecycle the task plane drives.
type Runner interface {
	// Find locates a task's container by label, whatever its phase.
	Find(ctx context.Context, taskID string) (Container, error)

	// Start launches a detached, labelled container and returns its id.
	Start(ctx context.Context, spec Spec) (string, error)

	// Follow streams logs from since (a docker RFC3339Nano timestamp, empty for
	// the beginning) until the container exits, then returns its exit code.
	// Safe to call on a container this process did not start.
	Follow(ctx context.Context, containerID, since string, onLog LogFunc) (int, error)

	// Remove deletes a container once its outcome has been harvested. Never
	// call it on a running one: the make targets are supplied rather than
	// authored here, so killing one part-way can leave state nothing can
	// reconstruct.
	Remove(ctx context.Context, containerID string) error
}

// Label keys. These are the contract that makes a surviving container findable
// by anything — the agent after a restart, or the control plane's SSH recovery
// path reaching in from outside.
const (
	LabelCluster = "autoglue.cluster"
	LabelRun     = "autoglue.run"
	LabelTask    = "autoglue.task"
)
