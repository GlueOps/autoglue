package store

// schema is applied on every Open. Every statement is CREATE ... IF NOT EXISTS,
// so this doubles as the migration path for now: the agent is the only writer
// and the file is local, so there is no coordination to do.
//
// The two planes are deliberately separate sets of tables. Config is
// declarative and converges toward a generation; tasks are imperative and run
// once. Mixing them into one "work" table is the mistake this layout exists to
// prevent.
var schema = []string{
	// ---- shared ----

	// agent_meta holds the scalars: generation pointers, the acked log
	// watermark, agent identity. Key/value rather than a one-row table so
	// adding a scalar is not a migration.
	`CREATE TABLE IF NOT EXISTS agent_meta (
		key   TEXT PRIMARY KEY,
		value TEXT NOT NULL
	)`,

	// ---- config plane ----

	// desired_resources is keyed by generation, so a snapshot is inserted
	// wholesale and older generations remain readable until pruned. The server
	// sends a full snapshot rather than a diff, which means the agent never has
	// to reconstruct state from an incomplete history.
	`CREATE TABLE IF NOT EXISTS desired_resources (
		generation      INTEGER NOT NULL,
		resource_type   TEXT    NOT NULL,
		resource_id     TEXT    NOT NULL,
		resource_name   TEXT    NOT NULL DEFAULT '',
		phase           INTEGER NOT NULL DEFAULT 0,
		depends_on_json TEXT    NOT NULL DEFAULT '[]',
		spec_json       TEXT    NOT NULL,
		spec_hash       TEXT    NOT NULL,
		required        INTEGER NOT NULL DEFAULT 1,
		PRIMARY KEY (generation, resource_type, resource_id)
	)`,

	// reconcile_status is the agent's own view of what it has applied. Keyed
	// without generation: there is one current status per resource, and it
	// carries the generation it was achieved at.
	`CREATE TABLE IF NOT EXISTS reconcile_status (
		resource_type      TEXT    NOT NULL,
		resource_id        TEXT    NOT NULL,
		desired_generation INTEGER NOT NULL,
		state              TEXT    NOT NULL,
		last_error         TEXT    NOT NULL DEFAULT '',
		updated_at         TEXT    NOT NULL,
		PRIMARY KEY (resource_type, resource_id)
	)`,

	// ---- task plane ----

	// container_id is singular on purpose, and it is why a task maps to exactly
	// one `docker run`. A task spanning two containers has no representable
	// state when the first succeeded and the second never started, and on
	// restart there is no way to say which one to adopt. The container is what
	// survives the agent dying, so it is what the unit of work is pinned to.
	//
	// logs_cursor is the docker timestamp of the last chunk ingested, so an
	// adopted container resumes its log stream rather than replaying it.
	`CREATE TABLE IF NOT EXISTS tasks (
		id           TEXT PRIMARY KEY,
		command      TEXT    NOT NULL,
		args_json    TEXT    NOT NULL DEFAULT '{}',
		state        TEXT    NOT NULL,
		container_id TEXT    NOT NULL DEFAULT '',
		logs_cursor  TEXT    NOT NULL DEFAULT '',
		exit_code    INTEGER,
		error        TEXT    NOT NULL DEFAULT '',
		assigned_at  TEXT    NOT NULL,
		started_at   TEXT,
		ended_at     TEXT
	)`,

	// At most one task may be in flight. Enforced in the schema rather than in
	// code because it is the invariant the whole task plane rests on: the make
	// targets are supplied rather than authored here, so two of them running
	// against one cluster is not something to discover at runtime.
	//
	// A unique index on a constant expression means every matching row collides
	// with every other, so the partial WHERE decides how many rows may exist.
	`CREATE UNIQUE INDEX IF NOT EXISTS one_task_in_flight
		ON tasks(1) WHERE state IN ('assigned','started')`,

	// ---- outbox ----

	// Everything the agent knows that the server does not yet. This is the
	// reason the store is durable at all: an agent that finishes a task while
	// the API is unreachable must not lose the exit code, because re-running to
	// rediscover it is exactly what non-idempotent targets forbid.

	// seq is a global autoincrement rather than per-task. The server's
	// idempotency key is (task_id, seq), which a global counter satisfies
	// trivially, and a single counter is one less thing to get wrong on resume.
	`CREATE TABLE IF NOT EXISTS log_outbox (
		seq        INTEGER PRIMARY KEY AUTOINCREMENT,
		task_id    TEXT NOT NULL,
		stream     TEXT NOT NULL,
		chunk      BLOB NOT NULL,
		created_at TEXT NOT NULL
	)`,

	`CREATE INDEX IF NOT EXISTS log_outbox_task ON log_outbox(task_id, seq)`,

	// One row per task whose terminal outcome has not been accepted by the
	// server. Keyed by task so a redelivered result cannot duplicate.
	`CREATE TABLE IF NOT EXISTS result_outbox (
		task_id    TEXT PRIMARY KEY,
		state      TEXT    NOT NULL,
		exit_code  INTEGER,
		error      TEXT    NOT NULL DEFAULT '',
		started_at TEXT,
		ended_at   TEXT    NOT NULL,
		created_at TEXT    NOT NULL
	)`,
}

// Task states. assigned and started are the in-flight pair guarded by
// one_task_in_flight; the rest are terminal and release the slot.
const (
	TaskAssigned     = "assigned"
	TaskStarted      = "started"
	TaskSucceeded    = "succeeded"
	TaskFailed       = "failed"
	TaskDeadLettered = "dead_lettered"
)

// Reconcile states for a single resource.
const (
	ReconcilePending = "pending"
	ReconcileApplied = "applied"
	ReconcileFailed  = "failed"
)

// agent_meta keys.
const (
	metaCurrentGeneration = "current_generation"
	metaAppliedGeneration = "applied_generation"
	metaAckedLogSeq       = "acked_log_seq"
)
