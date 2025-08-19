package orgs

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/glueops/autoglue/api/middleware"
	"github.com/glueops/autoglue/internal/db"
	"github.com/glueops/autoglue/internal/db/models"
	"github.com/google/uuid"
)

type InviteInput struct {
	Email string `json:"email"`
	Role  string `json:"role"`
}

// InviteMember godoc
// @Summary      Invite user to organization
// @Tags         Organizations
// @Accept       json
// @Produce      plain
// @Param        body body InviteInput true "Invite input"
// @Success      201 {string} string "invited"
// @Failure      403 {string} string "forbidden"
// @Failure 400 {string} string "bad request"
// @Router       /api/v1/orgs/invite [post]
// @Param X-Org-ID header string true "Organization context"
// @Security     BearerAuth
func InviteMember(w http.ResponseWriter, r *http.Request) {
	auth := middleware.GetAuthContext(r)
	if auth == nil || auth.OrgRole != "admin" || auth.OrganizationID == uuid.Nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var input InviteInput
	json.NewDecoder(r.Body).Decode(&input)

	var user models.User
	err := db.DB.Where("email = ?", input.Email).First(&user).Error
	if err != nil {
		http.Error(w, "user not found", http.StatusBadRequest)
		return
	}

	invite := models.Invitation{
		OrganizationID: auth.OrganizationID,
		Email:          input.Email,
		Role:           input.Role,
		Status:         "pending",
		ExpiresAt:      time.Now().Add(48 * time.Hour),
		InviterID:      auth.UserID,
	}

	if err := db.DB.Create(&invite).Error; err != nil {
		http.Error(w, "failed to invite", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
