// Package api is the contract between the agent and the autoglue control plane.
//
// It is an interface and a set of types, with no HTTP implementation yet: the
// server endpoints do not exist. Defining the contract first keeps the
// supervisor testable against a fake, and makes the read/write split explicit
// before any of it is wired up.
//
// The split is the design rule, not a style preference. A call is a GET when
// the mutation it causes lands on the bastion rather than on the API — the
// agent rewriting its own disk is no more a server-side change than flashing a
// firmware image you fetched. Anything that changes control-plane state is a
// POST.
package api

import "context"

// Client is what the supervisor talks to. Every method is safe to call again
// with the same arguments: the reads are idempotent by nature and the writes
// are keyed so a redelivery is recognised rather than duplicated.
type Client interface {
	// --- reads (GET) ---

	// Sync returns the desired-state snapshot when the control plane holds a
	// generation newer than currentGeneration, and nil when the agent is
	// already up to date.
	//
	// The generation travels as a query parameter, deliberately not as a
	// request body. A GET body has no defined semantics and is dropped by some
	// proxies, which would silently turn every sync into a full re-delivery.
	Sync(ctx context.Context, currentGeneration int64) (*Snapshot, error)

	// Assignment returns the task the agent should be working on, or nil.
	//
	// This states current assignment rather than popping a queue. That is what
	// makes it a safe method: an unacknowledged task is simply named again on
	// the next call, so a lost response needs no hand-out lease to recover, and
	// a control plane that thinks a task is running can be contradicted by an
	// agent that has nothing.
	Assignment(ctx context.Context) (*Assignment, error)

	// --- writes (POST) ---

	// ReportReconcile publishes convergence. The gap between current and
	// applied generation is what the control plane reads as "not yet
	// converged".
	ReportReconcile(ctx context.Context, r ReconcileReport) error

	// StartTask claims the assignment. Until this lands, Assignment keeps
	// naming the same task.
	StartTask(ctx context.Context, taskID, containerID string) error

	// AppendLogs delivers output chunks. Batches must be posted serially per
	// task and in seq order: the server stores them under an autoincrement that
	// readers page with `id > after`, so a chunk that arrives late lands behind
	// a cursor no reader will revisit. Seq is the idempotency key alongside the
	// task id.
	AppendLogs(ctx context.Context, chunks []LogChunk) error

	// FinishTask reports a terminal outcome. This is the only thing that
	// releases the agent's task slot.
	FinishTask(ctx context.Context, r TaskResult) error
}

// Snapshot is a whole desired-state generation, never a diff. An agent that has
// been switched off has no history to apply diffs against, and a snapshot is
// correct from any starting state including a wiped disk.
type Snapshot struct {
	Generation   int64      `json:"generation"`
	FullSnapshot bool       `json:"full_snapshot"`
	PollAfterSec int        `json:"poll_after_sec"`
	Resources    []Resource `json:"resources"`
}

type Resource struct {
	ResourceType string `json:"resource_type"`
	ResourceID   string `json:"resource_id"`
	ResourceName string `json:"resource_name"`
	Phase        int    `json:"phase"`
	DependsOn    string `json:"depends_on"`
	Spec         string `json:"spec"`
	SpecHash     string `json:"spec_hash"`
	Required     bool   `json:"required"`
}

type ReconcileReport struct {
	CurrentGeneration int64            `json:"current_generation"`
	AppliedGeneration int64            `json:"applied_generation"`
	Healthy           bool             `json:"healthy"`
	Resources         []ResourceStatus `json:"resources,omitempty"`
}

type ResourceStatus struct {
	ResourceType      string `json:"resource_type"`
	ResourceID        string `json:"resource_id"`
	DesiredGeneration int64  `json:"desired_generation"`
	State             string `json:"state"`
	LastError         string `json:"last_error,omitempty"`
}

// Assignment names the task the agent should be running. State lets the agent
// tell "you have not started this yet" from "you told me you started it", which
// is what surfaces drift on the next poll instead of at a lease expiry.
type Assignment struct {
	TaskID       string `json:"task_id"`
	Command      string `json:"command"`
	Args         string `json:"args"`
	State        string `json:"state"`
	PollAfterSec int    `json:"poll_after_sec"`
}

type LogChunk struct {
	TaskID string `json:"task_id"`
	Seq    int64  `json:"seq"`
	Stream string `json:"stream"`
	Chunk  string `json:"chunk"`
}

// TaskResult is terminal. DeadLettered is a real outcome rather than a kind of
// failure: it means nobody can say whether the work happened, which is not the
// same as knowing it failed, and it must never be retried automatically.
type TaskResult struct {
	TaskID   string `json:"task_id"`
	State    string `json:"state"`
	ExitCode *int   `json:"exit_code,omitempty"`
	Error    string `json:"error,omitempty"`
}
