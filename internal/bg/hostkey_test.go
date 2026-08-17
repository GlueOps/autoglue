package bg

import (
	"crypto/ed25519"
	"encoding/base64"
	"net"
	"strings"
	"testing"

	"github.com/glueops/autoglue/internal/models"
	"github.com/glueops/autoglue/internal/testutil/pgtest"
	"github.com/google/uuid"
	"golang.org/x/crypto/ssh"
	"gorm.io/gorm"
)

func testHostKey(t *testing.T, seed byte) ssh.PublicKey {
	t.Helper()
	raw := make([]byte, ed25519.SeedSize)
	for i := range raw {
		raw[i] = seed
	}
	pub, err := ssh.NewPublicKey(ed25519.NewKeyFromSeed(raw).Public())
	if err != nil {
		t.Fatalf("build host key: %v", err)
	}
	return pub
}

func encodeKey(k ssh.PublicKey) string {
	return base64.StdEncoding.EncodeToString(k.Marshal())
}

func seedBastion(t *testing.T, db *gorm.DB) models.Server {
	t.Helper()

	org := models.Organization{Name: "hostkey-org-" + uuid.NewString()[:8]}
	if err := db.Create(&org).Error; err != nil {
		t.Fatalf("create org: %v", err)
	}
	key := models.SshKey{}
	key.OrganizationID = org.ID
	key.Name = "k"
	key.PublicKey = "ssh-ed25519 AAAA"
	key.EncryptedPrivateKey = "e"
	key.PrivateIV = "iv"
	key.PrivateTag = "tag"
	key.Fingerprint = "fp-" + uuid.NewString()[:8]
	if err := db.Create(&key).Error; err != nil {
		t.Fatalf("create ssh key: %v", err)
	}

	pub := "1.2.3.4"
	srv := models.Server{
		OrganizationID:   org.ID,
		Hostname:         "bastion",
		PublicIPAddress:  &pub,
		PrivateIPAddress: "10.0.0.1",
		SSHUser:          "deploy",
		SshKeyID:         key.ID,
		Role:             "bastion",
		Status:           "pending",
	}
	if err := db.Create(&srv).Error; err != nil {
		t.Fatalf("create server: %v", err)
	}
	return srv
}

func TestHostKeyCallbackLearnsOnFirstConnect(t *testing.T) {
	db := pgtest.DB(t)
	srv := seedBastion(t, db)
	key := testHostKey(t, 1)

	if err := makeDBHostKeyCallback(db, &srv)("host", &net.TCPAddr{}, key); err != nil {
		t.Fatalf("first connect should TOFU: %v", err)
	}

	var stored models.Server
	if err := db.Where("id = ?", srv.ID).First(&stored).Error; err != nil {
		t.Fatalf("reload: %v", err)
	}
	if stored.SSHHostKey != encodeKey(key) || stored.SSHHostKeyAlgo != key.Type() {
		t.Fatalf("host key not persisted: %q/%q", stored.SSHHostKeyAlgo, stored.SSHHostKey)
	}
}

func TestHostKeyCallbackRefusesMismatch(t *testing.T) {
	db := pgtest.DB(t)
	srv := seedBastion(t, db)
	known := testHostKey(t, 1)

	if err := makeDBHostKeyCallback(db, &srv)("host", &net.TCPAddr{}, known); err != nil {
		t.Fatalf("first connect: %v", err)
	}
	err := makeDBHostKeyCallback(db, &srv)("host", &net.TCPAddr{}, testHostKey(t, 2))
	if err == nil {
		t.Fatal("a different host key must be refused")
	}
	if !strings.Contains(err.Error(), "mismatch") {
		t.Fatalf("unhelpful refusal: %v", err)
	}
}

// The race this fixes: two first-connects run concurrently, one stores its key,
// and the other's guarded UPDATE matches zero rows. Before the fix that second
// caller returned nil — accepting a key it never compared against what had just
// been stored, which is the one thing TOFU exists to prevent.
//
// Simulated by giving the callback a stale in-memory server (host key empty)
// while the row already carries a different key, which is exactly the state the
// losing goroutine is in.
func TestHostKeyCallbackRefusesLostRaceWithDifferentKey(t *testing.T) {
	db := pgtest.DB(t)
	srv := seedBastion(t, db)

	// Another connection got there first and stored key 1.
	winner := testHostKey(t, 1)
	if err := db.Model(&models.Server{}).Where("id = ?", srv.ID).
		Updates(map[string]any{
			"ssh_host_key":      encodeKey(winner),
			"ssh_host_key_algo": winner.Type(),
		}).Error; err != nil {
		t.Fatalf("seed winner key: %v", err)
	}

	// This caller still believes the server has no key, and is presented a
	// different one.
	stale := srv
	stale.SSHHostKey = ""
	stale.SSHHostKeyAlgo = ""

	err := makeDBHostKeyCallback(db, &stale)("host", &net.TCPAddr{}, testHostKey(t, 2))
	if err == nil {
		t.Fatal("losing a first-connect race must not silently accept an unverified key")
	}
	if !strings.Contains(err.Error(), "mismatch") {
		t.Fatalf("refused for the wrong reason: %v", err)
	}

	// And the winner's key must survive untouched.
	var stored models.Server
	if err := db.Where("id = ?", srv.ID).First(&stored).Error; err != nil {
		t.Fatalf("reload: %v", err)
	}
	if stored.SSHHostKey != encodeKey(winner) {
		t.Fatal("the loser overwrote the stored host key")
	}
}

// Losing the race with the *same* key is benign: two connections to the same
// host saw the same key, and refusing there would fail a legitimate handshake.
func TestHostKeyCallbackAcceptsLostRaceWithSameKey(t *testing.T) {
	db := pgtest.DB(t)
	srv := seedBastion(t, db)
	key := testHostKey(t, 1)

	if err := db.Model(&models.Server{}).Where("id = ?", srv.ID).
		Updates(map[string]any{
			"ssh_host_key":      encodeKey(key),
			"ssh_host_key_algo": key.Type(),
		}).Error; err != nil {
		t.Fatalf("seed: %v", err)
	}

	stale := srv
	stale.SSHHostKey = ""
	stale.SSHHostKeyAlgo = ""

	if err := makeDBHostKeyCallback(db, &stale)("host", &net.TCPAddr{}, key); err != nil {
		t.Fatalf("same key on both connections should be accepted: %v", err)
	}
	// The caller's in-memory copy must be refreshed, or its next comparison
	// runs against an empty string and re-enters this path.
	if stale.SSHHostKey != encodeKey(key) {
		t.Fatal("in-memory host key not refreshed after losing the race")
	}
}
