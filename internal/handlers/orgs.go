package handlers

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/glueops/autoglue/internal/api/httpmiddleware"
	"github.com/glueops/autoglue/internal/auth"
	"github.com/glueops/autoglue/internal/models"
	"github.com/glueops/autoglue/internal/utils"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ---------- Helpers ----------

func mustUser(r *http.Request) (*models.User, bool) {
	return httpmiddleware.UserFrom(r.Context())
}

func isOrgRole(db *gorm.DB, userID, orgID uuid.UUID, want ...string) (bool, string) {
	var m models.Membership
	if err := db.Where("user_id = ? AND organization_id = ?", userID, orgID).First(&m).Error; err != nil {
		return false, ""
	}
	got := strings.ToLower(m.Role)
	for _, w := range want {
		if got == strings.ToLower(w) {
			return true, got
		}
	}
	return false, got
}

func mustMember(db *gorm.DB, userID, orgID uuid.UUID) bool {
	ok, _ := isOrgRole(db, userID, orgID, "owner", "admin", "member")
	return ok
}

func randomB64URL(n int) (string, error) {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

// ---------- Orgs: list/create/get/update/delete ----------

type orgCreateReq struct {
	Name   string  `json:"name" example:"Acme Corp"`
	Domain *string `json:"domain,omitempty" example:"acme.com"`
}

// CreateOrg godoc
// @ID           CreateOrg
// @Summary      Create organization
// @Tags         Orgs
// @Accept       json
// @Produce      json
// @Param        body body orgCreateReq true "Org payload"
// @Success      201  {object} models.Organization
// @Failure      400  {object} utils.ErrorResponse
// @Failure      401  {object} utils.ErrorResponse
// @Failure      409  {object} utils.ErrorResponse
// @Router       /orgs [post]
// @ID           createOrg
// @Security     BearerAuth
func CreateOrg(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		u, ok := mustUser(r)
		if !ok {
			utils.WriteError(w, http.StatusUnauthorized, "unauthorized", "")
			return
		}

		var req orgCreateReq
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			utils.WriteError(w, 400, "invalid_json", err.Error())
			return
		}

		if strings.TrimSpace(req.Name) == "" {
			utils.WriteError(w, 400, "validation_error", "name is required")
			return
		}

		org := models.Organization{Name: req.Name}
		if req.Domain != nil && strings.TrimSpace(*req.Domain) != "" {
			org.Domain = req.Domain
		}

		if err := db.Create(&org).Error; err != nil {
			utils.WriteError(w, 409, "conflict", err.Error())
			return
		}

		// creator is owner
		_ = db.Create(&models.Membership{
			UserID: u.ID, OrganizationID: org.ID, Role: "owner",
		}).Error

		utils.WriteJSON(w, 201, org)
	}
}

// ListMyOrgs godoc
// @ID           ListMyOrgs
// @Summary      List organizations I belong to
// @Tags         Orgs
// @Produce      json
// @Success      200  {array} models.Organization
// @Failure      401  {object} utils.ErrorResponse
// @Router       /orgs [get]
// @ID           listMyOrgs
// @Security     BearerAuth
func ListMyOrgs(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		u, ok := mustUser(r)
		if !ok {
			utils.WriteError(w, http.StatusUnauthorized, "unauthorized", "")
			return
		}

		var orgs []models.Organization
		if err := db.
			Joins("join memberships m on m.organization_id = organizations.id").
			Where("m.user_id = ?", u.ID).
			Order("organizations.created_at desc").
			Find(&orgs).Error; err != nil {
			utils.WriteError(w, 500, "db_error", err.Error())
			return
		}
		utils.WriteJSON(w, 200, orgs)
	}
}

// GetOrg godoc
// @ID           GetOrg
// @Summary      Get organization
// @Tags         Orgs
// @Produce      json
// @Param        id   path string true "Org ID (UUID)"
// @Success      200  {object} models.Organization
// @Failure      401  {object} utils.ErrorResponse
// @Failure      404  {object} utils.ErrorResponse
// @Router       /orgs/{id} [get]
// @ID           getOrg
// @Security     BearerAuth
func GetOrg(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		u, ok := mustUser(r)
		if !ok {
			utils.WriteError(w, 401, "unauthorized", "")
			return
		}
		oid, err := uuid.Parse(chi.URLParam(r, "id"))
		if err != nil {
			utils.WriteError(w, 404, "not_found", "org not found")
			return
		}
		if !mustMember(db, u.ID, oid) {
			utils.WriteError(w, 401, "forbidden", "not a member")
			return
		}
		var org models.Organization
		if err := db.First(&org, "id = ?", oid).Error; err != nil {
			utils.WriteError(w, 404, "not_found", "org not found")
			return
		}
		utils.WriteJSON(w, 200, org)
	}
}

type orgUpdateReq struct {
	Name   *string `json:"name,omitempty"`
	Domain *string `json:"domain,omitempty"`
}

// UpdateOrg godoc
// @ID           UpdateOrg
// @Summary      Update organization (owner/admin)
// @Tags         Orgs
// @Accept       json
// @Produce      json
// @Param        id   path string        true "Org ID (UUID)"
// @Param        body body orgUpdateReq  true "Update payload"
// @Success      200  {object} models.Organization
// @Failure      401  {object} utils.ErrorResponse
// @Failure      404  {object} utils.ErrorResponse
// @Router       /orgs/{id} [patch]
// @ID           updateOrg
// @Security     BearerAuth
func UpdateOrg(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		u, ok := mustUser(r)
		if !ok {
			utils.WriteError(w, 401, "unauthorized", "")
			return
		}
		oid, err := uuid.Parse(chi.URLParam(r, "id"))
		if err != nil {
			utils.WriteError(w, 404, "not_found", "org not found")
			return
		}
		if ok, _ := isOrgRole(db, u.ID, oid, "owner", "admin"); !ok {
			utils.WriteError(w, 401, "forbidden", "admin or owner required")
			return
		}
		var req orgUpdateReq
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			utils.WriteError(w, 400, "invalid_json", err.Error())
			return
		}
		changes := map[string]any{}
		if req.Name != nil {
			changes["name"] = strings.TrimSpace(*req.Name)
		}
		if req.Domain != nil {
			if d := strings.TrimSpace(*req.Domain); d == "" {
				changes["domain"] = nil
			} else {
				changes["domain"] = d
			}
		}
		if len(changes) > 0 {
			if err := db.Model(&models.Organization{}).Where("id = ?", oid).Updates(changes).Error; err != nil {
				utils.WriteError(w, 500, "db_error", err.Error())
				return
			}
		}
		var out models.Organization
		_ = db.First(&out, "id = ?", oid).Error
		utils.WriteJSON(w, 200, out)
	}
}

// DeleteOrg godoc
// @ID           DeleteOrg
// @Summary      Delete organization (owner)
// @Tags         Orgs
// @Produce      json
// @Param        id   path string true "Org ID (UUID)"
// @Success      204  "Deleted"
// @Failure      401  {object} utils.ErrorResponse
// @Failure      404  {object} utils.ErrorResponse
// @Router       /orgs/{id} [delete]
// @ID           deleteOrg
// @Security     BearerAuth
func DeleteOrg(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		u, ok := mustUser(r)
		if !ok {
			utils.WriteError(w, 401, "unauthorized", "")
			return
		}
		oid, err := uuid.Parse(chi.URLParam(r, "id"))
		if err != nil {
			utils.WriteError(w, 404, "not_found", "org not found")
			return
		}
		if ok, _ := isOrgRole(db, u.ID, oid, "owner"); !ok {
			utils.WriteError(w, 401, "forbidden", "owner required")
			return
		}
		// Optional safety: deny if members >1 or resources exist; here we just delete.
		res := db.Delete(&models.Organization{}, "id = ?", oid)
		if res.Error != nil {
			utils.WriteError(w, 500, "db_error", res.Error.Error())
			return
		}
		if res.RowsAffected == 0 {
			utils.WriteError(w, 404, "not_found", "org not found")
			return
		}
		w.WriteHeader(204)
	}
}

// ---------- Members: list/add/update/delete ----------

type memberOut struct {
	UserID uuid.UUID `json:"user_id" format:"uuid"`
	Email  string    `json:"email"`
	Role   string    `json:"role"` // owner/admin/member
}

type memberUpsertReq struct {
	UserID uuid.UUID `json:"user_id" format:"uuid"`
	Role   string    `json:"role" example:"member"`
}

// ListMembers godoc
// @ID           ListMembers
// @Summary      List members in org
// @Tags         Orgs
// @Produce      json
// @Param        id   path string true "Org ID (UUID)"
// @Success      200  {array}  memberOut
// @Failure      401  {object} utils.ErrorResponse
// @Router       /orgs/{id}/members [get]
// @ID           listMembers
// @Security     BearerAuth
func ListMembers(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		u, ok := mustUser(r)
		if !ok {
			utils.WriteError(w, 401, "unauthorized", "")
			return
		}
		oid, err := uuid.Parse(chi.URLParam(r, "id"))
		if err != nil || !mustMember(db, u.ID, oid) {
			utils.WriteError(w, 401, "forbidden", "")
			return
		}
		var ms []models.Membership
		if err := db.Where("organization_id = ?", oid).Find(&ms).Error; err != nil {
			utils.WriteError(w, 500, "db_error", err.Error())
			return
		}

		// load emails
		userIDs := make([]uuid.UUID, 0, len(ms))
		for _, m := range ms {
			userIDs = append(userIDs, m.UserID)
		}
		var emails []models.UserEmail
		if len(userIDs) > 0 {
			_ = db.Where("user_id in ?", userIDs).Where("is_primary = true").Find(&emails).Error
		}
		emailByUser := map[uuid.UUID]string{}
		for _, e := range emails {
			emailByUser[e.UserID] = e.Email
		}

		out := make([]memberOut, 0, len(ms))
		for _, m := range ms {
			out = append(out, memberOut{
				UserID: m.UserID,
				Email:  emailByUser[m.UserID],
				Role:   m.Role,
			})
		}
		utils.WriteJSON(w, 200, out)
	}
}

// AddOrUpdateMember godoc
// @ID           AddOrUpdateMember
// @Summary      Add or update a member (owner/admin)
// @Tags         Orgs
// @Accept       json
// @Produce      json
// @Param        id   path string          true "Org ID (UUID)"
// @Param        body body memberUpsertReq true "User & role"
// @Success      200  {object} memberOut
// @Failure      401  {object} utils.ErrorResponse
// @Router       /orgs/{id}/members [post]
// @ID           addOrUpdateMember
// @Security     BearerAuth
func AddOrUpdateMember(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		u, ok := mustUser(r)
		if !ok {
			utils.WriteError(w, 401, "unauthorized", "")
			return
		}
		oid, err := uuid.Parse(chi.URLParam(r, "id"))
		if err != nil {
			utils.WriteError(w, 404, "not_found", "org not found")
			return
		}
		if ok, _ := isOrgRole(db, u.ID, oid, "owner", "admin"); !ok {
			utils.WriteError(w, 401, "forbidden", "admin or owner required")
			return
		}
		var req memberUpsertReq
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			utils.WriteError(w, 400, "invalid_json", err.Error())
			return
		}
		role := strings.ToLower(strings.TrimSpace(req.Role))
		if role != "owner" && role != "admin" && role != "member" {
			utils.WriteError(w, 400, "validation_error", "role must be owner|admin|member")
			return
		}
		var m models.Membership
		tx := db.Where("user_id = ? AND organization_id = ?", req.UserID, oid).First(&m)
		if tx.Error == nil {
			// update
			if err := db.Model(&m).Update("role", role).Error; err != nil {
				utils.WriteError(w, 500, "db_error", err.Error())
				return
			}
		} else if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			m = models.Membership{UserID: req.UserID, OrganizationID: oid, Role: role}
			if err := db.Create(&m).Error; err != nil {
				utils.WriteError(w, 500, "db_error", err.Error())
				return
			}
		} else {
			utils.WriteError(w, 500, "db_error", tx.Error.Error())
			return
		}

		// make response
		var ue models.UserEmail
		_ = db.Where("user_id = ? AND is_primary = true", req.UserID).First(&ue).Error
		utils.WriteJSON(w, 200, memberOut{
			UserID: req.UserID, Email: ue.Email, Role: m.Role,
		})
	}
}

// RemoveMember godoc
// @ID           RemoveMember
// @Summary      Remove a member (owner/admin)
// @Tags         Orgs
// @Produce      json
// @Param        id       path string true "Org ID (UUID)"
// @Param        user_id  path string true "User ID (UUID)"
// @Success      204 "Removed"
// @Failure      401 {object} utils.ErrorResponse
// @Router       /orgs/{id}/members/{user_id} [delete]
// @ID           removeMember
// @Security     BearerAuth
func RemoveMember(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		u, ok := mustUser(r)
		if !ok {
			utils.WriteError(w, 401, "unauthorized", "")
			return
		}
		oid, err := uuid.Parse(chi.URLParam(r, "id"))
		if err != nil {
			utils.WriteError(w, 404, "not_found", "org not found")
			return
		}
		if ok, _ := isOrgRole(db, u.ID, oid, "owner", "admin"); !ok {
			utils.WriteError(w, 401, "forbidden", "admin or owner required")
			return
		}
		uid, err := uuid.Parse(chi.URLParam(r, "user_id"))
		if err != nil {
			utils.WriteError(w, 400, "invalid_user_id", "")
			return
		}
		res := db.Where("user_id = ? AND organization_id = ?", uid, oid).Delete(&models.Membership{})
		if res.Error != nil {
			utils.WriteError(w, 500, "db_error", res.Error.Error())
			return
		}
		w.WriteHeader(204)
	}
}

// ---------- Org API Keys (key/secret pair) ----------

type orgKeyCreateReq struct {
	Name           string `json:"name,omitempty" example:"automation-bot"`
	ExpiresInHours *int   `json:"expires_in_hours,omitempty" example:"720"`
}

type orgKeyCreateResp struct {
	ID        uuid.UUID  `json:"id"`
	Name      string     `json:"name,omitempty"`
	Scope     string     `json:"scope"` // "org"
	CreatedAt time.Time  `json:"created_at"`
	ExpiresAt *time.Time `json:"expires_at,omitempty"`
	OrgKey    string     `json:"org_key"`    // shown once:
	OrgSecret string     `json:"org_secret"` // shown once:
}

// ListOrgKeys godoc
// @ID           ListOrgKeys
// @Summary      List org-scoped API keys (no secrets)
// @Tags         Orgs
// @Produce      json
// @Param        id   path string true "Org ID (UUID)"
// @Success      200  {array} models.APIKey
// @Failure      401  {object} utils.ErrorResponse
// @Router       /orgs/{id}/api-keys [get]
// @ID           listOrgKeys
// @Security     BearerAuth
func ListOrgKeys(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		u, ok := mustUser(r)
		if !ok {
			utils.WriteError(w, 401, "unauthorized", "")
			return
		}
		oid, err := uuid.Parse(chi.URLParam(r, "id"))
		if err != nil || !mustMember(db, u.ID, oid) {
			utils.WriteError(w, 401, "forbidden", "")
			return
		}
		var keys []models.APIKey
		if err := db.Where("org_id = ? AND scope = ?", oid, "org").
			Order("created_at desc").
			Find(&keys).Error; err != nil {
			utils.WriteError(w, 500, "db_error", err.Error())
			return
		}
		// SecretHash must not be exposed; your json tags likely hide it already.
		utils.WriteJSON(w, 200, keys)
	}
}

// CreateOrgKey godoc
// @ID           CreateOrgKey
// @Summary      Create org key/secret pair (owner/admin)
// @Tags         Orgs
// @Accept       json
// @Produce      json
// @Param        id   path string           true "Org ID (UUID)"
// @Param        body body orgKeyCreateReq  true "Key name + optional expiry"
// @Success      201  {object} orgKeyCreateResp
// @Failure      401  {object} utils.ErrorResponse
// @Router       /orgs/{id}/api-keys [post]
// @ID           createOrgKey
// @Security     BearerAuth
func CreateOrgKey(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		u, ok := mustUser(r)
		if !ok {
			utils.WriteError(w, 401, "unauthorized", "")
			return
		}
		oid, err := uuid.Parse(chi.URLParam(r, "id"))
		if err != nil {
			utils.WriteError(w, 404, "not_found", "org not found")
			return
		}
		if ok, _ := isOrgRole(db, u.ID, oid, "owner", "admin"); !ok {
			utils.WriteError(w, 401, "forbidden", "admin or owner required")
			return
		}

		var req orgKeyCreateReq
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			utils.WriteError(w, 400, "invalid_json", err.Error())
			return
		}

		// generate
		keySuffix, err := randomB64URL(16)
		if err != nil {
			utils.WriteError(w, 500, "entropy_error", err.Error())
			return
		}
		sec, err := randomB64URL(32)
		if err != nil {
			utils.WriteError(w, 500, "entropy_error", err.Error())
			return
		}
		orgKey := "org_" + keySuffix
		secretPlain := sec

		keyHash := auth.SHA256Hex(orgKey)
		secretHash, err := auth.HashSecretArgon2id(secretPlain)
		if err != nil {
			utils.WriteError(w, 500, "hash_error", err.Error())
			return
		}

		var exp *time.Time
		if req.ExpiresInHours != nil && *req.ExpiresInHours > 0 {
			e := time.Now().Add(time.Duration(*req.ExpiresInHours) * time.Hour)
			exp = &e
		}

		rec := models.APIKey{
			OrgID:      &oid,
			Scope:      "org",
			Name:       req.Name,
			KeyHash:    keyHash,
			SecretHash: &secretHash,
			ExpiresAt:  exp,
		}
		if err := db.Create(&rec).Error; err != nil {
			utils.WriteError(w, 500, "db_error", err.Error())
			return
		}

		utils.WriteJSON(w, 201, orgKeyCreateResp{
			ID:        rec.ID,
			Name:      rec.Name,
			Scope:     rec.Scope,
			CreatedAt: rec.CreatedAt,
			ExpiresAt: rec.ExpiresAt,
			OrgKey:    orgKey,
			OrgSecret: secretPlain,
		})
	}
}

// DeleteOrgKey godoc
// @ID           DeleteOrgKey
// @Summary      Delete org key (owner/admin)
// @Tags         Orgs
// @Produce      json
// @Param        id     path string true "Org ID (UUID)"
// @Param        key_id path string true "Key ID (UUID)"
// @Success      204 "Deleted"
// @Failure      401 {object} utils.ErrorResponse
// @Router       /orgs/{id}/api-keys/{key_id} [delete]
// @ID           deleteOrgKey
// @Security     BearerAuth
func DeleteOrgKey(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		u, ok := mustUser(r)
		if !ok {
			utils.WriteError(w, 401, "unauthorized", "")
			return
		}
		oid, err := uuid.Parse(chi.URLParam(r, "id"))
		if err != nil {
			utils.WriteError(w, 404, "not_found", "org not found")
			return
		}
		if ok, _ := isOrgRole(db, u.ID, oid, "owner", "admin"); !ok {
			utils.WriteError(w, 401, "forbidden", "admin or owner required")
			return
		}
		kid, err := uuid.Parse(chi.URLParam(r, "key_id"))
		if err != nil {
			utils.WriteError(w, 400, "invalid_key_id", "")
			return
		}
		res := db.Where("id = ? AND org_id = ? AND scope = ?", kid, oid, "org").Delete(&models.APIKey{})
		if res.Error != nil {
			utils.WriteError(w, 500, "db_error", res.Error.Error())
			return
		}
		if res.RowsAffected == 0 {
			utils.WriteError(w, 404, "not_found", "key not found")
			return
		}
		w.WriteHeader(204)
	}
}
