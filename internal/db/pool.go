package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

// OpenPool creates a pgx connection pool for the given DSN.
//
// River talks to Postgres through pgx rather than GORM, so the process keeps a
// pool alongside the *gorm.DB. Connections are established lazily, so this is
// cheap for commands that never touch the queue.
func OpenPool(ctx context.Context, dsn string) (*pgxpool.Pool, error) {
	cfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("parse database url: %w", err)
	}

	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("open pgx pool: %w", err)
	}
	return pool, nil
}
