package handlers

import (
	"testing"

	"github.com/glueops/autoglue/internal/handlers/dto"
	"github.com/glueops/autoglue/internal/models"
	"github.com/glueops/autoglue/internal/testutil/pgtest"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// sshKeyListColumns embeds a correlated subquery in a multi-line Select string.
// GORM quotes each argument when Select is given several column names, so the
// projection only survives as one raw string — and if that ever regresses the
// failure is a broken query, not a compile error. Hence a real database.
func TestSshKeyListColumnsCountsServers(t *testing.T) {
	db := pgtest.DB(t)

	org := models.Organization{Name: "ssh-count-org"}
	if err := db.Create(&org).Error; err != nil {
		t.Fatalf("create org: %v", err)
	}

	shared := createTestSshKey(t, db, org.ID, "shared-key")
	lonely := createTestSshKey(t, db, org.ID, "unattached-key")

	// Two servers on one key: the case that makes an unconditional "delete the
	// key when its server goes" wrong, and the case a LEFT JOIN without GROUP BY
	// would report as two rows.
	createTestServer(t, db, org.ID, shared.ID, "host-a")
	createTestServer(t, db, org.ID, shared.ID, "host-b")

	var out []dto.SshResponse
	if err := db.Model(&models.SshKey{}).
		Where("organization_id = ?", org.ID).
		Select(sshKeyListColumns).
		Order("created_at DESC").
		Scan(&out).Error; err != nil {
		t.Fatalf("list: %v", err)
	}

	counts := map[uuid.UUID]*int{}
	for _, k := range out {
		counts[k.ID] = k.ServerCount
	}

	if len(out) != 2 {
		t.Fatalf("got %d rows, want 2 — a duplicated row means the count is joining rather than subquerying", len(out))
	}
	if got := counts[shared.ID]; got == nil || *got != 2 {
		t.Errorf("shared key server_count = %v, want 2", intOrNil(got))
	}
	if got := counts[lonely.ID]; got == nil || *got != 0 {
		t.Errorf("unattached key server_count = %v, want 0", intOrNil(got))
	}
}

func TestSshKeyUnattachedFilter(t *testing.T) {
	db := pgtest.DB(t)

	org := models.Organization{Name: "ssh-unattached-org"}
	if err := db.Create(&org).Error; err != nil {
		t.Fatalf("create org: %v", err)
	}

	attached := createTestSshKey(t, db, org.ID, "attached-key")
	orphan := createTestSshKey(t, db, org.ID, "orphan-key")
	createTestServer(t, db, org.ID, attached.ID, "filter-host")

	var out []dto.SshResponse
	if err := db.Model(&models.SshKey{}).
		Where("organization_id = ?", org.ID).
		Select(sshKeyListColumns).
		Where(sshKeyUnattached).
		Scan(&out).Error; err != nil {
		t.Fatalf("list unattached: %v", err)
	}

	if len(out) != 1 {
		t.Fatalf("got %d unattached keys, want 1", len(out))
	}
	if out[0].ID != orphan.ID {
		t.Errorf("returned key %s, want the orphan %s", out[0].ID, orphan.ID)
	}
	if out[0].ServerCount == nil || *out[0].ServerCount != 0 {
		t.Errorf("orphan server_count = %v, want 0", intOrNil(out[0].ServerCount))
	}
}

// The projection must never leak key material: it is the response body for an
// endpoint any org member can call, and the private columns live on the same row.
func TestSshKeyListColumnsOmitPrivateMaterial(t *testing.T) {
	db := pgtest.DB(t)

	org := models.Organization{Name: "ssh-secrets-org"}
	if err := db.Create(&org).Error; err != nil {
		t.Fatalf("create org: %v", err)
	}
	createTestSshKey(t, db, org.ID, "secret-key")

	var out []dto.SshResponse
	if err := db.Model(&models.SshKey{}).
		Where("organization_id = ?", org.ID).
		Select(sshKeyListColumns).
		Scan(&out).Error; err != nil {
		t.Fatalf("list: %v", err)
	}

	if len(out) != 1 {
		t.Fatalf("got %d rows, want 1", len(out))
	}
	if out[0].EncryptedPrivateKey != "" || out[0].PrivateIV != "" || out[0].PrivateTag != "" {
		t.Error("list projection selected encrypted key material")
	}
}

func createTestServer(t *testing.T, db *gorm.DB, orgID, keyID uuid.UUID, hostname string) models.Server {
	t.Helper()

	srv := models.Server{
		OrganizationID:   orgID,
		Hostname:         hostname,
		PrivateIPAddress: "10.0.0.1",
		SSHUser:          "deploy",
		SshKeyID:         keyID,
		Role:             "worker",
	}
	if err := db.Create(&srv).Error; err != nil {
		t.Fatalf("create server %s: %v", hostname, err)
	}
	return srv
}

func intOrNil(p *int) any {
	if p == nil {
		return nil
	}
	return *p
}
