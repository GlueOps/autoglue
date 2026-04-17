package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/glueops/autoglue/internal/api/httpmiddleware"
	"github.com/glueops/autoglue/internal/models"
	"github.com/glueops/autoglue/internal/testutil/pgtest"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func TestCreateClusterMetadata_OrgScoping(t *testing.T) {
	db := pgtest.DB(t)
	migrateClusterMetadata(t, db)

	orgA := createTestOrg(t, db, "metadata-org-a")
	orgB := createTestOrg(t, db, "metadata-org-b")
	clusterInOrgB := createTestCluster(t, db, orgB.ID, "cluster-b")

	req := httptest.NewRequest(http.MethodPost, "/clusters/"+clusterInOrgB.ID.String()+"/metadata", strings.NewReader(`{"key":"network.service_cidr","value":"10.96.0.0/12"}`))
	req.Header.Set("Content-Type", "application/json")
	req = withOrgAndClusterID(req, orgA.ID, clusterInOrgB.ID)
	rr := httptest.NewRecorder()

	CreateClusterMetadata(db).ServeHTTP(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Fatalf("expected 404 for cross-org cluster access, got %d body=%s", rr.Code, rr.Body.String())
	}
}

func TestCreateClusterMetadata_NormalizesKeyAndTrimsValue(t *testing.T) {
	db := pgtest.DB(t)
	migrateClusterMetadata(t, db)

	org := createTestOrg(t, db, "metadata-normalize-org")
	cluster := createTestCluster(t, db, org.ID, "cluster-normalize")

	req := httptest.NewRequest(http.MethodPost, "/clusters/"+cluster.ID.String()+"/metadata", strings.NewReader(`{"key":"  Network.Service_CIDR  ","value":"  My Value  "}`))
	req.Header.Set("Content-Type", "application/json")
	req = withOrgAndClusterID(req, org.ID, cluster.ID)
	rr := httptest.NewRecorder()

	CreateClusterMetadata(db).ServeHTTP(rr, req)

	if rr.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d body=%s", rr.Code, rr.Body.String())
	}

	var out map[string]any
	if err := json.Unmarshal(rr.Body.Bytes(), &out); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if got := out["key"]; got != "network.service_cidr" {
		t.Fatalf("expected normalized lowercase key, got %v", got)
	}
	if got := out["value"]; got != "My Value" {
		t.Fatalf("expected trimmed value preserving case, got %v", got)
	}
}

func TestCreateClusterMetadata_RequiredFields(t *testing.T) {
	db := pgtest.DB(t)
	migrateClusterMetadata(t, db)

	org := createTestOrg(t, db, "metadata-required-org")
	cluster := createTestCluster(t, db, org.ID, "cluster-required")

	cases := []struct {
		name string
		body string
	}{
		{name: "missing key", body: `{"key":"   ","value":"ok"}`},
		{name: "missing value", body: `{"key":"network.calico_cidr","value":"   "}`},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/clusters/"+cluster.ID.String()+"/metadata", strings.NewReader(tc.body))
			req.Header.Set("Content-Type", "application/json")
			req = withOrgAndClusterID(req, org.ID, cluster.ID)
			rr := httptest.NewRecorder()

			CreateClusterMetadata(db).ServeHTTP(rr, req)
			if rr.Code != http.StatusBadRequest {
				t.Fatalf("expected 400, got %d body=%s", rr.Code, rr.Body.String())
			}
		})
	}
}

func TestCreateClusterMetadata_DuplicateKeyConflict(t *testing.T) {
	db := pgtest.DB(t)
	migrateClusterMetadata(t, db)

	org := createTestOrg(t, db, "metadata-duplicate-org")
	cluster := createTestCluster(t, db, org.ID, "cluster-duplicate")

	first := httptest.NewRequest(http.MethodPost, "/clusters/"+cluster.ID.String()+"/metadata", strings.NewReader(`{"key":"network.service_cidr","value":"10.96.0.0/12"}`))
	first.Header.Set("Content-Type", "application/json")
	first = withOrgAndClusterID(first, org.ID, cluster.ID)
	firstRR := httptest.NewRecorder()
	CreateClusterMetadata(db).ServeHTTP(firstRR, first)
	if firstRR.Code != http.StatusCreated {
		t.Fatalf("expected first create to succeed, got %d body=%s", firstRR.Code, firstRR.Body.String())
	}

	second := httptest.NewRequest(http.MethodPost, "/clusters/"+cluster.ID.String()+"/metadata", strings.NewReader(`{"key":"NETWORK.SERVICE_CIDR","value":"10.97.0.0/12"}`))
	second.Header.Set("Content-Type", "application/json")
	second = withOrgAndClusterID(second, org.ID, cluster.ID)
	secondRR := httptest.NewRecorder()
	CreateClusterMetadata(db).ServeHTTP(secondRR, second)

	if secondRR.Code != http.StatusConflict {
		t.Fatalf("expected 409 on duplicate key, got %d body=%s", secondRR.Code, secondRR.Body.String())
	}
}

func migrateClusterMetadata(t *testing.T, db *gorm.DB) {
	t.Helper()
	if err := db.AutoMigrate(&models.ClusterMetadata{}); err != nil {
		t.Fatalf("migrate cluster metadata: %v", err)
	}
}

func createTestOrg(t *testing.T, db *gorm.DB, namePrefix string) models.Organization {
	t.Helper()
	org := models.Organization{Name: namePrefix + "-" + uuid.NewString()}
	if err := db.Create(&org).Error; err != nil {
		t.Fatalf("create org: %v", err)
	}
	return org
}

func createTestCluster(t *testing.T, db *gorm.DB, orgID uuid.UUID, namePrefix string) models.Cluster {
	t.Helper()
	cluster := models.Cluster{
		OrganizationID: orgID,
		Name:           namePrefix + "-" + uuid.NewString(),
	}
	if err := db.Create(&cluster).Error; err != nil {
		t.Fatalf("create cluster: %v", err)
	}
	return cluster
}

func withOrgAndClusterID(r *http.Request, orgID, clusterID uuid.UUID) *http.Request {
	ctx := httpmiddleware.WithOrg(r.Context(), &models.Organization{ID: orgID})
	routeCtx := chi.NewRouteContext()
	routeCtx.URLParams.Add("clusterID", clusterID.String())
	ctx = context.WithValue(ctx, chi.RouteCtxKey, routeCtx)
	return r.WithContext(ctx)
}
