package orgs

import (
	"encoding/json"
	"github.com/glueops/autoglue/api/middleware"
	"github.com/glueops/autoglue/internal/db"
	"github.com/glueops/autoglue/internal/db/models"
	"net/http"
	"time"
)

type JoinInput struct {
	InvitationID string `json:"invitation_id"`
}

// JoinOrganization godoc
// @Summary      Accept org invitation
// @Tags         Organizations
// @Accept       json
// @Produce      json
// @Param        body body JoinInput true "Join input"
// @Success      200 {string} string "joined"
// @Failure      403 {string} string "invalid invite"
// @Router       /api/v1/orgs/join [post]
// @Security     BearerAuth
func JoinOrganization(w http.ResponseWriter, r *http.Request) {
	auth := middleware.GetAuthContext(r)
	if auth == nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var input JoinInput
	json.NewDecoder(r.Body).Decode(&input)

	var user models.User
	db.DB.First(&user, auth.UserID)

	var invite models.Invitation
	if err := db.DB.Where("id = ? AND email = ?", input.InvitationID, user.Email).First(&invite).Error; err != nil {
		http.Error(w, "invalid invitation", http.StatusForbidden)
		return
	}

	if invite.Status != "pending" || invite.ExpiresAt.Before(time.Now()) {
		http.Error(w, "invitation expired or used", http.StatusForbidden)
		return
	}

	db.DB.Create(&models.Member{
		UserID:         auth.UserID,
		OrganizationID: invite.OrganizationID,
		Role:           invite.Role,
	})

	db.DB.Model(&invite).Update("status", "accepted")
	w.WriteHeader(http.StatusOK)
}
