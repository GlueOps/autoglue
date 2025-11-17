package handlers

import (
	"testing"

	"github.com/glueops/autoglue/internal/models"
	"github.com/glueops/autoglue/internal/testutil/pgtest"
	"github.com/google/uuid"
)

func TestValidStatus(t *testing.T) {
	// known-good statuses from servers.go
	valid := []string{"pending", "provisioning", "ready", "failed"}
	for _, s := range valid {
		if !validStatus(s) {
			t.Errorf("expected validStatus(%q) = true, got false", s)
		}
	}

	invalid := []string{"foobar", "unknown"}
	for _, s := range invalid {
		if validStatus(s) {
			t.Errorf("expected validStatus(%q) = false, got true", s)
		}
	}
}

func TestEnsureKeyBelongsToOrg_Success(t *testing.T) {
	db := pgtest.DB(t)

	org := models.Organization{Name: "servers-org"}
	if err := db.Create(&org).Error; err != nil {
		t.Fatalf("create org: %v", err)
	}

	key := createTestSshKey(t, db, org.ID, "org-key")

	if err := ensureKeyBelongsToOrg(org.ID, key.ID, db); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestEnsureKeyBelongsToOrg_WrongOrg(t *testing.T) {
	db := pgtest.DB(t)

	orgA := models.Organization{Name: "org-a"}
	orgB := models.Organization{Name: "org-b"}

	if err := db.Create(&orgA).Error; err != nil {
		t.Fatalf("create orgA: %v", err)
	}
	if err := db.Create(&orgB).Error; err != nil {
		t.Fatalf("create orgB: %v", err)
	}

	keyA := createTestSshKey(t, db, orgA.ID, "org-a-key")

	// ask for orgB with a key that belongs to orgA → should fail
	if err := ensureKeyBelongsToOrg(orgB.ID, keyA.ID, db); err == nil {
		t.Fatalf("expected error when ssh key belongs to a different org, got nil")
	}
}

func TestEnsureKeyBelongsToOrg_NotFound(t *testing.T) {
	db := pgtest.DB(t)

	org := models.Organization{Name: "org-nokey"}
	if err := db.Create(&org).Error; err != nil {
		t.Fatalf("create org: %v", err)
	}

	// random keyID that doesn't exist
	randomKeyID := uuid.New()

	if err := ensureKeyBelongsToOrg(org.ID, randomKeyID, db); err == nil {
		t.Fatalf("expected error when ssh key does not exist, got nil")
	}
}
