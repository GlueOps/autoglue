package httpmiddleware

import (
	"math"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/glueops/autoglue/internal/auth"
	"github.com/glueops/autoglue/internal/common"
	"github.com/glueops/autoglue/internal/models"
	"github.com/glueops/autoglue/internal/testutil/pgtest"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// seedAgent builds a whole org/server/cluster/agent chain, because
// uniq_agents_live_cluster allows one live agent per cluster and every caller
// here wants its own.
func seedAgent(t *testing.T, db *gorm.DB) (models.Agent, models.Organization, auth.AgentCredential) {
	t.Helper()

	org := models.Organization{Name: "agentauth-" + uuid.NewString()[:8]}
	if err := db.Create(&org).Error; err != nil {
		t.Fatalf("org: %v", err)
	}
	key := models.SshKey{
		AuditFields:         common.AuditFields{OrganizationID: org.ID},
		Name:                "k",
		PublicKey:           "p",
		EncryptedPrivateKey: "e",
		PrivateIV:           "iv",
		PrivateTag:          "tag",
		Fingerprint:         "fp-" + uuid.NewString()[:8],
	}
	if err := db.Create(&key).Error; err != nil {
		t.Fatalf("sshkey: %v", err)
	}
	pub := "1.2.3.4"
	srv := models.Server{
		OrganizationID:   org.ID,
		PublicIPAddress:  &pub,
		PrivateIPAddress: "10.0.0.1",
		SSHUser:          "root",
		SshKeyID:         key.ID,
		Role:             "bastion",
	}
	if err := db.Create(&srv).Error; err != nil {
		t.Fatalf("server: %v", err)
	}
	cl := models.Cluster{OrganizationID: org.ID, Name: "c-" + uuid.NewString()[:8]}
	if err := db.Create(&cl).Error; err != nil {
		t.Fatalf("cluster: %v", err)
	}

	cred, err := auth.MintAgentCredential()
	if err != nil {
		t.Fatalf("mint: %v", err)
	}
	a := models.Agent{
		ID:             cred.ID,
		OrganizationID: org.ID,
		ClusterID:      cl.ID,
		ServerID:       srv.ID,
		KeyHash:        cred.KeyHash,
		SecretHash:     cred.SecretHash,
		Prefix:         cred.Prefix,
		Status:         models.AgentStatusActive,
	}
	if err := db.Create(&a).Error; err != nil {
		t.Fatalf("agent: %v", err)
	}
	return a, org, cred
}

// TestAgentAuthAcceptsOnlyTheWholeTuple walks every way a caller can hold part
// of a credential.
//
// What breaks in production if this fails: whichever case starts returning 204
// is a way into the control plane without the full tuple. The two that would
// be silent for months are `revoked` and `expired` — ValidateOrgKeyPair, which
// this was modelled on, checks neither, and gets away with it only because a
// sweeper deletes its rows. An agent row is never deleted (re-enrolment revokes
// it, precisely so its dead-lettered tasks outlive the credential), so a
// missing status predicate leaves every bastion ever decommissioned still able
// to authenticate and to take work for its old cluster.
func TestAgentAuthAcceptsOnlyTheWholeTuple(t *testing.T) {
	db := pgtest.DB(t)
	a, org, cred := seedAgent(t, db)

	var sawAgent *models.Agent
	var sawOrgID uuid.UUID
	var sawOrg, sawUser, sawRoles bool
	h := AgentAuth(db)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sawAgent, _ = AgentFrom(r.Context())
		sawOrgID, sawOrg = OrgIDFrom(r.Context())
		_, sawUser = UserFrom(r.Context())
		_, sawRoles = RolesFrom(r.Context())
		w.WriteHeader(http.StatusNoContent)
	}))

	call := func(id, key, secret string) int {
		sawAgent, sawOrg, sawUser, sawRoles = nil, false, false, false
		req := httptest.NewRequest(http.MethodGet, "/api/v1/agent/sync", nil)
		req.Header.Set("X-Agent-ID", id)
		req.Header.Set("X-Agent-KEY", key)
		req.Header.Set("X-Agent-SECRET", secret)
		rr := httptest.NewRecorder()
		h.ServeHTTP(rr, req)
		return rr.Code
	}

	if got := call(a.ID.String(), cred.Key, cred.Secret); got != http.StatusNoContent {
		t.Fatalf("happy path = %d, want 204", got)
	}
	if sawAgent == nil || sawAgent.ID != a.ID || sawAgent.ClusterID != a.ClusterID {
		t.Fatalf("agent missing from context: %+v", sawAgent)
	}
	if !sawOrg || sawOrgID != org.ID {
		t.Fatalf("org = %v/%v, want %v", sawOrg, sawOrgID, org.ID)
	}
	// An agent that also presented as a user principal would satisfy
	// RequireAuthenticatedUser and RequireRole, so a route composed with either
	// by mistake would admit a bastion to a human endpoint. Absence here is
	// what makes that mistake fail closed instead of silently working.
	if sawUser || sawRoles {
		t.Fatal("an agent must not present as a user principal")
	}

	for name, code := range map[string]int{
		"wrong secret":   call(a.ID.String(), cred.Key, "nope"),
		"unknown key":    call(a.ID.String(), "agt_nope", cred.Secret),
		"mismatched id":  call(uuid.New().String(), cred.Key, cred.Secret),
		"missing id":     call("", cred.Key, cred.Secret),
		"unparseable id": call("not-a-uuid", cred.Key, cred.Secret),
		"missing secret": call(a.ID.String(), cred.Key, ""),
		// The prefix is displayed in UIs and logs on purpose. If it were
		// accepted as the key, publishing it would be publishing half the
		// credential.
		"prefix as key": call(a.ID.String(), cred.Prefix, cred.Secret),
	} {
		if code != http.StatusUnauthorized {
			t.Errorf("%s = %d, want 401", name, code)
		}
	}

	if err := db.Model(&models.Agent{}).Where("id = ?", a.ID).
		Update("status", models.AgentStatusRevoked).Error; err != nil {
		t.Fatal(err)
	}
	if got := call(a.ID.String(), cred.Key, cred.Secret); got != http.StatusUnauthorized {
		t.Fatalf("revoked = %d, want 401 — a decommissioned bastion still authenticates", got)
	}

	past := time.Now().Add(-time.Hour)
	if err := db.Model(&models.Agent{}).Where("id = ?", a.ID).
		Updates(map[string]any{"status": models.AgentStatusActive, "expires_at": past}).Error; err != nil {
		t.Fatal(err)
	}
	if got := call(a.ID.String(), cred.Key, cred.Secret); got != http.StatusUnauthorized {
		t.Fatalf("expired = %d, want 401", got)
	}
}

// TestAgentAuthIgnoresCallerSuppliedOrg pins the scope to the credential.
//
// What breaks in production if this fails: the whole reason this middleware is
// separate from AuthMiddleware. That one resolves an org from X-Org-ID and from
// the URL, all of it user-membership reasoning with no agent analogue. If any
// of it is ever folded in here, a leaked bastion credential stops being worth
// one cluster and becomes worth an organization — which is exactly the org-wide
// key this design exists to retire. Nothing about that is visible in a
// response: the request succeeds either way, just against the wrong tenant.
func TestAgentAuthIgnoresCallerSuppliedOrg(t *testing.T) {
	db := pgtest.DB(t)
	a, org, cred := seedAgent(t, db)
	_, victim, _ := seedAgent(t, db)

	var sawOrgID uuid.UUID
	h := AgentAuth(db)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sawOrgID, _ = OrgIDFrom(r.Context())
		w.WriteHeader(http.StatusNoContent)
	}))

	req := httptest.NewRequest(http.MethodGet, "/api/v1/agent/sync", nil)
	req.Header.Set("X-Agent-ID", a.ID.String())
	req.Header.Set("X-Agent-KEY", cred.Key)
	req.Header.Set("X-Agent-SECRET", cred.Secret)
	req.Header.Set("X-Org-ID", victim.ID.String())
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)

	if rr.Code != http.StatusNoContent {
		t.Fatalf("status = %d %s", rr.Code, rr.Body.String())
	}
	if sawOrgID != org.ID {
		t.Fatalf("org = %v, want the credential's %v — a header re-scoped an agent", sawOrgID, org.ID)
	}
}

// TestAgentAuthTouchesLastSeenWithoutTouchingUpdatedAt.
//
// What breaks in production if this fails: two things, both quiet. If
// last_seen_at stops being written, an agent that is polling happily looks dead
// and whatever staleness alerting gets built on it fires on healthy fleets. If
// updated_at starts moving instead, it stops meaning "the control plane changed
// this row" and starts meaning "a bastion polled thirty seconds ago" — so every
// agent looks freshly modified forever and the column is useless for spotting
// an actual credential change.
func TestAgentAuthTouchesLastSeenWithoutTouchingUpdatedAt(t *testing.T) {
	db := pgtest.DB(t)
	a, _, cred := seedAgent(t, db)

	h := AgentAuth(db)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))
	req := httptest.NewRequest(http.MethodGet, "/api/v1/agent/assignment", nil)
	req.Header.Set("X-Agent-ID", a.ID.String())
	req.Header.Set("X-Agent-KEY", cred.Key)
	req.Header.Set("X-Agent-SECRET", cred.Secret)
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)
	if rr.Code != http.StatusNoContent {
		t.Fatalf("status = %d", rr.Code)
	}

	var fresh models.Agent
	if err := db.First(&fresh, "id = ?", a.ID).Error; err != nil {
		t.Fatal(err)
	}
	if fresh.LastSeenAt == nil {
		t.Fatal("last_seen_at not written — a live agent will look dead")
	}
	if !fresh.UpdatedAt.Equal(a.UpdatedAt) {
		t.Fatalf("updated_at moved on a poll: %v -> %v", a.UpdatedAt, fresh.UpdatedAt)
	}
}

// TestAgentSecretRejectionIsNotShortCircuited asserts the *shape* of a
// rejection, not just its verdict.
//
// What breaks in production if this fails: the secret becomes guessable. The
// verdict is identical either way — nil is returned, the request 401s, every
// other test in this file still passes — so the only observable difference is
// how long the answer took. A `len(secret) != len(stored)` guard before the
// comparison, or an `==` on the encoded hashes (a length-prefixed string
// compare that leaks a byte at a time), would both return orders of magnitude
// faster than a real verification and hand a remote attacker an oracle.
//
// The assertion is a ratio against a genuine verification rather than an
// absolute duration, so it means the same thing on a loaded CI box as on a
// laptop: whatever argon2id costs here, a wrong secret of any length must pay
// substantially the same. A short circuit is not a near miss on this bound —
// it is nanoseconds against tens of milliseconds — so the factor is left loose
// on purpose. `min` rather than a mean because scheduler noise only ever adds
// time, and the fastest observed run is the one closest to the real cost.
func TestAgentSecretRejectionIsNotShortCircuited(t *testing.T) {
	db := pgtest.DB(t)
	a, _, cred := seedAgent(t, db)

	fastest := func(secret string, wantValid bool) time.Duration {
		best := time.Duration(math.MaxInt64)
		for i := 0; i < 3; i++ {
			start := time.Now()
			got := auth.ValidateAgentKeyPair(a.ID.String(), cred.Key, secret, db)
			if d := time.Since(start); d < best {
				best = d
			}
			if (got != nil) != wantValid {
				t.Fatalf("secret %q: valid = %v, want %v", secret, got != nil, wantValid)
			}
		}
		return best
	}

	sameLen := strings.Repeat("A", len(cred.Secret))
	if sameLen == cred.Secret {
		t.Fatal("degenerate fixture: the wrong secret is the right one")
	}

	// The baseline is a full, successful verification: the row lookup plus the
	// argon2id KDF. Everything below has to cost about the same.
	baseline := fastest(cred.Secret, true)

	for name, wrong := range map[string]string{
		"same length as the real secret": sameLen,
		"one character":                  "x",
		"far longer than the real one":   strings.Repeat("z", 4096),
	} {
		got := fastest(wrong, false)
		if got*5 < baseline {
			t.Errorf("wrong secret (%s) rejected in %v against a %v verification: "+
				"something short-circuits before the constant-time comparison",
				name, got, baseline)
		}
	}
}
