package ssh

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/glueops/autoglue/api/middleware"
	"github.com/glueops/autoglue/internal/db"
	"github.com/glueops/autoglue/internal/db/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// GetSSHKey godoc
// @Summary      Get ssh key by ID (org scoped)
// @Description  Returns public key fields. Append `?reveal=true` to include the private key PEM.
// @Tags         ssh
// @Accept       json
// @Produce      json
// @Param        X-Org-ID header string true "Organization UUID"
// @Param        id path string true "SSH Key ID (UUID)"
// @Param        reveal query bool false "Reveal private key PEM"
// @Security     BearerAuth
// @Success      200 {object} sshResponse
// @Success      200 {object} sshRevealResponse "When reveal=true"
// @Failure      400 {string} string "invalid id"
// @Failure      401 {string} string "Unauthorized"
// @Failure      403 {string} string "organization required"
// @Failure      404 {string} string "not found"
// @Failure      500 {string} string "fetch failed"
// @Router       /api/v1/ssh/{id} [get]
func GetSSHKey(w http.ResponseWriter, r *http.Request) {
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

	var key models.SshKey
	if err := db.DB.Where("id = ? AND organization_id = ?", id, ac.OrganizationID).
		First(&key).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		http.Error(w, "fetch failed", http.StatusInternalServerError)
		return
	}

	if r.URL.Query().Get("reveal") != "true" {
		writeJSON(w, http.StatusOK, sshResponse{
			ID:             key.ID,
			OrganizationID: key.OrganizationID,
			PublicKey:      key.PublicKey,
			CreatedAt:      key.CreatedAt.UTC().Format(time.RFC3339),
			UpdatedAt:      key.UpdatedAt.UTC().Format(time.RFC3339),
		})
		return
	}

	writeJSON(w, http.StatusOK, sshRevealResponse{
		sshResponse: sshResponse{
			ID:             key.ID,
			OrganizationID: key.OrganizationID,
			PublicKey:      key.PublicKey,
			CreatedAt:      key.CreatedAt.UTC().Format(time.RFC3339),
			UpdatedAt:      key.UpdatedAt.UTC().Format(time.RFC3339),
		},
		PrivateKey: key.PrivateKey,
	})
}
