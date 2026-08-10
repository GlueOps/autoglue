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

	river.AddWorker(workers, &BastionBootstrapWorker{db: d.DB})
	river.AddWorker(workers, &ClusterActionWorker{db: d.DB, baseURL: d.BaseURL})
	river.AddWorker(workers, &DNSReconcileWorker{db: d.DB})
	river.AddWorker(workers, &DbBackupWorker{db: d.DB})
	river.AddWorker(workers, &OrgKeySweeperWorker{db: d.DB})
	river.AddWorker(workers, &TokensCleanupWorker{db: d.DB})

	return river.NewClient(riverpgxv5.New(pool), &river.Config{
		Logger:  newLogger(),
		Workers: workers,
		Queues: map[string]river.QueueConfig{
			QueueDefault:     {MaxWorkers: 10},
			QueueMaintenance: {MaxWorkers: 2},
			QueueClusters:    {MaxWorkers: 30},
		},
		PeriodicJobs: periodicJobs(),
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
		river.NewPeriodicJob(
			river.PeriodicInterval(interval("bastion.interval_seconds", 10*time.Second)),
			func() (river.JobArgs, *river.InsertOpts) {
				return BastionBootstrapArgs{}, nil
			},
			&river.PeriodicJobOpts{ID: "bootstrap_bastion", RunOnStart: true},
		),
		river.NewPeriodicJob(
			river.PeriodicInterval(interval("dns.interval_seconds", 10*time.Second)),
			func() (river.JobArgs, *river.InsertOpts) {
				return DNSReconcileArgs{}, nil
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
			dailyAt(3, 45),
			func() (river.JobArgs, *river.InsertOpts) {
				return TokensCleanupArgs{}, nil
			},
			&river.PeriodicJobOpts{ID: "tokens_cleanup"},
		),
	}
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

func interval(key string, def time.Duration) time.Duration {
	if s := viper.GetInt(key); s > 0 {
		return time.Duration(s) * time.Second
	}
	return def
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
