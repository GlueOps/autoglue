package orgs

import (
	"github.com/glueops/autoglue/api/middleware"
	"github.com/glueops/autoglue/internal/db"
	"github.com/glueops/autoglue/internal/db/models"
	"github.com/gorilla/mux"
	"net/http"
)

// DeleteOrganization godoc
// @Summary      Delete organization
// @Tags         Organizations
// @Param        orgId path string true "Organization ID"
// @Success      204 {string} string "deleted"
// @Failure      403 {string} string "forbidden"
// @Router       /api/v1/orgs/{orgId} [delete]
// @Security     BearerAuth
func DeleteOrganization(w http.ResponseWriter, r *http.Request) {
	auth := middleware.GetAuthContext(r)
	if auth == nil || auth.OrgRole != "admin" {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}

	orgId := mux.Vars(r)["orgId"]
	db.DB.Delete(&models.Organization{}, "id = ?", orgId)
	w.WriteHeader(http.StatusNoContent)
}
