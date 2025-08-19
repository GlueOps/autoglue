package credentials

import (
	"encoding/json"
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

// UpdateCredential godoc
// @Summary      Update credential (org scoped)
// @Description  Patch provider and/or rotate secret by supplying plaintext.
// @Tags         credentials
// @Accept       json
// @Produce      json
// @Param        X-Org-ID header string true "Organization UUID"
// @Param        id path string true "Credential ID (UUID)"
// @Param        body body updateCredentialRequest true "Fields to update"
// @Security     BearerAuth
// @Success      200 {object} credentialResponse
// @Failure      400 {string} string "invalid id / invalid json"
// @Failure      401 {string} string "Unauthorized"
// @Failure      403 {string} string "organization required"
// @Failure      404 {string} string "not found"
// @Failure      500 {string} string "encryption/update failed"
// @Router       /api/v1/credentials/{id} [patch]
func UpdateCredential(w http.ResponseWriter, r *http.Request) {
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
		http.Error(w, "fetch failed", http.StatusInternalServerError)
		return
	}

	var req updateCredentialRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	if req.Provider != nil {
		cred.Provider = *req.Provider
	}
	if req.Plaintext != nil {
		data, iv, tag, err := utils.EncryptForOrg(ac.OrganizationID, []byte(*req.Plaintext))
		if err != nil {
			http.Error(w, "encryption failed", http.StatusInternalServerError)
			return
		}
		cred.EncryptedData = data
		cred.IV = iv
		cred.Tag = tag
	}

	if err := db.DB.Save(&cred).Error; err != nil {
		http.Error(w, "update failed", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, credentialResponse{
		ID:             cred.ID,
		OrganizationID: cred.OrganizationID,
		Provider:       cred.Provider,
		CreatedAt:      cred.Timestamped.CreatedAt.UTC().Format(timeRFC3339),
		UpdatedAt:      cred.Timestamped.UpdatedAt.UTC().Format(timeRFC3339),
	})
}
