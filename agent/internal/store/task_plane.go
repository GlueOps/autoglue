package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

// ErrTaskInFlight is returned when accepting a task would break the
// one-at-a-time invariant. It is a normal condition, not a fault: the control
// plane restating an assignment the agent already holds hits this path.
var ErrTaskInFlight = errors.New("a task is already in flight")

type Task struct {
	ID          string
	Command     string
	Args        string // JSON object
	State       string
	ContainerID string
	// LogsCursor is the docker timestamp of the last line ingested, so an
	// adopted container resumes its stream instead of replaying it into the
	// outbox from the beginning.
	LogsCursor string
	ExitCode   *int
	Error      string
}

// SetLogsCursor advances the resume point. Written as output is ingested, so a
// crash costs at most the lines since the last write rather than the whole
// stream.
func (s *Store) SetLogsCursor(ctx context.Context, taskID, cursor string) error {
	if _, err := s.db.ExecContext(ctx,
		`UPDATE tasks SET logs_cursor = ? WHERE id = ?`, cursor, taskID); err != nil {
		return fmt.Errorf("set logs cursor for %s: %w", taskID, err)
	}
	return nil
}

// AcceptTask records a newly assigned task.
//
// Re-delivery is expected rather than exceptional: the call-home is a safe
// method that states the current assignment every time, so the agent will be
// told about its own running task repeatedly. Accepting the same id twice is a
// no-op; accepting a *different* id while one is in flight is refused, because
// that means the agent and the control plane disagree about reality and
// guessing which is right is how two make targets end up running at once.
func (s *Store) AcceptTask(ctx context.Context, t Task) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin accept tx: %w", err)
	}
	defer func() { _ = tx.Rollback() }()

	var existingID string
	err = tx.QueryRowContext(ctx,
		`SELECT id FROM tasks WHERE state IN (?, ?)`, TaskAssigned, TaskStarted).Scan(&existingID)
	switch {
	case err == nil && existingID == t.ID:
		return nil // already held, nothing to do
	case err == nil:
		return fmt.Errorf("%w: holding %s, offered %s", ErrTaskInFlight, existingID, t.ID)
	case !errors.Is(err, sql.ErrNoRows):
		return fmt.Errorf("check in-flight task: %w", err)
	}

	args := t.Args
	if args == "" {
		args = "{}"
	}
	if _, err := tx.ExecContext(ctx, `
		INSERT INTO tasks (id, command, args_json, state, assigned_at)
		VALUES (?, ?, ?, ?, ?)
		ON CONFLICT(id) DO NOTHING`,
		t.ID, t.Command, args, TaskAssigned, nowRFC3339()); err != nil {
		return fmt.Errorf("insert task %s: %w", t.ID, err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit accept: %w", err)
	}
	return nil
}

// StartTask moves a task to started and records the container carrying it. The
// container id is what makes the run adoptable after an agent restart, so it is
// written before the work is reported as underway rather than after.
func (s *Store) StartTask(ctx context.Context, id, containerID string) error {
	res, err := s.db.ExecContext(ctx, `
		UPDATE tasks SET state = ?, container_id = ?, started_at = ?
		WHERE id = ? AND state = ?`,
		TaskStarted, containerID, nowRFC3339(), id, TaskAssigned)
	if err != nil {
		return fmt.Errorf("start task %s: %w", id, err)
	}
	return expectOneRow(res, "start task "+id)
}

// FinishTask writes the terminal state and, in the same transaction, the
// outbox row that owes the server an answer.
//
// One transaction on purpose. If the task were marked terminal without the
// outbox row landing, the agent would free its slot having permanently
// forgotten an outcome it can never recompute — the make targets are supplied,
// not authored here, so re-running to rediscover an exit code is not available.
func (s *Store) FinishTask(ctx context.Context, id, state string, exitCode *int, errMsg string) error {
	switch state {
	case TaskSucceeded, TaskFailed, TaskDeadLettered:
	default:
		return fmt.Errorf("finish task %s: %q is not a terminal state", id, state)
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin finish tx: %w", err)
	}
	defer func() { _ = tx.Rollback() }()

	now := nowRFC3339()
	res, err := tx.ExecContext(ctx, `
		UPDATE tasks SET state = ?, exit_code = ?, error = ?, ended_at = ?
		WHERE id = ? AND state IN (?, ?)`,
		state, exitCode, errMsg, now, id, TaskAssigned, TaskStarted)
	if err != nil {
		return fmt.Errorf("finish task %s: %w", id, err)
	}
	if err := expectOneRow(res, "finish task "+id); err != nil {
		return err
	}

	if _, err := tx.ExecContext(ctx, `
		INSERT INTO result_outbox (task_id, state, exit_code, error, started_at, ended_at, created_at)
		SELECT id, state, exit_code, error, started_at, ended_at, ?
		FROM tasks WHERE id = ?
		ON CONFLICT(task_id) DO NOTHING`, now, id); err != nil {
		return fmt.Errorf("queue result for %s: %w", id, err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit finish: %w", err)
	}
	return nil
}

// InFlightTask returns the task currently held, or nil. On startup this is how
// the agent discovers what it was doing before it died, and the container id is
// what lets it adopt the still-running work rather than restarting it.
func (s *Store) InFlightTask(ctx context.Context) (*Task, error) {
	var t Task
	err := s.db.QueryRowContext(ctx, `
		SELECT id, command, args_json, state, container_id, logs_cursor, exit_code, error
		FROM tasks WHERE state IN (?, ?)`, TaskAssigned, TaskStarted).
		Scan(&t.ID, &t.Command, &t.Args, &t.State, &t.ContainerID, &t.LogsCursor, &t.ExitCode, &t.Error)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("read in-flight task: %w", err)
	}
	return &t, nil
}

// ---- outbox ----

// AppendLog queues a chunk of task output. The returned seq is the server's
// idempotency key alongside the task id, so a redelivered chunk is recognised
// rather than duplicated.
func (s *Store) AppendLog(ctx context.Context, taskID, stream string, chunk []byte) (int64, error) {
	res, err := s.db.ExecContext(ctx, `
		INSERT INTO log_outbox (task_id, stream, chunk, created_at) VALUES (?, ?, ?, ?)`,
		taskID, stream, chunk, nowRFC3339())
	if err != nil {
		return 0, fmt.Errorf("append log for %s: %w", taskID, err)
	}
	seq, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("read log seq: %w", err)
	}
	return seq, nil
}

type LogChunk struct {
	Seq    int64
	TaskID string
	Stream string
	Chunk  []byte
}

// PendingLogs returns up to limit unacked chunks in sequence order. Order is
// the contract: the server appends them under an autoincrement id that readers
// page through with `id > after`, so chunks arriving out of order would leave a
// hole no reader ever revisits.
func (s *Store) PendingLogs(ctx context.Context, limit int) ([]LogChunk, error) {
	acked, err := s.getMetaInt64(ctx, metaAckedLogSeq)
	if err != nil {
		return nil, err
	}

	rows, err := s.db.QueryContext(ctx, `
		SELECT seq, task_id, stream, chunk FROM log_outbox
		WHERE seq > ? ORDER BY seq LIMIT ?`, acked, limit)
	if err != nil {
		return nil, fmt.Errorf("read pending logs: %w", err)
	}
	defer func() { _ = rows.Close() }()

	var out []LogChunk
	for rows.Next() {
		var c LogChunk
		if err := rows.Scan(&c.Seq, &c.TaskID, &c.Stream, &c.Chunk); err != nil {
			return nil, fmt.Errorf("scan log chunk: %w", err)
		}
		out = append(out, c)
	}
	return out, rows.Err()
}

// AckLogs advances the watermark and drops what it covers.
//
// The watermark is stored rather than inferred from what remains, so an agent
// that restarts mid-drain resumes from the right place instead of replaying
// from zero.
func (s *Store) AckLogs(ctx context.Context, throughSeq int64) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin ack tx: %w", err)
	}
	defer func() { _ = tx.Rollback() }()

	if _, err := tx.ExecContext(ctx, `
		INSERT INTO agent_meta(key, value) VALUES(?, ?)
		ON CONFLICT(key) DO UPDATE SET value = excluded.value
		WHERE CAST(excluded.value AS INTEGER) > CAST(agent_meta.value AS INTEGER)`,
		metaAckedLogSeq, fmt.Sprint(throughSeq)); err != nil {
		return fmt.Errorf("advance log watermark: %w", err)
	}
	if _, err := tx.ExecContext(ctx,
		`DELETE FROM log_outbox WHERE seq <= ?`, throughSeq); err != nil {
		return fmt.Errorf("trim log outbox: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit ack: %w", err)
	}
	return nil
}

type PendingResult struct {
	TaskID   string
	State    string
	ExitCode *int
	Error    string
}

// PendingResults returns terminal outcomes the server has not accepted yet.
func (s *Store) PendingResults(ctx context.Context) ([]PendingResult, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT task_id, state, exit_code, error FROM result_outbox ORDER BY created_at`)
	if err != nil {
		return nil, fmt.Errorf("read pending results: %w", err)
	}
	defer func() { _ = rows.Close() }()

	var out []PendingResult
	for rows.Next() {
		var r PendingResult
		if err := rows.Scan(&r.TaskID, &r.State, &r.ExitCode, &r.Error); err != nil {
			return nil, fmt.Errorf("scan pending result: %w", err)
		}
		out = append(out, r)
	}
	return out, rows.Err()
}

// AckResult drops a delivered outcome.
func (s *Store) AckResult(ctx context.Context, taskID string) error {
	if _, err := s.db.ExecContext(ctx,
		`DELETE FROM result_outbox WHERE task_id = ?`, taskID); err != nil {
		return fmt.Errorf("ack result %s: %w", taskID, err)
	}
	return nil
}

func expectOneRow(res sql.Result, what string) error {
	n, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s: rows affected: %w", what, err)
	}
	if n == 0 {
		return fmt.Errorf("%s: no row in the expected state", what)
	}
	return nil
}
