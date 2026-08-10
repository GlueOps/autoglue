package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/glueops/autoglue/internal/api/httpmiddleware"
	"github.com/glueops/autoglue/internal/common"
	"github.com/glueops/autoglue/internal/handlers/dto"
	"github.com/glueops/autoglue/internal/models"
	"github.com/glueops/autoglue/internal/testutil/pgtest"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func seedJobLogs(t *testing.T, db *gorm.DB, orgID, subjectID uuid.UUID, subjectType string, n int) {
	t.Helper()
	for i := 0; i < n; i++ {
		row := models.JobLog{
			JobID:          1,
			OrganizationID: orgID,
			SubjectType:    subjectType,
			SubjectID:      subjectID,
			Stream:         models.JobLogStreamStdout,
			Chunk:          fmt.Sprintf("chunk-%d\n", i),
		}
		if err := db.Create(&row).Error; err != nil {
			t.Fatalf("seed job log: %v", err)
		}
	}
}

func serverLogsReq(orgID *uuid.UUID, serverID string, query string) *http.Request {
	r := httptest.NewRequest(http.MethodGet, "/servers/"+serverID+"/logs"+query, nil)

	ctx := r.Context()
	if orgID != nil {
		ctx = httpmiddleware.WithOrg(ctx, &models.Organization{ID: *orgID})
	}
	routeCtx := chi.NewRouteContext()
	routeCtx.URLParams.Add("id", serverID)
	ctx = context.WithValue(ctx, chi.RouteCtxKey, routeCtx)
	return r.WithContext(ctx)
}

func decodePage(t *testing.T, rr *httptest.ResponseRecorder) dto.JobLogPage {
	t.Helper()
	var page dto.JobLogPage
	if err := json.Unmarshal(rr.Body.Bytes(), &page); err != nil {
		t.Fatalf("decode page: %v (body=%s)", err, rr.Body.String())
	}
	return page
}

// seedServer creates the minimum server row the handler authorizes against.
func seedServer(t *testing.T, db *gorm.DB, orgID uuid.UUID, status string) uuid.UUID {
	t.Helper()

	org := models.Organization{ID: orgID, Name: "org-" + orgID.String()[:8]}
	if err := db.Create(&org).Error; err != nil {
		t.Fatalf("seed org: %v", err)
	}

	key := models.SshKey{
		AuditFields:         common.AuditFields{OrganizationID: orgID},
		Name:                "k",
		PublicKey:           "ssh-ed25519 AAAA",
		EncryptedPrivateKey: "x",
		PrivateIV:           "x",
		PrivateTag:          "x",
		Fingerprint:         uuid.NewString(),
	}
	if err := db.Create(&key).Error; err != nil {
		t.Fatalf("seed ssh key: %v", err)
	}

	id := uuid.New()
	if err := db.Exec(
		`INSERT INTO servers (id, organization_id, ssh_key_id, hostname, private_ip_address, ssh_user, role, status, created_at, updated_at)
		 VALUES (?, ?, ?, 'host', '10.0.0.1', 'ubuntu', 'bastion', ?, now(), now())`,
		id, orgID, key.ID, status,
	).Error; err != nil {
		t.Fatalf("seed server: %v", err)
	}
	return id
}

func TestGetServerLogs_PagesByCursor(t *testing.T) {
	db := pgtest.DB(t)
	orgID := uuid.New()
	serverID := seedServer(t, db, orgID, "provisioning")
	seedJobLogs(t, db, orgID, serverID, models.JobLogSubjectServer, 5)

	// First page, capped at 2.
	rr := httptest.NewRecorder()
	GetServerLogs(db)(rr, serverLogsReq(&orgID, serverID.String(), "?limit=2"))
	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200 (body=%s)", rr.Code, rr.Body.String())
	}
	first := decodePage(t, rr)
	if len(first.Items) != 2 {
		t.Fatalf("first page = %d items, want 2", len(first.Items))
	}
	if first.Done {
		t.Error("done should be withheld while the reader is still behind")
	}

	// Second page continues from the cursor rather than repeating.
	rr = httptest.NewRecorder()
	GetServerLogs(db)(rr, serverLogsReq(&orgID, serverID.String(),
		fmt.Sprintf("?limit=10&after=%d", first.NextCursor)))
	second := decodePage(t, rr)
	if len(second.Items) != 3 {
		t.Fatalf("second page = %d items, want the remaining 3", len(second.Items))
	}
	if second.Items[0].Chunk != "chunk-2\n" {
		t.Errorf("second page starts at %q, want chunk-2 (cursor did not advance)", second.Items[0].Chunk)
	}
}

func TestGetServerLogs_DoneOnlyWhenTerminalAndCaughtUp(t *testing.T) {
	db := pgtest.DB(t)
	orgID := uuid.New()
	serverID := seedServer(t, db, orgID, "ready")
	seedJobLogs(t, db, orgID, serverID, models.JobLogSubjectServer, 3)

	// Terminal status, but the reader has a full page: done must stay false or
	// a client that stops here loses the tail.
	rr := httptest.NewRecorder()
	GetServerLogs(db)(rr, serverLogsReq(&orgID, serverID.String(), "?limit=3"))
	if page := decodePage(t, rr); page.Done {
		t.Error("done = true while the page was full; client would miss later chunks")
	}

	// Caught up on a terminal server: now it is safe to stop.
	rr = httptest.NewRecorder()
	GetServerLogs(db)(rr, serverLogsReq(&orgID, serverID.String(), "?limit=10"))
	if page := decodePage(t, rr); !page.Done {
		t.Error("done = false for a ready server with the reader caught up")
	}
}

func TestGetServerLogs_DoesNotLeakAcrossOrgs(t *testing.T) {
	db := pgtest.DB(t)

	orgA := uuid.New()
	serverA := seedServer(t, db, orgA, "provisioning")
	seedJobLogs(t, db, orgA, serverA, models.JobLogSubjectServer, 3)

	orgB := uuid.New()
	_ = seedServer(t, db, orgB, "provisioning")

	// Org B asking for org A's server must 404, not return its output.
	rr := httptest.NewRecorder()
	GetServerLogs(db)(rr, serverLogsReq(&orgB, serverA.String(), ""))
	if rr.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want 404 for another org's server", rr.Code)
	}
}

func TestGetServerLogs_RequiresOrg(t *testing.T) {
	db := pgtest.DB(t)
	rr := httptest.NewRecorder()
	GetServerLogs(db)(rr, serverLogsReq(nil, uuid.NewString(), ""))
	if rr.Code != http.StatusForbidden {
		t.Fatalf("status = %d, want 403 without an org", rr.Code)
	}
}
