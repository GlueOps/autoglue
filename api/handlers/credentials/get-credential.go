package credentials

import (
	"errors"
	"net/http"

	"github.com/glueops/autoglue/api/middleware"
	"github.com/glueops/autoglue/internal/db"
	"github.com/glueops/autoglue/internal/db/models"
	"github.com/glueops/autoglue/internal/utils"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

// GetCredential godoc
// @Summary      Get credential by ID (org scoped)
// @Description  Redacted by default. Append `?reveal=true` to include decrypted value.
// @Tags         credentials
// @Accept       json
// @Produce      json
// @Param        X-Org-ID header string true "Organization UUID"
// @Param        id path string true "Credential ID (UUID)"
// @Param        reveal query bool false "Reveal decrypted secret (requires authorization)"
// @Security     BearerAuth
// @Success      200 {object} credentialResponse
// @Success      200 {object} credentialRevealResponse "When reveal=true"
// @Failure      400 {string} string "invalid id"
// @Failure      401 {string} string "Unauthorized"
// @Failure      403 {string} string "organization required"
// @Failure      404 {string} string "not found"
// @Failure      500 {string} string "failed to fetch/decrypt"
// @Router       /api/v1/credentials/{id} [get]
func GetCredential(w http.ResponseWriter, r *http.Request) {
	ac := middleware.GetAuthContext(r)
	if ac == nil || ac.OrganizationID == uuid.Nil {
		http.Error(w, "organization required", http.StatusForbidden)
		return
	}

	id, err := uuid.Parse(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	var cred models.Credential
	if err := db.DB.Where("id = ? AND organization_id = ?", id, ac.OrganizationID).First(&cred).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		http.Error(w, "failed to fetch", http.StatusInternalServerError)
		return
	}

	if r.URL.Query().Get("reveal") != "true" {
		writeJSON(w, http.StatusOK, credentialResponse{
			ID:             cred.ID,
			OrganizationID: cred.OrganizationID,
			Provider:       cred.Provider,
			CreatedAt:      cred.Timestamped.CreatedAt.UTC().Format(timeRFC3339),
			UpdatedAt:      cred.Timestamped.UpdatedAt.UTC().Format(timeRFC3339),
		})
		return
	}

	dec, err := utils.DecryptForOrg(ac.OrganizationID, cred.EncryptedData, cred.IV, cred.Tag)
	if err != nil {
		http.Error(w, "failed to decrypt", http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, credentialRevealResponse{
		credentialResponse: credentialResponse{
			ID:             cred.ID,
			OrganizationID: cred.OrganizationID,
			Provider:       cred.Provider,
			CreatedAt:      cred.Timestamped.CreatedAt.UTC().Format(timeRFC3339),
			UpdatedAt:      cred.Timestamped.UpdatedAt.UTC().Format(timeRFC3339),
		},
		Decrypted: dec,
	})
}
