package orgs

import (
	"encoding/json"
	"github.com/glueops/autoglue/api/middleware"
	"github.com/glueops/autoglue/internal/db"
	"github.com/glueops/autoglue/internal/db/models"
	"github.com/gorilla/mux"
	"net/http"
)

type UpdateOrgInput struct {
	Name string `json:"name"`
	Slug string `json:"slug"`
}

// UpdateOrganization godoc
// @Summary      Update organization metadata
// @Tags         Organizations
// @Accept       json
// @Produce      json
// @Param        orgId path string true "Org ID"
// @Param        body body UpdateOrgInput true "Organization data"
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

	orgId := mux.Vars(r)["orgId"]
	var input UpdateOrgInput
	json.NewDecoder(r.Body).Decode(&input)

	var org models.Organization
	db.DB.First(&org, "id = ?", orgId)

	org.Name = input.Name
	org.Slug = input.Slug
	db.DB.Save(&org)

	json.NewEncoder(w).Encode(org)
}
