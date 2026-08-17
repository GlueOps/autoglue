package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

// Task kinds. A container task is one `docker run` and nothing else; see the
// note on ContainerID below for why that is a hard rule rather than a default.
// Host tasks are the escape hatch for work that cannot be containerised, and
// carry the same one-at-a-time and terminal-report contract.
const (
	AgentTaskKindContainer = "container"
	AgentTaskKindHost      = "host"
)

// Task states. assigned and started are the in-flight pair; the other three are
// terminal and are the only things that release the agent's slot.
//
// dead_lettered is a third outcome, not a flavour of failure. It means nobody
// can say whether the work happened — which is not the same as knowing it did
// not — and it must never be retried automatically. Without it, "never
// auto-retry" plus "only a terminal report frees the slot" deadlocks an agent
// permanently on the first ambiguous container.
const (
	AgentTaskStateAssigned     = "assigned"
	AgentTaskStateStarted      = "started"
	AgentTaskStateSucceeded    = "succeeded"
	AgentTaskStateFailed       = "failed"
	AgentTaskStateDeadLettered = "dead_lettered"
)

// Dead-letter reasons. Recorded rather than folded into Error because the
// reason decides who can act on it: an agent-reported dead letter needs a human
// to inspect the cluster, a re-enrolment one needs nothing but a fresh run.
const (
	AgentTaskDeadLetterAgentReported = "agent_reported"
	AgentTaskDeadLetterReenrolled    = "agent_reenrolled"
	AgentTaskDeadLetterOperator      = "operator"
	// AgentTaskDeadLetterAgentLost is written by the sweeper when the agent
	// holding a task has stopped calling home. Completed-but-unreported and
	// half-applied are indistinguishable from the control plane, which is why
	// this is a dead letter and never a retry.
	AgentTaskDeadLetterAgentLost = "agent_lost"
)

// AgentTask is one unit of imperative work handed to one agent: exactly one
// container run, or one host command.
//
// A task is never popped from a queue. The assignment endpoint states this row
// as a pure function of durable state, so an unacknowledged task is simply
// named again on the next poll and a lost response costs nothing. Hand-out
// proves nothing; only a reported terminal state does.
type AgentTask struct {
	ID uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`

	OrganizationID uuid.UUID `gorm:"type:uuid;not null;index" json:"organization_id"`
	ClusterID      uuid.UUID `gorm:"type:uuid;not null;index" json:"cluster_id"`

	// AgentID plus the partial unique index is the one-task-at-a-time
	// invariant, mirroring the agent's own `one_task_in_flight` index. Enforced
	// here too rather than trusted from the agent: the make targets are
	// supplied rather than authored here, so two of them running against one
	// cluster is not something to discover at runtime.
	AgentID uuid.UUID `gorm:"type:uuid;not null;index;uniqueIndex:uniq_agent_tasks_in_flight,where:state = 'assigned' OR state = 'started'" json:"agent_id"`
	Agent   Agent     `gorm:"foreignKey:AgentID;constraint:OnDelete:CASCADE" json:"-"`

	// RunID is the ClusterRun this task advances, and it is unique: one run is
	// one make target is one container is one task. That uniqueness is what
	// makes assignment derivable — deriving it becomes an upsert that cannot
	// mint a second task for a run no matter how often it is polled.
	//
	// Deliberately no navigation field. A run row is the audit anchor a task
	// hangs off, so cascading a run delete into the task would destroy the
	// record of a dead letter; cluster deletion already reaches this row
	// through AgentID.
	RunID uuid.UUID `gorm:"type:uuid;not null;uniqueIndex" json:"run_id"`

	Kind string `gorm:"type:text;not null;default:'container'" json:"kind"`

	// Command is the make target, copied from Action.MakeTarget the same way
	// ClusterRun.Action is: it is what actually executes, and the catalogue row
	// may be edited or removed afterwards.
	Command string         `gorm:"type:text;not null" json:"command"`
	Args    datatypes.JSON `gorm:"type:jsonb;not null;default:'{}'" json:"args"`

	State string `gorm:"type:text;not null;default:'assigned';index" json:"state"`

	// ContainerID is singular on purpose, and it is why a task maps to exactly
	// one container. A task spanning two has no representable state when the
	// first succeeded and the second never started, and nothing on either side
	// can say which to adopt after a restart. The container is what survives
	// the agent dying, so it is what the unit of work is pinned to.
	ContainerID string `gorm:"type:text;not null;default:''" json:"container_id"`

	// LogSeq is the highest per-task chunk sequence accepted. It is the
	// idempotency key for log ingest: chunks arrive in seq order per task, so a
	// watermark rejects a redelivered batch without needing a unique index on
	// the high-volume job_logs table, and without a dedupe table to sweep.
	LogSeq int64 `gorm:"not null;default:0" json:"log_seq"`

	ExitCode         *int   `json:"exit_code,omitempty"`
	Error            string `gorm:"type:text;not null;default:''" json:"error"`
	DeadLetterReason string `gorm:"type:text;not null;default:''" json:"dead_letter_reason"`

	AssignedAt time.Time  `gorm:"type:timestamptz;not null;default:now()" json:"assigned_at"`
	StartedAt  *time.Time `gorm:"type:timestamptz" json:"started_at,omitempty"`
	EndedAt    *time.Time `gorm:"type:timestamptz" json:"ended_at,omitempty"`

	CreatedAt time.Time `json:"created_at,omitempty" gorm:"type:timestamptz;column:created_at;not null;default:now()"`
	UpdatedAt time.Time `json:"updated_at,omitempty" gorm:"type:timestamptz;autoUpdateTime;column:updated_at;not null;default:now()"`
}
