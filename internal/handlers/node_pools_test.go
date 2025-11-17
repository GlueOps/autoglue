package handlers

import (
	"os"
	"testing"

	"github.com/glueops/autoglue/internal/common"
	"github.com/glueops/autoglue/internal/models"
	"github.com/glueops/autoglue/internal/testutil/pgtest"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func TestMain(m *testing.M) {
	code := m.Run()
	pgtest.Stop()
	os.Exit(code)
}

func TestParseUUIDs_Success(t *testing.T) {
	u1 := uuid.New()
	u2 := uuid.New()

	got, err := parseUUIDs([]string{u1.String(), u2.String()})
	if err != nil {
		t.Fatalf("parseUUIDs returned error: %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("expected 2 UUIDs, got %d", len(got))
	}
	if got[0] != u1 || got[1] != u2 {
		t.Fatalf("unexpected UUIDs: got=%v", got)
	}
}

func TestParseUUIDs_Invalid(t *testing.T) {
	_, err := parseUUIDs([]string{"not-a-uuid"})
	if err == nil {
		t.Fatalf("expected error for invalid UUID, got nil")
	}
}

// --- ensureServersBelongToOrg ---

func TestEnsureServersBelongToOrg_AllBelong(t *testing.T) {
	db := pgtest.DB(t)

	org := models.Organization{Name: "org-a"}
	if err := db.Create(&org).Error; err != nil {
		t.Fatalf("create org: %v", err)
	}

	sshKey := createTestSshKey(t, db, org.ID, "org-a-key")

	s1 := models.Server{
		OrganizationID: org.ID,
		Hostname:       "srv-1",
		SSHUser:        "ubuntu",
		SshKeyID:       sshKey.ID,
		Role:           "worker",
		Status:         "pending",
	}
	s2 := models.Server{
		OrganizationID: org.ID,
		Hostname:       "srv-2",
		SSHUser:        "ubuntu",
		SshKeyID:       sshKey.ID,
		Role:           "worker",
		Status:         "pending",
	}

	if err := db.Create(&s1).Error; err != nil {
		t.Fatalf("create server 1: %v", err)
	}
	if err := db.Create(&s2).Error; err != nil {
		t.Fatalf("create server 2: %v", err)
	}

	ids := []uuid.UUID{s1.ID, s2.ID}
	if err := ensureServersBelongToOrg(org.ID, ids, db); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestEnsureServersBelongToOrg_ForeignOrgFails(t *testing.T) {
	db := pgtest.DB(t)

	orgA := models.Organization{Name: "org-a"}
	orgB := models.Organization{Name: "org-b"}

	if err := db.Create(&orgA).Error; err != nil {
		t.Fatalf("create orgA: %v", err)
	}
	if err := db.Create(&orgB).Error; err != nil {
		t.Fatalf("create orgB: %v", err)
	}

	sshKeyA := createTestSshKey(t, db, orgA.ID, "org-a-key")
	sshKeyB := createTestSshKey(t, db, orgB.ID, "org-b-key")

	s1 := models.Server{
		OrganizationID: orgA.ID,
		Hostname:       "srv-a-1",
		SSHUser:        "ubuntu",
		SshKeyID:       sshKeyA.ID,
		Role:           "worker",
		Status:         "pending",
	}
	s2 := models.Server{
		OrganizationID: orgB.ID,
		Hostname:       "srv-b-1",
		SSHUser:        "ubuntu",
		SshKeyID:       sshKeyB.ID,
		Role:           "worker",
		Status:         "pending",
	}

	if err := db.Create(&s1).Error; err != nil {
		t.Fatalf("create server s1: %v", err)
	}
	if err := db.Create(&s2).Error; err != nil {
		t.Fatalf("create server s2: %v", err)
	}

	ids := []uuid.UUID{s1.ID, s2.ID}
	if err := ensureServersBelongToOrg(orgA.ID, ids, db); err == nil {
		t.Fatalf("expected error when one server belongs to a different org")
	}
}

// --- ensureTaintsBelongToOrg ---

func TestEnsureTaintsBelongToOrg_AllBelong(t *testing.T) {
	db := pgtest.DB(t)

	org := models.Organization{Name: "org-taints"}
	if err := db.Create(&org).Error; err != nil {
		t.Fatalf("create org: %v", err)
	}

	t1 := models.Taint{
		OrganizationID: org.ID,
		Key:            "key1",
		Value:          "val1",
		Effect:         "NoSchedule",
	}
	t2 := models.Taint{
		OrganizationID: org.ID,
		Key:            "key2",
		Value:          "val2",
		Effect:         "PreferNoSchedule",
	}

	if err := db.Create(&t1).Error; err != nil {
		t.Fatalf("create taint 1: %v", err)
	}
	if err := db.Create(&t2).Error; err != nil {
		t.Fatalf("create taint 2: %v", err)
	}

	ids := []uuid.UUID{t1.ID, t2.ID}
	if err := ensureTaintsBelongToOrg(org.ID, ids, db); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestEnsureTaintsBelongToOrg_ForeignOrgFails(t *testing.T) {
	db := pgtest.DB(t)

	orgA := models.Organization{Name: "org-a"}
	orgB := models.Organization{Name: "org-b"}
	if err := db.Create(&orgA).Error; err != nil {
		t.Fatalf("create orgA: %v", err)
	}
	if err := db.Create(&orgB).Error; err != nil {
		t.Fatalf("create orgB: %v", err)
	}

	t1 := models.Taint{
		OrganizationID: orgA.ID,
		Key:            "key1",
		Value:          "val1",
		Effect:         "NoSchedule",
	}
	t2 := models.Taint{
		OrganizationID: orgB.ID,
		Key:            "key2",
		Value:          "val2",
		Effect:         "NoSchedule",
	}

	if err := db.Create(&t1).Error; err != nil {
		t.Fatalf("create taint 1: %v", err)
	}
	if err := db.Create(&t2).Error; err != nil {
		t.Fatalf("create taint 2: %v", err)
	}

	ids := []uuid.UUID{t1.ID, t2.ID}
	if err := ensureTaintsBelongToOrg(orgA.ID, ids, db); err == nil {
		t.Fatalf("expected error when a taint belongs to another org")
	}
}

// --- ensureLabelsBelongToOrg ---

func TestEnsureLabelsBelongToOrg_AllBelong(t *testing.T) {
	db := pgtest.DB(t)

	org := models.Organization{Name: "org-labels"}
	if err := db.Create(&org).Error; err != nil {
		t.Fatalf("create org: %v", err)
	}

	l1 := models.Label{
		AuditFields: common.AuditFields{
			OrganizationID: org.ID,
		},
		Key:   "env",
		Value: "dev",
	}
	l2 := models.Label{
		AuditFields: common.AuditFields{
			OrganizationID: org.ID,
		},
		Key:   "env",
		Value: "prod",
	}

	if err := db.Create(&l1).Error; err != nil {
		t.Fatalf("create label 1: %v", err)
	}
	if err := db.Create(&l2).Error; err != nil {
		t.Fatalf("create label 2: %v", err)
	}

	ids := []uuid.UUID{l1.ID, l2.ID}
	if err := ensureLabelsBelongToOrg(org.ID, ids, db); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestEnsureLabelsBelongToOrg_ForeignOrgFails(t *testing.T) {
	db := pgtest.DB(t)

	orgA := models.Organization{Name: "org-a"}
	orgB := models.Organization{Name: "org-b"}
	if err := db.Create(&orgA).Error; err != nil {
		t.Fatalf("create orgA: %v", err)
	}
	if err := db.Create(&orgB).Error; err != nil {
		t.Fatalf("create orgB: %v", err)
	}

	l1 := models.Label{
		AuditFields: common.AuditFields{
			OrganizationID: orgA.ID,
		},
		Key:   "env",
		Value: "dev",
	}
	l2 := models.Label{
		AuditFields: common.AuditFields{
			OrganizationID: orgB.ID,
		},
		Key:   "env",
		Value: "prod",
	}

	if err := db.Create(&l1).Error; err != nil {
		t.Fatalf("create label 1: %v", err)
	}
	if err := db.Create(&l2).Error; err != nil {
		t.Fatalf("create label 2: %v", err)
	}

	ids := []uuid.UUID{l1.ID, l2.ID}
	if err := ensureLabelsBelongToOrg(orgA.ID, ids, db); err == nil {
		t.Fatalf("expected error when a label belongs to another org")
	}
}

// --- ensureAnnotaionsBelongToOrg (typo in original name is preserved) ---

func TestEnsureAnnotationsBelongToOrg_AllBelong(t *testing.T) {
	db := pgtest.DB(t)

	org := models.Organization{Name: "org-annotations"}
	if err := db.Create(&org).Error; err != nil {
		t.Fatalf("create org: %v", err)
	}

	a1 := models.Annotation{
		AuditFields: common.AuditFields{
			OrganizationID: org.ID,
		},
		Key:   "team",
		Value: "core",
	}
	a2 := models.Annotation{
		AuditFields: common.AuditFields{
			OrganizationID: org.ID,
		},
		Key:   "team",
		Value: "platform",
	}

	if err := db.Create(&a1).Error; err != nil {
		t.Fatalf("create annotation 1: %v", err)
	}
	if err := db.Create(&a2).Error; err != nil {
		t.Fatalf("create annotation 2: %v", err)
	}

	ids := []uuid.UUID{a1.ID, a2.ID}
	if err := ensureAnnotaionsBelongToOrg(org.ID, ids, db); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestEnsureAnnotationsBelongToOrg_ForeignOrgFails(t *testing.T) {
	db := pgtest.DB(t)

	orgA := models.Organization{Name: "org-a"}
	orgB := models.Organization{Name: "org-b"}
	if err := db.Create(&orgA).Error; err != nil {
		t.Fatalf("create orgA: %v", err)
	}
	if err := db.Create(&orgB).Error; err != nil {
		t.Fatalf("create orgB: %v", err)
	}

	a1 := models.Annotation{
		AuditFields: common.AuditFields{
			OrganizationID: orgA.ID,
		},
		Key:   "team",
		Value: "core",
	}
	a2 := models.Annotation{
		AuditFields: common.AuditFields{
			OrganizationID: orgB.ID,
		},
		Key:   "team",
		Value: "platform",
	}

	if err := db.Create(&a1).Error; err != nil {
		t.Fatalf("create annotation 1: %v", err)
	}
	if err := db.Create(&a2).Error; err != nil {
		t.Fatalf("create annotation 2: %v", err)
	}

	ids := []uuid.UUID{a1.ID, a2.ID}
	if err := ensureAnnotaionsBelongToOrg(orgA.ID, ids, db); err == nil {
		t.Fatalf("expected error when an annotation belongs to another org")
	}
}

func createTestSshKey(t *testing.T, db *gorm.DB, orgID uuid.UUID, name string) models.SshKey {
	t.Helper()

	key := models.SshKey{
		AuditFields: common.AuditFields{
			OrganizationID: orgID,
		},
		Name:                name,
		PublicKey:           "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAITestKey",
		EncryptedPrivateKey: "encrypted",
		PrivateIV:           "iv",
		PrivateTag:          "tag",
		Fingerprint:         "fp-" + name,
	}

	if err := db.Create(&key).Error; err != nil {
		t.Fatalf("create ssh key %s: %v", name, err)
	}

	return key
}
