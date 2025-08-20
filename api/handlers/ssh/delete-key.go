package ssh

import (
	"net/http"
	"strings"

	"github.com/glueops/autoglue/api/middleware"
	"github.com/glueops/autoglue/internal/db"
	"github.com/glueops/autoglue/internal/db/models"
	"github.com/google/uuid"
)

// DeleteSSHKey godoc
// @Summary      Delete ssh keypair (org scoped)
// @Description  Permanently deletes a keypair.
// @Tags         ssh
// @Accept       json
// @Produce      json
// @Param        X-Org-ID header string true "Organization UUID"
// @Param        id path string true "SSH Key ID (UUID)"
// @Security     BearerAuth
// @Success      204 {string} string "No Content"
// @Failure      400 {string} string "invalid id"
// @Failure      401 {string} string "Unauthorized"
// @Failure      403 {string} string "organization required"
// @Failure      500 {string} string "delete failed"
// @Router       /api/v1/ssh/{id} [delete]
func DeleteSSHKey(w http.ResponseWriter, r *http.Request) {
	ac := middleware.GetAuthContext(r)
	if ac == nil || ac.OrganizationID == uuid.Nil {
		http.Error(w, "organization required", http.StatusForbidden)
		return
	}

	idStr := strings.TrimPrefix(r.URL.Path, "/api/v1/ssh/")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	if err := db.DB.Where("id = ? AND organization_id = ?", id, ac.OrganizationID).
		Delete(&models.SshKey{}).Error; err != nil {
		http.Error(w, "delete failed", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
