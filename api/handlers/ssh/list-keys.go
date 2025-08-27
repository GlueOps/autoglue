package ssh

import (
	"net/http"
	"time"

	"github.com/glueops/autoglue/api/middleware"
	"github.com/glueops/autoglue/internal/db"
	"github.com/glueops/autoglue/internal/db/models"
	"github.com/google/uuid"
)

// ListPublicKeys godoc
// @Summary      List ssh keys (org scoped)
// @Description  Returns ssh keys for the organization in X-Org-ID.
// @Tags         ssh
// @Accept       json
// @Produce      json
// @Param        X-Org-ID header string true "Organization UUID"
// @Security     BearerAuth
// @Success      200 {array}  sshResponse
// @Failure      401 {string} string "Unauthorized"
// @Failure      403 {string} string "organization required"
// @Failure      500 {string} string "failed to list keys"
// @Router       /api/v1/ssh [get]
func ListPublicKeys(w http.ResponseWriter, r *http.Request) {
	ac := middleware.GetAuthContext(r)
	if ac == nil || ac.OrganizationID == uuid.Nil {
		http.Error(w, "organization required", http.StatusForbidden)
		return
	}

	var rows []models.SshKey
	if err := db.DB.Where("organization_id = ?", ac.OrganizationID).
		Order("created_at DESC").Find(&rows).Error; err != nil {
		http.Error(w, "failed to list ssh keys", http.StatusInternalServerError)
		return
	}

	out := make([]sshResponse, 0, len(rows))
	for _, s := range rows {
		out = append(out, sshResponse{
			ID:             s.ID,
			OrganizationID: s.OrganizationID,
			Name:           s.Name,
			PublicKey:      s.PublicKey,
			Fingerprint:    s.Fingerprint,
			CreatedAt:      s.CreatedAt.UTC().Format(time.RFC3339),
			UpdatedAt:      s.UpdatedAt.UTC().Format(time.RFC3339),
		})
	}

	writeJSON(w, http.StatusOK, out)
}
