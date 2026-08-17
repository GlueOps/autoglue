package store

import (
	"strings"
	"testing"
)

// store.go claims the one-task-in-flight rule is enforced by the schema rather
// than by AcceptTask. If that is only true of the Go code, the claim is a
// comment that lies, and a future caller writing its own INSERT gets two make
// targets running against one cluster. So: bypass AcceptTask and go straight at
// the table.
func TestOneTaskInFlightIsEnforcedByTheIndexNotTheCode(t *testing.T) {
	s, ctx := open(t)

	insert := func(id, state string) error {
		_, err := s.db.ExecContext(ctx, `
			INSERT INTO tasks (id, command, args_json, state, assigned_at)
			VALUES (?, 'make', '{}', ?, '2026-01-01T00:00:00Z')`, id, state)
		return err
	}

	if err := insert("t1", TaskStarted); err != nil {
		t.Fatalf("first insert: %v", err)
	}

	// A second in-flight row, in either in-flight state, must be refused by the
	// database itself.
	for _, state := range []string{TaskAssigned, TaskStarted} {
		err := insert("t2-"+state, state)
		if err == nil {
			t.Fatalf("raw insert of a second %s task succeeded; the partial unique index is not constraining", state)
		}
		if !strings.Contains(strings.ToLower(err.Error()), "unique") {
			t.Fatalf("refused for the wrong reason: %v", err)
		}
	}

	// Terminal rows are outside the partial index, so any number may coexist —
	// otherwise task history would be capped at one row.
	if err := insert("done-1", TaskSucceeded); err != nil {
		t.Fatalf("terminal insert: %v", err)
	}
	if err := insert("done-2", TaskFailed); err != nil {
		t.Fatalf("second terminal insert: %v", err)
	}
	if err := insert("done-3", TaskDeadLettered); err != nil {
		t.Fatalf("dead-lettered insert: %v", err)
	}
}
