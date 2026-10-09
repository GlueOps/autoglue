package bg

import (
	"context"
	"encoding/json"
	"sync"
	"testing"
	"time"

	"github.com/glueops/autoglue/internal/testutil/pgtest"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/riverqueue/river"
	"github.com/riverqueue/river/riverdriver/riverpgxv5"
	"github.com/riverqueue/river/rivermigrate"
	"github.com/riverqueue/river/rivertest"
)

var (
	riverPoolOnce sync.Once
	riverPool     *pgxpool.Pool
	riverPoolErr  error
)

// testRiverPool returns a pgx pool on the embedded Postgres with River's
// schema migrated, so workers can be run through River's own executor.
func testRiverPool(t *testing.T) *pgxpool.Pool {
	t.Helper()
	url := pgtest.URL(t)
	riverPoolOnce.Do(func() {
		ctx := context.Background()
		riverPool, riverPoolErr = pgxpool.New(ctx, url)
		if riverPoolErr != nil {
			return
		}
		m, err := rivermigrate.New(riverpgxv5.New(riverPool), nil)
		if err != nil {
			riverPoolErr = err
			return
		}
		_, riverPoolErr = m.Migrate(ctx, rivermigrate.DirectionUp, nil)
	})
	if riverPoolErr != nil {
		t.Fatalf("river pool: %v", riverPoolErr)
	}
	return riverPool
}

func testTx(t *testing.T, pool *pgxpool.Pool) pgx.Tx {
	t.Helper()
	tx, err := pool.Begin(context.Background())
	if err != nil {
		t.Fatalf("begin: %v", err)
	}
	t.Cleanup(func() { _ = tx.Rollback(context.Background()) })
	return tx
}

// Through River's real executor: the metadata the bootstrap sets on a snooze
// must survive to the next attempt, or the session-lost budget and the
// known-host decision would reset every time.
func TestBastionBootstrapMetadataSurvivesSnooze(t *testing.T) {
	f := newBastionFixture(t)
	pool := testRiverPool(t)
	tx := testTx(t, pool)
	ctx := context.Background()

	var dropped bool
	var mu sync.Mutex
	listenFakeSSHD(t, onlyKey(f.pub), func() string {
		mu.Lock()
		defer mu.Unlock()
		if !dropped {
			dropped = true
			return execDrop
		}
		return execOK
	})

	tw := rivertest.NewWorker(t, riverpgxv5.New(pool), &river.Config{}, &BastionBootstrapWorker{db: f.db})
	res, err := tw.Work(ctx, t, tx, BastionBootstrapArgs{ServerID: f.server.ID}, nil)
	if err != nil {
		t.Fatalf("first attempt: %v", err)
	}
	if res.EventKind != river.EventKindJobSnoozed {
		t.Fatalf("first attempt event = %s, want snoozed", res.EventKind)
	}
	var meta map[string]any
	if err := json.Unmarshal(res.Job.Metadata, &meta); err != nil {
		t.Fatalf("metadata: %v", err)
	}
	if meta[bastionMetaKnownHost] != false || meta[bastionMetaSessionLost] != float64(1) || meta[bastionMetaSessionLostAt] == nil {
		t.Fatalf("metadata after lost session = %s", res.Job.Metadata)
	}
	if got := f.status(t); got != "provisioning" {
		t.Fatalf("status = %q, want provisioning", got)
	}

	res, err = tw.WorkJob(ctx, t, tx, res.Job)
	if err != nil {
		t.Fatalf("second attempt: %v", err)
	}
	if res.EventKind != river.EventKindJobCompleted {
		t.Fatalf("second attempt event = %s, want completed", res.EventKind)
	}
	if got := f.status(t); got != "ready" {
		t.Fatalf("status = %q, want ready", got)
	}
}

// A server set back to pending while its bootstrap is snoozed is claimed by
// the sweep, but the dispatch dedupes against the live job; it must be
// counted as a duplicate, not as dispatched. Uses a real client: rivertest
// switches unique enforcement off.
func TestBastionDispatchCountsDuplicates(t *testing.T) {
	f := newBastionFixture(t)
	pool := testRiverPool(t)
	ctx := context.Background()
	t.Cleanup(func() {
		_, _ = pool.Exec(ctx, `DELETE FROM river_job WHERE kind = $1 AND args->>'server_id' = $2`,
			BastionBootstrapArgs{}.Kind(), f.server.ID.String())
	})

	client, err := NewInsertClient(pool)
	if err != nil {
		t.Fatalf("insert client: %v", err)
	}
	// The live, snoozed bootstrap.
	if _, err := client.Insert(ctx, BastionBootstrapArgs{ServerID: f.server.ID},
		&river.InsertOpts{ScheduledAt: time.Now().Add(time.Hour)}); err != nil {
		t.Fatalf("insert live job: %v", err)
	}
	other := newBastionFixture(t)
	t.Cleanup(func() {
		_, _ = pool.Exec(ctx, `DELETE FROM river_job WHERE kind = $1 AND args->>'server_id' = $2`,
			BastionBootstrapArgs{}.Kind(), other.server.ID.String())
	})

	dispatched, duplicates := dispatchBastionBootstraps(ctx, client, f.db,
		[]uuid.UUID{f.server.ID, other.server.ID})
	if dispatched != 1 || duplicates != 1 {
		t.Fatalf("dispatched, duplicates = %d, %d; want 1, 1", dispatched, duplicates)
	}
	if got := f.status(t); got != "provisioning" {
		t.Fatalf("status = %q, want provisioning (the live job owns it)", got)
	}
}
