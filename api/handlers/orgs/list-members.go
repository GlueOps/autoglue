package orgs

import (
	"encoding/json"
	"net/http"

	"github.com/glueops/autoglue/api/middleware"
	"github.com/glueops/autoglue/internal/db"
	"github.com/glueops/autoglue/internal/db/models"
	"github.com/google/uuid"
)

// ListMembers lists all members of the authenticated org
// @Summary List organization members
// @Description Returns a list of all members in the current organization
// @Tags Organizations
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

	json.NewEncoder(w).Encode(members)
}
