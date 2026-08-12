package bg

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/riverqueue/river"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

// vacuumTables are the tables worth a scheduled vacuum: everything here is
// churned hard enough that autovacuum's default 20% scale factor lets a lot of
// dead tuples accumulate between passes.
//
// river_job is the obvious one — every tick job inserts a row and the cleaner
// deletes it a week later. The leader/queue/client tables are tiny but are
// UPDATEd on a heartbeat, so their dead-tuple ratio is far worse than their row
// count suggests. job_logs and refresh_tokens both take bulk deletes from the
// daily cleanup jobs, which is exactly the pattern that leaves a table full of
// holes.
//
// This list is a constant, never job arguments: VACUUM takes an identifier, and
// an identifier cannot be a bind parameter.
var vacuumTables = []string{
	"river_job",
	"river_leader",
	"river_queue",
	"river_client",
	"river_client_queue",
	"job_logs",
	"refresh_tokens",
	"cluster_runs",
}

type VacuumArgs struct{}

func (VacuumArgs) Kind() string { return "vacuum" }

func (VacuumArgs) InsertOpts() river.InsertOpts {
	// Not retried. A vacuum that failed will be attempted again next month, and
	// a failing one is usually failing for a reason a retry cannot fix
	// (ownership, disk) rather than a transient one.
	return river.InsertOpts{Queue: QueueMaintenance, MaxAttempts: 1}
}

// VacuumTableResult is one table's outcome, recorded on the job row so the
// River dashboard shows what actually ran.
type VacuumTableResult struct {
	Table   string `json:"table"`
	Skipped bool   `json:"skipped,omitempty"`
	Error   string `json:"error,omitempty"`
	Millis  int64  `json:"millis"`
}

type VacuumResult struct {
	Status string              `json:"status"`
	Tables []VacuumTableResult `json:"tables"`
}

// VacuumWorker runs VACUUM (ANALYZE) over the high-churn tables.
//
// Deliberately not VACUUM FULL. FULL rewrites the table under an ACCESS
// EXCLUSIVE lock, which would block every read and write to river_job for the
// duration — on a worker that is fetching jobs from that table, that is a
// self-inflicted outage. Plain VACUUM takes only SHARE UPDATE EXCLUSIVE, so
// normal traffic continues; it makes dead space reusable rather than returning
// it to the OS, which is the right trade for a table that will immediately
// refill it. The ANALYZE is the other half of the point: after a bulk delete
// the planner's row estimates are badly stale.
type VacuumWorker struct {
	river.WorkerDefaults[VacuumArgs]
	db *gorm.DB
}

// Timeout is generous because the cost scales with table size, and this is the
// one job whose whole purpose is to be slow occasionally rather than fast
// always. It stays well under RescueStuckJobsAfter.
func (w *VacuumWorker) Timeout(*river.Job[VacuumArgs]) time.Duration {
	return time.Hour
}

func (w *VacuumWorker) Work(ctx context.Context, _ *river.Job[VacuumArgs]) error {
	results := make([]VacuumTableResult, 0, len(vacuumTables))

	for _, table := range vacuumTables {
		if err := ctx.Err(); err != nil {
			return err
		}

		res := VacuumTableResult{Table: table}
		started := time.Now()

		switch exists, err := tableExists(ctx, w.db, table); {
		case err != nil:
			res.Error = err.Error()
		case !exists:
			// River's schema varies by migration version, and cluster_runs
			// predates some deployments. A table we do not have is not a
			// failure, it is just nothing to do.
			res.Skipped = true
		default:
			res.Error = w.vacuum(ctx, table)
		}

		res.Millis = time.Since(started).Milliseconds()
		results = append(results, res)

		switch {
		case res.Skipped:
			log.Debug().Str("table", table).Msg("[vacuum] table absent, skipped")
		case res.Error != "":
			// One table failing must not cost the rest their vacuum, so this is
			// recorded and the loop continues.
			log.Warn().Str("table", table).Str("error", res.Error).Msg("[vacuum] failed")
		default:
			log.Info().Str("table", table).Int64("ms", res.Millis).Msg("[vacuum] done")
		}
	}

	if err := river.RecordOutput(ctx, VacuumResult{Status: "ok", Tables: results}); err != nil {
		log.Warn().Err(err).Msg("[vacuum] could not record output")
	}
	return nil
}

func (w *VacuumWorker) vacuum(ctx context.Context, table string) string {
	// Sanitize even though the input is a package constant: it costs nothing,
	// and it means adding a name to vacuumTables can never become an injection.
	stmt := fmt.Sprintf("VACUUM (ANALYZE) %s", pgx.Identifier{table}.Sanitize())

	// VACUUM cannot run inside a transaction block, so this must go through
	// Exec on the session rather than any Transaction wrapper.
	if err := w.db.WithContext(ctx).Exec(stmt).Error; err != nil {
		return err.Error()
	}
	return ""
}

func tableExists(ctx context.Context, db *gorm.DB, table string) (bool, error) {
	var reg *string
	err := db.WithContext(ctx).Raw("SELECT to_regclass(?)::text", table).Scan(&reg).Error
	if err != nil {
		return false, err
	}
	return reg != nil, nil
}
