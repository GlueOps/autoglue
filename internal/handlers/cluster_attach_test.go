package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/glueops/autoglue/internal/api/httpmiddleware"
	"github.com/glueops/autoglue/internal/config"
	"github.com/glueops/autoglue/internal/models"
	"github.com/glueops/autoglue/internal/testutil/pgtest"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// The five cluster foreign keys managed by attach/detach handler pairs. Each is
// a distinct column on `clusters`, and Terraform drives several of them
// concurrently against the same row.
type fkAttachCase struct {
	name string
	// column is the real database column, which is not always what the Go field
	// name or the JSON tag suggests -- see glue_ops_load_balancer_id.
	column     string
	attach     func(*gorm.DB, config.Config) http.HandlerFunc
	detach     func(*gorm.DB, config.Config) http.HandlerFunc
	newTarget  func(*testing.T, *gorm.DB, uuid.UUID) uuid.UUID
	bodyKey    string
	targetCode string // error code when the target belongs to another org
}

func fkAttachCases() []fkAttachCase {
	return []fkAttachCase{
		{
			name:       "captain_domain",
			column:     "captain_domain_id",
			attach:     AttachCaptainDomain,
			detach:     DetachCaptainDomain,
			newTarget:  func(t *testing.T, db *gorm.DB, org uuid.UUID) uuid.UUID { return newTestDomain(t, db, org).ID },
			bodyKey:    "domain_id",
			targetCode: "domain_not_found",
		},
		{
			name:       "control_plane_record_set",
			column:     "control_plane_record_set_id",
			attach:     AttachControlPlaneRecordSet,
			detach:     DetachControlPlaneRecordSet,
			newTarget:  func(t *testing.T, db *gorm.DB, org uuid.UUID) uuid.UUID { return newTestRecordSet(t, db, org).ID },
			bodyKey:    "record_set_id",
			targetCode: "recordset_not_found",
		},
		{
			name:       "apps_load_balancer",
			column:     "apps_load_balancer_id",
			attach:     AttachAppsLoadBalancer,
			detach:     DetachAppsLoadBalancer,
			newTarget:  func(t *testing.T, db *gorm.DB, org uuid.UUID) uuid.UUID { return newTestLoadBalancer(t, db, org).ID },
			bodyKey:    "load_balancer_id",
			targetCode: "lb_not_found",
		},
		{
			name:       "glueops_load_balancer",
			column:     "glue_ops_load_balancer_id",
			attach:     AttachGlueOpsLoadBalancer,
			detach:     DetachGlueOpsLoadBalancer,
			newTarget:  func(t *testing.T, db *gorm.DB, org uuid.UUID) uuid.UUID { return newTestLoadBalancer(t, db, org).ID },
			bodyKey:    "load_balancer_id",
			targetCode: "lb_not_found",
		},
		{
			name:       "bastion_server",
			column:     "bastion_server_id",
			attach:     AttachBastionServer,
			detach:     DetachBastionServer,
			newTarget:  func(t *testing.T, db *gorm.DB, org uuid.UUID) uuid.UUID { return newTestServer(t, db, org).ID },
			bodyKey:    "server_id",
			targetCode: "server_not_found",
		},
	}
}

func TestClusterAttach_Authorization(t *testing.T) {
	db := pgtest.DB(t)
	cfg := config.Config{}

	for _, tc := range fkAttachCases() {
		t.Run(tc.name, func(t *testing.T) {
			org := createTestOrg(t, db, "attach-authz")
			other := createTestOrg(t, db, "attach-authz-other")
			cluster := newAttachCluster(t, db, org.ID)
			target := tc.newTarget(t, db, org.ID)
			body := fmt.Sprintf(`{%q:%q}`, tc.bodyKey, target.String())

			t.Run("403 without org", func(t *testing.T) {
				rr := httptest.NewRecorder()
				tc.attach(db, cfg).ServeHTTP(rr, clusterReq(http.MethodPost, body, nil, cluster.ID.String()))
				assertStatusCode(t, rr, http.StatusForbidden, "org_required")
			})

			t.Run("400 on malformed cluster id", func(t *testing.T) {
				rr := httptest.NewRecorder()
				tc.attach(db, cfg).ServeHTTP(rr, clusterReq(http.MethodPost, body, &org.ID, "not-a-uuid"))
				if rr.Code != http.StatusBadRequest {
					t.Fatalf("expected 400, got %d body=%s", rr.Code, rr.Body.String())
				}
			})

			t.Run("404 for a cluster in another org", func(t *testing.T) {
				foreign := newAttachCluster(t, db, other.ID)
				rr := httptest.NewRecorder()
				tc.attach(db, cfg).ServeHTTP(rr, clusterReq(http.MethodPost, body, &org.ID, foreign.ID.String()))
				// The code matters, not just the status: "not_found" is what
				// distinguishes "no such cluster for you" from a target lookup
				// failure, and it is why the pre-read before the write is kept.
				assertStatusCode(t, rr, http.StatusNotFound, "not_found")
			})

			t.Run("404 for a target in another org", func(t *testing.T) {
				foreignTarget := tc.newTarget(t, db, other.ID)
				rr := httptest.NewRecorder()
				b := fmt.Sprintf(`{%q:%q}`, tc.bodyKey, foreignTarget.String())
				tc.attach(db, cfg).ServeHTTP(rr, clusterReq(http.MethodPost, b, &org.ID, cluster.ID.String()))
				assertStatusCode(t, rr, http.StatusNotFound, tc.targetCode)
			})
		})
	}
}

func TestClusterAttach_WritesItsOwnColumn(t *testing.T) {
	db := pgtest.DB(t)
	cfg := config.Config{}

	for _, tc := range fkAttachCases() {
		t.Run(tc.name, func(t *testing.T) {
			org := createTestOrg(t, db, "attach-writes")
			cluster := newAttachCluster(t, db, org.ID)
			target := tc.newTarget(t, db, org.ID)

			rr := httptest.NewRecorder()
			body := fmt.Sprintf(`{%q:%q}`, tc.bodyKey, target.String())
			tc.attach(db, cfg).ServeHTTP(rr, clusterReq(http.MethodPost, body, &org.ID, cluster.ID.String()))
			if rr.Code != http.StatusOK {
				t.Fatalf("expected 200, got %d body=%s", rr.Code, rr.Body.String())
			}

			// Assert against the database, not the response body: the response
			// is rendered from an in-memory struct and can disagree with the row.
			if got := clusterColumn(t, db, cluster.ID, tc.column); got == nil || *got != target.String() {
				t.Fatalf("%s: expected column to be %s, got %v", tc.column, target, derefOrNil(got))
			}
			assertNeedsValidation(t, db, cluster.ID)
		})
	}
}

func TestClusterDetach_ClearsItsOwnColumn(t *testing.T) {
	db := pgtest.DB(t)
	cfg := config.Config{}

	for _, tc := range fkAttachCases() {
		t.Run(tc.name, func(t *testing.T) {
			org := createTestOrg(t, db, "detach-clears")
			cluster := newAttachCluster(t, db, org.ID)
			target := tc.newTarget(t, db, org.ID)

			// Seed the attachment directly. If the fixture did not actually
			// start attached, "the column is NULL" would pass for a handler
			// that does nothing at all.
			if err := db.Model(&models.Cluster{}).Where("id = ?", cluster.ID).
				Update(tc.column, target).Error; err != nil {
				t.Fatalf("seed %s: %v", tc.column, err)
			}
			if got := clusterColumn(t, db, cluster.ID, tc.column); got == nil || *got != target.String() {
				t.Fatalf("fixture did not attach %s (got %v) -- the assertion below would be vacuous",
					tc.column, derefOrNil(got))
			}

			rr := httptest.NewRecorder()
			tc.detach(db, cfg).ServeHTTP(rr, clusterReq(http.MethodDelete, "", &org.ID, cluster.ID.String()))
			if rr.Code != http.StatusOK {
				t.Fatalf("expected 200, got %d body=%s", rr.Code, rr.Body.String())
			}

			if got := clusterColumn(t, db, cluster.ID, tc.column); got != nil {
				t.Fatalf("%s: expected NULL after detach, got %q", tc.column, *got)
			}
			assertNeedsValidation(t, db, cluster.ID)
		})
	}
}

func TestUpdateCluster_ClusterProviderWritesProviderColumn(t *testing.T) {
	db := pgtest.DB(t)
	org := createTestOrg(t, db, "update-provider")
	cluster := newAttachCluster(t, db, org.ID)

	rr := httptest.NewRecorder()
	req := clusterReq(http.MethodPatch, `{"cluster_provider":"hetzner"}`, &org.ID, cluster.ID.String())
	UpdateCluster(db, config.Config{}).ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d body=%s", rr.Code, rr.Body.String())
	}

	// The request field is cluster_provider; the column is provider.
	var got models.Cluster
	if err := db.First(&got, "id = ?", cluster.ID).Error; err != nil {
		t.Fatalf("reload: %v", err)
	}
	if got.Provider != "hetzner" {
		t.Fatalf("expected provider=hetzner, got %q", got.Provider)
	}
	if got.Name != cluster.Name {
		t.Fatalf("name should be untouched, got %q want %q", got.Name, cluster.Name)
	}
}

func TestUpdateCluster_EmptyPatchPreservesEveryField(t *testing.T) {
	db := pgtest.DB(t)
	org := createTestOrg(t, db, "update-empty")
	cluster := newAttachCluster(t, db, org.ID)

	rr := httptest.NewRecorder()
	UpdateCluster(db, config.Config{}).ServeHTTP(rr,
		clusterReq(http.MethodPatch, `{}`, &org.ID, cluster.ID.String()))
	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d body=%s", rr.Code, rr.Body.String())
	}

	var got models.Cluster
	if err := db.First(&got, "id = ?", cluster.ID).Error; err != nil {
		t.Fatalf("reload: %v", err)
	}
	for _, f := range []struct{ name, got, want string }{
		{"name", got.Name, cluster.Name},
		{"provider", got.Provider, cluster.Provider},
		{"region", got.Region, cluster.Region},
		{"docker_image", got.DockerImage, cluster.DockerImage},
		{"docker_tag", got.DockerTag, cluster.DockerTag},
		{"random_token", got.RandomToken, cluster.RandomToken},
	} {
		if f.got != f.want {
			t.Errorf("%s: an empty PATCH wiped it: got %q want %q", f.name, f.got, f.want)
		}
	}
}

func TestClearClusterKubeconfig_EmptiesAllThreeColumns(t *testing.T) {
	db := pgtest.DB(t)
	org := createTestOrg(t, db, "clear-kubeconfig")
	cluster := newAttachCluster(t, db, org.ID)

	if err := db.Model(&models.Cluster{}).Where("id = ?", cluster.ID).
		Updates(map[string]any{
			"encrypted_kubeconfig": "ciphertext",
			"kube_iv":              "iv",
			"kube_tag":             "tag",
		}).Error; err != nil {
		t.Fatalf("seed kubeconfig: %v", err)
	}

	rr := httptest.NewRecorder()
	ClearClusterKubeconfig(db, config.Config{}).ServeHTTP(rr,
		clusterReq(http.MethodDelete, "", &org.ID, cluster.ID.String()))
	if rr.Code != http.StatusOK && rr.Code != http.StatusNoContent {
		t.Fatalf("expected 2xx, got %d body=%s", rr.Code, rr.Body.String())
	}

	// A revocation that returns success while leaving the ciphertext in place
	// is the failure this guards against.
	for _, col := range []string{"encrypted_kubeconfig", "kube_iv", "kube_tag"} {
		got := clusterColumn(t, db, cluster.ID, col)
		if got == nil || *got != "" {
			t.Errorf("%s: expected empty after clear, got %v", col, derefOrNil(got))
		}
	}
}

// --- helpers ---

// newAttachCluster creates a cluster in a state distinguishable from the
// post-handler state: a status that is not pre_pending and a non-empty
// last_error, so a handler that resets neither is visible.
func newAttachCluster(t *testing.T, db *gorm.DB, orgID uuid.UUID) models.Cluster {
	t.Helper()
	c := models.Cluster{
		OrganizationID: orgID,
		Name:           "cluster-" + uuid.NewString(),
		Provider:       "aws",
		Region:         "us-west-2",
		Status:         models.ClusterStatusReady,
		LastError:      "previous failure",
		DockerImage:    "ghcr.io/glueops/captain",
		DockerTag:      "v1.2.3",
		RandomToken:    uuid.NewString(),
	}
	if err := db.Create(&c).Error; err != nil {
		t.Fatalf("create cluster: %v", err)
	}
	return c
}

func newTestCredential(t *testing.T, db *gorm.DB, orgID uuid.UUID) models.Credential {
	t.Helper()
	c := models.Credential{
		OrganizationID: orgID,
		Provider:       "aws",
		Kind:           "static",
		ScopeKind:      "org",
		// Unique per credential: (org, provider, scope_kind, fingerprint) is a
		// unique index, and a test may need several credentials in one org.
		ScopeFingerprint: strings.ReplaceAll(uuid.NewString()+uuid.NewString(), "-", ""), // char(64)
		Name:             "cred-" + uuid.NewString(),
		EncryptedData:    "x",
		IV:               "x",
		Tag:              "x",
	}
	if err := db.Create(&c).Error; err != nil {
		t.Fatalf("create credential: %v", err)
	}
	return c
}

func newTestDomain(t *testing.T, db *gorm.DB, orgID uuid.UUID) models.Domain {
	t.Helper()
	d := models.Domain{
		OrganizationID: orgID,
		DomainName:     uuid.NewString()[:8] + ".example.com",
		CredentialID:   newTestCredential(t, db, orgID).ID,
	}
	if err := db.Create(&d).Error; err != nil {
		t.Fatalf("create domain: %v", err)
	}
	return d
}

func newTestRecordSet(t *testing.T, db *gorm.DB, orgID uuid.UUID) models.RecordSet {
	t.Helper()
	r := models.RecordSet{
		DomainID: newTestDomain(t, db, orgID).ID,
		Name:     "endpoint",
		Type:     "A",
		Values:   datatypes.JSON([]byte(`["1.2.3.4"]`)),
	}
	if err := db.Create(&r).Error; err != nil {
		t.Fatalf("create record set: %v", err)
	}
	return r
}

func newTestLoadBalancer(t *testing.T, db *gorm.DB, orgID uuid.UUID) models.LoadBalancer {
	t.Helper()
	lb := models.LoadBalancer{
		OrganizationID:   orgID,
		Name:             "lb-" + uuid.NewString(),
		Kind:             "network",
		PublicIPAddress:  "1.2.3.4",
		PrivateIPAddress: "10.0.0.4",
	}
	if err := db.Create(&lb).Error; err != nil {
		t.Fatalf("create load balancer: %v", err)
	}
	return lb
}

func newTestServer(t *testing.T, db *gorm.DB, orgID uuid.UUID) models.Server {
	t.Helper()
	key := models.SshKey{Name: "key-" + uuid.NewString()}
	key.OrganizationID = orgID
	if err := db.Create(&key).Error; err != nil {
		t.Fatalf("create ssh key: %v", err)
	}
	ip := "203.0.113.10"
	s := models.Server{
		OrganizationID:   orgID,
		Hostname:         "bastion",
		PublicIPAddress:  &ip, // required by Server.BeforeSave for role=bastion
		PrivateIPAddress: "10.0.0.10",
		SSHUser:          "cluster",
		SshKeyID:         key.ID,
		Role:             "bastion",
	}
	if err := db.Create(&s).Error; err != nil {
		t.Fatalf("create server: %v", err)
	}
	return s
}

// clusterReq builds a request carrying a clusterID route param, and an org in
// context only when orgID is non-nil.
func clusterReq(method, body string, orgID *uuid.UUID, clusterID string) *http.Request {
	r := httptest.NewRequest(method, "/clusters/"+clusterID, strings.NewReader(body))
	r.Header.Set("Content-Type", "application/json")

	ctx := r.Context()
	if orgID != nil {
		ctx = httpmiddleware.WithOrg(ctx, &models.Organization{ID: *orgID})
	}
	routeCtx := chi.NewRouteContext()
	routeCtx.URLParams.Add("clusterID", clusterID)
	ctx = context.WithValue(ctx, chi.RouteCtxKey, routeCtx)
	return r.WithContext(ctx)
}

// clusterColumn reads one column straight from the row. The column names are
// test constants, never user input.
func clusterColumn(t *testing.T, db *gorm.DB, clusterID uuid.UUID, column string) *string {
	t.Helper()
	var v *string
	q := fmt.Sprintf("SELECT %s::text FROM clusters WHERE id = ?", column)
	if err := db.Raw(q, clusterID).Scan(&v).Error; err != nil {
		t.Fatalf("read %s: %v", column, err)
	}
	return v
}

func assertNeedsValidation(t *testing.T, db *gorm.DB, clusterID uuid.UUID) {
	t.Helper()
	var c models.Cluster
	if err := db.First(&c, "id = ?", clusterID).Error; err != nil {
		t.Fatalf("reload cluster: %v", err)
	}
	if c.Status != models.ClusterStatusPrePending {
		t.Errorf("expected status %q, got %q", models.ClusterStatusPrePending, c.Status)
	}
	if c.LastError != "" {
		t.Errorf("expected last_error to be cleared, got %q", c.LastError)
	}
}

func assertStatusCode(t *testing.T, rr *httptest.ResponseRecorder, status int, code string) {
	t.Helper()
	if rr.Code != status {
		t.Fatalf("expected %d, got %d body=%s", status, rr.Code, rr.Body.String())
	}
	var body struct {
		Code string `json:"code"`
	}
	if err := json.Unmarshal(rr.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode error body %q: %v", rr.Body.String(), err)
	}
	if body.Code != code {
		t.Fatalf("expected error code %q, got %q", code, body.Code)
	}
}

func derefOrNil(s *string) any {
	if s == nil {
		return nil
	}
	return *s
}
