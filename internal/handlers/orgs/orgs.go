package orgs

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/glueops/autoglue/internal/db"
	"github.com/glueops/autoglue/internal/db/models"
	"github.com/glueops/autoglue/internal/middleware"
	"github.com/glueops/autoglue/internal/response"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

// CreateOrganization godoc
// @Summary      Create a new organization
// @Description  Creates a new organization and assigns the authenticated user as an admin member
// @Tags         organizations
// @Accept       json
// @Produce      json
// @Param        X-Org-ID header string false "Optional organization context (ignored for creation)"
// @Param        body body OrgInput true "Organization Input"
// @Success      200 {object} map[string]string "organization_id"
// @Failure      400 {string} string "invalid input"
// @Failure      401 {string} string "unauthorized"
// @Failure      500 {string} string "internal error"
// @Security     BearerAuth
// @Router       /api/v1/orgs [post]
func CreateOrganization(w http.ResponseWriter, r *http.Request) {
	authCtx := middleware.GetAuthContext(r)
	if authCtx == nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	userID := authCtx.UserID

	var input OrgInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil || strings.TrimSpace(input.Name) == "" {
		http.Error(w, "invalid input", http.StatusBadRequest)
		return
	}

	org := &models.Organization{
		Name: input.Name,
		Slug: input.Slug,
	}

	if err := db.DB.Create(&org).Error; err != nil {
		http.Error(w, "could not create org", http.StatusInternalServerError)
		return
	}

	member := models.Member{
		UserID:         userID,
		OrganizationID: org.ID,
		Role:           "admin",
	}

	if err := db.DB.Create(&member).Error; err != nil {
		http.Error(w, "could not add member", http.StatusInternalServerError)
		return
	}

	_ = response.JSON(w, http.StatusCreated, org)
}

// ListOrganizations godoc
// @Summary      List organizations for user
// @Tags         organizations
// @Produce      json
// @Success      200 {array} models.Organization
// @Failure      401 {string} string "unauthorized"
// @Security     BearerAuth
// @Router       /api/v1/orgs [get]
func ListOrganizations(w http.ResponseWriter, r *http.Request) {
	auth := middleware.GetAuthContext(r)
	if auth == nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var orgs []models.Organization
	err := db.DB.Joins("JOIN members m ON m.organization_id = organizations.id").
		Where("m.user_id = ?", auth.UserID).Where("organizations.deleted_at IS NULL").Find(&orgs).Error
	if err != nil {
		http.Error(w, "failed to fetch orgs", http.StatusInternalServerError)
		return
	}

	_ = response.JSON(w, http.StatusOK, orgs)
}

// InviteMember godoc
// @Summary      Invite user to organization
// @Tags         organizations
// @Accept       json
// @Produce      plain
// @Param        body body InviteInput true "Invite input"
// @Success      201 {string} string "invited"
// @Failure      403 {string} string "forbidden"
// @Failure 400 {string} string "bad request"
// @Router       /api/v1/orgs/invite [post]
// @Param X-Org-ID header string true "Organization context"
// @Security     BearerAuth
func InviteMember(w http.ResponseWriter, r *http.Request) {
	auth := middleware.GetAuthContext(r)
	if auth == nil || auth.OrgRole != "admin" || auth.OrganizationID == uuid.Nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var input InviteInput
	json.NewDecoder(r.Body).Decode(&input)

	var user models.User
	err := db.DB.Where("email = ?", input.Email).First(&user).Error
	if err != nil {
		http.Error(w, "user not found", http.StatusBadRequest)
		return
	}

	invite := models.Invitation{
		OrganizationID: auth.OrganizationID,
		Email:          input.Email,
		Role:           input.Role,
		Status:         "pending",
		ExpiresAt:      time.Now().Add(48 * time.Hour),
		InviterID:      auth.UserID,
	}

	if err := db.DB.Create(&invite).Error; err != nil {
		http.Error(w, "failed to invite", http.StatusInternalServerError)
		return
	}

	_ = response.JSON(w, http.StatusCreated, invite)
}

// ListMembers lists all members of the authenticated org
// @Summary List organization members
// @Description Returns a list of all members in the current organization
// @Tags organizations
// @Security BearerAuth
// @Produce json
// @Success 200 {array} models.Member
// @Failure 401 {string} string "unauthorized"
// @Router /api/v1/orgs/members [get]
// @Param X-Org-ID header string true "Organization context"
func ListMembers(w http.ResponseWriter, r *http.Request) {
	authCtx := middleware.GetAuthContext(r)
	if authCtx == nil || authCtx.OrganizationID == uuid.Nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var members []models.Member
	if err := db.DB.Preload("User").Preload("Organization").Where("organization_id = ?", authCtx.OrganizationID).Find(&members).Error; err != nil {
		http.Error(w, "failed to fetch members", http.StatusInternalServerError)
		return
	}
	_ = response.JSON(w, http.StatusOK, members)
}

// DeleteMember godoc
// @Summary      Remove member from organization
// @Tags         organizations
// @Param        userId path string true "User ID"
// @Success      204 {string} string "deleted"
// @Failure      403 {string} string "forbidden"
// @Router       /api/v1/orgs/members/{userId} [delete]
// @Security     BearerAuth
func DeleteMember(w http.ResponseWriter, r *http.Request) {
	auth := middleware.GetAuthContext(r)
	if auth == nil || auth.OrgRole != "admin" {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}

	userId := chi.URLParam(r, "userId")

	if err := db.DB.Where("user_id = ? AND organization_id = ?", userId, auth.OrganizationID).Delete(&models.Member{}).Error; err != nil {
		http.Error(w, "failed to delete", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// UpdateOrganization godoc
// @Summary      Update organization metadata
// @Tags         organizations
// @Accept       json
// @Produce      json
// @Param        orgId path string true "Org ID"
// @Param        body body OrgInput true "Organization data"
// @Success      200 {object} models.Organization
// @Failure      403 {string} string "forbidden"
// @Router       /api/v1/orgs/{orgId} [patch]
// @Security     BearerAuth
func UpdateOrganization(w http.ResponseWriter, r *http.Request) {
	auth := middleware.GetAuthContext(r)
	if auth == nil || auth.OrgRole != "admin" {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}

	orgId := chi.URLParam(r, "orgId")

	var input OrgInput
	json.NewDecoder(r.Body).Decode(&input)

	var org models.Organization
	db.DB.First(&org, "id = ?", orgId)

	org.Name = input.Name
	org.Slug = input.Slug
	db.DB.Save(&org)

	_ = response.JSON(w, http.StatusOK, org)
}

// DeleteOrganization godoc
// @Summary      Delete organization
// @Tags         organizations
// @Param        orgId path string true "Organization ID"
// @Success      204 {string} string "deleted"
// @Failure      403 {string} string "forbidden"
// @Router       /api/v1/orgs/{orgId} [delete]
// @Security     BearerAuth
func DeleteOrganization(w http.ResponseWriter, r *http.Request) {
	auth := middleware.GetAuthContext(r)
	if auth == nil {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}

	orgId := chi.URLParam(r, "orgId")
	orgUUID, err := uuid.Parse(orgId)
	if err != nil {
		http.Error(w, "invalid organization id", http.StatusBadRequest)
		return
	}

	var member models.Member
	if err := db.DB.
		Where("user_id = ? AND organization_id = ?", auth.UserID, orgUUID).
		First(&member).Error; err != nil {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	if member.Role != "admin" {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}

	if err := db.DB.Where("organization_id = ?", orgUUID).Delete(&models.Member{}).Error; err != nil {
		http.Error(w, "failed to delete members", http.StatusInternalServerError)
		return
	}
	if err := db.DB.Delete(&models.Organization{}, "id = ?", orgUUID).Error; err != nil {
		http.Error(w, "failed to delete org", http.StatusInternalServerError)
		return
	}

	response.NoContent(w)
}
