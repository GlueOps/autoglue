package handlers

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/glueops/autoglue/internal/api/httpmiddleware"
	"github.com/glueops/autoglue/internal/models"
	"github.com/glueops/autoglue/internal/testutil/pgtest"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func runLogsReq(orgID *uuid.UUID, clusterID, runID, query string) *http.Request {
	r := httptest.NewRequest(http.MethodGet,
		fmt.Sprintf("/clusters/%s/runs/%s/logs%s", clusterID, runID, query), nil)

	ctx := r.Context()
	if orgID != nil {
		ctx = httpmiddleware.WithOrg(ctx, &models.Organization{ID: *orgID})
	}
	routeCtx := chi.NewRouteContext()
	routeCtx.URLParams.Add("clusterID", clusterID)
	routeCtx.URLParams.Add("runID", runID)
	ctx = context.WithValue(ctx, chi.RouteCtxKey, routeCtx)
	return r.WithContext(ctx)
}

// seedRun creates an org, a cluster and a run on it.
func seedRun(t *testing.T, db *gorm.DB, orgID uuid.UUID, status string) (clusterID, runID uuid.UUID) {
	t.Helper()

	org := models.Organization{ID: orgID, Name: "org-" + orgID.String()[:8]}
	if err := db.Create(&org).Error; err != nil {
		t.Fatalf("seed org: %v", err)
	}

	cluster := models.Cluster{OrganizationID: orgID, Name: "c-" + orgID.String()[:8]}
	if err := db.Create(&cluster).Error; err != nil {
		t.Fatalf("seed cluster: %v", err)
	}

	run := models.ClusterRun{
		OrganizationID: orgID,
		ClusterID:      cluster.ID,
		Action:         "bootstrap",
		Status:         status,
	}
	if err := db.Create(&run).Error; err != nil {
		t.Fatalf("seed run: %v", err)
	}
	return cluster.ID, run.ID
}

func TestGetClusterRunLogs_PagesByCursor(t *testing.T) {
	db := pgtest.DB(t)
	orgID := uuid.New()
	clusterID, runID := seedRun(t, db, orgID, models.ClusterRunStatusRunning)
	seedJobLogs(t, db, orgID, runID, models.JobLogSubjectClusterRun, 5)

	rr := httptest.NewRecorder()
	GetClusterRunLogs(db)(rr, runLogsReq(&orgID, clusterID.String(), runID.String(), "?limit=2"))
	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200 (body=%s)", rr.Code, rr.Body.String())
	}
	first := decodePage(t, rr)
	if len(first.Items) != 2 {
		t.Fatalf("first page = %d items, want 2", len(first.Items))
	}

	rr = httptest.NewRecorder()
	GetClusterRunLogs(db)(rr, runLogsReq(&orgID, clusterID.String(), runID.String(),
		fmt.Sprintf("?limit=10&after=%d", first.NextCursor)))
	second := decodePage(t, rr)
	if len(second.Items) != 3 {
		t.Fatalf("second page = %d items, want the remaining 3", len(second.Items))
	}
	if second.Items[0].Chunk != "chunk-2\n" {
		t.Errorf("second page starts at %q, want chunk-2", second.Items[0].Chunk)
	}
}

func TestGetClusterRunLogs_DoneTracksRunStatus(t *testing.T) {
	db := pgtest.DB(t)

	// A running job must never report done, or a client stops tailing early.
	orgA := uuid.New()
	clusterA, runA := seedRun(t, db, orgA, models.ClusterRunStatusRunning)
	seedJobLogs(t, db, orgA, runA, models.JobLogSubjectClusterRun, 2)

	rr := httptest.NewRecorder()
	GetClusterRunLogs(db)(rr, runLogsReq(&orgA, clusterA.String(), runA.String(), "?limit=10"))
	if page := decodePage(t, rr); page.Done {
		t.Error("done = true for a running job")
	}

	// A finished job with the reader caught up may stop.
	orgB := uuid.New()
	clusterB, runB := seedRun(t, db, orgB, models.ClusterRunStatusSuccess)
	seedJobLogs(t, db, orgB, runB, models.JobLogSubjectClusterRun, 2)

	rr = httptest.NewRecorder()
	GetClusterRunLogs(db)(rr, runLogsReq(&orgB, clusterB.String(), runB.String(), "?limit=10"))
	if page := decodePage(t, rr); !page.Done {
		t.Error("done = false for a finished run with the reader caught up")
	}
}

func TestGetClusterRunLogs_RejectsRunFromAnotherCluster(t *testing.T) {
	db := pgtest.DB(t)
	orgID := uuid.New()

	clusterA, runA := seedRun(t, db, orgID, models.ClusterRunStatusRunning)
	seedJobLogs(t, db, orgID, runA, models.JobLogSubjectClusterRun, 2)

	// A second cluster in the same org. Asking for cluster B's URL with cluster
	// A's run must not work: the org filter alone would happily allow it.
	otherCluster := models.Cluster{OrganizationID: orgID, Name: "other"}
	if err := db.Create(&otherCluster).Error; err != nil {
		t.Fatalf("seed second cluster: %v", err)
	}

	rr := httptest.NewRecorder()
	GetClusterRunLogs(db)(rr, runLogsReq(&orgID, otherCluster.ID.String(), runA.String(), ""))
	if rr.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want 404 for a run on a different cluster", rr.Code)
	}

	// Sanity: the same run under its own cluster is fine.
	rr = httptest.NewRecorder()
	GetClusterRunLogs(db)(rr, runLogsReq(&orgID, clusterA.String(), runA.String(), ""))
	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200 under the correct cluster", rr.Code)
	}
}

func TestGetClusterRunLogs_DoesNotLeakAcrossOrgs(t *testing.T) {
	db := pgtest.DB(t)

	orgA := uuid.New()
	clusterA, runA := seedRun(t, db, orgA, models.ClusterRunStatusRunning)
	seedJobLogs(t, db, orgA, runA, models.JobLogSubjectClusterRun, 3)

	orgB := uuid.New()
	_, _ = seedRun(t, db, orgB, models.ClusterRunStatusRunning)

	rr := httptest.NewRecorder()
	GetClusterRunLogs(db)(rr, runLogsReq(&orgB, clusterA.String(), runA.String(), ""))
	if rr.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want 404 for another org's run", rr.Code)
	}
}

func TestGetClusterRunLogs_RequiresOrg(t *testing.T) {
	db := pgtest.DB(t)
	rr := httptest.NewRecorder()
	GetClusterRunLogs(db)(rr, runLogsReq(nil, uuid.NewString(), uuid.NewString(), ""))
	if rr.Code != http.StatusForbidden {
		t.Fatalf("status = %d, want 403 without an org", rr.Code)
	}
}
