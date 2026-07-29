package handlers

import (
	"sort"
	"sync"
	"testing"

	"github.com/glueops/autoglue/internal/models"
	"gorm.io/gorm/schema"
)

// The handlers in clusters.go address columns by name. GORM resolves a name it
// does not recognise as a column in one of two ways, and only one of them is
// loud:
//
//   - a name matching an association field (ControlPlaneRecordSet, one deleted
//     "ID" away from ControlPlaneRecordSetID) resolves to a field with no
//     column, so the assignment is dropped and the write reports success
//   - anything else is passed through as a quoted identifier and fails at the
//     database with SQLSTATE 42703
//
// The first is the dangerous one: it returns 200 having written nothing, which
// is indistinguishable from a successful detach until a foreign key rejects the
// delete later. These tests are what stands between that and production.
func parseClusterSchema(t *testing.T) *schema.Schema {
	t.Helper()
	s, err := schema.Parse(&models.Cluster{}, &sync.Map{}, schema.NamingStrategy{})
	if err != nil {
		t.Fatalf("parse cluster schema: %v", err)
	}
	return s
}

func TestClusterWritableColumnsAreRealColumns(t *testing.T) {
	s := parseClusterSchema(t)

	for _, col := range clusterWritableColumns {
		f := s.LookUpField(col)
		switch {
		case f == nil:
			t.Errorf("%q does not resolve to any field on models.Cluster", col)
		case f.DBName == "":
			// This is the silent-no-op case.
			t.Errorf("%q resolves to an association, not a column: writing it would be a silent no-op", col)
		case f.DBName != col:
			// Guards the two names that do not follow from the Go field or the
			// JSON tag, e.g. glueops_load_balancer_id vs glue_ops_load_balancer_id.
			t.Errorf("%q resolves to column %q; use the column name", col, f.DBName)
		}
	}
}

func TestClusterImmutableColumnsAreNotWritable(t *testing.T) {
	writable := map[string]bool{}
	for _, c := range clusterWritableColumns {
		writable[c] = true
	}
	for _, c := range clusterImmutableColumns {
		if writable[c] {
			t.Errorf("%q is listed as both writable and immutable", c)
		}
	}
}

// A new column on models.Cluster must be classified deliberately. Without this,
// a column added later is silently neither writable nor immutable, and the
// lists stop describing the table.
func TestClusterColumnListsCoverTheTable(t *testing.T) {
	s := parseClusterSchema(t)

	actual := map[string]bool{}
	for _, f := range s.Fields {
		if f.DBName != "" {
			actual[f.DBName] = true
		}
	}

	classified := map[string]bool{}
	for _, c := range append(append([]string{}, clusterWritableColumns...), clusterImmutableColumns...) {
		classified[c] = true
	}

	var missing, unknown []string
	for col := range actual {
		if !classified[col] {
			missing = append(missing, col)
		}
	}
	for col := range classified {
		if !actual[col] {
			unknown = append(unknown, col)
		}
	}
	sort.Strings(missing)
	sort.Strings(unknown)

	if len(missing) > 0 {
		t.Errorf("columns on models.Cluster that are neither writable nor immutable: %v"+
			" -- add each to clusterWritableColumns or clusterImmutableColumns", missing)
	}
	if len(unknown) > 0 {
		t.Errorf("columns listed but not present on models.Cluster: %v", unknown)
	}
}
