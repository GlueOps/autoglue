package handlers

import (
	"context"
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

func serverReq(orgID *uuid.UUID, serverID string) *http.Request {
	r := httptest.NewRequest(http.MethodPost, "/servers/"+serverID+"/reprovision", nil)
	ctx := r.Context()
	if orgID != nil {
		ctx = httpmiddleware.WithOrg(ctx, &models.Organization{ID: *orgID})
	}
	routeCtx := chi.NewRouteContext()
	routeCtx.URLParams.Add("id", serverID)
	ctx = context.WithValue(ctx, chi.RouteCtxKey, routeCtx)
	return r.WithContext(ctx)
}

func seedProvisionedServer(t *testing.T, db *gorm.DB, orgID uuid.UUID) models.Server {
	t.Helper()

	key := createTestSshKey(t, db, orgID, "reprov-"+uuid.NewString()[:8])
	pub := "203.0.113.10"
	srv := models.Server{
		OrganizationID:   orgID,
		Hostname:         "bastion",
		PublicIPAddress:  &pub,
		PrivateIPAddress: "10.0.0.9",
		SSHUser:          "deploy",
		SshKeyID:         key.ID,
		Role:             "bastion",
		Status:           "ready",
		SSHHostKey:       "AAAAC3NzaC1lZDI1NTE5stored",
		SSHHostKeyAlgo:   "ssh-ed25519",
	}
	if err := db.Create(&srv).Error; err != nil {
		t.Fatalf("create server: %v", err)
	}
	return srv
}

func reloadServer(t *testing.T, db *gorm.DB, id uuid.UUID) models.Server {
	t.Helper()
	var s models.Server
	if err := db.Where("id = ?", id).First(&s).Error; err != nil {
		t.Fatalf("reload server: %v", err)
	}
	return s
}

// Both effects or neither. Clearing the host key without queueing leaves the
// server trusting whatever answers next with nothing scheduled to reconnect;
// queueing without clearing leaves every attempt failing on mismatch forever.
func TestReprovisionQueuesAndClearsHostKeyTogether(t *testing.T) {
	db := pgtest.DB(t)

	org := models.Organization{Name: "reprov-org"}
	if err := db.Create(&org).Error; err != nil {
		t.Fatalf("create org: %v", err)
	}
	srv := seedProvisionedServer(t, db, org.ID)

	w := httptest.NewRecorder()
	ReprovisionServer(db)(w, serverReq(&org.ID, srv.ID.String()))

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, body = %s", w.Code, w.Body.String())
	}

	got := reloadServer(t, db, srv.ID)
	if got.Status != "pending" {
		t.Errorf("status = %q, want pending — nothing will re-bootstrap it", got.Status)
	}
	if got.SSHHostKey != "" || got.SSHHostKeyAlgo != "" {
		t.Errorf("host key not cleared: %q/%q — every reconnect will fail on mismatch",
			got.SSHHostKeyAlgo, got.SSHHostKey)
	}
}

// Org scoping: reprovisioning drops a trust anchor, so reaching across orgs
// would let one tenant disable another's host key verification.
func TestReprovisionIsOrgScoped(t *testing.T) {
	db := pgtest.DB(t)

	orgA := models.Organization{Name: "reprov-a"}
	orgB := models.Organization{Name: "reprov-b"}
	if err := db.Create(&orgA).Error; err != nil {
		t.Fatalf("create orgA: %v", err)
	}
	if err := db.Create(&orgB).Error; err != nil {
		t.Fatalf("create orgB: %v", err)
	}
	srv := seedProvisionedServer(t, db, orgA.ID)

	w := httptest.NewRecorder()
	ReprovisionServer(db)(w, serverReq(&orgB.ID, srv.ID.String()))

	if w.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want 404 for another org's server", w.Code)
	}

	got := reloadServer(t, db, srv.ID)
	if got.SSHHostKey == "" {
		t.Fatal("a cross-org request cleared the host key")
	}
	if got.Status != "ready" {
		t.Fatalf("a cross-org request changed status to %q", got.Status)
	}
}

func TestReprovisionRequiresOrg(t *testing.T) {
	db := pgtest.DB(t)
	w := httptest.NewRecorder()
	ReprovisionServer(db)(w, serverReq(nil, uuid.NewString()))
	if w.Code != http.StatusForbidden {
		t.Fatalf("status = %d, want 403 without an org", w.Code)
	}
}

// A worker has no stored host key and is not claimed by the sweep, so both
// halves would be no-ops. Refusing beats a 200 that reports a requeue which
// never happened.
func TestReprovisionRefusesNonBastion(t *testing.T) {
	db := pgtest.DB(t)
	org := models.Organization{Name: "reprov-worker"}
	if err := db.Create(&org).Error; err != nil {
		t.Fatalf("create org: %v", err)
	}
	srv := seedProvisionedServer(t, db, org.ID)
	if err := db.Model(&models.Server{}).Where("id = ?", srv.ID).
		Update("role", "worker").Error; err != nil {
		t.Fatalf("seed worker: %v", err)
	}

	w := httptest.NewRecorder()
	ReprovisionServer(db)(w, serverReq(&org.ID, srv.ID.String()))
	if w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400 for a non-bastion", w.Code)
	}
	if got := reloadServer(t, db, srv.ID); got.Status != "ready" {
		t.Fatalf("a refused request still changed status to %q", got.Status)
	}
}

func TestReprovisionRejectsBadID(t *testing.T) {
	db := pgtest.DB(t)
	org := models.Organization{Name: "reprov-badid"}
	if err := db.Create(&org).Error; err != nil {
		t.Fatalf("create org: %v", err)
	}
	w := httptest.NewRecorder()
	ReprovisionServer(db)(w, serverReq(&org.ID, "not-a-uuid"))
	if w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400", w.Code)
	}
}

// Recovery has to work on a server already stuck mid-provision — that is the
// case someone reaches for this in.
func TestReprovisionWorksFromProvisioning(t *testing.T) {
	db := pgtest.DB(t)
	org := models.Organization{Name: "reprov-stuck"}
	if err := db.Create(&org).Error; err != nil {
		t.Fatalf("create org: %v", err)
	}
	srv := seedProvisionedServer(t, db, org.ID)
	if err := db.Model(&models.Server{}).Where("id = ?", srv.ID).
		Update("status", "provisioning").Error; err != nil {
		t.Fatalf("seed provisioning: %v", err)
	}

	w := httptest.NewRecorder()
	ReprovisionServer(db)(w, serverReq(&org.ID, srv.ID.String()))
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, body = %s", w.Code, w.Body.String())
	}
	if got := reloadServer(t, db, srv.ID); got.Status != "pending" {
		t.Fatalf("a stuck server was not requeued: %q", got.Status)
	}
}
