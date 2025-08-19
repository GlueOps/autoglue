package credentials

import (
	"encoding/json"
	"net/http"

	"github.com/glueops/autoglue/api/middleware"
	"github.com/glueops/autoglue/internal/db"
	"github.com/glueops/autoglue/internal/db/models"
	"github.com/glueops/autoglue/internal/utils"
	"github.com/google/uuid"
)

type CredentialInput struct {
	Provider      string `json:"provider" validate:"required"`
	EncryptedData string `json:"encrypted_data" validate:"required"`
	IV            string `json:"iv" validate:"required"`
	Tag           string `json:"tag" validate:"required"`
}

// CreateCredential godoc
// @Summary      Create credential (org scoped)
// @Description  Encrypts and stores plaintext for the org given by X-Org-ID.
// @Tags         credentials
// @Accept       json
// @Produce      json
// @Param        X-Org-ID header string true "Organization UUID"
// @Param        body body createCredentialRequest true "Credential payload"
// @Security     BearerAuth
// @Success      201 {object} credentialResponse
// @Failure      400 {string} string "invalid json / missing fields"
// @Failure      401 {string} string "Unauthorized"
// @Failure      403 {string} string "organization required"
// @Failure      500 {string} string "encryption/create failed"
// @Router       /api/v1/credentials [post]
func CreateCredential(w http.ResponseWriter, r *http.Request) {
	ac := middleware.GetAuthContext(r)
	if ac == nil || ac.OrganizationID == uuid.Nil {
		http.Error(w, "organization required", http.StatusForbidden)
		return
	}

	var req createCredentialRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}
	if req.Provider == "" || req.Plaintext == "" {
		http.Error(w, "provider and plaintext are required", http.StatusBadRequest)
		return
	}

	data, iv, tag, err := utils.EncryptForOrg(ac.OrganizationID, []byte(req.Plaintext))
	if err != nil {
		http.Error(w, "encryption failed", http.StatusInternalServerError)
		return
	}

	cred := models.Credential{
		OrganizationID: ac.OrganizationID,
		Provider:       req.Provider,
		EncryptedData:  data,
		IV:             iv,
		Tag:            tag,
	}
	if err := db.DB.Create(&cred).Error; err != nil {
		http.Error(w, "create failed", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusCreated, credentialResponse{
		ID:             cred.ID,
		OrganizationID: cred.OrganizationID,
		Provider:       cred.Provider,
		CreatedAt:      cred.Timestamped.CreatedAt.UTC().Format(timeRFC3339),
		UpdatedAt:      cred.Timestamped.UpdatedAt.UTC().Format(timeRFC3339),
	})
}
