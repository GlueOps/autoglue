package orgs

import (
	"encoding/json"
	"github.com/glueops/autoglue/api/middleware"
	"github.com/glueops/autoglue/internal/db"
	"github.com/glueops/autoglue/internal/db/models"
	"github.com/google/uuid"
	"net/http"
	"strings"
)

type OrgInput struct {
	Name string `json:"name"`
	Slug string `json:"slug"`
}

// CreateOrganization godoc
// @Summary      Create a new organization
// @Description  Creates a new organization and assigns the authenticated user as an admin member
// @Tags         Organizations
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

	json.NewEncoder(w).Encode(map[string]uuid.UUID{
		"organization_id": org.ID,
	})
}
