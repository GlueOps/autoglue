package bg

import (
	"context"
	"log/slog"
	"os"
	"time"

	"github.com/glueops/autoglue/internal/config"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/riverqueue/river"
	"github.com/riverqueue/river/riverdriver/riverpgxv5"
	"github.com/riverqueue/river/rivertype"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

// Queue names. Long-running cluster work is kept off the maintenance queue so a
// week-long bootstrap can never starve the hourly sweepers.
const (
	QueueDefault     = river.QueueDefault
	QueueMaintenance = "maintenance"
	QueueClusters    = "clusters"
)

// Deps is everything the workers need to run. The API process does not build
// this: it only ever inserts.
type Deps struct {
	DB      *gorm.DB
	BaseURL string
}

// Client is the River client type used throughout. Aliased so callers do not
// have to spell the pgx transaction parameter.
type Client = river.Client[pgx.Tx]

// NewInsertClient builds an insert-only client for the API process. It has no
// workers and no queues, so it never fetches work: calling Start on it is
// neither required nor useful.
func NewInsertClient(pool *pgxpool.Pool) (*Client, error) {
	return river.NewClient(riverpgxv5.New(pool), &river.Config{
		Logger: newLogger(),
	})
}

// NewWorkerClient builds a fully configured client for the worker process:
// registered workers, queue concurrency, and the periodic schedule.
func NewWorkerClient(pool *pgxpool.Pool, d Deps) (*Client, error) {
	workers := river.NewWorkers()

	river.AddWorker(workers, &BastionSweepWorker{db: d.DB})
	river.AddWorker(workers, &BastionBootstrapWorker{db: d.DB})
	river.AddWorker(workers, &ClusterActionWorker{db: d.DB, baseURL: d.BaseURL})
	river.AddWorker(workers, &DNSReconcileWorker{db: d.DB})
	river.AddWorker(workers, &DbBackupWorker{db: d.DB})
	river.AddWorker(workers, &JobLogsCleanupWorker{db: d.DB})
	river.AddWorker(workers, &OrgKeySweeperWorker{db: d.DB})
	river.AddWorker(workers, &TokensCleanupWorker{db: d.DB})
	river.AddWorker(workers, &VacuumWorker{db: d.DB})

	return river.NewClient(riverpgxv5.New(pool), &river.Config{
		Logger:  newLogger(),
		Workers: workers,
		Queues: map[string]river.QueueConfig{
			QueueDefault:     {MaxWorkers: 10},
			QueueMaintenance: {MaxWorkers: 2},
			QueueClusters:    {MaxWorkers: 30},
		},
		PeriodicJobs: periodicJobs(),
		// The rescuer reclaims jobs whose attempt started longer ago than this,
		// and it is purely time based — it has no idea whether a worker is
		// still alive and working. The default is one hour, which would reap a
		// cluster_action mid-bootstrap and, because that job is not retryable,
		// discard it outright. Keep the window above the longest worker Timeout.
		RescueStuckJobsAfter: rescueWindow(),
		// Archer had a bespoke cleanup job for this. River expires finished
		// jobs itself, so the retention windows live here instead.
		CompletedJobRetentionPeriod: retention("river.completed_retain_days", 7),
		CancelledJobRetentionPeriod: retention("river.cancelled_retain_days", 7),
		DiscardedJobRetentionPeriod: retention("river.discarded_retain_days", 14),
	})
}

// periodicJobs replaces archer's self-rescheduling chain, where every worker
// re-enqueued itself as its last act. River runs this schedule on the elected
// leader only, so scaling the worker deployment past one replica does not
// double-insert.
func periodicJobs() []*river.PeriodicJob {
	return []*river.PeriodicJob{
		// The sweep only claims and dispatches; the bootstraps themselves run
		// as one job per server on the clusters queue, which is what restores
		// archer's 30-way concurrency after the port collapsed it to a single
		// serial tick.
		river.NewPeriodicJob(
			river.PeriodicInterval(interval("bastion.interval_seconds", 30*time.Second)),
			func() (river.JobArgs, *river.InsertOpts) {
				return BastionSweepArgs{}, &river.InsertOpts{UniqueOpts: tickUnique}
			},
			&river.PeriodicJobOpts{ID: "bastion_sweep", RunOnStart: true},
		),
		river.NewPeriodicJob(
			river.PeriodicInterval(interval("dns.interval_seconds", 30*time.Second)),
			func() (river.JobArgs, *river.InsertOpts) {
				return DNSReconcileArgs{}, &river.InsertOpts{UniqueOpts: tickUnique}
			},
			&river.PeriodicJobOpts{ID: "dns_reconcile", RunOnStart: true},
		),
		river.NewPeriodicJob(
			river.PeriodicInterval(interval("org_key_sweeper.interval_seconds", time.Hour)),
			func() (river.JobArgs, *river.InsertOpts) {
				return OrgKeySweeperArgs{}, nil
			},
			&river.PeriodicJobOpts{ID: "org_key_sweeper", RunOnStart: true},
		),
		river.NewPeriodicJob(
			river.PeriodicInterval(interval("backup.interval_seconds", time.Hour)),
			func() (river.JobArgs, *river.InsertOpts) {
				return DbBackupArgs{}, nil
			},
			&river.PeriodicJobOpts{ID: "db_backup_s3"},
		),
		river.NewPeriodicJob(
			dailyAt(4, 15),
			func() (river.JobArgs, *river.InsertOpts) {
				return JobLogsCleanupArgs{}, nil
			},
			&river.PeriodicJobOpts{ID: "job_logs_cleanup"},
		),
		river.NewPeriodicJob(
			dailyAt(3, 45),
			func() (river.JobArgs, *river.InsertOpts) {
				return TokensCleanupArgs{}, nil
			},
			&river.PeriodicJobOpts{ID: "tokens_cleanup"},
		),
		// First of the month, at an hour that clears the daily cleanups above:
		// their bulk deletes are exactly what this is tidying up after, so
		// running into them would waste the pass.
		river.NewPeriodicJob(
			monthlyAt(1, 2, 30),
			func() (river.JobArgs, *river.InsertOpts) {
				return VacuumArgs{}, nil
			},
			&river.PeriodicJobOpts{ID: "vacuum"},
		),
	}
}

// tickUnique stops reconcile ticks from stacking up. Under archer each worker
// enqueued its own successor, so exactly one tick existed at a time. River's
// PeriodicInterval inserts on a wall clock regardless of whether the previous
// tick finished, so an 8 minute bastion bootstrap would otherwise queue dozens
// of redundant ticks against the same queue cluster_action runs on.
var tickUnique = river.UniqueOpts{
	ByArgs:  true,
	ByState: []rivertype.JobState{rivertype.JobStateAvailable, rivertype.JobStateScheduled, rivertype.JobStateRunning, rivertype.JobStateRetryable, rivertype.JobStatePending},
}

// dailyAtSchedule fires once per day at a fixed local wall-clock time.
type dailyAtSchedule struct {
	hour   int
	minute int
}

func (s dailyAtSchedule) Next(t time.Time) time.Time {
	next := time.Date(t.Year(), t.Month(), t.Day(), s.hour, s.minute, 0, 0, t.Location())
	if !next.After(t) {
		next = next.AddDate(0, 0, 1)
	}
	return next
}

func dailyAt(hour, minute int) river.PeriodicSchedule {
	return dailyAtSchedule{hour: hour, minute: minute}
}

// monthlyAtSchedule fires once per month on a fixed day at a fixed local
// wall-clock time. Keep day at or below 28: time.Date normalizes an overflowing
// day into the next month, so day 31 in February would silently drift to March.
type monthlyAtSchedule struct {
	day    int
	hour   int
	minute int
}

func (s monthlyAtSchedule) Next(t time.Time) time.Time {
	next := time.Date(t.Year(), t.Month(), s.day, s.hour, s.minute, 0, 0, t.Location())
	if !next.After(t) {
		// Month+1 is normalized by time.Date, so December rolls into January of
		// the following year without special-casing.
		next = time.Date(t.Year(), t.Month()+1, s.day, s.hour, s.minute, 0, 0, t.Location())
	}
	return next
}

func monthlyAt(day, hour, minute int) river.PeriodicSchedule {
	return monthlyAtSchedule{day: day, hour: hour, minute: minute}
}

// maxWorkerTimeout is the longest Timeout any registered worker returns. Keep
// this in step with the Timeout methods; ClusterActionWorker is the outlier.
const maxWorkerTimeout = 168 * time.Hour

// rescueWindow keeps stuck-job rescue comfortably clear of legitimately
// long-running work, while still reclaiming genuinely abandoned jobs.
func rescueWindow() time.Duration {
	if h := viper.GetInt("river.rescue_after_hours"); h > 0 {
		return time.Duration(h) * time.Hour
	}
	return maxWorkerTimeout + time.Hour
}

func interval(key string, def time.Duration) time.Duration {
	if s := viper.GetInt(key); s > 0 {
		return time.Duration(s) * time.Second
	}
	return def
}

// jobLogRetainDays is how long a job transcript is kept. Much longer than
// River's own job retention: the job row is a status, but the transcript is the
// evidence, and it is what someone comes back for weeks after a bootstrap went
// wrong. The row it describes will already have been expired by then.
func jobLogRetainDays() int {
	if d := viper.GetInt("job_logs.retain_days"); d > 0 {
		return d
	}
	return 45
}

func retention(key string, defDays int) time.Duration {
	days := viper.GetInt(key)
	if days <= 0 {
		days = defDays
	}
	return time.Duration(days) * 24 * time.Hour
}

func newLogger() *slog.Logger {
	level := slog.LevelInfo
	if config.IsDebug() {
		level = slog.LevelDebug
	}
	return slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: level}))
}

// Insert is a thin helper so callers outside this package do not need to import
// River directly for the common case.
func Insert(ctx context.Context, c *Client, args river.JobArgs, opts *river.InsertOpts) error {
	_, err := c.Insert(ctx, args, opts)
	return err
}
