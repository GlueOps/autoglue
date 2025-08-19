package orgs

import (
	"net/http"

	"github.com/glueops/autoglue/api/middleware"
	"github.com/glueops/autoglue/internal/db"
	"github.com/glueops/autoglue/internal/db/models"
	"github.com/gorilla/mux"
)

// DeleteMember godoc
// @Summary      Remove member from organization
// @Tags         Organizations
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

	userId := mux.Vars(r)["userId"]

	if err := db.DB.Where("user_id = ? AND organization_id = ?", userId, auth.OrganizationID).Delete(&models.Member{}).Error; err != nil {
		http.Error(w, "failed to delete", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
