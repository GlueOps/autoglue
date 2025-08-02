package orgs

import (
	"encoding/json"
	"github.com/glueops/autoglue/api/middleware"
	"github.com/glueops/autoglue/internal/db"
	"github.com/glueops/autoglue/internal/db/models"
	"net/http"
)

type SwitchOrgInput struct {
	OrganizationID string `json:"organization_id"`
}

// SwitchOrganization godoc
// @Summary      Switch active organization
// @Tags         Organizations
// @Accept       json
// @Produce      json
// @Param        body body SwitchOrgInput true "Org to switch to"
// @Success      200 {object} map[string]string "active org id"
// @Failure      403 {string} string "not a member"
// @Security     BearerAuth
// @Router       /api/v1/orgs/switch [post]
func SwitchOrganization(w http.ResponseWriter, r *http.Request) {
	auth := middleware.GetAuthContext(r)
	if auth == nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var input SwitchOrgInput
	json.NewDecoder(r.Body).Decode(&input)

	var member models.Member
	if err := db.DB.Where("user_id = ? AND organization_id = ?", auth.UserID, input.OrganizationID).First(&member).Error; err != nil {
		http.Error(w, "not a member", http.StatusForbidden)
		return
	}

	// Respond with confirmation
	json.NewEncoder(w).Encode(map[string]string{
		"active_org": input.OrganizationID,
	})
}
