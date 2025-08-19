package credentials

import (
	"net/http"

	"github.com/glueops/autoglue/api/middleware"
	"github.com/glueops/autoglue/internal/db"
	"github.com/glueops/autoglue/internal/db/models"
	"github.com/google/uuid"
)

// ListCredentials godoc
// @Summary      List credentials (org scoped)
// @Description  Returns redacted credentials for the organization in X-Org-ID.
// @Tags         credentials
// @Accept       json
// @Produce      json
// @Param        X-Org-ID header string true "Organization UUID"
// @Security     BearerAuth
// @Success      200 {array}  credentialResponse
// @Failure      401 {string} string "Unauthorized"
// @Failure      403 {string} string "organization required"
// @Failure      500 {string} string "failed to list credentials"
// @Router       /api/v1/credentials [get]
func ListCredentials(w http.ResponseWriter, r *http.Request) {
	ac := middleware.GetAuthContext(r)
	if ac == nil || ac.OrganizationID == uuid.Nil {
		http.Error(w, "organization required", http.StatusForbidden)
		return
	}

	var creds []models.Credential
	if err := db.DB.Where("organization_id = ?", ac.OrganizationID).
		Order("created_at DESC").Find(&creds).Error; err != nil {
		http.Error(w, "failed to list credentials", http.StatusInternalServerError)
		return
	}

	out := make([]credentialResponse, 0, len(creds))
	for _, c := range creds {
		out = append(out, credentialResponse{
			ID:             c.ID,
			OrganizationID: c.OrganizationID,
			Provider:       c.Provider,
			CreatedAt:      c.Timestamped.CreatedAt.UTC().Format(timeRFC3339),
			UpdatedAt:      c.Timestamped.UpdatedAt.UTC().Format(timeRFC3339),
		})
	}
	writeJSON(w, http.StatusOK, out)
}
