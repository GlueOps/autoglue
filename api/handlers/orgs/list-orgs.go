package orgs

import (
	"encoding/json"
	"net/http"

	"github.com/glueops/autoglue/api/middleware"
	"github.com/glueops/autoglue/internal/db"
	"github.com/glueops/autoglue/internal/db/models"
)

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

	_ = json.NewEncoder(w).Encode(orgs)
}
