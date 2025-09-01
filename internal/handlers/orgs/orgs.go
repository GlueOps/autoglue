package orgs

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/glueops/autoglue/internal/db"
	"github.com/glueops/autoglue/internal/db/models"
	"github.com/glueops/autoglue/internal/middleware"
	"github.com/glueops/autoglue/internal/response"
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
