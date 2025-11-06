package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/glueops/autoglue/internal/api/httpmiddleware"
	"github.com/glueops/autoglue/internal/models"
	"github.com/glueops/autoglue/internal/utils"
	"gorm.io/gorm"
)

type meResponse struct {
	models.User   `json:",inline"`
	Emails        []models.UserEmail    `json:"emails"`
	Organizations []models.Organization `json:"organizations"`
}

// GetMe godoc
//
//	@ID			GetMe
//	@Summary	Get current user profile
//	@Tags		Me
//	@Produce	json
//	@Success	200	{object}	meResponse
//	@Router		/me [get]
//	@Security	BearerAuth
//	@Security	ApiKeyAuth
func GetMe(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		u, ok := httpmiddleware.UserFrom(r.Context())
		if !ok {
			utils.WriteError(w, http.StatusUnauthorized, "unauthorized", "not signed in")
			return
		}

		var user models.User
		if err := db.First(&user, "id = ? AND is_disabled = false", u.ID).Error; err != nil {
			utils.WriteError(w, http.StatusUnauthorized, "unauthorized", "user not found/disabled")
			return
		}

		var emails []models.UserEmail
		_ = db.Preload("User").Where("user_id = ?", user.ID).Order("is_primary desc, created_at asc").Find(&emails).Error

		var orgs []models.Organization
		{
			var rows []models.Membership
			_ = db.Where("user_id = ?", user.ID).Find(&rows).Error
			if len(rows) > 0 {
				var ids []interface{}
				for _, m := range rows {
					ids = append(ids, m.OrganizationID)
				}
				_ = db.Find(&orgs, "id IN ?", ids).Error
			}
		}

		utils.WriteJSON(w, http.StatusOK, meResponse{
			User:          user,
			Emails:        emails,
			Organizations: orgs,
		})
	}
}

type updateMeRequest struct {
	DisplayName *string `json:"display_name,omitempty"`
	// You can add more editable fields here (timezone, avatar, etc)
}

// UpdateMe godoc
//
//	@ID			UpdateMe
//	@Summary	Update current user profile
//	@Tags		Me
//	@Accept		json
//	@Produce	json
//	@Param		body	body		updateMeRequest	true	"Patch profile"
//	@Success	200		{object}	models.User
//	@Router		/me [patch]
//	@Security	BearerAuth
//	@Security	ApiKeyAuth
func UpdateMe(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		u, ok := httpmiddleware.UserFrom(r.Context())
		if !ok {
			utils.WriteError(w, http.StatusUnauthorized, "unauthorized", "not signed in")
			return
		}

		var req updateMeRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			utils.WriteError(w, http.StatusBadRequest, "invalid_json", err.Error())
		}

		updates := map[string]interface{}{}

		if req.DisplayName != nil {
			updates["display_name"] = req.DisplayName
		}

		if len(updates) == 0 {
			var user models.User
			if err := db.First(&user, "id = ?", u.ID).Error; err != nil {
				utils.WriteError(w, 404, "not_found", "user")
				return
			}
			utils.WriteJSON(w, 200, user)
			return
		}

		if err := db.Model(&models.User{}).Where("id = ?", u.ID).Updates(updates).Error; err != nil {
			utils.WriteError(w, 500, "db_error", err.Error())
			return
		}

		var out models.User
		_ = db.First(&out, "id = ?", u.ID).Error
		utils.WriteJSON(w, 200, out)
	}
}
