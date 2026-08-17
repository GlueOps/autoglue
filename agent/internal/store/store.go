// Package store is the agent's durable local state.
//
// It holds two planes that never mix: the config plane, which converges toward
// a desired generation, and the task plane, which runs one imperative task at a
// time. It also holds the outbox — the results and log chunks the agent has
// produced but the server has not yet accepted.
//
// The outbox is why this is SQLite rather than a JSON file. An agent that
// completes a task while the API is unreachable must not lose the exit code:
// re-running to rediscover it is precisely what non-idempotent make targets
// forbid, so the record has to survive a crash mid-write.
package store

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"time"

	// Pure-Go SQLite. Deliberately not mattn/go-sqlite3: this binary is
	// cross-compiled and shipped to bastions, and cgo would make that a
	// toolchain problem on every target.
	_ "modernc.org/sqlite"
)

type Store struct {
	db *sql.DB
}

// Open opens (creating if needed) the agent database at path and applies the
// schema.
func Open(ctx context.Context, path string) (*Store, error) {
	// WAL so a reader never blocks the writer, and busy_timeout so the two
	// plane loops queue instead of returning SQLITE_BUSY at each other.
	// foreign_keys is on for the day the schema grows one.
	dsn := path + "?_pragma=journal_mode(WAL)&_pragma=busy_timeout(5000)&_pragma=foreign_keys(on)"

	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, fmt.Errorf("open sqlite at %s: %w", path, err)
	}

	// One writer. SQLite serialises writes anyway, and letting the pool open
	// several connections only converts that into lock contention the busy
	// timeout then has to absorb.
	db.SetMaxOpenConns(1)

	s := &Store{db: db}
	if err := s.migrate(ctx); err != nil {
		_ = db.Close()
		return nil, err
	}
	return s, nil
}

func (s *Store) Close() error { return s.db.Close() }

func (s *Store) migrate(ctx context.Context) error {
	for _, stmt := range schema {
		if _, err := s.db.ExecContext(ctx, stmt); err != nil {
			return fmt.Errorf("apply schema: %w\nstatement: %s", err, stmt)
		}
	}
	return nil
}

// ---- meta scalars ----

func (s *Store) getMetaInt64(ctx context.Context, key string) (int64, error) {
	var raw string
	err := s.db.QueryRowContext(ctx, `SELECT value FROM agent_meta WHERE key = ?`, key).Scan(&raw)
	if err == sql.ErrNoRows {
		// Absent means "never set", which for every counter here is zero.
		return 0, nil
	}
	if err != nil {
		return 0, fmt.Errorf("read meta %s: %w", key, err)
	}
	n, err := strconv.ParseInt(raw, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("meta %s is not an integer (%q): %w", key, raw, err)
	}
	return n, nil
}

func (s *Store) setMetaInt64(ctx context.Context, key string, n int64) error {
	_, err := s.db.ExecContext(ctx, `
		INSERT INTO agent_meta(key, value) VALUES(?, ?)
		ON CONFLICT(key) DO UPDATE SET value = excluded.value`,
		key, strconv.FormatInt(n, 10))
	if err != nil {
		return fmt.Errorf("write meta %s: %w", key, err)
	}
	return nil
}

// CurrentGeneration is the newest desired snapshot the agent has stored.
func (s *Store) CurrentGeneration(ctx context.Context) (int64, error) {
	return s.getMetaInt64(ctx, metaCurrentGeneration)
}

// AppliedGeneration is the newest snapshot the agent has fully converged to.
// It trails CurrentGeneration while reconciliation is in progress, and the gap
// between the two is what the server reports as "not yet converged".
func (s *Store) AppliedGeneration(ctx context.Context) (int64, error) {
	return s.getMetaInt64(ctx, metaAppliedGeneration)
}

func (s *Store) SetAppliedGeneration(ctx context.Context, gen int64) error {
	return s.setMetaInt64(ctx, metaAppliedGeneration, gen)
}

func nowRFC3339() string { return time.Now().UTC().Format(time.RFC3339Nano) }
