package handlers

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/glueops/autoglue/internal/config"
	"github.com/glueops/autoglue/internal/models"
	"github.com/glueops/autoglue/internal/testutil/pgtest"
	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// TestClusterDetach_DoesNotClobberAConcurrentSiblingWrite is the regression
// test for the bug this file's handlers were changed to fix.
//
// OpenTofu destroys a cluster's attachments in parallel, so several detach
// handlers run against the same row at once. When each one saved the whole
// row, the last writer replayed every column from a snapshot taken before the
// others had written, silently restoring a foreign key another handler had just
// cleared. The delete of the referenced row then failed on that foreign key.
//
// Rather than race goroutines and hope, this drives the interleaving from a
// GORM callback: the competing write is issued at the exact point between the
// handler's read and its write.
func TestClusterDetach_DoesNotClobberAConcurrentSiblingWrite(t *testing.T) {
	shared := pgtest.DB(t)
	dsn := pgtest.URL(t)
	cfg := config.Config{}

	for _, tc := range fkAttachCases() {
		t.Run(tc.name, func(t *testing.T) {
			// Any column the handler under test does not own.
			sibling := colClusterCaptainDomainID
			if tc.column == sibling {
				sibling = colClusterBastionServerID
			}

			org := createTestOrg(t, shared, "clobber-"+tc.name)
			cluster := newAttachCluster(t, shared, org.ID)

			own := tc.newTarget(t, shared, org.ID)
			var siblingTarget uuid.UUID
			if sibling == colClusterCaptainDomainID {
				siblingTarget = newTestDomain(t, shared, org.ID).ID
			} else {
				siblingTarget = newTestServer(t, shared, org.ID).ID
			}

			if err := shared.Model(&models.Cluster{}).Where("id = ?", cluster.ID).
				Updates(map[string]any{tc.column: own, sibling: siblingTarget}).Error; err != nil {
				t.Fatalf("seed attachments: %v", err)
			}

			// A dedicated handle: callbacks are registered on the *gorm.DB and
			// would otherwise leak into every other test sharing the pgtest one.
			ded, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
			if err != nil {
				t.Fatalf("open dedicated handle: %v", err)
			}

			fired := 0
			err = ded.Callback().Update().Before("gorm:update").
				Register("test:interleave", func(tx *gorm.DB) {
					if tx.Statement.Table != "clusters" || fired > 0 {
						return
					}
					fired++
					// Issued on a different connection, so this is a genuinely
					// separate transaction. Never use tx here: that would run
					// inside the handler's own transaction and model nothing.
					q := fmt.Sprintf("UPDATE clusters SET %s = NULL WHERE id = ?", sibling)
					if e := shared.Exec(q, cluster.ID).Error; e != nil {
						t.Errorf("competing write: %v", e)
					}
				})
			if err != nil {
				t.Fatalf("register callback: %v", err)
			}

			rr := httptest.NewRecorder()
			tc.detach(ded, cfg).ServeHTTP(rr, clusterReq(http.MethodDelete, "", &org.ID, cluster.ID.String()))
			if rr.Code != http.StatusOK {
				t.Fatalf("expected 200, got %d body=%s", rr.Code, rr.Body.String())
			}

			// Without this the test silently proves nothing if the callback
			// stops firing (a GORM upgrade, a renamed processor, a handler that
			// stops issuing an UPDATE at all).
			if fired != 1 {
				t.Fatalf("interleave fired %d times, want 1 -- the test is vacuous", fired)
			}

			if got := clusterColumn(t, shared, cluster.ID, tc.column); got != nil {
				t.Errorf("%s: handler did not clear its own column, got %q", tc.column, *got)
			}
			if got := clusterColumn(t, shared, cluster.ID, sibling); got != nil {
				t.Errorf("%s: the concurrent write was clobbered -- the handler replayed a stale value (%q)",
					sibling, *got)
			}
		})
	}
}

// A detach arriving after the cluster row is gone must not recreate it.
// db.Save falls back to an INSERT when its UPDATE matches no rows, which
// resurrected deleted clusters with every stale foreign key restored.
func TestClusterDetach_DoesNotResurrectADeletedCluster(t *testing.T) {
	shared := pgtest.DB(t)
	dsn := pgtest.URL(t)
	cfg := config.Config{}

	org := createTestOrg(t, shared, "resurrect")
	cluster := newAttachCluster(t, shared, org.ID)

	ded, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("open dedicated handle: %v", err)
	}

	fired := 0
	err = ded.Callback().Update().Before("gorm:update").
		Register("test:delete-underneath", func(tx *gorm.DB) {
			if tx.Statement.Table != "clusters" || fired > 0 {
				return
			}
			fired++
			if e := shared.Exec("DELETE FROM clusters WHERE id = ?", cluster.ID).Error; e != nil {
				t.Errorf("delete underneath: %v", e)
			}
		})
	if err != nil {
		t.Fatalf("register callback: %v", err)
	}

	rr := httptest.NewRecorder()
	DetachCaptainDomain(ded, cfg).ServeHTTP(rr, clusterReq(http.MethodDelete, "", &org.ID, cluster.ID.String()))

	if fired != 1 {
		t.Fatalf("interleave fired %d times, want 1 -- the test is vacuous", fired)
	}
	if rr.Code != http.StatusNotFound {
		t.Errorf("expected 404 once the row is gone, got %d body=%s", rr.Code, rr.Body.String())
	}

	var count int64
	if err := shared.Model(&models.Cluster{}).Where("id = ?", cluster.ID).Count(&count).Error; err != nil {
		t.Fatalf("count: %v", err)
	}
	if count != 0 {
		t.Errorf("the deleted cluster was recreated by the detach")
	}
}
