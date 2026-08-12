package bg

import (
	"context"
	"database/sql"
	"os"
	"testing"
	"time"

	"github.com/glueops/autoglue/internal/testutil/pgtest"
)

func TestMain(m *testing.M) {
	code := m.Run()
	pgtest.Stop()
	os.Exit(code)
}

// The whole reason this test exists: VACUUM cannot run inside a transaction
// block, and it goes to Postgres through GORM, which is a layer built almost
// entirely around wrapping things in transactions. Whether the statement
// survives that trip is not something to take on faith.
func TestVacuumRunsThroughGorm(t *testing.T) {
	db := pgtest.DB(t)
	ctx := context.Background()
	w := &VacuumWorker{db: db}

	before := lastVacuum(t, "job_logs")

	if msg := w.vacuum(ctx, "job_logs"); msg != "" {
		t.Fatalf("vacuum job_logs: %s", msg)
	}

	// pg_stat is updated by the stats collector, so the timestamp can lag the
	// statement returning by a moment.
	deadline := time.Now().Add(5 * time.Second)
	for {
		if after := lastVacuum(t, "job_logs"); after != nil && (before == nil || after.After(*before)) {
			return
		}
		if time.Now().After(deadline) {
			t.Fatal("last_vacuum for job_logs never advanced")
		}
		time.Sleep(50 * time.Millisecond)
	}
}

// Every name in vacuumTables must be quotable and, if present, vacuumable. A
// typo here would otherwise only show up as a warning line once a month.
func TestVacuumTablesAreAllValid(t *testing.T) {
	db := pgtest.DB(t)
	ctx := context.Background()
	w := &VacuumWorker{db: db}

	for _, table := range vacuumTables {
		exists, err := tableExists(ctx, db, table)
		if err != nil {
			t.Fatalf("tableExists(%s): %v", table, err)
		}
		if !exists {
			// The river_* tables are migrated by River itself, which pgtest
			// does not run. Absence is the case the worker skips.
			continue
		}
		if msg := w.vacuum(ctx, table); msg != "" {
			t.Errorf("vacuum %s: %s", table, msg)
		}
	}
}

func TestTableExistsReportsMissingTables(t *testing.T) {
	db := pgtest.DB(t)
	ctx := context.Background()

	exists, err := tableExists(ctx, db, "job_logs")
	if err != nil {
		t.Fatalf("tableExists(job_logs): %v", err)
	}
	if !exists {
		t.Error("job_logs should exist after migration")
	}

	exists, err = tableExists(ctx, db, "no_such_table_here")
	if err != nil {
		t.Fatalf("tableExists(missing): %v", err)
	}
	if exists {
		t.Error("a table that was never created should not report as existing")
	}
}

func TestMonthlyAtRollsForwardAcrossYearBoundary(t *testing.T) {
	s := monthlyAtSchedule{day: 1, hour: 2, minute: 30}

	cases := []struct {
		name string
		from time.Time
		want time.Time
	}{
		{
			name: "later the same month",
			from: time.Date(2026, time.March, 10, 9, 0, 0, 0, time.UTC),
			want: time.Date(2026, time.April, 1, 2, 30, 0, 0, time.UTC),
		},
		{
			name: "earlier on the target day",
			from: time.Date(2026, time.March, 1, 1, 0, 0, 0, time.UTC),
			want: time.Date(2026, time.March, 1, 2, 30, 0, 0, time.UTC),
		},
		{
			name: "exactly on the target instant does not return itself",
			from: time.Date(2026, time.March, 1, 2, 30, 0, 0, time.UTC),
			want: time.Date(2026, time.April, 1, 2, 30, 0, 0, time.UTC),
		},
		{
			// December + 1 is month 13, which time.Date normalizes.
			name: "december rolls into january",
			from: time.Date(2026, time.December, 15, 0, 0, 0, 0, time.UTC),
			want: time.Date(2027, time.January, 1, 2, 30, 0, 0, time.UTC),
		},
		{
			// The 1st exists in February, unlike the 29th through 31st, which
			// is exactly why the schedule is pinned to it.
			name: "january rolls into february",
			from: time.Date(2026, time.January, 20, 0, 0, 0, 0, time.UTC),
			want: time.Date(2026, time.February, 1, 2, 30, 0, 0, time.UTC),
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := s.Next(tc.from); !got.Equal(tc.want) {
				t.Errorf("Next(%s) = %s, want %s", tc.from, got, tc.want)
			}
		})
	}
}

func lastVacuum(t *testing.T, table string) *time.Time {
	t.Helper()

	// NullTime rather than *time.Time: a table that has never been vacuumed
	// reports NULL here, which is the expected state at the start of this test
	// and not something database/sql will scan into a bare pointer.
	var at sql.NullTime
	err := pgtest.DB(t).
		Raw("SELECT last_vacuum FROM pg_stat_user_tables WHERE relname = ?", table).
		Scan(&at).Error
	if err != nil {
		t.Fatalf("read last_vacuum for %s: %v", table, err)
	}
	if !at.Valid {
		return nil
	}
	return &at.Time
}
